---
idea: runner-hardening-kindly
agent: claude
status: complete
date: 2026-06-12
---

## Summary of work

All six points implemented per consensus D1-D12, shipping as 1.24.0 together
with the sibling protocol idea.

**D1-D4 Supervision.** `internal/runner/supervision.go`: SupervisionConfig +
atomic activityTracker + countingWriter (wraps cmd.Stdout/Stderr in
execAgentProcess — zero healthy-path FS I/O) + `waitSupervised` (1s tick;
first-output guard, stall guard, heartbeats; watchdog hook runs BEFORE kill;
always drains waitErr). Exec path: runAgent gained an attempt loop
(`runExecAttempt`) — retry once only for no_first_output, attempt_id threads
through events and the procctl marker (runID:agentID:attemptID), an invalid
attempt-1 artifact moves to `<artifact>.attempt-1.invalid`. ACP path: the
initialize→session→prompt sequence moved into a goroutine supervised by the
same waitSupervised; handler marks activity per session update;
`agent.started` never satisfies first-output. Knobs `first_event_timeout_ms` /
`stall_timeout_ms` / `heartbeat_ms` on agents.Spec + TOML overrides (pointer-
typed; explicit 0 ⇒ -1 disabled) + run.created runtime payload.
`agent.heartbeat` payload per D4; excluded from runstate.Recent, the narrator
allowlist, and snapshot triggers.

**D5-D6 Classification.** `internal/runner/failclass.go`: ordered regex table
seeded with agy's taxonomy + exact hints; bounded 4 KiB tails;
`terminalFailureClass` gives watchdog/timeout/kill causes priority.
`agent.failed` gains failure_class, recovery_hint, exit_code, signal,
stderr_tail_bytes; one decision/payload path per mode (finalizeExecResult /
finishACP / RunFixup). runstate.AgentState += FailureClass/RecoveryHint/
AgentExit; SummarizeEvent renders classed failures, watchdog events, and
heartbeats; TUI narrator (`friendlyEventText`) + agent status header +
trigger/narrator allowlists updated.

**D7 Artifact-wins.** Exec: validated artifact + ordinary `*exec.ExitError` →
agent.finished with agent_exit/agent_exit_kind=exec, ExitError cleared;
watchdog finals, hard timeout, and user kills always fail; failures always
carry a non-empty reason. ACP: validation switched from hard-coded
ValidateRoundOneArtifact to `validateArtifactForPhase` (bug fix); a prompt
error AFTER session open with a valid artifact → agent_exit_kind=acp_error.
Fix-up: new `ValidateFixupArtifact` (frontmatter idea match + review-ready
status + "## Fix-up cycle" section) gates success; exit code alone no longer
decides. New `Result.Success()`; driverImplOps (3 sites) and pipeline_cmd
(3 sites) switched to it; RunRoundOne's ok-count uses it.

**D8 Hardening batch.** `cleanParticipantEnv` sheds CLAUDECODE/CLAUDE_CODE_*/
AI_AGENT* when spawning the claude CLI (exec path; ACP env comes from
acp.MergedEnv — see deviations). gitTreeClean probes run with
GIT_OPTIONAL_LOCKS=0 (asserted by a PATH-shim test). New
`docs/agent-cli-mechanics.md` (verified codex/claude/agy/hermes mechanics,
incl. the hermes silent-death mode and the codex `-o` documented-not-adopted
decision), linked from docs/agent-runtime-configuration.md.

**D9 Review snapshots.** `internal/runner/reviewsnapshot.go`:
CreateReviewSnapshot (shared clone --no-checkout on LOCAL tmp; clean tree →
detached HEAD; dirty tree → kindly temp-index snapshot commit; staged/worktree
divergence and non-git roots → snapshotUnavailable), pid+boot-id markers with
stale sweep and step-aside, MoveArtifactBack (copy+fsync+rename within the
target dir), Cleanup. Wired into runAgent for Phase==review via a
publishArtifact hook so the terminal event reports the LIVE path; events
review.snapshot_created/_fallback/_artifact_move_failed.

**D10 Consult.** `internal/runner/consult.go` (BuildConsultPrompt +
RunConsult under supervision/classification; consult-flavor artifact-wins for
a non-empty answer + ordinary nonzero exit) and `internal/app/consult.go`
(`parley consult [--dir] [--timeout] AGENT [QUESTION|stdin]` — stderr
progress, artifact `parley-deck/consults/<UTC>-<agent>-<slug>.md` with the
canonical frontmatter, written even on failure; index.jsonl via new
`fsutil.AppendLine` (O_APPEND + mkdir-claim, stuck-claim degradation);
`parley consults list` with a bare-filename FILE column). Usage updated.

## Implementation plan / checklist

Slices landed in consensus D11 order: (1) result decision + artifact-wins,
(2) supervision core, (3) classifiers + consumers, (4) env/git/docs,
(5) snapshots, (6) consult. Version bumped to 1.24.0 + CHANGELOG entry.

## Tests (D12)

New: internal/runner/hardening_test.go (artifact-wins end-to-end with a real
subprocess; no-first-output watchdog retry matrix incl. the event-before-kill
ordering assertion; stall guard; heartbeat payload; supervisionForAgent
defaults/disable/clamp; classifier table; cleanParticipantEnv;
ValidateFixupArtifact; review-snapshot lifecycle on a real git repo incl.
dirty-tree snapshot commit, untracked capture, move-back, cleanup, non-git
fallback), internal/driver/gitprobe_test.go (GIT_OPTIONAL_LOCKS=0 PATH-shim
assert), fsutil AppendLine (+stuck claim), config supervision-knob tri-state,
app consultSlug. Full suite green; `-race` green on runner+fsutil.

## Deviations from FINAL.md

- ACP marker shedding: the ACP path builds env via acp.MergedEnv before
  spawning; claude-as-ACP-participant is not in the active roster, so the
  exec-path shedding covers the consensus case. Extending MergedEnv is a
  follow-up if a claude ACP participant returns.
- `signal` on agent.failed is parsed from the wait error text ("signal: X")
  instead of platform-specific WaitStatus plumbing — same information, no
  per-OS files.
- Consult artifact frontmatter adds `agent_exit` (the consult-flavor
  artifact-wins exit) alongside the agreed fields.

## Fix-up cycle 1 (review/round-01 → review/consensus.md)

All seven agreed fixes applied:

1. ACP contract: initialize/session-open/prompt-complete now mark activity;
   attempt_id threads through the ACP procctl marker, agent.started, heartbeat,
   watchdog, and terminal payloads; ACP attempts share the exec
   retry-once-for-no_first_output loop (lifted into runAgent's ACP branch).
2. Snapshot retention: cleanup is conditional — a move-back failure flips
   keepForRecovery, Abandon() drops only the marker (the stale sweep skips
   marker-less dirs) and the snapshot copy survives; terminal events now report
   the LIVE canonical path even on publish failure (publishArtifact always
   returns it).
3. moveAsideInvalidArtifact: unique suffix when .attempt-1.invalid exists
   (never overwrite an earlier recovery file); rename failure removes the
   invalid artifact. Test: TestMoveAsideInvalidArtifact.
4. RunFixup runs through the hardened exec path (process group + marker
   runID:agentID:fixup + cleanParticipantEnv + counting writers +
   waitSupervised + watchdog/heartbeat events, phase "fixup"); no retry by
   design (code-mutating phase).
5. failEarly emits the classified payload (failure_class + recovery_hint).
6. Consult provenance: frontmatter + ledger gain session_id (written even when
   empty — one-shot CLIs expose none) and timeout_ms now records the EFFECTIVE
   timeout (ConsultResult.EffectiveTimeout).
7. TestClassifyFailure locks the exact class/hint contract for all 9 provider
   classes + 3 watchdog hints (the dismissed verbatim-strings finding's
   testable half).

## Fix-up cycle 2 (review/round-02 → review/consensus.md cycle 2)

Both agreed remainders applied:

1. `finishACP` now mirrors `finalizeExecResult`: whenever publishArtifact
   returns a non-empty live path, `result.OutputPath` carries it — the ACP
   terminal event reports the LIVE canonical path even when the snapshot
   move-back failed (closes codex's fix-2 remainder and agy's NIT).
2. `TestMoveAsideInvalidArtifact` gains the rename-failure case: a source
   basename near NAME_MAX makes the `.attempt-1.invalid` destination exceed
   it, the rename fails deterministically (ENAMETOOLONG), and the test asserts
   the invalid artifact is removed from the canonical path with no recovery
   file left behind.

## Notes for reviewers

Known sandbox artifact: `TestDurableKillEndToEndRealProcess` fails under
codex's seatbelt sandbox (sysctl kern.boottime restricted); it passes in a
normal shell. The new watchdog tests take ~7s combined (1s supervisor tick ×
real subprocesses).
