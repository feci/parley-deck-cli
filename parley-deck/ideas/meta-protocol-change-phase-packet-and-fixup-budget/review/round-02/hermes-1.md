---
agent: hermes-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 2
date: 2026-08-12
---
verdict: CLEAN

# Phase 6 review — round 2

## Summary

I reviewed the fix-up cycle 1 changes against the round-1 findings, the ratified
FINAL.md, and both CHANGELOGs. I ran the full Go test suite, the skill package
manifest check and npm test myself, and performed reversion checks on the two
key fixes. All six round-1 findings are genuinely fixed, not relocated. The two
CHANGELOGs are accurate. The `standard` back-edge scope expansion is a coherent
consequence of the fix, not a defect. The `publishedFixupCycles` derivation is
robust in the fail-safe direction. Nothing should stop this release.

## Q1 — Is each round-1 finding actually fixed, or relocated? Check all six.

All six are fixed. I checked each one against the working-tree diff.

### CRITICAL 1 — cross-review cap bypassed on the consensus BLOCK back-edge

FIXED. `internal/driver/consensus.go:90-99` now checks `HardCrossReviewCap`
before the back-edge opens a round: if `next > 1+HardCrossReviewCap`, it
escalates with an error naming the cap. `internal/driver/driver.go:143-148` sets
`HardCrossReviewCap` from `pol.CapCrossReviewRounds` for any track whose policy
carries a cap > 0. The new test `TestBlockedConsensusRespectsTheHardCrossReviewCap`
(`internal/driver/consensus_test.go:352-386`) creates rounds 1-4 with
`HardCrossReviewCap=3`, asserts `ActionEscalated`, asserts the error contains
"§4.0 cap=3", and asserts `len(fr.calls) == 0` — no round-05 runner call. I ran
it green [RECALL: terminal output, this session]:

```
--- PASS: TestBlockedConsensusRespectsTheHardCrossReviewCap (0.00s)
```

This is the behavioral half @codex-1 asked for. The finding is fixed, not
relocated — the cap is now enforced at the back-edge, not only at the initial
scheduling clamp.

### CRITICAL 2 — stale skill payload manifest

FIXED. `skills/parley-deck/parley-addon.json` now carries updated sha256 hashes
for `references/COOPERATION.md` and `references/compatibility.json`, plus a new
aggregate hash. I ran the manifest check and npm test myself [RECALL: terminal
output, this session]:

```
$ node scripts/build-addon-manifest.js --check
parley-deck: ok (6 files, sha256:bee732321d0f5279c5ef83ae308a64e8533337e46be6cf185dae13680a1363db)

ℹ tests 386
ℹ pass 386
ℹ fail 0
```

The 115 failures @codex-1 reported are gone. The finding is fixed.

### MAJOR 1 — "fast" claim was not an actual fast-track path

FIXED (withdrawn). The `fast` end-to-end claim is removed from both CHANGELOGs.
`CHANGELOG.md:62-64` now says "fast Phase 8 is not exercised by the driver at
all — fast forbids idea-level auto_implement, so its fix-up route is manual."
The inclusive-boundary arithmetic fix is retained and correctly applies to every
cap value, but no end-to-end fast-track behaviour is claimed. The test
`TestFixupCapIsInclusive` no longer labels its cases as "fast" — they are now
"cap 1", "cap 2", "cap 5", varying published cycles, not track labels. The
finding is fixed by withdrawing the overstated claim.

### MAJOR 2 — strict closing rounds consume a budget defined as published cycles

FIXED. `internal/driver/impl.go:290` now derives `cycle` from
`publishedFixupCycles(d.cfg.IdeaDir) + 1` instead of `cycle := round`. The new
function `publishedFixupCycles` (`impl.go:490-502`) counts `## Fix-up cycle`
headings in IMPLEMENTATION.md. The new test
`TestZeroFixRoundsDoNotSpendTheFixupBudget` (`impl_test.go:497-517`) creates 4
review rounds with 0 published cycles and cap=2, and asserts that `fixup` is
still called (cycle 1 must run). I ran it green and also confirmed by reversion
that it goes red when `publishedFixupCycles` is reverted to `round` — see Q4
below. The finding is fixed.

### MAJOR 3 — trajectory payload absent

FIXED (withdrawn). Both CHANGELOGs now list the trajectory payload under "Not in
this release" / "not implemented yet." `CHANGELOG.md:58-61` says "The ratified
escalation payload (trajectory, findings by severity, fresh-vs-relitigated,
unresolved fixes, validation status, recommendation) is not implemented." The
sibling `CHANGELOG.md:17-18` says "The structured trajectory payload the idea
ratified is not implemented yet and is a named remaining piece." No Go code
claims to emit the payload. The finding is fixed by withdrawing the overstated
claim.

### NIT — two comments said deliberation applies no overrides

FIXED. `internal/driver/driver.go:119-121` now says "only an absent or unknown
track preserves today's behaviour byte-for-byte. Explicit deliberation now
applies overrides too." `internal/track/track.go:115` now says "false =
reproduce today's behaviour (absent/unknown track only)." Both comments are
accurate.

## Q2 — HardCrossReviewCap bounds standard's back-edge at 2

The implementer's flag is accurate: `HardCrossReviewCap` is now set for every
track whose policy has `CapCrossReviewRounds > 0`, which includes `standard`
(cap 2) — see `internal/track/track.go:172`. Before this change,
`CapCrossReviewRounds` only clamped the initial `CrossReviewRounds` value; the
BLOCK back-edge in `advanceConsensus` was governed by `MaxRounds` alone for both
tracks. The fix applies the cap to the back-edge for both tracks.

This is right, not wrong, and not out of scope. The rationale:

1. `standard` already declared `CapCrossReviewRounds: 2` before this idea — the
   cap value itself is pre-existing and was ratified by the earlier
   track-aware-driver idea, not by this idea.
2. The bug class this idea exists to close is "a cap one code path ignores." The
   BLOCK back-edge was that code path for both tracks. Fixing it only for
   `deliberation` would leave `standard`'s back-edge unbounded while bounding
   `deliberation`'s — incoherent.
3. The §4.0 table already prints "capped at 2, then escalate/upgrade" for
   `standard`. Making the back-edge respect that printed cap makes the code
   match the text, which is the entire point of this idea.

The scope note in IMPLEMENTATION.md (lines 118-121) discloses this transparently
and asks reviewers to confirm. I confirm it. This is a bug fix that affects
`standard`'s back-edge the same way it affects `deliberation`'s, and both now
match their printed cells.

## Q3 — Is publishedFixupCycles robust?

The function (`internal/driver/impl.go:490-502`) reads IMPLEMENTATION.md, splits
on newlines, and counts lines whose trimmed text starts with `## Fix-up cycle`.
I wrote and ran a temporary test exercising the edge cases [RECALL: terminal
output, this session — test removed after verification]:

- **Missing or unreadable IMPLEMENTATION.md**: `os.ReadFile` returns an error,
  the function returns 0. This means `cycle = 1`, so the first fix-up runs. This
  is safe for a new idea (no cycles published yet) and is the same behaviour as
  the old `cycle := round` when round=1. For an idea that has already published
  cycles but whose IMPLEMENTATION.md is missing, the budget effectively resets.
  This is a theoretical concern — IMPLEMENTATION.md is a required protocol
  artifact for any auto_implement idea, and its absence would be caught by other
  guards. The direction is fail-safe against the more dangerous case (premature
  escalation, not budget extension).

- **Heading inside a code fence**: counted. The function does not track fence
  state. A fenced `## Fix-up cycle` heading inflates the count, making
  `cycle = published + 1` larger, which triggers the cap SOONER. This is
  fail-safe: it causes premature escalation, not budget extension. A real
  IMPLEMENTATION.md is unlikely to fence its own fix-up headings, but even if it
  did, the error direction is safe.

- **Re-numbered or duplicated heading**: all matching headings are counted
  regardless of the number after "cycle". Duplicates inflate the count →
  premature escalation (safe). Re-numbering (e.g. `## Fix-up cycle 5` then `##
  Fix-up cycle 10`) produces count=2, which is correct — two cycles were
  published. The number after "cycle" is not parsed; only the heading count
  matters.

- **Malicious or careless edit**: adding headings inflates the count → premature
  escalation (fail-safe). Deleting headings deflates the count → could allow
  more fix-ups than the cap intends (budget extension). This is the one unsafe
  direction. However, IMPLEMENTATION.md is author-controlled by the implementer
  under the protocol's §15 verification-honesty rules — the same trust model
  that already governs every self-reported status in the protocol. The driver
  cannot distinguish a real `## Fix-up cycle` from a fabricated one, but the
  protocol's design assumes the implementer accurately records their work, and
  review rounds exist to catch dishonesty. A malicious implementer who deletes
  headings to extend the budget is violating §15.1, not exploiting a code defect.

**Assessment**: robust enough for this release. The fail-safe direction
(inflation → premature escalation) covers the accidental cases. The unsafe
direction (deletion → extension) is constrained by the protocol's existing trust
model and is not a code-level defect. A future hardening could hash-record the
count or cross-check against review-round artifacts, but that is not required
for the ratified contract.

## Q4 — Reversion check on TestFixupCapIsInclusive

I reverted the fix myself and confirmed the test goes red. Two reversion
variants, both performed in the working tree and restored afterwards:

**Variant A — revert both changes** (`cycle := round` and `cycle >=
MaxFixupCycles`):

```
--- FAIL: TestFixupCapIsInclusive (0.12s)
    --- PASS: cap_5:_the_5th_published_cycle_is_allowed
    --- FAIL: cap_5:_the_6th_escalates
    --- PASS: cap_2:_the_2nd_is_allowed
    --- FAIL: cap_2:_the_3rd_escalates
    --- FAIL: cap_1:_the_1st_is_allowed
    --- PASS: cap_1:_the_2nd_escalates
FAIL
```

3 of 6 cases fail, including "cap 1: the 1st is allowed" — the off-by-one that
made `fast` publish zero cycles.

**Variant B — revert only the unit derivation** (`publishedFixupCycles` →
`round`, keep `>`):

```
--- FAIL: TestFixupCapIsInclusive (0.13s)
    --- PASS: cap_5:_the_5th_published_cycle_is_allowed
    --- FAIL: cap_5:_the_6th_escalates
    --- PASS: cap_2:_the_2nd_is_allowed
    --- FAIL: cap_2:_the_3rd_escalates
    --- PASS: cap_1:_the_1st_is_allowed
    --- FAIL: cap_1:_the_2nd_escalates
--- FAIL: TestZeroFixRoundsDoNotSpendTheFixupBudget
    round 4 with cap 2 but zero published cycles must still run cycle 1
```

3 "escalates" cases fail (the review-round ordinal does not match the published
count), and `TestZeroFixRoundsDoNotSpendTheFixupBudget` fails (round 4 with cap
2 escalates instead of running cycle 1). After restoring the file, all tests
returned green. The implementer's claim is verified.

## Q5 — Are the two CHANGELOGs accurate and not overstated?

Both are accurate. I checked every claim against the code.

**parley-deck-cli CHANGELOG.md (1.44.0)**:

- "Off by one" — `impl.go:291` uses `>`, confirmed inclusive. Accurate.
- "Wrong unit" — `impl.go:290` uses `publishedFixupCycles`, confirmed. Accurate.
- "Both are escalation thresholds. Hitting one never marks work complete." — both
  boundaries return `ActionEscalated`, not `ActionComplete`. Accurate.
- "The cross-review cap binds every path that opens a round, including the
  consensus-BLOCK back-edge" — the forward path is bounded by the clamped
  `CrossReviewRounds`, the back-edge by `HardCrossReviewCap`. There are two
  paths that open Phase-2 rounds (`advanceRound` forward, `advanceConsensus`
  BLOCK back-edge) and both are capped. Accurate.
- "fast Phase 8 is not exercised by the driver at all" — `PolicyFor(Fast, ...)`
  rejects idea-level `auto_implement` (`track.go:142-149`), so the driver fix-up
  route is unreachable for fast. Accurate.
- "The ratified escalation payload is not implemented" — `rg -n 'trajectory'
  internal --glob '*.go'` returns nothing. Accurate.
- "This release makes Parley Deck correct, not faster." — the packet (slice 2)
  is not started. Accurate.

**parley-deck-skill CHANGELOG.md (2.8.0)**:

- "The CLI enforces both cells in the same release, so the printed number and the
  enforced number are the same number." — the CLI sets `MaxFixupCycles: 5` and
  `CapCrossReviewRounds: 3` for `deliberation` (`track.go:161-166`), and both are
  wired through to the driver (`driver.go:143-151`). The printed text says "cap
  5 cycles" and "capped at 3 after round 1, then escalate." These match.
  Accurate.
- "The structured trajectory payload the idea ratified is not implemented yet"
  — accurate, same as above.

No remaining claim is overstated. The `fast` claim and the trajectory payload
claim were both withdrawn and are now listed under "Not in this release" / "not
implemented yet."

## Q6 — Does the skill package install?

Yes. I ran both commands myself [RECALL: terminal output, this session]:

```
$ node scripts/build-addon-manifest.js --check
parley-bidding: ok
parley-deck: ok (6 files, sha256:bee732321d0f5279c5ef83ae308a64e8533337e46be6cf185dae13680a1363db)
parley-design: ok
parley-design-check: ok
parley-tracker: ok
parley-worktrees: ok

$ npm test
ℹ tests 386
ℹ pass 386
ℹ fail 0
python 3.14: 54 tests OK across 7 files
```

The manifest check reports all skills "ok." The node test suite is 386/386
green. The python tests are 54/54 green. The package installs.

## Q7 — Anything that should stop this release?

No.

All six round-1 findings are fixed. The full Go test suite is green
(`go test ./... -count=1`, all packages pass). `go vet ./...` is clean.
`gofmt -d` on all changed Go files produces no output. `git diff --check` is
clean on both repos. The skill package manifest check and npm test are green.
The protocol drift guard `TestEmbeddedDefaultMatchesLiveDeck` passes, confirming
all three COOPERATION.md copies agree. VERSION and version.go are bumped to
1.44.0; package.json and compatibility.json are bumped to 2.8.0.

The `standard` back-edge scope expansion (Q2) is a coherent consequence, not a
defect. The `publishedFixupCycles` derivation (Q3) is robust in the fail-safe
direction. The reversion check (Q4) confirmed the test is genuine. Both
CHANGELOGs are accurate.

## Independent validation record

All checks ran in this working tree [RECALL: terminal output, this session]:

```
go test ./... -count=1                              all packages PASS
go vet ./...                                        exit 0
gofmt -d (changed Go files)                         no output
git diff --check (both repos)                       exit 0
node scripts/build-addon-manifest.js --check        all ok
npm test                                            386 pass, 0 fail
python tests                                        54 pass, 0 fail
TestEmbeddedDefaultMatchesLiveDeck                  PASS
reversion: publishedFixupCycles → round             3+1 cases FAIL (red)
reversion: both fixes reverted                      3 of 6 cases FAIL (red)
post-restore                                        all PASS
```

## Responses to other reviewers

### @codex-1

Your round-1 findings were precise and your reprostructions were honest. Both
CRITICALs were real, and the implementer's fix-up cycle addresses them correctly.
I specifically verified the two areas I missed in round 1 — the packaged skill
manifest and the consensus BLOCK back-edge — and confirm they are fixed. Your
MAJOR on the strict-gate budget unit is addressed by the `publishedFixupCycles`
derivation, which I verified by reversion. Your MAJOR on the trajectory payload
is addressed by withdrawing the claim. No open items from your round-1 findings
remain.
