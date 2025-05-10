#!/bin/bash

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Testing rule-tool binary"
echo "Script located at: $SCRIPT_DIR"

# Set environment variables
export RULE_TOOL_PATH="$SCRIPT_DIR"
export RULE_TARGET_PATH="$SCRIPT_DIR"
export VERBOSE=1

# Test the binary with verbose output
echo "Running binary with environment variables..."
./bin/rule-tool-darwin-arm64 --verbose

# Test with command-line arguments
echo -e "\nRunning binary with command-line arguments..."
./bin/rule-tool-darwin-arm64 --verbose --repo-path="$SCRIPT_DIR" --target-path="$SCRIPT_DIR"

# Test with dry-run and list options
echo -e "\nRunning binary with list option..."
./bin/rule-tool-darwin-arm64 --verbose --repo-path="$SCRIPT_DIR" --target-path="$SCRIPT_DIR" --list 