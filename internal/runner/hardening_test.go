package runner

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

// --- test-binary-as-agent helpers (same pattern as TestFakeAgentHelper) --------

// TestFakeExitAgentHelper writes a valid round-01 artifact, then exits with the
// code from FAKE_AGENT_EXIT (artifact-wins fixture).
func TestFakeExitAgentHelper(t *testing.T) {
	if !hasArg("parley-fake-exit-agent") {
		return
	}
	input, _ := io.ReadAll(os.Stdin)
	re := regexp.MustCompile(`(?m)^- Create exactly this file and no other protocol artifact: (.+)$`)
	match := re.FindStringSubmatch(string(input))
	ideaMatch := regexp.MustCompile(`(?m)^idea: (.+)$`).FindStringSubmatch(string(input))
	if len(match) != 2 || len(ideaMatch) != 2 {
		os.Exit(7)
	}
	body := "---\nagent: fake\nidea: " + ideaMatch[1] + "\nround: 1\ndate: 2026-06-12\n---\n\n## Summary\nx\n\n## Proposed approach\nx\n\n## Concerns / open questions\nx\n\n## Risks\nx\n"
	_ = os.WriteFile(match[1], []byte(body), 0o644)
	code, _ := strconv.Atoi(os.Getenv("FAKE_AGENT_EXIT"))
	os.Exit(code)
}

// TestFakeSilentAgentHelper produces no output and no artifact; it sleeps until
// killed (first-output watchdog fixture).
func TestFakeSilentAgentHelper(t *testing.T) {
	if !hasArg("parley-fake-silent") {
		return
	}
	_, _ = io.ReadAll(os.Stdin)
	time.Sleep(30 * time.Second)
	os.Exit(0)
}

// TestFakeStallAgentHelper prints one line immediately, then goes silent until
// killed (stall-guard fixture).
func TestFakeStallAgentHelper(t *testing.T) {
	if !hasArg("parley-fake-stall") {
		return
	}
	_, _ = io.ReadAll(os.Stdin)
	os.Stdout.WriteString("warming up\n")
	_ = os.Stdout.Sync()
	time.Sleep(30 * time.Second)
	os.Exit(0)
}

func fakeAgent(helper string, spec agents.Spec) agents.Discovery {
	spec.HeadlessArgs = []string{"-test.run=" + helper, "--", strings.TrimPrefix(helperArg(helper), "")}
	spec.PromptMode = agents.PromptStdin
	return agents.Discovery{Spec: spec, Path: os.Args[0], Found: true}
}

func helperArg(helper string) string {
	switch helper {
	case "TestFakeExitAgentHelper":
		return "parley-fake-exit-agent"
	case "TestFakeSilentAgentHelper":
		return "parley-fake-silent"
	case "TestFakeStallAgentHelper":
		return "parley-fake-stall"
	}
	return helper
}

func setupRunnerIdea(t *testing.T) (string, protocol.IdeaStatus, store.Store) {
	t.Helper()
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	idea, err := protocol.CreateIdea(root, "Hardening test task", []string{"fake"})
	if err != nil {
		t.Fatal(err)
	}
	return root, idea, store.New(filepath.Join(root, protocol.DeckDir, "runs", "test-run"))
}

func eventsOfType(t *testing.T, st store.Store, typ string) []store.Event {
	t.Helper()
	events, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	var out []store.Event
	for _, e := range events {
		if e.Type == typ {
			out = append(out, e)
		}
	}
	return out
}

// --- D7: artifact beats exit code ------------------------------------------------

func TestArtifactWinsOnNonzeroExit(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	t.Setenv("FAKE_AGENT_EXIT", "3")
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents:  []agents.Discovery{fakeAgent("TestFakeExitAgentHelper", agents.Spec{ID: "fake"})},
		Timeout: 10 * time.Second, Store: st,
	})
	if len(results) != 1 {
		t.Fatalf("results=%d", len(results))
	}
	r := results[0]
	if !r.Success() || !r.ArtifactOK || r.ExitError != "" {
		t.Fatalf("artifact-wins result not successful: %+v", r)
	}
	if r.AgentExit != 3 || r.AgentExitKind != "exec" {
		t.Fatalf("AgentExit=%d kind=%q, want 3/exec", r.AgentExit, r.AgentExitKind)
	}
	finished := eventsOfType(t, st, "agent.finished")
	if len(finished) != 1 {
		t.Fatalf("agent.finished=%d, want 1", len(finished))
	}
	if got := finished[0].Data["agent_exit"]; got != 3 && got != int64(3) && got != float64(3) {
		t.Fatalf("agent_exit=%v, want 3", got)
	}
	if rounds := eventsOfType(t, st, "round.completed"); len(rounds) != 1 {
		t.Fatal("artifact-wins must still complete the round")
	}
}

// --- D1/D3: first-output watchdog with one retry ----------------------------------

func TestNoFirstOutputWatchdogRetriesThenFails(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents: []agents.Discovery{fakeAgent("TestFakeSilentAgentHelper", agents.Spec{
			ID: "fake", FirstEventTimeoutMS: 100, StallTimeoutMS: -1, HeartbeatMS: -1,
		})},
		Timeout: 20 * time.Second, Store: st,
	})
	r := results[0]
	if r.Success() {
		t.Fatalf("silent agent must fail: %+v", r)
	}
	if r.FailureClass != "no_first_output" {
		t.Fatalf("failure_class=%q, want no_first_output", r.FailureClass)
	}
	started := eventsOfType(t, st, "agent.started")
	if len(started) != 2 {
		t.Fatalf("agent.started=%d, want 2 (one retry)", len(started))
	}
	watchdog := eventsOfType(t, st, "agent.no_first_output")
	if len(watchdog) != 2 {
		t.Fatalf("agent.no_first_output=%d, want 2", len(watchdog))
	}
	if watchdog[0].Data["action"] != "retrying" || watchdog[1].Data["action"] != "failed" {
		t.Fatalf("watchdog actions=%v,%v want retrying,failed", watchdog[0].Data["action"], watchdog[1].Data["action"])
	}
	failed := eventsOfType(t, st, "agent.failed")
	if len(failed) != 2 {
		t.Fatalf("agent.failed=%d, want 2 (each killed attempt is terminal)", len(failed))
	}
	// Ordering (D1): each watchdog event must precede its attempt's terminal failure.
	all, _ := st.Load()
	firstWatchdog, firstFailed := -1, -1
	for i, e := range all {
		if e.Type == "agent.no_first_output" && firstWatchdog == -1 {
			firstWatchdog = i
		}
		if e.Type == "agent.failed" && firstFailed == -1 {
			firstFailed = i
		}
	}
	if firstWatchdog == -1 || firstWatchdog > firstFailed {
		t.Fatalf("watchdog event must be appended before the kill's terminal event (wd=%d failed=%d)", firstWatchdog, firstFailed)
	}
}

// --- D1: stall guard after first output (no retry) -------------------------------

func TestStallGuardKillsAfterFirstOutput(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents: []agents.Discovery{fakeAgent("TestFakeStallAgentHelper", agents.Spec{
			ID: "fake", FirstEventTimeoutMS: -1, StallTimeoutMS: 100, HeartbeatMS: -1,
		})},
		Timeout: 20 * time.Second, Store: st,
	})
	r := results[0]
	if r.Success() || r.FailureClass != "stalled" {
		t.Fatalf("want stalled failure, got %+v", r)
	}
	if got := len(eventsOfType(t, st, "agent.started")); got != 1 {
		t.Fatalf("agent.started=%d, want 1 (stall is never retried)", got)
	}
	if got := len(eventsOfType(t, st, "agent.stalled")); got != 1 {
		t.Fatalf("agent.stalled=%d, want 1", got)
	}
}

// --- D4: heartbeats persisted, never counted as activity --------------------------

func TestHeartbeatEventsEmitted(t *testing.T) {
	root, idea, st := setupRunnerIdea(t)
	results := RunRoundOne(context.Background(), Options{
		Root: root, RunID: "test-run", Idea: idea, Task: "t",
		Agents: []agents.Discovery{fakeAgent("TestFakeStallAgentHelper", agents.Spec{
			ID: "fake", FirstEventTimeoutMS: -1, StallTimeoutMS: 2500, HeartbeatMS: 600,
		})},
		Timeout: 20 * time.Second, Store: st,
	})
	if results[0].FailureClass != "stalled" {
		t.Fatalf("fixture should end stalled, got %+v", results[0])
	}
	beats := eventsOfType(t, st, "agent.heartbeat")
	if len(beats) == 0 {
		t.Fatal("want at least one agent.heartbeat before the stall kill")
	}
	data := beats[0].Data
	for _, key := range []string{"elapsed_ms", "stdout_bytes", "stderr_bytes", "attempt_id", "phase"} {
		if _, ok := data[key]; !ok {
			t.Fatalf("heartbeat payload missing %s: %+v", key, data)
		}
	}
}

// --- D2: supervision config derivation -------------------------------------------

func TestSupervisionForAgent(t *testing.T) {
	def := supervisionForAgent(agents.Discovery{Spec: agents.Spec{ID: "x"}}, 30*time.Minute)
	if def.FirstEventTimeout != 120*time.Second || def.HeartbeatInterval != 60*time.Second {
		t.Fatalf("defaults wrong: %+v", def)
	}
	if def.StallTimeout != 30*time.Minute-time.Second {
		t.Fatalf("stall must clamp under the hard timeout: %v", def.StallTimeout)
	}
	off := supervisionForAgent(agents.Discovery{Spec: agents.Spec{ID: "x", FirstEventTimeoutMS: -1, StallTimeoutMS: -1, HeartbeatMS: -1}}, time.Hour)
	if off.FirstEventTimeout != 0 || off.StallTimeout != 0 || off.HeartbeatInterval != 0 {
		t.Fatalf("explicit -1 must disable: %+v", off)
	}
}

// --- D5: failure classification ---------------------------------------------------

func TestClassifyFailure(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	cases := []struct {
		stderr string
		want   string
	}{
		{"Error: authentication_failed — run login", "auth"},
		{"429 Too Many Requests, rate limit hit", "rate-limit"},
		{"upstream returned status 503 server error", "overloaded"},
		{"prompt is too long: contextWindowExceeded", "context-window"},
		{"something inscrutable happened", "unknown"},
	}
	for _, tc := range cases {
		p := write("stderr.log", tc.stderr)
		class, hint := classifyFailure(p, "", "exit status 1")
		if class != tc.want {
			t.Errorf("stderr %q → class %q, want %q", tc.stderr, class, tc.want)
		}
		if hint == "" {
			t.Errorf("class %q must carry a hint", class)
		}
	}

	// Lock the agreed class/hint contract (review consensus cycle-1 fix 7):
	// the table — not prose elsewhere — is the UX source of truth.
	exact := map[string]string{
		"rate-limit":      "Wait for reset or switch provider keys/endpoints.",
		"auth":            "Run the agent CLI's auth command (e.g. 'claude login', 'hermes auth') to refresh credentials.",
		"overloaded":      "Retry in a few minutes or choose a less busy model.",
		"context-window":  "Reduce the prompt size or prune file attachments/logs from scope.",
		"billing":         "Check your API account balance and credit card status.",
		"model-not-found": "Check the model spelling and access permissions in your API settings.",
		"sandbox":         "Adjust the local sandbox configuration or run with lower restriction.",
		"budget":          "Increase the session budget limit (e.g. raise spend caps in settings).",
		"invalid-request": "Verify the prompt structure and system constraints in config.",
	}
	for _, rule := range failureRules {
		want, ok := exact[rule.class]
		if !ok {
			t.Errorf("class %q is not in the locked contract table — update the test deliberately", rule.class)
			continue
		}
		if rule.hint != want {
			t.Errorf("class %q hint drifted: got %q want %q", rule.class, rule.hint, want)
		}
	}
	for class, want := range map[string]string{
		"no_first_output": "Verify the agent executable is not blocking or waiting for stdin.",
		"stalled":         "Check the process tree; the agent emitted no output within the stall window.",
		"timeout":         "The hard per-agent timeout elapsed; raise timeout_ms or split the task.",
	} {
		if got := watchdogHints[class]; got != want {
			t.Errorf("watchdog hint %q drifted: got %q want %q", class, got, want)
		}
	}
}

// TestMoveAsideInvalidArtifact (review fix 3): unique destination, never
// overwrite an earlier recovery file, remove-on-rename-failure.
func TestMoveAsideInvalidArtifact(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "agent.md")
	prior := out + ".attempt-1.invalid"
	if err := os.WriteFile(prior, []byte("earlier recovery"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(out, []byte("new invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	moveAsideInvalidArtifact(out)
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("invalid artifact must leave the canonical path, stat err=%v", err)
	}
	if data, err := os.ReadFile(prior); err != nil || string(data) != "earlier recovery" {
		t.Fatalf("pre-existing recovery file must never be overwritten: %q err=%v", data, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 2 { // prior + uniquely-suffixed new one
		t.Fatalf("want 2 recovery files, got %d", len(entries))
	}

	// Rename failure: a source basename near NAME_MAX makes the
	// ".attempt-1.invalid" destination exceed it, so the rename fails
	// (ENAMETOOLONG) while the source itself stays removable — the invalid
	// artifact must then be removed from the canonical path.
	failDir := t.TempDir()
	longOut := filepath.Join(failDir, strings.Repeat("x", 250)+".md")
	if err := os.WriteFile(longOut, []byte("invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	moveAsideInvalidArtifact(longOut)
	if _, err := os.Stat(longOut); !os.IsNotExist(err) {
		t.Fatalf("rename failure must remove the invalid artifact, stat err=%v", err)
	}
	failEntries, _ := os.ReadDir(failDir)
	if len(failEntries) != 0 {
		t.Fatalf("no recovery file should exist after a failed rename, got %v", failEntries)
	}
}

// --- D8: participant env shedding --------------------------------------------------

func TestCleanParticipantEnv(t *testing.T) {
	env := []string{
		"CLAUDECODE=1", "CLAUDE_CODE_SESSION_ID=abc", "CLAUDE_CODE_ENTRYPOINT=cli",
		"AI_AGENT=claude", "AI_AGENT_FOO=x", "PARLEY_KEEP=1", "PATH=/usr/bin",
	}
	got := cleanParticipantEnv("claude", env)
	for _, kv := range got {
		key, _, _ := strings.Cut(kv, "=")
		if key == "CLAUDECODE" || key == "AI_AGENT" || strings.HasPrefix(key, "CLAUDE_CODE_") || strings.HasPrefix(key, "AI_AGENT_") {
			t.Fatalf("marker %s survived shedding: %v", key, got)
		}
	}
	if len(got) != 2 { // PARLEY_KEEP + PATH
		t.Fatalf("kept=%v, want PARLEY_KEEP and PATH only", got)
	}
	if other := cleanParticipantEnv("codex", env); len(other) != len(env) {
		t.Fatalf("non-claude participants must keep their env untouched")
	}
}

// --- D7: fix-up artifact validation -------------------------------------------------

func TestValidateFixupArtifact(t *testing.T) {
	dir := t.TempDir()
	impl := filepath.Join(dir, "IMPLEMENTATION.md")
	valid := "---\nidea: demo\nstatus: fix-up-cycle-1\n---\n\n## Summary of work\nx\n\n## Fix-up cycle 1\nfixed things\n"
	if err := os.WriteFile(impl, []byte(valid), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateFixupArtifact(dir, "demo"); err != nil {
		t.Fatalf("valid fix-up rejected: %v", err)
	}
	if err := ValidateFixupArtifact(dir, "other"); err == nil {
		t.Fatal("mismatched idea slug must fail")
	}
	bad := strings.Replace(valid, "fix-up-cycle-1", "drafting", 1)
	_ = os.WriteFile(impl, []byte(bad), 0o644)
	if err := ValidateFixupArtifact(dir, "demo"); err == nil {
		t.Fatal("non-review-ready status must fail")
	}
	noSection := "---\nidea: demo\nstatus: implemented\n---\n\n## Summary of work\nx\n"
	_ = os.WriteFile(impl, []byte(noSection), 0o644)
	if err := ValidateFixupArtifact(dir, "demo"); err == nil {
		t.Fatal("missing fix-up section must fail")
	}
}

// --- D9: review snapshot lifecycle ---------------------------------------------------

func TestReviewSnapshotLifecycle(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	git("init", "--quiet")
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("v1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "-A")
	git("commit", "--quiet", "-m", "init")

	// Dirty the tree: a modified tracked file AND an untracked one — the
	// snapshot commit must capture both.
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("v2-dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "untracked.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	snap, err := CreateReviewSnapshot(repo, "demo", "review/round-01", "codex", "run-1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer snap.Cleanup()
	if snap.Mode != "snapshot-commit" {
		t.Fatalf("mode=%q, want snapshot-commit for a dirty tree", snap.Mode)
	}
	if data, err := os.ReadFile(filepath.Join(snap.Dir, "tracked.txt")); err != nil || string(data) != "v2-dirty\n" {
		t.Fatalf("snapshot must carry the dirty content: %q err=%v", data, err)
	}
	if _, err := os.Stat(filepath.Join(snap.Dir, "untracked.txt")); err != nil {
		t.Fatalf("snapshot must carry untracked files: %v", err)
	}

	// Artifact move-back: write inside the snapshot, publish to the live repo.
	rel := filepath.Join("review", "round-01", "codex.md")
	if err := os.MkdirAll(filepath.Join(snap.Dir, "review", "round-01"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(snap.Dir, rel), []byte("review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	canonical := filepath.Join(repo, rel)
	if err := snap.MoveArtifactBack(rel, canonical); err != nil {
		t.Fatalf("move back: %v", err)
	}
	if data, err := os.ReadFile(canonical); err != nil || string(data) != "review\n" {
		t.Fatalf("canonical artifact wrong: %q err=%v", data, err)
	}

	dir := snap.Dir
	snap.Cleanup()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("cleanup must remove the snapshot dir, stat err=%v", err)
	}

	// Non-git roots fall back loudly instead of erroring hard.
	if _, err := CreateReviewSnapshot(t.TempDir(), "demo", "review/round-01", "codex", "run-1"); err == nil {
		t.Fatal("non-git root must report snapshot unavailable")
	}
}
