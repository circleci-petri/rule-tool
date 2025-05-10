package linker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/circleci/llm-agent-rules/pkg/models"
)

func TestLinkRuleCreatesRelativePath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "linker-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the .cursor/rules directory
	rulesDir := filepath.Join(tmpDir, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("Failed to create rules directory: %v", err)
	}

	// Create a test rule file in a separate location
	ruleDir := filepath.Join(tmpDir, "repo", "rules")
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		t.Fatalf("Failed to create rule repo directory: %v", err)
	}

	rulePath := filepath.Join(ruleDir, "test-rule.mdc")
	if err := os.WriteFile(rulePath, []byte("test rule content"), 0644); err != nil {
		t.Fatalf("Failed to create test rule file: %v", err)
	}

	// Create a mock rule
	rule := &models.Rule{
		Name:    "test-rule",
		Path:    rulePath,
		Content: "test rule content",
	}

	// Create a linker pointing to the temp directory
	l := NewLinker(tmpDir)
	l.SetVerbose(true)

	// Link the rule
	if err := l.LinkRule(rule); err != nil {
		t.Fatalf("LinkRule failed: %v", err)
	}

	// Verify that the symlink was created
	linkPath := filepath.Join(rulesDir, "test-rule.mdc")
	if _, err := os.Stat(linkPath); os.IsNotExist(err) {
		t.Fatalf("Symlink was not created")
	}

	// Verify that the symlink is relative, not absolute
	linkTarget, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if filepath.IsAbs(linkTarget) {
		t.Errorf("Symlink target is absolute, expected relative path: %s", linkTarget)
	}

	// Verify that the symlink resolves to the correct file
	resolvedPath := filepath.Join(filepath.Dir(linkPath), linkTarget)
	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		t.Errorf("Symlink target does not resolve correctly: %s", resolvedPath)
	}
} 