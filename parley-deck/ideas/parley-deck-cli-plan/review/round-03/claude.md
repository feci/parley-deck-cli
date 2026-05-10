---
agent: claude
idea: parley-deck-cli-plan
review-round: 3
date: 2026-05-10
reviewing: implementation-fix-up-cycle-3
---

## Summary

Fix-up cycle 3 closes all three round-02 MAJORs. Participant selection and the HITL confirmation gate are wired (`--participants`, `--yes`, plus a TTY-driven prompt that fail-closes when stdin is non-interactive), every `Store.Append` call now propagates errors into `Result.ExitError`, and a new `failEarly` helper guarantees an `agent.failed` event for setup-time failures while populating `CompletedAt`/`Duration`. The append path is also serialized with a process-wide `sync.Mutex`, the dead `IsolatedHomeID` field is gone, `isolatedGeminiHome` now uses `os.UserHomeDir()` and fails fast when no creds and no API key env are present, and runner tests grew to cover the skip and child-failure cases. Cycle 3 is in good shape.

A handful of round-02 MINOR/NIT items remain (Hermes still in `DefaultSpecs`, post-round TUI ordering, hard-coded `cli-default` in the prompt template, `isolatedGeminiHome` hardcoded for Gemini behind a generic `IsolateHome`). All are explicitly declared in `IMPLEMENTATION.md` deviations or are non-blocking. Two small new observations are noted below as MINOR/NIT.

## Round-02 finding status

- **[MAJOR] #1 participant selection / HITL gate** — **Closed.** `app.go:146-180`: new `--participants` and `--yes` flags; `selectedParticipantIDs` (`app.go:280-306`) validates IDs against the installed set and errors before any idea is created; `confirmLaunch` (`app.go:308-322`) prompts on a TTY and fail-closes on non-TTY stdin so an unattended `parley run` cannot silently fan out paid agents. `--auto`/`--yes` are the documented escape.
- **[MAJOR] #2 store append errors swallowed** — **Closed.** Every call site now checks the error: `runner.go:77-92` (round roll-up appends a synthetic `runner` Result on failure), `runner.go:132-138` (skip path sets `ExitError`), `runner.go:182-193` (`agent.started` routes through `failEarly`), `runner.go:213-225` (`agent.finished`/`agent.failed` composes the error via `combineError`), and `app.go:189-201` (`run.created` returns 1). `printRunResults` surfaces the error to the user.
- **[MAJOR] #3 early-return paths emit no `agent.*` event** — **Closed.** `failEarly` (`runner.go:230-248`) is called from every pre-launch error path (`MkdirAll` × 2, `BuildRoundOnePrompt`, `commandFor`, both `os.Create`s, and the `agent.started` append failure), and emits a single terminal `agent.failed` event with populated `CompletedAt`/`Duration`. Verified by `TestRunRoundOneRecordsAgentFailure` (`runner_test.go:143-187`).
- **[MINOR] #4 dead `IsolatedHomeID`** — **Closed.** Field removed from `Spec`.
- **[MINOR] #5 `isolatedGeminiHome` silent on missing creds + non-portable `$HOME`** — **Closed.** `runner.go:325-355` uses `os.UserHomeDir()` and returns a clear error when neither OAuth file is present and no `GEMINI_API_KEY`/`GOOGLE_API_KEY` is set.
- **[MINOR] #6 Hermes still in `DefaultSpecs`** — **Acknowledged, not removed.** Now declared under cycle-3 deviations: "supported because the user explicitly requested Hermes after the original finalized plan." Acceptable given the gate above prevents silent fan-out.
- **[MINOR] #7 TUI opens after the synchronous round** — **Unchanged**, declared in deviations as live streaming deferred. Cycle-3 added a small win by reusing the discovery list (`app.go:220`, `runTUIViewWithDiscovery`) so the post-run TUI doesn't re-probe.
- **[MINOR] #8 concurrent writes to events.jsonl** — **Closed.** `store/events.go:24` adds `var appendMu sync.Mutex` and `Append` brackets the file open/write/close in the lock, so atomicity no longer depends on platform `O_APPEND` semantics.
- **[MINOR] #9 `BuildRoundOnePrompt` hardcodes `cli-default`** — **Unchanged.** Still printed verbatim at `runner.go:276-280`. Non-blocking; revisit when `Spec` carries model/effort fields.
- **[NIT] #10 `printUsage` missing `help` / `agents probe`** — **Closed.** `app.go:65-79` lists both (`agents discover|probe`, `help`).
- **[NIT] #11 `Duration`/`CompletedAt` zero on early-return** — **Closed** by `failEarly`.
- **[NIT] #12 `Store.Append` opens/closes per event** — **Partially.** Mutex was added (the more important guarantee); per-event open/close remains. Non-blocking.
- **[NIT] #13 `IsolateHome` per-spec but `isolatedGeminiHome` hardcoded** — **Unchanged.** Only Gemini sets `IsolateHome`, so latent. Non-blocking.

## Findings

### CRITICAL

None.

### MAJOR

None.

### MINOR

#### 1. Skip path can mask a real durability failure

`runner.go:127-140`: when the artifact already exists, the function emits `agent.skipped` and, on event-append failure, sets `result.ExitError` to "event append failed: …" and returns. It does **not** also emit a corresponding `agent.failed` event the way `failEarly` does for other early returns. `printRunResults` (`app.go:332-347`) renders this row as a *failure* (because `ExitError != ""` flips the switch table away from the `Skipped && ExitError == ""` arm), but `events.jsonl` shows neither `agent.skipped` (the append failed) nor `agent.failed` for that agent — the same gap #3 closed for the other paths. Either route the skip-append-failure through `failEarly` or accept that the user-visible result and the event log diverge.

#### 2. `confirmLaunch` fail-closes on non-TTY stdin without `--yes`/`--auto` documented as the only escape

`app.go:308-322`: when stdin is a pipe/heredoc/CI runner, `confirmLaunch` returns false and the command exits 0 with "No run started. Use `--yes` or `--auto` …". This is the right HITL default, but two small things:

- Exit code 0 for "no run started" is debatable — a script that pipes a task in and expects the run to proceed will treat a no-op as success. Either use a non-zero code (e.g. 2) for the unattended-without-`--yes` case, or note this in usage.
- The hint text says `--yes` *or* `--auto`, but `--auto` is documented as "automatic low-risk progression policy," not a confirmation skip. They happen to be equivalent here because `confirmLaunch` is gated by `!*auto && !*yes`, but conflating the two flags is going to be confusing once `--auto` carries real progression semantics. Consider keeping `--yes` as the only confirmation-skip and not coupling it to `--auto`.

### NIT

#### 3. No test exercises participant selection / confirmation gate

`runner_test.go` covers the runner contract (success, skip, child failure) but `app.go` selection and `confirmLaunch` paths are untested. A small table-test on `selectedParticipantIDs` (empty input, comma-separated, unknown ID, duplicate) would be cheap insurance for what is now the cost-control surface.

#### 4. `BuildRoundOnePrompt` still emits `cli-default` literals

Carry-over from round-02 #9. Honest today; re-flag when `Spec` grows model/effort fields.

#### 5. `isolatedGeminiHome` still dispatched generically

Carry-over from round-02 #13. Latent; only Gemini sets `IsolateHome`.

## Verification notes

- I read the source rather than running `go test ./...`. The new tests look right: `TestRunRoundOneSkipsExistingArtifact` (asserts `agent.skipped` then `round.completed`) and `TestRunRoundOneRecordsAgentFailure` (asserts `agent.failed` then `round.incomplete` with the runner test binary returning exit 7) directly cover the round-02 behavior gaps.
- The cycle-3 smoke run `parley run --no-tui 'non interactive hitl confirmation smoke'` is a valid HITL-default check: with non-TTY stdin (sandbox) the command should print the HITL message and exit without launching agents — no paid model traffic, exactly the safety property #1 was meant to deliver.
- I did not exercise real agents either; that gating is now safe to do behind `--participants codex --yes` once a contributor wants to validate the headless launch flags end-to-end.
- `head-commit: uncommitted` carries through; once the repo gets its first commit, future rounds can anchor findings to a SHA.

## Recommendation

**Approve cycle 3.** No CRITICAL or MAJOR findings remain. All three round-02 MAJORs are closed, the round-02 MINORs that mattered for safety/durability (#4, #5, #8, #11) are closed, and the still-open MINOR/NIT items are explicitly declared as deferred work or are non-blocking. MINOR #1 (skip-append-failure path) and the small UX points in MINOR #2 are worth a follow-up but should not block proceeding to live TUI streaming, token accounting, or HITL capture in the next slice.
