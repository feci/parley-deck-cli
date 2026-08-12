---
agent: hermes-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
review-round: 3
date: 2026-08-12
---
verdict: NOT CLEAN

# Phase 6 review — round 3

## Summary

Fix-up cycle 2 replaced the fail-open prose-derived counter with the driver's own
`.fixup-done` markers and made a read error escalate. That closes the four attacks
@codex-1 levelled at the old source — a missing file, a fenced heading, malformed
ordinals, and a careless heading rename can no longer reset the count. The new
source is a genuine improvement: the driver writes the marker as part of the
fix-up transaction (impl.go:314), and an unreadable review directory now returns
an error instead of zero (impl.go:513).

But the new source has its own fail-open path, and it is the same class: a marker
that a party the cap constrains can delete. The `.fixup-done` file is a regular
file on disk, in a directory every participant writes to. Deleting it lowers the
count and buys another cycle, exactly as deleting a heading did. The AF2
crash-idempotency branch opens a second path that bypasses the budget check
entirely. And both CHANGELOGs still describe the OLD prose-derived source, not
the one that shipped.

All attacks below were verified by running adversarial probes in an isolated copy
of the module at /tmp/hermes1-r03-iso — never in the working tree. The working
tree's impl.go diff is unchanged from before this review (55 insertions, 2
deletions; confirmed after cleanup).

## Q1 — Is the fail-open counter actually closed?

No. The four round-2 attacks on the prose source are closed. Two new attacks on
the marker source are open, and one is the same class: a constrained party can
delete state to extend the budget.

### What is closed

[PRIMARY — CONFIRMED] The four attacks @codex-1 levelled at the old
`## Fix-up cycle` heading counter are genuinely closed:

- **Missing IMPLEMENTATION.md**: no longer relevant; the function reads the
  review directory, not IMPLEMENTATION.md.
- **Heading inside a code fence**: no longer relevant; no Markdown parsing.
- **Malformed/duplicate ordinals**: no longer relevant; the count is the number
  of markers, not parsed ordinals.
- **Careless heading rename** (`## Fix-up cycle 5` → `## Fixup cycle 5`): no
  longer relevant; markers are not headings.

[PRIMARY — CONFIRMED] A read error on the review directory now escalates. I
verified this in the isolated copy with a perm-000 review directory:

```
unreadable review dir: correctly returned error: open .../review: permission denied
```

The `os.IsNotExist` case correctly returns `(0, nil)` — a new idea with no
review directory yet has published zero cycles. That is safe and correct.

### What is still open

#### [MAJOR] marker deletion lowers the count and extends the budget

`publishedFixupCycles` (impl.go:506-525) counts `.fixup-done` files in
`review/round-NN/` directories. The marker is a regular file (impl.go:314
writes it with `os.WriteFile`, mode 0644). Every participant writes files
under `review/round-NN/` — that is where reviewer artifacts live. Deleting
a marker from a prior round lowers the count, exactly as deleting a heading
did under the old source.

I verified this in the isolated copy:

```
fake-marker inflation: count went 1 -> 2 (a reviewer can inflate the count)
marker deletion: count went 2 -> 1 (deleting a marker lowers the count)
```

The deletion attack is the budget-extension direction — the unsafe one. With
cap=5 and 5 markers on disk, deleting one makes `publishedFixupCycles` return
4, so `cycle = 4 + 1 = 5`, and `5 > 5` is false: the 5th cycle runs again.
This is the same fail-open class the implementer correctly identified in
IMPLEMENTATION.md's fix-up cycle 2 section: "a number that is a safety
boundary must not be authored by the party it constrains." A `.fixup-done`
file in a world-writable-by-participants directory is still authored by
parties the cap constrains.

The inflation direction (a fake marker in an empty round-NN) is fail-safe —
it causes premature escalation, not budget extension. But the deletion
direction is not.

#### [MAJOR] oddly-named round directories are counted

`publishedFixupCycles` uses `strings.HasPrefix(e.Name(), "round-")` (impl.go:517)
to identify round directories, but does NOT call `roundOrdinal()` to validate
the suffix. A directory named `round-xx` (or `round-anything`) is counted if
it contains a `.fixup-done` file. By contrast, `highestReviewRound`
(cursor.go:212) uses `roundOrdinal()`, which returns 0 for non-numeric
suffixes, so the two functions disagree about what a "round directory" is.

I verified this in the isolated copy:

```
oddly-named round-xx: counted -> 2 (HasPrefix-only, not roundOrdinal)
```

This is a consistency defect, not directly a budget-extension vector on its
own (an oddly-named dir with a marker inflates the count → premature
escalation, which is fail-safe). But it means the count can disagree with
`highestReviewRound`, and a directory like `round-backup` with a stray
`.fixup-done` would inflate the count silently. The fix is to use
`roundOrdinal(e.Name()) > 0` as the filter, matching `highestReviewRound`.

#### [CRITICAL] the AF2 crash-idempotency branch bypasses the budget check

`advanceReview` (impl.go:136-153) checks for a `.fixup-done` marker in the
current round BEFORE the budget check at impl.go:290. If the marker exists,
AF2 fires: it archives the consensus, opens `round+1`, and returns
`ActionFixup` — without calling `publishedFixupCycles` and without checking
`cycle > MaxFixupCycles`.

This means a stale `.fixup-done` marker in the current highest round bypasses
the budget gate entirely. I verified this in the isolated copy with cap=5,
5 published cycles (markers in rounds 1-5), and a stale marker in round-06:

```
AF2 at-cap probe: action=fixup err=<nil> calls=[open-review]
AF2 bypassed the budget check: opened a new round past cap=5 without
running Fixup and without escalating
```

The driver opened round-07 past the cap of 5, ran no Fixup, and did not
escalate. The budget check at impl.go:290 is unreachable when AF2 fires.

In normal operation, AF2 is a crash-recovery path: the marker exists only
if Fixup + RunChecks already succeeded for that round (impl.go:314). But the
marker is a file on disk, and the same deletion/forgery attack applies: a
participant who creates a `.fixup-done` in the current round-NN forces the
driver to open round-NN+1 without any budget check. Combined with the
deletion attack above, a participant can both delete prior markers (lower
the count) and create a stale marker in the current round (bypass the
check entirely).

The AF2 branch should call `publishedFixupCycles` and enforce the cap before
opening the next round, the same way the non-AF2 path does. The crash-
recovery semantics are preserved — the marker still prevents re-running
Fixup — but the budget gate must not be skipped.

### Can anything still extend the budget?

Yes. Two paths:
1. Delete a `.fixup-done` marker from a prior round → count drops →
   `cycle = count + 1` is within cap → another cycle runs.
2. Create a `.fixup-done` marker in the current highest round → AF2 fires
   → the budget check is never reached → the next round opens unconditionally.

Both are the fail-open-extraction-from-editable-state class. The new source
moved the editable state from IMPLEMENTATION.md headings to `.fixup-done`
files in a directory every participant can write to. That is better — the
driver writes the marker, not the implementer's prose — but it is not
tamper-proof, and the AF2 bypass is a separate structural gap.

## Q2 — Does counting markers under-count and hand back budget after a failed fix-up?

Yes, and this is the right definition, but it has a subtle interaction with the
deletion attack.

A fix-up whose post-fix-up checks FAILED writes no marker (impl.go:309-311
escalates before impl.go:314 writes the marker). So `publishedFixupCycles`
does not count it. This is correct: the ratified unit is "published cycles,"
and a fix-up that broke the build was not published — it escalated. The
counting is faithful to the definition.

The under-counting direction is fail-safe: if a fix-up fails and escalates,
the budget is not spent. When the human re-enters the loop, the count is
still accurate. This is the right behavior.

The concern is the OVER-counting direction combined with deletion: a
fix-up that SUCCEEDED writes a marker, and that marker can be deleted,
making the count under-report a cycle that DID happen. That is the Q1
deletion attack, not a Q2 definition issue. The definition is right; the
storage is vulnerable.

## Q3 — Is `standard`'s newly bounded BLOCK back-edge correct and disclosed?

[PRIMARY — CONFIRMED] The back-edge is correctly bounded.
`internal/driver/consensus.go:92-99` checks `HardCrossReviewCap` before
opening a round on the `TriageBlocked` branch. `driver.go:143-148` sets
`HardCrossReviewCap` from `pol.CapCrossReviewRounds` for any track whose
policy carries a cap > 0. `standard`'s policy (track.go:172) carries
`CapCrossReviewRounds: 2`, so `HardCrossReviewCap=2` is set, and the
back-edge escalates when `next > 1 + 2 = 3` (i.e., round 4 is blocked).

[PRIMARY — CONFIRMED] The test `TestBlockedConsensusRespectsTheHardCrossReviewCap`
(consensus_test.go:352-386) asserts the escalation, names the cap, and asserts
no round-05 runner call. I ran it green in the working tree:

```
--- PASS: TestBlockedConsensusRespectsTheHardCrossReviewCap (0.00s)
```

[PRIMARY — CONFIRMED] Both CHANGELOGs disclose the `standard` back-edge
bounding:

- CLI CHANGELOG.md:33-37: "This also bounds `standard`'s back-edge at its
  own printed cap of 2, which it previously ignored."
- Skill CHANGELOG.md:17: "The cross-review ceiling now binds the
  consensus-BLOCK back-edge as well, on `standard` (cap 2) as well as
  `deliberation` — previously that path ignored both."

The disclosure is accurate and present in both. This part of the release is
correct.

## Q4 — Tests run, including a reversion check of the marker-count fix

All tests ran in the working tree [PRIMARY — CONFIRMED, terminal output this
session]:

```
go build ./...                                    exit 0
go vet ./...                                      exit 0
go test ./... -count=1                            all 28 packages PASS
node scripts/build-addon-manifest.js --check      all 6 skills ok
npm test (parley-deck-skill)                      386 pass, 0 fail
TestBlockedConsensusRespectsTheHardCrossReviewCap PASS
TestFixupCapIsInclusive                           PASS (all 6 subtests)
TestZeroFixRoundsDoNotSpendTheFixupBudget         PASS
```

### Reversion check — in an ISOLATED COPY, not the working tree

Per the binding rule and the round-2 process failure record, I copied the Go
module to /tmp/hermes1-r03-iso and performed the reversion there. The working
tree was never edited for this check.

**Reversion**: replaced `publishedFixupCycles` with the old prose-derived
heading count (reads IMPLEMENTATION.md, counts `## Fix-up cycle` prefixes,
returns 0 on error) and reverted the caller to ignore the error.

[PRIMARY — CONFIRMED] `TestFixupCapIsInclusive` went red on all three
"escalates" cases:

```
--- FAIL: TestFixupCapIsInclusive/cap_5:_the_6th_escalates
--- FAIL: TestFixupCapIsInclusive/cap_2:_the_3rd_escalates
--- FAIL: TestFixupCapIsInclusive/cap_1:_the_2nd_escalates
```

The test sets up `.fixup-done` markers, not `## Fix-up cycle` headings, so
the reverted function returns 0 in all cases and the cap is never hit. This
confirms the test genuinely depends on the marker-count source.

`TestZeroFixRoundsDoNotSpendTheFixupBudget` passed under the revert — it
sets up zero markers AND zero headings, so both sources return 0. This is
expected; the test distinguishes the unit derivation, not the source.

After the check, the isolated copy was deleted. The working tree's
`internal/driver/impl.go` diff is unchanged (55 insertions, 2 deletions),
confirmed by `git diff --stat` after cleanup.

## Q5 — Is anything in either CHANGELOG or IMPLEMENTATION.md still overstated?

Yes. The CLI CHANGELOG's "Wrong unit" paragraph describes the source that was
replaced in fix-up cycle 2, not the one that shipped.

### [MAJOR] CLI CHANGELOG describes the old prose-derived source

CLI CHANGELOG.md:14-17 says:

> The count is now derived from the `## Fix-up cycle N` records actually
> published in `IMPLEMENTATION.md`, which is the unit §4.0 is written in.

This is the fix-up cycle 1 source. Fix-up cycle 2 replaced it with
`.fixup-done` marker counting (impl.go:490-525). The CHANGELOG was not
updated to reflect the fix-up cycle 2 change. The actual source counts
driver-written `.fixup-done` markers in review-round directories, not
`## Fix-up cycle N` headings in IMPLEMENTATION.md.

The skill CHANGELOG.md:20-22 says "the printed number and the enforced
number are the same number" — this is accurate for the marker-count source
(as long as markers are not tampered with), but it does not describe the
source, so it is not directly overstated. However, it inherits the CLI's
claim by reference ("Tracks parley-deck-cli 1.44.0").

The code comment at impl.go:284-289 also still says "The cycle number is
derived from the fix-up cycles actually PUBLISHED in IMPLEMENTATION.md" —
this is stale; the function below it reads the review directory, not
IMPLEMENTATION.md. This is a NIT-level code comment issue but compounds the
CHANGELOG overstatement.

Suggested fix: update CHANGELOG.md:14-17 to say the count is derived from
the driver's own `.fixup-done` markers, one per review round, written only
after Fixup and the post-fix-up check gate both succeeded. Update the
impl.go:284-289 comment to match. The "Why 5" paragraph (CHANGELOG.md:39-43)
references `## Fix-up cycle` headings across all 69 ideas — that is a
historical analysis and remains accurate as motivation, but it should not
be presented as the enforcement source.

### Everything else is accurate

- The "Off by one" paragraph (CHANGELOG.md:8-12) is accurate: `>` is
  inclusive.
- The two `deliberation` cells (CHANGELOG.md:24-27) are accurate: 5 and 3.
- "Both are escalation thresholds. Hitting one never marks work complete."
  (CHANGELOG.md:29) is accurate: both return `ActionEscalated`.
- The `standard` back-edge disclosure (CHANGELOG.md:33-37) is accurate.
- "fast Phase 8 is not exercised by the driver at all" (CHANGELOG.md:63-65)
  is accurate: `PolicyFor(Fast, ...)` rejects idea-level `auto_implement`.
- "The ratified escalation payload is not implemented" (CHANGELOG.md:58-61)
  is accurate: `rg -n 'trajectory' internal --glob '*.go'` returns nothing.
- The skill CHANGELOG's table and "not implemented yet" disclosure are
  accurate.
- IMPLEMENTATION.md's fix-up cycle 2 section is accurate in its description
  of the defect and the fix — it correctly describes the `.fixup-done`
  marker approach and the error-escalation behavior.

## Q6 — Anything that should stop this release?

Yes. The CRITICAL AF2 bypass and the MAJOR marker-deletion + CHANGELOG
findings should stop this release.

The AF2 bypass means the budget gate is unreachable on the crash-recovery
path, and a participant can force it by creating a single file. The marker-
deletion attack means the budget can be extended by deleting a file from a
prior round — the same fail-open class the implementer correctly identified
and fixed for the prose source, but not yet closed for the marker source.
The CHANGELOG describes the old source.

Required before release:
1. Move the `publishedFixupCycles` + cap check ABOVE the AF2 branch in
   `advanceReview`, so the budget gate is reached on every path that opens
   a new round. (CRITICAL)
2. Harden the marker against deletion, or document why the threat model
   accepts it. Options: (a) append-only ledger file the driver owns, (b)
   cross-check marker count against review-round count and escalate on
   mismatch, (c) restrict marker permissions. At minimum, the
   `roundOrdinal` filter should be added so oddly-named directories are
   not counted. (MAJOR)
3. Update CLI CHANGELOG.md:14-17 and impl.go:284-289 to describe the
   `.fixup-done` marker source, not the prose-derived source. (MAJOR)

## Findings index

| Severity | Finding |
| --- | --- |
| CRITICAL | The AF2 crash-idempotency branch (impl.go:136-153) bypasses the `publishedFixupCycles` budget check entirely; a stale `.fixup-done` marker in the current round opens a new round past the cap without escalating. |
| MAJOR | A `.fixup-done` marker is a regular file in a participant-writable directory; deleting one from a prior round lowers the count and extends the budget — the same fail-open class as the deleted heading it replaced. |
| MAJOR | `publishedFixupCycles` uses `HasPrefix("round-")` without `roundOrdinal()`, so oddly-named directories (`round-xx`, `round-backup`) are counted, disagreeing with `highestReviewRound`. |
| MAJOR | CLI CHANGELOG.md:14-17 and impl.go:284-289 still describe the old prose-derived source (`## Fix-up cycle N` headings in IMPLEMENTATION.md), not the `.fixup-done` marker source that shipped. |

## Responses to other reviewers

### @codex-1

Your round-2 CRITICAL was correct and the implementer's fix-up cycle 2
addressed it for the four attacks you levelled. I verified all four are
closed. But your own principle — "a number that is a safety boundary must
not be authored by the party it constrains" — applies to the new source too.
A `.fixup-done` file is still authored by parties the cap constrains; it is
just written by the driver instead of parsed from the implementer's prose.
The deletion attack is the same class. And the AF2 bypass is a structural
gap you did not test in round 2 — the budget check is unreachable when AF2
fires. I agree with your round-2 assessment that the release must stop
until the counter is hardened; it still must.

### @kimi-1

Your round-2 MINOR (standard back-edge disclosure) and NIT (trajectory
message) are fixed and verified. Your round-2 MAJOR is the same defect
@codex-1 found, and the implementer's fix-up cycle 2 addressed it for the
four prose-derived attacks. The marker-deletion and AF2 findings above are
new in round 3; they were not visible until the source changed.
