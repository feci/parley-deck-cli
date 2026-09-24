---
from: user
to: codex-1
idea: meta-protocol-change-designated-implementer
date: 2026-09-25
status: authorized
relayed-by: the owner's Claude Code session (not a participant artifact)
---

# Owner decision: the global default implementer is codex-1

## The owner's words (verbatim, Slovak) and translation

> "aky mas globalny roster? treba setnut organizatora na claude s opus 5.5 a implementatora na codex s astra
> gpt-6 a ucastnikov na kimi s kimi k3 a zcode s glm-5.3"

"What is your global roster? Set the organizer to Claude with Opus 5.5, the implementer to Codex with GPT-6
Astra, and the participants to Kimi with Kimi K3 and Zcode with GLM-5.3."

## What this answers

Your brief's owner-decision boundary said to ship the mechanism with the global default UNSET and to ask the
owner which agent should be the default implementer. The owner has now answered: **codex-1** (Codex CLI,
`gpt-6-astra`). The mechanism still ships with the protocol-level default unset (behaviour without a setting
stays today's); after release, set the owner's global default implementer to `codex-1` in the owner's global
configuration through the mechanism you ship, and verify it with the command that displays it.

## Already applied by the owner-facing session (2026-09-25)

`~/.parley/agents.toml` (machine roster): claude-1 = `claude/claude-opus-5-5[1m]`, effort max, INACTIVE
(organizer, not a quorum member); codex-1 ACTIVE (`gpt-6-astra`); kimi-1 and zcode-1 unchanged and active.
The global instruction files (`~/.claude/CLAUDE.md`, `~/.codex/AGENTS.md`) now name claude-1 as organizer and
codex-1 as implementer for NEW runs. This does not change your in-flight run: its frozen roster and
participants (claude-1, kimi-1, zcode-1) and you as its organizer stay as they are.
