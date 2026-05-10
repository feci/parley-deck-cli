---
agent: claude
idea: parley-deck-cli-plan
review-round: 2
date: 2026-05-10
reviewing: implementation-fix-up-cycle-2
---

## Summary

Fix-up cycle 2 lands the first real `runner` layer and wires `parley run` end-to-end: workspace check → idea + `00-prompt.md` → run event → parallel agent launch with stdout/stderr capture and per-round artifact validation. The core round-01 review findings are addressed (no-op `run`, dispatcher-level 15s timeout, silent workspace init, undeclared `resume`/`answer`, alt-screen TUI, nanosecond run IDs, `IdeaStatus.Path` now used, telemetry rendered, `bufio.Scanner`, etc.). The new `runner_test.go` exercises the headless contract through a fake child process without spending model tokens — that is the right shape of test for this layer.

The new code introduces three concerns I would treat as MAJOR before cycle 3 closes:

1. `parley run` always launches **every installed agent** with no `--participants` flag, no interactive selection, and no confirmation gate. On a machine with codex+claude+gemini+hermes this fans out to four parallel paid agent runs per command, contradicting the "Selects participants" feature in `FINAL.md` and the HITL default.
2. The runner swallows every `store.Store.Append` error with `_ =`, so a failing event write is invisible — defeating the durability rationale of the JSONL store.
3. Early-return paths in `runAgent` (mkdir / prompt build failures) update the in-memory `Result` but emit no `agent.*` event at all, leaving `events.jsonl` with a `round.incomplete` roll-up and no per-agent record. That is the exact recovery gap the event log is supposed to prevent.

Nothing here is CRITICAL: the test passes, the runner's headless contract is sound, cross-builds work, and the structural FINAL.md gaps from round-01 are closed or honestly declared. But cycle 3 should not move on to live TUI streaming or token accounting before the participant-selection and event-durability gaps are tightened.

## Previous findings status

Round-01 findings, cycle-2 status:

- **[CRITICAL] `parley run` is a misleading no-op** — **Fixed**. Cycle 1 added idea/prompt/participants scaffolding; cycle 2 added real agent launch. `app.go:139-216` now requires a workspace, creates the idea via `protocol.CreateIdea`, records `run.created` with `idea` + `participants` in the payload, and invokes `runner.RunRoundOne`.
- **[MAJOR] 15-second hard timeout wraps every command** — **Fixed**. `app.go:31` replaces the dispatcher timeout with `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`. Per-operation timeouts live where they belong: `agents.probeVersion` (4s) and `runner.runAgent` (`opts.Timeout`, default 30m).
- **[MAJOR] `parley run` silently initializes a workspace** — **Fixed**. `app.go:154-161` now calls `protocol.ReadWorkspaceStatus` first and returns the same `run \`parley init\` first` hint on `os.ErrNotExist`.
- **[MAJOR] `parley status` thinner than `FINAL.md`** — **Documented**, not extended. Listed in `IMPLEMENTATION.md` cycle-1 deviations as "still frontmatter-oriented and does not yet compute full runtime phase, round, participant state, HITL inbox, or artifact completion from event logs." Acceptable for this slice; the gap is now visible to reviewers.
- **[MAJOR] `parley resume` and `parley answer` missing and undeclared** — **Fixed**. Both are explicit dispatcher stubs in `app.go:49-54` returning exit 1 with a clear "not implemented yet" message, and listed under cycle-1 deviations.
- **[MINOR] `parley agents probe` alias unimplemented** — **Fixed**. `app.go:98` accepts both `discover` and `probe`.
- **[MINOR] Status error message not actionable when workspace missing** — **Fixed**. `app.go:117-121` uses `errors.Is(err, os.ErrNotExist)` and prints the same hint the TUI uses.
- **[MINOR] TUI mis-renders on narrow terminals; leaks into scrollback** — **Fixed**. `tui/app.go:29` uses `tea.WithAltScreen()`; `tui/app.go:40-72` handles `tea.WindowSizeMsg` and computes box widths from terminal width with sensible 80-column fallback.
- **[MINOR] `DefaultSpecs` includes Hermes, not in `FINAL.md`** — **Not fixed**. `internal/agents/discover.go:73-81` still ships Hermes. Also not added to `FINAL.md`'s adapter list, so the divergence from the agreed roster persists. Re-flagged below.
- **[MINOR] Run IDs collide on sub-second concurrent starts** — **Fixed**. `store/events.go:28` uses `20060102T150405.000000000Z`.
- **[MINOR] `IdeaStatus.Path` is dead state** — **Fixed**. Used by `runner.go:112` (`outputPath`) and `runner.go:224` (read `00-prompt.md`).
- **[NIT] `Spec.Telemetry` never displayed** — **Fixed**. `agents/discover.go:130` prints `telemetry: %s`.
- **[NIT] Custom `splitLines` reinvents `bufio.Scanner`** — **Fixed**. `store/events.go:60` uses `bufio.NewScanner(bytes.NewReader(data))`.
- **[NIT] `Store.Append` opens/closes the file on every event** — **Not fixed**. Acceptable for now; flagged again below for the new concurrent-writer context.
- **[NIT] `printUsage` omits `version`/`help`** — **Partially fixed**. `version` is now in the usage block (`app.go:75`); `help`/`-h`/`--help` is dispatched but still not listed. `agents probe` alias is also not listed. Re-flagged as NIT.

## Findings

### CRITICAL

None.

### MAJOR

#### 1. `parley run` launches every installed agent with no user selection

`app.go:163-168`:

```go
discovered := agents.Discover(ctx, agents.DefaultSpecs())
participants := installedAgentIDs(discovered)
if len(participants) == 0 {
    fmt.Fprintln(stderr, "no installed headless agents found; run `parley agents discover` ...")
    return 1
}
idea, err := protocol.CreateIdea(*root, task, participants)
```

There is no `--participants codex,claude` flag, no interactive picker, no confirmation gate, and the idea's roster is set to "everything that happens to be on PATH." `FINAL.md:60-67` describes `parley run` as: *"Creates a new idea and kickoff prompt. **Selects participants.** Starts round 1 in parallel."* The current implementation skips the selection step.

Why this matters:

- On a typical contributor machine all four adapters (codex, claude, gemini, hermes) will run for every `parley run`, each spending real model tokens with the default 30-minute timeout. That is a real cost regression.
- It contradicts the HITL-default posture in `FINAL.md:177-189`. Without `--auto`, an unprompted invocation should not silently start four paid agents.
- `IMPLEMENTATION.md` cycle-2 verification notes acknowledge it implicitly ("the new command now launches all detected models, ... I did not run `parley run` against the real installed agents in this verification pass") but the gap is **not** listed under `Deviations from FINAL.md`.

Suggested fix: add a `--participants codex,claude` flag (comma-separated, defaults to a fail-closed empty list or to a per-workspace default) and either prompt or require explicit selection when more than one agent is installed. At minimum, declare the auto-fanout behavior under deviations and gate it behind `--auto` until selection is implemented.

#### 2. Runner swallows every `Store.Append` error

Every event write in `runner.go` uses `_ = opts.Store.Append(...)` (lines 77, 125, 179, 208). The whole rationale for the JSONL event store (`FINAL.md:101-103` — "append-only runtime events" feeding "fast resume") depends on those writes succeeding. A failing append (disk full, permissions revoked, parent dir gone) is invisible to the user and to `parley resume`'s future implementation.

Suggested fix: at minimum log the error to stderr; better, surface it through `Result` so `printRunResults` can warn the user that durable state may be incomplete. The store's own `Append` already returns the error — only the call sites discard it.

#### 3. Early-return paths in `runAgent` emit no `agent.*` event

`runner.go:107-146`. If `os.MkdirAll(agentDir, ...)` (line 133), `os.MkdirAll(filepath.Dir(outputPath), ...)` (line 137), `BuildRoundOnePrompt` (line 142), `commandFor` (line 151), or either `os.Create` for stdout/stderr (lines 161, 167) fails, the function sets `result.ExitError` and returns without ever appending `agent.started` or `agent.failed`. The roll-up at the end of `RunRoundOne` correctly classifies the round as `round.incomplete`, but `events.jsonl` never names the failing agent or its error.

Why this matters: this is exactly the durability gap the JSONL store is supposed to close. A future `parley resume` (or post-mortem of an automated run) will see `run.created` → `round.incomplete` with no agent records, even though one specific agent failed for a known reason. It is also surprising that `agent.skipped` (line 125) is emitted but `agent.failed` for setup failures is not.

Suggested fix: either centralize the early-return through a `defer` that always emits `agent.failed` with the captured `result.ExitError`, or bracket the function body with explicit `agent.started` (after stdin/stdout setup) and `agent.failed` writes so every code path emits exactly one terminal event per agent.

### MINOR

#### 4. `agents.Spec.IsolatedHomeID` is dead state

`internal/agents/discover.go:20` adds `IsolatedHomeID string` to the spec, but no caller reads or sets it (verified with a workspace grep — only the declaration matches). Same shape as the round-01 `IdeaStatus.Path` finding. Either wire it (e.g., as the temp-dir prefix in `isolatedGeminiHome`) or remove it — per `CLAUDE.md` "Surgical Changes," fields added by your changes that nothing reads should be deleted.

#### 5. `isolatedGeminiHome` silently proceeds when oauth files are missing

`runner.go:301-310`:

```go
for _, name := range []string{"oauth_creds.json", "google_accounts.json"} {
    source := filepath.Join(os.Getenv("HOME"), ".gemini", name)
    data, err := os.ReadFile(source)
    if err != nil {
        continue
    }
    ...
}
```

If `oauth_creds.json` does not exist (the user has not authenticated Gemini, or `$HOME` is unset in some sandboxed contexts), the runner still creates the temp home, writes only `settings.json`, and spawns Gemini against an effectively empty profile. Gemini will fail with an opaque error after launch, and the failure shows up downstream as "artifact was not created" rather than "Gemini is not authenticated." Also, `os.Getenv("HOME")` is empty on Windows; `os.UserHomeDir()` is the portable choice.

Suggested fix: fail fast with a clear "Gemini OAuth credentials not found at $HOME/.gemini/oauth_creds.json — run `gemini` once to authenticate" when neither credential file is present, and use `os.UserHomeDir()`.

#### 6. Hermes adapter still in `DefaultSpecs`

`internal/agents/discover.go:73-81` still ships Hermes despite the round-01 finding. `FINAL.md:97-99` lists `codex, claude, gemini` and nothing else; the `command name is not standardized yet` note is itself an admission that the entry is speculative. With the cycle-2 change to "launch every installed agent," any user with a `hermes` binary on PATH will now have it executed for every `parley run`. Either add Hermes to `FINAL.md`'s adapter list (with consensus) or drop it.

#### 7. TUI opens after the synchronous round completes

`app.go:208`. When `--no-tui` is not set, `runTask` runs all agents to completion (worst case `opts.Timeout = 30m` per agent in parallel), then calls `runTUIView`. The TUI is therefore a static post-mortem dashboard, not a live orchestration surface. `IMPLEMENTATION.md` does declare "Live TUI streaming is not implemented yet," so this is not undeclared, but the current ordering means a user who omits `--no-tui` stares at a frozen terminal for the duration of round-01. Consider either inverting (open the TUI first and run agents on a goroutine, even if the TUI shows "round in progress") or printing live `[codex] writing artifact...` lines on stdout while waiting.

#### 8. Concurrent writes to `events.jsonl` rely on POSIX `O_APPEND` atomicity

`runner.go:58-66` launches all selected agents in parallel; each goroutine repeatedly calls `opts.Store.Append`, which performs `OpenFile(O_APPEND) + Write + Close` per event (`store/events.go:38-50`). On POSIX, `write()` ≤ PIPE_BUF (typically 4 KB) on an `O_APPEND` fd is atomic — JSONL events here are small enough that this works. On Windows the guarantee is filesystem-dependent, and the cycle-2 verification claims include a Windows cross-build. Add a `sync.Mutex` to `Store` (or keep a long-lived `*os.File` guarded by one) so the contract doesn't depend on platform-specific append semantics.

#### 9. `BuildRoundOnePrompt` hardcodes `cli-default` for model/thinking/profile

`runner.go:240-243` writes:

```
Effective launch config:
- model: cli-default
- thinking/reasoning/effort/profile: cli-default
```

`Spec` carries no model/effort fields yet, so today this is honest — but it baked into the prompt template, which means future per-spec overrides will need to thread through `BuildRoundOnePrompt` as well. Either parameterize now (`opts.LaunchConfig`) or drop the section until there is something real to inject.

### NIT

#### 10. `printUsage` still misses `help` and the `agents probe` alias

`app.go:64-77` lists `init`, `agents discover`, `status`, `run`, `resume`, `answer`, `tui`, `version` — but the dispatcher also accepts `help`/`-h`/`--help` and `agents probe`. Add them for symmetry.

#### 11. `Result.Duration` and `CompletedAt` not populated on early-return paths

`runner.go:113-146`. Combined with finding #3, the structured result for setup-time failures has zero `Duration` and zero `CompletedAt`. `printRunResults` (`app.go:264-288`) currently only formats `Duration` on the success branch, so this is not user-visible today, but the `Result` struct's contract is inconsistent. Set `result.CompletedAt = time.Now().UTC()` and compute `result.Duration` in the early-return helper alongside the `agent.failed` event suggested in #3.

#### 12. `Store.Append` opens/closes per event

Carry-over from round-01 NIT. With cycle-2's parallel agent goroutines and per-agent stdout/stderr/event traffic, this is now a measurable cost rather than a theoretical one. A single long-lived `*os.File` on `Store` plus a `sync.Mutex` would also fix finding #8.

#### 13. `Spec.IsolateHome` is per-spec but `isolatedGeminiHome` is hardcoded for Gemini

`runner.go:280-287` checks `agent.IsolateHome` but unconditionally calls `isolatedGeminiHome()`. If a future adapter sets `IsolateHome: true`, it will get a Gemini-shaped temp home. Either rename the function/field (`isolateHome` → `geminiIsolatedHome` is fine) or dispatch by spec ID.

## Verification notes

- `IMPLEMENTATION.md` reports `go test ./...` passes and three cross-compile targets build (`darwin/arm64`, `linux/amd64`, `windows/amd64`). I did not re-run `go test` in this review pass; I read the test source. The test contract looks sound: `TestRunRoundOneCreatesArtifactWithHeadlessAgent` re-invokes the test binary via `os.Args[0]` with `-test.run=TestFakeAgentHelper`, which writes a synthetic round-01 artifact and exits 0. This validates: participant filtering through `selectedAgents`, prompt → stdin plumbing, output-path discovery, and `round.completed` event emission. It does **not** cover: store-append failure, mkdir failure, prompt-build failure, parent-context cancellation, or skipped-when-artifact-already-exists (although the helper deliberately doesn't pre-create the file). Worth adding at least the "skipped when artifact exists" and "agent.failed when binary exits non-zero" cases in cycle 3.
- Hermes is still discovered/launched in this environment per `IMPLEMENTATION.md`. The implementer's note "I did not run `parley run` against the real installed agents in this verification pass" is consistent with the participant-selection finding above; it is the right cautious choice given the current behavior.
- I did not exercise the runner against real Codex/Claude/Gemini/Hermes from this review; verification of the `headless launch specifications` (`agents/discover.go:41-82`) — particularly the `--ask-for-approval never` and `--permission-mode acceptEdits` flags — is best done in a controlled smoke run with cost limits.
- `head-commit: uncommitted` continues to make findings hard to anchor to a SHA. Once the repo gets its first commit, future review rounds should reference it.

## Recommendation

**Proceed to cycle 3 conditional on closing findings #1, #2, and #3.** None are CRITICAL; the runner contract is sound and round-01 issues are addressed. But cycle 3 (live TUI streaming, token accounting, HITL capture) would compound on top of these durability and selection gaps if left unfixed:

- #1 (participant selection) blocks safe real-agent verification — the same blocker the implementer cites.
- #2 and #3 (silent event-store failures, missing `agent.*` events on early returns) directly undercut the design goal of "rebuild state from durable events" before that goal is exercised.

Findings #4–#13 are non-blocking and can ride along with cycle 3.
