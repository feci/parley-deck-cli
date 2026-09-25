---
idea: meta-protocol-change-designated-implementer
status: complete
implementer: kimi-1
started: 2026-09-25
completed: 2026-09-25
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer#designated-implementer + /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill#designated-implementer
head-commit: 717f3de
skill-commit: a624318
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
- 2026-09-25 — **fix-up cycle 2 (Phase 8).** Phase-8 packet attested pre-edit (full mode,
  `b273af1e…f388`, `fallback_reason` absent — see `## Fix-up cycle 2`). Applied AF-12…AF-15 as
  signed, with both reviewer-carried corrections (AF-13's accurate `:42`/`driverImplOps` locator;
  AF-15 alternative (a), accepted and recorded in `## Fix-up cycle 2`). Source/tests commit
  `1bad263`; this file's update follows in the next commit on the same branch. No frozen FINAL,
  signoff, or peer-review edits; no skill delta; the product default ships UNSET.

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

- **Phase-8 packet attestation (PRE-EDIT — attested at `0893989` before the cycle's protocol
  hunks; relabelled at fix-up cycle 2 per AF-14(1), the original "this fix-up's authority"
  wording wrongly placed the pre-edit figures under `d238238`):** `parley protocol packet --dir .
  --phase 8 --track deliberation --idea meta-protocol-change-designated-implementer --flag
  auto_implement --flag protocol_change --audience participant --json` → `context_mode: "full"`,
  `fallback_reason` ABSENT (top-level keys enumerated: `body_path`, `context_mode`, `index`,
  `packet_sha256`, `request`, `shadow`, `source`, `source_sha256`), `source_sha256` =
  `packet_sha256` = `c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f` =
  `git show 0893989:parley-deck/COOPERATION.md` (1,400 lines / 114,771 bytes — the pre-edit
  authority, byte-identical to what all three cycle-1 signoff attestations used); phase-8
  `shadow` (86,336 bytes, 40 included / 29 omitted) not used. Phases 6–8, §15, §0/§9.0 and the
  shipped Phase-5/§4.0/§10 hunks were read from it. The POST-edit authority at `d238238` is
  `b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388` (1,400 lines / 115,166
  bytes — re-derived at cycle 2 from `git show d238238:parley-deck/COOPERATION.md`; the live tree
  matches). **All shadow byte figures in this bullet are rendered WITH `--flag auto_implement
  --flag protocol_change`**: the phase-8 shadow over the `d238238` source measures 86,716 bytes /
  40 included / 29 omitted (`fe0e4c04…ebcf`, measured twice at fix-up cycle 2 with parley 1.49.1),
  while the SAME phase-8 packet over the SAME `d238238` source rendered WITHOUT the two flags
  measures 86,701 bytes / 40 included / 29 omitted (`a0845614…f08efd`) — the cycle-2 plan's
  "86,701" figure is that flagless rendering, and it reproduces on demand (settled at fix-up
  cycle 3 by all three participants' own runs; VC-3.2 in `review/consensus.md`). A PHASE-7 packet
  over the same `d238238` source, also rendered with both flags, reports a different shadow
  (80,799 bytes / 41 included / 28 omitted, `c6d29141…d980` — both cycle-2 signoff attestations).
  Read the phase and the flags with the figure.
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
- 2026-09-25 · kimi-1 — fix-up cycle 2: accepted AF-15's carried correction (a)
  (`un[\s-]*confirmed` completing the signed AF-1 clause (iv) enumeration) with rationale
  recorded in `## Fix-up cycle 2` before applying it; fallback (b) not invoked. AF-13's accurate
  `:42`/`driverImplOps` locator carried transparently (the signed plan body is untouched).

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
  doubling is confined to `live()`-true dispatches and never fires on `none` or the
  fall-throughs (read with zcode-1's round-02 NIT-3 corner, recorded via AF-14(2): `live()` is
  also true for a pin-source dispatch when tier 3 is set — `present: true`, `source: pin` — so
  the kickoff `agent.model_diversity` emission can fire there too); event counts are pinned by
  the AF-4 tests.
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

**Cycle-1 record metadata (repaired openly at the fix-up cycle-2 publication, per AF-12):**
`status: fix-up-cycle-1`; `head-commit: d238238`; `skill-commit: a624318`; `record-commit:
e4d868a` — all PRIOR commits at repair time, so the cycle-1 self-hash gap is closed
retroactively with no self-naming. Admission: the per-publication frontmatter bump owed at
`e4d868a` (Phase 8, packet body :622 — "They also update the top-level frontmatter: bump
`status:` to `fix-up-cycle-N`, update `head-commit:`") was missed there; this section's "the
Phase-8 fix-up close will bump `head-commit` again" rationale conflated the per-cycle close with
the idea close. The bump is applied at this cycle-2 publication, not deferred again to idea
close.

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

## Fix-up cycle 2 (2026-09-25, kimi-1)

status: complete
completed: 2026-09-25
head-commit: 1bad263
skill-commit: a624318
record-commit: b850576

Commit-provenance note: the cycle-2 source commit was first committed as `bf3336d` with a
`[kimi-1]` prefix; per the owner's standing prefix override (inbox
`codex-1-to-kimi-1_meta-protocol-change-designated-implementer_commit-prefix.md`, 2026-09-25) its
message was corrected to the `[codex-1]` task prefix with `(authored by kimi-1)` retained — a
message-only rewrite: the tree is byte-identical, so every evidence item in this section ran
against exactly the committed tree. `1bad263` supersedes `bf3336d`; no test evidence is claimed
to have run at a different content hash.

Applies the signed cycle-2 fix plan (`review/consensus.md` — kimi-1 ✅ ACCEPT, claude-1 🟡
ACCEPT-WITH-RESERVATIONS, zcode-1 🟡 ACCEPT-WITH-RESERVATIONS; no BLOCK) with both carried
corrections, each explicitly nonblocking and reviewer-accepted. Phase-8 packet attested BEFORE
any edit (2026-09-25T04:44Z, parley 1.49.1, this worktree at pre-edit source commit `8d026d4` —
organizer-only commits since `d238238`; the code tree was and remains content-identical to the
reviewed `d238238`): `parley protocol packet --dir . --phase 8 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change
--audience participant --json`, exit 0 → `context_mode: "full"`, `fallback_reason` ABSENT
(top-level keys enumerated: `body_path`, `context_mode`, `index`, `packet_sha256`, `request`,
`shadow`, `source`, `source_sha256`), `source_sha256` = `packet_sha256` =
`b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388` = `shasum -a 256
parley-deck/COOPERATION.md` (1,400 lines / 115,166 bytes — the live post-cycle-1 authority at
`d238238`, unchanged by this cycle); phase-8 `shadow` `fe0e4c04…ebcf` (86,716 bytes, 40
included / 29 omitted) not used — re-run for determinism with identical figures. Phase 8 was
read from the body (the frontmatter-bump sentence at body :622 verbatim).

### Fixes applied

- **AF-12** — this file's frontmatter at this publication: `status: implemented` →
  `fix-up-cycle-2`; `head-commit: 0893989` → `1bad263` (the cycle-2 source commit — PRIOR and
  knowable under the source-then-record order); `skill-commit: bf7e049` → `a624318` (still the
  skill HEAD; cycle 2 has no skill delta). The cycle-1 omission is repaired openly inside
  `## Fix-up cycle 1` (metadata lines plus admission), at this publication — not deferred again
  to idea close.
- **AF-13** — `internal/app/driver_impl.go`, comments only, no behaviour: `:148`'s gate field
  now states the role-action scope ("roleErr path, escalated by every role action
  (R16/R25/R28)") and `:188` now reads "Hard gate on any run that reaches a role action";
  `grep -n "any run" internal/app/driver_impl.go` shows no unqualified claim. **Carried
  correction, recorded here without touching the signed plan (claude-1's Reservation 1
  self-correction, zcode-1 concurring):** the accurate scope model is the `roleErr` field comment
  at `:42` in the `driverImplOps` struct ("declared-facilitator role deadlock; every role action
  escalates") — not `:45` (the "All five stay zero on an undesignated deck" comment) and not
  "eleven lines above" (`:42` sits 106 lines above `:148`, in a different struct than
  `dispatchDesignation`). The signed plan's locator sentence was wrong — including inside a
  DRAFTER-PRIMARY-tagged clause, which is hereby narrowed: the two target lines were correctly
  re-read; the reference model was propagated unchecked. The fix direction was right, and `:42`'s
  wording is the model both rewrites follow.
- **AF-14** — record prose precisions in this record commit: (1) the cycle-1 Phase-8 attestation
  bullet under `## Validation evidence` is relabelled PRE-EDIT at `0893989`, with the POST-edit
  `d238238` authority (`b273af1e…f388`, 115,166 bytes / 1,400 lines) named, and all shadow byte
  figures labelled with their phase (phase-8: 86,336 B over the pre-edit source, 86,716 B over
  `d238238`; the phase-7 shadow over the same `d238238` source is 80,799 B per both signoff
  attestations) so later phase-7 readers do not mistake different packet sizes; (2) the AF-11
  Surprises & Discoveries sentence now carries the pin-source `live()` corner (zcode-1 NIT-3).
- **AF-15 — as widened by carried correction (a)** — `internal/protocol/implementer.go:185`, one
  regex line plus its comment: `negationMarker` is now
  `(?i)\b(?:un[\s-]*confirmed|not(?:[\s-]+yet)?[\s-]+confirmed)\b`. The signed AF-1 clause (iv)
  enumeration (`not confirmed`, `not yet confirmed`, `unconfirmed` — any casing) now fails closed
  on hyphenated/space-separated spellings; this completes the signed enumeration and is NOT a
  broader free-prose denylist — the bare-"not" extension stays dismissed, and zcode-1's NIT-2(a)
  multi-segment tolerance is untouched. **My acceptance of alternative (a), recorded before it
  was applied:** I accept claude-1's offered alternative (a) — widening the bare `unconfirmed`
  alternative to `un[\s-]*confirmed` — independently elected by zcode-1, because `unconfirmed` is
  enumerated verbatim in the signed AF-1 clause (iv) (no new banned word, no new enumeration),
  and both reviewers' independent probes show zero false positives — re-verified by my own read
  of the regex (`\b` blocks `un` inside `run`/`fun` and `not` inside `cannot`; the trailing `\b`
  blocks `unconfirmedness`; the em-dash is outside `[\s-]`). Fallback (b) (explicit documented
  tolerance) is NOT invoked. Tests (`internal/protocol/implementer_test.go`,
  `TestConfirmationRecordsAreFailClosed`): 14 new cases across BOTH parsers — rejects:
  `not-yet-confirmed`, `not-confirmed`, `un-confirmed`, `un confirmed` (both parsers) plus cased
  `UN-Confirmed` (waiver); boundary accepts: `cannot-confirmed`, `unconfirmedness`, `run
  confirmed the fix`, em-dash `reason—detail`. The existing 15 cases and the T-4/T-5 fixtures are
  unmodified and green.

### Deviations from agreed fixes

- **AF-15's regex is wider than the signed plan's literal text** — exactly the carried
  correction both reviewers signed for (claude-1 offered (a)/(b); zcode-1 elected (a), (b) as
  fallback only; my acceptance recorded above). No new banned word; the enumeration completed is
  the signed one.
- **AF-13's locator carried transparently** (see Fixes applied): the signed plan's `:45` /
  "eleven lines above it" reference is wrong; the fix targets and direction were right and are
  what changed. The signed consensus body and all three signatures are preserved unedited — the
  correction lives here and in the source commit message.
- **One plan figure differs by invocation flags, disclosed:** AF-14(1) names an 86,701-byte
  phase-8 shadow for the `d238238` source — that is the FLAGLESS phase-8 rendering, and it
  reproduces on demand (`a0845614…f08efd`, 40 included / 29 omitted; confirmed at fix-up cycle 3
  by all three participants' own runs, VC-3.2 in `review/consensus.md`). The relabelled bullet
  carries 86,716 bytes (`fe0e4c04…ebcf`) because the flagged command — `--flag auto_implement
  --flag protocol_change` — is what the plan and both cycle-2 signoffs attested with. The source
  hash, line count and byte count are identical on both invocations.
- Otherwise none: no item was infeasible; no protocol text, skill file, FINAL, consensus, or
  peer artifact was edited; no new owner question arises.

### Evidence this cycle (proportional, as ratified — no duplicate broad suite absent new concerns)

- `go build ./...` OK; `go vet ./...` OK; `gofmt -l` clean on all three touched Go files.
- `go test ./internal/protocol/ -count=1` ok — the full package;
  `TestConfirmationRecordsAreFailClosed` 30/30 PASS lines (29 subtests: 15 prior + 14 new).
- Valid-fixture boundary across the parsers' other consumer: `go test ./internal/app/ -run
  'TestTier2UnavailabilityGateAndExits|TestPinDesignationConflictEscalates|TestUnsetPathIsByteIdentical'
  -count=1` ok — the T-4/T-5 fixtures and the AC-4 byte-identical pin, unmodified.
- Mutation checks (each applied, run, restored, and re-verified green): **M1** — full
  separator-class revert to the pre-AF-15 regex → exactly the 9 new reject subtests FAIL;
  **M2** — un-branch-only revert (the signed plan's original shape without correction (a)) →
  exactly the 5 (a)-specific subtests FAIL. No boundary accept flipped under either mutation.
- **AC-17 standing, honestly:** claude-1's independent full suite at `d238238` (clean local-disk
  detached checkout, `git status --porcelain` empty, `go build ./...` / `go vet ./...` / `go test
  ./... -count=1 -timeout 25m` all exit 0, durable log
  `/tmp/parley-claude1-r2-1790307993/full-suite.log`, corroborated by zcode-1's own read of that
  log) stands EXACT-COMMIT scoped to `d238238` and is NOT re-claimed for `1bad263`; this cycle's
  evidence basis is the ratified proportional plan (the only behaviour change is AF-15's regex
  line in `internal/protocol`). The LE-7 fresh-invocation goal-done check (AC-21's second clause)
  remains separate and still owed. No completion, goal-done, or release is claimed here.

### Review-pointer

Round-03 re-review reassesses exactly: AF-12's machine effect (digest/wait output) and metadata
fields; AF-13's comment accuracy; AF-14's labels; AF-15's regex scope (enumerated negations
only) and its tests; the unmodified AC-4 pins — per the signed plan's verification plan step 3.

## Fix-up cycle 3 (2026-09-25, kimi-1)

status: complete
completed: 2026-09-25
head-commit: 717f3de
skill-commit: a624318
record-commit: ee8849c

Applies the signed cycle-3 fix plan (`review/consensus.md` — claude-1 ✅ ACCEPT, kimi-1 ✅
ACCEPT, zcode-1 ✅ ACCEPT; no BLOCK). Two one-clause documentary repairs: no behaviour change, no
protocol text, no skill delta, no FINAL edit, no new test. Phase-8 packet attested at this
publication (2026-09-25, parley 1.49.1, this worktree; the source delta against reviewed
`1bad263` is exactly the AF-17 comment line): `parley protocol packet --dir . --phase 8 --track
deliberation --idea meta-protocol-change-designated-implementer --flag auto_implement --flag
protocol_change --audience participant --json`, exit 0 → `context_mode: "full"`,
`fallback_reason` ABSENT (top-level keys enumerated: `body_path`, `context_mode`, `index`,
`packet_sha256`, `request`, `shadow`, `source`, `source_sha256`), `source_sha256` =
`packet_sha256` = `b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388` (1,400 lines
/ 115,166 bytes — this cycle moved no protocol text); phase-8 shadow WITH both flags 86,716
bytes / `fe0e4c04…ebcf` / 40 included / 29 omitted, not used. Phase 8 and §15 were read from the
body.

### Fixes applied

- **AF-17** (origin: claude-1/review/round-03 NIT-1) — in the cycle-3 SOURCE commit `717f3de`:
  `internal/app/driver_designation_test.go:448` comment now reads "Present-empty is an incomplete
  designation, on any run that reaches a role action." — the corrected AF-3/R-1 scope wording
  AF-13 applied at `driver_impl.go:188`, reused, not expanded. Comment only; no assertion,
  fixture or behaviour change. Post-fix sweep: `grep -rn "any run" --include="*.go" .` leaves
  only qualified or unrelated hits (`driver_impl.go:188` qualified; `internal/agents/naming.go:50`
  and `internal/budget/run_identity_inventory_test.go:351` unrelated English).
- **AF-16** (origin: claude-1/review/round-03 MINOR-1) — in THIS record commit, the two named
  locations only: (1) the cycle-1 Phase-8 attestation bullet under `## Validation evidence` now
  labels every shadow figure with its flag state — 86,716 B is the phase-8 shadow WITH `--flag
  auto_implement --flag protocol_change`; 86,701 B (`a0845614…f08efd`) is the same packet over
  the same source WITHOUT them and reproduces on demand; the "did not reproduce / superseded"
  clause is struck; the rule is completed to "Read the phase and the flags with the figure";
  (2) `### Deviations from agreed fixes` restated on the same terms. The AF-14 summary inside
  `## Fix-up cycle 2` → "Fixes applied" is deliberately NOT touched, per the signed plan — it
  describes what AF-14 delivered and stays literally true.
- **Standing per-publication bump (AF-12's restored obligation, not a new finding):** frontmatter
  at this publication — `status: fix-up-cycle-2` → `fix-up-cycle-3`; `head-commit: 1bad263` →
  `717f3de` (the cycle-3 source commit — PRIOR and knowable under the source-then-record order);
  `skill-commit: a624318` (unchanged; no skill delta). The cycle-2 section's `record-commit` is
  filled in: `b850576` (PRIOR and knowable at this publication).

### Deviations from agreed fixes

None. Both repairs are exactly the signed text; no item was infeasible; no protocol text, skill
file, FINAL, consensus, signature, or peer artifact was edited; no new owner question arises.

### Evidence this cycle (proportional, as ratified — no broad suite for a comment/prose-only delta)

- `go build ./...` exit 0; `go vet ./...` exit 0; `gofmt -l
  internal/app/driver_designation_test.go` clean (no output, exit 0).
- `go test ./internal/app/ -count=1 -run
  'TestMalformedTier2Gates|TestTier2UnavailabilityGateAndExits|TestPinDesignationConflictEscalates|TestUnsetPathIsByteIdentical'`
  → `ok  parley-deck-cli/internal/app  0.524s`, exit 0 — including `TestMalformedTier2Gates`, the
  test whose comment AF-17 touched.
- **AC-17 standing, honestly:** claude-1's independent full suite stays EXACT-COMMIT scoped to
  `d238238` and is NOT re-claimed; the targeted behaviour evidence (both reviewers' focused runs,
  M1=9 / M2=5, the probe sets) stays scoped to `1bad263`; this cycle's delta is one comment line
  plus record prose, so no suite result is claimed for it beyond the focused run above. The LE-7
  fresh-invocation goal-done check (AC-21's second clause) remains separate and still owed. No
  completion, goal-done, or release is claimed here.

### Review-pointer

Round-04 re-review reassesses exactly: AF-16's two locations (flag-state labels; the struck
supersession clause), AF-17's comment plus the tree-wide `any run` sweep, and the cycle-3
metadata bump (`status: fix-up-cycle-3`, `head-commit: 717f3de`, cycle-2 `record-commit:
b850576`) — per the signed plan's verification plan step 3.

## Outcomes & Retrospective

**Implementation complete — 2026-09-25, kimi-1 (sole implementer).** `status: complete` is set at
this publication — close-sequence step 3, after both predecessor gates landed:

- **Closing review consensus (cycle 4)** — `review/consensus.md`, `outstanding_agreed_fixes: 0`,
  `blocked: false`, signoffs claude-1 ✅ / kimi-1 ✅ / zcode-1 ✅, no reservations, preserved at
  `804522c`. Its one open item (claude-1 round-04 NIT-1, the blanket lead-in at
  `IMPLEMENTATION.md:325-326`) is a RECORDED DISMISSAL — an accepted residual standing UNREPAIRED
  with all three concurring; no acceptance criterion depends on it, and this closing edit does not
  touch that text. All seventeen agreed fixes (AF-1…AF-17) are applied as signed.
- **LE-7 goal-done check — PASS** (`goal-check-zcode-1.md`): a FRESH invocation of zcode-1, an
  existing non-implementer quorum member — never the organizer, never the implementer, never a
  fourth identity. Verdict `PASS`, no reservation: all twenty-one FINAL observable acceptance
  criteria verified against the current tree (CLI HEAD `804522c`, code tree `717f3de`, skill
  `a624318`, record `ee8849c`), each with PRIMARY current-tree evidence or exactly-scoped PRIOR
  evidence. Per LE-7/LE-11 the check could only withhold a close, never establish one — this
  status flip is my own act as implementer.

**Metadata at this publication:** `status: complete`, `completed: 2026-09-25`; `head-commit:
717f3de` and `skill-commit: a624318` retained (PRIOR and accurate — no code moved since); the
cycle-3 section's `record-commit` is filled in as `ee8849c` (PRIOR and knowable now). This record
names no commit that does not already exist — no self-hash.

**Achievements:** the owner-designated implementer mechanism shipped exactly as the frozen FINAL
specifies — four-state per-idea designation, one resolution chain (pin → per-idea → global
default → fallback), fail-closed validity gates, the tier-2 three-exit unavailability gate,
dormant tier-3 with honest fall-through, reviewer-diversity warning, and the global default
SHIPPED UNSET. Three protocol copies in fidelity (AC-1); findings converged 17 → 8 → 2 →
1-dismissed-NIT across four cycles, zero MAJOR/behaviour findings in the last two.

**Gaps and limitations — named, not waived (carried from the signed consensuses):** (1) no broad
`go test ./...` exists at `1bad263` / `717f3de` / HEAD — AC-17's whole-suite clause rests on the
exactly-scoped `d238238` independent full suite plus the goal check's live re-execution at HEAD;
(2) skill npm-level gates unexercised — owed to the organizer's release preflight; (3) no
end-to-end `parley run` against a real designated deck; no live §9.0 ping behind
`designeeAvailable`; TUI/pipeline-block surfaces remain FINAL deferrals F3/F10; Windows is DF-4;
(4) the NIT-1 residual stands unrepaired by recorded dismissal; (5) the pre-edit 86,336 B figure
is not absolutely re-derivable — its flag label rests on the bullet's own quoted command.

**Lessons for §13:** record the full invocation (flags included) with every measurement — one
unlabelled shadow figure cost two review cycles; tree-wide claims need tree-wide sweeps — the
third `any run` carrier sat in the delta's own test file; the per-publication frontmatter bump
fires at every cycle's record commit, not at idea close; carrying corrections in the record and
commit message — never by editing a signed body — was re-verified independently every cycle.

**Release work — still pending, none of it mine.** The owner/organizer-only sequence is unchanged
and unstarted here: `origin/main` integration (R56) in the ratified order (`release-1.49.1` done
→ this idea → `windows-portability`), release preflight incl. the skill npm gates, expected CLI
1.50.0 / skill-core 2.14.0 re-verified at staging, per-channel participant audit after each
channel exists (R51), `parley protocol publish` owner-attended (R54), and the owner's
post-release machine-default selection (`codex-1`) as post-release configuration only. No
participant merges, publishes, versions, installs, releases, or mutates any global default; the
product default ships UNSET. DF-1…DF-4 stand as recorded — DF-2's carrier
`meta-protocol-change-designation-launch-surfacing` stays NAMED and INACTIVE.
