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

## ⭐ Random File Generator Script

Use the create_random_file.sh script to generate a single large file filled with random data — useful for testing sync speed, integrity checks, and progress reporting.

```bash
# Usage: ./create_random_file.sh [output_path] [size_in_mib]
#
# Arguments:
#   output_path   Path to the output file (default: random_1gb.bin)
#   size_in_mib   Size of the file in MiB (default: 1024 MiB = 1 GiB)
#
# This script:
# - Ensures the parent directory exists
# - Uses /dev/urandom to fill the file with random data
# - Displays progress during file creation

OUTPUT_FILE=${1:-random_1gb.bin}
SIZE_MB=${2:-1024}

mkdir -p "$(dirname "$OUTPUT_FILE")"

echo "Creating $SIZE_MB MiB random file at: $OUTPUT_FILE"
dd if=/dev/urandom of="$OUTPUT_FILE" bs=1M count="$SIZE_MB" status=progress

echo "Done."
```

### Example Usage

```bash
./create_random_file.sh ./test/source/file1.bin 512   # Creates a 512 MiB file
./create_random_file.sh ./bigfile.bin                 # Creates a 1 GiB file
for i in {1..5}; do ./create_large_file.sh "test/file_$i.bin" 100; done
```

---

## 🛠 Development

* Language: Go
* Standard packages only: `os`, `io`, `crypto/sha256`, `path/filepath`, `sync`, `flag`
* Recommended Go version: 1.19+
* OS: Linux
