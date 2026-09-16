---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: implementation-candidate
status: candidate
implements: claude-1-history-admission-design-20260916.md as corrected by
  claude-1-history-admission-correction-20260916.md and kimi-1-history-admission-design-review-20260916.md
method: native Read/Write/Edit only; no shell, Git, build, test, CLI, browser, subagent or other provider
---

# MRW-1 candidate: request-scoped declared-unavailable worktrees

**Code candidate only — not built, not run, not reviewed, not an attestation, not an apply.** No
migration was applied, no count chosen, no cap changed, no source accepted. The coordinator runs
the checks; I could not, so every claim below about behavior is a claim about the code I wrote,
not about an observed execution.

## 1. Plan (written before the code)

Thread one **request-scoped** typed declaration through the existing migration — no new writer, no
new global authority store, no persisted declaration consulted by ordinary bootstrap.

| Step | File | Change |
| --- | --- | --- |
| 1 | `internal/budget/binding.go` | `UnavailableRoot{Path,Observation,History}`; `normalizeDeclaredUnavailable` (validate/sort/dedupe); `launchScopeDeclared(..., declared)` returning `missing []UnavailableRoot`; `launchScope` = nil wrapper (all existing call sites unchanged) |
| 2 | `internal/budget/launch_migration_history.go` | `LaunchMigrationInventory += UnavailableRoots []UnavailableRoot`, `HistoryCoverage string` (both `omitempty`); `InspectLaunchMigrationDeclared` |
| 3 | `internal/budget/protocol_migration_history.go` | `ProtocolMigrationInventory += LowerBoundBasis string` (`omitempty`); `InspectProtocolMigrationDeclared` |
| 4 | `internal/budget/protocol_migration.go` | `ProtocolMigrationRequest += DeclaredUnavailable []string` (`omitempty`); normalize beside `StartedAt.UTC()`; validity clauses in `protocolMigrationInitial`; three inspect call sites thread the declaration |
| 5 | `internal/budget/migration_recovery.go` | `recoveryBase.declared`, carried from the retained request, threaded into the recovery preview's re-inspection |
| 6 | `internal/app/budget_migrate*.go` | repeatable `--declare-unavailable-worktree`; protocol `inspect`/`apply` only; rejected for `--kind launch` and for every other verb |

Guards preserved as-is: ordinary `Ensure*`/`Configure*` keep passing a nil declaration and stay
fail-closed; counters, caps, activation order and unknown-history rules unchanged.

## 2. What the declaration binds (and what it refuses)

A declared row is admitted only when the *actual* stat class is unavailable and the registration
is still present in `git worktree list --porcelain`:

- `stat-enoent` / `stat-not-directory` — the only declarable classes, recorded verbatim;
- any other stat error — hard error, never declarable;
- available again (a real directory at apply time) — hard error `declared-unavailable worktree is
  available again`, so a migration never runs over now-readable evidence;
- declared but absent from the porcelain output — hard error (stale registry);
- undeclared and unavailable — the existing `historical worktree is unavailable` refusal, verbatim.

`History` is pinned to the string `"unknown"` and validated on every read: never `0`, never a
count. Declared paths are excluded from `Roots` and `Sources`, so no evidence, floor or count is
claimed from them; `HistoryCoverage = "declared-incomplete"` and `LowerBoundBasis =
"surviving-visible-worktrees-only"` make the incompleteness explicit in the digest. The canonical
form is sorted and deduped, so `nil` and `[]string{}` cannot split an exact replay.

## 3. Compatibility, disclosed in both directions

New binaries read old records: every new field is `omitempty`, `migrationShape` tolerates absent
`omitempty` tags, and an empty declaration marshals byte-identically to today, so old digests
still validate. **Old binaries reject new records**: `DisallowUnknownFields` refuses
`declared_unavailable` / `unavailable_roots` / `history_coverage` / `lower_bound_basis`. The
direction is not symmetric; this is an accepted, visible limitation, not a fixed one. My §8
`omitempty` tamper-evidence caveat (deleted-empty-list is indistinguishable from no declaration)
also stands, unresolved.

## 4. Departures from my proposal, and why

1. **`Roots`/`Sources` exclusion made structural rather than incidental.** `launchScopeDeclared`
   `continue`s before the candidate-path stat, so a declared path cannot enter `roots` at all.
2. **Declarations are refused unless history inspection is actually running against a Git
   repository** (`inspectHistory == true` and `git rev-parse` succeeded). My design left this
   implicit; fail-closed is the correct reading of "never repository state".
3. **`LowerBoundBasis` is set only when a declaration was admitted.** An undeclared inventory keeps
   its exact current bytes and digest — required for the compatibility claim above.
4. **Set equality is bidirectional in `protocolMigrationInitial`**: every declared path must appear
   as an unavailable row, and every unavailable row must have been declared. Kimi's review asked
   for (c)/(d); one-way containment would have let an inventory carry an undeclared unknown row.
5. **No `--declare-unavailable-worktree` on `recover`.** The recovery verb re-inspects with the
   declaration retained *inside the frozen request*, so accepting it again on the command line
   would offer an operator a second, divergent declaration for the same immutable decision.
6. **Added a structural exclusion check inside the validity clauses** (`checkDeclaredUnavailable`):
   no declared path may appear as, or as a parent of, a retained history root or source root. The
   inspect path already cannot produce one; the check makes a hand-edited or future-constructed
   record fail closed rather than quietly counting evidence from a path called unreadable.
7. **The three inspect call sites became four.** My §5 estimate named
   `InspectProtocolMigration` at `protocol_migration.go:196`, `:235` and `:317`, plus the recovery
   site; the pre-activation re-inspection at `:317` is the one that also re-checks availability,
   registration and stat class immediately before the activation witness is written.

Not departures, restated because they were explicit non-goals: no total-actions `N`, no
`--started-at`, no `--decision-id`, no operator reason is chosen here; no cap raised; no
attestation; no prune, recreate, zero-fill, reclassify or clone.

## 5. Tests written (`*_declared_test.go`)

Package `budget` (real `git init` + real linked worktree, directory actually removed — no
fabricated porcelain bytes):

1. undeclared inspect still refuses; declared inspect yields one `stat-enoent` row with
   `history:"unknown"`, `history_coverage:"declared-incomplete"`, `lower_bound_basis` set, and the
   path absent from `roots` and from every source root;
2. a declaration over a registration that is actually present refuses as available-again; an
   unregistered, relative or empty path refuses; a registration replaced by a regular file binds
   the distinct `stat-not-directory` class rather than the missing one;
3. exact replay (including `nil` vs `[]string{}` vs unsorted/duplicated input) is idempotent and
   preserves spend and epoch; a request differing only in one declared path conflicts;
4. an injected publication fault before `migration-active` leaves the import inactive; the
   recovery preview still re-inspects (proving the declaration threads through recovery) and the
   exact replay completes the activation;
5. ordinary bootstrap on the same broken repository still refuses — `EnsureCycleBinding` for
   another idea, and `ConfigureLaunchBudget` for the same one — proving no declaration leaked
   into repository state;
6. an undeclared inventory still carries none of the new fields, so an already-applied import
   keeps its bytes and digest in this binary.

Package `app`: repeatable flag accepted for protocol `inspect`/`apply`, rejected for
`--kind launch` (inspect and apply), for `migrate recover inspect|apply`, for an unknown `migrate`
verb and for `budget worktree inspect`; inspect output carries the declared rows with unknown
history; an attended apply replays exactly and an apply with the declaration removed refuses; no
accounting state is created by a refused or read-only command.

Exact files written: `internal/budget/protocol_migration_declared_test.go`,
`internal/app/budget_migrate_declared_test.go`. Both skip when `git` is unavailable.

## 6. Unfinished / honest gaps

- **Nothing was executed.** No `go build`, `go vet` or `go test` ran in this session (no Bash by
  instruction). Compile errors and test failures are both possible; the coordinator's run is the
  first real evidence. Treat §5 as "tests written", not "tests passing".
- The preview target correction is carried in text only: the v2 idea and the amendment worktree
  root are the intended `--idea`/`--dir`, and the launcher remains **one direct
  `runner.RunRound` call, not the driver**. No code here selects a target, and nothing here
  decides my §7 residual (driver-step activation would be a second, separately scoped decision).
- Step/launch scopes are untouched: a declared import for `DriverStep` is mechanically possible
  through the same flag, but no such decision is proposed, prepared or implied.
- `--total-actions` honesty over an unreadable region (my §8) is unchanged and still open.

**Own path:** `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/claude-1-history-admission-candidate-20260916.md`
