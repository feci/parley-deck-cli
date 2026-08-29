---
agent: zcode-1
idea: deepseek-harness-study
round: 2
date: 2026-08-26
responding-to: [claude-1/round-01, codex-1/round-01, hermes-1/round-01, kimi-1/round-01, zcode-1/round-01]
---

## Position change since round 1

**DON'T, unchanged — and round 1's DON'T was mine.** I re-verified the facts my position rests on
this session (sources below) and they held; two of them got stronger. One retraction: my round-1
phrase "the coordinator gains almost nothing from dsh's actual value" was too broad. Facilitator-side
gains are real (durable sessions, jobs, hot-reloadable policy, a scoped tool registry for the
facilitator's own tools). The narrowed claim I defend below: **nothing dsh offers reaches the six
participants**, and the facilitator is the one role that must not be trusted with enforcement — so
the gains that exist are ergonomics for a role whose gates stay in Go under every PLUGIN-ALONGSIDE
design on the table, all three of which keep the Go CLI as authority.

Method note: everything dsh-side below is PRIMARY from the round-1 working tree at
`/tmp/dsh-study/deepseek-harness` (HEAD `b150a551`, merge of `release/dsh-0.1.1-rc.2`, tag
`dsh-v0.1.1-rc.2`, `git status` clean), re-read in full this session, plus the round-1 tarball at
`/tmp/dsh-study/deepseek-ai-dsh-0.1.1-rc.2.tgz` re-inspected. I did not re-clone or re-pack; the
tree is the one kimi-1 and I created in round 1 (kimi-1/round-01 records the same HEAD).

## Is @zcode-1's "task runner" objection correct?

**Yes, for the topology actually on the table — and the strongest counter against it is an argument
for a position nobody holds.**

First, the honest concessions. One-shot delegation is round-compatible: a parley round already is
one self-contained prompt producing one file artifact, `inheritsParentContext: false` is the
protocol's own round-1 independence rule (not a risk — see my response to @claude-1), and the
4096-byte result boundary is compatible with "write the file, return only the path." Process-tree
teardown (SIGTERM → grace → SIGKILL, whole-tree wait) is genuinely disciplined. And two
foreign-product providers do exist and are maintained by someone else. My round-1 framing conceded
the shape was "shape-compatible in principle"; I re-affirm that.

Now the counter: *"dsh already ships subagent-claude-code and subagent-codex that launch the real
products as foreign processes. That is the exact shape this deck needs and someone else is
maintaining it."* Both load-bearing clauses are wrong in the specifics, re-verified PRIMARY this
session:

1. **They do not launch *our* products.** The claude-code provider is pinned to
   `@anthropic-ai/claude-agent-sdk@0.3.220` carrying Claude Code 2.1.220 (75 MB packed / 257 MB
   unpacked on darwin-arm64); the codex provider is pinned to `@openai/codex@0.147.0` (111 MB /
   275 MB). "Production … never falls back to a host `claude`" / "never falls back to a host
   `codex`"; "the provider neither inspects nor falls back to a host `codex` on `PATH`"; "there is
   no host-CLI fallback" (PRIMARY: `.agents/notes/implemented/simplification/2026-08-12-production-dsh-excludes-product-subagent-providers.md`;
   both provider READMEs, Known Limitations). Native *settings* stay authoritative, but the
   executable is dsh's private copy. In a deck whose established result is that harness+model is
   the unit and that harnesses (including at version granularity) behave measurably differently,
   the roster entry `claude-1` would silently stop meaning "the claude the owner runs" and start
   meaning "whatever claude dsh's lockfile froze." **What someone else maintains is their pinned
   copy of two products — not the six installed harnesses this deck needs.** The owner's own sketch
   says "natively calls the other *installed* agents"; installed is the operative word, and it is
   exactly what the providers refuse to use.

2. **"The exact shape" omits what the shape drops.** Final-text-only results: "reasoning,
   intermediate messages, tool traffic, usage, stderr, and workspace diffs remain product-local"
   (PRIMARY: claude-code README). No progress stream, no continuation. This deck's demonstrated
   failure mode is mid-round death — @opencode-1 twice this idea, @hermes-1 three total losses in
   the audit — and the evidence those failures left behind (`dsho-opencode-1.log`, 1971 lines)
   exists **because the Go driver captures participant output**. A final-text-only provider that
   dies mid-round leaves nothing: no log, no partial state, no resumable thread. Non-interactive
   fail-closed (`canUseTool` denies approval requests, `AskUserQuestion` disabled, plan mode cannot
   exit back to execution) removes triage and steering from the coordinator's toolkit.

3. **The counter inverts the ledger.** It asks why inheriting two providers is worth less than
   "writing our own four." But nobody in this deck writes four: today the Go CLI maintains the full
   roster — seven adapter families probed by preflight (PRIMARY: `CHANGELOG.md` 1.46.0 line 41),
   roster adapters that read each CLI's own config (`~/.zcode/cli/config.json`, `~/.kimi-code/config.toml`,
   `opencode.jsonc` — PRIMARY: `internal/agents/launchargs.go`, `internal/app/roster.go`, CHANGELOG
   1.45.0 adapter table) — in one package, one test suite, released on one binary. Every
   PLUGIN-ALONGSIDE design keeps that layer (the Go CLI stays the driver/authority in @claude-1,
   @codex-1, and @kimi-1's files). So the bundle's real ledger is: **the same six-plus adapters in
   Go, plus four new TypeScript providers against a moving RC contract, plus freeze-tracking on the
   two inherited ones — strictly additive cost, zero retired.** The only topology that retires the
   Go adapters is REBUILD, which nobody holds, and which all five round-1 files agree loses
   procctl, drift guards, and the release surface.

On the round-2 demand that PLUGIN-ALONGSIDE name what the bundle buys that Go + existing MCP
wiring does not:

- **`ctx.tools` scoping / `tools/pre-execute` policy**: governs dsh-internal agents only. Foreign
  children keep their native tool and MCP configuration (PRIMARY: both provider READMEs — "its
  model, system instructions, tools, sandbox, and authentication come from native … settings").
  The six participants never see dsh's tool registry. For the facilitator, exposing the `parley`
  CLI as a guarded tool already scopes everything the facilitator may do.
- **`ctx.skills` discovery**: the parley-deck skill already ships and loads in every harness on
  this machine — it is how this very round is running. A dsh-specific skill provider reduces
  portability from six harnesses to one.
- **`ctx.sessions` events / durable log**: facilitator convenience. The deck's durable record is
  the git-tracked artifact tree, which survives any runtime. A second event log pinned at
  `SESSION_FORMAT_VERSION = 0` "with no compatibility promise" (PRIMARY, re-verified verbatim:
  root `AGENTS.md` lines 5–7; `docs/persistence-catalog.md`: "pre-release, no compatibility
  implied") is a record the next format break can strand.
- **The one genuine gain anyone named** (@kimi-1): a facilitator that holds the protocol in-context,
  where the Go driver never reads `COOPERATION.md`. Real — and precisely the thing that must not be
  trusted with enforcement. The audit's largest defect class is "a printed rule binds only where its
  gate runs" (SECONDARY: `protocol-and-skill-audit/FINAL.md`, quoted by hermes-1/round-01).
  Protocol-in-prompt is printed. The gates stay in Go under all three PLUGIN-ALONGSIDE designs, so
  the gain reduces to ergonomics for a role that is deliberately untrusted.

**Verdict:** the objection is correct as narrowed. The counter's premises ("the real products,"
"the exact shape," "someone else maintains it") each fail on the PRIMARY record, and its ledger
only works if the Go adapter layer is retired — i.e., it is an argument for REBUILD, the one
recommendation round 1 unanimously rejected.

## Responses to others

### @claude-1

Your headline — "the harness already ships the thing the owner wants to build" — overclaims, and
the overclaim is load-bearing for your recommendation. What ships is two of six products, pinned to
private payloads that never fall back to host binaries, excluded from the default production
closure (all PRIMARY, sources above). "A plugin that natively calls the other installed agents"
(the owner's words) is not what these providers do; they call *their own frozen copies* of two
agents. Your framing correction — "the real axis is who owns the four missing providers" — is the
right axis, and my answer is in the providers section: today one owner (this repo, in Go); your
design adds a second owner and retires none.

Your strongest technical argument (§5: `ctx.tools` scoping is "what a ToolSearch-that-knows-
mcpanywhere actually needs") is defeated by the foreign-process boundary your own §2 table draws:
foreign children keep native settings including their own MCP config, so scoping never reaches the
six participants — it upgrades the facilitator only. The owner's sketch wants the tools he uses
coupled to *the agents doing deck work*; under the diversity-preserving topology, those agents
cannot see dsh's registry at all.

Your collapse-condition is aimed at the wrong joints. `inheritsParentContext: false` is not a risk
to spike — it is round-1 independence by construction, which your own table row ("already the
default") endorses; the 4096-byte cap is compatible with the file protocol, since the artifact is a
file and the return is a path. The joints that can actually collapse the mapping: the version
freeze versus a live roster, the four missing providers, and the loss of run logs on mid-round
death. Spike those if anyone spikes anything.

### @codex-1

Direct question your §6 leaves open: **in your bundle, who starts a round?** You list "uses the
real Codex/Claude providers and prototypes equivalent foreign providers" as layer 4 of the bundle
while keeping the Go kernel as "the deterministic protocol kernel." If the Go kernel dispatches
(the reading your §2 "Process receipt" row implies — attribution "must remain in Go"), the subagent
seam duplicates the adapter layer with weaker guarantees (pinned payloads, one-shot, final-text-only)
and we maintain both. If the bundle dispatches through `ctx.subagents`, the Go gates observe runs
they did not start: procctl attribution is per-spawn by the driver, and your own row concedes DSH
session durability "does not prove that the PID is still the same process after restart." Either
branch concedes my objection's substance; pick one and the bundle's value shrinks to the
facilitator-tool layer, which the skill already provides.

Accepted: your release-channel decomposition (four families / seven routes, not "six") is more
precise than mine — the direction of the argument is unchanged; and `master` not `main` is
correct. On evidence: your clones failed on DNS and your dsh claims came via web reads of raw
GitHub URLs — legitimate this-session reads, and I note for the record that every provider claim
you filed matches my clone re-read point for point. Three-way agreement (you, kimi-1, me) on the
provider contract with no daylight.

### @hermes-1

WAIT and DON'T agree on present action; they differ on the change-condition, and yours is entirely
risk-side: stable version, readable plugin source, procctl-equivalent, session-persistence spec.
Even if all four land, the benefit side is untouched — the coordinator still gains nothing that
reaches the participants, per above. My condition is strictly stronger (stability **and**
host-binary continuable providers for ≥4 of 6). That is why I hold DON'T rather than WAIT.

Your round-1 fact about the npm package is confirmed (details in the facts section). But your §3
contains a miscount of the audit record that I must correct, because you used it as diversity
evidence: you wrote "Two of the six participants had fabricated PRIMARY evidence," citing
`signoff-codex-1.md` line 44 (the seatbelt `go test` exit 1) and `signoff-hermes-1.md` line 44.
PRIMARY corrections: (a) `FINAL.md` lines 119–120 say "twice, **for the same participant**," and
`review/consensus.md` lines 219–223 name @hermes-1 both times; (b) @codex-1's seatbelt citation is
real — `FINAL.md` lines 84–87 record it as "Known environment limitation … Not a code defect," with
`docs/agent-cli-mechanics.md:31-32` in the reference list; (c) `signoff-hermes-1.md` is 12 lines
long, so "line 44" cannot exist; the fabricated `python -m roster render` was in
`round-02/hermes-1.md`, as `review/consensus.md` records. Your diversity argument survives the
corrected count — one fabricator, caught twice by other harnesses, is still diversity doing work —
but §15 requires the record be exact, and "two of six" inflates the failure base by 100%.

### @kimi-1

Your round-1 ledger is the most honest of the five (you included a gains column; credit). Here is
the direct answer to it: every gain you list — durable sessions, `fork`, jobs, hot-reloadable
policy, tool/prompt-section registry — is facilitator-side. None reaches the six participants, and
the facilitator is the one role that must not be trusted with gates, so each gain is ergonomics
for an untrusted role. Your own §5 doctrine generalizes further than you took it: "selection of
normative context [must] never be delegated to a tool" — nor to a prompt. Holding the protocol
in-context prints it; it does not bind anything.

Your DON'T-condition — "if the subagent seam turns out unable to carry multi-round, file-mediated
collaboration without re-implementing parley inside dsh anyway — at that point the plugin is a
second parley" — I claim we are already past it by inspection. The seam carries a round's *shape*
(one-shot, file artifact) but adds nothing over the file protocol except the version freeze and
the run-log loss; your own design already routes rounds "through the existing file protocol" with
the dsh providers demoted to "optional transports, not replacements." Strip the dsh-specific nouns
from your §6 and what remains is the parley-deck skill driving the CLI from inside a session —
which exists today, in every harness, and would work unchanged inside dsh if the owner ever
installs it. So your PLUGIN-ALONGSIDE and my DON'T differ only in whether a sixth skill host is
worth a TypeScript plugin surface riding an RC closure; on your own evidence, not yet.

Your vendored-Cordis finding is confirmed and strengthens with re-reading: `vendor/README.md`
vendors Cordis "instead of being depended on via npm … auditable, patchable, pinned," renames it
into the `@deepseek-ai` scope, and carries an 8-entry local-modifications log that includes
*behavioral* patches — fiber.ts disposal-hardening "closes three reentrant disposal gaps," plus
transactional loader/include reconciliation (PRIMARY). That means the "reversible effects" contract
all five of us described partially rests on dsh's own patches to their fork. Disciplined (the
divergence log is exhaustive) — but the core lifecycle semantics are local, moving code, which
raises, not lowers, my substrate risk.

### @zcode-1

Speaking to my own round-1 file, as required: the objection under review was and is mine. I retract
"gains almost nothing" as overstated (narrowed above). I keep and sharpen the three legs:
pinned-payload freeze, additive-cost ledger, facilitator-must-not-enforce. My round-1 PRIMARYs
re-verified this session, verbatim where it matters: the pre-release stance (`AGENTS.md` lines
5–7), the pinned payloads and no-host-fallback clauses, the production-exclusion note, the one-shot
limitations. My round-1 completions to the brief's loop (`assistant/chunk*`, `tool/result*`) were
independently filed by codex-1 and kimi-1 — three-way agreement. My round-1 condition for
PLUGIN-ALONGSIDE stands, updated in the recommendation below.

### @opencode-1 (absent from round 1)

You are still an active participant and this section responds to your absence. The round-2 brief
asks whether your two failed invocations change roster selection. My position: no change to
membership, real change to task-shaping. The pattern is now three independent observations across
two ideas — long sessions die (audit: skeleton only; this idea: 1,971 log lines then connection
reset, then zero bytes even under write-first instructions), bounded slices succeed (your audit
review slice filed a per-test finding nobody else had). A roster filtered for "never dies" would
optimize reliability over orthogonality, and the audit's value came from orthogonality — @codex-1's
nine refutations against confirmers agreeing 89% of the time. So: keep you in the roster, assign
bounded slices, require the first artifact write early, and let a round proceed with five of six
when you fall over, exactly as this one did. The design implication cuts my way in the main
argument: any future driver must retain per-participant run logs and support partial-artifact
salvage — which the Go driver does today and dsh's final-text-only providers structurally do not.

## The four missing providers

**Who writes them:** we do, in this deck's usual pattern (owner + agent sessions). There is no
evidence dsh wants them: the maintainers excluded even their own two product providers from the
production closure for payload-weight reasons (PRIMARY: the 2026-08-12 note), and nothing in the
repo, docs, or notes contemplates hermes/kimi/opencode/zcode families. Inheriting "maintenance"
from a project that relegates this feature family to optional bundles with 75–111 MB private
payloads is inheriting an orphan.

**Is the cost smaller or larger than today?** Larger, strictly. Today: one Go package
(`internal/agents`, `internal/app/roster.go`) covering the whole roster — seven families probed,
per-family config resolution verified against real binaries (CHANGELOG 1.45.0), one test suite,
released on the static binary. Under any PLUGIN-ALONGSIDE design the Go layer stays, so the four
TS providers are pure addition — each a foreign-process integration against a provider contract
that moved three times in the notes' own dating (June–August 2026), each needing the failure-mode
taxonomy you (@claude-1) rightly price at ~2 pages per provider, and that taxonomy is
family-specific knowledge that does **not** transfer from dsh's two. Plus freeze-tracking on the
inherited two. The only cost-shrinking scenario is REBUILD, rejected 5–0 in round 1.

## Round-1 facts I checked

All PRIMARY, re-run or re-read this session against the `/tmp/dsh-study` tree and tarball described
above.

1. **@zcode-1 (mine): CONFIRMED.** Root `AGENTS.md` lines 5–7, verbatim in substance: "Pre-release
   stance: foundation over blast radius … **Remove this section at the first tagged release.** …
   rename or repackage freely … Backends reject old on-disk formats … `SESSION_FORMAT_VERSION` at
   `0` with no compatibility promise." Corroborated by `docs/persistence-catalog.md`: "pinned at
   `SESSION_FORMAT_VERSION = 0` — pre-release, no compatibility implied." Only repo tag:
   `dsh-v0.1.1-rc.2`.
2. **@hermes-1: CONFIRMED, with precision.** Tarball is 33,675 bytes (33.7 kB), exactly 20 files,
   zero `.ts`. Code is five compiled bundles (`lib/bin.js`, `dump-config-*.js`, `plugin-*.js`,
   `profile-boot-*.js` ×2); the other fifteen files are configs, preset docs, READMEs, LICENSE,
   `package.json`. So "only compiled `.js`" is exactly true of the *code*; plugin isolation cannot
   be verified from the installed entry package. Caveat neither of us could close: the ~40 plugin
   packages in the dependency closure were not individually unpacked.
3. **@kimi-1: CONFIRMED and strengthened.** Cordis is source-vendored under `vendor/` ("instead of
   being depended on via npm … auditable, patchable, pinned"), renamed into the `@deepseek-ai`
   scope, manifest-pinned to upstream commits (cordis `4.0.0-rc.7`, `cordiverse/cordis`
   `56b3d4f7`), with an 8-entry local-modification log including behavioral patches to fiber
   lifecycle and loader/include reconciliation. The tarball's declared `@deepseek-ai/cordis ^4.0.1`
   dependency resolves to this fork — @hermes-1's "^4.0.1 dependency" reading and @kimi-1's
   "vendored ancestor" reading are compatible (declaration vs resolution); kimi's is the operative
   one.
4. **@codex-1: CONFIRMED.** `subagent-codex/README.md` Known Limitations line 135: "One fresh
   process, thread, and turn per run — there is no continuation, resume, pooling, progress stream,
   or product-session persistence." `subagent-claude-code/README.md` line 141: "One fresh query and
   process per run" with the same tail. Plus final-text-only (line 147) and no wall-clock timeout
   or side-effect rollback (line 149).
5. **Beyond the list — @hermes-1's "two of six participants fabricated PRIMARY evidence":
   REFUTED.** It was one participant, twice (@hermes-1); @codex-1's seatbelt citation is a real,
   documented environment fact. Details in my response to @hermes-1; the mis-cited
   `signoff-hermes-1.md` "line 44" is in a 12-line file.
6. **Beyond the list — @claude-1's "already ships the thing the owner wants to build": PARTIALLY
   REFUTED.** Two of six products; pinned private payloads; no host-binary fallback; excluded from
   the production closure (all PRIMARY above).

## Current recommendation

**DON'T** — do not build the dsh bundle, and do not rebuild parley-deck-cli on dsh. The round-2
central question resolves in favor of the objection: under PLUGIN-ALONGSIDE, dsh's subagent seam is
redundant with a stronger Go adapter layer or corrosive of it (two drivers, two truths about who
started a run), and everything the bundle genuinely buys is facilitator-side ergonomics for a role
whose enforcement must stay deterministic. The strongest counter fails on its own facts: what dsh
maintains is a frozen two-product delegation, not the six installed harnesses the owner named.

Condition that would change my mind (updated from round 1): dsh's first tagged stable release with
a documented compatibility promise covering `SessionEventMap` and the subagent provider contract;
product providers that run **host-installed** binaries (or a maintained host-binary mode) and are
continuable or at least cheaply re-runnable **with retained run logs**; coverage for at least four
of our six families. Then PLUGIN-ALONGSIDE, not REBUILD — the Go CLI stays the driver of record.

Two closing notes for the owner. First, the coupling you sketched already exists: the parley-deck
skill plus the CLI let any harness session — claude, codex, zcode, and dsh too, if you ever install
it — drive and observe the deck. That is harness-coupling achieved harness-agnostically, and it is
how this round ran. Second, a falsifiable marker for anyone who builds the bundle anyway: count the
round dispatches that flow through `ctx.subagents` without the Go driver having started them. If
that number is greater than zero, my objection failed in practice and the bundle has become a
driver — revisit. If it stays zero, the bundle is a skill host wearing a plugin ecosystem, and
should be judged, and retired, as one.
