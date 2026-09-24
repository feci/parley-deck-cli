---
agent: claude-1
idea: meta-protocol-change-lean-organizer
review-round: 2
date: 2026-09-24
reviewed-commit: 118b245
responding-to: [claude-1/review/round-01, kimi-1/review/round-01]
---

## Summary

Fix-up cycle 1 is substantively good work. I re-ran a **fresh full-scope refutation pass** — not a
fix checklist — over both worktrees at the reviewed commits, and every one of my 19 round-01
findings is genuinely repaired at the mechanism level, including the release-blocking CRIT-1: the
protocol's own Phase-3 template now triages `ready`, the deck corpus is back to its pre-idea
**9 of 80** malformed count, and `MissingConsensusSections` no longer exists in the tree. `parley
wait` returns **0** on this deck for the first time, `packet check` now catches hostile omits of the
phase-pinned §15 blocks, the brief resolves phase 5/8 instead of the §15-free phase-0 packet, and
the F4 measurement vector reproduces **to the byte** at the test path while my own Method-B render
reproduces the R1 path-length mechanism 1:1 at a third path (+16 B for a 16-character-longer deck
path). Both suites are green when I build and run them myself.

Against that, the fresh pass found three MAJOR issues that round 1 did not reach. Two are unmet or
mis-stated acceptance elements — the C.3/R-2 requirement to record the body **against the measured
floor** is satisfied by a number that is neither the floor nor a correct omission-set total, and the
new tri-platform CI leg runs `go test ./...` with no `-timeout` against a package that takes 589.7 s
of Go's 600 s default. The third is a defect the organizer independently hit in live use this round
and explicitly asked reviewers to settle: `parley wait --json` writes a human-readable trailer to
**stdout**, so its output is not parseable JSON on the exit-0 and exit-3 paths. I determined the
streams directly (stderr = 0 bytes) and record the answer below.

**I do not consider this ready for a zero-fix closing consensus,** but the distance is now small and
narrow: no CRITICAL, nothing that threatens the ratified gating properties, and no scope creep.

**Provenance (§15.2).** Every verdict below is `PRIMARY` unless tagged otherwise: each is a check I
executed in my own isolated worktrees, with the command and relevant output quoted. I ran no
release, version, publish, tag, merge or global-install action, and did not run `protocol publish`.
Inferences that are not executions are labelled as inferences.

**Protocol context attestation (this review):**

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase6-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

`shasum -a 256` of that file returns `8ce83cde…3a9db7` (109,928 B), and `shasum -a 256` of
`parley-deck/COOPERATION.md` at my reviewed tree returns the same value — the attested packet **is**
the live authority at the reviewed commit. PRIMARY.

**Tree provenance.**

| Repo | Review worktree | Commit | Tree | Dirty |
|---|---|---|---|---|
| CLI | `worktrees/lean-organizer-review-claude-1` | `118b2453ec0f12be063148fa575c2ccd7f2e8233` | `df104596b42cda10334197ca99788f0b4c1582db` | no |
| Skill | `worktrees/lean-organizer-review-claude-1-skill` | `a820dc7fbec80855b5112b97735c890eff1b6f83` | — | no |

Both worktrees are my own round-01 pair, refreshed with `git checkout --detach` to the round-2
commits; my round-01 evidence under `/tmp/claude1-review/` is preserved, round-2 evidence is under
`/tmp/claude1-r2/`, and I altered no other participant's worktree. CLI bases: `b37f7ef` (pre-impl)
→ `86d028b` (impl) → `3c97f44` (round-1 record) → `64a622c` (fix-up source) → `118b245` (this
record). Skill base `d1e57d5` → `0f513f6` → `a820dc7`.

My build: `go build -o /tmp/claude1-r2/parley ./cmd/parley` at `118b245`, clean tree → sha256
`33fdfcd2db59b42159db898940e46fbf9233a646e2fdb16099a99c1c9e2dcc86`, `go version -m` →
`vcs.revision=118b2453…`, `vcs.modified=false`, `go1.27.1 darwin/arm64`. Rebuilding from the same
tree a second time is **byte-identical** — see the DF-4 result under Refutation attempts.

Scope check first: `git diff --stat 3c97f44 64a622c` changes 22 source/test files plus deck
artifacts, and **every** source file maps to a signed fix (F1–F21). No unexplained source change,
no new scope.

## Position changes since prior review round

1. **CRIT-1 withdrawn as resolved, and my ❌ BLOCK verdict is withdrawn.** My round-01 verdict was
   "❌ BLOCK — CRIT-1 must be resolved before this can be called complete or released." It is
   resolved (evidence below), so the block falls. My round-2 position is 🟡: findings remain, none
   is a blocker.
2. **MIN-2 re-filed rather than closed.** I accepted F10 at signoff as written; executing it shows
   the *finding's substance* survives in two states, one of which is this idea's live state right
   now. I am not withdrawing MIN-2; I re-file the residue as `R2-MIN-1`. This is a change from the
   implicit expectation in my signoff that F10 would close it.
3. **My R1 reservation is discharged.** I wrote that "F4 should state the method and path beside the
   numbers". It does, in both methods, and I reproduced both. I withdraw R1 as satisfied.
4. **My R2 reservation is discharged.** DF-1..DF-4 exist as real, git-tracked `status: candidate`
   slugs. My conditional concurrence with F8 and the corrected phase-8 disposition is now
   unconditional for this release.
5. **On DF-4 I have moved from "unexplained" to "explained".** In round 01 I filed the binary-hash
   variance as an unreproducible-record problem (MIN-6) and the consensus carried the three-way
   hash divergence to DF-4 as an open question. I can now close the mechanism question with a
   direct experiment (below): it is the embedded build path, and `-trimpath` removes it.

## Refutation attempts

Refutation-default (LE-1): for each FINAL acceptance criterion and each fix-up claim I tried to
construct a failing case. Unsuccessful attempts are recorded as such — they are evidence the
property held. I did not rely on the shipped tests as proof of the properties they assert; where a
test mirrors the implementation I re-derived the property through the real entrypoint.

### A — pure organizer default

- **A.4 / CRIT-1 counterexample, re-run verbatim.** I rebuilt my round-01 fixture from scratch: a
  deck whose `consensus.md` is the protocol's own Phase-3 template (`COOPERATION.md:375-384`,
  de-indented, placeholders substituted, both participants ✅ ACCEPT).
  `parley consensus status --dir /tmp/claude1-r2/fx-template tpl-deck` → `Consensus: ready`, exit 0.
  At `3c97f44` the identical fixture returned `malformed / missing required consensus section(s):
  Drafter position changes, Alternatives disposition`. **Counterexample no longer reproduces.**
- **Corpus sweep, my own script, my own binary.** 80 `consensus.md` in this deck; **9 malformed** —
  exactly the pre-idea base count I measured at `b37f7ef` in round 01 (was 79/80 at `3c97f44`). The
  9 are the same historical set (`integrate-parley-bidding-addon`, `loop-engineering-research`,
  `meta-protocol-change-devx-speed`, `parley-design-skills`, `protocol-restructure-appendices`,
  `readme-skill-catalogue`, `skills-cli-install-path`, `sync-skill-protocol-fallback`,
  `track-aware-driver`). **Could not break. PRIMARY.**
- **Is the gate really gone, or moved?** `grep -rn "MissingConsensusSections" --include="*.go" .` →
  **zero hits**; the function no longer exists. `RequiredConsensusSections` has exactly two
  consumers: `buildConsensusDraftPrompt` (`driver_consensus.go:133`) and `designDraftTemplate`
  (`consensus.go:834`). `RequiredFinalSections` likewise feeds prompt + scaffold only;
  `ValidateFinal` (`finalsections.go:92`) checks status/slug/`FinalIsScaffold`, not sections. The
  substring-vs-heading issue I raised is moot — there is no gate to be lenient. **Verified. PRIMARY.**
- **Parity test is not mirroring.** `TestConsensusDraftPromptScaffoldParity`
  (`facilitator_test.go:121`) drives the real `consensus.Draft` entrypoint over a real fixture deck
  and asserts both directions — every constant section present as a heading line, **and** no `## `
  heading outside the constant set. That is a genuine drift gate. **Holds. PRIMARY.**
- **A.1 preflight matrix, four adversarial shapes.** facilitator ∈ participants without the flag →
  exit **3**, message naming `facilitator: codex-1`, `participants:` and
  `facilitator_participates: true`; **not waivable** with `--yes` (still 3); with the flag → 0;
  facilitator ∉ participants → 0; absent field → 0. **Could not break. PRIMARY.**
- **A.3 / F21 role ineligibility at the driver level.**
  `TestFixtureAutoDriveNeverLaunchesFacilitatorForCodeRoles` drives real `driver.Advance` over real
  fixture agent processes and asserts on the event log's `agent.started` trail. I read it to confirm
  it is not a unit stub: it launches, writes `IMPLEMENTATION.md` and `review/round-01/kimi-1.md`,
  and asserts `codex-1` never appears. This is the shape FINAL A named and kimi-1's K1-F3 asked for.
  **Holds. PRIMARY.**
- **A.5 absent-field byte-identical run plan.** `TestPlanByteIdenticalWithAbsentFacilitatorField`
  (`runplan_test.go:247`) — PASS. The named acceptance shape now exists. **Holds.**
- **A.5 §9.0 sentence (kimi-1's K1-F1).** Present exactly once in all three copies, correctly
  positioned inside §9.0 immediately before "Then proceed with the per-agent session-start
  checklist". **Verified. PRIMARY.**

### B — `parley wait` + PhaseDigest

- **MAJ-1 counterexample.** `parley wait --dir <review-worktree> --idea
  meta-protocol-change-lean-organizer --for round --timeout 8s` → full digest, `note: pre-existing
  unanswered to-user escalation for this idea (arrived before this wait; reported, not blocking):
  claude-to-user_…_driver-error.md`, `wait: boundary reached (round complete)`, **exit 0**. Nine
  `-to-user_` notes sit in this inbox; the six-week-old cross-idea
  `claude-1-to-user_fixup-budget_…` note is correctly invisible, and this idea's `blocking: no`
  core-publish note is correctly filtered. **Broken at `3c97f44`; fixed. PRIMARY.**
- **Was exit 4 simply disabled?** No. On my fixture deck I created a qualifying note
  (`idea: wfx`, `blocking: yes`) **three seconds into** a 12 s wait → **exit 4**, `wait: blocking
  escalation (new unanswered to-user inbox note for this idea): alpha-1-to-user_wfx_urgent.md`.
  **The safety property survives. PRIMARY.**
- **MAJ-2 (historical `driver.error`).** `driverErrorEventSince` snapshots the event count at wait
  start; the live deck's recovered historical error is now a `note:` annotation, not an exit.
  **Fixed. PRIMARY.**
- **MAJ-3 counterexample.** `--for implementation` against `status: ready-for-review` → **exit 0**,
  `wait: boundary reached (implementation published)`. An unrecognised status (`half-done`) → exit 3
  with the condition **named**: `outstanding: implementation status \`half-done\` is not a
  recognised ready state`. Both halves fixed. **PRIMARY.**
- **Fail-open probe (SUCCEEDED).** A `to-user` note for this idea with **no frontmatter**, touched
  during the wait, is neither blocking nor annotated — silently dropped. See `R2-MIN-2`.
- **`--json` parseability probe (SUCCEEDED).** See `R2-MAJ-2`; I isolated the streams and answer the
  organizer's open question there.
- **Next-action probe (SUCCEEDED).** Two states still resolve to `await implementation`. See
  `R2-MIN-1`.
- **B timeout contract.** `--timeout 31m` on `deliberation` → rejected naming the `30m0s` ceiling;
  `29m` accepted. `waitDefaultMS = 25*60*1000` (`wait.go:50`) with
  `[defaults.timeouts] wait_ms = 1500000` seeded (`runtime.go:645`), per-call override. `default <
  30m && configurable` holds, and the seed comment explicitly disclaims a provider-TTL fact.
  **Could not break. PRIMARY.**
- **B determinism / no writes / no model-written field.** Two `organizer brief` runs (which embed
  the digest) `cmp`-identical; `git status --porcelain` in my review worktree stayed **empty**
  through every `wait`, `brief` and `packet` invocation. `internal/driver` and `internal/app`
  digest/wait test sets green. **Could not break. PRIMARY.**
- **MIN-1 (`fell_back`).** Live digest now shows `fell_back=false` on all five valid rows (three
  round-02, two review round-01); it was `true` on every valid round-02 row at `3c97f44`.
  **Fixed. PRIMARY.**
- **`--for any`** → exit 0 on a reached boundary. **Holds.**

### C — audience packet, slim skill, brief

- **C.1 runtime never-cut floor, hostile map.** I wrote a map naming **eight** never-cut blocks for
  facilitator omission (`## 15.`, `### 15.1`, `### 15.2`, `### 15.3`, `### 15.4`, `### 15.7`,
  `## 6.`, `## 14.`) and rendered at every kernel phase {1,2,3,5,6,7,8} with line-anchored heading
  checks. **Every block stayed present at every phase.** **Could not breach the floor. PRIMARY.**
- **MAJ-6 counterexample (the `packet check` proof obligation).** Same hostile entries, one at a
  time, through `parley protocol packet check`: `## 6.`, `## 14.`, `## 15.`, `### 15.1`,
  `### 15.5`, `### 15.7`, `### Phase 1 — Round 1 …` and `## 7.` are **all now rejected**
  (`never-cut violations:`, exit 1), while the clean map — which legitimately omits `### 11.A`,
  `### 11.B`, `### 11.C` — stays **ok**. The allowance is narrowed to transport-conditional
  subsections exactly as F6 states. **Fixed; the negative test FINAL C.1 names now exists. PRIMARY.**
- **C.1 fallbacks.** `--audience banana` → `context_mode: full`,
  `audience_fallback_reason: unknown-audience:banana`; `--audience facilitator` → `packet`;
  `--audience participant` → `full`. **Could not break. PRIMARY.**
- **MIN-4 counterexample.** `--optimize --audience banana` → body header now reads
  `flags=- audience=-` (resolved), not `audience=banana`. **Fixed. PRIMARY.**
- **C.3 gating check (R-2).** `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` re-run at my tree
  → PASS; the named-omission-set absence is line-anchored, so it cannot be satisfied by an omission
  index mention. **Could not break the gating check. PRIMARY.**
- **MAJ-4 / F4 Method A, reproduced to the byte.** `go test ./internal/app/ -run
  'TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail|TestLiveDeckFacilitatorAcrossPhases' -count=1
  -v` in my worktree logs **53,870 / 59,206 / 59,307 / 60,567 / 63,485 / 61,129 / 62,211 / 70,086 /
  72,696 B** for phases 0–8 and `--optimize baseline=65,750 B` — **identical to every figure
  recorded in `IMPLEMENTATION.md`.** Method A is path-independent because the test passes the
  32-char relative authority path. **Record is accurate. PRIMARY.**
- **F4 Method B, reproduced at a third path.** My deck path is 117 chars against the record's 101.
  Direct CLI renders give phases 0/1/2/3/4/5/6/7/8 = 53,955 / 59,291 / 59,392 / 60,652 / 63,570 /
  61,214 / 62,296 / 70,171 / 72,781 B — **exactly the recorded Method-B vector + 16 B at every
  phase**, and 117 − 101 = 16. The R1 path-length mechanism reproduces 1:1 at a third independent
  path. **R1 discharged. PRIMARY.**
- **C.3 floor recording (SUCCEEDED).** The body is not recorded against the measured floor. See
  `R2-MAJ-1`.
- **MAJ-8 / F8 blind-spot (i), re-measured live.** Line-anchored over the Method-B bodies: phase 5
  and phase 8 carry `### 15.7` but neither `### 15.5` nor `### 15.6`; phase 7 carries all three;
  phase 0 carries no §15 at all. **Identical to the table recorded in `IMPLEMENTATION.md`.** The
  measurement FINAL required is performed and recorded with the whether-it-matters analysis.
  **Fixed as recorded. PRIMARY.**
- **MAJ-5 counterexample.** Six fixture decks across the real status vocabulary:
  `implementation`→phase **5**, `complete`→**8**, `fix-up-cycle-1`→**8**, `final`→4, `open`→0,
  `abandoned`→0. The live brief for this idea now renders **packet-phase8** (it rendered
  packet-phase4 during Phase 6 at `3c97f44`). A live Phase 5–8 idea can no longer be handed the
  §15-free phase-0 packet. **Fixed. PRIMARY.**
- **C.5 brief contract.** Live: **2,631 B** ≤ 8,192 B; two runs `cmp`-identical; `git status` in the
  deck stayed empty. The brief's phase-8 body is 2,984 B larger than the plain phase-8 render
  because this idea sets `flags=protocol_change`, pinning §7 — I diffed the two bodies and the only
  heading difference is `## 7. Changing this protocol`. Correct, not a defect. **Could not break.
  PRIMARY.**
- **MIN-3.** The brief's opening line now reads "never stored **in the deck**" and names the
  `.parley-runtime/protocol-packets/` cache it writes outside the deck. **Fixed. PRIMARY.**
- **C.4 SKILL.md, relocation re-derived independently.** Core **17,292 B** ≤ 20,000 B; all six
  driver commands named; all six `references/` entries linked. I re-extracted **all 76 headings at
  every level** from `git show d1e57d5:skills/parley-deck/SKILL.md` and checked each against core
  and each reference: **0 dropped, 0 duplicated across locations**, 9 retained in core. Frontmatter
  block and `## Core Rule` **byte-identical** to `d1e57d5`. **Could not break. PRIMARY.**

### D — handoff, ledger, telemetry

- **MAJ-7 / F7 mechanism.** `commitCursor` calls `d.touchRunManifestUpdatedAt()` at
  `driver.go:224`, which calls `runmanifest.TouchUpdatedAt` (`manifest.go:231`).
  `TestCommitCursorAdvancesRunManifestUpdatedAt` drives the **real `commitCursor`** and asserts the
  advance. The former synthetic 2-hour manifest is replaced by
  `TestUsageIngestAttributionWindows`, which builds records through `runmanifest.New/Write` +
  `TouchUpdatedAt` — the exact API the driver uses — and asserts `attributed` for the matching idea
  and `ambiguous` for another idea's window. **Mechanism fixed. PRIMARY.**
- **F7 live state, checked rather than assumed.** All four `run.json` records in the live deck are
  still `created_at == updated_at` (zero-width), and no run record has been created since the fix
  landed at `64a622c` (2026-09-24T01:48Z). So the live smoke's `attribution=ambiguous` is **honest
  and expected**, and the end-to-end path is currently proven by test rather than by a live
  artifact. I state that boundary rather than claiming more. **Disposition (iii) holds.**
- **D.2 ledger contract, end-to-end on my own fixture.** Rollout with two cumulative
  `token_count` events (12 then 120) plus a decoy `response_item` containing the literal strings
  `token_count` and `total_token_usage`: stdout **112 B** ≤ 1 KB; **events=2** (decoy filtered by the
  type check); exactly **one** ledger row; the six `total_token_usage` keys verbatim;
  `total_tokens=120` — the last cumulative value, not the first and not a sum; re-ingest →
  `idempotent no-op`, ledger unchanged; method header states explicit-args + run-window attribution
  and "A slug appearing inside a session file is never attribution evidence". **Could not break.
  PRIMARY.**
- **D.2 no count-accepting flag.** Eight invented flags (`--tokens --input-tokens --total-tokens
  --count --usage --input --output --total`) all rejected with `flag provided but not defined`.
  **Could not break. PRIMARY.**
- **D.4 kimi telemetry.** `case "kimi"` present at `telemetry/usage.go:307`; four tests pass
  including `TestKimiCoverageNoneWhenNoUsageEmitted` (the honest-degradation path) and
  `TestKimiUsageIgnoresForeignShapes`. **Holds. PRIMARY.**
- **MIN-5 / F13.** `BuildPhaseHandoffRecord` takes the run dir (the `"parley-deck"` derivation is
  gone); `TestPhaseHandoffRoundTrip`, `TestBuildPhaseHandoffRecordTakesRunDir` and
  `TestCommitCursorWritesPhaseHandoffRecord` all pass. **Fixed.**
- **NIT-1 / NIT-2 / K1-F6c.** `containsBytes`/`indexOfBytes` gone (0 hits), `bytes.Contains` used;
  `var _ = strconv.Itoa` and the `strconv` import gone (0 hits); `streamLines` has a `dropping`
  mode so an oversized line's tail is never yielded. **Fixed. PRIMARY.**
- **NIT-3 / F18.** `TestOptimizeAudienceBranchSharesReasonsSet` pins the merged branch's reason sets
  across kernel phases. **Fixed.**

### Cross-cutting — staged core, three copies, suites, binary

- **Staged core, verified independently and by decomposition.**
  `~/.parley/staging/COOPERATION-2.13.0.md` is **109,772 B**, sha256
  `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`, and is **byte-identical**
  (`cmp` silent) to the skill's template-form third copy at `a820dc7`. I then decomposed it rather
  than eyeballing hunks: `git show d1e57d5:skills/parley-deck/references/COOPERATION.md` (108,244 B)
  **plus exactly five diff hunks** — Quickstart facilitator row, §4 Phase 5, §4 Phase 6, §9.0
  sentence together with §9 item 1, §11 advisory — **equals the staged core**. And
  `diff(core 2.10.0, skill base)` is five hunks that are exactly the 2.11.0 (§15.6, §15.7) and
  1.48.0 (LE-7/LE-11 bullet, goal-check paragraph, §9 item 1) change sets. **Three change sets, 0
  unexplained.** My direct `diff(2.10.0, staged)` yields 9 hunks against the record's 11 — pure
  context-granularity, same content. **Could not break. PRIMARY.**
- **Staged core project zones.** Placeholder header byte-intact (`<workspace-name>`,
  `<transport-choice>`, `<YYYY-MM-DD>`); **zero** occurrences of `claude-1|kimi-1|zcode-1|codex-1`;
  §2 roster-view and host-handle tables carry header rows only, where the deck view carries six
  host-handle rows. The single `Protocol synced:` occurrence is §9.0 prose describing the sync
  mechanism, not a header value. **Could not break. PRIMARY.**
- **Three copies + staged core carry identical hunks.** All six hunk fragments present in all four
  files. The §15 region is **byte-identical** across deck view, `internal/protocol/defaults/` and
  the staged core (same sha256, 8,056 B each). **Could not break. PRIMARY.**
- **Both suites green, run by me.** CLI `go test ./... -count=1 -timeout 2400s` → **31/31 packages
  ok, 0 FAIL**. Skill `npm ci && npm test` → **399 pass, fail 0, exit 0**, including the
  `lean-organizer` hunk assertions and the addon-manifest check. **Holds — but see `R2-MAJ-3` for
  what happens without `-timeout`.**
- **MIN-6 / F14 binary provenance.** My clean build at `118b245` reports `vcs.modified=false`; the
  round-01 problem (a hash recorded from a dirty tree at the pre-implementation base) is gone.
  **Fixed.**
- **DF-4 mechanism, settled by experiment (a contribution, not a finding).** Two builds from the
  same clean tree at the same path are **byte-identical** (`33fdfcd2…` twice), so the toolchain is
  deterministic. But at the *same* commit `64a622c`, clean, same Go version, from **three different
  directories** I get **three different** hashes (`e9c1d85a…`, `96b27278…`, and the implementer's
  recorded `8bd06646…`). Adding `-trimpath` makes two different directories produce **byte-identical**
  binaries (`0ba8193b…` twice). **The entire hash variance is the embedded build path, and
  `go build -trimpath` removes it** — the same class of mechanism as R1's packet path-length effect.
  DF-4 can be answered in one line. **PRIMARY.**
- **Attended publish gate.** The escalation note still carries the exact command
  `parley protocol publish --version 2.13.0 --from ~/.parley/staging/COOPERATION-2.13.0.md`;
  `~/.parley/protocol/core/` still holds only `2.10.0`. I did not run it and did not manufacture a
  TTY. **Holds.**
- **R2 reservation discharged.** DF-1 `meta-protocol-change-consensus-duty-gates`, DF-2
  `meta-protocol-change-facilitator-integrity-phase-coverage`, DF-3
  `facilitator-packet-per-phase-bounds`, DF-4 `release-binary-reproducibility` all exist, are
  git-tracked, and are correctly inactive `status: candidate` records (LE-10) that staff no quorum.
  **Verified. PRIMARY.**
- **Edit hygiene.** `git diff 64a622c 118b245 -- review/consensus.md` touches only the four DF slug
  lines inside the drafter's own `## Deferred follow-ups` section; **no signoff block is modified**.
  Correct.
- **A false lead I chased and dropped.** Three untracked files appeared in my skill worktree during
  the suite run. They were **mine** — a word-splitting slip in one of my own verification commands
  writing through an unquoted path containing spaces — not a defect in the skill tests. I removed
  them; both worktrees end clean. I record this so the near-miss is auditable rather than silent.

## Findings

### [MAJOR] R2-MAJ-1 — the facilitator body is still not recorded against the measured floor, and the number that stands in for it is mislabelled and undercounted

FINAL's acceptance table for C requires: "map-derived byte count recorded against **both** the
measured floor (≈ 42 KB with §2) and the ≤ 70,000 B guardrail in the same test run", and open item
14 (R-2) repeats it. What `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` actually records
(`internal/app/facilitator_packet_live_test.go:96-119`) is:

```
R-2 measurement (phase 1 / deliberation / github-pr): facilitator body=59206 B;
named-omission-set bytes=15759 B; +§2(3977 B) reference=19736 B; guardrail=70000 B;
--optimize baseline=65750 B
```

Three separate problems:

1. **The floor is absent.** FINAL distinguishes two quantities explicitly — "named omission set
   27,420 B; floor ≈ 38.1 KB at phase 2 / ≈ 38.0 KB at phase 1 without §2; **≈ 42.1 KB with §2
   retained**". The number logged beside the guardrail is 19,736 B, which is neither.
2. **The code says the opposite of what it computes.** `:97` reads "the whole-section floor is the
   sum of the source bytes of the named omission set", and the variable is named `floor`. A floor is
   what the body *retains*; this sums what it *omits*. `:107` then concedes "claude-1's measured
   floor included the subsections; approximate the §2 retention by adding §2's bytes" — an
   approximation of a different quantity is recorded under the name of the required one.
3. **Even as an omission-set total it undercounts.** `byHeading[locator]` takes one `Parse` block,
   and `Parse` splits `### N.M` subsections into their own blocks. I probed this directly
   (temporary test in `internal/protocolpacket`, removed afterwards; tree verified clean):
   `## 12.` contributes **273 B** and `## 13.` **522 B**, while their subsections — which the map
   does omit — are **17 blocks / 11,903 B** that the sum silently excludes. My own whole-section
   measurement over the source gives 23,223 B for the seven top-level omitted sections alone.

The **gating** half of C.3 (named-omission-set absence) and the guardrail itself both hold — this
does not threaten them. What fails is the recorded comparison the quorum ruled on, which is exactly
the class of defect F4 was raised to fix.

**Suggested fix:** compute the floor as the retained whole-section total (or reuse the body's own
omission index, which already carries each omitted block's line span), log it as `floor`, and fix
the comment and the variable name. If the quorum prefers to keep the omission-set number as well,
log it under its own label and make it whole-section-correct.

### [MAJOR] R2-MAJ-2 — `parley wait --json` writes a non-JSON trailer to stdout, so the machine-readable path is unusable on the exit-0 and exit-3 routes

`internal/app/wait.go:173-179`: `printWaitDigest(out, digest, annotations, *jsonOut)` emits the JSON
envelope to `out`, and the very next line writes `wait: boundary reached (…)` / `wait: timeout after
…` to **`out`** as well. Only the exit-4 routes send their reason to `stderr`.

I isolated the streams (this is the question codex-1's
`inbox/codex-1-to-all_…_wait-observation-02.md` asks reviewers to settle):

```
$ parley wait --dir <fixture> --idea wfx --for consensus --timeout 3s --json 1>out.txt 2>err.txt
exit=3
stderr bytes: 0
stdout last line: wait: timeout after 3s; outstanding: awaited condition not yet reached — inspect the digest
$ python3 -c "import json;json.load(open('out.txt'))"
json.decoder.JSONDecodeError: Extra data: line 51 column 1 (char 1179)
```

**stderr is empty.** The trailer is on stdout. So the answer to the organizer's open question is:
the wrapper is not merging streams — `parley wait --json` genuinely emits invalid JSON on stdout on
both the exit-0 and exit-3 paths. The exit-4 path is clean and does parse.

This is not a fix-up regression — the same shape exists at `3c97f44` (`git show
3c97f44:internal/app/wait.go:161-162`) and neither reviewer filed it in round 1. The fix-up did
change the envelope from a bare `PhaseDigest` to `{notes?, digest}`, so any consumer written against
round-1 behaviour also needs updating; nothing has shipped, so that is a note rather than a break.
FINAL B.3 offers `--json` precisely so an organizer can branch mechanically, and codex-1 hit the
failure (`jq: parse error: Invalid numeric literal at line 124, column 5`) in ordinary use this
round.

**Suggested fix:** in `--json` mode either send the terminal status line to stderr, or fold it into
the envelope as a `result` field (`{"notes":[…],"digest":{…},"result":"timeout","outstanding":[…]}`)
and print nothing else on stdout. Add a test that pipes `--json` stdout through `encoding/json`
decode on all four exit paths.

### [MAJOR] R2-MAJ-3 — the new tri-platform CI leg runs `go test ./...` with no `-timeout`, against a package that uses 589.7 s of Go's 600 s default

`.github/workflows/tests.yml` (new in this idea) runs, on `ubuntu-latest`, `windows-latest` and
`macos-latest`:

```yaml
      - name: Test
        run: go test ./... -count=1
```

No `-timeout`, so Go's 10-minute per-package default applies. On this machine, in my green run:

| package | duration |
|---|---:|
| `internal/trajectory` | **589.692 s** |
| `internal/app` | **521.712 s** |
| `internal/runner` | 120.696 s |

Ten seconds of headroom. And running **exactly the command FINAL's validation plan names** —
`go test ./...`, no timeout flag — failed for me:

```
panic: test timed out after 10m0s
FAIL	parley-deck-cli/internal/app	600.146s
panic: test timed out after 10m0s
FAIL	parley-deck-cli/internal/trajectory	600.092s
```

I only got a green suite by passing `-timeout 2400s`, which is also what the implementer's record
uses. kimi-1 ran without the flag in round 1 and passed, so the outcome is load-dependent — which is
the point: the margin is ~2%.

FINAL risk 8 states the Windows leg exists "so the release ships Windows assets honestly", and the
cross-cutting acceptance row requires "both worktrees' test suites green incl. the new tests". A CI
leg configured to time out does not deliver honest Windows coverage. *(That GitHub's hosted runners
are slower than a local Apple-silicon machine is an inference, not an execution — I cannot measure
them from here. The 589.7 s measurement, the absent flag, and my own failing run are PRIMARY.)*

**Suggested fix:** add `-timeout 45m` (or similar) to the workflow's test step and to FINAL's named
validation command in `IMPLEMENTATION.md`; independently, `internal/trajectory` and `internal/app`
are worth a look — a 9-minute package is a latent CI failure on any runner.

### [MINOR] R2-MIN-1 — the next-action line still says `await implementation` in at least two states where that is not what is awaited, including this idea's current live state

`internal/driver/phasedigest.go:188-199`. F10 added one branch (`ReadyForReview && d.Review == nil`
→ `NextAwaitReviewArtifact`) and `NextAwaitImplementation` remains the unconditional fallback. Two
states reach it wrongly:

**(a) Fix-up published while the previous review round exists** — recorded by the implementer as an
F10 residual, and it is this deck right now:

```
implementation: present=true status=fix-up-cycle-1 implementer=zcode-1
next: await implementation
```

The implementation is published and awaiting re-review; I am the artifact it is waiting for. The
same line appears in `parley organizer brief`, the documented post-compaction re-orientation surface.

**(b) Rounds complete, no `consensus.md` at all** — unrecorded. `BuildPhaseDigest`
(`phasedigest.go:134`) only populates `d.Consensus` when `consensus.Status` succeeds, so an absent
`consensus.md` leaves it `nil`, `NextAwaitConsensus` at `:186` is unreachable, and a deck sitting at
Phase 2/3 is told to "await implementation" — skipping Phases 3–4 entirely. Reproduced on a fixture
with a complete round-01 and no consensus:

```
implementation: present=false status= implementer=
next: await implementation
```

The fixed enumeration (FINAL B.2) is respected; the *selection* is wrong. I accepted F10 as worded
at signoff, so I record this as my own mis-scoping as much as an implementation gap — but MIN-2's
substance is not closed.

**Suggested fix:** return `NextAwaitReviewArtifact` whenever the implementation is ready and the
latest review round is complete but the implementation is newer than it; and emit
`NextAwaitConsensus` when rounds are complete and no `consensus.md` exists (a nil consensus section
is a real state, not an absent one).

### [MINOR] R2-MIN-2 — F2 fails open: a `to-user` note with unparseable or `idea:`-less frontmatter is silently ignored, neither blocking nor annotated

`internal/app/wait.go:297-305`. `qualifyingEscalationNotes` does `if merr != nil { continue }` on a
frontmatter read failure, and `if id := …meta["idea"]; id != idea { continue }` — so a note with no
frontmatter has `id == ""`, never matches, and disappears. Reproduced: a note named
`alpha-1-to-user_wfx_nofm.md` containing `URGENT: stop, the release is wrong.` and touched during
the wait →

```
  no-frontmatter note -> EXIT=3  (3=ignored, 4=blocked)
```

At `3c97f44` the same note exited 4 (filename matching). The three qualifiers F2 implements are
right, and a well-formed note behaves correctly; the issue is the degraded input. It also sits
against B.2's own stated discipline — "fail-closed `unparsed`, never a guess" — and against the
digest's practice of surfacing `unparsed` rows rather than dropping them. The silent part is the
worst of it: there is not even a `note:` annotation.

**Suggested fix:** when a `*-to-user_*.md` cannot be parsed or carries no `idea:`, emit a `note:`
annotation naming the file and why it was not evaluated. Blocking on it would be defensible too, but
annotating is enough to remove the silence.

### [MINOR] R2-MIN-3 — the owner-facing core-publish escalation note still describes the pre-restage file

`parley-deck/inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_core-publish.md` is the
artifact the owner acts on for the one irreversible, write-once action in this release. It still
says:

```
**Staged file:** `~/.parley/staging/COOPERATION-2.13.0.md` (109,507 B)
**Composition (verified 2026-09-23):** …
```

F20 restaged that file on 2026-09-24; it is now **109,772 B** (sha256 `fc907e59…`, mtime
`Sep 24 03:35`), and the re-verification is recorded in `IMPLEMENTATION.md` as 2026-09-24. F20's
label correction *was* applied to the change-set list (it now reads "§9.0 pure-organizer default
sentence, §9 checklist item 1"), but the size and the verification date were not. No sha256 is
recorded anywhere in the note.

The command itself is correct and points at the right path, so the right bytes would be published —
the harm is that the note's own integrity statement no longer matches the file, and an owner who
checks cannot distinguish a legitimate restage from tampering.

**Suggested fix:** update the size to 109,772 B and the verification date to 2026-09-24, and add the
staged sha256 `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f` so the owner can
verify they are publishing the reviewed bytes.

### [MINOR] R2-MIN-4 — `## Deviations from FINAL.md` says "None in scope or mechanisms", but F1 deliberately dropped a mechanism FINAL A.4 names

FINAL A.4 requires "both section lists derived from one `protocol.RequiredConsensusSections` /
`RequiredFinalSections` constant and a **parity test** proving prompt **and gate** read the same
value". After F1 there is no gate reading either constant (`MissingConsensusSections` is gone;
`ValidateFinal` checks status/slug/scaffold only), and the parity test proves prompt ↔ scaffold ↔
constant instead.

This is the right outcome and I signed for it — the acceptance-table row for A asks only that the
*emitted prompt* carry the sections "derived from the shared constant (parity test)", which is met,
and the removal was directed by a unanimously signed fix. But `IMPLEMENTATION.md:418-421` opens with
"None in scope or mechanisms", and Phase 5 requires that "any unavoidable deviation is logged in
`IMPLEMENTATION.md` — not silently absorbed". The F1 entry explains the change; the deviations
section, which is where a reader looks for exactly this, does not mention it.

**Suggested fix:** one bullet under `## Deviations from FINAL.md` recording that A.4's gate-side
parity was removed by signed fix F1, that the acceptance-table row is met by prompt ↔ scaffold ↔
constant parity, and that hard-gating the duty sections is DF-1.

### [NIT] R2-NIT-1 — `outstandingAgents` names "IMPLEMENTATION.md not filed" but has no parallel for a missing `consensus.md`

`internal/app/wait.go:266-286`. The implementation branch names its blocking condition precisely
(F3's improvement), but the consensus branch only appends `d.Consensus.Missing` when
`d.Consensus != nil`; an absent `consensus.md` falls through to the generic `awaited condition not
yet reached — inspect the digest`. Same shape as `R2-MIN-1(b)`, same one-line fix: add
`"consensus.md not filed"` when the scope is consensus and the section is nil.

### [NIT] R2-NIT-2 — F2's "arrival" signal is file mtime, and the driver rewrites escalation notes in place

`internal/driver/loop.go:353-354` builds the note name as
`fmt.Sprintf("claude-to-user_%s_%s.md", ideaSlug, topic)` and writes it with `os.WriteFile` — an
unconditional overwrite. A second `driver.error` for the same idea replaces the first, still
unanswered, escalation and gives it a fresh mtime. That happened in this run:
`git diff 3c97f44 64a622c -- inbox/claude-to-user_…_driver-error.md` shows the round-02
historical-worktree blocker replaced by `draft FINAL.md: context canceled`, with `phase:` flipped
`round` → `consensus`.

Consequence for B: an in-place rewrite (or a bare `touch`) of a pre-existing escalation makes it
look "new" to `wait` and blocks; the original text is gone from the file. This fails loud, so it is
safe, and the canonical record survived — `organizer-notes.md:24` still carries the full
historical-worktree path and the mislabelled-author note. Worth a comment at the mtime comparison
recording that arrival is approximated by mtime and why that is acceptable; the `from: claude`
mislabel is already FINAL-recorded testimony and out of scope.

### [NIT] R2-NIT-3 — `organizer-notes.md` ends with literal `\n` escapes instead of newlines

The DF-1..DF-4 list on the last line reads
`- DF-1: ideas/…/00-prompt.md\n- DF-2: ideas/…/00-prompt.md\n- DF-3: …\n- DF-4: …\n` as a single
line with backslash-n escapes, and the file has no trailing newline. This is the organizer's own
artifact rather than the implementation, but it is part of the audit trail this idea will close on.

### [NIT] R2-NIT-4 — the recorded staged-core §15 region size does not reproduce

`IMPLEMENTATION.md` records the restage recheck as "§15 region (8,041 B) … byte-equal to the deck
view". Extracting from the `## 15. Verification integrity` heading line to EOF, I measure **8,056 B**
in the staged core, the deck view and `internal/protocol/defaults/` — all three with the identical
sha256, so the substantive byte-equality claim verifies. The 15 B delta is presumably a different
extraction boundary; given that stale recorded figures were round 1's central record finding, the
convention is worth stating beside the number.

## Updated findings

Disposition of every finding I filed in `claude-1/review/round-01`. All 19 are accounted for; none
is withdrawn for any reason other than being fixed, and each disposition is backed by a check I ran
at `118b245`.

| ID | round-01 severity | Disposition at `118b245` |
|---|---|---|
| CRIT-1 consensus-section gate | CRITICAL | **RESOLVED.** Template → `ready`; corpus 9/80 (base count); `MissingConsensusSections` absent from the tree. |
| MAJ-1 any `-to-user_` note blocks | MAJOR | **RESOLVED.** Live `wait --for round` → exit 0; new qualifying note still → exit 4. |
| MAJ-2 historical `driver.error` poisons wait | MAJOR | **RESOLVED.** Event-count snapshot at wait start; history demoted to `note:`. |
| MAJ-3 `ready-for-review` never ready; timeout names nobody | MAJOR | **RESOLVED.** Exit 0 on `ready-for-review`; unrecognised status named verbatim. Residue → `R2-NIT-1`. |
| MAJ-4 stale R-2 byte figures | MAJOR | **RESOLVED.** Method A reproduces to the byte; Method B reproduces at my third path +16 B. R1 discharged. |
| MAJ-5 brief resolves wrong phase ≥ 5 | MAJOR | **RESOLVED.** `implementation`→5, `complete`/`fix-up-*`→8; live brief now phase-8. |
| MAJ-6 `packet check` blind to phase-pinned never-cut | MAJOR | **RESOLVED.** All hostile §15.x omits rejected; §11 allowance stays legal. |
| MAJ-7 zero-width attribution windows | MAJOR | **RESOLVED as a mechanism.** `TouchUpdatedAt` wired into real `commitCursor`; synthetic test replaced by a real-API one. Live still `ambiguous` — honest; no post-fix run exists yet. |
| MAJ-8 blind-spot (i) unmeasured | MAJOR | **RESOLVED as recorded.** My independent re-measure matches the recorded table exactly. Materiality sustained; map pin → DF-2. |
| MIN-1 `fell_back` semantics | MINOR | **RESOLVED.** `fell_back=false` on all valid live rows. |
| MIN-2 next action "await implementation" | MINOR | **PARTIALLY RESOLVED — re-filed as `R2-MIN-1`.** |
| MIN-3 brief writes a packet body | MINOR | **RESOLVED.** Line names the `.parley-runtime` cache outside the deck. |
| MIN-4 rejected audience stamped into header | MINOR | **RESOLVED.** Header reads `audience=-`. |
| MIN-5 handoff `RunID` / partial round-trip | MINOR | **RESOLVED.** Run dir parameter; nine-field round-trip test. |
| MIN-6 binary from a dirty pre-impl tree | MINOR | **RESOLVED.** Clean-tree rebuild, `vcs.modified=false`; my own build agrees. Residual variance explained under DF-4. |
| MIN-7 `ResolveImplementer` unused | MINOR | **RESOLVED.** Export dropped; Decision Log and doc comment corrected. |
| NIT-1 hand-rolled byte search | NIT | **RESOLVED.** `bytes.Contains`. |
| NIT-2 strconv sentinel | NIT | **RESOLVED.** Sentinel and import gone. |
| NIT-3 optimize/audience branch seam | NIT | **RESOLVED.** Reasons-parity assertion added. |

**Dispositions the brief asked me to evaluate openly.** None of these narrowed what I inspected or
reported.

- **(i) Phase-1 59,206 B meets the guardrail; phases 7/8 at 70,086 / 72,696 B; per-phase policy is
  DF-3; method and path recorded. — I concur, with one correction of emphasis.** FINAL C.3 and the
  acceptance table scope the guardrail to the phase-1 / deliberation / github-pr body, which I
  measured at 59,206 B (Method A) and 59,291 B at my own path — ~10.8 KB of margin. The map, the
  ceiling and the never-cut floor are untouched, and open item 2's "show the bytes to the quorum"
  discipline was followed. The emphasis worth adding: **phase 7 moved from 179 B under to 86 B over**
  because F20 added the §9.0 sentence FINAL A.5 already required — a signed fix completing frozen
  scope crossed a non-binding figure at a non-ratified phase. That is exactly DF-3's question, and
  it strengthens rather than weakens the case for it. The method/path recording is accurate and I
  reproduced both methods. **Concur.**
- **(ii) §15.5/§15.6 phase-5/8 omission accepted for this release per signed F8, with measurement
  and DF-2. — I concur, and my materiality call stands.** My independent line-anchored re-measure
  matches the recorded table exactly. I still hold that showing a facilitator §15.7's table of
  duties while cutting the text of two of them is a coherence defect in a normative artifact, not a
  record-completeness nit — and kimi-1 self-corrected to substantially that position. But the
  operational impact at phases 5 and 8 is low (those duties bind the drafter at 3/6/7, where they
  are pinned), the measurement and rationale are now recorded, and DF-2 exists. **Concur for this
  release.**
- **(iii) F7's live smoke remains `ambiguous`; future transition windows claimed functional. — I
  concur, having checked the claim rather than the wording.** `touchRunManifestUpdatedAt` is called
  from the real `commitCursor` at `driver.go:224`, and the replacement test drives `commitCursor`
  itself. All four live `run.json` records are still zero-width and no run has been created since
  the fix landed, so `ambiguous` is the honest answer and the end-to-end path is currently
  test-proven rather than artifact-proven. I would rather that boundary be stated than blurred, and
  the record does state it. **Concur.**
- **(iv) F10's residual next-action is explicitly recorded. — I concur that it is recorded, and I
  do not concur that recording closes it.** The residue reproduces live on this deck and in the
  brief, and a second unrecorded state behaves the same way. Filed as `R2-MIN-1`. Recording a
  known-wrong output is the right transparency; it is not a disposition that retires the finding,
  and the review-brief rule says a disputed finding closes only by reviewer withdrawal, consensus,
  or an operator ruling. **Partially concur.**
- **(v) Binary rebuilt from a clean source commit with flags/VCS provenance. — I concur, and I can
  now close the question behind it.** `vcs.modified=false` at a clean commit is correct and fixes
  MIN-6. The remaining puzzle — five different hashes over behaviour-identical trees — is fully
  explained by the embedded build path: same commit from three directories gives three hashes, and
  `-trimpath` collapses them to one. **Concur, and DF-4 has its answer.**

## Responses to other reviewers

### @kimi-1

Your round-01 review was careful and adversarial where it counted, and one of your findings was the
single thing my pass structurally could not reach.

- **K1-F1 (the §9.0 sentence) — you were right and I missed it; verified fixed.** My cross-copy check
  asserted that the six hunks that exist are present in all three copies, which can never detect a
  *seventh* hunk FINAL required and nobody wrote. That is a real methodological gap in my round-1
  work, and it is worth generalising: presence checks over the implemented set cannot find omissions
  against the specified set. F20 landed the sentence in all three copies, correctly positioned in
  §9.0, and the restaged core carries it. **Fixed; I concur with your finding and your fix.**
- **K1-F2 ≡ my MAJ-4 — merged and resolved.** Our independently reproduced vectors agreed to the
  byte in round 1 and the corrected record now reproduces to the byte at the test path. Your reading
  that the substance was unchanged (phase 1 within guardrail; phase 8 above for a structural reason)
  was right.
- **K1-F3 (criterion-A test shapes) — you were right, and F21 delivered the real shape.** The new
  `TestFixtureAutoDriveNeverLaunchesFacilitatorForCodeRoles` is a genuine driver-level auto-drive
  with real fixture agent launches asserting on the event log, not an ops-boundary unit test, and
  `TestPlanByteIdenticalWithAbsentFacilitatorField` closes the run-plan half. I verified both by
  reading them, not just by their green result. **Concur.**
- **K1-F4 ≡ my MAJ-1/MAJ-2 — your §15.1 SELF-CORRECTION was the right call, and the outcome
  vindicates it.** You initially read the acceptance-table row as tolerating the broader
  any-unanswered-note semantics and filed NIT; I read FINAL B.3's "a **new** unanswered … arrives"
  as binding and filed two MAJORs. You reproduced my live-deck evidence and corrected. F2's
  implementation is now the stricter reading and `wait` returns 0 on this deck. I note this not to
  re-litigate but because it is the clearest example in this idea of a conflict closing the way §15.3
  says it should — by evidence, not by count.
- **K1-F5 ≡ my MAJ-8 — we converged, and I think your partial self-correction was correct.** Your
  "practical impact is nil" and my "coherence defect" were both partly right; the recorded F8
  analysis lands between them, and DF-2 carries the pin.
- **K1-F6a/b/c — all three fixed and verified** (run-dir parameter, `strconv` sentinel and import
  gone, `streamLines` drops oversized lines whole).
- **One round-01 inference of yours that did not hold, now settled.** You wrote: "I treat the hash
  delta as build-environment, not source, because behavior is byte-identical on a deterministic
  render." Your *conclusion* was right, but the reasoning could not establish it — behavioural
  agreement on one render does not show the source trees matched. In fact they did not: my MIN-6
  found the implementer's binary reported `vcs.revision=b37f7ef` and `vcs.modified=true`, i.e. the
  pre-implementation commit with uncommitted changes. That mattered, and F14 fixed it. The residual
  variance is now fully explained — same commit, clean, three directories, three hashes;
  `-trimpath` collapses them. I am recording the mechanism here rather than as a finding, since
  DF-4 owns it. No self-correction is being asked of you at this point; the record simply completes.
- **Your open question 3 (idea-scoped vs any-unanswered semantics)** is answered by F2 and the skill
  text now documents the implemented semantics accurately. My `R2-MIN-2` is the leftover: the
  qualifier set is right, but a note the parser cannot read escapes all three qualifiers silently.
- **Your open question 4 (reproducible-build story)** — see the `-trimpath` result; that is your
  question closed if the quorum wants it closed.

I have no unresolved conflict with you in this round. Where we differed in round 1, the differences
closed on evidence; nothing in your review needs a counter-proposal from me now.

## Open questions

1. **R2-MAJ-1 disposition.** Is the intent that the recorded companion number be the retained
   whole-section floor (FINAL's ≈ 42 KB) recomputed at the current HEAD, or does the quorum prefer
   to record the corrected whole-section *omission* total and amend the acceptance wording? The
   first matches FINAL as frozen; the second needs a recorded deviation. Either way the comment and
   the variable name should stop calling an omission sum a floor.
2. **R2-MAJ-3 scope.** Is adding `-timeout` to the CI step and to the FINAL-named validation command
   in-cycle work, or a DF? I lean in-cycle: it is a two-token change on a shipped file, and the
   release's Windows-honesty claim rests on that leg being able to finish. The separate question —
   why `internal/trajectory` takes nine minutes — is clearly a follow-up.
3. **R2-MAJ-2 shape.** Should the terminal status line move to stderr (smallest change, keeps the
   human output where it is for non-`--json` callers) or become a `result` field inside the envelope
   (better for scripted callers, changes the envelope a second time in one cycle)? I have no strong
   position; both satisfy B.3.
4. **R2-MIN-1 scope.** Both next-action states are arguably pre-existing digest-vocabulary gaps that
   B is simply the first consumer to expose — the same argument the quorum resolved *in-cycle* for
   MAJ-3/MIN-2 in round 1. Consistency suggests in-cycle again, but it is the quorum's call.
5. **Closing-round mechanics.** Given that everything I file here is MAJOR-or-below and none of it
   touches a gating property, does the quorum want a second narrow fix-up cycle followed by a third
   review round, or to carry `R2-MIN-3`/`R2-MIN-4`/the NITs as recorded deviations and close on a
   fix-up limited to the three MAJORs? Phase 8 requires zero Agreed fixes to close either way, so
   this is a question about how the remainder is dispositioned, not about relaxing the rule.

---

**Artifact path:**
`parley-deck/ideas/meta-protocol-change-lean-organizer/review/round-02/claude-1.md`

**Severity counts:** CRITICAL **0** · MAJOR **3** · MINOR **4** · NIT **4** — 11 findings
(`R2-MAJ-1..3`, `R2-MIN-1..4`, `R2-NIT-1..4`). Round 01 was 1/8/7/3; all 19 of those are
dispositioned above, 18 resolved and 1 (MIN-2) partially resolved and re-filed.

**Unresolved acceptance criteria.** One FINAL acceptance element is unmet and one is unmet at the
command FINAL names:

- **C, R-2 / open item 14** — "map-derived byte count recorded against **both** the measured floor
  (≈ 42 KB with §2) and the ≤ 70,000 B guardrail in the same test run": the guardrail and the
  `--optimize` baseline are recorded; the floor is not (`R2-MAJ-1`). The **gating** half of C.3
  (named-omission-set absence) and the guardrail itself are met.
- **Cross-cutting** — "both worktrees' test suites green incl. the new tests": green with
  `-timeout 2400s`; **not** green for me with the plan's own `go test ./...` (`R2-MAJ-3`).

Everything else in the A–D acceptance table verified at `118b245`: A's preflight gate,
driver-level role ineligibility with escalate-not-fallback, run-plan byte identity and prompt
parity; B's exit map, digest columns, determinism, timeout default/ceiling and
never-rewrites property; C's never-cut floor (runtime **and** the `packet check` proof), retention
and omission sets, fallbacks, attestation, SKILL.md core and relocation, and the brief contract;
D's handoff records, ledger contract, attribution mechanism and kimi telemetry; and cross-cutting's
three identical copies, the staged core's three change sets with project zones preserved, and the
untouched owner publish gate.

**Ready for a zero-fix closing consensus?** **No — but close.** Three MAJOR and eight lesser findings
are open, and Phase 8 closes only on a Phase-7 consensus listing zero Agreed fixes. None is a
blocker: there is no CRITICAL, no gating property is broken, no scope crept, and the trajectory is
the converging shape §4's stopping judgment describes — 19 findings in, 18 resolved, the new ones
smaller and confined to recording, output formatting and CI configuration rather than to mechanism.
I judge one more narrow fix-up cycle sufficient. My round-01 ❌ BLOCK is **withdrawn**; my position
for the Phase-7 consensus is 🟡 ACCEPT-WITH-RESERVATIONS on a fix-up cycle, not on closure. **I am
not calling this implementation complete, and I have released, published, merged, tagged and
installed nothing.**

**Tested commits:** CLI `118b2453ec0f12be063148fa575c2ccd7f2e8233` (fix-up source
`64a622ce48660725a59303befbd351962864c8e8`; round-1 record `3c97f44`; implementation `86d028b`;
base `b37f7ef`); skill `a820dc7fbec80855b5112b97735c890eff1b6f83` (base `d1e57d5`); staged core
`~/.parley/staging/COOPERATION-2.13.0.md` sha256
`fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f` (109,772 B).

**Review worktrees retained for further verification:**
`worktrees/lean-organizer-review-claude-1` @ `118b245` and
`worktrees/lean-organizer-review-claude-1-skill` @ `a820dc7`, both clean.
