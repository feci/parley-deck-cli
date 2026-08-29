---
agent: codex-1
idea: deepseek-harness-study
round: 3
date: 2026-08-26
responding-to: [round-02 files, owner ACP challenge]
---

## Does ACP change my position?

**No. My round-2 recommendation was my own: PLUGIN-ALONGSIDE, restricted to a DSH-native operator surface, with no DSH round driver or subagent dispatch in the first bundle.** I had already found the generic ACP provider and explicitly left real-product parity `UNVERIFIED`. The new evidence removes one bad cost argument, but it does not turn ACP into a reason to rebuild Parley on DSH. If anything, the live probes below strengthen my round-2 boundary: keep the Go CLI as the only driver of record.

I withdraw any reliance on the round-2 claim that Hermes, Kimi, OpenCode, and ZCode require four product-specific DSH providers. **SECONDARY:** that claim appears in `round-02/claude-1.md`, `round-02/hermes-1.md`, `round-02/kimi-1.md`, and is load-bearing in `round-02/zcode-1.md`. A generic host-binary ACP provider means an ACP-capable product needs a configuration row and a compatibility test, not a new provider implementation. The “frozen two-product delegation” and “four providers to write” arguments therefore do not survive for products whose installed binaries actually complete ACP initialization.

That correction does not remove the other two objections. First, the ACP child remains outside DSH's tool and prompt scopes. Second, routing ACP through DSH duplicates an ACP client Parley already owns, while dividing run ownership between two drivers.

Under §15.1, I already owned the broad round-2 claim that DSH tool restrictions do not govern a foreign child, so I do not issue a self-verdict on it here. **Fresh PRIMARY observation:** in the checked-out `dsh-v0.1.1-rc.2` source, `packages/subagent/subagent-acp/src/index.ts:141-149` sets `toolFilter: false`, the other optional start capabilities to false, and `inheritsParentContext = false`; its README lines 19-21 and 68-70 say the child has its own tools and receives no parent conversation. I continue to treat that boundary as load-bearing. A non-owner should issue the formal verdict for consensus.

## Verification of the ACP claims

### (a) The generic provider

**PRIMARY — source read this session.** I ran `git -C /private/tmp/dsh-study/deepseek-harness describe --tags --always --dirty` and got `dsh-v0.1.1-rc.2`, then read `packages/subagent/subagent-acp/README.md`, `src/index.ts`, `src/run.ts`, and `package.json` in that tree.

I owned the broad existence/configuration claim in round 2, so this is fresh observation rather than a self-verdict on that broad claim. The owner's narrower claim that this provider imposes no product-executable payload pin is **CONFIRMED, PRIMARY**:

- `README.md:25-34` and `src/index.ts:26-74` define required `command`, arbitrary `args`, optional `cwd`, permission policy, explicit `env`, and disposal graces. There is no product selector.
- `src/run.ts:210-215` passes `argv: [spec.command, ...spec.args]` directly to the subprocess seam.
- `package.json` pins `@agentclientprotocol/sdk` `0.25.1` and the DSH package closure, but declares no Claude, Codex, Hermes, Kimi, OpenCode, ZCode, or other agent executable payload. That is a protocol-library pin, not a frozen child product.

Therefore `subagent-acp` is not the two-product, private-payload mechanism described in the round-2 refutation. It can launch any executable the row names, provided that executable speaks compatible ACP. The dedicated Claude/Codex providers remain pinned-payload products, but they are not the relevant path when a host CLI serves ACP.

### (b) The four local CLIs

The compound claim “four of six already speak ACP on this machine” is **WRONG AS STATED, PRIMARY**. Parley knows four configured launch recipes, but one is stale against the installed binary, and I could live-confirm only two under this task's no-extra-write constraint.

**PRIMARY — catalog observation.** `parley agents list` printed these nominal routes: Claude `--experimental-acp`, Hermes `acp`, Kimi `acp`, and OpenCode `acp`. `internal/agents/acp_specs.go:27-48` hard-codes the same recipes; it sets Codex `ACPArgs: nil` and contains no ZCode entry. That proves Parley's catalog state, not executable compatibility.

**PRIMARY — live initialize probe.** I used a Python subprocess probe that wrote only this ACP request to each process's stdin, kept stdin open long enough to read a reply, and sent no `session/new` or `session/prompt`:

```json
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{},"clientInfo":{"name":"parley-round3-probe","version":"0"}}}
```

Results:

- **Hermes 0.20.5: CONFIRMED.** `/Users/tomasfecko/.local/bin/hermes acp` returned JSON-RPC result id 1 with `protocolVersion: 1` and `agentInfo.name: "hermes-agent"`. Separately, `hermes acp --check` printed `Hermes ACP check OK`.
- **Kimi 0.36.1: CONFIRMED.** `/Users/tomasfecko/.kimi-code/bin/kimi acp` returned JSON-RPC result id 1 with `protocolVersion: 1` and `agentInfo.name: "Kimi Code CLI"`.
- **Claude Code 2.1.246: WRONG for the claimed route.** `/Users/tomasfecko/.local/bin/claude --experimental-acp` exited 1 with `error: unknown option '--experimental-acp'`. `claude --help` advertises no ACP mode, and a read-only `strings` search of the installed binary found none of `experimental-acp`, `Agent Client Protocol`, `--acp`, `ACP server`, or `ACP mode`. The hard-coded Parley catalog row is stale or premature for this installed Claude. I did not discover another stable ACP launch route.
- **OpenCode 1.18.23: capability advertised, live initialization UNVERIFIED here.** `/opt/homebrew/bin/opencode acp --help` says `start ACP (Agent Client Protocol) server`. The initialize probe failed before reading stdin because the task sandbox denied its attempted open of `~/.local/share/opencode/log/opencode.log`. I did not redirect its data directory because that would create another file, contrary to the assigned write boundary. This is not evidence that OpenCode's ACP server is broken; it is also not a successful live verification.

The most important round-3 finding is therefore a configuration-drift defect: **the installed Parley binary confidently reports a Claude ACP recipe that the installed Claude binary rejects.** The owner is right about the generic DSH mechanism, but the stated four-of-six machine fact does not survive execution today.

## Does a dsh ACP row beat parley's own ACP client?

**No, not as Parley's coordinator. I concede the routing point.** A DSH row would put the same host executable behind a second ACP client and make DSH the process owner, but it would not add protocol coverage, tool control inside the child, or Parley semantics.

**PRIMARY — Parley source read this session.** `internal/acp/client.go`, `protocol.go`, `spawn.go`, `transport.go`, and `ringbuffer.go`, plus `internal/runner/acp.go`, already implement:

- JSON-RPC/NDJSON `initialize`, `session/new`, `session/prompt`, `session/cancel`, streamed updates, permission requests, and bounded framing;
- direct host-binary spawn and whole-process-group kill;
- persisted PID/PGID/boot/start/marker attribution through `procctl` (`internal/runner/acp.go:81-104`);
- first-output and stall watchdogs, heartbeats, retry classification, stdout/stderr capture, and hard timeouts (`internal/runner/acp.go:107-210`);
- durable events for message chunks, tool calls, plans, usage, and permission decisions (`internal/runner/acp.go:332-426`); and
- canonical artifact validation plus the artifact-wins recovery rule after a post-session ACP error (`internal/runner/acp.go:213-299`).

**PRIMARY — DSH source comparison.** `subagent-acp` adds three concrete implementation details that Parley does not currently match: a credential-scrubbed parent environment, default permission rejection, and a more graduated EOF → SIGTERM → SIGKILL disposal ladder (`README.md:31-34,58-60`; `src/run.ts:242-263`). Parley's `internal/acp/shellenv.go:78-115` begins with all of `os.Environ()`, `NoopHandler.RequestPermission` selects an allow option, and its live shutdown is EOF followed by a hard process-group kill after the timeout.

Those are useful hardening ideas, not a reason to surrender run ownership. They should be considered for direct implementation in Parley's ACP path. Against them, DSH exposes only final committed assistant text to its parent and deliberately consumes thoughts, tool calls, and plans without surfacing them (`subagent-acp/README.md:80-99`; `src/run.ts:242-250`), while Parley already persists those update classes and validates the file the participant was asked to own.

The only remaining DSH-specific addition is native registration on `ctx.subagents`, so a DSH facilitator can invoke the run as a DSH tool. That is facilitator ergonomics. It does not beat the client we own, and it recreates @zcode-1's surviving “two drivers, two truths about who started a run” problem.

## The two-of-six gap

The nominal gap is worse than two on this machine today.

**PRIMARY:** Codex has `ACPArgs: nil` in `internal/agents/acp_specs.go:32-33`; ZCode has no catalog entry; the configured Claude route fails live; Hermes and Kimi initialize successfully; OpenCode advertises ACP but was not live-confirmed in this constrained run. Thus the verified uniform path covers **two of six**, the advertised path covers a possible third, and three named participants are definitely outside the currently verified route.

Even if Claude's route is repaired and OpenCode is confirmed, four-of-six does not justify moving coordination into DSH. Parley's Go driver can already select ACP for capable agents and retain its established headless adapters for Codex, ZCode, and any temporarily incompatible ACP implementation. That gives one driver, one process/event ledger, and all-six reach. Putting the same mixed topology behind DSH would add a second owner without closing the remaining gaps.

ACP still helps Parley: it can be a per-agent transport upgrade that provides structured progress and protocol-level cancellation where available. It does not help by becoming a new coordinator boundary. A uniform four-agent path is less valuable than a single authority that reaches all six and can opportunistically use ACP for the subset that passes a real initialize/session/artifact parity test.

## Current recommendation

**PLUGIN-ALONGSIDE remains my recommendation, but ACP is not part of its justification or first deliverable. Do not rebuild the coordinator and do not route Parley rounds through DSH's ACP provider.** Keep the Go CLI as driver and use its own ACP client directly for each agent that passes parity; retain Go headless adapters for the rest.

The optional DSH bundle remains limited to the operator surface from my round-2 file: DSH-native discovery of the Parley skill, a small guarded tool facade over the Go CLI, and facilitator-only MCP visibility policy. It contains no subagent rows, no round dispatch, no quorum or artifact authority, and no replacement for `procctl` or drift guards.

Change this to **DON'T** if that operator-only bundle cannot independently demonstrate a materially smaller visible tool set and a fail-closed phase-inappropriate-call denial. Change it to **REBUILD** only after all six installed harnesses pass end-to-end parity under one driver and DSH has a stable compatibility promise plus equivalents for Parley's process attribution, supervision, artifact validation, offline operation, and release surface. The generic ACP provider removes the four-provider implementation story; it does not meet those conditions.
