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
	"parley-deck-cli/internal/store"
)

// acpClientName matches what ACP-aware CLIs expect to log; AionUi uses
// "AionUi" — parley advertises itself as "parley-deck".
const acpClientName = "parley-deck"

// runACPAgent runs a participant via the Agent Client Protocol instead of
// piping a prompt through one-shot text stdio. The agent is expected to write
// the canonical artifact file (outputPath) via its own filesystem tools.
// Streaming session/update notifications are appended to the run's event log.
func runACPAgent(parent context.Context, opts Options, agent agents.Discovery, result Result, outputPath, stdoutPath, stderrPath, prompt string) Result {
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

	if appendErr := opts.Store.Append(store.Event{
		Time: result.StartedAt,
		Type: "agent.started",
		Data: map[string]any{
			"agent":      agent.ID,
			"artifact":   outputPath,
			"stdout":     stdoutPath,
			"stderr":     stderrPath,
			"launch":     agents.LaunchACP,
			"command":    agent.Path,
			"acp_args":   agent.ACPArgs,
			"segment_id": opts.SegmentID,
		},
	}); appendErr != nil {
		return failEarly(opts, result, fmt.Errorf("event append failed: %w", appendErr))
	}

	initResult, err := client.Initialize(ctx, acp.ClientCapabilities{})
	if err != nil {
		return finishACP(opts, result, agent, process, stderrFile, teeOut, fmt.Errorf("initialize: %w", err))
	}
	_ = opts.Store.Append(store.Event{
		Type: "agent.acp.initialized",
		Data: map[string]any{
			"agent":            agent.ID,
			"protocol_version": initResult.ProtocolVersion,
			"agent_info":       initResult.AgentInfo,
		},
	})

	session, err := client.NewSession(ctx, acp.NewSessionParams{CWD: opts.Root})
	if err != nil {
		return finishACP(opts, result, agent, process, stderrFile, teeOut, fmt.Errorf("session/new: %w", err))
	}
	handler.sessionID = session.SessionID
	_ = opts.Store.Append(store.Event{
		Type: "agent.acp.session_opened",
		Data: map[string]any{"agent": agent.ID, "session_id": session.SessionID},
	})

	promptText := strings.TrimRight(prompt, "\n")
	promptRes, err := client.Prompt(ctx, session.SessionID, promptText)
	finishErr := err
	if finishErr == nil {
		_ = opts.Store.Append(store.Event{
			Type: "agent.acp.prompt_completed",
			Data: map[string]any{"agent": agent.ID, "stop_reason": promptRes.StopReason},
		})
	}
	return finishACP(opts, result, agent, process, stderrFile, teeOut, finishErr)
}

// finishACP centralises shutdown, stderr capture, artifact validation and
// the terminal agent.finished/agent.failed event.
func finishACP(opts Options, result Result, agent agents.Discovery, process *acp.Process, stderrFile *os.File, teeOut chan struct{}, runErr error) Result {
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
		if validateErr := ValidateRoundOneArtifact(result.OutputPath, agent.ID, opts.Idea.Slug); validateErr != nil {
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
			"artifact":    result.OutputPath,
			"artifact_ok": result.ArtifactOK,
			"duration_ms": result.Duration.Milliseconds(),
			"error":       result.ExitError,
			"launch":      agents.LaunchACP,
			"segment_id":  opts.SegmentID,
		},
	}); err != nil {
		result.ExitError = combineError(result.ExitError, fmt.Errorf("event append failed: %w", err))
	}
	return result
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

	mu         sync.Mutex
	messageBuf strings.Builder
	thoughtBuf strings.Builder
	chunkBatch int
	noop       acp.NoopHandler
}

func (h *acpRunnerHandler) SessionUpdate(update acp.SessionUpdate) error {
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
