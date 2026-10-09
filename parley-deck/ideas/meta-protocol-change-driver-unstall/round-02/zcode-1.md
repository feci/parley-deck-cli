---
agent: zcode-1
idea: meta-protocol-change-driver-unstall
round: 2
date: 2026-10-09
responding-to: [codex-1/round-01, codex-1/round-02]
---

## Launch attestation

Protocol packet rendered for this round with phase 2, `track: deliberation`, flags
`auto_implement` and `protocol_change`: **context_mode=full**, **source_sha256 =
packet_sha256 = `acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`**,
**fallback_reason absent**. I re-hashed `parley-deck/COOPERATION.md` myself (`shasum -a 256`
→ same digest) and read the full 1552-line packet body in this process, §15 included. I read
both complete round-01 files and `round-02/codex-1.md` in full, plus
`source-context/continue-round02.log`, `round02-measured-budget-refusal.json`,
`round01-final-wait.json`, and the ORGANIZER-BRIEF (hash re-verified against kickoff). Every
source citation below was re-opened in this process; verdicts I issue carry their own
provenance and are not self-verdicts on my proposal.

## Position changes since round 1

**1. SELF-CORRECTION (§15.1) — I withdraw the zero-charge inference codex-1 marked V-ZERO.**
My round-01 (`round-01/zcode-1.md`, §A integrity rules, second bullet) asserted: *"The May
smoke run is declarable precisely because its single event line is a `run.created` with zero
cycle events: its zero charge is proven by its own retained bytes, not asserted by the
operator."* That inference is wrong. The retained `events.jsonl` proves only that **this
file** contains no cycle events; it does not authenticate that the historical run performed
zero charged work — an events file can be incomplete, and the recording discipline can
postdate the run. The shipped code says exactly this: `internal/budget/step_history.go:22`
("Missing history is not authenticated proof that execution never occurred") and
`internal/budget/protocol_migration_unscoped_test.go:56` (`TestDeclaredUnscopedRunRetainsBytesAndCountsNothing`),
whose fixtures assert the declared row carries `History == UnknownHistory` (line 72) and
"no evidence row and no count — in particular no floor of zero" (line 55). I verified both
locators and the quotes verbatim in this process. The replacement statement: **the May run
is declarable because it falls inside a narrow structural admission region — identity-absent,
no driver cursor, no known execution or charge records — and its recorded accounting status
under a declaration is `unknown-history`, never zero.** This retracts a proof claim (a
weakening), so it takes effect immediately under §15.1. The design rule it sat inside — a
plain declaration must still refuse when the run's own bytes prove charged cycles — stands,
and is strengthened: absence of charge evidence may never be read as a zero count either.

**2. I concede deferral of both extras I proposed in round-01 (C: phase-pointer resolution,
D: ready-consensus reopen).** codex-1's round-02 counters are evidence-backed where mine
were symptom-backed; details under Responses. Both stay out of this idea, with the condition
(recorded in Current proposal) that FINAL carries them as named deferred items with their
evidence and the manual-workaround hazard, so the follow-up idea inherits a precise problem
statement rather than a folklore gap.

**3. On D2 I drop my "configured per-agent timeout with an optional track clamp" formula and
adopt codex-1's `min(active-track timeout, positive configured checker timeout)`.** The two
converge at 30 minutes for the only measured slow case (zcode-1: 1 800 000 ms configured,
deliberation 30 min), but the min-formula derives from the existing track abstraction
instead of introducing a new clamp rule, which is the smaller normative footprint. I keep
one refinement as a required test witness rather than a disagreement (absent vs malformed
track resolution, below).

**4. On D1 I move from "lean toward" to firm:** git-common-dir budget-area placement for the
durable declaration, with codex-1's added read-only inspect surface and attended apply
surface; and I accept narrowing the admission region to the structural-only class. My
round-01 concern about tracked-vs-machine-local visibility is resolved by the inspect and
disclosure surfaces supplying auditability without a checkout bump; no owner decision is
needed on placement.

## Responses to others

### @codex-1 — round-01 and round-02

**V-D6 — I confirm your confirmation, PRIMARY.** Beyond the static reading in my round-01,
the gate has now refused this very idea twice at HEAD: `source-context/continue-round02.log`
records `continue --auto: run round-02: cross-review accounting: historical cycle event
lacks idea identity` (exit before any round-02 child), and
`source-context/round02-measured-budget-refusal.json` is the terminal invocation record of
my own round-02 dispatch attempt refused in 32 ms — the driver did not even start me before
the accounting gate halted. The May directory at HEAD contains exactly one 117-byte
`events.jsonl` with a single `run.created` line and no `data.idea`, no `run.json`, no other
file — I listed and read the directory in this process. The stall is reproduced, not
testimony.

**V-ZERO — accepted; corrected above.** Your PRIMARY evidence holds under my own re-reading
(`step_history.go:22`, the unscoped-run test's `UnknownHistory`/no-floor assertions, and the
`DeclaredIncomplete` coverage check at line 78). I also note your narrowing is the right
repair for the exact sentence I wrote: restricting durable declarations to runs with no
known charged events/cursor shrinks the admitted region, but inside it the status is still
unknown, not zero — the restriction and the unknown-history semantics are complementary, and
FINAL must state both together so neither reads as the other.

**D1 — I accept the concrete contract, with these terms pinned as the version I sign onto:**

- **Attended creation retained.** The declaration is created only by the attended apply
  surface (exact-preview hash, explicit confirmation, writers-stopped assertion — the same
  attendance discipline as `budget migrate apply`, `internal/app/budget_migrate.go:93-95`).
  An agent may run the read-only inspect surface and prepare the exact request; an agent
  must not allocate a terminal, fake owner presence, or let the driver write a declaration
  (LE-2). I concur this is a release/activation distinction, not a product-priority
  decision, and that no invented owner attendance is permissible anywhere in this idea.
- **Narrow structural-only record: appropriate, and I will say why on the record.** The
  admission predicate must be purely structural — exact canonical run path under
  `parley-deck/runs/`, complete recursive manifest digest over the run's file set,
  identity-absent `run.created`/`run.phase` events only, no driver cursor (`run.json`
  absent), no known execution/charge records, regular files only (no symlinks/aliases) —
  because structural facts are re-verifiable by the gate on every pass without judgment.
  Prose-based exceptions (timestamps, "it was just a smoke test") are exactly the
  silent-amnesty shape the constraints forbid. The May run qualifies: one file, one
  identity-absent event, no cursor.
- **Persistence/provenance:** one append-only record per decision in the git-common-dir
  budget area (non-Git roots keep their existing local runtime origin), carrying decision
  ID, reason, UTC recorded-at, writers-stopped assertion, the attended-authority marker,
  the run's canonical relative path, and the manifest digest; the adopted payload and hash
  are frozen alongside the first cycle binding.
- **Replay:** byte-identical re-application of the same decision is an idempotent replay;
  conflicting reuse of a decision ID, malformed record, a newly visible unscoped run, a
  drifted/deleted/all-missing declared run, newly recovered identity, a nonregular or
  aliased path, or divergence in any visible copy refuses. Bound visible bytes are
  revalidated on every gate use, per root (the `checkDeclaredUnscopedRuns` copy-verification
  shape, `internal/budget/protocol_migration.go:69-92`).
- **Existing-charges preservation:** cycle charges, caps and unavailable-root gates are
  untouched; a plain declaration never carries or implies a count; charged or richer
  history stays on the attended per-idea migration; the declaration is disclosed in the new
  cycle policy/inspection as `unknown-history` with its decision/digest, never as a number.
- **Compatibility:** an older binary ignores the ledger; a newer binary without a ledger
  behaves exactly as today — the ledger only narrows refusals for declared, byte-identical
  runs.

**Remaining owner-only decisions: none in design, three owner-only acts.** If this contract
is accepted, design is settled by participant convergence. What remains are acts, not open
choices: (a) the owner attending the one-time activation apply after review — and until it
happens, every new idea still hits D6, so the release handoff must not claim the legacy gate
is activated; (b) the standing release publication commands (npm/core) the brief already
reserves to the owner; (c) the pre-authorized attended-close conditions check at close time.
I see no fourth.

**D2 — accepted with one test-level refinement.** The ceiling `min(active-track timeout,
positive configured checker timeout)`, missing configured timeout → track value,
5/15/30 min per fast/standard/deliberation, absent track → the protocol's standard default,
malformed track fails closed; resolved once per logical step and shared by both attempts;
`RunParticipantStep` for execution including protected checkers; valid FAIL final; invalid
own output or classified child/watchdog failure gets at most the one existing retry; a valid
PASS from an unsuccessful process still cannot close; membership/dropout policy unwidened —
all confirmed against the code I re-read (`driver_impl.go:838` literal, `wait.go`
`trackTimeoutCeiling`, the existing two-attempt wrapper at `driver_impl.go:854-890`). The
refinement: the existing `trackTimeoutCeiling` deliberately lumps *unknown* tracks into the
most permissive 30-minute branch (`wait.go:58-68`); the new goal-check resolver must not
inherit that branch for a malformed track. Absent (no field) → standard 15 min; unparseable
or out-of-vocabulary → fail closed. That distinction needs an explicit test, not a comment.

**V-GOAL — agreed.** FINAL should cite the independently supported >120 s observation; the
518 s figure stays marked testimony unless a verifier opens the origin/main receipt.

**D3 deferral — conceded.** Your two-layer counter is right and I verified both layers:
`consensus_request_signoffs.go:170` mints newer `consensus-signoff` runs that shadow the
driving run in `ResolveRun` (`runstate.go`, first idea-matching error-free run), while the
"No recoverable action" text comes from the read-only `runplan.Plan`
(`runplan.go:213`) — so my one-condition resolution fix would at best swap which run the
planner shrugs at, and at worst act on stale frozen roster/safety counters. A coordinated
resolver+planner fix is a behavior change in the runner core, not a small unstall. Deferral
is correct under "keep it small."

**D4 deferral — conceded.** My round-01 lean ("smaller diff") answered the wrong question:
the risk in reopening a *signed, ready* consensus is lifecycle semantics — dissolving an
accepted decision its signers already ratified — not diff size. The byte-identical rename
workaround both predecessors used preserves content but bypasses every verb, which is
precisely why the proper fix (a checked archive operation tied to a completed fix cycle and
fresh replacement review) needs its own idea. Deferred, with the hazard recorded.

## New concerns / questions

1. **Activation window.** Between release and the owner's attended apply, D6 continues to
   stall every new idea. The release note, `IMPLEMENTATION.md` and the released-inbox
   handoff must state this plainly (D6 fixed in product, activation pending owner act), and
   the prepared activation command with its exact preview must be in the handoff.
2. **Eligibility predicate precision.** "Structurally valid identity-absent history" must
   be enumerated in FINAL as a closed predicate (path shape, file inventory, event
   vocabulary, cursor absence, regular files), because "structural" without an enumeration
   drifts into judgment the gate cannot re-verify. If the enumeration and the May run ever
   disagree, the run is not declarable — the predicate wins.
3. **Fail-closed tests, both gates.** D1: undeclared unscoped run refuses; declared-but-
   mutated bytes in any root refuses; declared-but-deleted/all-missing refuses; conflicting
   decision-ID reuse refuses; a second unscoped run appearing refuses; a declared run that
   gains a cursor or charge event refuses (out of region); identical replay is idempotent;
   inspection discloses `unknown-history` + decision/digest; no driver-written declarations.
   D2: timed-out/stalled check leaves completion unverified and escalates; valid FAIL is not
   retried; watchdog/child failure gets exactly one retry; PASS from non-zero-exit process
   cannot close; buffered final-text transport keeps its hard deadline; absent track →
   15 min, malformed track fails closed; both attempts share the resolved ceiling; a
   protected checker's retry never becomes dropout authority; LE-5 budget hit escalates
   rather than completes. These are the witnesses I will hold Phase 6 review to.
4. **Convergence is now two-participant.** With kimi-1 excluded on the recorded owner-
   authorized dropout, every convergence here is correlated-agreement territory under
   §15.6(b): FINAL should say so and name what would falsify the design (the adversarial
   tests above are exactly that falsifier set).

## Current proposal

Adopt **D1 + D2 only**, as specified in my Responses section, with codex-1's round-02
contracts accepted as written and my pinned enumerations (structural admission predicate,
replay refusal list, D2 absent/malformed track rule) folded into FINAL verbatim or tighter:

1. **D1 durable unknown-history declaration** — attended apply + read-only inspect;
   append-only, one record per decision in the git-common-dir budget area; structural-only
   admission region; manifest-digest-bound with per-root revalidation on every gate use;
   frozen at first binding; disclosed as `unknown-history` with decision/digest; existing
   charges/caps/unavailable-root gates untouched; per-idea migration remains the path for
   charged history; older binaries unaffected.
2. **D2 derived goal-check ceiling** — `min(track, positive configured)` per D2, absent
   track → standard, malformed track → fail closed (not the permissive unknown branch);
   one resolved ceiling for both attempts through `RunParticipantStep` including protected
   checkers; valid-FAIL finality; LE-7/LE-11, current-tree evidence, dissent, reservations
   and reviewer gates unchanged; the two normative sentences (Phase 8 "two-minute product
   ceiling", §9.0 "goal120s … stay") replaced in all synchronized authority copies.
3. **Deferrals recorded in FINAL** with evidence: D3 resolver/planner (two-layer mismatch,
   stale-counter hazard), D4 ready/partial reopen (lifecycle semantics, manual-rename
   hazard), plus the round-01 deferred gap list. Each names its follow-up.
4. **Fail-closed test witnesses** as enumerated in New concerns §3, added as acceptance
   criteria so both gates carry regression-proof refutation attempts, not just happy paths.
5. **Release**: CLI 1.54.0, skill/core 2.18.0 staged from 2.17.0, protocol hunks
   synchronized across live, bootstrap, skill-reference and staged-core copies; handoff
   states the activation-pending status of the D6 fix and carries the prepared activation
   command.

No open owner-only design decision remains if codex-1 accepts this round's pinned
enumerations; the three owner-only acts are listed under Responses/@codex-1. If codex-1's
round-03 (if any) disputes an enumeration, that is the only surface I would still argue
rather than concede, because every enumeration above is a fail-closed boundary, and widening
any of them is the specific failure mode this idea exists to prevent.
