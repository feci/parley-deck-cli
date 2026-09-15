---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
reviewed-commit: 2ab9e82f0dd5a885717483a363de9482ae7c938a
review-kind: supporting-source-only
---

# Unchanged-source reconciliation: supporting source review

## Summary

I read the unchanged-source reconciliation at 2ab9e82 and its dependencies:
`trajectory reconcile-unchanged`, the typed `unchanged` preview branch and how it
plugs into reconciliation, history, parent recovery, snapshots, the CLI, the
runner lifecycle and budget receipts. I found no way, within the disclosed
same-UID trust boundary, to record any of these as an unchanged observation:

- a changed attempt;
- a substituted charge, action, epoch or invocation;
- a rewritten archived scope;
- a mixed verifier provenance;
- a rehashed preview.

These properties hold by reading, and none of them grants a retry, budget,
continuation or completion: charge/action/epoch binding, full archive validation
before any scope bytes are released, streak preservation, exact replay,
parent-recovery dispatch and concurrent-apply convergence.

Findings: two MAJOR, two MINOR, two NIT.

- **F1 (MAJOR).** A charged attempt with unchanged source that timed out, was
  cancelled, was stopped by a watchdog, hit a structured provider failure, was
  signalled, or was refused after charging halts the opt-in trajectory. No
  supported control resolves it. This is disclosed as unresolved and, per the
  implementer's testimony, not introduced by this slice.
- **F2 (MAJOR).** The contradictory-lifecycle tests cannot detect removal of the
  process-qualification predicate. The stale preview digest makes every such
  subtest fail regardless of that predicate.
- **F3 (MINOR).** "Original quorum" is bound to this attempt's before-archive,
  which an earlier charged patch can rewrite. It is not bound to anything frozen
  at activation.
- **F4 (MINOR).** Material equality, the only gate for skipping verification,
  cannot see untracked files hidden through non-tracked Git exclude sources. The
  Git-semantics premise is RECALL, so this finding is UNVERIFIED.
- **F5 (NIT).** The claim that provider failures are refused covers only
  structured provider failures.
- **F6 (NIT).** Rederiving every unchanged resolution under the guard adds work
  inside the runner's 30-second trajectory publication ceiling.

This is not a full-scope Phase-6 review, a signoff or an acceptance. I executed
nothing and issue no execution verdict.

## Scope and provenance

- Launch: claude-1, protocol attestation `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- Method: repository files at the frozen commit, read with file tools only. I
  ran no shell, Git, build, test, race run, mutant or overlay. Line numbers refer
  to this tree.
- Provenance (§15.2):
  - **PRIMARY (source):** a file, line and quotation I read.
  - **Tests:** cited as source I read, not as tests I ran.
  - **RECALL:** runtime behaviour of Go, Git or the OS stated from memory. It is
    UNVERIFIED.
  - **Testimony I did not verify:** the facilitator-reported focused, full,
    six-package race, vet, Windows cross-build and shared-volume results; the six
    removed-protection overlays; and the 401-file manifest
    `35558a3f60bc892c15d68cfb063dec9b94ad54e8cbc75573667fec80a9f7cd7a`.
- Read in full:
  - `internal/trajectory/{unchanged.go, unchanged_test.go, reconcile.go, continuation.go, parent_recovery.go, snapshot.go, state.go, captured.go, trajectory.go}`
  - `internal/app/{trajectory_unchanged.go, trajectory_unchanged_test.go, trajectory.go, trajectory_reconcile.go, trajectory_verify.go}`
  - `internal/telemetry/record.go`
  - `internal/runner/{launch_budget.go, telemetry.go, failclass.go, durablekill.go}`
  - `internal/budget/{cycle_observer.go, cycle_session.go, reservation_receipt.go, resource_guard.go}`
  - `internal/evidence/tree.go`, the implementation note and the relevant parts of
    FINAL.md, IMPLEMENTATION.md and `docs/agent-runtime-configuration.md`.
- Read in part: `runner/{runner.go, launch.go, protocol_context.go, supervision.go, cycle_budget.go}`,
  `telemetry/usage.go`, `procctl/procctl_unix.go`, `driver/{checks.go, impl.go}`,
  `app/driver_impl.go`, `budget/cycle_binding.go`, `protocol/workspace.go`, `fsutil`.
- Not claimed: current-tree acceptance, Windows runtime, real-model behaviour,
  the helper-ticket/orphan-reservation recovery slice, packet/pilot experiments,
  quorum decisions.

## Refutation attempts

**R1: Original charge, action and epoch binding. Held.**

- Tried: swap in another ledger entry, action identity, epoch or invocation, or
  rehash a preview around a different attempt.
- Every guarded read checks the charge against the live ledger:
  - `state.go:297` `if err := a.Charge.Check(ledger); err != nil {`
  - `reservation_receipt.go:36-37` compares scope, the `startedAt` epoch, entry,
    kind and `ReservedAt`, plus
    `actionIdentityDigest(entry.Action) != r.actionSHA256`.
- The preview pins the whole canonical attempt (charge, sources, archives, launch,
  terminal):
  - `unchanged.go:67` `data, err := canonical(a)`
  - `unchanged.go:71` `u.AttemptSHA256 != digest(data) || u.InvocationID != a.Launch.InvocationID`
  - The resolution charge key is checked at `continuation.go:238`
    `p.ChargeKey != s.Attempts[i].Charge.EntryKey`.
- The launch identity is recorded under the guard against the live session charge:
  `state.go:547` `if a.Charge.EntryKey != charge.EntryKey || a.Launch != nil || ...`.

**R2: Original archived scope and complete archive integrity. Held.** The quorum
aspect is F3.

- Tried:
  - editing the current prompt to match the policy while the archive differs;
  - a truncated footer or corrupt remainder;
  - a symlink or oversized member;
  - releasing selected bytes before validation.
- Scope comes from the before-archive:
  `unchanged.go:81` `readSnapshotMember(ctx, snapshotDirectory(b), a.BeforeArchive, a.Before, name, 1<<20)`.
- Bytes are returned only after the full inspection succeeds:
  - `snapshot.go:527` `if err := inspectSnapshotFileWithMember(...); err != nil {`
  - That inspection re-encodes the canonical tar and checks the archive hash,
    whole-tree digest and links (`snapshot.go:678-686`).
  - It also checks file stability (`687-690`), refuses symlink selection
    (`660-661`) and enforces the member bound (`634-638`).
- `ScopeSHA256` is `digest(raw)` of the archived bytes (`unchanged.go:111`).
- Each archived criterion must match the frozen policy name and command hash
  (`unchanged.go:100`). The current parsed contract must equal the archived one
  (`unchanged.go:108` `if !sameJSON(current, original) {`).

**R3: Timestamp and order. Held by reading.** None of the order clauses is tested
in isolation (F2).

- Tried: charge-before-request, request-before-charge, start before charge,
  trajectory terminal before start, telemetry completion before the trajectory
  terminal, and an invocation older than the policy epoch.
- The clauses are at `unchanged.go:150`:
  - `terminal.RequestedAt.Before(s.Policy.StartedAt)`
  - `terminal.StartedAt.Before(a.Charge.ReservedAt)`
  - `a.Terminal.At.Before(*terminal.StartedAt)`
  - `terminal.CompletedAt.Before(a.Terminal.At)`
- The production order I traced independently satisfies all of them in both
  charge orders:
  1. requested record (`runner/telemetry.go:94`);
  2. cycle charge and `trajectory.Begin` (`launch_budget.go:70,77`);
  3. process start;
  4. `l.trajectory.Finish` (`telemetry.go:179-181`);
  5. `l.invocation.Finish` (`telemetry.go:191`).
- These are wall-clock comparisons, so a backwards clock step can only cause a
  false refusal.

**R4: Ordinary zero and nonzero exit qualification. Held.** F5 limits the provider
claim.

- The two accepted outcomes:
  - `unchanged.go:148` `completed := o.Status == "process-exited" && o.FailureClass == nil && *o.ExitCode == 0`
  - `unchanged.go:149` `failed := o.Status == "failed" && o.FailureClass != nil && *o.FailureClass == "process_failure" && *o.ExitCode > 0`
- The telemetry outcome is cross-checked against the retained trajectory terminal
  (`o.Status != a.Terminal.Status || *o.ExitCode != *a.Terminal.ExitCode`).
- The runner passes the same pre-override status and exit to both records
  (`telemetry.go:179-191`). A later `trajectory_failure` override therefore causes
  a status mismatch and a refusal.

**R5: Rejected incomplete, signalled, timed-out and provider-failed launches.
Refused by reading.** The consequence is F1; the test gap is F2.

- Charged but never started: `unchanged.go:141` requires `terminal.StartedAt != nil`
  and `a.Terminal.ExitCode != nil`.
- Signalled: exit -1 fails `> 0`. RECALL: Go reports -1 for a signalled process.
- Hard timeout:
  - `supervision.go:160-163` `case <-done: kill(); <-waitErr; return ctxErr()`
  - maps to `"timeout"` at `telemetry.go:155-156`.
- Watchdogs: `no_first_output` and `stalled` (`telemetry.go:151-154`).
- Structured provider failures: `telemetry.go:159-160`.
- Unobserved handoff: status mismatch; handoffs also return before any charge
  (`launch_budget.go:66-68`).
- Unavailable capture: `unchanged.go:45` (`SnapshotError != ""`, nil after-state
  or archive).

**R6: Gate and streak, including regression → unchanged → regression. Held.**

- An unchanged observation does not reset the streak; only a regression advances
  the trigger:
  - `continuation.go:69-72` `case Inconclusive: if r.Preview.Unchanged == nil { h.Decision.Consecutive = 0 }`
  - `continuation.go:79-80` `if a.Outcome == Regression && h.Decision.Consecutive >= 2 { h.RequiredReviewSequence = i + 1`
- Walkthrough:
  - R(1): consecutive 1.
  - U(2): consecutive stays 1; inconclusive pending [2]. The next charge is refused
    (`state.go:395-396`) until an attended `--acknowledge-inconclusive`
    (`continuation.go:197-206`).
  - R(3): consecutive 2, required review 3. The next charge is refused until an
    attended `--acknowledge-review`.
  - U(4) after acknowledging through 3: the review is not reopened.
  - R(5): required review 5 exceeds the acknowledgment, so the review opens again.
- A changed-source inconclusive still resets the streak, matching `Evaluate`
  (`trajectory.go:292-293`).
- An unchanged observation can only be the fixed all-unresolved `Inconclusive`
  assessment (`unchanged.go:71`), so it cannot launder a regression.
- For states without the new field the algorithm is identical. Retained
  continuation histories revalidate through the prefix recomputation at
  `continuation.go:263-275`.
- This is consistent with FINAL D6: "Repeated unchanged criticism is not a new
  regression."

**R7: Original criteria and quorum matching. Criteria held; quorum is F3.**

- Count: `unchanged.go:89`. Order, name and command hash: `unchanged.go:100`.
- Both archived and current scope go through one parser, which trims name and
  command like `driver.ReadChecksContract` does (`reconcile.go:191` vs
  `checks.go:128-129`). Every divergence I traced can only refuse.

**R8: Exact historical replay and parent recovery. Held.**

- Replay compares every fact except the pre-transition state hash:
  - `reconcile.go:329-330` `func compareReconciledParent(actual, original ReconciliationPreview) error { actual.StateSHA256 = original.StateSHA256`
  - It requires the retained digest (`unchanged.go:229`) and only re-syncs
    (`unchanged.go:234`).
- `PreviewUnchanged` returns the retained preview and refuses omitted sequences
  (`unchanged.go:190-197`).
- Every guarded read rederives the evidence:
  `reconcile.go:360-361` `if p.Unchanged != nil { return unchangedPreview(ctx, b, s, p.Root, p.Sequence) }`.
- Parent recovery never reconstructs unchanged evidence:
  `parent_recovery.go:267` `if err != nil && old.Unchanged == nil && old.Root == root && ...`.
  All other resolutions go through normal comparison (`parent_recovery.go:284-289`).
- Consequence, by design: losing the gitignored invocation records, or moving the
  root recorded in `p.Root`, makes every guarded read refuse. This is the same
  stance as for parent evidence.

**R9: Persistence and output failure. Held.**

- A persist error surfaces whether or not the write landed (`unchanged.go:248-250`).
  A retry after a landed write takes the replay branch.
- After an output failure the CLI exits 1 with retry guidance
  (`trajectory_unchanged.go:39-41`). An exact retry returns the same digest.
- The publication-interruption test (`unchanged_test.go:291-330`) matches this by
  reading.

**R10: Concurrent apply. Held by reading.**

- Preview and apply both run inside the resource guard (`state.go:453-459`).
- A second applier sees the resolution and must match its digest
  (`unchanged.go:224-231`).
- A new charge (`state.go:395-396`) and `Continue` (`continuation.go:123-124`) are
  both refused while any attempt is unreconciled.
- Guard acquisition waits at most 30 s (`state.go:453`), so a slow holder yields a
  retryable error, not divergence.

**R11: Bounded resources and malicious or corrupt metadata. No unbounded read
found.** The cumulative cost is F6.

- Runtime records: 16 MiB regular-file bound, `SameFile`, unknown fields refused,
  canonical re-encode required (`reconcile.go:79-110`).
- Scope member: 1 MiB (`unchanged.go:81`).
- Archive, member and entry bounds: `snapshot.go:27-29, 608`.
- Telemetry files are immutable once created:
  `record.go:323` `os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)`.
- The CLI bounds the sequence to `1..MaxPatches` (`trajectory_unchanged.go:21`).

**R12: Mixed provenance and separation from the verifier path. Held.**

- An unchanged preview must carry no verifier fields (`unchanged.go:56`): empty
  `RunID`, `ParentSHA256` and `RecoverySHA256`.
- The parent branch requires its own bindings (`continuation.go:245-246`).
- `CapturedRequest` still refuses equal trees:
  `captured.go:129-130` `if a.Before.Tree.SHA256 == a.After.Tree.SHA256 { return CapturedRequest{}, errors.New("unchanged source is not a new patch regression")`,
  and again when hashing (`captured.go:90`).

**R13: Attendance boundary. No bypass found.**

- `reconcile-unchanged` is dispatched before the attendance probe
  (`app/trajectory.go:17-18`), as the existing parent `reconcile --yes` path also
  is (`trajectory_reconcile.go:48-66`).
- `continue --yes` still requires attendance (`trajectory_reconcile.go:70-71`).
- Publication only turns "unreconciled" into "inconclusive pending", which still
  blocks charges (`state.go:395-396`).

## Findings

### [MAJOR] F1: Abnormal or refused-after-charge unchanged attempts permanently halt the trajectory

**What is wrong (PRIMARY).** A charged attempt whose after-tree equals its
before-tree but whose lifecycle is not an ordinary exit cannot be reconciled by
either path:

- The unchanged path refuses it at `unchanged.go:141` or `unchanged.go:150`. The
  comment at `unchanged.go:145-147` says: "Retain them for explicit process
  recovery instead of inferring inactivity from silence."
- The patch-verification path refuses equal trees (`captured.go:129-130`, quoted in
  R12).

Once such an attempt exists:

- Every further fixup is refused:
  `state.go:395-396` `if h.Unreconciled > 0 || ... { return errors.New("trajectory awaits independent reconciliation or an attended review decision; further fixup is refused")`.
- Continuation is refused: `continuation.go:123-124`
  `errors.New("continuation requires every charged attempt to be reconciled")`.
- Completion escalates: `state.go:500-501` via `driver/impl.go:316-317`
  `return ActionEscalated, c, fmt.Errorf("patch trajectory blocks completion: %w", err)`,
  and `app/driver_impl.go:486`.

The trajectory CLI has no other control. Its complete dispatch is
`app/trajectory.go:17-33`, with usage at `40-43`. I did not audit every budget
command for a way to detach a frozen trajectory policy.

**Counterexample A: hard timeout or rate limit with no edits.**

1. A fixup launch hangs or fails without touching files.
2. The runner cancels and kills the process group
   (`procctl_unix.go:40` `syscall.Kill(-target, syscall.SIGTERM)`, then `:48`
   SIGKILL).
3. Telemetry records `failed`/`timeout` (`telemetry.go:155-156`). For a structured
   rate limit it records `failed`/`rate-limit` (`telemetry.go:159-160`).
4. The trajectory terminal gets the same status and exit (`telemetry.go:179-181`),
   and the after-tree equals the before-tree.

Result: `reconcile-unchanged` refuses at `unchanged.go:150`, `trajectory verify`
refuses at `captured.go:130`, and the trajectory is halted.

The same attempt with a single changed byte would still be verifiable, because
`capturedRequestAt` never checks the implementer's status (`captured.go:126`).

**Counterexample B: refused after the charge was spent.**

- *Protocol-context refusal.* `beginProtocolLaunch` computes the refusal first
  (`protocol_context.go:118`
  `prepared, protocolContext, contextErr := prepareProtocolPrompt(root, prompt, info)`).
  It then charges the cycle and records `trajectory.Begin` (`:121`
  `evidence, err := beginLaunch(ctx, root, runID, agent)`), and only then finishes
  with the refusal (`:125-126`). The trajectory terminal is `failed` with a nil
  exit, so `unchanged.go:141` refuses.
- *Budget refusal.* `ChargeStep` (`launch_budget.go:92`), the monetary binding
  (`:111`) and the launch reservation (`:139`) all run after the cycle charge and
  `trajectory.Begin` (`:70,77`). A refusal there is finished at
  `telemetry.go:105-106`, with the same outcome as above.
- *Before spawn.* A command-build, log-open, context or `cmd.Start` failure
  (`runner.go:1066-1107`) has the same outcome.

**Why it matters.** Rate limits and timeouts are routine; the previous attempt of
this very review was lost to a weekly limit. FINAL D6 says "A failed charged
attempt stays spent." Reading that as "the loop may continue within budget" is my
inference.

The known gaps do not name this class explicitly:

- The docs (`agent-runtime-configuration.md:1384-1386`) say such cases "remain
  unresolved".
- IMPLEMENTATION.md (`:333-335`) says "incomplete lifecycle/capture cases remain
  unresolved". A timeout is a complete lifecycle.
- The planned integration slice reported by the facilitator covers reservations
  interrupted *before* trajectory publication. Here the terminal *was* published.
- Nothing states that the consequence is a permanent halt with no operator path.

**Fairness.** The note and IMPLEMENTATION.md say unchanged attempts had no
reconciliation at all before this slice, so this is not a regression. That rests on
their testimony; I did not inspect the prior commit.

**Suggested fix.** Name these classes and the halt as an explicit open obligation
with an owner. Then choose one:

- an attended acknowledgment control for abnormal unchanged attempts that retains
  the lifecycle hashes, keeps the charge spent, grants nothing beyond the remaining
  budget, and still requires a later independently verified clean patch for
  completion; or
- avoid spending a cycle charge on a launch whose refusal is already known before
  `beginLaunch`, as in the protocol-context case, provided requested evidence is
  still written.

### [MAJOR] F2: The contradictory-lifecycle tests cannot detect removal of the qualification predicate

**What is wrong (PRIMARY, tests read, not run).**
`TestUnchangedReconciliationRejectsMissingAndContraryEvidence`:

1. builds one `exit 0` attempt (`unchanged_test.go:189`
   `r := unchangedProcess(t, root, b, "exit 0", false)`);
2. takes its preview digest `p`;
3. rewrites only `terminal.json` per mutation (`:209-222`: pid, identity, run,
   before-start, no-start, no-exit, signal, wrong-exit, timeout, cancelled,
   provider, handoff);
4. asserts only that this call errors:

```go
_, err := ReconcileUnchanged(ctx, root, "fixture", 1, p.SHA256())
...
if err == nil {
	t.Fatal("contrary process lifecycle accepted")
}
```

Any byte change to `terminal.json` changes `TerminalSHA256`, so the apply always
refuses at `unchanged.go:240`
`if sequence != len(s.Resolutions)+1 || p.SHA256() != expected {`. That holds even
if the status, failure-class, exit-sign and time-order clauses at
`unchanged.go:148-150` were deleted.

No other test covers this:

- The positive scenarios only use real `exit 0`, `exit 7` and empty-commit
  processes (`unchanged_test.go:108`; `trajectory_unchanged_test.go:82,105,144,196,215`).
- No test calls `unchangedPreview`/`PreviewUnchanged` with a contrary lifecycle.
  My grep of `*_test.go` found only the missing-file loop at
  `unchanged_test.go:196-206`.

Mutants I expect to survive every listed test (by reading, not executed):

- (m1) drop `*o.FailureClass == "process_failure"`;
- (m2) relax `*o.ExitCode > 0` to `!= 0`;
- (m3) drop `terminal.StartedAt.Before(a.Charge.ReservedAt)`;
- (m4) drop `terminal.RequestedAt.Before(s.Policy.StartedAt)`;
- (m5) drop `a.Terminal.At.Before(*terminal.StartedAt)`;
- (m6) delete the whole `(!completed && !failed)` qualification, while keeping the
  nil checks at `:141` so the nil-pointer subtests do not panic.

The note claims coverage of "missing or contradictory lifecycle" (note `:80-83`).
The six reported overlays do not name these clauses (note `:91-94`).

**Why the guard is load-bearing.** `KillGroup` sends SIGTERM and escalates to
SIGKILL only after 1.5 s (`procctl_unix.go:40-48`). An agent that handles SIGTERM
and exits with a positive code on a hard timeout therefore yields runner outcome
`failed`/`timeout` with that code. The trajectory terminal gets the same
status and exit (`telemetry.go:179-181`), so only m1's clause refuses it. RECALL:
Go's `ProcessState.ExitCode()` reports the process's own code after a handled
signal.

**Suggested fix.**

- Call `PreviewUnchanged` for each mutation, or assert the specific lifecycle
  error.
- Add consistent fixtures using the existing helper, which already writes both
  records (`unchanged_test.go:98-101`): `failed` plus `timeout`, `cancelled`,
  `stalled` or `provider-error` with a positive exit; `failed` plus
  `process_failure` with exit -1.
- Apply time-order mutations consistently across requested, started and terminal.
- Add overlays for m1–m6.

### [MINOR] F3: "Original quorum" is not bound to anything frozen at activation

**What is wrong (PRIMARY).**

- The frozen policy has no participants field (`state.go:27-35`: Version, Idea,
  Scope, StartedAt, Implementer, Baseline, Criteria).
- `unchangedScope` compares the quorum in *this attempt's* before-archive only with
  the current file (`unchanged.go:89-98, 104-110`).
- For attempt N>1, the before-archive is the previous attempt's after-archive or
  its continuation archive (`state.go:325-327`; the continuation override is at
  `state.go:289-291`). That previous attempt was an implementer-authored charged
  patch, and fixup processes run in the worktree root (`runner.go:1079`).
- Parent verification reads the quorum from the *current* prompt
  (`trajectory_verify.go:86-87`
  `if !found || !slices.Contains(current.Participants, implementer) || !slices.Contains(current.Participants, verifier) {`),
  pins it in the request (`:200`), and reconciliation compares the current prompt
  with that request (`reconcile.go:129`).

**Counterexample.**

1. Attempt 1 edits `source` and rewrites `participants: [builder, reviewer]` to
   `[builder, other]`, leaving `checks:` untouched so the policy hashes still match.
2. `trajectory verify --verifier reviewer` is now refused (`trajectory_verify.go:87`).
   `--verifier other` is accepted if `other` is a configured participant CLI.
3. Reconciliation accepts.
4. An unchanged attempt 2 reads `[builder, other]` from both the archive and the
   current file and is accepted.

The docs call this the original quorum (`agent-runtime-configuration.md:1377-1379`:
"Original quorum and named criterion commands come from the retained before-source
archive"). §5 says quorum "locks once Phase 0 completes".

**Scope and impact.** For the unchanged path the effect is a mislabelled
observation. The verifier-selection weakness is pre-existing in the parent path
(LE-3 diversity) and outside this slice.

**Suggested fix.** Either pin the quorum from the baseline archive, i.e. the
activation-time prompt, or record the participants in a future policy version.
Until then, qualify the docs wording.

### [MINOR] F4: Material equality can miss untracked files hidden through non-tracked excludes (premise RECALL, UNVERIFIED)

**What is wrong.** The only gate that lets an unchanged attempt skip independent
verification is tree-digest equality:
`unchanged.go:48-49`
`if a.Before.Tree.SHA256 != a.After.Tree.SHA256 { return Attempt{}, errors.New("changed material source requires independent patch verification")`.
Commit and status differences are deliberately tolerated.

The inventory and status honour Git's standard excludes:

- `snapshot.go:70` `gitOutput(ctx, root, "ls-files", "-c", "-o", "--exclude-standard", "-z")`
- `state.go:129` `gitOutput(ctx, root, "status", "--porcelain=v1", "--untracked-files=all")`
- The documented intent is "honoring .gitignore so dependency/cache directories
  stay out of scope" (`evidence/tree.go:155-157`).

**Premise (RECALL).** `--exclude-standard` and default `git status` also apply
`$GIT_DIR/info/exclude` and `core.excludesFile`, neither of which is tracked
material.

**Counterexample (UNVERIFIED).**

1. During a fixup the implementer creates an untracked `internal/x/zz_init.go`
   that affects every build of that package.
2. It appends that path to `.git/info/exclude`.
3. The digest and status are unchanged, so the attempt is reconcilable as
   unchanged with no verification.
4. The file persists invisibly in the live worktree but is absent from the archives
   that captured verification restores.

Editing the tracked `.gitignore` would change the digest, so only untracked
exclude sources form this channel.

**Scope and impact.** Pre-existing and shared with `evidence.TreeDigest`. The
outcome stays inconclusive and attended.

**Suggested fix.** Bind the exclude inputs (`info/exclude`, the effective
`core.excludesFile`) into the observation, or refuse when they differ from the
baseline. Alternatively, state the limit next to the "records material equality"
claim.

### [NIT] F5: "Provider failures refuse" holds only for structured provider failures

- The collector only parses when launch args are structured:
  `telemetry.go:99` `structured := metadata.LaunchMode != agents.LaunchACP && telemetry.StructuredArgs(args)`.
- It recognizes a failure only from JSON events:
  `usage.go:153` `if event["is_error"] == true || kind == "turn.failed" || kind == "error" {`.
- Otherwise a nonzero exit becomes `process_failure` (`telemetry.go:161-162`) and
  qualifies under `unchanged.go:149`.
- The regex classifier in `failclass.go` feeds run events (`runner.go:747`), not
  telemetry.

So a plain-text rate-limit or auth error with exit 1 is published as an "ordinary
process failure". The docs (`:1384-1386`) and note (`:38-39`) state the refusal
without this qualification. There is no acceptance impact.

Suggested fix: qualify the wording to "classified (structured) provider failures".

### [NIT] F6: Rederiving every unchanged resolution adds guarded work under the runner's 30 s publication ceiling

- Every guarded read, including the runner's `Finish` (`state.go:574` via `withState`,
  whose check runs at `state.go:487`), repeats this for each unchanged resolution
  (`reconcile.go:343-364`):
  - a full inspection of its before-archive, up to
    `MaxSnapshotBytes int64 = 256 << 20` (`snapshot.go:27`; `unchanged.go:81`);
  - three record reads of up to 16 MiB each (`unchanged.go:140,153,157`);
  - a `LoadCycleBinding` (`unchanged.go:120`).
- This comes on top of `checkSourceSnapshots` (`state.go:627-651`).
- The runner gives the trajectory publication 30 s in total
  (`telemetry.go:180` `context.WithTimeout(context.Background(), 30*time.Second)`).
- On failure it records `trajectory_failure` and leaves the attempt without a
  terminal (`telemetry.go:182-185`). Neither path accepts that (`unchanged.go:45`;
  `captured.go:126`).

The ceiling is pre-existing; this slice adds per-resolution cost. I did not measure
the magnitude, so it is UNVERIFIED.

Suggested fix: within one guarded read, avoid re-inspecting archives that were
already verified, or decouple the post-exit publication ceiling from history size.

## Remaining uncertainty

- **Executed evidence.** I executed nothing. Every "surviving mutant" and
  "counterexample" statement above is derived from source reading.
  Facilitator-reported results remain testimony.
- **Trust boundary.** Hashes detect mismatch, not coordinated same-UID rewriting of
  state, ledger, telemetry and archives, as the note discloses (`:66-70`).
- **Windows.**
  - `syncUnchangedState` (`unchanged.go:260-273`) opens read-only handles and calls
    `SyncFile`, which is `file.Sync()` on every non-darwin platform
    (`sync_other.go:1-8`).
  - The state writer, by contrast, is platform-split (`replace_windows.go:14-27`).
  - Whether exact replay or concurrent convergence succeeds on Windows is
    UNVERIFIED. RECALL suggests flushing a read-only handle fails there.
  - Whether a trajectory can be activated on Windows at all is also UNVERIFIED.
  - The unchanged tests skip on Windows (`unchanged_test.go:25-27`;
    `trajectory_unchanged_test.go:26-28`).
- **Idea path.** The scope readers hard-code
  `parley-deck/ideas/<idea>/00-prompt.md` (`unchanged.go:80`; `reconcile.go:150`).
  Cycle policies carry an `IdeaPath` constrained only to start with `parley-deck/`
  (`cycle_binding.go:74`), and the driver passes `d.cfg.IdeaDir`
  (`driver/cycle_budget.go:76`). I did not establish whether a trajectory-enabled
  binding can use a different path.
- **Launch modes.** The unchanged path does not check `LaunchMode`, unlike the
  verifier path (`reconcile.go:260`). I did not inspect ACP stop and exit semantics
  (`acp.go:262`). Whether ACP or interactive fixups ever qualify is unknown; from
  what I saw the risk is false refusal only.
- **Torn reads.** `readReconciliationJSON` accepts both the newline-less and the
  newline encodings (`reconcile.go:106`), while telemetry writes only the newline
  form, non-atomically (`record.go:323-327`). Pinning a torn, newline-less prefix
  would require both preview and apply to observe that exact prefix. I judge this
  negligible; not measured.
- **Design question, not a finding.** A latest unchanged attempt blocks
  `RequireResolved` even when its bytes equal a previously verified clean patch
  (`state.go:503-505`; documented at `agent-runtime-configuration.md:1407-1408`).
  At the cap this forces escalation. That is a fail-closed choice for the
  operator and participants.
- **Quorum parsers.** Verification-time parsing (`workspace.go:394-410`, `parseList`)
  and reconciliation-time YAML parsing (`reconcile.go:179-189`) differ. Every
  divergence I traced only refuses; this was not exhaustive.
- **Out of scope.** Real-model launch surfaces, independent concurrency and closure
  evidence, the helper-ticket/orphan recovery slices, the packet and pilot
  experiments, historical quorum decisions, owned signatures and final report gates
  are not addressed by this source-only review.
