package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type SyncActionType string

const (
	Add     SyncActionType = "Add"
	Modify  SyncActionType = "Modify"
	Missing SyncActionType = "Missing"
)

type SyncAction struct {
	Type   SyncActionType
	Source FileData
	Target *FileData
}

func CompareFileData(a, b []FileData) []SyncAction {
	mapA := make(map[string]FileData)
	mapB := make(map[string]FileData)
	var actions []SyncAction

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
				actions = append(actions, SyncAction{
					Type:   Modify,
					Source: fileA,
					Target: &fileB,
				})
			}
		} else {
			actions = append(actions, SyncAction{
				Type:   Add,
				Source: fileA,
				Target: nil,
			})
		}
	}

	for name, fileB := range mapB {
		if !seen[name] {
			actions = append(actions, SyncAction{
				Type:   Missing,
				Source: fileB,
				Target: nil,
			})
		}
	}

	return actions
}

func IsSyncRequired(deleteMissing bool, actions []SyncAction) bool {
	for _, action := range actions {
		if action.Type == Add || action.Type == Modify {
			return true
		}
		if action.Type == Missing && deleteMissing {
			return true
		}
	}
	return false
}

func ExplainSyncActions(actions []SyncAction) {
	fmt.Println("=== Differences ===")
	for _, action := range actions {
		switch action.Type {
		case Add:
			fmt.Printf("Add: %s\n", action.Source.Name)
		case Modify:
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
		case Missing:
			fmt.Printf("Missing: %s\n", action.Source.Name)
		}
	}
}

func SyncFiles(actions []SyncAction, sourceRoot, targetRoot string, deleteMissing bool) {
	var wg sync.WaitGroup
	actionChan := make(chan SyncAction)
	var completed int64
	total := int64(len(actions))
	workerCount := optimalWorkerCount()

	go func() {
		for {
			done := atomic.LoadInt64(&completed)
			percent := float64(done) / float64(total) * 100
			fmt.Printf("\rProgress: %.2f%% (%d/%d)", percent, done, total)
			if done == total {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Println()
	}()

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for action := range actionChan {
				switch action.Type {
				case Add, Modify:
					srcPath := filepath.Join(sourceRoot, action.Source.Name)
					dstPath := filepath.Join(targetRoot, action.Source.Name)

					if err := copyFileWithModTime(srcPath, dstPath); err != nil {
						fmt.Printf("\nERROR copying %s: %v\n", srcPath, err)
					}
				case Missing:
					if deleteMissing {
						targetPath := filepath.Join(targetRoot, action.Source.Name)
						if err := os.Remove(targetPath); err != nil {
							fmt.Printf("\nERROR deleting %s: %v\n", targetPath, err)
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
	wg.Wait()
	return
}

func copyFileWithModTime(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
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

	return os.Chtimes(dst, info.ModTime(), info.ModTime())
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
