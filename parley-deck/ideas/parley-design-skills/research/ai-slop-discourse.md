# AI Slop in UI/Visual Design — Empirical Definition

**Digest for:** `parley-design` (doctrine skill) + `parley-design-check` (enforcement skill)
**Compiled:** 2026-07-28
**Scope:** what "AI slop" empirically *is* in UI, how it is detected today, why it happens, where it is wrong to fight it.

---

## 0. Source register

Every claim below is tagged with a source ID. `[MINE]` marks my own analysis, not a source claim.

| ID | Source | Type | Why it counts |
|----|--------|------|---------------|
| **S1** | https://prg.sh/ramblings/Why-Your-AI-Keeps-Building-the-Same-Purple-Gradient-Website | Essay | Best statement of the causal mechanism |
| **S2** | https://www.925studios.co/blog/ai-slop-web-design-guide | Agency guide | Named tell list + copy clichés |
| **S3** | https://impeccable.style/slop/ | **Registry, 64 patterns** | Closest existing thing to a machine-checkable slop registry |
| **S4** | https://www.adriankrebs.ch/blog/design-slop/ + https://github.com/AdrianKrebs/design-slop-cop | **Empirical study, N=1,590** | Only quantitative study found; open-source detector |
| **S5** | https://www.developersdigest.tech/blog/ai-design-slop-and-how-to-spot-it | Article | Restates S4's 16 patterns with quotes |
| **S6** | https://dev.to/alanwest/why-every-ai-built-website-looks-the-same-blame-tailwinds-indigo-500-3h2p + https://dev.to/alanwest/how-to-fix-the-ai-generated-look-in-your-frontend-1ahh | Articles | Origin story + an actual ESLint ban-list |
| **S7** | https://news.ycombinator.com/item?id=48504912 | HN thread | Practitioner counter-arguments |
| **S8** | `/Users/tomasfecko/.claude/skills/hallmark/references/slop-test.md` (58 gates) and `/Users/tomasfecko/.claude/skills/hallmark/references/anti-patterns.md` | **Local prior art** | Most rigorous gate list available; already on this machine |
| **S9** | https://github.com/pbakaus/impeccable (`npx @impeccable/detect`) | **Shipping tool** | 60 deterministic rules, real CLI contract, Apache-2.0 |
| **S10** | https://axe-web.com/insights/ai-website-design-sameness/ | Article | Explicit "when it does NOT matter" |
| **S11** | https://arxiv.org/abs/2510.01171 (Verbalized Sampling; Zhang, Yu, Chong, Sicilia, Tomz, Manning, Shi) | Paper | Mechanism + the empirical case for DIVERGE |
| **S12** | https://en.wikipedia.org/wiki/AI_slop ; https://simonwillison.net/2024/May/8/slop/ | Etymology | Definitional discipline |
| **S13** | https://arxiv.org/pdf/2510.22954 (*Artificial Hivemind*), https://arxiv.org/pdf/2604.16027 (*Where does output diversity collapse in post-training?*) | Papers | Corroborating mode-collapse literature |
| **S14** | https://arxiv.org/pdf/2512.03373 (*LLM-Generated Ads*) | Paper | Quantified detection penalty (adjacent domain) |

---

## 1. Definitional discipline — what "slop" actually means

**FACT (S12).** Simon Willison's May 2024 post defined slop as *"AI-generated content that is both unrequested and unreviewed"*, and argued the term should function like "spam". His qualifier is load-bearing: *"Not all AI-generated content is slop. But if it's mindlessly generated and thrust upon someone who didn't ask for it, slop is the perfect term."*

**FACT (S12).** "Slop" was named 2025 Word of the Year by Merriam-Webster and the American Dialect Society; Macquarie Dictionary named "AI slop" its 2025 WOTY (both Committee's and People's Choice).

**FACT (S3).** The most mature registry explicitly splits its 64 patterns into two rule classes: **AI-slop rules** ("flag tells of AI-generated interfaces") and **quality rules** ("flag general design mistakes regardless of origin").

**[MINE] — the definition `parley-design` should adopt.** Three orthogonal failure axes are being conflated in the discourse. Keep them separate or the enforcement skill will be incoherent:

- **Axis A — Provenance tell.** The artifact *reads as machine-generated*. Cost: credibility/brand. (purple CTA, side-stripe card, gradient headline)
- **Axis B — Absence of position.** The artifact has no point of view; it would fit any brief. Cost: it is forgettable and interchangeable. (Hero → 3 cards → CTA → footer)
- **Axis C — Defect.** The artifact is simply broken or inaccessible. Cost: usability. (contrast failure, horizontal scroll, uncaught script error)

Axis C is *not* slop — it is a bug that correlates with slop. Axis B is the real target of `parley-design`. Axis A is the cheap machine-checkable proxy that `parley-design-check` can actually enforce. **Do not let the proxy become the goal** (see §6).

---

## 2. The base rate — the only quantitative study found

**FACT (S4).** Adrian Krebs analysed **1,590 Show HN submissions** with Playwright: *"A headless browser loads each site (Playwright)… A small in-page script analyzes the DOM and reads computed styles."* He explicitly rejected vision: *"I intentionally do not take screenshots and let the LLM judge them"* — instead every rule is a *"deterministic CSS or DOM check"*. Self-reported false-positive rate: **"maybe 5–10%"**.

**FACT (S4).** Distribution:

| Tier | Criterion | Share | Count |
|------|-----------|-------|-------|
| High | 4+ patterns | **22%** | 347 |
| Medium | 2–3 patterns | **32%** | 508 |
| Low | 0–1 patterns | **46%** | 735 |

**FACT (S4, repo).** Scoring formula: `score = round(100 × patternsFlagged / patternsTotal)`; tiers `5+ = High | 3–4 = Medium | 1–2 = Low | 0 = None`.

**[MINE] — three consequences for `parley-design-check`:**
1. **A single tell is not evidence.** 46% of a general sample trip 0–1 patterns; tripping one puts you inside the ordinary population. Gate on *count*, not on any individual rule.
2. **A 5–10% FP rate forbids hard-blocking on a single probabilistic rule.** Deterministic-but-heuristic ≠ correct. Severity tiers + an ignore/justify channel are mandatory.
3. **The threshold is empirically anchored:** ≥4 concurrent tells is where the study itself drew "High". That is a defensible default for a FAIL gate.

---

## 3. Ranked TELL REGISTRY (machine-checkable form)

Ranked by **corroboration** (how many independent sources name it) × **detectability** × **claimed strength as a tell**.

**Detectability codes** — `[MINE]`, a taxonomy synthesised from S3's CLI/Browser/LLM-only trichotomy and S4's `extract()`/`score()` split:

| Code | Meaning | Needs |
|------|---------|-------|
| **R** | Source regex — class strings, raw CSS/JSX text | nothing (grep) |
| **C** | Computed CSS — resolved styles on a rendered page | headless browser |
| **D** | DOM structure / text content | parsed HTML or browser |
| **P** | Pixel / screenshot — needs vision | render + VLM |
| **H** | Human or LLM judgement only | reviewer |

### Tier 1 — Confirmed tells (4+ independent sources, deterministic)

| # | ID | Tell | Detect | Concrete rule | Sources |
|---|----|------|--------|---------------|---------|
| 1 | `SLOP-COL-PURPLE` | **Indigo/violet/purple as the filled accent on CTAs and links** ("VibeCode Purple", "the Purple Problem") | **R+C** | S6's shipped regex: `/bg-(indigo\|violet\|purple)-600/`. S4's rule fires only when `filledAccentCount >= 1`, where a filled accent = a button/link with `bg.a >= 0.5` or a gradient fill; **outline/ghost purple and decorative purple are explicitly excluded** | S1,S2,S3,S4,S5,S6,S7,S8 |
| 2 | `SLOP-COL-GRADIENT` | **Purple→blue / cyan→magenta gradient** on hero background or button | **R+C** | S6: `/from-purple-\d+ to-(blue\|pink)-\d+/`; computed `background-image: linear-gradient(...)` with 2 stops crossing hue ~250–330 | S1,S2,S3,S4,S6,S7,S8 |
| 3 | `SLOP-TYPE-GRADTEXT` | **Gradient headline** (`background-clip: text` + gradient fill) | **C** | S8 gate 2: *"including a `background-clip: text` gradient headline… No genre allows gradient text."* | S3,S4,S6,S8 |
| 4 | `SLOP-CARD-STRIPE` | **Thick coloured stripe on one edge of a rounded card** (usually left, 3–6px) | **C+D** | S3 calls it *"The most recognizable tell of AI-generated UIs"*; S5 quotes a designer: *"colored left borders are almost as reliable a sign of AI-generated design as em-dashes for text"*; S8 gate 5 | S3,S4,S5,S7,S8 |
| 5 | `SLOP-HERO-CENTERED` | **Full-viewport centred hero, everything on one axis** | **C+D** | S8 gate 6 is the sharpest operationalisation: fail if hero is `min-height: 100vh` with everything centred **OR** eyebrow + title + lede + CTA all share one centred vertical axis. Auto-fail. | S1,S2,S3,S4,S5,S8 |
| 6 | `SLOP-LAY-3CARDS` | **Three equal columns, icon-above-heading, identical card heights** | **D+C** | S8 gate 3; S2: *"identical padding, identical border radius, and identical card heights"*; S4 "icon-topped feature cards" | S1,S2,S3,S4,S5,S8 |
| 7 | `SLOP-TYPE-INTER-ONLY` | **Inter (or Roboto/Open Sans/Poppins/Lato/system) as the *only* face, especially on a centred hero** | **C** | S4's rule #9 is a **conjunction**: "Centered + Inter". S8 gate 1 + gate 37 (≤3 families) + anti-pattern *"A one-font page is a template page."* | S1,S2,S3,S4,S5,S6,S7,S8 |
| 8 | `SLOP-HERO-BADGE` | **Pill badge / eyebrow floating directly above the H1** | **D** | S4 rule #13; S8 gate 54 hard-bans the *tag-left/header-right* variant outright and sets eyebrows **default-OFF** | S3,S4,S5,S8 |
| 9 | `SLOP-FX-GLASS` | **Glassmorphism** — `backdrop-filter: blur()` on translucent floating panels used decoratively | **C** | S4 rule #6; S7: *"the two dominant CSS fingerprints are shadcn/ui defaults and glassmorphism"*; S8: *"Glassmorphism can work when it communicates depth… It cannot work as decoration."* | S3,S4,S7,S8 |
| 10 | `SLOP-TYPE-ITALIC-ACCENT` | **One serif-italic accent word inside an otherwise-upright sans headline** | **C+D** | S8 gate 38a: *"Italic headers — above all the single italicised emphasis-word inside an upright headline — are a top AI tell"*; `font-style: italic` on any `h1–h6`/`.hero__title`/wordmark/stat/`<em>` inside a heading | S3,S4,S5,S8 |

### Tier 2 — Strong tells (2–3 sources, deterministic)

| # | ID | Tell | Detect | Concrete rule | Sources |
|---|----|------|--------|---------------|---------|
| 11 | `SLOP-SHAPE-RADIUS` | **Extreme / uniform border-radius** (`rounded-2xl`, `rounded-3xl`, "uniform 16px everywhere") | **R+C** | S6: `/rounded-(2xl\|3xl)/`; S2: *"card-based layouts with uniform 16px border radius everywhere"* | S1,S2,S3,S6 |
| 12 | `SLOP-FX-GLOW` | **Saturated coloured box-shadow glow** on buttons/cards, esp. on dark | **C** | S4 rule #7 "colored glows and colored box-shadows"; S8 *"Shadow-glow on dark"* → fix is elevation via lightness, not shadow | S3,S4,S5,S8,S9 |
| 13 | `SLOP-COL-DARKMUTED` | **Permanent dark mode + medium-grey body text**, contrast barely passing or failing | **C** | S4 rule #10; S8 gates 40–41 give real thresholds: body **WCAG 4.5:1 / APCA Lc ≥ 60**, large text/icons/focus rings **3:1 / Lc ≥ 45**; button-text-≈-button-fill fails if within **5% lightness AND 0.05 chroma in OKLCH** | S3,S4,S5,S8 |
| 14 | `SLOP-FONT-TEMPLATED` | **Templated display-font set**: Space Grotesk, Instrument Serif, Geist, Syne, Fraunces as page default | **C** | S4 rule #1, verbatim list. ⚠️ See §6.3 — this list *is* the anti-slop escape hatch of 2025 | S4,S5 |
| 15 | `SLOP-ICON-EMOJI` | **Emoji as feature-card / nav / step / pricing icon** (`✨ 🚀 ⚡ 🔥 🎯 ✅`) | **D+R** | S4 rule #8 "sidebar-emoji"; S8 gate 30 also fails **mixing ≥2 icon libraries** (Material + Heroicons + Lucide on one page) | S3,S4,S5,S8 |
| 16 | `SLOP-LAY-STEPS` | **Numbered "1 · 2 · 3" step sequence** / tiny numbered section labels | **D** | S4 rule #11; S3 "Tiny numbered section labels" | S3,S4,S5,S8 |
| 17 | `SLOP-PROOF-STATBAR` | **Stat banner row** — "10K+ users · 99.9% uptime · 4.9★" | **D+R** | S4 rule #12; **S8 gate 46 is the stronger form**: fail any quantitative claim the user did not supply. Fix: replace with `—` + "metric to confirm", or drop the section. Also: *"A stat is never the hero's sole headline."* | S3,S4,S5,S8 |
| 18 | `SLOP-TYPE-ALLCAPS` | **All-caps section labels / all-caps body**, often with wide tracking | **C** | S4 rule #16; S3 "All-caps body text" + "Wide letter spacing on body text". S8 gate 55 adds the failure mode: all-caps display with `line-height < 1.0` → cap-collision on wrap; floor is **1.0**, recommended **1.02–1.08** | S3,S4,S5,S8 |
| 19 | `SLOP-BG-ORB` | **Aurora blob / mesh gradient / radial halo / floating orbs** behind the hero, unmotivated | **C+P** | S8 gate 29: fail if >1 accent colour, >~5% footprint, or animating mesh-gradient page-wide. S2: *"abstract 3D blobs floating in space"* | S2,S3,S7,S8 |
| 20 | `SLOP-MOT-GENERIC` | **Motion without meaning** — `transition-all`, uniform `hover:scale-105`, bounce/elastic easing, fade-up-on-scroll on every section | **R+C** | S8 gates 10–14 (+ full microinteraction catalogue). S3 "Bounce or elastic easing", "Layout property animation". S2: *"Hover states that do nothing"*, *"buttons that snap instead of easing"*. S9 ships these as detector rules | S2,S3,S8,S9 |
| 21 | `SLOP-COL-PURE` | **Pure `#000` / `#fff` and zero-chroma neutrals** | **R+C** | S8 gate 7 + gate 22 (`oklch(... 0 ...)` banned; min **0.005 chroma**). S9: *"Don't use pure black/gray (always tint)"* | S3,S8,S9 |
| 22 | `SLOP-LAY-NESTCARD` | **Card inside a card** — containment with no semantic reason | **D** | S3 "Nested cards"; S8 gate 4; S9 *"Don't wrap everything in cards or nest cards inside cards"* | S3,S8,S9 |
| 23 | `SLOP-CHROME-NAV` | **The AI nav** — wordmark-left + 4–5 inline links + CTA-right + full width + 1px hairline bottom + white bg | **D+C** | S8 gate 42. Rationale is the sharpest in the corpus: *"the shape is genre-blind: it lands the same on a wedding photographer's portfolio, a bakery, a B2B SaaS, and a manifesto"* | S8 |
| 24 | `SLOP-CHROME-FOOTER` | **The AI footer** — 4 link columns (Product · Company · Resources · Legal) + social row + tiny copyright | **D** | S8 gate 43 | S8 |
| 25 | `SLOP-CHROME-REDRAWN` | **Re-drawn UI chrome** — hand-built fake browser bar (URL pill + traffic-light dots), fake phone notch, fake terminal/IDE window around a `<pre>` | **D+P** | S8 gate 47: *"the model invented a UI that already exists in the user's environment"* | S8 |
| 26 | `SLOP-COPY-CLICHE` | **Copy clichés** | **R** | Literal strings from S2: *"Build the future of work"*, *"Your all-in-one platform"*, *"Scale without limits"*, *"Empowering teams to build better products"*. S3 adds categories: "Marketing buzzword", "Aphoristic-cadence copy", "Theater framing copy", "Em-dash overuse", "Same text repeated in one container". Prose tells: *"In today's fast-paced…"*, *"In the ever-evolving landscape of…"*, *"Look no further"*, *"Let's dive in"*, *"Moreover/Furthermore"* | S2,S3 |
| 27 | `SLOP-COPY-PLACEHOLDER` | **Placeholder names & startup-bingo product names** | **R** | S8 gate 19: "Jane Doe / John Smith" + Acme, Nexus, Pulse, Unleash, Seamless, Supercharge | S8 |
| 28 | `SLOP-LAY-FAQ` | **"Frequently asked questions" + 3+ collapsible Q&As** appended to a landing page | **D** | S4 rule #14 (weakest of the 14 — very high legitimate-use rate) | S4 |
| 29 | `SLOP-BG-GRIDLINES` | **Decorative grid-line background** (repeating-linear-gradient graph paper) | **C** | S3 "Decorative grid-line background", "Repeating-gradient stripes" | S3 |
| 30 | `SLOP-SPACE-MONOTONE` | **Monotonous spacing / flat type hierarchy** — every section identical rhythm, no rule, no ornament, no colour shift | **C** | S8 gate 9 + gate 24 (any padding/gap off the named 4px scale — *"Arbitrary `padding: 17px` is a tell"*) + gate 25 (prose `max-width` outside **45–75ch**) | S3,S8 |

### Tier 3 — Defect class (Axis C: correlated with slop, but really bugs)

Included because S3 ships them in the same 64-pattern registry and they are the cheapest true-positives.

| # | ID | Tell | Detect | Rule |
|---|----|------|--------|------|
| 31 | `DEF-A11Y-CONTRAST` | Low-contrast text | **C** | WCAG 4.5:1 / APCA Lc 60 body; 3:1 / Lc 45 large+icons+focus (S8 gate 40) |
| 32 | `DEF-LAY-HSCROLL` | Horizontal scroll at any width 320–1920px | **C** | Fix is `overflow-x: clip` on **both** `html` and `body` — `clip` not `hidden`, to preserve sticky/fixed descendants (S8 gate 34) |
| 33 | `DEF-LAY-OVERFLOW` | Content overflowing container / text occluded / clipped by overflow parent | **P+C** | S3 "Content overflowing container", "Text occluded by overlapping element" |
| 34 | `DEF-RT-ERROR` | Uncaught script error on load; broken/placeholder image; content invisible at rest | **C (runtime)** | S3 "General Quality" set — the highest-signal, zero-FP rules |
| 35 | `DEF-STATE-INPUT` | Input states: border-width shifts between states, focus ring built from `border` not `outline`, input height ≠ button height (44px floor), helper-text slot collapses when empty, disabled signalled by `opacity` alone | **C** | S8 gate 39 — five sub-checks, all deterministic |
| 36 | `DEF-A11Y-STATES` | Any interactive element lacking `:focus-visible` / `:active` / `:disabled`; any `transform`/keyframe without a `prefers-reduced-motion` fallback; auto-rotating content without pause (WCAG 2.2.2) | **C** | S8 gates 26, 27, 18 |
| 37 | `DEF-RESP-WRAP` | Button/nav/CTA label wrapping to two lines at any width 320–1920px | **C** | S8 gate 49 |
| 38 | `DEF-GRID-MINMAX` | `grid-template-columns` with a bare `1fr` track containing an `<img>` → blows past viewport (`1fr` resolves to `minmax(auto, 1fr)`) | **R+C** | S8 gate 50 — one-character fix: `1fr` → `minmax(0, 1fr)` |

### Tell class 4 — Human/LLM-only (cannot be gated mechanically)

| ID | Tell | Why not mechanical |
|----|------|--------------------|
| `SLOP-IMG-AILOOK` | AI-illustration look: smooth-mesh-blob characters with no joint articulation, "3D abstract faceless humans holding glowing orbs", Midjourney-default symmetric lighting, the corporate-doodle person with one eye larger (S2, S8) | Needs vision + taste |
| `SLOP-STRUCT-TEMPLATE` | The generic AI macrostructure: Hero → 3 features → CTA → footer (S8 gate 8) | Detectable in the limit, but any real DOM rule over-fires on legitimate pages |
| `SLOP-POV-ABSENT` | **No point of view** — the page would fit any brief (S1, S8 axis D "Specificity") | This is the actual disease. Only judgement sees it |
| `SLOP-DECOR-UNMOTIVATED` | Decoration with no semantic anchor: a "42" in the corner with no edition meaning, a cursor floating beside a hero, a Pantone chip with no colour rationale (S8 gate 45) | Requires knowing what the content means |
| `SLOP-PROOF-FABRICATED` | Invented metrics/testimonials — the *number shape* is regex-able, its *truth* is not (S8 gate 46) | Needs ground truth from the user |

**Coverage arithmetic [MINE].** 30 slop tells + 8 defect rules are mechanically checkable; 5 are judgement-only. That is ≈88% machine coverage — consistent with S3's own ratio (**59 of 64 patterns have deterministic detectors; 5 remain "LLM only"**) and with S9's **60 deterministic rules**. This convergence is the strongest single justification for splitting `parley-design` (doctrine, covers the 12%) from `parley-design-check` (tooling, covers the 88%).

---

## 4. Structural causes — why the median wins

**FACT (S11).** *Verbalized Sampling* (arXiv 2510.01171, 83pp, 31 figures, 44 tables) attributes mode collapse to **typicality bias in preference data**: annotators *"systematically favor familiar text"*, so alignment amplifies it. Their training-free remedy — ask the model to *"Generate 5 X and their corresponding probabilities"* — **increases diversity 1.6–2.1× over direct prompting** with no loss of safety or factual accuracy, and larger models benefit more.

**FACT (S13).** Corroborating literature: RLHF's reverse-KL objective is *mode-seeking*, favouring single high-probability solutions; the standard KL penalty only partially alleviates collapse.

**FACT (S1).** The median-of-training-data mechanism, stated exactly: *"you're not getting design, but the median of every Tailwind CSS tutorial scraped from GitHub between 2019 and 2024"*, and *"LLMs don't have taste… Taste requires lived experience… They have statistical correlations between tokens."*

**FACT (S6).** The **default-library path**: Adam Wathan (Tailwind's creator) posted in **August 2025** that he was sorry for making every Tailwind UI button `bg-indigo-500` five years earlier, *"because now every AI-generated interface on earth is purple."* Chosen as *"a neutral, inoffensive placeholder that worked well in demos"* → thousands of tutorials copied it → the training corpus encoded "modern web design = purple button."

**FACT (S7).** The **agent-optimised-library path**: *"shadcn in particular is a library that is explicitly designed to be copy-pasted by AI agents, which means every AI-generated landing page without stylistic intervention converges on the shadcn visual."*

**FACT (S1, design-system corpus).** The **unconstrained-prompt path**: with no constraints, the highest-probability answer to an open design question is the most common look in training data. Constraint is the whole intervention.

**[MINE] — two causes the discourse under-names:**

- **Symmetry/safety bias.** Centring is the lowest-variance layout decision: it cannot be "wrong", it needs no editorial judgement about what leads, and it survives any content length. Every asymmetric layout is a *claim* about hierarchy; a centred one is an abstention. This is why `SLOP-HERO-CENTERED` and `SLOP-SPACE-MONOTONE` co-occur — they are the same abstention at two scales. Supported by S8's decision to make gate 6 an **auto-fail** rather than a warning.
- **Decoration as filler.** S2's *"Hover states that do nothing"* and S8's gate 45 describe the same thing from opposite ends: when the model has a slot and no content, it emits ornament. **The ornament is a marker of an unanswered content question.** Corollary: every unmotivated decoration is a place where the brief was thin.

---

## 5. Counter-evidence and nuance

### 5.1 Where generic is correct

**FACT (S10).** Explicitly names "Scenario A: The Good Enough Campaign" — conference landing pages, single ad campaigns, feature waitlists, internal tools, MVP tests — targeting *"70% Brand Match"* where *"speed is more important than perfect brand alignment"*.

**FACT (S7).** Practitioner counter-arguments from the HN thread:
- *"a site that's somewhat standardized and has things where we expect them is far more useful"*
- *"slop is what you get when you average every web style together. Qt works because there's really only one way Qt looks."* — i.e. **constraint, not novelty, is what removes the slop feeling**. Also: *"asking it to make it look like a Qt app… removed almost all feeling of slop."*
- *"If you don't have taste or want to be opinionated, why do you care?"*

**FACT (Jakob's Law / NN/g corpus).** Users spend most of their time on other sites and expect yours to work the same way; familiar patterns let users *recognise* rather than *recall*. ⚠️ **Caution:** secondary sources circulate a *"familiar patterns execute 3–5× faster than novel ones"* figure that I could **not** trace to an NN/g primary. Treat as unsourced. (Ironic self-application: this is exactly `SLOP-PROOF-FABRICATED`.)

**[MINE] — the classes where anti-slop rules must stand down:**
1. **Enterprise CRUD / admin / internal tools** — the design system *is* the brief; deviation is cost.
2. **Regulated / safety UI** — convention is a safety property.
3. **Accessibility-driven form** — visible focus rings, ≥44px targets, 4.5:1 contrast are not "template look"; a rule that penalises them is broken.
4. **Platform-native surfaces** — HIG/Material conformance is correct; "distinctive" is wrong.
5. **Conversion-critical checkout/auth** — familiarity dominates novelty.
6. **Documentation / reference** — scannability beats art direction.

### 5.2 The trust penalty is real but measured in an adjacent domain

**FACT (S14).** In advertising, correctly identifying content as AI-generated reduced preference by **21.2%** ("Detection Penalty") — yet AI ads still won overall, preferred by **50.3%** even when known to be AI-generated.
**[MINE] — INFERENCE, flagged:** this is ads with explicit disclosure, not UI with implicit tells. It supports "being *recognised* as AI costs something" but does **not** establish a conversion figure for UI. Do not quote it as a UI statistic.

### 5.3 Anti-slop becomes its own cliché — with hard evidence

This is the most important finding for `parley-design`'s architecture, and it is *not* speculation.

**FACT — direct contradiction between two 2026 sources.** S6's prescribed fix for the AI look is a specific palette and font stack, verbatim:
```js
fontFamily: {
  display: ['"Fraunces"', 'serif'],
  sans:    ['"IBM Plex Sans"', 'sans-serif'],
}
colors: { ink: {50:'#f6f5f1', 500:'#3d3a32', 900:'#1a1814'},
          ember: {400:'#e8775a', 600:'#c45530'} }
```
S4's detector rule #1 (`slop-fonts.js`) flags **Fraunces** — along with Space Grotesk, Instrument Serif, Geist, Syne — as a *slop* font: *"Templated display fonts … as page default."* **The 2025 escape hatch is the 2026 tell.**

**FACT (S8, local).** Hallmark — a mature anti-slop skill — had to add gates against **its own** favourite outputs:
- Gate 21: *"Did I default to the **Specimen** macrostructure … when the brief did not explicitly call for editorial / foundry / specimen energy? (Specimen fall-through is banned.)"*
- Anti-pattern *"Default-attractor sameness"*: *"Two consecutive Hallmark outputs in the same project use the same macrostructure… The page looks redesigned only because copy changed."*
- Gate 57: *"Studied DNA discarded for a catalog theme"* — reverting to the built-in theme list is *"the attractor pull."*
- Gate 32 requires a *different variation knob*: *"Two Bento Grids with `tiles=6, spans=irregular, accent=corner-only` are the same Bento."*
- Axis F of the pre-emit critique scores **"structural distance, not visual distance — colour-swaps don't count as variety."**

**FACT (2026 trend corpus).** The industry-level backlash is already legible as a style: "anti-design", "tactile brutalism", *"raw geometry, aggressive color contrast, and simulated physical textures to prove human authorship"*, *"friction, texture, glitch, and nostalgia"*. That is a nameable aesthetic — therefore a future tell.

**[MINE] — the design rule this forces.** Any fixed anti-slop list is a **moving target with a half-life**. Therefore:
- The registry must be **versioned and dated** (`registry: 1.0.0`, each rule carrying `added`/`deprecated`).
- Rules must express **prohibitions of defaults**, never **prescriptions of alternatives**. "Do not ship the model's first palette unmodified" survives; "use Fraunces" does not.
- `parley-design` must gate on **structural distance from the project's own prior outputs** (S8 axis F / gate 32), not only on the global tell list. Self-repetition is the failure mode that a global registry structurally cannot see.

---

## 6. What to do INSTEAD — the philosophical frame

**Six propositions, each grounded.**

**P1 — Slop is not ugliness; it is the absence of a decision.**
S1: models have *"statistical correlations between tokens"*, not taste. S8's axis A asks: *"Is there a clear *why* — a position the page is taking? Or is it just a layout?"* and axis D: *"Does this look like *this brief* — or a generic 'page that could be anyone'?"* → **Every token in the system must be traceable to a stated reason.** A design system without a rationale column is a slop generator with extra steps.

**P2 — Diverge before you converge — and the paper says how much it buys you.**
S11 measures **1.6–2.1× diversity** from a training-free prompt that asks for *N* candidates at once. Parley Deck's DIVERGE method (each agent proposes a *different* direction) is the multi-agent form of exactly this — with the added property that the agents are *different models*, so the modes are drawn from different training medians rather than one. **This is the empirical justification for the whole `parley-design` collaboration method.**

**P3 — One direction wins whole. Averaging *is* the disease.**
S1 and the whole corpus define slop as *the average of the training set*. A consensus mechanism that blends four visual directions **reconstructs the median in miniature** and produces a locally-generated slop. Therefore the owner's stated rule — *one direction wins whole, graft 2–3 concrete details from the losers* — is not stylistic preference, it is the only merge policy consistent with the diagnosis.
**[MINE] — make the graft *bounded and typed*.** Cap grafts at 3, require each to name (a) the losing direction it came from, (b) the exact token/component/rule it changes, (c) why it does not contradict the winner's thesis. An unbounded graft list is an average.

**P4 — Constraint, not novelty, is the cure.**
S7: *"Qt works because there's really only one way Qt looks."* S9's entire architecture is contextual setup — audience, brand positioning, voice, **anti-references**, colours, type, components — captured in `PRODUCT.md` + `DESIGN.md` *before* generation. S8 gate 48 bans **mid-render token improvisation**: every colour and font must resolve through `var(--token-name)`; an inline hex is *"the model picked the theme, then forgot it and freestyled."* → **The deliverable of `parley-design` is a token contract, and the first job of `parley-design-check` is to prove nothing escapes it.**

**P5 — Decoration must be motivated.**
S8 gate 45, the single best formulation in the corpus: decoration needs *"a semantic anchor in the content — a cursor inside a typed command (signals 'you'd type next'), a numeral that names an issue / year / version / chapter, a gradient that responds to interaction, a stamp that names an authorship or date. Random ornaments … are slop."* → Ship a `motivation:` field next to every decorative element. Unmotivated = removed.

**P6 — Distance from your own last output is a first-class constraint.**
S8 gate 8 + gate 32 + axis F. A project shipping ten pages that all pass the global registry but share one skeleton has failed. **Score structural distance, not visual distance. Colour-swaps do not count as variety.**

---

## Transferable to parley-design / parley-design-check

Ranked by value × implementability.

**1. Adopt S3's two-class rule taxonomy as the spec's top-level split.** `slop` rules ("tells of AI-generated interfaces") vs `quality` rules ("design mistakes regardless of origin"). *Why:* it is the only way to give the two skills clean ownership — `parley-design` owns intent and the slop class; `parley-design-check` owns the quality class outright and enforces the deterministic subset of the slop class. Source: https://impeccable.style/slop/

**2. Adopt S3's detection-method label as a required field on every rule.** Their trichotomy is `CLI` (deterministic on source files) / `Browser` (deterministic, needs real layout) / `LLM only` (no deterministic detector). Extend to my 5-code taxonomy `R|C|D|P|H` if finer granularity is wanted. *Why:* it makes the spec honest about what is gateable, and tells the CI which lane runs the rule. **59 of 64** of S3's patterns are deterministic; plan for that ratio.

**3. Steal S4's rule-module interface verbatim as the plugin contract.** Each pattern file exports `{ id, label, shortLabel, description, category, thresholds, extract(ctx), score(signal, T) }`, where **`extract()` is serialised into the browser and `score()` runs in Node**. *Why:* it cleanly separates *observation* from *judgement*, makes thresholds tunable without touching detection, and is already proven at N=1,590. Files: `src/patterns/<id>.js` in https://github.com/AdrianKrebs/design-slop-cop

**4. Adopt S4's conjunctive rule discipline — never gate on a bare aesthetic fact.** The purple rule fires only on `filledAccentCount >= 1` where filled = `bg.a >= 0.5` or gradient fill on a button/link, and **explicitly excludes** outline/ghost buttons and decorative purple. The Inter rule is `Centered + Inter`, not "Inter". *Why:* single-property rules are how a checker becomes an aesthetic tyrant. Every Tier-1 rule in §3 should be expressible as a conjunction of ≥2 observable conditions.

**5. Anchor the gate threshold in S4's measured base rate.** `0–1 tells = none/low` (46% of the wild population), `2–3 = medium` (32%), `4+ = high` (22%). Score = `round(100 × flagged / total)`. *Why:* it is the only empirically grounded threshold available, and it prevents a first-offence hard-fail.

**6. Take S9's CLI contract wholesale for `parley-design-check`.** `npx <tool> detect <path|file|url>`; `--json` for CI; `--fast` regex-only mode that skips jsdom; **exit `0` = clean, exit `2` = anti-patterns found**; scans `HTML, CSS, JSX, TSX, Vue, Svelte`; plus `ignores list` / `ignores add-file "<pattern>"` / `ignores add-value <rule> --reason "<justification>"`. *Why:* it is a shipping, Apache-2.0, vendor-neutral contract that already solves the "how do I legitimately override a heuristic rule" problem — and `--reason` is exactly the kind of typed artifact the owner's protocol style wants. Source: https://github.com/pbakaus/impeccable

**7. Port the numeric thresholds from S8 — they are the most specific in the corpus.** All directly checkable: accent ≤ **~5%** of any viewport by area; neutral chroma ≥ **0.005** (no zero-chroma greys); prose `max-width` **45–75ch**; spacing on a named 4px scale (`--space-3xs`…`--space-5xl`); **≤3 font families**, outlier face in **≤2 slots**; all-caps display `line-height` **≥1.0** (rec. 1.02–1.08); hero `padding-block-end` **≥1.3×** `padding-block-start`; hero must fit at **1280×800**; no horizontal scroll **320–1920px** via `overflow-x: clip` on both `html` and `body`; body contrast **4.5:1 / APCA Lc 60**, large+icons+focus **3:1 / Lc 45**; button-text-vs-fill fails within **5% lightness AND 0.05 chroma** in OKLCH; input **44px** floor; helper-text slot reserves **1lh**. File: `/Users/tomasfecko/.claude/skills/hallmark/references/slop-test.md`

**8. Make the honesty gate (S8 gate 46) a hard, non-waivable FAIL.** Any quantitative claim the user did not supply — "10× faster", "trusted by 50,000+ teams", "99.9% uptime", "+47% conversion" — fails. Three sanctioned fixes, in order: replace with `—` + labelled block; ask the user and pause; rebuild without the proof slot. *Why:* it is the one rule with an ethical rather than aesthetic basis, it is trivially regex-detectable in its number-shape, and it is the difference between "looks generated" and "lies".

**9. Ship the pre-emit self-critique as the doctrine skill's typed pre-FINAL artifact.** S8's six axes scored 1–5 — **P**hilosophy, **H**ierarchy, **E**xecution, **S**pecificity, **R**estraint, **V**ariety — run *before* the gate sweep, with **any axis < 3 forcing a revision pass**, and the result stamped into the artifact (`/* pre-emit critique: P5 H4 E5 S4 R5 V5 */`). *Why:* it is the machine-readable form of "did anyone actually decide anything", it maps 1:1 onto Parley's cross-review round, and the stamp is a versioned artifact — exactly the AG-UI-style protocol shape the owner asked for.

**10. Build the anti-self-repetition gate that no public tool has.** S8's `.hallmark/log.json` + CSS macrostructure stamp + gate 8 + gate 32 (*"Two Bento Grids with `tiles=6, spans=irregular, accent=corner-only` are the same Bento"*) + axis F (*"structural distance, not visual distance — colour-swaps don't count as variety"*). *Why:* every public detector is a *global* registry and is structurally blind to a project repeating itself. This is `parley-design`'s clearest differentiator, and it is the only defence against the anti-slop cliché documented in §5.3.

**11. Encode DIVERGE with S11's citation and a bounded graft rule.** Verbalized Sampling gives **1.6–2.1×** diversity from asking for *N* candidates at once; Parley's *N different models* strengthens it. Then: one direction wins whole; **cap grafts at 3**; each graft names its source direction, the exact token/component it changes, and why it does not contradict the winner's thesis. *Why:* it converts the owner's stated method into a protocol rule with a citation and a hard numeric bound, and it structurally forbids the "average of four directions" failure.

**12. Adopt S8's audit report shape as `parley-design-check`'s human output.**
```
[severity] Tell name — file:line
  why it's a tell (one line)
  → fix (one line)

Summary — N critical · M major · K minor
Verdict — [ships as slop | reads as AI-generated | close, fix the minors]
```
*Why:* three severities, a verdict enum, and a file:line anchor — trivially convertible to the `ReportFindings`/SARIF shape and to a PR annotation.

**13. Use the Tell Registry in §3 as the v1 seed.** 30 slop rules + 8 defect rules + 5 human-only entries, each already carrying an ID, a detection code, a concrete rule, and a source list. IDs are namespaced (`SLOP-COL-*`, `SLOP-TYPE-*`, `SLOP-LAY-*`, `SLOP-CHROME-*`, `SLOP-FX-*`, `SLOP-MOT-*`, `SLOP-COPY-*`, `DEF-*`) so categories can be enabled/disabled per project genre.

**14. Version and date the registry, and require a `sunset` review.** Each rule carries `added:`, `sources:`, `status: active|deprecated`, `confidence: confirmed|strong|weak`. *Why:* §5.3 proves rules expire — Fraunces went from *the fix* to *the tell* in roughly one year.

---

## Do NOT copy

**1. Do not ban fonts outright — Inter is not the tell.**
S4's own rule is the **conjunction** `Centered + Inter` (rule #9), and S8's is *"Inter used as **both** display and body, with no pairing face"*. Inter is an excellent UI typeface; a rule that fails it unconditionally will fire on every well-built dashboard in existence. *Reason:* the tell is *absence of a typographic decision*, not any particular face.

**2. Do not copy S6's prescribed palette and font stack (`ink`/`ember`, Fraunces + IBM Plex).**
S4's `slop-fonts.js` already flags Fraunces as a slop font. *Reason:* substituting one shared palette for another shared palette relocates the sameness rather than removing it. Prescribe *process*, never *values*.

**3. Do not copy "eliminate rounded corners and drop shadows" (S6's remediation step 3).**
*Reason:* border-radius and elevation are affordance and depth signals with real usability function; banning them is aesthetic dogma that produces a recognisable flat-brutalist look — i.e. a future tell — and degrades the perception of interactive vs static surfaces. The defensible rule is the *uniformity* one (`SLOP-SHAPE-RADIUS`: every card at the same extreme radius), not the property itself.

**4. Do not copy any detector as a hard blocker at its stated accuracy.**
S4 self-reports **5–10% false positives**. *Reason:* at 30+ rules a 5–10% FP rate makes a clean page fail routinely. Deterministic ≠ correct. Slop rules → warn + require justification; only the defect class (`DEF-*`, Tier 3) and the honesty gate should hard-fail.

**5. Do not adopt a severity-free flat list.**
S3 documents 64 patterns with **"no explicit severity hierarchy"** across them. *Reason:* a gate without severity cannot decide what blocks a merge. S8's critical/major/minor split is the minimum; borrow that instead.

**6. Do not use screenshot+LLM judging as the primary scorer.**
S4 rejected it deliberately: *"I intentionally do not take screenshots and let the LLM judge them… Deterministic checks against computed styles take the LLM out of the scoring loop."* *Reason:* an LLM scoring LLM output is a circular, non-reproducible, non-diffable gate. Reserve vision for the 5 human-only tells, and label those findings as advisory.

**7. Do not import Hallmark's genre/macrostructure/theme catalogue.**
Twenty-one named macrostructures, ~20 named themes (Specimen, Midnight, Brutal, Garden, Atelier, Newsprint, Terminal…), and a components cookbook (`N1`–`N13`, `Ft1`–`Ft8`, `H1`–`H9`…) live under `/Users/tomasfecko/.claude/skills/hallmark/references/`. *Reason:* (a) it is a specific house style, and the owner's requirement is **vendor-neutral**; (b) S8's own gates 21, 32 and 57 document that the catalogue *became* the attractor it was built to escape. Copy the **mechanism** (a named-alternatives table + a "you must differ from your last output" gate); do not copy the **contents**.

**8. Do not moralise the word "slop" inside the spec.**
S12's definition is about content that is *"mindlessly generated and thrust upon someone who didn't ask for it"* — and Willison's explicit caveat is *"Not all AI-generated content is slop."* *Reason:* a bespoke, reviewed page that happens to use a gradient is not slop. Rule names in the registry must be descriptive (`SLOP-COL-GRADIENT`) and findings must read as "this is a recognisable default", never as a verdict on the author.

**9. Do not ban the em-dash.**
S3's pattern is *"Em-dash **overuse**"*, not em-dash use, and the "em-dash = AI" heuristic is a text-detector superstition that punishes correct typography. *Reason:* S8 simultaneously **requires** `—` (U+2014) over `--` and `…` (U+2026) over `...`. Any rule here must be a *density* rule (e.g. per-1000-words rate against a corpus baseline), or be dropped.

**10. Do not apply the anti-slop gates to design-system-governed, regulated, or platform-native surfaces.**
S10 names the exemptions (conference pages, internal tools, MVP tests, *"70% Brand Match"*); S7: *"a site that's somewhat standardized and has things where we expect them is far more useful."* *Reason:* enterprise CRUD, checkout, auth, admin, HIG/Material surfaces and docs are cases where convention is the correct answer and distinctiveness is a defect. The spec needs a **`genre:` / `profile:` field that disables the Axis-A/B rule classes while keeping every `DEF-*` rule on.**

**11. Do not quote the AI trust/conversion numbers as UI evidence.**
S14's **21.2%** detection penalty / **50.3%** still-preferred figures are from *advertising with explicit source disclosure*, not from UI with implicit tells. The Jakob's-Law *"3–5× faster"* figure could not be traced to an NN/g primary at all. *Reason:* citing them would make the spec commit `SLOP-PROOF-FABRICATED` — the exact rule it enforces.

**12. Do not treat the tell list as terminal.**
S4's 14 rules, S5's 16, S3's 64, S9's 60 and S8's 58 overlap but do not agree, and §5.3 shows at least one item (Fraunces) inverting polarity within a year. *Reason:* ship the registry with a version, per-rule `added`/`confidence`/`sources`, and a stated review cadence. A frozen anti-slop list is a slop generator on a delay.

---

## Appendix — corroboration matrix for the Tier-1 tells

| Tell | S1 | S2 | S3 | S4 | S5 | S6 | S7 | S8 | Count |
|------|----|----|----|----|----|----|----|----|-------|
| Purple/indigo filled accent | ● | ● | ● | ● | ● | ● | ● | ● | **8** |
| Inter-only / centred + Inter | ● | ● | ● | ● | ● | ● | ● | ● | **8** |
| Purple→blue/pink gradient | ● | ● | ● | ● | – | ● | ● | ● | **7** |
| 3-up icon-topped feature cards | ● | ● | ● | ● | ● | – | – | ● | **6** |
| Full-viewport centred hero | ● | ● | ● | ● | ● | – | – | ● | **6** |
| Side-stripe card border | – | – | ● | ● | ● | – | ● | ● | **5** |
| Gradient headline (`background-clip:text`) | – | – | ● | ● | ● | ● | – | ● | **5** |
| Extreme/uniform border-radius | ● | ● | ● | – | – | ● | ● | – | **5** |
| Pill badge above H1 | – | – | ● | ● | ● | – | – | ● | **4** |
| Glassmorphism (decorative) | – | – | ● | ● | – | – | ● | ● | **4** |
| Serif-italic accent word | – | – | ● | ● | ● | – | – | ● | **4** |
| Coloured glow box-shadow | – | – | ● | ● | ● | – | – | ● | **4** |
