---
agent: hermes-1
idea: tui-round-summary
review-round: 2
date: 2026-07-04
reviewed-commit: 2b25f08
responding-to: [hermes-1/review/round-01, codex-1/review/round-01]
---

## Position changes since prior review round

My MINOR-1 and MINOR-2 are resolved. I checked codex-1's round-01 findings as well: the
MAJOR (Home digest pushing Recent runs off-screen) and the MINOR (UTF-8 rune-split in
firstSentence) are both in scope for this re-review and I assessed them against the same
2b25f08 tree.

Verification performed on commit 2b25f08:

- internal/tui/roundsummary.go:88 now renders the degraded-extraction marker as a
  leading `warnStyle.Render("~")` prefix placed BEFORE `truncateText(ln.Position, …)`.
  The `~` glyph is emitted unconditionally on `ln.Fell` rows and cannot be eaten by the
  width truncation, because truncation now applies only to the position text, not to the
  combined "position + tag" string. The honesty signal survives narrow terminals. MINOR-1
  is fixed.
- internal/driver/driver.go:263-267 sets `nextAction = "consensus-ready (manual draft)"`
  when `d.cfg.Consensus == nil`, and `"drafting consensus"` only on the wired path. This
  matches the driver's actual behavior on the next line (return `ActionConsensus` and
  halt for the nil case, `Consensus.Draft` for the wired case). MINOR-2 is fixed.
- `go test ./internal/driver ./internal/tui` → both packages PASS (0.857s / 0.334s).

## Updated findings

### MINOR-1 (was hermes-1/round-01) — degraded-extraction tag truncated away

Resolved. The `~` prefix (warnStyle) is structurally ahead of truncation and is present
on every `Fell` row regardless of terminal width or position length. No remaining
concern. Closed.

### MINOR-2 (was hermes-1/round-01) — next: says "drafting consensus" when gate unwired

Resolved. The nil-gate branch now reports "consensus-ready (manual draft)", which is an
accurate description of the slice-1 halt. The wired path still says "drafting consensus".
No remaining concern. Closed.

### NIT-1 (was hermes-1/round-01) — aggregate flags lack per-agent attribution

No change. The aggregate-counts trailer is the ratified design and remains acceptable for
v1. Ratified, no change required.

### MAJOR (was codex-1/round-01) — Home digest can push Recent runs off-screen

Not addressed in this fix-up cycle. codex-1 reported that `renderHome` inserts the digest
before Recent runs using `rows/3` as the budget and does not reserve rows for the runs
block below it, so on a short terminal with several ideas the digest can be the block that
clips Recent runs. I confirmed this is still the layout in internal/tui/live.go (digest
appended before the runs section, outer clip at `rows`). This is out of scope for the two
hermes-1 fix-ups this cycle targeted, but it is a real viewport-regression risk and should
be the next item picked up. I agree with codex-1's suggested fix: render Home from
reserved sections (reserve Recent runs header + 1 row + command hint, allocate digest from
the remainder) and add a small-height regression test asserting Recent runs survives.

### MINOR (was codex-1/round-01) — firstSentence byte-cap can split UTF-8 runes

Not addressed in this fix-up cycle. `firstSentence` still caps at 120 bytes
(internal/driver/digest.go:105-106), so a boundary landing inside a multi-byte rune would
produce invalid UTF-8 before marshal. The existing test uses ASCII only. This is a display
path so it is not a protocol blocker, but for non-English summaries it is worth a rune-cap
fix + a Unicode-boundary test. Lower priority than the Home-layout MAJOR.

### NIT-2 (was hermes-1/round-01) — digestView duplicates RoundDigest shape

No change. Structural duplicate is acceptable for v1; not worth a refactor now.

Summary: zero remaining issues from my round-01 list. Two carry-over items from codex-1's
round-01 list (the Home-layout MAJOR and the UTF-8 MINOR) remain open and are not blockers
for the hermes-1 fix-ups, but the Home-layout MAJOR should be the next thing addressed.

## Open questions

- Should the Home-layout MAJOR (codex-1/round-01) be scheduled as a follow-up slice
  before any further digest enrichment, given it is a viewport regression on small
  terminals?
- Is the UTF-8 rune-split MINOR worth fixing in this idea or folding into a broader
  "unicode-safe text helpers" pass across the TUI?
