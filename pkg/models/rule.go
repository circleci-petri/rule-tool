package models

import (
	"os"
	"path/filepath"
	"strings"
)

// Rule represents a single Cursor rule
type Rule struct {
	Name        string
	Description string
	Globs       []string
	Path        string
	Content     string
	Selected    bool
	Topic       string // Represents the subfolder/category the rule belongs to
	IsInstalled bool   // Tracks if the rule is already installed in the target
}

// NewRule creates a new Rule instance from a file path
func NewRule(path string) (*Rule, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	rule := &Rule{
		Name:     name,
		Path:     path,
		Content:  string(content),
		Selected: false,
		Topic:    "", // Default to empty topic
	}

	// Parse the content to extract description and globs
	rule.parseContent()

	return rule, nil
}

// parseContent extracts metadata from rule content
func (r *Rule) parseContent() {
	lines := strings.Split(r.Content, "\n")
	inFrontmatter := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "description:") {
				r.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			} else if strings.HasPrefix(line, "globs:") {
				globsStr := strings.TrimSpace(strings.TrimPrefix(line, "globs:"))
				if globsStr != "" {
					r.Globs = append(r.Globs, globsStr)
				}
			}
		}
	}
}
