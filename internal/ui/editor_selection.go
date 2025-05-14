package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/circleci/llm-agent-rules/internal/config"
)

// EditorOption represents an editor choice
type EditorOption struct {
	name        string
	description string
	value       config.Editor
}

// FilterValue implements list.Item interface
func (e EditorOption) FilterValue() string {
	return e.name
}

// Title returns the editor name
func (e EditorOption) Title() string {
	return e.name
}

// Description returns the editor description
func (e EditorOption) Description() string {
	return e.description
}

// EditorSelectionModel represents the editor selection UI
type EditorSelectionModel struct {
	list   list.Model
	choice config.Editor
	done   bool
	cfg    *config.Config
}

// NewEditorSelectionModel creates a new editor selection model
func NewEditorSelectionModel(cfg *config.Config) *EditorSelectionModel {
	// Define available editors
	editors := []list.Item{
		EditorOption{
			name:        "Cursor",
			description: "Use Cursor editor integration",
			value:       config.EditorCursor,
		},
		EditorOption{
			name:        "Windsurf",
			description: "Use Windsurf editor integration",
			value:       config.EditorWindsurf,
		},
	}

	// Create local style definitions
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#8A2BE2")).
		Padding(0, 1)

	paginationStyle := lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.Color("#FFFF00"))

	helpStyle := lipgloss.NewStyle().
		PaddingLeft(4).
		PaddingBottom(1).
		Foreground(lipgloss.Color("#00FFFF"))

	// Create a new list
	l := list.New(editors, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Your Editor"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &EditorSelectionModel{
		list:   l,
		choice: config.EditorCursor, // Default to Cursor
		cfg:    cfg,
	}
}

// Init initializes the model
func (m *EditorSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles user input
func (m *EditorSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// Get the selected editor
			if i, ok := m.list.SelectedItem().(EditorOption); ok {
				m.choice = i.value
			}
			m.done = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the UI
func (m *EditorSelectionModel) View() string {
	if m.done {
		return ""
	}
	// Create a document style
	docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)
	return docStyle.Render(m.list.View())
}

// GetSelectedEditor returns the selected editor
func (m *EditorSelectionModel) GetSelectedEditor() config.Editor {
	return m.choice
}

// RunEditorSelection runs the editor selection UI and returns the selected editor
func RunEditorSelection(cfg *config.Config) (config.Editor, error) {
	m := NewEditorSelectionModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return config.EditorCursor, fmt.Errorf("error running editor selection: %w", err)
	}

	return m.GetSelectedEditor(), nil
}
