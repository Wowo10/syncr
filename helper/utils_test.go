package helper_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"syncr/helper"
	"testing"
)

func TestIsDir(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"empty", "", false},
		{"not existing", "/tmp/notexisting", false},
		{"tmp", "/tmp", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.IsDirectory(tt.path); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDirWriteable(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"not existing", "/tmp/notexisting", false},
		{"tmp", "/tmp", true},
		{"root", "/root", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.IsDirectoryWritable(tt.path); got != tt.want {
				t.Errorf("IsDirWriteable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectFileData(t *testing.T) {
	dir, err := os.MkdirTemp("/tmp", "syncrtest")
	if err != nil {
		t.Fatalf("cannot create test dir: %v", err)
	}
	defer os.RemoveAll(dir)

	testFilePath := filepath.Join(dir, "testfile.txt")
	content := []byte("This is a test file")

	if err := os.WriteFile(testFilePath, content, 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}

	files, err := helper.CollectFileData(dir)
	if err != nil {
		t.Fatalf("cannot collect file data: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if files[0].Name != testFilePath {
		t.Errorf("expected file name %s, got %s", testFilePath, files[0].Name)
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(content))
	if files[0].Checksum != checksum {
		t.Errorf("expected checksum %s, got %s", checksum, files[0].Checksum)
	}

	if files[0].Size != int64(len(content)) {
		t.Errorf("expected file size %d, got %d", len(content), files[0].Size)
	}

	if files[0].Permissions != 0o644 {
		t.Errorf("expected file permission 0644, got %o", files[0].Permissions)
	}
}
