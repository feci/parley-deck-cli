---
agent: codex-1
idea: parley-design-skills
round: 1
date: 2026-07-28
---

## Summary

I propose `PDS/1.0.0` for the dependency-free design protocol and `PDC/1.0.0` for its optional reference checker.
`parley-design` should reuse Parley Deck's files and phases rather than create a second state machine: round 1 carries isolated directions, round 2 carries one adversarial critique, consensus records the Decider's choice, and `FINAL.md` carries the binding direction contract.
After Phase 5, the ordinary implementer documents the system actually present in the code as `design-system/DESIGN-SYSTEM.md` plus DTCG `2025.10` tokens; this descriptive system cannot retroactively legalize a contract violation.
`parley-design-check` should gate only reproducible artifact, token, source, accessibility, and geometry facts. Taste signals are recommendations, never blockers.
The multi-agent benefit is real but narrow: heterogeneous models generate a broader candidate set than one model's measured 30/35 argmax repetition. They are not an accuracy or taste amplifier.
That advantage is squandered by exposing directions before submission, feeding detector results before independent critique, majority voting, repeated debate, or merging the winners' token systems.
The negative evidence therefore changes the design: one critique round, a mechanical distinctness gate, a human Decider, preserved `ABSTAIN`, and no numeric aesthetic gate.

## Proposed approach

### 1. Protocol shape of `parley-design`

The add-on should contain exactly five Markdown files, at most 96 KiB total:

- `parley-deck-skill/addons/parley-design/SKILL.md` — entrypoint, routing, required read set, and fast-path classifier; at most 12 KiB.
- `reference/protocol.md` — scope, RFC 2119 terms, roles, phases, transition table, decision rules, and conformance; at most 24 KiB.
- `reference/artifacts.md` — the schemas and complete minimal examples; at most 18 KiB.
- `reference/rules.md` — surface-neutral Named Rules and machine-readable check links; at most 24 KiB.
- `reference/web-annex.md` — CSS, WCAG 2.2, viewport, and web evidence rules; at most 18 KiB.

`SKILL.md` and `protocol.md` are the always-read core. `artifacts.md` is required for D0-D4, `rules.md` for critique/audit, and `web-annex.md` only when `surface_profile: web`; every participant in a phase reads the same declared set.

`PDS/1.0.0` has these normative phases:

1. `D0 QUALIFY` — extend `00-prompt.md` with `design.spec`, `surface_profile`, `scope`, `decider`, `run_id`, `divergence_axes`, `required_evidence`, and the fast/full route.
2. `D1 DIVERGE` — each `round-01/<agent-id>.md` contains one complete direction, written in isolation.
3. `D2 DISTINCTNESS` — the driver/checker evaluates categorical fingerprints before peers may read the directions.
4. `D3 CRITIQUE` — one `round-02/<agent-id>.md` per participant; every participant critiques every direction except its own and emits `preserve`, `defect`, and `graft-candidate` records.
5. `D4 DECIDE+GRAFT` — the Decider chooses exactly one direction whole; agents sign that the process and objective constraints were followed, not that they share the Decider's taste.
6. `D5 APPLY` — ordinary Parley Phase-5 implementation; the add-on never owns UI code.
7. `D6 DOCUMENT` — derive the durable system from the implemented commit.
8. `D7 AUDIT` — run the available checker tiers once at a pinned commit, then use normal Parley Phases 6-8 for dispositions and fixes.

The round-1 direction extension fields are `artifact_kind: direction`, `spec: PDS/1.0.0`, `direction_id`, `run_id`, `seed_key`, `assigned_rank`, `fingerprint`, `evidence_capabilities`, and `content_sha256`.
Its fixed body is capped at 750 words plus token tables: `Thesis`, `Refused default`, `Audience/job`, `Macrostructure`, `Type voice`, `Colour strategy`, `Density`, `Motion posture`, `First viewport`, `Risky move`, `Candidate shortlist`, and `Critique requests`.

Conformance is cumulative:

- `PDS-C1` — valid D0-D4 artifacts, isolation attested, distinctness passed, one winner, and bounded grafts.
- `PDS-C2` — a valid direction contract is embedded in `FINAL.md` and pins all source hashes.
- `PDS-C3` — the post-build design-system artifacts resolve and reconcile against the contract.
- `PDS-C4` — an `AUDIT.json` pins commit, registry hash, engine coverage, and has no unresolved blocking finding or required review.

The doctrine can self-assess C1-C3 without tooling; C4 may only be claimed from `parley-design-check`.
Rule ids are permanent and never reused. Changing a rule's category, blocking authority, or meaning is a `PDS` major version; threshold tightening is at least a minor version.
Deprecated rules remain recognized for one full minor version, carry `deprecated_since`, `replacement`, and `remove_in`, and historical audits always use their pinned registry hash.

### 2. The physical design-system artifact set

The pre-build artifact is `FINAL.md` section `## Direction contract`, with:

- `contract_id`, `spec`, `winner_direction`, `winner_sha256`, `decider`, `decision_ref`, and `decided_at`.
- Five falsifiable Named Rules covering structure, type, colour/material, density/spacing, and interaction/motion.
- `must_share`, `may_vary`, `refused_defaults`, `grafts[]`, `evidence_requirements`, and `finish_condition`.
- Each graft names `source_direction`, `source_part`, `winner_rule_preserved`, and the winner token through which it will be expressed.

The post-build set is:

- `design-system/DESIGN-SYSTEM.md` — frontmatter plus Named Rules, foundations, content voice, shared/varying boundaries, component index, and reconciliation table.
- `design-system/tokens.tokens.json` — canonical DTCG `2025.10` snapshot with nested `primitive`, `semantic`, and `component` groups.
- `design-system/resolver.json` — DTCG Resolver `2025.10`, required only when the implementation has modes.
- `design-system/COMPONENTS.md` — each component has Purpose, When to use, When not to use, Anatomy, Variants, States, Behaviours, Interactions, Content, Accessibility, and Tokens used.
- `design-system/design-check.waivers.json` — present only when at least one waiver exists.
- `design-system/AUDIT.json` — stable machine report; `AUDIT.md` is generated only for human review and is not a second source of truth.

`DESIGN-SYSTEM.md` frontmatter includes `artifact_kind`, `spec`, `source_commit`, `contract_ref`, `contract_sha256`, `token_files`, `surface_profiles`, `documented_by`, and `documented_at`.
Use DTCG exactly for token identity, `$value`, `$type`, `$description`, aliases, `$deprecated`, and Resolver modes.
Raw values live only in `primitive`; `semantic` values alias primitives; `component` values alias semantic tokens. The checker rejects reversed or skipped reference edges.
Use DTCG `$extensions` only with reverse-domain keys such as `org.parley.design.provenance`; do not add Parley-only top-level token fields.
Deliberately keep Named Rules, component behaviour, content guidance, provenance, and governance in Markdown: DTCG models values, not those concerns.
Generated CSS variables, Tailwind config, Swift values, or terminal palettes are adapters and never canonical token sources.

Reconciliation is mandatory before Phase 6 opens. Every implemented difference is classified as `match`, `acceptable-adaptation`, `missing`, `contradicted`, or `added-without-approval`.
`missing`, `contradicted`, and `added-without-approval` are MAJOR until the implementation is fixed or a Parley review consensus plus the Decider explicitly accepts the deviation.
Because `FINAL.md` is immutable, an accepted deviation is cited from `IMPLEMENTATION.md` and review consensus; only then may the post-build system document it.
The documenter MUST NOT canonize a value merely to make a finding disappear, and a quality floor remains design-system-blind.

### 3. The ritual made mechanical

Each full-route participant first writes four one-line candidate handles in resonance order.
The selected rank is `1 + (uint32(sha256("PDS/1" | idea_slug | run_id | agent_id)[0:8]) mod 4)`; the participant commits that rank rather than its rank-1 default.
The artifact records all four handles, their hashes, the seed key, and the assigned rank, making the local roll reproducible without a service.

The five fingerprint axes are `macrostructure`, `type_voice`, `colour_strategy`, `density`, and `motion_posture`, each from the enums in `artifacts.md`.
`D2` fails when a pair matches all of the first three axes, or when at least three directions match on four of five axes.
Only the collapsed participants rerun, once, with driver-assigned positions on the primary differing axis.
If collapse remains, the Decider must either accept the convergence as intentional with a written brief-specific reason or return `ABSTAIN`; it is never treated as automatic consensus.

Round 2 is a silent, simultaneous critique with absolute brief-linked records, not pairwise ranking.
No participant scores its own direction. Each `defect` cites a brief goal, contract/rule id, evidence tier, and 0-4 impact anchor; proposed fixes are optional and non-binding.
Artifact presentation order is deterministically rotated per critic from the run seed.
One critique round is the default. A second requires a Decider-recorded factual ambiguity; it is not allowed merely because tastes differ.

The human is the default and only automatic Decider. Scores, heat marks, and agent preferences are advisory.
If the human explicitly delegates, the delegate is named in D0 and excluded from proposing and scoring.
If no eligible Decider can choose confidently, `ABSTAIN` pauses D4; Parley consensus cannot manufacture a winner by vote.

The Decider selects one winner whole and may choose no grafts or 2-3 grafts.
Grafts must come from round-2 `graft-candidate` records, must be discrete component/copy/interaction/motion details, and must be re-expressed through winner tokens.
Colour systems, type scales, spacing scales, grids, and overall component languages are never graftable.
Losers become `closed: not-selected`; later re-litigation requires a new idea or an explicit contract-invalidating finding.

I reject three parts of the DCP/1 strawman.
First, its Rumble branch contradicts binding D3 and would pull the add-on into prototype implementation; an incommensurable pair should yield `ABSTAIN/NEEDS-EVIDENCE` and a separately scoped experiment idea, after which a new design run still selects one whole.
Second, trimmed-mean aesthetic scoring over three peer scores is statistical theatre given 38.34% human agreement; retain anchored criterion vectors as evidence, but do not aggregate them into a gate.
Third, DCP/1 creates a parallel artifact state machine and ratifies `DESIGN.md` too early; reuse Parley rounds and write the durable system after Phase 5.

### 4. Slop doctrine and load model

The always-loaded doctrine should contain no themes, archetype catalogue, or house palette.
It should contain twelve short surface-neutral Named Rules covering hierarchy, complete interaction states, truthful content, legibility, focus, reduced motion, bounded effects, semantic token use, evidence honesty, responsive survival, one coherent visual world, and category-plus-avoidance predictability.
These are concrete obligations: each names required evidence, allowed override authority, and the check id or `agent-judgement`.

The web annex carries numbers and syntax: WCAG 2.2 text contrast 4.5:1 or 3:1 for large text, non-text contrast 3:1, 24x24 CSS-pixel minimum-target review, reflow at 320 CSS px, the 320/375/414/768/1280/1920 viewport sweep, no `transition: all`, and no unhandled motion without `prefers-reduced-motion`.
It also carries the overused-font list as a `slop` prior, not a ban, plus measured thresholds such as 15% primary-font share over at least 20 text elements.
Effect budget is one dominant material/effect family and one authored motion moment per surface unless the direction contract names and justifies more.
Free axes cannot remain implicit: every direction must explicitly choose all five fingerprint axes and state the category default it refuses.
This prevents the absence of a shipped theme from silently handing every free axis back to each model's training prior.

The 96 KiB/five-file ceiling is binding for `PDS/1.x`; exceeding it requires a major version and a measured load-cost justification.
New defects do not automatically become doctrine. A proposed rule begins `draft`, needs examples and counterexamples, and enters the stable core only through a normal Parley idea.

### 5. Split line and machine-readable contract

`parley-design` owns meanings, protocol phases, roles, artifact schemas, Named Rules, surface-neutral invariants, evidence requirements, the web annex, and manual degraded behavior.
It owns no executable detector, package install, browser launch, UI code, waiver mutation, or numeric taste score.

`parley-design-check` owns registry data, executable predicates, target adapters, fixtures, CLI/report schemas, tier capability discovery, deduplication, and waiver validation.
It MUST be deterministic, offline after installation, non-interactive, and runnable independently of any agent runtime.

Every doctrine rule has a fenced YAML contract record:

`id`, `authority: quality|contract|slop`, `check_id|null`, `minimum_tier`, `evidence_required`, `design_system_blind`, `since`, and `override`.

The checker registry repeats these bridge fields plus `doctrine_ref` and `doctrine_sha256`.
`bin/verify-contract.mjs` fails CI if an implemented check lacks a doctrine id, if authority differs, or if the doctrine hash is stale.
The doctrine remains useful without the checker: unavailable checks are recorded `UNJUDGEABLE`, never silently passed.

### 6. `parley-design-check`

Proposed files are `SKILL.md`, `bin/check.mjs`, `registry/rules.json`, `registry/rules.schema.json`, `schema/config.schema.json`, `schema/audit.schema.json`, generated `reference/rules.md`, and `fixtures/<rule-id>/{pass,fail,counterexample}.*`.
`registry/rules.json` is data only. Each record has:

`id`, `title`, `authority`, `status`, `default_outcome`, `parley_severity`, `tier`, `targets`, `scopes`, `detector`, `parameters`, `evidence_required`, `design_system_blind`, `yields_to`, `since`, `deprecated`, `replacement`, `doctrine_ref`, and fixture paths.

Detection functions live outside the registry and dispatch by `detector`; thresholds remain data.
Only `status: stable` plus `authority: quality|contract` may block.
A stable blocker requires deterministic pass/fail fixtures, at least three legitimate counterexamples, and no unresolved false-positive fixture.

Engine tiers and dependency cost are explicit:

- `T0 ARTIFACT/TEXT` — Node built-ins only; protocol shape, hashes, JSON, line scans, raw literals, `transition: all`, and seeded distinctness.
- `T1 STATIC-AST` — optional pinned PostCSS/HTML parser and colour library; token graph, CSS property/value membership, ΔE, selector and DOM facts; no computed layout.
- `T2 RENDERED` — optional existing Playwright/Chromium plus a declared URL/file target; computed styles, axe, target geometry, focus, overflow, and viewport sweeps.
- `T3 PIXEL` — not shipped in v1; any future pixel check needs a pinned render environment and may not become a greenfield quality score.

No tier silently falls back. `AUDIT.json` records `requested_tiers`, `executed_tiers`, `unavailable_tiers`, tool versions, viewport set, commit, registry version/hash, and every `UNJUDGEABLE` rule.

Initial stable ids and thresholds should include:

- `pdc:artifact-schema`, `pdc:content-hash`, `pdc:exactly-one-winner`, and `pdc:graft-budget` at T0.
- `pdc:direction-distinctness` using the D2 thresholds above and `pdc:direction-length` at 750 words.
- `pdc:dtcg-2025-10`, `pdc:token-reference-direction`, `pdc:undeclared-token`, and `pdc:raw-colour-literal`.
- `pdc:off-scale-dimension` with declared-scale membership, not divisibility by 8.
- `pdc:near-duplicate-colour`: ΔE2000 `<1.0` violation; `1.0-2.3` needs review.
- `pdc:transition-all`, `pdc:layout-property-animation`, and `pdc:missing-reduced-motion`.
- `pdc:wcag-text-contrast` at 4.5/3.0, `pdc:wcag-nontext-contrast` at 3.0, and `pdc:horizontal-overflow` at all six declared viewports.
- `pdc:target-size-24` as `NEEDS_REVIEW`, because WCAG's spacing/equivalent/inline/user-agent/essential exceptions must be adjudicated.
- `pdc:undersized-ui-text` below 11 px, softened to 10 px for non-interactive small print, as design-system-blind `NEEDS_REVIEW`.
- `pdc:overused-font` as `RECOMMENDATION` only, requiring at least 20 text elements and 15% primary-family share.

Findings use exactly `VIOLATION`, `NEEDS_REVIEW`, `RECOMMENDATION`, `UNJUDGEABLE`, or `PASS`.
They include rule/version, file/line or selector, observed value, threshold, evidence, engine, confidence basis, remediation, and suppression status.
Rules use `yields_to` and disjoint thresholds so one defect has one owning id.

Exit codes are `0` clean/recommendations only, `1` blocking violations, `2` unresolved needs-review, `3` invalid configuration or required tier unavailable, and `4` checker failure.
The Parley driver runs the checker once at the reviewed commit. Reviewers consume the same report rather than producing divergent scans.
A T2 layout verdict may only be originated by the recorded T2 run; agents who did not see rendered evidence may acknowledge coverage but may not assert a layout pass.
If the profile requires an unavailable tier, the audit is `DEGRADED` and only the human Decider may waive that coverage gap.

`design-check.waivers.json` records `waiver_id`, `rule_id`, narrow `scope`, optional `value`, `reason`, `evidence`, `requested_by`, `approved_by`, `created_at`, `expires_at`, and `artifact_sha256`.
Every waiver needs one non-requesting participant's countersignature; a project-wide waiver additionally needs the human Decider and two participants.
Inline comments do not suppress in v1. The checker never writes waivers.
WCAG contrast, hidden/occluded content, missing interaction states, and undersized UI text remain design-system-blind: adding a bad token cannot legalize the output.

### 7. Positions on F1-F8 and honest operating limits

- **F1 — split authority.** Reproducible `quality` and contract violations block; an evidenced quality defect permits one reviewer to BLOCK. `slop` never blocks unilaterally and becomes an agreed fix only when at least two independent reviewers concur, or all available non-implementer reviewers when fewer than two exist. Category changes require a major spec version.
- **F2 — convergence is an alarm under the exact D2 thresholds.** One forced-axis rerun is allowed; persistent convergence needs a brief-specific human acceptance or `ABSTAIN`, never an automatic pass.
- **F3 — rendered evidence is optional-with-declared-capability.** The vocabulary is T0/T1/T2/T3 plus `UNJUDGEABLE`; no one may sign a layout verdict they did not obtain from pinned T2 evidence.
- **F4 — deterministic local dice are required on every full route.** The SHA-256 rank assignment above is checkable and hosted-service-free. The fast path does not generate directions and therefore needs no roll.
- **F5 — five Markdown files and 96 KiB for `parley-design`.** Phase read sets are fixed; participants may not choose different doctrine subsets.
- **F6 — centralized, expiring, reason-carrying waivers with countersignature.** Narrow value/path scope is required; no unilateral or bare wildcard suppression; objective legibility and accessibility rules are design-system-blind.
- **F7 — fast path skips direction deliberation only for an existing-system change affecting one surface, at most three existing components, at most five files/300 LOC, with no new foundation, component archetype, brand direction, security/accessibility policy, or token family.** It still runs T0/T1 and one model-diverse review. Greenfield work or any new visual-world decision cannot use fast.
- **F8 — human Decider by default, agents advisory.** Delegation must be explicit and removes that agent from proposal/scoring. Existing Parley signoff verifies process and objective obligations; it does not turn taste into a quorum vote.

## Concerns / open questions

The only implementation choice I would leave to Phase 2 is whether the baseline T1 parser dependencies are bundled or adapters to tools already present in the target repository; the tier contract must remain identical either way.
The DTCG official schema should be vendored with its license/provenance so offline validation does not depend on a live URL, but the stable upstream URL and exact `2025.10` hash must remain recorded.
The D2 categorical enums need calibration fixtures across marketing, product, documentation, and terminal surfaces; otherwise a mechanically precise distinctness check can still measure the wrong axes.
`content_sha256` canonicalization must specify line endings, excluded frontmatter fields, and UTF-8 bytes before implementation; signatures over ambiguous serialization are worthless.
The design audit should be one input to normal review consensus, not a new signoff surface or a silent amendment to `COOPERATION.md`.

## Risks

Multi-agent work can make design worse through feasibility-biased selection, self-preference, position/verbosity bias, stance homogenization, and the false confidence of several agents sharing no rendered evidence.
One critique round, recusal, fixed length, rotated order, a human Decider, and `ABSTAIN` mitigate those failures; they do not prove better taste.
Human aesthetic agreement around 38% means any taste score used as a hard gate would be pseudo-measurement.
False positives will otherwise turn the checker into a conservatism machine. Stable blockers therefore need counterexample fixtures, exact evidence, tier honesty, and narrow waivers.
Waiver growth can hollow out the system; expiry, countersignature, scope ordering, and design-system-blind quality rules are essential.
Anti-slop doctrine can become a recognizable inverse house style. Rules should attack reflex and incoherence, not ban visual devices, and should be reviewed against category-plus-avoidance predictability.
Post-build documentation can launder defects into tokens; reconciliation must precede canonization and unresolved contradictions remain review findings.
Rendered checks can vary across environments. T2 should assert deterministic geometry and standards facts, not screenshot similarity; T3 is deliberately excluded from v1.
The largest scope risk is rebuilding impeccable's detector fleet. If scope must shrink, cut T2 first, then component prose linting, and keep PDS phases, the direction contract, DTCG validation, token-reference integrity, T0/T1 source checks, the audit schema, and waiver governance.
