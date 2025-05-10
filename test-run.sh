#!/bin/bash

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Testing cursor-rules binary"
echo "Script located at: $SCRIPT_DIR"

# Set environment variables
export CURSOR_RULES_PATH="$SCRIPT_DIR"
export CURSOR_TARGET_PATH="$SCRIPT_DIR"

# Test the binary with verbose output
echo "Running binary with environment variables..."
./bin/cursor-rules-darwin-arm64 --verbose

# Test with command-line arguments
echo -e "\nRunning binary with command-line arguments..."
./bin/cursor-rules-darwin-arm64 --verbose --repo-path="$SCRIPT_DIR" --target-path="$SCRIPT_DIR"

# Test with dry-run and list options
echo -e "\nRunning binary with list option..."
./bin/cursor-rules-darwin-arm64 --verbose --repo-path="$SCRIPT_DIR" --target-path="$SCRIPT_DIR" --list 