---
idea: meta-protocol-change-lean-organizer
author: kimi-1
role: independent diagnostic delta reviewer (non-implementer)
artifact: release-linux-diagnostic-review
date: 2026-09-24
delta: dd0a226..bc5ff29 (single commit, parent verified: merge-base = dd0a226)
inputs: release-linux-hosted-diagnostics-claude-1.md, release-linux-diagnostic-followup-zcode-1.md
verdict: PASS — safe to push for ONE untagged hosted ubuntu diagnostic run; NOT a U2 fix, NOT release v1.49.1 readiness
method: static source/test inspection only; no edits, no rerun, no push, per dispatch
---

# Linux diagnostic delta review dd0a226..bc5ff29 — kimi-1

## Verdict

**PASS.** The delta is exactly what was dispatched: one 15-line product hunk
(`internal/evidence/execute.go` +15/−0), one test-helper guard
(`internal/acp/spawn_test.go` +21/−2), one new 161-line test file
(`internal/evidence/execute_runerror_test.go`, `//go:build !windows`), plus
zcode-1's own record. Nothing else. Closed artifacts untouched; claude-1's
report left untracked for the organizer; VERSION/CHANGELOG/version.go not in
the delta. Provenance sha256s in zcode-1's record verify exactly against my
tree (`shasum -a 256`): execute.go `094d4f17…9e6417`, spawn_test.go
`374bcdf5…8c9501`, execute_runerror_test.go `819ecbd5…a265d`.

## U2 diagnostic preservation — each required property verified in source

| Property | Verdict | Provenance |
|---|---|---|
| Fires only on non-ExitError with zero command output | PASS | `execute.go:219-221`: `runErr != nil && len(out) == 0` plus `!errors.As(runErr, &exitErr)` — same `errors.As` idiom as the pre-existing complete-gate at :287 |
| Bounded, scrubbed error text | PASS | `execute.go:222` wraps the reason in `ScrubAndTruncate` (:636-654): credential patterns :623-632 redacted to `«redacted»`, tail-bounded to 100 lines / 4096 bytes + 3-byte `…` marker. Scrub happens at construction — verified `writeVerificationArtifact` (`internal/trajectory/verification.go:194-214`) adds only canonical-JSON + 1 MiB bound, no second scrub; zcode-1's layering claim is exact |
| OutputSHA256 remains emitted-output hash | PASS | `execute.go:197` hashes the streaming `hasher` fed by `io.MultiWriter` (:143-147), computed before the hunk and never touched by it. Test-pinned both ways: empty-string hash on the preserved branch (execute_runerror_test.go:42) and exact output hash when output exists (:138). I re-verified empty-sha256 = `e3b0c442…b855` locally — the hosted signature |
| No false successful/complete envelope | PASS | Status (:260-278) and complete (:279-289) untouched; any `runErr != nil` forces `StatusFail`, and non-ExitError forces `complete = false` at :288. Tests pin Complete=false/StatusFail/exit −1 on every preserved-branch test (:36, :74, :131). Downstream `Diagnostics` consumers are classification-neutral: `trajectory.go:156` is a size guard (8192; max preserved text 4099), `app/evidence_table.go:188` renders only. My own repo-wide grep finds no assertion of empty diagnostics outside the new file — zcode-1's claim holds |
| Lifecycle unchanged | PASS | Process lifecycle block `execute.go:184-191` (KillGroup/Wait) untouched; the hunk is a Diagnostics-only assignment after collection, before envelope parsing. Overflow prefix (:207-209) is unreachable when `len(out) == 0`; envelope/test2json parsers see empty output and stay inert, so the preserved text is never clobbered or prepended |
| Privacy tests meaningful, plausible error strings | PASS | Scrub test (:67-84) carries a plausible credential-bearing refusal (`bearer aBcD.eFgH-token123`, `api_key=ghp_…`×28) and asserts raw secrets absent AND `«redacted»` present AND label present — it fails in both directions, not vacuous. Plausible branch texts exercised: controller refusal (:34) and a real OS `cmd.Start()` failure via absent cwd (:55-63). Bounded test pins `len ≤ evidenceMaxBytes+3` + truncation marker (:94-99) |

Test-suite shape: six tests; the four preservation tests fail without the hunk
by construction (they assert text that only the hunk produces), the two
negative tests (output-present :105-141, ExitError-skip :146-161) pin the
pre-existing shape in both worlds — zcode-1's claimed bidirectional proof
(4 FAIL → restore → 6 PASS) is structurally sound. I did not re-execute it
(no-rerun dispatch); this is static concurrence, not independent execution.

Test 5's design note (:116-122) is correct: the unreaped supervisor is a
zombie that answers a signal-0 probe as alive, so the marker file is the only
clean completion signal, and `cmd.Wait` drain-to-EOF makes the piped bytes
deterministic. Sound reasoning, correctly applied.

## D1 READY guard — spawn_test.go:101-122

- **Guard correct.** Buffered channel (cap 1, :110) + reader goroutine +
  `select` against the file's established 10 s `time.After` idiom (sibling
  F-9 sites :152-159, :172-179, :197-204). Success path is behavior-identical
  to the old blocking `ReadString` except for the bound; the goroutine's send
  can never block, so the guard itself introduces no goroutine leak.
- **Success path: no leaked goroutines or child processes.** The reader
  goroutine exits with its send; both callers (:146, :166) `Stop`/`Wait` the
  child; the observer is released. `p.Stdout()` is the raw `cmd.StdoutPipe()`
  read end (`spawn.go:97-98`) with no competing internal copier — the race-
  detector pass zcode-1 reports is consistent with the shape.
- **Guard-firing path (unreachable on 64 KiB-pipe hosts):** the reader
  goroutine, the stderr copier (parked in `gatingObserver.Write`), and the
  unstopped child outlive the helper until test-binary exit — bounded, loud
  (named 10 s failure instead of the 45 m binary-wide panic), and matching
  the file's established F-9 guard idiom, which leaves the same residuals on
  its own timeout paths. The comment discloses the reader goroutine. A
  tighter cleanup is not trivially available: `Stop` drains `p.wg` before
  reaping (`spawn.go:138`), so stopping without releasing the observer would
  itself hang — the honest fix would add moving parts no available host can
  exercise. Accepted as-is.
- zcode-1's honest limitation stands: the stall cannot be induced on
  Linux/macOS, so no local mutation discriminates the guard; the bound is
  structural. No build tag touched — `spawn_test.go` stays Windows-eligible
  exactly as claude-1 found it. **No Windows claim made or implied; the
  owner-blocked Windows track is untouched.**

## Findings (all reportable; none blocking)

1. **Nit — record accounting:** zcode-1's "+19/−2" for spawn_test.go is the
   net; `git diff --numstat` shows **+21/−2**. Cosmetic only.
2. **Nit — tail-bound evicts the label:** a run error > 4096 bytes loses the
   `run error (no command output):` prefix (ScrubAndTruncate keeps the tail).
   Disclosed in zcode-1's record; the bounded test pins length, not label.
   All realistic branch reasons (attribution facets, exec errors) are far
   shorter. Accepted as proposed.
3. **Informational — pre-existing asymmetry:** `exitCodeOf` (:660) uses a
   bare type assertion while the new branch and the complete-gate use
   `errors.As`; a wrapped ExitError would keep the old exit −1/empty-
   diagnostics shape. Not reachable from os/exec's direct returns; no
   regression introduced.
4. **Informational — build tag:** the new test file's `!windows` follows the
   in-package precedent (`execute_cancellation_test.go`); the controlled-path
   tests it mirrors (`execute_control_test.go`) are untagged. No Windows
   contract change either way.

## NOT established / NOT claimed

- **U2 root cause remains UNKNOWN.** This change only stops discarding the
  reason; candidates (A)/(B) remain claude-1's deduction, now observable.
  A green hosted run would prove nothing about U2 (claude-1's own framing).
- The full suite was not re-run by zcode-1 (last full evidence:
  release-ci-fix-zcode-1.md) and not re-run by me (dispatch). Named
  dependents (trajectory trio, runner, app) reported green by zcode-1; my
  concurrence on their indifference is static (grep + consumer analysis).

## Safe-to-push judgment

**Safe to push for exactly one untagged hosted ubuntu diagnostic run** over
the touched packages (`evidence, acp, trajectory, runner, app`), per
zcode-1's hosted-next-step. The only product effect is that a previously
discarded error string is persisted — scrubbed, bounded — into an
already-failing, already-incomplete envelope; status, completeness, and the
output hash are provably untouched, so the change cannot mask, skip, or
green anything. On a U2 recurrence the retained `step-00N.json` now
discriminates candidate (A) (named attribution facet) from (B)
(artifact-write/controller text) with zero further guessing.

This is **not** release readiness: product source changed, so any publish is
a new immutable version by organizer/owner decision, and this review gates
only the diagnostic push. Windows scope remains with the owner-blocked track.

## Limitations

- Static review only: no build, vet, test, mutation, or rerun performed, per
  dispatch. zcode-1's validation numbers are cited, not reproduced.
- The D1 guard-firing path is unexercised on any available host (structural).
- I reviewed the hosted evidence only as quoted in claude-1's report; I did
  not independently pull run 35998165067's logs.
