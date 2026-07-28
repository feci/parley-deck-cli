---
agent: hermes-1
idea: parley-design-skills
round: 1
date: 2026-07-28
---

## Summary

I propose `parley-design` as a versioned specification (PDS/1.0) with RFC 2119
normative language, 8 numbered phases (D0–D7), typed canonical artifacts, and
a conformance section. The doctrine ships invariants only — no theme catalogue,
no house aesthetic (D6). Token artifacts use W3C DTCG `2025.10` verbatim; OKLCH
is mandatory for color; WCAG 2.2 ratios are blocking, APCA is advisory. The
slop doctrine is a two-class rule set (quality vs slop), ~30 named tells graded
critical/major/minor, with a ~15 KB always-loaded core and lazily-loaded web
annex. The split with `parley-design-check` is a machine-readable `rules.json`
that is the single source of truth — the markdown catalog is generated from it,
guarded by a drift test. The multi-agent ritual (D3) is justified as a diversity
generator, not an accuracy amplifier; one critique round is the default; ABSTAIN
is a preserved verdict; no agent scores its own direction; grafts are capped at
3 and must cite source + property + non-contradiction.

## Proposed approach

### Q1 — Protocol shape of `parley-design`

**Version:** `PDS/1.0` (Parley Design Spec). Semver of the spec, not the tooling.
Every artifact carries `spec: PDS/1.0` in frontmatter from day one — retrofitting
a version field is the hardest thing to add to a live protocol (AG-UI shipped
without one and already has two SDKs version-skewed).

**Normative language:** RFC 2119. Uppercase MUST/SHOULD/MAY exclusively for
normative statements. This makes rules greppable and unambiguous for both agents
and checkers.

**Phases (normative spine):**

| Phase | Name | Exit condition |
|-------|------|----------------|
| D0 | Brief | `DESIGN-BRIEF.md` written; ≥3 anti-goals; ≥1 anti-reference |
| D1 | Diverge | N `DIRECTION-<agent>.md` + `.tokens.json` produced; distinctness gate passed |
| D2 | Critique | N `CRITIQUE-<agent>.md` produced; each cites target + rule ids |
| D3 | Decide | `VERDICT.md` names exactly 1 winner + 2–3 grafts; re-litigation blocked |
| D4 | Graft | Winner's tokens + grafted details merged into `design-system.tokens.json` |
| D5 | Tokenize | DTCG token files validated against `2025.10/format.json` schema |
| D6 | Apply | UI implemented by Phase-5 implementer; `APPLICATION.md` maps tokens→components |
| D7 | Audit | `DESIGN-SYSTEM.md` written from built code; `AUDIT.md` produced by check |

**Conformance (§10):** Four levels — L1 Artifacts (required files exist with
correct frontmatter), L2 Phase-order (phases ran in order, exit conditions met),
L3 Token-integrity (DTCG schema valid, alias graph is a DAG, no raw literals
in applied code), L4 Applied-UI (contrast matrix passes, states complete,
no banned tells). An implementation conforms at level N iff levels 1..N all pass
`parley-design-check --strict`.

**Deprecation:** Deprecated rules keep validating for ≥1 minor spec version.
A deprecation table in the spec carries: rule id, version introduced, version
that removes it, replacement (if any). A retired rule's id is never reused.

### Q2 — Design-system artifact set

Physical layout under `parley-deck/ideas/<slug>/`:

```
DESIGN-BRIEF.md
directions/
  DIRECTION-<agent>.md
  DIRECTION-<agent>.tokens.json
VERDICT.md
design/
  tokens/
    primitives.tokens.json    # DTCG, raw OKLCH values, no semantic names
    semantic.tokens.json      # DTCG, aliases only ({...} / $ref)
    component.tokens.json     # DTCG, aliases to semantic only
    resolver.json             # DTCG Resolver 2025.10 — light/dark modes
  FOUNDATIONS.md
  TYPE-SCALE.md
  COLOR-SYSTEM.md             # full contrast matrix
  SPACING-LAYOUT.md
  MOTION.md
  VOICE-AND-CONTENT.md
  components/<c>.spec.md      # one per component
  PATTERNS.md
  ACCESSIBILITY.md
  DESIGN-DECISIONS.md         # ADR log
DESIGN-SYSTEM.md              # written AFTER build, from shipped code (D7)
design-check.config.json      # enforcement contract
```

**DTCG reconciliation — adopt verbatim, diverge deliberately in three places:**

1. **Adopt:** `$value`/`$type`/`$description`/`$extensions`/`$deprecated`;
   13 `$type` values; aliasing via `{group.token}` and `$ref` JSON Pointer;
   the official JSON Schema at `https://www.designtokens.org/schemas/2025.10/
   format.json` as the validation target. Zero invention here — this is the
   cheapest, hardest gate available.

2. **Mandatory `colorSpace: "oklch"`** for all color primitives, with a 6-digit
   `hex` fallback. Enforced by a `core/colorspace`-equivalent rule. DTCG lists
   `oklch` among 14 supported colorSpaces — this is standards-aligned, not a
   fashion. HSL/RGB are prohibited in token files.

3. **Three-tier layering** (primitives → semantic → component) from Material's
   `md.ref` / `md.sys` / `md.comp` — copy the *layering*, not the Google-branded
   prefixes. Use `primitive.` / `semantic.` / `component.` or Atlassian's
   `foundation.property.modifier` anatomy. Enforce direction of references as
   a graph assertion: component→semantic→primitive only, never skipping or
   reversing. A cycle is a hard fail.

**Divergence from DTCG:** None in the token file format itself. The divergence
is in *what we require on top of it*: (a) every semantic color token MUST have
a value in every declared resolver context (light + dark minimum); (b) every
component spec MUST reference a composite `typography` token (font-family +
size + weight + line-height + letter-spacing as a unit), never an isolated
`font-size` — this is Polaris's `--p-text-heading-lg-*` pattern expressed as
DTCG's composite type; (c) `core/duplicate-values` is error at ΔE2000 < 1.0,
needs-review at 1.0–2.3 (measured thresholds, not DTCG's concern).

**WCAG 2.2 as blocking constants (W3C Recommendation, 12 Dec 2024):**
4.5:1 body text, 3:1 large text (≥24px / 18.66px bold), 3:1 non-text/UI
components, 24×24 CSS px targets (SC 2.5.8), text-spacing survivability at
1.5×/2×/0.12em/0.16em (SC 1.4.12). APCA Lc is advisory only — WCAG 3 is still
a Working Draft, contrast was pulled from it in July 2023, algorithm "yet to
be determined," no Recommendation before 2028. Do not fail a build on APCA.

**Named Rules as the doctrine unit** (from impeccable, proven in practice):
`**The [Name] Rule.** [one forceful sentence]`, tagged with a section, 1–3
per section, mirrored into `DESIGN-SYSTEM.md` as machine-readable
`narrative.rules[]`. Citable, diffable, reviewable across agents. Examples:
*The OKLCH-Only Rule*, *The Weight-Inversion Rule*, *The Hairline-First Rule*,
*The One-Accent Rule* (accent ≤ ~5% of viewport, atmospheric genre relaxes
to ~20%), *The Honest-Copy Rule* (no fabricated metrics — labelled hole or
ask the user).

**DESIGN-SYSTEM.md is written AFTER the build** (D7), from shipped code, by
a dedicated documenter role — per impeccable's finding that "a rulebook
written before the build gets defended against reality instead of describing
it." The *direction contract* (D7's "before" artifact) rides in `VERDICT.md`
and binds implementers; the *design system* is the post-build description.
Where they diverge, the build wins and the prose notes the divergence.

### Q3 — The ritual, made mechanical

D3 maps onto phases D1–D4. The mechanics:

**D1 Diverge:** Each agent independently produces `DIRECTION-<agent>.md` +
`.tokens.json`. Agents MUST NOT read each other's directions (protocol rule 1,
enforced by phase isolation). Each direction carries: name, one-sentence
thesis, mood, typographic idea, color idea, spatial idea, **the one risky
move**, what it deliberately sacrifices, and a **Signature** field ("the single
unique element this page will be remembered by" — from frontend-design). Two
directions with the same Signature are the same direction.

**Distinctness gate (mandatory, before critique):** Compute each direction's
three-axis fingerprint: **paper-band** (dark L<30% / mid 30–85% / light >85%)
× **display-style** (serif/sans/mono/condensed/heavy) × **accent-hue**
(warm 10–60° / cool 200–300° / neutral / chromatic-other). If ≥3 of N agents
share all three axes, that is a **convergence alarm** (F2) — either a shared
training attractor or a correct-but-templated direction. The escape: the
convergent agents re-roll (see F4) or the human Decider flags it as a slop
signal. This gate is what prevents "critiquing a collapsed set" — the failure
mode that launders mode collapse into false consensus.

**D2 Critique:** One round is the default. A second round requires a stated
reason (deliberation degrades over rounds — factual attrition + stance
homogenization are documented pathologies). Each agent critiques every other
agent's direction using the anchored rubric (see below). Critique artifacts
cite target direction id + rule ids. No agent critiques its own direction.

**D3 Decide — the anti-averaging rule (binding):** Exactly one direction wins
whole. A synthesis of two directions is a protocol violation. 2–3 named
details may be grafted from losing directions; each graft MUST name (a) its
source direction, (b) the exact token/component/rule it changes, (c) why it
does not contradict the winner's thesis. An unbounded graft list is an average.

**Decision mechanism (F8):** Human-decides-by-default, agents advisory. No
pairwise "which is better, A or B?" — LLM-as-judge has position bias worst
when quality gap is small (our case), self-preference ranges −38% to +90%,
verbosity correlates with judge scores at r≈.87. Instead: **score absolutely
against an anchored rubric** (the six-axis pre-emit critique: Philosophy,
Hierarchy, Execution, Specificity, Restraint, Variety — 1–5 each, peer-scored
not self-scored). **Randomize presentation order.** **Discard self-scores**
(do not down-weight — discard). **Cap and normalize artifact length** before
scoring. **ABSTAIN is a preserved verdict** that escalates to the human rather
than being coerced into a vote — human aesthetic agreement is ~38%, barely
above floor, so forcing a vote manufactures false consensus.

**Re-litigation block:** Once `VERDICT.md` is signed, losing directions are
closed. A losing direction may be re-opened only by an explicit human
override with a stated reason, recorded in `DESIGN-DECISIONS.md`.

### Q4 — Slop doctrine

**Two-class rule taxonomy** (from impeccable.style/slop, the most mature
registry): `quality` rules (objective defects — contrast, occlusion, missing
state — any single agent can BLOCK) and `slop` rules (taste with a strong
prior — purple gradient, Inter-only, centered hero — requires quorum per F1).

**Core doctrine (~15 KB, always-loaded):** Invariants only, surface-agnostic
(D5). These are the rules that hold for any surface:

- **Contrast:** WCAG 2.2 ratios (4.5:1 body, 3:1 large/non-text) — blocking
- **Interaction states:** 8 mandatory states (default, hover, focus-visible,
  active, disabled, loading, error, success) — every interactive element
- **Motion:** 3 duration buckets (100–150ms micro, 200–300ms minor, 300–500ms
  major); 3 named easings only; animate `transform`+`opacity` only;
  `prefers-reduced-motion` fallback mandatory; ≤3 animation primitives per page
- **Honest copy:** no fabricated metrics (the one rule with an ethical basis);
  labelled hole or ask the user
- **Effect budget:** decoration must carry a `motivation:` field; unmotivated
  = removed; the effort/complexity hierarchy defaults to "nothing"
- **Overused-font list:** versioned, dated, with `added`/`deprecated`/`confidence`
  per entry — not a ban, a warning. The tell is *absence of a typographic
  decision*, not any particular face. Fraunces went from *the fix* to *the tell*
  in ~1 year; a frozen list is a slop generator on a delay.
- **Structural variety:** "structural sameness is the AI fingerprint, not visual
  sameness" — each direction names its macrostructure/shape; if 3 agents
  independently name "hero + 3 features + CTA + footer", that IS the finding
- **Token discipline:** every value resolves to a declared token; no mid-render
  improvisation (inline hex/oklch outside `:root` is a fail)

**Web annex (~10 KB, lazily-loaded, explicitly labelled):** CSS-shaped rules
that only apply to web surfaces: `overflow-x: clip` not `hidden`; `minmax(0,
1fr)` on image tracks; `100vw` banned; `dvh`/`svh` not bare `vh`; four
mandatory widths (320/375/414/768px); focus rings never animate; `font-variant-
numeric: tabular-nums` on data; one icon library per project; no emoji as
feature icons; no `z-index: 9999` (six named levels); no `transition: all`;
no bounce/elastic on UI. A future TUI/CLI/docs annex is addable without
touching the core.

**Size and token cost (F5):** Hard ceiling — core ≤15 KB, web annex ≤10 KB,
total references ≤50 KB. The economy inverts under N agents: five agents each
independently loading the same conditional files multiplies cost 5×. The
answer is one pre-digested, versioned doctrine artifact in the deck that every
agent reads once — not per-agent lazy loading.

**Anti-cliché mechanism:** Rules express **prohibitions of defaults**, never
**prescriptions of alternatives**. "Do not ship the model's first palette
unmodified" survives; "use Fraunces" does not. The registry is versioned with
a sunset review cadence. Structural distance from the project's own prior
outputs is a first-class constraint (Hallmark's axis F: "colour-swaps don't
count as variety").

### Q5 — Split line and contract

**`parley-design` owns:** the doctrine (invariants, named tells, rule
definitions, the ritual, the artifact schemas). Pure markdown. Zero deps. It
*names* checks it cannot run via the `enforced-by` annotation.

**`parley-design-check` owns:** the rule registry as executable data, the
detection engine, the gate semantics, the suppression protocol. May ship
Node scripts, may delegate to stylelint/axe when present.

**The machine-readable contract:** `rules.json` is the single source of truth.
Each entry: `{ id, spec_section, severity, class: "quality"|"slop", tier:
"T0"|"T1"|"T2"|"T3", targets[], predicate, autofix?, since, deprecated?,
enforced_by: "check#N"|"agent-judgement" }`. The markdown rule catalog in
`parley-design` §8 is **generated** from `rules.json`, guarded by a drift test
(exactly the pattern parley-deck-cli already uses for
`TestEmbeddedDefaultMatchesLiveDeck`). Never two hand-maintained copies —
AG-UI's `events.proto` (16 types) vs its TS/Python SDKs (33 types) is the
cautionary tale.

Every doctrine rule carries `enforced-by: check#N` or `enforced-by:
agent-judgement`. That single annotation keeps both skills honest: the doctrine
cannot claim mechanical enforcement it doesn't have, and the checker cannot
grade taste it cannot measure.

### Q6 — `parley-design-check` design

**Rule registry:** `rules.json` (data separate from detection logic, as
impeccable does). Generated into the markdown catalog. Drift-guarded.

**Engine tiers and dependency cost:**

| Tier | Engine | Needs | Cost | What it unlocks |
|------|--------|-------|------|-----------------|
| T0 | regex / line scan | nothing | µs–ms | hex-not-token, `transition: all`, banned fonts, emoji-as-icon, `100vw`, gradient-text |
| T1 | CSS parse (AST) | a CSS parser | ms | token conformance, off-scale spacing, specificity, z-index, near-duplicate colors (ΔE2000) |
| T2 | static DOM parse | HTML parser | ms | heading order, `alt`, decorative SVG without `aria-hidden`, duplicate `id` |
| T3 | headless browser | Chromium ~150 MB | ~0.5s/page | computed contrast, target size, horizontal overflow @320–1920px, `:focus-visible`, reflow |

T0+T1+T2 ship by default with zero runtime dependency — ~60% of the value at
~0% of the cost, runnable anywhere (Go, shell, Python, Node). T3 is optional,
auto-detected, with a clean degradation message: "headless browser not found;
checks 16–21 skipped; contrast unjudgeable." Never claim coverage the tier
cannot deliver.

**Severity semantics — tri-state (from IBM Equal Access):** `violation` (blocks
the phase gate), `needs-review` (routed to agents as required consult),
`recommendation` (advisory). A binary pass/fail on design work is a lie — axe
itself fails *open* on gradient backgrounds ("could not be determined"), so a
"zero violations" report on a gradient hero is a false all-clear. Emit
`needs-review` (never `pass`) whenever a check returns indeterminate.

**Exit codes:** 0 = clean, 2 = findings, 1 = command failed. Distinguishing
"clean" from "tool broke" is what makes it CI-safe.

**False-positive suppression (F6):** Narrowest-first ladder: `ignore-value
<id> <value>` → `ignore-value <id> "*" --file <glob>` → `ignore-file <glob>`
→ `ignore-rule <id>`. Each stored ignore carries `createdAt` + `reason`. A
bare `"*"` with no `--file` is refused. The checker itself never writes ignore
config — always through a `hook-admin`-equivalent command. A second agent
MUST counter-sign a rule-level suppression. **At least one rule class is
design-system-blind:** the honesty gate (no fabricated metrics) cannot be
waived by widening the design system — an implementer cannot legalise its own
output by declaring "50,000+ users" as a token.

**Gate in Parley review phase:** `parley-design-check --strict` runs as a
Phase-6 gate, analogous to `RunChecks`. `violation` findings block; `needs-
review` findings are routed to the adversarial agents as required consult;
`recommendation` is advisory. The standing steer on a clean check: "a green
mechanical pass does not mean the design is good — keep following the doctrine."

### Q7 — Honest risks

1. **Multi-agent does not improve taste.** The strongest published result is
   negative (debate scored 74.73% vs 80.73% CoT vs 82.13% self-consistency on
   MMLU). No study shows debate improves aesthetic outcomes. The defensible
   justification is that multi-agent is a **diversity generator** (impeccable's
   measured 30/35 identical concepts; Verbalized Sampling's 1.6–2.1× diversity
   from asking for N candidates), not an accuracy amplifier. Parley gets
   divergence for free because the four participants are genuinely different
   models — that advantage is real. What squanders it: (a) critiquing a
   collapsed set before the distinctness gate; (b) more than one critique round
   (deliberation homogenises stance); (c) averaging directions instead of
   picking one whole; (d) letting agents score their own work.

2. **38% human agreement on aesthetics.** An aesthetic score can never be a
   hard gate. Gates belong on mechanically checkable properties; taste stays
   advisory; ABSTAIN must be preserved. Forcing a vote manufactures false
   consensus at the ~38% floor.

3. **The doctrine becomes its own cliché.** A rule set that inverts current
   defaults (cream + serif + terracotta) produces a new, equally recognisable
   house style — S4's detector already flags Fraunces, which was S6's
   *prescribed fix* one year earlier. Defense: rules prohibit defaults, never
   prescribe alternatives; the registry is versioned with a sunset review;
   structural distance from prior outputs is first-class.

4. **What I would cut first:** `parley-design-check` in its full T3 form. The
   doctrine skill is standalone useful — it changes output even with zero
   tooling. Ship T0+T1 (regex + CSS-AST, zero deps) as v1; T3 (headless
   browser) as v2 when the roster has browser-capable agents. The doctrine
   names the checks; the checker arrives when it can.

## Concerns / open questions

- **D5 surface-neutral core vs platitudes:** The core invariants (contrast,
  states, motion timing, honest copy) are genuinely surface-agnostic. But
  "hierarchy" and "structural variety" are harder to express without
  referencing layout primitives. The web annex handles web; a TUI annex is
  plausible (Parley's own TUI is a candidate). The risk is a core so abstract
  it becomes platitudes. Mitigation: keep the core short and hard — every rule
  must be checkable or name the agent-judgement that owns it.

- **F3 rendered evidence:** Nothing in the roster renders by default. If an
  agent cannot see the rendered UI, it cannot sign a verdict about layout,
  overflow, or contrast-on-gradient. Proposed: `evidence-tier` declared per
  agent (`rendered` / `static` / `unjudgeable`); a verdict about layout from
  an `unjudgeable` agent is downgraded to `advisory` and the degradation is
  bannered in the critique artifact. No agent signs a verdict about layout it
  never saw.

- **F4 externalised dice:** Cross-model divergence is the primary mechanism
  (different models, different training medians). A deterministic seeded
  assignment forcing each agent to build its rank-k rather than rank-1 candidate
  is a supplementary defense — the seed derives locally from the run id, never
  from a hosted service (impeccable's API is a non-starter here). The roll is
  checkable: `deterministicRank(items, seed)` with the seed key recorded in the
  artifact. A contract with no seed key means the roll was skipped — that is
  a material finding ahead of any craft point.

- **F7 fast path:** Under `track: fast`, a single agent applies the core
  doctrine invariants directly — no divergence round, no critique, no verdict.
  Threshold: single surface, no design system needed, brief is ≤500 chars or
  explicitly flagged `fast`. The ritual is skipped; the invariants are not.

## Risks

- **Convergence-as-slop (F2):** If all four agents independently land on the
  same three-axis fingerprint, is that evidence the direction is right or
  evidence of a shared training attractor? I take the position: **convergence
  is a slop alarm until proven otherwise.** The escape is re-roll or human
  override. The alternative (treating convergence as validation) reconstructs
  the median in miniature — which is the exact disease this protocol exists to
  cure.

- **Rule re-categorisation (F1):** The quality/slop split must be right on day
  one because re-categorising a rule later invalidates earlier reviews. I
  propose: quality = objectively measurable defect (contrast ratio, missing
  state, horizontal overflow) — single agent can BLOCK. Slop = taste with a
  strong prior (purple gradient, centered hero, Inter-only) — requires quorum.
  The honesty gate is quality (ethical basis, regex-detectable number-shape),
  not slop.

- **Doctrine rot:** The anti-slop canon is a moving target with a half-life
  (Fraunces: fix → tell in ~1 year). The registry MUST be versioned and dated
  with per-rule `added`/`deprecated`/`confidence`/`sources` and a stated review
  cadence. A frozen anti-slop list is a slop generator on a delay.

- **Token cost under N agents:** Five agents each loading the same conditional
  references multiplies cost 5×. The doctrine must be pre-digested into one
  versioned artifact, not rely on per-agent lazy loading.

- **Disagreement with DCP/1 strawman:** The DCP/1 strawman in
  `design-critique-method.md` §7 is a reasonable starting point but I disagree
  on three points: (1) it does not make the distinctness gate mandatory before
  critique — without it, the protocol launders convergence into false consensus;
  (2) it does not discard self-scores — down-weighting is insufficient given
  the −38% to +90% self-preference range; (3) it does not preserve ABSTAIN as
  a first-class verdict — at 38% human agreement, forcing a vote is
  manufacturing false confidence. My counter-proposal: mandatory distinctness
  gate (three-axis, pre-critique), absolute scoring against an anchored rubric
  with self-scores discarded, and ABSTAIN as a preserved verdict that escalates
  to the human Decider.
