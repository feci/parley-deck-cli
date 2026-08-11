package runner

// Frontier context selection: send the previous round in full plus a carry-forward ledger, instead
// of re-sending every artifact from every prior round.
//
// From FINAL.md of idea `protocol-read-cost-regression`. Two runtime paths embedded history
// quadratically — `gatherPriorRounds` (design) and `gatherReviewContext` (review) — and both were
// STRICTER than the protocol they implement. Phase 2 requires addressing every other active agent
// and a counter-proposal for disagreement; it never requires re-reading every historical artifact.
//
// THE BOUNDARY, adopted verbatim from codex-1's signoff and binding on this file: this is an
// implementation-scoped CONTEXT OPTIMIZATION. It is not an artifact-validity rule and not a
// consensus rule. Nothing here may decide whether an artifact counts or whether consensus closes.
// If it ever does, the change needs a §7 protocol idea.
//
// WHY THE FALLBACK IS NOT BELT-AND-BRACES. Phase 2 rule 1 reads "Silence = implicit agreement", so
// the protocol converts an omission into consent. An objection dropped from the carry-forward is not
// a lost datum — it is agreement that was never given. Every uncertainty therefore selects full
// history, visibly.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// frontierContext assembles what a cross-review round receives.
//
// FIX-UP CYCLE 1. The first implementation derived a ledger by matching marker substrings and
// compacted whenever it matched anything. All three reviewers independently found the same
// CRITICAL: that is FAIL-OPEN. An objection whose wording contains no marker simply disappears,
// and because Phase 2 rule 1 is "Silence = implicit agreement", a disappeared objection is not a
// lost datum — it is recorded consent that was never given. Extending the marker list is
// whack-a-mole against unbounded natural language; the class of failure cannot be patched out.
//
// So the derived ledger is gone. Compaction now requires the AUTHORED ledger that codex-1's and
// hermes-1's signoffs specified, and its absence is one of the fallback conditions they required:
// "missing, invalid, ambiguous or challenged ledger state falls back to full history". Until decks
// carry authored ledgers, every round therefore gets full history and this file changes nothing —
// which is the correct behaviour for an optimization that cannot yet prove it is safe.
func frontierContext(dirFor func(int) string, round int, full func() (string, error)) (string, error) {
	// FIX-UP CYCLE 2 — compaction is HARD-DISABLED.
	//
	// Cycle 1 gated compaction on an authored ledger EXISTING. codex-1 (CRITICAL) and kimi-1
	// (MAJOR) both found that this is not fail-closed: any non-empty bytes at that pathname are
	// accepted as a participant-authored ledger with no parsing, no provenance, no
	// expected-participant coverage, no ownership, no lifecycle, no locator or hash check, and no
	// verdict-conflict detection. Creating one file would switch the optimization on with every
	// content protection still unbuilt — fail-open with an extra step.
	//
	// So the switch is a constant, not a file. The plumbing, the fallback and the tests ship and are
	// exercised; the behaviour shipped is byte-identical to the behaviour before this idea. It may
	// be flipped only when the signed ledger contract has a validator, including G3, G5 and G6.
	if !compactionEnabled {
		return full()
	}
	if round <= 2 {
		return full()
	}
	led, why := authoredLedger(dirFor, round-2)
	if led == "" {
		return fallbackTo(full, why)
	}
	prev, err := renderRound(dirFor(round-1), filepath.Base(dirFor(round-1)))
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(prev) == "" {
		return fallbackTo(full, "the previous round produced no readable artifacts")
	}
	return led + "\n" + prev, nil
}

// compactionEnabled is false until the participant-authored ledger has a VALIDATOR.
//
// Not a flag, not an env var, not a config key: a constant, so that nothing an agent or a deck can
// write turns it on. Turning it on is a source change that goes through review.
const compactionEnabled = false

// ledgerFileName is the participant-authored carry-forward ledger for a round.
const ledgerFileName = "_ledger.md"

// authoredLedger returns the concatenated authored ledgers for rounds 1..upTo, or "" plus the
// reason no compaction is permitted.
//
// It deliberately makes NO judgement about content. Deciding that a round's dissent is adequately
// represented is exactly the judgement the reviewers said an extractor must not make.
func authoredLedger(dirFor func(int) string, upTo int) (string, string) {
	var b strings.Builder
	found := 0
	for r := 1; r <= upTo; r++ {
		data, err := os.ReadFile(filepath.Join(dirFor(r), ledgerFileName))
		if err != nil {
			return "", fmt.Sprintf("no authored %s for round %d", ledgerFileName, r)
		}
		// A BOM is not content: TrimSpace alone lets a one-character file pass as an authored ledger.
		if strings.TrimSpace(strings.TrimPrefix(string(data), "\ufeff")) == "" {
			return "", fmt.Sprintf("the authored %s for round %d is empty", ledgerFileName, r)
		}
		found++
		fmt.Fprintf(&b, "\n===== CARRY-FORWARD LEDGER round %d =====\n%s\n", r, string(data))
	}
	if found == 0 {
		return "", "there are no earlier rounds to compact"
	}
	b.WriteString("\nAn item above is disposed of only by its own owner. If you see an objection you\n" +
		"did not raise and its owner has not withdrawn it, it is still live. Full artifacts remain\n" +
		"on disk; open any of them.\n")
	return b.String(), ""
}

func fallbackTo(full func() (string, error), why string) (string, error) {
	s, err := full()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"\n===== FULL HISTORY (carry-forward fallback) =====\n"+
			"Reason: %s.\nEvery prior-round artifact follows in full.\n%s", why, s), nil
}

func renderRound(dir, label string) (string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "_index.md" || e.Name() == ledgerFileName {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	var b strings.Builder
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "\n===== %s/%s =====\n%s\n", label, name, string(data))
	}
	return b.String(), nil
}

// roundContextSentence states what a cross-review prompt ACTUALLY contains.
//
// FIX-UP CYCLE 3. Cycle 2's wording said older rounds "appear either in full or as a carry-forward
// ledger; a banner above says which" — but with compaction disabled no banner is ever emitted, so
// the sentence an agent received was false. codex-1 raised it as MAJOR and kimi-1 as NIT; they
// agreed on the fact and disagreed only on severity, and a prompt that misdescribes its own
// contents is worth fixing at either severity.
//
// Deriving it from the same constant that controls the behaviour means the sentence cannot drift
// from what the agent is actually given.
func roundContextSentence() string {
	if compactionEnabled {
		return "Older rounds appear either in full or as a carry-forward ledger; a banner above says which."
	}
	return "Every prior round appears below in full."
}
