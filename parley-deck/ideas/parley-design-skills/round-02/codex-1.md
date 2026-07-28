---
agent: codex-1
idea: parley-design-skills
round: 2
date: 2026-07-28
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

I changed my registry position: `rules.json` plus generated prose is not one source in practice.
It makes the dependency-free doctrine uneditable without the checker toolchain and leaves generated prose as a stale installed artifact; the late brief's M3 analysis exposed that mistake.
The source of truth should be one literate `RULES.md` whose fenced YAML records obey a frozen grammar and whose adjacent prose is the human doctrine.
The checker parses that file directly; it carries no fallback copy.

I changed the `parley-design` ceiling from five Markdown files / 96 KiB to exactly four files / 64 KiB.
The late artifact inventory demonstrated how quickly "helpful" references become a 17-file protocol inside a protocol.
Four participants paying at most 64 KiB each is bounded; four participants selecting different files from a 96 KiB tree is not.

I changed F4 from deterministic rank-k selection on every full run to deterministic, distinct primary-axis assignment on every full run.
A model can backfill a four-handle shortlist after learning its hash-selected rank, so the claimed rank roll is structurally unverifiable without another commit phase.
An enumerated primary axis assigned one-to-one across participants is directly checkable in the submitted round-01 frontmatter and attacks the actual failure: directions occupying the same design position.

I narrowed the artifact set.
`COMPONENTS.md`, a mandatory Resolver file, and a generated Markdown audit are not universal v1 artifacts.
Component specs are required only for components changed or introduced; Resolver is required only for declared modes; `AUDIT.json` is canonical and review prose cites it.

I did not change my central position: `parley-design` is a profile over Parley Deck, not a second workflow engine.
Mechanical facts gate, taste does not, and the human is the default Decider.
The pre-build contract outranks an unapproved implementation deviation.

## Responses to others

### @claude-1

I adopt Claude's insistence on a protocol-shaped specification, one critique round, a human Decider, bounded grafts re-expressed through winner tokens, and explicit degraded evidence.
Those choices answer the negative evidence rather than pretending debate improves taste.

I also adopt Claude's three-way distinction among objective defects, project-system conformance, and taste tells.
Collapsing contract violations into generic quality loses the prerequisite that a ratified project contract must exist before the rule has authority.

I reject Claude's vendored registry fallback in `parley-design-check`.
A fallback plus a drift test is still two installed copies between test runs; an independently installed checker must receive `--registry <RULES.md>` or fail configuration, never silently switch doctrine.

I reject the suggestion that the build wins when contract and implementation disagree.
Code is evidence of what happened, not authority to amend `FINAL.md`.
The implementation must be fixed or the Decider must approve a recorded deviation through ordinary review consensus.

I reject a distinct P0-P6 artifact ladder.
Its useful labels can be profile aliases, but its transitions must be Parley transitions or they are unenforceable ceremony.

### @codex-1

This is a self-audit, not a peer response.
I retain my round-01 gate semantics, DTCG adoption, audit pinning, waiver countersignature, and refusal to ship T3 pixel judgement in v1.

I reject my own `rules.json` source-of-truth proposal, five-file / 96 KiB ceiling, and seeded rank-k ritual for the reasons above.
I also retract mandatory primitive-to-semantic-to-component layering for projects that have no component-token layer; skipping a nonexistent layer is not a defect.
Reference direction must be acyclic and may skip optional layers only when the contract declares them absent.

### @hermes-1

I adopt Hermes's use of RFC 2119 language, stable rule ids, DTCG `2025.10`, WCAG 2.2 as the blocking accessibility baseline, APCA as advisory, and Named Rules as the citable prose unit.
These are interoperable constraints with falsifiable checks.

I reject Hermes's D0-D7 ladder and its large predeclared `design/` tree.
It duplicates Parley phase state, moves design-system ratification before reality is audited, and creates files whose values will restate tokens.

I reject mandatory OKLCH, mandatory light and dark modes, and mandatory three-tier token layering as universal conformance requirements.
DTCG permits several color spaces, a brief may require one mode, and small systems may have no component tokens.
Those are per-project contract choices, not anti-slop invariants.

I reject "the build wins and prose notes the divergence."
That is contract laundering by the implementer and is actively harmful.

I reject the claim that T0+T1+T2 can ship with zero runtime dependency while T1 uses a CSS parser and T2 an HTML parser.
Dependency cost must be stated literally: Node built-ins at T0, pinned parsers at T1, a browser at T2.

### @kimi-1

I adopt Kimi's strongest point: the ritual is a profile over existing Parley rounds.
I also adopt findings rather than holistic scores, exactly one critique round, leaky-anonymity honesty, no self-review, and the rule that grafts cannot replace a winner's token system.

I reject seeded divergence only after collapse.
By then every participant has already paid for a failed direction and the facilitator is judging a potentially gameable fingerprint.
Every full run should start with unique, deterministic positions on an enumerated primary axis; cross-model heterogeneity handles the remaining axes.

I reject automatically calling a second collapse `verified-genuine`.
Repeated convergence can be repeated attraction.
After one rerun, only the human Decider may accept the convergence with a brief-specific reason; otherwise the outcome is `ABSTAIN`.

I reject delaying the shared mechanical report until after each reviewer writes judgement.
For implementation review, one audit at the pinned commit must precede Phase 6 so reviewers consume identical facts and do not waste four model calls rediscovering token or contrast failures.
Direction critique remains unseeded by taste-detector output; artifact-shape and distinctness checks are different.

I reject `x-` token groups in place of DTCG `$extensions`.
Use the standard extension mechanism with reverse-domain keys; do not create a local dialect merely because agents are the readers.

## Resolved disagreements A-G

### A. Reuse Parley's state machine

Reuse Parley Deck's state machine.
`PDS/1.0` defines a profile, required sections, gates, and transition predicates; it does not define another durable phase cursor.

The exact mapping is:

- Phase 0 / `00-prompt.md`: qualify scope, enumerate divergence axes, name the Decider, declare evidence and route.
- Parley round-01: isolated directions.
- Between rounds: deterministic schema and distinctness gate.
- Parley round-02: one adversarial critique.
- `consensus.md`: Decider's one-winner verdict, bounded grafts, dissent, evidence declaration, and participant process signoffs.
- `FINAL.md`: immutable Direction Contract.
- Parley Phase 5: ordinary implementer applies the contract and derives the as-built system at the review commit.
- Parley Phases 6-8: pinned audit, dispositions, fixes, reconciliation, and completion.

A second D0-D9 machine loses Parley's ownership rules, track classifier, quorum, terminal states, Phase-5 implementer boundary, review severities, and driver enforcement.
It also creates impossible mixed states such as Parley `round-02` while PDS says `D7`.

Reusing Parley loses separate Heat Map, Scorecard, Rumble, and micro-phase commits.
That loss is intentional: Heat Map is optional critique data, numeric aesthetic scoring is theatre, and Rumble contradicts binding D3 and D8 by commissioning two prototypes before one winner is chosen.

### B. Registry source of truth

Choose one literate `parley-design/references/RULES.md`.
Each rule is an H3, one fenced YAML record, then rationale, counterexamples, and remedy prose.
The YAML is the only machine source and the prose in the same file is the only human source.

Freeze the extraction grammar: exactly one `pds-rule` YAML fence beneath each rule H3; UTF-8; duplicate ids fatal; unknown keys warn unless prefixed `x-`; unknown rule ids pass through as `UNJUDGEABLE`, never crash.
Schema validity, unique ids, permanent meanings, fixtures, and registry digest are CI gates.

`rules.json` plus generated Markdown breaks zero-dependency authorship and creates stale installed prose.
A Markdown table parser breaks on wrapped cells, embedded examples, and prose edits.
Two hand-maintained copies break immediately.
A vendored fallback breaks later and more quietly.

### C. Size and file count

`parley-design` has exactly four Markdown files and a hard 64 KiB total ceiling for `PDS/1.x`:

- `SKILL.md` <= 8 KiB.
- `references/PDS.md` <= 20 KiB.
- `references/RULES.md` <= 24 KiB.
- `references/WEB-ANNEX.md` <= 12 KiB.

The always-read set is `SKILL.md` plus `PDS.md`, at most 28 KiB per agent.
Critique loads `RULES.md`; web work loads `WEB-ANNEX.md`.
The facilitator pins the same phase read-set and content digests for every participant; agents do not choose subsets.

At four agents, the absolute worst doctrine load is 256 KiB, versus megabytes for the source projects.
Checker code is not model context and is outside this ceiling, but v1 must not bundle Chromium and its own `SKILL.md` stays <= 8 KiB.
Crossing either doctrine ceiling requires `PDS/2.0` plus measured retrieval-cost evidence.

### D. Evidence-tier vocabulary

One vocabulary wins everywhere:

- `T0 ARTIFACT`: Markdown/frontmatter, JSON, hashes, raw text, and token graphs using runtime built-ins.
- `T1 SOURCE`: parsed CSS/HTML/JS and static DOM facts; pinned parser dependencies, no computed layout.
- `T2 RENDERED`: pinned browser, computed styles, accessibility tree, focus behavior, geometry, and viewport sweeps.
- `T3 PIXEL`: screenshot/pixel evidence; not shipped in v1 and never a generic aesthetic score.

`UNJUDGEABLE` is a verdict, not a tier.
`PASS`, `VIOLATION`, `NEEDS_REVIEW`, and `RECOMMENDATION` are the other verdicts.
Human or agent judgement is provenance, not `T4`.

`stated|source|rendered|measured`, `text|static-dom|layout|pixel`, and the research strawman's six-tier list are rejected.
They make the same evidence incomparable and confuse measurement with pixel inspection.

### E. The dice

Every full run requires a deterministic one-to-one assignment on the brief's enumerated primary divergence axis.
`assignment = rotate(sorted(primary_positions), uint32(sha256("PDS/1" || run_id)[0:8]))` mapped to sorted participant ids.
The checker verifies the mapping and each direction's declared position.

This is not rank-k theatre.
The brief must provide at least as many materially distinct primary positions as proposers or the full route cannot start.
Cross-model heterogeneity remains valuable on all unassigned axes but is insufficient alone because the models share training attractors.

G1 additionally requires each pair to differ on at least two declared axes and forbids duplicate Signatures.
One failed set may rerun once.
Persistent collapse requires human acceptance with a brief-specific reason or `ABSTAIN`.

### F. Split and machine-readable contract

`parley-design` owns meanings: protocol profile, artifact schemas, rule ids, authority, thresholds, evidence minimums, waiver policy, conformance, and human fallback.
It owns no executable detector, browser launch, generated report, waiver mutation, or UI code.

`parley-design-check` owns execution: registry parsing, detector implementations, target adapters, capability discovery, fixtures, deterministic reports, exit codes, and waiver validation.
It owns no rule prose, threshold, authority, or aesthetic decision.

The interface is the `pds-rule` record:

```yaml
id: core:a11y-contrast-text
authority: quality
minimum_tier: T2
detector: contrast-text
parameters: {body_min: 4.5, large_min: 3.0}
default_verdict: violation
design_system_blind: true
override: human-decider
status: stable
since: PDS/1.0
```

Allowed authority values are `quality`, `contract`, and `slop`.
Checker code exposes detector names only; it reads parameters and authority from `RULES.md`.
Invocation is `design-check --registry <RULES.md> --config <config.json> --target <path> --json`.
If the registry is absent or malformed, exit 3; there is no embedded fallback.

`AUDIT.json` pins `spec`, registry digest, source commit, config digest, target profiles, requested/executed/unavailable tiers, tool versions, rule verdicts, and waiver ids.
Exit codes are 0 clean/recommendations, 1 blocking violations, 2 unresolved needs-review, 3 invalid config or required tier unavailable, and 4 internal failure.

Without the checker, agents read the same records and emit manual findings.
With the checker, reviewers consume one audit from one pinned environment.
Unknown newer ids are preserved as `UNJUDGEABLE`; older checkers do not invent a pass.

### G. Actively harmful proposals

Hermes's and the research strawman's "build wins" rule is actively harmful because it lets the implementer amend an immutable contract by violating it.
The build is descriptive evidence; only an explicit Decider-approved review disposition can authorize divergence.

The research strawman's D0-D9 state machine is actively harmful because it duplicates an engine the add-on is forbidden to amend.
Its Rumble branch violates the owner's one-winner-whole decision and drags design into Phase-5 prototype ownership.

The strawman's trimmed-mean Scorecard is actively harmful statistical theatre.
Four proposers yield too few independent non-self scores to justify trimming, and numeric aggregation launders judge bias into precision.

Hermes's mandatory OKLCH/dark-mode/component-layer policy is actively harmful as universal doctrine.
It replaces training-prior slop with a Parley house implementation.

Kimi's post-judgement checker ordering is actively harmful in code review.
It makes each agent rediscover deterministic defects and allows contradictory claims against different source states.

Claude's fallback registry copy is actively harmful despite a drift test.
Failing configuration is safer than silently enforcing yesterday's doctrine.

Any T0 regex that labels purple, cream, fonts, or marketing cadence a blocking defect is actively harmful.
Those are review signals with counterexamples, never automatic violations.

## Final positions on F1-F8

- **F1:** `quality` and ratified-`contract` violations may support a unilateral block with reproducible evidence; `slop` is advisory until normal review quorum makes it an agreed fix.
- **F2:** Convergence is an alarm; one rerun is allowed, then only a human reasoned acceptance or `ABSTAIN`, never automatic consensus.
- **F3:** Rendered evidence is optional-with-declared-capability; only a pinned T2 run may originate rendered verdicts, and absent required T2 blocks the claimed conformance level.
- **F4:** Every full run gets checkable, deterministic, distinct primary-axis assignments; cross-model heterogeneity alone and rank-k self-reporting are insufficient.
- **F5:** Exactly four doctrine files, 64 KiB total, with identical phase read-set digests for every participant.
- **F6:** Waivers are centralized JSON, narrow, expiring, reasoned, hash-pinned, and counter-signed; quality-floor rules remain design-system-blind.
- **F7:** Fast is allowed only for one surface within an existing system, <=3 existing components and <=5 files/300 LOC, with no new foundation, token family, or visual direction.
- **F8:** A human Decider selects; explicit delegation must name a non-proposer/non-critic, while Parley signoff ratifies process and objective obligations rather than voting on taste.

## New concerns / questions

The literate registry grammar must be deliberately boring.
If YAML anchors, multiline executable snippets, or prose-derived thresholds are allowed, parsers will diverge.
The schema should forbid them in v1.

Primary-axis positions can be nominally different but visually equivalent.
Each position therefore needs a brief-authored falsifier: one sentence stating what observable choice would violate that position.

`quality` is easy to overclaim.
Every stable blocking rule needs pass, fail, and at least three legitimate-counterexample fixtures; category changes require a registry major version and never reinterpret historical digests.

T2 determinism needs a render contract: browser/version, fonts, viewport, device scale, locale, motion preference, color scheme, network policy, and readiness signal.
Without it, "same pinned commit" is not the same evidence.

The Decider bottleneck is real and acceptable.
Calling an agent-selected winner "provisional" still encourages implementation before authority exists; unattended full runs should stop at `ABSTAIN`, not auto-ratify.

## Current proposal

Ship `PDS/1.0` as a four-file doctrine profile and `parley-design-check` as an independent Node CLI that requires the doctrine registry path.
The checker ships T0 and T1, supports an external pinned-browser adapter for T2, and reserves T3 without implementing it.

The only PDS additions to existing canonical Parley artifacts are structured sections:

- `00-prompt.md`: `design_profile`, `run_id`, `decider`, target profiles, enumerated axes, primary positions, evidence requirements, and fast/full route.
- `round-01/<agent>.md`: one capped Direction with assignment, axis fingerprint, Signature, risky move, refusals, token sketch, and evidence declaration.
- `round-02/<agent>.md`: one non-self adversarial critique with rule/goal references, evidence tier, verdict, and graft candidates; no score.
- `consensus.md`: exactly one winner, 0-3 grafts, dissent, collapse disposition, evidence degradation, and ordinary Parley signoffs.
- `FINAL.md`: the immutable Direction Contract with winner hash, Named Rules, must-share/may-vary, refused defaults, token commitments, grafts, and finish conditions.

At the Phase-5 review commit, the ordinary implementer supplies `design-system/DESIGN-SYSTEM.md`, DTCG token files actually used, scoped component specs for changed components, optional Resolver data for declared modes, optional `waivers.json`, and canonical `AUDIT.json`.
`IMPLEMENTATION.md` carries the reconciliation table so it remains inside Parley's living execution record.

Reconciliation dispositions are `match`, `acceptable-adaptation`, `missing`, `contradicted`, and `added-without-approval`.
The last three are MAJOR until fixed or explicitly accepted by the Decider through review consensus.
An accepted deviation is documented in the as-built system and waiver ledger but never rewrites or retroactively satisfies the immutable contract.

Conformance is cumulative without a second state machine:

- `PDS-L1 PROFILE`: required profile sections, isolation, assignment, distinctness, critique, and one-winner verdict are valid.
- `PDS-L2 CONTRACT`: `FINAL.md` carries a valid hash-pinned Direction Contract.
- `PDS-L3 SOURCE`: the as-built system reconciles, DTCG/source checks ran, and no quality or contract violation remains.
- `PDS-L4 RENDERED`: the pinned T2 contract ran for every declared target and no rendered quality or contract violation remains.

Projects claim the highest level actually evidenced.
`UNJUDGEABLE` is honest lower-level conformance, not permission to claim L4.
This design keeps the portable doctrine useful alone, gives the checker a strict interface, and refuses both taste-by-committee and implementation-by-violation.
