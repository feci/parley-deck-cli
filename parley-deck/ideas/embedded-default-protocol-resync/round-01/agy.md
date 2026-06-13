---
agent: agy
idea: embedded-default-protocol-resync
round: 1
date: 2026-06-13
---

## Summary
The embedded cooperation protocol template (`internal/protocol/defaults/COOPERATION.md`) has drifted from the live `parley-deck/COOPERATION.md` (missing §12 and header sync metadata, and hardcoding project-specific configurations like the workspace name and agent roster). I propose to sync §12 verbatim, genericize the template header and roster to use clear placeholders, and establish an automated Go test that validates content equivalence between the template and the live deck after normalizing these project-specific zones.

## Proposed approach

### 1. §12 Propagation
- **Verbatim Carry:** Propagate §12 to the embedded default verbatim, including the closing provenance line ("ratified by idea `meta-protocol-change-end-to-end-pipeline` (2026-06-02)").
- **Rationale:** Keeping the audit history intact preserves the identity of the protocol version. It also simplifies the diff check and normalization rules in the test guard.

### 2. Genericization Scope
To make the embedded default a true template for bootstrapping new projects, we must genericize project-specific details:
- **Header Fields:**
  - `Workspace:` Replace `parley-deck` with a placeholder: `{{WORKSPACE_NAME}}`.
  - `Protocol synced:` Omit this line from the template entirely, as a fresh project is not yet synced to a specific skill snapshot.
  - `Parley deck:`, `Transport:`, `Created:`, `Status:`: Keep verbatim as protocol constants.
- **Section 2 Roster and Host Handles:**
  - Replace the parley-deck roster and host handle mapping with a standard two-agent placeholder roster:
    ```markdown
    | Agent ID       | Workspace dir                       | Role          |
    | -------------- | ----------------------------------- | ------------- |
    | `agent-1`      | `parley-deck/`                      | `participant` |
    | `agent-2`      | `parley-deck/`                      | `participant` |
    ```
    and:
    ```markdown
    | Agent ID       | Host handle    |
    | -------------- | -------------- |
    | `agent-1`      | not mapped     |
    | `agent-2`      | not mapped     |
    ```
  - **Rationale:** A placeholder roster gives new projects a clean, functional starting point without shipping parley-deck's active agent roster to unrelated codebases.

### 3. Anti-drift Mechanism
- **Automated Go Test:** Add a test in `internal/protocol/workspace_test.go` that asserts the equivalence of `internal/protocol/defaults/COOPERATION.md` and `parley-deck/COOPERATION.md`.
- **Normalization Strategy:** The test will read both files, extract/ignore the header lines and §2 tables (normalizing them to common placeholders), and perform a byte-by-byte comparison on the rest of the text.
- **Rationale:**
  - **Test vs. Generator:** A generator (`go generate`) adds build-step friction and can be forgotten. A test runs automatically during `go test ./...` and blocks CI/builds on drift.
  - **Test vs. CLI Command:** A subcommand like `parley protocol check` is user-facing, whereas drift detection is a internal developer maintenance concern.

## Concerns / open questions
- **Init-time Replacement:** Should `InitWorkspace` perform basic replacements beyond transport? For instance, replacing `{{WORKSPACE_NAME}}` with the actual directory name of the target repository so the bootstrapped file is immediately ready for use without manual editing? If so, this requires updating `internal/protocol/workspace.go`. I recommend keeping it simple: just transport replacement for now, leaving `{{WORKSPACE_NAME}}` as a clear manual checklist item for the human integrator (as per Appendix A step 3: "Fill in the header").

## Risks
- **Test Fragility:** Future changes to the formatting of §2 or the header block might break the regexes/parsers used in the test's normalization step. The normalization parser must be kept simple and robust (e.g. matching Markdown tables or lines by prefix).
- **Template Completeness:** If we replace the roster, the bootstrapping agent needs to know they must manually configure it. Appendix A already documents this, so the risk is low.
