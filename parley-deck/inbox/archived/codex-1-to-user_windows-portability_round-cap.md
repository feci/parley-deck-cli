---
from: codex-1
to: user
idea: windows-portability
phase: round-04
blocking: yes
date: 2026-09-25
---

## Question

May the quorum run one additional focused cross-review round to resolve Windows durability before consensus? Organizer recommendation: yes, limited to the prepared durability-followup-brief.md, with all material findings still reportable.

## Context

All three real CLI participants completed round 4 with exit 0; final artifacts pass parley 1.49.0 structural validation. No product implementation, FINAL, consensus, signoff or release has occurred on this idea. The organizer remains pure and has not verified product code.

Five design areas have converged. Durability remains materially contested: Claude withdrew its earlier O_SYNC create-entry inference, identifies thirteen barrier operations behind six call sites, and requires explicit per-operation mechanisms or named refusals including feature-level impact. Kimi and Zcode accept the earlier inference conditionally and a six-site contract. Their round-04 artifacts explicitly target round 3; neither answers Claude's new correction. A majority is not a resolution under section 15.3.

The exact proposed additional work is in ../ideas/windows-portability/durability-followup-brief.md. It requests an evidence-backed operation table, preservation of containment and anti-replay properties, and an explicit restructure-or-refuse design. It does not ask the owner to settle Win32 semantics or authorize a weaker guarantee.

The procedural reason for this request is the live COOPERATION.md section 4.0 deliberation cap: cross-review is "capped at 3 after round 1, then escalate". Rounds 2, 3 and 4 exhausted it. The parley-deck skill instructs the organizer to follow the live phase rules. Existing lifecycle/release authorization is preserved; it does not explicitly waive this round cap. The mechanical wait digest says "await consensus signoff", but raw Claude round 4 has an explicit remaining blocker and overrides that suggestion.

Release prerequisite 1 is present (CLI 1.49.1 predecessor handoff); prerequisite 2, the designated-implementer completion handoff, remains absent at this boundary. Windows is still experimental and CLI winget held. No Windows runtime acceptance claim is made.

## What I need from you

Authorize one additional focused cross-review round (round 5). Then the existing quorum can resolve the new evidence, Zcode can draft the agreed specification and implement, and Claude/Kimi can independently review. All existing hosted-CI, release sequencing and channel-verification gates remain in force.
