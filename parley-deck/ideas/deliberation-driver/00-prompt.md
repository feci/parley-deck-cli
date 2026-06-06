---
idea: deliberation-driver
author: user
created: 2026-06-05
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + state-machine / cursor design
  codex: Go idioms, concurrency/locking, test design
  agy: protocol-correctness vs COOPERATION.md, transport edge cases
  hermes: donor-pattern fidelity (kanban dispatcher ready-task semantics → Go port)
transport: local-dir
transport-note: >
  Project default Transport is github-pr, but the owner explicitly selected
  local-dir for THIS idea (the driver is itself gated to local-dir). This idea
  runs on branch feature/deliberation-driver with filesystem artifacts under
  parley-deck/; the global Transport header is unchanged. See
  inbox/claude-to-all_deliberation-driver_transport-local-dir.md.
cross_review_rounds: 1
status: final
---

## Problem / idea

`parley-deck-cli` orchestrates multi-agent deliberation but **STALLS after
round-01**. Every run ends identically: `run.json` status=`running` forever,
idea_status=`round-01`, participants `pending`; `events.jsonl` ends at
`round.completed` with nothing after. Agents start, write round-01 artifacts, the
round closes — then nothing advances the idea to round-02 / consensus / final.
The orchestrator is a **batch executor, not a driver**.

### Root cause — verified at file:line in v1.14.1 (confirm in round-01)
- `parley run` runs exactly one round and returns: `internal/app/app.go:1654`
  (RunRoundOne, --no-tui) / `:1671` (RunRoundOneAsync, TUI), returns at `:1698`.
- `--auto` today only skips the launch prompt + starts a HITL auto-answerer
  (`runcontrol.StartAutoAnswerer` at `app.go:1652`/`:1669`). It does NOT progress
  phases.
- `round.completed` is emitted (`internal/runner/runner.go:145`) but consumed by
  NOBODY.
- The planner jumps round-01 → consensus with no cross-review step:
  `internal/runplan/runplan.go:110-122` emits `KindDraftConsensus` right after
  round-01. There is NO `KindOpenNextRound` action
  (`internal/runaction/action.go:9-14`). A BLOCKED consensus only yields
  `KindInspect` "…before opening another round" (`runplan.go:152-162`) — it never
  actually opens round-02.
- No code path writes idea_status `round-02`: `grep -rn "round-02" internal/*.go`
  returns zero non-test hits.

### Key enabler (the hard part already exists; it is just never called)
- `internal/runner/runner.go:544 RunRound(N≥2)` — reseeds participants with prior
  rounds (`gatherPriorRounds`) and writes `round-NN/<agent>.md`. Today only called
  by `internal/app/pipeline_cmd.go:42` (launchBlockRound), never by `parley run`.
- `internal/runner/phase58.go:20 RunImplementation`, `:36 RunReviewRound` exist.
- consensus: Draft (`consensus.go:105` → updateIdeaStatus "consensus" `:133`),
  Finalize (`:196` → "final" `:236`), Reopen (`:246`, requires TriageBlocked →
  reverts to latestRound `:283`), triage states
  TriageReady/Reserved/Blocked/Partial/Malformed (`:23-27`).
- Durable-cursor template to COPY: `internal/pipeline/run.go PipelineRun.Save`
  (`:82`, atomic tmp+rename); ready-set pattern `internal/pipeline/dag.go
  ReadyBlocks` (`:11`); `internal/pipeline/executor.go Driver.Advance` (`:68`).
  Round-from-disk inference already exists in `internal/runstate`
  (`inferCurrentRound`) — reuse for cursor Rebuild.

## Proposed direction (a STARTING proposal — round-01 is independent; challenge it)

Harvest the **PATTERN** of Hermes' kanban dispatcher (dispatch_once +
recompute_ready + crash-safe re-entry), NOT its SQLite kernel, into a small new Go
package `internal/driver` (~300 LOC). Model each deliberation phase as a
dependency-gated "task" on a degenerate-linear graph
round-01 → round-N → consensus → final → impl → review(↔fix-up). Provide:

1. `internal/driver/cursor.go` — `Cursor{Phase, CurrentRound, IdeaStatus,
   RoundsRun, MaxRounds}` with atomic Save (copy `pipeline/run.go:82`) and
   `Rebuild(runDir)` that DERIVES the phase purely from disk (00-prompt.md status,
   round-NN dirs, consensus.Status, FINAL.md).
2. `internal/driver/driver.go` — `Advance(ctx,*Cursor)` ONE re-entrant tick:
   recompute ready phase (pure, disk-derived) → run one gated action → save cursor.
   When round-N is complete (all participant artifacts present AND round.completed)
   the gate PROMOTES round-N → round-(N+1) via RunRound; the consensus gate
   promotes consensus → final autonomously.
3. `internal/driver/loop.go` — outer poll loop (mirror pipeline auto) + an advisory
   `driver.lock` (PID) on the run dir — the ONLY concurrency control (no claim_lock).
4. Gates:
   - round gate: all round-N artifacts written + round.completed → open
     round-(N+1) when cross-review is wanted, else → consensus.
   - consensus gate (reuse consensus triage): TriageReady/Reserved → Finalize →
     FINAL.md; TriagePartial → request signoffs (auto in --auto; else stop for
     HITL); TriageBlocked → consensus.Reopen + CurrentRound++ → new cross-review
     round seeded by the BLOCK counter-proposal (bounded by MaxRounds).
5. Wiring:
   - `internal/runaction` + `internal/runplan`: add `KindOpenNextRound` so even
     read-only `parley continue` surfaces "open round-02".
   - `internal/app/app.go runTask`: when --auto AND transport==local-dir, after
     RunRoundOne hand control to `driver.Run(...)`. Non-auto keeps one-shot.
   - `internal/app/app.go runContinue` (`:1019`): with --auto, EXECUTE the next
     action instead of only printing it.
6. `internal/driver/driver_test.go`: table tests — round-01→round-02 promotion
   calls RunRound(2); crash test deletes/corrupts cursor and asserts Rebuild
   recovers the same phase; BLOCK test asserts Reopen + CurrentRound++.

## Constraints (non-negotiable; reviewers BLOCK violations)

- **Simplicity First** (CLAUDE.md): do NOT import or reimplement Hermes' SQLite
  task kernel, claim_lock/heartbeat/zombie-reaping, blackboard, or
  DAG/topological machinery. The graph is linear; a simple ordered switch computes
  the ready phase.
- **Disk is the single source of truth**; the cursor is a rebuildable cache. Every
  gate and `readyPhase` MUST be a pure function of on-disk artifacts → re-entry is
  idempotent. `RunRound` must skip agents whose artifact already exists
  (Overwrite=false).
- **Auto-advance gated behind --auto AND local-dir transport ONLY.** NEVER
  auto-drive GitHub/GitLab (B/C rely on humans switching PR/MR labels). Read
  transport before constructing the driver.
- **Signoffs MUST be real participant agent invocations** (each agent authors its
  own ✅/🟡/❌ via `consensus.AppendSignoff`). The driver MUST NOT fabricate ACCEPT
  signoffs.
- **MaxRounds circuit breaker (default 4):** a BLOCK→new-round loop must terminate
  and escalate to the user, never spin forever.
- **English-only** in every `parley-deck/` file. One file per agent per round;
  never edit another agent's file; append-only signoffs.

## Non-goals

- No SQLite, no daemon/server, no worker pool, no cross-process scheduler.
- No new storage engine; reuse run.json / events.jsonl / runstate / consensus.
- Not unifying with §12 pipeline blocks here (pipeline drives coarse whole-idea
  blocks; this driver advances rounds WITHIN one deliberation). Note the
  relationship; do not merge them here.
- No change to how agents are launched (discover.go / acp_specs.go untouched).

## Deliverables (in order)

1. **FINAL.md**: agreed driver design (cursor schema, Advance state machine, gate
   table, crash-recovery contract, file-by-file plan, MaxRounds + transport
   gating decisions).
2. **MINIMAL FIRST SLICE implemented + proven first:** ONLY PhaseRound (cursor +
   readyPhase + promote round-01→round-02 via RunRound), wired behind
   `parley run --auto --no-tui` on a local-dir workspace. Acceptance: a real run
   writes round-01/<agent>.md for all participants, THEN round-02/<agent>.md, sets
   00-prompt.md status=round-02, and events.jsonl shows a SECOND
   round.started/round.completed pair — the follow-on events absent today.
3. Then extend the same loop through consensus → final → impl → review/fix-up.

## Round-01 focus questions (answer independently)

Is the proposed driver the right shape? cursor-vs-disk authority; idempotency
holes; the cross-review-vs-straight-to-consensus decision; transport gating; how
to detect "cross-review wanted"; MaxRounds + escalation; what the minimal first
slice must prove.
