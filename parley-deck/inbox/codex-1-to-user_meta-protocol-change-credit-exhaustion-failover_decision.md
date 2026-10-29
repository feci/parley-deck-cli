---
from: codex-1
to: user
idea: meta-protocol-change-credit-exhaustion-failover
phase: round-02
blocking: yes
date: 2026-10-29
---

## Question

Authorize the native CLI fallback for this run, or repair the historical-worktree gate before peer cross-review.
I recommend the bounded native fallback. It preserves the brief's prohibition on worktree changes.

## Context

The resume steps are complete. The updated brief and exact succession ruling are in the idea.
Codex round 1 retains its original bytes. Both peer round-1 files are complete and valid.
Codex round 2 applies the owner ruling and responds to both peers.

Both measured peer round-2 launches refused before a child started. Their failure class is `budget_refused`.
Their PID and start time are null. Neither failure is provider exhaustion or a provider attempt.
Kimi invocation: `21e634a2-1db5-42ad-99c7-ab6b4d9cbb6d`.
Zcode invocation: `7c900416-094d-41b1-9763-f3a99337b3d4`.

Read-only cycle inspection reports no initialized cycle policy.
Read-only migration and legacy inspection both report:

```text
historical worktree is unavailable: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer
```

The unavailable worktree gate persists despite the brief's report that D6 was applied.
The brief states, "Do not prune them, and do not declare worktrees."
No registration or accounting state was changed. No peer round-2 child is active.

The prepared peer briefs are private task files:

- `.parley-runtime/credit-failover/round02-kimi-1-brief.txt`
- `.parley-runtime/credit-failover/round02-zcode-1-brief.txt`

The discovered launch forms are:

```text
/Users/tomasfecko/.kimi-code/bin/kimi --output-format stream-json -m omniroute/k3.1-pro -p <brief>
/Users/tomasfecko/.nvm/versions/node/v24.15.0/bin/zcode --prompt=<brief> --mode yolo --cwd <CLI worktree>
```

The fallback would use each agent's existing configuration and the same 1,800-second ceiling.
Kimi effort remains max from its native configuration. Zcode remains `zai/glm-5.3` with max effort from native configuration.
Each process would own its canonical file. Private logs and content-free receipts would record actual times, exit, and artifact hash.
The fallback would not claim driver accounting or automatic membership authority.
All later consensus, review, independent goal checks, and release requirements would remain binding.

The Parley Deck skill's Recovery And Partial Completion section requires a recorded failure and a user recovery choice.
Source: `/Users/tomasfecko/.codex/skills/parley-deck/SKILL.md`.
Its relevant instruction is:

> "For non-zero CLI exit, rate limit, auth failure, empty output, or timeout, capture the failure in an inbox note and ask the user whether to retry, replace the participant, extend the deadline, or continue under the protocol's silence/deadline rule."

The brief pre-authorizes retries for specified provider errors and timeouts. This accounting refusal is neither case.
An explicit native-fallback choice resolves the remaining authority question without a worktree declaration.

## What I need from you

Choose one option:

1. **Recommended:** authorize native CLI fallback for this idea wherever this recorded driver gate prevents a required participant action.
   Retain the brief's models, timeouts, retry limits, independent roles, private logs, and audit receipts.
2. Repair the historical-worktree gate through the owner's normal process, then tell codex-1 to resume measured dispatch.

No design consensus, product implementation, or release is complete. The design PR is https://github.com/feci/parley-deck-cli/pull/86.
Safe evidence is in the idea's `source-context/round02-*-refusal.json` and `source-context/round02-cycle-gate.txt`.
