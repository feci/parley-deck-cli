---
idea: meta-protocol-change-evidence-first-efficiency
author: zcode-1
created: 2026-09-16
status: tests-written
---

# zcode-1 — grouped launcher test completion (recovery, 2026-09-16)

Recovery session for the four missing tests in
`.parley-runtime/grouped-amendment-launcher-20260916/main_test.go`. Tests are
WRITTEN in this session; the coordinator executes them later. This session ran
no build, no test, no git operation and no subprocess (per task constraints),
so this note claims no green run and no acceptance.

## Prior failure retained (invocation 08cc206d)

- Previous invocation 08cc206d timed out at 900 s; no test was added there.
- Coordinator compile check FAILED in 0.896 s: `main_test.go:13: "time"
  imported and not used` (deadline-test stub imported `time` with no body).
- All four production corrections were already present in `main.go`; the
  prior note's checklist was left unfinished. That record stands; this note
  does not retroactively complete it.

## Tests written this session (appended to main_test.go)

1. `TestRefusesTamperedMigratedLedgerBeforeDiscovery` — builds the cycle
   binding via `budget.InspectProtocolMigration` + `budget.MigrateProtocolBudget`
   (total = max(1, visible floor), maximum 3, epoch not after earliest,
   explicit fixture reason/decision), removes the imported charge from
   `ledger/ledger.json` keeping valid JSON, asserts preview AND execute refuse
   before any discovery (spy must stay uncalled), with no source mutation and
   no new charge.
2. `TestAmbientConfigLayerBindsDigestAndExplicitOverrideRestoresEnv` — ambient
   `PARLEY_HEADLESS_AGENT_CONFIG` with no `-config`: env layer is hashed into
   the plan digest; editing only that file's model/comment bytes makes
   `-execute -expect-digest` refuse ("digest changed") before discovery
   without spending a cycle. Narrow tail: explicit `-config` override executes
   successfully and restores the former ambient env value afterwards.
3. `TestExecuteRunsDiscoveryUnderSharedDeadline` — injected `discoverAgents`
   captures the discovery context deadline; asserts it is set, at most 900 s,
   at least 10 minutes remaining; no sleeps; execute succeeds with one cycle.
4. `TestFailedPeerReportsUnsuccessfulWithoutLosingOrganizer` — one peer
   (kimi-1) uses gated `TestFakeFailedPeerHelper` exiting 7 with marker
   `grouped-failed-peer` (ordinary runs are a no-op); report must be code 2 /
   ok=false with organizer_preserved=true, organizer bytes unchanged, exactly
   one group cycle spent.

## Checklist (reflects the original timeout, not a successful implementation)

- [x] This note written first, before any test code.
- [x] Four tests appended to `main_test.go`; unused `time` import now used by
  the migration-epoch and deadline tests. Original tests and assertions were
  not altered (append-only, review-verified — no compile run in this session).
- [ ] `go test` executed — NOT in this session; the coordinator runs it later.
- [ ] Coordinator compile gate green again (main_test.go:13).

## Session log

- Read only the two allocated local files plus the exact budget migration
  functions and cycle/ledger/config layer loading needed by the tests.
- No production files touched. No product defect discovered.
