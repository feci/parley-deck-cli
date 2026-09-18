package trajectory

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/budget"
)

func quorumPatchFixture(t *testing.T, members string) (string, *budget.CycleBinding, State, CapturedRequest) {
	t.Helper()
	root, b := unchangedFixture(t, "true")
	name := "parley-deck/ideas/fixture/00-prompt.md"
	body := strings.Replace(string(snapshotRead(t, filepath.Join(root, name))), "[builder, reviewer]", members, 1)
	quoted := "'" + strings.ReplaceAll(body, "'", "'\\''") + "'"
	unchangedProcess(t, root, b, "printf '%s' "+quoted+" > "+name+"; printf changed > source", false)
	s, err := Inspect(context.Background(), root, "fixture")
	if err != nil {
		t.Fatalf("quorum edit prevented retaining the original attempt: %v", err)
	}
	if len(s.Attempts) != 1 || s.Attempts[0].Terminal == nil || s.Attempts[0].AfterArchive == nil || len(s.Resolutions) != 0 {
		t.Fatal("actual charged terminal/source was not retained")
	}
	// This structural request also models an old ticket issued before the new
	// permission gate. It is not an independently accepted result or new grant.
	r, err := capturedRequestAt(*s, "reviewer", 1)
	if err != nil {
		t.Fatal(err)
	}
	return root, b, *s, r
}

func requireQuorumRefusal(t *testing.T, err error, fragment string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), fragment) {
		t.Fatalf("expected quorum refusal %q, got %v", fragment, err)
	}
}

func TestActivationQuorumRejectsActualChargedMembershipChanges(t *testing.T) {
	for _, members := range []string{"[builder, other]", "[builder, reviewer, other]", "[reviewer, builder]", "[builder]", "[reviewer, other]", "[builder, reviewer, reviewer]"} {
		t.Run(members, func(t *testing.T) {
			root, b, s, r := quorumPatchFixture(t, members)
			stateBefore := snapshotRead(t, statePath(*b))
			ledgerBefore, err := b.Store.Inspect(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, err = checkCapturedActivationQuorum(context.Background(), *b, s, r)
			requireQuorumRefusal(t, err, "quorum")
			_, err = FreezeCaptured(context.Background(), root, "fixture", "reviewer")
			requireQuorumRefusal(t, err, "quorum")
			_, err = PrepareCapturedVerification(context.Background(), root, "fixture", "reviewer", "quorum-test")
			requireQuorumRefusal(t, err, "quorum")
			_, err = checkCapturedAuthority(context.Background(), root, r)
			requireQuorumRefusal(t, err, "quorum")
			journal := filepath.Join(filepath.Dir(b.Store.Dir), "trajectory-verifications", r.Charge.EntryKey)
			if _, err = os.Lstat(journal); !os.IsNotExist(err) {
				t.Fatalf("refused quorum consumed a helper ticket: %v", err)
			}
			ledgerAfter, err := b.Store.Inspect(context.Background())
			if err != nil || !sameJSON(ledgerBefore, ledgerAfter) || !bytes.Equal(stateBefore, snapshotRead(t, statePath(*b))) {
				t.Fatalf("quorum refusal rewrote accounting/history: %v", err)
			}
		})
	}
}

func TestActivationQuorumSelectedVerifierAndValidReplay(t *testing.T) {
	root, b, s, r := quorumPatchFixture(t, "[builder, reviewer]")
	members, err := activationQuorum(context.Background(), *b, s)
	if err != nil || !slices.Equal(members, []string{"builder", "reviewer"}) {
		t.Fatalf("original members lost: %v %v", members, err)
	}
	for _, verifier := range []string{"other", "builder"} {
		r.Verifier = verifier
		_, err := checkCapturedActivationQuorum(context.Background(), *b, s, r)
		requireQuorumRefusal(t, err, "independent activation quorum")
	}
	_, err = FreezeCaptured(context.Background(), root, "fixture", "other")
	requireQuorumRefusal(t, err, "independent activation quorum")
	_, err = PrepareCapturedVerification(context.Background(), root, "fixture", "other", "unlisted-verifier")
	requireQuorumRefusal(t, err, "independent activation quorum")
	first, err := FreezeCaptured(context.Background(), root, "fixture", "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	second, err := FreezeCaptured(context.Background(), root, "fixture", "reviewer")
	if err != nil || !sameJSON(first, second) {
		t.Fatalf("valid request replay changed: %v", err)
	}
	ticket, err := PrepareCapturedVerification(context.Background(), root, "fixture", "reviewer", "valid-quorum")
	if err != nil {
		t.Fatal(err)
	}
	if err := ReserveCapturedVerificationLaunch(context.Background(), ticket, "valid-quorum-launch"); err != nil {
		t.Fatal(err)
	}
	if err := StopCapturedVerification(context.Background(), ticket, "valid-quorum-launch"); err != nil {
		t.Fatal(err)
	}
}

func TestActivationQuorumLegacyTicketCanStopWithoutExecutionPermission(t *testing.T) {
	root, b, s, r := quorumPatchFixture(t, "[builder, other]")
	root, err := canonicalRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	r.Verifier = "other"
	ticket := VerificationTicket{1, root, "old-quorum-ticket", r}
	sha, err := ticket.SHA256()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := openVerificationDirectory(*b, r.Charge.EntryKey, true)
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	if _, err := writeVerificationArtifact(dir, "request.json", ticket); err != nil {
		t.Fatal(err)
	}
	if _, err := writeVerificationArtifact(dir, "launch.json", verificationLaunch{1, sha, "old-verifier-launch", time.Now().UTC()}); err != nil {
		t.Fatal(err)
	}
	stateBefore := snapshotRead(t, statePath(*b))
	requestBefore, err := dir.ReadFile("request.json")
	if err != nil {
		t.Fatal(err)
	}
	launchBefore, err := dir.ReadFile("launch.json")
	if err != nil {
		t.Fatal(err)
	}
	entered := false
	err = withVerification(context.Background(), ticket, func(*os.Root, string) error { entered = true; return nil })
	requireQuorumRefusal(t, err, "independent activation quorum")
	if entered {
		t.Fatal("old invalid quorum released execution callback")
	}
	err = ReserveCapturedVerificationLaunch(context.Background(), ticket, "new-verifier-launch")
	requireQuorumRefusal(t, err, "independent activation quorum")
	if err := StopCapturedVerification(context.Background(), ticket, "wrong-launch"); err == nil {
		t.Fatal("quorum cleanup bypassed original invocation binding")
	}
	if _, err := dir.Lstat("stop.json"); !os.IsNotExist(err) {
		t.Fatal("unbound cleanup published a stop")
	}
	if err := StopCapturedVerification(context.Background(), ticket, "old-verifier-launch"); err != nil {
		t.Fatalf("quorum gate blocked old-ticket stop: %v", err)
	}
	stopBefore, err := dir.ReadFile("stop.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := StopCapturedVerification(context.Background(), ticket, "old-verifier-launch"); err != nil {
		t.Fatal(err)
	}
	stopAfter, _ := dir.ReadFile("stop.json")
	requestAfter, _ := dir.ReadFile("request.json")
	launchAfter, _ := dir.ReadFile("launch.json")
	if !bytes.Equal(stopBefore, stopAfter) || !bytes.Equal(requestBefore, requestAfter) || !bytes.Equal(launchBefore, launchAfter) || !bytes.Equal(stateBefore, snapshotRead(t, statePath(*b))) {
		t.Fatal("stop replay changed original evidence")
	}
	if _, err := dir.Lstat("claim.json"); !os.IsNotExist(err) {
		t.Fatal("stop invented helper execution")
	}
	if len(s.Attempts) != 1 || len(s.Resolutions) != 0 {
		t.Fatal("legacy fixture lost its spent unresolved patch")
	}
	// No helper/process was claimed. This proves durable stop binding, not that
	// every possible historical descendant of an old ticket is inactive.
}

func TestActivationQuorumUnchangedLaterArchiveCannotReplaceBaseline(t *testing.T) {
	root, b, s, r := quorumPatchFixture(t, "[builder, other]")
	dir, err := os.OpenRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	// These are actual later archive bytes. Test the scope gate directly without
	// inventing an accepted first patch or a complete second attempt.
	a := Attempt{Before: r.After, BeforeArchive: r.AfterArchive}
	_, err = unchangedScope(context.Background(), *b, s, a, dir)
	requireQuorumRefusal(t, err, "unchanged before-source differs from the original activation quorum")
	// Both sides of a hypothetical later comparison may agree with each other;
	// original activation membership must still govern verifier selection.
	r.Before, r.BeforeArchive = r.After, r.AfterArchive
	r.Verifier = "other"
	_, err = checkCapturedActivationQuorum(context.Background(), *b, s, r)
	requireQuorumRefusal(t, err, "independent activation quorum")
}

func TestActivationQuorumRefusesMissingMalformedAndCorruptBaseline(t *testing.T) {
	for _, members := range []string{"missing-file", "[]", "[builder]", "[reviewer, other]", "[builder, reviewer, reviewer]", "[builder, ../other]", "corrupt-archive", "missing-archive"} {
		t.Run(members, func(t *testing.T) {
			root, b, s, _ := quorumPatchFixture(t, "[builder, reviewer]")
			if strings.HasSuffix(members, "-archive") {
				path := snapshotPath(snapshotDirectory(*b), s.BaselineArchive)
				if members == "missing-archive" {
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				} else {
					snapshotWrite(t, filepath.Dir(path), filepath.Base(path), []byte("corrupt"), 0600)
				}
			} else {
				name := "parley-deck/ideas/fixture/00-prompt.md"
				if members == "missing-file" {
					if err := os.Remove(filepath.Join(root, name)); err != nil {
						t.Fatal(err)
					}
				} else {
					snapshotWrite(t, root, name, []byte("---\nparticipants: "+members+"\n---\n"), 0600)
				}
				// Isolate archived membership decoding with a real, hash-verified
				// malformed-source archive; do not publish it as a valid state.
				source, err := Observe(context.Background(), root)
				if err != nil {
					t.Fatal(err)
				}
				ref, err := CaptureSnapshot(context.Background(), root, snapshotDirectory(*b), source)
				if err != nil {
					t.Fatal(err)
				}
				s.Policy.Baseline, s.BaselineArchive = source, ref
			}
			if got, err := activationQuorum(context.Background(), *b, s); err == nil || got != nil {
				t.Fatalf("unavailable/invalid activation quorum accepted: %v %v", got, err)
			}
		})
	}
}

func TestActivationQuorumRejectsWidenedParentRequestBeforePreviewHash(t *testing.T) {
	root, p := parentRecoveryFixture(t)
	b, err := budget.LoadCycleBinding(context.Background(), root, "fixture", budget.Fixup)
	if err != nil {
		t.Fatal(err)
	}
	s, err := Inspect(context.Background(), root, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(root, ".parley-runtime", "trajectory-verification", p.RunID)
	dir, err := os.OpenRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	var req HelperRequest
	if _, err := readReconciliationJSON(dir, filepath.Join(".parley-runtime", "trajectory-verification", p.RunID, "request.json"), &req); err != nil {
		t.Fatal(err)
	}
	req.Participants = append(req.Participants, "other")
	raw, err := canonical(req)
	if err != nil {
		t.Fatal(err)
	}
	snapshotWrite(t, base, "request.json", raw, 0600)
	name := "parley-deck/ideas/fixture/00-prompt.md"
	body := strings.Replace(string(snapshotRead(t, filepath.Join(root, name))), "[builder, reviewer]", "[builder, reviewer, other]", 1)
	snapshotWrite(t, root, name, []byte(body), 0600)
	// Defeat the old request-vs-current comparison on purpose. Only the frozen
	// baseline member check should reject, before any stale preview comparison.
	if err := checkReconciliationScope(dir, req); err != nil {
		t.Fatalf("fixture failed old scope gate: %v", err)
	}
	_, err = deriveParentEvidence(context.Background(), *b, *s, root, p.RunID)
	requireQuorumRefusal(t, err, "helper request differs from the original activation quorum")
}

func TestActivationQuorumValidUnchangedReconciliationStillReplays(t *testing.T) {
	root, b := unchangedFixture(t, "true")
	unchangedProcess(t, root, b, "exit 0", false)
	p, err := PreviewUnchanged(context.Background(), root, "fixture", 1)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		applied, err := ReconcileUnchanged(context.Background(), root, "fixture", 1, p.SHA256())
		if err != nil || !sameJSON(applied, p) {
			t.Fatalf("valid unchanged replay %d: %v", i, err)
		}
	}
	s, err := Inspect(context.Background(), root, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Resolutions) != 1 || s.Resolutions[0].Preview.Assessment.Outcome != Inconclusive {
		t.Fatal("quorum pinning changed inconclusive accounting")
	}
}

func TestActivationQuorumDoesNotReinterpretExistingCheckSyntax(t *testing.T) {
	root, b, s, _ := quorumPatchFixture(t, "[builder, reviewer]")
	for _, checks := range []string{"", "checks: 'legacy scalar command'\n", "checks: []\n"} {
		body := "---\nparticipants: [builder, reviewer]\n" + checks + "---\n"
		snapshotWrite(t, root, "parley-deck/ideas/fixture/00-prompt.md", []byte(body), 0600)
		source, err := Observe(context.Background(), root)
		if err != nil {
			t.Fatal(err)
		}
		ref, err := CaptureSnapshot(context.Background(), root, snapshotDirectory(*b), source)
		if err != nil {
			t.Fatal(err)
		}
		candidate := s
		candidate.Policy.Baseline = source
		candidate.BaselineArchive = ref
		members, err := activationQuorum(context.Background(), *b, candidate)
		if err != nil || !slices.Equal(members, []string{"builder", "reviewer"}) {
			t.Fatalf("quorum decoder imposed new check syntax %q: %v", checks, err)
		}
	}
}
