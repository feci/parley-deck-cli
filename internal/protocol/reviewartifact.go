package protocol

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ValidateReviewArtifact checks a review/round-NN/<agent>.md.
//
// It lives here, below both `runner` and `consensus`, because a review artifact must satisfy
// the same contract no matter which entry point produced it: the driver validates what an agent
// just wrote, and review-consensus drafting validates what is already on disk. Two copies of a
// rule are two rules (codex-1/F18 follow-up).
func ValidateReviewArtifact(path, agentID, ideaSlug string, round int) error {
	meta, err := ReadFrontmatter(path)
	if err != nil {
		return err
	}
	for key, want := range map[string]string{
		"agent":        agentID,
		"idea":         ideaSlug,
		"review-round": strconv.Itoa(round),
	} {
		if got := strings.Trim(strings.TrimSpace(meta[key]), `"'`); got != want {
			return fmt.Errorf("%s frontmatter %s=%q, want %q", path, key, got, want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	if !hasHeadingLine(content, "## Findings") {
		return fmt.Errorf("%s missing '## Findings' heading", path)
	}
	// LE-1 (refutation-default, review fix F5): require a real heading AND non-empty
	// content — a substring mention or an empty section is a rubber-stamp, not work shown.
	if !hasNonEmptySection(content, "## Refutation attempts") {
		return fmt.Errorf("%s missing a non-empty '## Refutation attempts' section (refutation-default: a review must record what it tried to break)", path)
	}
	// A review is a statement about a specific tree. Without `reviewed-commit` nobody can tell
	// which one, so a stale review is indistinguishable from a current one and "the tree does not
	// move while a review round is open" becomes unverifiable after the fact (codex-1/F18).
	//
	// Measured before enforcing: 348 of 539 review artifacts already carry it. The 191 that do not
	// are historical and are not revalidated; this binds new reviews.
	if commit := strings.Trim(strings.TrimSpace(meta["reviewed-commit"]), `"'`); commit == "" {
		return fmt.Errorf("%s frontmatter has no reviewed-commit: a review must name the tree it reviewed", path)
	}
	return nil
}

// hasHeadingLine reports whether content has a line that is exactly the given heading.
func hasHeadingLine(content, heading string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == heading {
			return true
		}
	}
	return false
}

// hasNonEmptySection reports whether content has the given level-2 heading followed by at
// least one non-blank line before the next level-2 heading.
func hasNonEmptySection(content, heading string) bool {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != heading {
			continue
		}
		for _, next := range lines[i+1:] {
			n := strings.TrimSpace(next)
			if n == "" {
				continue
			}
			if strings.HasPrefix(n, "## ") {
				break
			}
			return true
		}
		return false
	}
	return false
}
