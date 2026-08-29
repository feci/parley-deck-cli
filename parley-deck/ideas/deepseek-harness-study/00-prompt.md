---
idea: deepseek-harness-study
author: user
created: 2026-08-26
participants: [claude-1, codex-1, hermes-1, kimi-1, opencode-1, zcode-1]
track: deliberation
status: round-01
---

## Problem / idea

The owner is considering **rebuilding `parley-deck-cli` as a DeepSeek Harness (`dsh`) plugin
package** instead of a standalone Go CLI. Round 1 is a **study**, not a design: each participant
independently learns how the harness actually works and reports what it would and would not buy us.

The owner's sketch, in his words (translated):

> A core harness plus a plugin that natively calls the other installed agents. The whole skill
> built into the harness plugin, using the harness's other plugins. I'm thinking I would do it so
> that I do NOT use e.g. 6 models at once, because in the end the result of the work is
> **harness + model**, so the harness matters too. So I want to keep using agentic harnesses like
> claude code, codex, hermes, opencode, kimi and zcode. Unless that would be pointless. I'd also
> look at plugins for MCP, something like ToolSearch but aware of mcpanywhere, and plugins that
> work with graphify, openviking, cognee — bundle it all as one package under the parley-deck
> skill. If it makes sense. I'd like to couple the tools I actually use more tightly with the
> harness.

## What is already established (verified before this idea opened, treat as SECONDARY)

- `dsh` = `@deepseek-ai/dsh`, npm, TypeScript/Node, pnpm workspace, **developer preview**.
  Repo: https://github.com/deepseek-ai/deepseek-harness
- "Everything is a plugin", powered by **Cordis** (https://github.com/cordiverse/cordis):
  plugins contribute **services, typed events and reversible effects** to a shared context.
  There is **no privileged core**.
- **Bundles** distribute Cordis config rows + code, declared by a `dsh` field in `package.json`
  (`dsh.profile` lists a profile's bundles; `dsh.bundle` points at a patch file).
  **Profiles** are named compositions, kept in the Harness home dir.
  Composition order: bundles (profile order) → profile `cordis.patch.yml` → home patch →
  CLI `--patch`. Inspect with `dsh --profile web --dump-config`.
- Core services and their context keys: `core/session` → `ctx.sessions`;
  `core/system-prompt` → `ctx.systemPrompt`; `core/tools` → `ctx.tools`;
  `core/agent` → `ctx.agents`; `core/agent-loop` → `ctx.agentLoop`; `llm/llm` → `ctx.llm`.
- Loop shape: a **step** = one model request + its tool calls; a **turn** = zero or more steps.
  `turn/start → agent/pre-step (waterfall) → step/start → agent/request → llm/stream →
  assistant/message → tool/call → tools/pre-execute → tools/execute → tools/post-execute →
  step/end → agent/turn-stopping (serial) → turn/end`.
- Durable session events (`turn/*`, `step/*`, `user/message`, `assistant/*`, `tool/*`) survive
  reload. `ctx.sessions.fork(source, boundary?, childSessionId?)` exists.
- Docs to read: `/docs/architecture.md`, `/docs/development.md`, `/AGENTS.md`,
  `/docs/cookbook/extension-cookbook.md`.

**Do not take the above on trust — it was gathered by @claude-1 from the repo and one search.
Verify what you rely on, and correct it where it is wrong. Saying "this bullet is wrong" is one
of the more useful things you can do in round 1.**

## What each participant must produce

A study, grounded in the actual source or docs, with **PRIMARY** (you ran/read it this session,
name the file or command) vs **SECONDARY** (someone told you) tags on every claim. Cover:

1. **How the plugin model really works.** Lifecycle, mounting, isolation, versioning, what a
   plugin can and cannot reach. What does "reversible effects" mean concretely, and does it hold
   for a plugin that shells out to a foreign process?
2. **Where parley-deck's concepts would land.** Ideas, rounds, participants, artifacts, consensus,
   signoffs, quorum, phases. Which are `ctx.sessions` events, which are tools, which are agents,
   which have no home in this model at all. Be specific: name the service and the event.
3. **The load-bearing question the owner raised.** If parley-deck becomes a `dsh` plugin, is
   calling out to claude-code / codex / hermes / kimi / opencode / zcode still worth it — or does
   it collapse into "six models on one harness" and lose the thing that made it work?
   **This deck has hard evidence on that**: read
   `parley-deck/ideas/protocol-and-skill-audit/` (round-01, round-02, review/, signoff-*.md).
   In that audit the harnesses behaved *measurably differently* — one found defects the others
   missed, one fabricated PRIMARY evidence twice, one could not survive a long session, one
   re-derived every number by hand. Argue from that record, not from priors.
4. **What we would lose.** The Go CLI is one static binary with no runtime. `dsh` is Node,
   developer-preview, and iterating fast. Name what breaks: install story, offline use, the
   durable-kill process control, the drift guards, cross-platform releases, the six release
   channels we ship on today.
5. **The tool-plugin question.** MCP/ToolSearch aware of `mcpanywhere`; plugins for `graphify`,
   `openviking`, `cognee`. Does `ctx.tools` make these genuinely better than the MCP wiring we
   already have, or is it the same thing with a different config file?
6. **Your recommendation**, in one of: REBUILD / PLUGIN-ALONGSIDE (keep the Go CLI, add a `dsh`
   bundle) / WAIT (developer preview is too young) / DON'T. With the condition that would change
   your mind.

## Constraints

- **This is a study round. Do not write code and do not modify anything outside your own file.**
- The repository is READ-ONLY to you except `round-01/<your-agent-id>.md`.
- `dsh` is NOT installed on this machine. If you install it, do so in your own temp dir, never
  into the shared tree, and say exactly what you ran.
- English only.

## Non-goals

- Choosing a migration plan. That is a later idea, and only if this one says REBUILD or
  PLUGIN-ALONGSIDE.
- Rewriting the protocol. `COOPERATION.md` is not in scope.
- Any judgement about DeepSeek the *model*. This is about the harness.
