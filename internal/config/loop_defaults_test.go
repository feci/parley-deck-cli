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

// The seed template carries the [defaults.loop] block so `parley init` turns budgets on.
func TestCentralDefaultTemplateHasLoopBlock(t *testing.T) {
	tpl := centralDefaultTemplate()
	for _, want := range []string{"[defaults.loop]", "max_driver_steps", "max_wall_clock_ms", "max_cost_usd"} {
		if !strings.Contains(tpl, want) {
			t.Fatalf("seed template missing %q", want)
		}
	}
}
