package app

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMain isolates app tests from the developer's real ~/.parley/agents.toml
// (the PARLEY_HOME central-default layer) so agent-spec golden assertions see
// only built-in defaults.
func TestMain(m *testing.M) {
	// §D.9 test-binary re-exec: fixtures installed by writeReexecFixture are
	// copies of this binary (an extension-less #!/bin/sh script cannot be
	// exec'd on Windows); dispatch on the basename to the fake behavior
	// instead of entering the testing framework.
	if code, ok := runFixtureRole(filepath.Base(os.Args[0])); ok {
		os.Exit(code)
	}
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
