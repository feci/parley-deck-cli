---
idea: tui-live-steering
phase: final
status: final
drafter: claude
implementer: claude
date: 2026-06-06
participants: [claude, codex, agy, hermes]
supersedes: consensus.md
---

# FINAL — TUI live steering, kill, and command autocomplete

Ratified by claude, codex, agy, hermes (all ✅ ACCEPT in consensus.md). Build spec for
Phase 5. Touches `internal/{tui,runner,app}` (+ tests). No new deps; the normal round
path and `--no-tui` stay intact.

## A. Autocomplete (internal/tui/live.go)
- `commandSpecs` table — the single source for `runCommand`, `/help`, suggestions, hints:
  ```go
  type commandSpec struct{ Name, Usage string; RequiresRun, OpensPicker bool }
  var commandSpecs = []commandSpec{
      {Name: "/help"}, {Name: "/status"}, {Name: "/follow"},
      {Name: "/deck", Usage: "/deck <text>"},
      {Name: "/answer", Usage: "/answer [qid text]", OpensPicker: true},
      {Name: "/open", Usage: "/open [slug|run]", OpensPicker: true},
      {Name: "/home"}, {Name: "/quit"},
  }
  ```
- Dedicated NON-modal suggest state on `liveModel`: `suggest bool`, `suggestItems
  []commandSpec`, `suggestIndex int`. Recomputed after each input edit: active when
  `strings.HasPrefix(inputText,"/") && !strings.Contains(inputText," ")` and ≥1 match and
  `transcriptHeight()>=3`; else off. Printable/backspace still edit `inputText` (then
  recompute). Esc closes suggest first (before the existing esc behavior).
- `updateMain` key handling while `suggest` is on: `tab` → complete longest common prefix
  of matches; if exactly one match → full name + trailing space when `Usage` has args;
  ↑/↓ move `suggestIndex`; `enter` accepts the highlighted command (no-arg → execute via
  `runCommand`; `OpensPicker` → set `"/cmd "` and hand to the existing bare picker).
- Conditional Tab: the existing `tab`/`shift+tab` → `switchTab` ONLY when
  `!strings.HasPrefix(inputText,"/")`; `up/left/down/right` tab-nav unchanged.
- Render: `renderSuggest(width)` slim list above the input row (reuse box style); hint
  shows `Tab complete · ↑/↓ pick · Enter run`.

## B. Kill one agent
### runner (runner.go)
- Extend `Handle`: add `opts Options`, `rootCtx context.Context`, `segmentMu sync.Mutex`,
  `active map[string]*attempt`, `steerBusy map[string]bool` (guarded by existing `mu`).
  `attempt{agentID, segmentID, kind, steerID string; cancel context.CancelFunc; killed bool}`.
  `RunRoundOneAsync` stores `opts`+`ctx` and inits the maps.
- A shared agent-run helper registers the per-attempt `context.WithCancel(parent)` cancel
  into `active[agentID]` before `cmd.Run()` and deregisters on exit (under `mu`). The round
  path uses it so round agents are killable.
  ```go
  type KillResult struct{ AgentID string; Killed bool; SegmentID, Message string }
  func (h *Handle) KillAgent(agentID string) KillResult
  ```
  Locks `mu`, finds `active[agentID]`; if absent → `{Killed:false}` (no event). Else mark
  `killed`, call its `cancel`, append `agent.killed{agent,segment_id,kind,steer_id}`. The
  canceled attempt returns via the existing failure path; `Result.Killed=true`. Never calls
  the run-wide cancel.
### app (app.go)
- Build `KillAgent func(agentID string) error` bound to `handle.KillAgent` and set it on
  `LiveOptions` (≈1751) and `LaunchResult` (≈2039).
### tui
- `LiveOptions.KillAgent KillAgentFunc` + same on `LaunchResult`; `activateRun` copies it.
- Modal sub-mode `confirmKillAgentID string` — the FIRST interceptor in `updateMain` (after
  ctrl+c): renders `kill <agent>? (y/N)` warn-styled, blocks all other keys; `y`/`enter` →
  `opts.KillAgent(id)` (transient `inputErr`/`statusMsg` on err/ok), `n`/`esc` cancels.
  `ctrl+k` opens it when on an agent tab whose state is running. Badge `KILLED` for killed.

## C. Steer round-trip
### runner (runner.go / a new steer.go in runner)
```go
type SteerAttemptRequest struct{ AgentID, Text string }
type SteerAttemptResult struct{ SteerID, SegmentID, StdoutPath, StderrPath, ReplyPath, Status string; Accepted bool; Message string }
func (h *Handle) RunSteerAttempt(ctx context.Context, req SteerAttemptRequest) (SteerAttemptResult, error)
```
- Validate the agent is a participant. Under `mu`: if `steerBusy[agentID]` → reject
  `{Accepted:false, Message:"<agent> is already replying"}`. Else set busy. If a round
  attempt for that agent is currently `active`, run the steer AFTER it (depth-1 queue: a
  waiter goroutine); else start immediately.
- Start: under `segmentMu`, atomically `nextSegmentID`+`appendSegmentStarted(reason:"steer",
  targets:[agent], steer_id)`+create `runs/<runID>/agents/<id>/steers/<steerID>/`+publish.
  Append `steer.reply_started{id,agent,segment_id,stdout}`. Register the attempt
  (kind=steer) so it's killable. Build the prompt with `BuildSteerPrompt`, run via the
  shared helper in ATTEMPT MODE (stdout/stderr/reply under the steer dir; NO skip-if-exists,
  NO protocol-artifact fallback, NO artifact validation — only require non-empty
  stdout/reply). On finish append exactly one of `steer.replied{id,agent,segment_id,
  reply_path,stdout,duration_ms}` / `steer.reply_failed{id,agent,segment_id,error,stdout}`;
  clear `steerBusy` under `mu`.
- `BuildSteerPrompt(agent, opts, steerID, text, replyPath)`: steer text + idea
  `00-prompt.md` + the agent's own latest artifact (if present) + stdout tail (~4KB).
  Instruct: answer only this follow-up; do not touch protocol artifacts or other agents'
  files; print the reply to stdout and write it to `reply.md`.
### app (app.go)
- `SubmitSteer` impl: call `steer.Submit` first (keep `steer.requested`/`delivered`); for
  `TargetAgent` with a live handle, call `handle.RunSteerAttempt` and return
  `{ID,Status,SegmentID,StdoutPath}`; for deck/observational/CLI → record-only status. Set
  on `LiveOptions` + `LaunchResult`.
### tui
- `SteerRequest{Target steer.Target; AgentID, Text string}`, `SteerResult{ID, Status,
  SegmentID, StdoutPath string}`, `SteerFunc func(SteerRequest)(SteerResult,error)`;
  `LiveOptions.SubmitSteer` + `LaunchResult.SubmitSteer`; `activateRun` copies it.
- `submitSteer` calls `opts.SubmitSteer` when set; if a `StdoutPath` comes back, attach it
  as a steer-reply buffer for that agent (tail it like a transcript) and show a
  divider `── steer <id>: "<query>" ──`, reply in faint style, a "<agent> is replying…"
  status while running, and the queued/error hints. Falls back to the old record-only
  message when no seam.

## D. Decoupling & scope
`internal/tui` imports neither runner nor app — only func seams on `LiveOptions`/
`LaunchResult`. Kill/steer target the async round `Handle` path used by `parley tui`;
fix-up/impl phases and deck-steer fan-out are out of v1; CLI `parley steer` record-only.

## E. Tests
- runner (`runner` pkg, fake agent script via a temp executable / `CommandFor` seam):
  `RunSteerAttempt` happy path (reply captured, per-steer paths, event order); duplicate
  steer rejected; `KillAgent` kills one sleeper while a second survives;
  kill-vs-completion race; run `-race`.
- tui (`live_test.go`, drive `Update`): Tab LCP completion; suggest show/filter/clear +
  conditional-Tab; confirm-kill modal (y/n, blocks keys); `submitInput` steer via fake
  `SubmitSteer`; `activateRun` copies `SubmitSteer`+`KillAgent`.

## Acceptance
`go build ./... && go vet ./... && go test ./...` green (repo-local caches). Manual-smoke
note in IMPLEMENTATION.md. Then Phases 6-8 review by codex/agy/hermes to zero agreed fixes.
