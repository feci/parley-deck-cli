---
agent: kimi-1
idea: meta-protocol-change-lean-organizer
review-round: 1
date: 2026-09-24
reviewed-commit: 3c97f44
---

## Summary

The implementation is substantively sound: I independently re-derived every FINAL A–D acceptance
row against the reviewed trees and failed to break any of them — the preflight gate, driver role
ineligibility, `wait` exit map/ceilings, deterministic digest, audience packet retention/omission
(gating R-2 holds at phases 1 and 8), slim skill core, computed brief, usage ingest (streaming,
idempotent, explicit attribution), kimi telemetry (live-probed), the three protocol copies, and
the staged 2.13.0 core (exactly the three change sets). Three MINOR findings (a FINAL-A.5 §9.0
protocol-text gap, stale recorded byte measurements, and acceptance-test shapes that deviate from
the named forms) and three NITs are filed for Phase 7. I do not call this complete; nothing here
is a release or publish action, and none was taken.

## Protocol context attestation (this review)

```json
{"context_mode": "full", "source_sha256": "3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086", "packet_sha256": "3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase6-deliberation-3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086.md"}
```

Cross-check (PRIMARY): `shasum -a 256 parley-deck/COOPERATION.md` in my review worktree at the
reviewed commit = `3ad8a7eb…26086` — the attested packet IS the live authority at the reviewed tree.

## Reviewed trees, isolation, and tooling provenance

- CLI repo, reviewed commit: `3c97f445336a70bd7bd57eb15a0da78144bcd255` (final record; implementation
  `86d028b5558a60858534ac72f25537e52ed395fc`, base `b37f7ef9dd0941b21c3ba146c91d1a118259cbed`).
  FINAL frozen at `120a9bf757ad58794b1c6dd52105b2a301b23bae` — `git diff 120a9bf..3c97f44 --
  parley-deck/ideas/meta-protocol-change-lean-organizer/FINAL.md` is empty (PRIMARY).
- Skill repo, reviewed commit: `0f513f6c0fc78523d11321e60868942a08079d1d` (base
  `d1e57d5c8a56ac40fc0c8b4ac035069d2c0b9283`).
- My isolated review worktrees (parley-worktrees discipline; preserved until reviews/fixups close):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-review-kimi-1`
  (branch `review/meta-protocol-change-lean-organizer/kimi-1-20260924` @ `3c97f44`) and
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-review-kimi-1-skill`
  (same branch name @ `0f513f6`). All fixtures were disposable under `/tmp/kimi1-fixtures/`.
- Binary: independently rebuilt `go build ./cmd/parley` (go1.27.1 darwin/arm64) →
  `/tmp/parley-review-kimi-1/parley`, sha256 `6022c400dd346fdc12f20b004466ef02ea4f7085169bf0a603c76f505ffaf32c`
  (reproducible across two of my builds). The implementer's binary reports sha256
  `8d3e4c1c…` — different bytes, but identical behavior: both render the phase-1 facilitator
  packet with `packet_sha256 a95974cc…` (PRIMARY, both binaries run). I treat the hash delta as
  build-environment, not source, because behavior is byte-identical on a deterministic render.
- Test suites re-run by me at the reviewed commits: CLI `go build ./... && go test ./... -count=1`
  → exit 0, all 31 test packages ok (log `/tmp/kimi1-cli-full-suite.log`); skill `npm ci && npm test`
  → 399 pass / 0 fail (log `/tmp/kimi1-skill-test.log`). Both PRIMARY.

Every claim below is PRIMARY (a check I executed, with command and relevant output) unless
explicitly tagged otherwise. I am not an owner of the implementer's claims; I re-derived them
rather than endorsing them.

## Refutation attempts

### A — pure organizer default for declared facilitator runs

1. **Preflight fail-closed.** Built my own deck (`/tmp/kimi1-fixtures/preflight-conflict`) with
   `facilitator: codex-1` ∈ `participants:`. `parley preflight --dir … --no-ping` → **exit 3**,
   gate `[facilitator-declaration]` naming both `facilitator: codex-1` and
   `facilitator_participates: true`. Retried with `--yes` → still exit 3 (not waivable). With the
   exception flag → exit 0. Absent field → exit 0. Facilitator declared but NOT in participants →
   exit 0 (no gate). PRIMARY. Break attempt failed — the gate holds in all four adversarial shapes.
2. **Driver role ineligibility.** Read `internal/app/driver_impl.go:45-99` and
   `driver_consensus.go:84-122` at the reviewed tree: `newDriverImplOps` filters the declared
   facilitator out of implementer/reviewers/drafter (goal-done checker = drafter), and
   Implement/OpenReviewRound/Fixup/GoalCheck return the escalate-not-fallback error when the
   eligible set empties. `resolveImplementer` (`driver_impl.go:104-135`) validates recorded
   metadata against the ELIGIBLE list, so a stale facilitator-implementer record is rejected on
   re-entry. Tests `facilitator_roles_test.go` assert all of this including
   `escalated, not fallen back`. Break attempt: find any role path that still admits the
   facilitator — none found (see K1-F3 for the test-shape gap).
3. **Prompt repair parity.** `internal/protocol/consensussections.go:17-33` is read by the
   drafting prompt (`driver_consensus.go` `buildConsensusDraftPrompt`), the scaffold
   (`consensus.go designDraftTemplate`), and the gate (`consensus.go` →
   `protocol.MissingConsensusSections`) — one constant, parity test exists. The old review-cycle
   headings are gone from the design prompt. PRIMARY (code read + suite green).
4. **A-witness for this run.** This idea's own `00-prompt.md` declares `facilitator: codex-1`;
   `IMPLEMENTATION.md` records `implementer: zcode-1` ≠ facilitator. PRIMARY (files read).
5. **Break attempt that SUCCEEDED (K1-F1):** FINAL A.5 (`FINAL.md:84-88`) requires a permissive
   sentence in §4 Phase 5, §4 Phase 6, **§9.0**, and the Quickstart facilitator row in all three
   copies. §9.0 (`parley-deck/COOPERATION.md:847-881`) contains no pure-organizer sentence — the
   added line lives in §9 checklist item 1 (`:885`) and is D.5's audience+brief line. Same gap in
   `internal/protocol/defaults/COOPERATION.md` and the skill's `references/COOPERATION.md`.

### B — `parley wait` + PhaseDigest

All against my own disposable deliberation deck with two participants:

1. Missing artifacts → **exit 3** after `--timeout 12s`; partial digest names
   `claude-1 (round artifact), kimi-1 (round artifact)`. PRIMARY.
2. Present-but-invalid round file → **exit 4 immediately**, verbatim validator reason
   `… is missing a non-empty "## Existing alternatives" section (§15.6a)`. PRIMARY — the
   validity-with-failing-check column works (the concealment class from this run's notes).
3. Two valid round files → **exit 0**, `boundary reached (round complete)`. PRIMARY.
4. Usage errors → **exit 1** for missing `--idea`, unknown `--for`, bad `--timeout`, unknown idea.
   PRIMARY.
5. Ceiling: `--timeout 31m` on deliberation → exit 1 naming `30m0s`; after switching the idea to
   `track: fast`, `--timeout 10m` → exit 1 naming `5m0s`, `--timeout 4m` accepted. PRIMARY.
6. Missing ≠ invalid for consensus: no `consensus.md` → waits, exit 3 (not 4). PRIMARY.
7. Unanswered `to-user` inbox note → exit 4 naming the note file. PRIMARY.
8. Determinism: two `--json` runs byte-identical (`cmp` clean); digest carries
   agent/path/filed/bytes/owner/valid/stance_flags/unparsed/fell_back. PRIMARY.
9. Timeout default/config: `internal/config/runtime.go:644` seeds `wait_ms = 1500000` (25 m < 30 m)
   in `[defaults.timeouts]` with an explicit "NOT a verified provider cache fact" comment;
   `configLayers` (`runtime.go:391-392`) gives per-deck override via `parley-deck/agents.toml`;
   `--timeout` overrides per call. Grep for any provider-TTL-as-fact assertion → none. PRIMARY.
10. Break attempts that failed: `wait` contains no write calls (code read; `wait.go` never calls
    os.WriteFile); next-action is a fixed enumeration (`phasedigest.go:91-103`); adverse rows
    route to `open the raw artifact to adjudicate`. Scoping caveat recorded as K1-F4.

### C — audience packet, slim skill, computed brief

1. **Retention/omission, phase 1 (the guardrail's ratified scope).** Rendered
   `--audience facilitator --phase 1 --track deliberation --transport github-pr` on the reviewed
   deck: retention set 7/7 heading-verbatim (Quickstart, §2, §4, §5, §9, §11.B, §15); named
   omission set 9/9 absent (§1, §3, §8, §10, §12, §13, Appendix A, §11.A, §11.C) — the R-2 gating
   check. Body **58,941 B ≤ 70,000 B** (test `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`
   re-run by me; written file 59,024 B incl. 83 B file overhead). Omission index lists every
   omitted block with line span, content hash, classification, trigger. `source_sha256` is over the
   FULL authority (`3ad8a7eb…`). `--optimize` baseline measured by me in the same test run:
   **65,485 B**. PRIMARY.
2. **Phase 8.** Same render at `--phase 8`: omission set still 9/9 absent, retention 7/7, body
   **72,431 B** (test `TestLiveDeckFacilitatorAcrossPhases`; written file 72,514 B). The overage
   decomposes (my block-level diff phase 1 vs phase 8) into exactly the never-cut floor's phase-8
   pinning: `### Phase 5/6/7/8`, `#### Strict review gate`, `#### Stopping judgment`,
   `#### Review briefs and dispositions` — no omission-set block is retained. PRIMARY. See the
   disposition evaluation below.
3. **Never-cut floor.** `applicability.go:280-281` (`kernel={1,2,3,5,6,7,8}`,
   `fullFifteen={1,2,3,6,7}`) and `:283-310` — §15 unomittable at deliberation phases by
   construction; negative tests (`audience_test.go` hostile-map §15, `packet check` failing an
   audience that names an always-never-cut block) pass in my suite run. `parley protocol packet
   check --dir .` → `ok` across the audience overlay. PRIMARY.
4. **Fallbacks.** `--audience banana` → `context_mode: full` + `unknown-audience:banana`;
   `facilitator_participates: true` idea → `full` + `facilitator-participates`. Attestation gains
   additive `audience` / `audience_fallback_reason`, never `role` (`packet.go` Attestation struct,
   code read). PRIMARY.
5. **Slim skill.** `SKILL.md` = **17,038 B ≤ 20,000 B** (`wc -c`); names all six driver commands;
   links every `references/` file; description and Core Rule verbatim. My INDEPENDENT relocation
   check (not the implementer's hardcoded map): extracted all 29 `##` and 6 `###` headings from
   v2.12.1 `SKILL.md` (`git show d1e57d5:…`) and verified each survives in core or exactly one
   reference — 0 dropped, 0 duplicated. PRIMARY.
6. **Brief.** On my fixture deck: 1,267 B, byte-identical ×2, zero writes under `parley-deck/`
   (`find -newer` empty); on the live deck: 1,917 B with `context_mode: packet audience:
   facilitator`. Cap 8,192 B refuses rather than truncates (`organizer.go:119-122`); read-only-deck
   test uses chmod 0o555 + mtime snapshot. Break attempt (make it store a sidecar) failed. PRIMARY.
7. Break attempt that SUCCEEDED (K1-F2): the recorded byte vector in `IMPLEMENTATION.md` does not
   match the reviewed tree — see Findings.

### D — handoff records, usage ingest, kimi telemetry

1. **Handoff.** `commitCursor` → `writePhaseHandoff` (`driver.go` diff) writes
   `runs/<run-id>/handoff-phase-<phase>.md` atomically, embeds the schema doc ("recomputed view is
   authoritative"), emits `driver.phase_handoff`; a failed write is evented, never fatal. Unit
   tests (`phasehandoff_test.go`) verify schema-valid fields at a transition + round-trip. The two
   live `runs/` dirs in the implementer's worktree predate the feature (events.jsonl + run.json
   only) — no records expected there. PRIMARY (code read + suite).
2. **Ingest, codex.** My 3-event rollout fixture → one ledger row, LAST cumulative event verbatim
   (six fields: 7000/600/120/1100/60/9000), `parser_id codex-rollout/v1`, method header stating
   explicit-args + run-window attribution and no slug scanning; re-ingest → `idempotent no-op`;
   ledger still 1 row. PRIMARY.
3. **Ingest, claude.** 2-message fixture → summed 400/60/40/200/0/600 with the documented mapping.
   PRIMARY.
4. **No count flag.** `--tokens 9000` → `flag provided but not defined` (rejected). The flag set
   is dir/agent/source/path/idea/phase only (code read). PRIMARY.
5. **Attribution.** No `runs/` → `attribution=ambiguous` ("no run-record window available");
   after planting a `run.json` window covering the last event → `attributed`. PRIMARY.
6. **Streaming.** My generated 239,075,840 B rollout (with decoy lines containing the string
   `token_count` in user text) → peak RSS ~17 MB (`/usr/bin/time -l`: 17,154,048 B max resident),
   23 real events counted, decoys filtered by the JSON type check, stdout 121 B ≤ 1 KB. PRIMARY.
7. **kimi telemetry (live probe).** `kimi --output-format stream-json -m kimi-code/k3 -p "Reply
   with exactly: PONG"` on 0.42.0 → exit 0; stdout = `system.version` + `assistant "PONG"` +
   `session.resume_hint` envelopes, **zero `usage.record` lines** (honest `coverage: none` path);
   the on-disk `wire.jsonl` of that session carries 2 `usage.record` lines with the exact
   camelCase buckets (`inputOther`/`output`/`inputCacheRead`/`inputCacheCreation`) the new parser
   reads (`telemetry/usage.go` kimi case; fixture `testdata/kimi-stream.txt` matches reality).
   Structured argv is adapter-supported alongside `-p` — the D.4 premise, verified live. PRIMARY.
8. **kimi readiness repair (evidence class).** `parley preflight --dir <fixture>` (ping on) now
   reports `kimi-1 … yes` — the stream-json unwrap (`runner/kimistream.go`) repairs the exact
   parser rejection recorded in `organizer-notes.md`. PRIMARY.
9. Break attempts that failed: decoy `token_count` text lines cannot poison the codex parser
   (type check); a `{"role":"meta","type":"usage.record"}` wrapped record parses (envelope-tolerant);
   `UnwrapKimiStreamJSON` returns input unchanged when no assistant line exists (never a guess).

### Cross-cutting (three copies, staged core, release discipline)

1. **Three copies.** All five hunk fragments (Quickstart row, Phase 5, Phase 6, §9 item 1
   audience+brief, §11 wait advisory) present in deck view, `internal/protocol/defaults/`, and
   skill `references/` — and the K1-F1 §9.0 gap is consistent across all three. CLI drift test
   (`internal/protocol` package) green in my suite run; skill hunk assertions green. The staged
   2.13.0 core is byte-identical (sha256 `a8d3457a…dda5f4`) to the skill's bundled
   `references/COOPERATION.md` — the template-form copy. PRIMARY.
2. **Staged core `~/.parley/staging/COOPERATION-2.13.0.md` (109,507 B).** Verified independently:
   (a) placeholder header (`<workspace-name>`, `<transport-choice>`, `<YYYY-MM-DD>`) matching the
   2.10.0 template form; (b) stub §2 — empty roster view table AND empty host-handle table;
   (c) `diff` core 2.10.0 → staged = 10 hunks, each classifying into exactly one of the three
   change sets (3× 2.11.0 §15.6/§15.7; 3× 1.48.0 LE-bullet/goal-check/§9-item-1; 4× this idea
   incl. the shared §9-item-1 hunk) and nothing else; (d) 2.11.0 staged file's 14 added lines all
   present (faithful carry); (e) 1.48.0 fragments present in the released base deck and not the
   idea's; (f) changed sections (§15.6, §4 Phase 5, §9 item 1, goal-check paragraph) BYTE-equal
   the reviewed deck view. PRIMARY.
3. **Escalation note.** `inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_core-publish.md`
   carries the exact attended command `parley protocol publish --version 2.13.0 --from
   ~/.parley/staging/COOPERATION-2.13.0.md`. It was NOT run (no publish/merge/TTY action by me or,
   per the record, by anyone). PRIMARY (file read; `~/.parley/protocol/core/` still holds only
   2.10.0).
4. **Release discipline.** `VERSION` = 1.48.0, CLI `CHANGELOG.md` untouched, skill
   `package.json` = 2.12.1 — no premature bump; `parley-addon.json` regenerated for the new
   payload. Windows CI leg is new (`.github/workflows/tests.yml`, ubuntu+windows+macos,
   `go build` + `go test ./...` — covers `wait`/`usage`; this repo had no `.github/` before).
   `protocol-changelog.md` carries the idea entry. PRIMARY.

### Disposition evaluation (phase-8 facilitator body above the 70,000 B guardrail)

The implementer shows 71,168 B for phase 8 (my measurement at the reviewed tree: **72,431 B**)
while the ratified hard guardrail names the phase-1 body (58,190 B reported; **58,941 B**
measured). I evaluated whether the rationale holds:

- The FINAL's acceptance table scopes the guardrail to "the phase-1 / deliberation / github-pr
  facilitator body" — phase 1 measures 58,941 B, within guardrail. PRIMARY (my render + test).
- Open item 2's discipline is "show the bytes to the quorum before changing the map or the
  ceiling" — the bytes are shown (deviation note), map and ceiling untouched, never-cut floor
  untouched. The phase-8 overage is entirely floor-pinned §4 subsections (my block diff), not a
  retained omission-set block — the R-2 gating property holds at phase 8 too. PRIMARY.
- Residual: the recorded numbers are stale by ~751–1,263 B (measured before the protocol-text
  hunks landed, then not re-measured at the record commit), and the "hard-asserted" clause is
  wrong (the test asserts only ≤ 70,000 B). That is K1-F2, a record-accuracy finding; it does not
  change the substance — with corrected numbers the phase-1 body is still within the guardrail and
  phase 8 is still (further) above it for the same structural reason.

**My position: I concur with the disposition's rationale** (guardrail scope is phase 1; phase-8
overage is the ratified floor working as designed; the escalation discipline was followed) —
while requiring the record correction in K1-F2 so the quorum adjudicates with accurate numbers.

## Findings

### [MINOR] K1-F1 — FINAL A.5's §9.0 permissive sentence was never added (all three copies)

FINAL A.5 (`parley-deck/ideas/meta-protocol-change-lean-organizer/FINAL.md:84-88`) requires "one
permissive sentence each in §4 Phase 5, §4 Phase 6, **§9.0**, and the Quickstart facilitator row —
in all three COOPERATION.md copies". At the reviewed trees, §4 Phase 5
(`parley-deck/COOPERATION.md:443`), §4 Phase 6 (`:512`), and the Quickstart row (`:34`) carry the
pure-organizer default, but §9.0 (`:847-881`) carries no such sentence — in none of the three
copies (same regions in `internal/protocol/defaults/COOPERATION.md` and skill
`skills/parley-deck/references/COOPERATION.md`). The line added at `:885` is §9 checklist item 1
and is D.5's audience+brief line, not A.5's §9.0 default sentence. The record is also internally
inconsistent: `meta/protocol-changelog.md:2-6` says "§9 checklist item 1" while the owner
escalation note (`inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_core-publish.md`,
change-set item 3) calls it "§9.0 audience view". Impact is low (the default is enforced by
tooling and stated in three other locations) but the staged core would publish the gap.
Suggested fix: add the one-sentence pure-organizer default to §9.0 in all three copies, restage
the core, and align the escalation note's label.

### [MINOR] K1-F2 — Recorded facilitator-body measurements are stale at the reviewed commit

`IMPLEMENTATION.md:98-105` records "phase-1 … measures **58,190 B** (hard-asserted in tests)" and
the phase vector 52,854/58,190/58,291/59,551/62,469/59,791/61,005/68,880/71,168 B;
`IMPLEMENTATION.md:243-244` records `--optimize` baseline **64,734 B**. Re-running
`TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` and `TestLiveDeckFacilitatorAcrossPhases` at
`3c97f44` gives 53,605/58,941/59,042/60,302/63,220/60,864/61,946/69,821/72,431 B and an
`--optimize` baseline of 65,485 B — a uniform ~751 B shift at phases 0–4 and ~941–1,263 B at
phases 5–8, consistent with measurement taken before the protocol-text hunks landed and never
re-run for the record commit. Two clauses are wrong: the numbers, and "hard-asserted" (the test
asserts only `len(c.Body) > 70_000`; `facilitator_packet_live_test.go:118-122`). The substance is
unchanged (phase 1 within guardrail; phase 8 above it, further than recorded), so this is record
accuracy, not a guardrail failure. Suggested fix: re-run the two live tests at head and correct
the deviation note and validation evidence with the measured vector (58,941 B phase 1; 72,431 B
phase 8; 65,485 B optimize baseline).

### [MINOR] K1-F3 — Criterion A's named acceptance-test shapes are not the implemented shapes

The FINAL A acceptance row (`FINAL.md:434`) names "fixture auto-drive run launches the declared
facilitator for no drafter/implementer/reviewer/goal-done role (**event-log assertion**)" and
"absent-field deck → **byte-identical run plan** (regression test)". What exists:
`facilitator_roles_test.go` unit-tests `newDriverImplOps` fields and escalate-not-fallback at the
ops boundary (no driver run, no event-log assertion — the only event-log assertions added are
D.1's `driver.phase_handoff` in `phasehandoff_test.go`); and
`TestAbsentFacilitatorFieldKeepsV148RoleSelection` asserts role-selection identity, while
`internal/runplan/runplan.go:54` `Plan()` has no facilitator awareness and no byte-identity
regression test. The substance (facilitator never selected; escalate-not-fallback; absent-field
behavior preserved) IS covered by these unit tests plus my live preflight runs, so this is a
form/coverage deviation from the named acceptance shapes, not a behavioral gap. Suggested fix:
add a driver-level fixture test asserting role launches/escalation via the event log (and a
runplan byte-identity test or a recorded deviation explaining why the ops-boundary tests are the
refined equivalents).

### [NIT] K1-F4 — `wait` exit-4 fires on pre-existing escalations/driver errors, not only new ones

`internal/app/wait.go:250-266` (`blockingEscalation`) exits 4 on ANY `*-to-user_*.md` present in
`inbox/` — any idea, any age, including one that predates the wait — and
`driverErrorEvent` (`wait.go:281-299`) fires on any historical `driver.error` in the idea's
latest run, not only events arriving during the wait. FINAL B.3's wording is "a **new**
unanswered `to-user` escalation … arrives"; the acceptance-table row tolerates the broader read,
and fail-loud is the safe direction, but a stale unanswered note from another idea makes every
`wait` on the deck exit 4 forever. Suggested fix: scope to the waited idea (and/or to events
newer than wait start), or document the deliberate any-unanswered-escalation semantics in the
skill's wait section.

### [NIT] K1-F5 — Blind-spot (i): the phases-5/8 §15.5/§15.6 loss was measured but never analyzed

FINAL's context section (`FINAL.md:414-416`) asks implementation to "check what a facilitator view
loses there and whether it matters". At the reviewed tree the floor (`applicability.go:280-281,
308-309`: `fullFifteen={1,2,3,6,7}`) drops `### 15.5`/`### 15.6` from the phase-5 and phase-8
facilitator bodies (my phase-1-vs-8 block diff confirms), and `IMPLEMENTATION.md` records only the
phase byte counts — the "whether it matters" line is absent. Practical impact is nil (those
drafter duties are active at phases where they ARE pinned: 3, 6, 7), so this is a record
completeness nit. Suggested fix: one line in `IMPLEMENTATION.md` recording the check and its
conclusion.

### [NIT] K1-F6 — Cosmetic code issues in the new files

- `internal/driver/phasehandoff.go:58`: the `RunID` derivation
  (`filepath.Base(filepath.Dir(filepath.Clean(ideaDir + "/..")))`) yields `"parley-deck"` and is
  dead code — the only caller (`driver.go writePhaseHandoff`) always overwrites it. Remove the
  derivation (or the field from the builder signature).
- `internal/app/wait.go:371-372`: `var _ = strconv.Itoa` unused-import anchor; drop the import.
- `internal/app/usage_ingest.go:80-110` (`streamLines`): an oversized (>16 MiB) line is dropped,
  but the NEXT chunks begin a fresh accumulator, so the tail of an oversized line can be yielded
  as a partial "line". In practice this errors loudly (codex) or is skipped (claude unmarshal
  failure); a comment noting the partial-tail behavior would prevent a future parser from trusting
  it.

## Open questions

1. For the quorum on K1-F1: was §9.0 in FINAL A.5 meant literally (the readiness-check
   subsection), or was "§9" intended and item 1 accepted as satisfying it? If literal, the fix is
   three one-line additions plus a restage; if the latter, the FINAL record should note the
   interpretation before publish.
2. For the quorum on the disposition: with the corrected phase-8 figure (72,431 B), does anyone
   want the guardrail re-scoped or the phase-8 view slimmed before release? My position is no
   (floor behavior, opt-in view), but the corrected number is what the quorum should rule on.
3. K1-F4: is any-unanswered-escalation-blocks-all-waits the intended semantics, or should the
   exit-4 escalation path be idea-scoped? Both are defensible; the skill text currently implies
   the narrower read.
4. Binary provenance: my rebuild is reproducible but differs in sha256 from the implementer's
   binary while behaving byte-identically on deterministic renders. If the release wants a
   reproducible-build story (documented flags), that is release-step work, not Phase-6 scope —
   flagging it so it is not lost.

## Review record

- Canonical review file:
  `parley-deck/ideas/meta-protocol-change-lean-organizer/review/round-01/kimi-1.md` (in the
  original deck worktree `…/worktrees/lean-organizer`).
- Severity counts: **0 CRITICAL / 0 MAJOR / 3 MINOR / 3 NIT** (K1-F1…K1-F3; K1-F4…K1-F6).
- Verdict: all FINAL A–D acceptance criteria independently re-verified with real entrypoints and
  adversarial fixtures at the reviewed commits; the phase-8 disposition rationale holds (records
  corrected via K1-F2). Three MINOR and three NIT findings stand for Phase 7 adjudication. I do
  NOT call the implementation complete and I do not release anything; no publish/merge/TTY action
  was taken or is implied by this review.
- Tested commits: CLI `3c97f445336a70bd7bd57eb15a0da78144bcd255` (implementation
  `86d028b5558a60858534ac72f25537e52ed395fc`, base `b37f7ef9dd0941b21c3ba146c91d1a118259cbed`,
  FINAL frozen `120a9bf757ad58794b1c6dd52105b2a301b23bae`); skill
  `0f513f6c0fc78523d11321e60868942a08079d1d` (base `d1e57d5c8a56ac40fc0c8b4ac035069d2c0b9283`);
  staged core `~/.parley/staging/COOPERATION-2.13.0.md` sha256
  `a8d3457a48e32634f68c8eb1a3f4da9a171ba9ad2e0a2f4ba283795324dda5f4` (109,507 B).
- Isolated review worktrees (preserved for reuse/cleanup):
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-review-kimi-1` and
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer-review-kimi-1-skill`
  (branches `review/meta-protocol-change-lean-organizer/kimi-1-20260924`).
