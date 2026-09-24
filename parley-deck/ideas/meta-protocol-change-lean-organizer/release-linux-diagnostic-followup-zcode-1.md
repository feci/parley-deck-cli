---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: narrow Linux release repair implementer (diagnostic preservation follow-up)
artifact: release-linux-diagnostic-followup
date: 2026-09-24
base: dd0a226
inputs: release-linux-hosted-diagnostics-claude-1.md (U2 minimal repair proposal, D1 residual)
run-cited: 35998165067 (hosted ubuntu job 107627994770) — U1 PASS, U2 recurred localized
---

# Linux diagnostic follow-up — zcode-1

Dispatched scope: implement claude-1's minimal diagnostic preservation —
persist the scrubbed, bounded, non-ExitError run error into
`CommandEvidence.Diagnostics` when the child output is empty, with truthful
envelope/hash semantics and meaningful tests for privacy and nonempty
behavior — plus the independently identified D1 bounded READY-handshake
guard. No U2 fix attempted (root cause still unknown — this change only
stops discarding the reason), no Windows contract changes, no broad fixes,
no retries/skips/pass-criteria weakening, U3 and Windows outside scope.

**Scope check: no expansion was needed.** The change lands as one 15-line
product hunk, one 19-line test-helper hunk, and one new 161-line test file.
That is the whole delta.

## U2 diagnostic preservation — applied as proposed, with one refinement

`internal/evidence/execute.go`, immediately after the `capbuf.overflow`
prefix block in `RunCriterionControlled`:

```go
if runErr != nil && len(out) == 0 {
    var exitErr *exec.ExitError
    if !errors.As(runErr, &exitErr) {
        ce.Diagnostics = ScrubAndTruncate("run error (no command output): " + runErr.Error())
    }
}
```

This is claude-1's snippet with one deliberate restriction the dispatch
asked for: the branch is **non-ExitError only**. An `ExitError` with empty
output keeps the previous envelope shape (exit code + empty diagnostics) —
the exit code already tells that story, signal deaths keep their existing
classification, and the diagnostics field becomes the discriminator between
"run-error branch" (labeled text present) and "child waited/exited" (empty).

Branches that now persist their reason (all previously discarded):

| runErr source | now persisted |
|---|---|
| `procctl.Attributed` refusal — claude-1 candidate (A) | `criterion supervisor identity unavailable: <facet>` |
| control refusal before/after start | the controller's error text |
| `cmd.Start()` failure (fork/chdir/exec) | the exec error text |
| `criterion control did not release its supervisor` | that text |
| envelope-artifact write failure — candidate (B) via control | the controller's error text |
| ctx-deadline non-ExitError wait errors | the ctx error text |

`OutputSHA256`, `CommandSHA256`, `ExitCode`, `DurationMillis`, `Format`,
`Status`, and `Complete` are all computed before and independently of this
hunk — nothing about the envelope's semantics changes except that
`Diagnostics` is no longer empty. A hosted recurrence of the exact
`dd0a226` signature (exit −1, empty output, ~1500 ms kill grace) will now
carry the discriminating string, separating claude-1's candidates (A)
attribution refusal — with the exact facet — from (B) artifact write.

**Truthful envelope/hash semantics** (pinned by test, not just claimed):
`OutputSHA256` still hashes only the stream the command itself emitted —
for this branch, the empty-string hash `e3b0c442…b855`; the label `run error
(no command output):` marks the persisted text as executor-side and states
that the command emitted nothing. When output is non-empty, diagnostics
remain exactly the scrubbed captured output (test-pinned).

**Privacy, precisely layered.** claude-1's report says the envelope "is
already secret-scrubbed before persistence"; the exact layering (verified in
source) is: `writeVerificationArtifact` (`internal/trajectory/verification.go:194`)
applies only canonical-JSON + a 1 MiB bound — **no** scrub there. The
secret-scrubbing of `Diagnostics` happens at construction in `execute.go`,
which is exactly where this change applies `ScrubAndTruncate` to the run
error before it ever reaches the envelope. Same treatment as captured
output: credential patterns (`ghp_…`, bearer/api_key/token/password shapes)
redacted to `«redacted»`, tail-bounded to 100 lines / 4096 bytes. Known
property, inherited from `ScrubAndTruncate` and accepted as proposed: the
bound keeps the **tail**, so an error longer than the bound evicts the label
prefix; all realistic branch reasons (attribution facets, exec errors) are
far shorter. The bounded test pins the length, not the label position.

## U2 tests — `internal/evidence/execute_runerror_test.go` (new, `//go:build !windows`)

Six tests, each pinning one property:

1. `TestRunErrorPreservedWhenOutputEmpty` — control refuses after `start()`
   (the hosted U2 reaping path: live supervisor, kill grace, empty output).
   Pins: label + reason in diagnostics, `OutputSHA256 == sha256Hex([]byte{})`,
   `Format == shell`, material command never ran, StatusFail/−1/incomplete.
2. `TestRunErrorPreservedOnStartFailure` — non-controlled path: absent
   working directory → `cmd.Start()` failure persists its reason.
3. `TestRunErrorPreservedIsScrubbed` — error text carrying `bearer …` and
   `api_key=ghp_…` tokens: raw secrets absent, `«redacted»` present, label
   present.
4. `TestRunErrorPreservedIsBounded` — 6000-byte error: diagnostics ≤
   `evidenceMaxBytes+3`, truncation marker present.
5. `TestRunErrorNotPersistedWhenOutputPresent` — control releases, waits for
   the material command to complete (marker-file handshake), then refuses.
   Pins: diagnostics are the captured output **only** (no run-error text),
   `OutputSHA256 == sha256Hex(output)`. The nonempty guard is what keeps
   the change from ever masquerading executor text as child output.
6. `TestRunErrorBranchSkipsExitError` — `exit 3`, no output: exit 3, empty
   diagnostics, complete observation — the ExitError exclusion.

A measurement note baked into test 5's comment: waiting on process liveness
is wrong here — the unreaped supervisor is a **zombie**, and a zombie
answers a signal-0 probe (`procctl.Alive`) as if alive; the marker file is
the only clean completion signal, and `cmd.Wait`'s drain-to-EOF guarantees
the already-piped bytes reach the capture buffer.

**Bidirectional proof (claude-1's method):** deleting only the 15-line
product hunk → tests 1–4 FAIL (`--- FAIL:` ×4, 0.32s); restoring the file
byte-identical (`cmp`-verified) → all six PASS (0.25s). Tests 5–6 pass in
both worlds by design — they pin the pre-existing shape the change must not
touch.

## D1 — bounded READY handshake in `spawnGatedChild`

`internal/acp/spawn_test.go`: the unguarded
`bufio.NewReader(p.Stdout()).ReadString('\n')` is replaced by a reader
goroutine + `select` with the file's established 10 s `time.After` guard
(same idiom as F-9). A platform whose pipe buffer cannot hold the in-flight
16 KiB now fails in 10 s at a named site
(`child READY handshake did not complete within 10s …`) instead of hanging
to the 45 m binary-wide timeout. The parked reader goroutine outlives the
helper only on the already-failing path.

Honest limitation, recorded rather than hidden: the guarded stall **cannot
be induced on this host** — Linux/macOS 64 KiB pipe capacity means the
handshake never blocks here, so no local mutation discriminates the guard
(as F-1's mutation discriminated its line). The bound is structural. This
is deliberately NOT a Windows contract change: no build tag added or
removed, `spawn_test.go` stays Windows-eligible exactly as claude-1 found
it, and the owner-blocked Windows track is untouched and unwarranted. D1's
pipe-capacity contingency (kimi-1 K1) remains with that track.

## NOT established, NOT done

- **U2 root cause remains NOT established.** Nothing was fixed, skipped,
  retried, or weakened; candidates (A)/(B) are still claude-1's deduction,
  now observable instead of invisible. A green hosted run proves nothing
  about U2 (claude-1's own framing) — the purpose of this change is **one
  evidence-producing hosted run**, not a U2 solution.
- No U3, no W1–W10, no `.github`/workflow change, no version bump
  (VERSION/version.go/CHANGELOG untouched). **Product source changed here**:
  per protocol this delta publishes only as a new immutable version
  (v1.49.1), organizer/owner decision — `v1.49.0` (`06e563e`) never moves.
- No push, tag, release, install. Closed FINAL/IMPLEMENTATION/review/
  consensus records untouched; claude-1's report left untracked for the
  organizer; `organizer-usage.md` not touched.
- No full-suite rerun (last full evidence: release-ci-fix-zcode-1.md);
  validation below is the touched-code focus.

## Validation (go1.27.1 darwin/arm64)

- `go build ./...` ok; `go vet ./internal/evidence/ ./internal/acp/` clean;
  `gofmt -l` clean on all three touched files.
- Mutation proof above (4 FAIL → restore → PASS).
- `go test ./internal/evidence/` full package **ok 5.171s** (includes the
  six new tests).
- `go test ./internal/acp/ -count=3` **ok 8.178s**; `-race -count=1`
  **ok 3.870s** (exercises the guarded handshake).
- Named dependents: trajectory trio
  (`TestParentRecoveryPublicationInterruptionAndExactReplay`,
  `TestRecoveredParentObservationRequiresExactLineage`,
  `TestVerifierBudgetRefusalRecoversExplicitly`) **ok 74.197s**; runner
  `TestVerifierBudgetRefusalRecoversExplicitly` **ok 6.924s**; app
  `TestTrajectoryParentRecoveryRestoresPreviouslyBoundFacts`
  **ok 8.501s** — their retained envelopes and pass criteria are unchanged
  (only refused steps' diagnostics gain text, which nothing asserts empty;
  grep-verified across trajectory/runner/app/protocol tests).

## Provenance

| File | Delta | sha256 |
|---|---|---|
| `internal/evidence/execute.go` | +15 | `094d4f17…9e6417` |
| `internal/acp/spawn_test.go` | +19/−2 | `374bcdf5…8c9501` |
| `internal/evidence/execute_runerror_test.go` | new, 161 lines | `819ecbd5…a265d` |

Throwaway artifacts reverted and verified: the mutation run used a stashed
copy (`/tmp/execute.go.fixed`, `cmp`-identical after restore); the manual
supervisor timing probe lived in `/tmp/sup.sh` and never touched the repo.

## Hosted next step (organizer)

Independent delta review of this commit, then push for **one untagged
hosted ubuntu run** on the touched packages
(`go test ./internal/evidence/... ./internal/acp/... ./internal/trajectory/... ./internal/runner/... ./internal/app/... -count=1`,
the CI matrix equivalent). A U2 recurrence now self-identifies: the
`step-00N.json` envelope at the failing ordinal will carry either
`criterion supervisor identity unavailable: <facet>` (claude-1's A, facet
named) or the artifact-write/controller text (B) in
`command.diagnostics`, alongside the existing dump of the retained set. If
the run is green, the channel status is exactly what it was before this
change plus one closed residual — not a U2 claim.
