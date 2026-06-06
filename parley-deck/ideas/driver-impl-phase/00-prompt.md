---
idea: driver-impl-phase
author: user
created: 2026-06-06
participants: [claude, codex, agy, hermes]
roles:
  claude: facilitator + driver state-machine extension
  codex: Go idioms, the impl/review/fix-up orchestration + tests
  agy: protocol-correctness (Phase 5-8 rules) + SAFETY of auto-writing code
  hermes: loop-termination / fix-up convergence fidelity
transport: local-dir
cross_review_rounds: 1
status: final
---

## Problem / idea

`internal/driver` (shipped in parley-deck-cli 1.15.0) auto-drives a deliberation
through round → cross-review → consensus → signoffs → **FINAL.md** under
`parley run --auto --no-tui` (local-dir). It then STOPS at PhaseFinal (surface-only).
This idea implements **FINAL.md S4–S5**: extend the driver to auto-drive the
**implementation and review phases** (Parley Deck Phases 5–8) after FINAL.md, and
make `parley continue --auto` execute the next action instead of only printing it.

### The building blocks already exist (like RunRound did for slice 1)
`internal/runner/phase58.go`: `RunImplementation` (Phase 5 — one implementer writes
`IMPLEMENTATION.md` + code on a branch per FINAL.md), `RunReviewRound` (Phase 6 —
non-implementer reviewers write `review/round-NN/<agent>.md`), `RunFixup` (Phase 8 —
re-invokes the implementer to apply agreed fixes + update IMPLEMENTATION.md). The
`consensus` package already has review-mode Draft/Status/AppendSignoff. So this is
mostly ORCHESTRATION + gating, not new agent-launch infrastructure.

## Proposed direction (a STARTING proposal — challenge it in round-01)

Extend `Advance` with PhaseImpl/PhaseReview branches, fed by an injected `ImplOps`
interface (mirroring the slice-2 `ConsensusOps` seam, so the driver core stays
testable and never imports `internal/app`):
- **PhaseFinal** (valid FINAL.md, no IMPLEMENTATION.md) → `ImplOps.Implement` (one
  implementer agent) → PhaseImpl.
- **PhaseImpl** (IMPLEMENTATION.md present) → when implementation is "ready"
  (status frontmatter) → `ImplOps.OpenReviewRound` (non-implementer reviewers) →
  PhaseReview.
- **PhaseReview** (review/round-NN present + complete) → draft `review/consensus.md`
  → request review signoffs (real agents) → triage:
    - zero agreed fixes / all ACCEPT → **DONE** (status complete) → STOP.
    - agreed fixes present → `ImplOps.Fixup` (re-invoke implementer) → new review
      round, bounded by **MaxFixupCycles** (default 3).
- `runContinue` with `--auto` EXECUTES the next `runplan` action instead of printing.

## Round-01 focus questions (answer independently)

1. **SAFETY — auto-driving an agent that writes CODE.** This is materially riskier
   than writing markdown artifacts (slices 1–2). Is `--auto` + local-dir a
   sufficient opt-in, or does auto-implementation need an additional explicit gate
   (e.g. `--auto-implement`, or an idea-level opt-in)? The proposal is that the
   driver auto-drives the AUTHORING (implement on a branch → review → fix-up) but
   NEVER merges/pushes/releases — it stops at "ready to merge" and a human takes
   over. Is that the right safety boundary?
2. **Phase + done detection (disk-derived, like the cursor).** How does the driver
   tell "implementation in progress" from "ready for review"? IMPLEMENTATION.md
   `status:` frontmatter (implemented / fix-up-cycle-N / complete)? How does it
   detect "review complete" and "zero agreed fixes" — a machine-readable signal in
   `review/consensus.md` (an `## Agreed fixes` count, or review-mode triage)?
3. **Loop termination / convergence.** The fix-up loop (review → fixup → re-review)
   must terminate. MaxFixupCycles breaker + escalation. How to avoid re-running a
   review round whose artifacts already exist (idempotent re-entry)?
4. **Implementer identity.** Who implements — the FINAL drafter (protocol default)
   or a designated agent? Who reviews (non-implementer participants)? Who drafts
   review/consensus?
5. **`ImplOps` seam shape** + the disk-derived PhaseImpl/PhaseReview detection in
   `Rebuild`.

## Constraints (non-negotiable)
- Reuse `RunImplementation`/`RunReviewRound`/`RunFixup` + `consensus` review-mode;
  do NOT reimplement agent launching. Inject an `ImplOps` interface (app-side
  adapter) — the driver MUST NOT import `internal/app`.
- Disk is authoritative; every gate is a pure function of on-disk artifacts;
  re-entry idempotent. Auto-drive ONLY under `--auto` AND local-dir transport.
- The driver NEVER merges/pushes/releases — it stops at "ready to merge".
- Real agent invocations for implement/review/signoffs; never fabricate.
- MaxFixupCycles circuit breaker (default 3) → escalate, never spin.
- English-only; one file per agent per round; append-only signoffs.

## Non-goals
- No merge/push/release automation. No new storage. No change to how agents launch
  (discover.go untouched). Not unifying with the §12 pipeline.

## Deliverables (in order)
1. FINAL.md: the impl/review auto-drive design (ImplOps seam, Rebuild phase
   detection, gate table for PhaseFinal→Impl→Review→fix-up, done/termination
   contract, safety boundary).
2. Implement it (behind `parley run --auto`/`continue --auto`, local-dir), with
   unit tests using a fake ImplOps (no live agents), then a live acceptance that
   drives a tiny idea FINAL → implementation → review → complete.
