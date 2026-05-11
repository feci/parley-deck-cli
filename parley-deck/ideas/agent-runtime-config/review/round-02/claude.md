---
agent: claude
idea: agent-runtime-config
review-round: 2
date: 2026-05-11
reviewed-commit: a4f1765c9f4e500750a7711756c78539c6eeb76c
responding-to: [review/round-01/claude, review/round-01/gemini, review/round-01/hermes, review/consensus]
---

## Summary

Fix-up cycle 1 lands cleanly. All nine Agreed fixes from `review/consensus.md` are visible in code and exercised by tests; both MAJORs are properly resolved, and no new blockers or regressions appear. The two Deferred NITs (helper consolidation, `text/tabwriter`) remain as expected per consensus.

## Fix-up verification

- **[MAJOR] `isolated_home_env` ignored for gemini/hermes** — Fixed at `internal/runner/runner.go:429–465`. `isolatedAgentHome` now calls `isolatedHomeEnv(agent, home, "GEMINI_CLI_HOME")` / `isolatedHermesHome(...)`. The shared helper at lines 455–465 builds env from `agent.IsolatedHomeEnv` with `{tempdir}` expanded to the per-launch home, and falls back to the historical key only when no template is configured. Covered by `TestIsolatedHomeEnvUsesConfiguredTemplate` and `TestIsolatedHomeEnvFallsBackToHistoricalKey` (`internal/runner/runner_test.go:266–293`). Matches the suggested-fix shape from round-01.
- **[MAJOR] `HEADLESS` column always `not-probed` (Claude+Gemini)** — Fixed at `internal/agents/discover.go:220–272`. The column now reads `configured` when `len(result.HeadlessArgs) > 0` and `missing` otherwise; the `headless: <mode>` detail line surfaces the configured invocation. `TestAgentsListPrintsResolvedRuntime` asserts `configured` is present in the output (`internal/app/app_test.go:31`). The column now carries information that varies per agent and matches its header.
- **[MAJOR] Codex Git smoke runs directly in CLI instead of via agent (Gemini)** — Fixed at `internal/app/app.go:467–489`. `probePrompt` embeds the four Git smoke commands inside the Codex headless probe prompt; the CLI no longer `exec`s `git` directly (verified by grepping for `exec.CommandContext.*git` and `runCodexGitSmoke` — both gone). `TestCodexProbePromptIncludesGitSmoke` (`app_test.go:81–93`) freezes the contract. The Codex `on-failure` runtime is now the actor for any approval retry, per `FINAL.md`.
- **[MINOR] New TOML-declared agents lose source attribution** — Fixed at `internal/config/runtime.go:111–149`. Synthesized specs are seeded with `Sources: configDefaultSources(source)`, which stamps `source + ":default"` for every defaulted field. `TestLoadAgentSpecsLayersAndTracksSources` asserts the `extra` agent's `Sources["model"]` is neither empty nor `SourceDiscovered` (`runtime_test.go:84–86`).
- **[MINOR] `{tempdir}` placeholder split semantics** — Documented at `docs/agent-runtime-configuration.md:119–124` with both meanings called out explicitly (load-time OS temp vs per-launch isolated home). The existing tests (`TestExpandPlaceholders`, `TestLoadAgentSpecsLayersAndTracksSources` checking `IsolatedHomeEnv["EXTRA_HOME"] == "{tempdir}/extra"`) lock the dual behavior in.
- **[MINOR] `runFullVerification` fails fast** — Fixed at `internal/app/app.go:412–430`. Loop now collects per-agent failures, prints each `%s: full probe failed: %v`, and returns a single non-zero status reporting how many probes failed.
- **[MINOR] `parley agents verify` strict missing-agent behavior** — Documented at `docs/agent-runtime-configuration.md:95` and the CLI prints the hint `use --agent ID to verify one configured agent without requiring every built-in agent` when no `--agent` was supplied (`app.go:171–173`).
- **[MINOR] App-level test that `parley run` injects resolved runtime** — Added `TestRunRecordsResolvedRuntime` (`app_test.go:95–145`). It writes an `agents.local.toml` overriding `model = "local-model"`, invokes `Run([]string{"run", "--no-tui", "--yes", "--participants", "codex", ...})` against a fake codex CLI on `PATH`, loads the resulting `events.jsonl`, and asserts the `run.created` event's `runtime[0]` carries `model=local-model` and `approval_policy=on-failure`.
- **[NIT] Compatibility aliases not exercised** — Added `TestAgentsCompatibilityAliases` (`app_test.go:38–60`) covering both `agents discover` and `agents probe --agent codex`.

Dismissed/deferred items from consensus remain consistent with the code: `parley-deck/agents.toml` still re-declares `model = "cli-default"` etc. (intentional auditability, ack'd as non-blocking by my round-01 signoff); `valueOr` / `valueOrDefault` / `runtimeValue` still duplicated across packages (deferred NIT); `PrintRuntimeMatrix` still uses fixed-width padding rather than `text/tabwriter` (deferred NIT).

## Findings

### [CRITICAL] <short title>

None.

### [MAJOR] <short title>

None.

### [MINOR] <short title>

None.

### [NIT] <short title>

None.

## Open questions

None — all three open questions from round-01 were resolved in consensus and the resolutions are reflected in the code/docs. I am ready to sign off on the fix-up cycle.
