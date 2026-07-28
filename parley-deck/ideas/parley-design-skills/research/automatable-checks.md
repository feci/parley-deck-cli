# Digest — What design quality can be checked MECHANICALLY

**Purpose:** define the feasible scope of `parley-design-check` (the ENFORCEMENT skill) vs `parley-design` (the DOCTRINE skill).
**Date:** 2026-07-28. **Author:** research subagent (claude / Opus 5 1M).
**Method:** web research on primary sources (W3C, Deque, Playwright, Stylelint, Project Wallace) + **local empirical probes** (axe-core 4.12.1 driven by headless Chromium, @projectwallace/css-analyzer 9.9.3, stylelint 17.14.1, culori 4.0.2). Probe scratch dir: `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/probe/`.

Everything below is tagged **[FACT]** (verified from a primary source or reproduced locally — source given) or **[INFERENCE]** (my judgement, not verified).

---

## 0. The one-paragraph answer

**[INFERENCE]** Mechanical checking splits cleanly into three bands, and the band boundary — not the tool choice — is what should shape the protocol:

1. **Deterministic, near-zero-FP, no runtime** — *lexical + CSS-AST invariants* against a declared token set (hex-not-token, off-scale spacing, font-family count, `!important`, z-index sprawl, `transition: all`, near-duplicate colors by ΔE). These are ~80 % of "AI slop" tells and cost nothing.
2. **Deterministic but needs a real layout engine** — *computed-style + geometry checks* (contrast against the **computed** background, target size, horizontal overflow at 320 px, focus-ring presence/visibility, reflow). Needs headless Chromium. High value, moderate cost, low FP.
3. **Non-mechanical** — *taste, hierarchy, originality, "does it look like this brief"*. No tool decides this. This is exactly the adversarial-agent job that `parley-design` owns.

`parley-design-check` should ship band 1 and 2 and must **explicitly refuse** to pretend band 3 is mechanical.

---

## 1. Accessibility — exact thresholds and what tooling really covers

### 1.1 Standards status

**[FACT]** WCAG 2.2 is a **W3C Recommendation, published 12 December 2024** (that date is the current Recommendation republication; original Recommendation October 2023). Source: <https://www.w3.org/TR/WCAG22/>.
**[FACT]** SC **4.1.1 Parsing is "Obsolete and removed"** in WCAG 2.2. Source: <https://www.w3.org/TR/WCAG22/>.
**[FACT]** WCAG 3.0 is still a **Working Draft**; the most recent draft is **3 March 2026**, which renames "outcomes" to "requirements" and lists **~174 requirements**, with a scoring model replacing the pass/fail checklist. Candidate Recommendation is anticipated ~Q4 2027, Recommendation not before 2028. Sources: <https://accessibility.build/blog/wcag-3-march-2026-draft-outcomes-become-requirements>, <https://ratedwithai.com/blog/wcag-3-0-march-2026-update-timeline>.
**[FACT]** **APCA is NOT in the WCAG 3 draft.** The working group pulled APCA from the July 2023 WCAG 3 working draft; it was "only ever exploratory". WCAG 3 currently discusses a *visual contrast of text* requirement at the **Silver** tier (Bronze ≈ today's AA). Sources: <https://adrianroselli.com/2026/04/wcag3-contrast-as-of-april-2026.html>, <https://www.w3.org/WAI/GL/task-forces/silver/wiki/Visual_Contrast_of_Text_Subgroup>, GitHub `w3c/wcag3` issue #29.
> **[INFERENCE] Protocol consequence:** `parley-design` may *recommend* APCA Lc as a design-time heuristic (Hallmark does), but `parley-design-check` must **gate on WCAG 2.2 ratios**, because APCA has no normative status and no stable threshold table. Shipping an APCA gate would be gating on a moving, non-normative target.

### 1.2 Exact normative numbers (verbatim / near-verbatim from <https://www.w3.org/TR/WCAG22/>)

| SC | Level | Threshold — **[FACT]** |
|---|---|---|
| **1.4.3 Contrast (Minimum)** | AA | text + images of text: **contrast ratio ≥ 4.5:1**; **large-scale text ≥ 3:1**; exceptions: incidental, logotypes |
| **1.4.6 Contrast (Enhanced)** | AAA | **≥ 7:1**; large-scale text **≥ 4.5:1** |
| **1.4.11 Non-text Contrast** | AA | "User interface components and graphical objects" **≥ 3:1 against adjacent color(s)** |
| **1.4.4 Resize Text** | AA | text resizable to **200 %** without loss of content/functionality |
| **1.4.10 Reflow** | AA | no two-dimensional scrolling at **320 CSS px** width (vertical content) / **256 CSS px** height (horizontal content) |
| **1.4.12 Text Spacing** | AA | no loss at line-height **≥ 1.5×** font size, letter-spacing **≥ 0.12×**, word-spacing **≥ 0.16×** (para spacing ≥ 2×) |
| **2.4.7 Focus Visible** | AA | keyboard focus indicator is visible |
| **2.4.11 Focus Not Obscured (Min)** | AA | focused component **not entirely hidden** by author-created content |
| **2.4.12 Focus Not Obscured (Enh)** | AAA | **no part** of the focused component hidden |
| **2.4.13 Focus Appearance** | AAA | indicator area ≥ area of a **2 CSS px thick perimeter**, and **≥ 3:1 contrast** between focused/unfocused states |
| **2.5.8 Target Size (Minimum)** | AA (new in 2.2) | **"at least 24 by 24 CSS pixels"** |
| **2.5.5 Target Size (Enhanced)** | AAA | **44 by 44 CSS pixels** |

**[FACT] "Large text" definition:** "at least 18 point or 14 point bold or font size that would yield equivalent size for Chinese, Japanese and Korean (CJK) fonts" (= **24 px / 18.66 px bold** at the default 96 dpi mapping). Source: <https://www.w3.org/TR/WCAG22/#dfn-large-scale>.

**[FACT] The five exceptions to 2.5.8, verbatim** (source: <https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html>):
1. **Spacing** — "Undersized targets … are positioned so that if a 24 CSS pixel diameter circle is centered on the bounding box of each, the circles do not intersect another target or the circle for another undersized target"
2. **Equivalent** — "The function can be achieved through a different control on the same page that meets this criterion"
3. **Inline** — "The target is in a sentence or its size is otherwise constrained by the line-height of non-target text"
4. **User Agent Control** — "The size of the target is determined by the user agent and is not modified by the author"
5. **Essential** — "A particular presentation of the target is essential or is legally required"

> **[FACT — reproduced locally]** The **Spacing exception makes `target-size` a weak gate.** In my probe, an 18×18 px `<button>` **PASSED** axe `target-size` with the message: `Target has sufficient space from its closest neighbors. Safe clickable space has a diameter of 24px which is at least 24px.` So "all targets ≥ 24 px" is *not* what the automated rule enforces. If `parley-design-check` wants "24 px minimum, no exceptions" it must write its **own** geometry rule, not delegate to axe.

### 1.3 The "how much is automatable" claim — three different numbers, three different questions

**[FACT] The classic 20–30 % figure** counts *WCAG Success Criteria that are machine-testable*. Deque itself states its study result is "much higher than the widely-accepted belief that automated accessibility testing only provides 20-30 percent of coverage." Source: <https://www.deque.com/blog/automated-testing-study-identifies-57-percent-of-digital-accessibility-issues/>.

**[FACT] The GDS (UK Government Digital Service) empirical study**: 142 accessibility barriers were deliberately introduced into a page and run through **13 automated tools**. "The best performing tool in this category found **40 %** of the problems introduced, whereas the worst performing tool only picked up **13 %**." Commonly summarised as "automated testing finds around 30 %". Source: <https://alphagov.github.io/accessibility-tool-audit/results.html> (results page), overview at <https://alphagov.github.io/accessibility-tool-audit/index.html>.

**[FACT] Deque's counter-claim — 57 %**: from "over 2,000 audits … spanning over 13,000 pages (all first-time page assessments) covering nearly 300,000 issues", "**57 percent** of accessibility issues were completely covered by this automated testing." Deque explicitly flags the methodology switch: *"When we shifted the definition of 'accessibility coverage' beyond the number of WCAG Success Criteria to total volume of issues, the true impact of automated testing became clear."* Source: same Deque URL above; also <https://www.deque.com/automated-accessibility-coverage-report/>.

> **[INFERENCE] How to state this honestly in the skill:** the two numbers are **not in conflict** — they measure different denominators. Correct protocol wording: *"Automated tools cover roughly **30–40 % of WCAG success criteria** (GDS: 13–40 % per tool across 142 seeded barriers) but can account for up to **~57 % of issue volume** on a first-pass audit (Deque, 2,000+ audits / 300k issues), because the machine-detectable failures are the ones that repeat most often."* Cite both. Never cite 57 % alone — it flatters automation.

### 1.4 Engine-by-engine: what each tool actually does

#### axe-core — **[FACT, reproduced locally at v4.12.1]**
- `axe.getRules()` returns **105 rules**; `axe._audit.rules` → **96 enabled by default, 9 disabled by default**.
- **Experimental rules (7):** `css-orientation-lock`, `focus-order-semantics`, `hidden-content`, `label-content-name-mismatch`, `p-as-heading`, `table-fake-caption`, `td-has-header`.
- **Deprecated (5):** `aria-roledescription`, `audio-caption`, `duplicate-id-active`, `duplicate-id`, `landmark-complementary-is-top-level`.
- **`review-item` (2):** `frame-tested`, `hidden-content`.
- **Only ONE rule carries the `wcag22aa` tag: `target-size`.** (And it is **disabled by default** — see rule-descriptions: "Ruleset: Disabled (WCAG 2.2, not widely adopted)"; source <https://github.com/dequelabs/axe-core/blob/develop/doc/rule-descriptions.md>.)
- **Only THREE rules are `cat.color`:** `color-contrast` (`wcag2aa|wcag143|ACT`), `color-contrast-enhanced` (`wcag2aaa|wcag146`, disabled by default), `link-in-text-block` (`wcag2a|wcag141`).
- Other design-adjacent rules and their exact descriptions (source: rule-descriptions.md):
  - `meta-viewport` — "Ensure `<meta name='viewport'>` does not disable text scaling and zooming" (wcag2aa / 1.4.4, **default**)
  - `meta-viewport-large` — best-practice
  - `avoid-inline-spacing` — "Ensure text spacing through style attributes adjusts with custom stylesheets" (wcag21aa / 1.4.12, **default**)
  - `link-in-text-block` — "Ensure links are distinguished from surrounding text without relying on color" (wcag2a / 1.4.1, **default**)
  - `p-as-heading` — "Ensure bold, italic text and font-size is not used to style `<p>` as heading" (**experimental**)
  - `color-contrast-enhanced`, `identical-links-same-purpose` — AAA, **disabled**
- **[FACT — reproduced]** axe's `color-contrast` **catches the highest-value slop bug** — near-identical fg/bg. Probe output verbatim:
  `Element has insufficient color contrast of 1.01 (foreground color: #5c3ef6, background color: #5b3df5, font size: 10.0pt (13.3333px), font weight: normal). Expected contrast ratio of 4.5:1`
- **[FACT — reproduced]** axe's `color-contrast` **fails open on gradients**. Probe output verbatim, as an `incomplete` (needs-review) result, not a violation:
  `Element's background color could not be determined due to a background gradient`
  Corroborated: "axe-core's existing color contrast approach … assumes elements only have a single background color. If the tool can't come to a single color, it gives up." (GitHub `dequelabs/axe-core` issue #3390; Deque University color-contrast rule page.)
- **[FACT]** axe-core in **jsdom cannot do contrast at all** — "The color-contrast check doesn't work in JSDOM" (`dequelabs/axe-core` issue #595); `jest-axe` disables it. jsdom does no layout/rendering, so anything geometry- or paint-dependent is unavailable. Source: <https://github.com/dequelabs/axe-core/issues/595>, <https://www.npmjs.com/package/jest-axe>.
- **[FACT]** Deque's own doc `doc/accessibility-supported.md` states rules are gated on real AT/browser support (>1 % usage, 7 AT/browser combos tested), and that axe ships best-practice rules beyond strict WCAG.

#### Lighthouse — **[FACT]**
- "The Lighthouse Accessibility score is a weighted average of all accessibility audits. **Weighting is based on axe user impact assessments.**" — i.e. it is **axe-core underneath**, a *subset*.
- "Each accessibility audit is **pass or fail**. Unlike the Performance audits, a page doesn't get points for partially passing an accessibility audit."
- "Manual audits and low-impact / best-practices audits aren't included in the table because they don't affect your score."
- Source: <https://developer.chrome.com/docs/lighthouse/accessibility/scoring>.
- **[INFERENCE]** A Lighthouse a11y score of 100 is therefore *strictly weaker* than a clean `axe.run()` with WCAG 2.2 tags enabled. `parley-design-check` should **never** gate on the Lighthouse score — gate on axe directly.

#### pa11y — **[FACT]**
- v9.1.1 (npm). Two runners: **`htmlcs`** (HTML_CodeSniffer, the default) and **`axe`**. Runners are configurable per-run; `pa11y-ci` reads `.pa11yci` JSON, e.g. `"defaults": { "concurrency": 4, "standard": "WCAG2AA", "runners": ["axe"] }`.
- Standards: `Section508`, `WCAG2A`, `WCAG2AA` (default), `WCAG2AAA`.
- Has a **`-T, --threshold <number>`** option: "permits this number of errors, warnings, or notices, otherwise fail with exit code 2".
- Sources: <https://github.com/pa11y/pa11y>, <https://github.com/pa11y/pa11y-ci/issues/95>, <https://www.npmjs.com/package/pa11y>.
- **[INFERENCE]** pa11y adds a *second opinion* (HTML_CodeSniffer catches some things axe doesn't and vice-versa) but it still drives a headless browser, so it is the same engine tier as axe with no cost saving. Its one distinct feature worth stealing is the **numeric threshold + exit-code contract** — that is exactly the "machine-checkable gate" shape the owner wants.

#### IBM Equal Access (`accessibility-checker`, npm 4.0.29) — **[FACT]**
- Classifies findings into **Violation / Needs review / Recommendation** (+ a Hidden filter). "Needs review" is a first-class output type, not an afterthought.
- CLI: `npx achecker <file|dir|url|list.txt>`.
- IBM's own guidance: "If Needs Review issues appear during verification, we recommend pausing the verify process."
- Sources: <https://github.com/IBMa/equal-access/tree/master/accessibility-checker>, <https://www.ibm.com/able/toolkit/verify/automated>.
- **[INFERENCE]** The **three-state verdict (violation / needs-review / recommendation)** is the single best idea to lift from IBM. A binary pass/fail gate on design work produces either false alarms or false confidence; a tri-state maps perfectly onto Parley's "the agents adjudicate the ambiguous ones" model.

#### eslint-plugin-jsx-a11y (6.10.2) — **[FACT]**
- Self-described: "**Static AST checker** for a11y rules on JSX elements." ~37–38 rules (`alt-text`, `anchor-is-valid`, `click-events-have-key-events`, `no-static-element-interactions`, `label-has-associated-control`, `interactive-supports-focus`, `no-aria-hidden-on-focusable`, `no-autofocus`, `role-supports-aria-props`, `tabindex-no-positive`, …). Configs: `recommended`, `strict`.
- Documented limitation, verbatim from the README: "Static analysis tools cannot determine values of variables that are being placed in props before runtime … because it only catches errors in static code, use it in combination with `@axe-core/react`".
- Source: <https://github.com/jsx-eslint/eslint-plugin-jsx-a11y>.
- **[INFERENCE]** Framework-specific → **not vendor-neutral**. `parley-design-check` should *mention* it as an optional adapter, never depend on it.

---

## 2. Static analysis of styling — what you can prove from CSS text alone

### 2.1 Stylelint core (17.14.1) — **[FACT, verified locally: 149 core rules on disk]**
All of the following exist as core rules (verified by directory presence in `node_modules/stylelint/lib/rules/`), i.e. **no plugin needed**:

| Design-slop symptom | Exact stylelint core rule |
|---|---|
| hardcoded hex anywhere | `color-no-hex` |
| named colors (`red`, `rebeccapurple`) | `color-named` |
| inconsistent hex length (`#fff` vs `#ffffff`) | `color-hex-length` |
| `!important` | `declaration-no-important` |
| off-scale values for a given property | `declaration-property-value-allowed-list` / `declaration-property-value-disallowed-list` |
| unit sprawl (`px` + `rem` + `em` + `pt`) | `unit-allowed-list`, `unit-disallowed-list` |
| off-scale units per property (e.g. only `rem` for `padding`) | `declaration-property-unit-allowed-list` |
| banned properties (e.g. `float`, `zoom`) | `property-disallowed-list` |
| banned functions (e.g. `linear-gradient` on text) | `function-disallowed-list` |
| banned at-rules | `at-rule-disallowed-list` |
| missing generic fallback in font stack | `font-family-no-missing-generic-family-keyword` |
| duplicated names in a font stack | `font-family-no-duplicate-names` |
| specificity chaos | `selector-max-specificity`, `selector-max-id`, `no-descending-specificity`, `max-nesting-depth` |
| duplicate rules | `no-duplicate-selectors` |
| unknown / typo'd property values | `declaration-property-value-no-unknown` |
| `0px` noise | `length-zero-no-unit` |
| banned selectors (e.g. `[style*=]`) | `selector-disallowed-list` |
| stale vendor prefixes | `value-no-vendor-prefix` |

### 2.2 The token-enforcement plugin — **[FACT]**
- **`stylelint-declaration-strict-value`** (npm 1.11.1, rule id **`scale-unlimited/declaration-strict-value`**). "Specify properties for which a variable, function, keyword or value must be used." It enforces `$sass` / `namespace.$sass` / `@less` / `var(--custom-prop)` / css-loader `@value` / functions / keyword allowlists per property, with `ignoreValues` (e.g. `inherit`, `transparent`, `currentColor`) and a custom `message`.
- Source: <https://github.com/AndyOGo/stylelint-declaration-strict-value>, <https://www.npmjs.com/package/stylelint-declaration-strict-value>, listed in <https://stylelint.io/awesome-stylelint/>.
- Real-world precedent: Mozilla ships `stylelint-plugin-mozilla` rule **`no-base-design-tokens`** — <https://firefox-source-docs.mozilla.org/code-quality/lint/linters/stylelint-plugin-mozilla/rules/no-base-design-tokens.html>.
- **[INFERENCE]** This one rule is the *single highest-leverage* mechanical design-system gate that exists. "Every `color`, `background-color`, `border-color`, `box-shadow`, `font-family`, `font-size`, `padding`, `margin`, `gap`, `border-radius` MUST be `var(--token)`" collapses maybe half of all AI design slop into one deterministic, zero-FP rule.

### 2.3 CSS analytics / thresholds — **[FACT, reproduced locally]**
`@projectwallace/css-analyzer` **9.9.3** parses a stylesheet and returns metrics. Documented metric labels (source: <https://www.projectwallace.com/docs/metrics>): "Total colors", "Total unique colors", "Total color duplicates", "Total font-families / Total unique font-families", "Total font-sizes / Total unique font-sizes", "Total z-indexes / Total unique z-indexes", "Total `!important` declarations", "Ratio of `!important` declarations", "Maximum selector specificity", "Total selectors having maximum specificity", "Total box-shadows / unique", "Total text-shadows / unique", "Total animation durations / timing functions", "Total `@keyframes` rules", "Total `@media` atrules", "Total vendor prefixed properties/values", "Total properties / unique".

**Local probe** — I fed it an 8-rule "slop" stylesheet and it returned, correctly and instantly (no browser):
```json
{ "uniqueColors": ["#5b3df5","#6366f1","#a855f7","transparent","#fff","#0a0a0a","rgba(0,0,0,.2)","#5c3ef6"],
  "colorTotal": 13, "colorUniqueCount": 8,
  "fontFamilies": ["Inter,sans-serif","Poppins,sans-serif","\"Open Sans\",sans-serif"],
  "fontSizes": ["64px"],
  "zindex": ["10","9999","2147483647"],
  "important": { "total": 1, "ratio": 0.0294 },
  "maxSpecificity": [0,2,0],
  "units": ["px","vh","deg"] }
```
Notes: it **case-folds** (`#5B3DF5` → `#5b3df5`) so exact duplicates dedupe automatically, but it does **not** collapse *perceptually* identical colors (`#5b3df5` vs `#5c3ef6` stayed separate).
- **`constyble`** (npm 1.3.0) is the companion **CSS complexity linter** that turns those metrics into CI thresholds; **`wallace-cli`** is the CLI. Sources: <https://github.com/projectwallace/css-analyzer>, <https://github.com/projectwallace/wallace-cli>, <https://www.projectwallace.com/blog/new-online-css-code-quality-analyzer>.

### 2.4 Near-duplicate color detection — **[FACT, computed locally with culori 4.0.2]**
CIEDE2000 (`differenceCiede2000`) and OKLCH lightness deltas for representative pairs:

| pair | ΔE2000 | Δ L (OKLCH) | Δ chroma | verdict |
|---|---|---|---|---|
| `#5b3df5` vs `#5c3ef6` | **0.31** | 0.29 % | 0.0001 | duplicate token — hard fail |
| `#ffffff` vs `#fefefe` | **0.20** | 0.30 % | 0.0000 | duplicate token — hard fail |
| `#0a0a0a` vs `#000000` | **1.59** | 14.48 % | 0.0000 | duplicate-ish — warn |
| `#f8fafc` vs `#f1f5f9` (Tailwind slate-50/100) | **1.59** | 1.59 % | 0.0034 | **legitimate** adjacent scale step |
| `#111827` vs `#1f2937` (Tailwind gray-900/800) | **5.74** | 6.80 % | 0.0022 | legitimate |
| `#6366f1` vs `#818cf8` (indigo-500/400) | **11.90** | 9.47 % | 0.0458 | legitimate |

> **[INFERENCE] Recommended thresholds:** ΔE2000 **< 1.0 = ERROR** ("two tokens that are the same color — collapse them"); **1.0 ≤ ΔE2000 < 2.3 = NEEDS-REVIEW** (2.3 is the conventional just-noticeable-difference). Do **not** hard-fail above 1.0 — Tailwind's own slate-50/100 sits at 1.59 and is deliberate. Note the `#0a0a0a`/`#000000` row: ΔE2000 says 1.59 but OKLCH L differs by 14 %, so **use ΔE2000, not OKLCH L, for the duplicate test** (OKLCH L is non-linear near black).

### 2.5 Tailwind-specific static checks — **[FACT]**
- `eslint-plugin-tailwindcss` (npm 4.2.0) rules: **`classnames-order`**, **`enforces-shorthand`** (merge `mx-5 my-5` → `m-5`), **`no-arbitrary-value`** ("forbids using arbitrary values in classnames"), **`no-custom-classname`** ("detects classnames which do not belong to Tailwind CSS"), plus `no-contradicting-classname`, `enforces-negative-arbitrary-values`, `migration-from-tailwind-2`. Latest line is v4-targeted with **8 rules**; v3 users pin `@3.x.x`. Tailwind v4 support was a long-running gap (issue #325); alternatives exist (`@poupe/eslint-plugin-tailwindcss`, `oxlint-tailwindcss`).
- Sources: <https://github.com/francoismassart/eslint-plugin-tailwindcss>, <https://www.npmjs.com/package/eslint-plugin-tailwindcss>.
- **[INFERENCE]** `no-arbitrary-value` is *the* Tailwind equivalent of "off-scale spacing" — `p-[17px]` is the exact tell Hallmark gate 24 describes. But a **plain regex** `class="[^"]*\[[0-9.]+(px|rem|em|%)\]"` catches ~all of it with **zero node dependency**, works in `.html`, `.jsx`, `.vue`, `.svelte`, `.astro`, `.templ`, Go templates — and is vendor-neutral. Prefer the regex; offer the ESLint plugin as an optional adapter.

### 2.6 What you **cannot** get from CSS text alone — **[INFERENCE, but mechanically necessary]**
- **Contrast against the *computed* background.** A stylesheet says `color: var(--muted)`; only layout knows which of five stacked surfaces is behind it. This is precisely the "ink-on-ink" bug class (Hallmark gate 41). Text analysis produces both false positives and false negatives here.
- **Whether a rule ever matches an element.** Dead CSS, overridden focus rings, `:focus-visible` defined but out-specificity'd.
- **Line length in `ch`** when `max-width` is in `px` and the font is variable.
- **Actual accent-area footprint** (Hallmark gate 23, "accent > 5 % of viewport") — needs a rendered pixel count.
- **Whether the hero fits the fold** at 1280×800 — needs layout.

---

## 3. Visual / pixel checks

### 3.1 Playwright `toHaveScreenshot()` — **[FACT]**
Exact defaults (source: <https://playwright.dev/docs/api/class-pageassertions>, <https://playwright.dev/docs/test-snapshots>):

| option | default |
|---|---|
| `threshold` | **`0.2`** — "acceptable perceived color difference in the **YIQ color space**", 0 = strict, 1 = lax |
| `animations` | **`"disabled"`** |
| `caret` | **`"hide"`** |
| `scale` | `"css"` |
| `maskColor` | `#FF00FF` |
| `fullPage` | `false` |
| `omitBackground` | `false` |
| `maxDiffPixels` / `maxDiffPixelRatio` | unset |
| `stylePath` | unset (inject CSS to neutralise dynamic content) |

**[FACT]** Playwright's own warning, verbatim: *"Browser rendering can vary based on the host OS, version, settings, hardware, power source (battery vs. power adapter), headless mode, and other factors. For consistent screenshots, run tests in the same environment where the baseline screenshots were generated."* Update baselines with `npx playwright test --update-snapshots`.

### 3.2 Flakiness — **[FACT, secondary sources]**
- Documented flake causes: "font loading races, animations caught mid-transition, cursor blink state, lazy-loaded images that haven't finished rendering, and third-party widgets that load asynchronously"; "a screenshot taken on a developer's Mac will often look different than one taken on a Linux CI server due to font rendering and GPU differences." Source: <https://qtrl.ai/blog/visual-regression-testing-with-ai-2026>.
- **Percy** shipped an AI "Visual Review Agent" in late 2025 that "automatically filter[s] out 40 % of false positives, including anti-aliasing differences, sub-pixel rendering shifts, and operating system font variations."
- **Chromatic** relies on **TurboSnap** (only re-tests components changed between commits) plus "a built-in anti-flakiness layer that filters latency, animations, and resource-loading variability."
- **BackstopJS** (npm 6.3.25) is the OSS local option — no service, but you own baseline storage and the OS-drift problem.
- Sources: <https://qtrl.ai/blog/visual-regression-testing-with-ai-2026>, <https://saucelabs.com/resources/blog/comparing-the-20-best-visual-testing-tools-of-2026>, <https://crosscheck.cloud/blogs/best-visual-regression-testing-tools-2026/>.

> **[INFERENCE] Verdict for `parley-design-check`:** **do not ship pixel-diff visual regression.** It requires (a) a baseline store, (b) a pinned OS/container, (c) a human to adjudicate every diff, (d) an SaaS account for the good versions. It is a *regression* tool — it answers "did this change?", never "is this good?" — which is the wrong question for a design-quality gate on greenfield UI. **Do ship the deterministic geometry checks that use the same browser but produce a boolean, not an image**: horizontal overflow, element overlap, text clipping/truncation, target size, fold fit.

### 3.3 The pixel checks that *are* worth it — **[FACT, reproduced locally]**
Deterministic, boolean, no baseline needed:
- **Horizontal overflow** — one line of JS, no image: `document.documentElement.scrollWidth > document.documentElement.clientWidth`. My probe at a 375 px viewport returned `{"scrollWidth":1400,"clientWidth":375}` → unambiguous fail. This is Hallmark gate 34 (320–1920 px sweep), fully mechanised.
- **Contrast on gradient/image backgrounds** — the one place a *screenshot* genuinely beats the DOM. axe explicitly gives up ("Element's background color could not be determined due to a background gradient"). Deque notes "This problem is better handled with screenshot testing … not part of core axe." **[INFERENCE]** Sampling N pixels behind the text glyph bounding box from a screenshot and computing worst-case ratio is feasible but is a genuine research project; ship as `later`.
- **Text truncation / clipping** — `el.scrollWidth > el.clientWidth` on elements with `overflow: hidden` or `text-overflow: ellipsis`; **[INFERENCE]** high FP risk on intentionally-truncated table cells, so ship as needs-review not error.
- **Two-line clickable text** (Hallmark gate 49) — `el.getClientRects().length > 1` on `a`/`button`/nav items. **[INFERENCE]** deterministic, low FP, genuinely a slop tell.

---

## 4. Design-system conformance — how teams mechanically enforce "use the system"

**[FACT]** Named, real mechanisms:
1. **`stylelint-declaration-strict-value` / `scale-unlimited/declaration-strict-value`** — per-property "must be a variable/function/keyword" (see §2.2).
2. **Mozilla `stylelint-plugin-mozilla` → `no-base-design-tokens`** — bans raw/base tokens where semantic tokens are required. <https://firefox-source-docs.mozilla.org/code-quality/lint/linters/stylelint-plugin-mozilla/rules/no-base-design-tokens.html>
3. **ESLint core `no-restricted-imports`** — redirect `@mui/material/Button` → your wrapper; enforce module boundaries. <https://eslint.org/docs/latest/rules/no-restricted-imports>
4. **`vue/no-restricted-html-elements`** (eslint-plugin-vue) — ban raw `<a>` in favour of `<NuxtLink>`. <https://eslint.vuejs.org/rules/no-restricted-html-elements>
5. **Custom ESLint rules per design system** — the documented pattern ("Translating Your Design System Best Practices to ESLint", divRIOTS/Backlight). <https://backlight.dev/blog/best-practices-w-eslint-part-1>
6. **Stylelint `declaration-property-value-allowed-list`** — the vendor-neutral version of #1/#2, no plugin.
7. **DTCG / W3C Design Tokens Format Module** — **[FACT] first STABLE version `2025.10`, announced 28 October 2025** by the W3C Design Tokens Community Group; >20 editors/authors, contributors from Adobe, Google, Microsoft, Meta, Salesforce, Shopify, Figma, Framer, Tokens Studio, Penpot, Supernova, zeroheight. Sources: <https://www.w3.org/community/design-tokens/2025/10/28/design-tokens-specification-reaches-first-stable-version/>, spec draft <https://www.designtokens.org/tr/drafts/format/>, Style Dictionary support <https://styledictionary.com/info/dtcg/> (npm `style-dictionary` 5.5.0).
8. **Figma Code Connect** (npm `@figma/code-connect` 1.4.9) — maps repo components to Figma components so Dev Mode shows production snippets, "useful for when you have an existing design system and are looking to drive consistent and correct adoption". <https://github.com/figma/code-connect>, <https://help.figma.com/hc/en-us/articles/23920389749655-Code-Connect>
9. **Figma Variables REST API** — programmatic query/create/update/delete of variables; documented CI use for two-way token sync. <https://developers.figma.com/docs/rest-api/variables>

> **[INFERENCE] The single most important design decision for `parley-design-check`:** **the DTCG 2025.10 format is now stable and vendor-neutral — make it the canonical input artifact.** `parley-design` produces a `DESIGN-SYSTEM.tokens.json` in DTCG format; `parley-design-check` reads it and *derives* every scale check from it (allowed colors, allowed spacing, allowed font sizes, allowed radii). That is the protocol move: **the check has no hardcoded design opinions — it has a typed contract**. This is exactly the AG-UI/CopilotKit "typed canonical artifact + versioned spec" shape the owner asked for. Figma is then an *optional adapter*, not a dependency.
>
> **[INFERENCE] Do NOT depend on Figma.** Figma Code Connect and the Variables REST API both require an account, a token, and network. That kills portability and makes the skill un-runnable in the headless-CLI Parley environment.

---

## 5. Runtime / dependency cost — the engine tiers

**[FACT — measured locally]** headless Chromium launch + `file://` load + axe injection + full `axe.run()` on a trivial page: **launch+load 406 ms, total 537 ms** (Playwright `chromium_headless_shell-1228`, macOS arm64, Node v26.5.0). A second run with two rules: same order of magnitude. So the browser tier is **sub-second per page** once the binary exists — the cost is the **~150 MB binary + install step**, not the runtime.

| tier | engine | needs | typical cost | what it unlocks |
|---|---|---|---|---|
| **T0** | regex / line scan | **nothing** (any language: Go, POSIX shell, Python, Node) | µs–ms per file | hex-vs-token, `!important`, `transition: all`, `hover:scale-105`, arbitrary Tailwind `[17px]`, banned font names, emoji-as-icon, gradient-text, `min-height:100vh` hero, `100vh` on mobile |
| **T1** | CSS parse (AST) | a CSS parser. Node: `postcss` / `@projectwallace/css-analyzer` / `stylelint`. **Go alternative exists** (`tdewolff/parse/css`), Rust (`lightningcss`) | ms per stylesheet | token conformance per property, off-scale values, specificity, z-index set, unit sprawl, font-family count, near-duplicate colors (ΔE), duplicate selectors |
| **T2** | static DOM/JSX parse | HTML parser (Node `cheerio`/`jsdom`; **Go `golang.org/x/net/html`**) or an ESLint/Babel AST | ms per file | heading order, `alt`, `lang`, decorative `<svg>` without `aria-hidden`/`aria-label`, duplicate `id`, `tabindex` > 0, missing `<label for>`. **NO contrast, NO geometry** (jsdom has no layout — axe issue #595) |
| **T3** | headless browser | **Chromium ~150 MB**, Node (Playwright/Puppeteer) or CDP from any language. No network if `file://` | ~0.4 s launch + ~0.1 s/page (**measured**) | computed-style contrast, target size, horizontal overflow @320–1920 px, `:focus-visible` computed appearance, reflow, fold fit, two-line clickable text, full axe-core (96 default rules) |
| **T4** | screenshot pipeline | T3 + image decode/diff (`pixelmatch`/`odiff`) + **baseline store** + pinned OS/container; SaaS variants need **network + account** | seconds/page + human triage | gradient/image-background contrast, visual regression, accent-area % footprint |

**[INFERENCE] Portability ranking for a vendor-neutral Parley skill:** T0 and T1 are shippable as a **single self-contained binary or script with zero install**. T2 is nearly free. T3 is the first real dependency — but it is *already present* in most modern web repos (Playwright) and Parley agents run on dev machines that have it. T4 is where portability dies.

---

## 6. THE TABLE — check inventory for `parley-design-check`

`Tier` = §5. `FP` = false-positive risk. `Ship` = yes / no / later.

| # | Check | What it catches | Tier | FP risk | Ship |
|---|---|---|---|---|---|
| 1 | **Raw color literal outside token block** (`#hex`, `rgb()`, `hsl()`, `oklch()` not in `:root`/tokens file) | mid-render token improvisation (Hallmark gate 48); the #1 system-violation | T0/T1 | **very low** | **yes** |
| 2 | **Off-scale spacing** — any `padding`/`margin`/`gap`/`top`…/`border-radius` value not in the declared scale (`p-[17px]`, `padding: 23px`) | arbitrary-value tell (Hallmark gate 24) | T0/T1 | low (needs `ignoreValues`: `0`, `auto`, `100%`, `1px`) | **yes** |
| 3 | **Font-family count > 3** / banned default stacks (Inter, Roboto, Open Sans, Poppins, Lato, system-ui as *display*) | typography sprawl + the single loudest AI tell (Hallmark gates 1, 37) | T0/T1 | low; "banned font" list is opinionated → make it config, default warn | **yes** |
| 4 | **Near-duplicate colors, ΔE2000 < 1.0 error / < 2.3 needs-review** | two tokens that are the same color; `#5b3df5` vs `#5c3ef6` measured at ΔE 0.31 | T1 | low at <1.0; **medium** 1.0–2.3 (Tailwind slate-50/100 = 1.59) | **yes** (error <1.0, review 1.0–2.3) |
| 5 | **`!important` count / ratio** above threshold | cascade capitulation; Wallace reports `importants.total` + `.ratio` | T0/T1 | very low | **yes** |
| 6 | **z-index sprawl** — count of distinct values, any value > 100, any `2147483647` | z-index chaos; measured `["10","9999","2147483647"]` in a 8-rule sheet | T1 | low | **yes** |
| 7 | **`transition: all` / `transition-all`** | Hallmark gate 10, universal AI tell | T0 | very low | **yes** |
| 8 | **Uniform `hover:scale-105` / `transform: scale(1.05)` on ≥3 unrelated selectors** | Hallmark gate 11 | T0/T1 | low-medium (legit on cards) → warn | **yes (warn)** |
| 9 | **Animating `width`/`height`/`top`/`left`/`margin`/`padding`** | non-compositable animation (Hallmark gate 14) | T1 | low | **yes** |
| 10 | **Motion without `prefers-reduced-motion` fallback** | Hallmark gate 27; also an a11y obligation | T1 | low | **yes** |
| 11 | **Gradient text** (`background-clip: text` + `linear-gradient`) | Hallmark gate 2; also an unmeasurable-contrast surface | T0 | very low | **yes** |
| 12 | **Missing `:focus-visible` on interactive selectors** | SC 2.4.7; Hallmark gate 26 | T1 (declared) / T3 (effective) | **medium at T1** (may be defined in a shared reset) | **yes at T1 as needs-review; yes at T3 as error** |
| 13 | **Prose `max-width` outside 45–75ch** | line-length (Hallmark gate 25) | T1 if `ch`; **T3 if `px`** | medium | **yes (warn)** |
| 14 | **All-caps display with `line-height` < 1.0** | cap-collision on wrap (Hallmark gate 55) | T1 | low | **yes** |
| 15 | **Grid `1fr` track containing an image, not `minmax(0,1fr)`** | mobile blowout (Hallmark gate 50) | T1+T2 (needs to know a child is an img) | medium | **later** |
| 16 | **axe-core `color-contrast`** against computed backgrounds (4.5:1 / 3:1 large) | ink-on-ink, black-on-black; measured `1.01` ratio caught | T3 | **very low** (axe fails *open* on gradients → `incomplete`) | **yes — the flagship check** |
| 17 | **axe-core full WCAG 2.0/2.1 AA run** (96 default rules) | alt text, labels, landmarks, aria, contrast | T3 | very low | **yes** |
| 18 | **Own 24×24 CSS px target-size rule** (no spacing exception) | small tap targets — axe's `target-size` **passed** an 18×18 button via the Spacing exception | T3 | low-medium (icon buttons in dense toolbars) | **yes (needs-review), configurable strict** |
| 19 | **Horizontal overflow sweep 320/375/414/768/1280/1920 px** | Hallmark gate 34; measured `scrollWidth 1400 > clientWidth 375` | T3 | **very low** | **yes** |
| 20 | **Two-line clickable text** (`getClientRects().length > 1` on `a`/`button`/nav) | Hallmark gate 49 | T3 | low | **yes** |
| 21 | **Reflow @ 320 px / 400 % zoom (SC 1.4.10)** | two-dimensional scrolling | T3 | low | **yes** |
| 22 | **Text-spacing stress test (SC 1.4.12)**: inject `line-height:1.5;letter-spacing:.12em;word-spacing:.16em` then re-check overflow/clipping | brittle fixed-height layouts | T3 | medium | **later** |
| 23 | **Focus-ring appearance ≥ 2 px perimeter & ≥ 3:1 (SC 2.4.13 AAA)** | invisible focus | T3 | medium (outline vs box-shadow vs border) | **later** |
| 24 | **Text truncation / clipping** (`scrollWidth > clientWidth` on overflow-hidden) | clipped labels | T3 | **high** (intentional ellipsis) | **later (needs-review only)** |
| 25 | **Element overlap / collision detection** (bounding-box intersection of siblings) | two sticky-at-top-0 elements (Hallmark gate 56) | T3 | medium-high | **later** |
| 26 | **Accent-area footprint ≤ ~5 % of viewport** (Hallmark gate 23) | accent flooding | **T4** (needs pixel count) | high | **no** |
| 27 | **Contrast over gradient / image backgrounds** | the exact case axe declares `incomplete` | **T4** | medium | **later** |
| 28 | **Pixel visual regression (Percy/Chromatic/BackstopJS/Playwright `toHaveScreenshot`)** | *change*, not *quality*; `threshold` default 0.2 YIQ | **T4** | **very high** (OS fonts, GPU, antialiasing, animation frames) | **no** |
| 29 | **DTCG token file schema validation** (2025.10) | malformed / undeclared tokens; makes every scale check derivable | T0 (JSON) | very low | **yes — the contract** |
| 30 | **Unused / undeclared token cross-check** (tokens declared but never used; `var(--x)` used but never declared) | dead system, phantom tokens | T1 + T0 | low | **yes** |
| 31 | **Structural fingerprint: hero → 3-col icon cards → CTA → 4-col footer** | the generic AI template (Hallmark gate 8, 42, 43) | T2 (heuristic) | **high** — a 3-col feature grid is often correct | **no as a gate; yes as a `needs-review` signal for the agents** |
| 32 | **Icon-library mixing** (Material + Heroicons + Lucide imports on one page) / emoji-as-feature-icon | Hallmark gate 30 | T0/T2 | low | **yes** |
| 33 | **Re-drawn UI chrome** (fake browser bar / phone notch / terminal frame in HTML/CSS) | Hallmark gate 47 | T2 heuristic | high | **no — agent judgement** |
| 34 | **Invented metrics** ("10× faster", "trusted by 50,000+") not present in the brief | Hallmark gate 46 | T0 regex + **LLM judgement** | high without the brief | **no as a gate; yes as an agent prompt** |
| 35 | **Taste / hierarchy / originality / brief-fit** | the actual design quality | **none** | — | **no — this is `parley-design`'s job** |

---

## Transferable to parley-design / parley-design-check

Ranked, most valuable first.

1. **[→ design-check] Make `DESIGN-SYSTEM.tokens.json` in W3C DTCG `2025.10` format the typed canonical artifact.** The check ships **no hardcoded scales** — it derives allowed colors/spacing/type/radii from the token file. This is the protocol move the owner wants (typed artifact + versioned spec + machine-checkable). Stable since 28 Oct 2025, backed by Adobe/Google/Microsoft/Figma/Salesforce. `parley-design`'s Phase output = the token file; `parley-design-check`'s Phase input = the same file. (Sources: §4.7)
2. **[→ design-check] Adopt IBM Equal Access's tri-state verdict: `violation` / `needs-review` / `recommendation`.** A binary gate on design work is a lie. Map: `violation` blocks the Parley phase gate; `needs-review` is routed to the adversarial agents as a required consult; `recommendation` is advisory. (Source: §1.4 IBM)
3. **[→ design-check] Ship checks 1–11 + 29–30 as Tier 0/1 with zero runtime dependency.** These are ~60 % of the value at ~0 % of the cost, run anywhere (Go, shell, Python, Node), and never need Chromium. This is what makes the skill *portable*, which is the stated design goal.
4. **[→ design-check] Make Tier 3 (headless Chromium + axe-core) an OPTIONAL, auto-detected stage** with a clean degradation message. Measured cost is **537 ms/page** once the binary exists. Checks 16–21 are the ones that catch bugs Tier 0/1 provably cannot (ink-on-ink, overflow, target size).
5. **[→ design-check] Steal pa11y's contract shape**: `--threshold <n>` + **exit code 2** on breach. Simple, scriptable, vendor-neutral, and exactly what a Parley phase gate needs.
6. **[→ design-check] Write your OWN target-size rule at 24×24 CSS px.** Empirically, axe's `target-size` passes an 18×18 button via the WCAG Spacing exception. If the doctrine says "44 px floor" (Hallmark gate 39 does), the check must implement it, not delegate.
7. **[→ design-check] Use ΔE2000 for near-duplicate colors with the measured thresholds: <1.0 error, 1.0–2.3 needs-review.** Do not use OKLCH lightness delta (`#0a0a0a` vs `#000000` = 14 % ΔL but ΔE2000 1.59).
8. **[→ design] Write the doctrine so every rule states its enforcement tier.** Each doctrine rule in `parley-design` gets a machine tag: `enforced-by: check#N` / `enforced-by: agent-judgement`. That single annotation is what keeps the two skills honest and prevents `parley-design-check` from over-claiming.
9. **[→ design] Cite the coverage numbers correctly and prominently.** "Automated tools cover ~30–40 % of WCAG success criteria (GDS: best tool 40 %, worst 13 %, of 142 seeded barriers) but up to ~57 % of issue *volume* (Deque, 2,000+ audits / ~300k issues)." This is the sentence that justifies why `parley-design` (agents) exists at all.
10. **[→ design-check] Reuse Hallmark's slop-test as the doctrine source, but re-classify each of its 58 gates by tier.** File: `/Users/tomasfecko/.claude/skills/hallmark/references/slop-test.md`. Roughly: gates 10, 14, 22, 24, 25, 27, 30, 37, 48, 55 are T0/T1 mechanical; gates 26, 34, 39, 40, 41, 44, 49, 56 need T3; gates 6, 8, 20, 21, 32, 35, 45, 46, 47, 54, 57 are agent-judgement only. **[INFERENCE]**
11. **[→ design-check] Gate on `axe.run()` directly, never on the Lighthouse score.** Lighthouse *is* axe underneath, weighted and subsetted, pass/fail per audit with no partial credit — strictly weaker.
12. **[→ design-check] Emit `needs-review` (never `pass`) whenever axe returns `incomplete`.** Verbatim reproduced case: `Element's background color could not be determined due to a background gradient`. A gradient hero headline currently reports as neither pass nor fail — silently treating that as pass is how the ink-on-ink bug ships.
13. **[→ design-check] Where a Node ecosystem exists, delegate rather than reimplement**, via optional adapters with pinned rule ids: `stylelint` core (149 rules — `color-no-hex`, `declaration-no-important`, `unit-allowed-list`, `selector-max-specificity`, `font-family-no-missing-generic-family-keyword`, `declaration-property-value-allowed-list`, …), `scale-unlimited/declaration-strict-value`, `eslint-plugin-tailwindcss` (`no-arbitrary-value`, `no-custom-classname`, `enforces-shorthand`, `classnames-order`), `@projectwallace/css-analyzer` + `constyble`.
14. **[→ design] Adopt the DIVERGE→critique→one-wins framing explicitly in the artifact schema**, because *nothing in §6 can evaluate direction*. The table's bottom rows (31, 33, 34, 35) are the honest boundary: taste is un-mechanisable, so the protocol must name the human/agent step that owns it.

---

## Do NOT copy

1. **Pixel-diff visual regression (Percy / Chromatic / BackstopJS / Playwright `toHaveScreenshot`) as a design-quality gate.** *Reason:* it answers "did this change vs a baseline?", not "is this good?" — useless on greenfield UI where there is no baseline. Plus Playwright's own docs warn rendering "can vary based on the host OS, version, settings, hardware, power source (battery vs. power adapter), headless mode"; Percy needed an AI agent in 2025 just to filter out 40 % false positives from antialiasing and OS font variation. It also needs a baseline store, a pinned container, and (for the good tools) network + a paid account — all fatal to a portable, vendor-neutral skill.
2. **APCA Lc as a *gate* threshold.** *Reason:* APCA was pulled from the July 2023 WCAG 3 draft and is absent from the March 2026 draft; it "was only ever exploratory". Gate on WCAG 2.2's 4.5:1 / 3:1 / 3:1. APCA may appear in `parley-design` as a design-time heuristic, clearly labelled non-normative.
3. **WCAG 3.0 requirements / its scoring model.** *Reason:* Working Draft only (3 Mar 2026, ~174 requirements), CR anticipated ~Q4 2027, REC not before 2028. Building a gate on it means rewriting the gate twice.
4. **The Lighthouse accessibility score.** *Reason:* it is a weighted subset of axe with binary per-audit scoring and no partial credit; a 100 is compatible with real failures. Use axe directly.
5. **Any hard dependency on Figma** (Code Connect, Variables REST API, design-file diffing). *Reason:* requires an account, a personal access token, and network. Kills portability and cannot run in the headless-CLI Parley environment. Optional adapter at most.
6. **jsdom-based "accessibility checking" as a substitute for a browser.** *Reason:* jsdom does no layout or rendering — `color-contrast` "doesn't work in JSDOM" (axe issue #595) and returns `incomplete`, and `jest-axe` disables it outright. A jsdom pass that reports zero contrast violations is a **false all-clear**.
7. **The unqualified "57 % of accessibility issues are automatable" claim.** *Reason:* it measures issue *volume* on first-pass audits, not criteria coverage; quoting it alone implies automation covers a majority of *kinds* of problems, which the GDS 142-barrier study contradicts (13–40 % per tool).
8. **Framework-coupled linters as the core engine** (`eslint-plugin-jsx-a11y`, `eslint-plugin-tailwindcss`, `vue/no-restricted-html-elements`). *Reason:* the skills are explicitly vendor-neutral; these bind you to React/Tailwind/Vue. Also `eslint-plugin-jsx-a11y`'s own README: "Static analysis tools cannot determine values of variables that are being placed in props before runtime." Ship them as optional adapters behind a capability probe.
9. **Structural-fingerprint checks as hard gates** (check 31: hero → 3-col icon cards → CTA → 4-col footer; check 33: re-drawn chrome). *Reason:* very high false-positive rate — a 3-column feature grid is frequently the correct answer, and a legitimate product page may embed a real device frame. Emit as a `needs-review` signal that the adversarial agents must address in writing; do not block.
10. **Hallmark's 58-gate list copied verbatim.** *Reason:* it is (a) genre- and theme-coupled (Specimen / Midnight / Riso macrostructures, `.hallmark/log.json`, CSS stamp format), (b) opinionated about specific archetypes (N1b, Ft5, HP3), and (c) mixes mechanical and taste gates without marking which is which. Take the *rules*, re-derive them against the DTCG token contract, and tag each with its enforcement tier.
11. **`constyble` / Wallace thresholds as absolute numbers.** *Reason:* metrics like "unique colors ≤ 12" are only meaningful relative to a project's declared token set. Derive thresholds **from the DTCG file**, not from a global default, or the check becomes noise on any large codebase.
12. **A single binary pass/fail verdict.** *Reason:* see transferable #2 — every serious engine (IBM, axe, Deque) ships a needs-review state precisely because the pass/fail boundary is where automated design checking lies to you.

---

## Appendix — verified tool versions (npm, checked 2026-07-28)

| package | version |
|---|---|
| `axe-core` | 4.12.1 |
| `stylelint` | 17.14.1 (149 core rules) |
| `@projectwallace/css-analyzer` | 9.9.3 |
| `constyble` | 1.3.0 |
| `eslint-plugin-jsx-a11y` | 6.10.2 |
| `pa11y` | 9.1.1 |
| `accessibility-checker` (IBM) | 4.0.29 |
| `eslint-plugin-tailwindcss` | 4.2.0 |
| `stylelint-declaration-strict-value` | 1.11.1 |
| `backstopjs` | 6.3.25 |
| `@figma/code-connect` | 1.4.9 |
| `style-dictionary` | 5.5.0 |
| `culori` | 4.0.2 |

Local environment: Node v26.5.0, npm 11.17.0, Python 3.14.6, Go 1.26.5 (darwin/arm64); Playwright browsers cached at `~/Library/Caches/ms-playwright/` (`chromium-1217`, `chromium-1228`, `chromium_headless_shell-1217/1228`).

## Appendix — primary sources

- WCAG 2.2 — <https://www.w3.org/TR/WCAG22/>
- Understanding 2.5.8 Target Size (Minimum) — <https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html>
- Deque 57 % study — <https://www.deque.com/blog/automated-testing-study-identifies-57-percent-of-digital-accessibility-issues/>
- GDS accessibility tools audit (142 barriers, 13 tools) — <https://alphagov.github.io/accessibility-tool-audit/results.html>
- axe-core rule descriptions — <https://github.com/dequelabs/axe-core/blob/develop/doc/rule-descriptions.md>
- axe-core accessibility-supported policy — <https://github.com/dequelabs/axe-core/blob/develop/doc/accessibility-supported.md>
- axe-core jsdom contrast limitation — <https://github.com/dequelabs/axe-core/issues/595>
- axe-core background-image contrast limitation — <https://github.com/dequelabs/axe-core/issues/3390>
- Lighthouse accessibility scoring — <https://developer.chrome.com/docs/lighthouse/accessibility/scoring>
- pa11y — <https://github.com/pa11y/pa11y>
- IBM Equal Access — <https://github.com/IBMa/equal-access/tree/master/accessibility-checker>, <https://www.ibm.com/able/toolkit/verify/automated>
- Playwright visual comparisons — <https://playwright.dev/docs/test-snapshots>, <https://playwright.dev/docs/api/class-pageassertions>
- Project Wallace metrics — <https://www.projectwallace.com/docs/metrics>
- stylelint rules — <https://stylelint.io/awesome-stylelint/>
- stylelint-declaration-strict-value — <https://github.com/AndyOGo/stylelint-declaration-strict-value>
- Mozilla `no-base-design-tokens` — <https://firefox-source-docs.mozilla.org/code-quality/lint/linters/stylelint-plugin-mozilla/rules/no-base-design-tokens.html>
- DTCG stable spec announcement — <https://www.w3.org/community/design-tokens/2025/10/28/design-tokens-specification-reaches-first-stable-version/>
- DTCG Format Module draft — <https://www.designtokens.org/tr/drafts/format/>
- Figma Code Connect — <https://github.com/figma/code-connect>
- Figma Variables REST API — <https://developers.figma.com/docs/rest-api/variables>
- eslint-plugin-jsx-a11y — <https://github.com/jsx-eslint/eslint-plugin-jsx-a11y>
- eslint-plugin-tailwindcss — <https://github.com/francoismassart/eslint-plugin-tailwindcss>
- WCAG3 contrast status (Adrian Roselli, Apr 2026) — <https://adrianroselli.com/2026/04/wcag3-contrast-as-of-april-2026.html>
- Visual regression flakiness — <https://qtrl.ai/blog/visual-regression-testing-with-ai-2026>
- Local prior art: Hallmark slop test — `/Users/tomasfecko/.claude/skills/hallmark/references/slop-test.md`
