---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
reviewed-source: integration worktree `evidence-first-integration`, working tree as read on 2026-09-16
review-kind: supporting-source-only
---

# Blocked amendment cross-review: the human gate is unreachable, not merely ungranted

## Summary

- **The refusal is a first-binding problem, not a charging problem.** No cross-review
  cycle policy exists for this idea yet, so the launch path tries to *create* one, and
  creation enumerates every registered Git worktree. Two registrations are unavailable,
  so enumeration fails closed.
- **The three refused attempts spent nothing.** No cycle, no step, no child. Each left a
  requested/terminal pair with no start — the exact shape the importer classifies as
  `retained-pre-start-refusal`.
- **No existing supported attended operation can resume this.** Every documented route —
  `budget migrate inspect/apply`, `budget migrate recover`, `budget cycle extend` — funnels
  into the same enumeration or requires a binding that does not exist. `migrate apply`
  recomputes the inventory internally, so **an attended operator typing a correct
  `--expected-history-sha256` still fails.** This is the material finding: the gate is not
  waiting on a human; it is closed to one.
- **Operator input is indispensable, but no existing flag expresses what is needed.** The
  missing attestation is not a count. It is "these two registered paths are permanently
  unavailable and their history is unknown".
- **Three independent `agents exec` calls would spend three of the three cross-review
  charges** for one logical round, where the grouped runner spends one.

Findings: 0 CRITICAL, 1 MAJOR, 1 MINOR. This is source review within the scope below. It is
not Phase-6 acceptance, not a signoff, and it authors no amendment round or phase change.

## Scope and provenance

- **Launch attestation:** `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- **Method:** native Read/Grep/Glob over the working tree only. No shell, Git, build, test,
  CLI, or model call. I could not run `git worktree list` or `budget migrate inspect`.
- **PRIMARY (source):** file:line locator plus quotation, from files I opened. Test files are
  source I read, never runs I executed.
- **Unverified testimony (not mine, not owned by me):** the three `budget_refused` attempts,
  the observed `budget migrate inspect --kind cross-review` output, the two unavailable paths
  (`.../scratchpad/f2repo`, `/private/tmp/revert-test2`), and every timing/pass/fail in the
  coordinator notes. I did not reproduce any of them and assign them no verdict.
- **Self-owned (no self-verdict, per §15.1):** `claude-1-current-hardening-review-20260915.md`.
  I disposition its residuals below; I do not verify its claims.
- **Read in full:** `budget/{binding,cycle_binding,cycle_history,cycle_session,protocol_migration,
  protocol_migration_history,migration_recovery,step_binding,launch_migration_history}.go`;
  `runner/{cycle_budget,launch_budget,telemetry}.go`; `app/{agents_exec,budget_migrate,
  budget_migrate_protocol,budget_migration_recovery,budget_cycle,budget_origin,
  budget_attended_other}.go`; this idea's `00-prompt.md` and `FINAL.md`.
- **Read in part:** `runner/launch.go:150-250`, `budget/cycle_extension.go:74-93,163-181`,
  `track/track.go` (grep), `IMPLEMENTATION.md:1405-1446,3878-3893`.

## 1. The refusal chain (PRIMARY)

1. `agents_exec.go:113-115` builds `runner.LaunchInfo{... Phase: *phase ...}` and calls
   `runner.RunMeasured`.
2. `telemetry.go:89` `ctx, finishCycle = prepareLaunchCycle(ctx, root, info.Idea, info.Phase, info.RunID)`.
3. `cycle_history.go:21-24` — `if n, ok := cycleRound(phase); ok && n > 1 { return CrossReview, true }`.
   `round-02` is therefore a charged cross-review operation.
4. `cycle_budget.go:77` `LoadCycleBinding` uses `cycleScope(..., false)` (`cycle_binding.go:95`)
   and returns `nil, nil` for a missing scope dir (`:100-102`). **No worktree walk here.**
5. Because `b == nil`, `cycle_budget.go:40` calls `EnsureCycleBinding`, which at
   `cycle_binding.go:157` calls `cycleScope(ctx, root, idea, kind, true)` → `launchScope(..., true)`.
6. `binding.go:89` runs `git worktree list --porcelain -z`; `:101-104` is the failure:
   `if info, e := os.Stat(worktree); e != nil || !info.IsDir() {` →
   `err = fmt.Errorf("historical worktree is unavailable: %s", worktree)`.
7. `cycle_budget.go:60` `return budget.DeferCycleRefusal(ctx, kind, err), noop`.
8. `cycle_session.go:74-75` — `if err, ok := ctx.Value(cycleRefusalKey{kind}).(error); ok { return 0, err }`.
9. `launch_budget.go:70-75` wraps it; `telemetry.go:99-101` returns `nil, ...`; `telemetry.go:189-192`
   classifies `status, failure = "failed", "budget_refused"`.

The grep of `launchScope(`/`cycleScope(` call sites across `internal/` shows `binding.go:102` is the
**only** site raising that error, with callers at `binding.go:195` (`ConfigureLaunchBudget`),
`cycle_binding.go:157`, `step_binding.go:108`, and `launch_migration_history.go:136`.

## 2. What the three refusals actually spent: nothing (PRIMARY)

- The deferred error short-circuits `ChargeCycle` at `cycle_session.go:74-75`, **before**
  `s.binding.reserveWithReceipt(...)` at `:120`. No cross-review reservation exists.
- `JoinStepSession`/`ChargeStep` sit at `launch_budget.go:84-96`, **after** the cycle charge at
  `:70`. They were never reached. No step was spent.
- `recordLaunchRequest` runs at `telemetry.go:92`, before `reserveBudget` at `:99`, so each attempt
  did persist a requested record and a terminal with `StartedAt == nil` — consistent with the
  reported null `started_at`/PID.
- That shape is deliberately protected: `launch_migration_history.go:367-377` marks
  `classification = "retained-pre-start-refusal"` only when no start exists, and errors if any
  token, cost, exit code, or activity contradicts it; `:165-166` errors if a start later contradicts
  a pre-start refusal.

**The three refusals are cost-free and already correctly classified for any future import.** They
are evidence, not damage. Nothing needs to be "recovered" about them.

## 3. Can existing supported attended operations resume safely? No.

| Route | Source | Outcome |
| --- | --- | --- |
| `budget migrate inspect --kind cross-review` | `budget_migrate_protocol.go:46` → `InspectProtocolMigration` → `protocol_migration_history.go:80` → `InspectLaunchMigration` → `launch_migration_history.go:136` `launchScope(ctx, root, idea, true)` | Same worktree error. Matches the reported output. |
| `budget migrate apply --kind cross-review` | `protocol_migration.go:196` and `:235` both call `InspectProtocolMigration` | Fails identically. **The `--expected-history-sha256` value is recomputed, not trusted, so no operator input can satisfy it.** |
| `budget migrate recover inspect/apply --kind cross-review` | `migration_recovery.go:121-128` needs `readProtocolMigration(dir)`; no `migration.json` exists → NotExist. Even with one, `:472` calls `InspectProtocolMigration` | Dead twice over. |
| `budget cycle extend --kind cross-review` | `cycle_extension.go:180` `LoadCycleBinding` → nil; `InspectCycleBudget` returns `"cycle policy is not initialized"` (`:168-169`) | Cannot bootstrap a first binding. |
| grouped protocol runner (instead of manual exec) | `cycle_budget.go:71-81` `prepareLaunchCycle` → same `groupProtocolCycle` → same `EnsureCycleBinding` | Identical refusal. |

**[MAJOR] The documented attended recovery path for a cross-review binding is unreachable while a
registered worktree is unavailable.** Fail-closed is correct here; having *no* attended exit is not.

## 4. Is operator input indispensable, and what exactly would it attest?

Two separable things, and only one of them is a judgment call.

- **Attendance** is a kernel TTY check: `budget_attended_other.go:5-7`
  `return hasTTYSupported, hasTTYSupported && platformHasTTY()`. The code states its own limit —
  "this API cannot prove who made a decision" (`protocol_migration.go:165`), "participant metadata or
  a flag cannot supply attendance" (`budget_origin.go:52`). It is a mechanism-backed convention that
  stops an ordinary agent run, not authentication. I hold no position on whether it should be more.
- **The decision content** for `migrate apply` is genuinely operator-only: `--total-actions` (a total
  the tool cannot derive — "a separate explicit operator reconciliation", `protocol_migration_history.go:70`),
  `--started-at` (refused if later than observed history, `protocol_migration.go:40`), `--max-cycles`,
  `--decision-id`, `--reason`, `--writers-stopped`, `--yes`. `validCycleDecision`
  (`cycle_extension.go:74-77`) requires a nonempty id ≤128 chars, nonempty reason ≤1024, and a
  lowercase 32-byte hex digest.

So yes, operator input is indispensable by design. But **the attestation currently missing is not any
of those fields.** It is: *"registered worktree X is permanently unavailable; its history is unknown
and must stay unknown."* No flag expresses that, which is why the human gate is closed.

**Can we prepare a reviewable record first? Yes, and we should.** Everything except the final attended
apply is preparable now: exact command lines, the two registered paths verbatim, the reason text, a
decision id, and the predicted floor. The one field we cannot pre-compute is the history hash — and
that is the circular dependency itself.

## 5. Independent `agents exec` vs one grouped cross-review charge

- Grouped: `cycle_budget.go:62` opens one `OpenCycleSession`; `ChargeCycle` memoizes on
  `s.attempted` (`cycle_session.go:94-113`) and returns the same `s.ordinal` without a second
  reservation. One round = one charge, regardless of participant count.
- Manual: each `agents exec` is a separate process, so `beginLaunch` → `prepareLaunchCycle` opens a
  **fresh** session. There is no cross-process session. Three launches = three charges.
- The cap for this idea: `00-prompt.md:5` `track: deliberation`; `track.go:196-198`
  `CrossReviewRounds: -1`, `CapCrossReviewRounds: 3`; `cycle_budget.go:191-201` therefore yields
  `cross = 3`. The deck already records the same semantics: "Nested children of one synchronous
  operation share one charge; independent operations spend separate charges… The fifth deliberation
  fixup and third cross-review group are allowed; the sixth/fourth are refused"
  (`IMPLEMENTATION.md:1436-1440`).

**[MINOR] Running amendment round-02 as three independent manual launches would consume the entire
cross-review allowance on one logical round.** If the binding is later created with the floor derived
from existing `round-01`/`round-02` directories, `LegacyCycleFloor` (`cycle_history.go:90-94`,
`count = n-1`) yields carried = 1, leaving 2 — not enough for three independent launches. Recommend
the grouped runner for the amendment round, or an explicit operator decision that knowingly sets the
ceiling for per-participant charging. I did not read `cycleLimits` or `Store.Reserve`, so the exact
inclusive boundary above rests on `track.go` plus the recorded note, not on my own reading of the
comparison.

## 6. Proposed minimal recovery specification (MRW-1) — proposal only, no code, no acceptance

Modeled directly on the existing `budget origin inspect|apply` precedent (`budget_origin.go`), which
is already an attended, `--expected-sha256`-bound, decision-id'd repair of a *location* fact.

1. **`parley budget worktree inspect --dir DIR`** (read-only, unattended-allowed). Enumerates
   registrations and returns them **partitioned** into available/unavailable with each registered
   path verbatim, plus a digest over that partition. It must not prune, recreate, or drop any row.
2. **`parley budget worktree attest apply --dir DIR --unavailable PATH [--unavailable PATH] \
   --expected-sha256 H --decision-id ID --reason TEXT --writers-stopped --yes`** (attended). Writes
   a durable record naming each path and marking its history **`unknown`** — never a count, never zero.
3. **`launchScope(..., true)`** consults that record: an exactly-matching path is skipped from `roots`
   **and** sets an `incompleteHistory` flag on the returned scope.
4. **Ordinary bootstrap stays refused under `incompleteHistory`.** `EnsureCycleBinding`,
   `ConfigureLaunchBudget` and `EnsureStepBinding` must still reject a zero-carried start and demand
   the explicit migration path. The attestation converts "cannot enumerate" into "enumerated, with a
   declared unknown region" — it never converts unknown into zero.
5. **`InspectProtocolMigration` records the attested paths inside the inventory**, so they enter
   `HistorySHA256` and the immutable migration record, and `protocolMigrationInitial`'s
   `r.TotalActions < i.LowerBound` check (`protocol_migration.go:40`) keeps the operator's count
   carrying the unknown region.
6. **The frozen cap is untouched.** `cycle_binding.go:145-147` (`InitialMaximum() != maximum`) and
   `cycle_budget.go:47-49` (stricter-track refusal) continue to bind unchanged.

This touches one enumeration site, one new record, one CLI verb pair, and reuses the existing
attendance primitive, decision shape and replay discipline. It adds no way to lower a count and no
way to raise a cap.

## 7. Alternatives disposition

- **ALT-A `git worktree prune`** — REJECT. Destroys the registration and silently turns unknown
  history into no history. Also excluded by the task.
- **ALT-B recreate the missing directories** — REJECT, and the most dangerous option:
  `refuseUnmigratedCycles` would then read an empty history and bootstrap `carried: 0` *successfully*.
  A synthesized checkout looks like recovery while asserting zero.
- **ALT-C run from a non-Git root / independent clone** — REJECT. `binding.go:122-130` deliberately
  refuses a fresh local ledger when any ancestor `.git` exists; a clone resets accounting. Excluded.
- **ALT-D relabel the phase** (e.g. `--phase manual`) — REJECT, and I name it because it **would**
  mechanically work: `cycle_budget.go:73-74` returns `ctx, func(){}` when `CycleKindForPhase` is
  false, so no cycle is charged at all. An undocumented working bypass is worse than a documented
  refused one. Excluded by the task and by D6.
- **ALT-E grouped runner instead of three manual execs** — ADOPT *for accounting* (§5), REJECT as an
  unblocker: it reaches the identical refusal.
- **ALT-F `budget cycle extend`** — REJECT. Requires an existing binding (`cycle_extension.go:180`).

## 8. Dispositions of my prior residuals

- **N1 (cross-review/step precharge).** Narrowed, not closed. New PRIMARY: a *deferred* cycle refusal
  spends nothing (§2). The residual is unchanged where a binding exists — `reserveBudget` charges the
  cycle (`launch_budget.go:70`) and the step (`:84-96`) **before** `LoadLaunchBinding` (`:98`),
  `RequireMonetaryBinding` (`:111`) and the reservation (`:139-141`). A launch refused on monetary or
  binding grounds therefore still spends a cycle and a step without starting. Charge-before-work is
  D6's intent; the narrow question is whether *known-refusable* checks should precede the charge. Left
  open as previously filed, not re-filed.
- **R3 (helper scope window).** **Not implicated.** The blocked path never reaches
  `trajectory_verify.go` or the verifier ticket. R3 stands unchanged, still optional.

## 9. Evidence boundaries

- Not read: `cycleLimits`, `Store.Reserve`, `refuseUnmigratedSteps`, `lock_origin_migration.go`,
  `trackedCommandFor`. Claims depending on them are flagged above.
- Not established: that the two named paths are the *only* unavailable registrations; that `--idea`
  for the amendment resolves to this idea's directory rather than a v2 slug (which changes the
  `LegacyCycleFloor` carried value); any provider, timing or cost claim.
- Unverified premise (RECALL): that `git worktree list` still reports both registrations. I could not
  run it.

## 10. Next actions (proposed, not decided)

1. Prepare the full attested-unavailability record now — paths, reason, decision id, predicted floor —
   so an operator reviews a finished artifact rather than composing one.
2. Escalate MRW-1 (§6) as the minimal change; it needs an owner and a non-owner reviewer.
3. Decide grouped-vs-manual for amendment round-02 **before** any binding is frozen (§5).
4. Keep the three refusal records exactly as they are; they are the cleanest evidence in this thread.
