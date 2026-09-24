---
idea: meta-protocol-change-lean-organizer
author: kimi-1
role: independent reviewer (Linux release-channel repair, non-implementer)
artifact: release-linux-review
date: 2026-09-24
verdict: PASS
reviewed_commit: a2db799
base: 4e0c069
---

# Linux release-channel repair review — kimi-1

## Verdict

**PASS.** Commit `a2db799` is exactly the dispatched narrow repair: a correct U1
stderr-drain fix with behavior-proving regression tests, and diagnostic-only U2
wiring at the known failure sites. Scope, provenance and hygiene are clean.
**Ready for the organizer to push the untagged CI diagnostic run.** Pushing CI
is not publishing v1.49.1; the U2 root cause remains OPEN until hosted
evidence classifies it. This review endorses obtaining that evidence, not any
claim that the channel is green.

Independence: own detached tree at `a2db799`
(`worktrees/linux-review-kimi-1`, left in place as evidence). I read the
implementer report (`release-linux-repair-zcode-1.md`), the repair plan and
`release-ci-revalidation-claude-1.md`, then re-derived everything from source
and my own runs. I changed nothing: no fixes, edits, commits, pushes, installs
or CI reruns; all throwaway tests were deleted and `git status` verified empty
after each experiment. Local evidence: go1.27.1 darwin/arm64 — same non-Linux
limitation as the implementer's local runs; the hosted ubuntu leg remains the
classifier.

## U1 — `internal/acp/spawn.go`: fix is correct (verified)

Root cause confirmed in base source (`git show 4e0c069:internal/acp/spawn.go`):
pre-fix `Stop` ran `go func() { done <- p.cmd.Wait() }()` and `Wait` ran
`cmd.Wait()` first, both reaping before the `p.wg` copier drained; `os/exec`
closes the `StderrPipe` read end on `Wait` return, truncating mid-copy. This
matches claude-1 Finding E and the hosted record (`observed bytes: 0` of
16384, exit 7 correct, run 35987916696).

The reorder is sound by construction: the `Stop` goroutine now sends `done`
only after `p.wg.Wait()`, so removing the done-branch `p.wg.Wait()` loses no
guarantee; `Wait` drains to EOF before reaping with honest docs
(spawn.go:158-161). `acp.Spawn` has exactly one product call site
(`internal/runner/acp.go:81`) with `Stop` under 2s ctx at :96 and :257 and
`Kill` at :234; `Process.Wait()` has zero product callers (grep: only
`spawn_test.go:136`). Risk surface is `Stop` only, as the implementer claimed.

Lifecycle matrix — my own throwaway counterexamples (deleted after the run,
Unix re-exec children, results reproduced under `-race`):

| case | mechanism | result |
|---|---|---|
| grandchild in-group holding stderr past child exit | ctx 1s → group SIGKILL reaches writer → EOF → reap | bounded, 1.003s, ctx err |
| grandchild escaped via setsid holding stderr | ctx → group kill misses writer → new `p.stderr.Close()` interrupts copier Read → bounded abandon | bounded, 2.00s, ctx err, child exit 0 preserved |
| ordinary EOF, exit 0, 100-byte tail | drain completes at child exit, no added latency | `Stop` nil in 0.01s, tail complete, exit 0 |
| concurrent `Stop`+`Stop`+`Kill` (unsupported) | double `cmd.Wait` errors ("wait: no child processes"), no deadlock | both returned immediately (run w/o -race; identical pre/post-fix exposure) |
| failed observer (Write error) | `io.Copy` aborts — stderr loss is observer-caused, **pre-existing, unchanged** | `Stop` reaps, returns real exit 7, bounded |
| `Wait` with escaped writer closing 600ms after child exit | drain outlives child, bounded by writer lifetime | returned at 605ms, all 150 bytes observed, exit 0 |

Residual exposures are identical-or-better than pre-fix: a blocking observer
hangs `Stop` either way (contract at spawn.go:23-24 requires non-blocking; the
product observer is a file writer); a child escaping its own process group
could block `<-done` in both versions. Pre-fix was strictly worse — its
post-kill `p.wg.Wait()` could hang on any escaped writer with no interrupt;
the fix adds exactly the `stderr.Close()` that bounds that path.

## U1 — tests prove behavior and do not hang

The gated construction is deterministic for the code under test: READY on
stdout proves all 16384 bytes are in the pipes with the copier parked, so only
a mid-drain pipe close can lose them; the fixed ordering is correct regardless
of when the release lands (the 500 ms sleeps matter only for the pre-fix
demonstration, which is not CI-gating). The bounded-shutdown test carries an
explicit 10 s guard; the ordering tests cannot hang on fixed code (child exit
is independent of the release the test always closes).

Bidirectional proof independently reproduced in my tree:

- Fixed: `go test ./internal/acp/ -count=1` ok (2.5s), `-race` ok (3.6s),
  5/5 repeats ok.
- Base `spawn.go` (4e0c069) + new tests: **both drain tests FAIL,
  `observed bytes: 8192`** — the truncated second half, same mechanism as the
  hosted `observed bytes: 0` (the copier was parked at a different point).
  `TestSpawnObservesAllStderrAndRealExit` still passed locally — consistent
  with the load-dependent hosted signature. Restored; tree verified clean.

## U2 — diagnostics: diagnostic-only, correct sites, bounded, no secrets

- **No pass criteria altered.** All five sites add `t.Logf`-based dumps ahead
  of byte-identical failure assertions (the `if` conditions are untouched diff
  context). The app site rides `t.Cleanup` gated on `t.Failed()` because the
  failure surfaces inside the shared `t.Helper()`'d `reconciledHelper`
  (`trajectory_reconcile_test.go:21`) — minimal and correct. Mechanical scan
  of the full diff: no `t.Skip`, `testing.Short`, retry, filter or
  `continue-on-error`.
- **Correct locations.** Base-file lines match the hosted log exactly:
  `verification_parent_recovery_test.go:142` and `:509`,
  `verification_refused_recovery_test.go:229`, app `:153`
  (`reconciledHelper` call). Hosted `parent_recovery_test.go:124` is the
  subtest's fixture *call* line (t.Helper reporting); the dump was wired at
  the real failure branch inside `parentRecoveryFixture` plus its sibling
  preview-error branch. The `:509` rerun one-off's sibling (`:510`) is covered
  — same file, same family, zero extra surface.
- **Bounded.** Envelope contents are bounded structurally before persistence:
  `Command` ≤ 4096 and `Diagnostics` ≤ 8192 are enforced at
  `internal/trajectory/trajectory.go:156`; Diagnostics is the last-100-lines /
  4 KiB tail (`internal/evidence/execute.go:38-39,631-636`). Fixtures create
  exactly one journal; dumps fire only on failure; green runs verified inert
  (runner site PASS 6.6s, app site PASS 8.8s, trajectory trio ok 88.6s — zero
  `diagnostics:` lines).
- **No secrets.** `Diagnostics` is secret-scrubbed pre-persistence
  (`execute.go:605-636`), `Command` is `ScrubAndTruncate` ("raw text is never
  persisted", `execute.go:195`), raw output exists only as a hash
  (`execute.go:68-70`), `process-*.json` carries pid identity facts. The dump
  prints only `receipt.json`/`step-*`/`process-*`.
- **The discriminating datum will be present.** `verifyCaptured` retains a
  step envelope for every ordinal including incomplete executions
  (`internal/trajectory/captured.go:405-411`), so a recurring hosted failure
  dumps the criterion `exit_code` + `diagnostics`; `process-*.json`
  presence/absence then discriminates signal-death vs start-failure exactly as
  the plan's classification table requires.
- **Helper proof (my throwaway test, deleted):** fabricated journal → printed
  `step-001/002`, `process-001`, `receipt`; correctly excluded `launch.json`,
  `request.json`, `recovery.json`, `random.txt` and a directory named
  `step-003.json`; missing journal and empty sweep degrade to one log line
  each.

## Scope and provenance: clean

`git diff 4e0c069..a2db799` = exactly 9 source files (1 product + 8 test) +
the implementer's record file — matching the claimed inventory. `VERSION`,
`CHANGELOG.md`, `internal/app/version.go`, `go.mod`, `go.sum`, `.github/`
untouched (no version bump, no workflow change). Closed artifacts (FINAL.md,
IMPLEMENTATION.md, consensus.md, review/, round-01/, round-02/) untouched.
Product delta vs `v1.49.0` excluding `parley-deck`/`.github` = the same 9
files. `v1.49.0` = `06e563e` locally **and at origin** (`git ls-remote`).
`organizer-usage.md` local modification correctly excluded. U3 not implemented
(as dispatched); no Windows file touched; no Windows check skipped — Windows
scope remains owner-blocked. `go build ./...` ok; `go vet` on the 4 touched
packages clean; `gofmt` clean on all 9 touched files (pre-existing drift in
untouched `facilitator_test.go`/`protocol_test.go` confirmed, correctly left
alone). No full-suite rerun by me — per dispatch, focused tests only; the full
matrix is the hosted run's job.

## Findings

No blocking findings. Observations:

- **K1 (LOW, test-code, Windows-only, unverified):** the gated children write
  8192-byte chunks; if Windows anonymous-pipe buffering is smaller than the
  in-flight total while the copier is parked, the child could block before
  READY and the new ordering tests would time out loudly (never silently
  mis-pass). Linux/macOS pipe capacity (64 KiB) makes it a non-issue there.
  No Windows host available to confirm; hand to the owner-scoped Windows track
  (chunk ≤ 4096 if it ever bites).
- **K2 (INFO):** a blocking `StderrObserver` hangs `Stop` identically pre- and
  post-fix; the documented contract (spawn.go:23-24) already forbids it and
  the product observer is a file writer. No action.
- **K3 (INFO):** concurrent `Stop`/double-`Wait` is unsupported and errors
  noisily without deadlock, identically pre- and post-fix; product code never
  does this. No action.

## Readiness

**Ready for the organizer to push `a2db799` through the normal (unchanged)
Tests workflow as the untagged CI diagnostic run.** Expected outcomes:
ubuntu `internal/acp` ok with the deterministic race gone (U1 is proven
locally in both directions); any recurring trajectory/runner failure now
self-classifies via the dumped envelopes (`exit_code ≥ 128` / "signal: killed"
+ process identity = signal death; run error + no process identity = start
failure). If U2 recurs, its fix is follow-up work grounded in that evidence —
nothing may be labeled flaky, silenced, skipped or retried, and the channel is
not green until U2 has a root cause. v1.49.0 stays at `06e563e`; any eventual
publish of this product-source change is a new immutable v1.49.1, an
organizer/owner decision, not implied by this PASS.
