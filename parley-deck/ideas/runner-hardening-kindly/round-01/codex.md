---
agent: codex
idea: runner-hardening-kindly
round: 1
date: 2026-06-12
---

## Summary

Ship this as a runner/result contract change, not as six local patches. The current code has three separate wait/result paths: exec rounds in `internal/runner/runner.go:385-454`, ACP rounds in `internal/runner/acp.go:106-185`, and fix-up in `internal/runner/phase58.go:66-112`; P1-P3 must converge those paths or the first release will have inconsistent failure semantics.

## P1 - supervision hook points

Current exec seam: `runAgent` skips pre-existing artifacts at `internal/runner/runner.go:326-339`, builds prompts at `internal/runner/runner.go:348-356`, records durable `agent.started` at `internal/runner/runner.go:366-384`, then blocks in `execAgentProcess`; the actual wait is the `cmd.Wait()` select at `internal/runner/runner.go:730-739`. Put the watchdog in that select, not above `runAgent`, so process-group ownership stays with `procctl.KillGroup` and steer attempts keep the same semantics.

Concrete exec design:

- Add `internal/runner/supervision.go` with `SupervisionConfig{FirstEventTimeout, StallTimeout, HeartbeatInterval}` and an `activityTracker` that records stdout/stderr byte counts and last activity time.
- Extend `execAgentProcess` to accept `SupervisionOptions{Store, AgentID, SegmentID, AttemptID, Phase, Launch, Timeout, Config}`. Wrap `cmd.Stdout` and `cmd.Stderr` at `internal/runner/runner.go:718-719` in counting writers that still write to the same log files.
- Replace the select at `internal/runner/runner.go:732-739` with `waitSupervised(ctx, waitErr, kill func(){procctl.KillGroup(sp)}, activity)`. On watchdog kill, drain `waitErr` exactly as today so the child is reaped.
- First-output means first stdout/stderr byte for exec mode. A future optional artifact-size watcher can be added for buffered CLIs, but do not stat protocol paths every tick in v1; the write wrappers give zero extra filesystem I/O on healthy paths.

Current ACP seam: `runACPAgent` spawns at `internal/runner/acp.go:57-66`, records durable `agent.started` at `internal/runner/acp.go:81-104`, then runs `Initialize`, `NewSession`, and `Prompt` inline at `internal/runner/acp.go:106-138`. ACP activity already arrives through `acpRunnerHandler.SessionUpdate` at `internal/runner/acp.go:208-266`.

Concrete ACP design:

- Give `acpRunnerHandler` an `activity *activityTracker`; mark activity for every ACP session update and for runner-emitted `agent.acp.initialized`, `agent.acp.session_opened`, and `agent.acp.prompt_completed`.
- Move the initialize/session/prompt sequence into a goroutine returning `finishErr`. The caller waits with the same `waitSupervised`; the kill function should cancel the ACP context and call `process.Kill()`/`process.Stop()` from `internal/acp/spawn.go:104-139`.
- First event for ACP is any ACP protocol activity after `agent.started`; `agent.started` itself must not satisfy the first-output guard.

Config surface:

- Add `first_event_timeout_ms`, `stall_timeout_ms`, and `heartbeat_ms` to `agents.Spec` near `TimeoutMS` in `internal/agents/discover.go:34-43`, plus TOML overrides in `internal/config/runtime.go:22-52` and `applyOverride` around `internal/config/runtime.go:266-277`.
- Defaults: first event 120000 ms; heartbeat 30000 ms; stall `min(1800000, agent timeout)` with `0` disabling a guard only when explicitly configured. Existing `timeout_ms` remains the hard outer budget (`internal/agents/discover.go:62`, `internal/runner/runner.go:784-796`).
- Include the effective knobs in `run.created` runtime rows at `internal/runcontrol/runcontrol.go:149-169` so reattached TUI/status views can explain the policy used.

Retry:

- Retry once only for `no_first_output`, not for `stalled`, timeout, validation failure, auth/billing/config errors, or user kill.
- Keep the skip check at `internal/runner/runner.go:326-339` before the attempt loop. `Overwrite=false` still protects an artifact that existed before this invocation. Retry is an internal second attempt against the same requested artifact path.
- Track whether the artifact path did not exist before attempt 1. If attempt 1 creates an invalid file, move it to `<artifact>.attempt-1.invalid` before retry; never delete a valid artifact, and never overwrite an artifact that pre-existed the runner call.
- Emit `agent.failed` for the killed attempt before `agent.started` for the retry. `KillAgentDurable` treats `agent.failed` as terminal at `internal/runner/durablekill.go:121-123`, so a reattached kill will target the latest started attempt, not the already-killed one.
- Add `attempt_id` to `agent.started`, typed watchdog events, and terminal events. Use marker `runID:agentID:attemptID` when calling `procctl.MarkerEnv` at `internal/runner/runner.go:704` so logs and process identity line up, even though attribution remains PID/start/pgid/command based (`internal/procctl/procctl.go:100-157`).

Events:

- `agent.heartbeat`: `{agent, segment_id, attempt_id, phase, launch, elapsed_ms, timeout_ms, stdout_bytes, stderr_bytes, last_activity_ms_ago}`. Do not count heartbeat itself as activity.
- `agent.no_first_output`: `{agent, segment_id, attempt_id, grace_ms, elapsed_ms, stdout_bytes, stderr_bytes, pid, pgid, action:"retrying"|"failed"}` followed by terminal `agent.failed`.
- `agent.stalled`: `{agent, segment_id, attempt_id, stall_ms, elapsed_ms, stdout_bytes, stderr_bytes, pid, pgid}` followed by terminal `agent.failed`.

## P2 - failure classification

Add a bounded classifier over stderr tail, stdout tail, and ACP error events. Keep it data-driven: ordered regex rules for `rate_limit`, `overloaded`, `auth`, `billing`, `invalid_request`, `model_not_found`, `context_window`, `sandbox`, `budget`, and `unknown`, with a max byte/line budget similar to the reference script's tailing behavior.

`agent.failed` should carry: `failure_class`, `recovery_hint`, `exit_code`, `signal`, `stderr_tail_bytes`, and existing `error`. Watchdog terminal failures set `failure_class` to `no_first_output`, `stalled`, or `timeout` before regex classification. `failEarly` at `internal/runner/runner.go:459-477`, exec terminal append at `internal/runner/runner.go:441-454`, ACP terminal append at `internal/runner/acp.go:170-184`, and fix-up append at `internal/runner/phase58.go:107-112` must all use the same payload builder.

Consumers to adapt:

- `internal/runstate/runstate.go`: extend `AgentState` around `internal/runstate/runstate.go:50-65` with `FailureClass`, `RecoveryHint`, and `AgentExit`; set them in `applyAgentEvent` at `internal/runstate/runstate.go:444-485`; summarize watchdog and hint text in `SummarizeEvent` at `internal/runstate/runstate.go:409-441`.
- `internal/tui/protosnap.go:75-92`: add watchdog event types to narrator/snapshot trigger sets. Heartbeats should not be woven into every transcript by default; surface them in status/protocol views or verbose mode.
- `internal/tui/protocolui.go:216-264`: render typed watchdog lines and a concise recovery hint. Update the `agent.finished` comment at `internal/tui/protocolui.go:343-346` because finished can now carry a nonzero `agent_exit` note.
- `internal/tui/live.go:815-860` and `internal/app/app.go:2234-2261`: show recovery hints next to failed agents, while keeping raw log paths visible.

## P3 - artifact-wins decision table

The current terminal decision is `agent.failed` when `ExitError != "" || !ArtifactOK` in exec mode (`internal/runner/runner.go:437-440`) and ACP (`internal/runner/acp.go:166-169`). Replace this with a helper such as `finalizeAgentResult(result, runErr, validateErr, supervisionErr)`.

Decision table:

| Artifact validation | Process/protocol outcome | Terminal event | Result fields |
| --- | --- | --- | --- |
| valid | nil | `agent.finished` | `ArtifactOK=true`, `ExitError=""` |
| valid | ordinary nonzero exit | `agent.finished` with `agent_exit` | `ArtifactOK=true`, `ExitError=""`, `AgentExit` set |
| valid | ACP prompt error after session started | `agent.finished` with `agent_exit.kind="acp_error"` | same, unless the error is init/session setup |
| valid | timeout, no-first final failure, stalled, or user kill | `agent.failed` | watchdog/user cancellation wins |
| invalid/missing | nil | `agent.failed` | error is validation/missing artifact |
| invalid/missing | nonzero exit | `agent.failed` | combine validation and exit details |
| invalid/missing | watchdog/user cancellation | `agent.failed` | watchdog/user error plus validation if available |

Implementation phases:

- Phase 1/2 rounds, Phase 5 implementation, Phase 6 review, and Phase 7 review-consensus all use `runAgent`, so the helper applies once. Also fix ACP validation to call `validateArtifactForPhase` instead of hard-coded `ValidateRoundOneArtifact` at `internal/runner/acp.go:158-163`; otherwise ACP review/implementation artifacts can never validate correctly.
- Phase 8 `RunFixup` is separate (`internal/runner/phase58.go:46-112`). Do not call it successful merely because the process exits zero. Add `ValidateFixupArtifact(ideaPath, cycle)` that verifies `IMPLEMENTATION.md` still validates, contains a new fix-up section, and has a review-ready status. Then apply artifact-wins only for ordinary nonzero exits after that validation passes.
- `driverImplOps` currently treats any `Result.ExitError` as failure for implementation/review-consensus/fixup at `internal/app/driver_impl.go:99-105`, `internal/app/driver_impl.go:178-187`, and `internal/app/driver_impl.go:218-224`. The runner must clear `ExitError` on artifact-wins, or these adapters must switch to a `Result.Success()` method.
- Pipeline auto has the same direct checks at `internal/app/pipeline_cmd.go:780-784`, `817-820`, and `845-848`; include it in tests.

## P4 - small hardening

- Nested Claude host markers: `buildAgentInvocation` currently starts from `os.Environ()` at `internal/runner/runner.go:700-704` or isolated env at `internal/runner/runner.go:764-768`. Add `cleanParticipantEnv(agent.ID, env)` before appending `procctl.MarkerEnv`; for `claude`, remove `CLAUDECODE`, `CLAUDE_CODE_SESSION_ID`, `CLAUDE_CODE_ENTRYPOINT`, `AI_AGENT`, and `AI_AGENT_*`. Keep Parley's own `PARLEY_*` markers.
- Git read probes: `gitTreeClean` shells out at `internal/driver/impl.go:283-291`. Wrap read-only git commands with `GIT_OPTIONAL_LOCKS=0`; tests should assert the env. The repo search shows this is the only Go git probe today.
- CLI mechanics docs: add `docs/agent-cli-mechanics.md` and link it from `docs/agent-runtime-configuration.md`. Document codex stdin/final-output behavior, Claude `-p` binding and tool allowlists, agy `--print` value-taking, and Hermes `-z`/oneshot mechanics.
- Codex `-o`: do not switch the built-in spec at `internal/agents/discover.go:100-120` until the invocation schema has a `{final_output}` placeholder. For v1.24.0, document it and keep participant artifact validation as the authority.

## P5 - snapshot Phase 6 reviews

Introduce `internal/snapshot` or `internal/runner/reviewsnapshot.go`; do not bury git clone logic inside `RunReviewRound`.

Mechanics:

- Only Phase 6 review rounds use snapshots. Phase 5 and Phase 8 need the live worktree; deliberation rounds are protocol-only.
- Resolve the reviewed commit from `IMPLEMENTATION.md head-commit`, falling back to live `HEAD` only with a loud `review.snapshot_fallback` event. Verify with `GIT_OPTIONAL_LOCKS=0 git -C <live> rev-parse --verify <sha>^{commit}`.
- Create per-reviewer clones under local tmp, not the virtio-fs workspace: `<tmp>/parley-review-snapshots/<repo-hash>/<idea>/<round>/<agent>`. Use `fsutil.MkdirAllResilient` for parent creation.
- Use `git clone --shared --no-checkout <live-root> <snap-root>` and `git -C <snap-root> checkout --detach <sha>`. The clone's alternates file points back to the origin object store on virtio-fs; this should be read-only for origin, but it depends on origin objects remaining available until review completes.
- Copy the live idea protocol context into the snapshot before prompting, or gather review context before root swap. The reviewer must read code from `<snap-root>` but write the review artifact inside the snapshot, then the runner validates and moves it back to the live `parley-deck/ideas/.../review/round-NN/<agent>.md` with temp+rename.
- Terminal events should report the live artifact path, not the soon-to-be-deleted snapshot path. If move-back fails, leave the snapshot and emit `review.snapshot_artifact_move_failed` with recovery path.
- Snapshot lifecycle uses a `.pid` marker with run id, agent id, pid, and boot id. On create, remove stale markers whose process is gone; if a live marker exists, step aside to a suffixed path. Do not signal from snapshot healing; durable kill remains the `agent.started`/`procctl` path.
- On snapshot failure, emit `review.snapshot_fallback` and run against the live root. This is allowed but should be visible in TUI/status because it weakens the reviewed-commit pin.

Correctness risks: alternates can break if the origin object is pruned; submodules and ignored build artifacts are absent unless explicitly handled; sandboxed agents may not be able to write outside the snapshot, hence move-back is required; cwd-scoped resume makes stable per-agent paths useful but dangerous under concurrency unless marker healing is strict.

## P8 - consult command

Wire `parley consult` in the top-level switch at `internal/app/app.go:50-85` and usage at `internal/app/app.go:100-120`. Put implementation in `internal/app/consult.go`; it should reuse `discoverConfigured`, select one installed agent, and run a consult prompt that forbids protocol writes.

UX:

- `parley consult [--dir DIR] [--timeout DURATION] <agent> "<question>"`; if the question arg is absent, read stdin.
- The facilitator captures stdout/stderr and writes the durable artifact itself so consult remains read-only-ish for the agent.
- Artifact path: `parley-deck/consults/<YYYYMMDDTHHMMSSZ>-<agent>-<slug>.md`.
- Frontmatter: `artifact: consult`, `agent`, `question_slug`, `created`, `workspace_root`, `command`, `timeout_ms`, `stdout_log`, `stderr_log`, `quorum: false`.
- Append `parley-deck/consults/index.jsonl` with the same provenance. Consult artifacts are not ideas, rounds, signoffs, or quorum evidence.

## Required tests

- `internal/runner`: fake exec agents for no-first retry, stall kill, heartbeat payload, valid artifact plus exit 1 success with `agent_exit`, invalid artifact plus exit 0 failure, timeout/user kill not artifact-wins, overwrite skip before retry, and fix-up validation.
- `internal/runner` ACP tests: fake ACP server with no first protocol event, stalled updates, prompt error after valid artifact, and phase-specific artifact validation.
- `internal/agents` and `internal/config`: parse/default new supervision knobs; runtime matrix/event payload includes them.
- `internal/runstate`: project `failure_class`, `recovery_hint`, `agent_exit`, watchdog summaries, and heartbeat without corrupting terminal state.
- `internal/tui`: narrator/header output for recovery hints and watchdog events; heartbeat verbosity does not flood transcripts.
- `internal/driver` and `internal/app`: artifact-wins through `Result.Success`, `gitTreeClean` sets `GIT_OPTIONAL_LOCKS=0`, and driver `roundComplete` remains artifact-valid plus terminal-event based.
- Snapshot package: shared-clone command construction, stale/live pid-marker healing, fallback event, artifact move-back, and no deletion on move-back failure.
- `internal/app`: `consult` usage, stdin question, artifact/index schema, unknown agent, and no quorum/protocol-state mutation.

## Position

Accept all P1-P5 and P8, but slice them to avoid semantic drift:

1. Shared result decision + P3 artifact-wins across exec/ACP/implementation/review-consensus/fixup, with tests.
2. Supervision core for exec and ACP: first-event retry, stall, heartbeat, durable attempt ids.
3. P2 classifiers and consumer surfacing in runstate/TUI/app.
4. P4 env/git/docs hardening.
5. P5 review snapshots with move-back and fallback events.
6. P8 consult command and artifact/index format.
