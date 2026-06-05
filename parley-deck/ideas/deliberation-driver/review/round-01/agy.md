---
agent: agy
idea: deliberation-driver
review-round: 1
date: 2026-06-05
reviewed-commit: f8c880d
---

## Summary

This review evaluates the first implemented slice of the `deliberation-driver` (`internal/driver`) for protocol correctness, alignment with `FINAL.md` / `COOPERATION.md`, and robustness of transport gating. The implementation of Slice 1 successfully addresses the primary stall in `parley run` by orchestrating the promotion of completed independent rounds (`round-01`) to cross-review rounds (`round-02`) under the `--auto` flag and `local-dir` transport. The codebase is clean, tests are green, and the slice boundaries are safely respected. However, several findings regarding parsing robustness, locking concurrency, and error escalation need to be addressed.

## Findings

### [MAJOR] Fragile Transport Parsing in `EffectiveTransport`
* **What is wrong:** In `internal/driver/transport.go`, `EffectiveTransport` reads the global transport from `COOPERATION.md` by splitting lines and searching for `"Transport:"`, then passing the line to `backtickValue` which strictly parses content inside backticks. If a user edits `COOPERATION.md` to say `**Transport:** local-dir` without backticks (a common Markdown variation), `backtickValue` returns `""` and the parser fails.
* **Why it matters:** This makes the transport gate fragile under manual editing. When it fails to parse, it silently disables auto-advance and halts with `ActionSurfaceOnly`, which may confuse users who did configure `local-dir`.
* **Concrete suggested fix:** Align the parsing logic with the regex used in `internal/protocol/workspace.go:readTransport`, which handles optional backticks:
  ```go
  re := regexp.MustCompile(`^\*\*Transport:\*\*\s*` + "`?" + `([^` + "`" + `\s]+)`)
  ```
  Alternatively, write a robust parser in `transport.go` that strips both asterisks and backticks from the transport value.

### [MINOR] Missing Escalation on Agent/Runner Execution Failures
* **What is wrong:** If `d.runner.RunRound` fails during `Advance` (due to an agent crash, API timeout, or permissions error), the driver prints the error and returns `ActionEscalated`, but it does not write a blocking escalation file to the `inbox/` directory.
* **Why it matters:** `COOPERATION.md` §4 ("Recovery And Partial Completion") states that for non-zero CLI exits, rate limits, or auth failures, the facilitator/driver should capture the failure in an inbox note to alert the user.
* **Concrete suggested fix:** Update `Advance` or `Run` to write an inbox escalation note (similar to `escalateDeadline`) when a runner execution error is encountered.

### [MINOR] TOCTOU Window in Advisory `driver.lock` Acquisition
* **What is wrong:** `acquireLock` in `internal/driver/loop.go` uses a simple read-then-write check on the PID file without OS-level file locking or atomic creation flags (`O_EXCL`).
* **Why it matters:** Two concurrent driver invocations could read the stale or empty lock file simultaneously, both proceed to write their PIDs, and both run concurrently, violating the single-writer assumption.
* **Concrete suggested fix:** Use atomic file creation (`os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)`) or standard lockfile libraries, and handle PID cleanup appropriately.

### [NIT] Non-re-evaluated Transport Value in Loop
* **What is wrong:** `AutoLocalDir` is passed as a static configuration field to `driver.New` from `app.go`. The running poll loop in `Run` does not re-read the transport from disk on each tick.
* **Why it matters:** While a minor edge case, `FINAL.md` D8 states the transport choice is "re-evaluated every tick" so that a mid-run transport change can be detected.
* **Concrete suggested fix:** Re-evaluate the effective transport dynamically inside `Advance` or re-read the configuration file on each tick rather than relying on a static config field.

## Open questions

### 1. Deviation on `### @` Heading check
The implementation validates round completeness using the `responding-to` frontmatter field and `runner.ValidateRoundArtifact` rather than checking for per-agent `### @<other>` headings. This deviation is **protocol-acceptable** and highly pragmatic because the runner's built-in `BuildRoundPrompt` does not instruct agents to write per-agent subheadings; checking for them would reject valid runner output. 
* *Question:* Should we update the global `COOPERATION.md` to make the heading check optional or update the runner prompt in a future slice to align with the original protocol specification?

### 2. Driver recovery path on incomplete rounds
If the driver encounters an incomplete round, it halts and eventually escalates via `escalateDeadline`. 
* *Question:* Once the user manually fixes the missing artifact or adjusts the roster, is running `parley run --auto` the only supported recovery path, or should there be a dedicated command to resume/re-trigger the driver tick?
