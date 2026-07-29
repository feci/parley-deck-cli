---
idea: parley-design-skills
status: final
drafted-by: claude-1
date: 2026-07-28
consensus: consensus.md (C1-C13, U1, amendments A1/A2/A3; all four participants ACCEPT)
participants: [claude-1, codex-1, hermes-1, kimi-1]
implementation-target: parley-deck-skill (addons/parley-design, addons/parley-design-check)
---

## What ships

Two opt-in companion add-on skills in the `parley-deck-skill` package, discovered by the
installer's existing `discoverAddons()` glob and installed by default alongside
`parley-worktrees` and `parley-tracker`.

```
parley-deck-skill/addons/
├── parley-design/                 # doctrine + protocol, ZERO runtime dependencies
│   ├── SKILL.md                   # dispatcher + invariants (always loaded)
│   └── references/
│       ├── PDS.md                 # the protocol
│       ├── RULES.md               # the literate rule registry
│       └── WEB-ANNEX.md           # web-specific hard numbers (non-normative elsewhere)
└── parley-design-check/           # enforcement; MAY ship code
    ├── SKILL.md
    ├── bin/check.js               # the checker CLI
    ├── lib/registry.js            # literate-registry parser
    ├── lib/detectors/*.js         # one module per detector
    └── test/*.test.js
```

**Hard budget for `parley-design`: 64 KiB across the four markdown files**
(`SKILL.md` ≤8 KiB · `PDS.md` ≤20 KiB · `RULES.md` ≤24 KiB · `WEB-ANNEX.md` ≤12 KiB),
enforced by a test, not by a comment.

## Binding decisions carried from consensus

All of C1–C13, U1 as amended, and A1/A2/A3 are binding. The implementation-relevant core:

1. **`parley-design` is a profile over Parley's existing track.** No second phase cursor.
   The mapping table in C1 is normative and must appear in `PDS.md`.
2. **`RULES.md` is the single source of truth.** One `pds-rule` YAML fence per rule H3, prose
   in the same file. No generated view, no second copy anywhere — including inside the
   checker.
3. **Three rule classes with different authority:** `quality` (single-agent BLOCK on
   reproducible evidence) · `slop` (never blocks unilaterally; ≥2 independent non-author
   concurrences) · `system` (binding only after a contract is ratified).
4. **Evidence tiers are ordinal and always written number-plus-word:** `T0 ARTIFACT`,
   `T1 SOURCE`, `T2 RENDERED`, `T3 PIXEL`. `UNJUDGEABLE` is a verdict, not a tier. Engine
   names never appear in the core.
5. **G1 fails unless every pair of directions differs on ≥2 declared axes** (A1),
   unconditionally, plus banned-slop-signature and duplicate-Signature checks.
6. **An unattended full run records `ABSTAIN` and stops before CONTRACT / Phase 5** (A3).
   No agent-selected winner, provisional or otherwise, authorises implementation.
7. **No numeric aesthetic score, ever.** Findings ledger only.
8. **Selection, never averaging.** One direction wins whole; 0–3 grafts; a graft MUST NOT
   modify the winner's token file.
9. **Two artifacts with distinct authority:** the Direction Contract binds before the build;
   the Design System describes after it and **cannot retroactively legalise a violation**.
10. **Invariants ship; taste is ratified per deck.** No theme catalogue.

## `parley-design/SKILL.md` — required content

- Frontmatter: `name: parley-design`, a description that states it is a Parley Deck
  companion add-on for collaborative design-system work.
- **When to use / when NOT to use.** Explicitly: not for a component inside an already
  ratified system (that is the fast path), not as a taste oracle, not a theme catalogue.
- **The invariants**, short enough to hold in memory before editing anything: the honesty
  rule (no invented metrics, testimonials, logos or benchmarks — a labelled hole is honest,
  a fabricated number is slop), the interaction-state completeness rule, the contrast floor,
  the effect budget, and the precedence chain.
- **Precedence chain, stated once and binding:** `quality` rules > the ratified design
  system > the brief > parity with existing code > model habit. The top rules are explicitly
  **not** bypassable by "preserve structural parity", "mirror this reference" or "match the
  prior build" instructions.
- A dispatcher table: which reference file to load for which phase, with fixed read sets so
  every participant reads identical bytes.
- A pointer to `parley-design-check` and an honest statement that the doctrine is fully
  usable without it.

## `parley-design/references/PDS.md` — required structure

Frontmatter carries `spec: PDS/1.0`, `status: stable`, `conformance-language: RFC 2119`,
`registry: core-rules/1.0.0`, `registry-digest: <12-hex>`.

Normative-language rules, stated in §0: uppercase MUST/SHOULD/MAY only when normative;
every normative statement numbered **and** bold-named so agents can cite the name and tools
can cite the number; informative sections labelled `(informative)`; anything unlabelled is
normative.

Required sections:

| § | Content |
|---|---|
| 0 | Scope, non-goals, normative language, relationship to COOPERATION.md |
| 1 | The Parley mapping table (C1) — normative |
| 2 | Artifact kinds: `DESIGN-BRIEF`, `DIRECTION`, `CRITIQUE`, `VERDICT`, `CONTRACT`, `DESIGN-SYSTEM`, `AUDIT`, `WAIVERS`. Each gets an identically-shaped entry: name → one-line purpose → rationale paragraph → required-fields table → minimal example. **The shape never varies.** |
| 3 | Gates G1–G4 with exact failure conditions and the checker's error strings |
| 4 | The ritual: diverge (with the deterministic primary-axis assignment and the decline valve), one critique round, decide, bounded graft |
| 5 | Roles and invariants: Proposer, Critic, Facilitator, Decider; recusal; no self-scoring; no pairwise comparison in the deciding phase; length caps; declared degradation |
| 6 | Evidence tiers and verdicts |
| 7 | Rule classes and authority |
| 8 | Waivers |
| 9 | Conformance levels L1–L4 and what each requires |
| 10 | Extension policy: `core:` reserved; project rules `<project>:<slug>`; unknown rule ids and unknown token groups MUST NOT error |
| 11 | Versioning and deprecation; deprecated rules keep validating ≥1 minor version |
| 12 | Changelog (maintained, or the section is deleted — never left empty) |

**Conformance levels:** L1 artifacts exist and lint · L2 process order and gates recorded ·
L3 token integrity (DTCG `2025.10`, alias direction, no raw literals outside the token
layer) · L4 applied UI passes rendered-tier `quality` rules. L1+L2 need no runtime; L3 needs
a JSON validator; L4 needs a browser. A project declares the level it claims; the checker
verifies the claim.

**Nine defects not to reproduce** (from the AG-UI study, each with its countermeasure):
version on every artifact from the first commit; never two hand-maintained representations;
never write a count in prose; RFC 2119 uppercase reserved for normative use; maintain or
delete the changelog; zero placeholders; the checker runs standalone against files on disk;
fixtures are plain files and one script, offline; every generated view is drift-guarded.

## `parley-design/references/RULES.md` — the literate registry

Frozen extraction grammar, stated in the file itself:

- Each rule is an H3 whose text is the rule name.
- Immediately beneath it, **exactly one** fenced block tagged `pds-rule` containing YAML.
- Then prose: why it reads as machine-made or broken, a counterexample, and the remedy.
- UTF-8. Duplicate ids are fatal. Unknown keys warn unless `x-` prefixed. Unknown rule ids
  encountered by a consumer pass through as `UNJUDGEABLE` and MUST NOT crash.

Required YAML keys:

```yaml
id: core:focus-ring-animated
class: quality | slop | system
tier: T0 | T1 | T2 | T3
surface: core | web
enforced-by: check | agent-judgement | both
severity: 0-4                 # Nielsen anchors; only 4 (optionally 3) can block
added: 1.0.0
status: stable | draft | deprecated
sources: [...]                # where the rule came from
system-blind: true            # optional; cannot be waived by widening the system
```

**v1 registry content.** Ship a small, defensible set rather than a large borrowed one.
Required coverage, all written in our own words:

- **`quality`, surface `core`:** contrast floor; incomplete interaction states; honesty
  (invented metrics, testimonials, logos, benchmarks); unlabelled inference; text that
  cannot be read at the declared minimum size; motion without a reduced-motion path;
  focus indication absent or animated. At least the legibility floor is `system-blind`.
- **`slop`, surface `core`:** the aesthetic guessable from category alone *or from
  category-plus-avoidance*; decoration with no motivation; effect budget exceeded;
  structural sameness (the same macro-shape reused for a different brief); a direction whose
  Signature is absent or is a mood rather than a decision.
- **`system`:** value off the ratified scale; colour outside the ratified ramp; font outside
  the allowlist; token declared but unused, or used but undeclared.
- **`slop`/`quality`, surface `web`:** the specific CSS-shaped tells live in `WEB-ANNEX.md`
  and are referenced by id from here.

Rules are **append-only**: an id never changes meaning, and re-classification requires a
spec version bump.

## `parley-design/references/WEB-ANNEX.md`

Explicitly non-normative for non-web surfaces, stated in its first line. Carries the
tier→engine mapping for web (`T1 SOURCE` = CSS/HTML parse; `T2 RENDERED` = computed styles
and accessibility tree; `T3 PIXEL` = raster), the web-specific hard numbers (WCAG 2.2 ratios
blocking, APCA advisory; viewport widths; the overused-font list), and the web-shaped rule
bodies. A TUI/CLI annex is a deferred sibling and must be addable without touching the core.

## `parley-design-check` — required behaviour

- **Standalone.** Runs against files on disk. No agent runtime, no framework, no network at
  check time. Any schema it validates against is vendored with its upstream URL and hash
  recorded; it is never fetched during a run.
- **Registry loading.** Reads `RULES.md` from the installed `parley-design`. **If that skill
  is absent it MUST refuse rule checks and say so explicitly**; registry-independent
  structural and token checks may still run. It MUST NOT carry a bundled fallback copy.
- **Capability declaration is generated** from the detector modules present, never
  hand-maintained. A rule with `enforced-by: check` and no detector reports `UNJUDGEABLE`.
- **Every report carries** `implements: PDS/1.0` and the `registry-digest` it ran against.
- **Finding format:** `rule-id — violation — remedy`, always all three, on one line, stable
  and diffable across runs.
- **Verdicts:** `PASS | VIOLATION | NEEDS_REVIEW | UNJUDGEABLE`. Exit codes distinguish
  "clean", "violations found", and "the run itself failed". Findings are not errors.
- **v1 tiers:** `T0 ARTIFACT` and `T1 SOURCE` only. `T2` and `T3` are declared in the
  registry and reported `UNJUDGEABLE` — never silently skipped, never pretended.
- **Waivers:** read from one central file; each requires rule id, narrowest scope, reason,
  expiry and a counter-signature. `system-blind` rules cannot be waived by widening the
  system. Wildcards are rejected.

Commands: `check <paths...>` (rules), `check --level L1|L2|L3` (conformance), `check --json`.

## Testing

- Registry parses; ids unique; every id referenced by a detector exists; every
  `enforced-by: check` rule either has a detector or is reported `UNJUDGEABLE`.
- Fixture pairs per detector: one file that must fail, one that must pass. Offline, <5 s.
- Byte-budget test for the four doctrine files.
- Installer tests extended for two new addons (discovery, `--only`, `--no-addons`, markers).
- No placeholder text anywhere in shipped files — a design doctrine that ships
  `<TODO>` has no credibility.

## Attribution

hallmark is MIT, impeccable is Apache-2.0, this repo is Apache-2.0. **Independent
authorship is the rule**: study the prior art, then write our own rules in our own words. No
large blocks of their rule text. A `NOTICE.md` credits both projects as prior art studied.

## Open item carried forward

`DESIGN-SYSTEM.md` authorship remains unresolved: the winning author knows the intent and is
the worst party to describe what was actually built; a non-author lacks build context.
Current preference is the Phase-6 design reviewer. `PDS.md` MUST state the rule it adopts and
MUST NOT leave it implicit.

## Deferred (documented, not in v1)

`T3 PIXEL` detectors and visual regression · a TUI/CLI annex · empirical per-participant
prior measurement · a `LEDGER.md` the next design run reads and the §13 retro interrogates ·
any core-protocol carve-out in `COOPERATION.md` (a separate meta-protocol-change idea).
