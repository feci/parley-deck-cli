---
idea: meta-protocol-change-designated-implementer
status: implemented
implementer: kimi-1
started: 2026-09-25
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer#designated-implementer + /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill#designated-implementer
head-commit: 0893989
skill-commit: bf7e049
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

### Ratified at review cycle 1 (fix-up, 2026-09-25)

Each item below was presented for ratification in `review/consensus.md` (cycle 1, "FINAL
deviations & gap-fills presented for ratification" plus claude-1's reservations R-1/R-2/R-3) and
carries all three signoffs (kimi-1 ✅, claude-1 🟡 ACCEPT-WITH-RESERVATIONS, zcode-1 ✅).

1. **R16 wording deviation (AF-3, VC-B option (b)).** FINAL R16's literal "fire on any run,
   including a design-only run" is not implementable on the first-launch path by FINAL's own R25
   measurement (preflight at `internal/app/app.go:1922` precedes `runcontrol.Create` at `:1939`;
   `continueAuto` runs no preflight at all — claude-1 PRIMARY), so the shipped protocol sentence
   was corrected to the code: gates fire on any run that reaches an implementer, review-round,
   goal-check or fix-up action; a design-only run keeps the defect latent until the idea next runs
   with `auto_implement` on (R-3(i) wording). The substantive protection — a defective line stops
   any run before it dispatches, reviews, goal-checks or fixes up — is unchanged, and AC-1
   (identical hunks) is preserved. The unimplemented remainder of R16 (an `[A]`-tagged rule) is
   carried by the NAMED, INACTIVE follow-up idea `meta-protocol-change-designation-launch-surfacing`
   (DF-2, R-3(ii)) — recorded by name only; this consensus authorizes no work on it, and opening it
   is a post-close owner/organizer act.
2. **R29 `[D]`-tag interpretation (AF-2, VC-A).** Resolved: the rule body governs over the tag
   gloss — the comparison keys on the RECORDED source, so deleting a designation after a recorded
   dispatch escalates (the check is hoisted out of the `implDesignated` guard). No FINAL text
   changes; the interpretation is the consensus's, recorded openly.
3. **Gap-fill — pin-shadowed malformed tier-3 notice (AF-5).** R19 concession 3 states the hard
   failure without a pin exception and is silent on the pinned branch; the notice fills the silence
   fail-safe (visibility without a gate).
4. **Gap-fill — layered-config error notice (AF-6).** FINAL is silent on `LoadDefaults` failure on
   this path; the construction-time WARNING is a deliberate, disclosed R45-boundary judgment call —
   broken-config decks gain a warning line, healthy-config unset decks stay byte-identical (pinned
   by unmodified `TestUnsetPathIsByteIdentical`).
5. **Gap-fill — the reassignment exit is machine-read in the re-entry check (AF-2).** The shipped
   error message already named this exit; FINAL R29 does not specify its mechanics. Honouring it
   uses the AF-1-strict parser and changes no unset-path behaviour.
6. **AF-2 Load-error form (claude-1 R-2) — the narrower owner-compatible variant, NOT a
   deviation.** The draft's always-on `Store.Load()` fail-closed escalation was explicitly NOT
   ratified: it would have departed from the owner's boundary sentence ("Ship the mechanism with
   the global default UNSET, so behaviour without it is exactly today's"). As implemented, the
   escalation fires only when the current dispatch is genuinely designation-sourced (`implSource ∈
   {designation, global-default}` after the pin exit — exactly the set AF-4's `live` predicate
   selects there, exactly the set R29 protects), so the boundary holds by construction rather than
   by ratification. **Disclosed one-corner relaxation vs the pre-fix code (not silent):** pre-fix
   the check sat under the `present` guard, so a `none`/fall-through deck with a corrupt event
   store escalated at `Implement`; post-fix it proceeds — this aligns the failure surface with
   AF-4's designation/non-designation split (FINAL R2's table calls `none` an explicit
   non-designation), no designation-derived dispatch authority is in play on those decks, and the
   comparison itself still runs whenever Load succeeds. **Residual evidence-destruction limitation
   — recorded, not waived:** deleting the designation AND corrupting/destroying the event store
   defeats both branches (Load fails on a now-undesignated dispatch; the check returns nil; the
   reassignment proceeds with R29's durable evidence destroyed). This is inherent to a best-effort
   event file whose writes are already ignored by design (`_ = o.base.Store.Append(…)`); it
   requires destroying the run's durable record — a detectable, deck-visible act, not a one-line
   bypass; hardening past it (an always-on gate) crosses the owner's boundary and remains an owner
   decision, never a participant one. Pinned honestly by
   `TestReentryResidualDeletionPlusStoreDestructionProceeds`.
7. **Dismissed-finding documentation (zcode-1 source-only NIT — dismissal adopted the
   documentation alternative).** The re-entry comparison escalates on a source-only change (same
   id, different recorded source): R29's body makes the recorded SOURCE part of what the comparison
   protects, so a same-id tier change means the durable record no longer describes how the current
   dispatch is authorized. The owner-confirmed one-line exit is `implementer_reassigned: X to X —
   <reason> — confirmed <date>` (the same-id pair still matches the AF-1-strict parser), which is
   exactly what the escalation message prints. Comparing ids alone would let "designation deleted,
   global default set to the same id" pass unrecorded — a sibling of the class AF-2 closes.

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
- 2026-09-25 — **fix-up cycle 1 (Phase 8).** Phase-8 packet attested (full mode,
  `c7492182…568f`, `fallback_reason` absent). Applied AF-1…AF-11 exactly as amended by the
  R-1/R-2/R-3 ratification (see `## Fix-up cycle 1`): parser hardening, re-entry hoist with the
  narrower owner-compatible Load-error form, present/live split, two notices, protocol/changelog
  wording corrections across all three copies plus this idea's own changelog entry, and the skill
  companion line. Source/protocol/test commit `d238238` (CLI worktree); skill companion commit
  `a624318` (skill worktree); this file's update follows in the next commit on the same branch
  (a file cannot name its own commit). No frozen FINAL, signoff, or peer-review edits; no
  version/release/main-merge/global-config change; the product default ships UNSET.

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

### Fix-up cycle 1 evidence (at `d238238`, skill `a624318`)

- **Phase-8 packet attestation (this fix-up's authority):** `parley protocol packet --dir .
  --phase 8 --track deliberation --idea meta-protocol-change-designated-implementer --flag
  auto_implement --flag protocol_change --audience participant --json` → `context_mode: "full"`,
  `fallback_reason` ABSENT (top-level keys enumerated: `body_path`, `context_mode`, `index`,
  `packet_sha256`, `request`, `shadow`, `source`, `source_sha256`), `source_sha256` =
  `packet_sha256` = `c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f` =
  `shasum -a 256 parley-deck/COOPERATION.md` (1,400 lines / 114,771 bytes — the live amended
  authority, byte-identical to what all three signoff attestations used); `shadow` (86,336 bytes,
  40 included / 29 omitted) not used. Phases 6–8, §15, §0/§9.0 and the shipped Phase-5/§4.0/§10
  hunks were read from it.
- **Scoped suites:** `go test ./internal/protocol/ ./internal/config/ ./internal/consensus/
  -count=1` all ok; `go test ./internal/app/ -count=1 -timeout 900s` ok 534.5s (shared-mount fork
  profile); `go build ./...` + `go vet ./...` OK; `gofmt -l` clean on every file this cycle
  touched.
- **Mutation/bidirectional proofs (each reverted and re-verified green after):** AF-1 — substring
  identity mutation → 6 failing subtests; substring-"confirmed" marker mutation → 7 failing
  subtests. AF-2 — fail-closed Load branch removed →
  `TestReentryCorruptStoreFailsClosedOnlyWhenLive` leg 1 FAILs; broadened to the NOT-ratified
  always-on form → legs 2+3 and `TestReentryResidualDeletionPlusStoreDestructionProceeds` FAIL;
  hoist reverted behind `implDesignated` → `TestReentryEscalatesOnDesignationDeletion` (B1) and
  `TestReentryEscalatesOnGlobalDefaultCleared` (B2) FAIL. AF-4 — kickoff keyed on `present`
  instead of `live` → `TestNoneAndFallThroughsAreNotLive` FAILs. AF-5/AF-6 pin their
  NOTICE/WARNING text by assertion. All existing pins green unmodified:
  `TestUnsetPathIsByteIdentical` (AC-4), `TestReentryComparisonEscalatesOnDesignationChange`,
  T-2/T-4/T-5/T-6/T-7/T-9/T-10, AC-10/AC-14/AC-18, T-1/T-3/T-11/AC-15/AC-19.
- **AC-1 re-derived after the AF-3/AF-9 edits:** `tail -n +160 parley-deck/COOPERATION.md` and
  `tail -n +153` of both other copies all hash to
  `3621b7a637c5183bcf9120ff73e231e2a4503859fdec6cd6b84bf9713beef30e`; `diff` deck-vs-guarded still
  yields only the three pre-existing project-zone hunks (`3c3`, `6,7c6`, `154,159d152`).
- **AC-2 re-derived (R-1):** the corrected claims appear verbatim once per protocol copy and once
  in the changelog (`any run that reaches an implementer, review-round, goal-check or fix-up
  action`; `today's chain — FINAL.md's recorded implementer:/drafted-by:, else the first eligible
  participant`); `git diff e4640bf -- parley-deck/meta/protocol-changelog.md` stays inside this
  idea's own UNRELEASED entry (hunk `@@ -1,3 +1,30 @@`); residual-wording sweep across all four
  files is CLEAN.
- **AC-3 re-derived:** skill-worktree `grep -rn "default implementer" skills/ lib/ bin/` → 0 hits;
  no `FINAL drafter` residue in `skills/`. AC-19 template tests unaffected (internal/config green);
  the drift test over the embedded default is green (`internal/protocol` ok).

## Decision Log

- 2026-09-25 · kimi-1 — F2/F4/F5 decisions as recorded above (availability gate rides the dispatch
  actions; no preflight copy; names as listed; tier-3 malformed = whitespace-containing value).
- 2026-09-25 · kimi-1 — review cycle-1 fix-up (AF-1…AF-11 as amended by the R-1/R-2/R-3
  ratification) applied exactly as signed; decisions and disclosures are in
  `## Deviations from FINAL.md` → "Ratified at review cycle 1" and `## Fix-up cycle 1`.

## Surprises & Discoveries

(Living.)

- The deck's COOPERATION.md §0 sentence documenting `default_implementer` is itself a
  `default_implementer` string hit — any AC-19-style assertion must target config files, not
  prose (caught by the full `internal/config` run, fixed in the same session).
- On this host the suite's fork-heavy packages (`internal/app` 505s, `internal/trajectory` 605s)
  run an order of magnitude slower than on local disk; a 4-minute per-package timeout produces
  false "panic: test timed out" failures here. Both pass with `-timeout 900s`.
- **AF-11 (fix-up cycle 1):** under a live designation the always-on `agent.model_diversity`
  event is recorded TWICE per run — once at kickoff (R37's "runs unchanged … also runs at
  kickoff") and once at `OpenReviewRound` — versus once on an unset deck; a licensed consequence
  of R37, recorded here because consumers counting these events will see it. After AF-4 the
  doubling is confined to genuinely designated runs (it no longer fires on `none` or the
  fall-throughs); event counts are pinned by the AF-4 tests.
- **AF-6 test discovery (fix-up cycle 1):** the launch path has its own pre-existing fail-closed
  read of the layered config (the launch budget refuses when the deck `agents.toml` is
  unreadable), so the AF-6 notice test pins resolution + surfacing only; that refusal predates
  this delta and is not part of the designation path.

## Fix-up cycle 1 (2026-09-25, kimi-1)

Applies the signed cycle-1 fix plan (`review/consensus.md` — kimi-1 ✅, claude-1 🟡
ACCEPT-WITH-RESERVATIONS, zcode-1 ✅) exactly as amended by the R-1/R-2/R-3 ratification. Phase-8
packet attested before any edit (full mode, `fallback_reason` absent, authority `c7492182…568f`
byte-identical to the signoff attestations — see Validation evidence). Nothing else was changed:
no frozen FINAL edits, no signoff/peer-review edits, no silent design expansion.

**Fixes applied (all eleven; the count stayed `outstanding_agreed_fixes: 11`):**

- **AF-1** — `internal/protocol/implementer.go`: both confirmation-record parsers
  (`ImplementerWaived`, `ParseImplementerReassignment`) are now fail-closed via a shared
  `confirmationRecordTail`: exact subject identity (`kimi-10` no longer clears `kimi-1`; the
  `zz-impl-2`/`zz-impl` prefix collision closed), at least one non-empty reason segment, a strict
  trailing `confirmed <date>` (any casing) whose date parses as a valid calendar date (VC-D —
  `2026-02-31` rejected), and no negation (`not confirmed` / `not yet confirmed` / `unconfirmed`)
  in any earlier segment. Tests: `TestConfirmationRecordsAreFailClosed` (15 cases); the existing
  T-4/T-5 fixtures (`zz-impl — offline — confirmed 2026-09-25`, `aa-first to zz-impl — rotation —
  confirmed 2026-09-25`) pass unchanged, proving the recorded shape itself is untouched.
- **AF-2 (as amended by R-2 — the narrower owner-compatible form; the always-on form was NOT
  ratified)** — `internal/app/driver_impl.go`: `checkImplementerReentry()` is called
  unconditionally in `Implement` (hoisted out of the `implDesignated` guard, so DELETING a
  designation after a recorded dispatch escalates — VC-A: the body governs, the comparison keys on
  the recorded source); the escalation's named exit is now machine-read (a confirmed
  `implementer_reassigned: <recorded> to <current>` record parsed by the AF-1-strict parser clears
  it); a non-not-exist `Store.Load()` error escalates fail-closed ONLY under a genuinely
  designation-sourced dispatch (`implLive`, which after the pin exit is exactly `implSource ∈
  {designation, global-default}`) — unset/`none`/fall-through decks keep exactly today's
  behaviour, so the owner's unchanged-default boundary holds by construction. The one-corner
  relaxation and the residual evidence-destruction limitation are disclosed in `## Deviations from
  FINAL.md` item 6 and pinned by test. Tests: `TestReentryEscalatesOnDesignationDeletion` (B1,
  incl. negated-record and confirmed-record exits and idea scoping),
  `TestReentryEscalatesOnGlobalDefaultCleared` (B2),
  `TestReentryCorruptStoreFailsClosedOnlyWhenLive` (the R-2 three-leg set),
  `TestReentryResidualDeletionPlusStoreDestructionProceeds` (the residual pin); unmodified
  `TestUnsetPathIsByteIdentical` and `TestReentryComparisonEscalatesOnDesignationChange` stay
  green.
- **AF-3 (as amended by R-1 + R-3(i))** — all three `COOPERATION.md` copies, identical hunks:
  the validity-gate sentence now reads "fire on **any run that reaches an implementer,
  review-round, goal-check or fix-up action** … blocks that run before it dispatches" plus "A
  design-only run (`auto_implement` off) reaches none of these actions, so a defective designation
  stays latent until the idea is next run with `auto_implement` on." (the draft's "surfaces the
  defect at its first dispatching action" clause was withdrawn per R-3(i)); the §4.0 template
  comment reads "(blocked at the first role action)". R-1 hunk: this idea's own UNRELEASED
  changelog entry carries the same corrected claim (`parley-deck/meta/protocol-changelog.md:13`).
- **AF-4** — `internal/app/driver_impl.go`: `present` (R45 emission, unchanged) split from `live`
  (`present && source ∉ {none, fall-through-inapplicable, fall-through-unavailable}` — the
  boundary `designatedImplementerPreference` already applies for R36); `kickoffDesignationChecks()`
  is gated on `live`. `implementer: none`, `default_implementer = "none"`, and the waived
  fall-through keep the R45 line/event but run no kickoff checks and emit exactly one
  `agent.model_diversity` event (at `OpenReviewRound`), as on an unset deck. Tests:
  `TestNoneAndFallThroughsAreNotLive`; T-10/AC-14 unmodified and green.
- **AF-5** — `internal/app/driver_impl.go` (`DesignationAbsent`+`pinOK` branch only): a
  pin-shadowed malformed tier-3 value emits a one-line NOTICE naming the value and that it is
  ignored while the pin governs — never a gate. Test: `TestPinShadowedMalformedDefaultNoticesNotGates`;
  T-7/AC-10 (malformed tier-3 WITHOUT a pin hard-fails) unmodified and green.
- **AF-6** — `internal/app/driver_impl.go` + `internal/app/driver_consensus.go` (the third caller,
  claude-1's attribution correction adopted): `globalDefaultImplementer` carries the
  `LoadDefaults` error out; ops construction prints one `driver: WARNING` naming the config error
  with tier 3 treated as unset for this dispatch — never a gate, no event change; the R36 drafter
  preference lets a read error steer nothing. Test: `TestConfigErrorSurfacesAsNoticeNotGate`;
  `TestUnsetPathIsByteIdentical` (healthy config) pins the byte-identical side.
- **AF-7** — this file's frontmatter: `head-commit: e4640bf` → `0893989`, added
  `skill-commit: bf7e049`, `branch:` names both worktrees. The Phase-8 fix-up close bumps
  `head-commit` again per protocol; the fix-up HEADs are recorded below.
- **AF-8** — the duplicate empty `## Validation evidence` placeholder section is removed;
  `grep -c '^## Validation evidence' IMPLEMENTATION.md` returns 1.
- **AF-9 (as amended by R-1)** — all three copies, identical hunks: rank 4 is "(4) today's chain —
  `FINAL.md`'s recorded `implementer:` / `drafted-by:`, else the first eligible participant (list
  order)."; §10 TL;DR item 6 reads "then today's chain — `FINAL.md`'s recorded
  implementer/drafter, else the first eligible participant". R-1 hunk: the changelog's `:12`
  carries the same replacement. No behaviour change; R49's deferral untouched.
- **AF-10** — skill worktree `skills/parley-deck/references/ROSTER_AND_PROTOCOL.md` Phase-5 line
  (combined with AF-9's edit of the same line): restores the volunteer clause — "a claim does not
  override a live designation (**with no designation, a claim remains the normal volunteer
  route**)". No installer, manifest or channel file touched (R57's scope).
- **AF-11** — the double `agent.model_diversity` event under a live designation is recorded in
  `## Surprises & Discoveries` (licensed by R37; after AF-4 confined to genuinely designated runs).

**Deviations, gap-fills and limitations:** all recorded in `## Deviations from FINAL.md` →
"Ratified at review cycle 1" (items 1–7), including the R16 wording deviation, the R29 `[D]`-tag
interpretation, the two gap-fills, the R-2 narrower form with its one-corner relaxation and the
residual evidence-destruction limitation, and the dismissed-finding documentation for the
source-only strictness. No FINAL text was edited.

**Deferred follow-ups:** DF-1 (F6 carrier), DF-3 (organizer/owner carrier), DF-4
(`windows-portability`) unchanged; **DF-2's carrier is named** (R-3(ii)):
`meta-protocol-change-designation-launch-surfacing` — NAMED and INACTIVE; this records the name
and authorizes no work on it (no code, no idea directory, no activation; opening it is a
post-close owner/organizer act). It holds the unimplemented remainder of R16 and is
cross-referenced with FINAL register item F4.

**Exact commits:** CLI fix-up source/protocol/tests `d238238`; skill companion `a624318`; this
file's update is the next commit on the same branch (a file cannot name its own commit — no
self-hash claim). Pre-fix-up HEADs under review: CLI `0893989`, skill `bf7e049`.

**Tests this cycle:** scoped suites and mutation proofs as itemized in `## Validation evidence` →
"Fix-up cycle 1 evidence"; full `go test ./... -count=1 -timeout 900s` at `d238238` in a clean
local-disk detached checkout (`/private/tmp/parley-di-fullsuite-kimi1`, removed after) run by the
implementer SERIALLY (the only broad suite on the host at the time): **every package ok** —
longest: `internal/trajectory` 635.2s, `internal/app` 526.1s, `internal/runner` 136.1s,
`internal/budget` 90.9s (no speedup vs the shared mount was observed on this VM host; the run was
still deliberately serial and clean-checkout). `go build ./...`, `go vet ./...`, `gofmt` clean on
touched files, all in the same checkout. **AC-17 standing, honestly:** this is the IMPLEMENTER's
full-suite PASS at the fix-up HEAD (verification-plan step 1). The retained AC-17 gap is NOT
closed by it: step 2 — a NON-implementer reviewer running the whole suite serially in a clean
local-disk checkout, recorded in `review/round-02/` — remains outstanding, and AC-17's full-suite
clause may be cited as independently confirmed only after that run is green. The LE-7 goal-done
check by a fresh non-implementer remains owed after it. No completion is claimed here.

**No contradiction found:** nothing in the amended plan conflicted with the code or the owner
boundary, so no new consensus was needed; every fix is exactly the signed text.

## Outcomes & Retrospective

(At completion — not claimed here. `status: implemented` marks review-readiness only: Phase 6
review belongs to claude-1 and zcode-1, and the goal-done check belongs to a fresh
non-implementer.)
