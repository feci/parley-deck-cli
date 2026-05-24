package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runcontrol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/runplan"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/store"
)

func TestVersionCommandPrintsSemanticVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if got, want := stdout.String(), versionLine()+"\n"; got != want {
		t.Fatalf("version output=%q want %q", got, want)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"--version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if got, want := stdout.String(), versionLine()+"\n"; got != want {
		t.Fatalf("--version output=%q want %q", got, want)
	}
}

func TestHelpIncludesDescriptionsFlagsAndExamples(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{
		"Parley Deck keeps the cooperation trail",
		"Commands:",
		"Parameters and flags:",
		"--participants AGENTS",
		"--format markdown|json",
		"Examples:",
		"parley context repo-map --dir . --format markdown --max-files 50",
		"Exit codes:",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q:\n%s", want, out)
		}
	}
}

func TestVersionAllJSONIncludesSkillStatus(t *testing.T) {
	bin := t.TempDir()
	writeFakeParleyDeckSkill(t, bin)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"version", "--all", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	if !strings.Contains(stdout.String(), "\n  \"parley\"") {
		t.Fatalf("version json is not indented:\n%s", stdout.String())
	}
	if payload["ok"] != true {
		t.Fatalf("payload=%+v", payload)
	}
	parley := payload["parley"].(map[string]any)
	if parley["version"] != version {
		t.Fatalf("parley=%+v", parley)
	}
	skill := payload["parley_deck_skill"].(map[string]any)
	installer := skill["installer"].(map[string]any)
	if installer["version"] != "1.1.0" {
		t.Fatalf("installer=%+v", installer)
	}
}

func TestVersionAllUsesDirFlagForProjectStatus(t *testing.T) {
	bin := t.TempDir()
	writeFakeParleyDeckSkill(t, bin)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	root := t.TempDir()
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"version", "--all", "--json", "--dir", root}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	skill := payload["parley_deck_skill"].(map[string]any)
	project := skill["project"].(map[string]any)
	if project["projectArg"] != absRoot {
		t.Fatalf("project arg=%v want %s", project["projectArg"], absRoot)
	}
}

func TestVersionAllFallsBackToLegacySkillVersion(t *testing.T) {
	bin := t.TempDir()
	writeFakeLegacyParleyDeckSkill(t, bin)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	var stdout, stderr bytes.Buffer
	code := Run([]string{"version", "--all", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	skill := payload["parley_deck_skill"].(map[string]any)
	if skill["statusSupported"] != false {
		t.Fatalf("skill=%+v", skill)
	}
	installer := skill["installer"].(map[string]any)
	if installer["version"] != "1.0.8" {
		t.Fatalf("installer=%+v", installer)
	}
}

func TestVersionAllMissingSkillErrorIsNotDuplicated(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"version", "--all", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	message, _ := payload["parley_deck_skill_error"].(string)
	if strings.Contains(message, "version probe failed") {
		t.Fatalf("duplicated missing-command error: %q", message)
	}
}

func TestVersionFileMatchesBinaryVersion(t *testing.T) {
	data, err := os.ReadFile("../../VERSION")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(data)); got != version {
		t.Fatalf("VERSION=%q internal version=%q", got, version)
	}
	semver := regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	if !semver.MatchString(version) {
		t.Fatalf("version %q is not a major.minor.patch semantic version", version)
	}
}

func TestContextRepoMapJSON(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "sample"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cmd", "sample", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"context", "repo-map", "--dir", root, "--format", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	if payload["schema_version"] != float64(1) || payload["root"] != "." {
		t.Fatalf("payload=%+v", payload)
	}
	if !strings.Contains(stdout.String(), "cmd/sample/main.go") || !strings.Contains(stdout.String(), "\"name\": \"main\"") {
		t.Fatalf("repo map json missing expected file/symbol:\n%s", stdout.String())
	}
}

func TestContextRepoMapMarkdownAndValidation(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# readme\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"context", "repo-map", "--dir", root, "--format", "markdown", "--max-files", "1"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "# Repository Map") || !strings.Contains(stdout.String(), "README.md") {
		t.Fatalf("markdown output missing expected content:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"context", "repo-map", "--dir", root, "--format", "xml"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "invalid format") {
		t.Fatalf("stderr missing invalid format message: %s", stderr.String())
	}
}

func TestContextUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"context"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "usage: parley context repo-map") {
		t.Fatalf("stderr missing usage: %s", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"context", "bogus"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "usage: parley context repo-map") {
		t.Fatalf("stderr missing usage for bogus subcommand: %s", stderr.String())
	}
}

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
	parleyHome := t.TempDir()
	t.Setenv("PARLEY_HOME", parleyHome)
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
	if _, err := os.Stat(filepath.Join(parleyHome, "sessions.json")); err != nil {
		t.Fatalf("session registry was not written: %v", err)
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
	code = Run([]string{"continue", "--dir", root, runID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("continue code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Run: " + runID, "Recommended: Answer HITL question " + question.ID, "Command: parley answer " + runID + " " + question.ID, "Next actions:", "kind=answer-question"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("continue output missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"continue", "--dir", root, "--json", runID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("continue json code=%d stderr=%s", code, stderr.String())
	}
	var continuePayload struct {
		Actions []map[string]any `json:"actions"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &continuePayload); err != nil {
		t.Fatalf("invalid continue json: %v\n%s", err, stdout.String())
	}
	if len(continuePayload.Actions) == 0 || continuePayload.Actions[0]["kind"] != "answer-question" {
		t.Fatalf("continue payload=%+v", continuePayload)
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

func TestSessionsCLIListAndInspect(t *testing.T) {
	root := t.TempDir()
	parleyHome := t.TempDir()
	t.Setenv("PARLEY_HOME", parleyHome)
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"sessions", "list"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("empty list code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Sessions: none") {
		t.Fatalf("empty sessions list output:\n%s", stdout.String())
	}

	now := time.Date(2026, 5, 18, 9, 30, 0, 0, time.UTC)
	created, err := runcontrol.Create(runcontrol.CreateOptions{
		Root:         root,
		Task:         "Recoverable session task",
		Participants: []string{"codex"},
		Discovered: []agents.Discovery{{
			Spec:  agents.Spec{ID: "codex"},
			Found: true,
		}},
		Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"sessions", "list"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("list code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Session index:", created.RunID, "idea=" + created.Idea.Slug, "workspace=" + root, "status=running", "participants=codex"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("sessions list missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"sessions", "list", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("json list code=%d stderr=%s", code, stderr.String())
	}
	var listPayload struct {
		Sessions []map[string]any `json:"sessions"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &listPayload); err != nil {
		t.Fatalf("invalid json list: %v\n%s", err, stdout.String())
	}
	if len(listPayload.Sessions) != 1 || listPayload.Sessions[0]["run_id"] != created.RunID {
		t.Fatalf("list payload=%+v", listPayload)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"sessions", "inspect", created.RunID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("inspect code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Session: " + created.RunID, "Manifest:", "Manifest schema: 1", "Manifest status: running", "Mode: hitl", "Run: " + created.RunID, "Idea: " + created.Idea.Slug} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("sessions inspect missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"sessions", "inspect", "--json", created.RunID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("json inspect code=%d stderr=%s", code, stderr.String())
	}
	var payload struct {
		Session struct {
			RunID    string `json:"run_id"`
			IdeaSlug string `json:"idea_slug"`
		} `json:"session"`
		Manifest struct {
			RunID  string `json:"run_id"`
			Mode   string `json:"mode"`
			Status string `json:"status"`
		} `json:"manifest"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, stdout.String())
	}
	if payload.Session.RunID != created.RunID || payload.Session.IdeaSlug != created.Idea.Slug || payload.Manifest.Mode != "hitl" || payload.Manifest.Status != "running" {
		t.Fatalf("payload=%+v", payload)
	}

	legacyRoot := t.TempDir()
	if err := protocol.InitWorkspace(legacyRoot); err != nil {
		t.Fatal(err)
	}
	legacyRunID := "20260518T120000.000000000Z"
	legacyRunDir := filepath.Join(legacyRoot, protocol.DeckDir, "runs", legacyRunID)
	if err := store.New(legacyRunDir).Append(store.Event{
		Time: now,
		Type: "run.created",
		Data: map[string]any{
			"idea":         "legacy",
			"mode":         "hitl",
			"participants": []string{"codex"},
			"task":         "Legacy task",
		},
	}); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"sessions", "inspect", "--dir", legacyRoot, legacyRunID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("legacy inspect code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Session: " + legacyRunID, "Manifest: missing (legacy run;", "Run: " + legacyRunID, "Idea: legacy"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("legacy inspect missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"sessions", "inspect", "--dir", legacyRoot, "missing-run"}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("missing run unexpectedly succeeded:\n%s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "sessions inspect failed") {
		t.Fatalf("missing run stderr=%q", stderr.String())
	}
}

func TestConsensusCLIWorkflowAndIdeaStatus(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "sample")
	roundDir := filepath.Join(ideaDir, "round-01")
	if err := os.MkdirAll(roundDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte("---\nidea: sample\nparticipants: [codex]\nstatus: round-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(roundDir, "codex.md"), []byte("# codex\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "draft", "--dir", root, "--by", "codex", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("draft code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Consensus: partial") {
		t.Fatalf("draft stdout=%q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"consensus", "signoff", "--dir", root, "--agent", "codex", "--status", "accept", "--notes", "Accept.", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("signoff code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Consensus: ready") {
		t.Fatalf("signoff stdout=%q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"consensus", "status", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("status code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Consensus: ready") {
		t.Fatalf("consensus status stdout=%q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"status", "--dir", root, "--idea", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("status --idea code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{"Idea: sample", "Status: consensus", "Consensus: ready"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("status --idea output missing %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"consensus", "finalize", "--dir", root, "--by", "codex", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("finalize code=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(ideaDir, "FINAL.md")); err != nil {
		t.Fatal(err)
	}
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		t.Fatal(err)
	}
	if meta["status"] != "final" {
		t.Fatalf("status=%q, want final", meta["status"])
	}
}

func TestConsensusRequestSignoffsHappyPath(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha", "beta"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	beta := writeFakeSignoffCLI(t, bin, "beta", "accept", 0)
	writeAgentsLocalConfig(t, root,
		fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal},
		fakeAgentConfig{ID: "beta", Path: beta, Backend: agents.ExternalLocal},
	)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"Requesting signoff from alpha", "Requesting signoff from beta", "Requested signoffs complete: alpha,beta"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	summary, err := consensus.Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != consensus.TriageReady || len(summary.Signoffs) != 2 {
		t.Fatalf("summary=%+v", summary)
	}
}

func TestConsensusRequestSignoffsDryRunAndHostedGate(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalHosted})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "--dry-run", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("dry-run code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"Consensus signoff request dry-run", "Requires --yes: yes", "alpha backend=hosted"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("dry-run stdout missing %q:\n%s", want, stdout.String())
		}
	}
	summary, err := consensus.Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(summary.Signoffs) != 0 {
		t.Fatalf("dry-run wrote signoffs: %+v", summary.Signoffs)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("hosted gate code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "rerun with --yes or --dry-run") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestConsensusRequestSignoffsManualModeWritesHandoff(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{
		ID:         "alpha",
		Path:       alpha,
		Backend:    agents.ExternalHosted,
		LaunchMode: agents.LaunchManual,
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"Requesting signoff from alpha (manual)", "Manual handoff for alpha", "Requested signoffs pending: alpha"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	summary, err := consensus.Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(summary.Signoffs) != 0 {
		t.Fatalf("manual mode should not invoke headless signer: %+v", summary.Signoffs)
	}
	runsDir := filepath.Join(root, protocol.DeckDir, "runs")
	entries, err := os.ReadDir(runsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("runs=%d, want 1", len(entries))
	}
	handoff := filepath.Join(runsDir, entries[0].Name(), "agents", "alpha", "handoff.md")
	data, err := os.ReadFile(handoff)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Interactive handoff: alpha", runner.UsageCaveat, "Target artifact:"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("handoff missing %q:\n%s", want, string(data))
		}
	}

	if _, err := consensus.AppendSignoff(root, "sample", consensus.SignoffOptions{
		Agent:  "alpha",
		Status: "accept",
		Notes:  "manual signoff.",
	}); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"resume", "--dir", root, "--no-tui", entries[0].Name()}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("resume code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Validated pending signoff from alpha.") {
		t.Fatalf("resume stdout=%q", stdout.String())
	}
}

func TestResumeRejectsManualSignoffAfterExistingContentEdit(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)
	consensusPath := filepath.Join(root, protocol.DeckDir, "ideas", "sample", "consensus.md")
	data, err := os.ReadFile(consensusPath)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "## Signoffs", "## Context\nOriginal.\n\n## Signoffs", 1))
	if err := os.WriteFile(consensusPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{
		ID:         "alpha",
		Path:       alpha,
		Backend:    agents.ExternalLocal,
		LaunchMode: agents.LaunchManual,
	})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("manual code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	entries, err := os.ReadDir(filepath.Join(root, protocol.DeckDir, "runs"))
	if err != nil {
		t.Fatal(err)
	}
	tampered, err := os.ReadFile(consensusPath)
	if err != nil {
		t.Fatal(err)
	}
	tampered = []byte(strings.Replace(string(tampered), "Original.", "Changed.", 1))
	if err := os.WriteFile(consensusPath, tampered, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := consensus.AppendSignoff(root, "sample", consensus.SignoffOptions{
		Agent:  "alpha",
		Status: "accept",
		Notes:  "manual signoff.",
	}); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"resume", "--dir", root, "--no-tui", entries[0].Name()}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("resume code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "changed existing consensus content") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestConsensusRequestSignoffsRejectsAlreadySignedExplicitParticipant(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha", "beta"}, false, map[string]string{"alpha": "accept"})

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "--participants", "alpha", "sample"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "participant alpha already signed") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestConsensusRequestSignoffsReviewPath(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, true, nil)

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 0)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "--review", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	summary, err := consensus.Status(root, "sample", true)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != consensus.TriageReady || len(summary.Signoffs) != 1 {
		t.Fatalf("summary=%+v", summary)
	}
	data, err := os.ReadFile(filepath.Join(root, protocol.DeckDir, "ideas", "sample", "review", "consensus.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Signoff: alpha") {
		t.Fatalf("review consensus not signed:\n%s", string(data))
	}
}

func TestConsensusRequestSignoffsNonZeroAfterAppendFails(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "accept", 7)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "exited with error after appending valid signoff") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	summary, err := consensus.Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != consensus.TriageReady {
		t.Fatalf("summary=%+v", summary)
	}
}

func TestConsensusRequestSignoffsBlockStops(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeSignoffCLI(t, bin, "alpha", "block", 0)
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "alpha appended BLOCK signoff") {
		t.Fatalf("stderr=%q", stderr.String())
	}
	summary, err := consensus.Status(root, "sample", false)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triage != consensus.TriageBlocked {
		t.Fatalf("summary=%+v", summary)
	}
}

func TestConsensusRequestSignoffsRejectsForgedExtraSignoff(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha", "beta"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeForgedSignoffCLI(t, bin, "alpha", "beta")
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "--participants", "alpha", "sample"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "appended 2 signoff blocks; expected exactly one") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestConsensusRequestSignoffsRejectsExistingContentEdit(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)

	bin := t.TempDir()
	alpha := writeFakeRewriteSignoffCLI(t, bin, "alpha")
	writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})

	var stdout, stderr bytes.Buffer
	code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "changed existing consensus content outside the append-only suffix") {
		t.Fatalf("stderr=%q", stderr.String())
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

func TestActionCommandUsesActionRoundAndAvoidsHardcodedAgent(t *testing.T) {
	run := runstate.RunSummary{RunID: "run-1", IdeaSlug: "sample"}
	draft := actionCommand(run, runplan.NextAction{
		Kind:     runplan.KindDraftConsensus,
		IdeaSlug: "sample",
		Round:    "round-02",
	})
	if draft != "parley consensus draft --round 2 sample" {
		t.Fatalf("draft command=%q", draft)
	}
	if strings.Contains(draft, "codex") {
		t.Fatalf("draft command hardcodes agent: %q", draft)
	}

	finalize := actionCommand(run, runplan.NextAction{Kind: runplan.KindFinalize, IdeaSlug: "sample"})
	if finalize != "parley consensus finalize sample" {
		t.Fatalf("finalize command=%q", finalize)
	}
	if strings.Contains(finalize, "codex") {
		t.Fatalf("finalize command hardcodes agent: %q", finalize)
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

type fakeAgentConfig struct {
	ID         string
	Path       string
	Backend    string
	LaunchMode string
}

func writeAgentsLocalConfig(t *testing.T, root string, entries ...fakeAgentConfig) {
	t.Helper()
	var b strings.Builder
	for _, entry := range entries {
		backend := entry.Backend
		if backend == "" {
			backend = agents.ExternalLocal
		}
		fmt.Fprintf(&b, "[agents.%s]\n", entry.ID)
		fmt.Fprintf(&b, "command = %q\n", entry.Path)
		fmt.Fprintln(&b, "prompt_mode = \"stdin\"")
		fmt.Fprintf(&b, "external_backend = %q\n", backend)
		if entry.LaunchMode != "" {
			fmt.Fprintf(&b, "launch_mode = %q\n", entry.LaunchMode)
		}
		fmt.Fprintln(&b, "timeout_ms = 5000")
		fmt.Fprintln(&b)
	}
	path := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeConsensusIdea(t *testing.T, root, slug string, participants []string, review bool, signoffs map[string]string) {
	t.Helper()
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", slug)
	if err := os.MkdirAll(filepath.Join(ideaDir, "round-01"), 0o755); err != nil {
		t.Fatal(err)
	}
	prompt := fmt.Sprintf("---\nidea: %s\nparticipants: [%s]\nstatus: consensus\n---\n", slug, strings.Join(participants, ", "))
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(prompt), 0o644); err != nil {
		t.Fatal(err)
	}

	consensusDir := ideaDir
	if review {
		consensusDir = filepath.Join(ideaDir, "review")
		if err := os.MkdirAll(consensusDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "---\nidea: %s\n---\n\n## Signoffs\n", slug)
	for _, participant := range participants {
		status, ok := signoffs[participant]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "\n### Signoff: %s - 2026-05-13\nStatus: %s\nNotes: seeded signoff.\n", participant, status)
	}
	if err := os.WriteFile(filepath.Join(consensusDir, "consensus.md"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeFakeSignoffCLI(t *testing.T, dir, name, status string, exitCode int) string {
	t.Helper()
	path := filepath.Join(dir, name)
	counter := ""
	notes := name + " accepts."
	if status == "block" {
		notes = name + " blocks."
		counter = "Counter-proposal: revise the consensus.\n"
	}
	body := fmt.Sprintf(`#!/bin/sh
prompt=$(mktemp)
cat > "$prompt"
path=$(awk -F': ' '/^Consensus file to sign:/ {print $2; exit}' "$prompt")
if [ -z "$path" ]; then
  exit 3
fi
cat >> "$path" <<'SIGNOFF'

### Signoff: %[1]s - 2026-05-13
Status: %[2]s
Notes: %[3]s
%[4]sSIGNOFF
exit %[5]d
`, name, status, notes, counter, exitCode)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeFakeForgedSignoffCLI(t *testing.T, dir, name, forged string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	body := fmt.Sprintf(`#!/bin/sh
prompt=$(mktemp)
cat > "$prompt"
path=$(awk -F': ' '/^Consensus file to sign:/ {print $2; exit}' "$prompt")
cat >> "$path" <<'SIGNOFF'

### Signoff: %[1]s - 2026-05-13
Status: accept
Notes: %[1]s accepts.

### Signoff: %[2]s - 2026-05-13
Status: accept
Notes: forged signoff.
SIGNOFF
exit 0
`, name, forged)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeFakeRewriteSignoffCLI(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	body := fmt.Sprintf(`#!/bin/sh
prompt=$(mktemp)
cat > "$prompt"
path=$(awk -F': ' '/^Consensus file to sign:/ {print $2; exit}' "$prompt")
tmp=$(mktemp)
sed 's/## Signoffs/## Changed Signoffs/' "$path" > "$tmp"
mv "$tmp" "$path"
cat >> "$path" <<'SIGNOFF'

### Signoff: %[1]s - 2026-05-13
Status: accept
Notes: %[1]s accepts.
SIGNOFF
exit 0
`, name)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeFakeCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	path := filepath.Join(dir, name)
	body := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo '" + version + "'; exit 0; fi\ncat >/dev/null\nexit 0\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFakeParleyDeckSkill(t *testing.T, dir string) {
	t.Helper()
	path := filepath.Join(dir, "parley-deck-skill")
	body := `#!/bin/sh
if [ "$1" = "status" ]; then
  project=""
  while [ "$#" -gt 0 ]; do
    if [ "$1" = "--project" ]; then
      shift
      project="$1"
    fi
    shift
  done
  cat <<JSON
{
  "ok": true,
  "installer": {
    "version": "1.1.0",
    "source": "test"
  },
  "compatibility": {
    "status": "ok",
    "reasons": []
  },
  "project": {
    "metadataStatus": "valid",
    "projectArg": "$project"
  },
  "runtimeInstalls": []
}
JSON
  exit 0
fi
exit 2
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFakeLegacyParleyDeckSkill(t *testing.T, dir string) {
	t.Helper()
	path := filepath.Join(dir, "parley-deck-skill")
	body := `#!/bin/sh
if [ "$1" = "--version" ]; then
  echo '1.0.8'
  exit 0
fi
echo 'Unknown command: status' >&2
exit 1
`
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
idea=$(basename "$(dirname "$(dirname "$out")")")
cat > "$out" <<'ARTIFACT'
---
agent: codex
idea: REPLACE_IDEA
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
sed -i.bak "s/REPLACE_IDEA/$idea/" "$out"
rm -f "$out.bak"
exit 0
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}
