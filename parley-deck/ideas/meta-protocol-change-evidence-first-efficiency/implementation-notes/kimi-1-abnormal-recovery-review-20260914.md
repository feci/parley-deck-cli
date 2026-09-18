---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
reviewed-commit: 6f31201c987ecb085921b537d0345e134380ec05
review-kind: supporting-source-proposal
---

# Abnormal-recovery proposal: supporting source/proposal review (kimi-1)

## Summary

I reviewed codex-1's unimplemented abnormal-recovery proposal
(`implementation-notes/codex-1-abnormal-recovery-proposal-20260914.md`, source-base
83d0475) against the frozen source at 6f31201, after reading Claude F1 in its
original statement. The proposal's core separation — an explicitly inconclusive,
evidence-bound observation distinct from any continuation authority — is sound and
maps onto machinery that already exists. My principal conclusions:

1. **The observation needs no new custody assertion, and none is admissible.** For
   every coverable class the runner provably reaped the direct child before
   `terminal.json` exists, and every residual custody unknown (escaped descendants,
   PID reuse, hidden writes) fails closed at the next charge/continuation boundary.
   Recording the observation and authorizing continuation must stay two separate
   controls; the second already exists as the attended `continue
   --acknowledge-inconclusive` path.
2. **Two constraints the proposal does not state are load-bearing.** (a) The
   evidence variant must reuse the `Unchanged` preview field (version bump), or the
   regression streak silently resets — a laundering channel. (b) Attempts without a
   published terminal cannot be resolved by any evidence variant; state validation
   forbids it, so they need a separate state migration and must stay blocked here.
3. **The proposal's premise is partially stale at 6f31201.** The
   protocol-context-refusal class no longer spends a charge (it is recorded and
   returned before `beginLaunch`). The live refused-after-charge classes are the
   budget-layer refusals and start failures. A further sub-class — charge spent,
   `trajectory.Begin` refused, no launch and no terminal — is reachable, has no
   recovery, and is not named in the proposal.

Findings: two MAJOR, three MINOR, two NIT. This is a bounded supporting review of a
proposal and its source seams — not a Phase-6 review, a signoff, an acceptance, a
finding withdrawal, a quorum change, a pilot amendment, or a numerical ceiling. I
executed nothing.

## Scope and provenance

- Launch: kimi-1, protocol attestation `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- Method: repository files at the frozen commit read with file/search tools only.
  Two read-only directory listings (`ls`, `wc`) were used for orientation; no other
  shell, no Git operation, no build, test, race, overlay, or program execution.
  Line numbers refer to this tree.
- Provenance (§15.2):
  - **PRIMARY (source):** file, line and quotation I read in this tree.
  - **Claude F1** (`claude-1-reconciliation-review-20260914.md`): read for its
    original statement; its facilitator-reported test/race results are testimony I
    did not verify.
  - **Codex proposal:** the artifact under review.
  - **RECALL:** Go/OS process semantics stated from memory; marked where relied on.
- Read in full: `internal/trajectory/{unchanged.go, continuation.go, reconcile.go,
  state.go, verification.go, verify.go, parent_recovery.go, reservation_recovery.go}`,
  `internal/procctl/{procctl.go, procctl_unix.go, procctl_linux.go}`,
  `internal/telemetry/record.go`, `internal/runner/{telemetry.go, launch_budget.go,
  protocol_context.go, supervision.go}`, `internal/app/{trajectory.go,
  trajectory_reconcile.go, budget_attended_other.go}`,
  `internal/trajectory/captured.go:80-146`, `internal/budget/cycle_binding.go:230-287`,
  `docs/agent-runtime-configuration.md:1475-1494`, FINAL.md, IMPLEMENTATION.md
  known-gap passages, the codex proposal, Claude's F1 review.
- Read in part (grep-verified call sites): `internal/runner/{launch.go, runner.go,
  acp.go, handoff.go, cycle_budget.go}`, `internal/budget/cycle_observer.go`.
- Not claimed: any execution verdict, Windows behaviour, the packet/pilot
  experiments, quorum decisions, acceptance.

## Refutation attempts (against the proposal's shape)

**R1 — Can the observation launder a regression-streak break? Held only under a
field-reuse constraint (K2).** `trajectoryHistory` resets the consecutive-regression
count for any `Inconclusive` resolution whose preview lacks the `Unchanged` field:
`continuation.go:69-72` `case Inconclusive: if r.Preview.Unchanged == nil {
h.Decision.Consecutive = 0 }`. A changed-source inconclusive resets the streak by
design; an abnormal *unchanged* observation must not — it is not a patch, and D6
says repeated unchanged criticism is not a new regression, but nor is it evidence
that clears one. If the variant lives in a new sibling field (e.g.
`Preview.Abnormal`), the streak resets and a sequence R(1) → abnormal-U(2) → R(3)
would report `Consecutive == 1` instead of 2, suppressing the two-streak review
gate (`continuation.go:79-84`). Reusing `Unchanged` with `Version: 2` preserves the
rule untouched. Any other placement needs a deliberate amendment plus a
streak-preservation test.

**R2 — Can the observation itself grant continuation? Held, if existing machinery
is reused.** After any resolution recording an `Inconclusive`, the next charge is
refused at `state.go:407-409` ("trajectory awaits independent reconciliation or an
attended review decision; further fixup is refused") and `continuationPreview`
refuses at `continuation.go:123-124`. Release requires the attended `Continue`
control: platform TTY probe (`trajectory_reconcile.go:70-73`,
`budget_attended_other.go:5-7`), explicit `DecisionID` + bounded scrubbed `Reason`
(`continuation.go:168-170`), preview-digest pinning, per-decision-id replay
protection (`continuation.go:178-186`), and per-pending-sequence acknowledgment
(`continuation.go:197-206`). The control's own comment states it "never removes a
regression, grants budget, changes a verdict or closes an implementation"
(`continuation.go:162-165`). An abnormal observation that reuses this path grants
nothing beyond the remaining ledger budget; no new authority surface is needed.

**R3 — Can the variant resolve an attempt with no published terminal? Refused by
state validation (K4).** `validateTransitions` requires every resolved attempt to
have `s.Attempts[i].Terminal != nil` (`continuation.go:238`). Recording any
resolution for a terminal-less attempt makes `validateState` fail, and since
`withStateAuthority` runs `validateState` on every guarded read (`state.go:522`),
the entire trajectory would refuse all further operations. Terminal-less attempts
(`trajectory_failure` at `telemetry.go:193-200`, `telemetry_failure`, a crash
between `Begin` and `Finish`) therefore cannot be covered by an evidence variant;
they need a state-version migration and stay blocked in the bounded control.

**R4 — Does the refused-after-charge class still include protocol-context
refusals? No at 6f31201 (K3).** `beginProtocolLaunch` computes the refusal first
and, on refusal, records only request + terminal and returns — "Retain its
request/terminal before any fresh cycle, step, launch or helper-ticket
reservation" (`protocol_context.go:118-133`). `beginLaunch` (which performs
`budget.ChargeCycle` at `launch_budget.go:70` and `trajectory.Begin` at
`launch_budget.go:77`) is never reached. This implements the second arm of Claude
F1's suggested fix for counterexample B, first bullet. The proposal's source-base
83d0475 predates it; the live refused-after-charge classes are the budget-layer
refusals after the charge (`ChargeStep` `launch_budget.go:92`, monetary binding
`:111`, launch reservation `:139-141`) and spawn/build failures
(`launch.go:108-113`, `launch.go:340-346`, `runner.go:1105-1107`).

**R5 — Can a missing helper receipt be synthesized from completed steps? Refused,
correctly.** `readCapturedVerificationJournal` hard-fails on a missing receipt:
`verification.go:644-646` "verification terminal receipt unavailable; retain
partial observations". The receipt binds `Steps`/`LastSHA256` to the journal
(`verification.go:647`), binds the helper identity (`HelperPID` == claim PID), and
its `FinishedAt` must fall inside the helper invocation's observed terminal window
(`reconcile.go:291-293`). Steps alone cannot prove the helper exited, cannot bind a
completion time, and cannot be assessed — `AssessCaptured` runs only after receipt
validation (`verification.go:664`). A stopped or partial helper can never be
relabeled a verifier: a `FailureStage` receipt is refused (`verification.go:650-652`)
and a `stop.json` fence blocks claiming/starting (`verification.go:80-85`,
`:428-430`). Concur with the proposal; this case must remain permanently blocked
absent a separately designed retry-ticket API (out of scope here).

**R6 — Can leader absence, PID search, elapsed timeout, or a generic `--yes`
establish descendant inactivity? No.** `KillGroup` signals one process group
(`procctl_unix.go:40-48`); a descendant that called `Setsid` escapes the group and
is invisible to both the fence and the kill. The historical claim v2 records only
PID + time (`verification.go:49-55`, written at `:435`); `procctl.Attributed`
fails closed on any missing identity facet — boot ID, proc start, PGID, command
(`procctl.go:107-127`) — so no historical helper PID can ever be safely signalled
or declared dead. Registered criterion processes carry full identity
(`verification.go:66-71`, `:93-97`) and the stop route kills only those, with
strict attribution (`verification.go:134-152`) — a *narrower* guarantee than "the
tree is dead", as the proposal says. The residual risk is bounded structurally,
not by assertion: any late write to the checkout is caught fail-closed at the next
charge (`Observer.BeforeCycle`, `state.go:417-423` `actual != before`) and at
continuation (`continuation.go:131-136`: clean tree, exact retained after-bytes,
ancestry). Failing closed means a halt, never a false clean. (The F4
untracked-exclude channel is the standing exception — inherited, pre-existing, and
UNVERIFIED per Claude F4; RECALL-level Git semantics on my part as well.)

**R7 — Can an operator-supplied class or contradictory records enter the
observation? Must be refused.** The class must be a pure derived function of the
retained `terminal.json` (`status`, `failure_class`, `exit_code`) cross-checked
against the trajectory terminal exactly as `unchanged.go:157` does
(`o.Status != a.Terminal.Status || *o.ExitCode != *a.Terminal.ExitCode`), with
class-conditional presence rules: `started.json` required for process-started
classes and required-absent for never-started classes; exit code required for
exited classes and required-nil for refused classes (mirroring the nil-field
discipline already at `unchanged.go:148`, `:161`, `:165`). A started-record present
for `start_failure`, or an exit code present for `budget_refused`, is contradictory
evidence and must refuse.

**R8 — Can the observation spend, refund, or re-charge? No.** The unchanged
control already sets the shape: "It spends nothing, executes nothing and grants no
continuation" (`unchanged.go:214-215`). The charge stays spent in the ledger (D6:
"A failed charged attempt stays spent"); `len(s.Attempts) != len(ledger.Entries)`
is refused (`state.go:270`), so the observation cannot add or drop ledger facts.
The next charge, when acknowledged, consumes the same monotonic remaining budget.

## Findings

### [MAJOR] K1: Charged attempt with no launch and no terminal — a reachable, unrecoverable sub-class the proposal does not name

**What is wrong (PRIMARY).** The charge is spent and the attempt row is appended
inside `budget.ChargeCycle` (observer `AfterCycle`, `cycle_binding.go:281-284`;
"A failed AfterCycle never refunds the charge", `cycle_observer.go:14`) at
`launch_budget.go:70`. `trajectory.Begin` then runs at `launch_budget.go:77` and
can refuse: policy mismatch, charge mismatch, no charged attempt, implementer
mismatch, or "source changed between reservation and model launch"
(`state.go:589-602`). On that refusal `l.trajectory` stays nil, so `finish`
(`telemetry.go:99-100`, `:193-201`) never publishes a trajectory terminal. The
persisted attempt then has `Launch == nil`, `Terminal == nil`, `After == nil` —
valid state (`state.go:314-317` only constrains non-nil terminals) but
unreconcilable: unchanged refuses at `unchanged.go:45` (`a.Launch == nil || ...
|| a.Terminal == nil`), patch verification refuses at `captured.go:126`
(`a.Launch == nil || a.Terminal == nil || a.After == nil ...`), and
reservation-recovery's repairing branch does not apply because the row exists
(`reservation_recovery.go:353` `missing := len(s.Attempts) == len(i.Before.Attempts)`
is false; the else branch only re-syncs). The telemetry terminal records
`failed`/`budget_refused` (`telemetry.go:189-192`), mislabelling a trajectory
binding refusal.

**Reachability.** The window is the in-process gap between `ChargeCycle`
(`launch_budget.go:70`) and `Begin` (`launch_budget.go:77`); an external worktree
edit in that window trips `state.go:600-602`. Policy/charge mismatches are further
triggers. Narrow but real, and exactly the "unstarted refused-after-charge" family
the proposal claims to cover — yet this member has no terminal and therefore (K4)
cannot be covered by the proposed observation.

**Suggested fix.** Name the sub-class explicitly as blocked in the bounded
control, and either (a) hoist Begin's refusable checks into the pre-spend
reservation path (the Observer already re-observes source at `state.go:417-423`),
or (b) add a terminal-publication route for a charge-matched launch-less attempt,
reusing the new `withStateControl` authority. (a) is the smaller change. This is a
source defect/gap, distinct from the proposed API.

### [MAJOR] K2: Evidence-variant field placement is load-bearing for the regression streak

**What is wrong (PRIMARY, design constraint on the proposed API).** The streak rule
keys on field presence, not outcome semantics: `continuation.go:69-72`. Any variant
not carried in `ReconciliationPreview.Unchanged` silently resets
`Decision.Consecutive` (walkthrough in R1), suppressing the two-streak
`RequiredReviewSequence` gate (`continuation.go:79-84`) — a laundering channel
introduced by the recovery control itself. The proposal asks for "a distinct typed
evidence variant" but does not state this constraint.

**Suggested fix.** `UnchangedEvidence Version: 2` in the same field: add a derived
closed `Class` enum plus recorded exit code, make `StartedSHA256` class-conditional,
keep every v1 binding field with identical meaning (original hashes preserved:
`PolicySHA256`, `AttemptSHA256` = `digest(canonical(a))`, `ScopeSHA256`,
`InvocationID`, record SHAs). Old readers refuse `u.Version != 1`
(`unchanged.go:56`) — fail closed, which is the desired old-reader behaviour; name
it in the docs. Keep the fixed all-unresolved assessment (`unchangedAssessment`,
`unchanged.go:32-38`) so `RequireResolved` still demands a later `NoRegression`
(`state.go:546-548`) and independent clean completion is preserved. If a separate
field is ever preferred, amend `continuation.go:69-72` deliberately with the
streak-preservation fixture from the falsifier list (F-2 below).

### [MINOR] K3: Proposal premise partially stale at 6f31201 — protocol-context refusals no longer charge

R4 gives the mechanics (`protocol_context.go:118-133`). The proposal groups all
"refused after a charge" cases; at the frozen commit the protocol-context member is
already prevented prospectively. The remaining refused-after-charge members are:
step-charge refusal (`launch_budget.go:92`), monetary-binding refusal (`:111`),
launch-reservation refusal (`:139-141`), and pre-spawn build/log/context/`cmd.Start`
failures (`launch.go:108-113`, `:340-346`, `runner.go:1105-1107`). All produce
`failed` + nil exit + absent `started.json`, with the trajectory terminal published
by the new `Run.Finish` (`state.go:619-675`) — coverable by the observation as
never-started classes. The proposal should restate its target class list against
6f31201 so the control is not designed for an already-fixed case. (IMPLEMENTATION.md
:3365-3368 already records the prevention as prospective-only; historical precharged
refusals "remain spent, byte-identical and unresolved" — those remain in scope.)

### [MINOR] K4: Terminal-less attempts are structurally unresolvable by an evidence variant

R3 gives the mechanics (`continuation.go:238`, `state.go:522`). This bounds the
control's scope: `trajectory_failure` (publication failure, `telemetry.go:193-200`),
`telemetry_failure`, and crashes between `Begin` and `Finish` leave the actual
process outcome unknown — no observation can honestly classify them. The newest
`Run.Finish` change shrinks this class to true publication failures (it no longer
re-reads historical archive/intent/result contents: `withStateControl` skips
`checkSourceSnapshots`/`checkResolutions` while keeping ledger/policy/state-digest
binding, `state.go:482-534`, `:626-641`) but cannot eliminate it. Covering it is a
state-version migration (v3 → v4) with its own design; keep it out of the bounded
control and say so. The newest stop-only route (`StopCapturedVerification`,
`verification.go:106-155`) is consistent with this split: exact ticket/launch/
registered-process authority, durable stop fence first, signals only strictly
attributed identities, and it never touches the helper leader itself (the caller
kills the enclosing group) — prospective control, not historical recovery. I found
no defect in either newest change by reading; the skipped evidence checks fail
closed at the next full `withState` read.

### [MINOR] K5: No operator custody assertion is admissible — and none is needed for the coverable classes

The proposal asks which "explicitly scoped human assertion" may permit
continuation. My answer, from source: **the existing attended acknowledgment and
nothing else.** For killed classes the direct child is provably reaped before the
terminal exists — `waitSupervised` drains `waitErr` after every kill
(`supervision.go:161-163`, `:174-175`, `:182-183`), and `started.json` +
`terminal.json` with a classified failure bound the whole observed lifecycle. For
never-started classes no process ever existed for the invocation. Every residual
unknown (escaped `Setsid` descendants, PID reuse after reaping, hidden writes)
fails closed at the next boundary (`state.go:417-423`, `continuation.go:131-136`);
an operator's "writers are stopped" claim adds testimony without a process witness
(R6) and should be rejected as an input — concur with the proposal, and strengthen
it: the assertion is not merely weak, it is *unnecessary*, because the gate it
would unlock already exists with the right semantics (R2). Any future class where
the runner cannot prove reaping (e.g. a terminal published while the child may
still run) must NOT be admitted to the observation on the strength of an operator
claim.

### [NIT] K6: The class enum must be derived, closed, and keep the structured-provider qualification

Derive from (`status`, `failure_class`, `exit_code`) only. Suggested members:
`timeout`, `cancelled`, `no-first-output`, `stalled`, `provider-failure` (structured
only — collection requires structured launch args, `telemetry.go:119`, and the
classification reaches the outcome via `providerFailure`, `telemetry.go:173-174`;
a plain-text provider error with exit 1 remains ordinary `process_failure` and
qualifies under the existing v1 unchanged path, `docs :1484-1486`), `signal-exit`
(exit -1; RECALL: Go reports -1 for a signalled process), `budget-refused`,
`start-failure`. Exclude `trajectory-failure`, `telemetry-failure`,
`unobserved-handoff` (contradictory for a charged fixup), and any
`process-exited` status. Record the exact `failure_class` string and exit code
inside the hashed evidence so re-derivation in `checkResolutions`
(`reconcile.go:358-372`) is byte-exact and the record is self-describing.

### [NIT] K7: Docs/known-gap wording and per-resolution re-derivation cost

`docs/agent-runtime-configuration.md:1486-1488` ("These abnormal or
refused-after-charge unchanged attempts currently block further fixups and
continuation, with no supported recovery control for that class") and
IMPLEMENTATION.md `:333-335` must be amended only for the classes the control
actually covers — K1/K4 classes stay named as unresolved. Claude F6's cost concern
extends: every guarded read re-derives every resolution (`reconcile.go:358-372`)
under the runner's 30 s trajectory publication ceiling (`telemetry.go:194`); a v2
variant adds derivation work per abnormal resolution. Not measured; note only.

## Answers to the proposal's four review questions

1. **Independent observation vs one staged preview.** Separate. The observation is
   deterministic, replayable, evidence-bound, and useful on its own (it preserves
   the abnormal attempt as a first-class retained outcome, per D2/D7's
   failed-start and inconclusive-historical obligations); continuation authority
   already exists with attendance, decision-id, and reason. One staged control
   would fuse a deterministic derivation with an attended decision and invite the
   unattended path to borrow the attended one's authority.
2. **Facts supporting the inconclusive observation / authorizing assertion.**
   Facts: everything the v1 path already binds (policy, charge checked against the
   live ledger — `state.go:308`, archives validated whole before any member bytes,
   tree equality, launch identity, the three create-once invocation records —
   `record.go:323` `O_CREATE|O_EXCL` — cross-checked field-by-field), plus the
   derived class under R7's conditionals. Assertion: only the existing attended
   `Continue` with `--acknowledge-inconclusive` (R2). No custody assertion (K5).
3. **Cases that must remain blocked.** Terminal-less attempts (K4, including K1's
   no-launch sub-class); missing/torn/contradictory records; changed trees (patch
   verification only, `unchanged.go:48-49`, `captured.go:129-130`); snapshot/
   archive-unavailable attempts (`unchanged.go:45` analogue); old helpers with
   steps but no receipt (R5); historical helper-PID containment (R6); v1
   trajectories (`state.go:98-100`); anything whose class cannot be derived from
   retained records.
4. **Versioned variant.** K2 gives the shape: same field, `Version: 2`, derived
   closed class, class-conditional record presence, fixed all-unresolved
   assessment. Old readers fail closed (`unchanged.go:56`). Original hashes,
   spent charge, remaining budget, streak, and the clean-completion gate are all
   preserved by construction (R1, R2, R8). Exact replay mirrors
   `unchanged.go:197-203` and `:231-245`; concurrent apply keeps the
   expected-digest refusal (`unchanged.go:247-248`); re-derivation is a pure
   function of retained records (K6).

## Executable falsifiers a later implementer must satisfy

Fixtures below mirror the existing harness shapes (`unchanged_test.go`,
`verification_control_test.go` — read, not run). Each must fail when the named
mutant is applied.

- **F-1 timeout fixture.** Fake child hangs; hard timeout fires; after-tree equals
  before-tree. Preview derives class `timeout`; apply succeeds; history shows
  `Unreconciled == 0`, `InconclusivePending == [N]`. Next charge refused
  (`state.go:407-409`) until attended `--acknowledge-inconclusive`;
  `RequireResolved` still refuses (`state.go:546-548`).
- **F-2 streak fixture.** R(1) → abnormal-U(2) → R(3) yields `Consecutive == 2`
  and `RequiredReviewSequence == 3`. Mutant: variant placed outside
  `Preview.Unchanged`, or `continuation.go:70` amended wrongly, must fail this.
- **F-3 class conditionals.** `start_failure` with a `started.json` present →
  refuse; `budget_refused` with an exit code present → refuse; `timeout` with
  absent `started.json` → refuse.
- **F-4 signal-handled exit.** Child traps SIGTERM and exits 3 on the timeout kill
  (`procctl_unix.go:40-48` grace): outcome is `failed`/`timeout` with exit 3;
  observation must record class `timeout` + exit 3 and must never read it as
  `process-exited`. Mutant: dropping the `FailureClass == nil`/`process_failure`
  discrimination (Claude F2's m1/m2 applied to the new predicate) must fail.
- **F-5 no laundering.** One-byte tracked change during the abnormal attempt →
  refuse with "changed material source requires independent patch verification";
  a dirty-but-tree-equal attempt follows the existing v1 semantics and records the
  F4 limitation explicitly.
- **F-6 no synthesis.** Terminal-less attempt (simulated `trajectory_failure`) →
  refuse; no class may be derived from `requested.json` alone; K1's no-launch
  attempt → refuse with a named error.
- **F-7 replay and concurrency.** Second applier with a stale expected digest →
  refused; exact replay returns the retained preview byte-identical (SHA256
  equality); re-preview after apply equals the stored preview
  (`unchanged.go:197-203` shape); concurrent double-apply converges
  (`unchanged.go:247-248` shape).
- **F-8 old-reader fail-closed.** A v2 resolution validated by the v1 predicate
  (`u.Version != 1`) refuses.
- **F-9 no liveness consultation.** Derivation performs no PID/proc probe: a
  fixture with a live, reused PID at the recorded PID derives byte-identical
  evidence. (Guards against custody-by-probe creeping in.)
- **F-10 ledger preservation.** Observation leaves `len(Attempts) ==
  len(ledger.Entries)` and every entry untouched; the next acknowledged charge
  decrements the same remaining budget; no refund path exists.
- **F-11 completion gate.** After acknowledgment plus a subsequently verified
  clean patch, `RequireResolved` passes only with latest outcome `NoRegression`;
  the abnormal resolution never satisfies it alone.

## Recommended bounded next step

Implement, in one slice: (1) `UnchangedEvidence Version: 2` with the derived
closed class enum and class-conditional record presence (K2, K6); (2) the
derivation predicate inside `unchangedPreview`'s existing authority checks,
refusing K1/K4 cases with named errors; (3) preview/apply replay and concurrency
mechanics copied from the v1 control; (4) falsifiers F-1…F-11 as tests, including
the Claude-F2-style mutant overlays for the new predicate; (5) docs and
IMPLEMENTATION.md wording amended for covered classes only (K7). Explicitly out of
this slice: terminal-less recovery (state migration), helper retry tickets, any
custody assertion input, any change to the attended `Continue` control, and any
fix for K1 beyond naming it (recommend the (a) hoist as a separate small slice).
This review recommends boundaries only; it authorizes nothing.

## Remaining uncertainty

- **No execution.** Every fixture, mutant, and counterexample above is derived
  from source reading. Facilitator-reported test/race results cited by Claude and
  in IMPLEMENTATION.md remain testimony I did not verify.
- **RECALL items.** Go's exit-code reporting for signalled processes (-1), Git
  exclude-source semantics (`info/exclude`, `core.excludesFile` — the F4 channel),
  and `Setsid` group-escape semantics are stated from memory and UNVERIFIED.
- **Same-UID trust boundary.** Hashes detect mismatch, not coordinated same-UID
  rewriting of state, ledger, telemetry and archives (`unchanged.go:20` comment);
  the observation inherits this boundary unchanged.
- **ACP and interactive launches.** I did not audit ACP stop/exit semantics
  (`acp.go:262` call site read only); whether ACP fixups can produce a coverable
  abnormal class is unknown — from what I saw, the risk is false refusal only.
- **Windows.** `openVerificationDirectory` refuses non-POSIX hosts
  (`verification.go:251-253`); attendance and sync behaviour elsewhere are
  platform-split; nothing here was verified on Windows.
- **K1 fix choice.** Option (a) (hoist Begin's checks pre-spend) is recommended but
  unproven; the Observer/Begin check overlap needs a dedicated look so no new
  spend-without-record window is created.
- **Not addressed.** Quorum/exclude-source findings (Claude F3/F4), F6's general
  cost, live experiments, current-source participant acceptance, signatures, and
  final delivery gates. No finding is withdrawn by this review.
