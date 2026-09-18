---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-15
base-commit: 39107b138e77de064f21fec4aafa6073887c2783
kind: native-candidate
status: candidate-awaiting-independent-review
---

# R1/R2 native candidate: recovery releases the cycle guard during history validation

## Summary

- **R1 is implemented as a candidate.**
  - `recover-reservation` preview and apply now capture the complete structural authority under
    the guard, then release it.
  - The full history is validated without the guard.
  - Publication runs only after an exact guarded recheck.
  - A real `Begin`-derived `Run.Finish` succeeded while a recovery content check was paused. The same
    fixture fails on HEAD's recovery logic.
- **R2 is implemented.** The test failure now prints the content-free terminal JSON.
- **Retained results.**
  - One control passed (a coverage gap).
  - One existing test fails in this environment on HEAD too.
- **This is my implementer claim, not independent acceptance.**
  - R1 and R2 stay open until a non-owner reviews them.
  - R3 and the N1 cross-review and step residuals remain open.

## Scope and provenance

- **Launch:** claude-1, `context_mode=full`,
  `source_sha256=4519258c96a45515518f44d29f769a5510e32924d1e27ebcf6d04cf554b1937a`.
- **Environment:** native darwin/arm64, `go1.27.1`.
  - Go settings: `GOPROXY=off`, `GOFLAGS=-mod=readonly`, `GOTOOLCHAIN=local`.
  - `GOCACHE`, `GOTMPDIR` and `TMPDIR` are under `.parley-runtime/work/r1r2/`.
- **Not done:** no Git mutation, network access, other agents or model calls, packet or pilot trial, race detector or full suite.
- **Unedited:** no historical review, note or signature.
- **Two kinds of claim.** "Observed" means an executed result, cited by its log below. "By reading" means source reasoning that I did not execute.

## Files

| Change | Path |
| --- | --- |
| modified | `internal/trajectory/reservation_recovery.go` |
| new | `internal/trajectory/reservation_recovery_validation_test.go` |
| modified | `internal/app/agents_exec_test.go` |
| new | `parley-deck/ideas/meta-protocol-change-evidence-first-efficiency/implementation-notes/claude-1-reservation-validation-candidate-20260915.md` (this note) |

**Scratch files** exist only under `.parley-runtime/work/r1r2/`: `run.py`, `head/`, `overlays/`,
`logs/`, `results-*.json` and the Go cache.

**Side effect outside scratch:** the code under test creates its normal budget lock files under
`~/Library/Caches/parley/budget-locks/`, which a control's failure message shows. Existing fixtures do the same. I did not change that location.

## R1: what changed

The public entry points and the CLI are unchanged. `recoverReservation` keeps its signature,
because `reservation_recovery_test.go` calls it, and delegates to the new unexported
`recoverReservationChecked(..., check, persist)`. Both public entry points bind `check` to the
real `checkStateSnapshots`.

**Guarded observation.** `withReservationRecovery` holds the common cycle guard while it runs
every check the old body ran before its content read, in the same order:

1. load and reload the binding;
2. inspect the ledger;
3. read the bounded canonical state;
4. open and read the intent;
5. check the root identity;
6. run `validateIntentBefore` and `compareIntentPrefix`;
7. for a charged row, run `intentAttempt`, then either:
   - **missing row:** require `digest(raw) == BeforeSHA256` and ledger entries = attempts + 1;
   - **existing row:** check the row's identity;
8. run `validateState` on the history that publication would keep (a missing row is already appended).

It reads no archive or resolution content. The identity argument check still runs once, before any of this.

**Two stages.**

1. **Capture** the observation under the guard, then release the guard.
2. **Content check:** `check` runs without the guard. That is `checkStateSnapshots`: intents,
   every archive and every resolution.
3. **Exact guarded recheck:** repeat the full guarded observation and require all of these to be
   equal to the capture: store directory and scope, complete policy JSON, complete ledger JSON,
   exact state bytes, the missing flag, the candidate state and the preview digest.

Only after that equality, still under the same guard:

- apply's no-charge refusal;
- the expected-preview check;
- `syncIntent`;
- then either the missing-row `persist` or the existing-row durability sync.

**Retry.** Exactly one full restart from stage 1 is allowed, and only when the stage-3 equality
fails. It happens before any publication. A second inequality refuses with
`reservation recovery authority changed during full evidence validation`.

**Never retried:**

- identity errors;
- structural errors in either observation, including a disappeared binding, state or intent;
- content-check errors;
- a wrong expected preview;
- intent-sync or publication errors;
- context errors.

**Why the restart exists (observed).** Without it, a paused exact apply fails after a concurrent
apply publishes the missing row. The old guard serialized that case into the replay branch. See
the `no-drift-retry` control.

**Seam disclosure.** `check` is a parameter of an unexported function. It mirrors the existing
`persist` parameter and `withStateResolutionCheck`'s `check`. It is not a mutable hook or a
test-only branch. The canonical tests pass a wrapper that first runs the real
`checkStateSnapshots` and then pauses. No production hook exists; overlays are used only by the controls below.

## R1: new canonical tests

### `TestReservationRecoveryContentCheckDoesNotBlockLiveFinish` (decisive witness)

1. A real `ChargeCycle` creates the charged row, and the real `Observer` publishes its intent.
2. A real `Begin` handle publishes the launch.
3. The public `PreviewReservationRecovery` previews the live row.
4. A real child process changes the source and exits 7.
5. The preview content check is paused.
6. `run.Finish` runs with a 2 s context, then the pause is released.

Assertions:

- `Finish` succeeded;
- the preview restarted exactly once, and its second content check saw the terminal;
- the preview digest is unchanged;
- exit code 7 and the after-state are retained;
- an exact apply replay leaves the state bytes unchanged.

### `TestReservationRecoveryConcurrentApplyReplaysAfterContentCheck`

The fixture has a missing row. A real `CycleBinding.Reserve` publishes the intent and the charge;
the test observer then returns an error from `AfterCycle`, so the charge stays spent and the row
is missing.

1. Apply A pauses in its content check.
2. Apply B runs with a 2 s context and publishes exactly once. Inside B's `persist`, a 50 ms guard
   acquisition fails, so publication holds the guard.
3. A restarts once and takes the replay branch.

Assertions:

- A called `persist` zero times and returned the same digest;
- state bytes are unchanged after A;
- full `Inspect` passes;
- no launch or terminal was invented.

### `TestReservationRecoveryRechecksAuthorityAfterUnguardedContent`

Eight subtests, each applying against a missing row. In every content check, a guard acquisition
must succeed within 1 s.

| Subtests | Injected change | Asserted result |
| --- | --- | --- |
| policy, ledger (single) | `ExtendCycleBudget`; ledger `Settle` | one restart, one publication, same digest, full `Inspect` passes |
| policy, ledger, intent (repeated) | extension twice; `Settle` then `ReconcileUnknown`; intent rewritten with `PreparedAt` one nanosecond earlier each time (still structurally valid) | drift refusal after exactly 2 checks, 0 publications, state bytes unchanged |
| state, intent, binding removed | the file or directory disappears | refusal with no restart: 1 check, 0 publications |

**Fixture disclosure.** `failedAfterObserver` wraps the real `Observer` and replaces only
`AfterCycle`. It is used only to create the missing-row state. The decisive `Finish` witness uses
the unmodified observer and a `Begin`-derived handle. No handle is constructed by the fixture.

## Commands and decisive results (observed)

Runner: `.parley-runtime/work/r1r2/run.py`. It runs every command through `subprocess.run` with
`stdout=PIPE, stderr=PIPE`, writes a uniquely named log only after the command ends, and records
tested-file hashes before and after each stage. Every stage reports `sources_unchanged: true`.
All logs are under `.parley-runtime/work/r1r2/logs/`.

### 1. `python3 .parley-runtime/work/r1r2/run.py candidate`

Results: `results-candidate-20260915T230717.json`.

| Command | Result |
| --- | --- |
| `gofmt -l` on the three changed Go files | exit 0, no output |
| `go vet ./internal/trajectory/` | exit 0, 4.191 s |
| `go test -count=1 -v -run '^(TestReservationRecovery\|TestStateValidation\|TestPendingCycleRefusesBeforeHistoricalContent)' ./internal/trajectory/` | **exit 1**, 57.482 s |
| `go test -count=1 -v -run 'TestAgentsExecRecordsManualLaunch\|ReservationRecovery' ./internal/app/` | exit 0, 11.988 s |

**Trajectory run (exit 1):**

- New tests PASS: live `Finish` 1.48 s, concurrent apply 1.18 s, drift table 7.57 s (all 8 subtests).
- Existing recovery tests PASS: crash 2.39 s, refusal 7.63 s (10 subtests), publication failure 1.83 s, extension 1.12 s.
- Other existing tests PASS: `TestStateValidationDoesNotHoldControlGuard`, `RechecksPolicyLedgerAndExistence`,
  `PreservesFullChecksAndCallbackGuard`, `RetryIsBoundedAndRevalidatesEvidence`, and
  `TestPendingCycleRefusesBeforeHistoricalContent`.
- **FAIL:** existing `TestStateValidationDoesNotBlockRegisteredCriterionStop`, 21.39 s, at
  `state_validation_stop_test.go:43`: `registered criterion never started context deadline exceeded`.
  See stage 3.

**App run (exit 0):**

- `TestAgentsExecRecordsManualLaunch` PASS, 2.00 s.
- `TestTrajectoryReservationRecoveryCLIOutputFailureAndNoExecution` PASS, 1.07 s.
- `TestTrajectoryReservationRecoveryConcurrentProductionCLI` PASS, 2.68 s.
- `TestTrajectoryReservationRecoveryHistoricalReplayAfterLaterWork` PASS, 3.02 s.

### 2. `python3 .parley-runtime/work/r1r2/run.py controls`

Results: `results-controls-20260915T230940.json`.

Each recovery control runs `go test -count=1 -v -overlay <overlay.json> -run
'^TestReservationRecovery(ContentCheckDoesNotBlockLiveFinish|ConcurrentApplyReplaysAfterContentCheck|RechecksAuthorityAfterUnguardedContent)$'
./internal/trajectory/`. Overlay sources and `overlay.json` files are in
`.parley-runtime/work/r1r2/overlays/<name>/`.

| Control (overlay of `reservation_recovery.go`) | Result |
| --- | --- |
| `old-source-guarded`: HEAD logic (`7f62fbc6…`), with only its two guarded `checkStateSnapshots(ctx, *b, s)` calls replaced by the test's `check` | exit 1. `Finish` FAIL at `:110` `recovery preview content check blocked live terminal publication … budget lock contention exhausted the acquisition deadline … context deadline exceeded`. Concurrent apply FAIL at `:178` `recovery content check blocked a concurrent exact apply … context deadline exceeded`. Five drift subtests FAIL with `recovery content check holds the cycle guard`. The three removed-authority subtests pass; they do not test the guard. |
| `no-drift-retry` (`attempt == 0` changed to `attempt < 0`) | exit 1. `Finish` FAIL at `:113` (drift error, 1 check). Concurrent apply FAIL at `:182` `paused exact apply did not replay …`. Both single-drift cases FAIL at `:274`; the three repeated cases FAIL at `:270` with 1 check. |
| `no-policy-recheck` | exit 1. Only policy fails: single at `:274 … <nil> 1 1`; repeated at `:270 … <nil> 1 1`, meaning it published despite drift. |
| `no-ledger-recheck` | exit 1. Only ledger fails: single at `:274 … <nil> 1 1`; repeated at `:270 … <nil> 1 1`. |
| `no-state-recheck` (raw bytes, missing flag and state removed) | exit 1. `Finish` FAIL at `:113 … <nil> 1 false`. Concurrent apply FAIL at `:182 … <nil> 1 0`. The drift table passes. |
| `no-intent-preview-recheck` | **exit 0, all three tests PASS, 11.437 s. This control produced no counterexample.** |

**The old-source counterexample fails at the intended assertion, not by failing to compile.** The
adapter only moves HEAD's guarded content calls behind the test's `check`.

**Coverage gap from the passing control (by reading).** On a missing row, the appended attempt
records `ReservationIntentSHA256`, so the candidate-state comparison also catches intent drift. On
an existing row, the structural identity check refuses a changed intent. For an
`intent-without-published-charge` preview there is no row, so the preview-digest clause is the only
intent comparison. No test covers that case.

**R2 witness: `control-r2-failing-child`.**

- **Overlay:** `agents_exec_test.go` with only the fixture child changed to `cat >/dev/null; exit 3`. Command: `-run '^TestAgentsExecRecordsManualLaunch$' ./internal/app/`.
- **Result:** exit 1, FAIL 1.82 s at `agents_exec_test.go:46`:
  `code=1 stderr=agents exec: invocation failed; inspect private local evidence`,
  followed by `stdout={"schema_version":1,"type":"invocation.terminal",…,"outcome":{"status":"failed","exit_code":3,"failure_class":"process_failure","artifact_sha256":null,…},…}`.
- **Content:** the printed record holds identifiers, timestamps, PID, hashes and null usage. It holds no prompt text and no artifact bytes.

### 3. `python3 .parley-runtime/work/r1r2/run.py stop-diagnosis`

Results: `results-stop-diagnosis-20260915T231025.json`.

| Run | Result |
| --- | --- |
| `stop-isolated-candidate`: `go test -count=1 -v -run '^TestStateValidationDoesNotBlockRegisteredCriterionStop$' ./internal/trajectory/` | exit 1, FAIL 21.63 s at `state_validation_stop_test.go:43`: `registered criterion never started context deadline exceeded` |
| `stop-isolated-head-verbatim`: same command, with an overlay using HEAD `reservation_recovery.go` bytes and the new test file deleted, i.e. the exact HEAD trajectory package | exit 1, FAIL 21.93 s, same line and message |

**Observed:** the failure reproduces on the exact HEAD package in this environment, so this
candidate does not cause it. **Not established:** its cause. Codex-1's testimony reports a pass at
5.561 s elsewhere. A difference in this launch's sandbox or temporary-directory environment is an
untested hypothesis. The failure is retained, and this launch has no registered-`Stop` evidence.

## R2: what changed

In `TestAgentsExecRecordsManualLaunch`, the failure `Fatalf` now also prints `stdout=%s`, the
`--json` terminal record, with a two-line comment explaining why. Nothing else changed: the
fixture `timeout_ms`, every assertion and the branch structure are as before. The original
historical failure (parallel full suite, 6.28 s, generic invocation failure) stays unexplained.
Only a future failing run can classify it.

## Hashes

**HEAD** (`git show HEAD:<path>`):

| File | SHA-256 |
| --- | --- |
| `reservation_recovery.go` | `7f62fbc6288504941b1035e45c7f3b5ddd5a5ed0b9facc1c44e1fbd9c90b59db` |
| `agents_exec_test.go` | `bec450b19e45390eace35dc79472a2955e5f908c99074cf5bcf4a7e839bee0d6` |
| `reservation_recovery_test.go` | `5c37af49561d119ac715a923cac6c436b3d9b7f14807853791e4a1fdd4ea90ea` |
| `state.go` | `fc0d89cd164981c86e958c81584ae87bf8cb0b6f9bd0c1630e5cac43df682275` |
| `state_validation.go` | `eb09fa9990f4e5e0449762d8109d9622da9792619952a7c1da8db4a6b8539c7c` |

**Tested and frozen** (identical before and after all three stages):

| File | SHA-256 |
| --- | --- |
| `internal/trajectory/reservation_recovery.go` | `4f35b839bf8762c8a211c0efc7b60131fd28e16cdcdff0cf04a9c976b5b7ffbb` |
| `internal/trajectory/reservation_recovery_validation_test.go` | `c6fc513b89f501a8c278eaeae37119940bdfe9669890a5b1c0727f68e789a893` |
| `internal/app/agents_exec_test.go` | `28791fbbb50ea69b7aec58859cd1c9d795d8b6c219634f68e9ab40ab798a2732` |

`reservation_recovery_test.go`, `state.go` and `state_validation.go` are equal to HEAD.

**Runner:** `run.py` is `22092178f8ab94a3038e9f87ce854bcaff2ac897df9cc8148d6c5a19f0e83503`. The
`candidate` stage ran on an earlier revision that lacked only the later-appended `stop-diagnosis`
branch. Those earlier bytes were not hashed.

**Log SHA-256:**

| Log | SHA-256 |
| --- | --- |
| `trajectory-focused` | `baa61a7b…44e7e` |
| `app-focused` | `b7287154…0ccc3` |
| `old-source-guarded` | `f82692be…5166b` |
| `no-drift-retry` | `aad512d7…34419` |
| `no-policy-recheck` | `72dd2eb8…b2958` |
| `no-ledger-recheck` | `0ab700aa…b6bda` |
| `no-state-recheck` | `997629e4…b4df` |
| `no-intent-preview-recheck` | `9e01d3d3…03ef` |
| `r2-failing-child` | `892fba1c…dbb5` |
| `stop-isolated-candidate` | `a01e267f…9e5a` |
| `stop-isolated-head-verbatim` | `5c154758…ee30` |

Full digests are in the results JSON files.

## Limitations

1. **Uncovered intent drift.** No test covers intent drift on an `intent-without-published-charge`
   preview, and the preview-digest clause has no counterexample (see the passing control).
2. **No `Stop` evidence.** No registered-`Stop` contention fixture was written. The existing
   registered-stop test fails in this environment on HEAD too.
3. **Pause point.** The pause uses the unexported `check` seam. No overlay hook paused the public
   `PreviewReservationRecovery` itself. That the public entry points bind `checkStateSnapshots` is
   established by reading.
4. **Legitimate concurrent change (by reading, not exercised).** Such a change can land while
   `checkStateSnapshots` reads, because it re-inspects the ledger and intents without the guard.
   It then surfaces as a content-check error rather than drift, and is not retried. The old
   guarded path would have waited up to 30 s instead. `withValidatedState` makes the same trade.
5. **Equality assumption.** The recheck assumes the policy and `budget.Snapshot` JSON have no
   fields computed at read time. The struct definitions show none. I did not read the ledger
   `read` path.
6. **Timing.** The 2 s and 1 s windows are deterministic small-fixture witnesses. They are not
   production deadline guarantees or measurements with large histories.
7. **Not run:** race detector, app vet, shared-volume or Windows runs, full suite.
8. **Not established:** a fence against same-UID writers, process inactivity, or any general speedup.

## Next actions for a non-owner reviewer

1. Try to refute R1:
   - check the equality set and the retry bound;
   - check that no publication path runs on stale authority;
   - check the seam.
2. Decide whether to require an uncharged-intent drift case and a registered-`Stop` contention case.
3. Diagnose the environment-dependent registered-stop failure separately; it also fails on HEAD.
