---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: test-completion-and-correction
status: written-unexecuted
completes: claude-1-unscoped-history-candidate-20260916.md (that note's §10 points at a §11 it never
  received; this note is that missing completion record, written in a separate file so the original
  keeps its exact bytes)
method: native Read/Write over this worktree only. No shell, Git, build, test, CLI, migration, apply,
  attestation, provider, history-apply or subagent in this launch. Nothing here was compiled or
  executed by me.
---

# Unscoped historical runs — test completion record and correction

**Tests only. No production file edited in this launch. No apply, no attestation, no count chosen,
no identity assigned, no current-source acceptance or amendment.**

## 1. Correction first: the previous invocation timed out, it did not succeed

My previous invocation on this work (`de1b732d`) **ended in a 900-second timeout, not in a
completion.** The nine Go source files and the note
`claude-1-unscoped-history-candidate-20260916.md` it had written were retained by the coordinator,
but they are the state of a run that was cut off mid-task, and they must be read that way:

- That note still carries `status: in-progress`, and its **§10 refers to a "§11" that does not
  exist** in it. The completion record it promised was never written. That is the artifact of the
  timeout, not an editorial choice.
- I have **not** edited that note. It keeps its original bytes, dangling §11 reference included, so
  the record of what the timed-out run actually left behind stays intact. This separate note is the
  missing §11.
- No claim in that note about behaviour was verified by execution then, and none is now.

## 2. What this launch was, and what it deliberately was not

Scope given to me: finish the **missing tests only**, and write this note. Production source is
**immutable** in this launch and is under independent review by kimi-1 while I write; I changed no
production file, and none of the assertions below was weakened to accommodate one.

Verification status I can honestly report, per §15.2 — this is `PRIMARY` as to *what source exists
in this worktree* (I read it) and `RECALL` as to *whether any of it passes*, because I ran nothing:

- Coordinator testimony, **not mine and not verified here** (`SECONDARY`, and its execution
  provenance is the coordinator's, not mine): the immutable tree of 448 Go files / 456 files
  compiled in 2.390s, and 22 focused events passed in 3.653s. I did not run, observe or reproduce
  either. **No broad run, no `-race`, no `go vet` and no independent acceptance exists** for this
  work as of this note.
- The 7 pre-existing tests in `internal/budget/protocol_migration_unscoped_test.go` are reported
  passing by the coordinator. Again testimony; I read the file, I did not run it.
- **Everything I wrote in this launch is unexecuted.** I make no claim that any test below passes,
  compiles, or even that a helper I call resolves. They are tests I wrote and read, nothing more.

## 3. Tests written in this launch — exactly, and only, these

Three files touched, all test files. **Nine new test functions.** No production file was opened for
writing, and the 7 pre-existing unscoped tests and 2 pre-existing declared-app tests were left
byte-identical; my additions are appended after them.

| File | Status | Added |
|---|---|---|
| `internal/budget/run_identity_inventory_test.go` | **new** | 3 helpers + 6 tests |
| `internal/budget/protocol_migration_unscoped_test.go` | appended | 1 test (7 existing untouched) |
| `internal/app/budget_migrate_declared_test.go` | appended | 1 const + 1 helper + 2 tests (2 existing untouched) |

## 4. Per-test inventory

**`internal/budget/run_identity_inventory_test.go`** — `RunDirectoryManifest` / `InspectRunIdentities`
had no test file at all before this launch.

1. `TestInspectRunIdentitiesSeparatesAbsentFromEveryOtherClass` — the declarability boundary. Five
   run directories in one repository: absent, conflicting (two names), unreadable (unparseable
   line), empty-string identity, and bound. Asserts each class is reported as itself, that **only**
   the absent class carries a `Declaration`, that the emitted value is byte-equal to what the
   accounting boundary's own `declareRun` helper builds, that a bound row alone says `scoped` with
   a name while every other row says `unknown` with none, that the report keeps its non-count
   disclaimer, that two inspections of one unchanged tree produce one digest, and that the operator
   record's fields are exactly `{schema, scope, coverage, rows}`.
2. `TestInspectRunIdentitiesGroupsCopiesAndWithholdsDivergentDeclarations` — real Git worktrees.
   Identical copies group into one row with `Copies: 2`, sorted roots and one covering declaration;
   a one-byte divergence splits them into two rows, **withholds the declaration from both**, and
   discloses the divergence in `Uncertainty`.
3. `TestRunDirectoryManifestBindsContentNotLocation` — the property that makes one declaration cover
   every copy: the same file set at two different absolute locations digests identically. Also
   sorted slash-relative recursive paths, real `Bytes`/`SHA256` on a nested entry, and an unrelated
   added file changing the digest (fail-closed, deliberate).
4. `TestRunDirectoryManifestRefusesUnhashableEntriesAndRoots` — symlink to a file and symlink to a
   directory are both refused (`contains a symlink`), never followed and never skipped; a symlinked
   run directory and a regular file are both refused as roots (`is not a real directory`); a missing
   directory reports as missing. Skips cleanly where symlinks are unavailable.
5. `TestRunIdentityDiscoveryAndMigrationScannerDisagree` — the two divergences in §5, asserted as
   current behaviour on both surfaces. This is also the requested nil-declaration case: a
   `run.created` with no idea key in a run whose `driver.json` names one still refuses verbatim.
6. `TestInspectRunIdentitiesGrantsNothingToOrdinaryCallers` — discovery writes nothing (checked
   before any other call), changes no refusal for `InspectProtocolMigration` or `EnsureCycleBinding`,
   and enumerates an empty repository as empty without inventing a row.

**`internal/budget/protocol_migration_unscoped_test.go`**

7. `TestDeclaredUnscopedRunExemptsOnlyTheAbsentKey` — the declaration exempts exactly one case.
   Unparseable bytes refuse, an event with no `type` refuses with its exact text, and a declared
   directory holding an entry the walk cannot hash refuses at the digest check before any byte is
   read as evidence. None of these is an established absence of identity.

**`internal/app/budget_migrate_declared_test.go`** — the original file had no unscoped-run case.

8. `TestBudgetMigrateDeclaredUnscopedRunIsProtocolInspectApplyOnly` — real CLI plumbing through
   `runBudgetPlatformControl`. The flag value is computed from the fixture's actual file set by the
   exported `budget.RunDirectoryManifestDigest`, never fabricated. Asserts: worktree declaration
   alone still refuses (the run flag is load-bearing); both flags compose into
   `SurvivingScopedFloor` with `declared-incomplete` coverage and both row kinds retained; the row
   round-trips to the declared digest with `unknown` history; and the flag is refused with exit 2
   for `--kind launch` inspect and apply, `migrate recover inspect|apply`, `migrate declare` and
   `worktree inspect`. Malformed values (no `=`, absolute path, empty) are refused rather than
   repaired. No refused or read-only command creates accounting state.
9. `TestBudgetMigrateDeclaredUnscopedRunApplyReplaysExactly` — an attended apply carrying both
   declarations replays exactly twice, loads as a binding with the frozen cap, and refuses when
   **either** declaration is dropped. `--total-actions` is the observed floor from the preview: this
   fixture chooses no count and asserts no zero. The apply is local to a throwaway `t.TempDir()`
   repository — it is not, and cannot become, an amendment of any real history.

## 5. Problems the tests surfaced in production source

Recorded, not fixed. No assertion was weakened and no production file was edited to accommodate
these; test 5 asserts current behaviour exactly, so it documents rather than hides them.

1. **A cursor-identified run with a keyless `run.created` is neither scannable nor declarable.**
   `classifyRunIdentity` binds from `driver.json` first, so `InspectRunIdentities` reports the run
   `bound`/`scoped` and offers no declaration; `runs()` still requires the `idea` key on
   `run.created` (`launch_migration_history.go:705`) and refuses with `historical run identity is
   missing or conflicting`. The operator gets a refusal and no flag value that resolves it. Whether
   the scanner should accept the cursor's identity is a design question I am not answering here.
2. **Discovery can emit a declaration the scanner cannot honour.** For a run holding only
   `driver.json` and no `events.jsonl`, discovery reports `run-identity-absent` and emits a
   copy-pasteable `--declare-unscoped-run` value; `runs()` rejects the missing `events.jsonl` at
   `launch_migration_history.go:646-648`, *before* the declared branch at `:717`, so that exact
   declaration refuses with `historical driver cursor has no readable run identity`. The discovery
   surface's stated purpose is to emit values `migrate inspect` accepts; here it does not.
3. **NIT.** `repeatedFlag.Set` returns `declared worktree path is empty`
   (`internal/app/budget_migrate.go:26`) for an empty `--declare-unscoped-run` too. Shared message,
   wrong noun for the second flag.

Not asserted, deliberately, and disclosed rather than guessed: whether a protocol apply writes under
`.git/parley-launch-budgets`. I could run nothing, the pre-existing apply test makes no such claim,
and an unverifiable guess in either direction is worse than the gap. Also untested by choice: the
`maxRunManifestFile` (16 MiB), `maxRunManifestBytes` (64 MiB) and `maxRunManifestFiles` (4096)
ceilings — exercising them costs real I/O to test constants rather than a semantic boundary. The
`ModeSymlink` refusal covers the "entry this walk cannot hash" class; a non-regular, non-symlink
entry (FIFO, socket) is **not** covered, since creating one portably needs an import this package's
tests do not carry.

## 6. How to run these, for the coordinator

I ran none of this. Suggested order — focused first, so a compile error or a bad fixture surfaces
before a long run, then the final combined broad pass:

```
go test ./internal/budget/ -run 'RunIdentit|RunDirectoryManifest|DeclaredUnscopedRun' -count=1
go test ./internal/app/    -run 'BudgetMigrateDeclared'                                -count=1
go build ./... && go vet ./... && go test ./... -count=1        # then broad, and -race separately
```

Both new budget tests that need multiple roots use the existing `worktreeRepoFixture` and skip
without `git`; the manifest tests skip without symlink support.

## 7. What this explicitly does not do

No production file edited. No identity assigned to any run. No count inferred, including not zero.
No history edited, pruned, recreated or reclassified. No `<N>` selected, no cap changed, no
migration applied, no attestation, no signoff. No fixture in these tests performs an actual
amendment against real repository history.

**Own path:** `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/claude-1-unscoped-history-test-completion-20260916.md`
