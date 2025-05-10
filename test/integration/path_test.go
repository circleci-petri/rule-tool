package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/circleci/llm-agent-rules/internal/config"
)

// findBinary locates the cursor-rules binary in known locations
// Returns the path to the binary or an error if not found
func findBinary() (string, error) {
	// Determine build output path based on platform
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	binaryName := fmt.Sprintf("cursor-rules-%s-%s", goos, goarch)
	if goos == "windows" {
		binaryName += ".exe"
	}

	// Look for binary in common locations
	binaryPaths := []string{
		filepath.Join("..", "..", "bin", binaryName),                     // From project root /bin dir
		filepath.Join("..", "..", "cmd", "cursor-rules", "cursor-rules"), // Relative to test/integration
		filepath.Join(".", "cmd", "cursor-rules", "cursor-rules"),        // From project root
	}

	// Check common locations
	for _, path := range binaryPaths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath, nil
		}
	}

	// Check if binary exists in a CI-provided location
	ciBinaryPath := os.Getenv("CURSOR_RULES_BINARY_PATH")
	if ciBinaryPath != "" {
		if _, err := os.Stat(ciBinaryPath); err == nil {
			return ciBinaryPath, nil
		}
	}

	return "", fmt.Errorf("cursor-rules binary not found - please run 'task build' first")
}

// TestDefaultPathsIntegration tests that the application correctly defaults
// to the current working directory when no paths are provided
func TestDefaultPathsIntegration(t *testing.T) {
	// Find the binary (but don't build it)
	binaryPath, err := findBinary()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Using binary at: %s", binaryPath)

	// Save current env vars to restore later
	oldRulesEnv := os.Getenv(config.EnvRulesPath)
	oldTargetEnv := os.Getenv(config.EnvTargetPath)
	defer func() {
		os.Setenv(config.EnvRulesPath, oldRulesEnv)
		os.Setenv(config.EnvTargetPath, oldTargetEnv)
	}()

	// Clear environment variables
	os.Unsetenv(config.EnvRulesPath)
	os.Unsetenv(config.EnvTargetPath)

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary rules directory
	tempDir, err := os.MkdirTemp("", "cursor-rules-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a rules directory
	rulesDir := filepath.Join(tempDir, "rules")
	if err := os.Mkdir(rulesDir, 0755); err != nil {
		t.Fatalf("Failed to create rules directory: %v", err)
	}

	// Create a sample rule file
	rulePath := filepath.Join(rulesDir, "sample-rule.mdc")
	ruleContent := `---
description: A sample rule for testing
globs: "*"
---

# Sample Rule

## Rule
This is a sample rule.
`
	if err := os.WriteFile(rulePath, []byte(ruleContent), 0644); err != nil {
		t.Fatalf("Failed to write sample rule: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		envVars     map[string]string
		expectError bool
	}{
		{
			name:        "No paths provided - should use cwd",
			args:        []string{"--list", "--non-interactive"},
			envVars:     map[string]string{},
			expectError: true, // Will fail because cwd doesn't have a rules directory
		},
		{
			name:        "Repo path provided",
			args:        []string{"--list", "--non-interactive", "--repo-path", tempDir},
			envVars:     map[string]string{},
			expectError: false, // Should work with the valid tempDir/rules directory
		},
		{
			name:        "Repo path via env var",
			args:        []string{"--list", "--non-interactive"},
			envVars:     map[string]string{config.EnvRulesPath: tempDir},
			expectError: false, // Should work with env var
		},
		{
			name:        "Both paths provided",
			args:        []string{"--list", "--non-interactive", "--repo-path", tempDir, "--target-path", cwd},
			envVars:     map[string]string{},
			expectError: false, // Should work with both paths
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clear env vars after each test
				for k := range tc.envVars {
					os.Unsetenv(k)
				}
			}()

			// Run the command
			cmd := exec.Command(binaryPath, tc.args...)
			output, err := cmd.CombinedOutput()

			t.Logf("Command output: %s", output)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
