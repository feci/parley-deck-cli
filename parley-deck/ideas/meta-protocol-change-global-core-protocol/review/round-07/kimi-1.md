---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 7
date: 2026-08-07
reviewed-commit: 134c52a
verdict: FINDINGS
---

## Summary

Reviewing cycles 5 and 6 together (my round-6 run was interrupted). Cycle 6 makes the right
architectural call: the multiset was defeated five times because it cannot see order, and the LCS
sequence diff over exact lines is the correct shape. I verified the algorithm itself hard —
recurrence, backtrack, tie-break direction, multiplicity, CRLF handling, attribution, performance
— and it is sound, with two narrow exceptions below. The 134c52a heredoc-mangle restore touches
only IMPLEMENTATION.md prose and is faithful.

The cycle is not CLEAN, and the blocker is the same disease as cycle 3: **the tests do not pin the
fix.** I reverted `droppedContent` to the exact `8ed3c4b` multiset in a scratch archive of 134c52a
and ALL SIX `TestDroppedContent...` tests passed — the three new round-6 tests cannot distinguish
LCS from the defeated implementation. The task asked me to confirm each round-6 probe test fails
under that revert; none of them fails. Beyond that, I found one genuine (narrow) silent-erasure
path the new stamp rule introduces, and one normalization the narrowed claim does not disclaim.

Required commands at 134c52a (clean tree, repo-local caches; all repo-level runs predate the
uncommitted working-tree edit that appeared at 14:58 — see note at the end):

```text
$ go build ./...                  BUILD_OK (no output, exit 0)
$ GOOS=windows go build ./...     WIN_OK (no output, exit 0)
$ GOOS=linux go build ./...       LINUX_OK (no output, exit 0)
$ go vet ./...                    VET_OK (no output, exit 0)
$ go test ./... -count=1          25 packages ok, 2 no-test-files; internal/app 26.815s,
                                  internal/runner 9.458s; exit 0
$ go test ./internal/app -run '^(TestProtocol|TestPublish|TestCore|TestLoad|TestDroppedContent)' -count=1 -v
  --- PASS: 29    --- FAIL: 0
$ grep -c '^func Test' internal/app/protocol_test.go
  29
```

"29 protocol tests" in IMPLEMENTATION.md reproduces two ways. The runner failure codex-1 saw at
round 6 (`TestDurableKillEndToEndRealProcess`, boot-id verification) did NOT reproduce on this
host — runner green here, as it was in my round-5 run; consistent with his non-attribution of it.

## Prior findings: fixed or not

**codex-1 round-6 [MAJOR] — four meaning-changing transformations reported as no loss.**
FIXED IN THE IMPLEMENTATION, UNPINNED IN THE SUITE. I re-probed all four at 134c52a against the
LCS (unit probes in a scratch tree):

- Hard break: exact-line comparison treats `Core rules live here.␣␣` ≠ `Core rules live here.` —
  reported. (Covered by the new test only vacuously — see MAJOR below.)
- HTML-comment reorder: order-sensitive LCS leaves the moved prose unmatched — reported.
- Duplicate-heading context (NO permanent test exists for this one): prior
  `## Parent A / ### Repeated / Shared rule.` vs render with the rule under `## Parent B /
  ### Repeated` → report `[### Repeated — 2 lines not carried forward]`. LOUD. The fix works;
  IMPLEMENTATION.md's "All four of codex-1's cases are permanent tests" is nonetheless false —
  only three test functions were added (protocol_test.go:660-699).
- Second stamp-prefix header line: forgiven exactly one, second reported — but ONLY when the deck
  stamp differs from the regenerated one. The same-version case is my new MINOR below.

**codex-1 round-6 [MAJOR, carried] — G7b guarantees no test covers.** STILL OPEN, re-verified at
134c52a with two of his six mutants in fresh scratch archives (mutant verified landed by grep,
then the full 29-test selector run):

| Mutant at 134c52a | Result |
| --- | --- |
| remove final-file `O_NOFOLLOW` (core.go:167 → `O_WRONLY\|O_CREATE\|O_EXCL`) | 29/29 PASS |
| accept untrimmed versions (`ValidVersion(strings.TrimSpace(version))`) | 29/29 PASS |

The documented guarantees (IMPLEMENTATION.md:21-25,137-139; changelog "O_EXCL|O_NOFOLLOW") remain
invisible to the suite. Cycle 6 did not claim to address this; it must not be lost.

**codex-1 round-6 [MINOR, carried] — §7 absolute.** STILL PRESENT at 134c52a:
parley-deck/COOPERATION.md:777 ("What IS in force today … and no agent-accessible code path
writes a release"), pty concession two sentences earlier; embedded copy identical (§7 SHA-256
`8ff1d896…` matches the live deck and the skill SOURCE repo at parley-deck-skill@f57f114, clean
tree). protocol-changelog.md:25 "no agent-accessible write path" and IMPLEMENTATION.md:74-75
"durable guarantees … no agent-accessible write path" likewise survive.

**codex-1 round-6 [MINOR, carried] — Phase-8 record.** SAME LAG, THIRD CYCLE:
IMPLEMENTATION.md:8 now says `head-commit: 8ed3c4b (cycle 5); cycle 6 lands in the next commit`.
Reviewed commit is 134c52a.

**Rounds 3-5 CLEAN items — unbroken by cycle 6.** Full suite green at 134c52a (quoted above);
`TestEmbeddedDefaultMatchesLiveDeck` PASS; the cross-section, Load/Publish symlink, render
read-error, idempotence, and all three cycle-5 `TestDroppedContent` tests pass under the LCS.
Note one fragility: `TestDroppedContentDetectsSyncedPrefixOutsideTheHeader` passes only because
the fixture's old stamp differs from the regenerated stamp, so the real stamp consumes the single
forgiveness first (document order). Correct as designed, but it is the same-version sibling of my
new MINOR.

## New findings (by severity, or "none")

### [MAJOR] None of the round-6 probe tests distinguishes the LCS from the defeated multiset

`git archive 134c52a` → scratch; replaced ONLY `internal/protocolcore/render.go` with the exact
`8ed3c4b` (cycle-5 multiset) version; kept all 134c52a tests. Result:

```text
=== RUN   TestDroppedContentDetectsIndentationLoss            --- PASS
=== RUN   TestDroppedContentDetectsSyncedPrefixOutsideTheHeader --- PASS
=== RUN   TestDroppedContentTreatsDeeperHeadingsAsSections    --- PASS
=== RUN   TestDroppedContentDetectsMarkdownHardBreakLoss      --- PASS
=== RUN   TestDroppedContentDetectsReorderedContent           --- PASS
=== RUN   TestDroppedContentForgivesExactlyOneStamp           --- PASS
PASS
ok      parley-deck-cli/internal/app      0.324s
```

Per-test root causes (my analysis, matching what the revert run proves):

1. **HardBreakLoss** asserts only `Contains("## 3. Phases")`. The fixture adds `Next line.`,
   which is foreign to the render and reported under ANY implementation — so the heading is
   always present while the hard-break line itself still compares equal under the multiset's
   TrimRight. Under LCS the report is `## 3. Phases — 2 lines`; under the multiset `— 1 line`.
   The assertion cannot see the difference.
2. **ReorderedContent** is not codex-1's probe. His probe reordered THREE SHARED lines (deck
   `<!--, Core rules live here., -->` vs core `Core rules live here., <!--, -->`). The permanent
   test's fixture core contains no comment lines at all, so `<!--` and `-->` are simply foreign
   and both algorithms report `## 3. Phases — 2 lines`. The reorder essence — same lines,
   different order — is never exercised; even a count assertion could not save this fixture.
3. **ForgivesExactlyOneStamp** asserts `Contains("not carried forward")`, which the shared
   fixture's §99 (`## 99. Project-specific rule — 2 lines`) produces on EVERY render of any
   implementation. Tautological: it cannot fail on this fixture.

Combined with the false "All four of codex-1's cases are permanent tests" (IMPLEMENTATION.md,
cycle-6 section) and the missing duplicate-heading test, the net state is: **the LCS rewrite —
the cycle's entire production change — has zero regression protection.** A cycle-7 revert to the
multiset keeps the suite green. This is the cycle-3 tautology class and it blocks CLEAN.
(hermes-1's round-7 review reached the same MAJOR independently; the uncommitted working-tree
edit now in progress appears to be the fix — see note below.)

### [MINOR] The exactly-one-stamp rule mis-fires when no stamp was removed — a real silent path

`render.go:216-243`: `addedStamp` is true whenever the RENDER contains a stamp, and the first
UNKEPT stamp-prefixed deck line is then forgiven. When the deck's stamp already equals the
regenerated one (re-render at the same pinned version — the normal state of a converged deck),
the LCS CARRIES the stamp, so zero stamps were removed — yet the forgiveness is still spent on
the next stamp-prefixed line. Probe (scratch, command path): render a converged deck, hand-add
`**Protocol synced:** genuine project prose` under §3, `render --dry-run`:

```text
preserved from this deck: Workspace, Transport, Created, §2 roster table
would regenerate .../parley-deck/COOPERATION.md from core 1.0.0 (d9181eb71672).
Nothing was written. Re-run with --yes to apply.
--- FAIL: SILENT: stamp-prefixed prose in §3 vanished unreported on a same-version re-render
```

No removal block; `--yes` would delete the line without saying so. This is G1's silent-erasure
class, and a regression relative to cycle 5's per-section `sectionHasStamp`, which reported this
exact placement (prose in §3, stamp in the header). The docstring's own words — the forgiven
stamp is "the one that was regenerated" — are contradicted: here the regenerated one was carried,
not removed. I rate it MINOR, not MAJOR, for the narrow precondition (hand-added stamp-prefix
prose on an already-converged deck); codex-1's round-7 review rates the same finding MAJOR. Fix
direction: forgive only a removed line that occupies the structural stamp slot (or matches the
generated-stamp shape), not any first prefix match.

### [NIT] Blank-line-only divergence is silently normalized

`render.go:236` skips reporting any removed line with `TrimSpace(l) == ""`. Probe:
prior `## S\n\nfoo\n\nbar` vs rendered `## S\n\nfoo\nbar` → report `[]`. A paragraph split/join
(or tight↔loose list change) is Markdown-significant, and the narrowed claim ("reports what
regeneration will change") does not disclaim this exemption. Narrow trigger — the deck's
non-blank content must be an exact subsequence of the render, i.e. essentially a converged deck
hand-edited in blank lines only. codex-1 (round 7) additionally notes `splitLines`' bare
`TrimSuffix("\r")` equates a lone terminal CR outside any CRLF pair — a second undocumented
normalization; exotic, but the docstring says "Nothing else is normalized." Either document both
exemptions or report blank-line removals.

### [NIT] No input bound behind the quadratic table

Measured (scratch unit probe): 2,000 lines → 13 ms; 8,000 lines → 223 ms and ~256 MB of `int32`
cells; codex-1's benchmarks (102 MB at 5k) extrapolate to ~1.6 GB at 20k lines. The comment's
"~1-2k lines … a few MB … milliseconds" is accurate for the intended input, but nothing enforces
it and the deck is a file agents write — a runaway append would OOM the next `render`/`check`.
Not pathological at any realistic size; worth a line-count guard or a Hirschberg/Myers note for
later. codex-1 rates this MAJOR; I don't — the input is the project's own committed file, and
every preceding stage of the pipeline already trusts it.

## Test-quality assessment

29 protocol tests, all green at 134c52a; full suite green on this host including internal/runner.
The LCS core is well-behaved under direct probe: loud tie-break (on equal-length choices it
reports the deck line rather than hiding it — I failed to construct any silent MOVE; swaps and
cross-heading moves all report), multiplicity handled by sequence matching, removed-heading
attribution correct (`## Gone — 3 lines` for heading + two body lines), CRLF-only diffs clean,
performance fine at documented scale.

But the suite as a regression gate for THIS cycle is worth nothing: the three new tests pass
under the implementation they were written to replace (demonstrated above), the fourth case has
no test, and the two new edge silences (same-version stamp, blank lines) have no tests. What a
discriminating fixture looks like is on record from the probe work: isolate the transformation so
the deck differs from the core ONLY by it, and assert the exact section-and-count string (e.g.
`## 3. Phases — 2 lines` vs `— 1 line` distinguishes the hard-break case; `(document header)`
distinguishes the stamp case; a true three-shared-line reorder distinguishes the reorder case).

Carried, re-confirmed at 134c52a, none claimed by cycle 6: the six G7b mutants (two re-run by me
this round, both still invisible); the §7 absolutes in all three maintained copies plus changelog
plus IMPLEMENTATION.md; protocol.go:301 "hash detection"; the tautological O_NOFOLLOW half of
`TestPublishRefusesExistingReleaseDirAndSymlinks` (consistent with mutant 1 passing);
`TestPublishRejectsUnsafeVersions` lacking a padded case (consistent with mutant 2 passing);
`parley protocol` still absent from docs/cli-reference.md and the root CHANGELOG; no GOOS
cross-build check in scripts/CI. Environmental: the INSTALLED skill copy
(~/.kimi-code/skills/parley-deck/references/COOPERATION.md) is still stale — §7 hash
`cf7186b1…`, missing the entire core-store block; the three maintained copies agree at
`8ff1d896…`. No text claims the installed copy is current; restating my round-5 call that it
needs a re-sync.

**Process note:** during this review the working tree diverged from 134c52a
(`internal/app/protocol_test.go` modified 14:58, uncommitted — a rewrite of the three probe tests
with isolated fixtures, citing hermes-1's identical MAJOR). Every result quoted above was
measured on the pristine tree or a `git archive 134c52a` extraction before that edit landed. The
review targets the commit as instructed; the in-flight edit does not change any conclusion here —
it addresses exactly the MAJOR above, and its own reversion-sensitivity will need the same
revert-and-watch-it-fail proof next round.

Verdict: FINDINGS. The algorithm is right; the proof is absent, and one narrow silence is real.
