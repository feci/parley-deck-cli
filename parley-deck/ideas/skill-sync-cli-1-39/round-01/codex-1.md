---
agent: codex-1
idea: skill-sync-cli-1-39
round: 1
date: 2026-08-06
---

## Summary

The skill needs a narrow operational sync, not a 1.39.0 changelog: add the OpenCode launch shape, make the effective-argument check binding, remove the separate `writeModeArgs` launch model, and synchronize plus guard `compatibility.json`. Kimi's actionable launch rule is already present. Promotion history, probe versions and output, the Hermes incident narrative, and ACP availability should stay out because they do not change a facilitator's or participant's next action.

**Scope checked.** I read the complete live `COOPERATION.md` (including §15), the complete packaged `SKILL.md`, the idea brief, and `compatibility.json`. On the CLI side I inspected `internal/agents/discover.go` (the autonomous contract, Kimi/OpenCode/Hermes specs, and `AUTO` display), `internal/config/runtime.go` (the full config type/layer definitions and the override around lines 542–544), `internal/runner/runner.go` (round/steer invocation and `CommandFor`), `internal/app/consensus_request_signoffs.go` (the separate headless-signoff path), and ran a repository-wide Go-symbol search for `AutonomousWrite`, `MissingFrom`, `AutonomousEffective`, and `HeadlessArgs`. In the skill repo I also checked `package.json`, the installer use of the compatibility manifest, the existing test/prepack commands, and `RELEASING.md`. I did not read any other participant's `skill-sync-cli-1-39/round-01` file. I did not execute Kimi or OpenCode, and I did not inspect the historical machine-local Hermes override, so I issue no verdict on those live-probe or incident-history claims; neither is needed for this recommendation.

**Facilitator-owned claim verdict: CONFIRMED [PRIMARY], with a necessary narrowing.** The skill currently places `headlessArgs` and `writeModeArgs` in separate fields (`parley-deck-skill/skills/parley-deck/SKILL.md:361-363`) and says to append them as separate launch steps (`SKILL.md:807-814`). The CLI config has only `HeadlessArgs` (`parley-deck-cli/internal/config/runtime.go:97-131`), an override replaces that slice wholesale (`runtime.go:542-544`), and the common launcher constructs `args` exclusively by iterating `agent.HeadlessArgs` (`internal/runner/runner.go:1094-1108`). Round/steer execution launches those returned args (`runner.go:1023-1031`); headless signoff calls the same `CommandFor` path (`internal/app/consensus_request_signoffs.go:431-448`; `runner.go:1124-1133`). The autonomous args are instead compared with the effective headless args by `MissingFrom`/`AutonomousEffective` (`internal/agents/discover.go:111-139`) and feed the `AUTO` presentations (`discover.go:505-540`; `internal/app/roster.go:155`). Thus a separate write-argument list is not launch input in the CLI and can recreate the same omission when a facilitator translates the skill shape into Parley config.

Relevant primary passages are: `SKILL.md:810-811`, “Add `headlessArgs`” then “Add `writeModeArgs`”; `runtime.go:542-544`, `spec.HeadlessArgs = expandSlice(override.HeadlessArgs, root, tempdir)`; `runner.go:1098-1100`, `args = make([]string, 0, len(agent.HeadlessArgs))` followed by `for _, arg := range agent.HeadlessArgs`; and `discover.go:138-139`, `return s.AutonomousWrite.Declared() && len(s.AutonomousWrite.MissingFrom(s.HeadlessArgs)) == 0`.

The narrowing is that a standalone manual launcher following both documented append steps literally would append the write arguments and would not itself omit them. The documentation is wrong as a model of the current Parley CLI and dangerously ambiguous as a supposedly shared local-config contract; it is not a proof that every independent implementation with two arrays fails.

The brief's three directly inspectable staleness claims are also **CONFIRMED [PRIMARY]**: the skill table contains Kimi but no OpenCode (`SKILL.md:239-251`), its only explicit fail-closed check is the confinement sentence at line 251, and `package.json:3` says `2.3.0` while `references/compatibility.json:4` says `1.4.3`.

## Proposed approach

| Brief item | Decision | Agent-facing reason or existing text |
| --- | --- | --- |
| Kimi full headless adapter | **Already covered** for the action that matters. | `SKILL.md:249` already says plain `-p` is the autonomous mode and that `--yolo`/`--auto` cannot be combined with it. Leave the promotion label, probe version/output, and ACP alternative out. |
| OpenCode full headless adapter | **Carry into the skill.** | Add an OpenCode row whose mode is `run --auto <prompt>` and which says the prompt is argv, not stdin. The fact that unattended writes may currently work without `--auto` is implementation history; the facilitator's action is to pass the explicit flag. |
| Effective-argv fail-closed behavior | **Carry into the skill.** | A facilitator must inspect the command that will actually launch after all overrides, not a built-in label or a declared mode. When the Parley CLI is available, `parley agents list` is the relevant inspection surface because 1.39.0 prints effective argv and fails `AUTO` closed. |
| Historical Hermes failure | **Leave out.** | It explains the change but does not alter the reusable instruction once the effective-argv rule is stated. |
| ACP remains available for Kimi/OpenCode | **Leave out.** | ACP is an alternative transport/runtime path, not a requirement for authoring the one canonical artifact in the ordinary headless path. |
| Separate `writeModeArgs` field and assembly step | **Carry a correction, not the old shape.** | Delete `writeModeArgs` from the generic config and command-construction list. Define `headlessArgs` as the complete effective argv template, including every non-interactive/autonomous enabling argument and, for argv delivery, one `{prompt}` placeholder in the required position. |
| JSON config as a Parley runtime source | **Carry a clarification/correction.** | `PARLEY_HEADLESS_AGENT_CONFIG` is parsed by the current CLI as TOML (`runtime.go:16,22-27,134-165`), and the CLI layers `~/.parley/agents.toml`, `parley-deck/agents.toml`, and `parley-deck/agents.local.toml`; it does not load the shown JSON shape as runtime launch config. Either label the JSON explicitly as manual-facilitator-only state or, preferably, show the CLI TOML separately and stop presenting the two as one contract. |
| Prompt placement and model/thinking assembly | **Carry as part of the same correction.** | Do not patch only `writeModeArgs`. The common CLI launcher substitutes `{prompt}` inside `HeadlessArgs` and otherwise launches only that vector (`runner.go:1094-1108`); it does not synthesize the skill's later `modelFlag`/`thinkingFlag` steps. The skill should distinguish a manual assembly abstraction from CLI config, or state that all process arguments must be present in the effective argv template. Kimi and OpenCode make this material because their built-ins are `[-p, {prompt}]` and `[run, --auto, {prompt}]` (`discover.go:312-374`). |
| `compatibility.json` version | **Carry the bump and a guard.** | Set `skillVersion` to the package version released from this change. A one-time bump is insufficient after silent multi-release drift. Add a small assertion to the existing Node test suite that `skillPackage === package.json.name` and `skillVersion === package.json.version`, and run that assertion from the existing `prepack` path as well as `npm test`. This uses existing tooling; no new checker or runtime behavior is justified. |
| Bundled `references/COOPERATION.md` | **Leave unchanged.** | The brief says the protocol diff is clean and the constraint forbids hand-editing the generated fallback. No facilitator or participant action changed, so protocol release history does not belong in the skill. |

Proposed exact fail-closed wording, replacing the final paragraph of **Autonomous Execution** while retaining the existing confinement protection:

> The declared autonomous-write mode is a verification contract, not a second set of launch arguments. Before treating a headless participant as able to write its artifact, inspect the effective launch arguments after all configuration layers have been applied and verify that every argument required by the declared autonomous-write mode is present. If the effective arguments cannot be inspected, or any required argument is absent, treat autonomous write as unavailable (`AUTO=no`) and do not launch that participant as write-capable. Passing this check proves only that the autonomous mode is enabled; it does not prove workspace confinement. If workspace confinement cannot be demonstrated for an agent, treat its autonomous bit as unset rather than escalating to a full-filesystem bypass.

For the **Generic CLI Invocation Contract**, replace the split steps with this operational rule:

> `headlessArgs` is the complete argv template that will be launched. It MUST already contain the selected CLI's non-interactive and autonomous-write enabling arguments, plus exactly one `{prompt}` placeholder when the prompt is delivered through argv or a value-taking flag. Apply configuration layers, inspect the resulting argv, and only then add no further permission flags implicitly.

The local-config section should then either show a manual-only JSON example using only `headlessArgs`, or show the actual Parley TOML form. It should not attach the current JSON schema to `PARLEY_HEADLESS_AGENT_CONFIG` without an explicit translator.

## Concerns / open questions

The skill currently conflates three properties: autonomous permission, prompt delivery, and enforced workspace confinement. CLI 1.39.0 deliberately treats declaration/effective enabling and confinement as distinct (`discover.go:86-139`), while `SKILL.md:241,251` can be read as one Boolean eligibility bit. Consensus should preserve the existing safety floor but name the two checks separately so `AUTO=yes` is never read as proof of sandboxing.

The broader config drift predates 1.39.0, but the new argv-prompt adapters make it actionable now. A minimal OpenCode row plus one fail-closed sentence would leave the skill still claiming that a JSON file and later flag-synthesis steps configure a CLI that actually parses TOML and launches only its effective headless vector. I recommend fixing that in the same release rather than recording it as a changelog footnote.

The exact release version is not specified by the brief. The invariant should be equality with `package.json` at pack time, rather than hard-coding `2.3.0` in the design when the implementation may correctly bump the package for this release.

## Risks

- If `writeModeArgs` remains, a config override can again look autonomous in metadata while omitting the enabling flag from the actual command.
- If the version check runs only under `npm test`, a direct `npm publish` can bypass it; the lightweight equality assertion must also be on the existing `prepack` path.
- If the OpenCode row is copied as vendor history rather than an effective-command rule, it will drift on the next vendor default change without improving facilitator behavior.
- If `AUTO` and confinement remain linguistically merged, facilitators may either trust an unconfined permission mode as sandboxed or incorrectly discard a correctly enabled mode without explaining which safety property failed.
- If JSON, TOML, manual assembly, and CLI launch are not explicitly separated, future fixes may update the descriptive fields while leaving the sole launched argv unchanged.
