package agents

import (
	"strings"
	"testing"
)

// Component C (composite-agent-naming-and-roster-reinit): the built-in specs must
// ship an actually-autonomous write mode. kimi's round-01 found claude shipped
// acceptEdits and codex on-failure — neither is fully non-interactive.
func TestBuiltinsAreAutonomous(t *testing.T) {
	byID := map[string]Spec{}
	for _, s := range DefaultSpecs() {
		byID[s.ID] = s
	}
	wantMode := map[string]string{
		"claude": "bypassPermissions",
		"codex":  "approval_policy=never",
		"agy":    "dangerously-skip-permissions",
		"hermes": "yolo",
		// Promoted from ACP-only stubs to full adapters (2026-08-06). Both were probed
		// live before the mode names below were written: `kimi -p` and
		// `opencode run --auto` each wrote their file unattended, exit 0.
		"kimi":     "prompt",
		"opencode": "auto",
	}
	for id, mode := range wantMode {
		s, ok := byID[id]
		if !ok {
			t.Fatalf("missing built-in %q", id)
		}
		if !s.AutonomousWrite.Declared() {
			t.Errorf("%s: autonomous write not declared (%+v)", id, s.AutonomousWrite)
		}
		if s.AutonomousWrite.Mode != mode {
			t.Errorf("%s: mode=%q want %q", id, s.AutonomousWrite.Mode, mode)
		}
	}
	// The non-autonomous legacy defaults must be gone / replaced.
	joined := func(s Spec) string { return strings.Join(s.HeadlessArgs, " ") }
	if strings.Contains(joined(byID["claude"]), "acceptEdits") {
		t.Error("claude still ships the non-autonomous acceptEdits mode")
	}
	if strings.Contains(joined(byID["codex"]), "on-failure") {
		t.Error("codex still ships the non-autonomous on-failure policy")
	}
	if !strings.Contains(joined(byID["hermes"]), "--yolo") {
		t.Error("hermes is missing --yolo")
	}
}
