package agents

import "strings"

import "testing"

// FIX-PROVING. Before {model} substitution existed, the claude built-in carried
// "claude-opus-4-8[1m]" as a literal inside HeadlessArgs while Model was a separate field.
// Config layers set the field and never rewrote the args, so pinning a model changed only
// what was DISPLAYED — the process kept launching the built-in literal. This test fails on
// any spec that reintroduces a model literal in the args.
func TestBuiltinSpecsCarryNoModelLiteralInArgs(t *testing.T) {
	for _, s := range DefaultSpecs() {
		for i, a := range s.HeadlessArgs {
			if a != "--model" && a != "-m" {
				continue
			}
			if i+1 >= len(s.HeadlessArgs) {
				t.Errorf("%s: model flag %q is the last token, so it has no value", s.ID, a)
				continue
			}
			if v := s.HeadlessArgs[i+1]; v != ModelPlaceholder {
				t.Errorf("%s: model flag %q takes literal %q; it must take %s so config can reach the launch",
					s.ID, a, v, ModelPlaceholder)
			}
		}
	}
}

func TestResolveLaunchArgsSubstitutesConfiguredValues(t *testing.T) {
	s := Spec{
		HeadlessArgs: []string{"-p", "--model", ModelPlaceholder, "--effort", EffortPlaceholder, "--add-dir", "{root}"},
		Model:        "claude-opus-5[1m]",
		Reasoning:    "max",
	}
	got, status := s.ResolveLaunchArgs()
	want := "-p --model claude-opus-5[1m] --effort max --add-dir {root}"
	if strings.Join(got, " ") != want {
		t.Fatalf("resolved=%q want %q", strings.Join(got, " "), want)
	}
	if status.ModelUnbound || status.EffortUnbound {
		t.Fatalf("nothing should be unbound: %+v", status)
	}
	if m, ok := s.EffectiveModel(); !ok || m != "claude-opus-5[1m]" {
		t.Fatalf("effective model=%q ok=%v", m, ok)
	}
	if e, ok := s.EffectiveEffort(); !ok || e != "max" {
		t.Fatalf("effective effort=%q ok=%v", e, ok)
	}
}

// An unbindable placeholder must take its introducing flag with it. Leaving a value-taking
// flag as the last token makes the vendor CLI abort with a parse error and write nothing —
// a silent launch failure, which is worse than falling back to the CLI's own default.
func TestResolveLaunchArgsDropsFlagAndPlaceholderWhenUnbound(t *testing.T) {
	for _, model := range []string{"", CLIDefault, "   "} {
		s := Spec{
			HeadlessArgs: []string{"exec", "--sandbox", "workspace-write", "-m", ModelPlaceholder, "-"},
			Model:        model,
		}
		got, status := s.ResolveLaunchArgs()
		want := "exec --sandbox workspace-write -"
		if strings.Join(got, " ") != want {
			t.Fatalf("model=%q resolved=%q want %q", model, strings.Join(got, " "), want)
		}
		if !status.ModelUnbound {
			t.Fatalf("model=%q should be unbound", model)
		}
		if _, ok := s.EffectiveModel(); ok {
			t.Fatalf("model=%q must not report an effective model", model)
		}
	}
}

// A positional token before an unbound placeholder is NOT a flag and must survive.
func TestResolveLaunchArgsKeepsPositionalBeforeUnboundPlaceholder(t *testing.T) {
	s := Spec{HeadlessArgs: []string{"run", ModelPlaceholder, "{prompt}"}}
	got, _ := s.ResolveLaunchArgs()
	if strings.Join(got, " ") != "run {prompt}" {
		t.Fatalf("resolved=%q want %q", strings.Join(got, " "), "run {prompt}")
	}
}

// AUTO must be computed from the resolved argv. A spec whose autonomous flag survives
// resolution stays yes; the check must not be fooled by placeholder removal.
func TestAutonomousEffectiveUsesResolvedArgs(t *testing.T) {
	s := Spec{
		HeadlessArgs:    []string{"-p", "--model", ModelPlaceholder, "--permission-mode", "bypassPermissions"},
		AutonomousWrite: AutonomousWrite{Mode: "bypassPermissions", Args: []string{"--permission-mode", "bypassPermissions"}},
	}
	if !s.AutonomousEffective() {
		t.Fatal("autonomous args survive resolution, so AUTO must be yes")
	}
	stripped := s
	stripped.HeadlessArgs = []string{"-p", "--model", ModelPlaceholder}
	if stripped.AutonomousEffective() {
		t.Fatal("declared mode whose enabling args are absent must report AUTO=no")
	}
}

// FIX-PROVING (idea zcode-adapter). zcode has NO model flag: `--model` is absent from
// `zcode --help` and passing it exits 1. Its argv must therefore carry no model or effort
// placeholder at all — a future edit that adds one would make config appear to bind a value
// the process can never receive, which is the exact defect TestBuiltinSpecsCarryNoModelLiteralInArgs
// guards against in the opposite direction. The exact argv is locked so a spec change cannot
// silently append an option zcode rejects.
func TestZcodeSpecCarriesNoModelOrEffortPlaceholder(t *testing.T) {
	var zcode Spec
	for _, s := range DefaultSpecs() {
		if s.ID == "zcode" {
			zcode = s
			break
		}
	}
	if zcode.ID == "" {
		t.Fatal("built-in spec zcode is missing")
	}
	for _, a := range zcode.HeadlessArgs {
		if a == ModelPlaceholder {
			t.Errorf("zcode HeadlessArgs carries %s; zcode has no model flag, so nothing can bind it", ModelPlaceholder)
		}
		if a == EffortPlaceholder {
			t.Errorf("zcode HeadlessArgs carries %s; zcode has no effort flag, so nothing can bind it", EffortPlaceholder)
		}
		if a == "--model" || a == "-m" {
			t.Errorf("zcode HeadlessArgs carries %q; `zcode --model` exits 1", a)
		}
	}
	if got, want := strings.Join(zcode.HeadlessArgs, " "), "--prompt {prompt} --mode yolo --cwd {root}"; got != want {
		t.Fatalf("zcode argv=%q want %q", got, want)
	}
	// The autonomous-write flag must be present in the argv it claims to be enabled by.
	if missing := zcode.AutonomousWrite.MissingFrom(zcode.HeadlessArgs); len(missing) > 0 {
		t.Fatalf("zcode declares mode %q but its argv is missing %v", zcode.AutonomousWrite.Mode, missing)
	}
	// --cwd is a working directory, not an enforced sandbox: Scope must stay empty.
	if zcode.AutonomousWrite.Confined() {
		t.Error("zcode claims workspace confinement; --cwd is a cwd, not a sandbox")
	}
}
