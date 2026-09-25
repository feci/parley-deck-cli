package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"parley-deck-cli/internal/protocol"

	"parley-deck-cli/internal/agents"
	"parley-deck-cli/internal/consensus"
	"parley-deck-cli/internal/driver"
)

// driverConsensusOps is the production driver.ConsensusOps adapter (slice 2.2). It
// lives in internal/app so it can reuse the existing request-signoffs path and the
// headless agent invoker; the driver depends only on the driver.ConsensusOps
// interface, so internal/driver never imports internal/app (the import-direction
// guarantee D9 sought, achieved by injection instead of extraction).
type driverConsensusOps struct {
	root         string
	ideaSlug     string
	ideaDir      string
	participants []string
	discovered   []agents.Discovery
	out          io.Writer
}

func newDriverConsensusOps(root, ideaSlug, ideaDir string, participants []string, discovered []agents.Discovery, out io.Writer) driver.ConsensusOps {
	return driverConsensusOps{root: root, ideaSlug: ideaSlug, ideaDir: ideaDir, participants: participants, discovered: discovered, out: out}
}

func (o driverConsensusOps) Status() (consensus.Summary, error) {
	return consensus.Status(o.root, o.ideaSlug, false)
}

// Draft creates the consensus.md scaffold (+ sets idea status=consensus) if absent,
// then invokes a drafter agent to author the real synthesis into consensus.md.
func (o driverConsensusOps) Draft(ctx context.Context) error {
	path := filepath.Join(o.ideaDir, "consensus.md")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, derr := consensus.Draft(o.root, o.ideaSlug, consensus.DraftOptions{}); derr != nil {
			return derr
		}
	} else if err != nil {
		return err
	}
	return o.runDrafter(ctx, "consensus", buildConsensusDraftPrompt(o.ideaDir, path))
}

// RequestSignoffs invokes each missing participant through the existing
// request-signoffs path so every agent authors its own signoff.
func (o driverConsensusOps) RequestSignoffs(ctx context.Context, missing []string) error {
	return requestConsensusSignoffs(ctx, requestSignoffsOptions{
		Root:            o.root,
		IdeaSlug:        o.ideaSlug,
		ParticipantsRaw: strings.Join(missing, ","),
		Yes:             true,
	}, o.out, o.out)
}

// DraftFinal invokes a drafter agent to author the FINAL.md content directly. It
// deliberately does NOT call consensus.Finalize: Finalize sets idea status=final
// before the content exists, which would strand the idea at status=final with a
// scaffold if the drafter failed (AF1). The driver commits the idea status to
// "final" only AFTER validating the authored content (D7).
func (o driverConsensusOps) DraftFinal(ctx context.Context) error {
	path := filepath.Join(o.ideaDir, "FINAL.md")
	return o.runDrafter(ctx, "FINAL", buildFinalDraftPrompt(o.ideaDir, path))
}

func (o driverConsensusOps) Reopen(ctx context.Context, reason string) error {
	_, err := consensus.Reopen(o.root, o.ideaSlug, consensus.ReopenOptions{Reason: reason})
	return err
}

// runDrafter invokes the first available headless agent to author the target file
// per the given prompt. Drafting is a single-agent facilitator action (D6).
func (o driverConsensusOps) runDrafter(ctx context.Context, kind, prompt string) error {
	role := protocol.ReadFacilitatorRole(o.ideaDir)
	// R36 drafter separation: under a live designation, prefer an eligible drafter
	// that is not the designated implementer (the designee remains the fallback).
	preferNot := designatedImplementerPreference(o.root, o.ideaDir, o.participants, role)
	drafter, ok := firstEligibleHeadlessAgentPreferring(o.discovered, o.participants, rosterMappingFor(o.root), role, preferNot)
	if !ok {
		return fmt.Errorf("no headless idea participant available to draft %s", kind)
	}
	rootAbs, err := filepath.Abs(o.root)
	if err != nil {
		return err
	}
	fmt.Fprintf(o.out, "driver: drafting %s via %s ...\n", kind, drafter.ID)
	return runHeadlessSignoffAgent(ctx, rootAbs, drafter, prompt, o.out, o.out)
}

// designatedImplementerPreference returns the designation that drafter selection
// should avoid where an alternative exists (R36): the eligible tier-2 designation,
// else the eligible tier-3 global default, else "". Inert without a designation — and
// a defective or inapplicable designation steers nothing.
func designatedImplementerPreference(root, ideaDir string, participants []string, role protocol.FacilitatorRole) string {
	eligible := make([]string, 0, len(participants))
	for _, p := range participants {
		if role.IneligibleForRoles(p) {
			continue
		}
		eligible = append(eligible, p)
	}
	d := protocol.ReadImplementerDesignation(ideaDir)
	if d.State == protocol.DesignationSet {
		if memberOf(eligible, d.ID) {
			return d.ID
		}
		return ""
	}
	if d.State != protocol.DesignationAbsent {
		return "" // `none` and the gating present-empty state suppress tier 3
	}
	g := globalDefaultImplementer(root)
	if g == "" || strings.EqualFold(g, "none") || strings.ContainsAny(g, " \t\r\n") || !memberOf(eligible, g) {
		return ""
	}
	return g
}

// firstHeadlessAgent returns the first discovered headless agent that is also an
// idea participant (Parley Deck §4/§6: the facilitator-drafter must be a
// participant of the deliberation, not an arbitrary installed agent — AF2).
func firstHeadlessAgent(discovered []agents.Discovery, participants []string, mapping map[string]string) (agents.Discovery, bool) {
	return firstEligibleHeadlessAgent(discovered, participants, mapping, protocol.FacilitatorRole{})
}

// firstEligibleHeadlessAgent additionally refuses the declared facilitator of a
// declared-facilitator run (lean-organizer A): the facilitator never drafts consensus
// or FINAL unless facilitator_participates: true opted it back into participation.
// With no eligible drafter left it reports !ok so the caller escalates rather than
// silently falling back to the facilitator.
func firstEligibleHeadlessAgent(discovered []agents.Discovery, participants []string, mapping map[string]string, role protocol.FacilitatorRole) (agents.Discovery, bool) {
	return firstEligibleHeadlessAgentPreferring(discovered, participants, mapping, role, "")
}

// firstEligibleHeadlessAgentPreferring adds the R36 drafter-separation preference:
// preferNot is the designated implementer (empty when no designation is present) —
// the first pass skips it, the second pass falls back to it, so drafter ==
// implementer is never forbidden, only dispreferred where an alternative exists.
func firstEligibleHeadlessAgentPreferring(discovered []agents.Discovery, participants []string, mapping map[string]string, role protocol.FacilitatorRole, preferNot string) (agents.Discovery, bool) {
	// Iterate participants in order and resolve each (roster id via [roster.*] or a
	// bare family id) so a roster-id deck finds its drafter (composite-agent-naming).
	for pass := 0; pass < 2; pass++ {
		for _, p := range participants {
			if role.IneligibleForRoles(p) {
				continue
			}
			if pass == 0 && preferNot != "" && p == preferNot {
				continue
			}
			if agent, err := agents.ResolveParticipant(p, discovered, mapping); err == nil {
				if agent.Found && agents.LaunchModeOrDefault(agent.LaunchMode) == agents.LaunchHeadless {
					return agent, true
				}
			}
		}
		if preferNot == "" {
			break
		}
	}
	return agents.Discovery{}, false
}

// buildConsensusDraftPrompt is generated from protocol.RequiredConsensusSections and
// protocol.ConditionalConsensusSections — the SAME constants the consensus scaffold
// generator reads — so the prompt cannot diverge from the scaffold handed to the
// drafter (lean-organizer A.4; fix-up F1: prompt+scaffold are the constant's only
// consumers, no status gate reads it; the previous hardcoding emitted review-cycle
// headings and none of the §15.3/§15.5/§15.6 duty sections).
func buildConsensusDraftPrompt(ideaDir, path string) string {
	var sections strings.Builder
	for _, section := range protocol.RequiredConsensusSections {
		sections.WriteString(section)
		sections.WriteString("\n")
	}
	conditions := make([]string, 0, len(protocol.ConditionalConsensusSections))
	for _, section := range protocol.ConditionalConsensusSections {
		conditions = append(conditions, "`"+section+"`")
	}
	return fmt.Sprintf(`You are the Parley Deck facilitator drafting the consensus.

Read EVERY round artifact under %s/round-*/ (each participant's round files).
The file %s currently holds a scaffold. OVERWRITE it with the real synthesis of the
deliberation, keeping this exact structure and the YAML frontmatter:

%s
Under "## Signoffs", leave one HTML comment placeholder line per participant exactly
like "<!-- <agent-id> appends its signoff below -->" for every participant in the
idea's 00-prompt.md, and NOTHING else under that heading (the agents append their own
✅/🟡/❌ blocks later).

§15 duties that bind on every track (§15.7): "## Drafter position changes" records
EVERY material change in your position since your most recent round file (or "None"),
each with an exact prior quotation, the prior position, the new position, and the
source round path; "## Alternatives disposition" gives one ALT- id per alternative with
an adopt or reject and the decisive reason; "## Comparison & blind spots" keeps raw
disagreements visible instead of smoothing them into trade-offs. Additionally add %s
VERBATIM — quoting each contradictory verdict with its author, tag and evidence and the
resolution — ONLY when contradictory verdicts exist; absent any conflict that section
does not exist. Be concrete and concise. English only. Write the file now and
report only the path.`, ideaDir, path, sections.String(), strings.Join(conditions, " and/or "))
}

// buildFinalDraftPrompt tells the drafter exactly what the gate will require.
//
// The prompt asked for ONE section while `finalScaffoldReason` requires all seven from
// COOPERATION.md Phase 4, and it never mentioned the `idea:` frontmatter the slug check reads. So
// the driver instructed its own drafter to produce an artifact its own gate rejects (review round
// 1, @codex-1 MAJOR) — the same defect class this audit is about, committed while fixing it.
//
// The section list is generated from protocol.RequiredFinalSections rather than retyped, so the
// prompt cannot drift from the gate.
func buildFinalDraftPrompt(ideaDir, path string) string {
	var sections strings.Builder
	for _, section := range protocol.RequiredFinalSections {
		sections.WriteString(section)
		sections.WriteString("\n")
	}
	slug := filepath.Base(ideaDir)
	return fmt.Sprintf(`You are the Parley Deck facilitator drafting FINAL.md.

Read %s/consensus.md (the accepted consensus + signoffs) and the round artifacts.
WRITE (create or overwrite) %s.

YAML frontmatter MUST include:
  idea: %s
  status: final

The body MUST contain ALL of these headings, in this order:

%s
"## Final plan / specification" needs at least three concrete lines describing the agreed design.
The other sections may be "N/A" when the idea is trivial or design-only, but the HEADING must be
present — a heading that is absent cannot be answered N/A deliberately.

No placeholders, no unexpanded <...> tokens. Be concrete. English only. Write the file now and
report only the path.`, ideaDir, path, slug, sections.String())
}
