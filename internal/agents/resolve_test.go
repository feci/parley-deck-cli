package agents

import "testing"

func disco(id string, found bool) Discovery {
	return Discovery{Spec: Spec{ID: id}, Found: found}
}

func TestResolveExactSpecID(t *testing.T) {
	// Legacy deck: participant IS a family/spec id.
	got, err := ResolveParticipant("claude", []Discovery{disco("claude", true), disco("codex", true)}, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got.ID != "claude" || got.Adapter() != "claude" {
		t.Fatalf("id=%q adapter=%q", got.ID, got.Adapter())
	}
}

func TestResolveViaMapping(t *testing.T) {
	// New deck: participant is a roster ID mapped to a family.
	mapping := map[string]string{"claude-1": "claude", "antigravity-1": "agy"}
	got, err := ResolveParticipant("claude-1", []Discovery{disco("claude", true), disco("agy", true)}, mapping)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Identity is the roster ID; adapter is the family (for vendor/launch dispatch).
	if got.ID != "claude-1" {
		t.Errorf("identity = %q, want claude-1", got.ID)
	}
	if got.Adapter() != "claude" {
		t.Errorf("adapter = %q, want claude", got.Adapter())
	}
	// The prefix-heuristic-breaking case must work only via explicit mapping.
	agy, err := ResolveParticipant("antigravity-1", []Discovery{disco("agy", true)}, mapping)
	if err != nil || agy.ID != "antigravity-1" || agy.Adapter() != "agy" {
		t.Fatalf("antigravity-1 -> %+v err=%v", agy, err)
	}
}

func TestResolveFailsClosed(t *testing.T) {
	discovered := []Discovery{disco("claude", true)}
	// No exact id, no mapping -> hard error (never a prefix guess).
	if _, err := ResolveParticipant("claude-1", discovered, nil); err == nil {
		t.Error("expected error resolving claude-1 with no mapping")
	}
	// Mapping points at a family that is not installed/discovered -> error.
	if _, err := ResolveParticipant("kimi-1", discovered, map[string]string{"kimi-1": "kimi"}); err == nil {
		t.Error("expected error when mapped family is not discovered")
	}
	// Not-found discovery is ignored.
	if _, err := ResolveParticipant("codex", []Discovery{disco("codex", false)}, nil); err == nil {
		t.Error("expected error when the only matching agent is not installed")
	}
}

func TestResolveRejectsTraversal(t *testing.T) {
	discovered := []Discovery{disco("claude", true)}
	// Only genuine path-traversal ids are rejected (containment check).
	for _, bad := range []string{"../../tmp/x", "claude/1", "..", ".", "", "a..b", ".claude", "claude.", `a\b`} {
		if _, err := ResolveParticipant(bad, discovered, map[string]string{bad: "claude"}); err == nil {
			t.Errorf("ResolveParticipant(%q) must reject an unsafe participant id", bad)
		}
	}
	// Legacy safe ids (underscore, uppercase, single dot) still resolve (review MINOR).
	for _, ok := range []string{"my_cli", "Claude-1", "claude.v1"} {
		if _, err := ResolveParticipant(ok, []Discovery{disco(ok, true)}, nil); err != nil {
			t.Errorf("ResolveParticipant(%q) must accept a path-safe legacy id: %v", ok, err)
		}
	}
}

func TestResolvePreservesExplicitAdapter(t *testing.T) {
	// An exact match on a spec that already carries a distinct adapter must keep it.
	d := Discovery{Spec: Spec{ID: "claude-1", AdapterID: "claude"}, Found: true}
	got, err := ResolveParticipant("claude-1", []Discovery{d}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Adapter() != "claude" {
		t.Fatalf("adapter=%q, want claude (explicit adapter destroyed)", got.Adapter())
	}
}

func TestResolveDoesNotMutateInput(t *testing.T) {
	discovered := []Discovery{disco("claude", true)}
	if _, err := ResolveParticipant("claude-1", discovered, map[string]string{"claude-1": "claude"}); err != nil {
		t.Fatal(err)
	}
	// The source slice's spec identity must be untouched (resolver works on a copy).
	if discovered[0].ID != "claude" || discovered[0].AdapterID != "" {
		t.Fatalf("input mutated: %+v", discovered[0].Spec)
	}
}
