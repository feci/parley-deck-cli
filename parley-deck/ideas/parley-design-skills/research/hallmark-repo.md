# Hallmark — repo digest (everything EXCEPT the SKILL.md body)

Source root: `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/research/hallmark`
Repo: `Nutlope/hallmark` · npm name `hallmark` · **version 1.1.0** · MIT · "Made by Together AI" (author Hassan El Mghari / `@nutlope`, `@YoussefXLM`).
Studied: `README.md`, `ROADMAP.md`, `package.json`, `vercel.json`, `.gitignore`, `docs/recipes.md`, `docs/study-examples.md`, `docs/talk-slides.md`, `site/_tests/**` (README, 8 `brief.md`, 3 full HTML+CSS pairs, `all-themes.html`, `verbs/`, `custom/`), `site/css/tokens.css`, `site/js/main.js`, `site/index.html`, `site/examples/*/.hallmark/log.json`, plus the *non-SKILL.md* reference files that carry machinery (`slop-test.md`, `custom-theme.md`, `contract.md`, `design-md.md`, `export-formats.md`, `preview-examples.md`).

---

## 0. Repo shape at a glance

```
skills/hallmark/SKILL.md                 (67 KB — covered by the other analyst)
skills/hallmark/references/*.md          24 top-level refs, 5,630 lines total
skills/hallmark/references/components/   50 files (h1–h9, n1–n13+n1b, ft1–ft8, f1–f6, s1–s5, t1–t4, c1–c4)
skills/hallmark/references/macrostructures/  21 files (01-bento-grid … 21-component-playground)
skills/hallmark/references/genres/       4 files (atmospheric, editorial, modern-minimal, playful)
skills/hallmark/references/themes/       4 deep files ONLY (carnival, cobalt, hum, lumen)
skills/hallmark/references/verbs/        2 files (audit, redesign)
docs/{recipes.md, study-examples.md, talk-slides.md, screenshots/}
site/                                    marketing site + the whole eval corpus (_tests/, examples/)
package.json / vercel.json               packaging + static deploy (outputDirectory: "site")
```

Biggest reference files (lines): `custom-craft.md` 626 · `study.md` 509 · `hero-enrichment.md` 475 · `anti-patterns.md` 418 · `assets.md` 406 · `custom-theme.md` 367 · `export-formats.md` 329 · `component-cookbook.md` 265 · `microinteractions.md` 258 · `typography.md` 243 · `slop-test.md` 192 · `contract.md` 25.

---

## 1. Positioning — claims, scope, refusals

### The one-line claim (README.md:3)
> **A design skill for Claude Code, Cursor, and Codex that refuses to look AI-generated.**

### The mechanism claim (README.md:13, verbatim)
> "Hallmark picks a macrostructure for the brief, dresses it in one of twenty themes, runs **fifty-seven slop-test gates plus a pre-emit self-critique**, and refuses the on-distribution defaults every LLM was trained into. **Two pages by Hallmark for two different briefs feel like different sites, not colour-swaps of the same template.**"

Note the drift: README says "fifty-seven"; `references/slop-test.md` line 1 says **"Slop test — 58 gates + pre-emit self-critique"** (gates are numbered 1–57 plus an extra `38a`, so 58); `site/js/main.js` marketing copy says "57-gate slop test" and "Fifty component archetypes — fourteen navs, eight footers, four hero polish patterns". The count is a *marketing number that drifts*; the file is authoritative.

### The differentiator (SKILL.md:13 — quoted because it is the thesis)
> "The differentiator: Hallmark insists on **structural variety**, not just visual variety. Two pages by Hallmark for two different briefs should not share the same hero → 3-feature → CTA → footer rhythm."

### Four verbs (README.md table, verbatim)
| Verb | What it does |
| --- | --- |
| *(default)* | Build new UI. Picks a macrostructure, applies the rule-set, runs the slop test before handing back. |
| `hallmark audit <target>` | Score existing code against the anti-patterns. **Punch list, no edits.** |
| `hallmark redesign <target>` | Throw out the structure, keep copy + IA + brand, rebuild with a different fingerprint. |
| `hallmark study <screenshot \| URL>` | Extract the **DNA** from a design you admire: macrostructure, type-pairing, colour anchor. **Refuses pixel-clones and paid templates.** Optionally emits a portable `design.md` for handoff to other AI tools. |

(A fifth verb, `refine`, is referenced by the verb tests — `verbs/README.md` mentions "one worked example per verb (`audit`, `refine`, `redesign`, `study`)" and `audit-report.md` recommends "do NOT run `hallmark refine`" — but the `refine/` folder does not exist in the repo and `references/verbs/` holds only `audit.md` + `redesign.md`. **Documented-but-missing verb; a real drift bug.**)

### What it explicitly refuses to do (`references/contract.md` § Scope and limits, verbatim)
> "Hallmark is a *taste* skill. It will not:
> - Invent product copy. If the user hasn't given you the words, ask.
> - Pick a brand identity. It will follow one you give it.
> - Enforce a specific style (dark mode, glassmorphism, brutalism). It will execute whichever genre + tone the user committed to.
> - Build logic — state management, data fetching, business rules. It is a visual / interaction layer only.
>
> If a request falls outside taste — 'build the auth flow', 'wire up Stripe' — do the work, but apply Hallmark to the rendered surface."

Other hard refusals:
- `study` **refuses**: "paid-template-marketplace listings; copy-protected portfolios without permission" and **never names an exact font** ("visual font identification is unreliable") — `docs/study-examples.md` § "What `study` doesn't do", items 1–4.
- Gate 46: any quantitative claim the user did not supply is an **invented metric** → auto-fail.
- Gate 47: hand-built fake browser bars / phone frames / terminal chrome → auto-fail ("Re-drawn chrome is one of the strongest 'looks AI-generated' tells").
- `contract.md`: "**An existing global stylesheet is append-only**… Do a full rewrite only when the user explicitly asks for one: silently dropping a framework's CSS entry directives un-styles the entire app."
- `design-md.md`: "**No-overwrite policy.** If `design.md` already exists at the project root, do NOT overwrite. Refresh its `## Exports` section instead."

### The "Custom" branch positioning (README.md:75-90)
> "When a brief carries creative intent that no catalog theme fits, Hallmark switches to **Custom** and designs the page from scratch: a made-to-measure palette, type, and layout. Same 57 slop-test gates, no template underneath. **It stays a quiet branch; vanilla briefs never see it.**"

### The talk (`docs/talk-slides.md`) — positioning in a third-party frame
17 slides, 18 minutes, AI Engineer World's Fair 2026. Explicitly *not* a product pitch: "Hallmark is mentioned once, on slide 15, as one open-source example our team built." Slide 04 diagnosis: **"Models ship the mean of their training set."** Slide 16 thesis: **"Design taste is a layer in the agent stack"** — `MODELS / TOOLS / MEMORY / RAG / EVALS / TASTE ← here / ORCHESTRATION`. The seven fixes (compressed playbook, slide 13, verbatim):
```md
# system_prompt.md
You are a designer, not a templater. Before you write code:
  1. use shadcn/ui — never re-invent a primitive that already exists.
  2. read  DESIGN.md  for tokens. Use only those.
  3. pick ONE page shape (marquee, long-document, bento, stat-led, manifesto, catalogue).
  4. ban: purple-cyan gradients, Inter on Inter, 3-col icon grids, 100vh centered hero,
     card-in-card, background-clip text, pure #000 / #fff, invented metrics.
  5. ship 8 states for every interactive element.
  6. score  Distinctiveness · Hierarchy · Execution · Restraint  on 1–5.
     anything <3 triggers a revision pass.
```
Slide 12 rubric (four axes, 1–5, "**Anything below 3 → revision pass**") is the public/simplified version of the six-axis pre-emit critique in `slop-test.md`. Slide 11 ships the **eight states** table: `default · hover · :focus-visible · :active · disabled · loading · error · success`, with "If the agent emits two, your UI ships broken."

---

## 2. ROADMAP — what the author considers unfinished or hard

`ROADMAP.md` has three buckets. Verbatim headline items:

**Now (one item only)**
- **"Nanobanana hook for image-heavy briefs."** "Today the integration is recommend-only — Hallmark tells the user to go generate something and bring it back. Image-heavy briefs (e-commerce, travel, food, lookbook) route to typography-only and **feel underserved.**" Wants: write prompt → invoke API → ingest image → wire into build, **cache by prompt hash**. Plus a new image-led theme (working title *Plate*).

**Next (8 items)**
1. **Brand-first flow** — from a product description generate palette + type + voice + imagery and lock it into a `design.md`; then every page builds against that brand. "Closes the gap for users who have a product idea but no brand yet."
2. **Theme-aware motion tokens** — per-theme `--dur-micro` / `--dur-short` / `--dur-long`. "**Atelier should feel slower than Brutal; today they share durations.**"
3. **`hallmark variant`** — "produce three structurally distinct versions of the same brief side-by-side… **The biggest cause of 'AI feel' is users accepting the first output because they didn't know it could be different.**"
4. **Structural cookbook** — `structure.md` "catalogues the *axes* of variety but doesn't show what a left-margin-headed, hairline-divided, no-image page actually looks like assembled. Twelve to twenty worked fingerprints… **patterns are easier to reach for than principles.**"
5. **Tactile-rebellion reference** — controlled imperfection: "a 0.5° tilt on one mark is taste; on every word it's chaos."
6. **Charts reference for analytics pages** (`data-viz.md`) — "AI-generated charts are an obvious tell — rainbow palettes, dense gridlines, 3D donuts, dual-axis line spaghetti… **Half of every dashboard is chart-shaped, and Hallmark currently has nothing to say about it.**"
7. **Multi-page coherence** — "**The structural-variety rule is correct for variety, wrong for brand consistency inside a multi-page product.** Lock the brand axes (type, colour, divider language); vary the page-voice axes (heading placement, body composition, button voice)." ← This is the single most important admitted tension in the whole project.
8. **`study` reads your own codebase too** — third input mode: a path to your project → emits the same `design.md`.

**Later (5 items)**
- `hallmark explain` (narrate choices axis by axis — "The skill teaches; users start making the same calls themselves").
- **Negative-capability rules** — "for each anti-pattern, the perceptual or cognitive reason it fails. **Understanding it beats knowing it.**"
- **Emotion-first prompting** — "*nostalgic · optimistic · sceptical* instead of *editorial · brutalist · austere*. Today's tone words don't reach."
- Sound + haptic policy.
- **Live preview as an MCP server** — "watch the file, render in a sandbox, screenshot, feed the screenshot back for self-critique against the slop test. **Closes the loop between generation and audit.**" ← the admission that today there is **no rendering feedback loop at all**.

Deferred-forever items recorded in `site/_tests/README.md`: **#12 `study + redesign` combined verb** and **#14 auto-generated OG cards** ("Needs Node + Puppeteer utility") — both marked 🟡 across two release rounds.

---

## 3. Recipes + study examples — what a real invocation looks like end-to-end

### 3a. `docs/recipes.md` — 9 worked briefs (00–08)

The file's own spec of its shape (§ "How the recipes are organised", verbatim): every recipe has exactly four lines —
1. **Prompt** — verbatim, copy/paste-ready.
2. **Inferred trio** — what the design-context-gate produces (**audience · use · tone**). Marked `explicit` if the user provided all three; `the user opted out` if not.
3. **Picks** — **macrostructure · theme · enrichment archetype**.
4. **Excerpt** — one paragraph of the produced copy or layout.

Recipe 00 is the **canonical try-it prompt** — a self-test that the skill is wired correctly:
> *"Build me a landing page for Coffeebox — a small-batch coffee subscription. Roast on Sunday, ship on Monday, drink Tuesday. Audience: people who already buy good coffee and want fewer trips to the shop. Tone: warm, hand-set, editorial — like a small café's chalkboard."*
>
> **Why this is the canonical try-it:** … "If Hallmark produces *Linen-with-italic-Cormorant-and-warm-paper* for this prompt, **the skill is wired correctly**."

Recipe picks table (compressed from the file):

| # | Brief | Trio | Macrostructure · Theme · Enrichment |
| --- | --- | --- | --- |
| 00 | Coffeebox coffee subscription | explicit | Long Document · **Linen** · Tier-B hand-built SVG |
| 01 | Tide indie podcast ("just go ahead, pick the rest yourself") | opted out | Letter · **Salon** · none (typography only) |
| 02 | Streampipe CLI ("Use the Terminal theme") | explicit + theme requested | Long Document · **Terminal** · Tier-A inline CSS-art terminal |
| 03 | Maple Street Bread bakery | explicit | Catalogue · **Almanac** · Tier-A inline-SVG bread silhouettes |
| 04 | Meridian studio manifesto ("no flashy stuff") | partial | Quote-Led · **Brutal** · none |
| 05 | Tracejam SaaS observability | explicit | Workbench · **Midnight** · Tier-A pure-CSS sticky trace panel |
| 06 | Anya personal one-pager ("Don't ask, just figure it out") | opted out | Index-First · **Plain** · none |
| 07 | Foundry compliance SaaS | explicit | Bento Grid · **Newsprint** · Tier-A pure-CSS logo grid |
| 08 | Cohort courses platform | explicit | Stat-Led · **Linen** · none |

Closing rules (§ "What the recipes are *not*", verbatim):
- "**Not templates.** Hallmark's whole point is structural variety — duplicating a recipe verbatim is the **Specimen-fall-through** anti-pattern (gate 23 in SKILL.md)."
- "**Not fixed picks.** Two consecutive runs of recipe 00 (Coffeebox) on the same project should produce *different* macrostructures or themes — the `.hallmark/log.json` project memory enforces this."
- "**Not exhaustive.** The **21 macrostructures × 16 themes × 8 enrichment archetypes = 2,688 distinct fingerprints.** The 9 recipes here are a starter set."
  (Arithmetic drift: the catalog is now 20 themes → 21 × 20 × 8 = 3,360. The doc still says 16/2,688.)
- Usage instruction: "look for the closest match in this file and **note what's different** between your brief and the recipe. The macrostructure pick is usually inherited; the theme often shifts (per the diversification rule); enrichment may change tier."

### 3b. `docs/study-examples.md` — three worked `study` invocations

Purpose sentence (verbatim): "The whole point of `study`: extract **structural DNA** (macrostructure + archetypes + type-pairing role + colour anchor + rhythm), **never pixels**. If the user wants a 1:1 copy, they want a different tool."

Each example runs the same 4-beat protocol: **Refuse-or-proceed check → Diagnosis report → Confirmation question → Build**.

- **Example A · Pentagram-style editorial portfolio.** Diagnosis is 10 structured fields (see `verbs/study/diagnosis.md` for the canonical list): macrostructure, hero archetype, pitch archetype, footer archetype, display family role, body family role, surface lightness band, accent hue band + chroma, density verdict, type-pairing role. Names two *candidates* per role (one free / one paid) and never an ID: *"Fraunces italic display + Geist body"* free, *"Tiempos Headline + Söhne"* paid.
- **Example B · Klim-leaning type specimen.** The confirmation question actively *argues with the user*: "This DNA is for a specimen page, not a marketing page. If you're rebranding an agency, the glyph hero will be impressive on the first visit and pointless after. Suggested adaptation: Type Specimen as the homepage, then a Long Document or Workbench pattern for the about / case-studies / contact pages. Pick: keep all-Specimen, or hybrid?"
- **Example C · Rauno-Freiberg-style personal site.** Fires the **ambiguity-ask** branch: "Is this your own work, a public reference for inspiration, or someone else's live site? …If it's a paid portfolio template, I won't reproduce it."

The **anti-pattern-carry-over gate** is the most transferable idea here. Every diagnosis ends with *"N anti-patterns the screenshot has that you should NOT carry over"* + concrete numeric corrections, e.g.:
> "The hover-state … uses a slow ease-in-out **800 ms** colour fade. Per `microinteractions.md` § The timing canon, hover state should be **150–200 ms**. **Carry the *idea* (colour-shift on hover); shorten the *duration*.**"
> "The footer text is set at **11 px** with 0.06em tracking — under the **14 px floor** for body copy. Bump to **12-13 px**, keep the tracking."
> "grid lines … 1 px solid white at 0.4 opacity. On a 4K display this becomes a sub-pixel that disappears or aliases. Use 1 px solid + `color-mix` down to 35% lightness instead."

Example C also carries the inverse — a **"do carry this over"** note:
> "The italic name in the top-left has a baseline at the same y-position as the centred demo's vertical centre — a subtle horizontal-baseline alignment that makes the page feel composed rather than stacked. **Carry this; it's invisible until removed.**"

---

## 4. The `_tests/` corpus — what a brief contains, how quality is judged, whether there's an eval harness

### 4a. Corpus layout

```
site/_tests/README.md            the eval charter + findings + release notes
site/_tests/01..08-<slug>/       brief.md + index.html + style.css     (8 canonical tests, v0.6.0)
site/_tests/09..13-<slug>/       index.html only (self-contained, inline <style>, newer generation)
site/_tests/custom/{01..03}/     index.html + style.css (custom-route worked examples)
site/_tests/verbs/{audit,redesign,study}/   input + output + notes.md per verb
site/_tests/all-themes.html      20-theme live-preview review sheet (iframes into ../index.html?theme=KEY)
site/_tests/_thumbs/             1024×640 PNG/JPG screenshots + 2 MP4 loops for the gallery
site/examples/<slug>/            the *published* gallery builds (index.html + styles.css + tokens.css [+ script.js])
site/examples/<slug>/.hallmark/log.json    per-project memory (4 present: bananastudio, hyperlane, najm, wayfare)
```

### 4b. What a `brief.md` contains — the artifact is a **reasoning trace**, not a spec

Every `brief.md` is the *same six-step transcript* (`site/_tests/README.md`: "**`brief.md`** — the verbatim prompt + the new design flow (Step 0 Pre-flight, Step 1 Context, Step 2.5 Rotation, Step 5 Preview, Step 6 Stamp)"). Structure, verbatim from `01-tide-podcast/brief.md`:

1. `## The prompt (verbatim, same as v1)` — the raw user string in a blockquote.
2. `## Step 0 · Pre-flight` — `> "No pre-flight signals — proceeding with full Hallmark stack."` (i.e. no `package.json`, no `tailwind.config.*`, no existing CSS).
3. `## Step 1 · Design-context gate` — one of *fully answered / partial → inferred / skipped*. When skipped, the skill must state its inference in **one sentence**: `> "Going with: audience = … · use = … · tone = …. If any of those is wrong, redirect."`
4. `## Step 2.5 · Project memory rotation` — the rotation argument, in prose, citing the previous pick and naming the **candidate set** it chose from. Verbatim (test 05):
   > *"Previous run on this brief picked Bento Grid + Pastel + Tier-A flame chart. Picking from {Workbench, Stat-Led, Long Document} for the macro this time — Workbench wins because the brief says 'try it' and Workbench is structured to walk an SRE through a single workflow rather than show six surfaces in a tile grid."*
   > *"Theme rotation: Pastel (light · geometric-sans · cool-indigo) → Midnight (dark · geometric-sans · phosphor-cyan). Differs on paper band (light → dark) and accent hue (indigo → phosphor-cyan). **Two of three axes differ. Passes.**"*
   And the escape hatch (test 02): *"Theme: Terminal as requested. **No theme rotation when the user names one.**"*
5. `## Step 3 · Visual ruleset loaded` — an **explicit list of which reference files were loaded and which section of each** (e.g. `microinteractions.md` (single primitive: caret blink only)). This is a *token-budget audit trail*.
6. `## Step 4 · Hero enrichment` — a quoted justification, usually a refusal: *"Enrichment: none. Quote-Led + Brutal is already loud; adding any visual would dilute the polemic. The 'no flashy stuff' constraint maps directly here — voice is the design."*
7. `## Step 5 · Preview` — a fenced markdown block emitted **before any code**, always 8 rows: Macrostructure · Theme · Enrichment · Sections · Motion · Slop test (`38 / 38 ✓`) · Diversification.
8. `## Step 6 · Macrostructure stamp` — the CSS comment that will head the file.
9. `## What changed vs v1` / `## What stayed the same` — a **per-axis diff against the previous generation of the same brief**.

### 4c. How output quality is demonstrated / judged

There is **no automated harness**. Judgment is (i) hand-written prose findings, (ii) a self-reported gate score in the artifact, and (iii) a cross-run diff table.

**The charter (`site/_tests/README.md` lines 30-36, verbatim)** — four questions the corpus exists to answer:
> 1. Does the skill default to the same shape twice across the six tests? → No two should share a macrostructure or hero archetype.
> 2. Does the skill enrich every page even when it shouldn't? → Tests 01, 04, 06 should be typography-only.
> 3. When the user skips the design-context gate, does the skill state its inferences in one sentence at the top? → Yes, see brief.md headers.
> 4. Does the skill bend toward "Built for the modern team" voice when copy is left to it? → Each brief has its own voice; no template-soup.

**The pass criterion (verbatim):** "**Eight distinct fingerprints across eight prompts.** No macro repeats. No theme repeats. Every adjacent pair differs from its neighbour on at least one of: **macrostructure / display style / accent hue / paper band** — the v0.6.0 rotation rule firing."

**The most honest artifact in the repo — § Friction points I hit while generating** (author generating by hand, recording where the skill under-specified):
- "The 'domain → trio' table covers most cases but not all… I had to combine 'SaaS marketing' + 'developer tool' mentally to land on Bento Grid for Tracejam."
- "**Theme picking is still by feel.** The skill has a strong rule that 'the next theme can't be Specimen if the last one was Specimen' — but no rule like *'the next theme should be categorically different on at least one of: paper-band, display-style, accent-hue.'* I avoided repeating themes here by deliberate intuition; a documented rule would help."
- "**Copy generation bent toward 'Built for the modern team' twice — and I had to course-correct.** Tracejam's first draft headline was *'Trace what matters.'* — a near-template; I rewrote it to *'Distributed tracing that explains itself.'*"
- "**Free fonts only got me so far.** The bakery wanted Tiempos Headline (paid, Klim); I substituted Newsreader."
- "**Within-archetype variation knobs work in theory but I leaned on the same defaults.**"

**The findings → fix loop is the eval.** Improvement list is prioritised Tier 1 (5 items, "Real holes") / Tier 2 (5 items, "Quality polish") / Tier 3 (4 items, "Long bet"), then the *next section of the same file* records what shipped:
> "## v0.5.0 — Tier 1, 2, and most of Tier 3 — DONE. Items 1–11 and 13 … shipped in commit `b61f1ef`" with a ✅/🟡 checklist mapping each finding to the file it landed in (e.g. "✅ **#6 Different-knobs slop-test gate** (gate 34)", "✅ **#10 aria-label slop-test gate** (gate 35)", "✅ **#11 Project memory** `.hallmark/log.json`").
> Then "## v0.6.0 — Tracejam fixes + SaaS expansion + visual gallery" — **three new gates born from one page's regressions**: "Gate 36: no horizontal scroll on any viewport 320–1920 px. Gate 37: decorative effects on text must be visually verified to sit at x-height, not baseline. Gate 38: interactive bars must declare `align-items: center` + `line-height: 1`."

**This is the repo's core loop: a defect found in one generated page becomes a numbered permanent gate.** Gates 34–36, 50–57 all read as forensic post-mortems of specific bugs (e.g. gate 52 names the exact failure — "most visible on Sport: italic Anton title overlapping '02 / EXAMPLES'"; gate 55 names "SAME PROMPT, TWO / DIFFERENT OUTPUTS." as the observed comma+cap-D collision).

### 4d. The verb tests (`site/_tests/verbs/`)

`verbs/README.md`: "Each folder holds `input.html` (or `input-description.md` for `study`) · `output.*` · `notes.md` — one page explaining what fired, what reference files loaded, and what the durable artifact is. **The three verb tests together exercise every load path in the skill that the default-flow tests don't.**"

- **`audit/audit-report.md`** — the model output format: three severity bands (**CRITICAL "page ships as slop" · MAJOR "looks AI-generated" · MINOR "small taste issues"**), each a table of `# | Tell | Where (file · selector · line no.) | Fix`. Ends with `**5 critical · 6 major · 3 minor.** Total: 14.` and: "The page fails the slop test on **at least 14 of 38 gates.** It is not refinable in place." Then a **routing recommendation**: "Do **not** run `hallmark refine` here. Run `hallmark redesign` instead — `refine` preserves structure, and the structure itself is the primary problem."
  `audit/notes.md` states the discipline: "**Reads the file.** No edits. No new files. The verb is read-only… **`audit` is the verb that doesn't change anything. It's the safety verb**… That is the ideal output: **a verb that helps the user pick *another* verb**."
- **`redesign/notes.md`** — an explicit **preservation contract**: a bullet list of *every word preserved verbatim* (headline, sub-copy, button label, three feature headings, three feature bodies, CTA headline+sub, footer, brand name) and an 11-row `Axis | Input | Output | Why` replacement table. Plus the discipline note: "**Did not rewrite the copy.** The cliché phrases stayed. (That's `redesign`'s discipline — the user keeps the words.)" Scored `Before 24 / 38 ✓ → After 38 / 38 ✓`.
- **`study/`** — `input-description.md` (the screenshot rendered as prose, since the test artifact can't carry an image) → `diagnosis.md` (10 numbered fields + carry-over warnings + confirmation question + *"What the user answered"*) → `notes.md` (what loaded, what it did NOT do) → `output.html`/`output.css`.

### 4e. Self-reported scoring inside the artifact

Every emitted page stamps its own verdict in the CSS. Two generations exist:

- v0.6.0 stamp (tests 01–08) — 3–4 comment lines: macrostructure + hero knobs, theme + accent + enrichment, `studied: no · context: explicit|partial, inferred|skipped, inferred · v0.6.0`.
- v1.1 stamp (`site/examples/custom-0*`) — adds a **machine-greppable gate ledger**:
  ```css
  /* Hallmark · route: custom · structure: … · idea: "…" · paper: oklch(…) · accent: oklch(…)
   * display: Bodoni Moda · body: Spectral · gates: all-pass · v1.1 */
  /* Hallmark · pre-emit critique: P5 H5 E5 S5 R5 V5 */
  /* axes: … · contrast: pass (46–50) · nav: none (broadside) · footer: colophon (Ft-bespoke)
   * · slop: pass (51–55) · honest: pass (56) · chrome: pass (57) · tokens: pass (58)
   * · responsive: pass (59) · icons: pass (60) · mobile: pass (36, 59, 61–69) */
  ```
  (Those trailing gate numbers are from an *older numbering* than today's `slop-test.md` — another drift signal: the stamps are hand-written and go stale.)

Observed pre-emit scores across the corpus: `P5 H5 E5 S5 R5 V5` (garden-01, custom-01), `P5 H5 E5 S5 R4 V5` (riso-01, press-01, custom-02/03/04/05), `P5 H4 E4 S4 R5 V4` (09-slow-pour, 12-loafer), `P5 H4 E5 S5 R5 V4` (11-soroe). **Nothing in the corpus ever scores below 4** — i.e. the self-score is descriptive, not adversarial.

### 4f. `all-themes.html` — the only "harness"-shaped artifact

`site/_tests/all-themes.html` (132 lines, no build step): a static grid of 20 `<iframe src="../index.html?theme=KEY">` scaled to `0.359375` (460/1280) that renders the **same** landing page in all 20 themes side by side, each card tagged `Keep` / `Crowded / rework` / `Cut candidate` (legend: green/amber/red dots). Header comment: *"Internal review sheet… Not linked from the public site."* Sub-line: *"Tags are my recommendation, not a decision — yours to make."* Today all 20 are tagged `keep` (values: Keep ×16, Reworked ×2 (garden, bloom), New ×1 (cobalt), Default ×1 (hum)).

This is the closest thing to a **regression harness**: hold the content constant, vary the theme, eyeball 20 renders on one screen. It is manual and human-judged.

---

## 5. The theme catalog

### 5a. Current catalog — 20 named themes (authoritative list, `slop-test.md` gate 57 + `custom-theme.md` § "Two routes")

> "Hallmark's 20 themes (Specimen, Midnight, Brutal, Garden, Atelier, Newsprint, Terminal, Manifesto, Almanac, Sport, Studio, Riso, Bloom, Coral, Cobalt, Aurora, Editorial, Carnival, Lumen, Hum). **Each one is a fixed combination of paper-band, display-style, and accent-hue.** The rotation rule cycles through them so two consecutive runs don't read alike. **This is the default.**"

`site/css/tokens.css` header states the design constraint (verbatim):
> "Twenty-four themes. Each occupies a distinct point in OKLCH space. **Paper colours span dark / cream / saturated / cool / warm; never adjacent.** Display fonts span **category-level differences (serif / mono / condensed sans / italic); no theme is a minor variation of another.** Switched via `[data-theme="..."]` on `<html>`."
(The "Twenty-four" is stale — the file defines 20 + one `lumen[data-drop="day"]` variant.)

| # | Theme | `tokens.css` one-line identity | `--color-paper` | Genre (`main.js` THEME_GENRES) | Hero archetype (`main.js` ARCHETYPES) |
| --- | --- | --- | --- | --- | --- |
| 1 | **Specimen** | warm oat editorial workshop · serif default | `oklch(96% 0.018 80)` | editorial | marquee |
| 2 | **Midnight** | deep cool charcoal · technical | `oklch(15% 0.022 250)` | atmospheric | stat-led |
| 3 | **Brutal** | stark white · heavy condensed · bright red | `oklch(98% 0.001 0)` | editorial | marquee |
| 4 | **Garden** | modern botanical almanac (REWORKED — was flat sage + rose) | `oklch(95.5% 0.022 92)` warm oat cream, "never grey-sage" | editorial | letter |
| 5 | **Atelier** | luxury fashion-house | `oklch(94% 0.005 60)` | editorial | quote-led |
| 6 | **Newsprint** | broadsheet · salmon-pink paper · burgundy accent | `oklch(92% 0.045 50)` | editorial | split |
| 7 | **Terminal** | phosphor CRT · monospace everywhere | `oklch(11% 0.018 145)` | atmospheric | marquee |
| 8 | **Manifesto** | BLACK paper · ALL CAPS condensed · red colour-blocks | `oklch(10% 0.005 60)` | editorial | marquee |
| 9 | **Almanac** | encyclopaedic · **COOL** pale paper (not warm) · slate accent | `oklch(94% 0.008 245)` | editorial | stat-led |
| 10 | **Sport** | crisp white · italic display · burnt orange accent | `oklch(98% 0.003 250)` | editorial | stat-led |
| 11 | **Studio** | modern editorial agency · Fraunces italic display | `oklch(96.5% 0.005 200)` | editorial | letter |
| 12 | **Riso** | risograph print · Public Sans display, off-register | `oklch(91% 0.034 30)` | editorial | quote-led |
| 13 | **Bloom** | clean warm-light minimal · one restrained warm accent | `oklch(97% 0.010 72)` calm warm off-white | atmospheric | marquee |
| 14 | **Coral** | modern-minimal · warm-grey paper · single coral accent | `oklch(96.5% 0.005 50)` | modern-minimal | split |
| 15 | **Cobalt** | modern-minimal · cool engineered paper · electric cobalt | `oklch(98.5% 0.004 250)` | modern-minimal | cobalt (bespoke) |
| 16 | **Aurora** | atmospheric · cool blue-green gradient on near-black | `oklch(11% 0.025 200)` | atmospheric | marquee |
| 17 | **Editorial** | editorial-premium magazine voice | `oklch(94% 0.020 75)` "slightly cooler than Specimen" | editorial | split |
| 18 | **Carnival** | loud maximalist editorial · variable display | `oklch(92% 0.045 50)` warm pink-cream | editorial | marquee |
| 19 | **Lumen** | premium AI-tool · **apparatus, not orb** · three-register type | `oklch(13% 0.014 265)` late-night studio, violet tilt (+ `[data-drop="day"]` light variant `oklch(97% 0.008 265)`) | atmospheric | marquee |
| 20 | **Hum** | playful · vibrant · alive · multi-accent · rounded-sans | `oklch(97% 0.012 95)` cream, pear-yellow pull | **playful** (only member) | hum (bespoke) |

Per-theme *shape* tokens are a separate override block (`tokens.css:832` "PER-THEME COMPONENT SHAPE — radius, border weight, shadow"), e.g. `[data-theme="brutal"] { --rule-card: 2px }`, `[data-theme="riso"] { --radius-card: 2px }`, `[data-theme="bloom"] { --radius-card: 16px; --radius-pill: 999px }`, `[data-theme="editorial"] { --radius-card: 0; --rule-card: 0.5px }`.

### 5b. Four **genres** — a rule-set overlay above themes

`main.js` comment: "Each theme belongs to one of four genres — **a rule-set overlay that scopes which slop-test gates apply and which voice fixtures the skill picks from.**" Files: `references/genres/{editorial,atmospheric,modern-minimal,playful}.md`. Membership: editorial ×12, atmospheric ×5, modern-minimal ×2, playful ×1.

Genre overrides are written *inline into gate text*, e.g.:
- Gate 2 (gradients): "*atmospheric allows radial gradients on background only — never on text or pill buttons. **No genre allows gradient text.***"
- Gate 6 (centred hero): "*atmospheric and playful allow a centred hero when the canvas itself is the design (Suno-style)*"
- Gate 7 (`#fff`): "*modern-minimal allows pure `#fff` paper (the Stripe / ElevenLabs school)*"
- Gate 22 (zero-chroma neutrals): "*modern-minimal allows zero-chroma neutrals*"
- Gate 23 (accent ≤5% viewport): "*atmospheric allows accent-tinted radial blooms covering up to ~20 % of the canvas, since the bloom is the design*"
- Gate 21 (Specimen fall-through): "*atmospheric, modern-minimal, and playful never default to Specimen — only editorial does, and only when the brief signals it.*"

### 5c. Legacy / drifted theme names still live in the corpus

The 8 canonical tests use **Salon, Linen, Plain, Pastel, Quiet, Halo** — names not in today's 20-theme catalog (they came from a 16-theme era). Newer inline-CSS tests 09–13 use **Atelier, Aurora, Studio, Pastel, Halo**. The corpus was never re-generated after the rename. Practical lesson: **the corpus rots the moment the catalog changes, because nothing links them programmatically.**

### 5d. Custom route (the escape valve)

`references/custom-theme.md`. Two depths:
- **Tuned** — one-off OKLCH palette + font pairing *keeping* Hallmark's structures. "The freedom is the combination… but never the floor."
- **Bespoke** — designs "palette, type, **and** composition" from first principles. **Drops**: named-theme tokens, genre cluster routing, the macrostructure/archetype catalog, the diversification rotation. **Keeps** (the floor): "every universal slop-test gate", accessibility & contrast (APCA/WCAG), visible `:focus-visible`, `prefers-reduced-motion`, semantic landmarks, alt text, the **font ban-list**, OKLCH palette discipline, one orchestrated motion, the Step 5 preview before code, the Step 6 stamp + log.

Trigger discipline (verbatim): "Hallmark must **not** offer catalog-vs-custom on every prompt. That's friction, not discipline." Five signals only (explicit ask · named brand colour · 3+ off-catalog vibe words · brand-mood reference attached · a singular structural vision). "**One adjective ('warm', 'technical', 'playful') is not a signal** — that's a tone, the catalog already carries it." And: "**Default to catalog** — silence routes to catalog, not custom." Then exactly **one** follow-up question: *"describe the brand's vibe in 4–8 words"* + optional anchor colour. "**Do not ask anything else.**"

Palette construction is a deterministic recipe (`§ B`): accent first (**clamp chroma to 0.12–0.20**; hue derived from vibe — warmth 30–60°, technical/industrial 220–250°, botanical/moss 130–160°, late-night/neon 280–320°, sun-drenched 60–80°) → paper L band by vibe (bright 95–98 %, archival 92–95 %, technical 98–100 % near-white cool-tinted, dark 12–18 %) → ink (paper L<50 → ink L 88–96 %; else 16–24 %) → supporting greys stepped 6–10 % L, all tinted 0.005–0.018 chroma toward the anchor. **"Always tint paper toward the anchor hue with chroma 0.005–0.020."**

`site/_tests/custom/README.md` lists the five guards against over-invention: opt-in only · one question only · every existing rule applies · Step 5 preview surfaces the palette before code · **"Diversification is theme-route-blind"** (custom runs log their three axis values to `.hallmark/log.json` and rotate the same way catalog themes do).

---

## 6. Packaging & distribution

- **Install:** `npx skills add nutlope/hallmark` — "Re-run any time to update."
- **Manual install paths (README):** Claude Code `~/.claude/skills/hallmark/` · **Cursor `.cursor/rules/hallmark.mdc` (body of `SKILL.md`, no frontmatter)** · Codex `~/.codex/skills/hallmark/` (personal) or `.codex/skills/hallmark/` (project-scoped).
- **`package.json`** — the whole distribution contract:
  ```json
  "files": ["skills"],
  "skill": {
    "entry": "skills/hallmark/SKILL.md",
    "references": "skills/hallmark/references",
    "harnesses": ["claude-code", "cursor", "codex"]
  },
  "scripts": { "serve": "python3 -m http.server --directory site 4173" }
  ```
  **No build step, no tests, no lint, no dependencies.** `type: "module"` but zero JS is shipped in the package.
- **SKILL.md frontmatter** (3 fields only):
  ```yaml
  name: hallmark
  description: "Anti-AI-slop design skill for greenfield pages, audits, redesigns, and design extraction from URLs or screenshots. Use when the user asks to build a new app or landing page, wants to redesign something, invokes Hallmark by name, or uses audit/redesign/study."
  version: 1.1.0
  ```
  Version is duplicated in `package.json` (manual sync).
- **`vercel.json`** — `{"outputDirectory": "site", "framework": null, "buildCommand": null, "installCommand": null}`. The marketing site *is* the repo, statically served.
- **`.gitignore` notables** — `.hallmark/` ("Runtime state from Hallmark's own runs (per-project log, **never publish**)"), `skills-lock.json` ("generated whenever 'skills add' runs in this directory, not part of the skill source"), `.image-cache/`, `.claude/`, and a hand-curated allow/deny list of `site/examples/*` scratch builds with a `!site/examples/hum-07/` re-include.
- **Progressive-disclosure loading discipline** (SKILL.md:381, the only line quoted from SKILL.md because it is a packaging rule): "`slop-test.md` — **strictly Step 7, after Build.** The 58 gates are a post-emit check, not a pre-emit reference. **Pre-loading slop-test.md costs ~7K tokens for nothing** — the gates inform fixes, not generation. If a gate fails at Step 7, fix and re-test; do not consult the file earlier 'to know what to avoid' — that's what `anti-patterns.md` is for."
- **`preview-examples.md` opening line** is the same doctrine: "Load this file only when picking an unusual macrostructure or custom theme and the bullet-list spec in `SKILL.md § 5. Preview` doesn't give enough scaffolding on its own. **Most builds don't need to read this file.**"
- **Portable handoff artifacts** (`design-md.md` + `export-formats.md`): `tokens.css` is "the source of truth… **Always emitted alongside the page CSS**"; three translations live inline in `design.md` — Tailwind v4 `@theme`, DTCG `tokens.json` (Style Dictionary / Token Studio / Cobalt), shadcn/ui CSS variables. "**No new verb.** This is a side effect of every build." `design.md` targets "~45 lines — enough to seed a real app, not so much that it becomes a wiki to maintain", and fires **phrase-only** (`"lock the system"`, `"give me a design.md"`, `"extract the DNA"`, …) — never automatically.

---

## (a) DOCTRINE — design rules, taste, knowledge

Numbered/quantified rules worth stealing verbatim:

**Colour**
- Accent ≤ **~5 % of any single viewport** by area (gate 23); atmospheric genre may go to ~20 % *only* for radial blooms.
- Accent chroma clamp **0.12–0.20**.
- **No zero-chroma neutrals** — "minimum 0.005 chroma", tint every neutral toward the anchor hue (gate 22).
- No pure `#000` / `#fff` as a base colour (gate 7).
- Paper L bands: light 92–98 %, near-white 98–100 %, dark 12–18 %. Ink: 88–96 % on dark paper, 16–24 % on light.
- Contrast: body (<24 px regular / <18 px bold) **WCAG 4.5:1 or APCA Lc ≥ 60**; large text + icons + focus rings **3:1 / Lc ≥ 45**. Fast pre-check: `|L_text − L_bg| < 50 %` → likely fail.
- **Button text ≈ button fill fails if within 5 % lightness AND 0.05 chroma** (the black-on-black bug).
- Any section with background OKLCH L < 50 % **must swap its text colour in the same rule**.

**Typography**
- **Three faces is the ceiling** (`--font-display`, `--font-body`, ≤1 outlier). Gate 37.
- The outlier face may occupy **at most two slots** (gate 38) — "The outlier is a register, not a third surface."
- **Gate 38a: no italic headings, at all** — "Italic headers — above all the single italicised emphasis-word inside an upright headline — are a top AI tell. Headers are roman; emphasis comes from weight, accent colour, or a drawn underline."
- Body-copy floor **14 px**; prose `max-width` **45–75 ch**.
- All-caps display heads: **`line-height` floor 1.0, recommended 1.02–1.08** (gate 55, cap-collision on wrap).
- Banned display fonts (gate 1): Inter, Roboto, Open Sans, Poppins, Lato, system defaults.

**Motion**
- Hover **150–200 ms**; focus rings **instant, never fade in** (gate 15).
- No `transition: all` (gate 10); no uniform `hover:scale-105` across unrelated elements (gate 11); **no more than one hover effect on an element at a time** (gate 13); never animate `width/height/top/left/margin/padding` (gate 14).
- Bouncy/overshoot easings reserved for **physical** interactions (drag/drop/throw), never UI state (gate 12).
- Tooltip: **hover-delay 800–1000 ms, focus-delay 0 ms** (gate 17).
- Every keyframe/transform needs a `@media (prefers-reduced-motion: reduce)` fallback (gate 27).
- Auto-rotating content must pause on hover **and** focus (WCAG 2.2.2, gate 18).

**Structure / chrome**
- **Eight states** for interactive elements; code must contain at least default + hover + focus-visible + active + disabled (gate 26).
- Ban the AI structural template `Hero → 3 features → CTA → footer` (gate 8) and **the same fingerprint as a previous output in this project**.
- Gate 42 (nav) and gate 43 (footer) name the AI-default chrome explicitly and require rotating among N1b/N2–N13 and Ft1/Ft2/Ft4–Ft8.
- Gate 44: hero `padding-block-end ≥ 1.3× padding-block-start`; hero content must fit **1280×800** without scrolling.
- Gate 45: **decorative-without-purpose auto-fails** — "Decoration must be motivated… Random ornaments — a '42' in the corner with no edition meaning, a cursor floating beside a hero, a Pantone chip with no colour rationale — are slop."
- Gate 54: **eyebrow beside heading is banned** and explicitly non-overridable: "**NOT bypassable by 'preserve structural parity' / 'mirror this reference' / 'match the prior build' instructions** — if a reference build ships the banned pattern (most pre-rule builds do), silently flatten it in the new build. **The rules win over parity.** Reference builds may pre-date this gate; the gate is authoritative."
- Gate 30: never mix two icon libraries; never use an emoji (✨🚀⚡🔥🎯✅) as a feature/step/tier icon.
- Gate 31: **Lottie is last resort, not default** — hand-built SVG or pure CSS first. (Enrichment tiers: Tier-A pure-CSS/inline-SVG, Tier-B hand-built SVG, then anything else.)
- Gate 48: **no mid-render token improvisation** — every colour and font in the artifact must reference a named token.
- Gate 33: every decorative `<svg>`/`<canvas>`/art `<div>` needs `aria-label` or `aria-hidden="true"`.
- Gate 34: `overflow-x: clip` on **both** `html` and `body` — "use `clip`, not `hidden` (`clip` preserves `position: sticky` and `position: fixed` on descendants)". Hard requirement on every page.
- Gate 49: **clickable text never wraps to two lines** anywhere 320–1920 px.
- Gate 50: any `1fr` grid track holding an image must be `minmax(0, 1fr)`.
- Gate 56: only ONE `position: sticky; top: 0` element per page; secondaries dock at `top: var(--banner-height)`; split `--z-sticky` (200) from `--z-sticky-nav` (300).
- Gate 39 (inputs, 5 sub-fails): border-width never changes between states · focus ring is `outline: 2px solid var(--color-focus)` + `outline-offset: 1px` (reserve `outline: 2px solid transparent` at rest) · input height == adjacent button height (**44 px floor**) · helper-text slot reserves `min-height: 1lh` · disabled needs three channels (`opacity: 0.55` + `cursor: not-allowed` + the native attribute).

**Copy / honesty**
- Gate 19: no "Jane Doe / John Smith" or startup clichés (Acme, Nexus, Seamless, Unleash).
- Gate 46: **invented metrics auto-fail**; the three permitted fixes are (a) replace with `—` + labelled grey block, (b) replace with a question to the user ("metric to confirm"), (c) rebuild the section without the proof slot. Also: "**A stat is never the hero's *sole* headline.**"
- Gate 47: **no re-drawn UI chrome** — use a real screenshot in `<picture>`/`<figure>` or omit.

---

## (b) PROCESS — workflow / phases / state

**The 8-step default flow** (from SKILL.md headings + every `brief.md`):
`0. Pre-flight scan → 1. Design-context gate → 2. Pick a macrostructure FIRST → 2.5 Check project memory → 2.6 Theme route (studied-DNA | catalog | custom) → 3. Load the visual ruleset → 4. Decide on hero enrichment → 5. Preview → 6. Build → 7. The slop test`

Key process inventions:

1. **Macrostructure FIRST, dress second.** Shape is chosen before palette/type. Anti-pattern otherwise: "Specimen fall-through" (gate 21). Vague briefs must pick from the *first ten* of the 21 macrostructures — "deliberately the strongest non-Specimen shapes; they cover ~80% of briefs."
2. **The design-context gate** — three questions only (**audience · use · tone**), and the user is allowed to skip. On skip, the model must **state its inference in one sentence and invite redirect**: *"Going with: audience = … · use = … · tone = …. If any of those is wrong, redirect."*
3. **Pre-emit self-critique before the gate sweep** (`slop-test.md`, verbatim): "Run this **before** the gate list, not after. Score the planned output 1–5 on each axis. Anything **< 3 on any axis triggers a revision pass** before the gate sweep — don't bring known weakness into a fifty-eight-gate review. **Two passes is normal. Three is a sign the brief is wrong, not the design — re-read the brief.**"
   | Axis | What you're scoring |
   | --- | --- |
   | **Philosophy** | Is there a clear *why* — a position the page is taking? Or is it just a layout? |
   | **Hierarchy** | Can a reader tell, in 2 seconds, what's primary, secondary, tertiary? |
   | **Execution** | Are the details (rule weight, accent footprint, text-wrap, focus rings, contrast) in spec? |
   | **Specificity** | Does this look like *this brief* — or a generic "page that could be anyone"? |
   | **Restraint** | Have you removed everything that isn't earning its place? |
   | **Variety** | Does this share a structural fingerprint with a previous output? **"Score by structural distance, not visual distance — colour-swaps don't count as variety."** |
4. **Step 5 Preview — a plan emitted as text before any code.** Fixed 8-row block: Macrostructure · Theme · Enrichment · Sections · Motion · Slop test (`58 / 58 ✓`) · Diversification. This is the user's early-redirect window and is required even on the custom/bespoke route.
5. **Step 7 slop test is post-emit only.** Gates are a *check*, not a prompt (see § 6, the ~7K-token argument).
6. **Verb routing is an explicit output.** `audit` must recommend which verb to run next; `redesign` publishes a preservation contract before it rebuilds.
7. **Diff-against-previous-run is a first-class artifact.** Every `brief.md` ends with `What changed vs v1` / `What stayed the same`, per axis.

---

## (c) MACHINERY — scripts, detectors, enforcement, tooling

**There are no scripts.** Zero executables, zero tests, zero CI, zero dependencies. All enforcement is textual + convention. The machinery is:

1. **`.hallmark/log.json`** — per-project append-only memory, gitignored, "never publish". Real schema (from `site/examples/bananastudio/.hallmark/log.json`):
   ```json
   [{ "date": "2026-05-04",
      "macrostructure": "Marquee Hero",
      "theme": "Bloom",
      "enrichment": "Tier-A animated aperture + breathing portraits + marquee strip",
      "brief": "BananaStudio — redesign to dark mode, bold, smooth interactions, animated illustrations",
      "knobs": { "H7_marquee_hero": "title=display-italic-mix, alignment=left-bias, counterweight=animated-aperture",
                 "F2_gallery": "tiles=7, layout=feature+6-portraits, breathing-loop=on",
                 "F4_steps": "count=3, container=workbench-cards, divider=hairline-gradient" },
      "theme_axes": "dark / italic-serif / saturated-amber-bloom" }]
   ```
   Optional keys seen elsewhere: `polish`, `nav`, `footer`, `genre`, `vibe`. Newest entry first (wayfare has v2 then v1).
2. **The CSS macrostructure stamp** — a machine-greppable header comment in every emitted stylesheet. Gate 20 fails if it is missing. It is simultaneously (a) provenance, (b) the fallback memory when no `log.json` exists ("Read the file system: if a `.hallmark/log.json` entry **or** a CSS macrostructure stamp exists, this build's macrostructure must differ from the last"), and (c) a **behaviour switch for future runs**: `studied: yes` makes `audit` lenient on Specimen-fall-through but stricter on "did the page actually use the extracted DNA?"; `context: redesign` does the same for a deliberate restructure.
3. **The 3-axis diversification rule** — themes are reduced to `paper-band / display-style / accent-hue`; consecutive runs must differ on **at least one**. Stated in a `brief.md` as a pass/fail judgement: *"Two of three axes differ. Passes."* / *"One axis differs — passes (the rule requires at least one)."* Overridden when the user names a theme.
4. **Within-archetype knobs** (gate 32) — "Two Bento Grids with `tiles=6, spans=irregular, accent=corner-only` are the same Bento… **State the knob deltas in the stamp.**"
5. **The gate ledger in the stamp** — `· contrast: pass (40–41) · nav: N# · footer: Ft# · slop: pass (42–45) · honest: pass (46) · chrome: pass (47) · tokens: pass (48) · responsive: pass (49) · icons: pass (30) · mobile: pass (34, 49, 50–57)`. Each gate cluster names the stamp field it writes.
6. **Progressive disclosure by file** — 24 refs + 50 component files + 21 macrostructure files + 4 genre files, each loaded on demand with an explicit "load this only when…" preamble.
7. **`all-themes.html`** — the manual 20-up regression sheet (see § 4f).
8. **`site/js/main.js` as executable spec** — `THEMES` registry (20), `ARCHETYPES` (theme → `{hero, footer}` tuple, comment: "**switching themes literally rebuilds the page, not just recolours it**"), `THEME_GENRES` (theme → one of 4 genres), a **locked `HERO_TITLE`** ("A design skill that refuses to look AI-generated.") with per-theme copy fixtures for everything else. `STORAGE_KEY = "hallmark-theme"`; `?theme=KEY` querystring; press `T` to cycle.
9. **The token contract** (`export-formats.md`) — a fixed taxonomy (`--color-paper/-2/-3`, `--color-ink/-2`, `--color-rule/-2`, `--color-muted`, `--color-neutral`, `--color-accent`, `--color-accent-ink`, `--color-focus`, `--font-display/-body/-outlier`, `--space-3xs…5xl` on a 4-pt scale, `--text-xs…display` on a 1.25 major-third ratio, `--ease-*`, `--dur-micro/short/long` = `120ms/220ms/420ms`, `--rule-hair/fine` = `1px/2px`, `--radius-card/pill/input`, `--shadow-card`) + 4 export formats. "Don't make up token names downstream that aren't in `tokens.css` — the source of truth is the source of truth."

---

## (d) SINGLE-AGENT ASSUMPTIONS that would NOT transfer to a multi-agent protocol

1. **One conversation = one memory.** Gate 57 keys off "a `study` diagnosis emit**ted earlier in the conversation**"; `design-md.md` fires on the user having said a trigger phrase. In Parley there is no shared conversation — five agents each have their own. Any "earlier in this session" condition must become a *file*.
2. **The stamp is the only cross-run channel, and it's written by the same agent that reads it.** With N agents writing N artifacts, a single `style.css` stamp cannot be the coordination point; and `.hallmark/log.json` is a **single-writer append-only array with no ordering, no author field, and no conflict handling**. Concurrent agents would clobber it.
3. **Interactive mid-flow questions.** Step 1's gate, `custom-theme.md` § A's one follow-up, and `study`'s confirmation question all **block on a human reply**. Headless agents (`codex`, `hermes`, `agy`, `kimi` in `--print` mode) cannot ask. Every ask-branch must be pre-resolved in the idea's brief or converted into a declared assumption in the agent's artifact.
4. **Self-scored gates with no adversary.** The corpus shows no pre-emit score below 4 and every page self-reporting `38/38 ✓` or `58/58 ✓`. A single agent grading its own output is exactly what Parley's cross-review replaces — but the *rubric* is what transfers, not the self-grading.
5. **Diversification-against-self.** "Structural variety" is defined as *differing from my own previous output*. Multi-agent inverts the problem: N agents independently pick from the same catalog and will **converge**, not diverge (they share a training prior). A multi-agent variety rule must be **claim-based** (agent A claims Bento, agent B must pick differently) or **convergence must be reframed as consensus signal**, not as a failure.
6. **Visual verification is imagined, not rendered.** Gate 35: "The check is *visual*: **imagine the rendered output** and confirm the band lands in the right vertical zone." Gate 34 says "Open the rendered page; drag the dev-tools width slider." Nothing in the repo actually renders. The ROADMAP admits this ("Live preview as an MCP server… Closes the loop between generation and audit"). In Parley this is a real opportunity — one agent can render/screenshot while another reviews.
7. **Single artifact per run.** `index.html` + `style.css` + `tokens.css`, one page, one stamp. Parley's protocol is N artifacts under `parley-deck/ideas/<slug>/` that must be *reconciled*; Hallmark has no reconciliation concept at all.
8. **Marketing counts drift because nothing verifies them.** README "57 gates" vs `slop-test.md` "58 gates"; `tokens.css` "Twenty-four themes" vs 20 defined; `recipes.md` "16 themes → 2,688 fingerprints" vs 20 → 3,360; `verbs/README.md` references a `refine/` folder that does not exist; the v1.1 stamps carry gate numbers from an obsolete numbering. **A multi-agent protocol with a drift guard (Parley already has `TestEmbeddedDefaultMatchesLiveDeck`) should mechanically enforce these.**
9. **Corpus is manually re-generated and rots.** Tests 01–08 still name themes (Salon, Linen, Plain, Pastel) deleted from the catalog. Nothing links a `brief.md` to the catalog it was generated against.

---

## Transferable to parley-design

Ranked by expected value for a multi-agent design protocol.

**1. A numbered, permanent gate list as the shared objective function.** 58 gates, each phrased as a **yes/no question where "yes" is a failure** ("Run this list before handing back any output. **Every answer must be no.**"). This is the single highest-leverage artifact — it gives N heterogeneous agents one comparable scoring surface, converts subjective taste disputes into citable gate numbers, and makes cross-review mechanical: agent B's BLOCK is "gate 44(b) fails at 1280×800", not "the hero feels off". Adopt the numbering discipline (gates are append-only; a gate number never changes meaning) and the `38a` precedent for inserting without renumbering.

**2. The "defect → permanent gate" ratchet.** Every regression found in a generated page became a numbered gate forever (gates 34–36 from Tracejam; 50–57 from one responsiveness pass). Wire this into Parley's Phase-7/8 fix-up: any design defect a reviewer catches must land as a new numbered gate in the design ruleset, not just as a fix to one file. This is the retrospective-optimization loop (§13 RHO) applied to taste.

**3. Pre-emit self-critique on six named axes with a hard `< 3 → revision` threshold.** Philosophy · Hierarchy · Execution · Specificity · Restraint · Variety, 1–5, stamped as `P5 H4 E5 S4 R5 V5`. In Parley this becomes each agent's **self-declared score in its own artifact**, which cross-review then *contests* — solving Hallmark's "nothing ever scores below 4" problem for free. Keep the note "Two passes is normal. Three is a sign the brief is wrong, not the design."

**4. The machine-readable provenance stamp in the emitted artifact.** A header comment carrying macrostructure + knobs + theme + accent + enrichment + gate ledger + `studied:` / `context:` flags. Two properties matter: (a) it is **greppable state that survives the session**, (b) **flags in it change how a later reviewer grades the file** (`studied: yes` → lenient on fall-through, strict on DNA fidelity). Parley equivalent: stamp the FINAL.md-ratified design decisions into every emitted stylesheet/component so the code reviewer can check the code against the consensus without re-reading the deliberation.

**5. Plan-before-code: the Step 5 Preview block.** A fixed 7–8 row plan (shape · theme · enrichment · sections · motion · gate result · diversification) emitted **as text before a single line of code**, on every route including bespoke. Maps directly onto Parley: this *is* each agent's round-01 artifact. Standardise the row set so five agents' previews are diffable line-for-line. That diffability is what makes multi-agent design consensus tractable.

**6. Structure-first, dress-second, with a named catalog of shapes.** 21 macrostructures × 20 themes × 8 enrichment tiers, and the rule that vague briefs pick from *the first ten* shapes. For Parley: give agents a **finite named vocabulary** so five artifacts can be compared on the same axes instead of five prose essays. Naming is the interop layer.

**7. The 3-axis reduction (`paper-band / display-style / accent-hue`) as a similarity metric.** Any theme collapses to three comparable values. This is the cheapest possible "are these two designs actually different?" test and it is *computable from the stamp*. In multi-agent: use it to detect **convergence** across agents (all five picked `light / geometric-sans / cool` → the roster has a shared prior, flag it) rather than to force rotation.

**8. `audit` as a read-only verb whose best output is a routing decision.** "**`audit` is the verb that doesn't change anything. It's the safety verb**… That is the ideal output: **a verb that helps the user pick another verb**." Severity bands (CRITICAL / MAJOR / MINOR), a `# | Tell | Where (file:line) | Fix` table, a count line, and an explicit "do not run X, run Y". This is exactly the shape a Parley cross-review artifact should take for design work — and it maps onto Parley's existing dispositions.

**9. The `redesign` preservation contract.** Before rebuilding, publish the exhaustive list of what is preserved *verbatim* (every headline, label, body string) and an `Axis | Input | Output | Why` table for everything replaced. Explicitly out of scope: "Did not rewrite the copy… That's `redesign`'s discipline." In a multi-agent setting this is the anti-scope-creep device — it is checkable by another agent without taste judgement.

**10. `study`'s refuse-or-proceed + "anti-patterns to NOT carry over".** Three parts: (i) a refusal check before extraction (paid templates, copy-protected work, ambiguous ownership → ask), (ii) **name the *role*, not the ID** ("italic editorial serif at high optical size", two candidates — one free, one paid — never a font identification), (iii) a mandatory list of things in the reference that are *wrong* and must be corrected on rebuild, with numbers (800 ms → 150–200 ms; 11 px → 12–13 px). The "one thing it does better than most and you SHOULD carry over" inverse is equally valuable. This is a ready-made protocol for "here's a design I like, build like it" in Parley.

**11. Genre as a rule-set *overlay* with inline gate exceptions.** Four genres (editorial / atmospheric / modern-minimal / playful) scope which gates apply, written **inline into each gate's text** ("*Genre note: modern-minimal allows pure `#fff` paper*"). Cheaper and less error-prone than forking rulesets. Parley: one gate list, per-project overlays declared in the idea's FINAL.md.

**12. Progressive-disclosure loading with an explicit token argument.** "Pre-loading `slop-test.md` costs ~7K tokens for nothing — the gates inform fixes, not generation." Every reference file opens with "load this only when…". With 5 agents × N rounds, load discipline is a real cost line; state it as a rule, per-file, with the reason.

**13. Constraints-beat-instructions, i.e. the ban list.** Slide 10: "**Counterintuitive but true: telling the model what NOT to do is stronger than describing what you want.**" Ship a short, memorable ban list (8 items fit on a slide) *separately* from the long gate list, and put it where every agent reads it in round 0.

**14. A one-line portable design contract (`design.md`, ~45 lines).** "enough to seed a real app, not so much that it becomes a wiki to maintain", with a **no-overwrite policy** (refresh the Exports section only) and **phrase-only triggers** (never auto-emit). In Parley this is the natural bridge: FINAL.md ratifies the system → `design.md` is the durable, vendor-neutral handoff every implementer agent reads.

**15. The gate-54 "rules win over parity" clause.** "**NOT bypassable by 'preserve structural parity' / 'mirror this reference' / 'match the prior build' instructions**… The rules win over parity. Reference builds may pre-date this gate; the gate is authoritative." A multi-agent protocol needs exactly this to stop agent B from justifying a defect with "agent A's artifact does it too". Steal the wording.

**16. Cheap eval: hold content constant, vary the variable, render N-up.** `all-themes.html` — 20 iframes of the *same* page at `scale(0.359)`. Zero build. Parley variant: render each agent's proposal for the same brief side by side; the human sees the disagreement in one screen.

**17. `brief.md` as a reasoning trace with a per-axis diff against the previous run.** "What changed vs v1 / What stayed the same" is a better artifact than a changelog because it is *axis-indexed*. Parley already has per-agent artifacts; add the axis-diff so round N+1 is comparable to round N mechanically.

**18. `--color-accent-ink` as a required paired token.** "whenever `--color-accent` fills a surface that carries text, `--color-accent-ink` must be defined, verified ≥ APCA Lc 60 / WCAG 4.5:1 against `--color-accent`, and applied as the `color` on that fill." Small, concrete, catches the single most common shipped contrast bug.

---

## Do NOT copy

**1. The marketing counts, and the practice of hand-maintaining them.** "57 gates" (README) vs "58 gates" (`slop-test.md`) vs "57-gate" (`main.js`); "Twenty-four themes" (`tokens.css` header) vs 20 defined; "16 themes × … = 2,688 fingerprints" (`recipes.md`) vs 20 → 3,360; v1.1 stamps citing gate numbers 40–69 from a dead numbering. Parley must derive counts from the file or omit them. (Parley already has the pattern: a Go drift guard.)

**2. Self-reported gate scores as the quality signal.** Every artifact in the corpus claims `38/38 ✓` or `58/58 ✓` and no pre-emit axis ever scores below 4. A self-graded pass is a claim, not evidence. In Parley, a gate result is only meaningful when **a different agent** verified it, or a tool did.

**3. "Imagine the rendered output" as a verification method.** Gate 35's "The check is *visual*: imagine the rendered output"; gate 34's "Open the rendered page; drag the dev-tools width slider"; gate 44's "test at 1280×800". None of this happens. Either wire a real renderer/screenshot step (Parley can — one agent renders, one reviews) or **delete the gate**; a gate that cannot be executed is worse than no gate because it manufactures false passes.

**4. Blocking mid-flow questions to a human.** The design-context gate, custom's one follow-up, and `study`'s confirmation all halt on a human reply. Headless roster members cannot answer. Convert every ask into: resolve from the brief → else declare an explicit assumption in the artifact → else raise it as a consult under Parley's existing consult mechanism.

**5. `.hallmark/log.json` as-is.** Single-writer, append-only, no author field, no schema version, no ordering guarantee, no conflict handling, and gitignored ("never publish"). Five concurrent agents will clobber it. If Parley wants design memory, it must be per-agent files that a reconciliation step merges — same as the artifact model already in the protocol.

**6. "Diversify against my own last run" as the variety rule.** Correct for a single agent serving many briefs; **wrong** for N agents serving one brief, where convergence is signal, not slop. Also wrong for multi-page products — the ROADMAP already concedes it: "**The structural-variety rule is correct for variety, wrong for brand consistency inside a multi-page product.**" Do not import forced rotation; import the 3-axis *similarity metric* and decide per-context whether similarity is good or bad.

**7. A 20-entry theme catalog with only 4 documented themes.** `references/themes/` holds `carnival.md`, `cobalt.md`, `hum.md`, `lumen.md` — the other 16 exist only as token blocks in a site CSS file the skill does not ship (`site/` is excluded from the npm `files` array). The skill therefore *names* 20 themes it cannot fully define to the agent. Don't ship a catalog larger than the documentation behind it — with N agents, every undocumented entry is N different interpretations.

**8. The documented-but-missing `refine` verb.** `verbs/README.md` promises a worked example; `audit-report.md` routes users to it; `references/verbs/` has no `refine.md` and no `verbs/refine/` folder exists. Never let a doc reference a capability that isn't there — in a multi-agent run an agent will confidently invoke it.

**9. The stale corpus.** Tests 01–08 name themes (Salon, Linen, Plain, Pastel, Quiet, Halo) that no longer exist. Examples are the strongest teaching device Hallmark has and they now teach a dead vocabulary. Either regenerate the corpus on every catalog change (needs automation Hallmark doesn't have) or make examples reference the catalog by ID with a link-check.

**10. Vendor-specific / paid-service coupling.** ROADMAP "Now" is a Nanobanana API hook; slides recommend shadcn/ui as fix #1 ("Adopt it and your agent's output gets 80% better on day one"); `docs/talk-slides.md` and the site are Together-AI-branded. Parley is vendor-neutral by charter — take the *principle* ("pin a component system, don't invent primitives") and leave the vendor.

**11. Prescribing a house aesthetic wholesale.** Hallmark's taste is warm-paper editorial + OKLCH + hairline rules + Fraunces/Geist. That is *a* taste, and gate 38a (no italic headings, ever) is defensible as an anti-tell but is also a strong stylistic commitment presented as a universal law. `contract.md` claims "It will not… Enforce a specific style" while the gate list does exactly that. parley-design should separate **anti-slop invariants** (contrast, states, motion timing, honest copy, no re-drawn chrome) from **house-style preferences** (paper warmth, italic bans, hairline dividers) and let the deck's own FINAL.md ratify the latter.

**12. Prose-only enforcement with zero tooling.** No tests, no CI, no linter, no dependencies, no build. Admirable minimalism, but it means *nothing* in the repo is verified. Parley already runs `RunChecks` in the driver — design gates that can be expressed as a check (token-only colours, `overflow-x: clip` present, focus-visible on every interactive selector, no `transition: all`, no banned font families, three-family ceiling, `minmax(0, 1fr)` on image tracks) should be **executable**, and the rest should be explicitly labelled "human/agent judgement".
