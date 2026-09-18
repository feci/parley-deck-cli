---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
type: partial-source-review
reviewed-commit: 122f4d2
scope: internal/runner/launch.go (RunInteractive), internal/runner/launch_test.go,
  internal/runner/telemetry.go, internal/telemetry/record.go,
  internal/procctl/terminal_unix.go, internal/procctl/terminal_windows.go,
  internal/app/consensus_request_signoffs.go, docs/agent-runtime-configuration.md
signoff: none
---

# Partial source review — terminal-launch slice

## 1. Scope and what this is not

This is a **partial source review of one slice**. It is **not** Phase-6 full-scope
acceptance, not a review-consensus signoff, and not an AC closure for AC-T1/AC-T2/AC-P1.
No implementation edit, git mutation, config change or test execution was performed by
this reviewer.

Reviewed, by direct file reading: `internal/runner/launch.go`,
`internal/runner/launch_test.go`, `internal/runner/telemetry.go`,
`internal/runner/protocol_context.go`, `internal/runner/handoff.go`,
`internal/telemetry/record.go`, `internal/telemetry/usage.go`,
`internal/procctl/terminal_unix.go`, `internal/procctl/terminal_windows.go`,
`internal/procctl/procctl.go`, `internal/procctl/procctl_unix.go`,
`internal/app/consensus_request_signoffs.go`, `internal/agents/discover.go` (defaults
only), `docs/agent-runtime-configuration.md`, `.gitignore`, plus this idea's `FINAL.md`
and `IMPLEMENTATION.md`.

Explicitly **not** reviewed: `internal/protocolpacket/**`, `internal/fsutil/**`,
`internal/driver/**`, `internal/evidence/**`, the ACP/consult/steer launch paths, the
skill repository, and the evaluation/pilot work. `internal/protocolpacket` and the
source-instruction slice are **this reviewer's own allocation** and are excluded from
any verdict here (§15.1 — no self-verdicts).

## 2. Provenance (§15.2)

- **PRIMARY** — every claim below about *what the code says* is from the files named,
  read directly at the stated line locators in the integration worktree
  `worktrees/evidence-first-integration`. The two *negative* claims ("no other caller of
  `RunInteractive`", "`spawn-tty` is read in exactly two places") come from repo-wide
  regex searches for `RunInteractive\(|runInteractiveTTY\(` and
  `InteractiveInvokeSpawnTTY|InteractiveInvokeOrDefault` over all `*.go`; a caller reached
  by reflection or code generation would escape that search.
- **Environment-reported** — the mapping of that working tree to commit `122f4d2`
  ("measure terminal processes and preserve foreground ownership") comes from the
  session's git-status snapshot, not from a git command this reviewer executed. If the
  tree has moved since, re-anchor the locators.
- **Facilitator-reported, not verified here** — runner/app package tests PASS, race and
  vet PASS before the final prompt-delivery guard, the compiled runner test binary under
  Python `pty.fork` PASS (child TTY detection, child terminal read, restored parent
  terminal read), Windows cross-compiles with runtime unvalidated, full Go suite running.
  These remain the facilitator's results. This reviewer executed nothing.
- **RECALL (UNVERIFIED)** — three background facts are memory-only here because the Go
  toolchain sources are outside the permitted file scope: (a) `os/exec` passes an
  `*os.File` std stream straight to the child instead of creating a pipe; (b) the Go
  runtime blocks all signals across `fork`, which is what masks `SIGTTOU` around the
  child's `TIOCSPGRP`; (c) Linux `MAX_ARG_STRLEN` caps a *single* argv element at 128 KiB.
  Each is tagged where used. (a) and (b) are corroborated — not proven — by the
  facilitator-reported PTY run.

## 3. Refutation attempts (LE-1)

Each item states what was attempted against the claimed rationale and the result.

1. **"Missing prompt delivery refuses before spawn."** Tried to find a path that reaches
   `exec.Cmd` with no delivery. The guard at `internal/runner/launch.go:237-261` runs
   *before* `beginProtocolLaunch` (`:262`), so it refuses before the protocol is even
   rendered, let alone spawned. `placeholder` is `""` for mode `none` **and for any
   unrecognized mode**, and `deliversPrompt` then stays false — fail-closed in both
   cases. The refusal is a `*protocolContextError`, which `telemetry.go:136-139`
   classifies as `protocol_context_refused`; the fixture asserts exactly that
   (`launch_test.go:115`). **Could not break it.**
2. **"file mode uses `{prompt_path}`, arg mode uses `{prompt}`."** Tried substring
   confusion: in arg mode the guard looks for `{prompt}`, and `"{prompt_path}"` does not
   contain `"{prompt}"` (closing brace differs), so a file-mode template cannot satisfy
   an arg-mode contract. `ExpandInteractiveArgs` (`handoff.go:130-139`) expands only
   `{root}`, `{prompt_path}`, `{target_path}` — it never touches `{prompt}`, so the later
   `ReplaceAll` at `launch.go:289-293` cannot be pre-empted into a silent no-op.
   **Could not break it.** (Size/exposure of that substitution is MAJOR-3 below.)
3. **"Default print-only handoffs are unchanged."** Checked every built-in spec: all set
   `InteractivePromptMode: InteractivePromptNone` and `InteractiveInvoke:
   InteractiveInvokePrintOnly` (`internal/agents/discover.go:212, 235, 260, 284, 312,
   343, 375, 408`; `acp_specs.go:118-119`), and the zero-value defaults resolve the same
   way (`discover.go:646-660`). `RunInteractive` is reached only through
   `runInteractiveTTY`, gated on `invoke == spawn-tty`
   (`consensus_request_signoffs.go:481-485`). **Could not break it.**
4. **"The fresh rendered file gets a distinct process invocation; the previous handoff is
   not a process."** `prepareProtocolPrompt` re-resolves the live source per call
   (`protocol_context.go:74`) and `RunInteractive` writes its own prompt into the
   per-invocation directory (`launch.go:280-283`), not the handoff path. The fixture
   over-writes the handoff prompt with `"tampered handoff"` and then asserts the child saw
   the fresh bytes and not the tampered ones (`launch_test.go:42-52`), and that exactly two
   terminal records exist with the handoff one at `unobserved-handoff` / `StartedAt == nil`
   (`:56-63`). **Could not break it.**
5. **"Real file descriptors reach the child."** `cmd.Stdin, cmd.Stdout, cmd.Stderr =
   stdin, stdout, stderr` with `*os.File` values (`launch.go:297`); no wrapper, no
   `MultiWriter` (contrast the captured path, `launch.go:105-106`). That this avoids a
   pipe is **RECALL(a)**; the operative evidence is the PTY fixture asserting
   `test -t 0 && test -t 1 && test -t 2` plus a real `read` in the child
   (`launch_test.go:146`), which is facilitator-reported as passing. **Not broken, but
   the strongest evidence here is not mine.**
6. **"Streams are explicitly not-observed-terminal."** `evidence.directTerminal = true`
   (`launch.go:266`) forces `observation.StreamCoverage = "not-observed-terminal"`
   (`telemetry.go:108-111`), which is in the accepted enum (`record.go:290-294`). With no
   collector writes, `Collector.Result()` degrades usage to
   `source/cost_basis/coverage = "unavailable"` and leaves cost/model nil
   (`usage.go:277-285`, `record.go:267-271`). **Partially broken — see MAJOR-1**: the two
   byte counters are still emitted as `0`.
7. **"Foreground restoration uses a constant shell no-op … without altering parent-wide
   signal handlers."** The shim is a fixed `/bin/sh -c "exit 0"` with `Foreground: true,
   Pgid: group` (`terminal_unix.go:34-39`) — no task text, no model, so it is correctly
   *not* agent telemetry, and the parent's `sigaction`/mask is untouched. The masking it
   relies on is **RECALL(b)**. Tried to find a path where the child is still alive when
   restore runs: all three start paths wait first — failed `Start` returns with no child
   (`launch.go:315-317`), a failed `evidence.started` does `Cancel` + `Wait`
   (`:318-322`), and the normal path returns `cmd.Wait()` (`:323`). Defer order is
   restore → finish → cancel, i.e. the terminal is restored before telemetry is sealed.
   **Could not break it.**
8. **"Timeouts terminate the owned child group."** `cmd.Cancel` targets
   `Spawned{PID: pid, PGID: pid}` (`launch.go:312-314`). That PGID is right under *both*
   branches: `Foreground: true` with `Pgid: 0` makes the child its own group leader, and
   the ENOTTY fallback's `Setsid` does likewise (`terminal_unix.go:22`,
   `procctl_unix.go:14-19`). `KillGroup` signals `-pgid` with TERM→1.5 s→KILL and refuses
   to signal parley's own group (`procctl_unix.go:25-52`). **Not broken on inspection —
   but see MAJOR-2: the only cleanup fixture never exercises the foreground branch.**
9. **"Lifecycle accounting is honest."** Tried to make a refused or failed attempt
   disappear: the delivery refusal, an authority refusal, a failed start, a non-zero exit,
   a timeout, and a `started.json` write failure each produce exactly one terminal record
   with the expected failure class (`launch_test.go:73-121`). A terminal-telemetry write
   failure returns an integrity error that replaces the process result
   (`telemetry.go:143-146`). **Could not break it.**
10. **"Every supported launch path is covered" (D2, AC-T1) — for this knob.** Searched the
    whole tree for `RunInteractive` and `InteractiveInvokeSpawnTTY` callers. Only
    `consensus request-signoffs` honours `spawn-tty`; no other surface does. **Broken —
    see MINOR-6.**
11. **"A signoff refusal preserves the underlying launch error."** Traced
    `runInteractiveTTY` failure through `consensus_request_signoffs.go:186-197`: the
    artifact is validated *first*, and `runErr` is wrapped into the reported error rather
    than being replaced by a bare "did not append a signoff". **Could not break it.**

## 4. Findings

### CRITICAL

None found in this slice.

### [MAJOR] MAJOR-1 — unobserved terminal streams are serialized as measured zeros

`Observation.StdoutBytes` / `StderrBytes` are plain `int64` with no `omitempty`
(`internal/telemetry/record.go:60-66`), so a terminal invocation whose streams are by
construction *unobserved* still writes `"stdout_bytes": 0, "stderr_bytes": 0`. D2 in
`FINAL.md` says "Unknown values are null", and the sibling field in the same struct
(`FirstActivityMS *int64`) is nullable for exactly this reason. The launch plan for this
slice also says "do not … invent unobserved stream usage"; a literal `0` is an invented
observation of zero output.

`stream_coverage: not-observed-terminal` does disambiguate for a reader who looks at it —
but it carries `omitempty`, so the *handoff* record (`status: unobserved-handoff`, where
`directTerminal` is never set) emits the same two zeros with **no** coverage field at all.
Any consumer that sums observed bytes across invocations — the >= 20-attempt
reconciliation and the offline HTML are the obvious ones — will silently fold unobserved
launches in as zero-output launches.

The existing fixture checks `FirstActivityMS`, `CostUSD` and `ReportedModel` are nil
(`launch_test.go:67`) but not the byte counters.

Suggested fix: make both counters `*int64` (nil when not observed), or emit them only
when `stream_coverage == "captured"`, and extend `launch_test.go:67` to assert nil.

### [MAJOR] MAJOR-2 — no fixture exercises cleanup or restoration on the real-terminal branch

`TestInteractiveTimeoutKillsDescendants` (`launch_test.go:123-135`) — the only cleanup
fixture — passes a regular temp file as the "terminal" (`interactiveFixture`,
`launch_test.go:18-31`, `os.CreateTemp`). `AttachTerminal` therefore takes the
`ENOTTY` → `SetNewProcessGroup` early return (`terminal_unix.go:20-24`) and returns a
**no-op restore**. So that test proves descendant cleanup for the `Setsid` fallback and
proves nothing about the code this commit actually adds: `Foreground: true` + `Ctty`
(`terminal_unix.go:28`) and the `/bin/sh` restore shim (`:34-39`). The same applies to
`TestInteractiveProcessFailureEvidence`'s `timeout` case.

The one real-TTY fixture (`TestInteractiveRealTerminal`, `launch_test.go:140-162`,
skipped unless `PARLEY_TTY_TEST=1`) covers only the success path: child reads, then parent
reads. The riskiest new state transition — **killing a child that currently owns the
terminal foreground, then re-taking the foreground group from a background parent** — has
no fixture at all, under the PTY harness or otherwise. The failure mode if it is wrong is
severe and user-visible: a shell left with a dead foreground process group.

Suggested fix: add a PTY-harness case that sets a short `InteractiveTimeoutMS` against a
`sleep`-style child and asserts (a) the descendant did not survive and (b) the parent can
still read the terminal afterwards. This is a coverage gap, not a demonstrated defect —
inspection did not find an error in the logic.

### [MAJOR] MAJOR-3 — arg mode puts the entire rendered protocol into one argv element

`args[i] = strings.ReplaceAll(args[i], "{prompt}", prepared)` (`launch.go:289-293`) where
`prepared` is the full attestation + verbatim protocol + task (`protocol_context.go:111`).
Two consequences, no guard for either:

- **Size.** The live protocol for this launch is 108,400 bytes (`source_bytes` in this
  session's own shadow-packet audit). **RECALL(c)**: Linux caps a single argv element at
  `MAX_ARG_STRLEN` = 128 KiB, which would leave roughly 22 KB of headroom today and fail
  with a bare `E2BIG` after any protocol growth. That number is unverified here and should
  be confirmed before being relied on — but the *structural* point needs no constant: an
  unbounded document is interpolated into one argument with no size check and no
  actionable diagnostic. It does fail closed (recorded as `start_failure`), just opaquely.
- **Exposure.** argv is visible in the process table. How widely is platform-dependent and
  is **RECALL** here (Linux `/proc/<pid>/cmdline` is world-readable under the default
  `hidepid`; macOS restricts other users' argv) — but on any platform this is a strictly
  wider exposure than the file path, which is 0600
  inside a 0700 per-invocation directory (`writeHandoffPrompt` →
  `os.CreateTemp`, `handoff.go:144`; `telemetry.Begin` → `os.Mkdir(dir, 0o700)`,
  `record.go:174`).

`docs/agent-runtime-configuration.md:120-129` presents `file` and `arg` as equally valid
with no size or visibility caveat. Suggested fix: prefer/recommend file mode explicitly,
refuse above a bounded prompt size with a named failure class, and document the argv
exposure.

### [MINOR] MINOR-1 — `RunInteractive` has no launch-mode guard

`RunMeasured` refuses a non-headless agent (`launch.go:187-189`) and `WriteHandoffPacket`
deliberately rewrites a headless agent's mode to `manual` so the record cannot claim the
wrong mode (`handoff.go:42-44`). `RunInteractive` does neither: passed a headless (or
`acp`/`manual`) `Discovery` it will spawn it and record
`launch_mode: headless` (`telemetry.go:68` reads `agent.LaunchMode` verbatim). Today the
only caller checks `invoke == spawn-tty` first, so this is a latent mislabel, not a live
one. Suggested fix: mirror the `RunMeasured` guard.

### [MINOR] MINOR-2 — the terminal-restoration error is discarded

`launch.go:302-306`: `if err := restore(); err != nil { returnedErr =
errors.Join(returnedErr, errors.New("cannot restore terminal foreground process group")) }`
— `err` is bound and then dropped. The operator loses the only diagnosis available for a
class of failure that is hard to reproduce (EPERM from `setpgid` when the previous group
is gone, a missing `/bin/sh`, the 5 s shim deadline). Wrapping is safe here because the
shim's argv is a compile-time constant and its stdout/stderr are never captured, so no
task content can ride along. Suggested fix: `fmt.Errorf("cannot restore terminal
foreground process group: %w", err)`.

### [MINOR] MINOR-3 — the delivery contract is not validated at selection time

`validateLaunchModes` (`consensus_request_signoffs.go:340-361`) validates
`interactive_prompt_mode` and `interactive_invoke` shapes but not the
`spawn-tty` ⇒ matching-placeholder contract. So a misconfigured agent passes selection,
`writeSignoffHandoff` publishes a handoff packet and an `agent.handoff.started` event
(`:460-479`) that announces a session, and only then does `RunInteractive` refuse. Two
residues: an orphaned handoff, and a run-store trail with `agent.handoff.started` and no
corresponding failure or completion event (only the private invocation record carries the
refusal). Suggested fix: check the placeholder in `validateLaunchModes` for
`invoke == spawn-tty`, and/or append a handoff-failed event on refusal.

### [MINOR] MINOR-4 — the interactive timeout budget is applied twice

`RunInteractive` bounds the spawned session with
`InteractiveTimeoutMSOrDefault` (`launch.go:234`), and after it returns the caller starts
a *fresh* poll deadline of the same duration
(`consensus_request_signoffs.go:487-490`, `:589-591`). Worst case is 2× the configured
value before the request gives up, which is not what a reader of
`interactive_timeout_ms` would expect. Suggested fix: derive the poll deadline from the
remaining budget, or document the two-stage semantics.

### [MINOR] MINOR-5 — handoff prompts (full protocol + task text) land in a tracked path

`.gitignore:23` correctly ignores `/.parley-runtime/`, which is where the spawned
process's prompt lives (`launch.go:280`). But `WriteHandoffPacket` writes the same
rendered bytes to `parley-deck/runs/<run-id>/agents/<id>/handoff-prompt-<invocation>.md`
(`handoff.go:66-79`), and `parley-deck/runs/` is **not** ignored. Each handoff therefore
adds a ~108 KB copy of the protocol plus the verbatim task text to a committable path.
This predates `122f4d2` — flagged here because this commit makes the two-prompt split
load-bearing, so it is the natural moment to decide whether the run-store copy should be
the full render or a reference.

### [MINOR] MINOR-6 — `spawn-tty` is honoured on exactly one launch surface, silently

`InteractiveInvokeSpawnTTY` is read in exactly two places in the whole tree:
`validateLaunchModes` (`consensus_request_signoffs.go:355-357`) and the dispatch at
`consensus_request_signoffs.go:481-485`. `RunInteractive` has no other production caller
(only `internal/runner/launch_test.go`). So an agent configured
`launch_mode = "interactive"`, `interactive_invoke = "spawn-tty"` spawns a measured
terminal process for `consensus request-signoffs` and, for every other surface — rounds,
cross-review, implementation/review/fixup, steer, consult/goal-check — silently degrades
to a printed handoff with no process and no diagnostic. The handoff instructions make this
worse by printing `Invoke: spawn-tty` (`handoff.go:100`) on a packet that nothing will
spawn.

`IMPLEMENTATION.md` already records that launch coverage is incomplete, so this is not an
undisclosed gap in general — but the specific shape (a configuration knob that is a silent
no-op on most surfaces, while the handoff text asserts it is in effect) is worth either a
warning at validation time or a corrected `Invoke:` line until the other surfaces are
wired. It also bounds what this slice can contribute to AC-T1.

### [NIT]

- `launch.go:244-248`: the `placeholder != ""` test is loop-invariant and there is no
  `break`; an early return when `placeholder == ""` would read better.
- `terminal_unix.go:28` replaces any pre-existing `cmd.SysProcAttr` wholesale; the doc
  comment does not say so. Harmless today (no caller sets it first), fragile later.
- `launch.go:294-297` sets no `WaitDelay`, unlike `commandForLaunch` (`launch.go:63`).
  `KillGroup` escalates to SIGKILL so `Wait` should return, but the bound is implicit.
- Both refusal paths finish their evidence before `directTerminal` is set
  (`launch.go:249-261`, `:262-266`), so those records omit `stream_coverage` rather than
  labelling it `not-observed-terminal`. No stream existed, so nothing is misstated — but
  the labelling is inconsistent with the successful case.
- `terminal_unix.go:36` hardcodes `/bin/sh`.
- `docs/agent-runtime-configuration.md:120-129` is accurate on delivery, freshness,
  descriptors, unobserved streams and group restoration, but does not mention that
  `arg` mode is the one that carries the whole protocol on the command line (MAJOR-3).

## 5. Open questions for the owner

1. **`Ctty` semantics.** `terminal_unix.go:28` and `:38` pass `int(input.Fd())` — a
   *parent* descriptor number. Whether `Foreground` interprets `Ctty` in the parent's fd
   space (at fork, before the descriptor shuffle) or the child's is **RECALL** for this
   reviewer and could not be checked, since the Go sources are outside the permitted file
   scope. It is moot for every current caller, because they all pass `os.Stdin`
   (`consensus_request_signoffs.go:558`) where `Fd() == 0` and both readings coincide —
   which is also why the passing PTY harness does not discriminate between them. A future
   caller passing an explicitly opened `/dev/tty` at a higher fd may or may not work.
   Either document "input must be the child's stdin" on `AttachTerminal`, or confirm the
   semantics.
2. **Windows.** `terminal_windows.go:12-15` is a reasonable no-op (console inheritance +
   `KillGroup`), consistent with `IMPLEMENTATION.md`'s "cross-compiles; runtime
   unvalidated". No objection — just confirming that the untested status is recorded, not
   resolved.
3. **Non-zero exit after a valid append.** With spawn-tty, a user-driven session that
   exits non-zero after correctly appending a signoff still aborts the whole batch
   (`consensus_request_signoffs.go:194-197`), so later participants are never requested.
   That is pre-existing headless behaviour now reachable from a human terminal session,
   where a non-zero exit is a much more ordinary way to end. Worth an explicit decision
   rather than inheritance.

## 6. Standing obligations this review does not touch

Unchanged and still open, per `IMPLEMENTATION.md`: shared per-launch/per-action budgets,
the Kimi typed-evidence integration, the exact ratified packet experiment (AC-P2), the
twelve-task pilot (AC-X1/X2), the >= 20-attempt reconciliation (AC-T2), and the full
Phase-6/7/8 review and signoff cycle. **Nothing in this slice waives any of them**, and
nothing in this note is a signoff.
