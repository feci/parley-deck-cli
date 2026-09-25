---
idea: windows-portability
status: in-progress
implementer: zcode-1
started: 2026-09-25
branch: parley-deck-cli#windows-portability
head-commit: 356fbb8 (design/final freeze); Stage 0 code commit recorded in Progress
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

- [ ] **Stage 0** — contract/smoke/ledger: Mkfifo build tag; `app_test.go:155/:157`
      comma-ok; Ubuntu `GOOS=windows` cross-compile guard (amd64+arm64, build+vet);
      reconciliation ledger emitted (below); one diagnostic hosted cycle budgeted.
- [ ] **Early probe bundle** (before/alongside Stage 0–1): H3 rename mechanics + NTFS,
      H2 SyncFile raw error, H4 ReadDir class + toolchain, H1 ACL parent dump,
      H6 git config --show-origin, H7 read-only rename-over. Mechanics/diagnostic only.
- [ ] **Stage 1** — snapshot privacy/ACL (§D.1) incl. `denyRead` helper.
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

- (2026-09-25 ~18:20Z) Invocation 1 started: read full phase-5 packet, 00-prompt.md,
  FINAL.md, source-context; wrote this claim before any code edit (commit e9cf601). Beginning Stage 0.
- (2026-09-25 ~18:35Z) Stage 0 code complete (commit follows): mkfifo build tag
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

## Decision Log

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

## Hosted run register

| run id | commit | legs | outcome | notes |
|---|---|---|---|---|
| 35987916696 | 4e0c069 (main) | win FAIL / ubuntu FAIL / macos ok | historical | also failed Ubuntu; per-run claims separate |
| 36009912946 | (main) | 14 red vs 16 ok packages on Windows | historical | baseline for ledger refresh |
| 36170672078 | 255f1a5 (stage 0) | three-leg Tests matrix | in_progress at 2026-09-25T18:00Z | Stage 0 diagnostic cycle; green NOT assumed; refresh ledger from its Windows log (esp. wait/usage first execution per AC-BLD-1) |

## Reconciliation ledger (§D.9 — emitted at Stage 0, before the sweep; refreshed per hosted cycle)

Form: hosted failing class ↔ source-context class ↔ FINAL §D/§B item ↔ stage ↔ exclusion row.
**Exclusion rows: NONE exist yet.** A Windows test exclusion is admissible only when the
behavior under test cannot exist on Windows (§D.9); each needs its own reviewed row.

| # | Hosted phenomenon (run 35987916696 / 36009912946) | Class | FINAL item | Stage | Exclusion row |
|---|---|---|---|---|---|
| 1 | `internal/evidence` test build failure: `syscall.Mkfifo` undefined on Windows | W4 (TEST) | §D.8, AC-BLD-1 | 0 | none — build-tag fix, test keeps refusal assertion |
| 2 | `internal/app` aborts 3.499s: unchecked type assertion panic at `app_test.go:155`; no `wait`/`usage` results ever on Windows | W6 panic half (TEST) | §D.8, AC-BLD-1 | 0 | none — comma-ok fix |
| 3 | `writeFakeParleyDeckSkill` extension-less `#!/bin/sh` fixture unexecutable on Windows (W6 fixture half) | W5/W6 (TEST) | §D.9, AC-FIX-1 | 6 | none planned — re-exec/`cmd.exe /c` port |
| 4 | `#!/bin/sh`/extensionless fixtures: `launch_test.go`, `protocol_context_test.go`, `telemetry_test.go` (~13) | W5 (TEST) | §D.9, AC-FIX-1 | 6 | none planned — test-binary re-exec |
| 5 | Snapshot privacy guard `Perm()&0077` never passes (0777/0666 synthesis), 71 messages / ~69 tests | W1 (PRODUCT) | §D.1, AC-PRIV-1..6 | 1 | none — real ACL implementation |
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

Refresh rule: after every hosted cycle, re-diff failing vs ledger; new phenomenon ⇒ new row
before any code reaction; no `t.Skip` added without a row (AC-FIX-2).

## Deferred / follow-up naming (AC-PROC-4)

P-B (Job Objects `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `GetProcessTimes` attribution,
`QueryFullProcessImageName` in kill-safety gating) is **not implemented in this idea**;
named follow-up idea required with adversarial review (FINAL §D.2). The `mvdan.cc/sh/v3`
supervisor port is likewise a separate reviewed idea.

## Outcomes & Retrospective

(to be written at completion)
