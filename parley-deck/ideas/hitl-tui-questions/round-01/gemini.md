---
agent: gemini
idea: hitl-tui-questions
round: 1
date: 2026-05-11
---

## Summary
The goal is to introduce a Human-In-The-Loop (HITL) capability to Parley Deck, allowing agents or the system runner to pause execution and solicit user input. This requires a durable messaging mechanism for questions and answers, a specialized UI state in the Bubble Tea TUI, and a fallback mechanism for non-interactive environments.

## Proposed approach

### 1. Event-Driven Q&A Protocol
Extend the existing event system to include `QuestionRequested` and `AnswerProvided` event types.
- **QuestionRequested**: Contains a unique question ID, the source agent, the prompt text, and optional default/auto-answer values.
- **AnswerProvided**: Maps to a specific question ID and contains the user's response.

### 2. Durable Persistence
Questions and answers should be appended to the `events.jsonl` file in the run directory. This ensures that the state of the conversation is preserved even if the TUI is closed or the process is interrupted.

### 3. TUI Implementation (Bubble Tea)
- Introduce a `QuestionView` component that activates when a `QuestionRequested` event is encountered.
- This view should overlay or replace the main log view, displaying the prompt and an interactive text input.
- Key bindings (e.g., `Enter` to submit, `Esc` to potentially ignore or use default) should be clearly indicated.

### 4. Runner Logic
- The runner should block further execution for the requesting agent (or the entire round, depending on dependencies) until an `AnswerProvided` event is received.
- If `--auto` is provided, the runner should immediately emit an `AnswerProvided` event using the default value or a generic "auto-approved" response, bypassing the UI block.

### 5. CLI Fallback (`--no-tui`)
- In `--no-tui` mode, the runner can use a simple standard I/O prompt (e.g., using `fmt.Scanln`) if a TTY is detected.
- For truly non-interactive environments, we could implement a mechanism where the runner waits for a specific file (e.g., `run/questions/<id>.answer`) to appear, or simply fails if no answer is provided and `--auto` is not set.

## Concerns / open questions
- **Multi-Agent Questions**: If multiple agents ask questions simultaneously, how should the TUI prioritize or queue them? A simple first-come-first-served queue is likely best for this slice.
- **Agent Protocol**: How exactly does an agent (which might be a sub-process) communicate the need for a question? We may need a specific "magic string" or a structured JSON message on stdout that the runner intercepts.
- **Validation**: Should the system support basic validation (e.g., "yes/no" only) at the TUI level?

## Risks
- **Deadlocks**: The runner might wait indefinitely for an answer that never comes, especially in headless environments without `--auto`. We should consider a timeout or a clear "pending question" status in the logs.
- **TUI State Transitions**: Managing the transition between the scrolling log view and the interactive question prompt in Bubble Tea can be tricky; we must ensure the log continues to update in the background if applicable.
- **Persistence Bloat**: While `events.jsonl` is fine for now, we should ensure large answers (if supported) don't degrade performance.
