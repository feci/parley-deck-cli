package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

func TestAgentsListPrintsResolvedRuntime(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	writeFakeCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"agents", "list", "--dir", root}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"codex", "yes", "codex test 1.0", "configured", "workspace-write", "on-failure", "cli-default"} {
		if !strings.Contains(out, want) {
			t.Fatalf("list output missing %q:\n%s", want, out)
		}
	}
}

func TestAgentsCompatibilityAliases(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	writeFakeCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin)

	var stdout, stderr bytes.Buffer
	if code := Run([]string{"agents", "discover", "--dir", root}, &stdout, &stderr); code != 0 {
		t.Fatalf("discover alias code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "codex") {
		t.Fatalf("discover stdout=%q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"agents", "probe", "--dir", root, "--agent", "codex"}, &stdout, &stderr); code != 0 {
		t.Fatalf("probe alias code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}

func TestAgentsVerifyCheapPath(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	writeFakeCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"agents", "verify", "--dir", root, "--agent", "codex"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "codex: installed version=codex test 1.0") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestCodexProbePromptIncludesGitSmoke(t *testing.T) {
	prompt := probePrompt(agents.Discovery{Spec: agents.Spec{ID: "codex"}}, "/tmp/probe.md", "# sentinel")
	for _, want := range []string{
		"git status",
		"git branch tmp-codex-git-test",
		"git branch -D tmp-codex-git-test",
		"printf test | git hash-object -w --stdin",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestRunRecordsResolvedRuntime(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	localConfig := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	if err := os.WriteFile(localConfig, []byte(`
[agents.codex]
model = "local-model"
approval_policy = "on-failure"
`), 0o644); err != nil {
		t.Fatal(err)
	}

	bin := t.TempDir()
	writeFakeRoundAgentCLI(t, bin, "codex", "codex test 1.0")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"run", "--dir", root, "--no-tui", "--yes", "--participants", "codex", "Runtime task"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}

	runsDir := filepath.Join(root, protocol.DeckDir, "runs")
	entries, err := os.ReadDir(runsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("runs=%d, want 1", len(entries))
	}
	events, err := store.New(filepath.Join(runsDir, entries[0].Name())).Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 || events[0].Type != "run.created" {
		t.Fatalf("events=%+v", events)
	}
	runtimeRows, ok := events[0].Data["runtime"].([]any)
	if !ok || len(runtimeRows) == 0 {
		t.Fatalf("runtime event data=%+v", events[0].Data["runtime"])
	}
	row, ok := runtimeRows[0].(map[string]any)
	if !ok {
		t.Fatalf("runtime row=%+v", runtimeRows[0])
	}
	if row["agent"] != "codex" || row["model"] != "local-model" || row["approval_policy"] != "on-failure" {
		t.Fatalf("runtime row=%+v", row)
	}
}

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

func TestStatusAndResumeUseRunState(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "sample")
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte("---\nidea: sample\nparticipants: [codex]\nstatus: round-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runID := "20260512T100000.000000000Z"
	runDir := filepath.Join(root, protocol.DeckDir, "runs", runID)
	base := time.Date(2026, 5, 12, 10, 0, 0, 0, time.UTC)
	s := store.New(runDir)
	for _, event := range []store.Event{
		{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "mode": "hitl", "participants": []string{"codex"}, "task": "Sample task"}},
		{Time: base.Add(time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex", "stdout": "stdout.log", "stderr": "stderr.log"}},
	} {
		if err := s.Append(event); err != nil {
			t.Fatal(err)
		}
	}
	question, err := hitl.New(runDir).Create(hitl.Question{Agent: "codex", Prompt: "Which branch?", Risk: hitl.RiskNormal})
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"status", "--dir", root}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Transport:", "Ideas:", "sample", "Runs:", runID, "questions=1 open"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("status output missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"resume", "--dir", root, "--no-tui", runID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("resume code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Run: " + runID, "Idea: sample", "State: unverified", "Open HITL questions:", question.ID, "Next: parley answer " + runID + " " + question.ID} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("resume output missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"status", "--dir", root, "--idea", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("status --idea code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Run: " + runID, "Idea: sample", "State: unverified"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("status --idea output missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"status", "--dir", root, "--run", runID, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("json status code=%d stderr=%s", code, stderr.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	if payload["run_id"] != runID || payload["idea_slug"] != "sample" {
		t.Fatalf("payload=%+v", payload)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"status", "--dir", root, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("workspace json status code=%d stderr=%s", code, stderr.String())
	}
	var workspacePayload struct {
		Runs []map[string]any `json:"runs"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &workspacePayload); err != nil {
		t.Fatalf("invalid workspace json: %v\n%s", err, stdout.String())
	}
	if len(workspacePayload.Runs) != 1 || workspacePayload.Runs[0]["run_id"] != runID {
		t.Fatalf("workspace payload=%+v", workspacePayload)
	}
}

func TestResumeReportsKnownIdeaWithNoRuns(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "empty")
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte("---\nidea: empty\nparticipants: [codex]\nstatus: round-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"resume", "--dir", root, "--no-tui", "empty"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), `idea "empty" has no runs yet`) {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestAgentDurationUsesElapsedForRunningSnapshot(t *testing.T) {
	duration := agentDuration(runstate.AgentState{
		State:     runstate.StateRunning,
		StartedAt: time.Now().Add(-2 * time.Minute),
	})
	if duration <= 0 {
		t.Fatalf("duration=%s, want elapsed duration", duration)
	}
}

func writeFakeCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	path := filepath.Join(dir, name)
	body := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo '" + version + "'; exit 0; fi\ncat >/dev/null\nexit 0\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFakeRoundAgentCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	path := filepath.Join(dir, name)
	body := `#!/bin/sh
if [ "$1" = "--version" ]; then
  echo '` + version + `'
  exit 0
fi
out=$(awk -F': ' '/Create exactly this file and no other protocol artifact:/ {print $2; exit}')
if [ -z "$out" ]; then
  exit 3
fi
cat > "$out" <<'ARTIFACT'
---
agent: codex
idea: runtime-task
round: 1
date: 2026-05-11
---

## Summary
Fake artifact.

## Proposed approach
Use the test helper.

## Concerns / open questions
None.

## Risks
None.
ARTIFACT
exit 0
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}
