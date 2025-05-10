package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary runs 'task build' to build the rule-tool binary
func buildBinary() error {
	// Get the project root directory (2 levels up from test/integration)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	projectRoot := filepath.Dir(filepath.Dir(wd))

	// Run task build
	cmd := exec.Command("task", "build")
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Building binary using 'task build'...")
	return cmd.Run()
}

// TestMain is used to setup any prerequisites for integration tests
func TestMain(m *testing.M) {
	// Display information about running integration tests
	fmt.Println("Running integration tests...")

	// Check if binary exists
	_, err := findBinary()
	if err != nil {
		fmt.Println("Binary not found. Building it automatically...")
		if buildErr := buildBinary(); buildErr != nil {
			fmt.Printf("Failed to build binary: %v\n", buildErr)
			os.Exit(1)
		}

		// Verify the binary was built successfully
		binaryPath, verifyErr := findBinary()
		if verifyErr != nil {
			fmt.Printf("Binary still not found after build: %v\n", verifyErr)
			os.Exit(1)
		}
		fmt.Printf("Using binary at: %s\n", binaryPath)
	}

	// Run all tests
	os.Exit(m.Run())
}
