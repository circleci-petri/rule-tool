#!/bin/bash

# Test script for cursor-rules CLI
# This script demonstrates how to use the CLI in non-interactive mode for testing

# Set up test paths
REPO_PATH=$(pwd)
TARGET_PATH="/tmp/cursor-rules-test"

# Create test target directory
mkdir -p "$TARGET_PATH"
echo "Created test target directory: $TARGET_PATH"

# Build the CLI
echo "Building cursor-rules CLI..."
go build -o bin/cursor-rules ./cmd/cursor-rules

# List all available rules
echo -e "\n=== Listing available rules ==="
./bin/cursor-rules --repo-path="$REPO_PATH" --list

# Test dry run mode for linking rules
echo -e "\n=== Testing dry run mode ==="
./bin/cursor-rules --repo-path="$REPO_PATH" --target-path="$TARGET_PATH" --link="holodeck-engineering,conventional-commits" --dry-run

# Link a specific rule
echo -e "\n=== Linking holodeck-engineering rule ==="
./bin/cursor-rules --repo-path="$REPO_PATH" --target-path="$TARGET_PATH" --link="holodeck-engineering"

# Verify the rule was linked
echo -e "\n=== Verifying rule was linked ==="
ls -la "$TARGET_PATH/.cursor/rules/"

# Unlink the rule
echo -e "\n=== Unlinking holodeck-engineering rule ==="
./bin/cursor-rules --repo-path="$REPO_PATH" --target-path="$TARGET_PATH" --unlink="holodeck-engineering"

# Verify the rule was unlinked
echo -e "\n=== Verifying rule was unlinked ==="
ls -la "$TARGET_PATH/.cursor/rules/" || echo "Directory is empty"

# Clean up
echo -e "\n=== Cleaning up ==="
rm -rf "$TARGET_PATH"
echo "Removed test target directory: $TARGET_PATH"

echo -e "\n=== Test completed ===" 