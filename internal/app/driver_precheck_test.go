package app

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/protocolcore"
	"parley-deck-cli/internal/protocolpacket"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/telemetry"
)

// driverPrecheckRefusals returns the retained invocation terminals under root
// that match keep. A precheck refusal is retained WITHOUT a start; a
// successful precheck must retain nothing at all.
func driverPrecheckRefusals(t *testing.T, root string, keep func(telemetry.Record) bool) []telemetry.Record {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(root, ".parley-runtime", "invocations", "*", "terminal.json"))
	if err != nil {
		t.Fatal(err)
	}
	var out []telemetry.Record
	for _, p := range paths {
		var r telemetry.Record
		if _, err := readVerificationJSON(p, &r); err != nil {
			t.Fatal(err)
		}
		if keep(r) {
			out = append(out, r)
		}
	}
	return out
}

// tamperProtocolCacheRequest poisons the exact cached context body for an
// explicit render request and returns a restore func. It is the request-exact
// sibling of refuseAppProtocolCache for launches whose render scope is not a
// simple phase number (the consensus drafter renders with empty idea/phase).
func tamperProtocolCacheRequest(t *testing.T, root string, req protocolpacket.Request) func() {
	t.Helper()
	prepared, err := protocolpacket.Render(root, protocolcore.StoreAt(config.CentralHome()), req)
	if err != nil || prepared.BodyPath == "" {
		t.Fatal("cannot prepare exact cached context", req, err)
	}
	raw, err := os.ReadFile(prepared.BodyPath)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(prepared.BodyPath, []byte("tampered cached context"), 0600); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.WriteFile(prepared.BodyPath, raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func precheckImplOpsFixture(t *testing.T, runID string) (driverImplOps, string, agents.Discovery) {
	t.Helper()
	root, _, _, reviewer := trajectoryHelperFixture(t, "actual-helper")
	builder := reviewer
	builder.ID = "builder"
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea-x")
	o := driverImplOps{root: root, ideaSlug: "idea-x", ideaDir: ideaDir, implementer: "builder", drafter: "reviewer", out: io.Discard,
		base: runner.Options{Root: root, RunID: runID,
			Idea:   protocol.IdeaStatus{Slug: "idea-x", Path: ideaDir, Participants: []string{"builder", "reviewer"}},
			Agents: []agents.Discovery{reviewer, builder},
			Store:  store.New(filepath.Join(root, "parley-deck", "runs", runID))}}
	return o, root, reviewer
}

func assertSingleUnstartedRefusal(t *testing.T, records []telemetry.Record, agent, phase string) {
	t.Helper()
	if len(records) != 1 {
		t.Fatal("precheck did not retain exactly one unstarted refusal", len(records))
	}
	r := records[0]
	if r.Metadata.Agent != agent || r.Metadata.Idea != "idea-x" || r.Metadata.Phase != phase ||
		r.StartedAt != nil || r.Outcome == nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatalf("precheck refusal misattributed or started: %+v", r)
	}
}

// The production implementation precheck must refuse a known-tampered protocol
// before the implementer launch, attributed to the SELECTED implementer (not
// the base participant list), and a restored protocol must pass with no
// precheck launch or charge retained. The launch must be genuinely PENDING:
// with an existing IMPLEMENTATION.md the real RunImplementation is a skip-no-op
// (runAgent skips an existing artifact when Overwrite is false), so a refusal
// there would be a false precheck refusal.
func TestDriverPrecheckImplementationRefusalAttributionAndRestore(t *testing.T) {
	o, root, _ := precheckImplOpsFixture(t, "kimi-precheck-impl")
	// The trajectory scratch fixture (gateScratchRepo) ships an implemented
	// IMPLEMENTATION.md; remove it so the implementer launch is pending — the
	// only state in which the production Implement would actually launch and
	// render. The behavioral assertions below are unchanged; only the
	// precondition is established honestly.
	if err := os.Remove(filepath.Join(o.ideaDir, "IMPLEMENTATION.md")); err != nil {
		t.Fatal("establish pending-implementation precondition", err)
	}
	mine := func(r telemetry.Record) bool { return r.Metadata.RunID == "kimi-precheck-impl" }
	restore := refuseAppProtocolCache(t, root, 5) // phase 5 = implementation
	if err := o.PrecheckImplementation(context.Background()); err == nil {
		t.Fatal("tampered cached protocol was not refused before the implementation launch")
	}
	assertSingleUnstartedRefusal(t, driverPrecheckRefusals(t, root, mine), "builder", "implementation")
	restore()
	if err := o.PrecheckImplementation(context.Background()); err != nil {
		t.Fatal("restored protocol still refused the implementation precheck", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 1 {
		t.Fatal("successful precheck launched or charged a reservation", n)
	}
}

// The mirror image of the refusal case: with the fixture's IMPLEMENTATION.md
// still present, the real RunImplementation launches nothing (runAgent skips
// an existing artifact when Overwrite is false), so the app precheck must
// mirror that skip — nil error and zero retained records — even with a
// tampered cached protocol. A refusal here would refuse a launch production
// never performs.
func TestDriverPrecheckImplementationSkipsExistingArtifact(t *testing.T) {
	o, root, _ := precheckImplOpsFixture(t, "kimi-precheck-impl-skip")
	if _, err := os.Stat(filepath.Join(o.ideaDir, "IMPLEMENTATION.md")); err != nil {
		t.Fatal("fixture precondition drifted: IMPLEMENTATION.md must pre-exist", err)
	}
	mine := func(r telemetry.Record) bool { return r.Metadata.RunID == "kimi-precheck-impl-skip" }
	restore := refuseAppProtocolCache(t, root, 5) // phase 5 = implementation
	defer restore()
	if err := o.PrecheckImplementation(context.Background()); err != nil {
		t.Fatal("precheck refused a launch the real runner would skip", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 0 {
		t.Fatal("skipped implementation precheck retained a record", n)
	}
}

// The goal-check precheck mirrors GoalCheck's checker resolution: the
// non-implementer drafter is the checker; the self-checker and
// unknown-checker branches never launch, so they must stay nil (the real call
// fails them) and must retain no refusal. Tamper => attributed refusal;
// restore => clean pass with nothing launched.
func TestDriverPrecheckGoalCheckRefusalAndSkipBranches(t *testing.T) {
	o, root, _ := precheckImplOpsFixture(t, "kimi-precheck-goal")
	mine := func(r telemetry.Record) bool { return r.Metadata.RunID == "kimi-precheck-goal" }
	restore := refuseAppProtocolCache(t, root, -1) // goal-check renders with the unspecified phase
	if err := o.PrecheckGoalCheck(context.Background()); err == nil {
		t.Fatal("tampered cached protocol was not refused before the goal-check launch")
	}
	assertSingleUnstartedRefusal(t, driverPrecheckRefusals(t, root, mine), "reviewer", "goal-check")
	selfChecker := o
	selfChecker.drafter = "builder" // == implementer
	if err := selfChecker.PrecheckGoalCheck(context.Background()); err != nil {
		t.Fatal("self-checker branch must remain a real-call failure, not a precheck refusal", err)
	}
	unknownChecker := o
	unknownChecker.drafter = "ghost"
	if err := unknownChecker.PrecheckGoalCheck(context.Background()); err != nil {
		t.Fatal("unknown-checker branch must remain a real-call failure, not a precheck refusal", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 1 {
		t.Fatal("no-launch checker branches recorded a refusal or launch", n)
	}
	restore()
	if err := o.PrecheckGoalCheck(context.Background()); err != nil {
		t.Fatal("restored protocol still refused the goal-check precheck", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 1 {
		t.Fatal("successful goal-check precheck launched or charged a reservation", n)
	}
}

// The review-signoff precheck resolves the real review/consensus.md, selects
// the missing headless signer through the production request-signoffs
// resolution, and renders the review-consensus scope. Tamper => one unstarted
// refusal attributed to that signer/phase; restore => clean pass.
func TestDriverPrecheckReviewSignoffsRefusalAndRestore(t *testing.T) {
	_, root, _ := precheckImplOpsFixture(t, "kimi-precheck-review-signoff")
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea-x")
	reviewDir := filepath.Join(ideaDir, "review")
	if err := os.MkdirAll(reviewDir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := "---\nidea: idea-x\ndrafted-by: reviewer\ndate: 2026-09-16\n---\n\n## Agreed fixes\n\n## Signoffs\n"
	if err := os.WriteFile(filepath.Join(reviewDir, "consensus.md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	o := driverImplOps{root: root, ideaSlug: "idea-x", ideaDir: ideaDir, out: io.Discard}
	mine := func(r telemetry.Record) bool {
		return r.Metadata.Phase == "review-consensus" && r.Metadata.Agent == "reviewer"
	}
	restore := refuseAppProtocolCache(t, root, 7) // phase 7 = review-consensus
	if err := o.PrecheckReviewSignoffs(context.Background(), []string{"reviewer"}); err == nil {
		t.Fatal("tampered cached protocol was not refused before the review-signoff launch")
	}
	assertSingleUnstartedRefusal(t, driverPrecheckRefusals(t, root, mine), "reviewer", "review-consensus")
	restore()
	if err := o.PrecheckReviewSignoffs(context.Background(), []string{"reviewer"}); err != nil {
		t.Fatal("restored protocol still refused the review-signoff precheck", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 1 {
		t.Fatal("successful review-signoff precheck launched or charged a reservation", n)
	}
}

// The consensus drafter precheck mirrors runDrafter: firstHeadlessAgent over
// the idea participants and the "one-shot" run identity that CommandFor's
// launch boundary assigns, with the same empty phase/idea render scope.
// Tamper that exact scope => one unstarted refusal; restore => clean pass.
func TestDriverPrecheckConsensusDraftRefusalAndRestore(t *testing.T) {
	root, _, _, reviewer := trajectoryHelperFixture(t, "actual-helper")
	builder := reviewer
	builder.ID = "builder"
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea-x")
	o := driverConsensusOps{root: root, ideaSlug: "idea-x", ideaDir: ideaDir,
		participants: []string{"builder", "reviewer"},
		discovered:   []agents.Discovery{reviewer, builder}, out: io.Discard}
	mine := func(r telemetry.Record) bool { return r.Metadata.RunID == "one-shot" }
	restore := tamperProtocolCacheRequest(t, root, protocolpacket.Request{Phase: -1, Track: "unknown", IdeaSlug: ""})
	if err := o.PrecheckConsensusDraft(context.Background()); err == nil {
		t.Fatal("tampered cached protocol was not refused before the consensus drafter launch")
	}
	records := driverPrecheckRefusals(t, root, mine)
	if len(records) != 1 {
		t.Fatal("drafter precheck did not retain exactly one unstarted refusal", len(records))
	}
	r := records[0]
	if r.Metadata.Agent != "builder" || r.Metadata.Phase != "unspecified" ||
		r.StartedAt != nil || r.Outcome == nil || r.Outcome.FailureClass == nil || *r.Outcome.FailureClass != "protocol_context_refused" {
		t.Fatalf("drafter precheck refusal misattributed or started: %+v", r)
	}
	restore()
	if err := o.PrecheckConsensusDraft(context.Background()); err != nil {
		t.Fatal("restored protocol still refused the consensus drafter precheck", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 1 {
		t.Fatal("successful drafter precheck launched or charged a reservation", n)
	}
}

// Regression for the blocked-consensus mirror (MINOR-1): a review/consensus.md
// whose frontmatter reads `blocked: false` but which carries one ❌ BLOCK
// signoff — with the missing reviewer still awaited — is TriageBlocked by
// signoff derivation, and production's requestConsensusSignoffs aborts there
// BEFORE any agent is resolved, rendered or telemetered. Under a known-tampered
// protocol the precheck must surface the same blocked reason and retain NOTHING
// (no protocol_context_refused terminal, no launch) — not the false refusal the
// unmirrored gate produced — and a restored protocol must not reopen it either.
func TestDriverPrecheckReviewSignoffsBlockedConsensusWins(t *testing.T) {
	_, root, _ := precheckImplOpsFixture(t, "kimi-precheck-blocked-signoff")
	ideaDir := filepath.Join(root, "parley-deck", "ideas", "idea-x")
	reviewDir := filepath.Join(ideaDir, "review")
	if err := os.MkdirAll(reviewDir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := "---\nidea: idea-x\ndrafted-by: reviewer\ndate: 2026-09-16\nblocked: false\n---\n\n## Agreed fixes\n\n## Signoffs\n\n### Signoff: builder — 2026-09-16\nStatus: ❌ BLOCK\nNotes: Fix-up regressed the agreed criterion.\nCounter-proposal: Reopen review round 2.\n"
	if err := os.WriteFile(filepath.Join(reviewDir, "consensus.md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	o := driverImplOps{root: root, ideaSlug: "idea-x", ideaDir: ideaDir, out: io.Discard}
	mine := func(r telemetry.Record) bool { return r.Metadata.Phase == "review-consensus" }
	want := "target consensus is blocked; resolve the BLOCK before requesting more signoffs"
	restore := refuseAppProtocolCache(t, root, 7) // phase 7 = review-consensus
	if err := o.PrecheckReviewSignoffs(context.Background(), []string{"reviewer"}); err == nil || err.Error() != want {
		t.Fatal("blocked consensus must abort with the production blocked reason, not a protocol refusal", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 0 {
		t.Fatal("blocked consensus recorded a protocol-refusal terminal or launch", n)
	}
	restore()
	if err := o.PrecheckReviewSignoffs(context.Background(), []string{"reviewer"}); err == nil || err.Error() != want {
		t.Fatal("a restored protocol must not reopen a blocked consensus", err)
	}
	if n := len(driverPrecheckRefusals(t, root, mine)); n != 0 {
		t.Fatal("blocked consensus retained a record after restore", n)
	}
}
