package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/circleci/llm-agent-rules/internal/ui/components"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// Manager handles the management of modal overlays in the application
type Manager struct {
	currentModal tea.Model
	background   tea.Model
	overlay      *overlay.Model
	showModal    bool
}

// NewManager creates a new Manager instance
func NewManager(background tea.Model) *Manager {
	return &Manager{
		background: background,
		overlay: overlay.New(
			nil,
			background,
			overlay.Center,
			overlay.Center,
			0,
			0,
		),
		showModal: false,
	}
}

// Init initializes the Manager
func (m *Manager) Init() tea.Cmd {
	return nil
}

// Update handles updates for the Manager
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showModal && m.currentModal != nil {
		fg, fgCmd := m.currentModal.Update(msg)
		m.currentModal = fg

		if keyMsg, ok := msg.(tea.KeyMsg); ok && components.IgnoreKeyType(keyMsg.Type) {
			return m, fgCmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			if !m.showModal {
				m.currentModal = components.NewEditorModal()
				m.overlay.Foreground = m.currentModal
			}
			m.showModal = !m.showModal
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case components.CloseModalMsg:
		m.showModal = false
		return m, nil
	}

	bg, bgCmd := m.background.Update(msg)
	m.background = bg

	return m, bgCmd
}

// View renders the Manager's view
func (m *Manager) View() string {
	if m.showModal {
		// Return the appropriate view based on the current state
		return m.overlay.View()
	}
	return m.background.View()
}
