# Integration Tests

This directory contains integration tests for the cursor-rules application.

## Prerequisites

Before running these tests, make sure to:

1. Build the cursor-rules binary (optional, as tests will build it automatically if not found):

   ```bash
   task build
   ```

2. The tests will look for the binary in these locations (in order):
   - `bin/cursor-rules-{GOOS}-{GOARCH}` (platform-specific binary, e.g., cursor-rules-darwin-arm64)
   - `cmd/cursor-rules/cursor-rules` (relative to project root)
   - Environment variable `CURSOR_RULES_BINARY_PATH` (if set)

## Running Tests

To run the integration tests:

```bash
go test ./test/integration/...
```

Note: If the binary is not found, the test will automatically build it using `task build`.

## CI Integration

In CI environments, you can:

1. Build the binary in an earlier step using `task build`
2. Set the `CURSOR_RULES_BINARY_PATH` environment variable to point to the built binary
3. Run the integration tests

This approach ensures tests use the same binary that will be deployed.
