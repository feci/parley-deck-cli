---
idea: zcode-adapter
status: final
drafted-by: claude-1
date: 2026-08-19
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1]
rounds: 2
signoffs: [claude-1 ✅, codex-1 🟡, hermes-1 🟡, kimi-1 🟡]
---

# FINAL — `zcode` becomes a built-in adapter; `unknown` is the honest MODEL

## Decision

Add `zcode` to the built-in adapter set, in the CLI and the skill.

**Launch:** `zcode --prompt {prompt} --mode yolo --cwd {root}`
**Autonomous write:** `AutonomousWrite{Mode: "yolo", Args: ["--mode","yolo"], Scope: ""}`
**MODEL and EFFORT: `unknown` / `model-unbound`** — unanimous, and not a gap to paper over.

`--mode yolo` is passed explicitly even though `zcode --help` calls it the default for `--prompt`,
because `~/.zcode/cli/config.json` holds `permission.mode = "build"` and the two disagree; the
effective launch argv is the source of truth. `Scope` is empty because `--cwd` is a working
directory, not an enforced sandbox — the honesty rule at `internal/agents/discover.go:86-92` reserves
`"workspace"` for a CLI that enforces one.

**This satisfies the owner's requirement that every roster member support a yolo-equivalent.**
`zcode-1` is the only AUTO=no member today and this change removes it. No protocol amendment: writing
the requirement into `COOPERATION.md` would bind nothing, and this deck has just spent two ideas on
the cost of rules nothing enforces.

## Why MODEL stays `unknown`

zcode has **no model flag** — `--model` does not appear in `zcode --help` at all. @kimi-1 probed the
alternatives to exhaustion: argv, `--settings`, env knobs across the 13 MB runtime bundle, and
`--json`. All dead. The model lives only in `~/.zcode/cli/config.json → model.main`.

A roster cell carrying that value would violate the frozen contract (`internal/app/roster.go:188-192`)
— MODEL is what the launch *passes*, or `unknown`, never a configured value argv does not carry. Every
other adapter that displays a model carries `{model}` in its `HeadlessArgs`; **zcode is the first that
cannot**.

## The limitation that justifies the successor

zcode is the only roster member whose participating model **cannot be established from outside**.
`--json` returns sessionId, traceId, turnId, usage and projection but no model id, and the config can
change between rounds including via `/model` in the TUI. After a zcode round, no artifact records
which model deliberated.

That is a **§15 auditability limitation**, not a roster-contract one — and it, not the MODEL column,
is why the successor exists.

**Named required successor (@codex-1):** `zcode app-server` exposes ZCode Protocol methods including
`session/create`, `session/setModel` and `session/setThoughtLevel`, with a model object carrying
`providerId` / `modelId` / optional `variant` plus a thought level. It is the only known way to bind
zcode's model and effort without editing the user's own config. Deferred because it is a different
product — a persistent session client, not a one-shot launch, not ACP — needing its own runner and its
own proof obligation for a non-argv binding. **A named required successor, not an undocumented ACP
shortcut.**

## The change set

**Code**
- `internal/agents/discover.go` — new `withBuiltinSources(Spec{ID: "zcode", …})`: `Commands:
  ["zcode"]`, `VersionArgs: ["--version"]`, `LaunchMode: LaunchHeadless`, `HeadlessArgs:
  ["--prompt","{prompt}","--mode","yolo","--cwd","{root}"]`, `PromptMode: PromptArg`, `Model:
  CLIDefault`, `Reasoning: CLIDefault`, `ExternalBackend: ExternalHosted`, `AutonomousWrite` as
  above, `Notes` recording: no model flag exists; the model resolves from the agent's own config;
  `--json` carries no model id; rejected flags exit 1.
- `internal/agents/modelmeta.go` — map the `zai/` prefix (`z-ai` and `glm` exist; `zai/` does not).

**Tests** — the contested claims, resolved by measurement and reproduced three times:
- `internal/agents/acp_specs_test.go:62` and `internal/app/roster_test.go:134` carry **presence
  lists, not exhaustive ones**; both pass on a tree with no zcode spec. Adding zcode there is a
  **deliberate lock**, not a fix. (@codex-1 established the `acp_specs_test.go` characterisation,
  @kimi-1 the `roster_test.go` one, @claude-1 verified both.)
- `internal/agents/launchargs_test.go` — **named by @codex-1 in round 1**; the natural home for a lock
  asserting zcode's argv carries no `{model}` placeholder and resolves to the exact verified argv.
- `internal/app/app_test.go` — **@codex-1's two fake-zcode full-verify cases**: one fake that prints
  help and exits 0 writing no sentinel (full verify must return nonzero), one that accepts exactly
  `--prompt <prompt> --mode yolo --cwd <root>` and writes the sentinel, asserting exact token order
  so a future spec cannot silently append an unsupported option.
- **No test asserts a fixed adapter count or list**, so adding an adapter cannot break a count.

**Machine config, same change**
- Delete the whole hand-written `[agents.zcode]` block from `~/.parley/agents.toml` after a
  clean-profile proof (@claude-1's deletion, @codex-1's gate). A redundant wholesale override is how
  hermes once silently lost `--yolo`.
- Remove the unbindable `model` / `effort` keys from `[roster.zcode-1]`, keeping membership and
  `speed` (@kimi-1).
- **Correct the stale exit-0 note at `~/.parley/agents.toml:117`** (@kimi-1's locator), which still
  claims a bad zcode launch "LOOKS like a success".

**Skill**
- Add a `zcode` row to the autonomous-write table in `skills/parley-deck/SKILL.md` (`--mode yolo`,
  scoped by `--cwd <deck>`).
- Make zcode's runtime a **native installer target** (@kimi-1): *"A first-class adapter whose runtime
  cannot receive the cooperation skill through the supported installer is not first-class."*

**`--explain` provenance.** An operator staring at `model unknown` needs the answer one command away;
spec `Notes` print only in the agents inventory views, not in `--explain`. **Phase-5 default: a
static trailer** naming `~/.zcode/cli/config.json → model.main` as the source, without reading the
live value. @codex-1 prefers a labelled live read and filed a reservation, not a block; @kimi-1's
argument decides it — *honesty was never the objection, staleness is* — and static ships first because
it loses nothing and stays upgradeable. **Any deviation must be recorded in `IMPLEMENTATION.md`.**

## Release

Minor bump for CLI and skill; all channels (npm, Homebrew, winget, GitHub) with independent
per-channel verification; the skill ships with the CLI as always.

**Acceptance command, corrected and measured** (`selectDiscoveries` matches discovery IDs, not roster
IDs — `internal/app/app.go:2103-2113`; mechanism established by @kimi-1 in round 1, exact pair
measured by @codex-1):

```
parley agents verify --full --agent zcode --yes
```

Today: `--agent zcode-1` → **rc=2**, `unknown agent zcode-1`; `--agent zcode` → **rc=0**,
`zcode: installed version=unknown`.

Remaining acceptance: `parley agents list` prints a `headless:` argv line for zcode (today it prints
none — that absence is how we know it is unsupported); `roster show` reports AUTO=yes and drops
`not-in-roster`; and **a real deck round is driven by parley rather than by hand**.

## Deferred

ZCode Protocol / app-server binding (required successor); generic exit-0-with-no-artifact diagnosis
(surface exit code, byte count and stderr tail when an expected artifact is absent — generic runner
tooling); any protocol rule requiring roster members to declare an autonomous-write mode (moot).

## Corrections against the facilitator

Recorded because a reader should discount the brief accordingly.

**Two errors were in `00-prompt.md` itself and were served to all three participants as fact.**
The exit-0 hazard does not exist — `--model`, `--settings`, `--max-turns` all **exit 1**; the original
measurement piped zcode into `head`, so `$?` was `head`'s status. @kimi-1 measured it correctly and
was right. And `kimi-1` is AUTO=yes, not AUTO=no — stale since parley 1.38, with v1.39.0 promoting
kimi to a full adapter. Q3 and Q4 both rested on false premises.

**Five attribution errors in the consensus draft**, all upheld at signoff, all verified before being
accepted: `launchargs_test.go` was named by @codex-1 in round 1, not missed by it; the family-ID
mechanism was @kimi-1's, not @codex-1's; the fake-zcode verify cases were @codex-1's, not @kimi-1's;
whole-block deletion was @claude-1's; roster key removal and the clean-profile gate were @kimi-1's.
Two of the five were flagged independently by both objectors, and two corrected credit *away from* the
agent making the correction.

## Recorded, not decided — `hermes-1` on `fireworks/inkling`

The roster changed mid-idea by owner instruction. Tally across this idea: **2 successes, 1 failure**.
The failure was `Model generated invalid tool call: bash` — 40 bytes, exit 0, no artifact — on a
prompt that succeeded on retry; a short tool-call probe passed on inkling and on the old `glm-5p2`.
The machine roster records the same model, previously under opencode, failing as "no answer, not
wrong answer", fixed 3/3 by write-first prompting. Whether this is a roster problem, a prompt-shape
problem or noise is not this idea's decision.

## Not authorised by this idea

Implementation and release. Phase 5 onward is the owner's call.
