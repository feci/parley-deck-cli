---
idea: agent-runtime-config
drafted-by: codex
date: 2026-05-11
---

## Agreed decisions

- Implement this as the next `parley-deck-cli` slice after `hitl-tui-questions`.
- Add a minimal `cmd/parley/main.go` entrypoint that delegates to the existing `internal/app.Run`. Do not introduce Cobra in this slice; keep the command routing aligned with the current stdlib `flag` style.
- Add an `AgentRuntime` / `ResolvedAgent` layer over the existing `agents.Spec` and `agents.Discovery` structs instead of replacing the adapter model wholesale.
- Use project-local TOML configuration:
  - `parley-deck/agents.toml` for checked-in shared defaults.
  - `parley-deck/agents.local.toml` for gitignored user-local overrides.
- Resolve configuration in this order, highest precedence first:
  1. Explicit CLI flag where the command already has one.
  2. `PARLEY_HEADLESS_AGENT_CONFIG`.
  3. `parley-deck/agents.local.toml`.
  4. `parley-deck/agents.toml`.
  5. Built-in defaults.
- Track the source of every resolved runtime field so the CLI can explain where values came from.
- Runtime fields must cover at least: stable agent ID, CLI path/command candidates, headless invocation, prompt mode, sandbox mode, approval policy, model, reasoning/profile, timeout, isolated-home env template, external-backend class, telemetry note, and notes.
- Change the built-in Codex runtime default to `sandbox_mode = "workspace-write"` and `approval_policy = "on-failure"`.
- Preserve `cli-default` for model, reasoning, and profile when the CLI cannot prove supported options. Do not invent model names or reasoning levels.
- Keep the no-silent-alternate-workdir rule: if a configured agent cannot write in the target repository, the CLI must surface the failure instead of silently relocating protocol work elsewhere.
- Replace or alias the current `parley agents discover|probe` UX with:
  - `parley agents list`: cheap capability matrix and resolved-runtime display.
  - `parley agents verify [--agent ID] [--full] [--yes]`: verification command.
- `parley agents list` must remain cheap and non-spending: path lookup, version probe, resolved runtime, and clear columns for installed/version/headless status.
- `parley agents verify` without `--full` must stay cheap: path lookup and version probing only.
- `parley agents verify --full` may run behavioral headless probes and must require explicit `--yes` for agents whose `external_backend` is not `local`.
- The Codex full verification path must include the Git write smoke sequence in the target repository:
  - `git status`
  - `git branch tmp-codex-git-test`
  - `git branch -D tmp-codex-git-test`
  - `printf test | git hash-object -w --stdin`
- If the Codex Git smoke fails because Codex itself needs approval, the retry must be handled by Codex's native `on-failure` runtime behavior. The CLI must not implement a generic stderr-pattern "retry with approval" engine in this slice.
- Full verification probe artifacts must be written under `parley-deck/meta/runtime-probes/<run-id>/` and this path must be gitignored. Probe artifacts must include a sentinel that lets the verifier confirm the intended agent wrote the current probe output.
- `parley run` must show the resolved capability matrix before launch. `--yes` or `--auto` can skip any interactive confirmation, but the effective runtime should still be visible in stdout or the run event log.
- Update `docs/agent-runtime-configuration.md` as the user-facing reference for adding agents, Codex's required runtime contract, verification commands, consent rules, and isolated-home behavior.
- Do not update `parley-deck/meta/protocol-changelog.md`; this slice changes CLI behavior and documentation, not the Parley cooperation protocol.

## Agreed trade-offs

- TOML is chosen for human-edited configuration even though it adds a Go dependency. The implementation should add one focused TOML dependency rather than supporting both TOML and JSON config formats.
- Project-local configuration is enough for this slice. Machine-global configuration may be useful later, but it would make the effective runtime harder to audit now.
- A minimal `cmd/parley` entrypoint is in scope because otherwise the new `parley agents` commands are not runnable from a fresh clone. A larger CLI framework migration is out of scope.
- Behavioral verification is opt-in because it can spend hosted-agent tokens and write temporary files.
- The CLI should expose and verify agent-native approval/sandbox settings where known. It should not claim portable retry semantics across agents that do not expose equivalent runtime controls.
- The implementation should keep the diff small by layering resolved runtime data around current discovery and runner code.

## Open items deferred to implementation

- Exact TOML package choice.
- Exact table formatting for the capability matrix.
- Whether old `agents discover|probe` commands remain as compatibility aliases or are replaced outright.
- The smallest useful sentinel format for `meta/runtime-probes/<run-id>/`.

## Deferred follow-ups

- Interactive `parley agents add` wizard.
- Machine-global config under `~/.config/parley/`.
- Generic model discovery and automatic strongest-model selection.
- Per-run model/profile/timeout override flags beyond existing flags.
- Portable CLI-side approval retry engine.
- `parley resume` and richer runtime status/resume behavior.
- Automatic consensus progression policy.
- GitHub/GitLab transport automation.
- Token telemetry parsing and cost reporting.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-11
Status: ✅ ACCEPT
Notes: I accept the consensus as drafted. The final plan should preserve the minimal implementation shape: runtime config first, visible effective values, Codex `workspace-write` + `on-failure`, and no generic retry engine.
