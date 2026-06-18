---
agent: claude
idea: meta-protocol-change-fusion-execplans
round: 1
date: 2026-06-18
---

## Summary

parley-deck is already a stronger multi-agent deliberation than Fusion's *ephemeral*
panel→judge→answer (we keep durable per-agent artifacts, real disagreement, and a
human gate), so most of Fusion is "already owned." The genuinely *new* idea from
Fusion is one lens: **blind spots none addressed**. The bigger prize is on the
**ExecPlans** side — turning `FINAL.md` / `IMPLEMENTATION.md` into self-contained
*living* documents directly serves two current pain points: (1) headless / auto-drive
**resumption across invocations**, and (2) starving `parley retro` (§13) of structured
evidence. My top-3: ExecPlan-ify IMPLEMENTATION.md; add a blind-spots lens to
consensus; add behavior-focused acceptance + idempotence/recovery to FINAL.md.

## Per-concept verdicts

**Fusion**

1. **Judge "compare, don't merge" + structured analysis** — *Mostly already owned.*
   Our `consensus.md` is the judge; cross-review rounds already surface
   contradictions; signoffs already encode agreement. **Adapt the one missing lens:**
   add a short **"Blind spots / unaddressed"** subsection to `consensus.md` (and the
   review consensus) — the failure mode where *all* agents share the same gap is
   exactly what our independent-then-converge flow can't catch today. Cheap, real.
2. **Confidence-by-breadth-of-agreement** — *Adapt, minimally.* Signoffs are binary
   (✅/❌). A non-binding "consensus strength" note (e.g. unanimous vs. majority-with-
   reservations) could feed §13 retro and the driver, but turning it into a new
   *gate* risks eroding the all-participants-must-accept invariant. Note only, no gate.
3. **Synthesis-as-distinct-value (helps even without diversity)** — *Already owned,
   reframe.* The consensus *drafter* is our synthesizer. The useful takeaway: treat
   the consensus draft as a first-class analytical artifact, not just a signoff
   stapler. No structural change needed.
4. **Negative-weight / "confidently wrong"** — *Adapt into §13.* Add a "confidently
   wrong" failure-signal to retro mining (an agent that BLOCKED/asserted strongly and
   was later overturned). Small, fits the existing retro signal model. Not a review
   severity change.
5. Recursion-depth / cost / preset / web-tool mechanics — *Reject (N/A):* these are
   API plumbing irrelevant to a file-based protocol.

**ExecPlans**

6. **Self-contained living FINAL.md / IMPLEMENTATION.md** — *Adopt (highest leverage).*
   IMPLEMENTATION.md already logs deviations + fix-up cycles, but it is not guaranteed
   *resumable-from-alone*. Adopting the living sections — **Progress** (timestamped
   checklist), **Decision Log** (decision/rationale/date·author), **Surprises &
   Discoveries**, **Outcomes & Retrospective** — makes a fresh headless codex/agy/
   hermes (or the auto-drive **driver**) able to resume implementation from the
   artifact alone. Two-birds: it *also* hands §13 `parley retro` far richer structured
   signals than the thin density metrics we mine now.
7. **Behavior-focused acceptance criteria** — *Adopt (cheap).* Require `FINAL.md` to
   state observable acceptance ("after X, Y is true") that Phase 6 review and the
   driver can check. Reduces the "compiles but does nothing" failure ExecPlans warns
   about; tightens our review.
8. **Idempotence & Recovery section** — *Adopt (small, targeted).* For `auto_implement`
   ideas especially, an explicit idempotency/recovery section hardens the gated
   auto-drive (clean-tree, no-land) and gives the driver a defined rollback path.
9. **Self-containment / "define every term" / narrative-prose maximalism** — *Adapt
   selectively.* Self-containment for FINAL/IMPLEMENTATION: yes. The "prose over
   lists, avoid tables" maximalism partly **conflicts** with our structured frontmatter
   + severity tables (CRITICAL/MAJOR/MINOR/NIT) — keep our structure; borrow only the
   "each milestone reads as goal→work→result→proof" discipline.

## Prioritized recommendation (top-3)

1. **ExecPlan-ify `IMPLEMENTATION.md`** (living sections). Highest leverage: fixes
   real resumption pain for headless agents + the driver, and feeds §13 retro.
2. **"Blind spots / unaddressed" lens** in `consensus.md` + `review/consensus.md`.
   Fusion's one genuinely-new contribution; near-zero cost.
3. **Behavior-focused acceptance + Idempotence/Recovery in `FINAL.md`**, applied
   especially to `auto_implement` ideas, to harden auto-drive.

## Risks / what NOT to adopt

- **Ceremony tax.** parley-deck already pays a heavy multi-agent overhead. Make new
  sections *lightweight and conditional* (full living-doc rigor only for complex /
  `auto_implement` ideas; trivial ideas keep today's lean IMPLEMENTATION.md).
- **Don't** import Fusion's ephemeral-panel or confidence-gating as a *gate* — it
  conflicts with durable artifacts + all-participants-accept.
- **Don't** adopt ExecPlans' anti-list/anti-table stance wholesale — our structured
  artifacts are a feature, not a smell.
- Net: these are mostly *artifact-shape* refinements (FINAL/IMPLEMENTATION/consensus),
  not new phases — which is the cheap, safe place to borrow.
