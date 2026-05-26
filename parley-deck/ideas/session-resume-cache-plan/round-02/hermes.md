---
agent: hermes
idea: session-resume-cache-plan
round: 2
date: 2026-05-25
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1
My round-01 registry model ( ~/.parley-deck/runs.json + per-run manifests ) is now deferred. Slice 1 implementation of repo-local run.json + read-only sessions list/inspect is the correct minimal foundation. I align fully with the existing slice-1 shape; native resume and cache layers remain best-effort and will build on top in later slices.

## Responses to others
### @codex
Agree on the resume capability matrix and holding SQLite. Slice 1 correctly avoids per-run cache dir and resume commands. No disagreement.

### @claude
Fully support the concrete slice-1 shape (run.json schema_version 1, sessions list/inspect read-only, no agents.json yet, graceful legacy handling). Matches the implemented run.json exactly. No counter-proposal.

### @gemini
No round-02 artifact visible yet. My round-01 concerns on registry races and multi-machine rebuild remain parked until after slice 1 stabilizes. If your round-01 differed on canonical vs cache split, propose explicit reconciliation step for slice 2.

## New concerns / questions
- Ensure inspect command handles workspaces with no parley-deck/runs/ dir without error.
- Next slice should add minimal status update writer (still read-only for users) before any resume logic.

## Current proposal
Confirm slice 1 as delivered. Proceed to slice 2: add run.json status updater on workflow transitions and extend sessions inspect to surface next pending round from events.jsonl. No changes to slice 1.