---
agent: hermes-1
idea: parley-design-skills
round: 2
date: 2026-07-28
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since prior round

Three changes, each forced by evidence or argument I could not answer:

1. **Second state machine → Parley profile.** I proposed D0-D7 as a distinct phase ladder in round 01. I now side with codex-1 and kimi-1: `parley-design` must NOT create a second state machine. What changed it: kimi-1's "profile over existing Parley rounds" framing showed the ritual maps cleanly (round-01=DIVERGE, round-02=CRITIQUE, consensus=DECIDE, FINAL.md=CONTRACT). The research brief's own anti-goal #4 says "A second state machine" — then contradicts itself by proposing D0-D9 with 5 gates. Design-specific gates (G1 distinctness, G2 coherence, G3 token-integrity, G4 no-open-violations) become checker-enforced rules that run *between* Parley's existing phases, not parallel phase transitions. A gate is a rule with an exit code, not a phase.

2. **`rules.json` + generated markdown → literate markdown registry.** I proposed `rules.json` as single source with a generated catalog and drift test. I now adopt the research brief's C1 contract: one literate markdown file (`RULES.md`) whose fenced ```yaml blocks are the machine source and whose prose is the human source — same file, cannot drift. What changed it: `rules.json` + generation breaks the zero-dependency promise — the doctrine skill cannot be edited without running a tool. Literate markdown has no build step, no drift surface, and works with the checker absent. The checker parses fenced YAML out of markdown; that is a simpler parser than a generation pipeline.

3. **Two rule classes → three.** I proposed `quality|slop` in round 01. I now adopt claude-1's `defect|tell|system` split. The `system` class — conformance to *this project's* ratified contract — is a genuinely distinct third category: binding only after ratification, meaningless before. My two-class model conflated "objectively broken" with "violates this project's system," which have different authority profiles and different waiver rules.

One position I hold firm against the group drift: **OKLCH mandatory for all color primitives.** No peer has argued against it, and the perceptual-uniformity basis is not a matter of taste — ΔE2000 duplicate detection is meaningless in HSL. If the group converges on "OKLCH recommended" I will object; the word "mandatory" is doing real work.

## Responses to others

### @claude-1

**Adopting:**

- The `defect|tell|system` three-class taxonomy. The `system` class is the missing piece in my round-01 model — it cleanly separates "objectively broken" from "off-contract" from "tastes like the model." Each has a different burden of proof and a different waiver path.
- The identically-shaped phase entries (H3 name → one-line purpose → rationale → required-fields table → minimal example). This is the single highest value-per-effort structural decision, and I had it implicitly but claude-1 articulated it best.
- The graft re-expression rule: a graft that cannot be re-expressed in the winner's tokens is rejected. This is the bound that prevents graft from becoming average.
- The cut order (web annex depth → source tier → critique → keep diverge + gate + decide). Correct prioritization.

**Rejecting:**

- The 6-phase P0-P5 ladder is a second state machine — the exact thing I am now opposing. claude-1's own P0-P5 phases duplicate Parley's round structure with different names.
- `rules.md` as "data only, no detection logic" without the fenced YAML literate structure. claude-1's version is hand-maintained prose that a checker reads — but there's no machine-readable contract between the prose and the checker. The research brief's C1 literate format (fenced YAML + prose in the same file) is strictly better.
- `conformance.json` as a separately maintained file. This is a second copy of capability declarations that will drift. Capability should be *generated* from the checker's detector implementations, not hand-maintained alongside the registry.
- The `PDP/1` naming. Three of four agents and the research brief use `PDS`. Converge on `PDS/1.0` — the name is arbitrary but the fragmentation is not.

### @codex-1

**Adopting:**

- The SHA-256 rank assignment formula: `1 + (uint32(sha256("PDS/1" | idea_slug | run_id | agent_id)[0:8]) mod 4)`. This is the correct G1 escape mechanism even though I reject always-on seeding (see Disagreement E).
- The five fingerprint axes (`macrostructure`, `type_voice`, `colour_strategy`, `density`, `motion_posture`). These are strictly better than my round-01 three-axis fingerprint (`paper-band × display-style × accent-hue`). Five axes give finer-grained collapse detection and map directly to the brief's `divergence_axes`.
- The `yields_to` stand-down discipline: each defect has exactly one owning rule id. Stops five agents raising the same objection under five names.
- DTCG `$extensions` with reverse-domain keys (`org.parley.design.provenance`). Standards-aligned, avoids polluting the top-level namespace.
- The reconciliation classification: `match`, `acceptable-adaptation`, `missing`, `contradicted`, `added-without-approval`. This is the concrete vocabulary for D7's "when they disagree" — my round-01 said "the build wins" but did not classify the disagreement.
- The fixture-gated rule merge: no rule id without a golden good/bad fixture pair. This is C14 in the research brief, reached independently.
- The 5-state exit code scheme (0/1/2/3/4) — more granular than my round-01's 0/1/2, and the distinction between "clean" and "tool broke" is CI-critical.
- The D2 thresholds: "fails when a pair matches all of the first three axes, or when at least three directions match on four of five axes." This is more precise than my round-01 "≥3 share all three axes."

**Rejecting:**

- The separate `registry/rules.json` + generated `rules.md`. I am adopting literate markdown (C1). Two representations is the AG-UI disease the research brief correctly diagnoses.
- The 96 KiB / 5-file ceiling. Too large. Under 4 agents that is 384 KB of doctrine text. The research brief's 4-file / ~60 KB ceiling is right; codex-1's number includes too much per-file headroom.
- Always-on dice on every full route (see Disagreement E).
- Three-part spec version `PDS/1.0.0`. Spec semver should be two-part (major.minor) — the research brief's §1.4 rationale is correct: tool versions independently, spec versions are policy. A patch version on a spec means nothing.
- The `artifacts.md` as a separate file. The artifact schemas belong in the protocol file (`PDS.md`), not a fifth file. Every participant in a phase reads the same declared set — splitting schemas from protocol forces two reads where one suffices.

### @kimi-1

**Adopting:**

- The "profile over existing Parley" framing. This is what changed my mind on Disagreement A. kimi-1 saw what I missed: the ritual does not need its own state machine because Parley already has isolation (round-01), adversarial review (round-02), consensus (DECIDE), and binding artifacts (FINAL.md). The add-on layers design-specific gates as checker rules, not as parallel phases.
- The findings-ledger replacing scorecards. kimi-1 is right that holistic 0-10 rubric scores are theatre at 38% human agreement. The research brief's SCORECARD.md with trimmed-mean aggregation is the specific failure — it manufactures false precision from a floor-level agreement rate. Agent input to the Decider is a *findings ledger* (typed entries citing rule ids + evidence tiers), not numbers. This is the most important contribution in round 01.
- The one-seeded-re-diverge-then-accept-as-data escape for G1. Matches my F4 position exactly: cross-model divergence is the default, seeded assignment fires only on G1 failure, and a second collapse is accepted as `verified-genuine` with a Decider-recorded reason. Never a coin flip, never an infinite re-roll.
- The graft prohibition on touching the winner's token file: grafts may not add or alter tokens, only component-spec-level details re-expressed in existing tokens. This is stricter than my round-01 and correctly so — a graft that changes the token system is a merge, not a graft.
- The deterministic read-order derivation: `sha256(agent-id || run_id) mod 4!`. This is a concrete, checkable mechanism for order randomization that works on disk where a facilitator shuffle cannot. Adopt verbatim.
- The 4-file / 40 KB ceiling for the always-loaded core. Tighter than my round-01 and correct — the core is read by every agent on every run.

**Rejecting:**

- The tier vocabulary `text|source|static-dom|layout|pixel`. I am adopting the research brief's `text-regex|css-parse|dom|browser|screenshot|human` — it names the actual engine, not an abstraction. "layout" is ambiguous (layout from what?); "browser" is unambiguous. One vocabulary must win and this one is the most descriptive.
- The 400-word body cap on directions. Too tight — a real token table (scale, roles, radii, durations) consumes 100+ words alone. claude-1's 750-word cap is more realistic. I adopt 750 words + token tables.
- `PDS/1` without a minor version. Should be `PDS/1.0` — the minor version is where deprecation and rule additions land. A bare major is insufficient.

## Resolved disagreements A-G

### A. Second state machine, or reuse Parley's?

**Parley's phases win.** `parley-design` is a profile that layers design-specific gates onto Parley's existing rounds, not a parallel D0-D9 state machine. What is lost by the state-machine approach: alignment with Parley's isolation, signoff, artifact ownership, and track conditioning — plus the worktrees/tracker precedent that add-ons never duplicate canonical infrastructure. What is lost by the profile approach: nothing the checker cannot enforce. G1 (distinctness) is a checker rule that runs between round-01 and round-02. G2 (coherence) runs after graft. G3 (token-integrity) runs at systematize. G4 (no-open-violations) runs at audit. Each gate is `parley-design-check` emitting an exit code, not a phase transition. The research brief's D0-D9 ladder is documentation of the design workflow, not an enforceable state machine — and it contradicts the brief's own anti-goal #4.

### B. Registry source of truth

**Literate markdown wins.** One file (`RULES.md`), fenced ```yaml blocks as machine source, prose as human source, same file, cannot drift. What breaks under `rules.json` + generated markdown: the doctrine skill cannot be edited without running a generation tool, which breaks the zero-dependency promise. What breaks under a markdown registry the checker parses *without* fenced blocks (claude-1's version): there is no machine-readable contract between the prose and the checker — the checker would need to NLP-parse free text to extract rule metadata. The fenced YAML block is the minimal parser contract: regex for fenced blocks + YAML parse. No generation, no drift guard, no second copy.

### C. Size and file count

**4 files, ~60 KB hard ceiling for the skill.** SKILL.md (≤12 KB dispatcher) + PDS.md (≤24 KB protocol) + RULES.md (≤18 KB literate registry) + WEB-ANNEX.md (≤6 KB, surface-conditional). Per-project `design/` artifacts are separate and track-conditioned (fast: 4 files, standard: 7, deliberation: 9). codex-1's 96 KiB is too large — under 4 agents that is 384 KB of doctrine, and the research brief's anti-goal #14 is explicit: file sprawl is a failure mode under N agents. kimi-1's 40 KB is aspirationally tight but does not leave room for a real rule registry with 30+ entries. 60 KB is the honest floor.

### D. Evidence-tier vocabulary

**`text-regex|css-parse|dom|browser|screenshot|human` wins.** It names the actual engine, not an abstraction. My round-01 `stated|source|rendered|measured` is ambiguous (what is the difference between "stated" and "source"?). codex-1's T0-T3 is compact but opaque — a second implementer needs a glossary. kimi-1's `text|source|static-dom|layout|pixel` is close but "layout" is ambiguous. `screenshot` is reserved (not shipped in v1 — pixel checks are excluded per the research brief and my round-01). `human` is the agent-judgement tier for rules that cannot be mechanically decided (`pov-absent`, `decor-unmotivated`, `self-repetition`). The `unjudgeable: <tier>` verdict is adopted unanimously — no disagreement there. One vocabulary, four verdicts: `pass|violation|needs-review|unjudgeable`.

### E. The dice (F4)

**Cross-model divergence is the default; seeded forced-axis assignment fires ONLY on G1 failure.** Always-on seeding (codex-1, claude-1) pays the mechanism cost on every run to fix a failure mode that G1 already catches. Worse, forcing an agent to build its rank-k direction when it believes in its rank-1 can produce a direction it cannot defend — which weakens the distinctness that divergence is supposed to produce. The seed is `sha256(run-id)`, local, checkable, recorded in the artifact's `seed:` field. A direction that should carry a seed (G1 failed, re-diverge ordered) but does not is a material finding. What is checkable: G1 is computed over declared axes before any critique; the re-diverge seed is deterministic from the run id; the assignment is reproducible by any party. Cross-model heterogeneity alone is *correlated* divergence (shared training distribution), sufficient to break the single-model argmax rut, not sufficient to guarantee orthogonality — which is why G1 exists as the check, not the dice.

### F. Where the two skills split, and the exact contract

**`parley-design` owns everything that requires a decision; `parley-design-check` owns everything a script can decide with no model in the loop.** The contract is the literate registry (`RULES.md`):

- Each rule is an H3 heading + a fenced ```yaml block + prose. The YAML block carries: `id`, `class: defect|tell|system`, `tier`, `severity` (Nielsen 0-4), `targets[]`, `enforced-by: check|agent-judgement|both`, `yields-to: []`, `added:`, `confidence:`, `sources:`, `status:`, and optionally `deprecated:`/`replacement:`.
- The checker parses fenced blocks (regex + YAML parse), dispatches detection by `id`, uses `tier` to select the engine, and uses `class` to set exit-code semantics.
- The checker's `conformance.json` is *generated* from its detector implementations (not hand-maintained), declaring which rule ids it has detectors for and what tier each reaches. A doctrine rule with `enforced-by: check` that has no corresponding detector is reported as `unjudgeable`, never silently passed.
- Findings cite the rule `id` verbatim — the error string IS the rule text (C7 from the research brief). Findings are `rule-id — violation — remedy`, always all three, copied verbatim, diffable across runs.
- The doctrine is fully usable with the checker absent: every rule is human-readable prose an agent can apply by reading. The checker makes it cheap and CI-able; it is never the only path.

### G. Actively harmful

1. **The research brief's SCORECARD.md with trimmed-mean scoring.** At 38% human aesthetic agreement, aggregating four agents' rubric scores into a trimmed mean manufactures a number with false precision. The brief's own anti-goal #8 says "A numeric design score / a headline grade" is forbidden — then it proposes a scorecard. kimi-1's findings-ledger is the correct replacement: agent input to the Decider is typed findings citing rule ids + evidence tiers, not numbers.

2. **The research brief's RUMBLE branch (D5b).** It contradicts binding D3 ("one direction wins in its entirety — never an average, never a merge"). Building both directions is prototype implementation, which D8 explicitly excludes ("the skill does not own Phase-5 code"). It doubles cost and opens an escape hatch that will grow — "rare and justified" today is "default for hard calls" in six months. codex-1 correctly rejected this. If two directions are genuinely incommensurable, the answer is `ABSTAIN` and a separately scoped experiment idea, not a Rumble inside the design protocol.

3. **The research brief's 17-artifact set for the deliberation track.** It violates the brief's own anti-goal #14: "File sprawl... under N agents it is 5× the cost and a guaranteed drift surface." 17 artifacts × 4 agents = 68 reads. The brief collapses 9 design/ files into 2 (DESIGN-SYSTEM.md + LEDGER.md) in §2.5, then re-expands the artifact table to 17 entries in §2.3. Pick one. My position: standard track = 7 artifacts, deliberation = 9, fast = 4.

4. **claude-1's hand-maintained `conformance.json`.** A second copy of capability declarations that will drift from the registry. Generate it from detector implementations or do not ship it.

## Final positions on F1-F8

- **F1 — Three classes: `defect` (single-agent BLOCK, objective threshold), `tell` (quorum ≥2 agents citing same rule id, taste with strong prior), `system` (single-agent BLOCK, binding post-ratification only). Honesty rule is `defect` with ethical basis. At least one `defect` rule (contrast floor) is design-system-blind.**
- **F2 — Convergence is an ALARM. G1 fires on all-axis match or banned-slop signature match. One seeded re-diverge, then accept as `verified-genuine` with a Decider-recorded reason. Never an automatic pass.**
- **F3 — Optional with declared capability. Vocabulary: `text-regex|css-parse|dom|browser|screenshot|human`. `unjudgeable: <tier>` is compliant. Degradation banner mandatory. No agent signs a layout verdict it never saw.**
- **F4 — Cross-model divergence is the default. Seeded forced-axis assignment fires ONLY on G1 failure. Seed is `sha256(run-id)`, local, checkable, recorded in artifact. Always-on seeds over-engineer for a failure G1 already catches.**
- **F5 — 4 skill files, ~60 KB hard ceiling. Per-project design/ artifacts are track-conditioned (fast: 4, standard: 7, deliberation: 9). Identical core reads under N agents are the point (shared objective function).**
- **F6 — Waivers in `design/waivers.md`, counter-signed. `defect`-class: human ratification at FINAL. `tell`-class: second participant counter-signature. `system`-class: Decider + implementer. Narrowest-exception ladder enforced. All `defect`-class a11y rules + honesty rule are design-system-blind.**
- **F7 — Fast path: single agent, core invariants, existing checker. Allowed only when ratified system exists AND single surface AND no new tokens AND no new foundation. Everything else runs full ritual. The ritual's cost is justified at direction-creation time only.**
- **F8 — Human Decider by default, agents advisory. Agent input is findings ledger, not scores. ABSTAIN is preserved and escalates to human. In unattended runs, stall at DECIDE — do not auto-elect by vote. Voting is the documented failure mode.**

## New concerns / questions

1. **DTCG `2025.10` schema URL must be vendored.** The research brief references `designtokens.org/schemas/2025.10/format.json` as the validation target. Offline validation cannot depend on a live URL. codex-1 flagged this; I concur. Vendor the schema + license + provenance, record the upstream URL and hash, but never fetch at validation time.

2. **The 10-state component spec is right but heavy for `fast` track.** The expanded state list (default, hover, focus-visible, active, selected, disabled, read-only, loading, error, success) improves on my round-01's 8. But requiring all 10 on every component in a fast-track single-surface change is overkill. Proposal: fast track requires the 5 core states (default, hover, focus-visible, active, disabled); standard and deliberation require all 10.

3. **`self-repetition` on day zero.** Our clearest differentiator — structural distance from the project's own prior outputs — needs LEDGER.md history. On the first design run, there is no history. Day-zero behavior: `self-repetition` is `unjudgeable` (no history to compare against), not `pass`. This is honest — we cannot claim distinctiveness from ourselves when we have no prior self.

4. **Genre/profile exemption scoping.** The research brief's §3.7 disables the entire `tell` class for enterprise CRUD, regulated UI, accessibility-driven forms, etc. Correct in principle — convention is the right answer in those domains. But who declares the genre? BRIEF.md, Decider-ratified. What stops an agent from declaring `genre: enterprise` to escape all slop findings? The Decider ratifies the genre at D0, on the record, and `defect`-class rules are never exempt. A genre declaration that contradicts the brief's actual scope is a material finding.

5. **APCA stays advisory.** WCAG 3 is a Working Draft. Contrast was pulled from WCAG 3 in July 2023. The algorithm is "yet to be determined." No Recommendation before 2028. Do not fail a build on APCA. This is unchanged from round 01 and I will object if it drifts.

## Current proposal

**Protocol:** `PDS/1.0`, RFC 2119 normative language, profile over Parley rounds (not a second state machine). 4 skill files: SKILL.md + PDS.md + RULES.md + WEB-ANNEX.md, ~60 KB ceiling. Design-specific gates (G1-G4) are checker-enforced rules between Parley phases, not parallel phase transitions. Conformance: 4 levels (L1 artifacts, L2 process, L3 token-integrity, L4 applied-UI), each checkable. Spec versions two-part (major.minor); tool versions independently with `implements: PDS/1.x`.

**Registry:** Literate markdown (`RULES.md`). Fenced ```yaml blocks are machine source; prose is human source; same file, cannot drift. Each rule: `id`, `class: defect|tell|system`, `tier`, `severity`, `targets`, `enforced-by`, `yields-to`, `added`, `confidence`, `sources`, `status`. No generation step, no second copy. Checker parses fenced blocks; `conformance.json` is generated from detector implementations, not hand-maintained.

**Token schema:** W3C DTCG `2025.10` adopted verbatim for format. Three-tier layering: `primitive` (raw OKLCH values, no semantic names) → `semantic` (aliases only) → `component` (aliases to semantic only). Reference direction enforced as a graph assertion: never skip, never reverse. **OKLCH mandatory** for all color primitives, 6-digit hex fallback. HSL/RGB prohibited in token files. Resolver: `light` + `dark` required; `hc`, `density`, `reduced-motion` recommended. Schema vendored offline. Named Rules (`**The [Name] Rule.** <one sentence>`) for what DTCG cannot hold: motion posture, focus treatment, density, effect budget. Component spec: 13 sections (Purpose → When to use → When NOT to use → Anatomy → Variants → Sizes → States (10 rows) → Behaviors → Interactions → Content → Accessibility → Do/Don't → Tokens used), merged from Carbon + Polaris + GOV.UK patterns.

**WCAG 2.2 as blocking constants:** SC 1.4.3 (4.5:1 body, 3:1 large) — violation. SC 1.4.11 (3:1 non-text) — violation. SC 1.4.10 (reflow at 320px) — violation. SC 1.4.12 (text spacing survivability) — violation. SC 2.2.2 (carousel pause) — violation. SC 1.3.1 (heading order) — violation. SC 2.5.8 (24×24 target size) — needs-review (5 legitimate exceptions exist; flag, require disposition, do not hard-block). APCA — advisory only, never a gate. All `defect`-class a11y rules are design-system-blind: being on the token ramp never exempts contrast, target size, motion safety, or honest copy.

**Rule classes and authority:** `defect` (objective threshold, single-agent BLOCK, design-system-blind for a11y) · `tell` (taste with strong prior, quorum ≥2 agents to block) · `system` (conformance to ratified contract, single-agent BLOCK, binding post-ratification only). `yields_to` stand-down: one defect = one owning rule id. Slop class is disabled by genre exemption (enterprise CRUD, regulated UI, etc.); `defect` class is never exempt. Count-gated: ≥4 concurrent `tell` findings on one artifact = the empirically anchored fail threshold (Krebs, N=1,590).

**Ritual:** Round-01 = DIVERGE (isolation, 750 words + token table, one-word handle, Signature field, position on every declared axis). G1 = distinctness gate (5-axis fingerprint, checker-enforced, fires before round-02). Round-02 = CRITIQUE (one round, assigned lenses, typed findings citing rule ids + evidence tiers, author silent on own direction, deterministic read-order `sha256(agent-id||run-id) mod 4!`). Consensus = DECIDE (human Decider, findings ledger not scores, exactly one winner whole, 0-3 grafts from heat-map `like` clusters, grafts re-expressed in winner tokens only, grafts never touch token system). FINAL.md = CONTRACT (winner + Named Rules + token table + grafts). Post-Phase-5 = DOCUMENT (`DESIGN-SYSTEM.md` written from shipped code, reconciliation table classifies divergence as `match|acceptable-adaptation|missing|contradicted|added-without-approval`). Audit = checker at highest available tier, `AUDIT.md` + `LEDGER.md`, G4 blocks RATIFIED on open `defect` violations.

**Evidence tiers:** `text-regex|css-parse|dom|browser|screenshot|human`. Verdicts: `pass|violation|needs-review|unjudgeable`. `unjudgeable: <tier>` is compliant. Degradation banner mandatory. `screenshot` reserved (not v1). v1 ships `text-regex` + `css-parse` + `dom` (zero runtime deps); `browser` is optional auto-detected adapter.

**Dice:** Cross-model divergence is the default. Seeded forced-axis assignment (`sha256(run-id)`) fires only on G1 failure. One re-diverge, then accept as `verified-genuine` with Decider-recorded reason. Never always-on, never infinite re-roll.

**Decider:** Human by default. Agent input is findings ledger (typed entries, not scores). ABSTAIN is preserved and escalates. In unattended runs, stall at DECIDE — no auto-election by vote.

**Cut order if scope shrinks:** (1) browser tier (T2/T3 carry ~60% of value at ~0% of cost); (2) component-spec templates beyond tokens + floor + Named Rules; (3) lens assignment (keep one-round critique, drop the hats); (4) web annex depth. Never cut: the distinctness gate, the findings-ledger selection rule, the tier-declaration machinery, the literate registry, OKLCH mandate.
