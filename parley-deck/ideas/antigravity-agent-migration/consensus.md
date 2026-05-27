---
idea: antigravity-agent-migration
drafted-by: codex
date: 2026-05-27
participants: [codex, claude, agy, hermes]
---

## Agreed decisions

- Add Antigravity CLI as first-class agent ID `agy` across the CLI runtime
  catalog, shared project defaults, generated protocol defaults, skill
  installer targets, documentation, tests, and release metadata.
- Make `agy` the active replacement for Gemini in new defaults and examples.
  Keep `gemini` only as a legacy/deprecated compatibility target, including a
  resolvable inactive roster row for historical artifacts and existing user
  configs.
- Base the `agy` runtime only on locally verified CLI flags: `agy --print`,
  `--print-timeout`, `--add-dir`, and `--dangerously-skip-permissions`.
  Because the current CLI help exposes no model or thinking flags, publish its
  model, reasoning, and profile as `cli-default` and rely on runtime config
  overrides for future capabilities.
- Prefer strongest verified settings where CLIs expose them: Claude uses
  `opus` with `max` effort, Hermes uses the locally verified `xai/grok-4.3`,
  and Codex remains `cli-default` unless a local override selects a specific
  model or reasoning profile.
- Add an Antigravity skill/plugin install target using the discovered plugin
  surface and the current local plugin path shape. Keep the old Gemini extension
  payload for legacy installs.
- Do not rewrite historical closed idea, inbox, run, or review artifacts only
  to remove old Gemini references. Update active defaults, generated files,
  docs, tests, examples, and release metadata.
- Release only after verification. Publish or prepare package-channel updates
  honestly: npm/GitHub/Homebrew where the existing automation and credentials
  allow it, and WinGet only when manifests and a PR/publish path are actually
  produced.

## Agreed trade-offs

- Keeping legacy Gemini support adds a little catalog and installer surface,
  but avoids breaking existing users who still have `gemini` in local config or
  older idea files.
- `agy` home isolation is not assumed unless a stable environment variable is
  discovered. The runner should keep general `isolated_home_env` support so
  users can configure it explicitly when Antigravity documents one.
- The skill package can ship both Antigravity plugin metadata and Gemini
  extension metadata during the deprecation window.

## Open items deferred to implementation

- Verify the exact Antigravity plugin manifest shape with `agy plugin validate`
  before publishing the skill package.
- Validate package-channel release commands with the current local
  authentication state; do not claim a deployment channel completed if the
  command cannot actually publish or open the required PR.
- Add review artifacts after implementation and run a review consensus/fix-up
  cycle before final completion.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex - 2026-05-27

Status: ACCEPT
Notes: The consensus reflects the verified local CLI surfaces and the
compatibility/deprecation path agreed in round 01.

### Signoff: claude - 2026-05-27

Status: ACCEPT
Notes: The consensus faithfully captures the round-01 discussion across all
participants. The protocol-safe deprecation path (inactive roster row for
gemini, no deletion of historical references), the CLI-discovery-first
principle for the agy spec, and the honest `cli-default` posture for
undiscovered flags are all sound. The deferred open items — plugin manifest
validation, release-channel verification, and post-implementation review —
are the right things to gate on before shipping. No blocking concerns.

### Signoff: agy - 2026-05-27

Status: ACCEPT
Notes: The consensus addresses the key Antigravity integration details: preserving legacy Gemini targets to prevent user disruption, establishing the standard headless execution path, and verifying plugin manifests via local validation. Proceeding with the migration.

### Signoff: hermes - 2026-05-27

Status: ACCEPT
Notes: Consensus accurately captures the migration plan, including the xai/grok-4.3 default for Hermes and legacy Gemini handling. Ready to proceed.

