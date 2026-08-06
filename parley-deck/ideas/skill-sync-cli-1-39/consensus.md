---
idea: skill-sync-cli-1-39
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
rounds: 2
date: 2026-08-06
status: consensus
---

## What was decided

Six edits to `parley-deck-skill`. The brief asked about three; the run found **four defects in one
forty-line region of `SKILL.md`**, and three of the four were found by participants rather than by
the facilitator's brief.

The framing that settled the scope, proposed by codex-1 and applied by everyone including against
their own round-1 positions: **the skill is instruction text, so the test is not "is this true" but
"does a facilitator or participant behave differently because of it".**

### 1. `opencode` row in the Autonomous Execution table

Unanimous. Adopted in kimi-1's wording, which carries the argv detail the other versions omit:

> | opencode | `run --auto` — the prompt is an argv positional, not stdin. `opencode run` writes
> unattended even without `--auto`; pass `--auto` explicitly, because an implicit vendor default is
> what may change between versions. |

### 2. `SKILL.md:251` — the inverted source-of-truth sentence

**kimi-1's find, and the sharpest of the run.** The line *"The source of truth for each agent's mode
is the spec's `autonomous_write` field"* points readers at the field 1.39.0 established does **not**
launch anything. It sits three lines under the table a facilitator uses to pick a flag, so it
actively directs attention away from the value that decides whether the agent can write.

All three non-facilitator participants converged on one replacement, each building on the others.
Adopted text — kimi-1's, which is the superset (its lead sentence is the inversion fix, its body is
codex-1's contract framing, its mechanism clause is hermes-1's, and the existing confinement floor
is retained verbatim at the end):

> The source of truth for an agent's autonomous capability is the **effective launch argv**, not the
> declared mode. The declared autonomous-write mode is a verification contract, not a second set of
> launch arguments: before treating a headless participant as able to write its artifact, inspect
> the effective launch arguments after all configuration layers have been applied — the launch
> config recorded in the orchestration summary, or `parley agents list` when the parley CLI drives
> the agents — and verify that every argument required by the declared mode is present. A config
> override can replace the launch arguments wholesale and silently drop the enabling flag. If the
> effective arguments cannot be inspected, or any required argument is absent, treat autonomous
> write as unavailable (`AUTO=no`) and do not launch that participant as write-capable. Passing this
> check proves only that the autonomous mode is enabled; it does not prove workspace confinement. If
> workspace confinement cannot be demonstrated for an agent, treat its autonomous bit as unset
> (fail-closed) rather than escalating to a full-filesystem bypass.

This replaces **both** sentences at `:251` — the inverted one and the confinement-only fail-closed
one — rather than adding a third paragraph beside them.

### 3. The command-construction recipe — manual vs CLI branches

**codex-1's find.** Steps 4-5 tell a facilitator to add model/thinking flags and deliver the prompt
as separate assembly stages. The launcher does neither: it substitutes `{prompt}` and `{root}`
*inside* `HeadlessArgs` and launches that vector alone (`runner.go:1094-1108`). The recipe is a
legitimate description of hand-rolling an invocation and a false description of what `parley` does
with config, and the section currently blurs them into one numbered list.

Adopted: split "Generic CLI Invocation Contract" into an explicit **manual facilitator** branch and
a **Parley CLI** branch. The multi-step assembly stays in the manual branch only. The CLI branch
states that resolved `headless_args` is the complete argv template, `{prompt}` must already sit in
its required position, `prompt_mode=stdin` controls only stdin wiring, and **no** permission, model,
thinking, profile or prompt arguments are synthesized afterwards. The same boundary applies to
"Headless Agent Configuration" and `WORKED_EXAMPLES.md`.

### 4. `compatibility.json` — bump plus one assertion

**Unanimous only after hermes-1 withdrew its round-1 objection, and it withdrew on evidence rather
than on the count.** It checked whether the existing `parley-addon.json` sha256 pin already guards
this, found that the pin verifies file hashes and never reads `skillVersion` or compares it to
`package.json`, and concluded its "no new tooling" objection did not apply.

Adopted: set `skillVersion` to **whatever version actually ships from this change** (not a number
chosen now), and add exactly one equality assertion — `compatibility.skillVersion ===
package.version` — to the existing Node test harness, running under both `npm test` and the existing
`prepack` lifecycle. No new script, job, or checker.

Recorded cost, raised by claude-1 and not disputed: once the assertion lands, bumping `package.json`
without `compatibility.json` fails the suite and `prepack`. That is the intent, and it becomes a
release step rather than a surprise.

hermes-1 additionally flagged an implementation gap: the existing manifest test verifies hashes and
does not read either version field, so the assertion needs its own home rather than an extension of
that check.

### 5. `writeModeArgs` — see `## Verdict conflicts`

### 6. What stays out

Unanimous. Promotion history, probe versions and outputs, the `hermes` incident narrative, ACP
availability, kimi's exit-1 message, and every CLI struct name. None changes a next action.

`references/COOPERATION.md` is **untouched**: 1.39.0 changed no protocol text, verified by
`git diff v1.38.0..v1.39.0` over both copies and by a normalized diff against the live deck.

## Verdict conflicts

### VC-1 — Does `writeModeArgs` stay in the documented JSON shape? OPEN

Everyone agrees the CLI-facing contract must show no such field. The residue is the **manual**
facilitator JSON example.

- **hermes-1, kimi-1, claude-1 — delete it.** hermes-1: *"A facilitator copying the JSON example
  will still populate `writeModeArgs` and still believe it matters. Remove it from the JSON shape;
  add the invariant line. Both, not one."* kimi-1 **withdrew its round-1 "keep the split"** on the
  grounds that the field split itself teaches the launch model the CLI does not implement.
- **codex-1 — keep it, but only inside the explicitly manual branch**, on the grounds that once the
  manual and CLI branches are separated the field is no longer a claim about CLI config, and that
  keeping it avoids a schema bump for local configs that already set it.

**This is 3-to-1 and §15.3 forbids resolving it that way, so it is carried rather than closed.**
The unexamined fact on which it turns: **nobody checked whether any local config in the wild
actually sets `writeModeArgs`.** claude-1 raised that in round 2 and it went unanswered. If none
does, codex-1's compatibility argument loses its object and deletion is free; if one does, deletion
is a silent behaviour change for that user.

**Resolution path for signoff:** whoever can, check it. Otherwise the drafter's proposal is to
delete the field from the JSON shape and state in the manual branch that the write-enabling flag
belongs inside `headless_args`, with codex-1's dissent recorded in `FINAL.md` rather than smoothed
away.

## Drafter position changes

claude-1 is facilitator, participant and drafter. Required by §15.5.

| # | Prior position | Source | New position | Why |
|---|---|---|---|---|
| 1 | *"I have proposed naming `hermes` as the live case… If two participants think the clause should go, it goes."* | `round-01/claude-1.md` | Withdrawn | codex-1 and kimi-1 both ruled it out on the does-it-change-an-action test. The drafter had applied that test to exclude probe versions and promotion history, then made an exception for the one incident it had personally worked on |
| 2 | Item 4 framed the defect as `writeModeArgs` being a dead field | `round-01/claude-1.md` | The defect is the whole assembly recipe plus the inverted `:251` sentence | codex-1 and kimi-1. The original framing would have produced a patch fixing one symptom of three |
| 3 | Proposed a boundary marker sentence for the recipe | `round-02/claude-1.md` | codex-1's fuller manual/CLI branch split | codex-1's version also covers "Headless Agent Configuration" and `WORKED_EXAMPLES.md`, which a single marker sentence does not reach |
| 4 | Proposed the drafter's own `:251` replacement wording | `round-02/claude-1.md` | kimi-1's converged text | It is a superset; the drafter's merged the CLI surface and the incident into one sentence, which hermes-1 correctly identified as the weaker structure |

Four changes, all four forced by another participant.

## Comparison & blind spots

**Correlated-agreement caveat (§15.6).** The four participants are related models. The unanimous
`CONFIRMED` on the launch-path claim is a shared prior, not four independent confirmations —
mitigated here by the fact that three participants read `runner.go` themselves and quoted different
line ranges, and one (hermes-1) proved the negative the facilitator had only asserted.

**§15.6's trigger did not fire.** Round 1 was not unanimous — hermes-1 opposed the guard — so no
steelman round was required. It became unanimous only after an evidence-based withdrawal in round 2,
which is the mechanism working rather than convergence by anchoring.

**What would have to be true for this to be wrong.** The whole package rests on one claim: that
`HeadlessArgs` is the sole launch input. If any code path appends arguments elsewhere — a wrapper,
a transport-specific launcher, a future ACP path — then the "effective argv" instruction sends
facilitators to an incomplete value. Three participants greped for it independently and none found
one, but all three searched the same repository with the same tool.

**Where nominally independent findings are one family.** codex-1's recipe finding and kimi-1's
`:251` finding are the same defect at two altitudes: the skill treats the declared mode and the
launched command as independent values. The facilitator's `writeModeArgs` item is the third
instance. They were filed separately and should be fixed as one passage, not three patches.

## Signoffs

<!-- Each participant appends its own block. Do not edit another participant's block. -->
