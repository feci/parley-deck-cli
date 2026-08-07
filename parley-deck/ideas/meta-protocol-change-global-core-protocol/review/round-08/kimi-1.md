---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 8
date: 2026-08-07
reviewed-commit: c4a8b83
verdict: FINDINGS
---
## Summary

All four claimed fixes are behaviorally real, and I verified each one directly. Hirschberg is
not just small, it is correct: on 4,150 randomized inputs plus edge cases its mask always has
optimal LCS length and is always a genuine subsequence of the render. The stamp pairing is
correct in all four combinations. Whitespace-only lines are now reported. The three rewritten
probe tests DO fail when the whole diff is reverted to the multiset — hermes-1's round-7 MAJOR
is fixed.

The cycle nevertheless fails verification item 4 of this review ("confirm each round-6/7 test
fails when its specific fix is reverted") for the two NEW behavioral fixes: reverting only the
round-7 stamp pairing, or only the whitespace counting, leaves all 31 checked-in tests green.
Cycle 7 fixed the tests that pin the LCS against the multiset, then added two unpinned fixes on
top of it — the same defect class hermes-1 flagged, one layer up. There is also one residual
silent-deletion edge in the stamp pairing (genuine stamp-prefixed line ABOVE the stamp in the
header region). Verdict: FINDINGS.

Required commands at `c4a8b831eddb134817333ac51235346356e92063` (also reproduced in a clean
`git archive` scratch copy):

```text
$ go build ./...                  exit 0, no output
$ GOOS=windows go build ./...     exit 0, no output
$ GOOS=linux go build ./...       exit 0, no output
$ go vet ./...                    exit 0, no output
$ go test ./... -count=1          exit 0: 26 packages ok, 1 no-test-files (cmd/parley);
                                  internal/app 30.751s, internal/protocolcore 2.803s,
                                  internal/runner 9.352s (round-6/7 runner flake did not recur)
```

Test count: `grep -c '^func Test'` = 29 in `internal/app/protocol_test.go` + 2 in
`internal/protocolcore/render_test.go` = 31. Targeted run: `--- PASS` ×31, `--- FAIL` ×0.
The claimed "31 protocol/diff tests" reproduces.

## Prior findings: fixed or not

**[MAJOR hermes-1] round-6 probe tests not reversion-sensitive — FIXED.** In a scratch copy with
`render.go` replaced by the exact `8ed3c4b` multiset, all three rewritten tests fail:

```text
--- FAIL: TestDroppedContentDetectsMarkdownHardBreakLoss (0.00s)
--- FAIL: TestDroppedContentDetectsReorderedContent (0.00s)
--- FAIL: TestDroppedContentForgivesExactlyOneStamp (0.00s)
```

The `minimalCore` fixtures genuinely differ from the core only by the transformation under test.
The IMPLEMENTATION.md claim "verified in both directions — they now FAIL when the LCS diff is
reverted" is TRUE for the multiset reversion.

**[MAJOR codex-1] stamp exemption inherited by a genuine line — FIXED in behavior, NOT pinned by
any test.** White-box probes of `droppedContent` against a rendered body carrying
`**Protocol synced:** core 1.0.0 (abc123)`:

```text
A deck stamp == render stamp, no extra line      → no report          ✓
B equal + genuine stamp line BELOW the stamp     → 1 line reported    ✓ (the round-7 bug shape)
C deck stamp stale, no extra line                → no report          ✓ (old stamp forgiven)
D stale + genuine line below                     → 1 line reported    ✓ (only the genuine one)
genuine stamp line in a body section             → reported under it  ✓
```

The fix works. But reverting ONLY the round-7 pairing logic to the round-6 "first unmatched
stamp" form (keeping the LCS) leaves `TestDroppedContentForgivesExactlyOneStamp` PASSING — new
finding 1 below.

**[MAJOR codex-1] whitespace-only lines normalized away — FIXED in behavior, NOT pinned by any
test.** Probe: deck = rendered output + one extra blank line → `[## 3. Phases — 1 line not
carried forward]`. With the old `strings.TrimSpace(l) == ""` skip restored, the same probe
reports nothing AND all 31 checked-in tests still pass — new finding 2. The second half of
codex-1's round-7 note (splitLines' `TrimSuffix(..., "\r")` also equates a lone terminal CR with
no CR, beyond the documented CRLF normalization) was NOT addressed — carried NIT.

**[MAJOR codex-1] O(n·m) LCS memory — FIXED, and the algorithm is verified correct, not just
small.** Correctness: 4,000 randomized trials (n,m ≤ 80, alphabet including empty/whitespace/
heading/stamp lines, biased shared chunks) + 150 large duplicate-heavy trials (n,m 50–400,
alphabet of 5) + 9 deterministic edge cases (empty, singletons, all-equal, disjoint, reversed,
duplicates). For every input: marked count == reference full-DP LCS length, AND the marked
subsequence of `a` is a subsequence of `b`. Base cases are right (`len(a)==1` marks iff the
element occurs in `b`; empty either side returns); the split picks the maximal
`left[j] + right[len(b)-j]` with strict `>` (smallest maximizing j), which the fuzz's optimality
result confirms across recursion shapes. Memory measured at the 20k-line test size:

```text
Hirschberg (c4a8b83):   PASS 0.82s,  peak RSS    92,094,464 B (~92 MB)
full DP  (reverted):    PASS 1.08s,  peak RSS 1,653,211,136 B (~1.65 GB)
```

The ~1.6 GB claim reproduces; Hirschberg is an ~18× reduction at this size and linear. Caveat:
the 20k test does NOT fail when the DP is restored on a machine with enough RAM — it pins
correctness-at-scale, not the memory bound. Acceptable (the fix is the algorithm; a unit test
cannot portably assert RSS), but the record should not claim this test guards the bound.

## New findings (by severity)

### [MAJOR] The round-7 stamp-pairing fix has no reversion-sensitive test

`TestDroppedContentForgivesExactlyOneStamp`'s deck carries a STALE stamp (`old`) plus the genuine
line. Under both the round-6 and round-7 logic the stale stamp is the one forgiven and the
genuine line is reported — the outputs are identical, so the test cannot distinguish the fix
from the bug. Verified: with only the pairing reverted (LCS kept), the test passes:

```text
$ go test ./internal/app -run 'TestDroppedContentForgivesExactlyOneStamp' -count=1 -v
=== RUN   TestDroppedContentForgivesExactlyOneStamp
--- PASS: TestDroppedContentForgivesExactlyOneStamp (0.00s)
```

The round-7 bug shape — deck ALREADY carrying the current stamp, genuine stamp-prefixed line
below it, exemption unconsumed — is exercised by no checked-in test. Suggested pin: build the
deck by rendering once (so it carries the current stamp), append `**Protocol synced:** genuine
prose` after the stamp, and assert the report contains a header-region entry; that fixture fails
under the round-6 logic (nothing reported) and passes now. This is the same defect class as
hermes-1's round-7 MAJOR, on the cycle-7 fix itself, and it is a direct negative answer to this
round's verification item 4.

### [MAJOR] The whitespace-counting fix has no reversion-sensitive test

With `if keep[i] || strings.TrimSpace(l) == "" { continue }` restored in `droppedContent`, the
entire checked-in suite passes:

```text
ok  parley-deck-cli/internal/app            21.178s
ok  parley-deck-cli/internal/protocolcore    1.165s
```

`TestDroppedContentTreatsDeeperHeadingsAsSections` was loosened from `#### Parent — 3 lines` to
`#### Parent — `, which tolerates any count: current code reports `#### Parent — 5 lines`
(heading + two blanks + two body lines), the reverted code reports 3, and the assertion matches
both. The comment above it ("whitespace-only lines … are content too") describes behavior no
test enforces. Suggested pin: assert the exact count (`#### Parent — 5 lines`) or add codex-1's
blank-line probe (extra blank line in a freshly rendered deck → exactly 1 line reported).

### [MINOR] A genuine stamp-prefixed line ABOVE the stamp in the header is silently forgiven

`headerStamp` takes the FIRST stamp-prefixed line in the header region as the generated stamp.
That holds for the natural ordering (the render places the stamp immediately after
`**Created:**`), but a hand-edit can put a genuine stamp-prefixed line above it — e.g. directly
under the title. Probe, deck otherwise carrying the current stamp:

```text
prior header: # T … **Protocol synced:** genuine project prose
                     **Protocol synced:** core 1.0.0 (abc123)   ← current, kept by LCS
droppedContent(prior, rendered) → []        ← SILENT; apply would delete the genuine line
```

`deckStamp` binds to the genuine line, `forgive` is true (genuine ≠ render stamp), and the
unmatched genuine line equals `deckStamp`, so it consumes the exemption; the actual stamp is
kept by LCS and needs none. Same failure class as round-7's MAJOR (silent deletion of
stamp-prefixed prose), narrower trigger. The robust rule: bind the exemption to the structural
slot (the line immediately following `**Created:**`), not to the first prefix match in the
header region.

### [MINOR] IMPLEMENTATION.md still claims all four round-6 cases are permanent tests

`IMPLEMENTATION.md:295`: "All four of codex-1's cases are permanent tests." The duplicate-heading
case still has none (`grep -in duplicate internal/app/protocol_test.go` → 0 matches; the
cycle-7 diff added no such test). The BEHAVIOR is correct — my probe (content added under the
second of two identical `## S` headings) reports exactly `## S — 1 line not carried forward` —
but the text states coverage that does not exist. Either add the test or narrow the sentence.

## Test-quality assessment

The suite is green and the count is honest (31 = 29 + 2). The three rewritten probe tests are
now genuinely reversion-sensitive against the multiset, and the fixtures are properly minimal.
The LCS core is covered at unit level (`TestLCSSeesOrder`, `TestLargeDocumentDoesNotExhaustMemory`)
and my independent fuzz found no correctness defect in Hirschberg.

The remaining weakness is the one this fleet keeps re-finding: the two fixes cycle 7 actually
made in `droppedContent`'s exemption/reporting logic are guarded by no test that fails when THEY
are reverted. Cycle 7's IMPLEMENTATION.md says "every fix is now reverted-and-re-run rather than
declared" — for the stamp and whitespace fixes, whatever manual reversion was run was not encoded
in a permanent test, so the suite cannot hold them. Part 2 of the stamp test (single regenerated
stamp, expect NO report) is a good over-reporting guard and does pin the exemption's existence.

Reversion matrix (scratch copy of c4a8b83, one change at a time):

| Reversion                                  | Checked-in suite result        |
|--------------------------------------------|--------------------------------|
| whole diff → 8ed3c4b multiset              | 3 probe tests FAIL (correct)   |
| stamp pairing → round-6 "first unmatched"  | all 31 PASS (unpinned fix)     |
| whitespace counting → TrimSpace skip       | all 31 PASS (unpinned fix)     |
| Hirschberg → full DP                       | all PASS; 1.65 GB vs 92 MB RSS |

Nothing cycle 7 broke: no production file outside `internal/protocolcore/render.go` changed, the
full repo suite is green (including the previously flaky runner package), and rounds 3–5 CLEAN
behaviors still pass. One cosmetic note: header-region losses are now attributed to the level-1
title (e.g. `# T — 1 line`) rather than `(document header)`, because `heading()` matches `# `;
the `headerSection` constant is effectively unreachable. Pre-existing since cycle 5, harmless.

Carried text/test gaps (not claimed as fixed by cycle 7, restated so they are not lost): the
absolute "no agent-accessible code path writes a release" at `parley-deck/COOPERATION.md:777`
and `internal/protocol/defaults/COOPERATION.md:768` (pty concession is inline, absolute stands),
"no agent-accessible write path" at `meta/protocol-changelog.md:25` and `IMPLEMENTATION.md:74-75`,
the "hash detection" ambiguity in CLI output (`internal/app/protocol.go:301`), the untested
"the write is atomic and preserves the file mode" (`IMPLEMENTATION.md:38`), the lone-`\r` strip
in `splitLines`, and the `head-commit` field naming the previous cycle (`134c52a`) rather than
the commit under review — structurally inherent to a self-referential field, flagged since
round 5. The installed skill copy at `~/.kimi-code/skills/parley-deck/references/COOPERATION.md`
is older than both in-repo copies (no §7 blast-radius clause, no write-path paragraph) — the
three-copy identity asserted in rounds 3–4 does not hold on disk.
