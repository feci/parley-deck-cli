---
agent: agy
idea: tui-protocol-visibility
round: 2
phase: review
date: 2026-06-12
---

## Fix verification

- **Home Attention badge (D14):** VERIFIED. In [live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1063-L1072), [renderHome](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1046-L1094) successfully resolves the `Attention` status from the idea's latest run in the `homeRuns` list and appends `attentionBadge(r.Attention)` next to the phase chip.
- **Placeholder indentation (D12):** VERIFIED. In [protocolui.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L677-L680), the `writeLine` helper inside [renderSilentPlaceholder](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L673-L712) prepends two spaces (`"  "`) to every printed line, maintaining alignment.
- **Step-8 pipeline title:** VERIFIED. In [protocolui.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L606-L611), [renderProtocolPanes](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L587-L668) looks up the title using the snapshot's actual step name `p.StepName` when `i == p.Step`, allowing step 8 to render as "Fix-Up" during fix-up cycles.
- **Collapsed ribbon round format:** VERIFIED. In [protocolui.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L428-L430), [renderRibbon](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L410-L463) formats active rounds as `(R%02d)` without the total rounds denominator.
- **Narrator/activity noise:** VERIFIED. In [protocolui.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L216-L265), the [friendlyEventText](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/protocolui.go#L216-L265) function maps raw event types to concise action verbs or drops the type prefix entirely when the summary is self-describing.
- **Home double status:** VERIFIED. In [live.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1058-L1061), [renderHome](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/tui/live.go#L1046-L1094) renders a neutral `"—"` chip instead of duplicating the raw idea status when the phase chip is empty.

## New findings

None.

## Verdict

ACCEPT
