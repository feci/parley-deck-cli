package driver

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/trajectory"
)

func TestTrajectoryDriverReservationPersistsBeforeAdapterAndBlocksResume(t *testing.T) {
	ideaDir, runDir, parts := setupReviewPhase(t, "auto_implement: true\n")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root, "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid"}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Git: %v %s", err, out)
		}
	}
	// The driver fixtures own this repository; ignore their runtime cursor only.
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".parley-runtime/\nparley-deck/runs/\n"), 0600); err != nil {
		t.Fatal(err)
	}
	git("init", "-q")
	git("add", ".")
	git("commit", "-qm", "Trajectory driver baseline")
	d := newImplDriver(ideaDir, runDir, parts, true, &fakeImpl{})
	b, err := budget.EnsureCycleBinding(context.Background(), root, "demo", budget.Fixup, d.cfg.MaxFixupCycles, 0, runDir, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	p, expected, err := trajectory.NewPolicy(context.Background(), root, "demo", "codex", []trajectory.Criterion{{Name: "material", Command: "true"}})
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	ctx, finishStep := budget.GroupStepSession(context.Background(), root, "demo")
	defer finishStep()
	_, n, finish, err := d.reserveFixupCycle(ctx, 0)
	finish()
	if err != nil || n != 1 {
		t.Fatalf("driver did not precharge observed attempt: %d %v", n, err)
	}
	s, err := trajectory.Inspect(context.Background(), root, "demo")
	if err != nil || s == nil || len(s.Attempts) != 1 || s.Attempts[0].Launch != nil {
		t.Fatalf("pre-adapter interruption not retained: %+v %v", s, err)
	}
	// A fresh driver cannot infer zero history from its new run/cursor.
	next := newImplDriver(ideaDir, filepath.Join(filepath.Dir(runDir), "resumed"), parts, true, &fakeImpl{})
	nextCtx, finishNext := budget.GroupStepSession(context.Background(), root, "demo")
	defer finishNext()
	_, _, finish, err = next.reserveFixupCycle(nextCtx, 1)
	finish()
	if err == nil || !strings.Contains(err.Error(), "trajectory") {
		t.Fatalf("resumed driver admitted unverified patch: %v", err)
	}
	ledger, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(ledger) != 1 {
		t.Fatalf("driver retry changed charges: %+v %v", ledger, err)
	}
	if err = os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	fi := &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, ReviewerCount: 2}}
	closing := newImplDriver(ideaDir, filepath.Join(filepath.Dir(runDir), "closing"), parts, true, fi)
	action, _, err := closing.Advance(context.Background())
	if err == nil || !strings.Contains(err.Error(), "trajectory") || action != ActionEscalated || contains(fi.calls, "complete") {
		t.Fatalf("driver completion ignored pending trajectory: action=%s calls=%v error=%v", action, fi.calls, err)
	}

}
