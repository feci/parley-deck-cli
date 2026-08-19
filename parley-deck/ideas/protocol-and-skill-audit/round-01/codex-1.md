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

### F3 — Standard-track review consensus stays partial after every reviewer accepts
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --round 1 --by reviewer-a audit-review && /tmp/parley-audit-codex-1.DGmaGx/parley consensus signoff --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --agent reviewer-a --status accept audit-review && /tmp/parley-audit-codex-1.DGmaGx/parley consensus signoff --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --agent reviewer-b --status accept audit-review && /tmp/parley-audit-codex-1.DGmaGx/parley consensus status --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review audit-review 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Drafted consensus at /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl,reviewer-a,reviewer-b
Signoffs: none
Appended signoff for reviewer-a to /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl,reviewer-b
Signoffs:
  reviewer-a ✅ ACCEPT
Appended signoff for reviewer-b to /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl
Signoffs:
  reviewer-a ✅ ACCEPT
  reviewer-b ✅ ACCEPT
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl
Signoffs:
  reviewer-a ✅ ACCEPT
  reviewer-b ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:220-235`; `internal/consensus/consensus.go:92-102,370-426`; `internal/app/driver_impl.go:328-355`
why it matters: PRIMARY — The authoritative track table says `standard` review consensus is the reviewers who reviewed signing off. Both selected reviewers accepted, but the shared consensus validator always uses all idea participants and therefore demands `impl`. The auto-driver calls the same validator through `ReviewStatus`, so this is not confined to the manual command: the implemented standard-track quorum override does not bind at the close gate.

### F4 — `consensus finalize --by` accepts a non-participant who has no drafter claim
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus finalize --dir /tmp/parley-audit-codex-1.DGmaGx/repo --by mallory audit-final`
output:
```text
Finalized consensus and created /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/FINAL.md
Consensus: ready
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/consensus.md
Signoffs:
  alice      ✅ ACCEPT
  bob        ✅ ACCEPT
```
contradicts: `parley-deck/COOPERATION.md:398-400`; `internal/consensus/consensus.go:196-243,549-576`
why it matters: PRIMARY — `00-prompt.md` names `alice` as author and lists only `alice` and `bob` as participants; no artifact grants `mallory` a drafter handoff. The command nevertheless exits 0, prints “Finalized,” and writes `author: mallory`. The CLI's finalization authority is therefore an unchecked free-form flag rather than the initiator/accepted-volunteer rule.

### F5 — `consensus finalize` closes an idea with an empty `FINAL.md` outline
severity: MAJOR
evidence: PRIMARY
command:
`sed -n '1,120p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/FINAL.md && sed -n '1,12p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/00-prompt.md`
output:
```text
---
idea: audit-final
status: final
author: mallory
consensus-date: 2026-08-20
participants: [alice, bob]
---

## Final plan / specification

### Goal

### Scope

### Implementation details

### Tests

### Non-goals

### Verification

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
---
idea: audit-final
author: alice
created: 2026-08-20
track: standard
participants: [alice, bob]
status: final
---

## Problem / idea

Choose the better of two equally viable labels; this is a judgment call.
```
contradicts: `parley-deck/COOPERATION.md:402-430`; `internal/consensus/consensus.go:196-243,549-576`
why it matters: PRIMARY — The immediately preceding finalize command printed success and set `00-prompt.md` to `status: final`, but every substantive field in the generated source-of-truth artifact is empty, and the protocol-required Purpose, Context, Observable acceptance criteria, Idempotence, and Risks sections do not exist. The command's success message therefore overstates its own effect and permanently closes the idea around a scaffold instead of an authoritative specification.

### F6 — A standard-track unanimous judgment closes without §15.6's adversarial alternative
severity: MAJOR
evidence: PRIMARY
command:
`find /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final -maxdepth 1 -type d -name 'round-*' -print | sort && rg -n 'judgment call|Use the label blue|Choose blue|Adversarial alternative|related models|agreed position to be wrong' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/00-prompt.md /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/*.md /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/consensus.md; /tmp/parley-audit-codex-1.DGmaGx/parley status --dir /tmp/parley-audit-codex-1.DGmaGx/repo --idea audit-final 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/00-prompt.md:12:Choose the better of two equally viable labels; this is a judgment call.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/bob.md:10:Use the label blue.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/bob.md:14:Choose blue.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/alice.md:10:Use the label blue.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/alice.md:14:Choose blue.
Idea: audit-final
Status: final
Participants: alice,bob
Consensus: ready
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/consensus.md
Signoffs:
  alice      ✅ ACCEPT
  bob        ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:1339-1358`; `internal/consensus/consensus.go:196-243`; `internal/driver/consensus.go:39-72,162-206`
why it matters: PRIMARY — The fixture explicitly identifies its output as a judgment call; both independent files select blue without disagreement; only `round-01` exists; and the grep finds no adversarial-alternative or correlated-agreement text. Nevertheless the CLI closed the idea. Neither manual finalization nor the auto-driver's FINAL validator checks §15.6, so a mandatory verification-integrity close condition is decorative.

### F7 — `continue` tells a fast-track idea to open the cross-review round that fast skips
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley continue --dir /tmp/parley-audit-codex-1.DGmaGx/repo audit-fast-run 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Run: audit-fast-run
Idea: audit-fast-plan
State: completed
Recommended: Open round-02 (cross-review) before drafting consensus
Command: parley run --auto --dir . "continue audit-fast-plan"
Next actions:
  open-next-round.audit-fast-plan  kind=open-next-round  risk=normal  requires=yes
      Open round-02 (cross-review) before drafting consensus
      command: parley run --auto --dir . "continue audit-fast-plan"
exit=0
```
contradicts: `parley-deck/COOPERATION.md:218-235`; `internal/runplan/runplan.go:112-148,227-279`
why it matters: PRIMARY — The fixture's `00-prompt.md` explicitly says `track: fast`, both round-01 artifacts exist, and the run records round 1 complete. The continuation planner ignores `track:` and consults only the separate `cross_review_rounds:` key, defaulting that key to 1. Its printed recovery instruction therefore drives the user into a phase the binding fast-track route explicitly skips.

## What I checked and found clean

## What I could not check, and why
