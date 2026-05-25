---
agent: gemini
idea: tui-action-execution
round: 1
date: 2026-05-25
---

## Summary
The current TUI displays `run.NextActions` in a read-only fashion. The objective is to enable users to select and execute these actions directly within the TUI. This analysis focuses on the contract for action commands, their safe execution, and the corresponding test strategy. The solution should integrate action execution capabilities into the existing TUI framework, ensuring clarity for users, especially regarding unsupported or interactive actions, and providing immediate feedback through UI refreshes.

## Proposed approach
The approach involves extending the TUI's capabilities to allow users to navigate and trigger `run.NextActions`.
1.  **TUI Action Selection and Display:** Modify `internal/tui/app.go` to introduce state for tracking the currently focused/selected `run.NextAction`. Update `renderQuestions` to visually highlight the selected action and provide UI controls (e.g., keyboard shortcuts) for selection and execution.
2.  **`ActionRunner` Interface:** Define an `ActionRunner` interface in `internal/app/app.go` or a new dedicated package. This interface will expose a method (e.g., `ExecuteAction(action runaction.NextAction) error`) responsible for translating and running the appropriate CLI command.
3.  **Dependency Injection:** `internal/app/app.go`'s `runTUIViewWithDiscovery` will be responsible for instantiating an implementation of `ActionRunner` and injecting it into `tui.RunWorkspace` via `WorkspaceOptions`, similar to `RefreshRuns` and `StartRunFunc`.
4.  **Action Command Contract:** The `ActionRunner` implementation will utilize `internal/runaction/action.go`'s `actionCommand` mapping to identify executable actions. Only actions with a defined, non-empty command should be considered for direct execution. Actions like `retry-agent`, which have no executable command, should be explicitly handled as non-executable.
5.  **Execution and Feedback:** Upon successful execution by the `ActionRunner`, a `RefreshRuns` operation will be triggered to update the TUI state, reflecting the changes caused by the action. For actions that are unsupported, interactive, or lack a clear `actionCommand`, the TUI should display a user-friendly message or the direct CLI command for manual execution, rather than attempting silent advancement.
6.  **Test Strategy:**
    *   **Unit Tests:** Develop comprehensive unit tests for the `ActionRunner` implementation, verifying correct mapping of `NextAction` kinds to CLI commands and proper error handling. Test cases should cover all defined action types, including those expected to be executable and those explicitly marked as non-executable.
    *   **Integration Tests:** Write integration tests for the TUI to ensure that action selection correctly updates the UI, that triggering an action invokes the `ActionRunner` as expected, and that the TUI refreshes its state post-execution. Mock the `ActionRunner` in TUI tests to isolate UI behavior from actual command execution.
    *   **End-to-End Tests:** Consider a limited set of end-to-end tests that simulate a user interacting with the TUI, executing a "safe" action, and observing the system state change.

## Concerns / open questions
*   **Definition of "Safe" Actions:** The context mentions "supported safe actions." How precisely will "safe" be defined and enforced? Will there be an explicit whitelist, or will it implicitly be tied to the presence of an `actionCommand`? What constitutes a "safe" action for direct, non-interactive execution from the TUI without further user confirmation?
*   **Action Parameterization:** The current `actionCommand` mapping seems to imply fixed commands. How will actions requiring dynamic parameters be handled? (e.g., "answer-question with input X"). Does `runaction.NextAction` contain enough information for parameterized command construction?
*   **TUI Feedback for Long-Running Actions:** How will the TUI provide feedback for actions that might take a significant amount of time to complete? The current model implies immediate refresh.
*   **Error Reporting in TUI:** What is the UX for displaying execution errors from the `ActionRunner` to the user in the TUI?
*   **Concurrency with RefreshRuns:** If `ActionRunner` triggers `RefreshRuns`, what mechanisms will prevent race conditions or unexpected UI behavior if the user attempts further interaction before the refresh completes?
*   **Testability of CLI commands:** How can we effectively test the `ActionRunner`'s invocation of CLI commands without actually running them in a full shell environment during unit tests?

## Risks
*   **Unintended Command Execution:** Without clear delineation and safeguards for "safe" actions, there's a risk of the TUI executing commands that were meant to be interactive or require more user confirmation, potentially leading to data loss or incorrect state changes.
*   **Stale UI on Failure:** If `RefreshRuns` is not robustly integrated or fails after an action, the TUI could display stale information, misleading the user about the current state of the run.
*   **Complex `ActionRunner` Logic:** As more action types are introduced, the `ActionRunner` implementation could become overly complex, making it harder to maintain and test, especially if action parameters are involved.
*   **Testing Gaps:** Due to the interaction between TUI, `ActionRunner`, and CLI commands, there's a risk of insufficient test coverage, leading to subtle bugs in action execution or state management.
*   **Poor User Experience:** If unsupported/interactive actions are not clearly communicated, users may become frustrated by actions that appear to do nothing or lead to confusing error messages.
*   **Security Vulnerabilities:** If the `ActionRunner` is not carefully designed and secured, there could be a risk of command injection or privilege escalation if malicious `run.NextActions` are somehow introduced.