---
agent: codex-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

## Summary

The diagnosis is directionally credible but causally unproved: protocol growth is real, yet bytes read are not a wall-clock measurement and “read cost rather than execution time” is a false separation because model context processing is part of execution. My competing explanation is critical-path drift: deeper/slower models, participant changes, queueing, retries, and the slowest member of a parallel round may dominate perceived latency, with protocol growth contributing mainly on cold or uncached calls. I recommend measuring that distinction first, then replacing monolithic reads with a deterministic phase-and-role context compiler; physical reordering and free-form facilitator summaries are insufficient. I deliberately did not read all of `COOPERATION.md`: I scanned its headings and read Quickstart, the three requested phases, their direct dependencies, and the four named safety areas; the cost is residual uncertainty about unmarked cross-section dependencies, which the proposed dependency-checked compiler must make detectable.

## Proposed approach

### 1. Treat the diagnosis as a hypothesis, not a result

I independently reproduced the size trend at its endpoints:

```text
PRIMARY — quoted command output
$ git show 697cb66:internal/protocol/defaults/COOPERATION.md | wc -l -c
     721   49380
$ wc -l -c internal/protocol/defaults/COOPERATION.md
    1360  104480 internal/protocol/defaults/COOPERATION.md
```

That is 2.12x byte growth. It establishes exposure, not latency. The current round-one runner reads `00-prompt.md` and embeds it in the launch prompt; it does not embed `COOPERATION.md` (`internal/runner/runner.go:821`, `internal/runner/runner.go:849`). Separately, the installed Parley skill orders agents to load the full protocol:

```text
PRIMARY — quoted command output
$ rg -n 'Always read .*COOPERATION\.md.*first|Load the full cooperation protocol before acting' <installed-skill>/SKILL.md
12:Always read `parley-deck/COOPERATION.md` first. ...
24:Do not run this skill from the abbreviated workflow alone. Load the full cooperation protocol before acting:
```

This means protocol cost is currently imposed by agent instructions rather than by the round-one prompt builder. Optimizing only the Go prompt builder would not remove it. The installed skill is itself material context—`wc -l -c SKILL.md` in its installed directory returned `947 58754`—so measuring only repository file reads may undercount duplicated instruction context.

The competing critical-path explanation has a concrete repository signal: Kimi became the active fourth participant on 2026-07-18, and the configuration calls it “the slowest participant (~20 min/round)” while assigning a 40-minute timeout (`parley-deck/agents.toml:16`, `parley-deck/agents.toml:36`). That comment is not causal proof, but it is a confounder strong enough that the present diagnosis cannot dismiss execution/backend time from CLI microbenchmarks. `parley --help` timing says nothing about hosted model latency.

Run a randomized replay benchmark over a frozen idea and frozen outputs:

1. Hold agent, exact model, effort, output-token cap, artifact set, and concurrency constant.
2. Factor A: full protocol versus generated phase packet. Factor B: cold/uncached versus warm/prefix-cached input.
3. Repeat each cell enough times to expose provider variance and randomize order.
4. Capture provider-reported input, cached-input and uncached-input tokens; time to first token; generation time; retries/rate limits; per-agent duration; and round wall time (`max(agent durations)`, not their sum).
5. Replay once with the historical roster/profile and once with the current roster/profile if exact configurations can be reconstructed.

If scoped context materially lowers time to first token and round wall time under the same backend, the read-cost diagnosis wins. If the difference is small, disappears under warm caching, or is dwarfed by one participant, critical-path drift wins. My current position is that both contribute, but the claim that read cost is *the* cause is unverified.

### 2. Load a small invariant kernel plus the current phase card

The always-loaded kernel should contain identity and exact output path, active idea/track, canonical-file ownership, no-secrets and English-only rules, “never edit another agent's file,” provenance routing, and two conditional guards: attempting a protocol change loads §7; automated/standing operation loads §14. Everything else should be selected by phase, role, track, and idea flags.

| Phase | Load-bearing context | Reference opened only when triggered |
| --- | --- | --- |
| Phase 1 participant | `00-prompt.md`; the Phase 1 schema and independence rule (`parley-deck/COOPERATION.md:306`, `parley-deck/COOPERATION.md:324`); active-track invariants; §6 ownership/language rules; §15.1–15.4 for claims, provenance, conflicts, and exemption claims (`parley-deck/COOPERATION.md:1236`, `parley-deck/COOPERATION.md:1319`). No peer round-one file. | Roster/readiness, transport publication, later phases, pipelines, retrospectives, and other tracks. §7 only if the action would change protocol text; §14 only for an automated outer loop. |
| Phase 3 drafter or signer | Phase 3 schema/signoff rules (`parley-deck/COOPERATION.md:354`, `parley-deck/COOPERATION.md:394`); active-track quorum (`parley-deck/COOPERATION.md:719`); latest participant positions, unresolved objections, and their evidence; all applicable §15 duties, especially verdict conflicts, role concentration, and correlated agreement (`parley-deck/COOPERATION.md:1288`, `parley-deck/COOPERATION.md:1357`). A signer needs the proposed consensus plus evidence links and its own prior position; a drafter needs the full current-state set. | Historical superseded prose unless a position changed, a locator is challenged, or an unresolved item points to it; implementation/review mechanics; inactive transports. |
| Phase 6 reviewer | `FINAL.md`, `IMPLEMENTATION.md`, reviewed commit/diff and checks; Phase 6 schema, severity vocabulary, refutation-default, and no-suppression rule (`parley-deck/COOPERATION.md:501`, `parley-deck/COOPERATION.md:556`); relevant §15.1–15.4; own output contract. | Design rounds unless checking a deviation or prior disposition; consensus/fix-up mechanics until the review advances; pipeline and strict-gate material unless enabled. |

A deliberately conservative raw-text slice—without rewriting or compacting any selected rule—was already much smaller than the full body:

```text
PRIMARY — quoted command outputs from selected COOPERATION.md line ranges
Phase 1: 180 lines, 14583 bytes
Phase 3: 275 lines, 20472 bytes
Phase 6: 224 lines, 17909 bytes
```

Those are about 14%, 20%, and 17% of the current 104,480-byte default. They are a feasibility probe, not a safe production manifest: dependency closure has not yet been mechanically proved.

### 3. Cut repetition and routing detail, never safety semantics

Cut from default participant packets:

- inactive transport branches, bootstrap/roster mechanics, CLI invocation recipes, ratification history, and not-yet-in-force notes;
- phase templates and lifecycle prose unrelated to the current role;
- repeated rationale and examples once a stable normative rule ID and exact source locator exist;
- other tracks except the active row plus universal fail-closed invariants.

Do not delete the full authority or weaken §15, §7, §6 rule 3, or §14. “Never cut” should mean two things: the exact normative rule remains in the authoritative protocol, and it is never silently omitted when its applicability predicate is true. §6 rule 3 belongs in every participant packet. §15.1–15.4 belong in analysis and review packets; §15.5–15.6 join consensus packets. A one-line kernel guard must load exact §7 before any protocol mutation and exact §14 before any scheduled/standing loop acts. Also retain non-solo execution, round-one independence, canonical audit ownership, refutation-default/no-suppression review, append-only signoffs, no-secrets, and track fail-closed routing.

The already-ratified physical restructure is not the optimization: its final specification explicitly preserves every line and only moves §9 (`parley-deck/ideas/protocol-restructure-appendices/FINAL.md:11`, `parley-deck/ideas/protocol-restructure-appendices/FINAL.md:17`), and says a ≤200-line core was not delivered (`parley-deck/ideas/protocol-restructure-appendices/FINAL.md:50`). It improves navigation but cannot reduce input when agents are still told to read the whole file.

### 4. Build a deterministic context compiler

Recommend a phase-scoped generated view backed by a typed rule registry, not manual prompt-side scoping and not a hard token truncator.

Each normative rule should have a stable ID, exact source text/locator, `phase`, `role`, `track`, condition predicates, and explicit dependencies. Given protocol version/hash plus idea frontmatter and requested action, a command such as `parley protocol context --idea <slug> --phase round-01 --role participant` should emit:

- the immutable kernel first, followed by deterministic phase cards to preserve prefix-cache value;
- a manifest containing protocol/idea hashes, included rule IDs, source locators, and the evaluated predicates;
- task artifacts required by that phase, with full authority available on demand.

Generation must fail closed on any unclassified normative rule, dangling dependency, stale hash, unknown phase/role, or applicable safety rule omitted from the manifest. Golden tests should assert expected rule IDs per phase; mutation tests should flip each condition and prove §7, §14, strict-gate, pipeline, and §15 cards appear. Per-phase byte/token budgets should be CI regression alarms with reviewed waivers, never runtime truncation. The installed skill and participant prompt must then direct agents to the generated packet and stop mandating the monolith; otherwise the new view is additive overhead.

Costs are real: rules must be classified and maintained, cross-rule dependencies can be wrong, generated packets may hide useful background, and changing prefix shape can reduce provider caching. Stable ordering, source hashes, on-demand retrieval, and fail-closed coverage tests are the price of safe savings. This mechanism is nevertheless safer than trusting each facilitator to remember which paragraphs matter.

### 5. Separate prior-round compounding from protocol loading

The compounding is both separable and currently an implementation choice. `gatherPriorRounds` loops over rounds `1..N-1`, reads every Markdown participant artifact, and concatenates it (`internal/runner/runner.go:936`, `internal/runner/runner.go:965`); the generated prompt then orders the agent to read every prior artifact (`internal/runner/runner.go:989`). Phase 2 itself requires addressing every active peer and providing counter-proposals, but does not say every historical version of every position must be reread (`parley-deck/COOPERATION.md:347`, `parley-deck/COOPERATION.md:352`). Consensus drafting separately orders another full-history read (`internal/app/driver_consensus.go:110`, `internal/app/driver_consensus.go:128`).

Round 2 should still receive every peer's round-one file in full: it is the first non-anchored contact with their analysis. For round 3 and later, load the immediately previous round in full plus a structured carry-forward ledger containing every unresolved objection, reservation/block, material claim and provenance locator, position change, and exact source path/hash. Older resolved or explicitly superseded prose stays addressable on demand and is pulled in verbatim whenever a conflict, changed position, missing locator, or exemption claim refers to it. This changes repeated history from approximately quadratic growth toward “latest state plus unresolved history.”

A free-form digest is not enough. The existing digest takes only the first Summary sentence/paragraph, caps it at 120 characters, and explicitly treats keyword flags as hints rather than verdicts (`internal/driver/digest.go:10`, `internal/driver/digest.go:36`, `internal/driver/digest.go:74`, `internal/driver/digest.go:112`). Using it would erase arguments. Require each participant to maintain an owned, machine-validated carry-forward block with stable claim/objection IDs and `supersedes` links; the compiler selects from those blocks but never invents their meaning.

Without full peer text, the failure modes are omitted minority arguments, lost qualifiers, false supersession, broken §15 ownership/provenance, and inaccurate quotations. The mitigations are participant-owned state, verbatim inclusion of unresolved/conflicting evidence, hashes and locators, mechanical completeness checks against active participants, and one-command source expansion. If a carry-forward block is missing or invalid, fail back to the full source rather than silently summarize it.

## Concerns and open questions

- Provider token accounting and cache telemetry may not be exposed consistently. The benchmark needs an explicit “unknown” state rather than deriving tokens from bytes and pretending it measured model work.
- What does “slower” mean to the owner: time to first artifact, whole-round wall clock, entire idea completion, cost, or subjective response quality? The mechanisms differ, and the current report does not define the dependent variable.
- Can historical model/effort/roster settings be reconstructed exactly enough for a fair replay? If not, compare full versus scoped context on today's fixed roster and keep the historical-cause claim open.
- Who owns carry-forward correctness? It should be the participant whose position is summarized, not the facilitator. Mechanical validation can prove presence and locators, not semantic fidelity.
- The selective protocol read cost me confidence that no unrelated section contains an unmarked dependency. I mitigated this by scanning every heading and following direct cross-references, but a safe compiler needs a machine-checkable dependency graph rather than that human judgment.
- This recommendation would eventually require the §7 protocol-change path before changing normative text or global skill behavior. This round should diagnose and specify the experiment, not apply it.

## Risks

- A selector bug can remove a brake more dangerously than a large prompt wastes tokens. Fail-closed generation, complete rule classification, and source fallback are release gates, not optional polish.
- Smaller context can improve latency while worsening decisions through missing background; measure defect/refutation yield and signoff reversals alongside time and token metrics.
- Dynamic packets may defeat prefix caching. Keep the invariant kernel and phase cards byte-stable and ordered; put volatile task material after them.
- Token budgets invite semantic compression and gaming. Use them only as regression alarms, never as permission to truncate a safety rule.
- A centralized digest can amplify facilitator framing—the same conflict this round is testing. Participant-owned state plus verbatim unresolved evidence is essential.
- Changing cumulative-history behavior without updating consensus drafting would move, not eliminate, the quadratic cost.
