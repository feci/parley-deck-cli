---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: independent-bounded-source-review
status: complete
reviewed-candidate: claude-1-unscoped-history-candidate-20260916.md
reviewed-source: .parley-runtime/unscoped-history-review-source/ (immutable bundle, NOT the active top-level source)
method: native Read/Grep/Glob over the bundle only. No Bash, no build, no test, no Git, no apply,
  no subagent. Every verdict below is PRIMARY (§15.2): I read the cited lines myself. The
  coordinator's compile/focused-test PASS is testimony to me; I executed nothing.
---

# Unscoped-history declaration candidate — independent source review (kimi-1)

**Scope: the candidate diff + targeted functions/tests in the immutable bundle. No whole-repo
reread, no execution, no signoff, no whole-audit acceptance. Findings are concrete only.**

## 1. What was reviewed

`candidate.diff` first, then the full text of: `internal/budget/launch_migration_history.go`,
`protocol_migration_history.go`, `protocol_migration.go`, `migration_recovery.go`,
`run_identity_inventory.go`, `protocol_migration_unscoped_test.go`, plus targeted seams:
`step_history.go` (`inspectStepRun`, `readStepHistoryFile`, `stepHistoryDirs`),
`binding.go` (`launchScopeDeclared`, `validateHistoryJSON`), `cycle_extension.go`
(`validCycleDecision`), `ledger.go` (`key`), and the CLI files `internal/app/budget_migrate.go`,
`budget_migrate_protocol.go`, `budget_migration_recovery.go`. Old-anchor test files
(`protocol_migration_declared_test.go`, `budget_migrate_declared_test.go`) are unchanged in the
bundle. The candidate note's §3/§5/§7 claims were checked against the code.

## 2. Per-axis verdicts (all PRIMARY from bundle source)

### 2.1 Old nil-declaration refusal preserved — PASS

`launch_migration_history.go:705` — `bind(e.Data[name], !declared && name == "idea" && e.Type ==
"run.created")`. For a nil declaration `declaredRuns` is empty (`:292-295`), `declared` is always
false, and the expression reduces to the old one verbatim. Cursor is bound before events in both
old and new (`:666-679`), so a `run.created` missing `idea` still refuses at its exact old site
with the old text `"historical run identity is missing or conflicting"` (`:661`) even when the
cursor already supplied a recoverable name; the post-scan refusal `"historical run needs an
explicit recoverable idea identity"` (`:742-743`) is also untouched. A declared run that turns out
to have a recoverable identity refuses at `:721-723`. Probe (1): holds.

### 2.2 Absent vs null/empty/malformed/conflicting — PASS

Scanner `bind` (`:655-665`): missing key exempt only when declared AND `run.created`/`idea`;
`"idea": null`, `""`, whitespace, non-string, and two different names all refuse in both modes.
`classifyRunIdentity` mirrors this (`run_identity_inventory.go:251-263`): null/empty/non-string →
`run-identity-conflicting`; two names → conflicting; only genuinely key-absent identity is
`run-identity-absent`. Probe (2): holds. Declared-test witnesses: conflicting identity refuses both
undeclared and declared (`protocol_migration_unscoped_test.go`, `TestDeclaredUnscopedRunPermitsAbsentIdentityOnly`).

### 2.3 Whole-directory digest over every visible copy — PASS

`RunDirectoryManifest` (`run_identity_inventory.go:48-91`): recursive walk, symlink refused
(`:62-64`), non-regular refused (`:68-70`), bounded (4096 files / 64MB / 16MB per file, `:80`).
Per-copy verification at scan time (`launch_migration_history.go:633-643`, refusal `"...differs
at <root>"`), stale declaration refused when the path is in no visible root (`:321-326`), roots
enumerated from `git worktree list` (`binding.go:167-177`). Tests cover changed byte, added file
in an added subdir, and removal from all roots. Pre-activation re-inspection repeats both checks
(`protocol_migration.go:457-463`). Probe (3): holds. (Observation, not a finding: empty
directories are invisible to the digest — the manifest lists regular files only. That matches its
own "file set" wording and is inert for evidence.)

### 2.4 Source retention / reinspection — PASS

The declared run's `events.jsonl` (and `driver.json` when present) is read via `s.file` and stays
in `Sources` per copy (`launch_migration_history.go:644, :666`); `checkDeclaredUnscopedRuns`
requires the retained `events.jsonl` source per root (`protocol_migration.go:85-89`). The
same-call evidence derivation re-reads every source and hash-compares
(`protocol_migration_history.go:139-153`). Manifest-covered files that are not sources (e.g. agent
logs) are pinned by the per-copy digest at every full inspection — inspect, apply
(`protocol_migration.go:370, :457`), and recovery preview (`migration_recovery.go:480`) — and are
never evidence inputs, so the same-call window for them is benign. Probe (4): holds.

### 2.5 No action count, no zero inference, no epoch assignment — PASS

The declared branch ends in `continue` (`launch_migration_history.go:740`): no `agent.started`
harvest, no `s.starts` entry, no `s.earlier` call. In the evidence loop the run leaves with
`identity == ""` and is dropped by the pre-existing `identity != idea` filter
(`protocol_migration_history.go:296-303`) — no evidence row under any rule, and no
`published-action-floor` of zero. `Copies == len(Roots)` counts copies of one file set
(`launch_migration_history.go:193-194`), never actions; `row.Earliest` is computed per-run
(`:724-733`) and never folded into the inventory's `Earliest` (`:253-256` untouched for declared
runs). Tests assert `LowerBound == 0`, `i.Earliest == nil`, `row.Earliest != nil`, and zero start
floor with an `agent.started` event present. Probe (5): holds.

### 2.6 Loaded-record validation, bidirectional, explicit floor — PASS

`readProtocolMigration` → `protocolMigrationInitial` re-validates every loaded record: digests
recomputed (`protocol_migration.go:159`), `checkDeclarations` (`:168`) — declaration↔rows in both
directions for both kinds plus the exact combined coverage token (`:43-61`,
`protocol_migration_history.go:61-71`) — floor recomputed from evidence (`:171`), and
`TotalActions >= LowerBound` (`:162`). Recovery journal decisions are re-validated the same way
via `recoveredInitial` → `protocolMigrationInitial` with the frozen request
(`migration_recovery.go:255-257`). Canonicalization happens once at the apply boundary
(`protocol_migration.go:305-314`); replay requires exact request equality (`:320, :349`); a
narrower declaration is a conflict (test `...NormalizationAndReplay`). Probe (6): holds.

### 2.7 Ordinary bootstrap fails closed — PASS

`inspectStepRun` (`step_history.go:105-172`) is untouched; no persisted exclusion is consulted
anywhere — declarations are request-scoped arguments only. Test asserts `EnsureStepBinding` still
refuses with its verbatim text `"historical driver idea identity is missing or malformed"` after a
declared import on the same repository, and that another idea's undeclared inspect still refuses.

### 2.8 CLI verb/kind containment — PASS

`--declare-unscoped-run` exists only on `budget migrate inspect|apply`
(`budget_migrate.go:62-63`); refused for `--kind launch` (`:76-81`); for protocol kinds the
allowed-flag map admits it only on `inspect` and attended `apply`
(`budget_migrate_protocol.go:35-51`); the recovery CLI registers no declaration flags at all
(`budget_migration_recovery.go:19-32`) — recovery inherits declarations solely from the immutable
record. Declaration digest format enforces lowercase 64-hex via `validCycleDecision`
(`cycle_extension.go:74-77`) and a canonical single-segment path under `parley-deck/runs/`
(`launch_migration_history.go:143-152`).

### 2.9 Shape/digest compatibility — PASS

Both new fields are `omitempty`; `migrationShape` skips absent `omitempty` keys
(`launch_migration_history.go:414`), so nil-declaration and declared-unavailable-only records keep
byte-identical digests (test asserts a nil-declaration marshal contains no `unscoped_runs` key).
An older binary rejects a non-empty new field via `DisallowUnknownFields`
(`launch_migration_history.go:370`) — the candidate note §7 disclosure is accurate.
`SurvivingVisibleFloor` keeps its exact value; the two new tokens exist only for the new kind.

## 3. Findings

### [MINOR] `InspectRunIdentities` offers `Declaration` values the declared inspection deterministically refuses — one identity-absent class has no declarable remedy

`classifyRunIdentity`/`InspectRunIdentities` is the surface whose stated purpose is to "tell an
operator what there is to declare" (`run_identity_inventory.go:152-155`), and it sets
`row.Declaration = <path>=<digest>` for every `run-identity-absent` row (`:209-210`). But the
apply-side scanner accepts a narrower class than the report's absent classification, so the
report's verbatim remedy fails for at least these concrete fixtures (all fail **closed** — nothing
is admitted; the report mis-advises):

(a) **Cursor-only run** — `parley-deck/runs/<name>/driver.json` = `{}` (well-formed, no idea
    keys), no `events.jsonl`. Report: `run-identity-absent`, Declaration offered
    (`run_identity_inventory.go:269-283, :315-318, :209-210`). Apply: `s.file(events.jsonl)` is
    `IsNotExist` and `driver.json` Lstat succeeds, so the scanner hard-refuses `"historical driver
    cursor has no readable run identity"` (`launch_migration_history.go:644-649`, error at `:647`).
    The declaration's digest binding at `:636-642` does run first, but the row is never recorded,
    so the declaration can never be satisfied. This class also blocks every **undeclared** inspect, so the
    operator is left with no tool-offered remedy at all — the strongest instance.

(b) **Event without a usable time** — `events.jsonl` =
    `{"type":"run.phase","data":{"action":"fixup"}}` (no `time`, no idea). Report: classification
    never inspects event times (`run_identity_inventory.go:296-307`) → absent + Declaration.
    Apply: the declared branch refuses `"historical accounting timestamp is missing or in the
    future"` (`launch_migration_history.go:726-728`).

(c) **Strict-JSON divergence** — an events line with a duplicate key
    (`{"time":"2026-05-10T19:40:03Z","type":"run.created","type":"run.phase","data":{}}`) or an
    uppercase `"Type"` key. Report uses plain `json.Unmarshal` (`run_identity_inventory.go:300`),
    which tolerates both (last-wins / case-insensitive) → absent + Declaration. Apply uses
    `migrationJSON` → `validateHistoryJSON`, which refuses duplicate or non-lowercase fields
    (`binding.go:460-461`) with `"malformed or duplicate historical fields"`
    (`launch_migration_history.go:363-364`).

(d) **Empty `driver.json`** (0-byte crash artifact). Report: `len(cursor)==0` is treated as no
    cursor (`run_identity_inventory.go:273, :315`) → absent + Declaration. Apply: `s.file`
    succeeds on the empty file, then `migrationJSON` hits EOF in `validateHistoryJSON` → refusal
    (`launch_migration_history.go:666-679` via `:362-365`), in both undeclared and declared modes.

(e) **Oversized `driver.json` (1MB < size ≤ 16MB)**. Report reads it under `maxRunManifestFile`
    (16MB, `run_identity_inventory.go:269`) → absent + Declaration. Apply reads it under a 1MB
    limit (`launch_migration_history.go:666` → `readStepHistoryFile` bound,
    `step_history.go:181-182`) → refusal.

Direction (either is acceptable; the second is smaller): teach the scanner's declared branch to
record rows for the cursor-only/both-missing cases under the same no-count exclusion, **or** make
`InspectRunIdentities` withhold `Declaration` (keep the row, append an `Uncertainty` note, the
mechanism already used for divergent copies at `:222-224`) whenever its classification relied on
tolerance the apply path does not share — missing `events.jsonl`, unusable event times,
strict-JSON failure, empty or >1MB cursor.

Severity rationale: every divergence is a false refusal, never a false admission; the report
grants nothing and the apply path re-derives everything independently. Impact is
operator-tooling correctness and remediability of one real history class (a), not evidence
integrity. Reviewers may weigh upgrading (a) if cursor-only crash runs are judged common.

## 4. Conditional scope acceptance

Acceptance of this candidate's **scope** (source-only, additive, no apply) is reasonable,
conditional on exactly these missing items — I ran nothing myself:

1. Disposition of MINOR-1: either the scanner or the report fixed per §3, or an explicit recorded
   decision that the divergence is accepted with the report's `Uncertainty` mechanism documenting
   it. Instance (a) at minimum needs the report to stop emitting a guaranteed-refused Declaration.
2. A round-trip test: for every report row with `Observation == run-identity-absent` and
   single-digest copies, either `Declaration != ""` **and** `InspectProtocolMigrationDeclarations`
   accepts it, or `Declaration == ""` with an explanatory `Uncertainty` entry. This pins all five
   instances of MINOR-1 as regression tests.
3. Coordinator-run full verification of the touched packages: `go test ./internal/budget/...
   ./internal/app/...` (not only the focused files), plus `go vet` and `-race` on
   `internal/budget`, since none of full/vet/race has run on this extension.

## 5. Explicit non-claims

No whole-audit acceptance, no signoff, no merge/apply judgment, no verdict on the deployed (older,
444-Go-file) CLI. Design proposals were read as context only. I did not execute the code; claims
above are source readings with cited lines, each falsifiable against the bundle. Claude's original
invocation timeout and the coordinator's compile/focused-test results are testimony I did not
re-verify.

## 6. Update — adversarial-probe mapping (added after the first write)

Mapping the launch task's six probes to the sections above, so later readers need not re-derive
coverage: (1) old nil-declaration refusal → §2.1 PASS; (2) null/empty/malformed/conflicting ≠
absent → §2.2 PASS; (3) full-dir hash over added/removed files/dirs/symlinks, roots, copies,
reinspection → §2.3 PASS (empty-dir invisibility noted there as inert); (4) Sources retention vs
manifest-covered non-source bytes → §2.4 PASS with the mechanism stated; (5) copies ≠ action
count, unknown-run timing not assigned to this idea → §2.5 PASS; (6) loaded-record bidirectional
validation and explicit unknown-coverage floor → §2.6 PASS. The single open item is MINOR-1 (§3):
a read-only report/apply divergence, fail-closed, with five concrete fixtures and an exact
round-trip test proposed in §4. Precision fix in this update: §3(a) now records that the
declaration's digest binding runs before the cursor-only refusal fires; the finding is unchanged.
