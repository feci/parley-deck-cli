---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent non-implementer reviewer (post-release Linux repair)
artifact: release-linux-review
date: 2026-09-24
reviewed: a2db7998fae61d44fdb9786d3c12b4a5b9cc868f
base: 4e0c069
verdict: PASS (safe to push) — with one MEDIUM-HIGH correction recommended first (F-2)
---

# Linux release-channel repair review — claude-1

## Verdict

**PASS.** Commit `a2db799` is correct for U1, scope-clean, and safe for the organizer to
push as an **untagged CI diagnostic run**. U1's root cause, fix and regression proof all
hold under independent counterexample.

**One finding should be fixed before that push** (one line, test-only): **F-2** — the
`internal/trajectory/parent_recovery_test.go:124` diagnostic site emits **no envelopes at
all**, because it globs a directory the journals are not in. That is one of the four
dispatched U2 sites, and it is the site of the hosted `helper: exit status 19` family. The
other four wired sites work. Pushing without F-2 is still safe and still worth doing — it
would classify the `:142` / `:229` / `:153` / `:509` family — but a recurrence at `:124`
would yield "no retained verification journals" and cost a repeat trip.

Pushing this is **not** publishing v1.49.1. U2's root cause remains **unestablished** and
stays open until hosted evidence exists. Nothing here claims the channel is green.

**Source pin.** Reviewed `a2db799` against base `4e0c069` in my own detached worktree
(`worktrees/review-claude-1-linux`), go1.27.1 darwin/arm64, macOS 27.0. Every mutation I
made was reverted; `git status --porcelain` was empty (tree byte-identical to `a2db799`)
before the final confirmation runs. I changed nothing in the branch under review, made no
commit, push, install, tag, release or CI run.

## Severity index

| ID | Severity | Area | Summary |
|---|---|---|---|
| F-1 | MEDIUM | U1 tests | the newly added `p.stderr.Close()` line has **zero** discriminating coverage |
| F-2 | MEDIUM-HIGH | U2 wiring | `:124` site globs the wrong directory → dumps nothing at the `exit status 19` site |
| F-3 | PASS | U2 wiring | the other 4 sites each dump the full 9-envelope set (forced-failure verified) |
| F-4 | PASS | U2 privacy | output bounded ~12.5–13.3 KB per failing invocation, structurally |
| F-5 | PASS | U2 privacy | zero credential-shaped tokens; unscrubbed field cannot carry the payload |
| F-6 | PASS | U2 safety | pass criteria unaltered; diagnostics inert on green (0 lines, 5/5 tests) |
| F-7 | LOW (info) | U1 residual | a blocked observer still parks Stop — identical pre/post, unreachable in product |
| F-8 | LOW (info) | U1 portability | "Stop stays bounded" is a Unix-pollable-fd property, not a Windows one |
| F-9 | LOW | U1 tests | two new tests use unbounded ctx with no `time.After` guard |
| F-10 | PASS | scope | 9 Go files + 1 record; tag, workflow, version metadata, closed artifacts all untouched |

## U1 — `internal/acp/spawn.go` drain-before-reap

### Root cause and fix: confirmed correct

The reordering matches the `os/exec` contract (`Wait` closes `StderrPipe` on return, so all
reads must complete first). Lifecycle audit re-derived independently:

- **Product surface.** `acp.Process` is used at `internal/runner/acp.go:81` (Spawn), `:96`
  and `:257` — both `Stop(ctx)` with a 2 s timeout. **`Process.Wait()` has zero product
  callers**; I grepped all non-test `.Wait()` calls and the only hits are unrelated types
  (`handle.Wait()` in `internal/app/app.go`, `tracked.Wait()` in `internal/runner/launch.go`).
  The implementer's claim holds.
- **Nil-safety.** `&Process{` appears only at `spawn.go:75`, so the unguarded
  `p.stderr.Close()` cannot nil-deref. (It is asymmetric with the `p.stdin != nil` guard
  three lines above, but harmless.)
- **Ordinary EOF** and **failed observer**: a `MultiWriter` error aborts the copier early,
  after which `cmd.Wait()` proceeds; if the child then blocks on a full pipe, `Stop` is still
  bounded by ctx. Unchanged from pre-fix.
- **Concurrent callers.** Second `Stop` returns `exec: Wait was already called` promptly —
  measured, see F-3 evidence below. `p.wg.Wait()` and double `Close()` are both safe.

### Bidirectional regression proof: independently reproduced

Pre-fix `spawn.go` (from `4e0c069`) + the committed new tests, my tree:

```
--- PASS: TestSpawnObservesAllStderrAndRealExit (0.00s)     <- hosted observable passes LOCALLY pre-fix
    spawn_test.go:127: observed bytes: 8192
--- FAIL: TestSpawnStopDrainsStderrBeforeReaping (1.01s)
    spawn_test.go:142: observed bytes: 8192
--- FAIL: TestSpawnWaitDrainsStderrBeforeReaping (1.01s)
--- PASS: TestSpawnStopTimeoutIsBoundedAfterKill (0.30s)
```

Fixed code: `go test ./internal/acp/ -count=3` → `ok 7.183s`; `-race -count=1` → `ok 3.560s`.

This reproduces the implementer's numbers exactly and adds a point worth recording: the
**pre-existing hosted observable passes locally even pre-fix**, so the two new tests are a
genuine deterministic upgrade over what CI had, not a restatement of it.

### F-1 (MEDIUM) — the new `p.stderr.Close()` line is untested

Mutation experiment: I deleted only the line `_ = p.stderr.Close()` from the ctx branch,
leaving drain-first intact.

```
--- PASS: TestSpawnStopTimeoutIsBoundedAfterKill (0.30s)     <- still green without the line
    REVIEW C2: Stop did NOT return within 3s (group-escaping stderr writer held the copier)
--- FAIL: TestReviewStopBoundedWithGroupEscapingStderrWriter (3.00s)
```

`TestSpawnStopTimeoutIsBoundedAfterKill` passes identically **pre-fix, post-fix, and with the
Close removed**. Its child (`TestParkingStderrChild`) holds the only stderr write end and
dies with the group kill, so EOF arrives regardless — the test never reaches the condition
the line exists for. Of the three new tests, two discriminate; this one is a hang guard only.

My counterexample (reviewer-only, not committed): a re-exec'd child starts a grandchild with
its own `Setsid`, hands it the inherited stderr pipe, then exits — so `killProcessGroup`
cannot reach the writer and the copier stays blocked in `Read`. Result with the line present:

```
REVIEW C2: Stop returned after 302.290ms err=context deadline exceeded | observed 64 of 64 stderr bytes | ring=64
```

So the line does exactly what it claims. It is simply the riskiest new line in the change
and nothing in the commit would catch its removal. Not a correctness defect — a durability
gap, suitable for the same follow-up that carries F-2.

### Measured behavior change in the escaped-writer path (MEDIUM-LOW, record-only)

Same scenario, pre-fix code: `Stop returned after 3.638ms err=<nil> | observed 64 of 64`.

So in this path `Stop` now consumes the **full caller ctx budget** (2 s in product) and
returns `ctx.Err()` where it previously returned `nil` almost immediately. The record's
"pre-fix had the identical residual exposure — no regression" is accurate about the
*hang* exposure (F-7) but understates this: the honest statement is **a bounded latency
cost, up to the caller's own ctx, bought in exchange for not truncating**. That is the
right trade and the product is indifferent to it:

- both call sites discard the `Stop` error (`_ = process.Stop(...)`);
- `ExitCode()` is still populated on that path, which `finishACP` depends on — verified:
  `REVIEW C5: first Stop=context deadline exceeded ExitCode=0`, and
  `REVIEW C5: second Stop returned exec: Wait was already called`.

No change requested; it should just not be described as costless.

### F-7 (LOW, informational) — blocked-observer residual, genuinely unchanged

With an observer that parks forever inside `Write`, the copier is blocked **outside** `Read`,
so closing the read end cannot interrupt it and `Stop` does not return. Measured identically
on fixed and pre-fix code (`Stop did NOT return within 3s` in both arms) — **no regression**,
confirming the implementer's claim. Unreachable in product: the real observer is
`telemetry.streamWriter` (`internal/telemetry/usage.go:48`), which takes a short mutex and
returns `len(data), nil` unconditionally. The `StderrObserver` doc already requires
non-blocking behaviour.

### F-8 (LOW, informational) — the boundedness claim is Unix-scoped

`p.stderr.Close()` interrupts a pending `Read` because Unix `StderrPipe` files are pollable
(`pipe2(O_NONBLOCK)`, netpoll-registered). Windows anonymous pipes from `CreatePipe` are not
pollable, so the same `Close` does not unblock a pending read there. Windows is owner-escalated
and `internal/acp` passed on the Windows leg; I make no Windows claim. Only: do not generalise
the code comment's "Stop stays bounded" beyond POSIX.

### F-9 (LOW) — unbounded contexts in the two new drain tests

`TestSpawnStopDrainsStderrBeforeReaping` calls `p.Stop(context.Background())` and
`TestSpawnWaitDrainsStderrBeforeReaping` calls `p.Wait()` (no ctx at all), each awaited on a
bare channel receive with no `time.After` guard — unlike the sibling
`TestSpawnStopTimeoutIsBoundedAfterKill`, which guards with `10*time.Second`, and unlike the
pre-existing `TestSpawnObservesAllStderrAndRealExit`, which uses a 5 s ctx. Today they are
stable (1.01 s each, 3/3 + `-race`). But if a future change turns the drain into a hang, they
fail by consuming the `-timeout 45m` budget and panicking the whole test binary instead of
failing fast. The two 500 ms sleeps also add ~2 s to the package; acceptable.

## U2 — retained-envelope diagnostics

I verified wiring by **forcing each assertion to fail in my own tree** (e.g. `receipt.Steps
!= 4` → `!= 99`) and reading what actually got printed. Forcing was reverted; no forced
condition exists in the reviewed commit.

### F-2 (MEDIUM-HIGH) — the `:124` site dumps nothing

`dumpAllRetainedVerifications` globs

```
<root>/.parley-runtime/trajectory-verifications/*
```

but a filesystem probe of the live fixture shows the journals are at

```
<root>/.git/parley-launch-budgets/cycles-<hash>/trajectory-verifications/<charge>
```

`<root>/.parley-runtime/trajectory-verification` (singular, no `s`) **does** exist — it is
the *request* directory, a different artifact — which is the likely source of the mistake.
Forced failure output at that site, all four subtests:

```
parent_recovery_test.go:126: diagnostics: no retained verification journal... (glob err: <nil>)
```

Zero envelopes. That is exactly the hosted
`TestParentRecoveryPublicationInterruptionAndExactReplay/after-publication` /
`helper: exit status 19` site the plan named. It degrades gracefully and breaks no test, so
it is not a correctness regression — it is a diagnostic that will be silent precisely when
it is needed.

**Remedy is one line and already precedented inside this same change.**
`parentRecoveryFixture` calls `ticket, _ := journalFixture(...)` and discards the journal
path that `journalFixture` (`verification_test.go:58`) computes with the exact expression the
working `:142` and `:509` sites use. Binding it and calling `dumpRetainedVerificationSteps`
instead — verified in my tree:

```
parent_recovery_test.go:126: diagnostics: retained receipt.json
parent_recovery_test.go:126: diagnostics: retained step-001.json … step-004.json
parent_recovery_test.go:126: diagnostics: retained process-001.json … process-004.json
```

Full envelope set. This also removes the need for the `dumpAllRetainedVerifications` helper
(its only two callers are these two sites). I did **not** apply this to the branch.

### F-3 (PASS) — the other four sites all fire correctly

Each forced to fail; each printed `receipt.json` + `step-001..004.json` +
`process-001..004.json`, **zero** `unavailable` / `unreadable` messages:

| site | file:line reported | envelopes |
|---|---|---|
| trajectory `:142` (`recoveredParentLineageFixture`) | `verification_parent_recovery_test.go:142` | 9 |
| trajectory `:509` sibling (`refusedParentRecoveryFixture`) | `:510` | 9 × 6 subtests |
| runner `:229` | `verification_refused_recovery_test.go:229` | 9 |
| app `:153` (via `t.Cleanup` on `t.Failed()`) | `trajectory_parent_recovery_test.go:158` | 9 |

The `:142` site is the one I most doubted — its path is hand-written rather than
fixture-returned — and it resolves correctly. The app site's cleanup is registered *after*
`trajectoryHelperFixture`, so LIFO ordering runs the dump **before** the fixture's own
TempDir teardown; confirmed empirically.

**The discriminating datum is present.** `step-001.json` carries
`execution.complete`, `execution.record.command.exit_code`, `duration_ms` and `diagnostics`;
`process-00N.json` carries `pid`/`pgid`/`boot_id`/`proc_start`/`started_at`. That is enough
to separate signal death (exit ≥128) from a start failure (no process identity) from
ctx/overflow, which is what the repair plan asked for.

### F-4 (PASS) — output is bounded, structurally

Measured per failing invocation: **13,250 B** (trajectory `:142`), **12,556 B** (runner
`:229`), ~13.6 KB per subtest at `:509`. The bound is not incidental: both unbounded fields
in `CommandEvidence` — `Command` and `Diagnostics` — pass `ScrubAndTruncate`
(`internal/evidence/execute.go:620`: last 100 lines, then last 4096 bytes) **before
persistence**; every other field is a scalar or a hex digest. The dump prints the retained
file verbatim, so it inherits that bound and cannot exceed it.

### F-5 (PASS) — no secrets

Credential-pattern scan (bearer / `sk-` / `gh?_` / `AKIA` / labelled token=value) across
every dump I captured: **zero hits**. Stronger than a scan, the structure holds: the
criterion text appears only in `step-*.json`'s `command`, which is scrubbed. The one field
*not* passed through `ScrubAndTruncate` is `process-*.json` `identity.command` — and it
carries only the fixed captured-criterion wrapper
(`sh -c trap ':' TERM … parley_command=$PARLEY_CAPTURED_COMMAND; unset …`), which reads the
real command from an env var and unsets it. So the unscrubbed field structurally cannot
carry the criterion payload. Runner-local absolute temp paths and a `boot_id` appear; neither
is sensitive on a hosted runner.

### F-6 (PASS) — pass criteria unaltered, diagnostics inert on green

Every insertion is a statement placed *before* an existing `t.Fatal`/`t.Fatalf`; no assertion
condition was touched. The app site adds only a `t.Failed()`-guarded `t.Cleanup` and binds a
previously discarded `journal`. Green-path confirmation on the clean `a2db799` tree:

```
ok  internal/trajectory  78.627s   (3 named tests PASS)
ok  internal/runner       6.389s
ok  internal/app          8.121s
diagnostics: lines on green runs — 0 / 0 / 0
```

## F-10 (PASS) — scope and provenance

- `a2db799` touches **9 Go files + 1 record** (`release-linux-repair-zcode-1.md`). Nothing else.
- **`.github/` untouched** by this commit; `VERSION`, `CHANGELOG.md`,
  `internal/app/version.go` untouched. No skip, filter, retry, `-run`, `-short`,
  `continue-on-error` or suppression is introduced — there is no workflow change at all.
- Closed artifacts (`FINAL.md`, `IMPLEMENTATION.md`, `review/`, `consensus.md`) untouched;
  `organizer-usage.md` correctly left out of the commit.
- `v1.49.0` = `06e563e8b1fe8094132149b83e14da7be9aae51e`, unmoved.
- Product source vs the tag: `git diff --name-only v1.49.0 a2db799 -- . ':!parley-deck'
  ':!.github'` lists **one non-test file — `internal/acp/spawn.go`**. So the commit does
  publish only as a new immutable v1.49.1, exactly as the record states. No version bump was
  made and none is authorised here.
- **U3 correctly not implemented** (`refusalGit` identity, `internal/app/evidence_refusals.go`
  untouched); **W1–W10 correctly untouched**; no Windows fix attempted; no check skipped.

## Readiness for the organizer

**Ready to push an untagged CI diagnostic run.** Recommended order:

1. Apply **F-2** (one line, test-only, in `parent_recovery_test.go`) so the `exit status 19`
   family also yields envelopes. Optionally fold in **F-1** coverage. Both are test-only and
   do not change the product-source footprint or the v1.49.1 question.
2. Push. The existing `Tests` workflow gives the hosted ubuntu leg unchanged; the targeted
   equivalent is
   `go test ./internal/acp/... ./internal/runner/... ./internal/trajectory/... ./internal/app/... -count=1 -timeout 45m`.
3. Expect `internal/acp ok` (U1 proven in both directions locally). For U2, read
   `execution.record.command.exit_code` + `diagnostics` and the `process-*.json` identity to
   classify. **Absence of a recurrence classifies nothing** — it is not evidence the family
   is fixed, because nothing was fixed.

**Not authorised by this review and not performed by me:** any version bump, tag, publish,
merge, notes edit, workflow change, or CI trigger. `v1.49.0` stays immutable. U2's root
cause is still unknown and the channel is **not** green.

## Limitations

- macOS arm64 / go1.27.1 only. I have no Linux or Windows host; F-8 is reasoned from the
  `internal/poll` pollable-fd contract, not executed on Windows.
- F-2's on-disk path was probed in the `before-publication` subtest fixture; the other three
  subtests produced the same "no retained verification journals" message, and both
  `dumpAllRetainedVerifications` call sites share that fixture.
- I did not run the full 31-package suite: no finding required it, and the last full-suite
  evidence at this source baseline is recorded in `release-ci-fix-zcode-1.md`.
- I reviewed the U1/U2 repair only — not the closed A–D implementation, not Windows, not the
  published release notes.
- No credentials, tokens or log secrets appear in this report.
