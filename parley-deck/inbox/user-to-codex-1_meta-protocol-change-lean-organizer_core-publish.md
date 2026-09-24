---
from: user
to: codex-1
idea: meta-protocol-change-lean-organizer
date: 2026-09-23
status: authorized
relayed-by: the owner's Claude Code session (not a participant artifact)
---

# Owner decision: one combined protocol core at the end of this run

## Facts (measured 2026-09-23 with `parley protocol status` and `ls ~/.parley/protocol/core`)

- The only PUBLISHED global core is **2.10.0** (`~/.parley/protocol/core/2.10.0/COOPERATION.md`, 105,981 B).
- Core **2.11.0** (the §15.6 change, CLI 1.47.0 / skill 2.11.0) has been STAGED since 2026-08-31 at
  `~/.parley/staging/COOPERATION-2.11.0.md` (105,575 B) and was never published.
- The protocol changes released in CLI **1.48.0** / skill 2.12.1 (COOPERATION.md 105,731 -> 108,400 B:
  §9 item 1, the goal-check "can only withhold a close" paragraph, the LE-7/LE-11 bullet) are in NO core.

## The owner's words (verbatim, Slovak) and translation

Offered the choice "publish 2.11.0 now, or wait and publish one combined core at the end", the owner answered:

> "počkám, publikuj to celé na konci"

"I will wait; publish all of it at the end."

## What this means for the release step

- Do not publish or ask the owner to publish 2.11.0 separately.
- Stage ONE core that contains all three change sets on top of 2.10.0: the 2.11.0 hunks, the 1.48.0
  protocol hunks, and this idea's reviewed hunks. Build it from the core TEMPLATE (placeholder header,
  stub §2), not from a deck view. Suggested version: the skill version this run releases (2.13.0).
- Verify before escalating: the §15 region and every changed section match the released deck view at the
  release commit; the header is still placeholders; diff core 2.10.0 -> staged core contains exactly the
  three change sets and nothing else.
- `parley protocol publish` stays attended-only. Escalate the exact command to the owner; the owner runs it.
