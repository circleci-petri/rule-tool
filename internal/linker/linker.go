package linker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/circleci/llm-agent-rules/internal/config"
	"github.com/circleci/llm-agent-rules/pkg/models"
)

// Linker handles creating symlinks between rules repository and target project
type Linker struct {
	TargetDir string
	DryRun    bool
	Verbose   bool
	Config    *config.Config
}

// NewLinker creates a new linker for the specified target directory
func NewLinker(targetDir string) *Linker {
	return &Linker{
		TargetDir: targetDir,
		DryRun:    false,
		Verbose:   false,
		Config:    nil,
	}
}

// SetDryRun enables or disables dry run mode
func (l *Linker) SetDryRun(dryRun bool) {
	l.DryRun = dryRun
}

// SetVerbose enables or disables verbose output
func (l *Linker) SetVerbose(verbose bool) {
	l.Verbose = verbose
}

// SetConfig sets the configuration for the linker
func (l *Linker) SetConfig(cfg *config.Config) {
	l.Config = cfg
}

// GetRulesDirectory returns the appropriate rules directory path based on configuration
func (l *Linker) GetRulesDirectory() string {
	if l.Config != nil {
		return l.Config.GetRulesDirectory()
	}
	// Default to .cursor/rules if no config is set
	return filepath.Join(".cursor", "rules")
}

// EnsureTargetDirectory ensures the editor-specific rules directory exists in the target project
func (l *Linker) EnsureTargetDirectory() error {
	// Get the rules directory based on the selected editor
	rulesDir := filepath.Join(l.TargetDir, l.GetRulesDirectory())

	// Check if directory exists
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		if l.DryRun {
			if l.Verbose {
				fmt.Printf("Would create directory: %s\n", rulesDir)
			}
			return nil
		}

		err := os.MkdirAll(rulesDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	return nil
}

// LinkRule creates a symlink from the rule source to the target directory
func (l *Linker) LinkRule(rule *models.Rule) error {
	// Ensure target directory exists
	if err := l.EnsureTargetDirectory(); err != nil {
		return err
	}

	var targetFileName string

	// If the rule has a topic (is in a subfolder), convert path separators to underscores
	if rule.Topic != "" {
		// Convert path separators to underscores
		topicUnderscored := strings.ReplaceAll(rule.Topic, "/", "_")
		targetFileName = topicUnderscored + "_" + rule.Name + filepath.Ext(rule.Path)

		if l.Verbose {
			fmt.Printf("Converting path separators to underscores: %s -> %s\n",
				rule.Topic+"/"+rule.Name,
				targetFileName)
		}
	} else {
		// No topic, use the original filename
		targetFileName = filepath.Base(rule.Path)
	}

	// Get the rules directory based on the selected editor
	rulesDir := l.GetRulesDirectory()

	// Set the target path in the editor-specific rules directory
	targetPath := filepath.Join(l.TargetDir, rulesDir, targetFileName)

	// Check if the target already exists
	if _, err := os.Stat(targetPath); err == nil {
		// Remove existing link or file
		if l.DryRun {
			if l.Verbose {
				fmt.Printf("Would remove existing: %s\n", targetPath)
			}
		} else {
			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf("failed to remove existing rule: %w", err)
			}
		}
	}

	// Create symlink
	if l.DryRun {
		if l.Verbose {
			fmt.Printf("Would create symlink: %s -> %s\n", rule.Path, targetPath)
		}
		return nil
	}

	// Determine symlink path to use
	var symlinkPath string
	if filepath.IsAbs(rule.Path) {
		// For absolute paths, convert to relative path for the symlink
		targetDir := filepath.Dir(targetPath)
		relPath, err := filepath.Rel(targetDir, rule.Path)
		if err != nil {
			return fmt.Errorf("failed to create relative path for symlink: %w", err)
		}
		symlinkPath = relPath
	} else {
		// For paths that are already relative, use them directly
		symlinkPath = rule.Path
	}

	if l.Verbose {
		fmt.Printf("Creating symlink: %s -> %s\n", symlinkPath, targetPath)
	}

	if err := os.Symlink(symlinkPath, targetPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// LinkRules creates symlinks for all provided rules
func (l *Linker) LinkRules(rules []*models.Rule) error {
	for _, rule := range rules {
		if err := l.LinkRule(rule); err != nil {
			return err
		}
	}
	return nil
}

// UnlinkRule removes a symlink for a rule
func (l *Linker) UnlinkRule(ruleName string) error {
	// Handle rules with topic paths
	var targetFileName string

	// If the rule name includes a topic path (contains slashes)
	if strings.Contains(ruleName, "/") {
		// Convert slashes to underscores for the filename
		targetFileName = strings.ReplaceAll(ruleName, "/", "_") + ".mdc"
	} else {
		// No path separators, might be a base name or already in underscore format
		// Strip the extension if it was included
		baseName := strings.TrimSuffix(ruleName, filepath.Ext(ruleName))
		targetFileName = baseName + ".mdc"
	}

	// Get the rules directory based on the selected editor
	rulesDir := l.GetRulesDirectory()

	// Create the target path in the editor-specific rules directory
	targetPath := filepath.Join(l.TargetDir, rulesDir, targetFileName)

	// Check if the target exists
	if _, err := os.Stat(targetPath); err == nil {
		// Remove existing link or file
		if l.DryRun {
			if l.Verbose {
				fmt.Printf("Would remove: %s\n", targetPath)
			}
			return nil
		}

		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("failed to remove rule: %w", err)
		}

		return nil
	}

	// If we didn't find the file with the converted name,
	// try checking if it's a flat file without path conversion
	flatPath := filepath.Join(l.TargetDir, rulesDir, filepath.Base(strings.TrimSuffix(ruleName, filepath.Ext(ruleName)))+".mdc")
	if flatPath != targetPath {
		if _, err := os.Stat(flatPath); err == nil {
			if l.DryRun {
				if l.Verbose {
					fmt.Printf("Would remove: %s\n", flatPath)
				}
				return nil
			}

			if err := os.Remove(flatPath); err != nil {
				return fmt.Errorf("failed to remove rule: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("rule %s is not linked", ruleName)
}

// IsRuleLinked checks if a rule is already linked in the target directory
func (l *Linker) IsRuleLinked(rule *models.Rule) bool {
	var targetFileName string

	// If the rule has a topic (is in a subfolder), convert path separators to underscores
	if rule.Topic != "" {
		// Convert path separators to underscores
		topicUnderscored := strings.ReplaceAll(rule.Topic, "/", "_")
		targetFileName = topicUnderscored + "_" + rule.Name + filepath.Ext(rule.Path)
	} else {
		// No topic, use the original filename
		targetFileName = filepath.Base(rule.Path)
	}

	// Get the rules directory based on the selected editor
	rulesDir := l.GetRulesDirectory()

	// Check the target path in the editor-specific rules directory
	targetPath := filepath.Join(l.TargetDir, rulesDir, targetFileName)

	// Check if the target exists
	if _, err := os.Stat(targetPath); err == nil {
		return true
	}

	// Also check for the old-style flat path (for backward compatibility)
	flatPath := filepath.Join(l.TargetDir, rulesDir, filepath.Base(rule.Path))
	if flatPath != targetPath {
		if _, err := os.Stat(flatPath); err == nil {
			return true
		}
	}

	return false
}
