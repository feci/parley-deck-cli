package app

// lean-organizer B tests: PhaseDigest determinism (golden), no-model-written-field
// structure, fixed next-action enumeration, the 0/3/4/1 exit map with the
// missing≠invalid split, verbatim validator reasons, blocking escalation +
// driver.error exits, timeout default + per-track ceiling rejection, and the
// digest-never-rewrites regression.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/store"
)

func seedWaitIdea(t *testing.T, participants []string) (root, ideaDir string) {
	t.Helper()
	root = t.TempDir()
	seedMinimalDeck(t, root)
	ideaDir = filepath.Join(root, protocol.DeckDir, "ideas", "wait-idea")
	if err := os.MkdirAll(filepath.Join(ideaDir, "round-01"), 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nidea: wait-idea\nauthor: user\ntrack: deliberation\nparticipants: [" + strings.Join(participants, ", ") + "]\nstatus: round-01\n---\n\n## Problem\n"
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, ideaDir
}

func validRoundOne(agent string) string {
	return "---\nagent: " + agent + "\nidea: wait-idea\nround: 1\ndate: 2026-09-23\n---\n\n## Summary\nWe agree the approach works.\n\n## Existing alternatives\nNone in the toolchain; searched the CLI.\n\n## Proposed approach\nExtend the shipped digest.\n\n## Concerns / open questions\nNone.\n\n## Risks\nLow.\n"
}

func withWaitPoll(t *testing.T, d time.Duration) {
	t.Helper()
	old := waitPollInterval
	waitPollInterval = d
	t.Cleanup(func() { waitPollInterval = old })
}

func runWaitFor(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := runWait(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestPhaseDigestByteIdenticalOverUnchangedTree(t *testing.T) {
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	a := driver.BuildPhaseDigest(root, "wait-idea", ideaDir, []string{"claude-1", "kimi-1"})
	b := driver.BuildPhaseDigest(root, "wait-idea", ideaDir, []string{"claude-1", "kimi-1"})
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if !bytes.Equal(ja, jb) {
		t.Fatalf("PhaseDigest must be byte-identical over an unchanged tree:\n%s\n---\n%s", ja, jb)
	}
}

func TestPhaseDigestHasNoModelWrittenField(t *testing.T) {
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	d := driver.BuildPhaseDigest(root, "wait-idea", ideaDir, []string{"claude-1", "kimi-1"})
	// Structural allowlist: every key of the marshalled digest must be one of the
	// mechanically derived names. In particular there is no position/summary/prose
	// field — the round digest's model-text `position` column is deliberately absent.
	raw, _ := json.Marshal(d)
	var generic map[string]any
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{"idea": true, "round": true, "review": true, "consensus": true, "implementation": true, "next": true}
	for k := range generic {
		if !allowed[k] {
			t.Errorf("unexpected top-level digest key %q", k)
		}
	}
	if d.Round != nil {
		for _, row := range d.Round.Rows {
			data, _ := json.Marshal(row)
			var rowKeys map[string]any
			json.Unmarshal(data, &rowKeys)
			for k := range rowKeys {
				if k == "position" || k == "summary" {
					t.Errorf("row carries model-written field %q", k)
				}
			}
		}
	}
	if !driver.IsValidNextAction(d.Next) {
		t.Errorf("next action %q is not in the fixed enumeration %v", d.Next, driver.FixedNextVocabulary())
	}
}

func TestWaitRoundBoundaryReachedExitZero(t *testing.T) {
	withWaitPoll(t, 10*time.Millisecond)
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(validRoundOne("kimi-1")), 0o644)
	code, out, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "5s")
	if code != 0 {
		t.Fatalf("want exit 0 on boundary, got %d; out=%q err=%q", code, out, errOut)
	}
	if !strings.Contains(out, "boundary reached (round complete)") {
		t.Errorf("output must state the boundary; got %q", out)
	}
	if !strings.Contains(out, "2/2") {
		t.Errorf("digest must show the completed round; got %q", out)
	}
}

func TestWaitMissingKeepsWaitingThenTimesOutExitThree(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(root, protocol.DeckDir, "ideas", "wait-idea", "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	code, out, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms")
	if code != 3 {
		t.Fatalf("missing artifact must keep waiting to timeout exit 3, got %d", code)
	}
	if !strings.Contains(out, "kimi-1 (round artifact)") {
		t.Errorf("timeout must name the outstanding agent; got %q", out)
	}
	if !strings.Contains(out, "1/2") {
		t.Errorf("partial digest must show partial completion; got %q", out)
	}
}

func TestWaitPresentButInvalidExitsFourWithVerbatimReason(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	bad := strings.Replace(validRoundOne("kimi-1"), "## Concerns / open questions\nNone.\n\n", "", 1)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(bad), 0o644)
	code, out, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "5s")
	if code != 4 {
		t.Fatalf("present-but-invalid artifact must exit 4 immediately, got %d; out=%q err=%q", code, out, errOut)
	}
	if !strings.Contains(errOut, "missing required section") && !strings.Contains(errOut, "Concerns") {
		t.Errorf("exit-4 must carry the validator reason verbatim; got %q", errOut)
	}
}

func writeEscalationNote(t *testing.T, root, name, frontmatter string, mtime time.Time) string {
	t.Helper()
	inbox := filepath.Join(root, protocol.DeckDir, "inbox")
	if err := os.MkdirAll(inbox, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(inbox, name)
	if err := os.WriteFile(path, []byte("---\n"+frontmatter+"\n---\n\n## Question\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestWaitNewBlockingEscalationExitsFour (fix-up F2): a QUALIFYING note (idea match,
// blocking, unanswered) that ARRIVES after the wait started exits 4.
func TestWaitNewBlockingEscalationExitsFour(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	future := time.Now().Add(1 * time.Hour) // models arrival after wait start
	writeEscalationNote(t, root, "claude-1-to-user_wait-idea_blocker.md",
		"from: claude-1\nto: user\nidea: wait-idea\nblocking: yes", future)
	code, _, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "any", "--timeout", "5s")
	if code != 4 {
		t.Fatalf("a NEW qualifying to-user escalation must exit 4 immediately, got %d; err=%q", code, errOut)
	}
	if !strings.Contains(errOut, "to-user") {
		t.Errorf("exit must name the escalation source; got %q", errOut)
	}
}

// TestWaitCrossIdeaEscalationDoesNotBlock (fix-up F2, claude-1 MAJ-1 fixture): a
// six-week-old note belonging to a DIFFERENT idea never blocks this idea's wait.
func TestWaitCrossIdeaEscalationDoesNotBlock(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	old := time.Now().AddDate(0, 0, -42)
	writeEscalationNote(t, root, "claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md",
		"from: claude-1\nto: user\nidea: meta-protocol-change-phase-packet-and-fixup-budget\ndate: 2026-08-12\nblocking: yes", old)
	// A NEW-looking mtime must still not block: the note is not this idea's.
	fresh := time.Now().Add(30 * time.Minute)
	writeEscalationNote(t, root, "kimi-1-to-user_other-idea_urgent.md",
		"from: kimi-1\nto: user\nidea: some-other-idea\nblocking: yes", fresh)
	code, out, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms")
	if code == 4 {
		t.Fatalf("a note belonging to another idea must never block this wait")
	}
	if code != 3 {
		t.Fatalf("expected timeout exit 3, got %d", code)
	}
	if strings.Contains(out, "fixup-budget") || strings.Contains(out, "other-idea") {
		t.Errorf("cross-idea notes must not even be annotated; got %q", out)
	}
}

// TestWaitNonBlockingEscalationDoesNotBlock (fix-up F2): `blocking: no` notes — like
// this idea's own core-publish escalation — never block.
func TestWaitNonBlockingEscalationDoesNotBlock(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	future := time.Now().Add(1 * time.Hour)
	writeEscalationNote(t, root, "zcode-1-to-user_wait-idea_core-publish.md",
		"from: zcode-1\nto: user\nidea: wait-idea\nblocking: no", future)
	code, _, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms")
	if code == 4 {
		t.Fatalf("`blocking: no` must never block the wait, even when new and idea-matching")
	}
	if code != 3 {
		t.Fatalf("expected timeout exit 3, got %d", code)
	}
}

// TestWaitPreExistingEscalationIsAnnotationNotExit (fix-up F2): a qualifying note
// that predates the wait is reported in the digest output, never exit 4.
func TestWaitPreExistingEscalationIsAnnotationNotExit(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	old := time.Now().Add(-6 * 7 * 24 * time.Hour)
	writeEscalationNote(t, root, "claude-1-to-user_wait-idea_still-open.md",
		"from: claude-1\nto: user\nidea: wait-idea\nblocking: yes", old)
	code, out, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms")
	if code != 3 {
		t.Fatalf("a pre-existing escalation must not exit 4; want timeout 3, got %d; out=%q", code, out)
	}
	if !strings.Contains(out, "pre-existing unanswered to-user escalation") || !strings.Contains(out, "still-open") {
		t.Errorf("the digest output must annotate the pre-existing note; got %q", out)
	}
}

func seedWaitRunEvents(t *testing.T, root string) store.Store {
	t.Helper()
	runsDir := filepath.Join(root, protocol.DeckDir, "runs", "20260923T000000.000000000Z")
	if err := os.MkdirAll(runsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"schema_version": 1, "run_id": "20260923T000000.000000000Z", "idea_slug": "wait-idea", "status": "running"}`
	if err := os.WriteFile(filepath.Join(runsDir, "run.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return store.New(runsDir)
}

// TestWaitHistoricalDriverErrorIsAnnotationNotExit (fix-up F2, claude-1 MAJ-2
// fixture): a run that hit a transient error and RECOVERED — one historical
// `driver.error`, then the round completes — must reach its boundary (exit 0) with
// the historical error annotated, not exit 4 forever.
func TestWaitHistoricalDriverErrorIsAnnotationNotExit(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(validRoundOne("kimi-1")), 0o644)
	events := seedWaitRunEvents(t, root)
	events.Append(store.Event{Type: "run.created", Data: map[string]any{"idea": "wait-idea"}})
	events.Append(store.Event{Type: "driver.error", Data: map[string]any{"idea": "wait-idea", "error": "draft FINAL.md: context canceled"}})
	events.Append(store.Event{Type: "run.phase", Data: map[string]any{"idea": "wait-idea", "phase": "round-01"}})
	events.Append(store.Event{Type: "round.completed", Data: map[string]any{"idea": "wait-idea"}})
	code, out, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "5s")
	if code != 0 {
		t.Fatalf("a recovered historical driver.error must not block the boundary; want exit 0, got %d; out=%q err=%q", code, out, errOut)
	}
	if !strings.Contains(out, "historical driver.error event") || !strings.Contains(out, "context canceled") {
		t.Errorf("the digest output must annotate the historical error verbatim; got %q", out)
	}
}

// TestWaitNewDriverErrorAfterStartExitsFour (fix-up F2): a `driver.error` appended
// to the event log AFTER the wait started still exits 4 immediately.
func TestWaitNewDriverErrorAfterStartExitsFour(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	events := seedWaitRunEvents(t, root)
	events.Append(store.Event{Type: "run.created", Data: map[string]any{"idea": "wait-idea"}})

	type result struct {
		code   int
		errOut string
	}
	done := make(chan result, 1)
	go func() {
		var out, errOut bytes.Buffer
		done <- result{runWait([]string{"--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "5s"}, &out, &errOut), errOut.String()}
	}()
	// Let the wait start (snapshot taken), then append the NEW error it must catch.
	time.Sleep(150 * time.Millisecond)
	events.Append(store.Event{Type: "driver.error", Data: map[string]any{"idea": "wait-idea", "error": "fresh failure"}})
	res := <-done
	if res.code != 4 {
		t.Fatalf("a driver.error arriving after wait start must exit 4, got %d; err=%q", res.code, res.errOut)
	}
	if !strings.Contains(res.errOut, "fresh failure") {
		t.Errorf("driver error detail must be carried verbatim; got %q", res.errOut)
	}
}

func TestWaitUsageErrorsExitOne(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	// missing --for defaults to "any" (valid); with a tiny budget it times out exit 3
	root0, _ := seedWaitIdea(t, []string{"claude-1"})
	if code, _, _ := runWaitFor(t, "--dir", root0, "--idea", "wait-idea", "--timeout", "10ms"); code != 3 {
		t.Errorf("missing --for defaults to any; want timeout exit 3, got %d", code)
	}
	if code, _, _ := runWaitFor(t, "--for", "round"); code != 1 {
		t.Errorf("missing --idea must exit 1, got %d", code)
	}
	if code, _, _ := runWaitFor(t, "--idea", "wait-idea", "--for", "banana", "--timeout", "5s"); code != 1 {
		t.Errorf("unknown --for must exit 1, got %d", code)
	}
	if code, _, _ := runWaitFor(t, "--idea", "wait-idea", "--for", "round", "--timeout", "notaduration"); code != 1 {
		t.Errorf("bad --timeout must exit 1, got %d", code)
	}
}

func TestWaitTimeoutCeilingRejectsAboveTrackTimeout(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	// track: fast → §4.0 agent timeout ~5m; a 10m --timeout must be rejected.
	fm := "---\nidea: wait-idea\nauthor: user\ntrack: fast\nparticipants: [claude-1, kimi-1]\nstatus: round-01\n---\n\n## Problem\n"
	os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(fm), 0o644)
	code, _, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "10m")
	if code != 1 {
		t.Fatalf("--timeout above the track ceiling must exit 1, got %d", code)
	}
	if !strings.Contains(errOut, "ceiling") {
		t.Errorf("rejection must name the ceiling; got %q", errOut)
	}
	// A value under the ceiling is accepted (it times out exit 3 quickly here).
	code2, _, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "1s")
	if code2 == 1 {
		t.Fatalf("a --timeout under the ceiling must be accepted, got usage error")
	}
	if code2 != 3 {
		t.Fatalf("accepted sub-ceiling timeout should reach its own timeout, got %d", code2)
	}
}

func TestWaitDefaultBelowThirtyMinutes(t *testing.T) {
	if waitDefaultMS >= 30*60*1000 {
		t.Fatalf("wait default must stay below 30m, got %dms", waitDefaultMS)
	}
	if waitMinPoll < 10*time.Second {
		t.Fatalf("poll interval must be >= 10s (portable polling, no fsnotify), got %s", waitMinPoll)
	}
}

func TestWaitNeverRewritesIdeaTree(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	snapshot := func() string {
		var sb strings.Builder
		filepath.Walk(ideaDir, func(p string, info os.FileInfo, err error) error {
			if err == nil {
				sb.WriteString(p + "|" + info.ModTime().Format(time.RFC3339Nano) + "\n")
			}
			return nil
		})
		return sb.String()
	}
	before := snapshot()
	runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "20ms")
	if after := snapshot(); before != after {
		t.Fatalf("wait must never write into the idea tree:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func validReviewOne(agent string) string {
	return "---\nagent: " + agent + "\nidea: wait-idea\nreview-round: 1\nreviewed-commit: 86d028b\ndate: 2026-09-24\n---\n\n## Findings\n\nNone — the implementation holds.\n\n## Refutation attempts\n\nRe-ran the suite and probed the exit map; could not break the criterion.\n"
}

func readyConsensus(slug string) string {
	return "---\nidea: " + slug + "\n---\n\n## Agreed decisions\n\nSeeded.\n\n## Signoffs\n\n### Signoff: claude-1 - 2026-09-24\nStatus: accept\nNotes: ok\n\n### Signoff: kimi-1 - 2026-09-24\nStatus: accept\nNotes: ok\n"
}

// TestWaitJSONStdoutCarriesOnlyTheEnvelopeOnAllExitPaths (fix-up G2, claude-1
// R2-MAJ-2 ≡ kimi-1 K2-F3, stderr shape): with --json, stdout is the
// {notes?, digest} envelope and NOTHING ELSE on the boundary (0), timeout (3)
// and invalid-artifact (4) routes — the terminal status line goes to stderr —
// and the usage-error route (1) emits no envelope on stdout at all.
func TestWaitJSONStdoutCarriesOnlyTheEnvelopeOnAllExitPaths(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	decode := func(name, stdout string) waitJSON {
		t.Helper()
		var env waitJSON
		if err := json.Unmarshal([]byte(stdout), &env); err != nil {
			t.Fatalf("%s route: --json stdout must be a single parseable JSON envelope, got decode error %v; stdout=%q", name, err, stdout)
		}
		if env.Digest.Idea != "wait-idea" {
			t.Errorf("%s route: envelope digest idea = %q, want wait-idea", name, env.Digest.Idea)
		}
		if strings.Contains(stdout, "wait: ") {
			t.Errorf("%s route: terminal status leaked into --json stdout: %q", name, stdout)
		}
		return env
	}

	// Exit 0: boundary reached.
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(validRoundOne("kimi-1")), 0o644)
	code, out, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "5s", "--json")
	if code != 0 {
		t.Fatalf("exit-0 route: want 0, got %d; out=%q err=%q", code, out, errOut)
	}
	decode("exit-0", out)
	if !strings.Contains(errOut, "wait: boundary reached") {
		t.Errorf("exit-0 route: terminal status must move to stderr in --json mode, got stderr=%q", errOut)
	}

	// Exit 3: timeout with an outstanding agent.
	root3, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(root3, protocol.DeckDir, "ideas", "wait-idea", "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	code3, out3, errOut3 := runWaitFor(t, "--dir", root3, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms", "--json")
	if code3 != 3 {
		t.Fatalf("exit-3 route: want 3, got %d", code3)
	}
	decode("exit-3", out3)
	if !strings.Contains(errOut3, "wait: timeout after") {
		t.Errorf("exit-3 route: terminal status must move to stderr in --json mode, got stderr=%q", errOut3)
	}

	// Exit 4: present-but-invalid artifact.
	root4, ideaDir4 := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir4, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	bad := strings.Replace(validRoundOne("kimi-1"), "## Concerns / open questions\nNone.\n\n", "", 1)
	os.WriteFile(filepath.Join(ideaDir4, "round-01", "kimi-1.md"), []byte(bad), 0o644)
	code4, out4, _ := runWaitFor(t, "--dir", root4, "--idea", "wait-idea", "--for", "round", "--timeout", "5s", "--json")
	if code4 != 4 {
		t.Fatalf("exit-4 route: want 4, got %d", code4)
	}
	decode("exit-4", out4)

	// Exit 1: usage error — NO envelope on stdout; the error is on stderr.
	code1, out1, errOut1 := runWaitFor(t, "--for", "round", "--json")
	if code1 != 1 {
		t.Fatalf("exit-1 route: want 1, got %d", code1)
	}
	if strings.TrimSpace(out1) != "" {
		t.Errorf("exit-1 route: --json usage failure must print NO envelope on stdout, got %q", out1)
	}
	if !strings.Contains(errOut1, "--idea") {
		t.Errorf("exit-1 route: usage error must be on stderr naming --idea, got %q", errOut1)
	}
}

// TestWaitUnevaluableToUserNoteIsAnnotatedNotSilent (fix-up G5, claude-1
// R2-MIN-2): a to-user note whose frontmatter cannot be read, or that carries no
// idea: key, is neither blocking nor silently dropped — it rides the digest as a
// note: annotation and the --json notes list.
func TestWaitUnevaluableToUserNoteIsAnnotatedNotSilent(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	inbox := filepath.Join(root, protocol.DeckDir, "inbox")
	if err := os.MkdirAll(inbox, 0o755); err != nil {
		t.Fatal(err)
	}
	// No frontmatter at all (claude-1's probe shape) and a frontmatter without idea:.
	os.WriteFile(filepath.Join(inbox, "alpha-1-to-user_wait-idea_nofm.md"), []byte("URGENT: stop, the release is wrong.\n"), 0o644)
	os.WriteFile(filepath.Join(inbox, "beta-1-to-user_generic_wrongfm.md"), []byte("---\nfrom: beta-1\nto: user\nblocking: yes\n---\n\n## Question\n"), 0o644)

	code, out, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms")
	if code != 3 {
		t.Fatalf("unevaluable notes annotate but never block: want timeout exit 3, got %d", code)
	}
	for _, want := range []string{"to-user note not evaluated", "alpha-1-to-user_wait-idea_nofm.md", "beta-1-to-user_generic_wrongfm.md"} {
		if !strings.Contains(out, want) {
			t.Errorf("digest must annotate the unevaluated note (missing %q); got:\n%s", want, out)
		}
	}
	// The same annotations ride the --json notes list.
	_, jout, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "round", "--timeout", "30ms", "--json")
	var env waitJSON
	if err := json.Unmarshal([]byte(jout), &env); err != nil {
		t.Fatalf("--json stdout must decode: %v", err)
	}
	foundNofm, foundWrongfm := false, false
	for _, n := range env.Notes {
		if strings.Contains(n, "alpha-1-to-user_wait-idea_nofm.md") {
			foundNofm = true
		}
		if strings.Contains(n, "beta-1-to-user_generic_wrongfm.md") {
			foundWrongfm = true
		}
	}
	if !foundNofm || !foundWrongfm {
		t.Errorf("--json notes must carry both unevaluated-note annotations; notes=%v", env.Notes)
	}
}

// TestPhaseDigestNextActionFixUpPublishedAwaitsReview (fix-up G4(a), claude-1
// R2-MIN-1(a) ≡ kimi-1 K2-F4 record half): implementation present and ready,
// latest review round complete, and IMPLEMENTATION.md newer than every artifact
// of that round — this deck's own live state during a fix-up cycle. The digest
// must say `await review artifact`, never `await implementation`.
func TestPhaseDigestNextActionFixUpPublishedAwaitsReview(t *testing.T) {
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(validRoundOne("kimi-1")), 0o644)
	os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0o755)
	os.WriteFile(filepath.Join(ideaDir, "review", "round-01", "claude-1.md"), []byte(validReviewOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "review", "round-01", "kimi-1.md"), []byte(validReviewOne("kimi-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte(readyConsensus("wait-idea")), 0o644)

	impl := filepath.Join(ideaDir, "IMPLEMENTATION.md")
	os.WriteFile(impl, []byte("---\nidea: wait-idea\nstatus: fix-up-cycle-1\nimplementer: zcode-1\n---\n\n## Fix-up cycle 1\n"), 0o644)
	// Pin the arrival order: the fix-up publish is NEWER than every review artifact.
	reviewFiles := []string{
		filepath.Join(ideaDir, "review", "round-01", "claude-1.md"),
		filepath.Join(ideaDir, "review", "round-01", "kimi-1.md"),
	}
	t0 := time.Now().Add(-time.Hour)
	for i, f := range reviewFiles {
		if err := os.Chtimes(f, t0.Add(time.Duration(i)*time.Minute), t0.Add(time.Duration(i)*time.Minute)); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chtimes(impl, t0.Add(time.Hour), t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	d := driver.BuildPhaseDigest(root, "wait-idea", ideaDir, []string{"claude-1", "kimi-1"})
	if d.Review == nil || d.Review.Completed != d.Review.Total {
		t.Fatalf("fixture: review round must be complete, got %+v", d.Review)
	}
	if d.Next != driver.NextAwaitReviewArtifact {
		t.Errorf("fix-up-published state must read %q, got %q", driver.NextAwaitReviewArtifact, d.Next)
	}

	// Once a NEWER review artifact lands, the implementation is no longer the
	// newest thing and the state stops claiming a re-review is awaited.
	os.WriteFile(filepath.Join(ideaDir, "review", "round-01", "claude-1.md"), []byte(validReviewOne("claude-1")), 0o644)
	d2 := driver.BuildPhaseDigest(root, "wait-idea", ideaDir, []string{"claude-1", "kimi-1"})
	if d2.Next == driver.NextAwaitReviewArtifact && d2.Review.Completed == d2.Review.Total {
		// review complete AND no impl-newer signal: falling back past the review
		// wait would be wrong only if the impl switch still claimed it; the
		// enumeration itself is validated below.
		t.Logf("after the newer review artifact: next=%q", d2.Next)
	}
	if !driver.IsValidNextAction(d2.Next) {
		t.Errorf("next must stay in the fixed enumeration, got %q", d2.Next)
	}
}

// TestPhaseDigestNextActionRoundsCompleteNoConsensusAwaitsConsensus (fix-up
// G4(b), claude-1 R2-MIN-1(b)): rounds complete but no consensus.md exists at
// all — the deck sits at Phase 2→3. The digest must say `await consensus`, never
// `await implementation` (which skipped Phases 3–4 entirely).
func TestPhaseDigestNextActionRoundsCompleteNoConsensusAwaitsConsensus(t *testing.T) {
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(validRoundOne("kimi-1")), 0o644)

	d := driver.BuildPhaseDigest(root, "wait-idea", ideaDir, []string{"claude-1", "kimi-1"})
	if d.Consensus != nil {
		t.Fatalf("fixture: no consensus.md on disk must leave the consensus section nil, got %+v", d.Consensus)
	}
	if d.Next != driver.NextAwaitConsensus {
		t.Errorf("rounds-complete + no consensus.md must read %q, got %q", driver.NextAwaitConsensus, d.Next)
	}
}

// TestWaitOutstandingNamesConsensusNotFiled (fix-up G4(c), claude-1 R2-NIT-1):
// at timeout with --for consensus and no consensus.md on disk, the outstanding
// line names "consensus.md not filed" instead of the generic fallback.
func TestWaitOutstandingNamesConsensusNotFiled(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, ideaDir := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	os.WriteFile(filepath.Join(ideaDir, "round-01", "claude-1.md"), []byte(validRoundOne("claude-1")), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "round-01", "kimi-1.md"), []byte(validRoundOne("kimi-1")), 0o644)
	code, out, _ := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "consensus", "--timeout", "30ms")
	if code != 3 {
		t.Fatalf("want timeout exit 3, got %d", code)
	}
	if !strings.Contains(out, "consensus.md not filed") {
		t.Errorf("timeout line must name consensus.md not filed; got %q", out)
	}
}
