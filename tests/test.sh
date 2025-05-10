#!/bin/bash

# Test script for rule-tool CLI
# This script demonstrates how to use the CLI in non-interactive mode for testing

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_PATH="$SCRIPT_DIR/.."
TARGET_PATH="/tmp/rule-tool-test"

# Create test target directory
mkdir -p "$TARGET_PATH"
echo "Created test target directory: $TARGET_PATH"

# Clean up on exit
trap "rm -rf $TARGET_PATH" EXIT

# Build the CLI
echo "Building rule-tool CLI..."
go build -o bin/rule-tool ./cmd/rule-tool

# List all available rules
echo -e "\n=== Listing available rules ==="
./bin/rule-tool --repo-path="$REPO_PATH" --list

# Test dry run mode for linking rules
echo -e "\n=== Testing dry run mode ==="
./bin/rule-tool --repo-path="$REPO_PATH" --target-path="$TARGET_PATH" --link="holodeck-engineering,conventional-commits" --dry-run

# Link a specific rule
echo -e "\n=== Linking holodeck-engineering rule ==="
./bin/rule-tool --repo-path="$REPO_PATH" --target-path="$TARGET_PATH" --link="holodeck-engineering"

# Verify the rule was linked
echo -e "\n=== Verifying rule was linked ==="
if [ -d "$TARGET_PATH/rules/holodeck-engineering" ]; then
    echo "  ✅ Rule was linked successfully"
else
    echo "  ❌ Rule was not linked"
    exit 1
fi

# Unlink the rule
echo -e "\n=== Unlinking holodeck-engineering rule ==="
./bin/rule-tool --repo-path="$REPO_PATH" --target-path="$TARGET_PATH" --unlink="holodeck-engineering"

# Verify the rule was unlinked
echo -e "\n=== Verifying rule was unlinked ==="
ls -la "$TARGET_PATH/.cursor/rules/" || echo "Directory is empty"

# Clean up
echo -e "\n=== Cleaning up ==="
rm -rf "$TARGET_PATH"
echo "Removed test target directory: $TARGET_PATH"

echo -e "\n=== Test completed ===" 