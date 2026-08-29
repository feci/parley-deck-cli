---
agent: opencode-1
idea: protocol-generation-bias
round: 1
date: 2026-08-28
---

## Summary

A2 gives "this whole approach is wrong, here is another" the three things the closed review vocabulary does not: a **class** (`REFRAME`), a **route** (absorb before `FINAL.md` freezes), and a **destination** (`## Frames considered` inside `FINAL.md`, which Phases 5–8 can read). Generation of the other machine is not this axis; absorption is. B1 is an absorption failure. B2 is a generation failure A2 cannot fix alone.

## Proposed approach

### The gap this axis owns

Review severity is a closed set of four in-frame defect grades
(`internal/driver/impl.go:445`, PRIMARY: `case "CRITICAL", "MAJOR", "MINOR", "NIT":`).
A claim that the *frame* is wrong is not a more-severe `CRITICAL`. It has no tag,
so no gate, so no paragraph in a frozen `FINAL.md`. The implicit destination today
is a successor idea `<slug>-v2` (RECALL: brief finding 3). That destination sits
*after* Phase 5 has already shipped the old frame.

B1 is the existence proof that homelessness, not generation, is the failure.
In round 2, `claude-1` named Proxmox Backup Server's native S3-compatible
datastore and withdrew the round-1 design. `FINAL.md` of
`servers/parley-deck/ideas/2026-08-14T12-41-49-daily-backup-str` does not mention
PBS, Proxmox Backup, or native S3 (PRIMARY: `rg -n "PBS|Proxmox Backup|native S3"`
on that file → no matches). The other machine existed. It was never carried.
Finding 3 of the brief is therefore not "reframes are forbidden". It is
"reframes are homeless". Homeless findings lose to freeze.

Do not quote a protocol sentence "Severity tags are fixed" — that phrase is not
in `COOPERATION.md` (PRIMARY: `rg "Severity tags are fixed"` → 0). The closure
is in the driver, not in shared rule prose. A2 should put the class in the
prose *and* keep it *out* of that `case` switch.

### Class: `REFRAME` is not a fifth severity

Do **not** extend `impl.go:445`. Grading "wrong approach" as worse or better
than a bug is a type error. `REFRAME` is a parallel finding class, the way
provenance tags sit beside severity rather than inside it (RECALL: §15.2).

Closed payload, four fields, all required:

1. **Current frame** — one sentence naming the approach under discussion.
2. **Other frame** — one sentence naming a structurally different approach
   (different mechanism, different vendor, or "use the built-in / do nothing"),
   not a patch to (1).
3. **Witness** — a locator another participant can open without the author:
   man page, vendor doc URL, shipped subcommand, RFC, or an in-repo path.
   No witness → it is not a `REFRAME`; it is a `MAJOR` or an argument.
4. **Stay-condition** — what would have to be true for (1) to remain the
   design after (2) is on the table.

That payload is the whole rule. No simplicity essay. The witness *is* the
off-the-shelf check the protocol currently has zero words for (RECALL: brief
finding 2; whole-protocol counts of `simpler` / `YAGNI` / `off-the-shelf` = 0).

External analogs (all RECALL; locators are real, not invented):

- Michael Nygard, "Documenting Architecture Decisions" (2011),
  https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions —
  ADRs record options considered and why rejected ones lost. Destination
  pattern, not generation pattern. `## Frames considered` is an ADR stanza
  inside `FINAL.md`.
- Richards J. Heuer Jr., *Psychology of Intelligence Analysis*, CIA, 1999,
  Analysis of Competing Hypotheses. A hypothesis never entered on the matrix
  cannot win or lose. That is B1: PBS was spoken and then omitted from the
  matrix that shipped.
- Irving L. Janis, *Groupthink*, Houghton Mifflin, 1972 — "mindguards" keep
  dissenting frames off the written estimate. Freeze-without-absorption is
  the protocol's mindguard.
- Perez et al., "Discovering Language Model Behaviors with Model-Written
  Evaluations", arXiv:2212.09251 — sycophancy toward the user's (here: the
  first agent's, the brief's) frame.
- Liang et al., "Encouraging Divergent Thinking in Large Language Models
  through Multi-Agent Debate", arXiv:2305.19118 — debate elaborates a shared
  frame more often than it replaces it. A vocabulary that only grades patches
  will keep producing patches.
- Panickssery, Bowman, Perez, "LLM Evaluators Recognize and Favor Their Own
  Generations", arXiv:2404.13076 — self-preference. The drafter of `FINAL.md`
  is an evaluator of the round-1 frame; without a forced stanza, the native
  S3 option is dropped rather than scored.

### Route: absorption before freeze

Rounds 1–4: a `REFRAME` may appear in any participant file. It is not
resolved by count (RECALL: §15 verification half; out of scope to relitigate).

Phase 4 → 5 gate, one sentence of shared rule text:

> `FINAL.md` MUST NOT freeze while any `REFRAME` raised in this idea is
> absent from `## Frames considered`. Each entry MUST be `ADOPTED`,
> `REJECTED` (citing the stay-condition and the witness), or `SUPERSEDED`
> (citing the successor frame). Silence is not rejection.

That is the B1 fix. `claude-1`'s PBS withdrawal becomes a `REFRAME` whose
witness is the PBS datastore documentation. Freeze is illegal until
`FINAL.md` either adopts native S3 or records why `vzdump`+`rclone` still
meets the stay-condition against that witness. Shipping round-01 anyway
is then a gate failure, not a drafting preference.

Phase 5–8 destination:

- `IMPLEMENTATION.md` MAY NOT introduce a frame absent from
  `## Frames considered`. Implementation drift *is* a `REFRAME`, and it
  reopens only that section of `FINAL.md` — not a `<slug>-v2` idea.
- Phase 7 MAY emit `REFRAME` only for drift. "I have a better idea now"
  after freeze is a new idea, not a finding. This stops A2 becoming an
  infinite reopen.
- In-frame defects stay `CRITICAL|MAJOR|MINOR|NIT`. They do not change class.

No new opt-in flag. The class is default-on for every idea that produces
a `FINAL.md`. B1 was a mechanically decidable backup design; §15.6's
trigger would still have been off (RECALL: brief finding 1,
`COOPERATION.md:1341`). A2 must not live behind that trigger.

### What A2 does on the benchmarks

**B1 — the missing option that existed.** A2 is built for this. The
failure was destination, not generation. Round 2 produced the other
machine; `FINAL.md` dropped it (PRIMARY: no PBS string in that file).
With `REFRAME` + freeze gate + `## Frames considered`, native S3 is
either the shipped design or a recorded rejection against a witness.
"Shipped anyway, unmentioned" is ungrammatical. A2 does not require
PBS in round 1; it requires that once named, it cannot vanish.

**B2 — the unproposed option.** A2 does not generate `pnpm deploy`. A
finding class cannot summon a candidate nobody wrote. If a human or
another axis names it before freeze, A2 gives it the same landing path
as PBS. If nobody names it, A2 is silent and B2 still fails. That is
the bound of the axis, not a defect in the class.

### Cost

Ratified constraint (RECALL: `ideas/mas-research-mining/FINAL.md`):
shared-rule-text bytes must be net negative. A2 adds ~15 lines (payload
+ one freeze sentence + one Phase 5–8 paragraph).

What it deletes, so the addition is earned:

1. **Delete the `<slug>-v2` requirement for pre-close frame breaks.**
   Today's destination for "wrong approach" is a new idea — a full
   lifecycle of ceremony to park a finding that already exists. One
   in-band stanza is smaller than a successor idea.
2. **Do not add a fifth severity, a simplicity essay, or a seventh
   opt-in flag.** Brief finding 6 (RECALL): `require_model_diversity`
   is 0/88. Flags that are not used are this deck's own defect class.
   A2 lives inside existing reviews.
3. **Do not port PDS G1 into `COOPERATION.md`.** Distinctness-before-
   convergence is A1. A2 only absorbs.

If the byte budget still fails, delete the unused second trigger of
§15.6 rather than adding a flag to turn A2 off. That deletion is A3's
text; A2 only needs it as an offset.

## Concerns / open questions

This brief is the anchor it warns about. Six agents received the same
gap list at the same instant. Treating "generation bias" as the disease
and "a mechanism" as the cure is the overlay-local-extension failure
mode (RECALL: "one analysis with four signatures"). Attacks on the
frame, not only on the gap:

- The critic's `pnpm deploy` story may be a prompting failure (the first
  proposal entered the task as if it *were* the task) more than a
  vocabulary failure. That is A5's ground. A2 is still required for B1
  even if A5 is right about B2.
- "Give it a class" can become theater: every `MINOR` restated as
  `REFRAME` to stall freeze. The witness field is the anti-theater
  rule. No locator, no class.
- Assignment is advisory (§5) and still an anchor. The strongest A2
  is one that states it cannot satisfy B2.

Open: who writes `## Frames considered`? The drafter, with the REFRAME
author's four fields copied verbatim. Paraphrase is how PBS disappeared.

Open: may `REJECTED` cite cost or schedule alone? No. Only the
stay-condition plus witness. Cost-as-sole-reason is how a named better
machine loses to inertia.

## Risks

**Where A2 fails.** It does not make a structurally different candidate
exist. The corpus claim that agent-originated frame breaks were
subtractive, never "a different machine" (RECALL: brief finding 8), is
exactly what A2 cannot reverse: you cannot absorb a machine that was
never proposed. A2 also fails if the drafter writes `REJECTED`
boilerplate; the gate checks presence and shape, not honesty. A2 fails
late on purpose: a non-drift REFRAME in Phase 7 is ignored, and a true
late discovery then needs a new idea — same as today.

**Defect.** If forced to drop A2, defect to **A1-forced-divergence**.
B2 is the harder benchmark; only a generate-before-converge rule
addresses it. A3 without A2 still has nowhere to put the alternative
once the repaired trigger fires. A4 can appoint a frame-breaker who
then files a homeless finding. A6 (delete rules) does not give PBS a
paragraph in `FINAL.md`. A5 may blunt the first-proposal anchor but
does not route a late break. Keep A2 even if A1 is also adopted:
generation without a landing pad is B1.

**Gaming / cost.** Witness-free spam, or 15 lines that become 150 in
committee. Keep the payload closed at four fields. Do not put
"simplicity" in shared rule text; the word is a magnet for essays this
deck already measured at zero and should not start.

**Enforcement.** §15 currently has no Go surface (RECALL: brief finding
4). A2's freeze gate is one string check: `FINAL.md` contains
`## Frames considered` and every prior `REFRAME` id is listed. Smaller
than enforcing provenance tags. The skill template must gain that
heading or the documented path will again emit artifacts the gate
refuses (RECALL: brief finding 5, `SKILL.md` vs
`internal/protocol/reviewartifact.go:41`). Adding the heading repairs
the existing emit/enforce split; it is not a new ceremony. Offset: do
not add a `## Refutation attempts` essay to this class. `REFRAME`
already carries its other-frame.
