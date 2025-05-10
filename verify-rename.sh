#!/bin/bash

echo "Verifying Rule Tool build and run after renaming..."

# Build the binary
echo "Building rule-tool..."
go build -o bin/rule-tool ./cmd/rule-tool

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build succeeded!"

# Set environment variables
export RULE_TOOL_PATH=$(pwd)
export RULE_TARGET_PATH=$(pwd)

# Run the tool with --list option to check functionality
echo "Running rule-tool --list..."
./bin/rule-tool --list

if [ $? -ne 0 ]; then
    echo "❌ Test run failed!"
    exit 1
fi

echo "✅ Test run succeeded!"
echo "Rename verification completed successfully!" 