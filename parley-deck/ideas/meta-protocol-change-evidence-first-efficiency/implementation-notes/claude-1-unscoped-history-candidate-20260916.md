---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: minimum-additive-source-candidate
status: in-progress
supersedes-scope-of: claude-1-unscoped-history-remedy-20260916.md (design only; this note is the built subset)
builds-on: claude-1-history-admission-candidate-20260916.md (declared-unavailable worktrees, already in tree)
method: native Read/Write over this worktree only. No shell, Git, build, test, CLI, migration, apply,
  attestation, provider or subagent. Nothing here was compiled or executed by me.
---

# Unscoped historical runs — minimum additive source candidate

**Source only. No apply, no attestation, no count chosen, no identity assigned, no history mutated.**

## 0. Coordinator direction I am building to (recorded verbatim in substance)

Accepted as sound: the absent-vs-conflicting distinction, the whole-directory digest, retaining the
readable bytes, and asserting no zero count. Corrected: my four "current flaws" (remedy note §7) are
**extension constraints, not demonstrated failures** of the existing declared-worktree-only
semantics. I therefore do **not** replace `DeclaredUnavailable`, do **not** widen scalar
`LowerBoundBasis` to a list, and do **not** unify the declaration types. Existing saved records and
nil-declaration digests stay byte-identical.

I record one correction to my own earlier text, per §15.1: remedy §7.2 asserted the existing
`checkDeclaredUnavailable` coverage clause "would wrongly refuse every valid unscoped-only import".
That is a statement about a **hypothetical extension**, not about shipped behaviour — with no
unscoped declaration in existence the clause is exactly right. The constraint is real (the clause
must learn a second declaration kind); the "flaw" framing was wrong and is withdrawn.

## 1. What is built here

Additive only. Every existing signature, persisted field and nil-declaration digest is preserved.

| # | File | Change |
|---|---|---|
| 1 | `internal/budget/run_identity_inventory.go` *(new)* | `RunDirectoryManifest` / `RunDirectoryManifestDigest` (bounded recursive file-set digest, symlink/non-regular refused); `InspectRunIdentities` + `RunIdentityReport`/`RunIdentityRow`, read-only, tolerant of absent identity, grants nothing |
| 2 | `internal/budget/launch_migration_history.go` | `UnscopedRun` row type; `normalizeDeclaredUnscopedRuns` / `parseDeclaredUnscopedRuns` (`parley-deck/runs/<name>=<64hex>`); `LaunchMigrationInventory += UnscopedRuns []UnscopedRun` (`omitempty`); `inspectLaunchMigrationDeclarations` with both declarations; `runs()` learns the declared branch |
| 3 | `internal/budget/protocol_migration_history.go` | new `InspectProtocolMigrationDeclarations(ctx, root, idea, kind, ideaPath, declared, unscoped)`; `InspectProtocolMigrationDeclared` **keeps its signature** and delegates with `nil`; `LowerBoundBasis` gains two explicit combined tokens |
| 4 | `internal/budget/protocol_migration.go` | `ProtocolMigrationRequest += DeclaredUnscopedRuns []string` (`omitempty`), **alongside** `DeclaredUnavailable`; both canonicalized at the same boundary; `checkDeclarations` splits the coverage clause per kind; all three inspect sites carry both |
| 5 | `internal/budget/migration_recovery.go` | `recoveryBase.unscoped` read back from the frozen request; recovery preview re-inspects with both |
| 6 | `internal/app/budget_migrate*.go` | repeatable `--declare-unscoped-run`; protocol `inspect`/`apply` only; refused for `--kind launch` and every other verb |

**No `parley budget runs` command.** The coordinator obtains the authoritative declaration digest
from `InspectRunIdentities` / `RunDirectoryManifestDigest` as a source helper; both are testable
without any product endpoint. Complete inspect/apply plumbing took priority, as directed.

## 2. The distinction the remedy turns on (unchanged, and the reason retention is not optional)

| | declared-unavailable **root** | declared-unscoped **run** |
|---|---|---|
| Bytes | unreadable | **readable, and read** |
| Unknown | everything | **identity only** |
| Handling | excluded from `Roots`/`Sources` | **retained in `Sources`** — bytes are evidence of themselves |
| Binds to | registered path + stat class | relative run path + recursive manifest digest |

Retention is load-bearing: because the declared run's `events.jsonl` stays in `Sources`, the existing
re-read-and-compare in `InspectProtocolMigrationDeclared` covers it for free, and any later byte
change refuses with `errHistoryChanged`. **There is no structural exclusion of readable bytes** —
the exclusion is from *scope-count evidence* only, and it is recorded explicitly.

## 3. Admission rules — every failure is a hard error, never a silent skip

1. Declaration value is canonical relative `parley-deck/runs/<name>` `=` 64 lowercase hex, exactly
   three slash segments, no `..`, backslash, NUL or newline, `Clean`-stable, not absolute.
2. Canonical form is sorted and deduplicated; two digests for one path refuse.
3. In **every** visible root where the directory exists, the recursive manifest digest must equal the
   declared digest. An added, deleted or changed file in any copy refuses (`differs at <root>`).
4. Identity must be **actually** unrecoverable. A recovered name refuses — the exact mirror of
   `declared-unavailable worktree is available again`.
5. **Only absent identity is declarable.** Conflicting or malformed identity refuses, because a
   conflicting run *holds* identity evidence and declaring it away would discard evidence.
6. A declared path matched in no visible root is a stale declaration and refuses.
7. **Undeclared keeps the verbatim existing refusals**, both of them, at their existing sites.
8. Symlinked or non-regular entries anywhere under a declared run directory refuse.

## 4. What the declared run does *not* contribute

- No `Evidence` row under any rule, and therefore no floor — **and no `published-action-floor` of
  zero either.** Nothing asserts zero. (It reaches the protocol counter, gets `identity == ""`, and
  is skipped by the pre-existing `identity != idea` filter.)
- No `agent.started` harvest, no `s.starts` entry, no `LegacyStartFloor` bump.
- No contribution to the inventory's `Earliest`. The row records the run's own earliest event time
  because it is useful to the operator, and it is **never** read as this idea's epoch authority.
- `Copies` counts copies of one file set. It is never an action count, and every row carries
  `History: "unknown"`.

## 5. Coverage tokens (scalar preserved; old cases byte-identical)

| declared | `LowerBoundBasis` |
|---|---|
| neither | `""` *(unchanged)* |
| unavailable roots only | `surviving-visible-worktrees-only` *(unchanged, byte-identical)* |
| unscoped runs only | `scoped-run-history-only` *(new)* |
| both | `surviving-visible-worktrees-and-scoped-run-history-only` *(new)* |

`HistoryCoverage` is `declared-incomplete` whenever either kind is declared. The `checkDeclarations`
gate is bidirectional in both directions for both kinds, so a record cannot carry rows without the
declaration or a declaration without the rows.

## 6. Containment

Ordinary `EnsureCycleBinding` / `ConfigureLaunchBudget` and every nil-declaration caller still refuse
on the same repository, with the same text. The declaration is a request-scoped argument to one
decision: it is retained inside the immutable record for replay and recovery, and **no bootstrap
path consults a persisted exclusion.**

## 7. Disclosure — forward compatibility

`unscoped_runs` and `declared_unscoped_runs` are `omitempty`, so an inventory or request with no
unscoped declaration marshals exactly as before and keeps its digest: a record already applied by
today's binary still reads. **The reverse does not hold — an older binary rejects these fields as
unknown under `DisallowUnknownFields` once they are non-empty.** Same disclosure as the
declared-unavailable fields, now for a second pair.

Also disclosed: the target blob is a **tracked** file, so a `git checkout` between inspect and apply
can change the manifest digest and refuse a frozen declaration. Correct, fail-closed, and
operationally surprising.

## 8. Residuals (unresolved, and not resolved by this code)

1. **`<N>` is still not honestly chooseable, and the machine constrains it *less* where evidence is
   thinnest.** `TotalActions >= LowerBound` is the only machine constraint, and excluding evidence
   lowers the floor. Recording the basis makes that visible; it does not make it safe. **I choose no
   `N` and supply no value for one.**
2. The blob looks like a smoke fixture (tracked since the initial CLI commit, task
   `"smoke implementation run"`, no cursor, no starts). Its true count is plausibly `0`. This code
   **refuses to act on that plausibility**: absence of `agent.started` is not proof of zero.
3. Whole-directory manifests bind more than the scanner reads (`events.jsonl` + `driver.json`), so an
   unrelated log file invalidates a declaration. Deliberate and fail-closed.
4. The brief's 150/25 figures are testimony I could not verify from this restricted root; per §15.2
   they are `RECALL` to me and nothing here depends on them. The design scales by *distinct
   (relative path, manifest digest)*, which is the same whether the copies number 25 or 250.
5. **Nothing here was compiled or tested by me** (no Bash in this launch). Every claim about
   behaviour is a claim about source I wrote and read, not observed execution. The coordinator
   compiles and runs the tests after the terminal.

## 9. What this explicitly does not do

No identity assigned to any run. No count inferred, including not zero. No history edited, pruned,
recreated or reclassified. No old-schema exemption, silent or otherwise. No `<N>` selected, no cap
changed, no migration applied, no attestation, no signoff, no current-source acceptance or amendment.

## 10. Status

See §11 for the honest completion record, written after the code.

**Own path:** `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/claude-1-unscoped-history-candidate-20260916.md`
