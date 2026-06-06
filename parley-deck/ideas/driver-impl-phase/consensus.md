---
idea: driver-impl-phase
drafted-by: claude
date: 2026-06-06
rounds: [round-01, round-02]
---

Synthesis of round-01 (independent) + round-02 (cross-review) from claude, codex,
agy, hermes. Full convergence, no blockers. Extends `internal/driver` (shipped
1.15.0) to auto-drive Parley Deck Phases 5–8 (implementation → review → fix-up)
after FINAL.md. The runner already provides every agent-launch building block
(`RunImplementation`/`RunReviewRound`/`RunFixup`/`RunReviewConsensus`) and the
machine-readable DONE signal (`outstanding_agreed_fixes` frontmatter), so this is
orchestration + gating, not new launch infrastructure.

## Agreed decisions

### D1 — `ImplOps` seam (driver-owned interface, app-injected adapter)
A `driver.ImplOps` interface (mirroring `ConsensusOps`): `Implement`,
`ImplementationStatus`, `RunChecks`, `OpenReviewRound(round)`,
`ReviewRoundComplete(round)`, `DraftReviewConsensus(round)`, `ReviewStatus`,
`RequestReviewSignoffs(missing)`, `Fixup(cycle)`, `Complete`. The driver core keeps
pure disk helpers (phase detection, artifact validation) and never imports
`internal/app`; the app injects the production adapter (reusing phase58
RunImplementation/RunReviewRound/RunFixup + `consensus` review-mode).

### D2 — `Rebuild` precedence: most-terminal-first
`valid IMPLEMENTATION.md status=complete → PhaseDone`; else
`review/consensus.md OR latest review/round-NN present → PhaseReview`; else
`valid IMPLEMENTATION.md present → PhaseImpl`; else
`valid FINAL.md present → PhaseFinal`; else `consensus.md → PhaseConsensus`; else
`PhaseRound`. (Most-terminal-first so a valid FINAL.md never hides later impl/review
artifacts.) Cycle/round numbers derived from disk, not persisted cursor fields.

### D3 — Safety: opt-in + preconditions + no-land boundary
- Code-writing phases (Implement, Fixup) require **`--auto` AND idea-level
  `auto_implement: true`** in `00-prompt.md` (default false). No new `--auto-implement`
  CLI flag — the permission lives as durable, auditable project state. A
  `--no-implement` CLI override forces the driver to stop at FINAL even when
  `auto_implement: true`.
- `auto_implement` is checked before BOTH Phase 5 and Phase 8 (re-entry must not run
  a fix-up just because review artifacts exist if the opt-in is absent).
- **Clean git working tree** precondition before Implement/Fixup; dirty → halt +
  blocking escalation (never overwrite uncommitted user work).
- **No-land boundary:** the driver NEVER merges/pushes/tags/releases; it stops at
  driver-written `IMPLEMENTATION.md status: complete` = "ready to merge".

### D4 — `RunChecks` build/test gate (injected, app-owned)
`RunChecks` is injected (production default `go test ./...`); the driver runs it
after implementation and after each fix-up, BEFORE `OpenReviewRound`. A failing
gate **escalates immediately** (no implicit fix-up loop — there is no reviewer-agreed
fix list for a raw compile/test failure). The driver does not parse test output
beyond pass/fail but includes the result text in the escalation.

### D5 — DONE detection (no prose scraping)
DONE = exactly ONE valid `review/consensus.md` for the current cycle with: signoff
triage Ready/Reserved, `blocked != true`, and `outstanding_agreed_fixes == 0`
(frontmatter the runner already emits + validates, phase58.go:266-298). `complete`
is written to IMPLEMENTATION.md **only by the driver** after that condition, so an
implementer cannot short-circuit review. A single clean consensus terminates (no
"two consecutive clean rounds" — the protocol precedent from slices 1–2);
non-convergence is bounded by MaxFixupCycles.

### D6 — Gate table (Advance)
| Phase | Condition | Action |
|---|---|---|
| PhaseFinal | valid FINAL.md, no IMPLEMENTATION.md, `auto_implement` & clean tree | `Implement` → PhaseImpl |
| PhaseFinal | `auto_implement` false / `--no-implement` | surface-only (stop at FINAL) |
| PhaseFinal | dirty tree | escalate |
| PhaseImpl | IMPLEMENTATION status not review-ready/malformed | await/escalate |
| PhaseImpl | review-ready, `RunChecks` fails | escalate |
| PhaseImpl | review-ready, `RunChecks` passes, review/round-NN absent | `OpenReviewRound(next)` |
| PhaseReview | latest review round incomplete | await |
| PhaseReview | round complete, no review consensus for cycle | `DraftReviewConsensus` |
| PhaseReview | consensus signoffs missing | `RequestReviewSignoffs`, re-check; still missing → escalate |
| PhaseReview | `blocked` / malformed / bad fix count | escalate |
| PhaseReview | Ready/Reserved & `outstanding_agreed_fixes==0` | `Complete` (driver writes status=complete) → PhaseDone |
| PhaseReview | `outstanding_agreed_fixes>0`, cycle<MaxFixupCycles | `Fixup` → next review round |
| PhaseReview | cycle ≥ MaxFixupCycles (default 3) | escalate |

### D7 — Idempotent re-entry + crash invalidation
Every action checks disk before launching (no re-implement/re-review of existing
artifacts). The fix-up cycle = highest `review/round-NN` + IMPLEMENTATION.md
`## Fix-up cycle N`. A partial/scaffold IMPLEMENTATION.md or review consensus must
re-validate before being trusted (the slice-2 S2-AF1 lesson) — `complete` /
`outstanding_agreed_fixes` are honored only when the artifact validates.

### D8 — review/consensus.md across cycles
Keep one canonical `review/consensus.md`, overwritten for each active review cycle
with a `review-cycle` frontmatter field; cross-cycle history lives in
IMPLEMENTATION.md. No separate `review/consensus-cycle-N.md` files (3:1; agy's
archive-to-`review/round-NN/consensus.md` is a deferred option if auditability needs
grow).

### D9 — `runContinue --auto`
`runContinue` with `--auto` EXECUTES the next `runplan` action (constructing the
driver and ticking it) instead of only printing it.

### D10 — Roles
Implementer = the FINAL drafter (or IMPLEMENTATION.md frontmatter), restricted to
idea participants; reviewers = the non-implementer participants; review-consensus
drafter = a real participant action.

## Trade-offs accepted
- Single clean review consensus terminates (no consecutive-clean); MaxFixupCycles is
  the non-convergence breaker.
- RunChecks failures escalate (no auto-repair loop for un-reviewed failures).
- The driver trusts agent-set `status: implemented`/`outstanding_agreed_fixes` only
  after disk validation + the RunChecks gate.

## Deferred follow-ups
- agy's per-round consensus archiving (`review/round-NN/consensus.md`).
- hermes's consecutive-clean confirmation round (recommendation, not adopted).
- Auto-drive for github-pr/gitlab-mr transports (the no-land boundary keeps this out).

## Dismissed findings
- A separate `--auto-implement` CLI flag — rejected in favor of the disk-auditable
  idea-level `auto_implement` field.
- Prose-scraping the `## Agreed fixes` section — rejected; `outstanding_agreed_fixes`
  frontmatter is the contract.

## Signoffs

<!-- each participant appends its own ✅ / 🟡 / ❌ block -->

### claude — ✅ ACCEPT (2026-06-06)
This is the design. It's orchestration over existing runner building blocks +
`outstanding_agreed_fixes`; the safety model (idea-level `auto_implement` opt-in +
clean-tree + RunChecks gate + no-land boundary + MaxFixupCycles) contains the real
new risk of auto-writing code, and `complete`-written-only-by-driver prevents
short-circuiting review. No blockers.

### codex — ✅ ACCEPT (2026-06-06)
I accept the design: the `ImplOps` seam keeps the driver core clean while allowing the app adapter to reuse existing Phase 5–8 runner operations. The most-terminal-first Rebuild order, idea-level `auto_implement` plus clean-tree plus RunChecks plus no-land safety model, DONE via `outstanding_agreed_fixes`, and driver-only `complete` write give the automation clear re-entry semantics without letting implementers bypass review.

### agy — ✅ ACCEPT (2026-06-06)
I accept the proposed design. The `ImplOps` seam cleanly separates orchestration from implementation, while the most-terminal-first Rebuild logic ensures accurate phase re-entry. The safety model (combining `auto_implement` opt-in, clean-tree checks, `RunChecks`, and the no-land boundary) coupled with the driver-only `complete` write and `outstanding_agreed_fixes` contract provides a robust, risk-mitigated path for auto-driving implementation and review.

### hermes — ✅ ACCEPT (2026-06-06)
The design is sound: the `ImplOps` seam cleanly isolates driver orchestration, most-terminal-first Rebuild ensures correct phase detection, and the `auto_implement`+clean-tree+RunChecks+no-land safety model plus `outstanding_agreed_fixes` DONE signal with driver-only `complete` write fully mitigates auto-implementation risk.
