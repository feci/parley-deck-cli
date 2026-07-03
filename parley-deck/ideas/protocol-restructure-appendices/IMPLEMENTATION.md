---
idea: protocol-restructure-appendices
status: ready-for-review
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
  line count unchanged (1139).
- **Drift guard:** `go test ./internal/protocol/...` green (both copies byte-identical outside the allowlist).
- **Skill fallback:** `diff <(tail -n+7 embedded) <(tail -n+7 skill)` empty.
- **Cross-reference audit:** 19 distinct `§N`/`Appendix` refs, **0 dangling**.
- **Full suite:** `go test ./...` exit 0, no FAIL.
- **Positional prose audit:** the only "see below"/"later" hits are temporal (§4 "if later
  invalidated") or intra-§11 ("see below" within the transport table) — none falsified by moving §9.

### Deviations from consensus
None. Implemented exactly as ratified (minimal §9 move, no banner, keep numbers).

## Observable acceptance criteria status
1. Section order matches FINAL (§9 after §10; §11–§14 + AppA after §10) — **met**.
2. Pure move (empty sorted-line diff, both copies) — **met**.
3. Drift guard green + skill body-identical — **met**.
4. Zero dangling cross-references — **met**.
5. Full suite green; `protocolSha256` — **refreshed at release step**.
