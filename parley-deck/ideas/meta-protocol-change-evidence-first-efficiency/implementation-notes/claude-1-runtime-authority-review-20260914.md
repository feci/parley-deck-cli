---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
reviewed-commit: 6f31201c987ecb085921b537d0345e134380ec05
review-kind: supporting-source-only
---

# Runtime-authority boundaries: supporting source review

## Summary

I read the two new runtime-authority routes at 6f31201 and the code they depend on:

- **Terminal publication.** `Begin` freezes the canonical post-launch state, and `Finish` publishes
  through `withStateControl`.
- **Existing-ticket stop.** `StopCapturedVerification` reads state through
  `withVerificationAuthority(..., false, ...)`.
- **Prerequisite responses** to my earlier findings: F1 (known protocol refusal before fresh
  reservations), F3 (activation-archive membership) and F6 (publication coupling), plus bounded
  archive stability.

Within the disclosed same-UID trust boundary I could not break either control route by reading:

- `Finish` refuses changed canonical state, a substituted intent reference, a changed charge or policy,
  and terminal replay. It grants nothing: execution, reconciliation, continuation and completion still
  run full source/result validation.
- `StopCapturedVerification` refuses a rewritten ledger charge, a different ticket/request/launch
  binding and a changed process identity. It never releases an execution callback and signals only
  attributed processes.
- No legitimate guarded writer I traced can change trajectory state between `Begin` and `Finish`, so the
  frozen-state binding adds no false refusal among them.
- Verification membership is pinned to the hash-verified activation archive. My F3 counterexample now
  refuses before a ticket exists.
- A timestamp-only archive transition gets exactly one strict reread, pinned to the original file.

Findings: no CRITICAL, one MAJOR, two MINOR, three NIT.

- **N1 (MAJOR).** The known-refusal prevention does not reach the two production callers that spend
  one-time authority before the launch boundary:
  - driver-managed fixups charge the cycle and step first;
  - `trajectory verify` creates the charge's single ticket directory first.

  A known refusal on either path still leaves a spent trajectory attempt that can never be resolved.
- **N2 (MINOR).** `Finish` (≤30 s) and the runner's stop (≤20 s) still wait for the shared cycle
  guard. Full-history readers hold that guard for their whole validation: `trajectory inspect`, history
  and previews, and a refused concurrent fixup's `BeforeCycle`. Terminal and stop loss therefore remain
  coupled to history size through contention.
- **N3 (MINOR).** A membership difference between the live prompt and the activation archive is
  detected only at reconciliation, after the one-per-charge ticket and the verifier launch are consumed.
  A `parseList`/YAML divergence is one deterministic trigger.
- **N4 (NIT).** The "changed-original-ledger" terminal case rewrites state, not the ledger, and refuses
  before `Finish`'s own checks run.
- **N5 (NIT).** `Finish` can re-publish a deleted content-addressed historical archive from the live
  tree. The bytes are hash-identical, but "no old source is reconstructed" overstates the behaviour, and
  the retention test cannot observe it.
- **N6 (NIT).** A publication deadline can be persisted as `source-unavailable`.

I also give bounded feedback on the unimplemented abnormal-recovery proposal.

This is not a full-scope Phase-6 review, a signoff, a withdrawal of any earlier finding, or acceptance.
My 2ab9e82 review is unchanged. I executed nothing and issue no execution verdict.

## Scope and provenance

- **Launch.** claude-1, protocol attestation `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- **Method.** Native Read/Glob/Grep over the frozen tree only. No shell, Git, build, test, race run,
  mutant or overlay. Line numbers refer to this tree.
- **Provenance (§15.2).**
  - *PRIMARY (source):* a file, line and quotation I read. Tests are cited as source I read, never as
    tests I ran.
  - *RECALL:* Go, YAML, Git or OS runtime semantics stated from memory. These are UNVERIFIED.
  - *Testimony I did not verify:*
    - the 418-file manifest `2d4316bf1d6e43873bd0a508529081fb983ea4f14e6fa521cf2b9a12f094ed2f`;
    - the full-suite, six-package race, vet, Windows cross-build and compiled shared results;
    - the four terminal cases and three negatives;
    - the cleanup checkpoint's four tests, four negatives and disk-exhaustion retry;
    - the F6 envelope timings: 30.004 s original, 5.193 s current, 197.013 s full read;
    - the HTTP 429 result of my allocated activation-archive review (IMPLEMENTATION.md:3337-3349,
      3380-3389). That attempt produced no artifact, and nothing here relies on it.
- **Read in full:**
  - `internal/trajectory/{state.go, verification.go, verification_cleanup_test.go, terminal_authority_test.go, activation_quorum.go, unchanged.go, captured.go, reconcile.go, continuation.go, reservation_recovery.go, parent_recovery.go, archive_stability.go}`
  - `internal/runner/{protocol_context.go, telemetry.go, launch_budget.go, protocol_refusal_precharge_test.go, trajectory_verification.go, cycle_budget.go, durablekill.go}`
  - `internal/app/{trajectory_verify.go, trajectory.go}`, `internal/procctl/procctl.go`,
    `internal/budget/cycle_observer.go`, `internal/driver/cycle_budget.go`
  - FINAL.md; IMPLEMENTATION.md lines 2973-3613; my 2ab9e82 review; the codex-1 notes on:
    - terminal authority and cleanup authority;
    - the abnormal-recovery proposal;
    - activation quorum and archive stability;
    - protocol refusal and the history envelope.
- **Read in part:**
  - `internal/trajectory/snapshot.go` (1-300, 495-760), `activation_quorum_test.go` (1-300),
    `unchanged_test.go` (23-84), `state_test.go` (48-76)
  - `internal/runner/runner.go` (1040-1149), `internal/runner/phase58.go` (50-122)
  - `internal/budget/cycle_binding.go` (225-287), `cycle_session.go` (37-122),
    `reservation_receipt.go` (34-43), `resource_guard.go` (21-38)
  - `internal/protocol/workspace.go` (225-417), `internal/driver/impl.go` (355-391)
  - `internal/app/driver_impl.go` (440-499), `internal/app/trajectory_verify_test.go` (255-276)
  - `internal/telemetry/record.go` (`Record`), `docs/agent-runtime-configuration.md` (1150-1229)
- **Not claimed:**
  - current-tree acceptance, Windows runtime or real-model launches;
  - ACP/interactive spawn ordering;
  - escaped-descendant inactivity or general helper recovery;
  - packet/pilot experiments, quorum decisions or signatures.

## Refutation attempts

### Terminal publication (`Begin` / `Finish`)

**RT-1: Changed original state after launch. Held.**

- `Begin` hashes the exact bytes it writes:
  - `state.go:604` `frozen, err := canonical(s)`
  - `state.go:608` `if err = writeState(statePath(b), s); err != nil {`
  - `state.go:611` `run = &Run{root, idea, charge, invocation, expected, digest(frozen)}`
- `readState` accepts only canonical bytes (`state.go:260` `if !bytes.Equal(data, expected) {`). The
  file digest therefore equals the frozen digest.
- `Finish` refuses any difference:
  `state.go:635` `if err != nil || r.stateSHA256 == "" || digest(frozen) != r.stateSHA256 {`.
- The test rewrites `ReservationIntentSHA256` to a valid-format hash
  (`terminal_authority_test.go:85-89`). `validateState` alone would accept that value (`state.go:280`).
  The test asserts the specific refusal and unchanged bytes (`:99-107`).

**RT-2: Old intent reference or intent content substitution. Held.**

- The reference is inside the frozen attempt (RT-1).
- Intent *content* is deliberately not read on this route. `checkSourceSnapshots`, and with it
  `checkReservationIntents`, is skipped when `requireSourceEvidence` is false (`state.go:483, 525-532`).
- A substituted intent file therefore cannot block a live terminal. Every later full read catches it:
  `reservation_recovery.go:245-246`
  `if sha != a.ReservationIntentSHA256 || ... { return errors.New("charged trajectory changed its original reservation intent")`.
- Publication grants nothing.

**RT-3: Charge and ledger changes. Held.** The test caveat is N4.

- Every attempt's charge is checked against the live ledger before the closure runs:
  `state.go:308` `if err := a.Charge.Check(ledger); err != nil {`, reached via `state.go:522`.
- The predicate compares scope, epoch, entry, kind, `ReservedAt`, reserve amount and action identity
  (`reservation_receipt.go:36-40`).
- An added or removed ledger entry refuses at `state.go:270` (`len(s.Attempts) != len(ledger.Entries)`).
- The receipt-side `r.charge.Check(ledger)` (`state.go:631`) is an additional check.

**RT-4: Policy change. Held.**

- `state.go:629` `if b.Policy.TrajectorySHA256 != r.policySHA256 {`
- `state.go:270` `sha != b.Policy.TrajectorySHA256`

**RT-5: Missing older archives or results. Held.**

- `withStateControl` passes `false, nil` (`state.go:482-484`).
- `validateState` still runs `validateTransitions`. Its unchanged-resolution check is structural
  (`unchanged.go:54-75`) and reads no telemetry or archive.
- The test cases remove the baseline archive (`terminal_authority_test.go:77-80`) and an earlier
  `terminal.json` (`:81-84`).
- Both cases then assert that `Inspect` and `RequireResolved` still refuse (`:132-137`). For the archive
  case, see N5.

**RT-6: Post-source or capture failure. Held.**

- Observation failure records `source-unavailable` (`state.go:654-656`).
- Archive failure records `archive-unavailable` (`state.go:652`).
- `validateState` accepts only those two shapes (`state.go:322-331`).
- The state is written, then the error is returned (`state.go:663-668`). The runner marks
  `trajectory_failure` (`telemetry.go:197-199`).
- Nothing is reconstructed. None of the four new cases exercises this branch. See also N6.

**RT-7: Terminal replay. Held.**

- After the first publication the frozen digest differs, so replay refuses at `state.go:635`, before the
  `a.Terminal != nil` check (`state.go:639`).
- The test asserts refusal and byte identity (`terminal_authority_test.go:138-144`).
- The runner calls `Finish` at most once (`telemetry.go:147` `l.once.Do(func() {`).

**RT-8: Concurrent guarded writers between `Begin` and `Finish`. No false refusal found.**

- **A new charge.** `BeforeCycle` refuses while the live attempt is unreconciled (`state.go:406-409`).
- **Resolving the live attempt.** Both routes need a terminal (`unchanged.go:45`; `captured.go:126`).
- **Replaying an earlier resolution.** No state is written (`unchanged.go:231-245` only syncs;
  `reconcile.go:422-426` returns nil).
- **`Continue`.** Refuses (`continuation.go:123-125`).
- **Reservation recovery of an existing row.** Only syncs (`reservation_recovery.go:359-364, 377-389`).
- **Activation replay.** Writes nothing (`state.go:370-376`).
- **Parent recovery.** Writes `parent-recovery.json`, not state (`parent_recovery.go:386`).
- Contention on the guard itself is N2.

**RT-9: Lost live handles. Nothing is reconstructed; an ambiguity remains open.**

- `Run` has only unexported fields (`state.go:565-573`). It is assigned only after `writeState` returns
  nil (`state.go:608-611`). No persisted handle or replay API exists.
- If the runner dies after `Begin`, the attempt keeps `Launch` without `Terminal` and stays blocked
  (`unchanged.go:45`; `captured.go:126`; `state.go:406-409`).
- If `writeState` lands but reports an error, `run` stays nil and the launch aborts before spawn:
  `launch_budget.go:78-80` → `telemetry.go:99-101` → `protocol_context.go:134-137` →
  `runner.go:1051-1053`. The persisted shape is the same, with no process behind it.
- Whether `fsutil.ReplaceSyncedFile` can fail after its rename is RECALL (I did not read fsutil) and is
  UNVERIFIED.
- AP-4 and AP-5 cover the consequence for recovery design.

### Existing-ticket stop control

**RS-1: Changed ledger or charge. Held.**

- The control route still runs `validateState`, including `state.go:308` (`verification.go:356-360`).
- The coordinated state/request/launch rewrite fixture (`verification_cleanup_test.go:193-241`) asserts
  the original-ledger refusal and that no `stop.json` was written.

**RS-2: Changed ticket, request or launch. Held within the trust boundary.**

- The request is recomputed from state (`verification.go:362`
  `r, err := capturedRequestAt(s, ticket.Request.Verifier, ticket.Request.Sequence)`). It must match both
  the caller's ticket (`:374`) and the persisted `request.json` (`:384`).
- The launch binding is checked at `verification.go:111-114` and `408-417`, and the stop binding at
  `:123-125`.
- A coordinated same-UID rewrite can still pass: `validateState` checks the attempt's `After` only
  structurally (`state.go:326, 332`), so a new `After` plus matching `request.json` and `launch.json`
  gets through. That grants only a stop fence and attributed signals, never execution. This is within
  the disclosed boundary.

**RS-3: Missing earlier resolution terminal. Held.**

`verification_cleanup_test.go:109-162` does the following:

1. removes the first invocation's `terminal.json`;
2. confirms `Inspect` refuses (`:136-138`);
3. requires a stop and a byte-identical stop replay (`:139-158`);
4. requires that no helper claim exists (`:159-161`).

**RS-4: Stopped tickets. Held.**

- The first `stop.json` is preserved (`verification.go:116-125`).
- Every later start and journal read refuses it (`verification.go:428, 461, 653-655`).
- The runner stops only failed, cancelled or watchdog-killed verifier launches
  (`runner.go:1111-1117, 1129-1133`). Acceptance separately requires exit 0
  (`trajectory_verify.go:266-272`). A stop therefore cannot discard an acceptable result.

**RS-5: Changed process identity. Held within the stated scope.**

- Signalling is gated by `KillTreeAttributed` → `Attributed` (`procctl.go:175-183, 107-158`).
- The test asserts `process start time mismatch` while the fixture child stays alive
  (`verification_cleanup_test.go:73-84`).
- A dead leader is skipped (`verification.go:142-144`).
- A torn `process-NNN.json` aborts the loop (`verification.go:135-141`). The criterion group then goes
  unsignalled for two reasons:
  - it is its own group, because `verification.go:95` accepts only identities with `PGID == PID`
    (see also `runner.go:1112-1113`);
  - the model/helper `KillGroup` (`runner.go:1116`) therefore does not reach it.
- Both limits are disclosed as outside the claim (cleanup note :33-36).

**RS-6: Omitted quorum permission. Held.**

- Only execution routes require `checkCapturedActivationQuorum` (`verification.go:366-370`).
- An old ticket with an invalid quorum gets no execution callback or new launch
  (`activation_quorum_test.go:143-149`), while its stop still succeeds.

### F1 response: known protocol refusal

**RR-1: Standalone fresh launch. Held.**

- The refusal is recorded before `beginLaunch` (`protocol_context.go:121-133`).
- `prepareLaunchCycle` lives inside `beginLaunch` (`telemetry.go:86-91`). No session exists, so nothing
  can be charged.
- `finish` classifies the refusal as `protocol_context_refused` (`telemetry.go:175-188`), with no
  trajectory and no budget settlement.
- Test: `protocol_refusal_precharge_test.go:17-63`.

**RR-2: Driver-precharged fixup. Not prevented → N1.**

**RR-3: Prepared verifier ticket.**

- **Runner API: held.** `reserveCapturedVerification` is inside `beginLaunch` (`telemetry.go:96`). The
  test reuses the same in-memory ticket (`protocol_refusal_precharge_test.go:106, 135-141`).
- **Production caller: unreachable → N1.**

### F3 response: activation membership

**RQ-1: Charged replacement, widening, reordering, removal and duplication. Held.**

- `activation_quorum.go:70` `if !slices.Equal(current, members) {`
- Tests: `activation_quorum_test.go:46-73`, including that no journal directory is created (`:63-66`).
- My 2ab9e82 counterexample (attempt 1 rewrites membership, then `--verifier other`) now refuses in
  `FreezeCaptured` (`captured.go:159`), before a ticket exists. The implementer's executed reproduction
  is testimony.

**RQ-2: Verifier outside the activation quorum, or equal to the implementer. Held.**
`activation_quorum.go:59-61`.

**RQ-3: A later unchanged archive replacing the baseline. Held.**

- `unchanged.go:93` `if !slices.Equal(original.Participants, members) {`
- Test: `activation_quorum_test.go:182-200`.

**RQ-4: Widened parent request. Held.** `reconcile.go:249-251`; test `activation_quorum_test.go:243-279`.

**RQ-5: Terminal and stop remain independent of quorum. Held.**

- A quorum-changing patch still retains its terminal (`activation_quorum_test.go:22-29`).
- An old ticket can still be stopped (`:110-180`).

**RQ-6: Live prompt differs from activation membership before verify. Refused only after launch → N3.**

### F6 response: publication coupling

**RF-1: `Finish` performs no history reads. Held by reading.**

What remains under the guard:

- ledger inspection, state read and structural validation (`state.go:514-522`);
- two Git brackets and tree digests (`state.go:109-164`);
- one archive capture with full inspection (`snapshot.go:171-268`);
- revalidation and a state write (`state.go:657-662`).

The cost scales with the current source and the state JSON, not with archived history. The 5.193 s
figure is testimony.

**RF-2: `Begin` still uses full `withState` (`state.go:581`). Consistent with the design.** It runs
before spawn, so no completed outcome is at risk.

**RF-3: Guard contention. → N2.**

### Archive stability

**RA-1: Timestamp-only transition, replacement and repeated drift. Held.**

- Stability is checked only after the full archive, canonical, tree and link checks
  (`snapshot.go:681-689`, then `690-696`).
- Identity, mode and size must match across the initial, opened, ended and named views
  (`archive_stability.go:24-26`).
- A timestamp-only change permits one recheck with `origin=info, allowTimestampTransition=false`
  (`snapshot.go:693`).
- That recheck requires the same file (`snapshot.go:547-549`) and refuses a second transition
  (`archive_stability.go:31-33`).

**RA-2: Duplicate member selection or extraction in the reread. Held.**

- The second pass receives `nil, nil` for destination and selection (`snapshot.go:693`).
- Symlinks are created only after stability succeeds (`snapshot.go:697-715`).
- Restore removes its new directory on any error (`snapshot.go:738-743`).

## Findings

No CRITICAL finding.

### [MAJOR] N1: Known protocol refusals still spend one-time authority on both production callers

**What is wrong (PRIMARY).** The prevention sits at the runner's launch boundary
(`protocol_context.go:121-134`). Both production callers that commit one-time authority do so before
that boundary.

*Driver-managed fixups.*

1. The driver reserves before calling the implementer (`driver/impl.go:361`).
2. `reserveFixupCycle` installs the trajectory observer (`driver/cycle_budget.go:75`) and charges the
   step and the fixup cycle:
   - `driver/cycle_budget.go:84` `if err := budget.ChargeStep(ctx); err != nil {`
   - `driver/cycle_budget.go:87` `n, err := budget.ChargeCycle(ctx, budget.Fixup)`
   - its comment (`:72-73`) says: "the nested real runner reuses this session".
3. With an active trajectory, `AfterCycle` appends the charged attempt with no launch (`state.go:457`).
4. `Fixup` calls `runner.RunFixup` (`app/driver_impl.go:471`). That joins the existing session
   (`phase58.go:62`; `cycle_session.go:45-57`) and reaches `execAgentProcess` (`phase58.go:122`), then
   `beginProtocolLaunch` (`runner.go:1050`).
5. A known refusal now returns before `beginLaunch`, so `trajectory.Begin` (`launch_budget.go:76-81`)
   never runs.
6. The result is a spent charge with `Launch == nil` and `Terminal == nil`:
   - unchanged reconciliation refuses (`unchanged.go:45`);
   - verification refuses (`captured.go:126`);
   - the next fixup refuses (`state.go:406-409`);
   - continuation refuses (`continuation.go:123-125`);
   - the driver escalates (`driver/impl.go:370-371`).
7. The test that models this shape asserts exactly that halt: `Launch`/`Terminal` nil and completion
   refused (`protocol_refusal_precharge_test.go:93-102`).
8. The same ordering spends a driver cross-review cycle before its round runs
   (`driver/cycle_budget.go:66-69`). There is no trajectory halt there, but budget is spent.

*`trajectory verify`.*

1. My grep found one production caller of `PrepareCapturedVerification`: `app/trajectory_verify.go:226`.
   That caller creates the charge's single journal directory before launch:
   - `verification.go:322` `dir, err := openVerificationDirectory(b, r.Charge.EntryKey, true)`
   - `verification.go:283-285`
     `if err = base.Mkdir(charge, 0700); err != nil { return nil, errors.New("verification already reserved or unavailable; preserve it for explicit recovery")`
2. The known refusal is detected later, inside `runner.RunConsult` (`trajectory_verify.go:256`; then
   `consult.go:139` and `runner.go:1050`).
3. Every CLI run allocates a fresh run ID (`trajectory_verify.go:196`) and prepares again (`:226`).
   There is no flag to reuse a prepared ticket (`:35-39`).
4. The note says "A previously unconsumed verifier ticket remains available after restoring valid
   authority" (protocol-refusal note :24-25). That holds at the runner API, but no production caller can
   reach it:
   - a retry refuses at `Mkdir`;
   - the refused run's parent result (`trajectory_verify.go:216-225`) cannot reconcile, because its
     launch record is absent (`reconcile.go:267-270`).
5. Premise (RECALL): `os.Root.Mkdir` fails on an existing directory.

**Why it matters.** The note makes three statements:

- "a fresh refusal spends no cycle/step/launch budget" (:30-31);
- the launch boundary records refusals "before preparing a cycle or making any fresh reservation"
  (:16-17);
- the change is "prospective prevention of one known pre-spawn refusal class only" (:47).

It does not say that driver-managed fixups (the `auto_implement` path) and the verifier CLI still commit
their one-time authority before detection. For those callers, the refused-after-charge part of my F1
is not prevented for future launches either.

**Suggested fix.**

- Before `reserveFixupCycle` and before `PrepareCapturedVerification`, run the same shared renderer for
  the launch info that will follow (phase `fixup` or `trajectory-verification`).
- On refusal, record the unstarted request/terminal as `protocol_context.go:121-133` does, then stop
  without charging or preparing.
- Keep the launch-boundary check authoritative. A pre-check pass does not guarantee a launch pass, and
  that time-of-check window must stay disclosed.
- Add driver-level and CLI-level refusal tests.
- At minimum, narrow the note, docs and IMPLEMENTATION.md to standalone runner launches and name these
  two paths as open F1 cases.

### [MINOR] N2: Both new control routes still contend for the guard that unbounded history readers hold

**What is wrong (PRIMARY).**

1. Both routes acquire the same cycle guard, with a 30-second wait inside the caller's context
   (`state.go:494-496`:
   `wait, cancel := context.WithTimeout(ctx, 30*time.Second)` … `release, err := budget.AcquireResourceGuard(wait, filepath.Dir(b.Store.Dir))`).
   - `Finish` inherits a 30-second total budget (`telemetry.go:194`).
   - The runner's stop inherits a 20-second total budget (`runner/trajectory_verification.go:65`
     `context.WithTimeout(context.WithoutCancel(ctx), 20*time.Second)`).
2. Full readers hold that guard for their entire validation (`state.go:525-532`, then `fn` at `:533`):
   - `Inspect`, reachable from the unattended `parley trajectory inspect` (`app/trajectory.go:69-77`);
   - `InspectHistory`;
   - the unchanged, reconciliation, continuation and parent-recovery previews;
   - `RequireResolved` (`driver/impl.go:316`).
3. A concurrent fixup attempt also takes the guard (`cycle_binding.go:231`) and calls `BeforeCycle`
   (`:253`). That runs the full `checkStateSnapshots` (`state.go:403`) before the cheap refusal it will
   reach anyway (`state.go:406-409`).
4. Testimony: a full `Inspect` of the retained 128-attempt / 256 MiB fixture took 197.013 s
   (history-envelope note :15-16).

**Consequences (by reading; magnitude UNVERIFIED).**

- A `Finish` that cannot acquire the guard in time returns an error. The runner then records
  `trajectory_failure` with no trajectory terminal (`telemetry.go:197-199`), which is the F1-class shape.
  If the guard arrives just before the deadline, see N6.
- A stop that cannot acquire the guard leaves a registered criterion in its own session unsignalled.
  Only the model/helper group is killed (`runner.go:1112-1116`).

The terminal note disclaims "a universal deadline guarantee", but it does not name this contention
mechanism.

**Suggested fix.**

- In `BeforeCycle`, evaluate the `trajectoryHistory` refusal before `checkStateSnapshots`. That reorder
  can only refuse.
- Let read-only full validators inspect archives and results outside the guard. Before returning, they
  re-acquire it and require unchanged state and ledger digests.
- Otherwise, document that a concurrent inspection can still cost a live terminal or stop.

### [MINOR] N3: Membership mismatch is detected only after the one-use ticket and launch are consumed

**What is wrong (PRIMARY).**

- **Before launch.** `verifyTrajectoryWithAgent` pins archive membership through `FreezeCaptured`
  (`trajectory_verify.go:188`). The helper request's participants, however, come from the live prompt:
  - `trajectory_verify.go:86-87` `current, found := findIdeaStatus(status, idea)` /
    `if !found || !slices.Contains(current.Participants, implementer) || !slices.Contains(current.Participants, verifier) {`
  - `trajectory_verify.go:200` `Participants: participants`
- **Request-scope check.** It compares the live prompt with a request built from that same prompt
  (`:109-114`), and the helper repeats that (`:339`).
- **Ticket preparation.** `PrepareCapturedVerification` receives no participants (`verification.go:300`).
- **Reconciliation.** Only here is activation membership required (`reconcile.go:249-251`):
  `if !slices.Equal(req.Participants, members) { return derived, errors.New("helper request differs from the original activation quorum")`
- **Recovery.** Parent recovery reuses that derivation (`parent_recovery.go:219`), and a new ticket for
  the charge is refused (see N1). The charge becomes permanently unverifiable, after a verifier model run
  and every criterion execution.

**Triggers.**

- **(a) Live membership edit.** Membership in the live prompt changes after the patch's after-archive
  (removed, reordered or added) without a charged patch.
  - Before F3 this was accepted end-to-end; that was my F3.
  - It now refuses, but only after launch.
- **(b) Parser divergence (pre-existing).**
  - The live list is parsed by `parseList` (`workspace.go:352, 394-410`), which strips a trailing `]`
    only when it is the final byte (`:397`).
  - `participants: [builder, reviewer, other] # original quorum` therefore yields a last element of
    `other] # original quorum`.
  - yaml.v3 (`activation_quorum.go:27-33`; `reconcile.go:194`) yields `other`. RECALL: YAML treats the
    trailing text as a comment after a flow sequence.
  - Verifier `reviewer` passes the pre-launch check, and reconciliation refuses.
  - This divergence already refused late through `checkReconciliationScope` (`reconcile.go:129`).

SELF-CORRECTION (claude-1, 2ab9e82 review, "Quorum parsers" under Remaining uncertainty):

- It replaces "Every divergence I traced only refuses" with "Every divergence I traced refuses, but some
  refuse only after the ticket and verifier launch are consumed".
- The added consequence is a strengthening and needs a non-owner verdict.

**Coverage.**

- `activation_quorum_test.go:243-279` tests a rewritten request at reconciliation.
- `trajectory_verify_test.go:262-272` covers an unrostered agent.
- I found no test asserting a pre-launch refusal when live membership differs from activation membership.

**Suggested fix.**

- Return the activation members that `FreezeCaptured` already computes (`captured.go:159`), or add a
  guarded reader for them.
- Require `slices.Equal(participants, members)` before `PrepareCapturedVerification`
  (`trajectory_verify.go:192-226`) and again before `ExecuteCapturedVerification` (`:339-342`).
- Parse live membership with the same YAML decoder, or refuse when the two parsers disagree.

### [NIT] N4: The "changed-original-ledger" case rewrites state and never reaches `Finish`'s own checks

- The case edits the trajectory state's charge, not the published ledger
  (`terminal_authority_test.go:90-95`):
  `state.Attempts[len(state.Attempts)-1].Charge.ReserveMicros = &one`
- The error message it asserts comes from `reservation_receipt.go:40` via `state.go:308`. That runs in
  `withStateAuthority` (`state.go:522`), before the `Finish` closure. The frozen-state check
  (`state.go:634-637`) would also refuse it.
- By reading, deleting `state.go:631-633` survives all four cases.
- That deletion is harmless only while the session receipt equals the stored charge, and `Begin`
  compares only `EntryKey` (`state.go:593`).
- The cleanup ledger fixture has the same shape (`verification_cleanup_test.go:204-213`).

Suggestion: rename the case, or add one that mutates the published ledger entry while leaving state
unchanged.

### [NIT] N5: `Finish` can re-publish a deleted historical archive from the live tree

- **Deterministic addressing.** Archives are deterministic content addresses:
  - `snapshot.go:150` `ModTime: time.Unix(0, 0)`;
  - the metadata embeds the exact `Source` (`:214`);
  - a missing final path is (re)published (`:252-263`).
- **Mechanism.** If an attempt's after-`Source` equals its before-`Source` (same commit, status and
  tree), it produces the same reference. If that file was deleted, the control route now reaches
  `CaptureSnapshot` and restores it, and a later `Inspect` passes.
- **Effect.** The bytes are hash-identical, so no false evidence arises. However, the deletion signal
  disappears. The terminal note's "No old source or missing terminal is reconstructed from today's tree"
  (:21) is therefore broader than the code.
- **Test coverage.** The retention test cannot observe this. Its child always edits `source`
  (`terminal_authority_test.go:48`) before the test asserts that `Inspect` refuses (`:132-134`).

Suggestion: qualify the wording, or add the unchanged-source variant and decide the behaviour explicitly.

### [NIT] N6: A deadline can be persisted as `source-unavailable`

- `Observe(ctx, r.root)` (`state.go:642`) uses the context that guard acquisition may have nearly
  exhausted (`telemetry.go:194`; `state.go:494`).
- A context error is then persisted as `SnapshotError: "source-unavailable"` (`state.go:654-656`). No
  after-source is retained, even though the tree was available.
- This label pre-dates the slice, and the new route makes the case less likely. N2 is one way to reach
  it.

Suggestion: use a distinct deadline label, or document that `source-unavailable` means "after-source
not retained".

## Abnormal-recovery proposal: bounded feedback (unimplemented)

I read `codex-1-abnormal-recovery-proposal-20260914.md` as discussion input. Nothing below treats a
suggested control as implemented, or authorizes historical reconstruction.

**AP-1: Separate observation from permission in state, not only in the UI.**

- The existing attended `Continue --acknowledge-inconclusive` clears any sequence in
  `InconclusivePending` (`continuation.go:74-77, 197-206`).
- If an abnormal observation were appended as an ordinary `Inconclusive` resolution, that generic
  acknowledgment would silently become custody authority.
- Instead, keep abnormal sequences out of that set, or require a separate custody decision record that
  `validateTransitions` (`continuation.go:232-280`) checks for them. A hand-edited state then cannot skip
  that record.
- One staged preview may show both authorities, but with two applies, two decision IDs and two hashes.

**AP-2: Preserve streak and closure semantics.**

- The observation must not reset `Consecutive`, matching the rule for unchanged observations
  (`continuation.go:70-72`).
- It must never be `NoRegression`, so `RequireResolved` keeps completion blocked (`state.go:546-548`).

**AP-3: Account for old readers and mixed versions.**

- State decoding rejects unknown fields and non-canonical bytes, and accepts only versions 2 and 3
  (`state.go:248-262, 97-104`).
- A new evidence variant would make older binaries refuse the whole state, including their `Finish` and
  stop. Mixed-version runners would then lose terminals and stops.
- Bump the version with an explicit refusal message, and never roll out while an older runner holds a
  live handle.

**AP-4: Retained facts that could support an attended continuation.** Each fact would be rechecked under
the guard, with full evidence, at apply time.

1. **No instrumented spawn: `Launch == nil` on a charged attempt.**
   - On the path I read, `Begin` precedes spawn: `launch_budget.go:76-81` → `telemetry.go:99` →
     `protocol_context.go:134` → `runner.go:1105`.
   - N1's precharged refusals produce exactly this shape.
   - It needs no human custody assertion. A fresh clean observation equal to the retained before-source,
     a full archive check and intent validation are enough.
   - Promotion must reuse the before-source, because `continuationPreview` currently requires `After`
     (`continuation.go:131`).
   - Other launch surfaces (ACP, interactive) are UNVERIFIED.
2. **Complete, consistent lifecycle.** Either of:
   - a trajectory terminal with retained, equal before and after trees, plus telemetry
     requested/started/terminal records that match the trajectory status and exit. Those records must
     satisfy the identity and ordering clauses of `unchanged.go:147-166`, with an enumerated abnormal
     failure class in place of normal exit;
   - an explicitly unstarted terminal: no `started.json`, `StartedAt` nil, and a class in
     `start_failure`/`protocol_context_refused`/`budget_refused` (`telemetry.go:175-192`).
3. **Custody facts, strongest first.**
   - (i) Fact 1: no instrumented spawn.
   - (ii) A boot-ID change since the recorded start. `Attributed` already treats that as a prior-boot
     process (`procctl.go:127-129`), and it covers escaped descendants on that host.
     - Telemetry retains only a PID (`telemetry/record.go:86`).
     - A boot ID exists only in the run store's `agent.started` identity (`durablekill.go:104-111`), so
       future launches would need to bind it to the attempt.
   - (iii) Same boot, attributed leader gone and group empty. This is bounded evidence only: it cannot
     exclude escapees that started a new session.
   - (iv) A scoped human assertion, required whenever (i) and (ii) are unavailable.
     - It must name the author, time, invocation ID, PID/PGID, host and boot ID, worktree and helper
       roots, and the tool-generated residual unknowns.
     - It is recorded as unverified testimony. It is never inferred from `--yes`, elapsed time or an
       empty PID search.
     - It should state the practical risk: a surviving writer's later edit gets attributed to the next
       charged attempt, which can corrupt the D6 regression streak.
4. **The existing continuation predicate still holds.** A clean current tree equal to the retained after
   bytes, descending from the patch HEAD (`continuation.go:131-136`).

**AP-5: What must still block continuation, attended or not.**

- **`Launch` present, no `Terminal`.** Causes include a crash, a failed `Finish`, or a `Begin` write that
  landed but errored (RT-9). Both the outcome and the after-source are unknown.
- **`source-unavailable` or `archive-unavailable` terminals** (`state.go:322-331`).
  - The after-source is not retained.
  - The label can be a deadline artifact (N6).
  - Continuing from today's tree would reconstruct it.
- **Contradictory records.**
  - Trajectory and telemetry status disagree after a landed write followed by the `trajectory_failure`
    override (`telemetry.go:197-199`).
  - PID or ordering mismatches between records.
- **Changed material source.** Independent verification only (`captured.go:129-131`).
- **Consumed or incomplete helper tickets.**
  - Cases: a stopped helper, a missing or failed receipt, a torn process record, or a v1/v2 claim that
    carries only PID and time (`verification.go:49-55`).
  - There is one directory per charge (`verification.go:283-285`). A retry therefore needs its own
    versioned, linked design, not a continuation.
- **v1 digest-only histories** (`state.go:97-100`).

**AP-6: Refutation cases to add to Q4.**

- N1's driver and CLI refusal paths.
- Guard contention during `Finish` (N2).
- A deadline-labelled `source-unavailable` terminal (N6).
- A boot-ID change combined with PID reuse.
- A landed `Begin` write that reported failure.

## Remaining uncertainty

- **Executed evidence.** I executed nothing. Every counterexample, surviving-mutant and consequence
  statement comes from source reading. All facilitator- and implementer-reported results remain
  testimony.
- **Trust boundary.** Hashes detect mismatch, not coordinated same-UID rewriting of state, journals,
  telemetry and archives. RS-2 relies on that disclosed boundary.
- **Pre-existing versus new.** Without Git I cannot diff against 83d0475 or 2ab9e82. Statements about
  what a slice changed rest on the current tree plus the implementer's notes.
- **N2 magnitude.** Unmeasured. I did not read `lockWithReady` fairness or the ledger's own locking.
  `CaptureSnapshot` takes a nested snapshot-directory guard (`snapshot.go:178-184`); the lock order is
  consistent for the callers I traced (`Activate`, `Continue`, `Finish`), but that trace is not
  exhaustive.
- **N1 exact spend.** I did not verify whether the runner's own step charge (`launch_budget.go:84-97`)
  is idempotent inside the driver's step session. This affects only the spend count, not the halt.
- **Runtime semantics.** These premises are RECALL:
  - yaml.v3 comment handling (N3);
  - `os.Root.Mkdir` on an existing directory (N1);
  - `fsutil.ReplaceSyncedFile` ordering (RT-9).
- **Not assessed.**
  - Windows: verification journals refuse there (`verification.go:251-253`).
  - ACP and interactive spawn ordering.
  - `evidence.RunCriterionControlled` start/release semantics.
  - Escaped-descendant inactivity and general helper recovery.
- **Out of scope.**
  - Live launch, concurrency and closure evidence; packet and full-six experiments.
  - Historical quorum and signatures.
  - Final HTML and delivery follow-ups.
  - This review cannot change historical membership, participant artifacts or signatures.
