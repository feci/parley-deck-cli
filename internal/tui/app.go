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
	"parley-deck-cli/internal/runaction"
	"parley-deck-cli/internal/runstate"
)

type WorkspaceOptions struct {
	Root         string
	Status       protocol.WorkspaceStatus
	Agents       []agents.Discovery
	Runs         []runstate.RunSummary
	Context      context.Context
	RefreshRuns  func() ([]runstate.RunSummary, error)
	StartRun     StartRunFunc
	ActionRunner ActionRunner
}

type StartRequest struct {
	Task            string
	Participants    []string
	LaunchOverrides map[string]string
	Auto            bool
}

type StartRunFunc func(context.Context, StartRequest) (runstate.RunSummary, error)

type ActionRequest struct {
	Run    runstate.RunSummary
	Action runaction.NextAction
}

type ActionResult struct {
	Message string
	Command string
	Refresh bool
}

type ActionRunner func(context.Context, ActionRequest) (ActionResult, error)

type model struct {
	root            string
	status          protocol.WorkspaceStatus
	agents          []agents.Discovery
	runs            []runstate.RunSummary
	ctx             context.Context
	refreshRuns     func() ([]runstate.RunSummary, error)
	startRun        StartRunFunc
	actionRunner    ActionRunner
	width           int
	height          int
	focus           focusZone
	selectedIdea    int
	selectedAgent   int
	selectedAction  int
	confirmActionID string
	launchOverrides map[string]string
	startMode       bool
	starting        bool
	actionRunning   bool
	startText       string
	startErr        string
	actionText      string
	errText         string
}

type focusZone string

const (
	focusIdeas   focusZone = "ideas"
	focusActions focusZone = "actions"
	focusAgents  focusZone = "agents"
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

type refreshRunsMsg struct {
	runs []runstate.RunSummary
	err  error
}

type startRunMsg struct {
	run runstate.RunSummary
	err error
}

type actionRunMsg struct {
	result ActionResult
	err    error
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
	// Only surface agents that were actually found on this machine; missing
	// agents are not selectable or launchable, so they are not displayed.
	available := make([]agents.Discovery, 0, len(opts.Agents))
	for _, agent := range opts.Agents {
		if agent.Found {
			available = append(available, agent)
		}
	}
	m := model{
		root:            root,
		status:          opts.Status,
		agents:          available,
		runs:            append([]runstate.RunSummary(nil), opts.Runs...),
		ctx:             ctx,
		refreshRuns:     opts.RefreshRuns,
		startRun:        opts.StartRun,
		actionRunner:    opts.ActionRunner,
		width:           100,
		height:          defaultDashboardHeight,
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
		cmds = append(cmds, refreshRunsCmd(m.refreshRuns))
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.startMode {
			return m.updateStartMode(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.switchFocus(1)
			return m, nil
		case "shift+tab":
			m.switchFocus(-1)
			return m, nil
		case "enter":
			if m.focus == focusActions {
				return m.triggerSelectedAction()
			}
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
		case "a":
			m.setSelectedAgentMode(agents.LaunchACP)
			return m, nil
		case "m":
			m.setSelectedAgentMode(agents.LaunchManual)
			return m, nil
		case "x":
			m.clearSelectedAgentMode()
			return m, nil
		}
	case refreshRunsMsg:
		selectedRunID := ""
		selectedAgentID := ""
		selectedActionID := ""
		if run, ok := m.selectedRun(); ok {
			selectedRunID = run.RunID
			if len(run.State.Agents) > 0 {
				selectedAgentID = run.State.Agents[clampIndex(m.selectedAgent, len(run.State.Agents))].ID
			}
			if len(run.NextActions) > 0 {
				selectedActionID = run.NextActions[clampIndex(m.selectedAction, len(run.NextActions))].ID
			}
		}
		if msg.err != nil {
			m.errText = msg.err.Error()
		} else {
			m.errText = ""
			m.runs = msg.runs
			if selectedRunID != "" {
				m.selectedIdea = indexRun(m.runs, selectedRunID)
			}
			if selectedAgentID != "" {
				m.selectedAgent = indexAgent(m.selectedRunStateAgents(), selectedAgentID)
			}
			if selectedActionID != "" {
				m.selectedAction = indexAction(m.selectedRunActions(), selectedActionID)
			}
			m.clampSelections()
		}
		if m.refreshRuns != nil {
			return m, refreshTickCmd()
		}
		return m, nil
	case refreshTickMsg:
		if m.refreshRuns == nil {
			return m, nil
		}
		return m, refreshRunsCmd(m.refreshRuns)
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
	case actionRunMsg:
		m.actionRunning = false
		m.confirmActionID = ""
		if msg.err != nil {
			m.actionText = msg.err.Error()
			if msg.result.Command != "" {
				m.actionText += "\ncommand: " + msg.result.Command
			}
			return m, nil
		}
		m.actionText = valueOr(msg.result.Message, "action completed")
		if msg.result.Command != "" && !strings.Contains(m.actionText, msg.result.Command) {
			m.actionText += "\ncommand: " + msg.result.Command
		}
		if msg.result.Refresh && m.refreshRuns != nil {
			return m, refreshRunsCmd(m.refreshRuns)
		}
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	width := tuiWidth(m.width)
	height := tuiHeight(m.height, defaultDashboardHeight)
	if height < compactDashboardHeight {
		return m.renderCompactDashboard(width, height)
	}
	return m.renderDashboard(width, height)
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
		status := msg.status
		m.status = &status
		m.errText = ""
		dashboard := newModel(WorkspaceOptions{Root: m.root, Status: status, Agents: m.agents})
		dashboard.width = m.width
		dashboard.height = m.height
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

func (m model) renderDashboard(width, height int) string {
	bodyWidth := width - 4
	leftWidth := clampInt(width/3, 30, 44)
	rightWidth := bodyWidth - leftWidth - 4
	if rightWidth < 42 {
		rightWidth = 42
		leftWidth = bodyWidth - rightWidth - 4
	}
	if leftWidth < 28 {
		leftWidth = 28
	}

	left := panelStyle(m.focus == focusIdeas).Width(leftWidth).Render(clipLines(m.renderIdeas(), height-7))
	right := panelStyle(m.focus == focusAgents || m.focus == focusActions).Width(rightWidth).Render(clipLines(m.renderWorkspace(), height-7))
	footer := mutedStyle.Render(m.renderFooter())

	return strings.Join([]string{
		m.dashboardHeader("normal"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right),
		"",
		footer,
		"",
	}, "\n")
}

func (m model) renderCompactDashboard(width, height int) string {
	var b strings.Builder
	b.WriteString(m.dashboardHeader("compact"))
	b.WriteString("\n")
	if m.errText != "" {
		b.WriteString(warnStyle.Render(m.errText))
		b.WriteString("\n")
	}
	if m.startErr != "" {
		b.WriteString(warnStyle.Render(m.startErr))
		b.WriteString("\n")
	}
	if m.startMode {
		b.WriteString(sectionTitle("Start new idea"))
		b.WriteString("\n")
		b.WriteString(oneLine(m.renderStartBox()))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(m.renderCompactFooter()))
		return clipLines(b.String(), height) + "\n"
	}

	b.WriteString(sectionTitle("Sessions"))
	b.WriteString("\n")
	b.WriteString(m.renderCompactRuns(3))
	b.WriteString("\n")

	if run, ok := m.selectedRun(); ok {
		b.WriteString(sectionTitle("Run details"))
		b.WriteString("\n")
		b.WriteString(m.renderCompactRunSummary(run))
		b.WriteString("\n")
		b.WriteString(sectionTitle("Agents"))
		b.WriteString("\n")
		b.WriteString(m.renderCompactRunAgents(run, 3))
		b.WriteString("\n")
		b.WriteString(sectionTitle("Actions"))
		b.WriteString("\n")
		b.WriteString(m.renderCompactActions(run, 2))
		b.WriteString("\n")
		b.WriteString(sectionTitle("Questions"))
		b.WriteString("\n")
		b.WriteString(m.renderCompactQuestions(run, 2))
	} else {
		b.WriteString(sectionTitle("Agents"))
		b.WriteString("\n")
		b.WriteString(m.renderCompactDiscoveredAgents(4))
	}
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render(m.renderCompactFooter()))
	return clipLines(b.String(), height) + "\n"
}

func (m model) dashboardHeader(layout string) string {
	focus := string(m.focus)
	if focus == "" {
		focus = string(focusIdeas)
	}
	parts := []string{
		fmt.Sprintf("Parley Deck  transport=%s", m.status.Transport),
		"mode=HITL default",
		"focus=" + focus,
	}
	if run, ok := m.selectedRun(); ok {
		parts = append(parts, fmt.Sprintf("run=%s %s %s", run.IdeaSlug, attentionBadge(valueOr(run.Attention, runstate.Attention(run))), agentProgress(run)))
	}
	if layout == "compact" {
		parts = append(parts, "layout=compact")
	}
	return headerStyle.Render(strings.Join(parts, "  "))
}

func (m model) renderWorkspace() string {
	if m.startMode {
		return m.renderStartBox()
	}
	run, ok := m.selectedRun()
	if !ok {
		return m.renderAgents()
	}

	var b strings.Builder
	b.WriteString(m.renderRunOverview(run))
	b.WriteString("\n\n")
	b.WriteString(m.renderRunAgents(run, 5))
	b.WriteString("\n\n")
	b.WriteString(m.renderEventsLimit(run, 5))
	b.WriteString("\n\n")
	b.WriteString(m.renderActionsAndQuestions(run, 4, 4))
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderRunOverview(run runstate.RunSummary) string {
	var b strings.Builder
	b.WriteString("Run details\n")
	b.WriteString(fmt.Sprintf("idea: %s  run: %s\n", valueOr(run.IdeaSlug, "unknown"), run.RunID))
	b.WriteString(fmt.Sprintf("attention: %s  state: %s  agents: %s  q: %d\n",
		attentionBadge(valueOr(run.Attention, runstate.Attention(run))),
		displayRunState(run),
		agentProgress(run),
		run.OpenQuestions,
	))
	if run.Task != "" {
		b.WriteString(fmt.Sprintf("task: %s\n", truncateText(run.Task, 96)))
	}
	if !run.LastEventAt.IsZero() {
		b.WriteString(fmt.Sprintf("last event: %s ago\n", formatDuration(run.LastEventAge)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderRunAgents(run runstate.RunSummary, limit int) string {
	var b strings.Builder
	b.WriteString("Agents\n")
	if len(run.State.Agents) == 0 {
		b.WriteString(mutedStyle.Render("no agent state yet"))
		return b.String()
	}
	for i, state := range run.State.Agents {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more agent(s)\n", len(run.State.Agents)-i))
			break
		}
		marker := " "
		if i == m.selectedAgent {
			marker = m.selectionMarker(focusAgents, true)
		}
		b.WriteString(fmt.Sprintf("%s %-8s %-14s %-8s %s\n",
			marker,
			state.ID,
			stateBadge(state.State),
			formatDuration(agentDuration(state)),
			truncateText(state.LatestEvent, 64),
		))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderEventsLimit(run runstate.RunSummary, limit int) string {
	var b strings.Builder
	b.WriteString("Event stream\n")
	if len(run.State.Recent) == 0 {
		b.WriteString(mutedStyle.Render("waiting for events"))
		return b.String()
	}
	for i, event := range run.State.Recent {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more event(s)\n", len(run.State.Recent)-i))
			break
		}
		b.WriteString(fmt.Sprintf("%s  %-16s %s\n", event.Time.Local().Format("15:04:05"), event.Type, truncateText(event.Text, 96)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderActionsAndQuestions(run runstate.RunSummary, actionLimit, questionLimit int) string {
	var b strings.Builder
	b.WriteString("Actions\n")
	if len(run.NextActions) == 0 {
		b.WriteString(mutedStyle.Render("no planner actions"))
	} else {
		for i, action := range run.NextActions {
			if i >= actionLimit {
				b.WriteString(fmt.Sprintf("  ... %d more action(s)\n", len(run.NextActions)-i))
				break
			}
			marker := m.selectionMarker(focusActions, i == m.selectedAction)
			summary := truncateText(valueOr(action.Summary, action.ID), 88)
			if m.actionRunning && i == m.selectedAction {
				summary = "[running] " + summary
			}
			b.WriteString(fmt.Sprintf("%s %-18s %-10s %s\n", marker, action.Kind, riskBadge(action.Risk), summary))
		}
	}
	if m.actionText != "" {
		b.WriteString("\n")
		if m.actionRunning {
			b.WriteString(warnStyle.Render(m.actionText))
		} else {
			b.WriteString(mutedStyle.Render(m.actionText))
		}
	}
	b.WriteString("\n\nQuestions\n")
	if len(run.Questions) == 0 {
		b.WriteString(mutedStyle.Render("no questions"))
		return strings.TrimRight(b.String(), "\n")
	}
	for i, question := range run.Questions {
		if i >= questionLimit {
			b.WriteString(fmt.Sprintf("  ... %d more question(s)\n", len(run.Questions)-i))
			break
		}
		marker := " "
		if question.Status == hitl.StatusOpen {
			marker = warnStyle.Render("!")
		}
		b.WriteString(fmt.Sprintf("%s %-28s %-13s %-10s %s\n", marker, question.ID, question.Status, riskBadge(question.Risk), truncateText(question.Prompt, 90)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderCompactRuns(limit int) string {
	if len(m.runs) == 0 {
		if len(m.status.Ideas) > 0 {
			return mutedStyle.Render(fmt.Sprintf("no runs found  %d protocol idea(s) on disk", len(m.status.Ideas)))
		}
		return mutedStyle.Render("no runs found")
	}
	var b strings.Builder
	for i, run := range m.runs {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more session(s)\n", len(m.runs)-i))
			break
		}
		marker := m.selectionMarker(focusIdeas, i == m.selectedIdea)
		b.WriteString(fmt.Sprintf("%s %-18s %-28s agents:%s q:%d",
			marker,
			attentionBadge(valueOr(run.Attention, runstate.Attention(run))),
			truncateText(run.IdeaSlug, 28),
			agentProgress(run),
			run.OpenQuestions,
		))
		if !run.LastEventAt.IsZero() {
			b.WriteString(fmt.Sprintf(" last:%s", formatDuration(run.LastEventAge)))
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderCompactRunSummary(run runstate.RunSummary) string {
	parts := []string{
		fmt.Sprintf("%s %s", attentionBadge(valueOr(run.Attention, runstate.Attention(run))), valueOr(run.IdeaSlug, "unknown")),
		"state=" + displayRunState(run),
		"agents=" + agentProgress(run),
		fmt.Sprintf("q=%d", run.OpenQuestions),
	}
	if run.Task != "" {
		parts = append(parts, "task="+truncateText(run.Task, 56))
	}
	return strings.Join(parts, "  ")
}

func (m model) renderCompactRunAgents(run runstate.RunSummary, limit int) string {
	if len(run.State.Agents) == 0 {
		return mutedStyle.Render("no agent state yet")
	}
	var b strings.Builder
	for i, state := range run.State.Agents {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more agent(s)\n", len(run.State.Agents)-i))
			break
		}
		marker := " "
		if i == m.selectedAgent {
			marker = m.selectionMarker(focusAgents, true)
		}
		b.WriteString(fmt.Sprintf("%s %-8s %-14s %s\n", marker, state.ID, stateBadge(state.State), truncateText(state.LatestEvent, 48)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderCompactActions(run runstate.RunSummary, limit int) string {
	if len(run.NextActions) == 0 {
		return mutedStyle.Render("no planner actions")
	}
	var b strings.Builder
	for i, action := range run.NextActions {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more action(s)\n", len(run.NextActions)-i))
			break
		}
		marker := m.selectionMarker(focusActions, i == m.selectedAction)
		b.WriteString(fmt.Sprintf("%s %-18s %-10s %s\n", marker, action.Kind, riskBadge(action.Risk), truncateText(valueOr(action.Summary, action.ID), 60)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderCompactQuestions(run runstate.RunSummary, limit int) string {
	if len(run.Questions) == 0 {
		return mutedStyle.Render("no questions")
	}
	var b strings.Builder
	for i, question := range run.Questions {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more question(s)\n", len(run.Questions)-i))
			break
		}
		marker := " "
		if question.Status == hitl.StatusOpen {
			marker = warnStyle.Render("!")
		}
		b.WriteString(fmt.Sprintf("%s %-22s %-13s %s\n", marker, question.ID, question.Status, truncateText(question.Prompt, 64)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderCompactDiscoveredAgents(limit int) string {
	if len(m.agents) == 0 {
		return mutedStyle.Render("no agents configured")
	}
	var b strings.Builder
	for i, agent := range m.agents {
		if i >= limit {
			b.WriteString(fmt.Sprintf("  ... %d more agent(s)\n", len(m.agents)-i))
			break
		}
		marker := m.selectionMarker(focusAgents, i == m.selectedAgent)
		state := stateBadge("missing")
		if agent.Found {
			state = stateBadge("found")
		}
		b.WriteString(fmt.Sprintf("%s %-8s %-10s mode:%s version:%s\n", marker, agent.ID, state, m.effectiveLaunchMode(agent), valueOr(agent.Version, "-")))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m model) renderIdeas() string {
	var b strings.Builder
	b.WriteString(sectionTitle("Sessions"))
	b.WriteString("\n")
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
		b.WriteString(fmt.Sprintf("%s %-18s %s\n", marker, attentionBadge(badge), run.IdeaSlug))
		b.WriteString(fmt.Sprintf("   run: %s\n", run.RunID))
		b.WriteString(fmt.Sprintf("   agents: %s  q: %d\n", agentProgress(run), run.OpenQuestions))
		if !run.LastEventAt.IsZero() {
			b.WriteString(fmt.Sprintf("   last: %s ago\n", formatDuration(run.LastEventAge)))
		}
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
	b.WriteString(fmt.Sprintf("acp: %s\n", acpCommandLine(agent)))
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
	return m.renderFooterText(true)
}

func (m model) renderCompactFooter() string {
	return m.renderFooterText(false)
}

func (m model) renderFooterText(includeLifecycleWarning bool) string {
	var parts []string
	if m.startMode {
		parts = append(parts, "Start: type task  enter launch  esc cancel")
	} else {
		keys := "Keys: N new  r refresh  tab focus  j/k select"
		if m.focus == focusActions {
			keys += "  enter action"
		}
		if m.focus == focusAgents {
			keys += "  h/i/a/m set agent mode  x clear mode"
		}
		keys += "  q/esc quit"
		parts = append(parts, keys)
		if includeLifecycleWarning {
			parts = append(parts, "Quit closes the TUI and cancels TUI-started runs; it does not detach child processes.")
		}
	}
	if m.errText != "" {
		parts = append([]string{warnStyle.Render(m.errText)}, parts...)
	}
	if m.startErr != "" {
		parts = append([]string{warnStyle.Render(m.startErr)}, parts...)
	}
	return strings.Join(parts, "  |  ")
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

func (m *model) switchFocus(delta int) {
	previous := m.focus
	zones := []focusZone{focusIdeas}
	if m.hasActionFocusTargets() {
		zones = append(zones, focusActions)
	}
	if m.hasAgentFocusTargets() {
		zones = append(zones, focusAgents)
	}
	if len(zones) == 1 {
		m.focus = zones[0]
		return
	}
	current := 0
	for i, zone := range zones {
		if zone == m.focus {
			current = i
			break
		}
	}
	next := (current + delta) % len(zones)
	if next < 0 {
		next += len(zones)
	}
	m.focus = zones[next]
	if previous != m.focus && m.focus != focusActions {
		m.confirmActionID = ""
	}
}

func (m *model) moveSelection(delta int) {
	switch m.focus {
	case focusIdeas:
		m.selectedIdea = clampIndex(m.selectedIdea+delta, len(m.runs))
		m.selectedAgent = 0
		m.selectedAction = 0
		m.actionText = ""
		m.confirmActionID = ""
	case focusAgents:
		if run, ok := m.selectedRun(); ok && len(run.State.Agents) > 0 {
			m.selectedAgent = clampIndex(m.selectedAgent+delta, len(run.State.Agents))
			return
		}
		m.selectedAgent = clampIndex(m.selectedAgent+delta, len(m.agents))
	case focusActions:
		previous := m.selectedAction
		m.selectedAction = clampIndex(m.selectedAction+delta, len(m.selectedRunActions()))
		if previous != m.selectedAction {
			m.actionText = ""
			m.confirmActionID = ""
		}
	}
}

func (m *model) clampSelections() {
	m.selectedIdea = clampIndex(m.selectedIdea, len(m.runs))
	if run, ok := m.selectedRun(); ok && len(run.State.Agents) > 0 {
		m.selectedAgent = clampIndex(m.selectedAgent, len(run.State.Agents))
	} else {
		m.selectedAgent = clampIndex(m.selectedAgent, len(m.agents))
	}
	m.selectedAction = clampIndex(m.selectedAction, len(m.selectedRunActions()))
	if m.focus == focusActions && !m.hasActionFocusTargets() {
		m.focus = focusIdeas
	}
}

func (m model) selectedRun() (runstate.RunSummary, bool) {
	if len(m.runs) == 0 {
		return runstate.RunSummary{}, false
	}
	return m.runs[clampIndex(m.selectedIdea, len(m.runs))], true
}

func (m model) selectedRunActions() []runaction.NextAction {
	run, ok := m.selectedRun()
	if !ok {
		return nil
	}
	return run.NextActions
}

func (m model) hasActionFocusTargets() bool {
	return len(m.selectedRunActions()) > 0
}

func (m model) hasAgentFocusTargets() bool {
	if len(m.agents) > 0 {
		return true
	}
	run, ok := m.selectedRun()
	return ok && len(run.State.Agents) > 0
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
		b.WriteString(warnStyle.Render("No installed launchable participants available"))
	} else {
		b.WriteString(fmt.Sprintf("participants: %s", strings.Join(participants, ", ")))
	}
	return b.String()
}

func (m model) defaultParticipants() []string {
	var participants []string
	for _, agent := range m.agents {
		if !agent.Found {
			continue
		}
		switch m.effectiveLaunchMode(agent) {
		case agents.LaunchHeadless:
			if len(agent.HeadlessArgs) == 0 {
				continue
			}
		case agents.LaunchACP:
			if agent.ACPArgs == nil {
				continue
			}
		default:
			continue
		}
		if agent.ID == "gemini" {
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
			m.startErr = "no installed launchable participants available"
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
	if !launchModeAvailable(agent, mode) {
		m.errText = fmt.Sprintf("%s has no %s launch configuration", agent.ID, mode)
		return
	}
	m.launchOverrides[agent.ID] = mode
	m.errText = ""
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
	m.errText = ""
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

func acpCommandLine(agent agents.Discovery) string {
	if agent.ACPArgs == nil {
		return "not configured"
	}
	command := agents.CLIDefault
	if len(agent.Commands) > 0 {
		command = agent.Commands[0]
	}
	if len(agent.ACPArgs) == 0 {
		return command
	}
	return command + " " + strings.Join(agent.ACPArgs, " ")
}

func launchModeAvailable(agent agents.Discovery, mode string) bool {
	switch agents.LaunchModeOrDefault(mode) {
	case agents.LaunchHeadless:
		return len(agent.HeadlessArgs) > 0
	case agents.LaunchACP:
		return agent.ACPArgs != nil
	case agents.LaunchInteractive, agents.LaunchManual:
		return true
	default:
		return false
	}
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

func (m model) triggerSelectedAction() (tea.Model, tea.Cmd) {
	if m.actionRunning {
		m.actionText = "action already running"
		return m, nil
	}
	run, ok := m.selectedRun()
	if !ok || len(run.NextActions) == 0 {
		m.actionText = "no planner action selected"
		return m, nil
	}
	action := run.NextActions[clampIndex(m.selectedAction, len(run.NextActions))]
	if m.actionRunner == nil {
		m.actionText = "action execution is not available in this view"
		if command := runaction.Command(action, run.RunID, run.IdeaSlug); command != "" {
			m.actionText += "\ncommand: " + command
		}
		return m, nil
	}
	if action.RequiresYes && m.confirmActionID != action.ID {
		m.confirmActionID = action.ID
		m.actionText = "Action requires confirmation; press enter again to run: " + valueOr(action.Summary, action.Kind)
		if command := runaction.Command(action, run.RunID, run.IdeaSlug); command != "" {
			m.actionText += "\ncommand: " + command
		}
		return m, nil
	}
	m.confirmActionID = ""
	m.actionRunning = true
	m.actionText = "Running action: " + valueOr(action.Summary, action.Kind)
	return m, actionRunCmd(m.ctx, m.actionRunner, ActionRequest{Run: run, Action: action})
}

func actionRunCmd(ctx context.Context, runner ActionRunner, request ActionRequest) tea.Cmd {
	return func() tea.Msg {
		if runner == nil {
			return actionRunMsg{err: fmt.Errorf("action runner is not configured")}
		}
		result, err := runner(ctx, request)
		return actionRunMsg{result: result, err: err}
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

func (m model) selectedRunStateAgents() []runstate.AgentState {
	run, ok := m.selectedRun()
	if !ok {
		return nil
	}
	return run.State.Agents
}

func indexAgent(agents []runstate.AgentState, agentID string) int {
	for i, agent := range agents {
		if agent.ID == agentID {
			return i
		}
	}
	return 0
}

func indexAction(actions []runaction.NextAction, actionID string) int {
	for i, action := range actions {
		if action.ID == actionID {
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
