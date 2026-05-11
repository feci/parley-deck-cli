---
agent: hermes
idea: agent-runtime-config
round: 1
date: 2026-05-11
---

## Summary
Define explicit runtime config surface in parley-deck-cli for stable cooperation. Expose agent ID, CLI path, headless mode, workspace-write sandbox, approval=on-failure, cli-default model/reasoning, 30min timeout, isolated home, and backend disclosure. Codex profile: sandbox_mode=workspace-write, approval_policy=on-failure, verify Git writes in-repo, retry exact failed Git cmds via approval, no alternate workdir without consent. Trade-off: small explicit config vs current hidden env state.

## Proposed approach
1. Add `parley config` subcommand to print effective runtime values (agent-id, cli-path, sandbox, approval, model=cli-default, timeout=1800000ms, isolated-home).
2. Update docs/agent-runtime-configuration.md with Codex guidance and user-consent prompts for Git retry / no-workdir-workarounds.
3. Enforce narrow workspace-write + on-failure approval as defaults for new agents; surface external-backend flag.
4. Keep scope to one design + one impl PR; preserve Codex/Claude/Gemini/Hermes support.

## Concerns / open questions
- How to detect/verify Git writes without over-approving?
- Isolated-home enforcement across OSes (macOS vs Linux paths)?
- CLI-default model may need future per-agent override without violating "no guessed models".

## Risks
- Overly strict approval flow blocks legitimate Git ops if retry path is incomplete.
- Users may still prefer alternate workdirs for speed; explicit consent gate adds friction.