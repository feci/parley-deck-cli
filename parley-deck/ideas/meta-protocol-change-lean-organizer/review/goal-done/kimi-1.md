---
agent: kimi-1
idea: meta-protocol-change-lean-organizer
goal-check: true
date: 2026-09-24
reviewed-commit: 4df0855
---

# LE-7 goal-done check — meta-protocol-change-lean-organizer

**Verdict: PASS** — every observable acceptance criterion of the frozen FINAL (A–D rows
plus the cross-cutting row) was independently verified by me at the pinned closing HEAD
with my own commands, my own isolated trees, and my own binary. Residuals and test limits
are stated below; none is an unmet criterion, and none is a reservation on this verdict.
Per LE-7 this check can only withhold a close, never establish one: this PASS is
defense-in-depth evidence for the closing consensus, not the close itself.

**Commission:** `parley-deck/inbox/codex-1-to-kimi-1_meta-protocol-change-lean-organizer_goal-done.md`
(organizer → kimi-1, 2026-09-24, `blocking: no`). Checker eligibility: I am a
non-implementer participant (implementer: zcode-1) and not the declared facilitator
(codex-1); this is a fresh session launched for this check and the round-04 review.
Nothing in this check relies on implementer assertions or organizer testimony; every
criterion verdict below is PRIMARY (a check I executed, command and relevant output
quoted) unless explicitly tagged.

**Pinned closing HEAD (verified, not assumed):**

- CLI record tip `4df0855` (`git show --stat` → touches only `IMPLEMENTATION.md`), code
  = fix-up-3 source `c3baf09` (touches only `internal/driver/phasedigest.go`,
  `internal/app/wait_test.go`). FINAL frozen: `git diff 120a9bf..4df0855 --
  parley-deck/ideas/meta-protocol-change-lean-organizer/FINAL.md` → empty.
- Skill `b06a65a`: `git rev-parse lean-organizer` == `b06a65a` (unchanged since cycle 2).
- Frozen FINAL `120a9bf` is the criteria source; the A–D table rows are the lookup.

**Protocol context attestation (this check):**

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase6-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

Cross-check (PRIMARY): `shasum -a 256` over the packet body and over the live deck's
`parley-deck/COOPERATION.md` both return `8ce83cde…a9db7` — the attested packet IS the
live authority.

## Method and provenance

- Trees: my own isolated review worktrees
  `…/worktrees/lean-organizer-review-kimi-1` (new branch
  `review/meta-protocol-change-lean-organizer/kimi-1-20260924-r4` at `4df0855`) and
  `…/worktrees/lean-organizer-review-kimi-1-skill` (at `b06a65a`), plus my own clean
  clone `git clone --no-hardlinks` → `/tmp/kimi1-r4-clean` at `c3baf09`
  (`git status --porcelain` empty).
- Binary: built by me from the clean clone, `go build -o /tmp/kimi1-r4/parley
  ./cmd/parley` (go1.27.1 darwin/arm64) → sha256
  `0bdf69c3ef0ad0131307ec9a596c0fbc8bdc3ae6a09fa155dc2c95c103d3fc43`; `go version -m` →
  `vcs.revision=c3baf0977b82c2e9fc386b174f45de9e5a1dc1d9`, `vcs.modified=false`.
  Two-party reproduction: rebuilding from the implementer's clean clone
  `/tmp/fixup3-clean` yields sha256 `19831a24…06bdc`, byte-identical to the recorded
  task binary `/tmp/parley-lean-organizer-fixup3/parley`. Nothing installed globally.
- Suites run by me at the pinned commits: CLI `go build ./... && go test ./...
  -count=1 -timeout 2400s` → **exit 0, 31/31 packages ok, 0 FAIL**
  (`/tmp/kimi1-r4-cli-full-suite.log`); skill `npm test` → **exit 0, 399 pass / 0 fail**
  (`/tmp/kimi1-r4-skill-test.log`).
- Live-deck commands used my binary against the real deck read-only; `git status` of the
  live deck before/after shows no writes from this check. Fixtures disposable under
  `/tmp/kimi1-r4/`.

## Criterion-by-criterion evidence (frozen FINAL, "Observable acceptance criteria")

### A — pure organizer is the default for declared facilitator runs: PASS

- `parley preflight` exits non-zero naming both fields on conflict, exit 0 with the
  exception: my fresh fixture (`agents.toml` + `meta/version.json` + an idea declaring
  `facilitator: claude-1` ∈ `participants:`) + minimal `COOPERATION.md` → **exit 3**,
  the `[facilitator-declaration]` gate naming `facilitator:` and
  `facilitator_participates:` verbatim; the same fixture with
  `facilitator_participates: true` → the facilitator gate does not fire (0 hits in
  output; the residual exit 3 is my minimal fixture's unrelated unknown-freshness gate,
  and the exit-0 path is covered by the shipped
  `TestPreflightFacilitatorParticipatesFlagClearsConflict`, PASS in my run). The
  degenerate fail-closed shape (no COOPERATION.md) → **exit 1, stdout 0 B**, the read
  failure named. PRIMARY.
- Driver role-ineligibility + escalate-not-fallback + absent-field byte-identity +
  prompt parity: `TestDeclaredFacilitatorNeverSelectedForCodeRoles`,
  `TestFixtureAutoDriveNeverLaunchesFacilitatorForCodeRoles`,
  `TestFixtureAutoDriveFacilitatorOnlyEscalatesWithoutLaunching`,
  `TestAbsentFacilitatorFieldKeepsV148RoleSelection`,
  `TestFirstEligibleHeadlessAgentSkipsFacilitator`,
  `TestConsensusDraftPromptScaffoldParity`, `TestConsensusPromptNamesFifteenDuties` — all
  PASS in my own focused runs at `c3baf09`. PRIMARY.
- This run's `IMPLEMENTATION.md` records `implementer: zcode-1` ≠ facilitator. PRIMARY
  (read).

### B — `parley wait` + PhaseDigest: PASS

- Exit map, all reproduced by me this session with streams separated: **0** (live
  `--for round`: stdout decodes as the single `{notes, digest}` envelope, status line on
  stderr), **3** (live `--for review`: timeout, partial digest names `claude-1 (review
  artifact), kimi-1 (review artifact)`), **4** (my fixture: present-but-invalid review
  artifact → immediate exit 4 with the validator error verbatim on stderr, envelope-only
  stdout), **1** (`--for bogus`: usage error, stdout 0 B). PRIMARY.
- Digest columns per agent carry `path/filed/bytes/owner/valid(with failing
  check)/stance/unparsed/fell_back` — observed live in the exit-0/3 output. Determinism
  and no-model-written-field: `TestPhaseDigestByteIdenticalOverUnchangedTree` and
  `TestPhaseDigestHasNoModelWrittenField` PASS in my runs. PRIMARY.
- Timeout: seeded `[defaults.timeouts] wait_ms = 1500000` (25 min < 30 min) with the
  comment "NOT a verified provider cache fact" (no universal provider TTL asserted);
  per-call `--timeout` override; ceiling rejection covered by
  `TestWaitTimeoutCeilingRejectsAboveTrackTimeout` (green inside my full-suite run; the
  help text states the ceiling). PRIMARY.
- Digest never rewrites or replaces a round file: my live-deck runs left `git status`
  unchanged; the regression tests are green inside my suite run. PRIMARY.

### C — audience packet + slim SKILL.md + generated brief: PASS

- `parley protocol packet check` on the live map → `ok`, exit 0 (69 blocks). My own
  hostile-map fixture (facilitator audience omitting never-cut `## 15. Verification
  integrity`) → `packet check: FAILED`, **exit 1**, the violation named. PRIMARY.
- **Gating named-omission-set absence (R-2):** my own render (`--audience facilitator
  --phase 1 --track deliberation --transport github-pr`): all nine named omission
  headings (§1, §3, §8, §10, §12, §13, Appendix A, §11.A, §11.C) **absent**; all seven
  retention headings (Quickstart, §4, §5, §9, §2, §15, §11.B) **present**; never-cut
  §15.1/15.2/15.3/15.4/§15.7 present; `## Packet omission index` present with triggers;
  §5 extracted from body and source compare **byte-identical** (`cmp` silent);
  `source_sha256` over the full authority; additive `audience: facilitator` attestation
  field. PRIMARY.
- Byte recording in the same test run: `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`
  PASS in my clean clone, logging `body=59206 B; floor=52295 B; omission=30262 B;
  guardrail=70000 B; --optimize baseline=65750 B` — body against both the measured floor
  and the guardrail, plus the `--optimize` baseline, in one run. PRIMARY.
- Unknown audience → `context_mode: full` with `audience_fallback_reason:
  "unknown-audience:bogus"` (PRIMARY, my invocation);
  `TestFacilitatorParticipatesKeepsFullContext` PASS.
- SKILL.md core: **17,802 B ≤ 20,000 B** (`wc -c`, my review tree), names `parley run` /
  `continue` / `wait` / `status` / `consensus` / `preflight` (grep 1×/1×/2×/2×/1×/1×);
  relocation-parity and core tests green inside my `npm test` run. PRIMARY.
- Brief: **2,516 B ≤ 8,192 B**, `cmp`-identical across two runs, writes no file
  (live-deck `git status` unchanged across my invocations; the read-only contract tests
  green inside my suite run). PRIMARY.

### D — fresh session per phase + client-accounting ledger + telemetry: PASS

- Handoff records: `TestCommitCursorWritesPhaseHandoffRecord`,
  `TestPhaseHandoffRoundTrip`, `TestBuildPhaseHandoffRecordTakesRunDir` PASS in my
  focused run. PRIMARY.
- `parley usage ingest`: no count-accepting flag exists (PRIMARY, `--help` lists only
  `-agent -dir -idea -path -phase -source`; `TestUsageIngestHasNoCountAcceptingFlag`
  PASS). Live ingest by me of the real 2,438,827 B Codex rollout (the same file this
  idea's own ledger ingested) into a disposable deck: exit 0, **122 B stdout** (≤ 1 KB),
  one ledger row with idea/phase/agent/source-path/parser-id and the six
  `total_token_usage` fields verbatim (`total_tokens=1513379`, matching the idea ledger's
  independent row), `attribution=ambiguous` honestly recorded; re-ingest → `idempotent
  no-op (identical row exists)`, ledger stays one row. PRIMARY.
  `TestUsageIngestCodexRolloutSixFieldsVerbatim`, `TestUsageIngestIdempotent`,
  `TestUsageIngestAttributionWindows`, `TestUsageIngestStreamsLargeFixtureBoundedMemory`
  PASS in my runs.
- Kimi telemetry: `TestKimiUsageRecordParsesWireFixture`,
  `TestKimiCoverageNoneWhenNoUsageEmitted`, `TestKimiAdapterUsesStructuredArgvDetection`,
  `TestKimiUsageIgnoresForeignShapes` PASS in my runs — the honest `coverage: none` path
  included. PRIMARY.
- Protocol text: the permissive §9.0 audience/brief sentence is present in all three
  COOPERATION.md copies (grep 1× each) and the §11 wait-preference advisory line rides
  the same copies (§11 intro). PRIMARY.

### Cross-cutting: PASS

- All three protocol copies carry identical ratified hunks: the §15 region measures
  **8,056 B** in all three under the stated extraction convention (`sed -n '/^## 15\.
  Verification integrity/,$p' | wc -c`), the deck copy hashes to the attested packet
  (`8ce83cde…`), the CLI drift test is green inside my suite run, and the skill copy is
  byte-identical to the staged core (both sha256 `fc907e59…62c9f`). PRIMARY.
- Staged core exists at `~/.parley/staging/COOPERATION-2.13.0.md` (109,772 B, sha256
  `fc907e59…62c9f`, mtime 2026-09-24 03:35:27) matching the escalation note's facts
  (note read in full); `~/.parley/protocol/core/` holds only `2.10.0` — **no publish has
  occurred**; `parley protocol publish` remains the owner's attended-only action.
  PRIMARY.
- Both worktrees' test suites green at the pinned HEADs **run by me** (above). No
  release, merge, tag, global install, or publish was performed by anyone per the record
  (`VERSION` 1.48.0, CHANGELOGs untouched by the implementation/fix-up commits) — and
  none by me. PRIMARY.

## Original A–D goal — still holds

The owner-approved goal was to reduce organizer token usage **while preserving canonical
participant ownership and verification integrity**, shipped as mechanisms A–D. At the
closing HEAD the mechanisms exist, are enforced by tooling (not prose), and pass every
frozen observable criterion under independent hands. The integrity half is demonstrably
intact: this very check is a fresh non-implementer session verifying current-tree
evidence; the review cycle that preceded it caught and fixed a real silent-failure edge
(R3-MIN-1) instead of rubber-stamping. The savings half is honestly scoped: FINAL never
claimed a measured token reduction as a criterion — it ships the measurement ledger (D)
because the largest levers (skills catalogue, AGENTS.md read-in-full) are out of scope
and instructions alone did not move the rate in the study. The goal as approved holds;
the outcome is real but modest by design, and the ledger is the instrument that will
make the post-release effect visible rather than arguable.

## Residuals (disclosed, none blocking; not reservations on this verdict)

1. **GitHub-hosted runner wall-clock unmeasured.** G3's explicit `-timeout 45m` removes
   the known cliff; whether hosted runners need it is unverified. Local evidence:
   `internal/trajectory` measured 626.9 s in my run (over Go's 600 s default) — the flag
   is load-bearing locally.
2. **Live attribution windows remain test-proven only.** The deck's recent `runs/`
   records are `events.jsonl`-only or zero-width (re-verified this session: four
   events-only records, one zero-width manifest at `2026-09-23T20:25:01.377412Z`); no
   real driver transition has yet produced a non-zero-width window. The mechanism is
   unit/driver-tested; the live state is honestly recorded as `ambiguous`.
3. **Windows `wait`/`usage` portability is macOS-verified only** (the Windows CI leg
   ships in the workflow; it has not been observed green by me).
4. **Kimi structured usage** rests on fixtures captured this run; if the live kimi CLI
   stops emitting the structured envelope, the honest degradation is `coverage: none`
   (tested), never a fabricated number.
5. **Actual organizer token savings are unmeasured until D's ledger ingests real
   post-release runs.** This is a FINAL-declared measurement plan, not a gap in the
   shipped criteria.

## Test limits

- Full suites were run on macOS arm64 only; Windows/Linux legs exist in CI config but
  were not executed by me.
- The 228 MB streaming property is covered by the shipped bounded-memory test (green in
  my run); my live ingest used the real 2.4 MB rollout, not the 228 MB fixture.
- Equal-mtime verification used `git archive | tar -x` extraction (one of the operations
  the underlying finding named) plus the shipped `os.Chtimes`-pinned unit test; I did not
  re-test `git clone`/`rsync` variants this round (same mtime-equalizing class, covered
  by the test).
- This check verifies the frozen FINAL's observable criteria at the pinned HEAD. It does
  not certify the release channels (not yet executed — the organizer's/owner's post-close
  steps), and it does not close the idea: per LE-7 a textual pass never substitutes for
  the closing consensus's own current-tree evidence and signoff conditions.

## Closure statement

Decisive verdict: **PASS** at CLI `4df0855` (source `c3baf09`), skill `b06a65a`, frozen
FINAL `120a9bf`. Unresolved conditions outside this check's scope: the closing
consensus's own requirements remain owed (zero Agreed fixes; all-✅ signoffs or a recorded
operator ruling for any 🟡 under LE-11; citation of this verdict); the release and the
attended core publish remain the organizer's/owner's actions and are not performed. No
source, FINAL, peer, consensus, organizer, or inbox artifact was edited by this check;
this file is its only output. No commit, merge, tag, release, install, or publish action
was taken.
