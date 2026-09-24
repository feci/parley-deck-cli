---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent Linux CI diagnostic reviewer (non-implementer)
artifact: release-linux-hosted-diagnostics-02
date: 2026-09-24
run: 36002574211 (Tests, push, main@231d889) — ubuntu job 107642571416
verdict: U2 ROOT-CAUSED — process attribution, not journal persistence
---

# Linux hosted diagnostics 02 — claude-1

## Verdict

| Item | Status | Basis |
|---|---|---|
| **U2 classification** | **process attribution (candidate A)** — journal persistence (B) **excluded** | diagnostic string, below |
| **U2 root cause** | **ESTABLISHED** to a named probe: `linuxProbe.command` reads `/proc/<pid>/cmdline` as **zero bytes** on a live, stat-readable process | source + envelopes |
| Causal confidence | **HIGH** that the failing facet is the cmdline probe; **MEDIUM-HIGH** for the execve argv-publication window as its mechanism | see Confidence |
| U1 acp stderr drain | **PASS** (2nd consecutive hosted green) | `ok internal/acp 2.660s` |
| The diagnostic delta itself | **worked as designed** | `ok internal/evidence 10.827s`, `ok internal/procctl 0.107s` |
| Windows | **not examined, no claim** | per dispatch, ubuntu only |

Ubuntu job 12:59:06Z → 13:02:54Z, `failure`. **28 packages `ok`, 3 `FAIL`** (`app`,
`runner`, `trajectory`). Single bounded foreground wait; no rerun.

## The diagnostic string — U2 is attribution refusal

Both dump sites in `internal/trajectory` printed **byte-identical** text:

```
"diagnostics": "run error (no command output): criterion supervisor identity unavailable: no recorded command"
```

That resolves the (A)/(B) question my previous report left open:

- **(A) process attribution — CONFIRMED.** `execute.go:163-166` wraps
  `procctl.Attributed(sp)`; `"no recorded command"` is the facet at
  `procctl.go:126-128` (`s.Command == ""`).
- **(B) journal persistence — EXCLUDED.** `writeVerificationArtifact("process-NNN.json")`
  (`verification.go:620`) is **never reached**: `start()` returns the error first
  (`verification.go:616-618`). Consistent with `process-NNN.json` being absent at both sites.
- (C) control-protocol and (D) `withVerification` postamble strings — **not observed**.

Envelope signature at both sites, unchanged from run 35998165067 except now labelled:

| field | site `:142` | site `:510` |
|---|---|---|
| failing **ordinal** | **1** | **4** |
| `exit_code` | `-1` | `-1` |
| `duration_ms` | `1521` | `1520` |
| `output_sha256` | `e3b0c442…7852b855` (empty string) | same |
| `complete` | `false` | `false` |
| `process-<ordinal>.json` | absent | absent |
| retained | `receipt`, `step-001` | `receipt`, `step-001..004`, `process-001..003` |

`TestRecoveredParentResolutionRequiresExactLineage/replacement-lifecycle-moved` and
`TestRecoveredParentResolutionGuardsRevalidate/observation-deleted`.

## Why it is a probe race, not a logic bug

`"no recorded command"` is a **recorded**-identity facet, so the failure happened in
`procctl.CaptureByPID` (`procctl.go:82-98`), not in the live comparison. `Attributed`
checks facets in order, so reaching the `Command` check proves the three before it
**succeeded on the same pid, microseconds apart**:

| probe | mechanism | result |
|---|---|---|
| `bootID()` | `/proc/sys/kernel/random/boot_id` | ok |
| `pgid()` | `getpgid(2)` syscall | ok (> 0) |
| `procStart()` | `/proc/<pid>/stat` field 22, ≥ 20 fields parsed | **ok** |
| `command()` | `/proc/<pid>/cmdline` | **empty → `("", false)`** |

`linuxProbe.command` (`procctl_linux.go:48-55`) fails on `err != nil || len(data) == 0`.
The read-error branch is excluded twice over: `/proc/<pid>/stat` was readable at the
adjacent line, and the process then absorbed a **full 1500 ms `KillGroup` grace**
(`procctl_unix.go:41`; `duration_ms` 1520/1521 vs a 24 ms healthy baseline) — it was alive
and already past `trap ':' TERM`, line 1 of `criterionSupervisor` (`execute.go:101`).
So the branch is **`len(data) == 0` on a live, non-zombie, fully procfs-visible process.**

**Two different ordinals (1 and 4) in one run, different tests, identical facet** — with
`process-004.json` present-then-absent across runs. That is a per-spawn probabilistic
race, not state-dependent logic. It also explains every prior observation: load
sensitivity, migration between sibling tests, and the previously unexplained `-1`.

**Mechanism (inference).** Linux publishes argv late in `execve`: the CLOEXEC exec-status
pipe that unblocks Go's `forkExec` is closed in `begin_new_exec()`, while `mm->arg_start/
arg_end` are set afterwards in `create_elf_tables()`. `proc_pid_cmdline_read()` returns
**0 bytes** while `mm->arg_end` is unset. So `cmd.Start()` can return before argv is
visible, and `Capture` — called immediately after, by design — can read the empty window.

## Confidence

- **HIGH** — U2 is attribution refusal on the `Command` facet, and `Capture`'s cmdline
  probe returned empty on a live process. Directly evidenced by the persisted string plus
  `Attributed`'s check order; no inference about kernel internals required.
- **MEDIUM-HIGH** — the execve argv-publication window as the *reason* the read was empty.
  Reached by elimination (zombie, kernel thread, dead pid, self-zeroed argv all excluded)
  plus known kernel ordering. **Not instrumented on the runner.** The remaining gap is
  *why* the read was empty, not *that* it was.

## Repair — and a design question that should be answered first

`Attributed` is a **fail-closed safety gate** against signalling a reused pid. `Command` is
one of its facets. **The correct repair makes *recording* reliable; it must not relax the
*comparison*.** Two shapes:

- **R1 (recommended):** bound a short retry around the cmdline read at capture time, so
  `Capture` does not record an empty `Command` for a process it just started. The
  `Attributed` gate is untouched — recorded identity must still match the live probe.
- **R2 (reject):** making `Command` optional, or treating empty as "skip". This weakens the
  PID-reuse gate and must not be done.

**I am flagging a design call rather than handing over an implementation.** R1 inserts a
wait into the critical section between `Start()` and the kill-capable identity record — for
a supervisor that is already parked on `read`, and on a path whose whole purpose is durable
kill. The bound, and whether waiting there is acceptable at all versus capturing identity
behind the existing handshake, is an owner/product safety judgement. I have the evidence to
name the defect precisely; I do not have the evidence to pick that bound, and I decline to
guess one. Note the same probe is reachable from the ACP path via `CaptureByPID`, and the
symmetric live-read facet (`"cannot read live command"`) is the same class, unobserved so far.

**Validation.** Deterministic proof needs Linux: a `//go:build linux` test that spawns many
short-lived children and asserts `Capture` never yields an empty `Command`. It should
reproduce pre-fix under load and pass post-fix. macOS cannot validate this — `darwinProbe`
uses `ps`, a different mechanism with no such window. Do **not** treat a single green
ubuntu run as proof; this defect is probabilistic and 28/31 packages were already green here.

## Residual — dump coverage drifts with the failure

The `app` and `runner` failures produced **no** diagnostic string: they surfaced at sibling
sites (`trajectory_parent_recovery_test.go:200`, `verification_bound_recovery_test.go:233`)
rather than the wired ones (`:153`, `:229`). Same incompleteness family, no envelopes. This
is D-3 (site drift) realised, and it means the wired dumps are not a stable net. I record
it as an observation; widening them is not required now that the facet is named.

## Next action

**Owner/organizer decision on the R1 design question above, then one repair + one hosted
ubuntu leg.** Not a rerun of this tree: the defect is named, and re-running only resamples
the race. I do **not** recommend repeat-until-green.

## Limitations

- Ubuntu job `107642571416` of run `36002574211` only, read once. No local reproduction of
  U2 exists on any machine; the runner was not instrumented.
- The kernel-ordering mechanism is inference by elimination, stated at MEDIUM-HIGH, not a
  measurement. The facet identification does not depend on it.
- macOS leg still `in_progress`; the Windows leg completed `failure` and I **did not open
  it** — per dispatch this is ubuntu-only and I make no Windows claim. Nothing here waives
  the owner-blocked Windows track. D1 remains open and untested on Linux (acp passed).
- No source edit, rerun, push, tag, release, install or version bump was performed.
  `v1.49.0` (`06e563e`) is unmoved; closed A–D artifacts untouched. No secrets in this report.
