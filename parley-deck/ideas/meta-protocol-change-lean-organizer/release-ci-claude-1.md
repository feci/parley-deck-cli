---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent release incident reviewer
artifact: release-ci
date: 2026-09-24
verdict: FAIL
---

# Release CI incident review — claude-1

## Verdict

**FAIL.** The release channel is not healthy. Two of three hosted CI legs failed on
source byte-identical to the published tag `v1.49.0`, and one published claim in the
shipped `CHANGELOG.md` is falsified by the very run that was supposed to back it.

This verdict is about the **channel and its claims**, not about the A–D feature work,
which I did not re-review here. No source, tag, release or workflow was changed by me.

## Scope and method

Evidence is the hosted GitHub Actions logs pulled with `gh` (run `35983818751`, repo
`feci/parley-deck-cli`), plus local reproduction in a clean clone at tag `v1.49.0` and
in Linux containers. `gh run view --log` refuses while a run is in progress, so per-job
logs were fetched from `repos/.../actions/jobs/<id>/logs`.

Run `35983818751` is now **complete** (it was still in progress at dispatch):

| leg | job | result |
|---|---|---|
| ubuntu-latest | 107581701569 | **failure** — `Test` step, exit 1 |
| windows-latest | 107581701326 | **failure** — `actions/checkout@v4`; setup-go/Build/Test **skipped** |
| macos-latest | 107581701420 | success |

The dispatch brief stated "macOS/Linux were still running". Both have since finished:
**Linux failed**, macOS passed. The Windows failure is *not* in build or test — that leg
never reached them.

Source tested vs source released: `git diff v1.49.0 137b1c1 -- . ':(exclude)parley-deck'`
is **empty**. The CI result applies directly to the released source.

Skill channel: `feci/parley-deck-skill` "Tests" (`35983815600`) and "Release portable
binaries" (`35983892815`, v2.13.0 @ `8161e5e`) are both **success**. I found no skill-side
defect.

## Finding A — Linux `internal/app`: 15 subtests, deterministic, root-caused

`TestEvidenceVerifierProductionClosure` — all 15 subtests failed. 14 at
`evidence_verify_test.go:231` ("prior verification refusal requires exact-record recovery
and fresh checks"); `report-persistence-failure-recovery` at `:281` ("canonical refusal is
not verified committed: could not commit the exact refusal record; retained recovery is
required"). The split is itself diagnostic: `:230` skips the committed-refusal check for
exactly that one subtest, so it surfaces the same underlying failure one layer deeper.

**Root cause (product code, not test scaffolding).** `commitVerificationRefusal`
(`internal/app/evidence_refusals.go:83`) runs `git commit` through `refusalGit`, which sets
only `GIT_OPTIONAL_LOCKS=0` on top of `os.Environ()` and therefore depends entirely on an
**ambient git author identity**. The fixture `gateScratchRepo`
(`internal/app/driver_evidence_test.go:24`) runs `git init` and then commits with
*per-command* `-c user.email=… -c user.name=…`, which never persists into repo config. So
the scratch repo has a `HEAD` (hence `git ls-tree HEAD` succeeds and `committedRefusal`
returns `(false, nil)` rather than an error) but no configured identity. On the dev machine
git auto-detects an identity from a domain-qualified hostname and the commit succeeds; on a
stock hosted runner no usable identity is available and the commit fails, so the refusal
never lands in `HEAD` and every subtest reports it as uncommitted.

**Reproduced and fix-validated** (clean clone at `v1.49.0`):

- macOS, identity available → `ok internal/app 29.486s`.
- macOS, identity suppressed (`user.useConfigOnly=true`) → identical signature: same
  subtests, same lines 231/281, same messages.
- Linux container, non-root user, full repo, **no identity** → identical CI signature,
  including the `:281` variant.
- Same container, **with** `user.email`/`user.name` set → `ok internal/app`,
  `ok internal/trajectory`. **The fix clears it.**

**Provenance — preexisting, newly exposed, not introduced by this idea.**
`commitVerificationRefusal` was added **2026-09-11** in `9cef3fc`, under the earlier idea
`meta-protocol-change-evidence-first-efficiency`. It was never exposed because
`.github/workflows/tests.yml` did not exist: the workflow was added by *this* idea
(`86d028b`, amended `998346c`), and run `35983818751` is the **only Actions run in the
repository's entire history** (`actions/runs` `total_count` = 1; workflow registered
2026-09-24T09:51:44Z). This is CI's first execution, so it is a newly-revealed preexisting
condition, not a regression.

**Effect on released artifacts — real but narrow, and fail-closed.** This is product code,
reachable two ways: `retainDriverRefusal` (invoked from `VerifyCompletionEvidence`,
`internal/app/evidence_verify.go:362`) and the user-facing
`parley evidence refusals recover` (`internal/app/evidence_refusals.go:204`). On a host
with no resolvable git identity — containers, CI images, minimal images — the recover
command exits 1, and driver refusals degrade to "refusal retained; use evidence refusals
inspect/recover", a state that cannot then be cleared on that same host because recovery
uses the same commit path. Two mitigating facts: the path is only reached **when a
verification refusal has already occurred** (not on clean runs), and it is **fail-closed** —
it blocks completion rather than fabricating a pass. The released binaries are not corrupt
and behave correctly on ordinary developer machines, which have a git identity.

## Finding B — Linux `internal/trajectory`: 1 subtest, NOT reproduced, unresolved

`TestRecoveredParentResolutionRequiresExactLineage/missing-replacement-terminal` failed at
`verification_parent_recovery_test.go:142`: "recovered helper execution did not complete",
with `Steps:1` (4 expected) and `FailureStage:execution`, from
`internal/trajectory/captured.go:423`.

This is a **second, independent cause** — it is not the git-identity failure. I attempted
reproduction four ways and it passed every time: macOS with identity, macOS without, Linux
container in isolation (all 11 subtests pass, `missing-replacement-terminal` in 0.47s), and
the faithful Linux container arm with no identity. I therefore did **not** establish its
root cause and will not assert one.

What the code constrains it to: `Complete` (`internal/evidence/execute.go:275`) is false if
output overflowed, the envelope or go-test stream was invalid, `ctx.Err() != nil`, or — for
controlled executions — `ExitCode >= 128`, or the run error is not an `*exec.ExitError` with
a non-negative code. The latter two are the signal-termination cases. A child killed by a
signal under runner resource pressure fits the observed record, and CI was markedly slower
than any local run (`internal/app` 160.7s and `internal/trajectory` 171.6s on CI, against
5–6s per package in my container), which is consistent with contention. That is a
**hypothesis, not a finding** — one observation, no reproduction.

Provenance and effect: undetermined. Because it is unreproduced, I cannot state whether it
is deterministic on hosted Linux or load-dependent, and I will not certify it either way.

## Finding C — Windows: checkout fails, so the Windows leg has never run

```
error: unable to create file parley-deck/ideas/meta-protocol-change-evidence-first-efficiency-v2/
source-context/runtime-cap-evidence/parley-deck__ideas__meta-protocol-change-evidence-first-
efficiency__implementation-notes__codex-1-kimi-cap-controls-20260916.json: Filename too long
```

That repository path is **233 characters**; the runner workspace prefix
`D:\a\parley-deck-cli\parley-deck-cli\` is **37**; 233 + 37 = **270**, over the 260-character
`MAX_PATH` limit. `core.longpaths` is not configured anywhere in `.github/`. It is the
**only** path in the tree over 200 characters (longest = 233), so this is a single-file
problem.

**Provenance — preexisting, from an earlier idea.** Added **2026-09-16** in `5ecac7b` under
`meta-protocol-change-evidence-first-efficiency-v2`. Not introduced by lean-organizer;
newly exposed by the new workflow.

**Effect on released artifacts.** The released Windows binaries are **not** affected — this
is a clone/checkout-time limit on repository content, not a defect in the compiled artifact.
The real damage is evidentiary, and it is the reason this matters for delivery:

## Finding D — a shipped claim is falsified

`CHANGELOG.md` at the released tag states: *"Windows behavior of `wait`/`usage` is exercised
by the CI leg; local verification was macOS-only."* The Windows leg failed at checkout with
`setup-go`, `Build` and `Test` all **skipped**. No Windows test has ever executed for this
project. The workflow comment calls the Windows leg "load-bearing … so releases ship Windows
assets honestly" — and Windows assets shipped with zero Windows execution behind them.

So `v1.49.0` ships Windows binaries whose only stated evidence does not exist. Given this
idea's subject matter, that is the finding I would least want left unrecorded.

## Process observation

The GitHub Release `v1.49.0` was published at **09:52:11Z**, 21 seconds after the CI run
started and ~3.5 minutes before the Linux leg failed. Publication did not gate on CI. The six
release binaries were not produced by Actions either — the Tests run is the repository's only
run ever — so they were built and uploaded outside the hosted channel.

## Recommended remediation — minimal, and none of it moves a released tag

Nothing here requires changing released source, so **no new version is needed for A or C**;
`v1.49.0` must not be re-pointed.

1. **Finding A (CI-side, no source change).** Add an identity step to
   `.github/workflows/tests.yml` before `Test`:
   `git config --global user.email "ci@example.invalid" && git config --global user.name "CI"`.
   Validated above: it turns the failing packages green. *Tradeoff to decide, not for me to
   pick silently:* this fixes the channel but leaves the product's dependence on an ambient
   git identity intact. If `parley` is meant to run in containers/CI, harden
   `refusalGit` to supply an explicit fallback identity for its own internal commits — that
   **is** a source change and belongs in a **new version (v1.49.1)**, never in a moved tag.
2. **Finding C (CI-side).** Add `git config --global core.longpaths true` as a step
   *before* `actions/checkout@v4` on the Windows leg. I could not verify this on a hosted
   runner without pushing a workflow, which is outside my remit — so treat it as unverified.
   The guaranteed fallback, if longpaths proves insufficient, is to shorten that single
   233-character record path; it is deck-record content, not source.
3. **Finding B.** Re-run the Linux leg (`gh run rerun 35983818751 --job 107581701569`) to
   classify it as flaky or deterministic **before** anyone calls this channel green. Do not
   record it as resolved on the strength of my non-reproduction — I reproduced Finding A and
   could not reproduce this one, and those are different states.
4. **Finding D.** Do not move `v1.49.0`. Correct the Windows-coverage sentence in the next
   version, and — owner/organizer decision, not mine — consider amending the GitHub Release
   *notes* body, which is mutable metadata independent of the immutable tag. Until a Windows
   leg actually executes, the honest statement is that Windows is built but untested.

## Recommended next step

**Do not treat the release channel as delivered.** Apply remediation 1 and 2 to the workflow,
re-run the full matrix, and only then re-assess. The blocking question for delivery is
Finding D: `v1.49.0` is public with a coverage claim its own CI disproves, and that should be
corrected in the next version rather than by touching the published tag.

## Limitations

- Finding B is unresolved; I state a hypothesis and label it as such.
- The Windows `core.longpaths` remedy is reasoned from the log and path arithmetic, not
  executed on a hosted Windows runner.
- I reviewed the release channel, not the A–D feature implementation.
- Container reproduction used `linux/arm64`; CI is `linux/amd64`. Finding A reproduced
  identically regardless, so the architecture difference does not affect it.
- No source, tag, release, workflow or branch was modified. No credentials or log secrets
  are included here.
