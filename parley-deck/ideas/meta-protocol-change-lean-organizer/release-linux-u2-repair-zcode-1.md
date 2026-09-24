---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: narrow Linux release repair implementer (U2 recording fix)
artifact: release-linux-u2-repair
date: 2026-09-24
base: 231d889
inputs: release-linux-hosted-diagnostics-02-claude-1.md (U2 root cause + R1/R2), release-linux-diagnostic-followup-zcode-1.md, release-linux-diagnostic-review-claude-1.md, release-linux-repair-followup-zcode-1.md, release-linux-u2-review-claude-1.md, release-linux-u2-review-kimi-1.md
amended: 2026-09-24 (post-review corrections appended; body above left verbatim as reviewed at ef10cf3)
environment: go1.27.1 darwin/arm64 (host) + golang:1.26 container go1.26.8 linux/arm64 (local Linux verification, never organizer tests)
verdict: U2 REPAIRED AT THE RECORDING LAYER — bounded, fail-closed /proc cmdline poll; race reproduced locally pre-fix (first reproduction anywhere); every attribution facet and refusal preserved
---

# Linux U2 repair — zcode-1

Dispatched scope: the smallest safe **recording** fix for claude-1's root cause —
`linuxProbe.command` reads `/proc/<pid>/cmdline` as zero bytes inside the execve
argv-publication window, so `CaptureByPID` records an empty `Command` for a live
process it just started, and `Attributed`'s strict fail-closed `"no recorded
command"` facet refuses it (hosted signature: exit −1, empty output, empty hash,
~1500 ms kill grace, ordinals 1 and 4 of run 36002574211). `Attributed`'s
comparison must not be weakened; the organizer ruled the bound an engineering
parameter for participants to justify with evidence, not an owner permission gate.

**Scope check: the whole product delta is one hunk in one build-tagged file.**
`internal/procctl/procctl_linux.go`, +51/−5: the single-read `command` becomes a
bounded poll. No portable file, no probe interface, no supervisor script, no
control contract, no test scaffolding, no workflow, no version file is touched.

## New evidence: the race reproduces locally, pre-fix

Before designing blind, I reproduced claude-1's mechanism in the local Linux
container (golang:1.26, go1.26.8, linux/arm64). With the pre-fix single-read
`command()` restored as mutation W and the spawn test repeated 10× (1200
`sleep` spawns: 40+20·4 sequential+parallel per iteration):

```
procctl_cmdline_linux_test.go:170: capture recorded an empty command for just-started pid 26325
… 8 failures total: 6 sequential (:170), 2 parallel (:198), ~0.7% of spawns
```

This is the **first local reproduction of U2 on any machine** (claude-1: "No
local reproduction of U2 exists on any machine"). It confirms the facet
identification directly — `Capture` on a just-started, live, stat-readable
process can read zero cmdline bytes — and independently elevates claude-1's
MEDIUM-HIGH execve-window *mechanism* from inference-by-elimination to
observed-locally: the container has no journal/controller confounds (candidates
B/C/D excluded by the same strings claude-1 used, plus the fact that a bare
`exec.Command("sleep")`+`Capture` with no supervisor reproduces it).

Post-fix, the identical load is deterministically green: `-count=20` (2400
spawns) → **20/20 PASS, 0.342 s**. That is the structural guarantee, not a
probabilistic pass: the poll's only exits are publication-within-bound (record
the real value), read error (record empty, exactly pre-fix), or bound
exhaustion (record empty, exactly pre-fix).

## The change

```go
func (linuxProbe) command(pid int) (string, bool) {
	return pollCmdline(func() ([]byte, error) {
		return os.ReadFile("/proc/" + strconv.Itoa(pid) + "/cmdline")
	})
}

func pollCmdline(read func() ([]byte, error)) (string, bool) {
	deadline := time.Now().Add(cmdlinePublishBound)
	for {
		data, err := read()
		if err != nil {
			return "", false
		}
		if len(data) > 0 {
			return strings.TrimSpace(strings.ReplaceAll(string(data), "\x00", " ")), true
		}
		if !time.Now().Before(deadline) {
			return "", false
		}
		time.Sleep(cmdlinePublishPoll)
	}
}
```

with `cmdlinePublishBound = 100 ms`, `cmdlinePublishPoll = 1 ms`, and the
mechanism + bound documented in-file at the constants. Placement rationale:

- **This is the single shared choke point of the evidenced defect.** Every
  Linux capture path funnels through it — evidence's criterion supervisor
  (`execute.go:162`, the hosted U2 site), the runner's headless agent capture
  (`runner.go:1113`), the ACP path claude-1 flagged as same-class
  (`acp.go:118`, `launch.go:74`), and preflight (`preflight.go:916`) — and so
  does `Attributed`'s live read, whose empty-read facet ("cannot read live
  command") is the symmetric case claude-1 named "unobserved so far".
- **Build-tag scoping keeps every other platform byte-identical.** `git diff`
  shows exactly one product file changed; darwin (`ps`-based, no such window
  per claude-1) and Windows (owner-blocked track) are untouched by
  construction, and both still cross-compile (`GOOS=darwin/linux/windows go
  build ./...` all ok).

Behavioral delta, exhaustively (the only rows that change involve an
empty-read of a live-looking process; everything else is byte-identical logic):

| state of the read | before | after |
|---|---|---|
| first read non-empty (every healthy path) | return value | return value — **zero added latency, one read** |
| read error (process gone / unreadable) | `("", false)` now | `("", false)` now — dead pids are never polled (test-pinned) |
| empty read, then argv publishes (the U2 window) | `("", false)` — **the bug** | the published value, within the bound (test-pinned) |
| empty read, never publishes (bound exhausted) | `("", false)` | `("", false)` after ≤100 ms — **exactly the pre-fix outcome** (test-pinned) |
| non-empty value that does not match the record | refusal, unchanged | refusal, unchanged — **the poll never sees a mismatch**: it returns any non-empty value immediately and `Attributed` compares it strictly as before (test-pinned, exact reason) |

## Bound choice — evidence, not taste

The window is the remainder of `execve` after `begin_new_exec()` closes the
CLOEXEC exec-status pipe (which unblocks `cmd.Start()`): the kernel still has to
run `create_elf_tables()` to set `mm->arg_start/arg_end`. That is microseconds
of CPU; the variable is preemption under load, which is exactly when the hosted
runner (dozens of parallel test packages) saw it.

- **100 ms** is three-plus orders of magnitude beyond preemption jitter, and
  ~0.7%-of-spawns reproduced even on a lightly-loaded container with the
  immediate read — the window is real but shallow.
- It sits far inside every surrounding budget, so it cannot push any caller
  past an existing envelope: KillGroup grace **1500 ms**, acp WaitDelay **2 s**,
  evidence WaitDelay **10 s**. The only sequences that can observe the bound
  are ones that already failed (refusal → kill → exit −1), extended by ≤100 ms.
- **1 ms** poll interval: ≤ ~101 attempts of a three-syscall /proc read —
  negligible cost, no CPU spin, and small enough that exhaustion tracks the
  wall bound (test-pinned: ≥85 ms and ≤ ~105 attempts).
- This is not a sleep hiding a failure: after the bound the behavior is
  **identical to today's failure** — empty `Command` recorded, `Attributed`
  refuses `"no recorded command"`, the run fails with the diagnostics label my
  previous change persists. The bound converts a *transient kernel publication
  delay* into a recorded identity; it never converts an *absent* identity into
  a present one.

## Why this cannot turn untrusted identity into trusted identity

1. **Retry-on-empty only, return-on-first-value.** The loop re-reads solely
   when the kernel returned exactly zero bytes for a live-looking process. The
   instant any bytes exist it returns them and never reads again — it cannot
   shop for a different value. Mutation R (keep re-reading after a value) is
   failed by two tests specifically to pin this property.
2. **The comparison is untouched.** `Attributed`'s facet set, order, and
   strictness are byte-identical: boot id, pid alive, **exact** start time,
   **exact** process group, session-leader, command match (live extending
   recorded only). A value the poll waits for must still pass every one —
   argv publication is monotonic for a living pid, so what the poll waits for
   is the kernel publishing the process's *own* argv, not an attacker
   substituting a better match; a pid swap is caught by the start-time/pgid/
   session-leader exact comparisons that surround the command facet, exactly
   as before.
3. **Every failure branch fails closed, test-pinned to the pre-fix outcome.**
   Exhaustion → `("", false)` (mutation F, fabricating a value on exhaustion,
   fails two tests). Read error → immediate `("", false)`, no polling of a
   dead pid (mutation D fails three tests). A zombie — dead, unreaped, alive
   to signal-0, stat/pgid still matching — refuses at the live-command facet
   with the **exact pre-fix reason** `"cannot read live command"`, after
   spending the bound polling real `/proc` (test-pinned; this is also the
   mutation-W discriminator).
4. **PID reuse stays rejected.** Tampered command on a live captured process
   refuses `"command mismatch (pid was reused)"`; tampered start time refuses
   earlier at `"process start time mismatch (pid was reused)"` — exact strings
   pinned on Linux with the polling probe active, alongside the pre-existing
   `TestAttributedSelfAndRefusals`/`TestAttributedFailsClosedOnMissingFields`
   (all still green).

## Cancellation, deadline, lifecycle

The probe interface is deliberately context-free (portable, frozen); so the
poll takes no context either — external cancellation can only reach it as
process death, which arrives as a read error and exits immediately (pinned by
`TestPollCmdlineDeadProcessStopsPolling` + `TestLinuxCaptureByPIDOfReapedProcessIsPrompt`).
Concretely on the criterion path: ctx cancellation fires `cmd.Cancel` →
KillGroup → TERM (trapped by the supervisor, line 1 of `criterionSupervisor`)
→ KILL at the 1500 ms grace → death → ENOENT → the poll, if running, exits at
its next iteration. Worst case bound (100 ms) + grace (1500 ms) stays inside
the pre-existing envelope; the healthy path adds zero reads. No goroutine, no
defers, no signals, no `Attributed`/`KillGroup`/`CriterionStartControl`
signature is touched — the publish-then-release serialization of the control
contract is exactly as reviewed.

Live-read latency note for reviewers: `Attributed` (reattach/durablekill/TUI
liveness) now spends ≤100 ms extra **only** when refusing an alive-but-empty-
cmdline pid (zombie-like states); every established process fast-paths on the
first read. `durablekill.go` calls `Attributed` once per invocation, no loop.

## Alternatives, and why rejected

- **R2 (claude-1's reject): optional/empty-tolerant `Command` facet.** Weakens
  the PID-reuse gate. Not implemented; `Attributed` is byte-untouched.
- **Capture behind the supervisor handshake** (wait for supervisor liveness
  proof before recording): requires inverting the control contract's
  publish-identity-before-release serialization — that ordering *is* the
  durable-kill safety being protected — or adding a readiness pipe and
  changing `criterionSupervisor` + both control test scaffolds: a broader
  redesign of a reviewed protocol for no additional safety over the poll.
  Also the supervisor's own stdout is the hashed capture sink — a readiness
  write there would corrupt the output binding.
- **Polling another kernel observable** (`stat` state, `wchan`): nothing
  reliable marks "argv published"; `/proc/<pid>/cmdline` is itself the datum
  whose emptiness defines the window.
- **Portable retry in `CaptureByPID`:** would change darwin and Windows
  recording behavior; the evidenced race is Linux-execve-specific and the fix
  is scoped to it.
- **Fix at `execute.go` only:** same race class exists at the four other
  capture sites; the probe is the one shared choke point.

## Tests — `internal/procctl/procctl_cmdline_linux_test.go` (new, 301 lines, `//go:build linux`)

Nine tests; the dispatch's required coverage maps as:

| required case | test(s) |
|---|---|
| transient empty cmdline | `TestPollCmdlineRecoversTransientEmpty` (empty,empty,value → ok, 3 reads, parse) |
| permanent missing command | `TestPollCmdlinePermanentEmptyFailsClosedWithinBound` (fail-closed after exactly the bound, attempt-capped) + `TestLinuxZombieCmdlinePollsBoundThenFailsClosed` (real `/proc`, exact refusal facet) |
| changed identity / PID reuse | `TestLinuxAttributedRefusesIdentityChangesUnderPollingProbe` (exact refusal strings, both facets) + pre-existing self/refusal suite re-run green on Linux |
| cancellation / bound | `TestPollCmdlineDeadProcessStopsPolling`, `TestPollCmdlineImmediateError`, `TestLinuxCaptureByPIDOfReapedProcessIsPrompt` (death is the cancellation; prompt, no bound burned) |
| healthy fast path | `TestPollCmdlineHealthyFastPath` (exactly one read, zero wait) + `TestLinuxCaptureOfStartedProcessRecordsCommandEveryTime` (120 spawns; claude-1's validation shape) |

The zombie test is the deterministic real-`/proc` stand-in for "command never
publishes": a killed-but-unreaped child is alive to signal-0, its
stat/pgid still match every earlier facet, and its cmdline reads permanently
empty — so it proves both that the poll truly engages on real `/proc` (elapsed
≥85 ms) and that the gate still refuses at the same facet with the same reason.

### Bidirectional proof (claude-1's method; all inside the container, restored and sha256-verified after)

| mutation | expected discriminators | observed |
|---|---|---|
| **W** — pre-fix single-read `command()` | zombie bound test; (probabilistically) spawn test | count=1: only zombie FAIL; count=10: spawn test FAIL 8/10 iterations with `capture recorded an empty command` — **the race itself, reproduced** |
| **F** — exhaustion fabricates `("unknown", true)` | fail-closed pins | `PermanentEmpty…` + `Zombie…` FAIL |
| **R** — keep re-reading after a value (value shopping) | first-read-return pins | `HealthyFastPath` + `TransientEmpty` FAIL |
| **D** — drop the early error return (poll a dead pid) | promptness pins | `DeadProcessStopsPolling` + `ImmediateError` + `ReapedProcessIsPrompt` FAIL |

Restore verified: container file sha256 `a3f8eb16…040547` == host file, then
full new-test set green again. Post-fix load: spawn test `-count=20` = 2400
spawns, 20/20 PASS; `-count=3 -race` PASS (1.399 s).

## Validation

Host (go1.27.1 darwin/arm64): `go build ./...` ok; `go vet ./internal/procctl/`
clean; `gofmt -l` clean on both touched files (pre-existing drift in untouched
`procctl_test.go`/`procctl_windows.go` left alone, as ever); `go test
./internal/procctl/` ok 2.025 s; `go test ./internal/evidence/` ok 5.556 s;
`GOOS=linux` and `GOOS=windows` `go build ./...` ok (baseline green before the
change too — other platforms provably unaffected).

Container (golang:1.26, go1.26.8 linux/arm64, local Linux verification — never
organizer tests): full `internal/procctl` green **except**
`TestKillGroupReapsChild`, and the whole `internal/evidence` package green as
non-root — with both deviations proven **pre-existing and environmental**, not
mine:

- `TestKillGroupReapsChild` fails **identically on pre-fix code** in this
  container (revert-run shown): the container's PID 1 is `sleep infinity`,
  which never reaps orphaned zombies, so the TERMed grandchild stays a zombie
  that `unixAlive`'s signal-0 correctly reports; the killed pids were verified
  state `Z` in `/proc`. Hosted ubuntu (real init) passed `internal/procctl` in
  run 36002574211.
- `TestSaveUnwritableDirFails`/`TestFailedSaveLeavesNoReport` fail as
  **root** (permission bits don't stop root); the entire evidence package is
  `ok 8.834 s` when run as uid 1000.

Dependent hosted-failure sites, container, post-fix:
`TestRecoveredParentResolutionRequiresExactLineage` PASS 7.52 s,
`TestRecoveredParentResolutionGuardsRevalidate` PASS 2.20 s (the two hosted
U2 sites), runner `TestVerifierBudgetRefusalRecoversExplicitly` ok 0.902 s,
app `TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts` ok 1.550 s.
The evidence package exercises the criterion-supervisor `Capture→Attributed`
path through the real Linux probe on every controlled run.

No full-suite rerun (last full evidence: release-ci-fix-zcode-1.md; delta is
one probe function). Diagnostic-site drift (claude-1's D-3 residual) **not
widened**: with the recording layer fixed, the dumps' evidence job is met;
the sites stay as reviewed.

## NOT done

- No `Attributed`, `commandMatches`, `KillGroup`, `CriterionStartControl`,
  `criterionSupervisor`, or probe-interface change — the safety gates are
  byte-identical (`git diff` = one file, `internal/procctl/procctl_linux.go`).
- No Windows/darwin change (their files untouched; owner-blocked Windows
  track unaffected), no U3, no W1–W10, no `.github`/workflow, no
  skip/filter/retry/pass-criteria change anywhere.
- No version bump: VERSION/version.go/CHANGELOG untouched — per protocol this
  product delta publishes only as a new immutable version, organizer/owner
  decision; `v1.49.0` (`06e563e`) never moves. No push, tag, release,
  install. Closed FINAL/IMPLEMENTATION/review/consensus untouched;
  claude-1's report, the ledger and the runs/ dirs stay untracked for the
  organizer.

## Limitations

- Local Linux is arm64; the hosted ubuntu leg is x86_64. The mechanism
  (execve ordering) is kernel-generic and the reproduction is plain
  `forkExec`+`/proc`, but the x86_64 confirmation belongs to the hosted run.
- Reproduction rates are load-dependent: ~0.7%/spawn here vs two visible
  failures per hosted run; a single green hosted run still must not be read
  as proof by itself (claude-1's rule) — the claim rests on the reproduced
  mechanism plus the structural guarantee, of which the hosted leg is the
  confirmation.
- The 100 ms bound is margin-based (µs of residual exec work vs preemption
  jitter), not instrumented against kernel timings; exhaustion is designed to
  reproduce today's failure exactly, so an under-sized bound degrades to
  status-quo-ante, never to unsafety.

## Hosted next step (organizer)

Independent delta review of this commit, then one untagged hosted ubuntu run
on the touched and dependent packages (`internal/procctl`, `internal/evidence`,
`internal/trajectory`, `internal/runner`, `internal/app`). Expected: the
`"no recorded command"` signature disappears from the Command facet while
every other refusal facet — and every fail-closed behavior — remains exactly
as reviewed.

## Provenance

| File | Delta | sha256 |
|---|---|---|
| `internal/procctl/procctl_linux.go` | +51/−5 | `a3f8eb16067a949b603a343a7512cb5a9c93f77c8b6660f18b0341cdf5840547` |
| `internal/procctl/procctl_cmdline_linux_test.go` | new, 301 lines | `5b3cca33da5f15d34e8aa557f7437be085309724c565a22495db908cd2935b26` |

Throwaway artifacts reverted and verified: the four mutation variants lived
only in the container and in `/tmp/mut_*.go` (deleted); the pre-fix stash
`/tmp/procctl_linux_prefix.go` (deleted); the container `parley-u2` (removed).
No secrets in this report.

## Post-review corrections (2026-09-24, zcode-1)

Both independent reviews (`release-linux-u2-review-claude-1.md`,
`release-linux-u2-review-kimi-1.md`) approve ef10cf3's behavior for the hosted
leg; the organizer recorded them (9134c7a). Their findings C-1/C-2/C-3 and
K-1/K-2/K-3 are defects in *this report's reasoning and my in-repo comments*,
not in the gate. Corrections below preserve history: the body above stays
verbatim as reviewed; the same facts are corrected in the two source files'
comments by a comment-only commit (equivalence proven below). **All
measurements cited here are the reviewers', reproduced from their reports —
I did not run them; hosted x86_64 leg results belong to the independent
reviewer and are not claimed here.**

### Fact-correction map

| # | original claim (this report) | correction | finding |
|---|---|---|---|
| 1 | "external cancellation can only reach it as process death, which arrives as a read error and exits immediately (pinned by `TestPollCmdlineDeadProcessStopsPolling` + `TestLinuxCaptureByPIDOfReapedProcessIsPrompt`)" (Cancellation, deadline, lifecycle) | Exact **only for reaped pids** (ENOENT). A killed-but-**unreaped** process (zombie interim) reads empty with **no error**, so the poll is **bounded, not immediate**: it spends the full 100 ms and then fails closed exactly as pre-poll. The two cited tests model only the reaped case (fake error reader / actually-reaped pid); the zombie test pins the bound-spending case — the file's own tests already encoded what the prose overclaimed. claude-1 measured it directly (`zombie(Z): len=0 err=<nil>`), kimi-1 measured the full bound burned (0.10 s). | C-1, K-1 |
| 2 | "100 ms is three-plus orders of magnitude beyond preemption jitter" (Bound choice) | True at the **median** only (~15 µs p50). The **tail under CPU oversubscription is milliseconds to tens of milliseconds**: claude-1 instrumented the Start→first-non-empty window over 10,400 spawns — idle max 281 µs, loaded maxima 33.9/34.2/56.5 ms, **0 samples past the bound** — so 100 ms carries **~1.8× headroom at the measured maximum, not ~1000×**. The bound remains supported by those measurements (no sample exceeded it; exhaustion degrades to the status-quo refusal), and claude-1's unresolved-condition note stands: exhaustion is silent, so if U2 recurs hosted, raise the bound before re-diagnosing. | C-2 |
| 3 | "the ACP path claude-1 flagged as same-class (`acp.go:118`, `launch.go:74`)" named the sites but not the cancellation-path arithmetic; Live-read note said "`durablekill.go` calls `Attributed` once per invocation, no loop" | Three completions. **`launch.go:74`** is the one capture site on a *cancellation* path (`cmd.Cancel` → `KillGroup(CaptureByPID(…))`, `WaitDelay` 2 s at `launch.go:63`): worst case 100 ms poll + 1500 ms grace ≈ 1600 ms of 2000 ms — inside the envelope as claimed, but headroom drops 500 ms → 400 ms (unnamed before). **`preflight.go:916`** captures *before* `go cmd.Wait()` is armed (`:918`), so a fast-exiting probe is an unreaped zombie at capture time: claude-1 measured the full bound burned on **79/100** `sh -c "exit 1"` (mean 79.5 ms) and **73/100** `sh -c "true"` (mean 73.4 ms), vs **0/100** long-lived `sleep 5` (mean 477 µs) — the healthy preflight probe (the agent's real, long-lived invocation) is unaffected, failing probes pay ≈ one bound wall-clock (probes run concurrently per agent). **`durablekill.go`** is not loop-free as a file: `KillAgentDurable` (`:19`) is one-shot, but `AgentLiveness`'s `Attributed` call (`:66`) is invoked from TUI render/update paths — `protocolui.go:353`, `live.go:848`, and `live.go:1305` (which loops every agent) via `live.go:1290` → `app.go:2416` (`AgentLivenessAt`). `trajectory/verification.go:145` calls `KillTreeAttributed` in an ordinal loop (4×len(criteria)) but returns on the first refusal — at most one bound per stop, no accumulation (claude-1 checked specifically). | C-3 |
| 4 | (scope note added) Live-read note said liveness pays ≤100 ms "only when refusing an alive-but-empty-cmdline pid" — but framed the choke point as recording-side plus a symmetric afterthought | The shared probe also changes **`Attributed`'s live read** (`procctl.go:152`): an empty live read was an **instant** refusal ("cannot read live command") pre-poll; now it waits ≤100 ms first, then compares or refuses at the same facet — so TUI liveness renders (map row 3) can spend the bound per alive-but-empty pid instead of refusing at once. Exposure is narrow (parley-owned zombies are transient — `cmd.Wait`/init reaps; persistent cost needs a persistent argv-less process, already un-attributable pre-fix). kimi-1 judges the placement correct: the execve window is symmetric (the live read can hit it microseconds after capture), and fixing only the recording side would leave the mirrored flake. No facet, order, or strictness change — timing-only scope. | K-2 |
| 5 | (not previously recorded) | **Pre-existing re-exec window**: a process that re-execs can publish its **pre-exec** argv inside the window; first-non-empty records it and `Attributed` then refuses fail-closed ("command mismatch (pid was reused)"). kimi-1 observed it 1/60 re-exec spawns (`sh -c 'exec sleep 30'`). Identical window pre-poll (the single read behaved the same on non-empty data) — recorded so a hosted sighting is not misread as a poll regression. | K-3 |

K-1 is C-1 (same finding, both reviewers). claude-1's C-4 (behaviour table
omits the live-probe row) and C-5 (pre-existing `commandMatches` comment/test
gap) are INFO/"no action requested" and outside this dispatch; C-4's substance
is covered factually by map row 4.

### Comment-only equivalence record

The corrections above are written into the two owned source files as **comment
text only**; the code at ef10cf3 is untouched. Proofs (host, go1.27.1
darwin/arm64):

- **git diff**: every added/removed line in both files is a `//` comment line
  (asserted programmatically; zero code lines touched).
- **Token streams**: `go/scanner` over both versions of each file, COMMENT
  tokens excluded, compared as (type, literal) pairs — **identical**
  (`procctl_linux.go`: 434 tokens; `procctl_cmdline_linux_test.go`: 1975).
- **Compiled**: `gofmt -l` clean on both; `GOOS=linux go build ./...` ok;
  `GOOS=linux go vet ./internal/procctl/` ok. Compiled archives/test binaries
  are *not* byte-identical — comment lines shift subsequent source line
  numbers, which DWARF/runtime line tables embed — so instead:
  `go tool objdump` of `pollCmdline` from the before/after test binaries
  (`-trimpath -buildvcs=false`), line annotations and addresses normalized:
  **65 lines of disassembly, identical**. Test re-runs are unchanged by
  construction (same instructions; only failure-message line numbers differ).

New file sha256 after the comment corrections (the Provenance table above
still records ef10cf3, which is immutable):

| File | sha256 at ef10cf3 | sha256 after comment corrections |
|---|---|---|
| `internal/procctl/procctl_linux.go` | `a3f8eb16…040547` | `81c77dd4720563bb7c6025c328fff315ed523771de6f3a4ca491823867053dfa` |
| `internal/procctl/procctl_cmdline_linux_test.go` | `5b3cca33…935b26` | `caa7dc7929f70a3941f60131ca9b8fc3a2ad29409f29e349844761f3a67bcded` |

**Line-number drift for future readers** (the reviews and the ef10cf3 commit
message cite old line numbers): in the test file, old lines 11–111 shift +4,
old 112–286 shift +6, old 287–301 shift +7 — so the cited spawn-race failure
sites `:170`/`:198` become `:176`/`:204`. In `procctl_linux.go`, old 68–87
shift +7 and old 88+ shifts +19 (`pollCmdline` moves 89 → 108).

### Not done

No code, test-logic, product-logic, or version change; no push, tag, release,
install, or CI rerun; the reviews, organizer records, closed artifacts, ledger,
inbox and runs/ are read/referenced only, never edited (reviews remain
untracked for the organizer). Safety facets, refusal strings, and fail-closed
behavior are exactly as at ef10cf3 — the hosted leg and its results belong to
the independent reviewer.
