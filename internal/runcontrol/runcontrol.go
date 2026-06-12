package runcontrol

import (
	"context"
	"path/filepath"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runmanifest"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/sessionstore"
	"parley-deck-cli/internal/store"
)

type CreateOptions struct {
	Root         string
	Task         string
	Participants []string
	Discovered   []agents.Discovery
	Auto         bool
	Now          time.Time
}

type CreatedRun struct {
	Idea       protocol.IdeaStatus
	RunID      string
	RunDir     string
	Store      store.Store
	RunOptions runner.Options
}

// writeManifest is a seam so tests can exercise the best-effort manifest path.
var writeManifest = runmanifest.Write

func Create(opts CreateOptions) (CreatedRun, error) {
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	idea, err := protocol.CreateIdea(opts.Root, opts.Task, opts.Participants)
	if err != nil {
		return CreatedRun{}, err
	}
	transport := ""
	if status, err := protocol.ReadWorkspaceStatus(opts.Root); err == nil {
		transport = status.Transport
	}

	mode := ModeName(opts.Auto)
	runID := store.NewRunID(now)
	runDir := filepath.Join(opts.Root, protocol.DeckDir, "runs", runID)
	runStore := store.New(runDir)
	if err := runStore.Append(store.Event{
		Time: now.UTC(),
		Type: "run.created",
		Data: map[string]any{
			"task":         opts.Task,
			"mode":         mode,
			"idea":         idea.Slug,
			"participants": opts.Participants,
			"runtime":      RuntimeEventData(opts.Discovered),
		},
	}); err != nil {
		return CreatedRun{}, err
	}
	if err := writeManifest(opts.Root, runID, runmanifest.New(runmanifest.Options{
		Root:         opts.Root,
		RunID:        runID,
		IdeaSlug:     idea.Slug,
		Task:         opts.Task,
		Mode:         mode,
		Transport:    transport,
		Phase:        "round",
		IdeaStatus:   idea.Status,
		CurrentRound: "round-01",
		ActiveSteps:  initialRoundSteps(opts.Participants),
		Participants: opts.Participants,
		CreatedAt:    now,
		UpdatedAt:    now,
	})); err != nil {
		// Best-effort: a transient manifest-write failure on a weakly-coherent mount
		// (e.g. virtio-fs) must NOT orphan an already-created run. The run is defined by
		// events.jsonl; runstate degrades gracefully when run.json is absent. Record the
		// deferral for the audit trail and let the run proceed.
		_ = runStore.Append(store.Event{
			Time: now.UTC(),
			Type: "run.manifest_deferred",
			Data: map[string]any{"error": err.Error()},
		})
	}

	registerSession(opts.Root, idea, runID, opts.Task, opts.Participants, now)

	created := CreatedRun{
		Idea:   idea,
		RunID:  runID,
		RunDir: runDir,
		Store:  runStore,
	}
	created.RunOptions = runner.Options{
		Root:   opts.Root,
		RunID:  runID,
		Idea:   idea,
		Task:   opts.Task,
		Agents: opts.Discovered,
		Store:  runStore,
	}
	return created, nil
}

func initialRoundSteps(participants []string) []runmanifest.Step {
	steps := make([]runmanifest.Step, 0, len(participants))
	for _, participant := range participants {
		steps = append(steps, runmanifest.Step{
			ID:           "round-01." + participant,
			Kind:         "round",
			AgentID:      participant,
			ArtifactPath: filepath.ToSlash(filepath.Join("round-01", participant+".md")),
			Status:       "pending",
		})
	}
	return steps
}

func StartAutoAnswerer(ctx context.Context, runDir string) {
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = hitl.New(runDir).AutoAnswerOpen()
			}
		}
	}()
}

func ModeName(auto bool) string {
	if auto {
		return "auto"
	}
	return "hitl"
}

func RuntimeEventData(discovered []agents.Discovery) []map[string]any {
	data := make([]map[string]any, 0, len(discovered))
	for _, result := range discovered {
		data = append(data, map[string]any{
			"agent":                  result.ID,
			"installed":              result.Found,
			"version":                result.Version,
			"launch_mode":            agents.LaunchModeOrDefault(result.LaunchMode),
			"acp_args":               result.ACPArgs,
			"sandbox_mode":           result.SandboxMode,
			"approval_policy":        result.ApprovalPolicy,
			"model":                  result.Model,
			"reasoning":              result.Reasoning,
			"profile":                result.Profile,
			"speed":                  result.Speed,
			"timeout_ms":             result.TimeoutMS,
			"first_event_timeout_ms": result.FirstEventTimeoutMS,
			"stall_timeout_ms":       result.StallTimeoutMS,
			"heartbeat_ms":           result.HeartbeatMS,
			"isolate_home":           result.IsolateHome,
			"buffers_stdout":         result.BuffersStdout,
			"external_backend":       result.ExternalBackend,
			"sources":                result.Sources,
		})
	}
	return data
}

func registerSession(root string, idea protocol.IdeaStatus, runID, task string, participants []string, now time.Time) {
	sessionStore, err := sessionstore.Default()
	if err != nil {
		return
	}
	_ = sessionStore.Upsert(sessionstore.Session{
		WorkspaceRoot: root,
		RunID:         runID,
		IdeaSlug:      idea.Slug,
		Task:          task,
		Participants:  append([]string(nil), participants...),
		CreatedAt:     now.UTC(),
		UpdatedAt:     now.UTC(),
		LastEventAt:   now.UTC(),
	})
}
