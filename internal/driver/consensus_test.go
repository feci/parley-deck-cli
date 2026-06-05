package driver

import (
	"context"
	"os"
	"path/filepath"
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

const validFinal = "---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\nThe driver advances the deliberation through a disk-derived cursor and an\nordered switch over the round and consensus phases, gated by transport.\nThe consensus gate drafts, collects signoffs, and finalizes autonomously.\nThis paragraph is padded well beyond two hundred and fifty bytes so the\nnon-scaffold length check passes comfortably during the unit test run.\n"

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
}

func TestConsensusReadyButFinalScaffoldEscalates(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeConsensusDoc(t, ideaDir, "consensus")
	scaffold := "---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\n<fill in the agreed plan here>\n"
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
	dir := t.TempDir()
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
	if finalScaffoldReason(write("---\nidea: demo\nstatus: draft\n---\n\n## Final plan / specification\nplenty of words here to exceed the length threshold for the scaffold check, padding padding padding padding padding padding padding.\nline two\nline three\n")) == "" {
		t.Fatal("non-final status should be rejected")
	}
	if finalScaffoldReason(write("---\nidea: demo\nstatus: final\n---\n\n## Final plan / specification\n<placeholder> padding padding padding padding padding padding padding padding padding padding padding padding padding padding padding.\nline two\nline three\n")) == "" {
		t.Fatal("unexpanded <…> placeholder should be rejected")
	}
}
