# Peer AI design skills / plugins — state of the art survey

Survey date: 2026-07-28. Purpose: know what already exists before designing **parley-design**
(doctrine, pure markdown) and **parley-design-check** (enforcement, may ship scripts).

Legend: **[F]** = fact I verified by reading the artifact or its published page.
**[I]** = my inference / judgement.

---

## 0. Method + what was actually read

**[F]** Local filesystem enumeration:

```
find /Users/tomasfecko/.claude -maxdepth 5 -name 'SKILL.md'
→ /Users/tomasfecko/.claude/skills/parley-tracker/SKILL.md
  /Users/tomasfecko/.claude/skills/parley-deck/SKILL.md
  /Users/tomasfecko/.claude/skills/hallmark/SKILL.md
  /Users/tomasfecko/.claude/skills/parley-worktrees/SKILL.md
  /Users/tomasfecko/.claude/skills/graphify/SKILL.md
```

**[F]** `artifact-design` and `dataviz` are **not** on disk as `SKILL.md`. Their SKILL.md bodies are
embedded in the Claude Code binary (`/Users/tomasfecko/.local/share/claude/versions/2.1.220`,
Mach-O 64-bit arm64, 256,908,272 bytes). I read them by invoking the `Skill` tool. Only dataviz's
*payload* (`references/`, `scripts/`) is lazily extracted to
`/private/tmp/claude-501/bundled-skills/2.1.220/4fdadcd2f8243ffac061f198d7157521/dataviz/`.
**[I]** Consequence for us: a skill can be *pure prompt* (no disk footprint) or *prompt + payload*.
Anthropic ships both shapes in the same product.

**[F]** Plugin caches read:
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/vercel/0.45.1/skills/shadcn/SKILL.md` (596 lines)
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/frontend-design/unknown/skills/frontend-design/SKILL.md` (55 lines)
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/figma/2.2.81/skills/figma-generate-library/SKILL.md` (370 lines)
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/figma/2.2.81/figma-power/steering/create-design-system-rules.md` (510 lines)
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/figma/2.2.81/skills/figma-implement-motion/references/motion-lint-rules.md`

---

## 1. The comparison table

| # | Name | Source | Size / shape | Doctrine it encodes | Enforcement | One idea worth stealing |
|---|---|---|---|---|---|---|
| 1 | **hallmark** v1.1.0 | `/Users/tomasfecko/.claude/skills/hallmark/` · [github.com/nutlope/hallmark](https://github.com/nutlope/hallmark) (MIT, 19.1k★) | 106 files, 9,591 md lines, 932 KB. `SKILL.md` (558 lines) + `references/` (30 topic files + 50 component files + 21 macrostructure files + 4 genre files + 4 theme specs). **Zero executable code.** | "Structural variety, not just visual variety." Pick a *macrostructure* first (21 named page-shapes), then a theme from a catalog of 20, then nav/footer archetypes. Diversify against project memory. | **Self-enforcement in prose:** 58-gate slop test + 6-axis pre-emit self-critique, run at Step 7. No script. | **Machine-readable provenance stamp + rotation memory**: `/* Hallmark · macrostructure: <name> · tone: <tone> · anchor hue: <hue> */` written into the CSS, plus `.hallmark/log.json` (last 20 entries) read on the *next* run to force a different pick. |
| 2 | **impeccable** v4.x | [github.com/pbakaus/impeccable](https://github.com/pbakaus/impeccable) (Apache-2.0, ~52k★ per repo page) | 1 skill + **23 commands** + `cli/`, `scripts/`, `extension/`, `plugin/`, `tests/`; per-harness dirs `.claude/ .codex/ .cursor/ .gemini/ .grok/ …` | Design *language* + full lifecycle: `init → shape → craft → audit → polish`, with modifier passes (`bolder`, `quieter`, `distill`, `harden`, `typeset`, `layout`…). Two shared artifacts: `PRODUCT.md` (audience/brand) and `DESIGN.md`/`design.json` (voice, **anti-references**, colors, type, components). | **Real, deterministic:** `npx impeccable detect` = **60 deterministic detector rules**, no LLM, no API key. Exit codes **0** = no findings, **2** = findings, **1** = command failed. Flags `--json`, `--scope type|layout`, `--no-design-system`, `--no-config`. Config `.impeccable/config.json` → `detector.ignoreRules`, `detector.ignoreFiles`, `detector.ignoreValues`, `detector.designSystem.enabled`. Ships **edit-time hooks** (`hook.mjs`, `hook-before-edit.mjs`) wired via provider-native manifests (`.claude/settings.local.json`, `.cursor/hooks.json`, `.codex/hooks.json`, `.github/hooks/impeccable.json`). | **The doctrine/enforcement split already exists here and works**: the *same* detector powers the CLI, the pre-edit hook, `/impeccable audit`, and a browser extension. One rule engine, four surfaces. Also: **`anti-references`** as a first-class field in DESIGN.md (say what you must NOT look like). |
| 3 | **dataviz** (Anthropic, bundled) | Claude Code 2.1.220 binary; payload at `/private/tmp/claude-501/bundled-skills/2.1.220/…/dataviz/` | 9 files, 1,314 lines. `references/` ×7 + `scripts/validate_palette.js` (316 lines) + `validate_palette.py` (305 lines) | "A chart is read by people and executed by you… right by construction rather than by taste." **7-step ordered procedure; color comes LAST.** Design-system-**agnostic method** + swappable **parameter table** (ramps, categorical theme, sequential hue, diverging pair, status palette, texture, surfaces, filter controls). | **The strongest enforcement of anything surveyed.** A runnable validator with published numeric thresholds: `BAND = {light:[0.43,0.77], dark:[0.48,0.67]}` (OKLCH L), `CHROMA_FLOOR = 0.10`, `CVD_TARGET = 8.0` / `CVD_FLOOR = 6.0` (OKLab ΔE×100 under Machado-Oliveira-Fernandes 2009 severity-1.0 protan/deutan), `NORMAL_FLOOR = 15.0` (hard FAIL), `CONTRAST_MIN = 3.0` (WCAG vs surface), `ORDINAL_MIN_DL = 0.06`. Exit 1 on FAIL; WARN bands exit 0. Runs as CLI **or** as a `<script type="module">` inside the artifact page reading `data-palette` off `<body>`. | **Separate the computable from the judgeable, and say so explicitly.** The header states: *"Checks 1 (fixed hue order) and 6 (values are from the documented palette) are structural rules the skill enforces, not measurable from hexes alone."* That is exactly the parley-design / parley-design-check boundary, written down. Also: **"The single most important habit: the color part is computable, so compute it. Never eyeball…"** |
| 4 | **artifact-design** (Anthropic, bundled) | Claude Code 2.1.220 binary (prompt-only, no disk payload) | ~1 page of prose, no references, no scripts | **Treatment calibration before craft**: "Calibrate treatment, not whether to design." Utilitarian vs editorial branch. Hard precedence rule: *"the user's own words, then the project's existing system, then your choices."* Named default-cluster to avoid. Dual-theme token discipline (`prefers-color-scheme` + `:root[data-theme=…]` must win both directions). | None (pure prose). | **The precedence ladder** and **the "when it's a UI, not a document" pivot** ("scanned and operated, not read top-to-bottom… encode state in form as well as number… Semantic color (good/warning/critical) is separate from the accent hue and doesn't count as your accent"). |
| 5 | **frontend-design** (Anthropic official plugin) | `claude-plugins-official/frontend-design/…/SKILL.md` | **55 lines. One file. No references, no scripts.** | The ur-text the others descend from. "Design lead at a small studio… the client has already rejected proposals that felt templated." Process: **brainstorm, explore, plan, critique, build, critique again**. Plan = 4 fields: **Color / Type / Layout / Signature**. | None. Self-critique only ("Consider Chanel's advice: before leaving the house… remove one accessory"). | **The 4-field design plan with `Signature` as a required field** — "the single unique element this page will be remembered by". And the **calibration paragraph naming the three current AI clusters** ((1) cream #F4F1EA + serif + terracotta, (2) near-black + acid-green/vermilion, (3) broadsheet hairlines + dense columns) — a *dated* list, honestly labelled "right now". |
| 6 | **vercel:shadcn** | `claude-plugins-official/vercel/0.45.1/skills/shadcn/SKILL.md` | 596 lines, single file, no scripts | Component-library-native design doctrine: "In the Vercel stack it is the default interface language. Do not stop at 'the component works.'" Ships a **"Reach for this first" routing table** (9 use-cases → component sets → *why*), composition recipes, and a 9-item **anti-patterns list**. | **Frontmatter-declared static check** — the most portable enforcement mechanism seen: `validate: [{ pattern: '"base"\s*:\s*"base-ui"', message: '…', severity: warn }]`. Plus routing metadata `pathPatterns`, `bashPatterns`, `retrieval.{aliases,intents,entities}`, `metadata.priority: 6`. | **`validate:` blocks in SKILL.md frontmatter.** A skill can carry regex rules + message + severity *declaratively*, with no script and no runtime — the harness runs them. This is the cheapest possible parley-design-check v0. |
| 7 | **figma-generate-library** + `create-design-system-rules` | `claude-plugins-official/figma/2.2.81/` | 370 lines + 510 lines | The most **protocol-shaped** artifact surveyed. Mandatory phases **0 DISCOVERY → 1 FOUNDATIONS → 2 FILE STRUCTURE → 3 COMPONENTS → 4 INTEGRATION+QA**, with stable task IDs `P{phase}.{step}` (`P0.a`, `P3.d`), a posted `Phase N Checklist` before mutation, a `Phase N Summary` after, and explicit **exit criteria** per phase. Rule 1: "**Variables BEFORE components** — components bind to variables. No token = no component." | Procedural + auditable: Phase 4 = `4b. Accessibility audit`, `4c. Naming audit (no duplicates, no unnamed nodes, consistent casing)`, `4d. Unresolved bindings audit (no hardcoded fills/strokes remaining)`. Sibling `motion-lint-rules.md` defines severities **Error** (cannot generate) / **Warning** (generated with a known gap) and demands *"Copy the message verbatim — do not paraphrase, summarize, or rephrase it."* | **Phase gating with printed exit criteria and stable task IDs**, plus the *"print a gap analysis to chat: what exists in code but not Figma, what exists in Figma but not code, and every conflict with its resolution"* (0f). And verbatim-message discipline for lint findings. |
| 8 | **design-review** subagent | [OneRedOak/claude-code-workflows/design-review](https://github.com/OneRedOak/claude-code-workflows/tree/main/design-review) | `design-review-agent.md` (subagent, `model: sonnet`, ~40 tools incl. Playwright MCP) + `design-principles-example.md` ("S-Tier SaaS Dashboard Design Checklist", 7 top-level sections) | **"Live Environment First"** — assess the *rendered, interactive* experience before static analysis. 8 phases: 0 Preparation → 1 Interaction/User Flow → 2 Responsiveness (1440 / 768 / 375 px) → 3 Visual Polish → 4 Accessibility WCAG 2.1 AA → 5 Robustness → 6 Code Health → 7 Content & Console. | Browser-driven, not static: Playwright MCP drives the real page. Findings triaged **[Blocker] / [High-Priority] / [Medium-Priority] / [Nitpick]** (nits prefixed "Nit:"), reported under fixed headings. | **A severity taxonomy with a fixed report skeleton**, plus **"Live Environment First"** — the only surveyed artifact that insists design review happens against pixels, not source. Its `design-principles-example.md` is the *project-supplied* rubric the agent reads — i.e. the rubric is data, not baked into the agent. |
| 9 | **superdesign-skill** | [github.com/superdesigndev/superdesign-skill](https://github.com/superdesigndev/superdesign-skill) (MIT, 368★) | skill + `@superdesign/cli`, `DESIGN.md`, `AGENTS.md`, `CLAUDE.md`, `.claude/commands/`, `.codex-plugin/` | Three scenarios: "Help me design X" / "Set design system" / "Help me improve design of X". Flow: extract-or-create design system → `.superdesign/design-system.md` → replica HTML template of current UI → baseline draft → **branch variations**. | Requires a hosted CLI + login (`superdesign login`) — an availability liability. | **`iterate-design-draft --mode branch` — "each prompt = one variation"**, i.e. parallel divergent directions is already a productised primitive. Also the **"replica HTML template of the current UI" step** before designing: reproduce what exists so the new direction is compared against a real baseline, not a memory of one. |
| 10 | **mattbx/shadcn-skills** | [github.com/mattbx/shadcn-skills](https://github.com/mattbx/shadcn-skills) | 2 skills: `shadcn-component-discovery` (+`references/registries.md`) and `shadcn-component-review` (+`references/theme-styles.md`, `review-checklist.md`, `animation-patterns.md`). No scripts. | Reuse-before-build ("discover 1,500+ existing components before building custom") + a *review* skill scored against 5 axes: Structure (`data-slot`, composition), Spacing (`gap-*` vs `space-y-*`), Design tokens (semantic colors only), Composability (`className` via `cn()`, CVA), Responsive/a11y (mobile-first, touch targets, `min-w-0`, focus states). | Checklist-in-markdown only. | **The split into a *discovery* skill and a *review* skill** — the same author shipped exactly our two-skill shape. And **per-theme expectations** (Vega, Nova, Maia, Lyra, Mira → each has its own spacing/radius/density expectations), i.e. the check is *theme-relative*, not absolute. |
| 11 | **bergside/awesome-design-skills** | [github](https://github.com/bergside/awesome-design-skills) | **67** design systems, each a folder `skills/<name>/{SKILL.md, DESIGN.md}` + `skills/index.json` slug→path map | Style presets as pluggable data: Brutalism, Glassmorphism, Editorial, Neumorphism, Claymorphism, Bento, Dithered, Retro, Material, Minimal, … | None. | **The two-file convention: `SKILL.md` = instructions for agents, `DESIGN.md` = human-readable rationale/references/maintenance**, plus an `index.json` registry so a *slug* can be resolved and pulled. |

---

## 2. What the field agrees on (the consensus anti-slop canon)

**[F]** These specific tells appear in **three or more** independent artifacts, which makes them the
de-facto canon rather than any one author's taste:

| Tell | hallmark | impeccable | frontend-design | artifact-design | vercel:shadcn |
|---|---|---|---|---|---|
| Purple→blue / gradient hero **and** gradient-clip text | gate 2 + "The purple-gradient hero" | "purple gradients" detector | cluster (2)-adjacent | "purple-to-blue gradient hero on white" | "Large gradient backgrounds… on every surface" |
| Inter / Roboto / system default as display face | gate 1, "Inter-everywhere" | "Don't use overused fonts (Arial, Inter, system defaults)"; rule id `overused-font` | — | "Inter or Space Grotesk as the 'safe' face" | — |
| 3-equal-column icon-above-heading card grid | gate 3, "The 3-column feature grid" | — | — | "everything centered; `rounded-lg` everywhere" | — |
| Card-in-card nesting | gate 4 | "Don't wrap everything in cards or nest cards inside cards" | — | — | "Nested cards inside cards inside cards" |
| Thick coloured side-stripe on cards | gate 5 | "side-tab borders" detector | — | "accent bar/rail on rounded cards" | — |
| Pure `#000` / `#fff` | gate 7 | "Don't use pure black/gray (always tint)" | — | "A pure mid-grey reads as unconsidered" | — |
| Bounce / elastic easing on UI state | gate 12 | "Don't use bounce/elastic easing (feels dated)" | — | — | — |
| Full-viewport centred hero | gate 6, "Full-viewport centred hero" | — | hero-is-a-thesis | "everything centered" | — |
| Emoji as section/feature icons | gate 30(b) | — | — | "emoji as section markers" | — |
| Invented metrics / fake social proof | gate 46 + "Honest copy" discipline | — | — | — | — |

**[I]** Two of these are *ours to sharpen*: **invented metrics** (only hallmark gates it) and
**re-drawn UI chrome** (hallmark gate 47 — fake browser bars, phone frames, terminal frames — nobody
else names it). Both are trivially machine-detectable and are strong differentiators for
parley-design-check.

---

## 3. The three enforcement architectures that exist today

**[F]**

1. **Prose self-gate** (hallmark, frontend-design, artifact-design, mattbx). A numbered list the
   model is instructed to run against its own output. hallmark is the most elaborate: 58 gates +
   6-axis pre-emit critique (Philosophy / Hierarchy / Execution / Specificity / Restraint / Variety,
   scored 1–5, `< 3` triggers a revision pass, scores stamped as
   `/* Hallmark · pre-emit critique: P5 H4 E5 S4 R5 V5 */`).
   *Cost:* zero deps. *Weakness:* the checker is the author. **[I]** This is precisely the failure
   mode a multi-agent protocol removes.
2. **Declarative frontmatter rules** (vercel:shadcn `validate:` block). Regex + message + severity,
   run by the harness. *Cost:* zero deps, but harness-specific. **[I]** Not vendor-neutral —
   `validate:` is a Vercel-plugin/Claude-Code convention, not a standard.
3. **Real deterministic engine** (impeccable `detect`, dataviz `validate_palette.js`,
   plus the wider ecosystem: `ds-lint` (Rust), `stylelint-declaration-strict-value`,
   `stylelint-plugin-carbon-tokens`, Kong `design-tokens/stylelint-plugin`,
   Mozilla `no-base-design-tokens`, `stylelint-plugin-rhythmguard`).
   *Cost:* a runtime + install. *Strength:* CI-able, hook-able, exit-code-able, no LLM.

**[I]** For **parley-design-check** the right answer is a **tiered** architecture: tier 0 = a
declarative rule file (vendor-neutral, e.g. `design-rules.yaml`) that any agent can read and apply by
hand; tier 1 = a zero-dependency Node script that reads that same rule file and emits JSON +
exit codes; tier 2 = optional delegation to `impeccable detect` / stylelint when present. One rule
corpus, three execution surfaces — copying impeccable's "one detector, four surfaces" insight.

---

## 4. Protocol-shaped precedents (the AG-UI-style ask)

**[F]** Only two surveyed artifacts are genuinely *protocol*-shaped rather than advice-shaped:

- **figma-generate-library** — versioned phases with **exit criteria**, **stable task IDs**
  (`P{phase}.{step}`), a mandatory pre-phase checklist post, a mandatory post-phase summary, an
  explicit *"No setup exception: creating a new Figma file… all count as creation/mutation"*, and a
  hard non-negotiable ordering rule ("Variables BEFORE components. No token = no component.").
- **dataviz** — a 7-step ordered **procedure** with a declared invariant/parameter split
  (*"The method is invariant; only these parameters change per system"* + an 8-row parameter table),
  a runnable conformance checker, and named thresholds with units.

**[F]** hallmark is *almost* protocol-shaped — it has a versioned `version: 1.1.0` frontmatter, a
canonical machine-readable **stamp** format, a canonical **state file** (`.hallmark/log.json`, JSON
array newest-first, schema `{date, macrostructure, theme, enrichment, brief}`, trimmed to 20), and
numbered gates — but the artifacts are CSS comments, not typed documents, and nothing validates the
stamp.

**[I]** What none of them have: a **versioned spec document with a conformance section**, typed
artifact names, and a statement of what "conforming" means. That is the single largest opening.

---

## 5. What is MISSING across ALL of them — the multi-agent opening

**[I]** Every artifact surveyed is **single-author**. That produces five structural gaps a
Parley Deck protocol can uniquely close:

1. **No real divergence.** hallmark *rotates* one author through a catalog (macrostructure
   diversification, theme axes, nav/footer rotation) and superdesign *branches* variations from one
   model. Nobody has **N independent agents each committing to a different direction and then
   attacking each other's**. Rotation is anti-repetition; divergence is anti-consensus-collapse.
   Rotation cannot find the direction the author would never have picked.
2. **The critic is the author.** hallmark's 58 gates, impeccable's `critique`, frontend-design's
   "critique again", dataviz's checks — all self-scored (except OneRedOak's subagent, which is still
   one agent, and impeccable's detector, which only sees syntax). **No adversarial cross-review of a
   *visual direction* exists in the field.** A model grading its own taste is the exact failure the
   whole category was created to fix, reproduced one level up.
3. **No "one wins whole" rule.** Every artifact either picks one thing outright or blends. Nobody
   writes down the **anti-averaging** rule ("one direction wins whole; never an average; graft 2–3
   concrete details from the losers"). **[I]** This is the field's biggest unnamed failure mode:
   merging two coherent directions produces an incoherent third. Naming and enforcing it is genuinely
   novel.
4. **No signed, auditable decision record.** hallmark stamps a CSS comment; impeccable writes
   `DESIGN.md`; superdesign writes `.superdesign/design-system.md`. None records **who proposed
   what, who objected, what the objection was, and why the winner won.** A design system whose
   rationale is lost is a design system that gets silently violated six weeks later.
5. **Doctrine and enforcement are not separated on purpose anywhere except dataviz** — and dataviz
   only does it *within* one skill (its header line "Checks 1 and 6 are structural rules the skill
   enforces, not measurable from hexes alone"). Nobody ships the split as two artifacts with a
   defined interface between them. **[I]** Our two-skill split is the right shape; the interface
   (a typed, versioned rule corpus emitted by parley-design and consumed by parley-design-check) is
   the thing to specify.

**[I]** Two smaller gaps worth naming:
- **No genericity/portability discipline.** Almost everything is React+Tailwind+shadcn-coupled
  (vercel:shadcn, mattbx, superdesign) or Figma-coupled (figma-*). hallmark and frontend-design are
  stack-neutral but web-only. Nothing addresses TUI, native, print, or docs.
- **No "conditional rigor."** Every artifact runs its whole apparatus every time. hallmark bolts on
  a "Component-scope flow" that *skips* macrostructure/nav/footer/enrichment — the closest thing to
  a track system — but there is no declared fast/standard/deliberation track. Parley Deck already
  has `track: fast|standard|deliberation` (§4.0); a design skill that inherits it would be unique.

---

## Transferable to parley-design / parley-design-check

Ranked, most valuable first. Each is concrete and cites its source.

### → parley-design (doctrine, pure markdown)

1. **Steal the invariant/parameter split verbatim as the skill's spine** (dataviz SKILL.md).
   State up front: "The method is invariant; only these parameters change per system," then ship an
   8-row parameter table (ours: neutrals, accent, semantic set, type roles, scale ratio, spacing
   base, radius set, motion easings+durations). This makes the skill vendor-neutral *by
   construction* rather than by claim, and it's exactly the protocol posture the owner wants.
2. **Adopt phase gating with printed exit criteria and stable task IDs** (figma-generate-library:
   `P{phase}.{step}`, `Phase N Checklist` before mutation, `Phase N Summary` after, per-phase exit
   criteria, "do not move to the next phase until the current phase's acceptance checks are
   complete"). Map onto DIVERGE → CRITIQUE → DECIDE → GRAFT → APPLY. Give each phase a required
   output artifact and a machine-checkable exit condition.
3. **Require a `Signature` field in every direction proposal** (frontend-design: plan = Color / Type
   / Layout / **Signature** = "the single unique element this page will be remembered by"). In a
   diverge protocol this is the field that *forces* the directions apart — two proposals with the
   same Signature are the same proposal, and that's mechanically checkable.
4. **Require `anti-references` as a first-class field** (impeccable's `DESIGN.md` captures "voice,
   **anti-references**, colors, typography, and components"). Saying what the thing must NOT look
   like is cheaper and more discriminating than saying what it should.
5. **Write the anti-averaging rule down as a numbered protocol rule.** Nobody has it (§5.3). Suggested
   binding form: *"Exactly one direction wins whole. A synthesis of two directions is a protocol
   violation. 2–3 named details may be grafted from losing directions; each graft must name its
   source direction and the property it carries."* The named-graft requirement makes it auditable.
6. **Emit a machine-readable provenance stamp + a rotation/decision ledger** (hallmark's CSS stamp +
   `.hallmark/log.json`, newest-first, trimmed to 20). Ours should be richer: which agent proposed
   the winner, which agents objected, the surviving objections, and the grafts — i.e. the signed
   consensus Parley Deck already produces, but in a form the *next* design run reads.
7. **Adopt the treatment-calibration opener** (artifact-design: "Calibrate treatment, not whether to
   design"; utilitarian vs editorial). Wire it to Parley Deck's existing `track: fast|standard|
   deliberation` so a memo doesn't trigger a 4-agent divergence round.
8. **Adopt the precedence ladder verbatim** (artifact-design): "the user's own words, then the
   project's existing system, then your choices." Two lines, kills a whole class of arguments.
9. **Require a "replica baseline" step before redesign** (superdesign: reproduce the current UI as an
   HTML template first). Directions get compared against the real current state, not a recollection.
10. **Ship the two-file convention** (bergside): `SKILL.md` = agent instructions, a separate
    human-readable rationale doc. Parley Deck already has this instinct (FINAL.md vs IMPLEMENTATION.md).
11. **Component scope ≠ page scope** (hallmark's "When the brief is a component, not a page" —
    2 signals route component; component scope *skips* macrostructure/nav/footer/enrichment/project
    memory but keeps an **8-state** requirement: default · hover · focus-visible · active · disabled ·
    loading · error · success, plus a labelled 8-state demo wrapper). The 8-state list is the single
    most reusable concrete checklist in the whole survey.

### → parley-design-check (enforcement, may ship scripts)

1. **Copy the tiered "one rule corpus, many surfaces" architecture** (impeccable: same detector powers
   CLI + pre-edit hook + `/audit` + browser extension). Tier 0 declarative rule file → Tier 1
   zero-dep Node checker → Tier 2 delegate to existing linters.
2. **Copy the exit-code contract exactly** (impeccable detect): **0** no findings, **2** findings,
   **1** command failed. Plus `--json`. Distinguishing "clean" from "tool broke" is what makes it
   CI-safe; a 0/1 contract cannot.
3. **Publish numeric thresholds with units and a citation, not adjectives** (dataviz:
   `CVD_TARGET = 8.0` OKLab ΔE×100 under *Machado, Oliveira & Fernandes (2009) severity 1.0*, with
   the note that swapping in Viénot-1999 "would require recalibrating these"). Every rule we ship
   must name its metric, its threshold, and its measurement model.
4. **Adopt the WARN/FAIL two-band model with obligations** (dataviz): a WARN is *"legal ONLY with
   mandatory secondary encoding"*; a contrast WARN *"obligates visible labels or a table view — it is
   not dismissable."* A warning that carries a required remedy is far stronger than a warning that
   carries a shrug.
5. **Adopt a severity taxonomy + fixed report skeleton** (OneRedOak: `[Blocker]` / `[High-Priority]` /
   `[Medium-Priority]` / `[Nitpick]`, nits prefixed "Nit:", fixed headings). Parley Deck's review
   phase already has BLOCK/dispositions — align the vocabularies deliberately rather than inventing a
   third.
6. **Adopt verbatim-message discipline** (figma `motion-lint-rules.md`: *"Copy the message verbatim —
   do not paraphrase, summarize, or rephrase it."* with Error = cannot proceed, Warning = proceeds
   with a known gap). Stable finding text is what makes findings diffable across runs and across
   agents — essential when 4 agents report on the same artifact.
7. **Adopt the config surface** (impeccable `.impeccable/config.json`): `ignoreRules`, `ignoreFiles`,
   `ignoreValues`, `designSystem.enabled`, plus in-file inline ignores with a `--no-inline-ignores`
   escape. Suppression must be first-class or the tool gets deleted.
8. **Gate the high-value rules nobody else gates.** Both are cheap regex/AST work and both are
   differentiators: **invented metrics** (hallmark gate 46 — "+47 % conversion", "trusted by 50,000+
   teams", "10× faster" with no source; fix = `—` + labelled block or a different macrostructure;
   *"A stat is also never the hero's sole headline"*) and **re-drawn UI chrome** (gate 47 — fake
   browser bar with URL pill + traffic-light dots, fake phone notch, fake terminal frame, fake IDE
   tabs).
9. **Gate token discipline mechanically** (hallmark gate 48 "Mid-render token improvisation": any
   `#hex` / `oklch()` / `rgb()` / `hsl()` / bare `font-family` outside the `:root` token block is a
   fail). This is a pure regex and it is the highest-yield single rule in the entire survey —
   `ds-lint`, `stylelint-declaration-strict-value`, and the Carbon/Kong/Mozilla stylelint plugins all
   exist to enforce the same thing, so delegation targets are available.
10. **Make the check theme-relative, not absolute** (mattbx: per-theme expectations for Vega/Nova/
    Maia/Lyra/Mira). Our checker must read the design system parley-design produced and validate
    *against it*, not against a universal taste.
11. **Keep a "Live Environment First" tier** (OneRedOak): the static checker is tier 1; a rendered-page
    tier (viewports 1440 / 768 / 375, or hallmark's stricter 320 / 375 / 414 / 768) catches what
    source cannot — horizontal scroll, two-line CTAs, cap-collision on wrapped all-caps display heads.
    Ship it as optional, gated on a browser being available.
12. **Ship the checker as both CLI and in-page module** (dataviz `validate_palette.js` runs via `node`
    *or* as `<script type="module">` reading `data-palette` off `<body>`, logging a `console.table`).
    Cheap to do, doubles the places the rule can fire.

---

## Do NOT copy

1. **hallmark's 20 named themes / 21 macrostructures / 50 component archetypes / 4 genres.**
   **[F]** 106 files, 9,591 markdown lines, 932 KB. **[I]** This is a *content library*, not a
   protocol. It is the opposite of vendor-neutral (it hard-codes one studio's taste: Specimen,
   Atelier, Brutal, Newsprint, Riso, Carnival…), it ages, and it must be maintained forever. Worse,
   in a multi-agent divergence protocol a fixed catalog **defeats the point** — four agents drawing
   from the same 20 themes converge, they don't diverge. Steal the *stamp*, the *log*, and the
   *gate list mechanics*; leave the catalog.
2. **hallmark's "always ask three questions before designing" gate.** **[F]** SKILL.md §1: *"There is
   no 'the brief looks complete' exception… A long, detailed brief gets the same three-question
   prompt as a five-word one."* **[I]** In Parley Deck the agents are headless — there is no human in
   the loop mid-run to answer. A mandatory interactive gate is an outright deadlock. Parley Deck's
   §9.0 preflight is where user input belongs, once, up front.
3. **Prose-only self-gates as the *only* enforcement** (hallmark, frontend-design, mattbx,
   artifact-design). **[I]** 58 unenforced gates is a 58-item wish list; the model that wrote the
   page grades the page. This is the exact gap that justifies parley-design-check existing.
4. **superdesign's hosted-CLI dependency** (`npm i -g @superdesign/cli`, `superdesign login`).
   **[I]** Login-walled, network-dependent, single-vendor, and a hard-fail in a sandboxed or offline
   headless agent. Violates the vendor-neutral requirement outright.
5. **impeccable's 23 commands.** **[I]** `bolder` / `quieter` / `distill` / `delight` / `overdrive` /
   `colorize` / `typeset` / `layout` / `adapt` / `optimize`… is a command surface, not a doctrine, and
   most of them are one prompt-line apart. Two skills with a clean interface beats 23 verbs. (Steal
   `detect`; steal the `init`→`shape`→`craft`→`audit`→`polish` spine; skip the rest.)
6. **impeccable's per-harness directory sprawl** (`.claude/ .codex/ .cursor/ .gemini/ .grok/ .kiro/
   .opencode/ .pi/ .qoder/ .rovodev/ .trae/ .trae-cn/ .vibe/ .agents/`). **[I]** Genuine
   vendor-neutrality is one canonical file plus documented adapters, not fourteen copies to keep in
   sync. Parley Deck already has a drift-guard scar from maintaining *two* copies of COOPERATION.md.
7. **vercel:shadcn's stack lock-in as doctrine.** **[F]** *"Default to dark mode for dashboards, AI
   apps, internal tools… Use Geist Sans for interface text and Geist Mono for code… Prefer zinc,
   neutral, or slate… `--radius: 0.625rem` is a strong baseline."* **[I]** These are good defaults for
   *Vercel*, and they are exactly the "AI default look" our anti-slop doctrine exists to resist.
   Steal the **"Reach for this first" table shape** and the **`validate:` frontmatter mechanism**;
   do not steal the aesthetic prescriptions.
8. **`validate:` frontmatter as our *only* enforcement.** **[F]** It is a Claude-Code/Vercel-plugin
   convention (`pattern` / `message` / `severity`), not a cross-vendor standard. **[I]** Use it as an
   *optional adapter* on top of a portable rule file; never as the source of truth.
9. **A "current AI clusters" list without a dated, revisable home.** **[F]** frontend-design says
   *"AI-generated design **right now** clusters around three looks"*; artifact-design already lists
   ~8. **[I]** These lists rot fast (2022 aurora blobs → 2024 purple gradients → 2026 cream+serif+
   terracotta). Ship them in a **versioned, dated, separately-updatable** rule file with an explicit
   review cadence — never welded into the doctrine prose.
10. **OneRedOak's `design-principles-example.md` as a rubric to adopt.** **[F]** It is the "S-Tier
    SaaS Dashboard Design Checklist" — 7 sections of unnumbered, unmeasurable bullets ("Meticulous
    Craft", "Strategic White Space", "Purposeful Micro-interactions") plus module-specific tactics for
    *that* product's multimedia-moderation and config-panel modules. **[I]** Un-checkable and
    domain-bound. Steal the *pattern* (the rubric is project-supplied data the reviewer reads), not
    the file.
11. **Figma-coupled workflow** (`figma-generate-library`'s 20–100+ `use_figma` calls, Plugin API
    gotchas, `combineAsVariants`). **[I]** Steal the phase/exit-criteria/task-ID **protocol shape**;
    the Figma mechanics are dead weight for a vendor-neutral skill.
12. **Star-count-driven prioritisation.** **[F]** Reported star counts across this survey ranged wildly
    for the same repos on the same day (hallmark cited as 4.6k / 11.2k / 19.1k; frontend-design as
    "65,847"; impeccable as 15k / 52k). **[I]** Third-party blog/aggregator numbers here are
    unreliable; only counts read directly off the GitHub repo page are used above, and even those are
    noisy. Do not use popularity to pick which doctrine to copy.

---

## Sources

Local files (absolute paths):
- `/Users/tomasfecko/.claude/skills/hallmark/SKILL.md` and `references/{slop-test,anti-patterns,design-md}.md`
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/frontend-design/unknown/skills/frontend-design/SKILL.md`
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/vercel/0.45.1/skills/shadcn/SKILL.md`
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/figma/2.2.81/skills/figma-generate-library/SKILL.md`
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/figma/2.2.81/figma-power/steering/create-design-system-rules.md`
- `/Users/tomasfecko/.claude/plugins/cache/claude-plugins-official/figma/2.2.81/skills/figma-implement-motion/references/motion-lint-rules.md`
- `/private/tmp/claude-501/bundled-skills/2.1.220/4fdadcd2f8243ffac061f198d7157521/dataviz/{scripts/validate_palette.js,references/anti-patterns.md}`
- `artifact-design` and `dataviz` SKILL.md bodies: embedded in `/Users/tomasfecko/.local/share/claude/versions/2.1.220`, read via the Skill tool.

Web:
- https://github.com/nutlope/hallmark
- https://github.com/pbakaus/impeccable · https://impeccable.style/docs/detector/
- https://github.com/OneRedOak/claude-code-workflows/tree/main/design-review
- https://github.com/superdesigndev/superdesign-skill · https://superdesign.dev/skill
- https://github.com/mattbx/shadcn-skills
- https://github.com/bergside/awesome-design-skills
- https://ui.shadcn.com/docs/changelog/2026-03-cli-v4
- https://vercel.com/blog/how-to-prompt-v0
- https://lib.rs/crates/ds-lint · https://github.com/carbon-design-system/stylelint-plugin-carbon-tokens · https://github.com/Kong/design-tokens/blob/main/stylelint-plugin/README.md · https://firefox-source-docs.mozilla.org/code-quality/lint/linters/stylelint-plugin-mozilla/rules/no-base-design-tokens.html · https://stylelint.io/awesome-stylelint/
