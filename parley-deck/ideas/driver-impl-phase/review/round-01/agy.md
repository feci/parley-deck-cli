---
agent: agy
idea: driver-impl-phase
review-round: 1
date: 2026-06-06
reviewed-commit: f624c05
---

## Summary

This is the Round 1 code review of the Phase 5-8 auto-driver implementation (`driver-impl-phase`) by **agy**.

The implementation succeeds in auto-driving the implementation, review, consensus, and fix-up phases, successfully delegating tasks to agents via the `ImplOps` seam while keeping the core driver decoupled from `internal/app`. However, several critical safety-gate omissions, crash-recovery stale-state loops, lack of role separation, and brittle frontmatter parsers have been identified during this review. 

## Findings

### [CRITICAL] RunChecks Build/Test Gate Omitted in Phase 8 Fix-Up Path
* **What's wrong:** The `RunChecks` build/test gate is executed only in `advanceImpl` (after the initial implementation is ready for review). It is completely omitted in the `advanceReview` fix-up branch after a code-writing agent completes `Fixup`.
* **Why it matters:** If the implementer's fix-up fails compilation or breaks tests, the driver will proceed to invoke reviewer agents anyway. This wastes expensive LLM and API reviewer credits on broken code, directly violating the safety guarantees specified in `FINAL.md` ("RunChecks ... runs after implement + each fixup before review").
* **Concrete fix:** In [internal/driver/impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go#L181-L191), call `RunChecks` immediately after the `Fixup` execution and escalate on failure before archiving consensus and opening the next review round:
  ```go
  if err := d.cfg.Impl.Fixup(ctx, cycle); err != nil {
  	return ActionEscalated, c, fmt.Errorf("fix-up cycle %d: %w", cycle, err)
  }
  if ok, detail := d.cfg.Impl.RunChecks(ctx); !ok {
  	return ActionEscalated, c, fmt.Errorf("checks failed after fix-up cycle %d before review:\n%s", cycle, strings.TrimSpace(detail))
  }
  ```

### [MAJOR] Non-Idempotent Crash-Stranding / Stale-State Hazard in Fix-Up Loop
* **What's wrong:** If the driver process crashes or gets interrupted *after* `Fixup` executes but *before* the active `review/consensus.md` is archived and the next review round is opened, the driver will get stuck on re-entry. On the next tick, `Rebuild` will still derive `PhaseReview` at the same `highestReviewRound` (since the next round directory was not created), find the root `review/consensus.md` with outstanding fixes, and invoke `Fixup` again, resulting in duplicate code-writing agent calls.
* **Why it matters:** Non-idempotent re-entry wastes agent credits, executes code-writing multiple times, and pollutes git history.
* **Concrete fix:** Make `advanceReview` recovery-aware of the `IMPLEMENTATION.md` status cycle value. Before executing `Fixup`, parse `IMPLEMENTATION.md`'s status. If the status is already `"fix-up-cycle-N"` (where `N` matches the current `round`), we know the fix-up has already completed. Skip calling `Fixup` and directly execute the archiving and next round setup:
  ```go
  status, err := d.cfg.Impl.ImplementationStatus()
  if err != nil {
  	return ActionEscalated, c, fmt.Errorf("read implementation status: %w", err)
  }
  if status == fmt.Sprintf("fix-up-cycle-%d", round) {
  	// RunChecks should still be verified in the recovery path
  	if ok, detail := d.cfg.Impl.RunChecks(ctx); !ok {
  		return ActionEscalated, c, fmt.Errorf("checks failed after fix-up cycle %d before review:\n%s", round, strings.TrimSpace(detail))
  	}
  	if err := d.archiveReviewConsensus(round); err != nil {
  		return ActionEscalated, c, fmt.Errorf("archive review consensus: %w", err)
  	}
  	if err := d.cfg.Impl.OpenReviewRound(ctx, round+1); err != nil {
  		return ActionEscalated, c, fmt.Errorf("open review round %d: %w", round+1, err)
  	}
  	c.Phase = PhaseReview
  	c.UpdatedAt = nowRFC3339()
  	_ = c.Save(d.cursorPath())
  	return ActionFixup, c, nil
  }
  ```

### [MAJOR] Lack of Implementer/Reviewer Role Separation for Consensus Drafting
* **What's wrong:** In [internal/app/driver_impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go#L48), the review-consensus `drafter` is set directly to the `implementer`.
* **Why it matters:** Having the implementer draft `review/consensus.md` directly violates the separation of roles between implementation and review. The implementer could act biased and filter out critical/major reviewer findings from the final consensus list.
* **Concrete fix:** Set `drafter` to a non-implementer participant (e.g., the first agent in `reviewers` or a designated facilitator among reviewers) instead of `implementer`:
  ```go
  drafter := ""
  if len(reviewers) > 0 {
  	drafter = reviewers[0]
  }
  ```

### [MAJOR] Brittle Parsing of `outstanding_agreed_fixes` and `blocked` Fields in ReviewStatus
* **What's wrong:** In [internal/app/driver_impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go#L147-L152)'s `ReviewStatus()`, the frontmatter values read by `ReadFrontmatter` for `outstanding_agreed_fixes` and `blocked` are not stripped of surrounding quotes (`"` or `'`).
* **Why it matters:** If an LLM or YAML writer formats the frontmatter with quotes (e.g., `outstanding_agreed_fixes: "0"` or `blocked: "false"`), `strconv.Atoi` will fail with an error and the driver will escalate/crash.
* **Concrete fix:** Strip single and double quotes from `outstanding_agreed_fixes` and `blocked` before parsing, similar to `ImplementationStatus()`:
  ```go
  rawFixes := strings.Trim(strings.TrimSpace(meta["outstanding_agreed_fixes"]), `"'`)
  fixes, err := strconv.Atoi(rawFixes)
  ...
  blockedStr := strings.Trim(strings.TrimSpace(meta["blocked"]), `"'`)
  blocked := strings.EqualFold(blockedStr, "true")
  ```

### [MAJOR] Malformed/Invalid Reviewer Files Cause Infinite Retry Loop
* **What's wrong:** If a reviewer writes an invalid review file (e.g. missing YAML frontmatter), `ReviewRoundComplete` returns `false, nil`. In `advanceReview`, since `complete` is false, it calls `OpenReviewRound` again. However, because the reviewer file already exists on disk, the runner's `runAgent` skips it (when `Overwrite` is false), meaning it is never re-generated or corrected, causing the driver to spin endlessly (or until timeout) in `ActionAwait`.
* **Why it matters:** The auto-driver is stranded in an infinite loop without correcting the validation error or escalating/re-drafting correctly.
* **Concrete fix:** Before calling `OpenReviewRound` inside `advanceReview`'s retry path, remove any malformed review artifacts for the active round that failed validation so that `runAgent` is forced to re-execute the agent to correct the file.

### [MINOR] Brittle Opt-In Parsing for `auto_implement` and `cross_review_rounds`
* **What's wrong:** `ReadAutoImplement` and `ReadCrossReviewRounds` in [internal/driver/transport.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/transport.go#L45-L50) parse values directly from `00-prompt.md` without stripping surrounding quotes.
* **Why it matters:** If a user or agent creates `00-prompt.md` with `auto_implement: "true"`, the parsing fails to recognize the opt-in and disables code-writing.
* **Concrete fix:** Strip single/double quotes before checking `strings.EqualFold` or running `strconv.Atoi`:
  ```go
  clean := strings.Trim(strings.TrimSpace(v), `"'`)
  ```

### [MINOR] Git Status Command Failures Treated as Clean
* **What's wrong:** In [internal/driver/impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/driver/impl.go#L231-L238)'s `gitTreeClean(root)`, if the `git status --porcelain` command fails (e.g., due to a lock file like `.git/index.lock` or command timeout/permission errors), `err != nil` triggers and the function returns `true` (clean).
* **Why it matters:** Treating actual errors during dirty-checks as "clean" poses a safety risk where uncommitted human changes could be overwritten.
* **Concrete fix:** Check if the repository contains a `.git` directory first (or run `git rev-parse --is-inside-work-tree`). If it is a git repo, treat command failures as unsafe/dirty:
  ```go
  func gitTreeClean(root string) bool {
  	cmdCheck := exec.Command("git", "-C", root, "rev-parse", "--is-inside-work-tree")
  	if err := cmdCheck.Run(); err != nil {
  		// not inside a git tree or git command completely missing -> treat as clean
  		return true
  	}
  	cmd := exec.Command("git", "-C", root, "status", "--porcelain")
  	out, err := cmd.Output()
  	if err != nil {
  		return false // git error inside repo -> treat as dirty/unsafe
  	}
  	return strings.TrimSpace(string(out)) == ""
  }
  ```

## Open questions

1. **Reviewer Retry/Timeout Policies:** If a reviewer agent consistently fails to write a valid file, should the driver continue to retry up to the loop deadline, or should there be a distinct retry limit count to escalate earlier?
2. **Extensible Verification Command:** Currently, `RunChecks` hardcodes `go test ./...` if a `go.mod` is found. In workspaces that do not use Go (or use mixed languages), how should custom check/verify commands be declared (e.g., in `00-prompt.md` under a `verification_command` key)?
