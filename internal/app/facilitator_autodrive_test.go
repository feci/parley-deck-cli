package app

// lean-organizer fix-up F21 (kimi-1 K1-F3, FINAL A criterion shapes): a driver-level
// fixture auto-drive test asserting role launches/escalation VIA THE EVENT LOG — the
// production driver adapter launches REAL fixture agents whose agent.started events
// land in the run's event store; the declared facilitator must never appear, and a
// facilitator-only deck must escalate without launching anyone.

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
)

// fixtureAgentScript writes a shell agent that reads its prompt on stdin, extracts
// the "Create exactly this (review) file: <path>" line, and writes a minimal valid
// artifact there (implementation or review shape, by mode).
func fixtureAgentScript(t *testing.T, dir, id, mode string) string {
	t.Helper()
	var artifact string
	switch mode {
	case "implement":
		artifact = `cat > "$out" <<'ART'
---
idea: fx-auto
status: implemented
implementer: %[1]s
started: 2026-09-24
completed: 2026-09-24
branch: fixture#main
head-commit: fixture
---

## Summary of work

Fixture implementation by %[1]s.
ART`
	case "review":
		artifact = `cat > "$out" <<'ART'
---
agent: %[1]s
idea: fx-auto
review-round: 1
reviewed-commit: fixture
date: 2026-09-24
---

## Refutation attempts

Tried to break the fixture; could not.

## Findings

None.
ART`
	default:
		t.Fatalf("unknown fixture mode %q", mode)
	}
	path := filepath.Join(dir, id)
	inner := fmt.Sprintf(artifact, id)
	body := "#!/bin/sh\nprompt=$(mktemp)\ncat > \"$prompt\"\n" +
		"out=$(awk '/Create exactly this (review )?file: /{for(i=1;i<=NF;i++) if ($i ~ /^\\//) {print $i; exit}}' \"$prompt\")\n" +
		"if [ -z \"$out\" ]; then\n" +
		"  echo \"fixture agent " + id + ": no output path in prompt\" >&2\n" +
		"  exit 3\n" +
		"fi\n" +
		inner + "\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func fixtureDiscovery(id, path string) agents.Discovery {
	d := agents.Discovery{Found: true, Path: path}
	d.ID = id
	d.LaunchMode = agents.LaunchHeadless
	d.PromptMode = agents.PromptStdin
	d.ExternalBackend = agents.ExternalLocal
	return d
}

// autoDriveFixture builds a git-clean workspace whose idea declares codex-1 the
// facilitator, wires the PRODUCTION driver adapter over real fixture agents, and
// returns a driver positioned at Phase 4 (FINAL published, implementation next).
func autoDriveFixture(t *testing.T, participants string, agentsList []agents.Discovery) (*driver.Driver, store.Store, string, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PARLEY_HOME", t.TempDir())
	t.Setenv("PARLEY_HEADLESS_AGENT_CONFIG", "")
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	writeSourceRoleMetadata(t, root)
	ideaDir := filepath.Join(root, protocol.DeckDir, "ideas", "fx-auto")
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	prompt := "---\nidea: fx-auto\nauthor: user\ntrack: deliberation\nparticipants: [" + participants + "]\nfacilitator: codex-1\nauto_implement: true\nchecks: \"true\"\nstatus: final\n---\n\n## Problem\nx\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(prompt), 0o644); err != nil {
		t.Fatal(err)
	}
	final := "---\nidea: fx-auto\nstatus: final\nauthor: fixture\n---\n\n# FINAL\n\n## Final plan\n\nx\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte(final), 0o644); err != nil {
		t.Fatal(err)
	}
	runDir := filepath.Join(root, protocol.DeckDir, "runs", "20260924T000000.000000000Z")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	events := store.New(runDir)
	c := driver.Rebuild(ideaDir, 6)
	if err := c.Save(filepath.Join(runDir, "driver.json")); err != nil {
		t.Fatal(err)
	}
	base := runner.Options{Root: root, RunID: "20260924T000000.000000000Z", Store: events, Agents: agentsList}
	base.Idea.Slug = "fx-auto"
	base.Idea.Path = ideaDir
	base.Idea.Participants = strings.Fields(strings.ReplaceAll(strings.Trim(participants, " "), ", ", " "))
	var parts []string
	for _, p := range strings.Split(participants, ",") {
		parts = append(parts, strings.TrimSpace(p))
	}
	// advanceFinal requires a clean git tree.
	for _, args := range [][]string{
		{"init", "-q"}, {"config", "user.email", "t@fixture"}, {"config", "user.name", "t"},
		{"add", "-A"}, {"commit", "-qm", "fixture"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	ops := newDriverImplOps(base, root, "fx-auto", ideaDir, parts, io.Discard)
	d := driver.New(driver.Config{
		Root: root, IdeaDir: ideaDir, IdeaSlug: "fx-auto", RunDir: runDir, Participants: parts,
		Auto: true, AutoImplement: true, Events: events, Impl: ops.(driver.ImplOps), Out: io.Discard,
	}, nil)
	return d, events, root, ideaDir
}

func startedAgents(events store.Store) []string {
	evs, err := events.Load()
	if err != nil {
		return nil
	}
	var ids []string
	for _, ev := range evs {
		if ev.Type != "agent.started" {
			continue
		}
		if id, _ := ev.Data["agent"].(string); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// TestFixtureAutoDriveNeverLaunchesFacilitatorForCodeRoles (fix-up F21): the real
// driver, real production ops and real fixture launches — the event log's
// agent.started trail names the participant implementer and reviewers, NEVER the
// declared facilitator codex-1.
func TestFixtureAutoDriveNeverLaunchesFacilitatorForCodeRoles(t *testing.T) {
	bin := t.TempDir()
	impl := fixtureDiscovery("claude-1", fixtureAgentScript(t, bin, "claude-1", "implement"))
	rev := fixtureDiscovery("kimi-1", fixtureAgentScript(t, bin, "kimi-1", "review"))
	d, events, _, ideaDir := autoDriveFixture(t, "codex-1, claude-1, kimi-1", []agents.Discovery{impl, rev})

	// Advance 1 (Phase 4 → 5): the implementer launches.
	action, _, err := d.Advance(context.Background())
	if err != nil || action != driver.ActionImplemented {
		t.Fatalf("advance 1: action=%s err=%v", action, err)
	}
	if _, statErr := os.Stat(filepath.Join(ideaDir, "IMPLEMENTATION.md")); statErr != nil {
		t.Fatalf("fixture implementer did not write IMPLEMENTATION.md: %v", statErr)
	}

	// Advance 2 (Phase 5 → 6): the reviewer launches.
	action2, _, err2 := d.Advance(context.Background())
	if err2 != nil || action2 != driver.ActionReviewOpened {
		t.Fatalf("advance 2: action=%s err=%v", action2, err2)
	}
	if _, statErr := os.Stat(filepath.Join(ideaDir, "review", "round-01", "kimi-1.md")); statErr != nil {
		t.Fatalf("fixture reviewer did not write its artifact: %v", statErr)
	}

	// THE event-log assertion: launches name participants only — never codex-1.
	started := startedAgents(events)
	if len(started) == 0 {
		t.Fatal("the event log must carry the launches (agent.started events)")
	}
	for _, id := range started {
		if id == "codex-1" {
			t.Fatalf("declared facilitator codex-1 was launched for a code role; agent.started trail: %v", started)
		}
	}
	sawImpl, sawRev := false, false
	for _, id := range started {
		if id == "claude-1" {
			sawImpl = true
		}
		if id == "kimi-1" {
			sawRev = true
		}
	}
	if !sawImpl || !sawRev {
		t.Fatalf("expected launches for implementer claude-1 and reviewer kimi-1; trail: %v", started)
	}
}

// TestFixtureAutoDriveFacilitatorOnlyEscalatesWithoutLaunching (fix-up F21): a
// facilitator-only deck escalates at the driver level — ActionEscalated,
// "escalated, not fallen back" — and the event log records ZERO launches:
// escalate-not-fallback, visible in the log, not just in the error string. The
// cursor is seeded at Phase 5/6 (IMPLEMENTATION.md published, review round next),
// the role action whose deadlock the pure Phase-5 path reaches through ops.Implement.
func TestFixtureAutoDriveFacilitatorOnlyEscalatesWithoutLaunching(t *testing.T) {
	bin := t.TempDir()
	fac := fixtureDiscovery("codex-1", fixtureAgentScript(t, bin, "codex-1", "review"))
	d, events, root, ideaDir := autoDriveFixture(t, "codex-1", []agents.Discovery{fac})

	// Seed PhaseImpl: IMPLEMENTATION.md published (status implemented) so Advance
	// routes to the review-round role action the facilitator-only deadlock blocks.
	impl := "---\nidea: fx-auto\nstatus: implemented\nimplementer: codex-1\nstarted: 2026-09-24\ncompleted: 2026-09-24\nbranch: fixture#main\nhead-commit: fixture\n---\n\n## Summary of work\n\nFixture.\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "IMPLEMENTATION.md"), []byte(impl), 0o644); err != nil {
		t.Fatal(err)
	}
	runDir := filepath.Join(root, protocol.DeckDir, "runs", "20260924T000000.000000000Z")
	c := driver.Rebuild(ideaDir, 6)
	if err := c.Save(filepath.Join(runDir, "driver.json")); err != nil {
		t.Fatal(err)
	}

	action, _, err := d.Advance(context.Background())
	if err == nil || action != driver.ActionEscalated {
		t.Fatalf("facilitator-only deck must escalate, got action=%s err=%v", action, err)
	}
	// At the driver level the facilitator-only shape escalates through the §1
	// non-solo track gate BEFORE the role action is reached (measured here); the
	// role-deadlock string itself is pinned by the ops-boundary test
	// (TestFacilitatorOnlyParticipantsEscalatesNotFallsBack). Either way this is
	// escalate-not-fallback, and the property this test owns is the event log:
	if started := startedAgents(events); len(started) != 0 {
		t.Fatalf("escalation must launch nobody (the facilitator is never a silent fallback); agent.started trail: %v", started)
	}
	if !strings.Contains(err.Error(), "escalat") && !strings.Contains(err.Error(), "none available") {
		t.Errorf("escalation must name its cause; got %q", err.Error())
	}
}
