package hitl

import (
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/store"
)

func TestCreateListAnswerQuestion(t *testing.T) {
	runDir := t.TempDir()
	hitlStore := New(runDir)
	question, err := hitlStore.Create(Question{
		Agent:         "codex",
		Prompt:        "Which branch?",
		Details:       "Need target branch.",
		DefaultAnswer: "main",
		Risk:          RiskNormal,
	})
	if err != nil {
		t.Fatal(err)
	}
	if question.ID == "" || question.Status != StatusOpen {
		t.Fatalf("unexpected created question: %+v", question)
	}

	questions, err := hitlStore.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) != 1 || questions[0].Prompt != "Which branch?" {
		t.Fatalf("questions=%+v", questions)
	}

	answered, err := hitlStore.Answer(question.ID, "develop", false)
	if err != nil {
		t.Fatal(err)
	}
	if answered.Status != StatusAnswered || answered.Answer != "develop" || answered.AnsweredAt.IsZero() {
		t.Fatalf("unexpected answer: %+v", answered)
	}
	events, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := []string{events[0].Type, events[1].Type}, []string{"hitl.question", "hitl.answered"}; got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("events=%v, want %v", got, want)
	}
}

func TestAutoAnswerOnlyLowRiskWithDefault(t *testing.T) {
	hitlStore := New(t.TempDir())
	low, err := hitlStore.Create(Question{
		Agent:         "runner",
		Prompt:        "Use default formatter?",
		DefaultAnswer: "yes",
		Risk:          RiskLow,
	})
	if err != nil {
		t.Fatal(err)
	}
	high, err := hitlStore.Create(Question{
		Agent:         "runner",
		Prompt:        "Delete files?",
		DefaultAnswer: "no",
		Risk:          RiskHigh,
	})
	if err != nil {
		t.Fatal(err)
	}

	answered, err := hitlStore.AutoAnswerOpen()
	if err != nil {
		t.Fatal(err)
	}
	if len(answered) != 1 || answered[0].ID != low.ID || answered[0].Status != StatusAutoAnswered {
		t.Fatalf("answered=%+v, want only low-risk auto answer", answered)
	}
	questions, err := hitlStore.List()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Question{}
	for _, question := range questions {
		byID[question.ID] = question
	}
	if byID[high.ID].Status != StatusOpen {
		t.Fatalf("high-risk question status=%s, want open", byID[high.ID].Status)
	}
}

func TestListOrdersByCreatedAtThenID(t *testing.T) {
	hitlStore := New(t.TempDir())
	base := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	for _, question := range []Question{
		{ID: "later", Agent: "codex", Prompt: "Later?", CreatedAt: base.Add(time.Second)},
		{ID: "b", Agent: "codex", Prompt: "B?", CreatedAt: base},
		{ID: "a", Agent: "codex", Prompt: "A?", CreatedAt: base},
	} {
		if _, err := hitlStore.Create(question); err != nil {
			t.Fatal(err)
		}
	}

	questions, err := hitlStore.List()
	if err != nil {
		t.Fatal(err)
	}
	got := []string{questions[0].ID, questions[1].ID, questions[2].ID}
	want := []string{"a", "b", "later"}
	if got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("order=%v, want %v", got, want)
	}
}

func TestQuestionIDIncludesAgentSlug(t *testing.T) {
	id := NewQuestionID("Antigravity Agent", time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC))
	if !strings.Contains(id, "antigravity-agent") {
		t.Fatalf("id=%q missing agent slug", id)
	}
}

func TestQuestionIDFallsBackForEmptyAgentSlug(t *testing.T) {
	id := NewQuestionID("@#!", time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC))
	if !strings.Contains(id, "-agent-") {
		t.Fatalf("id=%q missing fallback agent slug", id)
	}
}

func TestRejectsUnsafeQuestionID(t *testing.T) {
	hitlStore := New(t.TempDir())
	if _, err := hitlStore.Create(Question{
		ID:     "../outside",
		Agent:  "codex",
		Prompt: "Unsafe?",
	}); err == nil {
		t.Fatal("expected unsafe id to fail")
	}
	if _, err := hitlStore.Answer("../outside", "no", false); err == nil {
		t.Fatal("expected unsafe id answer to fail")
	}
}
