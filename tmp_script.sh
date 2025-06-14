#!/bin/bash

# Exit on errors
set -e

# Find all *.gohtml files
find . -type f -name "*.gohtml" | while read -r file; do
    # Construct the new filename
    newfile="${file%.gohtml}.gohtml"

    # Use git mv to preserve history
    git mv "$file" "$newfile"

    echo "Renamed: $file → $newfile"
done
