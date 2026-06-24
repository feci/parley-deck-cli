---
idea: loop-engineering-research
drafter: claude-1
date: 2026-06-22
status: consensus-accepted
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Shared thesis (all four lenses agree)

**Parley Deck already implements the *inner* deliberation loop** (idea → rounds →
consensus → implement → review → fix-up) with durable on-disk artifacts, a
maker/checker split (Phase 5 implementer vs Phase 6 non-implementer reviewers),
worktrees, skills, and a CLI driver that auto-advances phases. Mapped onto loop
engineering's six blocks, Parley **covers** Worktrees, Skills, and (structurally)
Maker/Checker + Durable-State-per-idea; it **partially covers** Connectors
(`github-pr` is a write connector; no ingestion); and it **misses** the *outer*
loop: a **Trigger/heartbeat**, **explicit stopping conditions**, a **goal-done
gate**, and **cross-run durable goal state**.

The loop-engineering frame's real gift to us is **vocabulary + three concrete
sharpenings of things we already half-do**, plus one ambitious-but-gated outer loop.
The unifying safety principle (Osmani, echoed by every lens): an automated loop may
**discover and draft**, but a human/quorum must stay the gate before anything
**lands** — *"build it like someone who intends to stay the engineer, not just the
person who presses go."*

The deepest finding is not from the LE frame at all but from auditing our own engine
against it: **several gates we believe we have are inert.** `strict_gate` is protocol
prose the driver ignores; `RunChecks` is a no-op outside Go; the auto-loop closes on a
single agent's self-reported `outstanding_agreed_fixes == 0` with no objective
re-verification. Loop engineering gave us the lens that exposed these.

## Six-block mapping (consensus)

| # | LE block | Status | Seam / gap (exact) |
|---|---|---|---|
| 1 | Automations / Trigger | **MISSING** | `parley run` is human-triggered; driver auto-advances *within* an idea but nothing *opens* ideas. `pipeline watch` can auto-open remediation ideas but with `participants: []` (a non-solo Phase-0 violation — see LE-10). |
| 2 | Worktrees | **COVERED** | `parley-worktrees` addon. |
| 3 | Skills | **COVERED** | Parley *is* a SKILL.md; `COOPERATION.md` is externalized intent; `parley-tracker`. |
| 4 | Connectors | **PARTIAL** | `github-pr` writes PRs; no *ingestion* of issues/CI/logs; `parley-tracker` is mirror-only. §12 keeps side-effects gated (correct boundary). |
| 5 | Maker / Checker | **STRUCTURAL, gates inert** | Phase 5 vs Phase 6 by *agent id* (good; AF3 already moved the consensus drafter off the implementer). But: reviewer prompt is *confirmatory not refutational*; `strict_gate` unimplemented; model-diversity unenforced; `RunChecks` Go-only; close = one agent's self-report. |
| 6 | Durable State | **PER-IDEA yes / CROSS-RUN no** | Canonical artifacts + `runs/<id>/{run,events,driver}.json` are strong per-idea. No cross-run "what's next" ledger (CONTESTED — see LE-12). |

## Prioritized backlog (consensus)

Legend: **Adopt** = clear yes; **Adapt** = yes with stated modification; **Reject** =
no / out of scope. Each row names the file/§ it changes and a spin-off idea slug. This
idea is **design-only**; accepted rows become their own ideas.

### Tier 1 — verification honesty (cheap, unanimous, ship first)
| # | Recommendation | Target | LE principle | Effort | Call | Spin-off |
|---|---|---|---|---|---|---|
| LE-1 | **Refutation-default review posture** — `BuildReviewPrompt` instructs "assume it's wrong; try to construct a failing case / run each FINAL.md criterion; only report no-findings after stating what you tried"; `ValidateReviewArtifact` requires a `## Refutation attempts` section | `COOPERATION.md` §Phase 6 (≈L353-374) + `internal/runner/phase58.go` (`BuildReviewPrompt` L235-267, `ValidateReviewArtifact` L397-420) | Maker/Checker "defaults to refutation" | S | **Adopt** | `meta-protocol-change-refutation-default` |
| LE-2 | **Implement `strict_gate`** — the promised `ReadStrictGate` + `strict_gate_clean`/`closing_review_round` close fields were specified in `review-gate-honesty` and never built; under `strict_gate:true`, require a fresh full-scope zero-finding review round before `Complete()` | `internal/driver/impl.go` (`advanceReview` L182-197) + new `ReadStrictGate` (beside `ReadAutoImplement`) + `internal/app/driver_impl.go` `ReviewStatus`; `COOPERATION.md` §Phase 8 | Refutation re-verification pass | M | **Adopt** | `strict-gate-enforcement` |
| LE-3 | **Model-diversity guard** — warn (configurable: refuse auto-complete) when every reviewer/drafter shares the implementer's `Model`; no-op in today's 4-model roster, matters for 2-agent same-model decks | `internal/app/driver_impl.go` (`newDriverImplOps` L35-54, compares `agents.Discovery.Model`); one normative line in `COOPERATION.md` §Phase 6 | "ideally different models" | S | **Adopt** | folds into LE-1 |
| LE-4 | **Generalize `RunChecks`** beyond Go — read a `checks:` command from `00-prompt.md` frontmatter / `~/.parley [defaults]`; for code-writing `auto_implement` ideas with no `checks:` and no `go.mod`, **fail closed** instead of auto-passing. **Also tie the "artifact-wins" fix-up override to `RunChecks` passing**, not shape-validation alone (hermes #8): a fix-up that exits non-zero but wrote a valid-shaped IMPLEMENTATION.md must not count as success | `internal/app/driver_impl.go` (`RunChecks` L118-130); `internal/runner/phase58.go` (`RunFixup`/`ValidateFixupArtifact` L128-184); `COOPERATION.md` §Phase 4/5 document `checks:` | "verification is still on you"; conditional rigor | S-M | **Adopt** | `generalize-runchecks` |

### Tier 2 — stopping conditions / budgets (unanimous on need)
| # | Recommendation | Target | LE principle | Effort | Call | Spin-off |
|---|---|---|---|---|---|---|
| LE-5 | **Unified loop-budget contract** — one `LoopBudget` carrying max driver-steps, max rounds, max fixup-cycles, max wall-clock, max cost; **budget hit = escalate, never mark complete**. Consolidates today's fragmented caps (`loop.go:16` hard-coded 30m, `driver.go` `MaxRounds`/`MaxFixupCycles`, `pipeline auto` `maxCycles=3` + wave cap). Protocol states *that* limits exist; CLI/`[defaults]` state the numbers | Protocol `COOPERATION.md` §4/§9.0/§12; CLI `internal/driver/{driver,loop}.go`, `internal/runmanifest`, flags `--max-driver-steps/--max-wall-clock/--max-cost-usd/--max-rounds/--max-fixup-cycles`; `~/.parley [defaults.loop]` | Stopping conditions / hard limits | M | **Adopt** | `driver-loop-budgets` |
| LE-6 | **Best-effort cost telemetry** — emit `agent.usage` / `loop.budget` events where the CLI exposes tokens/cost; **strict mode treats unknown cost as a halt**; do not pass provider spend-caps that abort before an artifact lands | `internal/runner`, `internal/store/events.go`; `docs/agent-cli-mechanics.md` | Max-cost enforceability | M-L | **Adapt** | folds into LE-5 |

### Tier 3 — goal-done gate
| # | Recommendation | Target | LE principle | Effort | Call | Spin-off |
|---|---|---|---|---|---|---|
| LE-7 | **Objective goal-condition check before close** under `auto_implement`/`strict_gate`, derived from FINAL.md "Observable acceptance criteria," reusing the existing **advisory `consult.go`** as the "separate small model." Run **once, before `Complete()`** — not every tick (cost/fatigue) | `internal/driver/impl.go` (new step before `Complete()` L186-197) reusing `internal/runner/consult.go` (L51-63); `COOPERATION.md` §Phase 4 promotes acceptance criteria from advisory to gate *under auto mode* | `/goal` "separate model checks whether you're done" | M | **Adopt** (run-once form) | `goal-done-gate` |

### Tier 4 — outer loop + human brake (the paradigm piece — gated)
| # | Recommendation | Target | LE principle | Effort | Call | Spin-off |
|---|---|---|---|---|---|---|
| LE-8 | **Human-brake invariant** — codify that an automated loop may **discover + draft (Phase 0/1) only**; it must never push to quorum, implement, land/merge, finalize, modify the roster, or override/ bypass consensus without a recorded human or full-quorum gate | New `COOPERATION.md` § "Automation triggers / human brake"; references §2 roster, §3 consensus, §12 gates | "stay the engineer" / anti-cognitive-surrender | S | **Adopt** | `automation-human-brake` |
| LE-9 | **`parley loop tick`** — one-shot, scheduler-friendly (cron / GitHub Actions / MCP event), **not a resident daemon**: discover candidates (commits/CI/issues), dedupe, write trigger provenance into a **candidate** prompt, optionally call `parley run` within the LE-5 budget. Disabled by default; human-confirm mandatory before quorum | New `internal/loop` package; `docs/cli-reference.md`; protocol subsection under §12 | Automations / heartbeat | L | **Adapt** (opt-in, gated by LE-8) | `standing-loop-watch-mode` |
| LE-10 | **Fix `openRemediationIdea`** — it currently opens ideas with `participants: []`, violating the non-solo Phase-0 invariant; make watcher-created ideas a non-active *candidate* (outside `ideas/`) until they satisfy Phase-0 quorum | `internal/app/pipeline_cmd.go` (`openRemediationIdea`); `COOPERATION.md` §12.11 | Non-solo invariant; do before expanding triggers | S | **Adopt** | `remediation-idea-quorum-fix` |

### Tier 5 — HITL-fatigue guardrails
| # | Recommendation | Target | LE principle | Effort | Call | Spin-off |
|---|---|---|---|---|---|---|
| LE-11 | **HITL fatigue guardrails** — (a) batch/rate-limit driver-opened questions; (b) under `auto_implement`, do **not** auto-close on `TriageReserved` — escalate reservations to a human; (c) refuse auto-complete when `len(reviewers) < 2` | `internal/hitl/hitl.go`; `internal/driver/impl.go` (`advanceReview` L182, Reserved split by `AutoImplement`); `COOPERATION.md` §Phase 6 (extend "≥1 reviewer" to "≥2 for auto-complete") | HITL fatigue / cognitive surrender | S-M | **Adopt** | `hitl-fatigue-guardrails` |

### CONTESTED — cross-run durable goal state (LE-12)
| Lens | Position |
|---|---|
| claude-1, codex-1 | **Adopt**: a cross-run goal/backlog ledger — `parley-deck/STATE.md` (goal_id, status, done-condition, source, current idea/run, next-action, budget counters, last-checked) and/or `parley goals list\|add\|pause\|close\|check`; entries **append/update-disciplined** to avoid merge conflicts; goal *truth* lives in the deck, only *defaults* in `~/.parley`. |
| antigravity-1 | **Reject** as drafted: a shared mutable `STATE.md` adds context-window bloat + write-conflict risk; `parley-tracker` already maps external work — don't duplicate on-disk state. |
| hermes-1 | (out of lens; no position) |

**Consensus call:** the *need* (cross-run "what's next") is real and recurring (our
deferred follow-ups already scatter across memory/inbox/FINAL sections), but the
*home* is unresolved — new `STATE.md` vs reuse `parley-tracker` vs `runs/`-derived
view — and antigravity's write-conflict concern is legitimate. **Defer to its own idea
`durable-backlog-ledger`; do NOT build until the trigger/budget/brake scaffold
(LE-5/8) exists** (a backlog with no bounded loop to consume it is premature).

### REJECTED (unanimous)
- **Fully autonomous discover → implement → merge/deploy daemon.** Preserve §12 gates;
  production mutations stay non-bypassable; merge/release, roster edits, and consensus
  overrides are never automated.
- **Uniform refutation rigor on trivial / design-only ideas.** Keep **conditional
  rigor** — the default `outstanding_agreed_fixes == 0` close stays for low-risk ideas;
  refutation/strict-gate rigor scales with `auto_implement`/`strict_gate`.

### Framing (deferred)
Position Parley Deck as a **"loop-engineering substrate — the outer loop with a
human-gated consensus brake"** in README/marketing, borrowing Cherny/Osmani
vocabulary — but **only after** the safety work (LE-1/2/5/8) lands, so the claim is
true. Fold into the existing `readme-marketing-intro` idea.

## Recommended sequencing
1. **Tier 1 (LE-1..4)** — cheap, independently shippable, directly closes the
   false-green seams the audit found. Do first.
2. **LE-5 (budgets)** + **LE-10 (remediation quorum fix)** — the safety floor.
3. **LE-7 (goal-done) + LE-11 (HITL guardrails)** — close-decision integrity.
4. **LE-8 (human-brake §)** then **LE-9 (`parley loop tick`)** — the outer loop, only
   on top of the brake + budgets.
5. **LE-12** resolved in its own idea once 1-4 exist.

## Signoffs

**Outcome: consensus reached** — 4/4 accept (2 ACCEPT, 2 ACCEPT-WITH-RESERVATIONS,
0 BLOCK). hermes-1's reservation (fold #8 "artifact-wins" override into LE-4) was
incorporated above. codex-1 and antigravity-1's reservations are scoped deferments of
LE-12, already reflected in the CONTESTED section.

### Signoff: claude-1 — 2026-06-22
Status: ACCEPT
As drafter (synthesis lens) I confirm the backlog faithfully merges all four round-01
analyses, including the disagreement on LE-12 surfaced honestly rather than averaged
away. The verification-honesty cluster (LE-1..4) is the highest-leverage, lowest-risk
set and should ship first; the outer loop (LE-8/9) is correctly gated behind the
human-brake invariant and budgets. Design-only idea — accepted rows spin off their own
ideas.

### Signoff: codex-1 — 2026-06-22
Status: ACCEPT-WITH-RESERVATIONS
I accept the consensus sequencing: first close the false-green verification gaps, then
add budgets, the human brake, and only then expand the outer trigger loop. My
reservation is on LE-12: I still believe cross-run durable goal state is necessary, but
deferring its exact home to `durable-backlog-ledger` is acceptable because `STATE.md`
versus tracker/runs-derived state needs its own design. This should not be read as
rejecting the ledger need; it is a scoped deferment until LE-5 and LE-8 make the
consuming loop bounded and safe.
_(Facilitator note: codex-1 authored this signoff but its filesystem sandbox blocked
the write — "writing outside of the project ... rejected by user approval settings."
Transcribed verbatim from its run output per the protocol's tooling-failure exception;
source preserved at `tasks/bhc5t581e.output`.)_

### Signoff: hermes-1 — 2026-06-22
Status: ACCEPT-WITH-RESERVATIONS
My Tier-1 verification-honesty findings are faithfully represented: LE-1 captures my
refutation-default posture plus the machine-validated `## Refutation attempts` section
(#1/#9), LE-2 captures strict_gate enforcement exactly as I traced it (`ReadStrictGate`,
the `strict_gate_clean`/`closing_review_round` close fields, and a fresh full-scope
zero-finding review round before `Complete()`), and LE-3/LE-4 carry the model-diversity
guard and RunChecks generalization with the same targets and fail-closed semantics. The
REJECTED section correctly preserves my conditional-rigor position (#10). My one
reservation: recommendation #8 — the "artifact-wins" fix-up override
(`phase58.go:128-139`) counting a non-zero-exit fix-up as success on shape validation
alone — should fold into LE-4 once RunChecks is generalized. _(Incorporated into LE-4
above.)_

### Signoff: antigravity-1 — 2026-06-22
Status: ACCEPT
The consensus resolution to mark our LE-12 reject as contested and defer it to a
dedicated spin-off idea (`durable-backlog-ledger`) is acceptable. This deferral respects
our concerns regarding shared mutable files, context bloat, and write conflicts.
Gating its implementation on the completion of the safety budget and trigger scaffold
(LE-5/8) ensures we do not build a backlog ledger prematurely before its runtime
constraints exist.
