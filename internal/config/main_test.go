package config

import (
	"os"
	"testing"
)

// TestMain isolates this package's tests from the developer's real
// ~/.parley/agents.toml: unless a test sets PARLEY_HOME itself, point it at an
// empty scratch dir so the central-default layer is a no-op and built-in
// default assertions stay stable.
func TestMain(m *testing.M) {
	os.Exit(runIsolated(m))
}

func runIsolated(m *testing.M) int {
	if _, ok := os.LookupEnv(EnvParleyHome); !ok {
		if dir, err := os.MkdirTemp("", "parley-no-central"); err == nil {
			os.Setenv(EnvParleyHome, dir)
			defer os.RemoveAll(dir)
		}
	}
	return m.Run()
}
