---
agent: kimi-1
idea: skill-sync-cli-1-39
round: 2
date: 2026-08-06
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

**Scope checked (§15 declaration).** Read in full this round: `00-prompt.md` and all four
round-01 files (`claude-1.md`, `codex-1.md`, `hermes-1.md`, `kimi-1.md`). Executed fresh
primary checks in the skill repo: a repo-wide grep for `writeModeArgs|write_mode_args|
WriteModeArgs` (three documentation hits, zero code consumers), a grep of `test/` for
`compatibility.json`/`skillVersion` consumers, and a read of `test/manifest-coverage.test.js`
(header and structure). CLI-side inspections from round 1 are not re-litigated — the launch
model is settled and I do not re-argue it. Not done this round: running the skill test suite,
probing live binaries, re-running the COOPERATION.md currency diff (facilitator-owned; no
position here depends on it). Verdicts below are tagged and none is issued on a claim I own.

## Position changes since round 1

1. **`writeModeArgs`: I withdraw "keep the split."** Round-1 me argued to keep the two-list
   JSON shape and add only an invariant line, on the grounds that deleting a documented field
   churns the shape without adding the guard. That was wrong on the argument, for three
   reasons. (a) The settled PRIMARY finding is not merely that a check was missing — it is
   that the field split *itself teaches the launch model the CLI does not implement*; an
   invariant sentence bolted onto a wrong model leaves the wrong model standing.
   (b) The churn objection collapsed under evidence: my fresh repo-wide grep finds
   `writeModeArgs` only at `SKILL.md:363`, `SKILL.md:811`, and `WORKED_EXAMPLES.md:33` —
   nothing in `bin/`, `lib/`, `scripts/`, or `test/` reads it, so removal breaks no consumer
   in this repo. (c) A field nothing reads, which a human must remember to also mirror into
   `headlessArgs`, is the defect's structure under a doc-string name. I now converge with
   claude-1, codex-1, and hermes-1: **remove `writeModeArgs` from the documented shape** —
   conditioned on the boundary marking in item 2 below, which is what makes removal safe
   rather than merely cosmetic.
2. **The multi-step recipe (round-2 question 2): ruling — legitimate manual-assembly
   abstraction that needs a boundary marking, not a rewrite.** The steps are instructions for
   a facilitator assembling a command by hand; they become false only when read as a model of
   what the launcher does. The fix is to state the boundary: the launcher substitutes
   `{prompt}` inside the effective headless vector and synthesizes nothing — so
   model/thinking/profile flags, exactly like the write-enabling flag, reach the process only
   if the facilitator put them in the assembled argv. kimi (`[-p, {prompt}]`) and opencode
   (`[run, --auto, {prompt}]`) are the concrete proof. That is two sentences, not a rewrite;
   a rewrite would re-narrate a manual procedure that is not itself wrong.

## Responses to others

### @hermes-1

Your guard objection has three limbs. On the argument, each fails.

- **"The skill already has a drift-detection mechanism."** Verdict on that claim:
  **REFUTED [PRIMARY]** for this field. In round 1 I executed the comparison myself:
  `parley-addon.json` pins a sha256 of `references/compatibility.json`, and the pinned hash
  **matches the stale content** (`b2465b20…67abd`, round-1 file item D). The mechanism
  verifies byte integrity against a manifest that was re-pinned over the drifted bytes.
  Whatever the Protocol Drift Check covers, four releases of `skillVersion` drift passed
  through it unnoticed — that is not a mechanism that needs no backup; it is the measured
  proof that the existing one does not detect this drift class.
- **"`skillVersion` is informational; the policy says warn on drift, no lockstep."** The
  policy sentence governs *cross-component* lockstep (CLI version vs skill version). The
  proposed assertion is internal to one package: `compatibility.json.skillVersion` and
  `package.json.version` are two fields denoting the *same* fact, and hand-maintained
  duplicates drift — this one did, four times. If you hold the field is purely informational
  and nothing acts on it, the consistent positions are guard it or delete it; keeping a
  published-schema field (`schemaVersion: 1`) permanently wrong is strictly worse than both.
  Deletion was considered and rejected in round 1 because external tooling may read the file.
- **"Process discipline; the owner's standing rule now enforces it."** The standing rule says
  a CLI release ships the skill in the same turn. It does not name `skillVersion`, the field
  drifted under the previous process-discipline regime, and "remember to bump it" is the
  exact mechanism that produced four silent releases. The assertion converts the standing
  rule into a mechanical fact for the skill's own version field — one line in an existing
  `node --test` suite that already runs on `npm test` and `prepack`. No new command, script,
  dependency, or CI job: not new tooling by any reading of the constraint.

Your risk 1 (downstream consumers of `writeModeArgs`) is now answered with primary evidence:
the grep above — three doc hits, zero code. External consumers cannot be ruled out, but the
CLI never machine-read the field either, so the residual risk is documentation-only. Your
schema-bump question: no bump — `schemaVersion` belongs to `compatibility.json`'s own schema,
not to the SKILL.md worksheet, and the field was never validated.

On wording, two substantive disagreements stand. First, your fail-closed sentence merges the
two checks into one conditional ("If the enabling flag is absent …, *or if* workspace
confinement cannot be demonstrated …"). 1.39.0 deliberately separated declaration/effective
enabling from confinement; a fused conditional re-teaches one Boolean eligibility bit and
leaves the facilitator unable to say *which* property failed — codex-1's risk, and it is
real. Second, SKILL.md:251's first sentence is not "correct but incomplete": as a validation
pointer it is inverted. "The source of truth … is the spec's `autonomous_write` field"
directs the reader to the declaration as the authority to consult, and 1.39.0 made the
declaration the thing *under audit*. Your proposed text leaves that sentence standing above
the new check, teaching the pre-1.39.0 model inside the fixed paragraph.

Two smaller points. Your opencode row ends "ACP remains available via `opencode acp`" — drop
it: the skill's contract is one-shot headless invocation and mentions ACP nowhere, and you
named no facilitator action the clause changes. And on naming `agents list`: the skill already
names CLI surfaces where they *are* the check (SKILL.md:257). A rule naming no tool teaches
no check; the wording-rot risk is handled by stating the invariant first and the surface
second, which the converged text below does.

### @codex-1

- **Recipe finding:** agreed, and I rule for your own first option — boundary marking, not
  rewrite (see Position changes, item 2). Your row already leans there ("distinguish a manual
  assembly abstraction from CLI config"); I am making it explicit so the drafter does not
  rewrite steps that are only unmarked, not wrong.
- **JSON-as-runtime-source row: partial disagreement.** Carry the labeling half — mark the
  JSON shape as a manual-facilitator worksheet, never presented as what
  `PARLEY_HEADLESS_AGENT_CONFIG` parses. Do **not** carry "show the CLI TOML form": that
  embeds the CLI's config surface in vendor-neutral instruction text and goes stale on the
  next CLI config change — the very failure class this idea exists to fix. The env-var
  mismatch (SKILL.md:142/341 documents a JSON file; `runtime.go:440` parses TOML) predates
  1.39.0 and changes no 1.39.0 action; I hold it belongs in a follow-up idea, as I flagged in
  round 1. Your own scope discipline cuts both ways: the skill is not a changelog, and the
  TOML form is not a 1.39.0 behavior either.
- **Fail-closed wording:** I adopt your text as the base. It keeps the two checks separate,
  and the clause "passing this check proves only that the autonomous mode is enabled; it does
  not prove workspace confinement" is exactly the property hermes-1's merged conditional
  loses. Two amendments, both argued elsewhere in this file: prepend the source-of-truth
  replacement sentence (your file does not address SKILL.md:251's first half), and name the
  inspection surface in one subordinate clause so the check is performable.
- **`prepack`:** agreed — it closes the direct `npm publish` bypass at zero cost, and it
  strengthens the guard case against hermes-1's objection rather than adding to it.

### @claude-1

- Your item-4 hypothesis got the independent verification you asked for — three separate
  source inspections, one settled finding. Nothing there to contest.
- **Incident naming: confirmed — it goes.** You set the condition yourself ("if two
  participants say it goes, it goes"); codex-1 and I said so in round 1, and hermes-1's file
  independently leaves the narrative out. What stays is the mechanism clause — "a config
  override can replace the launch arguments wholesale and silently drop the enabling flag" —
  because it tells the facilitator what to look for. The agent's name and the war story
  change no next action.
- Your release-day risk (a red suite nobody expected) is the guard working as intended, but
  you are right it must be announced: the release notes for this skill release name the new
  assertion in one line. That is documentation, not tooling.
- Your table-column suggestion: agree with your own lean — no. Five rows do not need a column
  that says the same thing five times; the invariant sentence carries it once.

## New concerns / questions

1. **`WORKED_EXAMPLES.md:33` must be in the diff — still unclaimed.** My round-1 finding E.2
   was not picked up by any other round-01 file: the worked example carries the same
   two-list shape (`"writeModeArgs": ["--workspace-write"]`) and is the file a facilitator
   copies from. Whatever lands on the SKILL.md JSON shape must be mirrored there in the same
   edit — including the secondary defect that `--workspace-write` as a single token matches
   no shipped adapter (codex's is two tokens, `--sandbox workspace-write`, discover.go:196).
   If we remove the field from SKILL.md and forget the example, we fix the contract and
   leave the template that reproduces the bug.
2. **"Inspect the effective argv" must be performable by a manual facilitator** (hermes-1's
   risk 2, fairly raised). The skill already requires recording the effective launch config
   in the orchestration summary; the converged wording names that record as the manual
   inspection surface and `agents list` as the CLI one. Without that, we ship a check nobody
   can run without the CLI installed.
3. **Question for the drafter: the bump target.** `skillVersion` becomes whatever version
   this release ships (likely `2.4.0`), not `2.3.0`. The assertion makes this self-enforcing
   from then on, but the one-time correction still has to land in *this* diff — a release
   that ships the guard with the stale value fails its own suite.

## Current proposal

Five edits, everything else stays out.

1. **Autonomous Execution table — add the opencode row**, no ACP clause:
   > | opencode | `run --auto` — the prompt is an argv positional, not stdin. `opencode run`
   > writes unattended even without `--auto`; pass `--auto` explicitly, because an implicit
   > vendor default is what may change between versions. |
2. **Replace the SKILL.md:251 paragraph** (both sentences — the inverted source-of-truth
   sentence and the confinement-only fail-closed sentence) with the converged text, built on
   codex-1's base with my source-of-truth sentence prepended and hermes-1's mechanism clause
   retained:
   > The source of truth for an agent's autonomous capability is the **effective launch
   > argv**, not the declared mode. The declared autonomous-write mode is a verification
   > contract, not a second set of launch arguments: before treating a headless participant
   > as able to write its artifact, inspect the effective launch arguments after all
   > configuration layers have been applied — the launch config recorded in the
   > orchestration summary, or `parley agents list` when the parley CLI drives the agents —
   > and verify that every argument required by the declared mode is present. A config
   > override can replace the launch arguments wholesale and silently drop the enabling
   > flag. If the effective arguments cannot be inspected, or any required argument is
   > absent, treat autonomous write as unavailable (`AUTO=no`) and do not launch that
   > participant as write-capable. Passing this check proves only that the autonomous mode
   > is enabled; it does not prove workspace confinement. If workspace confinement cannot
   > be demonstrated for an agent, treat its autonomous bit as unset (fail-closed) rather
   > than escalating to a full-filesystem bypass.

   On the two competing round-01 texts: keep codex-1's, drop hermes-1's. hers fuses the two
   failure modes into one conditional and leaves the inverted 251 sentence standing; codex-1's
   keeps the checks separate and carries the "proves only … not confinement" clause, which is
   the one property the fused version cannot express. One text, not both — shipping two
   adjacent wordings of the same rule is how the paragraph got contradictory in the first
   place.
3. **Generic CLI Invocation Contract — remove `writeModeArgs` and mark the boundary.**
   Delete the field from the JSON shape (SKILL.md:363) and step 3 from the construction list
   (SKILL.md:811); label the JSON shape as a manual-assembly worksheet, not CLI runtime
   config; and state the single-list invariant in place of the removed step:
   > `headlessArgs` is the complete argv template that will be launched. It must already
   > contain every flag the command needs — the non-interactive and autonomous-write
   > enabling flags, and any model/thinking/profile flags — plus exactly one `{prompt}`
   > placeholder where the prompt is delivered through argv. The launcher substitutes
   > `{prompt}` inside this vector and synthesizes nothing else. After any config change,
   > verify the effective argv before trusting the launch to write unattended.
4. **Mirror the shape fix in `references/WORKED_EXAMPLES.md:33`** — drop the `writeModeArgs`
   field and correct the example flag to a shape a shipped adapter actually accepts (see New
   concerns, item 1).
5. **`compatibility.json` — bump plus guard.** Set `skillVersion` to the version this release
   ships; add one assertion to the existing `node --test` suite
   (`skillVersion === package.json.version`), also executed from the existing `prepack` path;
   name the new assertion in one line of the release notes.

Explicitly out: the hermes incident name, ACP availability, the CLI TOML form, the
JSON/TOML env-var mismatch (follow-up idea), and `references/COOPERATION.md` — byte-identical,
untouched, per the constraint.
