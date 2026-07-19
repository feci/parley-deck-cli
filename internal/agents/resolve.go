package agents

import (
	"fmt"
	"strings"
)

// pathSafeParticipant reports whether a participant id is safe to use as a single
// filepath segment (artifact paths, run-log dirs) — closing a path-traversal hole
// where a malicious deck could put `[roster."../../x"]` / a `..`-bearing participant
// into an idea and make a CLI write outside the deck (review CRITICAL, codex-1). It
// is deliberately a CONTAINMENT check, not the strict §2 grammar, so legacy custom
// spec ids with `_`/uppercase still resolve (review MINOR, kimi-1); the §2 reader
// (protocol/roster.go) enforces the stricter roster-id grammar where it matters.
func pathSafeParticipant(id string) bool {
	if id == "" || id == "." || id == ".." {
		return false
	}
	if strings.ContainsAny(id, `/\`) || strings.Contains(id, "..") {
		return false
	}
	if strings.HasPrefix(id, ".") || strings.HasSuffix(id, ".") {
		return false
	}
	return true
}

// ResolveParticipant maps a participant / roster ID (e.g. "claude-1") to a
// discovered agent, fail-closed. This closes the two-namespace schism: the driver
// used to match `participants:` against spec/family IDs by exact string equality,
// so a deck whose roster is `[claude-1, codex-1, …]` selected zero agents. Now:
//
//  1. exact spec-ID match — legacy decks whose participants ARE spec/family ids;
//  2. explicit mapping[participant] -> family, matched against a discovered spec ID
//     (written by `parley roster init` as `[roster.<id>] adapter = "<family>"`);
//  3. otherwise a hard error — never a prefix heuristic (antigravity-1 -> agy would
//     break it), never a silent guess.
//
// The returned Discovery carries the participant string as its identity (Spec.ID)
// and the family as its AdapterID, so downstream artifact paths, frontmatter, and
// signoffs use the roster ID while launch/vendor dispatch uses the family via
// Spec.Adapter().
func ResolveParticipant(participant string, discovered []Discovery, mapping map[string]string) (Discovery, error) {
	if !pathSafeParticipant(participant) {
		return Discovery{}, fmt.Errorf("agents: unsafe participant id %q (no `/`, `\\`, `..`, or leading/trailing dot)", participant)
	}
	// (1) exact spec-ID match. Preserve an already-explicit adapter (review MINOR).
	for _, d := range discovered {
		if d.Found && d.ID == participant {
			d.Spec.AdapterID = d.Spec.Adapter()
			d.Spec.ID = participant
			return d, nil
		}
	}
	// (2) explicit roster -> family mapping.
	if family, ok := mapping[participant]; ok {
		for _, d := range discovered {
			if d.Found && d.ID == family {
				d.Spec.AdapterID = family
				d.Spec.ID = participant
				return d, nil
			}
		}
		return Discovery{}, fmt.Errorf("agents: participant %q maps to family %q, which is not an installed/discovered agent", participant, family)
	}
	// (3) fail closed — no exact id, no mapping.
	return Discovery{}, fmt.Errorf("agents: cannot resolve participant %q (no matching agent id and no [roster.*] mapping); run `parley roster init`", participant)
}
