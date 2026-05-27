---
agent: agy
idea: antigravity-agent-migration
round: 1
date: 2026-05-27
---

# Summary

Gemini CLI is being phased out in favor of Antigravity CLI (`agy`). As the `agy` participant representing the Antigravity CLI lens, this proposal outlines the migration strategy to establish `agy` as a first-class agent target. This includes built-in agent discovery configurations, workspace sandboxing, installer target specifications in `parley-deck-skill`, default roster updates, and packaging coordinates.

# Proposed approach

### 1. Antigravity CLI Built-in Spec & Execution Details
- **Discovery Registration**: Add `agy` to `DefaultSpecs` in `internal/agents/discover.go` with ID `"agy"`.
- **Headless Invocation**: Configure `LaunchMode` as `LaunchHeadless` with the command `agy`.
- **Arguments**: Use `agy --print` (or `-p`) for non-interactive execution. The headless arguments will include:
  - `--add-dir {root}` to supply workspace context.
  - `--dangerously-skip-permissions` to bypass interactive tool approval prompts during automated runs.
- **Defaults**: Report `cli-default` for model, reasoning, and profile since `agy` does not expose these controls in its `--help` output. Keep `balanced` speed and a standard `1800000ms` (30-minute) timeout.
- **Home Isolation**: Map environment isolation to a temporary home directory when `isolate_home` is enabled. Since `agy` shares its core config structures with Gemini, copy existing user oauth credentials (`oauth_creds.json`, `google_accounts.json`) into the isolated path.

### 2. Parley Deck Skill Installer Target
- **Installer Spec**: Update `parley-deck-skill`'s `installer.js` to add an `agy` target.
- **Installation Path**: Direct the installation to `~/.gemini/config/plugins/parley-deck`, matching the target directory structure discovered via the `agy plugin` subcommands.
- **Plugin Metadata**: Create a `plugin.json` manifest distributed with the skill payload. It will specify:
  ```json
  {
    "name": "parley-deck",
    "version": "1.1.1",
    "description": "Run Parley Deck multi-agent idea, implementation, review, and consensus workflows.",
    "skills": ["SKILL.md"],
    "agents": ["agents/manifest.yaml"]
  }
  ```
- **Tests**: Expand `installer.test.js` to assert directory detection, installation validation, and marker existence for `agy`.

### 3. Repository Defaults & Roster Migration
- **Roster & Cooperation**: Update `COOPERATION.md` default agent roster to list `agy` instead of `gemini`. Mark `gemini` as deprecated/legacy rather than a primary participant.
- **Workflow Templates**: Update `parley init`, cli references, and README documentation examples to default to `--participants claude,agy` instead of `claude,gemini`.
- **CLI Tests**: Modify relevant test files (`runner_test.go`, `discover_test.go`, `live_test.go`) to ensure all test suites pass with the new default configurations.

### 4. Release Coordination
- Prepare matching updates for the Homebrew formulas and WinGet manifests.
- Update packaging metadata in both `parley-deck-cli` and `parley-deck-skill` to reflect `agy` integration.

# Concerns / open questions

- **CLI Config Keys**: Does `agy` support an explicit environment override (e.g. `ANTIGRAVITY_CLI_HOME` or `ANTIGRAVITY_HOME`) similar to `GEMINI_CLI_HOME`?
- **Interactive Hand-off**: Does `agy` need an interactive launch mode support in Parley (like `--prompt-interactive` or `-i`) or is headless `--print` sufficient for all cooperation rounds?
- **Legacy Deprecation Window**: Should `gemini` CLI support be completely removed or kept as a fallback for one minor release cycle?

# Risks

- **Permission Escapes**: If `agy` attempts to write files outside of the `--add-dir` workspace configuration or if the sandbox restrictions conflict with the local environment, automated runs may hang or error out. We must thoroughly verify behavior on clean systems.
- **Plugin Registry Contamination**: Since `agy` plugin directories are placed in `~/.gemini/config/plugins`, they share namespaces with other Gemini-based tools. We must ensure no naming collisions occur.
