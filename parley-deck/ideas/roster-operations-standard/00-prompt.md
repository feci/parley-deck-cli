---
idea: roster-operations-standard
author: user
created: 2026-08-06
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: round-01
track: standard
---

## Problem / idea

**The roster has no standard surface.** Ask "what is the current agent roster?" and you get a
different answer depending on which command ran and who ran it. Adding `opencode` yesterday left
sessions inconsistent — hermes works in some, not in others. The user wants roster operations
standardized across the CLI and the skill: one stable table, one way to update, one way to sync.

### Measured today (facilitator, `PRIMARY`, parley 1.39.0)

**1. Two commands, two incompatible tables, different row sets.**

`parley roster show` — 4 rows:
```
ROSTER-ID  FAMILY  DISPLAY-NAME             MODEL              EFFORT  SPEED     AUTO
claude-1   claude  claude_opus-4.8-1m_max   claude-opus-5[1m]  max     balanced  yes
```
`parley agents list` — 13 rows, none of the same columns:
```
AGENT  INSTALLED  VERSION  LAUNCH  HEADLESS  SANDBOX  APPROVAL  MODEL  TIMEOUT  HOME  AUTO  BACKEND
```
Neither is documented as *the* answer to "show me the roster".

**2. `opencode` is in `agents list` but absent from `roster show`.** It was promoted to a full
adapter in CLI 1.39.0. Becoming an adapter did not make it a roster member. Two different meanings
of "roster" coexist with no defined relationship.

**3. Three agents report a model they do not launch.** This is the same declared-vs-effective split
that §15 and CLI 1.39.0 fixed for autonomous write, now unfixed for `model`:

| agent | MODEL column | actually launched |
|---|---|---|
| claude | `claude-opus-5[1m]` | `--model claude-opus-4-8[1m]` |
| codex | `gpt-5.6-sol` | no `-m` flag at all |
| kimi | `kimi-code/k3` | no model flag at all |

Root cause for claude: `internal/agents/discover.go:219` embeds the model literal **inside**
`HeadlessArgs`, while `Model` (`:226`) is a separate field. Config layers set `Model` and never
rewrite the embedded flag, so both `~/.parley/agents.toml` and `parley-deck/agents.toml` pinning
`claude-opus-5[1m]` are silently ineffective at launch. `roster show` contradicts itself in one
row: DISPLAY-NAME says `opus-4.8`, MODEL says `opus-5`.

**4. There is no documented update or sync operation.** `parley roster init [--scope
session|machine]` exists but is not in `parley --help`. Nothing states how to change one agent's
model locally, how to change it globally, or what "sync the session roster with the global one"
means mechanically.

## What the user asked for

1. **A stable roster table with fixed columns**, the same every time, including at least:
   `agent`, `model`, `model family`, `model company`, `effort`, `speed`, `auto` — and whatever
   else the participants judge load-bearing (installed, roster membership, effective-vs-declared).
2. **A simple way to update the roster — local (session/deck) and global (machine).**
3. **A defined `sync local ← global` operation** that is always performed the same way, not
   re-invented per session.
4. **An explanation and a fix for the inconsistency** that adding `opencode` introduced.

## Constraints

- The table must be identical across `parley` output and anything the skill prints; the skill must
  not describe a second format.
- Where declared and effective values differ, the surface must show the **effective** one and flag
  the divergence — do not repeat the `AUTO` mistake for `model`.
- Column set must be stable across releases; adding a column is a deliberate, documented change.
- `model family` and `model company` must be derived, not hand-maintained per deck.
- No new configuration file unless justified; the layering
  (built-in → `~/.parley/agents.toml` → deck → `agents.local.toml` → env) already exists.
- Any CLI change ships with the skill in the same release.
- Backwards compatibility: existing decks must keep working without manual migration, or the
  migration must be automatic and stated.

## Non-goals

- Redesigning the agent-discovery mechanism itself.
- Adding new agents to the roster as part of this idea.
- Changing the protocol's quorum or signoff rules.
- A GUI or TUI redesign; this is about the command surface and the data contract behind it.
