---
idea: deliberation-driver
status: final
drafted-by: claude
date: 2026-06-05
participants: [claude, codex, agy, hermes]
implementer: claude
consensus: parley-deck/ideas/deliberation-driver/consensus.md
---

## Final plan / specification

### Problem & approach
`parley run` is a one-shot batch executor that stalls after round-01 (verified:
no `KindOpenNextRound`; `runplan` jumps a complete round to consensus; `RunRound`
N≥2 exists but is never called from `parley run`). The fix is a small new package
`internal/driver` that owns the next tick: a disk-derived **Cursor** + a pure
**readyPhase** + a re-entrant **Advance** that drives an idea through
round → cross-review → consensus → final → impl → review/fix-up. The graph is
degenerate-linear with one back-edge (consensus BLOCK → reopen round), so it is a
readable ordered switch, NOT a DAG/scheduler. Implements consensus decisions
D1–D15 (see `consensus.md`).

### Cursor schema (`internal/driver/cursor.go`)
```go
type Phase string
const (
    PhaseRound     Phase = "round"
    PhaseConsensus Phase = "consensus"
    PhaseFinal     Phase = "final"
    PhaseImpl      Phase = "impl"
    PhaseReview    Phase = "review"
    PhaseDone      Phase = "done"
    PhaseBlocked   Phase = "blocked" // halted, escalation written
)

type Cursor struct {
    Phase        Phase  `json:"phase"`
    CurrentRound int    `json:"current_round"`
    IdeaStatus   string `json:"idea_status"`
    RoundsRun    int    `json:"rounds_run"`
    MaxRounds    int    `json:"max_rounds"` // default 4; only config-derived field
    UpdatedAt    string `json:"updated_at"` // RFC3339; informational
}
```
- `Save(path)` — atomic tmp+rename, same dir (copy `internal/pipeline/run.go:82`).
- `Load(path)` — best-effort; a missing/corrupt cursor is non-fatal (warn) and
  triggers a full Rebuild.
- `Rebuild(ideaDir, cfg) Cursor` — **total**, derives Phase purely from disk
  (D2/D3). The persisted Phase is never trusted over disk.

### readyPhase precedence (pure, disk-derived) — D3
1. `FINAL.md` present → PhaseImpl/PhaseReview (per IMPLEMENTATION.md presence/status).
2. else `consensus.md` present → PhaseConsensus (gate on `consensus.Status` triage).
3. else highest `round-NN/` dir → PhaseRound (completeness gate D4).
4. else PhaseRound, CurrentRound=1.

### Round-complete predicate (D4) — `roundComplete(ideaDir, round) (bool, error)`
1. expected set = `participants:` from `00-prompt.md`.
2. every `round-NN/<agent>.md` exists as a regular file.
3. each validates via `runner.ValidateRoundArtifact(path, agentID, ideaSlug, round)`;
   for round ≥ 2 the driver additionally requires `responding-to` present AND a
   `### @<other>` heading for every other participant.
4. terminal event for idea+round in the current run: `round.completed` accepts,
   `round.incomplete` rejects (authoritative block).
5. reconciliation: artifacts all valid (step 3) but no terminal event present →
   append a reconstructed `round.completed` (`completed==total==len(participants)`,
   marker `reconstructed:true`) to the current run, then accept. Malformed event log
   → return error → escalate. Never promote on file presence alone.

### Advance state machine (`internal/driver/driver.go`)
`Advance(ctx, *Cursor) (Action, error)` — ONE re-entrant tick:
```
c.Rebuild(ideaDir, cfg)               // disk wins
if !transportAllowsAuto(ideaDir) { return ActionSurfaceOnly, nil }
switch readyPhase(c):
  PhaseRound:
     if !roundComplete(c.CurrentRound):           return ActionAwait        // missing/incomplete
     if c.CurrentRound < 1 + cfg.CrossReviewRounds: RunRound(c.CurrentRound+1, Overwrite=false)
     else:                                          consensus.Draft(...)
  PhaseConsensus: triageGate(...)                  // D6
  PhaseFinal:     RunImplementation(...)
  PhaseImpl:      RunReviewRound(...) ... fix-up loop
  PhaseBlocked:   return ActionEscalated
c.Save(cursorPath)
```
Every branch is a no-op if its output already exists (Overwrite=false + fileExists
guards), so a duplicated tick / crash-restart cannot double-produce.

### Gate table
| Phase / state | Condition | Action |
|---|---|---|
| Round N | incomplete | await (loop deadline D12 → escalate) |
| Round N | complete, `CurrentRound < 1+cross_review_rounds` | `RunRound(N+1, Overwrite=false)` |
| Round N | complete, no more rounds | `consensus.Draft` |
| Consensus | TriageReady/Reserved, no FINAL.md | invoke FINAL drafter agent → verify non-scaffold (D7) → advance; else halt+escalate |
| Consensus | TriagePartial | invoke missing signers as real agents (`internal/signoffs`) |
| Consensus | TriageBlocked | `consensus.Reopen` + invalidate stale consensus.md/FINAL.md (`*.bak`) + `nextRound=latestRound+1` ≤ MaxRounds |
| Consensus | TriageMalformed | halt + escalate |
| Final | FINAL.md valid | `RunImplementation` |
| Impl done | IMPLEMENTATION.md present | `RunReviewRound` → fix-up loop |

### Non-scaffold FINAL.md check (D7)
frontmatter `status: final` valid; length > 250 bytes; `## Final plan /
specification` present with ≥3 non-placeholder lines; no unexpanded `<...>` / `[...]`
template variables.

### Crash-recovery contract
Disk files (atomic writes) are the durable truth; the cursor is a rebuildable cache;
events.jsonl is a derived signal (reconstructable per D4.5). Re-entry is idempotent:
`Rebuild` + the Overwrite=false guards make any tick safe to repeat. On
`consensus.Reopen` the driver invalidates stale `consensus.md`/`FINAL.md` so Rebuild
cannot misclassify a reopened round.

### Concurrency (D10)
Mandatory advisory `<runDir>/driver.lock` (PID) for the loop; acquire-fail = clean
stop. Contract: single-driver + idempotent re-entry, NOT multi-writer. The os.Stat
TOCTOU window in `runner.runAgent` stays open (no claim_lock) and is documented.

### Transport gating (D8)
`transportAllowsAuto`: read `transport:` from idea `00-prompt.md` if present else
`COOPERATION.md` global; auto-advance only if `--auto` AND effective transport ==
`local-dir`; re-evaluated every tick. Otherwise return surface-only.

### Circuit breakers
- MaxRounds (D11) default 4 cross-review rounds (escalate at round-05) → blocking
  `inbox/claude-to-user_deliberation-driver_max-rounds.md`, halt.
- Partial-round deadline (D12): fixed 30m per round → blocking inbox escalation,
  never spin.

### Test seam (D14)
Inject `RoundRunner` and `ConsensusOps` interfaces; unit tests use fakes that record
calls (no live agents). Table tests = the six in consensus D14.

## Slice 1 — implement and prove FIRST (D15)
Scope: ONLY PhaseRound (cursor + readyPhase + round-01→round-02 promotion via
RunRound), wired behind `parley run --auto --no-tui` on local-dir. NO consensus/
final/signoff/extraction work in slice 1.

File-by-file:
1. `internal/driver/cursor.go` — Cursor + Save/Load/Rebuild (PhaseRound paths only).
2. `internal/driver/driver.go` — `readyPhase`, `roundComplete`, `Advance` (PhaseRound
   branch), `RoundRunner` interface + production adapter calling
   `runner.RunRound(ctx, Options{Round:n, RoundLabel, Overwrite:false})`.
3. `internal/driver/loop.go` — `Run(ctx, ideaDir, cfg)` poll loop + advisory
   `driver.lock`; stop when readyPhase leaves PhaseRound or on await-deadline.
4. `internal/driver/driver_test.go` — table tests 1–4 + 6 from D14 with a fake
   RoundRunner (no live agents).
5. `internal/runaction/action.go` + `internal/runplan/runplan.go` — add
   `KindOpenNextRound` (visibility only; runplan emits it instead of jumping straight
   to consensus when more cross-review rounds are wanted).
6. `internal/app/app.go runTask` — `--auto && local-dir` → after RunRoundOne call
   `driver.Run`; non-auto unchanged. Read `cross_review_rounds` from 00-prompt.

Acceptance: a real `parley run --auto --no-tui` with all participants completes
round-01, opens round-02, writes valid `round-02/<agent>.md` for every participant
(D4 round≥2 validation), sets `00-prompt.md` status=`round-02`, emits the second
`round.started`/`round.completed`, and a repeated tick does NOT re-dispatch.

## Later slices (after slice 1 proven)
- S2: extract `internal/signoffs` from `internal/app/consensus_request_signoffs.go`
  (D9) — prerequisite for any driver→signoff/FINAL-draft wiring.
- S3: consensus gate (D6/D7) — Draft, FINAL-drafter agent launch + non-scaffold
  verify, request-signoffs via `internal/signoffs`, Reopen + stale invalidation.
- S4: PhaseFinal → RunImplementation; PhaseImpl → RunReviewRound + fix-up loop.
- S5: `runContinue --auto` executes the next action; full MaxRounds/deadline
  escalation wiring.

## Non-goals (D, reaffirmed)
No SQLite/daemon/worker-pool/scheduler; no claim_lock/heartbeat/DAG; reuse
run.json/events.jsonl/runstate/consensus; do not merge with §12 pipeline blocks;
discover.go/acp_specs.go untouched.

## Decisions reference
Full rationale + dismissed alternatives in `consensus.md` (D1–D15, trade-offs,
deferred, dismissed). round-01 / round-02 artifacts under `round-01/`, `round-02/`.
