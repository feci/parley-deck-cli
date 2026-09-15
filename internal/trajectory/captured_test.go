package trajectory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

func capturedCriteria() []Criterion {
	return []Criterion{{Name: "material", Command: `if [ "$(cat source)" = original ] && [ -f removed ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; fi`}}
}

func capturedFixture(t *testing.T, criteria []Criterion, mutate func(string)) (string, *budget.CycleBinding, CapturedRequest) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("criterion execution requires POSIX")
	}
	root := snapshotFixture(t)
	snapshotWrite(t, root, ".gitignore", []byte("ignored/\n.parley-runtime/\n"), 0600)
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "fixture")
	snapshotWrite(t, root, "parley-deck/ideas/fixture/00-prompt.md", []byte("---\nidea: fixture\ntrack: deliberation\nparticipants: [builder, reviewer]\n---\n"), 0600)
	gitFixture(t, root, "add", ".")
	gitFixture(t, root, "commit", "-qm", "Captured verification baseline")
	_, err := budget.EnsureCycleBinding(context.Background(), root, "fixture", budget.Fixup, 5, 0, "", ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	p, expected, err := NewPolicy(context.Background(), root, "fixture", "builder", criteria)
	if err != nil {
		t.Fatal(err)
	}
	if err = Activate(context.Background(), root, expected, p); err != nil {
		t.Fatal(err)
	}
	b, err := budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	ctx := chargeFixture(t, root, b)
	run, err := Begin(ctx, root, "fixture", "builder", "actual-patch-fixture")
	if err != nil {
		t.Fatal(err)
	}
	if mutate != nil {
		mutate(root)
	}
	code := 7
	if err = run.Finish(context.Background(), "failed", &code); err != nil {
		t.Fatal(err)
	}
	r, err := FreezeCaptured(context.Background(), root, "fixture", "reviewer")
	if mutate == nil {
		if err == nil {
			t.Fatal("unchanged criticism became a new patch")
		}
		return root, b, CapturedRequest{}
	}
	if err != nil {
		t.Fatal(err)
	}
	return root, b, r
}

func dirtyCapturedChild(t *testing.T) func(string) {
	return func(root string) {
		t.Helper()
		cmd := exec.Command("sh", "-c", "printf 'broken\\n' > source; rm removed; printf 'added\\n' > added; exit 7")
		cmd.Dir = root
		err := cmd.Run()
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 7 {
			t.Fatalf("actual dirty failing child did not execute: %v", err)
		}
	}
}

func openCapturedFixture(t *testing.T, root string, r CapturedRequest) *CapturedWorkspace {
	t.Helper()
	w, err := OpenCaptured(context.Background(), root, r, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.Close(); err != nil {
			t.Error(err)
		}
	})
	return w
}

func TestCapturedVerificationExecutesArchivedDirtyPatchWithOriginalGitIdentity(t *testing.T) {
	criteria := capturedCriteria()
	root, b, r := capturedFixture(t, criteria, dirtyCapturedChild(t))
	if r.Before.Tree.Commit != r.After.Tree.Commit || r.After.Clean {
		t.Fatal("fixture lost the actual dirty same-HEAD attempt")
	}
	// Later live edits cannot change what this comparison executes.
	snapshotWrite(t, root, "source", []byte("original\n"), 0600)
	snapshotWrite(t, root, "removed", []byte("live file was restored\n"), 0600)
	index := snapshotRead(t, filepath.Join(root, ".git", "index"))
	head := snapshotRead(t, filepath.Join(root, ".git", "HEAD"))
	refs := gitFixture(t, root, "show-ref")
	w := openCapturedFixture(t, root, r)
	for _, item := range []struct {
		root string
		want Source
	}{{w.BeforeRoot(), r.Before}, {w.AfterRoot(), r.After}} {
		actual, err := Observe(context.Background(), item.root)
		if err != nil || actual != item.want {
			t.Fatalf("reproduced source is not the original observation: %+v %v", actual, err)
		}
	}
	oid := r.Before.Tree.Commit
	object := filepath.Join("objects", oid[:2], oid[2:])
	live, err := os.Stat(filepath.Join(root, ".git", object))
	if err != nil {
		t.Fatal(err)
	}
	copy, err := os.Stat(filepath.Join(w.BeforeRoot(), ".git", object))
	if err != nil || os.SameFile(live, copy) {
		t.Fatal("verification shares writable Git objects with the live repository")
	}
	o, err := VerifyCaptured(context.Background(), w, criteria)
	if err != nil {
		t.Fatal(err)
	}
	a, err := AssessCaptured(r, o)
	if err != nil || a.Outcome != Regression || len(a.Regressed) != 1 {
		t.Fatalf("actual retained regression was not derived: %+v %v", a, err)
	}
	if o.Pairs[0].Before[0].Record.Command.ExitCode != 0 || o.Pairs[0].After[0].Record.Command.ExitCode != 1 {
		t.Fatal("comparison did not execute the actual passing/failing source")
	}
	if !bytes.Equal(index, snapshotRead(t, filepath.Join(root, ".git", "index"))) ||
		!bytes.Equal(head, snapshotRead(t, filepath.Join(root, ".git", "HEAD"))) ||
		refs != gitFixture(t, root, "show-ref") ||
		string(snapshotRead(t, filepath.Join(root, "source"))) != "original\n" {
		t.Fatal("comparison changed the live source or Git state")
	}
	ledger, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(ledger) != 1 {
		t.Fatal("verification spent or reset a fixup charge", err)
	}
	if err = RequireResolved(context.Background(), root, "fixture"); err == nil {
		t.Fatal("core execution granted acceptance without an independent helper receipt")
	}
	if _, err = VerifyCaptured(context.Background(), w, criteria); err == nil {
		t.Fatal("same workspace replayed a completed helper execution")
	}
	data, err := json.Marshal(r)
	if err != nil || bytes.Contains(data, []byte(root)) || bytes.Contains(data, []byte(criteria[0].Command)) {
		t.Fatal("request contains raw paths or commands")
	}
}

func TestCapturedVerificationRunsExactABBAOrder(t *testing.T) {
	trace := filepath.Join(t.TempDir(), "trace")
	quote := "'" + strings.ReplaceAll(trace, "'", "'\\''") + "'"
	criteria := capturedCriteria()
	criteria[0].Command = "cat source >> " + quote + "; " + criteria[0].Command
	root, _, r := capturedFixture(t, criteria, dirtyCapturedChild(t))
	w := openCapturedFixture(t, root, r)
	if _, err := VerifyCaptured(context.Background(), w, criteria); err != nil {
		t.Fatal(err)
	}
	if string(snapshotRead(t, trace)) != "original\nbroken\nbroken\noriginal\n" {
		t.Fatal("paired execution did not use AB/BA order")
	}
}

func TestCapturedVerificationRefusesRequestSubstitutionBeforePreparingRoots(t *testing.T) {
	root, _, r := capturedFixture(t, capturedCriteria(), dirtyCapturedChild(t))
	for _, tc := range []struct {
		name string
		edit func(*CapturedRequest)
	}{
		{"self-verifier", func(r *CapturedRequest) { r.Verifier = r.Implementer }},
		{"charge-key", func(r *CapturedRequest) { r.Charge.EntryKey = digest([]byte("another charge")) }},
		{"charge-time", func(r *CapturedRequest) { r.Charge.ReservedAt = r.Charge.ReservedAt.Add(time.Second) }},
		{"invocation", func(r *CapturedRequest) { r.InvocationID = "another-invocation" }},
		{"policy", func(r *CapturedRequest) { r.PolicySHA256 = digest([]byte("another policy")) }},
		{"attempt", func(r *CapturedRequest) { r.AttemptSHA256 = digest([]byte("another attempt")) }},
		{"source", func(r *CapturedRequest) { r.After.Tree.SHA256 = digest([]byte("another source")) }},
		{"archive", func(r *CapturedRequest) { r.AfterArchive.SHA256 = digest([]byte("another archive")) }},
		{"criterion", func(r *CapturedRequest) { r.Criteria[0].CommandSHA256 = digest([]byte("true")) }},
		{"kind", func(r *CapturedRequest) { r.Kind = "clean-commit" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, _ := json.Marshal(r)
			var changed CapturedRequest
			if err := json.Unmarshal(raw, &changed); err != nil {
				t.Fatal(err)
			}
			tc.edit(&changed)
			parent := t.TempDir()
			if w, err := OpenCaptured(context.Background(), root, changed, parent); err == nil || w != nil {
				t.Fatal("changed request admitted")
			}
			entries, err := os.ReadDir(parent)
			if err != nil || len(entries) != 0 {
				t.Fatal("refused request allocated verification source")
			}
		})
	}
}

func TestCapturedVerificationRefusesIncompleteOriginalGitContext(t *testing.T) {
	for _, kind := range []string{"staged", "missing-commit", "alternates", "oversize", "nested-parent"} {
		t.Run(kind, func(t *testing.T) {
			mutate := dirtyCapturedChild(t)
			if kind == "staged" {
				mutate = func(root string) {
					snapshotWrite(t, root, "source", []byte("broken\n"), 0600)
					gitFixture(t, root, "add", "source")
				}
			}
			root, _, r := capturedFixture(t, capturedCriteria(), mutate)
			parent := t.TempDir()
			switch kind {
			case "missing-commit":
				oid := r.After.Tree.Commit
				object := filepath.Join(root, ".git", "objects", oid[:2], oid[2:])
				if err := os.Rename(object, filepath.Join(t.TempDir(), "preserved-object")); err != nil {
					t.Fatal(err)
				}
			case "alternates":
				other := snapshotFixture(t)
				snapshotWrite(t, root, ".git/objects/info/alternates", []byte(filepath.Join(other, ".git", "objects")+"\n"), 0600)
			case "oversize":
				f, err := os.Create(filepath.Join(root, ".git", "objects", "oversize"))
				if err != nil {
					t.Fatal(err)
				}
				err = f.Truncate(MaxCapturedGitBytes + 1)
				closeErr := f.Close()
				if err != nil || closeErr != nil {
					t.Fatal(err, closeErr)
				}
			case "nested-parent":
				parent = root
			}
			index := snapshotRead(t, filepath.Join(root, ".git", "index"))
			if w, err := OpenCaptured(context.Background(), root, r, parent); err == nil || w != nil {
				t.Fatal("unreproducible or unsafe Git context was silently substituted")
			}
			if !bytes.Equal(index, snapshotRead(t, filepath.Join(root, ".git", "index"))) {
				t.Fatal("refused reproduction rewrote the original index")
			}
			left, err := filepath.Glob(filepath.Join(parent, "trajectory-verifier-*"))
			if err != nil || len(left) != 0 {
				t.Fatal("failed preparation left partial verification roots")
			}
		})
	}
}

func TestCapturedVerificationRetainsPartialExecutionOnSourceDrift(t *testing.T) {
	criteria := []Criterion{{Name: "material", Command: `printf 'mutated by check\n' > source; printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'`}}
	root, _, r := capturedFixture(t, criteria, dirtyCapturedChild(t))
	w := openCapturedFixture(t, root, r)
	o, err := VerifyCaptured(context.Background(), w, criteria)
	if err == nil || len(o.Pairs) != 1 || o.Pairs[0].Before[0].Record.Command.OutputSHA256 == "" {
		t.Fatalf("source drift lost the actual execution: %+v %v", o, err)
	}
	e := o.Pairs[0].Before[0]
	if e.TreeAfterSHA256 == "" || e.TreeAfterSHA256 == e.TreeBeforeSHA256 {
		t.Fatal("expected digest was substituted for actual changed source")
	}
	if _, err = AssessCaptured(r, o); err == nil {
		t.Fatal("source-changing criterion confirmed a patch outcome")
	}
}

func TestCapturedVerificationRetainsInterruptedExecution(t *testing.T) {
	signal := filepath.Join(t.TempDir(), "started")
	quote := "'" + strings.ReplaceAll(signal, "'", "'\\''") + "'"
	criteria := []Criterion{{Name: "material", Command: `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; printf started > ` + quote + "; exec sleep 20"}}
	root, _, r := capturedFixture(t, criteria, dirtyCapturedChild(t))
	w := openCapturedFixture(t, root, r)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	type result struct {
		o   Observation
		err error
	}
	done := make(chan result, 1)
	go func() {
		o, err := VerifyCaptured(ctx, w, criteria)
		done <- result{o, err}
	}()
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
started:
	for {
		select {
		case <-ctx.Done():
			t.Fatal("fixture criterion did not start before its bound")
		case early := <-done:
			t.Fatalf("fixture completed before interruption: %+v %v", early.o, early.err)
		case <-tick.C:
			if _, err := os.Stat(signal); err == nil {
				break started
			} else if !os.IsNotExist(err) {
				t.Fatal(err)
			}
		}
	}
	cancel()
	actual := <-done
	o, err := actual.o, actual.err
	if err == nil || len(o.Pairs) != 1 || o.Pairs[0].Before[0].Record.Command.OutputSHA256 == "" || o.Pairs[0].Before[0].Record.Command.OutputSHA256 == digest(nil) || o.Pairs[0].Before[0].Complete {
		t.Fatalf("interrupted execution was lost or certified: %+v %v", o, err)
	}
	if _, err = AssessCaptured(r, o); err == nil {
		t.Fatal("interrupted source confirmed a patch outcome")
	}
}

func TestCapturedVerificationRechecksDurableAuthorityAndScope(t *testing.T) {
	for _, kind := range []string{"missing-authority", "changed-terminal", "changed-scope", "changed-root", "changed-head"} {
		t.Run(kind, func(t *testing.T) {
			criteria := capturedCriteria()
			root, b, r := capturedFixture(t, criteria, dirtyCapturedChild(t))
			w := openCapturedFixture(t, root, r)
			switch kind {
			case "missing-authority":
				dir := filepath.Dir(b.Store.Dir)
				if err := os.Rename(dir, dir+"-preserved"); err != nil {
					t.Fatal(err)
				}
			case "changed-terminal":
				s, _, err := readState(statePath(*b))
				if err != nil {
					t.Fatal(err)
				}
				s.Attempts[0].Terminal.At = s.Attempts[0].Terminal.At.Add(time.Second)
				if err = writeState(statePath(*b), s); err != nil {
					t.Fatal(err)
				}
			case "changed-scope":
				criteria[0].Command = "true"
			case "changed-root":
				snapshotWrite(t, w.AfterRoot(), "added", []byte("another patch\n"), 0600)
			case "changed-head":
				gitFixture(t, w.BeforeRoot(), "update-ref", "--no-deref", "HEAD", "HEAD~1")
			}
			o, err := VerifyCaptured(context.Background(), w, criteria)
			if err == nil || len(o.Pairs) != 0 {
				t.Fatalf("changed authority/source executed criteria: %+v %v", o, err)
			}
		})
	}
}

func TestCapturedVerificationCannotFreezeUnchangedCriticism(t *testing.T) {
	root, _, _ := capturedFixture(t, capturedCriteria(), nil)
	if err := RequireResolved(context.Background(), root, "fixture"); err == nil {
		t.Fatal("unchanged source auto-resolved a charged attempt")
	}
}

func TestCapturedVerificationPreservesCommittedPatchHistory(t *testing.T) {
	criteria := capturedCriteria()
	root, _, r := capturedFixture(t, criteria, func(root string) {
		snapshotWrite(t, root, "source", []byte("broken\n"), 0600)
		gitFixture(t, root, "add", "source")
		gitFixture(t, root, "commit", "-qm", "Actual committed fixture patch")
	})
	if !r.After.Clean || r.Before.Tree.Commit == r.After.Tree.Commit {
		t.Fatal("actual committed patch was not captured")
	}
	w := openCapturedFixture(t, root, r)
	if parent := gitFixture(t, w.AfterRoot(), "rev-parse", "HEAD~1"); parent != r.Before.Tree.Commit {
		t.Fatal("actual Git ancestry was replaced with synthetic history")
	}
	o, err := VerifyCaptured(context.Background(), w, criteria)
	if err != nil {
		t.Fatal(err)
	}
	a, err := AssessCaptured(r, o)
	if err != nil || a.Outcome != Regression {
		t.Fatalf("committed captured regression: %+v %v", a, err)
	}
}

func TestCapturedVerificationDerivesCleanAndInconclusiveOutcomes(t *testing.T) {
	for _, kind := range []string{"unrelated-change", "prior-failure", "opaque-output"} {
		t.Run(kind, func(t *testing.T) {
			criteria := capturedCriteria()
			if kind == "prior-failure" {
				criteria[0].Command = `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1`
			} else if kind == "opaque-output" {
				criteria[0].Command = "true"
			}
			root, _, r := capturedFixture(t, criteria, func(root string) {
				snapshotWrite(t, root, "unrelated", []byte("new source\n"), 0600)
			})
			w := openCapturedFixture(t, root, r)
			o, err := VerifyCaptured(context.Background(), w, criteria)
			if kind == "opaque-output" {
				if err == nil || len(o.Pairs) != 1 {
					t.Fatal("opaque output became a material confirmation or lost observations")
				}
				if _, err = AssessCaptured(r, o); err == nil {
					t.Fatal("opaque output independently assessed")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			a, err := AssessCaptured(r, o)
			want := NoRegression
			if kind == "prior-failure" {
				want = Inconclusive
			}
			if err != nil || a.Outcome != want {
				t.Fatalf("outcome %s: %+v %v", kind, a, err)
			}
		})
	}
}
