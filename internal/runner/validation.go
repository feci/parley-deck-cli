package runner

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"parley-deck-cli/internal/protocol"
)

// ValidateRoundArtifact validates a round-N artifact. Round 1 keeps the strict
// round-01 contract; later (cross-review) rounds require matching frontmatter
// and at least one section heading.
func ValidateRoundArtifact(path, agentID, ideaSlug string, round int) error {
	if round <= 1 {
		return ValidateRoundOneArtifact(path, agentID, ideaSlug)
	}
	meta, err := readRoundOneFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent": agentID,
		"idea":  ideaSlug,
		"round": strconv.Itoa(round),
	} {
		if got := strings.Trim(strings.TrimSpace(meta[key]), `"'`); got != want {
			return fmt.Errorf("%s frontmatter %s=%q, want %q", path, key, got, want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	body := string(data)
	if !strings.Contains(body, "## ") {
		return fmt.Errorf("%s has no section headings", path)
	}
	// Later rounds have no fixed section list, so the floor is that at least one section says
	// something — same defect as codex-1/F17, one level looser.
	if !anySectionHasContent(body) {
		return fmt.Errorf("%s has section headings but no content under any of them", path)
	}
	return nil
}

func anySectionHasContent(body string) bool {
	for _, part := range strings.Split(body, "\n## ")[1:] {
		if i := strings.Index(part, "\n"); i >= 0 {
			for _, line := range strings.Split(part[i:], "\n") {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" || strings.HasPrefix(trimmed, "<!--") {
					continue
				}
				return true
			}
		}
	}
	return false
}

func ValidateRoundOneArtifact(path, agentID, ideaSlug string) error {
	meta, err := readRoundOneFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent": agentID,
		"idea":  ideaSlug,
		"round": "1",
	} {
		if got := strings.Trim(strings.TrimSpace(meta[key]), `"'`); got != want {
			return fmt.Errorf("%s frontmatter %s=%q, want %q", path, key, got, want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	body := string(data)
	// §15.6(a): the acquisition duty is checked HERE, on the runtime path, not only in
	// protocol.ValidateRoundOneArtifact. A validator with no caller is the defect class this
	// rule exists to prevent, and shipping one inside its own fix is how it recurs.
	//
	// This section uses the non-empty check rather than strings.Contains: a substring mention or
	// a bare heading is a rubber-stamp, not an enumerated search.
	if !protocol.HasNonEmptySection(body, protocol.RoundOneRequiredSection) {
		return fmt.Errorf("%s is missing a non-empty %q section (§15.6a)", path, protocol.RoundOneRequiredSection)
	}
	for _, section := range []string{
		"## Summary",
		"## Proposed approach",
		"## Concerns / open questions",
		"## Risks",
	} {
		if !strings.Contains(body, section) {
			return fmt.Errorf("%s missing required section %q", path, section)
		}
		// A heading is not an answer. The validator used to check only that the four headings
		// existed, so an artifact with every required section EMPTY completed round 1 and the
		// auto-driver advanced on it (audit finding codex-1/F17). Measured before enforcing:
		// of 211 round-01 artifacts carrying all four sections, 0 have an empty one.
		if !sectionHasContent(body, section) {
			return fmt.Errorf("%s has an empty required section %q", path, section)
		}
	}
	return nil
}

// sectionHasContent reports whether the named heading is followed by anything but blank lines
// and HTML comments, up to the next `## ` heading.
func sectionHasContent(body, section string) bool {
	idx := strings.Index(body, section)
	if idx < 0 {
		return false
	}
	rest := body[idx+len(section):]
	if next := strings.Index(rest, "\n## "); next >= 0 {
		rest = rest[:next]
	}
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "<!--") {
			continue
		}
		return true
	}
	return false
}

func readRoundOneFrontmatter(path string) (map[string]string, error) {
	meta, err := protocol.ReadFrontmatter(path)
	if err != nil {
		return nil, err
	}
	if hasAnyFrontmatter(meta) {
		return meta, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseLeadingFrontmatterWithoutOpeningFence(string(data)), nil
}

func hasAnyFrontmatter(meta map[string]string) bool {
	for _, key := range []string{"agent", "idea", "round"} {
		if strings.TrimSpace(meta[key]) != "" {
			return true
		}
	}
	return false
}

func parseLeadingFrontmatterWithoutOpeningFence(data string) map[string]string {
	lines := strings.Split(data, "\n")
	meta := map[string]string{}
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if i == 0 || len(meta) == 0 {
				return map[string]string{}
			}
			return meta
		}
		if trimmed == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return map[string]string{}
		}
		key = strings.TrimSpace(key)
		if key == "" || strings.ContainsAny(key, " \t") {
			return map[string]string{}
		}
		meta[key] = strings.TrimSpace(value)
	}
	return map[string]string{}
}
