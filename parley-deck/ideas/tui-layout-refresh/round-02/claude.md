---
agent: claude
idea: tui-layout-refresh
round: 2
date: 2026-05-26
status: approved
---

## Signoff

Approved. The consensus captures the right scope.

My round-01 proposal pushed further toward a strip+detail-pane model that would have been a larger rewrite. The consensus instead keeps a two-column layout with normal and compact branches, which addresses the core height problem without the discoverability and scope-creep risks I flagged. Good call.

Two minor notes (non-blocking):

1. **Compact threshold.** The consensus says "height-aware compact rendering" but doesn't name a specific breakpoint. Suggest `height < 30` as the switch, consistent with my round-01 analysis. Implementation detail, not a design question.

2. **Question auto-focus.** The HITL question panel is the most time-sensitive element. The consensus preserves current keyboard behavior, which is fine, but during implementation it would be worth ensuring that an arriving open question surfaces visibly in both normal and compact modes (e.g., an amber badge on the questions section title) so the operator doesn't miss it while focused elsewhere. This is already implied by the semantic color rules but worth calling out as a test case.

No blocking concerns.
