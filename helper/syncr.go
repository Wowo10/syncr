package helper

import "fmt"

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
