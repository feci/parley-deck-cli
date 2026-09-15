package trajectory

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func gitFixture(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root, "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid"}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fixture Git operation: %v: %s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func fixtureSnapshots(t *testing.T, states [][2]string) []string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("criterion execution currently requires POSIX")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git unavailable")
	}
	repo := t.TempDir()
	gitFixture(t, repo, "init", "-q")
	var commits []string
	for i, state := range states {
		for name, value := range map[string]string{"first": state[0], "second": state[1], "patch-note": fmt.Sprint(i)} {
			if err := os.WriteFile(filepath.Join(repo, name), []byte(value), 0600); err != nil {
				t.Fatal(err)
			}
		}
		gitFixture(t, repo, "add", "--all")
		gitFixture(t, repo, "commit", "-qm", fmt.Sprintf("Fixture patch %d", i))
		commits = append(commits, gitFixture(t, repo, "rev-parse", "HEAD"))
	}
	var roots []string
	for _, commit := range commits {
		root := filepath.Join(t.TempDir(), "snapshot")
		gitFixture(t, repo, "clone", "-q", "--no-hardlinks", repo, root)
		gitFixture(t, root, "checkout", "-q", "--detach", commit)
		roots = append(roots, root)
	}
	return roots
}

func fixtureCriteria() []Criterion {
	var out []Criterion
	for _, name := range []string{"first", "second"} {
		cmd := fmt.Sprintf(`if [ "$(cat %s)" = good ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; fi`, name)
		out = append(out, Criterion{Name: name, Command: cmd})
	}
	return out
}

func fixtureTrail(t *testing.T, states [][2]string) ([]Request, []Observation) {
	t.Helper()
	roots := fixtureSnapshots(t, states)
	var requests []Request
	var observations []Observation
	previous := ""
	for i := 1; i < len(roots); i++ {
		r, err := Freeze(context.Background(), Request{Idea: "fixture", ID: fmt.Sprintf("patch-%d", i), Sequence: i, PreviousSHA256: previous, Implementer: "fixture-author", Verifier: "fixture-verifier"}, roots[i-1], roots[i], fixtureCriteria())
		if err != nil {
			t.Fatal(err)
		}
		o, err := Verify(context.Background(), r, roots[i-1], roots[i], fixtureCriteria())
		if err != nil {
			t.Fatal(err)
		}
		requests, observations = append(requests, r), append(observations, o)
		previous, err = r.SHA256()
		if err != nil {
			t.Fatal(err)
		}
	}
	return requests, observations
}

func TestRealPairedExecutionsTriggerAndRetainTrajectoryReview(t *testing.T) {
	r, o := fixtureTrail(t, [][2]string{{"good", "good"}, {"bad", "good"}, {"good", "bad"}, {"good", "good"}})
	d, err := Evaluate(r, o)
	if err != nil {
		t.Fatal(err)
	}
	if !d.ReviewRequired || d.TriggerSequence != 2 || d.Consecutive != 0 || d.Pending {
		t.Fatalf("wrong retained decision: %+v", d)
	}
	if d.Assessments[0].Regressed[0] != "first" || d.Assessments[1].Regressed[0] != "second" || d.Assessments[2].Outcome != NoRegression {
		t.Fatalf("wrong material patch outcomes: %+v", d.Assessments)
	}
}

func TestCleanInterveningPatchBreaksConsecutiveRegressions(t *testing.T) {
	r, o := fixtureTrail(t, [][2]string{{"good", "good"}, {"bad", "good"}, {"good", "good"}, {"good", "bad"}})
	d, err := Evaluate(r, o)
	if err != nil {
		t.Fatal(err)
	}
	if d.ReviewRequired || d.Consecutive != 1 || d.Pending {
		t.Fatalf("nonconsecutive patches triggered review: %+v", d)
	}
}

func TestUnchangedPriorFailureIsNotANewRegression(t *testing.T) {
	r, o := fixtureTrail(t, [][2]string{{"bad", "good"}, {"bad", "good"}, {"bad", "good"}})
	d, err := Evaluate(r, o)
	if err != nil {
		t.Fatal(err)
	}
	if d.ReviewRequired || d.Consecutive != 0 || !d.Pending {
		t.Fatalf("repeated prior failure became a regression: %+v", d)
	}
	for _, a := range d.Assessments {
		if a.Outcome != Inconclusive {
			t.Fatalf("prior failure was certified: %+v", a)
		}
	}
}

func TestVerificationRefusesUnboundIncompleteAndChangedSources(t *testing.T) {
	for _, tc := range []struct{ name, command string }{
		{"opaque", "true"},
		{"zero", `printf 'PARLEY-EVIDENCE {"executed_cases":0,"failed_cases":0}\n'`},
		{"malformed", `printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1,"failed_cases":0}\n'`},
		{"drift", `printf changed > first; printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			roots := fixtureSnapshots(t, [][2]string{{"good", "good"}, {"bad", "good"}})
			criteria := []Criterion{{Name: "acceptance", Command: tc.command}}
			r, err := Freeze(context.Background(), Request{Idea: "fixture", ID: "patch-1", Sequence: 1, Implementer: "author", Verifier: "verifier"}, roots[0], roots[1], criteria)
			if err != nil {
				t.Fatal(err)
			}
			o, err := Verify(context.Background(), r, roots[0], roots[1], criteria)
			if err == nil {
				t.Fatalf("%s produced a confirmation: %+v", tc.name, o)
			}
			if _, err := Assess(r, o); err == nil {
				t.Fatalf("partial %s observation was admitted", tc.name)
			}
		})
	}
}

func TestABBAInstabilityRemainsInconclusive(t *testing.T) {
	roots := fixtureSnapshots(t, [][2]string{{"good", "good"}, {"bad", "good"}})
	counter := filepath.Join(t.TempDir(), "calls")
	if err := os.WriteFile(counter, []byte("0"), 0600); err != nil {
		t.Fatal(err)
	}
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	cmd := `n=$(cat ` + quote(counter) + `); n=$((n+1)); printf '%s' "$n" > ` + quote(counter) + `; if [ "$n" -le 2 ]; then printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":0}\n'; else printf 'PARLEY-EVIDENCE {"executed_cases":1,"failed_cases":1}\n'; exit 1; fi`
	criteria := []Criterion{{Name: "unstable", Command: cmd}}
	r, err := Freeze(context.Background(), Request{Idea: "fixture", ID: "patch-1", Sequence: 1, Implementer: "author", Verifier: "verifier"}, roots[0], roots[1], criteria)
	if err != nil {
		t.Fatal(err)
	}
	o, err := Verify(context.Background(), r, roots[0], roots[1], criteria)
	if err != nil {
		t.Fatal(err)
	}
	a, err := Assess(r, o)
	if err != nil || a.Outcome != Inconclusive {
		t.Fatalf("unstable execution was confirmed: %+v %v", a, err)
	}
	if o.Pairs[0].Before[0].Record.Status != "pass" || o.Pairs[0].After[0].Record.Status != "pass" || o.Pairs[0].After[1].Record.Status != "fail" || o.Pairs[0].Before[1].Record.Status != "fail" {
		t.Fatal("actual invocation order was not ABBA")
	}
}
