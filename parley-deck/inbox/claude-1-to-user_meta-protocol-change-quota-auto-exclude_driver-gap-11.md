---
from: claude-1
to: user
idea: meta-protocol-change-quota-auto-exclude
phase: round-03
blocking: no
date: 2026-10-03
---

## Question

This run needs no action on the driver's new blocking note,
`claude-to-user_meta-protocol-change-quota-auto-exclude_driver-error.md` (19:09 CEST). Do you want to
make the accounting decision that lets the driver launch cross-review rounds in this repository again?

## Context

- The driver halted before it launched round 3, with "cross-review accounting: historical cycle event
  lacks idea identity".
- The cause is one legacy record: line 1 of `parley-deck/runs/20260510T194003Z/events.jsonl`. It is the
  `run.created` event of the May smoke run, committed in `3ec10ac`, and it has no `idea` field.
- The budget gate refuses that record whenever an idea has no cross-review cycle policy yet. So the first
  driver-launched cross-review round of every new idea in this repository stops here. Until your prune,
  the stale-worktree check ran first and hid this one.
- The CLI's remedy is an operator decision, run from an attended terminal:
  1. `parley budget migrate inspect --kind cross-review --idea <slug>`;
  2. `apply`, with `--declare-unscoped-run parley-deck/runs/20260510T194003Z=<manifest digest>`, a
     decision id, a reason, `--writers-stopped` and explicit ceilings.

  The declaration covers that one decision only and is not stored, so each new idea would need it again.
  A lasting fix, such as changing the legacy record or how the gate treats runs without an idea, would be
  a separate change. The organizer runs none of these.
- This run continues without the fix. Round 3 runs through the recorded runner-template fallback
  (launched 19:14 CEST). The driver will be retried for consensus, signoffs and FINAL, which do not pass
  this gate.

## What I need from you

Nothing for this run. Optionally, choose between the per-idea migration and a separate lasting fix. The
proposal note at the end of this run lists this as an owner action item.
