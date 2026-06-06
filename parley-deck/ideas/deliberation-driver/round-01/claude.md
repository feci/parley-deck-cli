---
agent: claude
idea: deliberation-driver
round: 1
date: 2026-06-05
---

## Summary

I confirm the root cause at v1.14.1: `parley run` is a one-shot batch executor.
Verified directly — `runaction/action.go:9-14` has no `KindOpenNextRound`;
`runplan.go:113-122` emits `KindDraftConsensus` the moment a round is complete and
no consensus file exists (so round-01 jumps straight to consensus, no
cross-review); `runplan.go` `TriageBlocked → KindInspect` "before opening another
round" never opens it; `runner.go:292` already gives idempotent skip
(`os.Stat(outputPath)==nil && !Overwrite`); `runner.go:544 RunRound(N≥2)` exists
and reuses RunRoundOne machinery; `pipeline/run.go:82` is a clean atomic-save
template; `runstate.go:205 inferCurrentRound` exists to reuse. The enablers are
all present; nothing *drives* them. The `internal/driver` proposal is the right
shape. My round-01 pushes hardest on three things: **disk-authoritative cursor**,
**making cross-review an explicit policy (not inferred)**, and **the signoff gate
being a real agent-launch, not a state flip**.

## Proposed approach

**1. Cursor is a pure cache; disk is authoritative — Rebuild is total.**
`Cursor.Rebuild(runDir)` derives phase with a strict, ordered precedence and never
trusts a persisted `Phase`:
```
FINAL.md present                      → PhaseImpl / PhaseReview (per IMPLEMENTATION.md)
else consensus.md present             → PhaseConsensus (gate on consensus.Status triage)
else highest round-NN dir present     → PhaseRound (gate: is round-N complete?)
else                                  → PhaseRound, round=1 (initial)
```
A round-N is "complete" iff **all participant artifacts `round-NN/<id>.md` exist**.
Treat the on-disk *files* as the gate, not the `round.completed` event: events.jsonl
can be truncated on crash, files cannot (atomic writes). The driver should *emit*
`round.completed` during reconciliation if the files are present but the event is
missing — i.e. the event becomes a derived signal, not the source of truth. This
directly fixes "round.completed consumed by nobody": we stop depending on it.

`Cursor` therefore persists almost nothing load-bearing — `{Phase, CurrentRound,
IdeaStatus, RoundsRun}` are all re-derivable; only `MaxRounds` is config. The
durable cursor is an *optimization* so a warm loop skips a full disk scan each
tick; correctness comes entirely from Rebuild. Save via the `pipeline/run.go:82`
atomic tmp+rename verbatim.

**2. Cross-review is explicit policy, not content inference.**
The crux question — "round-01 → round-02 or → consensus?" — must be deterministic
and disk-derivable, not a fuzzy "did anyone disagree?" scan. Add
`cross_review_rounds: N` to `00-prompt.md` frontmatter (default **1**). The round
gate: `if CurrentRound < 1 + cross_review_rounds → open round-(CurrentRound+1) via
RunRound; else → draft consensus`. `N=0` reproduces today's straight-to-consensus
for trivial ideas; `N=1` (default) gives one independent + one cross-review round.
This keeps `readyPhase` a pure function of `(CurrentRound, config, artifacts)`.

**3. `Advance(ctx, *Cursor)` is one re-entrant, idempotent tick:**
```
c.Rebuild(runDir)            // disk wins; cursor reconciled
switch readyPhase(c):
  PhaseRound, incomplete  → RunRound(Overwrite=false) for missing participants only
  PhaseRound, complete & more-rounds-wanted → RunRound(next) (skips existing → no-op if re-entered)
  PhaseRound, complete & done → consensus.Draft
  PhaseConsensus → triage gate (below)
  PhaseFinal     → RunImplementation
  PhaseImpl done → RunReviewRound ... fix-up loop
c.Save()
```
Re-entry safety: every branch is a no-op if its output already exists (Overwrite=false,
fileExists guards), so a duplicated tick or crash-restart cannot double-produce.

**4. Consensus triage gate (reuse `consensus` pkg, do not reimplement):**
- `TriageReady/Reserved` + no FINAL.md → `Finalize` → FINAL.md (autonomous).
- `TriagePartial` → **invoke each missing signer as a real agent** to author its own
  signoff via `consensus.AppendSignoff` (the driver's "request-signoffs" action is an
  agent launch through the runner, NOT a synthesized ACCEPT). In non-auto, stop and
  surface the action.
- `TriageBlocked` → `consensus.Reopen` + `CurrentRound++` → new cross-review round
  seeded by the BLOCK counter-proposal, **bounded by MaxRounds** (default 4
  cross-review rounds); on breach, write `inbox/claude-to-user_deliberation-driver_
  max-rounds.md` (blocking) and stop.
- `TriageMalformed` → stop + escalate (never auto-advance malformed).

**5. Transport gate, read first.** `driver.Run` refuses unless `--auto` AND
transport==`local-dir`. For github-pr/gitlab-mr it returns a clear "auto-advance
disabled for <transport>" and falls back to surfacing the next action (humans drive
labels). Read transport from the run/workspace config before constructing the driver.

**6. Wiring + minimal slice.** Add `KindOpenNextRound` to `runaction` + `runplan`
(so read-only `parley continue` surfaces "open round-02" — fixes the planner half).
`runTask`: `--auto && local-dir` → after RunRoundOne, hand to `driver.Run`. First
slice implements ONLY PhaseRound (cursor + readyPhase + round-01→round-02 promotion)
and must prove a real run produces `round-02/<id>.md` for all participants +
`status=round-02` + a second `round.started/round.completed` pair.

## Concerns / open questions

- **Signoff launch is the hard coupling.** The driver must spawn agent processes for
  signoffs/cross-review rounds, so it is an agent-launcher, not just a state machine.
  I think it should call the existing `runner` launch path (same as RunRound) and own
  only the *gating*, not new launch code. Confirm: does `request-signoffs` have a
  per-agent launch path the driver can reuse, or must one be added? (`consensus_request_signoffs.go` exists — codex/agy please map it.)
- **Gate on files vs `round.completed` event** — I argue files; want codex/hermes to
  stress-test whether any consumer truly needs the event as the gate.
- **`cross_review_rounds` default** — I propose 1. Does agy see a protocol-correctness
  reason COOPERATION.md implies a different default (it says round-01 independent, then
  "cross-review rounds until objections resolved" — which argues for *at least* 1)?
- **MaxRounds counts what?** I propose counting cross-review rounds (so round-01 is
  free): MaxRounds=4 ⇒ escalate at round-05. Confirm off-by-one with the BLOCK path.
- **Where does the cursor live on disk?** Proposal: `<runDir>/driver.json` next to
  run.json; `driver.lock` (PID) beside it. Confirm no collision with runstate.

## Risks

- **Runaway --auto loop spawning real agents** is the scariest failure (cost + hangs).
  Mitigations: MaxRounds breaker, per-round process timeout (reuse 30m default),
  advisory `driver.lock`, and fully idempotent re-entry so a stuck tick re-derives
  rather than double-spawns. The loop must check `ctx.Done()` and the breaker every tick.
- **Auto-finalize on weak/fabricated consensus** ships bad design. Hard-blocked by the
  "signoffs are real agent invocations" rule + malformed/blocked never auto-advance.
- **Scope creep** into DAG/topological/claim-lock machinery. The graph is linear; a
  switch suffices. Reviewers should BLOCK any import of pipeline DAG or a lock beyond
  the advisory PID file. Target ≤ ~300 LOC for the package.
- **Partial-round deadlock**: if a participant never writes round-N, the gate blocks
  forever. Needs the protocol's deadline/quorum escape (or a per-round timeout that
  escalates), else --auto hangs silently. Flag for the design: how does a missing
  participant artifact terminate, not stall, the driver?
