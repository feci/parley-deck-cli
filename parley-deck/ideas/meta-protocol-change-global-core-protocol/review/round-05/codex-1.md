---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 5
date: 2026-08-07
reviewed-commit: 1903e49
verdict: FINDINGS
---

## Summary

Cycle 4's four regression-test claims are real: each targeted test fails for the exact reverted
behavior, including the rewritten cross-section fixture. The cumulative implementation is still not
CLEAN. `droppedContent` has three reproducible silent-loss shapes outside that fixture; the corrected
third protocol copy is still only an uncommitted sibling-repository diff; several guarantees still
survive with all 23 protocol tests after their implementation is removed; the text still says there
is no agent-accessible release writer while admitting that a pty-allocating agent can invoke it; and
my round-02 Phase-8 metadata finding remains open.

The required commands produced:

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
?   parley-deck-cli/cmd/parley [no test files]
ok  parley-deck-cli/internal/app 55.845s
...
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got
    {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001
    Message:process verification failed (no recorded boot id); not killed}
FAIL parley-deck-cli/internal/runner 8.172s
FAIL
exit 1
```

The runner failure also reproduced alone. `02bb1e3..1903e49` has no `internal/runner` diff, and the
same failure is recorded at earlier reviewed commits, so I do not attribute it to this cycle. The
protocol-focused run is clean and the documented count now reproduces:

```text
$ go test ./internal/app -run '^(TestProtocol|TestPublish|TestCore|TestLoad)' -count=1 -v
23 tests, 23 PASS
ok  parley-deck-cli/internal/app 0.442s
```

## Prior findings: fixed or not

- **Cycle-4 cross-section test — FIXED.** With only per-section matching reverted to the old global
  multiset, `TestProtocolRenderDetectsCrossSectionErasure` fails and says the §2 loss was masked by
  the identical line under §3. This is now a genuine distinguishing fixture.
- **Load release-directory symlink test — FIXED.** Removing only `Load`'s release-directory
  `Lstat` makes `TestLoadRefusesASymlinkedReleaseDirectory` fail after reading `"PLANTED"`.
- **Publish store-component symlink test — FIXED.** Removing only Publish's
  `assertNoSymlinkComponents` call makes `TestPublishRefusesASymlinkedStoreComponent` fail with
  `published through a symlinked store component`.
- **Render read-error test — FIXED.** Restoring `prior, _ := os.ReadFile(path)` makes
  `TestProtocolRenderReportsDeckReadErrors` fail because dry-run reports a successful regeneration
  and no read error. The `--dry-run` shape avoids the earlier wrong-reason write failure.
- **Round-01/02 operational fixes — MOSTLY FIXED.** Production dispatch, Windows/Linux builds,
  existing-release refusal, path validation, Load/Publish store-component refusal, release-dir
  refusal, status/render read errors, TTY refusal wording, and the explicit rank-2/rank-4 deferral
  are present and the covered cases behave as documented.
- **Round-02 `droppedContent` finding — NOT fully fixed.** The exact cross-heading duplicate case is
  now caught, but lossy line normalization and incomplete heading indexing still permit silent
  loss; see the first MAJOR finding.
- **Round-02 Phase-8 record finding — NOT FIXED.** `IMPLEMENTATION.md` still says
  `status: implemented` and `head-commit: (see release commit)` after four fix-up cycles.
- **Agent-proof publishing wording — NOT FIXED.** The adjacent pty caveat is honest, but the
  surviving absolute statements remain false; see the MINOR finding.
- **Core-side section rename over-report — VERIFIED and accepted as recorded.** A fresh old render
  upgraded across `## 3. Old phases` -> `## 3. Phases` reports
  `## 3. Old phases — 2 lines not carried forward` while `Core rules live here.` survives in the
  new body. I do not elevate this acknowledged fail-loud trade-off to a new finding.

## New findings (by severity, or "none")

### [MAJOR] `droppedContent` still silently loses Markdown content

The current matcher is not a lossless per-section comparison. It trims every line on both sides
(`internal/protocolcore/render.go:215-225,249-265`), treats every line beginning with
`**Protocol synced:**` as generated regardless of location (`render.go:220-223`), and recognizes
only `##` and `###` headings (`render.go:269-274`). Three scratch tests against `1903e49` produced:

```text
TestProbeMarkdownIndentationLoss:
  Markdown-significant indentation was lost without a §3 report
  report named only §99

TestProbeSyncedPrefixOutsideHeader:
  project prose beginning with the stamp prefix was lost without a §3 report
  report named only §99

TestProbeLevelFourSubsectionLoss:
  got:  ### Parent — 1 line not carried forward
  want: ### Parent — 2 lines not carried forward
```

The first fixture puts `    Core rules live here.` (a Markdown code block) in the deck while the
core has unindented `Core rules live here.` prose. `strings.TrimSpace` calls them identical, so
preview omits the §3 loss and apply changes the Markdown meaning silently. The second puts genuine
project prose beginning with `**Protocol synced:**` under §3; the unconditional prefix skip omits
it from the report. The third moves `Shared rule.` between two different `####` subsections under
the same `### Parent`; because level-four headings are not section boundaries, only the old
subheading is counted and the moved rule is falsely treated as carried in place.

This violates G1 independently of the repaired §2/§3 fixture. Compare exact Markdown structure (or
the ratified stable block IDs), preserve Markdown-significant whitespace, use full heading paths and
occurrences at every supported level, and skip the generated stamp only at its exact header slot and
shape. Add end-to-end preview/apply tests for all three cases.

### [MAJOR] The corrected third protocol copy is not committed or installed

`IMPLEMENTATION.md:43-47` says the §7 clause was added to all three copies. The live deck and
embedded default carry the corrected rank-2/rank-4 disclaimer. The sibling skill repository does
not carry that correction in its HEAD:

```text
$ git -C ../parley-deck-skill status --short
 M skills/parley-deck/references/COOPERATION.md

$ git -C ../parley-deck-skill rev-parse HEAD
455aafe9f99fd6c01223b920a0768af2119e14a3

$ git -C ../parley-deck-skill show HEAD:skills/parley-deck/references/COOPERATION.md
...
An idea that is already open completes under the protocol version it was pinned to; the next idea
in that deck picks up the current one.
```

The honest disclaimer exists only in that repository's uncommitted working-tree diff. The installed
Codex 2.5.1 fallback does not contain the blast-radius clause at all. A commit/release from the
recorded sibling HEAD therefore restores the false rank-2 guarantee. Commit and release the third
copy correction, and add a release-time drift guard that actually includes it; the in-repo drift
test compares only the live and embedded copies.

### [MAJOR] G7b is still violated by documented guarantees that no protocol test guards

I independently removed each claimed fix in a fresh `1903e49` scratch copy and ran all 23 protocol
tests. Every mutation below compiled and returned `ok parley-deck-cli/internal/app`:

| Removed behavior | Text that claims it | Result |
| --- | --- | --- |
| `O_NOFOLLOW` from Publish's open flags | changelog:23; `IMPLEMENTATION.md:21-24,73-75` | 23/23 PASS |
| rejection of padded/untrimmed versions | `IMPLEMENTATION.md:138-139` | 23/23 PASS |
| CRLF-core normalization | `IMPLEMENTATION.md:148-150` | 23/23 PASS |
| corrected `protocol check` preamble | `IMPLEMENTATION.md:187-189` | 23/23 PASS |
| render's existing-file mode preservation | `IMPLEMENTATION.md:37-40` | 23/23 PASS |

There is also no end-to-end test of the atomic-write claim, and no test reaches the third skill
copy. The existing symlink test still short-circuits on the pre-existing release directory before
the final-file `O_NOFOLLOW` flag matters. `TestPublishRejectsUnsafeVersions` still has no padded
input. `TestProtocolRenderHandlesCRLFDecks` still supplies a CRLF deck but an LF core. No test
asserts the check preamble or file mode.

These are not requests for exhaustive testing of implementation trivia: the implementation and
changelog explicitly present them as shipped guarantees, while G7b says no such guarantee may be
documented without a real-entry end-to-end test. Either add distinguishing production-boundary
tests or narrow the claims. Also change `IMPLEMENTATION.md:69`: the cross-builds were run manually
in this review, but no repository check or CI configuration makes them "part of the check."

### [MINOR] The user-only/no-agent-writer guarantee still contradicts the documented pty behavior

All three corrected §7 working copies say both that an agent may not publish and that the TTY gate
does not stop an agent that allocates a pty, then conclude that "no agent-accessible code path
writes a release" (`parley-deck/COOPERATION.md:758-775`; same semantic text in the other two).
The changelog repeats "no agent-accessible write path" (`meta/protocol-changelog.md:23-27`), and
`IMPLEMENTATION.md:71-75` repeats it immediately after acknowledging the stock pty bypass. The CLI
refusal similarly says "Only the user may change" before its pty caveat and then advertises "hash
detection" as a durable guarantee even though release-tamper detection is explicitly rank 4
(`internal/app/protocol.go:296-301`).

The caveat makes the practical limit discoverable, but it does not make the adjacent absolutes
true. State only the property implemented and tested: a non-pty process is refused; a pty-capable
process is not prevented, and release-tamper attribution is not shipped.

### [MINOR] The Phase-8 implementation record still identifies neither the phase nor reviewed HEAD

`IMPLEMENTATION.md:1-10` remains `status: implemented` with `head-commit: (see release commit)`.
Phase 8 requires the top-level status and SHA to advance every cycle. This is my round-02 finding,
not a new mechanism, and it disproves cycle 4's statement that all codex-1 round-01/02 findings were
addressed. Set the current fix-up status and exact reviewed commit in the machine-readable
frontmatter; do not rely on four appended prose sections to recover current state.

## Test-quality assessment

Cycle 4 repaired exactly the four tests it claims to repair. All four fail for their intended
reverted behavior; none passes for a setup, later-write, or missing-fixture reason:

| Test | Surgical reversion | Observed failure |
| --- | --- | --- |
| `TestProtocolRenderDetectsCrossSectionErasure` | per-section -> global counts | §2 loss masked by §3 copy |
| `TestLoadRefusesASymlinkedReleaseDirectory` | remove release-dir `Lstat` | read `"PLANTED"` through link |
| `TestPublishRefusesASymlinkedStoreComponent` | remove Publish component check | publish succeeded through link |
| `TestProtocolRenderReportsDeckReadErrors` | swallow read error | dry-run exited successfully and reported no read error |

The protocol count is genuinely 23. The broader suite is not green on this machine because of the
unchanged runner boot-id failure quoted above. Test quality is therefore improved materially over
round 4, but the suite still does not enforce several statements the implementation labels as
guarantees, and G1 still has untested, reproduced silent-loss cases.
