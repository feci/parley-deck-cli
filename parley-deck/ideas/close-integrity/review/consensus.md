---
idea: close-integrity
phase: review-consensus
drafter: claude-1
review-round: 1
date: 2026-06-24
---

## Review consensus — round 1

Three non-implementer reviewers ran Phase 6 in refutation mode against
`git diff 97034cd..HEAD` (LE-7 goal-done gate + LE-11 close guardrails):

- **codex-1** — 1 MAJOR (`review/round-01/codex-1.md`)
- **hermes-1** — 1 MAJOR, 1 MINOR, 1 NIT (`review/round-01/hermes-1.md`)
- **antigravity-1** — 1 CRITICAL, 2 MAJOR, 1 MINOR (`review/round-01/antigravity-1.md`)

The dogfood worked: refutation review surfaced one CRITICAL race the solo build missed,
plus two independent parser-bypass paths for the goal-check fail-open. All severities were
verified against the actual code before acceptance (see "Verification" below).

## Agreed fixes (all applied in fix-up cycle 1)

| ID | Severity | Raised by | Fix |
|----|----------|-----------|-----|
| **CF1** | CRITICAL | antigravity-1 (CRITICAL) + codex-1 (MAJOR) | Dedupe to distinct non-implementer IDs in `newDriverImplOps` (a `seen` map). Resolves BOTH the count-bypass (codex: duplicate IDs inflate `ReviewerCount` past the LE-11 `< 2` guard) AND the concurrent-write race (antigravity: `RunRoundOne` launches one goroutine per list element, so `[rev, rev]` makes two goroutines write the same `agents/rev/stdout.log` and `review/round-NN/rev.md`). Verified real: `selectedAgents` (runner.go:327) does not dedupe; `RunRoundOne` (runner.go:204) fans out per element. |
| **CF2** | MAJOR | hermes-1 + antigravity-1 | Harden `parseGoalVerdict`: add `` ` ``, `"`, `'`, `_` to the leading `TrimLeft` cut set (hermes: `` `GOAL-CHECK: FAIL` `` / `"…"`-wrapped tokens) and `TrimLeft` the post-colon `rest` of `*`/`` ` ``/`"`/`'`/`_`/space (antigravity: `**GOAL-CHECK:** FAIL` left `** FAIL` in `rest`). A confident wrapped FAIL is no longer swallowed by the fail-open path. |
| **CF3** | MAJOR | antigravity-1 | Give the goal-check `RunConsult` an explicit `Timeout: 2 * time.Minute`. It previously inherited the agent's full 15–30 min timeout, so a hung checker would block the driver tick that long instead of failing open fast — defeating the advisory gate's purpose. |
| **CF4** | MINOR | antigravity-1 | Reset `verdict = ""` on every matched `GOAL-CHECK:` line so a trailing ambiguous line (e.g. `GOAL-CHECK: RE-EVALUATING`) clears a prior PASS/FAIL back to ambiguous — true "last verdict wins". |
| **CF5** | MINOR | hermes-1 | Add tests for the strict_gate design-only (`StrictGate: true, AutoImplement: false`) close path: Reserved + 1 reviewer + goal PASS → complete (LE-11 is auto-only); goal FAIL → escalate (LE-7 fires on StrictGate). Documents the conditional-rigor boundary. |
| **CF6** | NIT | hermes-1 + antigravity-1 (OQ#1) | Defensive guard at the top of `GoalCheck`: if `o.drafter == o.implementer`, fail open (advisory) and run no agent, rather than running the implementer as its own checker. The upstream guards already prevent reaching this, but the contract is now enforced locally too. |

New regression tests: `goal_check_test.go` (CF2/CF4 wrapper + reset cases),
`driver_impl_le_test.go` (CF1 dedupe, CF6 no-independent-checker),
`close_integrity_test.go` + `strict_gate_test.go` helper (CF5 strict design-only path).

## Deferred follow-ups (out of scope for close-integrity — tracked, not blocking)

- **DF1** — Reject duplicate participant IDs at the workspace/CLI load boundary
  (codex-1 OQ, antigravity-1 OQ#2). CF1 closes the close-integrity-relevant race by
  deduping `o.reviewers`; a global rejection in `parseList`/`workspace.go` is broader
  defense-in-depth touching every phase and belongs in its own hardening idea.
- **DF2** — Apply the LE-7/LE-11 close guards to the `pipeline auto` block-completion
  path (codex-1 OQ#1). Verified that path is a *different subsystem*: `internal/app/pipeline_cmd.go`
  completes blocks via the pipeline DAG driver (`pipeline.Advance` + `implementationComplete()`
  reading `status: complete`), not through `internal/driver.advanceReview` where LE-7/LE-11
  live. Wiring per-block goal-check + reviewer-count semantics into the §12 pipeline is a
  separate idea.

## Resolved (no action needed)

- **codex-1 OQ#2** — `internal/runner/TestDurableKillEndToEndRealProcess` failed in codex's
  sandbox. Re-ran in the implementer's environment: `ok parley-deck-cli/internal/runner`.
  The failure is the known codex-sandbox sysctl/process-group limitation, not a regression
  from this diff.

## Verification after fix-up cycle 1

`gofmt`, `go build ./...`, `go vet`, `go test -count=1 ./...` (all packages incl. the
embedded-default drift guard) — green.

## Round-02 re-review outcome

All three reviewers re-ran Phase 6 in refutation mode against the fix-up (`git show
9267034`) and the current functions. Round-02 is clean — **zero new agreed fixes**:

- **codex-1** (`review/round-02/codex-1.md`) — No findings. Verified CF1 collapse before
  `ReviewerCount`/`RunReviewRound`, CF2 (incl. the `GOAL-CHECK: ***` over-trim → stays
  ambiguous), CF3 override precedence, CF4 last-wins + reset, CF6 fail-open. (Re-noted the
  codex-sandbox durable-kill failure — already resolved.)
- **hermes-1** (`review/round-02/hermes-1.md`) — 0 CRITICAL/MAJOR/MINOR; 1 NIT (DF3 below).
  Traced the CF3 timeout plumbing through `timeoutForAgent` to confirm the 2m override is
  honored, not a no-op. Confirmed CF2's rest-trim cut set contains no letter, so it can
  never consume the P/F of a real verdict.
- **antigravity-1** (`review/round-02/antigravity-1.md`) — No findings, no regressions; all
  six fixes survive refutation.

### DF3 — deferred NIT (dismissed for this idea, by reviewer agreement)

hermes-1 flagged that `parseGoalVerdict` uses `HasPrefix(rest, "PASS"/"FAIL")`, so
off-spec `PASSED`→PASS / `FAILURE`→FAIL. This is **pre-existing** (the round-01 parser did
the same) — CF2 did not introduce it, only widened the input shapes marginally. The error
direction is asymmetric-safe (`FAILURE`→FAIL escalates to a human; `PASSED`→PASS requires a
non-spec word past an already-passed review consensus). Both hermes-1 ("not worth a code
change for this idea, out of scope") and the drafter agree: dismiss for close-integrity;
the exact-token verdict match belongs to a future parser-hardening idea.

## Signoffs

Reviewer verdicts transcribed by the facilitator from each agent's round-02 artifact
(the canonical per-agent review files cited above are the source of truth; codex writes to
the shared consensus file are sandbox-blocked, so all three are recorded uniformly here).

### Signoff: codex-1 — 2026-06-24
Status: ✅ ACCEPT
Round-02: "No findings." All six agreed fixes verified; only residual is the known
codex-sandbox durable-kill test limitation (out of fix-up surface).

### Signoff: hermes-1 — 2026-06-24
Status: ✅ ACCEPT
Round-02: "No CRITICAL, MAJOR, or MINOR findings. The six agreed fixes are correctly
implemented and survive refutation." One pre-existing NIT (DF3) recorded as out-of-scope.

### Signoff: antigravity-1 — 2026-06-24
Status: ✅ ACCEPT
Round-02: "No findings… The fix-up completely resolves all agreed points from round-01,
is robust against refutation, and introduces no regressions or new bugs."

### Signoff: claude-1 (implementer/facilitator) — 2026-06-24
Status: ✅ ACCEPT
All six round-01 agreed fixes applied and re-review-clean (0 new agreed fixes). DF1/DF2/DF3
deferred to their own ideas; the codex durable-kill failure is a documented sandbox
limitation (green locally). Marking close-integrity complete.
