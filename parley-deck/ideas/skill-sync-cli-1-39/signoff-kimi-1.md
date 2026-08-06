### kimi-1

**Verdict: accept with reservations** (one non-blocking reservation, item R1 below).

#### VC-1 — I accept the closure

The drafter made the measurement and benefits from it, so per §15.1 I re-checked what I could
reach myself instead of endorsing it. I split claim (a) into its two atomic parts, because my
ownership differs across them.

**(a1) Zero Go references to `writeModeArgs` / `write_mode_args` — CONFIRMED (PRIMARY).**
I executed a repo-wide ripgrep for `writeModeArgs|write_mode_args|WriteModeArgs` over **all**
`*.go` files under the repository root — a superset of the drafter's non-test scope, test files
included. Zero hits. Owner of the underlying claim is claude-1 (`round-01/claude-1.md` item 4);
I am a non-owner, and this re-confirms my round-1 verdict (`round-01/kimi-1.md` item A.4) with a
fresh run.

**(a2) The CLI never reads the JSON config file — no verdict from me; I co-own this claim**
(`round-01/kimi-1.md` refinement, lines 99–103; also `round-01/codex-1.md` row "JSON config as a
Parley runtime source"). §15.1 bars my verdict. Raw evidence, reported without classification: a
repo-wide ripgrep for `headless-agents|headless_agents|headlessAgents|HeadlessAgents` over all
`*.go` returned zero hits, and `internal/config/runtime.go:137-154` (`configLayers`) enumerates
exactly `~/.parley/agents.toml` (line 141), `parley-deck/agents.toml` (144),
`parley-deck/agents.local.toml` (145), and `$PARLEY_HEADLESS_AGENT_CONFIG` (147–151) — no JSON
path. The standing non-owner verdict is hermes-1's (PRIMARY, `round-02/hermes-1.md` @codex-1:
"`configLayers:137-154` reads `agents.toml`/`agents.local.toml`, never the skill's
`headless-agents.local.json`"). The claim rests on that, not on the drafter's or my say-so.

**(b) Decks setting the field — CONFIRMED for the one deck inside this repository (PRIMARY);
the aggregate is UNVERIFIED by me.** This repo's own deck,
`parley-deck/meta/headless-agents.local.json`, sets `writeModeArgs` for all four roster agents:
claude-1 (lines 28–31), codex-1 (78–81), hermes-1 (109–112), kimi-1 (154, empty array) —
consistent with the drafter counting it among the 11. It is **not** an exclusive case: for
claude-1, codex-1 and hermes-1 the same enabling flags are duplicated in `headlessArgs`
(lines 18–25, 68–75, 103–106) — consistent with its exclusion from the 8. The other 10 named
decks, the 23-deck denominator, and the `igm-app` exclusivity example are **outside this
repository**; my mandate for this signoff is limited to this repo, so I did not read them and
say so rather than guess. This does not weaken the closure: deletion rests on (a) — no loader,
validator, or consumer exists — and on the structural teaching argument; (b) illustrates that
argument but does not carry it. A substantially wrong count would not change the decision.

**Why the closure is legitimate under §15.3.** codex-1's minority position had one empirical
premise: keeping the field "avoids a schema bump for local configs that already set it"
(`round-02/codex-1.md`, position changes and lines 127–128). The measurement shows there is no
schema, loader, or consumer, so the cost that premise protects against does not exist; and the
majority's teaching argument, previously asserted, is now measured. codex-1's remaining ground —
a labelled manual branch detoxifies the field — is a design position, not a checkable claim, and
the draft engages it with argument (hermes-1's copying-facilitator point, `round-02/hermes-1.md`
@kimi-1) plus the measured exclusivity. That is resolution by evidence and argument, not by the
3-to-1 count; the dissent and the reopen condition are recorded in the draft. Closure accepted.

The codex-1-specific question is N/A to me: I was in the majority of three, having withdrawn my
round-1 "keep the split" (`round-01/kimi-1.md` item C, lines 174–181) in round 2
(`round-02/kimi-1.md`, position change 1) on argument, before any measurement existed. The
measurement neither defeats nor outflanks a position I no longer held; it corroborates the
withdrawal's reason (a).

#### Reservation R1 (non-blocking) — the `prepack` clause of decision 4 is imprecise as drafted

The draft has the equality assertion "running under both `npm test` and the existing `prepack`
lifecycle" (consensus.md line 82). Per hermes-1's PRIMARY finding, which I did not independently
re-run and cite as SECONDARY (`round-02/hermes-1.md` @codex-1; also `round-02/codex-1.md`):
`package.json:66` `prepack` runs only `node scripts/build-addon-manifest.js --check`, which reads
neither version field and does not invoke `node --test`. Landing the decision as intended
therefore requires either editing the existing `prepack` command or adding the check inside
`build-addon-manifest.js` (hermes-1's option b). Both are edits to existing files, so "no new
script, job, or checker" survives — but the draft records only the manifest-*test* half of the
gap (consensus.md lines 88–90), not the `prepack`-invocation half. Implementation must close that
explicitly, or a direct `npm publish` still bypasses the guard — the exact path that produced the
drift. The decision's intent is unambiguous, so this does not block.

#### Drafter position changes (§15.5 ratification)

Changes 1–4 match the round files I read (PRIMARY on the comparison): change 1 against
`round-01/claude-1.md` concern 1 and `round-02/claude-1.md` SELF-CORRECTION 1; change 2 against
SELF-CORRECTION 2; changes 3 and 4 against `round-02/claude-1.md` @codex-1 / @kimi-1. Change 5 is
drafter testimony about unfiled private reasoning: I ratify that it is filed as §15.5 requires;
its content I cannot check and do not confirm. Completeness: comparing `round-02/claude-1.md`
"Current proposal" (six items) against the draft's six decisions, the only deltas are changes 3
and 4, both filed — the table is complete to my reading (PRIMARY).

#### §15.6

The trigger analysis is accurate (PRIMARY, against the round files): round 1 was not unanimous —
hermes-1 opposed the guard (`round-01/hermes-1.md` item 3) — so no steelman section was required.
The §15.6(b) record (correlated-agreement caveat, what-would-have-to-be-true, one-family note)
is present at consensus.md lines 193–211.

I am the author of the adopted `:251` text and the opencode row wording; per §15.1 I issue no
verdict on their quality. My acceptance of the draft is a signoff decision, not a
self-verification.

#### Scope declaration

Read in full for this signoff: `COOPERATION.md` §15 (lines 1176–1316); `00-prompt.md`; all four
round-01 and all four round-02 files; `consensus.md` (all 215 lines). Executed fresh in this
repo: two repo-wide ripgreps over `*.go` (patterns and results above); read
`internal/config/runtime.go:130-159`; read `parley-deck/meta/headless-agents.local.json` (all 179
lines). **Not done:** read any deck outside this repository (the other 10 decks in claim (b), the
23-deck denominator); read the skill repo (`package.json:66`, the test harness — relied on the
named SECONDARY source for R1); ran no test suites and no binaries; re-ran no git diffs (the
bundled-`references/COOPERATION.md` currency claim is facilitator-owned and no element of this
signoff depends on it). I ran no git write commands and edited no file other than this signoff.
