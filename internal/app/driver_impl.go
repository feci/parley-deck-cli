package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/config"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
	"parley-deck-cli/internal/evidence"
	"parley-deck-cli/internal/protocol"
	"parley-deck-cli/internal/runner"
	"parley-deck-cli/internal/store"
	"parley-deck-cli/internal/track"
	"parley-deck-cli/internal/trajectory"
)

// driverImplOps is the production driver.ImplOps adapter (driver-impl-phase). It
// reuses the runner Phase 5-8 helpers + review-mode consensus; the driver depends
// only on the driver.ImplOps interface, so internal/driver never imports
// internal/app.
type driverImplOps struct {
	base            runner.Options // the round-01 runOpts (Root, RunID, Store, Agents, Idea, Timeout)
	root            string
	ideaSlug        string
	ideaDir         string
	implementer     string   // FINAL drafter / first participant
	reviewers       []string // non-implementer participants
	drafter         string   // review-consensus drafter (facilitator)
	out             io.Writer
	verificationCLI string // internal process-fixture seam; empty uses this running CLI
	roleErr         string // declared-facilitator role deadlock; every role action escalates
	facilitator     protocol.FacilitatorRole
	// Designated-implementer state (meta-protocol-change-designated-implementer).
	// All four stay zero on an undesignated deck, which keeps today's behavior
	// byte-identical (nothing emitted, nothing compared, no gate).
	implSource     protocol.ImplementerSource // resolution rank that produced implementer
	implDesignated bool                       // a tier-2/tier-3 designation is present (drives the R45 event/line)
	implLine       string                     // the extra designation-path stdout line
	dispatchErr    string                     // tier-2 availability gate; escalated by the dispatch actions only (R17)
}

func newDriverImplOps(base runner.Options, root, ideaSlug, ideaDir string, participants []string, out io.Writer) driver.ImplOps {
	// Declared-facilitator runs (lean-organizer A): the declared facilitator is
	// ineligible for every code-facing role — implementer, reviewer, review-consensus
	// drafter, and (via the drafter) goal-done checker — unless
	// facilitator_participates: true opted it back in. When the predicate empties the
	// eligible set the ops carry a role error that every role action escalates with;
	// the driver NEVER silently falls back to facilitator implementation.
	role := protocol.ReadFacilitatorRole(ideaDir)
	roleErr := ""
	eligible := make([]string, 0, len(participants))
	for _, p := range participants {
		if role.IneligibleForRoles(p) {
			continue
		}
		eligible = append(eligible, p)
	}
	if role.Declared && !role.Participates && len(eligible) == 0 && len(participants) > 0 {
		roleErr = fmt.Sprintf("declared facilitator %s is the only participant; a declared-facilitator run needs at least one non-facilitator participant to implement — add a participant or set facilitator_participates: true (escalated, not fallen back)", role.Facilitator)
	}
	// Designated implementer (meta-protocol-change-designated-implementer): resolve
	// through the four-rank chain — pin → per-idea designation → live layered global
	// default → today's chain. Validity and pin-conflict gates ride the roleErr path
	// so every role action escalates (R25); the tier-2 availability gate is a
	// dispatch gate only (R17/R18). With no designation present anywhere this is
	// exactly today's resolveImplementer call.
	designation := resolveDispatchDesignation(root, ideaDir, eligible, base.Agents)
	implementer := designation.implementer
	if designation.gate != "" {
		if roleErr != "" {
			roleErr += "; " + designation.gate
		} else {
			roleErr = designation.gate
		}
	}
	// Dedupe to distinct non-implementer IDs (review CF1). Duplicate participant
	// IDs (e.g. [impl, rev, rev]) would otherwise (a) inflate ReviewerCount so the
	// LE-11 `< 2` guard passes with a single real reviewer, and (b) make
	// RunReviewRound launch one goroutine per duplicate, all writing the same
	// agents/<id>/stdout.log and review/round-NN/<id>.md concurrently (a write race).
	var reviewers []string
	seen := make(map[string]bool)
	for _, p := range eligible {
		if p != implementer && !seen[p] {
			seen[p] = true
			reviewers = append(reviewers, p)
		}
	}
	// Track-aware reviewer cap (idea track-aware-driver): an EXPLICIT fast/standard
	// track reduces the reviewer set per §4.0 (fast=1, standard=2); absent and
	// deliberation keep all non-implementers. The truncation is deterministic
	// (participant order) and never drops below the non-solo floor of 1.
	if t, present := driver.ReadTrack(ideaDir); present {
		if pol, err := track.PolicyFor(t, present, len(reviewers), driver.ReadAutoImplement(ideaDir), driver.ReadStrictGate(ideaDir)); err == nil && pol.MaxReviewers > 0 && len(reviewers) > pol.MaxReviewers {
			reviewers = reviewers[:pol.MaxReviewers]
		}
	}
	// The review-consensus drafter MUST be a non-implementer so the implementer
	// cannot filter reviewer findings out of the consensus (AF3). Fall back to the
	// implementer only if there are somehow no reviewers.
	drafter := implementer
	if len(reviewers) > 0 {
		drafter = reviewers[0]
	}
	ops := driverImplOps{
		base: base, root: root, ideaSlug: ideaSlug, ideaDir: ideaDir,
		implementer: implementer, reviewers: reviewers, drafter: drafter, out: out,
		roleErr: roleErr, facilitator: role,
		implSource: designation.source, implDesignated: designation.present,
		implLine: designation.line, dispatchErr: designation.dispatchGate,
	}
	if designation.present {
		ops.kickoffDesignationChecks()
	}
	return ops
}

// dispatchDesignation is the designation-aware resolution of who implements, computed
// once per driverImplOps construction and frozen for that dispatch (R29).
type dispatchDesignation struct {
	implementer  string
	source       protocol.ImplementerSource
	present      bool   // a tier-2/tier-3 designation is present (R45 event/line gate)
	line         string // the extra designation-path stdout line
	gate         string // validity / pin-conflict gate — roleErr path, any run (R16/R25/R28)
	dispatchGate string // tier-2 availability gate — dispatch actions only (R17/R18)
}

// resolveDispatchDesignation implements the R13 chain: rank 1 pin, rank 2 per-idea
// designation, rank 3 live layered global default, rank 4 today's chain verbatim.
// Ranks 2 and 3 are dormant unless somebody set them; with neither set the result is
// exactly resolveImplementer(ideaDir, eligible) and present=false.
func resolveDispatchDesignation(root, ideaDir string, eligible []string, discovered []agents.Discovery) dispatchDesignation {
	legacy := resolveImplementer(ideaDir, eligible)
	pinID, _, pinOK := protocol.ResolveImplementerChain(ideaDir, protocol.PinImplementerCandidates(), eligible)
	d := protocol.ReadImplementerDesignation(ideaDir)
	resolvedLine := func(id string, src protocol.ImplementerSource) string {
		return fmt.Sprintf("driver: implementer resolved: %s (source: %s)", id, src)
	}
	fallThroughNotice := func(id, reason string) string {
		return fmt.Sprintf("driver: designated implementer %s %s; falling through to the default chain — to gate instead: set per-idea `implementer: <other-id>`, write `implementer: none`, or record `implementer_waived: %s — <reason> — confirmed <date>`", id, reason, id)
	}
	switch d.State {
	case protocol.DesignationEmpty:
		// R2/R16: present-empty is an INCOMPLETE designation — a typo, never the
		// opt-out. Hard gate on any run; both legal spellings named.
		return dispatchDesignation{implementer: legacy, gate: "incomplete designation: 00-prompt.md carries `implementer:` with an empty value — write `implementer: <agent-id>` naming an eligible participant, or `implementer: none` for an explicit opt-out"}
	case protocol.DesignationNone:
		// R2: explicit per-idea non-designation — tier 3 suppressed, today's chain, no gate.
		return dispatchDesignation{implementer: legacy, source: protocol.SourceImplementerNone, present: true,
			line: fmt.Sprintf("driver: implementer designation declined (implementer: none); using the default chain (source: %s)", protocol.SourceImplementerNone)}
	case protocol.DesignationSet:
		id := d.ID
		if !memberOf(eligible, id) {
			// R16/R20: at tier 2 a non-participant AND the idea's own declared
			// non-participating facilitator are the same hard gate — both are absent
			// from the facilitator-filtered eligible set. Malformed values (R4) land
			// here too: they are never repaired, they fail this membership check.
			return dispatchDesignation{implementer: legacy, gate: fmt.Sprintf("invalid designation: `implementer: %s` is not an eligible participant of this idea — edit it to an eligible participant or write `implementer: none`", id)}
		}
		if pinOK && pinID != id {
			// R28: pin and designation disagree — escalate, never silently honour
			// either. Only a confirmed owner/author record retires the pin.
			if meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md")); err != nil || !protocol.ParseImplementerReassignment(meta, pinID, id) {
				return dispatchDesignation{implementer: legacy, gate: fmt.Sprintf("implementer conflict: IMPLEMENTATION.md pins %s but 00-prompt.md designates %s — either restore the designation to match the pin, or record `implementer_reassigned: %s to %s — <reason> — confirmed <date>` (never the incoming implementer's own edit)", pinID, id, pinID, id)}
			}
		}
		src := protocol.SourceImplementerDesignation
		if pinOK {
			src = protocol.SourceImplementerPin // rank 1 wins when tiers 1 and 2 agree (R27)
		}
		if !designeeAvailable(root, id, discovered) {
			if meta, err := protocol.ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md")); err == nil && protocol.ImplementerWaived(meta, id) {
				// R18 exit 2, recorded: the confirmed waiver clears the gate.
				return dispatchDesignation{implementer: legacy, source: protocol.SourceImplementerFallThroughUnavailable, present: true,
					line: fallThroughNotice(id, "is waived by a recorded `implementer_waived:` line")}
			}
			// R18: valid but unavailable at the ping — blocking dispatch gate with
			// the Confirm pre-built (three recorded exits). Dispatch actions only.
			return dispatchDesignation{implementer: id, source: src, present: true, line: resolvedLine(id, src),
				dispatchGate: fmt.Sprintf("designated implementer %s is unavailable (§9.0 ping) — exits: (1) re-designate: edit `implementer:` to another eligible participant; (2) record `implementer_waived: %s — <reason> — confirmed <date>`; (3) write `implementer: none`", id, id)}
		}
		return dispatchDesignation{implementer: id, source: src, present: true, line: resolvedLine(id, src)}
	}
	// DesignationAbsent — rank 1 pin governs re-entry (R27/R30)…
	if pinOK {
		if g := globalDefaultImplementer(root); g != "" {
			// …and a present tier-3 designation is still recorded as present (R45).
			return dispatchDesignation{implementer: pinID, source: protocol.SourceImplementerPin, present: true, line: resolvedLine(pinID, protocol.SourceImplementerPin)}
		}
		return dispatchDesignation{implementer: pinID, source: protocol.SourceImplementerPin}
	}
	// …then rank 3: the live layered global default (R10/R11).
	g := globalDefaultImplementer(root)
	if g == "" {
		return dispatchDesignation{implementer: legacy} // unset path: byte-identical
	}
	if strings.EqualFold(g, "none") {
		// R9: "none" is the documented deck-wide suppressor — today's chain, recorded.
		return dispatchDesignation{implementer: legacy, source: protocol.SourceImplementerNone, present: true,
			line: fmt.Sprintf("driver: standing implementer default suppressed (default_implementer = \"none\"); using the default chain (source: %s)", protocol.SourceImplementerNone)}
	}
	if strings.ContainsAny(g, " \t\r\n") {
		// R19 concession 3 / AC-10: a malformed tier-3 value keeps its hard failure;
		// only unavailability and inapplicability fall through.
		return dispatchDesignation{implementer: legacy, gate: fmt.Sprintf("invalid default_implementer %q: not a single agent id — fix the config value or set it to \"none\"", g)}
	}
	if !memberOf(eligible, g) {
		// R20–R22: a standing preference legitimately predates this idea's roster and
		// facilitator declaration — inapplicable, not invalid. Notice + fall through.
		return dispatchDesignation{implementer: legacy, source: protocol.SourceImplementerFallThroughInapplicable, present: true,
			line: fallThroughNotice(g, "is not an eligible participant of this idea (inapplicable)")}
	}
	if !designeeAvailable(root, g, discovered) {
		// R19: tier-3 unavailability is a one-line notice and fall-through, no gate.
		return dispatchDesignation{implementer: legacy, source: protocol.SourceImplementerFallThroughUnavailable, present: true,
			line: fallThroughNotice(g, "is unavailable (§9.0 ping)")}
	}
	return dispatchDesignation{implementer: g, source: protocol.SourceImplementerGlobalDefault, present: true, line: resolvedLine(g, protocol.SourceImplementerGlobalDefault)}
}

func memberOf(list []string, id string) bool {
	for _, p := range list {
		if p == id {
			return true
		}
	}
	return false
}

// designeeAvailable is the dispatch-time availability signal behind the R18/R19 ping
// checks: the designee resolves to a discovered (Found) agent on this machine.
func designeeAvailable(root, id string, discovered []agents.Discovery) bool {
	d, err := agents.ResolveParticipant(id, discovered, rosterMappingFor(root))
	return err == nil && d.Found
}

// globalDefaultImplementer is rank 3's live layered read (R10/R11). A layered-config
// read error leaves tier 3 unset rather than gating dispatch on a parse failure that
// the config-loading paths already surface.
func globalDefaultImplementer(root string) string {
	defs, err := config.LoadDefaults(root)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(defs.DefaultImplementer)
}

// kickoffDesignationChecks runs the designation-only kickoff surfacing (R36–R39): the
// unchanged LE-3 model-diversity check runs early so an unsatisfiable roster costs an
// error message rather than a design cycle, and thin reviewer benches warn without
// blocking. Fires only under a present designation; the OpenReviewRound check at
// review time is unchanged.
func (o driverImplOps) kickoffDesignationChecks() {
	if o.out == nil {
		return
	}
	// R37: same check, same severity, no new gate class — an escalation here prints
	// early and still gates later at OpenReviewRound, unchanged.
	if err := o.checkModelDiversity(); err != nil {
		fmt.Fprintf(o.out, "driver: kickoff model-diversity check: %v\n", err)
	}
	switch len(o.reviewers) {
	case 0:
		// R39 surfaced earlier: the hard stop itself stays at OpenReviewRound.
		fmt.Fprintf(o.out, "driver: WARNING designated run has no non-implementer reviewer; review cannot open (the hard stop is unchanged)\n")
	case 1:
		// R38: a two-participant designated idea warns, never blocks.
		fmt.Fprintf(o.out, "driver: WARNING designated run leaves a single non-implementer reviewer (%s); %s implements and never reviews itself — consider a third participant for review depth\n", o.reviewers[0], o.implementer)
	}
}

// checkImplementerReentry is R29: before a pin exists, the recorded
// agent.implementer_resolved event is the durable record of what was dispatched, and a
// live tier-2/tier-3 change against it escalates rather than reassigning silently. The
// comparison fires only when the recorded source was designation or global-default; a
// pinned resolution governs every resume (R30) and never reaches here.
func (o driverImplOps) checkImplementerReentry() error {
	if o.base.Store == (store.Store{}) {
		return nil
	}
	// The comparison is a PRE-PIN instrument: once IMPLEMENTATION.md exists the pin
	// governs every resume (R30) and there is nothing to compare against.
	if o.implSource == protocol.SourceImplementerPin {
		return nil
	}
	events, err := o.base.Store.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("driver: cannot replay the run's dispatch record: %v", err)
	}
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		if e.Type != "agent.implementer_resolved" {
			continue
		}
		if idea, _ := e.Data["idea"].(string); idea != o.ideaSlug {
			continue
		}
		recSrc, _ := e.Data["source"].(string)
		if recSrc != string(protocol.SourceImplementerDesignation) && recSrc != string(protocol.SourceImplementerGlobalDefault) {
			continue
		}
		recID, _ := e.Data["implementer"].(string)
		if recID == o.implementer && recSrc == string(o.implSource) {
			return nil
		}
		return fmt.Errorf("driver: designation changed after a recorded dispatch (recorded: %s via %s; now: %s via %s) — restore the designation to match the record, or record `implementer_reassigned: %s to %s — <reason> — confirmed <date>`; escalating rather than reassigning silently", recID, recSrc, o.implementer, o.implSource, recID, o.implementer)
	}
	return nil
}

// resolveImplementer picks the implementer from durable role metadata (D10/AF6):
// IMPLEMENTATION.md `implementer` (on re-entry), else FINAL.md `implementer` /
// `drafted-by`, validated against participants; otherwise participants[0]. This is
// rank 1 + rank 4 of the R13 chain — today's behavior, byte-identical, pinned by
// TestResolveImplementerFromRoleMetadata (R33). The shared implementation lives in
// internal/protocol (R31); the eligibility list is the caller's explicit parameter
// (R34). Tiers 2–3 never enter this function.
func resolveImplementer(ideaDir string, participants []string) string {
	if id, _, ok := protocol.ResolveImplementerChain(ideaDir, protocol.LegacyImplementerCandidates(), participants); ok {
		return id
	}
	if len(participants) > 0 {
		return participants[0]
	}
	return ""
}

// withParticipants clones the base run options with a narrowed participant set.
func (o driverImplOps) withParticipants(ids ...string) runner.Options {
	opts := o.base
	opts.Idea.Participants = append([]string(nil), ids...)
	return opts
}

// modelOf returns the configured model id for an agent id, or "" if unknown. It
// matches the id directly OR its mapped family (claude-1 -> claude) so the LE-3
// model-diversity gate is not silently disabled for a roster-id deck (review MINOR,
// kimi-1). It does NOT require the agent to be installed — the configured model
// counts for the diversity comparison regardless of discovery.
func (o driverImplOps) modelOf(id string) string {
	family := rosterMappingFor(o.root)[id]
	for _, a := range o.base.Agents {
		if a.ID == id || (family != "" && a.ID == family) {
			return a.Model
		}
	}
	return ""
}

// reviewersShareImplementerModel reports whether every reviewer uses the same model as
// the implementer (LE-3). Returns (false, "") when the implementer's model is unknown
// or there are no reviewers — diversity can't be asserted, so it never fires spuriously.
func (o driverImplOps) reviewersShareImplementerModel() (bool, string) {
	implModel := o.modelOf(o.implementer)
	if implModel == "" || len(o.reviewers) == 0 {
		return false, ""
	}
	for _, r := range o.reviewers {
		if !strings.EqualFold(o.modelOf(r), implModel) { // review fix F6: case-insensitive
			return false, ""
		}
	}
	return true, implModel
}

// checkModelDiversity (LE-3) emits the agent.model_diversity event and either warns
// (default) or, with require_model_diversity, returns an escalation error. Extracted so
// OpenReviewRound stays thin and the event/warn/escalate logic is unit-testable without
// launching reviewers (review fixes F4 + F8).
func (o driverImplOps) checkModelDiversity() error {
	same, model := o.reviewersShareImplementerModel()
	if !same {
		return nil
	}
	required := driver.ReadRequireModelDiversity(o.ideaDir)
	// §4.0: the fast track's single reviewer MUST be model-diverse — a same-model
	// solo checker is the highest rubber-stamp risk (review-01 F4). Make it a hard
	// gate on fast regardless of the frontmatter flag.
	if t, present := driver.ReadTrack(o.ideaDir); present && t == track.Fast {
		required = true
	}
	action := "warn"
	if required {
		action = "escalate"
	}
	if o.base.Store != (store.Store{}) { // F4: best-effort durable event for TUI/state consumers
		_ = o.base.Store.Append(store.Event{
			Time: time.Now().UTC(),
			Type: "agent.model_diversity",
			Data: map[string]any{
				"idea": o.ideaSlug, "implementer": o.implementer,
				"reviewers": append([]string(nil), o.reviewers...),
				"model":     model, "required": required, "action": action,
			},
		})
	}
	if required {
		return fmt.Errorf("require_model_diversity: every reviewer shares the implementer's model %q; refusing to open review (LE-3)", model)
	}
	fmt.Fprintf(o.out, "driver: WARNING model-diversity — every reviewer shares the implementer's model %q; a same-model checker is more likely to rubber-stamp (LE-3). Set require_model_diversity: true to make this a hard gate.\n", model)
	return nil
}

func (o driverImplOps) Implement(ctx context.Context) error {
	if o.roleErr != "" {
		return fmt.Errorf("driver: %s", o.roleErr)
	}
	if o.dispatchErr != "" {
		return fmt.Errorf("driver: %s", o.dispatchErr)
	}
	fmt.Fprintf(o.out, "driver: implementing via %s ...\n", o.implementer)
	if o.implDesignated {
		// R29: a live tier-2/tier-3 change against a recorded designation dispatch
		// escalates rather than reassigning silently.
		if err := o.checkImplementerReentry(); err != nil {
			return err
		}
		// R45/R46: designation-only observability — one extra stdout line and the
		// durable resolved-source event, guarded exactly as checkModelDiversity's
		// event. On an undesignated deck neither fires and the line above stays
		// byte-identical to today.
		fmt.Fprintln(o.out, o.implLine)
		if o.base.Store != (store.Store{}) {
			_ = o.base.Store.Append(store.Event{
				Time: time.Now().UTC(),
				Type: "agent.implementer_resolved",
				Data: map[string]any{
					"idea": o.ideaSlug, "implementer": o.implementer,
					"source": string(o.implSource),
				},
			})
		}
	}
	r := runner.RunImplementation(ctx, o.withParticipants(o.implementer))
	if !r.Success() {
		return fmt.Errorf("implementer %s: %s", r.AgentID, r.ExitError)
	}
	return nil
}

func (o driverImplOps) ImplementationStatus() (string, error) {
	meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "IMPLEMENTATION.md"))
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.TrimSpace(meta["status"]), `"'`), nil
}

// RunChecks runs the verification gate (LE-4). Resolution order:
//  1. an explicit `checks:` command from 00-prompt frontmatter → `sh -c <command>`;
//  2. else `go test ./...` when the workspace is a Go module;
//  3. else, for a code-writing (auto_implement) idea with no checks → FAIL CLOSED;
//  4. else (design-only, non-Go) → nothing to check.
//
// Because advanceImpl (pre-review) and advanceReview (post-fix-up) both escalate when
// RunChecks fails, step 3 transitively ties the "artifact-wins" fix-up override to a
// real check: a fix-up that wrote a valid-shaped artifact but cannot be verified no
// longer auto-passes (hermes #8).
func (o driverImplOps) RunChecks(ctx context.Context) (bool, string) {
	run := func(name string, cmd *exec.Cmd) (bool, string) {
		fmt.Fprintf(o.out, "driver: running checks (%s) ...\n", name)
		cmd.Dir = o.root
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err := cmd.Run()
		return err == nil, buf.String()
	}
	// List-form `checks:` = the completion contract (completion-contracts-evidence-ledger):
	// run each criterion, write the evidence table into IMPLEMENTATION.md, fail closed on
	// any non-zero exit. A malformed list is a hard fail (present but invalid).
	if criteria, isList, err := driver.ReadChecksContract(o.ideaDir); err != nil {
		return false, "checks: contract invalid — " + err.Error()
	} else if isList {
		return o.runChecksContract(ctx, criteria)
	}
	checks := ""
	if meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "00-prompt.md")); err == nil {
		checks = strings.TrimSpace(strings.Trim(strings.TrimSpace(meta["checks"]), `"'`))
	}
	if checks != "" {
		return run(checks, exec.CommandContext(ctx, "sh", "-c", checks))
	}
	if _, err := os.Stat(filepath.Join(o.root, "go.mod")); err == nil {
		return run("go test ./...", exec.CommandContext(ctx, "go", "test", "./..."))
	}
	if driver.ReadAutoImplement(o.ideaDir) {
		return false, "no go.mod and no `checks:` configured for a code-writing (auto_implement) idea; set checks: in 00-prompt.md so a fix-up cannot pass unverified"
	}
	return true, "no go.mod and no checks configured; nothing to check"
}

func (o driverImplOps) OpenReviewRound(ctx context.Context, round int) error {
	if o.roleErr != "" {
		return fmt.Errorf("driver: %s", o.roleErr)
	}
	if len(o.reviewers) == 0 {
		return fmt.Errorf("no non-implementer reviewers available")
	}
	// LE-3 model-diversity: a checker that shares the implementer's model is more
	// likely to rubber-stamp. Default = warn (+event); require_model_diversity escalates.
	if err := o.checkModelDiversity(); err != nil {
		return err
	}
	fmt.Fprintf(o.out, "driver: opening review round %d (reviewers: %s) ...\n", round, strings.Join(o.reviewers, ", "))
	// AF5: drop any reviewer artifact that exists but fails validation, so
	// RunReviewRound (Overwrite=false) regenerates it instead of skipping a
	// malformed file forever (which would spin the driver to the deadline).
	dir := filepath.Join(o.ideaDir, "review", roundDirLabel(round))
	for _, reviewer := range o.reviewers {
		path := filepath.Join(dir, reviewer+".md")
		if _, err := os.Stat(path); err == nil {
			if runner.ValidateReviewArtifact(path, reviewer, o.ideaSlug, round) != nil {
				_ = os.Remove(path)
			}
		}
	}
	opts := o.withParticipants(o.reviewers...)
	opts.Round = round
	results := runner.RunReviewRound(ctx, opts)
	failed := 0
	for _, r := range results {
		if r.ExitError != "" {
			failed++
		}
	}
	if failed == len(results) && failed > 0 {
		return fmt.Errorf("all reviewers failed in review round %d", round)
	}
	return nil
}

func (o driverImplOps) ReviewRoundComplete(round int) (bool, error) {
	dir := filepath.Join(o.ideaDir, "review", roundDirLabel(round))
	for _, reviewer := range o.reviewers {
		path := filepath.Join(dir, reviewer+".md")
		if _, err := os.Stat(path); err != nil {
			return false, nil
		}
		if err := runner.ValidateReviewArtifact(path, reviewer, o.ideaSlug, round); err != nil {
			return false, nil
		}
	}
	return len(o.reviewers) > 0, nil
}

func (o driverImplOps) DraftReviewConsensus(ctx context.Context, round int) error {
	fmt.Fprintf(o.out, "driver: drafting review consensus via %s ...\n", o.drafter)
	strict := driver.ReadStrictGate(o.ideaDir)
	opts := o.withParticipants(o.drafter)
	opts.Round = round
	opts.StrictGate = strict // LE-2: emit the close fields under strict_gate
	r := runner.RunReviewConsensus(ctx, opts)
	if !r.Success() {
		return fmt.Errorf("review-consensus drafter %s: %s", r.AgentID, r.ExitError)
	}
	path := filepath.Join(o.ideaDir, "review", "consensus.md")
	if err := runner.ValidateReviewConsensusArtifact(path); err != nil {
		return err
	}
	// Review fix F7: under strict_gate, the drafter MUST emit the close fields. If they
	// are absent, fail fast here rather than spinning fresh rounds to MaxFixupCycles.
	if strict {
		meta, err := protocol.ReadFrontmatter(path)
		if err != nil {
			return err
		}
		if strings.TrimSpace(meta["strict_gate_clean"]) == "" || strings.TrimSpace(meta["closing_review_round"]) == "" {
			return fmt.Errorf("strict_gate: review/consensus.md must set strict_gate_clean and closing_review_round")
		}
	}
	return nil
}

func (o driverImplOps) ReviewStatus() (driver.ReviewStatus, error) {
	summary, err := consensus.Status(o.root, o.ideaSlug, true)
	if err != nil {
		return driver.ReviewStatus{}, err
	}
	meta, err := protocol.ReadFrontmatter(filepath.Join(o.ideaDir, "review", "consensus.md"))
	if err != nil {
		return driver.ReviewStatus{}, err
	}
	rawFixes := strings.Trim(strings.TrimSpace(meta["outstanding_agreed_fixes"]), `"'`)
	fixes, err := strconv.Atoi(rawFixes)
	if err != nil || fixes < 0 {
		return driver.ReviewStatus{}, fmt.Errorf("review/consensus.md outstanding_agreed_fixes=%q is not a non-negative integer", rawFixes)
	}
	blocked := strings.EqualFold(strings.Trim(strings.TrimSpace(meta["blocked"]), `"'`), "true")
	// LE-2 strict_gate close fields (absent/zero on non-strict ideas → harmless defaults).
	strictClean := strings.EqualFold(strings.Trim(strings.TrimSpace(meta["strict_gate_clean"]), `"'`), "true")
	closingRound := 0
	if v := strings.Trim(strings.TrimSpace(meta["closing_review_round"]), `"'`); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			closingRound = n
		}
	}
	return driver.ReviewStatus{
		Summary: summary, OutstandingAgreedFixes: fixes, Blocked: blocked,
		StrictGateClean: strictClean, ClosingReviewRound: closingRound,
		ReviewerCount: len(o.reviewers), // LE-11
	}, nil
}

// discoveryFor returns the discovered agent for an id, matching the id directly OR
// its mapped family (claude-1 -> claude) so the LE-7 goal-check is not silently
// skipped for a roster-id deck (review MINOR, kimi-1).
func (o driverImplOps) discoveryFor(id string) (agents.Discovery, bool) {
	family := rosterMappingFor(o.root)[id]
	for _, a := range o.base.Agents {
		if a.ID == id || (family != "" && a.ID == family) {
			return a, true
		}
	}
	return agents.Discovery{}, false
}

// GoalCheck (LE-7) runs a fresh non-implementer agent (the review drafter) to check the
// FINAL.md acceptance criteria before close, reusing the consult execution path with a
// verdict prompt. Missing, self, failed or ambiguous execution cannot establish
// completion. A textual pass remains defense in depth, not criterion evidence.
func (o driverImplOps) GoalCheck(ctx context.Context) (bool, string) {
	if o.roleErr != "" {
		return false, "goal-check unavailable: " + o.roleErr
	}
	checker := o.drafter
	// CF6: GoalCheck must use a non-implementer checker. The upstream guards
	// (ReviewerCount < 2 under auto; OpenReviewRound under strict) already prevent
	// the drafter==implementer fallback from reaching here, but enforce the contract
	// locally too — never run the implementer as its own goal checker.
	if checker == "" || checker == o.implementer {
		return false, "goal-check has no independent checker"
	}
	agent, err := agents.ResolveParticipant(checker, o.base.Agents, rosterMappingFor(o.root))
	if err != nil {
		return false, "goal-check checker unavailable"
	}
	fmt.Fprintf(o.out, "driver: goal-done check via %s ...\n", checker)
	dir := filepath.Join(o.root, protocol.DeckDir, "runs", o.base.RunID, "agents", checker)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, "goal-check cannot create its evidence directory"
	}
	ctx = runner.WithLaunchInfo(ctx, runner.LaunchInfo{RunID: o.base.RunID,
		Idea: o.ideaSlug, Phase: "goal-check", Store: o.base.Store})
	res := runner.RunConsult(ctx, runner.ConsultOptions{
		Root:  o.root,
		Agent: agent,
		// Keep the existing bounded one-shot deadline. A timeout leaves
		// completion unverified and halts instead of silently passing.
		Timeout:    2 * time.Minute,
		Prompt:     runner.BuildGoalCheckPrompt(agent, o.base.Idea),
		StdoutPath: filepath.Join(dir, "goal-check.stdout.log"),
		StderrPath: filepath.Join(dir, "goal-check.stderr.log"),
		Progress:   o.out,
	})
	if res.ExitError != "" || res.AgentExit != 0 {
		return false, "goal-check checker failed; completion is unverified"
	}
	switch parseGoalVerdict(res.Answer) {
	case "FAIL":
		return false, res.Answer
	case "PASS":
		return true, ""
	default:
		return false, "goal-check inconclusive; completion is unverified"
	}
}

// parseGoalVerdict accepts PASS only when every stated verdict is an exact PASS.
// A failure is sticky; any unknown verdict makes a pass ambiguous. A trailing
// format example therefore cannot erase an earlier failure or reservation.
func parseGoalVerdict(answer string) string {
	seenPass, seenFail, ambiguous := false, false, false
	for _, line := range strings.Split(answer, "\n") {
		t := strings.ToUpper(strings.TrimSpace(line))
		// CF2: strip leading markdown / quote wrappers (heading, bold, blockquote,
		// underscore, backtick, single/double quote) so a wrapped marker like
		// `GOAL-CHECK: FAIL` or **GOAL-CHECK:** still matches the prefix.
		t = strings.TrimLeft(t, "#*->_ \t`\"'")
		if !strings.HasPrefix(t, "GOAL-CHECK:") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(t, "GOAL-CHECK:"))
		// CF2: a bolded/quoted marker ("**GOAL-CHECK:** FAIL") leaves "** FAIL" in
		// rest — strip the leading wrapper run before the PASS/FAIL prefix check.
		rest = strings.Trim(rest, "*`\"'_ ")
		switch {
		case rest == "PASS":
			seenPass = true
		case strings.HasPrefix(rest, "FAIL"):
			seenFail = true
		default:
			ambiguous = true
		}
	}
	if seenFail {
		return "FAIL"
	}
	if seenPass && !ambiguous {
		return "PASS"
	}
	return ""
}

func (o driverImplOps) RequestReviewSignoffs(ctx context.Context, missing []string) error {
	return requestConsensusSignoffs(ctx, requestSignoffsOptions{
		Root:            o.root,
		IdeaSlug:        o.ideaSlug,
		Review:          true,
		ParticipantsRaw: strings.Join(missing, ","),
		Yes:             true,
	}, o.out, o.out)
}

func (o driverImplOps) PrecheckFixup(ctx context.Context) error {
	return runner.PrecheckFixup(ctx, o.withParticipants(o.implementer))
}

func (o driverImplOps) Fixup(ctx context.Context, cycle int) error {
	if o.roleErr != "" {
		return fmt.Errorf("driver: %s", o.roleErr)
	}
	if o.dispatchErr != "" {
		return fmt.Errorf("driver: %s", o.dispatchErr)
	}
	fmt.Fprintf(o.out, "driver: running fix-up cycle %d via %s ...\n", cycle, o.implementer)
	r := runner.RunFixup(ctx, o.withParticipants(o.implementer))
	if !r.Success() {
		return fmt.Errorf("fix-up implementer %s: %s", r.AgentID, r.ExitError)
	}
	return nil
}

// Complete writes IMPLEMENTATION.md status=complete. This is a deterministic file
// write by the orchestrator (NOT an implementer agent), so an implementer cannot
// short-circuit review (consensus D5).
func (o driverImplOps) Complete(ctx context.Context) error {
	return evidence.WithReportWriter(ctx, o.ideaDir, func(_ *evidence.ReportWriter) error { return o.completeWithWriter(ctx) })
}

func (o driverImplOps) completeWithWriter(ctx context.Context) error {
	if err := trajectory.RequireResolved(ctx, o.root, o.ideaSlug); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	pin := ""
	if o.base.Store.Enabled() {
		cursor, err := driver.LoadCursor(filepath.Join(o.base.Store.Directory(), "driver.json"))
		if err == nil {
			pin = cursor.ChecksContractSHA256
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	contract, err := driver.ObserveChecksContract(o.ideaDir, pin)
	if err != nil {
		return err
	} else if contract != "" {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("independent evidence completion requires a POSIX execution host; Windows runtime is not supported")
		}
		if err := o.requireAcceptedVerification(); err != nil {
			return err
		}
		gate := o.EvidenceCloseGate(o.drafter)
		if !gate.Allowed {
			return fmt.Errorf("independent evidence changed before completion: %s", strings.Join(gate.Reasons, "; "))
		}
	}
	path := filepath.Join(o.ideaDir, "IMPLEMENTATION.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if contract != "" {
		completed, _, err := evidence.TransitionStatusToComplete(data)
		if err != nil {
			return err
		}
		report, err := evidence.Load(o.ideaDir)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(o.root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if err := verifyValidationEvidence(data, rel, report); err != nil {
			return err
		}
		beforeRest, _, err := splitValidationEvidence(data)
		if err != nil {
			return err
		}
		if report.ExtraDigests[rel] != sha256Hex(string(beforeRest)) {
			return fmt.Errorf("implementation changed before completion write")
		}
		afterRest, _, err := splitValidationEvidence(completed)
		if err != nil {
			return err
		}
		if reasons := evidence.VerifyCompletionTransition(report, rel, report.ExtraDigests[rel], afterRest, o.drafter); len(reasons) > 0 {
			return fmt.Errorf("completion status was not independently authorized: %s", strings.Join(reasons, "; "))
		}
		return writeVerificationBytes(path, completed)
	}
	lines := strings.Split(string(data), "\n")
	inFrontmatter := false
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			if !replaced {
				lines = append(lines[:i], append([]string{"status: complete"}, lines[i:]...)...)
				replaced = true
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "status:") {
			lines[i] = "status: complete"
			replaced = true
		}
	}
	if !replaced {
		return fmt.Errorf("%s has no frontmatter status field", path)
	}
	return writeVerificationBytes(path, []byte(strings.Join(lines, "\n")))
}

func roundDirLabel(n int) string { return fmt.Sprintf("round-%02d", n) }
