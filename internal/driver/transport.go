package driver

import (
	"path/filepath"
	"strconv"
	"strings"

	"parley-deck-cli/internal/protocol"
)

// EffectiveTransport returns the transport governing this idea: the idea-level
// `transport:` in 00-prompt.md if present, else the project COOPERATION.md global
// (consensus D8). Returns "" if neither is found. This is why a local-dir idea
// can run inside an otherwise github-pr project without auto-driving the project.
// The global is read via protocol.ReadWorkspaceStatus, whose parser tolerates the
// `**Transport:** local-dir` (no-backtick) Markdown variation.
func EffectiveTransport(ideaDir, root string) string {
	if t, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "transport"); ok {
		if n := normalizeTransport(t); n != "" {
			return n
		}
	}
	if status, err := protocol.ReadWorkspaceStatus(root); err == nil {
		if n := normalizeTransport(status.Transport); n != "" {
			return n
		}
	}
	return ""
}

// ReadCrossReviewRounds reads cross_review_rounds from 00-prompt.md; default 1.
// N=0 is an explicit straight-to-consensus bypass.
func ReadCrossReviewRounds(ideaDir string) int {
	const def = 1
	if v, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "cross_review_rounds"); ok {
		if n, err := strconv.Atoi(strings.Trim(strings.TrimSpace(v), `"'`)); err == nil && n >= 0 {
			return n
		}
	}
	return def
}

// ReadAutoImplement reads the idea-level auto_implement opt-in from 00-prompt.md;
// default false. Code-writing phases (Implement/Fixup) require this true (D3).
func ReadAutoImplement(ideaDir string) bool {
	if v, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "auto_implement"); ok {
		return strings.EqualFold(strings.Trim(strings.TrimSpace(v), `"'`), "true")
	}
	return false
}

// ReadStrictGate reads the idea-level strict_gate opt-in from 00-prompt.md; default
// false. When true, the review loop completes (LE-2) only after a fresh full-scope
// closing review round is certified clean, not merely on outstanding_agreed_fixes == 0.
func ReadStrictGate(ideaDir string) bool {
	if v, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "strict_gate"); ok {
		return strings.EqualFold(strings.Trim(strings.TrimSpace(v), `"'`), "true")
	}
	return false
}

// ReadRequireModelDiversity reads require_model_diversity from 00-prompt.md; default
// false. When true, the driver escalates (instead of warning) if every reviewer shares
// the implementer's model (LE-3).
func ReadRequireModelDiversity(ideaDir string) bool {
	if v, ok := readFrontmatterField(filepath.Join(ideaDir, "00-prompt.md"), "require_model_diversity"); ok {
		return strings.EqualFold(strings.Trim(strings.TrimSpace(v), `"'`), "true")
	}
	return false
}

func normalizeTransport(raw string) string {
	v := strings.ToLower(strings.Trim(strings.TrimSpace(raw), "`'\"* "))
	switch v {
	case "local-dir", "github-pr", "gitlab-mr":
		return v
	}
	return ""
}
