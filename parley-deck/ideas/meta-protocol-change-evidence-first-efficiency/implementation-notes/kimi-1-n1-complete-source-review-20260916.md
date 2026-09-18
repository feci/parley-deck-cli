# kimi-1 — N1 complete precharge slice: independent bounded source review (2026-09-16)

**Allocation:** this note only. No source edits, no Bash/tests/Git, no subagents/provider calls.
**Source under review (PRIMARY):** immutable snapshot
`/private/var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-n1-combined-lu066y5v/kimi-app-corrected-source`
(451 Go/module files per launch task; bound by `source-manifest.json` in the integration runtime
dir — testimony, not independently re-hashed here). Root integration is recovery-only1041259;
root `internal/` was NOT consulted, per instructions.
**Protocol context:** attestation supplied with this launch, `context_mode=full`,
`source_sha256 == packet_sha256` (4519…1937a). All verdicts below are PRIMARY (files read
directly, locators quoted) unless marked testimony.
**Testimony (not re-run; Bash prohibited):** coordinator reports full affected driver/runner/app +
race/vet PASS, 28 events / 15.577s, including my corrected app tests; prior failure retained.
I did not execute anything; source is the primary evidence.
**Provenance:** original N1 author zcode-1 timed out; coordinator owns the driver/runner slice,
fixtures and skip correction; kimi-1 owns only the app test corrections (prior notes
`kimi-1-n1-app-review-tests-20260916.md`, `kimi-1-n1-app-fixture-correction-20260916.md`).
This review evaluates the newly found combined scope independently.

## Production entrypoints traced

- Driver: `driver/driver.go:276` (`withStepBudget` per Advance) → `driver/budget.go:12-41`
  (monetary preflight at :13-21, `EnsureStepBinding` :22, session :26) →
  `stepRoundRunner.RunRound` `driver/budget.go:49-57` (precheck :50 → `ChargeStep` :53 →
  nested runner :56) → `cycleRoundRunner.RunRound` `driver/cycle_budget.go:48-70`
  (`EnsureCycleBinding` :57 → `OpenCycleSession` :61 → `ChargeCycle` :66 → real launch :69).
- Driver precheck seam: `driver/launch_precheck.go:59-67` (unwraps exactly one
  `cycleRoundRunner`, consults the seam once, adapters without it pass through nil),
  `:48-54` (adapter mirrors `RunRound` normalization, `Overwrite=false`).
- Runner: `runner/launch_precheck.go:11-75` (typed prechecks + `precheckSelectedLaunch`);
  `runner/protocol_context.go:135-145` (`PrecheckProtocolLaunch`) → `:147-159`
  (`retainProtocolRefusal` → `recordLaunchRequest`, never `beginLaunch`);
  real boundary `beginProtocolLaunch` `:116-129` → `beginLaunch`
  `runner/telemetry.go:83-112` → `reserveBudget` `runner/launch_budget.go:71-158`.
- Fixup: `driver/impl.go:357` (`PrecheckFixup`) → `:365` `reserveFixupCycle`
  (`driver/cycle_budget.go:74-89`: `EnsureCycleBinding` :76 → session :80 → `ChargeStep` :84
  → `ChargeCycle` :87); runner twin `runner/phase58.go:55-64`.

## Mandated refutation targets — results

1. **Step/cycle ordering (driver): holds.** Every wrapper consults the precheck before
   `ChargeStep` (`driver/budget.go:50-53, 62-65, 71-74, 80-83, 98-101, 107-110, 116-119,
   125-128, 146-149`); the cycle charge is nested inside the step charge
   (`driver/cycle_budget.go:57-67`). `Reopen`/`Fixup`/`RunChecks`/`Complete` deliberately
   unprechecked (`driver/launch_precheck.go:34-37`). `VerifyCompletionEvidence` checks
   adapter capability before charging (`driver/budget.go:154-162`). Pinned by
   `driver/launch_precheck_test.go:123-170` (`assertNoStepCharge`) and `:176-258`.
2. **Monetary known-refusal preflight: holds, fail-closed in both entrypoints.** Driver:
   `driver/budget.go:13-21` before any binding/session. Runner:
   `runner/launch_budget.go:76-91` — binding load error refuses (:77-82), unreadable
   defaults refuse with a generic diagnostic that copies no config contents (:83-88),
   `RequireMonetaryBinding` (`budget/monetary_binding.go:39-56`: zero ceiling = no new
   policy and never disables a persisted one; nil binding + ceiling → `ErrUnknownCost`;
   ceiling must match frozen policy history) refuses before `ChargeCycle` at
   `launch_budget.go:113`. Pinned by `runner/launch_precheck_test.go:87-168` (zero
   step/cycle entries, one unstarted `budget_refused` terminal; restore → exactly 1+1).
3. **Late failure stays spent / no refund: holds.** Charges precede the launch in both
   entrypoints; `settleBudget` (`runner/launch_budget.go:164-181`) only settles, never
   refunds, and runs even on deadline/publication failure (:160-163). Driver test
   `launch_precheck_test.go:221-257` (failed real child: step=1, cycle=1, stays spent;
   incomplete round awaits without re-charge); `cycle_budget_test.go:22-72` (failed fixup
   survives new run + cursor deletion, no refund).
4. **Skip/overwrite mirror: holds.** Precheck predicate
   `runner/launch_precheck.go:66` is byte-equivalent to `runAgent` `runner/runner.go:406`
   (`os.Stat` success && `!Overwrite` → skip; stat errors proceed to render/launch).
   Review-consensus forces `Overwrite=true` in both (`launch_precheck.go:46`;
   `phase58.go:363` per my prior note). `firstOnly` truncation before the skip loop
   matches `selected[0]`-only real entries. Pinned by
   `runner/launch_precheck_skip_test.go:13-79` incl. the real-`RunRound` cross-check and
   byte-untouched artifacts, and my `app/driver_precheck_test.go:130-144`.
5. **Already-bound fixture vs unbound migration: fixture valid; production gap found —
   see MAJOR-1.** The corrected driver fixture pre-binds
   (`driver/launch_precheck_test.go:196-198`, comment :193-195) and is a sound
   already-bound test. Restoration does NOT heal an unbound scope (confirmed, and the
   launch task warned not to infer it).
6. **Shared grouped count: holds.** Driver and manual runner share one durable charge via
   the same frozen binding and nested-session reuse (`budget/cycle_session.go:45-57`,
   `budget/step_session.go:35-47`); `driver/cycle_budget_test.go:165-234` (driver child
   fails, manual `RunFixup` budget-refused, `Count==1`, terminals=2/started=1/denied=1).
7. **Nil/empty selections: holds.** `selectedAgents` (`runner/runner.go:378-387`) yields
   empty → precheck loop no-ops to nil (`runner/launch_precheck.go:59-60,74`);
   unresolved-participant noop pinned by `runner/launch_precheck_test.go:69-81`.
   `PrecheckFixup` errors on empty selection (`runner/phase58.go:57-59`) exactly like the
   real path.
8. **Unknown binding/default errors: holds, fail-closed.** `LoadLaunchBinding`
   (`budget/binding.go:137-162`) errors on unreadable/lost policy, nil only when absent;
   `LoadStepBinding`/`LoadCycleBinding` same shape. All propagate to refusal before
   charging in both entrypoints.
9. **Malformed explicit policy: holds.** With a frozen binding, a non-matching explicit
   `LaunchBudget` refuses (`runner/launch_budget.go:92-107`); the comparison covers store,
   full limit history and reserve micros. Unbound explicit policy is documented trusted
   orchestration input (:16-18) validated at `Store.Reserve`.
10. **Duplicate precheck refusal receipts: accumulation confirmed — NIT-1; poisonous in
    combination with MAJOR-1.**
11. **Parent recovery impact: no defect found.** The precheck path only calls
    `recordLaunchRequest` → `telemetry.Begin` (fresh ID; `runner/telemetry.go:115-151`);
    the bound-recovery exclusive allocation (`recordBoundLaunchRequest` via
    `boundRecoveryInvocation`, `telemetry.go:92-101`) is unreachable from
    `PrecheckProtocolLaunch` (`runner/protocol_context.go:135-145`), so a pinned recovery
    ID is never consumed or duplicated by a precheck. Trajectory capture
    (`runner/launch_budget.go:120-126`; observer `driver/cycle_budget.go:75`) attaches
    only on real fixup charges; a precheck refusal creates no trajectory run and no
    observer interaction. Residual: a `trajectory.Begin` failure after a successful
    `ChargeCycle(Fixup)` leaves a charged cycle without capture — pre-existing, inside the
    documented non-atomicity (`launch_budget.go:64-70`), folded into MINOR-1's fix.

## Findings

### [MAJOR-1] N1 precheck refusal receipts poison FIRST cycle binding on an unbound scope; every retry also burns a step

**Where:** interaction of `runner/protocol_context.go:147-159` + `runner/telemetry.go:115-151`
(receipt written with the real phase) with `budget/cycle_history.go:178-205` (first-binding
invocation scan; :198-201 errors on any matching idea + cycle-classified phase) reached from
`budget/cycle_binding.go:122-205` (`old == nil` → `refuseUnmigratedCycles` at :161/:187).

**Evidence chain (all PRIMARY):**
1. Driver order on a fresh idea: precheck (`driver/budget.go:50`) runs BEFORE
   `EnsureCycleBinding` (`driver/cycle_budget.go:57`). The step binding is created earlier
   (`driver/budget.go:22`), but no cycle binding exists yet.
2. A refusing precheck retains `.parley-runtime/invocations/<fresh-id>/requested.json` +
   `terminal.json` with `Metadata.Idea` and `Metadata.Phase` = e.g. `round-02`
   (`runner/protocol_context.go:147-158` → `runner/telemetry.go:127-151`).
3. After authority is restored, the next Advance passes the precheck, charges a step
   (`driver/budget.go:53`), then `EnsureCycleBinding` finds no binding and scans history:
   the receipt's phase classifies (`budget/cycle_history.go:17-25`: `round-NN`, n>1 →
   CrossReview; `fixup` → Fixup) and matches the kind → hard error "historical cycle
   requests require explicit accounting migration before first binding"
   (`cycle_history.go:200`), surfaced as "cross-review accounting: …"
   (`driver/cycle_budget.go:58-60`).
4. Every retry repeats step 3: one more durable step spent, binding still never created,
   until `MaxDriverSteps` exhausts and escalates. Fixup flavor identical: `PrecheckFixup`
   (`driver/impl.go:357`, phase `fixup` via `runner/phase58.go:60-63`) precedes
   `reserveFixupCycle`'s `EnsureCycleBinding` (`driver/cycle_budget.go:76`).

**Counterexample (deterministic trace):** fresh idea, `MaxDriverSteps=5`, cross-review cap
unbound, tampered protocol cache → Advance refuses for free (correct); operator restores
the cache → Advance: step 1/5 spent, then `EnsureCycleBinding` errors; three more Advances
spend steps 2–4 and fail identically; the idea can never cross-review until an operator
runs exact-replay migration against a record that charged nothing.

**Why this is N1-introduced:** pre-N1 the driver ensured the binding before any receipt
existed (`driver/cycle_budget.go:57` precedes the launch); the manual runner ensures it in
`prepareLaunchCycle` → `groupProtocolCycle` (`runner/cycle_budget.go:40,71-81`) before
`beginProtocolLaunch` renders (`runner/telemetry.go:89,108`). Only the N1 seam retains a
receipt before first binding. The corrected fixture dodges the window by pre-binding
(`driver/launch_precheck_test.go:196-198`) — production never does. Note the guard itself
is correct for REAL retained requests; the defect is that the FREE path manufactures them.

**Fix (preferred):** hoist charge-free binding creation ahead of the precheck — for
round > 1 ensure the CrossReview binding (same `cap`/`floor` values `cycleRoundRunner`
already computes at `driver/cycle_budget.go:52-57`) before `precheckRoundLaunch`, and move
`EnsureCycleBinding(Fixup)` above `PrecheckFixup` at `driver/impl.go:357`. This runs the
legitimate migration scan against pre-N1 history earlier (strictly safer) and lands any
later receipt in an already-bound scope, matching the fixture. Alternative: mark precheck
receipts in `requested.json` and exempt unstarted zero-charge refusals in
`budget/cycle_history.go:178-205` — more invasive; late real refusals DID charge and must
keep counting, so the mark must be unforgeable by the launch path.

### [MINOR-1] `reserveBudget` charges the cycle before the step, so step exhaustion spends a cycle with no launch

**Where:** `runner/launch_budget.go:113-141` — `ChargeCycle` :114, fixup `trajectory.Begin`
:120-126, `JoinStepSession` :128, `ChargeStep` :136.

**Counterexample:** step binding `MaxSteps=1`, fixup binding `Maximum=4`, count 0. Manual
`RunMeasured` fixup launch #1 spends 1+1. Launch #2: `ChargeCycle` succeeds (count 2/4),
then `ChargeStep` refuses (1/1) → `budget_refused`, no process, cycle count stays 2 with
one actual launch. Same shape if `trajectory.Begin` fails at :121-125 after the charge.
The driver nests in the opposite order (step `driver/budget.go:53` before cycle
`driver/cycle_budget.go:66`), so the waste is symmetric there (cycle exhaustion burns a
step); both are inside the documented non-atomicity (`launch_budget.go:64-70`), but
exhaustion is knowable read-only in advance.

**Fix:** extend the existing read-only preflight block (`launch_budget.go:76-112`) with a
non-mutating ceiling check of the step binding and the cycle binding (a `Check`-style
inspect) so known exhaustion refuses before ANY charge; optionally also reorder
step-before-cycle to match the driver nesting. No refund semantics needed.

### [NIT-1] Duplicate precheck refusal receipts accumulate unboundedly

Each precheck refusal allocates a fresh invocation ID (`runner/telemetry.go:115-116` →
`recordLaunch` → `telemetry.Begin`); repeated Advance polls while authority is broken
retain one unstarted receipt per poll with no dedup or throttle. Today each is an
independent MAJOR-1 migration blocker; after that fix, the growth is audit-trail noise.
**Fix:** after MAJOR-1, document one-receipt-per-refusal as intended, or dedup by
(scope, phase, source-sha) in `retainProtocolRefusal`.

### [NIT-2] `errors.New(...).Error()` round-trip

`driver/budget.go:157` builds an error solely to stringify it. Harmless; use a plain
string constant. Recorded, not gating.

### [NIT-3] Signoff precheck seam widens beyond render refusals

`app/driver_precheck.go:88-99` returns `consensus.Status`/discovery/selection errors
through the free precheck on a determinism rationale (:86-88). Config-shaped causes are
deterministic; transient IO errors are not — they would abort the transition free and
receipt-less (safe direction; the real call would fail identically after spending a step).
Recorded for completeness; no change required.

## Scope-limited conclusion

Within the bounded slice, the N1 precharge design is implemented faithfully for the
already-bound case: precheck-before-charge ordering, fail-closed monetary/binding/defaults
preflight, exact skip/overwrite mirroring, spent-stays-spent late failures, and shared
grouped counting all verified against production entrypoints with tests that pin them.
One MAJOR interaction defect (MAJOR-1) and one MINOR accounting-order waste (MINOR-1)
found, each with a deterministic counterexample and a concrete fix; three NITs recorded
without suppression. My prior app/skip findings stand; the coordinator's skip correction is
independently re-confirmed faithful here.

This is NOT a whole-audit acceptance and NOT a claim that the source is published: the
snapshot is an immutable input bound by the coordinator's manifest (testimony), and the
coordinator's PASS is testimony I could not re-run under the no-Bash constraint.

## Residual coverage (not reviewed; do not read as accepted)

- Renderer internals (`protocolpacket`/`protocolcore`), `telemetry.Invocation`/`Begin`,
  store ledger `Reserve`/`Settle` validation internals, `budget` launch/protocol migration
  machinery beyond the cited scans, `trajectory` package internals, `track.PolicyFor`,
  `config.LoadDefaults`.
- App helper wiring (`withParticipants`, `firstHeadlessAgent`, `requestSignoffTargets`,
  `rosterMappingFor`) read only at call sites.
- Untested app branches per my prior note: non-review `PrecheckConsensusSignoffs`,
  `PrecheckReviewLaunch`, `PrecheckReviewConsensus`, `PrecheckConsensusFinal`.
- race/vet and the full 451-file module: coordinator-owned runs (testimony).
- MAJOR-1's fix verification belongs to a follow-up: a production-order fixture (fresh
  idea, tamper → refuse → restore → Advance) proving no forced migration and no step
  drain.
