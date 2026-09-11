package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/trajectory"
)

func trajectoryRuntimeFixture(t *testing.T) (string, protocol.IdeaStatus) {
	t.Helper()
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		base := []string{"-C", root, "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid"}
		cmd := exec.Command("git", append(base, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Git: %v %s", err, out)
		}
	}
	git("init", "-q")
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	declareTestLaunchSource(t, root)
	idea, err := protocol.CreateIdea(root, "Trajectory runtime fixture", []string{"test-1", "reviewer"})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(idea.Path, "00-prompt.md"), "---\nidea: "+idea.Slug+"\nparticipants: [test-1, reviewer]\ntrack: deliberation\nstatus: implemented\n---\nFixture\n")
	mustWrite(t, filepath.Join(idea.Path, "IMPLEMENTATION.md"), "---\nidea: "+idea.Slug+"\nstatus: implemented\n---\n\n## Summary of work\nFixture\n")
	mustWrite(t, filepath.Join(root, ".gitignore"), ".parley-runtime/\nparley-deck/runs/\n")
	mustWrite(t, filepath.Join(root, "source"), "original\n")
	git("add", ".")
	git("commit", "-qm", "Runtime fixture")
	_, err = budget.EnsureCycleBinding(context.Background(), root, idea.Slug, budget.Fixup, 5, 0, "", idea.Path)
	if err != nil {
		t.Fatal(err)
	}
	p, expected, err := trajectory.NewPolicy(context.Background(), root, idea.Slug, "test-1", []trajectory.Criterion{{Name: "material", Command: "true"}})
	if err != nil {
		t.Fatal(err)
	}
	if err = trajectory.Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	return root, idea
}
func TestTrajectoryRuntimeManualAndGroupedFixupCaptureActualFailure(t *testing.T) {
	for _, mode := range []string{"manual", "grouped"} {
		t.Run(mode, func(t *testing.T) {
			root, idea := trajectoryRuntimeFixture(t)
			agent := telemetryShell("printf 'broken\\n' > source; exit 7", false)
			agent.Found = true
			if mode == "manual" {
				r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: agent, Prompt: "synthetic fixture", Timeout: 5 * time.Second, Info: LaunchInfo{Idea: idea.Slug, RunID: "fixture-1", Phase: "fixup"}})
				if err == nil || r.StartedAt == nil {
					t.Fatalf("actual failing model fixture did not start: %+v %v", r, err)
				}
			} else {
				idea.Participants = []string{"test-1"}
				r := RunFixup(context.Background(), Options{Root: root, Idea: idea, RunID: "fixture-1", Agents: []agents.Discovery{agent}, Timeout: 5 * time.Second, Store: store.New(filepath.Join(root, "parley-deck", "runs", "fixture-1"))})
				// Existing artifact-wins semantics may accept the phase artifact, but that
				// cannot erase the underlying process failure in the trajectory.
				if r.Duration == 0 {
					t.Fatalf("grouped fixture never ran: %+v", r)
				}
			}
			s, err := trajectory.Inspect(context.Background(), root, idea.Slug)
			if err != nil {
				t.Fatal(err)
			}
			if s == nil || len(s.Attempts) != 1 {
				t.Fatalf("missing attempt: %+v", s)
			}
			a := s.Attempts[0]
			if a.Launch == nil || a.Terminal == nil || a.Terminal.Status != "failed" || a.Terminal.ExitCode == nil || *a.Terminal.ExitCode != 7 || a.After == nil || a.After.Clean {
				t.Fatalf("actual failed patch not retained: %+v", a)
			}
			restoreRuntimeAttempt(t, root, idea.Slug, a, "broken\n")
			records := terminalRecords(t, root)
			if len(records) != 1 || records[0].InvocationID != a.Launch.InvocationID {
				t.Fatalf("trajectory not tied to actual instrumented invocation: %+v", records)
			}
			r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: agent, Prompt: "retry fixture", Timeout: 5 * time.Second, Info: LaunchInfo{Idea: idea.Slug, RunID: "fixture-2", Phase: "fixup"}})
			if err == nil || r.StartedAt != nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "budget_refused" {
				t.Fatalf("unverified retry started: %+v %v", r, err)
			}
			b, err := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.Fixup)
			if err != nil {
				t.Fatal(err)
			}
			ledger, err := b.Store.Inspect(context.Background())
			if err != nil || b.Count(ledger) != 1 {
				t.Fatalf("retry changed charge count: %v %+v", err, ledger)
			}
			if err = trajectory.RequireResolved(context.Background(), root, idea.Slug); err == nil {
				t.Fatal("unverified patch can close")
			}
		})
	}
}
func TestTrajectoryRuntimeInterruptedProcessRetainsDirtyState(t *testing.T) {
	root, idea := trajectoryRuntimeFixture(t)
	agent := telemetryShell("printf 'interrupted\\n' > source; exec sleep 20", false)
	r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: agent, Prompt: "timeout fixture", Timeout: 1500 * time.Millisecond, Info: LaunchInfo{Idea: idea.Slug, RunID: "interrupt-fixture", Phase: "fixup"}})
	if err == nil || r.StartedAt == nil {
		t.Fatalf("interrupted fixture: %+v %v", r, err)
	}
	s, err := trajectory.Inspect(context.Background(), root, idea.Slug)
	if err != nil {
		t.Fatal(err)
	}
	a := s.Attempts[0]
	if a.Terminal == nil || a.Terminal.Status != "failed" || a.After == nil || a.After.Clean {
		t.Fatalf("cancelled capture missing: %+v", a)
	}
	if content, err := os.ReadFile(filepath.Join(root, "source")); err != nil || string(content) != "interrupted\n" {
		t.Fatalf("interrupted patch changed: %q %v", content, err)
	}
	restoreRuntimeAttempt(t, root, idea.Slug, a, "interrupted\n")
}

func restoreRuntimeAttempt(t *testing.T, root, idea string, a trajectory.Attempt, want string) {
	t.Helper()
	if a.After == nil || a.AfterArchive == nil {
		t.Fatal("actual process output lacks a reconstructible archive")
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, idea, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "source"), "later live edit\n")
	dir := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-snapshots")
	restored, err := trajectory.RestoreSnapshot(context.Background(), dir, *a.AfterArchive, *a.After, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	actual, err := evidence.TreeDigest(restored)
	if err != nil || actual != a.After.Tree.SHA256 {
		t.Fatal("restored actual child output has another digest", err)
	}
	content, err := os.ReadFile(filepath.Join(restored, "source"))
	if err != nil || string(content) != want {
		t.Fatalf("lost actual process output after live edit: %q %v", content, err)
	}
	content, err = os.ReadFile(filepath.Join(root, "source"))
	if err != nil || string(content) != "later live edit\n" {
		t.Fatal("restoring evidence changed the live worktree")
	}
}

func TestTrajectoryRuntimeLostAuthorityCannotReportSuccessfulCapture(t *testing.T) {
	root, idea := trajectoryRuntimeFixture(t)
	b, err := budget.LoadCycleBinding(context.Background(), root, idea.Slug, budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(b.Store.Dir)
	agent := telemetryShell(`mv "$1" "$1-preserved"`, false)
	agent.HeadlessArgs = append(agent.HeadlessArgs, dir)
	r, err := RunMeasured(context.Background(), ExecOptions{Root: root, Agent: agent, Prompt: "lost authority fixture", Timeout: 5 * time.Second, Info: LaunchInfo{Idea: idea.Slug, RunID: "lost-authority", Phase: "fixup"}})
	if err == nil || r.StartedAt == nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "trajectory_failure" {
		t.Fatalf("zero exit hid required capture loss: %+v %v", r, err)
	}
	if _, err = os.Stat(filepath.Join(dir+"-preserved", "trajectory.json")); err != nil {
		t.Fatal("fixture did not preserve displaced evidence", err)
	}
}
