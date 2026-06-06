---
idea: deliberation-driver
review-cycle: 1
drafted-by: claude
date: 2026-06-05
reviewed-commit: f8c880d
---

Synthesis of Phase 6 review-round-01 from codex, agy, hermes (review/round-01/).
The implementer (claude) does not review its own code. Agreed fixes are applied in
fix-up cycle 1 on branch feature/deliberation-driver.

## Agreed fixes

### AF1 — Atomic driver.lock acquisition (codex MAJOR + agy MINOR)
`acquireLock` reads the PID file then `os.WriteFile`s it — two `parley run --auto`
processes can both observe no live lock and both become drivers, double-launching
agents (Overwrite=false does not close the `runner.runAgent` os.Stat TOCTOU that
D10 delegates to the lock). Fix: acquire with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)`
after stale-lock removal; release only if the file still holds this process'
ownership token. Add a concurrent-acquisition test proving only one holder wins.

### AF2 — Restore the D4 cross-review gate via the runner prompt (codex MAJOR)
The round≥2 gate was weakened to `responding-to` presence + `runner.ValidateRoundArtifact`
(any `## ` heading), so an artifact with `responding-to` + a generic `## Summary`
passes without per-agent engagement. Rather than weaken D4, strengthen the evidence:
update `runner.BuildRoundPrompt`'s required file shape to include a `### @<other
participant>` subheading under `## Responses to other participants` for each other
participant, then restore the driver's D4 check requiring a `### @<other>` heading
for every other active participant (keep `responding-to` as an additional metadata
check). Also (codex OQ) `terminalRoundEvent` must match BOTH `idea` and `round`, not
`round` alone. Update the real-run acceptance + fake test artifacts to carry the
`### @` evidence so tests would catch a regression.

### AF3 — Durable escalation on driver errors (codex MAJOR + agy MINOR)
A malformed `events.jsonl` makes `roundComplete` return an error that `loop.Run`
only prints to stderr — no blocking inbox note, unlike the deadline path. Runner
execution failures (`RunRound` error) likewise return `ActionEscalated` with no
durable artifact. COOPERATION.md §"Recovery And Partial Completion" requires
capturing such failures in an inbox note. Fix: generalize `escalateDeadline` into an
`escalate(reason, detail)` writer and route malformed-log errors AND runner failures
through it (idea, round, event-log path / agent error, recovery instructions).

### AF4 — Robust transport, re-evaluated per tick (agy MAJOR + agy NIT, D8)
`EffectiveTransport` only parses a backticked `**Transport:** \`local-dir\`` and
silently fails (→ surface-only) on the common no-backtick variation. And
`AutoLocalDir` is computed once in `app.go`, not per tick as D8 requires. Fix: reuse
the robust protocol transport reading (`protocol.ReadWorkspaceStatus(...).Transport`
or the `readTransport` regex which already handles optional backticks); compute the
effective transport from disk INSIDE `Advance` each tick (pass IdeaDir+Root, drop the
static `AutoLocalDir` bool) so a mid-run transport change is honored.

## Deferred follow-ups

- Full `readyPhase` precedence table for consensus/final/impl phases (hermes MINOR)
  — lands with later slices S2–S5; slice 1 correctly stops at the round boundary.
- Event-log indexing instead of full scan per `roundComplete` (hermes NIT) — fine at
  slice-1 volume.
- De-duplicate `setIdeaStatus`/`writeFileAtomic` vs `consensus`'s copies (hermes NIT).
- A dedicated `parley` resume/tick command vs re-running `parley run --auto` (agy OQ).
- Whether to document the `### @<other>` heading requirement in COOPERATION.md
  itself (agy OQ) — AF2 enforces it via the runner prompt for now.

## Dismissed findings

None — all findings are accepted as agreed fixes or scoped deferrals.

## Signoffs

<!-- each participant (incl. the implementer) appends its own ✅ / 🟡 / ❌ block -->

### claude — ✅ ACCEPT (2026-06-05)
Accept the synthesis. AF1 (atomic lock) and AF3 (durable escalation) are real
correctness/robustness gaps; AF2 is the right call — strengthen the runner prompt to
emit per-agent `### @` evidence and enforce D4 rather than weakening the gate; AF4
makes the transport gate robust and D8-compliant. Implementing all four in fix-up
cycle 1.

### codex — ✅ ACCEPT (2026-06-05)
AF1-AF4 are correctly applied: the driver lock is atomic and token-owned, the D4 cross-review gate is restored, driver errors create durable inbox escalations, and transport is re-read per tick through the protocol status path. `go build ./... && go vet ./... && go test ./...` is green with the requested cache settings.

### agy — ✅ ACCEPT (2026-06-05)
AF1-AF4 are correctly applied: lock acquisition is atomic, the D4 cross-review gate is restored with per-agent subheadings, driver/runner failures produce durable inbox escalations, and the transport is dynamically evaluated per tick. Build, vet, and all tests (including new concurrency/transport coverage) are green.

### hermes — ✅ ACCEPT (2026-06-05)
AF1–AF4 correctly applied (atomic lock, D4 cross-review gate+headings, durable escalate on errors, per-tick protocol transport); build/vet/tests green.

---

# Slice 2 review (cycle 2) — consensus gate

drafted-by: claude · date: 2026-06-05 · reviewed-commit: a83efa8 → fix-up pending.
Synthesis of review/round-03 (codex, agy, hermes) on slice 2 (the consensus gate).
All findings accepted as agreed fixes (applied in slice-2 fix-up cycle 1).

## Agreed fixes (slice 2)

### S2-AF1 — FINAL status committed only after non-scaffold validation (CRITICAL: all three)
`DraftFinal` called `consensus.Finalize`, which set idea `status: final` + a scaffold
FINAL.md BEFORE the drafter ran; a failed drafter stranded the idea at `final` with
scaffold content, and `Rebuild` (FINAL.md present → PhaseFinal) never revalidated.
Fix: the adapter no longer calls `consensus.Finalize` — the drafter agent authors
FINAL.md directly, and the DRIVER commits idea status to `final` only after
`finalScaffoldReason` passes. `Rebuild` treats only a VALID (non-scaffold) FINAL.md as
PhaseFinal; a scaffold stays in the consensus phase so the gate re-drafts it.

### S2-AF2 — Drafter must be an idea participant (MAJOR: agy)
`firstHeadlessAgent` picked the first installed headless agent regardless of the
roster. Fix: it is restricted to the idea's `participants` (COOPERATION.md §4/§6);
`newDriverConsensusOps` takes the participant list.

### S2-AF3 — Windows-safe process liveness (MAJOR: agy)
`processAlive` (signal 0) always failed on Windows → the driver lock always reclaimed
→ concurrent drivers. Fix: build-tagged `proclive_unix.go` (signal 0, EPERM→alive) +
`proclive_windows.go` (conservatively alive → refuses a second driver). Verified with
`GOOS=windows go build`.

### S2-AF4 — invalidateStale removes existing .bak, returns error, escalates (MAJOR: agy + codex)
`os.Rename(path, path+".bak")` discarded errors and fails on Windows / repeat BLOCK
cycles when `.bak` exists, leaving a stale FINAL.md. Fix: remove a pre-existing
`.bak` first, return the error, and `advanceConsensus` escalates rather than running
the next round with stale state.

### S2-AF5 — Reorder BLOCK: open round before Reopen/invalidate (MINOR: agy)
A `RunRound` failure after `Reopen`/invalidate forgot the BLOCK and re-drafted the
old consensus (loop). Fix: open the re-deliberation round FIRST; only on success run
`consensus.Reopen` + `invalidateStale` + commit the cursor.

### S2-AF6 — Ratify the ConsensusOps injection vs D9 extraction (MAJOR: hermes)
The driver injects `ConsensusOps` (app-side adapter) instead of extracting
`internal/signoffs`. codex/agy/hermes agree the import-direction guarantee is
preserved by injection. FINAL.md S2 updated to ratify the injection; the physical
extraction is deferred/optional.

## Deferred follow-ups (slice 2)
- Configurable facilitator-drafter in 00-prompt.md (agy/codex OQ) vs first
  participant.
- Human-in-the-loop signoffs on local-dir (agy OQ) — currently escalates under
  Partial; intended fallback.
- A failed consensus *content* enrich (vs FINAL) leaves a scaffold consensus that
  proceeds to signoffs — lower severity; revisit if it bites.

## Signoffs (slice 2)

### claude — ✅ ACCEPT (2026-06-05)
All six agreed fixes applied. S2-AF1 is the important one — status=final is now
committed only after the content validates, and `Rebuild` won't strand a scaffold at
PhaseFinal. Windows lock + invalidateStale are now correct; the drafter is a
participant; the BLOCK reorder prevents the loop. `go build/vet/test ./...` green;
`GOOS=windows go build` green. Live acceptance re-verified.

### codex — ✅ ACCEPT (2026-06-05)
S2-AF1..S2-AF6 are applied: final status is committed only after valid FINAL content, the drafter is participant-restricted, process liveness is build-tagged, stale invalidation escalates on failure, BLOCK reopen ordering is corrected, and FINAL.md ratifies `ConsensusOps` injection. `go build ./... && go vet ./... && go test ./...` plus `GOOS=windows GOARCH=amd64 go build ./...` are green with workspace caches.

### agy — ✅ ACCEPT (2026-06-05)
Verified that S2-AF1..S2-AF6 are applied correctly, with robust final status checks, restricted drafter selection, build-tagged liveness checks, clean stale file invalidation, corrected BLOCK reordering, and ratified dependency injection. Tests and Windows cross-compilation builds are clean and green.

### hermes — ✅ ACCEPT (2026-06-05)
S2-AF1..S2-AF6 all verified applied: final status only after valid FINAL, drafter restricted to participants, Windows-safe liveness, proper .bak handling+escalation, BLOCK reorder, ConsensusOps injection ratified. Builds (Unix+Windows), vet, tests green; live runs clean.
