---
idea: driver-impl-phase
review-cycle: 1
drafted-by: claude
date: 2026-06-06
reviewed-commit: f624c05
---

Synthesis of Phase 6 review-round-01 from codex, agy, hermes. hermes: no findings.
codex + agy converged on the same top issue (RunChecks after fix-up) plus hardening.
Agreed fixes applied in fix-up cycle 1.

## Agreed fixes

### AF1 — RunChecks after each fix-up before re-review (CRITICAL: agy + codex MAJOR)
`advanceReview` runs `Fixup` then opens the next review round with NO `RunChecks` in
between (D4 requires checks after implement AND each fix-up). Fix: after `Fixup`, run
`RunChecks`; failure → escalate (do not archive/open the next round). Regression test.

### AF2 — Fix-up loop crash-idempotency (MAJOR: agy)
A crash after `Fixup` but before opening round N+1 re-enters PhaseReview at round N
with the same consensus and re-runs `Fixup` (duplicate code-writing). Fix: a
driver-owned per-round marker `review/round-NN/.fixup-done`, written after a
successful Fixup+RunChecks. Checked at the top of the fix-up path so re-entry skips
straight to archive + open round N+1 without re-Fixup or re-draft. (Agent-set status
is unreliable — `BuildFixupPrompt` only says "keep status accurate".)

### AF3 — Review-consensus drafter must be a non-implementer (MAJOR: agy)
The adapter set the review-consensus drafter = implementer; a biased implementer
could filter reviewer findings. Fix: drafter = a reviewer (non-implementer).

### AF4 — Quote-tolerant ReviewStatus parsing (MAJOR: agy)
`outstanding_agreed_fixes`/`blocked` are `strconv.Atoi`/compared without stripping
quotes; `outstanding_agreed_fixes: "0"` would fail closed. Fix: strip `"'` first.

### AF5 — Malformed reviewer file must not infinite-retry (MAJOR: agy)
A reviewer artifact that exists but fails validation makes `ReviewRoundComplete`
false; re-opening the round skips the existing file (Overwrite=false) → spin to
deadline. Fix: before re-opening an incomplete round, remove reviewer artifacts that
fail validation so the agent regenerates them.

### AF6 — Implementer resolved from role metadata (MINOR: codex, D10)
The adapter always uses `participants[0]`. Fix: resolve the implementer from
IMPLEMENTATION.md `implementer` (re-entry), else FINAL.md `implementer`/`drafted-by`,
validate it is a participant, else fall back to `participants[0]`; reviewers = the
rest.

### AF7 — Not-ready implementation status awaits, not escalates (MINOR: codex, D6)
`advanceImpl` escalates any non-ready status. Fix: empty/unknown → escalate
(fail-closed), but known in-progress states (`in-progress`/`wip`/`draft`) → await.

### AF8 — Quote-tolerant opt-in parsing (MINOR: agy)
`ReadAutoImplement`/`ReadCrossReviewRounds` don't strip quotes; `auto_implement:
"true"` would silently disable the opt-in. Fix: strip `"'` before comparing/Atoi.

### AF9 — gitTreeClean: a git error INSIDE a repo is unsafe (MINOR: agy)
`gitTreeClean` returns true on any git error, so a `.git/index.lock`/error inside a
real repo reads as clean. Fix: probe `git rev-parse --is-inside-work-tree` first; if
inside a repo and `git status` errors → treat as dirty (unsafe).

## Deferred follow-ups
- A `verification_command` 00-prompt key for non-Go workspaces (agy OQ) — RunChecks
  currently hardcodes `go test ./...` gated on go.mod.
- A distinct reviewer-retry limit vs the loop deadline (agy OQ).

## Dismissed findings
None.

## Signoffs

<!-- each participant appends its own ✅ / 🟡 / ❌ block -->

### claude — ✅ ACCEPT (2026-06-06)
Accept the synthesis. AF1 (RunChecks after fix-up) and AF2 (driver-owned fix-up
marker for crash-idempotency) close the real correctness/safety gaps; AF3 restores
implementer/reviewer separation for the consensus draft; AF4/AF8 (quote tolerance)
and AF9 (git-error = unsafe) harden the parsers and the safety precondition; AF5
prevents an infinite reviewer retry; AF6/AF7 align roles/status with D6/D10.
Implementing all nine in fix-up cycle 1.

### codex — ✅ ACCEPT (2026-06-06)
AF1-AF9 are applied, including the post-fixup RunChecks gate, fix-up marker idempotency, non-implementer review-consensus drafter, parser hardening, malformed-review pruning, role metadata resolution, in-progress await behavior, and safer git cleanliness handling. `go build ./... && go vet ./... && go test ./...` plus `GOOS=windows GOARCH=amd64 go build ./...` are green with explicit Go caches.

### agy — ✅ ACCEPT (2026-06-06)
All agreed fixes (AF1–AF9) are verified successfully. Unit tests coverage added, and local builds (including Windows cross-compilation) are fully green.

### hermes — ✅ ACCEPT (2026-06-06)
AF1-AF9 verified in 3336b37 (RunChecks post-fixup, .fixup-done marker, non-implementer drafter, quote stripping, pruning, role metadata, in-progress await, git probe); native+Windows builds/tests green.
