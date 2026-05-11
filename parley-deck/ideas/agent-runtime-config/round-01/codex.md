---
agent: codex
idea: agent-runtime-config
round: 1
date: 2026-05-11
---

## Summary

The slice should introduce an explicit, local-first agent runtime configuration layer before M6 auto/consensus work. The current code has useful defaults in `internal/agents/discover.go`, but those defaults are hardcoded and cannot represent the user's required Codex profile, per-agent timeout policy, model/reasoning defaults, isolated-home needs, or external-backend disclosure. The design should keep discovery conservative: record only what can be configured or proven, and show `cli-default` where the CLI cannot prove model or thinking options.

## Proposed approach

Add a local configuration model loaded in this order:

1. explicit CLI flags for the current run;
2. `PARLEY_HEADLESS_AGENT_CONFIG` when set;
3. `parley-deck/meta/headless-agents.local.json` when present;
4. built-in adapter defaults;
5. discovered CLI defaults.

The local JSON file should be treated as machine-local config and should not be required in git. It should support defaults and per-agent entries:

```json
{
  "defaults": {
    "timeoutMs": 1800000,
    "speed": "balanced",
    "externalBackendDisclosure": "repo-context-without-secrets"
  },
  "agents": {
    "codex": {
      "cli": "codex",
      "headlessArgs": ["exec", "--cd", "{root}", "--sandbox", "workspace-write", "--ask-for-approval", "on-failure", "-"],
      "promptMode": "stdin",
      "sandboxMode": "workspace-write",
      "approvalPolicy": "on-failure",
      "model": "cli-default",
      "reasoning": "cli-default",
      "timeoutMs": 1800000
    }
  }
}
```

Implementation should add an `internal/config` package or equivalent small loader, then thread an effective `agents.Spec` or `agents.RuntimeConfig` into discovery and runner launch. `parley agents discover` should print a capability matrix with CLI path, installed status, version/probe result, headless mode, write mode, model/reasoning fields, timeout, isolated-home requirement, and warnings.

For Codex specifically, the CLI should prefer the user-requested `workspace-write` + `on-failure` profile when the configured CLI supports it. If the local Codex command cannot be probed or does not support the requested approval flag, the CLI should surface a warning rather than silently downgrading to `never` or a broad bypass mode.

Add a `parley agents check-git [--agent codex]` or fold the check into `agents discover --verify-writes` to run the non-destructive Git write smoke sequence in the target repository:

```sh
git status
git branch tmp-parley-git-test
git branch -D tmp-parley-git-test
printf test | git hash-object -w --stdin
```

If a command fails because of sandbox or permission restrictions, the recommended behavior is to retry the exact failed command through the runtime approval path. The CLI should not automatically move to another work directory as a workaround.

Keep implementation small:

- parse and merge config;
- update built-in specs through config;
- update discovery output and tests;
- update runner prompt's effective launch config from real values instead of hardcoded `cli-default`;
- add docs and examples for Codex, Claude, Gemini, and Hermes.

## Concerns / open questions

- The current `codex` binary on this machine is on PATH but `codex exec --help` fails because its vendored executable is missing. The design must distinguish "path exists" from "headless launch is operational".
- We should not commit private local config by default. A sample file can live in docs, while `headless-agents.local.json` should be local-only or explicitly ignored.
- Some CLIs may expose model lists or thinking controls only through authenticated commands. Discovery should avoid commands that spend model tokens or require network calls unless the user asks for deeper verification.
- The product requirement says GitHub/GitLab transport support matters, but this slice should only make runtime config transport-aware enough to not block future transport automation.

## Risks

- Over-designing the config schema could delay the next functional slice. Keep fields close to current `agents.Spec` plus the new runtime/safety fields.
- Hardcoding Codex `on-failure` without probing support could break users on older Codex CLIs.
- Displaying too much configuration may overwhelm users. The default matrix should be concise, with verbose JSON output available for debugging.
- If local config includes credential paths or home-directory details, accidental commits become a privacy risk. The implementation should document and enforce local-only defaults.
