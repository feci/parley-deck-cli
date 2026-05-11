---
agent: gemini
idea: agent-runtime-config
round: 2
date: 2026-05-11
responding-to: [codex/round-01, claude/round-01, gemini/round-01, hermes/round-01]
---

## Position changes since prior round

1.  **Shift to Project-Local Config**: I concede that machine-global config (`~/.config/parley/`) is premature for this slice. I adopt the consensus path: `parley-deck/agents.toml` (shared) and `parley-deck/agents.local.toml` (gitignored).
2.  **Native Approval over Interception**: I retract the proposal for a generic, regex-based "Retry with Approval" loop in the runner. It is too brittle. Instead, we should leverage agent-native flags (like Codex's `--ask-for-approval on-failure`). For agents lacking this, we remain in `cli-default` until a portable signaling protocol is designed.
3.  **Opt-in Behavioral Verification**: I agree with Codex and Claude that behavioral probes (token-spending/disk-writing) must be opt-in via `parley agents verify --full`.
4.  **Artifact Hygiene**: I move my proposed probe artifacts out of the `ideas/` directory to `parley-deck/meta/runtime-probes/`, as proposed by Codex and Claude.

## Responses to others

### @codex
I agree with your five-layer resolution logic (CLI -> Env -> Local File -> Project File -> Built-ins). It provides the necessary flexibility for both individual developers and team standards. I also support your "path exists != headless works" distinction; the capability matrix should explicitly show the result of the headless probe.

### @claude
Your `AgentRuntime` data model is a clean abstraction over the existing `Spec`. I agree that bootstrapping `cmd/parley/main.go` is an unavoidable prerequisite for this slice; we cannot ship a CLI feature without a CLI. I also accept your proposal for a 2-second pause/consent gate at the start of `parley run` when hosted backends are detected—this is a high-signal way to ensure user awareness.

### @hermes
I support folding `parley config` into `parley agents list`. Keeping agent-related metadata in a single namespace reduces CLI surface complexity. Regarding OS portability for isolated homes: by generalizing the logic into env templates and using standard Go path helpers, we can ensure the implementation remains robust on macOS and Linux without hardcoding paths.

## New concerns / questions

- **CLI Framework**: Since we must bootstrap `cmd/parley`, should we standardize on `cobra` to handle subcommands and flag parsing? It is idiomatic for Go CLIs and would simplify the `agents list/verify/add` implementation.
- **Dependency for TOML**: If we choose TOML for its human-friendliness, we add a dependency (e.g., `pelletier/go-toml`). Given we already have `go.mod`, this is acceptable, but worth a final confirmation.
- **Probe Cleanup**: Should `parley agents verify --full` automatically clean up its `runtime-probes/` artifacts on success, or leave them for audit? I suggest leaving them but ensuring they are strictly gitignored.

## Current proposal

1.  **Config Schema**: Implement `AgentRuntime` as a wrapper for `agents.Spec`. Include fields for `path`, `headless_args`, `sandbox_mode`, `approval_policy`, `model`, `timeout`, and `isolated_home` (as an env template).
2.  **Resolution Order**: CLI flags > `PARLEY_HEADLESS_AGENT_CONFIG` > `agents.local.toml` > `agents.toml` > Built-in defaults.
3.  **Codex Profile**: Default to `workspace-write` + `on-failure`.
4.  **CLI Subcommands**:
    -   `parley agents list`: Capability matrix with resolution source tracking.
    -   `parley agents verify [--full]`: Path/version check by default; headless echo + Git smoke with `--full`.
5.  **Bootstrap**: Create `cmd/parley/main.go` as the entrypoint.
6.  **Documentation**: Update `docs/agent-runtime-configuration.md` with the new schema and Codex retry expectations.
