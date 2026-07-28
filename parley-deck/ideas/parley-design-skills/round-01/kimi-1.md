---
agent: kimi-1
idea: parley-design-skills
round: 1
date: 2026-07-28
---

## Summary

My lens is the collaborative ritual and protocol shape. Central position: **`parley-design`
must be a profile over the existing Parley deliberation track, not a new phase state
machine.** DCP/1 (research/design-critique-method.md §7) gets the human-ritual mechanics
right — Heat Map, one critique round, Decider not vote, bounded graft, Rumble branch — but
as written it invents P0–P8 with its own artifacts and a parallel state machine, which is
anti-goal #3 in the brief. I adopt ~80% of DCP/1's mechanics and re-key them onto Parley's
rounds: DIVERGE = round-01 (isolation already exists there), CRITIQUE = round-02 (one
round), DECIDE = consensus.md (Decider act + signoff ratification), GRAFT = bounded edits
in consensus.md, CONTRACT rides in FINAL.md (D7), SYSTEM is documented post-Phase-8 from
shipped code (D7's second artifact). I diverge from DCP/1 on four points (Q3): (1) no
holistic scorecards — agent input to the Decider is a *findings ledger*, not 0–10 rubric
scores; (2) anonymity demoted to SHOULD — file-level authorship makes it leaky, recusal is
the real mechanism; (3) the distinctness escape is one seeded re-diverge, then convergence
is accepted as data; (4) grafts are mechanically forbidden from touching the winner's
token file.

On the negative evidence: multi-agent is justified **only** as a diversity generator plus
an adversarial finding-filer — never as a taste amplifier. Cross-model divergence is real
but correlated (shared training distribution): sufficient to break the single-model argmax
rut (the measured 30/35), not sufficient to guarantee orthogonality. Everything below
protects that one advantage and keeps selection rule-based and Decider-owned.

## Proposed approach

### Q1 — Protocol shape of `parley-design`: `PDS/1`

Spec id `PDS/1`, first version `PDS/1.0`, carried as `spec: PDS/1.0` in the frontmatter of
every artifact the skill produces (fixing AG-UI's no-version-on-the-wire defect).
Conformance language: RFC 2119 uppercase MUST/SHOULD/MAY, declared in SKILL.md frontmatter
as `conformance-language: RFC 2119`. Versioning: semver of the *spec*, independent of the
tooling; rule ids are append-only and never renumbered (insert like `38a`); a deprecated
rule keeps validating for ≥1 minor version and is listed in a deprecation table with
introduced/removed versions. Deprecation stance: rules may be narrowed only with a recorded
false-positive citation (the bug that forced the narrowing, verbatim — impeccable's FP-fix
convention), because a silently narrowed rule invalidates every past review that cited it.

Table of contents (the always-loaded core is §0–§6 only):

```
SKILL.md (the spec, ≤ 400 lines)
§0 Scope, non-goals, relationship to COOPERATION.md (addon, never amends it)
§1 Terminology (Direction, Finding, Graft, Named Rule, Decider, Profile, Tier)
§2 Design principles (≤6: one-winner-whole; findings-not-scores; absent=undeclared;
   unknown rule ids MUST NOT error; gate on mechanics, advise on taste)
§3 Roles and authority split (Proposer, Critic, Decider, Implementer, Documenter;
   Critics hold no decision authority; the Decider never critiques)
§4 The ritual as a Parley profile (phase→round mapping, entry/exit conditions,
   illegal transitions with the exact error strings the checker emits)
§5 Canonical artifacts and schemas (one identically-shaped entry each: H3 name,
   one-line purpose, rationale, REQUIRED frontmatter table, REQUIRED body sections)
§6 Collaboration contract rules (numbered, bold-named, citable — AG-UI interrupts.mdx style)
reference/rules.md          — the rule registry, doctrine form (~40 rules, see Q4)
reference/floor.md          — the craft floor, ≤ 50 lines, hard numbers, one absolute ban
reference/web-annex.md      — every CSS-shaped rule, explicitly labelled (D5)
reference/schemas.md        — artifact frontmatter JSON Schemas + component-spec template
reference/anchors.md        — worked anchor examples per severity level (lazy)
```

**Conformance** is defined in four levels, each checkable: **L1 artifact-shape** (required
artifacts exist, frontmatter validates — artifact linter); **L2 ritual-order** (round-01
files predated any round-02 read; exactly one critique round; distinctness gate evaluated —
attested per §4 plus transport branch structure); **L3 contract-integrity** (the built diff
contains no off-contract values — `parley-design-check` source tier); **L4
system-integrity** (DESIGN-SYSTEM.md matches shipped code — the post-build documentation
check). "An implementation conforms to PDS/1.0 at level N iff…" is written out per level;
`parley-design-check` is the reference runner for L1/L3/L4, L2 is attested + spot-audited.

### Q2 — The design-system artifact set (what exists on disk, per project)

```
design/
  tokens.json            — W3C DTCG 2025.10, THE machine source of truth for values
  DESIGN-SYSTEM.md       — human authority: Named Rules + MUST-share/MAY-differ sections
  components/<name>.spec.md — one per component, Carbon-status frontmatter
  DESIGN-DECISIONS.md    — ADR log; every entry ≥2 alternatives + why-not
  waivers.md             — the only place a rule is waived (see F6)
```

- `tokens.json` is DTCG 2025.10 **verbatim**: `$value`/`$type`/`$description`, aliases via
  `{group.token}`, three-tier layering `primitive.* → semantic.* → component.*` with
  reference direction enforced (component→semantic→primitive only, never reverse — a graph
  assertion). Colors authored `colorSpace: "oklch"` + 6-digit `hex` fallback. Resolver
  (modes) is v1-optional: `light` required, `dark` recommended; `hc`/`reduced-motion` later.
- Deliberate divergences from DTCG: (a) we ship **no** Style Dictionary/Terrazzo dependency —
  the doctrine names the format, the check skill validates against the official published
  JSON Schema; (b) non-spec token groups use `x-` prefix, not `$extensions` reverse-domain
  keys, because our consumers are agents not vendor tools; (c) `DESIGN-SYSTEM.md` is primary
  narrative and tokens.json is subordinate to it for *intent* — impeccable's split
  ("tokens normative, prose contextual") inverted for a multi-agent world: values live
  exactly once (tokens.json), rules live exactly once (Named Rules), never both.
- `DESIGN-SYSTEM.md` frontmatter: `spec`, `status: draft|preview|stable` (Carbon PDLC
  enum), `registry_version`, `derived_from` (the idea slug + content hash of the shipped
  diff). Body: `## Named Rules` (each `**The X Rule.** one sentence`, 1–3 per section),
  `## What surfaces MUST share` / `## What surfaces MAY differ on` (the parallelisation
  contract for Phase-5 implementers), `## Reconciliation` (see D7 below).
- Component spec: `status` + sections `Purpose / When to use / When NOT to use (must name
  the alternative) / Anatomy / Variants / States (all 8 rows, token refs only) / Behaviors /
  Content / Accessibility / Do-Don't (≥3 pairs) / Tokens used (must resolve)`.

**D7 reconciliation, concretely:** the FINAL.md `## Direction Contract` section (winner +
Named Rules + token table) binds implementers *before* the build; `DESIGN-SYSTEM.md` is
written *after* Phase 8 by a Documenter role whose ground truth is the merged diff. When
they contradict: **the build wins on values, the contract wins on intent.** Every
contradiction lands in `## Reconciliation` with one of two dispositions: `build-fixed`
(contract violation, implementer repairs) or `contract-amended` (a recorded waiver with
counter-signature). Silent divergence is itself a `quality`-class finding.

### Q3 — The ritual, made mechanical (D3 onto Parley rounds)

- **Phase 0.** `00-prompt.md` gains a required `## Design Brief` section when the addon is
  active: `goals:` (each with an id, cited later by findings), `constraints:`,
  `divergence_axes:` (3–5 named axes, e.g. `structure, type-voice, color-strategy, motion,
  density` — no axes, no round-01), `surface_mode: persuade|operate|read|experience`,
  `decider: human` (default; see F8), `design_track: full|fast`.
- **round-01 = DIVERGE.** Each agent's own round file carries the direction schema:
  frontmatter `design-direction: true`, `handle:` (one word, not the agent id),
  `axes: {structure: <position>, ...}` (a declared position on *every* brief axis),
  `seed:` (only if assigned, see F4); body hard-capped at ~400 words in five blocks
  `THESIS / WORLD / TOKENS / RULES / FIRST-SURFACE` (a block that reads like a mood = not
  decided), a token table, and a `## Graftable parts` list (the author's own 3 best
  details). Crazy-8 lives here as 8 one-line variants under `## Variants considered` —
  intra-agent divergence with no new artifact. Isolation is the protocol's existing
  round-01 rule; the skill only restates it as a MUST and forbids pre-reads.
- **Gate G1 DISTINCTNESS** (between rounds): mechanical over round-01 frontmatter.
  Collapse = ≥3 of 4 directions match on *all* declared axes, or all 4 match on the primary
  axis. Escape: exactly **one** seeded re-diverge (each agent assigned a distinct position
  on the primary axis, assignment derived from `sha256(run-id)` — recorded in `seed:`).
  Second collapse: accept, record `convergence: verified-genuine` in consensus.md. Never a
  coin flip, never an infinite re-roll.
- **round-02 = CRITIQUE, exactly one round.** Lenses assigned deterministically (rotate
  Black/Yellow/White/Green over `sha256(run-id)` — checkable, no facilitator artifact).
  Findings are typed, one per line, JSONL block:
  `{target, part, class: like|wish|what_if, severity: 0-4, ref: <rule-id|goal-id>, tier: text|source|static-dom|layout|pixel, note}`.
  A `wish` with no `ref` is dropped as taste. `like` marks are the graft shortlist.
  An agent never critiques its own direction; it gets one `## Rebuttal` section addressing
  misreadings only. Fixes are optional and labelled non-binding (Braintrust/plussing
  reconciliation: diagnosis owed, prescription not owned). **No holistic 0–10 scores
  anywhere** — see F8 for why the findings ledger replaces the scorecard. File transport
  note: per-reviewer order shuffling is impossible on disk, so each reviewer MUST read in
  the order `sha256(agent-id || run-id) mod 4!` — deterministic, checkable, and it kills
  the fixed-position bias DCP/1's facilitator shuffle cannot reach under github-pr.
- **consensus.md = DECIDE.** Contents: `winner:` (exactly one direction id — a
  discriminated union, mechanically not an average), `grafts:` (≤3, each
  `{source, part, why, reexpressed_as: <winner token>}`), `dissent:` (verbatim),
  `convergence: collapse|reseeded|verified-genuine|divergent`, `evidence:` (tiers reached,
  per signer), degradation banner if any tier/participant missing. Grafts may **not** add
  or alter tokens in the winner's token table — component-spec-level details only,
  re-expressed in existing tokens; a graft that cannot be re-expressed is rejected. G2
  COHERENCE: the check skill re-runs post-graft; new violations fail the graft, not the
  winner. **Re-litigation stop:** losers archived verbatim as `maybe-later`; the decision
  is pinned by a 12-hex content hash; re-opening requires a *new* idea whose prompt cites
  that hash. Frozen dissent, not recurring dissent.
- **FINAL.md** carries `## Direction Contract` (winner, Named Rules, token table, grafts)
  per D7 — unchanged FINAL semantics, addon-scoped section.
- **Post-Phase-8 = RATIFY.** The Documenter pass (one agent, assigned in the
  implementation idea) writes `DESIGN-SYSTEM.md` + `tokens.json` from the shipped diff.
  Multi-agent adds value at direction time and review time; at coherence time it is a
  liability — so the Documenter is singular by design.

### Q4 — Slop doctrine: depth and load model (F5 answered here)

Always-loaded core (every participant reads it once, round-01): SKILL.md ≤ 400 lines +
`rules.md` + `floor.md` + `schemas.md` — **hard ceiling 4 files / 40 KB**. Lazy references:
`web-annex.md`, `anchors.md` — ceiling 2 more files, read per surface type. v1 ships **~40
rules**, not 60: the subset defensible as cross-project invariants, each
`{id, category: quality|slop, authority: block|quorum|advise, tier, tell, why, fix}` at ~1
screen each. Anti-house-cliché defence: the registry ships *tells*, never a look; the
category-plus-avoidance calibration test is itself a judgement-tier rule; all free axes are
owned by the per-deck ratified system (D6) and kept off the training prior by the
`divergence_axes` mechanism + G1 (the axes force declared positions; the gate measures
them; the seeded escape forces them when models drift home).

### Q5 — The split line and the machine-readable contract

`parley-design` owns: ritual, registry *as doctrine*, schemas, floor, Named-Rule format.
`parley-design-check` owns: `rules.json` (data-only registry, single source of truth —
doctrine §rules table is *generated* from it, drift-guarded by a test), the artifact
linter, the DTCG schema validator, source-tier detectors, the optional browser adapter,
exit codes, the waiver validator. The contract between them, four pieces: (1) `rules.json`
schema `{id, category, authority, tier, severity, name, tell, fix, since, deprecated?}`;
(2) per-project `design-check.config.json` (token globs, profile, waivers path, own JSON
Schema); (3) finding format `{rule_id, file, line, tier, verdict, evidence}` with
`verdict: violation|needs-review|recommendation`; (4) namespace rule: `core:` reserved for
spec rules, project rules MUST be `<project>:<slug>`, consumers MUST NOT error on unknown
ids. The doctrine cites checks by `core:<slug>` it cannot run; a doctrine reference to a
nonexistent id fails the drift guard.

### Q6 — `parley-design-check` design

Registry as data, detection separate (impeccable's shape, our words). Tiers: **T0 text**
(zero-dep, ~19 rules: banned fonts, single-font, flat hierarchy ratio < 2.0, monotonous
spacing > 0.6 dominant over ≥10 samples, gradient-text, purple/cyan hue bands, cream
(`min≥209, r≥g≥b, 6≤r−b≤48`), bounce-bezier outside [−0.1,1.1], em-dash 8+/1-per-500,
buzzwords, raw hex outside token block, `transition: all`, `!important` ratio) — ships v1,
runs anywhere. **T1 CSS-AST** (off-scale values vs tokens.json, ΔE2000 < 1.0 duplicate
colors error / 1.0–2.3 needs-review, z-index sprawl, non-compositable animation, missing
reduced-motion fallback) — ships v1. **T2 static DOM** (heading order, alt, lang) — v1.
**T3 browser** (computed contrast, 24×24 target size *own rule*, overflow sweep
320/375/768/1280, two-line clickable text, axe run) — optional adapter, auto-detected,
honest `unavailable` report, never silently skipped. T4 pixel/baseline: refused. Severity:
`violation` blocks the phase gate; `needs-review` requires a written disposition by a named
agent in the review artifact; `recommendation` advisory. Exit codes `0 clean / 1 violation
/ 2 warn-only / 3 config error`. FP suppression: inline `pdc-disable-next-line <rule> --
<reason>` with the reason **required and preserved** (impeccable's syntax, not its
reason-discarding semantics); config waivers per F6. Gate semantics in a Parley review:
the driver or any reviewer runs it; findings enter the review artifact as typed findings
citing rule ids — mechanical output enters *after* each reviewer's judgment is written
(anti-anchoring ordering), never before.

### Open forks — my positions

- **F1 Rule authority: split, as data.** `authority:` lives in the registry per rule:
  `block` (objective, mechanical threshold — contrast, occlusion, missing state, token
  drift; single agent may BLOCK), `quorum` (taste with a strong prior; BLOCK needs ≥2 of 4
  concurring, a single agent's finding is advisory), `advise` (never blocks). Recategorising
  a rule = spec minor bump + changelog entry + `registry_version` cited in every review, so
  old reviews stay interpretable.
- **F2 Convergence = ALARM, with a bounded escape.** Trigger: the G1 definition above
  (≥3-of-4 all-axis match, or 4-of-4 primary-axis match). Escape: one seeded re-diverge;
  then accept as `verified-genuine`. Convergence is suspicious once, evidence twice.
- **F3 Rendered evidence: optional-with-declared-capability.** Tier vocabulary
  `text|source|static-dom|layout|pixel`; every critique/review artifact declares tiers
  reached; `unjudgeable` is a first-class per-rule verdict; a signature scopes to tiers
  (`signed: text,source` attests those tiers only — nobody signs a layout verdict they
  never saw). Missing tier or missing participant ⇒ `⚠ DEGRADED:` banner leads the
  consensus artifact. Silent degradation is a conformance error.
- **F4 Dice: free by default, seeded as the G1 escape only.** Cross-model divergence is the
  default randomiser; the seeded forced-axis assignment fires only on G1 failure, seed
  `sha256(run-id)` (local, checkable, recorded in `seed:`). Always-on seeds would pay the
  mechanism cost on every run to fix a failure mode that only sometimes occurs.
- **F5 Ceiling: 4 always-loaded files / 40 KB; ~40 rules; 2 lazy files.** Stated in Q4.
  Under N agents, identical reads of the core are the *point* (shared objective function);
  4 × 40 KB is cheap; divergent lazy reads are bounded by the per-surface trigger.
- **F6 Waivers: one file, counter-signed.** `design/waivers.md` entries
  `{rule_id, value, scope, reason, author, seconded_by, created, context_hash}`.
  `advise`-tier: single author. `quorum`-tier: second participant counter-signature.
  `block`-tier and any whole-file/whole-rule waiver: human ratification at FINAL.
  Legibility floors (min text size, contrast) are deliberately design-system-blind: being
  on the token ramp never exempts them — the anti-laundering rule.
- **F7 Fast path.** `design_track: fast` (or the idea's `fast` track): one agent, one
  direction, self-check against the registry, the track's existing refutation-default
  reviewer does a rule-cited review, human gate decides. Threshold: fast is allowed only
  when a ratified design system already exists *and* the surface count is 1 *and* no new
  tokens are needed; everything else runs the full ritual. The ritual's cost is justified
  at direction-creation time only.
- **F8 Decider: human by default; agents file findings, not scores.** Selection and
  ratification are different acts: **selection** is a Decider act (human; the 00-prompt MAY
  name one agent as Decider-Delegate for `fast` or low-risk surfaces, human override
  retained at the FINAL gate); **ratification** is the existing signoff quorum, unchanged.
  Agent input to the Decider is the findings ledger (net unresolved `block` findings,
  `like`-mark counts, dissent) — **not** holistic rubric scores. This survives the
  judge-bias evidence: recusal handles self-preference, findings-are-falsifiable handles
  the 38% ceiling, and feasibility bias is countered because the Decider sees the bold
  direction's findings *and* its likes, not a safety-weighted mean. `ABSTAIN` is a
  preserved verdict that escalates to the human, never coerced.

## Concerns / open questions

- L2 ritual-order conformance is attestation-based; file mtimes are noise across worktrees
  and clones. The honest v1 answer is attestation + branch structure; anything stronger
  wants driver support (a `parley design check` step), which I flag as a possible seam for
  the integration lens — as an addon we must not silently require it.
- Whether `consensus.md`/`FINAL.md` addon-scoped sections need a meta-protocol-change idea:
  I believe not (additive sections, ownership unchanged, worktrees/tracker precedent), but
  the synthesis round should confirm.
- The human Decider is a latency bottleneck for unattended runs; the Decider-Delegate path
  needs its risk wording nailed in round-02, especially who may appoint a delegate.
- Lens rotation and read-order derivation assume a deterministic roster ordering and a
  stable run-id — both true today, worth one contract-rule sentence so they stay true.
- Anonymity: I demote it to SHOULD (order-rotation + recusal carry the real weight), but if
  the facilitator can mechanically relabel direction files under github-pr transport, full
  anonymity is cheap to restore. Transport-dependent; needs the integration lens's answer.

## Risks

- **Where multi-agent makes design worse** (state plainly in the skill): averaging
  operators (mush — the mean of four directions is closer to the training mode than any
  single one); quorum ratchet toward the category standard (the option nobody objects to);
  round multiplication (factual attrition, stance homogenization — one critique round is
  the cap); five blind reviewers multiplying confidence, not evidence (the tier-declaration
  machinery exists precisely for this); too many voices at coherence time (Documenter is
  singular). The advantage is squandered by: any pre-round-01 visibility, a prescriptive
  doctrine that collapses all four directions onto one output, a majority-vote selection,
  homogeneous roster configs, and letting mechanical findings enter before judgment.
- **Where this fails:** the human Decider rubber-stamps the ledger (mitigation: the ledger
  shows dissent and `like`-clusters, not a ranking); G1 is gamed by superficial axis
  differences (mitigation: axes are declared in the brief, positions are specific values,
  not adjectives); rule-registry false-positive fatigue (mitigation: FP-citation narrowing
  rule + `needs-review` tier absorbs ambiguity); token cost creep (the 40 KB ceiling is
  policy, and policy erodes — the drift guard must also fail on core-size growth).
- **Cut first, in order, if scope must shrink:** (1) the Rumble branch (rare by design;
  escalate to human instead); (2) the T3 browser adapter (T0–T2 carry ~60% of the value at
  ~0% of the cost); (3) component-spec templates beyond tokens + floor + Named Rules;
  (4) lens assignment (keep one-round critique, drop the hats); (5) Crazy-8 (keep the
  schema, drop the intra-agent step). Never cut: the distinctness gate, the findings-ledger
  selection rule, the tier-declaration machinery — those three are the evidence answer.
