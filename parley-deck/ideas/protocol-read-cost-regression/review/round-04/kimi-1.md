---
agent: kimi-1
idea: protocol-read-cost-regression
review-round: 4
date: 2026-08-11
reviewed-commit: 41e6cd6 (v1.43.1)
---

verdict: CLEAN — no CRITICAL, no MAJOR, no MINOR, no NIT. The restoration is exact: the entire
`internal/runner` package is byte-identical to the pre-idea commit, both active input changes I
flagged in round 3 are fully reverted, and the protocol text in all three copies describes the
shipped overlay accurately.

## Summary

Fix-up cycle 4 claims that `frontier.go` and `frontier_test.go` are deleted and `runner.go` /
`phase58.go` are restored to their pre-idea form. The claim is true, and stronger than stated:
`git diff d4256a2 HEAD -- internal/runner/` is **empty** — the whole package, tests included, is
byte-identical to v1.42.1, the last pre-idea release commit. There is nothing left of the feature
to review, which is the point of adopting @codex-1's counter-proposal. @codex-1's Finding B is
closed by deletion, not relocated: no constant, no dormant path, and no source-grep guards remain
to carry the "compiled is not verified" risk. My own round-3 NITs are all resolved or moot (Q2).
The enablement gate survives as text in the idea artifacts, which is the right home for it — any
future compaction is a fresh implementation reviewed against the signed contract, not a flip of
unexecuted code. Build, vet, and the full suite are green on my own run.

Candidly: in round 3 I judged the tripwired constant sound and voted CLEAN on Finding B. The owner
adopted @codex-1's position over mine. The deletion moots that disagreement — the verification I
vouched for no longer has an object — and the record (IMPLEMENTATION.md cycle 4, CHANGELOG 1.43.1)
states the process failure plainly rather than burying it. I have no residual objection.

## Refutation attempts

- `git diff d4256a2 HEAD -- internal/runner/runner.go internal/runner/phase58.go` — empty. Then
  widened to the whole package: `git diff d4256a2 HEAD -- internal/runner/` — also empty. Not just
  the two named files: `phase58_test.go`, `runner_test.go`, `round_test.go` and every sibling carry
  zero frontier-era residue.
- `ls internal/runner/*.go` — `frontier.go` and `frontier_test.go` are absent.
- Repo-wide Go grep for `frontier`, `compactionEnabled`, `_ledger.md`, `roundContextSentence`,
  `ledgerFileName`, `carry-forward` — zero matches. No partial deletion, no orphaned symbol.
- Re-ran the toolchain myself: `go build ./...` green, `go vet ./...` green,
  `go test ./... -count=1` — every package ok, including `internal/runner` (9.7s) and
  `internal/protocol` (the drift guard, 2.8s).
- Read the round ≥ 2 prompt builder in the shipped source (runner.go:982-1011): the instruction is
  the pre-idea text — "READ every prior-round artifact below and respond to the other participants
  by name" — with no banner clause and no carry-forward sentence. The `review` and
  `review-consensus` dispatch cases (runner.go:913-924) call the pre-idea `gatherReviewContext`
  directly; `gatherReviewContextFull` no longer exists anywhere.
- `find parley-deck tmp-test-plugin -name '_ledger.md' -o -name 'protocol-overlay.md'` — neither
  file exists in any deck, so no walker behaviour around them is observable regardless.
- Compared the overlay paragraph across all three protocol copies byte-for-byte (Q4) and checked
  each claim in it against `internal/protocolcore/overlay.go`, `lock.go`, and `render.go`.
- Read `changeReport` (render.go:175ff) for the nil-overlay path: with `ov == nil` there are no
  payloads to witness, so the relocation reclassification cannot fire and the report reduces to the
  pre-idea `droppedContent` semantics; the rendered body is untouched when no operations compose.

## Q1 — Are runner.go and phase58.go genuinely back to pre-idea behaviour?

**Yes, exactly.** Byte-identical to d4256a2, not "equivalent". The dispatch in `buildPromptForRound`
is the pre-idea switch; `gatherPriorRounds` and `gatherReviewContext` are the pre-idea walkers; the
prompt templates are the pre-idea templates. There is no residue of the removed feature in the
package, and none anywhere else in the tree (the grep above). A residue hunt finds nothing to
report.

## Q2 — Is my own round-3 objection resolved, or merely relocated?

Resolved or moot, each of them:

| Round-3 finding | Disposition in 1.43.1 |
| --- | --- |
| NIT-1 — instruction references a banner that never exists | **Resolved by deletion.** The entire sentence is gone with the feature; the instruction is the pre-idea text. No banner, no ledger claim, nothing to misdescribe. |
| NIT-2 — "compiled and exercised" overstates the dormant path | **Moot.** There is no dormant path to describe. IMPLEMENTATION.md cycle 4 records the machinery as deleted, which is the accurate statement. |
| NIT-3 — no dispatch-level test for the design path | **Returned to the pre-idea baseline.** The indirection it guarded is deleted; runner.go:929 is again the direct `gatherPriorRounds` call, as it was before this idea was opened. The residual gap is the pre-idea state I already accepted, not something 1.43.1 introduces. |

On Finding B itself — @codex-1's, the one this cycle closes: **closed by deletion.** The risk was
that unreachable machinery plus source-level guards would give a one-line enablement unjustified
confidence. Nothing carries that risk now: there is no constant to flip and no unexecuted code to
rot. The one genuine service those guards performed — catching the `gatherReviewContextFull` ledger
leak before it shipped — dies with the ledger concept itself, so nothing of value is orphaned. The
enablement preconditions (validator with G3/G5/G6 end-to-end and mutation coverage) survive as
recorded text in IMPLEMENTATION.md and consensus.md; they bind a future implementation, which is
where a binding belongs.

## Q3 — Does 1.43.0/1.43.1 still change what an agent receives, compared to before this idea?

**No. Both active input changes I flagged in round 3 are fully reverted.**

1. **The `_ledger.md` exclusion is gone.** All walkers are the pre-idea walkers; no Go file
   references the name. One honest consequence, stated so the record is complete: a hand-planted
   `_ledger.md` in a round dir would now be rendered as if participant-authored again. That IS the
   pre-idea behaviour, so it is the intended end state of "restore pre-idea behaviour", not a
   regression; the safety property dies with the feature, and no such file exists anywhere
   (verified above), so the direction question is moot in practice.
2. **The instruction wording is gone.** Byte-identical package means exactly that — the prompt an
   agent receives in a round ≥ 2 deliberation, a review round, or a consensus draft is the
   pre-idea prompt.

The remaining 1.43.0 delta versus d4256a2 is the overlay slice (`internal/protocolcore`,
`internal/app/protocol.go`, the COOPERATION.md paragraph), which belongs to
`protocol-overlay-local-extension`'s record, not this idea's. For completeness: it is
operator-facing, not agent-prompt-facing; with no `protocol-overlay.md` in any deck and no lock in
this one, `Render` composes nothing and the nil-overlay path is pre-idea behaviour (verified in
refutation). The one agent-receivable surface it touches is the COOPERATION.md text — checked in
Q4.

## Q4 — Does the protocol text (all three copies) now describe the overlay accurately?

**Yes.** All three copies — `internal/protocol/defaults/COOPERATION.md`, `parley-deck/COOPERATION.md`,
and the bundled skill snapshot `parley-deck-skill/skills/parley-deck/references/COOPERATION.md` —
carry the identical paragraph, and each claim in it checks out against the shipped code:

- "the file grammar … exist and are extend-only" — overlay.go enforces exactly one operation of
  kind `extend`; `replace` is a parse-time refusal; unknown keys, aliases, second documents, and
  trailing bodies are refused.
- "the `parley.protocol-lock/v2` lock … exist" — lock.go implements it, including the deliberate
  refusal of pre-v2 flat locks with migration guidance.
- "composition at the terminal boundary" — render.go:130-141 appends the payload at the end of the
  normalized core body, which is the definition of the boundary the text implies.
- "the roster-annotation identity slot and the removal of prose-matched zone addressing do NOT" —
  confirmed absent: identity slots in render.go are the unchanged pre-identity mechanism, and zone
  addressing is still prose matching (`findLine` et al.). The "do not rely on the parts that are
  absent" instruction is therefore the safe and true one.

The `internal/protocol` drift guard (embedded == live deck) passes on my run. The pre-existing
byte-level lag between the skill snapshot and the embedded copy, recorded in this idea's round 1,
is a separate documented drift problem; the paragraph itself is identical in all three copies, so
the CHANGELOG claim "the protocol text in all three copies now says so" is true.

## Q5 — Anything that should make the owner yank or supersede 1.43.1?

**No.** 1.43.1 is strictly safer than 1.43.0 on every axis this idea touches: the measured
diagnosis and the signed contract are preserved in the artifacts, the prompt an agent receives is
the pre-idea prompt, there is no dormant machinery, and the suite is green. The process failure —
releasing 1.43.0 on a substituted gate after a MIXED verdict — is recorded in IMPLEMENTATION.md
cycle 4 and in the CHANGELOG's process note, not buried. Checked and deliberately not findings:

- `review/consensus.md` still describes the 1.43.0 ship state ("machinery inert by construction").
  It is a dated round-3 record that this round's consensus supersedes; the current truth lives in
  IMPLEMENTATION.md cycle 4. Bookkeeping, not a shipped defect.
- Pre-v2 flat protocol locks are now refused by `protocol check`/`render`. That is a deliberate,
  documented compatibility control of the overlay slice (a stale binary must be stopped), with
  migration guidance in the error text; it is operator-facing, belongs to the overlay idea's own
  review record, and this repo's deck has no lock at all, so nothing breaks here.
- The CHANGELOG's 1.43.0 entry still documents that release's contents. Correct changelog practice;
  the 1.43.1 entry directly above it records the removal and the reason.

## Findings

None.

## Open questions

- None for this release. For the eventual `protocol-ledger-validator` idea, the enablement gate as
  recorded (implemented validator; end-to-end and mutation coverage for G3, G5, G6 before any
  compaction ships) remains the binding precondition — now with the additional property that the
  implementation starts from a clean tree, so its review will cover executed code from day one.
