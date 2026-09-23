package protocol

import "strings"

// RequiredConsensusSections is the Phase-3 consensus.md template from COOPERATION.md,
// including the §15 drafter duties that bind on every track (§15.7): §15.5's
// `## Drafter position changes` and §15.6's `## Alternatives disposition`.
//
// It lives in protocol/ for the same reason RequiredFinalSections does: the
// auto-driver's consensus drafting prompt, the consensus scaffold generator, and the
// consensus status gate must all read ONE list. Before this constant existed the
// driver's prompt emitted review-cycle headings (`## Trade-offs accepted`, `## Deferred
// follow-ups`, `## Dismissed findings`) that the Phase-3 template never names, and
// carried none of the §15.3/§15.5/§15.6 duty headings — so the driver instructed its
// own drafter to produce an artifact the protocol's own template rejects
// (meta-protocol-change-lean-organizer FINAL A.4, verified at HEAD ffa4587).
var RequiredConsensusSections = []string{
	"## Agreed decisions",
	"## Agreed trade-offs",
	"## Open items deferred to implementation",
	"## Comparison & blind spots",
	"## Drafter position changes",
	"## Alternatives disposition",
	"## Signoffs",
}

// ConditionalConsensusSections are duty sections whose presence depends on the
// deliberation itself: `## Verdict conflicts` exists ONLY when contradictory verdicts
// were issued (§15.3 — "Absent any conflict the section does not exist"). The drafting
// prompt must name them and state the condition; the gate does not hard-require them.
var ConditionalConsensusSections = []string{
	"## Verdict conflicts",
}

// MissingConsensusSections returns the required Phase-3 headings absent from a
// consensus.md body, without the "## " prefix, in template order. It reads the same
// RequiredConsensusSections slice the drafting prompt is generated from, so the gate
// cannot drift from the prompt (the A.4 parity property).
func MissingConsensusSections(body string) []string {
	var missing []string
	for _, section := range RequiredConsensusSections {
		if !strings.Contains(body, section) {
			missing = append(missing, strings.TrimPrefix(section, "## "))
		}
	}
	return missing
}
