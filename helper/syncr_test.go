package helper_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"syncr/helper"
	"testing"
)

func TestCompareFileData(t *testing.T) {
	sourceDir, err := os.MkdirTemp("/tmp", "syncrsourcetest")
	if err != nil {
		t.Fatalf("cannot create test dir: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	testFilePath := filepath.Join(sourceDir, "testfile.txt")
	content := []byte("This is a test file")

	if err := os.WriteFile(testFilePath, content, 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}

	targetDir, err := os.MkdirTemp("/tmp", "syncrtargettest")
	if err != nil {
		t.Fatalf("cannot create test dir: %v", err)
	}
	defer os.RemoveAll(targetDir)

	sourceData, err := helper.CollectFileData(sourceDir)
	if err != nil {
		t.Fatalf("cannot collect file data: %v", err)
	}

	targetData, err := helper.CollectFileData(targetDir)
	if err != nil {
		t.Fatalf("cannot collect file data: %v", err)
	}

	actions := helper.CompareFileData(sourceData, targetData)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}

	if actions[0].Type != helper.Add {
		t.Errorf("expected action type Add, got %s", actions[0].Type)
	}

	if filepath.Base(actions[0].Source.Name) != filepath.Base(testFilePath) {
		t.Errorf("expected file name %s, got %s", testFilePath, actions[0].Source.Name)
	}

	if actions[0].Source.Checksum != fmt.Sprintf("%x", sha256.Sum256(content)) {
		t.Errorf("expected checksum %s, got %s", fmt.Sprintf("%x", sha256.Sum256(content)), actions[0].Source.Checksum)
	}

	if actions[0].Source.Size != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), actions[0].Source.Size)
	}

	if actions[0].Source.Permissions != os.FileMode(0644) {
		t.Errorf("expected permissions %v, got %v", os.FileMode(0644), actions[0].Source.Permissions)
	}
}
