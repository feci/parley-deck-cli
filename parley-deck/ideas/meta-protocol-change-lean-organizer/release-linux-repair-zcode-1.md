---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: participant implementer (Linux release-channel repair)
artifact: release-linux-repair
date: 2026-09-24
status: IMPLEMENTED — awaiting independent review + hosted ubuntu classification
base: 4e0c069 (lean-organizer branch; product source at base is byte-identical to v1.49.0)
---

# Linux release-channel repair — U1 fix + U2 diagnostics — zcode-1

Scope honored exactly: **U1 source fix with regression tests; U2 test-only
diagnostics at the known failure sites.** Not implemented here: U3 (stays
proposed/deferred — not needed for the current green channel), all Windows
W1–W10 (escalated to owner), any workflow change (none — `.github/` untouched),
any test skip/retry/filter/relaxation (none), any version bump (VERSION,
internal/app/version.go, CHANGELOG.md untouched). Closed FINAL /
IMPLEMENTATION / review / consensus artifacts untouched. v1.49.0 not moved;
any eventual publish of this product-source change is a new v1.49.1 per
protocol, decided by organizer/owner, not by this participant.

## U1 — `internal/acp/spawn.go`: drain stderr before reaping (Fix)

Root cause (verified in source, matching claude-1 Finding E and the repair
plan): `Stop` and `Wait` both called `p.cmd.Wait()` before the stderr copier
goroutine (`p.wg`) finished. `os/exec` closes the `StderrPipe` read end when
`Wait` returns, so a fast reap truncates the copier mid-`io.Copy` — hosted
ubuntu run 35987916696 recorded `observed bytes: 0` against 16384 written
while the exit code (7) was read correctly.

Lifecycle assessment (done before reordering, per dispatch — not mechanical):

- **Product usage audit.** `acp.Process` has exactly two product call sites:
  `internal/runner/acp.go:96` and `:257`, both `process.Stop(ctx)` with a
  2-second timeout; `kill()` additionally calls `Kill()`. Exported `Wait()`
  has zero product callers (only tests). So the reorder's risk surface in the
  product is `Stop` only.
- **Normal exit path.** The copier reaches EOF at the same kernel event
  (process exit closes the write end) that unblocks `cmd.Wait`. Draining
  first adds no latency and removes the truncation window entirely.
- **Grandchild holding the stderr write end.** A naive drain-first could
  block: EOF then arrives only when every writer exits. In `Stop` this is
  broken by the existing ctx-timeout branch: `killProcessGroup` (SIGKILL to
  the Setsid group) kills lingering writers, EOF follows, drain completes,
  then reap. To keep `Stop` bounded even when a writer escaped the process
  group (setsid grandchild), the timeout branch now also closes the stderr
  read end, which interrupts a copier still blocked in `Read`
  (pollable `*os.File`; pending-I/O unblock on Close). That branch is the
  abandon path: bytes still unread there are discarded explicitly, never
  silently truncated. Pre-fix code had the identical residual exposure
  (its post-kill `p.wg.Wait()` could block the same way) — no regression.
- **`Wait()` semantics.** Now documented honestly: it returns only after
  stderr drains to EOF, which can outlast the child's own exit when a
  spawned writer still holds the pipe; callers needing bounded shutdown use
  `Stop` (as all product code already does). No deadlock path was introduced:
  `Wait` blocks at most until writers close stderr or exit.

Diff summary (`internal/acp/spawn.go`): `Stop`'s goroutine becomes
`p.wg.Wait(); done <- p.cmd.Wait()`; the redundant `p.wg.Wait()` in the done
branch is removed (receiving implies the drain completed); the ctx-timeout
branch gains `p.stderr.Close()` between the group kill and `<-done`; `Wait`
becomes `p.wg.Wait(); return p.cmd.Wait()` with the doc above. ~20 changed
lines, one file, no API change.

### U1 regression tests (`internal/acp/spawn_test.go`)

Deterministic-by-construction, not load-dependent: a re-exec'd child
(`TestGatedStderrChild`) writes 8192 stderr bytes, waits for a "go" line on
stdin, writes 8192 more, announces READY on stdout, exits 7. A
`gatingObserver` parks the copier inside its first `Write`. READY proves all
16384 bytes are already in the pipes, so the only way bytes can go missing is
the reap closing the pipe under the parked copier — the bug being tested.

- `TestSpawnStopDrainsStderrBeforeReaping` — Stop ordering.
- `TestSpawnWaitDrainsStderrBeforeReaping` — Wait ordering.
- `TestSpawnStopTimeoutIsBoundedAfterKill` — guards the restructured
  shutdown (kill → interrupt → drain → reap → ctx error) against hangs;
  child (`TestParkingStderrChild`) sleeps holding stderr open.
- Existing `TestSpawnObservesAllStderrAndRealExit` unchanged (it remains the
  hosted observable).

The 500 ms sleeps in the ordering tests exist only to let a **pre-fix**
`Stop`/`Wait` reap the exited child and close the pipe under the parked
copier before the release; the fixed ordering is correct independent of where
the release lands (verified by reasoning and by 5/5 stable green runs).

**Bidirectional proof (local, go1.27.1 darwin/arm64, macOS 27.0):**

- Fixed code: `go test ./internal/acp/ -count=1` → ok; `-race` → ok;
  5 consecutive runs → 5× ok (2.4–2.6 s each).
- Pre-fix code (spawn.go stashed, tests kept):
  `TestSpawnStopDrainsStderrBeforeReaping` FAIL `observed bytes: 8192`;
  `TestSpawnWaitDrainsStderrBeforeReaping` FAIL `observed bytes: 8192` —
  the truncated first half, same mechanism as the hosted `observed bytes: 0`.
  Both restored green on `git stash pop`.

## U2 — retained-envelope diagnostics at the known failure sites (test-only)

Root cause of the trajectory/runner family remains NOT established; nothing
was fixed, skipped, retried or weakened. Per the repair plan's discriminating
datum, the retained `step-%03d.json` envelopes (with the criterion
`ExitCode` and `Diagnostics`) are now printed when the known assertions
fail, so the next hosted ubuntu failure classifies itself
(signal-death vs start-failure vs other) instead of yielding only
"Steps:1/4, FailureStage:execution".

New package-local test helpers (one per package, no product code):
`internal/{trajectory,runner,app}/verification_diagnostics_test.go` —
`dumpRetainedVerificationSteps(t, journal)` prints `receipt.json`,
`step-*.json`, `process-*.json` from the journal via `t.Logf`; trajectory
also gains `dumpAllRetainedVerifications(t, root)` sweeping
`.parley-runtime/trajectory-verifications/*`.

Privacy: the printed bytes are exactly the product's own retained records.
`Diagnostics` was secret-scrubbed (credential-shaped tokens redacted,
`execute.go` secretPatterns) and bounded (last 100 lines / 4 KiB,
`evidenceMaxLines/evidenceMaxBytes`) before it was ever persisted; `Command`
is the scrubbed representation, never the raw input. Nothing outside the
retained envelopes is logged; output happens only on failure paths.

Wired sites (5 assertions across 4 files — the 4 dispatched sites plus one
sibling in the same family/file):

1. `internal/trajectory/parent_recovery_test.go` — both failure paths inside
   `parentRecoveryFixture` (the `:124` site: helper exit and preview error)
   dump all journals under the fixture root.
2. `internal/trajectory/verification_parent_recovery_test.go:~142` —
   `recovered helper execution did not complete` in
   `recoveredParentLineageFixture` (the `:142` site) dumps the exact journal
   computed from the cycle binding.
3. `internal/trajectory/verification_parent_recovery_test.go:~510` — the
   same assertion in `refusedParentRecoveryFixture` (sibling of the hosted
   rerun's `TestRecoveredParentObservationRequiresExactLineage` one-off at
   `:509`; same file, same family, zero extra surface) dumps `path`.
4. `internal/runner/verification_refused_recovery_test.go:~229` — recovered
   verification did not complete → dumps `journal`.
5. `internal/app/trajectory_parent_recovery_test.go:~150` — the launch-stage
   site fails inside the shared `reconciledHelper` helper, so diagnostics
   ride a `t.Cleanup` that fires only when the test has failed, dumping the
   fixture's journal.

Self-check (throwaway test, run then deleted — not committed): populated a
journal with `step-001.json`/`receipt.json`/`launch.json`; helper printed the
step and receipt envelopes, excluded `launch.json`, and degraded gracefully
on a missing journal. On green runs the diagnostics are inert (verified: all
site tests pass locally with no diagnostic output).

## Source-change inventory (complete)

| file | kind | change |
|---|---|---|
| `internal/acp/spawn.go` | product | drain-before-reap in `Stop`+`Wait`, bounded-abandon close, docs |
| `internal/acp/spawn_test.go` | test | 3 regression tests + 2 child helpers + gating observer |
| `internal/trajectory/parent_recovery_test.go` | test | 2 dump calls (fixture failure paths) |
| `internal/trajectory/verification_parent_recovery_test.go` | test | 2 dump calls (:142 site + :509 sibling) |
| `internal/trajectory/verification_diagnostics_test.go` | test | new diagnostic helpers |
| `internal/runner/verification_refused_recovery_test.go` | test | 1 dump call (:229 site) |
| `internal/runner/verification_diagnostics_test.go` | test | new diagnostic helper |
| `internal/app/trajectory_parent_recovery_test.go` | test | failed-test cleanup dump (:153 site) |
| `internal/app/verification_diagnostics_test.go` | test | new diagnostic helper |
| this file | record | evidence |

No other file changed. `organizer-usage.md` (pre-existing local modification,
organizer-owned) was left untouched and is NOT part of this commit.

## Exact checks run (macOS arm64, go1.27.1, base 4e0c069)

- `go build ./...` → ok.
- `go vet ./internal/acp/ ./internal/trajectory/ ./internal/runner/ ./internal/app/` → clean.
- `gofmt -l` on all touched files → clean (pre-existing drift in untouched
  `internal/app/facilitator_test.go`, `protocol_test.go` left alone).
- `go test ./internal/acp/ -count=1` → ok (2.4–2.6 s, includes 3 new tests);
  `-race` → ok; 5× repeated → 5× ok.
- Pre-fix regression proof: spawn.go stashed → both drain tests FAIL
  (`observed bytes: 8192`); restored → ok.
- `go test ./internal/trajectory/ -run 'TestParentRecoveryPublicationInterruptionAndExactReplay|TestRecoveredParentPublicationRefusesNonRefusalOriginal|TestRecoveredParentObservationRequiresExactLineage' -count=1` → all PASS (68.9 s).
- `go test ./internal/runner/ -run 'TestVerifierBudgetRefusalRecoversExplicitly' -count=1` → PASS (5.5 s).
- `go test ./internal/app/ -run 'TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts' -count=1` → PASS (7.1 s).
- No full-suite rerun (per dispatch: not necessary; last full-suite evidence
  exists in release-ci-fix-zcode-1.md at this source baseline).

## Hosted next step (organizer, after independent review of this commit)

Push this commit through the normal channel (no workflow change exists or is
needed). The existing `Tests` workflow on push gives the hosted ubuntu leg;
the targeted local equivalent of what needs to classify is:

    go test ./internal/acp/... ./internal/runner/... ./internal/trajectory/... ./internal/app/... -count=1 -timeout 45m

Watch: `gh run watch <run-id> --interval 60` (or per-job logs from
`actions/jobs/<id>/logs`). Classification criteria for U2: if any
trajectory/runner failure recurs, the dumped `step-*.json` envelope's
`execution.record.command.exit_code` / `diagnostics` and `process-*.json`
identity discriminate signal-death (exit ≥128 / "signal: killed") from
start-failure (run error, no process identity) from ctx/overflow — that
evidence then grounds whatever fix follows; nothing may be labeled flaky or
silenced in the meantime. U1 is proven locally in both directions; the
hosted ubuntu leg should show `internal/acp ok` with the deterministic race
gone. This commit contains product source → it publishes only as a new
v1.49.1 per protocol; v1.49.0 stays untouched.
