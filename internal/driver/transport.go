package driver

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EffectiveTransport returns the transport governing this idea: the idea-level
// `transport:` in 00-prompt.md if present, else the project COOPERATION.md global
// (consensus D8). Returns "" if neither is found. This is why a local-dir idea
// can run inside an otherwise github-pr project without auto-driving the project.
func EffectiveTransport(ideaDir, root string) string {
	if t, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "transport"); ok {
		if n := normalizeTransport(t); n != "" {
			return n
		}
	}
	// COOPERATION.md is two levels up from …/ideas/<slug>.
	for _, candidate := range []string{
		filepath.Join(ideaDir, "..", "..", "COOPERATION.md"),
		filepath.Join(root, "parley-deck", "COOPERATION.md"),
	} {
		data, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, "Transport:") {
				if n := normalizeTransport(backtickValue(line)); n != "" {
					return n
				}
			}
		}
	}
	return ""
}

// ReadCrossReviewRounds reads cross_review_rounds from 00-prompt.md; default 1.
// N=0 is an explicit straight-to-consensus bypass.
func ReadCrossReviewRounds(ideaDir string) int {
	const def = 1
	if v, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "cross_review_rounds"); ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n >= 0 {
			return n
		}
	}
	return def
}

func normalizeTransport(raw string) string {
	v := strings.ToLower(strings.Trim(strings.TrimSpace(raw), "`'\" "))
	switch v {
	case "local-dir", "github-pr", "gitlab-mr":
		return v
	}
	return ""
}

func backtickValue(line string) string {
	start := strings.Index(line, "`")
	if start < 0 {
		return ""
	}
	rest := line[start+1:]
	end := strings.Index(rest, "`")
	if end < 0 {
		return ""
	}
	return rest[:end]
}
