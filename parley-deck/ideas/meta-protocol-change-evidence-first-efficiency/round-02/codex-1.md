---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
round: 2
date: 2026-09-05
---

# Cross-Review and Binding Corrections

I read all three peer round-01 artifacts after writing my own. The existing full
`go test ./...` baseline passed at code base `257ef8c` (no implementation edits;
app package 270.604s). This does not verify any proposed feature.

## Response to claude-1

Accept the single renderer, explicit applicability map, source-role distinction,
complete omission index, and keeping the comparative pilot on full context.
BLOCK the proposed disk-derived fixup count. PRIMARY: current
`internal/driver/cursor.go:47-62`, `driver.go:239-255`, and `impl.go:514-566`
explicitly preserve a monotonic CHARGED count outside participant artifacts.
Counting published headings was already rejected by review. Failed attempts cost
budget too. Reuse that safety counter and extend it to manual paths, never replace
it with headings. No self-authored "user" grant can buy capacity.

The skill source repository is in scope and writable in its own worktree. Package
release/global installation is not required to implement and test that source.
The three instruction paths must be changed together in source, with rollout
still gated. Source-role protocol changes are permitted through this meta idea;
global immutable-core publication remains human-only. Unknown applicability must
fall back to full, not merely include one unknown heading and drop dependencies.
Verbatim source binding and secret redaction cannot both hold: reject disclosure
of detected secret-bearing input, never redact it and claim the original hash.

## Response to hermes-1

Accept separation of provider, process, parser and observation outcomes. BLOCK
`RealHang` inferred from missing output and the exit-zero warning path that returns
readiness success. A buffered model can be alive beyond the observation window;
even a byte-producing process can be deadlocked. Name the observations honestly:
`deadline-no-output`, `deadline-after-output`, `process-exited-empty`,
`malformed-reply`, plus classified provider/process failures. A fake known-hung
fixture proves cleanup, not general diagnosis of all silent processes.

PONG in an echoed prompt is not successful agent readiness. Only exact semantic
assistant text extracted from recognized envelopes may satisfy PONG; preserve
the existing echo rejection. Unknown/malformed/empty/deadline outcomes require a
readiness-resolution gate with no automatic exclusion, not exit 0. Provider
failures remain blocking but distinct. Existing explicit user exclusion can stay
for definite unavailable failures; no silent quorum change. Use existing
`Spec.BuffersStdout` (PRIMARY `agents/discover.go:61-65`) to disable inappropriate
soft output watchdogs for declared buffered transports while retaining hard
timeouts. ACP activity is already tracked; do not count our heartbeat as progress.
No search of HOME for misplaced artifacts: absolute workspace paths and precise
searched-path diagnostics suffice without unrelated filesystem access.

## Response to kimi-1

Accept typed evidence, independent verifier, scope reconciliation, live six-agent
pilot, blinding, and telemetry-first sequencing. Correct two claims: existing
attempt_id is an integer reset per runner call, not globally unique; and a named
checks contract is not independent current-tree evidence by itself. PRIMARY:
`runner.go:539-568` and `app/driver_checks.go:51-97`. Add a UUID invocation ID and
retain the ordinal. Emit one normalized usage summary per attempt; dual-reading
both ACP aliases would double-count when both are emitted.

Bind evidence to reviewed commit plus tested code-tree digest, excluding only
defined evidence artifacts to avoid evidence-commit self-invalidation. Do not
pretend a generic command parser can detect a serial concurrency test. Require
structured executed-test counts where available and explicit independent
semantic verification for material claims. A serial-vs-barrier negative fixture
must demonstrate the distinction concretely. A manual launch THROUGH the shared
runner is measured telemetry, not honor-system reporting; imported assertions
remain separately labeled and cannot satisfy measured-cost gates.

## Common Integration Contract

Codex owns telemetry, runner/app shared integration, budgets, and pilot/report.
Claude owns a new `internal/protocolpacket` package, packet CLI helper/tests,
applicability map, skill source changes, and source protocol wording. Hermes owns
preflight/liveness and `runner/supervision.go` plus focused tests. Kimi owns new
`internal/evidence` and the driver checks/evidence integration, plus tests. Exact
claims are written in IMPLEMENTATION.md after consensus; shared app entrypoint
and runner prompt call sites are serialized by Codex.

Packet trial retains every preregistered threshold unchanged. The pilot's full
protocol context is identical across arms; experimental bounded roles and phase
caps are frozen and reported, never presented as unrestricted production runs.
Unknown usage never becomes zero; if a configured monetary ceiling cannot be
enforced, fail closed or reserve a documented conservative upper bound. The
two-confirmed-regression trajectory rule is opt-in and escalates without closure.
All five priorities, real experiments and the HTML remain required; no slice
alone completes this umbrella.

Additional PRIMARY inspection: `internal/runner/acp.go:387-396` emits only
`used`/`size` from ACP usage, and `internal/acp/protocol.go` defines no cost or
billed-token fields there. Those values describe context utilization, not a
monetary total. Do not rename that event and pretend it supplies price; preserve
it as context usage and normalize cost only when the provider actually reports
an attributable cost. This narrows the ACP migration claim in kimi-1 round 1.
