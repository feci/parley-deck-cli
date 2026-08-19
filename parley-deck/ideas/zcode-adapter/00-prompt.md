---
idea: zcode-adapter
author: claude-1
created: 2026-08-18
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: final
track: standard
---

## Problem / idea

**Make `zcode` a first-class Parley Deck adapter, in the CLI and in the skill, and release it.**

Owner instruction: *"ak ho parley nevie obsluzit, tak ten zcode samozrejme oprogramuj cez /parley-deck
a deployni novu verziu, zcode musime v parley deck podporovat, aj skill aj cli"* — and, on the
autonomous-write question: *"vsetci agenti ktory su v rostery by nieco take mali podporovat"*.

Today `zcode` is configured in `~/.parley/agents.toml` but parley cannot launch it: its adapter set is
fixed, there is no plugin path, `parley agents verify --agent zcode-1` returns "unknown agent", and
`roster show` reports `model-unbound, effort-unknown, metadata-unknown, AUTO=no`. It is a
manual-Bash participant. That must end.

**Output of this idea: a design that can be implemented in Phase 5**, not the implementation.

## Reconnaissance — all PRIMARY, verified 2026-08-18 on this machine, re-runnable

**The binary.** `/opt/homebrew/bin/zcode` → npm global `zcode-app-cli 3.7.7-13` / `zcode-runtime
0.16.3`. Vendor Z.AI (Anthropic-kind provider, `https://api.z.ai/api/anthropic`).

**Autonomous write — this is the answer to the owner's second question.**
`zcode --help` documents `--mode <mode>`: *"Permission mode for prompts: build, edit, plan, or yolo
(default: yolo for --prompt)"*, plus a legacy alias `--permission-mode <mode>` (default, build, edit,
plan, yolo). Confinement is `--cwd <path>`. There are also `--allowed-tools` / `--disallowed-tools`
for headless prompts.

**The documented default MUST NOT be relied on.** `~/.zcode/cli/config.json` holds
`permission.mode = "build"` (plus `autoApproveHighRisk: false`, `allowMediumRiskInAuto: false`), so
config and help disagree. Per the skill's rule that the effective launch argv is the source of truth,
`--mode yolo` is passed explicitly.

**Verified launch shape** (wrote its own file unattended in a scratch dir, exit 0):
```
zcode --prompt "<text>" --mode yolo --cwd <root>
```

**There is NO model flag.** `--model` does not appear in `zcode --help` at all. (A note in
`~/.parley/agents.toml` claims `--model`, `--settings` and `--max-turns` are all "listed in help and
rejected" — that is wrong about `--model`: it is not listed. `--settings` and `--max-turns` *are*
listed; whether they function was not re-tested today.) A rejected flag prints help and **exits 0**,
so a bad launch looks like a success.

**Consequence:** the model comes from `~/.zcode/cli/config.json` → `model.main` (currently
`zai/glm-5.3`; lite `zai/glm-5-turbo`). `--json` returns sessionId/traceId/turnId/usage/projection
and confirms contextWindow 1000000, but **does not report the model id** — so "which model answered"
is not verifiable from outside.

**The code this lands in.** Adapters are `withBuiltinSources(Spec{ID: "...", ...})` entries in
`internal/agents/discover.go`. The hermes spec (`:293-316`) is the closest structural analogue:

```go
ID: "hermes", Commands: []string{"hermes", ...}, VersionArgs: []string{"--version"},
LaunchMode: LaunchHeadless, HeadlessMode: "hermes --yolo --oneshot",
HeadlessArgs: []string{"--yolo", "--oneshot", "{prompt}", "--model", "{model}", "--accept-hooks"},
PromptMode: PromptArg, Model: "xai/grok-4.3", Reasoning: CLIDefault, ...
AutonomousWrite: AutonomousWrite{Mode: "yolo", Args: []string{"--yolo"}, Scope: ""},
```

**`AutonomousWrite` carries an explicit honesty rule** (`discover.go:86-92`), and it is binding here:

> Scope names the DEMONSTRABLE confinement and is set to "workspace" ONLY where the CLI enforces a
> real workspace sandbox (codex `--sandbox workspace-write`). For CLIs whose only confinement is a
> grant/bypass flag + cwd (claude `--add-dir`, hermes `--yolo`), Scope is left EMPTY — the honesty
> rule: never falsely assert workspace confinement that is not enforced.

**Model metadata.** `internal/agents/modelmeta.go` maps ids to MODEL-FAMILY / MODEL-COMPANY. It
already has `"z-ai": "Zhipu AI"` (`:43`) and `{"glm", "GLM", "Zhipu AI"}` (`:63`). zcode's id prefix
is `zai/` (no hyphen), which is why `roster show` reports `metadata-unknown`.

**Tests that enumerate adapters and will need extending:** `internal/agents/autonomous_test.go`,
`modelmeta_test.go`, `acp_specs_test.go`, `naming_test.go`, `internal/app/app_test.go`,
`preflight_test.go`.

## The questions this idea must answer

**Q1 — What does the MODEL column report, and is that honest?** The roster contract is explicit:
MODEL and EFFORT are *what the launch actually passes*, or `unknown` — never a configured value the
argv does not carry. zcode has no model flag, so nothing can be passed. Options include: leave
`Model: CLIDefault` and accept `model-unbound` forever; have the adapter *read*
`~/.zcode/cli/config.json` and report it with a status that says it was read, not passed; or add a
new status vocabulary term. **Whatever is chosen must not make the roster lie.** Note the STATUS
vocabulary is closed — adding a term is a contract change, so say so if you propose one.

**Q2 — `AutonomousWrite` for zcode.** Mode `yolo`, Args `["--mode", "yolo"]`, Scope — `""` or
`"workspace"`? `--cwd` is a working directory, not an enforced sandbox. Argue it against the honesty
rule quoted above, and check whether zcode enforces anything beyond cwd (`--allowed-tools`,
`--disallowed-tools`, `autoApproveHighRisk`, `allowMediumRiskInAuto`).

**Q3 — The owner's generalisation: "every agent in the roster should support something like yolo".**
Is that a statement about *this* adapter, or a rule the deck should adopt — e.g. an invariant that a
roster member must declare an autonomous-write mode, with `roster show` flagging members that do
not? Note kimi-1 currently reports AUTO=no and is a manual participant, so such a rule would have a
consequence today. If you propose it, it is a §7 protocol change and a **successor** idea, not this one.

**Q4 — The exit-0-on-bad-flag hazard.** A rejected zcode flag prints help and exits 0. A launch that
silently degraded to "printed help" would look like a successful round with an empty artifact. Does
the adapter need a post-launch assertion (e.g. artifact-exists is already checked — is that enough?),
and does `parley agents verify` need anything zcode-specific?

**Q5 — Scope of the change.** Enumerate exactly what must change: `discover.go`, `modelmeta.go`,
which tests, the skill's `SKILL.md` autonomous-write table (it lists claude/codex/hermes/agy/kimi/
opencode today), README, and anything else. State what you verified rather than assumed.

**Q6 — Release.** A CLI change ships with a skill release across every channel (npm, Homebrew,
winget, GitHub) with independent per-channel verification. Say what the version bump should be and
what the release checklist must include for a *new adapter* specifically.

## Constraints

- **Not a protocol change.** `COOPERATION.md` is untouched by this idea. If your answer needs one,
  propose it as a successor.
- Provenance (§15.2): PRIMARY locators or the command you ran, or a `RECALL` tag. **Do not assert a
  negative about this repo or about zcode without running the check** — and note that a zcode flag
  that "works" may simply have been ignored, because rejection exits 0.
- Do not invent zcode capabilities. Probe them. `zcode --help`, `zcode doctor`, `~/.zcode/`.
- ⚠ **Secrets:** `~/.hermes/config.yaml` and the zcode/LiteLLM configs contain live API keys. Do not
  print, quote, or copy them into any artifact. Redact if a command would surface one.
- The roster changed today: participants are claude-1, codex-1, **hermes-1 (now `fireworks/inkling`,
  not GLM 5.2)** and kimi-1. opencode-1 was removed.

## Non-goals

- Implementing it. This idea produces the design; Phase 5 implements it.
- Changing which model zcode runs, or fixing zcode itself.
- Re-litigating whether zcode belongs in the roster — the owner has decided it does.
