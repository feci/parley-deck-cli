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
	status          protocol.WorkspaceStatus
	agents          []agents.Discovery
	width           int
	focus           focusZone
	selectedIdea    int
	selectedAgent   int
	launchOverrides map[string]string
}

type focusZone string

const (
	focusIdeas  focusZone = "ideas"
	focusAgents focusZone = "agents"
)

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
	program := tea.NewProgram(newModel(status, discovered), tea.WithAltScreen())
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

func newModel(status protocol.WorkspaceStatus, discovered []agents.Discovery) model {
	m := model{
		status:          status,
		agents:          discovered,
		width:           100,
		focus:           focusAgents,
		launchOverrides: map[string]string{},
	}
	if len(discovered) == 0 {
		m.focus = focusIdeas
	}
	m.clampSelections()
	return m
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
		case "tab", "shift+tab":
			m.switchFocus()
			return m, nil
		case "j", "down":
			m.moveSelection(1)
			return m, nil
		case "k", "up":
			m.moveSelection(-1)
			return m, nil
		case "h":
			m.setSelectedAgentMode(agents.LaunchHeadless)
			return m, nil
		case "i":
			m.setSelectedAgentMode(agents.LaunchInteractive)
			return m, nil
		case "m":
			m.setSelectedAgentMode(agents.LaunchManual)
			return m, nil
		case "x":
			m.clearSelectedAgentMode()
			return m, nil
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

	left := panelStyle(m.focus == focusIdeas).Width(leftWidth).Render(m.renderIdeas())
	right := panelStyle(m.focus == focusAgents).Width(rightWidth).Render(m.renderAgents())
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
		dashboard := newModel(status, m.agents)
		dashboard.width = m.width
		return dashboard, nil
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

func panelStyle(focused bool) lipgloss.Style {
	style := boxStyle
	if focused {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return style
}

func (m model) renderIdeas() string {
	var b strings.Builder
	b.WriteString("Protocol\n")
	if len(m.status.Ideas) == 0 {
		b.WriteString(mutedStyle.Render("no ideas found"))
		return b.String()
	}
	for i, idea := range m.status.Ideas {
		marker := m.selectionMarker(focusIdeas, i == m.selectedIdea)
		b.WriteString(fmt.Sprintf("%s %s\n", marker, idea.Slug))
		b.WriteString(fmt.Sprintf("   status: %s\n", idea.Status))
		b.WriteString(fmt.Sprintf("   participants: %s\n", strings.Join(idea.Participants, ", ")))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderAgents() string {
	var b strings.Builder
	b.WriteString("Agents\n")
	if len(m.agents) == 0 {
		b.WriteString(mutedStyle.Render("no agents configured"))
		return b.String()
	}
	b.WriteString("  ID       STATE      MODE          VERSION\n")
	for i, agent := range m.agents {
		marker := m.selectionMarker(focusAgents, i == m.selectedAgent)
		state := warnStyle.Render("missing")
		version := "-"
		if agent.Found {
			state = okStyle.Render("found")
			version = agent.Version
			if version == "" {
				version = "unknown"
			}
		}
		mode := m.effectiveLaunchMode(agent)
		if m.hasLaunchOverride(agent.ID) {
			mode += "*"
		}
		b.WriteString(fmt.Sprintf("%s %-8s %-10s %-13s %s\n", marker, agent.ID, state, mode, version))
		if agent.Notes != "" {
			b.WriteString(mutedStyle.Render("  "+agent.Notes) + "\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(m.renderAgentDetails())
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderAgentDetails() string {
	agent, ok := m.selectedAgentDiscovery()
	if !ok {
		return mutedStyle.Render("Agent details unavailable")
	}

	configuredMode := agents.LaunchModeOrDefault(agent.LaunchMode)
	effectiveMode := m.effectiveLaunchMode(agent)
	modeLine := fmt.Sprintf("effective: %s", effectiveMode)
	if m.hasLaunchOverride(agent.ID) {
		modeLine += " (session only)"
	}

	var b strings.Builder
	b.WriteString("Agent details\n")
	b.WriteString(fmt.Sprintf("id: %s\n", agent.ID))
	b.WriteString(fmt.Sprintf("installed: %s\n", yesNo(agent.Found)))
	if agent.Version != "" {
		b.WriteString(fmt.Sprintf("version: %s\n", agent.Version))
	}
	if agent.Error != "" {
		b.WriteString(fmt.Sprintf("probe error: %s\n", agent.Error))
	}
	b.WriteString(fmt.Sprintf("configured launch: %s\n", configuredMode))
	b.WriteString(modeLine + "\n")
	b.WriteString(fmt.Sprintf("model: %s  reasoning/profile: %s\n", valueOrDefault(agent.Model), valueOrDefault(firstNonEmpty(agent.Reasoning, agent.Profile))))
	b.WriteString(fmt.Sprintf("sandbox: %s  approval: %s\n", valueOrDefault(agent.SandboxMode), valueOrDefault(agent.ApprovalPolicy)))
	b.WriteString(fmt.Sprintf("timeout: %dms  backend: %s  isolated home: %s\n", timeoutMS(agent), backendOrUnknown(agent.ExternalBackend), yesNo(agent.IsolateHome)))
	b.WriteString(fmt.Sprintf("headless: %s\n", headlessCommandLine(agent)))
	b.WriteString(fmt.Sprintf("interactive: %s\n", interactiveCommandLine(agent)))
	b.WriteString(fmt.Sprintf("prompt: %s  invoke: %s\n", agents.InteractivePromptModeOrDefault(agent.InteractivePromptMode), agents.InteractiveInvokeOrDefault(agent.InteractiveInvoke)))
	if agent.InteractiveNotes != "" {
		b.WriteString(fmt.Sprintf("interactive notes: %s\n", agent.InteractiveNotes))
	}
	if agent.Notes != "" {
		b.WriteString(fmt.Sprintf("notes: %s\n", agent.Notes))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderFooter() string {
	return "Keys: tab/shift+tab focus  j/k/up/down select  h/i/m set agent mode  x clear mode  q/esc/ctrl+c quit.  Mode overrides are session-only preview; this dashboard does not launch agents."
}

func (m model) selectionMarker(zone focusZone, selected bool) string {
	if !selected {
		return " "
	}
	if m.focus == zone {
		return ">"
	}
	return "."
}

func (m *model) switchFocus() {
	if m.focus == focusIdeas && len(m.agents) > 0 {
		m.focus = focusAgents
		return
	}
	m.focus = focusIdeas
}

func (m *model) moveSelection(delta int) {
	switch m.focus {
	case focusIdeas:
		m.selectedIdea = clampIndex(m.selectedIdea+delta, len(m.status.Ideas))
	case focusAgents:
		m.selectedAgent = clampIndex(m.selectedAgent+delta, len(m.agents))
	}
}

func (m *model) clampSelections() {
	m.selectedIdea = clampIndex(m.selectedIdea, len(m.status.Ideas))
	m.selectedAgent = clampIndex(m.selectedAgent, len(m.agents))
}

func clampIndex(index, length int) int {
	if length <= 0 {
		return 0
	}
	if index < 0 {
		return 0
	}
	if index >= length {
		return length - 1
	}
	return index
}

func (m model) selectedAgentDiscovery() (agents.Discovery, bool) {
	if len(m.agents) == 0 {
		return agents.Discovery{}, false
	}
	index := clampIndex(m.selectedAgent, len(m.agents))
	return m.agents[index], true
}

func (m model) effectiveLaunchMode(agent agents.Discovery) string {
	if override, ok := m.launchOverrides[agent.ID]; ok {
		return override
	}
	return agents.LaunchModeOrDefault(agent.LaunchMode)
}

func (m model) hasLaunchOverride(agentID string) bool {
	_, ok := m.launchOverrides[agentID]
	return ok
}

func (m *model) setSelectedAgentMode(mode string) {
	if m.focus != focusAgents {
		return
	}
	agent, ok := m.selectedAgentDiscovery()
	if !ok {
		return
	}
	if m.launchOverrides == nil {
		m.launchOverrides = map[string]string{}
	}
	m.launchOverrides[agent.ID] = mode
}

func (m *model) clearSelectedAgentMode() {
	if m.focus != focusAgents {
		return
	}
	agent, ok := m.selectedAgentDiscovery()
	if !ok || m.launchOverrides == nil {
		return
	}
	delete(m.launchOverrides, agent.ID)
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func valueOrDefault(value string) string {
	if strings.TrimSpace(value) == "" {
		return agents.CLIDefault
	}
	return value
}

func backendOrUnknown(value string) string {
	if strings.TrimSpace(value) == "" {
		return agents.ExternalUnknown
	}
	return value
}

func headlessCommandLine(agent agents.Discovery) string {
	switch {
	case agent.HeadlessMode != "":
		return agent.HeadlessMode
	case len(agent.HeadlessArgs) > 0:
		return strings.Join(agent.HeadlessArgs, " ")
	case len(agent.Commands) > 0:
		return agent.Commands[0]
	default:
		return agents.CLIDefault
	}
}

func interactiveCommandLine(agent agents.Discovery) string {
	command := valueOrDefault(agents.InteractiveCommandOrDefault(agent.Spec))
	if len(agent.InteractiveArgs) == 0 {
		return command
	}
	return command + " " + strings.Join(agent.InteractiveArgs, " ")
}

func timeoutMS(agent agents.Discovery) int {
	if agent.TimeoutMS > 0 {
		return agent.TimeoutMS
	}
	return agents.DefaultTimeoutMS
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
