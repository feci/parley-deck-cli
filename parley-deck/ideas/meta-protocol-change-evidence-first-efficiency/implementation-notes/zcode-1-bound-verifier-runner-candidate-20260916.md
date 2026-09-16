---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
base-commit: 6962f2af9b6e46548588c5152285e6808f3d6f7f
kind: implementation-candidate (code written; tests NOT executed in this launch)
extends: zcode-1-verifier-parent-lineage-candidate-20260916.md (and both earlier notes; all left byte-unchanged)
---

# Bound verifier runner candidate (zcode-1)

## Coordinator execution evidence (PRIMARY only for my reading of the retained logs)

Terminal check 195.018s, source unchanged: the trajectory package produced 94
pass events, but the recovered parent-lineage positive test and the clean
adversarial setups failed at the chronology check ("original started lifecycle
differs from the verifier terminal") — the observed defect this note's Task A
fixes. The runner package did not link (`no space left on device`, native disk
exhaustion) — environmental, NOT a runner source compile defect, and claiming
nothing about runner code. Full retained evidence:
`.parley-runtime/input/parent-focused/`. No checks passed for the newly changed
parent path; nothing here is PASS evidence.

## Task A — recovered chronology defect (internal/trajectory/reconcile.go)

The parent-lineage candidate applied the ordinary reservation window
(`requested <= launch.At <= started`) to a recovered replacement's authority
too. That is impossible by construction: `RecoverCapturedVerificationLaunch`
runs BEFORE the replacement invocation exists, so every lawful replacement is
requested AFTER `recovery.At` and the ordinary predicate refused the clean
case. The lawful recovered order is: original refusal completed <=
`recovery.At` <= replacement requested <= replacement started <= terminal
completed. Fix in `deriveParentEvidence`:

- Ordinary unrecovered predicate byte-identical to before
  (`!launch.At.Before(requested) && !launch.At.After(started)`).
- Recovered predicate: `!replacement.requested.Before(recovery.At)` — the
  replacement's request may not predate the recovery that authorizes it. The
  refusal-completion lower bound stays in the common authority reader
  (`readVerificationLaunch`), and `requested <= started <= completed` stay
  bound by the terminal record itself.
- No test fixture was changed to manufacture older timestamps; the existing
  real-now fixtures (requested after recovery) are now the lawful positive, and
  the `replacement-lifecycle-predates-recovery` variant still fails with the
  same message. `verification_recovery.go` and `verification.go` are untouched
  in this slice: recovery schema and authority readers unchanged.

## Task B — real runner bridge under the pinned invocation ID

**Telemetry primitive** (`internal/telemetry/record.go`): `BeginBound(directory,
id, metadata)` plus `ValidBoundInvocationID(id)`. `Begin` is refactored into
shared helpers (`openInvocationStore`, `newInvocation`) with unchanged
behavior; `BeginBound` requires a safe label that is also exactly one path
element (no separators, no dot segments — no path escape or aliasing), requires
the shared invocations directory to be a real non-symlink directory, and
exclusively creates the invocation directory fresh (`os.Mkdir`; a pre-existing
directory, file or symlink at the name is refused, never reused or
overwritten). Started/terminal records are written by the returned invocation
through the ordinary lifecycle the real process runner already uses.

**Trusted recovery context** (`internal/runner/trajectory_verification.go`):
`WithCapturedVerificationRecovery(ctx, ticket, intendedInvocation)` is the
narrow bridge. It freezes the ticket (same discipline as
`WithCapturedVerificationRecovery`'s `WithCapturedVerification`), requires a
valid bound ID, then validates the recovery authority BEFORE any launch exists
and with NO write: existence-only Lstat of the journal's `launch.json` and
`recovery.json` (never parsed here), followed by
`trajectory.ReserveCapturedVerificationLaunch` as the authority probe — under a
retained launch reservation this either fully revalidates the recovered
lineage and admits only the recovery's bound invocation, or fails closed; it
cannot write in that state, and the never-refused case is refused at the Lstat
before the probe (asserted by a dedicated no-write test). The context also sets
the captured-verification ticket binding.

**Launch wiring** (`internal/runner/telemetry.go`): `beginLaunch` resolves the
bound invocation through `boundRecoveryInvocation` — refusing, before any
telemetry directory exists, a bound launch without the exact verifier shape
(ticket match, phase `trajectory-verification`, run/idea/verifier match,
headless, never a handoff) — and records the request via
`recordBoundLaunchRequest`→`telemetry.BeginBound`. Every ordinary launch path
(including protocol-refusal retention, which keeps its unchanged
`recordLaunchRequest` signature) still generates a fresh ID; `LaunchInfo`,
`retry_of` and phase labels never select one. The launch then passes the
EXISTING scoped reservation checks
(`reserveCapturedVerification`→`ReserveCapturedVerificationLaunch` revalidates
the whole recovered authority at the actual reservation) and all normal budget
admission/settle (`reserveBudget`/`settleBudget` untouched; no generic
cycle/step budget changed; recovery is not a refund, reset, bypass or retry).
Exactly-once: the pinned directory is exclusively allocated, so concurrent or
replayed launches with one recovery ID cannot execute twice.

**Recovery runner API for the app (exact signature and requirements):**

    ctx, err := runner.WithCapturedVerificationRecovery(ctx, ticket trajectory.VerificationTicket, intendedInvocation string)

1. `ticket` is the prepared verification ticket of the refused launch;
2. `intendedInvocation` MUST be exactly the `invocation_id` the immutable
   `recovery.json` binds — the ID the operator/app chose when calling
   `trajectory.RecoverCapturedVerificationLaunch`;
3. use the returned ctx for exactly ONE verifier relaunch:
   `RunMeasured(ctx, ExecOptions{Root: ticket.Root, Agent: <headless agent whose ID is ticket.Request.Verifier>, Prompt: …, Timeout: …, Info: LaunchInfo{RunID: ticket.RunID, Idea: ticket.Request.Idea, Phase: runner.CapturedVerificationPhase}})`;
4. budget authorization is separate and mandatory:
   `runner.WithLaunchBudget(ctx, runner.LaunchBudget{Store: <fresh authorized store>, …})`
   — the recovery context grants no budget;
5. the captured helper half still executes under the pinned invocation during
   the observed launch window (the verifier process runs it in production; the
   tests drive `trajectory.ExecuteCapturedVerification` concurrently, which is
   the same journal semantics).

Attended app CLI wiring (surfacing the recovery to the operator, reading or
re-deriving the bound ID, invoking this API) is a SEPARATE subsequent
obligation and is NOT done or claimed here.

## Tests written (NOT executed; no Bash/toolchain ran in this launch)

`internal/telemetry/bound_invocation_test.go`: fresh exclusive allocation with
ordinary lifecycle and permissions; `ValidBoundInvocationID` table; every
reuse/collision refused (completed invocation, plain dir, file, symlink alias)
with store inventory unchanged and nothing written through the alias; unsafe
IDs and symlinked stores refused with nothing created; 8-way concurrent
`BeginBound` allocates exactly once.

`internal/runner/verification_bound_recovery_test.go`, all through real local
`RunMeasured` with /bin/sh fixture agents (no provider calls): the original
budget refusal produces no process; the positive relaunch executes under the
pinned ID with NATIVE requested/started/terminal (chronology asserted against
`recovery.At`), verifier work exactly once, the captured helper executing once
inside the observed window, retained launch/recovery/accounting byte-identical,
and fresh parent resolution over the native records (preview → explicit
`Reconcile` → idempotent replay → retained resolution); admission never writes
(never-refused ticket case); a six-variant contradiction table (stale, decoy,
unsafe ID, recovery rebound, evidence rewritten canonically, missing recovery)
fails closed with the journal byte-unchanged; a context admitted before a
same-UID recovery rewrite still fails at the launch boundary (no process, no
started record, no charge, pinned directory consumed fail-closed); concurrent
duplicate plus sequential replay execute exactly once with exactly one launch
reservation in the caller's budget store.

Compile and assertions are NOT verified here — the two earlier rounds of
coordinator-caught compile errors are exactly why every call signature was
re-read against source, but failures remain plausible and must be fixed before
integration. The runner package did not link in the last coordinator check
(disk exhaustion); re-running it is a precondition for any claim about these
tests.

## Limits and needed extensions (NOT done, not designed away)

1. **Attended app CLI wiring** for the recovery relaunch is unimplemented; the
   API contract above is what it must call.
2. **No exported reader for the bound intended invocation.** The app must
   retain the ID it chose at recovery time or read `recovery.json` itself; if
   re-deriving it from the journal is wanted, an exported read-only accessor
   (e.g. `trajectory.ReadCapturedVerificationRecovery`) is the specific needed
   extension — recorded here rather than evading the slice boundary
   (authority readers and the recovery schema are out of scope).
3. **A failed bound relaunch consumes the pinned invocation directory**
   (requested+terminal failed-start records, no started) with no second
   chance: exclusive allocation is fail-closed by design, and
   `RecoverCapturedVerificationLaunch` cannot re-recover it (it admits only the
   original refusal pinned in `launch.json`). Deliberate, documented.
4. **Same-UID forgery** of a fully self-consistent recovery+evidence+lifecycle
   set remains beyond byte evidence — unchanged honesty boundary; the runner
   probe trusts the trajectory authority reader, never `recovery.json` bytes
   directly.
5. The trajectory-level parent fixtures still reconstruct the replacement
   lifecycle bytes (derivation-reader isolation); native production of those
   records is covered at the runner level only.
6. This completes neither the audit slice nor D-series acceptance; no
   experiment was started, and nothing here is published or signed off.
