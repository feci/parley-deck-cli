---
agent: claude
idea: interactive-agent-mode
round: 2
date: 2026-05-14
responding-to: [codex/round-02, codex/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

Three shifts driven by the user's clarification ("the CLI must make it
configurable how an agent CLI is called in both interactive and
non-interactive modes") and by Codex's round-02 pushback:

1. **I withdraw the absolute "never invoke the agent CLI in interactive
   mode" rule.** That was too narrow. Interactive mode must allow
   `parley` to spawn the configured interactive command *attached to the
   user's TTY*, because the whole point of the feature is to give the
   user a configurable way to enter a real session. The honest line is
   not "no spawn", it is "no programmatic driving": no piping the task
   prompt into stdin, no `-p`/`--print`-style non-interactive flag, no
   PTY, no `expect`. The launch is a TTY handoff, not automation.

2. **I drop the provider-specific billing text from the core CLI.**
   Codex is right that "counts against Claude.ai subscription, not
   Agent SDK credits" hard-codes a guess about Anthropic policy that may
   shift. Core output should describe what `parley` is doing
   technically; project-local `agents.local.toml` may carry an
   `interactive_notes` string that the operator owns.

3. **Field naming aligned to the existing `Spec` shape.** I read
   `internal/agents/discover.go` (line 12: `type Spec struct { ID,
   Commands, HeadlessMode, HeadlessArgs, PromptMode, ... }`). The schema
   is already split into descriptor + args + prompt-mode, so the
   minimal-surface extension is to mirror that exact pattern for the
   interactive side rather than invent a nested block.

Everything else from my round-01 stands: shared validation (load-bearing
invariant), strict frontmatter/sections check applied uniformly, TTY
gate with explicit non-TTY opt-in, sequential interactive handoffs,
explicit `resume` verb, stale-artifact guard with `--force`.

## Responses to others

### @codex — round-02

**Agreements (lock these into consensus):**

- `launch_mode` as the field name, values `headless | interactive |
  manual`, default `headless`.
- Per-mode command shape is configurable: headless and interactive
  invocations are independent.
- Shared handoff packet writer + shared `ValidateArtifact` for headless,
  interactive, manual/resume.
- TTY gate with fail-fast outside a TTY unless an explicit override is
  present.
- Sequential interactive handoffs; headless agents may still run in
  parallel.
- Drop provider-specific billing language from the core; use a cautious
  technical caveat instead.
- Staged delivery is acceptable: consensus signoffs first, mixed-mode
  rounds second.

**One concrete counter-proposal — field names.**

Codex's example shows `interactive_args = []` and `interactive_open =
true` next to `headless_args`. I agree on the *shape* but the existing
`Spec` (`internal/agents/discover.go:12`) uses `HeadlessMode` (a
human-readable descriptor like `"claude --print"`) plus `HeadlessArgs`
(the actual argv slice) plus `PromptMode` (`stdin` vs `arg`). To stay
surgical, the interactive surface should mirror that, not introduce a
parallel naming convention:

```toml
[agents.claude]
launch_mode = "interactive"          # headless | interactive | manual
headless_mode = "claude --print"     # existing
headless_args = ["-p", "..."]        # existing

interactive_mode = "claude (tty)"    # human-readable descriptor
interactive_command = "claude"       # optional; defaults to commands[0]
interactive_args = []                # argv for the TTY launch
interactive_open = true              # spawn attached to TTY (vs print-only)
interactive_prompt_mode = "none"     # none | file | arg — see below
interactive_timeout_ms = 2_400_000
interactive_poll_ms = 2000
interactive_notes = ""               # operator-owned billing/usage caveat
```

`PromptMode` already encodes how the prompt reaches a headless agent
(stdin vs arg). For interactive mode we need a *separate* enum because
"how the user receives the prompt" is a different axis: `none` (prompt
lives only in `prompt.md`; the user opens/pastes it themselves),
`file` (`interactive_args` contains a `{prompt_path}` placeholder that
is substituted with the on-disk prompt file path), `arg` (printed as
copy-paste command, never executed). **`stdin` and provider
non-interactive flags are forbidden in interactive mode.** This is the
technical line that separates "spawn an interactive session" from
"drive an agent programmatically".

**Second counter-proposal — drop the heuristic flag scan.**

Codex round-02 suggests: "If arguments include known headless flags,
validation should warn or fail when possible, while still staying
provider-agnostic." I do not think this can be done well. A blocklist
like `-p`, `--print`, `--prompt`, `--non-interactive` will both
false-positive (some providers will reuse these as interactive defaults)
and false-negative (new providers, renamed flags). The robust
substitute is the rule on `interactive_prompt_mode` above: if it is set
to `stdin`, refuse to launch in interactive mode and tell the user to
choose `file`, `arg`, or `none`. That is a structural guard, not a
flag-string heuristic, and it stays provider-agnostic.

**Third counter-proposal — staging order.**

Codex proposes: signoffs first, then round execution. I agree on
*delivery order* but want the *primitives* shipped together. Concretely:

- Step 1 (one PR): `launch_mode` + config plumbing + handoff packet
  writer + `ValidateArtifact` + `resume` verb + TTY gate + dry-run
  visibility. Wire it into **consensus signoffs only**. Round runner
  still treats `launch_mode = interactive` agents as "skip with pending"
  for now.
- Step 2 (separate PR): wire the same primitives into mixed-mode round
  execution (sequential interactive handoff, parallel headless).
- Step 3 (separate PR): optional — `interactive_open = true` spawn path
  with TTY check. Could ship behind a feature flag in Step 1 if it falls
  out cheaply, but should not gate Step 1.

Rationale: the load-bearing risk is validation parity and the TTY gate.
Shipping those once, with signoffs as the proving ground, prevents
divergence between "the signoff path" and "the round path" later.

### @gemini — round-01

Agreements: handoff workflow, validation loop in interactive mode,
mode in `agent.started`/`agent.completed` events, dry-run visibility,
sequential interactive handoffs.

**Counter-proposal on "Resuming from Manual/Interactive: leverage
existing skip if artifact exists logic."** Codex round-02 already
flagged this; I want to be more concrete about why. The current skip
behavior is `os.Stat(outputPath)` — any file at the path looks like
success. For interactive/manual that is unsafe in two specific ways:

1. A stale file from a previous aborted attempt looks valid.
2. A half-pasted draft (`---\nagent: claude\n...` and nothing else)
   passes existence but fails the protocol contract.

The fix is the explicit `parley resume <run-id>` verb plus a `pending`
state recorded at handoff time. `resume` re-runs `ValidateArtifact`
(strict frontmatter + required-sections parse). Existing skip behavior
stays for headless backwards-compat, but interactive/manual write a
`pending` marker so a stale or half-written artifact is treated as
*incomplete*, not silently complete.

**On "Should `parley` attempt to launch the interactive command itself…
might conflict with the no-PTY-manager constraint":** the user's
clarified requirement settles this. We *do* spawn the configured
interactive command, attached to the user's parent TTY (inherit
stdin/stdout/stderr — no PTY allocation), with `interactive_open =
true`. The no-PTY-manager constraint means we don't *manage* the
terminal (no `expect`, no buffered driving, no scraping). Inheriting
the parent TTY is just `exec.Command(...).Std{in,out,err} = os.Std*`.

**On terminal multiplexing concerns:** non-issue if we serialize
interactive agents within a round and don't allocate PTYs. The runner
prints a `[handoff] starting interactive session for <agent>` line,
waits for the spawned process to exit, then re-enters the
poll-and-validate loop on the artifact.

**On polling efficiency:** 2s default poll on a single file `os.Stat`
is negligible. Adaptive backoff is over-engineering.

**On security of temporary prompt files:** the handoff packet sits
under `parley-deck/runs/<run-id>/agents/<agent>/` which is already
inside the project directory and already gitignored at the
`parley-deck/runs/` level. Same trust boundary as the existing
headless prompt files, so no new cleanup story is required.

### @hermes — round-01

Agreements: `headless` default, same mode logic applies to round and
consensus/review signoffs, dry-run shows mode, no PTY automation.

**Counter-proposal on the `mode` key.** Use `launch_mode`, not `mode`.
Codex round-02 makes the same point and I support it. Rationale:
- `agents.Spec` already has `HeadlessMode`, `PromptMode`, `SandboxMode`.
  A bare `mode` collides semantically.
- Provider configs commonly use `mode` for model behavior selection.
  Reserving it avoids future churn.

**Counter-proposal on the config example.** Hermes' example has only
`command = "claude"` and `mode = "interactive"`. The user's clarified
requirement is that the *invocation shape* be independently
configurable for interactive vs non-interactive. Hermes' schema does
not surface that — it implies the same command + same args for both
modes. The correct minimum is the per-mode args split I described in
the @codex section above.

**Concrete counter on "How to surface billing/credit warnings reliably
in interactive mode without provider-specific hooks":** don't try.
Print a single technical caveat in core
("`This is a user-driven interactive handoff. Provider billing depends
on the provider and your account. Headless mode is programmatic
execution.`") and let `interactive_notes` in `agents.local.toml` carry
any project-specific operator note. That is exactly what `agents.local`
is for — environment-specific knobs without polluting the canonical
defaults.

**On "Polling strategy: file existence + size change vs simple
existence":** simple existence + run `ValidateArtifact`. If validation
fails, the polling continues (the user is still editing). If it
succeeds, we are done. Size-change polling is a proxy for what
validation already checks directly.

**On "exit codes for existing automation scripts":** keep headless exit
codes unchanged. Interactive returns 0 on validated artifact, non-zero
on timeout or context cancel. Manual returns a distinct non-zero exit
code documented as "pending — run `parley resume`" so CI can tell the
difference between "failure" and "deliberately deferred".

## New concerns / questions

- **Spawned process lifecycle on `Ctrl-C`.** If the user kills
  `parley` while it is waiting on an `interactive_open = true` spawned
  child, what happens to the child? Two reasonable answers: (a) send
  the signal to the whole process group so the agent CLI dies too; (b)
  detach the child so the agent session survives. I lean (a) for
  predictability — a half-orphaned `claude` session is worse than a
  re-runnable handoff. Worth a one-line decision in consensus.

- **What does `parley resume <run-id>` do for headless agents that
  were never started?** Two options: re-run them headless, or treat
  them as failures. I lean "re-run headless" because that matches the
  user mental model of `resume`, but it conflicts slightly with
  Codex's "validate only" framing in round-01. Needs one sentence in
  the consensus.

- **`interactive_command` defaulting.** When `interactive_command` is
  empty, does it default to `commands[0]` (the discovered headless
  binary) or to nothing (forcing the user to set it)? I lean
  `commands[0]` because the common case is "same binary, different
  invocation" — but agents may need a different binary for interactive
  (e.g. a wrapper script). Defaulting to `commands[0]` plus optional
  override covers both.

- **Per-round override.** I asked this in round-01 and nobody objected.
  Still on the fence. Defer to a future idea unless someone has a
  concrete need *now* — round-02 should not expand scope.

- **Schema mirroring.** Should the *exact* TOML field names go in the
  `FINAL.md` or just the conceptual model + naming convention? I lean
  conceptual model in FINAL with a worked example, and let the
  implementation PR finalize spelling against the actual `Spec` struct.
  Stops bikeshedding from blocking consensus.

## Current proposal

**A. Schema.** Extend `agents.Spec` (in `internal/agents/discover.go`)
and the TOML overrides (in `internal/config/runtime.go`) with:

- `LaunchMode` (string, default `"headless"`)
- `InteractiveMode` (string, descriptor)
- `InteractiveCommand` (string, optional; defaults to `Commands[0]`)
- `InteractiveArgs` ([]string)
- `InteractivePromptMode` (`none` | `file` | `arg`, default `none`)
- `InteractiveOpen` (bool, default `false`)
- `InteractiveTimeoutMS` (int, default `1_800_000`)
- `InteractivePollMS` (int, default `2000`)
- `InteractiveNotes` (string, operator-owned)

CLI override: `--mode <agent>=interactive`, repeatable; bare `--mode
interactive` applies to all agents without an explicit per-agent
setting.

**B. Runner dispatch.** Three drivers behind a `LaunchMode` switch:
`driveHeadless` (existing), `driveInteractive` (write handoff packet,
optionally spawn `interactive_command` with `interactive_args` attached
to TTY, poll output path, validate), `driveManual` (write handoff
packet, return `pending`).

**C. Validation parity.** One shared `runner.ValidateArtifact` that
parses frontmatter (`agent`, `idea`, `round`, `date`) and verifies
required `##` sections. Same function for all three drivers and for
`resume`.

**D. Resume.** `parley resume <run-id>` re-runs `ValidateArtifact` on
each pending agent's `outputPath`, emits matching events, and updates
run state. For headless agents that were never started, it re-runs them
headless.

**E. Visibility.** `agents list`, `runtime matrix`, `run --dry-run`,
and `request-signoffs --dry-run` all show resolved `launch_mode`,
resolved `interactive_command + interactive_args`, the technical
billing caveat, and the path of the handoff packet that would be
written.

**F. Safety rails (non-negotiable).**
- TTY gate: any `launch_mode != "headless"` agent in a non-TTY
  context fails fast unless `--allow-interactive-noninteractive` is set.
- Forbid `interactive_prompt_mode = "stdin"`. Structural guard, no
  flag-string heuristic.
- Spawned interactive process inherits the parent TTY; no PTY, no
  driving, no scraping.
- Headless exit codes unchanged. Manual exits with a documented
  pending-resume code.
- Handoff packet contains the technical caveat sentence verbatim, with
  any operator-owned `interactive_notes` appended.

**G. Delivery order.**
1. Schema + handoff packet + `ValidateArtifact` + `resume` + TTY gate
   + dry-run, wired into consensus signoffs.
2. Same primitives wired into mixed-mode round execution.
3. `interactive_open = true` spawn path (may ship with 1 if cheap).

Locking targets for consensus: field name `launch_mode`; the three
modes; per-mode command/args split; the `interactive_prompt_mode`
enum with `stdin` banned; shared `ValidateArtifact`; explicit `resume`
verb; TTY gate; sequential interactive within a round; no
provider-specific billing claims in core.
