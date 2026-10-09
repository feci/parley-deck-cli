---
agent: zcode-1
idea: meta-protocol-change-driver-unstall
round: 1
date: 2026-10-09
---

## Launch attestation

Protocol context rendered with `parley protocol packet --dir . --phase 1 --track
deliberation --idea meta-protocol-change-driver-unstall --flag auto_implement --flag
protocol_change --json`: **context_mode=full**, **source_sha256 = packet_sha256 =
`acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`**, **fallback_reason
absent**. I re-hashed the packet body file myself (`shasum -a 256` → same digest) and it
matches the launch-supplied protocol verbatim, so the full body was read twice over. I read
`00-prompt.md` and `source-context/ORGANIZER-BRIEF.md` in full; the brief's SHA256
re-verified as `61bbd09b54a5afa753e4dcac0e03a5232e24cc6248c4bc3ef18ae655b19d5d1f`,
matching the kickoff. I read all three predecessor organizer notes in full and did not open
`round-01/codex-1.md` before finishing this file (it exists; independence rule respected).

Scratch evidence for this round (command transcripts, packet JSON) is retained under
`.parley-runtime/driver-unstall/zcode-round01/`. CLI worktree HEAD at analysis time:
`ea96877` (kickoff commit, past the briefed base `128e30b`); skill worktree at `8ce4dec`
as briefed.

## Summary

Both target stalls are confirmed in current source, not just in testimony. D6 is a
hard-coded refusal in the budget bootstrap: `cycleRunEvents`
(`internal/budget/cycle_history.go:251`) rejects any `run.created`/`run.phase` event
lacking `data.idea`, and the tracked May smoke run
(`parley-deck/runs/20260510T194003Z/events.jsonl`, one line, committed in `3ec10ac`) carries
exactly such an event, so `refuseUnmigratedCycles` blocks the first cross-review binding of
every new idea in every checkout. The shipped remedy (`budget migrate apply
--declare-unscoped-run`) is real but per-idea and request-scoped by design — the flag help
says "never stored as repository state" — so nothing durable ever resolves the recursion.
The 120 s goal-check ceiling is a literal `Timeout: 2 * time.Minute` at
`internal/app/driver_impl.go:838`, while the last real independent goal check ran ≈8 minutes
in a manual 1800 s process precisely because the product ceiling cannot fit it
(`meta-protocol-change-participant-dropout-review-gate-timing/goal-check-zcode-1.md`,
"Elapsed at artifact authorship ≈ 8 min (>120 s)"). I propose a durable, digest-bound,
deck-shared unscoped-run declaration consulted by the gate (never inventing identity or
zeroing provable charges), a goal-check ceiling derived from the checker's configured
per-agent timeout routed through the shipped two-attempt participant-step path with two
minimal protocol hunks, and the two small CLI-only gaps (phase-pointer resolution,
ready-consensus reopen) — deferring the rest.

## Proposed approach

### A. D6 — durable declaration for unscoped legacy runs

**What ships today (all PRIMARY, consulted directly):**

- `EnsureCycleBinding` calls `refuseUnmigratedCycles` when an idea has no cycle policy yet
  (`internal/budget/cycle_binding.go:161` and again at `:187` before first policy write).
- `refuseUnmigratedCycles` (`internal/budget/cycle_history.go:112-208`) enumerates
  `parley-deck/runs/*/events.jsonl` under **every** registered worktree root
  (`launchScopeDeclared`, `internal/budget/binding.go:127-177`, builds `roots` from
  `git worktree list --porcelain`).
- `cycleRunEvents` (`cycle_history.go:222-268`) hard-fails
  `"historical cycle event lacks idea identity"` at `:251` on any `run.created`/`run.phase`
  event without `data.idea` — the identity comparison that would skip another idea's runs
  happens only *after* this parse, so an unscoped run can never be excluded as foreign; it
  is a structural refusal.
- The only escape is `parley budget migrate apply --kind cross-review --declare-unscoped-run
  parley-deck/runs/<name>=<64-hex manifest digest>` (`internal/app/budget_migrate.go:63`).
  The declaration is retained inside that **idea's** immutable `migration.json`
  (`internal/budget/protocol_migration.go:26-38`, `:69-92`), and "no ordinary bootstrap ever
  consults either" declaration kind — so each new idea re-hits the gate. This matches all
  three predecessor notes: gap 11 first recorded 2026-10-03
  (`meta-protocol-change-quota-auto-exclude/organizer-notes.md`), then the identical refusal
  in `meta-protocol-change-participant-dropout/organizer-notes.md` ("D6 fallback —
  cross-review") and `meta-protocol-change-participant-dropout-review-gate-timing/organizer-notes.md`
  ("Round 1 complete / D6 fallback"). I did not re-run the driver to reproduce the halt in
  this process; the runtime claim rests on those three recorded runs plus the static reading
  above.

**Design — generalize the existing declaration, do not weaken the gate.**

1. New durable record: a deck-shared, append-only **unscoped-run declaration ledger**
   (e.g. one JSON record per declaration under the budget area already shared across
   worktrees — budget state lives under the git common dir, `binding.go:144,159` — or a
   single `parley-deck/meta/budget-unscoped-declarations.json`; FINAL should pick one,
   and I lean toward the git-common-dir budget area because it is machine-local by design
   and reaches every worktree root without a checkout bump, while a tracked copy would lag
   in roots pinned to older commits).
2. Record shape mirrors `ProtocolMigrationRequest`: decision id, reason, recorded-at,
   writers-stopped, run relative path (`parley-deck/runs/20260510T194003Z`), and the
   recursive manifest digest of the run directory — the same digest discipline
   `checkDeclaredUnscopedRuns` already enforces (`protocol_migration.go:69-92`), including
   per-root copy verification.
3. Gate consultation: in `refuseUnmigratedCycles`, when `cycleRunEvents` encounters an
   unscoped event, load the declaration ledger; a declared run whose current bytes re-hash
   to the recorded digest is treated as *unscoped-but-accounted*: excluded from
   identity-scoped counting with the recorded reason, never assigned an idea.
4. **Integrity rules (fail-closed preserved):**
   - Undeclared unscoped run → refuse, exactly as today.
   - Declared run whose bytes changed in any root (digest mismatch) or vanished from every
     root → refuse (reuses the `InspectProtocolMigrationDeclarations` re-verification
     pattern, `protocol_migration.go:452-463`).
   - **Never a zero count by fiat:** the declaration excludes the run from *identity*
     accounting only. If the run's own bytes prove charged cycles (`run.phase`
     `promoted`/`reopened` events), a plain declaration must still refuse — that case
     requires the per-idea migration with an explicit carried count, as today. The May
     smoke run is declarable precisely because its single event line is a `run.created`
     with zero cycle events: its zero charge is proven by its own retained bytes, not
     asserted by the operator. This settles the authority/replay contract: authority is the
   attended operator decision (same attendance gate as `migrate apply`,
   `budget_migrate.go:93-95`); replay is the digest-bound exact record, re-verified on
   every gate pass; no historical identity is ever invented and no provable charge is
   ever zeroed.
5. Bootstrap compatibility: an older binary ignores the ledger (it never reads it) and a
   newer binary reading a repo whose ledger is absent behaves exactly as today — the
   ledger only ever *narrows* refusals for durably declared, byte-identical runs.

**Why not the alternatives:** an idea-scoped gate (skip unscoped history that "isn't
mine") is precisely the silent ignoring of real unscoped history the constraints forbid,
and is indistinguishable from forgetting. Rewriting the May event to add `data.idea`, or
tombstoning it, edits append-only evidence and invents an identity — forbidden. The
per-idea migrate already exists and demonstrably does not stop the recursion.

### B. Goal-check ceiling — derive from configured per-agent timeout

**Confirmed (PRIMARY):** `Timeout: 2 * time.Minute` is literal at
`internal/app/driver_impl.go:838` inside `GoalCheck`'s consult closure. The mid-idea
dropout path wraps that closure in `RunParticipantStep` (`:854-890`) for the two-attempt
discipline, but the inner ceiling stays 120 s; every other path runs the single bare
consult (`:891-893`). The protocol pins the same number twice: §4 Phase 8 "its existing
two-minute product ceiling stays" (`parley-deck/COOPERATION.md:721-722`) and §9.0
"goal120s … stay" (`:920`), mirrored in the skill copy
(`worktrees/driver-unstall-skill/skills/parley-deck/references/COOPERATION.md:721-722,920`)
and the embedded default (`internal/protocol/defaults/COOPERATION.md`, same zone).

Real duration evidence: the review-gate-timing goal check states "Elapsed at artifact
authorship ≈ 8 min (>120 s); final duration is recorded by the process's own exit receipt
for the disclosed timeout-tuning follow-up" (`…review-gate-timing/goal-check-zcode-1.md`,
Launch attestation). The brief's specific "518 s" figure is organizer testimony I could
not independently confirm from artifacts in this checkout; the >120 s fact does not rest
on it. Corroborating pattern from the same notes: real zcode child processes routinely run
270–1520 s (e.g. round-01 775.7 s, review-02 1520.0 s), and zcode's first output can land
only near exit (775.662 s of 775.745 s) — silent-but-successful under a final-text-only
transport whose soft guards are disabled, leaving the hard ceiling as the only bound.

**Design:**

1. Replace the literal with the checker's **configured per-agent `timeout_ms`** from the
   resolved discovery/roster spec (the same source `agents exec` uses; zcode-1 and codex-1
   are configured 1 800 000 ms). This automatically honors §4.0's per-track timeout norms
   because roster timeouts are configured to them; an operator override is respected
   without a code change. A defensive track clamp may cap it (deliberation ≈30 min) so a
   misconfigured 6 h roster entry cannot pin a driver loop — FINAL should decide whether
   the clamp is worth the extra rule; I lean yes, one line, bounded by the §4.0 table.
2. Route **all** goal checks through `RunParticipantStep` (two attempts total, supervisor
   guards where the transport streams), not only the mid-idea dropout branch. This is the
   brief's "reuse shipped retry/watchdog machinery": first-output/stall/heartbeat defaults
   already exist (`internal/runner/supervision.go:33`,
   `defaultFirstEventTimeoutMS = 120_000`); final-text-only transports keep soft guards
   off and take the hard ceiling only — unchanged semantics, honest ceiling.
3. LE-7/LE-11 untouched in substance: fresh process, non-implementer checker
   (`driver_impl.go:810-817`), fail-closed on missing/failed/inconclusive results and
   PASS-with-reservations (`:901-917`), verdict parsing with sticky FAIL (`:920-949`) —
   all stay. A timeout still yields "goal-check checker failed; completion is unverified"
   and escalates; the check can still only withhold a close, never establish one. Longer
   ceiling ⇒ more wall-clock under LE-5 budgets, which still escalate rather than
   auto-complete — that is the correct trade.
4. **Minimal protocol hunks** (all three copies, identical normative replacement):
   - §4 Phase 8, replace "its existing two-minute product ceiling stays" with: "its
     ceiling is the checker's configured per-agent timeout under the active track's
     timeout norm, routed through the shipped two-attempt participant-step supervision; a
     timed-out or stalled check leaves completion unverified."
   - §9.0, replace "Readiness90s, goal120s and existing track/operation/configured
     ceilings stay" with: "Readiness90s and existing track/operation/configured ceilings
     stay; the goal check uses the checker's configured per-agent timeout through the
     shipped two-attempt participant-step path". Exact wording is FINAL's to fix; the
     scope is these two sentences only.

### C. Phase-pointer resolution (include, CLI-only)

`ResolveRun` resolves an idea slug to the **newest** error-free run
(`internal/runstate/runstate.go:353-357`), which after any signoff is an advisory
`consensus-signoff` run rather than the driving run; the plan for that run finds nothing
recoverable and the CLI answers "No recoverable action; inspect artifacts and logs"
(`internal/runplan/runplan.go:213`) — exactly what both predecessor notes record hitting
during implementation/review. Narrow fix: when the newest matching run yields no
recoverable action, fall back to the newest **driving** run for the slug (or prefer
driving runs at resolution time). No protocol text. This gap is demonstrated twice
(PRIMARY: both organizer notes, plus the locators above), small, and it stalls the driver
in exactly the phases this idea's owner wants unstalled.

### D. Ready-consensus reopen (include, CLI-only)

`Reopen` refuses any triage except `TriageBlocked`
(`internal/consensus/consensus.go:491-493`). In the participant-dropout run a signed,
ready review consensus could not be transitioned by any verb, forcing a byte-identical
manual rename workaround (`…participant-dropout/organizer-notes.md`, "Independent
review-02 and close-check preparation"); review-gate-timing hit the same wall and kept a
`consensus-cycle-01-proposed.md` outside the canonical name. Narrow fix: let
`consensus reopen --review` accept `ready` (and `partial`) with a mandatory recorded
reason, preserving the existing mechanics — rename the file byte-identically to the
aborted path, append the reopen reason, keep all signoffs. Content is never rewritten, so
the append-only audit trail survives; the refused-today behavior for other triages is
unchanged. Alternatively a dedicated `consensus adopt` verb for proposed files; I lean
toward widening `reopen` because the abort-and-reason machinery already exists and one
verb with one condition changed is the smaller diff. Both organizers' workarounds are
recorded evidence this stalls closes.

### Deferred (explicitly out of scope for this idea)

From the predecessor notes' gap lists: no existing-kickoff entry point for `parley run`
(gap 1) and the driver launching only the *next* round (gap 2) — real but architectural,
each a behavior change in the runner core; no roster freeze for seeded runs (gap 3);
single cross-review round default (gap 4, deliberate); stub placeholders passing
`wait`/`roundComplete` (gap 5 — validator policy, needs its own idea);
`steer` non-delivery (gap 9); thin drafting prompts (gap 10 — being improved piecemeal);
scaffold-signoff request after failed draft (gap 12); headless-organizer waiting
conventions (gap 8, procedural, already adopted in practice). None blocks every idea the
way D6 does, and "keep it small" binds.

## Concerns / open questions

1. **Declaration ledger placement** (git-common-dir budget area vs tracked
   `parley-deck/meta/`): machine-local shared state vs auditable-and-cloned. If the owner
   wants the declaration visible in repo history, the tracked location wins despite the
   stale-checkout lag; I lean machine-local because the gate already trusts per-root byte
   verification, not checkout state.
2. **Track clamp on the derived goal ceiling**: one extra rule vs trusting roster
   configuration. I propose the clamp; codex-1 may prefer pure config trust.
3. Does widening `reopen` need a protocol-text mention (the protocol's review-consensus
   flow implies reopen-for-blocked only), or is it safely below the text as a CLI
   affordance? My reading: §4 Phase 7/8 does not constrain the verb's triage domain, so
   CLI-only is defensible; flag for consensus.
4. The brief's 518 s figure — if FINAL cites it, it should be marked testimony or
   re-derived from the retained exit receipts, not presented as independently confirmed.
5. Goal-check duration evidence is zcode-heavy; codex (streaming, faster first output)
   may never have needed >120 s. The derivation should still be uniform — the ceiling is
   per-checker config, not per-adapter special-casing.

## Risks

1. **D6**: any bug in digest re-verification could silently exempt a *changed* legacy run
   — mitigated by reusing the shipped `checkDeclaredUnscopedRuns` copy-verification shape
   and by keeping the nonzero-charge refusal path untouched. A declaration is itself a
   security-relevant act; it must stay behind the same attended gate as `migrate apply`,
   and the driver (LE-2) must never write it.
2. **Goal ceiling**: longer ceilings increase exposure to hung checkers on
   final-text-only transports (no stall detection) — bounded by two attempts, the hard
   ceiling, and LE-5 wall-clock escalation; the failure mode remains a refused close, never
   a false one. Two 1800 s attempts also exceed some driver step budgets; hitting a loop
   budget must escalate (it does today) rather than retry silently.
3. **Reopen widening**: accepting `ready` could let a participant dissolve a signed
   consensus too casually; the mandatory reason plus byte-identical preservation keeps the
   audit honest, and signoff gates still bind the replacement cycle.
4. **Protocol-text risk**: two-sentence hunks in three synchronized copies invite drift;
   the embedded-template drift test (`TestEmbeddedDefaultMatchesLiveDeck`) and packet
   `check` cover it, and the release plan already stages core 2.18.0.
5. All runtime-behavior claims here rest on static reading plus the three predecessor
   notes' recorded runs; I did not execute the driver in this round. Implementation must
   add reproductions (a fixture run with an unscoped event; a goal check exceeding 120 s)
   before close — those become current-tree AC evidence, which this round cannot supply.

## Existing alternatives

Per §15.6(a), what the toolchain already ships versus what this proposal builds by hand:

- **Per-idea `budget migrate apply --declare-unscoped-run`** (ships:
  `internal/app/budget_migrate.go:63`, `internal/budget/protocol_migration.go`) — the
  durable ledger generalizes its digest-bound declaration to deck scope with gate
  consultation; today's verb remains the path for genuinely charged legacy counts.
- **`--declare-unavailable-worktree`** (ships: `binding.go:68-112`) — the analogous
  request-scoped declaration for unreadable worktrees; confirms the declaration grammar
  this proposal extends rather than invents.
- **Supervisor defaults + `RunParticipantStep`** (ships:
  `internal/runner/supervision.go:33`, used at `driver_impl.go:855`) — the goal-check fix
  only routes an existing path uniformly and swaps the ceiling source; no new watchdog
  machinery is built.
- **`agents exec` manual launches + `parley wait`** (ships) — the status-quo alternative
  every predecessor organizer fell back on; this idea exists to stop needing it.
- **`consensus reopen`** (ships, blocked-only: `consensus.go:491-493`) — widened by one
  condition rather than replaced.
- **§4.0 per-track timeout table** (protocol norm) — the derivation target for the
  ceiling; no new timeout regime is created.

No participant proposal I have read (I have read none — round-1 independence) was
consulted; this file was written from the kickoff materials, the three predecessor
organizer notes, and the current source only.
