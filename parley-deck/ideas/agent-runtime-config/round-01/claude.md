---
agent: claude
idea: agent-runtime-config
round: 1
date: 2026-05-11
---

## Summary

Turn agent runtime setup into a first-class, explicit part of the CLI surface. The current `internal/agents/discover.go` hard-codes a runtime profile per agent (commands, headless args, isolated-home, telemetry) and only checks `PATH` + `--version`. Cooperation reliability degrades the moment a user's installed CLI deviates from those baked-in defaults (auth model, sandbox flag, approval flag, model name, profile). This slice should expose the effective runtime per agent, persist user overrides in a project-local config, verify behavior end-to-end against the target repo (not just `--version`), and document the contract for adding a new agent. Codex specifically must use `sandbox_mode=workspace-write` + `approval_policy=on-failure`, verify Git writes in the target repo, and retry the exact failed Git command through the approval path on sandbox/permission denial — never silently relocate work.

## Proposed approach

Scope: one design PR (this idea) + one implementation PR. No M6, no `parley resume`, no transport automation, no release packaging.

1. **Explicit `AgentRuntime` data model.** Extend `agents.Spec` (or wrap it) into an `AgentRuntime` with the fields the prompt enumerates, every one resolvable from explicit config or `cli-default`:
   - `id` (stable kebab-case, matches §2 roster)
   - `path` (absolute, resolved via `LookPath` or override)
   - `headless.invocation` (mode + argv template, `{root}`/`{prompt}` placeholders)
   - `workspace_write` (narrow, repo-rooted: how this agent is told to write only inside the workspace — Codex `--sandbox workspace-write`, Claude `--add-dir {root}`, Gemini `--approval-mode auto_edit`, Hermes `--accept-hooks`)
   - `sandbox` + `approval` (e.g. Codex `workspace-write` / `on-failure`; others `cli-default` if not exposed)
   - `model` / `reasoning` / `effort` / `profile` (default `cli-default`; never invented)
   - `timeout_ms` (default `1_800_000`; override per agent; "balanced/fast/slow" speed alias separate field)
   - `isolated_home` ({path template, env vars} — already implicit for Gemini/Hermes; surface it)
   - `external_backend` (`hosted` | `local` | `unknown` + free-text disclosure string shown in the matrix)

2. **Config resolution order, project-local only.** Three layers, last wins:
   1. Built-in defaults (current `DefaultSpecs()`).
   2. Project file `parley-deck/agents.toml` (checked in; carries non-secret runtime defaults the team agrees on).
   3. User file `parley-deck/agents.local.toml` (gitignored; carries `path`, `model`, `profile`, `isolated_home`, secrets references).
   No machine-global file in this slice — explicit non-goal, keeps private paths out of the repo.

3. **`parley agents` subcommand.** Single new command with three modes, all read-only-by-default:
   - `parley agents list` — capability matrix (id, installed, version, headless ok, workspace-write ok, isolated-home, external-backend, timeout, model/profile shown as resolved value or `cli-default`).
   - `parley agents verify [--agent <id>]` — runs the full smoke per agent, exits non-zero on any verification failure (see step 4).
   - `parley agents add` — interactive scaffold that walks the user through every required field and appends to `agents.local.toml`; never writes credentials.

4. **Verification is behavioral, not just `--version`.** For each agent, in addition to `LookPath` + version probe:
   - **Headless echo probe:** send a fixed prompt instructing the agent to write exactly `parley-deck/ideas/_runtime-probe/<agent>.md` under the target repo and exit. Confirm the file exists, is non-empty, and contains the expected sentinel. This is the only honest check that headless invocation + workspace-write + isolated-home actually compose.
   - **Codex-specific Git probe:** run the four smoke commands from `docs/agent-runtime-configuration.md` (`git status`, ephemeral branch create/delete, `git hash-object -w --stdin`) inside `codex exec` against the real repo. If any fails due to sandbox/permission, retry the *exact same argv* through `--ask-for-approval on-failure` and surface the approval event. Do **not** offer an "alternate workdir" path without an explicit `--allow-alt-workdir` flag the user passes consciously.
   - Verification artifacts go under `parley-deck/ideas/_runtime-probe/` (gitignored except a `.gitkeep`) so probes do not pollute the audit trail.

5. **Codex profile baked in but overridable.** Default Codex `HeadlessArgs` become `["exec", "--cd", "{root}", "--sandbox", "workspace-write", "--ask-for-approval", "on-failure", "-"]` (current code uses `never`; that's the bug behind "silent relocations"). Document the change in `protocol-changelog.md`. Users can still override per-agent in `agents.local.toml` if their Codex build differs.

6. **Capability matrix shown at run start.** Every `parley run` / `parley round` first prints the resolved matrix and waits 2 s (or `--yes`) so the user sees what is about to be invoked, with the external-backend column visible — this is the consent point for hosted backends.

7. **Docs delta.** `docs/agent-runtime-configuration.md` becomes the user-facing reference; add an "Adding an agent" walkthrough that mirrors `parley agents add`, and a "Codex behavior contract" subsection that codifies the retry-on-failure + no-silent-relocation rule.

Success criteria (verifiable):
- `parley agents list` exits 0 with all four current agents shown, no invented model names.
- `parley agents verify` produces a probe artifact for each installed agent and exits non-zero if any required field is missing or any smoke fails.
- For Codex: forcing a sandbox denial (e.g. `chmod a-w` on a probe file) causes the runtime to issue an approval prompt and re-run the *exact* failed command, observable in the event log.
- New unit tests in `internal/agents/` cover config layering, placeholder expansion, and approval-retry plumbing.

## Concerns / open questions

- **Approval-retry plumbing for Codex.** `codex exec` with `--ask-for-approval on-failure` already retries internally; do we need any CLI-side retry logic, or only correct flags + an event observer? I lean toward: just set the flag, observe events, surface the approval to the human via HITL. Worth confirming against Codex's current event schema before locking the design.
- **Headless echo probe is intrusive.** It runs the real agent against the real repo and writes a real file. For hosted backends this also sends a token of repo context. Acceptable cost for "verify cooperation works", but the user must consent — should `parley agents verify` require `--yes` for any agent whose `external_backend != local`? I think yes.
- **`isolated_home` semantics differ per agent.** Gemini uses `GEMINI_CLI_HOME`; Hermes uses `HERMES_HOME`; Claude/Codex have no clean equivalent. I'd model it as `{enabled: bool, env: {VAR: path-template}}` rather than inventing a unified env var — simpler, no fiction.
- **`agents.toml` vs CLI flags.** Do we also want `--model`, `--profile`, `--timeout` flags on `parley run`? My instinct is no for this slice — config file only — to avoid two sources of truth. Mention as a follow-up.
- **External-backend disclosure granularity.** Free-text `disclosure` per agent is honest but ugly. Worth standardizing values (`hosted-anthropic`, `hosted-google`, `hosted-openai`, `local`, `unknown`) so the matrix can color-code? Probably yes, but bikeshed-y.
- **What "verify Git writes" means for non-Codex agents.** The Codex Git probe is specific. For others, the headless echo probe already exercises write capability; I don't think we need a separate Git probe per agent — but worth checking nobody disagrees.

## Risks

- **Scope creep into M6/auto-mode.** "Agent runtime" easily drifts into model selection policy, fallback ordering, consensus weighting. Keep this slice strictly to *describing and verifying* runtime; do not let it pick winners.
- **Breaking existing runs.** Switching Codex's default `--ask-for-approval` from `never` to `on-failure` changes runtime semantics for anyone already using the CLI headlessly. Mitigation: gate behind a clear changelog entry; provide `--ask-for-approval never` override in `agents.local.toml` for CI-style users.
- **Probe artifacts polluting `parley-deck/`.** If `_runtime-probe/` isn't gitignored cleanly, smoke runs leak into the audit trail and confuse consensus tooling. Mitigation: explicit `.gitignore` entry plus a sentinel marker the probe checks for, refusing to run outside the dedicated dir.
- **User confusion between `agents.toml` and `agents.local.toml`.** Two-file split is correct (team vs personal) but every two-file scheme breeds "why is my override ignored" issues. Mitigation: `parley agents list` shows the *source file* per resolved field.
- **Headless echo probe cost.** For hosted-backend agents, every `verify` costs a real model call. Default `parley agents verify` to discovery + version only; require `--full` to invoke the headless probe. Trade-off: weaker default verification, but cheap, and the full check is one flag away.
- **False sense of safety from "workspace-write".** Different agents implement workspace narrowing very differently (Claude's `--add-dir`, Gemini's `--approval-mode auto_edit`, etc.). The matrix should not claim equivalence — show the literal flag used, not a normalized "✓".
- **Codex retry path masks real failures.** If `on-failure` triggers approval prompts for every sandbox tick, the user habituates to clicking approve. Surface a counter ("3 approvals in this run") and consider promoting to a single warning rather than per-call prompts in a follow-up slice.
