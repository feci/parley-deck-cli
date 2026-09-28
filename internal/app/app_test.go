package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/hitl"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runaction"
	"parley-deck-cli/internal/runcontrol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/runplan"
	"parley-deck-cli/internal/runstate"
	"parley-deck-cli/internal/steer"
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

func TestDefaultParticipantSelectionSkipsLegacyGemini(t *testing.T) {
	discovered := []agents.Discovery{
		{
			Spec:  agents.Spec{ID: "agy", LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"--print"}},
			Found: true,
		},
		{
			Spec:  agents.Spec{ID: "gemini", LaunchMode: agents.LaunchHeadless, HeadlessArgs: []string{"--prompt"}},
			Found: true,
		},
		{
			Spec:  agents.Spec{ID: "manual", LaunchMode: agents.LaunchManual, HeadlessArgs: []string{"--run"}},
			Found: true,
		},
	}

	got, err := selectedParticipantIDs(discovered, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, ",") != "agy,manual" {
		t.Fatalf("default participants=%v, want [agy manual]", got)
	}

	got, err = selectedParticipantIDs(discovered, "gemini", nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, ",") != "gemini" {
		t.Fatalf("explicit participants=%v, want [gemini]", got)
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
	parley, ok := payload["parley"].(map[string]any)
	if !ok {
		t.Fatalf("parley missing or not an object: %v", payload["parley"])
	}
	if parley["version"] != version {
		t.Fatalf("parley=%+v", parley)
	}
	skill, ok := payload["parley_deck_skill"].(map[string]any)
	if !ok {
		t.Fatalf("parley_deck_skill missing or not an object: %v", payload["parley_deck_skill"])
	}
	installer, ok := skill["installer"].(map[string]any)
	if !ok {
		t.Fatalf("skill installer missing or not an object: %v", skill["installer"])
	}
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
	// Checked assertions: an unchecked type assertion panics and aborts the
	// whole test binary, taking the wait/usage results with it (AC-BLD-1).
	skill, ok := payload["parley_deck_skill"].(map[string]any)
	if !ok {
		t.Fatalf("parley_deck_skill missing or not an object: %v", payload["parley_deck_skill"])
	}
	project, ok := skill["project"].(map[string]any)
	if !ok {
		t.Fatalf("skill project missing or not an object: %v", skill["project"])
	}
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
	skill, ok := payload["parley_deck_skill"].(map[string]any)
	if !ok {
		t.Fatalf("parley_deck_skill missing or not an object: %v", payload["parley_deck_skill"])
	}
	if skill["statusSupported"] != false {
		t.Fatalf("skill=%+v", skill)
	}
	installer, ok := skill["installer"].(map[string]any)
	if !ok {
		t.Fatalf("skill installer missing or not an object: %v", skill["installer"])
	}
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
	for _, want := range []string{"codex", "yes", "codex test 1.0", "configured", "workspace-write", "never", "cli-default"} {
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
	declareAppTestSource(t, root)
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
	// --no-preflight: this test exercises runtime recording with a single fake
	// agent; the §1 non-solo hard-stop (which a 1-participant set would trip) is
	// covered separately in preflight_test.go.
	code := Run([]string{"run", "--dir", root, "--no-tui", "--no-auto", "--no-ping", "--no-preflight", "--yes", "--participants", "codex", "Runtime task"}, &stdout, &stderr)
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

// TestRunDefaultsToAutoMode locks the flipped default: `parley run` with neither
// --auto nor --no-auto records mode "auto", and --no-auto records "hitl". The
// fake agent is discoverable but fails round-01 so the auto-driver is skipped and
// the test only inspects the recorded run mode.
func TestRunDefaultsToAutoMode(t *testing.T) {
	modeFor := func(t *testing.T, extraArgs ...string) string {
		t.Helper()
		root := t.TempDir()
		t.Setenv("PARLEY_HOME", t.TempDir())
		if err := protocol.InitWorkspace(root); err != nil {
			t.Fatal(err)
		}
		localConfig := filepath.Join(root, protocol.DeckDir, "agents.local.toml")
		if err := os.WriteFile(localConfig, []byte("\n[agents.codex]\nmodel = \"local-model\"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		bin := t.TempDir()
		writeFailingRoundAgentCLI(t, bin, "codex", "codex test 1.0")
		t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

		var stdout, stderr bytes.Buffer
		// --no-preflight: this test only inspects the recorded run mode with a
		// single fake agent; the §1 non-solo hard-stop is covered in preflight_test.go.
		args := append([]string{"run", "--dir", root, "--no-tui", "--no-ping", "--no-preflight", "--yes", "--participants", "codex"}, extraArgs...)
		Run(append(args, "Task"), &stdout, &stderr)

		entries, err := os.ReadDir(filepath.Join(root, protocol.DeckDir, "runs"))
		if err != nil || len(entries) != 1 {
			t.Fatalf("runs=%v err=%v", entries, err)
		}
		events, err := store.New(filepath.Join(root, protocol.DeckDir, "runs", entries[0].Name())).Load()
		if err != nil || len(events) == 0 || events[0].Type != "run.created" {
			t.Fatalf("events=%+v err=%v", events, err)
		}
		mode, _ := events[0].Data["mode"].(string)
		return mode
	}

	if got := modeFor(t); got != "auto" {
		t.Fatalf("default run mode = %q, want \"auto\" (auto-drive must be the default)", got)
	}
	if got := modeFor(t, "--no-auto"); got != "hitl" {
		t.Fatalf("--no-auto run mode = %q, want \"hitl\"", got)
	}
}

// TestRunParticipantsSubsetHardStopsSolo locks Fix 1: `parley run --participants
// codex` on a machine with multiple installed agents must hard-stop (exit 1) on
// the §1 non-solo rule — it must NOT open a 1-participant idea. Preflight now
// evaluates the EXACT selected set, not every discovered agent.
func TestRunParticipantsSubsetHardStopsSolo(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	// Two installed agents are discoverable, but only codex is selected.
	bin := t.TempDir()
	writeFakeRoundAgentCLI(t, bin, "codex", "codex test 1.0")
	writeFakeRoundAgentCLI(t, bin, "claude", "claude test 1.0")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	var stdout, stderr bytes.Buffer
	// --no-ping keeps it hermetic (presence == available); preflight is NOT skipped.
	code := Run([]string{"run", "--dir", root, "--no-tui", "--no-auto", "--no-ping", "--yes", "--participants", "codex", "Solo task"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d want 1 (§1 non-solo hard-stop), stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	// No half-open idea: the run must stop before runcontrol.Create.
	if entries, _ := os.ReadDir(filepath.Join(root, protocol.DeckDir, "runs")); len(entries) != 0 {
		t.Fatalf("a run was created despite the §1 hard-stop: %v", entries)
	}
	ideas, _ := os.ReadDir(filepath.Join(root, protocol.DeckDir, "ideas"))
	if len(ideas) != 0 {
		t.Fatalf("a 1-participant idea was opened despite the §1 hard-stop: %v", ideas)
	}
}

// writeFailingRoundAgentCLI is discoverable (answers --version) but exits nonzero
// on a real invocation, so round-01 fails and the auto-driver is skipped.
func writeFailingRoundAgentCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	// §D.9 re-exec port: discoverable (--version) but exits nonzero on a real
	// invocation — the role spec carries the per-test version.
	writeRoleFixture(t, dir, name, "version-or-fail "+version)
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

func TestSteerQueuesAndContinueSurfacesIt(t *testing.T) {
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
	runID := "20260604T120000.000000000Z"
	runDir := filepath.Join(root, protocol.DeckDir, "runs", runID)
	base := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	s := store.New(runDir)
	for _, event := range []store.Event{
		{Time: base, Type: "run.created", Data: map[string]any{"idea": "sample", "participants": []string{"codex"}}},
		{Time: base.Add(time.Second), Type: "run.segment_started", Data: map[string]any{"segment_id": "segment-0001", "reason": "initial", "targets": []string{"codex"}}},
		{Time: base.Add(2 * time.Second), Type: "agent.started", Data: map[string]any{"agent": "codex", "segment_id": "segment-0001"}},
	} {
		if err := s.Append(event); err != nil {
			t.Fatal(err)
		}
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"steer", "--dir", root, "--agent", "codex", runID, "--", "focus on the parser"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("steer code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Recorded steer-0001-") || !strings.Contains(stdout.String(), "for codex") {
		t.Fatalf("steer output=%q", stdout.String())
	}

	queued, err := steer.List(runDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(queued) != 1 || queued[0].Text != "focus on the parser" || queued[0].SegmentID != "segment-0001" {
		t.Fatalf("queued=%+v", queued)
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"continue", "--dir", root, "--json", runID}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("continue code=%d stderr=%s", code, stderr.String())
	}
	for _, want := range []string{`"steers"`, "steer-0001", "focus on the parser"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("continue --json missing %q:\n%s", want, stdout.String())
		}
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
	// finalize is two steps since codex-1/F5: the first writes the scaffold and says the idea is
	// NOT closed; the second closes it once the scaffold has been written up.
	code = Run([]string{"consensus", "finalize", "--dir", root, "--by", "codex", "sample"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("finalize code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "NOT closed") {
		t.Fatalf("scaffold step must not read as a closure:\n%s", stdout.String())
	}
	finalPath := filepath.Join(ideaDir, "FINAL.md")
	if _, err := os.Stat(finalPath); err != nil {
		t.Fatal(err)
	}
	meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		t.Fatal(err)
	}
	if meta["status"] == "final" {
		t.Fatal("the idea was closed around an unwritten scaffold")
	}

	var written strings.Builder
	written.WriteString("---\nidea: sample\nstatus: final\nauthor: codex\n---\n\n")
	for _, section := range protocol.RequiredFinalSections {
		written.WriteString(section + "\n\nReal content.\nSecond line.\nThird line.\n\n")
	}
	if err := os.WriteFile(finalPath, []byte(written.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	if code = Run([]string{"consensus", "finalize", "--dir", root, "--by", "codex", "sample"}, &stdout, &stderr); code != 0 {
		t.Fatalf("second finalize code=%d stderr=%s", code, stderr.String())
	}
	meta, err = protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
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
	declareAppTestSource(t, root)
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
	writeSourceRoleMetadata(t, root)
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
	writeSourceRoleMetadata(t, root)
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
	declareAppTestSource(t, root)
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
	declareAppTestSource(t, root)
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
	logs, err := filepath.Glob(filepath.Join(root, protocol.DeckDir, "runs", "*", "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	observed := 0
	for _, path := range logs {
		events, err := store.New(filepath.Dir(path)).Load()
		if err != nil {
			t.Fatal(err)
		}
		for _, event := range events {
			if event.Type == "agent.signoff.artifact-present-after-failure" {
				observed++
				if event.Data["agent"] != "alpha" || event.Data["signoff_status"] != "accept" || event.Data["canonical_signoff_status"] != consensus.StatusAccept || event.Data["artifact_sha256"] == "" {
					t.Fatalf("artifact evidence: %+v", event)
				}
			}
		}
	}
	if observed != 1 {
		t.Fatalf("artifact evidence events = %d", observed)
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
	for _, exitCode := range []int{0, 7} {
		t.Run(fmt.Sprint(exitCode), func(t *testing.T) {
			root := t.TempDir()
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareAppTestSource(t, root)
			writeConsensusIdea(t, root, "sample", []string{"alpha"}, false, nil)
			bin := t.TempDir()
			alpha := writeFakeSignoffCLI(t, bin, "alpha", "block", exitCode)
			writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})
			var stdout, stderr bytes.Buffer
			code := Run([]string{"consensus", "request-signoffs", "--dir", root, "sample"}, &stdout, &stderr)
			if code != 1 || !strings.Contains(stderr.String(), "alpha appended BLOCK signoff") {
				t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
			summary, err := consensus.Status(root, "sample", false)
			if err != nil {
				t.Fatal(err)
			}
			if summary.Triage != consensus.TriageBlocked {
				t.Fatalf("summary=%+v", summary)
			}
			logs, err := filepath.Glob(filepath.Join(root, protocol.DeckDir, "runs", "*", "events.jsonl"))
			if err != nil || len(logs) != 1 {
				t.Fatalf("logs=%v err=%v", logs, err)
			}
			events, err := store.New(filepath.Dir(logs[0])).Load()
			if err != nil {
				t.Fatal(err)
			}
			if len(events) == 0 || events[0].Type != "run.created" || events[0].Data["idea"] != "sample" || events[0].Data["mode"] != "consensus-signoff" {
				t.Fatalf("missing headless run identity: %+v", events)
			}
			observed := 0
			raw, err := os.ReadFile(summary.Path)
			if err != nil {
				t.Fatal(err)
			}
			for _, event := range events {
				if event.Type == "agent.signoff.block-recorded" {
					observed++
					if event.Data["agent"] != "alpha" || event.Data["signoff_status"] != "block" || event.Data["canonical_signoff_status"] != consensus.StatusBlock || event.Data["process_failed"] != (exitCode != 0) || event.Data["artifact_sha256"] != sha256Hex(string(raw)) {
						t.Fatalf("wrong BLOCK evidence: %+v", event)
					}
				}
				if event.Type == "agent.signoff.artifact-present-after-failure" {
					t.Fatalf("BLOCK recorded as acceptance: %+v", event)
				}
			}
			if observed != 1 {
				t.Fatalf("BLOCK events=%d", observed)
			}
		})
	}
}

func TestConsensusRequestSignoffsFailedInvalidAppendHasNoArtifactEvent(t *testing.T) {
	for _, appendMode := range []string{"forged", "absent"} {
		t.Run(appendMode, func(t *testing.T) {
			root := t.TempDir()
			if err := protocol.InitWorkspace(root); err != nil {
				t.Fatal(err)
			}
			declareAppTestSource(t, root)
			writeConsensusIdea(t, root, "sample", []string{"alpha", "beta"}, false, nil)
			bin := t.TempDir()
			alpha := filepath.Join(bin, "alpha")
			if appendMode == "forged" {
				alpha = writeFakeForgedSignoffCLI(t, bin, "alpha", "beta")
				// The historical test edited the script's exit code to 7;
				// the re-exec port edits the role spec the same way.
				spec := filepath.Join(bin, "alpha.role")
				body, err := os.ReadFile(spec)
				if err != nil {
					t.Fatal(err)
				}
				body = bytes.ReplaceAll(body, []byte("forged-signoff "), []byte("forged-signoff-exit7 "))
				if err := os.WriteFile(spec, body, 0o644); err != nil {
					t.Fatal(err)
				}
			} else {
				// §D.9 re-exec port: drain stdin, exit 7.
				alpha = writeRoleFixture(t, bin, "alpha", "drain-exit 7")
			}
			writeAgentsLocalConfig(t, root, fakeAgentConfig{ID: "alpha", Path: alpha, Backend: agents.ExternalLocal})
			var stdout, stderr bytes.Buffer
			if code := Run([]string{"consensus", "request-signoffs", "--dir", root, "--participants", "alpha", "sample"}, &stdout, &stderr); code != 1 {
				t.Fatalf("invalid failure code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
			logs, err := filepath.Glob(filepath.Join(root, protocol.DeckDir, "runs", "*", "events.jsonl"))
			if err != nil {
				t.Fatal(err)
			}
			for _, path := range logs {
				events, err := store.New(filepath.Dir(path)).Load()
				if err != nil {
					t.Fatal(err)
				}
				for _, event := range events {
					if event.Type == "agent.signoff.artifact-present-after-failure" || event.Type == "agent.signoff.block-recorded" {
						t.Fatalf("invalid append certified: %+v", event)
					}
				}
			}
		})
	}
}

func TestConsensusRequestSignoffsRejectsForgedExtraSignoff(t *testing.T) {
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareAppTestSource(t, root)
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
	declareAppTestSource(t, root)
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
	draft := runaction.Command(runplan.NextAction{
		Kind:     runplan.KindDraftConsensus,
		IdeaSlug: "sample",
		Round:    "round-02",
	}, run.RunID, run.IdeaSlug)
	if draft != "parley consensus draft --round 2 sample" {
		t.Fatalf("draft command=%q", draft)
	}
	if strings.Contains(draft, "codex") {
		t.Fatalf("draft command hardcodes agent: %q", draft)
	}

	finalize := runaction.Command(runplan.NextAction{Kind: runplan.KindFinalize, IdeaSlug: "sample"}, run.RunID, run.IdeaSlug)
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

// AF1: detaching the TUI must wait for in-flight launched runs (so they finish
// and record their session), not cancel/abandon them when the command returns.
func TestLaunchReaperWaitsForInFlightRuns(t *testing.T) {
	var reaper launchReaper
	release := make(chan struct{})
	reaped := make(chan struct{})
	reaper.track(func() {
		<-release // stand in for an N-launched run that is still running
		close(reaped)
	})

	waited := make(chan struct{})
	go func() {
		reaper.waitForActive(&bytes.Buffer{})
		close(waited)
	}()

	select {
	case <-waited:
		t.Fatal("waitForActive returned before the in-flight run finished (detach abandoned the run)")
	case <-time.After(50 * time.Millisecond):
	}

	close(release)
	select {
	case <-waited:
	case <-time.After(2 * time.Second):
		t.Fatal("waitForActive did not return after the run finished")
	}
	select {
	case <-reaped:
	default:
		t.Fatal("tracked reap did not run to completion")
	}
}

// AF2: the consensus/FINAL drafter must be an idea participant, not just any
// installed headless agent.
func TestFirstHeadlessAgentRestrictedToParticipants(t *testing.T) {
	discovered := []agents.Discovery{
		{Spec: agents.Spec{ID: "hermes", LaunchMode: agents.LaunchHeadless}, Found: true},
		{Spec: agents.Spec{ID: "codex", LaunchMode: agents.LaunchHeadless}, Found: true},
	}
	got, ok := firstHeadlessAgent(discovered, []string{"codex", "agy"}, nil)
	if !ok || got.ID != "codex" {
		t.Fatalf("got %q ok=%v, want codex (the only headless participant)", got.ID, ok)
	}
	if _, ok := firstHeadlessAgent(discovered, []string{"agy"}, nil); ok {
		t.Fatal("no discovered agent is a participant; expected no drafter selected")
	}
}

// AF6: the implementer is resolved from durable role metadata (IMPLEMENTATION.md /
// FINAL.md), not blindly participants[0].
func TestResolveImplementerFromRoleMetadata(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "FINAL.md"), []byte("---\nidea: x\nimplementer: agy\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := resolveImplementer(dir, []string{"codex", "agy"}); got != "agy" {
		t.Fatalf("got %q, want agy (FINAL.md implementer, not participants[0])", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "IMPLEMENTATION.md"), []byte("---\nidea: x\nimplementer: codex\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := resolveImplementer(dir, []string{"codex", "agy"}); got != "codex" {
		t.Fatalf("got %q, want codex (IMPLEMENTATION.md takes precedence)", got)
	}
	_ = os.Remove(filepath.Join(dir, "IMPLEMENTATION.md"))
	if err := os.WriteFile(filepath.Join(dir, "FINAL.md"), []byte("---\nidea: x\nimplementer: hermes\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := resolveImplementer(dir, []string{"codex", "agy"}); got != "codex" {
		t.Fatalf("got %q, want codex (fallback; hermes not a participant)", got)
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
	fmt.Fprintf(&b, "---\nidea: %s\n---\n\n", slug)
	if !review {
		// The Phase-3 gate requires the canonical template sections incl. the
		// §15.5/§15.6 drafter duties (lean-organizer A.4).
		for _, section := range protocol.RequiredConsensusSections {
			if section == "## Signoffs" {
				continue
			}
			fmt.Fprintf(&b, "%s\n\nSeeded content.\n\n", section)
		}
	}
	fmt.Fprintf(&b, "## Signoffs\n")
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
	// §D.9 re-exec port: the parametric signoff writer is the "signoff"
	// role spec (status accept/block, per-test exit code); notes and the
	// block counter-proposal are derived exactly as the shell did.
	return writeRoleFixture(t, dir, name, fmt.Sprintf("signoff %s %s %d", name, status, exitCode))
}
func writeFakeForgedSignoffCLI(t *testing.T, dir, name, forged string) string {
	t.Helper()
	// §D.9 re-exec port: the "forged-signoff" role spec appends the agent's
	// and the forged signoff blocks to the consensus path extracted from
	// the stdin prompt.
	return writeRoleFixture(t, dir, name, "forged-signoff "+name+" "+forged)
}
func writeFakeRewriteSignoffCLI(t *testing.T, dir, name string) string {
	t.Helper()
	// §D.9 re-exec port: the "rewrite-signoff" role spec rewrites the seeded
	// marker and appends the agent's signoff block.
	return writeRoleFixture(t, dir, name, "rewrite-signoff "+name)
}
func writeFakeCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	// §D.9 re-exec port: --version answers; otherwise drains stdin, exit 0.
	writeRoleFixture(t, dir, name, "version-or-drain "+version)
}
func writeFakeParleyDeckSkill(t *testing.T, dir string) {
	t.Helper()
	// §D.9 rows 3/16: the historical fixture was an extension-less
	// #!/bin/sh script, which Windows cannot exec (the hosted
	// "executable file not found in %PATH%" signature). The re-exec port
	// installs a copy of this test binary; TestMain dispatches the
	// "parley-deck-skill" basename role below. Behavior is identical to the
	// shell script on every platform.
	writeReexecFixture(t, dir, "parley-deck-skill")
}

// runFixtureRole dispatches a re-exec'd fixture copy (this binary renamed to
// the fixture name) to its fake behavior. The bool reports whether the
// basename is a known fixture role.
func runFixtureRole(base string) (int, bool) {
	switch strings.TrimSuffix(base, ".exe") {
	case "parley-deck-skill":
		if os.Getenv("PARLEY_FAKE_SKILL") == "legacy" {
			return fakeLegacyParleyDeckSkillMain(), true
		}
		return fakeParleyDeckSkillMain(), true
	}
	// Generic role-spec fixtures (§D.9): the behavior spec lives beside the
	// binary copy as <name>.role, so parametric fixtures (per-test versions,
	// exit modes) need no per-name dispatch. Resolved from os.Executable() —
	// the actual fixture location — because argv[0] may be a bare name
	// resolved through PATH.
	//
	// RECURSION GUARD (runaway-fixture postmortem, 2026-09-28): this binary
	// has been RENAMED (its basename is neither a dedicated role nor the
	// standard .test name), so it is a fixture copy. A fixture copy that
	// cannot resolve its role MUST exit loudly — falling through to m.Run()
	// would run the entire suite inside the fixture, whose tests install
	// more copies: unbounded recursion (two organizer SIGTERM events).
	if exe, err := os.Executable(); err == nil {
		exeBase := filepath.Base(exe)
		if spec, rerr := os.ReadFile(filepath.Join(filepath.Dir(exe), strings.TrimSuffix(exeBase, ".exe")+".role")); rerr == nil {
			if code, ok := runRoleSpec(strings.TrimSpace(string(spec))); ok {
				return code, true
			}
		}
		// Exempt the REAL test binary by its platform-true name: strip the
		// Windows .exe suffix FIRST, then require the standard .test stem —
		// "app.test.exe" must pass (the hosted false-positive that killed
		// the whole package in 0.068s), while any renamed fixture copy
		// ("codex.exe", "parley-deck-skill.exe", bare "git") fails.
		stem := strings.TrimSuffix(exeBase, ".exe")
		if !strings.HasSuffix(stem, ".test") {
			fmt.Fprintf(os.Stderr, "fixture copy %q could not resolve its role spec; refusing to enter the test framework\n", exeBase)
			return 70, true
		}
	}
	return 0, false
}

// runRoleSpec executes a generic fixture behavior spec:
//
//	version-or-fail <v>   --version answers <v>; any other invocation exits 1
//	version-or-drain <v>  --version answers <v>; otherwise drain stdin, exit 0
//	drain-exit <n>        drain stdin, exit <n>
func runRoleSpec(spec string) (int, bool) {
	fields := strings.Fields(spec)
	if len(fields) == 0 {
		return 0, false
	}
	args := os.Args[1:]
	rest := ""
	if len(spec) > len(fields[0]) {
		rest = strings.TrimSpace(spec[len(fields[0]):])
	}
	switch fields[0] {
	case "version-or-fail", "version-or-drain":
		if rest == "" {
			return 0, false
		}
		if len(args) > 0 && args[0] == "--version" {
			fmt.Println(rest)
			return 0, true
		}
		if fields[0] == "version-or-fail" {
			return 1, true
		}
		_, _ = io.Copy(io.Discard, os.Stdin)
		return 0, true
	case "drain-exit":
		if len(fields) != 2 {
			return 0, false
		}
		code, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, false
		}
		_, _ = io.Copy(io.Discard, os.Stdin)
		return code, true
	case "round-agent":
		// The writeFakeRoundAgentCLI port: --version answers (the version
		// may contain spaces — parsed as the raw remainder, the runaway
		// postmortem finding); otherwise the round-01 prompt is read from
		// stdin, the artifact path is taken from the "Create exactly this
		// file..." line, and the artifact is written with the idea slug
		// (two dirs above the artifact) substituted.
		if rest == "" {
			return 0, false
		}
		if len(args) > 0 && args[0] == "--version" {
			fmt.Println(rest)
			return 0, true
		}
		prompt, err := io.ReadAll(os.Stdin)
		if err != nil {
			return 3, true
		}
		out := ""
		const marker = "Create exactly this file and no other protocol artifact:"
		for _, line := range strings.Split(string(prompt), "\n") {
			// The shell's awk matched the marker UNANCHORED and took the
			// text after the following ": " — replicate exactly (a
			// prefix-anchored match missed indented prompt lines).
			if idx := strings.Index(line, marker); idx >= 0 {
				out = strings.TrimSpace(strings.TrimPrefix(line[idx+len(marker):], ": "))
				break
			}
		}
		if out == "" {
			return 3, true
		}
		idea := filepath.Base(filepath.Dir(filepath.Dir(out)))
		artifact := fmt.Sprintf(`---
agent: codex
idea: %s
round: 1
date: 2026-05-11
---

## Summary
Fake artifact.

## Proposed approach
Use the test helper.

## Existing alternatives
Searched the stdlib and the lockfile; nothing ships this. Hand-built route is correct.

## Concerns / open questions
None.

## Risks
None.
`, idea)
		if err := os.WriteFile(out, []byte(artifact), 0o644); err != nil {
			return 3, true
		}
		return 0, true
	case "forged-signoff":
		// The writeFakeForgedSignoffCLI port: the signoff prompt is read
		// from stdin, the consensus path from its "Consensus file to
		// sign:" line, and two signoff blocks (the agent's and the forged
		// one) are appended.
		if len(fields) != 3 {
			return 0, false
		}
		prompt, err := io.ReadAll(os.Stdin)
		if err != nil {
			return 1, true
		}
		consensus := ""
		for _, line := range strings.Split(string(prompt), "\n") {
			if after, ok := strings.CutPrefix(line, "Consensus file to sign:"); ok {
				consensus = strings.TrimSpace(after)
				break
			}
		}
		if consensus == "" {
			return 1, true
		}
		block := fmt.Sprintf("\n### Signoff: %s - 2026-05-13\nStatus: accept\nNotes: %s accepts.\n\n### Signoff: %s - 2026-05-13\nStatus: accept\nNotes: forged signoff.\n", fields[1], fields[1], fields[2])
		f, err := os.OpenFile(consensus, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return 1, true
		}
		defer f.Close()
		if _, err := f.WriteString(block); err != nil {
			return 1, true
		}
		return 0, true
	case "signoff":
		// The writeFakeSignoffCLI port: extract the consensus path from the
		// stdin prompt, append the signoff block (accept or block with the
		// counter-proposal), exit with the per-test code.
		if len(fields) != 4 {
			return 0, false
		}
		prompt, err := io.ReadAll(os.Stdin)
		if err != nil {
			return 3, true
		}
		consensus := ""
		for _, line := range strings.Split(string(prompt), "\n") {
			if after, ok := strings.CutPrefix(line, "Consensus file to sign:"); ok {
				consensus = strings.TrimSpace(after)
				break
			}
		}
		if consensus == "" {
			return 3, true
		}
		name, status := fields[1], fields[2]
		code, cerr := strconv.Atoi(fields[3])
		if cerr != nil {
			return 0, false
		}
		notes := name + " accepts."
		counter := ""
		if status == "block" {
			notes = name + " blocks."
			counter = "Counter-proposal: revise the consensus.\n"
		}
		block := fmt.Sprintf("\n### Signoff: %s - 2026-05-13\nStatus: %s\nNotes: %s\n%s", name, status, notes, counter)
		f, err := os.OpenFile(consensus, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return 3, true
		}
		defer f.Close()
		if _, err := f.WriteString(block); err != nil {
			return 3, true
		}
		return code, true
	case "zcode-probe":
		// The writeFakeZcode port: verify the EQUALS-form argv exactly
		// (--prompt=<v> as one token, --mode yolo --cwd root, argc 5), then
		// read the prompt from stdin and write "<sentinel>\nheadless probe
		// ok\n" (line 5 sentinel, line 2 output path) to that path; the
		// help-exit0 mode prints usage and exits 0 writing nothing.
		if len(fields) != 2 {
			return 0, false
		}
		if len(args) > 0 && args[0] == "--version" {
			fmt.Println("zcode-app-cli 0.0.0-test")
			return 0, true
		}
		if len(args) != 5 {
			fmt.Fprintf(os.Stderr, "argc=%d want 5\n", len(args))
			return 64, true
		}
		promptText, perr := strings.CutPrefix(args[0], "--prompt=")
		if !perr {
			fmt.Fprintf(os.Stderr, "arg1=%s want --prompt=<value>\n", args[0])
			return 64, true
		}
		if args[1] != "--mode" || args[2] != "yolo" || args[3] != "--cwd" || args[4] == "" {
			fmt.Fprintf(os.Stderr, "argv mismatch: %v\n", args)
			return 64, true
		}
		if fields[1] == "help-exit0" {
			fmt.Println("Usage: zcode [options]")
			return 0, true
		}
		// The prompt travels in argv (the EQUALS form), not stdin — the
		// shell extracted $PROMPT from "${1#--prompt=}".
		lines := strings.Split(promptText, "\n")
		if len(lines) < 5 {
			return 64, true
		}
		out, sentinel := lines[1], lines[4]
		if out == "" || sentinel == "" {
			return 64, true
		}
		if err := os.WriteFile(out, []byte(sentinel+"\nheadless probe ok\n"), 0o644); err != nil {
			return 64, true
		}
		return 0, true
	case "forged-signoff-exit7":
		// The historical test edited the shell script's exit code; the
		// re-exec port edits the role spec the same way (see the caller).
		if code, ok := runRoleSpec("forged-signoff " + strings.Join(fields[1:], " ")); ok {
			if code == 0 {
				return 7, true
			}
			return code, true
		}
		return 0, false
	case "rewrite-signoff":
		// The writeFakeRewriteSignoffCLI port: replace the seeded marker in
		// the consensus file, then append the agent's signoff block.
		if len(fields) != 2 {
			return 0, false
		}
		prompt, err := io.ReadAll(os.Stdin)
		if err != nil {
			return 1, true
		}
		consensus := ""
		for _, line := range strings.Split(string(prompt), "\n") {
			if after, ok := strings.CutPrefix(line, "Consensus file to sign:"); ok {
				consensus = strings.TrimSpace(after)
				break
			}
		}
		if consensus == "" {
			return 1, true
		}
		raw, err := os.ReadFile(consensus)
		if err != nil {
			return 1, true
		}
		edited := strings.ReplaceAll(string(raw), "Seeded content.", "Edited by agent.")
		block := fmt.Sprintf("\n### Signoff: %s - 2026-05-13\nStatus: accept\nNotes: %s accepts.\n", fields[1], fields[1])
		if err := os.WriteFile(consensus, []byte(edited+block), 0o644); err != nil {
			return 1, true
		}
		return 0, true
	}
	return 0, false
}

// fakeParleyDeckSkillMain is the Go port of the historical shell script:
// "status" prints the version-status JSON (with the --project argument
// substituted) and exits 0; anything else exits 2.
func fakeParleyDeckSkillMain() int {
	args := os.Args[1:]
	if len(args) == 0 || args[0] != "status" {
		return 2
	}
	project := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--project" && i+1 < len(args) {
			project = args[i+1]
		}
	}
	fmt.Printf(`{
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
    "projectArg": %q
  },
  "runtimeInstalls": []
}
`, project)
	return 0
}

// writeRoleFixture installs a re-exec fixture (a copy of this test binary)
// with a sibling behavior spec at dir/name.role, returning the installed
// path (the .exe suffix on Windows — direct-exec config paths must use it).
func writeRoleFixture(t *testing.T, dir, name, spec string) string {
	t.Helper()
	path := writeReexecFixture(t, dir, name)
	if err := os.WriteFile(filepath.Join(dir, name+".role"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// writeReexecFixture installs a copy of the test binary at dir/name as the
// §D.9 test-binary re-exec port of an extension-less shell fixture (the .exe
// suffix on Windows is resolved by LookPath through PATHEXT), returning the
// installed path.
func writeReexecFixture(t *testing.T, dir, name string) string {
	t.Helper()
	src, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	installed := filepath.Join(dir, name)
	in, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.OpenFile(installed, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
	return installed
}
func writeFakeLegacyParleyDeckSkill(t *testing.T, dir string) {
	t.Helper()
	// §D.9 re-exec port, as writeFakeParleyDeckSkill; the legacy behavior is
	// selected by the env marker the child inherits from this test process.
	t.Setenv("PARLEY_FAKE_SKILL", "legacy")
	writeReexecFixture(t, dir, "parley-deck-skill")
}

// fakeLegacyParleyDeckSkillMain ports the historical legacy shell script:
// --version answers 1.0.8; anything else is an unknown command, exit 1.
func fakeLegacyParleyDeckSkillMain() int {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--version" {
		fmt.Println("1.0.8")
		return 0
	}
	fmt.Fprintln(os.Stderr, "Unknown command: status")
	return 1
}

func writeFakeRoundAgentCLI(t *testing.T, dir, name, version string) {
	t.Helper()
	// §D.9 re-exec port: the awk-over-stdin round-01 artifact writer is now
	// the "round-agent" role spec (prompt on stdin, artifact path from the
	// "Create exactly this file..." line, idea slug substituted).
	writeRoleFixture(t, dir, name, "round-agent "+version)
}
func declareAppTestSource(t *testing.T, root string) {
	t.Helper()
	meta := filepath.Join(root, "parley-deck", "meta")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "version.json"), []byte(`{"protocolRole":"source"}`), 0o600); err != nil {
		t.Fatal(err)
	}
}
