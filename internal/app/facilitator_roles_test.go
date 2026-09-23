package app

// lean-organizer A: the declared facilitator is never selected as drafter,
// implementer, reviewer, or goal-done checker, and a facilitator-only participant
// set escalates instead of silently falling back. The absent-field deck keeps the
// exact v1.48.0 role selection (regression row).

import (
	"context"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
)

func discoveryFor(ids ...string) []agents.Discovery {
	out := make([]agents.Discovery, 0, len(ids))
	for _, id := range ids {
		d := agents.Discovery{Found: true}
		d.ID = id
		d.LaunchMode = agents.LaunchHeadless
		out = append(out, d)
	}
	return out
}

func implOpsFor(t *testing.T, ideaDir string, participants ...string) driverImplOps {
	t.Helper()
	ops := newDriverImplOps(runner.Options{}, t.TempDir(), "slug", ideaDir, participants, nil)
	impl, ok := ops.(driverImplOps)
	if !ok {
		t.Fatalf("newDriverImplOps returned %T, want driverImplOps", ops)
	}
	return impl
}

func TestDeclaredFacilitatorNeverSelectedForCodeRoles(t *testing.T) {
	dir := t.TempDir()
	ideaDir := writeIdeaPrompt(t, dir, "declared-run",
		"idea: declared-run\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1, kimi-1]\nfacilitator: codex-1")
	ops := implOpsFor(t, ideaDir, "codex-1", "claude-1", "kimi-1")
	if ops.implementer == "codex-1" {
		t.Errorf("declared facilitator must never be implementer, got %q", ops.implementer)
	}
	for _, r := range ops.reviewers {
		if r == "codex-1" {
			t.Errorf("declared facilitator must never be a reviewer; reviewers=%v", ops.reviewers)
		}
	}
	if ops.drafter == "codex-1" {
		t.Errorf("declared facilitator must never be the review-consensus drafter / goal-done checker, got %q", ops.drafter)
	}
	if ops.roleErr != "" {
		t.Errorf("with non-facilitator participants present there is no role deadlock, got %q", ops.roleErr)
	}
}

func TestFacilitatorOnlyParticipantsEscalatesNotFallsBack(t *testing.T) {
	dir := t.TempDir()
	ideaDir := writeIdeaPrompt(t, dir, "fac-only",
		"idea: fac-only\nauthor: user\ntrack: standard\nparticipants: [codex-1]\nfacilitator: codex-1")
	ops := implOpsFor(t, ideaDir, "codex-1")
	if ops.roleErr == "" {
		t.Fatalf("facilitator-only participant set must escalate, got empty roleErr")
	}
	if !strings.Contains(ops.roleErr, "escalated, not fallen back") {
		t.Errorf("escalation must state escalate-not-fallback, got %q", ops.roleErr)
	}
	if err := ops.Implement(context.Background()); err == nil || !strings.Contains(err.Error(), "escalated, not fallen back") {
		t.Errorf("Implement must escalate the role deadlock, got %v", err)
	}
	if err := ops.OpenReviewRound(context.Background(), 1); err == nil || !strings.Contains(err.Error(), "escalated, not fallen back") {
		t.Errorf("OpenReviewRound must escalate the role deadlock, got %v", err)
	}
	if ok, reason := ops.GoalCheck(context.Background()); ok || !strings.Contains(reason, "escalated, not fallen back") {
		t.Errorf("GoalCheck must refuse under the role deadlock, got ok=%v reason=%q", ok, reason)
	}
}

func TestFacilitatorParticipatesKeepsRoleEligibility(t *testing.T) {
	dir := t.TempDir()
	ideaDir := writeIdeaPrompt(t, dir, "fac-participates",
		"idea: fac-participates\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1, kimi-1]\nfacilitator: codex-1\nfacilitator_participates: true")
	ops := implOpsFor(t, ideaDir, "codex-1", "claude-1", "kimi-1")
	if ops.implementer != "codex-1" {
		t.Errorf("facilitator_participates: true keeps default role selection (participants[0]), got implementer %q", ops.implementer)
	}
	if ops.roleErr != "" {
		t.Errorf("participating facilitator is not a deadlock, got %q", ops.roleErr)
	}
}

func TestAbsentFacilitatorFieldKeepsV148RoleSelection(t *testing.T) {
	dir := t.TempDir()
	ideaDir := writeIdeaPrompt(t, dir, "plain-run",
		"idea: plain-run\nauthor: user\ntrack: standard\nparticipants: [claude-1, kimi-1, zcode-1]")
	ops := implOpsFor(t, ideaDir, "claude-1", "kimi-1", "zcode-1")
	// v1.48.0: implementer = participants[0] (no durable metadata), reviewers = the
	// rest in order, drafter = reviewers[0]. Byte-identical selection, no facilitator
	// predicate interference.
	if ops.implementer != "claude-1" {
		t.Errorf("absent-field regression: implementer=%q want claude-1", ops.implementer)
	}
	if strings.Join(ops.reviewers, ",") != "kimi-1,zcode-1" {
		t.Errorf("absent-field regression: reviewers=%v want [kimi-1 zcode-1]", ops.reviewers)
	}
	if ops.drafter != "kimi-1" {
		t.Errorf("absent-field regression: drafter=%q want kimi-1", ops.drafter)
	}
}

func TestFirstEligibleHeadlessAgentSkipsFacilitator(t *testing.T) {
	dir := t.TempDir()
	ideaDir := writeIdeaPrompt(t, dir, "declared-run",
		"idea: declared-run\nauthor: user\ntrack: standard\nparticipants: [codex-1, claude-1]\nfacilitator: codex-1")
	role := readRoleForTest(t, ideaDir)
	agent, ok := firstEligibleHeadlessAgent(discoveryFor("codex-1", "claude-1"), []string{"codex-1", "claude-1"}, nil, role)
	if !ok || agent.ID != "claude-1" {
		t.Fatalf("consensus/FINAL drafter must skip the declared facilitator, got ok=%v agent=%q", ok, agent.ID)
	}
	_, ok = firstEligibleHeadlessAgent(discoveryFor("codex-1"), []string{"codex-1"}, nil, role)
	if ok {
		t.Errorf("facilitator-only set must report no eligible drafter (escalation), got one")
	}
}

func readRoleForTest(t *testing.T, dir string) protocol.FacilitatorRole {
	t.Helper()
	return protocol.ReadFacilitatorRole(dir)
}
