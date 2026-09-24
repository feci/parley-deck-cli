---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: narrow Linux repair implementer (post-review follow-up)
artifact: release-linux-repair-followup
date: 2026-09-24
base: a2db799
inputs: release-linux-review-claude-1.md (F-1/F-2/F-9), release-linux-review-kimi-1.md (K1/K2)
---

# Linux repair review follow-up — zcode-1

Dispatched scope: apply the minimal F-2 correction (wrong diagnostic
directory at the exit-19 site), add discriminating F-1 coverage for the
ctx-branch `p.stderr.Close()` line, bound the F-9 test waits, evaluate K1
honestly and record its disposition. No new product behavior beyond the
already reviewed U1; no Windows work, no U3, no skips, no scope expansion.
Both reviews' verdicts (PASS, safe to push untagged) stand; this follow-up
closes their actionable findings before that push.

## F-2 (MEDIUM-HIGH) — applied: the exit-19 site now dumps real envelopes

Claude established that `parentRecoveryFixture`
(`internal/trajectory/parent_recovery_test.go`) discarded the journal path
`journalFixture` returns and instead swept
`<root>/.parley-runtime/trajectory-verifications/*` — a directory that
does not hold the journals (they live at
`filepath.Dir(b.Store.Dir)/trajectory-verifications/<EntryKey>`, the
expression `journalFixture` computes and the working `:142`/`:509` sites
use). Correction, exactly the review's remedy:

- `parent_recovery_test.go:75` — `ticket, journal := journalFixture(...)`
  binds the fixture-retained journal path (previously `_`).
- Both failure branches in `parentRecoveryFixture` — the hosted
  `helper: exit status 19` site (`cmd.Wait` error) and its sibling
  `PreviewParentRecovery` error branch — now call
  `dumpRetainedVerificationSteps(t, journal)`, the same helper the four
  working sites use.
- `dumpAllRetainedVerifications` removed from
  `internal/trajectory/verification_diagnostics_test.go`: its only two
  callers were the two corrected sites, and its wrong-directory glob was
  the defect. No other references remain (grep).

Forced-failure proof (throwaway mutations, reverted and verified): forcing
each branch to fire made all four
`TestParentRecoveryPublicationInterruptionAndExactReplay` subtests print
the full 9-envelope set — `receipt.json`, `step-001..004.json`,
`process-001..004.json` — per subtest, at both branches. Green runs print
zero `diagnostics:` lines (verified on the exact test, count=0). No pass
criterion touched: the two dump calls remain statements ahead of
byte-identical `t.Fatal` assertions.

## F-1 (MEDIUM) — applied: discriminating coverage for `p.stderr.Close()`

Claude's mutation experiment showed `TestSpawnStopTimeoutIsBoundedAfterKill`
passes identically with the `_ = p.stderr.Close()` line deleted: its parking
child dies with the group kill, so EOF arrives regardless and nothing in
the commit guarded the riskiest new line. New file
`internal/acp/spawn_unix_test.go` (`//go:build !windows`),
`TestSpawnStopIsBoundedWhenStderrWriterEscapesGroup`:

- The re-exec'd child (`TestEscapedStderrChild`) writes a 64-byte stderr
  marker, spawns a grandchild (`TestEscapedStderrGrandchild`) with its own
  `syscall.SysProcAttr{Setsid: true}` and `Stderr: os.Stderr` — so the
  grandchild detaches into its own session and parks holding the pipe's
  only surviving write end — reports `GRAND <pid>` then `READY`, and parks.
- The test kills cleanup risk by `t.Cleanup` SIGKILL on the reported
  grandchild pid (it escapes the group kill by construction, so the test
  must reap it itself).
- `Stop` under a 300 ms ctx: the deadline branch kills the child's group
  (grandchild survives), and only `p.stderr.Close()` can interrupt the
  copier parked in `Read` — the exact condition the line exists for.

Bidirectional proof, following the review's method (stash → delete only
the `_ = p.stderr.Close()` line from `spawn.go` → run → restore → rerun;
`git diff` empty after restore):

- Line deleted: `--- FAIL: TestSpawnStopIsBoundedWhenStderrWriterEscapesGroup
  (10.01s) — Stop did not return after interrupting the group-escaping
  stderr writer`, while `TestSpawnStopTimeoutIsBoundedAfterKill` still
  PASSED (0.30s) — reproducing Claude's finding that the existing test does
  not discriminate the line.
- Line present (committed code): PASS in 0.31s; 5/5 repeats stable.

Platform boundary (inherent, explained in the file's doc comment): the test
needs POSIX process groups for the escape and a pollable pipe fd for the
interrupt (Claude F-8: Windows anonymous pipes are not pollable, so
`Close` does not unblock a pending Read there). It is build-tagged out on
Windows for that structural reason — not skipped, never executed there,
and it implies nothing about Windows, which remains owner-blocked and
unverified. `internal/acp/spawn.go` itself is untouched: the product
delta vs `v1.49.0` remains exactly the one reviewed U1 file, and F-8's
"do not generalise the boundedness comment beyond POSIX" is answered here
and in the test comment rather than by editing reviewed product source.

## F-9 (LOW) — applied: the two drain tests fail fast instead of hanging

`TestSpawnStopDrainsStderrBeforeReaping` and
`TestSpawnWaitDrainsStderrBeforeReaping` awaited `p.Stop(context.Background())`
/ `p.Wait()` on bare channel receives. Both now use the sibling
`TestSpawnStopTimeoutIsBoundedAfterKill` pattern — a 10 s `time.After`
guard with a named failure — so a future drain hang fails the test loudly
instead of consuming the `-timeout 45m` budget and panicking the binary.
No assertion, ctx value or sleep changed. (The F-1 mutation run doubles as
evidence the guard pattern fires: the new test's guard tripped at 10.01 s.)

## K1 (LOW) — evaluated, disposition: record + boundary note, no code change

Kimi's portability observation is correct in its own terms and its
failure-mode analysis holds: while the gating observer parks the copier,
16 KiB is in flight; Linux/macOS pipe capacity (64 KiB) makes the READY
handshake safe there, and any platform with a smaller buffer would block
the child before READY — a loud failure (handshake stall surfacing at
package timeout, or the 16384-byte assertion failing), never a silent
mis-pass. Disposition, recorded as a comment on `spawnGatedChild`:

- **No speculative chunk-shrink.** It cannot be made provably safe: a
  minimal anonymous-pipe buffer could be as small as the Windows default,
  and with partial first reads even 2×4096 in-flight could exceed it —
  there is no principled size without a Windows host to measure on.
- **The 16 KiB volume is load-bearing**: it mirrors the hosted truncation
  signature (0 of 16384 bytes observed, run 35987916696) that the tests
  were built to reproduce deterministically; halving it would weaken that
  mechanism-mirroring for a hypothetical platform.
- **The scenario has never executed**: these tests do not exist in any
  hosted Windows run (Windows is owner-blocked; no Windows claim is made
  or implied). Per Kimi's own framing, `chunk ≤ 4096` is a contingency for
  the owner-scoped Windows track *if it ever bites* — handed there in this
  record, not applied blind.

## K2 / F-7 / F-8 — acknowledged informational, no action

K2 = F-7: a blocking `StderrObserver` parks `Stop` identically pre- and
post-fix; the documented contract (`spawn.go:23-24`) forbids it and the
product observer is a non-blocking file writer. No change. F-8's
POSIX-scoping of the boundedness property is addressed by documenting the
boundary in the new unix-only test (see F-1) rather than touching reviewed
product source. No Windows architecture, U3, W1-W10, workflow, version,
tag, publish, push, merge, install, skip, filter or retry — none performed,
none implied.

## Validation (go1.27.1 darwin/arm64, focused per dispatch)

- `go build ./...` ok; `go vet ./internal/acp/ ./internal/trajectory/` ok;
  `gofmt -l` clean on both packages (pre-existing drift in untouched files
  left alone, as before).
- `go test ./internal/acp/ -count=3` → ok 8.179s; `-race -count=1` → ok
  3.861s; new F-1 test `-count=5` → 5/5 PASS (0.31s each).
- F-1 bidirectional mutation: line deleted → FAIL 10.01s at the guard,
  parking test still PASS; restored byte-identical (`git diff` empty) →
  both PASS.
- F-2 forced failures: exit-19 branch and preview branch each → 4/4
  subtests × 9 envelopes printed; reverted byte-identical; green rerun →
  zero `diagnostics:` lines.
- Trajectory trio (same set as the prior record):
  `TestParentRecoveryPublicationInterruptionAndExactReplay` 26.13s,
  `TestRecoveredParentObservationRequiresExactLineage` 52.78s,
  `TestRecoveredParentPublicationRefusesNonRefusalOriginal` 6.72s — all
  PASS, package ok 85.947s.
- No full-suite rerun: source delta is test-only on top of a2db799 and the
  last full-suite evidence at this baseline is in
  `release-ci-fix-zcode-1.md`; runner/app untouched this round (their dump
  sites are claude-1 F-3-verified and unmodified).
- All throwaway mutations reverted and verified; working tree before
  commit contained only the four owned source files plus this record.

## Owned paths in the follow-up commit

- `internal/acp/spawn_test.go` (F-9 guards, K1 boundary note)
- `internal/acp/spawn_unix_test.go` (new; F-1)
- `internal/trajectory/parent_recovery_test.go` (F-2)
- `internal/trajectory/verification_diagnostics_test.go` (F-2 helper removal)
- this record

Reviews (`release-linux-review-claude-1.md`,
`release-linux-review-kimi-1.md`), closed artifacts, `organizer-usage.md`,
runs/ and inbox/ are untouched and out of the commit. Organizer next:
independent delta review, then push for the hosted ubuntu untagged CI run;
a recurrence at `:124` now self-classifies like its siblings
(`exit_code ≥ 128`/diagnostics = signal death vs run-error/no process
identity = start failure). U2 root cause remains unestablished; the
channel is not green and this work claims nothing about Windows.
