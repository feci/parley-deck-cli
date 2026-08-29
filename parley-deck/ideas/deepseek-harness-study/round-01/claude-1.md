---
agent: claude-1
idea: deepseek-harness-study
round: 1
date: 2026-08-26
---

## Summary

**The single most important finding: the harness already ships the thing the owner wants to build.**
`packages/subagent/` contains `subagent-claude-code` and `subagent-codex` — providers that start a
**real Claude Code child through the official Claude Agent SDK** and a **real Codex `app-server`
child** — plus `subagent-acp`, `subagent-dsh-sdk`, and in-process spawn/fork drivers, all
registering on `ctx.subagents`, with `tool-subagent` exposing delegation to the model (PRIMARY:
`packages/subagent/README.md` in a clone at `/tmp/.../dsh-src`).

So "a plugin that natively calls the other installed agents" is not something we would invent. It
is an existing capability family with two of our six agents already implemented. **What is missing
is exactly four providers: hermes, kimi, opencode, zcode** — PRIMARY: `grep -ril
'hermes\|kimi\|opencode\|zcode' packages/` returns only unrelated hits (`llm-pi-ai` tests,
`bash-local`, `tool-web/search.ts`).

That reframes the question. It is not "should we rebuild parley-deck on dsh" but **"is parley-deck
a `dsh` bundle that adds four subagent providers and a protocol, or is it a CLI?"**

My recommendation is **PLUGIN-ALONGSIDE**, and the blocking reason is version `0.1.1-rc.2`.

## 1. How the plugin model really works

PRIMARY (`docs/architecture.md`, cloned repo):

- **Cordis** is the substrate: plugins contribute **services, typed events and reversible effects**
  to a shared context. There is no privileged core — the model adapter, tool registry, session log
  and *the agent loop itself* are all plugins.
- **Bundles** distribute Cordis config rows plus code, declared by a `dsh` field in `package.json`
  (`dsh.profile` lists a profile's bundles; `dsh.bundle` points at a patch file). **Profiles** are
  named compositions in the Harness home dir.
- Composition order is layered and inspectable: bundles (in profile order) → profile
  `cordis.patch.yml` → home patch → CLI `--patch`. `dsh --profile web --dump-config` prints the tree.
- Services and their context keys: `ctx.sessions`, `ctx.systemPrompt`, `ctx.tools`, `ctx.agents`,
  `ctx.agentLoop`, `ctx.llm`, and (PRIMARY, from package READMEs) `ctx.subagents`, `ctx.skills`.
- The loop is turn/step with **waterfall** extension points that require `next()` delegation:
  `turn/start → agent/pre-step → step/start → agent/request → llm/stream → assistant/message →
  tool/call → tools/pre-execute → tools/execute → tools/post-execute → step/end →
  agent/turn-stopping → turn/end`.

**On "reversible effects" and foreign processes** — the question the brief asked me to press.
I did **not** find a general transactional rollback for side effects, and I would not claim one.
What I did find is that the *process* side is taken seriously: `subagent-claude-code` publishes a
run only after the SDK's `spawnClaudeCodeProcess` hook supplies **a live CLI handle owned by
`dsh-subprocess`**, and `dispose()` "aborts the run, asks the SDK query to close, invokes the
shared process-tree termination escalation, and waits for whole-tree exit" (PRIMARY:
`packages/subagent/subagent-claude-code/README.md`). Reversibility there means *ownership and
teardown*, not undo. A plugin that writes a file still wrote the file.

**Scale, for honesty about what "a plugin" means here:** 50 package directories under
`packages/`, a pnpm workspace with `vendor/*`, `apps/*`, `native/landlock-run`, and a `python/`
tree (PRIMARY: `ls`, root `package.json` workspaces). This is not a small library.

## 2. Where parley-deck's concepts would land

| Parley concept | dsh home | Confidence |
|---|---|---|
| Participant (claude-1, codex-1, …) | `ctx.subagents` provider, one per agent | PRIMARY — two exist |
| Round dispatch | `ctx.subagents` delegation; `tool-subagent` if the model should drive it | PRIMARY |
| Round-1 independence | **already the default**: `subagent-claude-code` reports `inheritsParentContext: false` — the child gets a standalone text task and the parent cwd, *not* the parent conversation | PRIMARY |
| "Return only a short confirmation with the path written" | forced by the contract: the shared result boundary caps the whole result text at **4096 UTF-8 bytes** | PRIMARY |
| Idea / round / review artifacts | **stay files.** They are the audit trail; `ctx.sessions` is an append-only event log, which is a different thing | my judgement |
| Phase transitions, quorum, signoff triage | `ctx.tools` commands + our own logic; no native home | my judgement |
| Skill (SKILL.md) | `ctx.skills` with `skill-filesystem`, which "discovers skills from local filesystems" and `tool-skill` publishes the catalog + loader | PRIMARY |
| Durable kill across restarts | `dsh-subprocess` owns process trees **within a run**; our `internal/procctl` survives a *parley restart* by persisting identity. Not obviously the same guarantee | needs verification |

The mapping is good where the harness has a concept and poor where it does not. **Consensus,
quorum and signoff triage have no home** — they would be our code either way. That is fine; it is
also most of what `internal/consensus` is.

## 3. Harness diversity: worth keeping, or does it collapse?

**Keep it. The audit is evidence, not intuition** — and it points the same way the owner's
instinct does.

From `parley-deck/ideas/protocol-and-skill-audit/` (PRIMARY, I facilitated it and re-read the
files this session): six harnesses ran the *same* verification task over the *same* 47 findings.

- @codex-1 assessed 23 and returned **6 CONFIRMED / 8 PARTIAL / 9 REFUTED**. @kimi-1 assessed 42
  and returned **36 / 6 / 0**. @zcode-1 assessed 32 and returned **30 / 2 / 0**. *The same corpus
  drew 6 confirmations from one and 36 from another.* All nine refutations in the corpus were
  codex's.
- Behavioural differences were not noise: @hermes-1 filed **PRIMARY-tagged claims that did not
  survive execution, twice**; @opencode-1 could not survive a long session; @zcode-1 re-derived
  every number by hand and caught an off-by-one that three others had signed.
- The defects that actually mattered were found by **disagreement**: codex blocked twice on things
  four others had accepted, and kimi found that a commit did not contain the fix it claimed.

That is harness+model diversity doing work. **But the honest reading is narrower than "six is
good".** What produced the value was (a) *adversarial instruction* — every reviewer was told to
re-derive rather than read, and (b) *different failure modes*. Six identical harnesses on six
models would have given (b) partially. Six models on **one** harness would have given much less of
it, because the harness is what decides whether the agent re-runs the command or trusts the diff.

**So: `subagent-claude-code` and `subagent-codex` are exactly the right shape** — they start the
*real product*, with its own settings, sandbox and account state (PRIMARY: the claude-code provider
"deliberately omits the SDK `settingSources` option" so the host's normal user/project/local Claude
settings apply). That is harness diversity preserved inside dsh, not collapsed.

**The catch:** we would have to *write* the four missing providers ourselves, and each one is
non-trivial. The claude-code provider's README is ~2 pages of failure-mode taxonomy alone.

## 4. What we would lose

Measured (PRIMARY, `wc -l` this session): **35,040 lines of Go across 26 internal packages, plus
21,213 lines of tests.** The deck holds **1,519 markdown artifacts**. A rebuild is not a port.

Concretely at risk:

- **One static binary, no runtime.** `dsh` is Node + pnpm + a `native/landlock-run` component.
  Our install story is `brew install` / a single downloaded binary for 6 platforms.
- **Six release channels** we ship on today (GitHub, Homebrew ×2, npm, winget ×2) tuned around
  Go cross-compilation.
- **The drift guards** (`internal/protocol/drift_test.go`) and the whole §15 verification
  machinery — that is protocol enforcement, not harness plumbing, and it does not come across for
  free.
- **`internal/procctl` durable kill across a parley restart.** dsh owns process trees per run;
  whether it survives a harness restart is unverified by me.
- **Maturity.** `0.1.1-rc.2` (PRIMARY: root `package.json`). Developer preview, "iterating
  rapidly". We would be rebuilding a working tool on an interface that is explicitly unstable.

## 5. The tool-plugin question

`mcp-client` "registers external MCP server tools on `ctx.tools`" (PRIMARY:
`packages/mcp/README.md`). That is the same bridge every MCP host has. **On its own it is not
better than what we have** — `mcpanywhere` already routes 882 tools today.

The genuinely new thing is that `ctx.tools` is a **scoped registry with an execution pipeline**
(`tools/pre-execute`, `tools/execute`, `tools/post-execute`), so a plugin can *filter and rewrite*
tool exposure per agent without the MCP server knowing. **That is what a ToolSearch-that-knows-
mcpanywhere actually needs**: not another bridge, but a scoping layer that decides which of 882
tools a given participant sees this round. I rate this the strongest technical argument in the
owner's whole sketch.

graphify / openviking / cognee: these are retrieval-shaped. `ctx.tools` gives them a uniform
registration and the pipeline gives per-round scoping. Worth doing — but note they are *already*
reachable as CLIs and MCP servers, so the win is ergonomics and scoping, not capability.

## 6. Recommendation

**PLUGIN-ALONGSIDE.** Keep the Go CLI as the protocol authority and the thing that ships. Build a
`dsh` bundle that (a) registers our four missing subagent providers and (b) exposes the deck as a
skill via `skill-filesystem`. Do not rebuild 35k lines onto an `0.1.x-rc` interface.

**The condition that would change my mind to REBUILD:** dsh reaching a stable minor with a
documented plugin-API compatibility policy, *and* at least two of hermes/kimi/opencode/zcode
gaining upstream providers so we are not the sole maintainer of four foreign-process integrations.

**The condition that would change my mind to WAIT (do nothing yet):** if a spike shows the
4096-byte result boundary or `inheritsParentContext: false` cannot carry a Parley round, the whole
mapping in §2 collapses and there is nothing to build alongside.

## Corrections to the brief

The brief was written by me from one repo read and one search. Two things it got thin:

1. It said "no privileged core" and listed six services. It **omitted `ctx.subagents` and
   `ctx.skills`**, which are the two most relevant families to this idea. That omission would have
   led the round toward "we must build agent-calling ourselves", which is false.
2. It framed the choice as rebuild-vs-not. The real axis is **who owns the four missing providers**.

## What I could not verify

- I did **not** install or run `dsh`. Everything above is source and docs from a
  `git clone --depth 1` at `/tmp/.../dsh-src`, not observed behaviour.
- I did not verify that `subagent-claude-code` actually works on this machine.
- I did not read `docs/development.md`, `AGENTS.md`, or the extension cookbook.
- I did not verify whether dsh's process ownership survives a harness restart — the specific
  guarantee `internal/procctl` was built for.
- The 4096-byte result boundary is from the claude-code provider README; I did not confirm it is
  the shared contract for *all* providers.
