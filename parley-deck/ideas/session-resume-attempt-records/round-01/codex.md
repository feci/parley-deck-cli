---
agent: codex
idea: session-resume-attempt-records
round: 1
date: 2026-05-25
---

## Summary

Slice 2 should introduce a narrow attempt-record package and wire it only into
agent execution paths that already know the run ID, agent ID, expected artifact,
model/config, and process result. The first implementation should avoid broad
resume semantics and focus on durable diagnostics that `sessions inspect` can
show.

## Proposed approach

- Add `internal/attemptrecord` with schema version 1 and helpers to write
  `parley-deck/runs/<run-id>/attempts/<step-id>-<agent-id>-attempt-<n>.json`.
- Record:
  - run ID, idea slug, step ID, agent ID, expected artifact path;
  - prompt SHA-256 and byte length;
  - selected model/profile/effort when known;
  - redacted command path/args;
  - start/end timestamps, duration, exit code or timeout flag;
  - bounded stdout/stderr tails;
  - validation status: missing, exists, parse-error, valid.
- Store full rendered prompts in local cache only, under
  `~/.parley-deck/runs/<run-id>/input-packs/`, with the same prompt hash in the
  workspace record.
- Extend `sessions inspect` to list attempt records when present.
- Fold minimal run status updates into this slice only where needed to make
  inspect output coherent; do not attempt terminal lifecycle completeness yet.

## Concerns / open questions

- The existing runner may not have a stable step ID in every path. If so, use a
  deterministic first version such as `round-01` or `review-round-01`.
- Redaction should start conservative: strip obvious env assignments and avoid
  storing raw prompt text in workspace records.
- Validation should remain shallow in v1: expected file exists and frontmatter
  parses where applicable.

## Risks

- Over-capturing stdout/stderr can leak prompt or credential fragments.
- If attempt writes happen before run manifests are guaranteed, inspect output
  may become inconsistent.
- Naming churn in attempt records would be painful once users have local cache
  history, so the path scheme should be chosen carefully.
