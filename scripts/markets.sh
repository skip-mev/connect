#!/bin/bash

# Ensure the environment variable is set
if [ -z "$GENESIS_MARKETS" ]; then
  echo "Environment variable GENESIS_MARKETS is not set."
  exit 1
fi

# Print the value of the environment variable for debugging
echo "GENESIS_MARKETS is set to: $GENESIS_MARKETS"

# Check if the file exists
if [ ! -f "$GENESIS_MARKETS" ]; then
  echo "File not found: $GENESIS_MARKETS"
  exit 1
fi

# Print a message indicating the file exists
echo "File found: $GENESIS_MARKETS"

# Define the temporary file path for debugging purposes
temp_file="/tmp/extracted_json_content.txt"

# Extract the JSON content nested within the RaydiumMarketMapJSON variable and store it in the temporary file
awk '/RaydiumMarketMapJSON = `/,/`$/{ if ($0 !~ /RaydiumMarketMapJSON = `|`$/) {print $0} }' "$GENESIS_MARKETS" > "$temp_file"

# Check if the temporary file has content
if [ ! -s "$temp_file" ]; then
  echo "No content extracted to the temporary file. Please check the file and patterns."
  exit 1
fi

# Debug: Print the path of the temporary file and the first few lines of its content
echo "Temporary file created at: $temp_file"
echo "First few lines of the extracted content:"
head -n 10 "$temp_file"

# Use cat to display the extracted JSON content
cat "$temp_file"
