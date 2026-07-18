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

func TestAppLevelRosterIDResolution(t *testing.T) {
	// The deeper app-level paths must resolve a roster id (codex-1) via the mapping,
	// not just the runner (review MAJOR #5: preflight/drafter/signoff/steer/selection).
	discovered := []agents.Discovery{
		{Spec: agents.Spec{ID: "codex", LaunchMode: agents.LaunchHeadless}, Found: true},
	}
	mapping := map[string]string{"codex-1": "codex"}

	if sel, err := selectedParticipantIDs(discovered, "codex-1", mapping); err != nil || len(sel) != 1 || sel[0] != "codex-1" {
		t.Fatalf("selectedParticipantIDs(codex-1) = %v, %v", sel, err)
	}
	if pd := participantDiscoveries(discovered, []string{"codex-1"}, mapping); len(pd) != 1 || pd[0].ID != "codex-1" || pd[0].Adapter() != "codex" {
		t.Fatalf("participantDiscoveries(codex-1) = %+v", pd)
	}
	if d, ok := firstHeadlessAgent(discovered, []string{"codex-1"}, mapping); !ok || d.ID != "codex-1" || d.Adapter() != "codex" {
		t.Fatalf("firstHeadlessAgent(codex-1) = %+v, ok=%v", d.Spec, ok)
	}
	if sel, err := requestSignoffAgents([]string{"codex-1"}, discovered, mapping); err != nil || len(sel) != 1 || sel[0].ID != "codex-1" {
		t.Fatalf("requestSignoffAgents(codex-1) = %v, %v", sel, err)
	}
	// Without the mapping, a roster id is fail-closed (not silently a family).
	if _, err := selectedParticipantIDs(discovered, "codex-1", nil); err == nil {
		t.Fatal("selectedParticipantIDs(codex-1, nil) should fail closed")
	}
}

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

func setupRosterDeck(t *testing.T) string {
	t.Helper()
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
	return root
}

func TestRosterInitJSONAndScope(t *testing.T) {
	root := setupRosterDeck(t)
	// Invalid scope fails with exit 2 and writes nothing (review MAJOR).
	var out, errb bytes.Buffer
	if code := runRoster([]string{"init", "--dir", root, "--scope", "machien"}, &out, &errb); code != 2 {
		t.Fatalf("bad --scope should exit 2, got %d", code)
	}
	// --json --yes must ACTUALLY write and report the real outcome (review MAJOR: no-op bug).
	out.Reset()
	errb.Reset()
	if code := runRoster([]string{"init", "--dir", root, "--yes", "--json"}, &out, &errb); code != 0 {
		t.Fatalf("json init exit=%d stderr=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), `"outcome": "written"`) {
		t.Fatalf("json outcome not 'written':\n%s", out.String())
	}
	if m, _ := config.LoadRosterAdapters(root); len(m) == 0 {
		t.Fatal("--json --yes did not persist the mapping")
	}
}

func TestRosterInitRejectsInvalidMapping(t *testing.T) {
	root := setupRosterDeck(t)
	// A mapping to a non-existent family must NOT read as initialized (review MAJOR).
	deckCfg := filepath.Join(root, protocol.DeckDir, "agents.toml")
	if err := os.WriteFile(deckCfg, []byte("[roster.claude-1]\nadapter = \"claud\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	if code := runRoster([]string{"init", "--dir", root, "--yes"}, &out, &errb); code == 0 {
		t.Fatalf("init with a typoed adapter should fail, got exit 0:\n%s / %s", out.String(), errb.String())
	}
}

func TestRosterShowAndInit(t *testing.T) {
	root := setupRosterDeck(t)

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
