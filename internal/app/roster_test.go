package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
)

func TestProposeFamily(t *testing.T) {
	byFamily := map[string]agents.Spec{"claude": {}, "codex": {}, "agy": {}}
	cases := map[string]string{
		"claude":        "claude", // exact family
		"claude-1":      "claude", // trailing -N stripped
		"codex-2":       "codex",
		"antigravity-1": "agy", // alias (prefix would fail)
		"kimi-1":        "",    // no such family -> unresolved
	}
	for id, want := range cases {
		if got := proposeFamily(id, byFamily); got != want {
			t.Errorf("proposeFamily(%q) = %q, want %q", id, got, want)
		}
	}
}

func TestRosterShowAndInit(t *testing.T) {
	root := t.TempDir()
	t.Setenv(config.EnvParleyHome, filepath.Join(root, "no-central"))
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	// Give the deck a §2 roster whose IDs resolve to real families (claude-1->claude,
	// antigravity-1->agy via alias). The genericized init template has placeholders.
	coop := filepath.Join(root, protocol.DeckDir, "COOPERATION.md")
	if err := os.WriteFile(coop, []byte(`# Cooperation

## 2. Active agents (roster)

| Agent ID | Workspace dir | Role |
| -------- | ------------- | ---- |
|`+" `claude-1` "+`| ../claude/ | participant |
|`+" `codex-1` "+`| ../codex/ | participant |
|`+" `antigravity-1` "+`| ../antigravity/ | participant |
`), 0o644); err != nil {
		t.Fatal(err)
	}

	// show must run and render every §2 roster id with a family_model_effort name.
	var out, errb bytes.Buffer
	if code := runRoster([]string{"show", "--dir", root}, &out, &errb); code != 0 {
		t.Fatalf("roster show exit=%d stderr=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "DISPLAY-NAME") {
		t.Fatalf("show output missing header:\n%s", out.String())
	}

	// A dry-run proposes the mapping and writes nothing.
	out.Reset()
	errb.Reset()
	if code := runRoster([]string{"init", "--dir", root, "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("roster init --dry-run exit=%d stderr=%s", code, errb.String())
	}
	if m, _ := config.LoadRosterAdapters(root); len(m) != 0 {
		t.Fatalf("dry-run wrote a mapping: %v", m)
	}

	// A real init writes the [roster.*] mapping for every resolvable id.
	out.Reset()
	errb.Reset()
	if code := runRoster([]string{"init", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("roster init --yes exit=%d stderr=%s", code, errb.String())
	}
	mapping, err := config.LoadRosterAdapters(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(mapping) == 0 {
		t.Fatalf("init wrote no mapping; output=%s", out.String())
	}
	// Every written mapping must round-trip through the resolver (identity = roster
	// id, adapter = family). Build discovered from specs with Found=true so the test
	// does not depend on which CLIs happen to be installed.
	specs, _ := config.LoadAgentSpecs(root)
	discovered := make([]agents.Discovery, 0, len(specs))
	for _, s := range specs {
		discovered = append(discovered, agents.Discovery{Spec: s, Found: true})
	}
	for id, family := range mapping {
		got, rerr := agents.ResolveParticipant(id, discovered, mapping)
		if rerr != nil {
			t.Errorf("resolve %q failed: %v", id, rerr)
			continue
		}
		if got.ID != id || got.Adapter() != family {
			t.Errorf("resolve %q -> id=%q adapter=%q, want id=%q adapter=%q", id, got.ID, got.Adapter(), id, family)
		}
	}

	// Re-running is idempotent (no duplicate blocks).
	out.Reset()
	errb.Reset()
	if code := runRoster([]string{"init", "--dir", root, "--yes"}, &out, &errb); code != 0 {
		t.Fatalf("second init exit=%d stderr=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "already initialized") {
		t.Fatalf("second init not idempotent:\n%s", out.String())
	}
}
