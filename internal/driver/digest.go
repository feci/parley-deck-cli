package driver

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// digest.go builds a consolidated round digest (tui-round-summary): a deterministic,
// LLM-free position map of a completed round, emitted as a `round.digest` event for the
// Home tab. Stance is reported as keyword FLAGS (hints), never verdicts.

// AgentLine is one participant's row in a round digest.
type AgentLine struct {
	Agent    string `json:"agent"`
	Position string `json:"position"`  // first Summary sentence/paragraph, capped
	Fell     bool   `json:"fell_back"` // true when ## Summary was absent (degraded extraction)
	Present  bool   `json:"present"`   // artifact existed and parsed
}

// RoundDigest is the position map for a completed round.
type RoundDigest struct {
	Idea        string      `json:"idea"`
	Round       int         `json:"round"`
	Total       int         `json:"total"`
	Completed   int         `json:"completed"`
	Lines       []AgentLine `json:"lines"`
	FlagBlock   int         `json:"flag_block"`
	FlagCounter int         `json:"flag_counter"`
	FlagAccept  int         `json:"flag_accept"`
	FlagEscal   int         `json:"flag_escalate"`
	Next        string      `json:"next"`
}

const digestPositionCap = 120

var (
	summaryHeadingRe = regexp.MustCompile(`(?m)^##\s+Summary\s*$`)
	frontmatterRe    = regexp.MustCompile(`(?s)^\s*---.*?---\s*`)
	sentenceEndRe    = regexp.MustCompile(`[.!?]\s`)
)

// BuildRoundDigest reads each participant's round-NN artifact and builds the digest.
// It is pure over the filesystem (read-only) and never errors: a missing/unreadable
// artifact yields a not-present line so a display feature can never block advancement.
// next is the driver's already-decided next action (e.g. "opening round 03").
func BuildRoundDigest(ideaDir, ideaSlug string, round int, participants []string, next string) RoundDigest {
	d := RoundDigest{Idea: ideaSlug, Round: round, Total: len(participants), Next: next}
	roundDir := filepath.Join(ideaDir, roundLabel(round))
	for _, p := range participants {
		line := AgentLine{Agent: p}
		body, err := os.ReadFile(filepath.Join(roundDir, p+".md"))
		if err != nil {
			d.Lines = append(d.Lines, line)
			continue
		}
		line.Present = true
		d.Completed++
		text := string(body)
		pos, fell := extractPosition(text)
		line.Position = pos
		line.Fell = fell
		d.Lines = append(d.Lines, line)
		b, c, a, e := stanceFlags(text)
		d.FlagBlock += b
		d.FlagCounter += c
		d.FlagAccept += a
		d.FlagEscal += e
	}
	return d
}

// extractPosition pulls the first sentence/paragraph under `## Summary`, capped. When
// no Summary heading exists it falls back to the first prose paragraph after the
// frontmatter and reports fell=true so the UI can tag the extraction as degraded.
func extractPosition(text string) (pos string, fell bool) {
	if loc := summaryHeadingRe.FindStringIndex(text); loc != nil {
		return firstSentence(text[loc[1]:]), false
	}
	stripped := frontmatterRe.ReplaceAllString(text, "")
	return firstSentence(stripped), true
}

func firstSentence(s string) string {
	// First non-empty paragraph.
	para := ""
	for _, block := range strings.Split(s, "\n\n") {
		block = strings.TrimSpace(block)
		// Skip markdown headings and horizontal rules used as separators.
		if block == "" || strings.HasPrefix(block, "#") || block == "---" {
			continue
		}
		para = block
		break
	}
	para = strings.Join(strings.Fields(para), " ") // flatten newlines/indent
	if para == "" {
		return ""
	}
	// Truncate at the first sentence boundary if that keeps it under the cap.
	if m := sentenceEndRe.FindStringIndex(para); m != nil && m[0] < digestPositionCap {
		return strings.TrimSpace(para[:m[0]+1])
	}
	if len(para) > digestPositionCap {
		return strings.TrimSpace(para[:digestPositionCap]) + "…"
	}
	return para
}

// stanceFlags counts curated keyword mentions. These are HINTS, never verdicts — the
// UI must present them as mentions and the user drills into the tab to judge.
func stanceFlags(text string) (block, counter, accept, escalate int) {
	low := strings.ToLower(text)
	block = strings.Count(low, "block") // also covers "blocker" (superstring)
	counter = strings.Count(low, "counter-proposal")
	accept = strings.Count(low, "accept") + strings.Count(low, "agree")
	escalate = strings.Count(low, "escalat")
	return
}
