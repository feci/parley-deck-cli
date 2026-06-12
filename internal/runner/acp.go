package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"parley-deck-cli/internal/acp"
	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/procctl"
	"parley-deck-cli/internal/store"
)

// acpClientName matches what ACP-aware CLIs expect to log; AionUi uses
// "AionUi" — parley advertises itself as "parley-deck".
const acpClientName = "parley-deck"

// runACPAgent runs a participant via the Agent Client Protocol instead of
// piping a prompt through one-shot text stdio. The agent is expected to write
// the canonical artifact file (outputPath) via its own filesystem tools.
// Streaming session/update notifications are appended to the run's event log.
func runACPAgent(parent context.Context, opts Options, agent agents.Discovery, result Result, outputPath, stdoutPath, stderrPath, prompt string, attemptID int) Result {
	if agent.ACPArgs == nil && agent.Path != "" {
		// Defensive: an ACP-mode agent must declare its launch flags; AionUi
		// defaults to ["--experimental-acp"] for claude when unset. Requiring
		// nil-vs-empty intent keeps behavior predictable while still allowing
		// binaries that speak ACP with no arguments.
		return failEarly(opts, result, fmt.Errorf("agent %s has launch_mode=acp but ACPArgs is empty", agent.ID))
	}

	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return failEarly(opts, result, err)
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return failEarly(opts, result, err)
	}
	defer stderrFile.Close()

	ctx, cancel := context.WithTimeout(parent, timeoutForAgent(opts.Timeout, agent))
	defer cancel()

	// Register so Handle.KillAgent can cancel this ACP attempt's context (the same
	// per-agent kill path as headless agents); deregister on return. The KILLED
	// badge comes from the agent.killed event the kill emits, projected as sticky.
	opts.tracker.register(agent.ID, opts.SegmentID, "round", "", cancel)
	defer opts.tracker.finish(agent.ID)

	env := acp.MergedEnv(ctx, nil)
	process, err := acp.Spawn(ctx, acp.SpawnOptions{
		Command:    agent.Path,
		Args:       append([]string(nil), agent.ACPArgs...),
		WorkingDir: opts.Root,
		Env:        env,
	})
	if err != nil {
		return failEarly(opts, result, err)
	}

	teeStdout, teeOut := newTeeReader(process.Stdout(), stdoutFile)
	transport := acp.NewTransport(teeStdout, &lineWriter{w: process.Stdin(), copy: stdoutFile})
	handler := &acpRunnerHandler{
		store:     opts.Store,
		agentID:   agent.ID,
		ideaSlug:  opts.Idea.Slug,
		round:     opts.RoundLabel,
		artifact:  outputPath,
		stdoutTap: stdoutFile,
	}
	client := acp.NewClient(transport, handler, acp.ClientInfo{Name: acpClientName, Version: "0.1.0"})
	client.Start()

	// Capture the durable process identity so a restarted parley can re-attribute
	// and group-kill this ACP agent (same as the headless path).
	sp := procctl.CaptureByPID(process.PID(), fmt.Sprintf("%s:%s:%d", opts.RunID, agent.ID, attemptID))
	if appendErr := opts.Store.Append(store.Event{
		Time: time.Now().UTC(),
		Type: "agent.started",
		Data: map[string]any{
			"agent":       agent.ID,
			"artifact":    outputPath,
			"stdout":      stdoutPath,
			"stderr":      stderrPath,
			"launch":      agents.LaunchACP,
			"command":     sp.Command,
			"acp_args":    agent.ACPArgs,
			"segment_id":  opts.SegmentID,
			"attempt_id":  attemptID,
			"pid":         sp.PID,
			"pgid":        sp.PGID,
			"boot_id":     sp.BootID,
			"proc_start":  sp.ProcStart,
			"proc_marker": sp.Marker,
		},
	}); appendErr != nil {
		return failEarly(opts, result, fmt.Errorf("event append failed: %w", appendErr))
	}

	// The initialize→session→prompt sequence runs in a goroutine so the same
	// supervised wait (first-output watchdog, stall guard, heartbeats) covers
	// ACP agents (consensus D1). Activity arrives via the handler; the
	// agent.started event itself never satisfies the first-output guard.
	act := &activityTracker{}
	handler.activity = act
	cfg := supervisionForAgent(agent, timeoutForAgent(opts.Timeout, agent))
	hardTimeout := timeoutForAgent(opts.Timeout, agent)
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
					"launch":               agents.LaunchACP,
					"elapsed_ms":           elapsed.Milliseconds(),
					"timeout_ms":           hardTimeout.Milliseconds(),
					"stdout_bytes":         snap.StdoutBytes,
					"stderr_bytes":         snap.StderrBytes,
					"last_activity_ms_ago": activityAgeMS(snap),
				},
			})
		},
		onWatchdog: func(kind string, snap activitySnapshot, elapsed time.Duration) {
			// Appended BEFORE the kill fires (consensus D1).
			action := "failed"
			if kind == "no_first_output" && attemptID == 1 {
				action = "retrying" // ACP shares the exec retry-once contract
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

	acpCtx, acpCancel := context.WithCancel(ctx)
	defer acpCancel()
	finishCh := make(chan error, 1)
	go func() {
		initResult, err := client.Initialize(acpCtx, acp.ClientCapabilities{})
		if err != nil {
			finishCh <- fmt.Errorf("initialize: %w", err)
			return
		}
		act.MarkEvent() // a live initialize is agent activity (review fix 1a)
		_ = opts.Store.Append(store.Event{
			Type: "agent.acp.initialized",
			Data: map[string]any{
				"agent":            agent.ID,
				"protocol_version": initResult.ProtocolVersion,
				"agent_info":       initResult.AgentInfo,
			},
		})

		session, err := client.NewSession(acpCtx, acp.NewSessionParams{CWD: opts.Root})
		if err != nil {
			finishCh <- fmt.Errorf("session/new: %w", err)
			return
		}
		handler.sessionID = session.SessionID
		act.MarkEvent()
		_ = opts.Store.Append(store.Event{
			Type: "agent.acp.session_opened",
			Data: map[string]any{"agent": agent.ID, "session_id": session.SessionID},
		})

		promptText := strings.TrimRight(prompt, "\n")
		promptRes, err := client.Prompt(acpCtx, session.SessionID, promptText)
		if err == nil {
			act.MarkEvent()
			_ = opts.Store.Append(store.Event{
				Type: "agent.acp.prompt_completed",
				Data: map[string]any{"agent": agent.ID, "stop_reason": promptRes.StopReason},
			})
		}
		finishCh <- err
	}()
	kill := func() {
		acpCancel()
		_ = process.Kill()
	}
	finishErr := waitSupervised(ctx.Done(), ctx.Err, finishCh, kill, act, cfg, hooks)
	watchdog := ""
	switch {
	case errors.Is(finishErr, errNoFirstOutput):
		watchdog = "no_first_output"
	case errors.Is(finishErr, errStalled):
		watchdog = "stalled"
	}
	return finishACP(opts, result, agent, process, stderrFile, teeOut, finishErr, handler, watchdog, ctx.Err(), attemptID)
}

// finishACP centralises shutdown, stderr capture, artifact validation and
// the terminal agent.finished/agent.failed event, applying the consensus D7
// decision table: a validated artifact overrides an ACP prompt error that
// happened AFTER the session opened (agent_exit_kind=acp_error); initialize/
// session-setup errors, watchdog kills, and the hard timeout always fail.
func finishACP(opts Options, result Result, agent agents.Discovery, process *acp.Process, stderrFile *os.File, teeOut chan struct{}, runErr error, handler *acpRunnerHandler, watchdog string, ctxErr error, attemptID int) Result {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = process.Stop(shutdownCtx)
	<-teeOut
	if stderrFile != nil {
		_, _ = io.WriteString(stderrFile, process.Stderr())
	}

	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)
	if runErr != nil {
		result.ExitError = runErr.Error()
	}

	if _, statErr := os.Stat(result.OutputPath); statErr == nil {
		// Validate per the run's phase contract (was: hard-coded round-one
		// validation, which broke ACP review/implementation artifacts).
		if validateErr := validateArtifactForPhase(opts, result.OutputPath, agent.ID); validateErr != nil {
			result.ExitError = combineError(result.ExitError, validateErr)
		} else {
			result.ArtifactOK = true
		}
	}

	// Review snapshots (D9): publish a validated artifact to the canonical
	// path before the terminal event reports it.
	if result.ArtifactOK && opts.publishArtifact != nil {
		livePath, err := opts.publishArtifact(result.OutputPath)
		if livePath != "" {
			// Terminal events report the LIVE canonical path even when the
			// move-back failed (the snapshot path travels only as recovery
			// metadata in review.snapshot_artifact_move_failed).
			result.OutputPath = livePath
		}
		if err != nil {
			result.ArtifactOK = false
			result.ExitError = combineError(result.ExitError, fmt.Errorf("snapshot artifact move-back: %w", err))
		}
	}

	// Artifact-wins (D7): only a post-session prompt error qualifies — and
	// never a watchdog kill or the outer timeout/cancel.
	cancellation := watchdog != "" || errors.Is(ctxErr, context.DeadlineExceeded) || errors.Is(ctxErr, context.Canceled)
	if result.ArtifactOK && !cancellation && runErr != nil && handler != nil && handler.sessionID != "" {
		result.AgentExit = 1
		result.AgentExitKind = "acp_error"
		result.Warning = combineWarning(result.Warning, "acp prompt error overridden by validated artifact: "+runErr.Error())
		result.ExitError = ""
	}

	failed := result.ExitError != "" || !result.ArtifactOK
	if failed && result.ExitError == "" {
		result.ExitError = "artifact missing or invalid"
	}
	data := map[string]any{
		"agent":       agent.ID,
		"artifact":    result.OutputPath,
		"artifact_ok": result.ArtifactOK,
		"duration_ms": result.Duration.Milliseconds(),
		"error":       result.ExitError,
		"launch":      agents.LaunchACP,
		"segment_id":  opts.SegmentID,
		"attempt_id":  attemptID,
	}
	eventType := "agent.finished"
	if failed {
		eventType = "agent.failed"
		class, hint := terminalFailureClass(watchdog, ctxErr, false, result.StderrPath, result.StdoutPath, result.ExitError)
		result.FailureClass = class
		result.RecoveryHint = hint
		data["failure_class"] = class
		data["recovery_hint"] = hint
		data["stderr_tail_bytes"] = len(tailOfFile(result.StderrPath, failTailBytes))
	} else if result.AgentExitKind != "" {
		data["agent_exit"] = result.AgentExit
		data["agent_exit_kind"] = result.AgentExitKind
	}
	if err := opts.Store.Append(store.Event{Time: result.CompletedAt, Type: eventType, Data: data}); err != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", err))
	}
	return result
}

func combineWarning(existing, extra string) string {
	if existing == "" {
		return extra
	}
	return existing + "; " + extra
}

// acpRunnerHandler bridges ACP callbacks to the parley event store. Streaming
// chunks are aggregated for the stdout tap; permission requests are
// auto-allowed (matching headless behavior); fs/* requests are refused so
// the agent uses its own filesystem tools.
type acpRunnerHandler struct {
	store     store.Store
	agentID   string
	ideaSlug  string
	round     string
	artifact  string
	stdoutTap *os.File
	sessionID string
	// activity feeds the supervision watchdogs: any ACP session update counts
	// as agent activity (consensus D1); message text additionally counts bytes.
	activity *activityTracker

	mu         sync.Mutex
	messageBuf strings.Builder
	thoughtBuf strings.Builder
	chunkBatch int
	noop       acp.NoopHandler
}

func (h *acpRunnerHandler) SessionUpdate(update acp.SessionUpdate) error {
	if h.activity != nil {
		h.activity.MarkEvent()
		if update.Update.Content != nil {
			h.activity.Mark("stdout", len(update.Update.Content.Text))
		}
	}
	switch update.Update.SessionUpdate {
	case acp.UpdateAgentMessageChunk:
		if update.Update.Content != nil && update.Update.Content.Text != "" {
			h.mu.Lock()
			h.messageBuf.WriteString(update.Update.Content.Text)
			h.chunkBatch++
			batch := h.chunkBatch
			h.mu.Unlock()
			if h.stdoutTap != nil {
				_, _ = io.WriteString(h.stdoutTap, update.Update.Content.Text)
			}
			if batch%32 == 0 {
				h.flushChunkEvent(false)
			}
		}
	case acp.UpdateAgentThoughtChunk:
		if update.Update.Content != nil && update.Update.Content.Text != "" {
			h.mu.Lock()
			h.thoughtBuf.WriteString(update.Update.Content.Text)
			h.mu.Unlock()
		}
	case acp.UpdateToolCall, acp.UpdateToolCallUpdate:
		var raw map[string]any
		_ = json.Unmarshal(update.Update.Raw, &raw)
		_ = h.store.Append(store.Event{
			Type: "agent.acp.tool_call",
			Data: map[string]any{
				"agent":        h.agentID,
				"tool_call_id": update.Update.ToolCallID,
				"status":       update.Update.Status,
				"title":        update.Update.Title,
				"kind":         update.Update.Kind,
				"phase":        update.Update.SessionUpdate,
			},
		})
	case acp.UpdatePlan:
		entries := make([]map[string]any, 0, len(update.Update.Entries))
		for _, e := range update.Update.Entries {
			entries = append(entries, map[string]any{
				"content":  e.Content,
				"status":   e.Status,
				"priority": e.Priority,
			})
		}
		_ = h.store.Append(store.Event{
			Type: "agent.acp.plan",
			Data: map[string]any{"agent": h.agentID, "entries": entries},
		})
	case acp.UpdateUsage:
		_ = h.store.Append(store.Event{
			Type: "agent.acp.usage",
			Data: map[string]any{
				"agent": h.agentID,
				"used":  update.Update.Used,
				"size":  update.Update.Size,
			},
		})
	}
	return nil
}

func (h *acpRunnerHandler) flushChunkEvent(final bool) {
	h.mu.Lock()
	text := h.messageBuf.String()
	h.messageBuf.Reset()
	h.mu.Unlock()
	if text == "" {
		return
	}
	_ = h.store.Append(store.Event{
		Type: "agent.acp.message_chunk",
		Data: map[string]any{"agent": h.agentID, "text": text, "final": final},
	})
}

func (h *acpRunnerHandler) RequestPermission(req acp.PermissionRequest) (acp.PermissionResult, error) {
	result, err := h.noop.RequestPermission(req)
	_ = h.store.Append(store.Event{
		Type: "agent.acp.permission",
		Data: map[string]any{
			"agent":     h.agentID,
			"tool_call": req.ToolCall.ToolCallID,
			"kind":      req.ToolCall.Kind,
			"decision":  result.Outcome.Outcome,
			"option_id": result.Outcome.OptionID,
		},
	})
	return result, err
}

func (h *acpRunnerHandler) ReadTextFile(req acp.FSReadRequest) (acp.FSReadResult, error) {
	return acp.FSReadResult{}, errors.New("acp: client did not advertise fs.readTextFile")
}

func (h *acpRunnerHandler) WriteTextFile(req acp.FSWriteRequest) error {
	return errors.New("acp: client did not advertise fs.writeTextFile")
}

// lineWriter forwards stdin writes to both the agent's stdin and a tee sink
// so the request stream is captured in stdout.log alongside replies.
type lineWriter struct {
	mu   sync.Mutex
	w    io.Writer
	copy io.Writer
}

func (l *lineWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.copy != nil {
		_, _ = l.copy.Write([]byte("--> "))
		_, _ = l.copy.Write(p)
	}
	return l.w.Write(p)
}

// newTeeReader returns a reader that mirrors the source stream into the sink.
// The returned channel closes when the copy completes (EOF or error) so the
// caller can synchronise teardown with the tee goroutine.
func newTeeReader(src io.Reader, sink io.Writer) (io.Reader, chan struct{}) {
	pr, pw := io.Pipe()
	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := src.Read(buf)
			if n > 0 {
				_, _ = sink.Write(buf[:n])
				if _, werr := pw.Write(buf[:n]); werr != nil {
					_ = pw.CloseWithError(werr)
					return
				}
			}
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()
	return pr, done
}
