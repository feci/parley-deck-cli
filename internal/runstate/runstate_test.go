package runstate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestLoadRunDerivesCompletedOutcome(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex"})
	runID := "20260512T100000.000000000Z"
	runDir := RunDir(root, runID)
	base := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
	appendEvents(t, runDir,
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "mode": "hitl", "participants": []string{"codex"}, "task": "Sample task"}},
		store.Event{Time: base.Add(time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex", "artifact": "round-01/codex.md"}},
		store.Event{Time: base.Add(2 * time.Second), Type: "agent.finished", Data: map[string]any{"agent": "codex", "artifact": "round-01/codex.md", "duration_ms": float64(1000)}},
		store.Event{Time: base.Add(3 * time.Second), Type: "round.completed", Data: map[string]any{"completed": float64(1), "total": float64(1)}},
	)

	run, err := LoadRunAt(root, runID, base.Add(10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if !run.Terminal || run.Outcome != OutcomeCompleted || run.Liveness != "" {
		t.Fatalf("terminal=%v outcome=%q liveness=%q", run.Terminal, run.Outcome, run.Liveness)
	}
	if run.IdeaSlug != "sample" || run.Task != "Sample task" || len(run.State.Agents) != 1 {
		t.Fatalf("run=%+v", run)
	}
	if got := run.State.Agents[0].ArtifactPath; got != "round-01/codex.md" {
		t.Fatalf("artifact=%q", got)
	}
}

func TestLoadRunDerivesUnverifiedLivenessAndQuestions(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex"})
	runID := "20260512T100100.000000000Z"
	runDir := RunDir(root, runID)
	base := time.Date(2026, 5, 12, 10, 1, 0, 0, time.UTC)
	appendEvents(t, runDir,
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
		store.Event{Time: base.Add(time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex"}},
	)
	if _, err := hitl.New(runDir).Create(hitl.Question{Agent: "codex", Prompt: "Which branch?", Risk: hitl.RiskNormal}); err != nil {
		t.Fatal(err)
	}

	run, err := LoadRunAt(root, runID, base.Add(11*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if run.Terminal || run.Outcome != "" || run.Liveness != LivenessUnverified {
		t.Fatalf("terminal=%v outcome=%q liveness=%q", run.Terminal, run.Outcome, run.Liveness)
	}
	if run.OpenQuestions != 1 {
		t.Fatalf("open questions=%d", run.OpenQuestions)
	}
}

func TestLoadRunDerivesIncompleteOutcome(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex", "claude"})
	runID := "20260512T100200.000000000Z"
	runDir := RunDir(root, runID)
	base := time.Date(2026, 5, 12, 10, 2, 0, 0, time.UTC)
	appendEvents(t, runDir,
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex", "claude"}}},
		store.Event{Time: base.Add(time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex"}},
		store.Event{Time: base.Add(2 * time.Second), Type: "agent.finished", Data: map[string]any{"agent": "codex"}},
		store.Event{Time: base.Add(3 * time.Second), Type: "round.incomplete", Data: map[string]any{"completed": float64(1), "total": float64(2)}},
	)

	run, err := LoadRunAt(root, runID, base.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if !run.Terminal || run.Outcome != OutcomeIncomplete || run.Liveness != "" {
		t.Fatalf("terminal=%v outcome=%q liveness=%q", run.Terminal, run.Outcome, run.Liveness)
	}
}

func TestLoadRunDerivesFailedOutcome(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex"})
	runID := "20260512T100300.000000000Z"
	runDir := RunDir(root, runID)
	base := time.Date(2026, 5, 12, 10, 3, 0, 0, time.UTC)
	appendEvents(t, runDir,
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
		store.Event{Time: base.Add(time.Second), Type: "run.failed", Data: map[string]any{"error": "store unavailable"}},
	)

	run, err := LoadRunAt(root, runID, base.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if !run.Terminal || run.Outcome != OutcomeFailed || run.Liveness != "" {
		t.Fatalf("terminal=%v outcome=%q liveness=%q", run.Terminal, run.Outcome, run.Liveness)
	}
}

func TestLoadRunDerivesIdleLiveness(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex"})
	runID := "20260512T100400.000000000Z"
	base := time.Date(2026, 5, 12, 10, 4, 0, 0, time.UTC)
	appendEvents(t, RunDir(root, runID),
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
	)

	run, err := LoadRunAt(root, runID, base.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if run.Terminal || run.Outcome != "" || run.Liveness != LivenessIdle {
		t.Fatalf("terminal=%v outcome=%q liveness=%q", run.Terminal, run.Outcome, run.Liveness)
	}
}

func TestResolveRunChoosesNewestRunForIdea(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "sample", []string{"codex"})
	base := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
	appendEvents(t, RunDir(root, "20260512T100000.000000000Z"),
		store.Event{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
	)
	appendEvents(t, RunDir(root, "20260512T100500.000000000Z"),
		store.Event{Time: base.Add(5 * time.Minute), Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
	)

	run, err := ResolveRun(root, "sample")
	if err != nil {
		t.Fatal(err)
	}
	if run.RunID != "20260512T100500.000000000Z" {
		t.Fatalf("run=%s", run.RunID)
	}
}

func TestResolveRunReportsKnownIdeaWithNoRuns(t *testing.T) {
	root := t.TempDir()
	writeIdea(t, root, "empty", []string{"codex"})

	_, err := ResolveRun(root, "empty")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `idea "empty" has no runs yet`) {
		t.Fatalf("error=%v", err)
	}
}

func appendEvents(t *testing.T, runDir string, events ...store.Event) {
	t.Helper()
	s := store.New(runDir)
	for _, event := range events {
		if err := s.Append(event); err != nil {
			t.Fatal(err)
		}
	}
}

func writeIdea(t *testing.T, root, slug string, participants []string) {
	t.Helper()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nidea: " + slug + "\nparticipants: [" + strings.Join(participants, ", ") + "]\nstatus: round-01\n---\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
