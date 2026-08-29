---
agent: opencode-1
idea: deepseek-harness-study
round: 3
date: 2026-08-26
responding-to: [round-02 files, owner ACP challenge]
---

## Does ACP change my position?

**No.** Round-2 DON'T was mine (`round-02/opencode-1.md`). ACP destroys two supporting arguments I used; it does not destroy the recommendation.

Withdrawn, because they do not apply to `subagent-acp`:

- The four-provider weekend. I wrote "Cost is larger than today's six Go adapters" and "Do not write them" (`round-02/opencode-1.md`, mine). For an ACP-capable host CLI that is gone: a config row names `command`/`args`. I already flagged ACP as unverified then ("ACP is a maybe… I did not verify which of hermes/kimi/opencode/zcode do"). Now verified, with a correction below.
- Frozen two-product payload, **as applied to the ACP path**. That objection still holds for `subagent-claude-code` / `subagent-codex` (PRIMARY: `.agents/notes/implemented/simplification/2026-08-12-production-dsh-excludes-product-subagent-providers.md`). It does not hold for `subagent-acp`.

What is left, and is enough:

- dsh tool scoping still cannot reach a foreign child. ACP makes that *more* explicit, not less.
- We already drive ACP. A dsh row is a second, weaker client.
- Live uniform ACP coverage on this machine is **3 of 6**, not 4. Catalog ≠ handshake.
- Pre-release / no compatibility promise is untouched (PRIMARY: `/tmp/dsh-study/deepseek-harness/AGENTS.md` lines 5–7, HEAD `b150a551`, tag `dsh-v0.1.1-rc.2`).
- Two drivers, two truths about who started a run is untouched.

## Verification of the ACP claims

Clone: `/tmp/dsh-study/deepseek-harness`, HEAD `b150a551b8d465e31e418e1b2eaf5e79bbb7d28e`, tag `dsh-v0.1.1-rc.2` (PRIMARY: `git rev-parse` / `git describe` this session). I did not run `dsh`.

### (a) `subagent-acp` imposes no pinned product payload — CONFIRMED

PRIMARY, this session:

- `packages/subagent/subagent-acp/README.md` Configuration table: `command` **required** ("Executable spawned for each run"), `args` default `[]`, `env` default `{}` over a credential-scrubbed parent env. Opening sentence: "The ACP provider runs each subagent in a fresh subprocess and drives it as an Agent Client Protocol client… the child has its own runtime, session, model configuration, and tools."
- `packages/subagent/subagent-acp/src/index.ts`: `command: z.string().required()`, `args: z.array(z.string()).default([])`. `start()` copies `this.config.command` / `this.config.args` into the run spec. `inheritsParentContext = false`. `capabilities = { outputSchema: false, depthLimit: false, toolFilter: false, persona: false }`.
- `packages/subagent/subagent-acp/src/run.ts`: `ctx.subprocess.spawn({ argv: [spec.command, ...spec.args], ... })`. `newSession({ cwd: spec.cwd, mcpServers: [] })`. Collects only `agent_message_chunk` text.
- `packages/subagent/subagent-acp/package.json`: runtime deps are `@agentclientprotocol/sdk` `0.25.1` and schemastery. **No** `@anthropic-ai/claude-agent-sdk`, **no** `@openai/codex`.

Honest pin: the protocol SDK version is pinned. The **child executable is not**. That is the distinction the frozen-payload argument needed, and it holds.

`subagent-acp` is also absent from the default production closure. PRIMARY: `@deepseek-ai/dsh-base` `package.json` lists `dsh-subagent` + in-process spawn/fork, not `dsh-subagent-acp`; `packages/bundle/base/cordis.patch.yml` has no `dsh-subagent-acp` hit; published CLI tarball `/tmp/dsh-study/package/package.json` has no `dsh-subagent-acp`. Adding it is still a small package plus a config row, not a 75–111 MB private CLI.

### (b) Four of six already speak ACP — PARTIALLY REFUTED

The catalog claim is real. The live claim is not.

PRIMARY: `parley agents list` this session prints exactly the four ACP lines the owner quoted (`claude --experimental-acp`, `hermes acp`, `kimi acp`, `opencode acp`). PRIMARY: `internal/agents/acp_specs.go` lines 29–47 records the same four; `codex` is `{ACPArgs: nil}` with the "no stable ACP launch args" note; `zcode` is **not** in `ACPCatalog()`.

PRIMARY, live binaries this session (initialize JSON-RPC on stdio, protocolVersion 1):

| CLI | Declared | Help | Handshake |
|---|---|---|---|
| kimi 0.36.1 `/Users/tomasfecko/.kimi-code/bin/kimi acp` | yes | "ACP server over stdio" | `agentInfo.name=Kimi Code CLI` |
| opencode 1.18.23 `/opt/homebrew/bin/opencode acp` | yes | "start ACP server" | `agentInfo.name=OpenCode` |
| hermes 0.20.5 `/Users/tomasfecko/.local/bin/hermes acp` | yes | "Run Hermes Agent as an ACP" | `agentInfo.name=hermes-agent` (needs stdin kept open; closing stdin immediately yields empty stdout) |
| claude 2.1.246 `/Users/tomasfecko/.local/bin/claude` | **catalog yes** | `--help` contains no "acp" | **`--experimental-acp` → `error: unknown option`; `acp` → unknown command; `--acp` → unknown option** |
| codex | `ACPArgs: nil` | `mcp-server` only | no ACP launch |
| zcode 3.7.7-13 | absent from catalog | `app-server` (ZCode Protocol), no ACP | no ACP launch |

This is the most valuable finding this round: **`parley agents list` is a spec dump, not a handshake.** The owner's "four of six" is true of our catalog and false of the installed `claude`. Live uniform ACP coverage is **3 of 6**.

@claude-1's constraint checked against ACP specifically: it **survives, and the source states it as a capability hole**. PRIMARY: README "ACP advertises no start-time capabilities because this process cannot enforce the remote child's depth, **tool filter**, persona, or structured-output runtime"; `toolFilter: false`; `inheritsParentContext: false`; `mcpServers: []`. Parent `ctx.tools.restrict()` / `guard()` / `tools/pre-execute` do not enter the child. The child's tools are the child's.

## Does a dsh ACP row beat parley's own ACP client?

**No. Concede the point.** Routing through dsh adds a second client we do not control and subtracts things we already have.

PRIMARY, this repo, this session:

- `internal/acp/` — `client.go`, `protocol.go` (`ProtocolVersion = 1`, JSON-RPC 2.0 NDJSON), `spawn.go`, `transport.go`, `ringbuffer.go`, plus tests.
- `internal/runner/acp.go` `runACPAgent`: spawn **host** `agent.Path` + `ACPArgs`, `procctl.CaptureByPID`, initialize → `NewSession({CWD})` → `Prompt`, supervision (heartbeat / first-output / stall), event log (`agent.acp.initialized`, `session_opened`, `prompt_completed`, `tool_call`, `plan`, `usage`, `permission`, `message_chunk`), artifact validation with artifact-wins after session open, permission auto-allow, fs/* refused so the child uses its own tools.
- `internal/runner/runner.go` already branches `LaunchACP` into that path.

What a dsh ACP row adds over that:

- Nothing the coordinator needs. Same spawn-host-binary-and-speak-ACP shape.
- Final-text-only parent result (README Known Limitations: reasoning, tool activity, plans stay in the child). We already log those events.
- Empty `mcpServers`, no tool filter, no persona, no structured output.
- No `procctl` identity, no restart-kill attribution, no first-output watchdog.
- A Cordis/`ctx.subagents` start, so a second answer to "who started this run."
- An RC package that is not even in `@deepseek-ai/dsh-base`.

If the owner wants ACP launch, the switch is `launch_mode = "acp"` on agents we have already proven speak it. That is a config edit in a client we own.

## The two-of-six gap

Worse than the brief. Catalog gap is 2/6 (codex + zcode). **Live gap is 3/6** (those two plus claude 2.1.246).

A coordinator that can only reach three (or four on a machine whose `claude` still accepts `--experimental-acp`) over the uniform ACP path does **not** help. Parley's Go adapters already reach all six on the headless path that this idea is running on. Using dsh ACP for the speakers and Go for the rest is two drivers by construction. Using dsh ACP even for the speakers duplicates `runACPAgent`.

The gap would not flip this if it closed. Closing it would still leave a weaker client in front of a stronger one.

## Current recommendation

**DON'T.** Unchanged from my round 2. Do not rebuild. Do not add a parley `dsh` bundle. Do not register dsh ACP rows as a participant transport.

ACP is the way to talk to a foreign harness over a standard wire. We already took that way, in Go, against host binaries, with kill and logs. The owner's question refutes the "write four TypeScript providers / inherit two frozen CLIs" cost story. It does not refute "the coordinator is a process supervisor plus a file protocol," and it does not make dsh the better ACP client.

Condition that would change my mind to **PLUGIN-ALONGSIDE** (transport, not a skill-host): still the round-2 bar — tagged stable dsh with a compatibility promise, **and** continuable host-binary children for at least four of six **that beat `internal/acp` + `runACPAgent` on attribution, logs, and artifact validation**. A config row that loses to code we already ship is not that bar.

Condition that would change my mind to **WAIT**: none. There is no parley work to queue on this path.
