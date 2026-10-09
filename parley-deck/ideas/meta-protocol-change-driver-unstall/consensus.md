---
idea: meta-protocol-change-driver-unstall
drafted-by: codex-1
date: 2026-10-09
---

## Agreed decisions

Adopt D1 and D2 from both round-02 artifacts. codex-1 is organizer, participant, consensus/FINAL drafter and designated implementer; zcode-1 supplies independent review and fresh goal checking. Procedural convergence remains provisional until both current participants sign. Kimi's owner-authorized exclusion has no invented automatic-dropout history; the brief's attended-close fallback remains conditional.

### D1 — Durable declaration of narrowly bounded unknown history

Add read-only `budget legacy inspect` and attended `budget legacy apply` surfaces (equivalent discoverable naming is an implementation detail). The operator records a one-time repository/deck-scoped declaration in the existing git-common-dir budget area; non-Git decks use their existing local runtime origin. The driver only reads this authority and never creates it. Existing per-idea migration and all worktree-coverage gates remain unchanged.

The new authority admits only an exact canonical relative run path under this deck's `runs/` whose complete recursive manifest is bound by digest, containing exactly one regular `events.jsonl` file and no other entry. The file must contain nonempty, structurally valid JSONL events with the closed vocabulary `run.created`/`run.phase`, no recoverable idea identity, and no known execution/charge records. No `run.json`, symlinks, aliases, special files, malformed records, hidden execution files or prose-based exceptions qualify. Reuse existing strict decoding, bounded reads, path/manifest and identity checks; a stricter refusal takes precedence over the May example. Missing events never prove zero historical execution.

Persist checked-durable, append-only records, one per explicit decision ID, including a reason, UTC time, writers-stopped assertion, attended-authority marker, explicit unknown-history acknowledgement, path and manifest digest. Inspect supplies an exact preview/hash. Apply retains the existing attended operator boundary, explicit confirmation, writers-stopped assertion and exact-preview verification. No terminal allocation or fake attendance by agents. Identical replay is idempotent; conflicting ID reuse refuses. Do not infer authority from earlier per-idea migrations.

Validate every visible copy on adoption and on every gate use. A changed/missing/all-deleted declared run, recovered identity, nonregular/aliased entry, corrupt declaration, conflicting or divergent copy, unknown additional run, or unavailable root still refuses. Freeze the adopted payload and hash with the first cycle binding and retain replay provenance. Inspection discloses `unknown-history`, decision ID and digest, never a zero count. Preserve all historical bytes, known charges and cycle caps; richer/charged unknown history retains the existing attended per-idea migration path. A binary without a ledger behaves as before; older binaries gain no automatic admission.

### D2 — Derived ceiling and bounded goal-check execution

The goal-check hard ceiling is `min(active-track timeout, positive configured checker timeout)`. Track values are fast 5 minutes, standard 15 minutes, deliberation 30 minutes. Missing configured timeout uses the track bound; absent track means standard; malformed/out-of-vocabulary track refuses rather than inheriting the permissive helper fallback. Resolve once per logical step and use the same bound on both attempts.

Use existing `RunParticipantStep` execution for goal checks, including protected checkers, with durable idea/agent/logical-step identity across restarts/runs, original plus at most one retry, existing cleanup, buffering and watchdog semantics. Valid FAIL is final, not a retry trigger. Invalid own output or classified terminal child/watchdog failure may consume the existing second slot; PASS from a failed process cannot close. Cancellation, control-plane refusal, tampering or unresolved writers still fail closed. Retrying a protected checker never authorizes exclusion. Do not broaden membership policy, protected-role rules or reviewer counts. LE-7/LE-11, dissent/reservations and independent current-tree acceptance evidence stay binding.

Synchronize only the necessary timeout/retry normative wording in live protocol, bootstrap copy, skill references and staged core. Target CLI 1.54.0 and skill/core 2.18.0, staging core from 2.17.0. Run the full authorized release including independent channel verification.

### Required regression witnesses

- AC1: without authority the actual May-shaped fixture refuses; after one attended declaration at least two different new ideas bind without repeated per-idea declarations, retaining the exact history bytes.
- AC2: declared history remains unknown, with decision/digest and frozen adoption provenance visible; identical replay succeeds and conflicting ID reuse refuses.
- AC3: mutations/deletion/all-missing copies, corruption, symlinks/aliases, extra unknown history, recovered identities, driver cursors, hidden execution files or charge records refuse. Revalidate across worktree copies.
- AC4: known charges/caps, unavailable-root gates and existing per-idea migration behavior remain unchanged; no driver-created authority.
- AC5: default/valid/invalid track and missing/shorter/longer checker timeout tests prove the exact derived bound, including a bound beyond the former 120 seconds.
- AC6: original plus one attempt at the same resolved ceiling survives restart/run changes without a third attempt; child/watchdog failures and malformed own output exercise the retry, valid FAIL does not.
- AC7: process-failed PASS, hard timeout, cancellation, tamper/control-plane refusal and protected-checker failure cannot close or create dropout authority; buffered transports retain hard bounds.
- AC8: unchanged independent current-tree AC evidence and LE-5/LE-7/LE-11 gates remain required, with current participants' ACCEPTs and a fresh independent goal PASS under the brief's authorized attended close if the mechanical exception remains unavailable.
- AC9: CLI/skill checks and protocol drift/packet tests pass; releases/install hashes and fresh independent channel evidence are recorded. D6 remains locally activation-pending until the owner attends the one-time apply; handoff contains the exact inspected request/command and all still-unpublished npm/core commands in version order.

## Agreed trade-offs

The durable decision is machine-local shared-repository authority, not a tracked historical rewrite. Its exact digest and disclosure provide reviewability, at the cost of one attended activation. Unknown history remains unknown; this is an explicit scope decision, not proof of complete accounting. A conservative structural-only region is intentional; richer history is outside this new surface. Goal checks can take longer, bounded by existing track/configuration ceilings and a single retry.

## Open items deferred to implementation

No open design choice or reservation. Choose exact CLI flags/record layout using existing budget storage conventions and record them in IMPLEMENTATION. Prepare the owner's one-time activation and npm/core publication artifacts; do not claim those attended acts happened. Scope deferrals are listed below and must be carried into FINAL and the release handoff.

## Comparison & blind spots

Both independently authored proposals belong to the same family: explicit declaration of unknown history plus reuse of the existing timeout/supervisor machinery. Agreement by these two models is a shared prior, not independent evidence. Both have now read the same source and predecessor notes; no unavailable Kimi evidence is counted. No empirical proof of historical zero usage exists. Deployment on every possible filesystem and old binary interaction has not been exhaustively observed.

The design would be wrong if it silently counts unknown history as zero, admits charged/recovered/altered or undeclared history, resets a cycle cap, hides an unavailable root, retries a semantic FAIL, admits a failed PASS, amplifies attempts across restart, or turns retry eligibility into dropout/close authority. AC1–AC8 are concrete refutation witnesses for those risks. Phase 6 must try them, not merely inspect green happy paths.

The D6 gate was reproduced in this idea (`source-context/continue-round02.log` and measured budget refusal). The actual literal goal ceiling is source-backed in both rounds; historical checks exceeding 120 seconds motivate derivation. The exact historical 518-second figure is not relied upon. The zero-charge inference in zcode-1 round-01 was withdrawn in its round-02 SELF-CORRECTION after codex-1's V-ZERO WRONG/PRIMARY challenge. The existing test explicitly says "no evidence row and no count — in particular no floor of zero" and asserts UnknownHistory; absence of recorded events does not entail complete history. No conflicting verdict or dependent DISPUTED claim remains at consensus opening.

## Drafter position changes

The prior D1 sentence in `round-02/codex-1.md` was: "The narrow record may admit only structurally valid identity-absent run histories with no driver cursor and no known execution/charge records." This draft makes that admission closed and tighter: exactly one regular events.jsonl, only run.created/run.phase events, no other entries. This implements zcode-1 round-02's required enumeration and makes hidden execution metadata ineligible. No other material position changed since my round-02. D2, two-fix scope and attended activation are unchanged.

## Alternatives disposition

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

## Signoffs

<!-- Each participant appends only its own signoff. -->

### Signoff: codex-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: Drafter: yes. I accept D1/D2, the closed structural predicate and AC1–AC9, the recorded correction and deferrals. No material disagreement remains. Attended activation and existing close safeguards remain binding.

### Signoff: zcode-1 — 2026-10-09
Status: ✅ ACCEPT
Notes: I accept D1/D2 as drafted. The consensus folds in every enumeration I pinned in round-02 verbatim or tighter: the closed structural admission predicate (exactly one regular events.jsonl, run.created/run.phase vocabulary only, no other entries), the replay/conflict refusal list, and the D2 absent-track→standard / malformed-track→fail-closed rule that avoids the permissive unknown-track helper branch. The drafter's D1 narrowing is the closed enumeration my round-02 required, not a new position. My V-ZERO self-correction stands: declared history is unknown-history, never a zero count. Two-participant convergence is correlated agreement, disclosed as such; AC1–AC9 are the falsifier set, and I will hold Phase 6 review to them as refutation witnesses, not happy paths. Attended one-time activation, npm/core publication and the conditional attended close remain owner acts outside this signoff.
