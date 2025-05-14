package ui

import (
	"github.com/circleci/llm-agent-rules/pkg/models"
)

// item represents a rule in the list
type item struct {
	rule *models.Rule
}

// FilterValue implements list.Item
func (i item) FilterValue() string {
	return i.getRuleName()
}

// Title returns the item title
func (i item) Title() string {
	return i.getRuleName()
}

// Description returns the rule description
func (i item) Description() string {
	return i.rule.Description
}

func (i item) getRuleName() string {
	if i.rule.Topic != "" {
		return i.rule.Topic + "/" + i.rule.Name
	}
	return i.rule.Name
}
