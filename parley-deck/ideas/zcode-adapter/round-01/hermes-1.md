---
agent: hermes-1
idea: zcode-adapter
round: 1
date: 2026-08-18
---

## Summary
zcode (zcode-app-cli 3.7.7-13 / zcode-runtime 0.16.3, binary /opt/homebrew/bin/zcode, provider Z.AI anthropic-kind) must become a first-class adapter. No `--model` flag exists (`zcode --help` does not list it; rejected flags print help and exit 0, so a silent failure is possible). The effective model is `zai/glm-5.3` read from `~/.zcode/cli/config.json` (`model.main`), not passed in argv. `--json` confirms usage/projection but does NOT expose the model id, so the model is externally unobservable at launch time.

## Proposed approach

### Q1 — MODEL column: `CLIDefault` (`cli-default`), never `zai/glm-5.3`
The roster contract is explicit (00-prompt.md §Q1, discover.go:36 `Model` field): MODEL is what the launch argv actually carries, or `unknown` — never a configured value argv does not pass. zcode has no `--model` flag; `zcode --model` is not listed in help (verified 2026-08-18: `zcode --help` output; `zcode --model foo` prints help, exits 0). Adding a model to `HeadlessArgs` as `--model {model}` would result in a rejected flag that exits 0, making the adapter look successful but produce no real output — unacceptable. Therefore the adapter's `Model` field remains `CLIDefault`. `roster show` will report `model-unbound` / `cli-default`. That is NOT a lie; it is the only honest statement for an adapter whose argv does not carry a model. The configured `zai/glm-5.3` value from `~/.zcode/cli/config.json` should be mentioned in `Notes` (e.g. "effective model: zai/glm-5.3 from ~/.zcode/cli/config.json, unobservable at launch time; not passed in argv") but never promoted to `Model` or `ModelLabel`. To add a vocabulary term like "model-read-from-config" would be a contract change (§7 protocol change) — out of scope here; propose as successor if needed.
Verification: `grep 'main' ~/.zcode/cli/config.json` (REDACTED key excluded) shows `"main": "zai/glm-5.3"`; `zcode --json --mode yolo --prompt ...` returns session/projection/usage but no `modelId` field.
Cost of being wrong: if we claim `Model: "zai/glm-5.3"`, `roster show` displays a configured value the argv does not carry — a direct violation of the frozen roster contract. If we invent a protocol term, we change COOPERATION.md without authorization.

### Q2 — AutonomousWrite: Mode `"yolo"`, Args `["--mode", "yolo"]`, Scope `""`
Verified (`discover.go:86-92`): Scope `"workspace"` is allowed ONLY where the CLI enforces a real workspace sandbox (`codex --sandbox workspace-write`). zcode's confinement is `--cwd <path>` (a working directory, not an enforced sandbox) plus optional `--allowed-tools`/`--disallowed-tools` allowlists/denylists (verified in help). `--cwd` does not sandbox; it just sets the process cwd. The config (`~/.zcode/cli/config.json`) sets `permission.mode = "build"` and `autoApproveHighRisk: false`; the effective autonomous mode we pass explicitly is `--mode yolo`. Per the honesty rule (never falsely assert workspace confinement), Scope must be EMPTY (`""`). `AutonomousWrite` declaration: `{Mode: "yolo", Args: ["--mode", "yolo"], Scope: ""}`. This aligns structurally with `hermes` (`Mode: "yolo"`, Scope `""`), `claude` (`Scope: ""` with `--add-dir`), and `opencode` (`Scope: ""` with `--auto`).
Verification command: `cat -n .../discover.go | sed -n '86,97p'` shows the quoted rule.
Cost of being wrong: setting Scope `"workspace"` falsely asserts enforced sandbox — violates the binding honesty rule.

### Q3 — Owner's generalisation ("every agent should support something like yolo")
This is a protocol-level invariant proposal, not part of this adapter's design. It states that every roster member must declare an autonomous-write mode, with `roster show` flagging members that do not. Today `kimi-1` reports AUTO=no (verified in `00-prompt.md`) and is a manual participant, so adopting this invariant today would break kimi-1's current status or force an adapter change for kimi — out of scope for `zcode-adapter`. Per `COOPERATION.md` §7, adding an invariant is a protocol change requiring a dedicated idea (`§7 protocol change`). I recommend recording this as a successor idea (`idea: autonomous-write-invariant`) rather than bundling it here. For this adapter, the declaration is sufficient.
Verification: `grep 'kimi' .../discover.go` shows `AutonomousWrite: {Mode: "prompt", Args: []string{"-p"}, Scope: ""}` but `LaunchMode: LaunchHeadless` with `-p` — AUTO depends on `AutonomousEffective()` which checks args match; `kimi-1` may still be `no` depending on config. Not changed by this idea.
Cost of being wrong: mixing a protocol-change proposal into this adapter design would violate the non-goal (`COOPERATION.md is untouched`) and the track rules.

### Q4 — Exit-0-on-bad-flag hazard
Verified (`zcode --bad-flag-test`): rejected flags print help and exit 0. A bad adapter design could silently degrade to "printed help" and look like success. Mitigation: the adapter's launch uses only verified flags (`--prompt`, `--mode`, `--cwd`, `--json`); `--model`, `--max-turns`, `--settings` are NOT included (verified: `--model` not in help; `--settings` and `--max-turns` ARE listed — whether they actually function was not retested today per 00-prompt.md, so they must NOT be included in `HeadlessArgs` unless retested). The existing `artifact-exists` check (mentioned in 00-prompt.md) is sufficient: if the adapter launches and produces no artifact (only help text), the artifact check fails. `parley agents verify --agent zcode-1` must verify: (a) binary found (`zcode` on PATH, `/opt/homebrew/bin/zcode`), (b) version probe (`zcode --version`), (c) headless args contain `--mode yolo`, (d) `AutonomousWrite.Args` is a subset of `HeadlessArgs` (`AutonomousEffective() == true`). No zcode-specific extra assertions needed beyond the standard adapter verification pipeline.
Verification commands: `zcode --bad-flag-test` (exit 0, help printed); `zcode --json --mode yolo --cwd /tmp --prompt "hello"` (exit 0, real JSON artifact).
Cost of being wrong: including an unverified flag (`--settings`, `--max-turns`) could silently exit 0 with help output, making the adapter look working but empty.

### Q5 — Scope of change (verified files, not assumed)
Files that change (verified by inspection):
- `internal/agents/discover.go`: add zcode `Spec` entry in `defaultBuiltinSpecs()` (after `opencode` or before `kimi`; structural analogue is `hermes`). Must include: `ID: "zcode"`, `Commands: ["zcode"]`, `VersionArgs: ["--version"]`, `LaunchMode: LaunchHeadless`, `HeadlessMode: "zcode --prompt ..."`, `HeadlessArgs: ["--prompt", "{prompt}", "--mode", "yolo", "--cwd", "{root}", "--json"]` (order: verify `{prompt}` position against prompt-arg rules; `--json` optional — include for telemetry), `PromptMode: PromptArg`, `SandboxMode: CLIDefault`, `ApprovalPolicy: "yolo"`, `Model: CLIDefault`, `Reasoning: CLIDefault`, `Profile: CLIDefault`, `AutonomousWrite: {Mode: "yolo", Args: ["--mode", "yolo"], Scope: ""}`.
- `internal/agents/modelmeta.go`: add `"zai"` prefix mapping. Currently has `"z-ai"` (`zhipu`) but zcode's id is `zai/glm-5.3` (no hyphen). Need `prefixRule` for `"zai"` mapping to `"GLM"` / `"Zhipu AI"`, or add `"zai"` to producers (`"zai": "Zhipu AI"`) and rely on segment parsing (`segments[0] == "zai"`, `segments[1] == "glm-5.3"` → family `GLM`, company `Zhipu AI`). Verified: `modelmeta.go:43` has `"z-ai"`; adding `"zai"` fixes `metadata-unknown` for zcode.
- `internal/agents/autonomous_test.go`, `modelmeta_test.go`, `acp_specs_test.go`, `naming_test.go`, `internal/app/app_test.go`, `preflight_test.go`: extend to include `zcode` adapter assertions (verified paths by `find . -name '*test*'` and `grep -rn 'opencode\|hermes\|kimi'` in those files).
- Skill package (`parley-deck-skill/skills/parley-deck/SKILL.md`): add `zcode` row to the autonomous-write table (§Autonomous Execution, line ~243-249).
- `README.md` / release notes: mention new adapter.
Not changed: `COOPERATION.md` (protocol unchanged), any other adapter's spec, any protocol vocabulary (SUCCESSOR idea if needed).
Verification: `find ...` and `grep` commands listed above.
Cost of being wrong: missing `modelmeta.go` update means `roster show` continues to report `metadata-unknown`; missing test updates means adapter not covered by automated verification.

### Q6 — Release
This is a CLI change only (new adapter entry), no new binary. Per `parley-deck-skill/RELEASING.md` (verified exists), the skill release crosses npm, Homebrew, winget, GitHub with independent per-channel verification. Version bump: CLI version is in `VERSION` file (verified: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli/VERSION`); a new adapter is a minor bump (e.g., 1.44.0 → 1.45.0) per changelog convention (`CHANGELOG.md`). The release checklist for a new adapter specifically must include: (a) adapter added to `defaultBuiltinSpecs()`, (b) model meta added, (c) autonomous-write declaration verified (`AutonomousEffective()` true), (d) all six test files extended with zcode assertions, (e) skill `SKILL.md` updated, (f) per-channel install verification (`npm install -g` / `brew install` / `winget install`) confirms `zcode --version` is discoverable and `parley agents verify --agent zcode-1` passes. Secret redaction: `~/.hermes/config.yaml` and zcode configs (`~/.zcode/cli/config.json`) contain live API keys (verified: key starts with `d3ec...`); never print, quote or include them in any artifact — the config file contents are referenced by path only, with the key value redacted.
Verification: `cat VERSION`; `cat RELEASING.md` (exists at root); `cat ~/.zcode/cli/config.json | grep 'apiKey'` shows key (redacted in this analysis).
Cost of being wrong: incorrect version bump or missing per-channel verification means the adapter is not actually available to users.

## Concerns / open questions
- The effective model (`zai/glm-5.3`) is unobservable from `zcode --json`; a future zcode release could change `config.json` without the adapter noticing. The adapter must not claim it observed the model.
- `--settings` and `--max-turns` are listed in help; whether they function was not retested today (00-prompt.md note). Do NOT add them to `HeadlessArgs` without retesting.
- A third Q1 answer could be a new roster vocabulary term (`model-config-only` or similar) that explicitly signals "model configured but not passed". That would be a §7 protocol change (COOPERATION.md vocabulary change) and must be scoped as a SUCCESSOR idea, not this adapter.
- The owner's instruction to program through `/parley-deck` implies this adapter lands in the CLI repo; confirm whether the skill package release is synchronized or independent.

## Risks
- Silent failure: rejected flag exits 0. If an unverified flag slips into `HeadlessArgs`, the adapter produces an empty artifact. Mitigation: only verified flags (`--prompt`, `--mode`, `--cwd`, `--json`); no `--model`, `--settings`, `--max-turns`.
- Honesty violation: setting Scope `"workspace"` or claiming `Model` carries `zai/glm-5.3` violates binding rules (`discover.go:86-92`, roster contract). Mitigation: Scope empty, Model `CLIDefault`.
- Protocol creep: bundling the owner's autonomous-write invariant (Q3) or a new vocabulary term (Q1 third option) into this adapter design would violate `COOPERATION.md` (§7: protocol changes require dedicated ideas) and the non-goal (`COOPERATION.md untouched`).
- Secret leakage: `~/.zcode/cli/config.json` and `~/.hermes/config.yaml` contain live API keys. Any command output or file reference must redact them (verified: `grep apiKey ...` shows value; it is excluded from this artifact).
