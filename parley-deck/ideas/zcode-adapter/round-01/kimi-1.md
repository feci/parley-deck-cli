---
agent: kimi-1
idea: zcode-adapter
round: 1
date: 2026-08-18
---

## Summary

zcode can be a full built-in adapter today with the launch shape already proven —
`zcode --prompt "{prompt}" --mode yolo --cwd "{root}"` — but its MODEL column must stay
`unknown` (`model-unbound`), and that is the honest answer, not a gap to paper over. I probed
every conceivable argv/env mechanism for passing a model on zcode-runtime 0.16.3 and all are
dead: `--model` is rejected (exit 1), `--settings` and `--max-turns` are **listed in `--help`
but rejected by the parser** (exit 1), there is no `ZCODE_MODEL`-style env knob in the 13 MB
runtime bundle, and `--json` does not report the model id. The model lives only in
`~/.zcode/cli/config.json → model.main` (`zai/glm-5.3`). A roster cell carrying that value would
be exactly the "declaration wearing the effective cell" the frozen v1 contract forbids
(`internal/app/roster.go:188-192`).

Two probe findings correct stale machine notes and matter beyond Q1: (a) a rejected zcode flag
**exits 1, not 0**, on this version — the documented exit-0 hazard does not reproduce;
(b) `kimi-1` reports **AUTO=yes** today (the 00-prompt's premise for Q3 is stale; the
`~/.parley/agents.toml` note dates to parley 1.38, and v1.39.0 promoted kimi to a full adapter).

Design: new `zcode` spec in `internal/agents/discover.go` with
`AutonomousWrite{Mode: "yolo", Args: ["--mode", "yolo"], Scope: ""}` (Scope empty per the
honesty rule — `--cwd` is a working directory, not a sandbox), `"zai": "Zhipu AI"` added to
`modelmeta.go` producers, two test files extended, the skill's autonomous-write table gains a
zcode row, released as CLI **1.45.0** + skill **2.9.0**. No protocol change. The owner's
"every roster member must support yolo" rule and a config-read MODEL cell are both scoped as
successors.

## Proposed approach

### The spec (internal/agents/discover.go, appended to `defaultBuiltinSpecs`, :194-383)

```go
withBuiltinSources(Spec{
	ID:                    "zcode",
	Commands:              []string{"zcode"},
	VersionArgs:           []string{"--version"},
	LaunchMode:            LaunchHeadless,
	HeadlessMode:          "zcode --prompt ... --mode yolo",
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
	Telemetry:             "--json prints sessionId/traceId/usage/projection; no model id",
	Notes:                 "Z.AI (Zhipu AI) GLM coding CLI. No model/effort flag exists: --model is rejected, and --settings/--max-turns are listed in --help but rejected by the parser (all exit 1; probed 2026-08-18 on zcode-runtime 0.16.3). The model comes from ~/.zcode/cli/config.json model.main. config.json sets permission.mode=build, so --mode yolo must be explicit in argv.",
	// Scope EMPTY: --mode yolo is a permission grant plus --cwd; no enforced sandbox (honesty rule).
	AutonomousWrite: AutonomousWrite{Mode: "yolo", Args: []string{"--mode", "yolo"}, Scope: ""},
}),
```

No ACP catalog entry (`internal/agents/acp_specs.go:27-50`): `zcode app-server` speaks the
"ZCode Protocol" (`zcode --help`), not ACP, and nothing else claims ACP support. Headless only.

### Q1 — MODEL reports `unknown`, and that is honest. There is no third mechanism.

**Mechanism.** `EffectiveModel()` (`internal/agents/launchargs.go:97-108`) returns the value
after `--model`/`-m` in the resolved argv, else not-ok; `roster show` then prints `unknown` +
`model-unbound` (`internal/app/roster.go:369-379`). With `Model: CLIDefault` and no `{model}`
placeholder in `HeadlessArgs`, zcode reports `model-unbound`, `effort-unknown`,
`metadata-unknown` — permanently, within roster schema v1.

**What I verified (all PRIMARY, 2026-08-18, zcode-app-cli 3.7.7-13 / zcode-runtime 0.16.3):**

- `zcode --model zai/glm-5.3 --prompt … --mode yolo --cwd <scratch> --json` →
  `Unknown option '--model'`, help text, **exit 1**.
- `zcode --settings <file> --prompt …` → `Unknown option '--settings'`, help, **exit 1** —
  probed with two different settings files; the flag is *listed* in `--help` and rejected by
  the parser anyway. The 2026-08-17 agents.toml note was right about this; the 00-prompt left
  it open. The `--settings`-file third answer is **dead**.
- `zcode --max-turns 1 --prompt …` → `Unknown option '--max-turns'`, help, **exit 1** (also
  listed-but-rejected).
- `zcode --mode bogus --prompt …` → `Unsupported --mode value`, **exit 1**.
- Baseline `zcode --prompt "…" --mode yolo --cwd <scratch> --json` → exit 0, JSON with
  sessionId/traceId/turnId/response/usage/projection — **no model id** (re-verified).
- No env knob: `grep` of the real runtime (`vendor/zcode.cjs`, 13 MB, spawned per
  `bin/zcode.js:1239`) finds `ZCODE_API_KEY`, `ZCODE_BASE_URL`, `ZCODE_STORAGE_DIR`, etc. but
  **no model/effort variable**; `ZCODE_HOME` is used only for telemetry device-id persistence.
- `parley agents list` (run today): the zcode row shows MODEL `zai/glm-5.3` — the declared
  value from the machine config — alongside AUTO=no, BACKEND=unknown. So the declared model is
  already visible on the inventory surface; the roster contract rightly refuses to promote it
  to the *effective* column.
- `~/.zcode/cli/config.json` (dumped with keys redacted): `model.main = "zai/glm-5.3"`,
  `model.lite = "zai/glm-5-turbo"`, `permission.mode = "build"` — so the documented
  "yolo is default for --prompt" must not be relied on; `--mode yolo` stays explicit in argv.
- `metadata-unknown` today is a **consequence of** `MODEL=unknown`: `roster.go:386-390` derives
  metadata from `row.Model`, and `DeriveModelMeta("unknown")` matches nothing
  (`internal/agents/modelmeta.go:72-114`). The 00-prompt's attribution to the missing `zai`
  producer is incomplete — adding the producer does **not** change today's row. I still propose
  adding it (below): the only `DeriveModelMeta` call sites are the roster views
  (`roster.go:386`, `internal/app/roster_view.go:125`), so the day an *effective* `zai/…`
  reference reaches a row — a future schema v2, or a zcode that grows a model flag — the
  registry must already resolve it. Note also `modelmeta_test.go:62-71`: it skips
  `Model: CLIDefault` specs, but had the built-in shipped a literal `zai/glm-5.3` default,
  that test is the tripwire that forces the producer entry.

**Why `unknown` is acceptable for a first-class adapter.** The v1 contract exists to stop lies,
not to force every CLI into one shape (`roster.go:188-192`: "never a declaration wearing the
effective cell"). An `unknown` axis is already a normal, visible roster state — PRIMARY,
`parley roster show` today: codex-1, hermes-1 and kimi-1 all report `effort-unknown`; hermes-1
additionally `metadata-unknown`. First-class means: discovered, version-probed
(`zcode --version` → `zcode-app-cli 3.7.7-13 / zcode-runtime 0.16.3`, exit 0), launchable,
AUTO=yes, artifact-checked. None of that needs a model cell. The spec's `Notes` (surfaced by
`agents list`) tells the operator where the model actually lives and how to change it.

**The actual third answer is a successor, say so plainly.** A schema-v2 roster that carries a
config-read model *with provenance* (e.g. cell `zai/glm-5.3` + new STATUS `model-config-read`)
would be informative without lying — but the STATUS vocabulary is closed and the columns are
frozen (`roster.go:178-184`), so that is a contract change, proposed as successor idea
"roster schema v2: config-read model provenance", not here.

**Display-name divergence (flag, don't fix silently).** `RenderDisplayName`
(`internal/agents/naming.go:188-207`) derives the composite name from `spec.Model`, and
`resolveRoster` copies roster-entry fields onto the spec (`roster.go:354-364`). Today's machine
config has `[roster.zcode-1] model = "zai/glm-5.3"`, so the display name would read
`zcode_glm-5.3_max` while the MODEL column says `unknown`. Display is explicitly non-contract
(`roster.go:205-218`), but it re-creates the declaration-vs-effective gap one channel over.
Release checklist below drops `model`/`effort` from that roster entry (display then falls back
to a `cli-default` form — honest); keeping them is an operator choice to make with eyes open.

**Cost of being wrong.** If a model flag appears in a later zcode, the cost is one release
note's worth of follow-up: add `--model {model}` to `HeadlessArgs` and the cell fills itself —
the contract needs no change. If we lied now instead, the cost would be the exact regression
class (`model-drift`, hidden swaps) the roster was rebuilt to expose. Cheap to be honest,
expensive to be caught.

### Q2 — `AutonomousWrite{Mode: "yolo", Args: ["--mode", "yolo"], Scope: ""}`

The honesty rule (`internal/agents/discover.go:86-92`, read before answering) sets Scope to
`"workspace"` **only** where the CLI enforces a real sandbox; grant-flag + cwd CLIs stay empty
(precedent: hermes `discover.go:315`, kimi `:348`). zcode's confinement story, checked one
claim at a time:

- `--cwd <path>` is a working directory. Nothing in `zcode --help` or the config schema
  describes fs enforcement; it is the same class as hermes's cwd and claude's `--add-dir`.
- `--allowed-tools` / `--disallowed-tools` are tool allow/deny lists for headless prompts
  (listed in `--help`). I did **not** probe their runtime effect — flagged unverified — but
  even working perfectly they filter *tools*, not filesystem reach. Not a sandbox.
- `autoApproveHighRisk: false`, `allowMediumRiskInAuto: false` in `config.json` are approval
  gates the per-invocation `--mode yolo` sets aside (mode is the permission mode). A grant,
  not confinement.

So Scope is `""`. `Confined()` returns false; the skill's fail-closed rule (treat
un-demonstrable confinement as unset) is satisfied by construction. AUTO=yes holds:
`AutonomousEffective()` (`discover.go:138-146`) resolves the argv and finds both declared
tokens present; the fail-closed regression shape (override strips the flag) is already pinned
by `autonomous_test.go:105-126`, and zcode joins that test's contract table (Q5).

**Cost of being wrong.** Claiming `"workspace"` would be the single reviewed CRITICAL this
field's comment exists to prevent; claiming nothing costs only a `Confined()=false` reading,
which every non-codex adapter already reports.

### Q3 — The owner's generalisation is a deck rule, therefore a successor

"vsetci agenti ktory su v rostery by nieco take mali podporovat" is a statement about the
*roster*, not about this adapter — as an invariant ("an active roster member must be
autonomous-write effective") it changes roster semantics, and COOPERATION.md is untouched by
this idea. Scope it as successor idea "roster AUTO invariant" (natural shape: `roster show`
already exposes AUTO fail-closed; preflight/run warns or refuses when an *active* member reads
AUTO=no).

One correction that changes the successor's urgency: the 00-prompt says kimi-1 "currently
reports AUTO=no and is a manual participant". PRIMARY, `parley roster show --scope machine`
on parley 1.44.0, 2026-08-18: **kimi-1 reports AUTO=yes** (MODEL `kimi-code/k3`, STATUS
`effort-unknown` only). The machine note asserting AUTO=no/manual dates to parley 1.38;
v1.39.0 shipped the full kimi adapter (CHANGELOG:414-440). Once zcode ships, **all five**
machine roster members are AUTO=yes — the invariant would flag nobody today. It is still worth
ratifying for future ACP-only or manual members, but it is not load-bearing for this adapter.

### Q4 — The exit-0 hazard does not reproduce on 0.16.3; existing checks suffice

Probes above: rejected flag → help on stdout, **exit 1** (three flag shapes); invalid value →
exit 1. The "prints help and exits 0" claim in the 2026-08-17 agents.toml note is stale for
this version. Even if a future zcode regresses to exit 0, two existing layers catch a
degraded launch, and I read both rather than assuming:

- **Run path.** After process exit the runner stats and phase-validates the artifact
  (`internal/runner/runner.go:668-676`); missing/invalid → `agent.failed` with
  "artifact missing or invalid" (`:709`). Artifact-wins (`:699-706`) rescues only an ordinary
  exit error when a **valid** artifact exists — a help-printing launch produces none, so it
  fails regardless of exit code.
- **Verify path.** `agents verify --full` → `runHeadlessProbe` (`internal/app/app.go:2126-2157`):
  process error, artifact-exists, and sentinel-prefix all must pass. A rejected-flag launch
  fails at least two of the three.

So: no post-launch stdout heuristic ("looks like help") — redundant with the artifact check
and brittle across zcode versions. The one genuine gap I found is not zcode-specific:
`selectDiscoveries` (`app.go:2100-2113`) matches family IDs only, so
`agents verify --agent zcode-1` reports "unknown agent" — I verified `--agent hermes-1` fails
identically today while `--agent hermes` works. Fix within scope: document `--agent zcode`
(family ID) in cli-reference and the release note; optionally make `selectDiscoveries`
roster-aware (small, benefits every adapter; fine to defer).

**Cost of being wrong.** If some untested zcode path still exits 0 on rejection, the artifact
check is the backstop and the round fails visibly — the failure mode is a flagged round, never
a silently accepted empty artifact.

### Q5 — Exact change set (each entry verified, not assumed)

CLI repo:

- `internal/agents/discover.go` — append the spec above to `defaultBuiltinSpecs` (:194-383).
  The existing `[agents.zcode]` block in `~/.parley/agents.toml` then becomes an override on
  the built-in spec by ID (`internal/config/runtime.go:759` `applyOverride`), so the machine
  keeps working through the upgrade.
- `internal/agents/modelmeta.go` — `producers` (:35-44) += `"zai": "Zhipu AI"` (namespace form
  of the existing `"z-ai"` entry; Z.AI is Zhipu's brand, and `glm-5.3` also hits the `glm`
  prefix rule at :63, so family/company agree).
- `internal/agents/autonomous_test.go` — `wantMode` (:16-26) += `"zcode": "yolo"`;
  `TestPromotedAdaptersFullContract` cases (:60-67) += `{"zcode", []string{"--mode", "yolo"},
  []string{"--prompt", "{prompt}", "--mode", "yolo", "--cwd", "{root}"}}` — that test already
  locks Scope-empty, PromptArg, LaunchHeadless and AutonomousEffective for promoted adapters.
- `internal/agents/modelmeta_test.go` — cases += `zai/glm-5.3` and `zai/glm-5-turbo` →
  GLM / Zhipu AI / known (table pattern at :20).
- `internal/agents/acp_specs_test.go` — **no change**: assertions are required-subset and
  no-duplicates (:62-94); zcode enters no ACP list.
- `internal/agents/naming_test.go` — no change required (`RenderDisplayName` is spec-generic);
  an optional composite case (`zcode_glm-5.3_max`) documents the display derivation.
- `internal/app/app_test.go`, `internal/app/preflight_test.go` — **no change**: read both;
  the agents tests pin fake-codex behavior (:292-357), preflight pins protocol freshness, no
  adapter enumeration. Confirm with `go test ./...`.
- `docs/cli-reference.md` — add zcode to the verify example (:49) and a STATUS-table note that
  zcode's `model-unbound` is expected, not a fault. Light touch.
- `README.md` — no change: :33 already reads "Claude, Codex, Gemini, GLM, and more".
- `CHANGELOG.md` — new v1.45.0 entry (content requirements in Q6).

Skill package (sibling repo `parley-deck-skill`):

- `skills/parley-deck/SKILL.md` — autonomous-write table (:243-251) += row:
  `| zcode | `--mode yolo` (+ `--cwd <deck>` to scope it). No model flag exists; the model comes from zcode's own config. |`
- Skill `CHANGELOG.md`, version 2.8.0 → 2.9.0 ("tracks parley-deck-cli 1.45.0", matching the
  2.8.0-tracks-1.44.0 convention); `npm run manifest:addons` if the payload hash moves.

Operator config at release (not a repo change, part of the checklist):

- `~/.parley/agents.toml`: delete the custom `headless_args` from `[agents.zcode]` (the
  built-in now carries the identical argv; a wholesale override is the hermes landmine),
  delete `model`/`effort` from `[roster.zcode-1]` (Q1 display divergence), and remove the
  stale notes documenting 1.44.0 behavior (including the AUTO=no kimi note corrected in Q3).
- This deck's `parley-deck/agents.toml`: whether zcode-1 joins *this* roster is the owner's
  membership decision, not part of the adapter.

### Q6 — Release: CLI 1.45.0 + skill 2.9.0, with new-adapter-specific gates

Version: **minor**, CLI 1.44.0 → **1.45.0**. Precedent: v1.39.0 ("kimi and opencode are full
adapters", CHANGELOG:414) was a minor bump for exactly this change class. Skill 2.8.0 → 2.9.0.

CLI gates, in order:

1. `go test ./...` green — must include the extended `autonomous_test.go` and
   `modelmeta_test.go`.
2. Build; `parley agents list` shows zcode with the resolved headless argv and AUTO=yes;
   `parley agents verify --agent zcode --full --yes` passes the real headless probe
   (sentinel file written via `--cwd`).
3. `parley roster show --scope machine`: zcode-1 AUTO=yes; STATUS
   `model-unbound,effort-unknown,metadata-unknown` — expected and documented, not a failure.
4. VERSION=1.45.0; CHANGELOG entry **must** record: the probed launch shape; that `--settings`
   and `--max-turns` are listed-but-rejected (exit 1) on 0.16.3 — so nobody later "fixes" the
   model pin by adding them; and that `model-unbound` is contractual. Windows exes follow the
   existing `dist/parley-vX.Y.Z-windows-{x64,arm64}.exe` pattern.

Skill channel checklist (RELEASING.md), each channel verified independently:

1. Preflight: `npm test` (requires python3), `npm pack --dry-run`,
   `npm run build:portable:current`, install `--dry-run` + `doctor --json`.
2. npm: `npm publish --access public`; verify `npx -y parley-deck-skill@latest install`.
3. GitHub: tag `v2.9.0`, push; confirm release-portable workflow passed and both `.exe`
   assets are attached (`gh release view`).
4. WinGet: manifest hashes taken **from the GitHub release assets** (never a local build),
   `winget validate` + `winget install --manifest` on Windows.
5. Homebrew tap (`feci/homebrew-parley`): formula `url`/`sha256` from the release tarball;
   `brew style`, `brew audit --strict --online`, `brew upgrade`, `brew test`, `--version`.

New-adapter-specific additions:

- **Clean-profile smoke**: verify on an account *without* the old `[agents.zcode]` config
  block, proving the built-in stands alone; only then do the operator-config cleanup (Q5).
- **Post-cleanup AUTO check**: `roster show` again after deleting the custom `headless_args` —
  AUTO must stay yes (the fail-closed rule flips it to no if an override ever drops
  `--mode yolo`; this is the exact hermes regression, so check it explicitly).
- **Re-probe instruction**: the CHANGELOG and spec Notes pin the rejected-flag behavior to
  zcode-runtime 0.16.3 and instruct re-running the flag probes after any zcode upgrade.

## Concerns / open questions

- **Permanent unknowns.** MODEL/EFFORT/metadata for zcode-1 stay `unknown` until a roster
  schema v2 (successor proposed in Q1) or a zcode model flag exists. If another participant
  argues the cell should carry the config-read value now, my answer is: that is the frozen
  contract's definition of a lie; change the contract in the open, not the cell in the dark.
- **`--allowed-tools`/`--disallowed-tools` unprobed.** Listed in `--help`; runtime effect not
  verified. Irrelevant to Scope (tool filters are not fs confinement), but if a reviewer wants
  them in `HeadlessArgs` as hardening, that needs a probe first — and a rejected-flag probe at
  that, given what `--settings` turned out to be.
- **Roster-aware verify.** `agents verify --agent <roster-id>` fails for every adapter
  (verified hermes-1). Worth a small fix here or a named successor; documenting the family-ID
  form is the minimum.
- **Successors named, not started**: (1) roster schema v2 config-read model provenance
  (contract change); (2) roster AUTO invariant (§7 protocol change). Both are real, neither is
  this idea.
- **Why does `--help` list flags the parser rejects?** Unanswered — upstream zcode bug or
  intentional reservation. It is the strongest reason to re-probe after every zcode upgrade
  rather than trusting help text.

## Risks

- **zcode upgrades change flag behavior silently** (a listed-but-rejected flag is an upstream
  documentation bug; a future version may also change rejection exit codes). Mitigation: the
  artifact check is version-independent; CHANGELOG/Notes pin the probed version; release note
  instructs re-probing.
- **The model can drift under us**: zcode runs whatever `model.main` says, and the user (or a
  zcode update) can change it with no parley-visible signal. The roster will not show it —
  by design, today. This is the strongest argument for the schema-v2 successor; until then the
  spec Notes say where to look.
- **Operator config not cleaned up at release**: the stale `headless_args` override is harmless
  today (identical to the built-in) but is exactly the shape that stripped `--yolo` from hermes
  once. The release checklist makes its removal an explicit, verified step.
- **Display/column divergence** if `[roster.zcode-1] model` is kept: §2 tables and digests
  would read `zcode_glm-5.3_max` while `roster show` says `unknown`. Mitigated by the release
  checklist dropping those fields; flagged here so the choice is conscious.
- **Cost of a wrong Scope call**: asymmetric, as argued in Q2 — `""` is the only defensible
  value; anything else is the reviewed CRITICAL this field exists to prevent.
