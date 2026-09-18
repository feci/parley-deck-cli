---
agent: kimi-1
from: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-15
reviewed-commit: 39107b138e77de064f21fec4aafa6073887c2783
review-kind: supporting-source-only
---

# Recovery next slice: refused-launch ticket recovery and mixed-handle coexistence (kimi-1)

## Scope and provenance

- Launch: kimi-1, protocol attestation `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- Method: source reads only at frozen HEAD 39107b138e77de064f21fec4aafa6073887c2783
  (verified via `git rev-parse HEAD`); read-only `ls`/`find` for orientation. No
  builds, tests, overlays, Git mutations, subagents, or network. Line numbers refer
  to this tree.
- PRIMARY (source): every load-bearing claim below carries file:line and quotation
  from this checkout or from the retained candidate artifacts under
  `.parley-runtime/input/budget-refused-prototype-20260914/native-artifacts/`
  (cited as `candidate/<file>:<line>`). Candidate test source is not test
  execution; facilitator measurements in the checkpoint are unverified testimony.
- Read in full: FINAL.md, 00-prompt.md, current checkpoint,
  `candidate/{unchanged.go, refused_unchanged_test.go, runner_refused_unchanged_test.go}`,
  `next-compatibility-plan.md`, `budget-refused-compatibility-20260914/verification.json`,
  `internal/trajectory/{verification.go, captured.go, continuation.go, unchanged.go}`,
  `internal/runner/trajectory_verification.go`, `internal/app/trajectory_verify.go:150-336`.
- Read in part: `internal/trajectory/{state.go, reconcile.go, reservation_recovery.go}`,
  `internal/runner/telemetry.go:60-150`, IMPLEMENTATION.md:3067-4103 (tail plus
  targeted grep), my prior note `kimi-1-abnormal-recovery-review-20260914.md`.
- This is a bounded supporting-source plan: no signoff, no acceptance, no finding
  withdrawal, no code edits, no execution verdicts.

## Three essential conclusions

1. **Same-charge coexistence is impossible; stale-ticket coexistence is reachable
   with a precise legal ordering; a live in-flight Run handle is impossible.** A
   verification ticket can never exist for a v2-refused charge (proof below), so
   the only real old handle that can coexist with a v2 refusal resolution is a
   stale-but-valid `StopCapturedVerification` ticket from an earlier, already
   resolved sequence — and that stop is safe by construction.
2. **The real recovery gap is the verifier-launch budget refusal: it permanently
   bricks the trajectory today.** `beginLaunch` reserves the ticket launch before
   the budget check, so a refused verifier launch leaves `request.json` +
   `launch.json` pinned to a never-started invocation; every existing API then
   refuses, and the charge gate blocks all further fixup. This is recoverable
   without new custody or inactivity claims, using the candidate's own v2
   lifecycle predicate as the admission test.
3. **The slice is implementable as one new recovery artifact plus three small
   amendments in `internal/trajectory/verification.go`, four fixtures, and an
   attended-only app wiring point.** No true design blocker found; one dependency
   (F3 needs the v2 candidate integrated or overlaid) and one open reviewer
   question (O_EXCL vs preview/expected-SHA replay shape) are named below.

## Q1 — Coexistence resolution

**Impossible: a ticket handle for the v2-refused charge itself.**
`capturedRequestAt` refuses materially unchanged attempts:
`captured.go:129-131` `if a.Before.Tree.SHA256 == a.After.Tree.SHA256 { return
CapturedRequest{}, errors.New("unchanged source is not a new patch regression") }`.
The v2 class requires exactly that equality plus a pre-start lifecycle:
`candidate/unchanged.go:49-51` (`"changed material source requires independent
patch verification"`) and `:160-161` (`o.Status == "failed" && *o.FailureClass ==
"budget_refused" && terminal.StartedAt == nil && terminal.PID == nil &&
o.ExitCode == nil`). Both structs pin the same immutable attempt bytes
(`AttemptSHA256`, `captured.go:70`, `candidate/unchanged.go:25`); one attempt
cannot satisfy both predicates. Any fixture proposing a ticket on the refused
charge is not constructible through public APIs — this is the source-grounded
impossibility proof, not a parser-only observation.

**Impossible: a live in-flight Run handle at v2 publication.** The refused
attempt's own run never started (v2 predicate). An earlier attempt's run has
terminated — `validateTransitions` requires `r.At.Before(...)` to fail, i.e. a
resolution postdates its attempt's terminal (`continuation.go:238`). A later
fixup charge is refused while anything is unreconciled or inconclusive-pending:
`state.go:404-406` (`"trajectory awaits independent reconciliation or an attended
review decision; further fixup is refused"`), and an ack covers only the latest
pending sequence (`continuation.go:200-206`). A live verifier `RunConsult` for
sequence N keeps N unreconciled, which gates N+1 (same lines); a v2 resolution
for N is impossible by the first proof. No legal interleaving remains.

**Reachable: a stale, consumed ticket handle plus a later v2 resolution.**
Legal event ordering, all steps exercised by existing or candidate fixtures:

1. E1: attempt N charged, implementer run changes source, terminal retained.
2. E2: `PrepareCapturedVerification(N)` writes `request.json`
   (`verification.go:300-337`; journal dir keyed by `r.Charge.EntryKey`,
   `verification.go:322`).
3. E3: helper executes; `receipt.json` retained (`verification.go:424-528`).
4. E4: resolution N published through the verifier-run path (non-unchanged
   resolutions require `runtimeID(p.RunID)` and `validHash(p.ParentSHA256)`,
   `continuation.go:245`).
5. E5–E6: attended ack if needed; charge N+1 reserved (gate now passes).
6. E7: N+1's launch is budget-refused pre-start (production ordering:
   `telemetry.go:96-100` reserves before budget; candidate runner fixture
   asserts the exact refused record, `candidate/runner_refused_unchanged_test.go:61-63`).
7. E8: `ReconcileUnchanged(N+1)` publishes the v2 resolution
   (`candidate/unchanged.go:237-268`).
8. E9: the retained old handle calls `StopCapturedVerification(T(N), inv)`. It
   re-derives the request at the pinned sequence (`verification.go:362`
   `capturedRequestAt(s, ticket.Request.Verifier, ticket.Request.Sequence)`), so
   later attempts/resolutions do not invalidate it; it writes only `stop.json`
   in T(N)'s journal and kills only strictly attributed live processes
   (`verification.go:134-152`). Trajectory state, ledger, and the v2 resolution
   bytes are untouched. Stop legality survives the v2 publication — fixture F3.

**Old binary after a v2 publication:** any old-binary control (stop included)
fails at the whole-state read — `readState` uses `DisallowUnknownFields`
(`state.go:249`) and v1-only `UnchangedEvidence` has no `class` — before any
write or kill. The compatibility harness proved this for inspect/history/apply
(`verification.json`); the stop route needs fixture F4. This is a fail-closed
state-read refusal, stronger than a per-artifact parser refusal, and F4 must
prove the no-write property, not just the error.

## Q2 — The real gap and the recovery design

**MAJOR-1 — Verifier-launch budget refusal permanently bricks the trajectory.**
Ordering proof: `telemetry.go:96-99` calls `reserveCapturedVerification` (writes
`launch.json`, `verification.go:398-406`) *before* `reserveBudget`. On denial the
terminal is `budget_refused` with no `started.json`. Thereafter: retry of
`PrepareCapturedVerification` fails — `verification.go:283-285` `"verification
already reserved or unavailable; preserve it for explicit recovery"`; `Execute`
with a fresh invocation fails — `verification.go:414` `launch.InvocationID !=
invocation`; the charge gate blocks all fixup (`state.go:404-406`); `Continue`
cannot ack an unresolved attempt (`continuation.go:197-206` needs a resolution).
No recovery API exists. Consistent with FINAL D6 ("exhaustion escalates with
durable evidence and never closes"), the brick escalates — but nothing today
lets an operator *repair* it, so the escalation is terminal for the idea.

**Recovery design (no custody, no inactivity inference).** Admission is the
candidate's own v2 predicate applied to the refused *invocation* (not the
attempt): terminal `failed`/`budget_refused`, `StartedAt`/`PID`/`ExitCode` all
nil, `started.json` absent, chronology valid — plus `claim.json`/`receipt.json`/
`stop.json` absent. A recovered launch is not a retry of verification work: no
claim, step, or receipt can exist, and lineage stays enforced because
`executeTrajectoryHelper` requires the inherited process markers
(`app/trajectory_verify.go:327-330`).

**MAJOR-2 — Interrupted Prepare leaves the same brick, but has a safe quarantine
class.** If the process dies between the reservation `Mkdir` and the
`request.json` write, the same `:283-285` error bricks the charge. Safety proof
for quarantine: `PrepareCapturedVerification` returns a zero ticket on any error
(`verification.go:333-336`) and returns the ticket only after the durable
`request.json` write and a final SHA256 pass (`:327-332`). Therefore a
reservation lacking a canonical `request.json` provably has **no live ticket
handle anywhere** — quarantine cannot orphan a handle. This is an existence
proof, not a process-inactivity assumption.

## Q3 — Minimal source/test changes (design; no code written)

**S1 — `internal/trajectory/verification.go` (only production file with logic
changes):**
- New `RecoverCapturedVerificationLaunch(ctx, ticket, refusedInvocation,
  newInvocation) error`: runs under `withVerification` (full authority, same as
  `ReserveCapturedVerificationLaunch`, `verification.go:398-406`); requires
  `launch.json` pinned to `refusedInvocation`; requires absent
  `claim.json`/`receipt.json`/`stop.json`; reads
  `.parley-runtime/invocations/<refused>/terminal.json` under `ticket.Root` via
  `os.Root` and enforces the v2 predicate verbatim; writes `recovery.json`
  (O_EXCL, bounded, canonical — same discipline as `writeVerificationArtifact`,
  `verification.go:194-215`) binding `{ticketSHA, priorLaunchSHA,
  refusedTerminalSHA, newInvocation, at}`.
- Amend `readVerificationLaunch` (`verification.go:408-418`): resolve the
  *effective* launch — `launch.json`, else `recovery.json` when the invocation
  matches; return the effective artifact SHA so claim/receipt binding
  (`verification.go:647`) works unchanged.
- Amend `StopCapturedVerification` (`verification.go:106-155`): resolve the
  effective launch first; if the caller's invocation is superseded, return a
  dedicated error **before** writing `stop.json`. Named ordering rule:
  stop-before-recovery blocks recovery (`stop.json` present → refuse);
  recovery-before-stop supersedes the stale stop. Both directions fail closed.
- Amend the journal inventory (`verification.go:620-634`): admit `recovery.json`
  and validate its binding. Old readers reject such journals at `:632-634`
  (`"unexpected or out-of-order artifacts"`) — fail closed, asserted in F4.

**S2 — abnormal-reservation quarantine (same file or small new
`verification_recovery.go`):** explicit, attended-only
`QuarantineCapturedVerificationReservation`: refuse if `request.json` is
canonical; otherwise write `quarantine.json` (observed inventory, charge key,
timestamp), fsync, rename the reservation dir to a suffixed sibling (history
preserved, never deleted — honors `verification.go:281-283`), fsync parent.
Fresh `Prepare` then reserves cleanly.

**S3 — attended app wiring:** at `app/trajectory_verify.go:217-219`, surface the
"reserved but unrecoverable inline" error as a distinct class; recovery is
invoked only through an attended control (same attended-gate pattern as
`Continue`, per my prior note's citations of `trajectory_reconcile.go:70-73`).
No auto-retry anywhere.

**Fixtures and exact assertions:**
- **F1 (API, new `internal/trajectory/verification_recovery_test.go`)** —
  priority P1. Happy path: prepare ticket, `ReserveCapturedVerificationLaunch`,
  fabricate the refused invocation with the candidate's telemetry helper
  (`candidate/refused_unchanged_test.go:15-43`), recover with a new invocation,
  execute and read the receipt. Assertions: budget ledger snapshot byte-equal
  across recovery (AC-B1/D6); `RequireResolved` still fails before execution;
  stale stop with the refused invocation → superseded error and **no**
  `stop.json`; stop with the effective invocation succeeds post-execution;
  regression-streak semantics untouched. Negative table (mirrors
  `candidate/refused_unchanged_test.go:111`): started.json present; PID set;
  ExitCode set; wrong failure class; claim present; receipt present; stop
  present; second recovery; wrong refused id; execute with superseded
  invocation — each must refuse at its named assertion. Predicate-removal
  controls per IMPLEMENTATION.md:3089-3091 convention.
- **F2 (runner integration, `internal/runner/`)** — P2. Real
  `WithLaunchBudget`-denied launch on `CapturedVerificationPhase` (extends the
  candidate's runner fixture shape), then recovery and a real allowed run.
  Assert journal contents, terminal shape, and no spawn marker.
- **F3 (coexistence)** — P3; native overlay until the v2 candidate integrates.
  Ordering E1–E9; assert state and ledger byte-equality across E9, `stop.json`
  only in N's journal, v2 preview replay byte-stable, `RequireResolved` fails.
- **F4 (old-binary compatibility)** — P3; extend the
  `budget-refused-compatibility-20260914` harness with stale-ticket-stop and
  recovery-journal routes on the frozen v1 CLI; assert refusal and journal
  hash-set equality (no-write), not merely an error string.

## Findings

- **MAJOR-1** (above): verifier-launch refusal brick. Locators: `telemetry.go:96-100`,
  `verification.go:283-285`, `verification.go:408-418`, `state.go:404-406`,
  `continuation.go:197-206`.
- **MAJOR-2** (above): abnormal reservation brick; safe quarantine class exists.
  Locators: `verification.go:281-285`, `:327-336`.
- **MINOR-1**: naive recovery would let a stale stop durably kill a recovered
  launch (`stop.json` written before the claim check, `verification.go:115-133`);
  S1's effective-launch resolution and superseded-stop refusal are load-bearing.
- **MINOR-2** (carried, already recorded in `next-compatibility-plan.md:27-28`):
  candidate `UnchangedEvidence` comment still says "actually exited process"
  (production `unchanged.go:16-17`); qualify on integration.
- **NIT-1**: `docs/agent-runtime-configuration.md:1298-1337` documents the ticket
  lifecycle; add the recovery ordering rule when S1 integrates.
- **CRITICAL**: none found in this slice's scope.

## Bounded counterexamples

C1–C6 are the F1/F4 negative variants above. The strongest: C2 (claim exists →
recovery refuses; the attempt escalates with partial journal retained,
`verification.go:644-651`) and C6 (old-binary stop after v2 publication →
whole-state read failure, journal hash-set unchanged).

## Open question and blocker statement

- Open reviewer question (not a blocker): should recovery be single-shot O_EXCL
  with a deterministic error (recommended; matches
  `verification.go:192-193` "a torn/noncanonical artifact is an unresolved
  failure, never a retry") or preview+expected-SHA idempotent like
  `ReconcileUnchanged` (`candidate/unchanged.go:237-268`)? Recommend O_EXCL;
  request reviewer disposition.
- Dependencies (not blockers): F3 needs the v2 candidate integrated or run as a
  native overlay; F4 needs the compiled-old-CLI harness from
  `budget-refused-compatibility-20260914`.
- **True design blocker: none.** Every admission test reuses existing predicates
  and authorities; no custody, inactivity, or seamless mixed-version claim is
  required, and no FINAL.md acceptance criterion is weakened — recovery spends
  nothing, resets nothing, and closes nothing (D3/D6, AC-B1).

## Next actions (prioritized)

1. P1: implement S1 + F1 natively; run the negative table and predicate-removal
   controls; retain failures.
2. P2: F2 runner-level integration; then S3 attended wiring as a separate
   reviewed change.
3. P3: F3 (with the v2 candidate) and F4 (old-binary harness extension).
4. P4: S2 quarantine with C3 control; update the doc section (NIT-1) at
   integration time.
