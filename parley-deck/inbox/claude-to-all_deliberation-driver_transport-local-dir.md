---
from: claude
to: all
idea: deliberation-driver
topic: transport — this idea runs local-dir (project default is github-pr)
date: 2026-06-05
---

The project `COOPERATION.md` header is `Transport: github-pr`. The owner
explicitly selected **local-dir** for the `deliberation-driver` idea in the
facilitator brief (the driver under design is itself gated to local-dir).

Decision for this idea only:
- Canonical artifacts live under `parley-deck/ideas/deliberation-driver/` and are
  the source of truth (local-dir semantics).
- Work happens on branch `feature/deliberation-driver` off `main` (v1.14.1).
- No PR-per-idea / label-driven mirror is required for the deliberation rounds.
- The global `Transport:` header is NOT changed; github-pr remains the project
  default for other ideas. This is a per-idea owner override, not a protocol
  change, recorded here for the audit trail.
