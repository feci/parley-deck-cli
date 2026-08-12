---
agent: kimi-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 7
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 7 — release candidate 1.44.0 + 2.8.0

## Summary

Provenance note first: **there is no kimi-1 round-06 file.** Round-06 contains only @codex-1; my
most recent prior review is round-05, whose timeline section records why (fix-up cycles 5 and 6
landed while my round-5 review was running). This file is therefore my first full pass over the
cycle-6 candidate, run against a tree that held still throughout (mtimes: last candidate write
16:42; my session started later).

Both round-6 MAJORs are genuinely fixed at the behavioural level, and I verified the seam fix the
requested way — deletion in an isolated copy, plus a second reversion @codex-1 did not ask for
(Q1). The user-visible vocabulary is now consistent (Q2). But the release is **not clean**, on two
demonstrated, narrow, in-scope defects:

1. **[MAJOR] A fourth survive-its-own-removal behaviour, on a finding named twice already.** The
   corrected fix-up escalation message at `internal/driver/impl.go:307` has **no shipped
   assertion**. I reverted it to the round-4 defective form in an isolated copy: the full
   26-package suite exited 0. This is the class the brief says has appeared three times; this is
   the fourth, and it sits on the operator-facing text of the safety boundary itself.
2. **[MINOR] The round-6 vocabulary sweep did not reach the test comments it named.** Four of the
   eight quoted locations are **verbatim unchanged** (quotes in Q2), including "The ratified unit
   is published cycles". Cycle 6's own record — "three test comments" realigned — does not match
   the tree: I can find exactly one realigned test comment.

Both fixes are trivial: one test assertion and four comment edits. No redesign is requested, and
the settled trust-anchor scope call is not reopened.

## Q1 — are both round-6 findings actually fixed?

**The seam finding: yes, verified by deletion in an isolated copy — and the test has teeth on both
halves.**

**[PRIMARY — executed by me this session]** Isolated copy via `rsync` (no `.git`) to a `mktemp -d`
directory; the shared trees were never edited. Baseline: `TestTrackWiresTheHardCrossReviewCapThroughNew`
PASS (both subtests). Then two independent reversions:

```
REVERSION 1 — wiring line deleted (driver.go:144, `cfg.HardCrossReviewCap = pol.CapCrossReviewRounds`):
--- FAIL: TestTrackWiresTheHardCrossReviewCapThroughNew/deliberation
    consensus_test.go:422: track deliberation: HardCrossReviewCap=0, want 3 — the policy is not wired into the driver
--- FAIL: TestTrackWiresTheHardCrossReviewCapThroughNew/standard
    consensus_test.go:422: track standard: HardCrossReviewCap=0, want 2 — the policy is not wired into the driver
FAIL  parley-deck-cli/internal/driver      (internal/track: ok)

REVERSION 2 — wiring restored, guard neutered (consensus.go:95, condition prefaced with `false &&`):
--- FAIL: .../deliberation
    consensus_test.go:427: track deliberation: action=reopened err=<nil> — a BLOCK past the §4.0 cap must escalate
--- FAIL: .../standard                                                             (same)
FAIL  parley-deck-cli/internal/driver
```

Reversion 1 is the check the brief asked for: the wiring line cannot disappear silently anymore.
Reversion 2 answers the question the first leaves open — the test fails at the wiring assertion
before reaching the behavioural half, so does the behavioural half have teeth of its own? Yes:
with the wiring intact and only the guard bypassed, both subtests fail on `action=reopened`, i.e.
the "no runner call past the cap" assertion bites. The test drives a real `00-prompt` `track:` →
`New` → blocked consensus with `MaxRounds: 9`, so only the §4.0 cap can stop it — as claimed. Both
subtests red under each reversion, everything restored afterwards; the isolated copy was discarded.

**The vocabulary finding: fixed where it is user-visible, not fixed in the test comments it named**
— details and quotes in Q2.

Also confirmed green on the pristine isolated copy **[PRIMARY]**: `go build ./...`, `go vet ./...`,
`go test ./... -count=1` exit 0 (26 packages `ok`, zero FAIL). Skill repo **[PRIMARY]**:
`node scripts/build-addon-manifest.js --check` all ok; `npm test` exit 0 (386 node tests pass /
0 fail; 54 python tests OK across 7 files); skill working tree hash-identical before and after
(`git status --short | md5` unchanged).

## Q2 — is the vocabulary consistent everywhere?

**User-visible surfaces: yes, verified against the tree.** Source comments and the two error
messages: yes. **Test comments: no — four locations @codex-1 quoted in round 6 are still verbatim.**

Fixed and confirmed **[PRIMARY — read by me in the current tree]**:

- CLI `CHANGELOG.md` 1.44.0: "**`MaxFixupCycles = N` permits attempts 1..N and escalates when N+1
  would start**" and "**The unit is now a reserved fix-up attempt** … deliberately not 'a completed
  cycle'".
- Skill `CHANGELOG.md` 2.8.0: "The run halts with a blocking escalation **when another round or
  attempt would exceed either cap — reaching the last allowed one is not itself a halt**. The CLI
  charges a cycle when it is RESERVED, before the fix-up runs".
- `internal/driver/cursor.go:56-58`: "the driver's own monotonic count of **CHARGED** fix-up
  attempts — **reserved before the code-writing call** … (The name predates the corrected unit;
  `.fixup-done` markers are the record of cycles that actually completed.)"
- `internal/driver/impl.go:146` and `:303`: "cannot determine how many fix-up attempts have been
  **charged**" (was "cycles have been published").
- `internal/driver/impl.go:510-511` (`publishedFixupCycles` doc): "how many fix-up cycles have
  been **CHARGED — reserved, not necessarily completed**".
- `internal/driver/impl_test.go:100-103`: "The §4.0 budget counts driver-**CHARGED** attempts" —
  the one realigned test comment.
- `internal/driver/impl.go:307` now reports the spent count and names the refused cycle:
  `... after %d cycle(s); cycle %d would exceed MaxFixupCycles=%d ...` with `(published, cycle)`.
  (Its content is unpinned — that is Q3.)

Still using the superseded unit, quoted verbatim from the current tree **[PRIMARY]**:

- `internal/driver/impl_test.go:502-503` — "**The ratified unit is published cycles**, so rounds
  that published nothing must not spend any." A flat statement of the superseded unit, in the test
  that pins the zero-spend behaviour. Codex's round-6 bullet quoted exactly this.
- `internal/driver/impl_test.go:638` — "(**the cycle is spent when it runs**)" — fifteen lines
  below the (correct) "reservation is taken BEFORE the code-writing call" comment at :623-624.
  Spent-when-reserved is the whole point of cycle 4; this parenthetical still says when-it-runs.
- `internal/driver/impl_test.go:529` — "**every tamper direction** on the fix-up budget must be
  fail-safe." Universal, and false against this release's own disclosures: the reservation→marker
  single-record window (CLI CHANGELOG: "There is one window with a single record") and the settled
  trust-anchor limits are counterexamples. The test's two subtests cover only the post-marker
  two-record state.
- `internal/driver/impl_test.go:456-457` — "so `fast` (cap 1) **published none at all**" — still
  names the withdrawn end-to-end `fast` route as the example for a test that injects
  `auto_implement: true` into a route `fast` forbids. The historical arithmetic (cap 1 → 0 under
  the old guard) is fine; naming `fast` is the claim cycle 1 withdrew and the CHANGELOG's "Not in
  this release" section carefully does not make.

And the meta-defect: `IMPLEMENTATION.md` "Fix-up cycle 6" claims "Seven statements were realigned
across both CHANGELOGs, the Cursor field doc, two error messages and **three test comments**." The
tree shows **one** realigned test comment (`impl_test.go:100-103`); the other three flagged test
comments are unchanged. A fix-up record that overstates the fix is the claim-outran-code class
this idea has now hit five times — here recorded in a governance document that ships.

Severity calibration: these are comments, not behaviour and not user-facing docs; hence MINOR, not
MAJOR. But Q2's honest answer to "consistent everywhere" is no, and each line above was explicitly
named in the round-6 finding this cycle claims to have closed.

## Q3 — a fourth behaviour that survives its own removal?

**Yes, found and demonstrated: the corrected fix-up escalation message.**

`internal/driver/impl.go:307` is the text an operator reads when the fix-up budget fires — the
safety boundary this entire release is about. Its content was a review finding twice (my round-4
F2: at cap 5 it reported the *refused* cycle as spent — "after 6 cycle(s)"; @codex-1's round-6
probe re-reported it, predating the cycle-5 fix). Cycle 5 corrected the message; @codex-1's round-6
concrete fix asked to "add an assertion for the exact count in the escalation". **No such assertion
exists**: `impl_test.go` contains zero `err.Error()` assertions on this path.

**[PRIMARY — executed by me this session]** In the isolated copy I restored the defective form —
reporting the refused ordinal as the elapsed count:

```go
// reverted to the round-4 form:
fmt.Errorf("review still has %d agreed fixes after %d cycle(s) (MaxFixupCycles=%d); escalating",
    rs.OutstandingAgreedFixes, cycle, d.cfg.MaxFixupCycles)   // `cycle` = published+1
```

```
$ go test ./internal/driver -count=1
ok  parley-deck-cli/internal/driver  0.976s
$ go test ./... -count=1        # full module, message reverted
exit 0                          # 26 packages ok, zero FAIL
```

The defective diagnostic ships tonight and no shipped test can see it. This is the fourth instance
of the class, and the pattern is now precise: **a behaviour named in a review finding, fixed, and
left without the assertion the finding's author asked for.** The class rule the deck keeps paying
for applies to diagnostics on safety boundaries exactly as it applies to the boundaries themselves
— at an escalation, the message *is* the interface.

I re-swept the rest of the release surface for a fifth instance and found none **[PRIMARY —
read/reasoned, reversions where noted]**: the inclusive cap, the deliberation policy cells, the
`HardCrossReviewCap` guard and its wiring, the reservation ordering, the AF2 gate, the
corrupt-vs-absent cursor split, exact round-dir matching, and unreadable-marker propagation all
have shipped tests whose reversions were confirmed red either this round (Q1) or in recorded
earlier rounds; the BLOCK back-edge diagnostic's `next-2` arithmetic *is* pinned
(`consensus_test.go:380` asserts `"after 3 cross-review round(s)"`), which is exactly the
assertion the fix-up message lacks. The contrast between those two diagnostics is the finding in
miniature.

## Q4 — ship or not?

**Not yet.** Two specific defects in what ships, both demonstrated above, neither a preference for
a larger system and neither touching the settled trust-anchor scope call:

1. `internal/driver/impl.go:307`'s escalation message has no regression assertion; its defective
   form passes the full suite (Q3, demonstrated, exit 0). One assertion fixes it — the same shape
   as the existing `consensus_test.go:380` pin on the sibling diagnostic.
2. Four test comments named in the round-6 finding still carry the superseded unit verbatim
   (Q2, quoted), and cycle 6's record claims them realigned. Four comment edits — and the
   `IMPLEMENTATION.md` cycle-6 line corrected to match the tree — fix it.

Everything else about this candidate is in the best state the idea has ever been in: the seam is
tested end to end with teeth on both halves (Q1), the user-visible vocabulary is consistent, the
suite is green, the skill package is installable and tested, and the round-4 fix cluster is
permanently pinned. This is a one-assertion, four-comment cycle — but a release whose headline is
"the printed number and the enforced number are the same number" should not ship with "The
ratified unit is published cycles" in its own test file, or with its boundary diagnostic provably
invisible to its suite.

## Findings index

| Severity | Finding |
| --- | --- |
| MAJOR | `internal/driver/impl.go:307`'s corrected escalation message ships with no assertion; reverting it to the refused-cycle form leaves the full 26-package suite green (demonstrated, exit 0). Fourth instance of the survive-its-own-removal class; the assertion was explicitly requested in round 6. |
| MINOR | Round-6 vocabulary sweep incomplete: `impl_test.go:502-503` ("The ratified unit is published cycles"), `:638` ("spent when it runs"), `:529` ("every tamper direction … fail-safe"), `:456-457` (`fast` named in the cap-1 example) are verbatim as @codex-1 quoted them; `IMPLEMENTATION.md` cycle 6's "three test comments" realigned matches one realigned comment in the tree. |

## Validation record

**[PRIMARY — all commands run by me this session.]** Read-only in both shared trees; every
mutation confined to a `mktemp -d` rsync copy (no `.git`), deleted after use. No git write
commands. `git status --short` of this tree at session end shows the same 20-entry candidate set
as at session start; the skill tree's `git status --short | md5` was identical before and after
my checks (`58ac40fa…`).

```text
ISOLATED COPY (mktemp -d, rsync -a --exclude .git …):
  baseline: go build ./... && go vet ./... && go test ./... -count=1   exit 0, 26 ok
  baseline: -run 'TestTrackWires…|TestBlockedConsensusRespects…' -v    PASS (all)
  REVERT wiring (driver.go:144 deleted):
    go test ./internal/driver ./internal/track -count=1
    → BOTH seam subtests RED ("HardCrossReviewCap=0, want 3/2"); track ok
  RESTORE wiring; REVERT guard (consensus.go:95 → `false && …`):
    → BOTH seam subtests RED ("action=reopened … must escalate")
  RESTORE guard; REVERT message (impl.go:307 → refused-cycle form):
    go test ./internal/driver -count=1      ok
    go test ./... -count=1                  exit 0 — suite blind
  RESTORE message; final pristine re-run    exit 0; copy discarded
SKILL REPO (in place; check/test are read-only):
  node scripts/build-addon-manifest.js --check      all ok
  npm test                                          exit 0 (node 386 pass / 0 fail;
                                                    python 54 OK across 7 files)
```
