---
agent: antigravity-1
idea: verification-honesty
review-round: 1
date: 2026-06-24
---

## Summary

The implementation in `loop-engineering-impl` successfully introduces the core mechanisms for Tier 1 verification honesty (LE-1 through LE-4) and passes all unit tests. However, adopting a refutation posture exposes several critical gaps in the deterministic scan veto, model diversity guard, and command execution safety. These vulnerabilities allow the strict closing loop to complete without real verification in several scenarios and permit arbitrary command execution on the host system without validation.

## Refutation attempts

I attempted to break the implementation by constructing failing paths and tracing code execution for each of the core requirements in `FINAL.md`:

1. **Strict Closing Round Loop Unbounded Spin:**
   * *Attempt:* Checked whether a non-clean strict gate consensus can loop indefinitely.
   * *Result:* Failed to break. The driver checks `round >= d.cfg.MaxFixupCycles` inside both the strict closing check path and the outstanding agreed fixes path. Because `OpenReviewRound` is called with `round + 1` which results in the creation of `review/round-NN` directories (where `NN = round + 1`), the highest round number is guaranteed to increment. If the loop does not stabilize on a clean round before reaching `MaxFixupCycles` (defaulting to 3), the driver escalates and halts.
2. **Bypassing the Strict-Gate Veto (Finding Scan):**
   * *Attempt:* Tried to write a review with valid findings formatted in a way that escapes the deterministic veto scan.
   * *Result:* **Succeeded in breaking.** If a reviewer outputs a finding heading with lowercase or mixed-case severity (e.g. `### [critical]`) or writes a heading with the description on the next line (empty title on the header line, e.g. `### [CRITICAL]\nBug...`), the validator accepts it, but the scan misses it. The driver then auto-completes the round, letting a dirty implementation pass.
3. **Model Diversity Guard Failure to Warn/Escalate:**
   * *Attempt:* Checked if the model diversity check can fail open.
   * *Result:* **Succeeded in breaking.** If a reviewer's model returns `""` (e.g. due to discovery issues) or has casing differences compared to the implementer's model configuration, the comparison fails open.
4. **Bypassing RunChecks Fail-Closed for Code Ideas:**
   * *Attempt:* Tried to bypass the `auto_implement` checks gate without writing tests or Go code.
   * *Result:* **Succeeded in breaking.** Setting `checks: "true"` in the frontmatter bypasses the check gate completely.
5. **RunChecks Unsafe Command Execution:**
   * *Attempt:* Checked if the `checks:` command can execute arbitrary shell scripts.
   * *Result:* **Succeeded in breaking.** The command is passed directly to `sh -c` without sandboxing, posing a severe remote code execution risk if deliberation prompts are untrusted.

## Findings

### [MAJOR] Deterministic Scan Veto Bypassed by Case Casing or Extra Spacing in Finding Headings
* **What is wrong:** In `internal/driver/impl.go`, `scanHasRealFinding` does a case-sensitive check on the severity tag:
  ```go
  switch t[len("### ["):close] {
  case "CRITICAL", "MAJOR", "MINOR", "NIT":
  ```
  Additionally, the prefix matching `strings.HasPrefix(t, "### [")` is sensitive to extra spaces after the hashes (e.g., `###   [CRITICAL]`).
* **Why it matters:** LLM reviewers frequently output lowercase or mixed-case severity tags (e.g., `[critical]`, `[Major]`) or vary spacing in markdown headers. Since `ValidateReviewArtifact` does not enforce strict uppercase or heading formats (only checking that `"## Findings"` is present), these files will validate successfully. However, `scanHasRealFinding` will ignore them, causing the deterministic veto to fail open and allow a dirty round to auto-complete.
* **Concrete fix:** Normalize the severity tag to uppercase before comparison using `strings.ToUpper` or `strings.EqualFold`. Use a more robust regex or whitespace-tolerant parser to extract the severity tag and title from the heading.

### [MAJOR] Deterministic Scan Veto Bypassed by Empty Header Titles
* **What is wrong:** `scanHasRealFinding` explicitly ignores any heading where the title is empty:
  ```go
  title := strings.TrimSpace(t[close+1:])
  if title != "" && title != "<title>" {
      return true
  }
  ```
* **Why it matters:** LLM reviewers often format findings with the severity tag as the header line and the title or description on the next line (e.g., `### [CRITICAL]\nMemory leak in worker`). In this case, `title` is `""`, and the finding is skipped by the scan, bypassing the strict gate veto.
* **Concrete fix:** If a line matches `### [SEVERITY]` (where severity is one of the valid tags), treat it as a real finding regardless of whether there is text following the bracket on the same line.

### [MAJOR] Model Diversity Guard Event `agent.model_diversity` is Missing
* **What is wrong:** `FINAL.md` LE-3 states: "If all reviewers share the implementer's model, emit a warning (stdout + an `agent.model_diversity` event)." However, the event emission is completely missing from `internal/app/driver_impl.go`.
* **Why it matters:** Automated systems or TUI interfaces that listen to parley-deck events to monitor state will not receive notice that model diversity has been violated.
* **Concrete fix:** In `internal/app/driver_impl.go`, inside `OpenReviewRound`, emit the event using `o.base.Store.Append(store.Event{Timestamp: nowRFC3339(), Type: "agent.model_diversity", Data: map[string]any{"model": model}})` when a same-model deck is detected.

### [MINOR] Model-Diversity Guard Fails Open on Case Casing and Unknown Models
* **What is wrong:** In `internal/app/driver_impl.go`, the comparison `o.modelOf(r) != implModel` is case-sensitive. Additionally, if a reviewer's model is unresolved (returns `""`), the guard assumes models differ.
* **Why it matters:** If models are configured with different casing (e.g., `gpt-4o` vs `GPT-4o`) or discovery fails to resolve a model, the guard will fail open and skip warning or escalating, even if the agents are using the same underlying model.
* **Concrete fix:** Perform a case-insensitive comparison using `strings.EqualFold` and explicitly handle or warn if a reviewer's model is unresolved (`""`).

### [MINOR] Unsafe Shell Execution in RunChecks Command
* **What is wrong:** `RunChecks` executes arbitrary commands in the `checks` field of `00-prompt.md` via `sh -c` on the host environment without sanitization or sandboxing.
* **Why it matters:** If an idea's prompt is modified by an untrusted or compromised agent during deliberation, it could execute destructive or malicious shell commands (e.g., `checks: "rm -rf /"`).
* **Concrete fix:** Restrict `checks` to a safe whitelist of commands, or require explicit user confirmation before running a custom `checks:` shell script, or run it in a containerized/sandboxed environment.

### [NIT] Missing Validation of Strict Gate Fields in ValidateReviewConsensusArtifact
* **What is wrong:** Under `strict_gate: true`, the draft consensus must include `strict_gate_clean` and `closing_review_round`. However, `ValidateReviewConsensusArtifact` does not enforce their presence.
* **Why it matters:** If a drafter fails to output these fields, the driver won't escalate immediately. Instead, it will keep starting new review rounds until it hits `MaxFixupCycles` and escalates, wasting API tokens and time.
* **Concrete fix:** Update `ValidateReviewConsensusArtifact` to validate that these fields are present and valid when `strict_gate` is enabled.

## Open questions

1. **How should we handle non-standard markdown finding headers?** If an agent formats findings using bullet points (e.g. `- [CRITICAL] Issue`) instead of headings, the scan will miss them. Should we enforce heading structures in `ValidateReviewArtifact` to prevent this?
2. **Should we block custom `checks:` execution by default?** Since Parley Deck is run locally, running arbitrary commands from files written by LLMs poses a security risk. Should we require interactive user sign-off (HITL) for any custom `checks:` command before execution?
