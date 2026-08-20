package protocol

import "strings"

// RequiredFinalSections is the FINAL.md template from COOPERATION.md Phase 4.
//
// The protocol says the CONTENT of the sections below the specification may be `N/A` for trivial
// or design-only ideas. It does not say the headings may be absent — FINAL.md is the single source
// of truth an implementer works from, and a heading that is not there cannot be answered `N/A`
// deliberately.
//
// It lives in protocol/ because two independent gates need the same list: the auto-driver's
// consensus gate and the manual `consensus finalize` verb. They disagreed before, which is how an
// idea could be closed around a scaffold by one path and rejected by the other.
var RequiredFinalSections = []string{
	"## Final plan / specification",
	"## Purpose / user-visible outcome",
	"## Context & orientation",
	"## Observable acceptance criteria",
	"## Idempotence & recovery",
	"## Known risks / de-risking",
	"## References",
}

// MissingFinalSections returns the required headings absent from a FINAL.md body, without the
// "## " prefix, in template order.
func MissingFinalSections(body string) []string {
	var missing []string
	for _, section := range RequiredFinalSections {
		if !strings.Contains(body, section) {
			missing = append(missing, strings.TrimPrefix(section, "## "))
		}
	}
	return missing
}

// FinalScaffoldPlaceholders are the unexpanded template tokens that mark a FINAL.md as unwritten.
var FinalScaffoldPlaceholders = []string{
	"<...>", "<…>", "<slug>", "<agent-id>", "<agent>", "<fill", "<todo", "<tbd", "<your ", "<replace", "<insert",
}

// FinalIsScaffold reports why a FINAL.md body is still a scaffold, or "" when it is written.
// It is deliberately about CONTENT only; callers own frontmatter and slug checks.
func FinalIsScaffold(body string) string {
	if missing := MissingFinalSections(body); len(missing) > 0 {
		return "missing required section(s): " + strings.Join(missing, ", ")
	}
	lower := strings.ToLower(body)
	for _, ph := range FinalScaffoldPlaceholders {
		if strings.Contains(lower, ph) {
			return "contains an unexpanded placeholder " + ph
		}
	}
	if content := specContentLines(body); content < 3 {
		return "'## Final plan / specification' has fewer than 3 content lines"
	}
	return ""
}

func specContentLines(body string) int {
	const section = "## Final plan / specification"
	idx := strings.Index(body, section)
	if idx < 0 {
		return 0
	}
	rest := body[idx+len(section):]
	if next := strings.Index(rest, "\n## "); next >= 0 {
		rest = rest[:next]
	}
	n := 0
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "<!--") {
			continue
		}
		n++
	}
	return n
}
