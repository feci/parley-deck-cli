package app

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestRunAnswerUpdatesQuestionAndEventLog(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	runID := "answer-run"
	runDir := filepath.Join(root, protocol.DeckDir, "runs", runID)
	question, err := hitl.New(runDir).Create(hitl.Question{
		Agent:  "codex",
		Prompt: "Which branch?",
		Risk:   hitl.RiskNormal,
	})
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"answer", "--dir", root, runID, question.ID, "main branch"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Answered "+question.ID) {
		t.Fatalf("stdout=%q", stdout.String())
	}

	questions, err := hitl.New(runDir).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) != 1 || questions[0].Status != hitl.StatusAnswered || questions[0].Answer != "main branch" {
		t.Fatalf("questions=%+v", questions)
	}
	events, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := events[len(events)-1].Type; got != "hitl.answered" {
		t.Fatalf("last event=%s, want hitl.answered", got)
	}
}
