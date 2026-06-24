---
agent: antigravity-1
idea: close-integrity
review-round: 1
date: 2026-06-24
---

## Summary

In Phase 6 refutation mode, I reviewed the `close-integrity` design ([FINAL.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/close-integrity/FINAL.md) and [IMPLEMENTATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/close-integrity/IMPLEMENTATION.md)) and the implementation changes since `97034cd`. My focus was on security, robustness, and the possibility of new failure modes introduced by the close decision gates (LE-7 goal-done check and LE-11 HITL-fatigue guardrails).

I identified **four new findings** that were missed by the previous reviewers:
1. **CRITICAL**: Duplicate reviewer IDs lead to concurrent write race conditions, log corruption, and OS lockups.
2. **MAJOR**: The goal-check consult lacks a short timeout, causing it to block for 15+ minutes on hangs rather than failing open quickly.
3. **MAJOR**: Markdown bolding of the label prefix (`**GOAL-CHECK:** FAIL`) causes the parser to fail-open, silently swallowing a confident FAIL.
4. **MINOR**: The parser does not reset/clear prior verdicts when encountering subsequent ambiguous verdict lines, violating the "last verdict wins" contract.

**LE-11a Deferral Assessment:** The deferral of batching driver-opened HITL questions is **acceptable**. Since the driver halts execution on any escalation, it cannot flood the user with multiple concurrent questions for a single idea. Batching across ideas is a multi-idea coordination concern that properly belongs in the Tier-4 outer loop (`standing-loop-watch-mode`), not within the single-idea `Driver`.

---

## Refutation attempts

I systematically tried to break the implementation gates:

1. **Bypassing the `< 2` Reviewer Guard with Duplicate IDs:**
   - *Attempt:* I analyzed what happens if `participants` in the run contains duplicate reviewer IDs (e.g. `[impl, rev, rev]`). 
   - *Result:* `newDriverImplOps` puts duplicate entries in `o.reviewers`, causing `ReviewerCount` to be reported as `2`. This bypasses the gate. Furthermore, during reviewer execution, `runner.RunReviewRound` launches concurrent goroutines to execute the same agent ID (`rev`), which write to the exact same log files (`stdout.log`/`stderr.log`) and output artifact (`rev.md`) simultaneously, causing a severe concurrent-write race condition.

2. **Hanging/Unbounded Goal-Check Consults:**
   - *Attempt:* I traced the timeout configuration in `driverImplOps.GoalCheck`. 
   - *Result:* The call to `runner.RunConsult` leaves `opts.Timeout` unset. It defaults to the agent's full execution timeout (usually 15–30 minutes). If a model hangs or loops during the goal check, the auto-drive tick is blocked for this long instead of failing open quickly, risking hard context deadline failure for the driver execution.

3. **Bypassing `parseGoalVerdict` with Bolding (`**GOAL-CHECK:** FAIL`):**
   - *Attempt:* I evaluated how the parser handles common markdown bolding variations.
   - *Result:* If a model outputs `**GOAL-CHECK:** FAIL`, the leading asterisks are to be stripped by `TrimLeft`, resulting in `t = GOAL-CHECK:** FAIL`. The prefix `GOAL-CHECK:` matches, but `rest` becomes `** FAIL`. Since `** FAIL` does not start with `FAIL` (due to the trailing asterisks), the switch block skips it, and it fail-opens.

4. **Tricking the Parser's "Last Verdict Wins" Rule:**
   - *Attempt:* I analyzed what happens when the agent outputs multiple verdict lines, where the last one is invalid/ambiguous (e.g., `GOAL-CHECK: FAIL` followed by `GOAL-CHECK: RE-EVALUATING`).
   - *Result:* The parser matches the second line but does not match `PASS` or `FAIL` in the switch, leaving the previous `verdict` as `FAIL`. The function returns `FAIL` rather than resetting to ambiguous (`""`), violating the contract that the last verdict line determines the result.

5. **Running Goal-Check with the Implementer:**
   - *Attempt:* Checked if the implementer can act as checker when no reviewers are present.
   - *Result:* `newDriverImplOps` falls back to `drafter = implementer` when there are no reviewers. In an auto-drive run, this is blocked by the `ReviewerCount < 2` guard. In a strict-gate design-only run with no reviewers, `OpenReviewRound` fails earlier. Thus, while `GoalCheck` lacks local protection, upstream gates prevent the implementer from running as checker.

---

## Findings

### [CRITICAL] Concurrent execution of duplicate reviewers causes log corruption and filesystem race conditions

- **What is wrong:** If duplicate reviewer IDs are specified in the participants list (e.g., `[impl, rev, rev]`), `newDriverImplOps` builds `o.reviewers` as `[rev, rev]`, and `ReviewerCount` returns `2`. When `OpenReviewRound` is called, `runner.RunReviewRound` starts a goroutine for each element in the list. Two concurrent goroutines call `runAgent` for the same agent ID `rev`. Both goroutines concurrently write to:
  - `runs/<run-id>/agents/rev/stdout.log`
  - `runs/<run-id>/agents/rev/stderr.log`
  - `review/round-XX/rev.md`
- **Why it matters:** Operating systems and filesystems do not support concurrent, uncoordinated writes to the same files. This race condition leads to corrupted/interleaved logs, truncated artifacts, and intermittent write errors.
- **Concrete fix:** Normalize `o.reviewers` to a unique set in `newDriverImplOps` and reject duplicate participant lists at the API/CLI boundary:
  ```diff
  // internal/app/driver_impl.go
  func newDriverImplOps(base runner.Options, root, ideaSlug, ideaDir string, participants []string, out io.Writer) driver.ImplOps {
  	implementer := resolveImplementer(ideaDir, participants)
  	var reviewers []string
+ 	seen := make(map[string]bool)
  	for _, p := range participants {
- 		if p != implementer {
+ 		if p != implementer && !seen[p] {
+ 			seen[p] = true
  			reviewers = append(reviewers, p)
  		}
  	}
  ```

### [MAJOR] GoalCheck consult has no timeout and can run for 15+ minutes

- **What is wrong:** In `driverImplOps.GoalCheck` ([internal/app/driver_impl.go](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/app/driver_impl.go#L336-L367)), `runner.RunConsult` is called without specifying `Timeout` in `runner.ConsultOptions`. This causes the consult to fall back to the agent's default timeout (`agent.TimeoutMS` or `agents.DefaultTimeoutMS`), which is typically 15–30 minutes.
- **Why it matters:** The goal-check is a defense-in-depth, advisory gate designed to fail-open on timeout or failure. If the consult hangs or the model goes into a slow generation loop, it will block the driver for the entire 15+ minute agent duration. This wastes API cost and time, and can consume the entire parent context limit, causing the driver run to abort with a hard error rather than failing open gracefully.
- **Concrete fix:** Provide a short timeout (e.g. 2 minutes) for the goal-check consult inside `driverImplOps.GoalCheck`:
  ```diff
  // internal/app/driver_impl.go
  	res := runner.RunConsult(ctx, runner.ConsultOptions{
  		Root:       o.root,
  		Agent:      agent,
  		Prompt:     runner.BuildGoalCheckPrompt(agent, o.base.Idea),
+ 		Timeout:    2 * time.Minute,
  		StdoutPath: filepath.Join(dir, "goal-check.stdout.log"),
  		StderrPath: filepath.Join(dir, "goal-check.stderr.log"),
  		Progress:   o.out,
  	})
  ```

### [MAJOR] Markdown bolding wrapper (`**GOAL-CHECK:**`) bypasses the parser and fail-opens

- **What is wrong:** If the agent outputs a verdict bolding the marker itself (a common formatting practice), e.g., `**GOAL-CHECK:** FAIL`, the parser strips the leading `**` via `TrimLeft`, leaving `t = GOAL-CHECK:** FAIL`. This matches the prefix `GOAL-CHECK:`, but `rest` becomes `** FAIL`. The `switch` checks if `rest` has prefix `FAIL` or `PASS`. Because `** FAIL` starts with `*`, it is treated as ambiguous and fail-opens.
- **Why it matters:** Standard LLM behavior often bolds markdown prefixes. If a model outputs a confident `FAIL` formatted this way, it is silently swallowed by the fail-open check, defeating the goal gate.
- **Concrete fix:** Strip leading and trailing markdown markers (like asterisks) from `rest` inside `parseGoalVerdict` before evaluating the prefix:
  ```diff
  // internal/app/driver_impl.go
  func parseGoalVerdict(answer string) string {
  	verdict := ""
  	for _, line := range strings.Split(answer, "\n") {
  		t := strings.ToUpper(strings.TrimSpace(line))
  		t = strings.TrimLeft(t, "#*-> \t`\"'")
  		if !strings.HasPrefix(t, "GOAL-CHECK:") {
  			continue
  		}
  		rest := strings.TrimSpace(strings.TrimPrefix(t, "GOAL-CHECK:"))
+ 		rest = strings.TrimLeft(rest, "*`\"'")
  		switch {
  		case strings.HasPrefix(rest, "PASS"):
  			verdict = "PASS"
  		case strings.HasPrefix(rest, "FAIL"):
  			verdict = "FAIL"
  		}
  	}
  	return verdict
  }
  ```

### [MINOR] parseGoalVerdict fails to reset to ambiguous on subsequent invalid verdict lines

- **What is wrong:** If an agent outputs multiple lines starting with `GOAL-CHECK:`, and the last one does not start with `PASS` or `FAIL` (e.g. `GOAL-CHECK: FAIL` followed by `GOAL-CHECK: INCONCLUSIVE`), the parser matches `GOAL-CHECK:` on the second line but falls through the switch. Since `verdict` is not reset/cleared, the function returns the prior `FAIL`.
- **Why it matters:** This violates the "last verdict wins" rule. A disclaimer or correction line written by the agent will be ignored, leading to false-positives (unintended escalations) instead of failing open.
- **Concrete fix:** Reset the `verdict` value to `""` inside the loop whenever a new `GOAL-CHECK:` prefix is matched, before matching `PASS` or `FAIL`:
  ```diff
  // internal/app/driver_impl.go
  		if !strings.HasPrefix(t, "GOAL-CHECK:") {
  			continue
  		}
+ 		verdict = "" // reset on each verdict line so the latest one wins
  		rest := strings.TrimSpace(strings.TrimPrefix(t, "GOAL-CHECK:"))
  ```

---

## Open questions

1. **Local protection in GoalCheck:** Should `GoalCheck` have local checks asserting `o.drafter != o.implementer` and fail-open with a warning, rather than solely relying on the upstream `ReviewerCount < 2` and `OpenReviewRound` check? (Adding local robustness is recommended).
2. **Defensive checklist validation:** Should the workspace parser reject duplicate participant IDs during `parley run` configuration loading to prevent other potential race conditions?
