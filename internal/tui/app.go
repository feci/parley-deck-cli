package tui

// This file holds the first-run init wizard (shown by `parley tui` when the
// workspace is not yet initialized) plus the shared TUI styles and small string
// helpers used across the package. The old workspace dashboard model that used
// to live here was retired in favor of the unified tabbed Home view in live.go.

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("62")).Padding(0, 1)
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	// steerStyle (cyan) marks the steer input prefix so a message to an agent is
	// visually distinct from plain input and never sent by accident.
	steerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("44")).Bold(true)
)

type initModel struct {
	root         string
	agents       []agents.Discovery
	status       *protocol.WorkspaceStatus
	width        int
	height       int
	initializing bool
	errText      string
}

type initWorkspaceMsg struct {
	status protocol.WorkspaceStatus
	err    error
}

func RunInit(root string, discovered []agents.Discovery) error {
	program := tea.NewProgram(newInitModel(root, discovered), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func newInitModel(root string, discovered []agents.Discovery) initModel {
	return initModel{root: root, agents: discovered, width: 100}
}

func (m initModel) Init() tea.Cmd {
	return nil
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "i", "enter":
			if m.status != nil || m.initializing {
				return m, nil
			}
			m.initializing = true
			m.errText = ""
			return m, initWorkspaceCmd(m.root)
		}
	case initWorkspaceMsg:
		m.initializing = false
		if msg.err != nil {
			m.errText = msg.err.Error()
			return m, nil
		}
		// Init succeeded: quit the wizard. The caller re-reads the workspace
		// status and opens the unified Home view (live.go) with the full launch
		// path wired in, instead of the retired workspace dashboard.
		status := msg.status
		m.status = &status
		m.errText = ""
		return m, tea.Quit
	}
	return m, nil
}

func (m initModel) View() string {
	width := m.width
	if width < 80 {
		width = 80
	}
	bodyWidth := width - 4
	if bodyWidth < 76 {
		bodyWidth = 76
	}

	var body strings.Builder
	body.WriteString("Workspace is not initialized\n\n")
	body.WriteString(fmt.Sprintf("Path: %s\n\n", m.root))
	if m.initializing {
		body.WriteString(warnStyle.Render("Initializing workspace..."))
	} else {
		body.WriteString("Press i or enter to initialize Parley Deck here.")
	}
	if m.errText != "" {
		body.WriteString("\n\n")
		body.WriteString(warnStyle.Render("Init failed: " + m.errText))
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render("Parley Deck setup"))
	b.WriteString("\n\n")
	b.WriteString(boxStyle.Width(bodyWidth).Render(body.String()))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Keys: i/enter initialize  q/esc/ctrl+c quit"))
	b.WriteString("\n")
	return b.String()
}

func initWorkspaceCmd(root string) tea.Cmd {
	return func() tea.Msg {
		if err := protocol.InitWorkspace(root); err != nil {
			return initWorkspaceMsg{err: err}
		}
		status, err := protocol.ReadWorkspaceStatus(root)
		return initWorkspaceMsg{status: status, err: err}
	}
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func truncateText(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	if limit <= 3 {
		return value[:limit]
	}
	return value[:limit-3] + "..."
}
