#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
GIT_HOOKS_DIR="${REPO_ROOT}/.git/hooks"

# Make pre-commit hook executable
chmod +x "${SCRIPT_DIR}/pre-commit"

# Create symlink to the pre-commit hook
ln -sf "${SCRIPT_DIR}/pre-commit" "${GIT_HOOKS_DIR}/pre-commit"

echo "✅ Git hooks installed successfully!"
echo "Pre-commit hook will run automatically on commit." 