package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

type Options struct {
	Root       string
	RunID      string
	Idea       protocol.IdeaStatus
	Task       string
	Agents     []agents.Discovery
	Timeout    time.Duration
	Store      store.Store
	Round      int
	RoundLabel string
	// Phase selects the prompt + validation contract: "" or "deliberation"
	// (Phase 1-4 rounds), "review" (Phase 6 code-review rounds), or
	// "implementation" (Phase 5). Set via the RunReviewRound/RunImplementation
	// helpers; defaults to deliberation.
	Phase string
	// ArtifactName, when set, makes the agent write idea.Path/<ArtifactName>
	// (e.g. IMPLEMENTATION.md) instead of the per-agent round path.
	ArtifactName string
	// Overwrite allows re-running an agent even if its artifact already exists
	// (used for per-cycle review-consensus re-drafts). Off by default.
	Overwrite bool
	// SegmentID tags agent events with the current run.segment_started segment so
	// the projection can scope state to the active segment (fixes stale terminal
	// badges after continue/resume/retry). Set by the round-run entry points via
	// appendSegmentStarted; empty for legacy/unsegmented callers.
	SegmentID string
	// tracker, when set (by RunRoundOneAsync), registers each agent attempt's
	// cancel func so the live Handle can KillAgent one agent. nil on the sync path.
	tracker *Handle
}

type Result struct {
	AgentID     string
	OutputPath  string
	StdoutPath  string
	StderrPath  string
	StartedAt   time.Time
	CompletedAt time.Time
	ExitError   string
	ArtifactOK  bool
	Skipped     bool
	SkipReason  string
	Warning     string
	Duration    time.Duration
	// Killed is set when the attempt was terminated by Handle.KillAgent (vs a
	// timeout or self-exit), so projection can show a distinct "killed" badge.
	Killed bool
}

// attempt tracks one in-flight agent process so the live Handle can cancel just
// that agent (kill) without touching the run-wide context.
type attempt struct {
	agentID   string
	segmentID string
	kind      string // "round" | "steer"
	steerID   string
	cancel    context.CancelFunc
	killed    bool
}

type Handle struct {
	RunID  string
	RunDir string

	done    chan struct{}
	results []Result
	mu      sync.Mutex

	// Live attempt control (tui-live-steering). opts/rootCtx are captured so the
	// handle can spawn steer attempts; active maps agentID->in-flight attempt;
	// steerBusy enforces the depth-1 per-agent steer queue; segmentMu serializes
	// segment-id allocation so concurrent steer attempts never collide.
	opts      Options
	rootCtx   context.Context
	active    map[string]*attempt
	steerBusy map[string]bool
	steerSeq  int // fallback unique steer id when a caller supplies none
}

func RunRoundOneAsync(ctx context.Context, opts Options) *Handle {
	handle := &Handle{
		RunID:     opts.RunID,
		RunDir:    filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID),
		done:      make(chan struct{}),
		rootCtx:   ctx,
		active:    map[string]*attempt{},
		steerBusy: map[string]bool{},
	}
	opts.tracker = handle
	handle.opts = opts
	go func() {
		defer close(handle.done)
		defer func() {
			if recovered := recover(); recovered != nil {
				now := time.Now().UTC()
				result := Result{
					AgentID:     "runner",
					CompletedAt: now,
					ExitError:   fmt.Sprintf("runner panic: %v", recovered),
				}
				_ = opts.Store.Append(store.Event{
					Time: now,
					Type: "run.failed",
					Data: map[string]any{"error": result.ExitError},
				})
				handle.setResults([]Result{result})
			}
		}()
		handle.setResults(RunRoundOne(ctx, opts))
	}()
	return handle
}

func (h *Handle) Done() <-chan struct{} {
	return h.done
}

func (h *Handle) Wait() []Result {
	<-h.done
	return h.Results()
}

func (h *Handle) Results() []Result {
	h.mu.Lock()
	defer h.mu.Unlock()
	results := make([]Result, len(h.results))
	copy(results, h.results)
	return results
}

func (h *Handle) setResults(results []Result) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.results = append([]Result(nil), results...)
}

func RunRoundOne(ctx context.Context, opts Options) []Result {
	if opts.Round == 0 {
		opts.Round = 1
	}
	if opts.RoundLabel == "" {
		opts.RoundLabel = "round-01"
	}

	selected := selectedAgents(opts.Idea.Participants, opts.Agents)
	opts.SegmentID = appendSegmentStarted(opts, segmentReason(opts), agentIDs(selected))
	results := make([]Result, len(selected))
	var wg sync.WaitGroup
	for i, agent := range selected {
		i, agent := i, agent
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = runAgent(ctx, opts, agent)
		}()
	}
	wg.Wait()

	eventType := "round.completed"
	okCount := 0
	for _, result := range results {
		if result.ExitError == "" && (result.ArtifactOK || result.Skipped) {
			okCount++
			continue
		}
		eventType = "round.incomplete"
	}
	if indexPath, err := writeRoundIndex(opts.Idea, opts.RoundLabel, results); err != nil {
		now := time.Now().UTC()
		warning := "round index write failed: " + err.Error()
		results = append(results, Result{
			AgentID:     "runner/index",
			OutputPath:  indexPath,
			CompletedAt: now,
			Warning:     warning,
		})
		_ = opts.Store.Append(store.Event{
			Time: now,
			Type: "round.index_failed",
			Data: map[string]any{
				"idea":     opts.Idea.Slug,
				"round":    opts.RoundLabel,
				"artifact": indexPath,
				"error":    err.Error(),
			},
		})
	} else {
		_ = opts.Store.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "round.index_written",
			Data: map[string]any{
				"idea":     opts.Idea.Slug,
				"round":    opts.RoundLabel,
				"artifact": indexPath,
			},
		})
	}
	if err := opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: eventType,
		Data: map[string]any{
			"idea":      opts.Idea.Slug,
			"round":     opts.RoundLabel,
			"completed": okCount,
			"total":     len(results),
		},
	}); err != nil {
		results = append(results, Result{
			AgentID:     "runner",
			CompletedAt: time.Now().UTC(),
			ExitError:   "round event append failed: " + err.Error(),
		})
	}
	return results
}

// appendSegmentStarted records a run.segment_started boundary and returns the
// new monotonic segment id. The projection (runstate.ProjectEvents) resets the
// targeted agents to pending for this segment so a stale terminal badge from a
// prior segment never lingers after continue/resume/retry.
func appendSegmentStarted(opts Options, reason string, targets []string) string {
	seg := nextSegmentID(opts.Store)
	_ = opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "run.segment_started",
		Data: map[string]any{
			"segment_id": seg,
			"reason":     reason,
			"round":      opts.RoundLabel,
			"targets":    targets,
		},
	})
	return seg
}

// nextSegmentID returns the next monotonic segment-NNNN id by counting prior
// run.segment_started events. Round-runs are sequential per run, so the count is
// stable. Falls back to segment-0001 when the log cannot be read.
func nextSegmentID(s store.Store) string {
	events, err := s.Load()
	if err != nil {
		return "segment-0001"
	}
	n := 0
	for _, e := range events {
		if e.Type == "run.segment_started" {
			n++
		}
	}
	return fmt.Sprintf("segment-%04d", n+1)
}

// segmentReason classifies why a round-run is starting, for the audit trail.
// Only the first deliberation round is "initial"; everything else is "continue"
// (RunFixup overrides this with "retry").
func segmentReason(opts Options) string {
	if opts.Phase == "" && opts.Round <= 1 {
		return "initial"
	}
	return "continue"
}

// agentIDs extracts the stable IDs from discovered agents (segment targets).
func agentIDs(selected []agents.Discovery) []string {
	ids := make([]string, 0, len(selected))
	for _, a := range selected {
		ids = append(ids, a.ID)
	}
	return ids
}

func selectedAgents(participants []string, discovered []agents.Discovery) []agents.Discovery {
	byID := make(map[string]agents.Discovery, len(discovered))
	for _, agent := range discovered {
		if agent.Found {
			byID[agent.ID] = agent
		}
	}

	var selected []agents.Discovery
	for _, participant := range participants {
		if agent, ok := byID[participant]; ok {
			selected = append(selected, agent)
		}
	}
	return selected
}

func runAgent(parent context.Context, opts Options, agent agents.Discovery) Result {
	now := time.Now().UTC()
	agentDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "agents", agent.ID)
	stdoutPath := filepath.Join(agentDir, "stdout.log")
	stderrPath := filepath.Join(agentDir, "stderr.log")
	outputPath := filepath.Join(opts.Idea.Path, opts.RoundLabel, agent.ID+".md")
	if opts.ArtifactName != "" {
		outputPath = filepath.Join(opts.Idea.Path, opts.ArtifactName)
	}
	result := Result{
		AgentID:    agent.ID,
		OutputPath: outputPath,
		StdoutPath: stdoutPath,
		StderrPath: stderrPath,
		StartedAt:  now,
	}

	if _, err := os.Stat(outputPath); err == nil && !opts.Overwrite {
		result.Skipped = true
		result.SkipReason = "artifact already exists"
		result.CompletedAt = time.Now().UTC()
		result.Duration = result.CompletedAt.Sub(result.StartedAt)
		if err := opts.Store.Append(store.Event{
			Time: result.CompletedAt,
			Type: "agent.skipped",
			Data: map[string]any{"agent": agent.ID, "reason": result.SkipReason, "artifact": outputPath, "segment_id": opts.SegmentID},
		}); err != nil {
			result.ExitError = "event append failed: " + err.Error()
		}
		return result
	}

	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		return failEarly(opts, result, err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return failEarly(opts, result, err)
	}

	questionsDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "questions")
	prompt, err := buildPromptForRound(agent, opts, outputPath, questionsDir)
	if err != nil {
		return failEarly(opts, result, err)
	}

	if agents.LaunchModeOrDefault(agent.LaunchMode) == agents.LaunchACP {
		return runACPAgent(parent, opts, agent, result, outputPath, stdoutPath, stderrPath, prompt)
	}

	ctx, cancel := context.WithTimeout(parent, timeoutForAgent(opts.Timeout, agent))
	defer cancel()

	if err := opts.Store.Append(store.Event{
		Time: now,
		Type: "agent.started",
		Data: map[string]any{
			"agent":      agent.ID,
			"artifact":   outputPath,
			"stdout":     stdoutPath,
			"stderr":     stderrPath,
			"segment_id": opts.SegmentID,
		},
	}); err != nil {
		return failEarly(opts, result, fmt.Errorf("event append failed: %w", err))
	}

	// Register this attempt so Handle.KillAgent can cancel just this agent.
	opts.tracker.register(agent.ID, opts.SegmentID, "round", "", cancel)
	err = execAgentProcess(ctx, opts.Root, agent, prompt, stdoutPath, stderrPath)
	killed := opts.tracker.finish(agent.ID)
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)
	if err != nil {
		result.ExitError = err.Error()
		if ctx.Err() != nil {
			result.ExitError = ctx.Err().Error()
		}
	}
	if killed {
		result.Killed = true
		result.ExitError = "killed by user"
	}

	// stdout-capture fallback: some print-only CLIs (e.g. agy --print) emit the
	// artifact to stdout instead of writing the file. If the file is absent but
	// captured stdout is a plausible artifact (starts with the YAML frontmatter
	// fence), persist it as the agent-authored artifact. Strict "---" validation
	// keeps narration from becoming a protocol file.
	if _, statErr := os.Stat(outputPath); os.IsNotExist(statErr) {
		if data, readErr := os.ReadFile(stdoutPath); readErr == nil && firstLineIsFence(data) {
			// Validate a CANDIDATE before placing it at the protocol path, so a
			// malformed print-only response never poisons the artifact path.
			tmp := outputPath + ".stdout-candidate"
			if writeErr := os.WriteFile(tmp, data, 0o644); writeErr == nil {
				if validateArtifactForPhase(opts, tmp, agent.ID) == nil {
					if renameErr := os.Rename(tmp, outputPath); renameErr == nil {
						result.Warning = "artifact recovered from stdout (print-only agent)"
						_ = opts.Store.Append(store.Event{
							Time: time.Now().UTC(),
							Type: "agent.stdout_fallback",
							Data: map[string]any{"agent": agent.ID, "artifact": outputPath},
						})
					} else {
						_ = os.Remove(tmp)
					}
				} else {
					_ = os.Remove(tmp) // invalid candidate -> leave no protocol artifact
				}
			}
		}
	}

	if _, statErr := os.Stat(outputPath); statErr == nil {
		if validateErr := validateArtifactForPhase(opts, outputPath, agent.ID); validateErr != nil {
			result.ExitError = combineError(result.ExitError, validateErr)
		} else {
			result.ArtifactOK = true
		}
	}

	eventType := "agent.finished"
	if result.ExitError != "" || !result.ArtifactOK {
		eventType = "agent.failed"
	}
	if err := opts.Store.Append(store.Event{
		Time: result.CompletedAt,
		Type: eventType,
		Data: map[string]any{
			"agent":       agent.ID,
			"artifact":    outputPath,
			"artifact_ok": result.ArtifactOK,
			"duration_ms": result.Duration.Milliseconds(),
			"error":       result.ExitError,
			"segment_id":  opts.SegmentID,
		},
	}); err != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", err))
	}

	return result
}

func failEarly(opts Options, result Result, err error) Result {
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)
	result.ExitError = err.Error()
	if eventErr := opts.Store.Append(store.Event{
		Time: result.CompletedAt,
		Type: "agent.failed",
		Data: map[string]any{
			"agent":       result.AgentID,
			"artifact":    result.OutputPath,
			"artifact_ok": false,
			"duration_ms": result.Duration.Milliseconds(),
			"error":       result.ExitError,
			"segment_id":  opts.SegmentID,
		},
	}); eventErr != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", eventErr))
	}
	return result
}

func combineError(primary string, err error) string {
	if err == nil {
		return primary
	}
	if primary == "" {
		return err.Error()
	}
	return primary + "; " + err.Error()
}

func BuildRoundOnePrompt(agent agents.Discovery, idea protocol.IdeaStatus, task, outputPath, questionsDir string) (string, error) {
	promptData, err := os.ReadFile(filepath.Join(idea.Path, "00-prompt.md"))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`You are %s, a participant in a Parley Deck cooperation round.

Rules:
- Create exactly this file and no other protocol artifact: %s
- Immediately use your file-writing tool to create that exact file with its required frontmatter; do not explore the workspace first, and report a blocker instead of narrating.
- Do not edit any other agent's file.
- Do not overwrite the file if it already exists; report a blocker instead.
- Do not read or reference other agents' round-01 answers.
- Write the complete file, including YAML frontmatter. The first line of the file must be exactly "---".
- Return only a short confirmation with the path written.
- Be concrete, concise, and state trade-offs.
- If you are blocked by missing human input, create one JSON question file under: %s
- Question files use this shape: {"id":"<unique-id>","agent":"%s","prompt":"<question>","details":"<context>","default_answer":"<safe default if any>","risk":"low|normal|high","status":"open","answer":"","created_at":"<RFC3339 time>","answered_at":"0001-01-01T00:00:00Z"}
- If you choose to wait for an answer, poll your question file until status is answered or auto_answered. Otherwise proceed with an explicit assumption in your artifact.

Effective launch config:
- model: %s
- thinking/reasoning/effort/profile: %s
- speed: %s
- sandbox: %s
- approval: %s
- timeoutMs: %d

Idea prompt:
%s

Required file shape:
---
agent: %s
idea: %s
round: 1
date: %s
---

## Summary
## Proposed approach
## Concerns / open questions
## Risks
`, agent.ID, outputPath, questionsDir, agent.ID,
		runtimeValue(agent.Model),
		runtimeValue(firstNonEmpty(agent.Reasoning, agent.Profile)),
		runtimeValue(firstNonEmpty(agent.Speed, agents.DefaultSpeed)),
		runtimeValue(agent.SandboxMode),
		runtimeValue(agent.ApprovalPolicy),
		timeoutMSForAgent(agent),
		string(promptData), agent.ID, idea.Slug, time.Now().Format("2006-01-02")), nil
}

func roundNumber(opts Options) int {
	if opts.Round > 0 {
		return opts.Round
	}
	return 1
}

func roundLabel(n int) string {
	return fmt.Sprintf("round-%02d", n)
}

// firstLineIsFence reports whether the first line is exactly the YAML frontmatter
// fence "---". Unlike a prefix check it does not accept leading narration, so a
// stdout stream that merely contains "---" somewhere is rejected.
func firstLineIsFence(data []byte) bool {
	s := string(data)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return strings.TrimRight(s, "\r ") == "---"
}

// RunRound runs cross-review round N (N>=2) for an idea; round 1 uses
// RunRoundOne. It seeds each participant with the prior rounds' artifacts and
// writes round-NN/<agent>.md, reusing the same launch/validation machinery.
func RunRound(ctx context.Context, opts Options) []Result {
	if opts.Round < 2 {
		opts.Round = 2
	}
	if opts.RoundLabel == "" {
		opts.RoundLabel = roundLabel(opts.Round)
	}
	return RunRoundOne(ctx, opts)
}

func buildPromptForRound(agent agents.Discovery, opts Options, outputPath, questionsDir string) (string, error) {
	switch opts.Phase {
	case "implementation":
		return BuildImplementationPrompt(agent, opts.Idea, outputPath), nil
	case "review":
		ctx, err := gatherReviewContext(opts.Idea.Path, roundNumber(opts))
		if err != nil {
			return "", err
		}
		return BuildReviewPrompt(agent, opts.Idea, roundNumber(opts), outputPath, ctx), nil
	case "review-consensus":
		ctx, err := gatherReviewContext(opts.Idea.Path, roundNumber(opts)+1)
		if err != nil {
			return "", err
		}
		return BuildReviewConsensusPrompt(agent, opts.Idea, outputPath, ctx), nil
	}
	if roundNumber(opts) <= 1 {
		return BuildRoundOnePrompt(agent, opts.Idea, opts.Task, outputPath, questionsDir)
	}
	prior, err := gatherPriorRounds(opts.Idea.Path, roundNumber(opts))
	if err != nil {
		return "", err
	}
	return BuildRoundPrompt(agent, opts.Idea, roundNumber(opts), outputPath, questionsDir, prior), nil
}

// gatherPriorRounds concatenates every participant artifact from rounds 1..N-1
// so a cross-review round can respond to them.
func gatherPriorRounds(ideaPath string, round int) (string, error) {
	var b strings.Builder
	for r := 1; r < round; r++ {
		dir := filepath.Join(ideaPath, roundLabel(r))
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "_index.md" {
				continue
			}
			names = append(names, e.Name())
		}
		sort.Strings(names)
		for _, name := range names {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "\n===== %s/%s =====\n%s\n", roundLabel(r), name, string(data))
		}
	}
	return b.String(), nil
}

// BuildRoundPrompt builds the cross-review prompt for round N>=2. Unlike round
// 1, participants are given every prior-round artifact and asked to respond to
// each other and converge toward consensus.
func BuildRoundPrompt(agent agents.Discovery, idea protocol.IdeaStatus, round int, outputPath, questionsDir, prior string) string {
	others := make([]string, 0, len(idea.Participants))
	for _, p := range idea.Participants {
		if p != agent.ID {
			others = append(others, p)
		}
	}
	var headings strings.Builder
	for _, other := range others {
		fmt.Fprintf(&headings, "### @%s\n", other)
	}
	return fmt.Sprintf(`You are %s, a participant in a Parley Deck cross-review round %d.

Rules:
- Create exactly this file and no other protocol artifact: %s
- Immediately use your file-writing tool to create that exact file with its required frontmatter; do not explore the workspace first, and report a blocker instead of narrating.
- Do not edit any other agent's file.
- Do not overwrite the file if it already exists; report a blocker instead.
- READ every prior-round artifact below and respond to the other participants by name: where you agree, where you disagree, what you refine. Converge toward consensus.
- Under "## Responses to other participants", you MUST include one "### @<agent-id>" subsection for EACH other participant (%s) and address that participant specifically.
- Write the complete file, including YAML frontmatter. The first line of the file must be exactly "---".
- If you are blocked by missing human input, create one JSON question file under: %s
- Be concrete, concise, and state trade-offs.

Required file shape:
---
agent: %s
idea: %s
round: %d
date: %s
responding-to: [prior round artifacts]
---

## Summary
## Responses to other participants
%s## Refined position
## Remaining disagreements

Prior rounds (read these):
%s
`, agent.ID, round, outputPath, strings.Join(others, ", "), questionsDir, agent.ID, idea.Slug, round, time.Now().Format("2006-01-02"), headings.String(), prior)
}

// execAgentProcess builds the agent command, wires its stdout/stderr to the
// given log paths (and stdin for stdin-prompt agents), and runs it to
// completion. Shared by the round path (runAgent) and the steer path
// (runSteerAgent) so both honor the same exec + cancellation semantics.
func execAgentProcess(ctx context.Context, root string, agent agents.Discovery, prompt, stdoutPath, stderrPath string) error {
	cmd, cleanup, err := CommandFor(ctx, root, agent, prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return err
	}
	cmd.Dir = root
	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return err
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return err
	}
	defer stderrFile.Close()
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}
	return cmd.Run()
}

func CommandFor(ctx context.Context, root string, agent agents.Discovery, prompt string) (*exec.Cmd, func(), error) {
	args := make([]string, 0, len(agent.HeadlessArgs))
	for _, arg := range agent.HeadlessArgs {
		switch arg {
		case "{root}":
			args = append(args, root)
		case "{prompt}":
			args = append(args, prompt)
		default:
			args = append(args, arg)
		}
	}
	cmd := exec.CommandContext(ctx, agent.Path, args...)
	cleanup := func() {}
	if agent.IsolateHome {
		env, remove, err := isolatedAgentHome(agent)
		if err != nil {
			return nil, nil, err
		}
		cleanup = remove
		cmd.Env = append(os.Environ(), env...)
		if agent.ID == "hermes" {
			cmd.Env = append(cmd.Env, "HERMES_ACCEPT_HOOKS=1", "HERMES_SESSION_SOURCE=parley")
		}
	}
	return cmd, cleanup, nil
}

func timeoutForAgent(override time.Duration, agent agents.Discovery) time.Duration {
	if override > 0 {
		return override
	}
	return time.Duration(timeoutMSForAgent(agent)) * time.Millisecond
}

func timeoutMSForAgent(agent agents.Discovery) int {
	if agent.TimeoutMS > 0 {
		return agent.TimeoutMS
	}
	return agents.DefaultTimeoutMS
}

func runtimeValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return agents.CLIDefault
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func isolatedAgentHome(agent agents.Discovery) ([]string, func(), error) {
	switch agent.ID {
	case "gemini":
		home, err := isolatedGeminiHome()
		if err != nil {
			return nil, nil, err
		}
		return isolatedHomeEnv(agent, home, "GEMINI_CLI_HOME"), func() { _ = os.RemoveAll(home) }, nil
	case "hermes":
		home, err := isolatedHermesHome()
		if err != nil {
			return nil, nil, err
		}
		return isolatedHomeEnv(agent, home, "HERMES_HOME"), func() { _ = os.RemoveAll(home) }, nil
	default:
		if len(agent.IsolatedHomeEnv) == 0 {
			return nil, nil, fmt.Errorf("no isolated home strategy for %s", agent.ID)
		}
		home, err := os.MkdirTemp("", "parley-"+agent.ID+"-home.*")
		if err != nil {
			return nil, nil, err
		}
		return isolatedHomeEnv(agent, home, ""), func() { _ = os.RemoveAll(home) }, nil
	}
}

func isolatedHomeEnv(agent agents.Discovery, home, fallbackKey string) []string {
	if len(agent.IsolatedHomeEnv) == 0 && fallbackKey != "" {
		return []string{fallbackKey + "=" + home}
	}
	env := make([]string, 0, len(agent.IsolatedHomeEnv))
	for key, template := range agent.IsolatedHomeEnv {
		value := strings.ReplaceAll(template, "{tempdir}", home)
		env = append(env, key+"="+value)
	}
	return env
}

func isolatedGeminiHome() (string, error) {
	sourceHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base, err := os.MkdirTemp("", "parley-gemini-home.*")
	if err != nil {
		return "", err
	}
	geminiDir := filepath.Join(base, ".gemini")
	if err := os.MkdirAll(geminiDir, 0o700); err != nil {
		_ = os.RemoveAll(base)
		return "", err
	}
	copied := 0
	for _, name := range []string{"oauth_creds.json", "google_accounts.json"} {
		source := filepath.Join(sourceHome, ".gemini", name)
		data, err := os.ReadFile(source)
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(geminiDir, name), data, 0o600); err != nil {
			_ = os.RemoveAll(base)
			return "", err
		}
		copied++
	}
	if copied == 0 && os.Getenv("GEMINI_API_KEY") == "" && os.Getenv("GOOGLE_API_KEY") == "" {
		_ = os.RemoveAll(base)
		return "", fmt.Errorf("no Gemini OAuth files found in %s and no Gemini API key is set", filepath.Join(sourceHome, ".gemini"))
	}
	settings := []byte(`{"security":{"auth":{"selectedType":"oauth-personal"}}}` + "\n")
	if err := os.WriteFile(filepath.Join(geminiDir, "settings.json"), settings, 0o600); err != nil {
		_ = os.RemoveAll(base)
		return "", err
	}
	return base, nil
}

func isolatedHermesHome() (string, error) {
	sourceHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	sourceDir := filepath.Join(sourceHome, ".hermes")
	base, err := os.MkdirTemp("", "parley-hermes-home.*")
	if err != nil {
		return "", err
	}
	for _, dir := range []string{"logs", "sessions", "home"} {
		if err := os.MkdirAll(filepath.Join(base, dir), 0o700); err != nil {
			_ = os.RemoveAll(base)
			return "", err
		}
	}

	copied := 0
	for _, name := range []string{"config.yaml", ".env", "auth.json", "SOUL.md"} {
		source := filepath.Join(sourceDir, name)
		data, err := os.ReadFile(source)
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(base, name), data, 0o600); err != nil {
			_ = os.RemoveAll(base)
			return "", err
		}
		copied++
	}
	if copied == 0 {
		_ = os.RemoveAll(base)
		return "", fmt.Errorf("no Hermes config files found in %s", sourceDir)
	}
	return base, nil
}
