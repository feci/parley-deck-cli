package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

// --- protocol visibility (tui-protocol-visibility) ----------------------------
//
// The model caches one ProtocolSnapshot; View only ever renders the cache.
// Snapshots are built asynchronously (buildProtoCmd), gated by runToken +
// protoSeq, with at most one build in flight (protoDirty coalesces re-triggers).
// Reconcile cadence: 15s while attached and running, 60s when done/stale/
// detached. Nothing protocol-related runs on the 250ms event tick.

const (
	ribbonCollapsed = iota
	ribbonExpanded
	ribbonHidden
)

const (
	narrateProtocol = iota // default: woven protocol lines
	narrateVerbose         // + per-agent ACP tool activity
	narrateOff
)

const narratorRingCap = 32

type narratorEntry struct {
	seq  int
	line string
}

type growthInfo struct {
	stdout, stderr int64
	lastGrowth     time.Time
}

type protoMsg struct {
	snap  ProtocolSnapshot
	token int
	seq   int
}

type protoTickMsg struct{ token int }

type growthMsg struct {
	sizes map[string][2]int64 // agentID → {stdout, stderr} bytes
	at    time.Time
	token int
}

type growthTickMsg struct{ token int }

func protoTickCmd(token int, after time.Duration) tea.Cmd {
	return tea.Tick(after, func(time.Time) tea.Msg { return protoTickMsg{token: token} })
}

func growthTickCmd(token int) tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return growthTickMsg{token: token} })
}

// buildProtoCmd computes the snapshot off the render path. The input carries
// value copies; results are dropped when token or seq are stale.
func buildProtoCmd(in ProtocolSnapshotInput, token, seq int) tea.Cmd {
	return func() tea.Msg {
		snap, _ := BuildProtocolSnapshot(in)
		return protoMsg{snap: snap, token: token, seq: seq}
	}
}

// statGrowthCmd stats the given agents' stdout/stderr logs (≤2 stats per agent)
// so unvisited tabs can show activity glyphs without tailing their logs.
func statGrowthCmd(paths map[string][2]string, token int) tea.Cmd {
	return func() tea.Msg {
		sizes := make(map[string][2]int64, len(paths))
		for id, p := range paths {
			var s [2]int64
			for i := 0; i < 2; i++ {
				if p[i] == "" {
					continue
				}
				if info, err := os.Stat(p[i]); err == nil {
					s[i] = info.Size()
				}
			}
			sizes[id] = s
		}
		return growthMsg{sizes: sizes, at: time.Now(), token: token}
	}
}

// protoInterval picks the reconcile cadence: tight while a live attached run is
// in flight, relaxed for done/stale/detached runs.
func (m liveModel) protoInterval() time.Duration {
	if m.done || !m.attached() || m.staleAgentCount() > 0 {
		return 60 * time.Second
	}
	return 15 * time.Second
}

// scheduleProtoRefresh fires one async snapshot build, or marks the state dirty
// when one is already in flight (the completion handler re-fires once).
func (m *liveModel) scheduleProtoRefresh() tea.Cmd {
	if !m.hasRun() {
		return nil
	}
	if m.protoBusy {
		m.protoDirty = true
		return nil
	}
	m.protoBusy = true
	m.protoSeq++
	return buildProtoCmd(m.protoInput(), m.runToken, m.protoSeq)
}

// protoInput snapshots the model state the builder may use (value copies only).
func (m liveModel) protoInput() ProtocolSnapshotInput {
	ideaDir := strings.TrimSpace(m.opts.Idea.Path)
	if ideaDir == "" {
		for _, idea := range m.opts.Status.Ideas {
			if idea.Slug == m.opts.Idea.Slug && strings.TrimSpace(idea.Path) != "" {
				ideaDir = idea.Path
				break
			}
		}
	}
	if ideaDir == "" && m.opts.Idea.Slug != "" {
		ideaDir = filepath.Join(m.opts.Root, protocol.DeckDir, "ideas", m.opts.Idea.Slug)
	}
	return ProtocolSnapshotInput{
		Root:         m.opts.Root,
		RunID:        m.opts.RunID,
		RunDir:       m.opts.RunDir,
		IdeaSlug:     m.opts.Idea.Slug,
		IdeaDir:      ideaDir,
		Participants: append([]string(nil), m.opts.Participants...),
		MaxRounds:    4,
		Events:       m.events,
		Questions:    append([]hitl.Question(nil), m.questions...),
		State:        m.state,
		Previous:     m.proto,
		Now:          m.now,
	}
}

// growthPaths lists the running agents whose buffers are NOT loaded — those are
// the tabs the stat cache covers (loaded buffers track growth via tail cursors).
func (m liveModel) growthPaths() map[string][2]string {
	paths := map[string][2]string{}
	for i := range m.state.Agents {
		a := &m.state.Agents[i]
		if a.State != stateRunning {
			continue
		}
		if b := m.buffers[a.ID]; b != nil && b.loaded {
			continue
		}
		paths[a.ID] = [2]string{a.StdoutPath, a.StderrPath}
	}
	return paths
}

// noteRuntimeFlags caches declared per-agent runtime flags (buffers_stdout)
// from the run.created payload so renders never scan the event log.
func (m *liveModel) noteRuntimeFlags(events []store.Event) {
	for _, e := range events {
		if e.Type != "run.created" {
			continue
		}
		runtime, _ := e.Data["runtime"].([]any)
		for _, entry := range runtime {
			data, _ := entry.(map[string]any)
			agent, _ := data["agent"].(string)
			if agent == "" {
				continue
			}
			// Tri-state: record explicit true AND false — a declared false must
			// suppress the heuristic (review cycle-1 fix 4).
			if buffers, ok := data["buffers_stdout"].(bool); ok {
				if m.buffersStdout == nil {
					m.buffersStdout = map[string]bool{}
				}
				m.buffersStdout[agent] = buffers
			}
		}
	}
}

// --- narrator ------------------------------------------------------------------

// narratorLine renders one woven protocol rule-line.
func narratorLine(e store.Event) string {
	stamp := ""
	if !e.Time.IsZero() {
		stamp = e.Time.Local().Format("15:04:05") + " "
	}
	return "── " + stamp + friendlyEventText(e) + " ──"
}

// friendlyEventText turns an event into a concise human line: the raw type is
// dropped when the summary is self-describing, otherwise mapped to a verb
// (review cycle-1 fix 10).
func friendlyEventText(e store.Event) string {
	s := runstate.SummarizeEvent(e)
	round, _ := e.Data["round"].(string)
	switch e.Type {
	case "agent.started":
		return s.Agent + " started"
	case "agent.finished":
		if s.Text != "" && s.Text != s.Agent {
			return s.Text // already "codex wrote round-02/codex.md"
		}
		return s.Agent + " finished"
	case "agent.failed":
		if detail := strings.TrimSpace(strings.TrimPrefix(s.Text, s.Agent)); detail != "" {
			return s.Agent + " failed: " + detail
		}
		return s.Agent + " failed"
	case "agent.skipped":
		if detail := strings.TrimSpace(strings.TrimPrefix(s.Text, s.Agent)); detail != "" {
			return s.Agent + " skipped: " + detail
		}
		return s.Agent + " skipped"
	case "agent.killed":
		return s.Agent + " killed"
	case "round.completed":
		return strings.TrimSpace(round+" complete") + " (" + s.Text + ")"
	case "round.incomplete":
		return strings.TrimSpace(round+" incomplete") + " (" + s.Text + ")"
	case "run.phase":
		phase, _ := e.Data["phase"].(string)
		action, _ := e.Data["action"].(string)
		label, _ := e.Data["round_label"].(string)
		return strings.TrimSpace("phase → " + phase + " " + label + " (" + action + ")")
	case "run.created":
		return strings.TrimSpace("run created " + s.Text)
	case "run.segment_started":
		return strings.TrimSpace("segment started " + s.Text)
	case "run.failed":
		return "run failed"
	case "run.manifest_deferred":
		return "run manifest deferred (transient write failure)"
	case "hitl.question":
		return "question from " + s.Agent + " — /answer to reply"
	case "hitl.answered":
		return "question answered (" + s.Agent + ")"
	default:
		// Fallback: humanize the type and keep the summary text.
		t := strings.NewReplacer(".", " ", "_", " ").Replace(e.Type)
		return strings.TrimSpace(t + " " + s.Text)
	}
}

// appendProtocolEvents weaves allowlisted protocol events into every loaded
// transcript and the bounded replay ring (consensus D7). Verbose mode adds the
// agent's own ACP activity to its tab only.
func (m *liveModel) appendProtocolEvents(events []store.Event) {
	if m.narrateMode == narrateOff {
		return
	}
	for _, e := range events {
		if m.narrateMode == narrateVerbose && strings.HasPrefix(e.Type, "agent.acp.") && e.Type != "agent.acp.message_chunk" {
			agentID, _ := e.Data["agent"].(string)
			if b := m.buffers[agentID]; b != nil && b.loaded {
				m.appendNarratorTo(b, narratorLine(e))
			}
			continue
		}
		if !narratorTypes[e.Type] {
			continue
		}
		line := narratorLine(e)
		m.narratorTotal++
		m.narratorRing = append(m.narratorRing, narratorEntry{seq: m.narratorTotal, line: line})
		if len(m.narratorRing) > narratorRingCap {
			m.narratorRing = m.narratorRing[len(m.narratorRing)-narratorRingCap:]
		}
		for _, b := range m.buffers {
			if !b.loaded {
				continue
			}
			m.appendNarratorTo(b, line)
			b.narratorSeq = m.narratorTotal
		}
	}
}

func (m *liveModel) appendNarratorTo(b *agentBuffer, line string) {
	b.lines = append(b.lines, transcriptLine{Text: line, Stream: transcriptEvent})
	var capped bool
	b.lines, capped = capTranscriptLines(b.lines)
	if capped {
		b.trunc = true
	}
	if b.follow {
		b.scroll = m.bufferBottom(b)
	}
}

// replayNarrator backfills the ring into a freshly loaded buffer exactly once
// (per-buffer seq prevents duplicates).
func (m *liveModel) replayNarrator(b *agentBuffer) {
	for _, entry := range m.narratorRing {
		if entry.seq > b.narratorSeq {
			m.appendNarratorTo(b, entry.line)
			b.narratorSeq = entry.seq
		}
	}
}

// --- glyphs ----------------------------------------------------------------------

var spinnerFrames = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// agentGlyph is the conservative tab-strip activity glyph (consensus D8).
func (m liveModel) agentGlyph(agentID string) string {
	a := m.agentByID(agentID)
	if a == nil {
		return "○"
	}
	if a.State == stateRunning && m.agentLiveness(agentID) == "stale" {
		return "!"
	}
	switch a.State {
	case stateRunning:
		if m.agentOutputFlowing(agentID) {
			return spinnerFrames[(m.now.Unix())%int64(len(spinnerFrames))]
		}
		return "·"
	case stateFinished:
		// The runner downgrades a bad artifact to agent.failed, so finished
		// always means artifact_ok (runner.go eventType switch).
		return "✓"
	case stateFailed:
		return "✗"
	case stateKilled:
		return "x"
	case stateSkipped:
		return "-"
	default:
		return "○"
	}
}

// agentOutputFlowing reports recent stdout/stderr growth: tail cursors for
// loaded buffers, the 2s stat cache for unvisited tabs.
func (m liveModel) agentOutputFlowing(agentID string) bool {
	if b := m.buffers[agentID]; b != nil && b.loaded && !b.lastGrowthAt.IsZero() {
		return m.now.Sub(b.lastGrowthAt) < 5*time.Second
	}
	if g, ok := m.growth[agentID]; ok && !g.lastGrowth.IsZero() {
		return m.now.Sub(g.lastGrowth) < 5*time.Second
	}
	return false
}

// agentBuffersStdout: declared flag first (run.created runtime, tri-state — an
// explicit false suppresses the heuristic), heuristic fallback only when
// undeclared (running, zero stdout, >30s elapsed).
func (m liveModel) agentBuffersStdout(agentID string) bool {
	if declared, ok := m.buffersStdout[agentID]; ok {
		return declared
	}
	a := m.agentByID(agentID)
	if a == nil || a.State != stateRunning || a.StartedAt.IsZero() {
		return false
	}
	if b := m.buffers[agentID]; b != nil && b.loaded && b.stdout.offset == 0 {
		return m.now.Sub(a.StartedAt) > 30*time.Second
	}
	return false
}

// --- renders ---------------------------------------------------------------------

var stepTitles = map[string]string{
	"kickoff": "Kickoff", "round-01": "Round-01", "cross-review": "Cross-Review",
	"consensus": "Consensus", "final": "Final", "implement": "Implement",
	"review": "Review", "review-consensus": "Review-Consensus", "fix-up": "Fix-Up",
	"complete": "Complete", "unknown": "Unknown",
}

var pipelineShort = [9]string{"Kick", "R01", "XRev", "Cons", "Final", "Impl", "Revw", "RCon", "Fixp"}

// ribbonHeight is the number of rows the ribbon occupies in the current mode.
func (m liveModel) ribbonHeight() int {
	if !m.hasRun() || m.proto == nil || m.ribbonMode == ribbonHidden {
		return 0
	}
	if m.ribbonMode == ribbonExpanded {
		return 3
	}
	return 1
}

// renderRibbon renders the protocol awareness strip under the tab strip.
func (m liveModel) renderRibbon(width int) string {
	p := m.proto
	if p == nil {
		return ""
	}
	stale := m.staleAgentCount() > 0
	if m.ribbonMode == ribbonExpanded {
		return m.renderRibbonExpanded(width, stale)
	}

	var b strings.Builder
	if stale {
		b.WriteString("[STALE] ")
	}
	if p.Step < 0 {
		b.WriteString("◆ protocol state unavailable · reconcile retrying")
	} else {
		fmt.Fprintf(&b, "◆ Ph %d: %s", p.Step, stepTitles[p.StepName])
		if p.Step == 1 || p.Step == 2 {
			fmt.Fprintf(&b, " (R%02d)", p.CurrentRound) // D9 string: no /total denominator
		}
		if p.Blocked {
			fmt.Fprintf(&b, " BLOCKED → reopening round-%02d", p.CurrentRound+1)
		}
		if n, total := deliveredCount(p.Delivery); total > 0 {
			mark := ""
			if p.DiskFallback {
				mark = "?"
			}
			fmt.Fprintf(&b, " · Delivered %d/%d%s", n, total, mark)
		}
		if p.Signoffs != nil {
			signed := len(p.Signoffs.Participants) - len(p.Signoffs.Missing)
			fmt.Fprintf(&b, " · Signoffs %d/%d (%s)", signed, len(p.Signoffs.Participants), p.Signoffs.Triage)
		}
		if len(p.Waiting) > 0 {
			fmt.Fprintf(&b, " · Waiting: %s", strings.Join(p.Waiting, ", "))
		}
		if p.Next != nil {
			fmt.Fprintf(&b, " · Next: %s", p.Next.Kind)
		}
		if age := m.now.Sub(p.ReconciledAt); age > 30*time.Second {
			fmt.Fprintf(&b, " · reconciled %s ago", formatShortAge(age))
		}
		if p.Err != "" {
			b.WriteString(" · reconcile retrying")
		}
	}
	line := truncateText(b.String(), width-4) + "  " + "⌃P"
	if stale || p.Blocked {
		return warnStyle.Render(truncateText(line, width))
	}
	return mutedStyle.Render(truncateText(line, width))
}

func (m liveModel) renderRibbonExpanded(width int, stale bool) string {
	p := m.proto
	var pipe strings.Builder
	pipe.WriteString("Pipeline: ")
	for i, name := range pipelineShort {
		if i > 0 {
			pipe.WriteString(" ── ")
		}
		switch {
		case p.Step >= 0 && i < p.Step:
			pipe.WriteString(name + " ✓")
		case i == p.Step:
			pipe.WriteString(name + " ▶")
		default:
			pipe.WriteString(name)
		}
	}

	var deliv strings.Builder
	deliv.WriteString("Delivery: ")
	if len(p.Delivery) == 0 && p.Signoffs != nil {
		parts := make([]string, 0, len(p.Signoffs.Signoffs)+len(p.Signoffs.Missing))
		for _, s := range p.Signoffs.Signoffs {
			parts = append(parts, fmt.Sprintf("%s (%s)", s.Agent, s.Status))
		}
		for _, missing := range p.Signoffs.Missing {
			parts = append(parts, missing+" (— missing)")
		}
		deliv.WriteString(strings.Join(parts, " · "))
	} else if len(p.Delivery) == 0 {
		deliv.WriteString("—")
	} else {
		parts := make([]string, 0, len(p.Delivery))
		for _, row := range p.Delivery {
			parts = append(parts, deliveryCell(m, row))
		}
		deliv.WriteString(strings.Join(parts, " · "))
	}

	var system strings.Builder
	system.WriteString("System:   ")
	if stale {
		system.WriteString("[STALE] ")
	}
	if p.Next != nil {
		system.WriteString("Next: " + valueOr(p.Next.Summary, p.Next.Kind))
	} else {
		system.WriteString("Next: —")
	}
	source := "events"
	if p.DiskFallback {
		source = "disk fallback"
	}
	fmt.Fprintf(&system, " · Reconciled %s ago (%s)", formatShortAge(m.now.Sub(p.ReconciledAt)), source)
	if p.Err != "" {
		system.WriteString(" · reconcile retrying")
	}

	style := mutedStyle
	if stale || p.Blocked {
		style = warnStyle
	}
	return style.Render(truncateText(pipe.String(), width)) + "\n" +
		style.Render(truncateText(deliv.String(), width)) + "\n" +
		style.Render(truncateText(system.String(), width))
}

func deliveryCell(m liveModel, row AgentDelivery) string {
	switch row.State {
	case "delivered":
		mark := "✓"
		if row.Unvalidated {
			mark = "✓?"
		}
		if !row.At.IsZero() {
			return fmt.Sprintf("%s (%s %s)", row.ID, mark, row.At.Local().Format("15:04"))
		}
		return fmt.Sprintf("%s (%s)", row.ID, mark)
	case "running":
		if a := m.agentByID(row.ID); a != nil {
			return fmt.Sprintf("%s (⏸ %s)", row.ID, formatAgentDuration(*a, m.now))
		}
		return row.ID + " (running)"
	case "failed":
		return row.ID + " (✗ failed)"
	case "killed":
		return row.ID + " (x killed)"
	case "skipped":
		return row.ID + " (– skipped)"
	default:
		return row.ID + " (○ pending)"
	}
}

func deliveredCount(rows []AgentDelivery) (delivered, total int) {
	for _, row := range rows {
		if row.State == "delivered" {
			delivered++
		}
	}
	return delivered, len(rows)
}

// statusPhaseSegment is the compressed status-line phase grammar (D10), falling
// back to the legacy round status until the first snapshot lands.
func (m liveModel) statusPhaseSegment() string {
	if m.proto == nil || m.proto.Step < 0 {
		return "round=" + displayRoundStatus(m.state.RoundStatus, m.done, m.opts.Resume)
	}
	p := m.proto
	seg := fmt.Sprintf("ph=%d:%s", p.Step, shortPhaseName(p.Step, p.StepName, p.RoundLabel))
	if p.Blocked {
		seg += " BLOCKED"
	}
	if len(p.Waiting) > 0 {
		seg += " wait=" + strings.Join(p.Waiting, ",")
	}
	return seg
}

// renderProtocolPanes renders the PIPELINE / DELIVERY / SIGNOFFS / NEXT panes at
// the top of the Protocol (Status) tab.
func (m liveModel) renderProtocolPanes() string {
	p := m.proto
	if p == nil {
		return mutedStyle.Render("protocol state not reconciled yet (/refresh to force)")
	}
	var b strings.Builder
	b.WriteString(sectionTitle(fmt.Sprintf("Pipeline — %s", m.opts.Idea.Slug)))
	if p.Step >= 0 {
		b.WriteString(mutedStyle.Render(fmt.Sprintf("  phase: %s", p.StepName)))
	}
	b.WriteString("\n")
	for i := 0; i <= 8; i++ {
		marker := "[ ]"
		switch {
		case p.Step > i:
			marker = "[✓]"
		case p.Step == i:
			marker = "[▶]"
		}
		title := stepTitles[stepNames[i]]
		if i == p.Step {
			// The current row uses the snapshot's actual step name so step 8
			// reads "Fix-Up" during a fix-up cycle, "Complete" when done.
			title = stepTitles[p.StepName]
		}
		detail := ""
		if i == p.Step {
			if n, total := deliveredCount(p.Delivery); total > 0 {
				detail = fmt.Sprintf("delivered %d/%d", n, total)
			}
			if len(p.Waiting) > 0 {
				detail = strings.TrimSpace(detail + " · waiting: " + strings.Join(p.Waiting, ", "))
			}
			if p.Blocked {
				detail = strings.TrimSpace(detail + " · BLOCKED")
			}
		}
		line := fmt.Sprintf("  %s %d %-17s %s", marker, i, title, detail)
		if i == p.Step {
			b.WriteString(headerStyle.Render(strings.TrimRight(line, " ")))
		} else {
			b.WriteString(mutedStyle.Render(strings.TrimRight(line, " ")))
		}
		b.WriteString("\n")
	}
	if len(p.Delivery) > 0 {
		b.WriteString("\n")
		b.WriteString(sectionTitle(fmt.Sprintf("Delivery (%s)", p.RoundLabel)))
		b.WriteString("\n")
		for _, row := range p.Delivery {
			b.WriteString("  " + deliveryCell(m, row))
			if row.Note != "" {
				b.WriteString("  " + mutedStyle.Render(truncateText(row.Note, 48)))
			}
			b.WriteString("\n")
		}
	}
	if p.Signoffs != nil {
		b.WriteString("\n")
		schema := "consensus.md"
		if p.Signoffs.Review {
			schema = "review/consensus.md"
		}
		b.WriteString(sectionTitle(fmt.Sprintf("Signoffs — %s (triage: %s)", schema, p.Signoffs.Triage)))
		b.WriteString("\n")
		for _, s := range p.Signoffs.Signoffs {
			b.WriteString(fmt.Sprintf("  %-10s %s  %s\n", s.Agent, s.Status, s.Date))
		}
		for _, missing := range p.Signoffs.Missing {
			b.WriteString(mutedStyle.Render(fmt.Sprintf("  %-10s — missing", missing)))
			b.WriteString("\n")
		}
	}
	if p.Next != nil {
		b.WriteString("\n")
		b.WriteString(sectionTitle("Next"))
		b.WriteString("\n")
		b.WriteString("  " + truncateText(valueOr(p.Next.Summary, p.Next.Kind), 90))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// renderSilentPlaceholder replaces the bare "no output yet" body for a running
// agent: declared buffering, liveness, byte counters, and the agent's own recent
// protocol activity (all from memory — zero I/O).
func (m liveModel) renderSilentPlaceholder(agentID string, width int) string {
	var b strings.Builder
	// Two-space indent on every line, matching the transcript placeholders
	// (review cycle-1 fix 7).
	writeLine := func(text string) {
		b.WriteString(mutedStyle.Render(truncateText("  "+text, width-1)))
		b.WriteString("\n")
	}
	a := m.agentByID(agentID)
	buffers := m.agentBuffersStdout(agentID)
	if buffers {
		writeLine("◆ " + agentID + " buffers all stdout until exit; stderr is live.")
	} else {
		writeLine("◆ no output yet")
	}
	if a != nil {
		status := fmt.Sprintf("status: %s %s", strings.ToUpper(a.State), formatAgentDuration(*a, m.now))
		if lv := m.agentLiveness(agentID); lv != "" {
			status += " · proc:" + lv
		}
		writeLine(status)
	}
	if buf := m.buffers[agentID]; buf != nil && buf.loaded {
		counters := fmt.Sprintf("stdout: %s", formatByteCount(buf.stdout.offset))
		if buffers {
			counters += " (buffered)"
		}
		counters += fmt.Sprintf(" · stderr: %s", formatByteCount(buf.stderr.offset))
		writeLine(counters)
	}
	recent := m.recentAgentActivity(agentID, 5)
	if len(recent) > 0 {
		b.WriteString("\n")
		writeLine("recent activity:")
		for _, line := range recent {
			writeLine("· " + line)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// recentAgentActivity returns the agent's last N narrator-allowlisted event
// summaries from the in-memory event slice.
func (m liveModel) recentAgentActivity(agentID string, n int) []string {
	var out []string
	for i := len(m.events) - 1; i >= 0 && len(out) < n; i-- {
		e := m.events[i]
		if agent, _ := e.Data["agent"].(string); agent != agentID {
			continue
		}
		if !narratorTypes[e.Type] && !strings.HasPrefix(e.Type, "agent.acp.") {
			continue
		}
		if e.Type == "agent.acp.message_chunk" {
			continue
		}
		stamp := ""
		if !e.Time.IsZero() {
			stamp = e.Time.Local().Format("15:04:05") + " "
		}
		out = append(out, stamp+friendlyEventText(e))
	}
	// reverse to chronological order
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// ideaPhaseChip computes the Home phase chip for one idea (called only from
// refreshHomeRuns — never on a tick).
func ideaPhaseChip(ideaPath string) string {
	if strings.TrimSpace(ideaPath) == "" {
		return ""
	}
	detail, err := driver.RebuildDetail(ideaPath, 4)
	if err != nil {
		return ""
	}
	step, name := displayStep(detail)
	return fmt.Sprintf("Ph %d: %s", step, stepTitles[name])
}

func formatByteCount(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KiB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func formatShortAge(d time.Duration) string {
	switch {
	case d >= time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d >= time.Minute:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	default:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
}
