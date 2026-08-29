---
agent: kimi-1
idea: deepseek-harness-study
round: 1
date: 2026-08-26
---

## Summary

I cloned the harness source this session (`git clone --depth 1 https://github.com/deepseek-ai/deepseek-harness`
into `/tmp/dsh-study`, HEAD `b150a55`, merge of `release/dsh-0.1.1-rc.2`) and read its architecture,
contributor, Cordis-primer, cookbook, subagent, MCP, skill, and goal docs directly (PRIMARY for every
dsh claim below, file named per claim). I queried npm (`npm view @deepseek-ai/dsh version dist-tags
time engines bin`, PRIMARY). I did **not** install or run `dsh` — see "What I could not verify".

The brief's "already established" block is largely accurate — I confirmed almost every bullet against
source. What it omits matters more than what it gets wrong: **dsh already ships one-shot subagent
providers that delegate to Claude Code and Codex as foreign harnesses** (PRIMARY,
`packages/subagent/subagent-claude-code/README.md`, `packages/subagent/subagent-codex/README.md`),
and dsh is younger and less stable than "developer preview" suggests: first npm publish 2026-08-10,
ten RC releases in eleven days, currently `0.1.1-rc.2`, and its own AGENTS.md states there is no
tagged release and no on-disk compatibility promise (PRIMARY, npm `time` output; root `AGENTS.md`
"Pre-release stance"). This deck's own audit record (Q3) says harness heterogeneity is the mechanism
that catches defects, and this deck's own prior FINAL rejected cognee and graphify-as-selector
(Q5). Recommendation: **PLUGIN-ALONGSIDE**, with an explicit stability pin.

## 1. How the plugin model really works

**Composition.** A running `dsh` is a Cordis plugin tree composed at boot from ordered layers: each
bundle in the profile's listed order → the profile's `cordis.patch.yml` → the home-level patch → any
`--patch` overlay. A patch targets a row by id and replaces its whole config, or inserts rows
(PRIMARY, `docs/architecture.md` "Profiles and bundles"). Bundles and profiles declare themselves via
a `dsh` field in `package.json` (`dsh.profile`, `dsh.bundle`) — confirmed, brief bullet correct
(PRIMARY, same file). `dsh-base` is the first layer of every profile: model adapters, tools,
persistence, sandbox and approval policy, settings, credentials, telemetry (PRIMARY, same file).

**Plugin shape and lifecycle.** A plugin is a function with optional `inject` and `apply(ctx)`, or a
`Service` subclass. `inject` declares service dependencies; a plugin waits until its injected
services exist, so load order is expressed through service requirements, not boot sequencing
(PRIMARY, `docs/cordis-primer.md`). Function plugins named-export `name`/`inject`/`Config`/`apply`;
service packages default-export a class — mixing the forms makes the Loader discard the function
plugin's namespace, a real postmortem (PRIMARY, `packages/AGENTS.md`, citing
`docs/postmortem/0001-acp-default-export-drops-inject.md`).

**"Reversible effects," concretely.** Every contribution — prompt sections, tool schemas, adapters,
providers, listeners — is installed through `ctx.effect()` or `ctx.on()` and returns a disposer;
registrations unwind on reload/teardown, and registry contributions must *prove* disposal through an
HMR-safety test (PRIMARY, `docs/cordis-primer.md` "Practical Rules"; `packages/AGENTS.md`). The
cookbook claims "Plugin hot-reload: every registration is a `ctx.effect` → vendored HMR just works"
(PRIMARY, `docs/cookbook/extension-cookbook.md` feature map). The MCP bridge demonstrates the
boundary: editing a config row triggers disconnect+reconnect of the external server without a
process restart (PRIMARY, `packages/mcp/mcp-client/README.md`).

**Does reversibility hold for a plugin that shells out to a foreign process? Partially — and the
vendor says so itself.** The two foreign-CLI subagent providers own their child through
`dsh-subprocess`, which does tiered process-tree termination escalation and waits for whole-tree
exit on `dispose()` (PRIMARY, `subagent-claude-code/README.md` "Start and ownership";
`subagent-codex/README.md` same section). But both READMEs state verbatim under Known Limitations:
"**No wall-clock timeout or side-effect rollback** — the caller cancels long work, and files or
external systems changed before cancellation are not restored" (PRIMARY, both READMEs). So: Cordis
effects reverse *registrations*; process trees get reliable *killing*; side effects get *nothing*.
That is exactly parley-deck's situation today (a participant's edits are never rolled back either),
so this is parity, not regression — but anyone believing "reversible" covers foreign-process side
effects is wrong.

**Events.** Four dispatch modes — `emit`, `waterfall` (around-middleware; a listener that returns
without `next()` short-circuits), `parallel`, `serial` — and the dispatch mode is part of each
event's public contract (PRIMARY, `docs/cordis-primer.md`). The turn-flow sequence in the brief is
confirmed (PRIMARY, `docs/architecture.md` "Turn flow"), with one omission: the brief drops
`tool/result*` from the chain (see Corrections).

**Isolation.** There is no process isolation between plugins; isolation is *scoping*: the `core/scope`
library layers registrations (host layer vs per-preset scope chain, nearest layer wins a duplicate
name), agent presets compose per-session capability sets, and "a service row there needs an `isolate`
realm" (PRIMARY, `docs/architecture.md` mechanism table; `packages/skill/skill/README.md` describing
the host+per-scope layered registry). Spawned processes are confined separately via `ctx.sandbox`
backends (landlock/sandbox-exec; E2B as a remote POC) (PRIMARY, `docs/architecture.md`;
root `AGENTS.md` layout). A malicious or broken in-process plugin shares the host's memory and can
crash it — the docs I read promise no in-process fault isolation; I tag the crash consequence as
RECALL (standard for single-process plugin hosts), the absence of any documented isolation boundary
as PRIMARY.

**Versioning.** Plugins are npm packages installed per profile (`dsh plugin --profile <name> add
<pkg>`; PRIMARY, `subagent-claude-code/README.md`). The hard constraint for anyone storing state in
the session log: a `SessionEventMap` member is **required-on-read by default — a build that does not
know its type refuses the log** unless the event carries `ignorable: true`; only structural format
changes bump `SESSION_FORMAT_VERSION`, which currently sits at `0` "with no compatibility promise"
(PRIMARY, root `AGENTS.md` conventions + "Pre-release stance"). A third-party plugin that adds
durable session events therefore pins its sessions to itself until this stabilizes.

**What a plugin can reach.** Anything in its context: services by key (`ctx.tools`, `ctx.llm`,
`ctx.sessions`, `ctx.agents`, …), all typed events, `agent.inject()`/`followup()`/`steer()` into
live agents, `ctx.subprocess` for processes. There is no privileged core — the model adapter, tool
registry, session log and agent loop are themselves plugins (PRIMARY, `docs/architecture.md`
"Cordis"). Extension is deliberately wide; governance comes from convention (documented extension
points, invariant gates), not from a capability sandbox.

## 2. Where parley-deck's concepts would land

Mapped against the actual extension points (PRIMARY for the dsh side: `docs/architecture.md`
mechanism table, `docs/cookbook/extension-cookbook.md` feature map, named package READMEs; PRIMARY
for the parley side: this repo's `parley-deck/COOPERATION.md`-driven layout as evidenced in
`ideas/protocol-and-skill-audit/*` and `internal/*`):

- **Ideas, rounds, artifacts (`ideas/<slug>/round-NN/*.md`, `consensus.md`, signoffs).** These are
  files in a shared workspace that foreign CLI agents read and write with their own tools. In a dsh
  world the foreign participants keep doing exactly that — they are not dsh sessions, so the
  "model-visible ⟺ logged" invariant (PRIMARY, root `AGENTS.md`) does not bind their artifacts. The
  *orchestrator's* view of round state would land as **durable session events** ("Add durable session
  state → extend `SessionEventMap`", PRIMARY, architecture mechanism table) — with the versioning
  hazard from §1 — or deliberately stay as files, which is the safer choice at `SESSION_FORMAT_VERSION 0`.
- **Rounds as scheduled continuation.** The closest shipped analogue is `ctx.goals` +
  `dsh-goal-round-driver`: event-sourced same-session goal state with admitted-round counting,
  durable `goal/change` events, strict replay, and activation that is never persisted (PRIMARY,
  `packages/goal/goal/README.md`). It is single-agent, same-session — it schedules *one* agent's
  next round, not six harnesses' rounds. Useful shape to copy; not a home.
- **Participants (six foreign harnesses).** The `ctx.subagents` registry is the natural seam:
  multiple named provider implementations coexist, and `dsh-subagent-claude-code` / `-codex` already
  implement two of our six roster families (PRIMARY, `docs/subsystems/subagent.md`;
  `packages/subagent/` directory listing). But the shipped product providers are **one-shot**: "One
  fresh query and process per run — there is no continuation, resume, pooling, progress stream, or
  product-session persistence", `inheritsParentContext: false`, final-text-only results (PRIMARY,
  both product READMEs, Known Limitations). Parley rounds are multi-turn collaborations over shared
  files, so today only `subagent-spawn-in-process`/`-fork` (fresh child dsh agents) or new providers
  we write ourselves fit; kimi/hermes/opencode/zcode have no providers at all (PRIMARY,
  `packages/subagent/` listing — none exist).
- **Multi-agent coordination.** "Experimental Agent Teams" — a private opt-in seam on
  `ctx.agentTeams` with durable roster, task board, and mailbox over continuable subagents (PRIMARY,
  `docs/architecture.md` "Capability seams") — is the closest conceptual home for a rostered
  multi-agent run. It is explicitly experimental and excluded from official releases (PRIMARY,
  root `AGENTS.md` layout: `experimental/` = "private prototypes excluded from official releases").
- **Consensus, signoffs, quorum, drift gates, preflight, FINAL validation.** **No home.** Nothing in
  dsh validates cross-agent artifacts or counts signoffs; these would be a custom plugin exposing
  services + tools (`ctx.tools.register`) plus policy on `tools/pre-execute`/waterfall events — i.e.
  a TypeScript re-implementation of what `internal/consensus`, `internal/protocol` and
  `internal/driver` already do in Go (PRIMARY for the Go side: this repo's
  `internal/protocol/drift_test.go`, `internal/procctl/procctl.go`; RECALL for the full file list,
  though `docs/` and the audit CHANGELOG confirm the modules' roles).
- **Phases / driver loop.** A dsh-side facilitator would be a plugin listening on durable
  `session/event` and driving `agent.followup()`, plus `ctx.jobs` for background participant runs
  (PRIMARY, cookbook feature map: "Add background work → register on `ctx.jobs`"; UI/protocol-driver
  patterns). This is a genuine upgrade over the Go driver's blind spot: the Go runner never reads
  `COOPERATION.md` at all (PRIMARY, `ideas/speedup-tooling-evaluation/FINAL.md` "The structural
  finding" — @hermes-1 verified zero references in `internal/runner/runner.go`); a dsh facilitator
  would hold the protocol *in its context* as prompt sections rather than re-sending files.
- **The skill itself.** `ctx.skills` is a layered provider registry with model/user invocation
  policy and a shipped filesystem provider (PRIMARY, `packages/skill/skill/README.md`); the cookbook
  maps "Skills → section + tool registration; `inject()` skill content on invocation" (PRIMARY).
  The parley-deck skill lands here cleanly — as a runtime skill provider in the bundle.

## 3. Harness diversity: worth keeping, or does it collapse?

The audit is the controlled experiment: six harnesses, same task, same protocol, same corpus. Its
record, all PRIMARY from files I read this session:

- **Same corpus, 6 vs 36 confirmations.** Three verifiers adversarially checked 47 findings, told to
  default to REFUTED: @codex-1 confirmed 6 of 23 assessed, @kimi-1 36 of 42, @zcode-1 30 of 32; all
  nine REFUTED verdicts are @codex-1's (`ideas/protocol-and-skill-audit/consensus.md` §1). Whatever
  the structural explanation recorded there (codex-1 authored 24 of the 47 findings and was barred
  from assessing its own), the observable fact stands: **harnesses returned measurably different
  verdicts on identical inputs.**
- **One found defects the others missed.** @codex-1 wrote 24 of the 47 round-1 findings, including a
  23-item enforcement family nobody else surfaced (`consensus.md` §1–2). It also caught a fix with
  no regression test by mutation-checking (`FINAL.md` "What this idea proved about its own method").
- **One re-derived every number by hand — and that, not agreement, is what caught the ledger
  frauds.** The consensus ledger reported verdicts nobody cast (a row summing to 48 against 42
  assessed); @zcode-1 and @kimi-1 caught it independently; three reviewers caught the correction's
  own off-by-one "by re-deriving the numbers themselves" (`consensus.md` §1, §2 corrections;
  `FINAL.md`: "Two reviewers agreeing proved nothing here; a reviewer who ran the command proved
  everything").
- **One could not survive a long session.** @hermes-1 filed at the fourth attempt after three runs
  produced complete analyses and lost them all; what worked was a smaller slice plus
  write-first-append-as-you-go (`consensus.md` §4).
- **One fabricated PRIMARY evidence twice.** `review/consensus.md` (lines ~219-223, PRIMARY grep):
  "@hermes-1 filed a PRIMARY-tagged claim that does not survive execution — after an earlier round
  of this same idea discarded seven of its tie-break verdicts for fabricated PRIMARY evidence."
  Its conclusions were discarded both times and the fact recorded.

Conclusion: the audit's value came from *orthogonal failure modes* — one harness's blind spot
(leniency, session fragility, credulity) was another's catch. That is harness-level, not
model-level, diversity: the protocol, prompts, and corpus were identical, so the variance lives in
the harness+model combination, which is exactly the owner's intuition ("the result of the work is
harness + model", SECONDARY, 00-prompt.md). Collapsing to one harness with six model endpoints would
have turned the 47-finding audit into roughly the 24-finding audit (its best finder) with *no
independent verifier capable of refuting it* — @codex-1's 9 REFUTED verdicts would have gone
unchecked by the two harnesses that confirmed 66 of 74 assessed findings (~89%, `consensus.md` §1).

Does moving onto dsh collapse this? **No — dsh does not force the collapse; it already institutionalizes
foreign-harness delegation.** Its subagent seam keeps the child's harness authoritative: native
settings, model, tools, sandbox and authentication all stay with Claude Code/Codex; the provider
"does not select a model, create a product home, log in, or probe an account" (PRIMARY, both product
READMEs). A parley plugin on dsh can keep six real harnesses exactly as today. The collapse risk is
a *design choice* (writing providers that flatten everything into dsh subagent sessions), not a
consequence of the platform. The honest caveat: diversity also costs — three signoff rounds to fix
ledger arithmetic, two discarded signoffs, four hermes attempts (PRIMARY, above). The audit record
says the cost is paid back: every correction round found something real (`FINAL.md`).

## 4. What we would lose

PRIMARY unless tagged otherwise:

- **The one-binary install story.** Today: a static Go binary built with `go build -trimpath
  -ldflags "-s -w"` (`homebrew-parley/Formula/parley-deck-cli.rb`), shipped as raw binaries for six
  OS/arch targets (`dist/` — darwin/linux/windows × x64/arm64, verified by listing), a Homebrew tap
  (`Formula/parley-deck-cli.rb`, `parley-deck-skill.rb`), winget manifests
  (`../winget-pkgs/manifests/f/Feci/ParleyDeckCli/` up to 1.45.0), and `scripts/install-local.sh`.
  dsh requires Node `^22.19 || >=24` (dsh root `AGENTS.md` commands) plus pnpm for source work; its
  foreign-CLI bundles pull payloads of 75 MB packed / 257 MB unpacked (Claude Code darwin-arm64) and
  111 MB / 275 MB (Codex) — disclosed in dsh's own READMEs. "Six release channels" as phrased in my
  orders is SECONDARY; what I verified is four channel *types* (brew tap, winget, raw dist binaries,
  install-local) over six platform targets, plus GitHub source tags the formula downloads.
- **A moving foundation.** dsh has no tagged release; its AGENTS.md instructs "Remove this section
  at the first tagged release", prefers "the correct foundation over compatibility shims", warns
  "Backends reject old on-disk formats", and holds `SESSION_FORMAT_VERSION` at 0 (dsh root
  `AGENTS.md`). npm shows ten publishes between 2026-08-10 and 2026-08-21 (`npm view` time output).
  Rebuilding now means our durable artifacts ride a format with zero compatibility promise.
- **Durable kill.** `internal/procctl` captures a durable process identity (PID, PGID, boot id,
  proc-start, command, marker env) persisted in the `agent.started` event, and kills whole trees
  *across a parley restart*, refusing to kill when attribution fails (`internal/procctl/procctl.go`
  package doc and `Spawned`; the seatbelt/sysctl behavior is corroborated in the audit's
  `consensus.md` "One reproducibility fact"). dsh's `dsh-subprocess` does tiered tree termination and
  waits for whole-tree exit (dsh subagent READMEs) — real, but *session-lifetime* control; I found
  no documented cross-restart re-attribution-and-kill of a foreign process. Rebuild = port and
  re-prove this on three OSes.
- **The drift guards.** `internal/protocol/drift_test.go` byte-compares the embedded protocol
  template against the live deck with a deliberately narrow allowlist (read directly). Nothing
  comparable exists on the dsh side for our protocol; it would be a from-scratch TS port.
- **Offline/dev-loop properties.** The Go gates and `go test ./...` (27 packages, exit 0 — audit
  FINAL acceptance criteria) run fully offline; only the model calls need network. A dsh-based stack
  adds npm-mediated plugin distribution and fast-moving rc pins to every install and CI lane.
- **`.github/` CI — correction to my orders:** there is **no `.github/` directory in this repo**
  (`ls` failed; verified). So "the CI workflows we'd lose" cannot be enumerated from here; release
  automation presumably lives outside this tree or is manual into `dist/` — unverified either way.

What we would *gain* (to keep the ledger honest): a real orchestrator model in the loop with durable
sessions, fork (`ctx.sessions.fork(source, boundary?, childSessionId?)`, confirmed in
`docs/architecture.md`), jobs, hot-reloadable policy, and a tool/prompt-section registry — the
things the Go driver structurally cannot do because it never reads the protocol (Q2, last-but-one
bullet).

## 5. The tool-plugin question (MCP/mcpanywhere, graphify, openviking, cognee)

What dsh actually offers (PRIMARY, `docs/cookbook/extension-cookbook.md` feature map;
`packages/mcp/mcp-client/README.md`; `packages/skill/skill/README.md`): MCP tools arrive as raw
JSON-Schema registrations on `ctx.tools` (`mcp__<server>__<name>`, one plugin per server, stdio or
streamable-http, HMR reconnect); "ToolSearch / progressive disclosure" is a scoped
`ctx.tools.restrict()` replacement that keeps presentation, lookup and execution aligned; policy
waterfalls (`tools/pre-execute`, `ctx.tools.guard()`) give monotonic denials. This *is* better than
bolting MCP config onto each foreign CLI separately — for the **orchestrator's own** tools: typed,
scoped per agent, hot-swappable, guardable. But note the boundary the subagent READMEs draw: foreign
participants keep their *native* settings, including their own MCP config — a dsh-side MCP plugin
does **not** reach into claude-code/codex children (PRIMARY, both product READMEs: native settings
authoritative). So for the six participants, MCP stays per-harness exactly as today; `ctx.tools`
only upgrades the facilitator.

The named tools, with honesty about source quality:

- **mcpanywhere**: no presence in this repo outside this idea's prompt (PRIMARY, repo-wide grep).
  Web search identifies MCP Anywhere as a self-hosted open-source MCP gateway — one URL for MCP
  clients, Docker-isolated servers (SECONDARY, mcpanywhere.com and github.com/locomotive-agency/mcp-router).
  As a gateway it would sit *behind* `dsh-mcp-client` as one streamable-http server; dsh adds nothing
  gateway-specific, and the gateway adds nothing dsh-specific. Same wiring, different config file —
  with the scoped-registry upside above.
- **graphify**: installed and used here — `graphify-out/` holds 30+ dated report dirs with
  `graph.json`/`GRAPH_REPORT.md` (PRIMARY, listing). This deck already ruled on its *authoritative*
  use: "graphify is kept, outside the trusted path — useful for exploration and for *auditing* a
  phase packet, never for *selecting* one" (PRIMARY, `ideas/speedup-tooling-evaluation/FINAL.md`).
- **openviking**: ByteDance/volcengine's agent context database (memory + RAG + skills, `viking://`
  URIs, tiered loading), AGPL-3.0 (SECONDARY, web search: github.com/volcengine/OpenViking). Not
  present in this repo (PRIMARY, grep). The AGPL license alone needs an owner decision before
  bundling.
- **cognee**: already evaluated and **rejected** by this deck: "solves a problem we do not have —
  this deck *over*-remembers (quadratic history re-send); it does not under-remember" (PRIMARY,
  `ideas/speedup-tooling-evaluation/FINAL.md`, decision reached independently by three agents).

The deck's own binding doctrine for all four: "any tool that decides what an agent sees occupies the
exact position of the frontier machinery deleted in 1.43.1, which was removed because it could not
prove it never drops a participant objection… Selection of normative context [must] never be
delegated to a tool" (PRIMARY, same FINAL). A `ctx.tools` plugin does not change that doctrine —
it makes retrieval tools *nicer to wire*, not *admissible for normative selection*. Net: the
tool-plugin bundle is a reasonable convenience for the facilitator's exploration tools (graphify
read-only, MCP gateway), and a protocol violation if it ever selects what participants see.

## 6. Recommendation

**PLUGIN-ALONGSIDE** — keep the Go CLI as the protocol authority (consensus, quorum, drift guards,
procctl, releases), and build a thin `dsh` bundle that lets a dsh session act as *facilitator*:
parley-deck skill as a runtime skill provider, the `parley` CLI exposed as guarded tools, rostered
foreign harnesses driven through the existing file protocol (with `dsh-subagent-claude-code`/`-codex`
evaluated later as optional transports, not replacements).

Not REBUILD, on the PRIMARY evidence of §4: no tagged release, ten RCs in eleven days,
`SESSION_FORMAT_VERSION 0` with no compatibility promise, and our durable-kill + drift-guard +
six-target release surface would all need re-proving on a foundation that explicitly prefers
"foundation over blast radius" right now. Not WAIT in the passive sense: the bundle experiment is
cheap and reversible precisely because Cordis plugins uninstall cleanly (§1) — but pin an exact rc
and treat every upgrade as a migration. Not DON'T: the owner's actual want — a harness-coupled
facilitator with real tool seams — is legitimate, and dsh's architecture (events, effects, seams,
existing foreign-CLI subagents) is genuinely the best fit for it I have seen; the Go driver's
structural blindness (never reads the protocol, §2) is a real defect this would fix.

**Conditions that would change my mind.** To REBUILD: dsh ships a stable release with documented
on-disk compatibility, *and* continuable (multi-turn) foreign-CLI subagent providers exist or prove
trivial to write for our six roster families, *and* procctl-equivalent durable kill lands behind
`dsh-subprocess`. To DON'T: if the subagent seam turns out unable to carry multi-round,
file-mediated collaboration without re-implementing parley inside dsh anyway — at that point the
plugin is a second parley with worse release engineering.

## Corrections to the brief

All bullets below: SECONDARY = the brief's text; the correction is PRIMARY against the named source.

1. **Loop shape omits `tool/result*`.** The brief's chain goes `… tools/post-execute → step/end`;
   `docs/architecture.md` "Turn flow" has `tool/call* → tools/pre-execute → tools/execute →
   tools/post-execute → tool/result* → step/end`. Minor, but `tools/result` is the documented
   observation point for immutable final outcomes (cookbook hook-plugin section) and would matter
   for any audit/telemetry plugin.
2. **"Powered by Cordis" undersells it: Cordis is vendored and patched.** dsh carries Cordis as
   pinned source under `vendor/` with upstream SHAs and logged local modifications (root `AGENTS.md`
   vendoring policy + layout). Upstream cordiverse/cordis is the ancestor, not the dependency you
   get. I did not read the vendor diffs — unverified how far they diverge.
3. **"Developer preview" is weaker than stated.** No tagged release exists; AGENTS.md says to remove
   its pre-release stance "at the first tagged release", with "no external consumers" assumed and
   `SESSION_FORMAT_VERSION 0` "with no compatibility promise". npm: created 2026-08-10, latest
   `0.1.1-rc.2` (2026-08-21), ten versions in eleven days. Plan for format breakage, not just API
   churn.
4. **Missing fact, not wrong fact: foreign-harness subagents already exist.** The brief never
   mentions that `dsh-subagent-claude-code` and `dsh-subagent-codex` ship as optional profile
   bundles (one-shot, final-answer-only, native settings authoritative). This is the single most
   decision-relevant dsh fact for this idea.
5. **My orders' Q4 hints were partly wrong:** `.github/` does not exist in this repo (verified by
   `ls`); and "six release channels" is not what I can verify — I verify four channel types (brew
   tap, winget manifests, raw `dist/` binaries, install-local script) across six OS/arch binary
   targets. (`Formula/` and `winget-pkgs` do exist at `../homebrew-parley` and `../winget-pkgs` as
   stated.)
6. **Confirmed as accurate** (so nobody re-checks): the `dsh` package.json fields, composition order,
   the six core-service `ctx` keys, durable-event durability, and `ctx.sessions.fork(source,
   boundary?, childSessionId?)` all match `docs/architecture.md` verbatim.

## What I could not verify

- **I did not install or run `dsh`.** I only cloned the repo (into `/tmp/dsh-study`, as stated) and
  queried npm metadata. No `pnpm install`, no `dsh --dump-config`, no live plugin mount — the
  workspace build is heavy and real runs need a `DEEPSEEK_API_KEY` I do not have. Every runtime
  claim above is therefore documentation-PRIMARY, not execution-PRIMARY. Composition, hot-reload,
  and subagent behavior could differ from docs; I saw no evidence they do.
- **How far vendored Cordis diverges** from cordiverse/cordis (vendor manifest not read).
- **What release automation parley-deck-cli actually uses** (no `.github/` here; dist/ population
  mechanism unverified) — relevant to the true cost of losing the release story.
- **Whether any dsh persistence backend (SQLite) already offers cross-restart process re-attribution**
  comparable to procctl; I found none documented, which is not proof of absence.
- **mcpanywhere / openviking internals** — identified via web search only (SECONDARY); not installed,
  not tested, and absent from this repo.
- **The `experimental/agent-team` seam's actual API** (durable roster, task board, mailbox) — cited
  from one architecture paragraph; the subsystem doc was not read.
