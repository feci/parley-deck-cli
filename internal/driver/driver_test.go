package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"parley-deck-cli/internal/store"
)

// fakeRunner records RunRound calls and optionally writes round artifacts, so the
// driver core is exercised without launching live agents (consensus D14).
type fakeRunner struct {
	calls      []int
	err        error
	writeOnRun func(round int)
}

func (f *fakeRunner) RunRound(ctx context.Context, round int) error {
	f.calls = append(f.calls, round)
	if f.err != nil {
		return f.err
	}
	if f.writeOnRun != nil {
		f.writeOnRun(round)
	}
	return nil
}

// writeArtifact writes a valid round artifact. For round≥2 it includes a
// `### @<other>` heading for every other participant (the D4 cross-review evidence
// the runner prompt now emits and the driver gate now enforces).
func writeArtifact(t *testing.T, ideaDir, slug string, round int, agent string, participants []string) {
	t.Helper()
	dir := filepath.Join(ideaDir, roundLabel(round))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	var body string
	if round <= 1 {
		body = fmt.Sprintf("---\nagent: %s\nidea: %s\nround: 1\n---\n\n## Summary\nx\n## Proposed approach\nx\n## Concerns / open questions\nx\n## Risks\nx\n", agent, slug)
	} else {
		var headings strings.Builder
		for _, other := range participants {
			if other != agent {
				fmt.Fprintf(&headings, "### @%s\nx\n", other)
			}
		}
		body = fmt.Sprintf("---\nagent: %s\nidea: %s\nround: %d\nresponding-to: [round-01]\n---\n\n## Summary\nx\n## Responses to other participants\n%s## Refined position\nx\n", agent, slug, round, headings.String())
	}
	if err := os.WriteFile(filepath.Join(dir, agent+".md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// setupIdea creates a temp idea dir with a local-dir transport in 00-prompt so the
// per-tick transport gate resolves to local-dir without a COOPERATION.md.
func setupIdea(t *testing.T, participants []string, extraFrontmatter string) (ideaDir, runDir string) {
	t.Helper()
	root := t.TempDir()
	ideaDir = filepath.Join(root, "parley-deck", "ideas", "demo")
	runDir = filepath.Join(root, "parley-deck", "runs", "run1")
	if err := os.MkdirAll(ideaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := fmt.Sprintf("---\nidea: demo\nparticipants: [%s]\ntransport: local-dir\n%sstatus: round-01\n---\n\n## Problem\nx\n", strings.Join(participants, ", "), extraFrontmatter)
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"), []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	return ideaDir, runDir
}

func newTestDriver(ideaDir, runDir string, parts []string, crossReview int, auto bool, r RoundRunner) *Driver {
	return New(Config{
		IdeaDir:           ideaDir,
		IdeaSlug:          "demo",
		Participants:      parts,
		RunDir:            runDir,
		Root:              filepath.Dir(filepath.Dir(filepath.Dir(ideaDir))),
		Events:            store.New(runDir),
		CrossReviewRounds: crossReview,
		Auto:              auto,
	}, r)
}

func writeAll(t *testing.T, ideaDir string, round int, parts []string) {
	for _, p := range parts {
		writeArtifact(t, ideaDir, "demo", round, p, parts)
	}
}

func appendEvent(t *testing.T, runDir, typ, round string) {
	t.Helper()
	if err := store.New(runDir).Append(store.Event{Type: typ, Data: map[string]any{"idea": "demo", "round": round}}); err != nil {
		t.Fatal(err)
	}
}

func TestAdvancePromotesRound01ToRound02(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{writeOnRun: func(round int) { writeAll(t, ideaDir, round, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	action, c, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionPromoted {
		t.Fatalf("action=%s, want promoted", action)
	}
	if len(fr.calls) != 1 || fr.calls[0] != 2 {
		t.Fatalf("RunRound calls=%v, want [2]", fr.calls)
	}
	if c.CurrentRound != 2 {
		t.Fatalf("CurrentRound=%d, want 2", c.CurrentRound)
	}
	if got := readIdeaStatus(ideaDir); got != "round-02" {
		t.Fatalf("idea status=%q, want round-02", got)
	}
	if _, err := os.Stat(filepath.Join(runDir, "driver.json")); err != nil {
		t.Fatalf("cursor not saved: %v", err)
	}
}

func TestAdvanceNoDuplicateAfterRound02Complete(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	writeAll(t, ideaDir, 2, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	appendEvent(t, runDir, "round.completed", "round-02")
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionConsensus {
		t.Fatalf("action=%s, want consensus-ready", action)
	}
	if len(fr.calls) != 0 {
		t.Fatalf("RunRound called %v, want none (no duplicate dispatch)", fr.calls)
	}
}

func TestAdvanceAwaitsIncompleteRound(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeArtifact(t, ideaDir, "demo", 1, "codex", parts) // claude artifact missing
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionAwait {
		t.Fatalf("action=%s, want await", action)
	}
	if len(fr.calls) != 0 {
		t.Fatal("RunRound must not run for an incomplete round")
	}
}

func TestAdvanceRoundIncompleteEventBlocks(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.incomplete", "round-01")
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionAwait {
		t.Fatalf("action=%s, want await (round.incomplete is authoritative)", action)
	}
	if len(fr.calls) != 0 {
		t.Fatal("RunRound must not run when the terminal event is round.incomplete")
	}
}

func TestRound02RequiresCrossReviewHeadings(t *testing.T) {
	// A round-02 artifact with valid frontmatter + responding-to but WITHOUT the
	// per-agent ### @<other> heading must NOT count as complete (D4/AF2).
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	dir := filepath.Join(ideaDir, "round-02")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// codex responds with headings; claude omits the ### @codex heading.
	writeArtifact(t, ideaDir, "demo", 2, "codex", parts)
	weak := "---\nagent: claude\nidea: demo\nround: 2\nresponding-to: [round-01]\n---\n\n## Summary\nx\n## Responses to other participants\ngeneric prose, no per-agent heading\n"
	if err := os.WriteFile(filepath.Join(dir, "claude.md"), []byte(weak), 0o644); err != nil {
		t.Fatal(err)
	}
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	done, err := d.roundComplete(2)
	if err != nil {
		t.Fatal(err)
	}
	if done {
		t.Fatal("round-02 must be incomplete when a participant omits the ### @<other> cross-review heading")
	}
}

func TestAdvanceReconcilesMissingTerminalEvent(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	// No events.jsonl at all → reconciliation must reconstruct round.completed.
	fr := &fakeRunner{writeOnRun: func(round int) { writeAll(t, ideaDir, round, parts) }}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionPromoted {
		t.Fatalf("action=%s, want promoted (reconciled)", action)
	}
	events, err := store.New(runDir).Load()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range events {
		if e.Type != "round.completed" {
			continue
		}
		if r, _ := e.Data["round"].(string); r != "round-01" {
			continue
		}
		if rec, _ := e.Data["reconstructed"].(bool); rec {
			found = true
		}
	}
	if !found {
		t.Fatal("expected a reconstructed round.completed event for round-01")
	}
}

func TestRebuildDerivesPhaseFromDisk(t *testing.T) {
	parts := []string{"codex"}
	ideaDir, _ := setupIdea(t, parts, "")
	writeArtifact(t, ideaDir, "demo", 1, "codex", parts)
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseRound || c.CurrentRound != 1 {
		t.Fatalf("got phase=%s round=%d, want round/1", c.Phase, c.CurrentRound)
	}
	writeArtifact(t, ideaDir, "demo", 2, "codex", parts)
	if c := Rebuild(ideaDir, 4); c.CurrentRound != 2 {
		t.Fatalf("CurrentRound=%d, want 2", c.CurrentRound)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "consensus.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseConsensus {
		t.Fatalf("phase=%s, want consensus", c.Phase)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "FINAL.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if c := Rebuild(ideaDir, 4); c.Phase != PhaseFinal {
		t.Fatalf("phase=%s, want final", c.Phase)
	}
}

func TestCorruptCursorIgnoredRebuildRecovers(t *testing.T) {
	parts := []string{"codex"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeArtifact(t, ideaDir, "demo", 1, "codex", parts)
	writeArtifact(t, ideaDir, "demo", 2, "codex", parts)
	if err := os.WriteFile(filepath.Join(runDir, "driver.json"), []byte("{not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCursor(filepath.Join(runDir, "driver.json")); err == nil {
		t.Fatal("expected corrupt cursor to fail LoadCursor")
	}
	if c := Rebuild(ideaDir, 4); c.CurrentRound != 2 || c.Phase != PhaseRound {
		t.Fatalf("Rebuild got round=%d phase=%s, want 2/round", c.CurrentRound, c.Phase)
	}
}

func TestAdvanceSurfaceOnlyWhenAutoOff(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 1, false, fr) // Auto=false

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionSurfaceOnly {
		t.Fatalf("action=%s, want surface-only", action)
	}
	if len(fr.calls) != 0 {
		t.Fatal("RunRound must not run when --auto is off")
	}
}

func TestAdvanceSurfaceOnlyWhenTransportNotLocalDir(t *testing.T) {
	// Idea-level transport github-pr must disable auto-advance even with --auto.
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	// Override transport to github-pr.
	if err := os.WriteFile(filepath.Join(ideaDir, "00-prompt.md"),
		[]byte("---\nidea: demo\nparticipants: [codex, claude]\ntransport: github-pr\nstatus: round-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, 1, true, fr) // Auto=true but transport=github-pr

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionSurfaceOnly {
		t.Fatalf("action=%s, want surface-only (github-pr must never auto-drive)", action)
	}
	if len(fr.calls) != 0 {
		t.Fatal("RunRound must not run for github-pr transport")
	}
}

func TestCrossReviewRoundsZeroBypassesToConsensus(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "cross_review_rounds: 0\n")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	if got := ReadCrossReviewRounds(ideaDir); got != 0 {
		t.Fatalf("ReadCrossReviewRounds=%d, want 0", got)
	}
	fr := &fakeRunner{}
	d := newTestDriver(ideaDir, runDir, parts, ReadCrossReviewRounds(ideaDir), true, fr)

	action, _, err := d.Advance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if action != ActionConsensus {
		t.Fatalf("action=%s, want consensus-ready (N=0 explicit bypass)", action)
	}
	if len(fr.calls) != 0 {
		t.Fatal("RunRound must not run with cross_review_rounds=0")
	}
}

func TestAcquireLockIsExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "driver.lock")
	const n = 8
	var wins int32
	var mu sync.Mutex
	var releases []func()
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if rel, err := acquireLock(path); err == nil {
				atomic.AddInt32(&wins, 1)
				mu.Lock()
				releases = append(releases, rel)
				mu.Unlock()
			}
		}()
	}
	close(start)
	wg.Wait()
	if wins != 1 {
		t.Fatalf("acquireLock granted %d concurrent holders, want exactly 1", wins)
	}
	for _, rel := range releases {
		rel()
	}
	rel, err := acquireLock(path)
	if err != nil {
		t.Fatalf("re-acquire after release failed: %v", err)
	}
	rel()
}
