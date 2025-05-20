package synchronize

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syncr/models"
	"time"
)

func CompareFileData(a, b []models.FileData) []models.SyncAction {
	mapA := make(map[string]models.FileData)
	mapB := make(map[string]models.FileData)
	var actions []models.SyncAction

	for _, f := range a {
		mapA[f.Name] = f
	}
	for _, f := range b {
		mapB[f.Name] = f
	}

	seen := make(map[string]bool)

	for name, fileA := range mapA {
		seen[name] = true
		if fileB, ok := mapB[name]; ok {
			if fileA.Checksum != fileB.Checksum ||
				fileA.Size != fileB.Size ||
				!fileA.ModTime.Equal(fileB.ModTime) ||
				fileA.Permissions != fileB.Permissions {
				actions = append(actions, models.SyncAction{
					Type:   models.Modify,
					Source: fileA,
					Target: &fileB,
				})
			}
		} else {
			actions = append(actions, models.SyncAction{
				Type:   models.Add,
				Source: fileA,
				Target: nil,
			})
		}
	}

	for name, fileB := range mapB {
		if !seen[name] {
			actions = append(actions, models.SyncAction{
				Type:   models.Missing,
				Source: fileB,
				Target: nil,
			})
		}
	}

	return actions
}

func IsSyncRequired(deleteMissing bool, actions []models.SyncAction) bool {
	for _, action := range actions {
		if action.Type == models.Add || action.Type == models.Modify {
			return true
		}
		if action.Type == models.Missing && deleteMissing {
			return true
		}
	}
	return false
}

func ExplainSyncActions(actions []models.SyncAction) {
	fmt.Println("=== Differences ===")
	for _, action := range actions {
		switch action.Type {
		case models.Add:
			fmt.Printf("Add: %s\n", action.Source.Name)
		case models.Modify:
			fmt.Printf("Modify: %s\n", action.Source.Name)

			if action.Source.Checksum != action.Target.Checksum {
				fmt.Printf("  - Old checksum: %s\n  - New checksum: %s\n", action.Source.Checksum, action.Target.Checksum)
			}
			if action.Source.Size != action.Target.Size {
				fmt.Printf("  - Size changed: %d → %d\n", action.Source.Size, action.Target.Size)
			}
			if !action.Source.ModTime.Equal(action.Target.ModTime) {
				fmt.Printf("  - Modified time changed: %v → %v\n", action.Source.ModTime, action.Target.ModTime)
			}
			if action.Source.Permissions != action.Target.Permissions {
				fmt.Printf("  - Permissions changed: %s → %s\n", action.Source.Permissions, action.Target.Permissions)
			}
		case models.Missing:
			fmt.Printf("Missing: %s\n", action.Source.Name)
		}
	}
}

type SyncDone struct{}

func SyncFiles(actions []models.SyncAction, sourceRoot, targetRoot string, deleteMissing bool) {
	var wg sync.WaitGroup
	actionChan := make(chan models.SyncAction)
	syncDoneChan := make(chan SyncDone)
	var completed int64
	total := int64(len(actions))
	workerCount := optimalWorkerCount()

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				done := atomic.LoadInt64(&completed)
				percent := float64(done) / float64(total) * 100
				fmt.Printf("\rProgress: %.2f%% (%d/%d)", percent, done, total)
			case <-syncDoneChan:
				fmt.Printf("\rProgress: 100.00%% (%d/%d)\n", total, total)
				return
			}
		}
	}()

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for action := range actionChan {
				switch action.Type {
				case models.Add, models.Modify:
					srcPath := filepath.Join(sourceRoot, action.Source.Name)
					dstPath := filepath.Join(targetRoot, action.Source.Name)

					if err := copyFileWithModTimeAndPermissions(srcPath, dstPath); err != nil {
						log.Printf("\nERROR copying %s: %v\n", srcPath, err)
					}
				case models.Missing:
					if deleteMissing {
						targetPath := filepath.Join(targetRoot, action.Source.Name)
						if err := os.Remove(targetPath); err != nil {
							log.Printf("\nERROR deleting %s: %v\n", targetPath, err)
						}
					}
				}
				atomic.AddInt64(&completed, 1)
			}
		}()
	}

	for _, a := range actions {
		actionChan <- a
	}
	close(actionChan)
	close(syncDoneChan)
	wg.Wait()
	return
}

func copyFileWithModTimeAndPermissions(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	sourceInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := dstFile.Chmod(sourceInfo.Mode()); err != nil {
		return err
	}
	return os.Chtimes(dst, sourceInfo.ModTime(), sourceInfo.ModTime())
}

func optimalWorkerCount() int {
	n := runtime.NumCPU()
	if n < 2 {
		return 2
	} else if n > 8 {
		return 8
	}
	return n
}
