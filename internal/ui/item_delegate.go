package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom delegate for item rendering
type itemDelegate struct {
	styles struct {
		NormalTitle   lipgloss.Style
		NormalDesc    lipgloss.Style
		SelectedTitle lipgloss.Style
		SelectedDesc  lipgloss.Style
		CheckMark     lipgloss.Style
	}
}

func newItemDelegate() itemDelegate {
	d := itemDelegate{}

	d.styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")). // Hot pink
		Bold(true)

	d.styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")) // Sky blue

	d.styles.SelectedTitle = d.styles.NormalTitle.
		Background(lipgloss.Color("#4B0082")). // Indigo
		Foreground(lipgloss.Color("#FFFFFF"))

	d.styles.SelectedDesc = d.styles.NormalDesc.
		Foreground(lipgloss.Color("#FFD700")) // Gold

	d.styles.CheckMark = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")) // Bright green

	return d
}

func (d itemDelegate) Height() int                               { return 2 } // Each item takes up 2 lines (title + description)
func (d itemDelegate) Spacing() int                              { return 1 } // 1 line of spacing between items
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	selected := index == m.Index()
	rule := i.rule

	var title, desc string

	// Format the display name to include topic if present
	displayName := rule.Name
	if rule.Topic != "" {
		displayName = rule.Topic + "/" + rule.Name
	}

	// Add indentation to align with header
	indent := "    "

	if selected {
		title = indent + d.styles.SelectedTitle.Render(displayName)
		desc = indent + d.styles.SelectedDesc.Render(rule.Description)
	} else {
		title = indent + d.styles.NormalTitle.Render(displayName)
		desc = indent + d.styles.NormalDesc.Render(rule.Description)
	}

	// Add appropriate indicator based on rule status
	if rule.IsInstalled {
		title = title + " [INSTALLED]"
	} else if rule.Selected {
		title = title + " ✓"
	}

	_, _ = fmt.Fprintf(w, "%s\n%s", title, desc)
}
