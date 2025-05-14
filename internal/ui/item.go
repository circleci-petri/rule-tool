package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/circleci/llm-agent-rules/pkg/models"
)

// item represents a rule in the list
type item struct {
	rule *models.Rule
	// Add title and description styles directly
	titleStyle     lipgloss.Style
	descStyle      lipgloss.Style
	selectedStyle  lipgloss.Style
	checkmarkStyle lipgloss.Style
}

// FilterValue implements list.Item
func (i item) FilterValue() string {
	if i.rule.Topic != "" {
		return i.rule.Topic + "/" + i.rule.Name
	}
	return i.rule.Name
}

// Title returns the item title
func (i item) Title() string {
	if i.rule.Topic != "" {
		return i.rule.Topic + "/" + i.rule.Name
	}
	return i.rule.Name
}

// Description returns the rule description
func (i item) Description() string {
	return i.rule.Description
}
