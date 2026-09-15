package evidence

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/procctl"
)

func TestCriterionControlledPublicationBeforeExecution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("controlled criterion execution requires POSIX")
	}
	for _, mode := range []string{"reject-before-start", "reject-publication", "missing-release", "execute"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			var sp procctl.Spawned
			command := `printf started > material-started; printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; exit 0`
			r := RunCriterionControlled(ctx, root, "material", command, "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
				if mode == "reject-before-start" {
					return errors.New("stopped")
				}
				var err error
				sp, err = start()
				if err != nil {
					return err
				}
				if _, err = os.Stat(filepath.Join(root, "material-started")); !os.IsNotExist(err) {
					t.Error("material executed before process publication")
				}
				if strings.Contains(sp.Command, "material-started") || strings.Contains(sp.Command, "executed_cases") {
					t.Error("registered supervisor identity leaked the raw command")
				}
				if mode == "reject-publication" {
					return errors.New("identity write failed")
				}
				if mode == "missing-release" {
					return nil
				}
				return release()
			})
			if mode == "execute" {
				if !r.Complete || r.Record.Status != StatusPass || r.Record.Command.ExecutedCases != 1 || r.Record.Command.CommandSHA256 != sha256Hex([]byte(command)) {
					t.Fatalf("controlled command lost original semantics: %+v", r)
				}
			} else {
				if r.Complete || r.Record.Status != StatusFail {
					t.Fatalf("refused command accepted: %+v", r)
				}
				if _, err := os.Stat(filepath.Join(root, "material-started")); !os.IsNotExist(err) {
					t.Fatal("refused material command ran")
				}
			}
			if procctl.Alive(sp) {
				t.Fatal("controlled executor leaked its owned supervisor")
			}
		})
	}
}

func TestCriterionControlledSignalIsIncomplete(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX criterion execution")
	}
	for _, command := range []string{
		`printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; kill -KILL "$$"`,
		`printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 137`,
	} {
		r := RunCriterionControlled(context.Background(), t.TempDir(), "signal", command, "reviewer", func(start func() (procctl.Spawned, error), release func() error) error {
			if _, err := start(); err != nil {
				return err
			}
			return release()
		})
		if r.Complete || r.Record.Status != StatusFail {
			t.Fatalf("signal-like material exit certified as complete: %+v", r)
		}
	}
}
