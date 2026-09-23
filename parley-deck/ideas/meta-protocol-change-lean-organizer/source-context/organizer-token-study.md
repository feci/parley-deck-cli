# Organizer token study — 2026-09-23

Evidence base for the idea `meta-protocol-change-lean-organizer`. Produced by six independent readers and
three independent verifiers (each recomputed the headline numbers with its own scripts). Per-agent results
are in `results/*.json`; `verify-*.json` hold the verdicts and corrections. Where a verifier corrected a
reader, the corrected figure or range is used below. Token estimates from bytes use ~3.9-4 B/token.

## 1. The cost model

The organizer pays for re-reading its own context, not for writing. In the codex-1 organizer session
(`~/.codex/sessions/2026/09/10/rollout-2026-09-10T00-01-07-01a08830-*.jsonl`, gpt-6-astra, Codex CLI 0.153.4):

| Measure | Value |
|---|---|
| model requests (2026-09-09 .. 09-23) | 3,717 at freeze; 3,728 at 19:15Z |
| input tokens / cached share | 566.7M / 95.1% |
| output tokens (share of all tokens) | 3.82M (0.67%) |
| mean context re-read per request | 152k (p50 157.6k, p90 225k, window 258.4k) |
| compactions | 50; context grows ~21.7k -> ~243k, then compacts |
| API-equivalent cost at gpt-6-astra list price | ~$1,008 (cached reads ~54%, uncached ~27%, output ~19%) |

Organizer cost = requests x resident context. Output is under 1% of tokens but ~19-27% of cost when
price-weighted (reasoning ran at xhigh on all turns while `~/.codex/config.toml` says high).

Caveat: 98.1% of this session predates the v1.48.0 commit. It measures an organizer that was also author,
FINAL drafter, implementer and release manager of evidence-first-efficiency, not a pure facilitator.

## 2. Where it went (verified, with ranges)

1. **Self-implementation and self-verification: ~30-40% of input** (49.8% is the upper bound that counts all
   code-level work). After the owner's 2026-09-15 instruction "save tokens, just manage the others", Go edits
   fell from 452 to 2, but the organizer kept verifying participants' code itself (52.8M of code-level
   requests afterwards) and the hourly input rate did not fall (12.9M/h -> 15.0M/h). Instructions alone did not
   change the rate.
2. **Post-compaction re-orientation: ~33-39% of input.** The first 10-15 calls after each compaction read
   ~250-265 KB: parley-deck SKILL.md after 49 of 50 compactions, COOPERATION.md 49 of 51, AGENTS.md, and the
   graphify, openviking and worktrees skills. Standing instructions drive it (AGENTS.md tells the agent to read
   some skills in full). Because the session is compaction-saturated, a simulation shows that cutting reads
   OUTSIDE the re-orientation window saves ~0%. The levers are the post-compaction floor and the number of
   compactions, not read volume in general.
3. **Compaction itself: 42% of all uncached input.** 50 compactions of ~240k each; 38 had 0 cached tokens.
   One was fully cached, so the uncached cost is not intrinsic.
4. **Polling: 9.1-12.2% of input.** 1,488 of 1,539 `write_stdin` polls used a 1 s yield. Polled sessions were
   the organizer's own tests/harnesses 869 (51%), participants 403 (24%), release 106. Waiting on participants
   alone cost ~9M. Most polling disappears once implementation and verification are delegated.
5. **Large outputs.** Codex truncates exec output at ~40 KB; 227 outputs at the cap hold ~41% of resident
   tool-output text (~25% of input). 496 text outputs of 20 KB or more hold 68.6% of it.
6. **Fixed prefix.** ~15.5k tokens on the first request; a 22 KB skills catalogue (79 skills, 54 `vercel:*`)
   is present in every request (inferred ~5.6k tokens/request).

## 3. What the tooling already offers and the organizer did not use

- The organizer ran `parley run` once, `parley consensus` twice and `parley status` never; it wrote 70
  participant prompts by hand. SKILL.md 2.12.1 never mentions `parley run`, `continue`, `consensus`, `status`
  or `preflight`.
- The driver already launches rounds, validates artifacts, waits and escalates on deadlines, opens
  cross-review rounds, drafts consensus/FINAL through the first headless participant
  (`internal/app/driver_consensus.go:82-95`), collects append-only signoffs and runs Phases 5-8.
- Driver gaps: the consensus-drafting prompt omits §15.3/§15.5/§15.6 and "Comparison & blind spots";
  per-track human gates are not enforced; no PR/MR mirroring; no inbox digest; release is not covered.
- State digests exist and are small: `parley status` 10.7 KB for the whole deck (vs 443 KB for all
  `00-prompt.md` files plus 213 KB of inbox); `parley status --idea X --json` 2.2 KB;
  `parley consensus status --json` 1.8 KB vs a 22 KB consensus.md.
- The protocol already defines a facilitator reading set: `COOPERATION.md:26-33`
  ("Facilitator | Quickstart, §4, §5, §9, your §11 transport"). The packet renderer has no role input, and
  SKILL.md:22-53 makes the facilitator read the full body. A pure organizer needs ~43-50 KB of the 108.4 KB.
- The protocol does not require the organizer to be drafter or implementer: drafter = author, or with
  `author: user` the first round-1 writer or a volunteer (COOPERATION.md ~L407-409); implementer = FINAL
  drafter unless claimed (~L443); facilitator calls are provisional (§15.5 ~L1340-1343); the facilitator's
  verification duty is existence, ownership and validity (~L437, ~L890). Verify line numbers at HEAD.

## 4. What 1.48.0 did and did not do for cost

- Telemetry records only runner-launched processes, has no role field, and captures almost no tokens:
  claude-1 23/26 invocations with usage, zcode-1 2/6, kimi-1 0/20 (a parser gap even with structured flags),
  hermes-1 0/7. Nothing measures the organizer.
- `parley protocol packet` defaults to full context. The optimized packet is 54.8-71.7% of 108,400 B without
  flags; it was designed for participants and has no role input. The packet experiment and 12-task pilot never ran.
- 1.48.0 prepends the full 108 KB protocol to every runner-launched prompt (~27k tokens per launch),
  including drafter and signoff launches (participant cost, not organizer cost).
- No token budget exists; budgets count launches, driver steps, cycles, USD and wall clock.

## 5. Claude organizer session (for comparison)

`5dc331bd` (2026-05-28 .. 09-23, 1M window): 13,294 requests, 5.76B context tokens, mean 433k, ~$3,870
estimated. Facilitation ~31-37% (waiting ~5-11%, reading artifacts ~6-9%, launching ~5.5-7.7%, drafting 3.6%,
protocol 2.9%); implementation 43.7%. 95 full cache rewrites after gaps of an hour or more = 57% of cache
writes. Per recent idea: 46-62 requests, 15-21M context tokens, 43-49% waiting plus artifact reading. A fresh
context per idea would have saved an estimated 37-77% (simulation).

## 6. Prior decisions to respect (do not re-propose without new evidence)

Rejected before: LLM or lossy summarization as a replacement for artifacts; tools that select which protocol
text an agent sees (cognee, graphify, vector stores, `viking://`); stored summary sidecars; the prior-round
ledger behind a constant (deleted 1.43.1); asymmetric reviewer context; random mutation; new round-opening
mechanisms; speed profiles that downgrade model or effort; silent model swaps; self-reported usage records as
cost evidence. Invariants: non-solo; §15 no self-verdicts and provenance tags; conflicts never settled by
count; raw round files never hidden behind a summary; review briefs must not suppress findings; round 1 stays
independent; English only; §14 human brake; full protocol context stays the default until the packet trial
or a §7 protocol change says otherwise (this idea is such a §7 change for the organizer role).

## 7. Levers outside parley (owner environment, not part of the release)

- Trim the organizer's skill catalogue (54 `vercel:*` skills are irrelevant to organizing).
- AGENTS.md lines that require reading some skills "in full" drive post-compaction re-reads.
