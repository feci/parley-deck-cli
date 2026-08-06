---
agent: claude-1
idea: skill-sync-cli-1-39
round: 1
date: 2026-08-06
---

## Summary

My position: **three of the four items belong in the skill, one does not, and the most important
change is not on the brief's list at all.**

The skill is instruction text for agents. The test for every candidate fact is not "is it true"
but "does a facilitator or participant behave differently because of it". By that test:

- **`opencode` row** — carry it. A facilitator reading the table to launch a participant will
  otherwise have nothing for an agent that is now in the roster.
- **Fail-closed launch wording** — carry it, and this is the load-bearing one. It changes a check
  a facilitator must actually perform before trusting a participant can write.
- **`compatibility.json`** — bump it, but the bump is not the finding. The finding is that it went
  stale by four releases and nothing noticed.
- **Kimi's `--auto` incompatibility** — already covered; the existing row states it. It only needs
  the evidence attached.

And the item the brief missed: **the skill's own launch guidance now contradicts the CLI.**

## Proposed approach

### 1. The `opencode` row — carry

Add to the Autonomous Execution table:

> | opencode | `run --auto`. Its default permission mode already writes in-workspace, but pass
> `--auto` explicitly: an implicit default is what a vendor may change between versions. |

The second sentence earns its place. Without it a reader sees a `--auto` documented by its own
vendor as *"dangerous!"* and reasonably drops it, which is exactly the state the CLI review just
found live for `hermes`.

### 2. Fail-closed wording — carry, and it must name the check

The existing sentence is about *confinement*: "if workspace confinement cannot be demonstrated,
treat its autonomous bit as unset". The new case is different in kind — the mode is declared, the
confinement question does not arise, and the flag simply never reaches the process. Proposed text
to follow the existing paragraph:

> **Fail-closed applies to the launch itself, not only to confinement.** A config layer that
> overrides `headless_args` replaces them wholesale and does **not** touch `autonomous_write`, so an
> agent can declare a mode whose enabling flag the launched command never passes. Since CLI 1.39.0
> the `AUTO` column reports `no` in that case and names the missing arguments, and `parley agents
> list` prints the **effective** argv rather than a built-in label. **Read that line, not the
> label, before trusting that a participant can write its own artifact.** This was live for
> `hermes`: an override had dropped `--yolo` while it still reported `AUTO=yes`.

The concrete instruction is the middle sentence. Everything else is the justification a reader
needs to not undo it.

### 3. `compatibility.json` — bump, and add the cheapest possible guard

Bumping `skillVersion` to the release version is obvious and not the interesting part. The
interesting part: it drifted **four releases** without anyone noticing, because no test reads it.

I propose the minimum guard, not a new tool: a test in the existing suite asserting
`compatibility.json.skillVersion === package.json.version`. That is a few lines in a suite that
already runs on every `npm test` and every `prepack`. It converts a silent drift into a failed
release.

I would **reject** anything larger — a version-sync script, a release checklist file, a CI job.
The constraint says no new tooling without justification, and an equality assertion in an existing
test file is not new tooling.

### 4. What the brief missed — the skill's launch guidance now contradicts the CLI

`SKILL.md` "Generic CLI Invocation Contract" tells the facilitator to build the command itself:

> 1. Start with the configured `cli`. 2. Add `headlessArgs`. 3. Add `writeModeArgs` needed for the
> agent to write exactly one protocol artifact. …

And the JSON config shape it documents has **`writeModeArgs`** as a separate field from
`headlessArgs`.

**The CLI has no such field.** `AutonomousWrite.Args` is *never consumed at launch* — I verified
this by grepping every consumer in the repo: it is read in exactly two places, both of which render
the `AUTO` column, plus tests. The actual launch uses `HeadlessArgs` alone. That is precisely why
the defect was possible: the declared args and the launched args are separate values that nothing
reconciled.

So the skill instructs agents to compose a command from a field that does not exist, and implies
the autonomous flag is added independently of `headlessArgs` — which is the mental model that
makes the `hermes` bug invisible. This is worse than a missing row: it is guidance that would
reproduce the defect.

Proposed fix: state that the autonomous flag must be **inside** `headless_args`, and that
`autonomous_write` is a declaration used for auditing, not a second argument list appended at
launch. I would keep `writeModeArgs` in the documented JSON shape only if it genuinely still works
in the local config; if it does not, it should go.

**This needs verification by another participant** — I have checked the CLI source, not the skill's
own config path, and I own the claim, so under §15.1 I cannot verify it myself.

## Concerns / open questions

1. **How much of the incident belongs in the skill?** I have proposed naming `hermes` as the live
   case. It is one clause and it stops a reader treating the rule as hypothetical. But the skill is
   not a changelog and I may be over-weighting a defect I spent the afternoon on. If two
   participants think the clause should go, it goes.

2. **Does the kimi row need its evidence?** The row already says `--yolo`/`--auto` are mutually
   exclusive with `-p`. I would add the exit-1 message as proof, because the next person to doubt
   it will otherwise re-derive it. Minor either way.

3. **Should the table gain a column for "flag must be in `headless_args`"?** That would make item 4
   structural rather than prose. I lean no — five rows do not need a column that says the same
   thing five times — but a participant who disagrees should say so.

4. **Version discipline.** `skillVersion` will be `2.4.0` if we release that, but the value must be
   whatever version actually ships. If the guard lands, the value is enforced rather than chosen.

## Risks

- **Changelog creep.** Every fact in the brief is true and interesting, and most of them do not
  change agent behaviour. The `AutonomousWrite.MissingFrom` internals are explicitly a non-goal and
  I would resist any wording that smuggles them in as background.
- **The fix for item 4 could be wrong.** My claim that `writeModeArgs` does not exist rests on
  reading the CLI, not on testing the skill's documented config path. If a local config really does
  honour it, my proposed rewrite would break working guidance.
- **A guard that fails on release day.** If the `compatibility.json` test lands and someone bumps
  `package.json` without it, `npm test` and `prepack` both fail. That is the intent — but it means
  the release procedure gains a step that must be in the release notes, or the next release stalls
  on a red suite nobody expected.
- **My own record here.** I have been corrected five times today on claims I stated more strongly
  than the evidence supported, twice on this exact CLI change. Reviewers should treat item 4 as a
  hypothesis with a locator, not as an established fact.
