package trajectory

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/telemetry"
)

func unchangedFixture(t *testing.T, scopeCommand string) (string, *budget.CycleBinding) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("actual POSIX fixture")
	}
	root := t.TempDir()
	gitFixture(t, root, "init", "-q")
	for name, body := range map[string]string{
		".gitignore": ".parley-runtime/\n", "source": "original\n",
		"parley-deck/ideas/fixture/00-prompt.md": "---\nidea: fixture\ntrack: deliberation\nparticipants: [builder, reviewer]\nchecks:\n  - name: material\n    command: '" + scopeCommand + "'\n---\n",
	} {
		snapshotWrite(t, root, name, []byte(body), 0600)
	}
	gitFixture(t, root, "add", ".")
	gitFixture(t, root, "commit", "-qm", "Unchanged fixture")
	b, err := budget.EnsureCycleBinding(context.Background(), root, "fixture", budget.Fixup, 5, 0, "", filepath.Join(root, "parley-deck", "ideas", "fixture"))
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
	return root, b
}

func unchangedProcess(t *testing.T, root string, b *budget.CycleBinding, command string, precharged bool) telemetry.Record {
	t.Helper()
	ctx, finish, err := budget.OpenCycleSession(budget.WithCycleObserver(context.Background(), &Observer{Root: root}), b)
	if err != nil {
		t.Fatal(err)
	}
	defer finish()
	if precharged {
		if _, err = budget.ChargeCycle(ctx, budget.Fixup); err != nil {
			t.Fatal(err)
		}
	}
	i, err := telemetry.Begin(filepath.Join(root, ".parley-runtime", "invocations"), telemetry.Metadata{RunID: "unchanged-fixture", Idea: "fixture", Phase: "fixup", Agent: "builder", LaunchMode: "headless"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = budget.ChargeCycle(ctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	run, err := Begin(ctx, root, "fixture", "builder", i.ID)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = root
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err = i.Started(cmd.Process.Pid); err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		t.Fatal(err)
	}
	err = cmd.Wait()
	status, failure, code := "process-exited", "", 0
	if err != nil {
		var e *exec.ExitError
		if !errors.As(err, &e) {
			t.Fatal(err)
		}
		status, failure, code = "failed", "process_failure", e.ExitCode()
	}
	if err = run.Finish(context.Background(), status, &code); err != nil {
		t.Fatal(err)
	}
	if err = i.Finish(telemetry.Outcome{Status: status, FailureClass: telemetry.String(failure), ExitCode: &code}); err != nil {
		t.Fatal(err)
	}
	return i.Snapshot()
}

func TestUnchangedReconciliationPreservesSpentAttemptAndExactReplay(t *testing.T) {
	for _, command := range []string{"exit 0", "exit 7", "git -c user.email=t@t -c user.name=t commit --allow-empty -qm metadata-only"} {
		for _, precharged := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/precharged=%v", command, precharged), func(t *testing.T) {
				ctx := context.Background()
				root, b := unchangedFixture(t, "true")
				record := unchangedProcess(t, root, b, command, precharged)
				s, err := Inspect(ctx, root, "fixture")
				if err != nil {
					t.Fatal(err)
				}
				original, _ := canonical(s.Attempts)
				p, err := PreviewUnchanged(ctx, root, "fixture", 1)
				if err != nil {
					t.Fatal(err)
				}
				if p.Unchanged == nil || p.Unchanged.InvocationID != record.InvocationID || p.RunID != "" || p.ParentSHA256 != "" || p.Assessment.Outcome != Inconclusive || len(p.Assessment.Regressed) != 0 || !slices.Equal(p.Assessment.Unresolved, []string{"material"}) {
					t.Fatal("unchanged became acceptance or lost provenance", p)
				}
				if strings.HasPrefix(command, "git") && s.Attempts[0].Before.Tree.Commit == s.Attempts[0].After.Tree.Commit {
					t.Fatal("fixture did not change commit metadata")
				}
				if _, err = FreezeCaptured(ctx, root, "fixture", "reviewer"); err == nil {
					t.Fatal("unchanged source became a patch verification request")
				}
				if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, strings.Repeat("0", 64)); err == nil {
					t.Fatal("stale preview accepted")
				}
				if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err != nil {
					t.Fatal(err)
				}
				retained := snapshotRead(t, statePath(*b))
				for i := 0; i < 2; i++ {
					replay, err := PreviewUnchanged(ctx, root, "fixture", 1)
					if err != nil || !sameJSON(replay, p) {
						t.Fatal("preview changed after apply", err)
					}
					applied, err := ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256())
					if err != nil || !sameJSON(applied, p) || !bytes.Equal(retained, snapshotRead(t, statePath(*b))) {
						t.Fatal("exact replay rewrote history", err)
					}
				}
				if err = RequireResolved(ctx, root, "fixture"); err == nil {
					t.Fatal("unchanged observation granted completion")
				}
				if _, err = b.Reserve(budget.WithCycleObserver(ctx, &Observer{Root: root}), "unattended-next"); err == nil {
					t.Fatal("unchanged observation granted continuation")
				}
				c, err := PreviewContinuation(ctx, root, "fixture")
				if err != nil {
					t.Fatal(err)
				}
				if _, err = Continue(ctx, root, "fixture", c.SHA256(), "ack", "test-owned acknowledgment", false, false); err == nil {
					t.Fatal("missing inconclusive acknowledgment accepted")
				}
				if _, err = Continue(ctx, root, "fixture", c.SHA256(), "ack", "test-owned acknowledgment", false, true); err != nil {
					t.Fatal(err)
				}
				if err = RequireResolved(ctx, root, "fixture"); err == nil {
					t.Fatal("acknowledgment granted completion")
				}
				unchangedProcess(t, root, b, "exit 0", precharged)
				later := snapshotRead(t, statePath(*b))
				if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err != nil || !bytes.Equal(later, snapshotRead(t, statePath(*b))) {
					t.Fatal("historical replay damaged later attempt", err)
				}
				s, err = Inspect(ctx, root, "fixture")
				if err != nil {
					t.Fatal(err)
				}
				first, _ := canonical(s.Attempts[:1])
				ledger, err := b.Store.Inspect(ctx)
				if err != nil || len(s.Attempts) != 2 || len(s.Resolutions) != 1 || b.Count(ledger) != 2 || !bytes.Equal(original, first) {
					t.Fatal("original charge/attempt lost", err)
				}
			})
		}
	}
}

func TestUnchangedReconciliationRejectsMissingAndContraryEvidence(t *testing.T) {
	root, b := unchangedFixture(t, "true")
	r := unchangedProcess(t, root, b, "exit 0", false)
	ctx := context.Background()
	p, err := PreviewUnchanged(ctx, root, "fixture", 1)
	if err != nil {
		t.Fatal(err)
	}
	inv := filepath.Join(root, ".parley-runtime", "invocations", r.InvocationID)
	for _, name := range []string{"requested.json", "started.json", "terminal.json"} {
		file := filepath.Join(inv, name)
		original := snapshotRead(t, file)
		if err = os.Remove(file); err != nil {
			t.Fatal(err)
		}
		if _, err = PreviewUnchanged(ctx, root, "fixture", 1); err == nil {
			t.Fatal("missing lifecycle accepted", name)
		}
		snapshotWrite(t, inv, name, original, 0600)
	}
	terminalFile := filepath.Join(inv, "terminal.json")
	original := snapshotRead(t, terminalFile)
	mutations := map[string]func(*telemetry.Record){
		"pid":          func(r *telemetry.Record) { *r.PID++ },
		"identity":     func(r *telemetry.Record) { r.Metadata.Agent = "reviewer" },
		"run":          func(r *telemetry.Record) { r.Metadata.RunID = "other" },
		"before-start": func(r *telemetry.Record) { r.CompletedAt = &r.RequestedAt },
		"no-start":     func(r *telemetry.Record) { r.StartedAt = nil },
		"no-exit":      func(r *telemetry.Record) { r.Outcome.ExitCode = nil },
		"signal":       func(r *telemetry.Record) { *r.Outcome.ExitCode = -1 },
		"wrong-exit":   func(r *telemetry.Record) { *r.Outcome.ExitCode = 7 },
		"timeout":      func(r *telemetry.Record) { r.Outcome.FailureClass = telemetry.String("timeout") },
		"cancelled":    func(r *telemetry.Record) { r.Outcome.FailureClass = telemetry.String("cancelled") },
		"provider":     func(r *telemetry.Record) { r.Outcome.FailureClass = telemetry.String("rate_limit") },
		"handoff":      func(r *telemetry.Record) { r.Outcome.Status = "unobserved-handoff" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			var altered telemetry.Record
			if err := json.Unmarshal(original, &altered); err != nil {
				t.Fatal(err)
			}
			mutate(&altered)
			data, _ := canonical(altered)
			snapshotWrite(t, inv, "terminal.json", data, 0600)
			_, previewErr := PreviewUnchanged(ctx, root, "fixture", 1)
			_, err := ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256())
			snapshotWrite(t, inv, "terminal.json", original, 0600)
			if previewErr == nil || err == nil {
				t.Fatal("contrary process lifecycle accepted by fresh preview or apply")
			}
		})
	}
	if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err != nil {
		t.Fatal(err)
	}
	// Even an otherwise valid later usage observation cannot replace already
	// pinned lifecycle bytes. State reads and exact apply rederive this evidence.
	var altered telemetry.Record
	json.Unmarshal(original, &altered)
	altered.Warnings = append(altered.Warnings, "later-observation")
	data, _ := canonical(altered)
	snapshotWrite(t, inv, "terminal.json", data, 0600)
	if _, err = Inspect(ctx, root, "fixture"); err == nil {
		t.Fatal("accepted terminal changed without invalidating history")
	}
	if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err == nil {
		t.Fatal("exact replay accepted changed lifecycle")
	}
}

func TestUnchangedReconciliationRequiresOriginalArchivedScope(t *testing.T) {
	ctx := context.Background()
	root, b := unchangedFixture(t, "false")
	unchangedProcess(t, root, b, "exit 0", false)
	file := filepath.Join(root, "parley-deck", "ideas", "fixture", "00-prompt.md")
	// Current checks are changed to match the policy; this must not hide that
	// the archive did not contain that original material command.
	raw := snapshotRead(t, file)
	if err := os.WriteFile(file, bytes.Replace(raw, []byte("'false'"), []byte("'true'"), 1), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := PreviewUnchanged(ctx, root, "fixture", 1); err == nil || !strings.Contains(err.Error(), "archived material") {
		t.Fatal("current scope substituted for archived original", err)
	}
	root, b = unchangedFixture(t, "true")
	unchangedProcess(t, root, b, "exit 0", false)
	file = filepath.Join(root, "parley-deck", "ideas", "fixture", "00-prompt.md")
	raw = snapshotRead(t, file)
	for _, pair := range [][2]string{{"[builder, reviewer]", "[builder, other]"}, {"material", "changed"}, {"'true'", "'false'"}} {
		if err := os.WriteFile(file, bytes.ReplaceAll(raw, []byte(pair[0]), []byte(pair[1])), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := PreviewUnchanged(ctx, root, "fixture", 1); err == nil {
			t.Fatal("changed current contract accepted", pair)
		}
	}
	if err := os.WriteFile(file, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := PreviewUnchanged(ctx, root, "fixture", 1); err != nil {
		t.Fatal(err)
	}
}

func TestUnchangedReconciliationPublicationRecovery(t *testing.T) {
	for _, after := range []bool{false, true} {
		t.Run(fmt.Sprint(after), func(t *testing.T) {
			ctx := context.Background()
			root, b := unchangedFixture(t, "true")
			unchangedProcess(t, root, b, "exit 7", false)
			before := snapshotRead(t, statePath(*b))
			p, err := PreviewUnchanged(ctx, root, "fixture", 1)
			if err != nil {
				t.Fatal(err)
			}
			fault := errors.New("injected publication interruption")
			_, err = reconcileUnchanged(ctx, root, "fixture", 1, p.SHA256(), func(path string, s State) error {
				if after {
					if err := writeState(path, s); err != nil {
						return err
					}
				}
				return fault
			})
			if !errors.Is(err, fault) {
				t.Fatal("publication interruption hidden", err)
			}
			retained := snapshotRead(t, statePath(*b))
			if after == bytes.Equal(before, retained) {
				t.Fatal("wrong publication side of interruption")
			}
			preview, err := PreviewUnchanged(ctx, root, "fixture", 1)
			if err != nil || !sameJSON(preview, p) {
				t.Fatal("interruption changed exact preview", err)
			}
			if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err != nil {
				t.Fatal(err)
			}
			if after && !bytes.Equal(retained, snapshotRead(t, statePath(*b))) {
				t.Fatal("post-publication retry rewrote history")
			}
		})
	}
}

func TestUnchangedHistoryPreservesRegressionStreakAndAcknowledgedTrigger(t *testing.T) {
	for _, tc := range []struct {
		outcomes             []Outcome
		consecutive, trigger int
	}{
		{[]Outcome{Regression, Inconclusive, Regression}, 2, 3},
		{[]Outcome{Regression, Regression, Inconclusive}, 2, 2},
		{[]Outcome{Regression, Regression, Inconclusive, Regression}, 3, 4},
		{[]Outcome{Regression, Inconclusive, NoRegression}, 0, 0},
	} {
		s := State{}
		for i, outcome := range tc.outcomes {
			p := ReconciliationPreview{Sequence: i + 1, Assessment: Assessment{Outcome: outcome}}
			if outcome == Inconclusive {
				p.Unchanged = &UnchangedEvidence{}
			}
			s.Resolutions = append(s.Resolutions, Resolution{Preview: p})
			s.Attempts = append(s.Attempts, Attempt{})
		}
		h := trajectoryHistory(s)
		if h.Decision.Consecutive != tc.consecutive || h.RequiredReviewSequence != tc.trigger || len(h.InconclusivePending) != 1 {
			t.Fatal("unchanged observation altered regression history", tc, h)
		}
		if tc.trigger == 2 {
			s.Continuations = []Continuation{{Preview: ContinuationPreview{Sequence: 2}, ReviewThrough: 2}}
			h = trajectoryHistory(s)
			if h.ReviewPending || h.RequiredReviewSequence != 2 || h.Decision.TriggerSequence != 2 {
				t.Fatal("unchanged observation reopened acknowledged review", h)
			}
		}
	}
}

func TestUnchangedReconciliationRejectsRehashedStateCorruption(t *testing.T) {
	ctx := context.Background()
	root, b := unchangedFixture(t, "true")
	unchangedProcess(t, root, b, "exit 0", false)
	p, err := PreviewUnchanged(ctx, root, "fixture", 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256()); err != nil {
		t.Fatal(err)
	}
	original := snapshotRead(t, statePath(*b))
	for name, mutate := range map[string]func(*ReconciliationPreview){
		"invented-parent":     func(p *ReconciliationPreview) { p.RunID = "invented"; p.ParentSHA256 = strings.Repeat("b", 64) },
		"invented-recovery":   func(p *ReconciliationPreview) { p.RecoverySHA256 = strings.Repeat("b", 64) },
		"clean-verdict":       func(p *ReconciliationPreview) { p.Assessment.Outcome = NoRegression },
		"regression-verdict":  func(p *ReconciliationPreview) { p.Assessment.Outcome = Regression },
		"reduced-scope":       func(p *ReconciliationPreview) { p.Assessment.Unresolved = nil },
		"regressed-criterion": func(p *ReconciliationPreview) { p.Assessment.Regressed = []string{"material"} },
		"attempt-binding":     func(p *ReconciliationPreview) { p.Unchanged.AttemptSHA256 = strings.Repeat("b", 64) },
		"original-scope":      func(p *ReconciliationPreview) { p.Unchanged.ScopeSHA256 = strings.Repeat("b", 64) },
		"schema":              func(p *ReconciliationPreview) { p.Unchanged.Version++ },
	} {
		t.Run(name, func(t *testing.T) {
			var altered State
			if err := json.Unmarshal(original, &altered); err != nil {
				t.Fatal(err)
			}
			mutate(&altered.Resolutions[0].Preview)
			altered.Resolutions[0].SHA256 = altered.Resolutions[0].Preview.SHA256()
			if err := writeState(statePath(*b), altered); err != nil {
				t.Fatal(err)
			}
			_, err := Inspect(ctx, root, "fixture")
			if err == nil {
				t.Fatal("rehashed state replaced original evidence or assessment")
			}
			if err = os.WriteFile(statePath(*b), original, 0600); err != nil {
				t.Fatal(err)
			}
		})
	}
	if _, err = Inspect(ctx, root, "fixture"); err != nil {
		t.Fatal("fixture did not restore original history", err)
	}
}

func TestUnchangedSnapshotMemberRequiresCompleteArchive(t *testing.T) {
	root, dir := snapshotFixture(t), snapshotStoreFixture(t)
	source, ref := captureFixture(t, root, dir)
	original := snapshotRead(t, snapshotPath(dir, ref))
	data, err := readSnapshotMember(context.Background(), dir, ref, source, ".gitignore", 1024)
	if err != nil || string(data) != "ignored/\n" {
		t.Fatal("valid selection failed", err)
	}
	for _, mode := range []string{"missing", "limit", "footer", "corrupt-remainder", "wrong-whole-source", "link"} {
		t.Run(mode, func(t *testing.T) {
			want, changed, name, limit := source, append([]byte{}, original...), ".gitignore", int64(1024)
			switch mode {
			case "missing":
				name = "nonexistent"
			case "limit":
				limit = 1
			case "footer":
				changed = changed[:len(changed)-1024]
			case "corrupt-remainder":
				changed[bytes.Index(changed, []byte("retained at baseline"))] ^= 1
			case "wrong-whole-source", "link":
				reader := tar.NewReader(bytes.NewReader(original))
				reader.Next()
				var members []archiveFixtureMember
				for {
					h, err := reader.Next()
					if err == io.EOF {
						break
					}
					if err != nil {
						t.Fatal(err)
					}
					body, err := io.ReadAll(reader)
					if err != nil {
						t.Fatal(err)
					}
					members = append(members, archiveFixtureMember{h, body})
				}
				if mode == "wrong-whole-source" {
					want.Tree.SHA256 = strings.Repeat("a", 64)
				} else {
					members[0] = archiveFixtureMember{snapshotHeader("files/.gitignore", tar.TypeSymlink, 0777, 0, "source"), nil}
				}
				changed = archiveFixtureBytes(t, want, members)
			}
			bad := archiveFixtureStore(t, dir, changed)
			data, err := readSnapshotMember(context.Background(), dir, bad, want, name, limit)
			if err == nil || data != nil {
				t.Fatal("selected bytes escaped incomplete archive validation", mode, err)
			}
		})
	}
}
