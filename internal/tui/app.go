package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
)

type model struct {
	status protocol.WorkspaceStatus
	agents []agents.Discovery
	width  int
}

type initModel struct {
	root         string
	agents       []agents.Discovery
	status       *protocol.WorkspaceStatus
	width        int
	initializing bool
	errText      string
}

type initWorkspaceMsg struct {
	status protocol.WorkspaceStatus
	err    error
}

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("62")).Padding(0, 1)
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

func Run(status protocol.WorkspaceStatus, discovered []agents.Discovery) error {
	program := tea.NewProgram(model{status: status, agents: discovered}, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func RunInit(root string, discovered []agents.Discovery) error {
	program := tea.NewProgram(newInitModel(root, discovered), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func newInitModel(root string, discovered []agents.Discovery) initModel {
	return initModel{root: root, agents: discovered, width: 100}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf("Parley Deck  transport=%s  mode=HITL default", m.status.Transport)))
	b.WriteString("\n\n")

	width := m.width
	if width < 80 {
		width = 80
	}
	leftWidth := width/2 - 4
	rightWidth := width - leftWidth - 8
	if leftWidth < 34 {
		leftWidth = 34
	}
	if rightWidth < 34 {
		rightWidth = 34
	}

	left := boxStyle.Width(leftWidth).Render(m.renderIdeas())
	right := boxStyle.Width(rightWidth).Render(m.renderAgents())
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right))
	b.WriteString("\n\n")
	b.WriteString(boxStyle.Width(width - 4).Render(m.renderFooter()))
	b.WriteString("\n")
	return b.String()
}

func (m initModel) Init() tea.Cmd {
	return nil
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
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
		status := msg.status
		m.status = &status
		m.errText = ""
		return m, nil
	}
	return m, nil
}

func (m initModel) View() string {
	if m.status != nil {
		return model{status: *m.status, agents: m.agents, width: m.width}.View()
	}

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

func (m model) renderIdeas() string {
	var b strings.Builder
	b.WriteString("Protocol\n")
	if len(m.status.Ideas) == 0 {
		b.WriteString(mutedStyle.Render("no ideas found"))
		return b.String()
	}
	for _, idea := range m.status.Ideas {
		b.WriteString(fmt.Sprintf("%s  %s\n", okStyle.Render("●"), idea.Slug))
		b.WriteString(fmt.Sprintf("   status: %s\n", idea.Status))
		b.WriteString(fmt.Sprintf("   participants: %s\n", strings.Join(idea.Participants, ", ")))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderAgents() string {
	var b strings.Builder
	b.WriteString("Agents\n")
	b.WriteString("ID       STATE      VERSION\n")
	for _, agent := range m.agents {
		state := warnStyle.Render("missing")
		version := "-"
		if agent.Found {
			state = okStyle.Render("found")
			version = agent.Version
			if version == "" {
				version = "unknown"
			}
		}
		b.WriteString(fmt.Sprintf("%-8s %-10s %s\n", agent.ID, state, version))
		if agent.Notes != "" {
			b.WriteString(mutedStyle.Render("  "+agent.Notes) + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderFooter() string {
	return "Timeline: kickoff > rounds > consensus > implementation > review.  Keys: q quit.  Next slice: runner + live event streaming."
}
