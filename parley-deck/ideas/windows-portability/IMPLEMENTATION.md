---
idea: windows-portability
status: in-progress
implementer: zcode-1
started: 2026-09-25
branch: windows-portability
head-commit: cd82025 at this invocation's start (2026-09-25T18:23Z); prior checkpoints 356fbb8 (FINAL freeze) -> e9cf601 (claim) -> 255f1a5 (Stage 0 code) -> 9a96c2f/16824fd (stage 0 close) -> ff2ef9d (recon) -> ca5efef (fsacl unwired) -> cd82025 (ACL checkpoint)
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
- [ ] **Stage 1** — snapshot privacy/ACL (§D.1) incl. `denyRead` helper.
      IN PROGRESS (2026-09-25T18:23–18:38Z): `internal/fsacl` module landed AND
      WIRED this invocation: `privateSnapshotDirectory` (snapshot.go:158-162)
      routes through Ensure/VerifyPrivateStore; archive temp file protected via
      ProtectPrivateFile (snapshot.go:206) after CreateTemp; read-back at :546
      via VerifyPrivateFile — all three Windows-only effects, Unix byte-identical
      (AC-PRIV-6). Windows Ensure now distinguishes created-vs-preexisting stores:
      pre-existing failing stores get refuse-and-instruct (repair text names the
      user-invoked repair; DACL never rewritten in place, AC-PRIV-5); fresh stores
      get set+requery self-check. Reparse-point rejection added beyond symlinks.
      Adversarial suite `fsacl_windows_test.go` (AC-PRIV-1..6 windows half):
      own-creation round-trip with independent raw requery; BUILTIN\Users (S-1-5-32-545)
      and Administrators (S-1-5-32-544) grants refused naming the trustee;
      inherited-only refused; pre-existing refused+not-rewritten (marker file intact,
      DACL still grants Users); symlink/non-dir refused (symlink failure = fatalf,
      never skip); file policy round-trip. Shared test's Unix perm assertion split
      into fsacl_unix_test.go (hosted fact: dir perms synthesize 0777 on Windows).
      Hosted evidence for the UNWIRED module (run 36172828285): Ensure/verify/
      symlink-rejection/DenyRead all PASSED on real Windows; only the old perm-bit
      assertion failed (fixed). Remaining: wired-product hostile Windows execution
      (next push), restore-file privacy decision (see Progress note), AC-PRIV-4
      design-level refusal documentation (§N blind spot: no FAT volume on runners).
- [ ] **Stage 2** — gate-name encoding + raw-ID allowlist (§D.4), legacy fallback +
      shadow retirement.
- [ ] **Stage 3** — directory durability: §B table dispositions BEFORE code;
      `fsutil.SyncDir` named-type contract = audit of record + rooted refusal emitter.
- [ ] **Stage 4** — file sharing (§D.6) + open-site diagnostic cycle (H7).
- [ ] **Stage 5** — P-A liveness (§D.2) + W-SHELL missing-`sh` refusal (§D.3).
- [ ] **Stage 6** — fixture-portability sweep (§D.9) + ACP/AF_UNIX/CRLF (§D.7).
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

## Hosted run register

| run id | commit | legs | outcome | notes |
|---|---|---|---|---|
| 35987916696 | 4e0c069 (main) | win FAIL / ubuntu FAIL / macos ok | historical | also failed Ubuntu; per-run claims separate |
| 36009912946 | (main) | 14 red vs 16 ok packages on Windows | historical | baseline for ledger refresh |
| 36170672078 | 255f1a5 (stage 0) | win FAIL / ubuntu ok / macos ok; completed 18:13:58Z | Stage 0 diagnostic cycle 1: Windows 14 red pkgs / 16 ok / 111 failing funcs; evidence package runs (W4 fixed); internal/app panic moved to :185 (residual unchecked assertion, fixed in invocation 2); wait/usage still never executed |
| 36170742720 | 9c1db32 (docs-only over 255f1a5) | win FAIL / ubuntu ok / macos ok; completed 18:12:48Z | Diagnostic cycle 2 (triggered by the run-record push): Windows red set IDENTICAL to 36170672078 (same 14 pkgs, same 111 funcs) — denominator deterministic; ACP drain tests fail on Windows in BOTH runs (new row 17) |
| 36172430646 | 16824fd (probe bundle + app comma-ok) | win FAIL / ubuntu ok / macos ok; completed 18:29:37Z (win leg 12m34s) | H1–H7 bundle EXECUTED: winprobe 7/7 PASS (H3 NTFS+rename mechanics, H7 readonly-rename-over pair, H1 ACL chain — OBSERVED; H2/H4/H6 recorder values unknown-in-band without -v). wait/usage first hosted execution PASSED (row 2 closed). Same deterministic 14 red pkgs; rows 16 confirmed, 19–23 added |
| 36172828285 | ca5efef (fsacl unwired) | win FAIL / ubuntu ok / macos ok; completed ~18:32:50Z | fsacl FIRST hostile Windows execution: Ensure/verify round-trip, pre-existing refusal, symlink+non-dir rejection, DenyRead ALL PASS; os.Symlink works on runner (observed). 15 red = same 14 + fsacl perm-bit assertion (Unix mechanics, fixed this invocation in fsacl_unix_test.go split). W1 signature persists (81x) until wiring lands hosted |

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
| 3 | `writeFakeParleyDeckSkill` extension-less `#!/bin/sh` fixture unexecutable on Windows (W6 fixture half) | W5/W6 (TEST) | §D.9, AC-FIX-1 | 6 | none planned — re-exec/`cmd.exe /c` port |
| 4 | `#!/bin/sh`/extensionless fixtures: `launch_test.go`, `protocol_context_test.go`, `telemetry_test.go` (~13) | W5 (TEST) | §D.9, AC-FIX-1 | 6 | none planned — test-binary re-exec |
| 5 | Snapshot privacy guard `Perm()&0077` never passes (0777/0666 synthesis), 71 messages / ~69 tests | W1 (PRODUCT) | §D.1, AC-PRIV-1..6 | 1 | none — real ACL implementation; CONFIRMED hosted in both cycles (evidence, driver refusal families; "snapshot store must be a private real directory" signature) |
| 6 | Pipeline gate filenames with `>` → `ERROR_INVALID_NAME` (9×) | W2 (PRODUCT) | §D.4, AC-NAME-1/2 | 2 | none — universal encoding |
| 7 | Cross-process lock/open/rename "file is being used by another process" (~14: `ledger_test.go:57`, `review_test.go:438`, `verification_test.go:368`) | W3 (PRODUCT) | §D.6, AC-LOCK-1..4 | 4 | none |
| 8 | chmod-unreadable fixtures (~6: `consensus impl_test.go:807`, `phase_event_test.go:156/170/186`, `strict_gate_test.go:179`) | W7 (TEST) | §D.9 denyRead | 1 (helper) / 6 (sweep) | none — DACL deny-ACE helper |
| 9 | HOME-only fixtures (`agents/configmodel_test.go`, 2) | W8 (TEST) | §D.9 | 6 | none — USERPROFILE set |
| 10 | `/repo` POSIX-path worktree fixtures (~3-4) | W9 (TEST) | §D.9 | 6 | none — t.TempDir + volume-aware joins |
| 11 | CRLF (`hardening_test.go:456`, 1) | W10 (TEST) | §D.7, AC-CRLF-1..3 | 6 | none |
| 12 | ACP drain race (ubuntu; load-dependent) | U1 (PRODUCT, already repaired in base) | §D.7 preserve | n/a | none — spawn.go Stop/Wait untouched |
| 13 | Captured-execution family root cause not established (ubuntu) | U2 (undesignated) | §D.5 sibling discipline | probe/H4 | none — diagnose, never label flaky |
| 14 | refusalGit identity dependence | U3 (PRODUCT, already repaired in base) | preserved | n/a | none |
| 15 | `strict_gate` `:179` file-where-directory-belongs no-veto | §D.5 | §D.5, AC-STRICT-1/2 | with probe bundle | none — branch-independent fix ships |
| 16 | `internal/app` panics at `app_test.go:185` (fixed invocation 2). ROOT CAUSE NOW OBSERVED hosted (36172430646): `app_test.go:123` payload shows `parley_deck_skill_error: exec: "parley-deck-skill": executable file not found in %PATH%` — the extension-less fake-skill fixture (row 3) is unexecutable on Windows, so the key is absent; graceful failures now, no panic | W6 (TEST) | §D.8/§D.9, AC-BLD-1, AC-FIX-1 | 6 (fixture port) | none — re-exec/cmd.exe port at Stage 6 |
| 17 | NEW (both cycles, deterministic): `internal/acp` TestSpawnStopDrainsStderrBeforeReaping + TestSpawnWaitDrainsStderrBeforeReaping fail 10.01s on WINDOWS legs | ACP family (§D.7) — previously classified ubuntu-only (U1, repaired in base); on Windows this is a NEW phenomenon, cause not yet established | §D.7 split rule | 6 (§D.7 a/b split) + probe-informed | none — record before reaction; do NOT retune; U1 spawn.go Stop/Wait ordering stays untouched; not labeled flaky (it is deterministic in both runs) |
| 18 | NEW (both cycles): `internal/agents` TestZcodeResolvesModelAndEffortFromItsOwnConfig, TestKimiThinkingEffort; `internal/config` TestLoadAgentSpecsLayersAndTracksSources, TestExpandPlaceholders | class PROVISIONAL W8 (HOME/config-path fixtures) — root cause NOT yet verified; do not treat as settled until probed at Stage 6 | §D.9 | 6 | none — diagnose at sweep; no exclusion anticipated |

| 19 | NEW (36172430646): `trajectory/snapshot_read_test.go:261/:432` — os.Rename of a directory fails `The process cannot access the file because it is being used by another process` (open handle on source tree) | W3 sharing family (PRODUCT/TEST site) | §D.6, AC-LOCK-1..4 | 4 | none — recorded before reaction; structural handle discipline first (§G H7 mapping) |
| 20 | NEW (36172430646/36172828285): `evidence/refusal_test.go:299` `sync ...: Access is denied`; `tree_report_test.go:300/:311` TestSaveUnwritableDirFails + TestFailedSaveLeavesNoReport — chmod-unwritable/unreadable fixtures do not block writes/syncs on Windows | W7 chmod-fixture family (TEST) | §D.9 denyRead | 6 (sweep) / 1 (helper now exists) | none — DenyRead helper landed and hosted-verified; sweep converts fixtures |
| 21 | NEW (36172430646): `evidence/tree_report_test.go:87` TestTreeDigestModeChangeChanges — chmod 0600→0700 does not change the synthesized mode on Windows, so the tree digest does not change | mode-synthesis family, NEW distinct phenomenon (TEST + product semantics question) | §D.9 sweep; digest mode semantics need review | 6 | none — record first; the digest's mode-sensitivity on Windows needs a deliberate decision, not a silent fixture change |
| 22 | NEW (36172430646): `evidence/source_inventory_test.go:136` — fixture filename containing a newline (`line\nbreak`) fails to open: `The filename, directory name, or volume label syntax is incorrect` | W2-adjacent invalid-name class in FIXTURES (TEST) | §D.9 | 6 | none — fixture portability (t.TempDir-compatible names) |
| 23 | NEW (36172430646): `internal/app` `app_test.go:373` `code=1 stdout=codex: not installed`, `:421` agent-runtime resolution; correlates with row 18's agents/config failures | W8/agent-runtime family PROVISIONAL — root cause not yet verified (runner PATH lacks real agents; tests presumably fake them) | §D.9 | 6 | none — diagnose at sweep; no exclusion anticipated |

Refresh rule: after every hosted cycle, re-diff failing vs ledger; new phenomenon ⇒ new row
before any code reaction; no `t.Skip` added without a row (AC-FIX-2).

## Deferred / follow-up naming (AC-PROC-4)

P-B (Job Objects `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `GetProcessTimes` attribution,
`QueryFullProcessImageName` in kill-safety gating) is **not implemented in this idea**;
named follow-up idea required with adversarial review (FINAL §D.2). The `mvdan.cc/sh/v3`
supervisor port is likewise a separate reviewed idea.

## Outcomes & Retrospective

(to be written at completion)
