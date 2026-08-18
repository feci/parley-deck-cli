---
agent: claude-1
idea: zcode-adapter
round: 1
date: 2026-08-18
---

## Summary

A zcode adapter is straightforward everywhere except Q1, and Q1 has an answer that needs **no
contract change and no new STATUS term** — it uses a mechanism the CLI already ships.

I also correct a factual error in my own `00-prompt.md` (Q3's premise is false), and record an
orchestration failure that happened inside this very round.

## Proposed approach

### Q1 — MODEL stays `unknown`. The resolved value goes to `--explain`, not to the column.

**The constraint, verified.** Every existing adapter whose MODEL column shows a value carries
`{model}` in its `HeadlessArgs`. I checked the one I expected to be a counter-example and it is not:
kimi's spec is `HeadlessArgs: []string{"-m", "{model}", "-p", "{prompt}"}` with `Model: CLIDefault`
— the value from `[agents.kimi]` substitutes into `{model}`, so argv genuinely carries it.
**zcode would be the first adapter that cannot pass a model at all**, because `--model` does not
exist in `zcode --help`.

**Why "read the config and display it" is wrong.** It looks honest — we would be reading the very
file the process will use — but it fails on *time*. `~/.zcode/cli/config.json` can change between
discovery and launch, and `/model` in the TUI rewrites it. A value read at discovery and printed as
MODEL is exactly the "config says X, process runs Y" failure the contract exists to prevent; we would
have reproduced the claude Opus-4.8 defect with a different mechanism. The MODEL column answers one
question — *what will the launch pass?* — and for zcode the true answer is "nothing".

**So MODEL stays `unknown` / `model-unbound`, and that is a feature.** What is missing is not a
value in that column but *discoverability of where the model actually comes from*. The CLI already
has the right home for that: **`parley roster show --explain AGENT`**, documented as "per-field
provenance: which config layer set each value".

**Proposal:** `--explain zcode-1` reports, for the MODEL field, that the value is resolved by the
agent's own configuration rather than passed by parley — naming `~/.zcode/cli/config.json` →
`model.main` and its current value, explicitly labelled **read, not passed**. Zero contract change,
no new STATUS term, no lie in the frozen columns, and the operator gets the answer one command away.

**Cost of being wrong:** an operator reads `unknown` in the table and concludes the adapter is broken
rather than that the CLI cannot pass a model. Mitigated by the STATUS term `model-unbound` already
being in the closed vocabulary and by the `--explain` text saying so in words.

### Q2 — `AutonomousWrite{Mode: "yolo", Args: []string{"--mode", "yolo"}, Scope: ""}`

`Scope` must be **empty**, and the honesty rule at `discover.go:86-92` decides it: `"workspace"` is
reserved for a CLI that *enforces* a real sandbox (codex `--sandbox workspace-write`). zcode's only
confinement is `--cwd`, which is a working directory. Nothing in `zcode --help` enforces a boundary:
`--allowed-tools` / `--disallowed-tools` are opt-in allow/deny lists, and `autoApproveHighRisk:
false` / `allowMediumRiskInAuto: false` live in config, not argv. Same standing as claude, hermes and
kimi — all `Scope: ""`.

`--mode yolo` is passed **explicitly** even though help says it is the default for `--prompt`, because
`~/.zcode/cli/config.json` holds `permission.mode = "build"` and the two disagree. The skill's rule
is that the effective launch argv is the source of truth, and `AutonomousWrite.MissingFrom` exists
precisely to check the enabling flag is present in argv — a declared mode whose flag is absent is the
failure it guards.

### Q3 — Correcting my own prompt: the premise is false

`00-prompt.md` states "kimi-1 currently reports AUTO=no, so such a rule would have a consequence
today". **That is wrong.** `parley roster show` reports `kimi-1 ... yes` in the AUTO column. All four
current deck members are AUTO=yes; **zcode-1 is the only AUTO=no, and this idea removes it.**

So the owner's generalisation — *every agent in the roster should support something like yolo* —
would be **satisfied by this change alone**, with no protocol amendment needed. If the deck still
wants it written down as an invariant (a roster member MUST declare an autonomous-write mode), that
is a §7 change and belongs in a successor. I would not open it: it currently binds nothing, and this
deck has just spent two ideas on the cost of rules nothing enforces.

### Q4 — The exit-0 hazard is real, and it fired inside this round

A rejected zcode flag prints help and exits 0. **The same shape just happened for a different
reason:** hermes-1's first attempt at this round's prompt died with `Model generated invalid tool
call: bash`, produced 40 bytes and **exit 0**, and wrote no artifact.

That is the evidence for the answer: **artifact validation already catches it.** The round did not
advance, because the facilitator checks the file exists. Exit code was useless; the artifact check was
decisive. So no zcode-specific machinery is needed for correctness.

What is worth adding is diagnosis, not enforcement: `parley agents verify --agent zcode-1` should do
what verify does for other adapters — a real minimal write into a scratch dir, asserting the file
appears — because for zcode, exit 0 proves nothing at all. Verify must never treat exit 0 as success
for this adapter.

### Q5 — Scope of the change

- `internal/agents/discover.go` — new `withBuiltinSources(Spec{ID: "zcode", ...})`:
  `Commands: ["zcode"]`, `VersionArgs: ["--version"]`, `LaunchMode: LaunchHeadless`,
  `HeadlessArgs: ["--prompt", "{prompt}", "--mode", "yolo", "--cwd", "{root}"]`,
  `PromptMode: PromptArg`, `Model: CLIDefault`, `Reasoning: CLIDefault`, `ExternalBackend:
  ExternalHosted`, `AutonomousWrite` as in Q2, `Notes` recording the no-model-flag and exit-0 facts.
- `internal/agents/modelmeta.go` — the `zai/` prefix is unmapped (`z-ai` and `glm` are present).
  Worth adding for `--explain` output even though the MODEL column stays `unknown`.
- Tests: `internal/agents/autonomous_test.go`, `modelmeta_test.go`, `acp_specs_test.go`,
  `naming_test.go`, `internal/app/app_test.go`, `preflight_test.go` — I listed these from a grep for
  existing adapter names; the implementer must confirm which actually enumerate adapters rather than
  merely mention one.
- Skill: the autonomous-write table in `skills/parley-deck/SKILL.md` lists claude, codex, hermes, agy,
  kimi, opencode — it needs a `zcode` row: `--mode yolo` (+ `--cwd <deck>` to scope it).
- `~/.parley/agents.toml` — once the built-in adapter exists, the hand-written `[agents.zcode]`
  block's `headless_args` becomes a redundant wholesale override, and a redundant override is exactly
  how hermes once silently lost `--yolo`. It should be deleted in the same change.

### Q6 — Release

Minor bump: new adapter, backward compatible, nothing removed. CLI `1.45.0` + a skill minor. All
channels (npm, Homebrew, winget, GitHub) with independent per-channel verification, and the skill
ships with the CLI as always.

**New-adapter-specific checklist items** beyond the usual: `parley agents list` shows a `headless:`
argv line for zcode (today it prints none — that absence is how we know it is unsupported);
`parley agents verify --agent zcode-1` succeeds instead of "unknown agent"; `roster show` drops
`not-in-roster`/`unknown agent` and reports AUTO=yes; and a real deck round is driven by parley
rather than by hand.

## Concerns / open questions

1. **Is `--explain` actually able to carry a per-field note today, or does it only report config
   layers?** I proposed it from its documented description without reading its implementation. If it
   is layer-only, Q1's answer needs a different home and I would rather know that in round 2 than at
   review.
2. **`PromptMode`.** zcode takes `--prompt <text>` and the prompt is *not* required to be last
   (`--mode` and `--cwd` follow it in the verified shape). I set `PromptArg`, but the enum's exact
   semantics should be checked against how the runner substitutes `{prompt}`.
3. **Does an unknown-adapter block in `~/.parley/agents.toml` currently do anything at all?** If
   `[agents.zcode]` was inert, then adding the built-in adapter changes behaviour in ways nobody has
   observed yet, and the first driven round is the real test.

## Risks

- **hermes-1 is newly on `fireworks/inkling` and failed once in this round.** 1 failure / 1 success on
  the same prompt; a short tool-call probe passed, and the same probe passed on the old `glm-5p2`. So
  it is intermittent, not categorical — but a participant that silently produces nothing on exit 0 is
  the worst failure shape for this protocol, and the roster note for opencode records the *same
  model* previously failing as "no answer, not wrong answer", fixed by write-first prompting. Watch it
  this round; if it recurs, that is a roster decision, not a zcode one.
- **I wrote the prompt and got a fact in it wrong** (Q3). Other premises in it may be wrong too;
  participants should check rather than inherit.
- **Scope creep into a protocol change.** Q3 invites one. It should stay closed unless someone shows
  a case it would have caught.
