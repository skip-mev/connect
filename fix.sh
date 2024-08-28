#!/bin/bash

# Function to modify a single file
modify_file() {
    local file="$1"
    sed -i '' 's/\([[:alnum:]_]\+\)\.On(/\1.EXPECT()/g' "$file"
}

# Find all _test.go files and modify them
find . -name "*_test.go" | while read -r file; do
    echo "Processing $file"
    modify_file "$file"
done

echo "Modification complete."