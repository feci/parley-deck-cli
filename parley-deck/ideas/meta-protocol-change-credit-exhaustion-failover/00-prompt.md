---
idea: meta-protocol-change-credit-exhaustion-failover
author: user
created: 2026-10-29
participants: [codex-1, kimi-1, zcode-1]
implementer: codex-1
track: deliberation
auto_implement: true
require_model_diversity: true
checks: go test ./... -count=1 -timeout 45m
status: consensus
---

## Problem / idea

Exclude an agent at once when its provider reports exhausted account credits, tokens, quota, or usage allowance. Do not retry or reconnect that agent in this idea. Let the next eligible agent in order take over its role.

Read [the controlling brief](source-context/ORGANIZER-BRIEF.md) in full. It defines the design, process, close authority, and release contract.
The owner authorizes design, implementation, review, release, and local installation without a proposal stop.

## Owner direction

The owner wrote:

> "podme navrhnut update protokolu, ktory vynecha automaticky agenta, ktory ma model ktoremu providerovi dosli kredity, bez toho aby ho skusal znova pripojit, ak nejakemu agentovi dojdu tokeny alebo kredity, nech je hned vyradeny a dalsi v poradi preberie jeho ulohu, nastartuj /parley-deck a sprav novu verziu, pricom na zaver to aj vsade pushni a deployni cez vsetky kanaly a nainstaluj lokalne do vsetkych agentov"

The brief supplies the English translation. Preserve the quotation as source evidence under Appendix B.

## User direction

## Owner ruling: the succession order (binding)

The owner gave this ruling at about 19:30 CET on 2026-10-29, before round 1. The relay recorded it verbatim:

> "no ked vypadne organizer, tak implementer ho preberie, ak vypadne implementer tak organizer ho preberie,
> ak vypadne aj organizer aj implementer, tak ho preberie prvy participant a participanti ked vypadnu, tak sa
> nedej nic, len proste vypadnu"

Translation: "Well, when the organizer drops out, the implementer takes it over. If the implementer drops out,
the organizer takes it over. If both the organizer and the implementer drop out, the first participant takes
it over. When participants drop out, nothing happens; they just drop out."

## What this settles

1. The organizer drops out: the implementer takes over the organizer role.
2. The implementer drops out: the organizer takes over the implementer role.
3. Both drop out: the first participant takes over both roles. "First" is the first remaining agent in
   `participants:` order that is neither the dropped organizer nor the dropped implementer.
4. A participant (not the organizer, not the implementer) drops out: no succession. It only drops out.

This replaces the "recommended default order" in brief section 3. FINAL must implement this ruling. The
participants decide only the details that the ruling does not settle. Examples:

- one agent that is both organizer and implementer (as in this idea) counts as "both drop out";
- the rule "a declared facilitator does not implement" yields to rule 2;
- the drafter roles and the reviewer count after a succession;
- the floor and the single-reviewer exception after a succession.

## Constraints

- Codex organizes and implements. Kimi and Zcode provide independent analysis, reviews, and fresh goal checks.
- Use the existing github-pr transport and the owner-allocated credit-failover branches.
- Release CLI 1.56.0, skill 2.20.0, and stage core 2.20.0 from staged core 2.19.0 plus reviewed hunks.
- Retain the two-person usable non-facilitator floor, retained dissent, independent review, and current-tree acceptance evidence.
- Keep transient failures on the shipped two-attempt path. Decide explicit provider scope, reset semantics, and role succession in FINAL.
- Reuse the shipped exclusion reducer, immutable history, retry ledger, and cause-derived review gate where possible.
- Keep the 20,000-byte skill and 70,000-byte packet limits. Compact without rule loss. Do not worsen the meta-idea variant.
- Use ASD-STE100 through Appendix B for all new prose. Preserve evidence, code, and required syntax.
- Maximum five fix-up cycles. The controlling brief defines retry and close authority.
- Never prune or declare worktrees. Never change credentials or gateway settings.
- npm publication and global core publication remain owner-only. Do not attempt either gate.

## Non-goals

Generic retry services, account management, unrelated roster changes, Windows promotion, and changes to closed prior artifacts are outside scope.
Include only required parts of the two quota follow-ups. State remaining deferred parts.

## Readiness and launch context

The three selected agents returned PONG in fresh measured processes. See source-context/preflight-selected.json.
Codex uses gpt-6-astra/xhigh. Kimi uses omniroute/k3.1-pro/max. Zcode uses configured zai/glm-5.3/max.
These model labels describe configuration. They do not prove the provider's observed model.

The standalone presence preflight reports a historical facilitator declaration defect in another idea. It made no membership change.
The run command lacks a fixed-slug option. The organizer records both driver gaps and uses measured launch fallback where required.
The adapter defaults disagree with the machine roster. A task-local launch file applies the approved roster models without changing persistent configuration.
The first stale Kimi adapter probe failed authentication. The correct roster model returned PONG. No credential was changed.

Protocol context is full. Source and packet SHA256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec.
There is no fallback reason. The complete live authority and Appendix B were read.

## Sources

- source-context/ORGANIZER-BRIEF.md
- ../meta-protocol-change-quota-auto-exclude/FINAL.md and IMPLEMENTATION.md
- ../meta-protocol-change-participant-dropout/FINAL.md and IMPLEMENTATION.md
- ../meta-protocol-change-participant-dropout-review-gate-timing/FINAL.md and IMPLEMENTATION.md
- ../meta-protocol-change-driver-unstall/FINAL.md and IMPLEMENTATION.md
- ../quota-kickoff-reporting-and-alias-guidance/00-prompt.md
- ../quota-zcode-native-exhaustion-capture/00-prompt.md

## Round-2 source evidence

Read source-context/codex-native-source.md. These source excerpts arrived after codex-1 filed round 1.
They describe a candidate native JSON channel, not an accepted recognizer or a captured exhaustion event.
