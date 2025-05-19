#!/bin/bash

# Usage: ./create_random_file.sh [output_path] [size_in_mib]

OUTPUT_FILE=${1:-random_1gb.bin}
SIZE_MB=${2:-1024}

mkdir -p "$(dirname "$OUTPUT_FILE")"

echo "Creating $SIZE_MB MiB random file at: $OUTPUT_FILE"
dd if=/dev/urandom of="$OUTPUT_FILE" bs=1M count="$SIZE_MB" status=progress

echo "Done."
