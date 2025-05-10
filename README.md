# LLM Agent Rules Collection (WIP)

A collection of rules for LLM agents used within our organization. We use the cursor rules 'spec' to guide these, but any LLM should be able to take advantage if you tell them about it.

## What is this?

This repository contains standardized rules for LLM agents that establish consistent coding practices, workflows, and collaboration across our organization. These rules guide AI behavior when working on specific tasks or with particular technologies.

## How Rules Work with Cursor

Cursor follows specified rules or guidelines when generating or modifying code. When rules are provided to the AI, it adapts its responses to match established requirements.

Rules can be:

- Applied to specific file patterns using globs
- Used to enforce coding standards
- Applied to maintain consistency in testing approaches
- Used to standardize documentation and commit messages

## Available Rules

The repository contains the following rules:

| Rule Name                                                                       | Description                                                                                           | File Pattern             |
| ------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ------------------------ |
| [Conventional Commit Messages](./rules/conventional-commit-messages.md)         | Standards for writing commit messages following the Conventional Commits specification                | Git commits              |
| [Cypress to Playwright Conversion](./rules/cypress-to-playwright-conversion.md) | Guidelines for converting Cypress tests to Playwright tests with consistent organization and patterns | `*.spec.ts`, `*.test.ts` |
| [Multiline Commit Messages](./rules/multiline-commit-messages.md)               | Standards for properly formatting multiline git commit messages using multiple -m flags               | Git commits              |
| [One Rule to Rule Them All](./rules/one-rule-to-rule-them-all.md)               | Cursor Rules Location Rule                                                                            | \*                       |
| [Playwright Testing](./rules/playwright-testing.md)                             | Best practices for writing and maintaining Playwright tests                                           | `*.spec.ts`              |
| [Storybook Story Consolidation](./rules/storybook-story-consolidation.md)       | Guidelines for consolidating Storybook stories to reduce snapshot count                               | `*.stories.tsx`          |

## How to Use These Rules

With Cursor AI, rules can be referenced by name or file pattern:

1. Reference a rule when requesting AI assistance:

   ```
   "Create a test for this component following our playwright-testing rule"
   ```

2. Cursor automatically applies relevant rules based on file patterns.

3. Explicitly request application of a specific rule:
   ```
   "Apply the conventional-commit-messages rule to this commit"
   ```

## Important: Rule File Location

Cursor only processes rule files that are directly in the `.cursor/rules` directory. Files that are nested in subdirectories within `.cursor/rules` will **not** be processed. For example:

- ✅ `.cursor/rules/my-rule.md` - Will be processed
- ❌ `.cursor/rules/nested/my-rule.md` - Will NOT be processed

When linking rules to the `.cursor/rules` directory, path separators in the original filepath will be converted to underscores in the filename. For example:

- `rules/nested/my-rule.md` would be linked as `.cursor/rules/nested_my-rule.md`

## Contributing New Rules

To contribute a new rule:

1. Create a new markdown file in the `rules` directory
2. Use the following template format:

   ```markdown
   ---
   description: Brief description of what the rule does
   globs: file/pattern/*.extension
   ---

   # Rule Title

   Detailed explanation of the rule...
   ```

3. Submit a pull request with your new rule

## Best Practices

- Reference specific rules when working with Cursor AI
- Keep rules focused on a single concern
- Provide clear examples in rule documents
- Update rules as practices evolve

# Rule Tool CLI

## Introduction

A command-line tool for managing Cursor AI rules across projects. This tool helps you select and link rules from a central repository to your target projects.

## Installation

### From Source

```bash
git clone https://github.com/rule-tool/rule-tool-cli.git
cd rule-tool-cli
go build -o rule-tool ./cmd/rule-tool
```

Then, move the binary to your PATH or run it directly.

## Usage

The Rule Tool CLI can be used both interactively and non-interactively.

### Interactive Mode

```bash
# Run from your target project directory
rule-tool

# Specify path to rules repository
rule-tool --repo-path /path/to/rules/repo

# Specify a different target project
rule-tool --target-path /path/to/project

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

For scripting or automated usage, you can use these options:

```bash
# List all available rules
rule-tool --list

# Link specific rules
rule-tool --link "rule1,rule2,rule3"

# Unlink specific rules
rule-tool --unlink "rule1,rule2"

# Dry run mode (show what would happen without making changes)
rule-tool --link "rule1,rule2" --dry-run

# Force non-interactive mode
rule-tool --non-interactive
```

### Configuration

The CLI supports the following configuration methods in order of precedence:

1. Command-line flags (highest precedence)
2. Environment variables
3. Default values (lowest precedence)

| Setting     | Flag          | Environment Variable | Default     |
| ----------- | ------------- | -------------------- | ----------- |
| Rules Path  | --repo-path   | RULE_TOOL_PATH       | Current dir |
| Target Path | --target-path | RULE_TARGET_PATH     | Current dir |

### Environment Variables

- `RULE_TOOL_PATH`: Path to the rules repository
- `RULE_TARGET_PATH`: Path to the target project where rules will be linked

## Features

- Interactive selection of rules with descriptions
- Automatic linking of selected rules to your project
- Creation of `.cursor/rules` directory if it doesn't exist
- Management of existing rule links

## How It Works

1. The tool scans your rules repository for markdown files
2. Each rule is parsed to extract its name, description, and globs
3. You select which rules to apply to your project
4. The tool creates symbolic links from your rules repository to your project's `.cursor/rules` directory

## Development

### Requirements

- Go 1.16 or higher
- Libraries used:
  - [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI
  - [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling
  - [Huh](https://github.com/charmbracelet/huh) for interactive forms/selection

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

- Build the code to catch compilation errors
- Run tests to catch regressions
- Prevent commits if any of these checks fail

For more details, see [hooks/README.md](hooks/README.md).

## License

Apache License 2.0

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
