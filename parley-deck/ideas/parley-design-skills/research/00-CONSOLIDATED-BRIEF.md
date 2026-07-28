# Consolidated research brief — designing `parley-design`

**Inputs:** six analyst digests of `Nutlope/hallmark` (v1.1.0, MIT) and `pbakaus/impeccable` (v4.0.3 skill / v3.4.0 detector, Apache-2.0), read in full.
**Target:** a vendor-neutral companion add-on skill for Parley Deck, sibling to `parley-worktrees` and `parley-tracker`, where N headless CLI agents (claude, codex, hermes, agy, kimi) each write their own artifact, cross-review each other, sign a consensus, then implement and code-review.
**Audience:** the humans deciding what `parley-design` is before anyone writes a line of it.

---

## 1. What hallmark is

A single 67 KB doctrinal `SKILL.md` acting as a dispatcher over 24 lazily-loaded reference files, plus 21 macrostructures, 50 component archetypes, 20 themes, 4 genres. Zero scripts, zero tests, zero dependencies — enforcement is entirely textual and self-reported.

Its thesis is structural, not cosmetic: *"Structural sameness is the AI fingerprint, not visual sameness."* Most anti-slop guidance bans purple gradients and still ships hero → 3 features → CTA → footer; hallmark bans the shape.

Its single best idea is the **58-gate slop test**: a numbered, append-only list of yes/no questions where "yes" means fail, loaded strictly *after* generation ("pre-loading slop-test.md costs ~7K tokens for nothing"), genre-scoped with exceptions written inline in each gate, and grown by a hard ratchet — every defect found in one generated page becomes a permanent numbered gate (gates 34–36 came from one test page, 50–57 from one responsiveness pass).

That list is the only design artifact I have seen that turns a taste disagreement into a citation. Everything else hallmark does — the rotation memory, the self-scored six-axis critique, the "58/58 ✓" stamp — is the part a real adversary replaces.

## 2. What impeccable is

A 79-line entry file over 34 reference playbooks (~5,200 lines), ~15 Node scripts, a 60-rule deterministic HTML/CSS detector, 4 subagents with fixed output contracts, editor hooks in 4 harnesses, and a build that emits the same payload into 14 harness directories with CI enforcing `git diff --exit-code` on the generated trees.

Its slop theory is distributional: models cluster on the mode of their training set. Its craft floor is 45 lines of hard numbers (≥4.5:1, 65–75ch, −0.04em, 12–16px radii, one authored motion moment) with exactly one absolute ban (kicker/eyebrow above a heading).

Its single best idea is **process discipline that binds the reviewer, not the builder**: A/B isolation with judgment before mechanical evidence ("detector output is deterministic, but it still anchors judgment"), a mandatory `⚠️ DEGRADED: single-context (<reason>)` banner ("a silent degraded critique is a failed critique"), a fixed 5-section reviewer contract with ≤8 material fixes, and a scored verdict pass — `resolved | partial | unresolved`, two correction rounds maximum, *"the second verdict ends the work whatever it says."*

Its sharpest single sentence: *"if someone could guess your aesthetic from the category alone, **or from category-plus-avoidance**, rework until neither answer is obvious"* — which closes the loophole where anti-slop rules generate their own recognisable slop.

---

## 3. Side-by-side

| Axis | hallmark | impeccable |
|---|---|---|
| **Doctrine depth** | Very deep, very specific, very *marketing-page*. OKLCH 4-layer palette, 2+1 font rule, 1.25 ratio, 4pt nine-step scale, 6 z-levels, 3 easings, 8 states, 4 mandatory viewports, 21 named page shapes, 20 themes, 42,000 six-axis fingerprints. Taste is presented as law. | Deep but narrower and more mode-aware. 45-line craft floor + per-command playbooks + 4 visitor modes (Persuade/Operate/Read/Experience) + platform doctrine (HIG/Material 3). Taste is presented as *defaults you must beat*, with "the brief always wins" stated repeatedly. |
| **Enforcement** | None mechanical. 58 gates run by the same agent that generated the output; six-axis critique self-scored (nothing in the corpus ever scores below 4); verification is literally *"imagine the rendered output."* | Real. 60-rule detector (regex / static-HTML / Puppeteer / pixel engines), exit 0/2, no score. Editor hooks in 4 harnesses (Cursor blocks the write). A separate finish-reviewer subagent that cannot edit. CI drift gates on counts, versions, generated trees, and prose denylists. |
| **State / persistence** | Two things: a greppable CSS provenance stamp (macrostructure · knobs · theme · accent · gate ledger · `studied:`/`context:` flags that change how a later reviewer grades the file) and `.hallmark/log.json` — single-writer, gitignored, trimmed to 20, no author field, no schema version. | Four authority files with explicit *negative* scope lists: `PRODUCT.md` (schema-stamped), `DESIGN.md` (deliberately unstamped, written **after** the build from shipped code), `.impeccable/design.json` sidecar (`schemaVersion: 2`), per-surface briefs. Plus critique snapshots with a trend line, staleness findings typed `auto \| mention \| route`, 7-day renotify throttle in `~/.impeccable/`. |
| **Portability** | `npx skills add nutlope/hallmark`; 3 frontmatter fields; three named harnesses; no build step. Portable *handoff* via always-emitted `tokens.css` + an opt-in ~45-line `design.md` exporting Tailwind `@theme`, DTCG tokens, shadcn vars. | Build-time transform: 14-entry provider table, `<codex>`/`<claude>`/`<gemini>` conditional blocks (unknown tags preserved), six placeholders including `{{ask_instruction}}`, per-provider frontmatter gating, `HARNESSES.md` support matrix, and `degraded/<role>.md` single-sourcing so a harness without subagents runs the role inline **and must disclose it**. |
| **Scope** | Greenfield marketing/landing pages first, components second. Four verbs: build / `audit` / `redesign` / `study`. Explicitly refuses to invent copy, pick a brand, or build logic. | Whole product lifecycle: init → shape → build → 12 refinement verbs → critique/audit gates → finish review → documenter. 23 commands. Web + iOS + Android. Live browser variant mode. |
| **Single- vs multi-agent fit** | Actively hostile in places. Diversification is defined as "differ from *my own* last run"; `log.json` is a single-writer file; the stamp is a project singleton; three blocking questions are asked *always*. Its own ROADMAP concedes the variety rule is *wrong* for a multi-page product. | Closer, and self-aware about it. A/B isolation, degradation banners, anti-anchoring, fixed output contracts and bounded rounds are all directly reusable — but they are expressed as *parent spawns child*, which Parley does not have. Everything else assumes one writer per file, one browser, one human answering mid-flow. |
| **Maintenance burden** | Low code, high drift. No tests means nothing is verified: README says 57 gates and the file says 58; `tokens.css` says "Twenty-four themes" and defines 20; `recipes.md` computes 16 themes; the `refine` verb is documented and does not exist; the test corpus still names six deleted themes. | High but *guarded*. 15+ harness mirrors of the same content, a committed 8,250-line browser bundle, a 5,536-line `checks.mjs`, a hand-rolled CSS cascade — all held together by CI drift gates. Also carries its own live spec/tooling drift (`document.md` defines 8 DESIGN.md sections; the parser models 6). |

---

## 4. The transferable core

Ranked by expected value for Parley. Each entry: **what it is** / **why it works** / **inside a parley protocol**.

### 1. A numbered, append-only rule registry as the shared objective function

**What.** hallmark's 58 gates (yes/no, "yes" = fail, genre exceptions inline in the gate text) fused with impeccable's `antipatterns.mjs` — 60 rules as pure data, `{id, category, scopes, severity, name, description}`, **zero detection logic**, each description written as *the tell / why it reads as machine-made / the corrective move*.

**Why it works.** It converts "the hero feels off" into "gate 44(b) fails at 1280×800" and "`icon-tile-stack` at Hero.tsx:41". Numbers are citable, contestable, and diffable. Rules are append-only: a rule id never changes meaning (hallmark's `38a` is the precedent for inserting without renumbering).

**In parley.** One canonical, versioned file — `addons/parley-design/references/rules.md` — read once in round 0 by every participant. Every cross-review finding must cite a rule id **or** be explicitly labelled a judgment call with its own justification. This is the single highest-leverage import, and it is the artifact that makes design consensus arithmetic instead of rhetoric. Rules registry ≠ detector: ship the doctrine with no runtime, exactly as impeccable's registry does.

### 2. The `slop` vs `quality` split, with different burdens of proof

**What.** impeccable separates "this looks machine-made" (32 rules) from "this is objectively broken" (28 rules) with explicit comment banners in the source.

**Why it works.** They are different arguments. `low-contrast` and `text-occlusion` are defects; `cream-palette` and `italic-serif-display` are taste with a strong prior. Collapsing them means either taste findings block ships, or defects get argued down.

**In parley.** This is the decision that shapes the whole protocol: **a `quality` finding lets a single agent BLOCK; a `slop` finding requires quorum.** It maps straight onto Parley's existing dispositions and `strict_gate`, and it is the guard against a five-veto ratchet toward the category standard (see §5).

### 3. Evidence tiers, `unjudgeable`, and a mandatory degradation banner

**What.** impeccable's consistent "skip, don't guess" doctrine plus a literal third state (`fontSizeStepStatus` → `on-ramp | off-ramp | unjudgeable`), its four engines with strictly different sight (regex sees text and line numbers; static HTML sees no layout; only the browser sees layout; pixels see only 12 candidates), and the banner rule: `⚠️ DEGRADED: single-context (<reason>)` — *"a silent degraded critique is a failed critique."* Also: *"'Unavailable' means exactly one thing… It does not mean inconvenient."*

**Why it works.** It makes the *limits* of a verdict part of the verdict. hallmark's opposite choice — gates that say "imagine the rendered output" — manufactures false passes, and its own ROADMAP admits there is no rendering loop at all.

**In parley.** Every rule declares the evidence tier it needs (`text | source | rendered | layout | pixel`). A participant with no browser writes `unjudgeable: rendered` and is *compliant*, not silent. Any artifact produced with fewer participants than the roster, or without a capability the rule needed, leads with a degradation banner. Parley's roster is heterogeneous by construction (agy, hermes, kimi capabilities vary run to run) — this is the mechanism that keeps that honest instead of hidden.

### 4. Plan-before-code as the fixed-schema round-01 artifact

**What.** hallmark's Step-5 Preview (six bullets: macrostructure · theme · enrichment · sections · motion · diversification, *"Markdown bullets, not ASCII boxes — they render reliably across every chat client and terminal"*, and *"the preview is the durable summary; it's wrong to ship if it lies"*) plus impeccable's direction contract: five blocks, ≤150 words, in the emitted artifact's own head comment — `THESIS / OWN-WORLD / STORY / FIRST VIEWPORT / FORM` + a literal `FINISH:` line. Falsifiability rule: *"If a block reads like a mood, the direction is not decided yet."*

**Why it works.** Five prose essays cannot be compared. Five instances of the same six-field block can be diffed line for line. The `FINISH` line is a self-carried exit condition that survives a long multi-round run.

**In parley.** This is the round-01 artifact schema. Standardise the field set so consensus is a field-by-field comparison, put the ratified contract in the *emitted file* (not only in `FINAL.md`) so every implementer re-reads it on every edit, and make "if a block reads like a mood" a review check. `FINISH` maps to the phase gate: *"unreviewed and undocumented is unfinished."*

### 5. Named Rules as the grain of the ratified design system

**What.** impeccable's DESIGN.md format: `**The [Name] Rule.** [one forceful sentence]`, 1–3 per section — *The One Voice Rule*, *The Hairline First Rule*, *The Weight-Inversion Rule*. *"Much stickier for AI consumers than bullet lists."*

**Why it works.** A named rule can be cited, contested, versioned, and violated *by name*. A bullet list can only be vaguely gestured at. This is the difference between a reviewable artifact and mush.

**In parley.** Every durable decision in the ratified design artifact is a Named Rule. Cross-review findings cite them (`violates The Two-Face Rule`). Dissent at consensus attaches to a rule name, not to a paragraph. Pair with hallmark's `design.md` governance verbatim in shape: **"What pages MUST share" / "What pages MAY differ on"** — that section pair *is* the parallelisation contract for Phase-5 implementers touching different screens, and it inverts hallmark's own variety rule to consistency, which is the correct default for a product.

### 6. The deviation-classification vocabulary, and "an uncited deviation is a defect"

**What.** impeccable's finish reviewer classifies every salient element as `match | acceptable adaptation | missing | contradicted | added without approval`, with two mandatory matrix rows (TYPE and MATERIAL), and the binding clause: *"An adaptation counts as intentional only when it **cites** the user answer, surface brief, accessibility need, or product truth that forced it."*

**Why it works.** It forces evidence citation and it makes "the implementation drifted from the design" a mechanical finding rather than a vibe. It also names the two things reviewers always skip: whether the type is the right *character* and whether the *medium* was honoured.

**In parley.** This is the schema for the Phase-6/7 design review — the design-domain equivalent of "does the code match FINAL.md". hallmark supplies the same idea from the other end with its `stamp lies` check: *"if the stamp says Bento Grid but the page is a centred single-column hero, flag it as a critical structural finding."* Both are drift guards, and Parley already runs drift guards (`TestEmbeddedDefaultMatchesLiveDeck`).

### 7. Bounded verification with a scored verdict pass

**What.** ≤2 inspection rounds, ≤8 material fixes, ≤2 correction rounds, ≤3 regressions, fixes scored `resolved | partial | unresolved`, *"a fix answered mechanically, positions moved but the quality the finding named still absent, is **partial at best**"*, *"the reviewer's findings are the only list you work from, never your own re-opened hunt"*, and the closer: **"the second verdict ends the work whatever it says."** Plus: *"presenting mechanical confirmation as artistic success is how a failed build gets announced as a finished one."*

**Why it works.** It replaces an unbounded polish loop with a bounded one that ships with open items *stated*. hallmark has the same instinct in one line — *"two passes is normal, three means the brief is wrong, not the design."*

**In parley.** Parley has this disease documented: Tier 4 of the loop-engineering run took four fix-up cycles. Import the ceilings as protocol law, and import the verdict pass wholesale — it is exactly the missing "did the fix reach the quality the finding named" step that a re-run of `RunChecks` cannot supply.

### 8. Content-hash-pinned approvals

**What.** impeccable's catalog validator: `conceptContentHash = sha256(...).slice(0,12)`, carried in every review record, and *"an approval cannot silently survive a content edit: the validator rejects any review whose hash no longer matches the concept it points at."*

**Why it works.** A signature over content that has since moved is worse than no signature — it is a false attestation that a later agent will trust.

**In parley.** Highest-leverage import for the *consensus* half of the protocol, and it generalises beyond design: every signature under `parley-deck/ideas/<slug>/` should carry a 12-hex hash of the exact reviewed content, and a validator (`parley-tracker`'s gap-scan is the shape precedent) should reject a signed FINAL whose body changed after signing. Parley already has the disease this cures.

### 9. Per-model bias corrections, for the actual roster

**What.** impeccable ships model-attributed doctrine: `<claude>` *"your cream/serif/lamplight prior — treat that first palette as already spent"*; `<codex>` ghost card, sketchy SVG, `repeating-linear-gradient` stripes, grid backgrounds; `<gemini>` never animate an image on hover. Five of its 60 detector rules are model-attributed (`gpt-thin-border-wide-shadow`, `codex-grid-background`, …).

**Why it works.** The prior is per-model and measurable. Generic anti-slop advice is weaker than "here is *your* specific reflex, spend it before you start."

**In parley.** No analogue exists today and Parley is uniquely positioned to build one: it runs the same brief through five models on every idea, so it can *measure* each participant's priors empirically instead of guessing them. Ship a per-participant "known prior" block addressed at runtime (`if you are the codex participant…`) — **not** at build time, because one COOPERATION.md is read by all five at once. This is the most novel thing on this list.

### 10. The defect → permanent rule ratchet

**What.** hallmark's core loop: every regression found in a generated page becomes a numbered gate forever. Gates 34–36 from one page; 50–57 from one responsiveness pass; gate 52 names the exact failure ("italic Anton title overlapping '02 / EXAMPLES'").

**Why it works.** It compounds. A finding that stays a fix teaches nothing; a finding that becomes a rule protects every future run.

**In parley.** Wire it to §13 (RHO): a design defect surviving to Phase-7 must land as a new rule id in the registry, not just as a patch. The advisory-only, quorum-replaces-self-preference posture of §13 is exactly right for taste rules.

### 11. A three-axis fingerprint — used to measure agreement, not force difference

**What.** hallmark reduces any theme to `paper-band × display-style × accent-hue`, three comparable values computable from the stamp, and requires consecutive runs to differ on at least one.

**Why it works.** It is the cheapest possible "are these two designs actually different?" test, and it kills the colour-swap illusion (*"score by structural distance, not visual distance"*).

**In parley — inverted.** Compute the axes for each participant's proposal. **Convergence on all three across five independent models is not a pass, it is a slop alarm** — it means the roster hit a shared training attractor. Divergence is a real design fork, stated as axis deltas instead of vibes, and it gives the consensus round a coordinate system. Pair with impeccable's calibration test: *"if someone could guess your aesthetic from the category alone, or from category-plus-avoidance, rework."*

### 12. Anchoring control as protocol ordering

**What.** Two rules. (a) *"Assessment A must finish before detector findings enter the parent synthesis context. Detector output is deterministic, but it still anchors judgment."* (b) The reviewer's guard: *"inventory the comp's salient elements in your own words **before** reading the direction contract or any builder-authored summary: the contract is the builder's abstraction of the comp, and a review anchored on it inherits whatever that abstraction dropped."*

**Why it works.** Deterministic evidence and author summaries both collapse independent judgment into confirmation.

**In parley.** Parley gets isolation for free — separate processes, separate files — but it must *say so* and forbid the failure: no participant reads a peer's artifact before writing its own; mechanical check output enters only after each participant's judgment is written; a reviewer reads the artifact before the author's summary. Express as **round ordering and artifact visibility**, never as spawn instructions.

### 13. A ≤50-line craft floor with a stated override order

**What.** impeccable's `craft-floor.md`: two lists only. `Verify` (8 checks *on the built result, not an intention*, each with a number: ≥4.5:1, 65–75ch, −0.04em, 12–16px radii, one authored motion moment) and `Refuse` (category defaults, framed as *"reaching for one when the axis is free means you were not deciding"*) with exactly **one** absolute ban. Loaded just-in-time before editing, never for planning. Override order stated in the file: *"A pinned brief or the committed visual world overrides anything here; your own habit does not."*

**Why it works.** Short enough that five heterogeneous models parse it identically. The framing ("defaults, not bans") is what keeps it from becoming a house aesthetic, and the one hard ban is what keeps it from being toothless.

**In parley.** `DESIGN-FLOOR.md` in the deck, cited by rule name in cross-review, with hallmark's numbers folded in where they are genuinely invariant (contrast thresholds, 8 interaction states, `focus rings never animate`, no `transition: all`, no invented metrics) and its house taste left out (see §7). State the precedence chain explicitly — **rules > ratified deck design system > brief > parity with existing code > model habit** — and steal hallmark's gate-54 formulation for the top rules: *"NOT bypassable by 'preserve structural parity' / 'mirror this reference' / 'match the prior build' instructions… The rules win over parity"*, plus its "binds on content shape, not class-name allowlist" clause.

### 14. Honesty machinery: labelled holes, labelled inferences, labelled gaps

**What.** Three separate devices. (a) hallmark gate 46: an invented metric auto-fails; the fixes are a labelled `—` block, a question to the user, or removing the proof slot — *"the number-shaped hole is honest; the fabricated number is slop"*, *"the model is not allowed to invent specificity."* (b) impeccable's truth doctrine: *"Truth binds claims, not demonstrations… label it synthetic… What stays uninventable are commercial and factual claims… Refusing a bold direction because its demonstration data does not exist yet is the timidity reflex wearing honesty's clothes."* (c) The reviewer's budget rule: *"a review built from what you saw beats a perfect review that never arrives… by roughly the tenth turn stop reading and write. **Name whatever went unread**."*

**Why it works.** In a multi-agent run fabrication compounds: one agent invents "10× faster", the next three treat it as given. And a headless agent that runs out of budget silently produces a truncated review that reads exactly like a complete one.

**In parley.** All three become artifact requirements: no invented facts (leave a labelled hole), every inferred fact labelled inline with the disclosure leading the artifact rather than trailing it, and every review declaring what it did not read. impeccable's inversion of the unattended-user probe is the right frame — Parley agents are unattended by construction, so labelling is the only defence.

### 15. Reason-carrying waivers with a counter-signature

**What.** impeccable's narrowest-exception ladder: exact value → value-scoped-to-file → whole file (*"it silences every rule for that file permanently, including rules that have not been written yet"*) → whole rule. Every stored ignore carries `createdAt` and `reason`. *"The hook itself never writes ignore config."* One deliberate anti-laundering rule: `undersized-ui-text` is design-system-blind — *"adding 8px to the ramp launders the token but not the legibility problem, and that is exactly the escape hatch this rule closes."*

**Why it works.** Suppression is where every rule system dies. Impeccable's own bug is instructive: inline ignores parse the reason and **discard it**.

**In parley.** Copy the syntax, reject the semantics: a design waiver names rule, value, scope, and reason, lives in one reviewable place, and **requires a second participant's acknowledgement recorded in the deliberation**. Keep at least one rule class deliberately design-system-blind, because a multi-agent design system will absolutely be widened by an implementer to legalise its own output.

### 16. The effort ladder that defaults to nothing

**What.** hallmark's enrichment hierarchy: typography-only → Tier A pure CSS → Tier B hand-built SVG → Tier C generated → Tier D library → **Tier E Lottie, last resort** — *"reaching for Lottie when CSS would have built it is the new tell"*, *"when in doubt, no images"*, *"cut motion before adding it; most pages have too much, not too little."*

**Why it works.** It is trivially reviewable ("you shipped Tier D; justify why Tier A failed") and it is the best structural defence against decorative slop, because the default is nothing.

**In parley.** A generic complexity ladder with a stated default of "nothing" and a required one-line justification for every step up. Reviewable by any agent without taste judgment.

### 17. Write the durable design system *after* the build, from shipped code

**What.** *"A rulebook written before the build gets defended against reality instead of describing it, and it hands the design-system detector an unstable target that buries the build in noise."* Written by a dedicated documenter whose ground truth is the shipped artifact, with two named failure modes: *a prohibition that bans a device the world itself uses natively*, and *a value recorded to legitimize a defect*.

**Why it works.** It stops a speculative spec from becoming a thing agents defend.

**In parley.** This creates a real tension with item 5 (the design system as `FINAL.md` binding implementers) — see §6, question 9. The honest resolution is probably two artifacts with different names and different authority, not one.

---

## 5. What multi-agent adds that neither has — and where it doesn't

### Where independent adversarial agents genuinely beat one agent, for design specifically

**1. Convergence detection is the strongest slop detector anyone has built, and only a roster can build it.** impeccable measured its own disease — *"30/35 identical concepts across 16 prompt framings; the model cannot roll its own dice"* — and solved it with an external RNG. Parley does not need dice: it has five models trained on overlapping but non-identical distributions. If four of five independently reach for cream ground + high-contrast serif + terracotta accent, that is not consensus, that is the *distribution* talking, and it is measurable in a way no single agent can measure about itself. This inverts hallmark's rotation rule from a chore into a signal, and it is the one capability that justifies the whole add-on existing.

**2. Self-scoring becomes peer-scoring for free.** hallmark's corpus never scores below 4 on any of its six axes, and every page self-reports `38/38 ✓` or `58/58 ✓`. The rubric is fine; the scorer is the wrong party. Reassigning the six axes (Philosophy · Hierarchy · Execution · Specificity · Restraint · Variety) to peers and publishing the matrix turns a vanity stamp into a disagreement agenda. Divergent scores on one axis *are* the deliberation topic.

**3. Anchoring isolation, which impeccable pays subagent-spawn complexity to get, is structural in Parley.** Separate processes, separate files, no shared context. A/B isolation, "judgment before evidence", and "read the artifact before the author's summary" cost Parley nothing but a written rule.

**4. Capability federation instead of capability requirement.** impeccable's most valuable rules (occlusion, heading rhythm, hidden-at-rest, pixel contrast) need Chrome and a running app; it therefore makes the browser near-mandatory and the degraded path second-class. Parley can invert this: *one* participant with a browser attaches rendered evidence, and the other four consume it as evidence rather than pretending to see. Neither source can do that.

**5. False-positive triage.** A detector rule that misfires costs impeccable a wasted user interaction; in a five-agent review it costs a whole round — *unless* the protocol lets three participants contest a finding on the record. Contested-finding handling is something only a multi-agent protocol can even express, and it is the antidote to importing a rule corpus that was tuned on someone else's bug reports.

**6. Real alternatives, authored.** impeccable's whole `concept-seed` apparatus exists to stop one model shipping its argmax. Five models each committing fully to their own direction produces genuinely different designs with genuinely different rationales — a better version of hallmark's unbuilt `hallmark variant` roadmap item (*"the biggest cause of 'AI feel' is users accepting the first output because they didn't know it could be different"*).

### Where multi-agent does **not** help — state this plainly in the skill

**1. Taste by committee produces mush.** Design is not an averaging domain. The mean of five directions is, definitionally, closer to the training mode than any single committed direction — five agents converging on a compromise reproduce exactly the slop the skill exists to fight. **The merge operator for design must be *choose*, not *average*.** Consensus selects one authored direction with a named author and records dissent verbatim; it never blends palettes, never splits the difference on a structure, never adopts "the best parts of each". impeccable already knows this at single-agent scale: *"Refinement preserves; redesign replaces. Never split the difference into polish on the discarded look."*

**2. Quorum ratchets toward safe.** Five reviewers each holding a veto will converge on the category standard, because the category standard is the option nobody objects to. This is why item 2 of §4 matters: give `quality` findings a single-agent veto and force `slop` findings through quorum, or the protocol will systematically discard the bold direction. impeccable's counter-device — the "standing exit" (the category standard played straight, *"the user's door, never yours: never recommend it"*) — belongs here.

**3. Five blind reviewers are five times as confident and exactly as wrong.** None of the roster renders anything by default. Multi-agent multiplies *confidence* without multiplying *evidence*, and agreement among five agents that never saw the page is the single most dangerous artifact this add-on could produce. Every §4 honesty mechanism (evidence tiers, `unjudgeable`, degradation banners) exists to stop this specific failure.

**4. Cost and latency are indefensible for the fast path.** A landing page can be judged by a human in three seconds by looking at it. Running five agents and a consensus round to decide whether a hero is centred is absurd. `parley-design` must respect §4.0's `fast` track and, for small surfaces, produce a *screenshot or a preview block* rather than a deliberation.

**5. Consistency work wants fewer voices, not more.** hallmark's ROADMAP concedes it: *"The structural-variety rule is correct for variety, wrong for brand consistency inside a multi-page product."* When the job is "make these six screens feel like one product", N independent authors add pure entropy and the right answer is one ratified system and N obedient implementers. Multi-agent adds value at *direction* time and at *review* time; at *coherence* time it is a liability.

**6. Some rules are ungovernable by discussion.** Contrast ratios, focus states, token drift, banned fonts — these want a tool, not a debate. Anything mechanically checkable should be checked once and reported, not argued five times.

---

## 6. Hard design questions the humans must answer

Each is a genuine fork. The trade-off is named; the choice is not mine.

**Q1 — Surface scope: web-only, or surface-agnostic?**
Roughly 60% of hallmark's doctrine and ~100% of impeccable's detector are CSS-shaped (`overflow-x: clip`, `minmax(0,1fr)`, Tailwind class regexes, `oklch()`). *Web-only* inherits hundreds of hard numbers immediately and excludes TUIs, CLI output, docs, and slide decks — including Parley's own TUI, which is the most likely first customer. *Surface-agnostic* keeps every hard number in a clearly-labelled web annex and risks the generic layer degenerating into vibes. Pick before writing a single rule; this decision is not reversible cheaply.

**Q2 — Deliverable: decisions, or code?**
Does `parley-design` own the UI implementation in Phase 5, or does it produce only the ratified design artifact that ordinary implementers obey? *Owning code* gives it teeth and a checkable output; *decisions only* keeps it a clean companion add-on (the `parley-worktrees` / `parley-tracker` precedent: never change canonical artifact ownership) but leaves enforcement to whoever writes the CSS.

**Q3 — Rule authority: binding law, advisory doctrine, or split by category?**
*Binding* means one agent citing rule 17 can BLOCK a consensus — powerful, and a single false positive costs a full round. *Advisory* means rules inform but never stop, which is how hallmark works and why nothing in its corpus ever fails. *Split* (my recommendation: `quality` binding, `slop` quorum) needs the categorisation to be right on day one, because re-categorising a rule later invalidates past reviews.

**Q4 — Convergence semantics: PASS or ALARM?**
When five independent agents propose the same direction, is that (a) strong evidence the direction is correct, or (b) evidence of a shared training attractor? Both readings are defensible and they demand **opposite** protocol behaviour: (a) fast-tracks to consensus; (b) triggers a mandatory rework round. A hybrid ("convergence passes only if the converged direction survives a check against the ban list and the category-plus-avoidance test") is possible but must be specified precisely or it becomes a coin flip.

**Q5 — Consensus operator: select, or synthesise?**
*Select* one participant's direction wholesale, name the author, record dissent — coherent, but discards four agents' good ideas and makes 80% of the compute a review function. *Synthesise* a merged direction — uses everything, and risks the mush failure in §5. There is a third option worth considering: **select the structure, synthesise the rule list** (one author owns the shape; the Named Rules can be merged because they are individually falsifiable).

**Q6 — Executable checks, or prose only?**
hallmark is prose-only and nothing in it is verified. impeccable is script-heavy and needs Node ≥22.12, htmlparser2, Puppeteer, Chromium. `parley-tracker` already sets a precedent (a tool-enforced gap-scan validator), so the deck can host a checker — but what runtime may `parley-design` assume across a vendor-neutral roster, and who runs it (each agent? the driver, once, like `RunChecks`)? A ~19-rule zero-runtime static tier (text-only rules: banned fonts, single-font, flat hierarchy ratio < 2.0, monotonous spacing, gradient text, buzzword list, em-dash overuse) is available with no dependencies at all and is probably the right v1.

**Q7 — Rendered evidence: forbidden, optional-with-declaration, or required?**
*Forbidden* is honest and weak (no rule about layout can ever fire). *Required* makes agy/hermes/kimi structurally non-compliant. *Optional with declared capability* is obviously right in principle and requires the full evidence-tier + `unjudgeable` machinery from §4.3 to avoid agents claiming verification they did not perform. There is no cheap version of this.

**Q8 — Where does taste live: shipped canon, ratified per deck, or nowhere?**
*Shipped canon* (hallmark's 20 themes, impeccable's Neo Kinpaku) makes the skill immediately useful and makes every project that uses it look the same — which is the disease. *Ratified per deck* (humans sign a DESIGN.md once, agents obey it) is correct and means the skill does nothing useful on day one of a greenfield project. *Anti-slop invariants only* is the purist position and leaves agents to fall back on their training priors in every free axis — which is precisely what both source projects were built to prevent.

**Q9 — When is the design system written: before implementation, or after?**
impeccable is emphatic: after, from shipped code (*"a rulebook written before the build gets defended against reality"*). Parley's protocol is emphatic the other way: `FINAL.md` is ratified *before* Phase 5 and binds implementers. These are directly incompatible if it is one artifact. Either accept a pre-build binding contract and lose the "describes reality" property, or ship **two** artifacts (a ratified direction contract before; a documented design system after, authored from the merged diff) and pay for two.

**Q10 — Dice: endogenous divergence, or a seeded assignment?**
Is cross-model divergence sufficient randomisation, or does each participant additionally need a deterministic seeded assignment forcing it to develop its own rank-k candidate rather than rank-1? *Endogenous* is free and may be an illusion (all five are trained on the same web). *Seeded* is impeccable's proven answer but adds a mechanism, a seed key to record, and a reviewer check that the roll happened — and Parley must derive the seed locally from the run id, never from a hosted API.

**Q11 — Add-on shape: teaching skill, or protocol amendment?**
`parley-worktrees` and `parley-tracker` are opt-in add-ons that never change canonical artifact ownership. Does `parley-design` follow that (teaches a discipline, owns no phase) or does it amend COOPERATION.md with a binding design phase and its own gate? The former is safe and easily ignored; the latter has teeth and touches the core protocol the other add-ons deliberately do not.

**Q12 — Doctrine size and load model.**
hallmark: 24 references + 5 subtrees (~400 KB) with per-file lazy loading and a token argument. impeccable: 34 playbooks, exactly one loaded per request. Both economies *invert* under N agents — five agents lazily loading the same files costs 5× and guarantees divergent reads. The alternative is one pre-digested doctrine artifact everyone reads once. Name the ceiling now (my read: **≤4 files** — floor, rules registry, artifact schemas, web annex), because file count only ever grows.

---

## 7. Anti-goals — what `parley-design` must refuse to become

1. **A house aesthetic.** Not 20 named themes, not Neo Kinpaku, not warm-paper editorial with hairline rules. Both source projects ship a taste and call it a standard; hallmark's own `contract.md` claims *"It will not enforce a specific style"* while its gate list bans italic headings and `system-ui` outright. Ship **anti-slop invariants** (contrast, states, motion timing, honest copy, no re-drawn chrome) separately from **house preferences**, and let the deck's own ratified artifact own the latter, overridable on record.

2. **A CSS linter with a protocol wrapper.** No 5,536-line `checks.mjs`, no hand-rolled CSS cascade, no Puppeteer requirement, no framework port-sniffing matrix, no committed 8,250-line browser bundle. Non-deterministic layout evidence across five agents on five machines poisons consensus.

3. **A second state machine.** No 23-command menu, no four-verb mode surface. Parley already has phases, dispositions, consults, tracks, and a driver. `audit` is the review phase; `redesign` is an idea whose input is existing code; `study` is input preparation. Re-encoding them as skill verbs creates two competing state machines and imports the `/`-menu-pollution problem impeccable's own repo docs warn about.

4. **A blocking-question machine.** hallmark asks three questions **always**, with *"no length threshold below which asking is skipped"*; impeccable's "STOP and use the structured question tool" appears in ≥8 files. Five headless agents × three questions is fifteen blocking prompts and zero answers. Keep the *content* (audience · use · tone-as-an-extreme; *"'clean and modern' is not a tone"*), ask the human **once, before fan-out**, and convert every remaining ask into a labelled assumption or a Parley consult.

5. **A single-writer mutable state tree.** No `.hallmark/log.json` (gitignored, trimmed to 20, no author, no schema, no conflict handling — five agents will clobber it), no `.impeccable/` runtime tree, no `~/.parley-design/` throttle cache, no mtime-based staleness. Parley already has `parley-deck/ideas/<slug>/` as canonical state; design state is per-agent artifacts plus one reconciliation step, keyed by `(idea, agent, round)` and pinned by content hash.

6. **Self-scoring theatre.** No `58/58 ✓` written by the author, no `P5 H5 E5 S5 R5 V5` self-stamp, and above all no gate whose verification method is *"imagine the rendered output."* A gate that cannot be executed is worse than no gate, because it manufactures passes. Keep the tally *format*; move authorship to a peer or a tool.

7. **A forced-diversity generator.** Do not import "differ from the last run on at least one axis". With N agents on one idea, mandated difference produces incoherence and convergence is data. Import the axes as a *similarity metric*; decide per context whether similarity is good.

8. **A numeric design score.** impeccable deliberately has none, and it is right: a 0–100 invites agents to optimise the number. Calibrated rubrics with honest denominators (*"never print /40 over a partial set"*, *"most real interfaces score 20–32"*) are fine; a headline grade is not.

9. **File sprawl.** 24 references or 34 playbooks is a maintenance surface that only survives because exactly one file is read per request by one agent. Under N agents it is both expensive and a guaranteed drift surface — hallmark already fights drift between `structure.md`, `layout-and-space.md` and gate 54.

10. **Hand-maintained facts.** hallmark says 57 gates in the README and 58 in the file; "twenty-four themes" over 20 definitions; `2,688 fingerprints` computed from a dead catalog; a documented `refine` verb that does not exist; a test corpus naming six deleted themes. Anything stated twice must be generated or drift-guarded. Parley already owns the pattern (`TestEmbeddedDefaultMatchesLiveDeck`).

11. **Vendor coupling.** No Nanobanana, no shadcn-as-fix-#1, no hosted roll API, no telemetry ping, no Together-AI or Vercel branding, no `BRAND_FONT_DOMAINS` ("Geist is slop unless you are Vercel"). Local-file-first, reproducible from the run id.

12. **A committee that averages taste.** Restated as an anti-goal because it is the most likely failure: the consensus operator must never blend directions. If the protocol cannot select, it must escalate to the human — not compromise.

13. **A silent degrader.** No participant may sign a design verdict without declaring which evidence tiers it could not reach, which rules it marked `unjudgeable`, and what it did not read. impeccable's own `detectHtml` silently falls back to regex and loses ~40 of 60 rules with no diagnostic; that exact bug, in a consensus protocol, is a signed lie.

14. **An accessibility deferral with no backstop.** impeccable's finding is real (*"models over-cautious themselves into safe, underdesigned output when reminded about accessibility at design time"*) and its mitigation only works because one agent reliably reaches its own audit phase. Adopt the insight — keep the safety checklist out of the generative round — but bind it to a hard gate: **no consensus signature without the accessibility rules having been evaluated by someone.**
