# Rule Tool CLI Usage Guide

This document explains the basic usage of the Rule Tool CLI.

## Basic Workflow

1. Check out this project repository
2. Build the CLI tool using the task file
3. Run the CLI tool from the root of this project
4. Specify the path to your target project when prompted

The tool assumes:

- You're running it from the root of this rules project
- Rules are located in the `rules` directory

Example:

```bash
# Clone this repository
git clone https://github.com/your-org/llm-agent-rules.git
cd llm-agent-rules

# Build using Go Task
task build

# Run the tool or move it to your path
./bin/rule-tool
```

## Usage Options

### Interactive Mode (Default)

```bash
# Basic usage
rule-tool

# Specify target project
rule-tool --target-path /path/to/project

# Specify rules repository path
rule-tool --repo-path /path/to/rules/repo

# Use environment variable for rules repository path
export RULE_TOOL_PATH=/path/to/rules/repo
rule-tool

# Use environment variable for target project path
export RULE_TARGET_PATH=/path/to/project
rule-tool

# Use both environment variables
export RULE_TOOL_PATH=/path/to/rules/repo
export RULE_TARGET_PATH=/path/to/project
rule-tool
```

### Non-Interactive Mode

Non-Interactive mode really just here for testing. Don't know why you'd want to use it.

```bash
# List available rules
rule-tool --list

# Link specific rules
rule-tool --link "rule1,rule2,rule3"
```

### Environment Variables

The CLI supports the following environment variables:

- `RULE_TOOL_PATH`: Path to the rules repository (overridden by `--repo-path` flag if provided)
- `RULE_TARGET_PATH`: Path to the target project (overridden by `--target-path` flag if provided)

## How It Works

1. The tool scans for .mdc files in the rules directory/sub-directories
2. Each rule is parsed to extract metadata
3. You select which rules to apply
4. The tool creates symbolic links in your project's `.cursor/rules` directory
