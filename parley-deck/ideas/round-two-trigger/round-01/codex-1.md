---
agent: codex-1
idea: round-two-trigger
round: 1
date: 2026-09-01
---

## Summary

Do not claim to detect semantic convergence. **PRIMARY evidence:** the condition “no substantive
disagreement” has no executable match in `internal/**/*.go`; the inspected boundary code schedules
rounds by configured count (`internal/driver/driver.go:287-343`). The checkable problem is narrower:
who authorizes a one-round close, and whether that authorization is visible.

The facilitator should remain allowed to **propose** closing after round 1, because the frozen
measurement says the present judgment usually opens round 2 (52 of 80 ideas). The facilitator
should not be allowed to decide alone. For a `standard` or `deliberation` idea that attempts
consensus directly from round 1, every participant should make an explicit append-only
`CLOSE`/`OPEN` round-two decision in its existing consensus signoff. Any `OPEN` is the existing
`❌ BLOCK` plus a counter-proposal and therefore opens round 2 through the existing reopen path.
This adds no model call when everyone says `CLOSE`, adds no artifact class, and does not affect
`fast`.

This is a protocol-semantic change and therefore needs core 2.12.0 **and** matching prompt/parser
enforcement. It must not ship as a quiet CLI default.

## Proposed approach

### 1. Replace a convergence detector with a closure-authorization gate

Proposed normative rule for core 2.12.0:

> For an idea carrying `round_two_policy: quorum-veto-v1`, when `track` is `standard` or
> `deliberation` and design consensus is drafted from round 1, the facilitator's close call is a
> proposal only. Every required consensus signer MUST append `Round-02 decision: CLOSE|OPEN` and a
> non-empty `Round-02 reason:`. `OPEN` MUST use `Status: ❌ BLOCK` and a concrete counter-proposal
> naming what round 2 should test. Consensus is not ready while a decision is missing or
> inconsistent. `fast` and consensus drafted after round 2 or later are unchanged.

The validator checks only the enum, presence, signer coverage, and consistency with signoff
status. It does **not** certify that a reason is wise or that convergence exists. The protocol
should say that limitation as explicitly as current §15.6 says only clause (a) is machine-validated.

The authority split becomes:

- Facilitator: may open round 2 directly or propose a one-round close by drafting consensus.
- Quorum: authorizes the close; any participant may force cross-review using the veto it already
  possesses through `❌ BLOCK`.
- Owner: remains the attended protocol publisher and the destination for genuine escalation; the
  owner is not interrupted on every ordinary close.

### 2. Carry it in the existing signoff path

The implementation should change the protocol, the consensus signoff prompt, the CLI signoff
surface, and the parser together:

1. `consensus draft` records the selected source round and emits a `## Round-two transition
   record`. From round 1 it says `ballot-required`; from round 2+ it says `already-opened`.
2. The signoff prompt and `parley consensus signoff` accept the two exact round-two fields. An
   `OPEN` paired with ACCEPT/RESERVATIONS, or a `CLOSE` paired with a round-two-specific BLOCK, is
   malformed.
3. `consensus status` reports `pending`, `close-authorized`, `open-requested`, `already-opened`,
   `not-applicable-fast`, or `legacy-unmeasured`. It derives this state from canonical files; no
   new state file is needed.
4. An `OPEN` leaves normal triage at `blocked`. `parley consensus reopen --reason` preserves the
   blocked draft, and the existing driver/manual workflow opens round 2.

The audit trail is therefore complete in both directions. All `CLOSE` decisions remain in
`consensus.md`; an `OPEN` remains in the archived blocked consensus plus `round-02/`; a round opened
before consensus is visible as `already-opened`. No positive “evaluated” claim is inferred for
legacy ideas.

### 3. Make compatibility explicit

**PRIMARY evidence:** `parley protocol status` on 2026-09-01 reports `installed: 2.10.0`, no deck
pin, while `parley-deck/meta/version.json:2-8` describes this source deck as 2.11.0. The kickoff
records that 2.11.0 is staged but unpublished and that 40 of 41 decks lack §15.6
(`00-prompt.md:80-87`).

Therefore:

- Stage 2.12.0 on top of the 2.11.0 candidate, but do not publish either core from an agent run.
  The owner first resolves the 2.11.0 publication state, then publishes attended.
- A 2.12 prompt/template writes `round_two_policy: quorum-veto-v1` automatically for new explicit
  `standard` and `deliberation` ideas. It is a capability marker, not a user opt-in knob.
- Existing open ideas are not backfilled. This avoids changing their deliberation semantics
  midstream while general per-idea protocol pinning is still unimplemented.
- On a deck/idea without the marker, new CLI code preserves legacy behavior and reports
  `legacy-unmeasured`; it may warn, but it must not block, synthesize a ballot, or silently open a
  round. After an attended core sync, the policy applies to newly opened ideas.
- On a 2.12 deck, kickoff/preflight and `consensus draft` must fail closed for a newly created,
  eligible idea that omits the marker. Otherwise omission becomes the quiet opt-out that the
  carrier thesis predicts.

This dormant-until-marked behavior lets the CLI implementation land without changing any of the
40 legacy decks, while ensuring that prose and enforcement activate together.

### 4. Pre-register the corpus test

Use the first frozen cohort of 12 eligible, non-legacy attempts to close `standard` or
`deliberation` ideas from round 1. Report at corpus level:

- attempted one-round closes, complete ballots, `OPEN` ballots, and ballot-opened round 2s;
- participant position changes and alternatives/decisions changed after those rounds;
- generation-call count and wall time, with unknown telemetry reported as unknown rather than
  zero.

Do not report per-idea pairwise-distance success. Treat the gate as ceremonial and open a repeal
idea if the cohort contains zero `OPEN` ballots, or if every ballot-opened round 2 records no
position change and changes no alternative disposition or agreed decision. A positive cohort does
not authorize expansion to `fast` or a mandatory round 2; that would be another protocol decision.

## Existing alternatives

### Components this proposal would build by hand

| Hand-built component | Closest mechanism that already ships (verified locator) | Why not reuse unchanged | Origin |
| --- | --- | --- | --- |
| `round_two_policy` capability marker and source-round classification | `cross_review_rounds` is already read from idea frontmatter, defaulting to 1 (`internal/driver/transport.go:32-41`); track policy is applied in `internal/track/track.go:152-217` | These schedule a count; neither records who authorized a one-round close. | **Constraint-forced** by legacy compatibility and the prohibition on a quiet CLI semantic change. |
| Closure-ballot signoff fields | Existing append-only signoff blocks and status vocabulary are emitted/parsed by `internal/consensus/consensus.go:231-285,871-881` | The block has status, notes, and counter-proposal, but no explicit round-two decision. | **Constraint-forced** for an auditable evaluated/not-evaluated record; append-only placement is **inherited**. |
| Closure-ballot prompt carrier | `buildConsensusSignoffPrompt` already tells each signer to read context and append one canonical block (`internal/app/consensus_request_signoffs.go:711-757`) | It never asks whether round 2 is needed. An unenforced prose duty elsewhere would repeat the measured carrier failure. | **Constraint-forced**. |
| Ballot consistency validator and status projection | `Status` and `validateDocumentAwaiting` already derive ready/blocked/partial/malformed from signer coverage and canonical statuses (`internal/consensus/consensus.go:96-124,523-590`) | They cannot distinguish an considered one-round close from a signoff that never evaluated cross-review. | Validator is **constraint-forced**; the triage machine is **inherited**. |
| Round-two transition record | `Draft` already selects a source round and checks that its artifacts are complete (`internal/consensus/consensus.go:159-198`) | The design draft template does not preserve that source round or a transition decision (`internal/consensus/consensus.go:815-832`). | **Constraint-forced** by auditability. |
| `OPEN` routing | `Reopen` requires blocked triage, archives the consensus with a reason, and resets idea status (`internal/consensus/consensus.go:370-411`); the auto-driver opens the next round on blocked consensus (`internal/driver/consensus.go:91-131`) | This already supplies the desired back-edge. No new round-opening mechanism is needed. | **Merely inherited**. |
| Corpus evaluation fields | `retro.Scan` already counts design rounds and failure signals (`internal/retro/retro.go:19-35,55-80`) | Its current score treats rounds beyond one as friction and classifies one-round ideas as low-friction (`internal/retro/retro.go:124-164`); it does not read track, closure ballots, or premature-close evidence. | **Constraint-forced** for the frozen cohort measurement; extending `retro scan` is preferable to a new scanner. |

### Existing mechanisms considered

- **§15.6's old trigger language — not an executable gate.** **PRIMARY evidence:** the quoted
  “if round 1 closes with no substantive disagreement” text is the opening trigger paragraph in
  the installed 2.10.0 core at
  `/Users/tomasfecko/.parley/protocol/core/2.10.0/COOPERATION.md:1341-1354`; clause (b) itself is the
  correlated-agreement record at lines 1356-1358. The current repository/default 2.11 text has
  deleted that trigger and says only (a) is machine-validated
  (`internal/protocol/defaults/COOPERATION.md:1339-1354`). A scoped search of `internal/**/*.go`
  for `substantive disagreement` and `Adversarial alternative` returned no executable match; the
  only round-one §15.6 validator checks `## Existing alternatives`
  (`internal/protocol/roundartifact.go:8-30`). Nothing acts on semantic disagreement.

- **`parley consensus status|draft|signoff` — chosen substrate, not a detector.** **PRIMARY
  evidence:** the command dispatch and flags are at `internal/app/app.go:516-628`. `draft` checks
  completeness of the selected round, `status` parses signoffs, and `signoff` appends a participant
  status. None decides whether another round would improve the decision.

- **`parley consensus reopen --reason` — chosen back-edge.** **PRIMARY evidence:** the CLI requires
  a reason (`internal/app/app.go:661-683`), while the package refuses unless triage is already
  blocked and preserves the old draft (`internal/consensus/consensus.go:370-411`). It reacts to a
  block; it does not discover convergence or disagreement.

- **Auto-driver round-boundary policy — useful but insufficient.** **PRIMARY evidence:**
  `advanceRound` opens another round or drafts consensus solely from the configured
  `CrossReviewRounds` budget (`internal/driver/driver.go:287-343`). Its `roundComplete` gate checks
  artifact validity and events, not positions (`internal/driver/driver.go:346-407`). The default
  automated path already schedules one cross-review round, but manual facilitation and an explicit
  zero remain outside that protection.

- **`parley retro scan` — wrong current signal for this question.** It is read-only and scans the
  structured corpus (`internal/app/retro.go:21-59`), but currently scores extra design rounds as
  friction rather than detecting a suspicious absence of cross-review. It is an appropriate
  measurement carrier only after adding explicit closure fields.

- **Inbox escalation — fallback, not the routine trigger.** **PRIMARY evidence:** the protocol's
  human escalation path is `inbox/<from>-to-user_...` (`parley-deck/COOPERATION.md:690-727`), and
  auto-driver errors/budgets write blocking notes (`internal/driver/loop.go:202-227`). Routine
  escalation would move every close decision to the owner, add human latency, and still provide no
  convergence test. Use it only for malformed/conflicting state or an actual human-only choice.

- **`parley loop tick` — not applicable.** **PRIMARY evidence:** it is disabled by default and may
  only draft non-active candidate prompts; it never starts, advances, reopens, or finalizes a
  deliberation (`internal/app/loop_cmd.go:14-18,61-99`; `internal/loop/loop.go:166-192`). Reusing it
  would cross the §14 human brake.

- **`parley consult` / a special role — weaker than the existing quorum.** Consults are advisory
  and non-canonical (`parley-deck/COOPERATION.md:815-818`), and `roles:` cannot change quorum,
  ownership, or signoff weight (`parley-deck/COOPERATION.md:89-98`). A new consult adds a call but
  still needs promotion into canonical state. The existing participant signoffs already have the
  correct authority.

- **Keep the current facilitator judgment unchanged — credible but incomplete.** **PRIMARY frozen
  evidence:** 52 of 80 ideas did open round 2, so the judgment is mostly right
  (`00-prompt.md:56-60`). This proposal preserves it as the initial recommendation and adds no full
  round unless a participant disagrees. What it removes is unilateral, unrecorded closure.

- **Mandatory round 2 for every `deliberation` idea — deterministic fallback, rejected now.** It is
  easy to validate and the four high-rigour single-round cases make it tempting, but it flips the
  cost default on the strength of a small cohort. The closure ballot tests the narrower governance
  defect first. Mandatory round 2 remains the counter-proposal if the ballot is measured as ritual.

## Concerns / open questions

- The ballot may ritualize into four `CLOSE` lines. The frozen cohort and repeal condition are
  therefore load-bearing, not retrospective decoration.
- Related models voting separately are still one shared prior, not independent evidence. The
  existing §15.6 correlated-agreement record must remain; the ballot redistributes authority but
  does not manufacture epistemic independence.
- The facilitator still executes and pays for a veto-triggered round. This proposal prevents that
  cost conflict from authorizing closure, but it does not make the voter bear the bill. Existing
  wall-clock/step ceilings remain the proper safety brake.
- Manual appenders can bypass CLI ergonomics. The parser must validate the canonical file, so hand
  editing cannot produce a ready triage without the required fields on a marked idea.
- The implementation must define how a 2.12 deck proves that a newly created idea is eligible for
  the marker without pretending general per-idea protocol pinning exists. A narrow, explicit
  capability marker plus no retroactive backfill is the safest boundary I see.
- Should rollout initially cover only `deliberation`? I favor both explicit `standard` and
  `deliberation`: the old trigger named both, the added decision rides an already-required signoff,
  and only an `OPEN` incurs a round. `fast` remains excluded.

## Risks

- **False reassurance:** machine-valid ballot syntax may be presented as machine-detected
  convergence. The protocol and CLI output must use “closure authorization,” never “converged.”
- **False negatives:** all correlated participants may close early. This remains possible and must
  be visible as unanimous shared-prior judgment, not four independent confirmations.
- **False positives / cost:** one weak `OPEN` can force a full round. Requiring a concrete test in
  the existing BLOCK counter-proposal creates accountability without giving the facilitator power
  to overrule it.
- **Version drift:** enabling the parser by CLI version instead of an idea/deck capability marker
  would silently change old decks and violate the governing doctrine. Legacy must be reported, not
  coerced.
- **Split manual/automatic semantics:** the driver already schedules by count while manual
  facilitation uses judgment. The gate must bind at `consensus status/finalize`, the common
  canonical boundary, rather than only inside `parley run`.
- **Measurement inversion:** current `retro` scoring treats more design rounds as friction. Using it
  unchanged would make every successful safety-triggered round look worse; the new fields must
  separate deliberate safety cross-review from churn.
- **Telemetry overclaim:** call counts and wall time are the only allowed cost measures here;
  `internal/driver/loop.go:174-175` says provider usage events are currently absent in practice.
