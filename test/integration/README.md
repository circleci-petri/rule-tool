# Integration Tests

This directory contains integration tests for the rule-tool application.

## Requirements

To run the integration tests, you need:

1. Go 1.20 or higher
2. The rule-tool binary built and accessible in one of the following locations:
   - `bin/rule-tool-{GOOS}-{GOARCH}` (platform-specific binary, e.g., rule-tool-darwin-arm64)
   - `cmd/rule-tool/rule-tool` (relative to project root)

## Running Tests

From the project root directory:

```bash
go test -v ./test/integration/...
```

Or from the test/integration directory:

```bash
cd test/integration
go test -v ./...
```

## Test Cases

The integration tests verify:

1. Basic functionality
2. Command-line flags
3. Error handling
4. Rule linking and unlinking
5. Directory path handling
