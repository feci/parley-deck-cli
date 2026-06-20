package app

import (
	"os"
	"testing"
)

// TestMain isolates app tests from the developer's real ~/.parley/agents.toml
// (the PARLEY_HOME central-default layer) so agent-spec golden assertions see
// only built-in defaults.
func TestMain(m *testing.M) {
	os.Exit(runIsolated(m))
}

func runIsolated(m *testing.M) int {
	const envParleyHome = "PARLEY_HOME"
	if _, ok := os.LookupEnv(envParleyHome); !ok {
		if dir, err := os.MkdirTemp("", "parley-no-central"); err == nil {
			os.Setenv(envParleyHome, dir)
			defer os.RemoveAll(dir)
		}
	}
	return m.Run()
}
