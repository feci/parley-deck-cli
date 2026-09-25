---
from: codex-1
to: user
idea: windows-portability
phase: consensus
blocking: no
date: 2026-09-25
---

## Question

Do you explicitly approve the containment deviation described below? Until approval, the current design preserves containment and uses the explicit refusals authorized in the original brief. Consensus preparation proceeds on that current design; silence is not approval and no weakened operation will be implemented by default.

## Context

Round 6 completed sequentially, all three real CLIs exited 0 and all artifacts passed structural validation. Their current positions converge; no additional reconciliation round is requested. This note summarizes participant testimony, not an organizer code-verification verdict. Canonical sources: ideas/windows-portability/round-06/{claude-1,kimi-1,zcode-1}.md, especially Conflict 2 and Zcode's normalized table.

The current authorized design retains os.Root containment and refuses seven rooted publication/mkdir/rename sites before mutation (A2, A3, B3, B4, B5, C1, C2), plus dependent re-sync operations D2/D3/D4a/D4b. The participant-reported reachable Windows cost is precharge reservation-intent publication and parent-recovery / recovered-parent apply. Trajectory verification and captured journals already have Windows refusals. Plain-path sites can use the reviewed staging design without this deviation.

The proposed alternative would enable the rooted operations with path-based MoveFileEx(MOVEFILE_WRITE_THROUGH), using paths derived from Root.Name(), lexical sibling checks, product-constant or validHash-gated names, a publish-time key check at A3, deterministic retained class-A staging and randomized class-C staging, and explicit documentation. Zcode proposes no-REPLACE at C1/C2 to preserve publish-once intent; the exact contract still requires participant signoff.

The cost is explicit: publication loses os.Root's kernel-enforced per-open reparse resistance and handle-relative rename. A lexical path check plus later path-based move admits a time-of-check/time-of-use race if path components change. Root-derived paths and sibling checks are not equivalent containment. The proposed hosted escape test checks mechanics; it does not restore that guarantee or prove durability. A handle-derived-path alternative was named but not selected or designed; participants report it narrows rather than closes the race.

Zcode prefers the deviation to recover the affected features. All three accept the current refusal branch as a reviewed outcome within the original owner scope. The organizer recommends retaining the existing containment requirement unless the feature benefit justifies this expressly reduced guarantee to the owner.

## What I need from you

Approval is needed only to select the weaker-containment branch. The precise authority is your answer in user-to-codex-1_windows-portability_round5-unresolved.md: "Any proposed weakening of os.Root containment is not accepted by this answer: it must come back to the owner as an explicit deviation." This note is that explicit deviation request, not a claim that the answer already permits it. All ordinary consensus, implementation, review and release gates continue to apply to either branch.
