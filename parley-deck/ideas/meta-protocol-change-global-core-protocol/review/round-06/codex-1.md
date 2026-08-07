---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 6
date: 2026-08-07
reviewed-commit: 8ed3c4b
verdict: FINDINGS
---

## Summary

The three cycle-5 regression fixes are real, and each new test fails for its intended reason when
only that fix is reverted. The substantive §7 text also agrees across the live deck, embedded
default, and the clean committed skill repository.

The implementation is nevertheless not CLEAN. `droppedContent` still silently misses
meaning-changing transformations: Markdown hard-break whitespace can disappear; content can move
between duplicate-heading contexts; HTML-comment boundaries can change whether prose is visible;
and a second genuine `**Protocol synced:**` line in the document header is skipped merely
because the rendered header contains the generated stamp. All four command-path probes produced a
regeneration with no dropped-content entry at all.

Two round-5 findings also remain open despite not appearing in the claimed-fix list: documented
G7b guarantees still survive with all 26 protocol tests after their implementations are removed,
and the narrower user-only wording still ends with the absolute claim that no agent-accessible
write path exists. The Phase-8 status is now named, but its `head-commit` still identifies
`1903e49`, not reviewed HEAD `8ed3c4b`.

Required command results at `8ed3c4b2eb03cd8e527a1d5da0078a9a4698c3d8`:

```text
$ go build ./...
[no stdout/stderr]
exit 0

$ GOOS=windows go build ./...
[no stdout/stderr]
exit 0

$ GOOS=linux go build ./...
[no stdout/stderr]
exit 0

$ go vet ./...
[no stdout/stderr]
exit 0

$ go test ./...
...
ok      parley-deck-cli/internal/app       35.798s
...
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got
    {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001
    Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL    parley-deck-cli/internal/runner    7.350s
FAIL
exit 1
```

The runner failure reproduces alone. `1903e49..8ed3c4b` contains no runner change, and round
5 recorded the same failure at the baseline, so I do not attribute it to cycle 5. The protocol
suite itself is green:

```text
$ go test ./internal/app -run '^(TestProtocol|TestPublish|TestCore|TestLoad|TestDroppedContent)' -count=1 -v
...
PASS
ok      parley-deck-cli/internal/app       0.207s
PROTOCOL_TEST_COUNT runs=26 pass=26 fail=0
```

## Prior findings: fixed or not

- **PRIMARY — FIXED: indentation comparison.** Baseline
  `TestDroppedContentDetectsIndentationLoss` passes. In an isolated archive of
  `8ed3c4b`, changing only `lineKey` back to `strings.TrimSpace` makes it fail:

  ```text
  indentation-significant content was lost without a report
  FAIL parley-deck-cli/internal/app
  ```

- **PRIMARY — FIXED for the tested non-header location: synced-prefix scope.** Baseline
  `TestDroppedContentDetectsSyncedPrefixOutsideTheHeader` passes. Reverting only the
  conditional skip to `strings.HasPrefix(t, syncedPrefix)` makes it fail:

  ```text
  project prose beginning with the stamp prefix vanished unreported
  FAIL parley-deck-cli/internal/app
  ```

  The implementation remains incomplete for multiple same-prefix lines in the actual header; that
  is part of the new MAJOR below.

- **PRIMARY — FIXED for distinct deeper headings: all ATX levels.** Baseline
  `TestDroppedContentTreatsDeeperHeadingsAsSections` passes. Reverting only
  `heading` to the old `##`/`###` recognition makes it fail:

  ```text
  a level-4 subsection was not treated as its own section
  FAIL parley-deck-cli/internal/app
  ```

  Duplicate headings still collapse into one map key; that is part of the new MAJOR below.

- **PRIMARY — FIXED: the third protocol copy is committed and agrees.** The sibling skill
  repository is clean. `4b80468` commits the previously uncommitted bundled-snapshot
  corrections; its later commit `f57f114` commits the cycle-5 user-only rewording. The
  current working copy is byte-identical to that repository's `HEAD`. The live, embedded,
  and skill `## 7` sections have the same SHA-256:

  ```text
  8ff1d8966c4ebee6c1439b54f58fd51eed8e454912267ce4d897d58acddfe142
  ```

  The whole embedded/skill diff contains only their intended bootstrap `Transport` and
  `Created` identity values. The live/embedded drift guard also passes:

  ```text
  === RUN   TestEmbeddedDefaultMatchesLiveDeck
  --- PASS: TestEmbeddedDefaultMatchesLiveDeck (0.00s)
  PASS
  ```

- **PRIMARY — PARTIALLY FIXED, still open [MINOR]: user-only wording.** The rewritten opening at
  `parley-deck/COOPERATION.md:759-766` now accurately calls user control a rule backed by a
  limited mechanism. But the same section still says, as part of “What IS in force today,”
  “no agent-accessible code path writes a release” at line 777 after conceding that a
  pty-allocating agent can publish. The embedded and skill copies agree with that absolute.
  `parley-deck/meta/protocol-changelog.md:23-27` and `IMPLEMENTATION.md:71-75` retain
  the same overclaim; the latter calls it a “durable guarantee.” Reword or remove those surviving
  absolutes.

- **PRIMARY — PARTIALLY FIXED, still open [MINOR]: Phase-8 record.** The frontmatter now correctly
  says `status: fix-up-cycle-5`, but `IMPLEMENTATION.md:8` says
  `head-commit: 1903e49 (cycle 4); cycle 5 lands in the next commit`. The reviewed commit is
  `8ed3c4b`; the claimed “reviewed commit” fix therefore did not land.

- **PRIMARY — NOT FIXED, carried [MAJOR]: G7b test coverage.** Round 5 showed that several landed
  guarantees had no distinguishing end-to-end test. Cycle 5 added only the three
  `droppedContent` tests. In six independent `8ed3c4b` scratch archives I removed one
  behavior at a time and ran all 26 protocol tests:

  | Removed or falsified behavior | Current claim | Result |
  | --- | --- | --- |
  | remove final-file `O_NOFOLLOW` | `IMPLEMENTATION.md:21-24,71-75`; changelog line 23 | 26/26 PASS |
  | accept untrimmed versions | `IMPLEMENTATION.md:21-25,137-139` | 26/26 PASS |
  | remove CRLF-core normalization | `IMPLEMENTATION.md:147-150` | 26/26 PASS |
  | restore the false `check` preamble | `IMPLEMENTATION.md:186-189` | 26/26 PASS |
  | force render mode to `0644` | `IMPLEMENTATION.md:37-38` | 26/26 PASS |
  | replace atomic render with `os.WriteFile` | `IMPLEMENTATION.md:37-38` | 26/26 PASS |

  Each mutant compiled and returned `ok parley-deck-cli/internal/app`. Thus text still states
  guarantees no real-entry test covers, directly contrary to G7b. The in-repo drift test also
  covers only live/default, not the third skill copy; that copy was verified manually this round.

## New findings (by severity, or "none")

### [MAJOR] `droppedContent` still reports no loss for meaning-changing Markdown transformations

`droppedContent` is a multiset comparison keyed only by a literal heading string
(`internal/protocolcore/render.go:207-240,270-300`). It strips all trailing spaces
(`render.go:265-268`), ignores line order, merges every occurrence of the same heading, and
treats the presence of any stamp in a section as permission to skip every same-prefix prior line
in that section (`render.go:223-228,255-262`).

Four tests added only to an isolated `8ed3c4b` archive drove the
`runProtocol` `render --dry-run` command path:

1. **Markdown hard break:** before a following `Next line.`, the deck has
   `Core rules live here.  ` (two trailing spaces, a hard line break); the core has
   `Core rules live here.`. `lineKey` equates them and the renderer silently removes the
   hard break.
2. **Duplicate-heading context:** the deck places `Shared rule.` beneath
   `## Parent A` / `### Repeated`; the core moves it beneath `## Parent B` /
   `### Repeated` while moving a placeholder the other way. The map key is only
   `### Repeated`, so both multisets are identical and the rule's governing context changes
   silently.
3. **HTML-comment boundary:** the deck has `<!--`, `Core rules live here.`,
   `-->`; the core has the same three lines reordered as `Core rules live here.`,
   `<!--`, `-->`. The prose changes from hidden to visible, but the multiset is
   unchanged.
4. **Second header stamp-prefix line:** the deck header has its normal old stamp plus
   `**Protocol synced:** this is genuine project prose`. Because the rendered header contains
   one regenerated stamp, `sectionHasStamp` causes both prior prefix lines to be skipped and
   the project line vanishes.

Every probe failed with the same relevant output shape — there was no
`deck content NOT carried forward` block at all:

```text
preserved from this deck: Workspace, Transport, Created, §2 roster table
would regenerate .../parley-deck/COOPERATION.md from core ...
Nothing was written. Re-run with --yes to apply.
```

Representative failures were:

```text
Markdown hard-break whitespace was lost without a report
content moved between duplicate-heading contexts without a report
HTML-comment boundaries changed meaning without a report
a second same-prefix header line vanished without a report
```

This is G1's silent-erasure class, not a formatting preference. Trimming “only trailing
whitespace” is still lossy Markdown normalization; “same section” is not an exact generated-stamp
slot; and a heading label is neither a full heading path nor an occurrence identity. Fix the
report around a structure that preserves Markdown-significant bytes, ordering, heading
path/occurrence, fenced/comment state, and the exact generated stamp occurrence. A real Markdown
parser or the ratified stable block IDs is safer than another line-normalization exception. Add
all four probes as permanent preview/apply tests.

## Test-quality assessment

The three cycle-5 tests are reversion-sensitive and fail for their intended reasons. All 26
protocol tests pass at baseline, including the earlier cycle-3/4 regression tests:

- `TestProtocolRenderDetectsCrossSectionErasure`
- `TestLoadRefusesASymlinkedReleaseDirectory`
- `TestPublishRefusesASymlinkedStoreComponent`
- `TestProtocolRenderReportsDeckReadErrors`

The new tests are named `TestDroppedContent...`, so the earlier selector
`^(TestProtocol|TestPublish|TestCore|TestLoad)` still runs only 23 tests. The 26 count
reproduces when `TestDroppedContent` is included, or by counting the 26 `func Test...`
declarations in `internal/app/protocol_test.go`.

Cycle 5 did not regress the previously pinned production-dispatch, write-once, path-traversal,
Load/Publish store-component, Load release-directory, status/read-error, render/read-error,
cross-section, Windows-build, or Linux-build behavior. The required full suite is still not green
on this host because of the unchanged runner boot-ID failure quoted above.

Test quality remains insufficient for CLEAN: the new fixtures pin exactly their three fixes but do
not challenge the next semantic boundary, while six documented implementation guarantees remain
invisible to the entire 26-test protocol suite. The repository needs structural transformation
tests and real-entry tests that fail when each advertised boundary is removed, not only green
tests of nearby behavior.
