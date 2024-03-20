#!/bin/bash
# Usage: make diff-local-configs

echo "Diffing local configs against HEAD"
diff_output=$(git diff HEAD ./config/local/*.json)

if [ -z "$diff_output" ]; then
    echo "Local configs are up to date."
else
    echo "Local configs are not updated. Run make update-local-config to update."
    echo "$diff_output"
    exit 1
fi
