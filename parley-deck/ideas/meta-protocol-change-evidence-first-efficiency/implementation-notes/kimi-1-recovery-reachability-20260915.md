---
agent: kimi-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-15
reviewed-commit: 39107b138e77de064f21fec4aafa6073887c2783
review-kind: executed-reachability-probe (native fixture; NOT production recovery implementation)
status: final
---

# Recovery reachability probe: verifier-launch budget refusal and the no-handle inference (kimi-1)

## Scope, provenance, and limits

- Frozen HEAD `39107b138e77de064f21fec4aafa6073887c2783` (verified via `git rev-parse HEAD`).
- Continues `.parley-runtime/prior-proposal.md` with EXECUTED evidence. Owned files:
  `internal/runner/verification_refused_reachability_test.go` (new, only file added
  under `internal/`) and this note. Scratch (runner script, Go cache, logs) lives
  only under `.parley-runtime/work/`. `git status --porcelain` shows exactly these
  two owned paths as new; no other source/test/protocol/config file was touched.
- **This proves CURRENT behavior of the frozen tree. It is not an implemented fix,
  not whole-goal acceptance, and no finding is withdrawn.** The proposed source
  changes (S1/S2 in the prior proposal) require independent review before
  integration. No production overlays were used; no quarantine or retry was
  implemented; no other agents, model calls, network, browser, git mutations,
  package installs, or full-suite reruns occurred.

## PRIMARY executed evidence

Commands run through `python3` `subprocess.run(..., stdout=PIPE, stderr=PIPE)`;
each log was written only after the command terminated. Environment:
`GOCACHE=/…/checkout/.parley-runtime/work/gocache`, `TMPDIR`/`GOTMPDIR=/…/checkout/.parley-runtime/work/tmp`
(native filesystem, module cache left at the default to avoid downloads).
Toolchain: `go version go1.27.1 darwin/arm64`.

1. Compile check (runs no tests):
   `go test -count=1 -run ^$ ./internal/runner/` → EXIT 0
   (`.parley-runtime/work/logs/build.log`; one build error during authoring —
   struct-with-slice comparison — was fixed in my own new test only; retained in
   the first build.log run history of the work dir. No pre-existing failure was
   touched.)
2. Focused run (only the two new tests; accepted tests were not re-run):
   `go test -count=1 -v -run TestVerifierLaunchBudgetRefusalBricksTicketWithLedgersRetained|TestMissingRequestJSONDoesNotProveNoTicketHandle ./internal/runner/`
   → EXIT 0, elapsed 10.55s. Decisive output, quoted verbatim from
   `.parley-runtime/work/logs/focused-tests.log`:

       --- PASS: TestVerifierLaunchBudgetRefusalBricksTicketWithLedgersRetained (3.00s)
       --- PASS: TestMissingRequestJSONDoesNotProveNoTicketHandle (6.64s)
           --- PASS: TestMissingRequestJSONDoesNotProveNoTicketHandle/removed (2.17s)
           --- PASS: TestMissingRequestJSONDoesNotProveNoTicketHandle/corrupted (2.48s)
           --- PASS: TestMissingRequestJSONDoesNotProveNoTicketHandle/concurrent-prepare-executed (1.99s)
       ok  	parley-deck-cli/internal/runner	9.883s

Exact test source SHA256 (the file that produced this output):
`b614a41d2c97ad2db2fd2abfbc3f036cc6ef831c293296933165ce1ce290321f  internal/runner/verification_refused_reachability_test.go`
Runner script SHA256:
`6ad7c6e4db59c2bd8678fd11bfcd1fd0e3aabb39a02973076ae14050ad9477f8  .parley-runtime/work/run_tests.py`

## Task 1 — MAJOR-1 PROVEN as current behavior (brick confirmed end-to-end)

Test: `TestVerifierLaunchBudgetRefusalBricksTicketWithLedgersRetained`
(internal/runner/verification_refused_reachability_test.go). Real public
production paths only: `trajectoryRuntimeFixture` (existing helper), a real
changed-source attempt via public `RunMeasured` (phase `fixup`, local `/bin/sh`
fixture `telemetryShell`, never a real LLM), a real captured ticket via public
`trajectory.PrepareCapturedVerification`, then the verifier launch via public
`RunMeasured` with `WithCapturedVerification(ticket)` +
`WithLaunchBudget(Denied[budget.Launch])`. Executed findings:

- **Ordering witness.** The denial is a genuine pre-start budget refusal at the
  verifier launch: returned record has `StartedAt == nil`, `PID == nil`,
  `ExitCode == nil`, `FailureClass == "budget_refused"`, and the error contains
  `"launch is forbidden"` — pinning the refusal to the explicit launch policy, not
  the step session or the cycle gate. The child never spawned
  (`.parley-runtime/verifier-spawned` absent). Consistent with the source ordering
  `telemetry.go:96-100` (`reserveCapturedVerification` before `reserveBudget`).
- **Retained telemetry, no fabricated lifecycle.** The refused invocation's dir
  holds `requested.json` (original launch retained) and `terminal.json`
  (`status "failed"`, `failure_class "budget_refused"`, phase
  `trajectory-verification`, `StartedAt`/`PID` nil); `started.json` is absent.
  The central witness's lifecycle was never fabricated or hand-created.
- **Ticket pinned to the dead invocation.** Journal `launch.json` binds
  `invocation_id == <refused invocation>` and `ticket_sha256 == ticket.SHA256()`;
  `request.json` is byte-identical before/after.
- **Which budget was refused — distinguished.** The explicit launch store has NO
  `ledger.json` at all (a denied reservation persists nothing: the change error
  precedes any write, `internal/budget/ledger.go:342-356`), while the fixup cycle
  `ledger.json` is byte-identical pre/post with `Count == 1` (the original patch
  charge). A later unverified fixup retry is ALSO `budget_refused`, but at the
  cycle path — distinguished by `Metadata.Phase == "fixup"` (vs
  `trajectory-verification`) and by both stores remaining untouched.
- **Every current public recovery/charge path refuses explicitly.** Prepare retry
  → `"verification already reserved or unavailable"`; `ExecuteCapturedVerification`
  with a fresh invocation → `"verification invocation differs from its durable
  reservation"` with no `claim.json` written; the production verifier entrypoint
  shape — public `RunConsult` with the bound ticket (the path
  `internal/app/trajectory_verify.go:244` uses) → refused before spawn
  (`ExitError != ""`, empty answer, no spawn marker); `RequireResolved` fails.
- **Exact retention at the end.** Journal inventory is exactly
  `{request.json, launch.json}`, both byte-identical to their post-refusal
  capture. No recovery artifact appeared anywhere.

Conclusion: MAJOR-1 stands as executed fact — under current public APIs the
verifier-launch budget refusal permanently bricks the trajectory. The witness was
reachable under actual public ordering; nothing impossible was forced.

## Task 2 — the no-handle inference is REFUTED as stated; a weaker property survives

Test: `TestMissingRequestJSONDoesNotProveNoTicketHandle`.

- **Refutation (executed, subtests `removed` / `corrupted`).** A real ticket was
  constructed through public `Prepare` and the returned handle retained; then ONLY
  the fixture-local `request.json` was removed or corrupted. The retained handle
  still exists and is self-consistent (`ticket.SHA256()` unchanged). The prior
  proposal's inference "a reservation lacking a canonical `request.json` provably
  has no live ticket handle anywhere" therefore does NOT survive: the observed
  disk state cannot distinguish "interrupted Prepare, no handle escaped" from
  "handle escaped, file later lost". My earlier return-ordering argument
  (`verification.go:333-336` zero ticket on error) covers only handles that escape
  through Prepare's own return — it proves nothing about post-hoc file loss, and I
  no longer offer it as an existence proof.
- **What survives (executed).** The orphaned handle is INERT: `Reserve`,
  `Execute`, `Read`, and `Stop` with it ALL fail closed, and the journal
  byte-inventory is exactly what the corruption left (0 artifacts after removal;
  the corrupted `request.json` untouched) — every operation fails at the authority
  re-check (`verification.go:383-386`) before any write. The reservation still
  refuses a fresh `Prepare` (`"already reserved"`), so the MAJOR-2 brick stands.
- **Concurrent Prepare (EXECUTED variant, `concurrent-prepare-executed`).** Two
  goroutines racing `Prepare` with distinct run IDs: exactly one succeeded with a
  valid handle; the loser received a zero ticket (Version/Root/RunID/Request.Idea
  all zero) plus `"verification already reserved or unavailable"`; the published
  `request.json` binds the winner's run ID; the winner's handle stayed usable
  (`ReserveCapturedVerificationLaunch` OK).
- **Concurrent Prepare (SOURCE reasoning, labeled — not execution).**
  `withStateControl` holds the cross-process resource guard across the whole
  Prepare callback (`state.go:490-496`; `state_validation.go:43-50` invokes `fn`
  under the guard), so a second Prepare cannot interleave between the winner's
  `Mkdir` (`verification.go:283`) and its `request.json` write (`:327`); a crash
  between them leaves a request-less reservation whose Prepare returned a zero
  ticket (`:333-336`). So the only producer of a request-less reservation via
  public ordering is a crashed Prepare that published no handle — but, per the
  executed refutation above, an observer cannot distinguish that state from
  post-hoc loss under a live handle.

Consequence for the S2 quarantine design (prior proposal): its safety basis must
change from "provably no live handle anywhere" (refuted) to "any handle for this
reservation is provably inert — fail-closed at the authority re-check before any
write" (executed). The design outcome (quarantine cannot strand a usable handle)
survives on the weaker, proven basis. No quarantine was implemented here, and no
retry was granted.

## Failures retained / deviations

- One authoring-iteration build failure of my own new test (struct comparison
  against a type containing slices) — fixed in my test only; no pre-existing test
  or source was modified. No original failures were observed, fixed, or
  suppressed; the run above is the first and only execution of these tests.
- Deviation from plan: none. Both witnesses were constructible through public
  APIs; no impossible state was forced.

## Verdict summary

- MAJOR-1 (verifier-launch budget refusal brick): **PROVEN, current behavior** —
  strengthen the proposal from source-reasoned to executed.
- MAJOR-2 no-handle inference: **REFUTED as stated** — replace with the executed
  inert-orphan property; quarantine safety basis updated accordingly.
- Everything above describes the frozen tree `39107b138e77de064f21fec4aafa6073887c2783`
  as it behaves TODAY. Nothing here is an implemented recovery, an accepted fix,
  or a whole-goal acceptance signal. Proposed source changes require independent
  review before integration.
