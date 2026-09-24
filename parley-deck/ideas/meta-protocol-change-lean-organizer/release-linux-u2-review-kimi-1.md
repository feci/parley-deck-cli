---
idea: meta-protocol-change-lean-organizer
author: kimi-1
role: independent reviewer (non-implementer)
artifact: release-linux-u2-review
date: 2026-09-24
reviewed-commit: ef10cf3 (parent 231d889)
inputs: release-linux-hosted-diagnostics-02-claude-1.md, release-linux-u2-repair-zcode-1.md, full source at ef10cf3
environment: go1.27.1 darwin/arm64 (host) + own golang:1.26 container, go1.26.8 linux/arm64 (same toolchain as implementer)
verdict: PASS — merge-ready for exactly ONE untagged hosted x86_64 ubuntu validation leg; not a 1.49.1 publication decision
---

# Linux U2 repair review — kimi-1

## Verdict

| Item | Status | Basis |
|---|---|---|
| Product delta scope | **PASS** | `git diff 231d889..ef10cf3` = exactly `internal/procctl/procctl_linux.go` (+51/−5), the new 301-line linux-tagged test file, and the report. `procctl.go`, `procctl_unix.go`, darwin, windows: zero delta. |
| Provenance | **PASS** | sha256 of both files match the report's table (`a3f8eb16…0547`, `5b3cca33…5b26`); commit parent is 231d889; VERSION/CHANGELOG untouched; closed artifacts untouched; no push/tag. |
| Safety invariants | **PASS** | first-nonempty/no-shopping, permanent-empty fail-closed, read-error immediacy, PID/start/pgid strictness, zombie refusal — all hold; all independently falsified below. |
| Caller/lifecycle envelopes | **PASS with notes** | no real caller's envelope is breached; K-1/K-2 below. |
| Environmental-failure claims | **PASS** | both deviations bidirectionally reproduced as environmental, not product. |
| Bound (100ms/1ms) | **Accepted as engineering parameter** | justified by reproduced efficacy + fail-closed floor; residual recorded under Unresolved. |

## Independent reproduction (my own container, throwaway copies in /tmp, worktree never mutated)

**Post-fix (pristine export of ef10cf3, sha256-verified):**
- Full `internal/procctl` suite: **all 14 tests PASS, zero SKIP** — 9 new + all pre-existing, including `TestKillGroupReapsChild` (my container PID1 is bash, which reaps).
- New tests `-count=3` ok 1.024s; `-race -count=3` ok 2.051s.
- Spawn test `-count=20` = 2400 spawns: **20/20 PASS, 0.460s** (implementer: 20/20, 0.342s).
- Hosted U2 sites: `TestRecoveredParentResolutionRequiresExactLineage` PASS 7.40s, `TestRecoveredParentResolutionGuardsRevalidate` PASS 2.38s; runner `TestVerifierBudgetRefusalRecoversExplicitly` ok; app `TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts` ok.
- `internal/evidence`: **ok 9.029s as uid 1000**; as root exactly the two named permission tests FAIL (`TestSaveUnwritableDirFails`, `TestFailedSaveLeavesNoReport` — both chmod-0555 based; root bypasses permission bits). Claim confirmed verbatim.
- Host: `go build` darwin/linux/windows all ok; `GOOS=linux go vet ./internal/procctl/` clean; `gofmt -l` flags only the two pre-existing drift files the implementer named, not the touched ones; darwin suites `procctl` ok 0.562s, `evidence` ok 5.415s.

**Bidirectional falsification (my own mutations, matching the implementer's table):**
- **W** (pre-fix single-read restored): zombie test FAILs deterministically at count=1 ("poll must engage… took 6.375µs"); spawn race **reproduced: 5 empty captures in 1200 spawns** (~0.4%; implementer 8/1200 ≈ 0.7% — same order, load-dependent). Post-Start empty `/proc/<pid>/cmdline` on a live stat-readable process is real; candidates B/C/D stay excluded (bare `sleep`+`Capture`, no supervisor).
- **F** (exhaustion fabricates): exactly `PermanentEmpty` + `Zombie` FAIL (2).
- **R** (keep re-reading after a value): exactly `HealthyFastPath` + `TransientEmpty` FAIL (2).
- **D** (swallow read error, poll dead pids): exactly `DeadProcessStopsPolling` + `ImmediateError` + `ReapedProcessIsPrompt` FAIL (3).
- **PID1 claim**: `TestKillGroupReapsChild` FAILs identically (6.09s vs 6.08s, "still alive after 3s", defunct processes present) on **pre-fix and post-fix** under PID1=`sleep infinity`, and PASSes under a reaping PID1 (bash). Environmental, pre-existing — proven in both directions. No failing product test was skipped to claim green anywhere.

**My additional falsifiers (not in the shipped test set):**
- First-nonempty with two *distinct* values (empty→A→B): returns A, exactly 2 reads — no value shopping. PASS.
- Re-exec identity transition (`sh -c 'exec sleep 30'`, 60 spawns): 59 attribute, 1 records the pre-exec argv `"sh -c exec sleep 30"` and then refuses fail-closed `"command mismatch (pid was reused)"` against the post-exec live read. The gate never launders a transitioned identity; this window is pre-existing (first-nonempty behaves identically to the old single read on non-empty data) — see K-3.
- Zombie interim (kill, wait for real `Z` state, do not reap): `/proc/<pid>/cmdline` reads **empty, not error**; capture burns the full ~100ms bound (measured 0.10s), records empty Command, refuses `"no recorded command"`. See K-1.

## Dispatch-mandated checks

- **First-nonempty / no value shopping: HOLDS.** Return on first `len>0`, loop never re-reads; pinned twice, my A/B test concurs, mutation R dies.
- **Permanent-empty fails closed: HOLDS.** Exhaustion returns `("", false)` = byte-identical pre-fix outcome; attempt-capped (~≤105), elapsed ∈ [85ms, 5s]; real-/proc zombie pin; mutation F dies.
- **Read-error / death: HOLDS with precision K-1.** `err != nil` → immediate `("", false)`; dead-reaped pids never polled (`ReapedProcessIsPrompt`: elapsed < bound).
- **PID/start/pgid strictness: HOLDS.** `Attributed`/`commandMatches` byte-untouched; exact refusal strings reproduced live; pid swap still dies at start-time/pgid facets before the command facet is reached.
- **Cancellation/deadline bounds: HOLD** for every real caller (assessed individually, below).
- **Zombie: refuses at the same facet with the exact pre-fix reason** (`"cannot read live command"` for Attributed live-read; `"no recorded command"` for a zombie capture) after spending the bound. Reproduced.
- **Healthy path: zero added cost.** Exactly one read, no sleep; full suite 0.545s vs 0.107s hosted pre-fix baseline is noise-level.

## Real-caller assessment (the dispatch's "do not accept cancellation-means-death")

- `execute.go:162-163` (hosted U2 site): poll runs inside `start()` under the control closure; publish-before-release serialization untouched; on refusal the path already paid the 1500ms KillGroup grace (hosted signature 1520/1521ms) — the poll adds ≤100ms only to paths that already fail. ctx cancel → `cmd.Cancel` → KillGroup directly (no probe). ✓
- `runner.go:1113`: capture sits between `ctx.Err()` check and `waitSupervised`; worst case supervision entry is delayed ≤100ms in the empty case, then ctx watch + 1500ms grace govern. ✓
- `acp.go:118`: one-shot after ACP spawn, post-Start so post-fork-copy; established agents read once. ✓
- `preflight.go:916`: poll eats ≤100ms of an absolute probe deadline, only in the empty case. ✓
- `launch.go:74` (`cmd.Cancel` → `KillGroup(CaptureByPID(...))`): the one place the poll lands on a **cancellation** path — ≤100ms of probing before TERM (the probed identity is discarded; KillGroup needs only PID/PGID — pre-existing wasted work, mildly amplified). Worst case 100ms + 1500ms grace = 1600ms < `WaitDelay` 2s. Bounded, fail-safe direction. ✓
- `durablekill.go:33` (one-shot), `durablekill.go:66` + TUI (`AgentLiveness`), `verification.go:145` (`KillTreeAttributed` cleanup): each pays ≤100ms per alive-but-empty-cmdline process. TUI liveness is called synchronously on the Bubble Tea update/render path — see K-2.

## Findings

- **K-1 (LOW, prose-vs-behavior):** "cancellation arrives as process death = read error = immediate exit" is exact only for **reaped** pids (ENOENT). A killed-but-unreaped process (zombie interim) reads as **permanent-empty**, so the poll burns the full 100ms bound before failing closed — measured 0.10s in my zombie-interim falsifier. The implementer's own zombie test already encodes this; the prose is looser than the tests. Outcome unchanged: fail closed, ≤100ms, every envelope intact. No code change required; the framing should not be quoted as-is into any release note.
- **K-2 (LOW, scope note):** placing the poll at the shared choke point also slows `Attributed`'s live read by ≤100ms for alive-but-empty processes (zombies, self-argv-zeroing). `KillAgentDurable` is one-shot; the TUI calls liveness synchronously per render. Parley-owned zombies are transient (`cmd.Wait`/init reaps), so persistent cost requires a persistent argv-less process — rare and already un-attributable pre-fix. I judge the choke-point placement **correct**: the execve window is symmetric (the live read at `execute.go:163` can hit it microseconds after capture does), and fixing only the recording side would leave the mirrored flake. This is timing-only scope, not a broader safety redesign: no facet, order, strictness, goroutine, signal, or control-contract change.
- **K-3 (LOW, observed, pre-existing):** for a process that re-execs, the first non-empty read can record a **pre-exec** argv (observed 1/60). Attribution then refuses fail-closed (`"command mismatch (pid was reused)"`). Identical window pre-fix; the poll neither widens nor narrows it. Safe direction; noted so the hosted leg isn't misread if it ever surfaces.

## Bound assessment (engineering parameter, per the ruling)

100ms/1ms is justified by evidence, not taste: the window is post-`begin_new_exec` residual work (µs) plus preemption; pre-fix reproduction is ~0.4–0.7%/spawn, and post-fix 0/2400 here + 0/2400 reported (if the true residual rate were still ~0.5%, P(0/2400) ≈ 6e-6). The bound sits 15× under the KillGroup grace and 20–100× under both WaitDelays; exhaustion degrades to exactly the pre-fix refusal — never to a wrong attribution. It is margin-based, not kernel-instrumented; the residual risk is a too-small bound under extreme x86_64 load, which would resurface as the *same diagnosable U2 signature* (the D-1 diagnostics label persists it), not as unsafe behavior. Concur: not an owner gate.

## Readiness

**Ready for ONE untagged hosted x86_64 ubuntu run** over `internal/procctl`, `internal/evidence`, `internal/trajectory`, `internal/runner`, `internal/app`. Expected: `"no recorded command"` disappears from the Command facet; every other facet/refusal stays byte-identical; `TestKillGroupReapsChild` green under real init. Per claude-1's rule, that run is x86_64 *confirmation* of an arm64-reproduced mechanism plus the structural guarantee — not proof by itself. This review is not a v1.49.1 publication decision.

## Unresolved conditions

1. Hosted x86_64 leg outstanding — all local Linux evidence (mine and the implementer's) is arm64; the mechanism is kernel-generic but unconfirmed on x86_64.
2. The 100ms bound is margin-based, not kernel-instrumented; a bound undersized for an extreme hosted load degrades to status-quo U2 failures (visible via the persisted diagnostics label), never to unsafety. If U2 recurs hosted, re-measure before widening.
3. Owner Windows scope decision remains pending; this delta keeps windows/darwin byte-identical (cross-compile verified) and makes no Windows claim.
4. K-1/K-2/K-3 stand as documented LOW notes; none block the hosted leg.
5. No full-suite rerun exists post-fix (implementer's delta-scoped strategy; my container ran procctl full + dependent packages only). The hosted leg should be read as the full-suite checkpoint for this delta.

No source edits, commits, pushes, CI reruns, release, install, or closed-artifact edits were performed. All mutation trees and containers were throwaway under /tmp and docker; the worktree is untouched.
