---
agent: antigravity-1
idea: build-companion-skills
review-round: 1
date: 2026-06-22
reviewed-commit: working-tree
---

## Summary

This review evaluates the implementation of the `parley-worktrees` and `parley-tracker` companion skills against the design established in [FINAL.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/addon-skills-research/FINAL.md).

Both skills conform well to the core principles of vendor neutrality, tool-agnosticism, and keeping the core protocol slim. The directory layout, separation of mirror surfaces from canonical files, and basic lifecycle commands are solid. However, critical Git command edge cases must be addressed to prevent multi-agent failures, and the tracking of the Definition of Done (DoD) requires formalization in both schemas and templates.

## Findings

### [MAJOR] Unspecified or Missing Mechanism for "Ticking" DoD Checkboxes with Commit Refs
- **What is wrong**: The skill instructs that "every dod checkbox is ticked with a commit ref" for subtasks, and that the Definition of Done is a checklist of AC IDs. However, the YAML frontmatter schema only specifies a static list of AC IDs (e.g., `dod: [AC-1, AC-2, AC-3]`), and the templates contain no markdown checkboxes or fields in the frontmatter/body to record the "ticked" state or the associated commit references.
- **Why it matters**: Without a clear convention or field to record which AC is done and with which commit reference, agents will either invent their own inconsistent formats (violating the mechanical verification goal) or fail to record the required audit trail.
- **Concrete fix**: Define how completion is tracked. Either:
  1. Add a `dod_status` field to the frontmatter that maps AC IDs to commit SHA hashes:
     ```yaml
     dod_status:
       AC-1: "sha-123456"
       AC-2: "sha-789012"
     ```
  2. Or add a markdown checklist to the body of the ticket templates where agents can tick them off and append the commit reference:
     ```markdown
     ## Definition of Done / Verification
     - [x] AC-1 (Verify: `npm test`): complete in `sha-123456`
     - [ ] AC-2 (Verify: `npm run lint`): outstanding
     ```

### [MAJOR] Concurrent Branch Checkout Failures in `git worktree add -b`
- **What is wrong**: In `addons/parley-worktrees/SKILL.md` (section 3.3), the command to setup the integration worktree is:
  ```bash
  git worktree add -b integration/<slug> \
    "${PARLEY_WORKTREES_DIR:-../worktrees}/<slug>-integration" \
    origin/main
  ```
  If the `integration/<slug>` branch has already been created (e.g., by the facilitator or another concurrent agent in the shared clone), the `-b` flag (which attempts to create a *new* branch) will fail, aborting the command.
- **Why it matters**: In a multi-agent workflow where agents share the same base repository clone (local-dir transport) or fetch remote branches (GitHub/GitLab transport), the second agent trying to check out or interact with the integration worktree will run into fatal Git errors, blocking execution.
- **Concrete fix**: Update the command sequence to attempt checking out the branch if it exists, and only fallback to creating it if it doesn't:
  ```bash
  # Check out the integration branch if it already exists locally/remotely,
  # otherwise create a new branch from origin/main:
  git worktree add "${PARLEY_WORKTREES_DIR:-../worktrees}/<slug>-integration" integration/<slug> || \
  git worktree add -b integration/<slug> "${PARLEY_WORKTREES_DIR:-../worktrees}/<slug>-integration" origin/main
  ```

### [MINOR] Potential Deletion Failures with `git branch -d` for Feature Branches
- **What is wrong**: In `addons/parley-worktrees/SKILL.md` (section 3.8), the cleanup sequence includes:
  ```bash
  git branch -d feature/<slug>/<agent-id>
  ```
  By default, Git's `-d` flag checks if the branch is fully merged into the *current checked-out branch* (which is typically the default branch like `main` in the base repository). If the integration branch `integration/<slug>` has not yet been merged into `main` (e.g., because the PR is still open or undergoing review), `git branch -d` will fail with an error stating the branch is not fully merged.
- **Why it matters**: The cleanup script will fail, leaving stale local feature branches behind and causing the agent's cleanup phase to report an error or get stuck.
- **Concrete fix**: Update the cleanup instruction to suggest using `git branch -D` if the remote branch has been pushed and verified, or run the deletion from the integration worktree context where `integration/<slug>` is checked out:
  ```bash
  # Delete feature branch after verifying it is merged in origin or using forced delete
  git branch -D feature/<slug>/<agent-id>
  ```

### [MINOR] Conceptual Overlap between Story Splitting and Subtask Breakdowns
- **What is wrong**: In `addons/parley-tracker/SKILL.md` (section "Quality and decomposition"), the skill states:
  `Split by vertical slice → disjoint file sets. A story is split into subtasks each delivering a thin end-to-end slice...`
  In Agile standards, user stories are split into smaller user stories to maintain vertical slices (end-to-end user value). Subtasks, on the other hand, represent technical implementation steps (e.g., backend implementation, frontend implementation) that are typically horizontal slices and are not independently valuable to the user.
- **Why it matters**: Saying that "subtasks each deliver a thin end-to-end slice" introduces conceptual confusion. If a task delivers an end-to-end slice, it should be a story itself (per the INVEST criteria).
- **Concrete fix**: Clarify that *stories* are split vertically into smaller stories (vertical slices), whereas *subtasks* represent the discrete technical tasks (which should have disjoint file scopes to prevent worktree collision, but do not need to be vertical end-to-end slices).

### [NIT] Missing Intersection Computation Command/Helper
- **What is wrong**: In `addons/parley-worktrees/SKILL.md` (section 6), agents are instructed to "compute the intersection of its boundary with every already-claimed boundary" and refuse/warn if the intersection is non-empty. However, no command or shell script snippet is provided to help the agent automate this computation.
- **Why it matters**: Agents must write custom code/scripts or perform manual checks, which can be error-prone when dealing with complex glob patterns.
- **Concrete fix**: Mention that the agent can compare the list of globs in the `files` array of the respective tickets or `IMPLEMENTATION.md` rows, or suggest a simple helper command/script (e.g., a python or bash snippet) that can be run to check if glob patterns intersect.

## Open questions

1. **Concurrency Locks on `IMPLEMENTATION.md`**: Under the local-dir transport, how are concurrency locks on `IMPLEMENTATION.md` handled if two agents write to it at the exact same millisecond? Does Git's standard conflict resolution suffice, or is there a helper command needed?
2. **Submodule Integration Cohesion**: For submodules, if different worktrees check out different commits or branches of a submodule, does it affect the main repository's submodule pointer? How should submodule pointer updates be coordinated when merging to the integration branch?
3. **Markdown Ticket Body Divergence**: How does `--pull` tracking sync handle changes in the description body (like updated ACs)? If changes to non-mirror-owned fields surface as conflicts, how is this conflict presented and resolved? Is there a merge tool for markdown tickets?
