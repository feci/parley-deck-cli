package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

// KillResult reports the outcome of a per-agent kill request.
type KillResult struct {
	AgentID   string
	Killed    bool
	Cleared   bool // a stale "running" badge was cleared (no live process to kill)
	Failed    bool // the operation could not be completed (shown as an error)
	SegmentID string
	Message   string
}

// SteerAttemptRequest asks the live run to execute a steer as a fresh
// single-agent attempt whose stdout becomes the reply.
type SteerAttemptRequest struct {
	AgentID string
	Text    string
	SteerID string // correlates with the steer.requested/delivered events
}

// SteerAttemptResult is returned synchronously; the attempt itself runs async.
// StdoutPath is deterministic from SteerID so the TUI can tail it immediately.
type SteerAttemptResult struct {
	SteerID    string
	SegmentID  string
	StdoutPath string
	StderrPath string
	ReplyPath  string
	Status     string // "running" | "queued" | "rejected"
	Accepted   bool
	Message    string
}

// register records an in-flight attempt's cancel so KillAgent can target it.
func (h *Handle) register(agentID, segmentID, kind, steerID string, cancel context.CancelFunc) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.active == nil {
		h.active = map[string]*attempt{}
	}
	h.active[agentID] = &attempt{agentID: agentID, segmentID: segmentID, kind: kind, steerID: steerID, cancel: cancel}
}

// finish deregisters an attempt and reports whether it was killed by the user.
func (h *Handle) finish(agentID string) bool {
	if h == nil {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	att := h.active[agentID]
	killed := att != nil && att.killed
	delete(h.active, agentID)
	return killed
}

// agentRunning reports whether an attempt for agentID is currently registered.
func (h *Handle) agentRunning(agentID string) bool {
	if h == nil {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.active[agentID] != nil
}

// KillAgent cancels ONLY the named agent's in-flight attempt and records an
// agent.killed event. It never touches the run-wide context, so the other
// agents keep running. If the attempt already finished (normal completion won
// the race), it returns Killed:false and emits no event.
func (h *Handle) KillAgent(agentID string) KillResult {
	if h == nil {
		return KillResult{AgentID: agentID, Message: "no live run"}
	}
	h.mu.Lock()
	att := h.active[agentID]
	if att == nil {
		// Not in this process's memory — fall back to the durable path (handles
		// agents from a previous parley run, and stale badges).
		h.mu.Unlock()
		return KillAgentDurable(h.opts.Store, agentID)
	}
	if att.killed {
		// A concurrent kill already won; don't cancel or emit a second event.
		h.mu.Unlock()
		return KillResult{AgentID: agentID, Killed: false, Message: "already killing " + agentID}
	}
	att.killed = true
	seg, kind, steerID, cancel := att.segmentID, att.kind, att.steerID, att.cancel
	h.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	_ = h.opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "agent.killed",
		Data: map[string]any{"agent": agentID, "segment_id": seg, "kind": kind, "steer_id": steerID},
	})
	return KillResult{AgentID: agentID, Killed: true, SegmentID: seg, Message: "killed " + agentID}
}

// RunSteerAttempt executes a steer to one agent as a fresh single-agent attempt
// (depth-1 queue per agent). It returns synchronously with the deterministic
// reply paths; the process runs in a goroutine and emits steer.reply_* events.
func (h *Handle) RunSteerAttempt(ctx context.Context, req SteerAttemptRequest) (SteerAttemptResult, error) {
	if h == nil {
		return SteerAttemptResult{}, fmt.Errorf("no live run handle")
	}
	agent, ok := h.discoverAgent(req.AgentID)
	if !ok {
		return SteerAttemptResult{Accepted: false, Status: "rejected", Message: "unknown agent " + req.AgentID}, nil
	}
	if !containsString(h.opts.Idea.Participants, req.AgentID) {
		return SteerAttemptResult{Accepted: false, Status: "rejected", Message: req.AgentID + " is not a participant in this idea"}, nil
	}

	h.mu.Lock()
	if h.steerBusy == nil {
		h.steerBusy = map[string]bool{}
	}
	if h.steerBusy[req.AgentID] {
		h.mu.Unlock()
		return SteerAttemptResult{Accepted: false, Status: "rejected", Message: req.AgentID + " is already replying"}, nil
	}
	h.steerBusy[req.AgentID] = true
	queued := h.active[req.AgentID] != nil
	steerID := strings.TrimSpace(req.SteerID)
	if steerID == "" {
		h.steerSeq++
		steerID = fmt.Sprintf("steer-auto-%04d", h.steerSeq)
	}
	h.mu.Unlock()

	steerDir := filepath.Join(h.RunDir, "agents", req.AgentID, "steers", steerID)
	if err := fsutil.MkdirAllResilient(steerDir, 0o755); err != nil {
		h.clearSteerBusy(req.AgentID)
		return SteerAttemptResult{Accepted: false, Status: "rejected", Message: "mkdir: " + err.Error()}, nil
	}

	// A steer is a side conversation, not a round re-run, so it does NOT emit a
	// run.segment_started boundary — doing so would reset the targeted agent's
	// projected state (badge) and, for a QUEUED steer, reorder against the still-
	// running round attempt's segment. We use a steer-scoped segment label derived
	// from the unique steer id purely to correlate the steer.reply_* events; it is
	// known synchronously (no event-count race) so the result carries it.
	seg := "steer/" + steerID

	res := SteerAttemptResult{
		SteerID:    steerID,
		SegmentID:  seg,
		StdoutPath: filepath.Join(steerDir, "stdout.log"),
		StderrPath: filepath.Join(steerDir, "stderr.log"),
		ReplyPath:  filepath.Join(steerDir, "reply.md"),
		Accepted:   true,
		Status:     "running",
	}
	if queued {
		res.Status = "queued"
	}

	go h.runSteerAgent(ctx, agent, req, steerID, seg, res.StdoutPath, res.StderrPath, res.ReplyPath)
	return res, nil
}

// runSteerAgent waits for any in-flight round attempt to clear (depth-1 queue),
// then runs the agent in its pre-allocated steer segment and records the reply.
func (h *Handle) runSteerAgent(ctx context.Context, agent agents.Discovery, req SteerAttemptRequest, steerID, seg, stdoutPath, stderrPath, replyPath string) {
	defer h.clearSteerBusy(req.AgentID)
	start := time.Now().UTC()

	// Depth-1 queue: let any active round attempt for this agent finish first.
	// Bail if the caller's context or the run is canceled.
	for h.agentRunning(req.AgentID) {
		if ctx != nil && ctx.Err() != nil {
			h.appendSteerFailed(steerID, req.AgentID, seg, "canceled before start", stdoutPath)
			return
		}
		if h.rootCtx != nil && h.rootCtx.Err() != nil {
			h.appendSteerFailed(steerID, req.AgentID, seg, "run ended before steer started", stdoutPath)
			return
		}
		time.Sleep(150 * time.Millisecond)
	}

	_ = h.opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "steer.reply_started",
		Data: map[string]any{"id": steerID, "agent": req.AgentID, "segment_id": seg, "stdout": stdoutPath},
	})

	prompt := BuildSteerPrompt(agent, h.opts, steerID, req.Text, replyPath)

	parent := ctx
	if parent == nil {
		parent = h.rootCtx
	}
	attemptCtx, cancel := context.WithTimeout(parent, timeoutForAgent(h.opts.Timeout, agent))
	defer cancel()
	h.register(req.AgentID, seg, "steer", steerID, cancel)
	_, err := execAgentProcess(attemptCtx, h.opts.Root, h.opts.RunID, req.AgentID, h.opts.RunID+":"+req.AgentID+":"+steerID, agent, prompt, stdoutPath, stderrPath, nil, nil, supervisionForAgent(agent, timeoutForAgent(h.opts.Timeout, agent)), supervisionHooks{})
	killed := h.finish(req.AgentID)
	duration := time.Since(start)

	reply := readSteerReply(replyPath, stdoutPath)
	if err != nil || killed || strings.TrimSpace(reply) == "" {
		msg := "no reply produced"
		switch {
		case killed:
			msg = "killed by user"
		case err != nil:
			msg = err.Error()
			if attemptCtx.Err() != nil {
				msg = attemptCtx.Err().Error()
			}
		}
		h.appendSteerFailed(steerID, req.AgentID, seg, msg, stdoutPath)
		return
	}
	_ = h.opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "steer.replied",
		Data: map[string]any{
			"id": steerID, "agent": req.AgentID, "segment_id": seg,
			"reply_path": replyPath, "stdout": stdoutPath, "duration_ms": duration.Milliseconds(),
		},
	})
}

func (h *Handle) clearSteerBusy(agentID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.steerBusy, agentID)
}

func (h *Handle) appendSteerFailed(steerID, agentID, seg, msg, stdoutPath string) {
	_ = h.opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "steer.reply_failed",
		Data: map[string]any{"id": steerID, "agent": agentID, "segment_id": seg, "error": msg, "stdout": stdoutPath},
	})
}

func containsString(items []string, want string) bool {
	for _, s := range items {
		if s == want {
			return true
		}
	}
	return false
}

func (h *Handle) discoverAgent(agentID string) (agents.Discovery, bool) {
	// Resolve a roster id (claude-1) via the [roster.*] mapping as well as a bare
	// family id, so steering a roster-id participant finds its launch spec — with
	// identity = roster id and vendor dispatch = Adapter() (composite-agent-naming).
	if resolved, err := agents.ResolveParticipant(agentID, h.opts.Agents, resolveMapping(h.opts)); err == nil {
		return resolved, true
	}
	return agents.Discovery{}, false
}

// readSteerReply prefers the reply.md file; falls back to captured stdout.
func readSteerReply(replyPath, stdoutPath string) string {
	if data, err := os.ReadFile(replyPath); err == nil && strings.TrimSpace(string(data)) != "" {
		return strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile(stdoutPath); err == nil {
		return strings.TrimSpace(string(data))
	}
	return ""
}

// BuildSteerPrompt builds the one-shot prompt for a steer attempt: the owner's
// message plus enough context (the idea prompt + the agent's recent transcript)
// for a useful follow-up reply, with strict no-protocol-write instructions.
func BuildSteerPrompt(agent agents.Discovery, opts Options, steerID, text, replyPath string) string {
	var ideaPrompt string
	if data, err := os.ReadFile(filepath.Join(opts.Idea.Path, "00-prompt.md")); err == nil {
		ideaPrompt = string(data)
	}
	roundStdout := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "agents", agent.ID, "stdout.log")
	transcriptTail := tailFileBytes(roundStdout, 4<<10)

	return fmt.Sprintf(`You are %s, a participant in a Parley Deck run. The human owner has sent you a
direct follow-up message through the live TUI. Answer ONLY this message.

Rules:
- This is a side conversation, NOT a protocol round. Do NOT create or edit any
  protocol artifact (no round-NN/, FINAL.md, IMPLEMENTATION.md, review/…), and do
  NOT touch other agents' files.
- Write your reply to this exact file: %s
- Also print your reply to stdout so it shows in the owner's TUI.
- Be concise and directly useful. If you need to reference your prior work, use the
  context below; do not re-derive everything.

Owner's message (steer %s):
%s

Idea prompt for context:
%s

Your recent transcript (tail):
%s
`, agent.ID, replyPath, steerID, text, ideaPrompt, transcriptTail)
}

// tailFileBytes returns up to maxBytes from the end of a file (whole file if
// smaller), empty when the file is missing.
func tailFileBytes(path string, maxBytes int64) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	start := int64(0)
	if info.Size() > maxBytes {
		start = info.Size() - maxBytes
	}
	if _, err := f.Seek(start, 0); err != nil {
		return ""
	}
	data, err := readAllLimited(f, maxBytes)
	if err != nil {
		return ""
	}
	return string(data)
}

func readAllLimited(f *os.File, maxBytes int64) ([]byte, error) {
	buf := make([]byte, maxBytes)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return nil, err
	}
	return buf[:n], nil
}
