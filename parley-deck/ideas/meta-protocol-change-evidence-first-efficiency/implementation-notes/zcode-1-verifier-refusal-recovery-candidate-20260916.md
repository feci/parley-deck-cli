---
agent: zcode-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-16
base-commit: 6962f2af9b6e46548588c5152285e6808f3d6f7f
kind: implementation-candidate (code written; tests NOT executed in this launch)
---

# Verifier-launch budget-refusal recovery candidate (zcode-1)

## Scope and ownership

Owned files (nothing else written outside `.parley-runtime/work` scratch, which
this launch did not need): `internal/trajectory/verification_recovery.go` (new),
`internal/trajectory/verification.go` (three amendments),
`internal/trajectory/verification_recovery_test.go` (new),
`internal/runner/verification_refused_recovery_test.go` (new), this note. The
integrated R1 recovery, telemetry files, protocol text, config and all
historical artifacts are untouched. No Git mutations, no Bash (explicitly
unavailable), no network, no subagents, no model calls, no fabricated results.

Provenance: PRIMARY (source) — my own reads of `verification.go`,
`captured.go`, `state.go`, `state_validation.go`, `reconcile.go`,
`unchanged.go`, `runner/telemetry.go`, `runner/trajectory_verification.go`,
`runner/launch_budget.go`, `app/trajectory_verify.go`, `telemetry/record.go`
and the existing test fixtures at base commit. Kimi's retained reachability
probe and slice plan (`.parley-runtime/input/`) are quoted context/testimony,
not evidence I verified by execution; I did not copy the probe's PROVEN verdict
as independent support, and the probe file itself is not part of this tree.

## Implemented mechanics

**Admission (verification_recovery.go).** One explicit public operation,
`RecoverCapturedVerificationLaunch(ctx, ticket, refusedInvocation, newInvocation)`,
runs under the full `withVerification` authority (canonical `request.json`,
charge/state, activation quorum). It requires: the journal to be exactly
`{request.json, launch.json}` — no claim, receipt, prepared, stop, step,
process or unknown record (a stop written before recovery blocks it); the
launch reservation pinned to `refusedInvocation` with the ticket digest; and
the refused invocation's retained `requested.json`+`terminal.json` proving a
pre-start refusal: metadata bound to this ticket's run/idea/verifier with phase
`trajectory-verification`, headless mode, terminal `failed`/`budget_refused`,
`StartedAt`/`PID`/`ExitCode` null, no `started.json`, sane chronology. All
checks precede the single O_EXCL write of `recovery.json`, which binds ticket
SHA, prior launch SHA, refused requested/terminal SHAs, refused ID and the
intended new invocation. Replay is unsupported by design: repeated or
conflicting apply fails closed on the retained artifact; no expected-preview
authority exists (adopting Kimi's open preview+SHA variant would be a separate
reviewed change). Recovery executes no model, no criterion; it never deletes,
overwrites or reuses the reservation and never touches any ledger.

**Effective launch (verification.go).** `readVerificationLaunch` resolves the
effective authority: without `recovery.json`, exactly the old behavior; with
it, the binding (prior launch SHA, refused ID, digests) is validated, the
refused invocation returns a dedicated superseded error **before any write**
(stop/read/execute), and only the bound invocation resolves — to the recovery
artifact itself, so claims, receipts and stops bind the complete lineage.
`ReserveCapturedVerificationLaunch` becomes recovery-aware: after recovery the
artifact already reserves its bound invocation, so that one invocation is
admitted with no write (the subsequent launch is separately checked here and
charged by its own budget boundary afterwards); every other invocation is
refused as before. The journal inventory admits `recovery.json`; unknown names
still refuse the whole journal, so pre-recovery binaries (whose literal
allowlist lacks the name) fail closed on recovered journals — asserted by
simulating that allowlist, since no old binary exists here.

## Unexecuted test intentions (I wrote them; the coordinator must run them)

Trajectory API (`verification_recovery_test.go`): happy path — real
telemetry.Begin→Reserve→Finish ordering reproduces the refusal; recovery
retains request/launch bytes and the fixup ledger byte-identical (count 1);
stale stop/read/execute fail superseded with no `stop.json`/`claim.json`;
helper execution completes (4 steps, receipt validates, regression derived);
unknown `quarantine.json` refuses; effective stop still works; Prepare still
refuses reissue; `RequireResolved` still fails. A negative table (started
lifecycle, exit code, wrong class/run/idea/agent/phase, missing
request/terminal/launch, claim/receipt/stop/stray present, wrong refused ID,
same/implementer invocation, never-reserved). Repeated/conflicting/concurrent
apply (exactly one winner; loser cannot execute). A forged structurally
contradictory artifact refused by recovery and every reader. Old-reader
allowlist simulation. A seam test proving fresh `PreviewReconciliation` on a
recovered journal fails closed at the launch lineage ("successful observed
parent terminal binding is unavailable").

Runner (`verification_refused_recovery_test.go`, local `/bin/sh` fake CLI, no
model): a real `WithLaunchBudget`-denied verifier launch through public
`RunMeasured` (pre-start `budget_refused`, no spawn, reservation bound, cycle
ledger untouched); explicit recovery; a real allowed retry whose own generated
invocation differs is refused before its budget store and before spawn; stale
stop/read/execute refused; the helper half completes under the bound invocation
with the receipt validating and accounting unchanged.

**Not run:** `go build`, `go vet`, and every test above. No test of this
launch has any executed result; nothing here is PASS evidence.

## Observed tool limitations

Bash was unavailable, so no compilation, no test execution, no `gofmt`. All
code was checked only by reading; compile or assertion failures are plausible
and must be fixed before integration. The 30-minute ceiling and no-subagent
rule prevented deeper cross-package verification (e.g. `procctl.Alive`
behavioral checks).

## Unresolved scope (honest limits)

1. **The resolution seam is open.** `deriveParentEvidence` (reconcile.go, not
   owned) reads `launch.json` directly and requires that invocation's
   successful terminal. Post-recovery, launch.json still pins the refused
   invocation, so verifier-run resolution refuses — the trajectory stays
   unreconciled until that reader learns the recovery lineage. My seam test
   asserts this fail-closed behavior rather than weakening it; wiring it (and
   the attended app CLI control that hands the bound invocation to a real
   verifier relaunch) is the required follow-up.
2. **The positive "subsequent launch rides the recovery" leg is untested at
   the runner level** — the runner generates invocation IDs internally, so only
   the API-level idempotent reservation and the runner-level mismatch refusal
   are covered.
3. Custody/inactivity: admission binds retained bytes only; same-UID edits and
   forged-but-structurally-consistent digests are outside byte evidence (the
   runner's binding checks remain the defense). Step-session accounting for the
   refused launch is retained spent history, not refunded.
4. Kimi's withdrawn no-handle inference and the overwritten build-log
   limitation remain theirs to correct; I claim nothing about either beyond
   what my own reads show.

This candidate alone does not complete the audit slice, let alone D-series
acceptance; it is offered for inspection and actual test execution first.
