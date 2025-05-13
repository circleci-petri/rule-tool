package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Repository represents a Git repository
type Repository struct {
	URL       string // URL of the Git repository
	CloneDir  string // Directory where the repository is cloned
	RulesPath string // Path to the rules directory within the repository
}

// New creates a new Git repository instance
func New(url string) *Repository {
	return &Repository{
		URL:       url,
		CloneDir:  "",
		RulesPath: "",
	}
}

// Clone clones the Git repository to a temporary directory
// Returns an error if the git command fails or if the repository URL is not set
func (r *Repository) Clone() error {
	if r.URL == "" {
		return errors.New("git repository URL is not set")
	}

	// Create a temporary directory for the clone
	tempDir, err := os.MkdirTemp("", "rule-tool-git-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", r.URL, tempDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Clean up the temporary directory if the clone fails
		_ = os.RemoveAll(tempDir)
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	// Set the clone directory
	r.CloneDir = tempDir
	
	// Set the rules path
	r.RulesPath = filepath.Join(r.CloneDir, ".cursor", "rules")
	
	return nil
}

// Cleanup removes the temporary directory where the git repository was cloned
func (r *Repository) Cleanup() error {
	if r.CloneDir == "" {
		return nil
	}

	err := os.RemoveAll(r.CloneDir)
	if err != nil {
		return fmt.Errorf("failed to remove git clone directory: %w", err)
	}

	r.CloneDir = ""
	r.RulesPath = ""
	return nil
}

// GetRulesPath returns the path to the rules directory in the repository
func (r *Repository) GetRulesPath() string {
	return r.RulesPath
}

// IsCloned returns true if the repository has been cloned
func (r *Repository) IsCloned() bool {
	return r.CloneDir != ""
}
