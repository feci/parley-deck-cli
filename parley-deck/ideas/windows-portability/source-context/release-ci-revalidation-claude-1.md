---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent release incident reviewer (revalidation)
artifact: release-ci-revalidation
date: 2026-09-24
verdict: FAIL
---

# Release CI revalidation — claude-1

## Verdict

**FAIL.** Hosted run `35987916696` (commit `4e0c069`, branch `main`) concluded
**failure**: macOS success, **ubuntu failure**, **windows failure**.

The workflow fix is sound and did exactly what it claimed for its two targets, but the
channel is not green, and the Windows coverage claim in the shipped `CHANGELOG.md` is now
**falsified by execution** rather than merely unsupported. No green hosted matrix exists,
so nothing here discharges that claim.

I changed no code, workflow, tag, release or branch, and triggered no run.

## What I verified myself

**Run under review.** `35987916696`, workflow `Tests`, event `push`, head
`4e0c0699f445ef64b65cb4f681834ffce4353b08`, started 2026-09-24T10:33:29Z, watched to
completion with one blocking `gh run watch --interval 60` (watch exit status 1 = run
failed). Go `1.26.8` on all three legs.

| leg | job | conclusion | window |
|---|---|---|---|
| macos-latest | 107594829995 | **success** | 10:33:36 → 10:48:11Z |
| ubuntu-latest | 107594830258 | **failure** | 10:33:32 → 10:36:40Z |
| windows-latest | 107594830296 | **failure** | 10:33:31 → 10:37:50Z |

Per-job logs were pulled from `actions/jobs/<id>/logs` (`gh run view --log` refuses while
any leg of the run is still in progress).

**The workflow does not suppress failures — PASS.** `git diff 137b1c1 4e0c069` is two
files: `.github/workflows/tests.yml` (+17, a pure insertion, no existing line touched) and
`release-ci-fix-zcode-1.md`. The insertion is two `git config` steps placed before
`actions/checkout@v4`. Grepping the resulting file for suppression constructs returns
exactly two hits: `fail-fast: false` (pre-existing, and it *widens* reporting) and
`if: runner.os == 'Windows'`, which gates only the `core.longpaths` config step — no test
step is conditional. There is no `continue-on-error`, no retry, no `-run`/`-short` filter,
no `|| true`, no output redirection. `Build` (`go build ./...`) and `Test`
(`go test ./... -count=1 -timeout 45m`) are byte-identical to the incident run.

**Product source equals the released tag — PASS.** `git diff --name-only v1.49.0 4e0c069`
is four paths: three under `parley-deck/` and `.github/workflows/tests.yml`. With
`parley-deck` and `.github` excluded the diff is **empty** — the entire Go module,
`go.mod`, `VERSION`, `CHANGELOG.md` and `internal/app/version.go` are byte-identical to
`v1.49.0`. Tag `v1.49.0` is `06e563e8b1fe8094132149b83e14da7be9aae51e` both locally and at
`origin` (`git ls-remote --tags`); it has not moved. `origin/main` is `4e0c069`.
**The failures below are therefore failures of the released source**, not of a variant.

## Finding A (git identity) — RESOLVED on Linux, with a caveat

The identity step ran (`git config --global user.email "parley-ci@example.invalid"` /
`user.name "parley-ci"`, logged on every leg). `internal/app` on ubuntu is now
**`ok 136.801s`**, and `TestEvidenceVerifierProductionClosure` appears **zero** times in
the log. Against two prior hosted attempts that failed it with an identical 15-subtest
signature (14× `evidence_verify_test.go:231`, 1× `:281`), this is a clean fix, hosted-verified.

**Caveat, unchanged from my first report.** The fix greens the channel by giving the runner
what the product lacks. `refusalGit` (`internal/app/evidence_refusals.go:57`) still extends
`os.Environ()` with only `GIT_OPTIONAL_LOCKS=0` and commits with no identity of its own, so
CI can no longer detect that dependence. The product is internally inconsistent about this:
`internal/runner/reviewsnapshot.go:139` already pins `-c user.name=parley
-c user.email=parley@localhost` for its own internal commit. Making `refusalGit` match that
existing precedent is a **source change needing its own version** (v1.49.1), not a tag move.

## Finding C (Windows checkout) — RESOLVED, and no longer a reasoned guess

`core.longpaths true` before checkout worked on a hosted runner. `Run actions/checkout@v4`,
`actions/setup-go@v5` and `Build` all **succeeded** on windows-latest for the first time in
the repository's history. My earlier remedy was explicitly flagged as unverified; it is now
verified by execution. This is the one unambiguous gain of this run.

## Finding D (shipped Windows claim) — NOT discharged; now falsified by execution

`CHANGELOG.md` at `v1.49.0` states: *"Windows behavior of `wait`/`usage` is exercised by the
CI leg; local verification was macOS-only."* The Windows leg has now actually run:

- **13 packages FAIL, 17 ok**, **157 distinct failing tests**, plus one package that does
  not compile.
- `internal/app` **FAIL in 3.499s** — it did not fail slowly, it **aborted**:
  `panic: interface conversion: interface {} is nil, not map[string]interface {}` at
  `internal/app/app_test.go:155`, in `TestVersionAllUsesDirFlagForProjectStatus`. A panic
  kills the whole package test binary.
- `wait_test.go` (21 test funcs) and `usage_ingest_test.go` (6) live in that same package
  and sort *after* `app_test.go`. **No `wait` or `usage` test result appears anywhere in the
  Windows log.** Neither verb has ever executed on Windows.

So the claim is not "unsupported pending evidence" — the evidence now exists and
contradicts it. For comparison, that package runs 136.8s on ubuntu and 591.1s on macOS.

I must also record that **my own release preflight was wrong here**. `release-preflight-claude-1.md`
§4 asserted the Windows claim was "true" on the grounds that the matrix lists
`windows-latest`. That reasoned from configuration rather than execution, and CI's first two
runs disproved it. The matrix entry was never evidence.

## Finding E (NEW, root-caused) — `internal/acp` loses child stderr: a Wait-before-drain race

ubuntu: `--- FAIL: TestSpawnObservesAllStderrAndRealExit (0.01s)` /
`spawn_test.go:38: observed bytes: 0`.

The test re-execs the test binary as a child that writes 16384 bytes to stderr and exits 7.
The exit code assertion **passed** (7 was read correctly); only the byte count was 0. So the
child ran and wrote — the parent lost the bytes.

**Root cause (product code).** `internal/acp/spawn.go` is the only product code using
`cmd.StdoutPipe()`/`cmd.StderrPipe()` (lines 56/61). The stderr copier runs in a goroutine
tracked by `p.wg` (line 82). Both `Process.Stop` (line 123, via line 128) and `Process.Wait` (lines 144–145)
call `p.cmd.Wait()` **first** and `p.wg.Wait()` afterwards. `os/exec` documents that `Wait`
closes the pipes returned by `StderrPipe` as soon as the command exits, and that it is
therefore incorrect to call `Wait` before all reads from the pipe have completed. When
`Wait` wins the race the pipe is closed under the copier and `io.Copy` returns having
observed as little as **0 bytes** — exactly what CI recorded. A loaded 2-core hosted runner
loses this race far more often than a dev machine.

**Provenance:** preexisting. `internal/acp/spawn.go` was last touched by `5712a36`
(`meta-protocol-change-evidence-first-efficiency`), not by this idea.
**Narrow remedy:** drain `p.wg.Wait()` before `p.cmd.Wait()` in both `Stop` and `Wait`.
Source change → **v1.49.1**, not a tag move.
**Scope limit — this does NOT explain the other failures.** `internal/evidence/execute.go`
assigns `cmd.Stdout`/`cmd.Stderr` to an `io.Writer` (lines 146–147), which `os/exec` drains
inside `Wait` itself. The trajectory/runner failures below run through that safe path.

## Finding B family (trajectory/runner) — STILL UNRESOLVED, now identity-independent

ubuntu, with the identity fix active, produced three more failures — a **third** distinct
set, no test repeating across runs:

| package | test | detail |
|---|---|---|
| `internal/runner` | `TestVerifierBudgetRefusalRecoversExplicitly` | `verification_refused_recovery_test.go:229`, `Steps:1`, `FailureStage:execution` |
| `internal/trajectory` | `TestParentRecoveryPublicationInterruptionAndExactReplay/after-publication` | `parent_recovery_test.go:124`, `helper: exit status 19` |
| `internal/trajectory` | `TestRecoveredParentPublicationRefusesNonRefusalOriginal` | `verification_parent_recovery_test.go:142`, `Steps:1`, `FailureStage:execution` |

All carry the same message: *"captured criterion execution is incomplete; retain partial
observations."* Across three hosted ubuntu executions this area has now failed with six
different test names and never the same one twice; it has never reproduced locally
(my four attempts, zcode-1's controlled runs), and macOS ran the same packages green
(`internal/trajectory ok 809.899s`, `internal/runner ok 132.634s`).

**These are identity-independent** — confirmed, since this run had the identity present.
**Root cause is NOT established, and I will not call them flaky.** "Flaky" is a claim about
cause, and I have none: what I have is a repeated *shape* (a helper's captured execution
stops short — 1 of 4 steps, 3 of 4 steps, or no observed terminal) with no mechanism
identified. Finding E is a genuinely analogous bug in a *different* package, which raises but
does not establish the possibility of a second drain/observation defect on the
`internal/evidence` path; I checked that path and its stdio wiring is the safe kind, so
Finding E's mechanism is ruled out, not extended.

**Narrow next diagnostic** (not a fix, and not for me to run): re-run only the affected
packages on hosted ubuntu with `-v` and the helper's own stderr retained, so the child's
exit path is visible. `helper: exit status 19` is the single most specific datum available
and is worth tracing before anything else.

## Finding F (NEW) — what the first real Windows execution exposed

Beyond the aborted `internal/app`, the Windows failures cluster into a small number of
root-identifiable classes. These are portability defects, not noise:

1. **69 failures, one product line.** All carry *"snapshot store must be a private real
   directory"* from `internal/trajectory/snapshot.go:163`. The guard at line 162 is
   `info.Mode().Perm()&0077 != 0`. Go on Windows synthesizes mode bits and reports
   directories as `0777`/`0555`, so that expression is **always** non-zero and the check can
   never pass. The identical pattern recurs at `snapshot.go:544` for regular files (2 further failures, *"capture certified hidden Go input … snapshot store must be a private real directory"*), so a fix must cover both lines. This is **product source**, not test scaffolding: trajectory snapshot capture
   in the shipped Windows binaries is expected to refuse for this reason.
   *(Established from the code plus 69 identical hosted failures; I have no Windows host and
   did not run the binary there.)*
2. **`internal/evidence` does not compile on Windows.**
   `internal/evidence/tree_report_test.go:97` calls `syscall.Mkfifo`, which is undefined on
   Windows, and the file carries no build tag → `FAIL ... [build failed]`. It is the only
   untagged file in the tree using a Unix-only `syscall` symbol. Added by `7f2676f`
   (`meta-protocol-change-evidence-first-efficiency`) — preexisting.
3. **The `internal/app` panic chain.** `writeFakeParleyDeckSkill` (`app_test.go:1742`)
   writes a `#!/bin/sh` script named `parley-deck-skill` with **no extension**; Windows
   cannot resolve or execute it (`exec: "parley-deck-skill": executable file not found in
   %PATH%`). The payload then lacks `parley_deck_skill`, and the unchecked assertion at
   `app_test.go:155` panics on nil, aborting the package. Line 155 predates this idea
   (`1b1764a`). Two independent defects: a non-portable fixture, and an assertion that turns
   a clean failure into a suite-wide abort.
4. **Windows file-sharing semantics.** 6 budget-lock failures plus 8 open/rename failures
   with *"The process cannot access the file because it is being used by another process"*
   (e.g. `ledger_test.go:57`) — the locking model assumes POSIX semantics.
5. **10 fixtures hard-code `/bin/sh`** as the agent command (`launch_test.go:48`, `:133`).
   Test-side only: the sole `/bin/sh` in product source is
   `internal/procctl/terminal_unix.go:40`, which is build-tagged Unix.
6. **Line endings.** `hardening_test.go:456` expects `"v2-dirty\n"` and got `"v2-dirty\r\n"`.

## macOS — green, and the only green leg

31/31 packages `ok`, **zero** `--- FAIL` lines (`internal/app` 591.103s,
`internal/trajectory` 809.899s, `internal/runner` 132.634s). macOS has now passed twice
against this source.

## Does a green hosted execution discharge the Windows claim without a new version?

**Moot for this run — no green Windows execution occurred.** Stating the answer anyway,
because the question will recur:

- A green Windows leg at `4e0c069` *would* exercise Go source byte-identical to `v1.49.0`,
  so it would be legitimate evidence for the **source**, with two limits that must be stated
  rather than assumed away: the hosted runner is **windows x64 only**, so
  `parley-v1.49.0-windows-arm64.exe` would still have zero executed evidence; and the leg now
  runs with `core.longpaths true`, which is not the Windows default.
- Even so it could not repair the **published sentence**. `CHANGELOG.md` sits inside the
  immutable tag. A later run cannot make a claim true *as of publication*, when no Windows
  execution existed. Correcting the text requires either the next version's changelog or an
  edit to the mutable release notes.

**Precise remaining blockers**, in order:

1. **The Windows leg fails**; `wait` and `usage` have never executed on Windows. Blocking for
   the claim, and — via Finding F.1 — for the honesty of shipping Windows assets at all.
2. **The ubuntu leg fails** on an unresolved, unreproduced, now identity-independent
   trajectory/runner defect family, plus the root-caused Finding E.
3. **The release notes currently contradict themselves.** The body of the published
   `v1.49.0` notes still contains "Windows behavior of `wait`/`usage` is exercised by the CI
   leg" *and* an appended paragraph saying "Windows runtime validation is not yet claimed."
   The disclosure was added without correcting the sentence it disproves.

**Narrow next remedy** (owner/organizer actions — I performed none of them):

1. **No version, mutable metadata, do first:** edit the `v1.49.0` release notes so the
   Windows sentence in the body is corrected in place, not merely contradicted later. The
   honest form is that Windows binaries are built and shipped but not exercised.
2. **`v1.49.1`, smallest source set that makes a Windows leg meaningful:**
   `internal/trajectory/snapshot.go:162` **and `:544`** (POSIX-only private-permission checks),
   `internal/acp/spawn.go` (drain before `Wait`, Finding E), and `refusalGit` in
   `internal/app/evidence_refusals.go` (explicit fallback identity, mirroring
   `reviewsnapshot.go:139`). Test-side, in the same version: a build tag or Windows fallback
   for `tree_report_test.go`, a `.bat` fixture in `writeFakeParleyDeckSkill`, and a checked
   assertion at `app_test.go:155` so a missing fixture fails one test instead of aborting a
   package.
3. **Do not call the channel green** until the trajectory/runner family has a root cause.
   That is a diagnosis task, not a rerun task.
4. **Do not move `v1.49.0`.** Nothing above requires it.

## Process observation, restated with new evidence

`actions/runs` for this repository now totals **two**. The six `v1.49.0` binaries were
uploaded at 09:52:06–09:52:09Z, before the first run finished, by the owner account — they
were not produced by Actions. Independent corroboration that they match the source exists
but is weak: my preflight cross-builds from `06e563e` produced byte sizes identical to the
published assets (`windows/amd64` 10,789,376 B; `linux/arm64` 9,568,416 B). That is a size
match, not a hash match, and it says nothing about behavior.

## Limitations

- The trajectory/runner failure family is **unresolved**. I state a shape and a next
  diagnostic, and explicitly decline to classify it as flaky or deterministic.
- Finding F.1 is established from code reading plus 69 identical hosted failures. I have no
  Windows host and did not execute the shipped binary on Windows.
- Finding E's mechanism is established from the Go `os/exec` contract, the code ordering and
  the `observed bytes: 0` record; I did not reproduce it locally, and a 0.01s failure on a
  2-core runner is not something a dev machine reproduces on demand.
- The claim that `wait`/`usage` never ran on Windows rests on the package panicking at
  3.499s and on no `wait`/`usage` result appearing in the log — not on a per-test skip list.
- I reviewed the release channel, not the A–D feature implementation.
- No source, workflow, tag, release, branch or run was created, modified or triggered by me.
  No credentials or log secrets appear here.
