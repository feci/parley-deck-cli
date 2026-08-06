---
idea: skill-sync-cli-1-39
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: standard
rounds: 2
signoff-revisions: 2
date: 2026-08-06
status: final
---

# FINAL — sync `parley-deck-skill` to CLI 1.39.0

## Problem

CLI 1.39.0 changed what determines whether a headless participant can write its own artifact:
`headless_args` is the launch, and `autonomous_write` became a declaration audited against it.
`parley-deck-skill` still teaches the pre-1.39.0 model. Three of the four defects were found by
participants rather than by the brief, and all four sit inside roughly forty lines of `SKILL.md`.

## Decisions

### D1 — add the `opencode` row to the Autonomous Execution table

Unanimous. Adopted in kimi-1's wording:

> | opencode | `run --auto` — the prompt is an argv positional, not stdin. `opencode run` writes
> unattended even without `--auto`; pass `--auto` explicitly, because an implicit vendor default is
> what may change between versions. |

### D2 — replace the inverted source-of-truth sentence at `SKILL.md:251`

The line *"The source of truth for each agent's mode is the spec's `autonomous_write` field"*
points readers at the field that launches nothing, three lines below the table a facilitator uses
to pick a flag. Adopted replacement (kimi-1's converged text; replaces **both** the inverted
sentence and the confinement-only fail-closed sentence, rather than adding a third paragraph):

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

### D3 — split the command-construction recipe into manual and CLI branches

The numbered recipe describes hand-rolling an invocation and is a false description of what
`parley` does with config: the launcher substitutes `{prompt}` and `{root}` *inside* `headless_args`
and launches that vector alone. Adopted: "Generic CLI Invocation Contract" gains an explicit
**manual facilitator** branch (which keeps the multi-step assembly) and a **Parley CLI** branch
stating that resolved `headless_args` is the complete argv template, that `{prompt}` must already
sit in its required position, that `prompt_mode=stdin` controls only stdin wiring, and that **no**
permission, model, thinking, profile or prompt arguments are synthesized afterwards. The same
boundary applies to "Headless Agent Configuration" and `WORKED_EXAMPLES.md`.

### D4 — `compatibility.json`: bump plus one assertion, and close the `prepack` half

`skillVersion` sat at `1.4.3` while the package reached `2.3.0` — four releases — because nothing
reads it. Adopted: set `skillVersion` to the version actually shipped by this change, and add
exactly one equality assertion `compatibility.skillVersion === package.version`. No new script,
job, or checker.

**Implementation MUST also close the `prepack` half.** `package.json:60` runs `node --test` under
`test`; `package.json:66` runs only `node scripts/build-addon-manifest.js --check` under `prepack`.
A test-file assertion therefore does not gate `npm publish` — the exact path that allowed the
drift. Close it by either extending the `prepack` command or putting the equality check inside
`scripts/build-addon-manifest.js`. Shipping only the test-file assertion does not satisfy D4.

Accepted cost: after this lands, bumping `package.json` without `compatibility.json` fails the
suite and the pack step. That is the intent; it becomes a release step.

### D5 — delete `writeModeArgs` from the documented JSON shape, with a migration rule

The field is removed from the documented config shape. The manual branch states that the
write-enabling flag belongs inside `headless_args`, and carries the migration rule:

> When an existing `headless-agents.local.json` contains `writeModeArgs`, merge its arguments into
> that agent's `headless_args` and remove the field.

### D6 — what stays out

Promotion history, probe versions and outputs, the `hermes` incident narrative, ACP availability,
kimi's exit-1 message, and every CLI struct name. None changes a next action.

`references/COOPERATION.md` is **untouched**: 1.39.0 changed no protocol text. Verified by
`git diff v1.38.0..v1.39.0` over both copies and by a normalized diff against the live deck — the
three copies differ only in deck-instance header and roster placeholders.

## Recorded dissent (§15.3)

**codex-1 on D5.** codex-1's position was to keep `writeModeArgs` inside the explicitly manual
branch, on the ground that a labelled manual branch removes the field's false claim about CLI
config. VC-1 closes against it on the design argument — a two-field shape teaches a two-list launch
model the CLI does not implement — not on the 3-to-1 count and not on the workspace measurement,
which codex-1 successfully argued was inadmissible as the deciding evidence. codex-1's migration
condition was adopted in full as part of D5. This dissent is recorded rather than smoothed away.

**kimi-1's R3 — the tooling-defect explanation.** kimi-1 verdicted the drafter's two `grep -r`
measurements `WRONG / not reproduced`, and codex-1 likewise could not reproduce them. The
disagreement was resolved during revision 2 rather than carried: the cause is the facilitator's
shell aliasing `grep` to `ugrep --ignore-files`, which honours `.gitignore`; the participants run
without that alias, so their non-reproduction is the control case. The claim is scoped to the
facilitator's shell only. **§15.3 dependency check: no adopted decision depends on it.** D1-D6 rest
on (a), which carries three independent non-owner `rg`-based confirmations, and on the design
argument.

## Verification record (§15.2)

| claim | status | who |
|---|---|---|
| No Go code declares `writeModeArgs`/`write_mode_args` or reads `meta/headless-agents.local.json` | `CONFIRMED` | codex-1 `PRIMARY`, hermes-1 `PRIMARY`, kimi-1 `PRIMARY` on the field half (co-owner on the file half) |
| 23 deck configs, 12 contain the key, 10 non-empty, 9 exclusive, 2 empty | `CONFIRMED` | drafter (owner) + **kimi-1 `PRIMARY`** non-owner reproduction incl. the `igm-app` spot-check |
| `prepack` does not run `node --test` | `CONFIRMED` | hermes-1 raised, drafter `PRIMARY`, kimi-1 `PRIMARY` against `package.json` |
| Nothing under `test/`/`scripts/` reads `skillVersion` | `CONFIRMED` | drafter + kimi-1 `PRIMARY` |
| The facilitator's `grep` under-reports via `--ignore-files` | `CONFIRMED`, scoped to the facilitator's shell | drafter `PRIMARY` (`grep --version`, `git check-ignore -v`); codex-1 and kimi-1 supply the control cases |

The drafter issued **no verdict on any claim it owns**.

## Follow-ups (not in scope here)

1. **The 9 exclusive decks.** Their write-enabling flags sit in a field that works only on the
   manual path. No deck is edited by this idea; D5's migration rule is the instruction. Since
   1.39.0 a `parley` launch surfaces the gap as `AUTO=no` naming the missing arguments rather than
   failing silently.
2. **`grep`/`rg` as evidence of absence.** In the facilitator's shell both honour ignore files, so
   a negative result is not proof. Future measurements offered as evidence must use `find`-based
   enumeration, `rg -uuu`, or `/usr/bin/grep`, and must name the tool.

## Implementation

Implementer: claude-1 (FINAL drafter; no other participant claimed it). Target: a single skill
release carrying D1-D5. Deviations go to `IMPLEMENTATION.md`.
