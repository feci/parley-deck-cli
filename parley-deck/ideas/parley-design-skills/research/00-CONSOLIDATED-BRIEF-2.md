# Consolidated brief #2 — `parley-design` + `parley-design-check`

**Inputs:** six wave-2 analyst digests (AG-UI/CopilotKit protocol style · peer AI design skills · design-system standards · AI-slop discourse · automatable checks · design-critique method) plus the wave-1 brief and the hallmark/impeccable deep-dives. All read in full.
**Audience:** the humans deciding what these two skills are before a line is written.
**Posture:** decisive. Where analysts disagreed, this brief picks and says why. Open forks are isolated in §6.

---

# 1. Protocol style verdict

## 1.1 What "protocol-shaped" actually means (five mechanics, not a tone)

AG-UI's style reduces to five mechanics. Adopt all five; every one of them is cheap.

| # | Mechanic | AG-UI evidence | Our instantiation |
|---|---|---|---|
| M1 | **A small closed vocabulary of typed artifacts** | `EventType` enum, 33 members, `z.nativeEnum` discriminant + `.passthrough()` payload | 9 artifact kinds, closed. Fixed front matter, free-form body, `x-` prefix for extension keys |
| M2 | **An explicit phase state machine with typed terminal states** | `RunStarted → … → RunFinished{outcome} \| RunError`; "cannot finish while children are open" | `D0 → … → RATIFIED{winner, grafts[]} \| ABANDONED{reason}`; cannot RATIFY while any critique is unanswered |
| M3 | **One machine-readable source of truth; prose derived from it** | Zod schemas → `z.infer` types → doc tables | **One literate registry file** whose fenced blocks are the machine source and whose prose is the human source (§5.3). No second copy, ever |
| M4 | **Conformance is an executable check whose error strings ARE the rule text** | `verify.ts`: `"Cannot send 'RUN_FINISHED' while tool calls are still active: …"` | `parley-design-check` emits `rule-id — violation — remedy`, always all three |
| M5 | **Extensions are reserved-namespace or typed escape hatches; unknown values MUST NOT error** | `core:` reserved, `<framework>:<name>` custom, "for unknown reasons, clients must not error" | `core:` reserved for spec rule ids; project rules MUST be `<project>:<slug>`; consumers MUST NOT error on unknown rule ids or unknown token groups |

## 1.2 The nine defects of AG-UI we must not reproduce

Verified by the analyst, each with a one-line countermeasure that goes in the spec on day one.

| AG-UI defect | Countermeasure (binding) |
|---|---|
| No protocol version anywhere on the wire (0 grep hits) | `spec: PDS/1.0` in the front matter of **every** artifact, from the first commit |
| `events.proto` rotted to 16 types vs 33 in TS/Python | **Never two hand-maintained representations.** One literate registry; a drift-guard test if any generated view exists |
| README still says "~16 event types" (actual 33) | **Never write a count in prose.** Generate the index or omit the number |
| Zero RFC 2119 keywords; "should" vs "must" is guesswork | `conformance-language: RFC 2119`; uppercase MUST/SHOULD/MAY reserved *exclusively* for normative statements; lowercase elsewhere |
| Changelog has one entry; roadmap is an empty heading | Either maintain §16 Changelog per version, or delete the section. An empty Development section signals abandonment |
| Literal placeholders shipped (`> **Logo strip goes here**`) | Zero placeholders. Fatal in a *design* doctrine where credibility is the product |
| `verifyEvents` locked inside an RxJS runtime package | `parley-design-check` MUST run standalone against files on disk, no agent runtime, no framework |
| Conformance harness is a Next.js app + 27 services | Fixtures are plain files + one script, offline, < 5 s |
| Python models hand-mirrored, no parity gate | Every generated view is guarded by a failing test (the `TestEmbeddedDefaultMatchesLiveDeck` precedent) |

Two AG-UI choices we deliberately **invert**:

- **`.passthrough()` as default.** Correct for AG-UI (absorb arbitrary frameworks); wrong for us. Unknown top-level keys in `tokens.json` and artifact front matter SHOULD warn; `x-` is the explicit silent-extension prefix.
- **"Loose event format matching" as a headline principle.** "Close enough" is precisely what produces slop. Our loose seam is the **target profile** (loose about platform syntax, strict about token identity).

## 1.3 Normative-language rules (verbatim, go in §0.5)

1. **MUST / MUST NOT / SHOULD / SHOULD NOT / MAY** appear in uppercase **only** when normative (RFC 2119/8174).
2. Every normative statement is **numbered and bold-named**, in the `interrupts.mdx` shape: `3. **Cover all open critiques.** A VERDICT MUST address every CRITIQUE raised against the winning direction. Partial verdicts are not supported.` Name = citable by a reviewing agent. Number = citable by the checker.
3. Every artifact kind, phase, and rule gets an **identically-shaped entry**: H3 name → one-line purpose → a *rationale paragraph* (why it exists, who consumes it) → a required-fields table → a minimal example. The shape never varies. This is the single highest value-per-effort item in the whole research corpus.
4. **Informative** sections are labelled `(informative)` in the heading. Anything unlabelled is normative.
5. Interop concessions are recorded **inline with the target that forced them** (AG-UI's best practice: naming the offending producer in the schema comment).

## 1.4 Versioning, deprecation, drafts

- `spec: PDS/<major>.<minor>` — semver **of the spec**, never of the tooling. `parley-design-check` versions independently and declares `implements: PDS/1.0`.
- The rule registry versions separately: `registry: core-rules/1.4.0`, with a `registry-digest` (12 hex of sha256) carried in every AUDIT artifact so a signature cannot silently survive a registry edit (impeccable's content-hash-pinned approvals).
- **Deprecation lives in three synchronized places** and deprecated rules **keep validating for ≥1 minor spec version**: (a) `deprecated: since=1.2 removes=2.0 replaced-by=core:<slug>` in the rule's own metadata block, (b) a deprecation table in §13, (c) a migration note. Deprecated rules stay in the registry so old artifacts still validate.
- **Drafts use the 5-state ladder** — Draft → Under Review → Accepted → Implemented → Withdrawn — and appear **inline in the canonical catalog with a DRAFT badge**, run by the checker as **warn-only**. This is how a new anti-slop tell enters without being binding on day one.

## 1.5 Conformance levels

| Level | Name | "An implementation conforms at level N iff…" |
|---|---|---|
| **L1** | Artifacts | every required artifact for the declared track exists, carries `spec:`, and passes required-section lint |
| **L2** | Process | the phase order was respected, gates G0–G4 were evaluated and recorded, and VERDICT names exactly one winner with 0–3 typed grafts |
| **L3** | Token integrity | `*.tokens.json` validate against DTCG `2025.10`; alias direction holds; no raw literal outside the token layer; contrast matrix recomputes clean |
| **L4** | Applied UI | the rendered artifact passes the browser-tier `QUALITY` rules for every declared target profile |

L1+L2 need no runtime at all. L3 needs a JSON validator. L4 needs a browser. A project declares the highest level it claims; the checker verifies the claim.

## 1.6 THE PROPOSED TABLE OF CONTENTS

Skill file layout — **hard ceiling of four files**, because five agents each lazily loading references costs 5× and guarantees divergent reads:

```
parley-design/
├── SKILL.md                 ≤200 lines. Dispatcher + when-to-use + when-NOT-to-use + the 4 invariants.
├── references/PDS.md        THE PROTOCOL (the TOC below).
├── references/RULES.md      The literate rule registry (§3 of this brief).
└── references/WEB-ANNEX.md  Target-specific hard numbers (CSS/HTML). Explicitly non-normative for other targets.
```

`references/PDS.md`:

```
---
spec: PDS/1.0
status: stable
conformance-language: RFC 2119
registry: core-rules/1.0.0
registry-digest: <12-hex>
---

§0   Scope and Non-Goals                                            NORMATIVE
     0.1  What PDS governs
     0.2  What PDS does NOT govern (taste oracles, house aesthetics, code review)
     0.3  Relationship to parley-deck COOPERATION.md
          (§4.0 track · §9.0 preflight · §12 pipeline · §13 retro · §14 human brake)
     0.4  Relationship to parley-design-check (rules live here; the runner lives there)
     0.5  Conformance language and how to read this document

§1   Terminology                                                    NORMATIVE
     Brief · Direction · Signature · Divergence axis · Heat mark · Critique ·
     Verdict · Graft · Rumble · Token · Primitive/Semantic/Component ·
     Slop rule · Quality rule · Evidence tier · Target profile · Waiver · Ledger

§2   Principles                                                     INFORMATIVE (exactly 6)
     P1 Slop is the absence of a decision, not ugliness
     P2 Diverge before you converge (VS: 1.6–2.1× diversity, arXiv 2510.01171)
     P3 One direction wins whole — averaging IS the disease
     P4 Constraint, not novelty, is the cure — the token contract is the deliverable
     P5 Decoration must carry a semantic anchor
     P6 Distance from your own last output is a first-class constraint

§3   Architecture and Roles                                         INFORMATIVE (+1 mermaid)
     Proposer(n) · Critic(n) · Facilitator(deterministic) · Scribe · Decider(exactly 1)
     Authority split: critics have no decision power; the Decider does not critique

§4   Phase State Machine                                            NORMATIVE ← the spine
     4.1  Phases D0–D9 and the two terminal states
     4.2  Per-phase entry conditions, required artifacts, exit conditions
     4.3  Gates G0–G4 (each with the exact checker error string)
     4.4  Illegal transitions
     4.5  Track conditioning: fast | standard | deliberation (§4.0 inheritance)
     4.6  Re-entry rules (Phase-7 fix-up returns to D6 or D8, never to D1)
     4.7  Mermaid state diagram

§5   Canonical Artifacts                                            NORMATIVE ← typed vocabulary
     One identically-shaped entry per kind:
     H3 name · one-line purpose · rationale paragraph · required front matter ·
     required body sections · minimal example · who writes it · who consumes it
     5.1 BRIEF.md            5.2 DIRECTION-<agent>.md   5.3 HEATMAP.jsonl
     5.4 CRITIQUE-<agent>.md 5.5 SCORECARD.md           5.6 VERDICT.md
     5.7 design/ (the system: tokens + specs)           5.8 AUDIT.md
     5.9 LEDGER.md

§6   Token Contract                                                 NORMATIVE
     6.1  W3C DTCG 2025.10 adopted verbatim — Format, Color, Resolver
     6.2  Three tiers: primitive → semantic → component; reference direction is one-way
     6.3  Required groups; colors authored in colorSpace "oklch" with 6-digit hex fallback
     6.4  Modes via Resolver (light+dark REQUIRED; hc, density, reduced-motion RECOMMENDED)
     6.5  Snapshot vs RFC-6902 patch for token evolution
     6.6  Extension: x- prefixed groups; core: reserved; unknown groups MUST NOT error

§7   Collaboration Contract Rules                                   NORMATIVE (numbered + bold-named)
      1. **Independent divergence.**   6. **Bounded grafting.**
      2. **Declared axes.**            7. **Traceability.**
      3. **Distinctness.**             8. **Recusal.**
      4. **One critique round.**       9. **Advisory votes.**
      5. **No averaging.**            10. **Declared degradation.**
                                      11. **Counter-signed waivers.**
                                      12. **Abandonment.**

§8   Rule Catalog                                                   NORMATIVE (pointer + citation form)
     8.1  How to cite a rule in a critique (`core:<slug>` + evidence tier + location)
     8.2  Two classes and their different burdens of proof (slop vs quality)
     8.3  Severity: 0–4 Nielsen anchors, re-anchored for design violations
     8.4  Rule stand-down discipline (one defect → exactly one finding)
     8.5  → references/RULES.md

§9   Target Profiles and Capability Declaration                     NORMATIVE
     Discovery only, no negotiation. Absent = undeclared. `custom` escape hatch.
     css-vars | tailwind | swiftui | android-compose | email-html | terminal-tui | print
     Each profile declares which rules are checkable and which are N/A.

§10  Evidence Tiers and Verdicts                                    NORMATIVE
     10.1 Tiers: text-regex | css-parse | dom | browser | screenshot | human
     10.2 Verdicts: pass | violation | needs-review | unjudgeable
     10.3 "unjudgeable: <tier>" is COMPLIANT; a silent skip is not
     10.4 The degradation banner (a silent degraded review is a failed review)

§11  Conformance                                                    NORMATIVE
     11.1 Levels L1–L4        11.2 Fixtures       11.3 The reference runner

§12  Extension Points and Reserved Names                            NORMATIVE
§13  Versioning and Deprecation Policy                              NORMATIVE
§14  Draft Proposals                                                INFORMATIVE (5-state ladder)
§15  Worked Examples                                                INFORMATIVE (2 complete runs: 1 fast, 1 deliberation)
§16  Changelog                                                      append-only, one entry per spec version
```

**The load-bearing inversion vs AG-UI:** `VERDICT.outcome` is a discriminated union
`{ kind: "winner", direction: <id>, grafts: [0..3] } | { kind: "rumble", directions: [2], decided_by: <external evidence> } | { kind: "abandoned", reason: <string> }`.
There is no shape in that union that can express "the average of four directions". That is what makes P3 mechanical rather than exhortative.

---

# 2. Canonical artifact set

## 2.1 Design rules

1. Every artifact is either **machine-checkable** or **explicitly marked human-judgement**. No third category.
2. Artifacts live under the existing `parley-deck/ideas/<slug>/` so they inherit the deck's ownership and consensus rules. No new state tree, no `.parley-design/`, no gitignored mutable log.
3. **Tokens are the only source of truth for values.** Every prose artifact references tokens; none of them restates a value.
4. Artifact count is bounded and track-conditioned. `fast` produces 4 files; `standard` 9; `deliberation` 12.

## 2.2 The set

```
parley-deck/ideas/<slug>/
├── BRIEF.md                                 D0
├── directions/
│   ├── DIRECTION-<agent>.md                 D1  (one per proposer)
│   └── DIRECTION-<agent>.tokens.json        D1  (minimal DTCG proof-of-realness)
├── HEATMAP.jsonl                            D2
├── critiques/CRITIQUE-<agent>.md            D3
├── REBUTTAL-<agent>.md                      D3  (optional, misreadings only)
├── SCORECARD.md                             D4
├── VERDICT.md                               D5  (+ RUMBLE.md on the D5b branch)
├── design/                                  D7  ← the ratified system
│   ├── tokens/primitives.tokens.json
│   ├── tokens/semantic.tokens.json
│   ├── tokens/component.tokens.json
│   ├── tokens/resolver.json
│   ├── DESIGN-SYSTEM.md                     (Foundations index + Named Rules + MUST-share/MAY-differ)
│   ├── CONTRAST-MATRIX.md                   (generated, every text-on-bg pair, computed)
│   └── components/<component>.spec.md
├── design-check.config.json                 D7  ← the enforcement contract
├── AUDIT.md / audit.json                    D9  (machine-written)
└── LEDGER.md                                D9  (the decision record)
```

## 2.3 Artifact table

`MC` = machine-checkable. `T` = minimum track that requires it.

| # | File | Purpose | Required front matter | MC | How checked | T |
|---|---|---|---|---|---|---|
| 1 | `BRIEF.md` | The problem, audience, goals critique is judged against, hard constraints, **anti-goals**, **anti-references**, **`divergence_axes`**, `target_profiles`, `genre` | `spec`, `track`, `divergence_axes[]`, `target_profiles[]`, `genre` | structure | required headings; ≥3 anti-goals; ≥1 anti-reference; ≥3 divergence axes (gate **G0**) | fast |
| 2 | `DIRECTION-<agent>.md` | ONE committed visual direction: **Signature** (the single element it will be remembered by), thesis, position on **every** declared axis, the one risky move, what it deliberately sacrifices, `critique_requests[]` | `spec`, `handle` (one word, not the agent name), `axes{}`, `signature`, `word_budget_used` | structure + divergence | required headings; hard word cap enforced by truncation; axis coverage; **cross-direction distance** (gate **G1**) | standard |
| 3 | `DIRECTION-<agent>.tokens.json` | Minimal DTCG set proving the direction is real: ≥1 type ramp, ≥1 color ramp, radius, space | DTCG `$schema` | **full** | `designtokens.org/schemas/2025.10/format.json` | standard |
| 4 | `HEATMAP.jsonl` | Part-level `like` / `concern` marks on **other** agents' directions, committed before reveal | per-record `spec` | **full** | JSONL schema; recusal (no self-marks); budget 20–30 marks | standard |
| 5 | `CRITIQUE-<agent>.md` | Typed critique records against a named target + assigned lens | `spec`, `lens`, `targets[]` | **full** | every record `{target, part, class: like\|wish\|what_if, severity 0–4, tied_to_goal, rule_id?, evidence_tier, fix?}`; a `wish` without `tied_to_goal` is dropped as taste | standard |
| 6 | `SCORECARD.md` | Absolute (never pairwise) scores on the 4 weighted criteria + advisory straw poll | `spec`, `aggregation: trimmed-mean` | **full** | recusal; every direction scored by every non-author; trim high+low then mean | deliberation |
| 7 | `VERDICT.md` | **Exactly one winner**, 0–3 typed grafts, every losing direction listed with why it lost and marked `maybe-later`, override reason if the Decider ignored the poll | `spec`, `outcome{}` (discriminated union), `decider` | **full** | exactly 1 winner; grafts ≤3 each with `from:`, `part:`, `re-expressed-as:` (a winner token); all non-winners present in the rejected list | fast |
| 8 | `design/tokens/primitives.tokens.json` | Raw values only. **All colors `colorSpace: "oklch"`** + 6-digit hex fallback. No semantic names | DTCG | **full** | schema + colorspace + required-type + duplicate-value (ΔE2000) | standard |
| 9 | `design/tokens/semantic.tokens.json` | Roles only (`color.text`, `color.bg.surface`, …), Atlassian `foundation.property.modifier`, modifier omitted for defaults. **Every `$value` MUST be an alias** | DTCG | **full** | schema + "no literal `$value` in this file" + `$description` required | standard |
| 10 | `design/tokens/component.tokens.json` | Per-component tokens, aliasing **semantic only** | DTCG | **full** | schema + alias-direction graph assertion (never skips a tier, never reverses) | deliberation |
| 11 | `design/tokens/resolver.json` | Modes. `light` + `dark` REQUIRED; `hc`, `density`, `reduced-motion` RECOMMENDED | DTCG Resolver, `version: "2025.10"` | **full** | `schemas/2025.10/resolver.json`; every semantic color has a value in every declared context | standard |
| 12 | `design/DESIGN-SYSTEM.md` | Human-readable system. **Named Rules** (`**The Two-Face Rule.** <one forceful sentence>`), the Foundations index, the type-scale ratio **and** the rounded/pruned final table, the declared spacing scale as an explicit member list, motion stance, voice, and the **`## What surfaces MUST share` / `## What surfaces MAY differ on`** pair | `spec`, `status: draft\|preview\|stable`, `registry-digest` | structure + cross-ref | required headings; every value present must resolve to a token; line-height % 4px == 0; no un-rounded ratio output | standard |
| 13 | `design/CONTRAST-MATRIX.md` | Generated: every text-on-bg and non-text pair with its computed ratio | `generated-by`, `generated-at` | **full** | recomputed by the checker; fail <4.5:1 body / <3:1 large / <3:1 non-text | standard |
| 14 | `design/components/<c>.spec.md` | Merged Carbon+Polaris+GOV.UK component spec (§2.4) | `spec`, `status` (Carbon 4-value enum) | structure + tokens | required sections; all 10 state rows; Mouse+Keyboard subsections; `When NOT to use` names an alternative; every token resolves; **no literal hex, no bare px** | deliberation |
| 15 | `design-check.config.json` | The enforcement contract: token globs, target profiles, enabled rule ids, severities, thresholds, waivers | own JSON Schema, `implements: PDS/1.0` | **is the input** | its own published schema | fast |
| 16 | `AUDIT.md` + `audit.json` | Checker output. Fixed skeleton (§5.5), degradation banner, per-rule verdicts | `spec`, `registry`, `registry-digest`, `tiers_reached[]` | **full** | stable versioned schema | fast |
| 17 | `LEDGER.md` | Who proposed what, who objected, what the objection was, why the winner won, which grafts survived, `HOMOGENIZATION_WARNING` flags, waiver counter-signatures, content hashes | `spec`, `run_id` | structure | every phase artifact hashed; every override has a reason | standard |

## 2.4 `<component>.spec.md` required sections

```
# <Component>                          status: draft | preview candidate | preview | stable
## Purpose                             (1 sentence)
## When to use                         (bullets)
## When NOT to use                     (bullets; MUST name the alternative)
## Anatomy                             (numbered parts)
## Variants                            (table: name | when | token deltas)
## Sizes                               (table)
## States                              (10 rows: default, hover, focus-visible, active, selected,
                                        disabled, read-only, loading, error, success —
                                        every cell references a TOKEN, never a literal)
## Behaviors                           (responsive, overflow/reflow, truncation, expansion,
                                        scroll, empty, long-content)
## Interactions                        (### Mouse  ### Keyboard  ### Touch)
## Content guidelines                  (label rules, char budgets, capitalisation)
## Accessibility                       (role, name, ARIA, focus order, SR announcement,
                                        contrast results, target size)
## Do / Don't                          (≥3 pairs, one sentence each)
## Tokens used                         (explicit list; MUST all resolve)
## Related
```

## 2.5 Reconciliation notes (why these and not the analysts' longer lists)

- **Collapsed** `FOUNDATIONS.md` / `TYPE-SCALE.md` / `COLOR-SYSTEM.md` / `SPACING-LAYOUT.md` / `MOTION.md` / `VOICE-AND-CONTENT.md` / `PATTERNS.md` / `ACCESSIBILITY.md` / `DESIGN-DECISIONS.md` (9 files) into **`DESIGN-SYSTEM.md` + `LEDGER.md`** (2). Rationale: the wave-1 anti-goal "file sprawl" is real under N agents, and 8 of those 9 files would be prose restating token values. The ADR log's real job (≥2 alternatives per decision) is already discharged by `VERDICT.md` + `LEDGER.md`, which record four *authored* alternatives, not two hypothetical ones.
- **Kept** `CONTRAST-MATRIX.md` as a separate generated artifact because it is the one place a claim ("we pass AA") becomes a recomputable table.
- **Dropped** `EXHIBIT.md` as a file: anonymisation + slug assignment + order randomisation are deterministic facilitator work, recorded in `LEDGER.md`, not a document.
- **Adopted verbatim, zero invention:** DTCG `2025.10` for tokens (stable since 28 Oct 2025, official JSON Schemas live), Carbon's 4-value status enum, Carbon's six design-spec requirement rows, GOV.UK's five criteria (Useful/Unique → propose; Usable/Consistent/Versatile → publish) as the D5 gate questions, Atlassian's `foundation.property.modifier` naming, Material's three-tier layering (without the `md.*` names).
- **Rejected:** Atomic Design as the taxonomy (Frost: the labels "have never been the point"; none of the five major systems use it); `tr.designtokens.org` as a citation (301s); `@terrazzo/plugin-lint-*` (does not exist on npm); Carbon's "match the spec perfectly down to the pixel" (unfalsifiable) — replaced by "every implemented value resolves to a declared token".

---

# 3. The slop registry

The shared vocabulary between the two skills. **This table is the interface.**

## 3.1 Registry design rules

1. **Two classes, different burdens of proof** (impeccable's split, independently reached by the ai-slop analyst):
   - `SLOP` — "this reads as machine-generated". Arguable. **Requires quorum to block.** Warn + justification by default.
   - `QUALITY` — "this is broken". **A single agent MAY BLOCK.** Hard-fails.
2. **Every rule carries a `tier`** (cheapest engine that can decide it): `text-regex` · `css-parse` · `dom` · `browser` · `screenshot` · `human`. An agent that cannot reach a rule's tier writes `unjudgeable: <tier>` and is **compliant**.
3. **Rules are conjunctive.** Never gate on a bare aesthetic fact. The purple rule fires only on a *filled* accent; the Inter rule is `Centered + single-face`, not `Inter`.
4. **Gate on count, not on any single rule.** Empirically anchored (Krebs, N=1,590): 0–1 tells = 46 % of the wild population; 2–3 = 32 %; **4+ = 22 % = "high"**. Default FAIL threshold for the `SLOP` class = **≥4 concurrent tells**.
5. **Prohibitions of defaults, never prescriptions of values.** "Do not ship the model's first palette unmodified" survives; "use Fraunces" does not — Fraunces went from *the fix* (dev.to, 2025) to *a tell* (Krebs `slop-fonts.js`, 2026) inside a year.
6. **Every rule carries `added:`, `sources:`, `confidence: confirmed|strong|weak`, `status: active|draft|deprecated`.** The registry is dated and versioned; a frozen anti-slop list is a slop generator on a delay.
7. **Stand-down discipline.** Each defect has exactly one owning rule id; the rule's metadata names who it yields to. Stops five agents raising the same objection under five names.
8. **At least one rule class is deliberately design-system-blind** (impeccable's `undersized-ui-text` insight): widening the token ramp launders the token, not the problem.

## 3.2 Tier 1 — confirmed slop tells (≥4 independent sources)

| Rank | `core:` id | Tell | Tier | Firing rule (conjunctive) | Src |
|---|---|---|---|---|---|
| 1 | `color-accent-purple` | Indigo/violet as the **filled** accent on CTAs/links | css-parse + browser | filled accent (`bg.a ≥ 0.5` or gradient fill) on button/link **AND** hue ∈ [260,310]. Ghost/outline/decorative purple explicitly excluded | 8 |
| 2 | `type-single-face-centered` | One face, centred hero | browser (css-parse partial) | exactly 1 non-generic primary family (≥20 text elements) **AND** centred hero axis | 8 |
| 3 | `color-gradient-hero` | Purple→blue / cyan→magenta gradient | text-regex + css-parse | 2-stop gradient crossing hue 250–330 on hero bg or button fill | 7 |
| 4 | `layout-hero-fullviewport-centered` | Full-viewport centred hero | css-parse + dom | `min-height:100vh\|100dvh` **AND** eyebrow+title+lede+CTA share one centred vertical axis | 6 |
| 5 | `layout-3col-icon-cards` | 3 equal cols, icon-above-heading, equal heights | dom | 3 equal-width tracks **AND** icon-tile immediately above heading **AND** identical card heights | 6 |
| 6 | `card-side-stripe` | Thick coloured stripe on one card edge | css-parse + browser | dominant single-edge chromatic border ≥2px with `radius>0` (else ≥3px), **or** a 3–12px full-height chromatic `::before/::after`, **or** an inset box-shadow of the same geometry | 5 |
| 7 | `type-gradient-text` | Gradient headline | text-regex | `background-clip:text` (+`-webkit-`) **AND** a gradient in the same element's `background-image`. **No target profile permits this** | 5 |
| 8 | `shape-uniform-extreme-radius` | Uniform extreme radius | css-parse | `rounded-2xl\|3xl` or a single radius value ≥16px on ≥N card-like elements with zero variation | 5 |
| 9 | `type-hero-eyebrow-chip` | Pill badge / eyebrow directly above H1 | dom | H1 ≥48px, previous sibling non-heading, 2–60 chars, ≤14px, **AND** (tracked-caps ≥1.6px \| accent-bold ≥700 \| dash-prefix pseudo) | 4 |
| 10 | `fx-glass-decorative` | Decorative glassmorphism | browser | `backdrop-filter: blur()` on a translucent floating panel that carries no depth/layering function | 4 |
| 11 | `type-italic-accent-word` | Serif-italic accent word in an upright headline | css-parse | `font-style: italic` on `h1–h6` (or `<em>` inside a heading) **AND** serif family or generic `serif` fallback | 4 |
| 12 | `fx-colored-glow` | Saturated coloured box-shadow glow | css-parse | chromatic shadow (`chroma ≥ 30`), blur > 4, **AND** (zero offset **or** page luminance < 0.1) | 4 |

## 3.3 Tier 2 — strong slop tells (2–3 sources)

| `core:` id | Tell | Tier | Firing rule |
|---|---|---|---|
| `color-pure-neutral` | Pure `#000`/`#fff`, zero-chroma greys | text-regex | any `#000000`/`#ffffff` as base surface, or `oklch(L 0 H)` — min chroma **0.005** |
| `color-cream-default` | Reflex "tasteful" cream page | css-parse | page bg: `min(r,g,b) ≥ 209` **AND** `r ≥ g ≥ b` **AND** `6 ≤ r−b ≤ 48` |
| `layout-nested-cards` | Card inside a card | dom | card-like element with a card-like ancestor; only the innermost reports |
| `type-overused-font` | Monoculture face as primary | css-parse + browser | family ∈ registry list **AND** covers ≥15 % of text-bearing elements. **The list is data with `added:` dates, never welded into prose** |
| `type-flat-hierarchy` | No hierarchy | css-parse | ≥3 distinct sizes **AND** `max/min < 2.0` |
| `type-allcaps-body` | All-caps body / wide tracking on body | css-parse | non-heading, >30 chars, `text-transform: uppercase`; or tracking > 0.05em on non-caps body |
| `type-unrounded-ratio` | Raw `1.25^n` output | css-parse | any type-scale value with >1 decimal place (19.2px, 23.04px, 27.65px) — no mature system ships this |
| `space-monotonous` | Every gap identical | css-parse | ≥10 spacing samples **AND** dominant value > 60 % **AND** ≤3 unique values |
| `icon-emoji-as-ui` | ✨🚀⚡🔥🎯✅ as feature/nav/step icons | text-regex | emoji in an icon slot |
| `icon-library-mixing` | ≥2 icon libraries on one surface | text-regex | imports from ≥2 of lucide/heroicons/phosphor/material/fa |
| `layout-numbered-steps` | `01 / … 02 / …` section labels | dom | ≥2 numbered labels with ≥2 distinct indices, ≤13px, deliberately styled |
| `proof-stat-bar` | "10K+ users · 99.9 % uptime · 4.9★" | dom + text-regex | ≥3 stat items in one row above the fold |
| `bg-aurora-orb` | Aurora blob / mesh / radial halo / floating orbs | css-parse + screenshot | radial/mesh gradient with >1 accent hue, or footprint > ~5 % of viewport, or animated page-wide |
| `bg-gridlines` | Decorative graph-paper background | text-regex | ≥2 hairline stops **AND** a px tiling cell in the same declaration block |
| `motion-transition-all` | `transition: all` / `transition-all` | text-regex | literal |
| `motion-uniform-hover-scale` | `hover:scale-105` everywhere | text-regex | same scale transform on ≥3 unrelated selectors |
| `motion-bounce-easing` | Bounce/elastic on UI | css-parse | `cubic-bezier` with `y1` or `y2` outside `[-0.1, 1.1]`, or `animate-bounce` |
| `motion-scroll-fade-all` | Fade-up-on-scroll on every section | dom + text-regex | reveal animation on ≥N sequential sections |
| `chrome-ai-nav` | Wordmark-left + 4–5 links + CTA-right + hairline + white | dom | the full conjunction. Rationale to quote: *the shape is genre-blind* |
| `chrome-ai-footer` | 4 link columns Product·Company·Resources·Legal + social row | dom | the full conjunction |
| `chrome-redrawn-ui` | Fake browser bar / phone notch / terminal frame | dom (heuristic) | **needs-review only, never a gate** — a real device frame is legitimate |
| `copy-buzzword` | Marketing buzzwords | text-regex | ≥1 of the 30-phrase list |
| `copy-cliche-opener` | "Built for the modern team", "Unleash your…" | text-regex | the 10-phrase table |
| `copy-aphoristic-cadence` | "Not a X. Y." manufactured contrast | text-regex | ≥3 matches |
| `copy-placeholder-names` | Jane Doe / Acme / Nexus / Lorem | text-regex | literal list |
| `copy-emdash-density` | Em-dash **overuse** | text-regex | `count ≥ 8` **AND** `bodyText.length ≤ count × 500`. **Never ban the em-dash itself** |
| `type-oversized-h1` | Display headline carrying 100 chars | browser | `≥72px` **AND** `≥40 chars` **AND** dominates ≥25 % of viewport area |
| `layout-faq-appendix` | Bolted-on FAQ accordion | dom | advisory only — highest legitimate-use rate in the corpus |

## 3.4 Token / system rules (`QUALITY`, the highest-yield band)

| `core:` id | Rule | Tier | Verdict |
|---|---|---|---|
| `token-raw-literal` | Any `#hex` / `rgb()` / `hsl()` / `oklch()` / bare `font-family` / bare px outside the token layer | text-regex + css-parse | **violation** — the single highest-yield mechanical rule that exists |
| `token-offscale-space` | Any padding/margin/gap/radius not a **member of the declared scale** | css-parse | violation. **Assert membership, never `value % 8 == 0`** — GOV.UK ships a 5px grid, Carbon a 2px base |
| `token-offscale-typesize` | Font size not a step in the declared ramp (±0.5px) | css-parse | violation; abstain entirely if the system declares only fluid `clamp()` endpoints |
| `token-duplicate-color` | Two tokens that are the same colour | css-parse | **ΔE2000 < 1.0 = violation; 1.0–2.3 = needs-review.** Not OKLCH ΔL (`#0a0a0a` vs `#000` = 14 % ΔL but ΔE 1.59) |
| `token-alias-direction` | component → semantic → primitive only; never skip, never reverse | json | violation (graph assertion) |
| `token-schema-invalid` | DTCG `2025.10` schema failure | json | violation — pure schema, zero heuristics, the hardest gate we ship |
| `token-mode-missing` | A semantic colour with no value in a declared resolver context | json | violation |
| `token-no-description` | Semantic token without `$description` | json | recommendation — forces stated intent |
| `token-phantom` / `token-unused` | `var(--x)` never declared / declared never used | css-parse + text-regex | needs-review |
| `type-lineheight-off-grid` | `line-height % 4px != 0` | css-parse | recommendation. The one *arithmetic* rule that is universal in practice (all 15 M3 styles, all 8 Polaris line-heights) |
| `shadow-ramp-unsystematic` | Shadow tokens not sharing one hue, or non-monotonic alpha, or >7 steps | css-parse | recommendation |

## 3.5 Defect rules (`QUALITY`, single-agent BLOCK)

| `core:` id | Rule | Tier | Note |
|---|---|---|---|
| `a11y-contrast-text` | 4.5:1 body / 3:1 large (≥18pt or 14pt bold) | browser | the flagship check — reproduced catching `#5c3ef6` on `#5b3df5` at **1.01** |
| `a11y-contrast-nontext` | 3:1 for UI components and graphical objects (SC 1.4.11) | browser | catches invisible borders, ghost icons, low-contrast focus rings |
| `a11y-contrast-indeterminate` | axe returned `incomplete` ("background gradient") | browser | **MUST emit `needs-review`, never `pass`.** Treating incomplete as pass is how the bug ships |
| `a11y-target-size` | 24×24 CSS px, **our own rule, no spacing exception** | browser | axe's `target-size` passed an 18×18 button via the WCAG Spacing exception |
| `a11y-focus-visible-missing` | interactive element with no `:focus-visible` | css-parse (declared) / browser (effective) | needs-review at css-parse, violation at browser |
| `a11y-focus-ring-animated` | focus indicator animates in | css-parse | violation — must show instantly |
| `a11y-reduced-motion-missing` | transform/keyframe with no `prefers-reduced-motion` fallback | css-parse | violation |
| `a11y-heading-skip` | `level > prev + 1` | dom | violation |
| `a11y-carousel-nopause` | auto-rotating content without pause (SC 2.2.2) | dom | violation |
| `resp-horizontal-scroll` | `scrollWidth > clientWidth` at 320/375/414/768/1280/1920 | browser | violation. Fix is `overflow-x: clip` on **both** `html` and `body` — `clip`, not `hidden` |
| `resp-two-line-clickable` | `a`/`button`/nav label wraps (`getClientRects().length > 1`) | browser | violation |
| `resp-reflow-320` | two-dimensional scrolling at 320px (SC 1.4.10) | browser | violation |
| `layout-grid-1fr-image` | bare `1fr` track containing an image | css-parse + dom | violation. One-character fix: `1fr` → `minmax(0, 1fr)` |
| `state-incomplete` | any of the 8 states missing on a production element | dom + css-parse | violation |
| `state-input-discipline` | 5 sub-checks: border-width shifts between states · focus ring built from `border` not `outline` · input height ≠ button height (44px floor) · helper-text slot collapses when empty (`min-height: 1lh`) · disabled signalled by opacity alone | browser | violation |
| `defect-script-error` | uncaught error on load | browser | violation |
| `defect-broken-image` | `<img>` with no/empty `src` | dom | violation |
| `defect-hidden-at-rest` | >30 % of ≥200 chars invisible at rest after a reveal sweep | browser | violation |
| `defect-text-occlusion` | text overlapped/clipped by an opaque box | browser | violation |
| `defect-text-overflow` | `scrollWidth − clientWidth ≥ 16` outside a scroll region | browser | needs-review (intentional ellipsis) |

## 3.6 Human-only rules (`parley-design`'s exclusive territory)

| `core:` id | Rule | Why not mechanical |
|---|---|---|
| `honesty-invented-metric` | Any quantitative claim the user did not supply ("10× faster", "trusted by 50,000+", "+47 % conversion") | The *number shape* is regex-able; its *truth* needs the brief. **Hard-fail, non-waivable.** Three sanctioned fixes, in order: replace with `—` + a labelled block; ask and pause; rebuild without the proof slot |
| `honesty-fabricated-social-proof` | Invented testimonials, logos, awards | needs ground truth |
| `pov-absent` | The artifact would fit any brief | This is the actual disease. Only judgement sees it |
| `decor-unmotivated` | Ornament with no semantic anchor in the content | requires knowing what the content means |
| `img-ai-look` | Smooth-mesh-blob figures, symmetric default lighting, corporate-doodle person | needs vision + taste |
| `struct-generic-template` | Hero → 3 features → CTA → footer | detectable in the limit; **any DOM rule over-fires on legitimate pages.** `needs-review` signal to the agents, never a gate |
| `self-repetition` | This surface shares a structural fingerprint with the project's own prior output | **Our clearest differentiator.** Score **structural distance, not visual distance — colour-swaps do not count as variety.** Computed against `LEDGER.md`, not against a global registry |
| `direction-collapse` (gate G1) | Two DIRECTIONs share all declared axes | process rule, computed over artifacts, not over rendered UI |

## 3.7 Profile exemptions (mandatory)

A `genre:` / `profile:` field **disables the entire `SLOP` class while keeping every `QUALITY` rule on**, for: enterprise CRUD / admin / internal tools · regulated & safety UI · accessibility-driven forms · platform-native HIG/Material surfaces · checkout & auth · documentation & reference. In these classes convention is the correct answer and distinctiveness is the defect.

---

# 4. The ritual

## 4.1 Provenance — this is not invented

Our stated method is Google Ventures' Design Sprint **"Sticky Decision"** almost move-for-move: Art Museum (anonymous wall) → Heat Map (20–30 dots, silent, part-level) → Speed Critique (facilitator narrates, **author stays silent until the end**) → Straw Poll (1 vote, advisory) → **Supervote (the Decider gets 3 votes and MAY ignore the poll)**. Adopt the names verbatim — named rituals are memorable and auditable.

Four things the research forces us to **change** from the method as stated:

1. **Cap adversarial critique at ONE round.** More debate destroys the premise: factual attrition + stance homogenization across rounds (arXiv 2606.03032); diversity collapse from structural coupling (arXiv 2604.18005). A second round requires an explicit Decider instruction and a logged reason.
2. **Never decide by agent majority vote.** Idea *selection* is the documented failure: selectors show "a strong tendency to select feasible and desirable ideas, **at the cost of originality**" (Rietzschel/Nijstad/Stroebe 2010, Br. J. Psych.). That is the slop mechanism wearing a ballot. Votes are advisory; one Decider (the human by default) decides.
3. **Add the RUMBLE branch.** GV does not say "one winner whole, always". It branches: **conflicting winners → Rumble** (build both cheap, fake-brand them, decide by external evidence); **compatible winners → All-in-One** (our winner+grafts path). Without the branch we force a premature pick on genuinely incommensurable directions.
4. **Add a distinctness gate before critique.** Without it, four "different directions" will quietly be one direction and every downstream ritual is theatre. Critiquing a collapsed set **launders the collapse into a "consensus"**.

And one claim we must **not** make: *"multi-agent improves taste."* MAD underperforms plain single-agent on every benchmark where it has been measured (MMLU: CoT 80.73 / SC 82.13 / **MAD 74.73**), and all those benchmarks were factual/code, never aesthetic. The defensible claim is: **multi-agent is a diversity generator; a rule does the selecting.** Human–human agreement on aesthetic choices is **38.34 %** — an aesthetic score can never be a hard gate, and `ABSTAIN` is a legitimate verdict (the published route to human parity was abstaining on ~35 % of cases).

## 4.2 Roles and invariants

**Roles.** `Proposer` (each roster agent) · `Critic` (each roster agent, on others' work only) · `Facilitator` (the driver; deterministic, no model call) · `Scribe` (the driver; writes `LEDGER.md`) · **`Decider` (exactly one; the human by default)**.

**Invariant I-0 — Authority split.** Critics MUST have no decision authority. The Decider MUST NOT act as a Critic. All votes and scores are advisory and MUST be labelled as such in the ledger. *(Pixar Braintrust has no authority; Apple DRI; GV "the decider makes the call, not the group".)*

**Invariant I-1 — Recusal.** No agent may score, rank, or vote for its own direction. Self-scores MUST be **discarded, not down-weighted**. *(Self-preference bias measured at −38 %…+90 %.)*

**Invariant I-2 — Absolute, not pairwise.** Scoring MUST be absolute against the anchored rubric. Pairwise "which is better, A or B?" is FORBIDDEN in the deciding phase. *(Position bias is strongest exactly at small quality gaps — our situation with four good directions.)*

**Invariant I-3 — Length normalisation.** DIRECTION artifacts obey a hard section/word cap. Over-cap artifacts are **truncated before scoring, not rewarded**. *(Judge verbosity correlation r = .87 vs .44 for humans.)*

**Invariant I-4 — Declared degradation.** No participant may sign a design verdict without declaring which evidence tiers it could not reach and which rules it marked `unjudgeable`. *A silent degraded review is a failed review.*

## 4.3 The phases

Budgets are **ratios of the phase allowance**, never wall-clock minutes — agents have no clock, and "3 minutes" will be hallucinated as compliance.

### D0 · BRIEF → `BRIEF.md`
Facilitator drafts, Decider ratifies, once, up front. MUST contain: the problem; the audience; the **business + user goals critique will be judged against**; hard constraints; `target_profiles[]`; `genre`; **anti-goals** (≥3) and **anti-references** (≥1); and the **`divergence_axes`** — the named axes on which directions MUST differ (e.g. `structure`, `typographic voice`, `colour strategy`, `motion posture`, `density`).
**Gate G0:** no `divergence_axes` → no D1. Without declared axes, "be different" is unenforceable and G1 is uncheckable.
*Parley mapping:* §9.0 preflight + round-01 kickoff. This is the **only** place a human is asked anything. Every later question becomes a labelled assumption or a Parley consult.

### D1 · DIVERGE (isolated) → `DIRECTION-<agent>.md` + `.tokens.json`
Each Proposer works **without reading any other agent's output**. Isolation is mechanically enforced (separate dirs; no cross-reads until D2 opens). This is not hygiene — it is the whole product: the moment agents see each other's drafts, structural coupling starts and four directions become one.

Four sub-steps, mirroring GV's four-step sketch:
- **D1.a Notes** — harvest from `BRIEF.md` and provided references. 20 %
- **D1.b Ideas** — private candidate list; circle the most promising. 20 %
- **D1.c Crazy-8** — **8 one-line variations of your own best idea**, no elaboration. 10 %. *(Intra-agent divergence, which the method was missing.)*
- **D1.d Direction Sketch** — the committed artifact. 50 %

`DIRECTION.md` MUST be **self-explanatory without its author present**, MUST carry a **one-word handle** (not the agent's name), MUST declare its position on **every** declared axis, MUST name its **Signature** — the single element this surface will be remembered by — and MUST carry a token table so it is checkable, not merely describable.
**Two directions with the same Signature are the same direction.** That is mechanically checkable and it is what forces them apart.

**Exit gate G1 — DISTINCTNESS.** The Facilitator (deterministic, no model call) strips author identity, assigns stable slugs `A/B/C/D`, and **randomises presentation order per reviewer**. Then: if any two directions match on **all** declared axes, or share a banned-slop signature (identical font strategy + identical colour strategy + identical macrostructure), **D1 FAILED**. Re-run D1 with **forced-distinct axis assignments** (each agent is *assigned* a different position on the primary axis). MUST NOT proceed to critique on a collapsed set.
Checker error string: `G1 core:direction-collapse — DIRECTION-A and DIRECTION-C share all 5 declared axes. Re-run D1 with forced-distinct assignments on axis 'colour strategy'.`
*Parley mapping:* round-01 artifacts. Convergence here is **not** a pass — it is a slop alarm, and it is the one signal only a heterogeneous roster can produce.

### D2 · HEAT MAP (silent, parallel, part-level) → `HEATMAP.jsonl`
Every Critic independently emits typed marks against **parts** of **other** agents' directions, 20–30 marks total:
```jsonl
{"direction":"C","part":"nav.sticky-condense","mark":"like","intensity":2,"note":"…"}
{"direction":"C","part":"hero.gradient","mark":"concern","severity":3,"rule_id":"core:color-gradient-hero","evidence_tier":"css-parse","note":"…"}
```
Reveal is **simultaneous** — no Critic sees another's marks until all are submitted. *(GV Note-and-Vote: "commit your vote to paper" before revealing kills bandwagoning.)*
`mark: "like"` records are **the graft harvest**. This phase exists as much to find graftable parts as to find flaws — and it is the mechanism the method was missing for *choosing* which details to graft.

### D3 · CRITIQUE (one round, assigned lenses) → `CRITIQUE-<agent>.md` (+ optional `REBUTTAL-<agent>.md`)
- The Facilitator narrates each direction from the artifact and **names the dot clusters first**.
- Each Critic is **assigned a lens** — Black (risk/failure), Yellow (value/upside), White (facts/constraints/a11y), Green (adjacent alternatives). **Assigned, not chosen**, because de Bono's justification is the strongest sentence in the corpus: the roles are *"blatantly artificial, a feature which helps separate individual ego from the activity."* An agent asked for its opinion defends its own proposal; an agent assigned the Black hat discharges a role.
- **The author is silent during critique of its own direction** and gets exactly one `REBUTTAL.md` afterwards, addressing **only misreadings**.
- Every critique entry is typed: `{target, part, class: like|wish|what_if, severity: 0–4, tied_to_goal, rule_id?, evidence_tier, fix?}`.
  - `severity` uses the **Nielsen 0–4 anchors** re-anchored for design violations (0 "not a problem" → 4 "ships as slop / catastrophe"). Only 4 (optionally 3) blocks.
  - **`tied_to_goal` is REQUIRED for `wish`.** A `wish` with no goal link is dropped as taste. *(GV: "Good feedback is about how the design is meeting or missing the customer and business goals.")*
  - **`fix` is OPTIONAL and explicitly NON-BINDING.** This resolves the Pixar contradiction: diagnosis is owed (Braintrust "notes, not prescriptions"), prescription is not owned (GV "suggestions, not mandates"). Plussing as a blanket rule is rejected — in an LLM setting it pushes critics to *author* the fix, which is how a critique quietly becomes an unratified redesign.
- **HARD CAP: ONE round.** A second requires an explicit Decider instruction and a logged reason.
- The Facilitator MUST log **stance diversity before vs after** and flag `HOMOGENIZATION_WARNING` if it dropped below threshold.
*Parley mapping:* round-02 cross-review. `class`/`severity` map onto existing Parley dispositions; align the vocabularies deliberately rather than inventing a third.

### D4 · SCORE + STRAW POLL (advisory) → `SCORECARD.md`
Each Critic scores each **other** direction absolutely (I-1, I-2) on an anchored, weighted rubric. Weights derived from Awwwards' shape but **re-weighted to counteract feasibility bias**:

| Criterion | Weight | Anchored at 0 / 3 / 5 / 8 / 10 by worked examples |
|---|---|---|
| **Distinctiveness** (avoids the aligned attractor) | 30 | required |
| **Systemic coherence** (tokens actually generate the UI; no orphan values) | 25 | required |
| **Fitness to brief** (goals named in `BRIEF.md`) | 25 | required |
| **Craft & accessibility** (mechanical: contrast, targets, focus, scale conformance) | 20 | required |

Aggregation: **drop the highest and lowest score per direction, then mean** (the Awwwards trim, scaled to a 4-agent jury). Every criterion ships with **one worked anchor example per level** in the skill's references — the cheapest single accuracy improvement available, and the gap Awwwards itself leaves.
Then a **Straw Poll**: one vote each, committed before reveal, one paragraph of reason, **explicitly advisory**. `ABSTAIN` is legitimate and MUST be preserved, never coerced into a vote.
*Track note:* `fast` skips D4 entirely; the Decider reads `HEATMAP.jsonl` and decides.

### D5 · DECIDE (SUPERVOTE) → `VERDICT.md`
The **Decider** holds three supervotes and **MAY ignore the scorecard and the straw poll entirely**. `VERDICT.md` records the winning direction, every losing direction marked **`maybe-later` (retained, never deleted)**, and — if the poll was overridden — a one-line reason.

**Rule 5 — No averaging (verbatim, binding).**
> *Exactly one direction wins whole. A synthesis of two directions' visual systems is a protocol violation. 0–3 named details MAY be grafted from losing directions; each graft MUST name its source direction, the exact part it carries, and the winner token it is re-expressed in.*

The typed `outcome` union is what makes this mechanical: there is no shape in it that can express an average.

The five **GOV.UK gate questions** are asked here, verbatim: *Useful* and *Unique* to have proposed; *Usable*, *Consistent*, *Versatile* to publish.

### D5b · RUMBLE (branch) → `RUMBLE.md`
If the Decider judges the top two directions **genuinely incommensurable** (they answer the brief with conflicting *premises*, not conflicting details), the Decider MAY declare a Rumble: build both at cheap comparable fidelity, give each a **distinct fake handle** so they are not read as "version A / version B", and defer to **external evidence** (user test, stakeholder, metric).
Default remains All-in-One: *"If you think you can combine your winning sketches into one product, don't bother with a rumble."* A Rumble MUST be rare and MUST be justified in writing — it doubles cost.

### D6 · GRAFT (bounded, from the heat map) → section of `VERDICT.md`
The winning direction's author is the DRI for the graft (Braintrust: the director owns the fix).
- Graft candidates MUST come from `HEATMAP.jsonl` `like` clusters on **losing** directions — **not** from fresh invention in this phase.
- **Maximum 3.** Each MUST be a discrete, nameable detail — an interaction, a component treatment, a copy device, a motion rule.
- **A graft MUST NEVER be a token-system layer.** Never a colour ramp, never a type scale, never a grid. *Grafting a system layer is how you get a camel.*
- A graft that cannot be re-expressed in the winner's tokens is **rejected**.
- **Exit gate G2 — COHERENCE.** Re-run the mechanical checks. Any new orphan token, off-scale value, or contrast failure fails **the graft**, not the winner.
*Parley mapping:* D4+D5+D6 together produce the design half of `FINAL.md`. `VERDICT.md` is either embedded in or referenced by `FINAL.md` and is what binds Phase-5 implementers.

### D7 · SYSTEMATIZE → `design/`
The winner's author writes the token layers and component specs. Everything the system asserts must resolve to a token; `CONTRAST-MATRIX.md` is generated, not claimed.
**Gate G3:** `design-check --level L3` clean. Checker error strings name rule id, violation and remedy.

### D8 · APPLY → implementation
Ordinary Parley Phase-5 implementers obey the ratified system. The `## What surfaces MUST share` / `## What surfaces MAY differ on` pair in `DESIGN-SYSTEM.md` **is the parallelisation contract** for implementers working on different screens.
Precedence chain, stated once and binding: **QUALITY rules > the ratified design system > the brief > parity with existing code > model habit.** The top rules are explicitly **not bypassable** by "preserve structural parity" / "mirror this reference" / "match the prior build" instructions.

### D9 · AUDIT → `AUDIT.md` + `LEDGER.md` → **RATIFIED | ABANDONED**
The checker runs at the highest tier available and **declares which tiers it reached**. Agents adjudicate `needs-review` findings; a waiver requires a counter-signature.
**Gate G4:** no RATIFIED signature while any `QUALITY` violation or any unanswered CRITIQUE is open.
`LEDGER.md` is the artifact the *next* design run reads, and the artifact §13 retro interrogates: *did the direction that won actually ship, and did the grafts survive?*

## 4.4 Track conditioning (inherits COOPERATION.md §4.0)

| Track | Phases run | Proposers | Notes |
|---|---|---|---|
| `fast` | D0 → D1 → D5 → D7 → D9 | 1–2 | No heat map, no critique round, no scorecard. A human can judge a hero in three seconds; running five agents to decide whether it is centred is absurd |
| `standard` | D0 → D1 → D2 → D3 → D5 → D6 → D7 → D8 → D9 | roster | The default |
| `deliberation` | all, + D4, + D5b available, + component specs | roster | Reserved for a system that will bind many surfaces |

## 4.5 Where multi-agent does **not** help — state this in the skill

1. **Taste by committee produces mush.** The mean of five directions is definitionally closer to the training mode than any single committed direction. The merge operator MUST be *choose*, never *average*.
2. **Quorum ratchets toward safe.** Five reviewers each holding a veto converge on the category standard. This is exactly why `QUALITY` gets a single-agent veto and `SLOP` requires quorum.
3. **Five blind reviewers are five times as confident and exactly as wrong.** None of the roster renders anything by default. Agreement among five agents that never saw the page is the single most dangerous artifact this skill could produce — which is what I-4 and the evidence tiers exist to prevent.
4. **Consistency work wants fewer voices.** At *direction* time and *review* time multi-agent adds value; at *coherence* time it is a liability. One ratified system, N obedient implementers.
5. **Mechanically checkable things want a tool, not a debate.** Contrast, token drift, banned values — check once, report once, never argue five times.

---

# 5. Split line between the two skills

## 5.1 The boundary, in one sentence

**`parley-design` owns everything that requires a decision; `parley-design-check` owns everything that can be decided by a script with no model in the loop.** The moment the checker emits "this feels generic" it becomes an unfalsifiable LLM oracle and teams disable it. The moment the doctrine hard-codes a threshold it cannot recompute, it rots.

The precedent already exists and is written down: `dataviz` states in its own header that *"Checks 1 and 6 are structural rules the skill enforces, not measurable from hexes alone."* Nobody has yet shipped the split as **two artifacts with a defined interface**. That interface is the deliverable.

## 5.2 Ownership table

| Concern | `parley-design` (doctrine, pure markdown, zero deps) | `parley-design-check` (tooling, may ship scripts) |
|---|---|---|
| The protocol (phases, gates, roles, artifacts) | **owns** | reads it to know what to lint |
| The rule registry (ids, prose, rationale, examples, tiers, sources, dates) | **owns** — it is a markdown file | reads it; ships **no rule text of its own** |
| Rule *detectors* (regexes, AST queries, DOM probes, browser scripts) | **never** | **owns** |
| Thresholds and constants | **owns as data in the registry** | reads them; MUST NOT hardcode a second copy |
| Which rules are `SLOP` vs `QUALITY` | **owns** | reads it to set exit codes |
| Severity anchors (Nielsen 0–4) | **owns** | applies them |
| Evidence tiers per rule | **owns** | reports which tiers it actually reached |
| The token contract (DTCG adoption, alias direction, required groups) | **owns the requirement** | **owns the validation** |
| Design-system *values* (this project's palette, scale, faces) | **never** — those live in the deck's ratified `design/` | reads `design/` + `design-check.config.json` and validates **against them**, not against a universal taste |
| Divergence / critique / verdict rituals | **owns** | can only check *structure* (artifact presence, recusal, graft count, axis coverage) |
| Taste, hierarchy, originality, brief-fit | **owns** | **must explicitly refuse** — §0 of the check spec says so out loud |
| Waivers | **owns the policy** (counter-signature required, reason is the artifact) | **owns the mechanism** (parse, apply, report) |
| Reports | **owns the required shape** | **owns the emission** |

## 5.3 THE CONTRACT (numbered, binding on both skills)

**C1 — One registry, one copy, no generation step.**
The rule registry is a single markdown file, `parley-design/references/RULES.md`, in which **each rule is an H3 heading + a fenced ```yaml metadata block + prose**. The fenced blocks are the machine source; the prose is the human source; they are the same file, so they cannot drift.
```yaml
id: core:color-accent-purple
class: slop
tier: css-parse            # cheapest engine that can decide it
also-checkable-at: [browser]
severity: 2                 # Nielsen 0-4
targets: [css-vars, tailwind]
enforced-by: check          # check | agent-judgement | both
yields-to: []               # stand-down discipline
added: 2026-07-28
confidence: confirmed
sources: [S1, S3, S4, S6, S8]
status: active
```
This is the direct fix for AG-UI's two worst defects (`events.proto` rotted to 16 vs 33; README stale at "~16"). If a second representation is ever introduced, it MUST be generated and MUST be guarded by a failing test — the `TestEmbeddedDefaultMatchesLiveDeck` pattern already proven in parley-deck-cli.

**C2 — The doctrine never restates a threshold; it declares one.**
Numbers live in the registry's metadata (or in `design-check.config.json` for per-project values). Prose cites the id. **Never write a count in prose** — the number of rules is generated or omitted.

**C3 — Every rule declares its evidence tier, and the doctrine tells the agent what to write when it cannot reach it.**
This is how the doctrine references checks it cannot run. The doctrine says:
> *A participant that cannot reach a rule's declared tier MUST write `unjudgeable: <tier>` for that rule in its artifact. This is COMPLIANT. Guessing is not. `"imagine the rendered output"` is not a verification method and MUST NOT appear in any artifact.*

Four verdicts, not two: `pass` · `violation` · `needs-review` · `unjudgeable`. *(Tri-state from IBM Equal Access; `unjudgeable` from impeccable's `fontSizeStepStatus`.)*

**C4 — Declared degradation.**
Any artifact produced without reaching every tier its enabled rules require MUST lead with a banner:
`⚠️ DEGRADED — tiers reached: text-regex, css-parse. Not reached: browser (no Chromium). 14 rules marked unjudgeable.`
A silent degraded review is a failed review. This matters more here than anywhere else because Parley's roster is heterogeneous by construction.

**C5 — Capability federation, not capability requirement.**
The browser tier is **one participant's optional attachment**, consumed by the others as evidence — never a protocol requirement. Measured cost is ~537 ms/page once the ~150 MB binary exists; the cost is the install, not the runtime. Non-determinism across five machines would poison consensus if it were mandatory.

**C6 — Target profiles: discovery only, no negotiation.**
`BRIEF.md` declares `target_profiles[]`. The checker runs only the rules those profiles declare checkable. **Absent = undeclared.** No target list in the doctrine needs updating when a new stack appears, and this is the actual vendor-neutrality mechanism (not a claim of neutrality).

**C7 — Error strings ARE the rule text.**
Every finding is `rule-id — violation — remedy`, always all three, and the message is **copied verbatim, never paraphrased**, so findings are diffable across runs and across agents.
`core:token-raw-literal — Button.tsx:41 declares '#5b3df5' outside the token layer. Replace with var(--color-accent) or add the value to primitives.tokens.json and re-run D7.`

**C8 — Exit-code contract.**
`0` = clean (or advisory-only) · `2` = findings · `1` = the tool itself failed · `3` = config error. Plus `--json`, `--threshold <n>`, `--fast` (T0/T1 only, no browser). Distinguishing "clean" from "tool broke" is what makes it CI-safe; a 0/1 contract cannot.
**One severity key, not two.** Impeccable ships `severity: 'advisory'` *and* `advisory: true` as different things and 11 rules are labelled advisory while still exiting 2. We derive the gate from exactly one field.

**C9 — Report shape is owned by the doctrine, emitted by the tool.**
```
[severity] rule-id — Tell name — file:line
  why it's a tell (one line)
  → fix (one line)

Summary — N violations · M needs-review · K unjudgeable · J recommendations
Tiers   — reached: … / not reached: …
Verdict — ships as slop | reads as AI-generated | close, fix the minors | clean
```
Trivially convertible to SARIF, to a PR annotation, and to `ReportFindings`.

**C10 — Waivers are counter-signed and the reason is the artifact.**
Syntax borrowed from impeccable (`design-check-disable[-next-line] <rule> -- <reason>`, comment-syntax-agnostic); **semantics inverted**: the reason is parsed, stored, and surfaced in `AUDIT.md`, and a `SLOP`-class waiver requires a second participant's acknowledgement recorded in `LEDGER.md`. A bare waiver with no reason is itself a violation. The narrowest-exception ladder is enforced: exact value → value-in-file → file → rule.
At least one rule class stays **deliberately design-system-blind** — a multi-agent design system *will* be widened by an implementer to legalise its own output.

**C11 — The check never runs without a contract.**
`design-check.config.json` (with its own published JSON Schema and `implements: PDS/1.0`) is required input. It names the token globs, the target profiles, the enabled rule ids and any overridden thresholds. Consequence: **the checker ships zero design opinions.** All scale checks are *derived* from the deck's own `*.tokens.json`. That is what makes it theme-relative rather than an aesthetic tyrant.

**C12 — Registry pinning.**
Every `AUDIT.md` carries `registry: core-rules/x.y.z` + `registry-digest`. A signature under a registry version that no longer matches is invalid, not merely stale.

**C13 — Versions are independent, compatibility is declared.**
`parley-design` versions the spec (`PDS/1.0`); `parley-design-check` versions the tool and declares `implements: PDS/1.x`. Neither number is derived from the other. *(AG-UI already skewed: `@ag-ui/core` 0.0.57 vs `ag-ui-protocol` 0.1.19.)*

**C14 — Fixtures gate new rules.**
A rule id cannot be merged without a golden good/bad fixture pair. Fixtures are plain files + one script, run offline in < 5 s. No Next.js dojo, no API keys, no 27 services.

**C15 — `parley-design` is fully usable with `parley-design-check` absent.**
Every registry entry is human-readable prose an agent can apply by reading. The check makes it cheap, reproducible and CI-able; it is never the only path. Conversely `parley-design-check` MUST run standalone against files on disk with no agent runtime — a vendor in Go or Swift must be able to invoke it.

## 5.4 Tiered implementation of the check (one rule corpus, three surfaces)

Copying impeccable's proven "one detector, four surfaces" insight, minus its 23 commands and 14 harness directories:

| Tier | What | Dependencies | Coverage |
|---|---|---|---|
| **T-0** | The registry itself, applied by hand by any agent | none | all rules, at the agent's own evidence tier |
| **T-1** | Zero-dependency script reading the same registry: text-regex + CSS-parse + JSON-schema | a runtime, no browser, no network | ~60 % of the value at ~0 % of the cost |
| **T-2** | Optional browser stage, auto-detected, clean degradation message | headless Chromium (~150 MB) | contrast, target size, overflow sweep, focus, reflow — the rules T-1 provably cannot do |
| **T-3** | Optional delegation to installed tooling, pinned by rule id | whatever the repo already has | `stylelint` core (149 rules), `scale-unlimited/declaration-strict-value`, `axe-core`, `@projectwallace/css-analyzer`, `@terrazzo/cli`'s 27 token lint rules |

**Not shipped, ever:** pixel visual regression (answers "did it change", not "is it good"; Percy needed AI to filter 40 % of its own false positives), APCA as a gate (pulled from WCAG 3 in July 2023, "yet to be determined"), the Lighthouse score (a weighted subset of axe with no partial credit), jsdom a11y (a false all-clear — no layout, so contrast silently cannot run), any Figma dependency (account + token + network), framework-coupled linters as the core engine.

---

# 6. Open questions for the humans

Each is a genuine fork. The trade-off is named; where I have a recommendation it is marked, but the choice is not mine.

**Q1 — Surface scope: web-only, or surface-agnostic with a web annex?**
*Web-only* inherits ~200 hard numbers immediately (every threshold in the corpus is CSS-shaped) and excludes TUIs, native, docs and slide decks — including Parley's own TUI, the most likely first customer. *Surface-agnostic* keeps the generic layer honest but risks it degenerating into vibes, and the annex will carry 80 % of the substance anyway. **Recommendation: surface-agnostic core + `WEB-ANNEX.md`, with target profiles as the seam (C6).** Not cheaply reversible — decide before the first rule is written.

**Q2 — Add-on shape: teaching skill, or protocol amendment?**
`parley-worktrees` and `parley-tracker` are opt-in add-ons that **never change canonical artifact ownership**. Does `parley-design` follow that (teaches a discipline, owns no phase) or does it amend COOPERATION.md with a binding design phase and its own gate? The former is safe and easily ignored; the latter has teeth and touches the core protocol the other add-ons deliberately do not. *This determines whether D0–D9 are real phases the driver can enforce or a discipline agents are asked to follow.*

**Q3 — Registry source of truth: literate markdown (C1), or `rules.json` in the check with generated prose?**
*Literate markdown* = one copy, no build step, works with the check absent, and the doctrine stays "pure markdown" as specified. Cost: the check must parse fenced YAML out of markdown. *`rules.json` + generated prose* = a cleaner machine artifact, at the price of a generation step, a drift guard, and a doctrine skill that cannot be edited without running a tool. **Recommendation: literate markdown.** But it hard-codes a parser contract, so say yes or no now.

**Q4 — Convergence semantics: PASS or ALARM?**
When N independent models propose the same direction, is that (a) strong evidence the direction is right, or (b) evidence of a shared training attractor? Both are defensible and they demand **opposite** protocol behaviour: (a) fast-tracks to consensus; (b) triggers G1 failure and a forced-distinct re-run. A hybrid — *convergence passes only if the converged direction survives the SLOP registry and the "could you guess this aesthetic from the category alone, or from category-plus-avoidance?" test* — is possible but must be specified precisely or it is a coin flip. **The spec currently assumes (b).**

**Q5 — Who is the Decider when no human is present?**
The evidence says the Decider must be a *role*, not a vote, and preferably human. But Parley runs headless and auto-drives by default. Options: (i) the run **stalls** at D5 awaiting a human supervote (safe, breaks auto-drive); (ii) a **designated agent** holds the supervote and its override reason is logged for later human review (keeps auto-drive, reintroduces exactly the LLM-judge biases §4.2 exists to neutralise); (iii) **highest trimmed-mean score wins by default**, with the human able to override post-hoc (mechanical, but it is a vote, and voting is the documented failure mode). No option is clean.

**Q6 — `SLOP` quorum threshold.**
A `SLOP` finding needs quorum to block. Is quorum (a) a simple majority of non-authors, (b) ≥2 independent agents citing the same rule id, or (c) ≥4 concurrent tells on one artifact (the empirically anchored Krebs threshold)? (c) is the only number with evidence behind it, but it measures *artifacts*, not *agreement*, so it may need to be **both**: ≥2 agents **and** ≥4 tells.

**Q7 — Rendered evidence: forbidden, optional-with-declaration, or required?**
*Forbidden* is honest and weak (no layout rule ever fires). *Required* makes hermes/agy/kimi structurally non-compliant. *Optional with declared capability* is obviously right and requires the full evidence-tier + `unjudgeable` + degradation-banner machinery (C3, C4, C5). There is no cheap version of this. **The spec assumes optional-with-declaration.** Confirm you are willing to pay for the machinery.

**Q8 — When is the design system written: before implementation, or after?**
impeccable is emphatic: *after*, from shipped code — *"a rulebook written before the build gets defended against reality instead of describing it."* Parley is emphatic the other way: `FINAL.md` is ratified *before* Phase 5 and binds implementers. These are incompatible if it is one artifact. Either accept a pre-build binding contract and lose the "describes reality" property, or ship **two** artifacts with different names and different authority (`VERDICT.md` binds before; `design/DESIGN-SYSTEM.md` is authored after from the merged diff) — and pay for two. **The artifact set in §2 currently assumes one, written at D7.**

**Q9 — Where does taste live: shipped canon, ratified per deck, or nowhere?**
*Shipped canon* (hallmark's 20 themes, impeccable's house style) makes the skill immediately useful and makes every project using it look the same — the disease. *Ratified per deck* is correct and means the skill does nothing useful on day one of a greenfield project. *Anti-slop invariants only* is the purist position and leaves agents falling back on training priors in every free axis — which is what both source projects were built to prevent. **Recommendation: ratified-per-deck + invariants only.** It is the slow answer.

**Q10 — Per-model prior corrections: ship them, or measure them?**
impeccable ships model-attributed doctrine (`<claude>` "your cream/serif/lamplight prior — treat that first palette as already spent"; five detector rules are literally named after the model that emits them). Parley is uniquely positioned to **measure** priors instead of guessing: it runs the same brief through 4–5 models on every idea. Options: (i) ship a hand-written prior block per participant now (fast, opinionated, ages badly, and one COOPERATION.md is read by all agents at once so it must be injected at runtime); (ii) accumulate `LEDGER.md` data first and derive priors in a later version (slower, genuinely novel, and the only version that will still be true in a year).

**Q11 — Anti-self-repetition: what is the memory, and who owns it?**
`core:self-repetition` is our clearest differentiator and it needs a project-scoped history. hallmark's `.hallmark/log.json` is a gitignored, single-writer, trimmed-to-20 file with no author field and no merge semantics — five agents will clobber it. Options: (i) derive history from the committed `LEDGER.md` files under `parley-deck/ideas/*/` (no new state, but only sees Parley-run designs); (ii) a committed, append-only `parley-deck/design/HISTORY.jsonl` keyed by `(idea, agent, round)` (explicit, but a new canonical artifact to maintain). **Recommendation: (i).**

**Q12 — Does `parley-design-check` ship in the skill, in parley-deck-cli, or as its own package?**
In the skill = zero install, but a skill shipping executable code breaks the "companion add-on" precedent. In `parley-deck-cli` = one Go binary already installed everywhere, drift-guardable, and `RunChecks`-shaped — but couples the design registry to the CLI release cadence and forces every rule to be implementable in Go. Its own package = cleanest boundary, worst install story for a vendor-neutral roster. **This decides whether the browser tier is ever realistic.**

---

# 7. Anti-goals

What these two skills must refuse to become.

1. **A house aesthetic.** Not 20 named themes, not a signature palette, not warm-paper-editorial-with-hairline-rules. Ship anti-slop **invariants** separately from **preferences**, and let the deck's own ratified `design/` own the latter, overridable on record. Both source projects ship a taste and call it a standard.
2. **A prescriber of values.** Prohibitions of defaults, never prescriptions of alternatives. Fraunces was the prescribed 2025 fix and is a 2026 tell. Prescribe *process*; the moment we ship a palette we relocate the sameness instead of removing it.
3. **A CSS linter with a protocol wrapper.** No 5,500-line rule monolith, no hand-rolled CSS cascade, no committed 8,000-line browser bundle, no mandatory Puppeteer, no framework port-sniffing.
4. **A second state machine.** No 23-verb command surface. Parley already has phases, dispositions, consults, tracks and a driver. `audit` is the review phase; `redesign` is an idea whose input is existing code; `study` is input preparation.
5. **A blocking-question machine.** Five headless agents × three mandatory questions = fifteen blocking prompts and zero answers. Ask the human **once**, at D0, and convert every later ask into a labelled assumption or a Parley consult.
6. **A single-writer mutable state tree.** No gitignored rotation log, no `~/.parley-design/` throttle cache, no mtime staleness. Canonical state is `parley-deck/ideas/<slug>/`, keyed by `(idea, agent, round)` and pinned by content hash.
7. **Self-scoring theatre.** No `58/58 ✓` written by the author, no `P5 H5 E5 S4 R5 V5` self-stamp, and above all **no gate whose verification method is "imagine the rendered output."** A gate that cannot be executed manufactures passes. Keep the tally format; move authorship to a peer or a tool.
8. **A numeric design score / a headline grade.** Calibrated rubrics with honest denominators are fine; a 0–100 invites agents to optimise the number. Human–human aesthetic agreement is 38.34 % — a grade built on that is a coin flip wearing a lab coat.
9. **An aesthetic hard gate.** Taste findings are advisory, always. `ABSTAIN` is a legitimate verdict, and it is the published route to human-level judging performance.
10. **A forced-diversity generator.** Do not import "differ from the last run on at least one axis" as a build-time chore. With N agents on one idea, mandated difference produces incoherence; import the axes as a **similarity metric** and decide per context whether similarity is good.
11. **A committee that averages taste.** Restated because it is the most likely failure. The merge operator is *choose*. If the protocol cannot select, it escalates to the human — it never compromises.
12. **A silent degrader.** No participant signs a verdict without declaring which tiers it could not reach, which rules it marked `unjudgeable`, and what it did not read. impeccable's own `detectHtml` silently falls back to regex and loses ~40 of 60 rules with no diagnostic; that exact bug, inside a consensus protocol, is a signed lie.
13. **Hand-maintained facts.** hallmark says 57 gates in its README and 58 in the file; "twenty-four themes" over 20 definitions; a documented verb that does not exist. Anything stated twice must be generated or drift-guarded.
14. **File sprawl.** 24 references or 34 playbooks survives only because exactly one file is read per request by one agent. Under N agents it is 5× the cost and a guaranteed drift surface. **Four files, hard ceiling.**
15. **Vendor coupling.** No hosted CLI, no login wall, no telemetry, no per-harness directory sprawl, no Figma dependency, no brand-domain exemption lists ("Geist is slop unless you are Vercel"). Local-file-first, reproducible from the run id, offline.
16. **Moralising.** "Slop" is descriptive, not a verdict on the author. Willison's own caveat: *"Not all AI-generated content is slop."* Rule names are descriptive (`core:color-gradient-hero`); findings read as "this is a recognisable default", never as an accusation.
17. **Citing numbers we cannot source.** The "familiar patterns are 3–5× faster" Jakob's-Law figure traces to no NN/g primary; the 21.2 % AI-detection penalty is real but from *advertising with explicit disclosure*, not UI. Quoting either would make the spec commit `core:honesty-invented-metric` — the exact rule it enforces. The camel proverb is a 1954 *Reader's Digest* quip with a contested attribution: quote the provenance or don't quote it.
18. **An accessibility deferral with no backstop.** Keeping the a11y checklist out of the generative round is correct (it produces safe, underdesigned output). Bind it to a hard gate instead: **no RATIFIED signature without the `QUALITY` a11y rules having been evaluated by someone, at a declared tier.**
