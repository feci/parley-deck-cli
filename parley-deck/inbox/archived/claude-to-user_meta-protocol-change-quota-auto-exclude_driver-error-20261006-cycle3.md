---
from: claude
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: review
blocking: yes
date: 2026-10-06
---

## Question / blocker
The auto-driver halted with an error while advancing round-03:

    draft review consensus: review-consensus drafter claude-1: context canceled

Inspect the run (events.jsonl / agent logs), fix the cause, then re-run 'parley run --auto'.
