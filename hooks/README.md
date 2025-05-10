# Git Hooks

This directory contains Git hooks for the llm-agent-rules repository.

## Available Hooks

- `pre-commit`: Runs before a commit is created, ensuring that all Go code builds and passes tests.

## Installation

Run the install script to set up the hooks:

```bash
./hooks/install-hooks.sh
```

This will create symbolic links in your local `.git/hooks` directory, pointing to the hooks in this directory.

## How the Pre-Commit Hook Works

The pre-commit hook performs the following checks:

1. Checks if any `.go` files are being committed
2. If Go files are found, the hook:
   - Temporarily stashes your unstaged changes to ensure a clean working directory
   - Runs `go build ./...` to verify the code compiles
   - Runs `go test ./...` to run all tests
   - Restores your unstaged changes when done

If any of these checks fail, the commit will be aborted with an error message.

## Manual Installation

If the install script doesn't work for any reason, you can manually install the hooks:

1. Make the hook executable:

   ```bash
   chmod +x hooks/pre-commit
   ```

2. Create a symbolic link in your `.git/hooks` directory:
   ```bash
   ln -sf ../../hooks/pre-commit .git/hooks/pre-commit
   ```

## Bypassing the Hook

In rare cases, you may need to bypass the pre-commit hook. You can do this with the `--no-verify` flag:

```bash
git commit --no-verify -m "Your commit message"
```

**Note:** This should be used sparingly, as it defeats the purpose of having the hook.
