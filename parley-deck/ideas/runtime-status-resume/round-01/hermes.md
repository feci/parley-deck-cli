---
agent: hermes
idea: runtime-status-resume
round: 1
date: 2026-05-12
---

## Summary
After agent-runtime-config, the CLI already emits durable events under runs/<run-id>/events.jsonl and stores protocol artifacts. The immediate next slice should extend `parley status` (and the new `parley resume <idea-or-run>`) to load and present a compact, human-readable view of the current idea phase, active run, participant states, and next actionable step entirely from those durable files. This enables recovery after TUI detach or process exit without claiming to resurrect live subprocesses.

## Proposed approach
Implement a small `status` package (or extend protocol/workspace.go) that walks the parley-deck/ tree, reads frontmatter from the latest round/review artifact for the target idea (or infers from run-id via events), merges in the most recent events.jsonl lines for that run, and renders a plain-text table or key-value block to stdout. 

`parley status` (no arg) lists open ideas + their latest phase/round + last run id. 
`parley resume <slug-or-runid>` loads the same data, optionally starts a fresh TUI view in "resume" mode that replays recent events but does not attempt to re-attach the original runner PID. If the original process is gone, it simply shows the last known state and offers the next CLI command (e.g. `parley run --resume` or `parley answer`).

Store only a tiny supervisor pid file under runs/<run-id>/pid (optional, best-effort) so status can report "live" vs "detached". Never promise reattachment across machines or after SIGKILL; the conservative claim is "view + state restoration for the next human or agent step".

Keep `--no-tui` path identical to the TUI renderer (plain text only). Add unit tests that feed fixture events.jsonl + protocol .md files and assert deterministic output strings.

## Concerns / open questions
- How to reliably map a run-id back to its originating idea slug when the run was started from `parley run` inside an idea directory.
- Whether to surface pending HITL questions from the hitl/ subdir or only from events; the prompt mentions both.
- Exact output format that stays useful in CI logs yet rich enough for a human resuming work.

## Risks
- Over-promising reattachment could lead to confusing UX if the pid file is stale; keep the "live" indicator strictly advisory.
- Adding too many fields to the status view before we have agreement on what the matrix of agent-runtime-config actually emits.