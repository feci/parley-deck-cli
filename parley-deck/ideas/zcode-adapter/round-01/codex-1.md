---
agent: codex-1
idea: zcode-adapter
round: 1
date: 2026-08-18
---

## Summary

Ship `zcode` as a normal one-shot headless adapter now, with the effective launch

```text
zcode --prompt <prompt> --mode yolo --cwd <root>
```

and with `AutonomousWrite{Mode: "yolo", Args: ["--mode", "yolo"], Scope: ""}`. The roster must report `MODEL=unknown` and `EFFORT=unknown`. A configured `zai/glm-5.3` value is useful operator intent, but this launch does not carry it, so putting it in the effective MODEL cell would violate the frozen roster contract.

That limitation does not make the adapter second-class. First-class here means built-in discovery, a tested headless prompt path, an effective autonomous-write declaration, normal round/artifact validation, a full behavioral verification, runtime documentation, and a native skill-installer destination. It does not mean claiming a vendor capability that the one-shot CLI lacks.

There is a genuine third route for a later design: `zcode app-server` exposes ZCode Protocol methods including `session/create`, `session/setModel`, and `session/setThoughtLevel`; the `session/create` schema carries a model object with `providerId`, `modelId`, and optional `variant`, plus a thought level. That could bind model and effort without editing the user's config. It is not the current one-shot argv contract and it is not ACP, so adopting it requires a dedicated runner/client and a successor decision about how the roster proves non-argv bindings. It should not be smuggled into this additive adapter.

All ZCode observations below are PRIMARY from this machine on 2026-08-18. No credential-bearing config was printed. The safe config inspection selected only model identifiers and permission booleans.

## Proposed approach

### Q1 — MODEL and EFFORT: keep both `unknown`

**Mechanism.** Add this semantic contract to `DefaultSpecs()` in `internal/agents/discover.go` (ordinary Go formatting omitted here only for compactness):

```go
withBuiltinSources(Spec{
    ID:                    "zcode",
    Commands:              []string{"zcode"},
    VersionArgs:           []string{"--version"},
    LaunchMode:            LaunchHeadless,
    HeadlessMode:          "zcode --prompt ... --mode yolo --cwd ...",
    HeadlessArgs:          []string{"--prompt", "{prompt}", "--mode", "yolo", "--cwd", "{root}"},
    InteractivePromptMode: InteractivePromptNone,
    InteractiveInvoke:     InteractiveInvokePrintOnly,
    InteractivePollMS:     DefaultInteractivePollMS,
    PromptMode:            PromptArg,
    SandboxMode:           CLIDefault,
    ApprovalPolicy:        "yolo",
    Model:                 CLIDefault,
    Reasoning:             CLIDefault,
    Profile:               CLIDefault,
    Speed:                 DefaultSpeed,
    TimeoutMS:             DefaultTimeoutMS,
    ExternalBackend:       ExternalHosted,
    Telemetry:             "headless JSON has session/usage data but no externally reported model id",
    Notes:                 "ZCode one-shot headless mode; tested with zcode-app-cli 3.7.7-13 / runtime 0.16.3; model and effort are vendor-config defaults, not invocation-bound",
    AutonomousWrite:       AutonomousWrite{Mode: "yolo", Args: []string{"--mode", "yolo"}, Scope: ""},
})
```

There is deliberately no `{model}` or `{effort}` placeholder. `Spec.EffectiveModel()` and `Spec.EffectiveEffort()` therefore return false under `internal/agents/launchargs.go:93-121`; `resolveRoster` consequently renders `unknown` and the existing `model-unbound`, `effort-unknown`, and `metadata-unknown` statuses under `internal/app/roster.go:369-390`. No new STATUS word and no `COOPERATION.md` edit are needed.

Do add `"zai": "Zhipu AI"` beside `"z-ai"` in `internal/agents/modelmeta.go:35-44`, with a `zai/glm-5.3 -> GLM / Zhipu AI` golden case. This makes the metadata utility correct for a ZCode model reference wherever an actual reference is available; it must not be used to back-fill the effective roster cell. Since that cell is `unknown`, its metadata remains unknown in this design.

**What I verified.** `zcode --version` returned `zcode-app-cli 3.7.7-13` and runtime `0.16.3`; `zcode doctor` reported runtime `0.16.3`, process `zcode-cli`, Node `v26.7.0`, and `darwin/arm64`. `zcode --help | rg -- '--model|--effort|--reasoning|thought|/model'` found only the interactive `/model` command, not an invocation flag. A proper argv probe,

```text
zcode --prompt /help --mode yolo --cwd /private/tmp --model zai/glm-5.3
```

returned rc 1 and `Unknown option '--model'`. Safe `jq` selection from `~/.zcode/cli/config.json` verified `model.main=zai/glm-5.3`, `model.lite=zai/glm-5-turbo`, and `permission.mode=build`, but those are configuration values rather than launch arguments. The prompt's PRIMARY reconnaissance also established that ordinary headless JSON has session/trace/turn/usage/projection data but no model identifier.

The current roster code is unambiguous: `internal/agents/launchargs.go:93-107` says a spec with no model flag passes no model and must not claim its declared value; `internal/app/roster.go:369-378` writes `unknown`. Reading the config and displaying `zai/glm-5.3` in MODEL would therefore be a lie even if the config happened to be stable at inspection time.

The third route is real but out of scope. `zcode --help` names `app-server` as the “ZCode Protocol stdio app server”. In `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`, the `zr` method table contains `session/create`, `session/setModel`, and `session/setThoughtLevel`; the model schema (`$l`) is `{providerId, modelId, variant?}`, and the implementations call `app.setModel` / `app.setThoughtLevel`. The transport's requests and lifecycle are ZCode-specific, not the JSON-RPC 2.0 ACP implementation described by `Spec.ACPArgs` at `internal/agents/discover.go:69-73`. The slash commands `/model` and `/effort` also exist, but I did not verify a safe way to compose either command and the work prompt into one one-shot turn, so they are not a binding mechanism in this design.

**Cost of being wrong.** If Parley displays the configured model as effective, an operator can change the ZCode config between roster inspection and launch—or ZCode can select a fallback—and the audit record will falsely say the invocation pinned GLM 5.3. If Parley reports `unknown` and ZCode later adds a stable one-shot flag, the cost is only temporary loss of observability until the adapter is upgraded. If the app-server route is pulled into this slice without a dedicated, tested client, the cost is a much larger runner change whose session events, permissions, cancellation, and binding proof are not covered by the current contract.

### Q2 — autonomous write is `yolo`, with empty scope

**Mechanism.** Pass `--mode yolo` explicitly in `HeadlessArgs`, declare the identical two-token `AutonomousWrite.Args`, set `ApprovalPolicy: "yolo"`, and leave `Scope: ""`. Keep `--cwd {root}` because it selects the workspace, but do not describe it as confinement. Do not pass `--allowed-tools` or `--disallowed-tools` in the default adapter. The latter is a useful optional denylist, but neither is a filesystem sandbox, and the former is rejected by the installed runtime despite appearing in top-level help.

**What I verified.** The binding honesty rule is at `internal/agents/discover.go:86-92`: only a demonstrated workspace sandbox earns `Scope: "workspace"`; cwd plus an approval bypass earns an empty scope. The installed runtime source in `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs` has a `PermissionService.checkPermission` branch that allows normal tool calls in `mode === "yolo"` with reason `Yolo mode bypasses permission prompts`. Explicitly disallowed tools and tools requiring user interaction are checked before that branch. The `autoApproveHighRisk` check occurs later in build-mode logic, so the safe-config values `autoApproveHighRisk=false` and `allowMediumRiskInAuto=false` do not turn yolo into confinement.

The accepted launch

```text
zcode --prompt /help --mode yolo --cwd /private/tmp
```

returned rc 0 and the slash-command listing. On the same binary, properly tokenized probes showed `--allowed-tools Read` returning rc 1 / unknown option, while `--disallowed-tools Bash` returned rc 0 and executed `/help`. Source function `resolveDefaultSkillRoots` is unrelated to permissions; nothing in the permission source or help identifies an OS-level workspace sandbox.

**Cost of being wrong.** Claiming `Scope: "workspace"` would tell the facilitator that a full-filesystem permission bypass is confined when it is not. Leaving scope empty is conservative: it may cause policy to exclude ZCode where a sandbox is mandatory, but it never grants false assurance. Treat yolo as external-hosted, unrestricted execution and require the same explicit user confirmation Parley uses for other hosted, unconfined agents.

### Q3 — satisfy the owner for ZCode now; ratify a deck-wide invariant separately

**Mechanism.** In this idea, make ZCode's built-in adapter autonomous and extend the existing code-level invariant in `internal/agents/autonomous_test.go`. Do not change the protocol, quorum rules, roster STATUS vocabulary, or manual-agent semantics. Open a successor §7 idea for the broader rule: decide whether every *active roster member* must have effective autonomous argv, whether manual/interactive participants can carry a waiver, how preflight treats them, and whether a new status is needed. That successor should use the already computed `AutonomousEffective()` value, not merely a declaration.

**What I verified.** `internal/agents/autonomous_test.go:11-49` already asserts autonomous modes for the promoted built-ins, and `internal/app/roster.go:392-395` fails AUTO closed when the declared args are absent. On 2026-08-18, `parley roster show --all` reported `kimi-1 ... AUTO=yes`; its adapter source at `internal/agents/discover.go:318-348` declares `-p` as the headless auto-approve path. Thus the prompt's statement that Kimi currently reports AUTO=no is stale on this checkout. That does not settle the general policy: custom and manual roster participants can still be non-autonomous, and `COOPERATION.md` currently permits manual cooperation.

**Cost of being wrong.** Treating the owner's sentence as an unratified global hard stop in this slice could remove manual participants, alter quorum, and change protocol behavior. Treating it only as a ZCode concern forever could leave a future AUTO=no roster member able to stall unattended rounds. The bounded answer is adapter enforcement now and a visible successor protocol decision, not silent semantics in either direction.

### Q4 — artifact validation is the success oracle, not process exit

**Mechanism.** Do not add a ZCode-only post-launch parser to the round runner. The generic runner already requires a canonical artifact and validates its frontmatter and required sections. `internal/runner/runner.go:627-714` sets success only after `validateArtifactForPhase`; a missing/invalid artifact becomes `artifact missing or invalid` even when the process returned 0. That is exactly the correct defense against “printed help”.

Keep cheap `parley agents verify --agent zcode` as the documented installation/version check. Require `parley agents verify --full --agent zcode --yes` for adapter qualification and every release: `internal/app/app.go:2108-2163` launches the actual headless argv, requires the probe file, and checks its unique first-line sentinel. Static help parsing is not adequate because this ZCode launcher advertises options the current runtime rejects.

Add two full-verify integration cases in `internal/app/app_test.go`:

1. A fake `zcode` that returns a version, then prints top-level help and exits 0 for the task invocation, writing no sentinel file. `agents verify --full --agent zcode --yes` must return nonzero.
2. A fake `zcode` that accepts exactly `--prompt <prompt> --mode yolo --cwd <root>`, extracts the assigned probe path, and writes the sentinel. Full verify must pass. Assert the exact token order so a future spec cannot silently append an unsupported option.

The existing runner hardening tests already cover “no artifact means failure”; the ZCode integration test proves that the new adapter reaches that gate.

**What I verified.** The prompt records a rejected-flag/rc-0 build. The installed 3.7.7-13 launcher differs today: properly tokenized `--settings`, `--max-turns`, `--model`, `--permission-mode`, and `--allowed-tools` probes each returned rc 1 plus help, while `--disallowed-tools` was accepted. This discrepancy is itself the reason not to rely on exit code or the help surface. I inspected both `finalizeExecResult` in `internal/runner/runner.go:627-714` and `runHeadlessProbe` in `internal/app/app.go:2108-2163`; both fail when their required artifact is absent or invalid.

**Cost of being wrong.** Exit-code-only success can admit an empty participant and let later phases mistake absence for agreement. A brittle ZCode-specific stdout matcher can fail on localization or vendor wording while still missing another silent failure mode. Artifact plus sentinel/schema validation tests the capability Parley actually needs and works for both rc-0 and rc-1 rejection behavior.

### Q5 — exact implementation scope

**CLI repository (`parley-deck-cli`).** Change:

- `internal/agents/discover.go`: add the ZCode spec above. Do not enable home isolation: the verified login/config/session state is under the user's ZCode home, and this design has not established a minimum safe subset to copy.
- `internal/agents/modelmeta.go`: add the `zai` producer alias only; do not use config inspection to populate the roster.
- `internal/agents/autonomous_test.go`: add mode `yolo` and a full-contract case locking `HeadlessArgs`, `PromptArg`, empty scope, and `AutonomousEffective()==true`.
- `internal/agents/launchargs_test.go`: assert ZCode has no model/effort placeholder, reports neither effective value, and resolves to the exact verified argv.
- `internal/agents/modelmeta_test.go`: add the `zai/glm-5.3` family/company golden case.
- `internal/agents/acp_specs_test.go`: add ZCode to the `DefaultSpecs` inventory and assert `ACPArgs == nil`. Do not edit `internal/agents/acp_specs.go`; `zcode app-server` is ZCode Protocol, not ACP.
- `internal/agents/naming_test.go`: add ZCode to synthetic compose/parse/display cases using a configured `glm-5.3` model; document that display-name composition does not make the roster MODEL effective.
- `internal/app/roster_test.go`: map a `zcode-1` row with configured `zai/glm-5.3` and assert MODEL/EFFORT remain `unknown`, the existing unbound/unknown statuses remain, and AUTO is yes.
- `internal/app/app_test.go`: extend discovery/version coverage and add the two full-probe cases from Q4.
- `internal/app/preflight_test.go`: include a found ZCode discovery in the exact-selected-set test and assert selecting it does not get dropped merely because model/effort are unknown.
- `README.md`, `docs/agent-cli-mechanics.md`, and `docs/agent-runtime-configuration.md`: document tested versions, exact argv, unconfined yolo, unknown effective model/effort, and the requirement to use full verify for behavioral qualification. `docs/cli-reference.md` needs only a ZCode full-verify example; the existing roster vocabulary is unchanged.
- `CHANGELOG.md`, `VERSION`, and `internal/app/version.go`: release the additive adapter as CLI `1.45.0`.

No runner production file needs a ZCode branch; the generic artifact gate is sufficient. No status table, `COOPERATION.md`, protocol core, or ACP catalog entry changes.

**Skill repository (`parley-deck-skill`).** Change:

- `lib/installer.js`: add target `zcode`, kind `zcode`, command `zcode`, `detectByCommandOnly: true`, and `skillDir: .zcode/skills`. Insert it before `aionrs` so the fleet-failure fixtures keep `aionrs` last; update comments from 14 targets / 13 earlier targets / 78 earlier skill units to 15 / 14 / 84 (six skill units per target in the current default plan).
- `test/installer.test.js`: lock user destination `~/.zcode/skills/parley-deck`, project destination `<project>/.zcode/skills/parley-deck`, command detection, install, marker target, status/doctor, and uninstall.
- `test/manifest-coverage.test.js`: update the core-destination count from 14 to 15.
- `test/bidding-addon.test.js`: include `.zcode` in all-target atomicity coverage and update the last-target commentary from the fourteenth target after thirteen to the fifteenth after fourteen, and from 78 earlier units to 84; update any count-bound assertions on the same basis.
- `skills/parley-deck/SKILL.md`: add `zcode | --mode yolo; --cwd selects a directory but is not a sandbox` to the autonomous table, plus the MODEL/EFFORT unknown warning and full-verify command.
- `README.md`: change fourteen to fifteen named runtimes, add ZCode to the prose and `--target` grammar, and document the native user/project destinations.
- `RELEASING.md`: add ZCode target correctness and the adapter-specific live checks in Q6.
- `package.json` and `package-lock.json`: bump `2.8.0` to `2.9.0`; add the `zcode` keyword.
- `skills/parley-deck/references/compatibility.json`: bump `skillVersion` to `2.9.0`.
- `skills/parley-deck/agents/manifest.yaml`: refresh the compatibility hash after that file changes.
- `skills/parley-deck/parley-addon.json`: regenerate with `npm run manifest:addons` after the SKILL/compatibility/agent-manifest changes.
- `CHANGELOG.md`: record the new native target, launch contract, honest unknowns, and tested ZCode version.

**What I verified.** `zcode skills list --json` found user skills under `~/.agents/skills`. More decisively, the installed runtime's `resolveDefaultSkillRoots`/`skillRootsForBase` source returns both `<base>/.zcode/skills` and `<base>/.agents/skills` for user and every project directory up to the worktree root. Therefore `.zcode/skills` is a native, non-colliding target rather than a guessed path. The current skill `lib/installer.js:18-119` has 14 targets and no ZCode entry; `README.md:206-238` names 14 and omits ZCode; `skills/parley-deck/SKILL.md:239-250` omits it from the autonomous table. `test/manifest-coverage.test.js:314` hardcodes 14 destinations. The package version is 2.8.0 in `package.json`, `package-lock.json`, and compatibility JSON, and the core add-on manifest hashes the changing payload.

**Cost of being wrong.** Documentation without an installer target leaves ZCode users unable to install the skill through the promised native workflow. A guessed target could create an ignored tree; the runtime source removes that uncertainty. Missing count/manifest/version updates can make an all-target release fail, or worse, make packaged bytes and health metadata disagree. Adding app-server to ACP would create a transport mismatch and hang/fail protocol negotiation.

### Q6 — coordinated minor releases and per-channel proof

**Mechanism and versions.** This is additive user-visible functionality, so bump the Go CLI from `1.44.0` to `1.45.0` and the skill package from `2.8.0` to `2.9.0`. Publish the CLI release and the skill release as a coordinated pair; do not advertise ZCode support until the compatible CLI is available. The repositories have independent versions, so record both in the release notes and compatibility evidence rather than pretending they are one artifact.

The release gate is:

1. **CLI source/build:** run `gofmt`, `go test ./...`, build `./cmd/parley`, install into a clean prefix, and assert `parley version` is 1.45.0. Run fake-cli regression tests for valid argv, rc-0 help/no artifact, rc-1 rejection, and missing ZCode.
2. **Live adapter matrix:** on the measured ZCode 3.7.7-13/runtime 0.16.3, run cheap verify and `parley agents verify --full --agent zcode --yes`; run a scratch round that writes and validates its assigned artifact. Assert effective argv contains `--mode yolo` and `--cwd`, AUTO=yes, Scope is not workspace, and roster MODEL/EFFORT are `unknown` even when the roster/config declares `zai/glm-5.3`. Repeat with ZCode absent and require `not installed` rather than a false discovery.
3. **Skill package preflight:** run the commands already required by `parley-deck-skill/RELEASING.md`: `npm test`, `npm pack --dry-run`, `npm run build:portable:current`, installer all-target dry-run, and all-target doctor. Add `paths/install/status/doctor/uninstall --target zcode` against a temporary home and confirm the exact `.zcode/skills/parley-deck` marker/payload.
4. **npm (skill 2.9.0):** after publish, check the registry's exact version, install with a fresh npm cache, run `npx parley-deck-skill@2.9.0 install --target zcode`, then doctor and inspect the marker version. Do not infer npm health from the git tag.
5. **GitHub (both releases):** tag the exact reviewed commits, verify release notes name CLI 1.45.0 plus skill 2.9.0, and verify checksums/digests of every uploaded artifact. For the skill release, require the documented Windows x64/arm64 portable assets and a green `release-portable.yml`; download each asset and run `--version`/`paths --target zcode` on its target architecture. The CLI repository currently exposes no `.github` release workflow or `RELEASING.md` (`rg --files .github` and `find . -maxdepth 2 -iname '*release*'` found none), so the CLI tag/binary publication procedure must be written and proven before claiming a GitHub binary release.
6. **Homebrew (independent):** update and audit both tap formulas that are actually relevant: `parley-deck-cli` (documented in the CLI README) to 1.45.0 and `parley-deck-skill` (documented in the skill RELEASING guide) to 2.9.0, each from its final release checksum. In a clean prefix run `brew upgrade`, version checks, skill `install/doctor --target zcode`, CLI cheap/full verify, and the scratch artifact test.
7. **WinGet (skill 2.9.0):** derive hashes from the final GitHub skill-installer assets, not local builds; run `winget validate` and a clean Windows install, then `parley-deck-skill --version`, `paths --target zcode`, install, and doctor. Also run the actual Go CLI 1.45.0 plus ZCode adapter test on Windows before the release notes claim Windows cooperation. The present WinGet instructions package `Feci.ParleyDeckSkill`, not the Go `parley` binary, so passing WinGet proves the skill-installer target only; it must not be reported as proof that CLI 1.45.0 was delivered by WinGet.
8. **Post-release clean-machine audit:** independently install each channel, record command output, component version, artifact checksum, OS/architecture, ZCode version, full-probe result, roster JSON, and scratch artifact path. Run `parley version --all` and `parley-deck-skill status --target zcode --json` to detect cross-component drift. Check release logs/artifacts for accidental config or API-key disclosure.

**What I verified.** CLI `VERSION` and `internal/app/version.go` are both 1.44.0; the CLI README states semantic versioning and documents `feci/parley/parley-deck-cli` Homebrew installation. The skill package is 2.8.0. Its `RELEASING.md` requires npm, GitHub portable Windows assets, WinGet manifests sourced from final GitHub hashes, and a separate Homebrew tap formula, with independent verification commands. It also confirms that the WinGet assets are the skill installer. These are distinct distribution proofs, not one release copied four times.

**Cost of being wrong.** A patch bump understates a new supported runtime and weakens compatibility signaling. Publishing the skill first can advertise a native adapter that the installed CLI does not yet contain; publishing the CLI without the skill leaves ZCode unable to receive the cooperation instructions through the supported installer. Treating npm/WinGet skill-installer success as proof of Go CLI delivery creates a false cross-platform claim. Failing the live full probe can ship an adapter that is installed and version-detectable but incapable of doing a round.

## Concerns / open questions

- The compatibility floor is proven only for `zcode-app-cli 3.7.7-13` / runtime `0.16.3` on darwin/arm64. Linux and Windows, and earlier ZCode builds, need the release matrix above. Until then, documentation should say “tested with”, not “requires” or “supports all versions”.
- Top-level ZCode help and the actual runtime parser disagree for `--settings`, `--max-turns`, `--permission-mode`, and `--allowed-tools` on this installation. The adapter wisely uses none of them, but a vendor update can still change the accepted launch surface. Full verify is the compatibility gate.
- `zcode app-server` is promising because it can bind model and thought level. Before a successor adopts it, obtain or derive a stable protocol specification; test session creation, event subscription, permission requests, cancellation, artifact completion, model fallback, and error recovery; then decide whether the roster contract may accept a protocol message as effective-launch proof. Do not place its command in `ACPArgs`.
- ZCode's one-shot JSON not identifying the model means post-hoc model attestation is unavailable in this slice. “Unknown” is the honest result, not an implementation defect to conceal.
- The release requirement spans two repositories and two kinds of product. The skill has npm/GitHub-portable/WinGet/Homebrew procedures; the Go CLI currently documents Homebrew but has no repository-local release automation. Release notes must state exactly which channel delivered which component.

## Risks

- **Unconfined writes:** `--mode yolo` bypasses prompts and `--cwd` is not a sandbox. Mitigation: empty autonomous scope, explicit hosted/unconfined consent, no claim of workspace confinement, and optional operator-side deny rules only after live verification.
- **Configured-model drift:** ZCode can change or fall back from `model.main` without the launch carrying a pin. Mitigation: MODEL/EFFORT stay `unknown`; never snapshot config into the effective cells.
- **Flag-surface drift:** help text may be a compatibility facade rather than the parser's true option set. Mitigation: the minimal verified argv, exact-argv unit test, full sentinel probe, and artifact-schema validation.
- **False verify green:** the cheap verify checks version only. Mitigation: call it “cheap” everywhere and make `--full --yes` mandatory in adapter qualification and release gates.
- **Transport conflation:** ZCode Protocol resembles an agent transport but is not current ACP. Mitigation: no `ACPArgs`, no app-server code in this slice, and a separately ratified successor.
- **Skill destination/manifest drift:** a fifteenth native target touches fleet-wide atomicity assumptions and hashed payload metadata. Mitigation: path/detection/marker tests, count updates, manifest regeneration, all-target dry-run, doctor, and uninstall tests.
- **Partial release:** users can receive docs/skills without CLI support, or CLI support without the ZCode skill. Mitigation: coordinated versioned releases and independent post-install evidence for each channel and component.
- **Secret exposure during diagnostics:** ZCode and adjacent provider configs contain live credentials. Mitigation: tests use fake CLIs or safe key selection; release evidence records versions, argv, status, and hashes only—never whole config files or environments.
