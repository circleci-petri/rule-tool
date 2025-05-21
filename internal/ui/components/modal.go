package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("170")).
			Padding(1, 2).
			Width(50)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#5f5fd7")).
			Padding(0, 1)

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#d75fd7"))
)

type CloseModalMsg string

func CloseModalCmd(selected string) tea.Cmd {
	return func() tea.Msg {
		return CloseModalMsg(selected)
	}
}

type Modal struct {
	width   int
	height  int
	title   string
	message string
	buttons []string
	active  int
}

// NewModal creates a new modal with the given title and message,
// and buttons.
func NewModal(title, message string, buttons []string) *Modal {
	if len(buttons) == 0 || buttons == nil {
		buttons = []string{"Yes", "No"}
	}

	return &Modal{
		title:   title,
		message: message,
		buttons: buttons,
	}
}

func IgnoreKeyType(keyType tea.KeyType) bool {
	return keyType == tea.KeyEnter || keyType == tea.KeyUp || keyType == tea.KeyDown || keyType == tea.KeyLeft || keyType == tea.KeyRight
}

func (m *Modal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Modal) Next() {
	m.active++
	if m.active >= len(m.buttons) {
		m.active = 0
	}
}

func (m *Modal) Prev() {
	if m.active <= 0 {
		m.active = len(m.buttons) - 1
	} else {
		m.active--
	}
}

func (m *Modal) Init() tea.Cmd {
	return nil
}

func (m *Modal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, CloseModalCmd(m.buttons[m.active])
		case "up", "left":
			m.Prev()
			return m, nil
		case "down", "right":
			m.Next()
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Modal) View() string {
	var buttonRow strings.Builder

	for i, button := range m.buttons {
		if i == m.active {
			buttonRow.WriteString(activeButtonStyle.Render("• " + button))
		} else {
			buttonRow.WriteString(buttonStyle.Render("  " + button))
		}

		if i < len(m.buttons)-1 {
			buttonRow.WriteString("\n")
		}
	}

	// Layout modal content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.message,
		"",
		buttonRow.String(),
	)

	// Apply modal style
	modal := modalStyle.
		Width(m.width).
		Render(content)

	// Add title if provided
	if m.title != "" {
		modal = lipgloss.JoinVertical(
			lipgloss.Center,
			m.title,
			modal,
		)
	}

	return modal
}
