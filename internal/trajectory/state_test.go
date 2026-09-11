package trajectory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/budget"
)

func accountingFixture(t *testing.T) (string, *budget.CycleBinding, Policy) {
	t.Helper()
	root := t.TempDir()
	gitFixture(t, root, "init", "-q")
	idea := filepath.Join(root, "parley-deck", "ideas", "fixture")
	if err := os.MkdirAll(idea, 0700); err != nil {
		t.Fatal(err)
	}
	for path, data := range map[string]string{".gitignore": ".parley-runtime/\n", "source": "original\n", "parley-deck/ideas/fixture/00-prompt.md": "---\nidea: fixture\ntrack: deliberation\nparticipants: [builder, reviewer]\n---\n"} {
		if err := os.WriteFile(filepath.Join(root, path), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	gitFixture(t, root, "add", ".")
	gitFixture(t, root, "commit", "-qm", "Fixture baseline")
	b, err := budget.EnsureCycleBinding(context.Background(), root, "fixture", budget.Fixup, 5, 0, "", idea)
	if err != nil {
		t.Fatal(err)
	}
	p, expected, err := NewPolicy(context.Background(), root, "fixture", "builder", []Criterion{{"material", "true"}})
	if err != nil {
		t.Fatal(err)
	}
	if err = Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	b, err = budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	return root, b, p
}
func chargeFixture(t *testing.T, root string, b *budget.CycleBinding) context.Context {
	t.Helper()
	ctx := budget.WithCycleObserver(context.Background(), &Observer{Root: root})
	ctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(finish)
	if _, err = budget.ChargeCycle(ctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	return ctx
}
func TestPersistentTrajectoryCapturesDirtyFailedAttemptAndRefusesResume(t *testing.T) {
	root, b, _ := accountingFixture(t)
	ctx := chargeFixture(t, root, b)
	run, err := Begin(ctx, root, "fixture", "builder", "invocation-1")
	if err != nil {
		t.Fatal(err)
	}
	// An actual process changes source and fails. No commit/stash/discard is made.
	cmd := exec.Command("sh", "-c", "printf 'broken\\n' > source; exit 7")
	cmd.Dir = root
	err = cmd.Run()
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("expected actual failed child: %v", err)
	}
	code := exit.ExitCode()
	if err = run.Finish(context.Background(), "failed", &code); err != nil {
		t.Fatal(err)
	}
	s, err := Inspect(context.Background(), root, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	a := s.Attempts[0]
	if a.After == nil || a.After.Clean || a.After.Tree.SHA256 == a.Before.Tree.SHA256 || a.After.Tree.Commit != a.Before.Tree.Commit || a.Terminal.ExitCode == nil || *a.Terminal.ExitCode != 7 {
		t.Fatalf("lost actual failed dirty patch: %+v", a)
	}
	fresh, err := budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fresh.Reserve(budget.WithCycleObserver(context.Background(), &Observer{Root: root}), "resumed"); err == nil {
		t.Fatal("resume admitted unresolved patch")
	}
	if err = RequireResolved(context.Background(), root, "fixture"); err == nil {
		t.Fatal("completion admitted unresolved patch")
	}
	state, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(state) != 1 {
		t.Fatalf("refusal changed charges: %v %+v", err, state)
	}
	content, _ := os.ReadFile(filepath.Join(root, "source"))
	if string(content) != "broken\n" {
		t.Fatal("dirty patch was changed")
	}
	if err = run.Finish(context.Background(), "failed", &code); err == nil {
		t.Fatal("terminal replay admitted")
	}
}
func TestPersistentTrajectoryOrphanChargeIsNotEmptyHistory(t *testing.T) {
	root, b, _ := accountingFixture(t)
	// Simulate interruption after ledger publication but before observer publication.
	zero := int64(0)
	if _, err := b.Store.Reserve(context.Background(), budget.Request{ID: "orphan", Kind: budget.Fixup, ReserveMicros: &zero}, budget.Limits{}); err != nil {
		t.Fatal(err)
	}
	for _, fn := range []func() error{
		func() error { _, err := Inspect(context.Background(), root, "fixture"); return err },
		func() error { return RequireResolved(context.Background(), root, "fixture") },
		func() error {
			_, err := b.Reserve(budget.WithCycleObserver(context.Background(), &Observer{Root: root}), "next")
			return err
		},
	} {
		if err := fn(); err == nil {
			t.Fatal("missing charged history accepted")
		}
	}
}
func TestPersistentTrajectoryMissingPolicyStateAndExactCharge(t *testing.T) {
	for _, kind := range []string{"missing-state", "empty-attempts", "changed-time", "changed-entry", "duplicate-json", "removed-reference"} {
		t.Run(kind, func(t *testing.T) {
			root, b, _ := accountingFixture(t)
			chargeFixture(t, root, b)
			switch kind {
			case "missing-state":
				if err := os.Remove(statePath(*b)); err != nil {
					t.Fatal(err)
				}
			case "empty-attempts":
				s, _, err := readState(statePath(*b))
				if err != nil {
					t.Fatal(err)
				}
				s.Attempts = []Attempt{}
				if err = writeState(statePath(*b), s); err != nil {
					t.Fatal(err)
				}
			case "duplicate-json":
				data, err := os.ReadFile(statePath(*b))
				if err != nil {
					t.Fatal(err)
				}
				data = append([]byte("{\n  \"version\": 1,"), data[1:]...)
				if err = os.WriteFile(statePath(*b), data, 0600); err != nil {
					t.Fatal(err)
				}
			case "removed-reference":
				p := b.Policy
				p.TrajectorySHA256 = ""
				data, _ := canonical(p)
				if err := os.WriteFile(filepath.Join(filepath.Dir(b.Store.Dir), "policy.json"), data, 0600); err != nil {
					t.Fatal(err)
				}
			default:
				s, err := b.Store.Inspect(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				for key, e := range s.Entries {
					if kind == "changed-time" {
						e.ReservedAt = e.ReservedAt.Add(1)
						s.Entries[key] = e
					} else {
						delete(s.Entries, key)
						s.Entries[digest([]byte("replacement"))] = e
					}
					break
				}
				data, _ := json.MarshalIndent(s, "", "  ")
				if err = os.WriteFile(filepath.Join(b.Store.Dir, "ledger.json"), data, 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := RequireResolved(context.Background(), root, "fixture"); err == nil {
				t.Fatal("changed trajectory accepted")
			}
			fresh, err := budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
			if err == nil {
				_, err = fresh.Reserve(budget.WithCycleObserver(context.Background(), &Observer{Root: root}), "next")
			}
			if err == nil {
				t.Fatal("new work admitted with changed history")
			}
		})
	}
}
func TestPersistentTrajectoryPinsActorAndDoesNotTreatExitZeroAsVerification(t *testing.T) {
	root, b, _ := accountingFixture(t)
	ctx := chargeFixture(t, root, b)
	if _, err := Begin(ctx, root, "fixture", "other", "invocation-1"); err == nil {
		t.Fatal("wrong implementer admitted")
	}
	run, err := Begin(ctx, root, "fixture", "builder", "invocation-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = Begin(ctx, root, "fixture", "builder", "invocation-2"); err == nil {
		t.Fatal("second model invocation on the same charge")
	}
	code := 0
	if err = run.Finish(context.Background(), "process-exited", &code); err != nil {
		t.Fatal(err)
	}
	if err = RequireResolved(context.Background(), root, "fixture"); err == nil {
		t.Fatal("process success became verification")
	}
}
func TestPersistentTrajectorySourceUnavailableIsRetained(t *testing.T) {
	root, b, _ := accountingFixture(t)
	ctx := chargeFixture(t, root, b)
	run, err := Begin(ctx, root, "fixture", "builder", "invocation-1")
	if err != nil {
		t.Fatal(err)
	}
	// A dangling source symlink makes a complete code-tree digest unavailable.
	if err = os.Symlink("missing-target", filepath.Join(root, "unsupported")); err != nil {
		t.Skip(err)
	}
	if err = run.Finish(context.Background(), "failed", nil); err == nil {
		t.Fatal("unavailable source was reported captured")
	}
	s, err := Inspect(context.Background(), root, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	a := s.Attempts[0]
	if a.After != nil || a.Terminal == nil || a.Terminal.SnapshotError != "source-unavailable" {
		t.Fatalf("expected unavailable post-state retained: %+v", a)
	}
}
func TestPersistentTrajectoryActivationReplayAndNoObserverRefusal(t *testing.T) {
	root, b, p := accountingFixture(t)
	original := b.Policy
	original.TrajectorySHA256 = ""
	expected := budget.CyclePolicyDigest(original)
	if err := Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Reserve(context.Background(), "unobserved"); err == nil {
		t.Fatal("required observer bypassed")
	}
	s, err := b.Store.Inspect(context.Background())
	if err != nil || len(s.Entries) != 0 {
		t.Fatalf("refused unobserved work charged: %v %+v", err, s)
	}
	if err := os.WriteFile(filepath.Join(root, "source"), []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Activate(context.Background(), root, expected, p); err == nil {
		t.Fatal("activation replay silently changed baseline")
	}
	if _, err = b.Reserve(budget.WithCycleObserver(context.Background(), &Observer{Root: root}), "changed-baseline"); err == nil {
		t.Fatal("unobserved baseline change admitted")
	}
}

func TestPersistentTrajectorySeparateProcessHelper(t *testing.T) {
	root := os.Getenv("PARLEY_TRAJECTORY_FIXTURE_ROOT")
	if root == "" {
		return
	}
	s, err := Inspect(context.Background(), root, "fixture")
	if err != nil || s == nil || len(s.Attempts) != 1 {
		t.Fatalf("separate-process coverage: %+v %v", s, err)
	}
	if err = RequireResolved(context.Background(), root, "fixture"); err == nil {
		t.Fatal("fresh process allowed unverified closure")
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = b.Reserve(budget.WithCycleObserver(context.Background(), &Observer{Root: root}), "fresh-process"); err == nil {
		t.Fatal("fresh process launched a new unverified patch")
	}
}
func TestPersistentTrajectoryWorktreeAndSeparateProcessContinuity(t *testing.T) {
	root, b, _ := accountingFixture(t)
	chargeFixture(t, root, b)
	linked := filepath.Join(t.TempDir(), "linked")
	gitFixture(t, root, "worktree", "add", "--detach", linked, "HEAD")
	shared, err := budget.LoadCycleBinding(context.Background(), linked, "fixture", budget.Fixup)
	if err != nil || shared.Store.Dir != b.Store.Dir {
		t.Fatalf("worktree accounting split: %+v %v", shared, err)
	}
	status, err := budget.InspectCycleBudget(context.Background(), linked, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	_, err = budget.ExtendCycleBudget(context.Background(), linked, "fixture", budget.Fixup, budget.CycleExtensionRequest{DecisionID: "finite-extension-fixture", ExpectedPolicySHA256: status.PolicySHA256, Maximum: 6, Reason: "Verify that extending budget does not resolve trajectory evidence"})
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestPersistentTrajectorySeparateProcessHelper$", "-test.v")
	cmd.Env = append(os.Environ(), "PARLEY_TRAJECTORY_FIXTURE_ROOT="+linked)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("separate process failed: %v %s", err, output)
	}
	ledger, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(ledger) != 1 {
		t.Fatalf("worktree restart or extension changed charges: %+v %v", ledger, err)
	}
}
func TestPersistentTrajectoryIncompleteActivationRequiresExactReplay(t *testing.T) {
	root, b, p := accountingFixture(t)
	original := b.Policy
	original.TrajectorySHA256 = ""
	data, err := canonical(original)
	if err != nil {
		t.Fatal(err)
	}
	// Initialization was published but the policy reference was not acknowledged.
	if err = os.WriteFile(filepath.Join(filepath.Dir(b.Store.Dir), "policy.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if err = RequireResolved(context.Background(), root, "fixture"); err == nil {
		t.Fatal("incomplete activation became disabled policy")
	}
	fresh, err := budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fresh.Reserve(context.Background(), "unobserved"); err == nil {
		t.Fatal("incomplete activation permitted unobserved work")
	}
	if err = Activate(context.Background(), root, budget.CyclePolicyDigest(original), p); err != nil {
		t.Fatal(err)
	}
	if err = RequireResolved(context.Background(), root, "fixture"); err != nil {
		t.Fatal(err)
	}
}

func TestPersistentTrajectoryLiveHandleRejectsErasedAuthority(t *testing.T) {
	for _, stage := range []string{"before-launch", "after-launch"} {
		t.Run(stage, func(t *testing.T) {
			root, b, _ := accountingFixture(t)
			ctx := chargeFixture(t, root, b)
			var run *Run
			var err error
			if stage == "after-launch" {
				run, err = Begin(ctx, root, "fixture", "builder", "invocation-1")
				if err != nil {
					t.Fatal(err)
				}
			}
			dir := filepath.Dir(b.Store.Dir)
			if err = os.Rename(dir, dir+"-preserved-fixture"); err != nil {
				t.Fatal(err)
			}
			// The already-held runtime identity must remember observation was required.
			if stage == "before-launch" {
				_, err = Begin(ctx, root, "fixture", "builder", "invocation-1")
			} else {
				code := 0
				err = run.Finish(context.Background(), "process-exited", &code)
			}
			if err == nil {
				t.Fatal("live required capture became disabled after authority disappeared")
			}
		})
	}
}
