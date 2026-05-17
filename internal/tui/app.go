package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runstate"
)

type WorkspaceOptions struct {
	Root        string
	Status      protocol.WorkspaceStatus
	Agents      []agents.Discovery
	Runs        []runstate.RunSummary
	Context     context.Context
	RefreshRuns func() ([]runstate.RunSummary, error)
	StartRun    StartRunFunc
}

type StartRequest struct {
	Task            string
	Participants    []string
	LaunchOverrides map[string]string
	Auto            bool
}

type StartRunFunc func(context.Context, StartRequest) (runstate.RunSummary, error)

type model struct {
	root            string
	status          protocol.WorkspaceStatus
	agents          []agents.Discovery
	runs            []runstate.RunSummary
	ctx             context.Context
	refreshRuns     func() ([]runstate.RunSummary, error)
	startRun        StartRunFunc
	width           int
	focus           focusZone
	selectedIdea    int
	selectedAgent   int
	launchOverrides map[string]string
	startMode       bool
	starting        bool
	startText       string
	startErr        string
	errText         string
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

type refreshRunsMsg struct {
	runs []runstate.RunSummary
	err  error
}

type startRunMsg struct {
	run runstate.RunSummary
	err error
}

type refreshTickMsg time.Time

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("62")).Padding(0, 1)
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

func Run(status protocol.WorkspaceStatus, discovered []agents.Discovery) error {
	program := tea.NewProgram(newModel(WorkspaceOptions{Root: status.Root, Status: status, Agents: discovered}), tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func RunWorkspace(opts WorkspaceOptions) error {
	program := tea.NewProgram(newModel(opts), tea.WithAltScreen())
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

func newModel(opts WorkspaceOptions) model {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}
	root := opts.Root
	if root == "" {
		root = opts.Status.Root
	}
	m := model{
		root:            root,
		status:          opts.Status,
		agents:          opts.Agents,
		runs:            append([]runstate.RunSummary(nil), opts.Runs...),
		ctx:             ctx,
		refreshRuns:     opts.RefreshRuns,
		startRun:        opts.StartRun,
		width:           100,
		focus:           focusIdeas,
		launchOverrides: map[string]string{},
	}
	if len(m.runs) == 0 && len(m.agents) > 0 {
		m.focus = focusAgents
	}
	m.clampSelections()
	return m
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.refreshRuns != nil {
		cmds = append(cmds, refreshRunsCmd(m.refreshRuns), refreshTickCmd())
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		if m.startMode {
			return m.updateStartMode(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab":
			m.switchFocus()
			return m, nil
		case "N":
			if m.startRun == nil {
				m.startErr = "starting runs is not available in this view"
				return m, nil
			}
			m.startMode = true
			m.startText = ""
			m.startErr = ""
			return m, nil
		case "r":
			if m.refreshRuns == nil {
				return m, nil
			}
			return m, refreshRunsCmd(m.refreshRuns)
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
	case refreshRunsMsg:
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.errText = ""
			m.runs = msg.runs
			m.clampSelections()
		}
		return m, nil
	case refreshTickMsg:
		if m.refreshRuns == nil {
			return m, nil
		}
		return m, tea.Batch(refreshRunsCmd(m.refreshRuns), refreshTickCmd())
	case startRunMsg:
		m.starting = false
		if msg.err != nil {
			m.startErr = msg.err.Error()
			return m, nil
		}
		m.startMode = false
		m.startText = ""
		m.startErr = ""
		m.runs = upsertRunSummary(m.runs, msg.run)
		m.selectedIdea = indexRun(m.runs, msg.run.RunID)
		m.focus = focusIdeas
		return m, nil
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
	leftWidth := width/3 - 2
	midWidth := width/3 - 2
	rightWidth := width - leftWidth - midWidth - 10
	if leftWidth < 34 {
		leftWidth = 34
	}
	if midWidth < 34 {
		midWidth = 34
	}
	if rightWidth < 34 {
		rightWidth = 34
	}

	left := panelStyle(m.focus == focusIdeas).Width(leftWidth).Render(m.renderIdeas())
	middle := boxStyle.Width(midWidth).Render(m.renderEvents())
	right := panelStyle(m.focus == focusAgents).Width(rightWidth).Render(m.renderRunDetails())
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", middle, "  ", right))
	b.WriteString("\n\n")
	b.WriteString(boxStyle.Width(width - 4).Render(m.renderQuestions()))
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
		dashboard := newModel(WorkspaceOptions{Root: m.root, Status: status, Agents: m.agents})
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
	b.WriteString("Sessions\n")
	if len(m.runs) == 0 {
		b.WriteString(mutedStyle.Render("no runs found"))
		if len(m.status.Ideas) > 0 {
			b.WriteString("\n\n")
			b.WriteString(mutedStyle.Render(fmt.Sprintf("%d protocol idea(s) on disk", len(m.status.Ideas))))
		}
		return b.String()
	}
	for i, run := range m.runs {
		marker := m.selectionMarker(focusIdeas, i == m.selectedIdea)
		badge := valueOr(run.Attention, runstate.Attention(run))
		b.WriteString(fmt.Sprintf("%s %-7s %s\n", marker, badge, run.IdeaSlug))
		b.WriteString(fmt.Sprintf("   run: %s\n", run.RunID))
		b.WriteString(fmt.Sprintf("   agents: %s  q: %d\n", agentProgress(run), run.OpenQuestions))
		if !run.LastEventAt.IsZero() {
			b.WriteString(fmt.Sprintf("   last: %s ago\n", formatDuration(run.LastEventAge)))
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderEvents() string {
	run, ok := m.selectedRun()
	if !ok {
		return "Event stream\n" + mutedStyle.Render("select or start a run")
	}
	var b strings.Builder
	b.WriteString("Event stream\n")
	if len(run.State.Recent) == 0 {
		b.WriteString(mutedStyle.Render("waiting for events"))
		return b.String()
	}
	for _, event := range run.State.Recent {
		b.WriteString(fmt.Sprintf("%s  %-16s %s\n", event.Time.Local().Format("15:04:05"), event.Type, event.Text))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderRunDetails() string {
	run, ok := m.selectedRun()
	if !ok {
		return m.renderAgents()
	}
	var b strings.Builder
	b.WriteString("Run details\n")
	b.WriteString(fmt.Sprintf("idea: %s\n", valueOr(run.IdeaSlug, "unknown")))
	b.WriteString(fmt.Sprintf("run: %s\n", run.RunID))
	b.WriteString(fmt.Sprintf("attention: %s\n", valueOr(run.Attention, runstate.Attention(run))))
	b.WriteString(fmt.Sprintf("state: %s\n", displayRunState(run)))
	if run.Task != "" {
		b.WriteString(fmt.Sprintf("task: %s\n", truncateText(run.Task, 72)))
	}
	if len(run.State.Agents) == 0 {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render("no agent state yet"))
		return strings.TrimRight(b.String(), "\n")
	}
	agent := run.State.Agents[clampIndex(m.selectedAgent, len(run.State.Agents))]
	b.WriteString("\nAgents\n")
	for i, state := range run.State.Agents {
		marker := " "
		if i == m.selectedAgent {
			marker = ">"
		}
		b.WriteString(fmt.Sprintf("%s %-8s %-10s %s\n", marker, state.ID, state.State, state.LatestEvent))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Agent: %s\n", agent.ID))
	b.WriteString(fmt.Sprintf("duration: %s\n", formatDuration(agentDuration(agent))))
	if agent.ArtifactPath != "" {
		b.WriteString(fmt.Sprintf("artifact: %s\n", agent.ArtifactPath))
	}
	if agent.Error != "" {
		b.WriteString(warnStyle.Render("error: "+agent.Error) + "\n")
	}
	if preview := logPreview(agent); preview != "" {
		b.WriteString("\n")
		b.WriteString("Log tail\n")
		b.WriteString(preview)
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderQuestions() string {
	run, ok := m.selectedRun()
	if !ok {
		if m.startMode {
			return m.renderStartBox()
		}
		return "Questions\n" + mutedStyle.Render("no run selected")
	}
	var b strings.Builder
	if m.startMode {
		return m.renderStartBox()
	}
	b.WriteString("Questions\n")
	if len(run.Questions) == 0 {
		b.WriteString(mutedStyle.Render("no questions"))
		return b.String()
	}
	for _, question := range run.Questions {
		marker := " "
		if question.Status == hitl.StatusOpen {
			marker = "!"
		}
		b.WriteString(fmt.Sprintf("%s %-28s %-13s %-6s %s\n", marker, question.ID, question.Status, question.Risk, truncateText(question.Prompt, 90)))
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
	var parts []string
	if m.startMode {
		parts = append(parts, "Start: type task  enter launch  esc cancel")
	} else {
		parts = append(parts, "Keys: N new idea  r refresh  tab focus  j/k select  h/i/m set agent mode  x clear mode  q/esc quit")
		parts = append(parts, "Quit closes the TUI and cancels TUI-started runs; it does not detach child processes.")
	}
	if m.errText != "" {
		parts = append([]string{warnStyle.Render(m.errText)}, parts...)
	}
	if m.startErr != "" {
		parts = append([]string{warnStyle.Render(m.startErr)}, parts...)
	}
	return strings.Join(parts, "\n")
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
		m.selectedIdea = clampIndex(m.selectedIdea+delta, len(m.runs))
		m.selectedAgent = 0
	case focusAgents:
		if run, ok := m.selectedRun(); ok && len(run.State.Agents) > 0 {
			m.selectedAgent = clampIndex(m.selectedAgent+delta, len(run.State.Agents))
			return
		}
		m.selectedAgent = clampIndex(m.selectedAgent+delta, len(m.agents))
	}
}

func (m *model) clampSelections() {
	m.selectedIdea = clampIndex(m.selectedIdea, len(m.runs))
	if run, ok := m.selectedRun(); ok && len(run.State.Agents) > 0 {
		m.selectedAgent = clampIndex(m.selectedAgent, len(run.State.Agents))
	} else {
		m.selectedAgent = clampIndex(m.selectedAgent, len(m.agents))
	}
}

func (m model) selectedRun() (runstate.RunSummary, bool) {
	if len(m.runs) == 0 {
		return runstate.RunSummary{}, false
	}
	return m.runs[clampIndex(m.selectedIdea, len(m.runs))], true
}

func (m model) renderStartBox() string {
	var b strings.Builder
	b.WriteString("Start new idea\n")
	if m.starting {
		b.WriteString(warnStyle.Render("Starting run..."))
		return b.String()
	}
	b.WriteString("task> ")
	b.WriteString(m.startText)
	participants := m.defaultParticipants()
	b.WriteString("\n")
	if len(participants) == 0 {
		b.WriteString(warnStyle.Render("No installed headless participants available"))
	} else {
		b.WriteString(fmt.Sprintf("participants: %s", strings.Join(participants, ", ")))
	}
	return b.String()
}

func (m model) defaultParticipants() []string {
	var participants []string
	for _, agent := range m.agents {
		if !agent.Found || len(agent.HeadlessArgs) == 0 {
			continue
		}
		if m.effectiveLaunchMode(agent) != agents.LaunchHeadless {
			continue
		}
		participants = append(participants, agent.ID)
	}
	return participants
}

func (m model) updateStartMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		if m.starting {
			return m, nil
		}
		m.startMode = false
		m.startText = ""
		m.startErr = ""
		return m, nil
	case "enter":
		if m.starting {
			return m, nil
		}
		task := strings.TrimSpace(m.startText)
		if task == "" {
			m.startErr = "task is required"
			return m, nil
		}
		participants := m.defaultParticipants()
		if len(participants) == 0 {
			m.startErr = "no installed headless participants available"
			return m, nil
		}
		m.starting = true
		m.startErr = ""
		return m, startRunCmd(m.ctx, m.startRun, StartRequest{
			Task:            task,
			Participants:    participants,
			LaunchOverrides: copyStringMap(m.launchOverrides),
		})
	case "backspace", "ctrl+h":
		if len(m.startText) > 0 {
			runes := []rune(m.startText)
			m.startText = string(runes[:len(runes)-1])
		}
		return m, nil
	}
	if len(msg.Runes) > 0 && !m.starting {
		m.startText += string(msg.Runes)
	}
	return m, nil
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

func refreshRunsCmd(refresh func() ([]runstate.RunSummary, error)) tea.Cmd {
	return func() tea.Msg {
		runs, err := refresh()
		return refreshRunsMsg{runs: runs, err: err}
	}
}

func refreshTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return refreshTickMsg(t)
	})
}

func startRunCmd(ctx context.Context, start StartRunFunc, request StartRequest) tea.Cmd {
	return func() tea.Msg {
		if start == nil {
			return startRunMsg{err: fmt.Errorf("start callback is not configured")}
		}
		run, err := start(ctx, request)
		return startRunMsg{run: run, err: err}
	}
}

func upsertRunSummary(runs []runstate.RunSummary, run runstate.RunSummary) []runstate.RunSummary {
	if run.RunID == "" {
		return runs
	}
	next := append([]runstate.RunSummary(nil), runs...)
	for i := range next {
		if next[i].RunID == run.RunID {
			next[i] = run
			return next
		}
	}
	return append([]runstate.RunSummary{run}, next...)
}

func indexRun(runs []runstate.RunSummary, runID string) int {
	for i, run := range runs {
		if run.RunID == runID {
			return i
		}
	}
	return 0
}

func copyStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	copy := make(map[string]string, len(values))
	for key, value := range values {
		copy[key] = value
	}
	return copy
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func displayRunState(run runstate.RunSummary) string {
	if run.Error != "" {
		return "error: " + run.Error
	}
	if run.Terminal {
		return valueOr(run.Outcome, "unknown")
	}
	if run.Liveness != "" {
		return run.Liveness
	}
	return runstate.Attention(run)
}

func agentProgress(run runstate.RunSummary) string {
	total := len(run.State.Agents)
	if total == 0 {
		total = len(run.Participants)
	}
	if total == 0 {
		return "-"
	}
	done := 0
	for _, agent := range run.State.Agents {
		switch agent.State {
		case runstate.StateFinished, runstate.StateFailed, runstate.StateSkipped:
			done++
		}
	}
	return fmt.Sprintf("%d/%d", done, total)
}

func agentDuration(agent runstate.AgentState) time.Duration {
	if agent.Duration > 0 {
		return agent.Duration
	}
	if agent.State == runstate.StateRunning && !agent.StartedAt.IsZero() {
		if elapsed := time.Since(agent.StartedAt); elapsed > 0 {
			return elapsed
		}
	}
	return 0
}

func logPreview(agent runstate.AgentState) string {
	var parts []string
	if agent.StdoutPath != "" {
		if tail := tailLogFile(agent.StdoutPath, 8192, 8); strings.TrimSpace(tail) != "" {
			parts = append(parts, "stdout:\n"+tail)
		}
	}
	if agent.StderrPath != "" {
		if tail := tailLogFile(agent.StderrPath, 8192, 8); strings.TrimSpace(tail) != "" {
			parts = append(parts, "stderr:\n"+tail)
		}
	}
	return strings.Join(parts, "\n\n")
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
