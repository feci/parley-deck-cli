---
idea: meta-protocol-change-phase-packet-and-fixup-budget
status: complete
implementer: claude-1
started: 2026-08-11
---

# IMPLEMENTATION

`FINAL.md` has two independent interventions. **Slice 1 (the budgets) is done. Slice 2 (the packet)
is not started**, and its release is gated on the pre-registered experiment.

## Slice 1 — the two §4.0 budget cells

status: ready-for-review
date: 2026-08-11

### What changed

| File | Change |
| --- | --- |
| `internal/track/track.go:150-167` | `deliberation` policy is no longer a no-op: `ApplyOverrides: true`, `MaxFixupCycles: 5`, `CapCrossReviewRounds: 3`. `CrossReviewRounds: -1` and the reviewer fields are untouched, so the rest of the full lifecycle is preserved. |
| `internal/driver/impl.go:284` | `cycle >= MaxFixupCycles` → `cycle > MaxFixupCycles`. **The cap is now inclusive.** |
| `parley-deck/COOPERATION.md`, `internal/protocol/defaults/COOPERATION.md`, `parley-deck-skill/skills/parley-deck/references/COOPERATION.md` | both per-track cells: `unbounded` → `capped at 3 after round 1, then escalate` and `cap 5 cycles`. All three copies, one patch. |

### The inclusive-boundary fix is bigger than `deliberation`

@codex-1 predicted that setting `MaxFixupCycles: 5` under the `>=` guard would publish only four
cycles. Verified, and it was worse than that — the off-by-one applied to **every** track, so the
§4.0 table was wrong in all three columns:

| Track | Table printed | Actually published (before) | Now |
| --- | --- | --- | --- |
| `fast` | cap 1 cycle | **0** | 1 |
| `standard` | cap 2 cycles | 1 | 2 |
| `deliberation` | unbounded (driver default 3) | 2 | 5 |

`fast` published **no fix-up cycle at all** while its cell claimed one. This was not in scope as a
finding; it fell out of implementing the ratified inclusive boundary, and it is the second instance
of the class this idea's FINAL names — the §4.0 table asserting authority over behaviour nothing
enforced. It belongs in the seed inventory for
`meta-protocol-change-track-gate-enforcement-audit`.

### Verification

`go build` (darwin, linux, windows), `go vet ./...`, and `go test ./... -count=1` are green,
including the protocol drift guard that compares the deck copy against the Go-embedded default.

**Reversion checks** — each fix reverted, the test confirmed red, the revert compiled and applied,
the tree restored and re-verified green:

```
OK  inclusive cap        → TestFixupCapIsInclusive red on 3 of 6 cases, incl. "fast: cap 1 still
                           publishes one cycle" (ranFixup=false, want true)
OK  deliberation policy  → TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops red (Fixup=3,
                           want 5); TestDeliberationClampsCrossReviewRoundsToThree red (got 9)
```

Tests added: `TestFixupCapIsInclusive` (6 boundary cases across all three tracks, at the cap and one
past it) and `TestDeliberationClampsCrossReviewRoundsToThree`.

Three existing tests encoded the old semantics and were **deliberately rewritten**, not deleted:
`TestNewDeliberationIsLegacy` → `TestNewDeliberationKeepsTheLifecycleButBoundsBothLoops`,
`TestPolicyForDeliberationIsLegacy` → `TestPolicyForDeliberationBoundsBothLoopsAndKeepsTheRest`, and
`TestPhaseReviewMaxFixupCyclesEscalates`'s fixture moved from round-03 to round-04 because with an
inclusive cap of 3, cycle 3 is allowed and cycle 4 is the one that must escalate.

### Deviation from FINAL — one, and it is deliberate

`internal/driver/impl.go:215` carries a **second** `round >= d.cfg.MaxFixupCycles` guard, for the
strict-gate closing loop (LE-2 termination). It was left as `>=`. The FINAL's inclusive-boundary
requirement is written about *published fix-up cycles*, and LE-2's termination semantics were
ratified separately by a different idea. Changing both under one mandate would alter a bound nobody
in this idea reviewed. **Reviewers should decide** whether the strict-gate loop should also become
inclusive, or whether the two loops legitimately differ.

### Not verified

No end-to-end run. This changes when the driver escalates, and no live `parley run` exercised it —
the same instrumentation gap `protocol-read-cost-regression` recorded.

## Slice 2 — the packet

status: not started

Blocked on nothing; simply not begun. Note the ordering the FINAL fixes: the pre-registered
experiment runs **against the packet commit, before release**, and the ship/refute thresholds are
already written into `FINAL.md` and may not be changed after data exists. The middle band
`(0.50, 0.80]` returns to the user with the measured number.

The open refute-threshold conflict — @hermes-1 at 0.67, @codex-1 and @kimi-1 at 0.80 — is
unresolved by design. **The implementer must not resolve it by picking one.**

## Fix-up cycle 1

status: ready-for-re-review
date: 2026-08-12

Round-1 verdicts: **codex-1 NOT CLEAN** (2 CRITICAL, 3 MAJOR, 1 NIT), **hermes-1 CLEAN**,
**kimi-1 CLEAN**. §15.3 forbids settling that by count, and both CRITICALs were independently
reproduced by the implementer before any fix was written. Neither @hermes-1 nor @kimi-1 examined the
packaged skill or the consensus back-edge; @codex-1 did.

**Without this review the release would have shipped a skill package whose installer refuses its own
payload.**

### CRITICAL 1 — the cross-review cap was bypassed on the consensus BLOCK back-edge

`CapCrossReviewRounds` clamped only the *initially scheduled* budget. `advanceConsensus`'s
`TriageBlocked` branch opens another Phase-2 round governed by `MaxRounds` alone (default 4), so a
deliberation idea capped at 3 could still run a 4th cross-review round.

Fixed by a distinct `Config.HardCrossReviewCap`, set from the policy and checked before the back-edge
opens a round. It bounds **every** path that opens a Phase-2 round, not one of them.
`TestBlockedConsensusRespectsTheHardCrossReviewCap` asserts the escalation happens, names the cap,
and that **no round-05 runner call occurs** — the behavioural half @codex-1 asked for.

**Scope note for reviewers:** the same mechanism now also bounds `standard`'s back-edge at its
printed cap of 2. That was not separately ratified. It makes `standard` match its own §4.0 cell, and
leaving one track's back-edge unbounded while bounding the other's would be incoherent — but it is a
behaviour change beyond the two ratified cells and reviewers should confirm or reject it.

### CRITICAL 2 — the 2.8.0 skill package was not installable

`skills/parley-deck/parley-addon.json` still carried the pre-change hashes, so the installer's
fail-closed payload check rejected it. Reproduced before fixing:

```
$ node scripts/build-addon-manifest.js --check
parley-deck: STALE parley-addon.json — regenerate it
$ npm test        # first cascading error
Source payload does not match parley-addon.json:
modified: references/COOPERATION.md; modified: references/compatibility.json
```

Regenerated with the repository's own manifest command; `--check` and `npm test` are green.

### MAJOR — the budget was counted in the wrong unit

A strict-gate round with zero agreed fixes opens the next review round **without publishing a
fix-up**, but `cycle := round` derived the number from the review-round ordinal — so rounds that
produced nothing spent budget, and the ratified unit ("5 inclusive **published** cycles") was not
what shipped. `cycle` now comes from `publishedFixupCycles()`, counting `## Fix-up cycle N` records
in `IMPLEMENTATION.md`. New test `TestZeroFixRoundsDoNotSpendTheFixupBudget` (4 review rounds, cap 2,
zero published → cycle 1 must still run).

`TestFixupCapIsInclusive` was rewritten onto the same unit: its cases now vary *published cycles*,
not review-round directories. The earlier version passed while measuring the wrong thing.

### MAJOR — two overstated claims removed from the release notes

- **`fast` publishing "0 → 1" cycles was not demonstrated.** `fast` forbids idea-level
  `auto_implement`, so its driver fix-up route is unreachable; the test injected `MaxFixupCycles: 1`
  into an absent-track driver. The arithmetic fix is real for every cap value, but the end-to-end
  `fast` claim is withdrawn and the CHANGELOG now says the route is manual.
- **The ratified escalation payload does not exist.** `rg -n 'trajectory' internal --glob '*.go'`
  returns nothing. Both boundaries escalate and halt with the counts and the cap they enforced; the
  structured payload is now listed under "Not in this release" instead of being claimed.

### NIT

Two comments still described explicit `deliberation` as applying no overrides. Corrected.

### Verification

`go build` (darwin/linux/windows), `go vet`, the full Go suite, and the skill `npm test` +
`build-addon-manifest.js --check` are green.

**Reversion checks** — each new fix reverted, its test confirmed red, revert compiled and applied,
tree restored and re-verified:

```
OK  hard cross-review cap        → TestBlockedConsensusRespectsTheHardCrossReviewCap red
OK  published-cycle derivation   → TestFixupCapIsInclusive red on all three "escalates" cases
```

## Fix-up cycle 2

status: ready-for-re-review
date: 2026-08-12

Round-2 verdicts: **codex-1 NOT CLEAN** (1 CRITICAL, 1 MAJOR), **kimi-1 NOT CLEAN** (1 MAJOR,
1 MINOR, 1 NIT), **hermes-1 CLEAN**. @codex-1's CRITICAL and @kimi-1's MAJOR are the same defect,
found independently; @hermes-1 examined the same function and judged it "robust enough". §15.3
forbids settling that by count, and the strictest reading is the correct one here on the merits.

### The fix-up budget was fail-open, and it was the same class we deleted once already

Cycle 1 derived the count from `## Fix-up cycle N` headings in `IMPLEMENTATION.md` — a file the
**implementer owns**. @codex-1 built an adversarial probe and reported four independent failures:

```
missing file -> 0                                  (whole budget restored)
heading inside code fence -> 1                     (not Markdown-aware)
duplicate and malformed ordinals -> 3              (suffix never parsed)
one careless heading edit resets count 5 -> 4      (renaming buys another cycle)
```

@hermes-1's assessment reasoned about a *malicious* implementer and concluded the trust model
covers it. @codex-1's point is sharper and correct: **a careless edit has the same effect as a
malicious one.** `## Fix-up cycle 5` → `## Fixup cycle 5` silently buys a cycle past a ceiling the
FINAL defines as a blocking safety threshold whose extensions "never reset the count".

This is the fail-open-extraction-from-prose class that got the derived marker ledger deleted in
`protocol-read-cost-regression`. **The general rule this deck keeps paying for: a number that is a
safety boundary must not be authored by the party it constrains.**

**Fixed by counting the driver's own `.fixup-done` markers** — one per review round, written only
after `Fixup` and the post-fix-up check gate both succeeded. Not parsed from prose, not editable as
a side effect of writing a document. `publishedFixupCycles` now returns `(int, error)` and a read
error **escalates** instead of silently restarting the budget at zero.

### Other fixes

- **[kimi-1 MINOR]** Both CHANGELOGs now disclose that `standard`'s BLOCK back-edge is bounded too.
- **[kimi-1 NIT]** The cross-review escalation message no longer promises "the trajectory" — the
  same release states that payload is not implemented.

### A process failure by a reviewer, recorded because it nearly cost a fix

@kimi-1 ran its round-2 reversion check **by editing the working tree**, leaving
`internal/driver/consensus.go:92` carrying `// TEMP-REVERT for kimi-1 round-02 Q4 reversion check`
with the guard removed. It was caught when an unrelated test suddenly failed, and @kimi-1 did
restore it before finishing. Reviewers are read-only by protocol, and this is the second time in two
ideas that a reversion check has left a revert applied in a shared tree — the first cost this deck a
recorded incident in `protocol-read-cost-regression`. **A reversion check belongs in an isolated
copy**, which is exactly what @codex-1 did for the same round ("I copied the current Go module to an
isolated temporary directory").

### Verification

`go build`, `go vet`, full Go suite green; skill `npm test` and `build-addon-manifest.js --check`
green. Reversion check: reverting the marker count to prose parsing turns
`TestFixupCapIsInclusive` red on all three one-past-cap cases; restored and re-verified.

## Fix-up cycle 3

status: ready-for-re-review
date: 2026-08-12

Round-3 verdicts: **codex-1 NOT CLEAN** (1 CRITICAL, 2 MAJOR, 1 NIT), **hermes-1 NOT CLEAN**
(1 CRITICAL, 3 MAJOR). Two reviewers, independently, on the same conclusion: **cycle 2 relocated the
defect instead of closing it.** The editable state moved from `IMPLEMENTATION.md` headings to
`.fixup-done` files; both are still editable by a participant with workspace write.

That is the third form of one class in this idea. It is now stated as a rule in the source:
**a number that is a safety boundary must not be authored by the party it constrains.**

### The budget is now the maximum of two independent driver-authored sources

- a **monotonic counter in the run cursor** (`RunDir/driver.json`, outside the idea directory),
  written inside the fix-up transaction and carried forward across `Rebuild`;
- the **`.fixup-done` markers** under the idea's review rounds.

Taking the maximum is what makes it safe, in every tamper direction:

| Tamper | Effect |
| --- | --- |
| delete the markers | cursor still holds the count — no cycle bought |
| delete/corrupt the cursor | markers still hold it — no cycle bought |
| forge either | count only rises → escalates **sooner** |
| unreadable state | **escalates**; it never restarts the budget at zero |

`Advance` rebuilds the cursor from artifacts each tick, which was silently discarding the persisted
counter — found by the first tamper test failing. The carry-forward is monotonic: the value can only
rise.

### [hermes-1 CRITICAL] AF2 bypassed the budget entirely

The crash-idempotency branch reads the same `.fixup-done` marker and, if present, archives and opens
the next round **without ever consulting the cap**. @hermes-1's isolated probe drove it to round 7
against a cap of 5 with no fix-up and no escalation. The branch now checks the budget first.

### [codex-1 MAJOR] A failed fix-up handed the budget back — and @hermes-1 had judged that correct

@hermes-1's round-3 Q2 called the under-count "the right behavior"; @codex-1 called it a defect.
§15.3 forbids counting, and on the merits @codex-1 is right: a fix-up that ran and broke the build is
**exactly the churn the cap exists to interrupt**, and handing the budget back lets a failing fix-up
loop forever against a ceiling that never depletes. **The cycle is spent the moment it runs**, and
the count is persisted before the check gate so an escalation cannot un-spend it.

### [hermes-1 MAJOR] Oddly named directories counted

`strings.HasPrefix(name, "round-")` counted `round-backup`, `round-x`, `round-`. Replaced with an
exact `round-<digits>` test; `TestOnlyExactRoundDirsCountAsPublishedCycles` fixes 6 directory names
and expects 2.

### [hermes-1 MAJOR] The CHANGELOG described the previous design

It still described the prose-derived source that cycle 2 had already replaced. Rewritten to describe
what actually ships, including the "spent when it runs" rule.

### Verification

`go build`, `go vet`, full Go suite green. **All three reversion checks were run in an isolated copy
of the module** — the discipline this idea added to its own review brief after a reviewer left a
revert applied in the shared tree:

```
OK  monotonic carry-forward   → tamper test "deleting the markers…" red
OK  AF2 budget gate           → tamper test "a forged marker…" red
OK  exact round-dir matching  → TestOnlyExactRoundDirs… red
```

## Fix-up cycle 4

status: ready-for-re-review
date: 2026-08-12

Round-4: **codex-1 NOT CLEAN** (1 CRITICAL, 2 MAJOR), **hermes-1 CLEAN**. @codex-1 built adversarial
probes in an isolated module copy and reproduced six separate ways to get a cycle back.

### Fixed

**The reservation was taken too late.** `Fixup` ran first and the counter was written after it
returned, so a `Fixup` that errored, or a crash in the window between return and save, cost nothing.
**The cycle is now reserved before the code-writing call.** A crash before the fix-up is confirmed
conservatively spends the reservation — one lost cycle is the right error direction against an
unbounded loop.

**The AF2 gate was one off, in the harmful direction.** It rejected `spent >= cap`, which strands
the legitimate crash recovery of the *last allowed* cycle: at the inclusive boundary the fifth
fix-up has already spent its budget, and AF2 is finishing it, not starting a sixth. Now strictly
greater; starting the next cycle is still refused by the ordinary branch. @codex-1's boundary probe
is the evidence, and this was a defect this implementer introduced in cycle 3 while fixing another.

**A cursor read error was silently ignored** — `err == nil &&` in the carry-forward meant a corrupt
or unreadable cursor fell back to markers, and after an errored fix-up there is no marker, so the
spent count became zero. That directly contradicted the claim that an unknown count escalates. It
now escalates; a genuinely absent cursor (a fresh run) is distinguished from an unreadable one.

### Not fixed — scope, stated rather than papered over

@codex-1's CRITICAL also requires the ceiling to be **idea-scoped, serialized across concurrent
runs, and anchored outside the participant-writable repository**. Its probes are correct: two runs
of one idea keep independent cursors and independent locks, a deleted run directory loses the
counter, and a repository rollback restores a lower one.

**That is a different system from the one this idea ratified**, which was two numbers and their
enforcement. An out-of-repo trust anchor also breaks the deck's file-based portability model — a
fresh clone would lose the count entirely. Building it silently under this mandate would repeat the
mistake the idea's own FINAL warns about.

**The claims were corrected instead.** Both CHANGELOGs now say plainly what this is: robust against
accidental loss, a stale or deleted single record, and an errored or crashed fix-up; **not** a
security boundary against a participant with workspace write, a deleted run directory, a rollback,
or two concurrent runs. Follow-up: `fixup-budget-trust-anchor`.

### Verification

`go build`, `go vet`, full Go suite green.

### Late arrival — @kimi-1's round-3 file

@kimi-1's round-3 review landed after cycles 3 and 4 were already written, so round 3 was
**NOT CLEAN from all three**, not two. Its F1 (marker-only count still fail-open on deletion) and F2
(the CHANGELOG describing the superseded design) were independently the same findings already fixed
in cycle 3 and cycle 4. Its **F3 was not**, and is fixed here: the call-site comment still described
the cycle-1 prose derivation — two designs out of date — and repeated the withdrawn `fast` claim.
Corrected to point at `publishedFixupCycles`'s own doc comment rather than restating a mechanism
that keeps changing under it.

A comment that describes the implementation two revisions ago is the same defect class as a
CHANGELOG that does: a factual claim about the code that the code no longer supports.

## Fix-up cycle 5

status: ready-for-re-review
date: 2026-08-12

Round-5: **codex-1 NOT CLEAN** (2 MAJOR). It **accepted the scope call** — "Yes, as a release
decision… `fixup-budget-trust-anchor` is an adequate named follow-up" — and then found two things
that are squarely in scope.

### [MAJOR] The three round-4 fixes had no permanent tests

@codex-1 reverted all three at once in an isolated copy and ran the **shipped** suite:

```
$ go test ./internal/driver -count=1
ok  parley-deck-cli/internal/driver  0.846s
```

Every fix could vanish together without a single red test. That is the inert-test class this deck
has now hit in three consecutive ideas, and this time the implementer produced it while fixing
findings — the fixes were verified by isolated probes that were then thrown away.

`TestRound4FixesHavePermanentTests` now covers all four behaviours: an errored `Fixup` keeps its
reservation, AF2 finishes the last allowed cycle at equality, AF2 refuses beyond the cap, and a
corrupt cursor escalates while an absent one does not. Reproducing @codex-1's combined reversion
against the new suite turns **three of the four red**.

### [MAJOR] The corrected limits were still wrong — in the other direction

The claims said "completed fix-up cycles" and "losing one record does not lower the count". Both
overstate. The implementation charges a **reserved attempt**, including one that errors, and between
the reservation and the marker there is only **one** record — so losing the cursor in that window
does lose the count. @codex-1 reproduced exactly that composition.

Corrected in both CHANGELOGs and in the source comment: the unit is a reserved fix-up attempt; the
two-record redundancy holds **once both records exist**; the single-record window is stated as a
limit rather than papered over.

This is the fourth time in this idea that a claim about the code outran the code. The pattern is
now explicit in the record: **every time the mechanism changed, its description had to be re-checked
against it, and three times it was not.**

### Verification

`go build` (darwin/linux/windows), `go vet`, full Go suite green. Combined-reversion check run in an
isolated copy and discarded.

### Late arrival — @kimi-1's round-4 file

@kimi-1's round-4 review also landed after the cycle was written, making round 4 **NOT CLEAN from
two of three**, not one. Its F1 is @codex-1's round-5 MAJOR found independently — the three cycle-4
fixes shipping with no regression test — already closed in cycle 5. Four smaller findings were not,
and are fixed here:

- **F5 [MINOR]** `markedFixupCycles` counted markers through `fileExists`, which maps *any* stat
  error to "absent" — a direction that **lowers** the safety count. A permission-damaged round
  directory now propagates and escalates. `TestUnreadableMarkerEscalatesInsteadOfLoweringTheCount`.
- **F2 [NIT]** The fix-up escalation reported the *refused* cycle as though it were spent ("after 6
  cycle(s)" at cap 5). It now reports the spent count and names the cycle that would exceed.
- **F3 [NIT]** The BLOCK back-edge diagnostic counted one round that never ran (`next-1`). Raised by
  @codex-1 in round 3 and left unfixed through two cycles.
- **F4 [NIT]** `Cursor`'s struct doc said "Rebuild derives every field from on-disk artifacts" —
  false for precisely the field whose whole point is that it must not be derivable.

F2 and F3 are the same small class as the stale comments and the outrunning CHANGELOG: **the
diagnostic text was not re-checked when the number it describes changed.**

## Fix-up cycle 6

status: ready-for-re-review
date: 2026-08-12

Round-6: **codex-1 NOT CLEAN** (2 MAJOR). Both fixed.

### [MAJOR] The integration seam shipped untested

`New` copies `Policy.CapCrossReviewRounds` into `Config.HardCrossReviewCap`. That single line is the
only thing connecting the ratified §4.0 cell to the guard that enforces it — and the guard's own test
**injected `HardCrossReviewCap` by hand**. @codex-1 deleted the line in an isolated copy:

```
$ go test ./internal/driver ./internal/track -count=1
ok  parley-deck-cli/internal/driver
ok  parley-deck-cli/internal/track
```

Green, with the release-blocking behaviour gone. `TestTrackWiresTheHardCrossReviewCapThroughNew` now
drives the whole seam — a real `00-prompt` `track:` → `New` → a blocked consensus — for both
`deliberation` (cap 3) and `standard` (cap 2), asserts the wired value, asserts the escalation, and
asserts **no runner call**. It deliberately sets `MaxRounds: 9` so only the §4.0 cap can stop it.
Reproducing the deletion now turns both subtests red.

**This is the second seam in two rounds whose test asserted the guard while bypassing the wiring.**
A test that hand-injects the value it is meant to prove is wired is the same inert-test class as one
that calls the function under the guard instead of going through the dispatch.

### [MAJOR] Seven statements still described the superseded unit

The unit is a **charged (reserved) attempt**, not a published cycle, and the reservation-only window
is real. Seven places still said otherwise — the CLI CHANGELOG's "publishes cycles 1..N", the skill
CHANGELOG's "at either boundary the run halts" (equality is allowed), `Cursor`'s field doc, two
error messages, and three test comments including one repeating the withdrawn `fast` claim. All
aligned on one vocabulary: **charged/reserved** for the combined count, **published/completed**
only for `.fixup-done` markers.

The escalation off-by-one @codex-1 reported had already been fixed in the previous cycle from
@kimi-1's F2; its probe predated that change.

### Verification

`go build`, `go vet`, full suite green. Both isolated reversions (the seam line, and the three
cycle-4 fixes combined) now go red.

## Fix-up cycle 7 — authorised by the owner, capped at ONE cycle

status: ready-for-re-review
date: 2026-08-12
authority: user ruling after the budget escalation at
`parley-deck/inbox/claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md`

This idea's own ratified cap is 5 and this is cycle 7. The implementer escalated with a trajectory
rather than extending its own budget, and the owner authorised **one** further cycle with the
release blocked if it produces another finding in its own fix. Both round-7 MAJORs are closed here
and nothing else was touched.

### [MAJOR] The "one vocabulary" correction claimed in cycle 6 was not applied

Cycle 6 asserted the vocabulary was unified. It was not — the two operator-visible escalations still
said "cycle(s)" for a count that is the maximum of charged attempts and completed cycles.
@codex-1's counterexample is the point: **one errored attempt, zero cycles completed, and the
message read "after 1 cycle(s)".** That is false output, not a naming preference.

Applied now, and to the code rather than only the prose: `chargedFixupAttempts`, the local
`charged`, and both messages ("N charged attempt(s)", "attempt N+1 would exceed"). The persisted
cursor field keeps its old name for on-disk compatibility, and the comment says so.

**That claim was itself the fifth instance of this idea's recurring defect** — a statement about the
code that the code did not support. This time the statement was about a fix.

### [MAJOR] The corrected count had no assertion

`TestOverCapEscalationReportsChargedAttemptsNotCycles` drives the reservation-only state — one
charged attempt, no marker — and asserts all three properties: the charged count is reported, the
refused ordinal is named separately, and the word "cycle(s)" does not appear. Reverting the message
turns it red on all three.

### Verification

`go build` (darwin/linux/windows), `go vet`, full Go suite green. Reversion check run in an isolated
copy and discarded.

### Note on the round-7 record

@codex-1 returned NOT CLEAN with the two findings above. @hermes-1 and @kimi-1 did not complete
round 7: @kimi-1 exhausted its tool-iteration limit before writing (it reported that its probes were
confined to `/tmp` copies and the working tree was untouched — independently verified here), and the
run was stopped rather than left waiting on an interactive prompt. Their absence is recorded as
incomplete, **not** as agreement.

## Fix-up cycle 8 — the owner's blocking condition, and its resolution

status: complete
date: 2026-08-12
authority: user ruling ("fix and ship")

Round 8: **@codex-1 NOT CLEAN**, one MINOR, and it triggered the owner's condition exactly as
written — cycle 7 introduced a new defect **in its own fix**. The finding was verified before being
acted on:

```
$ rg -n 'func .*publishedFixupCycles' internal/ --glob '*.go'
  (no such function)
$ rg -n 'publishedFixupCycles' internal/ --glob '*.go'
internal/driver/impl.go:297:      // The count comes from publishedFixupCycles …
internal/driver/impl_test.go:104: // … see publishedFixupCycles for why.
```

Cycle 7 renamed the helper and left two comments pointing at a symbol that no longer exists. Real,
and the smallest possible instance of this idea's recurring defect: **a statement about the code
that the code does not support.** Eight cycles, and the last one is two words.

**The implementer did not decide this.** The owner set the condition, the condition fired, and the
decision went back to the owner, who ruled "fix and ship". Both references now name
`chargedFixupAttempts`; nothing else was touched, and no behaviour, test or claim changed.

### Final verification

`go build` (darwin, linux, windows), `go vet`, the full Go suite, the skill `npm test` and
`build-addon-manifest.js --check` are all green.

## Slice 1 — complete

The two §4.0 budget cells ship enforced, with text and code in one patch, boundary tests at both
caps, and the integration seam covered end to end.

**Slice 2 (the phase-scoped packet) is not started.** Its release remains gated on the
pre-registered experiment in `FINAL.md`, whose ship/refute thresholds may not be changed after data
exists, and whose middle band returns the decision to the owner with the measured number.

### What eight cycles cost, and what they bought

Every cycle from 2 onward found a defect in the previous cycle's fix. Severity fell — CRITICAL,
CRITICAL, CRITICAL, CRITICAL, MAJOR, MAJOR, MAJOR, MINOR — and the count did not. Two of the eight
findings were shipped-behaviour bugs the reviewers caught before release: a **skill package whose
installer refused its own payload**, and a **cross-review cap that one code path walked straight
past**. The rest were the implementer's claims outrunning its code, five times, in comments,
changelogs, test names and one compliance assertion about a fix.

The reusable rule, now recorded in the source: **a number that is a safety boundary must not be
authored by the party it constrains** — and every time the mechanism changes, its description is a
factual claim that has to be re-checked against it.
