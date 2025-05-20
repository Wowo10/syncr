package models

import (
	"os"
	"time"
)

type FileData struct {
	Name        string
	Checksum    string
	Size        int64
	ModTime     time.Time
	Permissions os.FileMode
}

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
