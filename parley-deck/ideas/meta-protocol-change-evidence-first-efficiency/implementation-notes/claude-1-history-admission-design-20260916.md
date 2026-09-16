---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: source-only-design
status: proposal-only
amends: claude-1-recovery-source-review-20260916.md (SELF-CORRECTION, section 0)
reviewed-source: integration worktree `evidence-first-integration`, working tree as read 2026-09-16
method: native Read/Grep only; no shell, Git, build, test, CLI or model call
---

# MRW-1: admitting declared-unknown history through the existing migration

**Proposal only — not implemented, not reviewed, not an attestation, not a signoff.** Every count
below is a placeholder. Implementation and independent review are still required.

## 0. SELF-CORRECTION to my own prior note

My earlier §6 step 3 had `launchScope(..., true)` consult a durable attestation record. I withdraw it: the same function backs `ConfigureLaunchBudget` (`binding.go:195`), `EnsureCycleBinding` (`cycle_binding.go:157`) and `EnsureStepBinding` (`stepScope(..., true)`, `step_binding.go:108`, wrapping `launchScope` at `:35`), so a durable record would quietly weaken ordinary bootstrap for every future idea and kind in the repo. Weakening of a self-owned claim, effective immediately (§15.1). The declaration must be **request-scoped** — an argument to one decision, never repository state.

## 1. Decision: no separate attestation writer

**Existing migration should take one typed operator decision.** `migrateProtocolBudget` already owns every property a writer would need: exact-replay-vs-conflict (`reflect.DeepEqual(prior.Request, r)`, `protocol_migration.go:185`, `:214`), an immutable digested record (`:249-261`), a serializing guard (`:204`), a stale-history veto (`:317-323`), an activation witness (`:324`, enforced at `checkProtocolMigrationPolicy:120-124`) and a recovery journal. A second writer rebuilds all of that — and would then need a *reader*, whose only candidate is `launchScope`, i.e. §0's mistake.

Second PRIMARY fact making one decision sufficient: the worktree walk runs only under `inspectHistory == true` (`binding.go:88`). Afterwards `LoadCycleBinding` uses `cycleScope(..., false)` (`cycle_binding.go:95`) and `checkProtocolMigrationPolicy` re-validates the *stored* inventory digests without re-enumerating. One attended apply unblocks the idea permanently; no later launch depends on the two roots existing.

## 2. Fields (proposed, exact)

```go
// binding.go — request-scoped; never persisted as repository state.
type UnavailableRoot struct{ Path, Observation, History string } // json: path/observation/history
//   Path = verbatim `git worktree list --porcelain` value; Observation = stat-enoent|stat-not-directory
//   History = always "unknown" — never a count, never 0
// LaunchMigrationInventory   += UnavailableRoots []UnavailableRoot `json:"unavailable_roots,omitempty"`
//                            += HistoryCoverage  string `json:"history_coverage,omitempty"`  // "declared-incomplete"
// ProtocolMigrationInventory += LowerBoundBasis  string `json:"lower_bound_basis,omitempty"` // "surviving-visible-worktrees-only"
// ProtocolMigrationRequest   += DeclaredUnavailable []string `json:"declared_unavailable,omitempty"`
```

All `omitempty`, so an already-applied migration keeps reading (`migrationShape` demands every non-`omitempty` tag, `launch_migration_history.go:236-242`); absent ⇒ no declaration ⇒ today's behavior. `Roots` is **not** extended — an unavailable path is never walked and contributes no evidence. `LowerBound` keeps its name; `LowerBoundBasis` is what distinguishes the surviving-visible floor from the operator's charged `TotalActions`.

New `launchScopeDeclared(ctx, root, idea, inspectHistory, declared)` also returns `missing []UnavailableRoot`; `launchScope` becomes a wrapper passing `nil`, leaving all existing call sites behaviorally identical. `binding.go:101-104` becomes: non-`ENOENT` stat error → hard error (never declarable); `ENOENT`/non-dir → hard error **unless** exactly declared; a declared row that is now a real directory → new hard error *"declared-unavailable worktree is available again"*; a declared path absent from the porcelain output → hard error. Declarations are sorted and deduped in `migrateProtocolBudget` beside `r.StartedAt = r.StartedAt.UTC()`, so `nil` and `[]string{}` cannot diverge under `DeepEqual`.

## 3. Exact command sequence

```
# 1. read-only preview — unattended allowed, writes nothing, spawns no agent. Emits history_sha256,
#    lower_bound(+basis), history.unavailable_roots[], history.launches[], evidence[], earliest.
parley budget migrate inspect --kind cross-review --dir <INTEGRATION-ROOT> \
  --idea meta-protocol-change-evidence-first-efficiency \
  --declare-unavailable-worktree '<PATH-A>' --declare-unavailable-worktree '<PATH-B>'
# 2. human reads that preview and decides --total-actions (§6). No command.
# 3. attended guarded apply — same declarations, digest from step 1.
parley budget migrate apply --kind cross-review --dir <INTEGRATION-ROOT> \
  --idea meta-protocol-change-evidence-first-efficiency \
  --declare-unavailable-worktree '<PATH-A>' --declare-unavailable-worktree '<PATH-B>' \
  --expected-history-sha256 <H-FROM-STEP-1> --total-actions <N> --max-cycles 3 \
  --started-at <RFC3339, not after inventory.earliest> \
  --decision-id <UNIQUE> --reason '<basis>' --writers-stopped --yes
```

Only after step 3 exits 0 may the grouped round start. The flag is rejected for `--kind launch` (alongside the protocol flags at `budget_migrate.go:54-58`) and for every other verb.

## 4. Invariants → holding mechanism

| Invariant | Mechanism |
| --- | --- |
| Rows retained, history unknown | `UnavailableRoots` inside `HistorySHA256`; `History` pinned to `"unknown"`, validated in `protocolMigrationInitial` |
| Visible floor ≠ charged count | `LowerBound` from `protocolEvidenceFloor` (surviving roots only) + `LowerBoundBasis`; `TotalActions` stays the operator's (`protocol_migration.go:40`) |
| Cap frozen, never raised | `protocolMigrationPolicy` still writes `Maximum: q.Maximum, Carried: 0`; `InitialMaximum()` (`cycle_binding.go:145`); stricter-track check (`cycle_budget.go:47-49`) |
| Idea/kind/shared-repo scope bound | `protocolMigrationScope`; identity checks `:185`/`:214`; scope keyed on the git-common-dir prefix (`binding.go:132`) |
| Ordinary bootstrap still refused | the three Ensure/Configure sites keep calling the `nil`-declaration wrapper; no record exists for them to read |
| Fail on stale registry/history/policy | declared-path registration check; `errHistoryChanged` re-inspect before activation (`:321`); `checkProtocolMigrationPolicy` |
| Files/counters preserved on failure | every refusal returns before `persist`; `migration-active` is written last (`:324`), so a partial apply leaves the policy inactive |
| Retry / concurrency | guard (`:204`) + post-guard re-read (`:209`) + *"another protocol migration decision won publication"* (`:215`) + `ledger.update` rejecting non-empty entries (`:266-270`) |
| Never prune / recreate / zero-fill / reclassify / clone / call providers uncharged | nothing here deletes a registration, creates a directory, writes a count into an unknown row, touches `CycleKindForPhase`, weakens `binding.go:122-130`, or adds a launch path |

## 5. Edit allocation (production code, estimate)

`binding.go` ≈45 · `launch_migration_history.go` ≈20 · `protocol_migration_history.go` ≈15 · `protocol_migration.go` ≈20 (field, normalization, three `InspectProtocolMigration` call sites, four validity clauses) · `migration_recovery.go` ≈8 (carry `declared` in `recoveryBase`, pass at `:472`, else a declared import can never be recovered) · `app/budget_migrate*.go` ≈25 (repeatable `flag.Value`, allow-lists, launch-kind rejection). **≈133 lines / 6 files.**

## 6. Tests

1. **Real missing-worktree witness** — repo + linked worktree, directory removed: undeclared inspect still errors `historical worktree is unavailable`; declared inspect yields one row `history:"unknown"`, `history_coverage:"declared-incomplete"`, path absent from `roots`.
2. **Stale digest** — apply with a digest predating a new file in a surviving root → `errHistoryChanged`; `migration-active` absent, `ledger.json` byte-identical.
3. **Observation change** — a retained pre-start refusal whose terminal gains tokens → inspect errors (`launch_migration_history.go:367-377`); via the journal, *"previously imported terminal observation changed"*.
4. **Replay vs conflict** — identical request twice (incl. `nil` vs `[]string{}`) → idempotent, one record, `N` entries once; a request differing only in one declared path → *"conflicting protocol migration replay"*.
5. **Crash recovery** — inject the `persist` seam (`:170`) to fail after `migration.json`, before `migration-active`: `readCyclePolicy` → *"not durably active"*, no charge grantable, exact replay completes.
6. **Unavailable → available** — recreate the declared directory between inspect and apply → refused; policy, ledger and counters untouched.
7. **No launch before a fully applied migration** — `migration.json` present, `migration-active` absent: `groupProtocolCycle` defers and `ChargeCycle` returns the refusal (`cycle_session.go:74-75`); ledger empty.
8. **Request-scoping** — `EnsureCycleBinding` on the same broken repo still refuses, proving no declaration leaked into ordinary bootstrap.
9. **Concurrency** — two identical concurrent applies → exactly `N` entries; two conflicting → one wins, the other gets *"another protocol migration decision won publication"*. No refund, overwrite or partial grant.

## 7. The human decision still required (no values invented)

1. **Attest the two registered paths as unavailable with unknown history.** Not yet given — roster change and pilot amendment were authorized; this was not.
2. **`--total-actions N`.** The tool enforces only `N ≥ lower_bound`, and that floor covers the *surviving* worktrees alone. Nothing here establishes that the declared roots' history falls within the cap, and this design must not be read as implying it.
3. **Consequence to see before choosing:** with `--max-cycles 3`, remaining grouped rounds `= 3 − N` (PRIMARY: `cycleLimits` `cycle_intent.go:29`; `Store.Reserve` `ledger.go:168-176`; `Carried: 0`). If `N ≥ 3`, round-02 is refused.
4. **`--max-cycles`** — I propose `3`, the existing deliberation cap; explicitly not an increase. **`--started-at` / `--decision-id` / `--reason`** — operator-authored, mechanically constrained.

**Collectable before asking:** the two verbatim porcelain paths and their stat observation; step 1's `history_sha256`, `lower_bound`, `evidence[]` rule names, `ungrouped_invocations[]`, `earliest`; the surviving `round-NN/` names under the idea (all `LegacyCycleFloor` reads, `cycle_history.go:90-94`); the refused launches' classifications; and whether any `policy.json` / `migration.json` / ledger already exists in the cycle scope dir (first-import vs recovery). The two paths and the refusal records are **unverified testimony** from codex-1's note — not reproduced, no verdict assigned. Grouped round, unchanged: one `runner.RunRound` session, skip codex-1's published round-02 via the existing owned-artifact path, keep all four members, no cap increase proposed.

## 8. Key unresolved issue

**Whether `--total-actions` can honestly be chosen at all while the two roots' history is unknown.** This design makes the unknown explicit and auditable; it cannot make it small. The operator picks a number that is a floor over surviving evidence plus a judgment about an unreadable region, and `3 − N` decides whether the amendment round can run. If originals or backups arrive, re-run step 1 with fewer declarations — an enumerable root beats any attestation. Secondary, for a reviewer to rule on: the `omitempty` compatibility choice makes "field deleted from an empty list" and "no declaration" indistinguishable — safe today, but a real weakening of the digest's tamper-evidence.

**Own path:** `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/claude-1-history-admission-design-20260916.md`
