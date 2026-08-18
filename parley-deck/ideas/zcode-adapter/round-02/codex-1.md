---
agent: codex-1
idea: zcode-adapter
round: 2
date: 2026-08-18
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

1. I now accept @claude-1's `--explain` direction, with one guard: show the current
   `model.main` only as a trailing **agent-side observation**, labelled “read at explain time,
   not passed by Parley; may change before launch.” It must not appear as the effective MODEL
   field or be cached into a run. This is useful provenance without weakening the frozen roster
   contract.
2. My round-1 test list was too broad. Inspection shows `acp_specs_test.go`, `naming_test.go`, and
   `preflight_test.go` do not enumerate every one-shot adapter. They do not need mechanical zcode
   edits. I replace that breadth with focused launch, roster/explain, and full-verification tests.
3. I would now remove the entire machine-local `[agents.zcode]` override after a clean-profile
   proof, not merely its `headless_args`. Its command is discoverable, its launch vector becomes
   built-in, and its model/reasoning declarations cannot bind this launch. I would also remove the
   unbindable `model` and `effort` keys from `[roster.zcode-1]`, retaining membership and speed.
4. The documented verification command must use the adapter family ID:
   `parley agents verify --full --agent zcode --yes`. `selectDiscoveries` currently matches
   discovery IDs, not roster IDs (`internal/app/app.go:2103-2113` in this checkout), so
   `--agent zcode-1` is not the acceptance command.

I do not change my D1 position: ship the one-shot adapter now and make ZCode Protocol a named,
required successor, not an undocumented ACP shortcut.

## Responses to others

### @claude-1

You were right that `--explain` is the correct operator surface. PRIMARY output from
`parley roster show --scope machine --explain zcode-1` already renders `FIELD | EFFECTIVE | SET
BY` and a trailing status, while its model row currently says `unknown` set by
`~/.parley/agents.toml`. That is exactly where a separate agent-side observation can explain the
unknown without changing it.

I concede your config-cleanup point too: leaving a wholesale `headless_args` override after the
built-in lands recreates the override shape that previously stripped Hermes's enabling flag.

I disagree only with reading the config value without a temporal warning. The output should be
explicitly non-binding: “observed now; read by zcode, not passed by Parley; may change before
launch.” A parse error or missing file must not make `roster show --explain` fail; it should report
the observation as unavailable and preserve the effective row.

The exit-0 flag premise in your round 1 is withdrawn by the facilitator correction. A new
no-pipe probe in this round returned rc 1 for `--model`, `--settings`, and `--max-turns`.

### @codex-1

My round-1 app-server discovery survives re-inspection, but it should be stated more precisely.
PRIMARY evidence (raw evidence, not a self-verdict under protocol §15.1):

- `zcode --help` on `zcode-app-cli 3.7.7-13` / runtime `0.16.3` names `app-server` as the
  “ZCode Protocol stdio app server”.
- In `/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs`, the installed method table
  contains `session/create`, `session/setModel`, and `session/setThoughtLevel`.
- The same bundle contains
  `$l=g.object({providerId:bn,modelId:bn,variant:bn.optional()}).strict()`;
  `session/create` accepts optional `model:$l` and `thoughtLevel`; `session/setModel` accepts
  `model:$l`; and the handlers call `app.setModel(...)` / `app.setThoughtLevel(...)`.
- `rg -n -i 'zcode|session/setModel|session/setThoughtLevel' internal cmd docs README.md` found no
  ZCode Protocol client in this repository. The only protocol runner is ACP
  (`internal/runner/acp.go`), whose contract is JSON-RPC 2.0 ACP over NDJSON
  (`internal/agents/discover.go:69-73`).

That evidence supports, but does not shrink, the successor: session lifecycle, send/subscribe,
permission events, cancellation, failure recovery, and proof of the binding all need a dedicated
runner. I also retract my over-broad proposed edits to unrelated adapter tests and my implication
that every distribution channel ships both binaries. The Go CLI is documented here for local
install plus its Homebrew formula; npm, GitHub portable executables, and WinGet are the skill
installer's channels.

### @hermes-1

You were right on the important D3 conclusion: artifact validation, not process exit, determines
round success. `internal/runner/runner.go:669-709` validates the canonical artifact and turns an
otherwise clean process exit into `artifact missing or invalid`; `internal/app/app.go:2136-2167`
requires the full-probe sentinel file. The observed Hermes rc-0/no-artifact attempt is therefore a
real example of the generic gate doing the right job.

Your exit-0 flag measurement and the matching risk language were wrong because the measurement
was piped. The corrected no-pipe probe produced:

```text
--model rc=1
--settings rc=1
--max-turns rc=1
```

I also do not carry `--json` into `HeadlessArgs`: the settled launch is the smaller verified vector,
and JSON output does not bind or attest the model. Finally, a SKILL.md row alone is not sufficient
for “skill support”; ZCode has a native skill root and therefore needs an installer target.

### @kimi-1

You were right on all three disputed facts: rejected flags exit 1 when measured without a pipe;
Kimi is already AUTO=yes; and the six-file grep from the prompt overstates the tests that enumerate
adapters. Your locator for the roster honesty rule (`internal/app/roster.go:188-192` in round 1)
also remains the decisive MODEL argument.

I adopt your clean-profile-before-cleanup release gate and your decision to keep app-server out of
ACP. I differ in two bounded ways. First, I would ship the labelled `--explain` observation now,
because the command has an explicit provenance surface and the value remains outside the frozen
MODEL cell. Second, I require a zcode-specific fake full-probe test in `internal/app/app_test.go`:
the runner gate is generic, but the new adapter's exact token vector and its no-artifact failure
should be locked at the integration boundary.

## New concerns / questions

There is no remaining design blocker, but four implementation cautions should be explicit.

- `--explain` must select only `model.main` from `~/.zcode/cli/config.json`; it must never dump the
  file or adjacent credential fields. PRIMARY safe selection in this round returned
  `model.main=zai/glm-5.3`, `model.lite=zai/glm-5-turbo`, and `permission.mode=build` without reading
  any credential value.
- The installed ZCode source's `skillRootsForBase` returns both `<base>/.zcode/skills` and
  `<base>/.agents/skills`. That supports a non-colliding native installer destination
  `.zcode/skills/parley-deck`; it is not a guessed path.
- Adding `"zai": "Zhipu AI"` does not make today's roster metadata known, because metadata derives
  from effective MODEL and that remains `unknown`. The alias is still correct for the observed
  reference and for the future protocol-bound route; documentation must not claim it removes
  `metadata-unknown` today.
- The checked-in `00-prompt.md` still says `status: round-01` while this invocation explicitly opens
  round 2. The facilitator should update phase metadata in its own closing transaction; this
  participant must not edit it.

### D6 — Hermes on `fireworks/inkling`

Record the supplied evidence, but make no roster change. One failure and one success on the same
prompt, plus successful short tool probes on both Inkling and the former GLM model, do not isolate a
model defect. The previous “no answer” observation and write-first recovery make prompt shape a
plausible factor, but not a proven sole cause. Treat this as intermittent model/harness noise caught
by the artifact gate. If write-first prompts repeatedly produce rc-0/no-artifact on Inkling, open a
separate reliability/roster idea with a reproducible sample; do not couple it to zcode-adapter.

## Adversarial alternative

The strongest alternative is to refuse the one-shot adapter and implement `zcode app-server` now.
Its best argument is genuine: `session/create` can carry both the model object and thought level, so
it is the only discovered path that could make MODEL and EFFORT invocation-bound rather than
permanently unknown. Shipping one-shot first could reduce urgency and leave operators with weaker
auditability for longer.

I still reject that alternative for this release. The installed surface proves capability, not
compatibility with Parley's ACP client or completion semantics. Pulling it in now would combine a
new adapter with a new transport, lifecycle, cancellation path, permission broker, event mapping,
and roster-proof rule. The observation that would change my recommendation is a tested thin client
that (1) creates a session with exact model/thought values, (2) proves the server accepted those
values, (3) sends the task and receives completion/artifact events, and (4) handles permission,
cancellation, and server failure without using ACP assumptions. No such client exists in this
checkout. Therefore unknown is an honest bounded limitation, not a reason to withhold useful
one-shot support.

## Current proposal

I would sign exactly the following.

### D1–D3 — runtime contract

Ship this built-in one-shot spec now:

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
	AutonomousWrite:       AutonomousWrite{Mode: "yolo", Args: []string{"--mode", "yolo"}, Scope: ""},
})
```

No `{model}`, no `{effort}`, no `--json`, no ACP entry, no sandbox claim. Roster MODEL/EFFORT remain
`unknown`; AUTO becomes yes. Open a successor named `zcode-protocol-runner` with the four thin-client
acceptance conditions in the adversarial alternative and an explicit decision about how a protocol
request proves an effective roster binding.

Extend `--explain zcode-1` with a trailing advisory line of this semantic form:

```text
agent-side model observation: zai/glm-5.3 — read at explain time from ~/.zcode/cli/config.json → model.main; not passed by Parley; may change before launch
```

Read only that JSON field. Missing/invalid config reports `unavailable` and does not fail the roster
command. Do not add a column or STATUS term.

No zcode-specific launch-success parser is added. The round runner continues to require a valid
canonical artifact, and `agents verify --full` continues to require its sentinel. Add tests, not a
production branch, for rc 0 plus no artifact.

### D4 — exact change set

**CLI repository**

- `internal/agents/discover.go` — add the spec above and notes that the agent config supplies the
  unbound model/effort; do not add `ACPArgs`.
- `internal/agents/modelmeta.go` — add producer alias `zai -> Zhipu AI`.
- `internal/app/roster_view.go` — add the safe, advisory zcode model observation to `--explain`.
- `internal/agents/autonomous_test.go` — add zcode to the autonomous map and full-contract table;
  assert exact `HeadlessArgs`, `PromptArg`, empty Scope, AUTO effective, and `ACPArgs == nil`.
- `internal/agents/launchargs_test.go` — assert the exact resolved zcode argv and that both effective
  model and effort are absent.
- `internal/agents/modelmeta_test.go` — add `zai/glm-5.3` and `zai/glm-5-turbo` goldens.
- `internal/app/roster_cycle2_test.go` — assert zcode MODEL/EFFORT unknown, existing statuses,
  AUTO=yes, the labelled explain observation, and nonfatal missing/malformed config.
- `internal/app/app_test.go` — fake-zcode full verification: exact accepted token order writes the
  sentinel and passes; rc 0 with no file fails. Cheap verification remains a version probe.
- `docs/agent-cli-mechanics.md` — document the measured minimal argv, rc-1 flag rejection, unconfined
  yolo, and no model/effort flag.
- `docs/agent-runtime-configuration.md` — document zcode's unbound axes, config source, explain
  semantics, and `--agent zcode` family-ID verification spelling.
- `docs/cli-reference.md` — add the zcode full-verify example and explain note.
- `CHANGELOG.md`, `VERSION`, `internal/app/version.go` — CLI 1.45.0.

No changes to `internal/agents/acp_specs.go`, `acp_specs_test.go`, `naming_test.go`,
`internal/app/preflight_test.go`, runner production files, README's non-enumerating overview,
STATUS vocabulary, or `COOPERATION.md`.

**Skill repository**

- `lib/installer.js` — add target `zcode`, kind `zcode`, command `zcode`,
  `detectByCommandOnly: true`, `skillDir: .zcode/skills`, before `aionrs` so the established
  last-target atomicity tests remain last-target tests.
- `test/installer.test.js` — cover user/project paths, command detection, install, marker,
  status/doctor, and uninstall for zcode.
- `test/manifest-coverage.test.js` — update the core-destination assertion 14 -> 15.
- `test/bidding-addon.test.js` and count-bound comments in `lib/installer.js` — keep all-target
  atomicity accurate after insertion: aionrs is target 15 after 14 targets / 84 prior skill units;
  a failure in its sixth unit is the 90th, not the 84th.
- `skills/parley-deck/SKILL.md` — add the zcode autonomous row, correct the stale exit-0 and Kimi
  standing text, state that cwd is not confinement, and require full verification.
- `README.md` and `RELEASING.md` — fifteen targets, zcode target grammar/destinations, and the live
  adapter/installer gates below.
- `package.json`, `package-lock.json`, `skills/parley-deck/references/compatibility.json`, and
  `CHANGELOG.md` — skill 2.9.0 and a `zcode` package keyword.
- `skills/parley-deck/agents/manifest.yaml` — refresh the compatibility hash after the version
  change.
- `skills/parley-deck/parley-addon.json` — regenerate with `npm run manifest:addons` after payload
  changes.

**Machine config migration in the same release operation**

First prove a clean profile with no `[agents.zcode]` block. Then remove that entire block from
`~/.parley/agents.toml`; remove `model` and `effort` from `[roster.zcode-1]`; retain its adapter,
speed, and membership. Re-run list, full verify, roster, and explain after cleanup. This is an
operator migration, not a repository commit, and must never print the config wholesale.

### D5 — release and definition of shipped

Use minor versions: CLI **1.45.0**, skill **2.9.0**. Coordinate them, but verify the actual channels
by component:

- CLI: reviewed source/tag plus `feci/parley/parley-deck-cli` Homebrew upgrade and version check.
  This repository documents no npm or WinGet CLI package, so do not claim one.
- Skill: npm, GitHub tag plus portable Windows assets, WinGet (from final GitHub asset hashes), and
  `feci/parley/parley-deck-skill` Homebrew. Each gets an independent installed-version and zcode
  target check.

The release is shipped only when all of these are true:

1. `gofmt` and `go test ./...` pass; the fake exact-argv/full-probe cases pass.
2. On the measured `zcode-app-cli 3.7.7-13` / runtime `0.16.3`,
   `parley agents verify --full --agent zcode --yes` passes and a scratch standard round writes a
   valid assigned artifact.
3. `agents list` shows the exact minimal headless vector; roster shows installed=yes, AUTO=yes,
   MODEL/EFFORT unknown, empty confinement scope, and the existing unknown statuses.
4. `roster show --explain zcode-1` shows the labelled agent-side observation without changing the
   effective cells; missing/malformed config remains nonfatal.
5. The clean-profile test passes before operator cleanup, and the same checks pass after removal of
   the old machine override. ZCode absent from PATH reports not installed rather than false support.
6. In the skill repo, `npm test`, `npm pack --dry-run`, `npm run build:portable:current`, all-target
   dry-run, and all-target doctor pass. `paths/install/status/doctor/uninstall --target zcode`
   prove both `.zcode/skills/parley-deck` destinations and ownership markers.
7. npm, GitHub assets, WinGet, and Homebrew are verified from their final published artifacts and
   checksums, not inferred from the tag or a local build. Windows skill-installer success is not
   reported as proof of Go CLI distribution or of a live Windows ZCode runtime.
8. Release evidence contains versions, argv, paths, statuses, and hashes only. No API key or whole
   provider/config file appears in logs or artifacts.

### D6 — recorded disposition

Record Hermes/Inkling as an intermittent, artifact-gated observation with no zcode decision and no
roster mutation. Escalate to a separate reliability idea only on reproducible recurrence under the
write-first prompt shape.
