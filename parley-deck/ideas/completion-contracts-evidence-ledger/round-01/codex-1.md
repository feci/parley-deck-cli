---
agent: codex-1
idea: completion-contracts-evidence-ledger
round: 1
date: 2026-07-04
---

## Summary

The current implementation has the right insertion points but the wrong data shape. `internal/driver/impl.go` calls `ImplOps.RunChecks` before opening the first review round and again after each fix-up, then later completes from review consensus, strict-gate fields, reviewer count, and goal-check. `internal/app/driver_impl.go` implements `RunChecks` as one unstructured command path: scalar `checks:` from `00-prompt.md`, else `go test ./...`, else fail closed for `auto_implement`, else no-op for design-only.

I would add `done_when` as an opt-in verification path, not as a replacement for review judgment. When absent, the driver must take the exact existing path: same `ReadTrack` and `ReadStrictGate` readers, same `RunChecks` behavior, same output, and no evidence file. The new path should be isolated behind a new contract reader and an evidence runner so absent `done_when` is behaviorally byte-for-byte today.

## Proposed approach

Add a structured contract reader in `internal/driver/transport.go` beside `ReadTrack` and `ReadStrictGate`, but do not retrofit those scalar readers to a general YAML parser. The current `readFrontmatterField` and `protocol.ReadFrontmatter` functions only understand top-level `key: value` lines; a frontmatter block like `done_when:` would currently read as an empty scalar. Use `gopkg.in/yaml.v3` only in a new `ReadCompletionContract(ideaDir string) (CompletionContract, bool, error)` helper that extracts the frontmatter block and unmarshals just this shape:

```go
type CompletionContract struct {
    Criteria []CompletionCriterion `yaml:"done_when"`
    Hash     string                 `yaml:"-"`
}

type CompletionCriterion struct {
    ID             string        `yaml:"id"`
    Run            string        `yaml:"run,omitempty"`
    PathExists     string        `yaml:"path_exists,omitempty"`
    CWD            string        `yaml:"cwd,omitempty"`
    Timeout        time.Duration `yaml:"timeout,omitempty"`
    ExpectExit     *int          `yaml:"expect_exit,omitempty"`
    StdoutContains []string      `yaml:"stdout_contains,omitempty"`
    StdoutRegex    []string      `yaml:"stdout_regex,omitempty"`
    MaxOutputBytes int64         `yaml:"max_output_bytes,omitempty"`
}
```

Validation should fail closed when `done_when` is present but malformed: duplicate or empty IDs, neither or both of `run` and `path_exists`, absolute `cwd` or paths escaping the repo, invalid regex, negative timeouts, or conflicting `checks:` plus `done_when`. My preference is to make `done_when` the authoritative completion gate when present; if a legacy `checks:` command still matters, copy it into `done_when` explicitly rather than silently running two different contract systems.

Keep `ImplOps.RunChecks(ctx)` unchanged for legacy ideas. Add a separate method, for example:

```go
type CheckReason string

const (
    CheckPreReview   CheckReason = "pre-review"
    CheckPostFixup   CheckReason = "post-fixup"
    CheckPreComplete CheckReason = "pre-complete"
)

type ImplOps interface {
    // existing methods...
    RunChecks(ctx context.Context) (bool, string)
    RunCompletionContract(ctx context.Context, reason CheckReason, cycle int) (bool, string)
}
```

`internal/driver/impl.go` then gates like this:

- In `advanceImpl`, if `ReadCompletionContract` is absent, call existing `RunChecks`; if present, call `RunCompletionContract(CheckPreReview, 0)` before `OpenReviewRound`.
- In `advanceReview` after `Fixup`, use the same branch before writing `.fixup-done` and opening the next review round.
- In `advanceReview` when `OutstandingAgreedFixes == 0`, run `RunCompletionContract(CheckPreComplete, round)` before `GoalCheck` and `Complete`. This makes stale evidence hard to close on; the last ledger entry is produced at the actual close boundary. Review consensus still controls judgment. A green ledger alone never opens, closes, or accepts a review.

Implement `RunCompletionContract` in `internal/app/driver_impl.go`. It should execute command criteria with `exec.CommandContext(ctx, "sh", "-c", criterion.Run)` from repo root or validated relative `cwd`, and evaluate file criteria using `os.Stat` under the repo. Each criterion gets a per-criterion timeout, defaulting to a conservative value such as 10 minutes and clamped by the driver context. On Unix, use a process group so `sh -c` children are killed on timeout; otherwise a hung child can survive after the shell is killed.

Write the ledger to `parley-deck/ideas/<slug>/review/evidence.md`. Use markdown for protocol consistency, but make each entry a fenced JSON object so the driver can parse it deterministically:

````markdown
---
idea: completion-contracts-evidence-ledger
artifact_kind: completion_evidence_ledger
schema: 1
---

## evidence-2026-07-04T12:00:00Z-pre-review

```json
{
  "schema": 1,
  "idea": "completion-contracts-evidence-ledger",
  "reason": "pre-review",
  "cycle": 0,
  "contract_sha256": "...",
  "workspace_head": "...",
  "workspace_state_sha256": "...",
  "all_passed": true,
  "criteria": [
    {
      "id": "go-tests",
      "status": "passed",
      "exit_code": 0,
      "duration_ms": 1234,
      "stdout_sha256": "...",
      "stderr_sha256": "...",
      "stdout_bytes": 12000,
      "stderr_bytes": 0,
      "timed_out": false,
      "matched_expectations": ["expect_exit", "stdout_contains:Usage:"]
    }
  ]
}
```
````

Do not store raw stdout or stderr. Store byte counts, SHA-256 digests, timeout status, exit code, duration, and which named expectations matched. Compute `contract_sha256` over the normalized contract and include a workspace state digest. In a git repo, include `HEAD`, `git status --porcelain`, tracked diff hash, and untracked file names plus content hashes. The pre-complete run is still the primary stale-evidence defense, but the digest gives reviewers a way to spot "evidence was for a different tree" cases.

Tests should cover: absent `done_when` uses only the old `RunChecks` branch and creates no ledger; valid command pass/fail; file existence pass/fail; timeout; malformed contract; output pattern mismatch; no raw command output in `evidence.md`; pre-complete reruns the contract before `Complete`; and a contract hash change invalidates older evidence.

## Concerns / open questions

The largest design question is whether `done_when` belongs in frontmatter or in a body section. Frontmatter is convenient for machine gates, but the current readers are intentionally scalar and line-oriented. A dedicated YAML frontmatter reader for this one block is acceptable because `yaml.v3` is already in `go.mod`, but it should not change `ReadTrack`, `ReadStrictGate`, `ReadAutoImplement`, or existing status readers.

`checks:` interaction should be explicit. I recommend rejecting `checks:` plus `done_when` together, because otherwise the driver has two verification authorities and reviewers may not know which one controls completion. A migration can be simple: `checks: go test ./...` becomes a `done_when` criterion with `id: go-tests`.

Output pattern checks and secret safety pull in opposite directions. If the driver never stores raw output, reviewers lose easy debugging context. I think that is the right default: evidence proves the check ran and passed, while debugging remains a local rerun. If truncated output is ever allowed, it should require an explicit per-criterion `capture: redacted_tail` opt-in and should still pass through secret-pattern redaction before writing.

Path criteria need a clear root. I would resolve all paths relative to workspace root unless `cwd` is supplied, reject `..` escapes, and record normalized relative paths in the ledger. Absolute paths make evidence non-portable and can leak local machine details.

## Risks

Hung commands are the most concrete runtime risk. The existing `RunChecks` relies on the driver context and `exec.CommandContext`; a `done_when` loop with multiple commands needs per-criterion timeout, process-group cleanup, and bounded output readers. Otherwise one bad criterion can hang the driver tick or leave child processes running.

Stale evidence is the main correctness risk. Running the contract before review is useful for reviewers, but it is not enough for Phase 8 completion if files changed later. The driver should run the contract again at the pre-complete boundary whenever `done_when` is present, and the ledger entry should carry contract and workspace hashes so stale evidence is visible rather than trusted by prose.

Secret leakage is the main audit-trail risk. Commands execute with the same trust model as `checks:` today, so they can print secrets from the environment. The ledger must never persist raw output, environment variables, full command-expanded values, or matched output text. It should record only command strings already visible in `00-prompt.md`, exit metadata, durations, byte counts, and digests.

Backward compatibility can regress accidentally if the contract parser is wired into common frontmatter reads. Keep the new reader and new runner behind a presence check. For an idea with no `done_when`, `advanceImpl`, post-fix-up, review consensus, strict gate, goal-check, and `Complete` should behave exactly as they do now, including the legacy `checks:` and `go test ./...` fallback behavior.
