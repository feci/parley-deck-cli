# hermes-1 fix handoff — correct readiness envelope acceptance

- idea: meta-protocol-change-evidence-first-efficiency
- agent: hermes-1
- branch: feature/meta-protocol-change-evidence-first-efficiency/hermes-1
- code commit: ff3df32 (this fix)
- date: 2026-09-05
- status: fix committed; awaiting codex-1 non-owner negative-test re-run on the integrated tree

## Why this fix

Codex-1 ran a non-owner negative test on the integrated draft tree
(`go test ./internal/app -run TestReadinessRejectsUntrustedOrErrorPONG -count=1`)
against four outputs that my D7 classifier wrongly marked ready:

1. `{"role":"user","content":"PONG"}` — a user-role echo of the probe prompt.
2. `{"answer":"PONG"}` — an arbitrary string-valued key, not a known envelope.
3. `{"type":"result","is_error":true,"subtype":"success","result":"PONG"}` — a provider
   error whose subtype claims success.
4. `{"data":{"role":"user","content":"PONG"}}` — a user-role echo wrapped one level deep.

All four FAILED (Class:ready, Ready:true). The cause: envelope recognition accepted any
recognized key name, ignored the `role` field, and did not know `is_error`, so it treated
arbitrary/user/error objects as assistant success.

## The correction (internal/app/preflight_liveness.go only)

- A recognized envelope must now pass explicit role and error-status checks before its
  content is trusted (see "Supported schemas" and "Error/status handling" below).
- Arbitrary string-valued keys (`answer`, `pong`, `output`, `reply`, `response`, …) are no
  longer in the accepted content-key set; an object carrying only such keys is unrecognized.
- A role-tagged user/system/tool message — at the top level or one wrapper level deep — is
  never ready.
- A structured error/status envelope is now classified `provider-failure`, NOT
  `process-failure`: a JSON error is a provider status field, not arbitrary process text,
  so it stays blocking and is never auto-excluded (see "Error/status handling").

Corrected classes for the four cases (now all non-ready, pinned in my test table):

| input | class | ready |
| --- | --- | --- |
| `{"role":"user","content":"PONG"}` | malformed-reply (SawSentinel) | false |
| `{"answer":"PONG"}` | malformed-reply (SawSentinel) | false |
| `{"type":"result","is_error":true,"subtype":"success","result":"PONG"}` | provider-failure | false |
| `{"data":{"role":"user","content":"PONG"}}` | malformed-reply (SawSentinel) | false |

## Supported readiness verdict schemas (explicit, fixture-pinned)

A ready verdict is produced ONLY by one of these shapes; every positive shape has an
explicit test case in `TestClassifyReadiness`. This is NOT a claim that all CLI output
formats are supported — only this subset is recognized.

1. Plain exact PONG (plain-output adapters). The whole stdout, whitespace-trimmed, is the
   single token `PONG` (only blank lines may surround it). `{"content":"PONG"}` and friends
   need the JSON path; plain adapters that emit only `PONG` stay supported.
2. JSON assistant/result envelope — a single, standalone JSON object (one object on stdout).
   Accepted content keys (string value = the assistant/result text):
   `content`, `message`, `text`, `result`.
   Accepted wrapper keys (one level deep, inner object carrying a content key):
   `result`, `data`, `payload`, `response`.
   (`result` is both a content key and a wrapper key.)
3. Assistant-role message. If a `role` field is present (top level or inside the wrapper),
   it must equal `assistant` (case-insensitive) for content to count. A role-less object is
   treated as a raw assistant/result envelope.

## Non-ready inputs (explicit)

- `role` = user / system / tool → malformed-reply (echoed input), even when `content` is `PONG`.
- Arbitrary string-valued keys not in the accepted set (`answer`, `pong`, `output`, `reply`,
  `response` as a bare key, …) → malformed-reply.
- Malformed JSON (e.g. `{"content":"PONG"`), fences, bullets, echoed instructions, surrounding
  commentary → malformed-reply; a `PONG` substring is SawSentinel, never ready.
- More than one JSON object on stdout / line-delimited JSONL → NOT parsed; falls through to
  the plain-text path (never ready). See "Streaming JSONL consequence".

## Error/status handling (never ready, never excludable)

An object is an error/status envelope if, at the top level or one wrapper level deep:

- `error` is present and non-empty; or
- `is_error` / `isError` is boolean true; or
- `type` / `status` equals error/failure/fault/exception/failed (case-insensitive); or
- `success` / `ok` is boolean false.

Such an envelope is `provider-failure` regardless of any nested `subtype:"success"` or
`success:true` companion: `is_error:true` wins whatever the subtype claims. provider-failure
raises a blocking gate and is never reduced to an excludable `process-failure`; when the
body text names a known provider subclass (rate-limit/auth/billing/overloaded/model-not-found)
`ProviderClass` is set, otherwise it stays empty (reason `provider-failure`). Provider status
fields (`error`/`is_error`/`type`/`status`/`success`/`ok`) are read as failure signals, never
as content.

## Streaming JSONL consequence

`recognizeEnvelope` parses a SINGLE standalone JSON object from stdout. A streaming JSONL
reply — several objects on one line, or line-delimited objects — fails to parse as one object
and falls to the plain-text path, where it is never ready (at most malformed-reply with
SawSentinel if a `PONG` substring survives). Unsupported streams therefore resolve to
unresolved/non-ready, and an echoed input line inside a stream can never yield ready. I did
not add JSONL reassembly; a streaming adapter that wants a ready verdict must emit a single
recognized object (or plain `PONG`).

## Open concerns from the prior handoff — still valid, re-reported

- Unbounded capture buffers. `hostedPONG` still writes child stdout/stderr into unbounded
  `bytes.Buffer`s (`cmd.Stdout = &out`). A hostile or runaway-stdout agent can grow memory
  without bound up to the probe deadline. Not addressed by this fix; out of scope for this
  correction.
- Tail sanitizer removes controls, not credentials. `sanitizeTail` truncates to 256 bytes and
  strips control bytes other than `\n`/`\t`. It does NOT scrub secrets (an API key or token
  echoed into output survives into StdoutTail/StderrTail). Raw tails MUST NOT be exported as
  public or "safe" evidence. This fix does not change that.

## What I verified (and what I did not)

Run in this worktree on branch feature/meta-protocol-change-evidence-first-efficiency/hermes-1:

- `gofmt -w` + `gofmt -l` clean on both changed files.
- `go vet ./internal/app` — clean.
- `go test ./internal/app -run 'TestClassifyReadiness|TestClassifyReadinessProviderSubclass|TestHostedPONGRealChildFixtures' -count=1 -v` — all PASS (42 classifier sub-cases + 10 child fixtures).
- `go test ./internal/app -count=1` — ok (~60s).

I did NOT run the codex-1-owned `TestReadinessRejectsUntrustedOrErrorPONG` itself (it lives
in the integration worktree's `internal/app/preflight_evidence_test.go`, which I do not own);
I pinned the same four inputs in my own `TestClassifyReadiness` table and they are non-ready.
I have not independently verified my own correctness against the whole tree, and I do not
claim independent overall verification — that is the non-owner reviewer's job.

## Files changed in this commit

- internal/app/preflight_liveness.go
- internal/app/preflight_liveness_test.go

Only these two files plus this note. I did not touch preflight.go, supervision.go, the shared
IMPLEMENTATION.md, codex-1's preflight_evidence_test.go, or any other worktree. The prior
handoff (implementation-notes/hermes-1.md) is unchanged.
