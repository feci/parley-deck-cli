---
agent: codex-1
idea: composite-agent-naming-and-roster-reinit
round: 1
date: 2026-07-18
---

## Summary

I support composite names as self-documenting **display names**, not as protocol IDs. The stable ID must remain the identity used by `participants:`, artifact paths, event keys, presets, and signoffs; model and effort are mutable launch properties and therefore cannot safely be part of that identity. I support an explicit, rerunnable roster wizard and mandatory autonomous artifact-write profiles, but both need fail-closed workspace confinement and immutable per-run profile snapshots.

The current Go model makes this separation urgent: `agents.Spec.ID` is used both as a runtime/vendor discriminator and as the artifact identity. `runner.selectedAgents`, `runAgent`, `BuildRoundPrompt`, `consensus.missingRoundArtifacts`, `driver.BuildRoundDigest`, and the signoff validator all assume that one string is stable. Meanwhile `runner.cleanParticipantEnv` and `isolatedAgentHome` switch on vendor strings such as `claude` and `hermes`. Those are different concepts and should not share a field.

## Proposed approach

### A. Stable base ID plus a composite display name

My position is option (b): retain a stable base ID such as `claude-1` and add the composite as `DisplayName`. The composite must not be an alias accepted anywhere identity-bearing. It is presentation and an immutable launch-profile snapshot, not a lookup key.

The Go split should live in `internal/agents`:

```go
type Spec struct {
    ID          string // stable protocol ID: claude-1
    AdapterID   string // launcher/discovery family: claude
    DisplayName string // claude-opus48-max
    Model       string // authoritative, unsanitized selected model/id
    Reasoning   string // authoritative selected effort/reasoning value
    // existing runtime fields...
    AutonomousWrite AutonomousWriteMode
}

type AutonomousWriteMode struct {
    Mode             string   // vendor-specific semantic name
    Args             []string // launch fragments, empty when implicit
    Implicit         bool
    Scope            string   // must be "workspace"
    ScopeEnforcement string   // cli-sandbox | outer-sandbox
}
```

Keep `Spec.ID` so the many artifact and event call sites remain keyed by stable identity. Change vendor-specific branches and catalogs to use `AdapterID`: notably `defaultBuiltinSpecs`/`mergeACPCatalog` in `internal/agents/discover.go` and `acp_specs.go`, plus `cleanParticipantEnv`, `isolatedAgentHome`, and Hermes environment handling in `internal/runner/runner.go`. `Discovery` can continue embedding `Spec`; matrix/TUI output should show both `ID` and `DisplayName`.

`LoadAgentSpecs` in `internal/config/runtime.go` should instantiate roster members by stable ID and apply an adapter template separately. The new config shape should be explicit:

```toml
[agents.codex-1]
adapter = "codex"
display_name = "codex-gpt55-xhigh"
model = "GPT-5.5"
reasoning = "xhigh"
```

Old `[agents.codex]`/`[agents.claude]` sections should be treated as legacy adapter defaults during migration. They must not cause historical `codex.md` artifacts to be renamed to `codex-1.md`. Resolution of an old open idea whose participant is literally `codex` should preserve `Spec.ID == "codex"` while borrowing the `codex` adapter; closed history is never rewritten.

The risk of making the composite the ID is not cosmetic. Changing `codex-gpt55-xhigh` to `codex-gpt56-max` during an idea would make `runner.selectedAgents` fail to match the participant, direct the next artifact to a new filename, make `consensus.missingRoundArtifacts` report the old participant missing, and make an otherwise valid new signoff an unknown participant. It also splits run events and TUI state across two identities. A migration or alias cannot repair that without making quorum ambiguous. With the proposed split, the base ID never changes; an open idea also retains the profile it started with.

At Phase 0, add a flat immutable snapshot such as:

```yaml
participant-profiles: [claude-1=claude-opus48-max, codex-1=codex-gpt55-xhigh]
```

and add exact structured profiles (`id`, `adapter`, `display_name`, exact `model`, exact `reasoning`, autonomous-write profile) to `runmanifest.Manifest`. Resume must prefer that snapshot over newly reinitialized defaults. Current `run.json` stores only `Participants`; the existing `run.created` runtime event is useful evidence but should not be the only replay contract.

#### Sanitization, parsing, and collisions

The display grammar is deliberately narrower than the user's hard character set:

```text
<agent-token>-<model-token>-<effort-token>[-<instance>]
agent-token, model-token := [a-z0-9]+
effort-token := low|medium|high|xhigh|max|ultracode|clidefault
instance := integer >= 2
whole name := [a-z0-9_-]+ (generated names use no underscore)
```

Generate tokens from user-visible, operator-confirmed labels, never by guessing from a vendor's raw opaque model ID. For each agent/model token: trim, lowercase ASCII, then delete ASCII spaces and punctuation (`.`, `-`, `_`, parentheses, brackets, slashes, and similar separators). Reject non-ASCII input or an empty result and ask for an explicit readable `name_token`; silently dropping non-ASCII letters would not be readable. Map `cli-default` to `clidefault`. If a model label has a terminal parenthesized tier equal to the selected/derived effort, remove that tier before model sanitization, so the information appears once rather than twice.

The required results are:

| Stable ID | Exact model / effort | Display name |
|---|---|---|
| `claude-1` | Opus 4.8 / max | `claude-opus48-max` |
| `codex-1` | GPT-5.5 / xhigh | `codex-gpt55-xhigh` |
| `hermes-1` | GLM 5.2 / high | `hermes-glm52-high` |
| `antigravity-1` | Gemini 3.5 Flash (High) / high | `agy-gemini35flash-high` |
| `kimi-1` | K3 / max | `kimi-k3-max` |

Parsing accepts exactly three tokens or four when the last is an integer at least 2. It returns the three sanitized tokens plus the optional instance. It does **not** reconstruct the authoritative model or reasoning because sanitization is lossy; those remain stored fields. A name that does not round-trip through grammar validation is invalid.

Collisions are resolved within the effective roster. Preserve every surviving stable-ID-to-display allocation on reinit. For new colliders, sort by stable ID, give the first unallocated candidate the unsuffixed form, then allocate `-2`, `-3`, and so on; never compact suffixes automatically after a member leaves. This makes repeated discovery order irrelevant and prevents a harmless reinit from renaming surviving displays. An explicit future compact operation may be offered, but it must not affect base IDs or open-idea snapshots.

For agents whose effort is not a launch flag, `Reasoning` remains the single semantic source of truth and the adapter says how it is materialized: `flag`, `model-variant`, `config-file`, or `implicit`. For `agy`, the wizard presents the model/tier pair and derives `high` from `(High)` rather than pretending there is a separate effort switch. Hermes and Kimi adapters write/use their supported config mechanism. The display always reflects the effective semantic value.

### B. One rerunnable roster command, two write scopes

Use this surface:

```text
parley roster init --scope session [--dir DIR] [--dry-run]
parley roster init --scope machine [--dry-run]
parley roster init --scope session|machine --from FILE --yes
```

`init` is intentionally rerunnable; do not add a second `--reinit` state machine. `session` means the current persistent deck/project, not an ephemeral process session, and help text must say so. Interactive operation always shows a final diff and asks once before writing. `--from FILE --yes` is the non-interactive path; `--yes` alone must not invent undiscoverable model choices.

Discovery should extend `internal/agents/discover.go`. `Discovery` needs `ModelOptions`, `ReasoningOptions`, their provenance, and nonfatal capability-probe errors. Each adapter may define a bounded read-only command/parser such as `agy models`; there is no credible universal output parser. When a CLI cannot enumerate a dimension, offer the current configured exact value and `cli-default`, plus explicit text entry. For coupled dimensions such as Agy's model tier, present valid pairs. After selection, the adapter must assemble the real invocation. Today changing `Spec.Model` alone often does not change execution because `HeadlessArgs` hard-code model/effort and `buildAgentInvocation` substitutes only `{root}` and `{prompt}`. Add structured adapter assembly or `{model}`/`{reasoning}` substitution and test the final argv, not only metadata.

Scope writes differ as follows:

- `machine` atomically updates the managed agent roster/profile portion of `~/.parley/agents.toml`. It preserves `[defaults]`, `[rosters.*]`, unknown keys, and unrelated comments; on the first conversion it creates a backup before establishing managed-block markers. It writes no deck file and never edits `COOPERATION.md`. It is the seed for future projects.
- `session` writes portable selected fields to `parley-deck/agents.toml`, writes resolved machine-local paths/argv and the same stable IDs to the gitignored `parley-deck/meta/headless-agents.local.json`, and performs a targeted update of only the §2 roster table, adding a `Display name` column. The JSON is generated local state, not a second authority. `LoadAgentSpecs` currently does not read it and instead reads `agents.local.toml`; implementation must either add the JSON as a documented local layer or migrate the protocol to `agents.local.toml`. It must not claim both are active. My preference is to retain `agents.local.toml` as the loader's local override and make the required JSON a generated compatibility/cache file with a schema/version and source fingerprint.

All session outputs should be prepared and validated before any rename; on failure, restore the prior set. Validate stable-ID uniqueness, display uniqueness, adapter existence, name grammar, autonomous workspace mode, and TOML/JSON round trips.

`parley init` should not duplicate the wizard. On a newly created deck, after transport/bootstrap scaffolding, it calls the same internal `RosterInit(scope=session)` service seeded by machine defaults. On an existing bootstrapped deck, `parley init` remains idempotent and does not rerun selection; the explicit roster command is the only rerun. `EnsureCentralDefault` may still seed a missing machine file, but it is not a substitute for confirmed session selection.

The §9.0 readiness check remains a liveness check only. It loads stable roster IDs, resolves each through `LoadAgentSpecs`, and PONGs the effective adapter invocation. It must not enumerate models, ask for effort, add agents, rename displays, or edit the roster. `preflight.rosterEntry.RosterID` remains the stable ID while `Runtime` comes from `AdapterID`/command. The explicit final confirmation in `roster init` satisfies confirmation of that roster edit; per-idea unavailable-agent exclusions still follow §9.0.

Reinitialization must scan nonterminal ideas and live runs before applying a session change. It never rewrites `participants:`, round/review filenames, `responding-to`, signoffs, or presets. Removing a base ID referenced by an open idea is deferred: retain it as inactive-for-new-ideas with its pinned runtime profile until the idea/run is terminal. Model/effort changes apply only to future ideas; open runs use their manifest snapshot. For legacy runs without a reconstructable snapshot, retain the old effective spec and warn rather than silently adopting a new profile. Machine-scope reinit cannot scan every project, which is another reason run snapshots are mandatory.

Concrete Go touch points include `internal/app/app.go` (new top-level `roster` command, reuse from `runInit`), `internal/config/runtime.go` (schema, layered merge, managed writes, `LoadAgentSpecs` instantiation), `internal/protocol/roster.go` (targeted §2 parser/editor; current regex accepts only lowercase hyphenated IDs), `internal/app/preflight.go` (stable roster ID versus adapter/runtime), `internal/runner/runner.go` (artifact path stays `ID`; vendor behavior moves to `AdapterID`; structured argv), `internal/runmanifest/manifest.go` (profile snapshots), and consensus/driver tests that prove display changes do not change artifact lookup.

### C. Mandatory autonomous artifact-write profiles, safely bounded

I support making this first-class and mandatory, but reject “full bypass where needed” for ordinary participant work. A participant is assigned one file under the workspace; needing machine-wide write permission is a configuration failure, not a reason to broaden authority. A vendor flag named “dangerous” is not evidence of workspace confinement. The resolved profile must declare and verify confinement (`cli-sandbox` or an outer sandbox); cwd plus prompt wording alone must not be reported as enforced scope.

The skill should say exactly:

> **Autonomous participant writes are required.** Before launching a headless participant, resolve and validate that agent's `autonomous_write` profile. The invocation must be non-interactive and able to create its assigned canonical artifact without a permission prompt, while writes remain confined to the selected workspace/deck. Use the vendor-specific mapping below; there is no universal `--yolo` flag. If a verified workspace-scoped profile is unavailable, do not launch that participant and report a blocker. Never replace workspace confinement with a machine-wide bypass. Continue to redact obvious secrets from prompts, argv, generated config, logs, and discovery output, and never send credentials, customer data, or unrelated private files merely because autonomous mode is enabled.
>
> - Claude Code: `--permission-mode bypassPermissions`, with effective workspace confinement and the selected root supplied to the launcher.
> - Codex: headless execution with `-s workspace-write` and a non-prompting approval policy; do not use a full-danger bypass for protocol artifacts.
> - Hermes: `--yolo`, under the workspace-confined launcher profile.
> - Antigravity (`agy`): `--dangerously-skip-permissions`, under the workspace-confined launcher profile.
> - Kimi Code: plain headless `-p`; this mode already auto-approves in-workspace writes. Do not add `--yolo` or `--auto`, because each is mutually exclusive with `-p`.

Generated config should store structured mode and scope, not merely bury the flags in `HeadlessArgs`. Launch assembly may still flatten `AutonomousWrite.Args` into argv, but verification should reject a selected headless agent when the effective invocation omits or contradicts its required profile. Add table-driven tests for every mapping, especially Kimi's prohibited combinations and Codex's workspace-only policy. Existing `SandboxMode` and `ApprovalPolicy` remain useful runtime facts, but they do not replace the first-class capability because some vendors express it as an implicit mode or a single combined flag.

## Concerns / open questions

1. The prompt requires writing `meta/headless-agents.local.json`, while the current Go loader consumes `agents.local.toml` and not that JSON. Consensus must name one runtime authority and a compatibility path; otherwise the wizard can appear successful while launch behavior ignores its output.
2. True workspace confinement for Claude/Hermes/Agy must be demonstrated per supported OS. Their bypass flags alone are broader permission policies, not filesystem sandboxes. If no portable outer confinement exists, the matrix should honestly report `ScopeEnforcement: unverified` and fail the mandatory autonomous profile rather than label it safe.
3. A lossless edit of user-owned TOML is required for machine scope. `go-toml/v2` unmarshal/marshal is not comment preserving; managed blocks or a lossless document editor are preferable to rewriting the whole file.
4. The canonical base ID for new configs should include the instance suffix (`claude-1`) even though current built-ins and CLI examples use `claude`. A migration warning period is needed before removing legacy adapter-key selection.

## Risks

- Conflating `ID`, adapter, and display will reintroduce path instability through an apparently harmless refactor. Tests should reinitialize model/effort mid-run and prove that `round-NN/<id>.md`, missing-artifact checks, cross-review headings, event projection, and signoffs remain keyed by the original base ID.
- Metadata-only selection is a serious false-success mode: unless adapter assembly changes argv/config, the chosen model and effort displayed in the roster may differ from the model actually invoked.
- Collision allocation based on discovery order will churn names across machines. Persisted allocations plus deterministic stable-ID ordering are necessary.
- Autonomous flags can enlarge authority and expose secrets through argv/logging. Workspace enforcement and the existing secret-redaction boundary must be independently tested; “yolo” must never imply consent to external data disclosure.
- Editing §2 during reinit touches a protocol file. The editor must be narrowly table-scoped and the protocol/default copies must remain synchronized according to the idea's drift-guard constraint.
