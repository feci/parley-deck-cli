---
idea: session-resume-cache-plan
author: user
created: 2026-05-18
participants: [codex, claude, gemini, hermes]
roles:
  codex: local CLI state and Parley runner design
  claude: session resume UX and durable metadata design
  gemini: cross-agent cache discovery and interoperability
  hermes: long-running session registry and recovery model
status: round-01
---

## Problem / idea

Design how `parley-deck-cli` should support durable continuation of older sessions after the CLI/TUI is closed or the machine restarts.

The user wants to be able to:

- stop `parley-deck-cli`;
- reopen it later;
- see older Parley sessions and their state;
- continue where the run stopped;
- preserve enough local context to resume agents or at least re-invoke missing work without restarting the whole idea;
- store this under local cache/state, likely under `~/.parley-deck`, while preserving the canonical `parley-deck/` artifacts.

## Local observations gathered before kickoff

Observed local agent persistence patterns, without reading secrets or transcript contents:

- Codex stores session index and transcripts under `~/.codex/session_index.jsonl` and `~/.codex/sessions/YYYY/MM/DD/*.jsonl`; top-level session records use `timestamp,type,payload`. It also has `state_5.sqlite` with tables such as `threads`, `agent_jobs`, `agent_job_items`, `thread_goals`, and `thread_spawn_edges`. The CLI exposes `codex resume` and `codex fork`.
- Claude Code stores project-scoped JSONL sessions under `~/.claude/projects/<escaped-project-path>/*.jsonl`; records include keys such as `type,permissionMode,sessionId`. It exposes `claude --continue`, `claude --resume [session-id/search]`, `--fork-session`, `--session-id`, and `--no-session-persistence`.
- Gemini CLI stores per-project chats under `~/.gemini/tmp/<project>/chats/session-*.jsonl` and project/trust metadata under `~/.gemini/projects.json`, `~/.gemini/state.json`, and `~/.gemini/trustedFolders.json`. It exposes `gemini --resume latest|<index>`, `--session-id`, `--list-sessions`, and `--delete-session`.
- Hermes stores sessions under `~/.hermes/sessions/session_*.json` with keys such as `session_id,model,base_url,platform,session_start,last_updated,system_prompt,tools,message_count,messages`. It exposes `hermes --resume <session-id>`, `hermes --continue [name]`, and `hermes sessions list/browse/rename/export/prune/delete`.
- Parley currently stores only an advisory session list under `~/.parley-deck/sessions.json` and canonical run events under `<workspace>/parley-deck/runs/<run-id>/events.jsonl`. This is enough for listing, but not enough for a full resumable orchestration contract.

## Constraints

- Keep canonical protocol artifacts under each workspace's `parley-deck/`.
- Local user-machine state may live under `~/.parley-deck`.
- Do not store credentials or agent auth material in Parley state.
- Do not assume all external CLIs can truly resume headless sessions; support capability detection.
- Preserve per-agent ownership: Parley may resume/re-invoke an agent, but must not edit another agent's canonical artifact as a normal path.
- Work with local-dir and github-pr transports.
- The plan should cover TUI and CLI behavior.

## Non-goals

- Do not implement in this idea.
- Do not copy external agent transcript contents into Parley unless explicitly needed and safe.
- Do not depend on undocumented private cache formats as the only source of truth.
