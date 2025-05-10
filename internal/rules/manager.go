package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/circleci/llm-agent-rules/pkg/models"
)

// Manager handles loading and managing rules
type Manager struct {
	Rules     []*models.Rule
	RulesPath string
}

// NewManager creates a new rules manager
func NewManager(rulesPath string) *Manager {
	return &Manager{
		Rules:     make([]*models.Rule, 0),
		RulesPath: rulesPath,
	}
}

// LoadRules loads all rules from the rules directory
func (m *Manager) LoadRules() error {
	// Clear existing rules
	m.Rules = make([]*models.Rule, 0)

	// Walk through the rules directory
	err := filepath.Walk(m.RulesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-markdown files
		if filepath.Ext(path) != ".mdc" {
			return nil
		}

		// Create a new rule from the file
		rule, err := models.NewRule(path)
		if err != nil {
			return err
		}

		// If the rule is in a subfolder, prepend the folder name to the display name
		relPath, err := filepath.Rel(m.RulesPath, path)
		if err == nil && filepath.Dir(relPath) != "." {
			// Get the folder structure
			folder := filepath.Dir(relPath)
			// Replace backslashes with forward slashes for consistency
			folder = strings.ReplaceAll(folder, "\\", "/")
			// Set the topic to the folder name
			rule.Topic = folder
		}

		// Add the rule to the list
		m.Rules = append(m.Rules, rule)
		return nil
	})

	return err
}

// GetRuleByName returns a rule by its name
func (m *Manager) GetRuleByName(name string) *models.Rule {
	for _, rule := range m.Rules {
		// Check both the plain name and the topic/name format
		if rule.Name == name {
			return rule
		}

		if rule.Topic != "" && rule.Topic+"/"+rule.Name == name {
			return rule
		}
	}
	return nil
}

// GetSelectedRules returns all selected rules
func (m *Manager) GetSelectedRules() []*models.Rule {
	selected := make([]*models.Rule, 0)
	for _, rule := range m.Rules {
		if rule.Selected {
			selected = append(selected, rule)
		}
	}
	return selected
}
