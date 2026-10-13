---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
role: independent supplemental source evidence check
date: 2026-10-13
reviewed-cli-commit: 99d04b3ecd0740fb834e70ec109c9384c7b7cd4d
current-cli-head: bcf1140ff64175dc21e89386a0c0bc8491c0ea72
process: fresh headless supplemental check after the recorded goal PASS
---

## Scope and freshness

I am zcode-1. This process checks the terminal Windows CI outcomes that both review signoffs require before source close. The goal PASS is already recorded in `source-context/goal-check-zcode-1.md`. This artifact does not replace that fresh goal process. It adjudicates the newly available Windows evidence only.

I read in full: `source-context/ORGANIZER-BRIEF.md`, `FINAL.md`, `review/consensus.md`, `source-context/goal-check-zcode-1.md`, and the delivery files `windows-evidence-brief.md`, `windows-final-comparison.json`, `windows-comparison-first-sample.json`. I read the four raw logs named by the comparison files. I also read the `windows-ci-watch.log` tail.

I changed no product file, signoff, FINAL, or credential. I did not merge, close, publish, install, apply legacy history, prune, or declare any worktree. I created this file only.

## Protocol context attestation

The launcher supplied one private prompt file. I hashed it before any task action.

```text
/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100/.parley-runtime/invocations/38c2fe70-18c1-4a3a-97da-ad6ec684232d/headless-prompt-3601243304.md
SHA256 99276d1d8b3c8b5ed21944f9c6d888e8451c0e506bd4192b70a21d97855317af
```

The hash matched the launcher value. I then read the complete file, all 1641 lines. Its attestation, verbatim:

```text
Protocol context attestation: {"context_mode":"full","source_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec","packet_sha256":"aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec"}
Shadow packet audit: {"packet_sha256":"9d8b2219d8a8a949761c9f46deb359f6471f38aecd165c9cbb7ce45527203f1f","packet_bytes":81966,"source_bytes":130042,"included_blocks":37,"omitted_blocks":33}
```

Mode `full`. Renderer reachable. No fallback reason. `Transport: github-pr`. The shadow audit line is an unapplied diagnostic. The deck copy hash equals the attested source hash (measured below).

## Tree identity and product bytes

Commands and output:

```text
git rev-parse HEAD          bcf1140ff64175dc21e89386a0c0bc8491c0ea72   (40 chars)
git diff --name-only 99d04b3..HEAD
  → 16 files under parley-deck/ideas/<slug>/ (IMPLEMENTATION, review/, source-context/),
    parley-deck/inbox/codex-1-to-all_<slug>_goal-launch.md,
    parley-deck/runs/20261013T1323…/events.jsonl, parley-deck/runs/20261013T1354…/events.jsonl
  → no product file
shasum -a 256 parley-deck/COOPERATION.md
  aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
```

The actual source HEAD is `bcf1140ff64175dc21e89386a0c0bc8491c0ea72`, as the launch brief states. One commit follows the reviewed commit `99d04b3`. That commit adds idea audit artifacts only. Product bytes remain at the reviewed `99d04b3`. The deck protocol copy matches the attested live source hash. Verdict: CONFIRMED (PRIMARY).

## Terminal Windows outcomes

Live state, fetched in this process:

```text
gh run view 38406035407 --repo feci/parley-deck-cli --json headSha,status,conclusion
  {"conclusion":"failure","headSha":"99d04b3ecd0740fb834e70ec109c9384c7b7cd4d","status":"completed"}
gh run view 38406028243 --repo feci/parley-deck-cli --json headSha,status,conclusion
  {"conclusion":"failure","headSha":"99d04b3ecd0740fb834e70ec109c9384c7b7cd4d","status":"completed"}
gh run view 37938864253 --repo feci/parley-deck-cli --json headSha,status,conclusion
  {"conclusion":"failure","headSha":"24fe38fd150af84147408968387816a9c0145562","status":"completed"}
```

Job matrix for the two current runs at `99d04b3`:

```text
run 38406035407  ubuntu-latest success 13:38:20Z; macos-latest success 13:46:01Z; windows-latest failure 14:20:58Z
run 38406028243  ubuntu-latest success 13:38:38Z; macos-latest success 13:51:44Z; windows-latest failure 14:21:02Z
gh api repos/feci/parley-deck-cli/actions/jobs/115325766440 → run_id 38406035407, conclusion failure
gh api repos/feci/parley-deck-cli/actions/jobs/115325741335 → run_id 38406028243, conclusion failure
```

I hashed all four named logs. Each hash equals the recorded value:

```text
logs/windows-99d04b3-115325766440.log   7251b0af013a82e61ec15938d9faa382ab5ab719b69fb860d7224804227d031e
logs/windows-99d04b3-115325741335.log   dfdd243f8a346c0e9d36c66098f00a264b9be2b214be4bff801abf6d95aa60b7
logs/windows-main-24fe38f.log           74e64781718ad83c87375a071c2d022605bd811d0d0fffe865eb408d5c66a715
2026-10-09-driver-unstall/logs/windows-fixup-pr.log
                                        ce898a3c02d75bf00e8a2ff7d61748aedce15754b7fa52df7a4b3450e9fb5f43
```

Both current Windows jobs completed as failures. Ubuntu and macOS passed at the same commit. Verdict: CONFIRMED (PRIMARY).

## Failure-name comparison

I extracted the failing names myself with the recorded pattern `--- FAIL: ([^\s]+)`. My independent counts:

```text
older pre-change log (windows-fixup-pr.log, product 134ac40)   356 names
latest pre-change main (windows-main-24fe38f.log)              356 names
union of both baselines                                        357 names
current job 115325766440                                       357 names
current job 115325741335                                       357 names
```

Set algebra, computed in this process:

- Added versus the older baseline, both jobs: `TestSeparateProcessesCannotExceedActionOrCostCap/false` only.
- Removed versus the older baseline, both jobs: none.
- Added versus the union of both baselines, both jobs: none. No current failure name is absent from the union. Each current set equals the union exactly.
- Added versus latest main `24fe38f`: `TestSeparateProcessesCannotExceedActionOrCostCap/true` only.

Budget test variants per log:

```text
older baseline   TestSeparateProcessesCannotExceedActionOrCostCap, …/true
main 24fe38f     TestSeparateProcessesCannotExceedActionOrCostCap, …/false
both current     TestSeparateProcessesCannotExceedActionOrCostCap, …/false, …/true
```

This reproduces every field of `windows-final-comparison.json`. The first-sample JSON carries the same per-job fields without the main-baseline union. Verdict: CONFIRMED (PRIMARY).

## Timeout and evidence-package compile failure

Each of the four logs contains exactly one timeout panic and exactly one evidence build failure:

```text
panic: test timed out after 45m0s
running tests:
	TestDropoutConcurrentAndPrivatePartial (44m40s)
##[error]internal\evidence\tree_report_test.go:97:20: undefined: syscall.Mkfifo
FAIL	parley-deck-cli/internal/evidence [build failed]
```

All four logs name the same hung test with the same 44m40s. The compile error is Unix-only `syscall.Mkfifo` in an evidence test. Both defects pre-date this change. All four logs also list the same 18 failing packages, including `internal/runner 2700.0xx s` at the 45-minute ceiling. Verdict: CONFIRMED (PRIMARY).

## Budget package boundary

The organizer claims the whole `internal/budget` package is unchanged. I verified the bytes and read the dependency boundary myself.

```text
git diff --name-only 24fe38f..99d04b3 -- internal/budget/   → 0 files
git diff --name-only 134ac40..99d04b3 -- internal/budget/   → 0 files
git diff --name-only 134ac40..d51e1c7 -- . (excl ideas/inbox/runs) → 0 files
go list -deps -test ./internal/budget | grep '^parley-deck-cli/'
  internal/fsutil, internal/procctl, internal/pidlease, internal/quota, internal/telemetry
```

None of those five dependency packages changed between `24fe38f` and `99d04b3`. The changed product packages are `internal/app`, `internal/protocol`, `internal/protocolpacket`, `internal/runner`, plus top-level release files. None enters the budget closure. So the budget test binary builds from byte-identical sources at both commits. Verdict: CONFIRMED (PRIMARY).

I read the test at `internal/budget/review_test.go:378-440`. It loops `cost := []bool{false, true}` and names the subtests `false` and `true`. Each subtest starts 12 child processes. Each child calls `Reserve` under a five-admission cap. A cap refusal exits with code 23. Line 438 is the child's `t.Fatal(err)` when `Reserve` fails with a different error. The parent then fails at line 409 with "process/acquisition failures prevent evaluating the admission count".

## File-sharing error at review_test.go:438

The same file-sharing error at `review_test.go:438` already exists in both historical variants:

```text
older baseline  …TestSeparateProcessesCannotExceedActionOrCostCaptrue3965516569\001\lock-origin:
                 The process cannot access the file because it is being used by another process.
main 24fe38f    …TestSeparateProcessesCannotExceedActionOrCostCapfalse1872208926\001\lock-origin:
                 The process cannot access the file because it is being used by another process.
```

Both current jobs report the same error under both subtests. One nuance differs from the historical form. In job `115325766440`, the `/true` failure cites a global lock path:

```text
review_test.go:438: open C:\Users\runneradmin\AppData\Local\parley\budget-locks\b60412e25cda175fdafa414341aec3c54b9824e022756e038c77293e948463ce.lock:
The process cannot access the file because it is being used by another process.
```

Same test line, same sharing-violation class, same unchanged `Reserve` path, different file. So the historical error exists in both variants, and the causes are similar in class but not byte-identical. I record this as a finding below. Verdict: CONFIRMED (PRIMARY), with the path variance stated.

## New policy and transport test coverage on Windows

The nonverbose logs do prove execution for packages that completed. Package summary lines give three groups:

1. `internal/protocolpacket` completed `ok` (5.507 s and 4.6 s) in both current jobs. Its changed test file carries the new on-demand appendix map test. So `TestLiveSTEAppendixIsOnDemandAndItsPointerIsRetained` ran and passed on Windows.
2. `internal/app` completed with failures (4.198 s, 4.659 s) and no timeout. All its tests ran. `TestOutputTextPolicyAcrossDriverPrompts` has no failure line. So it ran and passed on Windows.
3. `internal/runner` hit the 45-minute timeout. I verified the run order with `go test ./internal/runner/ -list '.*'`:

```text
#26  TestDropoutConcurrentAndPrivatePartial      (sequential; no t.Parallel in its body)
#66  TestInteractiveArgumentDeliveryAndBound
#76  TestOutputTextPolicyAcrossRunnerPrompts
#89  TestLargeHeadlessPromptReachesRealChildWithExactAttestedBytes
#90  TestHeadlessArgumentBoundsPreserveShortInputsAndRefuseUnsafeDelivery
#91  TestLargeStdinRemainsExactWithInheritedPromptArgument
#92  TestLargePromptDeliveryCoversEveryBuiltinArgumentTemplate
```

The sequential test at position 26 hung for 44m40s. The binary then panicked. Every test after position 26 never ran. That includes the new runner policy test and all four transport tests. All four logs, current and historical, hung the same way. The new test files contain no `t.Skip` or Windows build guard.

So Windows CI supplies no execution evidence for the runner-side policy test and the private-file transport tests. This is a concrete coverage gap, not a failing result. Verdict: CONFIRMED (PRIMARY).

## Adjudication against AC7 and the review consensus

AC7 requires recorded test, build, and skill evidence, and disclosure of any unrelated CI failure. The local full suite, build, vet, ubuntu, and macOS outcomes stand as recorded. The terminal Windows outcome is failure, and it matches the disclosed pre-change baseline in every part I checked: same 18 failing packages, same timeout test, same evidence compile error, and no failure name outside the two-baseline union. No new Windows defect is attributable to this change. The evidence therefore does not defeat AC7 or the zero-fix review consensus. It confirms the disclosed baseline shape.

New findings from this check, reported openly:

- [MINOR] The new runner policy test and the four transport tests never executed on Windows in any of the four runs. The 45-minute hang at test 26 precedes them in run order. Windows CI proves nothing about the transport repair. Linux, macOS, and local darwin evidence carry that criterion.
- [NIT] The `/true` budget failure in job `115325766440` cites the global `budget-locks` path, not the per-test `lock-origin` path seen in both historical logs. Same line and error class. The two-baseline union absorbs both variants either way.
- [NIT] The prior goal report's tree-identity bullet repeats `c0` in the full SHA. See the correction below.

## Signed prerequisite

Both review-consensus signoffs require the final Windows outcomes before source close. The brief's four close conditions already hold through the recorded goal PASS. The final Windows outcomes now exist and are recorded: both current jobs failed with the pre-existing timeout, the pre-existing evidence compile failure, and baseline failure names only. The completed Windows evidence satisfies the signed prerequisite. Native Windows validation remains unresolved and experimental. The CLI WinGet publication stays held, exactly as FINAL, the consensus, and the comparison file state. "No new names" does not establish full Windows correctness, complete coverage, or identical causes. This check claims none of those.

## Correction to the prior goal report

The tree-identity bullet in `goal-check-zcode-1.md` line 43 prints a 42-character SHA:

```text
bcf1140ff64175dc21e89386a0c0c0bc8491c0ea72   (bullet, line 43; repeats c0)
```

The frontmatter line 7 holds the correct 40-character value:

```text
current-cli-head: bcf1140ff64175dc21e89386a0c0bc8491c0ea72
```

I verified the correct value independently in this process with `git rev-parse HEAD`. It returns `bcf1140ff64175dc21e89386a0c0bc8491c0ea72`. The bullet's other facts hold: prefix `bcf1140`, branch `ste100`, one commit after `99d04b3`, no product byte changed. The typo is isolated. It changes no measurement, verdict, or close condition in that report. I keep the original file unchanged, as the launch brief directs. This section is the correction of record.

## Commands run

```text
shasum -a 256  <private prompt file; four Windows logs; parley-deck/COOPERATION.md>
git rev-parse HEAD; git log --oneline -8; git status --short
git diff --name-only 99d04b3..HEAD
git log --oneline -1 24fe38f; git log --oneline -1 134ac40
git diff --name-only 24fe38f..99d04b3 -- internal/budget/
git diff --name-only 134ac40..99d04b3 -- internal/budget/
git diff --name-only 134ac40..d51e1c7 -- . ':(exclude)parley-deck/ideas' ':(exclude)parley-deck/inbox' ':(exclude)parley-deck/runs'
git diff --name-only 24fe38f..99d04b3 -- . ':(exclude)parley-deck'
go list -deps -test ./internal/budget
go test ./internal/runner/ -list '.*'
python3  (failure-name extraction and set algebra over the four logs)
python3  (per-line SHA comparison in goal-check-zcode-1.md against git HEAD)
grep -c 'panic: test timed out' / 'internal/evidence [build failed]' / 'review_test.go:438'  (all four logs)
grep -E 'FAIL|ok' package summary lines  (all four logs)
gh run view 38406035407 / 38406028243 / 37938864253 --repo feci/parley-deck-cli --json …
gh api repos/feci/parley-deck-cli/actions/jobs/115325766440 / …/115325741335
```

WINDOWS-EVIDENCE: CONFIRMED — scoped to the organizer's recorded Windows comparison and the signed prerequisite. I reproduced every recorded field from the raw logs myself: the two current `99d04b3` Windows failures, the 356/356/357/357 name counts, the single added name per baseline pair, the empty union difference, the pre-existing timeout and evidence compile failure, and the unchanged `internal/budget` package with its clean dependency boundary. The completed evidence satisfies the signoffs' Windows prerequisite without defeating AC7 or the zero-fix consensus, while Windows stays experimental and CLI WinGet stays held. This verdict does not certify Windows correctness or coverage; the runner-side policy and transport tests never ran there, as recorded above.
