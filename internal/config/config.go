package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
)

// Environment variable constants
const (
	// EnvRulesPath is the environment variable name for specifying the rules path
	EnvRulesPath = "RULE_TOOL_PATH"
	// EnvTargetPath is the environment variable name for specifying the target project path
	EnvTargetPath = "RULE_TARGET_PATH"
)

// Config holds the global application configuration
type Config struct {
	// RulesRepoPath is the path to the rules repository
	RulesRepoPath string `env:"RULE_TOOL_PATH"`

	// TargetProjectPath is the path to the target project where rules will be linked
	TargetProjectPath string `env:"RULE_TARGET_PATH"`
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
	return filepath.Join(c.RulesRepoPath, "rules")
}
