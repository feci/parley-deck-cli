---
idea: meta-protocol-change-devx-speed
status: fix-up-cycle-1
implementer: claude-1
track: deliberation
started: 2026-07-03
head-commit: (pending)
---

## Progress

Implementing FINAL.md on the `deliberation` track (a protocol change is force-classified to
deliberation by its own new classifier). The change is **additive** to both protocol copies
plus a skill-fallback re-sync.

### What landed (protocol text — applied byte-identically to BOTH copies)

Edited `parley-deck/COOPERATION.md` **and** `internal/protocol/defaults/COOPERATION.md`
identically (drift guard `TestEmbeddedDefaultMatchesLiveDeck` green — `go test ./internal/protocol/...` ok):

1. **Quickstart & reading guide** (new block after the header, before §0): 5-minute start;
   "trivial reversible work needs no Parley" off-ramp; "Who are you? → read this" role table;
   and a **Core (§0–§8) vs reference appendices (§9/§11/§12/§13/§14)** progressive-disclosure map.
2. **§4.0 Track selection (conditional rigor)** (new subsection at the start of §4): the
   objective classifier table (deliberation-forcing triggers / fast / standard-default), the
   per-track behavior table (§9.0 ping, cross-review cap, consensus+FINAL collapse, reviewer
   count, fix-up cap, timeout, auto-advance), the all-track invariants, and binding / challenge
   / mid-idea-upgrade / roster-size-floor rules.
3. **`track:` field** added to the Phase 0 `00-prompt.md` frontmatter template.
4. **§10 TL;DR** gains a leading "Pick a track first (§4.0)" item.

### Deviations from FINAL.md (recorded per §4 Phase 5)

- **Physical appendix relocation + renumber → replaced by an in-place reading-guide map.**
  FINAL §4 / acceptance criterion 4 called for physically moving §9/§11/§12/§13/§14 to the
  back and a ≤~200-line core-before-first-appendix. I implemented the **functional** DevX
  outcome (a developer reads a short core and is told exactly which sections to skip) via the
  Quickstart's "Core vs reference" map, **without** physically moving ~460 lines or
  renumbering. Rationale: FINAL's own top risk is "restructuring breaks cross-references /
  consumer sync," and a 460-line move + renumber applied byte-identically across three copies
  maximizes exactly that risk while adding no capability. The additive approach keeps the diff
  reviewable, keeps every existing `§11.B`-style cross-reference valid, and keeps the drift
  guard tractable. **The pure physical reorganization is deferred to a dedicated follow-up
  idea** (`protocol-restructure-appendices`) where a section-move can be reviewed on its own.
  Acceptance criterion 4's "core ≤200 lines before first appendix" is therefore **partially
  met** (core-first signposting yes; physical appendix yes-by-reference, no-by-position).

- **CLI/tooling enforcement deferred.** FINAL §5 + acceptance criteria 1–3 reference tooling
  (a classifier command, `parley init/run` templating, per-track `agents.toml` timeout
  seeding, driver validation that rejects dropping non-solo/LE-1). The protocol is
  **self-enforcing through the skill** (manual facilitation — agents read and follow it), so
  the protocol-text change is functional today: the very next idea can set `track: fast`.
  The deterministic CLI enforcement is a substantial separate code effort and is deferred to
  a follow-up idea (`track-aware-driver`), noted so it is not lost.

### Follow-ups opened (deferred, documented)
- `protocol-restructure-appendices` — physical section move + renumber to appendices.
- `track-aware-driver` — classifier command, init/run templating, per-track timeout seeding,
  driver auto-advance + validation gates.

## Verification
- `go test ./internal/protocol/...` → ok (drift guard passes; both copies byte-identical
  modulo allowlisted header/roster zones).
- Skill fallback `parley-deck-skill/references/COOPERATION.md` re-synced (body-verbatim rule).

## Fix-up cycle 1
status: complete
completed: 2026-07-03

### Fixes applied (from Phase-6 review/round-01: codex-1, hermes-1, antigravity-1)
- **[CRITICAL codex-1 / concur hermes] Tracks not reconciled with existing phase rules** →
  added an explicit **single authoritative override clause** in §4.0 ("this table OVERRIDES
  the full-lifecycle defaults in §4/§5/§9.0/§11"), and added the missing **Phase-7 review-
  consensus row** to the per-track table. The old phase prose is now unambiguously overridden.
- **[MAJOR ×3] `§4.0` vs `Phase 0.0` heading mismatch / "Phase 0.0 before Phase 0"** → renamed
  the heading to `### 4.0 — Track selection`, matching every `§4.0` reference.
- **[MAJOR hermes/antigravity] LE-1…LE-11 not consolidated** → added `### 4.0.1 — Loop-
  engineering rules (LE-N), in plain English` (one block; inline `(LE-N)` tags now reference it).
  Criterion 4's LE-consolidation half is now met.
- **[MAJOR codex] Classifier not fail-closed** → added a normative "deliberation-first, then
  fast, else standard" ordering + a fail-safe rule (boundary/doubt → stricter track).
- **[MINOR codex] Quickstart over-compressed `fast`** → reworded to "round-1 + collapsed FINAL
  signoff + one refutation-default reviewer (≤1 fix-up)".
- **[MINOR hermes] mid-idea upgrade from `fast`** → clarified it reinstates skipped phases.
- **[MAJOR codex/hermes] changelog + metadata** → `meta/protocol-changelog.md` entry added
  (2026-07-03). `protocolSha256` in `meta/version.json` is refreshed at the release step (below).

Drift guard re-run after fixes: `go test ./internal/protocol/...` → ok. Skill fallback
re-synced again (body-identical).

### Deviations from agreed fixes
- **CLI/driver enforcement (criteria 1–3 tooling half) and physical appendix relocation
  (criterion 4 layout half)** remain **deferred to follow-up ideas** — stubs now exist at
  `ideas/track-aware-driver/00-prompt.md` and `ideas/protocol-restructure-appendices/00-prompt.md`.
  This scope narrowing is **proposed here and put to the Phase-7 review consensus (this round)
  for ratification** (not pre-ratified). The protocol text is self-enforcing via the skill, so
  the tracks are usable today.

## Observable acceptance criteria status (corrected)
1. `track:` field, default standard, deliberation force-triggers, fail-safe classifier — **met** (protocol text). Deterministic CLI classifier/defaulting — **deferred → `track-aware-driver`**.
2. Per-track reduced-ceremony behavior specified, authoritative-override clause added, usable by facilitator/agents — **met** (protocol text). CLI auto-enforcement — **deferred → `track-aware-driver`**.
3. All-track invariants (non-solo + refutation) stated as never-dropped + single authoritative gate — **met** (text). Driver rejection of dropped invariants — **deferred → `track-aware-driver`**.
4. Off-ramp + role table + reading guide + **LE consolidation now done** — **met**. Physical appendix move / ≤200-line core — **deferred → `protocol-restructure-appendices`** (reading-guide substitute in place; scope narrowing proposed for ratification in this Phase-7 review consensus).
5. Both copies byte-identical (drift guard green) + skill-fallback source re-synced in `parley-deck-skill/references/COOPERATION.md` (verified body-identical). Installed runtime skill copies refresh at the **skill release step** (they are a sibling repo, not this commit) — **met (source); runtime pending release**.
6. `meta/protocol-changelog.md` entry — **met**. `protocolSha256` bump in `meta/version.json` — **done at release step below**.
