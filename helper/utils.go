package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsDirectoryWritable(dir string) bool {
	testFile := filepath.Join(dir, ".perm_check.tmp")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

type FileData struct {
	Name     string
	Checksum string
}

func CollectFileData(dir string) (files []FileData, err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("cannot read file %s: %w", path, err)
			}
			defer file.Close()

			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			checksum := hex.EncodeToString(hash.Sum(nil))
			files = append(files, FileData{
				Name:     path,
				Checksum: checksum,
			})
		}
		return nil
	})

	return
}
