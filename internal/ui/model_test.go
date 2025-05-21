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

func TestGetEditorPathName(t *testing.T) {
	testCases := []struct {
		name   string
		editor string
		want   string
	}{
		{
			name:   "Default editor",
			editor: "editor (default)",
			want:   ".editor",
		},
		{
			name:   "Capitalized name should be lowercased",
			editor: "Windsurf",
			want:   ".windsurf",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				editor: tt.editor,
			}

			if got := model.getEditorPathName(); got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
