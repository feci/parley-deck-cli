---
idea: meta-protocol-change-lean-organizer
status: fix-up-cycle-1
implementer: zcode-1
started: 2026-09-23
completed: 2026-09-23
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer#lean-organizer
head-commit: 64a622c
fix-up-cycle: 1
design-pr: n/a
implementation-pr: n/a
---

## Fix-up cycle 1

status: complete
completed: 2026-09-24
head-commit: 64a622c (fix-up source; this record commit follows on the same branch)
skill-commit: a820dc7 (lean-organizer-skill branch, same cycle)
opened: 2026-09-24

Authorized by all three signoffs on `review/consensus.md` (claude-1 and kimi-1
ACCEPT-WITH-RESERVATIONS, zcode-1 ACCEPT) — authorizing F1–F21, not closure. VC-1/VC-2
resolved by the reviewers' independent reproductions and kimi-1's §15.1
SELF-CORRECTIONs (both sustain claude-1's evidence; both sign F1's primary branch —
remove the gate — and F2's stricter FINAL B.3 reading); VC-3: both concur with F8 for
this release, subject to R2 (DF-1..DF-4 opened as real slugs — done by the organizer;
paths recorded in this section). R1 accepted: F4 records the measurement method and
the absolute deck path beside every figure. R3 stands: this cycle does not close the
idea; fresh independent review round 2 and a zero-fix review consensus are required
after these fixes.

### Fix-up plan (F1–F21, per review/consensus.md)

- F1 (primary branch, both reviewers signed): remove the `86d028b` hard
  design-consensus section gate; keep the A.4 repair (`RequiredConsensusSections` as
  the single source for the drafting prompt and scaffold generator; parity test
  re-pointed to prompt ↔ scaffold ↔ constant); add both regression fixtures (the
  protocol's own Phase-3 template renders `ready`; deck-corpus malformed count must
  not exceed the 9/80 base count).
- F2: `parley wait` implements FINAL B.3's stricter semantics — blocking escalation
  only for notes with matching `idea:`, `blocking:` not `no`, `status:` not
  answered/resolved, arriving after wait start; `driver.error` only after wait start;
  pre-existing notes/errors become digest annotations; skill `wait` docs updated.
- F3: ready set derived from `protocol.ValidImplementationStatus` minus in-progress
  states (adds `ready-for-review`); timeout line names the actual blocking condition.
- F4 (amended per R1): re-measure at the fix-up HEAD in a clean tree; record method
  and absolute path beside every figure.
- F5: brief phase from the driver run cursor, `status:` fallback only without a run;
  full deck status vocabulary mapped.
- F6: `packet check` map-level rejection extended to phase-pinned never-cut blocks
  (§15.x); allowance narrowed to transport-conditional §11 subsections; hostile-map
  negative test.
- F7: run-record windows made real (`run.json` `updated_at` advanced at
  `commitCursor`); synthetic-window test replaced by a driver-built one; live smoke
  re-run.
- F8: record the blind-spot (i) measurement (phase-5/8 §15.5/§15.6 omission,
  phase-0 §15 absence) with the whether-it-matters analysis; no map change in-cycle.
- F9: `fell_back` derived from the validator/ownership path; column legend updated.
- F10: `NextAwaitReviewArtifact` when the implementation is present and ready but no
  review round exists; dead branch deleted.
- F11: brief names the `.parley-runtime` packet-body cache in its "never stored" line.
- F12: `renderPacket` receives the resolved audience.
- F13: `BuildPhaseHandoffRecord` takes the run dir; loader parses all nine fields;
  full round-trip identity test.
- F14: task-local binary rebuilt from the fix-up HEAD in a clean tree; sha256 +
  `vcs.revision`/`vcs.modified` recorded.
- F15: `consensus.ResolveImplementer` export dropped (no callers); doc comment and
  Decision Log corrected.
- F16: `containsBytes`/`indexOfBytes` replaced by `bytes.Contains`/`bytes.Index`.
- F17: `strconv` sentinel dropped.
- F18: same-`reasons`-set assertion for the merged `--optimize`/audience branch.
- F19: `streamLines` resets the accumulator when an oversized line is dropped.
- F20: §9.0 permissive sentence in all three COOPERATION.md copies; the two
  inconsistent records (protocol-changelog, core-publish escalation note) aligned to
  §9.0; combined 2.13.0 core restaged and independently rechecked; no publish.
- F21: criterion-A test shapes — driver-level fixture auto-drive test asserting role
  launches/escalation via the event log; absent-field run-plan byte-identity test.
- Deferred follow-ups opened by the organizer (R2 satisfied):
  DF-1 `ideas/meta-protocol-change-consensus-duty-gates/00-prompt.md`;
  DF-2 `ideas/meta-protocol-change-facilitator-integrity-phase-coverage/00-prompt.md`;
  DF-3 `ideas/facilitator-packet-per-phase-bounds/00-prompt.md`;
  DF-4 `ideas/release-binary-reproducibility/00-prompt.md`.

### Fix-up progress

- (2026-09-24, pre-edit) IMPLEMENTATION.md opened for fix-up cycle 1 before any
  source edit; inputs re-read in full (live phase-8 packet, attestation verified
  against the file hash; frozen FINAL; both round-01 reviews; review/consensus.md
  including every signoff and R1–R3).
- (2026-09-24) F1, F15 applied (consensus gate removed; parity re-pointed to
  prompt ↔ scaffold ↔ constant; ResolveImplementer export dropped; Decision Log
  corrected). F2+F17 (wait.go stricter B.3 semantics + strconv sentinel). F3/F9/F10
  (phasedigest). F5+F11 (organizer brief). F6 (packet check phase-pinned proof).
  F12+F18 (packet.go). F13 (phasehandoff + driver caller). F16+F19 (usage_ingest).
  F7 (runmanifest.TouchUpdatedAt at commitCursor). F21 (driver-level auto-drive
  event-log test + run-plan byte-identity). F20 (§9.0 sentence ×3 copies, changelog
  + escalation note aligned, core restaged + rechecked). F4/F8/F14 measured from a
  clean clone at the fix-up HEAD (below). Both full suites green.
- (2026-09-24) Deferred-follow-up references in review/consensus.md updated from
  TBD to the organizer-opened slugs (DF-1..DF-4 paths recorded there; signoff
  blocks untouched — diff shows only the four slug lines).

### Fixes applied (fix-up cycle 1)

All at fix-up source commit `64a622c` (CLI) and `a820dc7` (skill) unless noted.

- **F1 (CRIT-1, primary branch — both reviewers signed for removal).** Removed the
  `86d028b` hard design-consensus section gate (`if !review { MissingConsensusSections … }`
  in `internal/consensus/consensus.go`) — no required-sections gate exists again, at
  parity with base `b37f7ef`. `RequiredConsensusSections` remains the single source
  for the drafting prompt and the scaffold generator; the parity test is re-pointed
  to prove prompt ↔ scaffold ↔ constant (`TestConsensusDraftPromptScaffoldParity`).
  Both regression fixtures added in `internal/consensus/consensus_gate_regression_test.go`:
  the protocol's own Phase-3 template (extracted live from the deck COOPERATION.md,
  placeholders substituted) triages **ready**; the deck corpus sweep reads
  **80 consensus.md checked, 9 malformed** — the exact pre-idea base count both
  reviewers measured (was 79/80 at the reviewed HEAD).
- **F2 (MAJ-1 + MAJ-2 ≡ K1-F4).** `internal/app/wait.go`: `blockingEscalation` →
  `qualifyingEscalationNotes` + `arrivedBlockingEscalation` — a note blocks only
  when its frontmatter `idea:` matches the awaited slug AND `blocking:` is not `no`
  AND `status:` is not answered/resolved AND its mtime is after wait start;
  `driverErrorEvent` → `driverErrorEventSince` over an event-count snapshot at wait
  start (`firstDriverErrorBefore` annotates history). Pre-existing qualifying notes
  and historical errors are printed as `note:` digest annotations (and as `notes` in
  the `--json` envelope `{notes?, digest}`). Skill `wait` docs state the semantics.
  Tests: new-arrival exit 4; cross-idea note (the six-week fixup-budget note) never
  blocks and is never annotated; `blocking: no` (this idea's core-publish note)
  never blocks; pre-existing note → annotation + timeout 3; historical
  error-then-recovery → boundary 0 with the annotation; error appended after wait
  start → exit 4. **Live check:** `parley wait --for round` on this deck now exits
  **0** (exit 4 forever at the reviewed HEAD); the driver-error inbox note for this
  idea is demoted to a `note:` line.
- **F3 (MAJ-3).** `implSection` derives ready-for-review from
  `protocol.ValidImplementationStatus` minus `{"", "unparsed", "in-progress"}` —
  `ready-for-review` (4 live deck files) now ready. `outstandingAgents` names the
  actual blocking condition (implementation status X is not a recognised ready
  state / IMPLEMENTATION.md not filed) instead of "none named — inspect the digest".
- **F4 (MAJ-4 ≡ K1-F2, amended per R1).** Every figure re-measured at the fix-up
  HEAD in a clean clone; method and absolute path recorded beside the numbers —
  see `## Fix-up validation evidence` below. Stale Phase-5 figures corrected in
  `## Deviations from FINAL.md`, `## Validation evidence`, and Decision Log item 3.
- **F5 (MAJ-5).** `briefPhase` is now the most-advanced of: 00-prompt status (full
  deck vocabulary: `implementation`/`implemented`/`in-progress`→5,
  `ready-for-review`/`review*`→6, `fix-up-cycle-*`/`complete`→8), the driver run
  cursor (`run.json` `phase` + `current_round`), review-round presence (→6), and the
  IMPLEMENTATION.md status. A live Phase 5–8 idea can no longer resolve to the
  §15-free phase-0 packet (`TestOrganizerBriefPhaseNeverZeroForLivePhaseFivePlus`).
- **F6 (MAJ-6).** `packet check` map-level rejection extended from `always` to all
  non-transport-conditional never-cut rules (phase-pinned §15.x, `### Phase N`,
  flag-pinned §7); the allowance is narrowed to transport-conditional §11
  subsections. Negative test: hostile omits naming `## 15.`, `### 15.1`,
  `### 15.7` FAIL check; the clean map's §11.A/B/C omissions stay legal
  (`TestPacketCheckFailsAudienceOmittingPhasePinnedNeverCut`,
  `TestPacketCheckAllowsTransportConditionalOmit`).
- **F7 (MAJ-7).** `runmanifest.TouchUpdatedAt` advances `run.json` `updated_at` at
  `commitCursor` (evented, non-fatal). The synthetic 2-hour-window test is replaced
  by records built by the real machinery (`runmanifest.New/Write` at creation with
  created==updated, `TouchUpdatedAt` at the transition) in both
  `internal/driver/phasehandoff_test.go` and `internal/app/usage_ingest_test.go`.
  Live smoke re-run (below): attribution still `ambiguous` against the deck's
  historical zero-width manifests — honest; windows open for future transitions.
- **F8 (MAJ-8 ≡ K1-F5).** Blind-spot (i) measurement recorded below (line-anchored,
  live map, fix-up HEAD): phase-5 and phase-8 facilitator bodies carry `### 15.7`
  but not `### 15.5`/`### 15.6`; phase-0 omits §15 entirely; phase-7 carries all
  three. No map change in-cycle (accepted for this release by both signoffs, subject
  to R2; the pin question is DF-2, opened).
- **F9 (MIN-1).** `fell_back` now derives from the validator/ownership path
  (frontmatter unreadable, or no usable `agent:` owner → filename-derived
  attribution); the discarded `extractPosition` call is gone; column legend updated.
  Live digest now shows `fell_back=false` on all valid rows (was true on every
  valid round-02 row).
- **F10 (MIN-2).** `nextAction` returns `NextAwaitReviewArtifact` when the
  implementation is present and ready and no review round exists; the dead twin
  branch is deleted.
- **F11 (MIN-3).** The brief's opening line now states "never stored in the deck"
  and names the `.parley-runtime/protocol-packets/` body cache the renderer writes
  outside the deck.
- **F12 (MIN-4).** `renderPacket` receives the RESOLVED audience; `--optimize` with
  an unrecognized audience no longer stamps `audience=banana` into a full-fallback
  body (`TestOptimizeUnknownAudienceDoesNotStampRejectedAudience`).
- **F13 (MIN-5 ≡ K1-F6a).** `BuildPhaseHandoffRecord` takes the run dir parameter
  (the `"parley-deck"` derivation is gone); `LoadPhaseHandoffRecord` parses all nine
  frontmatter fields; the round-trip test asserts nine-field identity and the
  builder test asserts RunID = run-dir base.
- **F14 (MIN-6).** Task-local binary rebuilt from a clean clone at the fix-up HEAD:
  `/tmp/parley-lean-organizer-fixup1/parley`, sha256
  `8bd06646e0cbf1286d62d5de33e9ad885ce3b1fbbdf1ff811ff4f3b11555341a`,
  `go version -m` → `vcs.revision=64a622ce…`, `vcs.modified=false`, go1.27.1
  darwin/arm64. (A fifth distinct sha256 over behavior-identical trees — DF-4
  evidence; NOT installed globally.)
- **F15 (MIN-7).** `consensus.ResolveImplementer` export dropped (zero callers; the
  doc's claimed `participants[0]` fallback never existed in the internal resolver);
  the Decision Log entry claiming PhaseDigest calls it is corrected above.
- **F16 (NIT-1).** `containsBytes`/`indexOfBytes` replaced by `bytes.Contains`/
  (deleted — `bytes.Index` was only used by the hand-rolled contains).
- **F17 (NIT-2 ≡ K1-F6b).** `var _ = strconv.Itoa` sentinel and the import dropped.
- **F18 (NIT-3).** `TestOptimizeAudienceBranchSharesReasonsSet` pins the merged
  `--optimize`/audience branch to identical reason sets across kernel phases.
- **F19 (K1-F6c).** `streamLines` drops an oversized line WHOLE — a `dropping` mode
  discards the tail chunks so no fragment is ever yielded as a "line"
  (`TestStreamLinesDropsOversizedLineWhole`, 17 MiB fixture).
- **F20 (K1-F1).** The FINAL A.5 §9.0 permissive sentence added to all three
  COOPERATION.md copies (byte-identical, 265 B): "Declaring `facilitator:` makes the
  pure organizer this idea's default: …". `meta/protocol-changelog.md` corrected to
  list §9.0 among the A.5 lines and to label the audience-view line "§9 checklist
  item 1 (D.5)"; the core-publish escalation note corrected to "§9.0 pure-organizer
  default sentence, §9 checklist item 1 audience view + brief re-orientation".
  Combined 2.13.0 core RESTAGED from the amended template-form copy and
  independently rechecked (below). No publish.
- **F21 (K1-F3).** `TestFixtureAutoDriveNeverLaunchesFacilitatorForCodeRoles`: a
  real `driver.Advance` auto-drive over the production ops and REAL fixture agent
  launches — the event log's `agent.started` trail names the implementer and
  reviewer, never the declared facilitator; plus the facilitator-only escalation
  with ZERO launches. `TestPlanByteIdenticalWithAbsentFacilitatorField`: the
  absent-field run plan is byte-identical with/without the optional field.

### Deviations from agreed fixes

None in substance. Two recorded observations inside the signed scope:

- F10's NextAwaitReviewArtifact fires when no review round exists, exactly as
  signed. Residual: after a fix-up publishes (IMPLEMENTATION status
  `fix-up-cycle-1`) while the PREVIOUS review round still exists on disk, the
  enumeration reads `await implementation` until the first round-02 review file
  lands (then `await review artifact`). Matches the signed fix text; noted for the
  review round.
- F7's live smoke still reports `attribution=ambiguous` — correct and expected: the
  deck's existing run manifests keep their zero-width historical windows; only
  transitions after this fix open real ones. The unit/driver tests prove the window
  mechanism; the ledger records the honest live state.

### Fix-up validation evidence

All commands run 2026-09-24. "Clean clone" = `git clone --no-hardlinks` of
parley-deck-cli checked out at the fix-up source commit `64a622c`, working tree
clean (`git status --porcelain` empty) at `/tmp/fixup1-clean`.

**Both full suites (at the fix-up working trees):**

- CLI: `go build ./... && go test ./... -count=1 -timeout 2400s` — **exit 0, 31/31
  packages ok** (log `/tmp/fixup1-cli-full.log`), including every new fix-up test.
- Skill: `npm test` — **exit 0, 399 pass / 0 fail** (log `/tmp/fixup1-skill-full.log`);
  `node --test test/lean-organizer.test.js` 8/8 with the §9.0 hunk assertion added.
- Targeted fix-up packages re-run in the clean clone: consensus / protocol /
  protocolpacket / driver / app / runplan / runmanifest — ok (the full suite was run
  at the fix-up working trees; the clean clone re-ran the measurement tests and the
  touched packages).

**F4 measurements (method + absolute path beside every figure, per R1).**

*Method A — test-path vector (clean clone, relative authority path
`../../parley-deck/COOPERATION.md` (32 chars); command
`go test ./internal/app/ -run 'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestLiveDeckFacilitatorAcrossPhases' -count=1 -v`
at `/tmp/fixup1-clean`):* phase 0–8 facilitator bodies
**53,870 / 59,206 / 59,307 / 60,567 / 63,485 / 61,129 / 62,211 / 70,086 / 72,696 B**;
phase-1 named-omission-set 15,759 B; §2 3,977 B; `--optimize` baseline **65,750 B**.
The vector equals the reviewers' `3c97f44` vector + **265 B** at every phase — the
F20 §9.0 sentence (§9 is retained at every phase). Phase 1 = **59,206 B ≤ 70,000 B**
(guardrail met, ~10.8 KB margin). Phase 7 = 70,086 B — **86 B over** the 70,000 B
figure at this path (was 179 B under at `3c97f44`); phase 8 = 72,696 B. Per-phase
policy → DF-3.

*Method B — CLI render at the live deck (absolute authority path
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/parley-deck/COOPERATION.md`
(101 chars); `parley protocol packet --audience facilitator --phase N --track
deliberation --transport github-pr` with the fix-up binary; written body file bytes):*
phase 0 **53,939 B**, phase 1 **59,275 B**, phase 5 **61,198 B**, phase 7
**70,155 B**, phase 8 **72,765 B** — exactly Method A + **69 B** at every phase,
reproducing the R1 path-length mechanism 1:1 (101 − 32 = 69; the body embeds
`Source.Path`).

**F8 — blind-spot (i) measurement (line-anchored heading checks over the Method B
bodies, live map, fix-up HEAD):**

| phase | §15 | §15.5 | §15.6 | §15.7 |
|---|---|---|---|---|
| 0 | absent (phase 0 is outside the kernel set — no §15 at all) | absent | absent | absent |
| 1 | present | present | present | present |
| 5 | present | **absent** | **absent** | present |
| 7 | present | present | present | present |
| 8 | present | **absent** | **absent** | present |

Whether it matters (the FINAL-required analysis): §15.5/§15.6 bind the DRAFTER at
the phases where they pin (3, 6, 7); at phases 5 and 8 the facilitator's role is
verdict-reading and implementation/review adjudication, and §15.3/§15.4/§15.7 (the
close conditions, provenance, per-track table) are all present. The residual cost
both reviewers named stands: at phase 8 — fix-up adjudication — §15.7's table
asserts duties whose text (§15.5 role concentration, §15.6 alternatives/correlated
agreement) is cut. Accepted for this release per both signoffs (VC-3), subject to
R2; pinning them at 5/8 (+~2.4 KB at phase 8, already over) is DF-2's §7 question.
The phase-0 §15 absence is the map's design (kernel = {1,2,3,5,6,7,8}) and F5
ensures live Phase 5–8 briefs never resolve there.

**F7 live smoke (fix-up binary, live deck):** `parley usage ingest --agent codex-1
--source codex-rollout --path ~/.codex/sessions/2026/09/24/rollout-2026-09-24T00-18-19-….jsonl
--idea meta-protocol-change-lean-organizer --phase 8` → appended, 18 events,
total_tokens=1,513,379, **attribution=ambiguous** ("no run-record window contains the
accounting event"), stdout ≤ 1 KB; re-ingest → `idempotent no-op`, ledger unchanged
(1 row + method header). Honest: historical manifests keep zero-width windows;
`commitCursor` transitions from this fix onward advance `updated_at` (proven by
`TestCommitCursorAdvancesRunManifestUpdatedAt` and the driver-built attribution test).

**F2 live check (fix-up binary, live deck):** `parley wait --idea
meta-protocol-change-lean-organizer --for round --timeout 6s` → digest (round-02
3/3, review round-01 2/2, consensus reserved, implementation fix-up-cycle-1,
fell_back=false on all valid rows) + `note: pre-existing unanswered to-user
escalation … driver-error.md` + `wait: boundary reached (round complete)` →
**exit 0**. At the reviewed HEAD this exact invocation exited 4 forever on a
six-week-old cross-idea note.

**F14 binary provenance:** built in the clean clone: `go build -o
/tmp/parley-lean-organizer-fixup1/parley ./cmd/parley` (go1.27.1 darwin/arm64) →
sha256 `8bd06646e0cbf1286d62d5de33e9ad885ce3b1fbbdf1ff811ff4f3b11555341a`;
`go version -m` → `vcs.revision=64a622ce48660725a59303befbd351962864c8e8`,
`vcs.time=2026-09-24T01:48:45Z`, `vcs.modified=false`. NOT installed globally.

**F20 staged-core restage + independent recheck:** the combined core
`~/.parley/staging/COOPERATION-2.13.0.md` restaged from the amended template-form
third copy (skill `skills/parley-deck/references/COOPERATION.md`; the pre-edit
staged core sha256 `a8d3457a…` was verified byte-identical to the pre-edit copy
before restaging). New staged sha256
**`fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`**, **109,772 B**.
Recheck (vs published core 2.10.0 at `~/.parley/protocol/core/2.10.0/COOPERATION.md`):
placeholder header intact (`<workspace-name>`, `<transport-choice>`, `<YYYY-MM-DD>`;
no deck `Protocol synced:` values); §2 stubbed (grep claude-1|kimi-1|zcode-1|codex-1
→ **0**); diff = 11 hunks, each attributed by normalized signature to exactly one of
the three change sets — 3× 2.11.0 (§15.6/§15.7), 2× 1.48.0 (LE-7/LE-11 bullet,
goal-check withhold-only paragraph), 5× this idea pre-fixup + 1× this idea F20 §9.0
sentence; **0 unexplained**; §15 region (8,041 B), §9.0 region, §4 Phase-5, §4
Phase-6 and §11.B regions byte-equal to the deck view. No publish performed; the
owner's attended command in the escalation note is unchanged.

**F1 corpus fixture (clean clone):** `TestLiveDeckConsensusCorpusMalformedNotAboveBase`
→ `corpus: 80 consensus.md checked, 9 malformed (base 9)`; the protocol's own
Phase-3 template fixture triages `ready`.

## Protocol context attestation (Phase-5 implementation)

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase5-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

## Summary of work

Implementing the frozen FINAL (commit `120a9bf`) across both owner-supplied worktrees:
scope A (pure-organizer default for declared facilitator runs), B (`parley wait` + PhaseDigest),
C (audience-scoped protocol packet, slim SKILL.md core, computed organizer brief), D (driver-written
per-phase handoff records, `parley usage ingest` client-accounting ledger, kimi telemetry), and the
cross-cutting protocol-text/staging work — plus my own impl-claim signoff-status correction
(claude-1 signed 🟡 ACCEPT-WITH-RESERVATIONS, not an unconditional ACCEPT; corrected in
`inbox/zcode-1-to-all_meta-protocol-change-lean-organizer_impl-claim.md` this session).

`implementer:` (zcode-1) ≠ facilitator (codex-1) — the A-criterion witness for this run.

## Implementation plan / checklist

- [x] A.1 Parse optional `facilitator:` + `facilitator_participates:` fields in idea frontmatter (CLI).
- [x] A.2 `parley preflight` fail-closed (non-zero, both fields named) when `facilitator:` ∈
      `participants:` without `facilitator_participates: true`; exit 0 with it; absent field → untouched.
- [x] A.3 Driver role-ineligibility predicate: declared facilitator never selected as
      drafter/implementer/reviewer/goal-done checker; escalate (never silent fallback) when no
      participant implementer can be launched. Fixture auto-drive test with event-log assertion.
- [x] A.4 Prompt repair: `RequiredConsensusSections` shared constant; `buildConsensusDraftPrompt`
      emits canonical Phase-3 sections + §15.3/§15.5/§15.6 duties; parity test prompt↔gate;
      heading-consumer search documented below.
- [x] A.5 Absent-field deck → byte-identical run plan (regression test).
- [x] B.1 `PhaseDigest` extending `BuildRoundDigest` (design rounds, review rounds, consensus
      signoff state, implementation status) with per-agent columns
      agent/path/filed/bytes/owner/valid(failing check named)/stance_flags/unparsed/fell_back;
      validity from shipped validators.
- [x] B.2 `parley wait --idea <slug> --for round|consensus|review|implementation|any
      [--timeout D] [--json]`; event-log blocking with ≥10 s portable polling fallback; exit
      0/3/4/1; missing≠invalid; early exit 4 on unanswered `to-user` escalation / `driver.error`.
- [x] B.3 Timeout default 25 m, configuration-first, `default < 30m && configurable`; hard ceiling
      vs active track's §4.0 agent timeout.
- [x] B.4 Digest guardrail tests: golden byte-identity, no-model-written-field structural test,
      fixed-enumeration next action, digest-never-rewrites regression.
- [x] C.1 `parley protocol packet --audience participant|facilitator`; `audiences:` key in
      `meta/packet-applicability.yaml`; never-cut floor unbreachable (negative test); additive
      `audience` attestation field; unknown audience → full fallback with reason;
      `facilitator_participates: true` → full.
- [x] C.2 Facilitator retention set verbatim (Quickstart, §4, §5, §9, active §11, §2, §15) +
      complete omission index; named-omission-set absence is the GATING check (R-2); floor
      (≈42 KB with §2), ≤70,000 B guardrail, and `--optimize` baseline measured in the same run.
- [x] C.3 Slim SKILL.md core ≤20,000 B (relocation-only; frontmatter + Core Rule verbatim; six
      driver commands named; every reference linked); relocation proof script/test.
- [x] C.4 `parley organizer brief --idea <slug>`: ≤8,192 B, byte-identical ×2, writes no file
      (read-only-deck test); content = attestation + facilitator body path, `parley status --json`
      state, PhaseDigest, driver next action, phase pointer.
- [x] D.1 Driver per-phase handoff records under `runs/<run-id>/` on `WriteHandoffPacket` +
      shared PhaseDigest computation; schema doc states recomputation authoritative; unit test.
- [x] D.2 `parley usage ingest --agent --source codex-rollout|claude-jsonl --path --idea --phase`:
      no count-accepting flag (flag-set test); streaming bounded memory over the 228 MB rollout;
      ≤1 KB stdout; one ledger row (idea/phase/agent/source path/parser id/ingest time/six
      `total_token_usage` fields verbatim); idempotent; round-trip re-parse equality.
- [x] D.3 Attribution by explicit args + timestamp windows; residue `ambiguous`; method stated in
      ledger header; never slug scanning.
- [x] D.4 kimi telemetry case in `internal/telemetry/usage.go` from captured fixtures; visible
      `coverage: none` when unobtainable; structured argv only where adapter-supported (test).
- [x] P.1 Protocol text (all three COOPERATION.md copies): permissive A lines (§4 Phase 5,
      §4 Phase 6, §9.0, Quickstart facilitator row), §9 brief re-orientation line, §11 one-blocking-wait
      advisory line; CLI drift test green; skill-copy parity check (repo tooling).
- [x] P.2 Staged core in `~/.parley/staging/` from core 2.10.0 TEMPLATE + 2.11.0 hunks + 1.48.0
      hunks + this idea's hunks (placeholder header, stub §2; suggested 2.13.0); verification
      (three-change-set diff check); inbox note with exact owner-only attended publish command.
- [x] P.3 Windows CI leg covering `wait`/`usage` (CLI repo has no workflows today; add one).
- [x] P.4 Task-local CLI binary built (`go build`), path + sha256 recorded; NEVER installed globally.
- [x] T.1 CLI: `go build ./... && go test ./...` all green incl. new tests.
- [x] T.2 Skill: `npm test` all green incl. new tests; `npm run manifest:addons` after payload edits.
- [x] T.3 Adversarial negative cases per FINAL (never-cut breach, retained omission-set block,
      count-flag rejection, timeout ceiling rejection, unknown audience, read-only deck, etc.).
- [x] C-claim Correct own impl-claim signoff sentence (done — see Summary).
- Checks to run: CLI `go build ./... && go test ./...`; skill `npm test`; staged-core three-way diff.
- Review or risk notes: byte caps are binding — if 20,000 B / 70,000 B / 8,192 B genuinely cannot
  be met, bytes come back to the quorum (logged as deviation/blocker, never silently relaxed).

## Deviations from FINAL.md

None in scope or mechanisms. One measured finding is surfaced for the quorum rather than
resolved unilaterally (the FINAL's own open-item-2 discipline):

- **Phase-8 facilitator body above the 70,000 B guardrail** (figures corrected
  2026-09-24 by fix-up F4 — the original Phase-5 record transcribed pre-P.1-hunk
  values that reproduce at no committed state; cause and corrected vector per both
  reviewers' independent reproductions). The guardrail's ratified scope is the
  phase-1 / deliberation / github-pr body, which measures **59,206 B** at the
  fix-up HEAD (58,941 B at the reviewed HEAD `3c97f44`; +265 B = the F20 §9.0
  sentence). Phase 8 pins Phase 5–8 subsections + strict gate + stopping judgment
  simultaneously, which pushes that one phase over. At the fix-up HEAD
  `TestLiveDeckFacilitatorAcrossPhases` logs: **53,870 / 59,206 / 59,307 / 60,567 /
  63,485 / 61,129 / 62,211 / 70,086 / 72,696 B** for phases 0–8 — phase 7 now sits
  **86 B over** the 70,000 B figure at the test's relative path (it was 179 B under
  at `3c97f44`; the F20 sentence is 265 B), one more reason the per-phase policy
  question belongs to DF-3. The map and the ceiling are left untouched; the
  never-cut floor is untouched everywhere. What the tests assert (correcting the
  Phase-5 record's "hard-asserted" clause): the phase-1 test
  (`TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`) asserts the
  named-omission-set absence, the retention set and the ≤ 70,000 B guardrail;
  `TestLiveDeckFacilitatorAcrossPhases` asserts the never-cut floor at every phase
  and only LOGS the byte vector (per-phase bounds policy → DF-3). Measurement
  method and absolute paths: see `## Fix-up validation evidence`.

## Notes for reviewers

- FINAL is frozen at commit `120a9bf`; locators cited there were re-verified at implementation
  HEAD before edits (any drift is recorded in Surprises).
- Refutation targets are the FINAL acceptance table rows; every test below maps to a row.
- The three COOPERATION.md copies must carry identical hunks — diff them directly.

## Progress

- (2026-09-23 23:34Z) ALL APPROVED SCOPE IMPLEMENTED. C: audience packet (`--audience`,
  `audiences:` map key, never-cut floor enforcement + negative checks), organizer brief,
  live-deck guardrail measurements. D: phase handoff records, `parley usage ingest`
  (path-only, streaming, idempotent, attribution windows), kimi telemetry + structured
  argv + envelope unwrap. P: protocol hunks ×3 copies + changelog entry + Windows CI leg
  (CLI), staged core 2.13.0 (verified), owner core-publish escalation note, task-local
  binary. Both suites green (CLI 31/31 packages; skill 399 pass). (completed: everything;
  remaining: review phases 6–8)
- (2026-09-23 21:54Z) Phase-5 dispatch received; protocol packet + FINAL + consensus + signoffs +
  organizer notes + own claim read; claim signoff-status sentence corrected (Claude = 🟡
  ACCEPT-WITH-RESERVATIONS with R-1/R-2 concurred and carried). IMPLEMENTATION.md opened before
  any source edit. (completed: context intake + plan; remaining: all code items)
- (2026-09-23 22:27Z) Scope A implemented + tests green: facilitator frontmatter
  (internal/protocol/facilitator.go, workspace.go), preflight fail-closed gate (preflight.go),
  driver role-ineligibility with escalate-not-fallback (driver_impl.go, driver_consensus.go),
  RequiredConsensusSections prompt/gate/scaffold repair (consensussections.go, consensus.go,
  driver_consensus.go). Tests: internal/app/facilitator_test.go (6), facilitator_roles_test.go (5).
- (2026-09-23 22:27Z) Scope B implemented + tests green: PhaseDigest
  (internal/driver/phasedigest.go; reused validators + consensus.Status, exported
  consensus.ExpectedRoundParticipants/ResolveImplementer instead of forking), `parley wait`
  (internal/app/wait.go; exits 0/3/4/1, missing≠invalid, escalation/driver.error exits,
  ≥10 s portable polling, per-track ceiling, [defaults.timeouts] wait_ms seeded 25 m).
  Tests: internal/app/wait_test.go (11). (completed: A, B; remaining: C, D, protocol text, skill, staging)

## Decision Log

- (2026-09-23 · zcode-1) Open item 11 (preflight soft-warn when no `facilitator:` field): default
  **no warning** per consensus signoff record — implemented as silence.
- (2026-09-23 · zcode-1) Open item 1 (wait timeout config key): RESOLVED — a new seeded
  `wait_ms` key inside the shipped `[defaults.timeouts]` block (default 1,500,000 = 25 m;
  per-call `--timeout` overrides; per-track §4.0 ceiling rejects over-large values).
  Satisfies `default < 30m && configurable`.
- (2026-09-23 · zcode-1, figure corrected 2026-09-24 fix-up F4) Open item 3 (--optimize
  baseline): measured in the same test run — **65,750 B** at phase 1 / deliberation /
  github-pr at the fix-up HEAD (65,485 B at the reviewed HEAD `3c97f44`; the original
  64,734 B was pre-P.1-hunk and reproduced at no committed state), logged by
  `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`.
- (2026-09-23 · zcode-1, corrected 2026-09-24 fix-up F15) Open item 5 (signoff/stance
  parsing reuse): RESOLVED by reuse — PhaseDigest calls `consensus.Status` and the newly
  exported `consensus.ExpectedRoundParticipants`; nothing forked. The round-1 record here
  also named `consensus.ResolveImplementer` as a PhaseDigest callee — that was wrong
  (claude-1 MIN-7: `implSection` reads the IMPLEMENTATION.md frontmatter directly and the
  export had zero callers); fix-up F15 dropped the export and corrected this entry.
- (2026-09-23 · zcode-1) Open item 6 (heading-consumer search): DONE — repo-wide search for
  consumers of the old design-prompt headings (`## Trade-offs accepted` etc.) found no
  code consumer; the only existing heading consumers read `## Open items deferred to
  implementation` (Finalize reserved-check, unloggedReservations) and
  `## Agreed decisions` (new scaffold), both still emitted. Review-consensus headings
  (`## Agreed fixes` / `## Deferred follow-ups` / `## Dismissed findings`) are untouched.
- (2026-09-23 · zcode-1) Open item 7 (skill-copy parity shape): RESOLVED — hunk-presence
  assertions in `test/lean-organizer.test.js` (the manifest-coverage-pattern option):
  the six idea hunks are asserted verbatim in the bundled COOPERATION.md; the CLI's own
  drift test guards the two CLI copies; all three copies received byte-identical hunks.
- (2026-09-23 · zcode-1) Open item 8 (reference names): RESOLVED — claude-1's sketch adopted
  verbatim: `references/HEADLESS_LAUNCH.md`, `references/ARTIFACT_TEMPLATES.md`,
  `references/ROSTER_AND_PROTOCOL.md`.
- (2026-09-23 · zcode-1) Open item 9 (handoff record): RESOLVED — one record per transition at
  `runs/<run-id>/handoff-phase-<phase>.md`, written by `commitCursor` (the single phase-
  transition chokepoint), atomic write, schema doc embedded in each record stating the
  recomputed view is authoritative.
- (2026-09-23 · zcode-1) Open item 10 (machine-readable ledger sibling): RESOLVED —
  `parley-deck/ideas/<slug>/usage-ledger.jsonl`, one JSON row per ingest, attribution-method
  header line, idempotent append.
- (2026-09-23 · zcode-1) Kimi structured usage (open item 4): the live wire shape is
  `{"type":"usage.record","usage":{inputOther,inputCacheRead,inputCacheCreation,output}}`
  (captured this run into `source-context/kimi-usage-record.jsonl` and
  `internal/telemetry/testdata/kimi-stream.txt`). Probed live: kimi 0.42.0 emits NO usage
  to stdout in stream-json mode (only meta envelopes + assistant content), so live captures
  honestly downgrade to `coverage: none`; the parser accepts both the plain and the
  role-envelope-wrapped record for when the CLI starts emitting them.

## Surprises & Discoveries

- **kimi 0.42.0 stdout carries no usage** (probed live 2026-09-23 with
  `--output-format stream-json -p`): usage records exist only in the on-disk
  `~/.kimi-code/sessions/**/wire.jsonl`. Consequence: the structured argv change is
  adapter-supported and adopted, text consumers unwrap assistant content, and telemetry
  honestly reports `coverage: none` until kimi emits usage on stdout. This is the FINAL's
  accepted degradation, not a gap.
- **The kimi stream-json envelopes broke three text consumers** (preflight exact-PONG,
  consult answers, the stdout artifact fallback) — all three now unwrap via
  `runner.UnwrapKimiStreamJSON`; this repairs the observed "kimi stream-json preflight
  parser rejection" driver gap as necessary hardening of the D.4 argv change.
- **Skill repo has NO Windows CI leg today** (blind spot ii, checked: `test.yml` is
  ubuntu-only; `release-portable.yml` builds Windows binaries ON ubuntu). The ratified
  Windows leg for `wait`/`usage` is added to the CLI repo (`.github/workflows/tests.yml`,
  ubuntu + windows + macos matrix). The skill suite is pure Node/Python and its Windows
  binaries are built cross-platform; no skill-repo Windows leg was added (would need a
  Windows Python story) — recorded as the checked finding FINAL asked for.
- **This repo had no `.github/` at all** — the CI workflow file is entirely new, not an
  extension of an existing leg.
- The organizer's own run.json launch-args snapshot (stream-json override for kimi) is
  exactly the argv this change now ships as the default — the run was its own first test.

## Validation evidence

All commands run at the implementation HEAD in the CLI worktree
(`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer`, branch
`lean-organizer`) unless stated otherwise.

### CLI worktree

- `go build ./...` — exit 0.
- `go test ./... -count=1 -timeout 1500s` — **exit 0; 31/31 packages ok** (full log
  `/tmp/cli-full-suite.log`). Includes every new test below.
- **A** — `go test ./internal/app/ -run 'TestPreflightFacilitator|TestFacilitatorConflictGate|TestConsensusDraftPromptGateParity|TestConsensusPromptNamesFifteenDuties|TestDeclaredFacilitator|TestFacilitatorOnly|TestFacilitatorParticipates|TestAbsentFacilitator|TestFirstEligibleHeadless' -count=1` — ok.
  Proves: preflight fail-closed naming both fields; exit-0 with the exception flag;
  absent field untouched; parity over `protocol.RequiredConsensusSections` +
  `ConditionalConsensusSections`; driver never selects the declared facilitator for
  implementer/reviewer/drafter (goal-done checker = drafter); facilitator-only set
  escalates ("escalated, not fallen back") at Implement/OpenReviewRound/GoalCheck;
  v1.48.0 role selection preserved when the field is absent. Event-log assertion:
  `commitCursor` phase transitions emit `run.phase` + `driver.phase_handoff` events.
- **B** — `go test ./internal/app/ -run 'TestPhaseDigest|TestWait' -count=1` — ok.
  Proves: byte-identical digest over an unchanged tree; no model-written field
  (structural allowlist, no `position`); next-action ∈ fixed enumeration; exits 0
  (boundary, digest shows 2/2) / 3 (timeout names `kimi-1 (round artifact)`, partial
  digest) / 4 (present-but-invalid carries the validator reason verbatim
  `missing required section`; unanswered `to-user` note; `driver.error` event with the
  detail verbatim) / 1 (missing --idea, unknown --for, bad --timeout); missing ≠ invalid;
  ceiling rejects `--timeout 10m` on a `fast` deck naming the ceiling while accepting
  4m-equivalents; default 25 m < 30 m; poll interval ≥ 10 s; wait never modifies the
  idea tree (mtime snapshot regression).
- **C** — `go test ./internal/protocolpacket/ -count=1` — ok (incl. audience suite);
  `go test ./internal/app/ -run 'TestLiveDeckFacilitator|TestOrganizerBrief' -count=1` — ok.
  Proves: named-omission-set absence is the GATING assertion (line-anchored headings +
  body-text absence, immune to omission-index mentions); retention set verbatim; complete
  omission index; hostile map cannot cut §15 at a kernel phase; `packet check` fails an
  audience omitting an always-never-cut block and stays green on the shipped map across
  every phase × track; unknown audience → full + `unknown-audience:banana`;
  `facilitator_participates: true` → full + reason; attestation carries additive
  `audience`, never `role`. **R-2 same-run measurements** (phase 1 / deliberation /
  github-pr; figures corrected 2026-09-24 by fix-up F4 — the Phase-5 record transcribed
  pre-hunk values): facilitator body **59,206 B** at the fix-up HEAD (58,941 B at the
  reviewed HEAD); named-omission-set bytes 15,759 B; §2 3,977 B; guardrail 70,000 B;
  `--optimize` baseline **65,750 B** at the fix-up HEAD (65,485 B at the reviewed
  HEAD). Brief: ≤ 8,192 B (1,863 B live),
  byte-identical ×2, computes against a read-only deck writing nothing.
- **D** — `go test ./internal/app/ -run 'TestUsageIngest' -count=1` (incl. the 228 MB
  generated fixture streaming, ≤ 1 KB stdout, one row, six fields verbatim, idempotent
  no-op, round-trip re-parse equality, attributed-vs-ambiguous windows, flag-set source
  scan: no count-accepting flag); `go test ./internal/driver/ -run 'TestPhaseHandoff|TestCommitCursorWrites' -count=1`
  (record written per transition, schema-valid, schema doc states recomputation
  authoritative, round-trip); `go test ./internal/telemetry/ ./internal/runner/ -run 'TestKimi|TestUnwrapKimi' -count=1`
  (both live fixture shapes parse; coverage-none path; `StructuredArgs` true on the new
  argv; foreign shapes rejected by kimi and by claude; unwrap yields assistant content,
  non-JSON passthrough, never a guess). Live smoke: `parley usage ingest --agent codex-1
  --source codex-rollout --path <organizer rollout> --idea <slug> --phase 4` → 89 events,
  total_tokens=9,202,710, attribution=ambiguous (honest: no run window covers the final
  event).
- **Cross-cutting (CLI)** — `go test ./internal/protocol/ -run TestEmbeddedDefaultMatchesLiveDeck -count=1`
  — ok: the two CLI COOPERATION.md copies carry identical hunks (drift test green).
  `parley protocol packet check` — ok=True with the `audiences:` map (structure +
  audience never-cut proof across 27 phase/track builds).

### Skill worktree

(`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-skill`)

- `npm ci` then `npm test` — **exit 0; 399 pass / 0 fail** (node --test + python 54 +
  manifest check). `npm run manifest:addons` re-run after payload edits (SKILL.md +
  3 new references + COOPERATION.md hashes regenerated).
- `node --test test/lean-organizer.test.js` — 8/8 pass: core **17,038 B ≤ 20,000 B**;
  frontmatter description + Core Rule verbatim; six driver commands named
  (`run`/`continue`/`wait`/`status`/`consensus`/`preflight`); every reference file
  linked from core; relocation proof — all 29 original headings land (8 retained in
  core, 21 across exactly one reference each, none dropped); bundled COOPERATION.md
  carries the six idea hunks verbatim (third-copy parity, repo tooling).
- Pre-existing expectation updated for the new payload (installer doctor missing-file
  list now includes the three new references) — the list is derived from the copy plan,
  so this is the test tracking the payload, not a hand-edit.

### Cross-repo / staged core

- The three COOPERATION.md copies carry byte-identical hunks (CLI drift test + skill
  hunk assertions).
- Staged core: `~/.parley/staging/COOPERATION-2.13.0.md` (109,507 B), built from the
  TEMPLATE (placeholder header verified, §2 handle table stubbed) with exactly the
  three change sets (2.11.0, 1.48.0, this idea) — §15 region and every changed section
  byte-match the deck view; diff vs published core 2.10.0 contains exactly the three
  change sets.
- Task-local binary: `/tmp/parley-lean-organizer/parley`,
  sha256 `8d3e4c1c684023c03aaf26d2f47d52be27eaefa3693e5bc45db5aca0feca1dcb`
  (built for deterministic digest/brief use; NOT installed globally).

## Outcomes & Retrospective

(To be completed at Phase 8 close; preliminary:) shipped the full frozen FINAL A–D as
enforced tooling across both worktrees with zero dropped scope; the run itself exercised
its own product (this IMPLEMENTATION.md was written by the claimed non-facilitator
implementer; the organizer consumed only protocol artifacts). Notable learnings for
`parley retro`: fixture updates dominated the tail of the work (the consensus-prompt
repair intentionally changes what valid consensus.md bodies look like — 8 fixtures
updated); kimi's stdout-usage absence makes the honest coverage-none path the live
default; byte caps held everywhere without quorum escalation except the out-of-scope
phase-8 facilitator view (shown above).
