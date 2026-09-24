---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent Linux CI diagnostic reviewer (non-implementer)
artifact: release-linux-hosted-diagnostics
date: 2026-09-24
run: 35998165067 (Tests, push, main@dd0a226) — ubuntu job 107627994770
verdict: U1 PASS · U2 FAIL (recurred, now localized) · Windows PENDING (not waited on)
---

# Linux hosted diagnostics — claude-1

## Verdict

| Item | Status | Basis |
|---|---|---|
| **U1** acp stderr drain | **PASS** | `ok parley-deck-cli/internal/acp 2.659s` — first hosted-ubuntu green for this package |
| **U2** captured-verification incompleteness | **FAIL (recurred)** | `FAIL internal/trajectory 164.826s`, 2 subtests |
| U2 root cause | **NOT established** — but narrowed from "somewhere in execution" to **two candidate branches in one callback** | dumped envelopes, below |
| D1 test-guard residual | **still open**, untouched by this run | see separate section |
| Windows leg | **PENDING** — still `in_progress` when I classified Linux. **No waiver implied.** | per dispatch, Linux was not delayed for it |

Ubuntu job: 12:17:56Z → 12:21:31Z; Test step 12:18:19Z → 12:21:28Z (3m09s).
**30 packages `ok`, 1 `FAIL`.** `internal/runner` (56.772s) and `internal/app` (152.912s)
— both of which failed in earlier runs — are **green here**. `internal/evidence` and
`internal/procctl` also green.

## U1 — PASS, and it is the fix that did it

`TestSpawnObservesAllStderrAndRealExit` failed on the previous hosted ubuntu run with
`spawn_test.go:38: observed bytes: 0` (of 16384). At `dd0a226` the acp package is green in
2.659s. This is the only product-source change in the delta (`internal/acp/spawn.go`), it
targeted exactly that mechanism, and the failure it targeted is gone on the same platform
that produced it. **U1 is closed.**

U1 is also *causally independent* of U2: `internal/evidence` and `internal/trajectory` do
not import `internal/acp` (the only importer is `internal/runner/acp.go`). The U1 fix
neither caused nor masked the U2 recurrence.

## U2 — recurred; what the envelopes actually prove

Two failures, both in `internal/trajectory`:

- `TestParentRecoveryPublicationInterruptionAndExactReplay/rename` — `parent_recovery_test.go:126`
- `TestRecoveredParentObservationRequiresExactLineage/original-result-rewritten` — `verification_parent_recovery_test.go:510`

The new dump fired at both sites and printed the full retained set. **Both sites carry a
byte-for-byte identical signature:**

| Envelope | Site 1 | Site 2 |
|---|---|---|
| retained | `receipt`, `step-001..004`, `process-001..003` | same |
| **`process-004.json`** | **absent** | **absent** |
| `step-004.execution.complete` | `false` | `false` |
| `…command.exit_code` | **`-1`** | **`-1`** |
| `…command.duration_ms` | **`1519`** | **`1521`** |
| `…command.output_sha256` | `e3b0c442…7852b855` | same |
| `…command.diagnostics` | `""` | `""` |
| `…command.format` | `shell` | `shell` |
| `step-004.process_sha256` | absent/empty | absent/empty |
| receipt `steps` / `failure_stage` | `4` / `execution` | `4` / `execution` |

Contrast with the healthy `step-003` at the same site: `complete: true`, `exit_code: 1`,
`duration_ms: 24`, real `PARLEY-EVIDENCE` diagnostics, `process_sha256` present. And
`process-003.json` shows a fully attributed identity: `pid == pgid == 37959`, `boot_id`,
`proc_start`, `marker: captured-criterion`, full supervisor `command`.

### The deduction chain (each step from evidence + read source)

1. `output_sha256` = `e3b0c442…` is **SHA-256 of the empty string** (verified locally).
   The supervised criterion produced **zero bytes**. Combined with `format: shell` (no
   envelope parsed), the criterion command **never ran**.
2. `complete=false` with `exit_code = -1` excludes every other incompleteness branch in
   `evidence/execute.go:264-273`: not `capbuf.overflow` (zero bytes), not
   `envelopeInvalid` / `goTestInvalid` (no output to parse, `format` stayed `shell`), not
   the `control != nil && ExitCode >= 128` signal-death rule at `execute.go:268-270` (−1 is not ≥128), not a ctx
   deadline (subtests ran 2.95s/3.25s against a 60s helper ctx).
   The surviving branch is `execute.go:273` — `errors.As(runErr, &exitErr)` being **false** —
   `exitCodeOf` (`execute.go:641`) returns `-1` for any **non-`*exec.ExitError`**.
3. `duration_ms` of 1519/1521 ms — within 2 ms of each other across two independent tests,
   against a 24 ms healthy baseline — is **`procctl.KillGroup`'s full grace window**
   (`procctl_unix.go:41`, `1500 * time.Millisecond`, SIGTERM → poll → SIGKILL). That path
   runs **only when `cmd.Process != nil`** (`execute.go:184-191`). The window elapsed in
   full because the supervisor script opens with `trap ':' TERM`, i.e. it ignores SIGTERM.
   **Therefore `cmd.Start()` succeeded and the shell was live and already past its `trap`
   line.** This is not a spawn/fork failure and not a missing binary.
4. `process-004.json` is absent while `step-004.json` exists. `process-%03d.json` is
   written at `verification.go:620`, strictly **after** `start()` returns and strictly
   **before** `release()`. `release()` is what sends `go\n` to unblock the supervisor —
   consistent with (1), it was never reached.

**Conclusion: the error originates inside the `start()` callback or the line immediately
after it — after a successful `cmd.Start()`, before `release()`.** That is exactly two
candidates:

- **(A)** `procctl.Attributed(sp)` refused the freshly started supervisor →
  `"criterion supervisor identity unavailable: <reason>"` (`execute.go:163-166`).
  `Attributed` fails closed on 8 distinct facets (`procctl.go:107-150`).
- **(B)** `writeVerificationArtifact(dir, "process-004.json", …)` failed
  (`verification.go:620`).

Everything earlier in the closure (`checkAuthority`, `withVerification` preamble) is
**ruled out** by (3): those return before `start()`, leaving `cmd.Process == nil`, so the
1500 ms `KillGroup` could not have run.

### Why it is not established, and not "flaky"

The reason string is **discarded by product code**: `runErr` reaches
`RunCriterionControlled` and is used only for `exitCodeOf` and `status`; its text is never
written into `CommandEvidence.Diagnostics` (which holds captured child output only, here
empty). So (A) vs (B) — and, under (A), which of the 8 facets — is **not recoverable from
this run**. I am not labelling this flaky: it is load- and scheduling-sensitive, which is a
statement about *trigger*, not about *cause*, and the cause is a specific unlogged branch.

## Cross-run comparison

| Run | Commit | Ubuntu failures |
|---|---|---|
| 35983818751 | `137b1c1` | `TestEvidenceVerifierProductionClosure` ×15 (git identity — since fixed), trajectory `…ExactLineage/original-result-rewritten`, app `…RestoresPreviouslyBoundFacts` |
| 35987916696 | `4e0c069` | **U1** `TestSpawnObservesAllStderrAndRealExit` (`observed bytes: 0`), runner `:229`, trajectory `:124 after-publication`, trajectory `…RefusesNonRefusalOriginal :142`; receipts `Steps:1` |
| **35998165067** | **`dd0a226`** | trajectory only: `:126 rename`, `:510 original-result-rewritten`; receipts **`Steps:4`** |

Two things move every run: **which** test fails, and **which ordinal** fails (`Steps:1` →
`Steps:4`). The *signature* does not move: `FailureStage: execution`, incomplete captured
execution. The identity-fix and U1 regressions are monotonically retiring; the U2 family is
shrinking (3 packages → 1) but has not gone away. A green run would still not have
root-caused it — only this dump did any localization.

## Minimal repair proposal — persist the discarded error

**One change, product-side, ~4 lines, no behaviour change, no test weakening:** in
`internal/evidence/execute.go`, when `runErr != nil` and the captured output is empty,
record the scrubbed run error into `ce.Diagnostics` (the field is already
`ScrubAndTruncate`-bounded and the envelope is already secret-scrubbed before persistence):

```go
if runErr != nil && len(out) == 0 {
    ce.Diagnostics = ScrubAndTruncate("run error (no command output): " + runErr.Error())
}
```

This is the smallest change that names the exact branch. `Attributed` already returns a
precise `reason` string and `execute.go:165` already wraps it — the text simply has to
survive into the envelope. With it, the next hosted ubuntu failure distinguishes (A) from
(B), and under (A) names the facet, with zero further guessing.

**Recommended next diagnostic action: this change plus one untagged hosted ubuntu run.**
Not a rerun of the current tree — a rerun would reproduce the same unlabelled `-1` and
teach nothing new. **I do not recommend repeated whole-suite runs to chase green**; the
suite is already 30/31 and re-running it does not disambiguate the branch.

Because this touches product source, it publishes only as a new immutable version per
protocol — organizer/owner decision, not mine. `v1.49.0` (`06e563e`) stays unmoved.

## D1 residual (separate; not a Windows waiver)

D1 from my prior delta review is **unchanged and still open** — this run neither exercised
nor closed it:

1. `spawnGatedChild` (`internal/acp/spawn_test.go:101`) still reads the READY handshake via
   an **unguarded** `bufio.NewReader(p.Stdout()).ReadString('\n')`. On a platform whose pipe
   buffer is below the 16 KiB in flight, this blocks to the 45m global timeout as a
   binary-wide panic, rather than failing at a named guard — the exact mode F-9 removed two
   functions below in the same file. The same 10s guard would close it.
2. `spawn_test.go` carries **no build tag** (confirmed again at `dd0a226`), so the gated
   tests are Windows-**eligible** and D1 becomes live on the first Windows execution.

Linux cannot surface D1 (64 KiB pipe capacity) and did not: acp passed in 2.659s. **This is
a LOW durability nit, not a Windows result.** The Windows leg of this run was still
`in_progress` when I classified Linux; I make **no Windows claim, and nothing here waives
the owner-blocked Windows track.**

## Limitations

- I read the hosted ubuntu log for job 107627994770 only. No local reproduction of U2 was
  attempted and none exists on any machine to date.
- Candidates (A)/(B) are derived from the retained envelopes plus source reading. The
  discriminating string is absent from this run's evidence, so I assert the **branch**, not
  the **reason**.
- I considered and reject the Go `Setsid` race as an explanation for (A): Go's `forkExec`
  blocks the parent on the exec error-pipe, so `setsid` has completed before `Start()`
  returns. `process-003.json` (`pid == pgid`) is consistent with that.
- macOS and Windows legs of run 35998165067 had not completed; both are unclassified here.
- No code edit, push, rerun, release, install, tag or version bump was performed. No
  credentials or secrets appear in this report.
