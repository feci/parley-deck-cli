---
idea: meta-protocol-change-review-gate-honesty
agent: claude
status: implemented
date: 2026-06-12
---

## Summary of work

Protocol amendments landed per consensus D1-D6, in lockstep across both
in-repo copies:

1. **Phase 6 += "Review briefs and dispositions"** (inserted before Phase 7 in
   `parley-deck/COOPERATION.md` and
   `internal/protocol/defaults/COOPERATION.md`): no-suppression rule, the
   disposition shape (including agy's closing "Reviewer prompt" line), the
   reviewer-concurrence requirement, and the three closure paths for disputed
   findings.
2. **Phase 8 += "Strict review gate (optional)"**: `strict_gate: true`
   frontmatter semantics (exact case-insensitive true; default rule unchanged
   otherwise), fresh full-scope close standard with hermes's full-scope
   definition, the objective-finding boundary (subjective style is never a
   finding; NITs stay blocking), and consensus-or-operator mutability.
3. **Phase 8 += "Stopping judgment"**: trajectory rules
   (converging/churning/blocked) with agy's concrete triggers as illustrative
   examples; MaxFixupCycles framed as an escalation threshold, never a close
   criterion.
4. **§8 += "Consults"**: parley-deck/consults/ artifacts are advisory and
   non-canonical; promotion requires normal protocol artifacts.
5. **Changelog**: dated entry in `parley-deck/meta/protocol-changelog.md`.
6. **External snapshot flag**: inbox note
   `claude-to-all_review-gate-honesty_external-skill-snapshot-sync.md` records
   that the out-of-repo parley-deck-skill bundled fallback still needs a sync,
   and re-flags the embedded default's pre-existing §12 drift as a separate
   follow-up.

## Implementation plan / checklist

- [x] Both COOPERATION.md copies amended with identical text (anchor-based
      insertion before "### Phase 7", before "### Escalation to user", and
      before "## 9. Session-start checklist").
- [x] Protocol changelog entry.
- [x] Inbox note for the external skill snapshot.
- [x] Ships with the sibling `runner-hardening-kindly` implementation as 1.24.0.

## Deviations from FINAL.md

- None. No driver enforcement was added (per D5); `strict_gate` is normative
  for humans now and machine-readable later via the documented
  ReadAutoImplement-style parser.

## Notes for reviewers

The amendment text matches codex's round-01 blocks with the merges agreed in
round-02/consensus; check both copies stay byte-identical in the amended
sections.
