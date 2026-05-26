---
idea: meta-protocol-change-roster-headless-config
status: final
author: codex
consensus-date: 2026-05-25
participants: [codex, claude, gemini, hermes]
---

## Final plan / specification

Update `COOPERATION.md` section 2 so the project roster uses the active agent
IDs `codex`, `claude`, `gemini`, and `hermes` instead of placeholders. Use
portable logical workspace references and mark host handles as not mapped.

Add an advisory note that machines may keep ignored local launch configuration
such as `parley-deck/meta/headless-agents.local.json`. That file is optional,
machine-local, not canonical project state, and does not change quorum,
ownership, signoff weight, or transport rules.

Do not add a tracked schema or example file in this protocol change. Defer that
to a later idea if it becomes necessary.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
