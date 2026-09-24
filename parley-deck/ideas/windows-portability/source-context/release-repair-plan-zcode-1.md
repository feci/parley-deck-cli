---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: participant (release repair planning)
artifact: release-repair-plan
date: 2026-09-24
status: BLOCKED-ON-OWNER-SCOPE-DECISION
---

# Release repair plan — zcode-1

Plan only. No source, workflow, tag, release, notes or branch changed; no run
triggered; no consensus/implementation launched. Evidence: full read of
`release-ci-revalidation-claude-1.md`, `release-ci-claude-1.md`,
`release-ci-fix-zcode-1.md`; per-job logs of run `35987916696` pulled myself
(ubuntu job `107594830258`, windows job `107594830296`); source reading; one
local cross-compile experiment (`GOOS=windows go vet ./internal/evidence/` →
`undefined: syscall.Mkfifo` at `tree_report_test.go:97`, byte-matching hosted).
claude-1's findings were re-verified against logs and source; every family below
was re-derived from primary log lines, not inherited.

## What run 35987916696 settled

The workflow fix did its two jobs: Finding A resolved hosted (identity step ran;
`internal/app ok 136.801s` ubuntu; zero `TestEvidenceVerifierProductionClosure`
failures), Finding C resolved hosted (first-ever Windows checkout/setup/build).
macOS green (second time). Remaining: ubuntu fails 3 packages / 4 tests; Windows
fails 13 packages (12 FAIL + 1 test-build failure), ~157 tests (claude-1 count;
my grep: 156 `--- FAIL` lines, 98 unique top-level names). All product source
remains byte-identical to `v1.49.0`.

## Ubuntu families (the channel blocker)

**U1 — acp stderr lost (claude-1 Finding E). Root cause verified in source.**
`internal/acp/spawn.go`: copier goroutine (`:82`) drained via `p.wg`, but both
`Wait` (`:144-146`: `cmd.Wait()` then `p.wg.Wait()`) and `Stop` (`:133` goroutine
`done <- p.cmd.Wait()`, then `p.wg.Wait()`) wait the process before draining.
`os/exec` closes `StderrPipe` on `Wait`; the copier can read 0 bytes — hosted
ubuntu recorded exactly `spawn_test.go:38: observed bytes: 0` while exit code 7
was correct. Load-dependent: same package `ok 0.038s` on windows, green on macOS.
Fix: drain `p.wg` before `cmd.Wait()` (restructure the Stop goroutine to
`p.wg.Wait(); done <- p.cmd.Wait()`; copiers EOF at child exit, so no added
latency). Product source, ~6 lines, one file. Provenance `a637628` (v1.5.0,
2026-05-24) — preexisting.

**U2 — trajectory/runner captured-execution family. Root cause NOT established;
no fix proposed.** Failing tests this run: `TestVerifierBudgetRefusalRecoversExplicitly`
(runner, `verification_refused_recovery_test.go:229`, Steps:1/execution),
`TestParentRecoveryPublicationInterruptionAndExactReplay/after-publication`
(`parent_recovery_test.go:124`), `TestRecoveredParentPublicationRefusesNonRefusalOriginal`
(`verification_parent_recovery_test.go:142`, Steps:1/execution). Third distinct
set in three hosted ubuntu executions; never reproduced locally; identity
present this run, so identity-independent. **exit status 19 traced:** the helper
(`parent_recovery_test.go:55-70`) re-execs the test binary, runs
`ExecuteCapturedVerification`, prints its error and `os.Exit(19)` — the captured
stderr was `"captured criterion execution is incomplete"`, i.e. exit 19 is the
same shared condition surfacing through a helper, **not a distinct cause**.
Mechanism constrained by code (`internal/evidence/execute.go:271-277`): Complete
is false only on output overflow (>4 MiB — implausible for these fixtures),
invalid envelope, invalid go-test stream, `ctx.Err()` (tests failed in 2-3 s
against 60-90 s timeouts — implausible), controlled exit ≥128 (signalled child),
or a non-`ExitError` run error (e.g. `Start` failing under process/memory
pressure). The discriminating datum — the criterion envelope's ExitCode and
Diagnostics — **is recorded in retained `step-*.json` but never printed by the
failing assertions**. Smallest next step is test-only observability: at the four
failure sites (`parent_recovery_test.go:124`,
`verification_parent_recovery_test.go:142`,
`verification_refused_recovery_test.go:229`,
`trajectory_parent_recovery_test.go:153`) print the retained step envelope,
then one hosted ubuntu rerun of the three packages. Classification (signal
death vs start failure vs other) precedes any fix; nothing gets labeled flaky
and nothing gets skipped, retried or weakened.

**U3 — refusalGit identity (lead, confirmed in source).** Workflow identity
greens CI but masks the product dependence (`evidence_refusals.go` refusalGit
commits with no identity). Fix: per-command `-c user.name=parley -c
user.email=parley@localhost`, mirroring the existing precedent at
`internal/runner/reviewsnapshot.go:139-140`. ~2 lines, one file, product source.
Not required for a green channel; required for product honesty on
identity-less hosts.

## Windows families (first-ever execution; preexisting surface, none from A-D)

- **W1 snapshot privacy guard (PRODUCT).** `snapshot.go:162` (`Perm()&0077`)
  and `:544` (same expression for files). Go on Windows synthesizes 0777 (dirs)
  / 0666 (files), so the guard can never pass — 71 log messages, ~69 distinct
  tests across trajectory/budget/driver/runner. Snapshot capture in shipped
  Windows binaries refuses by construction. Provenance `4e39609` (2026-09-11).
  Smallest fix (Unix-only perm check) **weakens verification** on Windows; a
  real fix verifies ACLs — design work, not a patch.
- **W2 pipeline gate filename (PRODUCT).** `gate.go:40` `EdgeID = from+"->"+to`
  used in the filename at `:105`; `>` is invalid in Windows filenames — 9
  failures, every pipeline gate write. Provenance `d387de6` (v1.6.0, 2026-06-02).
  Smallest fix sanitizes only the path (keeps EdgeID identity) but creates
  per-OS on-disk names — a compatibility decision.
- **W3 cross-process lock/open/rename (PRODUCT).** ~14 failures
  (`ledger_test.go:57`, `review_test.go:438`, `verification_test.go:368`):
  "file is being used by another process". The budget locking model
  (`budget/lock.go` et al.) assumes POSIX advisory semantics. Needs a
  LockFileEx-style design — architecture-scale.
- **W4 evidence test build failure (TEST).** `tree_report_test.go:97`
  `syscall.Mkfifo`, no build tag; locally reproduced via `GOOS=windows go vet`.
  Fix: build tag + Windows fallback fixture.
- **W5 `/bin/sh` fixtures (TEST).** ~13 failures: `launch_test.go`,
  `protocol_context_test.go`, `telemetry_test.go` (internal/runner).
- **W6 app fixture + panic (TEST).** `writeFakeParleyDeckSkill`
  (`app_test.go:1742`) writes an extension-less `#!/bin/sh` script Windows
  cannot execute; the unchecked assertion at `app_test.go:155` then panics,
  aborting the package at 3.499s — **`wait`/`usage` have still never executed
  on Windows**. Two defects: portable fixture + checked assertion.
- **W7 chmod-unreadable fixtures (TEST).** ~6 (consensus `impl_test.go:807`,
  `phase_event_test.go:156/170/186`, `strict_gate_test.go:179`): chmod does not
  block owner reads on Windows.
- **W8 HOME-only fixtures (TEST).** 2 (`agents/configmodel_test.go`): needs
  `USERPROFILE` too.
- **W9 POSIX-path worktree fixtures (TEST).** ~3-4 (`worktree_inventory_test.go`):
  feeds `/repo` paths to a `filepath.IsAbs` product check.
- **W10 CRLF (TEST).** 1 (`hardening_test.go:456`).

U1 also exists on Windows in principle (race) but the package passed there —
load-dependent, no claim made.

## Smallest-fix set vs what it cannot do

**v1.49.1 (ubuntu channel repair, source → new version per protocol):**
`internal/acp/spawn.go` (U1, ~6 lines), `internal/app/evidence_refusals.go`
(U3, ~2 lines), test-only diagnostics patch (U2, 4 sites), plus
`VERSION`/`version.go`/`CHANGELOG.md` (correct the Windows sentence here — the
tag's text is immutable). ~5 files + metadata. Validation must be hosted:
ubuntu leg 31/31; U1 also re-observed on macOS; **if U2 recurs with
diagnostics, its fix is follow-up, not silence — the channel is not green while
it fails.** No skip, filter, retry or `continue-on-error` enters the workflow.

**Windows green leg: NOT honestly achievable inside release remediation.**
W1-W3 are product design decisions (ACL verification, lock semantics, filename
compatibility); W4-W10 is a ~12-file test-portability sweep. Full set ≈ 3-4
product files + ~12 test files + design — its own reviewed track with consensus.

**No-version, owner action available now:** edit the `v1.49.0` release NOTES
(mutable) to correct the falsified Windows sentence in place, per revalidation
§remedy 1. Tag stays untouched.

## A-D deployment remediation vs preexisting hardening

Authorized release remediation (this task's lineage): workflow environment
(done, hosted-verified), notes correction, v1.49.1 metadata + the two smallest
product fixes the channel itself exposed (U1 blocks a green ubuntu leg; U3 was
already flagged at dispatch as the v1.49.1 candidate), U2 diagnostics,
re-validation. **Everything Windows (W1-W10) is preexisting product hardening**
— provenance 2026-05-24/06-02/09-11, none introduced by lean-organizer — newly
exposed by this idea's first-ever CI. It is real debt, and it is not this
release's regression.

## Native Windows contract — honest position

Windows binaries are published, but snapshot verification, pipeline gates and
cross-process budget locking are broken-or-unvalidated by construction, and
`wait`/`usage` have never executed there (W6 abort). A native contract needs
architecture work: ACL-based privacy verification, a Windows lock strategy,
per-OS gate filenames, extension-aware exec, plus the fixture sweep; hosted
coverage is x64-only (`windows-arm64.exe` has zero executed evidence) and runs
with `core.longpaths`, not the OS default. That is beyond this task and beyond
a patch release.

## Residual unknowns

U2 root cause (discriminator identified, not yet run); whether U2 is one
mechanism or several; Windows ACL/lock/filename designs; whether W1's minimal
form is acceptable (it weakens verification); arm64-Windows behavior entirely.

## Owner decision required (scope exceeds authorized remediation)

Authorize **one** of: **(a)** a Windows-portability track (own idea + consensus,
v1.50-class, weeks) targeting a green Windows leg; or **(b)** descope — declare
Windows experimental/unvalidated in v1.49.1 changelog + corrected notes, fix
W1-W10 later. Either way v1.49.1 (U1+U3+diagnostics) and the notes correction
proceed independently of this choice.

## Blocking status and exact recommended action

**BLOCKED.** Channel is not green; Ubuntu fails U1+U2; Windows is a preexisting
portability surface beyond this dispatch. Exact recommended action, in order:
(1) owner corrects the `v1.49.0` release notes in place now; (2) owner/organizer
chooses (a) or (b) above — required before any Windows work, not before v1.49.1;
(3) implement v1.49.1 as the 5-file set above with the U2 observability patch,
validating on hosted ubuntu (and macOS), U2 classified from diagnostics before
any claim of green; (4) do not move `v1.49.0`, do not relabel U2 flaky, do not
skip anything to force green.
