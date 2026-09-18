---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: bounded-design
status: proposal-only
responding-to: the real v2 preview refusal "historical run identity is missing or conflicting"
builds-on: claude-1-history-admission-candidate-20260916.md (declared-unavailable worktrees; flaws listed in §7)
method: native Read only over this worktree's source tree; no shell, Git, build, test, CLI, migration,
  attestation or subagent. Nothing here was compiled or executed.
---

# Unscoped historical runs: minimum evidence-preserving remedy

**Design only. No implementation, no apply, no count, no attestation, no identity assigned.**

## 0. Recommendation, in four lines

1. **Fix the four §7 flaws in my declared-worktree candidate before anything is applied** — one of them (§7.2) would refuse every valid unscoped-run import outright, and two more become persisted-shape breaks the moment an apply lands.
2. **Ship §5.A alone first**: a read-only `budget runs inspect` that enumerates and classifies run identities and grants nothing. Today the scanner aborts on the first offender, so no product surface can tell the operator what to declare.
3. **Then §5.B**: a request-scoped `--declare-unscoped-run <rel-run-dir>=<manifest-sha256>`, bound to content, threaded through preview/apply/replay/recovery. Per the brief's figures this operation needs **one** such flag.
4. **Do not** adopt either rejected shortcut in §4, and do not let the floor be printed without its exclusions beside it.

## 1. Provenance separation (§15.1/§15.2)

`PRIMARY`, read by me in this worktree just now:

- `parley-deck/runs/20260510T194003Z/events.jsonl` is exactly one line — `{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"mode":"auto","task":"smoke implementation run"}}` — and `data` carries **no** `idea` or `idea_slug` key.
- That directory contains **only** `events.jsonl` (glob over `parley-deck/runs/**`): no `driver.json`, no `run.json`, no `agents/`. The five other run dirs here do carry `run.json` and/or `agents/*.log`.
- The refusal site: `internal/budget/launch_migration_history.go:487-497` (`bind`) invoked at `:533` with `required = (name == "idea" && e.Type == "run.created")`. With `e.Data["idea"] == nil` and `required == true`, `bind` falls through to `json.Unmarshal(nil, &name)`, which errors, and returns `errors.New("historical run identity is missing or conflicting")` — the exact observed text.
- `runs()` reads `events.jsonl` via `s.file(...)` at `:476` **before** that check, so the bytes are already in `Sources` when the refusal fires; the whole inventory is then discarded by `InspectProtocolMigrationDeclared` (`protocol_migration.go:96-99`).
- The scope filter `if identity != idea { continue }` sits at `:548`, **after** the refusal. A run naming *another* idea is skipped silently; a run naming *none* aborts the whole inspection. That asymmetry is deliberate and correct (see §4) and is the reason no ordering change fixes this.

`SECONDARY`/testimony, **not mine and not verified here** — the launch brief's figures: 150 visible run instances, 27 registrations, 25 issue rows, all one 117-byte blob at source hash `ee4f52b7…3c6906`, same bytes since CLI commit `3ec10ac`. The diagnostic at `/Volumes/.../visible-run-identities.json` is outside this restricted root; my Read was refused, so I do not own those numbers. The design below does not depend on any of them being exact — it scales by *distinct (relative path, content digest)*, and that count is `1` whether the copies number 25 or 250.

## 2. What actually broke, in one sentence

The two declared-unavailable worktrees were never the only blocker: a **tracked** file present in every checkout carries a `run.created` event with no idea key, and the scanner treats an unrecoverable identity as fatal for the entire repository. **Restoring the two missing worktrees would not unblock this migration** — the blob is in the restored checkouts too. The worktree problem and the run problem are independent.

## 3. The distinction the remedy turns on

Two different unknowns, and my existing candidate only models the first:

| | declared-unavailable **root** | unscoped **run** |
|---|---|---|
| Bytes | unreadable | **readable, and read** |
| Unknown | everything | **identity only** |
| Correct handling | exclude from `Roots`/`Sources` — it can supply no evidence | **retain in `Sources`** — the bytes are evidence of themselves; exclude only from identity-scoped counting |
| Binds to | registered path + stat class | relative run path + content digest |

Retention is the whole point. Keeping the declared run's bytes in `Sources` means the existing re-read-and-compare at `protocol_migration.go:108-121` covers them for free: any later change refuses with `errHistoryChanged`. Excluding them (the worktree pattern) would *lose* evidence we hold.

Checked against the current counting path: a retained unscoped run enters `runs[]` at `protocol_migration.go:221`, gets `identity == ""` from `protocolMigrationEvents`, and hits `continue` at `:264`. So it adds **no** `Evidence` row and **no** floor — importantly, it does not add a `published-action-floor` of `0` either. Nothing asserts zero. That is already the behaviour; the design must preserve it, and §6/T3 tests it.

## 4. Answering the two questions asked

**"Can the visible floor be presented before operator authorization, with explicit incomplete coverage?"**
Yes — and the confusion to clear first is that **a declaration is not an authorization.** `budget migrate inspect` writes nothing, spawns nothing, takes no `--yes`, and grants nothing; it is already a pre-authorization surface. Adding `--declare-unscoped-run` to it adds a *binding argument*, not a permission. Authorization remains exactly one act: `apply --yes --total-actions <N> --expected-history-sha256 …`.

What must **not** happen is a floor computed by a path that ignores unscoped runs *without* recording which ones — that is the silent old-schema exemption in report form. So the rule is: **the floor is presented only alongside the enumerated, digest-pinned exclusions that produced it.** Both live in one artifact.

**"Distinguish source implementation from a real unknown-history/count grant."**
Writing this code grants nothing. The grant is three separate later acts, none of them here: (i) the operator naming each exclusion verbatim, (ii) the operator choosing `<N>` over a region known to be incomplete, (iii) `apply --yes`. §5 is source; the grant is the operator's. My correction note's §3 still holds — cycle, step and launch authority remain three decisions, not one.

**Is there a safer simpler alternative?** I looked for one and reject the two candidates:

- *Move the refusal after the scope filter* (`:548` before `:533`). Rejected: a run naming another idea is *known* not to be ours; a run naming none *might* be ours. Skipping it is precisely the forbidden inference. The asymmetry must stay.
- *Avoid the migration by using a fresh scope.* Rejected: `ConfigureLaunchBudget`/`EnsureCycleBinding` refuse **because** they scan and find this, so a fresh scope does not avoid the scan.

There is no simpler shape that preserves zero-inference. But there **is** a safer *sequencing*: ship the non-granting discovery surface (§5.A) first, alone. It has no declaration, no write path and no grant, and without it the operator cannot even enumerate what to declare — today `runs()` aborts on the first offender, so the count of offenders is unknowable through any product surface. That is why the ad-hoc diagnostic had to be hand-built.

## 5. Exact changes

### 5.A — Ship first: non-granting discovery (`internal/budget/run_identity_inventory.go`, new)

```go
const (
    RunIdentityBound       = "bound"                     // an idea name was recovered
    RunIdentityAbsent      = "run-identity-absent"       // no idea key anywhere — the only declarable class
    RunIdentityConflicting = "run-identity-conflicting"  // two different names — never declarable
)

type RunIdentityRow struct {
    Path       string   `json:"path"`         // relative: parley-deck/runs/<name>
    SHA256     string   `json:"sha256"`       // directory manifest digest (below)
    Observation string  `json:"observation"`
    Idea       string   `json:"idea,omitempty"` // only when Observation == RunIdentityBound
    History    string   `json:"history"`      // UnknownHistory unless bound
    Copies     int      `json:"copies"`       // copies of *this file set*, never an action count
    Roots      []string `json:"roots"`        // sorted, verbatim
    Declaration string  `json:"declaration,omitempty"` // copy-paste flag value, absent class only
}

type RunIdentityReport struct {
    Schema, Registrations int
    Root, Idea, Coverage  string          // Coverage is always "incomplete-unknown-scope"
    Rows                  []RunIdentityRow
    Uncertainty           []string
    ReportSHA256          string
}

func InspectRunIdentities(ctx context.Context, root, idea string, d MigrationDeclaration) (RunIdentityReport, error)
func (r RunIdentityReport) OperatorRecord() ([]byte, error)
func runDirManifestDigest(dir string) (string, error)
```

Contract, copied deliberately from `WorktreeInventory.OperatorRecord` (`worktree_inventory.go:99-110`): **it grants nothing — no apply, attest or confirm path in this package consumes it.** It never aborts on an absent identity (that is its job); it still fails closed on a broken repo, an unreadable root, malformed JSON or a non-regular entry. It computes **no action floor and no `N`** — only identity classification and file-set digests. It takes `MigrationDeclaration` because `launchScopeDeclared(..., true, nil)` still refuses the two missing roots; the worktree declarations are already discoverable from `budget worktree inspect`, which enumerates registrations tolerantly today.

`runDirManifestDigest` is recursive over the run directory, sorted by relative slash path, `sha256` over `[{name,sha256},…]`, refusing symlinks and non-regular entries. For `20260510T194003Z` that is a one-entry manifest (`PRIMARY`, §1). Recursion is intentional: an `agents/` subdir or a later `driver.json` changes the digest and invalidates the declaration, which is the "refusing additional/changed files" requirement met by one rule instead of three.

New CLI `parley budget runs inspect --dir <root> --idea <slug> [--declare-unavailable-worktree …]` (`internal/app/budget_runs.go`). Read-only. Prints the report and, per absent-class row, the exact `--declare-unscoped-run` argument to copy. **Never prints an `N` and never prints a `--total-actions` value.**

This yields the ladder, each rung supplying the next rung's arguments, and only the last granting anything:
`worktree inspect` → `runs inspect` → `migrate inspect` → *operator decides `<N>`* → `migrate apply --yes`.

### 5.B — Then: the request-scoped unscoped-run declaration

| # | File | Change |
|---|---|---|
| 1 | `binding.go` | `type MigrationDeclaration struct { UnavailableWorktrees, UnscopedRuns []string }`; `type UnscopedRun struct { Path, SHA256 string; Copies int; Roots []string; History string; Earliest *time.Time }`; `normalizeDeclaredUnscopedRuns([]string) ([]unscopedDeclaration, error)` parsing `<relative-run-dir>=<64-hex>`, sorted/deduped, refusing two digests for one path |
| 2 | `binding.go` | `launchScopeDeclared(..., declared []string)` → `(..., d MigrationDeclaration)`; the run declaration is validated but not consumed here |
| 3 | `launch_migration_history.go` | `LaunchMigrationInventory += UnscopedRuns []UnscopedRun` (`omitempty`); extract `scanRunIdentity(events, cursor []byte) (name string, parsed []event, err error)` from `runs()` — **one** parser, returning `""` for absent and erroring on malformed/conflicting; the *caller* converts absent→refusal when undeclared |
| 4 | `launch_migration_history.go` | in `runs()`, after the `driver.json` read and **before** the event harvest: if the dir's manifest digest matches a declaration, require `name == ""` (else refuse `declared-unscoped run has a recoverable idea identity: <name>`), record/merge the `UnscopedRun` row, `continue`. No `agent.started` harvested, no `s.starts` entry, no `LegacyStartFloor` bump |
| 5 | `protocol_migration.go` | `ProtocolMigrationRequest.DeclaredUnavailable []string` → `Declaration MigrationDeclaration`; `LowerBoundBasis string` → `[]string` (sorted, `omitempty`) with tokens `SurvivingVisibleFloor` and new `ScopedRunsOnly = "scoped-run-history-only"`; `checkDeclaredUnavailable` → `checkDeclaration(i, d)` with the two coverage clauses **split** (§7.2) and bidirectional set equality for both kinds |
| 6 | `migration_recovery.go` | `recoveryBase.declared` → `recoveryBase.declaration MigrationDeclaration`, carried from the frozen request into the recovery preview's re-inspection |
| 7 | `budget_migrate*.go` | repeatable `--declare-unscoped-run '<rel-run-dir>=<manifest-sha256>'`; protocol `inspect`/`apply` only; rejected for `--kind launch`, `migrate recover inspect|apply`, unknown verbs, `budget worktree inspect`, `budget runs inspect` |

**Admission rules — every failure is a hard error, never a silent skip:**

1. Path is canonical relative, prefix `parley-deck/runs/`, exactly one further segment, no `..`, `\`, NUL or newline.
2. In **every** visible root where the directory exists, its manifest digest equals the declared digest; a differing root refuses (`declared-unscoped run differs at <root>`).
3. Identity must be **actually** unrecoverable — a recovered name refuses. This is the exact mirror of `declared-unavailable worktree is available again` (`binding.go:196`).
4. **Conflicting identity is never declarable**, only absent. The current error text conflates the two; a conflicting run *holds* identity evidence, so declaring it away would discard evidence. For this data all rows are absent-class, so the restriction costs nothing here and closes a real hole.
5. Declared but present in no visible root → stale declaration, refuse (mirrors `is not a retained registration`, `binding.go:221`).
6. **Undeclared unscoped run keeps the verbatim existing refusal.** No exemption, silent or otherwise.

One flag covers all copies of one blob, however many worktrees hold it. Per the brief's (unverified) figures, this operation needs **one** `--declare-unscoped-run`.

**Deliberate choice needing a reviewer's eye:** the declared run's event timestamp is recorded on the row as `Earliest` but is **not** folded into the inventory's `Earliest`. Folding it in would tighten the `StartedAt` epoch check at `protocol_migration.go:95` (conservative, fail-safe) but would derive a claim about *this idea's* history from a run we just declared unidentifiable. I chose record-don't-enforce; either direction is defensible and I do not want it settled silently.

### 5.C — The resulting preview, concretely

Every angle-bracket token is a placeholder. **I supply no value for any of them**, least of all `<N>`.
Steps 1–3 write nothing, spawn nothing and grant nothing; step 5 is the only authorization.

```
# 1. registrations, tolerant of the two missing roots   (exists today)
parley budget worktree inspect --dir <AMENDMENT-WORKTREE-ROOT>

# 2. NEW (§5.A): classify run identities; prints the exact flag values for step 3
parley budget runs inspect --dir <AMENDMENT-WORKTREE-ROOT> \
  --idea meta-protocol-change-evidence-first-efficiency-v2 \
  --declare-unavailable-worktree '<PORCELAIN-PATH-A>' --declare-unavailable-worktree '<PORCELAIN-PATH-B>'

# 3. NEW (§5.B): the preview that currently refuses. Writes nothing. Emits the floor
#    together with its enumerated, digest-pinned exclusions, and a history digest.
parley budget migrate inspect --kind cross-review \
  --dir <AMENDMENT-WORKTREE-ROOT> --idea meta-protocol-change-evidence-first-efficiency-v2 \
  --declare-unavailable-worktree '<PORCELAIN-PATH-A>' --declare-unavailable-worktree '<PORCELAIN-PATH-B>' \
  --declare-unscoped-run 'parley-deck/runs/20260510T194003Z=<MANIFEST-SHA256-FROM-STEP-2>'

# 4. human reads the preview and decides <N> over a knowingly incomplete floor. No command.
#    Not decided here, and §8.1 says the machine checks this less, not more.

# 5. attended guarded apply — identical declarations, digest from step 3
parley budget migrate apply  --kind cross-review  … same three --declare-* flags … \
  --expected-history-sha256 <DIGEST-FROM-STEP-3> --total-actions <N> --max-cycles 3 \
  --started-at <RFC3339-NOT-AFTER-EARLIEST> --decision-id <UNIQUE> --reason '<BASIS>' \
  --writers-stopped --yes
```

The run-directory name in step 3 is the one fact I verified myself (§1); its digest is not, because a
manifest digest is a property of all visible copies and I read only this worktree's.

Expected shape of the step-3 output, for review before it is built: `lower_bound` with
`lower_bound_basis: ["scoped-run-history-only","surviving-visible-worktrees-only"]`, one
`unavailable_roots` row per declared path with `history:"unknown"`, one `unscoped_runs` row with
`history:"unknown"` and `copies:<k>`, and `history_coverage:"declared-incomplete"`. No field of that
output is a total.

## 6. Semantic tests

Package `budget`, real `git init` + real linked worktrees, a run dir constructed in the test to the shape verified in §1:

- **T1** undeclared → the refusal text is unchanged byte-for-byte.
- **T2** declared → inspection completes; one `UnscopedRun` row, `History == UnknownHistory`, `Copies` = roots carrying it, `Roots` sorted verbatim, `LowerBoundBasis` contains `ScopedRunsOnly`.
- **T3** *(the zero-inference test)* the declared run contributes **no** `Evidence` row under any rule and `LowerBound` is byte-identical to the same repo with that dir absent — proving neither ≥1 counted nor 0 asserted.
- **T4** an `agent.started` event inside a declared run does not enter `s.starts` and does not bump `LegacyStartFloor`.
- **T5** a declared run carrying `{"idea":"…"}` refuses with the recoverable-identity error.
- **T6** a conflicting-identity run is not declarable, even though its undeclared error text is identical.
- **T7** manifest binding: adding `driver.json` in one root refuses (`differs at`); flipping one byte of `events.jsonl` refuses; removal from one root still binds the rest; removal from all roots refuses as stale.
- **T8** normalization refuses absolute paths, `..`, non-`parley-deck/runs/` prefixes, extra segments, non-hex/short digests, and two digests for one path.
- **T9** exact replay over `nil` / `[]` / unsorted / duplicated declarations is idempotent and preserves spend and epoch; a request differing only in one declared digest conflicts.
- **T10** injected publication fault before `migration-active`; the recovery preview re-inspects with **both** declaration kinds retained from the frozen request, and the exact replay completes activation.
- **T11** *(the containment test)* ordinary bootstrap still refuses on the same repository — `EnsureCycleBinding` for another idea and `ConfigureLaunchBudget` for the same one both hit the verbatim unscoped-run refusal. No declaration reached repository state.
- **T12** an undeclared inventory carries none of the new fields and keeps its exact bytes and digest.
- **T13** `InspectRunIdentities` enumerates every visible copy without refusing, separates absent from conflicting, and emits a `--declare-unscoped-run` value that `migrate inspect` accepts verbatim; a structural assertion that no apply/attest path consumes `RunIdentityReport`.
- **T14** no field of the report is readable as an action count: `Copies` is copies-of-a-file and every non-bound row carries `History: "unknown"`.

Package `app`: flag accepted for protocol `inspect`/`apply` only and rejected on every other verb listed in §5.B.7; no accounting state created by any refused or read-only command.

## 7. Flaws in my existing declared-worktree candidate (own work, correct before any apply)

1. **`LowerBoundBasis` is a scalar** (`protocol_migration_history.go:40`). The set of reasons an enumeration can be incomplete is open-ended, and this data already supplies a second. Must be `[]string`. Cost now: one test edit. Cost after an apply: a persisted-shape break.
2. **`checkDeclaredUnavailable` couples the two coverage facts** (`protocol_migration.go:54-59`): with `len(rows) == 0` it refuses any inventory whose `HistoryCoverage`/`LowerBoundBasis` is set. An unscoped-run-only import has exactly that shape, so **the current clause would wrongly refuse every valid unscoped-only import.** The clause must be split per declaration kind.
3. **`DeclaredUnavailable []string` is the wrong parameter shape** (`protocol_migration.go:30`). A second kind forces either a parallel field — with a second normalize call beside `StartedAt.UTC()` at `:241`, trivially forgettable in replay, where a missed normalization silently splits the digest — or signature churn across all four inspect sites plus recovery. One `MigrationDeclaration` type, now, while nothing is persisted.
4. **The candidate's framing generalized "unknown history" to mean "unreadable".** It does not: §3 shows the readable-bytes/unrecoverable-identity case needs the opposite retention rule. My candidate's structural-exclusion check (`checkDeclaredUnavailable` §63-78) would, if naively extended to runs, *delete* evidence we hold.

Unchanged from that candidate and still true: nothing there was compiled or executed either; §5 of it is "tests written", not "tests passing".

## 8. Residuals

1. **`<N>` is still not honestly chooseable.** The floor now excludes both an unreadable region and an unscoped run of unknown action count. `protocol_migration.go:95` enforces only `TotalActions >= LowerBound` — and since the floor drops when evidence is excluded, **the machine constrains `N` *less* exactly where the evidence is thinnest.** Recording `LowerBoundBasis` makes that visible; it does not make it safe. This is my §8 residual, unresolved and now sharper.
2. **The blob looks like a smoke fixture** — tracked since the initial CLI commit, task `"smoke implementation run"`, no cursor, no starts. Its true count is plausibly `0`. The design **refuses to act on that plausibility**: absence of `agent.started` is not proof of zero. A properly decided `fixture` classification (tracked in git, referenced by no cursor, never in an invocation record) could retire this class later; it is a separate decision and is not proposed here.
3. **Tracked-file coupling.** Because the blob is tracked, a `git checkout` between inspect and apply can change the manifest digest and refuse a frozen declaration. Correct, and operationally surprising. Document it; do not weaken it.
4. **Whole-directory manifests bind more than the scanner reads.** `runs()` reads only `events.jsonl` and `driver.json` (`:476`, `:498`), but the manifest covers the full tree, so an unrelated log file invalidates a declaration. Deliberate, fail-closed; for `20260510T194003Z` the tree is one file, so the cost here is zero.
5. **Old binaries still reject the new fields** under `DisallowUnknownFields`; the `omitempty` tamper-evidence gap (a deleted empty list is indistinguishable from no declaration) now applies to a second field. Both unchanged and open.
6. **The brief's 150/27/25 are testimony I could not verify** — the diagnostic is outside this restricted root. The design does not depend on them.

## 9. What this explicitly does not do

No identity assigned to any run. No count inferred, including not zero. No history edited, deleted, pruned, recreated or reclassified. No old-schema exemption, silent or otherwise — the undeclared path keeps its exact refusal. No `<N>` selected. No cap changed. No migration applied, no attestation, no signoff. Nothing built or executed: every claim in §5–§6 is a claim about code I have designed, not observed.

**Own path:** `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/claude-1-unscoped-history-remedy-20260916.md`
