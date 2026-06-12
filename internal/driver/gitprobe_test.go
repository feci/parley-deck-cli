package driver

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestGitTreeCleanSetsOptionalLocksOff (consensus D8): every git probe the
// driver spawns must run with GIT_OPTIONAL_LOCKS=0 so read-only status checks
// never write .git on a weakly-coherent mount. A PATH-shimmed fake git records
// the env it saw.
func TestGitTreeCleanSetsOptionalLocksOff(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PATH shim uses a shell script")
	}
	dir := t.TempDir()
	record := filepath.Join(dir, "seen-env")
	shim := filepath.Join(dir, "git")
	script := "#!/bin/sh\nprintf '%s\\n' \"$GIT_OPTIONAL_LOCKS\" >> " + record + "\nexit 0\n"
	if err := os.WriteFile(shim, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	_ = gitTreeClean(dir)

	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("the shim never ran: %v", err)
	}
	for i, line := range splitNonEmptyLines(string(data)) {
		if line != "0" {
			t.Fatalf("probe %d ran with GIT_OPTIONAL_LOCKS=%q, want 0", i, line)
		}
	}
}

func splitNonEmptyLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}
