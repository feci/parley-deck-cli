# Impeccable (pbakaus/impeccable) — philosophy & product-doc digest

Studied at `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/research/impeccable`.
Version at time of read: skill/plugin `4.0.3` (`.claude-plugin/plugin.json`). License Apache 2.0. Author Paul Bakaus.
Read in full: `README.md`, `DESIGN.md`, `PRODUCT.md`, `AGENTS.md`, `CLAUDE.md`, `NOTICE.md`, `README.npm.md`, `.claude-plugin/**`, `.github/workflows/**`, `.github/hooks/impeccable.json`.
Read additionally because the actual doctrine lives there, not in the root docs: `skill/SKILL.src.md`, `skill/reference/{craft-floor,routing,new-work,critique,audit,polish,document}.md`, `skill/agents/impeccable-finish-reviewer.md`, `docs/STYLE.md`, `docs/HARNESSES.md`, `scripts/lib/transformers/{providers,factory}.js`, `scripts/lib/utils.js`, `cli/engine/registry/antipatterns.mjs`.

**Headline shape of the thing:** one user-invocable skill (`impeccable`) with 23 sub-commands, plus 60 deterministic detector rules, plus a per-provider build that emits the same payload into 14 harness directories. The root `CLAUDE.md`/`AGENTS.md` are *repo-development* instructions (how to build/test/release Impeccable), **not** the agent design doctrine. The design doctrine lives in `skill/SKILL.src.md` + `skill/reference/*.md`. That split is itself a finding.

---

## 0. Corrections to the brief's assumptions

- `DESIGN.md` at repo root is **not** a template and **not** a spec of the skill. It is Impeccable's **own product's** design system (the impeccable.style website), written in the third-party [design.md format spec](https://raw.githubusercontent.com/google-labs-code/design.md/main/docs/spec.md) that Google Labs' "Stitch" linter validates. It is a dogfooding artifact and a worked example of the output shape the skill's `document` command produces.
- The honesty/evidence/verification rules the brief hoped to find in `CLAUDE.md`/`AGENTS.md` are mostly in `skill/reference/critique.md`, `skill/reference/new-work.md`, `skill/agents/impeccable-finish-reviewer.md`, and `docs/STYLE.md`.
- The multi-runtime story is a **build-time transform** (source-first `skill/` → 14 generated `dist/<provider>/` trees), not a runtime abstraction.

---

## (a) DOCTRINE — design rules, taste, knowledge

### a.1 The stated theory of AI slop

From `README.md` (verbatim):

> "Every model trained on the same SaaS templates. Skip the guidance and you get the same handful of tells on every project: Inter for everything, purple-to-blue gradients, cards nested in cards, gray text on colored backgrounds, the rounded-square icon tile above every heading."

The framing is **distributional**: slop is not "ugly", it is *the mode of the training distribution*. Hence the skill's self-description in `skill/SKILL.src.md`:

> "This skill gives you the tools and permission to create design that earns to be called **out-of-distribution craft**: Whereas before, your design work would have been safe, timid and measured, you now approach every design task as a award-winning design director…"

Core principles, `skill/SKILL.src.md` (verbatim):

> - Go all out. No hedging, no shortcuts. The deliverable must be complete (except assets the user must provide).
> - Dream big and bold. Distinct, beautiful, outstanding and highly inspiring work.
> - Verify in bounded passes, not a loop, and the ceiling covers the whole cycle…

### a.2 The calibration rule — the sharpest anti-slop rule in the repo

`skill/reference/new-work.md` §4 (verbatim, `rule:skill-calibration-saturated-looks`):

> "AI-generated interfaces cluster around a few looks regardless of subject: warm cream ground, high-contrast serif display, and a terracotta or signal-red accent; near-black with one neon accent and glowing edges; broadsheet-editorial hairlines, italic display serif, and small tracked mono labels. All are legitimate when the brief calls for them; the brief always wins. Where the brief leaves the aesthetic free, landing in one of them means the self-check failed: **if someone could guess your aesthetic from the category alone, or from category-plus-avoidance, rework until neither answer is obvious.**"

Note the "category-plus-avoidance" clause: avoiding the cliché in the obvious way is itself a predictable second-order cliché. Related, same file:

- `rule:skill-book-subject-not-cream-license` — "A bookish, warm, or child-facing subject does not soften the calibration… landing on cream plus serif for a book subject is the default wearing the subject's clothes."
- `rule:skill-pinned-world-not-default-rendition` — "A brief-pinned world pins the world, not its softest rendition… a rendition that matches what any model ships for that world failed the self-check at execution rather than selection."
- `rule:skill-constraints-rule-out-devices-not-energy` — "a brief's negative constraints (no gamification, no hype) rule out those devices, not exuberance".

### a.3 Named-face denylist (typography)

`new-work.md` §4, `rule:skill-typo-reflex-faces` — "these training-data defaults mean you stopped looking":

> Fraunces, Playfair Display, Cormorant, Lora, Crimson, Newsreader, Syne, Space Grotesk, Space Mono, IBM Plex, Inter-as-display, DM Sans, DM Serif, Outfit, Plus Jakarta Sans, Instrument Sans.

> "Naming one of these faces anyway requires a reason no other face could satisfy, and a subject association is never that reason: books wanting a serif, bookshops wanting hand-lettering, and tech wanting a mono are the associations the list exists to break."

### a.4 The craft floor (`skill/reference/craft-floor.md`, 5.5 KB, loaded only immediately before editing UI)

**Verify** (each "a check on the built result, not an intention", run in *batched* rounds sharing one render):

| Axis | Threshold (verbatim) |
|---|---|
| Contrast | "body and placeholder text ≥4.5:1, large text ≥3:1. On colored surfaces tint secondary text from that hue or the foreground; never gray." |
| Depth | "shadows carry an offset and a soft blur. A zero-offset colored halo is decoration." |
| Spacing | "tight groups, generous separation, more space above a heading than below it. Read the computed values." |
| Type | "body measure 65–75ch, display max 6rem, tracking floor -0.04em, balanced headings, obvious scale and weight steps." |
| Motion | "one authored moment, not scattered effects and not one identical entrance on every section. Exponential ease-out from an already-visible default." |
| States | "hover, disabled, loading, error, empty. Plus real content, working controls, responsive composition, keyboard focus." |
| Copy | "the product's own language. Controls name their action; errors name the problem and the recovery." |
| Coverage | "every brief requirement present and findable within seconds." |

**Refuse** — note the meta-rule framing, which is the interesting bit:

> "These are the category's defaults, not bans: the brief's own words can earn any of them. **Reaching for one when the axis is free means you were not deciding**; recognizing that means rewriting the element, not softening it."

Page scaffolds refused: identical icon+heading+text card grids ("Cards are the lazy container; nested cards are always wrong"); the hero-metric template; **kicker/eyebrow above a heading — "This one is a ban, not a default: no brief earns it back"**; section numbers `01/02/03` unless the sequence carries information; a modal for a task needing neither interruption nor protected focus.
Surface habits refused: gradient text; glass/blur as decoration; colored `border-left`/`border-right` above 1px; sparklines/progress rings/soft-shadowed rounded rectangles standing in for content; monospace "as a costume for 'technical'"; light-or-dark "picked by category… Pick it from the use scene: who, where, under what ambient light."

Closing line (`rule:skill-floor-not-ceiling`):

> "The floor holds the mechanics; it never picks the direction. With every check green, spend the page on the committed world, and when torn between refined and committed, **commit**."

### a.5 Per-model bias corrections (this is unusual and directly relevant to a multi-agent roster)

The craft floor and new-work carry **model-specific** blocks compiled in only for that provider:

`<codex>` block in `craft-floor.md` — GPT-family tells: tracking floor repeat, "Declare elevation once, border or shadow. A 1px border under a wide soft shadow is the **ghost card**. Card radii stay at 12–16px"; "Real illustration or none. Sketch-style SVG scenes, `loose-sketch` / `doodle` class names, and `feTurbulence` grain read as amateur"; `repeating-linear-gradient` stripes and grid overlays need a real canvas/map/blueprint under them; "Claims and configuration come from supplied truth; label illustrative values honestly."

`<gemini>` block — "Never animate an image on hover, directly or through its parent. It is not an action target. Give the container the feedback."

`<claude>` block in `new-work.md` §4 (verbatim):

> "Your measured rendition prior: warm, bookish, family, and child-facing subjects come out as cream grounds, serif display with italic accents, and lamplight, even when the assigned direction never asked for them. **Treat that first palette as already spent.** Before writing code, reread your OWN-WORLD block: when it says cream, paper, parchment, ivory, or lamplight for a Persuade surface the brief did not pin, the rendition failed and you rework it from the world's saturated materials first. The same subject renders as bookcloth, thread, jacket, and endpaper color on other models; nothing about the subject requires your default."

Four of the 60 detector rules are explicitly model-attributed: `gpt-thin-border-wide-shadow`, `repeating-stripes-gradient`, `codex-grid-background`, `theater-slop-phrase`, `image-hover-transform`.

### a.6 Two orthogonal axes replace one "register"

**Modes** (per *surface*, not per project) — `SKILL.src.md`:
- **Persuade** — "the visitor decides and acts; design is the product." Landing, marketing, pricing.
- **Operate** — "the visitor completes a task." App UI, dashboards, admin. "Scanability, consistency, native expectations, and the real usage scene outrank expression."
- **Read** — "the visitor understands something." Docs, guides, changelogs.
- **Experience** — "the visitor is inside the work itself." Portfolios, galleries.
"A tool's landing page is still Persuade; a fashion house's documentation is still Read; a docs index is Read, not Persuade."

**Platform** (`web` / `ios` / `android` / `adaptive`) — orthogonal; loads `reference/ios.md` (Apple HIG distilled) / `reference/android.md` (Material 3 distilled), both derived from MIT-licensed `ehmo/platform-design-skills` with attribution in `NOTICE.md`. `adaptive` loads **both**. Missing field defaults to `web`; an unrecognized value falls back to web **and** prints a `WARNING` directive naming the bad value "so a toolchain name or typo never silently gets web guidance."

### a.7 Truth / honesty doctrine about the artifact itself

`new-work.md`, `rule:skill-truth-binds-claims` (verbatim):

> "Truth binds claims, not demonstrations: in greenfield work, author whatever illustrative material the concept needs at full fidelity, **label it synthetic wherever a visitor could mistake it for the real thing**, and hand the user the list of what to replace with real material. What stays uninventable are commercial and factual claims: prices, customers, benchmarks, endpoints, capabilities the product does not have. **Refusing a bold direction because its demonstration data does not exist yet is the timidity reflex wearing honesty's clothes.**"

`new-work.md` §6: "**Prove, don't claim.** Show the subject doing its job… Sections that restate a claim in different words add length, not substance."
`init.md` never invents "testimonials, customers, benchmarks, pricing, licensing, or deployment claims"; PRODUCT.md has an `## Evidence on Hand` section that must "State absences that future work must not fabricate."

### a.8 Prose doctrine — `docs/STYLE.md` (the anti-AI-writing brief)

The bar (verbatim): **"for every paragraph, point to the sentence that makes it specifically yours. If you can't, the paragraph is AI by default, even if a human typed it."**

12 principles, sharpest ones verbatim:
1. "Open with the reader's wrong belief, your strongest claim, or the example."
2. "Take a position someone could disagree with. **If the paragraph could be inverted without changing meaning, it has no position.**"
3. "Name names. Use numbers… Cut 'lightweight'; write '54 KB'."
5. "Vary sentence length on purpose. Long, long, short. **Smooth uniform rhythm is the deepest AI tell.**"
11. "Concrete over comprehensive. **Coverage is an AI obsession.** Trade coverage for momentum. Leave things out."
12. "Close by handing off the next move. Don't summarize."

Enforced denylist (build fails): `load-bearing`, `highest-leverage`, `biggest unlock`, `reflex defaults`, `collapses into monoculture`, `data-driven`, `seamless`, `robust`, `elevate`, `empower`, `underscore`, `pivotal`, `tapestry`, `delve`, `in today's`, `gone are the days`, `whether you're`, `let's dive in`, `in summary`, `in conclusion`, `moreover`, `furthermore`, the em dash `—` (+ `&mdash;`/`&#8212;`/`&#x2014;`) and ` -- `. Em-dash rationale verbatim: **"Decision-avoidance: writer didn't pick a relationship between the clauses."** Each denylist row carries a *Why* and a *Use instead*.

Patterns the validator can't catch (require human judgment): **negation pivot** ("It's not just X, it's Y" — "now a stronger AI tell than any vocabulary item"), **triadic everything**, **the five-paragraph essay shape**, **uniform paragraph length**, **synthetic balance**, **hollow confidence**, **hedging stacks**, and **interchangeable copy**: "Swap 'Impeccable' for a competitor name. If nothing becomes false, the copy is generic."

### a.9 PRODUCT.md as the durable-truth artifact

Root `PRODUCT.md` is the dogfood instance and shows the schema: `## Register` (legacy, now deprecated), `## Users`, `## Product Purpose`, `## Brand Personality`, `## Anti-references`, `## Design Principles`, `## Accessibility & Inclusion`. Its own design principles, verbatim:

> 1. **Practice what you preach.** The site must pass its own anti-pattern tests with flying colors. If we ship anything we'd flag in an audit, we've lost.
> 2. **Show, don't tell.** … The site IS the demo.
> 3. **Expert confidence.** Direct, opinionated, decisive. No hedging.
> 4. **Editorial over marketing.** Feels like a design publication (Eye Magazine, It's Nice That, A List Apart) rather than a SaaS landing page.
> 5. **Purposeful restraint.** Every element earns its place.

Anti-references are a **named artifact field**, not an afterthought: "dark mode with purple gradients, neon accents, glassmorphism, glowing particles, cyan-on-black"; "hero-metric layouts, identical-card feature grids, sparkline decorations, 'boost your productivity' copy"; "Hedging language: 'might', 'could', 'consider', 'perhaps'"; "Over-decoration".

---

## (b) PROCESS — workflow, phases, state

### b.1 What DESIGN.md actually is

Root `DESIGN.md` = "Design System: Impeccable" (creative north star **"Neo Kinpaku"**). Structure:

- **YAML frontmatter** (machine-readable tokens; lines 1–249): `name`, `description`, `colors` (~90 entries, all OKLCH), `typography` (an enumerated px→rem `scale` from `"8": 0.5rem` to `"88": 5.5rem`, plus named roles `wordmark`/`display`/`headline`/`title`/`body`/`eyebrow`/`mono`), `rounded` (13 steps `none`→`pill: 999px`), `spacing` (7 steps `xs: 8px` → `3xl: 112px`), `components` (9 entries using `{colors.x}` token refs).
- **Markdown body**, 7 numbered sections: `1. Overview: Neo Kinpaku` · `2. The Kit: One Vocabulary For Every Page` · `3. Colors: Lacquer, Gold, Patina` · `4. Typography` · `5. Elevation and Material` · `6. Components` · `7. Do and Do Not`.

Two structural devices worth stealing:

1. **A source-of-truth pointer inside the frontmatter comment**: "All values below mirror `site/styles/kinpaku-tokens.css` verbatim. That file is the source of truth; this frontmatter is the portable export. If a token changes there, update both."
2. **Named, quotable rules** rather than prose. Every subsection ends in "The X Rule": *The Gold Carries Brand Rule*, *The Patina Has Meaning Rule*, *The Texture Budget Rule*, *The OKLCH-Only Rule*, *The Gold-By-Size-On-Paper Rule*, *The Weight-Inversion Rule*, *The Two-Face Rule*, *Tracked Labels Are Short Rule*, *Dark Type Needs Air Rule*, *Hairline First Rule*, *No Glass Rule*, *Texture Needs Contrast Rule*, *Asset-Led Material Rule*, *The Kit Consumption Rule*, *The Picker Is Brand Rule*. Each is one sentence a reviewer can cite.
3. **Rationale embedded as a number**: the type-scale comment — "The tail this ramp replaced had **86 distinct sizes**, including six near-identical steps between 13.7px and 15.4px that no reader could tell apart. Adding a step is a design decision, not a convenience."
4. **§7 "Do and Do Not"** — 8 Dos, 10 Do-Nots, each one line, each falsifiable against the built page.

The **canonical** DESIGN.md schema (from `skill/reference/document.md`) is the external design.md spec: optional YAML frontmatter + **up to eight markdown sections in fixed order**: `## Overview`, `## Colors`, `## Typography`, `## Layout`, `## Elevation & Depth`, `## Shapes`, `## Components`, `## Do's and Don'ts`. "**Tokens are normative; prose provides context for how to apply them.**" Component sub-tokens limited to 8 props (`backgroundColor`, `textColor`, `typography`, `rounded`, `padding`, `size`, `height`, `width`); anything else goes to a sidecar. "Skip anything the project doesn't have. Empty scale keys or fabricated tokens pollute the spec." "If a `DESIGN.md` already exists, **do not silently overwrite it**."

### b.2 The design-time ordering rule that matters most

`new-work.md`, `rule:skill-design-md-from-the-build`:

> "On a new or replacement world, DESIGN.md is written **at finish, from the built world**, by the shipped documenter; **a rulebook written before the build gets defended against reality instead of describing it**, and it hands the design-system detector an unstable target that buries the build in noise. A new world shipped with no DESIGN.md is still an incomplete run."

### b.3 The direction contract (design decisions frozen as a 150-word artifact header)

`new-work.md` §5, `rule:skill-decide-then-build` — before any code, five blocks, ≤150 words, in an **HTML comment in the emitted markup** ("never only a templating-frontmatter comment the compiler strips"):

- **THESIS** — "the one idea this surface owns and the category-default arrangement it refuses."
- **OWN-WORLD** — "the palette and component language, specific enough to be recognizable with all content removed."
- **STORY** — "what the visitor understands, believes, and does."
- **FIRST VIEWPORT** — "the exact composition, what is where and at what scale, and where the primary action sits."
- **FORM** — the chosen form, its rank on the ordered list, the staging, and **the seed key the concept-roll script printed**.
- plus **FINISH** — verbatim exit condition: *"unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, and DESIGN.md"*.

> "If a block reads like a mood, the direction is not decided yet; the finishing review audits the render against this contract." … "a page that looks complete with the FINISH line undischarged is not done, it is abandoned at the finish line."

### b.4 The anti-monoculture selection procedure ("the roll")

The single most transferable *process* idea. `new-work.md` §3:

1. Name the mechanism, audience scene, cultural home, and what the first surface must prove. **Name "the rut"**: "the page this category always ships and its predictable opposite; name both as the rut and keep them out of the seven-candidate list, so no die face is spent on the page the category already ships."
2. Derive **seven** concrete visual systems/artifacts/places/rituals from the audience's world, ordered by resonance. "**When more than three of the seven share one material family, the derivation stopped at the subject's most obvious artifact**… dig until the list spans at least three families."
3. Turn each into a complete direction (world + concrete first-surface experience).
4. Run `node scripts/concept-seed.mjs --scope direction --mode <mode>`. **The script assigns which direction gets built and deals catalog challengers.** Rationale (verbatim): "your top-ranked structure is what every run would ship, and **a single ranking is deterministic, so the dice come from outside**."
5. Present **one** direction fully committed + at most three dealt challengers with one-line cases including an honest "fuses poorly because X". "**What you never present is a ranked menu of your own grounded candidates; a lineup of those invites the safest card**, while dealt challengers carry no such rut."

Counterweights: **the standing exit** — every round offers the category standard played straight, "It is the user's door, never yours: never recommend it, never weigh it against the roll." Re-roll eliminates everything shown; after two consecutive re-rolls, ask what quality is missing. "You may re-roll on your own only on named factual grounds… **taste is never grounds.** The user may re-roll freely, and a user- or brief-pinned direction beats the roll, always."

Also: comparison sketches must use "**one shared frame**… identical framing across cards; **a candidate whose sketch looks more finished than the others has broken the comparison, not won it.**" Only legible text in a sketch is the real product name and one real headline; everything else greeked, because "a sketch that renders invented specs, prices, or dates puts claims in front of the user that PRODUCT.md never made."

### b.5 The bounded verification loop (explicit anti-infinite-polish)

`SKILL.src.md`: "Verify in **bounded passes, not a loop**… Build fully, inspect once with a batched round (desktop and mobile together), fix everything it shows in one batch, confirm with at most one more round, and stop polishing. **Open-ended self-QA burns the user's money doing worse what the finish handoffs do better.**"

`new-work.md` §7 (`rule:skill-verdict-bounds-the-finish`, `rule:skill-finish-separate-reviewer`) — the full finish protocol:

- Inspect desktop+mobile in **one batched screenshot round**; fix; confirm with one final round. **Two rounds is the ceiling.**
- "After the second inspection round **the build thread's polishing is over**: no further defect hunts, micro-edit scripts, or rebuilds here; whatever remains ships through the handoffs, where a fresh context does the finding better and cheaper."
- Spawn `impeccable-finish-reviewer` with request, confirmed answers, artifact path, screenshot paths, direction contract, hook findings, QUALITY BAR card and approved comp paths. "**This review never runs inside the build thread.**"
- "Verify its return carries the five contract sections; on an empty or thrashed return, **respawn once** with the same inputs before doing anything else."
- Apply material fixes **in one batch**, rebuild once, recapture the same viewports. "**A recapture measures positions, loading, and overflow; it cannot measure whether a fix reached the quality the finding named**", so send the recaptures back to the same reviewer for a **verdict** scoring each fix `resolved | partial | unresolved`.
- Partial/unresolved get exactly one more batch → recapture → verdict. "**Two correction rounds is the ceiling, the second verdict ends the work whatever it says**, and the reviewer's findings are the only list you work from, never your own re-opened hunt."
- "Report the final verdict table to the user as it stands, open items included: **presenting mechanical confirmation as artistic success is how a failed build gets announced as a finished one.**"
- Then spawn `impeccable-documenter` to write DESIGN.md "from the built world, **ground truth over intention**". "A clean detector pass is not finished; finished is the contract kept, the comp honored, the review closed, and the system recorded."

### b.6 Isolated dual assessment + mandatory degradation banner (`skill/reference/critique.md` Hard Invariants — verbatim)

> - Assessment A (design review) and Assessment B (detector/browser evidence) are both required.
> - Assessment A and B **MUST run as two isolated sub-agents** whenever a sub-agent/Task tool is exposed. Running them inline in this context is "possible" but is **NOT permitted**; it is a degraded run.
> - If you degrade for any reason, the report's first line MUST be a banner: `⚠️ DEGRADED: single-context (<reason>)`. **A silent degraded critique is a failed critique.**
> - Assessment A must finish before detector findings enter the parent synthesis context. **Detector output is deterministic, but it still anchors judgment.**
> - A skipped detector is a failed critique run unless `detect.mjs` is missing or crashes after a real attempt.
> - Do not claim a user-visible overlay exists unless script injection succeeded and the detector ran in the page.

Plus: "'Unavailable' means exactly one thing: no sub-agent/Task tool is exposed in this session (or… the user declined). **It does not mean inconvenient.**" And a required provenance header line: `Method: dual-agent (A: <agent-id> · B: <agent-id>)` or the degraded banner. "Skipping sub-agents without the banner is the most common failure of this command."

Synthesis rule: "Do NOT simply concatenate. Weave the findings together, noting **where the LLM review and detector agree, where the detector caught issues the LLM missed, and where detector findings are false positives.**"

### b.7 Scoring rubrics (calibrated, with anti-inflation language)

**Critique** — Nielsen's 10 heuristics, 0–4 each. "Be honest with scores. A 4 means genuinely excellent. **Most real interfaces score 20-32 out of 40.**" Mode-applicability: heuristics 7 and 10 may be `n/a` on Persuade/Experience surfaces; the total renormalizes ("**Never print `/40` over a partial set**"), and the snapshot must record the applicable max and which were n/a.

**Audit** — 5 dimensions × 0–4 = /20: Accessibility, Performance, Theming, Responsive Design, **Implementation Integrity (CRITICAL)**. Bands: `18-20 Excellent · 14-17 Good · 10-13 Acceptable · 6-9 Poor · 0-5 Critical`. Report opens with the **Implementation Integrity Verdict** — "**Start here.** Pass/fail: does the implementation express a coherent product-specific system? Cite verified evidence and detector findings."

**Severity, shared** (P0–P3):

| Priority | Name | Description | Action |
|---|---|---|---|
| P0 | Blocking | Prevents task completion entirely | Fix immediately; showstopper |
| P1 | Major | Causes significant difficulty or confusion | Fix before release |
| P2 | Minor | Annoyance, but workaround exists | Fix in next pass |
| P3 | Polish | Nice-to-fix, no real user impact | Fix if time permits |

Tie-breaker heuristic: "**Would a user contact support about this?** If yes, it's at least P1."

Audit's NEVER list: report issues without impact; generic recommendations; skip positive findings; forget to prioritize ("everything can't be P0"); "Report false positives without verification."
Critique's tone rules: "Don't soften criticism." "Cut 'consider exploring…' entirely." "Prioritize ruthlessly. If everything is important, nothing is."

### b.8 Reviewer role definition (`skill/agents/impeccable-finish-reviewer.md`)

Design points worth lifting wholesale:

- **Anti-anchoring**: "inventory the comp's salient elements **in your own words before** reading the direction contract or any builder-authored summary: the contract is the builder's abstraction of the comp, and **a review anchored on it inherits whatever that abstraction dropped**."
- **Budgeted reading**: "You run under a hard turn ceiling that ends the run without warning, and a run that ends before the five sections are written returns nothing; **a review built from what you saw beats a perfect review that never arrives.** So treat reading as an allowance, not a prerequisite… by roughly the tenth turn stop reading and write. **Name whatever went unread in the line above the sections.**" (`max-turns: 30`, `effort: high`, `tools: Read, Bash, Glob, Grep`, `model: inherit`.)
- **Explicit capability disclosure**: "You have no browser. Never attempt to render, screenshot, start a server, or open a page… When an expected input is missing, say so in one line at the top."
- **Fixed 5-section output contract**: `persistence` · `fidelity` · `ceiling` · `material_fixes` (≤8, ordered, fidelity failures ahead of craft) · `keep` ("one line naming what must not be diluted while fixing"). "**No praise, no summary prose.**"
- **Classification vocabulary**: every salient element is `match | acceptable adaptation | missing | contradicted | added without approval`. "An adaptation counts as intentional **only when it cites** the user answer, surface brief, accessibility need, or product truth that forced it; **an uncited deviation is a defect.**"
- **Two mandatory rows** in every matrix: TYPE and MATERIAL. "a face of a different character is contradicted however the layout matches"; "an element rendered as flat CSS or clean vector where the comp shows painted, textured, dimensional, or photographic material is contradicted regardless of placement, because **medium is part of the promise**."
- **Anti-token compliance**: "an asset applied at near-zero opacity or buried behind other paint is a **compliance token**, not a shipped material."
- **Verdict pass is scoring, not re-hunting**: "you are scoring, not re-hunting… **a fix answered mechanically, positions moved but the quality the finding named still absent, is partial at best.** Then name at most three regressions the fix batch itself introduced… **no new hunt, no new checks.**"
- **Non-duplication of machinery**: "Do not run a second detector pass; mechanical findings belong to the parent's hooks."

### b.9 Session/state model

- `PRODUCT.md` (durable product truth, root) — carries a provenance stamp `<!-- impeccable:product-schema N -->`. "**Stamps are schema versions, not release versions**". Retired fields go into `PRODUCT_DEPRECATED_SECTIONS` **with a reason**: "told only that a field is deprecated, models preserve it 'just in case', which is how a retired axis keeps steering current output."
- `DESIGN.md` (durable visual decisions, root) — deliberately carries **no** stamp, because it follows an external spec and every staleness signal is measurable without one.
- `.impeccable/surfaces/<...>` — per-surface briefs (mode lives here, not in PRODUCT.md).
- `.impeccable/critique/*.md` — persisted critique snapshots, consumed as a backlog by `polish` via `critique-storage.mjs latest <target>`. `.impeccable/critique/ignore.md` is "**the only prior-run input critique consumes**".
- Tracked vs ephemeral is documented in README with a `# impeccable-ignore-start` / `# impeccable-ignore-end` gitignore block; patterns deliberately **unanchored** for monorepos.
- Staleness: **Tier 1** boot findings (`collectBootFindings()`) with a hard performance contract — "**No directory walks, no git, no cross-workspace sweep**"; **Tier 2** on-demand `doctor.mjs`. Findings are data `{ id, artifact, path, severity, summary, fix }` where severity means *what should happen*: `auto` (fix silently on next write) / `mention` (state once) / `route` (name the owning command). `doctor --fix` applies only `auto`. Emission throttled to once a week per project.
- **`Never repair drift as a side effect of a design task.`** (`rule:skill-drift-not-a-side-quest`)

---

## (c) MACHINERY — scripts, detectors, enforcement, tooling

### c.1 The detector: 60 deterministic rules, no LLM, no API key

`cli/engine/detect-antipatterns.mjs` (thin facade over `cli/engine/{registry,rules,engines,shared}`) is the single source of truth. Two categories: `slop` (AI tells) and `quality` (real design/a11y issues). Exit codes: `0` clean, `2` findings.

Full rule roster as extracted from `cli/engine/registry/antipatterns.mjs` (60 entries):

**slop (32)** — `side-tab`, `border-accent-on-rounded`, `overused-font`, `single-font`, `flat-type-hierarchy`, `gradient-text`, `ai-color-palette`, `cream-palette`, `nested-cards`, `monotonous-spacing`, `bounce-easing`, `pulsing-dot`, `blinking-cursor`, `shape-assembled-illustration`, `dark-glow`, `radial-halo`, `radial-spotlight-glow`, `marquee`, `icon-tile-stack`, `italic-serif-display`, `hero-eyebrow-chip`, `kicker-above-heading`, `numbered-section-labels`, `marketing-buzzword`, `aphoristic-cadence`, `oversized-h1`, `extreme-negative-tracking`, `gpt-thin-border-wide-shadow`, `repeating-stripes-gradient`, `codex-grid-background`, `theater-slop-phrase`, `image-hover-transform`.

**quality (28)** — `broken-image`, `script-error`, `content-hidden-at-rest`, `edge-flush-cards`, `text-occlusion`, `first-viewport-column-overflow`, `gray-on-color`, `low-contrast`, `layout-transition`, `line-length`, `cramped-padding`, `body-text-viewport-edge`, `tight-leading`, `skipped-heading`, `heading-rhythm`, `justified-text`, `tiny-text`, `undersized-ui-text`, `all-caps-body`, `wide-tracking`, `text-overflow`, `repeated-container-text`, `clipped-overflow-container`, `design-system-font`, `design-system-color`, `design-system-radius`, `design-system-font-size` (+ small-touch-target family).

Note the last four: **`design-system-*` rules make DESIGN.md itself machine-enforceable** — a font/color/radius/font-size outside DESIGN.md is a finding.

Rule entries carry optional `skillSection` / `skillGuideline` fields — i.e. **a mechanical finding can point back at the prose rule it enforces**.

CLI surface: `npx impeccable detect src/ | index.html | https://example.com | --json . | --no-config src/`, plus `npx impeccable ignores list | add-file "src/legacy/**" | add-value overused-font Inter --reason "Brand font"`.

Waivers are three-tier: repo config (`.impeccable/config.json` → `detector.ignoreRules|ignoreFiles|ignoreValues|designSystem.enabled`), per-developer (`config.local.json`, gitignored), and **inline file-scoped comments** that travel with the file: `<!-- impeccable-disable overused-font: exported brand doc -->` (also `impeccable-disable-line` / `impeccable-disable-next-line`), bypassable with `--no-inline-ignores` / `--no-config`. **Every waiver carries a reason.**

### c.2 TDD order for a new rule — "non-negotiable" (`CLAUDE.md` + `AGENTS.md`, identical)

1. **Fixture** at `tests/fixtures/antipatterns/{rule-id}.html`, two columns (should-flag / should-pass), each case a unique heading. **≥4 flag cases and ≥5 false-positive shapes.** Explicit pixel dimensions in CSS (jsdom does no layout).
2. **Failing test** in `tests/detect-antipatterns-fixtures.test.mjs`, snippet-substring pattern (regex `/"([^"]+)"/` against `SHOULD_FLAG` / `SHOULD_PASS`). "Run it and watch it fail before implementing."
3. **Rule entry** (`id`, `category`, `name`, `description`, optional `skillSection`/`skillGuideline`).
4. **Pure check** `checkXxx(opts) → [{ id, snippet }]` — "No DOM access in the pure function."
5. **Two adapters** (browser `getComputedStyle`+`getBoundingClientRect`; jsdom `parseFloat(style.width)`), wired into **both** element loops. "Forgetting one loop is the most common mistake; symptom is 'test passes, live page silent' or vice versa."
6. **Verify on a live page** at `http://localhost:4321/fixtures/antipatterns/{rule-id}.html` and on the homepage (no false positives). "The two adapter paths can disagree."

Five artifacts must stay in sync per rule change; the required command is `bun run build && bun run build:browser && bun run build:extension && bun run test`.

### c.3 The hook — enforcement at edit time

Installed per-provider by `npx impeccable install` / `update`:

| Harness | Event | Manifest | Script |
|---|---|---|---|
| Claude Code | after edit (+Stop deep pass) | `.claude/settings.local.json` (or `settings.json` in place) | `${CLAUDE_PROJECT_DIR}/.claude/skills/impeccable/scripts/hook.mjs` |
| GitHub Copilot | `postToolUse`, matcher `edit\|create\|apply_patch`, `timeoutSec: 5` | `.github/hooks/impeccable.json` (committed) | `.github/skills/impeccable/scripts/hook.mjs` |
| Cursor | `preToolUse` — **blocks bad proposed writes before they land** | `.cursor/hooks.json` | `hook-before-edit.mjs` |
| Codex | `postToolUse` | `.codex/hooks.json` | `.agents/skills/impeccable/scripts/hook.mjs` |
| Grok Build | `postToolUse` + Stop | `.grok/hooks/impeccable.json` | — |
| everything else | none documented | — | skill + commands still ship |

Installer discipline: preserves unrelated hook entries; "If a hook manifest is malformed, install/update **aborts by default**; rerun with `--force` to back up the malformed file as `.bak`". Debug: `hook.auditLog` writes one NDJSON line per invocation.

Craft-floor defers to the hook: "When the design hook is active it already enforces the mechanical checks below as you edit: **act on its findings instead of re-auditing each rule.**"

### c.4 Prose enforcement in the build

`scripts/build.js` has **two** prose gates:
- `validateProse` — scans `README.md`, `README.npm.md`; full denylist; "Each rule prints a rationale and a suggested replacement when it fires. **Do not silently work around the regex.** If a banned word has earned a real meaning here, raise it as a `docs/STYLE.md` amendment."
- `validateSkillProse` — scans `skill/**/*.md` (markdown only, not `skill/scripts/**`) with a **narrower** subset: em dashes plus `load-bearing`, `highest-leverage`, `biggest unlock`, `reflex defaults`, `collapses into monoculture`, `data-driven`, `delve`, `tapestry`, `in today's`, `gone are the days`, `let's dive in`, `in summary`, `in conclusion`. Excluded words (`seamless`, `robust`, `elevate`) "are the ones with legitimate technical uses". Net: "an em dash in `skill/reference/*.md` fails `bun run build`; an em dash in a `skill/scripts/*.mjs` code comment does not."

### c.5 Count/version drift validators

- `generateCounts` fails the build if the "23 commands" number disagrees between the router table, `README.md`, `AGENTS.md`, `.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json`, and the site.
- `validatePluginVersions` fails if `marketplace.json`, the `./plugin` manifest, and the bundled `SKILL.md` frontmatter disagree with `plugin.json` ("this guards the marketplace install path against version drift (issue #274)").
- Rule markers: every instruction line in the skill source carries `<!-- rule:some-id -->`; `stripRuleMarkers()` removes them at staging. Rationale in `scripts/lib/utils.js`: "External eval tooling can pin each instruction line to a stable ID. Markers in the source keep that mapping verifiable in lock-step with the file."

---

## (d) MULTI-RUNTIME PORTABILITY — how one payload reaches 14 harnesses

**One source, N generated trees.** `skill/` is the only authoring surface: `SKILL.src.md` + `reference/*.md` + `scripts/*` + `agents/*`. `bun run build` emits `dist/<provider>/`; `bun run build:release` additionally syncs the **tracked** root harness folders (`.agents/ .claude/ .cursor/ .gemini/ .github/skills .grok/ .kiro/ .opencode/ .pi/ .qoder/ .rovodev/ .trae/ .trae-cn/ .vibe/ plugin/`).

Five portability mechanisms, in the order they run (`scripts/lib/transformers/factory.js` `createTransformer`):

1. **Per-provider config table** (`providers.js`) — 14 entries, each declaring `configDir` (`.claude`, `.codex`, `.agents`, `.github`, …), `displayName`, `providerTags`, `frontmatterFields`, `agentFormat`, `emitHooks`, `hooksManifestRel`, `placeholderProvider`.
2. **Conditional markdown blocks** (`compileProviderBlocks`) — `<codex>…</codex>`, `<claude>…</claude>`, `<gemini>…</gemini>` etc. Matching blocks keep their body and drop the tags; non-matching blocks are **removed entirely**; **unknown tags are preserved** so ordinary HTML is untouched. Allowed tags: `agents, claude, claude-code, codex, cursor, gemini, github, grok, kiro, opencode, pi, qoder, rovo-dev, trae, trae-cn, vibe`. Aliasing is explicit: Claude Code's `providerTags` is `['claude-code','claude']`; `agents` gets `['agents','codex']`.
3. **Placeholder substitution** (`replacePlaceholders`, `PROVIDER_PLACEHOLDERS`) — six tokens: `{{model}}`, `{{config_file}}`, `{{ask_instruction}}`, `{{command_prefix}}`, `{{available_commands}}`, `{{scripts_path}}`. Examples: Claude Code → `Claude` / `CLAUDE.md` / "STOP and call the AskUserQuestion tool to clarify." / `/`; Codex → `GPT` / `AGENTS.md` / "STOP and use Codex's structured user-input/question tool when available; if unavailable, ask directly in chat…" / `$`; OpenCode → "STOP and call the `question` tool"; Cursor/Gemini/Pi/Qoder/Trae/Rovo/Vibe → the generic "ask the user directly to clarify what you cannot infer." When `command_prefix !== '/'`, a lookbehind-guarded regex rewrites `/impeccable` → `$impeccable` **without** touching paths and URLs (`.github/hooks/impeccable.json` stays intact). **Placeholders are never run across JS source** — scripts get a single explicit marker via `replaceScriptProviderMarker`.
4. **Frontmatter field gating** — each provider declares which optional fields it can parse (`user-invocable`, `argument-hint`, `license`, `compatibility`, `metadata`, `allowed-tools`). Gemini and Codex get `frontmatterFields: []`. `docs/HARNESSES.md` holds the full support matrix (12 harness columns × 15 fields, values `Yes / No / Ignored / TBD`) plus a directory-fallback matrix (e.g. Cursor also reads `.agents/skills/` and `.claude/skills/`; Grok also reads `.agents/`, `.claude/`, `.cursor/`).
5. **Subagent emission + degraded single-sourcing** — the most transferable trick. Agents live once in `skill/agents/*.md`. They are emitted as (i) Claude-style markdown+frontmatter for `agentFormat: 'claude-md'` harnesses, (ii) Codex TOML (`name`, `description`, `model_reasoning_effort`, `nickname_candidates`, `developer_instructions = '''…'''`) nested **inside** the skill's own `agents/` folder for Codex, and (iii) **`reference/degraded/<role>.md`** for every provider — the same agent text prefixed with a generated preamble:

> "This harness has no subagent capability, so you are running this role inline. **Step fully out of the work you just finished**, adopt only this file's instructions for the pass, and **disclose the substitution in one line** when you report. Where the text below addresses a parent agent, **you are both parties**: produce the full output contract first, then act on it yourself."

An agent may declare `providers: [...]` to limit emission. Degraded files pass through the identical block-compile + placeholder pipeline, "so `<codex>` blocks and `{{placeholders}}` resolve identically."

Runtime capability handling is done **in the prose**, not in code: routing skips `live` and `detect.mjs` on native platforms; critique's `<codex>` block overrides the default sub-agent gate because "Codex's permission model requires asking before spawning" and explicitly states "**Asking is the normal path, not a degradation**"; the visual decision page self-detects headless/CI/remote-shell and "exits 2 with that advice, so **treat exit 2 as this fallback, never as an error to retry**."

Install/link surfaces: `npx impeccable install` (auto-detects harness folders, `--providers=`, `--scope=project|global`, `--no-hooks`), `npx impeccable update`, `npx impeccable link --source=.impeccable --providers=claude,cursor` (git-submodule vendoring; "leaves existing real skill directories untouched unless you pass `--force`"), Claude Code plugin marketplace, Grok `grok plugin install pbakaus/impeccable#plugin --trust`, website ZIPs, and raw `cp -r dist/<provider>/…`.

---

## (e) TESTING / CI — what is actually verified automatically

Local commands (`AGENTS.md`): `bun run test` (Bun + Node split — "Fixture tests (jsdom-based HTML detection) run via `node --test` because bun is too slow with jsdom"), `bun run test:live-e2e` (opt-in, ~2 min, 19 framework fixtures, real `npm install` per fixture, boots Vite/Next/SvelteKit/Astro/Nuxt dev servers, Playwright Chromium), `bun run test:skill-behavior` (opt-in, LLM-backed).

`.github/workflows/ci.yml`:
- `changes` job runs `scripts/ci-test-plan.mjs` and emits 9 boolean outputs (`core`, `detector`, `live`, `framework`, `cli_remote_e2e`, `live_e2e`, `live_e2e_accept_cleanup`, `skill_behavior`, `live_svelte_adapter_deepseek`) — **path-based selective test planning**.
- `test-matrix` on Node `22.12.0` and `24`: core tests always; detector/live/framework tests conditionally; then `bun run build`, `bun run build:extension`, `npx web-ext@8 lint` for the Firefox package.
- **Generated-output drift gate** (the key one): `git diff --exit-code -- .agents .claude .cursor .gemini .github/skills plugin cli/engine/detect-antipatterns-browser.js extension/detector` — CI fails if the committed generated harness trees don't match a fresh build.
- `live-e2e-smoke` (PR/push, 2 groups, 15 min cap) vs `live-e2e-full` (`workflow_dispatch` only, 5 groups covering 19 fixtures, 25 min cap), with npm + Playwright caches and failure-artifact upload.
- Provider-key-gated jobs skip cleanly with an explicit echo when `ANTHROPIC_API_KEY`/`DEEPSEEK_API_KEY` are absent.
- `skill-behavior` runs only on non-PR events (secrets), against `claude-sonnet-5`, `gpt-5.6-luna`, `gemini-3.5-flash`, `deepseek-v4-flash`.

`.github/workflows/sync-generated-output.yml`: on push to `main` touching `.claude-plugin/**`, `cli/engine/**`, `skill/**`, `scripts/**`, `package.json`, `bun.lock` → `bun run build:release`, detect drift, commit generated provider output **directly to main**, with a 5-attempt rebase-rebuild-retry loop and backoff ("A push can lose the race against a human commit… issue #388: ~10% of runs. Every attempt builds from a fresh origin/main, so pushed output always matches the source state it lands on.").

`.github/workflows/sheriff.yml`: daily cron `17 15 * * *` running `scripts/github/sheriff.mjs --apply --warning-days 7 --close-days 14 --maintainers pbakaus --regular-contributors pbakaus,abdulwahabone`.

**The skill-behavior suite is the single most interesting test idea.** From `CLAUDE.md`:

> "It inlines the source `skill/SKILL.src.md` into the system prompt of a real LLM, gives the agent `bash` / `read` / `write` / `list` tools scoped to a temp workspace, and **asserts on the tool-call trace — not on the model's free-form output. The trace is the source of truth.**"

> "**Every provider, every run.** … **Don't substitute Claude alone**: many of the most useful findings come from divergence between providers."

Cost is documented and budgeted: "~5 min, ~$0.50-1.50 across providers… **Keep it out of CI unless you really want it there.**" Missing keys skip cleanly. The harness **symlinks source, not built output**, deliberately: "reference files surface their raw `{{placeholders}}`, but the assertions key on tool calls rather than content, so it doesn't matter for correctness." The scenario baseline table lives in `tests/skill-behavior/README.md`, not in CLAUDE.md — "Duplicating it in this file is how it went stale before." `workflow-contract.test.mjs` asserts on **question order and artifact writes** for four flows: attended fresh init, initialized natural build request, replacement-world redesign, scope-preserving refinement.

Also verified: an eval framework in a separate private repo (`~/code/impeccable-evals`) that "measures whether the `/impeccable` skill improves or harms AI-generated frontend design by running the same brief through a model **with and without** the skill loaded."

---

## (f) AGENT-FACING OPERATING RULES worth quoting (from CLAUDE.md / AGENTS.md)

Repo-hygiene rules, several of which are protocol-shaped:

- "**Do not add standalone skills** unless there's a strong reason. The consolidation was deliberate: the `/` menu pollution problem is real and gets worse as users install more plugins."
- "**Do not reintroduce per-domain reference files.**" (v4 deleted `typography.md`, `color-and-contrast.md`, `spatial-design.md`, `motion-design.md`, `interaction-design.md`, `responsive-design.md`, `ux-writing.md`, `cognitive-load.md`, `personas.md`, `heuristics-scoring.md`, `build-floor.md`, `live-generation.md` — "Their content lives in the command references and `craft-floor.md`, where it is **loaded only when it applies**.")
- "**a11y lives in `audit.md`**, not in SKILL.md or the mode guidance. **Models over-cautious themselves into safe, underdesigned output when reminded about accessibility at design time.**"
- "`doctor` **is a utility command, not a design command**… **Keep maintenance tooling out of the design menu.**"
- "**Do not bump manifest versions or add changelog entries in a feature PR.** … a version in a feature branch conflicts with every other open branch, and a changelog entry describes a release that has not happened yet."
- "Do not delete them as 'dead code' — **I made that mistake once and broke 8 tests.**" (named exports kept solely for test spying)
- Sandbox honesty rules for Codex agents: `bun run build:release` can hit `EFAULT` in the sandbox — "**Rerun the release build outside the sandbox before treating it as a real build failure.**" Puppeteer tests "can hang in the sandbox… Run them outside the sandbox for authoritative results." "A direct `bun test tests/detect-antipatterns-fixtures.test.mjs` can time out and **is not the supported signal**."
- **AI-agent contribution policy** (identical in both files): "AI agents must disclose AI assistance in commits, PR descriptions, comments, and issue text. If an AI agent is not explicitly operating under instructions from `pbakaus` or `abdulwahabone`, it must not create GitHub issues or PRs… Instead, add a file named `AI_PR_NOTICE.txt` to the diff with exactly this text: *This contribution was prepared by an AI agent that tried to ship unchecked vibes across a human boundary. Impeccable asks for an issue and maintainer approval first.*"

---

## (g) SINGLE-AGENT ASSUMPTIONS that would NOT transfer to a multi-agent protocol

1. **One parent thread owns the whole lifecycle.** The finish protocol assumes a single builder that *spawns* subordinate reviewers/documenters and *applies* their fixes. Parley Deck has N peer authors with no parent. "The parent agent applies your fixes" and "This review never runs inside the build thread" both need re-anchoring: in Parley the reviewer is a peer round, not a child process.
2. **Subagent isolation is achieved by `Task`-tool spawning.** critique.md's `Hard Invariants` define availability as "a sub-agent/Task tool is exposed in this session". Parley already gets isolation for free (separate CLI processes writing separate files) — the *mechanism* is wrong even though the *invariant* is right.
3. **Interactive user in the loop, synchronously.** `{{ask_instruction}}`, `AskUserQuestion`, "STOP and call the … tool", "Ask one round of two or three related questions", the browser-served decision page (`serve-question.mjs --start --wait --key`), the standing exit, re-roll — all assume a human answering mid-run. Parley runs are typically headless/asynchronous.
4. **Browser + screenshots are the primary evidence channel.** Two batched screenshot rounds, live overlay injection, Playwright/Puppeteer, `.impeccable/*.png`. Headless CLI agents in Parley may have no browser at all; the "degraded" language exists but is treated as second-class.
5. **A single deterministic ranking is the failure mode being fixed.** "your top-ranked structure is what every run would ship, and **a single ranking is deterministic, so the dice come from outside**." In Parley, five *different* models each produce their own ranking — external dice are less necessary; **cross-model divergence is the natural randomizer**. Copying `concept-seed.mjs` verbatim would import a solution to a problem Parley doesn't have.
6. **Session-scoped, run-once bootstrap.** "Run `node scripts/context.mjs` **once per session** … follow its directives and **do not rerun it**"; weekly-throttled staleness caches in `~/.impeccable/`. Parley has many concurrent agent sessions per idea; a once-per-session cache keyed to `~` collides.
7. **Command-menu UX as the entry point.** 23 slash commands, `pin`/`unpin` shims, "Never auto-run a command", context-aware menu from `context-signals.mjs`. Parley's entry point is an idea slug and a protocol phase, not a `/` menu.
8. **Single mutable working tree.** Hooks that edit-gate a shared repo (Cursor `preToolUse` blocking writes) presume one writer. With N agents writing concurrently, a blocking pre-edit hook per-agent needs worktree isolation (which parley-worktrees already provides).
9. **Provider-conditional prose compiled at build time.** `<codex>` / `<claude>` blocks are stripped per output tree. A Parley deck ships **one** COOPERATION.md read by all five agents simultaneously, so the same idea has to be expressed as *runtime-addressed* sections ("if you are the codex participant…"), not build-time compilation — unless parley-design also gains a build step.
10. **Maintainer-gated contribution model.** The `AI_PR_NOTICE.txt` rule and issue-first policy are about an OSS repo's boundary, not about how peer agents cooperate.

---

## Transferable to parley-design

Ranked by value-per-unit-of-effort for a vendor-neutral, multi-agent design skill.

1. **The named-rule format for design artifacts.** Every durable decision in DESIGN.md is a one-sentence, quotable, falsifiable "**The X Rule**" (*The OKLCH-Only Rule*, *The Weight-Inversion Rule*, *The Texture Budget Rule*, …). In a multi-agent protocol this is the difference between a reviewable artifact and mush: a reviewer can cite `The Two-Face Rule` and a signer can accept or contest it. Adopt as the required grain of parley-design's DESIGN.md and of any cross-review finding.
2. **The direction contract (THESIS / OWN-WORLD / STORY / FIRST VIEWPORT / FORM / FINISH), ≤150 words, embedded in the emitted artifact.** This is the perfect Parley consensus object: small enough for five agents to converge on, specific enough to be auditable, and it lives in the deliverable rather than a side file. Map FINISH to the Parley phase gate ("unreviewed and undocumented is unfinished"). Add a FORM-equivalent that records *which* dissenting direction was chosen and why.
3. **`match | acceptable adaptation | missing | contradicted | added without approval`** as the fixed classification vocabulary for design review, with the rule "**an adaptation counts as intentional only when it cites** the user answer, brief, accessibility need, or product truth that forced it; **an uncited deviation is a defect**." This is directly implementable as a Parley cross-review verdict schema, and it forces evidence citation, which is exactly what a signed consensus needs.
4. **Bounded verification with a hard ceiling and a scored verdict pass.** Two inspection rounds max, then hand off; two correction rounds max; the second verdict ends the work whatever it says; fixes scored `resolved | partial | unresolved`; "a fix answered mechanically, positions moved but the quality the finding named still absent, **is partial at best**"; "**no new hunt, no new checks**" in the verdict pass. Parley's fix-up cycles currently have no principled stop condition — this supplies one, plus a vocabulary for "we shipped with open items", plus the killer line: "**presenting mechanical confirmation as artistic success is how a failed build gets announced as a finished one.**"
5. **The calibration rule + "category-plus-avoidance" self-check.** "If someone could guess your aesthetic from the category alone, **or from category-plus-avoidance**, rework until neither answer is obvious." Cheap to state, hard to game, and the single best anti-slop test for a design consensus round. Pair with the three named saturated clusters and the 16-face denylist as concrete, checkable content.
6. **Per-model bias corrections in the roster.** Impeccable ships `<codex>`, `<claude>`, `<gemini>` blocks that name each model's *specific* rendition prior (Claude's "cream / serif / lamplight" default; GPT's ghost card, sketchy SVG, stripe backgrounds; Gemini's image-hover). Parley has an explicit roster (claude, codex, hermes, agy, kimi) — a per-participant "your known prior, treat it as already spent" section is high-leverage and has no analogue in Parley today. Also lift the 5 model-attributed detector rules.
7. **Two-channel assessment with mandatory anti-anchoring and a degradation banner.** Independent design judgment (A) and deterministic evidence (B) must be produced in isolation, A before B enters synthesis ("**Detector output is deterministic, but it still anchors judgment**"); the synthesis names agreements, detector-only catches, and false positives; any degradation prints `⚠️ DEGRADED: <reason>` as the first line and "**a silent degraded critique is a failed critique**". Parley should require the same banner whenever a participant is excluded, times out, or runs without a capability.
8. **The finish-reviewer's anti-anchoring rule.** "Inventory the comp's salient elements **in your own words before** reading the direction contract or any builder-authored summary… a review anchored on it inherits whatever that abstraction dropped." Directly applicable to Parley cross-review: read the artifact before reading the author's FINAL.md summary.
9. **Budgeted review with declared gaps.** "Treat reading as an allowance, not a prerequisite… **a review built from what you saw beats a perfect review that never arrives**… Name whatever went unread in the line above the sections." Solves the headless-agent timeout problem honestly instead of silently truncating.
10. **Calibrated scoring with anti-inflation anchors.** "**Most real interfaces score 20-32 out of 40**"; renormalize when a heuristic is `n/a` and "**Never print `/40` over a partial set**"; P0–P3 with the "would a user contact support?" tie-breaker; audit's 5×0–4 with rating bands. Gives Parley participants a shared, comparable scale so five reviews can be merged rather than averaged into mush.
11. **Truth-binds-claims.** "in greenfield work, author whatever illustrative material the concept needs at full fidelity, **label it synthetic**… What stays uninventable are commercial and factual claims: prices, customers, benchmarks, endpoints, capabilities the product does not have. **Refusing a bold direction because its demonstration data does not exist yet is the timidity reflex wearing honesty's clothes.**" A crisp resolution of the honesty-vs-boldness tension that multi-agent consensus otherwise resolves toward timidity.
12. **A deterministic detector as a peer reviewer.** 60 rules, LLM-free, exit `0`/`2`, JSON output, `slop` vs `quality` split, rules pointing back at the prose guideline they enforce (`skillSection`/`skillGuideline`), and the four **`design-system-*` rules that make DESIGN.md itself enforceable**. Even a 15-rule v1 gives Parley a non-negotiable, model-independent vote in the design consensus. Keep the reason-carrying waiver system (repo config + per-dev + inline `<!-- impeccable-disable rule: reason -->`).
13. **`docs/STYLE.md` as a separate, enforced anti-AI-prose brief.** The bar ("point to the sentence that makes it specifically yours"), the denylist with *Why* + *Use instead* per row, the two-tier gate (strict for user-facing docs, narrower for LLM-facing instructions), and the judgment-only list (negation pivot, triadic everything, uniform rhythm, hollow confidence, interchangeable copy). Parley artifacts are markdown written by five LLMs; this is the closest thing to a slop test for the *documents themselves*.
14. **The seven-candidate derivation discipline with the "rut" exclusion and the three-material-family floor.** "Name both [the category default and its predictable opposite] as the rut and keep them out of the seven-candidate list"; "when more than three of the seven share one material family, the derivation **stopped at the subject's most obvious artifact**". Adapts cleanly to Parley: each participant proposes candidates, the union is checked for material-family diversity, and the rut is a shared exclusion list.
15. **Degraded-mode single-sourcing.** One role definition, emitted as (a) a native subagent and (b) a `degraded/<role>.md` with a generated preamble: "you are both parties: produce the full output contract first, then act on it yourself" + "**disclose the substitution in one line**". Perfect for Parley when a participant is missing: any agent can take an absent role inline, provided it declares it.
16. **Capability self-declaration in the role file.** "You have no browser. Never attempt to render, screenshot, start a server, or open a page." + "When an expected input is missing, say so in one line at the top." Parley's roster has wildly different capabilities; making each role state its own limits prevents fabricated verification.
17. **Never repair drift as a side effect.** `rule:skill-drift-not-a-side-quest` — findings are *reported*, not acted on, unless asked; only `auto`-severity findings ride along with a write already happening. Prevents scope-creep collisions when five agents share a tree.
18. **Severity that says what should happen, not how bad it is.** `auto | mention | route` for staleness findings, with `route` naming the command that owns the repair. A better shape than a raw severity number for machine-consumed findings.
19. **Layered, reason-carrying waivers.** Every ignore carries a `--reason`; inline waivers travel with the file; `--no-config` gives a raw scan. Multi-agent review needs an audit trail for "we decided this is fine".
20. **Two-tier performance contract on any boot-time check.** Tier 1 (`collectBootFindings`) may only spend what a boot already spends — "**No directory walks, no git, no cross-workspace sweep**"; everything expensive goes to an on-demand `doctor`. Worth copying verbatim as a rule for parley-design's session bootstrap.
21. **Trace-based skill testing across providers.** Inline the skill source into a real model's system prompt, give it scoped tools, and **assert on the tool-call trace, not the prose**; run all four/five providers every time because "**many of the most useful findings come from divergence between providers**". Parley is already multi-provider; this is a natural fit and would catch protocol text that only Claude actually follows.
22. **Named artifact fields for anti-references and evidence.** PRODUCT.md's `## Anti-references` and `## Evidence on Hand` ("State absences that future work must not fabricate"). Simple, high-signal additions to whatever product-truth artifact parley-design defines.
23. **Deprecation with a stated reason.** "told only that a field is deprecated, models preserve it 'just in case', **which is how a retired axis keeps steering current output**." Applies to every schema change in COOPERATION.md.
24. **Provenance stamps that are schema versions, not release versions.** `<!-- impeccable:product-schema N -->`; "a PRODUCT.md written by v4.0.0 is not stale under v4.0.1."
25. **Rule-ID markers in the source** (`<!-- rule:skill-ban-gradient-text -->`), stripped at staging, so external eval tooling can pin each instruction to a stable ID. Parley could keep them *unstripped* so cross-review findings cite rule IDs.
26. **Deliberate load-ordering to avoid over-cautious output.** "a11y lives in `audit.md`, not in SKILL.md… **Models over-cautious themselves into safe, underdesigned output when reminded about accessibility at design time.**" A real, non-obvious prompt-engineering finding: put the safety checklist in the review phase, not the design phase.
27. **Count/version drift validators.** A build gate that fails when a number stated in N places disagrees. Parley already has an embedded-default drift guard; this generalizes it to counts and versions.
28. **Multi-runtime portability mechanics** (only if parley-design ever needs a build step): the provider table + `<tag>` conditional blocks (unknown tags preserved) + six placeholders + frontmatter-field gating + a `HARNESSES.md` support matrix + a CI `git diff --exit-code` gate on generated trees. The `{{ask_instruction}}` placeholder in particular is the cleanest solution seen for "how do I say *ask the user* portably".

---

## Do NOT copy

1. **The `concept-seed.mjs` external dice.** Its whole justification — "a single ranking is deterministic, so the dice come from outside" — dissolves when five different models each rank independently. Parley's diversity is endogenous. Importing a roll would add a service dependency (`IMPECCABLE_CATALOG_DIR` → roll API → degraded seed), telemetry, and a "paid-service moat" catalog that is explicitly *not* in the OSS repo. Use cross-model divergence as the randomizer and spend the effort on merging divergent proposals instead.
2. **The parent/child spawn architecture as the isolation mechanism.** Do not encode "spawn a sub-agent via the Task tool" as the definition of independence. Parley's participants are already separate processes; define independence as *separate artifacts written before reading peers*, and reuse only the invariant ("A must finish before B enters synthesis"), not the mechanism.
3. **The 23-command menu.** `/impeccable bolder|quieter|distill|delight|overdrive|colorize|typeset|...` is a single-agent chat UX. Parley's surface is phases and artifacts. Copying the menu imports exactly the "`/` menu pollution problem" Impeccable's own CLAUDE.md warns about, one level up. Take the *vocabulary* (bolder/quieter/distill as review dispositions or steering verbs) without the command surface.
4. **The synchronous human decision page** (`serve-question.mjs --start`, browser-served option cards, shimmer-waiting sketch slots, exit codes 2/3/4). Beautiful, and wrong for headless multi-agent runs. Keep the *counterweight rules* it encodes (the standing exit is the user's door; taste is never grounds for re-roll; re-roll eliminates everything shown) and express them as artifact fields a human can act on asynchronously.
5. **Mandatory screenshot/browser evidence as the definition of verification.** Parley's roster includes agents with no browser. Make visual evidence a *declared capability*, not a gate; otherwise agents will claim verification they did not perform — the exact failure Impeccable itself guards against with "Do not claim a user-visible overlay exists unless script injection succeeded."
6. **Build-time provider-block compilation for the protocol text.** One COOPERATION.md is read by all participants at once, so `<claude>…</claude>` stripping cannot happen. Re-express per-model guidance as runtime-addressed sections. (Impeccable's *content* — the model priors — transfers; its *delivery mechanism* does not.)
7. **The Neo Kinpaku design system itself** (kinpaku gold, verdigris patina, lacquer black, Alumni Sans at weight 100, the `.ks-*` kit). It is one project's committed world. Shipping it as parley-design's default would be precisely the monoculture the skill exists to fight — and it would violate Impeccable's own *The brief wins* rule. Copy the *shape* of DESIGN.md, never its values.
8. **The a11y-only-in-audit placement, unexamined.** The rationale ("models over-cautious themselves into safe, underdesigned output") is a real finding, but Impeccable can defer a11y because a single agent reliably reaches its own audit phase. In a multi-agent protocol where a design round may be signed and shipped without every phase running, deferring accessibility risks it never being checked. Adopt the *insight* (don't front-load the safety checklist into the generative prompt) but bind it to a hard gate (no consensus signature without the audit round).
9. **`.impeccable/`-style ephemeral runtime state.** Session files, PNG caches, `hook.pending.json`, `manual-edit-events.jsonl`, `~/.impeccable/staleness-check.json`. Parley already has `parley-deck/ideas/<slug>/` as canonical state; a second, partly-gitignored, per-developer state tree would fight it and would collide across concurrent agents.
10. **Weekly-throttled, home-directory-cached notices.** `~/.impeccable/staleness-check.json` with a once-a-week-per-project throttle assumes one human at one machine in one session. N concurrent agent processes will race it, and a throttled notice that one agent consumed is invisible to the other four.
11. **The `AI_PR_NOTICE.txt` policy.** Repo-boundary governance for an OSS project, unrelated to how peer agents cooperate. (The narrower rule — *disclose AI assistance / disclose a substituted role in one line* — is worth keeping; the shaming file is not.)
12. **`doctor` / `pin` / `hooks` as design commands.** Impeccable itself insists these stay out of the design menu; do not reintroduce maintenance tooling into parley-design's surface.
13. **Prose denylists applied to code and to LLM-facing text at full strength.** Impeccable deliberately runs a *narrower* list on `skill/**` and none on `skill/scripts/**`, because `robust`/`seamless`/`elevate` have legitimate technical readings. A single global denylist over a Parley deck would fire constantly on protocol prose. Copy the two-tier design, not one blunt list.
14. **The em-dash ban as an absolute, without its rationale.** The rationale ("the writer didn't pick a relationship between the clauses") is the valuable part; a bare character ban invites the `--` substitute Impeccable had to ban next, and produces cargo-cult compliance rather than better sentences.
