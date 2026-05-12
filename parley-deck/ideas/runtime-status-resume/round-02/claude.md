---
agent: claude
idea: runtime-status-resume
round: 2
date: 2026-05-12
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since round 1

Three concrete changes:

1. **Adopt codex's `terminal` + `outcome` shape over my `Liveness=complete`.**
   The user-visible vocabulary becomes `outcome ∈ {completed, incomplete,
   failed}` for terminal runs and `liveness ∈ {unverified, idle}` for
   non-terminal runs. "complete" was overloaded with "succeeded".
2. **Drop the `parley run --participants <missing-ids> …` suggestion** from
   the next-action recommender. That command does not exist as a recovery
   path against an existing idea today, and recommending it pre-commits the
   next slice. The recommender becomes three cases: `parley answer …`,
   `parley resume <run-id>`, or "no recoverable action — see artifacts".
3. **Reference `hitl.Store.List()` rather than a `questions/*.json` glob**
   in FINAL. Confirmed by reading `internal/hitl/hitl.go` — files are
   indeed `.json`, but the design should depend on the package API, not
   the filename pattern.

The rest of my round-1 proposal stands.

## Responses to others

### @codex — round-01

Strong agreement. Codex's round-1 and my round-1 describe the same animal
from two angles; codex's round-2 already converges on the joint shape.
Three points to lock:

- **`terminal` + `outcome` vs `Liveness`.** I accept codex's framing.
  `RunSummary` exposes `Terminal bool`, `Outcome string` (empty when
  not terminal), and `Liveness string` (`unverified` | `idle`, empty when
  terminal). One value at a time, no ambiguity.
- **JSON schema.** Agree to mark it unstable in a code comment on the
  struct. I would still ship it in this slice — without it, our own
  `app_test.go` has to parse plain text, which is brittle. Counter-proposal
  to anyone who wants to defer JSON: ship it but cap it at the
  `RunSummary` fields plus an `agents:[]` and `questions:[]` array. No
  event log, no log bytes, no nesting beyond two levels.
- **TUI `Done=nil` resume path.** Codex is right that this is currently
  an accidental behavior of `waitDoneCmd`. The FINAL should require:
  (a) explicit handling in `tui.RunLive` for a nil-`Done` resumed view,
  (b) a unit test that constructs `LiveOptions{Done: nil}` and asserts
  the program exits cleanly on `q`. That keeps it from regressing.

No counter-proposal needed for codex's round-2 — we are aligned.

### @gemini — round-01

Three disagreements; each with a counter-proposal.

**1. Retry failed/skipped agents inside `resume`.**
Disagree for this slice. Appending a second `agent.started` for the same
agent into the same `events.jsonl` breaks the reducer's "one lifecycle
per agent per run" assumption — every consumer (`ProjectEvents`, the
TUI agent table, any future JSON consumer) would need to learn
"last-event-wins per agent" semantics. That is a real design surface
(artifact overwrite policy, what happens if an orphaned subprocess from
the previous run is still alive, how to surface the retry history in
the UI) and it deserves its own idea.

*Counter-proposal:* Retry is a separate follow-up idea
(`runtime-retry-failed`). For this slice, `resume` is read-only. When
the user wants to retry, they re-run with the same idea slug, which
creates a fresh `runs/<new-id>/` with clean event semantics. We accept
"two runs per idea is fine" as the cost.

**2. `pid` field on `agent.started` events.**
Disagree. A PID without a supervisor contract is a foot-gun: across
reboot, PID reuse, or the agent forking, "is this PID alive?" answers
"yes" for the wrong process. We would have to add at least
`(pid, start_time, hostname)` to be safe, and then we are halfway to a
supervisor design without admitting it.

*Counter-proposal:* No PID in events for this slice. Liveness for
non-terminal runs stays `unverified` and is labeled with the last-event
age. If a future supervised-run design needs PIDs, it can add a
`run.supervised` event (or sidecar file with a documented lifecycle).
Stale-detection is by event age only, default threshold 10 minutes,
configurable via `--stale-after`.

**3. `parley status --watch`.**
Disagree. A live-updating status view duplicates the surface that
`parley resume <run-id>` (TUI mode) already provides. Two TUI entry
points to the same data is one too many.

*Counter-proposal:* `parley status` prints a snapshot and exits.
`parley resume <run-id>` is the watch surface. If a user wants the
workspace-wide overview to update live, that is a separate UX feature
(`parley watch`?) and worth its own slice.

**Agreements:** run-id and idea-slug resolution, conservative process
claim, HITL question discovery via the durable files, and the warning
that filesystem drift can desynchronize from events. On filesystem
drift specifically: I would not have `status` check artifact presence
against events in this slice — that is protocol validation and creeps
into "did the round really happen". Keep `status` reporting only what
the event stream and HITL files say.

### @hermes — round-01

Two disagreements with counter-proposals; one alignment confirmation.

**1. Optional `runs/<run-id>/pid` supervisor file.**
Disagree for this slice. Same reasoning as my response to gemini's PID
field: a PID file without a clearly specified lifecycle (when written,
when cleaned up, what happens on SIGKILL or reboot, how PID reuse is
detected) creates a "live" indicator we cannot back up. Hermes
explicitly calls it "best-effort" and "strictly advisory" — that
qualification is exactly the problem: users read "live" and assume
correctness.

*Counter-proposal:* No PID file in this slice. The liveness vocabulary
is `outcome` (for terminal runs) and `unverified | idle` (for
non-terminal). If "is the run still being driven?" becomes a real
question, design a supervised-run mechanism in a separate idea — at
which point we can store PID, start_time, hostname, and a heartbeat in
a documented format.

**2. Reading frontmatter from "latest round/review artifact" to derive
phase.**
Disagree. `00-prompt.md` `status:` is the canonical phase marker per
COOPERATION.md §4 (`round-N | consensus | final | abandoned`). Inferring
phase from "the latest round/review file that happens to exist" can
silently disagree with `00-prompt.md` if files are produced
out-of-order or if an agent submits late.

*Counter-proposal:* For `parley status`, the idea-level phase comes
from `00-prompt.md` `status:` only. Artifact presence is reported
separately as a list (`round-01/codex.md ok`, `round-01/gemini.md
missing`) without inferring a phase from it. Deeper protocol validation
is a follow-up.

**Alignment:** "view + state restoration for the next human or agent
step" is the exact right framing for resume. I would put that sentence
verbatim into FINAL's user-facing description.

## New concerns / questions

- **`tui.RunLive` resume mode.** Codex correctly flagged that
  `Done=nil` is currently incidental. FINAL must call out: explicit
  `Resume bool` (or equivalent) on `LiveOptions`, plus one TUI test
  exercising the resumed-view path. Otherwise this regresses the next
  time someone refactors `waitDoneCmd`.
- **Where does the projection package live?** Codex suggests
  `internal/runtime` or `internal/store`; I suggested `internal/runstate`.
  Soft preference for `internal/runstate` — `runtime` collides with the
  Go stdlib package and `store` already exists for the events writer.
  Not a blocker; FINAL drafter picks one.
- **Run-id resolution edge case.** If `parley resume <slug>` finds zero
  matching runs (the idea exists in `ideas/<slug>/` but no run has been
  recorded yet), we should fail with "no runs for idea <slug>", not
  fall back to a listing. Spelled out so the implementer is not tempted
  to be clever.
- **`events.jsonl` re-read cost.** `parley status` reads every run's
  events file. At today's volumes this is fine; if it ever isn't, a
  per-run `summary.json` snapshot is the natural next step. Flag in
  FINAL "References / future optimizations"; do not pre-build it.

## Current proposal

A converged plan for FINAL. Bullet points marked `(locked)` already
have agreement across codex round-2 and this file; `(open)` items need
one more pass.

**Shared package: `internal/runstate/`** (locked)

- Move `ProjectEvents`, `RunState`, `AgentState`, `EventSummary` from
  `internal/tui/live.go` into `internal/runstate/`. The TUI imports from
  the new package; no behavior change. Move the relevant tests with the
  reducer.
- Add `LoadRun(runDir string) (RunSummary, error)` and
  `ListRuns(root string) ([]RunSummary, error)`. `LoadRun` reads
  `run.created` for idea/task/mode/participants, calls `ProjectEvents`,
  and uses `hitl.New(runDir).List()` for open questions.

**`RunSummary` shape** (locked, modulo field naming)

- `RunID string`
- `IdeaSlug string` (from `run.created.data.idea`; `"unknown"` if absent)
- `Mode string`
- `Participants []string`
- `Terminal bool`
- `Outcome string` — `"completed" | "incomplete" | "failed" | ""`
- `Liveness string` — `"unverified" | "idle" | ""` (empty when terminal)
- `LastEventAt time.Time`, `LastEventAge time.Duration`
- `OpenQuestions int`
- `Agents []AgentSummary` (state, last event, log paths, elapsed)
- A JSON-stable comment marking the struct as unstable for now.

**`parley status` surface** (locked)

- `parley status` — workspace overview: transport, idea table (phase
  from `00-prompt.md` `status:`), newest run per idea with outcome /
  liveness / open-questions / last-event age. Plain text, deterministic
  order, latest first.
- `parley status --run <run-id>` — detail view for one run: agent table
  with state/elapsed/last-event/log paths, the 8 most recent events,
  open HITL questions with IDs, suggested next action (one of
  `parley answer …`, `parley resume <run-id>`, or "no recoverable
  action").
- `parley status --idea <slug>` — detail view for the newest run of
  that idea; error if no runs.
- `parley status --json` — `RunSummary` (or array) plus agents and
  questions arrays, unstable-marked.
- `--dir` and `--stale-after` are honored.

**`parley resume RUN_OR_IDEA`** (locked)

- Resolution: exact `runs/<id>/` first; else newest run whose
  `run.created.data.idea` matches the slug; else clear error listing
  available run IDs.
- Default: open `tui.RunLive` in resume mode (explicit `Resume` flag in
  `LiveOptions`, `Done=nil` handled deliberately, `Cancel` is a no-op).
  TUI is a read-only tail of `events.jsonl`, logs, and `questions/`.
  The `a` keybinding for HITL answers stays functional because
  `hitl.Answer` writes durable files.
- `--no-tui` — print the same body as `parley status --run <id>` and
  exit.
- Not in this slice: subprocess re-attachment, retry of failed agents,
  PID files, lockfiles, supervisor daemon, multi-run-per-idea selection
  beyond newest, run-id prefix matching.

**Test plan** (locked)

1. `internal/runstate/runstate_test.go` — fixture events files exercising
   each `Outcome`/`Liveness` permutation (run completed, run incomplete,
   run failed, agent unverified with stale event, all agents terminal
   but no round event).
2. `internal/tui/live_test.go` — kept for rendering only; projection
   tests move to `runstate_test.go`.
3. `internal/app/app_test.go` — `parley status`, `parley status --run`,
   `parley status --idea`, `parley status --json` against temp
   workspaces; deterministic golden text.
4. `parley resume <run-id> --no-tui` text matches
   `parley status --run <run-id>`.
5. `parley resume nonexistent` exits non-zero with a listing of
   available run IDs (or "no runs yet").
6. New TUI test: `tui.RunLive` with `LiveOptions{Resume: true, Done:
   nil}` exits cleanly on `q` (regression guard for the resume path).

**Verification** (locked)

- `go test ./...` green.
- Manual: produce a run via `parley run --no-tui`, then exercise
  `parley status`, `--run`, `--idea`, `--json`, and both `--no-tui` and
  TUI `parley resume`.

**Explicitly deferred to later ideas** (locked)

- Retry failed/skipped agents (`runtime-retry-failed`).
- Supervised-run mechanism with PIDs, heartbeats, lockfiles.
- `parley status --watch` or `parley watch`.
- `--all-runs` and run-id prefix matching.
- Protocol-validation reporting beyond artifact presence listing.
- `state.json` summary cache.

**Open for round 3 if needed**

- Package name: `internal/runstate` vs `internal/runtime` vs other.
- Whether JSON ships in this slice (I say yes, unstable-marked; codex
  is neutral; gemini and hermes did not explicitly weigh in).
- Default value of `--stale-after` (I propose 10m).
