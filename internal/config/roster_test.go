package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeDeckConfig(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadRosterPresetsDeckOnly(t *testing.T) {
	root := t.TempDir()
	writeDeckConfig(t, root, "parley-deck/agents.toml", `
[rosters.pair]
participants = ["claude-1", "codex-1"]
[rosters.council]
participants = ["claude-1", "codex-1", "hermes-1", "antigravity-1"]
[defaults.track_rosters]
fast = "pair"
deliberation = "council"
`)
	rc, err := LoadRosterPresets(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(rc.Presets) != 2 {
		t.Fatalf("want 2 presets, got %d", len(rc.Presets))
	}
	if got := rc.Presets["pair"].Participants; len(got) != 2 || got[0] != "claude-1" {
		t.Fatalf("pair participants = %v", got)
	}
	if rc.TrackDefault["fast"] != "pair" || rc.TrackDefault["deliberation"] != "council" {
		t.Fatalf("track defaults = %v", rc.TrackDefault)
	}
}

func TestResolveRosterExplicitPreset(t *testing.T) {
	rc := RosterConfig{
		Presets: map[string]RosterPreset{
			"pair": {Name: "pair", Participants: []string{"claude-1", "codex-1"}, Source: "deck"},
		},
		TrackDefault: map[string]string{},
	}
	roster := map[string]bool{"claude-1": true, "codex-1": true, "hermes-1": true}
	res, err := ResolveRoster(rc, "pair", "", roster, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Participants) != 2 || res.Preset != "pair" {
		t.Fatalf("res = %+v", res)
	}
	if res.Provenance == "" || res.Provenance[:4] != "<!--" {
		t.Fatalf("provenance = %q", res.Provenance)
	}
}

func TestResolveRosterTrackDefault(t *testing.T) {
	rc := RosterConfig{
		Presets:      map[string]RosterPreset{"council": {Name: "council", Participants: []string{"claude-1", "codex-1"}, Source: "central"}},
		TrackDefault: map[string]string{"deliberation": "council"},
	}
	roster := map[string]bool{"claude-1": true, "codex-1": true}
	res, err := ResolveRoster(rc, "", "deliberation", roster, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Preset != "council" {
		t.Fatalf("track default not applied: %+v", res)
	}
	// No track default for standard → no preset.
	res2, err := ResolveRoster(rc, "", "standard", roster, nil)
	if err != nil || len(res2.Participants) != 0 {
		t.Fatalf("standard should yield no preset: %+v err=%v", res2, err)
	}
}

func TestResolveRosterFailClosed(t *testing.T) {
	rc := RosterConfig{
		Presets: map[string]RosterPreset{
			"pair":   {Name: "pair", Participants: []string{"claude-1", "codex-1"}, Source: "deck"},
			"empty":  {Name: "empty", Participants: nil, Source: "deck"},
			"dupe":   {Name: "dupe", Participants: []string{"claude-1", "claude-1"}, Source: "deck"},
			"ghost":  {Name: "ghost", Participants: []string{"claude-1", "gemini-1"}, Source: "deck"},
			"asleep": {Name: "asleep", Participants: []string{"claude-1", "codex-1"}, Source: "deck"},
		},
		TrackDefault: map[string]string{},
	}
	roster := map[string]bool{"claude-1": true, "codex-1": true}
	inactive := map[string]bool{"codex-1": true}

	if _, err := ResolveRoster(rc, "nope", "", roster, nil); err == nil {
		t.Fatal("unknown preset should error")
	}
	if _, err := ResolveRoster(rc, "empty", "", roster, nil); err == nil {
		t.Fatal("empty preset should error")
	}
	if _, err := ResolveRoster(rc, "dupe", "", roster, nil); err == nil {
		t.Fatal("duplicate member should error")
	}
	if _, err := ResolveRoster(rc, "ghost", "", roster, nil); err == nil {
		t.Fatal("member not in §2 should error")
	}
	if _, err := ResolveRoster(rc, "asleep", "", roster, inactive); err == nil {
		t.Fatal("inactive member should error")
	}
}
