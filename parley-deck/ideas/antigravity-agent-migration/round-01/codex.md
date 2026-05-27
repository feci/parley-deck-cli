---
agent: codex
idea: antigravity-agent-migration
round: 1
date: 2026-05-27
---

## Summary

The migration should replace Gemini as a default/runtime target with
Antigravity `agy`, while preserving the ability for advanced users to configure
legacy Gemini explicitly. The scope spans two repositories: `parley-deck-cli`
for agent runtime defaults, project bootstrap, docs, and tests; and
`parley-deck-skill` for install targets, package metadata, README, release
assets, Homebrew, and WinGet manifests.

## Proposed approach

- Add `agy`/Antigravity as a first-class built-in CLI agent in
  `parley-deck-cli`, using `agy --print` with explicit workspace access and a
  long print timeout.
- Remove Gemini from default participants and project shared defaults. If
  retaining compatibility, make Gemini legacy/deprecated documentation only, not
  part of default discovery examples.
- Update `parley-deck/agents.toml`, `COOPERATION.md`, protocol changelog, docs,
  README examples, CLI reference examples, and tests from `gemini` to `agy`.
- Update the embedded default `COOPERATION.md` used by `parley init` so new
  workspaces no longer inherit a Gemini roster.
- In `parley-deck-skill`, add an Antigravity target that installs into the
  current Antigravity skill/plugin location, prefers `agy` command detection,
  and keeps Gemini only as a legacy target if needed.
- Update skill package metadata, README, installer help, validation, and tests.
- Keep model/thinking defaults honest:
  - Claude can use explicit `--model opus --effort max`.
  - Hermes can use a configured strong model when the project has verified it.
  - Antigravity `agy --help` currently exposes no model/thinking flags, so its
    default should be `cli-default` unless local config overrides it.
  - Codex should remain configurable and should not invent a model string unless
    local discovery or user config proves it.

## Concerns / open questions

- Antigravity's durable skill/plugin directory is newer and may still be
  changing. The installer should use the documented current path and keep target
  validation focused on files it controls.
- `agy --print` works in this local environment, but third-party reports mention
  stdout capture differences on Windows. WinGet/Windows release notes should
  call for manifest validation on Windows rather than assuming macOS behavior.
- The request says "all Gemini mentions", but historical completed
  `parley-deck/ideas/` files are audit records and should not be mass edited.

## Risks

- Removing Gemini from built-ins could surprise users who still rely on it.
  Mitigate with release notes and local override documentation.
- Publishing both CLI and skill in one pass crosses Go, npm, Homebrew, and
  WinGet surfaces. Keep verification explicit per package.
- If Antigravity changes its plugin/skill layout again, the installer target may
  need a follow-up release.
