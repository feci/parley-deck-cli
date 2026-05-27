---
agent: claude
idea: antigravity-agent-migration
review-round: 1
date: 2026-05-27
---

## Findings

### MAJOR

1. **`installedAgentIDs` quietly tightens default-selection beyond gemini exclusion**
   `internal/app/app.go:1937` — The original function included any agent where
   `result.Found == true`. The new version adds `LaunchMode == LaunchHeadless`,
   `len(HeadlessArgs) > 0`, and `ID != "gemini"`. The gemini gate is the stated
   goal; the other two conditions are a silent behavioral change that could drop
   previously auto-selected agents from defaults. Both `internal/app/app.go:1937`
   and `internal/tui/app.go:1022` hardcode the string `"gemini"` — consider a
   `Deprecated bool` field on the Spec so future deprecations don't require
   scattering magic strings.

2. **Hermes `--yolo` approval policy is a significant security posture change**
   `internal/agents/discover.go:192` — Adding `--yolo` to hermes HeadlessArgs
   and setting `approval_policy = "yolo,accept-hooks"` goes beyond updating the
   model to `xai/grok-4.3`. `--yolo` bypasses all permission prompts in hermes,
   which is a materially different safety stance than the previous
   `--accept-hooks`. This change is not called out in the FINAL.md plan or the
   CHANGELOG entry and is not reflected in `docs/agent-runtime-configuration.md`.

### MINOR

1. **Claude HeadlessArgs duplicates Model/Reasoning fields**
   `internal/agents/discover.go:125` — `"--model", "opus", "--effort", "max"` is
   baked into HeadlessArgs while `Model: "opus"` and `Reasoning: "max"` are set
   as separate fields. Two sources of truth that can drift independently.

2. **Runtime configuration docs not updated for hermes changes**
   `docs/agent-runtime-configuration.md:26-28` — The gemini example was replaced
   with the agy section, but the doc doesn't reflect hermes gaining `--yolo`,
   `xai/grok-4.3`, or the changed `approval_policy`. Readers will see stale
   hermes defaults if they look at the docs.

3. **`agy` shares `~/.gemini/` tree but has no `isolate_home`**
   `internal/agents/discover.go:140-161` — The agy plugin path is
   `~/.gemini/config/plugins/parley-deck`. Gemini and hermes both use
   `isolate_home = true` to prevent config state leakage between runs. If agy
   writes any runtime state under `~/.gemini/`, concurrent or sequential runs
   could interfere. Confirm this is intentional.

### NIT

1. **Whitespace-only reformatting in unrelated map literal**
   `internal/agents/acp_specs_test.go:7-22` — Trailing-space cleanup in the
   `wantBinaries` map is unrelated to the migration.

2. **`Telemetry: "unknown"` on agy**
   `internal/agents/discover.go:159` — Acceptable placeholder, but worth a
   follow-up once agy documents its output/telemetry modes.
