---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: supporting-source-design-review
status: review-only
reviewed-source: integration worktree `evidence-first-integration` working tree as read 2026-09-16 (internal/budget, internal/runner/runner.go, internal/runner/cycle_budget.go); plan under review: claude-1-history-admission-design-20260916.md
method: native Read/Grep only; no shell, Git, build, test, CLI or model call
---

# Review: claude-1 MRW-1 history-admission design (bounded, five questions)

Boundaries: source reads only; nothing executed. That the blocked operation is v2 round-02 in sibling `evidence-first-amendment` is task-supplied context; the two unavailable paths and the three refusals remain codex-1 testimony I did not verify. No signoff, no whole-audit acceptance. Backup question and the count decision stay unanswered.

## Q1 — wrong-idea binding: CONFIRMED, owner correction required (PRIMARY)

Scope identity is per-idea: `scope = "parley-launch/v1:" + key(prefix+"\x00"+idea)` (binding.go:132), with the cycle scope adding kind (cycle_binding.go:44). A policy stored for `meta-protocol-change-evidence-first-efficiency` is never loaded for `...-v2`: the store dir is the scope hash and `LoadCycleBinding` rejects mismatches (`"cycle budget scope mismatch"`, cycle_binding.go:113-115). With no v2 policy, the grouped round re-enters `EnsureCycleBinding` (the `b == nil` branch, cycle_budget.go:39-40) → `cycleScope(ctx, root, idea, kind, true)` (cycle_binding.go:157) → the same `"historical worktree is unavailable"` refusal (binding.go:101-104). Applying as written would also freeze an immutable decision — wrong evidence floor, `idea_path`, `earliest`, `TotalActions` — onto a scope nobody is blocked on. Correction: `--idea meta-protocol-change-evidence-first-efficiency-v2` and `--dir` = the amendment worktree root (the blocked runner's root, guaranteeing identical git-common-dir and prefix), for inspect and apply alike. No command has been run; this is a plan-text correction.

## Q2 — one migration suffices for the grouped round only; the permanence sentence is overbroad (PRIMARY)

Sequencing in the blocked path: `RunRound` delegates to `RunRoundOne` (runner.go:925-933), which opens `budget.GroupStepSession` FIRST (runner.go:203) and `groupProtocolCycle(..., budget.CrossReview)` SECOND (runner.go:209-213). The step session is load-only: `JoinStepSession` → `LoadStepBinding` → `stepScope(ctx, root, idea, false)` (step_binding.go:62), and a nil binding is a free no-op (step_session.go:29-33). In production only the driver calls `EnsureStepBinding` (internal/driver/loop.go:39, internal/driver/budget.go:22). So round-02 scans history exactly once — in the cycle grouping — and ONE cross-review apply does unblock this operation; Claude's after-apply mechanics are exact: the walk runs only under `if inspectHistory {` (binding.go:88), `LoadCycleBinding` uses `cycleScope(..., false)` (cycle_binding.go:95), and `checkProtocolMigrationPolicy` re-reads only stored files (protocol_migration.go:97-125). But "unblocks the idea permanently; no later launch depends on the two roots" fails for first-time activations of the other two scopes sharing the identical scan: driver-step activation (`stepScope(..., true)`, step_binding.go:108) and launch configure (`launchScope(..., true)`, binding.go:195 — the monetary-gate scope per my owned audit kimi-1-experiment-cap-source-audit-20260916.md). Each is idea×kind-scoped; each would need its own declared import and separately scoped count decision, or restored roots. Cycle and step authority are NOT one decision.

## Q3 — request-scoped design preserves bootstrap; what the digest must bind (PRIMARY + design)

Ordinary bootstrap stays refused iff the three history-scanning call sites pass nil declarations: binding.go:195, cycle_binding.go:157, step_binding.go:108 (all `true` today); the wrapper claim is implementation work, not yet verifiable. Old records remain readable: `migrationShape` tolerates absent `omitempty` fields (launch_migration_history.go:236-242). The reverse is NOT compatible: `DisallowUnknownFields` (:193-195) makes pre-change binaries reject records carrying the new fields — this one-way direction must be recorded before apply. For stale-rejection the digest must bind: (a) verbatim porcelain path strings; (b) observation class — binding.go:101-104 conflates stat-error and non-dir today, so the declaration must pin enoent vs not-directory; (c) continued registration at apply (absent from porcelain ⇒ stale registry); (d) continued unavailability (a real directory reappearing ⇒ refuse; never migrate over now-readable evidence); (e) exclusion from `Roots`/`Sources` so no evidence is claimed from declared paths; (f) `History == "unknown"` pinned by a new `protocolMigrationInitial` clause (beside protocol_migration.go:37) — never 0, never a count; (g) sorted/deduped canonical form so `DeepEqual(prior.Request, r)` replay (:185, :214) cannot split nil from `[]string{}`. Non-ENOENT stat errors remaining non-declarable hard errors is correct. This preserves my stage-one invariant: directory-exists-only means no writer may infer verified history.

## Q4 — replay / partial-publication / recovery claims: supported by actual helpers (PRIMARY)

Verified against code, not function names: exact-replay-vs-conflict via `reflect.DeepEqual(prior.Request, r)` (:185, :214-215); serializing guard plus post-guard re-read (:204, :209-216); `migration.json` written only on first import (:243, :249-259); ledger update rejecting non-empty entries (`"concurrent charges prevent protocol import"`, :265-268); inactive-policy byte equality (:286-316); stale-history veto `errHistoryChanged` (:317-323); activation witness persisted LAST (:324) and enforced as `"protocol migration is not durably active; replay the exact decision before work"` (:120-123) — a crash before :324 leaves an inactive policy that the exact replay completes, with files and counters preserved. The recovery journal exists and chains (`protocolRecoveryBase`, migration_recovery.go:98-104; chain validation :142-182; `recoveredMigrationState` :282-294; `refusePendingRecovery` :296-303). Claude's recovery gap is real: an inactive import's preview re-inspects protocol history (`InspectProtocolMigration` at migration_recovery.go:472), which re-scans — the declared list must thread through that call plus all three migrate-side sites (protocol_migration.go:196, :235, :317), or a declared import can never be recovered.

## Q5 — minimum correction set (no cap increase, no unknown-as-zero, no prune/recreate/relabel/clone)

1. Retarget `--idea` + `--dir` per Q1; re-collect the stage-one inventory against the amendment root if the registered paths differ there.
2. Replace the permanence sentence per Q2 with per-scope sufficiency plus the two named future gates (driver-step activation, launch configure).
3. State pre-apply requirements: the v2 cycle scope dir holds no `policy.json`/`ledger/ledger.json`/`migration-active` (first-import precondition, protocol_migration.go:191-195); `--max-cycles` consistent with v2's own 00-prompt track (`"current track is stricter than the frozen cycle policy"`, cycle_budget.go:47-49); the visible floor computed from v2 evidence — `LegacyCycleFloor` counts an existing round-02 as ONE completed prior round (`count = n - 1`, cycle_history.go:90-94), and the three refusals add nothing while classified `retained-pre-start-refusal` (excluded at protocol_migration_history.go:195-199).
4. Show the operator the arithmetic: the ceiling is `Maximum − Carried` (cycle_intent.go:29-33) and `Store.Reserve` refuses at `count >= ceiling` over same-kind entries (ledger.go:168-176); with `--max-cycles 3`, N must satisfy N ≥ LowerBound (protocol_migration.go:40) AND N ≤ 2 for round-02 to reserve. N itself, the two verbatim paths, and their observation class are operator knowledge — not fabricable; backup availability is unanswered.
5. Sequence: inspect (read-only) → human N decision → attended apply → grouped round. Never the round first.
6. Record as accepted limitations before approval: the one-way compatibility (Q3) and Claude §8's `omitempty` tamper-evidence caveat.

## Main blockers

1. Wrong scope target (Q1) — blocks any apply as written; plan-text owner correction.
2. Overbroad permanence claim (Q2) — blocks the operator-approval wording; cycle, step and launch scopes are separate decisions.
3. Count N, the two verbatim paths/observations, and the backup answer — operator inputs, outstanding; not establishable from source.

No test or execution claim; no signoff; no acceptance of the whole audit.

**Own path:** `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/kimi-1-history-admission-design-review-20260916.md`
