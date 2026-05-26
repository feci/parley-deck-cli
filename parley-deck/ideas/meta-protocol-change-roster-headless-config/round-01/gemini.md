---
agent: gemini
idea: meta-protocol-change-roster-headless-config
round: 1
date: 2026-05-25
---

## Summary
I agree with the problem statement and the proposed solution. Standardizing headless agent configuration via a checked-in default file and a git-ignored local override is a robust pattern that will improve facilitator reliability and make project onboarding easier. This change decouples the facilitator's logic from the specific CLI flags of each agent, which is a necessary step for a healthy multi-agent ecosystem.

## Proposed approach
I endorse the core proposal of using `headless-agents.json` and `headless-agents.local.json`. I recommend we use JSON for its simplicity and ubiquity.

For the configuration structure itself, I propose a format that is expressive enough for different runtimes but still simple. My role is to consider cross-runtime implications, so I suggest we include support for environment variables and templatized arguments.

A proposed structure for an agent's configuration in `agents.json`:
```json
{
  "<agent-id>": {
    "cli": "path/to/executable",
    "env": {
      "API_KEY_NAME": "VALUE_FROM_FACILITATOR_ENV" 
    },
    "args": [
      "--headless",
      "--output-file",
      "%OUTPUT_FILEPATH%"
    ],
    "profiles": {
      "fast": {
        "model": "small-model",
        "args": ["--effort", "1"]
      },
      "deep": {
        "model": "large-model",
        "args": ["--effort", "5"]
      }
    },
    "defaultProfile": "deep"
  }
}
```
Key features of this approach:
- **`env` block:** Allows the facilitator to pass through necessary environment variables (like API keys) from its own environment. This avoids storing secrets in the config file.
- **`args` with placeholders:** Using placeholders like `%OUTPUT_FILEPATH%` allows the facilitator to inject dynamic values into the argument list in a structured way, accommodating different CLI syntaxes.
- **`profiles`:** This abstracts away the specific flags for model, speed, and thinking into named profiles. The facilitator can just request a profile (e.g., 'deep') without needing to know the underlying flags. This further reduces brittleness.

The facilitator's logic would be:
1. Load `headless-agents.json`.
2. Deep-merge `headless-agents.local.json` over it.
3. When invoking an agent, select the profile.
4. Construct the command line by combining the base `args` and the profile's `args`.
5. Replace placeholders like `%OUTPUT_FILEPATH%`.
6. Set the environment variables specified in the `env` block.
7. Execute the command.

## Concerns / open questions
1.  **Secret Management:** My proposal suggests referencing environment variable names from the facilitator's context rather than storing secret values. Is this sufficient? Should the protocol explicitly forbid storing secrets in `headless-agents.local.json`? I believe it should.
2.  **Merge Strategy:** The protocol should specify that the local file is deep-merged into the base file. This provides the most flexibility for local overrides.
3.  **Schema Definition:** Should we provide a formal JSON schema for this file? It would help with validation and documentation.
4.  **Placeholder Syntax:** Is `%PLACEHOLDER%` the right syntax? It's simple and unlikely to conflict with shell syntax. We should standardize it.

## Risks
1.  **Complexity Creep:** The biggest risk is making this configuration format too powerful or complex. We must strictly limit it to invocation data. It should not become a scripting language.
2.  **Configuration Fragmentation:** If agents have CLIs, environment variables, and now this config file, it might become confusing where to set a particular value. The protocol update should clearly state the order of precedence (e.g., CLI flags provided by the facilitator at runtime override the config file).
3.  **Stale Configuration:** The config files could become out of sync with the agent CLIs they are meant to configure. There is no easy way to automatically validate this, so it will rely on human maintenance.
