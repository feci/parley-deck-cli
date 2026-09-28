//go:build windows

package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// §D.3 W-SHELL hosted pin, checks path: with `sh` forced off PATH the LE-4
// verification gate blocks BEFORE any work with the named Git-for-Windows
// refusal — the gate never passes by default and exits non-zero (ok=false).
func TestRunChecksMissingShRefusesPreWork(t *testing.T) {
	dir := t.TempDir()
	idea := filepath.Join(dir, "ideas", "x")
	if err := os.MkdirAll(idea, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(idea, "00-prompt.md"),
		[]byte("---\nidea: x\nchecks: exit 0\nstatus: round-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", t.TempDir()) // no sh.exe anywhere on PATH
	o := driverImplOps{root: dir, ideaDir: idea, out: io.Discard}
	ok, msg := o.RunChecks(context.Background())
	if ok {
		t.Fatalf("checks passed without sh on PATH: %q", msg)
	}
	for _, want := range []string{"requires the POSIX shell", "Git for Windows", "refusing rather than executing unverified"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("refusal missing %q: %q", want, msg)
		}
	}
}
