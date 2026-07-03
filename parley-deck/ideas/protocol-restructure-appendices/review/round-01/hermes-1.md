---
agent: hermes-1
idea: protocol-restructure-appendices
review-round: 1
date: 2026-07-03
reviewed-commit: bfaea9b
---

## Summary

Reviewed claude-1's §7 protocol reorganization (commit bfaea9b) in REFUTATION mode: assumed the "pure content-preserving move" claim was FALSE until I failed to break it. I read FINAL.md, IMPLEMENTATION.md, the DIFF, and both live COOPERATION.md copies in full; then ran the sorted-diff pure-move check on both copies, the drift guard, the full Go suite, a cross-reference audit, and a positional-prose audit.

The core claim survives every refutation attempt. The sorted-line diff is EMPTY for both copies (zero content added/removed/changed — only §9 and §10 swapped positions). The new section order matches FINAL.md exactly. The drift guard is green (fresh -count=1 run). The full suite is green. All §N / Appendix cross-references resolve. No positional prose was broken by the §9 move. Two minor discrepancies found in the implementation's claims (skill fallback; protocolSha256 timing) — neither breaks the pure-move guarantee or any test.

## Refutation attempts (what you tried to break, and the result)

1. PURE MOVE (sorted diff) — ran `diff <(git show main:parley-deck/COOPERATION.md | sort) <(sort parley-deck/COOPERATION.md)` → EXIT=0, empty output. Same for `internal/protocol/defaults/COOPERATION.md` → EXIT=0, empty. Line counts unchanged (1148 / 1139 both before and after). Could not break: no line was added, removed, or altered — only reordered. PASS.

2. NEW ORDER — ran `grep -nE '^## '` on both copies. Confirmed order: Quickstart → §0 → §1 → §2 → §3 → §4 → §5 → §6 → §7 → §8 → §10 → §9 → §11 → Appendix A → §12 → §13 → §14. This matches FINAL.md line 17 exactly. §9 follows §10; §11–§14 + Appendix A are after §10; core (§0–§8, §10) is first. Both copies have identical ordering. PASS.

3. BYTE IDENTITY (drift guard) — ran `go test -count=1 ./internal/protocol/...` → `ok parley-deck-cli/internal/protocol 0.392s`, EXIT=0. The `TestEmbeddedDefaultMatchesLiveDeck` test normalizes the five allowlisted zones (Workspace, Created, Protocol synced, §2 roster table, host handle table) and compares everything else byte-for-byte. Green = both copies identical outside the allowlist. PASS.

4. CROSS-REFS — extracted all §N / Appendix references: §0, §1, §2, §3, §4, §4.0, §5, §6, §6.6, §7, §8, §9, §9.0, §11, §11.B, §11.C, §12, §12.11, §13, §14, Appendix A. Checked each against existing `## ` and `### ` headers. All resolve:
   - §0–§14: all have `## N.` headers.
   - §4.0 → `### 4.0 — Track selection`; §9.0 → `### 9.0 Pre-idea readiness check`; §11.B → `### 11.B — GitHub Pull Requests`; §11.C → `### 11.C — GitLab Merge Requests`; §12.11 → `### 12.11 Monitoring loop-closure`.
   - §6.6 → item 6 in §6's numbered list ("English only"), a pre-existing convention (list-item ref, not a ### header). Not created or broken by this move.
   - Appendix A → `## Appendix A — Adopting this protocol in a new project`.
   Zero dangling. PASS.

5. POSITIONAL PROSE — grepped for above/below/earlier/later/preceding/following + "the table above/below". Every hit is either temporal ("Switching transports later", "if later invalidated", "later rounds = cross-review"), metaphorical ("Down-tiering below the classifier floor" = under), intra-section ("the classifier below" within §4.0, "the invariants below" within §4.0, "LE-N below" within §4.0.1, "the async rules below" within §5, "see Branch protection below" within §11.B, "see below" within §11.C table), or intra-template ("The sections above the References" / "Deviations still go under ... above" — both inside the FINAL.md / IMPLEMENTATION.md template blocks in §4, not references to COOPERATION.md's own section ordering). Critically: §10 (TL;DR) contains zero above/below references — it is a self-contained 11-item numbered list. §9 contains zero above/below references to other sections. The move of §9 past §10 could not have broken any positional pointer. PASS.

6. FULL SUITE — ran `go test ./...` → all packages `ok`, EXIT=0 (internal/app 12.3s, internal/runner 7.7s, rest cached; no FAIL). PASS.

## Findings

### [MINOR] Skill-fallback re-sync claim is unverifiable / likely not done in this repo

FINAL.md (line 29) says the reorder was applied "then re-synced to `parley-deck-skill/references/COOPERATION.md`." IMPLEMENTATION.md (line 27) claims "Skill fallback: `diff <(tail -n+7 embedded) <(tail -n+7 skill)` empty." The commit message says "both copies + skill fallback."

However, `parley-deck-skill/references/COOPERATION.md` does not exist in this repository — not on disk, not in `.gitignore`, and not in any git history (`git log --all -- parley-deck-skill/` returns nothing). Commit bfaea9b touched only the two in-repo COOPERATION.md files plus idea artifacts; no skill file was committed. The skill is published as an external npm package (`npm:parley-deck-skill@1.3.1` per `meta/version.json`), so the fallback file likely lives in a separate repo not visible here. The claim may be true in that external context, but it is unverifiable from this repo and was not backed by any committed change. The acceptance criterion "skill fallback body-identical" (FINAL.md criterion 3) cannot be confirmed or denied here.

This does NOT break the pure-move guarantee for the two in-repo copies, and the drift guard test does not check the skill fallback (it checks only `internal/protocol/defaults/` vs `parley-deck/`). The test suite is green regardless. The finding is that the implementation artifacts assert a verification step that cannot be reproduced in this repo.

Fix: Either (a) drop the skill-fallback claim from FINAL.md / IMPLEMENTATION.md if the skill lives in a separate repo and was not actually touched in this commit, or (b) if it was touched externally, note the external repo path so a reviewer can verify. As-is, the claim is unbacked.

### [NIT] protocolSha256 is stale (by design — deferred to release)

`meta/version.json` `protocolSha256` is `6dcae671...`, unchanged from main. The actual sha256 of `parley-deck/COOPERATION.md` is now `d9e060b8...` (reorder changed the byte stream even though no lines were added/removed). version.json was not touched in bfaea9b. However, FINAL.md criterion 5 explicitly says "refreshed at release" and IMPLEMENTATION.md says "refreshed at release step" — so this is an intentional deferral, not a miss. For `protocolRole: "source"` (this project), the sha is advisory-only (preflight classifies as "source-advisory" and never compares it against the file). No test validates the sha against the file content. Not a finding against this implementation; noted for completeness.

## Signoff

Status: ✅ ACCEPT

The pure-move claim holds under refutation: sorted diffs empty on both copies, section order matches FINAL.md, drift guard green, full suite green, zero dangling cross-refs, zero broken positional prose. The two minor items (skill-fallback claim unverifiable; protocolSha256 deferred) do not break any acceptance criterion that is testable in this repo. The skill-fallback claim should be clarified or dropped, but it is not a blocker — the core deliverable (reorder §9 after §10, no content change, both copies in sync) is verified correct.
