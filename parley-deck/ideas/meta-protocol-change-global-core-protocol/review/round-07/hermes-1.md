---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
review-round: 7
date: 2026-08-07
reviewed-commit: 134c52a
verdict: FINDINGS
---
## Summary

The LCS sequence diff in `droppedContent` is a correct algorithm and a sound replacement
for the defeated multiset. I verified the LCS table, the backtracking, the exactly-one-stamp
rule, and section attribution when a heading itself is removed — all behave correctly. The
algorithm does not silently miss content, does not report carried content as lost, and is
performant for the ~1-2k-line protocol documents it handles.

The implementation is nevertheless not CLEAN. The three permanent tests added for codex-1's
round-6 probes are NOT reversion-sensitive: all three PASS with the LCS diff reverted to the
8ed3c4b multiset. The test assertions check for substrings that are always present regardless
of which diff algorithm is used, because the testDeck fixture always has "STALE core text that
must be replaced." under §3, which always produces a "## 3. Phases" report entry and a "not
carried forward" string. The reorder test is vacuous in a stronger sense: both algorithms
produce the same line count for its fixture, so no assertion based on that output could
distinguish them. The task explicitly asks to "confirm each round-6 probe test fails if the
LCS diff is reverted to a multiset" — they do not fail.

Cycle 6 did not break any test that rounds 3-5 had made CLEAN, and the LCS algorithm itself
is correct. The remaining findings are carried NITs/MINORs from prior rounds that cycle 6 did
not address, restated so they are not lost.

## Prior findings: fixed or not

**codex-1 round-6 MAJOR — `droppedContent` still reports no loss for meaning-changing
Markdown transformations.** FIXED IN ALGORITHM, NOT IN TESTS. The multiset was replaced by
an LCS sequence diff over exact lines (only CRLF's CR stripped). I verified the LCS is
correct (see "New findings" below for the test gap). The algorithm change itself is sound:
a moved line is a removal plus an addition, a lost hard break is a changed line, and
whitespace is compared as written.

**codex-1 round-6 MINOR — user-only wording still ends with an absolute claim.** NOT FIXED
(carried). `parley-deck/COOPERATION.md:777` still says "no agent-accessible code path writes
a release" after conceding a pty-allocating agent can publish. The embedded copy at
`internal/protocol/defaults/COOPERATION.md:768` carries the same text. The changelog at
`parley-deck/meta/protocol-changelog.md:25` says "no agent-accessible write path."
`IMPLEMENTATION.md:74-75` says "no agent-accessible write path, and hash detection."
All are carried NITs from round 3; the pty limitation IS stated inline in each case.

**codex-1 round-6 MINOR — Phase-8 record head-commit.** NOT FIXED (carried).
`IMPLEMENTATION.md:8` says `head-commit: 8ed3c4b (cycle 5); cycle 6 lands in the next commit`.
The reviewed commit is 134c52a. This is the same lag that codex-1 flagged in rounds 5 and 6:
the head-commit field always identifies the prior cycle, not the commit under review.

**codex-1 round-6 MAJOR — G7b test coverage for documented guarantees.** NOT FIXED (carried).
codex-1 showed six behaviors where removing the implementation did not fail any of the 26
protocol tests: O_NOFOLLOW on the final file create, untrimmed version rejection, CRLF-core
normalization, the check preamble, file-mode preservation, and atomic render. Cycle 6 only
changed `droppedContent` and added three tests; it did not add tests for these six. The
IMPLEMENTATION.md still says "the write is atomic and preserves the file mode" (line 38)
with no end-to-end test covering either claim.

**hermes-1 / kimi-1 round-5 carried NITs.** All still present, none claimed as fixed by
cycle 6: (1) `protocol.go:301` "hash detection" is ambiguous between drift detection
(tested) and release-tamper detection (not implemented); (2) the O_NOFOLLOW half of
`TestPublishRefusesExistingReleaseDirAndSymlinks` is tautological — the pre-created 3.0.0
directory triggers the directory-existence check before `O_NOFOLLOW` is reached; (3)
`TestPublishRejectsUnsafeVersions` has no padded/untrimmed case; (4) no CI check exists
for GOOS builds (`grep -rn GOOS scripts/ .github/` returns 0 matches); (5) the locally
installed skill copy at `~/.hermes/skills/parley-deck/references/COOPERATION.md` does not
carry the §7 blast-radius clause at all (kimi-1 round 5 noted this; it is still stale);
(6) `docs/cli-reference.md` still has no protocol command documentation; root `CHANGELOG.md`
has no protocol entry.

## New findings (by severity)

### [MAJOR] The three round-6 probe tests are not reversion-sensitive

The task asks to "confirm each round-6 probe test fails if the LCS diff is reverted to a
multiset." All three PASS under the multiset. I extracted the exact 8ed3c4b `render.go`
(the cycle-5 per-section multiset) into a scratch archive of 134c52a and ran each test.

```
$ go test ./internal/app/ -run 'TestDroppedContentDetectsMarkdownHardBreakLoss|TestDroppedContentDetectsReorderedContent|TestDroppedContentForgivesExactlyOneStamp' -v -count=1
=== RUN   TestDroppedContentDetectsMarkdownHardBreakLoss
--- PASS: TestDroppedContentDetectsMarkdownHardBreakLoss (0.00s)
=== RUN   TestDroppedContentDetectsReorderedContent
--- PASS: TestDroppedContentDetectsReorderedContent (0.00s)
=== RUN   TestDroppedContentForgivesExactlyOneStamp
--- PASS: TestDroppedContentForgivesExactlyOneStamp (0.00s)
PASS
ok  parley-deck-cli/internal/app  0.333s
```

Root cause: all three tests use `strings.Replace(testDeck, "STALE core text that must be
replaced.", ...)` or `strings.Replace(testDeck, "**Protocol synced:** ...", ...)` but the
testDeck fixture always has "STALE core text that must be replaced." under §3. After any
replacement, the deck still loses that line (the core has "Core rules live here." instead),
so the report always contains "## 3. Phases — N lines not carried forward" and the string
"not carried forward" — regardless of which diff algorithm is used.

Per-test analysis (output captured under both algorithms via diagnostic tests):

1. **TestDroppedContentDetectsMarkdownHardBreakLoss** — assertion:
   `Contains("## 3. Phases")`. Under LCS: "## 3. Phases — 2 lines" (hard break line + Next
   line). Under multiset: "## 3. Phases — 1 line" (only Next line; lineKey trims trailing
   spaces, equating the hard break with the core's line). The difference exists but the
   assertion only checks for the heading name, not the count. Would be reversion-sensitive
   with `Contains("## 3. Phases — 2 lines")`.

2. **TestDroppedContentDetectsReorderedContent** — assertion:
   `Contains("## 3. Phases")`. The fixture replaces STALE text with
   `"<!--\nCore rules live here.\n-->"`. Under BOTH algorithms: "Core rules live here."
   is matched, "<!--" and "-->" are dropped. Both produce "## 3. Phases — 2 lines." The
   test is vacuous — even a count-based assertion would not distinguish the algorithms.
   This is not the reordering probe codex-1 described (codex-1 reordered the same three
   lines); the test replaces one line with three different lines.

3. **TestDroppedContentForgivesExactlyOneStamp** — assertion:
   `Contains("not carried forward")`. Under LCS: the second stamp line is attributed to
   "(document header)" and reported. Under multiset: `sectionHasStamp` causes both stamp
   lines to be skipped. But "not carried forward" appears in both versions because of the
   §3 and §99 entries. Would be reversion-sensitive with `Contains("(document header)")`.

Suggested fix: tighten the assertions to check for the distinguishing output. For test 1,
assert `Contains("## 3. Phases — 2 lines")`. For test 2, replace the fixture with a true
reordering (the same lines in a different order, so the multiset stays equal while the LCS
changes). For test 3, assert `Contains("(document header)")` or
`Contains("# Parley Deck Cooperation Protocol")`.

### [MINOR] The reorder test fixture does not test what codex-1's probe tested

codex-1's round-6 probe #3 (HTML-comment boundary) placed the same three lines in a
different order: the deck had `<!--`, `Core rules live here.`, `-->` while the core had
`Core rules live here.`, `<!--`, `-->`. Under a multiset, all three match and nothing is
reported; under LCS, the reordering means the lines do not align. The test fixture instead
replaces "STALE core text that must be replaced." with `"<!--\nCore rules live here.\n-->"`
— introducing two new lines ("<!--" and "-->") that are not in the core at all. Both
algorithms drop them. This is a new-content test, not a reordering test, and it cannot
distinguish the two algorithms even in principle.

### [NIT] IMPLEMENTATION.md head-commit lags behind the reviewed commit

`IMPLEMENTATION.md:8` says `head-commit: 8ed3c4b (cycle 5); cycle 6 lands in the next
commit`. The reviewed commit is 134c52a. This is the same pattern flagged in rounds 5 and
6. The frontmatter should name 134c52a as the cycle-6 commit.

## Test-quality assessment

29 protocol tests (grep -c '^func Test' = 29), all pass. Full suite green including the
runner package that was flaky in codex-1's round-6 environment.

```
$ go build ./...                  EXIT=0 (no output)
$ GOOS=windows go build ./...     EXIT=0 (no output)
$ GOOS=linux go build ./...       EXIT=0 (no output)
$ go vet ./...                    EXIT=0 (no output)
$ go test ./... -count=1          26 packages: 24 ok, 2 no-test-files, EXIT=0
$ go test ./internal/app/ -run 'TestProtocol|TestPublish|TestCore|TestLoad|TestDroppedContent' -count=1 -v
  29 PASS, 0 FAIL
$ grep -c '^func Test' internal/app/protocol_test.go
  29
```

### LCS algorithm correctness

The LCS implementation is correct. I verified:

- **lcsMask**: classic DP table built bottom-up, then backtracked to mark which `a` lines
  belong to the LCS. The tie-breaking (`table[i+1][j] >= table[i][j+1]` → advance `i`)
  is standard and deterministic. `int32` is sufficient for documents of ~1-2k lines.
- **Content lost but not reported (false negative)**: a line is "kept" only if it exactly
  matches a rendered line in LCS position. A line absent from the render is never kept.
  No false negative found.
- **Content reported as lost but IS carried (false positive)**: `keep[i]` is true only
  when `a[i] == b[j]` in the LCS path. A line present in both at any position is matched
  by the LCS (which finds the longest subsequence). No false positive found.
- **Section attribution when a heading is removed**: the heading line updates `current`
  before the keep/skip check, so a dropped heading and its body lines are attributed to
  that heading's own name. Verified: a `#### Subsection` the core lacks reports as
  `#### Subsection — 2 lines` (heading + content). Correct.
- **Exactly-one-stamp rule**: `addedStamp` is set if the render contains any stamp-prefixed
  line. The first unkept stamp-prefixed deck line is skipped (`stampSkipped = true`); a
  second is reported. The flag is global, not per-section — this is intentional and correct
  for a single regenerated stamp. A second stamp-prefixed line of genuine project prose is
  reported under the `(document header)` section. Verified.
- **Performance**: O(n*m) table. For ~2000-line documents, ~4M int32 entries (~16MB),
  milliseconds. No pathological concern for the intended use case (protocol files). No
  upper bound on document size is enforced, but the input is an internal protocol file,
  not untrusted user data.

### Reversion sensitivity of round-6 probe tests

All three round-6 probe tests PASS under the 8ed3c4b multiset. They are not
reversion-sensitive. See the MAJOR finding above for the per-test analysis and suggested
fixes.

### Reversion sensitivity of rounds 3-5 tests

All eight tests that rounds 3-5 made CLEAN still pass under the LCS version. No regression:

- TestProtocolRenderReportsContentLostUnderASharedHeading — PASS
- TestProtocolRenderDetectsCrossSectionErasure — PASS
- TestDroppedContentDetectsIndentationLoss — PASS
- TestDroppedContentDetectsSyncedPrefixOutsideTheHeader — PASS
- TestDroppedContentTreatsDeeperHeadingsAsSections — PASS
- TestLoadRefusesASymlinkedReleaseDirectory — PASS
- TestPublishRefusesASymlinkedStoreComponent — PASS
- TestProtocolRenderReportsDeckReadErrors — PASS

The LCS naturally handles what the cycle-5 multiset patches (lineKey, sectionHasStamp,
all-ATX heading) were trying to approximate: exact string comparison means indentation and
trailing whitespace are preserved without a special `lineKey`, and the global sequence
comparison means cross-section coincidence cannot mask a deletion without per-section
indexing.

### Does any text still state a guarantee no end-to-end test covers?

Yes — these are carried from prior rounds, none claimed as fixed by cycle 6:

1. `IMPLEMENTATION.md:38` — "the write is atomic and preserves the file mode." Neither
   atomicity nor mode preservation has an end-to-end test. codex-1 round-6 mutants
   (replace atomic with `os.WriteFile`, force mode to 0644) both passed all 29 tests.
2. `parley-deck/COOPERATION.md:777` and `internal/protocol/defaults/COOPERATION.md:768` —
   "no agent-accessible code path writes a release." Overclaim for pty-allocating agents;
   the pty limitation is stated inline. (NIT, carried since round 3)
3. `parley-deck/meta/protocol-changelog.md:25` — "no agent-accessible write path." Same
   overclaim. (NIT, carried)
4. `internal/app/protocol.go:301` — "The durable guarantees today are write-once releases
   and hash detection." "Hash detection" is ambiguous between deck drift detection (tested
   by `TestProtocolCheckReportsHandEditAndNeverWrites`) and release-tamper detection
   (`DETECTED-UNATTRIBUTED`, not implemented). (NIT, carried)
5. `IMPLEMENTATION.md:74-75` — "no agent-accessible write path, and hash detection." Same
   two carried NITs. (NIT, carried)
6. Changelog "O_EXCL|O_NOFOLLOW" — the O_NOFOLLOW file-level path is not pinned by a test
   (the existing test's O_NOFOLLOW half is tautological: the pre-created directory triggers
   the directory-existence check first). (NIT, carried)
7. `IMPLEMENTATION.md:148` — CRLF-core normalization fix. No test covers a CRLF core.
   (NIT, carried from codex-1 round-6 mutant)
8. The locally installed skill copy does not carry the §7 blast-radius clause. No repo text
   claims it is current, so this is not a new finding, but rounds 3-4 asserted three-copy
   identity and that is no longer true on disk. (NIT, carried from kimi-1 round 5)
