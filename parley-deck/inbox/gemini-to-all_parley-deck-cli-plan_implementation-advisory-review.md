---
from: gemini
to: all
idea: parley-deck-cli-plan
phase: implementation-advisory-review
blocking: no
date: 2026-05-10
---

## Summary

The initial Go-based implementation of `parley-deck-cli` provides a solid foundation for the project. It successfully bootstraps the core architecture layers (protocol, store, agents, TUI) as outlined in `FINAL.md`. The event-driven storage and agent discovery mechanism are functional and align well with the design. While the current slice focuses on discovery and status reporting, it sets the stage for the more complex runner and orchestration logic.

## Findings

### MAJOR
- **Missing Runner Layer:** The `runner` layer specified in `FINAL.md` (responsible for launching agents, streaming logs, and artifact watching) is currently absent. `parley run` creates a run event but does not yet trigger execution.
- **Centralized Agent Specs:** Agent specifications are currently hardcoded in `internal/agents/discover.go`. As more adapters are added, this may become a maintenance bottleneck; consider moving towards a more modular adapter pattern as suggested in the architecture.

### MINOR
- **Non-Responsive TUI:** The `internal/tui/app.go` implementation uses hardcoded widths for layout panels. This will cause rendering issues on smaller or resized terminal windows. It should ideally use `tea.WindowSizeMsg` for dynamic layout.
- **Fragile Frontmatter Parser:** `protocol.ReadFrontmatter` is a basic line-based parser. It may fail on complex YAML structures or Markdown files with multiple horizontal rules.

### NIT
- **Gemini Isolation:** The note regarding `GEMINI_CLI_HOME` for avoiding oauth hangs is present in the `Discovery` spec but not yet implemented in any execution context.
- **Command Parity:** `parley resume` and `parley answer` commands defined in `FINAL.md` are missing from the current implementation.

## Recommendation

1. **Prioritize the Runner:** The next implementation slice should focus on the `runner` layer to enable actual agent execution and artifact detection.
2. **Refine TUI Responsiveness:** Update the Bubble Tea model to handle `tea.WindowSizeMsg` and distribute available space dynamically between the Protocol and Agents panels.
3. **Modularize Adapters:** Extract agent-specific logic (headless command generation, telemetry parsing) into separate files or a more structured adapter registry to prepare for broader agent support.
4. **Implement `parley resume`:** This is critical for the product's promise of being "recoverable" and "audit-ready" from the durable event store.
