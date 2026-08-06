---
idea: kimi-opencode-full-adapters
phase: review-consensus
review_round: 1
drafter: claude-1
participants: [claude-1, codex-1, hermes-1, kimi-1]
date: 2026-08-06
status: fixes-agreed
---

## Outcome

**One `MAJOR`, four `MINOR`, three `NIT`. Not ready to release.**

Two of the three reviewers (`hermes-1`, `kimi-1`) returned *ready to release*; `codex-1` blocked on
a `MAJOR` that neither of the others found. **The implementer independently confirmed the `MAJOR`
and it stands** — so the two release verdicts fall with it, not because they were outvoted but
because they rest on a premise the evidence contradicts.

### Drafter retraction, recorded first because it changes who was right

The implementer initially believed it had refuted `codex-1`'s `MAJOR` with a runtime probe showing
`opencode` effective args as `[run --auto {prompt}]`. **That probe was invalid.** It ran inside
`go test`, where `PARLEY_HOME` is isolated, so the real `~/.parley/agents.toml` was never read —
the measurement described a different environment than the one under review. Verdict on the
implementer's own refutation: **withdrawn**. `codex-1`'s source-level reading stands unrefuted.

This is the same failure mode this deck recorded against `hermes-1` earlier today at T1: a
measurement taken against the wrong environment and reported as a refutation.

## Agreed fixes

### AF-1 — Declared autonomous mode can be silently absent from the launch (MAJOR)

Raised by `codex-1`. Confirmed by the implementer with a direct scan of both config layers:

```
deck     hermes    headless_args set, declared '--yolo' present: *** NO ***
central  hermes    headless_args set, declared '--yolo' present: *** NO ***
central  opencode  headless_args set, declared '--auto' present: *** NO ***
central  kimi      headless_args set, declared '-p'     present: YES
```

`internal/config/runtime.go:542` replaces `spec.HeadlessArgs` wholesale when an override supplies
them, but never touches `AutonomousWrite`. So a config layer can strip the very flag that the
declared mode names, while `agents list` still prints `AUTO=yes` and a `headless:` line taken from
the built-in `HeadlessMode` **label** rather than the effective argv.

**The defect is wider than this change.** `hermes` is in this state **today**, before this idea:
it runs without `--yolo` in this deck while reporting `AUTO=yes`. The promotion of `opencode`
merely produced a third instance.

**Fix, three parts:**

1. **Fail closed.** When effective `HeadlessArgs` do not contain every `AutonomousWrite.Args`
   token, the agent MUST NOT report `AUTO=yes`. Report `no` and surface why.
2. **Show the truth.** `agents list` must derive its `headless:` line from effective
   `HeadlessArgs`, not from the built-in `HeadlessMode` label. That label is what misled the
   implementer.
3. **Repair the local config** (`~/.parley/agents.toml`, machine-local and gitignored): add
   `--auto` to the `opencode` override and `--yolo` to `hermes`, or drop the redundant overrides.

Part 1 is the honesty rule the `AutonomousWrite` type already states for `Scope`, applied one
level up: never assert an autonomous mode that the launched command does not enable.

### AF-2 — The autonomous contract test checks only `Mode` (MINOR)

`codex-1`. `autonomous_test.go` would still pass if `Args`, `HeadlessArgs`, `PromptMode`, `Scope`
or merged `ACPArgs` regressed. Extend the table to assert, for both new IDs: exact autonomous
args, exact headless argv template, `PromptArg`, empty `Scope`. Add a layered-config case proving
AF-1's fail-closed behaviour.

### AF-3 — Nothing locks the ACP-merge outcome for the promoted IDs (MINOR)

`kimi-1`. `TestDefaultSpecsMergesACPCatalog` covers claude/hermes/codex but not kimi/opencode. A
future edit to `mergeACPBackend` could silently drop `kimi acp` as a selectable mode — the exact
constraint `00-prompt.md` names — with every test green. Extend the merge test with both IDs
(`ACPArgs`, `LaunchMode`, one `HeadlessArgs` element).

### AF-4 — kimi is undiscoverable at its default install prefix (MINOR)

`kimi-1`, verified `PRIMARY`: `command -v kimi` returns rc=1; the binary lives only at
`~/.kimi-code/bin/kimi`. `Commands: ["kimi"]` resolves through `exec.LookPath`, i.e. PATH only, so
a fresh user who installed kimi normally gets `INSTALLED=no` and the newly promoted adapter never
engages. It works here **only** because the central config sets `command`.

Not a regression — every family here is reached through a central `command` — but for a newly
promoted adapter it means the feature does not work out of the box. **Fix:** state it in the
spec `Notes` ("set `command` if kimi is not on PATH"). Probing well-known prefixes at discovery
time is the larger alternative and is **not** adopted here: it expands discovery behaviour beyond
this idea's scope.

### AF-5 — `opencode` Telemetry value understates its output (NIT)

`hermes-1`, verified against its own probe: opencode streams build/tool events before the final
text. Change `"final text on stdout"` to something that says so.

### AF-6 — Stale guidance in the in-repo deck config (NIT)

`kimi-1`. `parley-deck/agents.toml` still says *"parley 1.36 treats kimi as an ACP backend
(AUTO=no) — invoke it manually, and do NOT run `parley roster init`"* and pins
`headless_args = ["-p", "{prompt}"]`, now redundant with the built-in. The roster-init half is
already false: the implementer reproduced the warned scenario (pre-existing `[roster.kimi-1]`,
one id unmapped, real `--yes` write) on **both** the installed 1.38.0 binary and the new build —
the mapping was appended, `kimi-1` was preserved, `antigravity-1` was not re-added. Refresh the
comment and drop the redundant pin.

## Dismissed

**`kimi-1` NIT-1 — surface `Scope`/`Confined()` in `agents list`.** Sound and worth doing, but it
changes the runtime matrix for every agent, not just the two promoted here. Recorded as a
follow-up rather than an agreed fix, so this idea does not grow a display change nobody asked for.

## Follow-ups

1. Surface `Confined()` in the runtime matrix (`kimi-1` NIT-1).
2. Probe well-known install prefixes at discovery time (the larger half of AF-4).
3. `~/.parley/agents.toml` carries further stale claims the implementer verified as wrong today:
   "no headless autonomous-write contract for kimi", "K3 default_effort = max", kimi 0.32.0,
   opencode 1.18.13. Machine-local; out of repo scope.

## Note on the review

All three reviewers declared the scope they checked, per §15.5. The `MAJOR` was found by exactly
one of them, and the other two — plus the implementer — reached the opposite conclusion. Under
§15.3 that was resolved by evidence, not by count: the config scan above is reproducible, and the
implementer's contrary probe was withdrawn once its environment was shown to be wrong.

## Signoffs

<!-- Each participant appends its own block. Do not edit another participant's block. -->
