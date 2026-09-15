---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-15
reviewed-commit: 39107b138e77de064f21fec4aafa6073887c2783
review-kind: supporting-source-only
---

# N1/N2/N3 dispositions at 39107b1 and remaining hardening defects

## Summary

- **N1 is narrowed.** A known refusal now happens before the driver charges a fixup and before the verifier CLI creates its ticket. Budget spent before launch on cross-review rounds and driver steps is still open.
- **N2 is narrowed.** The shared validator releases the guard while it reads history. By reading, its retry is bounded, uses fresh evidence and runs its callback at most once. One route, `recover-reservation`, still holds the guard across a full-history read (R1).
- **N3 is resolved** for mismatches that exist when the check runs. A narrower window remains (R3).
- **`TestAgentsExecRecordsManualLaunch`.** Source identifies which code path failed, but not why. The test throws away the record that would have classified the failure (R2).

Findings: 0 CRITICAL, 0 MAJOR, 2 MINOR, 1 NIT.

This result covers source only, within the scope below. It is not whole-audit acceptance, not a full Phase-6 review, and not a signoff.

## Scope and provenance

- **Launch attestation:** `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- **Method:** native Read, Grep and Glob over the frozen tree only.
  - No shell, Git, test, build, overlay or model call.
  - I did not open other participants' outputs from this allocation.
- **PRIMARY (source):** a file:line locator plus a quotation. Tests are source I read, not runs I executed.
- **Unverified testimony:** every timing and every pass, fail and control result in the checkpoint, the three codex-1 notes and IMPLEMENTATION.md:3685-4104.
- **Read in full:**
  - `trajectory/{state,state_validation,helper_scope,verification}.go`
  - the three `state_validation*_test.go` files and `helper_scope_test.go`
  - `runner/protocol_context.go` and all three `protocol_precheck_test.go` files
  - `app/{trajectory_verify,agents_exec}.go` and `app/trajectory_quorum_test.go`
  - `driver/cycle_budget.go`
- **Read in part:**
  - `phase58.go`, `driver/{impl,budget}.go`, `app/driver_impl.go`
  - `reservation_recovery.go`, `parent_recovery.go`, `snapshot.go`, `reconcile.go`, `captured.go`, `activation_quorum.go`
  - `launch.go`, `runner.go:1224-1236`
  - `cycle_binding.go`, `cycle_extension.go`
  - `agents_exec_test.go`, `app_test.go:1587-1609`

## Dispositions

### N1 [MAJOR]: narrowed; the budget-spend residual stays open

**Held by reading (PRIMARY).**

- **Driver.**
  - `impl.go:357` `if err := d.cfg.Impl.PrecheckFixup(ctx); err != nil {` escalates before `:365`
    `d.reserveFixupCycle(ctx, charged)`. The charges it would make are `ChargeStep` and `ChargeCycle` at `cycle_budget.go:84,87`.
  - `stepImplOps` (`driver/budget.go:79`) defines no `PrecheckFixup`. The promoted method therefore charges no step.
- **The precheck sees the same inputs as the launch.**
  - `driver_impl.go:470` `runner.PrecheckFixup(ctx, o.withParticipants(o.implementer))` mirrors `:475`.
  - `phase58.go:56,60-61` selects the same agent with `Phase: "fixup"`, exactly as `RunFixup` does (`:78,97-98`).
  - `prepareProtocolPrompt` reads only `info.Phase` and `info.Idea` (`protocol_context.go:38-73`). Its refusal checks inspect the render result, not the prompt (`:77-94`).
- **Verifier CLI.** The precheck at `trajectory_verify.go:212` runs before `:217`. That call consumes the ticket at `verification.go:283`
  `if err = base.Mkdir(charge, 0700); err != nil {`.
- **Test source (not runs).**
  - Driver: no fixup call, no step entry, and an unstarted terminal (`driver/protocol_precheck_test.go:101,127,142`). Restoration: `:151-169`.
  - App: no journal and no request (`app/protocol_precheck_test.go:79-83`). Restoration: `:110-117`.
  - A late refusal still leaves `Launch == nil` (`runner/protocol_precheck_test.go:66`).

**Still open.**

- `driver/cycle_budget.go:66` `if _, err := budget.ChargeCycle(ctx, budget.CrossReview); err != nil {` still runs before `RunRound` (`:69`).
- The step wrappers also charge before they launch (`driver/budget.go:46,82,106`).
- A known refusal on those paths spends budget but strands no one-time authority.
- The window between the precheck and the real render still includes the guarded `AfterCycle` history checks (`state.go:461`) and `Begin`. This window is already disclosed.

### N2 [MINOR]: narrowed to R1

**Mechanism.**

1. Structural snapshot taken under the guard: `state_validation.go:22-26`.
2. Content reads run without the guard: `:36`, `:39`.
3. The guard is reacquired and equality is checked at `:45`:
   `... !sameJSON(b.Policy, originalBinding.Policy) || !sameJSON(ledger, originalLedger) || !sameJSON(s, originalState) {`
4. Only then does the callback run, still under the guard: `:49` `return fn(b, ledger, s)`.

Separately, the pending-attempt refusal at `state.go:404` now runs before `checkStateSnapshots` (`:407`).

**Refutation attempts (all held by reading).**

- **Callback runs at most once.** `changed` is set only on the inequality branch (`:46-47`). A callback error is therefore returned, not retried. Test source: `state_validation_retry_test.go:62-64`.
- **Retry is bounded.** `:51` `if changed && attempt == 0 {`. A second drift returns an error. Test source: `:50-53`.
- **Evidence is fresh.** Every iteration rereads all evidence (test source `:54-56`). The callback receives the freshly read, equal state, so no update is lost.
- **Retried checks have no side effects.** `parentRecoveryResolutionCheck` only reads (`parent_recovery.go:262-293`). Its publication stays in the guarded callback (`:386`). A concurrent second apply takes the replay branch for the retained file (`:355-373`).
- **No archive swap under an unguarded reader.** `CaptureSnapshot` never replaces an archive that already exists:
  `snapshot.go:253-257` `if _, err = os.Lstat(final); err == nil { if err = CheckSnapshot(...`.
- **The guard is scoped per idea and kind** (`cycle_binding.go:95,113-116`).
- **Positive `BeforeCycle` and `AfterCycle` still read history under the guard.** They do so only after `state.go:404` has proved that no attempt is unreconciled. No live terminal or ticket can therefore be waiting on them.

**Coverage, and what I did not read.**

- The public schedule of two overlapping applies refuted the prototype without a retry. That schedule exists only as runtime instrumentation (IMPLEMENTATION.md:4012-4013).
- Canonical tests cover drift in state, policy and ledger (`state_validation_test.go:83-92,146-148`).
- I did not read `budget.Snapshot`. The equality check assumes it has no fields computed at read time.

### N3 [MINOR]: resolved for mismatches present at check time

**Mechanism.**

- The helper request must equal the charged attempt's original request: `helper_scope.go:31` `if err != nil || !sameJSON(r, req.Ticket.Request) {`.
- Membership must equal the activation quorum: `:38` `if !slices.Equal(req.Participants, members) {`.
- The live scope must match: `:46` `checkReconciliationScope(dir, req)` requires equal live membership and exact criteria (`reconcile.go:129,133`).
- The app calls this check at three points:
  - `trajectory_verify.go:182`, before the runtime directory, the precheck and the ticket;
  - `:261`, before acceptance;
  - `:331`, at helper entry.

**Refutation attempts.**

- **YAML trailing comment.** The mis-parsed `parseList` element enters `req.Participants` (`:173,181`). The check refuses at `:182`. Test source: `trajectory_quorum_test.go:23-24,55-63`.
- **Parser disagreement.** The `parseList` result, the strict live YAML parse and the activation archive must all agree pairwise.
- **Changed original invocation.** Refused (`helper_scope_test.go:31`).
- **Helper replay.** Prevented by the exclusive claim (`verification.go:436-438`).

## Findings

### [MINOR] R1: `recover-reservation` holds the cycle guard across a full-history read of a live charged row

**Evidence (PRIMARY).**

- **Reachable from the CLI.** `app/trajectory.go:17-18` → `trajectory_reservation_recovery.go:29`
  `p, err = trajectory.PreviewReservationRecovery(ctx, *root, *idea, *entry)`.
- **Takes the guard directly.** The shared body acquires it at `reservation_recovery.go:296`
  `release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir))` and holds it until return (`:300`).
- **Nothing refuses a row that has since launched.**
  - `validateIntentBefore` checks policy, the ledger prefix and the gates of the intent's own `Before` state (`:150-174`).
  - `compareIntentPrefix` compares only the first `len(i.Before.Attempts)` attempts (`:196-205`).
  - An existing row is compared only by `:361`:
    `if old.ReservationIntentSHA256 != sha || !sameJSON(old.Charge, a.Charge) || old.Before != a.Before || old.BeforeArchive != a.BeforeArchive {`
- **Then reads full history, still under the guard:** `:368` `if err = checkStateSnapshots(ctx, *b, s); err != nil {`.

**Counterexample (by reading; magnitude not measured).**

1. Driver fixup attempt k is live and has its intent (`state.go:443-449` requires one).
2. An operator previews `recover-reservation` for k's entry key.
3. The preview checks every archive and every resolution while holding the guard.
4. `Run.Finish` waits in `withStateControl` for at most 30 s (`state.go:490-492`). It then fails without recording the live terminal. This is the consequence I reported in N2.
5. A stop for a live ticket contends in the same way.

The limits in the N2 note name only "Positive cycle reservation and individual callbacks" (`codex-1-state-validation-20260914.md:63-64`). This route is neither.

**Minimal fix.**

- In `internal/trajectory/reservation_recovery.go`:
  1. keep the intent and structural checks under the guard;
  2. run `checkStateSnapshots` without it;
  3. reacquire the guard and require equal policy, ledger and state before returning, or before writing the missing row (`:380-381`).
- Add a test beside `state_validation_test.go:17-98` that pauses a preview and calls `Finish`.
- Update the N2 note and the docs.
- At minimum, disclose this route.

### [MINOR] R2: The manual-launch test throws away the record that would classify its failure

**Evidence (PRIMARY).**

- **Which branch failed.** The generic message at `agents_exec.go:130` is
  `"agents exec: invocation failed; inspect private local evidence"`.
  - It needs a recorded invocation (`:116-119`) and `launchErr != nil` (`:129`).
  - It prints after the content-free `--json` record has already gone to stdout (`:120-124`).
  - It is not the artifact branch (`:133-135`).
- **Where `launchErr` can come from.** `RunMeasured` can fail in three places:
  - setup (`launch.go:206-208`);
  - opening a private log (`:210-215`);
  - `cmd.Run()` (`:224`), under `context.WithTimeout(parent, timeoutForAgent(...))` (`:191`). With no override, `timeoutForAgent` uses `agent.TimeoutMS` (`runner.go:1225-1233`), and the fixture writes `timeout_ms = 5000` (`app_test.go:1602`).
- **The cause is not identifiable from source.** Source cannot choose among these three, and I assign no cause.
- **Why the evidence was lost.** On failure the test prints only `t.Fatalf("code=%d stderr=%s", code, errOut.String())` (`agents_exec_test.go:44`) and drops `out`. The fixture root is `t.TempDir()` (`:18`). This explains why the evidence is gone, not why the test failed.
- **Timing does not help.** If the testimony is accurate (a 6.638 s case against a 371 ms invocation), elapsed test time cannot show whether a deadline was hit.

**Fix.** Include `out.String()` in that `Fatalf`; the test at `:56-58` already asserts that this output carries no content. Leave the original failure unresolved until a failing record is captured. File: `internal/app/agents_exec_test.go`.

### [NIT] R3: Helper scope is checked before the long precheck and history read, not where the ticket is consumed

**Evidence (PRIMARY).**

- The order in `trajectory_verify.go` is:
  1. the scope check at `:182`;
  2. the protocol render at `:212`;
  3. `PrepareCapturedVerification`'s full `withState` read (`verification.go:310`);
  4. the ticket `Mkdir` (`:283`).
- `PrepareCapturedVerification` takes no scope (`:300`). An edit to the live `participants:` or `checks:` inside that window still consumes the ticket. Refusal then arrives only at `trajectory_verify.go:331` or `:261`.

**Fix (narrows the window, does not close it).** Pass the `HelperRequest` into preparation. Run `checkReconciliationScope` inside its guarded callback, before `openVerificationDirectory(..., true)` (`verification.go:322`). Files: `verification.go`, `trajectory_verify.go`, `trajectory_quorum_test.go`.

## Next actions

1. **R1:** fix it, with a test that pauses a `recover-reservation` preview while `Finish` or stop runs. Otherwise, disclose the route.
2. **R2:** make the test print the terminal record on failure.
3. **N1 residual:** keep the budget spent before cross-review and step launches open.
4. **R3:** optional.

## Remaining uncertainty

- **Candidate, not a finding: `Begin` after a charge.** The one-time retry could be exhausted after the charge.
  - Trigger: two policy extensions land during `Begin`'s unguarded read (`state.go:569`). `cycle_extension.go:218-235` shows no refusal for a pending charge.
  - Not established: I did not read `CheckPolicy` (`reservation_recovery.go:151`), which might refuse an extended policy on its own.
- **Not read:** `budget.Snapshot`, `trackedCommandFor`, `trajectoryHistory`.
- **Premises from recall (UNVERIFIED):**
  - `t.TempDir` directories are removed at cleanup.
  - TOML `timeout_ms` maps to `TimeoutMS`.
- **Not established:**
  - a fence against same-UID writers;
  - descendant inactivity;
  - any general speedup.
