---
idea: antigravity-agent-migration
status: final
author: codex
consensus-date: 2026-05-27
participants: [codex, claude, agy, hermes]
---

## Final plan

Implement the Gemini-to-Antigravity migration as a compatibility-preserving
default change across both maintained packages.

### CLI repository

- Add built-in `agy` agent support using the locally verified headless command:
  `agy --print --print-timeout 30m --add-dir {root}
  --dangerously-skip-permissions`.
- Keep `gemini` discoverable as a deprecated legacy agent, but remove it from
  active defaults, examples, and new-workspace rosters.
- Update project and embedded protocol defaults so `agy` is active and
  `gemini` is inactive/legacy.
- Update shared `parley-deck/agents.toml` defaults:
  - `claude`: strongest verified default is `opus` with `max` effort.
  - `agy`: `cli-default` model/reasoning/profile because current help exposes
    no supported model or thinking flags.
  - `hermes`: strongest locally verified default is `xai/grok-4.3`.
  - `codex`: remain `cli-default` unless local config overrides it.
- Update CLI docs, help examples, runtime configuration docs, tests, and
  release notes.

### Skill repository

- Add an Antigravity install target named `agy` using the current plugin
  location shape `~/.gemini/config/plugins/parley-deck`.
- Add Antigravity plugin metadata and validate it with `agy plugin validate`.
- Keep the Gemini extension payload and target as legacy compatibility during
  the deprecation window.
- Update README, installer help output, tests, npm metadata, release notes,
  Homebrew metadata, and WinGet manifests/descriptions where applicable.

### Release/deploy

- Verify both repos before publishing.
- Publish package releases only when the local credentials and channel tooling
  actually succeed.
- For WinGet, update manifests and open or prepare the required PR. Do not
  report WinGet deployment as complete unless that publish path is actually
  executed.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
