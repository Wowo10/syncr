#!/bin/bash

OUTPUT_FILE=${1:-random_1gb.bin}

mkdir -p "$(dirname "$OUTPUT_FILE")"

echo "Creating 1 GiB random file at: $OUTPUT_FILE"
dd if=/dev/urandom of="$OUTPUT_FILE" bs=1M count=1024 status=progress

echo "Done."
