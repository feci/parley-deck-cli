---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: participant implementer (release CI remediation)
artifact: release-ci-fix
date: 2026-09-24
---

# Release CI fix — zcode-1

Remediates Findings A and C of `release-ci-claude-1.md` under the dispatch scope:
**CI environment remediation only.** One commit on `lean-organizer` changing
`.github/workflows/tests.yml` and this record. No source change, no version bump,
no tag/release/asset change, no push (organizer pushes and collects hosted evidence).

## Change

`.github/workflows/tests.yml` — two environment steps inserted before
`actions/checkout@v4` (pure insertion, +17 lines, no other step touched):

1. **`Configure fixture git identity`** (all three legs): sets a single
   deterministic fixture identity (`user.email parley-ci@example.invalid`,
   `user.name parley-ci`) into the runner's global git config. Two plain
   `git config` lines; identical behavior under bash (ubuntu/macos) and pwsh
   (windows).
2. **`Enable long paths (Windows)`** (`if: runner.os == 'Windows'`):
   `git config --global core.longpaths true` **before checkout**, so the
   233-character deck record path fits under `MAX_PATH` with the
   `D:\a\parley-deck-cli\parley-deck-cli\` workspace prefix (claude-1 Finding C
   arithmetic, 233+37=270>260).

No test is skipped, filtered, retried, or suppressed anywhere. `Build` and `Test`
commands are byte-identical to the workflow that ran in `35983818751`.

## Why this validates the SAME released product source

- Tag `v1.49.0` (06e563e) is untouched: verified identical locally and on the
  remote before and after this work. The skill tag `v2.13.0` (8161e5e, separate
  repo) is untouched by this task.
- Before my commit, `git diff v1.49.0 HEAD -- . ':(exclude)parley-deck'` was
  empty — incident run `35983818751` (on `main` @ 137b1c1) tested product source
  byte-identical to the released tag (claude-1 verified; I re-verified at HEAD).
- My commit adds only `.github/workflows/tests.yml` content and this record.
  Post-commit check recorded below: `git diff v1.49.0 HEAD -- .
  ':(exclude)parley-deck' ':(exclude).github'` is empty — the entire Go module
  is byte-identical to `v1.49.0`. The workflow file is not Go source, is not
  compiled by `go build ./...`, and is not consulted by `go test ./...`.
- The incident channel failed on the *environment* (missing git identity on
  stock runners; Windows `MAX_PATH` at checkout), not on product source, per
  claude-1's root causes, which I re-verified in code and locally (below). A
  matrix run at this commit therefore exercises the released product source
  under a corrected environment. What changed between the incident run and the
  next hosted run is exactly the two defects claude-1 root-caused — nothing else.

## Exact checks performed (local, this machine, macOS arm64)

**Code re-verification (Finding A).** `refusalGit`
(`internal/app/evidence_refusals.go:59`) extends `os.Environ()` with only
`GIT_OPTIONAL_LOCKS=0` — no identity; `commitVerificationRefusal` (`:83`)
commits through it. The `gateScratchRepo` fixture (`driver_evidence_test.go:24`)
inits a scratch repo whose commits carry only per-command `-c` identity that
never persists into config. Confirmed as reported.

**Controlled reproduction harness.** Two throwaway global configs under
`GIT_CONFIG_GLOBAL`, plus `GIT_CONFIG_SYSTEM=/dev/null`, both with
`user.useConfigOnly=true` so hostname auto-detection cannot synthesize an
identity — the second file adds only the fixture identity the workflow now sets:

- Control (no identity): `go test ./internal/app -run TestEvidenceVerifierProductionClosure
  -count=1` → FAIL with the hosted signature exactly: 14 subtests at
  `evidence_verify_test.go:231` + 1 (`report-persistence-failure-recovery`) at
  `:281` with "could not commit the exact refusal record; retained recovery is
  required". Matches hosted attempt 1 (claude-1) and attempt 2 (below) line-for-line.
- Fixture identity only: same command → `ok internal/app 27.291s`, exit 0.
  The workflow's exact identity values clear the failure.

**No authorship assumptions broken.** No test asserts on git commit
author/committer identity; 24 test files that need an identity pass it
per-command via `-c`, which overrides global config. The global fixture identity
cannot change any test-visible commit except those that today depend on ambient
identity (the Finding A path).

**Full suite under the fixture identity** (proves the global identity breaks no
other package on a machine whose ambient identity is masked):
`GIT_CONFIG_GLOBAL=<fixture> GIT_CONFIG_SYSTEM=/dev/null go test ./...
-count=1 -timeout 45m` → **31/31 packages ok, exit 0** (internal/app 507.246s,
internal/trajectory 580.427s).

**Workflow file.** YAML re-parsed after edit (6 steps, expected order:
identity → longpaths → checkout → setup-go → Build → Test); `git diff` is a pure
17-line insertion.

## Hosted rerun — Finding B classification evidence

Rerun of only the ubuntu job: run `35983818751`, attempt-2 job `107589340137`
(2026-09-24T10:15:57Z–10:19:37Z), same commit (`main` @ 137b1c1), same workflow
(the rerun executes the original workflow — my fix is not in it; no push was
involved). Windows and macOS entries were re-registered by the rerun but did
**not** re-execute (their timestamps remain the attempt-1 window 09:51:55–10:04:21);
I claim no new Windows or macOS evidence.

Attempt-2 ubuntu results:

- `internal/app` **FAIL 160.090s** — Finding A reproduced **exactly**: same 15
  subtests of `TestEvidenceVerifierProductionClosure`, 14× `:231` + 1× `:281`,
  same messages. Finding A is deterministic under identity-less hosted Linux,
  2/2 attempts, and is fully explained by the missing-identity root cause.
- `internal/trajectory` **FAIL 168.811s** — but **not** Finding B's signature.
  `TestRecoveredParentResolutionRequiresExactLineage/missing-replacement-terminal`
  (failed at `verification_parent_recovery_test.go:142`, `Steps:1`, attempt 1)
  produced no failure line in attempt 2, i.e. it passed. Instead, a sibling
  one-off: `TestRecoveredParentObservationRequiresExactLineage/original-…`
  failed at `verification_parent_recovery_test.go:509` — the same shared
  assertion ("recovered helper execution did not complete") with `Steps:3`
  (4 expected), `FailureStage:execution`.
- Additionally `internal/app` `TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts`
  failed at `trajectory_parent_recovery_test.go:153` — launch-stage helper record
  (`FailureStage:launch`, `TrajectoryPending:true`), "no matching successful
  observed terminal".
- Remaining 29 of 31 packages `ok`.

**Classification of the trajectory-area failures (Finding B and its siblings).**
Nondeterministic on the hosted ubuntu runner, root cause **not established**.
Across two hosted executions, three distinct tests in the recovered/captured
trajectory-verification execution area failed once each — never the same test twice, and none reproduced locally (claude-1: four ways, including a
no-identity Linux container where the trajectory package passed; me: the
controlled runs above). The failures share one shape: a recovered verifier's
child execution terminates partway (Steps 1/4, 3/4) or never observes a
terminal. This is **consistent with** hosted child-process instability under
load (packages ran 160–171s hosted vs 5–30s locally), which remains a
hypothesis, not a finding. They are **not** attributed to the missing git
identity: claude-1's identity-less container passed this package, and attempt 2's
identity-dependent failures are exactly and only the Finding A signature. My
workflow fix therefore does **not** claim to fix these, and nothing is skipped
or suppressed. The organizer's post-push matrix is the next data point: identity
present, and if these recur there, they are confirmed identity-independent.

## What was NOT done here

- No source change to any Go file; no version bump; `VERSION`, `version.go`,
  `CHANGELOG.md` untouched (they are byte-identical to `v1.49.0`).
- Tags `v1.49.0` and `v2.13.0` untouched; no release, asset, notes, or branch
  other than `lean-organizer`; no push, merge, or publish.
- Closed artifacts (`IMPLEMENTATION.md`, `FINAL.md`, review/) unchanged.
- No adjudication of the changelog's Windows-coverage sentence (claude-1
  Finding D): whether green runs against this identical product source discharge
  it or a patch release is required is for the independent reviewer to decide,
  not asserted here.
- No new hosted Windows execution: Finding C's fix is reasoned (claude-1's log
  + path arithmetic) and awaits the organizer's push for its first hosted test.

## Proposed follow-up (inactive — not implemented, outside A–D)

`refusalGit` in `internal/app/evidence_refusals.go` supplies no identity for its
own internal commits, so on identity-less hosts (containers, minimal images) the
released product degrades: `parley evidence refusals recover` exits 1 and driver
refusals retain-but-cannot-clear (fail-closed, per claude-1's effect analysis).
Hardening it to pass an explicit fallback identity for these internal commits is
a **source change**, needs its own review and a **new version** (claude-1
suggests v1.49.1), and is outside this dispatch. Proposed only; nothing changed.

## Local validation limits

- Local suite execution is macOS arm64 with the fixture identity and ambient
  identity masked — not hosted ubuntu (amd64), not Windows, not the runner
  images' system configs.
- Hosted attempt-2 evidence covers only the ubuntu leg re-executing the
  *original* workflow; my edited workflow has zero hosted executions yet.
- The trajectory-area nondeterminism has no local reproduction, so no local
  claim about it is possible either way.

## Hosted next step (organizer)

Push `lean-organizer` (workflow fix + this record). The push triggers `Tests`
over the full matrix: first-ever hosted Windows checkout+build+test with long
paths, ubuntu with a deterministic identity, macOS re-run. Collect per-leg
conclusions and logs as the channel evidence; the trajectory-area failures,
if any recur, get classified against the two attempts recorded here.
