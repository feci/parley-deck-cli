---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 2
date: 2026-10-04
reviewed-commit: 40e44d4
responding-to: [claude-1/review/round-01]
---

> **SUPPLEMENTAL STRUCTURED-CHANNEL EVIDENCE REVIEW: narrow early Phase-5 feedback.** This is not a full
> Phase-6 acceptance review and not a fix-up cycle.
>
> **Scope.** It reviews codex-1's new zcode evidence
> (`source-context/provenance-review/codex-1-structured-channel-evidence.md` and
> `zcode-native-turn-failed-scrubbed.json`) against FINAL §4, §13.8, AC2, AC3, AC4 and AC18. It also
> disposes of my round-01 findings. **All other ACs remain unreviewed.**
>
> **What it is not.** No consensus signoff is given here. No AC is deemed passed because a candidate channel
> exists. codex-1 supplied evidence, not a verdict. Only codex-1 and I participate.
>
> **Protocol attestation.** `context_mode=full`,
> `source_sha256=8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`,
> `packet_sha256=8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`,
> `fallback_reason=null`. I read the body (1484 lines) in full, and `shasum -a 256` matches.
>
> **Tree.**
> - Reviewed commit `40e44d4`.
> - The two evidence files were untracked at launch. They were then committed in `8bdf8d8` with identical
>   blobs (`1b49d7a51c8b` and `ce99f4374f27`).
> - HEAD moved to `236084b` during this review. `git diff 40e44d4 236084b -- internal/ cmd/` is empty, so
>   every code statement below holds at both commits.
>
> **Vendor source.** `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`, sha256
> `3e3433d90fa502e5d02498dfde6c2090df898331359bcfe5f3dbc9a1d00b685f` (the same file codex-1 cites).
> - Package: `zcode-app-cli` 3.7.7-13, `cliVersion` 0.16.3. Its files are dated 2026-08-17, and it is the
>   only `zcode` on PATH.
> - Offsets are character offsets from a Node substring scanner. In these ASCII regions they equal byte
>   offsets.
> - I read the original surrounding source, not only codex-1's extracts.
>
> **Log access.** I read only `~/.zcode/cli/log/zcode-2026-10-02.jsonl`, through scripts that print
> whitelisted structural fields, salted-hash pseudonyms for ids, and zcode's own fixed strings.
> - Not printed or copied: upstream message bodies, headers, prompts, configs, credentials.
> - Not done: invoking zcode, kimi or any provider; changing any source, config, provider, model or
>   credential setting; using helpers or a browser.
> - Scripts are kept in `.parley-runtime/quota-implementation/provenance-review/claude-1-round-02/`
>   (gitignored).
> - I hit no quota, credit or authentication error, so no blocking owner note is owed.

## Summary

The per-invocation JSONL log improves on stderr for framing and attribution. Each record is one JSON-encoded
line, and the terminal `turn.failed` record carries a zcode-assigned `source: provider`, provider and model.

It still fails the narrow question, because it does not keep the decisive provider data:

- zcode's business-error summarizer keeps only the *names* of `reset_at` and `retry_after`, never their
  values.
- For AC2's own recorded HTTP-429 shape, the terminal record contains no exhaustion text at all. I ran this
  check on the real log: every `turn.failed` from 21:00 to 21:07Z reads `model_rate_limited | 429 →
  AI_APICallError`, with an empty message.
- Binding the root turn by traceId fails, because subagents share the traceId. A fail-closed binding needs
  pinned negative rules.

**Position:** MAJOR-1 remains and now covers this channel too. MAJOR-2 changes: my round-01 fix is not allowed
by FINAL as written. It needs an owner interpretation, whose boundary is stated below.

## Refutation attempts

All verdicts here are mine as a non-owner of codex-1's evidence, of FINAL's ACs and of the implementation.
Results that are my own new claims carry their provenance tag and **no** verdict word (§15.1). They are open
to a non-owner verdict.

### Q1: can the unchanged `--prompt` invocation's per-invocation JSONL supply attribution and binding?

**S1. Is the log reachable per invocation?** PRIMARY, located source.

- **Default directory.** `c2` (≈7542547) resolves `o=e.logDir??e.env?.ZCODE_LOG_DIR??Qti()`. `Qti`
  (≈7543188) is `~/.zcode/cli/log`.
- **Env flow.**
  - The CLI main builds `let c=Cje(t.env??process.env)` (≈13125592) and then
    `logger:t.logger??c2({env:c})` (≈13126430).
  - `runPrompt` passes the same env to `createZCodeApp`, whose bootstrap also calls `c2({env:e.env})`.
  - The sanitizer list `QJt` (≈471315) removes only proxy, CA, broker and telemetry keys, not
    `ZCODE_LOG_DIR`.
  - The launcher copies `process.env` into the runtime child (`bin/zcode.js:1408`, spawn at `:1477`).
- **Constraints a reader must handle.**
  - The file name uses the **local** date (`tri` ≈7543637, `oH` ≈7538494), so one invocation can write two
    files.
  - `log()` wraps `appendFileSync` in `catch{}` (≈7544854), so a missing record is silent.
  - The redactor replaces anything deeper than 8 levels with `[Redacted:DepthLimit]` (≈7541889). The native
    record already nests the provider message 6 to 7 levels deep.
- **Verdict on codex-1's claim** ("installed source also supports a scoped `ZCODE_LOG_DIR`"): **CONFIRMED**.
  The root process honors it.
- **Scope note.** It is a launch-environment change to the adapter; the argv is unchanged. Absence,
  truncation or a depth-capped value must gate, never default.

**S2. Unframed message.** PRIMARY, located source.

- Each record is `JSON.stringify(entry)` plus one newline (`WVr` ≈7540265, `log()` ≈7544854). Newlines
  inside messages are escaped.
- Result: the round-01 R3(c) forged-line attack does not carry over into the file.
- Residual: there can be several writers to one file, and write errors are swallowed. A reader must parse
  every line strictly, and any unparseable or partial line gates.

**S3. Tool origin.** PRIMARY, located source.

- The tool executor's catch (≈10559134) returns a failed tool result and does not rethrow. It logs
  `event:"tool.call.failed"` (module `core.tool.executor`) with the **full error chain**.
- So tool output does not become the root `turn.failed` cause on this path. But provider-attributed chains
  do enter the same file through non-terminal records (S4).
- Result: only `event=="turn.failed" && module=="core.runtime"` may carry semantics. A reader that takes
  semantics from any other record is refuted.

**S4. Nested call.** PRIMARY, located source.

- **Subagents keep the parent's traceId.** The subagent lifecycle `Ebn` (≈11196979) builds:
  - `p=dl(t.trace,…)`;
  - `m=dl(p,{sessionId:Uu("subagent_"+id),turnId:t.turnId,attributes:{agentId,agentType,parentSessionId,parentToolCallId}})`.

  `dl` (≈651594) keeps `traceId:e.traceId`.
- **Every turn logs under that trace.** Each runtime sets `rootTraceContext=o.traceContext??hm(…)`
  (≈11786475). Each turn uses `_d(this.rootTraceContext,…)` (`R4n` ≈11644990). `PL` logs `turn.failed`
  (≈11652306). So a subagent's turn writes a genuine `source: provider` `turn.failed` into the same file,
  **with the root's traceId**.
- **Failed subagents reach the log twice.** A failed subagent's core error is re-thrown unchanged
  (`yi(b)?b:Ae(ye.ToolExecutionFailed,"Explore subagent failed",…)`, ≈11192846). S3 then logs it as
  `tool.call.failed`.
- **stderr gives only the traceId.** On failure, `runPrompt` prints only `Error: … (traceId: ${m})`, with
  `m=f.traceId` (≈13084978). It prints no sessionId.
- **Result: binding by traceId alone fails.** A fail-closed binding exists only as pinned negative rules on
  undocumented internals. All of these must hold, and anything else gates:
  1. The private directory holds exactly one bootstrap trace (`bootstrap.app.startup.*`,
     `startupKind:"zcode_app"`).
  2. The root session is the one whose sessionId lacks the `subagent_` prefix and whose context lacks
     `agentId`, `parentSessionId` and `parentToolCallId`.
  3. There is exactly one root `turn.failed`.
  4. The exit code is nonzero.
  5. Any stderr traceId equals that trace.
- **UNVERIFIED:** whether a nested `zcode` started by the Bash tool inherits `ZCODE_LOG_DIR`. Rule 1 would
  gate it.

**S5. Unrelated failure.** PRIMARY, located source.

- **(a)** A subagent hits exhaustion, then the root turn fails for an unrelated reason. The root record
  carries the unrelated cause. Only a reader that scans all records misattributes it.
- **(b)** Attempts are exhausted, then the terminal error has a different class. Per-attempt
  `model.request.failed` (`FFr` ≈7100989) and `model.retry.delay.resolved` (`sI` ≈7151754) records carry
  their own status, so they cannot supply semantics.
- **(c)** The adapter context carries `modelId` and `providerId`, so the recognizer must also require them
  to equal the configured model. Whether an auxiliary model call can fail a root turn is UNVERIFIED; require
  the match anyway.

**S6. What "provider" means.** PRIMARY, located source.

- `aYo` (≈7152924) returns `provider` when `t.statusCode!==void 0||NS(e)`, and also as its final default.
  It labels the HTTP peer.
- OmniRoute is in the path: the native record lists `x-omniroute-route-class`. So the label does not prove
  upstream text. FINAL §4 requires the upstream message for a gateway pass-through, and source evidence for
  gateway phrases.
- **Observation:** the "(reset after 53h 41m 9s)" suffix has the same format as the gateway-local
  "Unavailable (reset after 55m 29s)" texts of incidents 4 and 7. Whether OmniRoute composes or decorates
  this text is UNVERIFIED (open question 3).

**Q1 result.**

- **Attribution** is zcode-assigned and structurally present, within S6's limit.
- **Binding** is achievable only with S4's pinned rules plus adversarial fixtures.
- Neither rescues AC2 (Q2 below).

### Q2: does this channel keep all decisive provider data?

**D1. The summarizer drops reset values by construction.** PRIMARY, located source.

- `tzr` (≈4437060) keeps `keys` (up to 20) plus the primitive values of `success`, `code`, `error_code`,
  `msg`, `message`, `request_id` and `requestId`. Inside `error` it also keeps `type`.
- Both business-error paths use it:
  - the SSE `event: error` path (`t7o` ≈4435177 → `YDr` ≈4434807);
  - the non-ok JSON path (`H9o` ≈4431884 → `GDr` ≈4432436 → `jV` ≈4432879).
- `bD` (≈4436800) collapses whitespace in `providerMessage` and truncates it at 1000 characters
  (`ghe=1e3`).
- **Verdict on codex-1's claim** ("lists reset field names without necessarily preserving their values"):
  **CONFIRMED**, and the located source goes further. The values of `reset_at` and `retry_after` are
  **never** kept for a business error.

**D2. The native record as captured.** PRIMARY, executed (`logprobe.js line 16560`, plus a field comparison).

- **Shape and contents.**
  - The raw line has the same key tree as codex-1's scrubbed fixture, with the same response-header names.
  - Fifteen compared fields are equal, including every message, `source`, `providerId`, `modelId`,
    `responseStatus`, `transport`, `attempt` and `timestamp`.
  - The raw message contains "53h 41m 9s".
- **Absent from the raw line:**
  - `22:14:57.003`;
  - a numeric `retry_after`;
  - `176930`;
  - a `"reset_at":"` value;
  - a `retry-after` header.

  The inner summary holds only `keys` and `message`. The final `model.retry.delay.resolved` of that trace
  reports `retryAfterSource:"missing"`.
- **Chain.**
  1. Wrapper `UNKNOWN_ERROR`.
  2. `AiSdkModelAdapterError`: `model_request_failed`, `source: provider`, `reason: unknown`, glm-5.3,
     anthropic kind.
  3. `ProviderBusinessError`.

  The context has no `parentSessionId`, so this is a root session.
- **Verdicts on codex-1's claims (all CONFIRMED):**
  - The record exists at `:16560`.
  - The fixture keeps the keys and nesting and invents no error fields.
  - The record reads 53h 41m 9s, not AC2's 49h 8m 50s.
  - It carries `source=provider`, provider and model identity, and the provider-supplied message. Its
    upstream status is subject to S6.
  - The values were missing before codex-1 scrubbed anything.
- **Arithmetic (mine):** 16:33:48Z + 53h 41m 9s = 2026-10-04T22:14:57Z, which is the `reset_at` of AC2's
  body. The record itself cannot show this agreement.

**D3. AC2's recorded shape leaves no exhaustion text.** PRIMARY, executed plus located source. This is my
claim.

- **Only three records match.** In the whole 2026-10-02 file, only three records contain "Limit Exhausted":
  lines 16557, 16558 and 16560, all at 16:33:48Z.
- **The AC2 window.** From 21:00 to 21:07Z, which covers the AC2 preflight (23:04 CEST) and the direct
  probe (23:06 CEST), more than ten invocations with `startupKind zcode_app` and `mode yolo` hit HTTP 429.
- **Every `turn.failed` in that window reads:**
  1. Wrapper: 21 characters, "Turn execution failed".
  2. `AiSdkModelAdapterError`: `model_rate_limited`, `source: provider`, 429, `rate_limited`, no body
     summary, message "Provider rate limited the model request.".
  3. `AI_APICallError`: message length **0**.
- **The only machine value.** Per-attempt `model.retry.delay.resolved` records carry `retryAfterHeader`
  (for example "177256"), with `retryAfterSource: provider_header`.
- **Why the text disappears:**
  - The anthropic provider's error schema requires `type:"error"` and `error.type` (`SAo` ≈3039064). The
    recorded body has neither.
  - So `createJsonErrorResponseHandler` (≈2666643) falls back to `message:o.statusText`, which is empty
    under HTTP/2.
  - `Yb` (≈4973045) maps 429 to the fixed text.
  - `rw` and `lYo` (≈7152369, ≈7153871) add a body summary only for a business error.
  - `mCt` (≈7540975) serializes only name, message, code, type, a plain `context` and `cause`. It never
    serializes `responseBody` or `responseHeaders`.
- **Result.** For AC2's recorded shape, this channel holds a bare 429 plus a long Retry-After in non-terminal
  records. FINAL §4 says both "never qualify".
- **Not established:** which trace was the direct probe. No record in the window carries the text.

**D4. Producing version.** PRIMARY, executed. This is my claim.

- Both shapes come from traces whose `mcp.server.connected` records report `mcpClientVersion 0.16.3`. That
  is the CLI constant `AE="0.16.3"` (≈13119757).
- `~/.zcode/cli/version.json` records `checkedVersion 3.7.7-13`. The separate desktop app
  `/Applications/ZCode.app` is 3.14.4.
- So the record is consistent with the installed CLI runtime. It does not prove the build is byte-identical
  (residual).
- The two shapes reflect two different gateway responses (an HTTP 200 SSE error versus HTTP 429 JSON), not
  two versions.
- **Verdict on codex-1's claim** ("version at capture is not established solely by the installed-source
  hash"): **CONFIRMED** as stated. D4 adds a consistent record-level marker.

**D5. Other existing channels.**

- **App-server payload.** PRIMARY, located source.
  - `EA` (≈8519056) builds the `TurnError` session event, which is also the app-server payload. It keeps
    messages (up to 500 characters), `detail` strings and attribution.
  - The detail strings come from `rbi` and `nbi` (≈8521520, ≈8521653): provider, model, request, code,
    reason, status, retryable.
  - So it carries strictly less than the JSONL.
- **`--json` and stderr.** PRIMARY, located source.
  - `runPrompt --json` writes only on success.
  - The stderr SDK dump (round-01 R2/R3) is the only observed carrier of the shape-B `reset_at`, and it
    stays rejected.
- **Status sink.** UNVERIFIED. `Pp` (≈7098996) passes the raw failure object to an in-process
  `statusSink.publishFailure`. Whether any persisted, invocation-bound consumer exposes it is not
  established.
- **Model-IO debug recorder.** PRIMARY, located source. `qwt` (≈7142959) serializes errors as
  `{name,message,stack}` (`i7r` ≈7149927) and records full request messages. It is out of scope.

**D6. Correlation.** PRIMARY, located source.

- The terminal adapter context carries `requestId` and `traceId` (`rw`). Per-attempt records carry the
  `rI` fields (≈7099967): sessionId, turnId, requestId and attempt.
- So correlating the final attempt with the terminal record inside one process trace is **provable**.
- It adds only Retry-After (shape B) or nothing (shape A), so it is **not sufficient**.

**Q2 result (mine): no.**

- In shape A the `reset_at` and `retry_after` values are lost. A contradiction between the display text or
  duration and the machine resets therefore cannot be detected.
- In shape B, AC2's own event, there are no exhaustion semantics at all.
- No correlatable record restores either.

### AC-level attempts (this channel only)

- **AC2:** fails on this channel (D3). Shape A would also need MAJOR-2's interpretation, and even then fails
  (D1).
- **AC3:**
  - The channel's own data puts shape B in the "bare 429" and "long Retry-After alone" negatives.
  - The assistant-quoted and tool-emitted negatives map to S3 and S4.
  - Native negative fixtures now exist: shape A at line 16560; shape B at lines 16702, 16883 and 17066.
  - Positive and adversarial fixtures are still required.
- **AC4:** unchanged at `40e44d4`. zcode is "unsupported, diagnostic-only"
  (`internal/telemetry/testdata/quota/README.md:9`). `internal/telemetry/quota.go:60–62` states that no
  production recognizer constructs `nativeQuotaError`. This holds for the reviewed tree only.
- **AC18:** a terminal-only, root-bound reader (S3 to S5) keeps tool output and subagent records out.
  Elapsed time and disagreement are not in play here. AC18 is not deemed passed.

### Prior findings

| Round-01 item | Now | Basis |
|---|---|---|
| MAJOR-1 | **Remains**, extended to the JSONL channel | D1–D3, S4 |
| MAJOR-2 | **Changed.** The suggested fix is withdrawn as an implementer choice, an owner interpretation is needed, and the parser overrun is split out as MINOR-2 | FINAL §4.5, ratified §9.0, the `quota.go` comment |
| MINOR-1 | **Remains**, extended | D1, D3, S4 |
| R4 (MCP stderr injection refuted) | Unchanged | — |
| R6 | **SELF-CORRECTION** (below) | D2, D3 |

**SELF-CORRECTION (R6).**

- **Replaced statement.** Round-01 R6 said that `grep -rlF 'Weekly/Monthly Limit Exhausted'` over
  "`~/.zcode/cli/logs` and `~/.zcode/logs` found no raw capture".
- **Correction.** I searched the wrong directory. The log lives in `~/.zcode/cli/log` (singular, `Qti`), and
  it holds a native structured record of a *different* event (shape A, 16:33:48Z).
- **Weakening, effective now.** I withdraw the claim that no native structured record exists.
- **New claim.** No record in AC2's own event window carries the exhaustion text, so a byte-faithful native
  positive for AC2 still cannot be rebuilt from recorded evidence. Provenance: D3, PRIMARY. It is my own
  claim and needs a non-owner verdict.

## Findings

### [MAJOR] MAJOR-1 (remains, extended): AC2 is unattainable on the configured zcode path, the JSONL channel included; the owner scope decision is still required

- **What is wrong.** Round-01 refuted the stderr text path, and the per-invocation JSONL candidate fails
  too:
  - For AC2's recorded HTTP-429 shape it records a bare 429 with an empty upstream message (D3).
  - For the business-error shape it drops `reset_at` and `retry_after` by construction (D1, D2).
  - Binding by traceId fails (S4).
  - The app-server payload carries even less (D5), so switching transport to it does not help with the
    installed zcode.
- **Does it fit FINAL §13.8?** §13.8 lets the implementer establish provenance per adapter, or ship the
  adapter unsupported.
  - A private `ZCODE_LOG_DIR` with a terminal-only, root-bound reader could be an adapter detail. But it
    cannot produce AC2's decisive data, so it is **not a viable route to AC2**.
  - At most it can enrich diagnostics, such as the attribution shown in the owner note. Even that needs
    S4's rules and fixtures.
- **Why it blocks.**
  - AC2 names a zcode recognizer and this exact message.
  - AC19 and AC21 need both stages plus current-tree criterion evidence.
  - So the gap cannot be absorbed silently (§4 Phase 5).
- **Smallest concrete scope choice (for the owner).** Ship zcode as diagnostic-only in this idea. Record an
  owner ruling on AC2 for zcode as a deviation, either:
  - defer AC2 to a linked follow-up idea, whose precondition is a zcode-side channel that keeps the provider
    body or reset fields for the terminal failure (an upstream zcode change, or a separately reviewed
    transport); or
  - amend AC2's adapter or message.

  I treat neither as approved. It is the owner's call.
- **Preconditions if a JSONL recognizer is ever pursued.** All of these are needed before any eligibility:
  - S4's binding rules;
  - the model and provider match (S5c);
  - gateway pass-through evidence (S6);
  - treating "reset key present, value unavailable" as unparseable;
  - native positive fixtures.

  It also needs adversarial fixtures for:
  - a subagent exhaustion followed by an unrelated root failure;
  - a `tool.call.failed` record carrying a provider chain;
  - two bootstrap traces in one directory;
  - truncated, partial or missing records;
  - two local-date files;
  - a depth-capped value;
  - a model mismatch;
  - the AC3 quoted and tool-emitted texts.

### [MAJOR] MAJOR-2 (changed): AC2's recorded message gates under FINAL §4.5 as written; my round-01 fix needs an owner interpretation

- **Reconciliation.**
  - FINAL §4.5: "A contradictory, past or unparseable reset gates. It is never relabeled `unknown` to gain
    eligibility."
  - Ratified §9.0: "Past, contradictory or unparseable resets gate; they are never relabeled unknown."
  - FINAL does not say whether a zone-less wall-clock time is "unparseable".
  - The committed classifier takes the strict reading: "Absolute values require a timezone; no
    machine-local timezone inference" (`internal/telemetry/quota.go` at `40e44d4`).
  - Under that reading, AC2's message ("reset at 2026-10-05 06:14:57 (reset after …)") gates on **every**
    channel. That contradicts AC2's positive, which is a tension inside FINAL that only the owner can
    resolve.
- **Withdrawn.** My round-01 suggestion to "ignore absolute times that carry no zone as non-decisive" is not
  allowed by FINAL as written. Applying it would be a silent reinterpretation unless the owner records an
  interpretation. If the ratified text itself must change, it needs a §7 amendment.
- **Boundary for the owner, if AC2's message should count as positive.** A zone-less display clock is
  non-decisive only when **all** of these hold:
  1. The same native terminal record carries at least one zoned machine reset whose value is present: RFC3339
     with an offset or `Z`, or `retry_after` / `Retry-After` in numeric seconds.
  2. Every decisive value agrees within the fixed tolerance, including any "reset after" duration measured
     from observation.
  3. The display clock equals the machine instant under some real UTC offset, from −12:00 to +14:00 in
     15-minute steps. Here it is +08:00.
  4. The raw display string is kept in the evidence and in the marker.

  In every other case it gates:
  - there is no machine value;
  - machine reset keys are present but their values are unavailable;
  - no offset reconciles the clock, which is a contradiction;
  - the clock is the only stated reset.

  Under this boundary nothing becomes "unknown" or "no reset". The 24-hour no-reset branch never applies
  while a reset is stated.
- **The boundary does not help on the JSONL channel** (D1, D2). Shape A lists `retry_after` and `reset_at` as
  keys without values. So the clock stays undecidable, the hidden values cannot be checked for contradiction,
  and the event gates.

### [MINOR] MINOR-1 (remains, extended): the zcode support row omits the decisive reasons

- **Where.** `internal/telemetry/testdata/quota/README.md:9` at `40e44d4`.
- **What to add:**
  - round-01 R1, R3(b) and R3(c);
  - `tzr` dropping the reset values (≈4437060);
  - shape B carrying no exhaustion text (≈3039064, ≈2666643, ≈4973045, ≈7152369; recorded 21:00–21:07Z);
  - the traceId shared with subagents (≈651594, ≈11196979).
- **Why.** Without them, the diagnostic-only status will be re-argued on weaker grounds.

### [MINOR] MINOR-2 (new, split from round-01 MAJOR-2): `resetAtText` runs into the parenthetical

- **Where.** `internal/telemetry/quota.go:74` at `40e44d4`:
  `` `(?i)\breset(?:s| will reset)? at ([^);\n]+)` ``.
- **What is wrong.** The pattern stops at `)` but not at `(`. On the recorded message it captures
  `2026-10-05 06:14:57 (reset after 49h 8m 50s` (round-01 executed check, `claude-1-reset.go.txt`).
- **Effect.**
  - It fails closed, so there is no false positive, but it records the wrong gate reason.
  - It would also gate a zoned reset that is followed by a parenthetical.
- **Fix.** Stop the capture at `(`, or parse each reset expression with its own anchored pattern.
- **Fixtures to add:**
  - "reset at 2026-10-04T22:14:57Z (reset after 49h 8m 50s)": both values are parsed and compared.
  - The recorded message: gates with the reason "zone-less absolute reset", unless the owner adopts
    MAJOR-2's boundary.

## Open questions

1. **Owner, MAJOR-1.** Is AC2 for zcode deferred to a linked follow-up idea with a zcode-side channel as its
   precondition, or amended? Either way, record it as an owner-authorized deviation.
2. **Owner, MAJOR-2.** Should the stated boundary for zone-less display clocks be adopted, or should the
   strict reading stay? Under the strict reading, AC2 as written cannot pass for this message on any channel.
3. **Gateway.** Does OmniRoute pass the z.ai message through verbatim, or does it compose "[model] [status]:
   … (reset after …)"? The suffix format matches incidents 4 and 7. FINAL §4 needs source evidence before
   either part counts as a positive form.
4. **UNVERIFIED:**
   - which consumers receive `statusSink.publishFailure` (≈7098996);
   - whether children started by the Bash tool inherit `ZCODE_LOG_DIR`;
   - whether an auxiliary model call can fail a root turn.
5. **Other reviewers.** No other reviewer file exists for this round. codex-1 is the implementer, and its
   material is assessed above as evidence, not as a verdict.
