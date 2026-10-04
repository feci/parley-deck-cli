# Provider quota recognizer support (2026-10-04)

Zcode is supported under the explicit owner deviation in
`parley-deck/inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_scope-reset-answer.md`.
This is a bounded stderr adapter, **not a native root-channel provenance claim**.

| Adapter | Automatic exclusion | Evidence and limits |
| --- | --- | --- |
| zcode | supported, owner deviation | Nonzero process exit, no valid artifact/later success, at least one JSON `responseBody`, every provider error record explicitly 429 allowance/account exhaustion with machine reset, agreeing resets, terminal turn failure, and reset >= 60 minutes. Malformed, duplicate, mixed, contradictory, truncated, quoted assistant/tool envelopes and incomplete records gate. Residual subagent-source ambiguity is expressly accepted by the owner. |
| codex | diagnostic-only | CLI 0.159.3: no established provider-only terminal discriminator. |
| kimi stream-json | diagnostic-only | Installed 0.42.0: role envelopes alone do not establish terminal provider provenance. |
| claude/text | diagnostic-only | Configured text output cannot distinguish assistant quotations from provider errors. |

`zcode-incident-2.stderr` is the four decisive lines quoted verbatim in
`source-context/EVIDENCE-quota-incidents.md`, incident 2. The retained evidence is an
excerpt, not a complete native terminal capture. Incident 1 retains the readiness
rows and gates only; its raw stderr is not present in the preserved sources. The
incident-1 replay fixture explicitly pairs those recorded readiness facts with the
recorded incident-2 provider body. No second primary stderr capture is claimed.

A zone-less display clock is accepted only alongside `reset_at` or `retry_after`
in the same JSON record. All interpreted instants must fit within one second;
possible display UTC offsets are restricted to -12:00 through +14:00 in 15-minute
steps. The provider display text is retained with the machine reset. The canonical
UTC reset remains the earliest machine/reset-duration instant. Observation times
come from receipt of the provider record, not eventual process teardown.

The assistant/tool quotation, mixed-error, malformed, short-reset, display-only and
contradictory fixtures are derived adversarial captures. Success/artifact precedence
is tested through invocation facts and the runner, since it cannot be inferred
from a string fixture. Synthetic fixtures and process helpers do not invoke a
provider, change CLI flags/models/credentials, or establish native provenance.

Source examined: `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`,
SHA256 `3e3433d90fa502e5d02498dfde6c2090df898331359bcfe5f3dbc9a1d00b685f`,
app 3.7.7 / installed package 3.7.7-13 / runtime 0.16.3.
The `AI_APICallError` SDK dump includes status/body/headers. Its untagged stderr
sink also receives callback errors. The terminal `Error: Turn execution failed`
wrapper does not authenticate the origin; this is precisely the accepted exception.

The native JSONL alternative remains unsupported: `turn.failed` / `core.runtime`
can carry provider causes but response summaries drop machine reset values; HTTP
429 summaries drop the allowance semantics; and nested subagent failures can share
trace/session/turn lineage. These are the decisive limitations identified in
review/round-02/claude-1.md. No native root-channel field was invented.
