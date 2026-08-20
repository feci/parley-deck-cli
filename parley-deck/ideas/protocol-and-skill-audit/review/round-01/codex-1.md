---
agent: codex-1
idea: protocol-and-skill-audit
review-round: 1
reviewed-commit: 0bb99031aa1ddc1027605424848daa3a696e9f00
date: 2026-08-21
---

## Summary

Review in progress.

## Refutation attempts

Review in progress.

## Findings

### [MAJOR] The pipeline treats the first-step FINAL scaffold as a completed block

PRIMARY — In an isolated copy I drove a one-participant deliberation to ready, called `consensus.Finalize` once, and then invoked the production `blockCompleteFunc` for that block. The first result reported `Scaffolded=true`, while the completion predicate returned `true`. The focused reproduction passed with: `go test -count=1 -v ./internal/app -run '^TestAuditReproPipelineTreatsFirstFinalizeScaffoldAsComplete$'`.

SECONDARY — `internal/app/pipeline_cmd.go:736-741` calls `Finalize` only once, ignores `Summary.Scaffolded`, prints “finalized,” and returns success. `internal/app/pipeline_cmd.go:1263-1277,1323-1329` then declares any non-action block complete when its artifact contains a line equal to `status: final`; the newly written scaffold contains exactly that frontmatter even though the idea intentionally remains open.

SECONDARY — This breaks the two-step invariant and permits the pipeline driver to seed or enter the next block from an unwritten specification. The concrete fix is to branch on `summary.Scaffolded` and stop as incomplete after step one, and to make block completion validate the written FINAL contract (and/or the idea status) instead of scanning for one frontmatter-looking line. Add an integration test covering `autoDriveDeliberationBlock` followed by pipeline advancement.

### [MAJOR] The automatic FINAL drafter is instructed to create an artifact its new gate rejects

PRIMARY — A focused test compared `buildFinalDraftPrompt` with `protocol.RequiredFinalSections`: the prompt omitted six of seven mandatory headings and never instructed the drafter to emit the required `idea:` frontmatter. The reproduction passed with: `go test -count=1 -v ./internal/app -run '^TestAuditReproFinalDrafterPromptCannotPassNewGate$'`.

SECONDARY — `internal/app/driver_consensus.go:131-142` requests only `status: final` and `## Final plan / specification`; the post-draft gate requires the matching idea slug and all seven protocol sections. The edited driver tests hand-built seven-section fixtures, so they do not exercise the actual prompt-to-gate contract.

SECONDARY — A compliant drafter following the shipped instruction can therefore be rejected and escalated even though it did exactly what the caller requested. The concrete fix is to generate the drafter instruction from `protocol.RequiredFinalSections`, explicitly require `idea: <slug>` plus `status: final`, and add a test that a literal prompt-compliant output passes the same production validator used after drafting.

### [MAJOR] Manual `consensus finalize` closes the requested idea around another idea's non-final artifact

PRIMARY — In an isolated copy I completed step one, replaced the scaffold with substantive content containing all seven headings but frontmatter `idea: other-idea` and `status: draft`, and ran `Finalize` again. It returned success and changed `sample/00-prompt.md` to `status: final`. The reproduction passed with: `go test -count=1 -v ./internal/consensus -run '^TestAuditReproFinalizeAcceptsWrongSlugAndStatus$'`.

SECONDARY — `internal/consensus/consensus.go:257-278` applies only `protocol.FinalIsScaffold` to an existing FINAL. `internal/protocol/finalsections.go:42-57` explicitly validates content only and says callers own frontmatter and slug checks; this caller performs neither. The auto-driver has a stricter, separate gate, leaving manual and automatic finalization inconsistent.

SECONDARY — The final status can therefore authenticate the wrong artifact and the wrong lifecycle state. The concrete fix is one shared FINAL validator that accepts the expected idea slug and requires `status: final`, all seven exact sections, non-placeholder content, and the specification content floor; call it from both manual finalization and the driver.

### [MAJOR] Manual review consensus still accepts review artifacts without `reviewed-commit`

PRIMARY — In an isolated idea with implementer `impl` and reviewer `rev-a`, I wrote a structured review containing no `reviewed-commit` and invoked `consensus.Draft(... Review:true)`. It created `review/consensus.md` without error. The reproduction passed with: `go test -count=1 -v ./internal/consensus -run '^TestAuditReproManualReviewDraftAcceptsMissingReviewedCommit$'`.

SECONDARY — `internal/consensus/consensus.go:129-142` gates review drafting through `missingRoundArtifacts`, which checks only body content and a heading. The new `runner.ValidateReviewArtifact` requirement is not used by this manual Phase 7 entry point. This leaves the reviewed-commit fix path-dependent.

SECONDARY — A manual review consensus can synthesize reviews whose target revision is unknown, defeating the provenance field this audit intended to require. The concrete fix is to share the review-artifact validator below `runner` (or add an equivalent frontmatter check in `consensus`) and apply it to every expected review file before drafting review consensus. Add a manual-command regression test, not only runner tests.

### [MAJOR] Implementer exclusion silently weakens deliberation review-consensus quorum

PRIMARY — In an isolated `track: deliberation` idea with participants `impl` and `rev-a`, I filed the permitted reviewer artifact and only `rev-a`'s review-consensus signoff. `consensus.Status(..., review=true)` reported `ready` with no missing participants. The reproduction passed with: `go test -count=1 -v ./internal/consensus -run '^TestAuditReproDeliberationReviewConsensusDropsImplementer$'`.

SECONDARY — `internal/consensus/consensus.go:100-113` reuses `expectedRoundParticipants` for review-consensus signoffs, and `internal/consensus/consensus.go:626-650` always removes the resolved implementer without reading the track. The new test `TestReviewConsensusDoesNotAwaitTheImplementersSignoff` creates no explicit track and incorrectly generalizes the standard-track quorum to every track.

SECONDARY — Phase 6 artifact authorship and Phase 7 approval quorum are different rules: the implementer must not review its own work, while deliberation consensus still requires every active participant's signoff. The concrete fix is to split review-round authors from review-consensus voters; exclude the implementer for Phase 6, use reviewers for standard/fast Phase 7, and use all active participants for deliberation. Add explicit tests for all three tracks.

### [MAJOR] A freshly initialized deck is trapped behind an unconfirmable freshness gate

PRIMARY — In an isolated workspace I ran `protocol.InitWorkspace`, then called the production freshness classifier with `Yes:true`. It still returned classification `unknown-freshness`, one `unknown-freshness` gate, and a confirm command that merely repeats `parley preflight --yes`. The reproduction passed with: `go test -count=1 -v ./internal/app -run '^TestAuditReproFreshInitYesCannotClearUnknownFreshness$'`.

SECONDARY — Fresh initialization writes consumer metadata without either protocol hash. `internal/app/preflight.go:426-433` returns the same gate whenever a hash is absent and contains no `opts.Yes` recovery branch, unlike the adjacent role gates. The added test verifies only that the gate appears; it never runs the printed confirmation path.

SECONDARY — This corrects F24's false “in sync” claim by replacing it with a permanent block on a legitimate fresh artifact. The concrete fix is to write truthful initial hashes during `init`, or make the explicit confirmation hash the live and packaged protocols and persist the result. In either design, an end-to-end test must prove that the displayed confirmation command actually clears the gate.

## Open questions

Review in progress.
