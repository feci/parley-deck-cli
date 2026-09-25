---
idea: meta-protocol-change-designated-implementer
status: implemented
implementer: kimi-1
started: 2026-09-25
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer#designated-implementer
head-commit: e4640bf
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Phase-5 implementation of the frozen FINAL (commit `e4640bf`, unanimous signoffs) by kimi-1 —
the agreed implementer (claim: `../../inbox/kimi-1-to-all_meta-protocol-change-designated-implementer_impl-claim.md`,
filed before any Phase-5 work). The protocol gains an owner-set designation of the participant who
executes FINAL: per-idea `implementer:` in `00-prompt.md` plus optional `[defaults].default_implementer`,
one resolution chain, fail-closed validity gates, completely dormant when neither is set. All locators
re-measured at HEAD `e4640bf` (VC-4), not transcribed from FINAL.

## Protocol context attestation

`parley protocol packet --dir . --phase 5 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change --audience
participant --json`:

```json
{
  "context_mode": "full",
  "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7",
  "body_path": ".parley-runtime/protocol-packets/full-phase5-deliberation-8ce83cde….md",
  "source": { "role": "source", "transport": "github-pr", "bytes": 109928 },
  "shadow": { "packet_sha256": "271f79bb…", "packet_bytes": 70749, "included_blocks": 34, "omitted_blocks": 35 }
}
```

- `fallback_reason` is ABSENT, verified by enumerating top-level keys (`body_path`, `context_mode`,
  `index`, `packet_sha256`, `request`, `shadow`, `source`, `source_sha256`).
- Rendered body verified byte-identical to the live authority: `shasum -a 256` of the body path and
  of `parley-deck/COOPERATION.md` both `8ce83cde…9db7`, 1,386 lines / 109,928 bytes. The Phase-5
  body was read from it (Phase 4/5 sections, §0, §9.0, §10, §7).
- The `shadow` block describes the optimized packet that was not used; the full body was used.
- HEAD at implementation start: `e4640bf2840249db0c1a1ecab7493813f4dacfdb`.

## Implementation plan / checklist

- [x] **`internal/protocol/implementer.go` (new)** — `ImplementerKey`, `ImplementerWaivedKey`,
      `ImplementerReassignedKey` constants; four-state parser `ImplementerDesignationFromMeta` /
      `ReadImplementerDesignation` (Absent / Empty / None / Set; R2–R5: presence survives empty,
      `none` any casing, quoted forms trimmed, malformed never repaired); `ImplementerSource` enum
      with distinct `pin` / `designation` / `global-default` / `none` / `fall-through-unavailable` /
      `fall-through-inapplicable` (R47); shared `ResolveImplementerChain(ideaDir, candidates,
      eligible) (id, source, ok)` taking the ordered source list + eligibility list (R31);
      `LegacyImplementerCandidates()` = [IMPLEMENTATION.md{implementer}, FINAL.md{implementer,
      drafted-by}]; `PinImplementerCandidates()` = rank 1 only; `ParseImplementerReassignment` /
      `ImplementerWaived` record readers (R18/R28 shapes).
- [x] **`internal/protocol/facilitator.go`** — code comment at `FacilitatorRoleFromMeta` recording
      the exclusionary-vs-appointive divergence (R6). No behaviour change.
- [x] **`internal/config/runtime.go`** — `DefaultImplementer` on `globalDefaults`
      (`toml:"default_implementer"`) and `CentralDefaults` (R7); non-empty-string merge in
      `mergeDefaults` — NOT a pointer (R8, ALT-19 rejected); `centralDefaultTemplate` emits the key
      COMMENTED OUT with one sentence recording the deliberate deviation (R12, K-2).
- [x] **`internal/app/driver_impl.go`** — dispatch resolution in `newDriverImplOps` after `eligible`
      is computed (R13): validity gates (present-empty, id-not-in-eligible at tier 2, malformed
      tier-3) and the R28 pin-vs-designation conflict carried on the existing `roleErr` path (R25);
      tier-2 availability gate carried separately and escalated only by the dispatch actions
      `Implement`/`Fixup` (R17/R18 — see Decisions); tier-3 inapplicable/unavailable → one-line
      notice + fall-through (R19–R22); kickoff diversity check + two-participant warning under a
      designation (R37/R38, `checkModelDiversity` itself unchanged); designation-only
      `agent.implementer_resolved {idea, implementer, source}` event + one extra stdout line in
      `Implement` immediately before `RunImplementation`, guarded by `o.base.Store != (store.Store{})`
      (R45/R46); R29 re-entry comparison via `store.Load()` firing only when the recorded source was
      `designation`/`global-default`. The `driver: implementing via %s ...` line stays byte-identical;
      the app-level `resolveImplementer` keeps its exact legacy behaviour as rank 4 (R33).
- [x] **`internal/app/driver_consensus.go`** — `firstEligibleHeadlessAgent` gains a `preferNot`
      parameter (the resolved designation) preferring an eligible non-designee drafter with
      degenerate fallback to the designee (R36); `runDrafter` supplies it; `firstHeadlessAgent`
      passes `""` (unchanged behaviour).
- [x] **`internal/consensus/consensus.go`** — `resolveImplementer` reimplemented over
      `protocol.ResolveImplementerChain` with the legacy candidates and the RAW participants list
      (R32, R34); `expectedRoundParticipants` observable output unchanged (R35).
- [x] **Tests (new, participant-owned; each fails if its rule is removed)** —
      `internal/protocol/implementer_test.go`: T-1 (four states + malformed inline-comment +
      quoted forms), chain order/pin/eligibility.
      `internal/config/`: T-3 (higher non-empty wins across two real layers; `""` does not clear;
      `"none"` suppresses) + AC-19 (template emits commented-out; deck gets no generator write).
      `internal/app/driver_designation_test.go`: T-2 (`none` suppresses tier 3), T-4 (tier-2
      unavailability gate + all three exits), T-5 (pin conflict escalates; `implementer_reassigned:`
      re-pins only with a confirmed record), T-7 (all five dormant tier-3 states: absent /
      non-participant / `none` / own declared facilitator / ping-failed — no gate, fall-through,
      distinguishing source), T-8 (no event on unset deck; legacy tests unmodified), T-9 (gate
      scoping across two ideas), T-10 (two-participant kickoff warns, not blocks), AC-10 (malformed
      tier-3 hard-fails), AC-18 (end-to-end fixture dispatch: per-idea designation dispatches the
      designee; two-layer `default_implementer` dispatches the higher layer's id and records
      source), T-6 (degenerate drafter fallback).
      `internal/consensus/`: T-11/AC-15 (review exclusion reads pin, not designation;
      `ExpectedRoundParticipants` unaffected by designation in every state) + AC-4 parity test.
- [x] **Protocol text, identical hunks in all three copies** (AC-1): `parley-deck/COOPERATION.md`,
      `internal/protocol/defaults/COOPERATION.md`, skill worktree
      `skills/parley-deck/references/COOPERATION.md` — §0 `[defaults]` sentence (R12 sentence), §4.0
      Phase-0 template (`implementer:` field), §4 Phase 4 (drafter separation, R36, plus the R40
      self-containment trigger-list sentence), Phase 5 (the mechanism: four states, chain,
      gates/exits, fall-throughs, claim subordination R15, scope sentence R42, `excluded:` note
      R24, drafter==implementer concentration record R36), §10 TL;DR item 6, §9.0 readiness note.
      Plus `parley-deck/meta/protocol-changelog.md` §7 entry (AC-2).
- [x] **Skill worktree** — `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md` Phase-5 line
      (R57), same branch; then AC-3 grep over the skill tree.
- [x] **Health** — `go build ./...`, `go vet ./...`, `go test ./...`, `gofmt -l` clean on changed
      files (AC-17); inertness greps (AC-16); three-copy tail hash check (AC-1).
- [x] **Release runbook** — `implementation-release-plan-kimi-1.md` (advisory; no publication).

## Decisions FINAL leaves to the implementer (register F2/F4/F5)

- **F5 names.** File `internal/protocol/implementer.go`; parser `ImplementerDesignationFromMeta`;
  chain `ResolveImplementerChain`; source enum `ImplementerSource` with constants
  `SourceImplementerPin/SourceImplementerDesignation/SourceImplementerGlobalDefault/SourceImplementerNone/
  SourceImplementerFallThroughUnavailable/SourceImplementerFallThroughInapplicable`.
- **F2 (`--no-implement` visibility).** Trace: the flag is defined at `internal/app/app.go:1807`
  (`parley run`) and `:1159` (`parley continue`), consumed into `AutoImplement` at `:1985`, `:2039`,
  `:1253`. The gate authority is `newDriverImplOps` construction, where the flag is NOT visible —
  and changing the constructor signature would break three existing test call sites
  (`driver_impl_le_test.go:89`, `facilitator_autodrive_test.go:146`, `facilitator_roles_test.go:31`)
  that must stay unchanged. Resolution: validity gates (R16) live on `roleErr` (any run, all four
  role actions — R25); the tier-2 AVAILABILITY gate (R18) is carried on a separate field escalated
  only by the dispatch actions `Implement` and `Fixup`. A run that cannot reach Phase 5 never calls
  either, so availability can never gate it — R17's scoping holds by construction and no stall is
  created. The waiver line remains available as an operator exit, not as cover for a mis-scoped gate.
- **F4 (preflight copy).** NOT shipped. R25's measurement holds at HEAD (preflight at `app.go:1922`
  runs before `runcontrol.Create` at `:1939`), so an idea-scoped preflight copy could only cover
  re-runs of existing ideas; the `roleErr` authority is complete without it. Recorded, not built.
- **Tier-3 "invalid" definition (AC-10).** FINAL R19 concession 3 requires invalid tier-3 values to
  keep hard failure while R20/R22 make not-in-`eligible` inapplicable. The line drawn: a tier-3 value
  containing whitespace after trimming (e.g. `"kimi-1 # note"`, the R4 malformation class arriving
  via TOML) is malformed → gate; a whitespace-free id not in `eligible` → inapplicable → fall-through.

## Deviations from FINAL.md

(Living. None planned; every deviation will be recorded here with rationale.)

- **Locator note, not a deviation:** FINAL R25 cites `Complete` (`:504`) as the fourth roleErr-
  escalating role action; re-measured at `e4640bf`, `:504` is `Fixup`'s `roleErr` check
  (`Complete` at `:518` never launches an agent and has no roleErr check). The four role actions
  that escalate `roleErr` are `Implement` (`:214`), `OpenReviewRound` (`:278`), `GoalCheck`
  (`:407`), `Fixup` (`:504`). Designation gates ride exactly those four.
- **FINAL's disclosed unreviewed item (K-8):** R57's skill-companion edit
  (`ROSTER_AND_PROTOCOL.md:64`) and the three locator corrections were not peer-reviewed before
  freezing. They are implemented as written and flagged here for reviewers.

## Notes for reviewers

- The unset path (neither `implementer:` nor `default_implementer` set) is byte-identical: no event,
  no extra stdout line, same dispatch, same review expectations (AC-4; pinned by unmodified
  `app_test.go` `TestResolveImplementerFromRoleMetadata`, `roundgate_test.go`
  `TestUnresolvableImplementerExpectsEveryone` + `TestFinalDrafterIsTheFallbackImplementer`).
- Reviewers: claude-1 and zcode-1 (non-implementers). The goal-done check belongs to a fresh
  non-implementer; this file claims no completion.
- `internal/tui/protosnap.go`, pipeline-block dispatch, and `excluded:` cross-checks are FINAL's
  deferred register items (F3/F10) and are untouched.

## Progress

- 2026-09-25 — claim filed (`inbox/kimi-1-to-all_..._impl-claim.md`); phase-5 packet attested
  (full mode, byte-identical authority); all code locators re-measured at `e4640bf`; this plan
  written before any source edit.
- 2026-09-25 — (first process, timed out at 60 min) all Go code landed: shared parser/chain in
  `internal/protocol/implementer.go`, `default_implementer` in `internal/config/runtime.go`,
  driver wiring in `driver_impl.go`/`driver_consensus.go`, consensus resolver swap in
  `consensus.go`; all new tests (T-1…T-11 + AC pins) written and green individually; protocol
  hunks H1–H3 (§0 defaults sentence, §4.0 template field, §4 Phase-4 drafter separation) applied
  identically to all three COOPERATION.md copies.
- 2026-09-25 — (continuation process) reconciled the timed-out state (no surviving child
  test/build processes; organizer's codex process left alone; stale driver-run dirs under
  `parley-deck/runs/` left untouched). Applied remaining hunks identically to all three copies:
  R40 self-containment trigger sentence, the Phase-5 mechanism paragraph pair (chain, four
  states, gates/exits, tier-3 fall-throughs, R28 escalation, R15 claim subordination, R42 scope
  sentence, R24 `excluded:` note, R36 concentration record), §10 TL;DR item 6, §9.0 readiness
  sub-bullet. Changelog §7-shape entry prepended. Skill companion `ROSTER_AND_PROTOCOL.md:64`
  Phase-5 line updated (R57). Advisory release runbook written. Fresh phase-5 packet attested
  against the amended authority (see Validation evidence).

## Validation evidence

(Living — hand-recorded command evidence; this idea sets no `checks:` list contract.)

- **AC-1 three-copy fidelity (after all hunks):** `tail -n +160 parley-deck/COOPERATION.md` and
  `tail -n +153` of both other copies all hash to
  `da9704d57888efe549a2c781ba3bcdf34012c3caa81de9977fff3a1bc7e8eaf7`; `diff` deck-vs-guarded shows
  only the three pre-existing project-zone hunks (`3c3`, `6,7c6`, `154,159d152`).
- **AC-3 skill grep:** in the skill worktree, `grep -rn "default implementer" skills/ lib/ bin/`
  → 0 hits; remaining `FINAL drafter` hits are the guarded copy's amended text and the updated
  ROSTER line, both consistent. `SKILL.md` restates no Phase-5 rule (only the facilitator
  exclusion at :121).
- **AC-16 inertness:** `grep -rn "impl-claim" --include="*.go" .` → 0 hits. `00-prompt.md` writes
  remain only `CreateIdeaFull` (`internal/protocol/workspace.go:225`), `updateIdeaStatus`
  (`internal/consensus/consensus.go:884`), `pipeline.SeedBlockPrompt`
  (`internal/pipeline/executor.go:363`) and the pre-existing pipeline candidate seeder
  (`internal/app/pipeline_cmd.go:1024`); the delta adds no writer (diff scan: only
  `ReadFrontmatter` reads and markdown text).
- **Unit-level evidence (first process, before timeout):** `go test ./internal/config/`
  (T-3 + AC-19, 3 tests PASS), `go test ./internal/app/` designation suite (T-2, T-4, T-5, T-6,
  T-7, T-8, T-9, T-10, AC-9/10/14/18 + `TestUnsetPathIsByteIdentical` — all PASS after the
  fixture authority fix via `protocol.InitWorkspace`), `go test ./internal/consensus/` (T-11,
  AC-15, AC-4 parity PASS); `gofmt -l` clean on new/changed Go files.
- **Post-edit phase-5 packet attestation (live, not refused):**
  `parley protocol packet --dir . --phase 5 --track deliberation --idea
  meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change
  --audience participant --json` → `context_mode: "full"`, `fallback_reason` ABSENT (top-level
  keys enumerated: `body_path`, `context_mode`, `index`, `packet_sha256`, `request`, `shadow`,
  `source`, `source_sha256`), `source_sha256` = `packet_sha256` =
  `c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f` = `shasum -a 256
  parley-deck/COOPERATION.md` (the amended authority, 114,771 bytes); `shadow` describes the
  unused optimized packet (74,769 bytes, 34 included / 35 omitted). The implementation-start
  attestation above (sha `8ce83cde…`) recorded the pre-edit authority; both were full-mode with
  no refusal.
- **Whole-tree health (AC-17):** `go build ./...` OK; `go vet ./...` OK; `go test ./...` — every
  package `ok`, with two environmental long-runners: `internal/app` ok in 505.6s and
  `internal/trajectory` ok in 605.6s (both fork `git` plumbing heavily and this repo sits on the
  `/Volumes/My Shared Files` mount; `internal/trajectory` does not depend on any package this
  delta touches — `go list -deps` shows only fsutil/telemetry/budget/procctl/evidence — so its
  slowness is provably independent of the change). `gofmt -l` clean on every changed file (the
  five files `gofmt -l` still lists — `facilitator_test.go`, `protocol_test.go`,
  `driver/checks_test.go`, `procctl_test.go`, `procctl_windows.go` — are pre-existing and
  untouched by this delta).
- **Fix recorded during continuation:** `TestDeckCreatedByToolingHasNoDefaultImplementer`
  (AC-19) over-asserted by walking the whole generated deck; once the §0 protocol sentence
  legitimately documented `default_implementer`, the walk tripped on prose. Rescoped to deck
  CONFIG files (`*.toml`): a created deck carries no TOML with the key at all. Full
  `internal/config` package green after the fix.

## Decision Log

- 2026-09-25 · kimi-1 — F2/F4/F5 decisions as recorded above (availability gate rides the dispatch
  actions; no preflight copy; names as listed; tier-3 malformed = whitespace-containing value).

## Surprises & Discoveries

(Living.)

- The deck's COOPERATION.md §0 sentence documenting `default_implementer` is itself a
  `default_implementer` string hit — any AC-19-style assertion must target config files, not
  prose (caught by the full `internal/config` run, fixed in the same session).
- On this host the suite's fork-heavy packages (`internal/app` 505s, `internal/trajectory` 605s)
  run an order of magnitude slower than on local disk; a 4-minute per-package timeout produces
  false "panic: test timed out" failures here. Both pass with `-timeout 900s`.

## Validation evidence

(Living — filled as checks run; driver-populated only under a `checks:` list contract, which this
idea does not set, so entries here are hand-recorded command evidence.)

## Outcomes & Retrospective

(At completion — not claimed here. `status: implemented` marks review-readiness only: Phase 6
review belongs to claude-1 and zcode-1, and the goal-done check belongs to a fresh
non-implementer.)
