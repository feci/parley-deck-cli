---
idea: protocol-generation-bias
author: claude-1
created: 2026-08-28
track: deliberation
participants: [claude-1, codex-1, hermes-1, kimi-1, zcode-1]
excluded: [opencode-1]   # §9.0 — see inbox/claude-1-to-all_protocol-generation-bias_opencode-excluded.md
roles:
  claude-1: A4-adversarial-appointment
  codex-1: A5-anchor-hygiene
  hermes-1: A6-subtract-nothing-new
  kimi-1: A1-forced-divergence
  opencode-1: A2-reframe-vocabulary
  zcode-1: A3-gate-trigger-repair
require_model_diversity: true
status: final
---

## Problem / idea

The protocol has a deep, well-evidenced theory of **verification** bias — no self-verdicts,
provenance-capped verdicts, conflicts never resolved by count even when unanimous. It has almost
nothing for **generation** bias: nothing forces a structurally different candidate to exist
before convergence, and nothing ever asks whether a simpler solution exists at all.

An external critic put it this way (paraphrased from a real conversation, 2026-08-28):

> When one agent comes up with a proposal, the others mostly just amend *his* proposal — they
> don't come up with a completely different one. AI cannot do the thing you actually need:
> judge whether a better and simpler solution exists.

His evidence: agents converged on hand-written build scripts copying dependencies out of
`node_modules`; the obvious `pnpm deploy` route only appeared after a human said it out loud.

**Design this away, or prove it cannot be designed away.**

### Findings established before this idea opened

These are measurements, not a proposal. They are stated so nobody re-derives them; they are
**not** a solution to react to. Verify any of them you intend to rely on — §15.2 applies, and
anything you take from this list unverified is `RECALL`, not `PRIMARY`.

1. **§15.6 is the only forced-alternative rule in the protocol**, and its trigger excludes this
   entire class of task: it fires only when the idea's output is *"primarily a judgment rather
   than a mechanically decidable artifact"* (`COOPERATION.md:1341`). A build/bundling fix is
   mechanically decidable, so the one anti-convergence rule was off for the critic's case. Its
   second condition — round 1 closing with *no substantive disagreement* — also switches it off
   when participants argue about details **within** one frame.
2. **No simplicity lens exists.** Whole-protocol counts: `simpler` 0, `simplicity` 0, `YAGNI` 0,
   `over-engineer` 0, `smallest` 0, `off-the-shelf` 0, `built-in` 0, `do nothing` 0.
   `alternative` appears 3 times, all inside §15.6.
3. **The review vocabulary is closed** to `CRITICAL|MAJOR|MINOR|NIT`
   (`internal/driver/impl.go:445`). Grep for `ALTERNATIVE|APPROACH|SIMPLER|REDESIGN` across all
   Go outside `defaults/`: zero hits. A "the whole approach is wrong" finding is not forbidden —
   it simply has no category, no route, and no destination, because `FINAL.md` is frozen once
   Phase 5 starts (a reframe requires a new `<slug>-v2` idea).
4. **§15 has no enforcement surface at all.** `/usr/bin/grep -rl` over `internal/` + `cmd/` for
   `PRIMARY|SECONDARY|RECALL|DISPUTED|Verdict conflicts|Drafter position changes|Adversarial
   alternative|correlated|blind spot` matches **zero** `.go` files.
5. **The installed skill carries none of it.** `skills/parley-deck/SKILL.md` has 0 occurrences of
   `§15`, `PRIMARY`, `RECALL`, `steelman`, `correlated`, `Refutation`. Its review template
   (`SKILL.md:768-786`) omits `## Refutation attempts` — the exact section the Go gate rejects a
   review for lacking (`internal/protocol/reviewartifact.go:41`). The documented hand-driven path
   emits artifacts the enforced path refuses.
6. **Opt-in gates are not used.** Across 88 ideas in this deck: `track:` declared 32,
   `checks:` 4, `strict_gate:` 3, `auto_implement:` 1, **`require_model_diversity:` 0**.
   Separately, the rule that later-round reviewers must address every other reviewer runs at
   **7% compliance across 348 reviewer files** — every one of those ideas closed green.
7. **`parley-design` (PDS/1.0) already solves the generation half, enforced in JS, for visual
   design only**: proposers are *assigned* distinct positions on the brief's primary axis by a
   seeded rotation (`PDS.md:367`), and gate **G1 DISTINCTNESS** *"MUST fail if any pair of
   directions differs on fewer than two of the brief's axes… Persistent convergence never
   auto-passes"*. None of this is in `COOPERATION.md`.
8. **Corpus measurement.** Where the final mechanism did not exist in round 1, a human put it
   there in 8 of 10 confirmed frame reversals. Three agent-originated frame breaks did land, and
   all three were subtractive or epistemic ("delete the invalid rule", "require a witness"),
   never "here is a different machine". 1,389 later-round files carry a `## Position changes`
   section and only ~40 declare "none" — the protocol produces the *appearance* of movement at
   very high volume.

### The two benchmarks your proposal must be judged against

**B1 — the missing option.** `servers/…/ideas/2026-08-14T12-41-49-daily-backup-str`. In
round 2, `claude-1` wrote: *"Proxmox Backup Server has a native S3-compatible datastore
backend… **Nobody proposed the option that actually exists.** I am withdrawing my own round-1
proposal."* `FINAL.md:18` shipped the round-01 design (`vzdump` + `rclone`) **anyway**. An agent
broke the frame and the frame won. Any mechanism that only *generates* alternatives fails B1;
B1 is about what happens to one after it exists.

**B2 — the unproposed option.** The critic's `pnpm deploy` case: the cheap, documented,
off-the-shelf route was never on the table until a human named it. Any mechanism that only
*adjudicates between candidates already on the table* fails B2.

**A proposal that does not say what it does on B1 and on B2 is incomplete.**

## Constraints

- **You are assigned one solution axis** (see `roles:` above). Develop *that* axis as strongly
  as you can. The assignment is deterministic, not editorial: agents sorted lexicographically,
  axes rotated by `sha256("PDS/1|protocol-generation-bias")[0:8] mod 6 = 3`. Per §5 the `roles:`
  map is advisory — it does **not** change quorum, ownership, signoff weight or drafter
  eligibility. It constrains only what you are asked to *generate* in round 1.
  - `A1-forced-divergence` — make N structurally distinct candidates *exist* before any
    convergence is permitted (the PDS port and anything better).
  - `A2-reframe-vocabulary` — give "this whole approach is wrong, here is another" a finding
    class, a route, and a destination that Phase 5–8 can actually absorb.
  - `A3-gate-trigger-repair` — repair §15.6's trigger and the 0–5% adoption problem; decide what
    must become default-on and what that costs.
  - `A4-adversarial-appointment` — appoint the frame-breaking work to someone: a standing
    red-team, a rotating simplicity auditor, a devil's advocate with defined success criteria.
  - `A5-anchor-hygiene` — attack the anchor at its source: what may and may not enter
    `00-prompt.md`, blinding, presentation order, staged disclosure, isolated staging.
  - `A6-subtract-nothing-new` — **argue the null position honestly**: that the fix is deleting
    rules rather than adding them, that no new gate would have fired on B1 or B2, and that a
    seventh unused opt-in flag is the fourth instance of this deck's own defect class. This is a
    real position with real evidence behind it, not a strawman. If you conclude a mechanism is
    needed, say so — but make the protocol earn it.
- **Bring external evidence.** The user explicitly asked for published work and approaches from
  outside this repository. Cite real papers, systems, or documented practices with enough
  locator that another participant can find them — multi-agent debate and its failure modes,
  conformity and sycophancy measurements, LLM-as-judge self-preference, diversity/ensembling,
  red-teaming protocols, structured analytic techniques from intelligence analysis (ACH,
  devil's advocacy, team A/team B), design-thinking divergence, Nominal Group Technique, or
  anything else that actually bears. **Tag every claim per §15.2.** Do not invent citations: an
  unverifiable reference is worse than none, and will be treated as a defect, not a contribution.
- **State where your own axis fails**, and name the axis you would defect to if forced. A file
  that only advocates is a weaker contribution than one that also bounds itself.
- **Cost matters.** `ideas/mas-research-mining/FINAL.md` ratified that shared-rule-text bytes
  must be net negative. A proposal that adds ceremony to every idea to catch a rare failure must
  say what it deletes, or defend the addition explicitly.
- Repository is **read-only** to you for this round. Do not modify code, protocol text, or any
  file other than your own round-01 artifact.

## Non-goals

- Not deciding today whether to implement. Round 1 is generation, not selection.
- Not re-litigating §15's verification half — it works and is not the subject.
- Not a `parley-design`/PDS rewrite. Porting a mechanism out of PDS is in scope; changing PDS is
  not.
- Not a defence of the protocol against the critic. If he is right, the correct output is a
  design that fixes it, or a documented reason it cannot be fixed.

## A warning about this brief

This document is itself the anchor risk it describes. All six of you read it at the same
instant, which is exactly the condition `protocol-overlay-local-extension` recorded as *"one
analysis with four signatures, not four confirmations"*. It deliberately contains **no candidate
solution** — the findings above are measurements of the gap, not a design. If you notice
yourself agreeing with its framing, say so explicitly in `## Concerns / open questions`, and
attack the framing rather than the gap.
