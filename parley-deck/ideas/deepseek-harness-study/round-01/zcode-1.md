---
agent: zcode-1
idea: deepseek-harness-study
round: 1
date: 2026-08-26
---

## Summary

I cloned `deepseek-harness` at `dsh-v0.1.1-rc.2` into `/tmp/dsh-study` (PRIMARY:
`git clone --depth 50 https://github.com/deepseek-ai/deepseek-harness`) and unpacked the
published CLI (`npm pack @deepseek-ai/dsh` → `deepseek-ai-dsh-0.1.1-rc.2.tgz`, PRIMARY). The
brief's technical picture of the plugin model is accurate — I found no errors, only omissions.
The omissions matter more than I expected:

1. **dsh is explicitly pre-release with no compatibility promise.** Root `AGENTS.md` carries a
   "Pre-release stance: foundation over blast radius" section: rename/repackage freely, backends
   reject old on-disk formats, `SESSION_FORMAT_VERSION` stays `0` with no promise, and the section
   is to be deleted only "at the first tagged release". The only tag in the repo is
   `dsh-v0.1.1-rc.2` (PRIMARY: `git tag`). (PRIMARY: `AGENTS.md` lines 5–7)
2. **dsh already ships one-shot delegation to Claude Code and Codex** as subagent providers —
   but they are excluded from the default production install, pin their own private CLI payloads
   instead of the host's installed binaries, and are strictly one-shot (fresh process,
   non-resumable conversation, no follow-up, no progress streaming, no persona). The other four
   parley-deck harnesses (hermes, kimi, opencode, zcode) have no provider at all. (PRIMARY:
   `docs/subsystems/subagent.md`; `.agents/notes/implemented/feature/2026-08-04-claude-code-and-codex-subagent-backends.md`;
   `.agents/notes/implemented/simplification/2026-08-12-production-dsh-excludes-product-subagent-providers.md`)
3. **The tool-plugin half of the owner's sketch does not exist.** `mcpanywhere`, `graphify`,
   `openviking`, and `cognee` have zero occurrences in the entire repo; ToolSearch is one row in
   the cookbook describing a *pattern* (`ctx.tools.restrict()`), not a shipped plugin. (PRIMARY:
   `grep -ri "mcpanywhere|graphify|openviking|cognee" …` → no hits; `docs/cookbook/extension-cookbook.md` line 113)

My recommendation is **DON'T** rebuild parley-deck as a dsh plugin — not because dsh is bad (it
is the most disciplined plugin harness I have read), but because the audit record in this very
deck shows the *harness* is a load-bearing variable in parley-deck's results, and a dsh plugin
would either collapse that diversity or merely re-implement today's Go driver inside a
preview-stage runtime. Condition that would change my mind is in §6.

## 1. How the plugin model really works

All of this is PRIMARY from the cloned tree unless tagged otherwise.

- **Cordis, and what a plugin is.** A plugin is an object with `name`, optional `inject`, and
  `apply(ctx)`; it mounts into a context that is a repository of services addressed by stable
  keys (`ctx.tools`, `ctx.llm`, …). Dependencies are declared via `inject` and resolved by
  waiting, so load order is expressed as service requirements, not boot sequencing.
  (`docs/cordis-primer.md`, "Cordis In Five Ideas")
- **"Reversible effects" concretely.** Every contribution — prompt section, tool schema,
  adapter, provider, event listener — is installed through `ctx.effect()` / `ctx.on()`, and the
  registration call returns a disposer. Unload/HMR unwinds registrations in reverse; the repo
  convention is literally "Registrations are effects" (`AGENTS.md` line 103), and hot reload is
  listed as "every registration is a `ctx.effect` → vendored HMR just works"
  (`docs/cookbook/extension-cookbook.md` feature map). What reversibility covers is **registry
  and event state inside the Node process**. It does not cover external processes or files.
- **Does it hold for a plugin that shells out to a foreign process? No — and dsh says so
  itself.** The claude-code/codex provider note states the providers "do not resume sessions,
  stream progress, accept new human interaction, **roll back tool or file side effects**, or
  impose a wall-clock timeout" (Agent Note, Consequences). Process cleanup is handled not by
  effect-reversibility but by the `dsh-subprocess` seam: SIGTERM → `graceMs` → SIGKILL
  escalation on the process tree, waiting for whole-tree exit (`packages/subprocess/subprocess/src/types.ts`
  lines 181–188). And `packages/subprocess/subprocess-local/README.md` line 31 states that
  SIGKILL, fatal OOM, native crashes, and power loss "require an external supervisor, container
  init, or equivalent OS owner" — i.e. there is **no kill that survives the host's own death**
  (contrast §4, `procctl`). A parley plugin shelling out to six CLIs would get exactly the
  weakest guarantee at the point where parley-deck today has its strongest one.
- **Isolation and reach.** There is no sandbox between plugins. Providers are "trusted
  same-process implementations" (`docs/subsystems/subagent.md`, provider contract JSDoc); any
  plugin can reach any context service. Per-agent scoping exists (`core/scope`, an agent's
  `agent.ctx`, `tools.restrict()`), and child agents get "a new flat scope rather than inheriting
  parent registrations" (subagent.md, "In-process backends"). Isolation is *organizational*
  (scopes, seams, typed events), not enforced.
- **Versioning.** Every package is `@deepseek-ai/dsh-<name>` with vendored, rescoped Cordis as a
  peerDependency (`AGENTS.md` line 101). The published CLI depends on 62 `^0.1.1-rc.2` packages
  (PRIMARY: `package/package.json` from the tarball) — a floating caret on an RC closure. With the
  pre-release stance above, plugin ABI churn is expected, not exceptional.
- **Composition.** Confirmed as the brief states: profiles in `$DSH_HOME/profiles/<name>` list
  bundles in `dsh.profile.bundles`; each bundle's `dsh.bundle` patch applies in profile order,
  then profile `cordis.patch.yml`, then home-level patch, then `--patch`; `dsh --profile web
  --dump-config` prints the composed tree; out-of-tree plugins install per-profile via
  `dsh plugin` (forwards to pnpm in the profile directory). (PRIMARY: `docs/architecture.md`
  §"Profiles and bundles"; `package/README.md` of the tarball)
- **Events.** The brief's loop shape is confirmed (`docs/architecture.md` §"Turn flow"); two
  small completions: the stream includes `assistant/chunk*` before `assistant/message` and
  `tool/result*` after `tools/post-execute`; `agent/pre-step`, `agent/request`, `llm/stream` and
  the three `tools/*` events are waterfalls (must call `next()`), `agent/turn-stopping` is serial
  without `next()`. `turn/*`, `step/*`, `user/message`, `assistant/*`, `tool/*` are durable
  session events. `ctx.sessions.fork(source, boundary?, childSessionId?)` exists
  (`docs/architecture.md` extension-point table).
- **The invariant that will matter to us.** "Model-visible means logged": anything reaching a
  model request must be reconstructable from the session log; a runtime invariant asserts it, and
  a new model-visible input requires a new `SessionEventMap` event (PRIMARY:
  `docs/architecture.md` §"Session log"). A `SessionEventMap` member is required-on-read by
  default — builds that do not know an event's type refuse the log unless it carries
  `ignorable: true` (`AGENTS.md` line 105). This is the mechanism (and the cost) of adding
  parley-durable state.

## 2. Where parley-deck's concepts would land

Mapping each parley-deck concept onto a named dsh service/event (all service names PRIMARY from
`docs/architecture.md` tables and `docs/subsystems/subagent.md`):

| Parley-deck concept | Best dsh home | Fit |
|---|---|---|
| Idea (slug, directory, 00-prompt.md) | None. Files in the workspace, managed by our plugin's own code | No harness home; children of `ctx.subagents` run in the parent Session cwd, so the shared-tree file convention survives unchanged |
| Participants (agent roster, six harnesses) | `ctx.subagents` named provider instances (the named-instance decision supports multiple named codex/claude-code instances) | 2 of 6 products have providers; hermes/kimi/opencode/zcode would each need a custom provider (the `-acp` provider is a route only if those CLIs speak ACP — unverified) |
| Rounds | One `ctx.subagents.start()` one-shot run per participant per round, or `ctx.jobs` for background collection | Shape-compatible in principle (a round IS one self-contained prompt → one artifact), but pinned product versions and no progress streaming (see §3) |
| Round artifacts (round-NN/*.md) | Files; the run's `SubagentResult.output` is final text only, and `outputSchema`/`persona` capabilities are advertised as unsupported by the product providers | Artifacts stay a file convention; nothing in dsh validates or stores them |
| Consensus, quorum, signoffs, review, FINAL | A plugin-defined service (`ctx.parley.*`) plus durable session events (extend `SessionEventMap`) | Expressible — plugins can define services and durable events — but it is 100% our logic; the harness contributes only plumbing, and every event addition is a typed, versioned format change under a format version with no compat promise |
| Phases (the driver: idea → rounds → review → signoff → final) | A driver plugin orchestrating `ctx.subagents` + `ctx.jobs`; closest built-ins are `dsh-goal-round-driver` (same-session rounds) and the experimental Agent Teams (`ctx.agentTeams`: durable roster, task board, mailbox) | The orchestrator role survives, but it is the same code we have in Go, rewritten in TypeScript against a moving target |
| Durable kill / restart attribution | Nothing; `dsh-subprocess` is in-process only | Lost (§4) |
| Triage/verify commands (`parley consensus status` etc.) | `ctx.commands` (human commands that dispatch without a model turn) | Direct fit |

The honest summary: ideas/artifacts/consensus are *our* concepts and stay *our* code; dsh gives
them a process model (sessions, events, jobs, one-shot delegation) but no semantic home. The one
genuinely attractive native idea is the durable session log with "model-visible means logged" —
but parley-deck's durability today is the git-tracked artifact tree, which is stronger for our
purposes (auditable, diffable, survives any runtime).

## 3. Harness diversity: worth keeping, or does it collapse?

The brief is right that this deck holds controlled evidence. I read the audit record this
session (PRIMARY: `parley-deck/ideas/protocol-and-skill-audit/FINAL.md` in full; `consensus.md`
§4–§5; targeted reads of `signoff-hermes-1.md`, `signoff-kimi-1.md`, `signoff-zcode-1.md`,
`review/consensus.md`). The record:

- **One harness produced half the corpus.** @codex-1 wrote 24 of 47 findings, and cast all nine
  REFUTED verdicts in adversarial verification (PRIMARY: `FINAL.md` "Context & orientation";
  `consensus.md`). No other participant filed anything close.
- **One harness could not survive the session format.** @hermes-1 "filed at the fourth attempt,
  after three runs produced complete analyses and lost them all" (PRIMARY: `consensus.md` §4).
  @opencode-1 "created its skeleton file — progress on three prior total losses — but filed no
  finding" (same section). These are harness-level failures, not model failures.
- **One harness fabricated PRIMARY evidence — twice.** `review/consensus.md` line 220 records
  that an earlier round "discarded seven of its tie-break verdicts for fabricated PRIMARY
  evidence"; `round-02/hermes-1.md` cites commands against paths that do not exist in the repo
  (`python -m roster render`, `docs/protocol-appendix-a.md`) (PRIMARY: grep hits in
  `signoff-zcode-1.md` line 15 and `review/consensus.md` lines 220–221). `FINAL.md` lines 119–120
  generalize it: a reviewer reported PRIMARY evidence it had not executed, twice, and its
  signoff conclusions were discarded.
- **Cross-checking between different harnesses is what caught ledger errors.** The ledger
  mis-recorded verdicts nobody cast; "@codex-1, @opencode-1, @zcode-1 and @kimi-1 all flagged the
  arithmetic" (PRIMARY: `consensus.md` §5 correction note). @kimi-1 re-derived every number,
  found a disposition citing a commit that did not contain the fix, and mutation-checked a fix
  itself (PRIMARY: `signoff-kimi-1.md` Notes; `FINAL.md` "What this idea proved").
- **FINAL.md's own conclusion:** "The method that worked was not agreement… Two reviewers
  agreeing proved nothing here; a reviewer who ran the command proved everything." (PRIMARY)

Argument from that record: parley-deck's measured value comes from **independent harness+model
stacks producing independently verifiable evidence**. The owner's own framing — "the result of
the work is harness + model" — is confirmed by the audit: the harness variable changed outcomes
(finding volume, session survival, evidence honesty, verification style), and the safety net
(other harnesses re-deriving results) only exists *because* the stacks are independent.

Would a dsh plugin keep that? Two scenarios:

1. **Six models on dsh's loop** (dsh-native subagents, spawn/fork providers): collapses
   completely. Same loop, same tools pipeline, same session format, same failure modes. The
   audit's cross-harness safety net is gone by construction.
2. **dsh as coordinator, six external harnesses as one-shot delegations**: preserves diversity
   in principle, but today only codex and claude-code have providers, and those providers **pin
   dsh's own product payloads** (`@openai/codex@0.147.0`; Claude Agent SDK 0.3.220 / CLI 2.1.220
   selected from the SDK's platform packages) and "never fall back to a host `codex`/`claude`"
   (PRIMARY: provider Agent Note; production-exclusion note). The harness half of "harness +
   model" would be frozen by dsh's release cadence, not by the owner's live environment the
   audit ran on. Non-interactive-by-design (permission prompts denied/fail-closed), no follow-up
   or steering, fresh process per call — so no triage conversations, no re-prompting a
   participant, opaque long rounds.

So: harness diversity is worth keeping, and the realistic dsh topology weakens it in exactly the
dimension (live, heterogeneous, independently-configured harnesses) the audit proved load-bearing.

## 4. What we would lose

All PRIMARY from this repository this session unless noted.

- **The one-binary, no-runtime property.** `parley` is a 9.4 MB static Go binary in the repo
  root; `dist/` holds six OS/arch targets (darwin/arm64, darwin/x64, linux/arm64, linux/x64,
  windows/arm64, windows/x64). `dsh` requires Node `^22.19.0 || >=24.0.0` (PRIMARY: dsh repo
  root `package.json` engines) and installs a 62-package dependency closure (PRIMARY: tarball
  `package/package.json`). Offline use, copy-to-remote-host use, and "no runtime" all break.
- **Durable process kill.** `internal/procctl/procctl.go`: agents are spawned in their own
  process group with a durable identity (`Spawned{PID, PGID, BootID, ProcStart, Command,
  Marker}`) persisted in the `agent.started` event, so **a restarted parley can re-attribute and
  kill the tree**, gated by a strict attribution check so a reused PID is never killed by
  mistake; OS specifics in build-tagged darwin/linux/windows files. dsh's equivalent is
  in-process-only escalation and its own README externalizes SIGKILL/OOM/power-loss to "an
  external supervisor" (§1). Orphaned agent CLIs after a host crash are precisely the case
  parley-deck handles and dsh does not.
- **The drift guards.** `internal/protocol/drift_test.go`: the embedded bootstrap
  `COOPERATION.md` and the live deck must be byte-identical outside five allowlisted
  project-specific zones; anchors match column signatures rather than padding. This test also
  *records* a real caught defect (claude-1/F2: the protocol told the bootstrap to run a command
  that broke the repo's own guard). In a plugin world this becomes a TS test against a moving
  `SESSION_FORMAT_VERSION = 0` format — portable in principle, but the guarantee stays ours to
  rebuild either way; dsh contributes nothing here.
- **The release channels.** Verified present: `dist/` six-target binaries; Homebrew tap
  `feci/parley` with two formulas (`../homebrew-parley/Formula/parley-deck-cli.rb`,
  `parley-deck-skill.rb`); winget manifests `../winget-pkgs/manifests/f/Feci/ParleyDeckCli` and
  `ParleyDeckSkill`; the npm-packaged skill (`../parley-deck-skill`, per RECALL of my own past
  verification setup and its presence as a sibling with npm tests). **Correction to the brief:
  there is no `.github/` directory in this repo** (I listed the root), so I could not verify any
  Actions-based release automation; "six release channels" is SECONDARY (brief) — what I can
  verify is six *platform targets* plus tap/winget/npm distribution surfaces.
- **Stability of the substrate.** `CHANGELOG.md` 1.46.0 documents a disciplined release process
  (37 confirmed audit findings, 33 fixed, every fix with a revert-check test). Against that,
  dsh's own AGENTS.md instructs contributors to "rename or repackage freely", its backends
  "reject old on-disk formats", and its session format carries "no compatibility promise"
  (PRIMARY). A parley plugin would track that churn indefinitely.

## 5. The tool-plugin question (MCP/mcpanywhere, graphify, openviking, cognee)

- **MCP: same protocol, genuinely better ops, same substance.** `@deepseek-ai/dsh-mcp-client`
  is one plugin instance per MCP server; it discovers tools and registers them on `ctx.tools`
  under `mcp__<serverName>__<rawName>` — "the same server-qualified shape Claude Code and Codex
  use" (PRIMARY: `packages/mcp/mcp-client/README.md`). The real deltas over our current per-harness
  MCP wiring: HMR hot-swap without process restart, a budgeted reconnect supervisor
  (exponential backoff, max attempts, crash-loop exhaustion), deterministic collision-proof
  naming, and all-or-nothing generation rollback. Nice. But it is the same MCP servers with the
  same tool surface — "the same thing with a different config file" is fair, with an ops
  improvement we do not need for a coordinator that rarely restarts.
- **ToolSearch aware of mcpanywhere: does not exist.** ToolSearch appears once in the docs, as a
  row in the cookbook's feature map — "replace a scoped `ctx.tools.restrict()` registration as
  the visible set changes" (PRIMARY: `docs/cookbook/extension-cookbook.md` line 113). `mcpanywhere`
  has zero occurrences in the repo (PRIMARY: repo-wide grep). Writing an mcpanywhere-aware
  ToolSearch would be the same work it is today; the one thing dsh genuinely offers is the
  `ctx.tools.restrict()` seam that keeps presentation, lookup, and execution aligned in one
  registration — a clean primitive, not a delivered feature.
- **graphify / openviking / cognee: zero occurrences** anywhere in the repo (PRIMARY: repo-wide
  grep over md/ts/json, node_modules excluded). There are no first-party plugins to adopt. Each
  would be our own integration; the natural seams exist (`ctx.tools` for model-facing
  capabilities, `dsh-mcp-client` if the tool speaks MCP, `ctx.attachments` for image-bearing
  results). To be fair, dsh's capability-seam discipline (Service Definition / Provider /
  Consumer) would shape such plugins well — but that is architecture we can apply anywhere,
  including in the six harnesses we already use. `docs/graph-atlas.md` sounds related to
  "graphify" but is dsh's own documentation-graph index (PRIMARY: file header) — unrelated.

## 6. Recommendation

**DON'T** — do not rebuild parley-deck-cli as a dsh plugin package.

Reasons, in order of weight:

1. The audit record in this deck shows the harness is a load-bearing variable and that
   independence between stacks is the safety net (§3). A dsh plugin either collapses the stacks
   onto one loop or — in the coordinator scenario — freezes two of them at dsh's pinned product
   versions and leaves four without any provider.
2. The coordinator role gains almost nothing from dsh's actual value (its model loop, tool
   pipeline, prompt assembly) because parley-deck's participants are foreign processes. We would
   be using an agent harness as a task runner, while giving up the properties that make the
   current CLI reliable in that role (static binary, durable kill, drift guards, six-target
   releases) — §4.
3. The substrate is an explicit no-compatibility-promise preview (0.1.1-rc.2) iterating on
   exactly the seams we would depend on (session format, subagent contracts — the provider
   notes are dated weeks ago).

**Condition that would change my mind:** dsh reaches a tagged stable release with a compatibility
promise for `SessionEventMap`/subagent provider contracts, AND continuable (not one-shot)
providers that run the *host-installed* claude/codex/hermes/kimi/opencode/zcode binaries exist
or can be built cheaply for at least four of the six. At that point the right move is
**PLUGIN-ALONGSIDE**, not REBUILD: keep the Go CLI as the driver of record, and add a thin dsh
bundle exposing `parley` as a tool/commands so a dsh session can launch and observe rounds —
that is the only coupling that adds capability without subtracting independence. If the owner
wants one thing from dsh *today*, steal the idea, not the dependency: the "model-visible means
logged" invariant and the named-provider registry are excellent design patterns for parley-deck's
own adapter layer.

## Corrections to the brief

The brief's "what is already established" block is essentially correct — I verified every bullet
against `docs/architecture.md`, the cordis primer, and the tarball, and found no false claim.
Corrections and material omissions:

1. **Omitted (material): the pre-release stance.** Root `AGENTS.md` "Pre-release stance:
   foundation over blast radius": free renaming/repackaging, backends reject old on-disk
   formats, `SESSION_FORMAT_VERSION = 0` with no compatibility promise, until the first tagged
   release. The only tag is `dsh-v0.1.1-rc.2`; npm publishes `0.1.1-rc.2`. (PRIMARY)
2. **Omitted (material): codex/claude-code subagent providers exist** — one-shot, pinned private
   payloads (`@openai/codex@0.147.0`, Claude Agent SDK 0.3.220 / CLI 2.1.220), excluded from the
   default production closure, non-interactive, no follow-up/progress/persona. (PRIMARY)
3. **Omitted: `core/scope`** is in the core-packages table (the per-agent scoped-registration
   primitive), and the loop's waterfalls include `assistant/chunk*` and `tool/result*` events
   the brief's sequence omits. (PRIMARY)
4. **Wrong in the task framing (not the established block): "`.github/`" as something this repo
   has.** There is no `.github/` directory in parley-deck-cli. (PRIMARY: root listing)
5. **"Six release channels we ship on today" is unverifiable as stated.** What exists: six
   OS/arch binary targets in `dist/`, two Homebrew formulas, two winget manifests, the npm skill
   package. I could not find in-repo release automation. (PRIMARY for the artifacts; the "six
   channels" count itself remains SECONDARY.)

## What I could not verify

- **I never ran `dsh`.** I `npm pack`ed the published tarball and read the cloned source, but
   did not install or boot it (no profile init, no `--dump-config` output, no headless run —
   that would need `DEEPSEEK_API_KEY`). All runtime-behavior claims above are from source and
   docs, not execution.
- **Whether hermes/kimi/opencode/zcode CLIs speak ACP** (which would make `dsh-subagent-acp` a
   viable provider route for them). Not checked; would need each CLI's protocol documentation.
- **The upstream cordiverse/cordis repo** — I read dsh's vendored Cordis via the primer and
   conventions but did not clone cordis separately.
- **GitHub release/CI automation for parley-deck-cli** — no `.github/` in the repo; releases may
   be driven from elsewhere (the sibling `homebrew-parley`/`winget-pkgs` trees are forks/taps),
   but I found no workflow files to confirm.
- **Iteration speed beyond the tag**: my clone was `--depth 50` and resolved to the single merge
   commit of `dsh-0.1.1-rc.2`, so I could not measure commit cadence; the young tag, the rc
   suffix, and the Agent Note dates (2026-06 through 2026-08-18) are the cadence evidence I have.
