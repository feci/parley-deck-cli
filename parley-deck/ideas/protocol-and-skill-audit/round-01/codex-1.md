---
agent: codex-1
idea: protocol-and-skill-audit
round: 1
date: 2026-08-20
---

## Findings

### F1 — `consensus draft` treats blank participant files as a completed round
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --round 1 --by alice audit-empty 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Drafted consensus at /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-empty/consensus.md
Consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-empty/consensus.md
Missing signoffs: alice,bob
Signoffs: none
exit=0
```
contradicts: `parley-deck/COOPERATION.md:308-322`; `internal/consensus/consensus.go:118-129,486-493`
why it matters: PRIMARY — Both `round-01/alice.md` and `round-01/bob.md` in the isolated fixture contain only a newline, but the command printed “Drafted consensus” and exited 0. Phase 1 requires each artifact's frontmatter and four named sections, while this command's gate checks only whether each expected pathname exists. The manual CLI can therefore advance an idea with no participant analysis at all.

### F2 — `consensus draft --review` requires the implementer to review their own work
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --round 1 --by reviewer-a audit-review 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
consensus draft failed: review/round-01 is incomplete; missing impl.md
exit=1
```
contradicts: `parley-deck/COOPERATION.md:503-525`; `internal/consensus/consensus.go:118-129,486-493`
why it matters: PRIMARY — The isolated `standard` fixture lists `impl`, `reviewer-a`, and `reviewer-b`; both non-implementers supplied valid review files, and `impl` correctly supplied none. The command nevertheless applies the full idea participant list to `review/round-01` and refuses Phase 7. A protocol-compliant manual Phase 6 cannot reach review consensus without fabricating the expressly prohibited implementer review.

## What I checked and found clean

## What I could not check, and why
