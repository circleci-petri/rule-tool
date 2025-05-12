# Rule Tool CLI

## Introduction

The Rule Tool CLI is a command-line interface application designed to manage rules for projects. It allows users to select and link rules from a central rules repository to their target projects, facilitating consistent application of guidelines and standards.

This tool simplifies managing and keeping rules consistent across many projects by using a central repository.

## Installation

### From Source

To install the Rule Tool CLI from source, follow these steps:

```bash
git clone https://github.com/rule-tool/rule-tool-cli.git
cd rule-tool-cli
go build -o rule-tool ./cmd/rule-tool
```

After building, move the `rule-tool` binary to a directory in your system's PATH or run it directly from the build location.

## Usage

The Rule Tool CLI supports both interactive and non-interactive modes.

### Interactive Mode

Run the `rule-tool` command without any specific flags to enter the interactive mode. This mode provides a Text User Interface (TUI) for selecting and managing rules.

```bash
# Run from your target project directory
rule-tool

# Specify path to rules repository
rule-tool --repo-path /path/to/rules/repo

# Specify a different target project
rule-tool --target-path /path/to/project
```

You can also use environment variables to set the rules repository path and target project path:

```bash
# Use rules path from environment variable
export RULE_TOOL_PATH=/path/to/rules/repo
rule-tool

# Use target path from environment variable
export RULE_TARGET_PATH=/path/to/project
rule-tool

# Use both environment variables
export RULE_TOOL_PATH=/path/to/rules/repo
export RULE_TARGET_PATH=/path/to/project
rule-tool
```

### Non-Interactive Mode

For scripting or automated workflows, use the following flags:

```bash
# List all available rules
rule-tool --list

# Link specific rules (comma-separated names)
rule-tool --link "rule1,rule2,rule3"

# Unlink specific rules (comma-separated names)
rule-tool --unlink "rule1,rule2"

# Dry run mode (show what would happen without making changes)
rule-tool --link "rule1,rule2" --dry-run

# Force non-interactive mode
rule-tool --non-interactive

# Enable verbose output
rule-tool --verbose [command]
```

### Configuration

The CLI determines the rules repository path and target project path based on the following precedence:

1.  Command-line flags (`--repo-path`, `--target-path`)
2.  Environment variables (`RULE_TOOL_PATH`, `RULE_TARGET_PATH`)
3.  Default value (current directory)

| Setting     | Flag          | Environment Variable | Default     |
| ----------- | ------------- | -------------------- | ----------- |
| Rules Path  | `--repo-path`   | `RULE_TOOL_PATH`       | Current dir |
| Target Path | `--target-path` | `RULE_TARGET_PATH`     | Current dir |

### Environment Variables

-   `RULE_TOOL_PATH`: Specifies the path to the rules repository.
-   `RULE_TARGET_PATH`: Specifies the path to the target project where rules will be linked.

**Note on Rules Repository Structure:** The `rule-tool` expects rules to be located in a directory named `rules` within the specified rules repository path. Rule files should have the `.mdc` extension.

## Features

-   Interactive selection and management of rules via a TUI.
-   Non-interactive command-line options for listing, linking, and unlinking rules.
-   Support for specifying rules repository and target project paths via flags or environment variables.
-   Automatic linking of selected rules to the target project (typically by creating symbolic links in a designated directory like `.cursor/rules`).
-   Management of existing rule links.
-   Dry run mode to preview changes before applying them.
-   Verbose output for detailed information.

## Development

### Requirements

-   Go 1.16 or higher

### Project Structure

```
/
├── cmd/
│   └── rule-tool/      # Main application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── rules/             # Rules loading and management
│   ├── linker/            # Symlink creation and management
│   └── ui/                # Terminal UI components
├── pkg/
│   └── models/            # Core data models
├── hooks/                 # Git hooks for development
└── README.md
```

### Building

```bash
go build -o rule-tool ./cmd/rule-tool
```

### Development Workflow

#### Git Hooks

This repository includes pre-commit hooks to ensure code quality. To install them:

```bash
./hooks/install-hooks.sh
```

The pre-commit hook will:

-   Build the code to catch compilation errors
-   Run tests to catch regressions
-   Prevent commits if any of these checks fail

For more details, see [hooks/README.md](hooks/README.md).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
