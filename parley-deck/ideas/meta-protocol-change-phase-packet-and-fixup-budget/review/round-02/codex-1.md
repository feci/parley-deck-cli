---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 2
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 2

## Summary

**[PRIMARY — CONFIRMED]** Five of my six round-1 findings are fixed or accurately withdrawn; the
source locators and independent checks are recorded under Q1. The published-cycle finding is
not fixed; it has been relocated from a review-round counter to a fail-open Markdown prefix count.
That counter can be reset by making `IMPLEMENTATION.md` unavailable or editing a prior heading, and
it counts code-fenced, duplicated, re-numbered, and malformed headings as real cycles. The release
must stop until the driver obtains a validated, monotonic count and treats unavailable or
inconsistent durable state as an escalation.

## Q1 — disposition of all six round-1 findings

1. **CRITICAL — consensus BLOCK bypassed the cross-review cap: fixed.**
   **[PRIMARY — CONFIRMED]** `internal/driver/driver.go:64-73,143-147` stores the policy ceiling as
   `HardCrossReviewCap`; `internal/driver/consensus.go:90-101` checks it before `RunRound`. I ran
   `go test ./internal/driver -count=1 -run
   '^TestBlockedConsensusRespectsTheHardCrossReviewCap$' -v`; the test passed and asserts no runner
   call occurs beyond the cap. An isolated standard-track probe also produced
   `action=escalated`, `§4.0 cap=2`, and `runnerCalls=[]` at the blocked round-3 consensus.

2. **CRITICAL — stale packaged-skill manifest: fixed.**
   **[PRIMARY — CONFIRMED]** `../parley-deck-skill/skills/parley-deck/parley-addon.json:7-11`
   carries the current protocol/compatibility hashes and aggregate
   `sha256:bee732321d0f5279c5ef83ae308a64e8533337e46be6cf185dae13680a1363db`.
   The independent package checks in Q6 pass.

3. **MAJOR — claimed real fast Phase-8 behavior from a synthetic driver path: fixed by withdrawing
   the claim.** **[PRIMARY — CONFIRMED]** `CHANGELOG.md:62-64` now states that `fast` forbids
   idea-level `auto_implement`, its fix-up route is manual, and no end-to-end fast claim is made.
   The rewritten test at `internal/driver/impl_test.go:448-476` is cap arithmetic over published
   cycle fixtures, not labeled as an actual fast-track integration test.

4. **MAJOR — review-round ordinals spent a published-fix-up budget: relocated, not fixed.**
   **[PRIMARY — WRONG]** `internal/driver/impl.go:290,490-501` now counts trimmed lines whose prefix
   is `## Fix-up cycle` and returns zero on every read error. It no longer spends budget on zero-fix
   review rounds, but it does not robustly count published cycles; see the CRITICAL finding below.

5. **MAJOR — ratified trajectory payload was claimed but absent: fixed by withdrawing the claim.**
   **[PRIMARY — CONFIRMED]** Both CHANGELOGs now say the structured payload is not implemented
   (`CHANGELOG.md:58-61`; sibling `CHANGELOG.md:16-18`). My exact search
   `rg -n 'fresh-vs-relitigated|findings by severity|trajectory payload|unresolved_fixes|validation_status'
   internal --glob '*.go'` returned no output, exit 1, consistent with that disclosure.

6. **NIT — comments described explicit deliberation as a no-op: fixed.**
   **[PRIMARY — CONFIRMED]** `internal/driver/driver.go:119-124` now reserves legacy behavior for an
   absent/unknown track, and `internal/track/track.go:115` says the same. The deliberation policy
   comment at `internal/track/track.go:151-159` explicitly separates the preserved full lifecycle
   from the two newly enforced cells.

## Q2 — `HardCrossReviewCap` and the standard back-edge

**Right for this release.** **[PRIMARY — CONFIRMED]** Standard's Phase-2 ceiling was separately
ratified before this idea: `parley-deck/ideas/meta-protocol-change-devx-speed/FINAL.md:38-46` says
standard is “capped at 2 rounds → escalate/upgrade,” and
`parley-deck/ideas/track-aware-driver/IMPLEMENTATION.md:77-79` plus that idea's
`review/consensus.md:23` record `CapCrossReviewRounds: 2` as an agreed fix. The prior implementation
clamped only the initially scheduled path; allowing the consensus-BLOCK back-edge to exceed the
same cap would repeat the round-1 defect under standard.

**[PRIMARY — CONFIRMED]** The current mechanism does change observable standard behavior: an
isolated `track: standard` setup with rounds 1–3 and a BLOCKED consensus constructed
`HardCrossReviewCap=2`, escalated before round 4, and made zero runner calls. That is enforcement of
the already-ratified cell, not a new standard policy number. The change is therefore in scope as a
correctness repair to a shared enforcement path.

## Q3 — robustness of `publishedFixupCycles()`

It is not robust enough to enforce a release gate.

### [CRITICAL] The published-cycle budget is fail-open and editable

**[PRIMARY — WRONG]** `internal/driver/impl.go:490-501` has four independent integrity failures:

- Missing or unreadable `IMPLEMENTATION.md`: `os.ReadFile` error returns `0`, so the next cycle is
  numbered 1 and the full budget is restored. The driver discards the error instead of escalating.
- Code fences: the function is line-prefix matching, not Markdown-aware. `## Fix-up cycle 99`
  inside a fenced example counts as a published cycle.
- Re-numbered, duplicated, or malformed headings: the suffix is never parsed. Two
  `## Fix-up cycle 9` headings plus `## Fix-up cycle banana` count as three cycles. Gaps and order
  are likewise unchecked.
- Budget extension by edit: changing a prior heading from `## Fix-up cycle 5` to
  `## Fixup cycle 5` changes the count from 5 to 4 and permits another cycle. A careless edit has
  the same effect as a malicious one: the next calculation is `4 + 1`, and `5 > cap(5)` is false.

**[PRIMARY — CONFIRMED]** I copied the current Go module to an isolated temporary directory, added
a same-package adversarial probe, and ran
`go test ./internal/driver -count=1 -run '^TestPublishedFixupCyclesAdversarialProbe$' -v`. Relevant
output:

```text
missing file -> 0
heading inside code fence -> 1
duplicate and malformed ordinals -> 3
one careless/malicious heading edit resets count 5 -> 4
unreadable file -> 0
PASS
```

**[PRIMARY — CONFIRMED]** The existing runner validation does not close this hole.
`internal/runner/phase58.go:159-183` accepts any review-ready frontmatter plus any occurrence of
`## Fix-up cycle`; it does not require a newly appended, exact, sequential heading. An isolated
probe with frontmatter `status: fix-up-cycle-2` but only an old `## Fix-up cycle 1` section passed
`ValidateFixupArtifact`:

```text
status says cycle 2, but validator accepted a document containing only cycle 1
PASS
```

Why this is CRITICAL: the ratified ceiling is a blocking safety threshold and extensions must be
finite, recorded, and must never reset (`FINAL.md:67-75`). An implementer controls
`IMPLEMENTATION.md`; with the current source it can reset or lower the enforcement count and keep
the auto-driven fix-up loop running past the cap. The round-1 problem has changed location, not
been eliminated.

Suggested fix: make the count API return `(int, error)` and escalate on missing/unreadable or
inconsistent state during Phase 8. Parse only exact, out-of-fence headings with positive decimal
ordinals; require a unique contiguous `1..N` sequence; cross-check frontmatter status and the newly
produced cycle in `ValidateFixupArtifact`. Prefer a driver-authored monotonic ledger/event or
marker whose update is part of the successful fix-up transaction; if headings remain the source,
add adversarial tests for every case above and for deletion/renaming after the cap is reached.

## Q4 — independent reversion check for `TestFixupCapIsInclusive`

**[PRIMARY — CONFIRMED]** The rewritten test genuinely detects loss of the published-cycle
derivation. In an isolated copy, I changed only
`cycle := publishedFixupCycles(d.cfg.IdeaDir) + 1` back to `cycle := round` and ran:

```text
$ go test ./internal/driver -count=1 -run '^TestFixupCapIsInclusive$' -v
--- FAIL: TestFixupCapIsInclusive
--- FAIL: .../cap_5:_the_6th_escalates
--- FAIL: .../cap_2:_the_3rd_escalates
--- FAIL: .../cap_1:_the_2nd_escalates
FAIL
exit 1
```

**[PRIMARY — CONFIRMED]** The current tree passes all six subtests. The three allowed-at-cap cases
happen to use review round 1 and therefore also pass under the reverted expression; all three
one-past-cap cases go red. This is a valid mutation check for the intended unit, even though the
source it tests is unsafe as implemented.

## Q5 — CHANGELOG accuracy

**[PRIMARY — CONFIRMED]** The `fast` and trajectory claims are now honestly withdrawn at
`CHANGELOG.md:54-64` and sibling `CHANGELOG.md:16-18`. Both CHANGELOGs still overstate enforcement
through the unsafe counter.

### [MAJOR] Both CHANGELOGs overstate enforcement by the unvalidated prefix count

**[PRIMARY — WRONG]** `CHANGELOG.md:14-17` says the count is derived from `## Fix-up cycle N`
records actually published. The implementation neither parses `N` nor establishes that a matching
line is a real published record, as the fenced/duplicate/malformed probe demonstrates. It also
silently resets to zero on read failure. `CHANGELOG.md:11-12`'s unconditional statement that
`MaxFixupCycles=N` publishes 1..N and escalates at N+1 is therefore false for corrupted or edited
durable state. The sibling `CHANGELOG.md:20-22` likewise says the printed and enforced numbers are
the same; that is not guaranteed while the enforcement count can reset or be lowered.

**[PRIMARY — CONFIRMED]** The remainder of the two release entries is appropriately scoped:

- `CHANGELOG.md:54-64` says the packet and payload are not shipped and fast is manual.
- The sibling `CHANGELOG.md:16-18` explicitly and correctly says the structured payload is absent.
- The three protocol copies differ only in their expected bootstrap headers outside the identical
  changed table rows; `go test ./internal/protocol -run '^TestEmbeddedDefaultMatchesLiveDeck$'
  -count=1` passed.

Suggested fix: after hardening the counter, describe the validated monotonic source. Until then,
replace “records actually published” with a narrower statement in the CLI entry and withdraw the
unconditional enforcement claim in both entries.

## Q6 — skill installation and test results

**[PRIMARY — CONFIRMED]** In `../parley-deck-skill` I ran the requested checks:

```text
$ node scripts/build-addon-manifest.js --check
parley-deck: ok (6 files, sha256:bee732321d0f5279c5ef83ae308a64e8533337e46be6cf185dae13680a1363db)
...all six packaged skills ok...
exit 0

$ npm test
tests 386; pass 386; fail 0
python 3.14: 54 tests OK across 7 files
...all six packaged manifests ok...
exit 0
```

The package now installs past the integrity gate that failed in round 1.

**[PRIMARY — CONFIRMED]** The relevant CLI packages are green:

```text
$ go test ./internal/driver ./internal/track ./internal/protocol -count=1
ok parley-deck-cli/internal/driver
ok parley-deck-cli/internal/track
ok parley-deck-cli/internal/protocol
```

**[PRIMARY — CONFIRMED]** The full CLI suite did not complete green in this sandbox: every listed
package passed except `internal/runner/TestDurableKillEndToEndRealProcess`; its output was
`process verification failed (no recorded boot id); not killed`. `git diff --quiet --
internal/runner internal/procctl` returned exit 0. The environment-specific failure is disclosed
and does not explain the counter defect above.

## Q7 — release gate

**Stop this release.** The CRITICAL counter defect permits the binding Phase-8 cap to reset or be
miscounted, and the MAJOR changelog finding describes that unsafe source too strongly. Fix the
counter/validation transaction, add the adversarial cases, correct both changelogs, then re-run
the focused tests and both skill-package gates. **[PRIMARY — CONFIRMED]** The other five round-1
findings are resolved as checked in Q1, and the standard hard-cap behavior should remain.

## Findings index

| Severity | Finding |
| --- | --- |
| CRITICAL | `publishedFixupCycles()` is fail-open and editable, so the binding budget can reset or be miscounted. |
| MAJOR | Both CHANGELOGs state enforcement more strongly than the editable raw prefix count supports. |

## Responses to the other round-1 reviewers

### @hermes-1

I agree with the preserved lifecycle/LE-11 analysis and with leaving the strict-close `>=` operator
alone. The round-1 CLEAN verdict did not inspect the later `publishedFixupCycles()` fix, so it does
not resolve this round's counter-integrity finding.

### @kimi-1

I agree with the round-1 observations that the previous arithmetic tests were synthetic and that
the package/versions/protocol copies otherwise align. The current fix correctly withdraws the
fast claim. I disagree with CLEAN for this revised tree because the new durable-state counter was
introduced after round 1 and the adversarial probes above show that it can reset the cap.

## Open questions

None. The required remediation is code-grounded and does not require an operator policy choice.
