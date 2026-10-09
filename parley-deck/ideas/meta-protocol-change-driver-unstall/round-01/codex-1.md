---
agent: codex-1
idea: meta-protocol-change-driver-unstall
round: 1
date: 2026-10-09
---

## Summary

Keep the change centered on auditable treatment of legacy unscoped history and a realistic independent goal-check ceiling. I favor a durable, exact-manifest declaration over pretending that the May record has a known identity or performed zero work. Reuse the existing track timeout and participant-step supervision for the goal check. Phase selection after auxiliary signoff runs is a plausible small third fix; reopening accepted consensus has more lifecycle implications and should remain deferred unless the peer supplies a narrow safe contract.

I wrote this independent analysis before reading any other participant's round-01 file. codex-1 is the organizer and implementer, as the owner directed; independent verification belongs to another participant.

## Proposed approach

### D1 — Durable legacy declaration, with explicit uncertainty

The bootstrap scanner currently calls `cycleRunEvents` on every visible run before it can discriminate this idea's history (`internal/budget/cycle_history.go`, `refuseUnmigratedCycles`). The actual legacy `parley-deck/runs/20260510T194003Z/events.jsonl` contains exactly:

```json
{"time":"2026-05-10T19:40:03.126637Z","type":"run.created","data":{"mode":"auto","task":"smoke implementation run"}}
```

No run.json exists in that directory at HEAD `128e30b`. These are my source observations for the independent verifier, not self-issued verification verdicts. The existing `cycleRunEvents` refuses with `historical cycle event lacks idea identity` on a run.created or run.phase without an idea.

Proposed direction: add a small repository-scoped recorded legacy declaration, binding the exact relative run path and existing `RunDirectoryManifestDigest` of its complete file set, with decision ID, reason and explicit owner authority. Retain the original run bytes; state `unknown-history`, never zero. At first cycle binding, read and validate that declaration, validate every visible copy using existing declaration-eligibility rules, disclose its use, and freeze its hash/provenance in the cycle binding so a later edit cannot quietly change the decision. Missing, malformed, conflicting, changed, unreadable, symlinked, non-eligible or newly unscoped history still blocks. Do not change worktree coverage rules or import unrelated charged history.

The intended witness is a pair of new ideas bootstrapping through the same exact declared May file set while an added unknown run or changed legacy byte causes both to refuse. This is an explicit recorded scope decision, not proof that the unknown history was empty. A declaration that is merely an unstored flag fails the lasting-fix criterion. A repository-specific exception hardcoded into the binary is also undesirable. The exact record/command and authority boundary must be settled before FINAL; do not invent owner reconciliation or broaden the exception to identity-bearing runs.

### D2 — Track-derived goal ceiling; preserve the close predicate

`internal/app/driver_impl.go`, `driverImplOps.GoalCheck`, presently passes `Timeout: 2 * time.Minute` to `runner.RunConsult`. It already has a `RunParticipantStep` path for saved enabled mid-idea dropout policy and non-protected checker roles. `runner.RunConsult` uses `supervisionForStep`; the buffered Zcode transport can legitimately produce no output until completion. The protocol explicitly repeats the two-minute goal ceiling in Phase 8 and section 9.0.

Use the existing per-track timeout value (5/15/30 minutes, `internal/app/wait.go`, `trackTimeoutCeiling`) with any shorter configured agent hard bound honored. Specify absent/unknown track behavior explicitly, preserving fail-closed invalid-track handling. Feed that same resolved ceiling through both attempts of the existing original-plus-one step mechanism, keeping first/stall/heartbeat behavior and truthful buffering declarations. No second retry budget, no retry on a valid FAIL or reserved/inconclusive semantic result merely to obtain PASS, no new same-agent exception, and no implication that the textual PASS establishes completion. Existing role protections, process-success requirement, current-tree AC evidence, dissent, reservations and reviewer gates remain.

Tests should use short controlled subprocesses/seams to demonstrate a check beyond the old bound can pass under a larger allowed ceiling, that timeout/retry remains bounded, a valid FAIL is not retried, buffered output retains its hard deadline, and failed/self/missing/ambiguous checks never pass. Synchronize the normative sentence in the live authority, bootstrap authority, skill reference and staged core; update skill guidance without unrelated rule edits.

### D3 — Other driver gaps

`internal/runstate/runstate.go:ResolveRun` returns the first idea-matching run from `ListRuns`; auxiliary consensus-signoff runs are newer and can shadow the actual driver run. A narrow continuation resolver should prefer a valid driving run for an idea target while keeping an explicit run ID exact, preserving safety counters and refusing ambiguity/corruption. It must not synthesize a cursor, drop pending state, or reinterpret completion. Include this only if the independent analysis finds it small and establishes the exact acceptable run classifications.

`internal/consensus/consensus.go:Reopen` currently accepts only `TriageBlocked`. A blanket ready/partial relaxation could erase an accepted decision or evade a signer obligation. Prefer deferral unless a checked completed-fix-cycle archive operation can be specified and tested without adding a second workflow.

Other recorded gaps to defer: fixed-slug/existing-kickoff startup, launch of missing current-round artifacts, frozen snapshot bootstrap for seeded runs, placeholder acceptance before child exit, generic focused-brief command support, broader status/wait phase reporting, unrelated Windows behavior and quota reporting.

## Existing alternatives

- **ALT-1, current attended per-idea migration**: `internal/budget/protocol_migration.go`, `ProtocolMigrationRequest.DeclaredUnscopedRuns`, and `internal/budget/protocol_migration_unscoped_test.go`. It already binds all visible copies, retains readable sources and records unknown history rather than zero. Adopt its manifest/eligibility primitives; reject repeating an unstored declaration for every idea as the lasting UX.
- **ALT-2, idea-only scan**: `refuseUnmigratedCycles` already skips positively identified other ideas. Skipping genuinely unidentified events would lose the very evidence this gate protects. Reject a bare skip; require a demonstrable provenance rule if a narrower safe variant is proposed.
- **ALT-3, legacy record migration**: adding an idea to the May event or moving/deleting it would be a historical alteration. The available line and absence of a manifest do not prove an identity. Reject absent independent primary evidence locating the real identity.
- **ALT-4, arbitrary larger goal constant**: raising 120 seconds to another number is small but repeats policy outside the existing track abstraction. Prefer track/config-derived bounds and the already shipped supervisor/retry path.
- **ALT-5, manual configured CLI workarounds**: predecessor organizer notes record them and the skill permits them for concrete driver gaps. They keep this idea moving but do not fix product behavior.

## Concerns / open questions

The durable declaration must have a precise authority and replay contract before implementation; an explanatory markdown sentence alone must not disable a safety gate. Determine whether a tracked reviewed record under this brief's authority is sufficient or whether the existing attended migration mechanism must persist a new reusable scope decision. If a new human-only reconciliation is required, that is a single owner decision with concrete options, not an invented zero count.

The current named-idea launch cannot mechanically create kickoff failure history: `parley continue --dir . --json meta-protocol-change-driver-unstall` returned `has no runs yet`; `parley run` always mints a timestamped slug. Dogfood the real 1.53.0 reviewer proof, but never fabricate automatic history from the two manual readiness invocations. The brief explicitly provides an attended close if that exception does not qualify.

## Risks

The strongest risk is making an accounting exception too broad. Adversarial tests must exercise exact manifest changes, unknown additional runs and existing charged history. Goal-check scheduling risks retry amplification or treating failure as exclusion/close authority. Auxiliary-run selection risks choosing a historical driver run without preserving current state. Keep the scope small enough that each risk has a specific regression witness.

## Protocol context

context_mode: full. source_sha256 = packet_sha256 = `acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`. fallback_reason: absent. Full live source read, including section 15. Material observations above are owned claims offered for independent verification; no self-verdict is issued.
