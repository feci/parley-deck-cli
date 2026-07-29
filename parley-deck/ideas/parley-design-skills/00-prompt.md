---
idea: parley-design-skills
author: user
created: 2026-07-28
participants: [claude-1, codex-1, hermes-1, kimi-1]
roles:
  claude-1: synthesis, parley-deck integration seams (addon boundary, packaging, portability across runtimes)
  codex-1: parley-design-check — enforcement tooling, rule-registry format, engine tiers, false positives, gate semantics
  hermes-1: the doctrine content itself — which design rules belong, artifact/token schema, standards alignment (DTCG, WCAG)
  kimi-1: the collaborative ritual and protocol shape — diverge/critique/decide/graft, anti-design-by-committee, conformance
status: final
track: deliberation
---

## Problem / idea

Add **two new companion add-on skills** to the Parley Deck ecosystem, alongside the
existing `parley-worktrees` and `parley-tracker` add-ons:

1. **`parley-design`** — a **doctrine + protocol** skill, pure markdown, **zero runtime
   dependencies**. It teaches a team of agents how to collaboratively produce a
   **design system** (tokens, type scale, colour, spacing, component specs) and then
   **apply it** to real UI — *without AI slop*. This is the "what good looks like and
   how we decide it together" skill.

2. **`parley-design-check`** — a separate **enforcement** skill that MAY ship real
   tooling (scripts, a rule registry, detectors). It mechanically finds design slop and
   design-system violations, and turns them into gate findings. This is the "prove it"
   skill.

The two are deliberately split so the doctrine stays portable to every agent runtime
while the enforcement layer is free to require Node, a parser, or a browser.

### Decisions already made by the project owner (binding, not open for re-litigation)

These were answered explicitly before this idea opened. Treat them as constraints; you
may propose *how* to satisfy them, not *whether* to.

- **D1 — Two skills, split by dependency weight.** `parley-design` = doctrine, zero
  dependency, maximum portability. `parley-design-check` = the heavier enforcement
  system. The owner explicitly rejected a single middle-ground skill.
- **D2 — Scope is system → then application.** Phase A produces the design system as
  canonical artifacts; Phase B applies it to concrete surfaces *and* audits that the
  implementation still obeys the system. Not "pages only", not "system only".
- **D3 — The collaborative ritual is DIVERGE → ADVERSARIAL CRITIQUE → ONE WINNER
  WHOLE → GRAFT.** Each participating agent independently proposes a *different*
  visual direction. They then cross-critique adversarially. **One direction wins in its
  entirety — never an average, never a merge of all four.** Then 2–3 *concrete, named*
  details may be grafted from the losing directions. The explicit purpose is to prevent
  taste-by-committee mush.
- **D4 — Protocol-shaped, not prose-shaped.** The owner likes how CopilotKit specifies
  AG-UI as a *protocol*: typed artifacts, explicit phases, normative language, a version,
  and a defined notion of conformance. `parley-design` should read like a specification
  a second implementer could conform to, not like a blog post of design advice.
- **D5 — Surface-agnostic core + a clearly-marked web annex.** The invariants
  (hierarchy, contrast, interaction states, honest copy, effect budget) are written so
  they hold for any surface. Every CSS-shaped rule (`oklch`, Tailwind classes,
  `overflow-x`, viewport widths) lives in a separate, explicitly-labelled **web annex**.
  A later TUI/CLI/docs annex must be addable without touching the core — Parley's own
  TUI is a plausible first non-web customer. Guard against the obvious failure: a
  surface-neutral core that degenerates into platitudes. Keep it short and hard.
- **D6 — Ship invariants, ratify taste per deck.** The skill ships **anti-slop
  invariants only** (contrast, states, motion timing, honest copy, the overused-font
  list, the effect budget). It ships **no theme catalogue** and no house aesthetic. The
  concrete visual world is decided and ratified *per project*, recorded in that project's
  own artifact, and is overridable on the record. Accept the known cost: on day zero of a
  greenfield project the skill supplies discipline, not a look. Say how the free axes are
  kept from silently collapsing back onto each model's training prior.
- **D7 — Two artifacts, not one: a contract before, a system after.** Before
  implementation, a short ratified **direction contract** (the winning direction plus its
  Named Rules) binds implementers and rides in `FINAL.md` — consistent with the Parley
  protocol. After implementation, a **design system** artifact is documented *from the
  code that was actually built*, per impeccable's finding that a rulebook written before
  the build gets defended against reality. Define both artifacts, their distinct
  authority, and the reconciliation step between them, including what happens when the
  built system contradicts the ratified contract.
- **D8 — Decisions and audit only; the skill does not own Phase-5 code.** `parley-design`
  owns the design deliberation and the design review/audit. Writing the UI code stays
  with the ordinary Phase-5 implementer, who obeys the ratified contract. This preserves
  the `parley-worktrees` / `parley-tracker` precedent that an add-on never changes
  canonical artifact ownership.

### Prior art you MUST study before writing your round-01 file

Full research digests are checked in next to this file, so you do not need network
access:

- `research/hallmark-doctrine.md` — the complete doctrine extraction from
  `Nutlope/hallmark` (MIT): its extensional definition of slop (named tells graded
  critical/major/minor), the **58-gate slop test**, the four-layer style system
  (genre → macrostructure → component archetype → theme), the three-axis diversification
  rule, the 8-step design flow, the CSS stamp and `design.md` system file.
- `research/hallmark-repo.md` — positioning, roadmap, recipes, the `_tests/` example
  corpus, theme catalog, packaging.
- `research/impeccable-commands.md` — `pbakaus/impeccable` (Apache-2.0): its ~35-command
  surface, routing, lifecycle, and the `craft-floor` / `critique` quality primitives.
- `research/impeccable-detector.md` — **the single most valuable artifact**: the complete
  60-rule antipattern registry with exact thresholds, the four detection engine tiers
  (regex → static HTML+CSS cascade → live browser → screenshot pixels), severity /
  advisory semantics, inline ignores, and the design-system conformance rules 57–60.
- `research/impeccable-state.md` — the persisted artifact model (`PRODUCT.md`,
  `DESIGN.md` + sidecar, surface briefs, critique snapshots), the staleness/drift model,
  sub-agents, hooks, and several ideas worth stealing outright: **Named Rules**, the
  **craft floor**, **"DESIGN.md is written AFTER the build, from the built world"**,
  **bounded verification with a hard ceiling of two correction rounds**, and the
  **"externalised dice"** finding.
- `research/impeccable-philosophy.md` — stated philosophy, honesty/evidence rules, and
  how one skill payload is made portable across 13 agent runtimes.
- `research/00-CONSOLIDATED-BRIEF.md` and `research/00-CONSOLIDATED-BRIEF-2.md` — two
  synthesis briefs. Brief 2 additionally covers: how CopilotKit/AG-UI specifies a
  protocol, peer AI design skills already in the wild, the W3C DTCG design-token
  standard and real design-system practice (Material/Polaris/Carbon/GOV.UK, Atomic
  Design), the empirical "AI slop" tell list from the design discourse, exactly which
  checks are mechanically automatable at which engine tier, and grounded human design
  critique practice (design-sprint decider ritual, crit formats, why committees fail).

**One measured finding from the prior art is the strongest argument for this whole
idea, and every round-01 file should reckon with it:** impeccable ships a script whose
header records that a single model, left alone, always builds its own #1 ranked concept —
*"Measured: 30/35 identical concepts across 16 prompt framings; the model cannot roll its
own dice."* impeccable had to bolt on an external random seed to escape the rut. **Parley
Deck gets that divergence for free, because the four participants are genuinely different
models.** Say clearly in your file whether you believe that advantage is real, and what
would squander it.

### What round-01 must answer

Write your independent proposal for the two skills. At minimum:

1. **The protocol shape of `parley-design`** — its table of contents, its normative
   phases, the exact canonical artifacts it makes agents produce (file names, schemas,
   frontmatter), and what "conformance" means. Propose a version number and a
   deprecation stance.
2. **The design-system artifact set** — what a Parley-produced design system physically
   *is* on disk. Reconcile with the W3C DTCG design-token format where sensible; say
   where you deliberately diverge and why.
3. **The ritual, made mechanical** — how D3 maps onto Parley Deck phases, rounds, and
   files. Who decides. What the decision rule is when the critique does not converge.
   How the graft is bounded so it does not silently become an average. What stops the
   losing directions from being re-litigated forever.
4. **The slop doctrine** — what actually goes in the doctrine, at what depth. hallmark
   is ~67 KB of SKILL.md plus 24 reference files and ~400 KB of references; impeccable is
   bigger still. State your position on size and token cost explicitly: what is the
   minimum doctrine that still changes output, and what belongs in lazily-loaded
   references versus the always-loaded core?
5. **The split line and the contract** — exactly what belongs in `parley-design` versus
   `parley-design-check`, and the *machine-readable contract* between them, so the
   zero-dependency doctrine skill can name a check it cannot itself run.
6. **`parley-design-check` design** — the rule registry format (data separate from
   detection logic, as impeccable does it), engine tiers and their dependency cost,
   severity/advisory semantics, false-positive suppression, and what the gate does in a
   Parley review phase.
7. **Honest risks** — where this fails, where multi-agent makes design *worse*, and what
   you would cut first if the scope has to shrink.

### Open forks the deliberation must settle (not pre-decided)

The consolidated research brief raised twelve forks. Four were answered by the owner and
are now D5–D8 above. These remain genuinely open — take a position on each, with reasons:

- **F1 — Rule authority.** Binding law, advisory doctrine, or split by category? The
  brief's proposal is a split: a `quality` finding (an objective defect — contrast,
  occlusion, missing state) lets a single agent BLOCK, while a `slop` finding (taste with
  a strong prior) requires quorum. The categorisation must be right on day one, because
  re-categorising a rule later invalidates earlier reviews.
- **F2 — Convergence semantics.** When several independent models propose the *same*
  direction, is that strong evidence the direction is right, or evidence of a shared
  training attractor? The two readings demand opposite protocol behaviour. D3 leans
  toward treating convergence as a slop alarm; specify precisely what triggers, and what
  the escape is, or it becomes a coin flip.
- **F3 — Rendered evidence.** Forbidden, required, or optional-with-declared-capability?
  Nothing in the roster renders by default. If evidence tiers are adopted, specify the
  tier vocabulary, the `unjudgeable` state, and the degradation banner, and say who may
  sign a verdict about layout they never saw.
- **F4 — Externalised dice.** Is cross-model divergence sufficient randomisation, or does
  each participant additionally need a deterministic seeded assignment forcing it to build
  its rank-k rather than rank-1 candidate? If seeded, the seed must derive locally from
  the run id — never from a hosted service — and the roll must be checkable.
- **F5 — Doctrine size and load model.** Name a hard ceiling on file count and bytes now.
  Both source projects rely on lazy per-request loading by a single agent; that economy
  inverts under N agents, who each pay the load and may read different subsets.
- **F6 — Waivers and suppression.** How a rule is waived, by whom, with what recorded
  reason, and whether a second participant must counter-sign. Suppression is where every
  rule system dies; at least one rule class should be deliberately design-system-blind so
  an implementer cannot legalise its own output by widening the system.
- **F7 — Fast path.** Running four agents plus a consensus round to decide whether a hero
  is centred is indefensible. Specify what `parley-design` does under the protocol's
  `fast` track, and the size threshold below which the full ritual is skipped.
- **F8 — Who is the Decider?** The critique research is emphatic that the deciding role
  must be a *role*, not a vote, and it proposes the human as the default Decider with
  agent scores explicitly advisory. Parley's own consensus model is agent signoffs. Take a
  position: human-decides-by-default with agents advisory, an agent Decider with a human
  override, or the existing signoff quorum. Whatever you choose must survive the judge-bias
  evidence below.

### Evidence that constrains this design — read before proposing

Wave-2 research surfaced published evidence that **contradicts the naive case for
multi-agent design**. Every round-01 file must reckon with it rather than assume more
agents means better taste. Details and citations are in
`research/design-critique-method.md` §5–§6.

1. **Multi-agent debate does not reliably beat a single agent.** The strongest published
   result is negative and is on *verifiable* tasks — on MMLU with GPT-4o-mini, debate
   scored **74.73 %** against **80.73 %** for plain chain-of-thought and **82.13 %** for
   self-consistency. No study shows debate improves *aesthetic* outcomes.
   **Consequence:** `parley-design` must NOT be justified as "more agents → better taste".
   Its defensible justification is that **multi-agent is a diversity generator, not an
   accuracy amplifier** — the value is producing genuinely different directions (which a
   single model demonstrably will not) and then selecting by rule.
2. **Deliberation degrades over rounds.** Documented pathologies: *factual attrition*
   (facts are progressively lost or misstated across rounds) and *stance homogenization*
   (positions converge regardless of correctness). **Consequence:** round count is a
   hazard, not a virtue. Default to **one** adversarial critique round; a second needs a
   stated reason. Stance diversity before vs after is a loggable metric.
3. **Diversity collapse is the mechanism of slop.** Multi-agent systems homogenise through
   structural coupling; aligned models exhibit mode collapse toward a narrow attractor.
   **Consequence:** isolation during divergence is not hygiene, it is the product — and a
   **mechanical distinctness gate is mandatory before any critique**. Critiquing a
   collapsed set launders the collapse into a "consensus".
4. **LLM-as-judge is measurably biased.** Position bias is worst exactly when the quality
   gap is small (our case); self-preference ranges roughly **−38 % to +90 %**; verbosity
   correlates with judge scores at **r ≈ .87** versus **.44** for humans.
   **Consequence:** no pairwise "which is better, A or B?" in the deciding phase — score
   absolutely against an anchored rubric; randomise presentation order; **no agent scores
   its own direction** (discard, do not down-weight); cap and normalise artifact length
   before scoring.
5. **Human agreement on aesthetics is ~38 %.** Two human labelers pick the same option on
   only about a third of aesthetic comparisons — barely above the task's floor. A judge
   reached that same human level *only* by abstaining on ~35 % of cases.
   **Consequence:** an aesthetic score can never be a hard gate. Gates belong on
   mechanically checkable properties; taste stays advisory; and **`ABSTAIN` must be a
   legitimate, preserved verdict** that escalates rather than being coerced into a vote.

`research/design-critique-method.md` §7 additionally contains a fully drafted candidate
protocol (`DCP/1`) with typed phases and artifacts, derived from real practice (the design
sprint's diverge/heat-map/speed-critique/straw-poll/supervote sequence, Pixar's Braintrust,
d.school's I-like/I-wish/What-if, Nielsen severity anchors, the Awwwards trimmed-mean
rubric). **Treat it as one participant's strawman, not as the answer** — attack it,
improve it, or replace it, but do not ignore it.

## Constraints

- **Vendor-neutral.** Both skills must work for any participant CLI in the roster
  (claude, codex, hermes, kimi, agy) and be installable into every runtime the
  `parley-deck-skill` installer targets. No skill may assume a specific vendor's tools,
  model, or IDE.
- **`parley-design` has zero runtime dependencies.** Markdown only. It may *describe*
  checks and *reference* the check skill, but must remain fully useful with no Node, no
  browser, and no network.
- **`parley-design-check` may ship code**, but must degrade honestly: state which tier
  each check needs, and never claim coverage it cannot deliver in the tier available.
  Follow the existing precedent set by `parley-tracker`'s `bin/validate.js` +
  `bin/claim.js` gap-scan.
- **Add-on layout.** Both live at `parley-deck-skill/addons/<name>/` with `SKILL.md`
  (plus `reference/`, `bin/`, `templates/` as needed), discovered by the installer's
  `discoverAddons()` glob, installed by default with `--no-addons` / `--only` opt-outs.
- **Licensing / originality.** hallmark is MIT, impeccable is Apache-2.0, this repo is
  Apache-2.0 — derivation is legally permitted with attribution. Nevertheless the
  default is **independent authorship**: study the prior art, then write our own rules in
  our own words. Any material genuinely derived from either project must be attributed in
  a `NOTICE.md`. Do not paste large blocks of their rule text.
- **Protocol boundary.** These are add-on skills. They must not silently amend
  `COOPERATION.md`. If the core protocol genuinely needs a seam (as the worktrees/tracker
  add-ons did), name it as a deferred meta-protocol-change idea rather than doing it here.
- **English only** for every file under `parley-deck/`.
- **Token cost is a first-class constraint.** hallmark explicitly treats over-eager
  reference loading as its largest avoidable cost. Any design that loads hundreds of KB
  on every invocation is a failed design.
- **The doctrine must not become its own cliché.** A rule set that merely inverts the
  current defaults produces a new, equally recognisable house style. Say how you avoid
  that.

## Non-goals

- Not a Figma plugin, not a design-tool integration, and not dependent on any hosted
  design service.
- Not a full re-implementation of impeccable. We are not shipping a browser-based
  detector fleet, a live-editing server, hooks, or a 35-command surface in v1.
- Not a component library, not a CSS framework, and not a set of ready-made themes to
  copy. The output is a *system for a specific project*, not a catalogue.
- Not a replacement for human taste or for the user's brand. The user's stated brand,
  constraints, and preferences always outrank the doctrine.
- Not a change to the core Parley Deck protocol in this idea.
