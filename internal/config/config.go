package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/circleci/llm-agent-rules/internal/git"
	"github.com/sethvargo/go-envconfig"
)

// Environment variable constants
const (
	// EnvRulesPath is the environment variable name for specifying the rules path
	EnvRulesPath = "RULE_TOOL_PATH"
	// EnvTargetPath is the environment variable name for specifying the target project path
	EnvTargetPath = "RULE_TARGET_PATH"
	// EnvGitRepoURL is the environment variable name for specifying the git repository URL
	EnvGitRepoURL = "RULE_GIT_REPO_URL"
	// EnvEditor is the environment variable name for specifying the editor
	EnvEditor = "RULE_EDITOR"
)

// Editor represents the code editor being used
type Editor string

// Supported editors
const (
	EditorCursor   Editor = "cursor"
	EditorWindsurf Editor = "windsurf"
)

// Config holds the global application configuration
type Config struct {
	// RulesRepoPath is the path to the rules repository
	RulesRepoPath string `env:"RULE_TOOL_PATH"`

	// TargetProjectPath is the path to the target project where rules will be linked
	TargetProjectPath string `env:"RULE_TARGET_PATH"`

	// GitRepoURL is the URL of the git repository containing the rules
	GitRepoURL string `env:"RULE_GIT_REPO_URL"`

	// Editor is the code editor being used
	Editor Editor `env:"RULE_EDITOR"`

	// UseGitRepo indicates whether to use the git repository URL instead of local path
	UseGitRepo bool

	// GitRepo holds the git repository instance when using a git repo
	GitRepo *git.Repository
}

// New creates a new configuration with default values
func New() *Config {
	var cfg Config

	// Process environment variables
	_ = envconfig.Process(context.Background(), &cfg)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Default to current working directory for rules repo if not set
	if cfg.RulesRepoPath == "" {
		cfg.RulesRepoPath = cwd
	} else if !filepath.IsAbs(cfg.RulesRepoPath) {
		// Convert relative path to absolute path
		cfg.RulesRepoPath = filepath.Join(cwd, cfg.RulesRepoPath)
	}

	// Default to current working directory for target project if not set
	if cfg.TargetProjectPath == "" {
		cfg.TargetProjectPath = cwd
	} else if !filepath.IsAbs(cfg.TargetProjectPath) {
		// Convert relative path to absolute path
		cfg.TargetProjectPath = filepath.Join(cwd, cfg.TargetProjectPath)
	}

	// Default to Cursor editor if not set
	if cfg.Editor == "" {
		cfg.Editor = EditorCursor
	}

	// Set UseGitRepo flag if GitRepoURL is provided and initialize Git repo
	if cfg.GitRepoURL != "" {
		cfg.UseGitRepo = true
		cfg.GitRepo = git.New(cfg.GitRepoURL)
	}

	return &cfg
}

// SetRulesRepoPath sets the path to the rules repository
// Command line flags take precedence over environment variables
// Converts relative paths to absolute paths
func (c *Config) SetRulesRepoPath(path string) {
	// Convert relative paths to absolute paths
	if !filepath.IsAbs(path) {
		// Get current working directory
		cwd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(cwd, path)
		}
	}
	c.RulesRepoPath = path
}

// SetTargetProjectPath sets the path to the target project
// Command line flags take precedence over environment variables
// Converts relative paths to absolute paths
func (c *Config) SetTargetProjectPath(path string) {
	// Convert relative paths to absolute paths
	if !filepath.IsAbs(path) {
		// Get current working directory
		cwd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(cwd, path)
		}
	}
	c.TargetProjectPath = path
}

// ValidateRulesRepoPath checks if the rules repository path is valid
func (c *Config) ValidateRulesRepoPath() bool {
	// Check if path exists
	if _, err := os.Stat(c.RulesRepoPath); os.IsNotExist(err) {
		return false
	}

	// Check if it's a directory
	info, err := os.Stat(c.RulesRepoPath)
	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

// ValidateTargetProjectPath checks if the target project path is valid
func (c *Config) ValidateTargetProjectPath() bool {
	// Check if path exists
	if _, err := os.Stat(c.TargetProjectPath); os.IsNotExist(err) {
		return false
	}

	// Check if it's a directory
	info, err := os.Stat(c.TargetProjectPath)
	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

// GetRulesDir returns the expected directory for rules within the repository
func (c *Config) GetRulesDir() string {
	if c.UseGitRepo && c.GitRepo != nil && c.GitRepo.IsCloned() {
		return c.GitRepo.GetRulesPath()
	}
	return filepath.Join(c.RulesRepoPath, "rules")
}

// SetGitRepoURL sets the URL of the git repository containing the rules
func (c *Config) SetGitRepoURL(url string) {
	c.GitRepoURL = url
	if url != "" {
		c.UseGitRepo = true
		c.GitRepo = git.New(url)
	} else {
		c.UseGitRepo = false
		c.GitRepo = nil
	}
}

// SetEditor sets the editor to use
func (c *Config) SetEditor(editor string) {
	switch editor {
	case string(EditorCursor):
		c.Editor = EditorCursor
	case string(EditorWindsurf):
		c.Editor = EditorWindsurf
	default:
		// Default to Cursor if invalid editor is provided
		c.Editor = EditorCursor
	}
}

// GetRulesDirectory returns the directory where rules should be linked based on the editor
func (c *Config) GetRulesDirectory() string {
	switch c.Editor {
	case EditorCursor:
		return filepath.Join(".cursor", "rules")
	case EditorWindsurf:
		return filepath.Join(".windsurf", "rules")
	default:
		return filepath.Join(".cursor", "rules")
	}
}
