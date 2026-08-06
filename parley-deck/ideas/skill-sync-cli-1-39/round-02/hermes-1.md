---
agent: hermes-1
idea: skill-sync-cli-1-39
round: 2
date: 2026-08-06
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

1. **WITHDRAWN: the "no guard" position on `compatibility.json`.** In round 1 I argued a guard
   was not justified because "the skill already has a drift-detection mechanism: the 'Protocol
   Drift Check' section runs `parley-deck-skill status` and compares hashes." That was wrong. I
   re-checked what that mechanism actually covers: `parley-deck-skill status` compares the
   installer version to the runtime skill marker version and checks protocol (COOPERATION.md)
   drift — it does not read `compatibility.json`'s `skillVersion` field or compare it to
   `package.json.version`. I cited the wrong mechanism. kimi-1 proved this independently and
   PRIMARY: the one existing guard on `compatibility.json` (`parley-addon.json`'s sha256 pin)
   verifies byte integrity, not version freshness, which is exactly why four releases drifted
   silently. My factual premise was false, so the conclusion does not stand. I now support the
   guard: one equality assertion (`skillVersion === package.json.version`) in the existing Node
   test suite. An assertion inside a suite that already runs is not new tooling; the constraint
   is "no new tooling unless justified," and four releases of silent drift is the justification.
   See @codex-1 for one implementation detail I think the proposal misses.

2. **REVISED: the "source of truth" sentence (SKILL.md:251).** In round 1 I called it "still
   correct but now incomplete." kimi-1's round-1 finding (item E.1) corrects me: the sentence is
   not merely incomplete, it is **inverted**. "The source of truth for each agent's mode is the
   spec's `autonomous_write` field" points readers at the field that 1.39.0 established is not
   what launches. A facilitator who follows it checks the declaration and stops; the effective
   argv is never inspected. That is the pre-1.39.0 mental model, left standing inside the fixed
   text. I should have caught this in round 1 and did not. It must be rewritten, not annotated.

3. **CONFIRMED (PRIMARY), newly: codex-1's finding that the launcher synthesizes no
   model/thinking/profile flags.** I did not verify this in round 1. I have now checked it
   myself. Scope checked: `internal/runner/runner.go:1097-1122` (`buildAgentInvocation`) and
   `internal/agents/discover.go:196-367` (all six built-in specs' `HeadlessArgs`). The launcher
   builds args exclusively from `agent.HeadlessArgs` with `{root}`/`{prompt}` substitution; no
   code appends model, thinking, or profile flags. The `Model`/`Reasoning`/`Profile` fields on
   the spec struct are read only for the "Effective launch config" display summary
   (runner.go:841-867) and for roster naming (naming.go:191-193), never for argv construction.
   For claude, hermes, and agy, the model flag happens to be pre-baked into `HeadlessArgs` in
   the built-in spec (e.g. claude: `["-p", "--model", "claude-opus-4-8[1m]", "--effort", "max",
   ...]`); for kimi and opencode it is not, and nothing adds it. So the skill's step 4 ("Add
   model/thinking/profile flags only when discovered or configured") describes a synthesis the
   CLI does not perform. This makes codex-1's finding correct and widens the fix beyond
   `writeModeArgs` — see Q2 below.

## Responses to others

### @claude-1

**On the guard (Q1):** I concede. Your argument that "an equality assertion in a suite that
already runs is not new tooling" is right, and my objection was built on a false premise (cited
above). I accept the guard.

**On naming the hermes incident (Q4):** You said in round 1: "If two participants think the
clause should go, it goes." The count is one (you) for naming it, three (codex-1, kimi-1, me)
against. Under your own threshold it does not go. But I will not rest on the count (§15.3). The
argument: the skill is instruction text, not a changelog. The general rule — "a config override
can replace `headless_args` wholesale without touching `autonomous_write`, so inspect the
effective argv" — is what changes facilitator behavior. The clause "this was live for `hermes`"
is the war story that motivates the rule. It does not tell a reader to do anything differently
than the rule already does. A reader who treats the rule as hypothetical despite a concrete
mechanism description is not going to be persuaded by an agent name. I contest naming it.

**On your fail-closed wording:** Your proposed text names `parley agents list`, the `AUTO`
column, and the `hermes` incident inline. I prefer kimi-1's structure: state the
vendor-neutral invariant first ("inspect the effective command"), then name the CLI surface
second as the tool that performs the check. Your version merges them, which couples the rule to
CLI output formatting that may change. See the convergence proposal below.

**On your own record ("corrected five times today"):** §15.2 says provenance controls
admissibility, not the count. Your item 4 claim (`writeModeArgs` not consumed at launch) is
PRIMARY and every other participant confirmed it. The self-flagging is noted but does not
downgrade the claim; the evidence stands on its own.

### @codex-1

**On the guard — one implementation gap you should address.** You proposed adding the assertion
to the `prepack` path as well as `npm test`. I checked: `prepack` runs only
`node scripts/build-addon-manifest.js --check` (package.json:66), which verifies byte-integrity
hashes and does not read `package.json.version` or `compatibility.json.skillVersion`. It does
not run `node --test`. So an assertion in `test/manifest-coverage.test.js` covers `npm test` but
NOT `npm publish` (which triggers `prepack`). To close the gap you have two options: (a) change
`prepack` to `node scripts/build-addon-manifest.js --check && node --test` (or a narrower
single-test invocation), or (b) add the version-equality check into
`scripts/build-addon-manifest.js` itself so the existing `prepack` path enforces it. I prefer
(b): it puts the check where the manifest already runs, avoids widening `prepack` to the full
test suite, and catches the `npm publish` without `npm test` case. Your call, but the gap
should be named in the final proposal either way. This is my own finding (PRIMARY:
package.json:60,66; scripts/build-addon-manifest.js:26).

**On Q2 (the recipe):** Your finding is correct and I confirmed it PRIMARY. Your proposed fix
is to "distinguish a manual assembly abstraction from CLI config, or state that all process
arguments must be present in the effective argv template." I agree with the second branch. See
my ruling below.

**On the fail-closed wording:** Your proposed text is strong. Your first sentence — "The
declared autonomous-write mode is a verification contract, not a second set of launch
arguments" — is the cleanest framing of the invariant in any round-1 file. I prefer it to my
own round-1 wording. But your text does not name the CLI surface that performs the check
(`agents list`), which kimi-1's does. I think the final text should take your first sentence as
the lead and add kimi-1's CLI-surface sentence. See convergence below.

**On the JSON-vs-TOML config mismatch you and kimi-1 both raised:** This is real (PRIMARY:
runtime.go:440 parses TOML; configLayers:137-154 reads `agents.toml`/`agents.local.toml`, never
the skill's `headless-agents.local.json`). But it predates 1.39.0 and you both agree it should
not enlarge this idea's diff. I support deferring it to a follow-up idea or inbox note. It
should be recorded but not fixed here.

### @kimi-1

**On SKILL.md:251 (Q3):** You are right and I was wrong. I called it "incomplete"; you called it
"inverted," and that is the more accurate characterization. The sentence actively teaches the
pre-1.39.0 model. Your rewrite in item C subsumes mine. I support your replacement wording as
the base, with one reservation below.

**On keeping the two-list JSON shape (`headlessArgs` + `writeModeArgs`):** You argue the defect
is the missing invariant and the missing check, not the field split, so deleting
`writeModeArgs` "churns the shape without adding the guard." I disagree. The field documents a
config input the CLI never reads — codex-1 confirmed this PRIMARY (zero grep hits for
`writeModeArgs` in Go source) and so did I. Keeping a dead field in the documented JSON shape
perpetuates the exact mental model this idea is fixing: that `writeModeArgs` is a separate
launch input. The guard and the invariant line you propose are necessary, but they are not
sufficient if the shape still shows two fields where the CLI has one. A facilitator copying the
JSON example will still populate `writeModeArgs` and still believe it matters. Remove it from
the JSON shape; add the invariant line. Both, not one.

**On WORKED_EXAMPLES.md:33:** Good catch, and I missed it. The example carries the same
`writeModeArgs` field and uses `--workspace-write` as a single token, which matches no shipped
adapter (codex's is two tokens: `--sandbox workspace-write`). Whatever lands in the JSON shape
must be mirrored there. If we remove `writeModeArgs`, the example must be updated in the same
edit. This is a second stale site the brief missed — I confirm it PRIMARY
(WORKED_EXAMPLES.md:31-33; discover.go:196).

**On naming CLI surfaces in vendor-neutral text:** You note precedent exists (SKILL.md:257
names `parley roster show`) and offer a fallback: if consensus prefers fully vendor-neutral
wording, the generic first two sentences stand alone. I think naming the surface is correct
here — a check described with no tool that performs it is what the pre-1.39.0 text had, and it
is why the defect went unnoticed. The invariant should be stated first (so it survives a CLI
rename), with the tool named second as the way to perform it. This matches the existing
precedent.

## New concerns / questions

1. **Q2 ruling — rewrite the whole recipe, or mark a boundary?** The brief asks everyone to
   rule. My ruling: **mark a boundary, do not rewrite the whole recipe.** The multi-step
   construction (steps 1-6) is a legitimate manual-assembly abstraction for a facilitator who
   launches agents by hand without the parley CLI. But it needs a boundary that marks where CLI
   config takes over. Specifically:

   - Steps 2 and 3 must collapse: `headlessArgs` is the complete launch vector, including the
     write-enabling flag. There is no separate `writeModeArgs` appended at launch. (This is
     the round-1 consensus, settled.)
   - Step 4 ("Add model/thinking/profile flags") must be bounded: when the parley CLI drives the
     agents, it launches `HeadlessArgs` verbatim and synthesizes no model/thinking/profile
     flags — those fields are display and naming only. For manual assembly, the facilitator
     may add them, but they must end up inside the launched argv. The skill should say this
     explicitly rather than presenting step 4 as something the CLI does.

   The reason I prefer a boundary over a full rewrite: the recipe has two audiences — a
   facilitator using the parley CLI (where the recipe is wrong about synthesis) and a
   facilitator assembling a command by hand (where the recipe is a useful abstraction). A full
   rewrite that says only "everything must be in `headlessArgs`" serves the first audience but
   removes the only guidance the second audience has. A boundary serves both. codex-1's
   proposed wording ("`headlessArgs` is the complete argv template ... plus exactly one
   `{prompt}` placeholder") is close to this, but it does not distinguish the two audiences; I
   would add one sentence marking where CLI config takes over.

2. **The `prepack` gap.** See @codex-1 above. The final proposal should name whether the
   version-equality guard runs on `prepack` (the `npm publish` gate) or only on `npm test`. If
   only on `npm test`, a direct `npm publish` without `npm test` bypasses it — which is the
   exact release path that produced the four-release drift.

3. **Schema bump for removing `writeModeArgs`?** My round-1 open question. The
   `compatibility.json` `schemaVersion` is 1. Removing a field from the documented (but not
   machine-validated) JSON config shape does not change runtime behavior — the CLI never read
   the field. I lean toward no schema bump: the field was dead, not deprecated. But if the
   guard assertion lands in the same release, the test itself becomes the proof that the shape
   is consistent. Worth confirming in round 3 or consensus.

## Current proposal

### Fail-closed wording — convergence of hermes-1 and codex-1

My round-1 wording and codex-1's round-1 wording both work but from different directions. I now
prefer codex-1's first sentence as the lead, kimi-1's CLI-surface naming as the body, and the
existing confinement sentence retained at the end. Converged text, replacing SKILL.md:251:

> The declared autonomous-write mode is a verification contract, not a second set of launch
> arguments. Before treating a headless participant as able to write its artifact, inspect the
> effective launch arguments after all configuration layers have been applied and verify that
> every argument required by the declared mode is present. A config override can replace the
> launch arguments and silently drop the enabling flag. When the parley CLI drives the agents,
> `parley agents list` performs this check and fails closed: it reports `AUTO=no`, names the
> missing flags, and prints the effective argv. If the enabling flag is absent, or the effective
> arguments cannot be inspected, treat the agent as non-autonomous until the config is fixed —
> do not assume the declared mode applies. Passing this check proves only that the autonomous
> mode is enabled; it does not prove workspace confinement. If workspace confinement cannot be
> demonstrated for an agent, treat its autonomous bit as unset (fail-closed) rather than
> escalating to a full-filesystem bypass.

This takes codex-1's lead sentence, kimi-1's CLI-surface sentence (with the invariant stated
first so it survives a CLI rename), and retains the existing confinement floor. It does not
name the hermes incident. claude-1's wording is functionally equivalent on the invariant but
merges the CLI surface and the incident into the same sentence; I prefer the separated
structure above.

### Summary of positions for this round

- **opencode row: CARRY.** No change from round 1.
- **kimi row: ALREADY COVERED.** No change.
- **Fail-closed wording: CARRY** (converged text above). Replaces SKILL.md:251 including the
  inverted "source of truth" sentence.
- **`compatibility.json`: BUMP + GUARD.** Withdrawn from round-1 "bump only." One equality
  assertion in the existing Node test suite. The `prepack` gap (see @codex-1) should be
  resolved: preferably by adding the check to `scripts/build-addon-manifest.js` so the existing
  `prepack` path enforces it, not by widening `prepack` to the full test suite.
- **`writeModeArgs` in JSON config shape: REMOVE.** From both SKILL.md:363 and
  WORKED_EXAMPLES.md:33. The field is dead (CLI never reads it); keeping it perpetuates the
  two-list mental model. kimi-1 disagrees (keep the shape, add the invariant); I ask kimi-1 to
  address the copying-facilitator problem in round 3.
- **Invocation contract steps 2-3: COLLAPSE** into one step (round-1 consensus, settled).
- **Invocation contract step 4: MARK A BOUNDARY.** State that the parley CLI launches
  `HeadlessArgs` verbatim and synthesizes no model/thinking/profile flags; for manual assembly,
  those flags must end up inside the launched argv. Do not rewrite the whole recipe.
- **hermes incident narrative: LEAVE OUT.** One against three; the rule carries without it.
- **JSON-vs-TOML config mismatch: DEFER.** Real, pre-dates 1.39.0, should not enlarge this
  diff. Record in an inbox note or follow-up idea.
- **Bundled `references/COOPERATION.md`: UNCHANGED.** No change from round 1.
