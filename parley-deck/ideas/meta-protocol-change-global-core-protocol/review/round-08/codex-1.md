---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 8
date: 2026-08-07
reviewed-commit: c4a8b83
verdict: FINDINGS
---
## Summary

Cycle 7 fixes the three defective round-6 fixtures: each now fails when `render.go` alone is
restored to the `8ed3c4b` multiset. The Hirschberg implementation is also a correct linear-space
LCS reconstruction. An independent full-DP oracle found the same optimum on 50,000 randomized
input pairs, with every returned mask forming a valid subsequence of the rendered input. The four
requested stamp combinations are correct when the genuine stamp-prefixed line follows the
structural stamp, and an end-to-end whitespace-only probe is now reported.

The result is not CLEAN. Stamp pairing still identifies the first prefix-matching header line,
not the structural generated-stamp slot, so genuine stamp-prefixed prose placed before an already
current stamp is silently erased. More broadly, none of the three cycle-7 implementation fixes is
reversion-pinned: reverting only stamp pairing or only whitespace reporting leaves all 31 scoped
tests green, and the 20k test passes with the old 1.6-GB full DP matrix on this host. The carried
G7b documentation/test gap is also still open: five independent safety/reporting mutants leave
all 31 tests green, while all three protocol copies and the changelog retain a safety absolute
that a real agent-driven PTY publish disproves.

Required command evidence at `c4a8b831eddb134817333ac51235346356e92063` (`PRIMARY`):

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
```

`go test` is not green on this machine. Its one failing package/test is the same runner boot-ID
failure recorded in rounds 6 and 7; `git diff --quiet 134c52a c4a8b83 -- internal/runner` exits 0,
so I do not attribute it to cycle 7.

```text
$ go test -count=1 ./...
?    parley-deck-cli/cmd/parley                 [no test files]
ok   parley-deck-cli/internal/acp                0.352s
ok   parley-deck-cli/internal/agents             0.635s
ok   parley-deck-cli/internal/app                50.765s
ok   parley-deck-cli/internal/config             1.360s
ok   parley-deck-cli/internal/consensus          1.635s
ok   parley-deck-cli/internal/driver             5.019s
ok   parley-deck-cli/internal/fsutil             1.065s
ok   parley-deck-cli/internal/hitl               2.126s
ok   parley-deck-cli/internal/loop               2.981s
ok   parley-deck-cli/internal/pipeline           2.716s
ok   parley-deck-cli/internal/procctl            2.532s
ok   parley-deck-cli/internal/protocol           0.794s
ok   parley-deck-cli/internal/protocolcore       4.040s
ok   parley-deck-cli/internal/repomap            3.192s
ok   parley-deck-cli/internal/retro              3.058s
ok   parley-deck-cli/internal/runaction          3.044s
ok   parley-deck-cli/internal/runcontrol         3.051s
ok   parley-deck-cli/internal/runmanifest        3.023s
--- FAIL: TestDurableKillEndToEndRealProcess (0.03s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL parley-deck-cli/internal/runner             10.206s
ok   parley-deck-cli/internal/runplan            2.921s
ok   parley-deck-cli/internal/runstate           3.021s
ok   parley-deck-cli/internal/sessionstore       2.989s
ok   parley-deck-cli/internal/steer              2.783s
ok   parley-deck-cli/internal/store              2.772s
ok   parley-deck-cli/internal/track              2.829s
ok   parley-deck-cli/internal/tui                3.079s
FAIL
exit 1
```

The repository contains 635 top-level Go `Test*` functions. The implementation's narrower
protocol/diff count does reproduce:

```text
$ go test -count=1 -json ./internal/app ./internal/protocolcore \
    -run '^(TestProtocol|TestPublish|TestCore|TestLoad|TestDroppedContent|TestLargeDocument|TestLCS)' \
  | awk '/"Action":"pass".*"Test":/ {p++} /"Action":"fail".*"Test":/ {f++} \
         /"Action":"skip".*"Test":/ {s++} END {printf \
         "PROTOCOL_DIFF_TEST_COUNT pass=%d fail=%d skip=%d total=%d\\n", p, f, s, p+f+s}'
PROTOCOL_DIFF_TEST_COUNT pass=31 fail=0 skip=0 total=31
exit 0
```

## Prior findings: fixed or not

- **[MAJOR hermes-1] Three round-6 probe tests were not reversion-sensitive — FIXED.** In a
  `c4a8b83` archive, I replaced only `internal/protocolcore/render.go` with the exact `8ed3c4b`
  multiset and ran the three tests. All three failed for their intended transformation:

  ```text
  --- FAIL: TestDroppedContentDetectsMarkdownHardBreakLoss
      ... a lost Markdown hard break was not reported
  --- FAIL: TestDroppedContentDetectsReorderedContent
      ... reordering that changes meaning was not reported
  --- FAIL: TestDroppedContentForgivesExactlyOneStamp
      ... a second stamp-prefixed line was silently swallowed
  FAIL
  exit 1
  ```

- **[MAJOR codex-1] Stamp forgiveness was inherited by a following genuine line — PARTIALLY
  FIXED.** Real command-path probes against a freshly rendered deck establish the four requested
  combinations when the extra line follows the structural stamp:

  | Deck stamp | Extra genuine line | Result |
  | --- | --- | --- |
  | equal | absent | already matches; no removal |
  | equal | present | one header line reported |
  | different | absent | old stamp forgiven; no removal |
  | different | present | one header line reported |

  This fixes the exact round-7 reproducer. It is incomplete when the genuine line precedes the
  structural stamp; see the new MAJOR below.

- **[MAJOR codex-1] Whitespace-only lines were normalized away — FIXED in behavior, UNPINNED in
  tests.** Adding a spaces-only line to an otherwise converged generated deck now reports:

  ```text
  deck content NOT carried forward by core r8-probe:
    - ## 0. Choose the transport — 1 line not carried forward
  ```

  Reverting only `if keep[i]` to the old `if keep[i] || strings.TrimSpace(l) == ""` nevertheless
  leaves all 31 scoped tests green. The same prior finding also identified lone terminal-CR
  normalization; that part is not fixed (MINOR below).

- **[MAJOR codex-1] Full O(n*m) DP memory — FIXED in implementation, UNPINNED in tests.** The
  recursion is correct:

  - empty `a`/`b` returns without marking;
  - `len(a)==1` marks the one `a` position iff any `b` element matches;
  - `mid := len(a)/2` always makes progress;
  - `left[j] + right[len(b)-j]` is the correct prefix/suffix split score;
  - the first maximizing `j` is a valid deterministic tie choice; and
  - recursive bases `base` and `base+mid` address the two disjoint ranges of the original mask.

  The independent oracle output was:

  ```text
  randomized_trials=50000
  invalid_length_or_subsequence=0
  valid_lcs_masks=50000
  exact_mask_tie_variants_vs_prior_full_dp=10122
  first tie variant:
    a=[b d c b a b] b=[b]
    hirschberg=[false false false false false true]
    fullDP=[true false false false false false]
  recursion branches exercised:
    empty=81933 a1-hit=96958 a1-miss=9738
    split-0=62283 split-m=19650 split-interior=65243
  ```

  The exact-mask differences are not correctness failures: duplicated symbols admit multiple
  optimum LCS masks, and every Hirschberg mask had the oracle length and was a subsequence of `b`.
  But the new resource test does not reject the old implementation:

  ```text
  # render.go restored from 134c52a; current render_test.go retained
  $ go test -count=1 -v ./internal/protocolcore -run '^TestLargeDocumentDoesNotExhaustMemory$'
  === RUN   TestLargeDocumentDoesNotExhaustMemory
  --- PASS: TestLargeDocumentDoesNotExhaustMemory (1.21s)
  PASS
  ok  parley-deck-cli/internal/protocolcore  1.512s
  ```

- **[MAJOR carried G7b] Text still states guarantees no end-to-end test pins — NOT FIXED.** At
  `c4a8b83`, each of these independent scratch mutants left all 31 scoped tests green:

  ```text
  remove final-file O_NOFOLLOW                         PASS 31/31
  accept untrimmed/padded versions                    PASS 31/31
  remove CRLF-core normalization                      PASS 31/31
  remove existing-file-mode preservation              PASS 31/31
  replace fsutil.WriteFileAtomic with os.WriteFile    PASS 31/31
  ```

  Yet `IMPLEMENTATION.md` states all five behaviors as landed (lines 21–24, 37–38, 138–150), and
  the changelog states `O_EXCL|O_NOFOLLOW` as a shipped guarantee (lines 23–24). This is the G7b
  mismatch requested in verification item 5.

- **[MINOR carried] “No agent-accessible write path” — NOT FIXED and disproved.** The identical
  §7 block in the live deck, embedded default, and skill source has SHA-256
  `797f8d5250bc7de9d0b2fe7d4b5f4280a964ade151db54b2f4f55ba91e29a679` and ends with “no
  agent-accessible code path writes a release.” The changelog says the same, and
  `IMPLEMENTATION.md` lines 71–75 repeat it after conceding the PTY route. Driving the real CLI
  from this agent with `tty:true` published successfully:

  ```text
  $ PARLEY_HOME=/private/tmp/... parley protocol publish \
      --version r8-probe --from internal/protocol/defaults/COOPERATION.md
  Published core r8-probe (73efe37fd237) to /private/tmp/.../protocol/core/r8-probe
  exit 0
  ```

  The CLI refusal also calls “hash detection” a durable guarantee, while `Store.Load` explicitly
  verifies no expected release hash. Deck-drift detection is covered; release-tamper detection is
  not implemented, so the unqualified phrase remains misleading.

- **[MINOR carried] Phase-8 reviewed HEAD — NOT FIXED.** `IMPLEMENTATION.md` names `134c52a` and
  says cycle 7 lands in the next commit; the reviewed commit is `c4a8b83`. This remains one commit
  behind the Phase-8 record required by the protocol.

## New findings (by severity, or "none")

### [MAJOR] Header stamp pairing still forgives genuine prose when it comes first

`headerStamp` (`render.go:352-365`) returns the first prefix-matching line before the first `##`
heading. That is not the structural generated-stamp slot: the renderer creates that slot
immediately after `**Created:**`. If genuine stamp-prefixed project prose precedes an already
current generated stamp, `headerStamp(prior)` selects the prose, `headerStamp(rendered)` selects
the current stamp, and `droppedContent` forgives the prose while LCS keeps the current stamp.

Real command-path reproduction from a freshly converged deck:

```text
# Insert before the current stamp:
**Protocol synced:** genuine project prose
**Protocol synced:** core r8-probe (73efe37fd237)

$ parley protocol render --dir <fixture> --dry-run
preserved from this deck: Workspace, Transport, Created, §2 roster table, host-handle table
would regenerate <fixture>/parley-deck/COOPERATION.md from core r8-probe (73efe37fd237).
Nothing was written. Re-run with --yes to apply.
```

There is no `deck content NOT carried forward` block; `--yes` deletes the genuine line silently.
Pair the line at the actual structural slot (immediately after the `Created` identity line), not
the first prefix match, and cover genuine prefix lines on both sides of that slot for equal and
different stamps.

### [MAJOR] Three cycle-7 fixes have no deterministic reversion-sensitive test

The round-6 fixture repair is good, but cycle 7 repeats the same proof gap for its own changes:

1. Reverting only the header-pairing block to the `134c52a` first-unmatched-prefix logic leaves all
   31 scoped tests green. The rewritten `TestDroppedContentForgivesExactlyOneStamp` exercises only
   a changed old stamp followed by extra prose, not the equal-stamp case that motivated the fix.
2. Reverting only the whitespace skip leaves all 31 scoped tests green.
   `TestDroppedContentTreatsDeeperHeadingsAsSections` now comments that blank lines count but asserts
   only that the heading string appears.
3. Restoring the full matrix leaves `TestLargeDocumentDoesNotExhaustMemory` green on a host with
   enough memory. The test is an environmental stress probe, not a deterministic assertion of the
   algorithm's space bound; `TestLCSSeesOrder` passes under either LCS implementation as expected.

Add direct cases for equal stamp + extra line, whitespace-only removal, and an algorithmic
allocation/implementation guard that fails deterministically when the quadratic table returns.

### [MINOR] The exact-line diff still erases a lone terminal carriage return silently

`splitLines` says it strips only the CR of a CRLF pair, but line 270 applies
`strings.TrimSuffix(line, "\r")` to every split segment. A lone CR at end-of-file is therefore
normalized too. Replacing only the generated deck's final LF with a lone CR makes `render
--dry-run` announce regeneration but print no removal block:

```text
preserved from this deck: Workspace, Transport, Created, §2 roster table, host-handle table
would regenerate <fixture>/parley-deck/COOPERATION.md from core r8-probe (73efe37fd237).
Nothing was written. Re-run with --yes to apply.
```

The public `Render` path already normalizes actual `\r\n` pairs before calling `droppedContent`,
so `splitLines` should not strip an unpaired final CR. Add an end-to-end terminal-CR case or narrow
the exact-line claim explicitly.

## Test-quality assessment

The three purpose-built round-6 fixtures are now high quality: minimal, transformation-isolated,
and demonstrated to fail under the defeated multiset. The Hirschberg code is correct beyond the
small checked-in examples; the randomized full-DP comparison exercised empty inputs, both
single-element outcomes, endpoint and interior splits, duplicate-heavy ties, and all recursion
offsets without an invalid optimum.

Cycle 7's own regression protection is inadequate. The stamp and whitespace changes have no test
that notices their removal, and the linear-space claim is guarded only by whether the current
machine can survive a 1.6-GB allocation—which this machine can. The five surviving G7b mutants
also mean the three protocol copies, changelog, CLI output, and `IMPLEMENTATION.md` continue to
outrun the real-entry-point suite. No rounds 3–6 scoped regression failed at the target commit;
the only repository-wide test failure is in unchanged runner code.
