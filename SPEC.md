# Rule Tool CLI Tool Specification

## Overview

A CLI tool written in Go using the Charm.sh suite of tools to manage rules across projects. The tool will help users select and link rules from the checked out repo to their target projects.

## Core Features

### Phase 1

- [x] Load rules from the local checked out repository
- [x] List available rules with descriptions when executed in a target project
- [x] Interactive selection UI for choosing rules
- [x] Create symlinks from selected rules to target project's `.cursor/rules` directory
- [x] Create `.cursor/rules` directory if it doesn't exist
- [x] Basic error handling and validation

Phase 1 is complete.
## Technical Implementation

### Tools & Dependencies

- Go programming language
- Charm.sh libraries:
  - Bubble Tea for terminal UI
  - Lip Gloss for styling
  - Huh for interactive forms/selection

### Modes of Operation

- **Interactive Mode**: Terminal UI with full styling and interactivity
- **Non-Interactive Mode**: Plain text output with no styling or color for use in scripts and automated environments

### Directory Structure

```
/
├── cmd/
│   └── rule-tool/
│       └── main.go
├── internal/
│   ├── config/
│   ├── rules/
│   ├── linker/
│   └── ui/
├── pkg/
│   └── models/
├── tests/
├── go.mod
├── go.sum
└── README.md
```

### Core Components

- **Rules Manager**: Handles loading and parsing rules from the local repository structure
- **Rule Linker**: Creates symlinks between the repository and target project
- **UI Layer**: Interactive terminal UI for rule selection and management
- **Configuration**: Manages tool configuration and settings

## Work Checklist

### Setup & Foundation

- [x] Initialize Go module
- [x] Add Charm.sh dependencies
- [x] Create basic project structure
- [x] Write README.md with usage instructions

### Rules Management

- [x] Create rule model structure
- [x] Implement rule parsing and loading from local directories
- [ ] Add validation for rule format


### Testing

- [x] Unit tests for rule parsing
- [x] Integration tests for linking functionality
- [x] Mock testing for UI components

### Refinement

- [ ] Add configuration options
- [ ] Implement error reporting
- [ ] Add logging
- [ ] Performance optimizations

## Testing Strategy

- Unit tests for core logic components
- Integration tests for file system operations
- Mock testing for UI components
- End-to-end testing with sample rules

## Future Considerations

- Rule updating mechanism
- Custom rule creation helpers
- Rule search and filtering
- Rule templates and customization

## Future Feature Ideas

### Expansion Features

- **Rule Format Transformation**: Ability to transform rules into other formats for integration with other tools
- **Rule Synchronization**: Sync rules across multiple repositories automatically
- **Rule Analytics**: Track which rules are most used/helpful across projects
