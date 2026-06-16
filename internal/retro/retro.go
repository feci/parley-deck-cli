// Package retro implements the read-only "diverse, hard cases" mining for
// Parley Deck's retrospective optimization (COOPERATION.md §13). It scans the
// deck's own structured artifacts, scores ideas by failure density, selects a
// type-diverse coreset, and produces a diagnosis. It NEVER mutates the harness or
// protocol; the `parley retro` command may at most scaffold a single new
// ideas/<slug>/00-prompt.md (see internal/app/retro.go).
package retro

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"parley-deck-cli/internal/protocol"
)

// IdeaSignals is the deterministic failure-density profile of one past idea,
// derived only from its canonical artifacts (§13.2 evidence corpus = structured
// artifacts; no raw session transcripts in v1).
type IdeaSignals struct {
	Slug         string   `json:"slug"`
	Rounds       int      `json:"rounds"`        // round-NN design rounds
	ReviewRounds int      `json:"review_rounds"` // review/round-NN cycles
	FixupCycles  int      `json:"fixup_cycles"`  // "## Fix-up cycle" in IMPLEMENTATION.md
	NotFixed     int      `json:"not_fixed"`     // NOT-FIXED occurrences in review files
	Dismissed    int      `json:"dismissed"`     // dismissed findings recorded
	Escalations  int      `json:"escalations"`   // inbox *to-user* notes for this idea
	Blocked      bool     `json:"blocked"`       // a ❌ BLOCKER signoff or BLOCK verdict
	Abandoned    bool     `json:"abandoned"`     // status: abandoned in IMPLEMENTATION/00-prompt
	RunFailures  int      `json:"run_failures"`  // agent.failed/watchdog/driver.error in runs/*/events.jsonl
	Score        float64  `json:"score"`         // weighted failure density
	FailureType  string   `json:"failure_type"`  // dominant failure-mode bucket (diversity key)
	Reasons      []string `json:"reasons"`       // human-readable contributors
}

var (
	reRoundDir   = regexp.MustCompile(`^round-\d+$`)
	reFixupCycle = regexp.MustCompile(`(?m)^## Fix-up cycle`)
	reNotFixed   = regexp.MustCompile(`NOT-FIXED`)
	// Blocker is detected ONLY from structural lines, never from prose that merely
	// discusses blockers (e.g. a consensus that says "reviewer files use
	// `Verdict: BLOCK`"). reBlockerSignoff matches a real ❌ signoff line
	// (Status: ❌ …, line-anchored); a BLOCK *verdict* is detected separately by
	// hasBlockVerdict (a `## Verdict` heading whose value line is BLOCK/❌).
	reBlockerSignoff = regexp.MustCompile(`(?m)^Status:.*❌`)
	reVerdictHeading = regexp.MustCompile(`(?m)^##+\s*Verdict\b`)
	// status frontmatter value (IMPLEMENTATION.md / 00-prompt.md).
	reStatus    = regexp.MustCompile(`(?m)^status:\s*(.+)$`)
	reRunFailed = regexp.MustCompile(`"type"\s*:\s*"(?:agent\.failed|agent\.no_first_output|agent\.stalled|driver\.error)"`)
	reRunIdea   = regexp.MustCompile(`"idea"\s*:\s*"([^"]+)"`)
)

// Scan walks parley-deck/ideas/* and builds one IdeaSignals per idea. It is
// strictly read-only and tolerant of missing/partial artifacts.
func Scan(root string) ([]IdeaSignals, error) {
	ideasDir := filepath.Join(root, protocol.DeckDir, "ideas")
	entries, err := os.ReadDir(ideasDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	inboxNotes := readInboxToUser(filepath.Join(root, protocol.DeckDir, "inbox"))
	runFailures := scanRuns(filepath.Join(root, protocol.DeckDir, "runs"))

	var out []IdeaSignals
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		s := scanIdea(filepath.Join(ideasDir, e.Name()), e.Name(), inboxNotes)
		s.RunFailures = runFailures[e.Name()]
		score(&s)
		out = append(out, s)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out, nil
}

func scanIdea(dir, slug string, inboxNotes []string) IdeaSignals {
	s := IdeaSignals{Slug: slug}
	if subs, err := os.ReadDir(dir); err == nil {
		for _, sub := range subs {
			if sub.IsDir() && reRoundDir.MatchString(sub.Name()) {
				s.Rounds++
			}
		}
	}
	if subs, err := os.ReadDir(filepath.Join(dir, "review")); err == nil {
		for _, sub := range subs {
			if sub.IsDir() && reRoundDir.MatchString(sub.Name()) {
				s.ReviewRounds++
			}
		}
	}
	impl := readFile(filepath.Join(dir, "IMPLEMENTATION.md"))
	s.FixupCycles = len(reFixupCycle.FindAllString(impl, -1))
	s.Abandoned = statusIs(impl, "abandoned") || statusIs(readFile(filepath.Join(dir, "00-prompt.md")), "abandoned")

	// Signals are anchored to their structural home so prose that merely discusses
	// them never false-positives (the round-01 review of this very tool tripped on
	// a consensus that quoted "Verdict: BLOCK").
	consensusText := readFile(filepath.Join(dir, "consensus.md")) + "\n" + readFile(filepath.Join(dir, "review", "consensus.md"))
	roundReviews := concatRoundReviews(filepath.Join(dir, "review"))
	// NOT-FIXED is a re-review fix-verification marker: count it only in review
	// round files, not in consensus/IMPLEMENTATION prose.
	s.NotFixed = len(reNotFixed.FindAllString(roundReviews, -1))
	// Blocker: a real ❌ signoff line (in consensus/impl/round files) or a BLOCK
	// verdict in a review round file — never a prose mention.
	s.Blocked = reBlockerSignoff.MatchString(impl+"\n"+consensusText+"\n"+roundReviews) || hasBlockVerdict(roundReviews)
	s.Dismissed = countDismissed(consensusText)

	for _, note := range inboxNotes {
		if strings.Contains(note, slug) {
			s.Escalations++
		}
	}
	return s
}

// score is the deterministic failure-density weighting and failure-type bucketing.
func score(s *IdeaSignals) {
	add := func(cond bool, w float64, reason string) {
		if cond {
			s.Score += w
			s.Reasons = append(s.Reasons, reason)
		}
	}
	// Extra design/review rounds beyond the minimum (1 design, 1 review) are friction.
	add(s.Rounds > 1, float64(s.Rounds-1)*1.0, "extra design rounds")
	add(s.ReviewRounds > 1, float64(s.ReviewRounds-1)*2.0, "multiple review cycles")
	add(s.FixupCycles > 0, float64(s.FixupCycles)*2.0, "fix-up churn")
	add(s.NotFixed > 0, float64(s.NotFixed)*1.5, "NOT-FIXED re-reviews")
	add(s.Dismissed > 0, float64(s.Dismissed)*0.5, "dismissed findings")
	add(s.Escalations > 0, float64(s.Escalations)*3.0, "user escalations")
	add(s.RunFailures > 0, float64(s.RunFailures)*1.5, "run failures (agent.failed/watchdog/driver.error)")
	add(s.Abandoned, 4.0, "abandoned work")
	add(s.Blocked, 4.0, "blocker signoff")

	s.FailureType = classify(s)
}

// classify buckets an idea by its dominant failure-mode signal (the diversity key
// for coreset selection). Order reflects severity/specificity.
func classify(s *IdeaSignals) string {
	switch {
	case s.Blocked || s.Abandoned:
		return "blocked-or-abandoned"
	case s.Escalations > 0:
		return "escalation"
	case s.RunFailures > 0:
		return "runtime-failure"
	case s.FixupCycles >= 2 || s.NotFixed > 0:
		return "fix-up-heavy"
	case s.ReviewRounds > 1:
		return "review-churn"
	case s.Rounds > 1:
		return "design-churn"
	default:
		return "low-friction"
	}
}

// Select returns a type-diverse coreset of the highest-scoring ideas: greedily
// take the top scorer of each not-yet-covered failure type first, then fill the
// remaining slots by score. Deterministic (input is score-sorted).
func Select(signals []IdeaSignals, k int) []IdeaSignals {
	if k <= 0 || len(signals) == 0 {
		return nil
	}
	var picked []IdeaSignals
	seenType := map[string]bool{}
	// Pass 1: one representative per failure type (skip the residual low-friction
	// bucket so the coreset is genuinely "hard").
	for _, s := range signals {
		if len(picked) >= k {
			break
		}
		if s.FailureType == "low-friction" || s.Score <= 0 {
			continue
		}
		if !seenType[s.FailureType] {
			seenType[s.FailureType] = true
			picked = append(picked, s)
		}
	}
	// Pass 2: fill remaining slots by score, no duplicates.
	inCoreset := map[string]bool{}
	for _, p := range picked {
		inCoreset[p.Slug] = true
	}
	for _, s := range signals {
		if len(picked) >= k {
			break
		}
		if inCoreset[s.Slug] || s.FailureType == "low-friction" || s.Score <= 0 {
			continue
		}
		picked = append(picked, s)
		inCoreset[s.Slug] = true
	}
	return picked
}

// Diagnose produces a read-only, deterministic diagnosis grouped by failure type.
func Diagnose(coreset []IdeaSignals) string {
	if len(coreset) == 0 {
		return "No hard cases found: no past idea shows review churn, fix-up churn, escalations, or blockers.\n"
	}
	byType := map[string][]IdeaSignals{}
	var order []string
	for _, s := range coreset {
		if _, ok := byType[s.FailureType]; !ok {
			order = append(order, s.FailureType)
		}
		byType[s.FailureType] = append(byType[s.FailureType], s)
	}
	var b strings.Builder
	b.WriteString("Retrospective diagnosis (deterministic, from structured artifacts):\n\n")
	for _, t := range order {
		b.WriteString("## Failure mode: " + t + "\n")
		for _, s := range byType[t] {
			b.WriteString("- " + s.Slug + " — " + strings.Join(s.Reasons, ", "))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("These are hypotheses, not findings (§13): a retro pass proposes; the normal idea gate accepts.\n")
	return b.String()
}

// --- helpers -------------------------------------------------------------------

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// concatRoundReviews concatenates only the review ROUND files
// (review/round-NN/*.md), excluding review/consensus.md, so prose in the review
// consensus never feeds the per-finding signals (NOT-FIXED, BLOCK verdicts).
func concatRoundReviews(reviewDir string) string {
	rounds, err := os.ReadDir(reviewDir)
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, rd := range rounds {
		if !rd.IsDir() || !reRoundDir.MatchString(rd.Name()) {
			continue
		}
		files, _ := os.ReadDir(filepath.Join(reviewDir, rd.Name()))
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
				b.WriteString(readFile(filepath.Join(reviewDir, rd.Name(), f.Name())))
				b.WriteString("\n")
			}
		}
	}
	return b.String()
}

// hasBlockVerdict reports whether any "## Verdict" heading in text has a BLOCK (or
// ❌) value — either inline ("## Verdict: BLOCK") or on the next non-blank line
// ("## Verdict\nBLOCK"). A prose mention like "use `Verdict: BLOCK`" is not a
// heading and never matches.
func hasBlockVerdict(text string) bool {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if !reVerdictHeading.MatchString(line) {
			continue
		}
		val := strings.TrimLeft(strings.TrimSpace(reVerdictHeading.ReplaceAllString(line, "")), ": \t")
		if val == "" {
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) != "" {
					val = strings.TrimSpace(lines[j])
					break
				}
			}
		}
		// Match on the verdict's LEADING token (BLOCK/BLOCKER) or an explicit ❌,
		// not the substring "block" — so "ACCEPT … no blocking issues" and
		// "REQUEST-CHANGES: block on X" (a fixes-requested verdict, captured by the
		// fix-up signals) do not count as a hard BLOCK.
		switch strings.ToUpper(firstToken(val)) {
		case "BLOCK", "BLOCKER":
			return true
		}
		if strings.Contains(val, "❌") {
			return true
		}
	}
	return false
}

// firstToken returns the leading word of s, up to the first space, tab, colon,
// period, or comma.
func firstToken(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, " \t:.,"); i >= 0 {
		return s[:i]
	}
	return s
}

func readInboxToUser(inboxDir string) []string {
	var notes []string
	entries, err := os.ReadDir(inboxDir)
	if err != nil {
		return notes
	}
	for _, e := range entries {
		if e.IsDir() || !strings.Contains(e.Name(), "to-user") {
			continue
		}
		notes = append(notes, e.Name()+"\n"+readFile(filepath.Join(inboxDir, e.Name())))
	}
	return notes
}

// statusIs reports whether the document's `status:` frontmatter value equals want
// (case-insensitive). Used to detect abandoned work (D4).
func statusIs(text, want string) bool {
	m := reStatus.FindStringSubmatch(text)
	if m == nil {
		return false
	}
	return strings.EqualFold(strings.Trim(strings.TrimSpace(m[1]), `"'`), want)
}

// scanRuns reads structured run event logs (parley-deck/runs/*/events.jsonl —
// NOT raw session transcripts) and returns, per idea slug, the count of
// failure-class events (agent.failed / watchdog / driver.error). The idea slug is
// taken from the run's run.created event (the first "idea" field in the log).
func scanRuns(runsDir string) map[string]int {
	out := map[string]int{}
	entries, err := os.ReadDir(runsDir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data := readFile(filepath.Join(runsDir, e.Name(), "events.jsonl"))
		if data == "" {
			continue
		}
		m := reRunIdea.FindStringSubmatch(data)
		if m == nil {
			continue
		}
		if n := len(reRunFailed.FindAllString(data, -1)); n > 0 {
			out[m[1]] += n
		}
	}
	return out
}

// countDismissed counts entries under a "## Dismissed findings" heading across the
// scanned text (bounded to that section per document occurrence).
func countDismissed(text string) int {
	n := 0
	lines := strings.Split(text, "\n")
	in := false
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "## Dismissed findings"):
			in = true
		case in && strings.HasPrefix(line, "## "):
			in = false
		case in && (strings.HasPrefix(strings.TrimSpace(line), "1.") || strings.HasPrefix(strings.TrimSpace(line), "- ")):
			if !strings.Contains(line, "None") {
				n++
			}
		}
	}
	return n
}
