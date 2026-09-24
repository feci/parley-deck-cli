---
idea: meta-protocol-change-lean-organizer
author: claude-1
role: independent reviewer (non-implementer) of the narrow diagnostic delta
artifact: release-linux-diagnostic-review
date: 2026-09-24
delta: dd0a226..bc5ff29 (4 files: execute.go, execute_runerror_test.go, spawn_test.go, zcode-1 record)
inputs: release-linux-hosted-diagnostics-claude-1.md (my proposal), release-linux-diagnostic-followup-zcode-1.md
toolchain: go1.27.1 darwin/arm64
verdict: PASS — safe to push as ONE untagged hosted diagnostic run; NOT a release
---

# Narrow diagnostic delta review — claude-1

## Verdict

| Item | Status |
|---|---|
| U2 diagnostic preservation (`execute.go:219-224`) | **PASS** — correct guard, truthful envelope, no lifecycle change |
| No false successful/complete envelope | **PASS — structurally impossible**, not merely untested |
| `OutputSHA256` still the emitted-output hash | **PASS** |
| Privacy tests meaningful | **PASS**, with a scoping note (INFO-1) |
| D1 READY timeout guard (`spawn_test.go:101-122`) | **PASS** |
| No leaked goroutine / child on the green path | **PASS** (verified, incl. `-race` and post-run `pgrep`) |
| Abandoned child on the failure path | **LOW, pre-existing** (D1-3) — not a push blocker |
| Delta scope discipline | **PASS** — no version bump, no workflow, no closed artifact touched |
| U2 root cause | **still UNKNOWN**, correctly claimed as such by the implementer |

Readiness is for **one untagged hosted ubuntu run whose purpose is evidence**. A green
run would make no U2 claim. This is not readiness for v1.49.1.

## U2 preservation — findings

**V-1 · No false complete/successful envelope is structurally impossible.**
The guard at `execute.go:219-221` (`runErr != nil && len(out) == 0 && !errors.As(runErr,
&exitErr)`) is the exact complement of the only condition under which `complete` can stay
true with a run error: `execute.go:288`, `complete = complete && errors.As(runErr,
&exitErr) && exitErr.ExitCode() >= 0`. Every envelope that can carry preserved text is
therefore `Complete=false`, and `execute.go:261-262` (`case runErr != nil: status =
StatusFail`) forces `StatusFail`. This holds by construction, independent of the tests.

**V-2 · `OutputSHA256` unchanged.** It is computed at `execute.go:199` from
`hasher.Sum(nil)`; the hasher is fed only by the `io.MultiWriter` sink at
`execute.go:145-148` wired to `cmd.Stdout`/`cmd.Stderr`. The hunk mutates `ce.Diagnostics`
alone, after the struct literal closes at `:206`. The executor text can never reach the
hash. Pinned twice: `execute_runerror_test.go:41` (hash of the empty string) and `:138`
(hash of exactly the emitted bytes).

**V-3 · The assignment cannot clobber the overflow warning.** The hunk *assigns* where
`:208` *prepends*, so the two would conflict if both could fire. They are mutually
exclusive: `captureMax = 4<<20` (`execute.go:44`) and `cappedWriter.Write` only sets
`overflow` when a chunk exceeds pre-write room, which requires the buffer already filled to
4 MiB — so `overflow == true` implies `len(out) == 4 MiB ≠ 0`. Verified by reading `Write`,
not assumed.

**V-4 · No later writer overwrites the preserved text.** The three downstream diagnostics
writers (`:234`, `:248`, `:255`) all prepend, and all are gated behind parsing a non-empty
`out` — unreachable on this branch. Empirically confirmed: `Format` stays `shell`
(`execute_runerror_test.go:44`).

**V-5 · Pass criteria unchanged.** The only non-test consumers of the field are
`app/evidence_table.go:188` (rendering; re-scrubs through the same function) and
`trajectory/trajectory.go:156`, where `classify` rejects `len(cmd.Diagnostics) > 8192`. The
preserved text is bounded to `evidenceMaxBytes+3 = 4099`, i.e. 2× headroom — and `classify`
already rejects on `!e.Complete` first (V-1). Nothing reparses `Diagnostics` for semantics.

**V-6 · Lifecycle unchanged.** The hunk is pure field mutation, placed after the whole
Start/Wait/KillGroup sequence (`execute.go:180-196`) and before any parsing. No new
syscalls, no new defers, no control-flow change; `errors` and `exec` were already imported.

**V-7 · The diagnostic will actually discriminate — and better than my proposal stated.**
I re-read `procctl.Attributed`: it returns **14** distinct constant reason literals (my
earlier report said 8 facets; the real count is 14). All are fixed English text carrying no
paths, pids or secrets. Propagation is verbatim: `verification.go:613-616` does `sp, err :=
start(); if err != nil { return err }`, which returns through `withVerification` as the
control function's error and lands in `runErr`. So candidate **(A)** will print
`run error (no command output): criterion supervisor identity unavailable: <one of 14>`,
naming the facet outright. Candidate **(B)** (`writeVerificationArtifact`,
`verification.go:617-620`) yields a path-bearing os error — also discriminating. The
branch-vs-branch question I could not answer from run 35998165067 becomes answerable from
the next one.

**V-8 · Negative tests carry weight.** `TestRunErrorNotPersistedWhenOutputPresent`
(`:101-141`) asserts *exact equality* against `ScrubAndTruncate(out)` plus the exact output
hash; `TestRunErrorBranchSkipsExitError` (`:146-160`) asserts `Complete=true`, `ExitCode=3`
and empty diagnostics. zcode-1 notes these pass with the hunk deleted, which is true and is
the right disclosure. Adding to it: they *do* discriminate the plausible mis-implementations
of this hunk — prepend-always, or dropping the `errors.As` restriction — so they are not
decorative.

**V-9 · Bidirectional proof, by a non-mutating route.** The dispatch forbids source edits,
so I did not re-run zcode-1's deletion experiment. Static equivalent:
`grep -rn "run error (no command output)" --include='*.go' .` returns exactly one producer
(`execute.go:222`) and three test consumers (`:39`, `:60`, `:81`). The four preservation
assertions have no other source for that text and must fail if the hunk is removed. Same
conclusion, different route; I am not restating zcode-1's experiment as reproduced.

### INFO-1 · Privacy tests: routing is proven, realism is deliberately not

`TestRunErrorPreservedIsScrubbed` (`:65-84`) drives two distinct pattern families
(`secretPatterns[1]` bare `bearer …`, `secretPatterns[2]` `api_key=…`) and asserts both
literals absent, `«redacted»` present, and the label retained. The scoping is correct: the
scrubber has its own coverage (`app/driver_checks_test.go:16 TestScrubAndTruncate`), so this
test's job is to prove the new text *routes through* it — which it does.

The dispatch asked whether the tests use plausible error strings. They do not, and that is
defensible rather than a gap, because I checked the realistic side directly (V-7): the
strings this branch will really carry are 14 fixed attribution literals and os path errors,
none credential-shaped. Paths in candidate-(B) and start-failure text are not secrets under
this scrub and are already persisted in `ce.Command` and `process-NNN.json`. The exposure
class is unchanged by this delta. Recorded as INFO, no action requested.

## D1 READY guard — findings

**D1-1 · PASS.** `spawn_test.go:101-122` matches the file's established idiom (10s
`time.After`, as at `:155`, `:176`, `:198`); `time` was already imported. `t.Fatal` stays on
the test goroutine in both select branches — required, since calling it from the reader
goroutine would be a misuse. **No build tag added or removed**: both `dd0a226` and `bc5ff29`
open `spawn_test.go` with `package acp`, so Windows eligibility is byte-identical to what I
found in my prior review. Windows scope remains with the owner; nothing here touches it.

**D1-2 · No leak on the green path — verified, not argued.** `ready` is buffered (cap 1),
so the reader goroutine never blocks on send; on success it has already returned from
`ReadString` before `spawnGatedChild` returns. `p.stdout` is assigned once in `Spawn` and
never reassigned, and the goroutine is its only reader, finishing before any caller reaches
`Stop`/`Wait` — no race window. Confirmed: `go test ./internal/acp/ -count=1 -race` → ok
3.873s, and `pgrep -fl 'gated-stderr-child|stderr-observer-child|parking-stderr-child'`
after the run returned nothing.

**D1-3 · LOW, pre-existing, now reachable by one more route — not a push blocker.**
`spawnGatedChild` registers no `t.Cleanup` for the child it spawns. Both `t.Fatal` sites in
the helper (the pre-existing handshake-failure branch at `:118` and the new timeout branch
at `:121`) abort before `return p`, so the caller never receives the handle and `p.Stop()`
is never called: the child is abandoned for the rest of the test binary. The new branch adds
a second route into that gap; it does not create it. Impact is bounded — `Spawn` gives the
child pipes for all three std fds (`StdinPipe`/`StdoutPipe`/`StderrPipe`), so an abandoned
child cannot hold a CI step's stdout open and hang the job; the classic runner-hang mode is
unavailable here. The in-file remedy already exists one file over
(`spawn_unix_test.go:48`, `t.Cleanup(func(){ _ = syscall.Kill(grandPID, syscall.SIGKILL) })`).
Suggested, not required: `t.Cleanup(func(){ _ = p.Kill() })` right after `Spawn`. I
deliberately do **not** ask for it now — it is orthogonal to the diagnostic, and this
delta's value is that it is small.

**D1-4 · INFO.** The guard is structurally unfalsifiable on Linux/macOS (64 KiB pipe
capacity means the child never stalls, so no local mutation discriminates it). zcode-1
records this honestly in-file (`:101-105`) and in the record; no verification claim is
overstated. I concur — it buys a named 10s failure site instead of a 45m binary-wide
timeout, and nothing more.

### DOC-1 · One inaccurate ordering phrase in the implementer's record

The commit message and record say `Status` and `Complete` are "computed before the hunk".
Precisely: `OutputSHA256`/`CommandSHA256`/`ExitCode`/`DurationMillis` and `Format`'s initial
value are computed before (`:197-206`); `status` (`:260`) and `complete` (`:279-289`) are
computed **after** it, and `Format` can still be reassigned at `:231`/`:243`. The conclusion
is nonetheless correct — none of those later computations read `ce.Diagnostics`; they read
`out`, `runErr`, `capbuf.overflow` and `ctx.Err()` — so "independently of this hunk" is the
accurate half of the claim. No code change needed; flagged so a later review does not lean
on the ordering wording.

## Checks I ran (read-only; no edits, no commits, no push, no installs)

| Check | Result |
|---|---|
| `gofmt -l` on all three changed files | clean |
| `go vet ./internal/evidence/... ./internal/acp/...` | clean |
| `go test ./internal/evidence/ -run TestRunError -count=1 -v` | **6/6 PASS** |
| `go test ./internal/evidence/ -count=1` | ok 5.022s |
| `go test ./internal/acp/ -count=1 -race` | ok 3.873s |
| `go test ./internal/trajectory/ -count=1 -run '<3 U2-site tests>'` | ok 71.704s |
| `git diff --name-only dd0a226..bc5ff29` | exactly 4 files |
| same, filtered to `VERSION CHANGELOG.md internal/version .github` | **empty** |

No full-suite rerun (unnecessary for a delta this narrow; last full evidence stands in
`release-ci-fix-zcode-1.md`). No hosted rerun, tag, release, install or version bump. No
source edit — including no deletion mutation, hence V-9's static route. No closed
FINAL/IMPLEMENTATION/review/consensus artifact read or written.

## Safe-to-push judgment

**Safe to push for one untagged hosted ubuntu diagnostic run.** The product delta is 15
lines confined to a failure path that already produced `StatusFail` + `Complete=false`; it
adds information to retained evidence and removes none; it cannot manufacture a passing or
complete envelope (V-1); it cannot alter a hash (V-2); it cannot change classification
(V-5). The test delta is test-only and Windows-neutral. Worst case on the hosted run is the
status quo ante plus a labelled reason string.

**Explicitly not established by this review:** U2's root cause, which remains **UNKNOWN**.
A green hosted run proves nothing about U2 and must not be read as closing it. The channel
is not green. The Windows leg stays owner-blocked and untouched. Because product source
changed, publication would require a new immutable version — organizer/owner decision, not
mine, and outside what this push is for; `v1.49.0` (`06e563e`) stays unmoved.

## Limitations

- Everything above is source reading plus local darwin/arm64 execution. U2 has never
  reproduced on any local machine, and this review did not attempt it.
- V-7's propagation claim is read from source; `withState`'s wrapping was not exhaustively
  traced, so the reason string may arrive wrapped by outer context — wrapped text still
  names the facet, which is what the diagnostic needs.
- D1 remains unexecuted on the only platform where its stall hypothesis is live.
- No credentials or secrets appear in this report.
