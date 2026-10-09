---
idea: meta-protocol-change-driver-unstall
status: final
author: codex-1
implementer: codex-1
consensus-date: 2026-10-09
participants: [codex-1, zcode-1]
---

## Final plan / specification

Implement exactly D1 and D2 below, accepted by both current participants. The controlling brief authorizes implementation and release; it does not authorize an agent to impersonate an attended operator. Kimi was excluded under that brief after two actual HTTP403 readiness failures; no cause-derived automatic reviewer exception is fabricated. codex-1 organizes, drafts and implements; zcode-1 independently reviews and goal-checks in separate processes.

### D1 — Durable declaration of narrowly bounded unknown history

Add read-only `budget legacy inspect` and attended `budget legacy apply` surfaces (equivalent discoverable naming is an implementation detail). The operator records a one-time repository/deck-scoped declaration in the existing git-common-dir budget area; non-Git decks use their existing local runtime origin. The driver only reads this authority and never creates it. Existing per-idea migration and all worktree-coverage gates remain unchanged.

The new authority admits only an exact canonical relative run path under this deck's `runs/` whose complete recursive manifest is bound by digest, containing exactly one regular `events.jsonl` file and no other entry. The file must contain nonempty, structurally valid JSONL events with the closed vocabulary `run.created`/`run.phase`, no recoverable idea identity, and no known execution/charge records. No `run.json`, symlinks, aliases, special files, malformed records, hidden execution files or prose-based exceptions qualify. Reuse existing strict decoding, bounded reads, path/manifest and identity checks; a stricter refusal takes precedence over the May example. Missing events never prove zero historical execution.

Persist checked-durable, append-only records, one per explicit decision ID, including a reason, UTC time, writers-stopped assertion, attended-authority marker, explicit unknown-history acknowledgement, path and manifest digest. Inspect supplies an exact preview/hash. Apply retains the existing attended operator boundary, explicit confirmation, writers-stopped assertion and exact-preview verification. No terminal allocation or fake attendance by agents. Identical replay is idempotent; conflicting ID reuse refuses. Do not infer authority from earlier per-idea migrations.

Validate every visible copy on adoption and on every gate use. A changed/missing/all-deleted declared run, recovered identity, nonregular/aliased entry, corrupt declaration, conflicting or divergent copy, unknown additional run, or unavailable root still refuses. Freeze the adopted payload and hash with the first cycle binding and retain replay provenance. Inspection discloses `unknown-history`, decision ID and digest, never a zero count. Preserve all historical bytes, known charges and cycle caps; richer/charged unknown history retains the existing attended per-idea migration path. A binary without a ledger behaves as before; older binaries gain no automatic admission.

### D2 — Derived ceiling and bounded goal-check execution

The goal-check hard ceiling is `min(active-track timeout, positive configured checker timeout)`. Track values are fast 5 minutes, standard 15 minutes, deliberation 30 minutes. Missing configured timeout uses the track bound; absent track means standard; malformed/out-of-vocabulary track refuses rather than inheriting the permissive helper fallback. Resolve once per logical step and use the same bound on both attempts.

Use existing `RunParticipantStep` execution for goal checks, including protected checkers, with durable idea/agent/logical-step identity across restarts/runs, original plus at most one retry, existing cleanup, buffering and watchdog semantics. Valid FAIL is final, not a retry trigger. Invalid own output or classified terminal child/watchdog failure may consume the existing second slot; PASS from a failed process cannot close. Cancellation, control-plane refusal, tampering or unresolved writers still fail closed. Retrying a protected checker never authorizes exclusion. Do not broaden membership policy, protected-role rules or reviewer counts. LE-7/LE-11, dissent/reservations and independent current-tree acceptance evidence stay binding.

Synchronize only the necessary timeout/retry normative wording in live protocol, bootstrap copy, skill references and staged core. Target CLI 1.54.0 and skill/core 2.18.0, staging core from 2.17.0. Run the full authorized release including independent channel verification.


## Purpose / user-visible outcome

After one explicit attended declaration of eligible legacy history, successive new ideas can pass the first cycle-accounting gate without repeating an unstored per-idea declaration. Historical uncertainty remains visible and future mutation/unknown history still refuses. Independent goal checks receive a hard ceiling appropriate to their track and configured agent, with bounded retry and unchanged completion requirements. Until the owner attends activation, this repository's D6 gate remains blocked; release of the product fix is not local activation.

## Context & orientation

Base CLI 1.53.0 at 128e30b, skill 2.17.0 at 8ce4dec. Main locations: internal/budget/cycle_binding.go and cycle_history.go (cycle bootstrap), binding.go (shared repository/deck-prefix origin), run_identity_inventory.go and protocol_migration.go (manifest/eligibility and current request-scoped declarations); internal/app/driver_impl.go (GoalCheck), internal/runner/dropout.go and consult.go (durable attempts and supervised execution). The May source is parley-deck/runs/20260510T194003Z/events.jsonl, originally committed at 3ec10ac, with one identity-absent run.created event. Its bytes must remain intact.

The gate was reproduced by this idea's continue command and measured round02 launch before dispatch. Both source-backed cross-reviews identify the literal 120-second goal timeout. The historical >120-second execution observation motivates replacing the bound; no precise 518-second assertion is required. Zcode's round01 inference that the absence of cycle events proves zero work was challenged as WRONG/PRIMARY and withdrawn by its round02 SELF-CORRECTION. The existing unknown-history test and source explain why incomplete bytes cannot prove zero. No unresolved DISPUTED claim supports this plan.

Nominally independent proposals form one family: explicit declaration plus existing supervision/track policy. The two-model convergence is a shared prior, not independent evidence, and now shares the same source context. Acceptance depends on observable adversarial witnesses rather than unanimity as proof.

## Observable acceptance criteria

- AC1: without authority the actual May-shaped fixture refuses; after one attended declaration at least two different new ideas bind without repeated per-idea declarations, retaining the exact history bytes.
- AC2: declared history remains unknown, with decision/digest and frozen adoption provenance visible; identical replay succeeds and conflicting ID reuse refuses.
- AC3: mutations/deletion/all-missing copies, corruption, symlinks/aliases, extra unknown history, recovered identities, driver cursors, hidden execution files or charge records refuse. Revalidate across worktree copies.
- AC4: known charges/caps, unavailable-root gates and existing per-idea migration behavior remain unchanged; no driver-created authority.
- AC5: default/valid/invalid track and missing/shorter/longer checker timeout tests prove the exact derived bound, including a bound beyond the former 120 seconds.
- AC6: original plus one attempt at the same resolved ceiling survives restart/run changes without a third attempt; child/watchdog failures and malformed own output exercise the retry, valid FAIL does not.
- AC7: process-failed PASS, hard timeout, cancellation, tamper/control-plane refusal and protected-checker failure cannot close or create dropout authority; buffered transports retain hard bounds.
- AC8: unchanged independent current-tree AC evidence and LE-5/LE-7/LE-11 gates remain required, with current participants' ACCEPTs and a fresh independent goal PASS under the brief's authorized attended close if the mechanical exception remains unavailable.
- AC9: CLI/skill checks and protocol drift/packet tests pass; releases/install hashes and fresh independent channel evidence are recorded. D6 remains locally activation-pending until the owner attends the one-time apply; handoff contains the exact inspected request/command and all still-unpublished npm/core commands in version order.


## Idempotence & recovery

Read-only inspection never creates authority. Apply binds an exact preview, explicit decision ID/reason, stopped writers and owner attendance; the same request replays idempotently, conflicting reuse refuses. Preserve complete historical bytes and checked-durable provenance. A partial/corrupt declaration or policy publication must refuse until the exact original decision can be replayed or recovered; do not silently initialize over it. Frozen adoption and every visible run copy are revalidated on use. New undeclared history and missing roots cannot borrow an earlier exception. Existing per-idea migrations retain their semantics and no worktree declaration is performed in this idea.

Goal attempts retain existing durable logical step identity across restarts and run changes; do not create a third child from an interrupted/failed second attempt. Use the same resolved ceiling for the step's two attempts. Control-plane/tamper/cancellation and unresolved writer state retain their original fail-closed behavior. A valid semantic FAIL is a completed check that does not establish the goal, not a request for another answer. Current-tree independent AC evidence remains separate from a successful textual PASS.

## Known risks / de-risking

The accounting exception could silently erase history if broadened or trusted without current copy checks. The closed admission predicate, frozen digest/payload and AC1–AC4 target that risk. Runtime append durability, path aliases, conflicting copies and concurrent writers need negative tests. Retry logic could amplify attempts or conflate protected roles with dropout authority; AC5–AC8 target those risks. Filesystem coverage is bounded by the existing visible-root inventory and preserves unavailable-root refusals.

One-time activation, npm publication and core publication retain their shipped owner-only attendance boundaries. Prepare exact commands and truthful pending status after independent review. The controlling brief separately pre-authorizes the attended one-reviewer close only when the final independent review has no open CRITICAL/MAJOR, both current participants ACCEPT, a fresh independent goal check passes and current-tree AC evidence is recorded. Maximum five fix-up cycles.

## Alternatives disposition and deferred scope

- ALT-1: adopt existing per-idea migration's strict manifest/identity/unknown-history primitives; reject repeated request-scoped declarations as the lasting solution. Keep that command for richer/charged history.
- ALT-2: reject an idea-only skip of unscoped events; no identity exists to justify ignoring them. An exact attended structural declaration supplies explicit authority instead.
- ALT-3: reject rewriting/moving/deleting the May history or manufacturing its identity; the evidence does not establish one.
- ALT-4: reject another arbitrary goal constant; adopt existing track values clamped by positive checker configuration.
- ALT-5: keep skill-authorized recorded native launch fallback for this blocked run, but reject manual workarounds as the product fix.
- ALT-6: adopt existing RunParticipantStep/supervisor machinery; reject new watchdog/retry machinery and widening protected-role membership policy.
- ALT-7: reject unavailable-worktree declarations as part of this surface; existing gates are unaffected and the brief forbids making such declarations here.
- ALT-8: defer auxiliary run resolver/planner together: the latest consensus-signoff can shadow the driving run, but the read-only planner also lacks implementation actions. Choosing an older run alone risks stale roster/counters/pending state. Follow-up must specify both layers.
- ALT-9: defer ready/partial consensus reopen: archived bytes alone do not authorize dissolving a signed decision or hiding a BLOCK. Follow-up needs a checked completed-fix-cycle archive tied to the implementation and new review. Manual renames bypass that missing lifecycle verb.
- ALT-10: defer fixed-slug/existing-kickoff startup, missing-current-round launches, seeded snapshot bootstrap, placeholder-before-exit validation, generic focused briefs, wider status/wait reporting, unrelated Windows behavior and quota reporting. Each belongs in a separate scoped driver-gap follow-up; no implementation promise is added here.


## References

- Consensus: ./consensus.md (both current ACCEPT blocks).
- Rounds: ./round-01/, ./round-02/.
- Controlling source: ./source-context/ORGANIZER-BRIEF.md.
- Actual D6 refusal: ./source-context/continue-round02.log and ./source-context/round02-measured-budget-refusal.json.
- Process evidence: ./source-context/round02-native-result.json; independent artifacts retain protocol attestations.
- Design PR: https://github.com/feci/parley-deck-cli/pull/80.
