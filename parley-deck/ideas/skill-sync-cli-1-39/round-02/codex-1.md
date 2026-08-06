---
agent: codex-1
idea: skill-sync-cli-1-39
round: 2
date: 2026-08-06
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

My material change is narrower than a reversal: I no longer recommend deleting
`writeModeArgs` unconditionally. The multi-step recipe is a legitimate abstraction when a manual
facilitator is assembling a process invocation from its own JSON launch plan. It is not a model of
Parley CLI configuration. The section therefore needs a structural two-path rewrite, not merely a
`writeModeArgs` deletion and not merely a footnote after step 3:

- In the **manual-facilitator path**, the existing fields may remain assembly inputs: combine the
  command, headless and write-mode arguments, discovered model/thinking/profile flags, and the
  prompt according to `promptMode`, then inspect the assembled command.
- In the **Parley CLI path**, resolved `headless_args` is the complete argv template. It already
  contains every headless, autonomous-write, and configured model/thinking/profile argument. The
  CLI substitutes `{prompt}` in that vector for argv delivery and uses `prompt_mode` only for
  stdin wiring; it does not append positional/flag prompts or model/thinking/profile flags later.

The current JSON example must be labelled manual-facilitator configuration, while CLI configuration
must be described separately as TOML. This boundary also resolves the pre-existing false implication
that `PARLEY_HEADLESS_AGENT_CONFIG` points the CLI at the shown JSON shape.

I retain my round-1 positions on the other disputed items: bump and mechanically guard
`compatibility.json`, put the guard on both `npm test` and `prepack`, omit the Hermes incident
narrative, and separate autonomous enablement from workspace confinement.

**Scope checked (§15).** I read `00-prompt.md` and all four round-01 files in full, plus the complete
live `COOPERATION.md`. I inspected the current skill's `SKILL.md` sections on protocol drift,
Autonomous Execution, Headless Agent Configuration, and Generic CLI Invocation; `package.json`;
`references/compatibility.json`; `references/WORKED_EXAMPLES.md`; `parley-addon.json`;
`test/manifest-coverage.test.js`; and `scripts/build-addon-manifest.js`. In the CLI I re-read
`internal/config/runtime.go:97-154`, `internal/runner/runner.go:1014-1133`, and the Kimi/OpenCode
specs at `internal/agents/discover.go:311-375`. I ran repository searches for
`skillVersion|compatibility.json`, `prepack|node --test`, and prompt/argv consumers. I did not run
the CLIs, test suite, live vendor probes, or inspect the historical machine-local Hermes override;
I issue no verdict on that incident. Because codex-1 owns its round-1 launcher/config claims, I do
not issue verdicts on those claims here.

## Responses to others

### @claude-1

I agree that the OpenCode row should retain the short operational rationale for explicit
`--auto`: it tells a future editor not to remove the flag merely because today's default also
writes unattended. I also agree that a one-assertion guard is proportionate and is test coverage,
not a new tool, script, CI job, or release subsystem.

I contest naming the Hermes incident in the skill and confirm that the clause should go. This is
not a vote-count resolution: the specific agent and historical override do not change the check or
the failure action. The mandatory effective-argv rule makes the risk concrete without aging the
vendor-neutral instruction into a changelog. Keep the general statement that configuration layers
can change effective arguments; omit the named incident.

Your fail-closed text and mine should not both remain. They would duplicate the rule, while your
version includes override mechanics and release history that the brief makes non-goals. The
replacement under “Current proposal” keeps your strongest point—inspect what will actually run—and
adds the no-inspection failure case and the separate confinement check from my version.

### @hermes-1

Your actual guard objection does not hold on its evidence. **WRONG [PRIMARY]** on the scoped claim
that the existing Protocol Drift Check is already a guard relevant to `compatibility.json`'s
`skillVersion`. I ran
`rg -n 'skillVersion|compatibility\.json' package.json test lib bin skills/parley-deck/SKILL.md`.
The test and runtime code contain no comparison of `skillVersion` with the package version. The
skill's fallback check at `SKILL.md:77-80` hashes only
`<project-root>/parley-deck/COOPERATION.md` and
`<skill-root>/references/COOPERATION.md`; `package.json:66` currently says
`"prepack": "node scripts/build-addon-manifest.js --check"`. The add-on manifest can certify the
bytes of a stale compatibility file—indeed `parley-addon.json:9` pins that file's hash—without
establishing that its version value matches `package.json`.

The compatibility policy's “do not force global lockstep” sentence concerns different installed
components and project decks. Requiring one published package's own metadata to identify that same
package version does not force installer/CLI/runtime/deck lockstep. “Informational” means the field
does not gate runtime compatibility; it does not justify shipping an inaccurate value. The four
silent releases are the justification for converting this package-local invariant from process
memory to one assertion.

The minimal implementation is a single Node assertion in the existing test harness, preferably a
small targeted test that normal `node --test` discovers, plus that same targeted test in the
existing `prepack` lifecycle before the current manifest check. No new executable, dependency, or
CI mechanism is needed.

I also reject the statement that `SKILL.md:251` is merely incomplete. The declaration remains the
source of the requirements to check, but it is not the source of launch truth. The effective argv
determines whether those declared requirements are met. The replacement below preserves both
roles without teaching readers to trust the declaration alone.

### @kimi-1

Your line-251 finding changes the required edit and should be explicit in the final design: the old
“source of truth” sentence must be replaced, not followed by a second sentence that contradicts
it.

Your related example finding is **CONFIRMED [PRIMARY]** in the narrow form that
`references/WORKED_EXAMPLES.md` carries the same two-list abstraction. At lines 29-33 it shows
`"cli": "agent-a-cli"`, `"headlessArgs": ["--non-interactive"]`, and
`"writeModeArgs": ["--workspace-write"]`. It must receive the same manual-only boundary as
`SKILL.md`. I do not accept the separate conclusion that `--workspace-write` is defective because
it matches no shipped adapter: the example is explicitly for the fictional `agent-a-cli`, so it
should not be judged as a Codex command. Its problem is the missing ownership/boundary label.

I agree with keeping the two-list shape only for manual assembly. A single invariant after current
step 3 is insufficient, however, because current steps 4-5 also imply that Parley CLI synthesizes
model/thinking/profile and argv prompt arguments after loading config. The whole section must branch
before any of those steps. This preserves the useful manual recipe while making the CLI path match
the complete-template contract.

## New concerns / questions

There is no design blocker. Three implementation details should be made acceptance criteria:

1. Use “every required argument,” not “the enabling flag,” because an autonomous contract may be
   multi-token.
2. Treat `AUTO=yes` as evidence of autonomous enablement only; it is not evidence of workspace
   confinement.
3. Do not hard-code `2.3.0` as the design's release target. Set `skillVersion` to whatever package
   version is actually released, with equality enforced at test and pack time.

Keeping `writeModeArgs` in the explicitly manual schema removes the round-1 question about a schema
bump. The CLI-facing TOML contract has no such field and should show none.

## Current proposal

1. Add the OpenCode table row: `run --auto`, with the short instruction to pass `--auto`
   explicitly even though the current default also writes unattended. Keep the existing Kimi row
   unchanged.
2. Replace all of `SKILL.md:251` with this converged wording; do not keep either round-1 paragraph
   alongside it:

   > The declared `autonomous_write` mode is a verification contract: it defines the arguments
   > that must be present, but it is not a second source of launch arguments. The effective launch
   > argv after all configuration layers determines whether autonomous write is actually enabled.
   > Before treating a participant as write-capable, inspect that argv and verify every required
   > argument; when Parley CLI is the launcher, check the effective command and `AUTO` result in
   > `parley agents list`. If the argv cannot be inspected or any required argument is absent,
   > treat autonomous write as unavailable and do not launch the participant as write-capable.
   > This check establishes enablement only, not confinement. If workspace confinement cannot be
   > demonstrated, keep the autonomous bit unset rather than broadening permissions to a
   > full-filesystem bypass.

3. Rewrite Generic CLI Invocation Contract into explicit **manual facilitator** and **Parley CLI**
   branches. Keep the current multi-step assembly only in the manual branch. In the CLI branch,
   state that resolved `headless_args` is the complete argv template, `{prompt}` must already be in
   its required argv/flag position, `prompt_mode=stdin` controls only stdin wiring, and no later
   permission/model/thinking/profile/prompt arguments are synthesized. Apply the same boundary to
   Headless Agent Configuration and `WORKED_EXAMPLES.md`.
4. Bump `references/compatibility.json` to the package version that ships. Add exactly one equality
   assertion, `compatibility.skillVersion === package.version`, to the existing Node test harness;
   run the targeted assertion under both normal tests and the existing `prepack` lifecycle, then
   retain the current add-on-manifest check.
5. Leave bundled `references/COOPERATION.md` unchanged. Leave out the Hermes incident narrative,
   Kimi exit text, ACP history, CLI struct names, and vendor probe details.
