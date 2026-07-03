---
idea: meta-protocol-change-devx-speed
status: ready-for-review
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

## Observable acceptance criteria status
1. `track:` supported in `00-prompt.md`, default standard, deliberation force-triggers — **met** (protocol text).
2. Per-track reduced-ceremony behavior specified and usable by facilitator/agents — **met** (protocol text); CLI auto-enforcement — **deferred** (follow-up).
3. All-track invariants (non-solo + refutation) stated as never-dropped — **met**; driver rejection — **deferred**.
4. Core-first + appendix signposting + off-ramp + role table + LE gloss — **met**; physical appendix move — **deferred** (deviation above).
5. Both copies byte-identical (drift guard) + skill fallback re-synced — **met**.
6. protocol-changelog entry + protocolSha256 handling — **pending release step**.
