---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
kind: implementation-note
status: candidate
date: 2026-09-16
scope: MRW-1 stage one (read-only worktree inventory) only
---

## Summary

Read-only first stage of MRW-1: a `budget` API and a `parley budget worktree inspect --dir DIR`
CLI surface that enumerate **every** Git worktree registration — including registrations whose
working directory is gone — without pruning, recreating paths, mutating history, or changing
`launchScope(..., true)`'s existing refusal. The output carries a clock-free observation digest
so that two identical previews hash identically, and it states its own source scope and residual
uncertainty in-band.

This note is written **before** the code is final and refined afterwards. It is my own artifact;
claims in it that rest on my own earlier review are **self-owned** and carry no verdict from me
(§15.1 — an owner must not verdict a claim it owns).

## Provenance and a blocking source limitation

- `SOURCE UNAVAILABLE` — the launch task names my prior source review at
  `…/parley-worktree-inventory-h5gr6sdd/runtime/prior-review.md`. That path is **outside** this
  worktree and the file tools are `--restricted` to the checkout, so I **could not read it**. A
  repo-wide search for `MRW-1` / `MRW-` returned no matches, and no `prior-review.md` exists in
  the checkout. Everything below is therefore built from (a) the launch task's own written
  requirements and (b) `PRIMARY` reads of the code in this checkout — **not** from the prior
  review. Any requirement that lived only in that document may be unimplemented here. I do not
  claim otherwise, and I did not reconstruct the item's content from memory.
- `PRIMARY` — `internal/budget/binding.go:59-135` (`launchScope`), read in this worktree:
  `inspectHistory == true` runs `git worktree list --porcelain -z` and, at lines 101-104, returns
  `fmt.Errorf("historical worktree is unavailable: %s", worktree)` when a registered path fails
  `os.Stat` or is not a directory. That refusal is the behavior stage one must leave intact.
- `PRIMARY` — `internal/budget/binding.go` is **not** one of my six owned paths, so `launchScope`
  is read-only to me and cannot have been altered by this change. The preservation is structural,
  not merely intended.

## Design

### Why a separate enumeration instead of reusing `launchScope(..., true)`

`launchScope`'s history mode is deliberately fail-closed: an unavailable registration aborts the
whole scope resolution. Stage one needs the opposite disposition for the *same* wire data — retain
the row, mark it unknown, keep going. Reusing the refusing function would require weakening it,
which the task forbids. So `worktree_inventory.go` re-runs the same standard Git command
(`git worktree list --porcelain -z`) and reuses the same canonicalization *approach*
(`filepath.Abs` → `filepath.EvalSymlinks`, and `rev-parse --path-format=absolute
--git-common-dir`) while keeping its own non-refusing availability classification. `launchScope`
is untouched.

**Reuse of existing helpers, and where it did not apply.** The digest calls the package's own
`key()` (`ledger.go:131`, `sha256` → hex) rather than inlining a hash. The bounded regular-file
reader (`readStepHistoryFile`) and the strict JSON-input validator (`validateHistoryJSON`) are
**not** used, because this stage reads no file and parses no JSON input: its only inputs are Git's
`-z` wire bytes and `os.Lstat` results, and its JSON is output-only. Those helpers guard untrusted
on-disk state, so applying them here would be decoration, not safety.

### Exported API (`internal/budget/worktree_inventory.go`)

- `InspectWorktreeInventory(ctx, root string) (WorktreeInventory, error)` — the only entry point.
  Read-only: it runs two `git rev-parse` probes and one `git worktree list --porcelain -z`, then
  `os.Lstat`s each registered path. It creates no directory, acquires no lock, and writes nothing.
- `WorktreeInventory` — `schema`, `root`, `git` (bool), `git_common_dir`, `repo_prefix`,
  `source_scope`, `covers_repository`, `registrations[]`, `uncertainty[]`,
  `observation_sha256`.
- `WorktreeRegistration` — `registered_path` (**raw**, exactly as Git emitted it),
  `resolved_path` (canonicalized, empty when unresolvable), `head`, `branch`, `detached`, `bare`,
  `locked`, `lock_reason`, `prunable`, `prunable_reason`, `availability`, `history_coverage`.
- `WorktreeInventory.OperatorRecord() ([]byte, error)` — returns the canonical, clock-free bytes
  that a *future* attested operator record would have to quote, plus nothing else. It is a pure
  function of the inventory. There is **no** apply, attest, confirm, `--yes`, or executable path
  anywhere in this change.

**The seam for a future operator record, and nothing more.** `inspect` emits
`observation_sha256`, which is the value a later attested control would pin — the same shape the
existing controls already use (`--expected-policy-sha256` in `budget_policy.go`,
`--expected-history-sha256` in `budget_migrate.go`). That is the whole enabling step: an operator
can capture a digest now, and a later stage can require it. To keep this honest, test 9 asserts
that `--expected-observation-sha256` is **rejected today** — the intended seam is documented by a
refusal, not by a half-built flag.

### Availability classification (unknown history, never zero)

`os.Lstat` on the raw registered path yields exactly one of four states:

| `availability`  | condition                              | `history_coverage` |
| --------------- | -------------------------------------- | ------------------ |
| `available`     | exists, is a directory                 | `local`            |
| `missing`       | `os.IsNotExist`                        | `unknown`          |
| `not-directory` | exists, is not a directory             | `unknown`          |
| `indeterminate` | `Lstat` failed for any other reason    | `unknown`          |

`history_coverage` is a **coverage classification, never a count**. No count, cap, floor, or
attestation field exists in the output, for available or unavailable rows alike — so an
unavailable registration can never be read as "zero", and nothing here infers the authorization
that the launch task withholds. `missing`, `not-directory` and `indeterminate` rows are all
**retained** in `registrations[]` and each appends a line to `uncertainty[]`.

`indeterminate` is a deliberate fourth state rather than a refusal: an `EACCES` on one registered
path is neither malformed registration input nor broken Git authority, and refusing there would
defeat "enumerate ALL registrations". It is never reported as `available`.

### Refusals (fail closed)

Rejected as malformed / ambiguous / out-of-scope registration input:

- a record not led by a `worktree ` attribute, or a `worktree` attribute with an empty value;
- a registered path that is not absolute, or that is byte-identical to an earlier one;
- a repeated attribute key inside one record;
- an **unknown** attribute key (fail-closed, see Limitations);
- `branch` together with `detached`; `bare` together with `HEAD` or `branch`;
- a record with neither `bare` nor `HEAD`;
- a `HEAD` value that is not a 40- or 64-character lowercase hex object name;
- a `branch` value that is not a `refs/` ref;
- trailing non-empty data after the final record separator;
- output above a 16 MiB bound, or more than 4096 registrations.

Rejected as broken Git authority, reusing `launchScope`'s own discipline: when
`git rev-parse --git-common-dir` fails but a `.git` entry exists at `root` or any ancestor, the
call returns an error instead of degrading to a non-Git answer. `ctx.Err()` is checked first so a
cancelled context is reported as cancellation.

### Non-Git root

A root with **no** `.git` at it or any ancestor is not an error and not a repository: the
inventory returns `git: false`, `covers_repository: false`, `source_scope:
"single-directory/non-git"`, **zero** registrations, and two explicit `uncertainty[]` lines saying
that no Git worktree authority was found and that this result describes only the one directory
inspected — it does **not** establish repository-wide history coverage. `git` and
`covers_repository` are both inside the digest, so a non-Git preview cannot collide with a Git
one.

### Observation digest

`observation_sha256 = key(canonical JSON of the load-bearing facts)`, where the canonical form is
the closed struct `observedWorktrees`: `schema`, `git`, `covers_repository`, `source_scope`,
`git_common_dir`, and every registration field above — with the registration list **sorted by raw
registered path** so the digest does not depend on Git's enumeration order.

Three observed values are deliberately **excluded**:

- the observation **time** — there is no clock anywhere in the digest input, so re-inspecting
  unchanged state yields a byte-identical digest;
- `root` and `repo_prefix` — both describe *where the caller stood*, not what the repository
  registers. Including `repo_prefix` would have made the same repository hash differently when
  previewed from a subdirectory, which is exactly the instability the requirement rules out.
  `git_common_dir` carries the repository identity instead.

Consequently `repo_prefix` is orientation only, so a failing `rev-parse --show-prefix` (a root
with no work tree, i.e. a bare repository) is **not** fatal: the field is left empty, an
`uncertainty[]` line records it, and the registrations are still enumerated. Authority is already
established by the `--git-common-dir` probe, and the `worktree list` read below must still
succeed or the call refuses. Emitted `registrations[]` keeps Git's original emission order; only
the digest input is sorted.

## Exact files

| Path | Status | Content |
| --- | --- | --- |
| `internal/budget/worktree_inventory.go` | new | `WorktreeInventory`, `WorktreeRegistration`, `InspectWorktreeInventory`, `OperatorRecord`, porcelain `-z` parser, availability classifier, digest |
| `internal/budget/worktree_inventory_test.go` | new | package-level fixtures incl. the real-Git missing-worktree witness |
| `internal/app/budget_worktree.go` | new | `runBudgetWorktree` — `inspect` only, `--dir` only |
| `internal/app/budget_worktree_test.go` | new | CLI integration + refusal of mutation verbs and extra arguments |
| `internal/app/budget.go` | edited | dispatch `worktree` in `runBudgetControl`; one added usage line |
| `parley-deck/…/claude-1-worktree-inventory-candidate-20260916.md` | new | this note |

No other file is touched. `internal/budget/binding.go` in particular is unmodified and unowned.

## Expected tests

In `internal/budget/worktree_inventory_test.go`:

1. `TestWorktreeInventoryRetainsRemovedLinkedWorktree` — **the real-Git witness**. `git init`,
   an empty commit, `git worktree add` a linked worktree, then `os.RemoveAll` the linked path
   (leaving the registration in the admin directory). Asserts: the removed registration is still
   present with `availability: "missing"`, `history_coverage: "unknown"`, its raw
   `registered_path`, its `HEAD` and `branch: refs/heads/linked` intact, and no `resolved_path`;
   the surviving row is `available`/`local`; `uncertainty` is non-empty; the emitted key sets of
   the inventory and of the retained row are **exactly** the documented ones, so no count, cap,
   floor or attestation field exists; `git worktree list --porcelain` still reports the removed
   registration afterwards (no prune); the removed path was not recreated; the admin directory
   still holds the same number of entries; no `.parley-runtime` or `parley-launch-budgets` was
   created; and — separately — `ConfigureLaunchBudget` on the same root still fails with
   `historical worktree is unavailable`, proving `launchScope(..., true)`'s refusal survives.
2. `TestWorktreeInventoryDigestIsStableAndClockFree` — two inspections of unchanged state return
   an identical `observation_sha256` and identical `OperatorRecord()` bytes; the digest input
   type's key set is asserted to be exactly the six documented fields (a structural check, so no
   clock and no caller location can be added unnoticed); and inspecting the same repository from
   a linked worktree and from a subdirectory yields the same digest.
3. `TestWorktreeInventoryDigestTracksRegistrationAndAvailability` — adding a worktree changes the
   digest and grows the registration list; renaming a registered directory away changes the
   digest again and flips that row to `missing`; renaming it back restores the exact prior
   digest. (A `rename` round-trip, not `RemoveAll` + `MkdirAll` — recreating a bare directory
   would leave Git still reporting the registration `prunable`, so it is not the prior state.)
4. `TestWorktreeRegistrationRefusesMalformedAuthority` — 17 table cases over the unexported
   parser with fabricated wire bytes: unknown attribute, duplicate attribute, duplicate path,
   relative path, empty path, valueless `worktree`, `branch`+`detached`, `bare`+`HEAD`, neither
   `bare` nor `HEAD`, non-hex HEAD, 39-char HEAD, non-`refs/` branch, attribute outside a
   record, a flag carrying a value, a second `worktree` inside a record, empty output, and an
   unterminated trailing record.
5. `TestWorktreeRegistrationAcceptsStandardAttributes` — the counter-test to 4: the strict parser
   still accepts every shape real Git emits (bare, branch, detached, `locked` with and without a
   reason, `prunable` with a reason, a 64-char SHA-256 object name) and preserves the lock and
   prunable metadata.
6. `TestWorktreeInventoryRefusesBrokenGitAuthority` — `root/.git` written as a garbage regular
   file so `rev-parse` fails while a `.git` entry exists → error naming the authority, and no
   fallback inventory.
7. `TestWorktreeInventoryNonGitRootIsExplicit` — a plain temp dir yields `git: false`,
   `covers_repository: false`, no registrations, the explicit single-directory `source_scope`, a
   non-empty `uncertainty`, and a digest distinct from a Git root's. (It skips rather than fails
   if the temp dir is itself inside a repository.)

In `internal/app/budget_worktree_test.go`:

8. `TestBudgetWorktreeInspectCLIReportsRetainedRegistration` — real Git root with a removed
   linked worktree, run with `supported=false, attended=false` (so it also proves no attended
   gate was introduced); exit 0; stdout is exactly one JSON object that round-trips into
   `budget.WorktreeInventory` under `DisallowUnknownFields`, contains the retained `missing`
   row and a non-empty `uncertainty`, did not recreate the path, and is byte-identical on a
   second run.
9. `TestBudgetWorktreeRefusesMutationAndExtraArguments` — `apply`, `attest`, `prune`, `repair`,
   `recover`, `remove` and bare `worktree`, plus `inspect` with a positional argument, `--yes`,
   `--apply`, `--attest`, `--prune`, `--decision-id`, `--idea` and
   `--expected-observation-sha256`, all exit 2 **even with `supported=true, attended=true`**,
   emit nothing on stdout, and create no runtime or policy state; and the usage text itself is
   asserted not to name any verb this stage does not implement.
10. `TestBudgetWorktreeInspectRefusesBrokenAuthority` — the CLI exits 1 and writes nothing on
    stdout for a corrupt `.git`.

I did **not** run any of these. The coordinator runs the tests; §15.1 also bars me from
verdicting my own output, so this section is an expectation, not a result.

## Limitations

1. **The prior review was unreadable** (see Provenance). Stage-one scope here is reconstructed
   from the launch task text alone.
2. **No test was executed.** No `PASS`/`FAIL` claim is made, and no acceptance is claimed. Any
   statement about what the code does is a `PRIMARY` read of source I just wrote — which, being
   self-owned, I cannot verdict.
3. **An unknown Git attribute is a refusal, not a warning.** If a future Git adds an attribute to
   `worktree list --porcelain`, `InspectWorktreeInventory` fails closed rather than silently
   dropping a fact from a digest that is meant to be load-bearing. That is a deliberate
   availability/compatibility trade and a real upgrade hazard; the fix is to extend the parser.
4. **No recovery.** Nothing here recovers, repairs, prunes, migrates, counts, caps, attests, or
   re-creates a missing root. `OperatorRecord()` only *renders* bytes; no flag consumes them.
5. **`indeterminate` is genuinely unresolved**, not a benign default. An inventory containing one
   is not a clean preview and must not be treated as an availability attestation.
6. **The digest covers what this stage observes, not the repository's content.** It says nothing
   about commits, reflogs or branch reachability inside any worktree, available or not.
7. **Symlinked registered paths.** `availability` is decided by `os.Lstat` on the raw path, so a
   registration pointing at a symlink is `not-directory`, not `available`, even if the target is
   a usable directory. `resolved_path` is left empty in that case. This is intentionally stricter
   than `launchScope`'s `os.Stat`, and it means the two functions can disagree about one such
   path — a divergence a later stage should reconcile deliberately rather than by accident.
8. **Concurrency.** The inventory is an unlocked point-in-time read; a worktree added or removed
   between the `git` call and the `Lstat` loop produces a digest for a state that no longer holds.
   It takes no lock by design (a read-only preview must not block writers).
9. **The 16 MiB output bound is checked after buffering, not during.** `exec.Cmd.Output()` reads
   the whole of Git's output into memory first, so the bound refuses to *parse* an oversized
   registration list but does not prevent it being read. This matches the existing
   `launchScope` call, so it is not a new exposure, but it is not a true streaming bound either.
10. **Availability and prunability are independent, and Git decides prunability.** A restored
    directory that lacks its `.git` pointer file is `available` yet still `prunable`; the row
    reports both verbatim rather than reconciling them. This stage has no opinion about which is
    authoritative — a later stage must decide deliberately.
11. **`repo_prefix` is outside the digest**, so two previews that agree on `observation_sha256`
    may still have been taken from different directories. The digest answers "is this the same
    registration state?", never "was this taken from the same place?".
12. **The three amendment refusals and the pending user backup clarification are untouched**, and
    nothing in this change depends on or anticipates their resolution. The full original FINAL
    remains binding.
