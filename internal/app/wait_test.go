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

func TestWaitBlockingEscalationExitsFour(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	inbox := filepath.Join(root, protocol.DeckDir, "inbox")
	if err := os.MkdirAll(inbox, 0o755); err != nil {
		t.Fatal(err)
	}
	esc := filepath.Join(inbox, "claude-1-to-user_wait-idea_blocker.md")
	if err := os.WriteFile(esc, []byte("---\nfrom: claude-1\nto: user\nblocking: yes\n---\n\n## Question\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "any", "--timeout", "5s")
	if code != 4 {
		t.Fatalf("unanswered to-user escalation must exit 4 immediately, got %d; err=%q", code, errOut)
	}
	if !strings.Contains(errOut, "to-user") {
		t.Errorf("exit must name the escalation source; got %q", errOut)
	}
}

func TestWaitDriverErrorEventExitsFour(t *testing.T) {
	withWaitPoll(t, 5*time.Millisecond)
	root, _ := seedWaitIdea(t, []string{"claude-1", "kimi-1"})
	runsDir := filepath.Join(root, protocol.DeckDir, "runs", "20260923T000000.000000000Z")
	os.MkdirAll(runsDir, 0o755)
	manifest := `{"schema_version": 1, "run_id": "20260923T000000.000000000Z", "idea_slug": "wait-idea", "status": "running"}`
	if err := os.WriteFile(filepath.Join(runsDir, "run.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	events := store.New(runsDir)
	events.Append(store.Event{Type: "run.created", Data: map[string]any{"idea": "wait-idea"}})
	events.Append(store.Event{Type: "driver.error", Data: map[string]any{"idea": "wait-idea", "error": "context canceled"}})
	code, _, errOut := runWaitFor(t, "--dir", root, "--idea", "wait-idea", "--for", "any", "--timeout", "5s")
	if code != 4 {
		t.Fatalf("driver.error event must exit 4 immediately, got %d; err=%q", code, errOut)
	}
	if !strings.Contains(errOut, "context canceled") {
		t.Errorf("driver error detail must be carried verbatim; got %q", errOut)
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
