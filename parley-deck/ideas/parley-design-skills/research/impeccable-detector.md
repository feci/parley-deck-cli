# Impeccable — mechanical anti-slop enforcement machinery (detector)

**Studied**: `pbakaus/impeccable` v3.4.0 (`package.json` `"version": "3.4.0"`, Apache-2.0, © 2026 Paul Bakaus).
**Root read**: `<scratch>/research/impeccable/.agents/skills/impeccable/scripts/` — `detector/**`, `detect.mjs`, `detect-csp.mjs`, `doctor.mjs`.
**Purpose of digest**: extract what a vendor-neutral multi-agent design skill (`parley-design`) can and cannot borrow.

Line counts (real, `wc -l`):

| file | lines | role |
|---|---:|---|
| `detector/registry/antipatterns.mjs` | 627 | the rule registry (60 rules) — **doctrine, no logic** |
| `detector/rules/checks.mjs` | 5536 | every threshold/heuristic — **the machinery** |
| `detector/detect-antipatterns-browser.js` | 8250 | prebuilt bundle of the browser engine |
| `detector/browser/injected/index.mjs` | 2023 | browser engine source |
| `detector/engines/static-html/css-cascade.mjs` | 1183 | hand-rolled CSS cascade for the static engine |
| `detector/design-system.mjs` | 983 | DESIGN.md → allowlist + drift rules |
| `detector/engines/regex/detect-text.mjs` | 785 | text/regex engine |
| `detector/cli/main.mjs` | 438 | CLI, config filtering, advisory partition, exit codes |
| `detector/engines/browser/detect-url.mjs` | 340 | Puppeteer driver |
| `detector/engines/static-html/detect-html.mjs` | 249 | static-HTML orchestrator |
| `detector/node/file-system.mjs` | 212 | walker, import graph, dev-server sniffing |
| `detector/engines/visual/screenshot-contrast.mjs` | 189 | pixel-diff contrast |
| `detector/profile/profiler.mjs` | 166 | timing/telemetry only |
| `detector/shared/inline-ignores.mjs` | 148 | eslint-style in-file waivers |
| `detector/shared/color.mjs` | 124 | color math |
| `detector/shared/constants.mjs` | 112 | font lists, WCAG constants |
| `detector/findings.mjs` | 18 | finding shape |
| `detector/shared/page.mjs` | 7 | `isFullPage()` |

---

## (a) DOCTRINE — design rules / taste / knowledge

Doctrine lives almost entirely in **one file**: `detector/registry/antipatterns.mjs`. It is a flat `const ANTIPATTERNS = [...]` array of 60 objects with **no detection logic whatsoever**. Each entry:

```js
{
  id: 'side-tab',
  category: 'slop',              // 'slop' | 'quality'
  scopes: ['type'],              // optional; only 'type' and 'layout' exist
  severity: 'error'|'advisory',  // optional; default 'warning' (findings.mjs)
  advisory: true,                // optional; ONLY on em-dash-overuse
  name: 'Side-tab accent border',
  description: '...why it is wrong AND what to do instead...',
  skillSection: 'Visual Details',   // back-link into the prose skill
  skillGuideline: 'colored accent stripe',
}
```

Two doctrine categories, explicitly separated in the source with comment banners:

- `// ── AI slop: tells that something was AI-generated ──` → `category: 'slop'`
- `// ── Quality: general design and accessibility issues ──` → `category: 'quality'`
- `// ── Common generated-UI tells ───` → 5 more `slop` rules, all `severity: 'advisory'`

Every `description` is written as **taste, then prescription**. Verbatim examples worth stealing wholesale:

- side-tab: *"Thick colored border on one side of a card — the most recognizable tell of AI-generated UIs."*
- overused-font: *"Inter, Roboto, Fraunces, Geist, Plus Jakarta Sans, and Space Grotesk are used on so many sites they no longer feel distinctive. Each new wave of AI-generated UIs converges on the same handful of faces."*
- cream-palette: *"A warm cream or beige page background has become the default 'tasteful' AI surface, reached for by reflex."*
- icon-tile-stack: *"A small rounded-square icon container above a heading is the universal AI feature-card template — every generator outputs this exact shape."*
- kicker-above-heading: *"…is banned outright, repeated or not. Generated kickers never earn their place: the heading carries its own weight."*
- undersized-ui-text: *"Being ON the DESIGN.md size ramp does not exempt a value here: adding 8px to the ramp launders the token but not the legibility problem, and that is exactly the escape hatch this rule closes."*
- italic-serif-display: *"…reads as taste in isolation but has become the universal AI-startup landing page hero… Editorial / magazine register may legitimately want this — judge by context."*

The rest of the doctrine is a **font monoculture list** in `detector/shared/constants.mjs`:

```js
const OVERUSED_FONTS = new Set([
  // Older monoculture (still ubiquitous):
  'inter','roboto','open sans','lato','montserrat','arial','helvetica',
  // Newer monoculture (the Anthropic-skill / Vercel / GitHub default wave):
  'fraunces','instrument sans','instrument serif',
  'geist','geist sans','geist mono','mona sans',
  'plus jakarta sans','space grotesk','recoleta',
]);
```

plus `KNOWN_SERIF_FONTS` (28 faces: fraunces, recoleta, newsreader, playfair display, cormorant, tiempos, lora, spectral, source serif, ibm plex serif, merriweather, libre caslon/baskerville, georgia, times new roman, dm serif display/text, instrument serif, gt sectra, ogg, canela, freight display/text …), `GENERIC_FONTS` (serif, sans-serif, monospace, cursive, fantasy, system-ui, ui-*, -apple-system, blinkmacsystemfont, segoe ui, inherit/initial/unset/revert) and a brand-exemption map:

```js
const BRAND_FONT_DOMAINS = {
  'roboto': GOOGLE_DOMAINS, 'google sans': …, 'product sans': …,
  'geist'|'geist sans'|'geist mono': VERCEL_DOMAINS,  // vercel.com, nextjs.org, v0.app
  'mona sans': GITHUB_DOMAINS,                        // github.com, githubnext.com
};
```
i.e. **"Geist is slop — unless you are Vercel."** `isBrandFontOnOwnDomain()` checks `location.hostname` suffix; browser-only.

A copy-doctrine list lives in `detect-text.mjs` (`BUZZWORDS`, 30 phrases): `streamline your, empower your, supercharge your, unleash your, unleash the power, leverage the power, built for the modern, trusted by leading, trusted by the world, best-in-class, industry-leading, world-class, enterprise-grade, next-generation, cutting-edge, transform your business, revolutionize, game-changer, game changing, mission-critical, best of breed, future-proof, future proof, seamless experience, seamlessly integrate, drive engagement, drive growth, drive results, harness the power`.

---

## 1. THE COMPLETE ANTIPATTERN REGISTRY (60 rules)

Severity column: registry value, else **warning** (default in `findings.mjs`: `severity: ap.severity || 'warning'`).
`ADV` column = `advisory: true` (the *only* flag that removes a finding from the failure count). **Note the trap: 11 rules carry `severity: 'advisory'` but NOT `advisory: true` — they still fail the build.** See §3.

### Category `slop` — "this was AI-generated"

| # | id | sev | ADV | scopes | What it detects — exact threshold / heuristic |
|---|---|---|---|---|---|
| 1 | `side-tab` | warning | | — | Dominant chromatic single-edge border. Element gate: browser rect ≥ 20×20. Per side: `w ≥ 1` AND color non-neutral AND `w ≥ 2 && (maxOtherSide ≤ 1 \|\| w ≥ maxOther*2)`. **Left/Right**: fires if `radius > 0` (any w≥2), else `w ≥ 3`. **Top/Bottom**: `3 ≤ w ≤ 12` unless `opts.tabContext`. `BORDER_SAFE_TAGS` exempt (SAFE_TAGS minus `label`). Exempt if `role=status\|alert\|alertdialog\|log` or `aria-live=polite\|assertive` (`isStatusContextElement`). `<span>` only via `badgeLike` (own bg α>0.1), and then top/bottom only. **Also detected in 4 other shapes**: (a) `scanCssTextForPseudoStripe` — `::before/::after` + `position:absolute\|fixed`, width 3–12px full-height (top/bottom both 0, or `height:100%`, or top&bottom ∈ [0,20]px) or height 3–12px full-width, chromatic fill (RGB spread ≥ 30, α ≥ 0.1); (b) `checkElementPseudoStripeDOM` — same in live DOM, host rect ≥ 40×20, `h ≥ rect.height-44 && h ≥ rect.height*0.5`, edge offset within ±2px; (c) `scanCssTextForInsetStripe` / `scanInsetStripeCss` — `box-shadow: inset` with `\|x\| or \|y\| ∈ [3,12]`, other axis 0, blur 0, spread 0, chroma ≥ 30, declared width > 40px; (d) regex matchers on `border-l/r-N` Tailwind (`≥2` if rounded else `≥4`), `border-left/right: Npx solid` (`≥2` if border-radius on line else `≥3`). |
| 2 | `border-accent-on-rounded` | warning | | — | Same dominant-edge gate, **Top/Bottom** side, `radius > 0 && w ≥ 2`. Tailwind: `border-t/b-N` with `rounded*` on same line and `N ≥ 1`. |
| 3 | `overused-font` | warning | | type | Primary (first non-generic) family ∈ `OVERUSED_FONTS`. **Browser**: only if ≥ 20 text elements on page AND that font covers ≥ **15%** (`PRIMARY_THRESHOLD = 0.15`) of text-bearing elements, and not `isBrandFontOnOwnDomain`. **Static/regex**: any `font-family:` declaration or Google Fonts `?family=` URL hit — no share threshold. |
| 4 | `single-font` | warning | | type | Exactly one distinct non-generic primary family. Browser: `fontUsage.size === 1` with ≥20 text elements. Static: `fonts.size===1 && document.querySelectorAll('*').length >= 20`. Regex: `fonts.size===1 && content.split('\n').length >= 20`. |
| 5 | `flat-type-hierarchy` | warning | | type | Collect distinct rounded font sizes (0.1px granularity, 8px ≤ s < 200px). If `sizes.size >= 3` and `max/min < 2.0` → fire. Regex path also seeds sizes from `clamp()` endpoints and a hardcoded Tailwind map `text-xs:12 … text-9xl:128`. |
| 6 | `gradient-text` | warning | | — | `background-clip: text` (or `-webkit-`) with `gradient` in the same element's `background-image` / within ±200 chars of raw HTML. Tailwind: `bg-clip-text` + `bg-gradient-to-`. |
| 7 | `ai-color-palette` | warning | | — | (a) text `hasChroma(c, 50)` with hue ∈ **[260,310]** on `h1/h2/h3` or `fontSize ≥ 20`; (b) gradient background stop `hasChroma(c,50)` with hue ∈ [260,310] (purple) or **[160,200]** (cyan); (c) neon text `hasChroma(c, 80)` in either band on a bg with `relativeLuminance < 0.1`; (d) Tailwind `text-(purple\|violet\|indigo)-N` on heading/`text-[2-9]xl`, or `from-(purple\|violet\|indigo)-N` + `to-(purple\|violet\|indigo\|blue\|cyan\|pink\|fuchsia)-N`; (e) raw-HTML hex list `#7c3aed \|8b5cf6\|a855f7\|9333ea\|7e22ce\|6d28d9\|6366f1\|764ba2\|667eea`. |
| 8 | `cream-palette` | warning | | — | Page (body, else html) bg passes `isCreamColor`: `min(r,g,b) ≥ 209` AND `r ≥ g ≥ b` AND `6 ≤ (r-b) ≤ 48`. Tailwind fallback map `bg-amber-50/100, bg-orange-50/100, bg-yellow-50, bg-stone-50/100/200` plus arbitrary `bg-[…]`, each re-filtered through `isCreamColor`. |
| 9 | `nested-cards` | warning | | layout | A card-like element with a card-like ancestor. `isCardLikeFromProps(shadow,border,radius,bg) = (shadow\|\|border) && (radius\|\|bg)`. Skips SAFE_TAGS, form controls, media, `position:absolute\|fixed`, class matching `dropdown\|popover\|tooltip\|menu\|modal\|dialog`, text < 10 chars, rect < 50×30. **Only the innermost** flagged element is reported. |
| 10 | `monotonous-spacing` | warning | | layout | Harvest padding/margin px + rem(×16) + `gap:` + Tailwind `[pmg]{axis}-N` (×4). Round each to nearest 4. Need ≥ **10** samples. Fire if `dominantValue/total > 0.6` AND `uniqueValues.length <= 3`. |
| 11 | `bounce-easing` | warning | | — | `animation-name` matching `/bounce\|elastic\|wobble\|jiggle\|spring/i`; Tailwind `animate-bounce`; any `cubic-bezier(x1,y1,x2,y2)` with **`y1 < -0.1 \|\| y1 > 1.1 \|\| y2 < -0.1 \|\| y2 > 1.1`** (overshoot). |
| 12 | `pulsing-dot` | warning | | — | Merged per-selector CSS decls (across rule blocks, `prefers-reduced-motion: reduce` blocks stripped first). Needs: infinite animation whose `@keyframes` vary `opacity`/`box-shadow`/`transform: …scale` (rotation-only = spinner, exempt), or name matching `/pulse\|blink\|ping/i`; size `2 ≤ w,h ≤ 16`px; round (`border-radius ≥ 40%` or `≥ 999px` or `≥ 0.4*min(w,h)`). Tailwind: `animate-ping\|animate-pulse` + `rounded-full` + `[wh\|size]-(1\|1.5\|2\|2.5\|3\|3.5\|4)`. **Severity promoted to `error`** when the selector resolves inside a `<header>`/`<nav>` source range. |
| 13 | `blinking-cursor` | **advisory** | | — | Infinite animation named `/blink\|caret\|cursor/i` OR whose keyframes ONLY toggle `opacity ≤ 0.15` / `visibility:hidden` (any other animated prop disqualifies). Exempt: `input/textarea/select/img/svg/script/style`, `isContentEditable`, `[contenteditable]`, `[role=textbox]`. Position gate: `rect.top + scrollY ≤ 1200px` (`CURSOR_FIRST_VIEWPORT_PX`). Shape: single glyph in `/^[_\|▀-▟■▮❙❚｜]$/`, OR empty box with fill/border, vertical `1≤w≤24, 6≤h≤48, h≥w` or underscore `1≤h≤6, 4≤w≤24`, and `borderRadius < 0.4*min(w,h)`. **Promoted to `warning`** if `pageTop ≤ 900` or inside `header,nav,[role=banner],[role=navigation]`. |
| 14 | `shape-assembled-illustration` | **advisory** | | — | Inline `<svg>` with ≥ **8** `<rect\|circle\|ellipse\|polygon>`, ≥ **3** distinct non-`none/transparent/currentcolor/inherit` fills, intrinsic **≥ 200×200** (width/height attrs, else viewBox), ≤ **2** `<text\|tspan>` nodes, and no `<pattern>`. |
| 15 | `dark-glow` | warning | | — | Per shadow layer: color `hasChroma(c, 30)`, blur (3rd length) `> 4`. Fires if **(a) zero offset** `x===0 && y===0` on **any** background, or **(b)** any chromatic blurred shadow when effective bg `relativeLuminance < 0.1`. Applies to `box-shadow` and `text-shadow` (text-shadow skipped when inherited unchanged from parent). Text-scan dark-page heuristic (`cssTextHasDarkRootBg`): dark hex/rgb literal, Tailwind `bg-(gray\|slate\|zinc\|neutral\|stone)-(9xx\|800)`, or a `body/html/:root`-scoped (or `<body style>`) background resolving via `var()` to α>0.5 and luminance < 0.1. |
| 16 | `radial-halo` | warning | | — | On a dark root page only. Non-repeating `radial-gradient`, no `url()` layer, ≥ 2 color stops, **first stop** chromatic (RGB spread ≥ **24**) and α ≥ **0.7**, **last stop** α ≤ **0.05**, and no stop position ≤ 24px (dot/texture exemption). |
| 17 | `radial-spotlight-glow` | warning | | — | Translucent sibling of #16, disjoint alpha band. Last stop α ≤ 0.05; visible stops (α > 0.05) count ≤ **2**; **every** visible stop α < **0.45**; ≥1 visible stop `hasChroma(c, 24)`; surface **≥ 240×160**px. |
| 18 | `marquee` | warning | | — | `<marquee>` element, OR an `animation … infinite` bound to a `@keyframes` whose `translateX/translate/translate3d` **percentage** travel ≥ **20%** (pixel-travel loops deliberately exempt; a single X sample that also varies `scale`/`opacity` = pulse, exempt). |
| 19 | `icon-tile-stack` | warning | | layout | Heading's `previousElementSibling` is: 32–128px on **both** axes; aspect ratio 0.7–1.4; has visible bg (α>0.1) / bg-image / border; `borderRadius < width/2` (excludes circles/avatars); contains `svg`/`i[data-lucide]`/`i[class*=fa-]`/`i[class*=icon]` (or is childless with emoji-only text) whose width < 0.95×tile; and `siblingBottom ≤ headingTop + 4`. |
| 20 | `italic-serif-display` | warning | | type | `font-style: italic` AND (`h1`, or `h2` with size ≥ 48) AND `fontSize ≥ 48` AND primary face ∈ `KNOWN_SERIF_FONTS` **or** stack ends in generic `serif`. |
| 21 | `hero-eyebrow-chip` | warning | | type | `h1` with `fontSize ≥ 48`, not inside `[role=tabpanel\|dialog\|application], dialog`. `previousElementSibling` non-heading, text 2–60 chars, `0 < fontSize ≤ 14`. Then any of **three branches**: (A) tracked-caps — uppercase and `letterSpacing ≥ 1.6`px; (B) accent-bold — `fontWeight ≥ 700` AND `isAccentColor` (rgb/hex spread ≥ 40, oklch C ≥ 0.05, hsl S ≥ 20%); (C) dash-prefix — a `::before/::after` sized 8–80 × 1–6px with chromatic fill (spread ≥ 30). |
| 22 | `kicker-above-heading` | warning | | type | **Outright ban — every candidate is a finding, no repetition count.** Candidate gate (`isKickerCandidate`): heading level ≤ 4 (h1-h4 or `role=heading` w/ `aria-level`), heading text ≥ 3 chars, heading `fontSize ≥ 20`; kicker tag ∈ `p/span/div/small`; kicker text 2–34 chars; uppercase or `small-caps`; `0 < fontSize ≤ 14`; `letterSpacing ≥ max(1, fontSize*0.08)`. Rejects: `^step\s*\d+`, bare 1–2 digits, meta text `/[·•\|]\|\s[\/›»>]\s\|\b(19\|20)\d{2}\b/`, legal numbering `^(§\|\d+(\.\d+)+\|(section\|article\|clause\|appendix\|exhibit\|schedule\|chapter\|part\|rule\|title)\s+…)`. Skip selectors: `nav, form, table, thead/tbody/tfoot, figure, figcaption, ol, ul, li, [role=navigation], [aria-label*=breadcrumb], [class*=breadcrumb], [aria-hidden=true], [data-impeccable-allow-kickers]`, plus card contexts `article/button/a/li/[role=listitem]/[role=option]`, plus `[role=tabpanel\|dialog\|application], dialog`. **Stands down** when h1 ≥ 48px and tracking ≥ 1.6 (that is #21's territory). |
| 23 | `numbered-section-labels` | **advisory** | | type | Label before an `h2/h3/h4` (or before the wrapper the heading leads). Text shape: bare `^(\d{2})$` or `^(\d{1,2})[non-word sep]\S`, index ≤ 40, total ≤ 40 chars. Label `0 < fontSize ≤ 13`; heading ≥ 1.3× label size (when resolvable); label tag ∈ `span/p/div/small/em/strong/b`; deliberate styling required: mono family **or** weight ≥ 600 **or** tracking ≥ 0.5 **or** uppercase **or** accent color. **Needs ≥ 2 candidates on the page with ≥ 2 distinct indices.** |
| 24 | `em-dash-overuse` | warning | **YES** | — | Two gates, both must hold: `count ≥ EM_DASH_FLOOR (8)` and `bodyText.length ≤ count * EM_DASH_CHARS_PER_DASH (500)`. Counts `—` and `--(?=\S)`. Text path decodes `&mdash;` / `&#0*8212;` / `&#x0*2014;` first. En-dashes deliberately not counted. |
| 25 | `marketing-buzzword` | warning | | — | ≥ **1** occurrence of any of the 30 `BUZZWORDS` phrases in HTML-stripped body text. |
| 26 | `aphoristic-cadence` | warning | | — | ≥ **3** total matches of `\bNot an? [a-z][^.!?]{1,40}[.!]\s+[A-Z][^.!?]{1,60}[.!]` (manufactured contrast) + `\b[A-Z][^.!?]{4,80}[.!]\s+(No\|Just)\s+[a-z][^.!?]{2,60}[.!]` (short rebuttal). |
| 27 | `oversized-h1` | warning | | type | `h1` with `fontSize ≥ 72` (`OVERSIZED_H1_FONT_PX`) AND `textLen ≥ 40` (`OVERSIZED_H1_MIN_CHARS`). When a rect is available it must **also** dominate: `rect.height/vh ≥ 0.28` OR `area/viewportArea ≥ 0.25`. |
| 28 | `extreme-negative-tracking` | warning | | type | Direct text > 20 chars, `letterSpacingPx/fontSize ≤ -0.05` em. |
| 29 | `gpt-thin-border-wide-shadow` | **advisory** | | — | ≥ **2** visible thin borders (`0 < width ≤ 1.5px` AND border-color α ≥ **0.28**) AND max shadow blur ≥ **16px** among layers with color α ≥ **0.12**. |
| 30 | `repeating-stripes-gradient` | **advisory** | | — | Bare presence of `repeating-(linear\|radial\|conic)-gradient(` anywhere in the HTML. |
| 31 | `codex-grid-background` | **advisory** | | — | **Both signals in the SAME declaration block** (CSS rule body or one inline `style="…"`): hairline stop `\d{1,3}px\s*,\s*transparent\s+\d{1,3}px` (or inverted `transparent calc(100% - Npx)`) counted **only inside `background`/`background-image` values**, plus a px tiling cell. Fires on `(hairlineCount ≥ 2 && hasPxCell)` or `hasPxPairCell` (`background-size: Npx Npx` / shorthand `/ Npx Npx`). |
| 32 | `theater-slop-phrase` | **advisory** | | — | `\b(\w+)\s+theater\b` in HTML-stripped body text. |
| 33 | `image-hover-transform` | **advisory** | | — | CSS `img…:hover { transform: scale\|rotate\|translate\|matrix\|skew }` or Tailwind `hover:(scale\|rotate\|translate\|skew)-` on an `<img class>`. |

### Category `quality` — general design / a11y defects

| # | id | sev | scopes | Threshold / heuristic |
|---|---|---|---|---|
| 34 | `broken-image` | warning | — | `<img>` with no `src` attribute, or `src` ∈ `"" \| " " \| "#"`. |
| 35 | `script-error` | **error** | — | Puppeteer `pageerror` events (listener attached **before** `goto` to catch parse errors). Deduped by message, first line, truncated 160 chars, **max 3 reported**. |
| 36 | `content-hidden-at-rest` | **error** | layout | After a **reveal sweep** (instant scroll in steps of `max(200, 0.7*innerHeight)` to `scrollHeight`, rAF+40ms per step, back to top, 700ms settle): fires when `totalChars ≥ 200` AND `hiddenChars ≥ 150` AND `hiddenChars/totalChars > 0.30`. "Hidden" = computed `opacity ≤ 0.02` or `visibility: hidden/collapse` inherited down the chain. `display:none`, `[hidden]`, `aria-hidden`, `content-visibility:hidden`, and `script/style/noscript/template/title/head/meta/link/option/optgroup/select/datalist/dialog` are **excluded from the denominator**, not counted as invisible. |
| 37 | `edge-flush-cards` | warning | layout | Horizontal scroller (`overflow-x: auto\|scroll`) with `scrollWidth > clientWidth + 8`, `scrollLeft ≤ 4` (at rest only), rect ≥ 120×60, `top + scrollY ≤ 2*vh`. Card: own bg α > 0.5 or ≥ 2 borders, rect ≥ 80×40, attributed to nearest scroller. Fires when `leftGutter ≥ 6 && -24 < rightGap < 8` (or mirror). The `-24` floor exempts deliberate peeking next-cards. |
| 38 | `text-occlusion` | warning | layout | Three ground-truth paths, all browser-only, viewport-bound. **(i) elementFromPoint probe**: grid of `cols = clamp(6..30, width/12)` × `rows = clamp(1..4, height/14)` samples; occluded if the top element is an "opaque decorated box" (bg α > 0.6 or ≥2 borders with color α > 0.3) or has its own text. Fire threshold `occFrac ≥ 0.30` for boxes, `≥ 0.45` for text. Skips floated, marquee-ish (`marquee\|ticker\|scroller\|carousel\|conveyor`), fixed/sticky overlays, `img/video/canvas/picture` (contrast territory). Text-on-text requires ≥1 side out of flow (`absolute\|fixed\|sticky` on the chain) — otherwise it is line-box bleed. **(ii) headline overhang**: `fontSize ≥ 40`, card bg α > 0.7 (no gradient/url), card 100px ≤ w ≤ 0.8vw, h ≥ 60; intersection `ix ≥ 8 && iy ≥ 0.5*lineHeight`, headline centerX outside the card, `ix ≤ 0.5*textWidth`. **(iii) inline padding leak**: `display: inline`, bg α > 0.6, `paddingTop+paddingBottom ≥ 24`, rect ≥ 12×24, `rect.height ≥ 2.2 * lineHeight`. |
| 39 | `first-viewport-column-overflow` | warning | layout | A grid/flex container ≥ 0.5vw whose `pageTop < 0.9vh` and `pageBottom > vh`; ≥ 2 direct children with `0.25 ≤ widthShare ≤ 0.9`, height ≥ 40; tops within `0.25*vh` of each other; tallest child's **content extent** > `1.4*vh` while the shortest's ≤ `vh`. |
| 40 | `gray-on-color` | warning | — | Text `!hasChroma(c, 20)` with `0.05 < luminance < 0.85`, on a bg where **every** stop `hasChroma(bg, 40)`. Tailwind: `text-(gray\|slate\|zinc\|neutral\|stone)-N` + `bg-(16 hues)-N` on the same line. |
| 41 | `low-contrast` | warning | — | WCAG AA. `threshold = 3.0` if `fontSize ≥ WCAG_LARGE_TEXT_PX (18*96/72 = 24px)` or (`fontSize ≥ WCAG_LARGE_BOLD_TEXT_PX (14*96/72 ≈ 18.67px)` and `weight ≥ 700`), else `4.5`. Gradient bgs measured against the **worst** stop. Also a **:hover-state pass** (`checkHoverContrast`) for styled controls (own bg α > 0.5). Also a **pixel-diff pass** — see §2. |
| 42 | `layout-transition` | warning | — | `transition-property` (not `all`/`none`) naming any of `width, height, padding, margin, max/min-height, max/min-width`, or any longhand padding-*/margin-*. |
| 43 | `line-length` | warning | type,layout | Browser-only (needs rect). `hasDirectText`, tag ∈ `p,li,td,th,dd,blockquote,figcaption`, `textLen > lineMax`. `charsPerLine = rect.width / (fontSize * 0.5)`; fires when `> lineMax + 5`. `lineMax` default **80**, overridable via `window.__IMPECCABLE_CONFIG__.lineLengthMax`. |
| 44 | `cramped-padding` | warning | layout | **Two shapes.** (a) *Own text*: rect gate `textLen > 20 && width > 100 && height > 30`, needs ≥2 borders or a visible bg boundary; thresholds **`vThresh = max(4, fontSize*0.3)`**, **`hThresh = max(8, fontSize*0.5)`**; one finding per element (worse axis wins). Inline `<code>` (not in `<pre>`) exempt. (b) *Flush children*: container with **no** direct text, ≥1 child, not `absolute/fixed`, tag not in a 22-item skip set (HTML/BODY/MAIN/HEADER/FOOTER/NAV/ARTICLE/ASIDE/BUTTON/A/LABEL/SUMMARY/CODE/PRE/INPUT/TEXTAREA/SELECT/FORM/FIGURE/TABLE/TBODY/THEAD/TR/TD/TH); visible boundary (border/outline/bg); `pad[side] ≤ PAD_THRESHOLD (2)`; **no child insulates that side** (`CHILD_INSULATE_THRESHOLD = 4`px of child padding *or* margin *or* rect gap); text descendants actually flush (`TEXT_EDGE_THRESHOLD = 4`px); ≥1 child with > 4 chars. Full-bleed bg bands (`rect.width ≥ 0.94*viewport`, bg-only) don't bound left/right. |
| 45 | `body-text-viewport-edge` | warning | layout | `<p>`/`<li>` only, `textLen > 40`, `rect.width/viewport > 0.5`, `rect.left < 16` or `rect.right > viewport - 16`; not inside `nav`/`header`, no own bg-color, not positioned. |
| 46 | `tight-leading` | warning | type | Non-heading, direct text, `textLen > 50`, `lineHeightPx/fontSize < 1.3`. |
| 47 | `skipped-heading` | warning | type | Document-order walk of `h1…h6`; fires when `level > prevLevel + 1`. |
| 48 | `heading-rhythm` | warning | layout,type | Browser-only, measures real rects. For `h2/h3/h4`: computes nearest content edge above the **heading cluster** (up to 3 label-like previous siblings absorbed: gap < 28px, sibling height ≤ 60, label-like = `sibFontSize < 0.75*headingFontSize` or text ≤ 40 chars, text ≤ 80 chars) and below. Requires x-overlap ≥ 8px. `6 ≤ below ≤ MAX_BELOW_PX (160)`. Violation: `above < below * 0.75 && (below - above) ≥ MIN_DEFICIT_PX (12)`. **Fires only at `MIN_VIOLATIONS = 2`.** Exempt: inside a card < `CARD_EXEMPT_HEIGHT (200)`px; container drawing its own top boundary (bg α>0.05, top border, shadow). |
| 49 | `justified-text` | warning | type | `text-align: justify` with `hyphens !== 'auto'`. |
| 50 | `tiny-text` | warning | type | Direct text, `textLen > 20`, `fontSize < 12`. Excluded tags `sub,sup,code,kbd,samp,var,caption,figcaption`; excluded contexts (huge selector) `button,a,label,summary,pre,[role=button/link/tab/menuitem/option],nav,footer,[aria-hidden=true],[class*=badge/caption/chip/code/console/diff/label/meta/mock/pill/preview/tag/terminal/writes]`; uppercase exempt; `isNonRenderedText` exempt. |
| 51 | `undersized-ui-text` | warning | type | **Deliberately ignores the design system.** Direct text ≥ 2 chars, `0 < fontSize < 11`. Floor **11px**, softened to **10px** only for non-interactive smallprint (`small, footer, [class*=legal/copyright/fineprint/fine-print/smallprint/small-print/disclaimer/disclosure/footnote]`). Fires when interactive (`a[href],button,summary,label,select,textarea,[role=button/link/tab/menuitem/menuitemcheckbox/menuitemradio/option/checkbox/radio/switch/treeitem],[tabindex]`) OR furniture (`nav,[role=navigation],td,th,[role=gridcell/cell],caption,figcaption,dt,dd,footer,[class*=meta/label/badge/chip/pill/tag/kicker/eyebrow/breadcrumb/timestamp/category/caption/nav]`) OR direct text ≤ 20 chars. Exempt: `sub,sup,option`; `pre,code,kbd,samp,var,svg,[aria-hidden],[class*=terminal/console/code/mock/editor/syntax/diff]`; sr-only (`isVisuallyHidden`: class match on 13 sr-only spellings, or absolute/fixed + `clip: rect(0…)` / `clip-path: inset(50%\|99\|100%)` / 1px box + overflow hidden). |
| 52 | `all-caps-body` | warning | type | Non-heading, direct text, `textLen > 30`, `text-transform: uppercase`. |
| 53 | `wide-tracking` | warning | type | Direct text, `textLen > 20`, not uppercase, `letterSpacingPx/fontSize > 0.05` em. |
| 54 | `text-overflow` | warning | layout | Browser-only. Element owns direct text; not `pre,code,textarea,svg,canvas,select,option,marquee`; not sr-only; neither it nor any ancestor is a scroll region. Fires when `scrollWidth - clientWidth ≥ 16`. Inline fallback (`clientWidth === 0`): compare inline rect right against nearest block container's padding-box right, `spill ≥ 16`; bails if any `transform` on the path. |
| 55 | `repeated-container-text` | warning | — | Same literal 4–48-char text with ≥1 letter appearing **≥ 3 times at ≥ 3 distinct structural signatures** inside one bounded container. Signature = element path (`tag.sortedClasses>…`) from occurrence to container — so parallel/templated repetition collapses to one signature and never fires. Container: tag ∈ `div,section,article,aside,main,figure,form,fieldset,details,li`, passes `isCardLikeFromProps` with `hasBorder = ≥3 sides ≥ 1px`; ≤ **250** descendants. Skips `table,select,datalist,nav,menu,[role=navigation/menu/menubar/listbox/grid/tablist/radiogroup],[aria-hidden]` and icon-font classes. Text attributed to innermost container only. |
| 56 | `clipped-overflow-container` | warning | layout | `overflow(-x/-y): hidden\|clip` with **no** `auto\|scroll` axis; contains an `absolute\|fixed` descendant. Exempt: `aria-hidden`, `role=none/presentation`, `img/svg/canvas/video`, decorative name tokens (`art\|bg\|background\|badge\|blob\|crop\|decor\|dot\|glow\|grain\|image\|mask\|ornament\|overlay\|photo\|scrim\|shadow\|shine\|texture`) without substantive content; container names `carousel\|comparison\|compare\|fisheye\|marquee\|preview\|scroller\|slider\|slideshow\|split\|viewport\|demo-area\|demo-stage\|demo-viewport` or `aria-roledescription` carousel/slider. Requires geometric escape (> 2px past the parent rect) or, when rects are unmeasurable, a negative/100% inset value. |
| 57 | `design-system-font` | warning | type | Primary font not in the DESIGN.md typography allowlist. |
| 58 | `design-system-color` | **advisory** | — | Color literal not within `COLOR_CHANNEL_TOLERANCE = 6` per-channel of any allowlisted color (frontmatter `colors.*` + sidecar `extensions.colorMeta.{canonical,tonalRamp[]}`). `var()`, `transparent/currentcolor/inherit/initial`, and α ≤ 0.05 abstain. Source scan additionally requires `isProbablyColorLiteral` context (a CSS color-ish property, a gradient/color-mix function, or a JS color key) and rejects `&#…` entities and `>#155<` prose. |
| 59 | `design-system-radius` | **advisory** | — | `border-radius` token off the allowlist by more than `RADIUS_TOLERANCE_PX = 0.5`. `var()`, `%`, `0/none/initial/inherit`, and `px ≤ 0.5` abstain. If the system declares a pill token (`full\|pill\|round\|rounded-full`), any `px ≥ 99` is allowed. |
| 60 | `design-system-font-size` | **advisory** | type | Literal `px`/`rem` font-size (`FONT_SIZE_LITERAL_RE = /^-?[\d.]+(?:px\|rem)$/`) off the ramp by more than `FONT_SIZE_TOLERANCE_PX = 0.5`. `clamp()` values are judged **only at their two fixed endpoints**. Abstains entirely (`hasFontSizes = false`) when the system declares *only* clamp endpoints and no enumerated `typography.scale` — "a fully fluid system has no discrete ramp; abstain instead of flagging every intermediate size". em/%/calc/var abstain. Source-scan only (computed font sizes cascade into off-ramp px in a browser). |

---

## (b) PROCESS — workflow / phases / state

There is barely any. The detector is **stateless and single-shot**: `impeccable detect [flags] [targets…]`. The only "process" artifacts:

- **Target routing** (`detector/cli/main.mjs`): `^(?:https?|file)://` → browser engine; `.html/.htm` → static-HTML engine; everything else → regex engine; no target + non-TTY stdin → `handleStdin` (parses a Claude-Code-hook JSON payload `tool_input.file_path`, else treats stdin as source text).
- **Design-system resolution is per-target, not per-cwd** (`findDesignRoot`): walk up from the *target's own* directory; a dir with `DESIGN.md`/`Design.md`/`design.md` (or in `.agents/context/` or `docs/`) IS the design root; a dir with `.git`/`package.json`/`.impeccable` but no DESIGN.md is a **boundary** — the walk stops with *no* design system; reaching `$HOME` or `/` returns null. Memoized by resolved root. Comment: *"This is the fix for cross-project contamination."*
- **Dev-server nudge**: `detectFrameworkConfig` finds `next/svelte/nuxt/vite` configs, probes the port with an HTTP fingerprint (`x-powered-by: next`, `x-sveltekit-page`, …), and prints to stderr *"For more accurate results, scan the running site: npx impeccable detect http://localhost:3000"* — an explicit statement that **static analysis is second-best**.
- **Large-scan consent**: > 50 files on a TTY → interactive `Continue? [Y/n]`.
- **`doctor.mjs`** is the only real state machine: a drift/staleness pass over the project's own artifacts with a three-tier severity ladder `{ auto: 'automatic', mention: 'worth saying', route: 'needs a command' }`. `--fix` applies **only** `severity: 'auto'` items — *"the ones with no judgment in them"*. Everything else is reported and left alone. Notably it refuses to delete `legacy-live-state` even though it is `auto`: *"losing session state to a doctor run is a worse outcome than a stale file."* Exit code is 0 unless the run itself failed — *"findings are not errors."*
- **`detect-csp.mjs`** is a pure classifier: greps ≤ 64KB of the first 6 directory levels and returns one of `append-arrays | append-string | middleware | meta-tag | null`, prioritized *"structured patches are safer than string splices"*. It **detects but never patches**; the agent writes the patch after a user consent prompt.

---

## (c) MACHINERY — engines, merge, suppression, scoring

### 2. THE FOUR DETECTION ENGINES — what each can and cannot see

Declared capability matrix (`RULE_ENGINE_SUPPORT` in the registry):

```js
const RULE_ENGINE_SUPPORT = {
  regex:         new Set(['source', 'page-analyzer']),
  'static-html': new Set(['element', 'page']),
  browser:       new Set(['element', 'page', 'layout']),
  visual:        new Set(['visual-contrast']),
};
```

| engine | entry | sees | **cannot** see |
|---|---|---|---|
| **regex** | `engines/regex/detect-text.mjs` | Raw text of `.css/.scss/.sass/.less/.jsx/.tsx/.js/.ts/.vue/.svelte/.astro`. `REGEX_MATCHERS` (line-scoped, ±3-line context for CSS-like files) + `REGEX_ANALYZERS` (7 whole-file passes). Extracts `<style>` blocks from Astro/Vue/Svelte and CSS-in-JS from `` styled.x`…` `` / `` css`…` `` templates, with line-offset bookkeeping. | Any cascade. Any computed value. Any layout. Whether a rule matches an element that exists. `var()` beyond a single-level in-file resolution. **This is the only engine that reports line numbers**, which is why inline ignores are precise here. |
| **static-html** | `engines/static-html/detect-html.mjs` + `css-cascade.mjs` (1183 lines, hand-rolled) | Parses with `htmlparser2` + `css-select` + `css-tree` + `domutils` (jsdom was **removed**). Follows and inlines linked CSS. Resolves `var()` chains up to 8 levels, oklch/oklab/hsl/hwb/color-mix→sRGB, composites translucent layers. Runs 14 `STATIC_ELEMENT_RULES` + 8 page rules + text-content analyzers. Simulates a `:hover` pass (`window.getHoverStyle`), a full-cover pseudo-surface pass (`window.getPseudoSurface`), and an accent-dash pseudo pass (`window.hasAccentDashPseudo`). | **No layout at all.** `rect: null` is passed explicitly, disabling `line-length` and the rect half of `cramped-padding`. Dimensions come from explicit CSS `width`/`height` only. Comments record: `@media`/`@supports` wrapped rules ignored, only `:root`-level custom properties resolved, "unsupported base selector" swallowed silently, `::before/::after` "can't join the element cascade". Findings carry **line 0**. Falls back to `detectText` if the four parser packages are missing. |
| **browser** | `engines/browser/detect-url.mjs` + `browser/injected/index.mjs` | Real Chrome via **Puppeteer** (`headless: true`, viewport default `1280×800`, `waitUntil: 'networkidle0'`, `timeout: 30000`). Only engine with real layout, `getBoundingClientRect`, `scrollWidth`, `elementFromPoint`, `elementsFromPoint`, live CSSOM `@keyframes`, `pageerror` events, and the ability to scroll. **9 rules are browser-exclusive**: `heading-rhythm`, `edge-flush-cards`, `text-occlusion`, `first-viewport-column-overflow`, `text-overflow`, `blinking-cursor`, `content-hidden-at-rest`, `script-error`, and the rect gates of `line-length` / `oversized-h1`'s viewport-dominance test. | Cross-origin stylesheets (`try { sheet.cssRules } catch { continue }`). Anything below the fold for the `elementFromPoint` probes (explicitly "viewport-bound"). Source line numbers. Needs a running server or a `file://` URL. |
| **visual (pixel)** | `engines/visual/screenshot-contrast.mjs` | Screenshots a clipped region (`height` hard-capped at **320px**), then re-screenshots with `color: transparent !important; -webkit-text-fill-color: transparent !important; text-shadow: none !important` applied to the target, and diffs the two PNGs pixel-by-pixel in a `<canvas>`. A pixel counts as glyph ink when the 4-channel absolute delta `≥ 10`. Needs ≥ **8** glyph pixels and ≥ **8** ratio samples. Reports the **p10 ratio** (not the min — noise-robust) against the same 3.0/4.5 WCAG threshold, with the median printed alongside. | Anything not already a *candidate*. Candidates are gated to ≤ **12** (`visualContrastMaxCandidates`) elements that have a "reason" to be unmeasurable statically: `background-clip text`, `text shadow`, image/gradient background ancestor, `opacity < 0.99` stack, `mix-blend-mode`, `filter`, `backdrop-filter`, or an `img/picture/video/canvas/svg` underlay found via `elementsFromPoint`. |

**Order of operations for a URL scan** (`detectUrl`): attach `pageerror` listener → `setViewport` → `goto` → optional `settleMs` → set `window.__IMPECCABLE_CONFIG__ = {autoScan:false, designSystem}` → inject the 8250-line bundle → `window.impeccableDetect({decorate:false, serialize:true})` → **reveal sweep + `content-hidden-at-rest`** (deliberately *after* the main scan, which must see the true at-rest state) → push ≤3 `script-error` → **visual contrast fallback last** (the sweep restores scroll to top).

### How findings are merged / deduped

There is **no cross-engine merge**. Each target is scanned by exactly one engine; `allFindings` is a flat concatenation. Dedupe is local and layered:

1. **Regex engine** (`detectText`): drop a finding if an existing one has same `antipattern` + identical `snippet` + `|Δline| ≤ 2`.
2. **Design-system findings** (`mergeDesignSystemFindings`): canonical key per rule — font → normalized name + google/non-google context; color → `r,g,b` after parsing; radius/font-size → px rounded to 2 decimals. Static (line 0) and source (real line) results merge, and **the merged entry inherits the real line number**. `dedupeDesignFindings` additionally keys on `antipattern + line + normalized ignoreValue`.
3. **Per-scan `seen` sets**: `seenFonts`/`seenColors`/`seenRadii` in `collectStaticDesignSystemFindings` report each violating *value* once per file, not once per element. Same in the browser (`designSeen`).
4. **Rule stand-downs** (the cleanest idea here — rules explicitly yield to each other so one defect gets one finding):
   - `kicker-above-heading` skips when h1 ≥ 48px & tracking ≥ 1.6 → that is `hero-eyebrow-chip`.
   - `tiny-text` (long body copy) and `undersized-ui-text` (functional/short text) partition the space so they never double-flag.
   - `radial-halo` (first stop α ≥ 0.7) and `radial-spotlight-glow` (all stops α < 0.45) use **disjoint alpha bands** — *"so they never double-report the same declaration"*.
   - `blinking-cursor` yields round dots to `pulsing-dot` (`radiusPx ≥ 0.4*min(w,h)` bail).
   - Visual-contrast candidates already flagged `low-contrast` by the DOM pass are filtered out by selector.
   - Nested cards: only the **innermost** flagged element reports.
   - `text-shadow` glow is skipped when the parent has the identical value (inheritance).
5. **Browser grouping**: findings are keyed by element in a `Map` (`addBrowserFindings`), then serialized as `{selector, tagName, rect, isPageLevel, isHidden, findings[]}`. `hero-eyebrow-chip` is deliberately **re-attributed to `el.previousElementSibling`** so the highlight lands on the eyebrow, not the heading.

### 3. SEVERITY, CONFIDENCE, FALSE-POSITIVE SUPPRESSION

**Severity** is 3-valued and mostly cosmetic: `finding()` stamps `severity: ap.severity || 'warning'`. Values observed: `'error'` (2 rules), `'advisory'` (11 rules), `'warning'` (default, 47 rules). Two rules **promote severity at runtime based on position**:
- `pulsing-dot` → `'error'` when its selector resolves inside a `<header>`/`<nav>` source range (`selectorHitsLandmark` / `landmarkSourceRanges`).
- `blinking-cursor` → `'warning'` (up from advisory) when `pageTop ≤ 900` or inside `header,nav,[role=banner],[role=navigation]`.
`detectUrl` honors the promotion: `if (f.severity && f.severity !== item.severity) item.severity = f.severity;`

**⚠️ A real inconsistency worth knowing.** `severity: 'advisory'` and `advisory: true` are *different things*:
```js
const ADVISORY_RULE_IDS = new Set(ANTIPATTERNS.filter(rule => rule.advisory === true).map(r => r.id));
```
Only **`em-dash-overuse`** sets `advisory: true`. So the 11 rules labeled `severity: 'advisory'` (`blinking-cursor`, `shape-assembled-illustration`, `numbered-section-labels`, `design-system-color/radius/font-size`, `gpt-thin-border-wide-shadow`, `repeating-stripes-gradient`, `codex-grid-background`, `theater-slop-phrase`, `image-hover-transform`) are **still counted as failures and still exit 2**. The comment claims the set is "derived from the registry so a rule only needs `advisory: true`" — but the registry mostly uses the other key.

**Confidence**: there is no numeric confidence field. Confidence is expressed structurally, via a consistently stated doctrine: **when a value cannot be resolved, skip rather than guess.** Verbatim:
- *"Unresolvable colors (currentColor, external vars): don't guess."*
- *"Skipping a rule beats a false positive here — the true painted contrast can't be measured from `color`."*
- *"a dropped stop can't manufacture a false finding, and skipping beats a wrong ratio."*
- *"conservative: no promotion without placement evidence."*
- `fontSizeStepStatus` returns a literal third state **`'unjudgeable'`** for `var()`/`calc()`/`%`/`em`.
- The one deliberate inversion: `isNeutralColor()` returns `false` (= chromatic = detect) for unknown formats — *"err on the side of DETECTING rather than silently skipping. This is the opposite of the previous default, which was the root cause of the oklch bug."*

**Suppression — three independent layers:**

1. **Inline ignores** (`shared/inline-ignores.mjs`), eslint-flavored and **comment-syntax-agnostic** (a raw token matched anywhere on a line; works in `//`, `/* */`, `<!-- -->`, `#`, `{/* */}`, `{# #}`):
   ```
   impeccable-disable <rule>[, <rule>…] [-- reason]     whole file
   impeccable-disable-line <rule>…      [-- reason]     same line
   impeccable-disable-next-line <rule>… [-- reason]     following line
   impeccable-disable                                   bare / `*` = every rule
   ```
   Reason separator is `--` (eslint) or `:` (biome); the reason is parsed off and **discarded** (*"self-documenting in the diff"*). Trailing comment closers (`*/`, `*/}`, `-->`, `*}`, `#}`, `%>`, `}}`) are stripped so they don't leak into the rule list. Rationale (verbatim): *"A config ignore is the right default for repo-wide policy. This complements it for the one case config can't cover: a waiver that belongs to a single file and needs to follow that file when it leaves the repo."* Because static-HTML findings carry line 0, **only whole-file directives apply there**.
2. **Project config** `.impeccable/config.json` + `config.local.json`, key `detector`: `ignoreRules[]` (rule ids), `ignoreFiles[]` (globs, matched against raw/absolute/relative paths), `ignoreValues[]` (rule + *value* pairs, optionally file-scoped, `"*"` wildcard for rules with no extractable value like `side-tab`), `designSystem.enabled`, `advisoryRules: 'include'|'exclude'`.
3. **CLI flags**: `--no-config` (kills all three), `--no-inline-ignores`, `--no-design-system`, `--no-advisory`, `--scope type|layout`.

**Structural FP suppression** is the bulk of `checks.mjs` and is worth studying as a genre. Recurring devices: tag allowlists (`SAFE_TAGS`, `BORDER_SAFE_TAGS`, `QUALITY_TEXT_TAGS`, `NON_RENDERED_TAGS`, `HIDDEN_TEXT_EXCLUDE_TAGS`, `TEXT_OVERFLOW_SKIP_TAGS`, `FLUSH_SKIP_TAGS`); ARIA-role carve-outs (status/alert regions, tabpanel/dialog/application, aria-selected="true", aria-current); class-name heuristics (`badge|chip|pill|terminal|console|mock|editor|syntax|diff`); an explicit opt-in escape hatch attribute `[data-impeccable-allow-kickers]`; screen-reader-only detection in two idioms; extension-noise skips (`id.startsWith('claude-')`, `'cic-'`, `[id^="impeccable-live-"]`); and `prefers-reduced-motion: reduce` block stripping so an a11y fallback can't mask the default experience.

Every FP fix is documented **with the bug that caused it** — e.g. *"the shipped miss: a `<span>` severity chip whose white text lost a specificity fight and rendered muted-on-red at 1.2:1"*, *"issue #407 — every Shopify product form ships an `<input name="id">"*, *"issue #408: dozens of phantom '10px body text' findings on every Shopify page"*, *"issue #409 Case A/B"*, *"issue #394"*, *"issue #303"*.

### 4. THE SCORING / PROFILE MODEL

**There is no score. There is no grade. There is no 0–100.** `detector/profile/profiler.mjs` is *not* a design profile — it is a **performance profiler**: `createDetectorProfile()` returns `{events: []}`, and `summarizeDetectorProfile()` groups by `engine|phase|ruleId|target` and emits `{calls, totalMs, avgMs, p50, p95, findings}` sorted by `totalMs`. Purely for `bun run bench:detector`.

The gate is a **binary count-and-exit**:

```js
const { primary, advisory } = partitionAdvisory(allFindings);
if (allFindings.length > 0) { …print…; process.exit(primary.length > 0 ? 2 : 0); }
process.exit(0);
```

- **exit 2** = at least one non-`advisory:true` finding.
- **exit 0** = clean, or advisory-only.
- Output: text (stderr, grouped by file, `line N: [rule-id] snippet` + `→ description`, plus a dimmed `── Advisory (not counted as failures) ──` section) | `--json` (stdout, full array with `advisory: true` stamped) | `--quiet` (count only).
- No severity weighting. `error` and `warning` are identical for exit-code purposes.

---

## 5. RUNTIME DEPENDENCIES — what a vendor-neutral protocol could NOT assume

| requirement | where | hard? |
|---|---|---|
| **Node ≥ 22.12.0** | `package.json` `engines` | Hard for everything. |
| ESM only (`"type": "module"`) | package.json | Hard. |
| `htmlparser2 ^12`, `css-select ^7`, `css-tree ^3.2.1`, `domutils ^4` | `dependencies` | **Required for the static-HTML engine.** `detectHtml` wraps the dynamic import in `try/catch` and **silently degrades to `detectText`** (regex) if any is missing — i.e. HTML files lose ~40 of 60 rules with no warning. |
| **`puppeteer ^25.1.0`** + a downloadable Chromium | `optionalDependencies` | **Required for URL scanning and 9 browser-exclusive rules.** Missing → hard throw: `'puppeteer is required for URL scanning. Install: npm install puppeteer'`. `process.env.CI` triggers `['--no-sandbox','--disable-setuid-sandbox']`. |
| A **running dev server** or `file://` URL | `detectUrl` | The browser engine cannot scan a bare filesystem path; the CLI prints an explicit nudge to start the dev server. |
| Canvas 2D + `getImageData` in-page | `screenshot-contrast.mjs` | For the pixel pass only. |
| Network | — | **None**, except whatever the scanned page itself loads. No API calls, no telemetry, no model calls anywhere in the detector. `detect-csp.mjs` is explicitly *"Mechanical (grep-based) — no network, no dev server, no JS evaluation."* |
| Writable FS | `doctor.mjs --fix` only | Detector is read-only. |
| `marked`, `fflate` | package.json | Not used by the detector path. |

---

## (d) SINGLE-AGENT ASSUMPTIONS that would NOT transfer

1. **One writer, one artifact.** Every finding is `{antipattern, name, description, severity, category, file, line, snippet}` — there is **no author/agent field**, no provenance, no "who wrote this". A protocol where 5 agents each write their own artifact cannot attribute or route findings with this shape.
2. **Findings are advice to a single agent inside one session.** The output is stderr prose + an exit code. There is no consensus notion, no vote, no signature, no "two of five reviewers flagged this".
3. **DESIGN.md is a single project-scoped truth**, resolved by walking up to one root, memoized once. Multi-agent deliberation would need *competing* proposed design systems and a merge/ratify step; `findDesignRoot` explicitly hard-stops at the first project boundary specifically to prevent multiple systems being in play.
4. **Suppression is unilateral.** Any writer can add `impeccable-disable overused-font` and the finding vanishes for everyone, with no record beyond a git diff and an optional free-text reason that the tool **discards at parse time**. In a cross-review protocol a waiver is exactly the thing that needs a counter-signature.
5. **Runtime-heavy verification.** The highest-value rules (occlusion, heading rhythm, hidden-at-rest, edge-flush, column overflow, pixel contrast) need Chrome + a running app. A vendor-neutral protocol whose participants are headless CLI agents on arbitrary machines cannot assume any of that, and cannot assume they all see the *same* rendered page.
6. **Node/JS/web monoculture.** Every threshold is CSS-shaped (`px`, `oklch()`, `border-radius`, Tailwind class regexes, `@keyframes`). Nothing in the machinery generalizes to a TUI, a native app, a slide deck, or a design doc.
7. **Hook-shaped I/O.** `handleStdin` parses `{tool_input: {file_path}}` — a Claude-Code PostToolUse hook payload. That's a single-harness integration.
8. **`--fix` philosophy is inherently single-actor**: *"applies only the migrations that carry no decision"*. In a multi-agent protocol, "no decision" is itself a decision someone must ratify.
9. **The extensive FP-suppression corpus is empirical and site-specific** (Shopify, Tailwind v4, jsdom bugs, Chrome oklch serialization). It encodes years of one maintainer's bug reports; a fresh multi-agent skill has no such corpus and would ship the false positives.
10. **Scale gates assume a whole page**: `≥ 20 text elements` for `overused-font`, `≥ 10 spacing samples`, `≥ 200 chars` for hidden-text, `isFullPage()` (`<!doctype|<html|<head`) gating all page analyzers. Agents reviewing a *component* or a *design proposal document* would trip none of them.

---

## Transferable to parley-design

Ranked by value-to-effort for a vendor-neutral, multi-agent, no-runtime protocol.

1. **Steal the registry as a doctrine artifact, verbatim in shape.** A flat, logic-free list of `{id, category: slop|quality, scopes, severity, name, description}` where the description states *the tell, why it reads as AI, and the corrective move*. 60 entries, ~1 screen each. This is the single most portable artifact in the repo — it needs **no runtime at all** and is directly usable as a shared checklist that every agent reads before writing, and as the vocabulary for cross-review comments. Reproduce §1 of this digest as `parley-design/references/antipatterns.md`.
2. **Adopt the `slop` vs `quality` split.** "This looks machine-made" and "this is objectively broken" are different arguments and deserve different burdens of proof in a consensus protocol. Slop findings are taste (arguable, quorum-able); quality findings are defects (a single agent should be able to BLOCK on `low-contrast` or `text-occlusion`).
3. **Adopt the "rules stand down for each other" discipline.** Explicit, documented yields (`kicker` yields to `hero-eyebrow`; `tiny-text` and `undersized-ui-text` partition the space; `radial-halo`/`radial-spotlight` use disjoint alpha bands) mean one defect produces exactly one finding. In a 5-agent review the same discipline stops five agents raising the same objection under five names — encode it as a **finding-taxonomy rule**: each defect has exactly one owning rule id, and the rule doc names who it yields to.
4. **Steal the "skip, don't guess" doctrine + the `unjudgeable` third state.** Findings should be `flag | pass | cannot-judge`, never a guess. This maps perfectly onto a multi-agent protocol where a participant that lacks a browser must say "I cannot judge this rule" rather than abstain silently or guess. `fontSizeStepStatus` → `'on-ramp' | 'off-ramp' | 'unjudgeable'` is the exact pattern.
5. **Steal the escape-hatch-closing rule (`undersized-ui-text`) as a protocol principle.** *"Being ON the DESIGN.md size ramp does not exempt a value here: adding 8px to the ramp launders the token but not the legibility problem."* Multi-agent design systems will absolutely be gamed by an implementer widening the token set to legalize its own output. **Have at least one rule class that is deliberately design-system-blind.**
6. **Steal the DESIGN.md → allowlist model, with the tolerances.** Frontmatter `typography` (per-role `fontFamily`/`fontSize` + an enumerated `scale`), `colors`, `rounded`, plus a JSON sidecar (`.impeccable/design.json`) carrying tonal ramps. Tolerances: color ±6 per channel, radius ±0.5px, font-size ±0.5px, pill radius = any ≥ 99px. And the **abstention rule**: if the system declares only fluid `clamp()` endpoints and no enumerated ramp, `hasFontSizes = false` and the rule does not run. This gives parley-design a machine-checkable, human-authored contract that every agent can read — no runtime.
7. **Steal inline-ignore syntax, but require a counter-signature.** `impeccable-disable-next-line <rule> -- <reason>`, comment-syntax-agnostic (raw token match, trailing-closer stripping). Change one thing: **do not discard the reason** — in a Parley protocol the reason is the artifact, and the waiver should require a second agent's ack recorded in the deliberation.
8. **Steal the ~15 rules that are checkable with zero runtime**, as a "static tier" any agent can run by reading source: `overused-font`, `single-font`, `flat-type-hierarchy` (ratio < 2.0), `monotonous-spacing` (≥10 samples, >60% dominant, ≤3 unique), `gradient-text`, `ai-color-palette` (hue 260-310 / 160-200), `cream-palette` (`min≥209, r≥g≥b, 6≤r-b≤48`), `bounce-easing` (cubic-bezier y outside [-0.1,1.1]), `marquee`, `em-dash-overuse` (8 + 1-per-500), `marketing-buzzword` (the 30-phrase list), `aphoristic-cadence` (≥3), `theater-slop-phrase`, `repeating-stripes-gradient`, `codex-grid-background`, `broken-image`, `justified-text`, `layout-transition`, `side-tab` regex forms. **These need only text, and every headless agent has text.**
9. **Steal the doctor's three-tier severity ladder** (`auto` = mechanical / `mention` = worth saying / `route` = needs a command) and the principle *"findings are not errors"* (`doctor` exits 0 regardless). It maps cleanly to a deliberation: auto-fix, note-in-review, escalate-to-consult.
10. **Steal the "gate on the running artifact when you can" nudge.** The CLI actively tells the user that static analysis is second-best and prints the exact command to do better. parley-design should encode the same honesty: a rule doc should state which tier of evidence it needs (text / DOM / layout / pixels) so an agent knows the confidence of its own verdict.
11. **Steal the position-based severity promotion idea.** Same defect, worse crime in the hero: `pulsing-dot` in a `<header>` → error; `blinking-cursor` in the first 900px → warning. A design review protocol should weight findings by prominence, not treat all surfaces equally.
12. **Steal the FP-fix documentation convention.** Every suppression carries the bug that motivated it, verbatim and searchable. In a multi-agent setting this becomes the audit trail for why a rule was narrowed, which is exactly the thing a later agent will otherwise re-litigate.
13. **Steal the exact WCAG constants and the "measure the worst gradient stop" rule.** `WCAG_LARGE_TEXT_PX = 18*96/72`, `WCAG_LARGE_BOLD_TEXT_PX = 14*96/72`, thresholds 3.0/4.5, gradient backgrounds measured at their worst stop. Cheap, correct, uncontroversial.
14. (Low priority) **The p10-not-min statistic** from the pixel pass. If parley-design ever gets screenshot evidence, report the 10th-percentile contrast with the median alongside — the min is noise, the median lies.

## Do NOT copy

1. **The 5536-line `checks.mjs`.** It is a single-file, browser-and-Node-dual-mode monolith with 90+ exported functions, per-engine adapters for every rule, and jsdom/Chrome/htmlparser2 workarounds interleaved with actual doctrine. A multi-agent protocol should never own this much machinery; the marginal rules it buys are not worth the maintenance surface or the version skew between agents.
2. **The Puppeteer/Chromium dependency, and every rule that requires it.** `heading-rhythm`, `text-occlusion`, `edge-flush-cards`, `first-viewport-column-overflow`, `content-hidden-at-rest`, `text-overflow`, `blinking-cursor`, `script-error`, and the pixel-contrast pass all need real layout. A vendor-neutral protocol whose participants are `claude`/`codex`/`hermes`/`agy`/`kimi` on arbitrary machines **cannot assume a headless browser, cannot assume the same Chrome build, and cannot assume they all render the same page.** Non-determinism across agents would poison consensus. If layout evidence is wanted, make it *one participant's optional attachment*, never a protocol requirement.
3. **The hand-rolled static CSS cascade (`css-cascade.mjs`, 1183 lines).** Its own comments enumerate its holes: `@media`/`@supports` ignored, only `:root` custom properties, pseudo-elements can't join the cascade, "unsupported base selector" silently swallowed, no layout. It is a large, subtly-wrong reimplementation of a browser. Do not rebuild it; declare the rule `unjudgeable` instead.
4. **The `severity: 'advisory'` vs `advisory: true` split.** It is an outright bug surface: 11 rules are *labeled* advisory but still exit 2. If parley-design has an advisory tier, it must be **one** key, and the gate must be derived from it mechanically.
5. **The silent degradation path.** `detectHtml` catching a failed import and falling through to `detectText` means an HTML file can lose ~40 rules with **no diagnostic**. In a consensus protocol, a participant that silently ran a weaker rule set and then signed off is a correctness hazard. Any tier drop must be *reported in the artifact*.
6. **`detect-antipatterns-browser.js` (8250-line prebuilt bundle).** A committed build artifact that must be regenerated (`npm run build:browser`) whenever the source changes — guaranteed drift, and unreviewable in a PR.
7. **The `profiler.mjs` "profile" model.** It measures milliseconds, not design quality. Do not mistake it for a scoring system, and do not build one on it. Also: **do not invent a 0–100 design score.** Impeccable deliberately has none, and a numeric score in a multi-agent protocol invites agents to optimize the number rather than the design.
8. **Brand-domain font exemptions (`BRAND_FONT_DOMAINS`).** Depends on `location.hostname` (browser-only), encodes three specific companies, and is unmaintainable as a general list. If needed, express it as a per-project DESIGN.md allowlist entry instead.
9. **Framework dev-server port sniffing + HTTP fingerprinting** (`FRAMEWORK_CONFIGS`, `isPortListening`) and **`detect-csp.mjs`'s framework-shape classifier.** Both are Next/Nuxt/SvelteKit/Vite-specific and will rot. Vendor-neutrality means not shipping a framework matrix.
10. **The interactive `Continue? [Y/n]` prompt on > 50 files.** Headless agents will hang on it. Any gate in a multi-agent protocol must be non-interactive by default.
11. **Waiver-without-record.** Inline ignores discard the reason at parse time and take effect unilaterally. Copy the syntax, reject the semantics.
12. **The very narrow rules with fragile heuristics** — `theater-slop-phrase` (a single-word regex), `repeating-stripes-gradient` (fires on the mere *presence* of a CSS function anywhere in the file), `shape-assembled-illustration` (8 primitives / 3 fills / 200px), `codex-grid-background` (vendor-named!). They will misfire on legitimate work, and in a multi-agent review a misfiring rule costs a whole deliberation round to argue down. If adopted at all, adopt them as *prompts for human judgment*, never as gates.
