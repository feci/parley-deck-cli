---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent reviewer of the U2 Linux attribution repair (non-implementer)
artifact: release-linux-u2-review
date: 2026-09-24
subject: ef10cf3 (parent 231d889) — internal/procctl/procctl_linux.go +51/−5, procctl_cmdline_linux_test.go new
inputs: release-linux-hosted-diagnostics-02-claude-1.md (my root cause), release-linux-u2-repair-zcode-1.md (implementer)
environment: own isolated worktree at ef10cf3 (no source edits on disk) + two throwaway golang:1.26 containers (go1.26.8 linux/arm64, 16 CPU) + host go1.27.1 darwin/arm64
verdict: PASS — fix is correct, fail-closed and in-scope; READY for ONE hosted x86_64 ubuntu validation. Three defects found, all in the *reasoning/comments*, none in the gate.
---

# U2 repair review — claude-1

## Verdict

| Question from dispatch | Result | Basis |
|---|---|---|
| Return on **first non-empty**, no value shopping | **PASS** | mutation R reproduced: 2 FAILs |
| **Permanent-empty fails closed** | **PASS** | mutation F reproduced: 2 FAILs |
| **Read error / death** → immediate exit | **PASS (code)** / **C-1 (claim)** | mutation D reproduced: 3 FAILs; but "death ⇒ read error" is false for unreaped processes |
| **PID / start-time / pgid** changes still refuse | **PASS** | `Attributed` byte-identical (git-verified); exact-reason tests green |
| **Cancellation / deadline bounds** | **PASS (bounded)** / **C-1, C-3** | inside every WaitDelay, but by a different mechanism than claimed |
| **Zombie** behaviour | **PASS** | refuses at the same facet, same string, after the bound |
| **Healthy path adds no delay** | **PASS** | one read; measured mean 477 µs on long-lived spawns |
| Scope: Windows/darwin unchanged | **PASS** | `git diff` on all other procctl files is empty |
| Actual weakening of the gate | **NONE FOUND** | see Safety |
| Ready for ONE hosted x86_64 leg | **YES**, with the arm64 caveat (L-1) | see Readiness |

**The delta is right.** `pollCmdline` retries only on a zero-byte read, returns the first
value it sees, exits immediately on a read error, and on exhaustion returns `("", false)`
— bit-for-bit the pre-fix refusal. Every `Attributed` facet is untouched. The three
findings below are all in the *explanation* of the change (two of them written into the
repo), not in its behaviour.

## Independent reproduction

All runs are mine, in my own containers, from a tarball of my own worktree at `ef10cf3`.
Both file hashes matched the implementer's provenance table **before** I ran anything:
`a3f8eb16…040547` and `5b3cca33…2935b26`, and again after every mutation was reverted.

**The race, pre-fix (mutation W = single-read `command()` restored):** reproduced on the
**first attempt** — `capture recorded an empty command for just-started pid 2973` at
`:170`, the *sequential* path, no concurrency involved. Two further `-count=10` runs:
**6/10 and 5/10 iterations failed**. Same order as the implementer's 8-per-`count=10`.
The zombie test also failed as designed (`took 13.417µs` — no poll). U2 is real and now
reproduces reliably on demand.

**Post-fix determinism:** spawn test `-count=30` = **3600 spawns, 30/30 PASS**; a further
`-count=10` (1200 spawns) **10/10 PASS under 8× CPU oversubscription** (loadavg 13.3) —
a load condition the implementer did not test; full package `-count=5` ok; `-count=3
-race` ok 2.676s; `go vet` clean; `gofmt` clean on both touched files (the only drift is
pre-existing, in untouched `procctl_test.go`/`procctl_windows.go`).

**Mutations F / R / D reproduced exactly as reported** (2 / 2 / 3 discriminating FAILs,
with the same test names and messages). The safety pins are genuine, not decorative.

**`TestKillGroupReapsChild` is environmental — confirmed by a different method than the
implementer's.** Rather than re-running it on pre-fix code, I ran the *same post-fix
code* in a second container started with `--init`. With tini reaping orphans the **whole
`internal/procctl` package is green, including that test (PASS 0.10s)**; without `--init`
it fails identically. That isolates the variable to PID-1 reaping, not to the delta.

**Dependent sites, my container, as a non-root user:** `internal/evidence` ok 9.074s;
`TestRecoveredParentResolutionRequiresExactLineage` PASS 7.86s and
`TestRecoveredParentResolutionGuardsRevalidate` PASS 2.14s (the two hosted U2 sites);
runner `TestVerifierBudgetRefusalRecoversExplicitly` ok; app
`TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts` ok. Host darwin: procctl ok,
evidence ok. **No test was skipped or filtered to produce any green above.**

## Findings

### C-1 — MEDIUM (reasoning; LOW behaviour) · "cancellation arrives as a read error" is false, and that claim is in the repo

`internal/procctl/procctl_cmdline_linux_test.go:4-10` states: *"external cancellation
reaches it only as process death — a read error — which returns immediately instead of
polling (pinned below)"*. `procctl_linux.go:63-64` echoes it ("A read error (process
gone) returns immediately"). The implementer's report builds its whole lifecycle section
on it.

**It is wrong for the dominant case.** I measured a killed-but-unreaped child directly:

```
alive      : len=0 err=<nil>          <- also caught the U2 race live, first sample
zombie(Z)  : len=0 err=<nil>  signal0-alive=true
reaped     : len=0 err=open /proc/19192/cmdline: no such file or directory
```

A dead process only produces a read error **after it is reaped**. Until then it is a
zombie that reads empty *with no error*, so the poll spends the **full 100 ms**. The file
contradicts itself: `TestLinuxZombieCmdlinePollsBoundThenFailsClosed` (same file, ~200
lines below) asserts exactly that the bound is spent on that state. The pinning test the
header cites, `TestPollCmdlineDeadProcessStopsPolling`, uses a *fake* error-returning
reader — it models only the reaped case and cannot detect this.

This matters because the one capture site that *is* a cancellation handler is
`runner/launch.go:74` — `cmd.Cancel = func() { KillGroup(CaptureByPID(...)) }`. Cancel
where the agent has already exited is precisely the zombie case.

**Outcome is still safe and bounded** (≤100 ms, then the pre-fix refusal), so this is a
documentation defect, not a behaviour defect. **Ask:** correct both comments to say
"bounded by `cmdlinePublishBound`; a dead-but-unreaped process reads empty without error
and spends the bound". No code change needed; the correction should not be allowed to
turn into a poll-the-`Attributed`-facets redesign.

### C-2 — LOW · the bound is adequate, but the stated margin is wrong by ~3 orders of magnitude

The report justifies 100 ms as *"three-plus orders of magnitude beyond preemption
jitter"*. I instrumented the window the poll must cover (time from `Start()` returning to
the first non-empty `/proc/<pid>/cmdline`), 10,400 spawns total:

| condition | samples | empty first read | p50 | p99 | p999 | **max** |
|---|---|---|---|---|---|---|
| idle | 2000 | 3.95% | 15 µs | 98 µs | 213 µs | **281 µs** |
| 4× CPU oversubscription | 2000 | 2.45% | 11 µs | 4.0 ms | 37 ms | **56.5 ms** |
| 4× oversubscription, 16 spawners | 3200 | 2.16% | 12 µs | 8.3 ms | 27 ms | **33.9 ms** |
| 8× oversubscription | 3200 | 2.28% | 11 µs | 5.1 ms | 23 ms | **34.2 ms** |

The µs claim holds for the *median*; the **tail under load is milliseconds to tens of
milliseconds**. Worst observed sample is 56.5 ms — the bound has about **1.8× headroom at
the measured maximum, not 1000×**. No sample exceeded the bound (0 of 10,400, none even
at the 5 s ceiling), so **100 ms is supported by this evidence** and exhaustion degrades
to status quo ante — but the hosted runner ("dozens of parallel test packages") is a
noisier environment than my load, and the margin is not what the report says it is.
**Ask:** restate the justification with measured numbers. Raising the bound is defensible
but not required; I do not ask for it.

### C-3 — LOW · the real-caller survey is incomplete, and "only extends paths that already fail" needs one qualifier

The report says *"`durablekill.go` calls `Attributed` once per invocation, no loop"*.
That is true of `KillAgentDurable`, and incomplete for the file and the codebase:

- **`runner/launch.go:74`** — inside `cmd.Cancel`, `WaitDelay = 2s`. Worst case now
  100 ms (poll) + 1500 ms (`KillGroup` grace, pre-existing) = ~1600 ms of a 2000 ms
  budget. **Inside the envelope as claimed**, but the headroom goes 500 ms → 400 ms and
  this site is not named.
- **`app/preflight.go:916`** — `Capture` runs *before* `go cmd.Wait()` is armed, so a
  probe that exits fast is an unreaped zombie at capture time. Measured in that exact
  shape: `sh -c "exit 1"` burned the full bound on **79/100** spawns (mean 79.5 ms),
  `sh -c "true"` on **73/100** (mean 73.4 ms), while long-lived `sleep 5` burned it on
  **0/100** (mean 477 µs). Severity stays LOW because the healthy preflight probe is the
  agent's real invocation (long-lived, unaffected), the affected probes are failing ones,
  and probes run concurrently per agent — so wall-clock cost is ≈ one bound, not N.
- **`durablekill.go:66 AgentLiveness`** is called from TUI *render* code
  (`protocolui.go:353`, `live.go:848`, and `live.go:1305` which loops every agent), not
  once per invocation. Only an agent in zombie state pays, so exposure is narrow — but
  "no loop" is not accurate for the file.
- **`trajectory/verification.go:145`** — `KillTreeAttributed` in an ordinal loop of
  `4×len(criteria)`; it returns on the first refusal, so at most one bound per stop. No
  accumulation. I checked this specifically and it does **not** multiply.

**Ask:** name `launch.go:74` and `preflight.go:916` in the record. No code change.

### C-4 — INFO · the "exhaustive" behaviour table omits the live-probe row

The table covers the *recording* read. `Attributed`'s **live** read polls too, and one
row is missing: an empty live read was previously an **unconditional refusal**
(`"cannot read live command"`); it is now wait-then-compare. **This is not a weakening** —
reaching that facet already requires exact boot id, alive, exact start time, exact pgid
and session-leader match, and whatever the poll returns is handed to the untouched
`commandMatches`. It removes a *timing-accident* refusal and replaces it with the designed
comparison; it grants no capability an attacker did not already have (argv publication is
monotonic for a living pid). Worth one row in the table, not a change.

### C-5 — INFO · pre-existing, no action requested

`procctl.go:160-171`: the comment says `commandMatches` *"does NOT accept a shorter
recorded command matching a longer live one"*, but `strings.HasPrefix(live, recorded)`
accepts exactly that, and the parenthetical about "truncation of the live read" describes
the opposite direction from the code. There is **no test for `commandMatches` at all**
(grep of `procctl_test.go`). Unchanged by this commit and **not exercised differently by
it** — pre- and post-fix both return the first non-empty read, so partial-read exposure is
identical. Recorded only because this fix's safety argument rests on that gate.

## Safety and scope assessment

**No weakening found.** Verified, not taken on trust:

- `git diff 231d889..ef10cf3` touches **three files**: one product file, one new test file,
  one report. `procctl.go`, `procctl_unix.go`, `procctl_darwin.go`, `procctl_windows.go`
  and `procctl_test.go` are **byte-identical** to the parent — so `Attributed`,
  `commandMatches` and `KillGroup` are untouched, as claimed.
- `VERSION`, `CHANGELOG.md`, `internal/version`, `.github` untouched. No version bump, no
  push, tag, release or install. `v1.49.0` (`06e563e`) unmoved. This review is **not**
  a 1.49.1 publication and takes no position on one.
- `GOOS=darwin|windows|linux go build ./...` all ok. The build tag keeps darwin (ps-based)
  and Windows byte-identical; the **owner decision on the Windows track remains open and
  nothing here touches or implies it**.
- The bound cannot convert an absent identity into a present one: exhaustion is the
  pre-fix refusal (F), a value is never re-read (R), a dead pid is never polled (D). All
  three reproduced by me.
- Worst case is bounded everywhere: ≤100 ms, inside `KillGroup`'s 1500 ms grace and every
  `WaitDelay` (2 s at runner/launch, 10 s at evidence) — confirmed by reading the
  constants, not by accepting the report's table.

## Readiness for ONE hosted x86_64 ubuntu validation

**Ready.** Recommend the single untagged run the implementer proposes, over
`internal/procctl`, `internal/evidence`, `internal/trajectory`, `internal/runner`,
`internal/app`.

Pass condition: `"no recorded command"` disappears from the Command facet while every
other refusal facet and every fail-closed behaviour is unchanged. Per my earlier rule,
**a single green run is not by itself proof** — this defect is probabilistic, and 28/31
packages were already green in run 36002574211. What makes the claim credible is the
reproduced mechanism plus the structural guarantee; the hosted leg is confirmation of
those on the real architecture, not the evidence base.

## Limitations

- **L-1 — nobody has run this on x86_64.** My Linux leg is arm64 (the docker host is
  `arm64 linux`), the same architecture class as the implementer's. The execve ordering is
  kernel-generic and the reproduction is plain `forkExec` + `/proc`, but the x86_64
  confirmation genuinely belongs to the hosted run and **neither of us has closed it**.
- **L-2** — my publication-window numbers are from containers on a shared 16-CPU host
  under synthetic `while :` load. They bound the window as I could induce it; they are
  not kernel-instrumented and are not the hosted runner's load profile.
- **L-3** — I did not run the full test suite, and make no claim about packages outside
  those listed. Windows was not examined; no Windows claim is made or implied.
- **L-4** — all mutation work happened inside throwaway containers, both of which are now
  removed; my worktree is clean at `ef10cf3` with both file hashes unchanged. No source
  fix, commit, push, CI rerun, release, install or closed-artifact edit was performed. No
  secrets in this report.

## Unresolved conditions

1. **C-1 comment correction is unresolved** — two in-repo comments assert something my
   measurement disproves, and the test cited as pinning it cannot pin it. This is the only
   finding I would want addressed before the change is treated as reviewed-and-explained,
   and it is a comment edit, not a code edit.
2. **C-2 bound justification is unresolved** — the numeric bound stands on my measurements,
   not on the margin argument in the report. Whether 100 ms is *sufficient on the hosted
   x86_64 runner under real CI load* is exactly what one hosted leg cannot settle either:
   exhaustion is silent (it reproduces today's failure), so an under-sized bound would look
   like U2 recurring rather than like a bound problem. If U2 recurs hosted, **raise the
   bound before re-diagnosing**.
3. **C-3 caller record is unresolved** — `launch.go:74` and `preflight.go:916` should be
   named in the record; neither needs a code change.
4. **U2 remains unconfirmed-as-fixed on the hosted runner.** Everything above is local
   arm64 plus a structural argument.
5. **Windows track is untouched and still owner-blocked**; D-3 diagnostic-site drift from
   my previous report remains open and was deliberately not widened here.
6. **No version/publication decision is made or implied by this review.**
