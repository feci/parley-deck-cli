---
agent: codex
idea: tui-live-steering
round: 2
date: 2026-06-06
---

## Cross-review

### @claude

I agree with the round-02 synthesis and consider it ready to become the FINAL contract. The important runner/app decisions are now correct: autocomplete is a dedicated non-modal suggest state; kill is a per-attempt `Handle` operation; steer is a fresh single-agent attempt in a per-steer directory; TUI access is only through injected seams; CLI and observational runs are record-only.

The `Handle` contract is the right spine, but the implementation should make one distinction explicit: `RunRoundOneAsync` already returns a handle-like object today for the live async path, but it must be expanded into the authoritative owner of attempt state, not wrapped by an additional app-side registry. Add the missing fields there: a mutex, segment mutex, active attempts keyed by agent id, depth-1 pending/running steer state, root context, run metadata, and the per-attempt cancel funcs. `RunRoundOne` can remain synchronous or delegate internally, but the live path must return a `*Handle` that has `KillAgent` and `RunSteerAttempt`.

The segment-id lock must cover the complete critical section: compute next segment id, append `run.segment_started`, and publish the segment id to the attempt. Locking only `nextSegmentID` is insufficient because two goroutines can still append the same computed id. For in-process TUI steers, a `Handle.segmentMu` is enough. Cross-process CLI remains record-only, so it does not participate in steer segment allocation.

Per-steer directories plus no artifact validation cover the main clobber paths, with one extra caution: do not call unchanged `runAgent` if it still derives `AgentDir`, `StdoutPath`, `StderrPath`, output artifact path, skip-if-artifact-exists, or stdout-fallback behavior from the normal round path. The implementation needs an attempt-mode/path override so a steer cannot truncate `agents/<id>/stdout.log`, cannot skip because `round-01/<id>.md` already exists, and cannot write or validate a protocol artifact.

### @codex

I agree with my round-01 runner design, with Claude's refinements adopted. The concrete seams should stay plain and stable:

```go
type SteerRequest struct {
    AgentID string
    Text    string
}

type SteerResult struct {
    ID         string
    Status     string
    SegmentID  string
    StdoutPath string
    StderrPath string
    ReplyPath  string
}

type KillResult struct {
    AgentID   string
    Killed    bool
    SegmentID string
    Message   string
}
```

`LiveOptions.SubmitSteer` and `LaunchResult.SubmitSteer` can expose the TUI-shaped result, while `runner.Handle.RunSteerAttempt` can use runner-native request/result types. `activateRun` must copy both `SubmitSteer` and `KillAgent`; this is not optional and needs a model-level regression test because otherwise Home-launched runs lose the new behavior silently.

Event ordering should be deterministic:

1. TUI submit calls app seam.
2. App records `steer.requested` and `steer.delivered`.
3. If record-only, return a record-only status.
4. If executable, `Handle.RunSteerAttempt` accepts or rejects under lock.
5. When the attempt actually starts, append locked `run.segment_started`, then `steer.reply_started`, then register/launch the process.
6. On process completion, append exactly one terminal steer event: `steer.replied` or `steer.reply_failed`.

Do not emit `agent.killed` unless `KillAgent` finds a live registered attempt and wins the race to cancel it. If normal completion has already deregistered the attempt, return `Killed:false` and do not synthesize a killed event. If a killed process later exits through the normal failure path, projection should prefer the prior `agent.killed` marker over a generic context-canceled failure label.

### @agy

I agree with your UX choices as adopted by Claude: the slim suggestion overlay, `ctrl+k` with modal y/N confirmation, distinct steer prompt, inline steer divider, spinner, and same-agent-tab reply are the correct user-visible behavior.

The only runner/app constraint I want preserved in the UX copy is that "queued" means depth 1. A steer to a busy agent is accepted only if no steer is already pending/running for that agent. A second steer must be rejected clearly rather than replacing the first, because replacement creates audit ambiguity unless we add a cancellation event. v1 should reject with a visible "<agent> is already replying" message.

For kill, the UI can show a `KILLED` badge, but the runner should still treat the attempt as terminal failed/killed rather than successful. That avoids accidentally promoting an incomplete round artifact while still allowing the rest of the run to continue and allowing a later steer attempt.

## Counter-proposals

No counter-proposal against the synthesis. The only additions are implementation constraints:

- Lock `nextSegmentID + appendSegmentStarted + attempt segment assignment` as one critical section.
- Make steer attempt paths explicit instead of relying on normal `runAgent` defaults.
- Emit one terminal steer reply event per accepted steer attempt.
- Keep CLI and observational `/open` as record-only unless a live in-process `Handle` is available.

## Confirmed for FINAL

- `RunRoundOneAsync` returns the live `Handle`; extend that handle with `KillAgent` and `RunSteerAttempt` rather than creating an app-side process registry.
- `KillAgent` cancels only the selected attempt's child context, emits `agent.killed` only when it finds a live attempt, and never calls run-wide cancel.
- `RunSteerAttempt` uses a fresh segment with `reason:"steer"` and `steer_id`, allocated under the handle segment lock.
- Steer attempts write only under `runs/<runID>/agents/<agentID>/steers/<steerID>/` and bypass round artifact skip, fallback, and validation.
- Keep `steer.requested`/`steer.delivered`; add `steer.reply_started`, `steer.replied`, and `steer.reply_failed`.
- Queue depth is 1 per agent; second steer is rejected.
- `LiveOptions` and `LaunchResult` both carry `SubmitSteer` and `KillAgent`; `activateRun` copies both.
- `internal/tui` imports neither runner nor app.
- Autocomplete is the dedicated non-modal suggest sub-mode from the synthesis, with no Tab cycling.

## Remaining risks

- Process descendants may survive `exec.CommandContext` cancellation if an agent command spawns children. v1 can use context cancellation, but the runner tests should expose the limitation and process-group signaling should be tracked as follow-up.
- Projection must distinguish steer segments from normal round segments so a steer cannot make the main round artifact look newly pending or overwrite the visible round status.
- Accepted steers that start after the main round `done` channel closes need a clear handle lifetime rule. I recommend accepted steers continue to terminal reply/failure, but new steers after the handle context is canceled are record-only or rejected before execution.
- Fake-agent tests must cover kill-vs-normal-completion, two-agent survival after one kill, duplicate steer rejection, per-steer stdout path use, and the `activateRun` seam copy.
