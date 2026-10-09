---
idea: meta-protocol-change-driver-unstall
implementer: codex-1
status: implemented
started: 2026-10-09
updated: 2026-10-09
fix-up-cycle: 0
---

## Implementation plan / checklist

Before product edits, implement the signed FINAL from design PR80 (merged). codex-1 is organizer, participant, FINAL drafter and implementer; zcode-1 independently reviews and goal-checks in separate processes. The controlling brief authorizes all transitions and release; no standing machine organizer overrides this idea.

- [x] D1: add a strict repository/deck-scoped append-only legacy decision store and read-only preview/attended apply commands. Bind closed structural-only histories by full manifest digest; reuse budget scope/path/JSON/guard/durability primitives. Preserve exact replay and unknown-history disclosure.
- [x] D1: adopt/freeze declaration payloads and hashes in first cycle policies, revalidate all visible declared copies on every use, reject additional unknown history and preserve charges/caps/coverage and per-idea migrations.
- [x] D2: strict track/config timeout resolution and uniform existing participant-step execution for goal checks, including protected checkers, without widening membership/close policy. Freeze effective timeout for restart of the same step if necessary; preserve semantic FAIL finality and failed PASS refusal.
- [ ] Add AC1–AC8 adversarial regression tests; run focused tests, full Go suite/build/vet and relevant race tests.
- [ ] Synchronize minimal normative timeout/retry wording across CLI live/bootstrap, skill reference/guidance and staged core2.18.0 based on staged2.17.0. Update metadata/changelog/version CLI1.54.0 and skill2.18.0; run skill/drift/packet checks.
- [ ] Open implementation CLI/skill PRs, obtain Zcode independent review, signed review consensus and at most five fix-up cycles. Both participants ACCEPT, current-tree AC evidence and a fresh Zcode goal PASS are required. Record the mechanical reviewer gate outcome and brief-authorized attended close if needed.
- [ ] Merge with merge commits; GitHub releases; both Homebrew formulae; skill-only WinGet PR; install all managed plus four generic skill targets and verify every SKILL.md hash. Zcode verifies channels in a separate process.
- [ ] Prepare exact one-time D6 activation request/command and every still-unpublished npm/core command in version order; record pending status, deferrals, usage and verified shared memory outcome in release handoff.

## Summary of work

Plan recorded before product changes. Design consensus includes both current ACCEPTs. Scope is D1 and D2 only. Product fix does not imply local activation: D6 remains blocked until the owner attends the exact one-time apply. The agent must not allocate a terminal or fake attendance. npm/core publishing also retains shipped owner boundaries.

## Deviations from FINAL.md

None.

## Decisions & rationale

2026-10-09 · codex-1: Keep existing request-scoped migration semantics intact. New durable authority is a separate operator surface. Use an immutable retained adoption alongside policy so lost/changed provenance refuses rather than becoming an undeclared bootstrap.

## Progress / current state

D1 and D2 implemented. Focused D1/D2 checks passed, including a real 121-second independent-process fixture crossing the former 120-second ceiling. Full Go and skill suites are running; build/vet passed. Next: finish checks, publish the review commit and obtain independent review. Full protocol context was read and source SHA256 acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137 attested in round artifacts. The final normative patch will be reviewed under that ratified design.

## Validation / AC evidence

Pending independent review. Focused D1 passed in 9.095s; D2 app cases passed in42.742s; real long-process witness passed in121.281s. Build and go vet passed. Baseline budget witnesses passed in source-context/baseline-budget-tests.log. Actual D6 reproduced; no real repository legacy authority has been applied.

## Outcomes & surprises

The native fallback was required for round02 because D6 blocked both continue and measured cross-review dispatch. Driver consensus draft/signoff/finalize worked. Kimi had two genuine HTTP403 auth failures and was manually excluded under the brief; immutable automatic kickoff history was not fabricated, so cause-derived reviewer relaxation is not claimed.

## Recovery / next steps

Resume from this checklist and canonical FINAL/consensus, not the process transcript. Preserve runtime receipts and owner worktrees. No worktree declarations/pruning or credential edits. Read raw independent artifacts after process exit before accepting them. Provider429/503/timeout waits and retries follow source-context/ORGANIZER-BRIEF.md.

## Implementation details for review

The new surface is `budget legacy inspect|apply --run parley-deck/runs/NAME`; apply requires the exact `--expected-preview-sha256`, decision ID/reason, `--writers-stopped`, `--acknowledge-unknown-history` and `--yes`, plus the existing attended platform check. The shared store is `legacy-<scope-hash>/records/<decision-id-hash>.json` in the existing budget area. Frozen cycle adoption is `legacy.json` plus `legacy_history_sha256`, disclosed by cycle inspect. Exact event-data vocabulary is intentionally tighter: run.created descriptive mode/task strings; run.phase a known non-cycle action. Every other data field refuses. Single-link regular files also prevent hard-link aliases.

Goal execution uses existing RunParticipantStep for every checker. `participant_timeout_ns` in immutable invocation metadata retains the first resolved ceiling across run/configuration changes. Successful legacy output without this new field can replay; an unfinished old step lacking its original bound cannot launch a new child with an invented ceiling. No policy or participant set is rewritten. Protocol copies and staged core2.18.0 carry identical normative hunks; the staged core remains unpublished.

AC9 spans delivery after code review. The authorized release sequence necessarily publishes channels after reviewed source merges. Functional pre-merge goal verification must cover AC1–AC8 and AC9's build/drift portion; AC9's live channel/install evidence remains a separately binding post-publication verification before the final released handoff. No channel is described as delivered before verification, and pending attended activation/npm/core acts are disclosed. Reviewers may challenge this sequencing; it changes no product scope or final delivery criterion.

Review candidate published with full suites still running. Initial skill runs exposed missing local dev dependency commonmark and an unregenerated manifest; npm ci installed the existing locked dependencies, the manifest was regenerated, and a clean skill-suite run is now in progress. No test was suppressed. The Go version file was brought into lockstep with internal/app/version.go before review. Full test results remain a merge gate.
