---
agent: hermes-1
idea: meta-protocol-change-devx-speed
review-round: 2
date: 2026-07-03
reviewed-commit: 74bbdc5
---

## Summary

Fix-up cycle 1 (commit 74bbdc5) addresses every round-01 finding I raised. The
implementer added: (1) a single authoritative per-track override clause in §4.0
that explicitly names §4/§5/§9.0/§11 as overridden; (2) a Phase-7 review-consensus
row in the per-track table; (3) the heading rename from "Phase 0.0" to "### 4.0";
(4) a §4.0.1 "Loop-engineering rules (LE-N), in plain English" consolidation block;
(5) a normative fail-safe classifier-ordering paragraph; (6) corrected Quickstart
fast-track wording; (7) the mid-idea-upgrade reinstatement clause; (8) a
protocol-changelog entry. All changes are applied byte-identically to both the
live deck and the embedded default (drift guard green, fresh `go test -count=1`
run). The skill-fallback source at `parley-deck-skill/references/COOPERATION.md`
is now verified body-identical to the live deck modulo the five allowlisted
zones — my round-01 MAJOR-2 complaint (unverifiable "re-synced" claim) is
resolved by the corrected IMPLEMENTATION.md wording that distinguishes source
(met) from installed-runtime (pending skill-release step).

I tried to break this again. I could not break the override clause, the
classifier fail-safe, the invariant block, the drift guard, or the LE
consolidation. The only residual is a NIT: the §10 TL;DR still says "fast (one
review)" while the Quickstart now says "round-1 + collapsed FINAL + one
refutation-default reviewer (≤1 fix-up)" — a minor wording inconsistency in the
TL;DR's terse summary, not a safety or correctness issue.

## Verification of round-01 findings

- CRITICAL/reconcile (tracks not reconciled with §4/§5/§9.0/§11): **RESOLVED** —
  §4.0 line 205 adds "This table is the single authoritative per-track gate. It
  OVERRIDES the full-lifecycle defaults stated in the rest of §4 and in §5
  (quorum), §9.0 (readiness), and §11 (transport)" with a reader-redirect for
  "every participant / all reviewers / consensus" phrasing; Phase-7 row present
  at line 200 in both copies. The contradiction is removed — the table is now
  the explicit authority and old phase prose is read through it.

- §4.0-vs-Phase-0.0 heading mismatch: **RESOLVED** — heading is now "### 4.0 —
  Track selection (conditional rigor)" at line 172 in both copies (diff confirms
  rename from "Phase 0.0"); every §4.0 reference now points at the correct
  heading.

- LE-1…LE-11 consolidation (MAJOR-1): **RESOLVED** — §4.0.1 "Loop-engineering
  rules (LE-N), in plain English" at line 228 in both copies restates
  LE-1/2/3/4/5/7/11/10 in one block; inline `(LE-N)` tags now have a single
  plain-English reference point. Criterion 4's LE-consolidation half is met.

- Classifier fail-closed / normative ordering / tie-break (MINOR-1 + MINOR-2):
  **RESOLVED** — "Classifier ordering is normative and fail-safe" paragraph at
  line 185 in both copies states: evaluate deliberation first, then fast, else
  standard; on any doubt/boundary "fail closed to the stricter track: standard
  over fast, and deliberation over standard." Covers both the tie-break rule and
  the normative ordering in one block.

- Quickstart fast wording (MINOR codex, concur): **RESOLVED** — Quickstart line
  21 now reads "fast = round-1 + a collapsed FINAL.md signoff + one
  refutation-default reviewer (≤1 fix-up cycle)" instead of the over-compressed
  "one review, then done."

- Mid-idea upgrade from fast reinstates skipped phases (MINOR-4): **RESOLVED** —
  line 223 in both copies adds "Upgrading from fast also reinstates any phase
  fast skipped (cross-review, a separate consensus/FINAL step) for the remainder
  of the idea."

- Skill fallback re-sync (MAJOR-2): **RESOLVED (source) / PARTIAL (runtime, by
  design)** — the parley-deck-skill source at
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill/references/COOPERATION.md`
  is verified body-identical to the live deck (diff shows only the five
  allowlisted header/roster zones differ; grep confirms Quickstart + §4.0.1
  present). The installed runtime copy at
  `~/.hermes/skills/parley-deck/references/COOPERATION.md` is stale (sha
  0e986e…, no Quickstart/§4.0.1) — but that is a sibling repo that refreshes at
  the skill-release step. IMPLEMENTATION.md now correctly says "met (source);
  runtime pending release," resolving my round-01 complaint that the claim was
  unverifiable. The correction is honest and the source is verifiably synced.

- Changelog entry in meta/protocol-changelog.md (criterion 6): **RESOLVED** —
  2026-07-03 entry present with full summary, participant list, and deferred
  follow-up references. `protocolSha256` in `meta/version.json` remains stale
  (20b98556… vs current 3284a86…) but is correctly deferred to the release step,
  matching my round-01 "pending release step" assessment.

- go test ./internal/protocol/... — **GREEN** — `go test -count=1
  ./internal/protocol/...` → ok (0.294s, fresh non-cached run). Drift guard
  `TestEmbeddedDefaultMatchesLiveDeck` passes; both copies byte-identical modulo
  allowlisted zones.

## Scope decision

**ACCEPT**. The protocol is self-enforcing via the skill (agents read and follow
the §4.0 table manually); the tracks are usable today — the very next idea can
set `track: fast`. Deferring (a) CLI/driver enforcement to `track-aware-driver`
and (b) the 460-line physical appendix relocation to
`protocol-restructure-appendices` is a sound risk trade-off: the physical move
maximizes cross-reference/drift risk for zero capability gain, and deterministic
enforcement is a substantial separate code effort. The design was unanimously
ratified (consensus.md: ✅ ×4, zero reservations) with this follow-up approach
named. The protocol-text change delivers the functional DevX outcome now; the
follow-ups deliver mechanical enforcement and physical reorganization later on
their own reviewable diffs.

## New findings

### [NIT] §10 TL;DR "fast (one review)" slightly understates the fast track

§10 TL;DR item 0 (line 787) says "`fast` (one review)" while the Quickstart
(line 21) was corrected to "round-1 + a collapsed `FINAL.md` signoff + one
refutation-default reviewer (≤1 fix-up cycle)." The TL;DR's terse phrasing is
arguably acceptable for a one-line summary, but "one review" could be read as
"one look and done," omitting the collapsed-FINAL and refutation-default
requirements. Suggest harmonizing to "`fast` (1 reviewer, collapsed FINAL)" for
accuracy. Not blocking — the authoritative definition is in §4.0 and the
Quickstart, both of which are correct.

No other new findings. No CRITICAL, MAJOR, or MINOR issues found.

## Signoff

Status: ✅ ACCEPT

All round-01 findings are resolved. The two MAJOR gaps (LE consolidation,
skill-fallback claim) are fixed. The CRITICAL reconciliation gap is fixed with an
explicit override clause + Phase-7 row. The classifier is normatively ordered
and fail-safe. Tests are green. The scope narrowing to protocol-text + two
ratified follow-ups is a sound engineering trade-off. The one residual NIT
(TL;DR fast wording) is cosmetic and non-blocking. This idea is ready to close.
