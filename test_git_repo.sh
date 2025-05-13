#!/bin/bash

# Test script for the Git repository URL functionality

echo "Testing rule-tool with Git repository URL..."

# Test with a public Git repository URL
./rule-tool --git-repo=https://github.com/example/rules-repo --verbose --list

# Test with dry run
./rule-tool --git-repo=https://github.com/example/rules-repo --verbose --dry-run --link=example-rule

echo "Tests completed."
