---
agent: hermes-1
idea: parley-learn-playbooks
review-round: 2
date: 2026-07-04
reviewed-commit: 8825dd3
responding-to: [hermes-1/review/round-01, codex-1/review/round-01]
---

## Position changes since prior review round

My round-01 verdict was ACCEPT with two MINOR findings and one NIT, none
blocking. All three were addressed in fix-up cycle 1 (commit 8825dd3) and I
have re-verified each on this branch. I also confirmed codex-1's MAJOR
write-boundary fix landed with a test. My position is unchanged: ACCEPT, now
clean.

## Updated findings

Zero remaining.

The three items I raised in round-01 are resolved:

1. [MINOR, resolved] §13.5 now reads "`parley learn <closed-idea-slug>`
   scaffolds a reusable `parley-deck/playbooks/<topic>.md` from a COMPLETED
   idea — a deterministic skeleton (track, roster, phase checklist, plus
   prompts for gotchas + fixes and the verification pattern) that the author
   refines into transferable, idea-agnostic prose before committing." This
   matches my suggested "seeds a skeleton" framing and no longer oversells
   the v1 tool as full auto-capture. The text is byte-identical in both
   `parley-deck/COOPERATION.md:1113` and
   `internal/protocol/defaults/COOPERATION.md:1104` — confirmed by reading
   both and by `TestEmbeddedDefaultMatchesLiveDeck` (drift guard) PASS.

2. [MINOR, dismissed] `playbooks/` absent from the §3 directory-layout tree:
   dismissed as consistent with `consults/` (also advisory, also not in the
   §3 tree). Confirmed — `consults/` is named in §8 but absent from §3 in
   both copies; adding `playbooks/` without `consults/` would be
   inconsistent, and adding both is scope creep. No action.

3. [NIT, resolved] IMPLEMENTATION.md now states the skill fallback "lives in
   the sibling `parley-deck-skill` repo (not this repo); it is re-synced
   there as part of the release" (IMPLEMENTATION.md:65-66). This resolves the
   unverifiable-from-this-branch ambiguity I flagged.

Cross-check on codex-1's MAJOR (not my finding, but I am responding to both
round-01 reviews): the symlinked-parent write boundary is hardened —
`learn.go:76-79` Lstats the `playbooks/` parent and refuses a symlink;
`learn.go:93` creates the playbook with `os.OpenFile(O_CREATE|O_EXCL|O_WRONLY)`
(atomic exclusive create, closes the TOCTOU); `TestLearnRejectsSymlinkedPlaybooksDir`
covers the symlinked-parent case. `go test ./internal/app ./internal/protocol`
both pass.

## Open questions

None. My round-01 open questions are both settled by the §13.5 rewording
(q1) and by the advisory/skeleton framing the author-refines language now
makes explicit (q2).
