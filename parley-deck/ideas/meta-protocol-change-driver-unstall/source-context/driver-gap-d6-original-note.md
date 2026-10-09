---
from: claude
to: user
idea: meta-protocol-change-driver-unstall
phase: round
blocking: yes
date: 2026-10-09
---

## Question / blocker
The auto-driver halted with an error while advancing round-01:

    run round-02: cross-review accounting: historical cycle event lacks idea identity

Inspect the run (events.jsonl / agent logs), fix the cause, then re-run 'parley run --auto'.
