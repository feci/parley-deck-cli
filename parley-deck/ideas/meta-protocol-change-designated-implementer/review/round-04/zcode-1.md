---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
review-round: 4
date: 2026-09-25
reviewed-commit: 717f3debe3aabb0f7a702d689e8d10de1e04fa8c (cycle-3 source, vs reviewed 1bad2634380dd785009987a5c50427050efb87c3) + ee8849c9ba6f75748fb471fe9c910bf362300e8d (cycle-3 record — IMPLEMENTATION.md only, no Go content; on-disk file byte-identical, sha256 86fd4cb58a5b4f61fd94b46efd61ae488c2db13030f641aac463ae85150c44ba) + a624318dcda02c47ecaa859d987efee08dba3104 (skill, unchanged this cycle; skill worktree HEAD re-verified at the full SHA)
---

## Summary

Round-04 re-review of the Phase-8 fix-up cycle 3 delta per the signed cycle-3 plan's
verification-plan step 3. **Both agreed documentary fixes (AF-16, AF-17) and the standing
per-publication metadata bump are applied exactly as signed, on my own primary evidence** — the
exact-commit git diffs, my own flagged AND flagless phase-8 packet runs, the tree-wide `any run`
sweep at HEAD, and the signed focused checks re-run by me in an isolated local-disk checkout
detached at `717f3de` (build/vet/gofmt exit 0; the four named `internal/app` tests PASS exit 0,
including `TestMalformedTier2Gates`, whose comment AF-17 touched, and `TestUnsetPathIsByteIdentical`,
the default-UNSET pin). The delta is one test-comment line plus record prose; no behaviour change
exists to re-derive, and per the ratified proportional plan **no broad suite and no new tests are
owed absent a new concern — I raised none**. I filed no new finding: a scoped null, with scope
named honestly below. Confinement and every frozen boundary verified: `FINAL.md` untouched since
`e4640bf`, `00-prompt.md` and all peer/review artifacts untouched, the signed cycle-3 consensus
frozen, the product default still ships UNSET. This review does not close the idea, signs no
consensus, and waives no coverage; the zero-fix consensus and the LE-7 fresh-invocation goal-done
check remain separately owed after this round.

## Protocol context attestation (Phase-8 re-review)

`parley protocol packet --dir . --phase 8 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change
--audience participant --json`, run by me this session in the original CLI workspace (the
designated-implementer worktree), parley 1.49.1, exit 0:

- `context_mode: full`, `fallback_reason` ABSENT (top-level keys: `body_path, context_mode,
  index, packet_sha256, request, shadow, source, source_sha256`).
- `source_sha256 = packet_sha256 = b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`
  (`COOPERATION.md`, role source, github-pr, 115,166 B / 1,400 lines) — unchanged from every
  cycle-2, cycle-3-signoff and round-03 attestation: this cycle moved no protocol text, exactly
  as the plan requires.
- Shadow audit: phase-8 WITH both flags → **86,716 B / 40 included / 29 omitted
  (`fe0e4c046897…`)**; the identical command WITHOUT the two flags (also run by me this session,
  same source hash) → **86,701 B / 40-29 (`a0845614a15c…`)**. Phase 8 and §15 were read from the
  full body.
- `strict_gate` is not set in `00-prompt.md` (re-read: `facilitator: codex-1`, `participants:
  [claude-1, kimi-1, zcode-1]`, no `strict_gate` key) — the default Phase-8 close rule applies.
- Flags attested precisely: this launch used `--flag auto_implement --flag protocol_change`
  (both required for this auto-implementing §7 protocol-change idea), matching every prior
  cycle attestation; the flagless variant was run deliberately, as evidence, and is labelled as
  such wherever cited below.

## Method and provenance

All evidence below is my own PRIMARY unless another participant is named; organizer notes were
orientation only and no organizer statement is a verdict. Unique local-disk isolated checkout —
`/private/tmp/zcode1-r04-cli`, `git worktree --detach` at the exact full SHA `717f3de…fa8c`,
`git status --porcelain` empty, on local disk, never the shared mount — **removed after the run**
(`git worktree remove --force`, verified gone). Executed there: `go build ./...`, `go vet ./...`,
`gofmt -l internal/app/driver_designation_test.go` (all exit 0 / clean), and `go test
./internal/app/ -count=1 -run 'TestMalformedTier2Gates|TestTier2UnavailabilityGateAndExits|
TestPinDesignationConflictEscalates|TestUnsetPathIsByteIdentical'` → `ok 0.478s`, exit 0, all
four PASS. In the shared worktree (read-only git + packet commands): the exact-commit diffs cited
below, the two packet runs, the tree-wide sweep, and hash verifications. Nothing was committed,
tagged, merged, published or released by me; no global default was mutated; I edited nothing
outside this file. Prior evidence is reused with exact-commit labels only (see AC-17 standing).

## The cycle-3 delta, reassessed

### AF-16 — flag-sensitive shadow labels in BOTH record locations, supersession claim struck — APPLIED AS SIGNED, verified live

- **Location 1** (`IMPLEMENTATION.md` cycle-1 Phase-8 attestation bullet, `## Validation
  evidence`, now `:322-336`): the "did not reproduce … and is superseded by this measurement"
  clause is gone; 86,716 B is labelled the phase-8 shadow **WITH** `--flag auto_implement --flag
  protocol_change`; 86,701 B (`a0845614…f08efd`) is labelled the SAME phase-8 packet over the
  SAME source **WITHOUT** them, "reproduces on demand"; the rule is completed to "Read the phase
  **and the flags** with the figure." The blanket sentence "All shadow byte figures in this
  bullet are rendered WITH …" covers the other figures and is TRUE for each: 86,336 B (pre-edit
  phase-8) was rendered by the verbatim cycle-1 command quoted in the same bullet — which
  includes both flags — and 80,799 B (phase-7) is separately marked "also rendered with both
  flags".
- **Location 2** (`### Deviations from agreed fixes`, now `:613-620`): retitled "One plan figure
  differs by invocation flags, disclosed"; states AF-14(1)'s 86,701 B is the FLAGLESS rendering,
  reproduces on demand, and that the record carries 86,716 B because the flagged command is what
  the plan and both cycle-2 signoffs attested with. Disclosure kept; no supersession claim.
- **My own live verification (the decisive check):** flagged phase-8 → 86,716 B / `fe0e4c04…` /
  40-29; flagless phase-8 → **86,701 B / `a0845614…` / 40-29**; identical source hash on both.
  Both labels in the record are correct for their own invocation, and the plan's observable
  holds — a reader running the flagless command finds 86,701 documented, not contradicted. A
  tree grep of the record finds "did not reproduce / superseded" only at `:686`, inside the
  cycle-3 section's description of striking the clause — a quotation of what was removed, not an
  assertion. This also closes my round-03 indeterminacy correctly: with the flagless run
  performed, the flag pair is confirmed as the discriminator (I had already confirmed
  claude-1's claim in my cycle-3 signoff; the record now states it as settled by all three
  participants' runs, VC-3.2 — accurate).

### AF-17 — the third `any run` carrier qualified — APPLIED AS SIGNED

- `internal/app/driver_designation_test.go:448` now reads "Present-empty is an incomplete
  designation, on any run that reaches a role action." — exactly the AF-3/R-1 scope wording
  reused from the corrected `driver_impl.go:188`, not expanded. The diff `1bad263..717f3de` on
  `internal/` and `cmd/` is this single comment line and nothing else (one hunk, one file,
  inside `TestMalformedTier2Gates`); no assertion, fixture or behaviour change.
- **Tree-wide sweep at HEAD, run by me:** `grep -rn "any run" --include="*.go" .` returns
  exactly the four hits the plan named — `driver_designation_test.go:448` (now qualified),
  `driver_impl.go:188` (qualified), `internal/agents/naming.go:50` and
  `internal/budget/run_identity_inventory_test.go:351` (unrelated English). I additionally
  checked `graphify-out/` contains no `.*.go` files, so the observable is exact, not an artifact
  of my grep scope.

### Standing Phase-8 per-publication bump (AF-12's restored obligation) — APPLIED AS SIGNED

- Frontmatter at `ee8849c` (disk byte-identical to it): `status: fix-up-cycle-3`,
  `head-commit: 717f3de` (the cycle-3 source commit — PRIOR and knowable under the
  source-then-record order), `skill-commit: a624318`. Short-SHA form remains the file's
  established convention — the round-03 dismissal stands; no re-filing.
- `## Fix-up cycle 3` section: names its source commit `717f3de`; states its own record commit
  "follows on the same branch and is named at the next touch" — **no self-hash claim** (grep for
  `ee8849c` in the record: zero hits), exactly as signed.
- The cycle-2 section's `record-commit:` is filled in: `b850576` (PRIOR and knowable now).
- The cycle-3 section's own packet attestation is flag-labelled ("phase-8 shadow WITH both flags
  86,716 bytes … not used") and matches my own run; its evidence block reports build/vet/gofmt
  and the focused app run at exit 0 — every one of which I reproduced independently at the exact
  source commit (same commands, same outcomes, `ok 0.478s` vs the record's `0.524s` — timing
  only). The organizer's note that the exact resumed implementer session exited 0 is consistent
  with the record; my verdict rests on my own re-runs, not on that testimony.

### Deliberately untouched, verified untouched

The AF-14 summary inside `## Fix-up cycle 2` → "Fixes applied" (`:576-580`) is byte-consistent
with the `b850576` text (`:574-578` there) — outside every diff hunk of the record commit, per
the plan's "no third edit" boundary. Behaviour pins: `internal/protocol/implementer.go` and
`internal/app/driver_impl.go` are unchanged in this cycle (the source diff touches only the one
test-comment line), so the AF-15 regex and AF-13 comments stand exactly as round-03 verified
them; `TestUnsetPathIsByteIdentical` (at `driver_designation_test.go:125`, same file) is
untouched by the single-line `:448` change and PASSES in my run — the default-UNSET pin holds.

## Refutation attempts

- **R-A — "AF-16's labels could be as wrong as the clause they replaced."** Refuted live: I ran
  both invocations myself this session over the same live source — flagged 86,716 / flagless
  86,701, identical block split 40-29, identical source hash `b273af1e…f388`. Every figure–label
  pair in both locations matches a rendering I produced; nothing in the record is contradicted
  by any command a reader can run.
- **R-B — "the blanket 'all figures rendered WITH flags' sentence could over-claim the pre-edit
  86,336 figure."** Refuted: the cycle-1 attestation's own command, quoted verbatim inside the
  same bullet, includes both flags; the sentence is true for every figure it covers, and 86,701
  is explicitly excepted as the flagless rendering rather than silently lumped in.
- **R-C — "AF-17's comment edit could have changed test behaviour or left another unqualified
  carrier."** Refuted: the `1bad263..717f3de` code diff is one comment line in one file;
  `TestMalformedTier2Gates` passes at the exact commit; my tree-wide sweep (including a check
  that `graphify-out/` holds no Go files) returns only qualified or unrelated hits.
- **R-D — "the metadata bump could be cosmetic or self-referential."** Refuted: the record names
  only PRIOR commits (`717f3de`, `a624318`, `b850576`); its own hash appears nowhere in it; the
  source-then-record order held (source commit `717f3de` precedes record commit `ee8849c`, whose
  direct diff is `IMPLEMENTATION.md` only).
- **R-E — "something could have moved under cover of the organizer commits."** Refuted by
  per-commit confinement, derived not assumed: `git diff --name-only 1bad263..HEAD -- internal/
  cmd/` = `driver_designation_test.go` only; `d238238..HEAD` = the four reviewed files and
  nothing else; `717f3de..ee8849c` = `IMPLEMENTATION.md` only; `ee8849c..357d1b5` (organizer
  handoff) = `organizer-notes.md`/`organizer-usage.md`/`usage-ledger.jsonl` only; the organizer
  commits between `b850576` and `717f3de` touch no code; `FINAL.md` is untouched from `e4640bf`
  and `00-prompt.md` from `d238238`; `review/consensus.md` is frozen since `770eb9b` (all three
  cycle-3 signoffs preserved) and the round-03 artifacts frozen since `0232ede`; the skill
  worktree HEAD is still the full `a624318…3104`.
- **R-F — "the record could be quietly re-claiming prior suite evidence at the new hash."**
  Refuted by reading: the record keeps AC-17's full suite **exact-commit scoped to `d238238`**,
  the targeted behaviour evidence **scoped to `1bad263`**, and claims for this cycle nothing
  beyond the focused run — the honest exact-commit labels the plan mandates, matching my own
  labels in this file.

## Findings

**None — a scoped null.** Within the scope the signed plan assigned round-04 (AF-16's two
locations; AF-17's comment plus the tree-wide sweep; the cycle-3 metadata bump) plus the
confinement/frozen-boundary checks I re-derived, I found no defect of any severity: nothing
incorrect, nothing undisclosed, no scope creep, no suppressed reservation, no new concern. Two
below-finding observations, for traceability only: (1) the record's `:686` quotation of the
struck "did not reproduce / superseded" clause is the phrase's only surviving occurrence — it
describes the removal and is accurate; (2) the frontmatter `head-commit` continues the file's
7-char short-SHA convention (standing round-03 dismissal, unchanged).

## AC-17 whole-suite standing — not re-claimed, not re-run, per the ratified plan

Claude-1's independent full suite stands **exact-commit scoped to `d238238`** (durable log cited
in round-03); the targeted behaviour evidence (focused runs, M1=9 / M2=5, probe sets) stays
**scoped to `1bad263`**; neither is a test at `717f3de` and the record does not present either
as one. This cycle's delta is one comment line plus record prose — no behaviour surface exists
to exercise beyond the signed focused run, which I reproduced at the exact commit. No broad
suite is owed or claimed; the LE-7 fresh-invocation goal-done check remains separate and owed.

## Deferrals — concurred, unchanged

DF-1 (drafter-precheck eligibility divergence → FINAL register F6, a later idea), DF-2
(launch-surfacing → the NAMED, INACTIVE slug, still absent from the tree), DF-3 (durable tier-3
fall-through notices → organizer release/done report plus the owner's post-release configuration
step), DF-4 (Windows residual → owner-deferred to `windows-portability`; this delta adds no
platform surface). All four carried unchanged by the cycle-3 consensus; this delta widens or
closes none; I concur with all four exactly as carried.

## Limitations (named, not waived)

No broad `go test ./...` was run by me at any commit (ratified proportional plan; the signed
focused checks are my whole execution set, plus the two packet runs and read-only git). I
re-measured the phase-8 flagged and flagless renderings but not the phase-7 rows of claude-1's
VC-3.2 table this session — those rest on claude-1's and my cycle-3-signoff runs. `ee8849c`
contains no Go content, so no suite result extends to it; I verified it by reading and diffing,
which is how both AF-16 clauses were checked. The cycle-1/2 fixes were not re-derived from
scratch: confinement proves this delta cannot touch them, and their standing evidence is the
round-02/03 record plus the exact-commit suites named above. Blind spots unchanged from cycle 3:
skill npm-level gates, end-to-end `parley run` on a designated deck, live §9.0 ping behind
`designeeAvailable`, TUI/pipeline-block surfaces, Windows — none waived. I did not read
claude-1's round-04 artifact (not filed at my run time; I was instructed not to wait on the
peer) and this review is independent of it.

## Disposition concurrences and verdict

I concur with every disposition of the signed cycle-3 consensus as executed: VC-3.1 (my round-03
null was a scope divergence, not a contradiction) and VC-3.2 (settled by all three
participants' runs, now stated accurately in the record); AF-16 and AF-17 applied exactly as
signed; the standing Phase-8 bump as AF-12's restored obligation rather than a new finding; the
proportional verification plan as applied (no broad suite for this comment/prose-only delta);
both dismissed findings standing; DF-1…DF-4 carried unchanged; the unchanged close gates — a
zero-fix consensus, then the LE-7 goal-done check by a fresh invocation of an existing
non-implementer quorum member, then the owner/organizer-only release sequence. The
stopping-judgment trajectory (17 → 8 → 2 → 0 findings from this reviewer) supports closure.

**Verdict: PASS — zero new agreed fixes from zcode-1 for this delta.** From my side the cycle-3
output is ready for the zero-fix consensus step. This verdict is mine as reviewer; it signs no
consensus, closes no idea, and authorizes no release act.
