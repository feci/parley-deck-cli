---
idea: meta-protocol-change-lean-organizer
status: implemented
implementer: zcode-1
started: 2026-09-23
completed: 2026-09-23
branch: /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer#lean-organizer
head-commit: 2b5e160
design-pr: n/a
implementation-pr: n/a
---

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

- **Phase-8 facilitator body = 71,168 B, above the 70,000 B guardrail.** The guardrail's
  ratified scope is the phase-1 / deliberation / github-pr body, which measures **58,190 B**
  (hard-asserted in tests). Phase 8 pins Phase 5–8 subsections + strict gate + stopping
  judgment simultaneously, which pushes that one phase over; the map and the ceiling are
  left untouched and the bytes are shown here for the quorum (test
  `TestLiveDeckFacilitatorAcrossPhases` logs every phase: 52,854 / 58,190 / 58,291 / 59,551 /
  62,469 / 59,791 / 61,005 / 68,880 / 71,168 B for phases 0–8). The never-cut floor is
  untouched everywhere.

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
- (2026-09-23 · zcode-1) Open item 3 (--optimize baseline): measured in the same test run —
  64,734 B at phase 1 / deliberation / github-pr (logged by
  `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`).
- (2026-09-23 · zcode-1) Open item 5 (signoff/stance parsing reuse): RESOLVED by reuse —
  PhaseDigest calls `consensus.Status` and the newly exported
  `consensus.ExpectedRoundParticipants` / `consensus.ResolveImplementer`; nothing forked.
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
  github-pr): facilitator body **58,190 B**; named-omission-set bytes 15,759 B; §2 3,977 B;
  guardrail 70,000 B; `--optimize` baseline **64,734 B**. Brief: ≤ 8,192 B (1,863 B live),
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
