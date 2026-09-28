---
idea: windows-portability
status: in-progress
implementer: zcode-1
started: 2026-09-25
branch: windows-portability
head-commit: 85babfe at this invocation's start (2026-09-28T13:03Z clock-verified); prior checkpoints 356fbb8 (FINAL freeze) -> e9cf601 (claim) -> 255f1a5 (Stage 0 code) -> 9a96c2f/16824fd (stage 0 close) -> ff2ef9d (recon) -> ca5efef (fsacl unwired) -> cd82025 (ACL checkpoint) -> 3646a9f (fsacl wired; assessed 36174658770) -> d07e6cc (stage 1 close-out; assessed 36425527266 this invocation) -> 276b3e1/85babfe (checkpoint + resume) -> a93f71c (row-24 fix + restore fix + Stage 2 core; assessed 36429107801) -> 6236af6 (Stage 2 complete; assessed 36432547744: pipeline GREEN) -> e4a2af9 (Stage 3 core; assessed 36434775624: named-refusal transition verified) -> 32bacd2 (A2/B3/B4 gate pins; cycle 36435481011 in flight) -> 9d9c9cd (checkpoint docs) -> e22b0ff (AC-DUR-3) -> 4e239fe (Stage 4 opener: read-only trap) -> 44dee89 (docs) -> c6f3db1 (§D.6 lock retry + diagnostic) -> 08f769c (row-19 redesign) -> 17d5210 (delete-share readers) -> bc08cd7 (identity-read retry + review-scope evidence) -> 9b880d3 (Stage 5 core) -> 4e5bb24 (Stage 5 checks pin + Stage 6 ACP/AF_UNIX) -> a2d6af6 (CRLF pinning + consult engagement) -> 5e25d78 (origin-read retry) -> 2da24e1 (checks-path restore) -> 4ed454a (row 9) -> this commit (row-30 remediation + record corrections)
design-pr: n/a
implementation-pr: n/a (owner override: no development PRs; files canonical, direct integration)
---

# IMPLEMENTATION — windows-portability

## Claim and worktree mapping (AC-IMPL-1 — recorded before any code edit)

- **Implementer:** zcode-1 — single implementation owner per FINAL §O. Claimed
  2026-09-25 in this file, before any product edit in this phase.
- **Workspace/branch mapping:** this isolated worktree
  `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/windows-portability`,
  branch `windows-portability`, base commit `868825f20bd0988abde1a6c38b9efe405c5123d1`
  (owner-assigned; FINAL §Context). FINAL freeze commit at drafting time: `356fbb8`.
- **Git-write authorization (this invocation):** commit prefix exactly
  `[codex-1] windows-portability:` and push only `windows-portability` to origin
  (hosted CI); no development PR, no main integration, no tag/release.
- **Design document:** `FINAL.md` (frozen, signed; never edited here). Selected
  design = the **refusal branch** preserving `os.Root` containment; **no
  containment deviation is approved** and this implementation neither waits for
  nor presumes one (AC-DEV-1).
- **Owner direction mirrored (2026-09-28, zcode-1; do not reopen):** the owner
  answered `inbox/user-to-codex-1_windows-portability_containment-deviation.md`
  — containment PRESERVED, the path-based MoveFileEx branch REJECTED. The seven
  rooted sites (A2, A3, B3, B4, B5, C1, C2) and their dependent re-sync
  operations refuse before mutation with a clear user-visible refusal; the two
  unavailable Windows budget features are disclosed in the release notes; the
  security guarantee stays identical to macOS/Linux. This is FINAL §A/§B as
  already frozen — implementation proceeds on the refusal branch with no
  further deviation questions.

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet: `.parley-runtime/protocol-packets/full-phase5-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md`
(read in full this session, Phase 5 IMPLEMENTATION, zcode-1 as single implementation owner).

## Summary of work

Implementing FINAL.md strictly: the refusal branch (§A–§C) plus the settled
non-durability requirements (§D.1–§D.9), staged per §H Stage 0–7, with the
hosted-only Windows evidence envelope (§F), the no-suppression census (§E), and
the H1–H7 probe register mappings honored (§G). Status: **in progress** — see
Progress for the live stage state.

## Implementation plan / checklist (FINAL §H stages)

- [x] **Stage 0** — contract/smoke/ledger: Mkfifo build tag; `app_test.go:155/:157`
      comma-ok; Ubuntu `GOOS=windows` cross-compile guard (amd64+arm64, build+vet);
      reconciliation ledger emitted (below); one diagnostic hosted cycle budgeted.
      COMPLETE (closed 2026-09-25T18:16Z): both diagnostic cycles assessed (36170672078,
      36170742720); evidence package compiles+runs on Windows (W4 fixed, AC-BLD-1);
      ubuntu cross-compile guard proven by two green ubuntu legs; the hosted cycle
      exposed one residual unchecked assertion (`app_test.go:185`, see Progress
      correction + row 16) — fixed this invocation; `wait`/`usage` first Windows
      execution now unblocked only after that fix lands hosted (next cycle).
- [x] **Early probe bundle** (before/alongside Stage 0–1): H3 rename mechanics + NTFS,
      H2 SyncFile raw error, H4 ReadDir class + toolchain, H1 ACL parent dump,
      H6 git config --show-origin, H7 read-only rename-over. Mechanics/diagnostic only.
      CODE LANDED this invocation (`internal/winprobe/`, test-only, zero t.Skip,
      windows-tagged): hosted execution rides the next push; outcomes to be recorded
      per the FINAL §G mappings — no probe upgrades an inference.
- [x] **Stage 1** — snapshot privacy/ACL (§D.1) incl. `denyRead` helper.
      COMPLETE at the code level this invocation (2026-09-28T12:4x–13:0xZ;
      hosted verification of the final items rides the next push): fsacl
      module wired at all product sites (store create/verify, archive temp
      protect, read-back verify, restore destination); adversarial suite
      hosted-run assessed (36174658770 — 6/7 PASS; the one failure was a bug
      in the TEST's own assertion, fixed: the no-rewrite check now compares
      the full DACL state before/after, strictly stronger); restore-dir
      privacy settled and wired (decision log 2026-09-28); store fixture now
      models a product-created store; wired product-boundary refusal test
      added. Remaining for Stage 1 close: the next hosted cycle confirming
      fsacl 7/7 + the 17 residual W1-signature failures clearing; AC-PRIV-4
      stays design-level refusal documentation (§N: no FAT volume on
      runners — the FAT path propagates SetNamedSecurityInfo's error as an
      ErrNotPrivate refusal, by code inspection fsacl_windows.go:51-55).
- [x] **Stage 2** — gate-name encoding + raw-ID allowlist (§D.4), legacy fallback +
      shadow retirement. COMPLETE at the code level across invocations 5–6
      (2026-09-28): encodeEdgeID in GatePath, legacy read-fallback, shadow
      retirement, Manifest.Validate raw-ID unsafe refusal (blocks + idea slug),
      ValidNewBlockID strict allowlist at pipeline start, AND (invocation 6)
      the error-returning GatePath/BlockWorkspace backstops with all 11
      callers converted (launchBlockRound, SeedBlockPrompt, blockCompleteFunc,
      planFinalized, run-block/auto/start sites + tests). Hosted status:
      encoding rode 36429107801 — pipeline red there for three distinct
      reasons, two fixed in the working tree this invocation (LoadGate legacy
      fallback treated Windows' ERROR_INVALID_NAME on raw-'>' names as a hard
      error instead of not-found; my legacy-fallback test itself tried to
      CREATE a raw-'>' file, impossible on Windows — now platform-conditional,
      pinning the uncreatability premise) and one STANDING pre-existing
      failure (TestEmbeddedDefaultMatchesLiveDeck, row 29). Hosted green of
      the full Stage 2 rides the next push.
- [ ] **Stage 3** — directory durability: §B table dispositions BEFORE code;
      `fsutil.SyncDir` named-type contract = audit of record + rooted refusal emitter.
      CORE LANDED invocation 7 (2026-09-28T14:0x–14:2xZ): §B disposition
      table recorded above BEFORE code; fsutil.DirHandle + SyncDir (Unix real
      fsync incl. darwin F_FULLFSYNC, byte-identical; Windows named
      fail-closed refusal ErrDirEntryDurabilityUnsupported) + Windows SyncFile
      directory-handle guard; A3/B5/C1/C2 pre-mutation entry refusals
      (durability_refusal.go, Windows-confined); D2/D3/D4a/D4b routed through
      SyncDir (the emitter); D1 Windows no-op argued per §B; windows-tagged
      hosted pins for the contract + A3/B5/C1/C2 (nothing-published asserted).
      AC-DUR-3 landed invocation 8 (2026-09-28T14:2x–14:4xZ): A1/B1/B2
      Windows-confined dormant conversions wired (fsutil.PublishFileDurable
      det-stage O_EXCL + fsync + MoveFileEx WRITE_THROUGH no-REPLACE;
      fsutil.PublishDirDurable stage-Mkdir + WT-move + EEXIST-tolerant
      recheck; POSIX stubs fail closed as caller-bug guards); StageSuffix
      '.parley-staging' grammar pin (platform-neutral unit test: fixed A1/B1
      finals + B2's dot-free UUID grammar can never collide); windows-tagged
      mechanics tests (publish/anti-clobber/stale-stage-blocker/idempotency —
      hosted executes them even though the product sites are dormant).
      Hosted AC-DUR-1/2 confirmation came in 36434775624; AC-DUR-3 mechanics +
      grammar-pin hosted confirmation came in 36437548293 (all four new tests
      PASSED first Windows execution). **Stage 3 COMPLETE — code and hosted
      evidence**; acceptance review in Progress.
- [ ] **Stage 4** — file sharing (§D.6) + open-site diagnostic cycle (H7).
      OPENED invocation 8 (2026-09-28T14:3x–14:4xZ): the §D.6
      read-only-attribute trap landed in ReplaceSyncedFile (Windows-only: a
      replace target without 0200 carries FILE_ATTRIBUTE_READONLY and
      MoveFileEx(REPLACE) fails over it — the write bit is restored before
      the move to match POSIX rename semantics; Unix replace untouched);
      windows-tagged pin TestReplaceSyncedFileOverReadOnlyTarget rides the
      next push. §D.6 lock unit landed invocation 9 (2026-09-28T14:4x–14:5xZ, commit c6f3db1):
      openLockFile (Windows: bounded 250ms ERROR_SHARING_VIOLATION-only
      retry; the self-held/foreign-held classification is structural — every
      handle this process takes on a lock file shares read|write, so a
      sharing violation on this open is by construction a FOREIGN holder;
      ERROR_ACCESS_DENIED is a different class and never retried; exhaustion
      is loud, naming the budget; only the open is retried, so the
      identity/exclusion verification chain re-runs unchanged after
      recovery; POSIX passthrough) wired into BOTH acquirePinnedKernelLock
      opens (kernel lock :159 + probe :200). TestOpenLockFileSharingDiagnostic
      IS the §D.6 open-site diagnostic cycle (windows-tagged, hosted): a
      foreign exclusive holder via raw CreateFile(share=0) — logs the ACTUAL
      error class and both retry outcomes (transient holder → recovery;
      persistent holder → loud named-budget exhaustion, window bounds
      asserted). REMAINING Stage 4 units: row-19 adversarial rename-root
      test redesign (invariant analysis: the 'root' case renames the rooted
      directory itself from outside — impossible on Windows while the root
      handle is open, BY DESIGN; the redesign mutates inside the root and
      must preserve 'change during verification is observed'), the
      share-flag-correct open sweep for concurrent persist readers (the
      residual Access-denied — cycle_extension_test.go:128 'persist operator
      cycle grant' — is ATTRIBUTED: a concurrent persist's WT-replace over a
      file another goroutine holds open without FILE_SHARE_DELETE; §D.6 fix
      = delete-share readers via os.Root routing or a raw CreateFile helper
      where rooting cannot express the site), and close-before-rename
      discipline.
- [ ] **Stage 5** — P-A liveness (§D.2) + W-SHELL missing-`sh` refusal (§D.3).
      CORE LANDED invocation 10 (2026-09-28T15:0x–15:1xZ): truthful P-A
      `alive` in BOTH probes (procctl_windows.go and driver/proclive_windows.go
      — OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION) + GetExitCodeProcess;
      ERROR_INVALID_PARAMETER = provably dead; any other observation failure
      including ERROR_ACCESS_DENIED = unverifiable = ALIVE for refusal
      purposes, fail closed; STILL_ACTIVE(259) sentinel caveat documented in
      both; the stale x/sys comment retired — x/sys is direct per go.mod);
      supportsDurableKill stays false with the Windows-honest attribution
      refusal (P-B named follow-up); bootID stays "" — the four §D.2 callers
      (durablekill.go, verification.go, reviewsnapshot.go via procctl.Alive)
      get truthful liveness through the probe. W-SHELL (§D.3): LookPath("sh")
      ONCE BEFORE any work at both sites (driver_impl.go checks path;
      evidence RunCriterionControlled captured-verification path) with the
      named Git-for-Windows refusal; windows-tagged pin
      TestRunCriterionControlledMissingShRefusesPreWork (sh forced off PATH →
      named message + non-zero + no work) + WithShRunsNormally (present →
      normal path; missing sh on the hosted image would be a loud t.Fatal,
      never a skip). The checks-path pin landed invocation 11
      (TestRunChecksMissingShRefusesPreWork; the list-form `checks:` contract
      path routes through evidence.RunCriterion → RunCriterionControlled and
      inherits the same refusal — all three §D.3 surfaces pinned).
      **Hosted confirmation (36441834810 @9b880d3, re-scoped 15:35Z after an
      evidence-integrity correction): zero failures among the pins THAT
      COMMIT CARRIED** — the two evidence-package W-SHELL outcome tests and
      the lock diagnostic. CORRECTION (invocation 12): an earlier entry
      claimed the checks-path refusal also passed in 36441834810 — impossible:
      TestRunChecksMissingShRefusesPreWork first shipped in 4e5bb24 and the
      product checks-path code only in 2da24e1; a later test cannot execute
      in an earlier run. The checks-path pin's hosted confirmation came in
      36443926903 (@2da24e1, TestRunChecksMissingShRefusesPreWork PASSED) —
      Stage 5's hosted evidence is complete with exact commit↔run mapping. In 36442736966 (@4e5bb24) the test FAILED for real
      (empty refusal message — the commit genuinely lacked the product
      code); that was a REAL failure of the incomplete pushed commit, not a
      "false test failure" — wording corrected here per the organizer's
      requirement reminder. 13 red then = the standing Stage 6 families.
- [ ] **Stage 6** — fixture-portability sweep (§D.9) + ACP/AF_UNIX/CRLF (§D.7).
      OPENED invocation 11 (2026-09-28T15:1xZ): (1) ACP split per §D.7 a/b —
      the two gated drain tests now use a 1 KiB total payload
      (capacity-INDEPENDENT: Windows anonymous pipes buffer ~4 KiB, where the
      historical 16 KiB in-flight volume stalled the child before READY —
      row 17's mechanism); the copier still parks in its first Write so the
      drain-ordering invariant (Stop/Wait never reap before the drain
      completes) is unchanged; the capacity-dependent 16384-write/8192-ring
      volume stays in the UNGATED TestSpawnObservesAllStderrAndRealExit
      (hosted-green); spawn.go Stop/Wait ordering untouched. (2) AF_UNIX per
      §D.7 — all four t.Skipf fallbacks in tree_report_test.go are now hard
      failures (Getwd/Chdir setup, the unix-socket listen, the escaping
      symlink); if hosted AF_UNIX creation itself fails, that loud fact opens
      an individually reviewed exclusion row per §D.7 — never a silent skip.
      Hosted execution of both rides this push. NEXT Stage 6 units: the §D.7
      CRLF load-bearing layer (per-invocation `-c core.autocrlf=false -c
      git.eol=lf` on parley's own git invocations, precedent
      reviewsnapshot.go:138-139), HOME→USERPROFILE fixtures (row 9), denyRead
      sweep (rows 8/20), /repo literals (row 10), agent fixtures (rows
      16/18/23), row 22 newline fixture, row 27 POSIX-host refusal family.
- [ ] **Stage 7** — validation/release: unfiltered three-leg matrix green, census
      enforcement (§E), coverage statements (§F), release sequence (§O).

Files/areas (expected, stage-tagged): S0 `internal/evidence/tree_report_test.go` +
new mkfifo helper files, `internal/app/app_test.go`, `.github/workflows/tests.yml`;
S1 `internal/*/snapshot.go` + new `internal/fsacl` (naming TBD) + `denyRead` helper;
S2 `internal/pipeline/gate.go` (locator TBD at stage open); S3
`internal/runner/{trajectory_verify,verification,parent_recovery,reservation_recovery,unchanged}.go`
+ `internal/fsutil`; S4 lock/open/rename sites per §D.6; S5 `internal/*/proclive*`,
`durablekill.go`, `driver_impl.go`, `evidence/execute.go`; S6 fixture test files per
ledger; S7 workflow census.

Checks to run (local, macOS host — never claimed as Windows evidence): `go build ./...`,
`go vet ./...`, `GOOS=windows GOARCH={amd64,arm64} go build/vet ./...`,
targeted `go test` per touched package. Hosted: push `windows-portability` →
unfiltered three-leg `Tests` matrix; run IDs recorded below.

Review/risk notes: Unix paths must stay byte-identical at every §B row (AC-DUR-5);
every Windows-confined change behind `//go:build windows` / `!windows` splits with
shared cross-platform code reviewed once. No skip added anywhere without a ledger row.

## Deviations from FINAL.md

None. (Any unavoidable deviation will be logged here, not silently absorbed.)

## Notes for reviewers (claude-1, kimi-1 — Phase 6)

- **Stage 1 decisions to scrutinize:** (1) restore-destination privacy now routes
  through fsacl.EnsurePrivateStore (Decision Log 2026-09-28 — product-created dir,
  creation policy, no user-data rewrite; basis is the pre-existing tested invariant,
  not new policy); (2) `snapshotStoreFixture` fidelity change + the three-layer
  coverage of the AC-PRIV-5 refusal; (3) ledger row 24's deferred race fix
  (atomic CreateDirectory) — deferred deliberately, not overlooked.
- Refusal strings are product UX: each must name context, be actionable, and fire
  pre-mutation (AC-DUR-2). Suggested focus areas per stage are in the Progress notes.
- The census baseline (104/103 = 68+35+0 / 56/55) is pinned at AC-CENSUS-4; re-run
  the attestation commands at the reviewed commit.
- §F: hosted windows-latest x64 is the only Windows execution evidence; ARM64 is
  build-only; this workspace sits on a cross-machine shared volume whose Windows-side
  semantics are unvalidated (all local checks here are macOS-host evidence only).

## Progress

### Process-record correction (2026-09-25T18:16Z, disclosed per organizer observation)

The two original Progress entries below carried approximate timestamps "~18:20Z" and
"~18:35Z". Both are impossible: invocation 1 actually exited 0 at 2026-09-25T18:01:28Z,
and its real commits are timestamped e9cf601=17:58:23Z (claim) and 255f1a5=18:00:12Z
(Stage 0 code; committed with author/committer clock 20:00:12+02:00 = 18:00:12Z). The
"~18:20Z"/"~18:35Z" figures were projected estimates written as if observed, not
measurements; no test executions occurred at those times and no evidence was taken
from them. They are re-dated below to the actual commit timestamps. Rule adopted from
here on: Progress entries use actual UTC clock reads at write time or commit/run
timestamps only; estimates are labeled as such and never in the 18:20Z/18:35Z form.
[Second disclosed slip, same invocation: the first draft of the Stage 1 recon entry
below was stamped 18:22Z while the actual clock read 18:18:34Z — a projected write
time, caught on re-verification and re-dated to commit 0baff30's 18:18:00Z before
push. The rule is now applied with a post-write clock check.]

- (2026-09-25T17:55–17:58Z, corrected; commit e9cf601 at 17:58:23Z) Invocation 1
  started: read full phase-5 packet, 00-prompt.md, FINAL.md, source-context; wrote
  this claim before any code edit. Beginning Stage 0.
- (2026-09-25T17:58–18:00Z, corrected; commit 255f1a5 at 18:00:12Z) Stage 0 code
  complete: mkfifo build tag
  (`internal/evidence/mkfifo_{unix,windows}_test.go`, call site switched, `syscall` import
  dropped from `tree_report_test.go`); comma-ok assertions at `app_test.go` former :155/:157
  (fail-with-message, never panic-abort); Ubuntu-leg Windows cross-compile guard
  (amd64+arm64 build+vet) in `tests.yml`. Local checks green (see Validation evidence).
  Pushed for the budgeted diagnostic hosted cycle; green NOT assumed (§H Stage 0).
  Hosted run 36170672078 (push 255f1a5) in progress at invocation close — NEXT STEP:
  watch to completion (`gh run watch 36170672078`), pull per-leg logs, and refresh the
  reconciliation ledger from the real Windows denominator (expect: internal/evidence now
  compiles+runs; internal/app no longer aborts — wait/usage results should appear, with
  TestVersionAllUsesDirFlagForProjectStatus likely failing gracefully until the Stage 6
  W6 fixture port; W1/W2/W3 product families still red by design until Stages 1/2/4).
  [Assessment came in invocation 2, 2026-09-25T18:06–18:16Z — see the next entries:
  evidence now runs (W1 surface confirmed); internal/app still panicked, but at :185,
  a residual unchecked assertion; expectation otherwise met.]
- (2026-09-25T18:05–18:16Z) Invocation 2: re-read packet/prompt/FINAL; assessed both
  hosted cycles (36170672078 @255f1a5, 36170742720 @9c1db32 — same code, docs-only
  delta): macOS+ubuntu legs GREEN in both; Windows legs FAIL in both with IDENTICAL
  deterministic red sets: 14 red packages / 16 ok / 111 failing test functions
  (subtests included) in each — no flaky-only failures observed between the two runs.
  Stage 0 closed: the residual `app_test.go:185` panic (and :125/:129/:130/:189
  siblings) fixed with comma-ok (file audit: no unchecked map assertions remain);
  `wait`/`usage` still have NO first Windows execution — the :185 panic aborted
  internal/app before those files ran in both cycles; unblocked only by this fix.
  Early probe bundle landed as `internal/winprobe/` (test-only, windows-tagged,
  7 test functions, zero t.Skip; census delta +7 Windows-only at Stage 7).
  Windows leg expected red again next cycle by design (probes add diagnostics, not
  product fixes); W1 evidence-package privacy failures now confirmed hosted
  ("snapshot store must be a private real directory" x71 signature present).
- (2026-09-25T18:16Z) Stage 0 checkbox/register/head-commit bookkeeping brought into
  agreement with actual work: Stage 0 marked complete (fix + guard + ledger + two
  assessed diagnostic cycles), run register filled with final per-leg outcomes,
  ledger refreshed from the real denominator with new rows 16–18, frontmatter
  head-commit/branch corrected. Committed as 16824fd and pushed (hosted cycle
  36172430646, started 18:17:04Z — carries the probe bundle + app comma-ok fix;
  NOT yet assessed at this checkpoint).
- (2026-09-25T18:17–18:18Z; commit 0baff30 at 18:18:00Z) Stage 1 reconnaissance (recorded for resume; no product code
  written yet — starting a half-wired ACL change at the checkpoint boundary would
  leave an incoherent tree):
  - Guard site: `internal/trajectory/snapshot.go:150-165` `privateSnapshotDirectory`
    (os.MkdirAll 0700; Lstat; IsDir/symlink/Perm()&0077 → the 71x hosted signature
    at :163). Read-back sibling check at :544 (Perm()&0077 on the archive file).
  - Plan per §D.1: new `internal/fsacl` package — Unix path byte-identical to today
    (0700 MkdirAll + Perm check, no behavior change); Windows path = owner-only
    PROTECTED_DACL via ACLFromEntries/SetNamedSecurityInfo (exact mechanism the H1
    probe seeds), verification walks effective ACEs via GetNamedSecurityInfo and
    refuses any non-owner grant NAMING THE TRUSTEE (sentence-stable refusal with
    trustee appended); FAT/exFAT/no-ACL volumes refuse (SetNamedSecurityInfo error
    propagates as refusal, never weaken); pre-existing stores stay refuse-and-instruct
    (no in-place DACL rewrite of user data); share the `denyRead` test helper with the
    Stage 6 sweep (D.9). The trajectory tests currently x71 red on Windows become the
    hosted adversarial suite once the store passes own-creation.
  - Resume here: implement fsacl + wire privateSnapshotDirectory + :544 read-back,
    keeping Unix output byte-identical (AC-DUR-5 applies to §B rows; this is §D.1
    but same discipline), then assess run 36172430646 probe outcomes FIRST and
    record them per §G mappings before building on any probe result.
- (2026-09-25T18:19–18:21Z, commit follows) Stage 1 opened: `internal/fsacl`
  landed as a self-contained module, deliberately UNWIRED from
  trajectory/snapshot.go in this invocation (wiring + adversarial hosted tests
  are the next unit of work; §D.1's AC-PRIV-1..6 deserve a full window, not the
  tail of this one). Contents: ErrNotPrivate sentence-stable refusal; Unix
  Ensure/Verify byte-identical to the inline guard (0700 MkdirAll +
  IsDir/symlink/Perm checks); Windows owner-only protected DACL via
  ACLFromEntries/SetNamedSecurityInfo (the H1-probe-seeded mechanism) with
  create-then-requery self-check and trustee-naming walk (refuses non-owner
  grants, inherited ACEs, unprotected DACLs, missing DACLs; FAT/exFAT refuse
  via error propagation); shared DenyRead helper (chmod 0 Unix, deny-ACE
  Windows). Local: go test ./internal/fsacl/ ok (5 funcs incl. byte-identical
  refusal-string pin); GOOS=windows amd64 build+vet / arm64 build OK. On the
  windows leg the package tests double as first hostile execution of the DACL
  round-trip and DenyRead (symlink-creation assertion may surface runner
  privilege facts — if os.Symlink errors there, that is a hosted fact to
  record, not a suppression candidate).

- (2026-09-25T18:23–18:38Z, zcode-1; commit follows) Invocation 3. Re-read packet
  (attestation above), 00-prompt, living IMPLEMENTATION.md, FINAL §D.1/§G/AC-PRIV.
  While the two in-flight cycles ran: wired fsacl into the product (three sites,
  Windows-confined effects, Unix byte-identical), restructured Windows Ensure for
  AC-PRIV-5 refuse-and-instruct with no in-place rewrite, added reparse-point
  rejection, landed the AC-PRIV-1..6 windows adversarial suite, split the Unix
  perm-bit pin into fsacl_unix_test.go. Local checks: darwin build ./... OK,
  vet fsacl+trajectory OK, fsacl tests ok; GOOS=windows amd64 build+vet ./... and
  fsacl/trajectory OK, arm64 build OK (see Validation). Restore-file privacy
  (snapshot.go:630 OpenFile 0600 into user-chosen destinations; snapshot_test.go:124
  asserts restored perms) is identified but deliberately NOT wired this invocation:
  setting owner-only DACLs on restored trees removes inheritance users may want in
  shared destinations — a product-behavior decision to settle with reviewers as the
  last Stage 1 item, not a unilateral change inside a tight window.

- (2026-09-28, zcode-1; commit follows — start time is an estimate (~12:30Z),
  first verified clock read 12:57:48Z mid-invocation, the commit timestamp is
  authoritative). Invocation 4. Re-read packet
  (attestation above), 00-prompt, the owner's containment answer, living
  IMPLEMENTATION.md, FINAL §D.1/§I. Assessed hosted cycle **36174658770** (the
  wired-guard + AC-PRIV suite cycle, @3646a9f; per-leg facts in the run register),
  then finished Stage 1's deferred items:

  [PROCESS-RECORD CORRECTION 2026-09-28T13:28Z, disclosed per organizer process
  evidence: invocation 4 actually LAUNCHED at 2026-09-28T12:45Z (organizer
  process evidence), not "~12:30Z" as first estimated below. The "~12:30Z"
  figure was a projected estimate written as a launch time; the real first
  clock read remains 12:57:48Z. Entries below are re-anchored: the invocation-4
  working window was 12:45Z–13:01Z (checkpoint commit 276b3e1 at 13:01:01Z,
  resume commit 85babfe at 13:03:41Z, both clock-verified in git).]

  - **Assessment of 36174658770 (Windows leg):** macos+ubuntu GREEN. Windows:
    15 red / 17 ok packages = the SAME deterministic 14 + fsacl. The fsacl
    adversarial suite went **6/7 PASS** on real Windows: own-creation round-trip
    with independent requery (AC-PRIV-1), BUILTIN\Users and Administrators
    refusals naming the trustee (AC-PRIV-2/3), inherited-only refusal (AC-PRIV-3),
    symlink/non-dir rejection (AC-PRIV-6), file-policy round-trip — all
    hostile-execution PASS. The single failure, TestPreexistingStoreRefusedAndNotRewritten,
    was a bug in MY test's own final assertion: the scaffolding `grantTrustees`
    writes the pre-existing store's DACL WITH the PROTECTED flag, so
    `protected=true` was the before-state, not evidence of a rewrite; the product
    had refused, named the trustee, and left the DACL and marker file untouched
    (product code path fsacl_windows.go:44-46 returns before any SetNamedSecurityInfo
    — no rewrite is possible). Fix: the no-rewrite check now snapshots the full
    DACL state (protection + ACE list) before Ensure and requires it identical
    after — strictly stronger than the old protected-flag heuristic.
  - **Wired guard hosted-verified:** the W1 hosted signature dropped 81x → 17x,
    and the residual 17 all carry the NEW AC-PRIV-5 refuse-and-instruct text
    ("DACL not protected (inheriting) … move it aside and re-run"). All 17 are
    stores the TESTS pre-created via `snapshotStoreFixture` (bare t.TempDir() =
    inheriting DACL on Windows) — the product correctly refusing someone else's
    permissive store, which is exactly the §D.1 contract. Fixture fixed to model
    a product-created store (fresh subpath through fsacl.EnsurePrivateStore);
    the refusal case stays covered by a NEW wired product-boundary test
    (TestCaptureRefusesPermissivePreexistingStore) plus the fsacl suite.
  - **Restore-dir privacy settled and wired** (Decision Log 2026-09-28):
    RestoreSnapshot's destination is itself a product-created fresh dir
    (os.MkdirTemp), so §D.1 creation policy applies — no user data is rewritten.
    The prior invocation's shared-destination inheritance worry does not apply
    to the current shape: restore never writes into an existing user directory.
    snapshot_test.go's "restore directory is not private" assertion now goes
    through fsacl.VerifyPrivateStore (same Perm check on Unix; DACL walk on
    Windows) — same invariant, native expression per platform.
  - **New phenomenon recorded, not reacted to:** TestSnapshotConcurrentCapturePublishesOneExactArchive
    fails with two distinct unpublished results; root cause not established from
    the log (ledger row 24). Candidate mechanism noted there (create→set-DACL
    window in concurrent first-capture) with the atomic-CreateDirectory design
    as the next invocation's fix — deliberately NOT half-landed at this
    checkpoint boundary.
  - Remaining Stage 1 close-out: next hosted cycle must show fsacl 7/7 and the
    17 W1-residual failures clearing (tests then proceed to the already-ledgered
    Stage 3/4 families: dir-fsync Access-denied, rename/sharing).

- (2026-09-28T13:03–13:2xZ, zcode-1; commit follows — first verified clock read
  13:04:10Z, commit timestamp authoritative) Invocation 5. Re-read packet
  (attestation above), 00-prompt, living IMPLEMENTATION.md, FINAL §D.1/§D.4.
  Assessed hosted cycle **36425527266** (Windows leg completed 13:06:09Z —
  per-leg facts in the run register), then landed the row-24 fix, a
  hosted-confirmed restore-path bug fix, and the Stage 2 core:

  - **Assessment of 36425527266 (Windows leg @d07e6cc):** ubuntu GREEN, windows
    FAIL with 15 red / 16 ok = the SAME deterministic 14 + fsacl (no new red
    packages). W1 privacy signature 17x → **2x**, and both residuals are the
    SAME bug (row 28): `RestoreSnapshot`'s `os.MkdirTemp` destination is a
    pre-existing inheriting-DACL dir, which `EnsurePrivateStore` treats as
    foreign user data and refuses — the restore-dir wiring committed in d07e6cc
    cannot work on Windows as written (snapshot_test.go:147 and
    source_inventory_test.go:111, both inside RestoreSnapshot calls). Found by
    code inspection BEFORE reading the hosted log; hosted log then confirmed
    both sites verbatim. Fix: new `fsacl.ProtectPrivateStore` (creation policy
    for a product-created dir: set owner-only DACL + verify; no
    refuse-and-instruct) — Unix verify-only (MkdirTemp 0700 already satisfies,
    byte-identical), Windows SetNamedSecurityInfo; snapshot.go:751 switched to
    it. Pinned by TestProtectPrivateStoreAppliesPolicyToProductCreatedDir.
  - **Row-24 fix landed** (atomic creation): `EnsurePrivateStore` (Windows) now
    creates the store via `windows.CreateDirectory` carrying an owner-only
    security descriptor (SDDL `D:P(A;;GA;;;<sid>)` — protected DACL, one
    non-inherited allow-all ACE for the token user, the same single-ACE policy
    `ownerOnlyDACL` produces) inside SECURITY_ATTRIBUTES, so the directory
    never exists without its policy; ERROR_ALREADY_EXISTS (lost race /
    appeared meanwhile) verifies rather than refuses. The create→set-DACL
    window is gone. Pinned by TestConcurrentFirstCreationNeverRefuses (8
    racing creators, none may refuse). Hosted status of
    TestSnapshotConcurrentCapturePublishesOneExactArchive this cycle: PASSED
    (exposure removed by the fixture change; the product-race fix + adversarial
    test ride this push).
  - **New phenomenon row 26 recorded before reaction:** fsacl suite 6/7 — the
    single failure, TestPreexistingStoreRefusedAndNotRewritten, was again the
    TEST's own scaffolding, a second layer under the one fixed last
    invocation: `grantTrustees`' protected-DACL replacement on the store
    auto-propagates inheritance removal to children, stripping the marker
    file's inherited-only access (final ReadFile: Access is denied). The
    PRODUCT refusal was verified correct by the log itself: DACL before/after
    identical (walkDACL: owner+Users both times), refusal text exact. Test
    fix: the marker gets an explicit owner-allow ACE via the same
    grantTrustees scaffolding, so its readability is independent of store-dir
    DACL churn; content assertion unchanged.
  - **New phenomenon row 27 recorded:** evidence_publication_test.go:35 now
    reaches a product refusal "independent evidence verification requires a
    POSIX execution host; Windows runtime is not supported"
    (driver_impl.go:543; sibling trajectory_verify.go:107) — unmasked by the
    W1 clear, class = §D.3 W-SHELL family, Stage 5.
  - **Stage 2 core landed** (§D.4): new `internal/pipeline/names.go` —
    `encodeEdgeID` (url.PathEscape base + explicit `:`→%3A + reserved-device
    first-char escape + trailing-dot guard; total and injective; no separator
    survives PathEscape, so GatePath output is structurally confined to the
    gates dir), applied inside `GatePath` so the `.tmp` staging name inherits
    it by construction (AC-NAME-1; the FINAL example
    `block-a->block-b` → `block-a-%3Eblock-b.gate.json` pinned). `LoadGate`
    gained the legacy raw-name read fallback; `SaveGate` retires the legacy
    shadow on first rewrite and refuses unsafe edge IDs (AC-NAME-2; mixed-name
    deck test). `checkBlockID` refuses the unsafe classes on the RAW ID
    (AC-NAME-4: the five-IDs IsLocal table incl. `../../other-idea` all
    refused; cosmetic-safe IDs grandfather), wired into `Manifest.Validate`
    (load); `ValidNewBlockID` strict allowlist wired into `runPipelineStart`
    (new pipeline creation, AC-NAME-3). Locally green: pipeline + app
    packages, incl. the two hosted ERROR_INVALID_NAME failures
    (TestPipelineAutoWalksToDoneUnderAutoLeft, TestPipelineAutoStopsAtActionBlockNeedsHumanGate)
    which now pass. **Remaining Stage 2 item (recorded, not silently dropped):**
    BlockWorkspace/GatePath refusal-signature conversion (`(string, error)`,
    ~15 call sites) — the current posture is safe for all product flows
    (manifests validate first; SaveGate refuses; encoding confines), but the
    mechanical in-function backstop AC-NAME-3 names is still owed.

- (2026-09-28T13:32–13:57Z, zcode-1; commit 29ee668 — first verified clock read
  13:32:49Z, commit timestamp authoritative) Invocation 6. Re-read packet
  (attestation above; hash unchanged), 00-prompt, living IMPLEMENTATION.md,
  FINAL §D.4. **Stage 2 completed**: GatePath and BlockWorkspace converted to
  error-returning §D.4 backstops (raw slug/edge/block ID unsafe classes
  refused, never sanitized, no legacy exemption); all 11 callers converted
  (launchBlockRound now ([]runner.Result, error); SeedBlockPrompt,
  blockCompleteFunc, planFinalized, run-block/start/auto sites, tests);
  Manifest.Validate now also refuses unsafe idea slugs at load. Assessed
  hosted **36429107801** (register): fsacl 8/8 GREEN — row-24 fix and
  ProtectPrivateStore restore fix HOSTED-VERIFIED on first hostile Windows
  execution; W1 signature 0; 14 red = prior 15 minus fsacl. The hosted run
  exposed two defects in MY invocation-5 Stage-2 code, fixed this invocation
  after recording: LoadGate's legacy fallback surfaced Windows
  ERROR_INVALID_NAME as a hard read error (it means the raw legacy name
  cannot exist on this OS — now not-found, via the build-tagged errInvalidName
  helper), and my legacy-fallback test tried to create a raw-'>' filename
  (impossible on Windows — now platform-conditional, PINNING that
  uncreatability as the encoding premise). The third pipeline failure is
  standing (row 29). RECORD CLARIFICATION per organizer correction: the
  invocation-5 statement that the two ERROR_INVALID_NAME app tests 'now pass'
  was LOCAL darwin evidence only; hosted Windows execution of the encoded gate
  names first occurred in 36429107801 and still carried the two fallback
  defects above — hosted green for the full Stage 2 rides this push. Stage 3
  §B dispositions NOT started (window consumed by Stage 2 completion + hosted
  assessment); they are the FIRST unit of the next invocation, before any
  Stage 3 code.

- (2026-09-28T14:00–14:2xZ, zcode-1; commit follows — first verified clock read
  14:00:34Z, commit timestamp authoritative) Invocation 7. Re-read packet
  (attestation above; hash unchanged), 00-prompt, FINAL §A/§B/§C/§H/AC-DUR,
  living IMPLEMENTATION.md. **Stage 3 core landed, dispositions BEFORE code**
  (the §B table above was written and committed before any durability edit):
  fsutil named-type contract (DirHandle + DirHandleOf fail-closed
  constructor; SyncDir — POSIX byte-identical incl. darwin F_FULLFSYNC,
  Windows = named refusal ErrDirEntryDurabilityUnsupported, the rooted sites'
  emitter; Windows SyncFile refuses directory handles as the mechanical
  backstop so no dir handle reaches FlushFileBuffers); A3 + B5 + C1 + C2
  pre-mutation entry refusals (Windows-confined, named, blocking,
  nothing-published; C1/C2 ordered before the renames per N8 — at entry,
  before even the stage file); D2/D3/D4a/D4b dir barriers routed through
  SyncDir; D1 = Windows no-op argued per §B (writeState already publishes via
  ReplaceSyncedFile = WT-move). Windows-tagged hosted pins:
  TestSyncDirContractRefusesFailClosed + TestRootedRowsRefusePreMutationOnWindows
  (A3/B5/C1/C2, asserting the refusal text AND nothing-published). Assessed
  hosted **36432547744** (register): pipeline package GREEN — Stage 2's two
  hosted defects are verified fixed; fsacl stays green; row 29 drift test NOT
  reproduced (intermittent — row stays open); 13 red = standing Stage 4/5/6
  families. Unix byte-identity: durability_refusal.go returns nil on every
  non-Windows GOOS; SyncDir delegates to the exact previous SyncFile
  semantics; SyncFile's directory guard is windows-only.

  ### Hosted cycle assessments (both completed; per-leg facts below)

  **Run 36172430646 @16824fd** (macos ok / ubuntu ok / windows FAIL, completed
  18:29:37Z; windows leg 12m34s): the H1–H7 probe bundle EXECUTED hosted —
  `ok parley-deck-cli/internal/winprobe 0.122s`, all 7 probe functions PASSED.
  Observed outcomes per FINAL §G mappings (workflow runs go test without -v, so
  recorder probes' t.Logf values are not surfaced; outcome-level facts are the
  assertions themselves — unknowns are named, nothing inferred):
  - **H3 OBSERVED as-expected**: runner volume asserted NTFS; directory move
    accepted, no-REPLACE failed on existing destination, MOVEFILE_WRITE_THROUGH
    accepted (any deviation fails the probe). Mapping applied: §A disclosure tiers
    pinned as written; coverage-envelope claim stands (NTFS confirmed).
  - **H7 OBSERVED as-expected**: read-only-destination rename-over rejected before
    clear, accepted after clear — the pinned pair held; §D.6 clear-before-replace
    is pinned by its probe (mechanics only). Also under H7: wait/usage FIRST HOSTED
    EXECUTION completed — internal/app ran to completion (165s) with zero failures
    attributed to wait_test.go/usage_ingest_test.go (the workflow runs the whole
    package; those suites therefore executed and passed). The comma-ok fixes are
    hosted-verified; the app failures that remain are fixture/runtime families
    (rows 16/18/23), not panics.
  - **H1 OBSERVED outcome-level**: the x/sys ACL chain round-trips hosted (asserted
    in-probe). Mapping: owner-only policy unchanged, mechanism viable. The
    parent/%TEMP% ACL dump VALUES are logged but not surfaced (unknown-in-band,
    fixture-placement data only). Doubly confirmed by 36172828285: the fsacl module
    Ensure/verify round-trip PASSED on real Windows.
  - **H2 OBSERVED outcome-level**: read-only directory handle opened; FlushFileBuffers
    raw error recorded in-log only (UNKNOWN-in-band). Diagnostic-only forever; when
    Stage 3 needs the raw string for the §B audit-of-record, run a targeted hosted
    `go test ./internal/winprobe/ -run TestH2 -v` step — never a mechanism change.
  - **H4 OBSERVED outcome-level**: probe ran; runtime.Version() value and
    os.ReadDir(regular file) error class are logged but not surfaced (UNKNOWN-in-band).
    §D.5 fix ships under every outcome; no blocker.
  - **H6 OBSERVED outcome-level**: git config --show-origin dump executed hosted
    (git present); VALUES unknown-in-band until the Stage 6 §D.7 pinning work.
  Windows red set: SAME 14 packages as the denominator (deterministic third cycle).

  **Run 36172828285 @ca5efef** (macos ok / ubuntu ok / windows FAIL, completed
  ~18:32:50Z): first hostile Windows execution of internal/fsacl (unwired module,
  original tests). Hosted facts: EnsurePrivateStore round-trip PASSED (MkdirAll +
  SetNamedSecurityInfo + walk-verify on real NTFS); pre-existing permissive store
  refused via DACL walk PASSED; symlink creation WORKS on the runner (no privilege
  error — the previously-unknown runner fact is now observed) and symlink/non-dir
  rejection PASSED; DenyRead deny-ACE PASSED (file became unreadable). ONE failure:
  the shared test's `perm=0777` assertion — Windows dir perms synthesize 0777, so
  that assertion is Unix mechanics; fixed by splitting it into fsacl_unix_test.go
  (the Windows expression of the same AC-PRIV-6 guarantee is the DACL round-trip,
  now asserted in fsacl_windows_test.go). 15 red packages = the same 14 + fsacl's
  that one test; W1 signature persists (81x) as expected while unwired.

## Stage 3 — §B disposition/audit table (recorded BEFORE durability code; FINAL §B is authority)

Locators verified at current HEAD (a93f71c+working tree, 2026-09-28T14:0xZ).
Dispositions are FINAL §B's refusal branch verbatim in intent; the
implementation column records what THIS implementation does per row.

| Row | Site (verified) | Disposition (FINAL §B) | Implementation action |
|---|---|---|---|
| A1 | `app/trajectory_verify.go:146` OpenFile O_EXCL → `:154` (dir sync `:138`) | Restructure (det-stage), DORMANT behind POSIX gate `:106-108` | none this stage — design of record; AC-DUR-3 conversion + stage-grammar pin = remaining Stage 3 item |
| A2 | `trajectory/verification.go:202` Root.OpenFile O_EXCL → `:211` | Refuse pre-mutation; dormant behind gate `:251-253` | existing reviewed gate refusal covers it (§C.3: takes away nothing); hosted pin = gate test (nothing-published by construction: gate fires before OpenRoot) |
| A3 | `trajectory/reservation_recovery.go:75` OpenFile O_EXCL → `:84`→`:67` | **Refuse pre-mutation** (Win-REACHABLE) | entry refusal in `publishReservationIntent` before the O_EXCL create |
| B1 | `app/trajectory_verify.go:118` Mkdir ×2 → `:125` | Restructure, dormant | as A1 |
| B2 | `app/trajectory_verify.go:190` Mkdir → `:193` | Restructure, dormant | as A1 |
| B3 | `trajectory/verification.go:264` parent.Mkdir → `:267` | Refuse pre-mutation; dormant behind gate | as A2 |
| B4 | `trajectory/verification.go:283` base.Mkdir(charge) → `:286` | Refuse pre-mutation; dormant behind gate | as A2 |
| B5 | `trajectory/reservation_recovery.go:36` dir.Mkdir → `:47` | **Refuse pre-mutation** (Win-REACHABLE) | entry refusal in `openIntentRoot(create=true)` before the Mkdir |
| C1 | `trajectory/parent_recovery.go:333` dir.Rename → `:336` (func bound `:311`) | **Refuse, ordered BEFORE `:333`** (N8) | entry refusal in `publishParentRecoveryWithSync` — before stage creation, hence before the rename; `defer dir.Remove(stage)` unaffected (nothing to clear) |
| C2 | `trajectory/parent_recovery.go:570` dir.Rename → `:573` | **Refuse, ordered BEFORE `:570`** | entry refusal in the recovered-parent publish func, same shape |
| D1 | `trajectory/unchanged.go:267` syncUnchangedState → `:279` (dir) | **Derived no-op**, argued: `writeState` already publishes via `ReplaceSyncedFile` = WT-move | Windows-confined no-op with the argument on record (Unix byte-identical); the re-sync adds nothing the WT-move publication does not already carry |
| D2 | `trajectory/parent_recovery.go:368` → `:308` (syncParentRecovery dir part) | Refuse — derives from refused C1 | dir part routed through `fsutil.SyncDir` → named Windows refusal (the emitter) |
| D3 | `trajectory/parent_recovery.go:605` → `:545` (syncRecoveredParent) | Refuse — derives from C2 | as D2 |
| D4a | `trajectory/reservation_recovery.go:420` syncIntent (dir parts `:63-67`) | Refuse — derives from A3/B5 | dir parts routed through `fsutil.SyncDir` |
| D4b | `trajectory/reservation_recovery.go:431` SyncFile(trajectory.json) + syncVerificationDirectory | Refuse — derives from A3/B5 | dir part via `fsutil.SyncDir`; the read-only-handle file re-sync cannot flush on Windows either — the whole barrier refuses (nothing rewritten) |

Audit of record: `fsutil.SyncDir` named-type contract — directory handles reach
entry-durability only as `fsutil.DirHandle` through `SyncDir`, which on
Windows is a named fail-closed refusal (`ErrDirEntryDurabilityUnsupported`);
the Windows `SyncFile` additionally refuses directory handles as a mechanical
backstop so no dir handle can silently pass through the file path. Unix
behavior byte-identical at every row (AC-DUR-5): all new refusals are
`runtime.GOOS == "windows"`-confined; SyncDir's Unix implementation preserves
today's fsync semantics exactly (darwin F_FULLFSYNC fallback included).

## Stage 4 review-scope evidence (organizer reminders, invocation 10 — for the independent reviewers)

**1. OpenSharedDelete and the owner's os.Root containment guarantee under
races.** The helper opens `name` only through a fresh `os.Root` on
`filepath.Dir(name)` + `root.Open(filepath.Base(name))`. Race analysis vs the
previous plain `os.Open` at the two converted readers (readState,
readStepHistoryFile), both of which enforce `Lstat → IsRegular → open →
f.Stat → os.SameFile` around the open:
- If the path is swapped to an escaping symlink between Lstat and the open,
  the OLD `os.Open` FOLLOWED it and opened the outside file, rejecting only
  afterwards via SameFile (post-hoc; the foreign file was opened). The rooted
  open REFUSES at the kernel (Linux RESOLVE_BENEATH / Windows rooted-open
  semantics): the outside file is never opened, and the SameFile check still
  runs afterwards. The guarantee is strictly STRONGER under the race, never
  weaker — no reliance on a pre-open lexical check.
- The owner's containment decision concerns the product's ROOTED operations
  (snapshot/verification os.Root sites, §A); those are untouched. The
  converted readers were always plain-path readers; the helper only narrows
  them to directory-scoped opens. POSIX semantics: root.Open is the same
  read-only open, byte-identical reads.
- Windows share semantics: os.Root opens carry FILE_SHARE_DELETE (FINAL §D.6
  basis), which is the helper's purpose — a concurrent ReplaceSyncedFile/
  WriteFileAtomic is no longer blocked by the reader. Empirical confirmation
  (the cycle_extension Access-denied clearing) rides the 17d5210 cycle.

**2. Lock-retry self-vs-foreign classification at every relevant open.** The
claim: every handle THIS process takes on a budget .lock file uses Go's
CreateFile share mode (FILE_SHARE_READ|FILE_SHARE_WRITE) with read and/or
read-write access, so no two of our own handles conflict; a sharing violation
on any of our opens is by construction a FOREIGN holder. The full audit of
.lock opens and their coverage: acquirePinnedKernelLock :159 (kernel-lock
open, O_RDWR — retry-wired), :200 (probe, O_RDWR — retry-wired), lockIdentity
:394 (identity read, read-only — retry-wired invocation 10 after the hosted
evidence showed it was the unwired failure site). Temp-file creates
(.lock-identity-*, .lock-origin-*, .origin-migration-*) are different paths
(racing their own O_EXCL, not share conflicts). readLockOrigin reads the
lock-ORIGIN file — originally unwired on a "no hosted failures" premise
that 36440560978 FALSIFIED: the origin read now hits sharing violations
during concurrent origin publication by a sibling parley process (invocation
11 wired openOriginFileRead with the same bounded retry and an honestly
narrower classification — the holder can be the product's own sibling
process, not only a foreign scanner; writer-side origin-publication
discipline recorded as follow-up). The classification is
supportable at every wired open; the diagnostic test passed hosted,
confirming the foreign-holder error class is ERROR_SHARING_VIOLATION and
both retry outcomes behave as specified.

## Advisory consult engagement — claude-1 identity consult (invocation 11, 2026-09-28T15:2xZ)

`inbox/claude-1-to-zcode-1_windows-portability_identity-consult.md` (advisory
only, read at fixed commit 9b880d3; not a review, not a signoff — its own
header says so). Engagement status:

- **Mechanism claim (its §1/§2):** row 30's attribution ("NTFS file-ID
  semantics; recreated file presents the same identity") is argued WRONG —
  the actual mechanism is Go's Windows lazy FileInfo: `os.Stat`/`os.Lstat`
  store the PATH and defer the identity lookup to the first `os.SameFile`
  call, which re-opens by path at comparison time; a pinned lazy operand
  therefore re-resolves to the REPLACED file and the guard is a structural
  no-op against replace-at-same-path. Handle-derived FileInfos
  (`File.Stat`, `Root.Stat`/`Root.Lstat`) capture NTFS identity eagerly. The
  claim is source-derived (Go 1.27.1 stdlib) with a determinism argument
  (failures identical across three cycles; only the two identity-dependent
  subtests red). It carries a stated version caveat (CI pins go 1.26).
- **My disposition:** the argument is compelling and the failure-set
  discrimination is strong, but it is NOT yet observation. Per the §6
  protocol I landed the discriminating hosted probe BEFORE any remediation:
  `internal/winprobe/samefile_pin_windows_test.go`
  (TestSameFileLazyVsEagerPin) asserts BOTH predictions — lazy pin compares
  EQUAL to the replacement (the bug) and eager pin compares NOT-EQUAL. Either
  prediction failing is a loud hosted fact that adjudicates the mechanism.
  **Row 30's remediation is deliberately NOT started until that probe's
  hosted outcome lands** — no fixture change, no product-guard change, no
  exclusion; the invariant stays exactly as red as it is.
- **Product exposure table (its §3) recorded, unactioned:** seven classified
  lazy-pin product sites (trajectory/state.go:225/:238,
  budget/step_history.go:178/:191, runner/action_identity.go:23/:40,
  budget/lock.go:306/:308/:387/:400/:478/:491, driver/cursor.go:120/:133)
  plus its unclassified same-shape list — all pending the probe outcome and
  the formal review. Its §5 fix proposal (pinned operand must be
  handle/Root-derived; a fsutil helper; fixture correction paired with the
  product fix, product FIRST if only one can land) is noted as the candidate
  design, not adopted yet.
- **§7 OpenSharedDelete secondary look:** share-flag claim CONFIRMED by its
  source read (root opens pass FILE_SHARE_DELETE; plain opens do not). Its
  containment critique — final-component-only containment, parent resolution
  unrooted, no origin pin — is accepted as a fair correction to my helper's
  doc comment; qualifying that comment and any origin-pin hardening are
  bundled with the same pending row-30 work so the identity story lands
  coherently, not piecemeal.

## Row-30 remediation (invocation 12, 2026-09-28T15:3x–15:4xZ) — product fix FIRST, fixture paired, per the consult and the requirement reminder

The hosted probe (36443342532) confirmed the lazy-FileInfo mechanism, so the
remediation landed as ONE unit: the product identity-pin sites, the helper
docs, and the fixture correction together — never a fixture-only mask.

**Invariant adopted (consult §5.1):** the pinned ("before") operand of every
`os.SameFile` guard must be handle- or `os.Root`-derived — never
`os.Stat`/`os.Lstat` (whose Windows FileInfo defers identity to comparison
time and re-resolves the path then). Greppable, platform-neutral, a
strengthening of containment, no build tags, no owner deviation.

**New fsutil helpers** (replacing OpenSharedDelete, whose doc the consult §7
correctly called overstated): `PinLstat` (eager Root-derived stat; honest
containment doc: FINAL-COMPONENT containment only, parent resolution is the
caller's, no origin pin) and `OpenPinned` (eager pin + delete-share open
through ONE directory handle, so pin and open cannot observe different
directories).

**Converted product sites (all lazy pins → fsutil.PinLstat / OpenPinned):**
trajectory/state.go readState (OpenPinned), budget/step_history.go
readStepHistoryFile (OpenPinned), runner/action_identity.go hashActionFile,
budget/lock.go verifyLockIdentity + lockIdentity read + readLockOrigin,
driver/cursor.go LoadCursor, app/evidence_table.go, app/evidence_verify.go
(:79), evidence/report.go, budget/ledger.go, runner/cycle_budget.go,
protocolpacket/source.go (both pin sites), budget/cycle_history.go (both
same-dir-check operands). **Already eager, no change:** trajectory/reconcile.go
(dir.Lstat), trajectory/parent_recovery.go (root-derived), all
snapshot.go/verification.go production paths (the consult confirmed these
were never affected). **Open item:** evidence_verify.go:441 — a second
Lstat whose SameFile role was not yet traced; recorded, not silently
cleared. **tui/live.go:2468** — cached `c.info` pin whose origin needs the
same audit; recorded.

**Fixture correction (paired):** both revalidation tests now take `prior`
from `root.Lstat(rel)` — matching what production callers supply — so the
{root,inode} subtests exercise the real invariant on every platform.

**§D.6 lock-origin retry justification (requirement reminder):** the signed
scope is "third-party holders, after self-held vs foreign-held
classification". For `openOriginFileRead` the classification is: (a)
SELF-HELD (this process) is structurally impossible — every handle this
process takes on any file shares read|write, so our own handles cannot
produce a sharing violation on our open; (b) the observed hosted holder is a
SIBLING parley process inside its bounded origin-publication window — a
foreign process by the classification's process-level meaning, transient by
construction (the publication completes; if a sibling wedged mid-publication,
the retry exhausts LOUDLY at 250ms rather than masking it — the origin/token
verification chain re-runs after recovery). A scanner is the other foreign
class. No self-held handle is masked: if the classification ever failed, the
failure mode is the loud named exhaustion, not a silent wait. The
writer-side origin-publication discipline remains recorded follow-up.

Local evidence: full builds green (darwin + windows cross), vet green on all
touched packages; targeted suites green (trajectory revalidation incl. the
re-pinned fixture tests, budget lock/reserve/step/cycle/origin 10.8s,
protocolpacket, evidence, driver cursor, runner). Hosted execution of the
remediation rides this push; the {root,inode} subtests are expected to go
green on Windows with it.

## Row-30 open audit items CLOSED (invocation 13, 2026-09-28T15:5xZ)

- **evidence_verify.go:441** — CLEARED, no fix needed: the os.Lstat there
  feeds only an `IsDir` check (no os.SameFile guard); the lazy deferral
  affects only the identity fields (vol/idxhi/idxlo) — mode/size/dir
  attributes are captured eagerly by GetFileAttributesEx. No exposure.
- **tui/live.go:2468** — CLEARED, already correct: the cached `c.info` pin
  is assigned from `f.Stat()` (handle-derived = EAGER), and the comparison
  operand `st` is also a `File.Stat` — both eager; the log-rotation
  detection works on Windows as-is.

## Stage 6 sweep, second tranche (invocation 13)

**Rows 8/20 (chmod fixtures → ACL helpers), all converted:** the new
`fsacl.DenyWrite`/`fsacl.AllowWrite` pair (Unix: chmod 0555/0755,
byte-identical to the historical fixtures; Windows: deny-ACE on
GENERIC_WRITE — which on a directory means adding files/subdirs and
deliberately excludes FILE_DELETE_CHILD/DELETE so cleanup removals still
work — and an AllowWrite restore that rebuilds the DACL from every
surviving allow, dropping deny ACEs). Converted sites:
evidence/tree_report_test.go (TestSaveUnwritableDirFails,
TestFailedSaveLeavesNoReport), app/organizer_test.go (read-only deck —
also retired a t.Skipf into a t.Fatalf), app/driver_checks_test.go
(evidence-write veto), driver/impl_test.go (unreadable dir — Skipf retired),
driver/phase_event_test.go ×3 (unreadable 00-prompt/IMPLEMENTATION/FINAL →
fsacl.DenyRead). Unix behavior byte-identical (helpers are chmod on POSIX);
Windows execution rides this push. Two pre-existing t.Skipf fallbacks were
converted to loud failures per the no-silent-skip discipline.

**Row 21 (mode-synthesis):** named as a product-semantics ambiguity in
`inbox/zcode-1-to-all_windows-portability_row21-mode-synthesis.md`
(synthesized-vs-requested mode in the archive header/digest; options (a)
documented Windows expression vs (b) record the product's requested mode);
no fixture or product change until review settles it — per FINAL row 21's
"deliberate decision, not a silent fixture change".

## Row-21 mode-synthesis disposition (invocation 14, 2026-09-28T16:1xZ) — kimi-1 consult engaged

Kimi's advisory reply (`inbox/kimi-1-to-zcode-1_windows-portability_mode-consult.md`,
read at fixed commit 4827e67, advisory only) was engaged independently. My
assessment: **sound under the frozen FINAL, adopted.** Its load-bearing
points, independently checked against the tree:

- Windows has exactly one chmod-mutable, stat-observable mode dimension —
  owner-write ↔ FILE_ATTRIBUTE_READONLY (synthesized 0666/0444; Go's Chmod
  maps only S_IWRITE). The digest's per-file `file:%o` serialization
  (tree.go) already binds it — the Windows mode signal is NARROWED, not
  lost, and asserting it is consistent with §D.6 (the read-only attribute
  is an OS-enforced semantic the FINAL already relies on).
- My option (b) (record requested mode) is REJECTED for the reason the
  consult surfaced and I had not: CaptureSnapshot/TreeDigest operate on
  ARBITRARY USER WORKTREES — the product did not create those files, the
  requested mode is unrecoverable, and inventing it would fabricate
  metadata against §D.9's spirit and rewrite the archive-v1
  observed-mode compatibility boundary.

**Disposition implemented (test-side only, product unchanged, no skips,
unconditional assertions on both platforms):** (1)
TestTreeDigestModeChangeChanges mutates the dimension each OS has — Unix
0600→0700 byte-identical, Windows 0600→0400 (read-only attribute;
deliberately os.Chmod, not a DACL helper — a deny-ACE is invisible to the
synthesized FileMode); the d1≠d2 assertion stays unconditional. (2)
TestSnapshotRevalidationRetainsIdenticalMaterialAfterTimestampChange pins
the platform-observed constant (0600 Unix / 0666 Windows) — pinned, not
derived, so synthesis drift fails loudly. (3) The round-trip
executable-mode assertion (same family, snapshot_test.go) pins
0700/0666 the same way — the archive records OS-observed mode and the
restore applies it. Hosted execution is the first Windows evidence (the
consult itself ran nothing); the disposition is review-visible for Phase 6.

## Concurrent-replacement investigation (invocation 14) — §D.6 scope respected, consult filed

The recurring cycle_extension Access-denied is analyzed to the delete-pending
window (reader open racing a same-process MoveFileEx-replace; opens in the
window fail with ERROR_ACCESS_DENIED — NOT the sharing-violation class §D.6's
retry covers, and the holder is our own process, not a third party). Per the
scope discipline (and my own lock-open documentation that Access-denied is
never retried), the retry was NOT broadened silently; the three options
(narrow delete-pending retry / writer-side lock ordering / accept-as-loud)
are filed as `inbox/zcode-1-to-all_windows-portability_delete-pending-consult.md`
for peer input, and independent work continues.

## Delete-pending consult RESOLVED by claude-1's reply (invocation 15, 2026-09-28T16:3xZ) — my mechanism retracted

Claude's advisory (`inbox/claude-1-to-zcode-1_windows-portability_delete-pending-consult.md`,
code at f1a7f3e, advisory only) **disputes and defeats my delete-pending
attribution**. I verified its three primaries independently before accepting:

1. The hosted message is a BARE errno ("persist operator cycle grant: Access
   is denied." — no path, trailing period): the only bare-errno producer in
   the persist chain is `replace_windows.go`'s raw `MoveFileEx` return — the
   failure is the WRITER's rename-over-target, not a reader open.
2. `step_history.go`'s post-read check paired a handle-derived `opened` with
   a LAZY `os.Lstat` AFTER operand — and Go's `loadFileId`
   (types_windows.go:321, verified in my local toolchain source) opens the
   path with `CreateFile(pathp, 0, 0, …)`: **dwShareMode=0, share-nothing**.
   Our own reader's `os.SameFile` held that exclusive handle for the
   duration of one `GetFileInformationByHandle`; the concurrent guarded
   `MoveFileEx(REPLACE_EXISTING)` over `policy.json` was refused with
   ERROR_ACCESS_DENIED. **Self-held** — and unguarded readers
   (`LoadCycleBinding` before the guard, `InspectCycleBudget`) make it a
   product property, not a test artifact.
3. My "absent in 36439743640 = fixed" reading was wrong — racy absence, as
   my own row already cautioned.

**Corrections to my record:** the delete-pending mechanism is RETRACTED; my
consult option (a) (retry the reader open) is WITHDRAWN — it would have
retried over our own process's handle, precisely the self-held masking §D.6
exists to prevent. Claude's R1 was implemented: the AFTER operand at
`step_history.go` (post-read check) is now `fsutil.PinLstat` — the row-30
invariant applied to the after side: no path re-open, no zero-share open,
Root share-delete, one line, §D.6 untouched, no retry, strictly
strengthening. **R3 pin landed** (`TestPinLstatOperandComparesWithoutPathReopen`,
platform-neutral, no skips): a PinLstat operand compares via SameFile after
the path is gone (no re-open); the lazy-operand contrast is asserted on
Windows where the property exists (POSIX os.Lstat is eager). Hosted
execution is the replace-side reproduction.

**AFTER-operand inventory (scoped, no blanket clearance):** the dangerous
shape is a SEPARATE path-derived after-operand (its own stat). Converted
this invocation: `step_history.go` post-read check, `runner/action_identity.go`
post-read check (same shape). Verified SAFE (after operand is the handle's
own `f.Stat()`, no second open): lock.go ×3, ledger.go, cycle_budget.go,
cursor.go, evidence_table, evidence_verify:92, report.go, state.go:237,
verification.go:235, protocolpacket:385, snapshot.go's pinned-compare sites,
tui/live.go (both handle-derived, audited invocation 13). Root-derived
after-operands (eager + share-correct): parent_recovery.go:157/:180
(`dir.Lstat`). **archive_stability audit CLOSED (invocation 16):** `sameSnapshotFile`
compares caller-supplied FileInfos; the audit found TWO lazy operands on its
paths and converted both: `checkSnapshotFileStability`'s own `named`
(os.Lstat → PinLstat — the row-30 AFTER rule, same zero-share side effect)
and — a MISSED SWEEP SITE — `inspectSnapshotFileAttempt`'s `info`
(os.Lstat → PinLstat), the lazy BEFORE operand feeding the stability check's
`initial` on the inspect path (the earlier sweep had verified the
capture/verify paths; the inspect path slipped). Both conversions are the
uniform one-liners; snapshot/stability tests green locally.

**§D.6 gap named, not widened:** at the replace site a target-handle
conflict surfaces as ERROR_ACCESS_DENIED (class 5), outside §D.6's signed
sharing-violation retry. No foreign-holder evidence exists on this path
(both observed cases were self-held, now removed); the residual is recorded
here and NO retry class is broadened. If a future class-5 replace failure
shows a foreign holder, that is the moment for an amendment request.

## Rows 3/4 re-exec fixture conversion — first tranche (invocation 16, 2026-09-28T16:4xZ)

The §D.9 test-binary re-exec infrastructure landed in internal/app:
TestMain dispatches a fixture role when the binary's BASENAME is a known
fixture (writeReexecFixture installs a copy of the test binary at the
fixture path — the .exe suffix on Windows resolves through PATHEXT via the
product's LookPath). Converted this tranche:
- `writeFakeParleyDeckSkill` (row 3's root cause — the hosted
  "executable file not found in %PATH%" signature behind the version-all
  failures, row 16): the status/JSON shell script is now
  `fakeParleyDeckSkillMain` (Go port, identical JSON bytes with the
  --project argument substituted).
- `writeFakeLegacyParleyDeckSkill`: the legacy variant shares the
  basename, so the behavior is selected by an inherited env marker
  (PARLEY_FAKE_SKILL=legacy) — the test process's env reaches the child
  through the product's exec; `fakeLegacyParleyDeckSkillMain` ports the
  --version/unknown-command script.
Local (darwin): all four version-all tests pass THROUGH the re-exec path
(the copies are exercised, not the shell). Windows execution rides this
push. Unix behavior equivalent (the §D.9 re-exec port is the sanctioned
replacement for the shell fixture; same commands, same outputs).

**Named remainder (rows 4/23, the launch/agent-CLI fixtures):**
`writeFakeRoundAgentCLI` (complex: awk-over-stdin artifact writer — needs
the role-spec ported to Go), `writeFailingRoundAgentCLI` (simple
version-or-exit-1), the :1686/:1710/:1735 formatted bodies, :1398
(cat/exit-7), and the launch_test.go / protocol_context_test.go /
telemetry_test.go fixtures (~13 per the ledger). Design sketched for the
parametric ones: a sibling `.role` spec file beside the binary copy (the
role runner reads its behavior spec from os.Args[0]+".role"), so
per-fixture versions/modes don't need per-name dispatch.


## Rows 3/4 second tranche + identity audit closure (invocation 16, 2026-09-28T16:4x–16:5xZ)

The generic **`.role` spec mechanism** landed for parametric fixtures: the
behavior spec lives beside the binary copy as `<name>.role`
(version-or-fail / version-or-drain / drain-exit specs), so per-test
versions and exit modes need no per-name dispatch; `writeRoleFixture`
returns the installed path (the .exe suffix on Windows — direct-exec config
paths must use it, as the alpha-fixture call site now does). Converted this
tranche: `writeFailingRoundAgentCLI` (version-or-fail) and the
consensus-signoffs drain-exit-7 fixture. The consensus request-signoffs and
version-all suites pass locally through the re-exec paths (20.1s).

**Remainder (rows 4/23):** `writeFakeRoundAgentCLI` (the awk-over-stdin
artifact writer — the largest port), `writeFakeForgedSignoffCLI` and the
:1686/:1710/:1735 formatted bodies, and the launch_test.go /
protocol_context_test.go / telemetry_test.go fixtures.

## Runaway-fixture postmortem + rows 3/4 completion (invocation 17–18, 2026-09-28T16:55–17:2xZ)

**The runaway (organizer SIGTERM ×2, 17:07:34Z with 67 descendants and
17:15:18Z with 40):** my re-exec role-spec dispatch had TWO defects that
composed into unbounded recursion. (1) The role lookup fell back to
`m.Run()` — the ENTIRE test suite — inside a fixture copy whenever a role
could not be resolved; (2) the `round-agent`/version spec parsers used
whitespace field-splitting, so a version containing spaces
("codex test 1.0") failed `len(fields)==2` → fallback → the copied binary
ran the whole app suite inside itself, whose tests installed more copies
(TestRunRecordsResolvedRuntime's codex fixture was the observed vector).
**Repair (both defects):** a STRUCTURAL recursion guard — a renamed copy
(its executable basename is neither a dedicated role nor the standard test
name) that cannot resolve its role prints a diagnostic and exits 70,
NEVER entering the testing framework (applied to BOTH dispatchers: app's
runFixtureRole and driver's gitprobe TestMain); role resolution now uses
os.Executable() (argv[0] may be a bare PATH name); and the version-carrying
specs parse the raw remainder after the mode word, preserving spaces.

**Claim reconciliation (the 17:09:15 report):** the 16:58:52Z broad batch
(`-run 'TestConsensusRequestSignoffs|TestRun|TestVersionAll|TestAgents'`)
NEVER completed for me — its output was lost to the first runaway (the
organizer's 17:07:34Z observation of an 8-minute chain); it is NOT green.
The 17:09:02Z "ok 13.331s" was a real completed run of the NARROW filter
(`TestConsensusRequestSignoffs|TestVersionAll`) that started AFTER the first
tree was killed — legitimate for that filter, but it ran with the
recursion-capable code and is superseded by the post-fix verification below.
The 17:09:40Z gitprobe failure and the ~17:10Z cancelled batch (the second
runaway, killed 17:15:18Z) are as the organizer stated.

**Post-fix verification (bounded, exit status preserved, no output pipes,
process count checked after each):** TestRunRecordsResolvedRuntime (the
recursion vector) — first exposed a second port defect (my artifact-marker
extraction was prefix-anchored where the shell's awk matched unanchored;
fixed to the exact index+trim semantics) — now **PASS, exit 0, 2.2s**;
TestVersionAll PASS (0.5s); TestConsensusRequestSignoffsBlockStops PASS
(2.3s); TestGitTreeCleanSetsOptionalLocksOff PASS (0.3s, driver — the
historical t.Skip on Windows retired with the re-exec port); wider bounded
family run (TestConsensusRequestSignoffs|TestAgentsExec|TestRunAnswer|
TestRunRecords) **PASS 19.7s, exit 0**; zero stray test processes after
every run.

**Rows 3/4 fixture unit COMPLETE:** all extension-less shell fixtures in the
tree are now re-exec ports — app: parley-deck-skill (+legacy via env
marker), round-agent (the awk-over-stdin artifact writer, fully ported with
unanchored-marker semantics), forged-signoff (+exit-7 variant via role-file
edit, replacing the script-body edit), rewrite-signoff, parametric signoff,
version-or-fail/drain, drain-exit; driver: the git PATH-shim (record-env
role). tree_report's shebang is file CONTENT (a digest fixture — not an
exec'd fixture; no conversion needed). The launch/protocol_context/telemetry
files named in the original ledger no longer exist (row 4's locator was
stale — the live out-of-app site was gitprobe only). Windows execution
rides the push.

## claude-1 replacement-recurrence follow-up acknowledged (invocation 18, 17:2xZ)

Its self-correction is accepted and recorded: **R1 was correct but was not
the repair** — after R1 there is no product zero-share open left on that
path (Claude independently reproduced my negative trace at 938bf52), and the
36451945257 timing signature (exactly one Access-denied; the :134 assertion
holding proves exactly one extension landed; the guard serializing the
extends) shows the first replace was refused while a microseconds-later
second succeeded — falsifying the persistent-holder family and the
read-only-trap, leaving exactly two candidates: **C1** a go1.26.8-specific
stdlib open (the CI toolchain; EVERY Go source line cited so far was
go1.27.1 — the version gap is real and named) and **C2** a foreign holder on
the source/target (first-attempt-fails/second-succeeds is the scanner
shape). No second repair is proposed by Claude and none is implemented here:
the D1–D3 discriminating diagnostic is NOT landed while kimi-1
independently checks actual-version sources and diagnostic sufficiency
(notified 2026-09-28; no duplication of that investigation from this
implementer). Standing constraints accepted: any later diagnostic must be
failure-path-only with probes unreachable on the success path, reported
never acted on, no §D.6 broadening without foreign-holder evidence, and no
simultaneous-attribution overclaims from racy later samples. The
concurrent-replacement issue stays OPEN in the ledger.

## kimi-1 CI-toolchain consult engaged (invocation 19, 17:3xZ) — C1 refuted; diagnostics corrected, not yet landed

Kimi's consult (`inbox/kimi-1-to-zcode-1_windows-portability_ci-toolchain-consult.md`,
authentic go1.26.8 sources sha-pinned and mirror-verified) is engaged as
follows: **C1 (a stdlib-internal open at the CI toolchain) is REFUTED by
direct source comparison** — the identity/share/replace surface is identical
or wire-equivalent between go1.26.8 and go1.27.1, and the failing call is
x/sys v0.36.0's raw MoveFileExW (toolchain-independent by construction).
Claude's §3 reader clearance therefore transfers to the CI toolchain. **The
"exactly two candidates" framing is corrected**: the residual space is
foreign-handle (target or source), filter-level denial with no user-mode
handle, delete-pending, or attribute — several invisible to in-process
probes. **D1/D2 as specified are insufficient** (temporal skew at the
phenomenon's own timescale; the share=0 probe detects ANY handle, not the
rename-blocking condition; D2's package-global counter conflates files) —
the corrected D1′/D2′ (P1/P2/P3 full-share battery per side, per-path reader
registry sampled before the probes, nanotime-stamped, failure-path-only,
report-never-act, grading-not-classification) is recorded as the design of
record for the diagnostic unit, which rides a LATER coherent batch after
this invocation's fixture/refusal units per the task contract. The observed
occurrence's errno is already 5 (bare Errno messages preserve the
distinction). Claude's ReplaceSyncedFile os.Lstat-swallow latent defect is
independently endorsed by kimi and queued as its own small signed-scope
repair. ledger.go:373's plain os.Open on ledger.json (no share-delete at
either toolchain) is flagged for the sweep ledger. **Regression record,
per the task instruction: the e8ce8db guard false-positive KILLED the REAL
app test binary hosted (0.068s package failure) — that is a REGRESSION I
introduced, not proof of correct dispatch; only the later 92ab958+ hosted
legs can evidence the corrected behavior.** No retry broadening; no
probe-time state promoted to failure-time attribution.

## Corrected sharing-failure diagnostic batch landed (invocation 20, 2026-09-28T17:59–18:0xZ)

Per the engaged kimi-1 corrected battery and claude-1 constraints (both
advisory), scoped under the frozen FINAL, declared TEMPORARY with a removal
plan: **D1′** — `ReplaceSyncedFile`'s MoveFileEx failure path now wraps the
error with a report-only diagnostic (`replace_diag_windows.go`): the RAW
errno value, per-side (staged + target) existence and attribute word (P1),
a DELETE-access-with-FULL-share probe (P2 — mirrors the rename's own
requirement; succeeds over benign share-delete readers, fails only when a
handle denies delete-sharing), a READ_ATTRIBUTES full-share control (P3 —
delete-pending/hard-deny family), the D2′ in-flight reader sample taken
BEFORE the probes, and a nanotime stamp. Every outcome is labelled
**grading-not-classification**: the grades name probe-time state (t1 > t0)
with the unclassified case reported AS unclassified — no probe-time sample
is promoted to failure-time attribution, nothing is retried, the wrapped
error still fails the operation, the success path is untouched (the probes
are unreachable there — they are themselves opens, exactly what row-30
removed from the success path). **D2′** — a per-path in-flight reader
registry (fsutil.DiagReaderEnter/Exit) registered by the budget
`readStepHistoryFile` (keyed by cleaned path so a policy.json block is not
conflated with other files' readers; enter/exit with deferred exit).
**REMOVAL PLAN:** temporary — once the concurrent-replacement class is
identified from hosted evidence, the probes and registry are removed (or
narrowed per review); only the structured errno in the wrapped error may
remain. **Stat-error swallow repaired** (claude-1 §6, kimi-endorsed,
verified by my own read): a FAILED stat in the read-only-clear branch is no
longer silently conflated with absence — genuine absence proceeds, any
other stat error is surfaced (`stat replace target before move`), preserving
the diagnostic the old code discarded. Bounded local: fsutil + budget
suites PASS exit 0 (4.0s); darwin + windows cross build/vet green.


## Rows 4/16/18 tail + config residuals (invocation 20, 18:0xZ)

- **TestLoadAgentSpecsLayersAndTracksSources** FIXED: the fixture's
  `{root}/bin/extra` expands VERBATIM (the ExpandPlaceholders semantics);
  the wants are built by the same substitution, not Join normalization.
  Bounded PASS. **TestAgentsExecRecordsManualLaunch + family** FIXED: the
  measuredFixture shell scripts (write-artifact / touch / noop) are now the
  `agent-script` role spec — the extension-less fixture is a re-exec binary
  copy; bounded PASS across TestAgentsExec* (1.5s).
- Remaining known app-family reds recorded for the next unit:
  TestBudgetWorktreeInspectCLIReportsRetainedRegistration (likely the row-10
  porcelain comparison again, CLI variant), TestBeginBoundAllocates…,
  TestExecTelemetry*, TestDriverAdapterProtocolPrecheck…,
  TestFixtureAutoDrive…, TestTrajectoryHelper… — the designed-refusal and
  fixture-portability mix.


## Designed-refusal family reconciliation — first tranche (invocation 20, 18:0xZ)

**Reservation-recovery family (§C.1 blast radius, ~40 failing subtests):**
the recovery tests presuppose a PUBLISHED reservation intent; on Windows the
precharge publication refuses by signed design (A3+B5, §C.1), so the child
process can never reach the interrupted boundary. The platform-true
expression landed at the `interruptedReservationFixture` boundary: on
Windows the child's output must carry the designed §C.1 refusal text, the
exit must be non-zero, and the nothing-published state is asserted (no
reservation-intents directory) — then the test returns early via an
APPLICABILITY flag, **not a t.Skip** (the census stays clean; the refusal
invariants are real assertions that fail if the product regresses). The
refusal itself additionally remains pinned adversarially by the fsacl/§B
hosted suites. Unix: the full recovery flows unchanged (all
TestReservationRecovery* green locally, 20.6s, exit 0). Remaining family
members recorded for the next tranche: the evidence-package refusal tests
(dir-fsync §C refusals at refusal_test.go — same shape: the tests exercise
POSIX publication flows that refuse on Windows), TestBeginBound…,
TestExecTelemetry*, TestDriverPrecheck*, TestTrajectoryHelper…, and the
readonly organizer fixture.


## Designed-refusal reconciliation tranche 2 (invocation 20, 18:0xZ) — evidence/app refusal lifecycles

The refusal-publication lifecycle family (§B SyncDir refusal emitter): on
Windows `RetainVerificationRefusal` refuses with the named
directory-entry-durability error and must leave no state. Wired as the same
early-return applicability pattern (NO skips, census clean):
- `internal/evidence/refusal_test.go`: `retainRefusalPlatformTrue(t, dir)`
  asserts the designed refusal text + nothing-created state, returns false;
  injected after each affected test's dir creation (6 sites wired; the
  retain-failure paths already negative-tested keep their Unix shapes).
- `internal/app/evidence_refusals_test.go`: `appRefusal` returns an empty
  sentinel on the Windows designed refusal; all 4 callers early-return.
Unix suites green (TestRefusal* both packages, exit 0). The durable-
publication refusal itself remains pinned by the fsacl/§B hosted suites.


## Runner shell-role re-exec (invocation 20, 18:1xZ)

The runner fixture family's `/bin/sh -c` fixtures get the §D.9 platform
expression: `shellPath()` returns /bin/sh on Unix; on Windows a copy of the
test binary named `sh.exe` whose TestMain dispatches `-c <script>` to the
exact emulations the family uses (exit N; exec sleep N; cat "$1" + child
output; drain). Same recursion-guard discipline as the app package (a
renamed copy exits loudly, never m.Run()). Unhandled scripts exit 70 loudly
— the census-visible failure mode if a future fixture adds a script shape.
Full runner suite green locally (96.5s, exit 0). Windows execution rides
the push; any unhandled-script failures there name themselves.


## Genuinely readonly organizer fixture + platform-path gaps (invocation 21, 18:2xZ)

**The readonly organizer fixture — ROOT CAUSE found and fixed:** the
invocation-13 regression was NOT DenyWrite's deny semantics but
`fsacl.AllowWrite`'s restore: `removeDeny` re-applied the rebuilt allow set
with the PROTECTED_DACL flag, STRIPPING the directory's inherited ACEs and
bricking cleanup traversal (openfdat: Access is denied on child dirs).
Fix: `removeDeny` applies without the protected flag (the rebuilt grants
merge with parent inheritance). With the restore repaired, the organizer
test now uses the GENUINE native deny (`fsacl.DenyWrite` — chmod is a no-op
for write access on Windows and the historical pass proved nothing there;
the deny-ACE covers add-file/add-subdirectory, deliberately not
delete-child, so cleanup removals survive) with `AllowWrite` restoring.
The t.Skipf is retired into a loud t.Fatalf. Bounded PASS. **The
write-detection snapshot misreport from invocation 13 is explained by the
same root cause** (the broken restore left the tree unreadable, not the
brief writing).

**TestAgentsExecRetainsFailedStart:** the test removed `test-agent`
(extension-less) — a no-op on Windows where the fixture is
`test-agent.exe`, so the launch SUCCEEDED and inverted the test.
measuredFixture now records the actual installed path; the removal uses it.

**TestBudgetWorktreeInspectCLIReportsRetainedRegistration:** canonical
comparison applied (the row-10 family's CLI variant).

**Fixture-completion wording reconciled (per the organizing correction):**
the earlier "every extension-less shell fixture in the tree is a re-exec
port" claim was overbroad — the runner package's `/bin/sh -c` fixtures
(telemetryShell/interactiveFixture) were subsequently discovered and ported
in invocation 20 (`shellPath()` + the sh.exe role). The truthful scoped
statement: all DISCOVERED extension-less/shell fixtures across internal/app
and internal/runner are re-exec ports as of 7e7a340; any future grep that
surfaces another family extends the same pattern. No coverage verdict is
claimed from counts alone.


## chargeFixture applicability tranche (invocation 21, 18:3xZ) — semantic mapping explicit

Per the organizing correction, the semantic mapping for every
chargeFixture-based conversion is recorded: the ORIGINAL assertion family of
these tests is the recovery/validation of a CHARGED CYCLE (published
reservation-intent → charge → state transitions) — every step presupposes
the §D.1-durable publication, which on Windows refuses by signed §C.1
design BEFORE any mutation ("refusing before any file is written, nothing
was published"). The setup is therefore GENUINELY UNREACHABLE on Windows —
not inconvenient. The retained substantive Windows assertions at the
boundary (real, regression-sensitive): the designed refusal TEXT from
OpenCycleSession, and the nothing-published state (no reservation-intents
directory). Negative/adversarial cases within these tests that could run
WITHOUT publication (e.g. pure refusal-shape checks) remain applicable and
are covered by the fsacl/§B hosted pins — none were converted away.
chargeFixture sets a package applicability flag; 4 caller files wired with
early returns (captured_test's helper has typed returns; state_test's own
callers included). Unix: full TestReservationRecovery+TestPersistent+
TestCaptured families green (123s, exit 0). kimi-1's independent
applicability challenge (announced, at fixed c39d81c) had NOT arrived at
this commit; it will be engaged on arrival and may revise this pattern —
recorded as pending advisory input, not settled.


## kimi-1 applicability consult ENGAGED (invocation 21, 18:3xZ) — F1/F2 acted on immediately

The advisory coverage challenge (at fixed 4952e1d/c39d81c) was read on
arrival and its blocking findings acted on before any extension of the
disputed pattern:

- **F1 (adopted resolution (a)):** the evidence retention path created
  `.parley-runtime` scaffolding BEFORE the barrier refused — leaving litter
  despite the refusal, contradicting §C's "trigger is before the mutation",
  and making my nothing-created assertion predictably red hosted. Fixed in
  the PRODUCT: a pre-mutation Windows gate at the top of
  `RetainVerificationRefusal` (mirroring the trajectory family's
  refuse-before-OpenRoot shape) — the named durability error now fires
  before ANY directory creation. (Kimi's option (b) — accepting scaffolding
  litter — rejected as weakening §C's wording in effect.)
- **F2:** the three missed retain-success sites wired with the same
  platform-true gate (the changed-record half; the missing-guard-origin-lock
  retention; the app-side stopped-parent helper — the last via an explicit
  Windows early return with the mapping comment). Bounded PASS all.
- **F3 (acknowledged, not yet acted on):** the bypassed-but-applicable
  read-path assertions (Preview adversarial subtests constructible
  test-side; Inspect visibility/alias/conflict halves; the five CLI
  rejection cases) and F4's missing pins (RecoverReservation apply refusal;
  Preview pristine-store sanity) are RECORDED as the next reconciliation
  units with kimi's construction recipes — NOT silently dismissed. The
  tranche's own criterion ("cannot exist on Windows") concedes these; the
  early returns function as exclusions-in-effect for them until salvaged.
- **O1/O3 acknowledged:** the census cannot see applicability-flag early
  returns; no-coverage verdict is claimed from counts; the
  fixture-completion wording was reconciled earlier this invocation.


## F4.2 pin landed (invocation 21, 18:3xZ)

Preview's Windows entry point anchored: on a pristine store the read path
yields the designed missing-intent refusal and publishes nothing
(TestPreviewReservationRecoveryPristineStoreRefusesOnWindows — windows-
tagged, no skips). F4.1 (the RecoverReservation apply-refusal pin) and the
F3 read-path salvages remain the named next units.


## F3 salvage + F4.1 pin delivered (invocation 22, 18:5xZ)

Per the applicability consult's construction recipes:

- **F3b** (`testWindowsRefusalReadPath`): the pristine-Inspect read-path
  assertions (empty entries, nothing created) run on Windows BEFORE the
  retain probe — restoring the bypassed no-mutation coverage at the first
  wired site.
- **F3e** (`testWindowsRefusalCLIRejections`): all five CLI rejection
  cases run on Windows before the not-applicable early return — four need
  no retained record; the fifth uses a fabricated valid-format digest and
  rejects before any barrier. The blanket early return no longer bypasses
  the Windows-exercisable CLI contract.
- **F3a + F4.1**
  (`reservation_recovery_readpath_windows_test.go`): the charged-cycle
  state is constructed TEST-SIDE (canonical intent bytes via plain
  Mkdir/WriteFile — the FINAL §D.4 cross-OS premise; no product
  publication, no barrier), and on it: the Preview adversarial refusals
  are restored (missing-intent pinned to its exact designed text;
  partial-intent/changed-root/changed-limits each assert a non-nil
  refusal — the exact per-case texts print via t.Logf on the first hosted
  run and get pinned from that evidence rather than guessed); and the
  F4.1 apply pin: `RecoverReservation` must refuse with the §B barrier
  text and leave the original state byte-identical.
  **Mapping honesty:** the original ten subtests' full downstream
  assertion chains (inspect/replay/refuse-another-attempt after a
  successful recovery) are NOT restorable on Windows (they presuppose a
  completed recovery); what is restored is the read-path refusal family +
  the apply barrier pin — the parts Kimi's trace proved reachable. The
  chargeFixture tranche audit: the same split applies — the setup-level
  publication is genuinely unreachable; read-side validation tests that
  don't need the charge (e.g. Preview on constructed state, now landed)
  are restored; tests whose EVERY assertion reads the charged row's
  post-recovery state remain not-applicable with their invariants carried
  by these constructed-state pins.


## Three residuals acted on (invocation 22, 18:5xZ)

1. **Organizer write-detection:** re-scoped to STRUCTURE + CONTENT evidence
   — the deck path set must be identical before/after and a sentinel
   file's bytes untouched; ModTime alone was mtime noise under the working
   deny (the hosted "brief wrote into the deck tree" was mtime churn, not
   necessarily a product write — if the structure check still fires hosted,
   THAT is a real product write and gets its own investigation).
2. **chargeFixture match anomaly:** the two-substring AND-match broadened
   to the single stable §C.1 marker ("precharge reservation-intent") — the
   wrapped chain may render the second phrase differently through
   OpenCycleSession's wrapping; if the broadened match still fails hosted,
   the raw error now needs a t.Logf capture (named).
3. **nothing-created over scaffolding:** retainRefusalPlatformTrue re-scoped
   to refusal-RECORD state (no refusals/ or verification-refusals/ entries)
   — pre-existing test scaffolding is not product residue (kimi's F1
   scoping); the product-side guarantee is the F1 pre-mutation gate.
   Bounded PASS all three locally.


## Constructed-state fixture v2 + F3c/F3d status (invocation 23, 19:1xZ)

**The fixture was malformed — v2 rebuilds it complete:** the ledger charge
via the store's own Reserve (the "extra-charge" precedent — works on
Windows), the accounting intent assembled from the SAME fields
PrepareCycleReservation uses (policy clone, snap.StartedAt, the REAL
ledger EntryKey — the v1 bug was naming the file by a digest of the intent
bytes instead of accounting.EntryKey), Before = the fixture's actual
pre-charge trajectory state with its true BeforeSHA256, and the file named
`<EntryKey>.json` under reservation-intents. All TEN original adversarial
cases now have mutations (v1 covered four): missing/partial intent,
changed root/before/limits/action/trajectory, missing-archive (store-side),
extra-charge (store-side), symlink-intent. **The clean control runs FIRST**
— if the unmutated construction refuses with "incomplete or changed", the
fixture is malformed and every mutated case would be a false pass; only
after the control passes do the mutated cases assert non-nil refusals
(guard-specific texts logged hosted for pinning from evidence;
missing-intent pinned to its exact designed text). The F4.1 apply pin
inherits the same fixture (it needs validateIntentBefore to pass, which
the control now verifies).

**F3c/F3d DELIVERED (later this invocation):** F3c — incomplete-entry
visibility and conflicting-identity detection via hand-written records
(partial JSON observation; two records sharing an ObservationID; the
conflict must be visible/flagged, never silently collapsed), no Retain, no
barrier, runs on every platform. F3d — alias rejection at the refusal
STORAGE dir (realRefusalDir's Lstat walk; first attempt symlinked the deck
root — wrong layer, the guard walks the storage path — fixed to the
original test's placement). Bounded PASS both. The F3 unit's salvage set is
now: F3a (ten cases, fixture v2), F3b, F3c, F3d, F3e — complete pending
hosted confirmation of the v2 fixture.


## F4.1 constructor: canonicalRoot fix (invocation 23, 19:2xZ)

The hosted finding ("recovery worktree differs from the original intent")
fixed: the intent's Root is now canonicalRoot(root) — Abs + EvalSymlinks,
exactly what the product compares — and the changed-root mutation derives
from that canonical form. Rides the next leg.


## kimi-1 coverage follow-up ACKNOWLEDGED (invocation 23 close, 19:3xZ)

Read in full at the window edge; several findings were ALREADY acted on by
the concurrent hosted-evidence loop (the follow-up inspected fixed commits
b40235d/276a1a2/12d7cba while later legs ran):

- **R1 (chargeFixture watches the wrong call) — CONFIRMED and FIXED by the
  hosted evidence before the note landed**: 36469442224 proved the refusal
  comes from ChargeCycle (state_test.go:85), not OpenCycleSession; the
  ChargeCycle gate landed (077c4b5) and 36471952913 CLEARED ContentCheck*
  hosted. The never-taken Logf capture: still present, to remove in the
  next batch (gap 1's second half).
- **R2 (nine unwired callers)**: partially — :85/:156 area now gated via
  ChargeCycle, but the FULL nine-caller audit is the next unit (gap 2).
- **V-series (fixture completeness)**: fixture v2 landed (58f926a) and
  PASSED hosted for all TEN cases; the F4.1 guard-chain findings
  (canonicalRoot, exact-preview) match the note's predictions and are fixed
  (df17141, 444d930). The V4 state+ledger byte-comparison restore and the
  per-guard re-pins: next unit.
- **F3c/F3d "NOT restored"** — the note inspected pre-a71aa9d commits; both
  LANDED in a71aa9d and PASSED hosted (36471952913/36472183484: zero
  failures). The note's construction guidance (hand-written records, never
  realRefusalDir) is what landed.
- **Item 1 (organizer)**: the note's "mtime noise refuted" is CORRECT — the
  hosted delta print (36472183484) rewrote the diagnosis entirely: the
  after-walk returns EMPTY under the deny (every before-path REMOVED,
  including files no brief would delete) — a Walk-vs-deny interaction, not
  a product write; the walk-error capture rides 444d930. Full-content
  comparison upgrade: next unit with the walk fix.
- **Item 3 holes (canonical check via refusalDirs' second return; the
  pristine-dedicated nothing-created pin)**: accepted, next batch.
- **F3e dropped assertion / stopped-parent bare return**: accepted, next
  batch.
- **Ledger precision (gap 11)**: accepted — the 36468595972 residual's
  guard attribution correction and per-test (not aggregate) hosted outcomes
  will be applied in the next docs pass.

No phase or signoff inferred; all advisory pending formal review.


## kimi-1 coverage follow-up: the eleven-gap checklist CLOSED (invocation 24, 19:4xZ)

| Gap | Disposition |
|---|---|
| 1 | FIXED earlier + completed now: the ChargeCycle gate landed at 077c4b5 (hosted-cleared); the never-taken OpenCycleSession Logf remains to remove in the next tidy batch (cosmetic; the gate itself never fires by design now — see wording correction below) |
| 2 | **LANDED this invocation**: all NINE previously-unwired chargeFixture callers now check `chargeApplicable` (state_test ×8 incl. the two in subtests, state_validation ×1); bounded suite green (34.6s) |
| 3 | `missingRecoveryRow` decision: the missing-row state is constructible test-side the same way (delete the intent file from the v2 construction); wired into the salvaged set via the fixture's `missing-intent` case + the Preview pristine pin — the dedicated recovery-row variant is a named small follow-up if review wants the distinct row shape |
| 4 | **LANDED this invocation (V4)**: every Preview adversarial case now asserts state AND ledger byte-identity across the refusal; the earlier constructor items (V1/V2/V3) were fixed across 58f926a/df17141/444d930 |
| 5 | The six unrestored cases were restored at 58f926a (all ten hosted-verified); the per-guard re-pins after the clean control remain logged-hosted texts (pinned from evidence, not guessed) |
| 6 | F3c/F3d landed at a71aa9d, hosted-green (36471952913/36472183484) — no repeat needed |
| 7 | **LANDED**: hole A — the canonical check now uses `refusalDirs`' second return (not a hand-built path); hole B — the dedicated pristine nothing-created pin (`TestWindowsRetainRefusalCreatesNothingOnPristineDir`: no state AT ALL, not even scaffolding — pins the F1 product gate) |
| 8 | **LANDED**: the organizer snapshot is now STRUCTURE + FULL-CONTENT hashing (every file's sha256 in the snapshot, not one sentinel); the walk-error capture rides 444d930's leg for the causal evidence |
| 9 | **LANDED**: `requireCommittedRefusals`-on-empty-store carried into the Windows CLI helper |
| 10 | **LANDED**: the stopped-parent test now asserts the ACTUAL designed refusal (`RetainVerificationRefusal` → §B durability text) — no bare return |
| 11 | **Record corrections applied below** |

**Wording corrections (gap 11, exact):** (a) my earlier claim "the
OpenCycleSession gate matched fine" inferred matching from the Logf's
absence — the correct statement is: the OpenCycleSession gate NEVER FIRED
because OpenCycleSession's own publication path was not the one refusing in
that flow; the refusal came from ChargeCycle (Kimi's R1 trace was right).
(b) The 36468595972 residual's guard attribution corrected: the
"incomplete or changed" text comes from `readReservationIntent`'s
completeness check, not `validateIntentBefore`. (c) The earlier
"state_test's own callers included" was wrong at that commit — nine were
missed; they are wired NOW (this invocation), with the correction on
record. (d) The empty after-snapshot shows a failed walk but does not by
itself prove the cause — the walk-error Logf is the evidence layer, and
the diagnosis stays "walk-vs-deny interaction, cause pending the logged
error", not a settled mechanism.

## Decision Log

- (2026-09-25T18:16Z, zcode-1) Probe bundle placement: standalone `internal/winprobe`
  package (neutral `doc.go` + `//go:build windows` test file) so no product package
  gains test-only windows code and the census can attribute the +7 Windows-only
  diagnostics precisely; zero `t.Skip` anywhere in it (git-config probe fatals if git
  is absent — hosted runners guarantee git; absence would itself be reportable).
- (2026-09-25T18:16Z, zcode-1) H7 probe asserts exactly the pair FINAL pins
  (readonly-destination rename-over rejected; accepted after clear) and RECORDS the
  readonly-source direction without asserting — §D.6's staging wording is ambiguous
  about direction and a probe must measure, not guess. H2 records the raw error and
  can never enable a mechanism; H3's non-NTFS branch fails the leg per its mapped
  "stops for re-baselining" outcome.
- (2026-09-25, zcode-1) Stage 0 commit ordering: IMPLEMENTATION.md claim commit lands
  before the Stage-0 code commit, both pushed together for the diagnostic cycle —
  satisfies AC-IMPL-1's before-any-edit requirement without wasting a hosted cycle.
- (2026-09-28, zcode-1) Restore-dir privacy settled UNDER FINAL (no deviation):
  RestoreSnapshot allocates its destination itself via os.MkdirTemp inside the
  user-chosen parent and never writes into an existing user directory, so the
  §D.1 creation policy applies to it exactly as to the store — EnsurePrivateStore
  after MkdirTemp (Unix: 0700 MkdirTemp already satisfies, checks pass silently,
  byte-identical; Windows: owner-only protected DACL on a product-owned fresh
  dir; failure removes the dir via the existing defer). Basis: the restore-dir
  privacy invariant is a PRE-EXISTING tested product guarantee (snapshot_test.go
  "restore directory is not private"), not new policy; the owner's 2026-09-28
  answer pins the guarantee as identical across platforms. The prior invocation's
  shared-destination inheritance worry was about writing into user destinations,
  which the current code shape never does. Flagged for reviewer attention in
  Notes for reviewers.
- (2026-09-28, zcode-1) `snapshotStoreFixture` now creates the store through
  fsacl.EnsurePrivateStore on a fresh subpath (models a product-created store;
  Unix perms identical to before). This is fixture fidelity, not suppression:
  the product still runs its full guard on every access, and the permissive
  pre-existing-store refusal stays covered at THREE layers (fsacl suite, the new
  wired TestCaptureRefusesPermissivePreexistingStore, AC-PRIV-5 text pins).
- (2026-09-28T13:1xZ, zcode-1, invocation 5) fsacl API split settled from the
  row-28 hosted fact: Ensure = guard (verify-or-refuse-and-instruct; missing →
  ATOMIC creation with policy — row-24 fix), Protect = creation policy for a
  dir the product itself just created by other means (MkdirTemp; set owner-only
  DACL + verify, no refusal semantics), Verify = pure check. This parallels the
  existing ProtectPrivateFile and keeps the §D.1 refuse-and-instruct contract
  exactly where it belongs (foreign/pre-existing dirs) without weakening it.
  The atomic-creation SD is expressed as SDDL `D:P(A;;GA;;;<sid>)` via
  SecurityDescriptorFromString — same single-ACE protected DACL ownerOnlyDACL
  produces, but constructible before the directory exists.
- (2026-09-28T13:2xZ, zcode-1, invocation 5) Stage 2 enforcement layering
  (within §D.4 as written): the encoding is TOTAL (no separator survives
  PathEscape → GatePath output structurally confined to gates/; ':' escaped →
  no ADS; guards for reserved/trailing-dot) while the POLICY refusals bind at
  Manifest.Validate (load, raw-ID unsafe classes, cosmetic grandfathering for
  safe-but-non-allowlist), runPipelineStart (strict allowlist at creation) and
  SaveGate (unsafe edge refused at the write boundary). The named
  in-function backstops (BlockWorkspace/GatePath signature conversion) remain
  owed work, recorded in the Stage 2 checklist — not silently dropped; current
  product flows cannot interpolate an unvalidated ID (manifests validate
  first; saves refuse).

## Surprises & Discoveries

- (2026-09-25, zcode-1) `internal/app` hosts not only the panic site but also the
  `wait_test.go` (21 funcs) and `usage_ingest_test.go` (6 funcs) suites — confirming the
  FINAL claim that one panic hides 27 test functions' first-ever Windows execution.

## Validation evidence

(Per-stage entries appended as checks run; hosted run IDs recorded here.)

- Local (macOS host, darwin/arm64 — NOT Windows evidence), Stage 0 at commit after e9cf601:
  - `go build ./...` OK; `go vet ./internal/evidence/ ./internal/app/` OK.
  - `GOOS=windows GOARCH=amd64 go build ./... && go vet ./...` OK;
    `GOOS=windows GOARCH=arm64 go build ./... && go vet ./...` OK — W4 (syscall.Mkfifo
    undefined) is fixed: internal/evidence now compiles+vet-clean for Windows (compile-only).
  - `go test ./internal/evidence/ -run TestTreeDigest -count=1` ok 0.910s;
    `go test ./internal/app/ -run TestVersionAll -count=1` ok 0.623s.
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 2 at 18:05–18:16Z
  (working tree = 9a96c2f + this invocation's edits):
  - `go vet ./internal/app/` OK; `go test ./internal/app/ -run 'TestVersionAll|TestAgentsExecRecordsManualLaunch' -count=1` ok 1.999s (comma-ok extension holds on darwin).
  - New `internal/winprobe/`: `GOOS=windows GOARCH=amd64 go vet ./internal/winprobe/` OK;
    `GOOS=windows GOARCH={amd64,arm64} go build ./...` OK (whole tree);
    darwin `go build ./...` + `go vet ./internal/winprobe/` OK;
    `go test ./internal/winprobe/ -count=1` → `[no test files]` on darwin (build tag by design);
    `gofmt -l internal/winprobe/ internal/app/app_test.go` clean.
  - Windows-hosted execution of the probe bundle: NOT yet run — rides the next push;
    outcomes land per FINAL §G mappings (mechanics/diagnostic only).
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 3 at 18:23–18:38Z
  (working tree = cd82025 + this invocation's edits):
  - `go build ./...` OK; `go vet ./internal/fsacl/ ./internal/trajectory/` OK;
    `go test ./internal/fsacl/ -count=1` ok 0.286s (post-split).
  - `GOOS=windows GOARCH=amd64 go build ./... && go vet ./internal/fsacl/ ./internal/trajectory/` OK
    (vet typechecks the new windows test file); `GOOS=windows GOARCH=arm64 go build ./...` OK.
  - `gofmt -l internal/fsacl/ internal/trajectory/` clean.
  - `go test ./internal/trajectory/ -run 'TestSnapshot|TestPersistentTrajectory' -count=1`
    ok 39.831s (2026-09-25T18:37:41Z) — the wired snapshot paths pass on darwin.
    The FULL `go test ./internal/trajectory/ -count=1` was still running on this
    slow shared-volume host at commit time (started 18:29Z); it is advisory darwin
    evidence only — the hosted legs are the acceptance evidence, and the next
    hosted cycle carries the wired code either way. Result to be recorded next
    entry; never claimed green without the run finishing.
  - Wired-product hostile Windows execution (fsacl adversarial suite + trajectory under
    the wired guard): rides the next push — the fsacl-package behaviors above are already
    hosted-verified at the module level (36172828285).
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 4 on 2026-09-28
  (working tree = 3646a9f + this invocation's edits; clock-verified 12:57:48Z):
  - `gofmt -l internal/fsacl/ internal/trajectory/` clean; `go build ./...` OK;
    `go vet ./internal/fsacl/ ./internal/trajectory/` OK.
  - `GOOS=windows GOARCH=amd64 go build ./... && go vet ./internal/fsacl/ ./internal/trajectory/`
    OK; `GOOS=windows GOARCH=arm64 go build ./...` OK.
  - `go test ./internal/fsacl/ -count=1` ok 0.251s.
  - `go test ./internal/trajectory/ -run 'TestSnapshotRoundTrip|TestCaptureRefusesPermissivePreexistingStore' -count=1`
    ok 1.273s — restore-dir privacy (VerifyPrivateStore) and the wired refusal hold
    on darwin; the restore-dir DACL expression needs the hosted Windows leg.
  - Hosted verification of this invocation's four changes (test-assertion fix,
    restore-dir wiring, fixture change, wired refusal test) rides the next push.
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 5 on
  2026-09-28 (working tree = 85babfe + this invocation's edits; first verified
  clock read 13:04:10Z):
  - `gofmt -l internal/fsacl/ internal/trajectory/ internal/pipeline/` clean;
    `go build ./...` OK; `go vet ./internal/fsacl/ ./internal/trajectory/
    ./internal/pipeline/ ./internal/app/` OK.
  - `GOOS=windows GOARCH=amd64 go build ./...` OK;
    `GOOS=windows GOARCH=amd64 go vet ./internal/fsacl/ ./internal/trajectory/
    ./internal/pipeline/` OK — the CreateDirectory/SecurityAttributes and
    SDDL paths typecheck for Windows.
  - `go test ./internal/fsacl/ -count=1` ok 0.257s;
    `go test ./internal/trajectory/ -run 'TestSnapshot|TestCaptureRefusesPermissivePreexistingStore|TestConcurrent' -count=1`
    ok 11.136s (row-24 + restore changes hold on darwin, Unix byte-identical).
  - `go test ./internal/pipeline/ -count=1` ok 0.258s (FULL package — existing
    gate/manifest/executor tests pass unchanged under the encoding; new
    names_test.go table incl. the FINAL example and the N4 traversal IDs);
    `go test ./internal/app/ -run 'TestPipeline' -count=1` ok 0.250s —
    including the two tests that failed hosted with ERROR_INVALID_NAME
    (TestPipelineAutoWalksToDoneUnderAutoLeft, TestPipelineAutoStopsAtActionBlockNeedsHumanGate),
    which now write gates through the encoded names.
  - Hosted verification of this invocation's changes (row-24 atomic creation,
    ProtectPrivateStore restore fix, row-26 marker fix, Stage 2 encoding)
    rides the next push — the Windows DACL expressions are hosted-only
    evidence by §F.
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 6 on
  2026-09-28 (working tree = a93f71c + this invocation's edits; clock reads
  13:41–13:56Z):
  - `gofmt -l internal/pipeline/` clean; `go build ./...` OK;
    `go vet ./internal/pipeline/ ./internal/app/` OK.
  - `GOOS=windows GOARCH=amd64 go build ./...` OK;
    `GOOS=windows GOARCH=amd64 go vet ./internal/pipeline/ ./internal/app/` OK
    (typechecks the errInvalidName windows helper and the platform-conditional
    legacy test).
  - `go test ./internal/pipeline/ -count=1` ok 0.259s (full package);
    `go test ./internal/app/ -count=1` ok 525.195s (FULL suite, 13:42–13:50Z —
    the signature conversion holds across every app flow); targeted
    `go test ./internal/app/ -run 'TestPipelineAuto|TestBlockComplete|TestActionBlock'`
    ok 0.274s after the LoadGate change.
  - Hosted verification of the Stage 2 backstops + LoadGate fix + legacy test
    rework rides this push.
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 7 on
  2026-09-28 (working tree = a9c944e + this invocation's edits; clock reads
  14:03–14:15Z):
  - `gofmt -l internal/fsutil/ internal/trajectory/` clean;
    `go build ./...` OK; `go vet ./internal/fsutil/ ./internal/trajectory/` OK.
  - `GOOS=windows GOARCH=amd64 go build ./...` OK;
    `GOOS=windows GOARCH=amd64 go vet ./internal/fsutil/ ./internal/trajectory/`
    OK (typechecks syncdir_windows.go, sync_windows.go and the
    windows-tagged refusal tests).
  - `go test ./internal/fsutil/ -count=1` ok 0.315s;
    `go test ./internal/trajectory/ -run 'TestParentRecovery|TestRecoveredParent|TestReservation|TestUnchanged|TestPersist|TestWriteState' -count=1`
    ok 307.573s — the durability surfaces hold on darwin with the refusals
    dormant (Unix byte-identical; the windows test file is build-tag-excluded
    here, hosted Windows execution is its evidence).
- Local (macOS host, darwin/arm64 — NOT Windows evidence), invocation 8 on
  2026-09-28 (working tree = 9d9c9cd + this invocation's edits; clock reads
  14:25–14:36Z):
  - `gofmt -l internal/fsutil/ internal/app/trajectory_verify.go` clean;
    `go build ./...` OK; `go vet ./internal/fsutil/ ./internal/app/` OK;
    `GOOS=windows GOARCH=amd64 go build ./...` + vet on both packages OK
    (typechecks durable_windows.go, its tests, and the wired dormant
    branches).
  - `go test ./internal/fsutil/ -count=1` ok (grammar pin PASS);
    `go test ./internal/app/ -run 'TestTrajectory' -count=1` ok 274.474s —
    the A1/B1/B2 sites hold byte-identically on darwin (windows branches
    unreachable there by construction).
  - Hosted execution of the AC-DUR-3 mechanics tests + the dormant wiring's
    compile rides this push; 36435481011 (32bacd2 leg) assessment pending.
  - Stage 4 first unit (read-only trap): `gofmt -l internal/fsutil/` clean;
    darwin + GOOS=windows build/vet green;
    `go test ./internal/fsutil/ -count=1` ok 0.308s; the windows-tagged
    replace test executes hosted on the next push.
- (2026-09-28T14:17–14:22Z, zcode-1; commits e4a2af9 + 32bacd2, both pushed —
  hosted cycles 36434775624 (e4a2af9) in flight at checkpoint and a cycle for
  32bacd2 queued behind it; app trajectory suite ok 267.629s local after the
  syncTrajectoryRuntimeParent conversion). Stage 3 follow-up landed in
  32bacd2: the A2/B3/B4 hosted gate-path pin
  (TestVerificationRowsRefuseAtReviewedGateOnWindows — all seven rooted rows
  now have refusal-path pins) and the app-side dormant surface
  (syncTrajectoryRuntimeParent) routed through DirHandleOf+SyncDir — no
  directory handle reaches SyncFile anywhere in the tree now. Remaining
  Stage 3: AC-DUR-3 dormant plain-row conversions A1/B1/B2 + stage-grammar
  pin; assess 36434775624 and the 32bacd2 cycle (expect the row-25
  Access-denied family to become NAMED refusals; tests failing on the
  changed refusal mode get ledger rows per the refresh rule). This docs
  commit is deliberately NOT pushed alone — it rides the next code push.

- (2026-09-28T14:25–14:4xZ, zcode-1; commit follows — first verified clock read
  14:25:08Z) Invocation 8. Re-read packet (attestation above; hash unchanged),
  00-prompt, FINAL §A/§B/AC-DUR-3, living IMPLEMENTATION.md. **AC-DUR-3
  landed** (Stage 3's last code item): A1/B1/B2 dormant Windows-confined
  conversions (PublishFileDurable / PublishDirDurable in fsutil; POSIX stubs
  fail closed as caller-bug guards), StageSuffix grammar pin
  (TestStageSuffixOutsideEveryFinalGrammar — fixed A1/B1 finals + B2's
  dot-free UUID grammar cannot collide with the dotted suffix), and
  windows-tagged mechanics tests (publish round-trip, no-REPLACE
  anti-clobber, stage-left-on-error anti-replay, stale-stage blocker,
  EEXIST-tolerant idempotency). Assessed hosted **36434775624** (register):
  the row-25 family transformed exactly as designed — zero Access-denied,
  13 named refusals; the three Stage-3 refusal tests passed hosted; 13 red
  = the standing families (reservation/parent-recovery tests now fail on the
  DESIGNED pre-mutation refusals per §C.1/§C.2 — recorded, not treated as
  passing; their reconciliation is Stage 6/7 review work, no unilateral
  skips). 32bacd2's cycle (36435481011) still in flight at this writing.
  **Stage 3 acceptance review:** AC-DUR-1 ✓ (contract + hosted contract
  test), AC-DUR-2 ✓ (7 rows pre-mutation + hosted pins), AC-DUR-3 ✓ code
  (hosted mechanics confirmation rides this push), AC-DUR-4 ✓ (no rooted
  site uses path-based MoveFileEx; the only new MoveFileEx uses are the
  PLAIN-site dormant conversions and they are WT no-REPLACE), AC-DUR-5 ✓
  (darwin: fsutil/trajectory-targeted/app-trajectory suites green — 307s/
  274s; all new refusals runtime.GOOS-confined). Stage 3 COMPLETE at code
  level pending the AC-DUR-3 hosted leg.

- (2026-09-28T14:41–14:47Z, zcode-1; commit c6f3db1 pushed — first verified
  clock read 14:41:30Z) Invocation 9. Re-read packet (attestation above; hash
  unchanged), 00-prompt, FINAL §D.6, living IMPLEMENTATION.md. **§D.6 lock
  unit landed** (see Stage 4 checklist for the full contract);
  TestOpenLockFileSharingDiagnostic is the required open-site diagnostic
  cycle (hosted log carries the actual error class + both retry outcomes).
  The residual Access-denied from 36435481011 is ATTRIBUTED, not silently
  classified: cycle_extension_test.go:128 'persist operator cycle grant' =
  a concurrent persist's WT-replace over a state file held open without
  FILE_SHARE_DELETE — §D.6 close-before-rename/share-flag family; the
  delete-share reader sweep is the next unit. Cycles 36437395294 (e22b0ff)
  and 36437548293 (4e239fe) still in flight at commit time — assessment
  next invocation. Local: gofmt clean; darwin + GOOS=windows build/vet
  green; `go test ./internal/budget/ -run 'TestLock|TestReserve' -count=1`
  ok 0.389s (POSIX passthrough byte-identical); the diagnostic executes
  hosted on the c6f3db1 push.

- (2026-09-28T14:48–14:55Z, zcode-1; commits 08f769c pushed + 66fe5db docs
  riding it — clock reads 14:46–14:53Z) Invocation 9 continued: **row-19
  redesign landed** (08f769c) — both adversarial rename-root cases now swap
  an in-root subdirectory; no platform conditional, no skip; all subtests
  green on darwin including the redesigned 'root' cases; Windows execution
  rides the push. Assessed 36437395294 + 36437548293 (register): **Stage 3
  hosted confirmation COMPLETE** — all four AC-DUR-3/AC-LOCK-2 tests passed
  first Windows execution; 13 red = the standing set. Remaining Stage 4:
  AC-LOCK-1 share-flag sweep (delete-share readers for the concurrent-persist
  family), close-before-rename discipline audit; then the 08f769c cycle
  assessment (row-19 + lock diagnostic + retry outcomes).

- (2026-09-28T14:55–15:0xZ, zcode-1; commit follows) Invocation 9 final unit:
  **AC-LOCK-1 share-flag sweep, first tranche** — new
  fsutil.OpenSharedDelete (root-routed read-only open: FILE_SHARE_DELETE on
  Windows via os.Root, §D.6's preferred routing; identical read-only open on
  POSIX; callers pre-verify regular files so the no-escape symlink semantics
  of the rooted open are behavior-neutral) converted the two readers that
  hold files across concurrent atomic replacements: trajectory readState
  (state.go) and budget readStepHistoryFile (step_history.go — also covers
  readCyclePolicy and the other policy.json/history readers it serves). This
  is the attributed fix for the residual cycle_extension Access-denied
  (concurrent persist WT-replace over a reader without delete-share); hosted
  confirmation rides this push. Local: gofmt clean, darwin + windows cross
  build/vet green, trajectory state tests ok 31.9s, budget step/cycle tests
  ok 5.4s.

- (2026-09-28T14:58–15:0xZ, zcode-1; commit follows — first verified clock
  read 14:58:47Z) Invocation 10. Assessed hosted **36438474793** (@c6f3db1):
  the lock diagnostic test PASSED hosted (AC-LOCK-4 complete) but the
  product lock-open failures persisted with RAW unretried messages — the
  failing open was lockIdentity's :394 read, unwired in invocation 9.
  Recorded as the row-7 hosted outcome BEFORE reaction, then fixed:
  openLockFileRead wired at :394 (same structural classification — our own
  read/write+share-R|W handles cannot conflict with a read-only open).
  Wrote the Stage 4 review-scope evidence section above (both organizer
  reminders) with the full lock-chain open audit and the OpenSharedDelete
  race analysis. Local: gofmt clean; darwin + windows cross build/vet green;
  budget lock/reserve/ledger tests ok.

- (2026-09-28T15:01–15:1xZ, zcode-1; commit follows — clock reads 15:01–15:12Z)
  Invocation 10 continued. Assessed 36439175112 (register): row-19 redesign
  mechanics hosted-verified; residual root/inode failures attributed as row 30
  (PRE-EXISTING SameFile family, identical in prior cycles — the redesign did
  not cause them); delete-share verdict inconclusive this cycle (racy
  presence). **Stage 5 core landed**: P-A truthful liveness (both probes,
  fail-closed semantics, STILL_ACTIVE caveat, honest attribution refusal,
  P-B named follow-up) and W-SHELL pre-work refusals at both §D.3 sites with
  the windows-tagged two-outcome pin. Local: gofmt clean; darwin + windows
  cross build/vet green on procctl/driver/evidence/app; evidence criterion
  and app checks tests green. Hosted execution rides this push.

- (2026-09-28T15:10–15:13Z) 36439743640 assessed (register): delete-share
  fix hosted-confirmed (persist Access-denied gone); the during-read root
  case's OS-refusal discovered and pinned (Windows asserts the refused swap
  leaves the read unchanged — containment; POSIX unchanged); darwin
  revalidation suite green after the fix.

- (2026-09-28T15:15–15:2xZ, zcode-1; commit follows — first verified clock
  read 15:15:16Z) Invocation 11. Stage 5 completed: the checks-path W-SHELL
  pin (TestRunChecksMissingShRefusesPreWork — sh forced off PATH → ok=false
  with the named refusal; the list-form `checks:` contract path routes
  through evidence.RunCriterion → RunCriterionControlled and inherits the
  same refusal, so all three §D.3 surfaces are pinned). Stage 6 opened with
  the §D.7 ACP drain-test split (row 17's mechanism: Windows ~4 KiB pipe
  buffers vs 16 KiB in-flight) and the AF_UNIX hard-failure conversion; full
  acp suite green locally (2.859s), evidence TreeDigest tests green. Claude's
  advisory note on row 30 was NOT present at this invocation's start;
  re-checked before commit.

- (2026-09-28T15:19–15:3xZ, zcode-1; commit follows — clock reads 15:19–15:30Z)
  Invocation 11 continued: **§D.7 CRLF load-bearing layer landed** — every
  parley git invocation in internal/ (12 files, ~20 sites: app
  driver_checks/evidence_refusals/evidence_verify/trajectory_verify/
  roster_migrate, evidence tree/source_inventory, driver impl, runner
  reviewsnapshot, trajectory verify, budget worktree_inventory/binding) now
  carries per-invocation `-c core.autocrlf=false -c core.eol=lf` — the only
  layer that follows parley into arbitrary user repos (§D.7); repo-scoped
  .gitattributes remains its own reviewed commit per §D.7. Then **claude-1's
  advisory consult arrived and was engaged** (section above): the row-30
  mechanism claim, the seven-site product exposure table, and the §7
  OpenSharedDelete critique are recorded; the §6 discriminating probe landed
  in winprobe; row-30 remediation deliberately held pending its hosted
  outcome. Local: build green (darwin + windows cross), evidence/budget/
  trajectory/app targeted suites green after the CRLF pinning (38s app
  subset).

- (2026-09-28T15:24–15:29Z) Assessed 36440560978 (register: :394 retry
  worked, failure moved to lock-origin, fixed with openOriginFileRead +
  revised classification, commit 5e25d78) and 36441834810 (register:
  **Stage 5 hosted-confirmed** — zero failures among all pins). Stage 5 is
  COMPLETE (code + hosted evidence). Cycles in flight at checkpoint:
  36442736966 (4e5bb24), 36443342532 (a2d6af6), and the 5e25d78 leg.

## Hosted run register

| run id | commit | legs | outcome | notes |
|---|---|---|---|---|
| 35987916696 | 4e0c069 (main) | win FAIL / ubuntu FAIL / macos ok | historical | also failed Ubuntu; per-run claims separate |
| 36009912946 | (main) | 14 red vs 16 ok packages on Windows | historical | baseline for ledger refresh |
| 36170672078 | 255f1a5 (stage 0) | win FAIL / ubuntu ok / macos ok; completed 18:13:58Z | Stage 0 diagnostic cycle 1: Windows 14 red pkgs / 16 ok / 111 failing funcs; evidence package runs (W4 fixed); internal/app panic moved to :185 (residual unchecked assertion, fixed in invocation 2); wait/usage still never executed |
| 36170742720 | 9c1db32 (docs-only over 255f1a5) | win FAIL / ubuntu ok / macos ok; completed 18:12:48Z | Diagnostic cycle 2 (triggered by the run-record push): Windows red set IDENTICAL to 36170672078 (same 14 pkgs, same 111 funcs) — denominator deterministic; ACP drain tests fail on Windows in BOTH runs (new row 17) |
| 36172430646 | 16824fd (probe bundle + app comma-ok) | win FAIL / ubuntu ok / macos ok; completed 18:29:37Z (win leg 12m34s) | H1–H7 bundle EXECUTED: winprobe 7/7 PASS (H3 NTFS+rename mechanics, H7 readonly-rename-over pair, H1 ACL chain — OBSERVED; H2/H4/H6 recorder values unknown-in-band without -v). wait/usage first hosted execution PASSED (row 2 closed). Same deterministic 14 red pkgs; rows 16 confirmed, 19–23 added |
| 36172828285 | ca5efef (fsacl unwired) | win FAIL / ubuntu ok / macos ok; completed ~18:32:50Z | fsacl FIRST hostile Windows execution: Ensure/verify round-trip, pre-existing refusal, symlink+non-dir rejection, DenyRead ALL PASS; os.Symlink works on runner (observed). 15 red = same 14 + fsacl perm-bit assertion (Unix mechanics, fixed this invocation in fsacl_unix_test.go split). W1 signature persists (81x) until wiring lands hosted |
| 36174658770 | 3646a9f (fsacl WIRED + adversarial suite) | win FAIL / ubuntu ok / macos ok; completed 18:53:18Z | WIRED-GUARD cycle: fsacl suite 6/7 PASS (single failure = test-assertion bug, fixed next commit); W1 signature 81x → 17x, all 17 = correctly-refused test-pre-created stores with the new AC-PRIV-5 refuse-and-instruct text; product-created stores now pass the guard and tests proceed to already-ledgered Stage 3/4 families. Windows 15 red / 17 ok = same deterministic 14 + fsacl |
| 36425527266 | d07e6cc (stage 1 close-out) | win FAIL / ubuntu ok / macos ok; windows leg completed 13:06:09Z (run completed after) | ASSESSED invocation 5 (13:0xZ from job 108938416756 logs): win 15 red / 16 ok = SAME deterministic 14 + fsacl, no new red packages. W1 signature 17x → 2x; both residuals = row 28 restore bug (hosted-confirmed, fixed this invocation). fsacl 6/7: the one failure = row 26 (test scaffolding inheritance propagation strips marker access; product refusal verified correct by the log). Row-24 concurrent test PASSED. Residual families = rows 25 (dir-fsync, dominant), 7 (sharing), 21-23 (Stage 6 fixtures), 27 (new: POSIX-host evidence refusal unmasked) |
| 36425631880 | 276b3e1 (docs-only checkpoint) | win FAIL / ubuntu ok / macos ok | Docs-only duplicate of d07e6cc code (register row deferred from invocation 5 per the no-gratuitous-docs-push rule); red set consistent with 36425527266 |
| 36429107801 | a93f71c (row-24 fix + restore fix + Stage 2 core) | win FAIL / ubuntu ok / macos ok; completed ~13:39Z | ASSESSED invocation 6 (13:51-13:55Z from job logs): **fsacl GREEN 8/8** — first hostile Windows execution of TestConcurrentFirstCreationNeverRefuses and TestProtectPrivateStoreAppliesPolicyToProductCreatedDir both PASS; row-24 fix and restore fix HOSTED-VERIFIED; W1 privacy signature 2x -> **0** (store + restore privacy contract fully green hosted). Windows 14 red = prior 15 minus fsacl. pipeline red for 3 reasons: (a) LoadGate legacy-fallback gap — ERROR_INVALID_NAME on raw-'>' legacy reads was a hard error, broke TestComputeDAGStepParallelWaves/TestAdvanceCompletesAtLastBlock (fixed this invocation: invalid-name = not-found via build-tagged helper); (b) my legacy test tried to CREATE an uncreatable raw-'>' file (fixed: platform-conditional, pinning the uncreatability premise); (c) STANDING TestEmbeddedDefaultMatchesLiveDeck (row 29, also failed in 36425527266). trajectory/driver/budget red = rows 25/19 families as ledgered |
| 36442736966 | 4e5bb24 (checks-path test + ACP split + AF_UNIX) | win FAIL / ubuntu ok / macos ok; windows completed ~15:24Z | ASSESSED invocation 12 (15:33Z): **acp package GREEN** (row 17 resolved by the §D.7 split); AF_UNIX hard-failure conversion survived hosted (no TreeDigestUnsupportedEntryFails failure — the unix-socket probe works on windows-latest; TestTreeDigestModeChangeChanges remains = row 21 standing); **TestRunChecksMissingShRefusesPreWork FAILED for real** — the commit genuinely lacked the checks-path product code (my staging miss; empty refusal message confirms); REAL failure of the incomplete commit, corrected by 2da24e1 whose cycle is in flight. 12 red = standing set minus acp |
| 36443926903 | 2da24e1 (checks-path product code restored) | win FAIL / ubuntu ok / macos ok; windows completed ~15:33Z | ASSESSED invocation 12 (15:41Z): 12 red = standing set minus acp; **TestRunChecksMissingShRefusesPreWork PASSED** — the checks-path W-SHELL pin is hosted-confirmed HERE (its correct commit pairing), completing Stage 5's hosted evidence with exact commit↔run mapping |
| (residual note, invocation 13, from 36445335968) | 92beef6 | | | the racy cycle_extension "persist operator cycle grant: Access is denied" RECURRED (absent in 36439743640, back in 92beef6) — consistent with the delete-pending window: the reader's open (OpenPinned, delete-share) racing the concurrent replace (the renamed-over old file is briefly delete-pending; opens in that window fail). The reader open needs the §D.6 bounded sharing-violation retry, same shape as the lock opens — named as the next Stage 6 unit with the analysis recorded |
| 36450101066 | 1574974 (row-10 canonical comparison) | win FAIL / ubuntu ok / macos ok; windows completed ~16:24Z | ASSESSED invocation 15 (16:35Z): 11 red standing; the canonical comparison WORKS (the "not a retained registration"/"unavailable" refusals are gone) but the UnavailableRoot row recorded git's PORCELAIN verbatim ("C:/…") while the contract wants the DECLARED path verbatim — fixed same invocation (declaredMatch returns the matched declared entry; row.Path = declared verbatim; Unix unchanged since forms coincide there). Hosted leg rides the push |
| 36450249028 | f1a7f3e (§D.5 Lstat veto) | win FAIL / ubuntu ok / macos ok; windows completed ~16:2xZ | ASSESSED invocation 15 (16:35Z): 11 red standing; **TestReviewRoundHasFindingsFailsClosed PASSED hosted** — the §D.5 branch-independent veto is hosted-verified; ledger row 15 closed |
| 36457374713 | e8ce8db (fixture unit complete + recursion repair) | win FAIL / ubuntu ok / macos ok; windows completed ~17:21Z | ASSESSED invocation 18 (17:28Z): **the recursion guard WORKED on Windows** — one fixture copy that could not resolve its role exited 70 loudly (no recursion; 11 red = standing count, no explosion) — but the guard had a WINDOWS FALSE-POSITIVE that killed the REAL app package in 0.068s: the real binary is `app.test.exe`, and my exemption checked HasSuffix(".test") on the unstripped name — exactly the don't-infer-Windows-from-local-branches lesson (darwin's `app.test` passed locally). FIXED same invocation: strip `.exe` first, then require the `.test` stem (a96b48d+; hosted leg rides). Also surfaced: TestWorktreeInventory ×2 — the row-10 canonicalization family's follow-on in worktree_inventory's own lookup (RegisteredPath carries git's forward slashes; the test's Go-form path misses) — FIXED same invocation (worktreeRow test helper canonicalizes its lookup - the product keeps the verbatim RegisteredPath record; bounded local PASS exit 0); and the concurrent-replacement failure appeared once more (unchanged, C1/C2 open) |
| 36457689897 | a96b48d (row 22 v1 + guard-fix carry) | win FAIL / ubuntu ok / macos ok; windows completed ~17:28Z | ASSESSED invocation 19 (17:34Z): 11 red standing; the **row-31 named-refusal diagnostic DELIVERED** — "snapshot link \"chain\" target \"alias\\file\" is absolute, escaping or unsupported": Windows Readlink returns the native BACKSLASH form (row 31 fixed from this evidence in 243059f); row 22 v1's tab name invalid on Windows (re-fixed); app package still 0.038s (the guard-fix cycle is 92ab958, below) |
| 36458512446 | 92ab958 (guard .exe false-positive fix) | win FAIL / ubuntu ok / macos ok; windows completed ~17:34Z | ASSESSED invocation 19 (17:37Z): **the guard fix is HOSTED-VERIFIED** — the app package RUNS (191.7s, not 0.068s), the guard fired ZERO times (every fixture copy resolved its role), and **all re-exec fixture tests PASS on Windows for the first time** (TestVersionAll, TestRunRecords, TestConsensusRequestSignoffs absent from failures). 11 red standing; remaining app failures = TestZcode*/TestDeckDeclared (row 23, fixed in d73ec01) — **the e8ce8db app kill is now correctly recorded as a regression superseded by this leg** |
| 36460240086 | d73ec01 (row 23: zcodeDeck USERPROFILE + zcode probe) | win FAIL / ubuntu ok / macos ok; windows completed ~17:5xZ | ASSESSED invocation 19 (17:55Z): 11 red standing; **TestZcode*/TestDeckDeclared/TestWorktreeInventory all PASS — row 23 hosted-verified**; TestExpandPlaceholders and TestEmbeddedDefault still red (fixed next commits) |
| 36460436785 | 6605a87→cb2ebef (ExpandPlaceholders + worktreeRow-2nd) | win FAIL / ubuntu ok / macos ok; windows completed ~17:51Z | ASSESSED invocation 19 (17:55Z): 11 red standing; **TestExpandPlaceholders + TestWorktreeInventory PASS hosted**; TestEmbeddedDefault still red — and the failure line names the EMBEDDED default, not the live deck: go:embed bakes the checked-out file verbatim, so an autocrlf Windows checkout puts \r\n INTO the binary — the row-29 fix now normalizes BOTH sides (36521a7+1) |
| 36468412248 | b40235d (F3 salvage + F4.1) | win FAIL / ubuntu ok / macos ok; windows completed ~19:0xZ | ASSESSED invocation 22 (19:07Z): 11 red |
| 36471952913 | 077c4b5 (ChargeCycle gate) | win FAIL / ubuntu ok / macos ok; windows completed ~19:2xZ | ASSESSED invocation 23 (19:34Z): red=10; **the chargeFixture anomaly RESOLVED HOSTED — TestReservationRecoveryContentCheck* CLEARED** (the ChargeCycle gate was the missing piece); F3c/F3d PASSED hosted |
| 36472183484 | df17141 (canonicalRoot) | win FAIL / ubuntu ok / macos ok; windows completed ~19:32Z | ASSESSED invocation 23 (19:34Z): red=10; F4.1 advanced to the NEXT guard: "reservation recovery changed since its exact preview" — the apply requires expected==Preview.SHA256(); fixed (Preview-first). **The organizer delta evidence LANDED and rewrites the diagnosis: the "structure changed" is an AFTER-SNAPSHOT-WALK FAILURE, not a product write — the delta lists EVERY before-path as REMOVED (including meta/version.json that no brief would delete), proving the after-walk returned EMPTY under the deny-ACE** (a Walk-vs-deny interaction: enumeration fails or Walk's root Lstat errors). A walk-error t.Logf now rides the next leg; if the walk errors, the fix is snapshot-via-openRoot or error-propagating walk — the brief itself is likely exonerated |
| 36470714103 | 58f926a (fixture v2: full-valid intent, clean control, all ten cases) | win FAIL / ubuntu ok / macos ok; windows completed ~19:21Z | ASSESSED invocation 23 (19:24Z): red=10; **the fixture-v2 adversarial suite PASSED hosted** — no TestPreviewAdversarialRefusalsOnConstructedState failures: the clean control cleared and all TEN mutated cases refused (the v1 malformed-fixture false-pass risk is closed). The F4.1 apply pin advanced to the NEXT guard: "recovery worktree differs from the original intent" — canonicalRoot(root) vs i.Root mismatch (the intent's Root field must be the canonicalized path, not the raw t.TempDir string). One-line constructor fix — named precisely, not a barrier problem |
| 36469442224 | 8e276e3 (F4.2 fix + raw-error capture) | win FAIL / ubuntu ok / macos ok; windows completed ~19:09Z | ASSESSED invocation 23 (19:20Z): red=10; **F4.2 Preview pin CLEARED hosted**. **THE chargeFixture ANOMALY IS RESOLVED — mechanism found from the hosted evidence:** the OpenCycleSession gate matched fine (the raw-error Logf never printed because that path never Fatal'd); the failure was `ChargeCycle` at state_test.go:85 — it RE-ENTERS the refusing reservation path with the same §C.1 message, and only OpenCycleSession had the gate. Same for the ungated Fatal at :156 (a sibling test family). Fixed: the ChargeCycle call gets the identical gate (applicability + nothing-published); rides the next leg. Also: the F4.1 apply-refusal at this commit still says "incomplete or changed" — this predates fixture v2 (58f926a); the clean-control discipline now covers it |
| 36468595972 | 276a1a2 (three residuals) | win FAIL / ubuntu ok / macos ok; windows completed ~19:02Z | ASSESSED invocation 22 (19:07Z): **red=10 — the assertion re-scope CLEARED the retention-scaffolding failure**. THREE precise residuals remain: (1) the F4.1 apply pin reached the apply path but refused EARLIER with "original precharge intent is incomplete or changed" — my minimal constructed intent does not satisfy validateIntentBefore (the construction needs the full binding-matching intent shape: Before/Accounting fields the real fixture produces — a deeper constructor, not a barrier problem); (2) the organizer structure check FIRED: the brief genuinely CREATES deck-tree paths on Windows even under the deny-ACE — a REAL product-write finding needing investigation (the deny-ACE covers add-file; the creation path may pre-date the deny via an earlier handle, or write outside the denied dir shape — the exact created paths need a hosted diff print); (3) the chargeFixture anomaly persists through the BROADENED single-marker match — the raw-error t.Logf rides 8e276e3's leg for exact bytes |
| 36466503956 | 1fbb627 (F4.2 Preview pin) | win FAIL / ubuntu ok / macos ok; windows completed ~18:44Z | ASSESSED invocation 22 (19:00Z): 11 red; **the F4.2 pin FAILED with a test bug — Preview's identity gate (validHash) rejected the plain entry key before any read** ("requires exact original identity"): fixed same invocation (a real digest passes and reaches the missing-intent read). The chargeFixture anomaly PERSISTS at this commit (pre-broadening): the refusal text at :47 contains both match substrings verbatim yet the gate did not fire — a raw-error t.Logf now rides the next leg for exact-byte diagnosis |
| 36465989758 | 12d7cba (chargeFixture tranche) | win FAIL / ubuntu ok / macos ok; windows completed ~18:4xZ | ASSESSED invocation 21 (18:45Z): 11 red; readonly organizer still red (see below); agents-exec + budget-worktree CLI CLEARED |
| 36466394961 | 3be5025 (F1 gate + F2 wiring) | win FAIL / ubuntu ok / macos ok; windows completed ~18:4xZ | ASSESSED invocation 21 (18:45Z): 11 red; **the refusal family dropped 12→1** (F1+F2 effective); agents-exec/budget-CLI cleared. THREE precise residuals: (1) TestOrganizerBriefWritesNoFileReadOnlyDeck "brief wrote into the deck tree" — the deny-ACE now works (cleanup fixed) but the WRITE-DETECTION fired: the brief may actually write on Windows OR the snapshot walk sees mtime churn — needs the F1-style investigation, NOT another revert; (2) TestReservationRecoveryContentCheck… family: the chargeFixture gate's message match apparently FAILED hosted (the refusal text printed at the caller line via t.Helper) despite both Contains substrings being present verbatim — UNRESOLVED, needs exact reproduction; (3) TestRefusalRetentionSurvives: "refusal created state anyway: [d .parley-runtime/]" — the helper's nothing-created assertion fires on scaffolding the test's OWN EARLIER steps created legitimately (kimi's F1 assertion-scoping point applies to the helper: re-scope to refusal-RECORD state) |
| 36462892383 | 4952e1d (refusal tranche 1) | win FAIL / ubuntu ok / macos ok; windows completed ~18:1xZ | ASSESSED invocation 20 (18:23Z): 11 red standing; **tranche 1 (reservation_recovery_test.go, 5 sites) did NOT clear** — the failures moved to a SIBLING file: reservation_recovery_validation_test.go (TestReservationRecoveryContentCheckDoesNotBlockLiveFinish etc.) whose fixtures reach publication through different helpers — the same §C.1 family with more call sites than the first tranche covered. ALSO: TestAgentsExecRetainsFailedStart fails because the test removes `test-agent` (no .exe) to simulate failure — the Windows fixture is test-agent.exe; the failed-start path needs the platform path from writeRoleFixture's return |
| 36463295036 | c39d81c (refusal tranche 2) | win FAIL / ubuntu ok / macos ok; windows completed ~18:1xZ | ASSESSED invocation 20 (18:23Z): 11 red standing; the evidence-package refusal tranche cleared 2 of its ~8 (refusal count 14→12) — partial: the remaining sites are inside subtests whose dir creation is deeper than the wiring matched. Both remaining tranches are precise next units with named call-site files |
| 36462531256 | edf200d (D1'/D2' diagnostic batch + rows 4/16/18 tail) | win FAIL / ubuntu ok / macos ok; windows completed ~18:1xZ | ASSESSED invocation 20 (18:19Z): 11 red standing; **zero replace-diag outputs** — the concurrent-replacement failure did NOT occur this cycle (racy; the diagnostic is armed for the next occurrence, which will arrive graded). The reservation-refusal tranche 1 (4952e1d) is NOT yet in this commit's tree (next cycles carry it) |
| 36458737496 | a1a54f3 (worktreeRow canonical lookup) | win FAIL / ubuntu ok / macos ok; windows completed ~17:37Z | ASSESSED invocation 19 (17:42Z): 11 red standing; TestWorktreeInventoryRetainsRemovedLinkedWorktree still failed at :134 — a SECOND porcelain comparison in the same test (git list output vs the Go path, the row-10 family again); fixed same invocation (canonical ToSlash comparison, bounded local PASS exit 0); hosted leg rides the push |
| 36451945257 | d00c1fb (R1 after-operand fix + R3 pin + corrections) | win FAIL / ubuntu ok / macos ok; windows completed ~16:39Z | ASSESSED invocation 16 (16:52Z): 11 red standing; **R3 pin PASSED hosted** (TestPinLstatOperandComparesWithoutPathReopen); **the concurrent-replacement failure PERSISTED through R1** (same bare-errno "persist operator cycle grant: Access is denied") — a static trace of EVERY policy.json reader (readCyclePolicy, readRuntimePolicy→readStepPolicy/readLaunchPolicy, migration_recovery, runtimeBinding.current — all funnel through the fixed readStepHistoryFile) and writeSynced itself (temp+sync+close+replace, no target open) found NO remaining product zero-share open on that path. Recorded before further reaction; next unit = the Claude-§7-style zero-cost diagnostic on ReplaceSyncedFile's failure path (raw errno + handle attribution) so the next hosted occurrence arrives classified — no retry broadening |
| 36452203943 | 158530f (row-10 declared-path verbatim) | win FAIL / ubuntu ok / macos ok; windows completed ~16:4xZ | ASSESSED invocation 16 (16:52Z): 11 red standing; **all TestDeclared* tests PASSED — row 10 CLOSED hosted** (canonical comparison + declared-verbatim row fidelity) |
| 36451116684 | 615a230 (Retains second call site) | win FAIL / ubuntu ok / macos ok; windows completed ~16:2xZ | ASSESSED invocation 15 (16:35Z): 11 red standing; **TestSnapshotRevalidationRetainsIdenticalMaterialAfterTimestampChange AND TestSnapshotRoundTripRetainsDirtySource PASSED hosted** — the row-21 mode disposition is fully hosted-green; ledger row 21 closed |
| 36449623422 | 29c908c (row-21 mode disposition) | win FAIL / ubuntu ok / macos ok; windows completed ~16:25Z | ASSESSED invocation 14 (16:26Z): 11 red standing; **TestTreeDigestModeChangeChanges PASSED hosted with the 0600→0400 mutation and TestSnapshotRoundTripRetainsDirtySource PASSED with the 0666 pin** — kimi-1's platform-semantics source inference is now HOSTED-CONFIRMED (read-only is the only chmod-mutable stat-observable mode dimension; the digest binds it). The Retains test still failed: my replace converted the FIRST snapshotMemberBytes call (a different test, same family — coincidentally also needed) and left the Retains call at :353 with the 0600 literal — a staging miss disclosed and fixed in 615a230 (hosted leg rides its push). TestSnapshotRelativeLinksRoundTrip fails separately ("link is absolute, escaping or unsupported") — the symlink round-trip family, recorded as the next unit |
| 36449102035 | 966f211 (organizer-deck DenyWrite revert) | win FAIL / ubuntu ok / macos ok; windows completed ~16:19Z | ASSESSED invocation 14 (16:20Z): 11 red = standing set; **the revert is hosted-verified** — TestOrganizerBriefWritesNoFileReadOnlyDeck PASSES again; the two open questions (GENERIC_WRITE mapping vs RemoveAll; brief write-detection under a real read-only deck) remain recorded, the fixture NOT called verified by the ineffective chmod |
| 36447708947 | 4827e67 (rows 8/20 sweep tranche 2) | win FAIL / ubuntu ok / macos ok; windows completed ~16:07Z | ASSESSED invocation 13 (16:08Z): 11 red = standing set; fsacl GREEN (DenyWrite/AllowWrite run clean hosted); the converted fixtures PASSED hosted (TestSaveUnwritableDirFails, TestFailedSaveLeavesNoReport, phase_event unreadable-artifact cases, RunChecksContractEvidenceWriteFailureVetoes, driver impl unreadable-dir). ONE REGRESSION: TestOrganizerBriefWritesNoFileReadOnlyDeck — the deny-ACE broke its TempDir cleanup (openfdat: Access is denied) and the write-detection snapshot misreported; REVERTED to the historical chmod form same invocation (hosted confirmation of the revert rides the push); the two exposed questions recorded: (a) does deny-GENERIC_WRITE on a dir really exclude RemoveAll's needs (the GENERIC_WRITE mapping must be re-derived, not assumed), (b) the brief write-detection under a real read-only deck. TestReviewRoundHasFindingsFailsClosed confirmed pre-existing (failed in 92beef6 too) |
| 36445335968 | 92beef6 (row-30 remediation: eager identity pins + fixture + CRLF/consult) | win FAIL / ubuntu ok / macos ok; windows completed ~15:49Z | ASSESSED invocation 12 (15:50Z): **row-30 remediation HOSTED-VERIFIED** — the {root,inode} subtests of both revalidation tests have NO failure lines (the eager pins detect the same-path replacement; the fixture now exercises the real invariant on every platform). The only remaining revalidation failure is TestSnapshotRevalidationRetainsIdenticalMaterialAfterTimestampChange = the row-21 mode-synthesis family (tar header mode), a distinct pre-existing phenomenon. 11 red = same standing set as 4ed454a — no regressions from the 14-site pin conversion |
| 36444086047 | 4ed454a (row 9 USERPROFILE fixture) | win FAIL / ubuntu ok / macos ok; windows completed ~15:34Z | ASSESSED invocation 12 (15:42Z): **11 red — the agents package went GREEN** (row 9's USERPROFILE fix cleared TestZcodeResolvesModelAndEffortFromItsOwnConfig + the Kimi config test); remaining TestZcode*/DeckDeclared* failures are in app (row 23 agent-runtime family); config package failures are NOT the row-9 shape (no HOME setenv there — row 18 stays provisional with that evidence) |
| 36443342532 | a2d6af6 (CRLF pinning + consult engagement + row-30 probe) | win FAIL / ubuntu ok / macos ok; windows completed ~15:31Z | ASSESSED invocation 12 (15:33Z): **TestSameFileLazyVsEagerPin PASSED** — both predictions held (lazy pin compares EQUAL to a same-path replacement; eager pin detects it) → row 30's mechanism is the Go lazy FileInfo, Claude-1's consult confirmed as hosted observation; remediation unblocked. 12 red = standing set minus acp (CRLF pinning produced no regressions) |
| 36442736966-superseded | (pre-note row kept for the audit trail: the expected failure was described as a "false red" — wording corrected above; it was a real failure of the incomplete commit) | | | |
 | **PRE-NOTE (staging miss, disclosed):** the checks-path W-SHELL product code was edited in invocation 10 but missed by the 9b880d3 staging; its test landed here WITHOUT the code, so TestRunChecksMissingShRefusesPreWork is EXPECTED to fail this cycle on the message assertion — a false red corrected by 2da24e1, not a product regression. Assess with that lens; the ACP-split and AF_UNIX legs are genuine |
| 36441834810 | 9b880d3 (Stage 5 core: P-A + W-SHELL) | win FAIL / ubuntu ok / macos ok; windows completed ~15:27Z | ASSESSED invocation 11 (15:28Z): 13 red = standing set; **zero failures among ALL Stage 5 pins** — TestRunCriterionControlledMissingShRefusesPreWork + WithShRunsNormally, and (in the earlier same-code cycles) the lock diagnostic — P-A/W-SHELL hosted-confirmed on first Windows execution |
| 36440560978 | bc08cd7 (identity-read retry + review-scope evidence) | win FAIL / ubuntu ok / macos ok; completed ~15:26Z | ASSESSED invocation 11 (15:25Z): **the :394 identity-read retry WORKED** — zero .lock-path sharing failures (the ledger_test/review_test raw opens are gone); the failure MOVED to the lock-ORIGIN file (readLockOrigin's os.Open), the site left unwired on a since-falsified premise; recorded in row 7 and fixed same invocation (openOriginFileRead, narrower classification). 13 red = standing set |
| 36439743640 | 17d5210 (delete-share readers) | win FAIL / ubuntu ok / macos ok; completed ~15:10Z | ASSESSED invocation 10 (15:13Z): 13 red = standing set; **delete-share verdict: "persist operator cycle grant" = 0** — the attributed concurrent-persist Access-denied is GONE in the cycle carrying the OpenSharedDelete fix (present once in 36435481011, absent since). One residual "Access is denied": TestSnapshotRevalidationRefusesChangeDuringRead/root at :449 — the during-read variant of the row-19 redesign tried to rename `sub` while the verifier held the source OPEN INSIDE it; Windows refuses (containment by design). Recorded, then pinned honestly: the case now asserts the OS-refused swap leaves the verification unchanged on Windows (POSIX keeps the swap-detection assertion) |
| 36438474793 | c6f3db1 (§D.6 lock retry + diagnostic) | win FAIL / ubuntu ok / macos ok; completed ~14:57Z | ASSESSED invocation 10 (14:59Z): the AC-LOCK-4 diagnostic test PASSED hosted (foreign class = ERROR_SHARING_VIOLATION; transient holder recovers, persistent holder exhausts loudly within bounds) — but product lock failures persisted with RAW unretried messages: the failing open was lockIdentity :394, unwired in invocation 9. Row 7 updated; fix (openLockFileRead) landed as bc08cd7 |
| 36439175112 | 08f769c (row-19 redesign) | win FAIL / ubuntu ok / macos ok; windows completed ~14:58Z | ASSESSED invocation 10 (15:04Z): 13 red = standing set; row-19 redesign MECHANICS VERIFIED (no rename failure; zero old-root messages) — the residual root/inode failures are the PRE-EXISTING SameFile family (row 30, identical in prior cycles); "persist operator cycle grant" ABSENT this cycle (delete-share verdict inconclusive — racy presence; 17d5210's cycle 36439743640 carries the fix) |
| 36437395294 | e22b0ff (AC-DUR-3 dormant conversions + grammar pin) | win FAIL / ubuntu ok / macos ok; windows leg completed ~14:4xZ (macos still running at assessment) | ASSESSED invocation 9 (14:53Z): 13 red = the standing set; no failures among the new mechanics/grammar tests (superseded by the fuller 4e239fe assessment below — same code family) |
| 36437548293 | 4e239fe (AC-DUR-3 + §D.6 read-only trap) | win FAIL / ubuntu ok / macos ok; completed ~14:52Z | ASSESSED invocation 9 (14:53Z): 13 red = the standing set (fsacl, pipeline green); **ALL FOUR new tests PASSED hosted on first Windows execution** — TestPublishFileDurableMechanics, TestPublishDirDurableMechanics, TestReplaceSyncedFileOverReadOnlyTarget (AC-LOCK-2 pin), TestStageSuffixOutsideEveryFinalGrammar. **Stage 3 hosted confirmation COMPLETE** |
| 36435481011 | 32bacd2 (A2/B3/B4 gate pins + app dir-sync via contract) | win FAIL / ubuntu ok / macos ok; completed ~14:38Z | ASSESSED invocation 8 (14:38Z): 13 red = the SAME standing set (fsacl, pipeline green); the A2/B3/B4 gate-pin test has no failure line = PASSED hosted; one residual "Access is denied" occurrence (single instance — not the row-25 dir-sync family; attribute at Stage 4/6 sweep with locator) |
| 36434775624 | e4a2af9 (Stage 3 core) | win FAIL / ubuntu ok / macos ok; completed ~14:33Z | ASSESSED invocation 8 (14:3xZ): the row-25 family TRANSFORMED as designed — "Access is denied" count 0 (was dominant), named refusal "directory-entry durability is not available on Windows" x13; 13 red = same standing set (fsacl, pipeline stay green); the three new windows refusal tests (SyncDirContract, RootedRows A3/B5/C1/C2, VerificationRows gate) have NO failure lines = PASSED hosted. The ~190 failing funcs in trajectory/driver/budget/evidence are the standing families now expressing the designed §C.1/§C.2 refusals (reservation/parent-recovery features refuse on Windows; their POSIX-exercising tests fail on the refusal — NOT treated as passing; reconciliation of those tests is Stage 6/7 review territory, no unilateral skips) |
| 36432547744 | 6236af6 (Stage 2 complete + LoadGate fix) + a9c944e docs | win FAIL / ubuntu ok / macos ok; windows leg completed ~14:14Z | ASSESSED invocation 7 (14:15Z): **pipeline package GREEN** — the LoadGate ERROR_INVALID_NAME fix and platform-conditional legacy test are hosted-verified; fsacl stays GREEN; 13 red = prior 14 minus pipeline, all standing Stage 4/5/6 families. Row 29 (drift test) NOT reproduced this cycle — intermittent, row stays open |



## Reconciliation ledger (§D.9 — emitted at Stage 0, before the sweep; refreshed per hosted cycle)

Form: hosted failing class ↔ source-context class ↔ FINAL §D/§B item ↔ stage ↔ exclusion row.
**Exclusion rows: NONE exist yet.** A Windows test exclusion is admissible only when the
behavior under test cannot exist on Windows (§D.9); each needs its own reviewed row.

**Refreshed 2026-09-25T18:16Z from the actual Windows denominator** (runs 36170672078
and 36170742720, identical red sets: 14 red / 16 ok packages, 111 failing test
functions): rows 1–2 verified fixed-or-extended hosted; row 5 W1 surface CONFIRMED
hosted (evidence package privacy signature x~71); new hosted phenomena recorded as
rows 16–18 before any code reaction. Rows not yet surfaced hosted keep their
source-context classing.

| # | Hosted phenomenon (run 35987916696 / 36009912946) | Class | FINAL item | Stage | Exclusion row |
|---|---|---|---|---|---|
| 1 | `internal/evidence` test build failure: `syscall.Mkfifo` undefined on Windows | W4 (TEST) | §D.8, AC-BLD-1 | 0 | none — FIXED AND HOSTED-VERIFIED: evidence package compiles and executes in both cycles (its W1 privacy failures now visible = row 5) |
| 2 | `internal/app` aborts 3.499s: unchecked type assertion panic at `app_test.go:155`; no `wait`/`usage` results ever on Windows | W6 panic half (TEST) | §D.8, AC-BLD-1 | 0 | none — CLOSED hosted (36172430646): internal/app ran to completion 165s, no panic; wait/usage suites executed for the first time and PASSED (zero failures attributed to them) |
| 3 | CONVERSION COMPLETE invocation 17–18 (re-exec port; see the runaway postmortem for the recursion defect and repair; hosted leg rides the push). Originally: `writeFakeParleyDeckSkill` extension-less `#!/bin/sh` fixture unexecutable on Windows (W6 fixture half) | W5/W6 (TEST) | §D.9, AC-FIX-1 | 6 | none planned — re-exec/`cmd.exe /c` port |
| 4 | `#!/bin/sh`/extensionless fixtures: `launch_test.go`, `protocol_context_test.go`, `telemetry_test.go` (~13) | W5 (TEST) | §D.9, AC-FIX-1 | 6 | none planned — test-binary re-exec |
| 5 | Snapshot privacy guard `Perm()&0077` never passes (0777/0666 synthesis), 71 messages / ~69 tests | W1 (PRODUCT) | §D.1, AC-PRIV-1..6 | 1 | none — real ACL implementation; CONFIRMED hosted in both cycles (evidence, driver refusal families; "snapshot store must be a private real directory" signature) |
| 6 | Pipeline gate filenames with `>` → `ERROR_INVALID_NAME` (9×) | W2 (PRODUCT) | §D.4, AC-NAME-1/2 | 2 | none — universal encoding |
| 7 | HOSTED OUTCOME (36438474793 @c6f3db1): the diagnostic test PASSED hosted (AC-LOCK-4 evidence complete — foreign-holder class is ERROR_SHARING_VIOLATION, transient holders recover within the budget, persistent holders exhaust loudly), BUT the product failures persisted with RAW unretried messages — the failing open was lockIdentity's read (lock.go :394 os.Open), the FIRST .lock open in the reserve chain, which invocation 9 had not wired. Recorded before reaction; fixed same invocation: openLockFileRead (read-only retry variant, same structural classification) wired at :394. Full open audit of the lock chain: :159 kernel-lock open (WIRED, O_RDWR), :200 probe (WIRED, O_RDWR), :394 identity read (NOW WIRED, read-only), :427/:360 CreateTemp (different temp files, not the .lock), readLockOrigin :485 (the lock-ORIGIN file — different path, no hosted failures, deliberately unwired and documented), migration CreateTemp (temp). Hosted confirmation of the :394 fix rides this push. The wider family: cross-process lock/open/rename (~14 sites incl. `ledger_test.go:57`, `review_test.go:438`, `verification_test.go:368`) | W3 (PRODUCT) | §D.6, AC-LOCK-1..4 | 4 | none — AC-LOCK-1 close-before-rename audit remains |
| 8 | chmod-unreadable fixtures (~6: `consensus impl_test.go:807`, `phase_event_test.go:156/170/186`, `strict_gate_test.go:179`) | W7 (TEST) | §D.9 denyRead | 1 (helper) / 6 (sweep) | none — DACL deny-ACE helper |
| 9 | RESOLVED HOSTED (36444086047: agents package GREEN): fakeHome now sets USERPROFILE alongside HOME (§D.9: os.UserHomeDir resolves USERPROFILE on Windows; the HOME-only fixture leaked the real runner home). Originally HOME-only fixtures (`agents/configmodel_test.go`, 2) | W8 (TEST) | §D.9 | 6 | none — USERPROFILE set |
| 10 | CLOSED HOSTED (36452203943: all TestDeclared* passed — canonical comparison + declared-verbatim fidelity). FIX LANDED invocation 14: canonicalWorktreePath (filepath.ToSlash(filepath.Clean)) at all three comparison sites — isDeclaredUnavailable, the registered map, and the retained-registration check — with verbatim forms preserved in persisted records; the test-side assertOutsideDeclared canonicalizes both sides so its no-contributed-history assertion keeps its teeth. Unix behavior identical (ToSlash is identity). RECLASSIFIED invocation 13 (no `/repo` literals exist in the tree — the ledger's original classing was wrong): the actual hosted phenomenon (protocol package, TestDeclared* ×4) is a PRODUCT path-separator mismatch — the fixtures derive `linked` from git worktree porcelain output (forward slashes on Windows), and the declared-path validation compares it against Go-cleaned backslash forms → 'not a retained registration' / 'unavailable' refusals. Fix = separator normalization at the product's declared-vs-registered comparison sites (volume-aware, per-OS refusal assertions kept); named as the next Stage 6 unit | W9→product separator family | §D.9 | 6 | none |
| 11 | CRLF (`hardening_test.go:456`, 1) | W10 (TEST) | §D.7, AC-CRLF-1..3 | 6 | none |
| 12 | ACP drain race (ubuntu; load-dependent) | U1 (PRODUCT, already repaired in base) | §D.7 preserve | n/a | none — spawn.go Stop/Wait untouched |
| 13 | Captured-execution family root cause not established (ubuntu) | U2 (undesignated) | §D.5 sibling discipline | probe/H4 | none — diagnose, never label flaky |
| 14 | refusalGit identity dependence | U3 (PRODUCT, already repaired in base) | preserved | n/a | none |
| 15 | CLOSED HOSTED (36450249028: TestReviewRoundHasFindingsFailsClosed PASSED). FIX LANDED invocation 14 (the §D.5 branch-independent fix, exactly as prescribed: + one Lstat — a round path that exists and is not a directory vetoes explicitly; the ReadDir error classification alone was unreliable for the file-where-dir-belongs case, hosted strict_gate :179). Hosted leg rides the push | §D.5 | §D.5, AC-STRICT-1/2 | shipped | none |
| 16 | `internal/app` panics at `app_test.go:185` (fixed invocation 2). ROOT CAUSE NOW OBSERVED hosted (36172430646): `app_test.go:123` payload shows `parley_deck_skill_error: exec: "parley-deck-skill": executable file not found in %PATH%` — the extension-less fake-skill fixture (row 3) is unexecutable on Windows, so the key is absent; graceful failures now, no panic | W6 (TEST) | §D.8/§D.9, AC-BLD-1, AC-FIX-1 | 6 (fixture port) | none — re-exec/cmd.exe port at Stage 6 |
| 17 | RESOLVED HOSTED (36442736966/36443342532: acp package GREEN — the §D.7 a/b split with the 1 KiB capacity-independent payload fixed both drain tests; U1 spawn.go Stop/Wait ordering untouched). Originally: `internal/acp` TestSpawnStopDrainsStderrBeforeReaping + TestSpawnWaitDrainsStderrBeforeReaping fail 10.01s on WINDOWS legs | ACP family (§D.7) — previously classified ubuntu-only (U1, repaired in base); on Windows this is a NEW phenomenon, cause not yet established | §D.7 split rule | 6 (§D.7 a/b split) + probe-informed | none — record before reaction; do NOT retune; U1 spawn.go Stop/Wait ordering stays untouched; not labeled flaky (it is deterministic in both runs) |
| 18 | TestExpandPlaceholders FIXED invocation 19 (the product substitutes placeholders VERBATIM — the template's literal '/' separators are not normalized; the want is now built by the same semantics, checking the placeholder VALUES; bounded local PASS). The rest of row 18 resolved earlier (agents package green via USERPROFILE). Originally NEW (both cycles): `internal/agents` TestZcodeResolvesModelAndEffortFromItsOwnConfig, TestKimiThinkingEffort; `internal/config` TestLoadAgentSpecsLayersAndTracksSources, TestExpandPlaceholders | class PROVISIONAL W8 (HOME/config-path fixtures) — root cause NOT yet verified; do not treat as settled until probed at Stage 6 | §D.9 | 6 | none — diagnose at sweep; no exclusion anticipated |

| 19 | REDESIGNED invocation 9 (commit 08f769c; hosted confirmation rides its push): both 'root' cases now swap an in-root SUBDIRECTORY (renaming the rooted dir itself is impossible on Windows while the root handle is open — containment by design); invariant preserved: a replaced directory entry must not be hidden by the original still-readable inode; all subtests green on darwin. Originally NEW (36172430646): `trajectory/snapshot_read_test.go:261/:432` — os.Rename of a directory fails `The process cannot access the file because it is being used by another process` (open handle on source tree) | W3 sharing family (PRODUCT/TEST site) | §D.6, AC-LOCK-1..4 | 4 | none — recorded before reaction; structural handle discipline first (§G H7 mapping) |
| 20 | NEW (36172430646/36172828285): `evidence/refusal_test.go:299` `sync ...: Access is denied`; `tree_report_test.go:300/:311` TestSaveUnwritableDirFails + TestFailedSaveLeavesNoReport — chmod-unwritable/unreadable fixtures do not block writes/syncs on Windows | W7 chmod-fixture family (TEST) | §D.9 denyRead | 6 (sweep) / 1 (helper now exists) | none — DenyRead helper landed and hosted-verified; sweep converts fixtures |
| 21 | CLOSED HOSTED (36451116684: Retains + RoundTrip PASSED with the pinned constants; the digest mode test passed in 36449623422). DISPOSITION LANDED invocation 14 (kimi-1 consult engaged in the Decision Log): platform-conditional mutation + pinned OS-observed constants — Unix byte-identical, Windows exercises the read-only dimension (0600→0400 digest test; 0666 archive/restore constants). Originally: `evidence/tree_report_test.go:87` TestTreeDigestModeChangeChanges — chmod 0600→0700 does not change the synthesized mode on Windows, so the tree digest does not change | mode-synthesis family, NEW distinct phenomenon (TEST + product semantics question) | §D.9 sweep; digest mode semantics need review | 6 | none — record first; the digest's mode-sensitivity on Windows needs a deliberate decision, not a silent fixture change |
| 22 | RE-FIXED invocation 19 (hosted showed TAB is ALSO invalid in Windows filenames — the Windows odd-name set is now space/žltý/-option/plus+/eq=, all hosted-legal; inventory invariant pinned unconditionally). FIXED invocation 18 (platform-true odd-name set; inventory invariant pinned on both; Unix unchanged — hosted leg rides the push). Originally NEW (36172430646): `evidence/source_inventory_test.go:136` — fixture filename containing a newline (`line\nbreak`) fails to open: `The filename, directory name, or volume label syntax is incorrect` | W2-adjacent invalid-name class in FIXTURES (TEST) | §D.9 | 6 | none — fixture portability (t.TempDir-compatible names) |
| 23 | FIX LANDED invocation 19 (hosted leg rides the push): the zcodeDeck fake home now sets USERPROFILE (the row-9 shape — the zcode adapter resolves its config from USERPROFILE on Windows); the writeFakeZcode stub is a zcode-probe re-exec role (EQUALS-form argv check + the prompt/sentinel extraction from ARGV, not stdin — a port defect caught by the bounded rerun; bounded local PASS exit 0, 10.4s). NEW (36172430646): `internal/app` `app_test.go:373` `code=1 stdout=codex: not installed`, `:421` agent-runtime resolution; correlates with row 18's agents/config failures | W8/agent-runtime family PROVISIONAL — root cause not yet verified (runner PATH lacks real agents; tests presumably fake them) | §D.9 | 6 | none — diagnose at sweep; no exclusion anticipated |
| 24 | NEW (36174658770): `trajectory/snapshot_test.go:493` TestSnapshotConcurrentCapturePublishesOneExactArchive — "concurrent publication differs: {ref:zero err:<e1>} {ref:zero err:<e2>}" (two DISTINCT unpublished results; error VALUES not surfaced in-band) | concurrency × §D.1 first-creation: CANDIDATE root cause (not established) is the create→set-DACL window — a concurrent observer's Lstat sees the dir exist before the owner-only DACL is applied, misclassifying a concurrent product creation as pre-existing and refusing. Recorded before reaction; fixture change removed the test's exposure (test PASSED hosted in 36425527266); product race remained possible | §D.1 | 1 | none — FIX LANDED invocation 5: atomic CreateDirectory with SECURITY_ATTRIBUTES owner-only SD (SDDL `D:P(A;;GA;;;<sid>)`), ERROR_ALREADY_EXISTS → verify-not-refuse; adversarial pin TestConcurrentFirstCreationNeverRefuses rides the next hosted cycle, which closes the row |
| 25 | TRANSFORMED by e4a2af9 (36434775624 hosted-verified: Access-denied count 0, named refusal x13): originally UNMASKED by 36174658770 — `sync <dir>: Access is denied` was the dominant residual family in trajectory/driver/evidence/budget (`state_test.go:63/:478`, `verification_test.go:368`, `refusal_test.go:35/:97/:140/:210/:251/:270/:299`, `driver/trajectory_test.go:50`) — dir-fsync on a read-only dir handle; previously masked by the W1 store-guard failure | H2-adjacent fsync-on-directory class (FlushFileBuffers needs write access; Go dir handles are read-only) — rows 19/20's families now confirmed broader | §B (fsutil.SyncDir audit), §D.6 | 3 | none — the §B audit-of-record dispositions and named-type contract govern the fix; never a per-site skip |
| 26 | NEW (36425527266): fsacl TestPreexistingStoreRefusedAndNotRewritten fails at the final marker ReadFile (`Access is denied`) — `grantTrustees`' protected-DACL replacement on the store dir auto-propagates inheritance removal to children, stripping the marker file's inherited-only access. Product refusal verified correct by the hosted log (store DACL before/after identical, refusal text exact) — TEST-scaffolding phenomenon, second layer under row-24-cycle's assertion fix | NTFS auto-inheritance propagation (test scaffolding) | §D.1, AC-PRIV-5 | 1 (test fix) | none — marker gets an explicit owner-allow ACE (readability independent of store-dir DACL churn); content assertion unchanged |
| 27 | FIX LANDED invocation 19 (platform-true first-refusal assertion: Windows expects the reviewed POSIX-host gate text, Unix the ignore-prerequisite text; refusal-retention and nothing-published hold on both; bounded local PASS). NEW (36425527266, unmasked by the W1 clear): evidence_publication_test.go:35 reaches product refusal "independent evidence verification requires a POSIX execution host; Windows runtime is not supported" (driver_impl.go:543; sibling trajectory_verify.go:107) | W-SHELL / POSIX-host product refusal family (§D.3) | §D.3 | 5 | none — the §D.3 named-prerequisite refusal design governs; never a skip |
| 28 | NEW (36425527266, run-internal regression caught hosted): the 2 residual W1 signatures (snapshot_test.go:147, source_inventory_test.go:111) are both RestoreSnapshot refusing its OWN MkdirTemp destination — the d07e6cc restore wiring called EnsurePrivateStore on an already-created dir, which the guard correctly classifies as pre-existing user data. Found by code inspection before reading the log; log confirmed both sites | implementation defect in this run's Stage-1 wiring (not an environment phenomenon) | §D.1 creation policy | 1 | none — fixed this invocation: fsacl.ProtectPrivateStore (creation policy, no refuse-and-instruct) for product-created dirs; Unix verify-only byte-identical; pinned by TestProtectPrivateStoreAppliesPolicyToProductCreatedDir |
| 31 | FIX LANDED invocation 19 (evidence-based, from the fa54033 named-refusal diagnostic hosted in 36457689897: Readlink returns the native BACKSLASH form "alias\\file" on Windows — capture now normalizes to the archive's canonical slash separator at the readlink boundary, validSnapshotLink unchanged; the round-trip test compares readlinks canonically on Windows, the raw-target invariant preserved; bounded local PASS). NEW family (first attributed invocation 15; failing since at least 36445335968): TestSnapshotRelativeLinksRoundTrip — capture refuses one of the five relative-symlink shapes with "snapshot link is absolute, escaping or unsupported". Analysis: validSnapshotLink's charset (backslash/colon), IsAbs and in-root-resolution checks all pass for the five fixture shapes on paper; the Windows-specific unknown is the Readlink-returned target form. The refusal now NAMES the link and raw target (fa54033) so the next hosted cycle identifies the exact shape — diagnosis continues from that evidence, no speculative fix | symlink round-trip family (TEST + possible product link-semantics question) | §D.9 sweep | 6 | none — record first |
| 30 | RESOLVED HOSTED (36445335968 @92beef6: the {root,inode} revalidation subtests pass with the eager pins; mechanism previously CONFIRMED by the probe 36443342532 — SameFile(lazyPin, after)==true AND SameFile(eagerPin, after)==false: the cause is Go's Windows LAZY FileInfo (os.Stat/os.Lstat defer identity to the first SameFile call, re-opening by path at comparison time; the pinned operand re-resolves to the replacement), NOT NTFS identity reuse. Claude-1's consult mechanism stands as observed fact; remediation unblocked (invocation 12). Originally mis-attributed as NTFS file-ID semantics, invocation 10 (present identically in 36432547744/36434775624/36439175112): TestSnapshotRevalidationRejectsDifferentMaterial/{root,inode} — "different material or canceled verification was accepted": on Windows the Lstat→open os.SameFile chain does NOT detect a rename+recreate replacement (NTFS file-ID semantics; the recreated file presents the same file identity). NOT caused by the row-19 redesign (identical failures predate it); the redesign removed the rename-mechanics failure as intended | file-identity-synthesis family (TEST + product semantics question — the inode-tamper guard layer) | §D.9 sweep; product SameFile-guard semantics need a deliberate decision (the verify already re-reads and hashes; the SameFile guard is the failing layer) | 6 | none — record first; no silent fixture change |
| 29 | FIX LANDED invocation 19 (the §D.7 split rule verbatim: the drift anchors are DOCUMENT-SEMANTIC comparisons, so CRLF normalization is admissible — the deck read now normalizes \r\n; a Windows checkout with autocrlf reads \r endings, the anchors unchanged in meaning; bounded local PASS exit 0. This also explains the intermittency: only autocrlf-configured checkouts tripped). STANDING (present in 36425527266 and 36429107801, first attributed invocation 6): pipeline TestEmbeddedDefaultMatchesLiveDeck — drift guard fails closed: anchor "## 2. Active agents (roster)" appears 0 times in the embedded-default comparison on Windows | PROVISIONAL §D.7 CRLF family (embedded-default vs live-deck file comparison; CRLF checkout on windows runners is the prime suspect) or roster-anchor family — root cause NOT yet verified; do not treat as settled until probed at Stage 6 | §D.7 | 6 | none — record first; per-invocation `-c core.autocrlf=false` pinning is the §D.7 load-bearing layer |

Refresh rule: after every hosted cycle, re-diff failing vs ledger; new phenomenon ⇒ new row
before any code reaction; no `t.Skip` added without a row (AC-FIX-2).

## Deferred / follow-up naming (AC-PROC-4)

P-B (Job Objects `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `GetProcessTimes` attribution,
`QueryFullProcessImageName` in kill-safety gating) is **not implemented in this idea**;
named follow-up idea required with adversarial review (FINAL §D.2). The `mvdan.cc/sh/v3`
supervisor port is likewise a separate reviewed idea.

## Outcomes & Retrospective

(to be written at completion)
