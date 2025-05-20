package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

func NewEditorModal() tea.Model {
	return NewModal(
		"Select your editor",
		"Please select your preferred editor:",
		[]string{"Cursor", "Windsurf"},
	)
}
