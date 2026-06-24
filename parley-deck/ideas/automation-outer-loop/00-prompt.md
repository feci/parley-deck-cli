---
idea: automation-outer-loop
author: claude-1
created: 2026-06-24
status: round-01
participants: [claude-1, codex-1, hermes-1, antigravity-1]
auto_implement: true
checks: "go test -count=1 ./internal/loop/ ./internal/app/"
derived_from: [ideas/loop-engineering-research/FINAL.md]
---

## Problem / idea

Tier 4 of the loop-engineering backlog: the **outer loop**, built only on top of the
verification-honesty (Tier 1), loop-budget (Tier 2), and close-integrity (Tier 3) safety
floor. Two backlog items:

- **LE-8 — Human-brake invariant** (`automation-human-brake`). A new protocol section that
  binds any automated/standing loop: it may **discover + draft (Phase 0/1 candidates)
  only**. It must NEVER push to quorum, implement, land/merge, finalize, modify the roster,
  or override/bypass consensus without a recorded human or full-quorum gate. This generalizes
  the §12.11 monitoring-watcher candidate invariant into a protocol-wide brake.

- **LE-9 — `parley loop tick`** (`standing-loop-watch-mode`). A one-shot, scheduler-friendly
  command (cron / GH-Actions / MCP can call it) — **not a daemon**. It discovers candidate
  signals (commits / CI / issues / a signals file), dedupes them against existing ideas, and
  writes provenance into `status: candidate` idea prompts. It is **disabled by default**,
  and even when enabled it only drafts candidates — it never stands up a quorum, never runs
  `parley run`, never pushes or merges. Any promotion to an active deliberation is a separate,
  human-gated step (LE-8).

## Constraints

- LE-8 is binding protocol text and must land in BOTH `COOPERATION.md` copies (drift guard).
- LE-9 must be disabled-by-default and fail-safe: with no config it writes nothing and exits 0.
- Reuse the established `status: candidate` candidate-idea shape (the §12.11 watcher pattern):
  no `participants:` quorum claim, a `## Promotion` note, human/manifest promotion required.
- No daemon, no background scheduling inside the CLI; `tick` is one-shot and returns.
- `tick` must never invoke `parley run`, push, merge, finalize, or edit the roster.

## Non-goals

- Auto-running deliberations from discovered candidates (that is the human-gated promotion).
- Real connector integrations (GitHub issues API, CI webhooks) — the MVP reads a signals file;
  live connectors are a follow-up.
- LE-12 cross-run durable goal state (its own contested idea).
