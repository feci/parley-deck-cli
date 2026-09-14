---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: db50dad638fc227c6d8d41e9de3ebef057c3d65c
source-manifest: 8938e8612351bd7efbbf85df589d0f4fdd541cff178daf9b42f8f6dec195abc2
status: observed-unresolved
---

# Full-size synthetic history exceeds the publication-sized deadline

Production Inspect over a native fixture with 128 charged attempts and a snapshot
of exactly 268,435,456 bytes refused after 30.002206667 seconds:
`after archive unavailable: context deadline exceeded`.
The same complete state passed Inspect with a separate 600-second deadline in
197.012724208 seconds and returned all 128 attempts. The 30-second attempt remains
failed; the longer read establishes that complete verification of this fixture
can succeed, without converting that timeout into a passing publication.

## What was executed

The Go test binary was compiled with a single native-only test-file overlay,
history_envelope_probe_test.go. No production source was overlaid. All 416
Go/module source hashes remained identical to the manifest above. The fixture
uses budget.Store.Reserve for 128 real local reservations and launches 128 real
`sh -c 'exit 0'` processes, recording each requested/started/terminal lifecycle
through telemetry. It then assembles resolution and continuation structures in
the test fixture and subjects them to production structural checks and Inspect.

This is not the production end-to-end 128-cycle workflow, a model treatment,
an independent participant's verification, a cold-cache benchmark or an
uncontended throughput claim. Other Go validation ran during the measurement.
No monetary/task-arm experiment ceilings are inferred from this synthetic test.

The source tree's four ballast files were sized so the canonical tar reaches
the exact 256 MiB maximum while every member remains within 64 MiB. The history
reuses identical source/archive references across its unchanged attempts and
acknowledged continuations. Distinct large archives and other resolution shapes
could have different costs; this is one valid maximum-count/maximum-size case,
not an exhaustive performance bound over every permitted history.

## Why this matters and what remains unmeasured

internal/trajectory/state.go:withStateResolutionCheck always calls
checkSourceSnapshots and the resolution checker before the requested operation.
checkSourceSnapshots inspects the baseline, every before/after archive and each
promoted continuation archive. internal/trajectory/unchanged.go:unchangedScope
also reads each before-scope and activation quorum through full archive readers.
The current code repeats these validations for identical references.

internal/runner/telemetry.go supplies a 30-second context to trajectory Finish.
The probe exercises production Inspect with that same duration, not a direct
runner Finish call or complete process-to-publication failure. Transfer to that
larger workflow still needs an explicit boundary check. This result does not
prove the original shared-volume timestamp cause or any process inactivity.

Potential changes need to preserve original charge/state authority and current
acceptance checks. Reusing validated archives within one guarded read requires
careful file-identity/stability handling. Separating factual terminal/stop
recording from permission to execute or accept may avoid requiring all historical
evidence merely to retain a real outcome. Neither option is implemented here;
increasing a deadline alone is not claimed to resolve the issue.

## Retained evidence

- .parley-runtime/history-envelope-probe-20260914/probe-result.json and probe.log.
- Native test source, overlay, compiled binary and original fixture directory:
  /var/folders/yt/p2sr23f12_qcfx_w2z5c1p4r0000gn/T/parley-history-envelope-fyfjg2_i.
- retained-fixtures.tar.gz preserves the complete fixture tree after termination,
  including source, archives, ledger and lifecycle records: 922,940 compressed
  bytes, SHA256 9dda0159795dac5aac0ca299f69133eec1de635ea418ab24bbd269af8cc58257.
- Native/shared log and retained-fixture archive bytes were compared equal.

The whole test passed in 245.771 seconds because it explicitly permits a proven
30-second timeout followed by complete verification. Its PASS is not a claim
that the 30-second production-sized deadline succeeded. F6 remains unresolved.
