---
idea: protocol-restructure-appendices
status: fix-up-cycle-1
implementer: claude-1
track: deliberation
started: 2026-07-03
head-commit: (pending)
---

## Progress

Pure content-preserving reorder applied to BOTH `COOPERATION.md` copies + the skill fallback,
on branch `protocol-restructure-impl`.

### What landed
- Relocated **§9** (session-start checklist) to sit after **§10** (TL;DR) via a deterministic
  reorder script (split on `## ` headers, re-emit in the ratified order). §11–§14 + Appendix A
  were already after §10 (no move needed). Section numbers unchanged → all `§N` refs resolve.
- Applied identically to `internal/protocol/defaults/COOPERATION.md` and
  `parley-deck/COOPERATION.md`; re-synced `parley-deck-skill/references/COOPERATION.md`.
- No `## Appendices` banner (per consensus D3 — keeps it a zero-addition move).

### Verification (all green)
- **Pure move:** `diff <(sort pre) <(sort post)` EMPTY for BOTH copies (zero content added/removed);
  line count unchanged (1139). (Confirmed independently by all three reviewers.)
- **Drift guard:** `go test ./internal/protocol/...` green (both copies byte-identical outside the allowlist).
- **Skill fallback:** re-synced in the **sibling repo** `parley-deck-skill/references/COOPERATION.md`
  (NOT inside parley-deck-cli — reviewers running in the CLI repo cannot see it); verified
  body-identical via `diff <(tail -n+7 embedded) <(tail -n+7 skill)` empty. This reorg changes the
  bundled fallback, so it ships as a skill patch release alongside the CLI release.
- **Cross-reference audit:** 19 distinct top-level `§N`/`Appendix` refs — **0 dangling**.
  Note (review round-01, codex-1): `§6.6` is a **sub-item convention** ("English-only rule" is list
  item 6 under `## 6. Conflict-avoidance mechanics"), not a markdown header — a header-only audit
  flags it, but the pure move PRESERVED it byte-for-byte (it is pre-existing, not move-created), so
  no reference regressed.
- **Full suite:** `go test ./...` exit 0 in the implementer's environment (and re-run green by
  hermes-1 and antigravity-1). In the **codex sandbox** the single test
  `internal/runner/TestDurableKillEndToEndRealProcess` fails with "no recorded boot id" — the
  **known, documented codex-sandbox limitation** (recurs on every review; unrelated to this
  reorder, which touches no runner code). Recorded as the standing sandbox exception, not a defect.
- **Positional prose audit:** the only "see below"/"later" hits are temporal (§4 "if later
  invalidated") or intra-§11 ("see below" within the transport table) — none falsified by moving §9.

### Deviations from consensus
None. Implemented exactly as ratified (minimal §9 move, no banner, keep numbers). The review-01
findings were claim-accuracy corrections (above), not changes to the reorder itself.

## Observable acceptance criteria status
1. Section order matches FINAL (§9 after §10; §11–§14 + AppA after §10) — **met**.
2. Pure move (empty sorted-line diff, both copies) — **met**.
3. Drift guard green + skill body-identical — **met**.
4. Zero dangling cross-references — **met**.
5. Full suite green; `protocolSha256` — **refreshed at release step**.
