---
agent: hermes-1
idea: parley-learn-playbooks
review-round: 1
date: 2026-07-04
reviewed-commit: 553ef14
---

## Summary

The implementation stayed minimal and faithful to the one-paragraph §13.5
consensus I championed. §13.5 is a single advisory paragraph in both
COOPERATION.md copies (byte-identical — diff confirmed), sitting under §13
beside the existing retro tooling note. It introduces no new phase, gate, or
quorum class. Playbooks are explicitly "advisory and non-canonical — like
consults (§8)" and `parley learn` is explicitly "a tooling command... NOT a
Parley round." The implementation mirrors `parley retro propose`'s safety
envelope (strict slug, Lstat fail-closed, single-file write, read-only over
the idea) and does not over-reach: no auto-application, no LLM distillation —
`distillPlaybook` is a deterministic skeleton with honest fill-in prompts for
the sections a human must refine. Tests pass, build/vet/gofmt clean, drift
guard green, protocolSha256 matches the live file.

Verdict: ACCEPT. Two MINOR findings and one NIT below; none block merge.

## Findings

### [MINOR] §13.5 protocol text oversells the v1 skeleton's coverage

§13.5 says `parley learn` "distills a COMPLETED idea into a reusable
`parley-deck/playbooks/<topic>.md` capturing the proven shape (track, roster,
phase checklist, gotchas + fixes, verification pattern)." The implementation
(`learn.go:124-169`) populates track, participants, and fix-up count from
frontmatter, and emits a static phase checklist. But "Gotchas & fixes" and
"Verification pattern" are left as parenthetical fill-in prompts
(`learn.go:160-161, 165-166`) — the tool does not extract them from the idea's
review consensus or IMPLEMENTATION.md. The protocol sentence reads as though
the command produces all five listed elements; in practice it produces two and
scaffolds three. This is the kind of gap my round-01 "garbage-in" concern
flagged: if the playbook is seen as auto-generated, the empty sections get
rubber-stamped. The skeleton labeling in the code comment (`learn.go:120-123`)
is honest; the protocol text is slightly less so. Suggested fix: qualify the
§13.5 sentence — e.g. "seeds a skeleton capturing the proven shape... with
gotchas and verification pattern left for the human to fill from the idea's
review artifacts" — or leave the text and accept that the tool comment carries
the honesty. Not blocking; the advisory status means a thin playbook is low
harm.

### [MINOR] `playbooks/` absent from §3 directory layout tree — but consistent with `consults/`

§13.5 introduces `parley-deck/playbooks/<topic>.md` as a directory convention,
but the §3 directory layout tree (COOPERATION.md:138-162) does not list a
`playbooks/` entry. This is the same precedent consults already set:
`parley-deck/consults/` is named in §8 (COOPERATION.md:734) but is also absent
from the §3 tree. So the omission is consistent with established practice, not
a new gap. I note it only because the FINAL.md (line 32) calls §13.5 "a single
paragraph + the `playbooks/` directory convention" — if the directory
convention is part of the protocol delta, one could argue it belongs in the
layout diagram. But adding it without also adding `consults/` would be
inconsistent, and adding both is scope creep beyond this idea. Leave as-is;
the §13.5 paragraph is sufficient.

### [NIT] IMPLEMENTATION.md claims skill fallback re-synced, but no skill fallback file is in this repo/worktree

IMPLEMENTATION.md:28 states "skill fallback `references/COOPERATION.md`
re-synced (body-identical from line 7)." The commit 553ef14 stat does not
include any `references/COOPERATION.md` path, and no such file exists in this
worktree (only `parley-deck/COOPERATION.md` and
`internal/protocol/defaults/COOPERATION.md` are present). The skill fallback
likely lives in the published npm package or a separate skill repo not checked
into this worktree, so this may be accurate from the drafter's perspective but
is unverifiable from this branch. No action needed if the re-sync happened out
of tree; flagging only so a future reviewer doesn't waste time searching for
the file here.

## Open questions

1. Should the §13.5 text be tightened to say the v1 tool "seeds a skeleton"
   rather than "distills... capturing" — to match the honest labeling in
   `learn.go:120-123` and set the right expectation that a human fills the
   gotchas/verification sections before the playbook is useful? (See MINOR 1.)

2. The `distillPlaybook` lifecycle line (`learn.go:146-147`) hardcodes
   "round-01 → cross-review (if divergent) → consensus + signoffs → FINAL →
   implement → refutation review → fix-up → complete." This is a generic
   description, not mined from the actual idea. For a `fast`-track idea that
   skipped cross-review, this line would be inaccurate in the generated
   playbook. Is that acceptable for v1 (human refines before commit), or
   should the skeleton detect the track and adjust? Low priority given the
   advisory/skeleton framing.
