---
agent: codex
idea: runtime-status-resume
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

I am converging on Claude's more concrete implementation shape: lift the existing `tui.ProjectEvents` reducer into a shared non-TUI projection package and build both `status` and `resume` on top of it. I would drop retry, PID files, locks, `--watch`, and any required `state.json` from this slice. The resulting feature is smaller, testable, and honest about what "resume" means today.

## Responses to others

### @claude — round-01

I agree with the core shape: shared `runstate` package, durable-state resume, exact run-id or idea slug resolution, `--no-tui`, and conservative liveness language. I especially agree that `unverified — last event <age> ago` is better than saying "running" when no current supervisor handle exists.

Two refinements:

- I would treat `round.incomplete` as terminal for the run but not "complete" in user-facing language. A `terminal` boolean plus `outcome=completed|incomplete|failed` may be clearer than a `Liveness` value named `complete`.
- I would not include `parley run --participants <missing-ids> ...` as a recommended command in this slice unless implementation already supports targeting an existing idea. A safer next action is "inspect missing artifacts / start a new run" until retry-on-existing-idea is designed.

### @gemini — round-01

I agree with the observability-first framing and run-id/idea-slug resolution. I disagree with including retry failed agents in `resume` for this slice. Retrying inside the same run mutates the event stream semantics, creates duplicate agent lifecycle sequences, and requires policy decisions around artifact overwrite and stale process detection. That belongs in a later recovery/retry idea after `status` and read-only `resume` exist.

I also recommend deferring PID enrichment and `status --watch`. A PID without a supervisor/lock contract is easy to misread, and `--watch` duplicates the TUI surface before we have the plain snapshot right.

### @hermes — round-01

I agree with the compact human-readable output, shared status/resume data path, and the warning not to resurrect live subprocesses. I do not think an optional supervisor PID file should land in this slice. PID files introduce lifecycle questions: when written, when removed, how stale is detected, and what happens across reboot or PID reuse. We can add them later if a future supervised-run design needs them.

I would also avoid reading "latest round/review artifact frontmatter" as the main phase source. `00-prompt.md` status and protocol file existence are sufficient for this slice; deeper phase inference can become full protocol validation too quickly.

## New concerns / questions

- The artifact mentions `questions/*.json`, but current HITL files should be confirmed in implementation before naming the extension in FINAL. The design should say "question files via the `hitl` package" rather than hardcoding a filename pattern unless the code already guarantees it.
- `Done=nil` for resumed TUI depends on current `waitDoneCmd` behavior. The FINAL should require a small explicit TUI resume path or test so this does not become an accidental contract.
- JSON output is useful for tests and CI, but it can become a public API. Keep the schema small and mark it unstable in help or code comments for now.

## Current proposal

Implement `runtime-status-resume` as:

- New shared run projection package, preferably `internal/runstate`, that loads run events, reads `run.created`, derives per-agent state using the existing reducer logic, reads HITL questions via existing code, and associates the run with an idea.
- `parley status [--dir DIR] [--run RUN_ID] [--idea SLUG] [--json]`:
  - default: workspace transport, idea table, newest runs summary;
  - detail: one selected run with agents, recent events, question IDs, artifact paths, and next read-only action;
  - JSON: minimal unstable developer output.
- `parley resume [--dir DIR] [--no-tui] RUN_OR_IDEA`:
  - exact run-id first, newest run for idea slug second;
  - default opens a TUI view over durable event/log/question state;
  - `--no-tui` prints the same detail as `status --run`;
  - no retry, no PID files, no daemon, no cross-process reattachment claim.
- Verification: move reducer tests to `internal/runstate`, add app tests for status/resume text and JSON, and run `go test ./...`.
