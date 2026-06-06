---
idea: deliberation-driver
drafted-by: claude
date: 2026-06-05
rounds: [round-01, round-02]
---

Synthesis of round-01 (independent) and round-02 (cross-review) from claude, codex,
agy, hermes. All four converged with no remaining blockers. Decisions below are the
agreed design; `FINAL.md` expands them into the implementation spec.

## Agreed decisions

### D1 — Package shape
New package `internal/driver`: `cursor.go` (Cursor cache + atomic Save/Load/Rebuild),
`driver.go` (pure `readyPhase` + `Advance(ctx, *Cursor)` as a readable ordered
switch), `loop.go` (poll loop + mandatory advisory `driver.lock`). Linear ordered
switch only — **no DAG, no topological scheduler, no worker pool, no claim_lock, no
SQLite, no heartbeat/zombie-reaping**. Target ≤ ~300 LOC.

### D2 — Disk is authoritative; cursor is a rebuildable cache
`Cursor{Phase, CurrentRound, IdeaStatus, RoundsRun, MaxRounds, UpdatedAt}`. `Rebuild`
is total and derives phase purely from disk; the persisted `Phase` is never trusted
over disk. Atomic Save copies `internal/pipeline/run.go:82` (tmp+rename, same dir).
Location: `<runDir>/driver.json`; lock `<runDir>/driver.lock` (PID).

### D3 — readyPhase precedence (pure, disk-derived)
`FINAL.md` present → PhaseImpl/PhaseReview (per IMPLEMENTATION.md); else `consensus.md`
present → PhaseConsensus (triage gate); else highest `round-NN/` dir → PhaseRound
(completeness gate); else PhaseRound round=1.

### D4 — Round-complete gate (two-signal + reconciliation) — codex's predicate
A round N is complete iff:
1. expected participant set = `participants:` from `00-prompt.md`;
2. every `round-NN/<agent>.md` exists as a regular file;
3. each validates via `runner.ValidateRoundArtifact(path, agentID, ideaSlug, round)`
   (frontmatter agent/idea/round; round-1 sections Summary/Proposed approach/
   Concerns/Risks; round≥2 frontmatter + ≥1 H2). **Driver-side stricter for round≥2:**
   `responding-to` present AND a `### @<other>` heading for every other participant;
4. the current run's event stream has a terminal event for that idea+round:
   `round.completed` accepts, `round.incomplete` rejects (authoritative block);
5. **reconciliation:** if all artifacts are valid (step 3) but no terminal event
   exists (events.jsonl missing/truncated), append a reconstructed `round.completed`
   (`completed==total==len(participants)`, marker `reconstructed: true`) **only after
   full validation**, scoped to the current run, then proceed. A malformed event log →
   escalate, never guess. Re-emission is never triggered by file presence alone.

### D5 — Round gate (cross-review policy is explicit, not inferred)
`cross_review_rounds: N` in `00-prompt.md` frontmatter, **default 1**. Gate: if
`CurrentRound < 1 + cross_review_rounds` → `RunRound(next, Overwrite=false)` (skips
already-written artifacts → idempotent); else → `consensus.Draft`. `N=0` reproduces
straight-to-consensus but ONLY as a deliberate, explicit bypass — never the silent
default.

### D6 — Consensus triage gate (reuse `internal/consensus`)
- `TriageReady/Reserved` + no FINAL.md → **invoke the FINAL drafter agent** (00-prompt
  author / agreed drafter) to author FINAL.md content → verify non-scaffold (D7) →
  advance. Drafter failure/scaffold → halt + escalate.
- `TriagePartial` → invoke each missing signer as a **real agent** (via `internal/
  signoffs`) to author its own signoff; never fabricate ACCEPT. Non-auto: surface.
- `TriageBlocked` → `consensus.Reopen` + **invalidate stale `consensus.md` and draft
  `FINAL.md`** (rename `*.bak` or delete) so Rebuild cannot misclassify the reopened
  round + `nextRound = latestRound + 1`, bounded by MaxRounds (D11).
- `TriageMalformed` → halt + escalate; never auto-advance.

### D7 — Non-scaffold FINAL.md check (agy)
Valid frontmatter with `status: final`; file length > 250 bytes; `## Final plan /
specification` section present with ≥3 lines of non-whitespace/non-placeholder text;
no unexpanded template variables (`<slug>`, `<agent-id>`, bracketed placeholders).

### D8 — Transport gate (per tick)
Read `transport:` from the idea's `00-prompt.md` if present, else fall back to the
`COOPERATION.md` global; re-evaluate every tick. Auto-advance only if `--auto` AND
effective transport == `local-dir`. github-pr/gitlab-mr → surface the next action,
never auto-drive (humans own PR/MR labels).

### D9 — Real agent launches via an extracted shared service
Cross-review rounds, signoff requests, and FINAL drafting are all real agent launches.
**Prerequisite (before any driver→signoff wiring):** extract the reusable
request-signoffs logic from `internal/app/consensus_request_signoffs.go` into a new
`internal/signoffs` package (options, target selection, configured-agent discovery,
launch-mode validation, signoff prompt construction, `runSignoffAgent`, before/after
consensus validation, pending-handoff, structured results). `internal/app` keeps only
flag parsing, exit-code mapping, stdout/stderr presentation, usage text. `internal/
driver` and `internal/app` both depend on `internal/signoffs`; neither imports the
other.

### D10 — Concurrency contract
Mandatory advisory `driver.lock` (PID) for the loop; failure to acquire = clean stop.
Promise is **single-driver + idempotent re-entry**, explicitly NOT multi-writer
correctness. The `os.Stat` TOCTOU window in `runner.runAgent` stays open (no
claim_lock) and is documented; acceptable for slice 1.

### D11 — MaxRounds circuit breaker
Default 4 (counts cross-review rounds; round-01 is free; escalate at round-05). On
breach: write a blocking `inbox/<facilitator>-to-user_deliberation-driver_max-rounds.md`,
halt; user direction is quoted into the next round.

### D12 — Partial-round termination
Fixed per-round deadline (reuse 30m). On expiry: blocking inbox escalation rather than
spin. Process-liveness polling is deferred (too brittle for the first cut).

### D13 — Wiring
Add `KindOpenNextRound` to `internal/runaction` + `internal/runplan` for visibility
(`parley continue` surfaces "open round-02"); execution stays in `internal/driver`.
`internal/app/app.go runTask`: `--auto && local-dir` → after RunRoundOne hand to
`driver.Run`; non-auto keeps today's one-shot. `runContinue`: `--auto` → execute the
next action instead of only printing it.

### D14 — Test seam
Inject `RoundRunner` and `ConsensusOps` interfaces (function fields / small
interfaces); unit tests use fakes that record calls — **no live agents**. Table tests:
(1) round-01 complete → RunRound(2); (2) existing valid round-02 + round.completed →
no duplicate RunRound(2); (3) missing artifact or round.incomplete → no promotion
(await); (4) missing/corrupt cursor → Rebuild derives same phase; (5) TriageBlocked →
Reopen → RunRound(current+1) with MaxRounds enforced; (6) valid artifacts + missing
terminal event → re-emit round.completed; valid artifacts + round.incomplete → no
promotion.

### D15 — Minimal first slice (implement + prove FIRST)
ONLY PhaseRound: cursor + readyPhase + round-01→round-02 promotion via RunRound, wired
behind `parley run --auto --no-tui` on a local-dir workspace. **Acceptance:** a real
run with all participants completes round-01, opens round-02, produces valid
`round-02/<agent>.md` for every active participant (passing D4 round≥2 validation),
updates `00-prompt.md` status to `round-02`, emits the second `round.started`/
`round.completed` pair, and a repeated driver tick does NOT re-dispatch the round.

## Trade-offs accepted

- Single-driver only (no multi-writer); the advisory lock covers CLI double-start, not
  the Rebuild↔first-write race. Revisit only if the CLI model proves it needs it.
- Fixed 30m partial-round deadline instead of process-liveness polling.
- Event reconciliation scoped to the current run store; artifact truth from the idea dir.

## Deferred follow-ups

- Phases beyond round-promotion (consensus→final→impl→review/fix-up) — implemented
  after slice 1 is proven, same loop.
- Claim-level multi-writer dispatch; process-liveness polling.
- Unifying with §12 pipeline whole-idea blocks — explicitly OUT of scope here.

## Dismissed findings

- Files-only round-complete gate (claude round-01) → superseded by D4
  (two-signal + reconciliation).
- Transport read from idea-level frontmatter only (agy round-01) → refined to
  idea-level-first-with-global-fallback (codex), D8.
- `internal/driver` importing `internal/app` → rejected; extract `internal/signoffs` (D9).

## Signoffs

<!-- each participant appends its own ✅ ACCEPT / 🟡 ACCEPT-WITH-RESERVATIONS / ❌ BLOCK block -->

### claude — ✅ ACCEPT (2026-06-05)
This is the design I will implement. D4 (codex's validated two-signal gate), D6/D7
(agy's FINAL content gate + non-scaffold check), D9 (codex's `internal/signoffs`
extraction to fix the import direction), and D11/D12 (MaxRounds + deadline) close the
real correctness and protocol holes. Slice 1 (D15) is scoped to round promotion only
and is provable with a real run. No blockers.

### codex — ✅ ACCEPT (2026-06-05)
D4 captures my required validated two-signal round gate with reconciliation rather than file-only promotion. D9 keeps real-agent signoff/FINAL launch logic behind an extracted `internal/signoffs` seam, D14 makes the driver testable without live agents, and D15 keeps slice 1 limited to provable round-01→round-02 promotion.

### agy — ✅ ACCEPT (2026-06-05)
I confirm that the FINAL content gate (D6/D7), transport gate (D8), Reopen invalidation, and MaxRounds escalation (D11) are protocol-correct and resolve my concerns. This design ensures that the driver handles boundary conditions and escalations safely and predictably.

### hermes — ✅ ACCEPT (2026-06-05)
Design preserves donor-pattern fidelity: dispatch_once gated by D4's validated two-signal round gate + reconciliation, readyPhase recomputed purely from disk (D2/D3), and the single BLOCK back-edge via D6 Reopen with stale consensus/FINAL invalidation. No reservations.
