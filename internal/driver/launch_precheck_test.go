package driver

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/budget"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

// precheckRunner is a fake RoundRunner that also implements the optional
// round precheck seam, so the wrapper's consult-before-charge order is
// observable without live agents.
type precheckRunner struct {
	*fakeRunner
	failRound bool
	prechecks []int
}

func (p *precheckRunner) PrecheckRound(ctx context.Context, round int) error {
	p.prechecks = append(p.prechecks, round)
	if p.failRound {
		return errors.New("precheck: round refused")
	}
	return nil
}

// precheckConsensus is a fake ConsensusOps implementing the consensus seam.
type precheckConsensus struct {
	calls     []string
	failDraft bool
}

func (c *precheckConsensus) Status() (consensus.Summary, error) { return consensus.Summary{}, nil }
func (c *precheckConsensus) Draft(ctx context.Context) error {
	c.calls = append(c.calls, "draft")
	return nil
}
func (c *precheckConsensus) RequestSignoffs(ctx context.Context, missing []string) error {
	c.calls = append(c.calls, "signoffs")
	return nil
}
func (c *precheckConsensus) DraftFinal(ctx context.Context) error {
	c.calls = append(c.calls, "final")
	return nil
}
func (c *precheckConsensus) Reopen(ctx context.Context, reason string) error {
	c.calls = append(c.calls, "reopen")
	return nil
}
func (c *precheckConsensus) PrecheckConsensusDraft(ctx context.Context) error {
	c.calls = append(c.calls, "precheck-draft")
	if c.failDraft {
		return errors.New("precheck: consensus draft refused")
	}
	return nil
}
func (c *precheckConsensus) PrecheckConsensusFinal(ctx context.Context) error { return nil }
func (c *precheckConsensus) PrecheckConsensusSignoffs(ctx context.Context, missing []string) error {
	return nil
}

// precheckImpl is a fake ImplOps implementing the impl seam.
type precheckImpl struct {
	*fakeImpl
	fail      map[string]bool
	consulted []string
}

func (p *precheckImpl) PrecheckImplementation(ctx context.Context) error {
	p.consulted = append(p.consulted, "implement")
	if p.fail["implement"] {
		return errors.New("precheck: implement refused")
	}
	return nil
}
func (p *precheckImpl) PrecheckReviewLaunch(ctx context.Context, round int) error {
	p.consulted = append(p.consulted, "review")
	return nil
}
func (p *precheckImpl) PrecheckReviewConsensus(ctx context.Context, round int) error {
	p.consulted = append(p.consulted, "review-consensus")
	return nil
}
func (p *precheckImpl) PrecheckReviewSignoffs(ctx context.Context, missing []string) error {
	p.consulted = append(p.consulted, "review-signoffs")
	return nil
}
func (p *precheckImpl) PrecheckGoalCheck(ctx context.Context) error {
	p.consulted = append(p.consulted, "goal-check")
	return nil
}

func assertNoStepCharge(t *testing.T, root string) {
	t.Helper()
	b, err := budget.LoadStepBinding(context.Background(), root, "demo")
	if err != nil {
		t.Fatalf("load step binding: %v", err)
	}
	if b == nil {
		return
	}
	s, err := b.Store.Inspect(context.Background())
	if err != nil || len(s.Entries) != 0 {
		t.Fatalf("precheck refusal spent a step reservation: %+v %v", s, err)
	}
}

// The wrappers must consult the optional precheck BEFORE the step charge, and
// a refusing precheck must stop the real operation with zero reservations.
func TestStepWrappersConsultPrecheckBeforeCharging(t *testing.T) {
	t.Run("round", func(t *testing.T) {
		parts := []string{"codex", "claude"}
		ideaDir, runDir := setupIdea(t, parts, "")
		writeAll(t, ideaDir, 1, parts)
		appendEvent(t, runDir, "round.completed", "round-01")
		fr := &precheckRunner{fakeRunner: &fakeRunner{}, failRound: true}
		d := newTestDriver(ideaDir, runDir, parts, 1, true, fr)
		if _, _, err := d.Advance(context.Background()); err == nil || !strings.Contains(err.Error(), "precheck: round refused") {
			t.Fatalf("refusing precheck ignored: %v", err)
		}
		if len(fr.calls) != 0 || len(fr.prechecks) != 1 || fr.prechecks[0] != 2 {
			t.Fatalf("round ran or precheck mis-scoped: calls=%v prechecks=%v", fr.calls, fr.prechecks)
		}
		assertNoStepCharge(t, d.cfg.Root)
	})
	t.Run("consensus-draft", func(t *testing.T) {
		parts := []string{"codex", "claude"}
		ideaDir, runDir := setupIdea(t, parts, "")
		writeAll(t, ideaDir, 1, parts)
		appendEvent(t, runDir, "round.completed", "round-01")
		fc := &precheckConsensus{failDraft: true}
		d := newTestDriver(ideaDir, runDir, parts, 0, true, &fakeRunner{})
		d.cfg.Consensus = fc
		if _, _, err := d.Advance(context.Background()); err == nil || !strings.Contains(err.Error(), "precheck: consensus draft refused") {
			t.Fatalf("refusing precheck ignored: %v", err)
		}
		if len(fc.calls) != 1 || fc.calls[0] != "precheck-draft" {
			t.Fatalf("draft ran or precheck not consulted first: %v", fc.calls)
		}
		assertNoStepCharge(t, d.cfg.Root)
	})
	t.Run("implement", func(t *testing.T) {
		parts := []string{"codex", "agy"}
		ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
		writeFinalValid(t, ideaDir)
		pi := &precheckImpl{fakeImpl: &fakeImpl{}, fail: map[string]bool{"implement": true}}
		d := newImplDriver(ideaDir, runDir, parts, true, pi)
		action, _, err := d.Advance(context.Background())
		if err == nil || action != ActionEscalated || !strings.Contains(err.Error(), "precheck: implement refused") {
			t.Fatalf("refusing precheck ignored: %s %v", action, err)
		}
		if len(pi.calls) != 0 || len(pi.consulted) != 1 || pi.consulted[0] != "implement" {
			t.Fatalf("implement ran or precheck not consulted first: calls=%v consulted=%v", pi.calls, pi.consulted)
		}
		assertNoStepCharge(t, d.cfg.Root)
	})
}

// N1 with the PRODUCTION adapter and a real (fake CLI) child: a known
// invalid-protocol refusal must spend no step and no cross-review cycle and
// start no process; restoring the exact source permits one grouped attempt
// with the original count semantics, which then stays spent.
func TestN1RoundPrecheckRefusesBeforeStepAndCycleReservations(t *testing.T) {
	parts := []string{"builder", "reviewer"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "parley-deck", "meta"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "version.json"), []byte(`{"protocolRole":"source"}`), 0600); err != nil {
		t.Fatal(err)
	}
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", HeadlessArgs: []string{"-test.run=TestN1RoundChildHelper", "--", "parley-n1-round-child"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}
	opts := runner.Options{Root: root, RunID: filepath.Base(runDir), Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: []string{"builder"}}, Agents: []agents.Discovery{agent}, Timeout: time.Second, Store: store.New(runDir)}
	// This fixture exercises an already admitted scope. A retained pre-start
	// request correctly requires explicit migration before FIRST binding;
	// restoration alone is not permission to forget that history.
	if _, err := budget.EnsureCycleBinding(context.Background(), root, "demo", budget.CrossReview, 4, 0, runDir, ideaDir); err != nil {
		t.Fatal(err)
	}
	warm := opts
	warm.Round = 2
	if err := runner.PrecheckRound(context.Background(), warm); err != nil {
		t.Fatal(err)
	}
	restore := refuseDriverProtocolCache(t, root)
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), CrossReviewRounds: 1, Auto: true, MaxDriverSteps: 5}, NewRunnerAdapter(opts))
	if _, _, err := d.Advance(context.Background()); err == nil || !strings.Contains(err.Error(), "protocol context refused") {
		t.Fatalf("known refusal charged through: %v", err)
	}
	assertNoStepCharge(t, root)
	if b, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.CrossReview); err != nil {
		t.Fatal(err)
	} else if b != nil {
		if state, err := b.Store.Inspect(context.Background()); err != nil || b.Count(state) != 0 {
			t.Fatalf("known refusal spent a cross-review cycle: %+v %v", state, err)
		}
	}
	started, refused := scanDriverInvocationTerminals(t, root)
	if started != 0 || refused != 1 {
		t.Fatalf("known refusal lifecycle: started=%d refused=%d", started, refused)
	}
	restore()
	if _, _, err := d.Advance(context.Background()); err == nil {
		t.Fatal("child fixture exits 7; the real attempt must fail")
	}
	step, err := budget.LoadStepBinding(context.Background(), root, "demo")
	if err != nil || step == nil {
		t.Fatalf("expected configured step binding: %v, %v", step, err)
	}
	stepState, err := step.Store.Inspect(context.Background())
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("grouped attempt did not spend exactly one step: %+v %v", stepState, err)
	}
	cyc, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.CrossReview)
	if err != nil || cyc == nil {
		t.Fatalf("expected retained cycle binding: %v, %v", cyc, err)
	}
	cycState, err := cyc.Store.Inspect(context.Background())
	if err != nil || cyc.Count(cycState) != 1 {
		t.Fatalf("grouped attempt did not spend exactly one cycle: %+v %v", cycState, err)
	}
	started, refused = scanDriverInvocationTerminals(t, root)
	if started != 1 || refused != 1 {
		t.Fatalf("grouped attempt lifecycle: started=%d refused=%d", started, refused)
	}
	// The incomplete round awaits without a re-render or a new charge; the
	// failed actual attempt stays spent.
	if a, _, err := d.Advance(context.Background()); err != nil || a != ActionAwait {
		t.Fatalf("incomplete round re-dispatched: %s %v", a, err)
	}
	stepState, err = step.Store.Inspect(context.Background())
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("await re-spent the step: %+v %v", stepState, err)
	}
	cycState, err = cyc.Store.Inspect(context.Background())
	if err != nil || cyc.Count(cycState) != 1 {
		t.Fatalf("await re-spent the cycle: %+v %v", cycState, err)
	}
}

func TestN1RoundChildHelper(t *testing.T) {
	for _, arg := range os.Args {
		if arg == "parley-n1-round-child" {
			os.Exit(7)
		}
	}
}

// refuseDriverProtocolCache corrupts the cached protocol publication after a
// warm render, so the next shared-renderer launch refuses with a hash
// mismatch while the material source is unchanged.
func refuseDriverProtocolCache(t *testing.T, root string) func() {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(protocolpacket.RuntimeDir(root), "*.md"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("expected one prepared protocol context: %v %v", paths, err)
	}
	original, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(paths[0], []byte("tampered cached context"), 0600); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.WriteFile(paths[0], original, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func scanDriverInvocationTerminals(t *testing.T, root string) (started, refused int) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var r telemetry.Record
		if err := json.Unmarshal(data, &r); err != nil {
			t.Fatal(err)
		}
		if r.StartedAt != nil {
			started++
		}
		if r.Outcome.FailureClass != nil && *r.Outcome.FailureClass == "protocol_context_refused" {
			refused++
		}
	}
	return started, refused
}

// N1 MAJOR-1 fresh-scope production order (round path). NO manual prebinding —
// the prebound TestN1RoundPrecheckRefusesBeforeStepAndCycleReservations above
// remains the control for an already admitted scope. Here the fresh,
// never-bound cross-review scope must be admitted by the driver itself,
// charge-free, BEFORE the protocol precheck: a tampered cached protocol
// refuses with no charge and no start, and restoring the exact source permits
// one grouped fake launch and one charge with no artificial migration.
func TestN1FreshScopeAdmittedBeforeProtocolPrecheck(t *testing.T) {
	parts := []string{"builder", "reviewer"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "parley-deck", "meta"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "version.json"), []byte(`{"protocolRole":"source"}`), 0600); err != nil {
		t.Fatal(err)
	}
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", HeadlessArgs: []string{"-test.run=TestN1RoundChildHelper", "--", "parley-n1-round-child"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}
	opts := runner.Options{Root: root, RunID: filepath.Base(runDir), Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: []string{"builder"}}, Agents: []agents.Discovery{agent}, Timeout: time.Second, Store: store.New(runDir)}
	// Fresh scope: deliberately NO budget.EnsureCycleBinding prebind here —
	// admission must come from the production order, not from the fixture.
	warm := opts
	warm.Round = 2
	if err := runner.PrecheckRound(context.Background(), warm); err != nil {
		t.Fatal(err)
	}
	restore := refuseDriverProtocolCache(t, root)
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), CrossReviewRounds: 1, Auto: true, MaxDriverSteps: 5}, NewRunnerAdapter(opts))
	if _, _, err := d.Advance(context.Background()); err == nil || !strings.Contains(err.Error(), "protocol context refused") {
		t.Fatalf("known refusal charged through: %v", err)
	}
	assertNoStepCharge(t, root)
	// Charge-free admission: the binding EXISTS (admitted before the precheck)
	// with zero cycles spent — the refusal cost nothing.
	b, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.CrossReview)
	if err != nil || b == nil {
		t.Fatalf("fresh scope was not admitted before the precheck: %v %v", b, err)
	}
	if state, err := b.Store.Inspect(context.Background()); err != nil || b.Count(state) != 0 {
		t.Fatalf("known refusal spent a cross-review cycle: %+v %v", state, err)
	}
	started, refused := scanDriverInvocationTerminals(t, root)
	if started != 0 || refused != 1 {
		t.Fatalf("known refusal lifecycle: started=%d refused=%d", started, refused)
	}
	restore()
	// Restoration alone permits the grouped attempt: the already admitted fresh
	// scope needs no artificial migration.
	if _, _, err := d.Advance(context.Background()); err == nil {
		t.Fatal("child fixture exits 7; the real attempt must fail")
	}
	step, err := budget.LoadStepBinding(context.Background(), root, "demo")
	if err != nil || step == nil {
		t.Fatalf("expected configured step binding: %v, %v", step, err)
	}
	stepState, err := step.Store.Inspect(context.Background())
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("grouped attempt did not spend exactly one step: %+v %v", stepState, err)
	}
	b, err = budget.LoadCycleBinding(context.Background(), root, "demo", budget.CrossReview)
	if err != nil || b == nil {
		t.Fatalf("expected admitted cycle binding: %v %v", b, err)
	}
	cycState, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(cycState) != 1 {
		t.Fatalf("grouped attempt did not spend exactly one cycle: %+v %v", cycState, err)
	}
	started, refused = scanDriverInvocationTerminals(t, root)
	if started != 1 || refused != 1 {
		t.Fatalf("grouped attempt lifecycle: started=%d refused=%d", started, refused)
	}
	// The incomplete round awaits without a re-render or a new charge; the
	// failed actual attempt stays spent.
	if a, _, err := d.Advance(context.Background()); err != nil || a != ActionAwait {
		t.Fatalf("incomplete round re-dispatched: %s %v", a, err)
	}
	stepState, err = step.Store.Inspect(context.Background())
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("await re-spent the step: %+v %v", stepState, err)
	}
	cycState, err = b.Store.Inspect(context.Background())
	if err != nil || b.Count(cycState) != 1 {
		t.Fatalf("await re-spent the cycle: %+v %v", cycState, err)
	}
}

// N1 MAJOR-1 fresh-scope production order (fix-up path). Same contract as the
// round path on the fix-up scope, with NO manual prebinding: tampered cached
// protocol => no charge and no start; restored => one grouped fake launch and
// one charge, no artificial migration. The admission lives in impl.go before
// reserveFixupCycle; the prebound fixup cycle tests in cycle_budget_test.go
// remain controls.
func TestN1FreshFixupScopeAdmittedBeforeProtocolPrecheck(t *testing.T) {
	parts := []string{"builder", "reviewer"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	writeFinalValid(t, ideaDir)
	writeImplWithCycles(t, ideaDir, "implemented", 0)
	if err := os.MkdirAll(filepath.Join(ideaDir, "review", "round-01"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ideaDir, "review", "consensus.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	if err := protocol.InitWorkspace(root); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "parley-deck", "meta"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "parley-deck", "meta", "version.json"), []byte(`{"protocolRole":"source"}`), 0600); err != nil {
		t.Fatal(err)
	}
	agent := agents.Discovery{Spec: agents.Spec{ID: "builder", HeadlessArgs: []string{"-test.run=TestCycleDriverChildHelper", "--", "parley-cycle-child"}, PromptMode: agents.PromptStdin}, Path: os.Args[0], Found: true}
	opts := runner.Options{Root: root, RunID: filepath.Base(runDir), Idea: protocol.IdeaStatus{Slug: "demo", Path: ideaDir, Participants: []string{"builder"}}, Agents: []agents.Discovery{agent}, Timeout: time.Second, Store: store.New(runDir)}
	// Fresh fix-up scope: deliberately NO budget.EnsureCycleBinding prebind —
	// admission must come from the production order, not from the fixture.
	if err := runner.PrecheckFixup(context.Background(), opts); err != nil {
		t.Fatal(err)
	}
	restore := refuseDriverProtocolCache(t, root)
	fi := &cycleProcessImpl{fakeImpl: &fakeImpl{roundComplete: true, checksOK: true, review: ReviewStatus{Summary: consensus.Summary{Triage: consensus.TriageReady}, OutstandingAgreedFixes: 1}}, opts: opts}
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 1, MaxDriverSteps: 5, Impl: fi}, &fakeRunner{})
	if _, _, err := d.Advance(context.Background()); err == nil || !strings.Contains(err.Error(), "protocol context refused") {
		t.Fatalf("known refusal charged through: %v", err)
	}
	if len(fi.calls) != 0 {
		t.Fatalf("refusing precheck still launched the fix-up: %v", fi.calls)
	}
	assertNoStepCharge(t, root)
	// Charge-free admission: the fix-up binding EXISTS with zero cycles spent.
	b, err := budget.LoadCycleBinding(context.Background(), root, "demo", budget.Fixup)
	if err != nil || b == nil {
		t.Fatalf("fresh fix-up scope was not admitted before the precheck: %v %v", b, err)
	}
	if state, err := b.Store.Inspect(context.Background()); err != nil || b.Count(state) != 0 {
		t.Fatalf("known refusal spent a fix-up cycle: %+v %v", state, err)
	}
	started, refused := scanDriverInvocationTerminals(t, root)
	if started != 0 || refused != 1 {
		t.Fatalf("known refusal lifecycle: started=%d refused=%d", started, refused)
	}
	restore()
	// Restoration alone permits the grouped attempt: no artificial migration.
	if _, _, err := d.Advance(context.Background()); err == nil {
		t.Fatal("child fixture exits 7; the real attempt must fail")
	}
	if !contains(fi.calls, "fixup") {
		t.Fatalf("grouped attempt did not reach the fix-up launch: %v", fi.calls)
	}
	step, err := budget.LoadStepBinding(context.Background(), root, "demo")
	if err != nil || step == nil {
		t.Fatalf("expected configured step binding: %v, %v", step, err)
	}
	stepState, err := step.Store.Inspect(context.Background())
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("grouped attempt did not spend exactly one step: %+v %v", stepState, err)
	}
	b, err = budget.LoadCycleBinding(context.Background(), root, "demo", budget.Fixup)
	if err != nil || b == nil {
		t.Fatalf("expected admitted fix-up binding: %v %v", b, err)
	}
	cycState, err := b.Store.Inspect(context.Background())
	if err != nil || b.Count(cycState) != 1 {
		t.Fatalf("grouped attempt did not spend exactly one fix-up cycle: %+v %v", cycState, err)
	}
	started, refused = scanDriverInvocationTerminals(t, root)
	if started != 1 || refused != 1 {
		t.Fatalf("grouped attempt lifecycle: started=%d refused=%d", started, refused)
	}
	// The failed actual attempt stays spent: re-entry at the inclusive cap
	// escalates before any re-render, re-launch or new charge.
	fi.calls = nil
	if _, _, err := d.Advance(context.Background()); err == nil {
		t.Fatal("charged attempt at the cap was refunded")
	}
	if len(fi.calls) != 0 {
		t.Fatalf("cap re-entry re-launched the fix-up: %v", fi.calls)
	}
	stepState, err = step.Store.Inspect(context.Background())
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("cap re-entry re-spent the step: %+v %v", stepState, err)
	}
	cycState, err = b.Store.Inspect(context.Background())
	if err != nil || b.Count(cycState) != 1 {
		t.Fatalf("cap re-entry re-spent the cycle: %+v %v", cycState, err)
	}
}

// MINOR-1 (round path): an admitted cross-review scope already at its frozen
// maximum is a KNOWN exhaustion. The driver must refuse BEFORE the step
// charge — no step spent, no dispatch, no extra cycle — while the binding and
// its history stay exactly as admitted.
func TestN1ExhaustedCrossReviewCycleRefusesBeforeStep(t *testing.T) {
	parts := []string{"codex", "claude"}
	ideaDir, runDir := setupIdea(t, parts, "")
	writeAll(t, ideaDir, 1, parts)
	appendEvent(t, runDir, "round.completed", "round-01")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	ctx := context.Background()
	b, err := budget.EnsureCycleBinding(ctx, root, "demo", budget.CrossReview, 1, 0, runDir, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	cctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := budget.ChargeCycle(cctx, budget.CrossReview); err != nil {
		t.Fatal(err)
	}
	finish()
	fr := &fakeRunner{}
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), CrossReviewRounds: 1, MaxRounds: 1, Auto: true, MaxDriverSteps: 5}, fr)
	if _, _, err := d.Advance(ctx); err == nil || !strings.Contains(err.Error(), "budget exhausted") {
		t.Fatalf("known exhausted cycle did not refuse before the step: %v", err)
	}
	if len(fr.calls) != 0 {
		t.Fatalf("known exhausted cycle still dispatched the round: %v", fr.calls)
	}
	assertNoStepCharge(t, root)
	state, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(state) != 1 {
		t.Fatalf("known refusal re-spent the cycle: %+v %v", state, err)
	}
}

// MINOR-1 (fix-up path): reserveFixupCycle must check the KNOWN exhausted
// fix-up cycle BEFORE ChargeStep — the doomed transition spends no step.
func TestN1ReserveFixupCyclePreflightsExhaustionBeforeStep(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	ctx := context.Background()
	b, err := budget.EnsureCycleBinding(ctx, root, "demo", budget.Fixup, 1, 0, runDir, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	cctx, finish, err := budget.OpenCycleSession(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := budget.ChargeCycle(cctx, budget.Fixup); err != nil {
		t.Fatal(err)
	}
	finish()
	steps, err := budget.EnsureStepBinding(ctx, root, "demo", 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	sctx, sfinish, err := budget.OpenStepSession(ctx, steps)
	if err != nil {
		t.Fatal(err)
	}
	defer sfinish()
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 1}, &fakeRunner{})
	if _, _, _, err := d.reserveFixupCycle(sctx, 0); err == nil || !strings.Contains(err.Error(), "budget exhausted") {
		t.Fatalf("known exhausted fix-up cycle did not refuse before the step: %v", err)
	}
	state, err := steps.Store.Inspect(ctx)
	if err != nil || len(state.Entries) != 0 {
		t.Fatalf("known exhausted fix-up cycle spent a step: %+v %v", state, err)
	}
}

// MINOR-1 reuse invariant (fix-up path): a grouped session ALREADY charged
// for this logical action stays reusable at the cap — the preflight must not
// refuse it (a naive current-count-only check would), and re-entry through
// the full wiring spends nothing new on either budget.
func TestN1ReserveFixupCycleChargedSessionReusableAtCap(t *testing.T) {
	parts := []string{"codex", "agy"}
	ideaDir, runDir := setupIdea(t, parts, "auto_implement: true\n")
	root := filepath.Dir(filepath.Dir(filepath.Dir(ideaDir)))
	ctx := context.Background()
	b, err := budget.EnsureCycleBinding(ctx, root, "demo", budget.Fixup, 1, 0, runDir, ideaDir)
	if err != nil {
		t.Fatal(err)
	}
	steps, err := budget.EnsureStepBinding(ctx, root, "demo", 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	sctx, sfinish, err := budget.OpenStepSession(ctx, steps)
	if err != nil {
		t.Fatal(err)
	}
	defer sfinish()
	d := New(Config{Root: root, IdeaDir: ideaDir, IdeaSlug: "demo", RunDir: runDir, Participants: parts, Events: store.New(runDir), Auto: true, AutoImplement: true, MaxFixupCycles: 1}, &fakeRunner{})
	rctx, n, finish, err := d.reserveFixupCycle(sctx, 0)
	if err != nil || n != 1 {
		t.Fatalf("first grouped charge: ordinal=%d err=%v", n, err)
	}
	defer finish()
	// Re-enter on the same live grouped context at the cap: reuse, no new spend.
	if _, n2, finish2, err := d.reserveFixupCycle(rctx, 1); err != nil || n2 != 1 {
		t.Fatalf("charged grouped session not reusable at the cap: ordinal=%d err=%v", n2, err)
	} else {
		defer finish2()
	}
	stepState, err := steps.Store.Inspect(ctx)
	if err != nil || len(stepState.Entries) != 1 {
		t.Fatalf("grouped reuse duplicated the step charge: %+v %v", stepState, err)
	}
	cycState, err := b.Store.Inspect(ctx)
	if err != nil || b.Count(cycState) != 1 {
		t.Fatalf("grouped reuse duplicated the cycle charge: %+v %v", cycState, err)
	}
}
