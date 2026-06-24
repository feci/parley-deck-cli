# round-01 — claude-1 (synthesis lens)

**Thesis:** Parley Deck is already ~70% a loop-engineering system, but it implements
the *inner* deliberation loop (idea → rounds → consensus → implement → review) with a
**human-gated consensus brake**. The three things loop engineering has that we lack are
all in the **outer** loop: a *trigger* (heartbeat), *explicit stopping conditions*, and
a *cross-run durable goal ledger*. The single highest-leverage thing we already half-do
and should make explicit is **refutation-default, different-model verification**.

I deliberately did not read the other round-01 files before writing this.

## 1. Mapping the six blocks onto Parley's surface

| # | LE block | Parley status | Seam / gap |
|---|---|---|---|
| 1 | **Automations / Trigger (heartbeat)** | **MISSING** | The CLI `internal/driver` auto-advances *within* an open idea and `parley run` auto-drives by default — but **nothing opens ideas**. Parley is human-triggered per idea. The harness has cron primitives; they are not wired to "discover work → draft a 00-prompt". |
| 2 | **Worktrees** | **COVERED** | `parley-worktrees` addon (sibling `../worktrees/`, lock manifest). Maps 1:1. |
| 3 | **Skills (externalize intent)** | **COVERED** | Parley Deck *is* a SKILL.md; `COOPERATION.md` is the externalized intent; `parley-tracker` externalizes ticket intent. Maps 1:1. |
| 4 | **Plugins / Connectors (side-effect authority)** | **PARTIAL** | `github-pr` transport is a PR connector (write side). But there is **no ingestion connector** — we cannot read issues/CI-failures/Slack to *discover* work, and `parley-tracker` is mirror-only (no automated write-back to the tracker). |
| 5 | **Maker / Checker sub-agents** | **COVERED but implicit** | Phase 5 implementer (default = FINAL drafter) vs Phase 6 reviewers (every *non-implementer*); Phase 7/8 fix-up loop until zero agreed fixes. We do NOT mandate (a) verifier model ≠ implementer model, nor (b) a "default-to-refute" stance. |
| 6 | **Durable State (memory on disk)** | **COVERED per-idea / MISSING cross-run** | Canonical artifacts + append-only signoffs + inbox are textbook durable state *inside one idea*. There is **no cross-idea/cross-run "what's next" ledger** — discovered work and deferred follow-ups scatter across memory files, inboxes, and FINAL.md "deferred" sections. |

**Read of the map:** our gaps are concentrated in blocks 1, 4 (ingestion half), 6
(cross-run half), and the *explicitness* of block 5. That is the backlog.

## 2. What loop engineering teaches us that we are missing

### 2a. Stopping conditions are not a feature — they are a missing invariant (CLI)
LE mandates max-iterations / max-cost / max-duration + a goal-checker, "whichever
triggers first halts and escalates." Our driver auto-advances Phases 5–8 with the only
hard gates being `auto_implement` opt-in + clean-tree + RunChecks + no-land. It has **no
explicit iteration ceiling, no wall-clock ceiling, and no cost ceiling** on the fix-up
loop (Phase 8 "repeat review until zero agreed fixes" is an *unbounded* loop on paper).
A flapping reviewer pair could in principle cycle. This is the clearest, lowest-risk
borrow.

### 2b. The outer-loop trigger is the real paradigm shift — and our biggest risk surface
"My job is to write loops" = the human stops being the per-idea trigger. For Parley
that means a `parley loop` watch mode: on a cron/event, read new commits / CI failures /
open issues, and **draft** a 00-prompt for a new idea. But Parley's whole value is the
*deliberate human-and-quorum gate*. So the trigger may **discover and draft**; it must
**not** push an idea to quorum, land, or finalize without a human confirm. That single
rule is what separates "loop engineering that keeps you the engineer" (Osmani) from
cognitive surrender.

### 2c. Refutation-default verification — we already paid for this lesson
The Phase-6 false-green (a build's self-reported "tests pass" that multi-agent review
caught) IS the "model too nice grading its own homework" failure. We solved it
*structurally* (non-implementer reviewers) but never wrote down the *stance*. LE gives
us the words: verifier **defaults to refutation** and **uses a different model**. We
should make both explicit in §6.

## 3. Where each change belongs (protocol invariant vs CLI vs ~/.parley)

- **Stance/semantics → `COOPERATION.md`** (refute-default; human-brake-before-land;
  the *existence* of stopping conditions as an invariant).
- **Mechanism/numbers → `parley-deck-cli` driver + `~/.parley [defaults]`** (the actual
  max-iters / wall-clock / budget values; the watch-mode wiring).
- **Rule of thumb:** the protocol says *that* a loop must terminate and *that* a
  verifier refutes; the CLI/`[defaults]` says *with what numbers* and *how*.

## 4. Prioritized recommendations

| # | Recommendation | Protocol/CLI | LE principle | Effort | Risk | Call | Spin-off slug |
|---|---|---|---|---|---|---|---|
| 1 | **§6: make verification "default-to-refute" + require verifier model ≠ implementer model** (soft invariant; record the exception when only one model is available) | Protocol | Maker/Checker | S | Low | **Adopt** | `meta-protocol-change-refute-default-verification` |
| 2 | **Driver stopping conditions**: explicit max-phase-iterations, max wall-clock, and (where measurable) token/cost budget on the auto-drive + Phase-8 fix-up loop; on breach → halt + `inbox/<driver>-to-user`. Numbers in `~/.parley [defaults].limits`, per-idea override | CLI + `[defaults]` | Decide / hard limits | M | Low | **Adopt** | `driver-stopping-conditions` |
| 3 | **§ "Human brake" invariant**: an automated outer loop may DISCOVER + DRAFT a 00-prompt, but must never push to quorum, land/merge, or finalize without a recorded human or full-quorum gate | Protocol | "stay the engineer" / anti-cognitive-surrender | S | Low | **Adopt** | folds into #5's idea |
| 4 | **Cross-run goal ledger**: a `parley-deck/STATE.md` (or `backlog/`) that consolidates discovered work + deferred follow-ups + "what's next", read at session start and each loop cycle; CLI `parley backlog` to list/add | Protocol-light + CLI | Durable State (cross-run) | M | Low | **Adopt** | `durable-backlog-ledger` |
| 5 | **`parley loop` watch mode**: cron/event trigger reads commits/CI/issues → drafts an idea 00-prompt for human confirm (gated by #3) | CLI | Automations / heartbeat | L | **Med-High** | **Adapt** (opt-in, human-confirm mandatory) | `standing-loop-watch-mode` |
| 6 | **Ingestion connector**: let a transport/connector *read* issues/CI to feed #5 (write-back to tracker via `parley-tracker` later) | CLI | Connectors (ingest half) | L | Med | **Adapt** (later; depends on #5) | `ingestion-connectors` |
| 7 | **Framing**: position Parley Deck in README as a *loop-engineering substrate — the outer loop with a human-gated consensus brake*; borrow Cherny/Osmani vocabulary | Docs | Naming | S | Low | **Adopt** (fold into `readme-marketing-intro`) | — |

## 5. What Parley already mitigates (don't over-engineer)
- **Verifier mis-grading** → multi-agent *quorum* review structurally beats single
  self-grade. Recommendation #1 only makes the *stance* explicit; we don't need a new
  mechanism.
- **Runaway cost** → today bounded by the human being the trigger. The moment we add #5
  (auto-trigger) this protection vanishes — which is exactly why #2 (stopping conditions)
  must land *before or with* #5, not after.
- **New risk introduced by automation: HITL approval fatigue.** If #5 drafts many ideas,
  the human confirm becomes a mechanical click. Mitigate: batch/rate-limit drafts, and
  make the confirm show a 3-line "why this idea" summary, not just a yes/no.

## 6. Sequencing
#1 and #2 are cheap, high-leverage, and independently shippable — do them first. #3+#4
are the *safety scaffold* that must exist before #5 is safe. #5/#6 are the ambitious
outer-loop; only attempt once the brake (#3) and the ledger (#4) and the limits (#2)
exist. Reject any version of #5 that can reach quorum/land without a human gate.
