package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syncr/models"
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

func CollectFileData(dir string) (files []models.FileData, err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("cannot stat file %s: %w", path, err)
			}

			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("cannot read file %s: %w", path, err)
			}
			defer file.Close()

			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			files = append(files, models.FileData{
				Name:        info.Name(),
				Checksum:    hex.EncodeToString(hash.Sum(nil)),
				Size:        info.Size(),
				ModTime:     info.ModTime(),
				Permissions: info.Mode().Perm(),
			})
		}
		return nil
	})

	return
}
