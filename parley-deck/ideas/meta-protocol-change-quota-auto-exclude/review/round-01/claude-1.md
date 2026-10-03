---
agent: claude-1
idea: meta-protocol-change-quota-auto-exclude
review-round: 1
date: 2026-10-04
reviewed-commit: 73ab364
---

> **EARLY, FOCUSED source-provenance review.** This is not the full-scope Phase-6 review. It covers only
> the zcode `--prompt` text-path recognizer question (FINAL §4.3 and §13.8, AC2, AC3's adversarial negatives,
> AC4 and AC18) and the frozen prototype's reset parsing that AC2 depends on. All other ACs are
> **unreviewed**. A full implementation review is still required. No consensus signoff is given here.
>
> Protocol attestation: `context_mode=full`,
> `source_sha256=8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`,
> `packet_sha256=8e9213bd45059069d484bd10e5ca1a1c509297039dfd8fc67d5e9ebda7590416`, `fallback_reason=null`.
> I read the body (`.parley-runtime/protocol-packets/full-phase6-deliberation-8e92…0416.md`, 1484 lines) in full.
> Its `shasum -a 256` matches the attested hash.
>
> Inputs: `quota-prototype.go` (frozen copy under `.parley-runtime/quota-implementation/provenance-review/`),
> and `source.json`. `shasum -a 256 /opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs` gives
> `3e3433d9…b685f`, which matches `source.json`. The package is `zcode-app-cli` 3.7.7-13; per `vendor/extraction.json`
> it is appVersion 3.7.7 and cliVersion 0.16.3. Byte offsets below refer to that file, located with `grep -bo` and
> a Node substring scanner.

## Summary

From installed zcode-app-cli 3.7.7-13 source, a fail-closed native terminal provider-error recognizer
**cannot** be established for the configured `zcode --prompt=<text> --mode yolo --cwd <root>` text path. The
path's only terminal error is one generic line, `Error: Turn execution failed`. It has no provider attribution,
and FINAL §4.4 and §4.5 need exhaustion semantics that this line does not carry. The provider's message (the
"Weekly/Monthly Limit Exhausted" text and `reset_at`) reaches the process only through the AI SDK's default
`console.error` side effect. That output is not the terminal error. It is not bound to the failing turn, it has
no framing, and other error sources in the same process write to it too. So the prototype's diagnostic-only
classification of zcode is **CONFIRMED** for this path, and it is not unnecessarily unsupported. AC2 cannot be met
on this path without an explicit scope decision or an adapter transport change. Separately, the prototype's own
reset parser would gate the recorded AC2 message even if a native recognizer existed (MAJOR-2).

## Refutation attempts

Every verdict below is mine, as a non-owner of the prototype and of FINAL's claims. Each carries its own provenance tag.

**R1. Which terminal error does the configured path emit?** — PRIMARY, located source.
- Configured invocation: `internal/agents/discover.go:412`, `HeadlessArgs: {"--prompt={prompt}", "--mode", "yolo", "--cwd", "{root}"}`.
  It passes no `--json` and no `--verbose`.
- The CLI's `runPrompt` is at offset ≈13085993 (`s(…,"runPrompt")`). On success in text mode it prints only
  `e.stdout.write(\`${W.response}\n\`)` and returns 0. On any failure it does exactly this:
  `catch(k){let T=k instanceof Error?k.message:String(k);return e.stderr.write(\`Error: ${T}${m?\` (traceId: ${m})\`:""}\n\`),o.verbose&&(…\`Cause: ${k.cause}\`…stack…),1}`.
- Turn failures are normalized at offset ≈11651409 and ≈11652317: `catch($){let D=CL($,f,"Turn execution failed")…throw … D}`.
  `CL` is at ≈11352536. It returns the error unchanged only for a core `ModelContextExceeded` error, or a core
  `ModelError` (it also handles cancellation and context overflow). Every other error becomes
  `Ae(ye.UnknownError,r,{cause:…,context:XI(void 0,nv.Wrapper)})`, where `r="Turn execution failed"`.
- The provider path does not produce a core `ModelError`. `b7r` at ≈7177124 builds the terminal error as
  `terminalError:new OP(rw(t,a,o,e.attempt,l.context))`. `OP` is a plain `Error` subclass named
  `TerminalStreamChunkError` (≈7154252). For a 429, `Yb` (≈4973045) sets
  `message:"Provider rate limited the model request."`. `CL` therefore wraps it.
- The recorded probe agrees: its final line is `Error: Turn execution failed`, with exit 1 (EVIDENCE §2).
- Verdict: the terminal error on this path is the generic wrapper. Without `--verbose` it carries no cause.
  With `--verbose` it would add only `Cause: TerminalStreamChunkError: Provider rate limited the model request.`,
  which is zcode's own text, not the upstream message. **CONFIRMED**: the terminal channel lacks the distinction
  FINAL §4.3 to §4.5 requires.

**R2. Where does the provider text come from?** — PRIMARY, located source.
- streamText (`SUr`, ≈4900685) defaults to `onError:A=s(({error:q})=>{console.error(q)},"onError")`.
- It is invoked for every stream part of type `error` (≈4930869): `at.type==="error"&&await U({error:rUr(at.error)})`.
- zcode calls `e.runtime.streamText(O)` (≈7169743), with `O=w9r(…)` (`createStreamTextOptions`, ≈7119632). Its
  `R9r({...})` object passes no `onError` and sets `maxRetries:0`, so the default `console.error` stays active.
- zcode retries itself: `b7r` computes a delay with `mW(...)` and the run continues. That is why the probe saw the
  same body several times.
- Verdict: the `statusCode: 429 / 'retry-after' / responseBody` lines in EVIDENCE §2 are Node `util.inspect` output
  of the SDK's APICallError. They are printed once per failed attempt as an SDK side-effect diagnostic. They are not
  the terminal error. **CONFIRMED**.

**R3. Refuting "unsupported is unnecessary".** I tried to build a recognizer that keys on the APICallError dump:
header `APICallError [AI_APICallError]:`, `statusCode: 429`, a parseable `responseBody`, and the
`Symbol(vercel.ai.error.AI_APICallError)` marker, followed by `Error: Turn execution failed (traceId: …)` and exit 1.
Three independent, source-grounded reasons defeat it.
- (a) **Other sources share the sink.** `runToolsTransformation` (`gHo`) sends tool-execution promise rejections
  and its outer handler to the same error part, and so to `console.error`. At ≈4899238 and ≈4899423:
  `.catch(F=>{h.enqueue({type:"error",error:F})})` and `catch(O){h.enqueue({type:"error",error:O})}`.
  Internal parts (`text part … not found`, ≈4931114) go there too. This confirms the prototype's stated limitation.
  It also bears directly on AC18: tool-originated errors are printed through the same channel.
- (b) **The dump is not bound to the failing turn.** Nested model calls in the same process (background or detached
  agents) use the same adapter and streamText default. See the
  `console.error("Detached background agent lifecycle failed",e…)` sites at ≈11193862 and ≈11200284. The dump does
  not carry zcode's traceId, while the terminal line carries only the traceId. A 429 dump from a nested call
  (possibly another provider or model) followed by an unrelated wrapped failure produces the same stderr shape. The
  text channel cannot show that the dump *is* the terminal error, as FINAL §4.3 requires ("captured the terminal error").
- (c) **The output has no framing.** Executed check (PRIMARY), offline and deterministic, with no provider and no zcode
  invocation: `/tmp/zc-framing/t2.js` under `node v26.10.0`. It builds an APICallError-shaped object with the recorded
  body, then a different `Error` whose *message* embeds that text. Output:
  `every genuine line appears verbatim as a whole line in forged output: true`.
  `util.inspect` does not escape newlines inside messages. Any error message that contains this text therefore
  yields stderr lines identical to a real dump, including the header and the `Symbol(...)` marker lines. A strict
  line-shape recognizer cannot tell origin. Limitation: zcode runs on electron-node, and I tested system Node's
  `util.inspect`. Frames were filtered out because they differ per build.
- Result: the refutation of "unsupported is unnecessary" fails. A strict-shape stderr recognizer would rest on
  correlation, not captured terminal provenance.

**R4. Refuting "a false-positive recognizer is needed for safety".** — PRIMARY, located source.
- I checked whether tool processes can write into zcode's stderr directly. The MCP SDK transport defaults to
  `stdio:["pipe","pipe",this._serverParams.stderr??"inherit"]` (≈8346016). zcode overrides it with
  `stderr:"pipe"` in `createTransport` (≈8374040), so the MCP-stderr injection path I hypothesized is **refuted**.
- I did **not** verify the Bash-tool child stdio, or the WASI `fd_write` shim that writes to `process.stderr`
  (≈1893466). These stay UNVERIFIED. They are not needed for the verdict, because R3(a) to R3(c) already suffice.
- Assistant-quoted text never reaches stdout on failure: `runPrompt` prints `W.response` only on success.

**R5. Is there an alternative structured channel?** — PRIMARY, located source; scope noted.
- zcode does build a structured error payload. `EA` / `Wxi` at ≈8518953 derive
  `attribution:{source:"provider"|"runtime"|"tool"|"network", providerId, modelId, statusCode, providerErrorCode, retryable}`.
  The chain walk is `Zcn`, which follows `cause ?? lastError ?? error` up to 12 levels.
- `rw` (≈7152369) sets `source` with `aYo`, which returns provider when `statusCode` is defined. `PL(...)` emits the
  failure event with `logEvent:"turn.failed"`.
- `runPrompt` never prints any of this in text mode, and its `--json` branch writes only on success.
- Whether that payload reaches a consumable channel (app-server or ZCode Protocol session events, or a persisted
  session log) is **not established** here. I also have not checked whether its `message`/`detail` preserve the
  upstream "Limit Exhausted… reset_at" text. Using it would be a transport change, not the configured path.

**R6. AC2's recorded positive.** — PRIMARY, located evidence plus executed check.
- EVIDENCE §2 omits headers, session ids and correlation ids "on purpose". The header line, `url`,
  `requestBodyValues`, Symbol markers and the `(traceId: …)` suffix are also absent.
- `grep -rlF 'Weekly/Monthly Limit Exhausted'` over both `parley-deck/runs` trees, `~/.parley`, `~/.zcode/cli/logs`
  and `~/.zcode/logs` found no raw capture. No credentials were read.
- A byte-faithful native positive fixture therefore cannot be reconstructed from recorded evidence. Any envelope
  would be synthesized.
- See MAJOR-2 for the prototype predicate's verdict on the recorded message itself.

**R7. AC4 and AC18 against the prototype.** — PRIMARY, located source (`quota-prototype.go`).
- `ClassifyQuota` (`:43-58`) never sets `Eligible`.
- `refineNativeQuota` is private and has no production caller (`:60-62` comment; `Collector.QuotaEvidence`, `:205-213`,
  calls only `ClassifyQuota`).
- On the frozen copy, zcode therefore cannot auto-exclude, so AC4 and AC18 hold for the zcode path as written. This
  covers the prototype only. The live tree was not reviewed.

## Findings

### [MAJOR] MAJOR-1: AC2 is unattainable on the configured zcode text path; a recorded scope decision is needed
- **What is wrong.** FINAL AC2 needs "a provenance-verified zcode recognizer". Per R1 to R3, the configured
  `--prompt` text path cannot provide one. The terminal error is the generic wrapper `Turn execution failed`.
  The provider text is a non-terminal SDK diagnostic that is unbound and unframed, and other error sources share it.
- **Why it blocks.** The prototype keeps zcode diagnostic-only, which is correct under §13.8 ("An adapter that cannot
  be verified ships unsupported"). But then AC2 cannot pass, and AC19 and AC20 completion inherit that gap. If no
  decision is recorded, the gap is silently absorbed, which §4 Phase 5 forbids.
- **Suggested fix.** Pick one, record it in `IMPLEMENTATION.md` `## Deviations from FINAL.md`, and get an owner
  decision, since this is protocol work:
  1. Ship zcode as diagnostic-only. Mark AC2 as not met on this adapter, with this provenance. Move the positive path
     to a follow-up idea.
  2. Change the zcode adapter transport to a channel whose *terminal* failure carries attribution. R5 is the
     candidate. That needs its own located behavior, plus recorded positive and adversarial negative fixtures
     before any eligibility.
- **Rejected option.** A strict stderr dump recognizer cannot be made fail-closed (R3).

### [MAJOR] MAJOR-2: the prototype's reset parser gates AC2's recorded message, so it cannot reach eligibility
- **Where.** `quota-prototype.go:74`, `resetAtText = (?i)\breset(?:s| will reset)? at ([^);\n]+)`, together with
  `:141` (an RFC3339Nano-only parse) and `:172-175` (any failed parse returns `invalid()`).
- **Executed check** (PRIMARY): `/tmp/zc-reset/main.go`, with the regexes copied verbatim, run on the recorded
  message. Output:
  - `resetAfter capture: "49h 8m 50s"`
  - `resetAtText capture: "2026-10-05 06:14:57 (reset after 49h 8m 50s" -> RFC3339Nano err: … cannot parse " 06:14:57 (reset after 49h 8m 50s" as "T"`
- **Result.** The capture is unparseable, so `quotaReset` returns `ok=false` and `refineNativeQuota` ends with
  "contradictory, past or unparseable stated reset". That is AC2's exact positive (49 h away) classified as a gate.
  The three reset expressions actually agree: `reset_at` 2026-10-04T22:14:57Z, the "reset after" duration from the
  observation, and the human time read as UTC+8.
- **Two compounding defects.** First, `[^);\n]` stops at `)` but not `(`, so the capture runs into the parenthetical.
  Second, a timezone-less human-readable time is treated as a stated, invalid reset.
- **Suggested fix.** Keep the no-local-timezone rule. Do not let a timezone-less human clock time void a message whose
  machine-readable resets (`reset_at` and/or "reset after") agree. For example, ignore absolute times that carry no
  zone as non-decisive and record them raw, while still gating explicit contradictions.
- **Fixtures to add:**
  - Positive: the recorded body above, including `reset_at`, with observation 2026-10-02T21:06:07Z.
  - Negatives: the same body with `reset_at` contradicting "reset after" by more than 1 s; and with "reset after 30m".
  - Negative: a timezone-less time as the *only* reset. It should stay gated.
- **Scope.** This affects every future native recognizer, not only zcode.

### [MINOR] MINOR-1: the support-table limitation for zcode omits the decisive reason
- **Where.** `quota-prototype.go:20`.
- **What is wrong.** "SDK error stream also accepts tool callback failures; recorded excerpt omits its native envelope"
  is accurate (R3a, R6) but not decisive. The decisive facts are R1 (the generic terminal wrapper with no attribution)
  and R3b and R3c (no binding to the turn, no framing).
- **Suggested fix.** Record those with the offsets above in `testdata/quota/README.md`, so the diagnostic-only status
  is not re-litigated on the weaker ground.

## Open questions

1. **Owner or organizer.** Which MAJOR-1 option governs stage 1? If option 2, is a zcode transport change in scope
   under FINAL, or a deviation needing owner confirmation?
2. **Source checks for option 2.** Does zcode's structured failure payload (R5) reach a consumable channel in
   `zcode app-server` or in the persisted session events? Does its `message`/`detail` keep the upstream body text?
   This needs located source plus a recorded native capture. I cannot settle it here.
3. **Unverified side channels.** Bash-tool child stdio and the WASI `fd_write` shim (≈1893466) remain unverified
   possible writers to zcode's stderr. They only matter if a stderr-based design is ever reconsidered.
