package protocol

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Designated-implementer fields (idea meta-protocol-change-designated-implementer).
//
// `implementer:` in 00-prompt.md is a DESIGNATION — an owner instruction about who
// should execute FINAL. `implementer:` in IMPLEMENTATION.md is an OUTCOME RECORD — who
// actually did the work (the re-entry pin). The same key name carries the asymmetry on
// purpose; the artifact it sits in disambiguates it.
const (
	ImplementerKey = "implementer"
	// ImplementerWaivedKey records a confirmed waiver of a per-idea designation whose
	// designee is unavailable: `implementer_waived: <agent-id> — <reason> — confirmed
	// <date>` (the §9.0 exclusion record shape). It clears the tier-2 availability
	// gate and falls through to the default chain for that run.
	ImplementerWaivedKey = "implementer_waived"
	// ImplementerReassignedKey records a confirmed reassignment over a stale pin:
	// `implementer_reassigned: <old-agent-id> to <new-agent-id> — <reason> — confirmed
	// <date>`. Only with that recorded confirmation does a disagreeing designation
	// supersede the IMPLEMENTATION.md pin (no self-appointment).
	ImplementerReassignedKey = "implementer_reassigned"
)

// DesignationState is the four-state parse of the per-idea `implementer:` field. Key
// presence survives an empty value because ReadFrontmatter splits on the first colon
// and stores TrimSpace(value): a bare `implementer:` line is present-with-empty, which
// is a DIFFERENT state from absent and from `none`.
type DesignationState int

const (
	// DesignationAbsent: no `implementer:` key — tier 3 (global default) may fire.
	DesignationAbsent DesignationState = iota
	// DesignationEmpty: key present, value empty/whitespace — an INCOMPLETE
	// designation (a typo, not an intention); a blocking gate, never silent.
	DesignationEmpty
	// DesignationNone: value `none` (trimmed, any casing) — an explicit per-idea
	// non-designation: tier 3 suppressed, the default chain runs, no gate.
	DesignationNone
	// DesignationSet: value is an agent id — a designation, validated fail-closed.
	DesignationSet
)

// ImplementerDesignation is the parsed per-idea designation.
type ImplementerDesignation struct {
	State DesignationState
	// ID is the designated agent id, trimmed and quote-stripped, only when State ==
	// DesignationSet. A malformed value (e.g. `kimi-1  # from the default`) is kept
	// LITERAL here and fails closed downstream — it is never repaired into an id.
	ID string
}

// ImplementerDesignationFromMeta parses the four states from one frontmatter map.
func ImplementerDesignationFromMeta(meta map[string]string) ImplementerDesignation {
	raw, present := meta[ImplementerKey]
	if !present {
		return ImplementerDesignation{State: DesignationAbsent}
	}
	id := strings.Trim(strings.TrimSpace(raw), `"'`)
	if id == "" {
		return ImplementerDesignation{State: DesignationEmpty}
	}
	if strings.EqualFold(id, "none") {
		return ImplementerDesignation{State: DesignationNone}
	}
	return ImplementerDesignation{State: DesignationSet, ID: id}
}

// ReadImplementerDesignation reads the designation from an idea's 00-prompt.md. A
// missing or unreadable file is DesignationAbsent (the absent-field deck is
// untouched), matching ReadFacilitatorRole's posture.
func ReadImplementerDesignation(ideaDir string) ImplementerDesignation {
	meta, err := ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return ImplementerDesignation{State: DesignationAbsent}
	}
	return ImplementerDesignationFromMeta(meta)
}

// ImplementerSource identifies which rank of the resolution chain produced the
// implementer. The distinctions are protocol-fixed (so no consumer can confuse them);
// the constant names are the implementer's choice.
type ImplementerSource string

const (
	// SourceImplementerPin: rank 1 — the IMPLEMENTATION.md outcome record.
	SourceImplementerPin ImplementerSource = "pin"
	// SourceImplementerDesignation: rank 2 — the per-idea `implementer:` designation.
	SourceImplementerDesignation ImplementerSource = "designation"
	// SourceImplementerGlobalDefault: rank 3 — the layered `[defaults].default_implementer`.
	SourceImplementerGlobalDefault ImplementerSource = "global-default"
	// SourceImplementerNone: an explicit opt-out (`none`) at either designation layer;
	// resolution fell through to the default chain by the owner's deliberate choice.
	SourceImplementerNone ImplementerSource = "none"
	// SourceImplementerFallThroughUnavailable: a designation was present but its
	// designee was unavailable (or waived) and the run fell through to the default
	// chain. Distinct from an opt-out and from an unset deck.
	SourceImplementerFallThroughUnavailable ImplementerSource = "fall-through-unavailable"
	// SourceImplementerFallThroughInapplicable: a tier-3 designation was present but
	// its id is not eligible on this idea (non-participant, or the idea's own declared
	// non-participating facilitator) and the run fell through to the default chain.
	SourceImplementerFallThroughInapplicable ImplementerSource = "fall-through-inapplicable"
	// SourceImplementerLegacy: rank 4 — today's chain (FINAL.md metadata → positional
	// tail). Never emitted as an event source: on a deck with no designation nothing
	// new is emitted at all.
	SourceImplementerLegacy ImplementerSource = "legacy"
)

// ImplementerCandidate is one ordered source of implementer identity: either a
// frontmatter read (File + Keys, relative to the idea directory) or a pre-resolved
// explicit id (ID). Candidates are tried in order; the first whose id is non-empty
// (after the shared space/quote trim) and a member of the eligibility list wins.
type ImplementerCandidate struct {
	Source ImplementerSource
	File   string
	Keys   []string
	ID     string
}

// PinImplementerCandidates is rank 1 alone: the IMPLEMENTATION.md re-entry pin.
func PinImplementerCandidates() []ImplementerCandidate {
	return []ImplementerCandidate{
		{Source: SourceImplementerPin, File: "IMPLEMENTATION.md", Keys: []string{"implementer"}},
	}
}

// LegacyImplementerCandidates is today's chain, verbatim (rank 1 + the rank-4
// metadata read): IMPLEMENTATION.md `implementer`, then FINAL.md `implementer` /
// `drafted-by`. The positional/`""` tails stay at the two call sites, byte-identical.
func LegacyImplementerCandidates() []ImplementerCandidate {
	return []ImplementerCandidate{
		{Source: SourceImplementerPin, File: "IMPLEMENTATION.md", Keys: []string{"implementer"}},
		{Source: SourceImplementerLegacy, File: "FINAL.md", Keys: []string{"implementer", "drafted-by"}},
	}
}

// ResolveImplementerChain walks the ordered source list and returns the first
// candidate id that is a member of eligible. The eligibility list is an explicit
// parameter BECAUSE the two call sites legitimately differ (the driver passes the
// facilitator-filtered list; the consensus side passes the raw participants) — the
// divergence is visible and testable here rather than baked in.
func ResolveImplementerChain(ideaDir string, candidates []ImplementerCandidate, eligible []string) (string, ImplementerSource, bool) {
	isEligible := func(id string) bool {
		for _, p := range eligible {
			if p == id {
				return true
			}
		}
		return false
	}
	for _, c := range candidates {
		if c.ID != "" {
			if id := strings.Trim(strings.TrimSpace(c.ID), `"'`); id != "" && isEligible(id) {
				return id, c.Source, true
			}
			continue
		}
		meta, err := ReadFrontmatter(filepath.Join(ideaDir, c.File))
		if err != nil {
			continue
		}
		for _, k := range c.Keys {
			if id := strings.Trim(strings.TrimSpace(meta[k]), `"'`); id != "" && isEligible(id) {
				return id, c.Source, true
			}
		}
	}
	return "", "", false
}

// confirmationMarker is the §9.0 confirmation-record tail: exactly `confirmed
// <date>` (any casing, one space, ISO date shape). The date must additionally parse
// as a valid calendar date (review VC-D) — `confirmed 2026-02-31` is the "typo
// indistinguishable from intention" class R3 names.
var confirmationMarker = regexp.MustCompile(`(?i)^confirmed \d{4}-\d{2}-\d{2}$`)

// negationMarker matches the negations that must never appear in a segment before
// the confirmation marker: `not confirmed`, `not yet confirmed`, `unconfirmed`
// (any casing). A negated record is not a confirmation, whatever else it carries.
var negationMarker = regexp.MustCompile(`(?i)\b(?:unconfirmed|not(?:\s+yet)?\s+confirmed)\b`)

// confirmationRecordTail enforces the §9.0 confirmation-record shape shared by the
// waiver and reassignment readers (review AF-1). parts is the record value split on
// "—": parts[0] is the subject (the caller's concern), parts[1..n-2] are reason
// segments, parts[n-1] is the confirmation marker. The record is valid only with at
// least one non-empty reason segment, a strict trailing marker whose date parses as
// a valid calendar date, and no negation in any earlier segment.
func confirmationRecordTail(parts []string) bool {
	if len(parts) < 3 { // subject + reason + marker
		return false
	}
	for _, seg := range parts[:len(parts)-1] {
		if negationMarker.MatchString(seg) {
			return false
		}
	}
	for _, seg := range parts[1 : len(parts)-1] {
		if strings.TrimSpace(seg) == "" {
			return false // a reason segment must carry a reason
		}
	}
	marker := strings.TrimSpace(parts[len(parts)-1])
	if !confirmationMarker.MatchString(marker) {
		return false
	}
	date := strings.TrimSpace(marker[len("confirmed"):])
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// ImplementerWaived reports whether 00-prompt.md carries a confirmed waiver for id:
// `implementer_waived: <agent-id> — <reason> — confirmed <date>`. The subject
// segment must name the id EXACTLY (substring containment would let `kimi-10` clear
// kimi-1's gate), the tail must satisfy confirmationRecordTail; anything weaker is
// not a waiver.
func ImplementerWaived(meta map[string]string, id string) bool {
	v := strings.Trim(strings.TrimSpace(meta[ImplementerWaivedKey]), `"'`)
	if v == "" || id == "" {
		return false
	}
	parts := strings.Split(v, "—")
	if !confirmationRecordTail(parts) {
		return false
	}
	return strings.TrimSpace(parts[0]) == id
}

// ParseImplementerReassignment parses `implementer_reassigned: <old> to <new> —
// <reason> — confirmed <date>`. It returns ok only when the record names exactly the
// (old → new) pair under dispute and its tail satisfies confirmationRecordTail — the
// owner/author-recorded act that retires the abandoned attempt, never the incoming
// implementer's own edit.
func ParseImplementerReassignment(meta map[string]string, oldID, newID string) bool {
	v := strings.Trim(strings.TrimSpace(meta[ImplementerReassignedKey]), `"'`)
	if v == "" || oldID == "" || newID == "" {
		return false
	}
	parts := strings.Split(v, "—")
	if !confirmationRecordTail(parts) {
		return false
	}
	pair := strings.TrimSpace(parts[0])
	for _, sep := range []string{" to ", " → ", "->"} {
		if before, after, ok := strings.Cut(pair, sep); ok {
			return strings.TrimSpace(before) == oldID && strings.TrimSpace(after) == newID
		}
	}
	return false
}
