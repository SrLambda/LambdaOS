#!/bin/bash
set -euo pipefail

SPECS_DIR="Specs"

if [ ! -d "$SPECS_DIR" ]; then
    echo "Error: $SPECS_DIR directory does not exist" >&2
    exit 1
fi

subdirs=$(find "$SPECS_DIR" -mindepth 1 -maxdepth 1 -type d)
if [ -z "$subdirs" ]; then
    echo "Error: $SPECS_DIR has no subdirectories" >&2
    exit 1
fi

md_files=$(find "$SPECS_DIR" -name '*.md')
if [ -z "$md_files" ]; then
    echo "Error: no .md files found in $SPECS_DIR" >&2
    exit 1
fi

# Report counts per subdirectory
echo "$subdirs" | while read -r dir; do
    count=$(find "$dir" -name '*.md' | wc -l)
    printf "  %s: %d .md file(s)\n" "$(basename "$dir")" "$count"
done

total=$(echo "$md_files" | wc -l)
echo "Total markdown files: $total"
exit 0
