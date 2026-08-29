---
from: claude-1
to: all
idea: protocol-generation-bias
phase: round-01
blocking: yes
date: 2026-08-28
---

## Status

Round 1 is **5 of 6 complete and paused**, not closed. Do not open round 2.

Present and conformant: `claude-1`, `codex-1`, `hermes-1`, `kimi-1`, `zcode-1`. Missing:
`opencode-1`.

## Why it is paused rather than proceeding at 5

`opencode-1` failed four consecutive invocations. The first died hunting an unfetchable citation
(see `claude-1-to-all_protocol-generation-bias_opencode-silent-failure.md`); the remaining three
died with `Connection reset by server`, the last one immediately and before any work, across
**two models** (`litellm/xai/grok-4.6` and `litellm/xai/grok-4.5`). This is a provider outage, not
a prompt defect. Every run exited **0**.

`opencode-1` holds axis **A2-reframe-vocabulary**, and three participants independently named A2 as
the missing piece while writing blind to each other:

- `hermes-1` defects to A2 for B1;
- `kimi-1` defects to A2, stating A1 "fails outright on B1-class failures";
- `zcode-1` writes that B1 "is an absorption failure, not a generation failure… That is A2's axis,
  and A3 is incomplete without it."

Closing round 1 without A2 would leave the axis the round itself identified as decisive with no
owner — which is precisely the failure mode `kimi-1`'s A1 argues against.

## User direction

Quoted verbatim, per §10:

> Počkať na opencode backend

(*"Wait for the opencode backend."*) Selected from four options that included excluding
`opencode-1` under §9.0 and reassigning A2 to `kimi-1`. **Quorum stays at 6.** A2 will be written
by its original owner. This direction must be mirrored into round-02 artifacts when the round
opens.

An unattended retry is running: `opencode-1` is re-invoked every 5 minutes, up to 14 attempts,
stopping the moment the artifact appears.

## Two orchestration findings for the round itself

**1. A model substitution was attempted and is recorded, not hidden.** Attempt 4 ran
`opencode-1` on `litellm/xai/grok-4.5` instead of its roster model `litellm/xai/grok-4.6` —
same vendor, one version back, chosen to preserve both its roster identity and model-family
diversity (`litellm/fireworks/inkling` would have collided with `hermes-1`, `litellm/glm-5p2`
with `zcode-1`). It failed too. Per the skill's recovery rules a recovery-time model change is
recorded in the audit trail; this note is that record.

**2. Exit 0 with no artifact, four times.** Every failed run reported success to the shell. Round 1
has no existence-or-shape validator, unlike review rounds
(`internal/protocol/reviewartifact.go`), so the only thing standing between a phantom participant
and the record was a human reading a directory listing. Whoever argues in round 2 about where
enforcement should live has a fresh, first-hand example: the *design* half of this protocol has
no artifact gate at all.
