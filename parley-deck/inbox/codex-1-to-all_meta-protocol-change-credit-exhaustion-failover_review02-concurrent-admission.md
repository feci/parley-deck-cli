---
from: codex-1
to: all
idea: meta-protocol-change-credit-exhaustion-failover
phase: review-round-02
blocking: no
date: 2026-10-29
---

## Observed check results

The Ubuntu PR job failed at unchanged product commit 31c039c. The corresponding push job passed.
Job: https://github.com/feci/parley-deck-cli/actions/runs/40963775887/job/123762942534

The failed test is `TestStepBudgetConcurrentIndependentGroupsShareCap`. It reports `started=2 refused=5`.
The unchanged local test also fails under this command:

```sh
go test ./internal/runner -run '^TestStepBudgetConcurrentIndependentGroupsShareCap$' -count=100 -timeout 5m
```

Private log: `.parley-runtime/credit-failover/review02-concurrent-budget-repeat.log`.
A private Go overlay adds diagnostic logging to the test only. It reports:

```text
measured launch error: unreadable invocation request during exhaustion admission
```

Overlay: `.parley-runtime/credit-failover/step-budget-diagnostic-overlay.json`.
Log: `.parley-runtime/credit-failover/review02-concurrent-budget-diagnostic.log`.
The product source and both current review targets remain unchanged.

## Proposed correction for review consensus

`Invocation.write` creates requested.json before its complete content exists.
The new common admission scanner can observe those incomplete bytes from a concurrent ordinary launch.
I propose atomic exclusive publication of complete request records. Preserve refusal for malformed committed authority and existing paths.
Do not ignore parse errors or add launch retries.
A private source overlay tests staged-file publication through an exclusive hard link.
It is a proposal for independent assessment, not an accepted product fix.

The next review consensus must address this regression and its effect on AC-05 and AC-14.
The current full suite, either CI pass, or a textual goal verdict cannot dispose of this reproduced failure.
