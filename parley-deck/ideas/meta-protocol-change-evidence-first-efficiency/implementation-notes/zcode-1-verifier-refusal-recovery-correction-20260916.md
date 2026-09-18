---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
base-commit: 6962f2af9b6e46548588c5152285e6808f3d6f7f
kind: correction (code written; tests NOT executed in this launch)
corrects: zcode-1-verifier-refusal-recovery-candidate-20260916.md (left byte-unchanged)
---

# Verifier-launch refusal recovery correction (zcode-1)

## Coordinator execution evidence (PRIMARY, quoted from retained logs)

The coordinator executed the unmodified candidate. Native JSON logs retained at
`.parley-runtime/input/first-check/` show, in `stdout.log`: trajectory build
failure at `verification_recovery_test.go:325:6`, `:326:22`, `:327:66`
(`undefined: err`) and `:447:10` (`declared and not used: path`); the runner
package failed separately with `compile: writing output: write $WORK/b148/
_pkg_.a: no space left on device` (disk exhaustion, environmental — it excuses
nothing and proves nothing about the runner code). `result.json` records
exit 1, both packages `[build failed]`, no test executed. Provenance: PRIMARY —
my own reads of those files in this checkout.

## What changed (code; compile-checked by reading only — Bash stayed disabled)

**Compile fixes** (`verification_recovery_test.go`): the repeated/conflicting
apply loop now uses a scoped `applyErr :=` (the outer `err` was if-scoped at
the first apply — the exact 325–327 failure);
`TestFreshResolutionValidatesRecoveredLineageAndFailsClosed` discards the
unused journal path (`ticket, _ :=` — the 447 failure). No assertion was
weakened; both tests keep their original checks.

**Correctness fix — recovered reads revalidate the exact refused evidence**
(`verification.go`, `verification_recovery.go` comment only). Review finding
accepted: `readVerificationLaunch` validated only the SHAPE of
`RefusedRequestedSHA256`/`RefusedTerminalSHA256` and otherwise trusted
`recovery.json`. It now takes the `VerificationTicket` (internal signature;
all five call sites in `verification.go` updated — `Stop`, `Reserve`,
`Execute` claim, `Execute` authority guard, journal read; no caller outside
this file existed), and when a recovery is present it requires:

- the retained `recovery.json` to be structurally bound to the canonical
  ticket SHA and the ACTUAL current `launch.json` digest and invocation, with
  `At` not before the launch reservation (reversed chronology refused), and
  the bound invocation distinct from both the refused invocation and the
  implementer's (`ticket.Request.InvocationID`) — malformed or contradictory
  retained artifacts are refused even though admission would never write them;
- the refused invocation's retained `requested.json`+`terminal.json` to be
  revalidated in full (same semantics as admission: metadata bound to this
  ticket's run/idea/verifier, phase `trajectory-verification`, headless,
  `failed`/`budget_refused`, no PID/exit code, no `started.json`, sane
  chronology against `launch.At`) AND their CURRENT digests to equal the two
  digests the recovery bound. Mutation, deletion, or a fabricated lifecycle
  after recovery fails every recovered reserve, stop, read, execute and
  receipt — including stale refused-invocation handles, which now surface the
  lineage error rather than superseded when the lineage itself is broken. The
  superseded/invocation resolution happens only after the evidence holds;
  without `recovery.json` behavior is byte-identical to before.

The admission path (`RecoverCapturedVerificationLaunch`) is unchanged; the
recovery artifact schema is unchanged; original refusal records, launch
reservation and accounting remain untouched by all of this.

**Tests** (written, NOT executed): the forged-artifact block in
`TestRepeatedAndConflictingRecoveryApplyFailsClosed` became a five-variant
table (foreign prior launch; correct structure with wrong evidence digests —
previously claimed beyond byte evidence, now refused by digest equality;
reversed chronology; self-bound and implementer-bound invocation), each from a
genuinely admitted recovery so only the varied binding is wrong.
New `TestRecoveredRefusedEvidenceMutationFailsEveryHandle`: four post-recovery
lineage mutations (canonical self-consistent rewrite changing only an ordinal
no rule inspects — catchable solely by digest equality; requested deleted;
terminal deleted; started lifecycle added), each asserting reserve/stop/read/
execute fail for BOTH the current recovered handle and the stale refused
handle, with no journal write and inventory unchanged. New
`TestRecoveredReceiptRevalidatesRefusedEvidence`: a completed, journal-valid
receipt cannot be read over later-mutated refusal evidence. New runner-level
`TestVerifierRecoveryRefusedEvidenceMutationFailsHandles`: evidence produced
by the real public `RunMeasured` refusal path (not the telemetry fixture) is
rewritten canonically after recovery; stale and current handles fail closed
before any write or charge and the evidence stays retained.

## Not executed in this launch

`go build`, `go vet`, `gofmt`, and every test above. I verified the two
reported compile errors against the source and re-read every edited region for
consistency (signatures, call sites, imports, variable use), but no toolchain
ran here: nothing in this note is PASS evidence, and the disk-exhaustion
runner result from `first-check` remains an environmental non-result.

## Explicitly incomplete obligations (unchanged from the candidate)

1. **Resolution seam still open.** `deriveParentEvidence` (reconcile.go, not
   owned) still reads `launch.json` directly and requires that invocation's
   successful terminal; `TestFreshResolutionValidatesRecoveredLineageAndFailsClosed`
   still asserts the fail-closed behavior. Wiring the resolution reader (and
   the attended app control that hands the bound invocation to a real verifier
   relaunch) is NOT done and NOT designed away by anything here.
2. **A real positive bound-runner relaunch riding the recovery is still
   untested** — the runner generates invocation IDs internally; only the
   API-level idempotent reservation and the runner-level mismatch refusal are
   covered. These remain separate obligations; I do not claim them complete.

## Honesty boundary (unchanged)

Digests bind retained bytes; they do not authenticate writers. A same-UID
forger able to rewrite both the evidence and `recovery.json` consistently can
still fabricate a self-consistent lineage — that remains outside byte
evidence, as the candidate stated. What the fix removes is trusting
`recovery.json` over diverging retained evidence, and shape-only validation of
the bound digests. Original refusal accounting is preserved, never refunded.
