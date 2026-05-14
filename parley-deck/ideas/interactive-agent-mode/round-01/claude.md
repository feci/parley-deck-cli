---
agent: claude
idea: interactive-agent-mode
round: 1
date: 2026-05-14
---

## Summary

Add a `launch_mode` field per agent (`headless` | `interactive` | `manual`)
with `headless` as the default for full backwards compatibility. The runner
gains a small mode dispatcher that, for `interactive`/`manual`, writes a
self-contained handoff packet (prompt + target path + validation contract +
billing/usage note), then either polls for the artifact (interactive) or exits
with next-step instructions (manual). Validation of the final artifact uses
the *same* code path as headless, so the protocol guarantees do not weaken.
No PTY, no session driving, no fake authorship.

## Proposed approach

**1. Configuration surface (smallest viable change)**

Extend `agentOverride` in `internal/config/runtime.go` and `agents.Spec` with:

- `launch_mode` (string, default `headless`)
- `interactive_timeout_ms` (int, default 1_800_000; distinct from headless
  `timeout_ms` because human latency is different from model latency)
- `interactive_poll_ms` (int, default 2000)
- `interactive_open_command` (optional string, e.g. `"claude"`; if set and a
  TTY is attached, the runner *suggests* launching it but never pipes a
  prompt into it)

`parley-deck/agents.toml` change for `claude` would become:

```
[agents.claude]
launch_mode = "interactive"
interactive_timeout_ms = 2400000
```

CLI override on `parley run` / `parley consensus`:
`--mode <agent>=interactive` (repeatable). A bare `--mode interactive`
applies to all agents lacking an explicit per-agent setting.

**2. Runner dispatch (`internal/runner/runner.go`)**

Refactor `runAgent` to delegate after prompt construction and pre-flight
(stdout/stderr files, `agent.started` event). Three drivers:

- `driveHeadless` — current `cmd.Run()` path, unchanged.
- `driveInteractive` — writes `runs/<id>/agents/<agent>/prompt.md` and
  `handoff.txt` (human-readable: target path, validation contract,
  billing/usage note, "open your CLI now" instructions). Emits
  `agent.handoff.started`. Polls `outputPath` every `interactive_poll_ms`
  until it exists *and* passes validation, or `interactive_timeout_ms`
  elapses, or context is cancelled. Emits `agent.handoff.completed` /
  `agent.handoff.timeout`.
- `driveManual` — same packet write as interactive, then returns a
  `Skipped`-like result with `SkipReason = "manual handoff pending"`. Emits
  `agent.handoff.pending`. The round event becomes `round.incomplete` with
  `pending_manual` count surfaced.

**3. Validation parity**

Extract the current artifact check (`os.Stat(outputPath)`) plus a new
frontmatter/sections sanity check into `runner.ValidateArtifact(outputPath,
agent.ID, idea.Slug, round)`. Both `driveHeadless` and `driveInteractive`
call the same function. This is the single load-bearing invariant: the
mode only changes *who runs the agent*, never *what counts as a valid
artifact*.

**4. Resume path (`parley resume <run-id>`)**

For `manual` mode (and for interactive runs the user closed early), add
`parley resume`: re-loads the run, re-checks each pending agent's artifact,
runs validation, and writes the matching `agent.finished` event. This makes
`manual` a first-class flow rather than a dead-end — the same mechanism
also lets a user recover from an aborted interactive session without re-
running anything.

**5. Visibility (dry-run, launch, TUI)**

- `parley agents` / runtime matrix: add `MODE` column and a
  `billing-note` line per agent (e.g., `claude.interactive: counts against
  Claude.ai subscription, not Agent SDK credits`).
- `parley run --dry-run`: print, per agent, the resolved mode + billing
  note + the path of the handoff packet that would be written.
- TUI `runstate`: surface `waiting (interactive)` and `pending (manual)`
  states distinct from `running`.

**6. Concurrency**

Headless agents continue to run in parallel. Interactive agents are
serialized within a single `RunRoundOne` call (one open handoff at a
time) — driving two terminals at once is a footgun, not a feature. Manual
agents do not block; they finish immediately with `pending_manual`.

**7. Safety rails (non-negotiable)**

- If `launch_mode != headless` and stdin/stdout is not a TTY *and*
  `--allow-interactive-noninteractive` is not set, fail fast with a clear
  error. CI must never silently hang waiting for a human.
- The handoff packet must contain the verbatim sentence:
  `This handoff is interactive. Provider billing applies to your
  interactive session, not to parley's automation budget.` This is the
  audit trail that the mode was explicit.
- Never invoke the agent CLI in `interactive` or `manual` mode — only
  *suggest* a command line. No PTY, no `expect`, no stdin pipe.

## Concerns / open questions

- **Stale-artifact detection.** Current `runAgent` skips when the
  `outputPath` already exists. For interactive mode this hides the case
  where a previous aborted attempt left a stale file. Proposal: in
  interactive/manual mode, refuse to start if the file exists and has
  mtime older than the run's `StartedAt`; print the path and require
  `--force` to overwrite expectations. Want sign-off on this rule.
- **Validation strictness.** Right now headless considers any file at
  the path "OK". Interactive opens a much wider door for a user pasting
  half a draft. Should we promote validation to actually parse the
  frontmatter (`agent`, `idea`, `round`, `date`) and verify the four
  required `##` sections exist? I lean yes — and apply it uniformly so
  headless gets the same guard for free. Worth confirming.
- **Consensus signoffs.** Signoffs have a different filename pattern and
  may live in `consensus/` rather than `round-NN/`. The mode dispatcher
  needs the same generality there — proposing we thread `outputPath`
  through `consensus.go` the same way the runner does, so one
  `driveInteractive` works for both. Has anyone already started this?
- **Per-round vs per-agent override.** Should `launch_mode` also be
  configurable per-round (e.g., headless for round-01, interactive for
  consensus signoff)? Adds complexity; would only do it if there is a
  concrete need.
- **Resume UX.** Is `parley resume <run-id>` the right verb, or should
  `parley run` itself idempotently re-validate pending agents on a
  second invocation? I prefer the explicit verb for auditability.
- **Open-command field.** Worth supporting at all? It is borderline — it
  blurs the line between "we don't drive the session" and "we launched
  it for you." If we keep it, it must spawn a *detached* terminal with
  no stdin pipe, and the prompt must be the user's responsibility to
  paste/load. Could also drop it entirely and just print the suggested
  command.

## Risks

- **Billing misclassification.** Anyone marking `claude` as `headless`
  in `agents.local.toml` after June 15, 2026 will draw from the
  programmatic budget without a fresh visual cue. Mitigation: print the
  mode and billing note at the *top* of every `parley run` output, not
  only in `--dry-run`, so the choice is always visible.
- **CI hang.** An interactive mode leaking into CI would block forever.
  Mitigation: TTY check + explicit opt-in flag, and emit
  `run.failed` immediately if the gate trips.
- **Validation gap.** If we keep the current "any file is OK" check,
  interactive mode will accept malformed artifacts more often than
  headless does (because humans are messier than agents). Mitigation:
  ship stricter validation in the same change so both modes benefit.
- **Concurrency confusion.** Two interactive agents racing for the same
  user's attention will produce one valid artifact and one timeout.
  Mitigation: serialize interactive agents and surface a clear
  `waiting on <agent>` status.
- **Stale-artifact false positives.** Re-running a round with an
  existing artifact will silently "succeed" today. Interactive mode
  amplifies the surprise. Mitigation: the mtime/`--force` rule above.
- **Drift between headless and interactive prompts.** If
  `BuildRoundOnePrompt` evolves only for headless callers, the handoff
  packet will fall out of sync. Mitigation: both drivers consume the
  exact same `prompt` string and the packet writer adds only the
  human-facing wrapper around it.
