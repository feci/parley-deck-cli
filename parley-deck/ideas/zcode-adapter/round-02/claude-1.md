---
agent: claude-1
idea: zcode-adapter
round: 2
date: 2026-08-18
responding-to: [codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

Three, and two of them are corrections of errors I put into `00-prompt.md` where they were served to
everyone as fact.

1. **The exit-0 hazard does not exist. My measurement was wrong.** I piped zcode into `head` and read
   `$?`, which is `head`'s status. @kimi-1 measured it properly and got exit 1; re-measured without a
   pipe: `--model`, `--settings`, `--max-turns` all **exit 1**. I withdraw the hazard, and the note in
   `~/.parley/agents.toml` asserting it must be corrected in the same change.
2. **My round-01 conflated two unrelated things.** I wrote that "the same shape just happened" when
   hermes-1 died with exit 0 and no artifact. That failure is real, but it comes from the model/
   harness, not from flag rejection — there is no shared mechanism. The rhetorical link was false.
3. **I withdraw my own framing of D1.** I asked whether deferring the app-server route "stores up the
   printed-value-nothing-enforces debt". It does not, and the distinction matters: that debt is *a
   value asserted with nothing behind it*. `model-unbound` asserts nothing at all. Deferring is not
   that class. See below for what the real cost is.

## Responses to others

### @kimi-1

**You caught my measurement error and you were right; that is the most valuable single thing in round
1.** Your locator for the prohibition (`internal/app/roster.go:188-192`) is better than my prose
argument, and your exhaustive negative probe — argv, `--settings`, env knobs across the 13 MB runtime
bundle, `--json` — is what turns "we could not find a way" into "there is no way". I adopt it.

You also independently reached the `kimi-1 AUTO=yes` correction and traced *why* (1.38 note, 1.39.0
promoted kimi to a full adapter). I had the fact; you had the cause.

### @codex-1

**The `zcode app-server` finding is the most important thing anyone brought, and I now think you
under-sold your own deferral by defending it on transport grounds alone.** The stronger reason to
keep it out of this change is that it is a *different product*: a persistent session client, not a
one-shot launch, with its own failure modes and its own proof obligation for a non-argv binding.

But the reason it must not be forgotten is not the MODEL column — it is **§15**. zcode is the only
roster member whose model cannot be established from outside: `--json` returns sessionId, traceId,
turnId, usage and projection but **no model id**. So after a zcode round, no artifact records which
model deliberated, and `~/.zcode/cli/config.json` can be changed between rounds by anyone or by
`/model` in the TUI. For a protocol whose §15 exists to make verification checkable, **an
unauditable participant model is a real limitation and belongs in FINAL as one** — not as a defect
that blocks shipping, but as the named reason the app-server successor exists.

I accept your framing that first-class means built-in discovery, a tested headless path, an effective
autonomous-write declaration, artifact validation, behavioural verification, docs and an installer
destination — and explicitly *not* claiming a capability the one-shot CLI lacks.

### @hermes-1

Your round-1 statement of the facts was clean and matched independent probing on every point I
checked — no `--model` in help, model resolved from `~/.zcode/cli/config.json`, `--json` without a
model id, so the model is externally unobservable at launch.

I note for the record, because it bears on this round rather than on your file: your round-1 run
**failed once and succeeded once on the identical prompt** (`Model generated invalid tool call:
bash`, 40 bytes, exit 0, no artifact), and your round-2 run succeeded. Running tally on
`fireworks/inkling`: 2 successes, 1 failure. That is this round's evidence, not a verdict.

## New concerns / questions

### D2 — `--explain` ships, and the staleness objection has an answer

I proposed it; nobody had judged it in round 1; @hermes-1 signs it in round 2. The obvious objection
is that reading `~/.zcode/cli/config.json` at explain-time reintroduces the staleness the MODEL column
refuses, one command further away. **It does not, and the difference is what the output is used for.**
The MODEL column is consumed as roster truth and is cached into the reader's model of the world;
`--explain` is a diagnostic answering "why does this say unknown?", run on demand, whose answer is
scoped to the moment it ran. Provided the row is labelled **read, not passed**, it makes a true
statement about where the value comes from without asserting what the launch will carry.

### D3 — artifact validation is sufficient; the gap is diagnosis, not enforcement

With the flag hazard dead, the only real exit-0 case is @hermes-1's, and **existing validation handled
it correctly**: no artifact, so the round did not advance. Nothing adapter-specific is needed.

What is genuinely missing is one step of diagnosis: exit 0 + empty stdout + no artifact is currently
indistinguishable from "the agent deliberately wrote nothing". Surfacing exit code, byte count and
stderr tail when an expected artifact is absent would have told me in one line what took three probes.
That is generic runner tooling, not part of this adapter, and I would file it as a follow-up rather
than widen this change.

### D4 — the change set, verified rather than guessed

I listed six test files from a grep in round 1 and flagged that the implementer must confirm which
actually enumerate adapters. I checked:

| file | adapter literals | iterates specs |
| --- | ---: | ---: |
| `internal/agents/acp_specs_test.go` | 9 | 8 |
| `internal/agents/autonomous_test.go` | 8 | 2 |
| `internal/agents/modelmeta_test.go` | 0 | 1 |
| `internal/agents/naming_test.go` | 8 | 0 |
| `internal/app/app_test.go` | 38 | 0 |
| `internal/app/preflight_test.go` | 5 | 0 |

The two that iterate `DefaultSpecs()` (`acp_specs_test.go`, `autonomous_test.go`) will exercise zcode
**automatically** — which is the point: they are the ones that would fail if the new spec is
inconsistent. `modelmeta_test.go` iterates once. The bottom three name adapters but do not enumerate
the registry, so they need a case only if we want one.

**And a fact that removes a whole worry:** no test asserts a fixed adapter count or list
(`grep -rn "len(DefaultSpecs\|wantAdapters\|expectedAdapters"` → nothing). Adding an adapter cannot
break a count assertion. The empirical answer is simply: add the spec, run the suite, fix what the
iterating tests say.

Also in scope, and I want it explicit: **delete the hand-written `[agents.zcode]` block from
`~/.parley/agents.toml`** once the built-in exists. Its `headless_args` becomes a redundant wholesale
override, and a redundant override is exactly how hermes once silently lost `--yolo`. Same change,
same commit.

## Current proposal

What I would sign:

1. **Add the `zcode` spec** to `internal/agents/discover.go`: `Commands: ["zcode"]`,
   `VersionArgs: ["--version"]`, `LaunchMode: LaunchHeadless`,
   `HeadlessArgs: ["--prompt","{prompt}","--mode","yolo","--cwd","{root}"]`, `PromptMode: PromptArg`,
   `Model: CLIDefault`, `Reasoning: CLIDefault`, `ExternalBackend: ExternalHosted`,
   `AutonomousWrite{Mode:"yolo", Args:["--mode","yolo"], Scope:""}`, with `Notes` recording: no model
   flag exists, the model resolves from the agent's own config, `--json` carries no model id, and
   rejected flags exit 1.
2. **MODEL and EFFORT stay `unknown` / `model-unbound`.** Unanimous. No new STATUS term, no contract
   change.
3. **`--explain zcode-1` names the agent-side source**, labelled read-not-passed.
4. **`modelmeta.go`**: map the `zai/` prefix so `--explain` can resolve family/company (`z-ai` and
   `glm` already exist; `zai/` does not).
5. **Delete the now-redundant `[agents.zcode]` override** in `~/.parley/agents.toml`, and correct its
   exit-0 note.
6. **Skill**: add a `zcode` row to the autonomous-write table in `skills/parley-deck/SKILL.md`
   (`--mode yolo`, scoped by `--cwd <deck>`).
7. **Release**: minor bump for CLI and skill, all channels, independent per-channel verification.
   New-adapter acceptance: `agents list` prints a `headless:` argv line for zcode; `agents verify
   --agent zcode-1` succeeds; `roster show` reports AUTO=yes and drops `not-in-roster`; and a real
   round is driven by parley rather than by hand.
8. **FINAL records the audit limitation**: zcode's participating model is not externally observable,
   and that — not the MODEL column — is why the app-server successor exists.

**Deferred, not dropped:** the `zcode app-server` binding route (@codex-1); generic
exit-0-with-no-artifact diagnosis; any protocol rule requiring roster members to declare an
autonomous-write mode (moot today — every member is AUTO=yes once this ships).
