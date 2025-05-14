package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/circleci/llm-agent-rules/internal/config"
	"github.com/circleci/llm-agent-rules/internal/linker"
	"github.com/circleci/llm-agent-rules/internal/rules"
	"github.com/circleci/llm-agent-rules/pkg/models"
)

// Define styles with vibrant colors
var (
	// UI element styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#8A2BE2")).
			Padding(0, 1)

	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4).
			Foreground(lipgloss.Color("#FFFF00"))

	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(lipgloss.Color("#00FFFF"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1).
			Bold(true)
)

// Model represents the UI model
type Model struct {
	list           list.Model
	rulesManager   *rules.Manager
	linker         *linker.Linker
	config         *config.Config
	selectedRule   *models.Rule
	err            error
	width          int
	height         int
	successMessage string
	showingSuccess bool
	successTimer   int
}

// New creates a new UI model
func New(cfg *config.Config, rulesManager *rules.Manager, linker *linker.Linker) *Model {
	// Convert rules to list items with styles
	items := []list.Item{}

	// Define our item styles
	normalTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")). // Hot pink for visibility
		Bold(true)

	normalDescStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")) // Sky blue

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")) // Gold

	checkmarkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")) // Bright green

	for _, rule := range rulesManager.Rules {
		items = append(items, item{
			rule:           rule,
			titleStyle:     normalTitleStyle,
			descStyle:      normalDescStyle,
			selectedStyle:  selectedStyle,
			checkmarkStyle: checkmarkStyle,
		})
	}

	// Create custom delegate
	delegate := newItemDelegate()

	// Create the list with custom styling
	l := list.New(items, delegate, 20, 20) // Start with reasonable defaults
	l.Title = "Available Rules"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	return &Model{
		list:           l,
		rulesManager:   rulesManager,
		linker:         linker,
		config:         cfg,
		width:          80, // Default width
		height:         24, // Default height
		successMessage: "",
		showingSuccess: false,
		successTimer:   0,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Request initial window size
	return nil
}

// Update handles user input and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.setListHeight(m.height)
		m.list.SetWidth(m.width)
		return m, nil

	case tea.KeyMsg:
		// If we're showing success message, clear it on any key press
		if m.showingSuccess {
			m.showingSuccess = false
			m.successMessage = ""
			return m, nil
		}

		// Skip hotkey handling if we're currently setting a filter
		if m.list.SettingFilter() {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selectedRule = i.rule
				m.selectedRule.Selected = !m.selectedRule.Selected
				return m, nil
			}

		case "a":
			// Select all *visible* rules (respecting filter)
			visibleItems := m.list.VisibleItems()
			for _, listItem := range visibleItems {
				if i, ok := listItem.(item); ok {
					i.rule.Selected = true
				}
			}
			return m, nil

		case "d":
			// Deselect all *visible* rules (respecting filter)
			visibleItems := m.list.VisibleItems()
			for _, listItem := range visibleItems {
				if i, ok := listItem.(item); ok {
					i.rule.Selected = false
				}
			}
			return m, nil

		case "l":
			// Link selected rules
			selected := m.rulesManager.GetSelectedRules()
			if len(selected) > 0 {
				err := m.linker.LinkRules(selected)
				if err != nil {
					m.err = err
				} else {
					// Update installation status for the newly linked rules
					for _, rule := range selected {
						rule.IsInstalled = true
					}

					// Show success message
					m.successMessage = fmt.Sprintf("✓ Successfully linked %d rules!", len(selected))
					m.showingSuccess = true

					// Clear success message after a short delay
					return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
						return tickMsg{}
					})
				}
			}
			return m, nil
		}
	}

	// Handle timer tick for clearing success message
	if _, ok := msg.(tickMsg); ok {
		m.showingSuccess = false
		m.successMessage = ""
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Add a tick message type for handling the timer
type tickMsg struct{}

// View renders the UI
func (m *Model) View() string {
	m.setListHeight(m.height)

	// Create header with title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		PaddingRight(2).
		Width(m.width)

	headerTitle := titleStyle.Render("Rule Tool CLI")

	// Status section
	var status string

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(0, 1)
		status = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	} else if m.showingSuccess && m.successMessage != "" {
		// Show success message with a highlighted style
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(0, 1)
		status = successStyle.Render(m.successMessage)
	} else {
		status = statusStyle.Render(m.updateStatusText())
	}

	// Get the list view (main content)
	listView := m.list.View()

	// Calculate widths for the bottom panels
	bottomWidth := max(m.width-4, 40)

	leftWidth := bottomWidth / 2
	rightWidth := bottomWidth - leftWidth

	// Style for the help/controls section (left side)
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666666")).
		Padding(1, 2).
		Width(leftWidth)

	// Style for the repository info section (right side)
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0000AA")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666666")).
		Padding(1, 3, 1, 3).
		Width(rightWidth)

	// Get paths for display and ensure they're not empty
	rulesRepoPath := m.config.RulesRepoPath
	if rulesRepoPath == "" {
		rulesRepoPath = "Not set"
	}

	targetPath := m.config.TargetProjectPath
	if targetPath == "" {
		targetPath = "Not set"
	}

	// Create content for both bottom panels
	helpContent := "Controls:\n" +
		"• Enter: Toggle selection\n" +
		"• a: Select all\n" +
		"• d: Deselect all\n" +
		"• l: Link selected rules\n" +
		"• /: Filter rules\n" +
		"• q: Quit"

	infoContent := "Repository Info:\n" +
		"• Rules: " + rulesRepoPath + "\n" +
		"• Target: " + targetPath + "\n\n" +
		"Indicators:\n" +
		"• [INSTALLED]: Rule is already installed\n" +
		"• ✓: Rule is selected for installation"

	// Render both panels
	helpSection := helpStyle.Render(helpContent)
	infoSection := infoStyle.Render(infoContent)

	// Join the bottom panels horizontally
	bottomSection := lipgloss.JoinHorizontal(lipgloss.Top, helpSection, infoSection)

	// Build the UI by explicitly stacking the components with fixed spacing
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerTitle,   // Header
		"\n",          // Empty line for spacing
		listView,      // List view (main content)
		status,        // Status bar
		bottomSection, // Bottom help section
	)

	return mainContent
}

// UpdateRuleInstallStatus updates the status line to show installed vs. selected counts
func (m *Model) updateStatusText() string {
	installedCount := 0
	newlySelectedCount := 0

	for _, rule := range m.rulesManager.Rules {
		// Count installed rules
		if rule.IsInstalled {
			installedCount++
		}

		// If the rule is selected but not already installed, count it as newly selected
		if rule.Selected && !rule.IsInstalled {
			newlySelectedCount++
		}
	}

	return fmt.Sprintf("%d rules already installed • %d new rules selected",
		installedCount, newlySelectedCount)
}

func (m *Model) setListHeight(height int) {
	// Reserve space for header (1 line)
	headerHeight := 1
	// Reserve space for status bar (1 line)
	statusHeight := 1
	// Reserve space for bottom section (about 10 lines)
	bottomSectionHeight := 6
	// Calculate remaining space for list view
	listHeight := max(height-headerHeight-statusHeight-bottomSectionHeight-4, 5)

	// Set the list height dynamically
	m.list.SetHeight(listHeight)
}
