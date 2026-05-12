---
agent: gemini
idea: runtime-status-resume
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

I am dropping my earlier proposals for "Retry Failed Agents" and "PID-based process checking" to favor a smaller, more honest observability-first slice. I agree with the emerging group consensus that `parley resume` should initially be a read-only re-attachment to durable state. Retrying and supervision are better handled in subsequent design cycles once the core projection layer is stable. I am also dropping the `--watch` proposal to prioritize getting the snapshot views right first.

## Responses to others

### @parley-deck/ideas/runtime-status-resume/round-02/codex.md — round-01

I strongly support your convergence on Claude’s approach in your Round 2 position. Dropping the PID and retry features for this slice simplifies the implementation and avoids over-promising liveness. Your refinement on treating `round.incomplete` as terminal but distinguishing it from `complete` in the UI is a valuable UX detail that we should include in the FINAL plan to clarify why a run stopped.

### @parley-deck/ideas/runtime-status-resume/round-01/claude.md — round-01

I agree with the `internal/runstate` package lift. It is the most robust way to ensure the TUI and CLI views stay in sync. Your `Liveness` model (`complete`, `idle`, `unverified`) is the right conservative approach given we don't have a supervisor. I also appreciate your observation that HITL answers still work in "read-only" resume mode—this is a powerful side-effect of our durable file architecture that we should explicitly call out as a feature.

### @parley-deck/ideas/runtime-status-resume/round-01/hermes.md — round-01

I support your push for compact human-readable output. While I am now siding against including the PID file in this slice, I agree that our liveness claims must remain "strictly advisory." Your suggestion to walk the `parley-deck/` tree and merge protocol state with event data is exactly what the `internal/runstate` package should do. I agree with your concern about mapping runs to ideas; relying on `run.created` is the right path, but we should ensure a graceful fallback.

## New concerns / questions

- **Event Log Performance:** As `events.jsonl` grows, reading the whole file for every `status` call might become slow. While we shouldn't over-optimize now, we should ensure the `runstate` package is designed to eventually support an offset cache or partial tailing.
- **Run-to-Idea Mapping Fallback:** We should verify that all runs in the wild already have the `idea` field in `run.created`. If not, we should design the fallback (e.g., checking artifact headers or prompt files) to ensure `status` doesn't show "unknown" for recent work.

## Current proposal

I propose we converge on a FINAL plan that implements the following:

- **Shared Projection:** Create `internal/runstate` to project `events.jsonl`, `run.created`, and HITL questions into a unified `RunSummary`. This package will be used by both `status` and `resume`.
- **`parley status`:**
    - Provide a workspace overview showing the active Transport and a count of open HITL questions across all runs.
    - List recent runs (default newest 5-10) with their idea slug, phase, agent matrix summary, and liveness.
    - Support `--run <id>` and `--idea <slug>` for detailed views (full agent state, log paths, recent events, and next actionable commands).
- **`parley resume`:**
    - Resolve target via exact run-id or the latest run for an idea slug.
    - Support TUI (read-only replay of events) and `--no-tui` (mirroring the `status --run` view).
    - Allow HITL answers via the existing `a` keybinding/command, as it writes through to durable storage.
- **Conservative Liveness:** Adopt Claude's `complete/idle/unverified` model to avoid over-promising process reattachment.
- **Testing:** Use deterministic event stream fixtures to verify the `runstate` reducer and CLI output.
