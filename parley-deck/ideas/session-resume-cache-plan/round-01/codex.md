---
agent: codex
idea: session-resume-cache-plan
round: 1
date: 2026-05-18
---

## Summary

Parley should treat its own state as the durable orchestration source of truth, not as a thin mirror of each vendor CLI's private cache. Existing agent CLIs prove the pattern: keep a local session index, append-only event/transcript records, stable session IDs, project/workspace scoping, and explicit resume commands. For Parley, this means adding a first-class run ledger under `~/.parley-deck` that links workspace runs, protocol artifacts, agent subprocess attempts, HITL questions, and optional external CLI resume handles.

## Proposed approach

Introduce a local Parley state store with three layers:

1. Global session index:
   - Path: `~/.parley-deck/sessions.json` or `~/.parley-deck/state.sqlite`.
   - Purpose: fast TUI startup, workspace/run discovery, recent sessions, attention state, last event time.
   - Store only Parley metadata: workspace root, idea slug, run ID, transport, participants, status, timestamps, and summary.

2. Per-run durable orchestration directory:
   - Path: `~/.parley-deck/runs/<run-id>/` for cross-workspace cache, plus existing canonical `<workspace>/parley-deck/runs/<run-id>/`.
   - Files:
     - `manifest.json`: schema version, workspace root, idea slug, run ID, created/updated timestamps, command/config snapshot.
     - `events.jsonl`: append-only Parley orchestration events, optionally mirrored from workspace run events.
     - `agents/<agent-id>/attempts/<attempt-id>.json`: command path, args template, cwd, env whitelist, model/profile, PID if active, exit state, stdout/stderr log paths, artifact path, validation result.
     - `resume.json`: per-agent resume capability and handle, e.g. `claude.session_id`, `gemini.session_id`, `hermes.session_id`, `codex.thread_id` when discoverable.
     - `questions/`: mirrored HITL question index, with canonical question files still under workspace run state.

3. Canonical protocol artifacts:
   - Continue using `<workspace>/parley-deck/ideas/<slug>/...` as the authoritative outcome.
   - Resume logic reads missing/invalid artifacts from the workspace, not from cached transcripts.

Add an explicit resume capability model:

- `none`: cannot resume; re-invoke a fresh agent attempt with prior prompt plus current artifacts.
- `interactive-only`: can resume in TUI/manual mode, but headless continuation is not reliable.
- `headless-session`: can pass a stable session ID or continue flag in headless mode.
- `fork`: can resume or fork previous context.

Populate this from CLI discovery and local config:

- Claude: `--resume`, `--continue`, `--session-id`, `--fork-session`.
- Gemini: `--resume`, `--session-id`, `--list-sessions`.
- Hermes: `--resume`, `--continue`, `sessions list/export`.
- Codex: `resume`, `fork`, plus local `session_index.jsonl` and session JSONL; headless `exec` may not equal interactive resume.

CLI/TUI behavior:

- `parley sessions list`: list all known local sessions across workspaces.
- `parley sessions inspect <run-id>`: show manifest, agents, attempts, questions, artifacts, and resume support.
- `parley resume <run-id>`: compute the next recoverable action:
  - if a process is still running, reattach/log-tail where possible;
  - if an artifact is missing, resume/retry only that participant;
  - if an artifact is malformed, ask owning agent to repair its own file if the agent is reachable;
  - if all round files are valid, advance to consensus or next phase.
- TUI startup loads the global index, validates it against workspace run directories, and marks stale/running/failed/action states.
- TUI can continue a selected run with an explicit action menu: retry missing agent, resume external session, open manual handoff, answer HITL question, advance phase.

## Concerns / open questions

- External CLI cache formats are private and can change. Parley should record handles returned or configured at launch time rather than scraping private caches as a dependency.
- Headless mode and interactive mode often have different persistence behavior. Claude explicitly has `--no-session-persistence` for print mode; Gemini and Hermes expose resume flags, but exact headless continuation semantics need probes.
- Codex in this machine has multiple installations; `/opt/homebrew/bin/codex` is broken while the NVM Codex binary exposes `resume/fork`. Parley should capture the exact executable path used per attempt.
- We need a schema migration policy before moving from JSON to SQLite.
- Logs and prompts can contain sensitive information, so retention, redaction, and pruning need to be designed from the start.

## Risks

- Over-promising "continue exactly where the model left off" is risky. The robust promise is "continue the Parley workflow without losing state"; exact model conversation continuation is best-effort per agent.
- If Parley stores too much raw transcript content, it becomes another sensitive data vault. Prefer metadata and canonical artifacts.
- If resume mutates canonical files automatically, it can violate Parley ownership rules. Resume must operate through the owning agent or explicit user override.
- Parallel resumes can corrupt JSON files without locking. Use file locks for JSON or move to SQLite with transactions.
