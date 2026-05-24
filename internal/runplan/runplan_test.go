package runplan

import (
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
)

func TestPlanPrioritizesOpenQuestions(t *testing.T) {
	root := newWorkspace(t, "sample", []string{"codex"})
	actions := Plan(root, Input{
		RunID:        "run-1",
		IdeaSlug:     "sample",
		Participants: []string{"codex"},
		Questions:    []hitl.Question{{ID: "q1", Agent: "codex", Risk: hitl.RiskHigh, Status: hitl.StatusOpen}},
		Agents:       []AgentState{{ID: "codex", State: "running"}},
	})
	if len(actions) == 0 || actions[0].Kind != KindAnswerQuestion || actions[0].ID != "answer-question.q1" || actions[0].Risk != RiskHigh {
		t.Fatalf("actions=%+v", actions)
	}
}

func TestPlanRetriesFailedMissingRoundArtifact(t *testing.T) {
	root := newWorkspace(t, "sample", []string{"codex", "claude"})
	writeFile(t, filepath.Join(root, protocol.DeckDir, "ideas", "sample", "round-01", "codex.md"), "# codex\n")

	actions := Plan(root, Input{
		RunID:        "run-1",
		IdeaSlug:     "sample",
		Participants: []string{"codex", "claude"},
		Terminal:     true,
		RoundStatus:  "incomplete",
		Agents:       []AgentState{{ID: "codex", State: "finished"}, {ID: "claude", State: "failed", Error: "auth"}},
	})
	if len(actions) != 1 || actions[0].Kind != KindRetryAgent || actions[0].AgentID != "claude" {
		t.Fatalf("actions=%+v", actions)
	}
}

func TestPlanUsesCurrentRoundForMissingArtifacts(t *testing.T) {
	root := newWorkspace(t, "sample", []string{"codex", "claude"})
	writeFile(t, filepath.Join(root, protocol.DeckDir, "ideas", "sample", "round-02", "codex.md"), "# codex\n")

	actions := Plan(root, Input{
		RunID:        "run-1",
		IdeaSlug:     "sample",
		Participants: []string{"codex", "claude"},
		Terminal:     true,
		RoundStatus:  "incomplete",
		CurrentRound: "round-02",
		Agents:       []AgentState{{ID: "codex", State: "finished"}, {ID: "claude", State: "failed", Error: "auth"}},
	})
	if len(actions) != 1 || actions[0].Round != "round-02" || actions[0].ArtifactPath != "round-02/claude.md" {
		t.Fatalf("actions=%+v", actions)
	}
}

func TestPlanDraftsConsensusWhenRoundArtifactsExist(t *testing.T) {
	root := newWorkspace(t, "sample", []string{"codex"})
	writeFile(t, filepath.Join(root, protocol.DeckDir, "ideas", "sample", "round-01", "codex.md"), "# codex\n")

	actions := Plan(root, Input{
		RunID:        "run-1",
		IdeaSlug:     "sample",
		Participants: []string{"codex"},
		Terminal:     true,
		Agents:       []AgentState{{ID: "codex", State: "finished"}},
	})
	if len(actions) != 1 || actions[0].Kind != KindDraftConsensus {
		t.Fatalf("actions=%+v", actions)
	}
}

func TestPlanRequestsSignoffsForPartialConsensus(t *testing.T) {
	root := newWorkspace(t, "sample", []string{"codex", "claude"})
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "sample")
	writeFile(t, filepath.Join(ideaDir, "round-01", "codex.md"), "# codex\n")
	writeFile(t, filepath.Join(ideaDir, "round-01", "claude.md"), "# claude\n")
	writeFile(t, filepath.Join(ideaDir, "consensus.md"), `---
idea: sample
drafted-by: codex
date: 2026-05-23
---

# Consensus

## Signoffs

### Signoff: codex - 2026-05-23
Status: accept
`)

	actions := Plan(root, Input{
		RunID:        "run-1",
		IdeaSlug:     "sample",
		Participants: []string{"codex", "claude"},
		Terminal:     true,
	})
	if len(actions) != 1 || actions[0].Kind != KindRequestSignoffs || actions[0].Summary == "" {
		t.Fatalf("actions=%+v", actions)
	}
}

func TestPlanFinalizesReadyConsensus(t *testing.T) {
	root := newWorkspace(t, "sample", []string{"codex"})
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "sample")
	writeFile(t, filepath.Join(ideaDir, "round-01", "codex.md"), "# codex\n")
	writeFile(t, filepath.Join(ideaDir, "consensus.md"), `---
idea: sample
drafted-by: codex
date: 2026-05-23
---

# Consensus

## Signoffs

### Signoff: codex - 2026-05-23
Status: accept
`)

	actions := Plan(root, Input{
		RunID:        "run-1",
		IdeaSlug:     "sample",
		Participants: []string{"codex"},
		Terminal:     true,
	})
	if len(actions) != 1 || actions[0].Kind != KindFinalize {
		t.Fatalf("actions=%+v", actions)
	}
}

func newWorkspace(t *testing.T, slug string, participants []string) string {
	t.Helper()
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if err := os.MkdirAll(filepath.Join(ideaDir, "round-01"), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nidea: " + slug + "\nparticipants: [" + join(participants) + "]\nstatus: round-01\n---\n"
	writeFile(t, filepath.Join(ideaDir, "00-prompt.md"), content)
	return root
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func join(values []string) string {
	out := ""
	for i, value := range values {
		if i > 0 {
			out += ", "
		}
		out += value
	}
	return out
}
