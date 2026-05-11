---
agent: claude
idea: agent-runtime-config
review-round: 1
date: 2026-05-11
reviewed-commit: d38e7f738471caa61042ca087c445b4df5bd2606
---

## Summary

The implementation lands the FINAL.md scope well: `cmd/parley/main.go`, TOML layering with source tracking, Codex `workspace-write` + `on-failure` defaults, `agents list|verify`, consent-gated `--full` with probe sentinel, Codex Git smoke, threaded runtime values into the runner, and updated `.gitignore` + docs. Tests cover the load order, source attribution, Codex defaults, the cheap verify path, and the Git-smoke failure case. The one substantive defect is that for `gemini` and `hermes` — the only built-in agents that actually use isolated homes — the configured `isolated_home_env` template is silently ignored by the runner. There are also some smaller display/source-tracking warts worth tightening before fix-up closes.

## Findings

### [CRITICAL] <short title>

None.

### [MAJOR] Configured `isolated_home_env` is ignored for gemini and hermes

`internal/runner/runner.go:429–458` (`isolatedAgentHome`) switches on `agent.ID` and unconditionally returns a hardcoded `GEMINI_CLI_HOME=…` / `HERMES_HOME=…` pair built from `isolatedGeminiHome()` / `isolatedHermesHome()`. The `default:` branch (which is the only path that reads `agent.IsolatedHomeEnv`) is unreachable for the two agents that actually declare `isolate_home = true` in the built-ins (`internal/agents/discover.go:111–116, 131–136`) and in `parley-deck/agents.toml:28–32, 40–44`.

Consequences:

- A user-local override in `agents.local.toml` such as
  ```toml
  [agents.gemini.isolated_home_env]
  GEMINI_CLI_HOME = "{tempdir}/parley-gemini"
  OTHER_ENV = "{deck}/cache"
  ```
  is parsed, threaded into `Spec.IsolatedHomeEnv`, reported by `parley agents list`, and recorded in the `run.created` event — but has zero effect at launch. The runner only emits the hardcoded `GEMINI_CLI_HOME=` env var.
- This contradicts `docs/agent-runtime-configuration.md:117–122` ("Configured isolated-home environment values may use placeholders: `{root}`, `{deck}`, `{tempdir}`") and FINAL.md's requirement that the resolved runtime — including the isolated-home env template — drive launches.
- It also re-creates the "hidden runtime" trap the idea was supposed to fix: the user sees a configured value, the CLI ignores it.

Suggested fix: keep the gemini/hermes credential-copy helpers (`isolatedGeminiHome` / `isolatedHermesHome` — they do non-trivial work copying OAuth/config files that pure env vars cannot replace), but build the returned `env` slice from `agent.IsolatedHomeEnv` using the temp `home` path the helper just created, instead of returning a hardcoded `{KEY}=home` pair. If no template is configured for the agent ID, fall back to the historical hardcoded var name. That preserves the config-as-source-of-truth contract while keeping the special credential bootstrap.

### [MAJOR] `HEADLESS` column in `parley agents list` is always `not-probed`

`internal/agents/discover.go:220–267` (`PrintRuntimeMatrix`) declares `headless := "not-probed"` at line 225 and never writes anything else to that variable before printing. The column header (line 221) is `HEADLESS`. The column adds noise without information and gives the false impression that `agents list` ran some kind of headless check.

This is the column the slice was explicitly meant to clarify: round-02 and FINAL.md both call out distinguishing `installed`, `version probed`, and `headless probed`. As implemented, the matrix only reports `installed` and version status — `headless` is unconditionally a placeholder. If headless-probed state is genuinely out of scope for the cheap `list` path, the column should be removed (and the headless invocation string moved to a sources/notes line) rather than printed as a constant string that contradicts the user's reading of the header.

Suggested fix (minimal): drop the `HEADLESS` column from `PrintRuntimeMatrix`, or rename it to `HEADLESS-PROBE` and only emit a non-`not-probed` value after a `verify --full` artifact exists for the current run.

### [MINOR] Shared `parley-deck/agents.toml` re-declares built-in defaults, polluting source attribution

`parley-deck/agents.toml:4–44` sets `model = "cli-default"`, `reasoning = "cli-default"`, `profile = "cli-default"`, `speed = "balanced"`, `timeout_ms = 1800000`, and (for codex) `sandbox_mode = "workspace-write"` / `approval_policy = "on-failure"` — all values that already match `DefaultSpecs()` in `internal/agents/discover.go`. Because `applyOverride` in `internal/config/runtime.go:174–193` updates `Sources["model"]` etc. whenever the override value is non-empty, the source for these fields shifts from `built-in` to `parley-deck/agents.toml` purely because the same string was re-stated.

Net effect: `parley agents list` reports the wrong provenance — users following a source trail to understand "why is my codex `approval_policy` `on-failure`?" land on `parley-deck/agents.toml` rather than the built-in default the slice was supposed to install. The TOML file then looks load-bearing when it isn't.

Suggested fix: either delete the redundant fields from `parley-deck/agents.toml` (leaving only fields the project genuinely wants to lock in, e.g. nothing for now), or change `applyOverride` to only stamp the source when the resolved value actually differs from the current spec value. The first is simpler.

### [MINOR] New TOML-declared agents lose source attribution for unset fields

In `internal/config/runtime.go:111–124`, when a TOML file introduces an agent ID that isn't in `DefaultSpecs()`, the new `Spec` is seeded with defaults (`PromptStdin`, `CLIDefault`, `DefaultSpeed`, `DefaultTimeoutMS`, `ExternalUnknown`) but `Sources` is initialised to just `{"id": source}`. Fields the TOML file doesn't override therefore have no `Sources` entry, and `sourceFor()` (`internal/agents/discover.go:332–337`) falls back to the literal string `"discovered"`. For a freshly added agent, the displayed source is `discovered`, which is misleading — those values came from `LoadAgentSpecs` config-defaults, not from `Discover()`.

Suggested fix: when synthesising a new spec for an unknown ID, populate `Sources` for every defaulted field with a sentinel like `"config-default"` (or the source file itself, since the agent was declared there). Add a test that asserts the source for an untouched field of a TOML-introduced agent is not `"discovered"`.

### [MINOR] `{tempdir}` placeholder has two different meanings depending on the field

`internal/config/runtime.go:135` sets `tempdir := os.TempDir()` and uses it for `Commands`, `VersionArgs`, `HeadlessArgs`, etc., which substring-expand `{tempdir}` at config-load time to the OS-wide temp directory. The same file at line 201 deliberately passes the literal `"{tempdir}"` for `IsolatedHomeEnv` so the runner (`internal/runner/runner.go:452–454`) can later expand it to the per-agent isolated home.

So `{tempdir}` in `[agents.foo.headless_args]` resolves to `os.TempDir()` once and never changes per run, while the same placeholder in `[agents.foo.isolated_home_env]` resolves to a fresh `parley-<id>-home.*` directory per launch. `docs/agent-runtime-configuration.md:121` documents only the second meaning. A user templating a path into `headless_args` will silently get the wrong directory.

Suggested fix: pick one. Either (a) preserve `{tempdir}` literally in all fields and expand it at launch time inside `CommandFor`, or (b) document the split explicitly in `docs/agent-runtime-configuration.md` and rename the load-time placeholder (e.g. `{ostempdir}`) so the two are distinguishable.

### [MINOR] `runFullVerification` fails fast on the first error

`internal/app/app.go:418–429`: as soon as the Codex Git smoke or any single headless probe returns an error, the loop bails out. With `verify --full` (no `--agent`), a user with three configured agents only ever sees the first failure; the other two remain unverified until that one is fixed. Given that probes already write artifacts under `parley-deck/meta/runtime-probes/<run-id>/`, it would be cheap to accumulate errors and report all of them at the end (still exiting non-zero).

Suggested fix: collect errors per agent, print each, and return a single non-zero status at the end of the loop.

### [MINOR] `parley agents verify` (no `--agent`, no `--full`) returns non-zero whenever any configured agent is missing from `PATH`

`internal/app/app.go:166–181`: `failed = true` whenever a selected agent has `!result.Found`. With four built-in specs, a developer who only has `codex` installed will always get exit code 1 from `parley agents verify`, even though their codex is healthy. This may be intentional ("verify checks the whole quorum"), but the behaviour isn't documented in `docs/agent-runtime-configuration.md` and feels surprising next to `agents list`, which cheerfully reports `installed=no` rows without failing.

Suggested fix: either document the strict semantics in the docs (and ideally print a one-line hint pointing the user at `--agent ID` for partial verification), or treat missing agents as informational unless `--agent` is supplied.

### [MINOR] No app-level test that `parley run` actually injects resolved runtime into the launched command

The runner-level `TestBuildRoundOnePrompt` (`internal/runner/runner_test.go:18–68`) covers the prompt text and `TestRunRoundOneCreatesArtifactWithHeadlessAgent` covers the fake-agent integration, but FINAL.md lists "`parley run` using resolved runtime values rather than hardcoded defaults" as a test target and there's no app-layer test that exercises `Run([]string{"run", ...})` end-to-end against a fake CLI to confirm the resolved sandbox/approval/model are present in the run event log or in the constructed `exec.Cmd`. The behaviour is plumbed (`internal/app/app.go:548–568` `runtimeEventData`, `internal/runner/runner.go:319–369` prompt builder), but the test surface stops at the runner package.

Suggested fix: add a small `app` test that invokes `Run([]string{"run", "--no-tui", "--yes", "--dir", root, "test task"})` against a fake codex CLI on `PATH`, then asserts that the `run.created` event in `events.jsonl` carries the resolved `sandbox_mode`, `approval_policy`, and `model` (and that they came from configured overrides rather than built-ins).

### [NIT] Compatibility aliases `agents discover|probe` are not exercised by tests

`internal/app/app.go:108–117` keeps `discover` and `probe` as aliases for `list` and `verify`, but neither alias is touched by `internal/app/app_test.go`. Compatibility aliases tend to bit-rot silently. A one-line test (`Run([]string{"agents", "discover", "--dir", root}, ...)` returns 0 with the same matrix line) would freeze the alias contract.

### [NIT] `runtimeEventData` and the prompt builder duplicate small helpers

`internal/app/app.go:570–575` (`valueOr`), `internal/agents/discover.go:318–323` (`valueOrDefault`), and `internal/runner/runner.go:413–418` (`runtimeValue`) all implement the same "empty-string → fallback" logic with slightly different fallbacks and signatures. Not a bug, but consolidating into one helper exported from `internal/agents` would reduce drift the next time `cli-default` semantics change.

### [NIT] `agents list` output relies on space padding rather than `text/tabwriter`

`internal/agents/discover.go:238–251` uses fixed-width `%-Ns` columns. Long version strings get truncated by `truncate(...)`, and very short ones leave wide gaps. `text/tabwriter` would auto-size and keep the matrix readable when an agent ID, version, or sandbox value grows. Not load-bearing, but cheap to switch.

## Open questions

- Is the strict exit-1 semantics of `parley agents verify` (without `--agent`) intentional, or a side-effect of the simple loop? If intentional, the docs should say so explicitly; if accidental, the easy fix is to soft-warn.
- For the [MAJOR] isolated-home gap on gemini/hermes: do we want to thread the configured env through the existing credential-copy helpers in this slice, or accept it as a known limitation and document the constraint clearly in `docs/agent-runtime-configuration.md` and `IMPLEMENTATION.md` so users know the template is informational for those two agents? I lean toward the former, but it's a judgment call worth surfacing in review-consensus.
- Should `{tempdir}` be unified (run-time per-agent) or kept split (load-time OS, run-time isolated)? Either is defensible, but the current state is the worst of both — the same token means two different things depending on which TOML key it appears under, and only one meaning is documented.
