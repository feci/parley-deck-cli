### hermes-1

Verdict: accept

## Scope declared

Read in full: COOPERATION.md §15 (lines 1176-1316), 00-prompt.md, all four round-01 files, all four round-02 files, consensus.md (all 215 lines).

Ran for this signoff (PRIMARY):
- A repo-wide search for `headless-agents|headless_agents|writeModeArgs|write_mode_args|WriteModeArgs` across every `*.go` file in `parley-deck-cli` — zero hits. This independently reproduces the drafter's claim (a) grep.
- Read `internal/config/runtime.go:137-154` (`configLayers`) to confirm the enumerated agent-config inputs are exactly: `~/.parley/agents.toml`, `parley-deck/agents.toml`, `parley-deck/agents.local.toml`, and `$PARLEY_HEADLESS_AGENT_CONFIG`. `meta/headless-agents.local.json` is not among them.

Did not run: the skill test suite, the parley CLI, live vendor probes, or the `v1.38.0..v1.39.0` COOPERATION.md currency diff (facilitator-owned claim; no position of mine depends on it).

Could not reach: VC-1 claim (b) — the 23 workspace decks under `AI_WORKSPACE` are outside this repository. I do not verify (b) and do not guess at it.

## VC-1 closure — accepted

**Claim (a) — "the CLI never reads the file the field lives in": CONFIRMED [PRIMARY].**

My search for `writeModeArgs`/`write_mode_args`/`WriteModeArgs`/`headless-agents`/`headless_agents` across all non-test `*.go` files returned zero hits. My read of `configLayers` (runtime.go:137-154) confirms the only agent-config inputs are the three TOML layers plus the env-var path; `meta/headless-agents.local.json` is not read by `parley` at launch. This matches the drafter's measurement verbatim. The drafter both made the measurement and benefits from it (it supports the majority's delete position, which the drafter shared), so §15.1 makes this independent check meaningful — and the check passes.

**Claim (b) — "11 of 23 decks set the field, 8 exclusively": UNVERIFIED.**

Those decks live under `AI_WORKSPACE`, outside this repository. I cannot reach them from here and will not guess. I note that (b) is additive weight for the majority's "the field teaches" argument, but the closure does not stand or fall on (b) alone — see below.

**The drafter's correction (§15.5, position change #5) is sound.** The drafter withdrew the private alarm that the 11 decks would "silently lose their flags at launch." That correction is logically entailed by (a) alone: if the CLI never reads the file, nothing is dropped, because nothing is consumed. It does not depend on (b), and I can evaluate it on (a), which I verified. The two facts being stated together — "(a) nothing is read, (b) the field is widely set" — is the right framing; (b) without (a) reads as an alarm, and the drafter correctly says so.

**codex-1's position: the measurement defeats the argument, it does not outflank it.**

codex-1's ground was "keeping it avoids a schema bump for local configs that already set it" (consensus.md:118-119, quoting codex-1 round-01). Claim (a) shows there is no schema — no validator, no loader, no consumer in the CLI. Deleting the field from the documentation changes the runtime behaviour of zero decks. The cost codex-1 was protecting against does not exist; the argument's factual premise is removed, not sidestepped. That is a direct defeat.

codex-1's structural insight — that manual facilitator assembly is a different activity from CLI config, and the recipe should branch — is honored by the consensus (item 3 adopts the manual/CLI branch split). What codex-1 does not get is keeping the field inside the manual branch. That is the right call: even on the manual path, a two-field shape teaches that the write-enabling flag can live separately from `headlessArgs`, which is the exact model that made the defect invisible. The branch split accommodates codex-1's distinction; the field deletion addresses the teaching model. Both can be true.

I was in the majority (delete) in rounds 1 and 2. The measurement adds empirical weight I could not produce (I did not survey the workspace decks), but my position did not depend on (b) — it rested on (a), which I verified myself in round 1 (PRIMARY: `agentOverride` at runtime.go:97-132 has no `writeModeArgs` field) and have now re-verified at the repo-wide level for this signoff.

## The six edits

1. **opencode row** — agree. Unanimous; kimi-1's wording carries the argv detail. No reservation.
2. **SKILL.md:251 replacement** — agree. I supported the converged text in round 2 (round-02/hermes-1.md:185-208). kimi-1's version is the superset: its lead sentence is the inversion fix, its body is codex-1's contract framing, its mechanism clause is mine, and the existing confinement floor is retained. The consensus correctly replaces both sentences at :251 rather than adding a third.
3. **Manual/CLI branch split** — agree. I proposed a boundary marker in round 2; the consensus adopts codex-1's fuller split. I acknowledged in round 2 that codex-1's version "also covers 'Headless Agent Configuration' and WORKED_EXAMPLES.md, which a single marker sentence does not reach" (round-02/hermes-1.md:169). The broader split is within the range I accepted.
4. **compatibility.json — bump plus one assertion** — agree. I withdrew my round-1 "no guard" objection in round 2 on evidence (round-02/hermes-1.md:11-24). The consensus runs the assertion under both `npm test` and `prepack` (consensus.md:81-82), which closes the prepack gap I raised (round-02/hermes-1.md:80-92). My implementation-gap note — that the existing manifest test verifies hashes and does not read version fields, so the assertion needs its own home — is recorded in the consensus (consensus.md:88-90).
5. **writeModeArgs deleted** — agree, per VC-1 closure above.
6. **What stays out** — agree. Unanimous. The hermes incident narrative, promotion history, probe versions, ACP availability, kimi's exit-1 message, and CLI struct names all fail the "does a facilitator behave differently" test. `references/COOPERATION.md` untouched per the byte-identical constraint.

## §15.5 drafter position changes — ratified

The five changes in the consensus table (consensus.md:179-189) are accurate and complete from my reading of claude-1's round-01 and round-02 files. Change #5 (the withdrawn "silently lose flags" alarm) is correctly filed under §15.5 rather than dropped, and the correction is sound per (a) above.

## §15.6 correlated agreement

The consensus correctly notes the §15.6 trigger did not fire: round 1 was not unanimous (I opposed the guard), so no steelman round was required (consensus.md:198-200). The correlated-agreement caveat — that the four participants are related models and the unanimous CONFIRMED on the launch-path claim is a shared prior — is accurately stated and mitigated by the three independent `runner.go` readings and my proof of the negative (consensus.md:194-196).

No reservations. The consensus is sound, the VC-1 closure is evidence-based and independently checked on the part I can reach, and the drafter's position changes are transparent and correct.
