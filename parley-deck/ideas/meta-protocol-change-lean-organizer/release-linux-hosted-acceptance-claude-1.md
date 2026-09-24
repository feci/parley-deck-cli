# Hosted acceptance — Linux U1/U2 release repairs (claude-1)

role: independent hosted acceptance reviewer (not author of any code under review)
date: 2026-09-24
subject: hosted run **36009912946**, branch `main`, head **9134c7a** (`ef10cf3` code + `9134c7a` record-only)
inputs: run 36009912946 job logs (ubuntu 107667649288, macos 107667648980, windows 107667649460);
prior failing runs 35987916696 / 35998165067 / 36002574211; release-linux-u2-review-claude-1.md (my
reviewed hashes); release-linux-u2-review-kimi-1.md; zcode-1 commit 519951e (comment-only, unpushed)

---

## VERDICT: **PASS** — Linux U1 and U2 acceptance met on hosted x86_64

Both originally-failing signatures are gone from a full, unsuppressed hosted Ubuntu suite, the source
that ran is byte-identical to what I reviewed, and `v1.49.0` has not moved. macOS is green. **Windows is
red and no Windows claim is made.** Release 1.49.1 is *not* unblocked by this report alone — conditions
in §8.

---

## 1. What actually ran

| leg | job | result | window | toolchain |
|---|---|---|---|---|
| ubuntu-latest | 107667649288 | **success** | 14:02:43Z → 14:06:39Z | go1.26.8 linux/amd64 |
| macos-latest | 107667648980 | **success** | 14:02:49Z → 14:19:11Z | go1.26.8 darwin/arm64 |
| windows-latest | 107667649460 | **failure** (Test step) | 14:02:45Z → 14:07:38Z | — |

Run event `push` on branch `main` at `9134c7a4e134…` — identical to the reviewed worktree HEAD at the
time of the push. (Organizer note: this was pushed to `main`, not to a `lean-organizer` remote branch;
`git ls-remote --heads origin` returns exactly one head, `main` @ 9134c7a.)

Ubuntu leg per-package times for the packages that carry U1/U2: `acp 2.633s`, `procctl 0.593s`,
`evidence 10.767s`, `runner 58.670s`, `app 156.451s`, `trajectory 165.125s`.

## 2. The original failures are gone

| run | head | original ubuntu failure | same leg now |
|---|---|---|---|
| 35987916696 | 4e0c069 | **U1** `FAIL internal/acp` — `spawn_test.go:38: observed bytes: 0` (of 16384); plus `FAIL trajectory`, `FAIL runner` | `ok acp 2.633s`, `ok trajectory`, `ok runner` |
| 35998165067 | dd0a226 | `FAIL trajectory` — `TestRecoveredParentObservationRequiresExactLineage`, `TestParentRecoveryPublicationInterruptionAndExactReplay` | `ok trajectory 165.125s` |
| 36002574211 | 231d889 | **U2** `FAIL trajectory` + `FAIL runner` + `FAIL app` carrying, 2×, `"exit_code": -1` with `"diagnostics": "run error (no command output): criterion supervisor identity unavailable: no recorded command"` | all three `ok` |

Grepped over the new ubuntu log, the count of `no recorded command` / `0 of 16384` /
`cannot read live command` / `--- FAIL` is **0**.

**These are genuine passes, not silent skips.** I inspected each previously-failing test body.
`TestVerifierBudgetRefusalRecoversExplicitly`, `TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts`,
`TestParentRecoveryPublicationInterruptionAndExactReplay` and `TestSpawnStop/WaitDrainsStderrBeforeReaping`
contain **no skip of any kind**. `TestRecoveredParentResolutionRequiresExactLineage` and
`…GuardsRevalidate` each open with exactly one guard, `if runtime.GOOS == "windows" { t.Skip("POSIX
captured verification") }` — **unreachable on the Ubuntu leg**. So for all six, a green package
necessarily means they executed and passed on Linux.

**Independent evidence the new probe code executed on x86_64:** `internal/procctl` took `0.106s / 0.105s /
0.107s` in the three pre-fix hosted runs and **`0.593s`** here. The ~0.49s increase is consistent with the
new bounded-poll suite (two tests that must spend the full 100 ms bound, plus 40 sequential + 4×concurrent
spawns). The decisive U2 regression test
`TestLinuxCaptureOfStartedProcessRecordsCommandEveryTime` has **no skip guard** and therefore ran.

## 3. Suppression audit — clean, with one stated limit

- Test step ran `go test ./... -count=1 -timeout 45m`, byte-identical to `.github/workflows/tests.yml:52`,
  under `shell: /usr/bin/bash -e`. No `-short`, no `-run` filter, no `|| true`, no `continue-on-error`,
  no retry, no rerun-on-failure anywhere in the workflow.
- Build was a separate successful step (`go build ./...`).
- **Package coverage complete:** 31 `ok` + 1 `[no test files]` (`cmd/parley`, a main package) = **32**,
  exactly `go list ./...` = 32. Set-differencing repo packages against the log leaves the empty set — no
  package silently absent.
- Zero `FAIL`, `panic:`, `DATA RACE`, or `test timed out` in the log.
- The only skipped *step* was "Enable long paths (Windows)", correctly gated by `if: runner.os == 'Windows'`.

**Stated limit (residual):** non-verbose `go test` does not print `--- SKIP`, so the log cannot enumerate
skipped subtests. Two new procctl tests carry environment-conditional skips —
`TestLinuxZombieCmdlinePollsBoundThenFailsClosed` (:217) and
`TestLinuxAttributedRefusesIdentityChangesUnderPollingProbe` (:268), both `identity probe unavailable in
this environment`. If those skipped, the log would look identical. They are *corroborating* tests; the
race regression test and the fail-closed bound test are unguarded, so the U2 claim does not rest on them.

## 4. Source identity and tag immutability — both hold

- At the hosted head `9134c7a`, `internal/procctl/procctl_linux.go` = `a3f8eb16067a…5840547` and
  `internal/procctl/procctl_cmdline_linux_test.go` = `5b3cca33da5f…2935b26` — **exactly** the two hashes
  recorded in `release-linux-u2-review-claude-1.md`. The hosted run tested the bytes I reviewed.
- `v1.49.0` = `06e563e8b1fe…` locally **and** at `git ls-remote --tags origin`; lightweight tag, unmoved.
- `VERSION` still `1.49.0`; nothing matching `VERSION|CHANGELOG|version` appears in `git diff --name-only
  v1.49.0..HEAD`. **No version bump occurred**, as protocol requires.
- Closed artifacts `FINAL.md`, `IMPLEMENTATION.md`, `consensus.md`, `review/` are **untouched** across the
  entire `v1.49.0..HEAD` delta.
- Product delta vs the tag remains the reviewed set: `acp/spawn.go`, `evidence/execute.go`,
  `procctl/procctl_linux.go`, plus tests and the `.github` CI fix.

## 5. zcode-1's comment corrections — commit exists, behavioral equivalence **verified**

Commit `519951e` landed on the branch while I was waiting. It is **local only — not pushed**, so it was
not in run 36009912946. Touches `procctl_linux.go` (+25/−6), `procctl_cmdline_linux_test.go` (+12/−5),
`release-linux-u2-repair-zcode-1.md` (+73/−1).

I verified equivalence myself rather than accepting the claim:

- **Token-stream identity, comments excluded** (`go/scanner`, mode 0, over the `9134c7a` and `HEAD` blobs):
  `procctl_linux.go` **434 tokens both sides, streams byte-identical**;
  `procctl_cmdline_linux_test.go` **1975 tokens both sides, byte-identical**; 0 scan errors.
  (Independently reproduces zcode's 434/1975 counts.)
- **Every changed source line is a `//` comment or blank** — diff filtered for non-comment changes is empty.
- `GOOS=linux go build ./...` ok, `GOOS=linux go vet ./internal/procctl/...` ok, `gofmt -l` clean on both
  touched files. (`procctl_test.go` and `procctl_windows.go` show gofmt drift, but both are untouched
  across the whole `v1.49.0..HEAD` delta — pre-existing, not introduced here.)
- New hashes: `81c77dd4…053dfa`, `caa7dc79…bcded`.
- Record edit is effectively append-only: the single deleted line is the `inputs:` header, re-added with
  the two review files appended, plus an `amended:` line and a new `§ Post-review corrections`; the
  reviewed body is verbatim.

**Conclusion: the hosted green transfers to current HEAD `519951e`.** Compiled behavior is identical;
only comments and line numbers differ (cited spawn-race sites shift :170/:198 → :176/:204).

The corrections are substantive and, importantly, they *weaken* two claims in the direction my review
asked for — cancellation is **bounded, not immediate** for killed-but-unreaped pids, and the 100 ms bound
is **~1.8× the worst observed sample** (56.5 ms across 10,400 spawns), not "orders of magnitude". Both
now match what I measured. No safety facet, refusal string, or fail-closed path changed.

## 6. macOS and Windows

- **macOS: green.** Same command, 31 `ok` + 1 no-test-files, zero FAIL, go1.26.8 darwin/arm64.
  `trajectory 897.873s`, `app 662.019s` — the leg takes 16 min, which also explains why the 4-min Ubuntu
  leg is runner speed, not truncated work. `procctl 0.225s` confirms the Linux-only tests are build-tagged
  out and darwin is unaffected, as designed.
- **Windows: FAILED and is not claimed to pass.** Test step failed: `internal\evidence\tree_report_test.go:97:20:
  undefined: syscall.Mkfifo` (`[build failed]`), plus 14 failing packages / 16 ok, including
  `TestSpawnStop|WaitDrainsStderrBeforeReaping` (10.01s each, the D1 guard firing), `internal/agents`,
  `internal/app`, `internal/budget`, `internal/config`, `internal/driver`. I did not wait on, investigate,
  or attempt this track — it is owner-blocked and out of scope.

## 7. What this green run does and does not establish

It **is** x86_64 confirmation of a fix whose causal mechanism was independently reproduced twice (my
pre-fix reproduction and kimi-1's), with local bidirectional mutation proof and 2400-spawn post-fix
determinism. It is **not** proof on its own: the pre-fix race was probabilistic (~0.7% per spawn locally),
so a single green run could not by itself exclude regression — the weight is carried by mechanism plus
mutation evidence, and this run removes the "never observed on hosted x86_64" gap. Remaining limits, all
already recorded by zcode-1 and unchanged by this run: bound exhaustion is **silent**, the 100 ms margin
is measured-tail-based rather than kernel-instrumented, and one run is one sample.

## 8. Release 1.49.1 readiness — conditions

1. **Windows is red.** v1.49.0's CHANGELOG states "Windows behavior of `wait`/`usage` is exercised by the
   CI leg". The first actual hosted Windows execution contradicts that. Either fix the leg, or correct the
   claim in the 1.49.1 entry. Shipping Windows assets on this evidence would not be honest. Owner scope +
   npm approval still pending — unchanged by this report.
2. **Push or pin.** HEAD `519951e` is unpushed; remote `main` is `9134c7a`. Release either from `9134c7a`
   (tested) or push the comment-only commit first — it needs no new CI run to be trusted (§5), but the
   remote should be authoritative for whatever gets tagged.
3. **New immutable version only.** Product source changed vs `v1.49.0`, so this publishes as **v1.49.1**;
   `v1.49.0`/`06e563e` must never move. `VERSION`/`CHANGELOG` edits + tag are organizer/owner actions — I
   made none.
4. **Record the silent-exhaustion caveat** in the 1.49.1 notes, with zcode's operational rule: if U2
   recurs hosted, **raise `cmdlinePublishBound` before re-diagnosing**. Making exhaustion observable
   (a diagnostic on bound exhaustion) is a source change and would need its own version — not this one.
5. **Optional, not blocking:** one `-v` hosted ubuntu run would close the §3 skip-visibility limit by
   showing whether the two env-guarded procctl tests execute there.

Conditions 1–4 are organizer/owner decisions. Nothing in the Linux U1/U2 engineering track is blocking.

---

**NOT done (scope discipline):** no source change, no test/CI rerun, no re-run of the hosted workflow, no
release, tag, publish, install, or push; `v1.49.0` untouched; closed A–D artifacts, `FINAL.md`,
`IMPLEMENTATION.md`, `consensus.md`, `review/` untouched; no Windows investigation or claim; peer reports
and organizer records left as found. Local verification was read-only plus a `GOOS=linux` build/vet and a
throwaway token-comparison tool under `/tmp`.
