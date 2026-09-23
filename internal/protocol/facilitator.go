package protocol

import (
	"path/filepath"
	"strings"
)

// Facilitator fields of an idea's 00-prompt.md frontmatter (lean-organizer A).
//
// `facilitator: <agent-id>` is OPTIONAL. A deck that does not declare it keeps
// byte-identical v1.48.0 behavior. A deck that does declare it runs a "declared
// facilitator run": the default is the pure organizer — the declared facilitator does
// not implement and does not verify code; participants own drafting, implementation,
// tests, and code verification.
//
// `facilitator_participates: true` opts the facilitator back into participation (full
// protocol context, no audience narrowing, role eligibility). Without it, a
// `facilitator:` that also names a `participants:` member is a preflight failure —
// fail-closed, both fields named — never a silent downgrade.
const (
	FacilitatorKey             = "facilitator"
	FacilitatorParticipatesKey = "facilitator_participates"
)

// FacilitatorRole is the parsed facilitator declaration of one idea.
type FacilitatorRole struct {
	Declared     bool
	Facilitator  string
	Participates bool
}

// ReadFacilitatorRole reads the facilitator declaration from an idea's 00-prompt.md.
// A missing file or missing field is an undeclared role (zero value), never an error:
// the absent-field deck must stay untouched.
func ReadFacilitatorRole(ideaDir string) FacilitatorRole {
	meta, err := ReadFrontmatter(filepath.Join(ideaDir, "00-prompt.md"))
	if err != nil {
		return FacilitatorRole{}
	}
	return FacilitatorRoleFromMeta(meta)
}

// FacilitatorRoleFromMeta maps parsed 00-prompt frontmatter onto the role.
func FacilitatorRoleFromMeta(meta map[string]string) FacilitatorRole {
	id := strings.Trim(strings.TrimSpace(meta[FacilitatorKey]), `"'`)
	if id == "" {
		return FacilitatorRole{}
	}
	return FacilitatorRole{
		Declared:     true,
		Facilitator:  id,
		Participates: isTrueValue(meta[FacilitatorParticipatesKey]),
	}
}

// Conflict returns the fail-closed preflight message when the declared facilitator is
// also listed in participants without facilitator_participates: true, or "" when the
// declaration is consistent.
func (r FacilitatorRole) Conflict(participants []string) string {
	if !r.Declared || r.Participates {
		return ""
	}
	for _, p := range participants {
		if p == r.Facilitator {
			return "idea frontmatter declares `facilitator: " + r.Facilitator + "` but also lists it in `participants:`" +
				" without `facilitator_participates: true`" +
				" — either remove the agent from participants (pure organizer) or set facilitator_participates: true (full participation)"
		}
	}
	return ""
}

// IneligibleForRoles reports whether the given agent id is barred from the
// drafter/implementer/reviewer/goal-done roles in this idea: exactly the declared
// facilitator of a declared run who has not opted back into participation.
func (r FacilitatorRole) IneligibleForRoles(id string) bool {
	return r.Declared && !r.Participates && id == r.Facilitator
}

// isTrueValue matches the exact case-insensitive "true" convention used by
// strict_gate / auto_implement parsing.
func isTrueValue(raw string) bool {
	return strings.EqualFold(strings.TrimSpace(raw), "true")
}
