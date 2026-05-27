---
agent: agy
idea: antigravity-agent-migration
review-round: 1
date: 2026-05-27
reviewed-commit: 697cb66
---

# Summary

The current implementation of the Gemini-to-Antigravity migration successfully establishes the Antigravity CLI (`agy`) as a first-class headless agent in the `parley-deck-cli` runtime catalog and updates the `parley-deck-skill` installer to support it as a native plugin target. 

The implementation aligns well with the specifications in [FINAL.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/antigravity-agent-migration/FINAL.md) and [IMPLEMENTATION.md](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/ideas/antigravity-agent-migration/IMPLEMENTATION.md). Specifically, `gemini` has been deprecated and removed from active defaults, examples, and TUI rosters, while remaining discoverable for backward compatibility. The new `agy` spec matches the locally verified headless execution flags. Claude and Hermes configurations have been successfully bumped to their strongest verified defaults, and the skill installer correctly generates the required `plugin.json` metadata, copies the skill artifact to the expected subfolder (`skills/SKILL.md`), and validates cleanly with `agy plugin validate`.

No **CRITICAL** or **MAJOR** findings were identified. All unit tests pass across both repositories. A few minor quality items, open questions, and release coordination steps are detailed below.

# Findings

No **CRITICAL** or **MAJOR** findings.

### 1. Hard fail on `isolate_home` for `agy` [MINOR]
- **File & Line**: [internal/runner/runner.go:L482-491](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/runner/runner.go#L482-L491)
- **Description**: While `gemini` and `hermes` have custom isolated home creation logic (such as copying global OAuth creds), `agy` does not have a built-in isolation strategy. If a user sets `isolate_home = true` for `agy` (e.g. in `agents.local.toml`), the execution will hit the `default` switch case. Because `agy` has no default `IsolatedHomeEnv` mapped, `len(agent.IsolatedHomeEnv)` is `0`, leading to a hard execution failure with `no isolated home strategy for agy`.
- **Suggested Fix**: Fall back to creating a standard temporary directory without environment mapping instead of throwing an error, or document this limitation clearly in the runtime docs so users know they must supply a custom `IsolatedHomeEnv` mapping if they enable `isolate_home` for `agy`.

### 2. Argument parsing order sensitivity [MINOR]
- **File & Line**: [internal/agents/discover.go:L146](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/agents/discover.go#L146)
- **Description**: In the built-in `agy` specification, `HeadlessArgs` is defined as:
  ```go
  []string{"--print-timeout", "30m", "--dangerously-skip-permissions", "--add-dir", "{root}", "--print", "{prompt}"}
  ```
  The `agy` CLI requires `--print` to take an argument (the prompt itself). Placing `{prompt}` at the very end as the argument value for `--print` is correct and allows Go's flag library to parse it cleanly. However, if any other arguments are appended or overridden locally, `--print` might consume subsequent arguments.
- **Suggested Fix**: Keep the prompt flag at the end as implemented, but add a brief code comment in `discover.go` to document that `--print` is flag-value-order sensitive.

### 3. Inactive Roster Role Mapping [NIT]
- **File & Line**: [internal/protocol/defaults/COOPERATION.md:L77](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/internal/protocol/defaults/COOPERATION.md#L77) & [parley-deck/COOPERATION.md:L77](file:///Volumes/My%20Shared%20Files/AI_WORKSPACE/parley-deck/parley-deck-cli/parley-deck/COOPERATION.md#L77)
- **Description**: The roster table correctly marks `gemini` as `inactive legacy`. However, the host accounts mapping table below it in both files still retains the `gemini | not mapped` row. While harmless, the legacy row could be annotated to match the inactive status or removed if host mapping is irrelevant for legacy runtimes.
- **Suggested Fix**: Update comments or notes in the document to clarify that the mapping is retained solely for legacy compatibility.

# Open questions

1. **Future model control for `agy`**: Currently, `agy` does not expose model or thinking flags in its `--help` output, so it is locked to `cli-default`. If the CLI adds support for model selection in a future version (e.g. `--model` or similar), will we update the built-in specs, or is the expectation that users will override these flags via `agents.local.toml`?
2. **Legacy `gemini` cleanup window**: How long do we intend to keep the legacy `gemini` target in the skill installer and CLI catalog? Will it be completely removed in the next major version bump (e.g. `v2.0.0` / `2.0.0`)?

# Residual test/release risks

1. **OAuth credential dependency**: Because `agy` runs in the user's home environment by default (no `isolate_home` configuration), its execution is coupled with the global `~/.gemini/oauth_creds.json` and `~/.gemini/google_accounts.json` state. If a user's credentials expire, the headless run will hang or fail without warning.
2. **WinGet and Homebrew publishing sequence**: Updating Homebrew formulas and WinGet manifests requires the package tarball and Windows executable hashes. This creates a release risk if those updates are forgotten or delayed post-publish. We should ensure the release pipeline checklists (in `RELEASING.md`) enforce that these registry updates are executed immediately after release assets are built.
