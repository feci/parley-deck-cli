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

// AF-2 (review MINOR, codex-1): the Mode-only contract above would still pass if
// Args, HeadlessArgs, PromptMode or Scope regressed — the fields that make an
// adapter usable and truthful. Lock them for the two promoted adapters.
func TestPromotedAdaptersFullContract(t *testing.T) {
	byID := map[string]Spec{}
	for _, s := range DefaultSpecs() {
		byID[s.ID] = s
	}
	cases := []struct {
		id       string
		autoArgs []string
		headless []string
	}{
		{"kimi", []string{"-p"}, []string{"-m", "{model}", "-p", "{prompt}"}},
		{"opencode", []string{"--auto"}, []string{"run", "--auto", "-m", "{model}", "{prompt}"}},
	}
	for _, c := range cases {
		s, ok := byID[c.id]
		if !ok {
			t.Fatalf("missing built-in %q", c.id)
		}
		if got := strings.Join(s.AutonomousWrite.Args, " "); got != strings.Join(c.autoArgs, " ") {
			t.Errorf("%s: autonomous args=%v want %v", c.id, s.AutonomousWrite.Args, c.autoArgs)
		}
		if got := strings.Join(s.HeadlessArgs, " "); got != strings.Join(c.headless, " ") {
			t.Errorf("%s: headless args=%v want %v", c.id, s.HeadlessArgs, c.headless)
		}
		if s.PromptMode != PromptArg {
			t.Errorf("%s: prompt mode=%q want %q", c.id, s.PromptMode, PromptArg)
		}
		if s.LaunchMode != LaunchHeadless {
			t.Errorf("%s: launch mode=%q want %q", c.id, s.LaunchMode, LaunchHeadless)
		}
		// Scope must stay EMPTY: neither CLI enforces a workspace sandbox, and the
		// type forbids claiming confinement that is not enforced.
		if s.AutonomousWrite.Scope != "" {
			t.Errorf("%s: scope=%q want empty — only an enforced sandbox may claim one", c.id, s.AutonomousWrite.Scope)
		}
		if s.AutonomousWrite.Confined() {
			t.Errorf("%s: Confined() must be false", c.id)
		}
		// The declared mode must be enabled by the spec's own launch argv.
		if !s.AutonomousEffective() {
			t.Errorf("%s: declared mode %q is not enabled by its own HeadlessArgs %v",
				c.id, s.AutonomousWrite.Mode, s.HeadlessArgs)
		}
	}
}

// AF-1 (review MAJOR, codex-1): a config layer may replace HeadlessArgs without
// touching AutonomousWrite. The AUTO signal must then fail CLOSED. This is the
// regression that was live in production for hermes (--yolo stripped by an
// override) and for opencode (--auto stripped) while both reported AUTO=yes.
func TestAutonomousFailsClosedWhenLaunchDropsTheFlag(t *testing.T) {
	s := Spec{
		HeadlessArgs:    []string{"run", "-m", "some/model", "{prompt}"},
		AutonomousWrite: AutonomousWrite{Mode: "auto", Args: []string{"--auto"}},
	}
	if s.AutonomousWrite.Declared() != true {
		t.Fatal("precondition: the mode is declared")
	}
	if s.AutonomousEffective() {
		t.Error("AUTO must be false when the declared enabling arg is absent from the launch")
	}
	if got := s.AutonomousWrite.MissingFrom(s.HeadlessArgs); len(got) != 1 || got[0] != "--auto" {
		t.Errorf("MissingFrom=%v want [--auto]", got)
	}
	s.HeadlessArgs = []string{"run", "--auto", "{prompt}"}
	if !s.AutonomousEffective() {
		t.Error("AUTO must be true once the launch passes the declared arg")
	}
	if got := s.AutonomousWrite.MissingFrom(s.HeadlessArgs); len(got) != 0 {
		t.Errorf("MissingFrom=%v want empty", got)
	}
}
