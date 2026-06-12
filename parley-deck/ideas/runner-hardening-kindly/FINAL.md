---
idea: runner-hardening-kindly
status: final
drafted-by: claude
date: 2026-06-12
---

## Summary

Adopt six hardening mechanisms from the MIT-licensed "kindly" skill into
parley-deck-cli's runner and CLI: three-layer agent supervision (first-output
watchdog, stall guard, persisted heartbeats), data-driven failure
classification with operator recovery hints, artifact-beats-exit terminal
semantics, an environment/git/docs hardening batch, snapshot-checkout isolation
for Phase 6 reviews, and a `parley consult` advisory command. Ships as release
1.24.0 together with the sibling protocol idea
`meta-protocol-change-review-gate-honesty`. Consensus: all four participants
✅ ACCEPT.

## Final plan / specification

The authoritative decision record is consensus.md D1-D12; the build contract:

**Supervision (D1-D4).** `internal/runner/supervision.go` with
SupervisionConfig + in-process activityTracker; counting writers around
cmd.Stdout/Stderr for exec mode, acpRunnerHandler activity for ACP;
`waitSupervised` replaces the wait select in execAgentProcess. Watchdog events
(`agent.no_first_output`, `agent.stalled`) are appended BEFORE the kill tree
fires; kill remains procctl.KillGroup + drained Wait. Retry once only for
no_first_output with attempt_id threading and invalid-attempt move-aside.
Knobs: first_event_timeout_ms 120000 / stall_timeout_ms min(1800000, timeout)
/ heartbeat_ms 60000 on agents.Spec + TOML + run.created runtime payload.
`agent.heartbeat` persisted but excluded from narrator, snapshot triggers, and
review/consensus context; shown only in live status surfaces.

**Classification (D5-D6).** `internal/runner/failclass.go` seeded with agy's
12-class taxonomy and exact recovery-hint strings; bounded tails; one terminal
payload builder for exec/ACP/failEarly/fixup; agent.failed gains
failure_class, recovery_hint, exit_code, signal, stderr_tail_bytes;
runstate/AgentState + SummarizeEvent + TUI narrator/header + parley status
surface them.

**Artifact-wins (D7).** `finalizeAgentResult` decision table: a VALIDATED
artifact with an ordinary nonzero exit is agent.finished + agent_exit
(Result.Success() helper; driverImplOps and pipeline_cmd adapt); timeouts,
watchdog finals, and user kills always fail; ACP validation switches to
validateArtifactForPhase; fix-up gains ValidateFixupArtifact before
artifact-wins applies.

**Hardening batch (D8).** cleanParticipantEnv sheds Claude host markers for
claude participants; GIT_OPTIONAL_LOCKS=0 on gitTreeClean git probes;
`docs/agent-cli-mechanics.md` documents verified per-CLI mechanics (codex -o
documented, not adopted).

**Review snapshots (D9).** Phase 6 reviewers run in a disposable shared-clone
checkout on LOCAL tmp: clean tree → detached HEAD; dirty tree → kindly
temp-index snapshot commit; staged/worktree divergence or creation failure →
loud live-tree fallback. Resolved sha in review.snapshot_created + the review
prompt; artifact written inside the snapshot and moved back via
copy+fsync+rename; .pid markers with stale sweep and step-aside; events for
created/fallback/move-failure.

**Consult (D10).** `parley consult [--dir] [--timeout] <agent> "<question>"`:
facilitator captures stdout into
`parley-deck/consults/<YYYYMMDDTHHMMSSZ>-<agent>-<slug>.md` with the canonical
frontmatter (quorum: false), progress on stderr, fsutil.AppendLine-backed
index.jsonl with a cross-process claim, and `parley consults list`.

**Slices (D11):** result-decision+P3 → supervision core → classifiers →
hardening batch → snapshots → consult. Tests per consensus D12.

## Implementer

claude (FINAL drafter; default per protocol Phase 5). Reviewers: codex, agy,
hermes.

## Deferred follow-ups

- Counting-writer micro-details and consults claim shape decided in
  implementation (reviewed in Phase 6).
- Artifact-size watcher for buffered CLIs; codex {final_output} placeholder;
  driver ReadStrictGate enforcement (flag semantics in the sibling idea).
