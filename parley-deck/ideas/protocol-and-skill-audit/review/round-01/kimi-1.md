---
agent: kimi-1
idea: protocol-and-skill-audit
review-round: 1
reviewed-commit: 0bb99031aa1ddc1027605424848daa3a696e9f00
date: 2026-08-21
---

## Summary

(review in progress)

## Refutation attempts

(pending)

## Findings

### [MAJOR] F2/F3 fix: the implementer's own signoff on a review consensus now makes the whole consensus `malformed`

`expectedRoundParticipants` drops the resolved implementer from the participant list before
`validateDocument` runs (`internal/consensus/consensus.go:102-113`), and `validateDocument`
errors on any signoff from an agent not in that list (`consensus.go:455-456`, "unknown
participant"). Malformed outranks every other triage (`consensus.go:490-491`).

That turns a signoff the old binary accepted — and which this deck uses constantly, the
implementer's "fix-up cycle N applied the agreed fixes" report on review/consensus.md — into a
hard error. PRIMARY: built the pre-fix binary from a1926ae (`git archive`, /tmp) and this commit
(`go build`), ran `consensus status --review` for all 78 ideas with both. Nine review consensuses
flip from `ready`/`partial` to `malformed`, e.g. `2026-06-02T21-14-49-finish-the-6-rem`
(ready→malformed, "line 56: unknown participant claude"), `embedded-default-protocol-resync`,
`meta-protocol-change-rho-retrospective-optimization`, `named-roster-presets`, `live-run-tui`,
`interactive-agent-mode`, `parley-deck-cli-plan`, `meta-protocol-change-agent-teams-patterns`,
`meta-protocol-change-review-gate-honesty`.

Two of the flipped ideas are IN FLIGHT, so this blocks real work, not just history:
- `rho-retro-tooling` (status: implementation): review/consensus.md line 63 `### Signoff: claude —
  2026-06-16` (the implementer's fix-up report) → `malformed`. The review-close gate escalates on
  anything but ready/reserved (`internal/driver/impl.go:210-211`), so resuming this idea under
  `--auto` dead-ends until someone hand-edits the implementer's signoff out of a signed artifact.
- `build-companion-skills` (status: round-01): review/consensus.md line 89 `### Signoff: claude-1
  — 2026-06-22` → was `partial` (missing antigravity-1), now `malformed`.

The quorum fix itself is right (the implementer must not be REQUIRED, and its signoff must not
COUNT toward review quorum). The defect is only that an extra implementer signoff is treated as
contraband instead of inert: validate signoff identity against the FULL participant list, and
compute `Missing`/quorum against the reduced one. Ideas whose review consensus predates this
change then stay `ready`/`partial` exactly as before, and the §6 exclusion still binds.

### [MAJOR] F5 two-step finalize: the `pipeline auto` caller never takes the second step, and the scaffold itself satisfies the pipeline's finality check

`autoDriveDeliberationBlock` (`internal/app/pipeline_cmd.go:736`) calls `consensus.Finalize` once,
ignores the new `Summary.Scaffolded` flag, and prints `auto: block %q finalized.` / returns 0
whatever happened. Nothing in that function (rounds → draft → signoffs → finalize,
`pipeline_cmd.go:677-741`) ever writes FINAL.md content, so under `pipeline auto` Finalize ALWAYS
takes the new scaffold branch: FINAL.md scaffold written, idea left OPEN, pipeline reports
"finalized".

Worse, the scaffold carries `status: final` in its frontmatter (finalTemplate,
`internal/consensus/consensus.go:765-768`), and the pipeline's own finality check is a line scan
for exactly `status: final` (`isFinalized`, `pipeline_cmd.go:1323-1329`; used by `planFinalized`,
`pipeline_cmd.go:1300-1311`). So from the next wave on, the block reads as plan-finalized forever:
`autoDriveDeliberationBlock` is only invoked `if !planFinalized(...)` (`pipeline_cmd.go:456`), the
second finalize step never runs, no drafter ever fills the scaffold, and an action block can reach
`pipeline execute` with an unwritten plan while the idea itself stays open.

Before this commit the same path closed the idea around the scaffold (the F5 bug itself); after
it, the idea stays open but the pipeline's report and gating are unchanged — the F5 protection is
void on one of the two paths that close ideas, and the manual two-step contract the fix introduced
is silently dropped by this caller. Concrete fix: in `autoDriveDeliberationBlock`, check the
returned summary; if `Scaffolded`, either invoke a drafter to write FINAL.md and re-run Finalize,
or stop with a truthful "scaffold written, block NOT finalized" non-zero/status message. And make
`isFinalized`/`planFinalized` reject a scaffold (e.g. via `protocol.FinalIsScaffold`) so a
template can never satisfy the pipeline's finality check.

## Open questions

(pending)
