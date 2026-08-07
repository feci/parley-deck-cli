---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 4
date: 2026-08-07
reviewed-commit: 02bb1e3
verdict: FINDINGS
---

## Summary

Cycle 3 is the first cycle whose CODE does what its documentation says — I verified every claimed
fix behaviorally against the real binary built from `git archive 02bb1e3` in a scratch tree. The
per-section `droppedContent` genuinely is per-section now: my round-03 erasure fixture (deck
carrying `- Deploy production.` ONLY under `## Requires explicit user approval`, core carrying it
under `## Allowed automatically`) is now reported in preview AND apply
(`## Requires explicit user approval — 1 line not carried forward`), the apply moves the line under
the core's section, and a following `check` converges (exit 0). Multiplicity is exact (3 identical
local lines → "3 lines"; with 1 same-section copy in the render → "2 lines"). `Load` refuses both a
symlinked store component and a symlinked release directory, deck unpolluted both times. The
changelog is honest about `DETECTED-UNATTRIBUTED`. The `check` preamble no longer contradicts its
entries. `render` on an unreadable deck exits 1 and leaves the file untouched. The test count
reproduces: 20.

And yet the cycle ships the same disease in its purest form yet: **the headline fix is pinned by a
test that is provably tautological, and IMPLEMENTATION.md:175 asserts "A new test constructs
exactly the cross-section case" — it does not.** I reverted `droppedContent` to the cycle-2 global
multiset in the scratch tree and the ENTIRE suite stayed green: 20 PASS, 0 FAIL, including
`TestProtocolRenderDetectsCrossSectionErasure`. Three cycles of documentation outrunning reality on
this one function; this time the code is right and the test is the lie. If the per-section logic
regresses tomorrow, nothing notices.

Command output (this machine, 02bb1e3, clean tree):

```
$ go build ./...                  EXIT=0
$ GOOS=windows go build ./...     EXIT=0
$ GOOS=linux go build ./...       EXIT=0
$ go vet ./...                    EXIT=0
$ go test ./... -count=1          25 packages ok, 2 no-test-files (incl. internal/app 46.081s,
                                  internal/runner 11.120s), EXIT=0
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore|TestLoad' -count=1 -v
  20 tests, 20 PASS, 0 FAIL (--- PASS lines counted: 20)
```

Revert checks below ran in `/tmp/pd-r4/src` (`git archive 02bb1e3`), each revert verified landed
by grep/compile before trusting any result. The project repo was not modified.

Counts: 1 MAJOR, 1 MINOR, 1 NIT (plus one prior finding downgraded to a NIT residual).

## Prior findings: fixed or not

1. **[MAJOR hermes-1 + kimi-1] cycle-2 "per-section" was global — FIXED IN CODE, PINNED BY ZERO
   TESTS.** The fix is real this time: `indexBySection` (render.go:250-267) keys the render's
   counts by section and the deck walk (render.go:215-228) consumes only within the same section;
   `heading()` is applied symmetrically on both sides, so the keys cannot drift. Verified
   behaviorally: the round-03 cross-section fixture reports in preview and apply; multiplicity
   3→3 and 3-with-1-carried→2; a fully-carried deck renders "already matches" with no false
   report; idempotent re-render; `check` exit 0. I could not construct a residual MISS for the
   demonstrated erasure classes. **But the test written to pin it cannot catch a revert** — new
   finding 1 — and one over-report class exists — new finding 3.

2. **[MAJOR hermes-1] `assertNoSymlinkComponents` write-only — FIXED, both halves verified through
   the binary.** Load path: `<home>/protocol` as a symlink → `protocol render` exits 1
   ("protocolcore: /tmp/pd-r4/home2/protocol is a symlink; the core store must not be reached
   through one"), deck untouched; `<store>/safe-1` itself a symlink (my exact round-03 repro) →
   exit 1 ("release directory … is a symlink"), deck untouched. Publish path: symlinked
   `<home>/protocol` → refused, outside store gained nothing; a pre-planted symlink at
   `<store>/8.8.8` → refused ("already exists and releases are write-once"), the outside directory
   stayed empty. **Residual, probed and judged NOT a finding:** a symlink at the release FILE
   (`<store>/<v>/COOPERATION.md` → outside) is still followed by `Load` — I watched the outside
   bytes land in the deck, exit 0. No text claims file-level refusal (the claims name store
   components and the release directory), and the capability is a strict subset of what the same
   local attacker already has (the release dir is 0755; `rm` + rewrite the 0444 file; the lock
   pins no hash — the acknowledged rank-4 gap). Flagging it so the next cycle doesn't rediscover
   it. The release-dir half of the Load fix has no test — new finding 2.

3. **[MAJOR hermes-1 + kimi-1] changelog cited `DETECTED-UNATTRIBUTED` as shipped — FIXED.**
   protocol-changelog.md:23-28 now lists what ships and marks the tamper signal and pinning
   "ratified but NOT implemented … ranks 2 and 4", with the pty limitation stated. Verified:
   `DETECTED-UNATTRIBUTED` still has 0 Go references, and the changelog no longer says otherwise.

4. **[MINOR kimi-1] §7 "in force" sentence overclaimed — PARTIALLY FIXED; the absolute clause I
   named survives.** The pty qualifier is now inline ("stops an ordinary agent run, not one that
   allocates a pty") — but the same sentence still ends "…and **no agent-accessible code path
   writes a release**" (parley-deck/COOPERATION.md:772-774; embedded and skill copies
   byte-identical, diff-verified deck==embedded==skill-wt). I published through `pty.spawn` again
   this cycle (exit 0): for a pty-allocating agent, `parley protocol publish` IS an
   agent-accessible code path that writes a release — the sentence's own parenthetical says so two
   commas earlier. My round-03 finding suggested the fix ("no code path an ordinary agent run can
   reach…"); the cycle qualified the refusal clause and left the absolute. Residual is wording, not
   substance — all three copies carry the admission inline now — so NIT, not MINOR.

5. **[MINOR kimi-1] render swallowed deck read errors — FIXED behaviorally, unpinned.** chmod 000
   deck → `protocol render: open …: permission denied`, exit 1, file byte-untouched (verified).
   `check` already erred; still does. No test — new finding 2.

6. **[MINOR hermes-1] check preamble contradicted its entries — FIXED.** Verified against the
   shared-heading fixture: `check` now prints "deck content not carried forward by this core:"
   above `## Requires explicit user approval — 1 line not carried forward`, exit 1. No falsehood.

7. **[NIT kimi-1] test count did not reproduce — FIXED.** 20 `--- PASS` lines from an actual run
   (quoted above); `grep -c '^func Test' internal/app/protocol_test.go` = 20. Cycle 1's "Fifteen"
   and cycle 2's "Nineteen" were also corrected in place.

**Still open from earlier rounds, not claimed by this fix-up** (restating so they are not lost):
the O_NOFOLLOW half of `TestPublishRefusesExistingReleaseDirAndSymlinks` remains unreachable dead
code (pre-created dir fires first; round-02 finding, untouched); the refusal text's "hash
detection" (protocol.go:301) still overstates what exists (round-03 NIT); a deck line of genuine
prose starting with `**Protocol synced:**` is still silently skipped (render.go:220, round-03 NIT);
`docs/cli-reference.md` still never documents the `protocol` command (its one "protocol" hit is
the word in an `init` description) and the root `CHANGELOG.md` has no
`protocol publish/render/check/status` entry; preflight.go:584-596 still owns the legacy stamp
format (latent until adoption); the mode-preserving write branch and `check --json`'s non-zero
exit remain untested; IMPLEMENTATION.md frontmatter still says `head-commit: (see release
commit)`; IMPLEMENTATION.md:69 still claims the cross-builds "are now part of the check" — there
is still no check (`grep -rn GOOS scripts/ .github/` → 0).

## New findings (by severity)

### [MAJOR] The cross-section test is a tautology; the cycle's headline fix is pinned by zero tests

`TestProtocolRenderDetectsCrossSectionErasure` (protocol_test.go:505-521) constructs NO
cross-section case. The fixture adds "Core rules live here." to the deck's `## 3. Phases` section
— and the render's copy of that line is ALSO under `## 3. Phases` (it comes from the core). A
same-section coincidence is handled identically by the per-section index and by the reverted
global multiset; the only line the assertion can ever catch is the STALE line, which both
implementations report. The test's own comment half-notices this ("in the deck it sat under
'## 3. Phases' too, so this must NOT be reported… The erasure to catch is the deck's OTHER line")
— but the OTHER line is caught by the reverted code too, so the test distinguishes nothing.

Demonstrated, not argued: in the scratch tree I replaced `internal/protocolcore/render.go` with
its 4a5c447 version (revert verified landed: `grep -c indexBySection` → 0) and ran the suite:

```
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore|TestLoad' -count=1 -v
--- PASS lines: 20    --- FAIL lines: 0
--- PASS: TestProtocolRenderDetectsCrossSectionErasure (0.00s)
```

Every test green with the fix fully reverted. So IMPLEMENTATION.md:175's "A new test constructs
exactly the cross-section case" is false, and the fix — correct as it is — could be deleted
tomorrow with the suite staying green. This is the third consecutive cycle in which
`droppedContent`'s documentation outran its reality (cycle 2: documented per-section, code
global; cycle 3: documented tested, test tautological), and it is exactly the class this re-review
exists to catch. A real test: put the deck's copy of the coincidental line under a section the
render does NOT carry it in (e.g. `## 99. Project-specific rule`), drop the STALE line so it
cannot satisfy the assertion, and assert the report names that section with the right count —
that fails under the global multiset and passes now (I confirmed the underlying behavior via the
binary: the round-03 fixture reports correctly).

### [MINOR] Three more claimed fixes have no test that notices their absence — each revert demonstrated green

Same revert method, one claim at a time, full 20-test protocol suite each run:

- **Load's release-directory symlink refusal** (core.go:96-99; fix-up bullet claims it alongside
  the store-component half): removing the `Lstat(dir)` block → 20/20 green. The behavior is real
  (binary probe above); only `TestLoadRefusesASymlinkedStore` exists and it covers the
  store-component half — which IS real: removing the `assertNoSymlinkComponents(s.Root)` call from
  Load fails it ("Load read through a symlinked store component").
- **Publish's store-component symlink refusal** (core.go:148): removing the call → 20/20 green.
  The changelog now asserts "symlinked store components refused on read and write" — the read half
  is pinned, the write half is not (carried from cycle 2, where I noted the same).
- **Render's read-error report-and-exit** (protocol.go:180-186; fix-up claims "Now reports and
  exits non-zero"): restoring `prior, _ := os.ReadFile(path)` → 20/20 green. A `chmod 000` fixture
  test is trivially writable — the same shape as the existing status read-error test.

Per the project's own G7b posture (cycle 2 disclosed the untested publish-cleanup rather than
claiming it — the right move), these three should have had tests or the same disclosure. Bundled
as one MINOR because the pattern is one pattern.

### [NIT] A core-side section rename reports a surviving line as "not carried forward"

Version bump that renames a section but keeps its body lines verbatim (core 3.0.0 `## Phases` /
`- Do the thing.` → 3.1.0 `## Stages` / `- Do the thing.`), deck freshly rendered from 3.0.0:

```
deck content NOT carried forward by core 3.1.0:
  - ## Phases — 2 lines not carried forward
```

One of those 2 lines (`- Do the thing.`) IS in the output, under `## Stages`. This is the direct
cost of the per-section strictness the reviewers demanded, and it errs on the loud side — nothing
is silent, the section genuinely is removed, and when the core MOVES a line between
differently-governed sections (my other probe: approval → allowed) the report is exactly what G1
wants. I judge it a report-fidelity nit, not a defect: the count conflates "the section is gone"
with "these lines vanished", and on a pure rename a user finds the "lost" lines alive in the new
section. If the report ever gets a polish pass, "section removed; N of its M lines survive
elsewhere" would be the honest shape. Noting it because the brief asked me to construct precisely
this input; it is not a reason to loosen the section matching.

## Test-quality assessment

20 protocol tests (+1 drift test in internal/protocol), all green on `-count=1`. Revert method as
before: scratch tree from `git archive 02bb1e3`, surgical reverts, each revert verified landed
before trusting a PASS.

New tests this cycle:

- `TestProtocolRenderDetectsCrossSectionErasure` (:505) — **TAUTOLOGY**, proven above: passes with
  the per-section fix fully reverted; its fixture's coincidental match is same-section, and its
  assertion is satisfied by the STALE-line report under both implementations. It adds zero
  distinguishing power over `TestProtocolRenderReportsContentLostUnderASharedHeading`.
- `TestLoadRefusesASymlinkedStore` (:525) — **REAL.** Removing Load's `assertNoSymlinkComponents`
  call fails it ("Load read through a symlinked store component"). It pins only the
  store-component half of the claimed Load fix.

Changed claims with NO test that notices (each revert left all 20 green): the per-section rewrite
itself (MAJOR above); Load's release-dir refusal; Publish's store-component refusal; render's
read-error exit (MINOR above). Carried, unchanged: the dead O_NOFOLLOW half of
`TestPublishRefusesExistingReleaseDirAndSymlinks` (round-02).

**G7b sweep — any text claiming a guarantee with no end-to-end test:**

- IMPLEMENTATION.md:175 "A new test constructs exactly the cross-section case" — **false**
  (demonstrated). The one new false claim this cycle, and a MAJOR one because it guards the
  headline fix.
- Changelog "symlinked store components refused on read and write" — read pinned, **write
  unpinned**; "O_EXCL|O_NOFOLLOW" — O_NOFOLLOW still dead in its test. Behavior of all four
  verified through the binary this cycle.
- §7 "no agent-accessible code path writes a release" — still absolute, still contradicted by the
  same sentence's pty admission (all three copies; byte-identical).
- IMPLEMENTATION.md "Now reports and exits non-zero" (render read errors) — true, verified,
  untested.
- IMPLEMENTATION.md:69 "GOOS=windows and GOOS=linux builds are now part of the check" — still no
  check exists.
- protocol.go:301 "hash detection" — still one word past honest (round-03 NIT, untouched).
- Honest-and-pinned this cycle: write-once (tested), attended refusal incl. its own limits
  (tested), Load store-component refusal (tested), status read errors (tested), traversal guard
  (tested, both layers), the 20-test count (reproduces), the changelog's ratified-not-implemented
  hedges (match the code: 0 Go references).

The trajectory is genuinely improving — this is the first cycle where every code claim I probed
behaved as documented. What remains is the same lesson in its third form: a fix is not done when
the code is right; it is done when a test can prove it was wrong before. The suite currently
cannot prove that for the fix this cycle exists to ship.
