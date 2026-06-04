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
	questionsDir := filepath.Join(root, protocol.DeckDir, "runs", "test-run", "questions")

	prompt, err := BuildRoundOnePrompt(agents.Discovery{
		Spec: agents.Spec{
			ID:             "fake",
			Model:          "test-model",
			Reasoning:      "test-reasoning",
			Speed:          "test-speed",
			SandboxMode:    "workspace-write",
			ApprovalPolicy: "on-failure",
			TimeoutMS:      1234,
		},
	}, idea, "Test task", output, questionsDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"Create exactly this file",
		output,
		"Do not edit any other agent's file.",
		"The first line of the file must be exactly \"---\".",
		"agent: fake",
		"idea: " + idea.Slug,
		questionsDir,
		"status\":\"open",
		"model: test-model",
		"thinking/reasoning/effort/profile: test-reasoning",
		"speed: test-speed",
		"sandbox: workspace-write",
		"approval: on-failure",
		"timeoutMs: 1234",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q\n%s", want, prompt)
		}
	}

	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("prompt generation should not create artifact, stat err=%v", err)
	}
}

func TestValidateRoundOneArtifactAcceptsMissingOpeningFence(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hermes.md")
	data := `agent: hermes
idea: sample
round: 1
date: 2026-05-17
---

## Summary
Summary.

## Proposed approach
Approach.

## Concerns / open questions
None.

## Risks
None.
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRoundOneArtifact(path, "hermes", "sample"); err != nil {
		t.Fatalf("validation should tolerate missing opening fence: %v", err)
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
	indexPath := filepath.Join(idea.Path, "round-01", "_index.md")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(indexData), "| fake | ok |") {
		t.Fatalf("index missing fake ok row:\n%s", string(indexData))
	}

	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := events[len(events)-1].Type, "round.completed"; got != want {
		t.Fatalf("last event=%s, want %s", got, want)
	}
	var indexWritten bool
	for _, event := range events {
		if event.Type == "round.index_written" {
			indexWritten = true
		}
	}
	if !indexWritten {
		t.Fatalf("events missing round.index_written: %+v", events)
	}
}

func TestRunRoundOneAsyncClosesWithResults(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Async runner test task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "async-run"))

	handle := RunRoundOneAsync(context.Background(), Options{
		Root:  root,
		RunID: "async-run",
		Idea:  idea,
		Task:  "Async runner test task",
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

	select {
	case <-handle.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("async runner did not finish")
	}
	results := handle.Results()
	if len(results) != 1 || !results[0].ArtifactOK {
		t.Fatalf("unexpected async results: %+v", results)
	}
	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	var started, finished bool
	for _, event := range events {
		if event.Type == "agent.started" {
			started = true
		}
		if event.Type == "agent.finished" {
			finished = true
		}
	}
	if !started || !finished {
		t.Fatalf("events missing started/finished: %+v", events)
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
	// A run.segment_started boundary now leads every round-run (Slice 1) so the
	// projection can reset the targeted agents; the skip event follows it.
	if got, want := events[0].Type, "run.segment_started"; got != want {
		t.Fatalf("first event=%s, want %s", got, want)
	}
	if got, want := events[1].Type, "agent.skipped"; got != want {
		t.Fatalf("second event=%s, want %s", got, want)
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
	var agentFailed bool
	for _, event := range events {
		if event.Type == "agent.failed" {
			agentFailed = true
		}
	}
	if !agentFailed {
		t.Fatalf("events missing agent.failed: %+v", events)
	}
	if got, want := events[len(events)-1].Type, "round.incomplete"; got != want {
		t.Fatalf("last event=%s, want %s", got, want)
	}
	indexData, err := os.ReadFile(filepath.Join(idea.Path, "round-01", "_index.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(indexData), "| fake | failed |") {
		t.Fatalf("index missing failed row:\n%s", string(indexData))
	}
}

func TestRunRoundOneIndexWriteFailureIsWarning(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Index warning task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(idea.Path, "round-01", "_index.md")
	if err := os.Mkdir(indexPath, 0o755); err != nil {
		t.Fatal(err)
	}
	runStore := store.New(filepath.Join(root, protocol.DeckDir, "runs", "index-warning-run"))

	results := RunRoundOne(context.Background(), Options{
		Root:  root,
		RunID: "index-warning-run",
		Idea:  idea,
		Task:  "Index warning task",
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

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2: %+v", len(results), results)
	}
	if !results[0].ArtifactOK || results[0].ExitError != "" {
		t.Fatalf("participant should still succeed: %+v", results[0])
	}
	if results[1].AgentID != "runner/index" || results[1].Warning == "" || results[1].ExitError != "" {
		t.Fatalf("index failure should be warning-only: %+v", results[1])
	}
	events, err := runStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	var indexFailed bool
	for _, event := range events {
		if event.Type == "round.index_failed" {
			indexFailed = true
		}
	}
	if !indexFailed {
		t.Fatalf("events missing round.index_failed: %+v", events)
	}
	if got, want := events[len(events)-1].Type, "round.completed"; got != want {
		t.Fatalf("last event=%s, want %s", got, want)
	}
}

func TestIsolatedHomeEnvUsesConfiguredTemplate(t *testing.T) {
	env := isolatedHomeEnv(agents.Discovery{
		Spec: agents.Spec{
			ID: "gemini",
			IsolatedHomeEnv: map[string]string{
				"GEMINI_CLI_HOME": "{tempdir}/gemini",
				"EXTRA_CACHE":     "{tempdir}/cache",
			},
		},
	}, "/tmp/parley-home", "GEMINI_CLI_HOME")

	got := strings.Join(env, "\n")
	for _, want := range []string{
		"GEMINI_CLI_HOME=/tmp/parley-home/gemini",
		"EXTRA_CACHE=/tmp/parley-home/cache",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("env missing %q: %v", want, env)
		}
	}
}

func TestIsolatedHomeEnvFallsBackToHistoricalKey(t *testing.T) {
	env := isolatedHomeEnv(agents.Discovery{Spec: agents.Spec{ID: "gemini"}}, "/tmp/parley-home", "GEMINI_CLI_HOME")
	if len(env) != 1 || env[0] != "GEMINI_CLI_HOME=/tmp/parley-home" {
		t.Fatalf("env=%v", env)
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
	ideaRe := regexp.MustCompile(`(?m)^idea: (.+)$`)
	ideaMatch := ideaRe.FindStringSubmatch(string(input))
	if len(ideaMatch) != 2 {
		t.Fatalf("could not find idea slug in prompt:\n%s", string(input))
	}
	body := `---
agent: fake
idea: ` + ideaMatch[1] + `
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
