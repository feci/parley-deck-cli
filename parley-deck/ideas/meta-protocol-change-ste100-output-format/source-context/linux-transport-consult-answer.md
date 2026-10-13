All verification is done. Here is the consult.

---

# Consult: Linux E2BIG on full-context argv prompt — meta-protocol-change-ste100-output-format

**Recommendation in one line:** accept the proposed shape — write the exact prepared prompt into the launch's private invocation directory with its SHA256, and when the prompt exceeds the existing 120 KiB inline threshold, substitute every argv `{prompt}` element (standalone and embedded) with a short, explicit read-and-verify instruction; leave stdin delivery and all protocol bytes untouched; add the tests listed below plus the one-line 69→70 fixture bump. This is the smallest repair that keeps full-context launch alive on Linux without changing one byte of protocol or prompt content.

## What I verified (evidence)

- **The limit is real and already binding.** `internal/protocol/defaults/COOPERATION.md` is 129,796 bytes (hunks verified: 583 − 694 + 4,601 + 1 = +4,491 over the 125,305 base; whole-file SHA256 `d304f0a6…cbf494b` matches FINAL). `prepareProtocolPrompt` (`internal/runner/protocol_context.go:111`) adds ~610 bytes of framing before the task text, so the prepared prompt is ~130,410 bytes before any task. Linux `MAX_ARG_STRLEN` is 131,072 including NUL — a plain signoff task (~1,150 bytes, `buildConsensusSignoffPrompt`) lands ~131,560 and E2BIGs. Every full-context argv launch on Linux now fails, not just large tasks.
- **Single choke point.** `buildAgentInvocation` (`internal/runner/runner.go:1292`) substitutes `{prompt}` into one argv element for all five argv adapters — agy (`--print {prompt}`), hermes (positional), kimi (`-p {prompt}`), opencode (positional), zcode (embedded `--prompt={prompt}`) — and its only two callers (`execAgentProcess` runner.go:1180, `commandForLaunch` launch.go:52) both run *after* `beginProtocolLaunch`, so the invocation directory already exists at substitution time. Custom operator `headless_args` templates inherit the same path. Stdin adapters (codex, claude, gemini) carry no `{prompt}` in argv and are unaffected.
- **The failing fixture is the hybrid case.** `writeAgentsLocalConfig` (app_test.go:1597) sets `prompt_mode = "stdin"` but the fake `zcode` inherits the built-in zcode arg template, and `buildAgentInvocation` substitutes argv **regardless of PromptMode** (PromptMode only wires stdin in addition). Hence `fork/exec … argument list too long` in both Linux logs (ci-pr-linux.log:221, ci-failed-115311350336.log:206).
- **The 69/70 failure is independent and stale.** `Check` already passes (the map has the Appendix B on-demand locator); only the hardcoded `rep.Blocks != 69` at `internal/protocolpacket/packet_test.go:107` is stale. One-line fix: 70.

## The repair, concretely

1. **Extract the existing constant.** `launch.go:366` already uses `120<<10` as the conservative per-element bound below Linux's 128 KiB. Hoist it to one named constant (e.g. `maxInlinePromptArg`) shared by the interactive guard and the new headless fallback — one number, two users.
2. **Engage on prompt size, per launch.** In `buildAgentInvocation`, when `len(prompt) > maxInlinePromptArg`: stage the exact prompt bytes to `<invocation-dir>/prepared-prompt.md` (0600; reuse the atomic staged-write in `writeHandoffPrompt`, handoff.go:151 — consider renaming it `writePrivatePrompt`) and write a plain-text sidecar `prepared-prompt.sha256` in `sha256sum` format. Signature change: pass `invocationDir string` in (no test callers exist today; only two call sites to update). **Fail closed**: if the file cannot be staged, return the error — never silently fall back to inline substitution, which would diverge platforms again.
3. **Substitute every `{prompt}`-bearing element** — standalone and embedded (zcode's equals form) — with the instruction, so all five argv adapters and custom specs get identical treatment. Proposed text (~350 bytes, vendor-neutral, no shell dependence):

   > Launch prompt delivery: the complete prompt for this launch (protocol context and task, N bytes, SHA256 `<hex>`) is the file `<absolute path>`. Read that file fully and treat it as your entire prompt. Before following it, verify the file's SHA256 equals the hash above; if you cannot read the file or the hash does not match, stop and report that failure. This notice replaces the inline prompt because the prompt exceeds the operating-system per-argument limit.

4. **Stdin stays exact.** `cmd.Stdin = strings.NewReader(prompt)` (runner.go:1214, launch.go:78) is untouched — stdin-mode configs still receive the byte-exact prepared prompt; the argv notice is then a redundant pointer plus verification aid.
5. **Provenance.** Keep the file and sidecar in the invocation directory (it already persists `stdout.log`/`stderr.log`, `interactive-prompt.md`, `handoff-prompt.md` — this is established shape). `.parley-runtime` is gitignored; the canonical content evidence remains the attestation hashes in `requested.json`, which the sidecar ties to the on-disk bytes. No parser is added (AC8-clean: plain text, nothing reads it back).

**Platform-agnostic threshold — flagged tradeoff.** At the current protocol size the fallback engages on *every* full-context argv launch on *all* platforms, including darwin where the exec would succeed. I recommend that deliberately: identical model-facing bytes and identical tests everywhere, instead of darwin-inline/Linux-file divergence. This is the honest deviation to record: darwin argv launches change behavior even though darwin alone would still exec.

## Narrower existing facility (flag) and why headless diverges

`RunInteractive` (launch.go:352–369) already writes the prepared prompt into the invocation dir and has both a `{prompt_path}` file mode and the 120 KiB guard — but for arg-mode interactive it **refuses** (`interactive-prompt-argument-too-large-use-file`). Headless cannot refuse: that would kill every Linux full-context launch, i.e. the entire deck on Linux. Automatic fallback is justified; document the divergence from the interactive stance in the same comment block.

**Native `--attach` comparison:** no argv adapter in the tree declares any native file-attach flag (grepped; none exists for zcode/kimi/hermes/opencode/agy). Adopting per-vendor flags would require new probes and per-adapter grammar. The vendor-neutral instruction needs neither. If a vendor ships one later, `{prompt_path}` (already reserved for interactive) is the natural headless opt-in — note it as future work, don't build it now.

**Honesty limit (per your framing):** file delivery does not prove the model read the file. The instruction asks for hash verification; the only real evidence of reading remains the participant's output artifact, exactly as the protocol already treats context. Record this limitation verbatim in `IMPLEMENTATION.md` (Deviations + Surprises), with the CHANGELOG 1.55.0 entry noting argv-prompt launches over 120 KiB now deliver via private file reference.

## Tests

1. `internal/runner` — oversize fallback: zcode spec + >120 KiB prompt → exactly one element references `prepared-prompt.md`; no element exceeds the threshold; file bytes equal the prompt exactly; sidecar hash matches; **fail-closed** when staging fails (read-only dir → error, no inline fallback).
2. Table test over `DefaultSpecs()`: every spec whose resolved args contain `{prompt}` falls back identically (embedded and standalone forms); stdin specs without `{prompt}` produce no file and unchanged argv.
3. Under-threshold control: byte-exact inline substitution, no file written (no regression for probes/pings).
4. Stdin exactness: hybrid shape (zcode args + `prompt_mode=stdin`) — `cmd.Stdin` yields the exact full prompt while argv carries only the short notice.
5. The Linux-failing fixture passes; optionally harden the fake CLIs with `for a in "$@"; do [ "${#a}" -le 131071 ] || exit 99; done` so any future E2BIG regression fails the fixture on every platform.
6. `packet_test.go:107`: 69 → 70.

## Alternatives rejected

Switching argv adapters to stdin (changes vendor semantics per CLI; zcode `--prompt=` doesn't read stdin — unsafe); defaulting to packet-mode context (violates the frozen full-context default); truncation/summarization (violates exact-bytes); raising the limit (kernel constant); unconditional file delivery (larger behavior change; probes and small prompts keep today's semantics).

No files were created or modified; the source tree is untouched. This is consult advice for the implementer's fix-up disposition, not a review signoff — the independent review proceeds separately.
