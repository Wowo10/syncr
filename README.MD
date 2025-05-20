# Syncr
## Go Directory Sync Tool

A simple command-line tool written in Go to sync files from a source directory to a target directory.  
It compares files using SHA256 checksums and metadata, reports differences, and performs efficient multi-threaded copying with progress tracking.

**Tested only for Linux**

---

## ❗ Features

- Compare files by:
  - SHA256 checksum
  - Size
  - Modification time
  - Permissions
- Detect:
  - Modified files
  - Missing files in target
  - New files in source
- Multi-threaded file copying
- Live progress tracking with percentage
- Optional skipping of deletions (`--delete-missing`)
- Works with standard library only — **no external dependencies**

---

## 📦 Installation

1. Clone the repository:

```bash
git clone https://github.com/Wowo10/syncr
cd syncr
````

2. Build the binary:

```bash
go build -o syncr .
```

---

## 🎵 Usage

```bash
./syncr [options] <src> <target>
  --delete-missing
        Delete Missing files in target directory
```

### Example

```bash
./syncr --delete-missing ./test/source ./test/target
```

---

## 🛠 Development

* Language: Go
* Standard packages only: `os`, `io`, `crypto/sha256`, `path/filepath`, `sync`, `flag`
* Recommended Go version: 1.19+
* OS: Linux
