//go:build windows

package evidence

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// §D.3 W-SHELL hosted pin, both outcomes: with `sh` forced off PATH the
// captured-verification criterion refuses BEFORE any work with the named,
// actionable message naming Git for Windows (never a green unverified run);
// with the runner's normal PATH (Git for Windows sh.exe present) the normal
// execution path proceeds.
func TestRunCriterionControlledMissingShRefusesPreWork(t *testing.T) {
	empty := t.TempDir()
	t.Setenv("PATH", empty) // no sh.exe anywhere on PATH
	start := time.Now()
	exec := RunCriterionControlled(context.Background(), t.TempDir(), "c", "exit 0", "executor-x", nil)
	if exec.Complete {
		t.Fatal("criterion completed without sh on PATH — unverified execution")
	}
	diag := exec.Record.Command.Diagnostics
	if !strings.Contains(diag, "requires the POSIX shell") || !strings.Contains(diag, "Git for Windows") ||
		!strings.Contains(diag, "refusing rather than executing unverified") {
		t.Fatalf("refusal must name the prerequisite and the refusal stance: %q", diag)
	}
	if exec.Record.Command.ExitCode != -1 {
		t.Fatalf("refused criterion must record a non-zero exit code, got %d", exec.Record.Command.ExitCode)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatalf("refusal must fire before any work: took %s", time.Since(start))
	}
	if _, err := os.Stat(filepath.Join(t.TempDir(), "anything")); err == nil {
		t.Fatal("refusal must not create work artifacts")
	}
}

func TestRunCriterionControlledWithShRunsNormally(t *testing.T) {
	// GitHub's windows-latest ships Git for Windows sh.exe on PATH; this leg
	// proves the normal path proceeds when the prerequisite is present. A
	// missing sh here is a reportable hosted fact, never a skip.
	if _, err := exec.LookPath("sh"); err != nil {
		t.Fatalf("hosted Windows image unexpectedly lacks sh on PATH: %v", err)
	}
	exec := RunCriterionControlled(context.Background(), t.TempDir(), "c", "exit 0", "executor-x", nil)
	if !exec.Complete || exec.Record.Command.ExitCode != 0 {
		t.Fatalf("normal path with sh present must run: complete=%v exit=%d diagnostics=%q",
			exec.Complete, exec.Record.Command.ExitCode, exec.Record.Command.Diagnostics)
	}
}
