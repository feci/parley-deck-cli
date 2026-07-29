# Digest — The standard, professional way design systems are specified

**Purpose:** give `parley-design` a *documented, existing* canonical-artifact format so it is not invented from scratch, and give `parley-design-check` a *mechanically checkable* substrate.
**Researched:** 2026-07-28. All web facts re-verified live on that date unless marked INFERENCE.
**Method note:** everything under "FACT" was read from the primary source (spec text, published npm tarball, or repo source file) and the source is named inline. Everything under "INFERENCE" is my judgement and is labelled.

---

## 0. Executive framing (INFERENCE)

There are **three** independent layers that mature design systems specify, and they are specified in *different formats by different standards*. Conflating them is the #1 cause of vague design docs:

| Layer | What it is | Standard that exists TODAY | Machine-checkable? |
|---|---|---|---|
| **L1 Tokens** | the *data* (color, dimension, typography values + aliases + modes) | **W3C DTCG Design Tokens 2025.10** — a real, stable, JSON-Schema-backed spec | **YES — fully.** JSON Schema + linters |
| **L2 Component specs** | anatomy, states, variants, behaviors, a11y, content | **No cross-industry standard.** But Carbon/Polaris/GOV.UK publish near-identical *templates*, which is a de-facto standard | **Partially** — structure yes, prose no |
| **L3 Governance** | how something enters/leaves the system | GOV.UK contribution criteria + Carbon PDLC statuses are the two best-documented models | **YES** — status enum + checklist gates |

`parley-design` should adopt **DTCG verbatim for L1** (zero invention), a **merged Carbon/Polaris/GOV.UK template for L2** (light invention: pick sections), and a **Carbon-PDLC-style status enum + GOV.UK criteria for L3**.

---

## 1. FACT — W3C DTCG Design Tokens format

### 1.1 Status and version (verified 2026-07-28)

- The DTCG spec **reached its first stable version, `2025.10`, announced 28 Oct 2025**. Source: https://www.w3.org/community/design-tokens/2025/10/28/design-tokens-specification-reaches-first-stable-version/
- The stable release lives at **https://www.designtokens.org/tr/2025.10/** (HTTP 200, verified) and contains **three modules, all marked stable**: **Format**, **Color**, **Resolver**. The stable page states: *"This specification is considered stable. Further updates will be provided in superseding specifications."*
- The **living draft** is at https://www.designtokens.org/tr/drafts/format/ — title **"Design Tokens Format Module 2025.10"**, status **"Draft Community Group Report"**, dated **17 June 2026**, editors **Louis Chenais, Kathleen McMahon, Drew Powers, Matthew Ström-Awn, Donna Vitan**. The drafts carry the warning *"Do not attempt to implement this version of the specification."*
- ⚠️ **`tr.designtokens.org` 301-redirects to `www.designtokens.org/TR/drafts/`** — cite `www.designtokens.org/tr/2025.10/` for the stable text, not `tr.designtokens.org`.
- Versioning scheme is **date-based `YYYY.MM`**, not semver. No published cadence.

### 1.2 The exact JSON structure

Verified against https://www.designtokens.org/tr/drafts/format/ and the machine schema (below).

**A token** = a JSON object that **has a `$value`**. Its **parent key is its name**.

```json
{
  "Button background": {
    "$type": "color",
    "$description": "The background color for buttons in their normal state.",
    "$value": {
      "colorSpace": "srgb",
      "components": [0.467, 0.467, 0.467],
      "alpha": 1,
      "hex": "#777777"
    }
  }
}
```

Reserved `$`-prefixed properties:

| Property | Applies to | Meaning |
|---|---|---|
| `$value` | token only | **required**; presence of `$value` is what makes an object a token |
| `$type` | token or group | category string, **case-sensitive**; **inherits** token → parent group → up the hierarchy |
| `$description` | token or group | plain text; tools may render as tooltip/comment |
| `$extensions` | token or group | vendor data, **reverse-domain-notation keys** |
| `$deprecated` | token or group | boolean **or** string explanation |
| `$extends` | group | deep-merge inheritance, semantics equivalent to JSON Schema `$ref` |
| `$root` | group | a group may carry a **root token named `$root`** alongside children |

**Name restrictions (FACT, directly checkable):** token and group names **must not begin with `$`**, and the characters **`{`, `}`, `.`** are prohibited **anywhere** in a name.

**A group** = a JSON object **without** `$value`. Groups nest arbitrarily and propagate `$type`.

**Aliasing — two forms, both normative:**

1. Curly-brace, references a whole token value:
   `{ "$value": "{colors.blue}" }`
2. **JSON Pointer (RFC 6901), support required**, references *any* location incl. sub-properties:
   `{ "$ref": "#/colors/blue/$value/components/0" }`

**Token `$type` list (13):** `color`, `dimension`, `fontFamily`, `fontWeight`, `duration`, `cubicBezier`, `number`, `strokeStyle`, `border`, `transition`, `shadow`, `gradient`, `typography`.
- `dimension` = object `{value: number, unit: "px" | "rem"}` — **not a bare string**.
- `duration` = object `{value, unit: "ms" | "s"}`.
- `fontWeight` = number `[1..1000]` **or** a named string (`thin`…`black`).
- `cubicBezier` = array of exactly 4 numbers.
- `strokeStyle` / `border` / `transition` / `shadow` / `gradient` / `typography` are **composite** types.
- `shadow` may be a **single object or an array** (multi-shadow).

### 1.3 Color module (FACT)

Source: https://www.designtokens.org/tr/drafts/color/ (Draft CG Report, 17 June 2026).

`$value` object: **`colorSpace` (required)**, **`components` (required)**, `alpha` (optional), `hex` (optional).
**14 supported `colorSpace` values:** `srgb`, `srgb-linear`, `hsl`, `hwb`, `lab`, `lch`, `oklab`, `oklch`, `display-p3`, `a98-rgb`, `prophoto-rgb`, `rec2020`, `xyz-d65`, `xyz-d50`.
Spec wording on fallback: *"The fallback color MUST be formatted in 6 digit CSS hex color notation format"* (so alpha is never smuggled into `hex`).

**This is the single most important fact for `parley-design`:** the standard *already* supports **OKLCH natively as a first-class colorSpace**. You do not need to invent a perceptual-color convention.

### 1.4 Resolver module — modes/themes (FACT)

Source: https://www.designtokens.org/tr/drafts/resolver/ (stable as of 2025.10 per the TR index; the *drafts* copy is a preview).

Solves light/dark/high-contrast/density/reduced-motion **without combinatorial file explosion**. Canonical shape:

```json
{
  "$schema": "https://www.designtokens.org/schemas/2025.10/resolver.json",
  "version": "2025.10",
  "sets": {
    "foundation": { "sources": [{ "$ref": "foundation.json" }] }
  },
  "modifiers": {
    "theme": {
      "contexts": {
        "light": [{ "$ref": "themes/light.json" }],
        "dark":  [{ "$ref": "themes/dark.json" }]
      },
      "default": "light"
    }
  },
  "resolutionOrder": [
    { "$ref": "#/sets/foundation" },
    { "$ref": "#/modifiers/theme" }
  ]
}
```

Required root keys: **`version` (MUST be `"2025.10"`)**, **`resolutionOrder`**; optional `sets`, `modifiers`. A set *"MUST contain a `sources` array"*. Resolution = combine sets, then apply the selected modifier context, e.g. input `{ "theme": "dark" }`.

### 1.5 ⭐ The machine-checkable hook (FACT, verified by HTTP 200 + body inspection)

Official JSON Schemas are published and live:

- `https://www.designtokens.org/schemas/2025.10/format.json` → `{"$schema":"http://json-schema.org/draft-07/schema#","$id":"https://www.designtokens.org/schemas/2025.10/format.json","title":"DTCG Format Schema", …}`
- `https://www.designtokens.org/schemas/2025.10/resolver.json`
- sub-schemas exist, e.g. `…/2025.10/format/tokenType.json` (referenced from the root schema)

**Consequence for `parley-design-check`:** token-file validity is a **pure JSON-Schema assertion** — `ajv`/any validator, zero heuristics, zero LLM judgement. This is the hardest, least-arguable gate you can ship.

---

## 2. FACT — Style Dictionary / Terrazzo: the canonical token→platform pipeline

### 2.1 Style Dictionary (verified via `npm view`, 2026-07-28)

- **Current version: `5.5.0`**. **`engines: { node: ">=22.0.0" }`.** (`npm view style-dictionary version` / `engines`)
- Docs: https://styledictionary.com/info/architecture/

**Canonical pipeline, exact stage order (per platform):**

1. Parse config
2. Find token files (`include` + `source` globs)
3. Parse token files (custom parsers; JSON/JS automatic)
4. **Deep merge** into one object
5. **Preprocessors**
6. **Transforms** (per token; value transforms skipped for tokens that are references)
7. **Transitive transforms** (post-resolution)
8. **Resolve references** (`"{size.font.base}"` → value)
9. **Filters**
10. **Formats**
11. **File headers**
12. **Actions** (asset copy, image gen, …)

**Predefined transform groups (exact names):** `web`, `css`, `scss`, `less`, `js`, `json`, `html`, `assets`, `react-native`, `ios`, `ios-swift`, `ios-swift-separate`, `android`, `compose`, `flutter`, `flutter-separate`.
Examples of composition: `web` = `attribute/cti, name/kebab, size/px, color/css`; `js` = `attribute/cti, name/pascal, size/rem, color/hex`; `react-native` = `name/camel, size/object, color/css`.
Source: https://styledictionary.com/reference/hooks/transform-groups/predefined/

### 2.2 Terrazzo (verified via `npm view`, 2026-07-28)

- **`@terrazzo/cli` `2.5.0`**, binaries **`tz` / `terrazzo`**, self-described: *"CLI for managing design tokens using the Design Tokens Community Group (DTCG) standard and generating code for any platform via plugins."* Published 2 days before this digest (very live).
- Plugins that exist right now: `@terrazzo/plugin-css` 2.5.0, `@terrazzo/plugin-sass` 2.5.0, `@terrazzo/plugin-js` 2.5.0, `@terrazzo/plugin-tailwind` 2.5.0, `@terrazzo/plugin-swift` 0.3.3.
- ❌ **`@terrazzo/plugin-lint-a11y` and `@terrazzo/plugin-lint-core` do NOT exist as separate packages** (`npm view` → 404). Linting is **built into the CLI**, not a plugin. Do not write install instructions for those package names.

**⭐ Terrazzo ships 27 built-in lint rules** (https://terrazzo.app/docs/linting/) — this is the closest thing to an off-the-shelf "design slop detector" for tokens:

- **Type validity (16):** `core/valid-color`, `core/valid-dimension`, `core/valid-font-family`, `core/valid-font-weight`, `core/valid-duration`, `core/valid-cubic-bezier`, `core/valid-number`, `core/valid-link`, `core/valid-boolean`, `core/valid-string`, `core/valid-stroke-style`, `core/valid-border`, `core/valid-transition`, `core/valid-shadow`, `core/valid-gradient`, `core/valid-typography`
- **Structure/style (8):** `core/colorspace` (force e.g. oklch), `core/consistent-naming` (kebab/camel), `core/duplicate-values` (**← direct anti-slop: catches two tokens with the same value**), `core/descriptions` (**← forces intent to be written down**), `core/max-gamut`, `core/required-children`, `core/required-modes` (**← forces every token to have dark mode**), `core/required-type`
- **Accessibility (2):** **`a11y/min-contrast`** (WCAG 2.2 contrast between declared color pairs), **`a11y/min-font-size`**

Config shape:
```js
lint: { rules: { "a11y/min-contrast": ["error", { /* opts */ }] } }  // "error" | "warn" | "off"
```

### 2.3 Other verified enforcement tooling (npm, 2026-07-28)

| Package | Version | Use for `parley-design-check` |
|---|---|---|
| `stylelint-declaration-strict-value` | 1.11.1 | **ban raw hex/px in CSS; require `var(--token)`** — the single best "no magic numbers" gate |
| `@shopify/stylelint-polaris` | 16.0.7 | proof that a major system ships a *stylelint plugin* as its enforcement arm |
| `axe-core` | 4.12.1 | runtime a11y assertions |
| `eslint-plugin-jsx-a11y` | 6.10.2 | static a11y assertions |
| `pa11y` | 9.1.1 | CI a11y crawling |
| `@adobe/spectrum-tokens` | 14.15.0 | real-world token package to diff against |
| `@carbon/themes` | 11.77.0 | ditto |

---

## 3. FACT — Mature systems' documented methodology

### 3.1 IBM Carbon — *the single best model to copy*

Source read in full: https://raw.githubusercontent.com/carbon-design-system/carbon-website/main/src/pages/contributing/component-checklist/index.mdx (saved locally at `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/research/carbon-checklist.mdx`) and `.../contributing/documentation.mdx` (saved as `carbon-doc-templates.mdx`).

**(a) Explicit "Definition of done" + a 4-value status enum tied to a lifecycle (PDLC):**

| Status | PDLC Phase | Description (verbatim) |
|---|---|---|
| `Draft` | Discovery | "Partially complete, ready for validation." |
| `Preview candidate` | Discovery | "Partially complete, with measurable results, stakeholders, and clear business value." |
| `Preview` | Delivery | "Mostly complete, changes possible based on feedback, available to use in production." |
| `Stable` | Launch and scale | "Complete across code, kit, docs, design, and ready for production use." |

**(b) The "Design spec" is named as the contract.** Verbatim: *"The design specification (spec) is the blueprint used by developers to build the component in code and for designers making the component in Figma. It is referenced as the source of truth for the visual appearance and functionality of a component."*

Its **6 mandatory requirement rows** (each with "Details" + "Why this matters"):
1. **Color tokens** — *"Design specs only use color tokens available in the system."* / *"Design specs only contain colors that are tokenized."*
2. **Type tokens** — *"Design specs only use type tokens available in the system."*
3. **Structure and measurements** — *"Design specs only use spacing tokens available in the system."* + *"Clearly annotate spacing and alignment for all design elements."* + *"Design specs include all possible configurations such as sizes and content configurations."*
4. **Interaction states** — *"Designs include specs for states such as hover, focus, selected, disabled, read-only, error, warning, etc."*
5. **Behaviors** — *"Designs include specs for behaviors such as responsiveness, content overflow or reflow, expansion, scrolling, etc."*
6. **Accessibility** — *"All text colors pass 4.5:1 color contrast with the exceptions of disabled states."* / *"All interactive non-text elements meet 3:1 contrast."* / *"Flow of focus is clearly documented."*

**(c) Code requirements (9 rows):** API guiding principles (prioritize end user / interoperability / stability / composition / developer experience), **Built to spec** (*"The design spec should match the implementation perfectly down to the pixel."*), **Tokens** (*"Component styles do not contain magic numbers or colors that are not tokenized."*), Globalization, Responsiveness (*"works on all device sizes from very large to ~320px wide"*), Storybook, Documentation, Fully Typed/JSDoc, Codemods.

**(d) Testing requirements (4 rows, all mechanically checkable):**
- **Unit testing** — jest + testing-library, *"coverage should meet and exceed 80% of functions, lines, statements"*
- **Visual regression tests (VRT)** — *"at least one test on the default story … using Percy"*
- **Accessibility verification tests (AVT)** — default state **and** *"all additional 'complex' states (open, closed, highlighted, expanded, focused, hovered, clicked, etc)"* checked by **IBM Equal Access Accessibility Checker**
- **Screen reader/voiceover** — manually tested in **JAWS, VoiceOver, and NVDA**

**(e) Four separate documentation artifacts are MANDATORY.** Verbatim: *"All components and patterns require usage, style, code, and accessibility guidance published on a Carbon ecosystem website."* Carbon publishes literal **templates** for each, with variants for single- vs multi-variant components. Their section skeletons (extracted from `documentation.mdx`):

- **Usage template:** `Live demo` → `Overview` (`When to use`, `When not to use`, `Variants`, `Feature flags`) → `Formatting` (`Anatomy`, `Styling`, `Sizing`, `Alignment`, `Placement`) → `Content` (`Main elements` → `Label`, `Body copy`; `Overflow content`; `Further guidance`) → `Universal behaviors` (`States`, `Interactions` → `Mouse`/`Keyboard`, `Validation`, `Responsive behavior`, `Default selection`, `Clickable areas`, `Loading`) → `Modifiers` → `AI presence` → `Related` → `References` → `Feedback`
- **Style template:** `Color` (+ `Interactive state color`) → `Typography` → `Structure` → `Size` → `AI presence` → `Feedback`. Each section is prescribed as **Table + Image**.
- **Code template:** `Overview` (`Skeleton`) → use-case sections → `Component API` (`[Component] props`) → `References` → `Feedback`
- **Accessibility template:** `What Carbon provides` (`Keyboard interaction`) → `Design recommendations` (`Headings`) → `Development considerations`
- **Pattern templates** (separate from components): `Overview` (`Anatomy`, `When to use`, `When not to use`) → `Designing with [pattern]` → `Behaviors` → `Best practices` → `Visual guidance` → `Other use cases` → `Accessibility` → `Related` → `References`

**(f) Design-kit governance:** Figma components must follow IBM Figma Best Practices + naming convention, be **"Built to spec … perfectly down to the pixel"**, and be **published to a library**; peer-reviewed by IBM's Figma guild, described as *"crucial partners in our governance model."*

**Verified Carbon numbers:**
- Spacing scale, from the published `@carbon/layout@11.55.0` tarball (`scss/`): `$spacing-01: 0.125rem` (2px), `-02: 0.25rem`, `-03: 0.5rem`, `-04: 0.75rem`, `-05: 1rem`, `-06: 1.5rem`, `-07: 2rem`, `-08: 2.5rem`, `-09: 3rem`, `-10: 4rem`, `-11: 5rem`, `-12: 6rem`, `-13: 10rem`. → **13 steps, 2px base, NOT a pure 8pt grid.**
- Type scale, from `@carbon/type` `scss/_scale.scss`: **not a modular ratio.** Verbatim formula comment: `// Yn = Yn-1 + {INT[(n-2)/4] + 1} * 2`, `$step == 1 → 12px`, 23 steps, doc comment *"supports sizes from 12px to 92px"*. Computed: **12, 14, 16, 18, 20, 24, 28, 32, 36, 42, 48, 54, 60, 68, 76, 84, 92**.

### 3.2 Shopify Polaris

- Docs site now redirects `polaris.shopify.com/design/*` → **`polaris-react.shopify.com/*`** (301 verified) — Polaris has moved to **Polaris Web Components**; the React docs are on a separate host. Cite carefully.
- Sections: **Foundations, Content, Design, Components, Patterns** (+ Experience). Content section covers **Voice and Tone, Naming, Product Content** — i.e. **content/voice guidelines are a first-class canonical artifact**, not an appendix.
- **Token naming convention:** CSS custom properties prefixed **`--p-`**; pattern is `--p-<group>-<property>-<role>[-<state>]`, e.g. `--p-color-bg-surface`, `--p-color-bg-surface-secondary`, `--p-color-bg-surface-tertiary`, with state suffixes `-hover`, `-active`, `-selected`, `-disabled`.
- **Token groups (extracted directly from `@shopify/polaris-tokens@9.4.2` `dist/css/styles.css`):** `--p-border-`, `--p-breakpoints-`, `--p-color-`, `--p-font-`, `--p-height-`, `--p-motion-`, `--p-shadow-`, `--p-space-`, `--p-text-`, `--p-width-`, `--p-z-`.

**Verified Polaris numbers (from the published tarball, not from docs prose):**

- **Space (18 steps, base 4px = 0.25rem, numbering ×100 of a 0.25rem unit):** `space-0: 0`, `025: 0.0625rem`, `050: 0.125rem`, `100: 0.25rem`, `150: 0.375rem`, `200: 0.5rem`, `300: 0.75rem`, `400: 1rem`, `500: 1.25rem`, `600: 1.5rem`, `800: 2rem`, `1000: 2.5rem`, `1200: 3rem`, `1600: 4rem`, `2000: 5rem`, `2400: 6rem`, `2800: 7rem`, `3200: 8rem`. ⭐ Note the **deliberate gaps** (no 700, 900, 1100…) — the scale is *pruned*, not uniform.
- **Semantic space aliases exist:** `--p-space-card-padding: var(--p-space-400)`, `--p-space-card-gap`, `--p-space-button-group-gap: var(--p-space-200)`, `--p-space-table-cell-padding: var(--p-space-150)`. → **component-level tokens alias primitive tokens.**
- **Font size (13 steps):** `275: 0.6875rem` (11px), `300: 0.75rem`, `325: 0.8125rem` (13px), `350: 0.875rem`, `400: 1rem`, `450: 1.125rem`, `500: 1.25rem`, `550: 1.375rem`, `600: 1.5rem`, `750: 1.875rem`, `800: 2rem`, `900: 2.25rem`, `1000: 2.5rem`.
- **Line height (8 steps, all multiples of 4px):** `300: 0.75rem`, `400: 1rem`, `500: 1.25rem`, `600: 1.5rem`, `700: 1.75rem`, `800: 2rem`, `1000: 2.5rem`, `1200: 3rem`.
- **Border radius (10 steps):** `0, 050: 0.125rem, 100: 0.25rem, 150, 200, 300, 400, 500, 750: 1.875rem, full: 624.9375rem`.
- **Shadow (7 elevation steps + insets + per-component):** `shadow-0: none` … `shadow-600: 0 1.25rem 1.25rem -0.5rem rgba(26,26,26,0.28)`. Note the **single hue `rgba(26,26,26,α)` used across all elevations with α rising 0.07 → 0.28** — a systematic shadow ramp, not ad-hoc blurs.
- **⭐ Composite `text-*` tokens** — Polaris ships full typographic *styles*, not just sizes: `--p-text-heading-3xl-{font-family,font-size,font-weight,font-letter-spacing,font-line-height}` down through `heading-2xl / xl / lg / md / sm / xs` and `body-lg / md / sm / xs`. Every one is defined **only as aliases** to primitives (e.g. `--p-text-heading-lg-font-size: var(--p-font-size-500)`). This is exactly DTCG's composite `typography` type, expressed in CSS.

### 3.3 GOV.UK Design System — *the best governance model*

Source: https://design-system.service.gov.uk/community/contribution-criteria

**Criteria for proposing (verbatim):**
- **Useful** — *"There is evidence that this component or pattern would be useful for many teams or services."*
- **Unique** — *"It does not replicate something already in the Design System."*

**Criteria for publishing (verbatim):**
- **Usable** — *"It has been tested in user research and shown to work with a representative sample of users, including those with disabilities."*
- **Consistent** — *"It reuses existing styles and components in the Design System where relevant."*
- **Versatile** — *"The implementation is versatile enough that the component or pattern can be used in a range of different services that may need it."*

**Governance body:** a **Design System working group** — *"a multidisciplinary panel of representatives from across government"* that reviews proposals against the criteria, makes improvement recommendations, and reviews *new components/patterns, recognisable changes, changes affecting when to use a component, and content affecting how services meet the Service Standard.* Process: **Propose → Review → Develop → Publish.**
Sources: https://design-system.service.gov.uk/community/design-system-working-group, https://team-playbook.design-system.service.gov.uk/how-we-work/contribution-model

**⭐ Verified GOV.UK numbers (from `alphagov/govuk-frontend` source, not docs):**
- `settings/_spacing.scss` → `$govuk-spacing-points: (0: 0, 1: 5px, 2: 10px, 3: 15px, 4: 20px, 5: 25px, 6: 30px, 7: 40px, 8: 50px, 9: 60px)`. **This is a 5px grid, not 4/8pt.** Plus a **`$govuk-spacing-responsive-scale`** where each point can differ per breakpoint (e.g. point 4 = `15px` mobile / `20px` tablet; point 5 = `15px` / `25px`).
- `settings/_typography-responsive.scss` → `$govuk-root-font-size: 16px`; `$govuk-typography-scale` points **80, 48, 36, 27, 24, 19, 16, 14**, each with **mobile / tablet / print** font-size *and* line-height. e.g. point 80 = 53px/55px mobile, 80px/80px tablet, 53pt/1.1 print. Point 36 = 27px/30px mobile, 36px/40px tablet.
- **This is the strongest single counter-example to "always use an 8pt grid and a modular ratio."** A world-class, evidence-driven system uses a 5px grid and a hand-picked, breakpoint-responsive type scale.

### 3.4 Material Design 3 (+ M3 Expressive)

- **Three-tier token architecture:** `md.ref.*` (reference — raw palette) → `md.sys.*` (system — semantic roles) → `md.comp.*` (component). e.g. `md.ref.palette.primary0`, `md.sys.color.primary`, `md.comp.button.container.color`. Source: https://m3.material.io/foundations/design-tokens/overview
- **M3 Expressive** announced **May 2025**, shipped on eligible Pixel devices with **Android 16 QPR1, September 2025**; adds components, variants, and updates to **shape, motion and typography**; typography gains **"emphasized"** styles alongside baseline — *"15 baseline and 15 emphasized"* type styles across `display / headline / title / body / label` × `large / medium / small`.

**⭐ Verified M3 type scale numbers** — read from the generated token source `material-components/material-web` → `tokens/versions/v0_192/_md-sys-typescale.scss` (m3.material.io is JS-rendered and unfetchable; the generated tokens are authoritative):

| Role | size | line-height | weight |
|---|---|---|---|
| display-large | 3.5625rem (57px) | 4rem (64) | regular |
| display-medium | 2.8125rem (45px) | 3.25rem (52) | regular |
| display-small | 2.25rem (36px) | 2.75rem (44) | regular |
| headline-large | 2rem (32px) | 2.5rem (40) | regular |
| headline-medium | 1.75rem (28px) | 2.25rem (36) | regular |
| headline-small | 1.5rem (24px) | 2rem (32) | regular |
| title-large | 1.375rem (22px) | 1.75rem (28) | regular |
| title-medium | 1rem (16px) | 1.5rem (24) | medium |
| title-small | 0.875rem (14px) | 1.25rem (20) | medium |
| body-large | 1rem (16px) | 1.5rem (24) | regular |
| body-medium | 0.875rem (14px) | 1.25rem (20) | regular |
| body-small | 0.75rem (12px) | 1rem (16) | regular |
| label-large | 0.875rem (14px) | 1.25rem (20) | medium (+ `weight-prominent`: bold) |
| label-medium | 0.75rem (12px) | 1rem (16) | medium (+ prominent: bold) |
| label-small | 0.6875rem (11px) | 1rem (16) | medium |

**Every single line-height is a multiple of 4px.** Sizes are 57/45/36/32/28/24/22/16/14/16/14/12/14/12/11 — **not** a clean ratio; they are hand-tuned and then snapped to a 4px leading grid.

### 3.5 Atlassian Design System

Source: https://atlassian.design/foundations/tokens/design-tokens

**Naming anatomy = `foundation.property.modifier`.** Verbatim definitions:
- **Foundation** — *"The type of visual design attribute or foundational style, such as color, elevation, or space"*
- **Property** — *"The UI element the token is being applied to, such as a border, background, shadow, or other property"*
- **Modifier** — *"Additional details about the token's purpose, such as its color role, emphasis, or interaction state"*

Examples: `color.icon.success`; `color.text` (default body text — **modifier omitted when it is the default**). Foundations listed: Color, Elevation, Opacity, Space, Typography, Border, Radius. Framing quote: *"Design tokens are the single source of truth to name and store design decisions."*
Color rule (verbatim, from the color foundation page): all color tokens start with `color`, then property (background/border/icon/text), then **one or more modifiers representing color role, emphasis level, and interaction state**.

### 3.6 Adobe Spectrum (brief)

- Two orthogonal sizing concepts: **scale** (whole-page: **desktop scale** vs **mobile scale**, for cursor vs touch) and **t-shirt sizes** (per-component variant). *"A component with t-shirt sizing is still affected by scale."*
- Token names encode the t-shirt size: `spectrum-button-primary-textonly-height` → `spectrum-button-m-primary-textonly-height`.
- `@adobe/spectrum-tokens@14.15.0` on npm ships tokens for *all* t-shirt sizes *"whether or not they change between sizes"* — deliberate redundancy for API stability.
Sources: https://spectrum.adobe.com/page/platform-scale/, https://github.com/adobe/spectrum-tokens/blob/main/packages/tokens/README.md

---

## 4. FACT + INFERENCE — Atomic Design and competing models

**FACT.** Brad Frost, *Atomic Design* (2013, book at https://atomicdesign.bradfrost.com/): **atoms → molecules → organisms → templates → pages**. Atoms = HTML tags (label, input, button); molecules = simple functional groups (label+input+button = search form); organisms = distinct interface sections; templates = layout/structure with placeholder content; pages = templates with real content.

**FACT — the criticism comes from the author himself.** Brad Frost has stated the **specific labels "have never been the point"** and that he doesn't really use them in his own work, though they remain a useful mental model; and explicitly: *"atomic design is not a linear process, but rather a mental model."* His guidance: *"whatever taxonomy you choose to work with should help you and your organization communicate more effectively."*

**FACT — practitioner criticism themes** (e.g. https://www.qt.io/software-insights/atomic-design-systems-why-the-labels-dont-matter, "From Template to Atoms"): the atom/molecule/organism boundary is **arbitrary and endlessly re-litigated**; `templates`/`pages` overlap with routing frameworks; the chemistry metaphor gives **no guidance on tokens, states, or accessibility**.

**FACT — what the real systems actually use instead.** None of Carbon, Polaris, GOV.UK, Material, Atlassian organises its public docs by atoms/molecules/organisms. They all use a flatter, **purpose-based** taxonomy:

> **Foundations** (color, type, space, motion, elevation, icons, grid) → **Tokens** → **Components** → **Patterns** (multi-component flows) → **Content/voice** → **Accessibility**

**INFERENCE.** For `parley-design`, Atomic Design is worth **one paragraph as vocabulary**, and worth **zero** as the artifact taxonomy. The Foundations→Tokens→Components→Patterns→Content→A11y spine is what five independent billion-dollar systems converged on, and it maps 1:1 onto files. Use that.

---

## 5. FACT — Typography / spacing / color, with numbers

### 5.1 Modular type scales
Standard named ratios: **1.125 major second, 1.200 minor third, 1.250 major third, 1.333 perfect fourth, 1.414 augmented fourth, 1.500 perfect fifth, 1.618 golden ratio.** Origin in the design-web canon: Tim Brown, *"More Meaningful Typography"*, A List Apart (https://alistapart.com/article/more-meaningful-typography/). Common practitioner guidance: 1.125–1.25 subtle/dense, 1.25–1.333 standard UI, 1.5–1.618 dramatic/editorial.

**⚠️ Counter-fact (important, verified above):** **none** of Carbon, Material 3, Polaris, or GOV.UK actually uses a pure geometric ratio. Carbon uses an explicit arithmetic formula (`Yn = Yn-1 + {INT[(n-2)/4]+1}*2`); Material 3 and GOV.UK use hand-tuned tables; Polaris uses a pruned ×0.0625rem numeric scale. **INFERENCE:** the professional practice is *"generate with a ratio, then round to the leading grid and prune the steps you don't need"*, not *"emit `1.25^n` and ship it."* Un-rounded ratio output (`19.2px`, `23.04px`, `27.65px`) is itself an AI-slop tell.

### 5.2 Spacing grids
- **4pt/8pt grid** is the dominant convention (Material's 4dp/8dp grid; Polaris base 4px; Carbon base 2px with everything ≥`spacing-03` on 4px).
- **Counter-example:** GOV.UK uses **5px** points (5/10/15/20/25/30/40/50/60) and makes them **responsive per breakpoint**.
- **INFERENCE / rule worth encoding:** the checkable invariant is **"every spacing value is a member of the declared scale"**, *not* "every value is divisible by 8". Assert membership, not arithmetic.
- **Leading grid:** M3's 15 line-heights are **all multiples of 4px**; Polaris's 8 line-heights are all multiples of 4px. This is a real, checkable rule: `line-height % 4px == 0`.

### 5.3 Color — perceptual uniformity
- **CSS Color Module Level 4** (https://www.w3.org/TR/css-color-4/) defines `lab()`, `lch()`, `oklab()`, `oklch()`; **Candidate Recommendation since 5 July 2022**; **shipping in all evergreen browsers** and included in the stable CSS definition as of the 2026 snapshot.
- **Oklab** was created by **Björn Ottosson (2020)** specifically to fix LCH's hue shift; **OKLCH = Oklab in polar form (L, C, H)** with better perceptual uniformity, notably in blues/purples. (https://en.wikipedia.org/wiki/Oklab_color_space, https://evilmartians.com/chronicles/oklch-in-css-why-quit-rgb-hsl)
- **DTCG lists `oklch` and `oklab` among its 14 `colorSpace` values** (§1.3) — so "author in OKLCH" is standards-aligned, not a fashion.
- **INFERENCE:** the practical rule is *author ramps in OKLCH holding L constant across hues* → equal-lightness palettes → contrast behaves predictably. HSL does not have this property (HSL `L` is not perceptual lightness), which is why HSL ramps look muddy at the same "lightness".

### 5.4 Contrast and accessibility numbers
Source: https://www.w3.org/TR/WCAG22/ — **W3C Recommendation, dated 12 December 2024** (original 5 Oct 2023).
- **SC 1.4.3 Contrast (Minimum), AA:** *"The visual presentation of text and images of text has a contrast ratio of at least 4.5:1"*. **Large text = "at least 18 point or 14 point bold"** (≈24px / 18.66px) and needs only **3:1**.
- **SC 1.4.11 Non-text Contrast, AA:** *"at least 3:1 against adjacent color(s)"* for **User Interface Components** and **Graphical Objects** required for understanding. (This is what catches invisible borders, low-contrast focus rings, and ghost icons — classic slop.)
- **SC 2.5.8 Target Size (Minimum), AA (new in 2.2):** targets **at least 24×24 CSS px**, or sufficient spacing.
- **SC 1.4.12 Text Spacing, AA:** content must survive **line-height ≥ 1.5× font-size; paragraph spacing ≥ 2× font-size; letter-spacing ≥ 0.12em; word-spacing ≥ 0.16em.**
- **APCA / WCAG 3 status (FACT, as of 2026):** **WCAG 3.0 is still a Working Draft**; **the visual-contrast section was moved OUT of the Working Draft in July 2023** and the current draft says the contrast algorithm is **"yet to be determined."** APCA is the leading candidate (`Lc` values, e.g. `Lc 60`, factoring size/weight/polarity). WCAG 3 **is not expected to reach Recommendation before 2028**. Sources: https://adrianroselli.com/2026/04/wcag3-contrast-as-of-april-2026.html, https://yatil.net/blog/wcag-3-is-not-ready-yet
- **INFERENCE / rule to encode:** **WCAG 2.2 ratios are the gate (blocking); APCA `Lc` is advisory (warning).** Do not fail a build on APCA.

### 5.5 Optical sizing
- `opsz` is an **OpenType registered variation axis**; the OpenType spec states *"The scale for the Optical size axis is text size in points"* and that applications may auto-select an optical-size variant based on text size.
- CSS: **`font-optical-sizing: auto | none`** (only two values); **`auto` is the default when the font has an `opsz` axis**. Finer control via `font-variation-settings: 'opsz' <num>`, where you normally match `opsz` to `font-size`.
- Effect: lower `opsz` → wider glyphs, wider spacing, thicker strokes, taller x-height (better at small sizes); higher `opsz` → tighter, finer (better for display).
Sources: https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Properties/font-optical-sizing, https://fonts.google.com/knowledge/glossary/optical_size_axis

---

## 6. RECOMMENDED CANONICAL ARTIFACT SET for `parley-design` (INFERENCE, grounded in §1–§5)

Design goal: **every artifact is either machine-checkable or explicitly marked as human-judgement.** Files live under the existing Parley Deck idea directory so they inherit the deck's ownership/consensus rules.

```
parley-deck/ideas/<slug>/
├── DESIGN-BRIEF.md                 # phase 0 input
├── directions/
│   ├── DIRECTION-<agent>.md        # one per agent — the DIVERGE artifact
│   └── DIRECTION-<agent>.tokens.json
├── DESIGN-FINAL.md                 # the winning direction + grafts + rejected-with-reasons
├── design/
│   ├── tokens/
│   │   ├── primitives.tokens.json  # DTCG, $type-annotated, OKLCH
│   │   ├── semantic.tokens.json    # DTCG, aliases only ({...} / $ref)
│   │   ├── component.tokens.json   # DTCG, aliases only
│   │   └── resolver.json           # DTCG Resolver 2025.10 — modes
│   ├── FOUNDATIONS.md
│   ├── TYPE-SCALE.md
│   ├── COLOR-SYSTEM.md
│   ├── SPACING-LAYOUT.md
│   ├── MOTION.md
│   ├── VOICE-AND-CONTENT.md
│   ├── components/
│   │   └── <component>.spec.md     # one file per component
│   ├── PATTERNS.md
│   ├── ACCESSIBILITY.md
│   └── DESIGN-DECISIONS.md         # ADR log: decision, alternatives, why, date
└── design-check.config.json        # the enforcement contract → parley-design-check
```

### 6.1 Artifact table

| # | File | Contains | Machine-checkable? | How |
|---|---|---|---|---|
| 1 | `DESIGN-BRIEF.md` | product, audience, tone adjectives (3–5), **anti-goals** (what this must NOT look like), constraints, reference & anti-reference URLs, success criteria | **Structure only** | required-headings lint; ≥3 anti-goals; ≥1 anti-reference |
| 2 | `directions/DIRECTION-<agent>.md` | ONE named visual direction: name + one-sentence thesis, mood, typographic idea, color idea, spatial idea, **the one risky move**, what it deliberately sacrifices | **Structure + divergence** | required headings; **cross-direction distance check** — reject if two directions share type family AND hue family AND radius bucket |
| 3 | `directions/DIRECTION-<agent>.tokens.json` | a *minimal* DTCG token set proving the direction is real (≥1 type ramp, ≥1 color ramp, radius, space) | **YES, fully** | DTCG JSON Schema `2025.10/format.json` |
| 4 | `DESIGN-FINAL.md` | **which direction won whole**, the 2–3 named grafts with their source direction, **explicit list of losing directions and why each lost**, signatures | **Structure** | must name exactly 1 winner; 2–3 grafts each with `from:`; every non-winning direction must appear in the rejected list |
| 5 | `design/tokens/primitives.tokens.json` | raw values only. **All colors `colorSpace: "oklch"`.** No semantic names. | **YES, fully** | JSON Schema + `core/colorspace` + `core/required-type` + `core/duplicate-values` |
| 6 | `design/tokens/semantic.tokens.json` | roles only (`color.bg.surface`, `color.text.default`, …). **Every `$value` MUST be an alias.** | **YES, fully** | JSON Schema + custom rule: no literal `$value` in this file + `core/descriptions` |
| 7 | `design/tokens/component.tokens.json` | per-component tokens, aliases to semantic only | **YES, fully** | same; plus "no component token references a primitive directly" |
| 8 | `design/tokens/resolver.json` | DTCG Resolver: `light`/`dark` (+ optional `hc`, `density`, `reduced-motion`) | **YES, fully** | JSON Schema `2025.10/resolver.json`; `core/required-modes` |
| 9 | `design/FOUNDATIONS.md` | the index: what each foundation is and the one rule that governs it | Structure | required headings |
| 10 | `design/TYPE-SCALE.md` | the ratio *used to generate*, the **rounded** final table (size, line-height, weight, letter-spacing, `opsz` if variable), and the role→step mapping | **YES** | table must match `primitives.tokens.json`; **`line-height % 4px == 0`**; no un-rounded values |
| 11 | `design/COLOR-SYSTEM.md` | ramp construction in OKLCH (L steps held constant across hues), semantic role map, **full contrast matrix** of every text-on-bg pair with computed ratios | **YES** | recompute every ratio; fail <4.5:1 body / <3:1 large / <3:1 non-text |
| 12 | `design/SPACING-LAYOUT.md` | the declared spacing scale (explicit member list), grid, breakpoints, container widths, radius scale, elevation/shadow ramp | **YES** | membership assertion against the token file |
| 13 | `design/MOTION.md` | duration + easing tokens per role, and the **`prefers-reduced-motion` rule** | **YES** | every animation token used must have a reduced-motion counterpart |
| 14 | `design/VOICE-AND-CONTENT.md` | voice & tone, capitalisation rule, button/label/error-message patterns, terminology list, do/don't pairs | Partial | required headings; ≥5 do/don't pairs; terminology table parseable |
| 15 | `design/components/<c>.spec.md` | **see §6.2** | **Structure fully; tokens fully** | required-section lint + token-reference resolution |
| 16 | `design/PATTERNS.md` | multi-component flows (empty state, error, loading, destructive confirm, form) | Structure | required headings |
| 17 | `design/ACCESSIBILITY.md` | keyboard model, focus order, focus-visible spec, landmark/heading structure, target sizes, SR expectations | **Partial** | assert 24×24 targets, focus ring ≥3:1, text-spacing survivability |
| 18 | `design/DESIGN-DECISIONS.md` | ADR log: `ID / date / decision / alternatives considered / why / consequences` | Structure | each entry needs ≥2 alternatives + a "why not" |
| 19 | `design-check.config.json` | the contract `parley-design-check` executes: token file globs, allowed colorspaces, contrast thresholds, banned-value regexes, required component sections, status enum | **IS the checker input** | JSON Schema of its own |

### 6.2 `<component>.spec.md` — required sections (merged Carbon + Polaris + GOV.UK)

```
# <Component>            status: draft | preview | stable     # Carbon PDLC enum
## Purpose               (1 sentence)
## When to use           (bullets)                            # GOV.UK + Carbon
## When NOT to use       (bullets, MUST name the alternative) # GOV.UK "unique"
## Anatomy               (numbered parts list)                # Carbon
## Variants              (table: name | when | token deltas)
## Sizes                 (table)                              # Spectrum t-shirt
## States                (table: default, hover, focus-visible, active,
                          selected, disabled, read-only, error, warning, loading)
                          — every cell references TOKENS, never literals   # Carbon
## Behaviors             (responsive, overflow/reflow, truncation,
                          expansion, scroll, empty, long-content)          # Carbon
## Interactions          (### Mouse  ### Keyboard  ### Touch)              # Carbon
## Content guidelines    (label rules, char budgets, capitalisation)       # Polaris
## Accessibility         (role, name, ARIA, focus order, SR announcement,
                          contrast results, target size)                  # Carbon
## Do / Don't            (≥3 pairs, each one sentence)
## Tokens used           (explicit list — MUST resolve)
## Related               (links)
```

**Machine-checkable in that file:** every `##` heading present; `status` ∈ enum; **States table contains all 10 required rows**; Interactions has Mouse+Keyboard subsections; **every token named resolves in the token files**; **no literal hex / no bare `px` outside the tokens list**; ≥3 do/don't pairs; When-NOT-to-use names an alternative.

### 6.3 Which artifacts are hard gates vs advisory (INFERENCE)

- **BLOCKING (mechanical, zero judgement):** #3, #5, #6, #7, #8 (JSON Schema); contrast matrix in #11; token-resolution + required sections + banned literals in #15.
- **BLOCKING (structural):** required headings in #1, #2, #4, #10, #12, #13.
- **ADVISORY (warn):** APCA `Lc`; `core/duplicate-values`; prose quality in #14, #18.
- Rationale: this mirrors Carbon's split between the hard "Definition of done" checklist and its softer prose templates, and keeps `parley-design-check` from becoming an LLM taste oracle.

### 6.4 Protocol shape (matches the CopilotKit/AG-UI style the owner asked for) (INFERENCE)

Spec `parley-design` as a **versioned protocol**, mirroring how DTCG itself is written:

- **`spec_version: "1.0.0"`** declared in `design-check.config.json`; artifacts carry it too.
- **Typed artifacts** — each of the 19 files above is a *type* with a required shape (exactly how DTCG types tokens).
- **Explicit phases**, named and numbered, so the driver can report them:
  `D0 Brief → D1 Diverge → D2 Cross-critique → D3 Selection (one wins whole) → D4 Graft → D5 Tokenize → D6 Component specs → D7 Apply to UI → D8 Check → D9 Ratify`
- **Normative language** — use RFC 2119 **MUST / SHOULD / MAY** exactly as DTCG does (*"A set MUST contain a `sources` array"*).
- **Conformance section** — "an implementation conforms if it produces artifacts 4, 5, 6, 8, 15 that pass `parley-design-check --strict`."
- **A JSON Schema of your own** for `design-check.config.json`, published in-repo, so the enforcement contract is itself validatable.

---

## Transferable to parley-design / parley-design-check

Ranked, most valuable first. Each is concrete and sourced.

1. **Adopt DTCG `2025.10` verbatim as the token format — do not invent one.** Files are `*.tokens.json`; `$value`/`$type`/`$description`/`$extensions`/`$deprecated`; aliases via `{group.token}` and `$ref` JSON Pointer; the 13 `$type`s. → *parley-design writes them, parley-design-check validates them.* (§1.2)
2. **Validate token files against the official JSON Schema `https://www.designtokens.org/schemas/2025.10/format.json` (and `resolver.json`).** Verified live, HTTP 200, draft-07. This is `parley-design-check`'s cheapest, hardest, most defensible gate — pure schema, no heuristics. (§1.5)
3. **Author all color in `colorSpace: "oklch"`, with a `hex` fallback in 6-digit form.** Standards-aligned (DTCG lists `oklch` among 14 spaces; CSS Color 4 ships everywhere). Enforce with a `core/colorspace`-equivalent rule. (§1.3, §5.3)
4. **Steal Carbon's "Definition of done" + 4-value status enum** (`Draft / Preview candidate / Preview / Stable`) and put `status:` in the front-matter of every component spec. Gives `parley-design-check` a first-class thing to gate on. (§3.1a)
5. **Steal Carbon's six design-spec requirement rows as the checker's rule set:** color tokens only, type tokens only, spacing tokens only + annotated, **all interaction states specified**, **behaviors specified**, **4.5:1 text / 3:1 non-text / documented focus flow**. Verbatim-quotable and each one is mechanically testable. (§3.1b)
6. **Steal Carbon's "no magic numbers" rule as literal lint:** *"Component styles do not contain magic numbers or colors that are not tokenized."* Implement with `stylelint-declaration-strict-value@1.11.1` (verified on npm) — ban raw hex/px, require `var(--token)`. This is the single highest-yield anti-slop check that exists today. (§3.1c, §2.3)
7. **Steal Carbon's four-artifact documentation rule:** *"All components and patterns require usage, style, code, and accessibility guidance."* Merge into one `<component>.spec.md` with the required-section list in §6.2, and lint the headings. (§3.1e)
8. **Steal GOV.UK's five contribution criteria verbatim as the acceptance gate** — proposing: **Useful, Unique**; publishing: **Usable, Consistent, Versatile**. They map perfectly onto Parley Deck's diverge/critique/consensus flow: *Unique* is the divergence check, *Consistent* is the token-reuse check, *Versatile* is the "does it survive other contexts" check. (§3.3)
9. **Adopt the three-tier token layering** (primitives → semantic → component) from Material's `md.ref` / `md.sys` / `md.comp`, and enforce the direction of references: **component→semantic→primitive only, never skipping or reversing.** This is a graph assertion, trivially checkable. (§3.4, §6.1 rows 5–7)
10. **Adopt Atlassian's `foundation.property.modifier` naming anatomy** for semantic tokens (`color.icon.success`, `color.text`), including the rule that **the modifier is omitted for the default**. Enforce with a `core/consistent-naming`-equivalent. (§3.5)
11. **Use DTCG Resolver for modes instead of duplicating token files per theme.** `version`, `sets`, `modifiers.contexts`, `resolutionOrder`. Require at minimum `light` + `dark`; recommend `hc` and `reduced-motion`. Enforce "every semantic color token has a value in every declared context." (§1.4)
12. **Run Terrazzo's 27 built-in lint rules** — especially `core/duplicate-values` (two tokens, same value = the system is lying), `core/descriptions` (forces stated intent), `core/required-modes`, `a11y/min-contrast`, `a11y/min-font-size`. `@terrazzo/cli@2.5.0`, config `lint.rules["rule"] = "error"|"warn"|"off"`. **Note: linting is in the CLI, not in `@terrazzo/plugin-lint-*` packages (those don't exist).** (§2.2)
13. **Encode the WCAG 2.2 numbers as blocking constants:** 4.5:1 body, 3:1 large (≥18pt / 14pt bold), 3:1 non-text/UI/graphics, 24×24 CSS px targets, text-spacing survivability at 1.5 / 2× / 0.12em / 0.16em. Cite W3C Recommendation **12 December 2024**. Emit a full **contrast matrix** artifact, not a claim. (§5.4)
14. **Assert scale MEMBERSHIP, not arithmetic.** "Every spacing value is in the declared scale" beats "divisible by 8" — GOV.UK's 5px grid proves 8pt is a convention, not a law. Same for type: assert the value is a step in `TYPE-SCALE.md`. Add one true arithmetic rule that *is* universal in practice: **line-height must be a multiple of 4px** (holds for all 15 M3 styles and all 8 Polaris line-heights). (§5.2, §3.4)
15. **Generate the type scale from a named ratio, then ROUND and PRUNE, and record both.** `TYPE-SCALE.md` states the ratio used *and* the final rounded table. Un-rounded ratio output (19.2px, 23.04px) is a checkable slop signature. Carbon's `Yn = Yn-1 + {INT[(n-2)/4]+1}*2` and Polaris's pruned scale (no 700/900/1100) are the proof that pros don't ship raw `1.25^n`. (§5.1, §3.1, §3.2)
16. **Ship composite typography tokens, not just font sizes** — Polaris's `--p-text-heading-lg-{font-family,size,weight,letter-spacing,line-height}` pattern, i.e. DTCG's composite `typography` `$type`. Enforce: a component spec MUST reference a text *style* token, never an isolated font-size. (§3.2, §1.2)
17. **Ship a systematic elevation ramp, not ad-hoc shadows.** Polaris uses one hue `rgba(26,26,26,α)` with α climbing 0.07→0.28 across `shadow-100..600`. Checkable rule: all shadow tokens share one hue; alpha is monotonic; count ≤7. (§3.2)
18. **Require a `## When NOT to use` section that names an alternative** — GOV.UK "unique" + Carbon usage template. Cheap to lint, and it is the section that most reliably exposes an unmotivated component. (§3.3, §3.1e)
19. **Require an ADR log (`DESIGN-DECISIONS.md`) with ≥2 considered alternatives per decision.** This is the artifact that makes the DIVERGE→one-wins-whole method auditable after the fact, and it is what the losers' grafts get recorded in. (§6.1 #18)
20. **Write the skill as a versioned spec with RFC 2119 MUST/SHOULD/MAY, numbered phases (D0–D9), typed artifacts, and a Conformance section** — exactly how DTCG writes itself (*"A set MUST contain a `sources` array"*). Publish a JSON Schema for `design-check.config.json` so the enforcement contract is itself machine-validated. (§6.4)
21. **Two-layer contrast policy:** WCAG 2.2 ratios BLOCK; APCA `Lc` WARNS. WCAG 3 is still a Working Draft, contrast was pulled out in July 2023, algorithm "yet to be determined", no Recommendation expected before 2028. (§5.4)
22. **`font-optical-sizing: auto` (or matched `opsz`) whenever the chosen family is variable with an `opsz` axis** — a small, checkable, high-signal craft marker. (§5.5)
23. **Taxonomy = Foundations → Tokens → Components → Patterns → Content/Voice → Accessibility.** Five independent major systems converged on it; it maps 1:1 to directories. (§4)

---

## Do NOT copy

1. **Do NOT structure the artifacts by Atomic Design (atoms/molecules/organisms/templates/pages).** *Reason:* Brad Frost himself says the labels *"have never been the point"* and that he doesn't use them; the boundaries are endlessly re-litigated; and **none** of Carbon, Polaris, GOV.UK, Material, or Atlassian organises its docs that way. It also gives zero guidance on tokens, states, or a11y — the three things `parley-design-check` needs. Mention it once as vocabulary; do not build the file tree from it. (§4)
2. **Do NOT cite or link `tr.designtokens.org`.** *Reason:* verified 301 → `www.designtokens.org/TR/drafts/`. And do not cite the `/tr/drafts/` copies as normative — they carry *"Do not attempt to implement this version of the specification."* **Cite `https://www.designtokens.org/tr/2025.10/`** for the stable text. (§1.1)
3. **Do NOT hard-code an 8pt grid as a universal law.** *Reason:* GOV.UK Frontend ships a **5px** scale (0/5/10/15/20/25/30/40/50/60) and Carbon's base step is **2px**. Asserting `value % 8 == 0` would fail two world-class systems. Assert membership in the declared scale instead. (§5.2, §3.3)
4. **Do NOT ship a raw geometric modular scale (`1.25^n`) as the type scale.** *Reason:* no mature system does. Carbon uses an arithmetic formula, M3 and GOV.UK use hand-tuned tables, Polaris prunes steps. Raw ratio output produces 19.2px/23.04px/27.65px — itself a slop tell. Generate → round → prune → record. (§5.1)
5. **Do NOT make APCA / WCAG 3 a blocking gate.** *Reason:* WCAG 3.0 remains a Working Draft; visual contrast was **removed** from the draft in July 2023; the current draft says the algorithm is **"yet to be determined"**; Recommendation not expected before 2028. Blocking on an undetermined algorithm would be indefensible. Advisory only. (§5.4)
6. **Do NOT reference `@terrazzo/plugin-lint-a11y` or `@terrazzo/plugin-lint-core`.** *Reason:* verified — they **do not exist on npm** (404). Linting ships inside `@terrazzo/cli@2.5.0`. Writing install instructions for phantom packages is exactly the kind of confident error §13 of the protocol is meant to catch. (§2.2)
7. **Do NOT copy Carbon's vendor-specific tooling requirements** — Percy for VRT, **IBM Equal Access Accessibility Checker** for AVT, IBM Figma guild review, `@carbon/upgrade` codemods, JAWS/NVDA manual passes. *Reason:* `parley-design` must be vendor-neutral and runtime-agnostic (same constraint as `parley-worktrees` / `parley-tracker`). Copy the **requirement** ("every complex state is machine-checked for a11y"), not the **vendor**. Let `design-check.config.json` name whichever checker the repo has.
8. **Do NOT copy Carbon's 80% unit-coverage number into the design skill.** *Reason:* it is a code-quality threshold from a specific org's CI, not a design-system fact, and Parley Deck already has `RunChecks` for code gates. Importing it invites arguments that have nothing to do with design.
9. **Do NOT copy Adobe Spectrum's "emit tokens for every t-shirt size whether or not they change."** *Reason:* deliberate redundancy that only pays off at Adobe's API-stability scale; for everyone else it directly triggers `core/duplicate-values` and inflates the surface `parley-design-check` has to reason about.
10. **Do NOT copy Material 3's `md.ref` / `md.sys` / `md.comp` NAMES.** *Reason:* copy the three-tier *layering* (which is the transferable idea) but not the Google-branded prefixes — vendor-neutrality again. Use `primitive.` / `semantic.` / `component.` or Atlassian's `foundation.property.modifier`.
11. **Do NOT copy Polaris's `--p-` prefix or its per-component shadow tokens** (`--p-shadow-button-primary-critical-hover`, etc.). *Reason:* `--p-` is Shopify-branded, and that shadow set is a bespoke skeuomorphic button treatment — 8 inset shadows for one button is the opposite of the systematic ramp we want to teach.
12. **Do NOT copy Carbon's requirement that the design spec match implementation "perfectly down to the pixel."** *Reason:* unfalsifiable without a Figma-to-DOM diffing pipeline that `parley-design-check` will not have. Replace with the checkable form: *every value in the implementation resolves to a declared token*. Keep the spirit, drop the unmeasurable phrasing.
13. **Do NOT copy GOV.UK's human working-group as a required process step.** *Reason:* Parley Deck's adversarial cross-review **is** the working group. Copy the **criteria** (Useful/Unique/Usable/Consistent/Versatile) as gate questions; do not invent a standing human panel the owner does not have.
14. **Do NOT let `parley-design-check` grade taste.** *Reason:* every mature system's enforcement layer is mechanical (stylelint plugin, schema, axe, VRT). The moment the checker emits "this feels generic", it becomes an unfalsifiable LLM oracle and teams will disable it. Every BLOCKING rule must be reproducible by a script with no model in the loop; put judgement in `parley-design`'s critique phase (D2), where adversarial agents already argue.

---

## Appendix — local files produced

- `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/research/carbon-checklist.mdx` — full Carbon component checklist source (159 lines)
- `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/research/carbon-doc-templates.mdx` — full Carbon usage/style/code/a11y/pattern documentation templates (~2600 lines)
- `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/ptok/package/dist/css/styles.css` — `@shopify/polaris-tokens@9.4.2` extracted (all token values)
- `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/clay/package/scss/` — `@carbon/layout@11.55.0` spacing scale
- `/private/tmp/claude-501/-Volumes-My-Shared-Files-AI-WORKSPACE-parley-deck/5dc331bd-5ddf-45e0-b6c2-d519d8c05128/scratchpad/ctype/package/scss/_scale.scss` — Carbon type-scale formula

## Sources

- https://www.designtokens.org/tr/2025.10/ · https://www.designtokens.org/tr/drafts/format/ · https://www.designtokens.org/tr/drafts/color/ · https://www.designtokens.org/tr/drafts/resolver/
- https://www.designtokens.org/schemas/2025.10/format.json · https://www.designtokens.org/schemas/2025.10/resolver.json
- https://www.w3.org/community/design-tokens/2025/10/28/design-tokens-specification-reaches-first-stable-version/
- https://styledictionary.com/info/architecture/ · https://styledictionary.com/reference/hooks/transform-groups/predefined/
- https://terrazzo.app/docs/linting/ · https://www.npmjs.com/package/@terrazzo/cli
- https://carbondesignsystem.com/contributing/component-checklist/ · https://carbondesignsystem.com/contributing/documentation/ · https://github.com/carbon-design-system/carbon
- https://design-system.service.gov.uk/community/contribution-criteria · https://design-system.service.gov.uk/community/design-system-working-group · https://team-playbook.design-system.service.gov.uk/how-we-work/contribution-model · https://github.com/alphagov/govuk-frontend
- https://polaris-react.shopify.com/ · https://www.npmjs.com/package/@shopify/polaris-tokens
- https://m3.material.io/foundations/design-tokens/overview · https://m3.material.io/styles/typography/type-scale-tokens · https://github.com/material-components/material-web
- https://atlassian.design/foundations/tokens/design-tokens · https://atlassian.design/foundations/color
- https://spectrum.adobe.com/page/platform-scale/ · https://github.com/adobe/spectrum-tokens
- https://atomicdesign.bradfrost.com/chapter-2/ · https://www.qt.io/software-insights/atomic-design-systems-why-the-labels-dont-matter
- https://www.w3.org/TR/WCAG22/ · https://www.w3.org/TR/css-color-4/ · https://en.wikipedia.org/wiki/Oklab_color_space · https://evilmartians.com/chronicles/oklch-in-css-why-quit-rgb-hsl
- https://adrianroselli.com/2026/04/wcag3-contrast-as-of-april-2026.html · https://yatil.net/blog/wcag-3-is-not-ready-yet
- https://alistapart.com/article/more-meaningful-typography/ · https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Properties/font-optical-sizing · https://fonts.google.com/knowledge/glossary/optical_size_axis
