---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: 7b351a2860d47a3e640955d550d93065ebbd6836
status: integrated-validated-awaiting-independent-review
---

# Refuse known protocol failure before higher-caller reservations

Claude N1 identified two gaps above the standalone runner's existing refusal
boundary: Driver.Advance charged its fixup step/cycle before invoking the runner,
and the verifier CLI created its one-use per-charge ticket before protocol refusal.

PrecheckProtocolLaunch shares the renderer and unstarted refusal retention with
beginProtocolLaunch. Success creates no invocation or reservation; failure records
an actual request and unstarted terminal through existing telemetry. It returns
no attestation for a later caller to reuse. Actual launch still renders again.

The mandatory internal ImplOps.PrecheckFixup method runs before reserveFixupCycle,
its ChargeStep/ChargeCycle and FixupCyclesPublished. The production adapter uses
the same withParticipants(implementer) and selectedAgents/resolveMapping path as
RunFixup. The embedded step wrapper promotes the precheck without charging it.
Initial zero-use bindings and contract/cursor initialization may still exist.

The verifier CLI prechecks before PrepareCapturedVerification. On refusal its
separate failed parent result retains the actual observed invocation and exact
terminal hash; its intended request_path has no request file/hash, and the
per-charge ticket directory remains absent. Restoring valid protocol permits a
new ordinary verification invocation for the original captured patch.

This is prospective prevention, not historical recovery. A source change after
a successful precheck still refuses at actual launch and can leave a prior caller
charge or prepared ticket unusable. Existing precharges are not refunded or
rewritten; no missing terminal, helper retry, receipt, custody or acceptance is
invented. Higher cross-review charging and N2 guard contention remain separate.

## Development evidence

The initial focused command FAILED (5.356s). Runner late-refusal and production
adapter tests passed, but two fixture assumptions prevented the driver and CLI
cases from reaching the intended boundary: the driver fixture's explicit cycle
maximum disagreed with the declared deliberation policy, and the app helper
incorrectly expected only one phase-specific cached full-context file.
Original source and logs are retained. A second command FAILED (9.662s): the
cache helper was corrected to use the renderer's exact BodyPath and both app
tests passed, but the driver fixture still used a mismatched policy maximum.
The next fixture correction uses the source-defined deliberation maximum of
five; it changes no production policy or gate. Accepted app/runner tests are not
rerun solely for that driver fixture correction.

Development source/logs: .parley-runtime/caller-precheck-development-20260914/.
The corrected driver check passed in 2.474s. This note supplies implementer
evidence, not independent acceptance or withdrawal of Claude's finding.

## Frozen validation

Manifest a9400d5630b351cc646ddf45ed6f54c8a2f8ae6ad27c02c465d9dd167e81450a
pins 424 Go/module files, including the exact ten changed source/test paths.
Full Go suite passed in 385.297s (31 passing packages plus CLI no-test-files skip).
Six-package race passed in 443.866s; vet passed. All four new top-level tests
passed in full and race. Compiled shared-volume runner passed in 7.176s, driver
in 2.947s and app in 79.724s. Windows CLI/trajectory/runner/driver/app cross-builds
passed; PE amd64 verified, Windows runtime unverified.

Five intended source-overlay failures catch omitting the driver precheck (1.782s),
CLI precheck (4.841s), actual refusal record (1.702s), actual-launch refusal
(2.138s), and using different production adapter participants (2.592s). The CLI
control's GOFLAGS overlay reaches its real child CLI build. None replaces a
passing case or weakens production checks. Both original fixture failures stay
retained with their own source/logs; accepted cases were not rerun solely for a
subsequent driver-fixture correction.

Final-verification.json under caller-precheck-final-validation-20260914 verifies
source/control identity, all four tests, the five negatives, exact native/shared
terminal logs, three unchanged participant reviews and historical September-5 HTML.
All commands captured output through native pipes and copied it after completion.
N2 historical-reader contention, old/racing precharges, helper recovery, full
participant acceptance/signatures and experiment/delivery gates remain open.
