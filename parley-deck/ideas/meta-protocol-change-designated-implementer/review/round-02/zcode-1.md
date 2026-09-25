---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
review-round: 2
date: 2026-09-25
reviewed-commit: d2382388415f8dc120a71ddbf6077d395dd18b8e (CLI source, baseline e4640bf; canonical record at e4d868a is IMPLEMENTATION-only) + a624318dcda02c47ecaa859d987efee08dba3104 (skill, baseline 8161e5e)
---

## Summary

The fix-up cycle-1 implementation is a faithful application of the signed plan — all eleven agreed fixes as amended by kimi-1's R-1/R-2/R-3 ratification are present and behave as signed, on my own primary evidence: isolated local-disk checkouts pinned to the exact SHAs, a full re-derivation of the three-copy/changelog/skill text state, and 50 adversarial parser cases of my own (calendar validity, negation placement/casing, subject-identity collisions, separator and shape boundaries, reassignment pair identity) which the AF-1-hardened parsers survive on every boundary that matters. The AF-2 re-entry protection now escalates on designation deletion and global-default clearing, machine-reads its own named reassignment exit, and confines corrupt-store fail-closure to genuinely designation-sourced dispatches exactly as R-2's narrower owner-compatible form requires, with the residual evidence-destruction limitation honestly pinned by test. Findings below are one MINOR (a protocol-required frontmatter bump not performed — traceability only, must be corrected before close) and three NITs; nothing CRITICAL or MAJOR. The AC-17 whole-suite clause remains independently UNVERIFIED by me per the cycle's serial-allocation: this review does not close the idea and drafts/signs nothing.

## Protocol context attestation (Phase-8 re-review launch)

`parley protocol packet --dir . --phase 8 --track deliberation --idea
meta-protocol-change-designated-implementer --flag protocol_change --json`, run by me in the CLI
source worktree against the live authority:

```json
{"context_mode": "full", "source_sha256": "b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388", "packet_sha256": "b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388", "source": {"path": "parley-deck/COOPERATION.md", "role": "source", "transport": "github-pr", "bytes": 115166}}
```

- `fallback_reason` ABSENT (top-level keys enumerated: `body_path, context_mode, index,
  packet_sha256, request, source, source_sha256`).
- The hash differs from the fix-up record's Phase-8 attestation (`c7492182…568f`, 114,771 bytes)
  for the correct reason, verified: that packet was rendered BEFORE the AF-3/AF-9 protocol-text
  edits; mine is rendered after them (115,166 bytes = the post-fix-up live authority). Not a
  discrepancy. Phases 4.0/4 (Phases 5–8), §15, §0/§9.0 were read from the full body.

**Method.** Unique local-disk isolated checkouts — `/private/tmp/zcode1-r02-cli` detached at
`d2382388415f8dc120a71ddbf6077d395dd18b8e`, `/private/tmp/zcode1-r02-skill` detached at
`a624318dcda02c47ecaa859d987efee08dba3104` — created by fresh clone, never the shared worktrees;
no peer fixtures or run caches used. My temporary adversarial test file lived only in my checkout
and was deleted after the run (tree clean, package re-verified green without it). Per the cycle's
review allocation I ran NO full suite and NO full `internal/app` suite: `go build ./...` OK,
`go vet` OK on all four changed packages, `go test` on `internal/protocol`, `internal/config`,
`internal/consensus` complete (all ok), and focused `-run` selections in `internal/app` covering
the entire designation/reentry/notice test inventory (19 + 7 + 2 named tests, all PASS). No test
processes left running; exit statuses preserved as quoted. All evidence below is my own PRIMARY
unless another participant is named.

## Refutation attempts

### AF-1 — confirmation-record parsers (identity / calendar / negation boundaries)

Attempted to break `ImplementerWaived` and `ParseImplementerReassignment`
(`internal/protocol/implementer.go`) with 50 cases of my own, run in my checkout:

- **Calendar validity (VC-D):** `2026-02-31`, `2025-02-29` (non-leap), `2026-00-10`, `2026-13-01`,
  `2026-01-32`, `0000-00-00` all REJECTED; `2024-02-29` (valid leap) ACCEPTED; non-padded
  `2026-9-25` rejected by the marker regex before `time.Parse`. The two-layer strictness
  (`(?i)^confirmed \d{4}-\d{2}-\d{2}$` then `time.Parse("2006-01-02")`) holds.
- **Negation:** `NOT confirmed`, `not yet confirmed`, `UNCONFIRMED` (any casing) in ANY earlier
  segment reject the record; negation in the marker position is impossible by regex.
- **Subject identity:** `kimi-10` vs `kimi-1` (prefix), `my-kimi-1` (suffix), `waived for kimi-1`
  (extra words) all REJECTED — exact trimmed equality enforced; `zz-impl-2`/`zz-impl` collision
  (claude-1 MAJOR-3) closed. Reassignment pair identity: reversed pair rejected, prefix
  collisions on either side rejected, `A to B`/`A → B`/`A -> B` separators accepted with exact
  trims, same-id `X to X` accepted (the documented source-only exit).
- **Shape:** no-reason (2 segments), empty reason, empty middle segment, marker-not-last,
  ASCII-hyphen separators (wrong dash), double-space marker all REJECTED; multi-segment reason,
  uppercase `CONFIRMED`, trailing-space marker (trimmed), quoted whole value ACCEPTED per the
  signed text.

Result: **not broken** on any authorizing boundary. Two tolerances documented as NIT-3 below.
Shipped `TestConfirmationRecordsAreFailClosed` + `TestDesignationRecordsRequireConfirmation` and
the existing T-4/T-5 fixtures (`zz-impl — offline — confirmed 2026-09-25`, confirmed
reassignment) PASS unmodified — the recorded shape itself is unchanged.

### AF-2 — re-entry hole (deletion / none / corrupt store / owner boundary)

Code trace (PRIMARY) + shipped tests (run by me, all PASS):

- `Implement` prints the byte-identical `driver: implementing via %s ...` line BEFORE calling
  `checkImplementerReentry()` unconditionally (hoisted out of the `implDesignated` guard) —
  deletion of a designation after a recorded dispatch escalates (B1:
  `TestReentryEscalatesOnDesignationDeletion` PASS, including the negated-record non-exit, the
  confirmed-record exit, and idea scoping), as does clearing the global default (B2:
  `TestReentryEscalatesOnGlobalDefaultCleared` PASS).
- The comparison keys on the RECORDED source ∈ {designation, global-default} (VC-A body-governs);
  a `none` or fall-through dispatch with a still-loadable store still compares against a
  designation-sourced record and escalates without a confirmed record — I traced this path; it is
  the ratified behavior, not a hole.
- Corrupt store (`Store.Load()` non-not-exist error): escalates fail-closed ONLY when
  `implLive` (exactly `implSource ∈ {designation, global-default}` after the pin exit — R-2's
  narrower form); unset/`none`/fall-through decks proceed (`TestReentryCorruptStoreFailsClosedOnlyWhenLive`
  3-leg PASS). The owner's unchanged-default boundary therefore holds by construction; the
  always-on form is demonstrably absent (the record's mutation evidence says broadening breaks
  legs 2+3 — consistent with the code I read).
- Residual evidence-destruction limitation (delete designation AND corrupt store → proceeds) is
  real, disclosed in `## Deviations from FINAL.md` item 6, and honestly pinned by
  `TestReentryResidualDeletionPlusStoreDestructionProceeds` PASS.
- `TestUnsetPathIsByteIdentical` and `TestReentryComparisonEscalatesOnDesignationChange` PASS
  unmodified; `TestPinDesignationConflictEscalates` (R28/AF-1-strict re-pin only after the
  confirmed record) PASS.

### AF-3 / AF-9 / R-1 / R-3 — protocol text, changelog, residual sweep

- All three `COOPERATION.md` copies carry verbatim: "fire on **any run that reaches an
  implementer, review-round, goal-check or fix-up action** … blocks that run before it
  dispatches" + "A design-only run (`auto_implement` off) reaches none of these actions, so a
  defective designation stays latent until the idea is next run with `auto_implement` on"
  (R-3(i) latent-until wording — the withdrawn "surfaces at its first dispatching action" clause
  is gone); §4.0 template comment "(blocked at the first role action)"; rank 4 "today's chain —
  `FINAL.md`'s recorded `implementer:` / `drafted-by:`, else the first eligible participant (list
  order)"; §10 TL;DR item 6 "then today's chain — …".
- R-1 hunks: `parley-deck/meta/protocol-changelog.md:12-13` carry both corrected claims (verified
  in my checkout); the changelog diff `e4640bf..d238238` is a single hunk `@@ -1,3 +1,30 @@`
  confined to this idea's own `**Status: UNRELEASED.**` entry; §7 shape (date, `Idea:`,
  `Drafted by:`, `Summary:`) present.
- Residual sweep for `fire on any run` (unqualified), `FINAL-drafter fallback`,
  `blocks the launch`/`blocked at launch`, `surfaces the defect at its first dispatching action`
  across all five files: CLEAN (no hits).

### AC-1 — three-copy fidelity (re-derived at d238238/a624318)

`tail -n +160 parley-deck/COOPERATION.md`, `tail -n +153` of `internal/protocol/defaults/COOPERATION.md`
and of the skill's `references/COOPERATION.md` all hash to
`3621b7a637c5183bcf9120ff73e231e2a4503859fdec6cd6b84bf9713beef30e` (matches the fix-up record's
claimed hash exactly). Pairwise diffs yield ONLY the three pre-existing project-zone hunks
(`3c3`, `6,7c6`, `154,159d152` deck-vs-defaults; same plus `5,7c5,6` header shape deck-vs-skill —
the same three logical hunks). **Confirmed.**

### AF-4 / AF-5 / AF-6 — present/live split and the two notices

- `live()` is exactly the signed predicate; `kickoffDesignationChecks()` is gated on it;
  `TestNoneAndFallThroughsAreNotLive` and `TestKickoffModelDiversityUnderDesignation` PASS;
  T-10/AC-14 (`TestTwoParticipantDesignatedKickoffWarnsNotBlocks`) PASS unmodified.
- AF-5: pin-shadowed malformed tier-3 → one-line NOTICE appended to the resolved line, no gate
  (`TestPinShadowedMalformedDefaultNoticesNotGates` PASS); malformed tier-3 WITHOUT a pin still
  hard-fails (`TestMalformedTier3HardFails` PASS).
- AF-6: `globalDefaultImplementer` carries the `LoadDefaults` error; ops construction prints one
  `driver: WARNING` naming it (`TestConfigErrorSurfacesAsNoticeNotGate` PASS); the third caller
  `designatedImplementerPreference` (`driver_consensus.go:122`) deliberately discards it — per
  the amended plan ("a read error steers nothing"), correct attribution.

### Remaining ACs / test-requirement matrix re-run

T-1/T-2/T-3 (two REAL layers verified: central+deck `agents.toml`)/T-4–T-11, AC-5/6/7/8/9/10/11/12/13/14/15/18
— every named test exists and PASSes in my checkout (full inventory listed in Method). AC-4's
legacy tests are byte-unmodified (`git diff e4640bf d238238 -- internal/app/app_test.go
internal/consensus/roundgate_test.go` → empty). AC-16: `grep -rn "impl-claim" --include="*.go"`
→ nothing; the delta adds no production `00-prompt.md` writer (diff scan shows reads and protocol
text only). AC-19: `centralDefaultTemplate` emits the key COMMENTED OUT
(`internal/config/runtime.go:654`); deck-creation test green. AC-20: no pty/terminal allocation
anywhere in the source delta. AC-21: round-01 carries both non-implementer files with populated
`## Refutation attempts`; the fresh-invocation goal-done check remains owed separately (below).

### Frozen-artifact integrity

`d238238` touched exactly the eight in-scope source/protocol/test files; `e4d868a` touched only
`IMPLEMENTATION.md`; `git diff 5710534..d238238` over `FINAL.md`, `consensus.md` and `review/` is
empty — no frozen FINAL, signoff, or peer-review edit. R16's launch-time surfacing is correctly
NOT implemented: no code, no idea directory for
`meta-protocol-change-designation-launch-surfacing` (named, INACTIVE, exactly as ratified); the
Windows track is untouched (no Windows file in the delta — owner-deferred DF-4).

### Canonical record (e4d868a) metadata/evidence evaluation

Every load-bearing claim I spot-checked reproduces: the AC-1 fix-up hash matches my re-derivation;
the AC-2/AC-3 re-derivations match; the mutation-proof inventory is consistent with the code
shapes I read; the packet-attestation hash difference is explained above; test names cited all
exist and pass. Two metadata defects found — F-1 (MINOR) and F-2 (NIT) below. The record honestly
marks its own full-suite PASS as implementer testimony and keeps step 2 (non-implementer serial
run in `review/round-02/`) outstanding, and claims no completion.

## Findings

### [MINOR] IMPLEMENTATION.md frontmatter not bumped for fix-up cycle 1 (protocol Phase-8 letter)

PRIMARY: `parley-deck/COOPERATION.md` Phase 8 requires, at each fix-up completion, "update the
top-level frontmatter: bump `status:` to `fix-up-cycle-N`, update `head-commit:`". At the
canonical record `e4d868a` (and still at branch HEAD) the frontmatter reads `status: implemented`,
`head-commit: 0893989`, `skill-commit: bf7e049` — the PRE-fix-up values. The fix-up HEADs
(`d238238` / `a624318`) are recorded only in prose ("Exact commits"). The record discloses a
deferral rationale ("The Phase-8 fix-up close will bump `head-commit` again per protocol"), and
AF-7's own rationale — "a fresh agent resuming from FINAL + IMPLEMENTATION alone would check out
the wrong tree" — is exactly the hazard the un-bumped machine-readable field partially recreates:
a frontmatter-only reader checks out `0893989`, not the fix-up tree. Not a correctness defect in
the shipped mechanism and not blocking re-review; it MUST be corrected at the next
`IMPLEMENTATION.md` touch and certainly before `status: complete`. One-line fix plus the
skill-commit bump.

### [NIT] AF-7/AF-8 landed in the record commit, not the fix-up source commit

PRIMARY: at `d238238` the duplicate `## Validation evidence` heading is still present
(`grep -c '^## Validation evidence'` = 2) and `head-commit:` still reads `e4640bf`; both are
correct only from `e4d868a` (= 1; = `0893989`). The record's stated split ("this file's update is
the next commit — a file cannot name its own commit") discloses this, and the AF-8 evidence claim
("grep -c returns 1") is true where it stands. Recorded for traceability of which commit carries
what; no action beyond F-1's correction.

### [NIT] Two documented parser tolerances inside the signed AF-1 shape

PRIMARY (my adversarial tests): (a) multi-segment reasons are accepted when every segment is
non-empty (`kimi-1 — reason—detail — confirmed 2026-09-25` is valid) — the signed text requires
"a non-empty reason segment" without capping segment count; (b) a reason consisting solely of the
word "not" does not trip the negation regex (which matches only `not [yet] confirmed` /
`unconfirmed` phrases) — the record still requires the strict trailing confirmed-date marker, so
no authorization is obtainable that the strict shape denies. Both are within the ratified
specification; recording them so future hardening is deliberate.

### [NIT] `live()` is true for a pin-source dispatch with tier 3 set

PRIMARY: in the `DesignationAbsent`+pin branch a present global default sets `present: true` with
`source: pin`, so `live()` — which excludes only {none, ft-inapplicable, ft-unavailable} — is
true, and kickoff designation checks (including the second `agent.model_diversity` event) fire on
a pinned re-entry whose owner has a global default set. This is exactly the signed AF-4 predicate
and the R45 present-notion, so it is implemented-as-signed; noting the semantic corner so the
AF-11 "confined to genuinely designated runs" sentence is read with it.

## Open questions

1. **AC-17 whole-suite clause — independently UNVERIFIED by me, per allocation.** I ran no full
   suite and no full `internal/app` suite (the cycle's serial broad execution is allocated to
   claude-1; no `review/round-02/claude-1.md` artifact existed at my write time, so there is no
   later artifact to cite). The implementer's serial full-suite PASS at `d238238` remains
   testimony I neither corroborate nor contradict. My scoped/verification clauses: `go build
   ./...` OK, `go vet` OK (changed packages), `internal/protocol`/`internal/config`/`internal/consensus`
   complete suites ok, focused `internal/app` designation inventory ok, `gofmt -l` clean on all
   ten changed Go files. The organizer reconciles the retained independent full-suite gate before
   any closure; AC-17's full-suite clause may be cited as independently confirmed only after that
   run is green.
2. **Still owed after this review (unchanged from the signed plan):** the retained AC-17 serial
   non-implementer run; the Phase-7 zero-fix consensus that closes the cycle; the LE-7 goal-done
   check by a fresh non-implementer invocation (AC-21's second clause — claude-1's C-1
   clarification governs: fresh invocation of an existing quorum member, no fourth participant).
3. **Carried blind spots (disclosed in the consensus; unchanged, no coverage waived):** no live
   §9.0 ping behind `designeeAvailable` (read-only equivalence); TUI and pipeline-block dispatch
   surfaces (FINAL F3/F10); Windows (DF-4, owner-deferred); skill npm-level gates belong to the
   organizer's release preflight; no end-to-end `parley run` against a real designated deck —
   code-trace plus the test inventory is the accepted evidence basis, as in cycle 1.

**Verdict.** All eleven agreed fixes as amended (R-1/R-2/R-3) are confirmed applied and behaving
as signed, on my own primary evidence; every prior round-01 finding remains visibly
dispositioned — none disappeared silently. Findings: 1 MINOR (record frontmatter, must be fixed
before close) + 3 NITs. This review does not approve closure, does not draft or sign any
consensus, and leaves the AC-17 full-suite clause and the goal-done check explicitly outstanding.
