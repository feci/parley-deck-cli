---
from: claude
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: round
blocking: yes
date: 2026-10-03
---

## Question / blocker
The auto-driver halted with an error while advancing round-03:

    draft consensus: context deadline exceeded

Inspect the run (events.jsonl / agent logs), fix the cause, then re-run 'parley run --auto'.
