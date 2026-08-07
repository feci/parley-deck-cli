---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 7
date: 2026-08-07
reviewed-commit: 134c52a
verdict: FINDINGS
---
## Summary

The LCS itself is a correct classic dynamic-programming implementation, and it fixes the four
minimal transformations from round 6: hard-break loss, HTML-comment reordering, movement between
duplicate-heading contexts, and a changed generated stamp followed by a second stamp-prefixed
project line. In an isolated `134c52a` archive, all four minimal probes passed with the LCS and all
four failed after replacing only `render.go` with the `8ed3c4b` multiset. A removed-heading probe
also produced the correct attribution, `## Removed section — 2 lines not carried forward`.
Order-preserving carried content is not reported; moved content is deliberately represented as a
removal plus an addition, as the narrowed contract states.

The implementation is nevertheless not CLEAN. The exactly-one-stamp exemption misfires once the
generated stamp already matches, blank/whitespace-only line removals bypass the report, the
quadratic LCS table has unbounded memory growth, and the claimed permanent regression coverage is
not real. All three checked-in round-6 tests still PASS with the old multiset, while the fourth
(duplicate-heading) test does not exist. Carried G7b/documentation findings and the one-commit-late
Phase-8 `head-commit` also remain open.

Required commands at `134c52a7de645bf25be73cd23ec78d34c9b7ee47`:

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
ok      parley-deck-cli/internal/app       40.757s
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL    parley-deck-cli/internal/runner     7.579s
FAIL
exit 1
```

The runner failure reproduces alone with the same message. `8ed3c4b..134c52a` contains no runner
change, and round 6 recorded the same failure, so I do not attribute it to cycle 6. The protocol
suite is green and its claimed count reproduces:

```text
$ go test ./internal/app -run '^(TestProtocol|TestPublish|TestCore|TestLoad|TestDroppedContent)' -count=1 -v
...
PASS
ok      parley-deck-cli/internal/app       0.403s
PROTOCOL_TEST_COUNT runs=29 pass=29 fail=0
```

## Prior findings: fixed or not

- **Round-6 `droppedContent` behavior — FIXED for the four reported transformations.** Four
  isolated direct probes pass with the current LCS and fail with the old per-section multiset.
  The checked-in tests do not establish this; that is a new MAJOR below.
- **Round-5/6 G7b end-to-end coverage — NOT FIXED [MAJOR].** The cycle added only three
  `droppedContent` tests. In six isolated `134c52a` mutants, all 29 protocol tests still passed
  after independently removing final-file `O_NOFOLLOW`, accepting untrimmed versions, removing
  CRLF-core normalization, restoring the false `check` preamble, forcing render mode to `0644`,
  or replacing the atomic render write with `os.WriteFile`.
- **User-only/no-agent-write wording — NOT FIXED [MINOR].** All three protocol copies still end
  the qualified pty discussion with the absolute “no agent-accessible code path writes a
  release.” The changelog says “no agent-accessible write path,” and CLI refusal output still
  calls “hash detection” a durable guarantee even though `Load` explicitly performs no expected-
  hash verification. The pty path and exported `Publish` make the absolute false, not merely
  untested.
- **Phase-8 reviewed HEAD — NOT FIXED [MINOR].** `IMPLEMENTATION.md` now says
  `head-commit: 8ed3c4b (cycle 5); cycle 6 lands in the next commit`; this review is of
  `134c52a`. The protocol requires the top-level field to name the new fix-up HEAD.
- **Rounds 3–5 CLEAN behavior — NOT REGRESSED.** The 29-test baseline keeps the cross-section,
  indentation, non-header stamp-prefix, deeper-heading, read-error, path/version, symlink,
  dispatch, Windows, and Linux regressions green. Cycle 6 changed no production file outside
  `internal/protocolcore/render.go`; the previously recorded core-side rename remains loud, and
  removed-heading attribution is correct under the new sequence matcher.

## New findings (by severity, or "none")

### [MAJOR] Exactly one generated stamp is not exactly one consumed stamp

`droppedContent` consumes the exemption only when it encounters an **unmatched** prior stamp
(`render.go:231-243`). If a deck already carries the exact generated stamp, LCS marks that line
kept, so the exemption remains unused. A following genuine project line beginning
`**Protocol synced:**` is then the first unmatched stamp-prefixed line and is incorrectly
forgiven.

An end-to-end scratch probe first rendered a clean deck, inserted one genuine stamp-prefixed line
immediately after the now-current generated stamp, and ran `render --dry-run`. The file would be
regenerated, but there was no removal block:

```text
preserved from this deck: Workspace, Transport, Created, §2 roster table
would regenerate .../parley-deck/COOPERATION.md from core 1.0.0 (d9181eb71672).
Nothing was written. Re-run with --yes to apply.
```

Applying would silently delete the project line. Identify the one structural generated-stamp slot
in the document header and consume that allowance whether LCS matched it or classified an old
stamp as removed; every other prefix-matching line must be reported.

### [MAJOR] The “exact lines” report still silently normalizes whitespace-only lines

After computing the exact-line LCS, `render.go:236` discards every removed line for which
`strings.TrimSpace(l) == ""`. A blank line can separate Markdown paragraphs, lists, code blocks,
or HTML blocks, so deleting it is a real regeneration change under the deliberately
syntax-agnostic exact-line contract. A command-path probe added one extra blank line to an
otherwise freshly rendered deck; `render --dry-run` printed the same no-removal output quoted
above and would delete it on apply.

There is a second normalization beyond the documented CRLF exception: `splitLines` uses
`TrimSuffix(..., "\r")` on every split line (`render.go:265-274`), so a lone terminal carriage
return is also equated with no carriage return even though it was not part of CRLF. Remove the
blank/whitespace-only skip from reporting and strip carriage returns only as part of explicit
CRLF normalization.

### [MAJOR] None of the checked-in round-6 regression tests distinguishes LCS from the multiset

Only three tests were added at `protocol_test.go:660-699`; there is no duplicate-heading test,
despite `IMPLEMENTATION.md:294-295` saying all four cases are permanent tests. More importantly,
after replacing current `render.go` with the exact `8ed3c4b` multiset, all three checked-in tests
still passed:

```text
=== RUN   TestDroppedContentDetectsMarkdownHardBreakLoss
--- PASS: TestDroppedContentDetectsMarkdownHardBreakLoss (0.00s)
=== RUN   TestDroppedContentDetectsReorderedContent
--- PASS: TestDroppedContentDetectsReorderedContent (0.00s)
=== RUN   TestDroppedContentForgivesExactlyOneStamp
--- PASS: TestDroppedContentForgivesExactlyOneStamp (0.00s)
PASS
```

The hard-break fixture also adds `Next line.`, which the multiset legitimately reports. The
HTML fixture adds `<!--` and `-->`, which the multiset also reports. The stamp assertion checks
only the generic phrase `not carried forward`, while the shared fixture already loses unrelated
stale and project-section content. Isolate each transformation so the old/new documents have the
same unrelated content and assert the exact affected section/count. Add the missing duplicate-
heading test. The four isolated probes used in this review already demonstrate the required shape:
all pass with LCS and all fail with the multiset.

### [MAJOR] The full LCS matrix has pathological unbounded memory growth

`lcsMask` allocates `(n+1) × (m+1)` `int32` cells plus one slice allocation per row
(`render.go:280-301`) before it can report anything. No file-size or line-count bound enforces the
comment's “~1-2k lines” assumption. One-shot benchmarks on this host measured:

```text
BenchmarkRound7LCS/lines-1369-12    4.861542 ms/op     8,459,648 B/op    1,372 allocs/op
BenchmarkRound7LCS/lines-5000-12   60.945792 ms/op   102,548,848 B/op    5,004 allocs/op
```

At 20,000 × 20,000 lines the cells alone require about 1.6 GB; at 50,000 they require about
10 GB. A large or adversarial committed protocol can therefore OOM `protocol check`/`render`.
Use a linear-space reconstruction (for example Hirschberg), a Myers-style sequence diff, or an
explicit defensible input bound; common-prefix/suffix trimming is also a useful fast path but does
not by itself remove the worst case.

## Test-quality assessment

The LCS recurrence and backtracking are sound for exact line sequences, and direct probes confirm
the intended order-sensitive behavior. There was no false “lost” report for an order-preserving
subsequence; moved lines are intentionally removals under the stated contract. Section attribution
is correct when the heading itself is removed.

The committed suite is not an adequate regression gate. Its total of 29 reproduces, but the three
new tests are vacuous against the exact implementation they claim to replace, the fourth case is
missing, neither newly reproduced silent-change edge is covered, and no resource-bound test or
algorithmic guard exists. The six surviving G7b mutants also show that the three protocol copies,
`meta/protocol-changelog.md`, CLI output, and `IMPLEMENTATION.md` still state guarantees the real
entry-point suite does not distinguish. A green protocol suite therefore cannot support CLEAN.
