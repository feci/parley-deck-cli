package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

type Options struct {
	Root   string
	RunID  string
	Idea   protocol.IdeaStatus
	Task   string
	Agents []agents.Discovery
	// RosterMapping is the roster-ID -> family map (from `[roster.*] adapter` in the
	// deck config) used to resolve participants that are roster IDs (e.g. claude-1)
	// rather than bare family ids. nil/empty falls back to exact spec-ID matching.
	RosterMapping map[string]string
	Timeout       time.Duration
	Store         store.Store
	Round         int
	RoundLabel    string
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
	// StrictGate, when set, tells the Phase-7 review-consensus prompt to also emit the
	// machine-readable closing_review_round + strict_gate_clean fields the driver's
	// strict_gate close path checks (LE-2). Set by driverImplOps.DraftReviewConsensus.
	StrictGate bool
	// SegmentID tags agent events with the current run.segment_started segment so
	// the projection can scope state to the active segment (fixes stale terminal
	// badges after continue/resume/retry). Set by the round-run entry points via
	// appendSegmentStarted; empty for legacy/unsegmented callers.
	SegmentID string
	// tracker, when set (by RunRoundOneAsync), registers each agent attempt's
	// cancel func so the live Handle can KillAgent one agent. nil on the sync path.
	tracker *Handle
	// publishArtifact, when set (review snapshots, consensus D9), moves a
	// VALIDATED artifact from the snapshot to its canonical path. It ALWAYS
	// returns the live canonical path (terminal events must report it even on
	// failure); a publish failure turns the result into a failure and the
	// snapshot is retained for recovery (review fix 2).
	publishArtifact func(snapshotPath string) (string, error)
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
	// AgentExit records an ordinary nonzero exit that was OVERRIDDEN by a
	// validated artifact (artifact-wins, consensus D7): the step succeeded,
	// the exit code is preserved for display. 0 when unset.
	AgentExit int
	// AgentExitKind distinguishes the artifact-wins source: "exec" for a
	// process exit code, "acp_error" for an ACP prompt error after the session
	// opened. Empty when AgentExit is unset.
	AgentExitKind string
	// FailureClass/RecoveryHint classify a failure for operators (consensus D5).
	FailureClass string
	RecoveryHint string
}

// Success reports whether the step counts as succeeded for gating purposes:
// skipped (artifact pre-existed) or finished with a validated artifact —
// including artifact-wins completions that carry a nonzero AgentExit.
func (r Result) Success() bool {
	if r.Killed {
		return false
	}
	if r.Skipped {
		return true
	}
	return r.ExitError == "" && r.ArtifactOK
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

	selected := selectedAgents(opts.Idea.Participants, opts.Agents, resolveMapping(opts))
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
		if result.Success() {
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

// RosterMappingLoader, when set by the app at startup, lazily loads the
// roster-ID -> family map for a root so run paths that do not set
// Options.RosterMapping still resolve roster-ID participants. Injected (rather than
// importing config here) to avoid a runner -> config import cycle.
var RosterMappingLoader func(root string) map[string]string

func resolveMapping(opts Options) map[string]string {
	if opts.RosterMapping != nil {
		return opts.RosterMapping
	}
	if RosterMappingLoader != nil {
		return RosterMappingLoader(opts.Root)
	}
	return nil
}

// selectedAgents resolves each participant/roster ID to a discovered agent via
// agents.ResolveParticipant (exact spec-ID -> [roster.*] mapping -> fail closed),
// carrying the participant string as the agent identity and the family as its
// adapter. Unresolvable participants are skipped (matching the pre-split behavior
// of an absent participant); `parley roster init` + preflight surface the gap.
func selectedAgents(participants []string, discovered []agents.Discovery, mapping map[string]string) []agents.Discovery {
	var selected []agents.Discovery
	for _, participant := range participants {
		if resolved, err := agents.ResolveParticipant(participant, discovered, mapping); err == nil {
			selected = append(selected, resolved)
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

	if err := fsutil.MkdirAllResilient(agentDir, 0o755); err != nil {
		return failEarly(opts, result, err)
	}

	// HITL questions are polled from the LIVE runs dir; capture before any
	// snapshot root swap.
	questionsDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", opts.RunID, "questions")

	// Phase 6 review isolation (consensus D9): the reviewer reads a disposable
	// shared-clone checkout on local tmp; its artifact is moved back to the
	// canonical deck path after validation. Any unavailability falls back to
	// the live tree with a loud event.
	if opts.Phase == "review" {
		snap, snapErr := CreateReviewSnapshot(opts.Root, opts.Idea.Slug, opts.RoundLabel, agent.ID, opts.RunID)
		if snapErr != nil {
			_ = opts.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "review.snapshot_fallback",
				Data: map[string]any{"agent": agent.ID, "reason": snapErr.Error(), "segment_id": opts.SegmentID},
			})
		} else {
			keepForRecovery := false
			defer func() {
				if keepForRecovery {
					snap.Abandon() // retain the dir; the artifact inside is the recovery copy
				} else {
					snap.Cleanup()
				}
			}()
			_ = opts.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "review.snapshot_created",
				Data: map[string]any{"agent": agent.ID, "sha": snap.SHA, "mode": snap.Mode, "dir": snap.Dir, "segment_id": opts.SegmentID},
			})
			liveRoot, liveOutput := opts.Root, outputPath
			relOutput, relErr := filepath.Rel(liveRoot, liveOutput)
			relIdea, relIdeaErr := filepath.Rel(liveRoot, opts.Idea.Path)
			if relErr == nil && relIdeaErr == nil {
				opts.Root = snap.Dir
				opts.Idea.Path = filepath.Join(snap.Dir, relIdea)
				outputPath = filepath.Join(snap.Dir, relOutput)
				result.OutputPath = outputPath
				opts.publishArtifact = func(snapshotPath string) (string, error) {
					if err := snap.MoveArtifactBack(relOutput, liveOutput); err != nil {
						keepForRecovery = true // the snapshot copy is now the recovery artifact
						_ = opts.Store.Append(store.Event{
							Time: time.Now().UTC(),
							Type: "review.snapshot_artifact_move_failed",
							Data: map[string]any{"agent": agent.ID, "snapshot_path": snapshotPath, "recovery_path": liveOutput, "error": err.Error(), "segment_id": opts.SegmentID},
						})
						return liveOutput, err
					}
					return liveOutput, nil
				}
			}
		}
	}

	if err := fsutil.MkdirAllResilient(filepath.Dir(outputPath), 0o755); err != nil {
		return failEarly(opts, result, err)
	}
	prompt, err := buildPromptForRound(agent, opts, outputPath, questionsDir)
	if err != nil {
		return failEarly(opts, result, err)
	}

	// preexisted guards the retry's move-aside: never disturb an artifact that
	// existed before this runner call (consensus D3; Overwrite runs land here).
	_, statErr := os.Stat(outputPath)
	preexisted := statErr == nil

	if agents.LaunchModeOrDefault(agent.LaunchMode) == agents.LaunchACP {
		// ACP attempts share the exec retry contract (review fix 1c): retry
		// ONCE, only for a first-output watchdog kill.
		for attemptID := 1; ; attemptID++ {
			res := runACPAgent(parent, opts, agent, result, outputPath, stdoutPath, stderrPath, prompt, attemptID)
			if attemptID == 1 && res.FailureClass == "no_first_output" && !res.Killed {
				if !preexisted {
					if _, err := os.Stat(outputPath); err == nil && validateArtifactForPhase(opts, outputPath, agent.ID) != nil {
						moveAsideInvalidArtifact(outputPath)
					}
				}
				continue
			}
			return res
		}
	}

	// Attempt loop: retry ONCE, and only for a first-output watchdog kill.
	for attemptID := 1; ; attemptID++ {
		attempt := runExecAttempt(parent, opts, agent, result, outputPath, stdoutPath, stderrPath, prompt, attemptID)
		if attemptID == 1 && attempt.watchdog == "no_first_output" && !attempt.result.Killed {
			// The killed attempt's terminal agent.failed is already appended
			// (before this retry's agent.started — durable kill targets the
			// newest attempt). Move an invalid attempt-1 artifact aside.
			if !preexisted {
				if _, err := os.Stat(outputPath); err == nil && validateArtifactForPhase(opts, outputPath, agent.ID) != nil {
					moveAsideInvalidArtifact(outputPath)
				}
			}
			continue
		}
		return attempt.result
	}
}

// moveAsideInvalidArtifact relocates an invalid attempt-1 artifact before the
// retry. The destination is guaranteed not to exist (unique suffix when the
// stable name is taken — never overwrite an earlier recovery file); when the
// rename itself fails, the invalid artifact — created by THIS attempt, the
// preexisted guard excludes foreign files — is removed so it cannot linger on
// the canonical path (review fix 3).
func moveAsideInvalidArtifact(outputPath string) {
	dest := outputPath + ".attempt-1.invalid"
	if _, err := os.Stat(dest); err == nil {
		dest = fmt.Sprintf("%s.%d", dest, time.Now().UnixNano())
	}
	if err := os.Rename(outputPath, dest); err != nil {
		_ = os.Remove(outputPath)
	}
}

// execAttempt is one supervised exec-mode invocation's outcome.
type execAttempt struct {
	result   Result
	watchdog string // "no_first_output" | "stalled" | "" — which guard killed it
}

func runExecAttempt(parent context.Context, opts Options, agent agents.Discovery, base Result, outputPath, stdoutPath, stderrPath, prompt string, attemptID int) execAttempt {
	result := base
	result.StartedAt = time.Now().UTC()
	hardTimeout := timeoutForAgent(opts.Timeout, agent)
	ctx, cancel := context.WithTimeout(parent, hardTimeout)
	defer cancel()

	// Register this attempt so Handle.KillAgent can cancel just this agent (the
	// cancel triggers the supervised wait's kill → group kill).
	opts.tracker.register(agent.ID, opts.SegmentID, "round", "", cancel)
	marker := fmt.Sprintf("%s:%s:%d", opts.RunID, agent.ID, attemptID)
	// agent.started is emitted once the process is live, enriched with its durable
	// process identity so a restarted parley can re-attribute and kill it.
	onStarted := func(sp procctl.Spawned) {
		_ = opts.Store.Append(store.Event{
			Time: time.Now().UTC(), // stamp at actual process start, not before setup
			Type: "agent.started",
			Data: map[string]any{
				"agent":       agent.ID,
				"artifact":    outputPath,
				"stdout":      stdoutPath,
				"stderr":      stderrPath,
				"segment_id":  opts.SegmentID,
				"attempt_id":  attemptID,
				"pid":         sp.PID,
				"pgid":        sp.PGID,
				"boot_id":     sp.BootID,
				"proc_start":  sp.ProcStart,
				"proc_marker": sp.Marker,
				"command":     sp.Command,
			},
		})
	}
	act := &activityTracker{}
	cfg := supervisionForAgent(agent, hardTimeout)
	hooks := supervisionHooks{
		onHeartbeat: func(snap activitySnapshot, elapsed time.Duration) {
			_ = opts.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "agent.heartbeat",
				Data: map[string]any{
					"agent":                agent.ID,
					"segment_id":           opts.SegmentID,
					"attempt_id":           attemptID,
					"phase":                phaseOrDefault(opts.Phase),
					"launch":               agents.LaunchHeadless,
					"elapsed_ms":           elapsed.Milliseconds(),
					"timeout_ms":           hardTimeout.Milliseconds(),
					"stdout_bytes":         snap.StdoutBytes,
					"stderr_bytes":         snap.StderrBytes,
					"last_activity_ms_ago": activityAgeMS(snap),
				},
			})
		},
		onWatchdog: func(kind string, snap activitySnapshot, elapsed time.Duration) {
			// Appended BEFORE the kill fires (consensus D1): the event log names
			// the killer ahead of any durable-kill attribution race.
			action := "failed"
			if kind == "no_first_output" && attemptID == 1 {
				action = "retrying"
			}
			_ = opts.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "agent." + kind,
				Data: map[string]any{
					"agent":        agent.ID,
					"segment_id":   opts.SegmentID,
					"attempt_id":   attemptID,
					"elapsed_ms":   elapsed.Milliseconds(),
					"stdout_bytes": snap.StdoutBytes,
					"stderr_bytes": snap.StderrBytes,
					"action":       action,
				},
			})
		},
	}
	_, err := execAgentProcess(ctx, opts.Root, opts.RunID, agent.ID, marker, agent, prompt, stdoutPath, stderrPath, onStarted, act, cfg, hooks)
	killed := opts.tracker.finish(agent.ID)
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	watchdog := ""
	switch {
	case errors.Is(err, errNoFirstOutput):
		watchdog = "no_first_output"
	case errors.Is(err, errStalled):
		watchdog = "stalled"
	}
	finalizeExecResult(opts, &result, agent, err, killed, ctx.Err(), watchdog, attemptID, outputPath, stdoutPath, stderrPath)
	return execAttempt{result: result, watchdog: watchdog}
}

// finalizeExecResult applies the stdout fallback, validates the artifact, and
// makes the single terminal decision (consensus D7): a VALIDATED artifact with
// an ordinary nonzero exit finishes with agent_exit; watchdog finals, the hard
// timeout, and user kills always fail; everything else fails with a class+hint.
func finalizeExecResult(opts Options, result *Result, agent agents.Discovery, runErr error, killed bool, ctxErr error, watchdog string, attemptID int, outputPath, stdoutPath, stderrPath string) {
	if runErr != nil {
		result.ExitError = runErr.Error()
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

	// Review snapshots: a validated artifact is published (moved back) to its
	// canonical path BEFORE the terminal event, which must report the live
	// path (consensus D9). A publish failure fails the step; the snapshot is
	// retained for recovery.
	if result.ArtifactOK && opts.publishArtifact != nil {
		livePath, err := opts.publishArtifact(outputPath)
		if livePath != "" {
			// Terminal events report the LIVE canonical path even when the
			// move-back failed (the snapshot path travels only as recovery
			// metadata in review.snapshot_artifact_move_failed).
			outputPath = livePath
			result.OutputPath = livePath
		}
		if err != nil {
			result.ArtifactOK = false
			result.ExitError = combineError(result.ExitError, fmt.Errorf("snapshot artifact move-back: %w", err))
		}
	}

	// Artifact-wins: an ordinary process exit error is overridden by a
	// validated artifact. Watchdog finals, the hard timeout, and user kills
	// are cancellations — they always win over the artifact.
	cancellation := killed || watchdog != "" || errors.Is(ctxErr, context.DeadlineExceeded) || errors.Is(ctxErr, context.Canceled)
	var exitErr *exec.ExitError
	if result.ArtifactOK && !cancellation && runErr != nil && errors.As(runErr, &exitErr) {
		result.AgentExit = exitErr.ExitCode()
		result.AgentExitKind = "exec"
		result.ExitError = ""
	}

	failed := result.ExitError != "" || !result.ArtifactOK
	if failed && result.ExitError == "" {
		result.ExitError = "artifact missing or invalid"
	}
	data := map[string]any{
		"agent":       agent.ID,
		"artifact":    outputPath,
		"artifact_ok": result.ArtifactOK,
		"duration_ms": result.Duration.Milliseconds(),
		"error":       result.ExitError,
		"segment_id":  opts.SegmentID,
		"attempt_id":  attemptID,
	}
	eventType := "agent.finished"
	if failed {
		eventType = "agent.failed"
		class, hint := terminalFailureClass(watchdog, ctxErr, killed, stderrPath, stdoutPath, result.ExitError)
		result.FailureClass = class
		result.RecoveryHint = hint
		data["failure_class"] = class
		data["recovery_hint"] = hint
		if exitErr != nil {
			data["exit_code"] = exitErr.ExitCode()
		}
		if sig := signalFromError(result.ExitError); sig != "" {
			data["signal"] = sig
		}
		data["stderr_tail_bytes"] = len(tailOfFile(stderrPath, failTailBytes))
	} else if result.AgentExitKind != "" {
		data["agent_exit"] = result.AgentExit
		data["agent_exit_kind"] = result.AgentExitKind
	}
	if err := opts.Store.Append(store.Event{Time: result.CompletedAt, Type: eventType, Data: data}); err != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", err))
	}
}

// terminalFailureClass picks the failure class: watchdog/timeout/kill causes
// are authoritative (consensus D5); everything else goes through the regex
// classifier over the bounded log tails.
func terminalFailureClass(watchdog string, ctxErr error, killed bool, stderrPath, stdoutPath, exitError string) (string, string) {
	switch {
	case killed:
		return "killed", "Killed by user request."
	case watchdog != "":
		return watchdog, watchdogHints[watchdog]
	case errors.Is(ctxErr, context.DeadlineExceeded):
		return "timeout", watchdogHints["timeout"]
	}
	return classifyFailure(stderrPath, stdoutPath, exitError)
}

func signalFromError(errText string) string {
	if i := strings.Index(errText, "signal: "); i >= 0 {
		rest := errText[i+len("signal: "):]
		if j := strings.IndexAny(rest, " \n("); j > 0 {
			return rest[:j]
		}
		return rest
	}
	return ""
}

func activityAgeMS(snap activitySnapshot) int64 {
	if snap.LastActivity.IsZero() {
		return -1
	}
	return time.Since(snap.LastActivity).Milliseconds()
}

func phaseOrDefault(phase string) string {
	if phase == "" {
		return "deliberation"
	}
	return phase
}

func failEarly(opts Options, result Result, err error) Result {
	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)
	result.ExitError = err.Error()
	// Setup failures get the same classified payload as terminal failures
	// (review fix 5): the regex runs over the error text plus whatever log
	// tails already exist.
	result.FailureClass, result.RecoveryHint = classifyFailure(result.StderrPath, result.StdoutPath, result.ExitError)
	if eventErr := opts.Store.Append(store.Event{
		Time: result.CompletedAt,
		Type: "agent.failed",
		Data: map[string]any{
			"agent":         result.AgentID,
			"artifact":      result.OutputPath,
			"artifact_ok":   false,
			"duration_ms":   result.Duration.Milliseconds(),
			"error":         result.ExitError,
			"segment_id":    opts.SegmentID,
			"failure_class": result.FailureClass,
			"recovery_hint": result.RecoveryHint,
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
		return BuildReviewConsensusPrompt(agent, opts.Idea, outputPath, ctx, opts.StrictGate), nil
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
// execAgentProcess spawns an agent in its OWN process group (so a kill reaps the
// whole tree), captures its durable identity (pid/pgid/boot/start/command) via
// procctl and hands it to onStarted (which records agent.started), then owns
// cancellation: one goroutine Waits, and on ctx cancel the whole group is killed
// (fixing orphan-on-timeout). Shared by the round path and steer attempts.
func execAgentProcess(ctx context.Context, root, runID, agentID, marker string, agent agents.Discovery, prompt, stdoutPath, stderrPath string, onStarted func(procctl.Spawned), act *activityTracker, cfg SupervisionConfig, hooks supervisionHooks) (procctl.Spawned, error) {
	path, args, env, cleanup, err := buildAgentInvocation(root, agent, prompt)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return procctl.Spawned{}, err
	}
	cmd := exec.Command(path, args...)
	if env == nil {
		env = os.Environ()
	}
	cmd.Env = append(cleanParticipantEnv(agent.Adapter(), env), procctl.MarkerEnv(runID, agentID, marker)...)
	cmd.Dir = root
	procctl.SetNewProcessGroup(cmd)

	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return procctl.Spawned{}, err
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return procctl.Spawned{}, err
	}
	defer stderrFile.Close()
	if act == nil {
		act = &activityTracker{}
	}
	// Counting writers attribute output bytes in-process — supervision adds no
	// filesystem probing on the healthy path (consensus D1).
	cmd.Stdout = &countingWriter{w: stdoutFile, t: act, stream: "stdout"}
	cmd.Stderr = &countingWriter{w: stderrFile, t: act, stream: "stderr"}
	if agent.PromptMode == agents.PromptStdin {
		cmd.Stdin = strings.NewReader(prompt)
	}
	if err := cmd.Start(); err != nil {
		return procctl.Spawned{}, err
	}
	sp := procctl.Capture(cmd, marker)
	if onStarted != nil {
		onStarted(sp)
	}
	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()
	kill := func() { _ = procctl.KillGroup(sp) } // reap the whole tree, not just the direct child
	return sp, waitSupervised(ctx.Done(), ctx.Err, waitErr, kill, act, cfg, hooks)
}

// cleanParticipantEnv sheds nested host-session markers when the spawned
// participant is the claude CLI: a participant is independent by definition
// and must not inherit the facilitator's session identity (consensus D8).
// Parley's own PARLEY_* markers are kept.
func cleanParticipantEnv(family string, env []string) []string {
	if family != "claude" {
		return env
	}
	out := env[:0:0]
	for _, kv := range env {
		key, _, _ := strings.Cut(kv, "=")
		switch {
		case key == "CLAUDECODE", key == "AI_AGENT":
			continue
		case strings.HasPrefix(key, "CLAUDE_CODE_"), strings.HasPrefix(key, "AI_AGENT_"):
			continue
		}
		out = append(out, kv)
	}
	return out
}

// buildAgentInvocation resolves the agent's command path, args (with {root}/
// {prompt} substitution), and isolated-home env, without binding a context. Used
// by both CommandFor (ctx-bound, for one-shot helpers) and execAgentProcess.
func buildAgentInvocation(root string, agent agents.Discovery, prompt string) (path string, args, env []string, cleanup func(), err error) {
	args = make([]string, 0, len(agent.HeadlessArgs))
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
	cleanup = func() {}
	if agent.IsolateHome {
		e, remove, err := isolatedAgentHome(agent)
		if err != nil {
			return "", nil, nil, nil, err
		}
		cleanup = remove
		env = append(os.Environ(), e...)
		if agent.Adapter() == "hermes" {
			env = append(env, "HERMES_ACCEPT_HOOKS=1", "HERMES_SESSION_SOURCE=parley")
		}
	}
	return agent.Path, args, env, cleanup, nil
}

func CommandFor(ctx context.Context, root string, agent agents.Discovery, prompt string) (*exec.Cmd, func(), error) {
	path, args, env, cleanup, err := buildAgentInvocation(root, agent, prompt)
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.CommandContext(ctx, path, args...)
	if env != nil {
		cmd.Env = env
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
	switch agent.Adapter() {
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
