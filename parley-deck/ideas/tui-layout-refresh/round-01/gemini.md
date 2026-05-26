---
agent: gemini
idea: tui-layout-refresh
round: 1
date: 2026-05-26
---

## Summary

The current TUI layout lacks a clear visual hierarchy and fails to use terminal space efficiently, especially on shorter displays. My proposal is to introduce a more deliberate, state-aware layout system that is responsive to terminal size. By defining clear regions, establishing a purposeful color-to-state mapping, and creating distinct layouts for standard and compact viewports, we can significantly improve scannability and usability without a full application rewrite.

## Proposed approach

I propose a focused refactoring of the TUI layout in `app.go` and `live.go`, centered around a more responsive and state-driven presentation.

### 1. Centralized Style and Layout Management

Instead of ad-hoc styling and width calculations within `View()` methods, we should centralize this logic.

- **Create `internal/tui/styles.go`:** This file will define the application's color palette and `lipgloss.Style` objects. This makes the UI theme consistent and easier to modify.
- **Define a Color-to-State Hierarchy:**
  - **Blue/Cyan:** Active states (`running`, `pending`, `starting`).
  - **Green:** Success states (`finished`, `ok`, `found`, `complete`).
  - **Red:** Failure states (`failed`, `error`, `missing`).
  - **Yellow/Orange:** Attention states (`blocked`, `confirm`, `question`, `warning`).
  - **Grey:** Inactive or informational states (`skipped`, `muted`, `unknown`).
- **Create a `badge(style, text)` helper:** This will render consistent, fixed-width status tags (e.g., `[RUNNING]`, `[ FAILED ]`) to improve alignment and scannability.

### 2. Dashboard View (`app.go`) Redesign

The current three-column layout is a good start, but it becomes cramped and vertically wasteful. I propose two layout modes based on terminal height.

- **Standard Layout (Height > 40 lines):**
  - A flexible three-column design.
  - **Column 1 (Context):** `Sessions` list. Use status badges and color to indicate which runs need attention.
  - **Column 2 (Activity):** A combined `Agent Status & Event Stream` panel for the selected session. This merges the current middle and right panels, showing agent progress and the latest events together, which is more cohesive.
  - **Column 3 (Detail):** A dynamic `Details / Actions` panel. This panel will show the `Log Preview` for the selected agent by default, but can be toggled (e.g., with `a`) to show the `Questions & Actions` list. This avoids having two large, separate panels for actions at the bottom.

- **Compact Layout (Height <= 40 lines):**
  - The layout switches to a two-column, horizontally-focused view.
  - **Column 1 (Sessions):** A highly condensed list of sessions, showing only the slug and a status badge.
  - **Column 2 (Workspace):** A tabbed view that takes the remaining width. The user can toggle between tabs for **"Events"**, **"Agents"**, and **"Actions"**. This eliminates all vertical stacking and makes the app usable on short terminals.

### 3. Live Run View (`live.go`) Redesign

The current view is "too vertical." I propose a persistent two-column layout that feels more like a traditional dashboard.

- **Left Column (Context - "The Who"):** The `Agents` table. This view is static and always visible, providing the core context for the run. The selected agent in this table dictates the content of the right column. Use color and badges for agent state.

- **Right Column (Detail - "The What"):** This column should be sectioned or tabbed to show details for the *selected agent*.
  - **Section 1: Log Preview:** Shows the `stdout`/`stderr` tail for the selected agent. This is the most frequently needed information and should be prominent.
  - **Section 2: Questions:** Shows the list of open HITL questions. If a question is selected, the answering prompt appears inline here, rather than at the bottom of the screen.

This redesign co-locates an agent with their corresponding output and actions, creating a much clearer cause-and-effect relationship in the UI.

## Concerns / open questions

1.  **Layout Switching Logic:** Should the switch between `Standard` and `Compact` layouts be based purely on height, or a combination of width and height? I propose starting with height, as that is the primary constraint mentioned.
2.  **Discoverability:** How will users discover the togglable panel in the dashboard's third column or the tabs in the compact view? The footer help text must be context-aware and clearly state the available keybindings (`a` to toggle actions, `tab` to switch tabs, etc.).
3.  **Information Density:** In the compact live view, will putting logs and questions in the same column feel too cramped? We may need to dynamically adjust the height allocated to each section.

## Risks

- **Scope Creep:** Refactoring the layout logic could tempt a broader rewrite. We must remain focused on reorganizing the existing components within the `View()` methods first, using a simple layout struct, rather than building an abstract layout engine.
- **Visual Noise:** A more liberal use of color can backfire if not applied with discipline. The proposed color-to-state mapping must be adhered to strictly to ensure colors add clarity rather than chaos.
- **Performance:** A more dynamic layout that recalculates on every `tea.WindowSizeMsg` could introduce lag or flickering. We should leverage Lip Gloss's built-in layout features, which are generally efficient.
