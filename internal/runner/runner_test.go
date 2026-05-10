package runner

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func TestBuildRoundOnePrompt(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Test task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(idea.Path, "round-01", "fake.md")

	prompt, err := BuildRoundOnePrompt("fake", idea, "Test task", output)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"Create exactly this file",
		output,
		"Do not edit any other agent's file.",
		"agent: fake",
		"idea: " + idea.Slug,
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q\n%s", want, prompt)
		}
	}

	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("prompt generation should not create artifact, stat err=%v", err)
	}
}

func TestRunRoundOneCreatesArtifactWithHeadlessAgent(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Runner test task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "test-run"))

	results := RunRoundOne(context.Background(), Options{
		Root:  root,
		RunID: "test-run",
		Idea:  idea,
		Task:  "Runner test task",
		Agents: []agents.Discovery{
			{
				Spec: agents.Spec{
					ID:           "fake",
					HeadlessArgs: []string{"-test.run=TestFakeAgentHelper", "--", "parley-fake-agent"},
					PromptMode:   agents.PromptStdin,
				},
				Path:  os.Args[0],
				Found: true,
			},
		},
		Timeout: 5 * time.Second,
		Store:   runStore,
	})

	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if !results[0].ArtifactOK || results[0].ExitError != "" {
		t.Fatalf("unexpected result: %+v", results[0])
	}
	if _, err := os.Stat(filepath.Join(idea.Path, "round-01", "fake.md")); err != nil {
		t.Fatal(err)
	}

	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := events[len(events)-1].Type, "round.completed"; got != want {
		t.Fatalf("last event=%s, want %s", got, want)
	}
}

func TestRunRoundOneSkipsExistingArtifact(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Skip test task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(idea.Path, "round-01", "fake.md")
	if err := os.WriteFile(output, []byte("already done\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "skip-run"))

	results := RunRoundOne(context.Background(), Options{
		Root:  root,
		RunID: "skip-run",
		Idea:  idea,
		Task:  "Skip test task",
		Agents: []agents.Discovery{
			{Spec: agents.Spec{ID: "fake"}, Path: os.Args[0], Found: true},
		},
		Timeout: 5 * time.Second,
		Store:   runStore,
	})

	if len(results) != 1 || !results[0].Skipped || results[0].ExitError != "" {
		t.Fatalf("unexpected result: %+v", results)
	}
	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := events[0].Type, "agent.skipped"; got != want {
		t.Fatalf("first event=%s, want %s", got, want)
	}
	if got, want := events[len(events)-1].Type, "round.completed"; got != want {
		t.Fatalf("last event=%s, want %s", got, want)
	}
}

func TestRunRoundOneRecordsAgentFailure(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Failure test task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "failure-run"))

	results := RunRoundOne(context.Background(), Options{
		Root:  root,
		RunID: "failure-run",
		Idea:  idea,
		Task:  "Failure test task",
		Agents: []agents.Discovery{
			{
				Spec: agents.Spec{
					ID:           "fake",
					HeadlessArgs: []string{"-test.run=TestFakeFailingAgentHelper", "--", "parley-failing-agent"},
					PromptMode:   agents.PromptStdin,
				},
				Path:  os.Args[0],
				Found: true,
			},
		},
		Timeout: 5 * time.Second,
		Store:   runStore,
	})

	if len(results) != 1 || results[0].ExitError == "" || results[0].ArtifactOK {
		t.Fatalf("unexpected result: %+v", results)
	}
	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := events[len(events)-2].Type, "agent.failed"; got != want {
		t.Fatalf("terminal agent event=%s, want %s", got, want)
	}
	if got, want := events[len(events)-1].Type, "round.incomplete"; got != want {
		t.Fatalf("last event=%s, want %s", got, want)
	}
}

func TestFakeAgentHelper(t *testing.T) {
	if !hasArg("parley-fake-agent") {
		return
	}
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`(?m)^- Create exactly this file and no other protocol artifact: (.+)$`)
	match := re.FindStringSubmatch(string(input))
	if len(match) != 2 {
		t.Fatalf("could not find output path in prompt:\n%s", string(input))
	}
	body := `---
agent: fake
idea: runner-test-task
round: 1
date: 2026-05-10
---

## Summary
Fake artifact for runner test.

## Proposed approach
Use the runner contract.

## Concerns / open questions
None.

## Risks
None.
`
	if err := os.WriteFile(match[1], []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func TestFakeFailingAgentHelper(t *testing.T) {
	if !hasArg("parley-failing-agent") {
		return
	}
	os.Exit(7)
}

func hasArg(want string) bool {
	for _, arg := range os.Args[1:] {
		if arg == want {
			return true
		}
	}
	return false
}
