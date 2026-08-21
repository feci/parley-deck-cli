package protocol

import (
	"fmt"
	"strings"
)

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

// ValidateFinal is the ONE gate for a FINAL.md, used by manual finalization and by the driver.
//
// They used to differ: the driver checked the idea slug and `status: final` while manual
// `consensus finalize` checked content only, so the manual path closed an idea around an artifact
// declaring a DIFFERENT idea and a non-final status (review round 1, @codex-1 MAJOR). A final
// status that can authenticate the wrong artifact is worse than none.
//
// wantSlug is the idea the artifact must claim; pass "" to skip that check.
func ValidateFinal(body, wantSlug string) string {
	if status := frontmatterValue(body, "status"); status != "final" {
		return fmt.Sprintf("frontmatter status=%q, want final", status)
	}
	if wantSlug != "" {
		declared := frontmatterValue(body, "idea")
		if declared == "" {
			return "frontmatter has no idea slug"
		}
		if declared != wantSlug {
			return fmt.Sprintf("frontmatter idea=%q but it closes idea %q", declared, wantSlug)
		}
	}
	return FinalIsScaffold(body)
}

// frontmatterValue reads one scalar from a leading YAML frontmatter block.
func frontmatterValue(body, key string) string {
	if !strings.HasPrefix(body, "---") {
		return ""
	}
	rest := strings.TrimPrefix(body, "---")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	for _, line := range strings.Split(rest[:end], "\n") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), key+":"); ok {
			return strings.Trim(strings.TrimSpace(v), `"'`)
		}
	}
	return ""
}
