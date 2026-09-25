---
idea: meta-protocol-change-designated-implementer
goal-check: LE-7
agent: zcode-1
invocation: fresh (not a resumed review, not a fourth identity)
date: 2026-09-25
cli-head: 804522c778140ad5cc0da5658c23c7014514ecba
code-tree: 717f3debe3aabb0f7a702d689e8d10de1e04fa8c
skill-tree: a624318dcda02c47ecaa859d987efee08dba3104
record-commit: ee8849c9ba6f75748fb471fe9c910bf362300e8d
verdict: PASS
---

# LE-7 goal-done check — zcode-1, fresh invocation

I am zcode-1, an existing non-implementer quorum participant of this idea (one of its two
reviewers), launched fresh solely for this check. codex-1 remains pure organizer; kimi-1 remains
sole implementer. Under LE-7/LE-11 this check **can only withhold a close, never establish one**;
a PASS below is a criterion verification, not a completion act. `status: complete` on
`IMPLEMENTATION.md` remains kimi-1's later act; every release gate stays downstream and unclaimed.

## Close precondition — verified before proceeding

`review/consensus.md` (cycle 4, sha256 `1fa3c1c8bab1975c6d624b7efd2d32e5c2e79e142d41e017a79382973f91e75a`,
read in full): signoffs **claude-1 ✅ ACCEPT, kimi-1 ✅ ACCEPT, zcode-1 ✅ ACCEPT**;
`outstanding_agreed_fixes: 0`; `blocked: false`. The one open item is NIT-1, **dismissed on the
record as an accepted residual** with all three participants concurring — and I verified the
dismissal is honest: `IMPLEMENTATION.md:325-326` still carries the blanket "All shadow byte
figures in this bullet are rendered WITH …" lead-in, no "Except where stated" repair exists
anywhere in the record, and the on-disk record is byte-identical to the `ee8849c` blob. **No
hidden repair; the residual stands exactly as signed.**

## Live protocol context — attested

`parley protocol packet --dir . --phase 8 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change
--audience participant --json`, run by me this session (parley 1.49.1, exit 0):
`context_mode: full`; `fallback_reason` ABSENT (top-level keys enumerated: `body_path,
context_mode, index, packet_sha256, request, shadow, source, source_sha256`);
`request` = phase **8**, track **deliberation**, flags **["auto_implement","protocol_change"]**,
audience participant, optimize false, transport github-pr; `source_sha256 = packet_sha256 =
b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388` = live
`parley-deck/COOPERATION.md` (115,166 B / 1,400 lines; body re-hashed by me, byte-identical);
shadow 86,716 B / 40-29 (`fe0e4c04…`) not used — matching every cycle-3/4 attestation. Phase 8
and the LE-7/LE-11 close-integrity rule were read from this body: a textual verdict never
substitutes for current-tree criterion evidence; no self-issued verdict, stale tree, skipped or
no-execution report may close.

## Current tree — hashes and identity (all PRIMARY, this session)

- CLI HEAD `804522c`; `git diff --name-only 717f3de HEAD -- internal/ cmd/` is **empty** — the
  code tree at HEAD is byte-identical to the reviewed cycle-3 source `717f3de`; the four later
  commits are organizer artifact-preservation only (`717f3de..ee8849c` = `IMPLEMENTATION.md`
  alone; `357d1b5`/`645faba`/`804522c` = review artifacts only).
- `git diff 1bad263 717f3de -- internal/ cmd/` is exactly **one comment line**
  (`driver_designation_test.go:448`, AF-17) — re-verified by me.
- Skill worktree clean at full SHA `a624318dcda02c47ecaa859d987efee08dba3104`.
- Frozen `FINAL.md` sha256 `6e4db4734a960e7b6cad69684ea3d73884f5dfa1931cf16e8836495058f79bf1`
  (0 diff vs `e4640bf`); `00-prompt.md` 0 diff vs `e4640bf` (`status: final`, `auto_implement:
  true`, `track: deliberation`, no `strict_gate`); on-disk `IMPLEMENTATION.md` sha256
  `86fd4cb58a5b4f61fd94b46efd61ae488c2db13030f641aac463ae85150c44ba` = the `ee8849c` blob.
- Round-04 artifact hashes match the consensus header: claude-1 `3947eb4d…7f27`, zcode-1
  `f3632b97…58ff`.

## Criterion check — every FINAL observable acceptance criterion, current tree

PRIMARY = my own command this session at the current tree. PRIOR = recorded evidence, cited with
exact-commit scope; I re-read each record before relying on it.

| AC | Result | Evidence (provenance) |
|---|---|---|
| AC-1 three-copy fidelity | **PASS** | PRIMARY: `tail -n +160` deck + `tail -n +153` embedded/skill all hash `3621b7a637c5183bcf9120ff73e231e2a4503859fdec6cd6b84bf9713beef30e`; deck-vs-embedded diff = exactly `3c3`, `6,7c6`, `154,159d152`; deck-vs-skill = the same three zones. |
| AC-2 changelog | **PASS** | PRIMARY: head entry dated 2026-09-25 with `Idea:`, `Drafted by:`, `Summary:`, marked UNRELEASED. |
| AC-3 skill companion | **PASS** | PRIMARY (skill worktree, clean at `a624318`): `grep -rn "default implementer" skills/ lib/ bin/` → 0 hits; `ROSTER_AND_PROTOCOL.md` Phase-5 line states the amended chain; no contradicting restatement. |
| AC-4 unset-path invariance | **PASS** | PRIMARY: the three cited legacy tests pass **unmodified** (`app_test.go:1557`, `roundgate_test.go:102`/`:111` — both test files 0-diff vs baseline `e4640bf`); `TestUnsetPathIsByteIdentical` and `TestExpectedRoundParticipantsUnaffectedByDesignation` pass live. |
| AC-5 four-state parse | **PASS** | PRIMARY: `TestImplementerDesignationFourStates` passes; test body read — four states, `none` any casing/padded, both quoted forms accepted, R4 malformed inline-comment value kept literal (never trimmed); downstream hard failure pinned by `TestMalformedTier3HardFails`/`TestMalformedTier2Gates`, both pass live. |
| AC-6 `none` opt-out | **PASS** | PRIMARY: `TestNoneSuppressesTier3` passes (distinguishes `none` from absent). Mutation pin is PRIOR at `d238238` (round-01). |
| AC-7 config precedence | **PASS** | PRIMARY: `TestDefaultImplementerLayerPrecedence` passes (two real layers; `""` does not clear; `"none"` suppresses); `mergeDefaults` non-empty-string pattern read at `runtime.go:557`. |
| AC-8 tier-2 gate + exits | **PASS** | PRIMARY: `TestTier2UnavailabilityGateAndExits` passes (gate + all three exits); fail-closed record parsers pinned by `TestConfirmationRecordsAreFailClosed` (pass). |
| AC-9 dormant tier-3 paths | **PASS** | PRIMARY: `TestDormantTier3StatesFallThrough` passes — all five states, no gate, fall-through, distinguishing source; absent emits nothing. |
| AC-10 invalid tier-3 fails | **PASS** | PRIMARY: `TestMalformedTier3HardFails` passes (malformed keeps hard failure, no fall-through). |
| AC-11 pin conflict escalates | **PASS** | PRIMARY: `TestPinDesignationConflictEscalates` + the five `TestReentry*` tests pass (escalation; re-pin only with a confirmed record). |
| AC-12 degenerate drafter fallback | **PASS** | PRIMARY: `TestDrafterSeparationPreference` passes. |
| AC-13 gate scoping | **PASS** | PRIMARY: `TestDesignationGateIsScopedToTheIdea` passes; `preflight.go` is not in the delta file list — no all-ideas iteration added. |
| AC-14 two-participant warning | **PASS** | PRIMARY: `TestTwoParticipantDesignatedKickoffWarnsNotBlocks` passes (warns, proceeds). |
| AC-15 exclusion reads record | **PASS** | PRIMARY: `TestReviewExclusionReadsThePinNotTheDesignation` and `TestExpectedRoundParticipantsUnaffectedByDesignation` pass; consensus side passes legacy candidates + raw list to the shared chain (code read). |
| AC-16 inertness | **PASS** | PRIMARY: `grep -rn "impl-claim" --include="*.go" .` → 0 hits; delta added-line scan shows `00-prompt.md` writes only in test fixtures — production writers remain `CreateIdeaFull` (`workspace.go:225`), `updateIdeaStatus` (`consensus.go:884`), `pipeline.SeedBlockPrompt` (`executor.go:302`) plus the pre-existing pipeline seeder (`pipeline_cmd.go:1024`, outside the delta); no production file in the delta adds a write call. |
| AC-17 whole-tree health | **PASS** (scoped) | PRIMARY: `go build ./...` exit 0; `go vet ./...` exit 0; `gofmt -l` clean on all ten changed Go files (tree-wide `gofmt -l` lists exactly the five pre-existing untouched files); `go test ./internal/protocol/ ./internal/config/ ./internal/consensus/ -count=1` all ok; 22-test focused designation suite in `internal/app` all PASS. Whole-suite clause: PRIOR — claude-1's independent full suite at exactly `d238238` (durable log read by me this session: build/vet/`go test ./...` all exit 0, every package ok, `dirty: []`) and kimi-1's serial full suite at `d238238`. **No broad suite has been run at `1bad263`, `717f3de` or HEAD, and none is claimed** — the code delta since `d238238` is one regex line plus comments, under the ratified proportional plan; see Limitations. |
| AC-18 end-to-end designation | **PASS** | PRIMARY: `TestTier2DesignationDispatchesDesignee` and `TestTwoLayerGlobalDefaultDispatchesHigherLayer` pass (fixture dispatch, event payload, higher-layer source recorded). |
| AC-19 ships UNSET | **PASS** | PRIMARY: `runtime.go:654` emits the key commented out with the deviation sentence; `EnsureCentralDefault` writes central-only (no deck generator); `TestCentralDefaultTemplateShipsImplementerUnset` + `TestDeckCreatedByToolingHasNoDefaultImplementer` pass live. |
| AC-20 attended boundaries | **PASS** | PRIMARY (structural): delta added lines contain no pty allocation, publication or channel action (only substring false positives — "pty" inside "empty"); zero `exec.Command` additions; the publish TTY gate (`internal/app/protocol.go`, "requires a controlling terminal") untouched — no commit since baseline touches it. |
| AC-21 independent review happened | **PASS** | PRIMARY: `review/round-01/claude-1.md` (`:331`) and `review/round-01/zcode-1.md` (`:25`) both carry populated `## Refutation attempts` sections (read live); the goal-done check is this fresh non-implementer invocation itself. Neither clause is satisfiable by implementer or organizer; neither was. |

## Prior-evidence provenance (exact-commit scoped, none re-labelled)

- **Full independent suite** — belongs to `d238238` (claude-1, non-implementer; log
  `/tmp/parley-claude1-r2-1790307993/full-suite.log`, read by me now; corroborated in round-02 by
  zcode-1) and kimi-1's implementer serial run at `d238238`.
- **Focused behaviour + mutation/adversarial evidence** (M1=9 / M2=5 regex mutations, reentry
  hoist/present-vs-live mutations, probe sets) — belongs to `1bad263` (cycle-2 record, both
  reviewers).
- **Comment-only AST-identity and focused evidence** — belongs to `717f3de` (claude-1 round-04
  AST hash; both reviewers' focused runs), plus my own live one-line-diff verification.
- Nothing broader is claimed for any later commit.

## Limitations — named, not waived

1. **No broad `go test ./...` exists at `1bad263`, `717f3de` or HEAD.** AC-17's whole-suite clause
   rests on the exactly-scoped `d238238` prior evidence under the ratified proportionality. My
   current-tree execution is meaningful but scoped: three full packages, the 22-test designation
   suite, build/vet/gofmt. I judged no new concern requires a broad re-run: the entire code delta
   since `d238238` is one regex line (whose 30-case fail-closed test passes live at HEAD) plus
   comment lines, and every designation pin is green at HEAD.
2. Skill npm-level gates (npm test, manifest `--check`, `npm pack --dry-run`, installer
   integrity) are unexercised by anyone — carried blind spot, release-preflight territory.
3. No end-to-end `parley run` against a real designated deck beyond the fixture tests; no live
   §9.0 ping behind `designeeAvailable` (tests fake it); TUI and pipeline-block surfaces remain
   FINAL deferrals F3/F10 — named unexamined, no defect claimed.
4. NIT-1 stands as an unrepaired accepted residual (blanket lead-in at `IMPLEMENTATION.md:325-326`)
   — dismissed with all three concurrences; no acceptance criterion depends on it.
5. The pre-edit 86,336 B shadow figure is not absolutely re-derivable (deck state gone) — carried
   limitation (a) of the signed consensus.
6. AC-5…AC-15, AC-18 are verified by live re-execution of their pinning tests at HEAD, not by
   fresh adversarial re-derivation of each; the adversarial depth (mutations) is prior evidence at
   `d238238`/`1bad263`.

## Verdict

**PASS.** Independently, on my own current-tree evidence plus exactly-scoped prior evidence under
the ratified proportionality, every one of FINAL's twenty-one observable acceptance criteria holds
against the current technical tree (`804522c`, code tree `717f3de`, skill `a624318`, record
`ee8849c`), with the three-signature zero-fix closing consensus intact and its one recorded
dismissal standing unrepaired. I found no missing criterion, no stale-tree citation, no
self-issued verdict and no hidden repair. Per LE-7 this PASS withholds nothing; it does not by
itself complete anything. Release gates (origin/main integration, preflight, publication,
audits, post-release owner configuration) remain downstream, unperformed and unclaimed; the
product default ships UNSET.

## What this invocation did and did not do

Did: read 00-prompt.md (all owner directions), frozen FINAL.md, IMPLEMENTATION.md, the complete
cycle-4 closing consensus and the round-01/round-04 review evidence; obtained the live phase-8
protocol packet (full mode, byte-verified); ran the current-tree checks and test suites recorded
above; wrote this one file. Did not: publish, version, install, merge, commit, change globals or
roster, modify any source/record/signature/frozen artifact, mark anything complete, or launch any
other agent. This file is the canonical and only output; nothing was committed.
