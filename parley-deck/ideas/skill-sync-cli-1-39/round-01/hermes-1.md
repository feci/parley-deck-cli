---
agent: hermes-1
idea: skill-sync-cli-1-39
round: 1
date: 2026-08-06
---

## Summary

1.39.0 changed two things that matter for what a facilitator or participant DOES: a new adapter (opencode) needs an "Autonomous Execution" table row, and the fail-closed rule now covers a second failure mode (declared mode whose enabling flag the effective launch does not pass) that the skill's existing fail-closed sentence does not cover. The kimi row is already correct. compatibility.json's `skillVersion` is stale by four releases and should be bumped; a guard is not justified under the no-new-tooling constraint. The facilitator's claim about the "Generic CLI Invocation Contract" is CONFIRMED by primary source inspection: the skill teaches a two-step launch model (headlessArgs + writeModeArgs) that the CLI does not implement — the launch consumes `HeadlessArgs` alone, and `AutonomousWrite.Args` is read only by the fail-closed display path, never at launch.

## Proposed approach

### Item 1 — kimi and opencode promoted to full built-in adapters

**opencode table row: CARRY.** The "Autonomous Execution" table (SKILL.md lines 243-249) has rows for claude, codex, hermes, agy, and kimi — no opencode. Add:

```
| opencode (OpenCode) | `run --auto` — `--auto` is the explicit auto-approve flag; pass it explicitly. ACP remains available via `opencode acp`. |
```

**kimi row: ALREADY COVERED.** The existing row (line 249) already says: "plain `-p` — its print mode already auto-approves in-workspace writes. NOTE: `--yolo`/`--auto` are mutually exclusive with `-p`, so `-p` IS kimi's yolo-equivalent." This matches the CLI spec exactly (`HeadlessArgs: ["-p", "{prompt}"]`, `AutonomousWrite{Mode: "prompt", Args: ["-p"]}`). No change needed.

**LEAVE OUT — the exit-code detail.** The brief notes `kimi --auto -p` exits 1 with "Cannot combine --prompt with --auto." The skill already states the mutual exclusivity. The specific exit code and error string are CLI-internals that do not change what a facilitator does; the skill's instruction ("`-p` IS kimi's yolo-equivalent") is sufficient.

**LEAVE OUT — "opencode writes unattended even without --auto."** This is a changelog fact about vendor behavior. What changes what the facilitator does is: pass `--auto` explicitly. The skill row should say that. Why the CLI chose to pass it explicitly (vendor may change implicit defaults between versions) is rationale, not instruction.

### Item 2 — The AUTO signal now fails closed

**The fail-closed sentence: CARRY (expand it).** The current sentence (SKILL.md line 251) reads: "If workspace confinement cannot be demonstrated for an agent, treat its autonomous bit as unset (fail-closed) rather than escalating to a full-filesystem bypass." This covers only the confinement case. 1.39.0 introduced a second fail-closed case: a declared mode whose enabling flag the effective launch argv does not pass. The skill must cover both.

**Proposed exact wording** (replace the sentence at line 251):

> Before trusting that a participant can write its artifact, verify that the agent's **effective launch argv** actually contains the enabling flag for its declared autonomous-write mode — not merely that the mode is declared. A config override can replace `headlessArgs` wholesale without touching `autonomous_write`, so a spec may declare a mode whose flag the launched command never passes. If the enabling flag is absent from the effective argv, or if workspace confinement cannot be demonstrated for an agent, treat its autonomous bit as unset (fail-closed) rather than escalating to a full-filesystem bypass.

This tells the facilitator exactly what to check: the effective argv, not the declared mode. It does not document `AutonomousWrite.MissingFrom` or spec structs (that is a non-goal).

**LEAVE OUT — the hermes defect specifics.** "An override had dropped `--yolo` while it still reported `AUTO=yes`" is the bug narrative. The skill needs the general rule, not the specific incident. The facilitator checks the effective argv regardless of which agent or which flag.

**LEAVE OUT — "`agents list` prints the effective argv."** This is a CLI feature. The skill should tell the facilitator WHAT to check (effective argv contains the enabling flag); it should not document HOW the CLI surfaces it (the `agents list` command format). An agent following the skill manually can check the argv in the orchestration summary or the capability matrix.

### Item 3 — compatibility.json

**Bump `skillVersion`: CARRY.** The file says `1.4.3` while the package is `2.3.0` — stale by four releases. Bump it to `2.3.0` as part of this skill release. This is a mechanical fix.

**Guard: LEAVE OUT (not justified).** The staleness went unnoticed for four releases, which argues for a guard. But the constraint is "no new tooling unless justified," and the skill already has a drift-detection mechanism: the "Protocol Drift Check" section (SKILL.md lines 52-80) runs `parley-deck-skill status` and compares hashes. The `skillVersion` field is informational metadata, not a gate — the `compatibilityPolicy` itself says "Report actual component versions and warn on drift; do not force global lockstep." Adding a guard that compares `skillVersion` to `package.json` would be new tooling for a field whose own policy says it is informational. The fix is process discipline (bump it every release, which the owner's standing rule now enforces), not a guard.

### Item 4 — The facilitator's claim about the "Generic CLI Invocation Contract"

**CONFIRMED (PRIMARY).**

The facilitator claims the skill's "Generic CLI Invocation Contract" and its JSON config shape teach a mental model that would reproduce the hermes defect: they tell the facilitator to add `writeModeArgs` as a step separate from `headlessArgs`, but the CLI launch never consumes `AutonomousWrite.Args`.

I checked this against the source. Scope checked: `internal/runner/runner.go` (`buildAgentInvocation`, lines 1097-1122), `internal/agents/discover.go` (`AutonomousWrite` struct, `MissingFrom`, `AutonomousEffective`, `PrintRuntimeMatrix`), `internal/config/runtime.go` (`agentOverride` struct, lines 97-132, and `applyOverride`, lines 514+), and every `.go` file matching `AutonomousWrite\.Args` or `writeModeArgs` across the CLI.

Findings:

1. **The launch uses `HeadlessArgs` alone.** `buildAgentInvocation` (runner.go:1098-1108) builds the args slice by iterating `agent.HeadlessArgs` only, substituting `{root}` and `{prompt}`. There is no reference to `AutonomousWrite.Args`, `WriteModeArgs`, or `writeModeArgs` anywhere in the launch path.

2. **`AutonomousWrite.Args` is consumed only by the fail-closed display path.** The field is read by `MissingFrom` (discover.go:115-130), `AutonomousEffective` (discover.go:138-140), and `PrintRuntimeMatrix` (discover.go:538-541) — the two `AUTO` display sites and the warning. The tests in `autonomous_test.go` exercise it. No launch site reads it.

3. **The config layer has no `writeModeArgs` field.** The `agentOverride` struct (runtime.go:97-132) has `HeadlessArgs` (`toml:"headless_args"`) but no `WriteModeArgs`, `write_mode_args`, or `autonomous_write` field. A config override can replace `HeadlessArgs` (runtime.go:542-544) but cannot touch `AutonomousWrite` at all. So `writeModeArgs` in the skill's JSON config shape is a field the CLI never reads.

4. **The built-in specs embed the write-enabling flag in `HeadlessArgs`.** Every spec has its `AutonomousWrite.Args` as a subset of `HeadlessArgs`: hermes has `--yolo` in both; claude has `--permission-mode bypassPermissions` in both; codex has `-c approval_policy="never"` in both; kimi has `-p` in both; opencode has `--auto` in both. The opencode spec comment (discover.go:354) states this explicitly: "Declaring it keeps AutonomousWrite.Args a subset of HeadlessArgs, as for claude/codex/hermes."

5. **The skill teaches the opposite model.** SKILL.md "Generic CLI Invocation Contract" (lines 807-816) says: "1. Start with the configured `cli`. 2. Add `headlessArgs`. 3. Add `writeModeArgs` needed for the agent to write exactly one protocol artifact." Step 3 tells the facilitator that `writeModeArgs` is a separate set of flags appended at launch. The JSON config shape (lines 348-375) reinforces this by listing `headlessArgs` and `writeModeArgs` as separate fields. But the CLI never appends `writeModeArgs` — the write-enabling flag must already be IN `headlessArgs`, or it will not be in the launched command.

**Conclusion: the claim is correct.** The skill teaches the exact mental model that made the hermes defect invisible. A facilitator following the skill's step-2-then-step-3 construction could put `--yolo` in `writeModeArgs` (or believe it is carried separately), override `headlessArgs` to drop it, and think the agent is still autonomous because the mode is "declared." The CLI's new `AutonomousEffective` fail-closed check now catches this at display time, but the skill's launch model is still wrong.

**Proposed fix:** Replace the two-step construction (steps 2 and 3) with a single step that states the reality:

> 2. Add `headlessArgs` — these are the complete launch flags, including the agent's autonomous-write enabling flag. The enabling flag must be present in `headlessArgs`; there is no separate `writeModeArgs` step at launch. A config override that replaces `headlessArgs` replaces the write-enabling flag too — verify the effective argv still carries it.

And remove `writeModeArgs` from the JSON config shape (lines 362-363), or repurpose it as an informational field that mirrors which flag in `headlessArgs` is the write-enabling one (not a separate launch input). Removing it is cleaner: it is a field the CLI never reads, and keeping it perpetuates the two-step mental model.

### Item 5 — Anything 1.39.0 invalidates that the brief missed

**The JSON config shape's `writeModeArgs` field is dead.** This is a consequence of the confirmed claim above but is a separate staleness item: the skill's documented JSON config (SKILL.md lines 348-375) includes `"writeModeArgs": ["<arg>", "<arg>"]` as a field a facilitator should fill. The CLI's config layer (`agentOverride` struct) has no such field — it would be silently ignored if a facilitator put anything there. This is not explicitly listed in the brief's "What is stale or missing" section, and it is invalidated by 1.39.0's fail-closed work (which made the disconnect between `AutonomousWrite.Args` and `HeadlessArgs` visible). Fix: remove `writeModeArgs` from the JSON config shape, or document that the write-enabling flag belongs in `headlessArgs`.

**The "source of truth" sentence is still correct but now incomplete.** SKILL.md line 251 says: "The source of truth for each agent's mode is the spec's `autonomous_write` field; a vendor flag change is a config edit, not a skill revision." After 1.39.0, the `autonomous_write` field is the source of truth for what is *declared*, but `AutonomousEffective` (the fail-closed check) is the source of truth for what is *actually enabled at launch*. The sentence is not wrong, but it is now incomplete in a way that matters: a facilitator who reads it as "check `autonomous_write` and you're done" would miss the effective-argv check. The proposed fail-closed wording above addresses this.

## Concerns / open questions

1. **Removing `writeModeArgs` from the JSON config shape is a schema change.** The `compatibility.json` `schemaVersion` is 1. Does removing a field from the documented (but not machine-validated) JSON config shape require a schema bump? The field was never read by the CLI, so no runtime behavior changes — but a facilitator who has been populating it would see it disappear from the docs. I lean toward no schema bump (the field was dead, not deprecated), but this is worth confirming in round 2.

2. **Should the skill name the opencode adapter's model/effort discovery?** The CLI spec says model comes from `~/.config/opencode/opencode.jsonc` and reasoning effort is `--variant <level>`. The skill's "Agent Capability Discovery" section already says to use `<cli> --help` / `models` / `model list` and not assume they exist. Is that sufficient, or does opencode need a discovery note? I lean toward sufficient — the skill is vendor-neutral and should not carry per-vendor discovery recipes.

3. **The fail-closed wording says "verify that the effective launch argv actually contains the enabling flag."** For a facilitator running manually (not through the CLI), how do they check this? The skill's "Headless Agent Configuration" section says to "Record the effective launch config in the orchestration summary." Is recording it enough, or should the skill say to compare the recorded argv against the autonomous-write flag? I think the proposed wording is sufficient — "verify" implies the comparison — but this may need tightening.

## Risks

1. **Over-fixing the JSON config shape.** Removing `writeModeArgs` from the documented JSON config is the clean fix, but if any downstream tooling or facilitator workflow reads that field (even though the CLI does not), removing it from the docs could break an undocumented contract. Risk is low — the CLI never reads it — but worth a grep of the skill repo in round 2.

2. **The fail-closed wording adds a check the facilitator may not know how to perform.** "Verify the effective launch argv" is clear in principle but requires the facilitator to know where to find the effective argv. If the CLI is not installed or the `agents list` command is not run, the facilitator must reconstruct it from the config. The skill should not document CLI commands (non-goal), but the wording should be testable by a facilitator with only the skill and the config file.

3. **Scope creep into CLI internals.** The proposed fail-closed wording and the `headlessArgs` fix stay at the agent-instruction level. But the temptation to explain *why* the config override can drop the flag (the `applyOverride` / `HeadlessArgs` replacement mechanism) is a slide into documenting CLI internals, which is a non-goal. The wording must state the rule without explaining the mechanism.
