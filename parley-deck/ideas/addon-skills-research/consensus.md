---
idea: addon-skills-research
drafted-by: claude-1
date: 2026-06-22
status: awaiting-signoffs
---

## Agreed decisions

Four round-01 analyses converged with no substantive blockers. Both skills follow
the **core-stays-thin, skill-carries-mechanics** principle and stay vendor/tracker/
runtime-agnostic.

### Track A — `worktree-sessions`
1. **Model:** isolated working dir per session/implementer over one shared `.git`;
   converts silent concurrent runtime corruption into normal git merge conflicts.
2. **Worktree location:** sibling dir **`../worktrees/<slug>-<role>-<agent-id>`**
   (configurable, e.g. `$PARLEY_WORKTREES_DIR`). **NOT** inside `.git/worktrees/`
   (that is git's own admin metadata — antigravity-1's `.git/worktrees/<…>` path is
   rejected as a footgun; everyone else used a sibling dir).
3. **Branch discipline:** unique branch per session/agent (`idea/<slug>` design,
   `feature/<slug>/<agent-id>` impl, `integration/<slug>` staging); never check out
   the same branch in two worktrees.
4. **Isolation:** per-worktree `.env.local` (gitignored); sparse-checkout (cone) to
   scope files (collision + token savings); per-worktree submodule init.
5. **Merge:** per-agent branch → `integration/<slug>` → test → main (PR under
   github-pr transport). Disjoint **file sets** make parallel worktrees safe by
   construction; the skill refuses/【warns on】 concurrent sessions with intersecting
   file sets (override explicit).
6. **Cleanup gates:** `git worktree remove` + `git worktree prune` (+ `repair`);
   stale-worktree and IDE dual-index hazards documented (skill writes editor excludes).
7. **Runtime overlap:** detect-and-adopt a runtime-created worktree (Claude Code /
   Cursor / Copilot CLI) instead of creating a second; the skill adds the
   claim + file-set-disjointness + lifecycle layer those runtimes don't.
8. **Coordination manifest:** parley-native runs use a **"Worktree allocation" table
   in `IMPLEMENTATION.md`** (reuse the canonical file as the lock manifest — codex-1/
   antigravity-1); when the `ticket-tracker` skill is in play, each subtask ticket
   also carries a `worktree: {path,branch,base}` field (hermes-1). Same claim concept,
   two layers; `IMPLEMENTATION.md` stays canonical for Parley's own implementation.
9. **Core seam (separate meta-protocol-change):** one §6/Phase-5 note —
   "parallel implementers MAY be isolated in git worktrees; coordination uses the
   shared `IMPLEMENTATION.md` lock manifest; merge via an integration branch;
   worktrees never change canonical artifact ownership."

### Track B — `ticket-tracker`
1. **Tracker = MIRROR, files canonical.** Tickets live as markdown files
   (`tickets/<epic>/<story>/<subtask>.md` or under the idea dir). Sync is **one-way
   file→tracker** by default; optional `--pull` only writes back fields flagged
   `mirror-owned` (default `[status, assignee]`) to prevent silent drift.
2. **One file, three audiences:** audience-tagged sections **`[B]` business /
   `[T]` technical / `[A]` AI directives** + a mandatory **"At a glance"** 2-4 line
   cross-audience block + a **YAML frontmatter** machine block. Layered, not duplicated.
3. **Metadata schema (neutral, mapped to trackers at the edge):** `id` (immutable),
   `type` (epic|story|subtask), `title`, `status`, `priority`, `parent`, `labels`,
   `assignee`, `files`, `apis`, `arch`, `worktree`, `dod`, `mirror-owned`,
   `canonical_source` (FINAL.md path + git sha). Graceful degradation when a tracker
   lacks a field (drop from mirror, never from file).
4. **Acceptance criteria — hybrid + tagged:** behavioural → **Gherkin** (Given/When/
   Then), one scenario per `AC-N`; non-functional → **measurable bullet with a
   mandatory `Verify:` command**. Every AC carries an audience tag. Required: ≥1
   happy-path AC **and** ≥1 edge/error/offline AC (or explicit `n/a (reason)`).
5. **No-assumption enforcement (the highest-leverage rule):** a **gap-scan** the agent
   runs **before claiming** a ticket; any required slot missing and not `n/a` →
   agent returns `BLOCKED` / opens an `Open questions [needs:…]`, never invents.
   Enforced by the tool (`claim` runs gap-scan + refuses `in-progress` on failure),
   not just the template. (AI agents do not fill gaps the way humans do.)
6. **DoD = checklist of AC ids**, not prose (verifiable). INVEST governs story quality;
   story-splitting by vertical slice → disjoint file sets (ties to worktree parallelism).
7. **FINAL.md → tracker mapping:** idea/`00-prompt` → epic `[B]`; `FINAL.md` design +
   acceptance → stories + story ACs; `IMPLEMENTATION.md` checklist → subtasks;
   review-consensus deferred items → follow-up tickets. `FINAL.md` canonical for
   *design*, tickets canonical for *work-state* — no overlap.
8. **Readiness lint / `validate`:** a checklist/validator that rejects vague tickets
   ("improve/support/handle" with no observable criteria) and AI-unready tickets.
9. **Authoring:** vendor-neutral; optional MCP/API connector is an add-on, not core —
   the skill authors the text + emits a payload; actual create/update is opt-in.

### Distribution (from prior discussion, recorded here)
- core + opt-in plugins; install-all-by-default with `--no-addons` / `--only <skill>`;
  transparent install output; independent per-skill versions; advertised identity
  stays "parley-deck, the protocol".

## Trade-offs / resolved divergences
- Worktree path: sibling `../worktrees/` (correct) over `.git/worktrees/` (footgun).
- Coordination: `IMPLEMENTATION.md` table (parley-native) **and** per-ticket `worktree`
  field (ticket-tracker) — complementary layers, not competing.
- Gherkin vs bullets: Gherkin for behaviour, measurable+`Verify:` for NFR (don't force
  NFRs into Gherkin).

## Deferred follow-ups (non-blocking)
- `--pull` conflict-resolution policy (tracker-wins / file-wins / surface-as-conflict).
- Dependency lock-file merge conflicts on integration (serialize dep changes or
  regenerate on the integration branch) — document as a known pitfall.
- Tracker minimum-capability vs graceful-degradation threshold.
- `paused` ticket/worktree status so a worktree can outlive an in-progress subtask.
- The two thin core carve-outs (§6/Phase-5 worktree note; tracker-mirror note) →
  a separate meta-protocol-change idea, not this research idea.

## Signoffs

<!-- Each participant appends its own block. -->

### Signoff: claude-1 — 2026-06-22
Status: ✅ ACCEPT
Synthesis captures all four lenses; the worktree-path correctness call and the
two-layer coordination reconciliation are the only real edits over round-01.
FINAL.md will carry the full templates + command sequences verbatim from codex-1/
hermes-1 with these decisions applied.

### Signoff: codex-1 — 2026-06-22
Status: ACCEPT
Worktrees verdict: agreed commands, branch discipline, and cleanup gates are correct and sufficient to author the skill.

### Signoff: hermes-1 — 2026-06-22
Status: ACCEPT
Ticketing verdict: three-audience [B]/[T]/[A] + "At a glance" + frontmatter (§B2), tool-enforced gap-scan no-assumption rule (§B5), and Gherkin-for-behaviour / measurable+Verify: split (§B4) all captured correctly.

<!-- Facilitator note (claude-1, 2026-06-22): antigravity-1 signoff hung (known agy write-hang). Quorum 3/4 ACCEPT; agy round-01 input incorporated and the .git/worktrees->../worktrees correction was the facilitator call it was asked to accept. antigravity-1 WAIVED per the established per-idea pattern. -->
Status (antigravity-1): WAIVED — signoff append hung; round-01 input incorporated.
