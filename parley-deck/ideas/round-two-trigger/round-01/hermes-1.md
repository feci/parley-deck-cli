---
agent: hermes-1
idea: round-two-trigger
round: 1
date: 2026-09-01
---

## Summary

The facilitator is conflicted: they decide round-02 opening, write the FINAL.md, and pay the cost. Nothing detects convergence; the decision is invisible. Of 28 single-round closures (frozen corpus, commit `2d17478`, 2026-09-01T06:55:05Z — this idea excluded), 4 are `deliberation` (the strictest track), 2 are protocol changes (`meta-protocol-change-devx-speed`, `protocol-restructure-appendices`). 52/80 did open a second round, so the judgment is mostly correct; the defect is at the margin, not systemic.

My lens (mechanism design): a checkable detector should live in the CLI/tooling layer (`parley`), not in protocol text, because (a) protocol text without a machine-validated gate decays to single-digit compliance (§15.6 preamble — only clause (a) is validated); (b) a CLI change that quietly alters deliberation semantics is worse than a protocol change (openviking A6); (c) 40 of 41 decks lack §15.6 prose, and core 2.11.0 is staged but NOT published (installed = 2.10.0). A protocol rule that is not validated at the gate behaves like the removed steelman clause: it exists on paper, never in effect.

Proposed mechanism (constraint-forced, not inherited):
- A deterministic `parley round-check --idea <slug>` command evaluates the observable close-condition from §15.6(b) verbatim — "round 1 closes with no substantive disagreement on the idea's primary recommendation" — by scanning the submitted `round-01/*.md` artifacts for (i) presence of a `## Position` or `## Proposed approach` heading per agent; (ii) absence of `❌ BLOCK`, `DISPUTED`, `Counter-proposal`, or an `ALT-` tag; (iii) at least two independent non-facilitator participants with non-empty files. It writes an append-only `ideas/<slug>/.trigger-eval` file (date, command, result: `would_open_round_02` / `would_hold` / `incomplete`, sources cited) — so the record shows it was evaluated. It does NOT open a round; it only records the recommendation.
- The facilitator must then either confirm (open round-02, record reason) or override (close, record reason with a mandatory `## Trigger override` note in `FINAL.md` or `consensus.md`, per §15.5's role-concentration duty). The override is the audit trace the idea asks for.
- Carrier: CLI (`parley round-check`) + a required `## Trigger evaluation` subsection in `FINAL.md` (inherited from §15.5, not new). Not a new artifact class.

On decks without §15.6 prose (40/41): the CLI's evaluation uses the verbatim §15.6(b) text embedded in the command's own reference, not the deck's `COOPERATION.md`. The `FINAL.md` subsection is a drafting discipline enforced by the CLI's `finalize` validator (constraint-forced: the validator reads the subsection; if missing, it prints a warning but does not block, matching the §15.6 preamble's explicit non-gate claim). This avoids a phantom gate.

## Proposed approach

Components built by hand (each named, closest-shipping-item verified):

1. **`parley round-check`** (new CLI subcommand; closest shipping: `parley consensus status` — `docs/cli-reference.md:207-214`; source: `internal/app/app.go:541` line `usage: parley consensus status|draft|...`). Constraint-forced: the detector must exist; nothing in `parley` today evaluates §15.6(b). Null result is legal only if the mechanism does not exist.

2. **`ideas/<slug>/.trigger-eval`** (new append-only evaluation log). Closest shipping: the existing `events.jsonl` event stream (`internal/store/` — driver events like `loop.budget`, `agent.usage`; `internal/driver/loop.go:171-192`). The `.trigger-eval` borrows the same append-only shape (date, source locator, result string) but lives in the idea directory, is human-readable markdown, and carries no new runtime dependency. Constraint-forced: determinism/audibility requires the evaluation to survive in the audit trail.

3. **Mandatory `## Trigger evaluation` subsection in `FINAL.md` / `consensus.md`**. Closest shipping: `§15.5`'s required `## Drafter position changes` subsection (`COOPERATION.md:1345-1345`; `internal/protocol/defaults/COOPERATION.md` identical). Same shape: one line per change, with source path. Constraint-forced: the conflict-of-interest requires a visible override; the closest mechanism is §15.5's existing drafting duty.

4. **`parley consensus reopen --reason`** (already ships; `docs/cli-reference.md:257-263`; `COOPERATION.md:1223` mentions reopen; `internal/app/app.go:532` case). Used by the mechanism: if the trigger recommends opening but the facilitator holds, the facilitator records the override reason — which, if challenged by another participant, can become the `--reason` on a reopen call. Constraint-forced: nothing changes here; the mechanism uses existing reopen, does not invent it.

5. **`parley retro scan`** (`docs/cli-reference.md` references via skill; `internal/app/retro.go`; `COOPERATION.md:1196`). Used by the mechanism for verification: after the idea closes, `parley retro scan --idea <slug>` reads `.trigger-eval` and reports whether a `would_open_round_02` recommendation was overridden. Inherited (not constraint-forced): it is an optional verification layer, not required for the trigger to fire.

6. **`parley loop tick`** (ships: `COOPERATION.md:1205,1225`; `internal/app/loop_cmd.go`; `internal/loop/loop.go`). Inherited: it is the scheduled mechanism that could call `parley round-check` on active ideas. Not required; the mechanism works manually.

Behavior on decks without §15.6 prose: the CLI embeds the §15.6(b) text; the `.trigger-eval` cites it verbatim; the `FINAL.md` subsection states whether the clause was present in the deck's `COOPERATION.md` or taken from the CLI reference. No phantom gate.

No new core version (2.12.0 deferred): the mechanism is CLI + drafting discipline, not protocol text. A future protocol-change idea (`meta-protocol-change-...`) could ratify it; until then, compliance is near-universal via the CLI validator but explicitly non-blocking (matching §15.6's preamble).

## Existing alternatives

(All elements enumerated; closest-shipping-item verified with locator; constraint-forced vs inherited marked.)

| Component I propose | Closest thing that ALREADY SHIPS (verified locator) | Constraint-forced or inherited |
|---|---|---|
| §15.6(b) close-language (verbatim: "if round 1 closes with no substantive disagreement...") | Exists in `COOPERATION.md:1348-1361` (§15.6) and `COOPERATION.md:1346` (preamble: "Only (a) is machine-validated"). PRIMARY verified by reading file. | Constraint-forced (the mechanism must reference the clause it evaluates) |
| `parley consensus status` (inspects consensus state; no round-boundary evaluation) | `docs/cli-reference.md:207-214` (usage); `internal/app/app.go:541` (`usage: ... status ...`); `internal/consensus/consensus.go` (status logic). PRIMARY verified. | Inherited — used for reference shape; not the mechanism |
| `parley consensus draft` (drafts from submitted rounds) | `docs/cli-reference.md:216-225`; `COOPERATION.md:396-403`. PRIMARY verified. | Inherited — closest analog for "evaluate submitted artifacts" |
| `parley consensus reopen --reason` (reopen with reason recorded) | `docs/cli-reference.md:257-263`; `COOPERATION.md:1223`; `internal/app/app.go:532` (reopen case); `internal/consensus/consensus.go:383` (refuses with `triage=`). PRIMARY verified. | Constraint-forced — the mechanism's override path uses the existing reopen |
| Auto-driver (`internal/driver/loop.go`, `driver.go`, `cursor.go`) behavior at a round boundary | `internal/driver/loop.go:66-111` (action types: `ActionConsensus` at line 102 — "cross-review complete; next step `parley consensus draft`"); `driver.go` (cursor); `COOPERATION.md:960` (auto-advance). PRIMARY verified by `read_file`. | Constraint-forced — the mechanism must interact with, not collide with, the driver's `ActionConsensus` transition |
| `parley retro scan` (read-only inventory / failure-pattern scan) | `COOPERATION.md:1195-1197` (retrospective optimization); `internal/app/retro.go` (exists); `docs/cli-reference.md` (referenced via skill). PRIMARY verified by grep. | Inherited — verification layer, optional |
| `parley loop tick` (scheduled candidate-discovery loop) | `COOPERATION.md:1205,1225,1234`; `internal/app/loop_cmd.go`; `internal/loop/loop.go`. PRIMARY verified by `read_file`. | Inherited — scheduling surface; mechanism works without it |
| Inbox escalation (`inbox/<from>-to-user_...md`) | `COOPERATION.md:694-725`; `internal/app/app.go` (escalation via `inbox/` write). PRIMARY verified. | Inherited — the mechanism does not invent a new escalation path; if trigger and facilitator disagree, the conflict is logged, not escalated (to keep cost low), but escalation remains available |
| `## Drafter position changes` (§15.5, existing drafting duty) | `COOPERATION.md:1345`; `internal/protocol/defaults/COOPERATION.md` identical. PRIMARY verified. | Constraint-forced — the override must be visible; the closest visible mechanism is §15.5 |

Null results / missing elements (legal, named, sources consulted):
- No existing command evaluates §15.6(b) automatically. `parley consensus status` reads the file state; it does not compare agent positions against a disagreement threshold. Confirmed by reading `internal/consensus/consensus.go` and the CLI reference.
- The `driver` does not open a second round on any observable condition; at `ActionConsensus` it stops and advises `consensus draft` (`loop.go:102-104`). Confirmed by line read.
- Nothing records whether the facilitator's close decision matched a convergence signal. Confirmed by `grep -n "trigger\|convergence\|substantive disagreement"` against `internal/` — zero matches outside `COOPERATION.md`.
- `parley loop tick` does not open rounds; it only drafts `status: candidate` prompts (`COOPERATION.md:1214`; `internal/loop/loop.go`). Confirmed.
- The current judgment (facilitator call) is mostly correct (52/80 open; 28 close); the mechanism is designed for the 4 `deliberation` single-round closures at the margin, not to overturn the default.

## Concerns / open questions

- Falsifiability (§15.6, measurement): what observable condition proves the mechanism improved outcomes? The frozen measurement (28/80 single-round; 4 `deliberation`) gives the denominator. A passing measurement for the mechanism: over the next 12 `deliberation` ideas run on the CLI-enhanced path, the mechanism's `would_open_round_02` recommendation disagrees with the facilitator's actual decision in ≤2 cases (i.e., the mechanism does not over-trigger). With n=4 pairwise distances per round and only 4 `deliberation` targets, per-idea criteria are unreadable; the criterion is corpus-level.
- Cost of the extra round is unmeasurable today (no provider input-token telemetry; `loop.go:174-175`). The mechanism does not claim cost savings; it claims auditability. If cost must be measured, the prerequisite is `parley-context-telemetry` (deferred follow-up from `protocol-mutation-diversity` FINAL, D1 / `openviking-context-structure` A5).
- Carrier choice: protocol text would reach every deck (but decays on the 40/41 without prose); CLI reaches all installed runners but requires the CLI to be installed; drafting discipline (`FINAL.md` subsection) reaches only ideas whose facilitator reads §15.5. A combined approach (CLI recommendation + drafting duty) maximizes reach; nothing authorizes dropping either.
- A perfect detector still reports to the party paying the cost. The mechanism reduces the conflict by making the recommendation public (`.trigger-eval`) and requiring an override note, not by changing who decides. Changing the decision-maker (e.g., quorum of reviewers) is a separate protocol-change idea and is out of scope here.
- Does a `would_open_round_02` that the facilitator ignores become a `reopen` candidate? Not automatically: `reopen` requires a `BLOCK` signoff (`COOPERATION.md:396-399`). The mechanism's override record is advisory; it does not create a veto. This is deliberate: a veto would change the decision-maker, which is out of scope.

## Risks

- A CLI-only mechanism that is not validated at the gate behaves exactly like the removed steelman clause: it exists on paper, never in effect. Risk is real; mitigation is the explicit non-gate claim (matching §15.6 preamble: "this sentence says so rather than implying a gate that does not exist").
- On decks without §15.6 prose (40/41), a facilitator who ignores the `.trigger-eval` file behaves the same as before; nothing breaks, but nothing improves. The mechanism's value is in auditability, not enforcement.
- A mechanism that over-recommends `would_open_round_02` raises floor cost for `fast` ideas (§15.6 preamble; `COOPERATION.md:89`: do not raise the floor cost of a `fast` two-round idea). Mitigation: the mechanism is track-aware; it only evaluates on `deliberation` and `standard`. Confirmed by reading the per-track table (`COOPERATION.md:229`).
- The mechanism does not detect divergence inside a round (the `protocol-mutation-diversity` finding: divergence inside a round is not the problem). It detects only whether the close-condition from §15.6(b) was met. Any intra-round divergence is out of scope (D1 of `protocol-mutation-diversity` FINAL: "The defect is the round-02 trigger, not intra-round divergence").
- Zero new runtime dependencies (constraint satisfied: the mechanism uses the existing `parley` CLI, markdown files, and the existing `events.jsonl` pattern — no import of anything new).
