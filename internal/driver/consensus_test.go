package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/store"
)

// fakeConsensus is an injected ConsensusOps that records calls and returns a
// scripted triage, so the consensus gate is tested without live agents.
type fakeConsensus struct {
	statusSeq    []string // triage returned per Status() call (last value repeats)
	idx          int
	missing      []string
	statusErr    error
	calls        []string
	onDraft      func()
	onDraftFinal func()
}

func (f *fakeConsensus) Status() (consensus.Summary, error) {
	if f.statusErr != nil {
		return consensus.Summary{}, f.statusErr
	}
	tri := ""
	if len(f.statusSeq) > 0 {
		if f.idx < len(f.statusSeq) {
			tri = f.statusSeq[f.idx]
			f.idx++
		} else {
			tri = f.statusSeq[len(f.statusSeq)-1]
		}
	}
	return consensus.Summary{Triage: tri, Missing: f.missing}, nil
}
func (f *fakeConsensus) Draft(ctx context.Context) error {
	f.calls = append(f.calls, "draft")
	if f.onDraft != nil {
		f.onDraft()
	}
	return nil
}
func (f *fakeConsensus) RequestSignoffs(ctx context.Context, missing []string) error {
	f.calls = append(f.calls, "signoffs")
	return nil
}
func (f *fakeConsensus) DraftFinal(ctx context.Context) error {
	f.calls = append(f.calls, "final")
	if f.onDraftFinal != nil {
		f.onDraftFinal()
	}
	return nil
}
func (f *fakeConsensus) Reopen(ctx context.Context, reason string) error {
	f.calls = append(f.calls, "reopen")
	return nil
}

const validFinal = "---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\nThe driver advances the deliberation through a disk-derived cursor and an\nordered switch over the round and consensus phases, gated by transport.\nThe consensus gate drafts, collects signoffs, and finalizes autonomously.\nThis paragraph is padded well beyond two hundred and fifty bytes so the\nnon-scaffold length check passes comfortably during the unit test run.\n\n\n## Purpose / user-visible outcome\nN/A\n\n## Context & orientation\nN/A\n\n## Observable acceptance criteria\nN/A\n\n## Idempotence & recovery\nN/A\n\n## Known risks / de-risking\nN/A\n\n## References\nN/A\n"

func writeConsensusDoc(t *testing.T, ideaDir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newConsensusDriver(ideaDir, runDir string, parts []string, fc ConsensusOps, fr RoundRunner) *Driver {
	return New(Config{
		IdeaDir:           ideaDir,
		IdeaSlug:          "demo",
		Participants:      parts,
		RunDir:            runDir,
		Root:              filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))),
		Events:            store.New(runDir),
		CrossReviewRounds: 1,
		MaxRounds:         4,
		Auto:              true,
		Consensus:         fc,
	}, fr)
}

func TestRoundBudgetSpentDraftsConsensus(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	writeAll(t, ideaDir, 2, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	appendEvent(t, runDir, "round.completed", "round-02")
	fc := &fakeConsensus{statusSeq: []string{consensus.TriagePartial}, onDraft: func() { writeConsensusDoc(t, ideaDir, "consensus") }}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionConsensusDrafted {
		t.Fatalf("action=%s, want consensus-drafted", action)
	}
	if c.Phase != PhaseConsensus {
		t.Fatalf("phase=%s, want consensus", c.Phase)
	}
	if len(fc.calls) != 1 || fc.calls[0] != "draft" {
		t.Fatalf("calls=%v, want [draft]", fc.calls)
	}
}

func TestConsensusReadyDraftsFinalAndAdvances(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	fc := &fakeConsensus{
		statusSeq:    []string{consensus.TriageReady},
		onDraftFinal: func() { os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte(validFinal), 0o644) },
	}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionFinalized {
		t.Fatalf("action=%s, want finalized", action)
	}
	if c.Phase != PhaseFinal {
		t.Fatalf("phase=%s, want final", c.Phase)
	}
	if len(fc.calls) != 1 || fc.calls[0] != "final" {
		t.Fatalf("calls=%v, want [final]", fc.calls)
	}
	wantLastRunPhase(t, runDir, "finalized", PhaseFinal, PhaseConsensus)
}

func TestConsensusReadyButFinalScaffoldEscalates(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	scaffold := "---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\n<fill in the agreed plan here>\n\n\n## Purpose / user-visible outcome\nN/A\n\n## Context & orientation\nN/A\n\n## Observable acceptance criteria\nN/A\n\n## Idempotence & recovery\nN/A\n\n## Known risks / de-risking\nN/A\n\n## References\nN/A\n"
	fc := &fakeConsensus{
		statusSeq:    []string{consensus.TriageReady},
		onDraftFinal: func() { os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte(scaffold), 0o644) },
	}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, _, err := d.Advance(context.Background())
	if err == nil {
		t.Fatal("expected an error for a scaffold FINAL.md")
	}
	if action != ActionEscalated {
		t.Fatalf("action=%s, want escalated", action)
	}
	// AF1: a failed/scaffold draft must NOT commit the idea to status=final.
	if got := readIdeaStatus(ideaDir); got == "final" {
		t.Fatalf("idea status must not be final after a scaffold draft, got %q", got)
	}
}

func TestConsensusReadyRevalidatesExistingScaffoldFinal(t *testing.T) {
	// A stale scaffold FINAL.md from a prior failed draft must be re-drafted, not
	// blindly accepted, and status is committed only after the content validates
	// (AF1 — the failure mode where status=final was set before the content).
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	if err := os.WriteFile(filepath.Join(ideaDir, "FINAL.md"),
		[]byte("---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\n<todo>\n\n\n## Purpose / user-visible outcome\nN/A\n\n## Context & orientation\nN/A\n\n## Observable acceptance criteria\nN/A\n\n## Idempotence & recovery\nN/A\n\n## Known risks / de-risking\nN/A\n\n## References\nN/A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fc := &fakeConsensus{
		statusSeq:    []string{consensus.TriageReady},
		onDraftFinal: func() { os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte(validFinal), 0o644) },
	}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionFinalized {
		t.Fatalf("action=%s, want finalized after re-draft", action)
	}
	if len(fc.calls) != 1 || fc.calls[0] != "final" {
		t.Fatalf("expected a re-draft of the stale scaffold, calls=%v", fc.calls)
	}
	if got := readIdeaStatus(ideaDir); got != "final" {
		t.Fatalf("idea status=%q, want final after a valid re-draft", got)
	}
}

func TestConsensusPartialRequestsSignoffsThenProceeds(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	// First Status: Partial (→ request signoffs); after request: Ready.
	fc := &fakeConsensus{statusSeq: []string{consensus.TriagePartial, consensus.TriageReady}, missing: []string{"claude"}}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionSignoffsRequested {
		t.Fatalf("action=%s, want signoffs-requested", action)
	}
	if len(fc.calls) != 1 || fc.calls[0] != "signoffs" {
		t.Fatalf("calls=%v, want [signoffs]", fc.calls)
	}
}

func TestConsensusPartialStillMissingEscalates(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	fc := &fakeConsensus{statusSeq: []string{consensus.TriagePartial}, missing: []string{"claude"}} // always Partial
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, _, err := d.Advance(context.Background())
	if err == nil {
		t.Fatal("expected escalation when signoffs remain missing after request")
	}
	if action != ActionEscalated {
		t.Fatalf("action=%s, want escalated", action)
	}
}

func TestConsensusBlockedReopensRound(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	writeAll(t, ideaDir, 2, parts)
	writeConsensusDoc(t, ideaDir, "consensus")
	fc := &fakeConsensus{statusSeq: []string{consensus.TriageBlocked}}
	fr := &fakeRunner{}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, fr)

	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionReopened {
		t.Fatalf("action=%s, want reopened", action)
	}
	if len(fc.calls) != 1 || fc.calls[0] != "reopen" {
		t.Fatalf("consensus calls=%v, want [reopen]", fc.calls)
	}
	if len(fr.calls) != 1 || fr.calls[0] != 3 {
		t.Fatalf("RunRound calls=%v, want [3] (reopened round-03)", fr.calls)
	}
	if c.Phase != PhaseRound || c.CurrentRound != 3 {
		t.Fatalf("got phase=%s round=%d, want round/3", c.Phase, c.CurrentRound)
	}
	if !fileExists(filepath.Join(ideaDir, "consensus.md.bak")) {
		t.Fatal("stale consensus.md should be invalidated to consensus.md.bak")
	}
	if fileExists(filepath.Join(ideaDir, "consensus.md")) {
		t.Fatal("consensus.md should be gone after reopen invalidation")
	}
	wantLastRunPhase(t, runDir, "reopened", PhaseRound, PhaseConsensus)
}

func TestConsensusBlockedMaxRoundsEscalates(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	// Already at the MaxRounds ceiling: rounds 1..5 exist (1 + MaxRounds=4 → 5).
	for r := 1; r <= 5; r++ {
		writeAll(t, ideaDir, r, parts)
	}
	writeConsensusDoc(t, ideaDir, "consensus")
	fc := &fakeConsensus{statusSeq: []string{consensus.TriageBlocked}}
	fr := &fakeRunner{}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, fr)

	action, _, err := d.Advance(context.Background())
	if err == nil {
		t.Fatal("expected MaxRounds escalation")
	}
	if action != ActionEscalated {
		t.Fatalf("action=%s, want escalated", action)
	}
	if len(fr.calls) != 0 {
		t.Fatal("RunRound must not run past MaxRounds")
	}
}

func TestConsensusMalformedEscalates(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	fc := &fakeConsensus{statusSeq: []string{consensus.TriageMalformed}}
	d := newConsensusDriver(ideaDir, runDir, parts, fc, &fakeRunner{})

	action, _, err := d.Advance(context.Background())
	if err == nil {
		t.Fatal("expected escalation for malformed consensus")
	}
	if action != ActionEscalated {
		t.Fatalf("action=%s, want escalated", action)
	}
}

func TestConsensusGateUnwiredStopsAtBoundary(t *testing.T) {
	// Without a ConsensusOps, the round budget being spent stops at the consensus
	// boundary (slice-1 behavior), never auto-driving consensus.
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	writeAll(t, ideaDir, 2, parts)
	appendEvent(t, runDir, "round.completed", "round-02")
	d := newTestDriver(ideaDir, runDir, parts, 1, true, &fakeRunner{}) // no Consensus

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionConsensus {
		t.Fatalf("action=%s, want consensus-ready (gate unwired)", action)
	}
}

func TestFinalScaffoldReason(t *testing.T) {
	// FINAL.md lives at ideas/<slug>/FINAL.md in production, and the gate now checks that the
	// frontmatter slug matches the directory it closes, so the fixture directory must be the slug.
	dir := filepath.Join(t.TempDir(), "demo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(body string) string {
		p := filepath.Join(dir, "FINAL.md")
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	if r := finalScaffoldReason(write(validFinal)); r != "" {
		t.Fatalf("valid FINAL.md rejected: %s", r)
	}
	if finalScaffoldReason(write("too short")) == "" {
		t.Fatal("short FINAL.md should be rejected")
	}
	if finalScaffoldReason(write("---\nidea: demo\nstatus: draft\n---\n\n## Final plan / specification\nplenty of words here to exceed the length threshold for the scaffold check, padding padding padding padding padding padding padding.\nline two\nline three\n\n\n## Purpose / user-visible outcome\nN/A\n\n## Context & orientation\nN/A\n\n## Observable acceptance criteria\nN/A\n\n## Idempotence & recovery\nN/A\n\n## Known risks / de-risking\nN/A\n\n## References\nN/A\n")) == "" {
		t.Fatal("non-final status should be rejected")
	}
	if finalScaffoldReason(write("---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\n<...> padding padding padding padding padding padding padding padding padding padding padding padding padding padding padding.\nline two\nline three\n\n\n## Purpose / user-visible outcome\nN/A\n\n## Context & orientation\nN/A\n\n## Observable acceptance criteria\nN/A\n\n## Idempotence & recovery\nN/A\n\n## Known risks / de-risking\nN/A\n\n## References\nN/A\n")) == "" {
		t.Fatal("unexpanded <...> placeholder should be rejected")
	}
	// Legitimate angle-bracket content (e.g. a help-text example) must be allowed,
	// not treated as a scaffold placeholder (false positive found in the live run).
	if r := finalScaffoldReason(write("---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\nReword the error to `Unknown option '<option>'` and point users at --help for the path `<path>`.\nKeep the change to wording only and verify the help output renders.\nThis line pads the section beyond the length threshold for the scaffold check comfortably.\n\n\n## Purpose / user-visible outcome\nN/A\n\n## Context & orientation\nN/A\n\n## Observable acceptance criteria\nN/A\n\n## Idempotence & recovery\nN/A\n\n## Known risks / de-risking\nN/A\n\n## References\nN/A\n")); r != "" {
		t.Fatalf("legitimate <option>/<path> content rejected: %s", r)
	}
}

// Review round-01 CRITICAL: the §4.0 cross-review cap was clamped only on the initially
// scheduled budget, so a BLOCKed consensus reopened rounds under MaxRounds alone and
// walked past the printed cap. With HardCrossReviewCap=3, round 5 must escalate and no
// round-05 runner call may happen.
func TestBlockedConsensusRespectsTheHardCrossReviewCap(t *testing.T) {
	dir := t.TempDir()
	ideaDir := filepath.Join(dir, "parley-deck", "ideas", "demo")
	runDir := filepath.Join(dir, "runs", "r1")
	os.MkdirAll(runDir, 0o755)
	for r := 1; r <= 4; r++ {
		os.MkdirAll(filepath.Join(ideaDir, roundLabel(r)), 0o755)
	}
	os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte("---\nidea: demo\ntrack: deliberation\n---\n"), 0o644)
	os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte("x"), 0o644)

	fr := &fakeRunner{}
	d := New(Config{
		IdeaDir: ideaDir, IdeaSlug: "demo", Participants: []string{"a", "b"}, RunDir: runDir,
		Root: dir, Events: store.New(runDir), Auto: true,
		MaxRounds: 4, HardCrossReviewCap: 3,
		Consensus: &fakeConsensus{statusSeq: []string{consensus.TriageBlocked}},
	}, fr)

	action, _, err := d.advanceConsensus(context.Background(), Cursor{Phase: PhaseConsensus})
	if action != ActionEscalated || err == nil {
		t.Fatalf("action=%s err=%v — want escalated at the §4.0 cap", action, err)
	}
	if !strings.Contains(err.Error(), "§4.0 cap of 3") || !strings.Contains(err.Error(), "after 3 cross-review round(s)") {
		t.Errorf("escalation must name the cap it enforced, got: %v", err)
	}
	if len(fr.calls) != 0 {
		t.Fatalf("a 4th cross-review round was opened past the cap: rounds run = %v", fr.calls)
	}
}

// Round-06: the guard above was only ever tested with HardCrossReviewCap injected by
// hand, so deleting the one line in New that wires Policy.CapCrossReviewRounds into the
// config left every test green while the BLOCK back-edge walked past the printed cap.
// This drives the whole seam: a real 00-prompt track → New → a blocked consensus.
func TestTrackWiresTheHardCrossReviewCapThroughNew(t *testing.T) {
	for _, tc := range []struct {
		track   string
		cap     int
		atRound int // rounds already on disk; opening atRound+1 must be refused
	}{
		{"deliberation", 3, 4},
		{"standard", 2, 3},
	} {
		t.Run(tc.track, func(t *testing.T) {
			dir := t.TempDir()
			ideaDir := filepath.Join(dir, "parley-deck", "ideas", "demo")
			runDir := filepath.Join(dir, "runs", "r1")
			os.MkdirAll(runDir, 0o755)
			for r := 1; r <= tc.atRound; r++ {
				os.MkdirAll(filepath.Join(ideaDir, roundLabel(r)), 0o755)
			}
			os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"),
				[]byte("---\nidea: demo\ntrack: "+tc.track+"\n---\n"), 0o644)
			os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte("x"), 0o644)

			fr := &fakeRunner{}
			d := New(Config{
				IdeaDir: ideaDir, IdeaSlug: "demo", Participants: []string{"a", "b", "c"},
				RunDir: runDir, Root: dir, Events: store.New(runDir), Auto: true,
				MaxRounds: 9, // deliberately generous: only the §4.0 cap may stop this
				Consensus: &fakeConsensus{statusSeq: []string{consensus.TriageBlocked}},
			}, fr)

			if d.cfg.HardCrossReviewCap != tc.cap {
				t.Fatalf("track %s: HardCrossReviewCap=%d, want %d — the policy is not wired into the driver",
					tc.track, d.cfg.HardCrossReviewCap, tc.cap)
			}
			action, _, err := d.advanceConsensus(context.Background(), Cursor{Phase: PhaseConsensus})
			if action != ActionEscalated || err == nil {
				t.Fatalf("track %s: action=%s err=%v — a BLOCK past the §4.0 cap must escalate", tc.track, action, err)
			}
			if len(fr.calls) != 0 {
				t.Fatalf("track %s: opened another cross-review round past the cap: %v", tc.track, fr.calls)
			}
		})
	}
}
