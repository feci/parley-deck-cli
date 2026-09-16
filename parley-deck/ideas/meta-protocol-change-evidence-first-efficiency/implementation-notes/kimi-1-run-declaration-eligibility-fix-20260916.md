---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
kind: scoped-remedy
status: complete
addresses: kimi-1-unscoped-history-source-review-20260916.md §3 MINOR-1 (review06c70e21, retained unchanged)
allocated-writes:
  - internal/budget/run_identity_inventory.go
  - internal/budget/run_identity_eligibility_test.go
  - parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/kimi-1-run-declaration-eligibility-fix-20260916.md
method: native Read/Edit/Write only. No Bash, no tests run, no Git, no subagents, no providers, no
  apply. Code claims below are source readings (PRIMARY); compile/test PASS is NOT claimed — I ran
  nothing.
---

# Run-declaration eligibility fix — MINOR-1 smaller safe-direction remedy (kimi-1)

**Scope: only the smaller remedy from my MINOR-1 — make the read-only report withhold `Declaration`
plus record an explicit `Uncertainty` wherever its tolerant classification exceeds what the strict
declared scanner admits. Scanner/migration/grant files unchanged. Claude concurrently owns
`run_identity_inventory_test.go`, `protocol_migration_unscoped_test.go`,
`app/budget_migrate_declared_test.go`; I do not touch them.**

## 1. What is being changed (plan, then result)

`internal/budget/run_identity_inventory.go`:

1. New unexported helper `runDeclarationEligibility(root, name) (reason string, eligible bool)`
   that re-validates one identity-absent run under exactly the strict declared scan's own machinery
   and bounds — `migrationScanner.runs`' declared branch mirrored per run: parent-chain Lstat +
   `readStepHistoryFile` (events 16MiB, cursor 1MiB), `migrationJSON` (duplicate/non-lowercase keys
   refused, required-field shape), untyped-event refusal, strict identity bind, and the
   missing/future event-time refusal. No permissive JSON logic duplicated; the scanner itself is
   not relaxed and declaration eligibility is not widened.
2. `InspectRunIdentities` offers `row.Declaration` only when that per-copy eligibility holds;
   otherwise it keeps the row (classification and unknown history preserved, no false zero/count)
   and appends an explanatory `Uncertainty`, the same mechanism already used for divergent copies.
   Eligibility is evaluated per actual run/file set (per copy root), never by a whole-root
   inspection that could hide other valid rows.

`internal/budget/run_identity_eligibility_test.go` (new): round-trip behavioral tests for all five
MINOR-1 fixtures plus the valid control — for every identity-absent row, either the report offers
a Declaration and `InspectProtocolMigrationDeclarations` accepts exactly that value, or
`Declaration` is empty, `Uncertainty` explains, and the strict scan keeps its original refusal
verbatim. The known-valid 117-byte May 10 event (`unscopedRunEvent`) must remain declarable.

## 2. Fixtures covered (from MINOR-1 §3 (a)–(e))

(a) cursor-only / missing events (both missing included), (b) missing/future event time,
(c) duplicate/non-lowercase JSON keys, (d) empty `driver.json`, (e) `driver.json` > 1MiB.
Control: valid `unscopedRunEvent` (117 bytes, 2026-05-10) stays declarable and round-trips.

## 3. Honest interactions reported up front

- `run_identity_inventory_test.go` `TestRunIdentityDiscoveryAndMigrationScannerDisagree` (Claude's
  file) currently asserts the OLD divergence for the cursor-only case
  (`bareRow.Declaration == ""` fails the test). My change flips that expectation; Claude's
  concurrent update of that test is required. Its first half (cursor-identified bound run) is
  unaffected.
- The passing read-only v2 preview (9.817s, lower_bound 1 over 322 sources, 25 copies of one
  unknown run + 2 unknown roots, registry/budget 0 unchanged) is a valid observation of one happy
  path; it does not erase fixtures (a)–(e), which this remedy withholds declarations for.
- Eligibility is per-run: inventory-global caps (64MiB total sources, `maxMigrationItems` sources)
  can still refuse an otherwise-eligible declaration at apply; the report never claimed otherwise.
- No compile, no test run, no vet by me; coordinator verification still required.

## 4. Result

Landed exactly as planned in §1, in the two allocated code paths.

**`internal/budget/run_identity_inventory.go`** (production helper, only file changed):

- `runDeclarationEligibility(root, name) (string, bool)` — new unexported helper mirroring, for one
  run directory only, the declared branch of `migrationScanner.runs`
  (`launch_migration_history.go:622-778`): parent-chain Lstat + `readStepHistoryFile` at the same
  bounds (events 16MiB, cursor 1MiB — the bound that makes fixture (e) diverge); missing
  `events.jsonl` resolved exactly as the scan does (cursor present → "cursor with no readable run
  identity"; both missing → "not present in any visible root"); cursor and every event line parsed
  with the scan's own `migrationJSON` (duplicate/non-lowercase keys, required-field shape,
  trailing data — fixtures (c) and (d)); the same untyped-event refusal; the same strict
  `idea`/`idea_slug` bind with a recovered identity refusing outright; the same missing/future
  event-time refusal (fixture (b)); the same `maxMigrationItems` event cap. No permissive JSON
  logic was duplicated (the tolerant `classifyRunIdentity` read is untouched and stays the
  enumerator), declaration eligibility was not widened (a recovered identity still refuses), and
  the scanner was not relaxed (no migration/grant/scanner file was edited).
- `InspectRunIdentities` now offers `row.Declaration` only when eligibility holds for the copy
  being enumerated (re-checked per copy root, AND-combined across copies). Otherwise the row is
  kept with its classification, detail and `UnknownHistory` intact, `Declaration` is empty, and
  `Uncertainty` gains one entry naming the run path and the concrete reason — the mechanism
  already used for divergent copies. No zero, no count, no coverage change; an eligible row's
  offered value is byte-identical to before (`path=digest`), so the valid control and every
  already-declarable row are unchanged. Doc comments on `InspectRunIdentities` and
  `RunIdentityRow.Declaration` now state the agreement rule.

**`internal/budget/run_identity_eligibility_test.go`** (new): three behavioral tests —
`TestRunDeclarationEligibilityWithholdsWhatTheStrictScanRefuses` (all five MINOR-1 fixtures, with
both-missing, future-time and uppercase-key variants: report stays tolerant, classification stays
`run-identity-absent`/unknown history, `Declaration` empty, `Uncertainty` explains, and
`InspectProtocolMigrationDeclarations` on the would-be value keeps its original refusal verbatim —
asserted per fixture); `TestRunDeclarationEligibilityAdmitsTheValidControl` (the 117-byte May 10
`unscopedRunEvent` stays declarable, the offered value is accepted by the declared inspect, row
retained with unknown history on its actual root, floor 0); and
`TestRunDeclarationEligibilityJudgesEachRunOnItsOwnFiles` (a withheld run never hides or renames a
valid row beside it, and the explanation names only its own run).

**Honest late notes:** (1) Claude's concurrent edit of
`run_identity_inventory_test.go:284-324` (`TestRunIdentityDiscoveryAndMigrationScannerDisagree`)
must flip the cursor-only half to the new expectation — with my change it fails as written; its
cursor-identified-bound half is unaffected, as are the other three tests in that file. (2) The
read-only v2 preview PASS (9.817s, lower_bound 1, 322 sources, 25 copies + 2 unknown roots,
registry/budget unchanged) remains a valid happy-path observation; fixtures (a)–(e) are distinct
shapes it never exercised. (3) Inventory-global caps (64MiB / `maxMigrationItems` sources) can
still refuse an otherwise-eligible declaration at apply — per-run eligibility cannot see them and
the report never claimed otherwise. (4) Shared test helpers (`runIdentityRoot`, `writeRunFile`,
`writeRunDirectory`, `rowsByPath`, `declareRun`, `unscopedRunEvent`) are reused from the existing
package test files, not redeclared.

**Non-claims:** no compile, test, vet or race run by me (no Bash permitted); coordinator full
`go test ./internal/budget/... ./internal/app/...` + vet/race verification is still owed, along
with Claude's concurrent test-file update. No product endpoint or type changes; no apply; no
signoff; no whole-audit acceptance.
