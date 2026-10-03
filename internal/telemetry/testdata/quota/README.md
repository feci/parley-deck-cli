# Provider quota recognizer support (2026-10-04)

No adapter currently authorizes automatic exclusion. This is a release blocker for
FINAL AC2, not a substituted positive provenance fixture. Tests exercise the policy
refinement using explicitly synthetic native evidence; those are not CLI fixtures.

| Adapter | Automatic exclusion | Located evidence and remaining gap |
| --- | --- | --- |
| zcode | unsupported, diagnostic-only | Installed `zcode-app-cli` 3.7.7-13, runtime 0.16.3, source below. Raw SDK errors share an untagged stderr sink with callback errors. The recorded response is only an excerpt; no complete source-authenticated terminal envelope was recovered. |
| codex | unsupported, diagnostic-only | Installed CLI 0.159.3; no native source establishing a provider-only `turn.failed` discriminator was located in this invocation. |
| kimi stream-json | unsupported, diagnostic-only | Installed 0.42.0 binary and existing role-tagged stream parser; terminal provider error provenance is not established. |
| claude/text | unsupported, diagnostic-only | Configured text stdout does not distinguish an assistant quotation from a provider error. No argv changes made. |

Located zcode source: `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`.
SHA256: `3e3433d90fa502e5d02498dfde6c2090df898331359bcfe5f3dbc9a1d00b685f`.
Its extraction metadata names app 3.7.7, CLI 0.16.3 and upstream
`https://cdn-zcode.z.ai/zcode/electron/releases/3.7.7/linux-x64/ZCode-3.7.7-linux-x64.deb`.
Offsets below are character offsets in UTF-8-decoded source (Python `str`):

- 2536158: `AI_APICallError` / `APICallError` constructor assigns `statusCode`,
  `responseBody`, `responseHeaders` and symbol marker. A message substring alone
  does not establish this class.
- 4900711: `streamText` default `onError: ... console.error(q)` writes arbitrary
  received error objects. Its callback does not emit provider provenance tags.
- 4897609: stream forwarding includes provider `error`, but tool callback rejection
  also enqueues `{type:"error",error:F}`. Treating the error sink itself as
  provider-only would weaken FINAL's predicate.
- 10559036: tool executor errors use `tool.call.failed` / `core.tool.executor`;
  the generic catch and stderr alone do not bind the recorded excerpt to a model.
- 11353517: runtime logs `turn.failed` via `core.runtime`; this is a turn error,
  not proof that its origin was a provider.
- 13082700–13086200: `runPrompt` prints a successful `W.response` to stdout; catch
  prints `Error: ${T}` to stderr and returns 1. `Error: Turn execution failed`
  does not identify the terminal provider cause.

`zcode-recorded-response.json` is copied from the response-body excerpt in
`parley-deck/ideas/meta-protocol-change-quota-auto-exclude/source-context/EVIDENCE-quota-incidents.md:32`.
It is a message/reset fixture, NOT native terminal provenance. Its floating wall
clock has no timezone; until a verified recognizer establishes how it relates to
`reset_at`, strict reset handling also refuses the ambiguous expression.

No zcode/kimi participant was invoked, no model/provider/credential setting was
changed, and no credentials or full environment were copied. Resolving AC2 needs a
source-backed provider-specific terminal channel (or a fully attributable capture
whose source path rules out assistant/tool content), plus a positive and matching
adversarial native-channel fixture. Do not just enable a regex over these bytes.

Independent focused review: `parley-deck/ideas/meta-protocol-change-quota-auto-exclude/review/round-01/claude-1.md`
(R1–R3, MAJOR-1 and MINOR-1). claude-1 located the decisive limits: `runPrompt` prints the generic
terminal wrapper without upstream attribution; the SDK dump is not bound to the terminal turn and
is unframed. Its offline Node check demonstrated a different error message can embed the same dump
lines. These are the reviewer's source-backed findings; the organizer did not issue a self-verdict.
The owner must decide between a separately justified terminal-attributed transport and an explicit
AC2 deferral. Neither the all-diagnostic implementation nor synthetic policy fixtures satisfy AC2.
