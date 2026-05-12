---
idea: runtime-status-resume
author: codex
created: 2026-05-12
participants: [codex, claude, gemini, hermes]
status: round-01
---

## Problem / idea

Design the next `parley-deck-cli` slice: runtime status and resume.

After `agent-runtime-config`, the CLI can surface effective agent setup, and `parley run` already creates durable runtime events under `parley-deck/runs/<run-id>/events.jsonl`. The next product gap is recovery and observability. Improve `parley status` and implement useful `parley resume <idea-or-run>` so a user can recover after TUI detach, process exit, interruption, or context compaction.

The CLI should rebuild state from durable protocol files and runtime event files. It should show the current idea, phase, round, run IDs, participants, agent states, log paths, pending HITL questions, artifact completion, and the next recoverable command or protocol action. If live process reattachment is feasible in the current architecture, design it narrowly; otherwise be explicit that resume restores the view and recoverable state rather than resurrecting dead subprocesses.

## Constraints

- Active Parley Deck transport is `github-pr`; canonical files remain under `parley-deck/`.
- Build on the current Go + Bubble Tea architecture and the merged `agent-runtime-config` runtime matrix.
- Keep the implementation scope small enough for one design PR and one implementation PR.
- Preserve `--no-tui` and plain terminal output; status and resume must work in CI/plain terminals.
- Use durable files already present or easy to derive: protocol frontmatter, round/review artifacts, `runs/<run-id>/events.jsonl`, logs, and HITL question files.
- Make reattachment claims conservative. Do not promise cross-process control unless the implementation can support it with stored supervisor/process state.
- The implementation must be testable with deterministic fixture events and protocol files.
- Every file under `parley-deck/` and every PR comment/review summary must be in English.

## Non-goals

- Do not implement full auto consensus or generic phase advancement automation in this slice.
- Do not implement GitHub/GitLab transport automation in this slice.
- Do not implement token/cost accounting beyond displaying fields already present in events.
- Do not add a hidden long-running daemon unless the design concludes it is the smallest necessary mechanism.
- Do not implement release packaging, npm, Homebrew, or CI release workflows in this slice.
