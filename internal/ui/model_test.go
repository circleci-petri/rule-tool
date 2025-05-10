package ui

import (
	"testing"

	"github.com/circleci/llm-agent-rules/pkg/models"
)

func TestFilterIncludesTopic(t *testing.T) {
	// Create a test rule
	rule := &models.Rule{
		Name:        "testrule",
		Description: "Test rule",
		Topic:       "testfolder",
	}

	// Create an item with this rule
	testItem := item{
		rule: rule,
	}

	// Verify FilterValue includes the topic
	expected := "testfolder/testrule"
	if actual := testItem.FilterValue(); actual != expected {
		t.Errorf("FilterValue() = %q, want %q", actual, expected)
	}

	// Test rule without a topic
	rule.Topic = ""
	expected = "testrule"
	if actual := testItem.FilterValue(); actual != expected {
		t.Errorf("FilterValue() = %q, want %q", actual, expected)
	}
}
