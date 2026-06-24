package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/protocol"
)

// LE-5: [defaults.loop] parses and merges into CentralDefaults.
func TestLoadDefaultsLoopBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)
	if err := os.WriteFile(filepath.Join(home, "agents.toml"), []byte(`
[defaults.loop]
max_driver_steps = 50
max_wall_clock_ms = 600000
max_cost_usd = 12.5
`), 0o644); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	defs, err := LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.MaxDriverSteps != 50 || defs.MaxWallClockMS != 600000 || defs.MaxCostUSD != 12.5 {
		t.Fatalf("[defaults.loop] not loaded: %+v", defs)
	}
}

// F-T2-1: a deck's explicit `= 0` overrides a central seed to unlimited; an unset deck
// field falls through to the central seed (presence-aware merge).
func TestLoopDefaultsDeckZeroOverridesCentralSeed(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvParleyHome, home)
	if err := os.WriteFile(filepath.Join(home, "agents.toml"),
		[]byte("[defaults.loop]\nmax_driver_steps = 200\nmax_wall_clock_ms = 7200000\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, protocol.DeckDir, "agents.toml"),
		[]byte("[defaults.loop]\nmax_driver_steps = 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	defs, err := LoadDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if defs.MaxDriverSteps != 0 {
		t.Fatalf("deck explicit 0 must override the central seed to unlimited; got %d", defs.MaxDriverSteps)
	}
	if defs.MaxWallClockMS != 7200000 {
		t.Fatalf("an unset deck wall-clock must fall through to the central seed; got %d", defs.MaxWallClockMS)
	}
}

// The seed template carries the [defaults.loop] block so `parley init` turns budgets on.
func TestCentralDefaultTemplateHasLoopBlock(t *testing.T) {
	tpl := centralDefaultTemplate()
	for _, want := range []string{"[defaults.loop]", "max_driver_steps", "max_wall_clock_ms", "max_cost_usd"} {
		if !strings.Contains(tpl, want) {
			t.Fatalf("seed template missing %q", want)
		}
	}
}
