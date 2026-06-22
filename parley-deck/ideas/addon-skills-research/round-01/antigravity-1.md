---
agent: antigravity-1
idea: addon-skills-research
round: 1
date: 2026-06-22
---

## Summary

This contribution approaches the two proposed companion skills—`worktree-sessions` and `ticket-tracker`—through the lens of established standards, industry prior-art, and conventions. For **Track A (Worktrees)**, we examine how parallel-agent orchestration platforms isolate execution states using Git worktrees, converting soft concurrent errors into resolvable merge conflicts. For **Track B (Ticketing)**, we apply the **INVEST** framework, **Gherkin / Given-When-Then** acceptance criteria, rigorous **Definitions of Done (DoD)**, and **story-splitting** techniques to make requirements readable by business stakeholders, technical developers, and automated AI agents. We map these standards to concrete mechanics, workflows, and templates for the Parley Deck ecosystem.

---

## Proposed approach

### Track A — Worktrees (Parallel-Agent Workflows & Conventions)

#### 1. Standards and Prior-Art
Industry prior-art in late-2025 and 2026 (e.g., Cursor, Claude Code, Augment, Parallel Code, Upsun) leverages Git worktrees to manage multi-agent concurrent execution on a single repository. The standard conventions established by these platforms include:
*   **Path Isolation:** Spawning individual agents into distinct workspace paths (e.g., `/worktrees/<session-id>/<agent-id>`) sharing a single parent `.git` directory.
*   **Concurrency Resolution:** Rather than risking concurrent write collisions, file clobbering, or database lock errors when multiple agents run tests or write code in the same working tree, worktrees ensure that overlap is isolated to individual branches. Concurrency issues are thus deferred to standard git branch merges and standard merge conflicts.
*   **Environment Isolation:** Making all local environment configurations (`.env.local`, `.npmrc`) untracked and localized strictly to the worktree path.
*   **Sparse Checkouts:** Constraining the directory scope visible to the agent inside its worktree (via `git sparse-checkout`) to decrease token overhead during file scans and search tool invocations.
*   **Staging/Integration Branches:** Requiring parallel branches to merge into an integration/staging branch (`staging/<session>`) where test suites are validated prior to merging into `main`.

#### 2. Command Sequences for Parley Deck Case
The `worktree-sessions` companion skill will automate the following lifecycle commands:

##### Step 1: Provisioning a Worktree
When a Parley session or a Phase-5 implementer claims a task in `ideas/<slug>/IMPLEMENTATION.md`, the CLI skill provisions a dedicated worktree:
```bash
# Define paths and branch naming conventions
WORKTREE_PATH="./.git/worktrees/<slug>-<agent-id>"
BRANCH_NAME="parley/impl/<slug>-<agent-id>"

# Add a new worktree checked out to a new branch tracking origin/main
git worktree add -b "$BRANCH_NAME" "$WORKTREE_PATH" origin/main
```

##### Step 2: Environment Setup and Sparse-Checkout
To avoid token bloat and environment collisions, initialize a sparse checkout and copy required local environment profiles:
```bash
# Navigate to the worktree directory
cd "$WORKTREE_PATH"

# Initialize sparse checkout in cone mode (limits files to root and specified dirs)
git sparse-checkout init --cone

# Define the paths relevant to the agent's task domain (read from task hints)
git sparse-checkout set internal/pkg/config cmd/helper/

# Copy the local environment profile (which is in .gitignore in the parent repo)
cp ../../.env.local .env.local
```

##### Step 3: Lifecycle Management and Listing
To check the status of running sessions and their paths:
```bash
git worktree list
```

##### Step 4: Verification and Teardown
Once the implementation is complete, the branch is pushed to origin, tested via CI, and the worktree is safely deleted:
```bash
# Navigate back to parent repository root
cd ../../

# Remove the worktree directory and associated metadata
git worktree remove --force "$WORKTREE_PATH"

# Prune any stale worktree administrative files
git worktree prune
```

#### 3. Pitfalls and Mitigations
*   **Branch Locks:** Git prevents checking out the same branch in multiple worktrees. *Mitigation:* The skill must enforce a unique branch name per agent-session (`parley/impl/<slug>-<agent-id>`).
*   **Untracked / Ignored Files:** Configuration files (like `.env`) do not carry over to the new worktree. *Mitigation:* The provisioning hook must explicitly detect and copy defined local configuration templates to the worktree path.
*   **Stale Metadata:** If a worktree directory is deleted using system commands (e.g. `rm -rf`), Git's internal state remains locked. *Mitigation:* The skill must run `git worktree prune` periodically at startup and teardown.
*   **IDE Confusion:** IDEs can lose track of references or index multiple copies of the same files. *Mitigation:* The skill should write IDE exclude settings (e.g., adding the `./.git/worktrees/` directory to VS Code's `files.exclude` or `.gitignore` in the parent) to prevent dual indexing.

#### 4. The Parley Deck Seam
*   **Skill Boundary (`worktree-sessions`):** Handles the low-level Git mechanics (adding, configuring sparse checkouts, copying env configs, removing, and pruning).
*   **Protocol Seam (in `COOPERATION.md`):** A thin amendment under **§6 Conflict-avoidance mechanics** stating: *"Implementers executing in parallel MAY be isolated in local Git worktrees. When worktree isolation is active, task coordination MUST use the shared `IMPLEMENTATION.md` file as the lock manifest, and changes MUST merge through a shared integration branch before landing on main."*

---

### Track B — Ticketing (INVEST, Gherkin, DoD, & Story-Splitting)

#### 1. Standards and Prior-Art
*   **INVEST Framework:** User stories must be **I**ndependent (allowing parallel worktrees to operate without blocks), **N**egotiable (allowing the agent to refine implementation detail in Phase 5's `IMPLEMENTATION.md`), **V**aluable (focusing on end-user/system utility), **E**stimable (fully scoped), **S**mall (fitting inside an agent's context window), and **T**estable (defined by clear assertions).
*   **Gherkin / Given-When-Then:** Establishes structured, machine-parseable acceptance criteria. This removes ambiguity and provides a direct blueprint for AI-generated tests.
*   **Definition of Done (DoD):** A standardized, cross-audience quality checklist ensuring that files remain the source of truth, tests pass, and peer reviews are completed.
*   **Story-Splitting:** Breaking down epics by vertical slice (e.g., by workflow step, variations in data, or interface channels) to permit parallel worktree execution without overlapping file modifications.

#### 2. The "Mirror" Rule
The ticket tracker (Jira, Linear, GitHub Issues) is a **mirror**, not the source of truth. The canonical specifications live as Markdown files under `parley-deck/ideas/<slug>/tickets/`. The `ticket-tracker` companion skill pushes updates from these Markdown files to the tracker APIs, ensuring that local files remain the primary authority.

#### 3. Multi-Audience Templates

##### Template 1: Epic (Epic Ticket Markdown)
```markdown
# Epic: [Title of the Epic]

## Business Value (Why)
* **Goal:** [High-level business goal]
* **Target Audience:** [Business stakeholders, users]
* **User Value:** [What value does this bring to the business/users?]

## Technical Architecture & Constraints (How)
* **System Boundaries:** [Affected services / packages]
* **Architectural Touchpoints:** [e.g., OpenAPI endpoints, DB schemas]
* **Non-Functional Requirements (NFRs):** [Performance, scalability, security limits]

## Epic Definition of Done (DoD)
- [ ] Core system integration tests passing.
- [ ] OpenAPI documentation synchronized and verified.
- [ ] Security audit and lint checking completed.
- [ ] Local markdown tickets mirror the ticket tracker state.
```

##### Template 2: User Story (User Story Ticket Markdown)
```markdown
# Story: [Title of the User Story]
**Epic Link:** [Link or ID of Epic]
**Status:** Backlog | **Priority:** High/Medium/Low
**Labels:** backend, api, security

## Value Statement
* **As a** [user role]
* **I want** [to perform some action]
* **So that** [some value is achieved]

## Technical Specifications
* **Schema Modifications:** [DB columns, API requests/responses]
* **Affected Packages/Paths:** [Path hints, e.g. `internal/auth/`, `pkg/token/`]
* **Dependencies:** [List of dependent stories or subtasks]

## Acceptance Criteria (Gherkin Format)
```gherkin
Scenario: Successful action execution
  Given [precondition state]
  And [additional context]
  When [user triggers action]
  Then [expected state modification occurs]
  And [response payload matches schema]

Scenario: Error handling under invalid state
  Given [invalid precondition state]
  When [user triggers action]
  Then [system returns error code X]
  And [no state modification occurs]
```

## Definition of Done
- [ ] Gherkin acceptance criteria mapped to automated unit/integration tests.
- [ ] Code builds without warnings and passes all linter checks.
- [ ] Code reviewed and approved by at least one peer agent (Phase 6).
- [ ] Zero secret values or uncommitted configurations leaked.
```

##### Template 3: Technical Subtask (Subtask Ticket Markdown)
```markdown
# Subtask: [Title of the Subtask]
**Parent Story:** [Link or ID of User Story]
**Status:** Todo | **Priority:** Medium

## Technical Instructions
* **Target Files:** [Specific files to modify, e.g., `internal/auth/provider.go`]
* **Implementation Logic:**
  1. [Step-by-step Technical execution instruction]
  2. [Specific function signature or interface implementation]
* **Test Case Requirement:** [Unit test target and coverage goal]
```

---

## Concerns / open questions

1.  **Shared Dependency Lock Files:** When parallel worktrees modify dependencies (e.g., adding packages to `go.mod`/`go.sum` or `package.json`), merging these branches back to `main` via staging will cause lock file conflicts. How should `worktree-sessions` handle auto-resolving these structured merge conflicts?
2.  **External Ticket Drift:** If a human business stakeholder changes a story priority or description directly inside the Jira/Linear web UI, how does the `ticket-tracker` skill reconcile this change back to the local canonical Markdown files without overwriting local-first edits?
3.  **IDE Context and Path Mapping:** Modern IDE search indexes can get confused by multiple copies of the code inside the hidden `.git/worktrees/` directory. Can we guarantee that all major editors (VS Code, JetBrains) respect the path exclusions automatically generated by our skill?

---

## Risks

1.  **Out-of-Sync Split-Brain:** A breakdown in the mirroring sync loop between the local canonical markdown files and the remote tracker could lead to business stakeholders reading outdated details while agents implement the true canonical local code.
2.  **Stale Worktree Disk Bloat:** If agents abort sessions mid-way or fail to clean up their worktrees upon task completion, the local machine will accumulate duplicate dependencies (e.g., massive `node_modules` folders) and stale git administrative data, exhausting disk space.
3.  **Branch Checkout Blockers:** If an agent script crashes or exits abnormally without releasing a branch from a worktree, subsequent command runs will fail with `fatal: 'branch-name' is already checked out...`, requiring manual human intervention to clean up.
4.  **Incorrect Gherkin Implementations:** AI agents might write Gherkin test assertions that pass syntactically but do not map to actual code hooks, creating false positives in test reports.
