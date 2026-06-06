---
agent: codex
idea: tui-live-steering
round: 1
date: 2026-06-06
---

## Summary

The clean runner/app design is to treat a steer reply as a new single-agent attempt, not as live stdin into the existing process. That matches the current one-shot model in `runAgent` and `CommandFor` (`internal/runner/runner.go:275`, `internal/runner/runner.go:659`) and avoids pretending there is an interactive session. The runner should expose this through the existing async `Handle`, while the TUI sees only injected function seams on `LiveOptions`.

I recommend queue depth 1 per agent for steer execution. If the target agent is currently running, the steer is accepted durably and runs after that agent's active attempt exits; if a steer attempt is already pending/running for that agent, reject a second one with a clear "already replying" error. Parallel invocations against the same agent are the highest-risk option because today's paths are per-agent and overwrite-prone (`agents/<id>/stdout.log`, `stderr.log`, and normal protocol artifacts at `internal/runner/runner.go:277-282`). "Next round only" is too invisible for the owner.

Per-agent kill should be a run-handle operation that cancels only the selected agent's child context. It should never call the run-wide `cancelRun` from `app.go:1744-1758`; that remains ctrl+c / whole-run cancellation. Killed agents should produce a durable `agent.killed` event and terminal failed/killed state, while the round WaitGroup continues waiting for the remaining goroutines.

Autocomplete is mostly TUI-owned. The runner/app contribution should be a small, shared command specification table for command names and metadata; agy/claude should own the UX details.

## Proposed approach (with concrete Go signatures/seams)

Add runner execution control around each agent process:

```go
type AgentRunKind string

const (
    AgentRunRound AgentRunKind = "round"
    AgentRunSteer AgentRunKind = "steer"
)

type AgentAttempt struct {
    AgentID   string
    SegmentID string
    Kind      AgentRunKind
    SteerID   string
    Cancel    context.CancelFunc
}

type Handle struct {
    RunID  string
    RunDir string
    done   chan struct{}

    opts Options
    rootCtx context.Context

    mu       sync.Mutex
    results  []Result
    active   map[string]AgentAttempt
    steerBusy map[string]bool
    wg       sync.WaitGroup
}
```

`RunRoundOneAsync(ctx, opts)` at `internal/runner/runner.go:71` should initialize the handle with `opts`, `rootCtx`, `active`, and `steerBusy`, then call an internal `handle.runAgent(ctx, opts, agent, attemptOptions)`. The synchronous `RunRoundOne` can keep the current simpler behavior or use a short-lived handle internally. The important part is that every async TUI-run attempt registers before `cmd.Run()` and unregisters with `defer`, including ACP mode.

Concrete kill seam:

```go
type KillResult struct {
    AgentID string
    Killed  bool
    Message string
}

func (h *Handle) KillAgent(agentID string) KillResult
```

`KillAgent` locks `h.mu`, finds `active[agentID]`, marks it killed or records a `killed` flag on the attempt, unlocks, calls only that attempt's cancel func, then appends:

```json
{"type":"agent.killed","data":{"agent":"agy","segment_id":"segment-0003","kind":"round|steer","steer_id":"..."}}
```

The running attempt then exits through the existing `agent.failed` path, but should set `Result.Killed bool` and prefer `ExitError = "killed"` over generic `context canceled`. Projection can map either `agent.killed` or failed-with-killed to a killed/failed terminal badge. Race safety: if normal completion wins and unregisters first, `KillAgent` returns `Killed:false` and does not synthesize a killed event.

Concrete steer seam:

```go
type SteerAttemptRequest struct {
    AgentID string
    Text    string
}

type SteerAttemptResult struct {
    SteerID    string
    SegmentID  string
    StdoutPath string
    StderrPath string
    Accepted   bool
    Message    string
}

func (h *Handle) RunSteerAttempt(ctx context.Context, req SteerAttemptRequest) (SteerAttemptResult, error)
```

The app wires that to the TUI without importing runner into `internal/tui`:

```go
// internal/tui/live.go
type SteerRequest struct {
    Target steer.Target
    AgentID string
    Text string
}
type SteerResult struct {
    ID string
    Status string
    SegmentID string
    StdoutPath string
}
type SteerFunc func(SteerRequest) (SteerResult, error)
type KillAgentFunc func(agentID string) error

type LiveOptions struct {
    ...
    SubmitSteer SteerFunc
    KillAgent   KillAgentFunc
}

type LaunchResult struct {
    ...
    SubmitSteer SteerFunc
    KillAgent   KillAgentFunc
}
```

`activateRun` must copy these fields along with `Cancel` at `internal/tui/live.go:380-388`, otherwise runs launched from the Home tab via `newLaunchFunc` (`internal/app/app.go:2016-2048`) will silently lose live steering/kill.

`submitSteer` in the TUI (`internal/tui/live.go:1106`) should call `m.opts.SubmitSteer` when present. The app-level implementation should still call `steer.Submit` first so the durable `steer.requested` and `steer.delivered` events remain exactly as established at `internal/steer/steer.go:87-123`; then it calls `handle.RunSteerAttempt` for agent-targeted steers. For deck-level steers, keep recording only unless FINAL defines fan-out semantics.

Runner entry:

```go
func (h *Handle) RunSteerAttempt(ctx context.Context, req SteerAttemptRequest) (SteerAttemptResult, error) {
    // validate agent exists in h.opts.Agents and h.opts.Idea.Participants
    // record/inspect busy state under h.mu
    // if active round attempt exists, set one pending slot and spawn waiter goroutine
    // if no active attempt, spawn immediately in h.wg
}
```

For execution, do not reuse the normal output artifact path. Add an internal options struct or extend `Options` carefully:

```go
type attemptPaths struct {
    AgentDir   string
    OutputPath string
    StdoutPath string
    StderrPath string
    ValidateArtifact bool
}
```

For a steer attempt, use:

```text
parley-deck/runs/<runID>/agents/<agentID>/steers/<steerID>/stdout.log
parley-deck/runs/<runID>/agents/<agentID>/steers/<steerID>/stderr.log
parley-deck/runs/<runID>/agents/<agentID>/steers/<steerID>/reply.md
```

This avoids clobbering `stdout.log` in `internal/runner/runner.go:277-279` and avoids validating/writing `round-01/<agent>.md` again. Emit a new segment before starting:

```go
opts.RoundLabel = "steer/" + steerID
opts.SegmentID = appendSegmentStarted(opts, "steer", []string{agent.ID})
```

Add event fields to `run.segment_started`: existing `segment_id`, `reason`, `round`, `targets` are enough, but `steer_id` would make correlation easier if appendSegmentStarted accepts optional metadata.

Prompt:

```go
func BuildSteerPrompt(agent agents.Discovery, opts Options, steerID, steerText, replyPath string) string
```

The prompt should say: you are the same Parley participant; answer only this follow-up; do not edit protocol artifacts or other agents' files; write a concise reply to `replyPath` and also print the reply to stdout. Context should include the idea `00-prompt.md`, the agent's latest own artifact if it exists, and the tail of that agent's current stdout log. It should not include other participants' round-01 files during Phase 1, because that violates independent-analysis isolation. For later rounds, context can be relaxed by phase rules, but the first implementation can stay conservative.

Reply events:

```json
{"type":"steer.reply_started","data":{"id":"steer-...","agent":"agy","segment_id":"segment-0003","stdout":"..."}}
{"type":"steer.replied","data":{"id":"steer-...","agent":"agy","segment_id":"segment-0003","reply":"...","reply_path":"...","stdout":"...","duration_ms":1234}}
{"type":"steer.reply_failed","data":{"id":"steer-...","agent":"agy","segment_id":"segment-0003","error":"...","stdout":"..."}}
```

Keep existing `steer.delivered` but update its status vocabulary: `queued` when accepted behind an active attempt, `running` when started immediately, `failed` on failure, `replied` on success. If changing `steer.Submit` is too invasive, append a second `steer.delivered` for status transitions and let projection use the latest event.

Reply surfacing: the TUI should not tail many arbitrary files directly by guessing paths. The event projection should learn stdout/reply paths from `agent.started` or `steer.reply_started`; when the active tab is `agent:<id>`, show the normal stdout plus steer reply entries for that agent in chronological order. A first implementation can append a synthetic line to that agent's transcript buffer when `steer.replied` arrives:

```text
[steer steer-0002 replied]
<reply text>
```

While running, show status from `steer.reply_started`: "agy is replying to steer-0002..." in the status row and/or the agent tab badge.

Autocomplete mechanism:

```go
type commandSpec struct {
    Name string
    Usage string
    RequiresRun bool
    OpensPicker bool
}

var commandSpecs = []commandSpec{
    {Name: "/help"}, {Name: "/status"}, {Name: "/follow"},
    {Name: "/deck", Usage: "/deck <text>"}, {Name: "/answer", OpensPicker: true},
    {Name: "/open", OpensPicker: true}, {Name: "/home"}, {Name: "/quit"},
}
```

`runCommand` and help text should consume the same table rather than duplicating names. A dedicated suggest row is cleaner than reusing `pickerState` for autocomplete because `pickerState` currently clears input on `openPicker` (`internal/tui/live.go:1238-1245`) and owns typing as a filter; command suggestion needs to preserve and edit `inputText`. Tab while `inputText` starts with `/` should complete to the longest common prefix; if exactly one match, complete the command plus a trailing space when it accepts args. Arrow/Enter can choose from the suggest row. Normal tab switching remains when input is empty or not slash-prefixed; this changes the `tab` branch at `internal/tui/live.go:770-772`.

Focus-question answers:

1. Steer execution model: queue depth 1 per agent, run after the current attempt if the agent is active, otherwise spawn a fresh single-agent attempt immediately. Prompt gets steer text, idea prompt, the target agent's own latest artifact, and stdout tail only. Reply is `steer.replied` plus `reply.md`/stdout in a per-steer directory.
2. Per-agent kill: register per-attempt context cancel funcs, not just `*exec.Cmd`. `exec.CommandContext` already observes the context (`internal/runner/runner.go:671`), and ACP uses a context too (`internal/runner/acp.go:44-55`). On Unix, if later tests show child subprocess leaks, add process-group signaling under `internal/acp/sysproc_*` style platform split, but start with context cancellation. Killed state is terminal failed/killed for that attempt; the run continues.
3. Autocomplete: use a dedicated suggest state backed by `commandSpecs`; Tab completes common prefix; arrows/Enter choose; argument completion is only for bare `/open` and `/answer` via the existing picker, not arbitrary arg completion in this slice.
4. Reply surfacing: append reply markers into the agent's existing transcript view, backed by event data and per-steer logs. Show "replying..." while `steer.reply_started` is present without a terminal reply/failed event.
5. Concurrency and safety: no parallel same-agent attempts. Use handle mutex for `active`/`steerBusy`, store append for durable events, and per-steer dirs to avoid stdout/artifact overwrite. If the run ended but the app still has a live handle, allow a steer attempt only while the parent context is not canceled; otherwise return "run is not active; steer recorded only" or reject before `steer.Submit` depending on desired UX. Observational `/open` runs have no handle seam, so they should record only.
6. Seams and testability: TUI gets `SubmitSteer` and `KillAgent` funcs on `LiveOptions`/`LaunchResult`; app binds those to the runner handle; runner exposes `Handle.RunSteerAttempt` and `Handle.KillAgent`. Tests: fake agent command that sleeps for kill; fake print command for steer reply; TUI model tests for Tab completion, kill command/key dispatch, and `submitSteer` using a fake seam; race test with `go test -race` around kill-vs-completion if possible.

## Concerns / open questions

The largest app seam question is whether the CLI `parley steer` should also execute a steer when the run is active. A CLI invocation cannot easily find the in-process `runner.Handle`, so I would leave CLI as durable record-only for now and make that explicit. TUI-created active runs can execute because app owns the handle.

The projection model needs a small extension. Existing `agent.started` points the TUI at one stdout path per agent, and segments reset targeted agents to pending (`internal/runner/runner.go:203-207`). A steer segment targeting one agent may make that agent look like it is in a new "round" unless runstate distinguishes `reason:"steer"`. The projection should not let a completed steer overwrite the main round's artifact status in confusing ways.

`RunFixup` duplicates command execution instead of using `runAgent` (`internal/runner/phase58.go:49-113`). If kill is implemented only inside `runAgent`, fix-up kills will not work. The control wrapper should either be shared by both paths or `RunFixup` should be folded back into the common execution helper.

ACP validation currently calls `ValidateRoundOneArtifact` in `finishACP` regardless of phase (`internal/runner/acp.go:121-130`), while stdio uses `validateArtifactForPhase`. That is an existing correctness smell and matters more once steer attempts intentionally have no normal round artifact. Steer attempts should bypass artifact validation and validate only that a reply file or stdout text exists.

For queue depth 1, "replace pending steer" might feel better than rejection, but replacement is less auditable unless it appends a cancellation event for the superseded steer. I prefer rejection for the first implementation.

## Risks

Race on segment IDs: `nextSegmentID` counts events without locking (`internal/runner/runner.go:225-236`). Normal rounds are sequential today, but steer attempts can be submitted while a round is running. Two concurrent steer attempts for different agents could both compute the same next segment. This needs either a store-level append lock around `nextSegmentID + Append`, a handle-level segment mutex for in-process TUI attempts, or a random suffix for steer segment IDs. Without this, runstate projection can corrupt segments.

Log clobbering: any design that reuses `agents/<id>/stdout.log` for steers will truncate the live round transcript because `os.Create` is used (`internal/runner/runner.go:336`, `internal/runner/phase58.go:78`). Per-steer stdout/stderr paths are mandatory.

Artifact clobbering/validation: calling `runAgent` unchanged for steer execution will skip if `round-01/<agent>.md` exists (`internal/runner/runner.go:292-304`) or may try stdout fallback into protocol artifacts (`internal/runner/runner.go:376-395`). Steer execution needs a distinct attempt mode with no protocol artifact fallback.

Whole-run cancellation bleed-through: `KillAgent` must never touch `LiveOptions.Cancel` or `cancelRun` from `app.go:1744`. It must cancel only the child attempt context. Tests should assert another sleeping fake agent finishes after killing one.

WaitGroup lifetime: steer goroutines should be tracked separately from the round `done` channel. The existing `Handle.done` closes when `RunRoundOne` returns (`internal/runner/runner.go:77-96`). If a steer is queued behind a running agent, it may outlive or start just as the main round completes. The handle needs a separate attempts WaitGroup or a policy that no steer starts after `done` is closed. I recommend: active TUI steers may continue after round completion if already accepted, but new steers after `done` are rejected unless the app deliberately creates a new handle mode.

TUI key collision: `tab` currently switches tabs (`internal/tui/live.go:770-772`). Slash autocomplete must intercept Tab only when `inputText` begins with `/`; otherwise it will regress existing navigation.
