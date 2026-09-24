---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent non-implementer reviewer (Linux repair delta)
artifact: release-linux-delta-review
date: 2026-09-24
reviewed: 56da6faf82110494e6b977ed443b348ae1a8e0f4
base: a2db7998fae61d44fdb9786d3c12b4a5b9cc868f
inputs: release-linux-repair-followup-zcode-1.md, release-linux-review-claude-1.md (F-1/F-2/F-9), release-linux-review-kimi-1.md (K1/K2)
verdict: PASS — ready to push the untagged hosted diagnostic run
---

# Linux repair delta review — claude-1

## Verdict

**PASS.** `a2db799..56da6fa` closes all three of my actionable findings and records K1
honestly. Every claim in the follow-up record that I could test, I reproduced
independently; none failed. The delta is test-only, the product source is untouched, and
the commit is ready for the organizer to push as an **untagged hosted diagnostic run**.

One new **LOW** finding (D-1), non-blocking for the ubuntu push. Nothing here makes the
channel green: **U2's root cause remains unestablished**, and a clean hosted run would
*not* root-cause the earlier failures — it would only fail to reproduce them.

**Independence.** Reviewed in my own detached worktree
(`worktrees/lean-organizer-review-claude-1` at `56da6fa`), go1.27.1 darwin/arm64,
macOS 27.0. Every mutation was reverted and `git status --porcelain` verified empty, with
`git diff 56da6fa` empty, before each confirmation run and at the end. I made no edit to
source, to either prior review, or to any closed artifact; no commit, push, install, tag,
release, version bump or CI run.

## Delta findings

| ID | Severity | Item | Result |
|---|---|---|---|
| D-A | PASS | F-2 exit-19 diagnostic rewiring | both branches dump the full 9-envelope set; classifying datum present |
| D-B | PASS | F-1 escaped-writer discriminator | bidirectional mutation reproduced exactly |
| D-C | PASS | F-9 bounded waits | guards fire on a forced real hang, not just present in source |
| D-D | PASS | scope / product source | test-only; `spawn.go` byte-identical; tag unmoved |
| D-1 | LOW (new) | K1 disposition precision | "fails loudly" is really "blocks until the global timeout"; gated tests are Windows-**eligible** |
| D-2 | INFO | U2 dump reading | envelope richness depends on how far the helper got — absence is a weaker signal |
| D-3 | INFO | hosted line mapping | hosted `:124` is now `:126`; match by test name, not line |

## D-A (PASS) — F-2: the exit-19 site now dumps real envelopes

The remedy is exactly the one my review derived: `parentRecoveryFixture` binds the journal
path `journalFixture` returns (`parent_recovery_test.go:75`) and both failure branches call
`dumpRetainedVerificationSteps`. `dumpAllRetainedVerifications` and its wrong-directory
glob are gone, with no dangling reference (grep across the repo: six call sites, all
`dumpRetainedVerificationSteps`; `filepath` still used, so no orphaned import).

I did not take the record's word for it. I forced **each branch independently** in my own
tree and read what was actually printed:

| forced branch | subtests | diagnostics lines | distinct envelopes | unavailable/unreadable |
|---|---|---|---|---|
| `cmd.Wait` — the hosted `helper: exit status 19` site | 4/4 FAIL | 36 | `receipt.json`, `step-001..004`, `process-001..004` (9 × 4) | **0** |
| `PreviewParentRecovery` sibling | 4/4 FAIL | 36 | same 9 × 4 | **0** |

**The classifying datum is genuinely there**, not just the filenames. Field counts across
one forced run: `exit_code` 16, `diagnostics` 16, `complete` 16, `duration_ms` 16;
`pid`/`pgid`/`boot_id`/`proc_start`/`started_at` 16 each; `failure_stage` 4. That is what
separates signal death (`exit_code` ≥ 128) from a start failure (no process identity) from
ctx/overflow — the discrimination the repair plan asked for.

**Privacy and bound re-verified at the newly wired site** (not inherited from my earlier
measurements): 49,025 B across the 4 failing invocations ≈ **12.3 KB each**, consistent
with the 12.5–13.3 KB I measured at the sibling sites; credential-pattern scan (bearer /
`sk-` / `gh?_` / `AKIA` / PEM / labelled secret-key-token) → **0 hits**.

**Inert on green** and pass criteria untouched: the trio
(`TestParentRecoveryPublicationInterruptionAndExactReplay` 23.29s,
`TestRecoveredParentObservationRequiresExactLineage` 46.56s,
`TestRecoveredParentPublicationRefusesNonRefusalOriginal` 5.75s) → `ok 75.896s`,
**0** `diagnostics:` lines. The diff shows only the dump call swapped; every assertion is
byte-identical.

**One thing the record under-claims.** `parentRecoveryFixture` has a second caller —
`activation_quorum_test.go:244`, in
`TestActivationQuorumRejectsWidenedParentRequestBeforePreviewHash` — which silently
inherits the same correction. Strictly a gain, no added risk, but the follow-up describes
the fix as covering only the four `…InterruptionAndExactReplay` subtests.

## D-B (PASS) — F-1: the discriminator genuinely discriminates

Mechanism confirmed in source before testing: `setSysProcAttr` gives the spawned child
`Setsid` (pgid == pid), and the new grandchild takes its **own** `Setsid`, so
`killProcessGroup(child)` structurally cannot reach the only surviving writer.

Bidirectional mutation, my tree, deleting *only* `spawn.go:152` (`_ = p.stderr.Close()`):

```
line deleted:   --- PASS: TestSpawnStopTimeoutIsBoundedAfterKill  (0.30s)
                --- FAIL: TestSpawnStopIsBoundedWhenStderrWriterEscapesGroup (10.01s)
                    spawn_unix_test.go:62: Stop did not return after interrupting
                                           the group-escaping stderr writer
restored:       ok  internal/acp  (5/5 repeats, 1.632s total; single run 0.31s)
```

This reproduces the implementer's numbers exactly and, in the same run, re-confirms my
original F-1 finding — the pre-existing bounded-after-kill test is *still* green with the
line gone. The FAIL is itself the proof that the writer escaped the group: had it not,
the copier would have hit EOF at the kill and `Stop` would have returned without the
`Close`. F-1 is closed.

The Unix build tag is structural, not a skip: the escape needs POSIX process groups and
the interrupt needs a pollable pipe fd. The in-file doc says so and claims nothing about
Windows. Correct, and it is the right way to answer F-8 without editing reviewed product
source.

## D-C (PASS) — F-9: guards verified against a real hang, not just read

Presence in source is weak evidence, so I forced an actual permanent drain hang
(suppressing `close(observer.release)`, leaving the copier parked forever):

```
--- FAIL: TestSpawnStopDrainsStderrBeforeReaping (11.01s)
    spawn_test.go:139: Stop did not return after the copier was released
--- FAIL: TestSpawnWaitDrainsStderrBeforeReaping (11.01s)
    spawn_test.go:159: Wait did not return after the copier was released
```

Both fail fast with their named messages instead of consuming the 45m budget and panicking
the binary. That is precisely what F-9 asked for. Reverted, both PASS (1.01s each).

## D-D (PASS) — scope and product source

- Delta touches **4 `_test.go` files + 1 record**. Nothing else.
- `git diff a2db799..56da6fa -- internal/acp/spawn.go` is **empty** — `spawn.go` is
  byte-identical, as claimed.
- Product delta vs the tag (`git diff --name-only v1.49.0 56da6fa -- . ':!parley-deck'
  ':!.github'`) still lists exactly **one non-test file: `internal/acp/spawn.go`** — the
  one U1 file both reviews already cleared.
- `.github/`, `VERSION`, `CHANGELOG.md`, `internal/app/version.go`: **untouched**.
- Closed artifacts (`FINAL.md`, `IMPLEMENTATION.md`, `consensus.md`, `review/`) and **both
  prior reviews**: untouched. `organizer-usage.md` correctly left out of the commit.
- `v1.49.0` = `06e563e8b1fe8094132149b83e14da7be9aae51e`, **unmoved**.
- **Zero** `t.Skip` / `testing.Short` / `continue-on-error` / retry / suite `-run` filter
  added. The `-test.run=^…$` strings in the new file are child re-execs, the existing
  `TestGatedStderrChild` pattern, anchored.
- Health on my tree: `go build ./...` ok, `go vet ./...` clean, `gofmt -l` clean on both
  touched packages, `internal/acp -count=3` **ok 8.225s**, `-race` **ok 3.874s** (record:
  8.179s / 3.861s). No stray `escaped-stderr` / `gated-stderr` processes survived any run.

## K1 — disposition is honest; one imprecision (D-1, LOW, new)

**The substance is right and I endorse the decision.** Refusing a speculative chunk-shrink
is correct: there is no principled size without a Windows host to measure on, and the
16 KiB volume is load-bearing — it mirrors the hosted `0 of 16384` truncation signature the
tests exist to reproduce. Shrinking it to satisfy an unverifiable hypothetical would trade
a real mechanism for an imagined one. Handing the contingency to the owner-scoped Windows
track matches kimi-1's own framing. No Windows pass is claimed anywhere, and none exists.

Two points where the wording outruns what I can verify:

1. **"fail the handshake below loudly" overstates the mode.** On the hypothesised
   small-buffer platform the outcome is a three-way block — child blocked writing its
   second chunk, copier parked in the observer, test blocked in the *unguarded*
   `bufio…ReadString('\n')` in `spawnGatedChild` — with nothing to release it. It surfaces
   only at the global test timeout, as a binary-wide panic. That is noisy, and kimi-1's
   "never silently mis-pass" holds, but it is the **exact failure mode F-9 removed two
   functions below in the same file**. The same 10 s guard on that handshake would make the
   comment's claim true, costs nothing, and does not touch the 16 KiB volume the
   disposition rightly wants to keep. Not required before the ubuntu push.
2. **"these tests have never run on Windows" is true, but does not mean they cannot.**
   `spawn_test.go` carries **no build tag** and compiles for Windows — I confirmed with
   `GOOS=windows go vet ./internal/acp/` (clean). So the gated tests are Windows-*eligible*
   the moment a Windows leg executes, unlike the new `spawn_unix_test.go`, which is
   structurally excluded. K1 is dormant, not neutralised; it becomes live on the first
   Windows run.

K2 / F-7 / F-8 dispositions (no action, informational) are correct as recorded.

## Residuals for reading the hosted run

- **D-2.** A forced-failure proof necessarily forces the branch *after* a successful helper
  run, so 9 envelopes is the **best** case, not the guaranteed one. A genuine failure
  before any retention will print `verification journal … unavailable`. That is still the
  "no process identity = start failure" classification, but it is a weaker signal than nine
  envelopes — read its absence as data, not as the diagnostic having failed again.
- **D-3.** The hosted `:124` site maps to the fixture call
  `root, p := parentRecoveryFixture(t)`, which is line **126** at `56da6fa` (both the
  fixture and the `Fatalf` are inside `t.Helper()`, so the subtest call site is what gets
  reported). Match a recurrence by **test name**, not by line number.

## Readiness

**Ready to push the untagged hosted diagnostic run.** No precondition remains from my prior
review: F-1, F-2 and F-9 are closed and independently verified. D-1 is a LOW durability
nit on a Windows-only hypothetical and should not gate the ubuntu leg.

Expect `internal/acp ok` (U1 proven in both directions, twice, on two machines). For U2,
read `execution.record.command.exit_code` + `diagnostics` in `step-*.json` and the identity
fields in `process-*.json`, and classify. **A green run classifies nothing** — it is not
evidence that the `:124`/`:142`/`:229`/`:153`/`:509` family is fixed, because nothing about
that family was fixed; only its observability was.

**Not authorised by this review and not performed by me:** version bump, tag, publish,
merge, workflow change, release-notes edit, install or CI trigger. `v1.49.0` stays
immutable at `06e563e`; no new version exists. U2's root cause is **unknown** and the
release channel is **not green**.

## Limitations

- macOS arm64 / go1.27.1 only. No Linux or Windows host. The Windows reasoning in D-1 is a
  compile-eligibility fact (`GOOS=windows go vet`) plus the pollable-fd contract — not an
  executed Windows result, and I make no Windows behavioural claim.
- I forced failures at both corrected branches and read the output; I did not reproduce the
  hosted `exit status 19` condition itself, which no one has reproduced locally.
- Delta review only (`a2db799..56da6fa`). I did not re-review the U1 fix or the four sibling
  dump sites beyond confirming they are unmodified — they are cleared in my prior review —
  and I ran no full 31-package suite: nothing in this delta required it.
- The `runner` and `app` dump sites are untouched by this delta; their F-3 verification
  stands from the prior review.
- No credentials, tokens or log secrets appear in this report.
