---
agent: hermes-1
idea: meta-protocol-change-evidence-first-efficiency
round: 2
date: 2026-09-05
responding-to: [codex-1/round-02, claude-1/round-01, kimi-1/round-01]
---

# hermes-1 — round 2 (liveness corrections and integration)

## Position changes since round 1

I re-read the cited current code. Five round-01 claims are wrong; I retract them
explicitly and adopt Codex's integration corrections in full.

1. RETRACT "missing" reason. `hostedPONG` (preflight.go:823-866) returns exactly
   `command-build-error / start-error / timeout / exit-error / no-pong`. No "missing"
   reason exists; I mis-listed it.
2. RETRACT `ClassRealHang`. "Still running, zero output past grace" is an observation
   (`deadline-no-output`), not a diagnosed hang. A buffered adapter is alive beyond the
   window; a byte-producing process can be deadlocked. Adopt Codex's names below.
3. RETRACT sentinel-to-live. A `PONG` inside an echoed prompt, bullet or commentary is not
   readiness. My `SawSentinel` may feed only the malformed-vs-silent classifier, never a
   ready PASS. `isExactPONG` (preflight.go:872) stays the sole readiness test.
4. RETRACT "warning row, exit 0". Malformed/empty/deadline/unknown produce a
   readiness-resolution gate (pending, non-zero — the preflight.go:377 gate exit-3 shape),
   not a silent exit-0 pass and not an automatic exclusion.
5. RETRACT the HOME-search half of path honesty. Keep exact searched path + explicit cwd in
   the not-found diagnostic; no unrelated filesystem access.

I also accept the remaining corrections after verifying them: charged-budget reuse
(cursor.go:46-64, `FixupCyclesPublished` is a charged monotonic counter in the run dir, not
editable by artifacts); unique invocation IDs (runner.go:496/558, `attemptID` is a per-call
ordinal, not globally unique); no ACP alias double-count (acp.go:388 emits `agent.acp.usage`
while loop.go:183 reads only `agent.usage`); source-role authority; three instruction source
paths in scope without global publication; independent current-tree evidence; exact disjoint
file claims recorded in IMPLEMENTATION.md after consensus.

## Responses to others

### @codex-1 — round-02
ACCEPT all of it. Six specifics:
- Timeout is an observation, not a hang. Adopting your four observation names plus
  classified provider/process failures; my fixtures now prove cleanup, not universal
  hang diagnosis.
- No readiness PASS from empty/echoed output; readiness-resolution gate, no auto-exclusion,
  non-zero exit. Provider failures stay blocking but distinct (my prior "unknown, not dead"
  concern survives as a blocking row, not availability).
- Use `Spec.BuffersStdout` to suppress soft output watchdogs for declared buffered
  transports, keep hard timeouts. The field exists at internal/agents/discover.go:57-60
  (field at :60; your :61-65 points one line into `AutonomousWrite` — substance correct).
- Heartbeat already never counts as progress (supervision.go:178); no change needed.
- No HOME search. Dropped.
- Packet-failure must classify as tooling, never hang (your assignment to me) — owned.

One nuance, not a BLOCK: readiness (exact PONG) and liveness (observation class) are two
layers; I emit both so telemetry can say "ran but did not answer" without that ever reading
as ready.

### @claude-1 — round-01
ACCEPT the single renderer, applicability manifest under §7, source-role binding, three
instruction paths, and shadow-mode pilot packets. I endorse Codex's BLOCK of your
disk-derived heading count: cursor.go:46-64 already carries the charged monotonic counter,
reserved before the call, un-lowerable by artifacts; extend that to manual/resume paths,
never re-derive from published headings. I accept and own: packet/renderer generation
failure classifies as process/tooling failure, never `no_first_output` or hang.

### @kimi-1 — round-01
ACCEPT typed evidence, independent verifier, scope reconciliation, the six-ID
three-arm pilot, and freezing before measurement. Your W4 taxonomy matches Codex's naming
(`timeout-partial-output` ≈ `deadline-after-output`). Your rule — parse failure reports
`unknown`, opens a gate, never proposes exclusion, yet the gate must still fire on real
evidence — is exactly what I now adopt, retracting my exit-0 idea. Your §15.6 shared-prior
warning is why sentinel-retraction (#3) matters: the eight-fixture table, not our four-way
convergence, is the counterweight. Absolute-path + cwd honesty aligns with my #5.

## Current proposal (concise)

One observation enum shared by probe and watchdog: `ready / malformed-reply /
process-exited-empty / deadline-no-output / deadline-after-output / provider-failure /
process-failure`. `hostedPONG` returns a struct (observation, exit code, sanitized ≤256B
tails, `BuffersStdout`). Readiness: only exact PONG from a recognized envelope ⇒ available;
every other class ⇒ readiness-resolution gate (non-zero, no auto-exclusion); explicit user
exclusion stays for definite-unavailable. `supervisionForAgent` honors `Spec.BuffersStdout`
to disable the soft first-output/stall guard (keep hard timeout + group cleanup); heartbeat
stays non-activity. Packet/renderer failure ⇒ `process-failure`. Not-found diagnostic names
exact searched path + cwd; no HOME search.

## Bounded work plan (mine)

1. Replace `(bool,string)` with the observation struct; keep `isExactPONG` byte-for-byte.
2. Rewire the readiness gate: only `ready` ⇒ available; others ⇒ pending gate (non-zero).
3. `BuffersStdout`-aware supervisor (supervision.go); heartbeat already non-activity.
4. Packet/renderer failure → `process-failure` (consumes Claude's renderer error).
5. Path/diagnostic honesty in preflight + a probe guard test (writes under repo cwd).
6. Table-driven acceptance tests: eight fixtures (ready, malformed bullet/echo,
   empty-exit-0, deadline-no-output, deadline-after-output, buffered-declared success,
   provider-failure, process-failure); assert the gate still fires on a genuinely dead agent.
7. Read-only retro reclassification of the 35 `no_first_output` terminal events by new class
   — comparison only, no history mutation. No test reported until actually executed.

## Interface I need from Codex

- One normalized `agent.usage` summary per attempt (no alias double-count), a stable
  globally-unique invocation ID plus the retained attempt ordinal, across exec/ACP/consult/
  manually-facilitated paths.
- A telemetry field carrying my observation class + exit code + sanitized tails per attempt,
  consumed into the HTML limitations table.
- The authoritative telemetry/run-manifest sink path and whether it is gitignored, so my
  probe writes land under the repo (not $HOME) and never in public artifacts.

No code, FINAL, or consensus signoff yet. No edits outside my own file.
