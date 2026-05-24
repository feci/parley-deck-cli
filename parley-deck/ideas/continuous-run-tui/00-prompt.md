---
idea: continuous-run-tui
author: user
created: 2026-05-23
participants: [codex, claude, gemini, hermes]
roles:
  codex: current CLI/TUI architecture and implementation slicing
  claude: TUI UX and continuation semantics review
  gemini: durable state and cross-agent interoperability review
  hermes: long-running run recovery and stale process model review
status: final
---

## Problem / idea

Design the next `parley-deck-cli` slice so Parley can be used continuously from the TUI.

Today the CLI can create new ideas and start new round-01 runs, but the user cannot comfortably continue work on a specific existing run. The TUI should let the user switch between several ideas or runs, intervene when needed, and then continue the interrupted run without restarting the whole idea.

The user specifically wants:

- continue a concrete run or idea from the TUI;
- see several active ideas/runs at once;
- switch to another idea while one run is waiting or running;
- answer or intervene when a run needs human input;
- resume, retry, or advance only the missing part of a workflow;
- avoid waiting idly for one idea before working on another.

## Existing context

Relevant prior Parley ideas:

- `live-run-tui`: implemented live viewing for one active run; resume/reattach was explicitly out of scope.
- `runtime-status-resume`: implemented durable read-only status/resume over `events.jsonl`; resume means state recovery and view restoration, not OS process reattachment.
- `tui-workspace-sessions`: implemented the workspace session console, session list, TUI start-new-idea mode, and in-process parallel TUI-started runs.
- `session-resume-cache-plan`: implemented slice 1 with `run.json`, `~/.parley-deck/sessions.json`, `parley sessions list`, and `parley sessions inspect`; deferred native resume handles, prompt packs, terminal status updates, stale recovery, and TUI retry/resume actions.

## Constraints

- Preserve canonical Parley artifacts under `parley-deck/ideas/<slug>/`.
- Do not edit another agent's artifact as a normal recovery path.
- Keep local cache/state under `~/.parley-deck` or `PARLEY_HOME`.
- Do not store credentials or full private transcripts in Parley state.
- Do not promise true OS-process reattachment unless a later supervisor/daemon proves it.
- Work with existing Go packages where practical: `runcontrol`, `runmanifest`, `sessionstore`, `runstate`, `runner`, `tui`, and `app`.
- Keep the MVP small enough to implement and test in slices.

## Non-goals

- No hosted service.
- No mandatory daemon for the first implementation slice.
- No scraping private agent transcript formats as the only source of truth.
- No automatic phase advancement that bypasses Parley consensus rules.
- No cross-run implementation file locking beyond the concrete files each action mutates.
