# Impeccable (pbakaus/impeccable) — command surface & routing

Study digest for the design of `parley-design`. Source root:
`/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/research/impeccable`
Skill version studied: `.agents/skills/impeccable/SKILL.md`, `version: 4.0.3`.

Files read in full: `SKILL.md` (79 lines), `reference/routing.md` (18), `scripts/command-metadata.json` (94),
and all 34 `reference/*.md` + `reference/degraded/finish-reviewer.md`. Total reference corpus ≈ 5,242 lines.
Machinery inspected: `scripts/` (41 entries), `scripts/detector/`, `scripts/lib/`, `scripts/concept-seed.mjs` header.

---

## 0. The one-paragraph shape of the thing

Impeccable is a **single skill with one entry file (79 lines) and 34 lazily-loaded reference playbooks**,
plus ~15 deterministic Node scripts and a rule-based HTML/CSS detector with ~63 named rules. SKILL.md is
deliberately tiny: it holds core principles, 4 "modes", one Commands table, and 8 lines of routing.
Everything else is `load exactly one playbook that owns the request`. Quality is enforced by three
separate mechanisms that do not overlap: (a) `craft-floor.md`, the taste/reflex floor loaded immediately
before any UI edit; (b) `detect.mjs`, a deterministic rule engine wired into editor hooks; (c) a
**subagent finish reviewer with a 5-section output contract and a max-2-correction-round ceiling**.

---

## 1. COMPLETE command table

`Category` values are Impeccable's own (Build / Evaluate / Refine / Enhance / Fix / Iterate). "When it runs"
is the lifecycle position. Verbatim purpose strings are from `scripts/command-metadata.json` where quoted.

| Command | Cat | One-line purpose | Inputs / preconditions | Outputs / artifacts | When it runs |
|---|---|---|---|---|---|
| `init` (alias `teach`) | Build | "Sets up a project for impeccable… writes PRODUCT.md (strategic: users, brand, principles)" | none; runs when `NO_PRODUCT_MD`. Explicitly forbidden to ask about aesthetics | `PROJECT_ROOT/PRODUCT.md` with verbatim `<!-- impeccable:product-schema 1 -->` comment; optionally `.impeccable/live/config.json` | **Entry point / gate.** Once per project. Has an explicit **Completion gate**: "verify that PRODUCT.md exists at the resolved path… If the file is absent, init is incomplete. Do not substitute interview notes, a planning packet, or later design prose for the file." |
| `shape [feature]` | Build | "Plan UX and UI before code. Runs a required multi-round discovery interview… produces a user-confirmed design brief" | PRODUCT.md; a named feature/surface | A 7-part brief (job/audience, outcome/proof, selected direction, scope/boundaries, states/ranges, interaction/layout, constraints/open decisions). **No code, no direction contract** | Entry point for ambiguous/multi-screen work. "shape never writes code or a direction contract." |
| `new-work.md` (not a command; the flow behind build/redesign) | Build | Make a new surface or replace a visual identity | PRODUCT.md required; DESIGN.md optional | Direction contract comment in the artifact; surface brief; built surface; DESIGN.md written **at finish** by the documenter | The real build spine. 7 numbered steps (see §3). |
| `craft [feature]` | Build | "Deprecated compatibility alias… It adds no behavior" | — | — | Never; routes to init → new-work. "Do not tell users they need to invoke `craft`." |
| `document` | Build | "Generate a DESIGN.md file that captures the current visual design system… Follows the Google Stitch DESIGN.md format" | code with tokens/components (Scan mode) **or** PRODUCT.md only (Seed mode) | root `DESIGN.md` (YAML frontmatter + 8 canonical sections) + `.impeccable/design.json` sidecar `schemaVersion: 2` | After a build (documenter subagent), or standalone to capture an incumbent system. Must not silently overwrite. |
| `extract [target]` | Build | "Pull reusable patterns, components, and design tokens into the design system" | an existing design system or user-chosen location | new shared components/tokens + migrated call sites + docs | Refinement/consolidation, after drift accumulates. |
| `critique [target]` | Evaluate | "Evaluate design from a UX perspective… quantitative scoring, persona-based testing, automated anti-pattern detection" | one resolvable file path or URL; `.impeccable/critique/ignore.md` if present | chat report (primary) + persisted snapshot in `.impeccable/critique/<slug>` with frontmatter `total_score`, `max_score`, `na_heuristics`, `p0_count`, `p1_count` | **Gate / diagnostic.** See §5. |
| `audit [target]` | Evaluate | "Run technical quality checks across accessibility, performance, theming, responsive design, and anti-patterns… P0-P3" | web project (native → `audit.native.md`) | Audit Health Score `??/20` table over 5 dimensions + severity-tagged findings + Recommended Actions | **Gate.** "Don't fix issues; document them for other commands to address." |
| `polish [target]` | Refine | "Performs a final quality pass fixing alignment, spacing, consistency, and micro-detail issues before shipping" | DESIGN.md/tokens; optionally the latest critique snapshot via `critique-storage.mjs latest` | edited source + final source diff cleanup | **Terminal refinement.** Every other Refine/Enhance/Fix command ends with "hand off to `$impeccable polish` for the final pass." |
| `bolder [target]` | Refine | "Amplify safe or boring designs" | **must know which section is the target and what stays untouched** | edits to the named target only | Taste op; see §4. |
| `quieter [target]` | Refine | "Tones down visually aggressive or overstimulating designs" | visitor mode; what's working | edits across colour/weight/simplification/motion/composition | Taste op; see §4. |
| `distill [target]` | Refine | "Strip designs to their essence by removing unnecessary complexity" | the ONE primary user goal | removals across IA/visual/layout/interaction/content/code + a written record of removed complexity | Refinement. |
| `harden [target]` | Refine | "Make interfaces production-ready: error handling, i18n, text overflow, edge case management" | — | overflow/i18n/error/empty/permission/concurrency handling | Pre-ship refinement. |
| `onboard [target]` | Refine | "Design onboarding flows, first-run experiences, and empty states" | the "aha moment"; users' experience level | welcome/setup/empty-state/tooltip/tour implementations | Feature-level work. |
| `animate [target]` | Enhance | "enhance it with purposeful animations, micro-interactions, and motion" | performance constraints | a written **motion thesis** (focal moment / continuity / feedback / budget) then implementation | Enhancement. |
| `colorize [target]` | Enhance | "Add strategic color to features that are too monochromatic" | existing brand colors; DESIGN.md | colour **roles** (canvas, text, action, borders, semantic, data) + tokens | Enhancement. |
| `typeset [target]` | Enhance | "Improves typography by fixing font choices, hierarchy, sizing, weight, and readability" | established families; `detect.mjs --scope type` | a stated type system (roles, contrast, measure, authoritative faces) + edits | Enhancement. Runs **two isolated assessments**. |
| `layout [target]` | Enhance | "Improve layout, spacing, and visual rhythm" | `detect.mjs --scope layout` | a stated **spatial thesis** + structural edits | Enhancement. Runs **two isolated assessments**. |
| `delight [target]` | Enhance | "Add moments of joy, personality, and unexpected touches" | the brand's emotional range | one **delight thesis** sentence + the smallest system that delivers it | Taste op; see §4. |
| `overdrive [target]` | Enhance | "Pushes interfaces past conventional limits… shaders, spring physics, scroll-driven reveals, 60fps" | **must propose 2-3 directions and get the user's pick before any code** | technically ambitious implementation w/ progressive enhancement | Late-stage, opt-in. Prints a literal banner: `──────────── ⚡ OVERDRIVE ─────────────` / `》》》 Entering overdrive mode...` |
| `clarify [target]` | Fix | "Improve unclear UX copy, error messages, microcopy, labels" | audience knowledge and emotional state | rewritten copy + terminology glossary when inconsistency spans the product | Fix. |
| `adapt [target] [context]` | Fix | "Adapt designs to work across different screen sizes, devices, contexts, or platforms" | target platforms/devices; web-only (native → `adapt.native.md`) | breakpoints/fluid layouts/touch targets/print/email variants | Fix. |
| `optimize [target]` | Fix | "Diagnoses and fixes UI performance across loading speed, rendering, animations, images, and bundle size" | measured baseline (LCP/INP/CLS) | perf edits + before/after metrics | Fix. |
| `live` | Iterate | "Interactive live variant mode. Select elements in the browser, pick a design action, and get AI-generated HTML+CSS variants hot-swapped via HMR" | running dev server w/ HMR **or** static HTML; `.impeccable/live/config.json` | 3 variants per Go, each with 0–4 declared parameters; on accept → **carbonize** into real source | Iteration loop; web-only. |
| `hooks <on\|off\|status\|ignore-rule\|ignore-file\|ignore-value\|reset>` | meta | Manage the design-detector hook per project | — | writes `.impeccable/config.json` (shared) / `.impeccable/config.local.json` (private) via `hook-admin.mjs` | Maintenance. |
| `doctor` | meta | "Report and repair drift between this project's Impeccable artifacts and what this version reads" | — | findings list `{id, artifact, path, severity, summary, fix}`; `--fix` applies `auto` findings | Maintenance. |
| — `ios.md` / `android.md` | platform | Native platform doctrine (HIG / Material 3) | `## Platform` is `ios`/`android`/`adaptive` | — | Loaded by `init` right after it records the platform. |
| — `operate.md` | mode depth | Extended Operate/Read doctrine | — | — | Loaded for Operate/Read surfaces. |
| — `visualize.md` | build sub-step | Direction comps & asset production | image generation available; PRODUCT.md + DESIGN.md | 3 comps under `.impeccable/mocks/`; written **implementation-fidelity inventory** | Between direction lock and build. "This step is proven to produce the most compositional and ambitious work." |
| — `craft-floor.md` | quality floor | Verify/Refuse lists | direction settled | — | "load immediately before editing UI… Do not load it for planning-only work." |
| — `new-work.md`, `routing.md`, `degraded/*.md` | infra | see §2, §3 | — | — | — |

Native variants: `audit.native.md` (5 dims: Accessibility, Performance, Appearance & Theming,
**Platform Conformance (CRITICAL)**, Adaptivity — same `??/20` skeleton, explicitly "keep the two in sync")
and `adapt.native.md` (phone→tablet, orientation/foldables, iOS↔Android idiom translation table, web→native).

---

## 2. The routing logic, exactly

Three-branch dispatch in SKILL.md lines 65–72:

- **No argument** → read `reference/routing.md` and present a **context-aware menu**; "never auto-run a command."
- **Explicit or clearly implied command** → "load its reference (native variant on native platforms) and follow it.
  **Ask once if two commands fit.**"
- **Otherwise** → general design work. "Missing PRODUCT.md routes a new surface or replacement world through
  init, then new-work; a narrow refinement of existing code proceeds on the incumbent implementation… offering
  init afterward rather than blocking on it."
- Aliases: `teach` → `init`. `craft` → deprecated new-work alias. `shape` owns task discovery, then enters
  new-work only for visual-world/surface-concept decisions.

**Setup is mandatory and runs once per session** (SKILL.md §Setup):
1. `node .agents/skills/impeccable/scripts/context.mjs` (add `--target <path>`) — loads PRODUCT.md, DESIGN.md,
   the matching surface brief, native guidance. "follow its directives and **do not rerun it**."
2. Load **exactly one** playbook that owns the request.
3. Load `craft-floor.md` immediately before editing UI (never for planning-only work).

**The no-argument menu is data-driven**, not static (`routing.md`):
- `NO_PRODUCT_MD` → lead with `$impeccable init` (one line on why), still show the rest; "don't silently jump into init."
- Otherwise run `node scripts/context-signals.mjs` once, read its JSON, lead with the **2-3 highest-value next
  commands** each with a one-line reason, then the full menu grouped by category.
- Explicit signal→command mappings (verbatim rules):
  - `setup.hasDesign` false while `setup.hasCode` true → `document`
  - `critique.latest === null` → offer `$impeccable critique <surface>` as "a strong default"
  - `critique.latest` with low `score` or non-zero `p0`/`p1` → `polish` ("it reads that snapshot as its backlog"), or re-run `critique` if stale
  - `git.changedFiles` pointing at one surface → scope `audit`/`polish` to those files, naming them
  - `devServer.running` true → `live` available; if false, don't lead with it. "`live` and the bundled `detect.mjs` are **web-only**"; on `ios`/`android`/`adaptive` lead with neither
- If `scan.targets` non-empty and platform is not native → run `node scripts/detect.mjs --json <targets>` once.
  `scan.via` ∈ `git-changes` | `source-dir` | `html` | `root`. Fold hits into picks: "many quality / contrast hits
  → `audit` or `polish`; a specific slop family → the matching command (gradient text or eyebrows → `quieter` /
  `typeset`, flat or gray palette → `colorize`…)". If detect errors or the tree is large, **skip and never block**.
- "Keep it to 2-3 pointed picks with the exact command to type."

**Cross-command routing back-links.** Both `critique` and `audit` end with a `Recommended Actions` list whose
allowed vocabulary is a fixed, verbatim-repeated 19-command allowlist (`adapt, animate, audit, bolder, clarify,
colorize, critique, delight, distill, document, harden, layout, onboard, optimize, overdrive, polish, quieter,
shape, typeset`), with the rules "Map each Priority Issue to the appropriate command", "Skip commands that would
address zero issues", and **"End with `$impeccable polish` as the final step if any fixes were recommended."**

Drift rule (SKILL.md, last line): **"Never repair drift as a side effect of a design task."** A `CONTEXT_STALE`
finding is reported, not acted on — except findings marked `auto`.

---

## 3. The lifecycle these commands imply

```
                 ┌──────────────────────────────────────────────── maintenance ───┐
                 │  doctor  ·  hooks on/off/ignore-*                              │
                 └───────────────────────────────────────────────────────────────┘

 ENTRY          PLAN            BUILD (new-work 1-7)                REFINE/ENHANCE/FIX      GATE      SHIP
 init  ──────►  shape  ──────►  1 decide what's true                bolder  quieter        critique  (finish
 (PRODUCT.md)   (brief,         2 ask 1 round of 2-3 Qs             distill harden         audit      review
   │            no code)        3 choose amount of invention        onboard animate        (scored,   verdict
   │                              (extend | whole surface |         colorize typeset        no fixes)  + DESIGN.md)
   │                               create/replace world)            layout  delight        │
   │                              → concept-seed.mjs roll           overdrive clarify      │
   │                            4 commit the world (colour           adapt   optimize      │
   │                              strategy / faces)                  extract               │
   │                            5 record DIRECTION CONTRACT ─────────────┐                 │
   │                            6 build with full commitment            │                 │
   │                            7 inspect (≤2 rounds) → finish reviewer  │                 │
   │                              → documenter → DESIGN.md              ▼                 ▼
   └────────────────────────────────────────────────► live (iterate)  polish ────────► ship
```

**Entry points:** `init` (the only gate on product truth), `shape` (planning-only), `new-work` (implied by any
"build/redesign this" request), `document` (capture an incumbent system), `live` (iterate in-browser).

**Refinements (all terminate in polish):** `bolder`, `quieter`, `distill`, `harden`, `onboard`, `animate`,
`colorize`, `typeset`, `layout`, `delight`, `clarify`, `adapt`, `optimize`, `adapt.native`. Each file literally
ends "hand off to `$impeccable polish` for the final pass." `overdrive` and `extract` do not (overdrive ends on
its own 4 tests; extract ends on documentation).

**Gates (produce verdicts, never edits):**
- `audit` — "Don't fix issues; document them for other commands to address."
- `critique` — "The chat response is the primary deliverable; the snapshot is an archive/backlog."
- the **finish reviewer** subagent inside new-work §7 — "You do not edit anything; the parent agent applies your fixes."
- `init`'s **Completion gate** — PRODUCT.md must exist as a file.
- `visualize`'s **one approval point** — "Do not begin code until the user approves a direction or explicitly delegates."

**The build spine (new-work §1–§7), condensed with its hard numbers:**
1. Decide what is already true: redesign / established world / incomplete brand / no visual authority.
2. One round of 2-3 related questions, mode-specific. "Do not ask for CSS values or canned aesthetic lanes."
3. Amount of invention:
   - *Extend an existing surface* → inherit, no tournament, no DESIGN.md change.
   - *Whole surface in an established world* → derive **5–7 materially different structures**, then
     `concept-seed.mjs --scope surface --mode <mode>` assigns which gets built.
   - *Create/replace the world* → name the rut and its opposite and exclude both; list **7 concrete visual
     systems** from the audience's culture, spanning **≥3 material families** ("When more than three of the seven
     share one material family, the derivation stopped"); turn them into complete directions; roll
     `concept-seed.mjs --scope direction --mode <mode>`; present **one direction fully committed** plus **at most 3
     dealt challengers**, plus a permanent **"standing exit"** (the category standard played straight, never
     recommended by the agent), plus re-roll ("after two consecutive re-rolls, ask what quality is missing").
     Decision is served as a real HTML page via `serve-question.mjs --start --payload <file>` (exit 3 = keep
     waiting, exit 4 = closed without answer → re-ask once then proceed unattended, exit 2 = headless → use the
     structured question tool).
4. Commit the world: colour strategy is one of **Restrained / Committed (one saturated colour carries 30-60% of
   the surface) / Full palette (3-4 named roles) / Drenched**. Dark-vs-light is decided by "one sentence of
   physical scene (who uses this, where, under what light)."
5. **Record the direction contract** — an HTML comment at the top of the emitted markup, **5 blocks, ≤150 words**:
   `THESIS`, `OWN-WORLD`, `STORY`, `FIRST VIEWPORT`, `FORM` (+ the seed key), closed by a literal
   `FINISH:` line reading verbatim *"unreviewed and undocumented is unfinished; this build ends with the finish
   review, the verdict, and DESIGN.md"*. "If a block reads like a mood, the direction is not decided yet."
   DESIGN.md is written **at finish, from the built world** — "a rulebook written before the build gets defended
   against reality instead of describing it."
6. Build with full commitment ("a stock component inside a committed form is a lapse").
7. Inspect desktop+mobile in **one batched screenshot round**, fix, confirm with **one more round; two rounds is
   the ceiling**. Then hand to the finish reviewer (see §5), apply fixes in one batch, recapture, get a **verdict**
   scoring each fix `resolved | partial | unresolved`; partial/unresolved get **exactly one more batch**;
   **two correction rounds is the ceiling and the second verdict ends the work whatever it says.** Then the
   documenter writes DESIGN.md + sidecar.

---

## 4. Taste operations, and how they're made repeatable

The "taste" commands are `bolder`, `quieter`, `delight` (+ arguably `distill`, `overdrive`, `colorize`).
None of them is left as vibes. Five distinct de-vibing devices are used:

**(a) A declared missing-input line.** Several taste files open with a literal
`> **Additional context needed**: …` header naming exactly what must be known before acting:
- `bolder`: "which section is the target, and what must stay untouched."
- `delight`: "the brand's emotional range."
- `animate`: "performance constraints." `clarify`: "audience knowledge and emotional state."
- `onboard`: "the 'aha moment' you want users to reach, and users' experience level."
- `adapt` / `adapt.native`: "target platforms/devices and usage contexts."

**(b) Visitor-mode conditioning.** `quieter`, `delight`, `animate`, `colorize`, `typeset`, `layout` each open with a
`## Visitor mode` block that redefines the operation per mode. E.g. `quieter`:
> Persuade + Experience: "quieter" means more restrained palette, more whitespace, more typographic air. Drama is
> reduced, not eliminated; the POV stays intact.
> Operate + Read: "quieter" means reducing visual noise. Fewer background accents, flatter cards, less color,
> less motion. The tool should disappear more completely into the task.

**(c) Numeric dimensions instead of adjectives.** `quieter` specifies:
saturation "fully saturated → **70-85%** saturation"; weights "**900 → 600, 700 → 500**"; motion distances
"**10-20px instead of 40px**"; "accent as **10% rule**"; "**Never bounce or elastic**" / "use ease-out-quart";
"Never gray on color — use a darker shade of that color or transparency instead."
`animate` ships a duration table: 100–150 ms immediate feedback / 150–300 ms routine state change /
300–500 ms layout, overlay, view transition / 500–800 ms authored focal entrance, "exit faster than entrance",
`cubic-bezier(0.16, 1, 0.3, 1)`.

**(d) A written thesis before editing.** Each taste op forces one declarative sentence/plan first:
- `delight`: "**Define one delight thesis** — State in one sentence what the user should feel and why that
  feeling belongs to this product."
- `animate`: "**Set the motion thesis**" (focal moment / continuity / feedback / budget). "A generic
  fade-and-rise, hover lift, parallax layer, or scroll reveal is **not a thesis**."
- `layout`: "**Set the spatial thesis**." `typeset`: "**Set the system**." `colorize`: "**Choose a strategy**"
  (emotional temperature, dominant relationship, contrast range, dosage) then build **roles, not a bag of swatches**.

**(e) A falsifiable exit test.** Each taste op ends with a checkable verification list, several of which are
named tests:
- `bolder` — **the skeleton test**: "Strip the copy out of your planned section and study the bare structure.
  Does the skeleton still say what this section is and why it matters…? If it only works once the words return,
  the boldness is in the text size, not the design." Plus a 4-item "Before you finish" checklist whose first
  item is "Everything outside the named target is unchanged."
- `layout` / `critique` — **the squint test**: "With detail blurred, can you still identify the primary element,
  the secondary element, and the major groups in order?"
- `delight` — "The moment is specific enough that a neighboring product could not use it unchanged."
- `overdrive` — **the wow test / the removal test / the device test / the context test**.
- `new-work` — **the memory test**: "if someone left after one viewport, what would they describe an hour later?
  If the honest answer is a mood, the concept has not committed yet."

**(f) Scope sovereignty.** `bolder` has a section literally titled `## Scope is sovereign`:
"'Everything else stays' is a literal instruction. Touch only the named target. Do not restyle its neighbors…
do not add colors, fonts, radii, shadows, or system primitives the surface does not already own."
And the anti-reflex opener: "The reflex answer, reaching for more effects, is the opposite of bold; reject it first."
Plus the diagnostic: "A flat section is typically one that quietly opts out of the system's own strongest moves."

**(g) In `live` mode, taste ops become enumerated axes.** `live.md` §4 Phase D forces per-action variance:
`bolder` = "amplify a different dimension per variant (scale / saturation / structural change). Not three
'slightly bigger' variants." `quieter` = "pull back a different dimension (color / ornament / spacing)."
`delight` = "different flavor of personality (unexpected micro-interaction / typographic surprise / illustrated
accent / sonic-or-haptic moment / easter-egg interaction)." `colorize` = "different hue family each (not shades
of one hue)." `typeset` = "different type pairing AND different scale ratio each." Etc. for all 11 actions.
Live also gives each taste op a **MUST parameter**: `colorize` → `color-amount` (range 0–1, step 0.05, default 0.5,
authored as `var(--p-color-amount, 0.5)`); `typeset` → `scale` (0.85–1.3, step 0.05); `layout` → `density`
(0.6–1.4, step 0.05). Budget by element weight: leaf **0 params**, small composition **0–1**, medium **target 2**,
large **target 2–3, up to 4**; **hard cap 4**.

---

## 5. `craft-floor` and `critique` — the two quality-gate primitives

### 5.1 `craft-floor.md` (45 lines) — the taste floor, loaded just-in-time

Framing (line 3): *"Load this after the direction is settled, and build without announcing the checklist. A pinned
brief or the committed visual world overrides anything here; your own habit does not. When the design hook is
active it already enforces the mechanical checks below as you edit: act on its findings instead of re-auditing
each rule."*

Two lists.

**`## Verify` — 8 checks "on the built result, not an intention", run together in the batched inspection round:**

| Check | Verbatim threshold |
|---|---|
| Contrast | "body and placeholder text ≥4.5:1, large text ≥3:1. On colored surfaces tint secondary text from that hue or the foreground; **never gray**." |
| Depth | "shadows carry an offset and a soft blur. A zero-offset colored halo is decoration." |
| Spacing | "tight groups, generous separation, **more space above a heading than below it**. Read the computed values." |
| Type | "body measure **65–75ch**, display **max 6rem**, tracking floor **-0.04em**, balanced headings, obvious scale and weight steps. Run the real copy at every breakpoint." |
| Motion | "**one authored moment**, not scattered effects and not one identical entrance on every section. Exponential ease-out from an already-visible default. Reach past transform and opacity: blur, backdrop-filter, clip-path, mask, and shadow belong to the palette when they stay smooth." |
| States | "hover, disabled, loading, error, empty. Plus real content, working controls, responsive composition, keyboard focus." |
| Copy | "the product's own language. Controls name their action; errors name the problem and the recovery." |
| Coverage | "every brief requirement present and findable within seconds." |

**`## Refuse` — "These are the category's defaults, not bans: the brief's own words can earn any of them.
Reaching for one when the axis is free means you were not deciding."**

*Page scaffolds:* same-size icon+heading+text cards as page structure ("Cards are the lazy container; **nested
cards are always wrong**"); the hero-metric template (big number, small label, supporting stats, accent);
**a kicker or eyebrow above a heading — "This one is a ban, not a default: no brief earns it back"**; section
numbers (01/02/03) unless the sequence carries information; a modal for a task needing neither interruption
nor protected focus.

*Surface habits:* gradient text ("Emphasis comes from weight or size"); glass/blur as decoration; a coloured
`border-left`/`border-right` above 1px on cards/list items/callouts/alerts; sparklines, progress rings, and
soft-shadowed rounded rectangles standing in for content; monospace as a costume for "technical"; light-or-dark
picked by category rather than by use scene.

*Hard numbers:* "Tracking stops at **-0.04em**. **-0.02 to -0.03em** usually reads better." "Declare elevation
once, border **or** shadow. A 1px border under a wide soft shadow is **the ghost card**. Card radii stay at
**12–16px**; pills are for small controls." "Real illustration or none. Sketch-style SVG scenes,
`loose-sketch` / `doodle` class names, and `feTurbulence` grain read as amateur." "Backgrounds are surfaces…
`repeating-linear-gradient` stripes and two-axis grid overlays need an actual canvas, map, blueprint, or
measuring tool under them."

Closing line: *"The floor holds the mechanics; it never picks the direction. With every check green, spend the
page on the committed world, and when torn between refined and committed, **commit**."*

**Companion "calibration" doctrine lives in `new-work.md` §4** — the three AI-slop attractors named explicitly:
1. "warm cream ground, high-contrast serif display, and a terracotta or signal-red accent";
2. "near-black with one neon accent and glowing edges";
3. "broadsheet-editorial hairlines, italic display serif, and small tracked mono labels."
And the falsification test: *"if someone could guess your aesthetic from the category alone, or from
category-plus-avoidance, rework until neither answer is obvious."*
Plus a **named font blocklist** ("these training-data defaults mean you stopped looking"): Fraunces, Playfair
Display, Cormorant, Lora, Crimson, Newsreader, Syne, Space Grotesk, Space Mono, IBM Plex, Inter-as-display,
DM Sans, DM Serif, Outfit, Plus Jakarta Sans, Instrument Sans — "Naming one of these faces anyway requires a
reason no other face could satisfy, and **a subject association is never that reason**."

### 5.2 `critique.md` (812 lines) — a two-assessor, scored, persisted review

**Purpose (verbatim):** "Resolve one stable target, run two independent assessments, synthesize a design critique,
persist a snapshot, and ask the user what to improve next."

**`### Hard Invariants` (7, verbatim-tight):**
- Assessment A (design review) and Assessment B (detector/browser evidence) are **both required**.
- "A and B **MUST run as two isolated sub-agents** whenever a sub-agent/Task tool is exposed. Running them inline
  in this context is 'possible' but is **NOT permitted**; it is a degraded run."
- "If you degrade for any reason, the report's first line MUST be a banner: `⚠️ DEGRADED: single-context (<reason>)`.
  A silent degraded critique is a failed critique."
- "**Assessment A must finish before detector findings enter the parent synthesis context.** Detector output is
  deterministic, but it still anchors judgment."
- "A skipped detector is a failed critique run unless `detect.mjs` is missing or crashes after a real attempt."
- Viewable targets require browser inspection when available.
- Any local server started for visualization must run in background, have a recorded stop method, and be stopped
  before final reporting.

**Provenance header (mandatory first line):** `Method: dual-agent (A: <agent-id> · B: <agent-id>)` or the
degraded banner. "'Unavailable' means exactly one thing: no sub-agent/Task tool is exposed in this session
(or the user declined). **It does not mean inconvenient.**"

**Setup:** resolve the target to a concrete path or URL ("Prefer a source path over a dev-server URL… ports drift,
paths do not"), slug it with `critique-storage.mjs slug`, "never hand-write a slug", and read
`.impeccable/critique/ignore.md` — "it is the only prior-run input critique consumes."

**Assessment A (design review)** evaluates: **Design specificity** ("Is the composition, interaction, and visual
language grounded in this product, or could an unrelated product use it unchanged? **Make this judgment before
seeing detector output**"), holistic design, cognitive load, emotional journey (peak-end rule), and Nielsen's 10
heuristics scored 0–4.

**Assessment B (detector + browser evidence):** `node scripts/detect.mjs --json [target]`; "Pass markup
files/directories; do not pass CSS-only files"; "For very large trees (**500+ scannable files**), narrow scope or
ask"; "Exit code 0 = clean; 2 = findings." Browser overlay flow is a 5-step protocol including a **mutability
preflight** ("setting `document.title` and appending a `<script>` tag. Read-only evaluate APIs do not count") and
"Do not claim a user-visible overlay exists unless script injection succeeded and the detector ran in the page."
Multi-view targets: inject on **3-5 representative pages**.

**Synthesis:** "Do NOT simply concatenate. Weave the findings together, noting where the LLM review and detector
agree, where the detector caught issues the LLM missed, and **where detector findings are false positives**."

**Scoring machinery:**
- Nielsen 10 heuristics × 0–4, each with its own 5-row 0/1/2/3/4 rubric table (all ten spelled out in the file).
- "Be honest with scores. A 4 means genuinely excellent. **Most real interfaces score 20-32 out of 40.**"
- **Mode applicability:** heuristics 7 (Flexibility & Efficiency) and 10 (Help & Documentation) may be `n/a` on
  Persuade/Experience surfaces. "The applicable maximum is 4 times the number of heuristics you actually scored:
  **/40** when all ten apply, **/32** when two are `n/a`. **Never print `/40` over a partial set.**"
- Rating bands: 36–40 Excellent · 28–35 Good · 20–27 Acceptable · 12–19 Poor · 0–11 Critical; renormalised by
  percentage when `max_score < 40` (90%+/70%+/50%+/30%+).
- **Issue severity P0–P3:** P0 Blocking (prevents task completion; showstopper) · P1 Major (fix before release) ·
  P2 Minor (workaround exists) · P3 Polish. Tie-break heuristic: *"Would a user contact support about this?
  If yes, it's at least P1."*
- **Cognitive load:** 3 types (intrinsic/extraneous/germane), an **8-item checklist**, and scoring
  "0–1 failures = low (good). 2–3 = moderate. **4+ = high (critical fix needed)**", plus the working-memory rule
  "**Humans can hold ≤4 items in working memory** (Miller's Law revised by Cowan, 2001)" with bands
  ≤4 manageable / 5–7 pushing / 8+ overloaded, and 8 named violations (The Wall of Options, The Memory Bridge,
  The Hidden Navigation, The Jargon Barrier, The Visual Noise Floor, The Inconsistent Pattern, The Multi-Task
  Demand, The Context Switch).
- **Personas:** 5 fixed archetypes — Alex (Impatient Power User), Jordan (Confused First-Timer), Sam
  (Accessibility-Dependent), Riley (Deliberate Stress Tester), Casey (Distracted Mobile User) — each with
  Profile / Behaviors / Test Questions / Red Flags, plus a selection table by interface type (landing page →
  Jordan, Riley, Casey; dashboard → Alex, Sam; e-commerce → Casey, Riley, Jordan; onboarding → Jordan, Casey;
  data-heavy → Alex, Sam; form-heavy → Jordan, Sam, Casey). "Be specific… **write what broke for them**."

**Persistence & trend:** metadata is passed as env JSON
`IMPECCABLE_CRITIQUE_META='{"target":…,"total_score":n,"max_score":n,"na_heuristics":"…","p0_count":n,"p1_count":n}'`
then `critique-storage.mjs write <target> <body-file>`; then `critique-storage.mjs trend <target> 5` → last 5
frontmatter entries, printed as one line:
`> **Trend for <slug> (last 5 runs): 24 → 28 → 32 → 29 → 32 (out of 40)**`. When maxima differ, print each score
with its own denominator (`24/32 → 30/40`) and note the runs are **not like-for-like**. "This is fire-and-forget…
Failures here should not block the rest of the flow."

**Then:** 2–4 targeted questions ("Every question must reference specific findings from the report. **Never ask
generic 'who is your audience?' questions.**" "If findings are straightforward… skip questions"), then a
prioritized `Recommended Actions` list from the 19-command allowlist, ending in `polish`.

**Tone rules (verbatim):** "Be direct… Be specific. 'The submit button,' not 'some elements.'… Give concrete
suggestions. **Cut 'consider exploring…' entirely.** Prioritize ruthlessly. If everything is important, nothing
is. **Don't soften criticism.**"

### 5.3 The third gate — the finish reviewer (`degraded/finish-reviewer.md`)

Not a user-facing command, but the strictest gate in the system and the most transferable.
- Role: "fresh eyes on a done artifact, outside the build thread's attention gravity. **You do not edit anything**."
- "You have **no browser**. Never attempt to render, screenshot, start a server, or open a page… **screenshots you
  fail to pass are checks it cannot run**" (stated on the caller's side in new-work §7).
- **Turn budget as doctrine:** "You run under a hard turn ceiling… a review built from what you saw beats a
  perfect review that never arrives… by roughly the **tenth turn** stop reading and write. Name whatever went
  unread in the line above the sections."
- **Anchoring guard:** "inventory the comp's salient elements in your own words **before** reading the direction
  contract or any builder-authored summary: the contract is the builder's abstraction of the comp, and a review
  anchored on it inherits whatever that abstraction dropped."
- 5 ordered checks: **Persistence**, **Fidelity** (an element matrix classifying every salient element as
  `match | acceptable adaptation | missing | contradicted | added without approval`, with **two mandatory rows,
  TYPE and MATERIAL**; "an uncited deviation is a defect"), **Ceiling** (unused native devices of the world),
  **Contract promise-by-promise** (first verify FORM carries the seed key — "a contract with no seed key… means
  the roll was skipped and that is a material fix ahead of any craft point"), **Truth** (synthetic data labeled,
  no invented commercial claims, "an asset applied at near-zero opacity… is a compliance token, not a shipped
  material").
- **Output contract: exactly five sections** — `persistence`, `fidelity`, `ceiling`, `material_fixes`
  (ordered, **at most eight**, fidelity failures ahead of craft), `keep` (one line naming what must not be diluted
  while fixing). "**No praise, no summary prose.**"
- **Verdict pass:** scoring only, not re-hunting; each material fix → `resolved | partial | unresolved`
  ("a fix answered mechanically, positions moved but the quality the finding named still absent, is **partial at
  best**"); at most **three** regressions; returns exactly two sections `verdict` and `remaining`.
- Caller-side ceiling (new-work §7): "**two correction rounds is the ceiling, the second verdict ends the work
  whatever it says, and the reviewer's findings are the only list you work from, never your own re-opened hunt.**
  Report the final verdict table to the user as it stands, open items included: presenting mechanical confirmation
  as artistic success is how a failed build gets announced as a finished one."

---

## 6. (a) DOCTRINE — design rules, taste, knowledge

- **Core principles (SKILL.md):** "Go all out. No hedging, no shortcuts." / "Dream big and bold." /
  "**Verify in bounded passes, not a loop**… Build fully, inspect once with a batched round (desktop and mobile
  together), fix everything it shows in one batch, confirm with at most one more round, and stop polishing.
  Open-ended self-QA burns the user's money doing worse what the finish handoffs do better."
- **The brief wins.** "Honor pinned aesthetics, eras, materials, fonts, and palettes even when they conflict with
  a saturated-pattern warning. Redirecting a clear brief toward your taste is failure."
- **Refinement preserves; redesign replaces.** "Never split the difference into polish on the discarded look."
- **Visual authority is evidence, not a filename.** "Missing DESIGN.md alone does not make a project greenfield."
- **4 Modes** = what the visitor's success looks like: **Persuade / Operate / Read / Experience**. Chosen "from the
  requested surface, not the product" ("A tool's landing page is still Persuade; a fashion house's documentation
  is still Read") and persisted only in that surface brief. Modes condition ~8 of the playbooks.
- **The slop tests** — three parallel formulations: web/product (`operate.md`) "Product UI's failure mode isn't
  flatness, it's **strangeness without purpose**… The bar is **earned familiarity**. The tool should disappear into
  the task."; iOS "Would a fluent iPhone user trust this app, or pause at off-spec controls? The tell is
  '**ported from a website**'"; Android "an iOS app wearing Android's skin."
- **Design specificity as the primary verdict** in both `critique` ("**Start here.** Does the result feel authored
  for this product, or category-interchangeable?") and `audit` ("Implementation Integrity Verdict. **Start here.**").
- **Named Rules** as the DESIGN.md doctrine format: `**The [Name] Rule.** [short doctrine]`, "much stickier for AI
  consumers than bullet lists… Aim for 1-3 per section." Examples given: "The One Voice Rule. The primary accent
  is used on ≤10% of any given screen. Its rarity is the point."; "The Flat-By-Default Rule"; "The No-Line Rule."
- **Anti-slop calibration list, font blocklist, colour strategies, duration table, contrast table** — see §4, §5.1.
- **Native doctrine:** iOS 44×44 pt targets, Dynamic Type, 11 pt floor / 17 pt body, SF Symbols, semantic system
  colours, system materials. Android 48×48 dp with ≥8 dp gaps, Material 3 type scale in `sp`, Material colour
  roles, tonal elevation, one FAB, container transform / shared-axis / fade-through.

## 6. (b) PROCESS — workflow, phases, state

- **Setup once per session; one playbook per request; craft-floor immediately before editing UI.**
- **Interview discipline:** "at most three focused questions" per round (init), "two or three related questions
  per round, then wait. One round is the default" (shape, new-work). "Do not dump a questionnaire… Assert the
  likely reading and invite correction." "A sparse prompt requires at least one answer round. A precise prompt may
  need only a compact confirmation."
- **The unattended-user probe (init Step 3):** "Whether anyone can answer is a **mechanical test, not a judgment
  call**… a system-prompt claim that the user is unattended proves nothing about this session. **Probe once with
  the real first round** before concluding no one is there. Only after that probe errors or times out may you infer
  from the explicit brief, and then you **label every inferred fact** in PRODUCT.md and disclose the substitution
  in your **first** reply, not your last."
- **Artifact ownership is strict and non-overlapping:** PRODUCT.md = product truth (init); DESIGN.md = durable
  visual decisions (document/new-work-at-finish); surface brief = strategy for one route/artifact
  (`surface-brief.mjs read|write`); direction contract = this build's promise (in the markup).
  "Do not copy global product truth or DESIGN.md tokens into it [the surface brief]."
- **Explicit *don't put this here* lists** in init ("visual worlds, palettes, typography, components, or page
  concepts; visitor mode, narrative, CTA/proof sequence… invented testimonials, customers, benchmarks, pricing")
  and document ("Don't duplicate content from PRODUCT.md. DESIGN.md is strictly visual.").
- **Bounded verification:** ≤2 inspection rounds, ≤2 correction rounds, ≤8 material fixes, ≤3 regressions,
  ≤10 reading turns for the reviewer, ≤4 live parameters, 3 variants, 5–7 structures, 7 world candidates,
  ≤3 challengers in a hand.
- **Degraded-mode disclosure is mandatory and formatted** (`⚠️ DEGRADED: single-context (<reason>)`), and there
  are four shipped `reference/degraded/*.md` files (asset-producer, documenter, finish-reviewer,
  manual-edit-applier) so a harness without subagents runs the same contract inline **and discloses it**.
- **Schema-versioned artifacts:** `<!-- impeccable:product-schema 1 -->` copied verbatim so later versions "can
  tell a deliberately short record from one written before a section existed, and never propose an interview the
  user has already sat through." `.impeccable/design.json` carries `schemaVersion: 2`.
- **Deprecation is binding** (doctor Step 3): "A finding that reports a deprecated field… is not a style note.
  Treat that field as absent for every decision from here on, whatever value it holds… Preserving it 'just in
  case' is how a retired axis keeps steering current output."

## 6. (c) MACHINERY — scripts, detectors, enforcement

| Script | Role |
|---|---|
| `context.mjs` (59 KB) | Session boot: loads PRODUCT.md, DESIGN.md, surface brief, native refs; emits directives incl. `NO_PRODUCT_MD`, `CONTEXT_STALE`, `UPDATE_AVAILABLE`, `MANUAL_DETECTOR_REQUIRED`. Run once, never rerun. |
| `context-signals.mjs` | JSON for the no-arg menu: `setup.{hasDesign,hasCode,platform}`, `critique.latest`, `git.changedFiles`, `devServer.running`, `scan.{targets,via}`. |
| `detect.mjs` + `detector/` | Deterministic rule engine over HTML/CSS. Exit 0 clean / 2 findings. `--json`, `--scope type|layout`, `--no-config`, `--no-inline-ignores`. ~63 rule ids incl. `gradient-text`, `kicker-above-heading`, `hero-eyebrow-chip`, `nested-cards`, `cream-palette`, `ai-color-palette`, `gpt-thin-border-wide-shadow`, `codex-grid-background`, `radial-halo`, `dark-glow`, `overused-font`, `extreme-negative-tracking`, `low-contrast`, `gray-on-color`, `monotonous-spacing`, `flat-type-hierarchy`, `em-dash-overuse`, `aphoristic-cadence`, `marketing-buzzword`, `theater-slop-phrase`, `design-system-{color,font,font-size,radius}`, `text-overflow`, `clipped-overflow-container`, `broken-image`. |
| `hook.mjs` / `hook-lib.mjs` (83 KB) / `hook-before-edit.mjs` / `hook-admin.mjs` | Editor hooks. **Two tiers**: per-edit "immediate tier" = mechanical, unambiguous (broken images, overflow/clipping, contrast, gradient text, glow shadows, design-system drift); everything else (copy cadence, palette/typography taste, layout rhythm) deferred to a **deep pass on the `Stop` hook event**, deduplicated against what per-edit already said. Cursor uses `preToolUse` to **block bad proposed writes**; Claude Code / Codex / Copilot post-edit reminders only. |
| `critique-storage.mjs` | `slug` / `write` / `latest` / `trend <n>` over `.impeccable/critique/`. |
| `concept-seed.mjs` (28 KB) | External dice. Header states the measured problem: "Left alone, it then always builds its #1 — and a single model's resonance ranking is deterministic… **Measured: 30/35 identical concepts across 16 prompt framings; the model cannot roll its own dice.**" Provides ASSIGNED INDEX (which of the model's own ranked shortlist gets built), 6 CHALLENGERS from 3 tiers, `--reroll n` reproducible chains, and ratings-weighted draws. |
| `serve-question.mjs` (59 KB) | Serves the direction-decision page; exit codes 2 (headless → use structured tool) / 3 (keep waiting) / 4 (closed, no answer). |
| `surface-brief.mjs` | `read`/`write` per-target strategy briefs. |
| `doctor.mjs` | Drift report `{id, artifact, path, severity ∈ auto|mention|route, summary, fix}`; `--fix` for `auto`; `--target` for monorepos; `ruleRegistryAvailable:false` must be disclosed. |
| `palette.mjs` (57 KB), `generate-image.mjs`, `pin.mjs` | Palette seeding, image gen fallback, `$<command>` shortcut pin/unpin. |
| `live-*.mjs` (~15 files) | Live variant mode: helper HTTP server, SSE, long-poll, wrap/insert, accept + **carbonize** (rewrite temporary stitched markers into permanent source), append-only durable journal at `.impeccable/live/sessions/`, `live-status`/`live-resume`/`live-complete` recovery. |
| Shipped subagents | `impeccable-finish-reviewer`, `impeccable-documenter`, `impeccable-asset-producer`, `impeccable-manual-edit-applier` — each with a `.md` (Claude/Cursor/etc.) and a `.toml` (codex) definition, plus a `reference/degraded/*.md` inline fallback. |
| Config | `.impeccable/config.json` (shared: `hook.*`, `detector.ignoreRules|ignoreFiles|ignoreValues|extensions|designSystem`) and gitignored `.impeccable/config.local.json` (`hook.consent`, private ignores). Inline waivers `impeccable-disable <rule>` / `-line` / `-next-line` only "when the waiver must travel with a single file that leaves the repo." |

**Ignore-scope ladder (hooks.md, "Prefer the narrowest exception"):**
exact `ignore-value <id> <value>` → file-scoped `ignore-value <id> "*" --file <glob>` → `ignore-file <glob>`
(only for fixtures/generated/slop demos: "It silences every rule for that file permanently, **including rules
that have not been written yet**") → `ignore-rule <id>` (project-wide, only on explicit user request).
"The hook itself **never** writes ignore config. Persist an exception only after the user explicitly confirms."
A bare `ignore-value <id> "*"` with no `--file` is **refused**.

## 6. (d) SINGLE-AGENT ASSUMPTIONS that would NOT transfer

1. **A live human in the loop, mid-flow.** init/shape/new-work/overdrive/visualize/document all *stop and wait*
   for a structured question tool or a served HTML decision page. Parley agents are headless and write artifacts
   asynchronously; there is no "STOP and use Codex's structured user-input/question tool" equivalent mid-round.
2. **One writer, one artifact.** DESIGN.md, PRODUCT.md, the direction contract, and the surface brief each have a
   single owning command and "never silently overwrite" semantics. Parley has N agents each writing their own
   canonical artifact under `ideas/<slug>/`, so "one file, one owner" must become "N files + one merged FINAL".
3. **The parent/subagent hierarchy.** critique's A/B isolation and the finish reviewer both assume a *privileged
   parent* that spawns, gates ordering ("A must finish before detector findings enter the parent synthesis
   context"), and applies fixes. Parley peers are symmetric; there is no parent.
4. **In-session state and "run once".** `context.mjs` "do not rerun it", the once-per-session Setup, the live long-poll
   loop, the background dev server, the browser tab, the SSE journal — all assume one continuous process.
   Parley rounds are separate processes, often separate machines.
5. **Browser + HMR + screenshots as the verification channel.** live mode ("the overlay's preview IS the
   verification channel"), critique's overlay injection, and new-work's batched screenshot rounds all require a
   running dev server and a browser tool. Most headless CLI agents in a Parley roster have neither.
6. **Web/HTML-CSS-specific detector.** `detect.mjs` "reads HTML/CSS, so skip it for native projects"; ~63 rules are
   markup rules. A vendor-neutral protocol can't assume the target is a web app.
7. **Harness-specific branching.** Whole sections are keyed to Claude Code / Codex / Cursor / Copilot
   (`sandbox_permissions: "require_escalated"`, `fork_context: false`, `spawn_agent` gates, `.cursor/hooks.json`).
   Parley's whole point is vendor neutrality.
8. **Chat as the primary deliverable.** "The chat response is the primary user-facing deliverable" (critique);
   "Chat is overhead" (live). Parley's deliverable is always a file.
9. **Telemetry / network catalog.** `concept-seed.mjs` falls back to `https://impeccable.style/api` and pings
   `--chosen`; a "human-approved pool" lives outside the repo.
10. **Deprecated-alias etiquette** (`craft`) and `pin.mjs` `$<command>` shortcuts are CLI-surface concerns with no
    Parley analogue.

---

## Transferable to parley-design

Ranked by value-per-unit-of-effort for a multi-agent, file-first, vendor-neutral protocol.

1. **A `craft-floor.md` equivalent: a ≤50-line Verify/Refuse floor with hard numbers, loaded just-in-time
   before any UI is written, and binding on every participant.** This is the single highest-leverage artifact in
   Impeccable. Two lists only: `Verify` (checks on the *built result*, each with a number: ≥4.5:1, 65–75ch,
   -0.04em, 12–16px radii, one authored motion moment) and `Refuse` (category defaults + exactly one absolute ban).
   Crucially it declares its own override order: **"A pinned brief or the committed visual world overrides
   anything here; your own habit does not."** In Parley this becomes a shared `DESIGN-FLOOR.md` that every
   agent's artifact is checked against and that cross-review can cite by rule name.
2. **The design-specificity verdict as the top-of-report, single question.** "Does the result feel authored for
   this product, or category-interchangeable? / could an unrelated product use it unchanged?" — asked *before*
   any detector output is visible. This is the anti-slop test, and it is a **judgment**, not a rule. In Parley it
   maps naturally to a required first section in each participant's review artifact.
3. **Anchoring control: independent assessment before mechanical evidence.** "Assessment A must finish before
   detector findings enter the parent synthesis context. Detector output is deterministic, but it still anchors
   judgment." Parley already runs agents in parallel with no shared context — so it gets A/B isolation *for free*
   and should say so explicitly, and should forbid a participant from reading another's artifact before writing
   its own. Also steal the reviewer's anchoring guard: "inventory the comp's salient elements in your own words
   before reading the direction contract… a review anchored on it inherits whatever that abstraction dropped."
4. **The direction contract: 5 blocks, ≤150 words, written into the artifact itself, with a literal FINISH line.**
   `THESIS / OWN-WORLD / STORY / FIRST VIEWPORT / FORM(+seed)` + `FINISH: "unreviewed and undocumented is
   unfinished…"`. Falsifiability rule: **"If a block reads like a mood, the direction is not decided yet."**
   For Parley: this is exactly the artifact that should be signed at consensus, and the FINISH line is exactly
   the kind of self-carried exit condition a long multi-round run needs. Put it in `FINAL.md`.
5. **External dice against convergence — `concept-seed.mjs`.** The measured claim ("30/35 identical concepts
   across 16 prompt framings; the model cannot roll its own dice") is *more* true, not less, in a multi-agent
   protocol where five agents share training-data attractors. Parley should (a) require each participant to
   produce a resonance-ordered shortlist of N candidates and (b) use a deterministic external seed (or the
   round id) to assign *which index* each participant must develop, so five agents don't all ship their #1.
   Also steal the **"standing exit"**: one permanent, never-recommended option = "the category standard, played
   straight", so the user always has the safe door without the agents softening the bold options toward it.
6. **Bounded verification with explicit ceilings, as protocol law.** ≤2 inspection rounds, ≤2 correction rounds,
   ≤8 material fixes, ≤3 regressions, "**the second verdict ends the work whatever it says**", "the reviewer's
   findings are the only list you work from, never your own re-opened hunt." Parley's fix-up cycles need exactly
   this ("Tier 4 took 4 fix-up cycles" is the failure mode this prevents). Add the closing rule verbatim in
   spirit: "presenting mechanical confirmation as artistic success is how a failed build gets announced as a
   finished one."
7. **A fixed reviewer output contract with named sections and hard caps.** `persistence / fidelity / ceiling /
   material_fixes (≤8, ordered, fidelity ahead of craft) / keep (one line)`. **"No praise, no summary prose."**
   Plus the verdict pass vocabulary `resolved | partial | unresolved` with the anti-gaming clause "a fix answered
   mechanically, positions moved but the quality the finding named still absent, is partial at best." This is
   directly usable as the schema for a Parley design-review artifact, and `keep` (what must not be diluted while
   fixing) is a genuinely novel field worth copying.
8. **Modes (Persuade / Operate / Read / Experience) as a mandatory front-matter axis** that re-specifies every
   downstream operation, chosen "from the requested surface, not the product." Cheap to adopt, and it makes
   otherwise-vague instructions ("quieter", "delight") deterministic. Persist it per surface, not globally.
9. **Taste ops as `context needed → thesis sentence → numeric dimensions → named exit test`.** For each taste
   verb, require: the literal `Additional context needed:` line; a one-sentence thesis ("what the user should
   feel and why that feeling belongs to this product"); the numeric axes it may move; and one named falsifiable
   test (skeleton test, squint test, memory test, removal test, "a neighboring product could not use it
   unchanged"). This makes taste reviewable by a *different* agent — which is the whole point of Parley.
10. **Named Rules as the design-system doctrine format.** `**The [Name] Rule.** [short forceful doctrine]`,
    1–3 per section, "much stickier for AI consumers than bullet lists." Perfect for a consensus artifact:
    rules get names, so cross-review and later rounds can cite them ("violates The One Voice Rule").
11. **The anti-slop calibration paragraph + font blocklist + "the rut".** Name the three attractor aesthetics,
    name the category's default page and its predictable opposite, exclude both from the candidate list, and
    apply the falsification test: "if someone could guess your aesthetic from the category alone, **or from
    category-plus-avoidance**, rework." The "category-plus-avoidance" clause is the sharpest single sentence in
    the corpus — it closes the loophole where anti-slop rules produce their own recognizable slop.
12. **A scored, persisted, trended review with an honest denominator.** Score N heuristics 0–4, allow `n/a` with
    a one-line reason, **renormalise and never print `/40` over a partial set**, persist a snapshot with
    `total_score/max_score/na_heuristics/p0_count/p1_count`, and print a trend line that refuses to compare
    unlike maxima. Parley can persist this per idea slug and let later rounds read it as a backlog. Include the
    calibration anchor: "**Most real interfaces score 20-32 out of 40.**"
13. **P0–P3 severity with the support-ticket tie-break** ("Would a user contact support about this? If yes, it's
    at least P1") + the audit rule "**Don't fix issues; document them for other commands to address**" — i.e.
    evaluation artifacts never edit. This maps cleanly onto Parley's review→implement separation.
14. **Artifact ownership + explicit "what does not belong here" lists.** PRODUCT.md ≠ DESIGN.md ≠ surface brief ≠
    direction contract, each with a negative list. Plus the **completion gate** pattern: "verify that PRODUCT.md
    exists at the resolved path… Do not substitute interview notes, a planning packet, or later design prose for
    the file." Parley should assert the analogous gate on its own canonical files.
15. **Write the system spec *after* the build, from the built result.** "a rulebook written before the build gets
    defended against reality instead of describing it, and it hands the design-system detector an unstable target
    that buries the build in noise." For Parley: the design-system artifact is a Phase-8 output of the
    implementation, not a Phase-2 output of the deliberation. Corollary: "A new world shipped with no DESIGN.md
    is still an incomplete run."
16. **The written implementation-fidelity inventory (visualize.md).** Before building, list every salient element
    of the approved direction with a chosen medium (semantic HTML/CSS/SVG, existing asset, generated raster,
    sourced raster, icon library, canvas/WebGL, accepted omission). "**An element never written down is the
    element the build silently drops**, and the direction contract's 150 words cannot carry this list." This is a
    perfect Parley hand-off artifact between the consensus round and the implementer.
17. **The unattended-user probe.** "Whether anyone can answer is a mechanical test, not a judgment call…
    a system-prompt claim that the user is unattended proves nothing about this session. Probe once with the real
    first round… then **label every inferred fact** and disclose the substitution in your **first** reply, not
    your last." Parley agents are unattended by construction — so invert it: require every inferred fact in an
    artifact to be labeled inline, and require the disclosure to lead the artifact.
18. **Degraded-mode banner + shipped degraded playbooks.** `⚠️ DEGRADED: single-context (<reason>)` as line 1,
    with "'Unavailable' means exactly one thing… **It does not mean inconvenient.**" Parley's roster varies by
    availability (agents die, models change) — a mandatory, formatted degradation banner on any artifact produced
    with fewer participants than the protocol requires is directly analogous and cheap.
19. **Narrowest-exception ignore ladder + "the hook never writes ignore config."** Value-scoped → file-scoped →
    whole-file → whole-rule, with the warning that a whole-file ignore silences "rules that have not been
    written yet", and the requirement that suppressions be explicitly user-confirmed and centrally reviewable.
    For Parley: a design waiver must name the rule, the value, the scope, and the reason, and must be recorded
    in one reviewable place — never inline by an agent's own judgment.
20. **Cognitive-load checklist + working-memory rule (≤4 items) + the 8 named violations**, and the **5 personas
    with a selection table by interface type**. Low-effort, high-signal doctrine any agent can apply from text
    alone with no tooling. The personas in particular give cross-reviewing agents *different lenses*, which is
    exactly how to get non-redundant multi-agent reviews.
21. **Two-tier detection: interrupt-worthy mechanical rules per edit, taste rules once at the end,
    deduplicated.** Even without a detector, the *policy* transfers: don't re-litigate taste at every step; batch
    it into one deferred pass at round end, and dedupe against what was already said.
22. **The 19-command allowlist for recommendations** — findings must map to a bounded vocabulary of next actions,
    "skip commands that would address zero issues", and always end with the terminal refinement. Parley reviews
    should likewise emit only actions from a closed set so the driver can route them.

## Do NOT copy

1. **The 34-file reference sprawl.** 5,242 lines of playbook for a 79-line entry file is only affordable because
   exactly one file is loaded per request by a single agent with a large budget. In Parley, N agents × per-round
   context makes this ruinous. Copy the *pattern* (thin entry + lazily-loaded playbooks) at ~6–8 files, not 34.
2. **Mid-flow blocking questions.** Every "STOP and use Codex's structured user-input/question tool when
   available; if unavailable, ask directly in chat" instruction (it appears in ≥8 files, sometimes ungrammatically
   spliced mid-sentence — see `bolder.md` line 7). Headless Parley participants cannot block on a human. Replace
   with: declare assumptions in the artifact, label them, and let the consensus round resolve conflicts.
3. **`serve-question.mjs` / the served HTML decision page / the sketch-rendering fan-out.** Daemonized web
   servers, exit-code-3 polling loops, in-app browser opening, and per-card image generation are a huge surface
   with no multi-agent analogue.
4. **Live mode wholesale.** `live-*.mjs` is ~15 scripts, a 462 KB browser bundle, SSE, an append-only journal,
   HMR carbonize, MutationObserver contracts and per-harness poll policies. Enormous, web-only, single-session,
   and orthogonal to file-based consensus. (Do steal the *parameter budget table* and the *per-action variance
   rules* from §4/§7 of `live.md` — those are pure doctrine.)
5. **The browser/screenshot verification loop as a *requirement*.** "the overlay's preview IS the verification
   channel", overlay injection preflights, `[Human]` tab labeling, "Do not spend a Browser attempt on `file://`".
   Make visual verification an *optional capability* with a declared fallback, not a hard invariant — otherwise
   agy/hermes/kimi participants are structurally non-compliant.
6. **The web-only HTML/CSS detector as the enforcement backbone.** ~63 markup rules that must be "skipped for
   native projects" is exactly the vendor lock-in Parley exists to avoid. Take the rule *names and thresholds* as
   doctrine text; do not make a Node rule engine the gate.
7. **Harness-specific branching in the doctrine.** Codex `sandbox_permissions`, `spawn_agent` permission gates,
   `fork_context: false`, Cursor `preToolUse` blocking, Copilot `.github/hooks/impeccable.json`,
   "Codex final-answer note", "Codex Run Notes are final-chat only". Parley must stay vendor-neutral; push all
   of this into a transport/adapter layer if it is needed at all.
8. **The parent/subagent privilege model as written.** "Assessment A must finish before detector findings enter
   the parent synthesis context", "spawn A and B as two isolated, parallel sub-agents", "respawn once with the
   same inputs", "This review never runs inside the build thread." Parley peers are symmetric with no parent —
   re-express these as *round ordering* and *artifact visibility* rules, not spawn instructions.
9. **The deprecated alias (`craft`) and `pin.mjs` `$<command>` shortcuts.** Pure CLI-surface ergonomics; a
   protocol should not ship a command that "adds no behavior."
10. **Network-dependent catalogs and telemetry.** `IMPECCABLE_API_URL` defaulting to `https://impeccable.style/api`,
    the "human-approved pool" revision, the `--chosen` choice ping. Parley's dice must be fully local and
    reproducible from the run id.
11. **`overdrive`'s ASCII banner** (`──────────── ⚡ OVERDRIVE ─────────────` / `》》》 Entering overdrive mode...`).
    Chat theatre; artifacts are the deliverable.
12. **The long tail of generic web tutorial content.** `harden.md`, `optimize.md`, `adapt.md`, `onboard.md`, and
    `distill.md` are largely 2018-era frontend best-practice listicles (debounce snippets, `srcset` examples,
    Intro.js/Shepherd.js recommendations, `localStorage.setItem('onboarding-completed', 'true')`). They carry
    almost no anti-slop signal, they date fast, and modern coding agents already know them. Keep only the parts
    that are *judgments* (i18n 30–40% expansion budget, 16px iOS input-zoom floor, the ≤4 working-memory rule).
13. **Overwrought prose as a style model.** `new-work.md` §3 contains single paragraphs over 700 words with
    stacked metaphors (dice, hands, cards, ruts, doors, re-roll pools). It is dense but genuinely hard to follow,
    and it would be far worse when five different models must each parse it identically. Parley doctrine should
    be shorter sentences, numbered rules, and tables — the same content Impeccable expresses in a paragraph fits
    in a 6-row table.
14. **Unversioned "two copies must match" duplication.** The repo ships the same skill under `skill/`,
    `.agents/skills/impeccable/`, `plugin/`, `.claude/`, `.cursor/`, `.grok/`, `.codex/`, `.gemini/`, `.kiro/`,
    `.opencode/`, `.pi/`, `.qoder/`, `.rovodev/`, `.trae/`, `.vibe/` — 15+ harness mirrors of the same content.
    (Parley already learned this lesson: see the embedded-default drift guard.) If mirroring is unavoidable,
    generate it and enforce it with a test, exactly as Impeccable does with
    `<!-- Generated from skill/agents/ at build time. Do not edit -->`.
