---
agent: zcode-1
idea: deepseek-harness-study
round: 3
date: 2026-08-26
responding-to: [round-02 files, owner ACP challenge]
---

Method note: dsh-side claims are PRIMARY from the round-1 working tree
`/tmp/dsh-study/deepseek-harness` (HEAD `b150a55`, tag `dsh-v0.1.1-rc.2`, `git status` clean —
re-confirmed this session), read fresh this round. parley-side claims are PRIMARY from this
repository (files named inline) and from commands run this session, including a live ACP
handshake probe of the four claimed CLIs (script `/tmp/acp_probe.py`, outside the repo). What
follows is a verdict on my own round-2 position where I say so (§15).

## Does ACP change my position?

**Yes — the owner's challenge breaks my round-2 decisive argument, and I am the one who has to
say so.** My round-2 file closed with: "what dsh maintains is a frozen two-product delegation,
not the six installed harnesses the owner named." That is false as a statement about dsh's
capability set. `subagent-acp` spawns whatever host binary you put in its `command` config key
and drives it over ACP; it pins no payload (verification below). The frozen two-product
delegation is `subagent-claude-code` and `subagent-codex`; the ACP provider beside them is a
different thing, and round 2 (my file included) argued as if the product providers were the whole
seam. They are not.

This also guts the cost side of the round-2 consensus I joined: @claude-1's "writing four of
those to that standard is not a weekend" (SECONDARY: round-02/claude-1.md line 111),
@kimi-1's "dead on arrival" (SECONDARY: round-02/kimi-1.md line 270), and @hermes-1's "trades a
live cross-platform Go adapter family for a pinned RC payload" (SECONDARY:
round-02/hermes-1.md line 14) are all wrong **as applied to the ACP path**: there are no four
TypeScript providers to write and no RC payload to inherit — there is one plugin to install and
one config row per ACP-capable agent. I concede this plainly.

What survives of my round-2 position, re-verified this round:

1. **The substrate risk is untouched.** Pre-release stance, `SESSION_FORMAT_VERSION = 0`, no
   compatibility promise, single rc tag (PRIMARY, re-confirmed: root `AGENTS.md`; `git tag`).
   The ACP provider itself documents staying pinned on `@agentclientprotocol/sdk` 0.25.1 while
   0.28 exists, upgraded only "without improving this backend" as a separate change (PRIMARY:
   `.agents/notes/implemented/feature/2026-06-22-acp-subagent-backend.md`, "Why stay on SDK
   0.25.1?") — protocol-churn management is a live activity there, not a hypothetical.
2. **"Two drivers, two truths about who started a run" is untouched.** A round dispatched through
   `ctx.subagents` is a process parley's `procctl` never captured; restart-safe attribution and
   group-kill exist only for runs the Go driver started (PRIMARY:
   `internal/runner/acp.go` lines 81–105 — and note the ACP runner registers even its ACP
   children with `procctl.CaptureByPID`, so parley already extends durable kill across this exact
   protocol boundary).
3. **The run-log loss argument survives ACP specifically — I checked, as instructed.**
   `subagent-acp` collects only committed `agent_message_chunk` text; "the automation server
   keeps reasoning, tool activity, plans, and other trace data in the child session log rather
   than emitting them on ACP", stderr is "inherited to the parent's own stream" rather than
   captured per run, and there is no wall-clock timeout — teardown is the dispose ladder
   (PRIMARY: `packages/subagent/subagent-acp/README.md` lines 60, 96–100; feature note
   2026-06-22, "Minimal client stub" and "Consequences"). And @claude-1's round-2 constraint —
   dsh tool scoping cannot reach a foreign child — is *confirmed against ACP in dsh's own words*:
   the provider "cannot enforce the remote child's depth, tool filter, persona, or
   structured-output runtime", so it advertises no capabilities and the service rejects requests
   needing them (PRIMARY: README lines 21, 70, 98).
4. **The coverage gap is real and worse than the brief states** — three of six, not four, on the
   live binaries (verification below).

So the shape of my position changes: round 2 leaned on "dsh cannot cheaply reach the installed
agents." That leg is gone. What carries the question now is entirely §3 below — what the dsh ACP
row adds over the ACP client parley already owns — and there the evidence is concrete and points
the other way.

## Verification of the ACP claims

**(a) `subagent-acp` spawns any host binary you name: TRUE.** PRIMARY, all read this session in
`/tmp/dsh-study/deepseek-harness`:

- `packages/subagent/subagent-acp/README.md` — the Configuration table matches the brief exactly:
  `command` **required** ("Executable spawned for each run"), `args` `[]`, `env` `{}` over a
  credential-scrubbed parent environment, `permission` default `reject`. Line 5 is verbatim what
  the brief quotes: "The ACP provider runs each subagent in a fresh subprocess and drives it as an
  Agent Client Protocol client… the child has its own runtime, session, model configuration, and
  tools." Line 21: `inheritsParentContext: false`, no conversation context crosses the boundary.
- `packages/subagent/subagent-acp/package.json` — the only non-workspace runtime dependencies are
  `@agentclientprotocol/sdk` `0.25.1` and `@deepseek-ai/schemastery`. No product CLI, no platform
  payload. Contrast the product providers' pinned `@openai/codex@0.147.0` / Claude Agent SDK
  0.3.220 closures (SECONDARY: round-2 record, my own file; the exclusion note below).
- `packages/subagent/subagent-acp/src/run.ts` line 210 — `argv: [spec.command, ...spec.args]`,
  handed straight to the `dsh-subprocess` seam. There is no vendor check, no payload resolution,
  no host-binary refusal anywhere in the package.
- One precision the brief omitted: `subagent-acp` is **not in the production CLI closure either**
  — `apps/cli/package.json` depends on `dsh-subagent`, `dsh-tool-subagent`,
  `dsh-tool-subagent-control`, but not on `subagent-acp` (PRIMARY: grep of `apps/cli/package.json`).
  It is a per-profile add, like the product providers. But the 2026-08-12 exclusion note is
  explicitly about *product payload weight* and names only the two product providers (PRIMARY:
  `.agents/notes/implemented/simplification/2026-08-12-production-dsh-excludes-product-subagent-providers.md`);
  `subagent-acp` is weight-free and documented in the config catalog (PRIMARY:
  `docs/config-catalog.md` lines 2168–2219), so adding it is one package plus config rows.

**(b) Four of six already speak ACP: FALSE as stated — it is three of six on the installed
binaries, and the discrepancy is this round's most useful finding.**

- `internal/agents/acp_specs.go` matches the brief exactly (PRIMARY, read in full): claude
  `["--experimental-acp"]`, hermes/kimi/opencode `["acp"]`, codex `ACPArgs: nil` with the
  no-stable-args note, zcode absent. The file's own header says the catalog "mirrors AionUi's
  ACP_BACKENDS_ALL table". `parley agents list` (run this session, this repo's binary) prints the
  four `acp:` rows. So the *catalog* claim is exactly as the brief says.
- **Live probe** (PRIMARY: `/tmp/acp_probe.py` — spawn each CLI with its catalog args under
  `/tmp`, send one ACP `initialize` NDJSON request with `protocolVersion: 1`, read one response
  line, kill; no prompt, no model call, no repo writes):
  - `hermes acp` → **ACP-OK**: full `initialize` result, `agentInfo` "hermes-agent 0.20.5",
    session capabilities including fork/list/resume.
  - `kimi acp` → **ACP-OK**: `protocolVersion: 1`, prompt capabilities, MCP http/sse.
  - `opencode acp` → **ACP-OK**: `protocolVersion: 1`, session fork/list/resume.
  - `claude --experimental-acp` → **FAIL**: exit 1 in 0.1 s, stderr `error: unknown option
    '--experimental-acp'`. Re-probed with stderr and exit code to be sure. `claude --version` is
    2.1.246; its `--help` contains no ACP or server flag of any kind (only `--print`,
    `--output-format stream-json`, background-agent flags). `~/.parley/agents.toml` sets no
    `acp_args` override for claude (only the headless model workaround), so the row comes purely
    from the static catalog. Whether the flag existed in some other Claude Code version I did not
    establish; on the installed one it does not. Note the consistency check: dsh's own
    claude-code provider routes through the Agent SDK, not a CLI ACP mode (SECONDARY: round-2
    record) — nobody's shipping path currently relies on a claude CLI ACP flag.
  - `codex --help` → `mcp-server`, `app-server [experimental]`, no `acp` (PRIMARY, this session):
    the `ACPArgs: nil` caution is confirmed live.
  - `zcode --help` → `app-server` (ZCode Protocol), no `acp` (PRIMARY, this session): absence
    confirmed live.

The claude row is a latent defect, not an active one — every current roster agent launches
headless (`agents list` LAUNCH column), and parley fails loudly rather than silently if anyone
sets `launch_mode=acp` for claude (`internal/runner/acp.go` line 29 refuses empty ACPArgs, and
the binary itself exits 1 immediately). But it matters twice for this round. First, the brief's
premise "four of our six already speak ACP" is one agent too generous on the machine the deck
actually runs on. Second, and more important: **parley's own uniform-protocol table drifted from
a live binary within one product update.** A static catalog mirrored from a third-party table
(AionUi's) asserted a launch flag that the installed CLI rejects. Uniform-protocol claims rot at
version granularity exactly like the adapter knowledge round 2 was dismissed for — the protocol
does not remove the drift, it only moves where you discover it.

## Does a dsh ACP row beat parley's own ACP client?

**No — and I can be concrete, because parley's client is not a plan, it is running code with the
seam-by-seam comparison all PRIMARY this session** (`internal/acp/{client,protocol,spawn,transport,ringbuffer,shellenv}.go`
+ tests; `internal/runner/acp.go` read in full).

Both sides drive the same wire: `initialize` → `session/new` → `prompt` → `cancel`, with
`session/update` notifications and `session/request_permission` callbacks. `subagent-acp` adds a
stop-reason mapping table, an auto-permission policy, a credential-scrubbed child env, and a
dispose ladder. parley's runner has the stop reason (`PromptResult.StopReason`), auto-allows
permissions via `NoopHandler` **and persists every permission decision as an `agent.acp.permission`
event**, and takes a different but equally defensible env stance (allow-list PATH/TLS inheritance
from the login shell, `internal/acp/shellenv.go`). On the wire itself there is nothing to inherit
that we lack.

What the dsh row **subtracts**, relative to what parley's ACP runner already does per run:

| Concern | `subagent-acp` | parley `internal/runner/acp.go` |
|---|---|---|
| Result surface | final `agent_message_chunk` text only (README line 99) | streamed chunks aggregated, flushed every 32 chunks + final as durable `agent.acp.message_chunk` events; tool_call, plan, and usage updates each become events |
| Wire/stderr evidence | stderr inherited to parent's own stream; chunks dropped (README line 60, note 2026-06-22) | full bidirectional wire tee into `stdout.log` (requests prefixed `--> `), stderr ring buffer (8 KiB, `spawn.go`) flushed to `stderr.log` |
| Process control | in-process dispose ladder; host death externalized (round-1 finding) | `procctl.CaptureByPID` on the ACP child — durable identity in `agent.started`, restart-safe re-attribution and group-kill |
| Supervision | no wall-clock timeout | first-output watchdog, stall guard, heartbeats (consensus D1), retry-once on `no_first_output` |
| Outcome reconciliation | stop-reason mapping | D7 artifact-wins table (validated artifact overrides a post-session ACP error), D9 review snapshot publish, phase-aware artifact validation |

That table is the whole answer. The one thing the dsh row offers that parley's client does not is
**not ours to maintain** — someone else keeps the client code working. That is a real gift to a
project that has no ACP client. parley has one, with tests, integrated into the same supervision,
kill, and artifact machinery as the headless path, selectable per agent via `launch_mode=acp`,
with a 16-entry catalog for future ACP CLIs (`acp_specs.go`). Routing parley's participants
through dsh's row would trade a strictly richer client for a strictly poorer one, and would
re-introduce the two-drivers problem for good measure.

The honest concession, since the brief demands one: if this question had been asked in round 1,
before I read `internal/runner/acp.go` closely, the dsh row would have looked like the cheap way
to get ACP. It is a cheap way to get ACP. It is not a way to get what parley already has on top
of ACP, and the audit record says the part on top — retained evidence through mid-round death,
durable kill, watchdogs — is the load-bearing part. The deck's own recent history is the example:
this idea's opencode-1 failures and the audit's hermes-1 losses were diagnosable *because* the Go
driver captures participant output (SECONDARY: round-2 record, my file §"the exact shape omits
what the shape drops"; the log lengths are quoted there). A final-text-only ACP child that dies
mid-round leaves nothing.

## The two-of-six gap

Revised to **three-of-six live**: hermes, kimi, opencode answer ACP `initialize` on the installed
binaries; claude's row fails at spawn; codex has no ACP mode (its experimental `app-server` and
`mcp-server` are different protocols); zcode has none (ZCode Protocol `app-server` is its own
successor route). All PRIMARY this session (probe + `--help` sweeps).

Against the brief's premise that "parley's Go adapters reach all six today": confirmed for the
headless path — today's `parley agents list` prints a working headless argv for all six installed
agents, including zcode (the 1.44.0-era "manual-bash participant" limitation documented in
`~/.parley/agents.toml` is resolved in the current tree; SECONDARY for the history, PRIMARY for
today's output). And parley reaches the *same three* over ACP that dsh's row would. So the real
comparison is:

- **Go driver**: headless 6/6 + ACP 3/6 (same protocol, same three) + catalog rows for any future
  ACP CLI.
- **dsh-coordinated deck**: ACP 3/6 +, for claude and codex only, the pinned-payload product
  providers round 2 unanimously rejected + nothing for zcode — plus, per the claude drift above,
  the same version-granularity maintenance burden, now in YAML rows instead of Go adapters.

A coordinator that reaches half the roster over the uniform path, two more only through frozen
private payloads, and the sixth not at all, is not a uniform coordinator — it is a worse copy of
what the Go driver already ships, with the run-log loss on top. The gap does not merely fail to
help; it forces the dsh topology back onto exactly the non-uniform mechanisms (headless adapters
for zcode, pinned providers or nothing for claude/codex) whose existence was the argument for
keeping the Go layer.

## Current recommendation

**DON'T — held, with the argument honestly rebuilt.** The owner's challenge destroyed my round-2
decisive line, and the round-2 cost consensus with it; I record that above and I do not defend
the wreckage. The position now rests on three legs, none of which is the frozen-payload argument:

1. **§3**: dsh's ACP row is a strict downgrade to an ACP client parley already owns — final text
   versus retained evidence, in-process disposal versus durable kill, no timeout versus
   watchdogs. Inheriting it buys exactly one thing (someone else's client code) at the cost of
   the properties the audit record proved load-bearing.
2. **§2/§4**: the uniform path reaches three of six live, and the stale claude row just
   demonstrated inside parley itself that protocol uniformity does not survive contact with
   version drift — an adapter layer maintained against live binaries is the absorber. A dsh
   config row is the same kind of knowledge with fewer places to test it.
3. The substrate is unchanged: rc-only, no compatibility promise, and the ACP provider's own
   notes show protocol-version churn being actively managed.

What the challenge *did* change, and the owner deserves this stated plainly: if parley ever wants
an ACP coordinator inherited rather than owned, `subagent-acp` is the right shape to inherit —
host binaries, no payload pinning, honest capability reporting, fail-closed permissions. My
change-condition is therefore closer than round 2 argued, but its remaining legs are still open:
a tagged stable release with a compatibility promise covering the subagent contract; retention
of run evidence (or an explicit acceptance of the loss); and a live-coverage story for the
non-ACP half of the roster. Until those land, the Go driver is the only coordinator on the table
that reaches all six and keeps the evidence.

One actionable finding for the repo itself, found by this round and independent of dsh: the
claude ACP catalog row is stale against the installed 2.1.246 and would fail any
`launch_mode=acp` round at spawn. Either drop the row until a shipping Claude Code carries the
flag, or make `parley agents verify` probe ACP args against the live binary so drift is caught in
preflight rather than mid-round. That is a small, real defect fix in parley-deck-cli — and a
better use of the ACP insight than a plugin package.
