//go:build !windows

package evidence

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/procctl"
)

// The preserved run-error branch exists for the hosted U2 signature: a
// non-ExitError run failure with zero captured output persisted only an
// unexplained exit_code of -1 with empty diagnostics, discarding the only
// text that names the failing branch. These tests pin the preservation, its
// privacy bounds, and — just as important — the cases it must NOT touch.

func TestRunErrorPreservedWhenOutputEmpty(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := RunCriterionControlled(ctx, root, "refused", "printf never > never-ran", "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
		if _, err := start(); err != nil {
			return err
		}
		// Refusing after start walks the same reaping path as the hosted U2
		// signature: runErr set with a live supervisor that ignores TERM, so
		// collection pays the full KillGroup grace with empty captured output.
		return errors.New("controller rejected before release")
	})
	if r.Complete || r.Record.Status != StatusFail || r.Record.Command.ExitCode != -1 {
		t.Fatalf("refused run misclassified: %+v", r)
	}
	if got := r.Record.Command.Diagnostics; !strings.Contains(got, "run error (no command output): controller rejected before release") {
		t.Fatalf("run error reason not preserved: %q", got)
	}
	if r.Record.Command.OutputSHA256 != sha256Hex([]byte{}) {
		t.Fatalf("preserved reason must not enter the output hash: %s", r.Record.Command.OutputSHA256)
	}
	if r.Record.Command.Format != FormatShell {
		t.Fatalf("empty output must keep shell format: %+v", r.Record.Command)
	}
	if _, err := os.Stat(filepath.Join(root, "never-ran")); !os.IsNotExist(err) {
		t.Fatal("refused material command ran")
	}
}

// A non-controlled start failure (absent working directory) exercises the
// same branch without a control callback.
func TestRunErrorPreservedOnStartFailure(t *testing.T) {
	rec := RunCriterion(context.Background(), filepath.Join(t.TempDir(), "absent"), "start-failure", "echo started", "reviewer")
	if rec.Status != StatusFail || rec.Command.ExitCode != -1 {
		t.Fatalf("failed start misclassified: %+v", rec)
	}
	if got := rec.Command.Diagnostics; !strings.Contains(got, "run error (no command output):") {
		t.Fatalf("start-failure reason not preserved: %q", got)
	}
}

// The preserved reason passes the same redaction as captured output:
// credential-shaped tokens never reach the envelope.
func TestRunErrorPreservedIsScrubbed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	secret := "ghp_" + strings.Repeat("A", 28)
	r := RunCriterionControlled(ctx, t.TempDir(), "secret", "printf never", "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
		return errors.New("controller refused: bearer aBcD.eFgH-token123 and api_key=" + secret)
	})
	if r.Complete || r.Record.Status != StatusFail || r.Record.Command.ExitCode != -1 {
		t.Fatalf("scrub probe misclassified: %+v", r)
	}
	d := r.Record.Command.Diagnostics
	if !strings.Contains(d, "«redacted»") || strings.Contains(d, secret) || strings.Contains(d, "aBcD.eFgH-token123") {
		t.Fatalf("run error reason not scrubbed: %q", d)
	}
	if !strings.Contains(d, "run error (no command output):") {
		t.Fatalf("label missing from preserved reason: %q", d)
	}
}

// The preserved reason is bounded like any other diagnostics payload.
func TestRunErrorPreservedIsBounded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := RunCriterionControlled(ctx, t.TempDir(), "flood", "printf never", "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
		return errors.New(strings.Repeat("x", 6000))
	})
	d := r.Record.Command.Diagnostics
	if len(d) > evidenceMaxBytes+3 {
		t.Fatalf("run error reason not bounded: %d bytes", len(d))
	}
	if !strings.Contains(d, "…") {
		t.Fatalf("truncation marker missing from bounded reason: %q", d)
	}
}

// When the run fails with a non-ExitError but the command DID emit output,
// diagnostics stay the captured output alone — the run error is not
// prepended, and the output hash covers exactly what the command emitted.
func TestRunErrorNotPersistedWhenOutputPresent(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r := RunCriterionControlled(ctx, root, "post-execution", "printf 'material-output\\n'; touch material-done", "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
		if _, err := start(); err != nil {
			return err
		}
		if err := release(); err != nil {
			return err
		}
		// Wait for the material command to COMPLETE before refusing, so the
		// refusal-triggered group kill cannot race the printf. The shell runs
		// printf before touch, so the bytes are already in the pipe once the
		// marker appears; cmd.Wait's drain-to-EOF then guarantees they reach
		// the capture buffer. (Polling process liveness here would be wrong:
		// the unreaped supervisor is a zombie, and a zombie answers a signal-0
		// probe as if alive.)
		for i := 0; i < 500; i++ {
			if _, err := os.Stat(filepath.Join(root, "material-done")); err == nil {
				return errors.New("controller refused after execution")
			}
			time.Sleep(10 * time.Millisecond)
		}
		return errors.New("material command never completed")
	})
	if r.Complete || r.Record.Status != StatusFail || r.Record.Command.ExitCode != -1 {
		t.Fatalf("post-execution refusal misclassified: %+v", r)
	}
	const out = "material-output\n"
	if got := r.Record.Command.Diagnostics; got != ScrubAndTruncate(out) || strings.Contains(got, "run error") {
		t.Fatalf("diagnostics must be the captured output only: %q", got)
	}
	if r.Record.Command.OutputSHA256 != sha256Hex([]byte(out)) {
		t.Fatalf("output hash must cover exactly the command output: %s", r.Record.Command.OutputSHA256)
	}
}

// An ExitError with empty output keeps the previous envelope shape — exit
// code plus empty diagnostics: the preservation branch is for non-ExitError
// run failures only.
func TestRunErrorBranchSkipsExitError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := RunCriterionControlled(ctx, t.TempDir(), "silent-exit", "exit 3", "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
		if _, err := start(); err != nil {
			return err
		}
		return release()
	})
	if !r.Complete || r.Record.Status != StatusFail || r.Record.Command.ExitCode != 3 {
		t.Fatalf("silent exit misclassified: %+v", r)
	}
	if r.Record.Command.Diagnostics != "" {
		t.Fatalf("ExitError with no output must keep empty diagnostics: %q", r.Record.Command.Diagnostics)
	}
}
