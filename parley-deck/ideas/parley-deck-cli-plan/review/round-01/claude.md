---
agent: claude
idea: parley-deck-cli-plan
review-round: 1
date: 2026-05-10
reviewed-commit: uncommitted
---

## Summary

The slice delivers what `IMPLEMENTATION.md` claims for M0 (`go.mod`, `cmd/parley`, command dispatch, workspace init), the basic JSONL portion of M1 (`internal/store`), and the discovery/probe portion of M2 (`internal/agents`). The code is small, readable, and the focused tests for `protocol` and `store` pass. Compared against `FINAL.md`, however, three things stand out as worth pushing back on now, before the next slice locks them in:

1. `parley run` reports success and creates a "run" that is not a run — no idea, no kickoff prompt, no participants, no agent launch. This is more than "does not yet launch agents for real protocol rounds"; the command's whole semantics in `FINAL.md` are absent.
2. A 15-second `context.WithTimeout` wraps every command, including `run` and `tui`. It is harmless today and will silently break the moment agents are launched.
3. The `Initial feature set` in `FINAL.md` lists `resume` and `answer` as v1 commands; they are not implemented and are not listed under `Deviations from FINAL.md`.

Everything else is small. Adapter scope, status detail, error messages, TUI ergonomics.

## Findings

### [CRITICAL] `parley run` is a misleading no-op

`internal/app/app.go:123` (`runTask`) prints `Created run <id>` and writes a single `run.created` event to `parley-deck/runs/<id>/events.jsonl`. It does not create a new idea directory, does not write `00-prompt.md`, does not select participants, does not launch round 1, and does not record the chosen participants anywhere. Compare `FINAL.md:60-67`:

> `parley run "<task>"` — Creates a new idea and kickoff prompt. Selects participants. Starts round 1 in parallel. Opens the TUI by default unless `--no-tui` is passed.

Why this matters: the smoke run cited in `IMPLEMENTATION.md:53-56` produced `runs/20260510T194003Z/events.jsonl` with only `{"type":"run.created","data":{"mode":"auto","task":"smoke implementation run"}}`. There is no idea on disk to associate with that run, no participants in the event payload, and no way to `resume` it later. A user reading the success message would reasonably assume a run was started; nothing was started, and nothing recoverable was persisted beyond the bare task string and mode.

`IMPLEMENTATION.md:35` says "The implementation does not yet launch agents for real protocol rounds." That covers the *agent launch* gap, but it does not cover the *idea-creation* and *participant-selection* gap, which are independent and required to make the runtime event store usable.

Suggested fix: either
- (a) implement the idea/participant scaffolding portion of `parley run` (create idea slug from task, write `00-prompt.md` with frontmatter, copy participants from `agents discover` output or a flag, append `idea.created` and `participant.selected` events), or
- (b) downgrade the command to a clearly partial form: print `parley run is not yet implemented; M2 only supports agents discover` and exit non-zero, and remove the JSONL write so no orphan runs accumulate.

Option (b) is the surgical choice for this slice; option (a) is the next-slice work. Either is fine, but the current middle ground (silent partial success) is the worst option.

### [MAJOR] 15-second hard timeout wraps every command

`internal/app/app.go:28-29`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
```

This `ctx` is threaded into `runAgents`, `runTask`, `runTUI`, and `runTUIView`. Today it is harmless because the only consumer that respects it is `agents.probeVersion` (`internal/agents/discover.go:115`), which already imposes its own 4s child timeout. But:

- The next slice that wires `parley run` to actually launch Codex/Claude/Gemini will inherit this 15s ceiling and kill the agents mid-stream.
- `parley tui` will also inherit it. The TUI does not currently read `ctx`, so today it survives, but any future "live event tail driven by ctx" pattern will be cut off at 15s.

Why it matters: latent timeouts that "happen to be unreachable today" are exactly the kind of bug that surfaces as a vague "agent died after 15 seconds" in the next slice, with no obvious culprit because the timeout is at the dispatcher level, not the runner level.

Suggested fix: drop the dispatcher-level timeout. Push timeouts into the operations that actually need them — `agents.probeVersion` already does this correctly. `runTask` and `runTUI` should use `context.Background()` (or accept a cancellation signal from `os.Interrupt`).

### [MAJOR] `parley run` silently initializes a workspace

`internal/app/app.go:138-141` calls `protocol.InitWorkspace(*root)` from inside `runTask`. A user who runs `parley run "fix typo"` in the wrong directory (or in an unrelated repo) will get a `parley-deck/` tree, a `COOPERATION.md`, and a `runs/<id>/` folder created with no warning.

Why it matters: per the global `CLAUDE.md` ("Surgical Changes") and `FINAL.md:48-51` (where `parley init` is the explicit workspace-creation command), workspace creation should be opt-in. Auto-init in `run` makes the system harder to reason about and means typos against an uninitialized directory leave debris.

Suggested fix: in `runTask`, call `protocol.ReadWorkspaceStatus(*root)` first; if it returns `os.ErrNotExist`, print the same hint that `runTUIView` prints (`internal/app/app.go:177-178`: "no parley-deck workspace found; run `parley init` first") and return non-zero. Reserve `InitWorkspace` for the `init` command.

### [MAJOR] `parley status` is much thinner than `FINAL.md` describes, and the gap is not declared

`FINAL.md:68-71`:

> `parley status` — Prints current idea, phase, round, participant states, pending questions, and artifact completion. Works in CI and plain terminals.

`internal/app/app.go:96-121` (`runStatus`) prints transport, idea slug, idea status (the `status:` frontmatter field, which is "final" for the existing idea — i.e., a phase label, not the live runtime state), and participants. There is no phase, round, per-participant state, no pending questions, no artifact completion. `IMPLEMENTATION.md:34-41` ("Deviations from FINAL.md") does not list this gap.

Why it matters: `status` is the CI-friendly observability surface in `FINAL.md`. Reviewers (and CI consumers) reading "status command implemented" will assume more than was delivered.

Suggested fix: either
- (a) extend `runStatus` to compute phase/round/participant state from the `runs/*/events.jsonl` store and the protocol tree (which would also validate the event-store design end-to-end), or
- (b) add an explicit deviation to `IMPLEMENTATION.md` saying "`parley status` currently prints workspace transport and idea frontmatter only; runtime phase/round/participant state and HITL inbox are deferred to M3+." Option (b) is the smaller, honest move for this slice.

### [MAJOR] `parley resume` and `parley answer` are missing and undeclared

`FINAL.md:72-78` lists `parley resume <idea-or-run>` and `parley answer <question-id>` in the "Initial feature set." Neither command is wired in `internal/app/app.go` (the dispatcher in `Run` only handles `init`, `agents`, `status`, `run`, `tui`). `IMPLEMENTATION.md:34-41` does not list either as a deviation.

Why it matters: a reviewer comparing `FINAL.md` to the implementation needs the deviation list to be exhaustive. Today, the absence is invisible from `IMPLEMENTATION.md`.

Suggested fix: add both to the `Deviations from FINAL.md` section, ideally with a one-line note about which milestone they map to (`resume` ≈ M1/M5, `answer` ≈ M5). No code change required for this slice.

### [MINOR] `parley agents probe` alias is unimplemented

`FINAL.md:53-58` lists `parley agents discover / parley agents probe` as a single feature. `internal/app/app.go:85-94` only accepts `discover` and prints `usage: parley agents discover` for anything else. Either drop the alias from `FINAL.md` (out of scope for a review) or accept `probe` as a synonym in `runAgents` — a one-line change.

### [MINOR] Status error message is not actionable when the workspace is missing

`internal/app/app.go:104-108` returns `status failed: open <path>/parley-deck/COOPERATION.md: no such file or directory` when no workspace exists. `runTUIView` (same file, lines 175-183) already does the right thing: it detects `os.ErrNotExist` and prints `no parley-deck workspace found; run `parley init` first`. The status path should match.

Suggested fix: in `runStatus`, wrap the `ReadWorkspaceStatus` error with `errors.Is(err, os.ErrNotExist)` and print the same hint.

### [MINOR] TUI rendering will mis-render on narrow terminals and leaks into scrollback

`internal/tui/app.go:53-57` hardcodes `Width(46)`, `Width(58)`, `Width(108)`. On a sub-110-column terminal, the boxes will wrap or overflow. Also, `tea.NewProgram(model{...})` (line 28) is created without `tea.WithAltScreen()`, so the dashboard stays in the user's scrollback after `q`/Ctrl+C. `FINAL.md:148-164` shows an aspirational fixed layout, so the hardcoding is understandable for a shell, but both issues are cheap to fix now.

Suggested fix: switch to `tea.WithAltScreen()` and respond to `tea.WindowSizeMsg` to compute box widths from the actual terminal width.

### [MINOR] `DefaultSpecs` adds Hermes, which is not in FINAL.md

`internal/agents/discover.go:54-62` registers a Hermes spec with `Commands: []string{"hermes", "hermes-agent", "hermesagent"}` and a `command name is not standardized yet` note. `FINAL.md:97-99` lists `codex`, `claude`, `gemini` as the first adapters and nothing else. The global `CLAUDE.md` rule ("No features beyond what was asked") applies.

Suggested fix: drop Hermes for v1 unless there is a stated requirement, or move it behind an explicit `--include-experimental` flag. The current speculative entry will produce noise in `agents discover` output for users who do not care about Hermes.

### [MINOR] Run IDs collide on sub-second concurrent starts

`internal/store/events.go:25-27` uses `t.UTC().Format("20060102T150405Z")` — second resolution, no nonce. Two `parley run` invocations in the same second will share a directory and append to the same `events.jsonl`. Today this is "two writers, one file" (acceptable due to `O_APPEND`), but the *logical* run identity is lost.

Suggested fix: either add a short random suffix (`...Z-a3f1`) or include nanoseconds, e.g., `t.UTC().Format("20060102T150405.000000000Z")`.

### [MINOR] `IdeaStatus.Path` is dead state

`internal/protocol/workspace.go:26` declares `Path string` and `readIdeas` populates it (`workspace.go:119`), but no caller reads it. Either start using it (e.g., to surface the idea root in `status`) or remove the field. Per `CLAUDE.md`, dead orphans created by your own changes should be removed.

### [NIT] `Spec.Telemetry` is captured but never displayed

`internal/agents/discover.go:18` defines `Telemetry string` and `DefaultSpecs` populates it for every adapter. `PrintDiscovery` (lines 92-113) never prints it. Either drop the field or include it as a column / extra line.

### [NIT] Custom `splitLines` reinvents `bufio.Scanner`

`internal/store/events.go:73-86` walks the byte slice manually. `bufio.NewScanner(bytes.NewReader(data))` with the default `ScanLines` splitter is equivalent and the standard pattern.

### [NIT] `Store.Append` opens and closes the file on every event

`internal/store/events.go:29-50` performs `OpenFile` + `Write` + `Close` per call. Acceptable today; will become a measurable cost once the runner streams events at high rate. Consider keeping a long-lived `*os.File` handle on the `Store` once a real runner exists.

### [NIT] `printUsage` omits `version` and `help`

`internal/app/app.go:55-66` documents `init`, `agents discover`, `status`, `run`, `tui` but not `version`/`--version` or `help`/`-h`/`--help`, all of which `Run` accepts at lines 32-37. Add them to the printed usage block for symmetry.

## Open questions

1. Was the lack of idea-creation in `parley run` (CRITICAL above) intentional for this slice, or was it overlooked? If intentional, please add it to `Deviations from FINAL.md` so reviewers do not have to infer it.
2. Is the Hermes adapter a product requirement we missed in `FINAL.md`, or was it added speculatively because the binary was on disk? If the former, please update `FINAL.md`'s adapter list. If the latter, please drop it.
3. The 15-second dispatcher timeout (MAJOR) — was it a placeholder for a future cancellation channel, or a leftover from prototyping? Knowing the intent affects whether we replace it with `context.Background()` or with a real interrupt-driven `signal.NotifyContext`.
4. For `parley status`, do we want it to read from `runs/*/events.jsonl` to compute live phase/round (option (a) above), or stay at "frontmatter-only" until M3+ (option (b))? The choice affects whether the event-store schema needs to grow now or later.
