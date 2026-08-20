---
idea: tui-round-summary
author: user
created: 2026-07-04
track: standard
participants: [claude-1, codex-1, hermes-1]
roles:
  claude-1: facilitation + TUI architecture fit (live.go patterns)
  codex-1: event/state correctness (round completion detection)
  hermes-1: information design — what belongs in a digest
status: final
---

## Problem / idea

Inspired by Hermes Agent v0.18.0 background fan-out: parallel subagents return ONE
consolidated turn instead of the user stitching results from separate windows.

Parley TUI analogue: during a live run, round agents work in parallel and the user
watches per-agent tabs. When a round completes there is no single consolidated view —
the user assembles the picture by flipping tabs. The Home tab shows protocol chips and
run state, but not "what just happened in round N".

Proposal to deliberate:

- **Consolidated round summary in the Home tab**: when the driver detects a round is
  complete (all expected round-NN files exist), render a digest block:
  - per agent: one-line position extract (from its round file Summary section),
  - convergence signal: counts of ACCEPT/counter-proposal/blockers where detectable,
  - what happens next (per driver state: next round opening / consensus drafting).
- **History**: keep the last few round digests visible (scrollable) so a returning
  user can catch up without tab-flipping.
- Reuse the existing driver events (run.phase etc.) and narrator infrastructure from
  the protocol-visibility work; no new protocol artifacts — this is a pure TUI/driver
  presentation feature.

## Constraints

- Read-only over canonical files: the digest derives from round files; it never
  writes protocol artifacts.
- Extraction must be robust to imperfect files (missing Summary section ⇒ fall back
  to first prose paragraph, capped length).
- No LLM calls for summarization in v1 — mechanical extraction only (deterministic,
  free, offline).
- Must not regress existing Home tab content (chips, roster, runs list).
- CLI-only change; no protocol text change.

## Non-goals

- No cross-round analytics or charts.
- No LLM-generated executive summaries (possible follow-up).
- No change to per-agent tabs or transcript view.
