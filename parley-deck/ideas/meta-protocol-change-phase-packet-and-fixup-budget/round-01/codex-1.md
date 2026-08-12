---
agent: codex-1
idea: meta-protocol-change-phase-packet-and-fixup-budget
round: 1
date: 2026-08-11
---

## Summary

Adopt generated, ephemeral phase packets, but only with a non-optional safety kernel, an omission ledger, and a fail-open-to-the-full-protocol fallback. The packet must be an exact extraction from the resolved live protocol, not a paraphrase or a checked-in derivative. The source mapping should live as machine-readable applicability metadata beside the normative source blocks; unclassified normative text is included in every packet and fails the coverage check.

Set the Phase 8 fix-up caps to `fast=1`, `standard=2`, and `deliberation=6`. Reaching a cap with work still open creates a blocking user escalation and halts automation; it never means complete. A user may authorize a new explicit finite ceiling after reviewing the trajectory.

Evidence used:

- **PRIMARY** — stable locator `parley-deck/COOPERATION.md`, §4.0, table **Per-track behavior**: the present caps are fast 1, standard 2, and deliberation unbounded; the same section says §14 is an invariant on every track.
- **PRIMARY** — stable locator `parley-deck/COOPERATION.md`, §4 Phase 2, **Rules**: “Silence = implicit agreement.”
- **PRIMARY** — stable locator `parley-deck/COOPERATION.md`, §6 rule 3: an agent must never edit another agent's file except under the recorded direct-user exception.
- **PRIMARY** — stable locator `parley-deck/COOPERATION.md`, §15.2, **Provenance**: an unlocated `PRIMARY` is malformed and reads as `RECALL`.

## Proposed approach

### 1. Section-to-phase mapping

Every packet starts with a generated envelope containing: protocol source SHA-256, protocol/deck version if available, packet schema and generator version, idea path and `00-prompt.md` SHA-256, phase, track, active transport, participant/quorum list, current `status`, and all present routing flags (`strict_gate`, `auto_implement`, `require_model_diversity`, and `checks`). The packet generator must not inspect another participant's round-01 file while producing a Phase 1 packet.

The following is the initial conservative mapping. “Load” means exact source blocks are included. “On demand” means the packet's omission ledger names the stable locator and the machine-readable trigger that requires regeneration with that section; it does not mean the rule is silently absent.

| Phase | Load-bearing protocol content | Conditional or reachable on demand | Reason |
| --- | --- | --- | --- |
| 1 — independent analysis | §4.0; §4 Phase 1; §5; §6 rules 1 and 3–6; §4 **Escalation to user** plus the canonical-inbox part of §8; §14; §15; active transport's §11 Phase 1 paragraph | §7 is loaded for a protocol-change idea; §12 is loaded for a pipeline/action-block idea; §13 is reference-only | Independence, ownership, status, quorum, evidence, and escalation are rules the participant can violate while writing the first artifact. §15 is explicitly bound by Phase 1. |
| 2 — cross-review | §4.0; §4 Phase 2; §5; §6 rules 1–6; escalation/inbox; §14; §15; active §11 Phase 2 paragraph | §7 for protocol changes; §12 for pipeline/action-block propositions | “Address everyone,” silence, counter-proposals, source sharing, verdict conflict, and async quorum all affect whether apparent convergence is real. |
| 3 — design consensus | §4.0; §4 Phase 3; the Phase 4 drafter/close preconditions; §5; §6 rules 2, 3, 5, and 6; escalation/inbox; §14; §15; active §11 Phase 3 paragraph | §7 for protocol changes; the rest of Phase 4 is reachable when finalization begins | Signoff semantics, append-only ownership, disputed claims, role concentration, correlated agreement, and the correct quorum are close gates, not background reading. |
| 5 — implementation | §4.0; applicable LE-2/4/5/7/11 definitions; §4 Phase 5; FINAL's frozen/source-of-truth and recovery contract from Phase 4; §6 rules 3, 5, and 6; escalation/inbox; §14; active §11 Phase 5 paragraph | §12 becomes load-bearing for pipeline/action/driver-managed work; §15 is fetched before assigning a material verification verdict; §13 is fetched for a retrospective | The implementer needs the authority boundary, plan/deviation/validation duties, human gates, and publication mechanics. Full review rules are not yet needed. |
| 6 — code review | §4.0; LE-1 and LE-3; §4 Phase 6 including **Review briefs and dispositions**; Phase 2's later-round response rules; §5; §6 rules 1 and 3–6; escalation/inbox; §14; §15; active §11 Phase 6 paragraph | §7 for review of a protocol change; §12 for a pipeline/action implementation | Refutation attempts, no-suppression, fixed severities, reviewer independence, provenance, and cross-review ownership directly determine review validity. |
| 7 — review consensus | §4.0; §4 Phase 7; Phase 3's exact signoff block and status meanings; §4 Phase 8's cap precondition; §5; §6 rules 2, 3, 5, and 6; escalation/inbox; §14; §15; active §11 Phase 7 paragraph | §7 for a protocol change; §12 when pipeline state/effects are being dispositioned | The drafter must not turn disagreement into a dismissed finding, originate somebody else's verdict, miscount signoffs, or start a fix-up beyond budget. |
| 8 — fix-up | §4.0; applicable LE-2/4/5/7/11 definitions; all of §4 Phase 8, including strict gate, stopping judgment, budgets, and close-decision integrity; §6 rules 3, 5, and 6; escalation/inbox; §14; active §11 Phase 8 paragraph | §12 for pipeline/action/driver-managed work; §15 is fetched before assigning or relying on a material completion verdict; §13 after completion | Fix scope, validation, cycle accounting, strict-close conditions, escalation, and transport completion are inseparable. Omitting the stopping rules would recreate the unbounded loop through the packet itself. |

Treatment of the four historically important sections is therefore explicit:

| Section | 1 | 2 | 3 | 5 | 6 | 7 | 8 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| §15 verification integrity | load-bearing | load-bearing | load-bearing | on-demand verdict trigger | load-bearing | load-bearing | on-demand verdict trigger |
| §7 protocol change | load-bearing when the idea is a protocol change; otherwise on demand | same | same | same | same | same | same |
| §6 rule 3, no cross-editing | load-bearing | load-bearing | load-bearing | load-bearing | load-bearing | load-bearing | load-bearing |
| §14 human brake | load-bearing | load-bearing | load-bearing | load-bearing | load-bearing | load-bearing | load-bearing |

For this idea, §7 is therefore loaded in every listed phase. §14 stays loaded even in a human-driven session: the protocol itself calls the human brake an invariant on every track, and a packet should not rely on an agent correctly inferring whether its invocation came from a standing loop.

Reference-only material is: Quickstart and §10 after routing; §0's transport-selection discussion after bootstrap (the active transport value still stays in the envelope); §1's general purpose prose, sizing advice, and helper explanation after participants lock (the non-solo close guard remains in the safety kernel); §2's roster view after `participants` locks; §3's complete directory tree when the phase block already supplies exact artifact paths; the full §9 checklist after a separate startup packet has run; inactive §11 transport subsections; Appendix A; §12 unless a pipeline/action/driver trigger is present; and §13 unless running a retrospective. Reference-only locators remain visible in the omission ledger.

### 2. Generation and honesty

Add a proposed instruction-layer command such as:

```text
parley protocol packet --idea <slug> --phase <1|2|3|5|6|7|8> [--explain]
```

It should use the same resolved protocol composition and drift check as `parley protocol check`, parse Markdown structurally, and emit exact source blocks. It must not create another checked-in protocol file. A cache, if needed, is content-addressed and disposable; the packet is prompt input, not authority.

Applicability metadata belongs adjacent to the canonical source block and names stable rule/block IDs, applicable phases, feature triggers, dependencies, and whether the block is safety-kernel or on-demand. The generator computes the transitive closure. It never maintains copied prose. A normative block with no classification is `UNCLASSIFIED_INCLUDED`: it is included in every packet and makes `packet check` fail until a protocol-change review classifies it. Parser failure, unknown phase/track/flag, source drift, an unresolved dependency, or a packet/source hash mismatch falls back to the complete resolved protocol and records `context-mode=full-fallback` plus the reason.

`--explain` and the packet footer must report both sides of the cut:

- every included stable locator, block hash, and inclusion reason;
- every omitted stable locator, block hash, classification, and on-demand trigger;
- unclassified blocks, missing dependencies, and stale annotations;
- the safety-kernel set and whether it is complete;
- the exact source and idea hashes used.

`parley protocol packet check --all` should build every phase × track × transport × relevant feature-trigger combination and enforce: no unknown normative block is omitted; every dependency is present; excerpts hash to the resolved source; the four-section matrix above holds; only the selected §11 transport is included; and a packet can be reproduced from its envelope. Golden tests should assert the stable IDs, not copied prose.

### 3. Instruction changes

All three instruction paths must change together:

1. The skill's standing rule becomes: load the startup envelope, generate the current phase packet, verify its hashes/coverage result, and read the full resolved protocol on any failure or unknown trigger.
2. §9 step 1 becomes a small startup packet: determine transport, open idea, status, phase, track, flags, inbox obligations, and packet hash; then load the phase packet. §9.0 remains a facilitator/pre-idea packet rather than being repeated in participant packets.
3. Facilitator participant/reviewer/implementer prompts are rendered from one template that requires `packet-id`, source SHA, phase, track, and `context-mode`. The prompt embeds or attaches the generated packet; it does not contain a hand-maintained section list.

For official `parley` launches, missing packet attestation is a hard pre-launch error. The prompt renderer also detects an unconditional “read all of COOPERATION.md” directive and requires the launch to be explicitly marked `full-fallback` with a reason. Launch summaries and telemetry expose `context-mode=packet|full-fallback`, so a cost regression cannot be silent.

This cannot stop a person or arbitrary direct CLI invocation from hand-writing a prompt that requests the full protocol. That path is outside an instruction-layer tool's enforcement boundary. It can, however, be made visible and non-default: no official launch without attestation, and no unlabelled full-read fallback.

### 4. Fix-up budget

| Track | Maximum published Phase 8 fix-up cycles | Result at cap with unresolved work |
| --- | ---: | --- |
| `fast` | 1 | Blocking user escalation; no auto-close |
| `standard` | 2 | Blocking user escalation; no auto-close |
| `deliberation` | 6 | Blocking user escalation; no auto-close |

Six is a safety threshold, not a claim that later findings are trivial. It is deliberately just above the measured 5.1-round average while preventing the 19–24-round tail from continuing without a human reading the trajectory. Fresh late MAJORs are a reason to escalate, not a reason to accept the implementation.

One cycle is counted when the implementer publishes a new `## Fix-up cycle N` at a new HEAD. Retries that do not publish a cycle consume the separate driver budgets. Changing implementer, model, branch, track, or driver does not reset the count. A stricter-track upgrade preserves cycles already spent and applies the larger cap; if the spent count is already at that cap, escalation happens immediately.

Before starting cycle `cap + 1`, the driver/facilitator writes a durable blocking `to-user` escalation containing: cycle/round chronology and HEADs, finding counts by severity and whether each touched new or unchanged code, unresolved agreed fixes, repeated/reopened findings, validation status, elapsed/cost telemetry when available, and a concrete recommendation. Automation halts. The user may choose a revised plan, abandon/defer, issue a scoped ruling, or authorize a new explicit finite ceiling. Silence never extends the budget, waives a finding, or marks completion.

### 5. Content that must never be cut and detection

The irreducible safety kernel is: packet/source identity; idea status/phase/track/transport/quorum and routing flags; §4.0 overrides, invariants, and force-upgrade rule; the exact current phase block; non-solo and files-canonical close guards; §6 rule 3, status re-read, English-only, and no-secrets; the escalation mechanism; §14; the active transport's current-phase mechanics; and every close/cap/strict-gate condition applicable to the phase. §15 joins that kernel for Phases 1, 2, 3, 6, and 7. §7 joins it for every phase of a protocol-change idea.

Omission is detected structurally by the source/block hashes, applicability coverage, dependency closure, safety-kernel assertion, omission ledger, and attested launch record. Any failed detector selects the full protocol, never a partial packet. The packet also says that an action or claim outside its declared action/trigger set requires regeneration or a full read; the agent is not authorized to improvise from an apparently similar rule.

There is one limit that must be stated plainly: no generator can prove that a human-authored applicability annotation is semantically correct. A rule can be misclassified even when every hash and dependency check passes. Therefore applicability changes are themselves normative protocol changes under §7, the first release should run in shadow/audit mode with full-protocol comparison, and the conservative default is always-include. This design makes every intentional omission inspectable and every unclassified or mechanically inconsistent omission fail closed; it does not pretend to solve semantic completeness automatically.

## Concerns and open questions

- The protocol currently states that per-idea version pinning is not in force (`parley-deck/COOPERATION.md`, §7, **Not yet in force**). Until pinning ships, the packet envelope must bind a source hash and refuse to continue silently if the live protocol changes mid-idea. The group should decide whether regeneration under the new hash needs quorum acknowledgement for an in-flight idea.
- “Normative block” needs a deterministic structural definition. Keyword scanning for `MUST` is insufficient because tables, exceptions, and uncapitalized rules are normative too. Stable source annotations plus a conservative section-level fallback are safer than prose heuristics.
- The §15 on-demand trigger in Phases 5 and 8 must be tested against real implementation narratives. If agents routinely make completion verdicts there, §15 should simply become load-bearing in those packets too.
- A six-cycle deliberation cap should be reviewed after enough completed ideas to compare escalation rate, human decisions, and eventual defect closure. Changing the number should require evidence, but the first cap must not wait for perfect data while the current value is unbounded.

## Risks

- A wrong applicability annotation is the highest risk because it can produce a validly hashed but semantically incomplete packet. Conservative inclusion, §7 review, shadow rollout, and a visible omission ledger reduce but do not eliminate it.
- A packet may save reading time while its generation or prompt embedding adds tool latency. Measure total agent wall clock and prompt bytes, not generator runtime alone.
- Feature flags can change during implementation. If `strict_gate`, `auto_implement`, checks, track, or transport differs from the packet envelope, the packet is stale and must be regenerated before action.
- Explicit full-protocol fallback may become the habitual path and erase the benefit. Telemetry should report fallback rate and reasons, but must never convert fallback into a correctness failure.
- Renewable budgets can become “unbounded by repeated approval.” Requiring a user-authored explicit new ceiling and trajectory summary preserves human ownership and makes each extension auditable; it does not guarantee the user will refuse repeated extensions.
