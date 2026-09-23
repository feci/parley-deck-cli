---
agent: claude-1
idea: meta-protocol-change-lean-organizer
round: 2
date: 2026-09-23
responding-to: [kimi-1/round-01, zcode-1/round-01]
---

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "packet_sha256": "12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase2-deliberation-12e4b31cd3f6c1066a97caa7c7cb84213e16429dd43b26de06d4b1c32977aa18.md"}
```

Read in full (1,380 lines). I recomputed `shasum -a 256` over the packet body: it matches both
hashes above, and equals the live `parley-deck/COOPERATION.md` (108,400 B). All new measurements in
this file were taken this session at CLI-worktree HEAD `4004a99`.

## Position changes since prior round

Seven changes, five of them because I measured something after round 1 that my round-01 file had
either estimated from a hand-built table or missed entirely. Each names the round-01 claim it
replaces, per the §15.1 self-correction discipline.

**PC-1 — C1 is worth more than I said. I was measuring the wrong baseline.**
Round-01 (`## Proposed approach`, C1): *"The audience dimension is a real but **modest increment**
on top of it"*, derived from my own section table. I have now measured the shipped optimizer's real
output instead of modelling it. At `--phase 2 --track deliberation` the full body is **109,167 B**
and `--optimize` returns **65,516 B** — but that optimized body still contains, verbatim headings
present, `## 1.`, `## 2.`, `## 3.`, `## 8.`, `## 10.`, `## 12.`, `## 13.`, `Appendix A` **and
`### 11.C`** (GitLab, on a `github-pr` deck). None of those is in the protocol's facilitator reading
set. Their source sizes sum to 31,398 B; excluding §2, which I argue below should stay, 27,420 B.
So a facilitator view can plausibly reach **≈38 KB** where the shipped optimizer stops at 65,516 B.
That moves me toward @kimi-1 and @zcode-1 and away from my own round-01 caution. Stated honestly:
**38 KB is a floor, not a forecast** — `--optimize` evidently also cuts *within* sections, so
removing a section saves *at most* its source bytes, never more.

**PC-2 — The "does the facilitator view include §15?" question is already answered by shipped code,
not by our preferences.** @kimi-1 asked it as an open question and leaned "plus §15"; @zcode-1's set
omitted it; my round-01 open question 2 circled the same ground. `internal/protocolpacket/applicability.go:208-215`
puts `## 15.` and `### 15.1`–`### 15.4`, `### 15.7` in the never-cut list at `kernel` phases
`{1,2,3,5,6,7,8}`, and `### 15.5`/`### 15.6` at `fullFifteen` `{1,2,3,6,7}` (`:190-191`). Since the
audience dimension may only omit what is *already* omittable, §15 is **not omittable at phase 2 for
anyone**. The right resolution is mechanical: the role map must be unable to express the omission,
and `packet check` must fail if it tries. No vote required.

**PC-3 — `parley wait`'s job is smaller than I described.** Round-01 (alternatives table) said the
digest's *"only consumer in the tree is `internal/tui/roundsummary.go`"*. That is true of consumers,
but I under-described persistence: the digest is already written to
`parley-deck/runs/<run-id>/events.jsonl` as a `round.digest` event with the whole JSON embedded. I
parsed this run's copy — per-agent keys are exactly `agent`, `position`, `fell_back`, `present`.
So B's delta is now precisely four things, not a build: add `path`, `owner`, `validity`; extend past
design rounds; block; print to stdout.

**PC-4 — I withdraw my proposed "exit 4 on any invalid artifact" and replace it with a sharper
rule.** Round-01 concern 6 proposed `parley wait` exit `4` whenever any artifact is invalid.
@kimi-1 folds validity into the wait *condition* instead, which subsumes my intent — but produces a
worse failure: an invalid artifact means the condition never holds, so the organizer blocks for the
full timeout. I verified the concealment that makes this expensive: `parley status --idea
meta-protocol-change-lean-organizer --json` is 1,789 B and carries **no `path`, no `owner`, no
validity field**; `state.agents[]` reports `"state": "pending"` for all three participants while
`round_status` is `"completed"`. New proposal in `## Current proposal` B-3.

**PC-5 — I concede the `--role` flag name to the 2-1 majority and narrow my objection to the emitted
JSON key, where the collision is real.** Round-01 said *"A new flag must therefore not be called
`--role`"*. The flag is fine. What is not fine is the *attestation* both peers propose to add a
top-level `role` field to: `parley protocol packet --json` already emits
`source.role: "source"` with authority text reading `"(protocolRole: source)"`. A top-level `role`
would sit one level above a different `role` meaning something else.

**PC-6 — Correction to my §15.6(a) alternatives table: I missed `WriteHandoffPacket`.** @zcode-1
located `internal/runner/handoff.go`; I verified it (`HandoffPacket{InvocationID, Dir, PromptPath,
InstructionsPath}`). My round-01 alternatives row for "session resume state" named only
`internal/sessionstore` and should have named this too. Having read it, I do not agree it is the
thing to extend — see `### @zcode-1` — but the omission was mine and the correction is recorded.

**PC-7 — I can now answer the timeout question with a shipped mechanism instead of a constant.**
Round-01 claim 16 recorded the 30-minute prompt-cache lifetime as `UNVERIFIED`/`RECALL` and proposed
"configuration with a conservative default". I have since found the configuration that already
exists. See `## Current proposal` B-4; this is the organizer's explicit round-02 question and I
believe it is now fully resolved without any provider-TTL claim.

## Responses to others

### @kimi-1

**Where you are right and I am adopting your position.**

*The driver-prompt repair is the sharpest single finding in round 1, and it is in scope.* You
located `buildConsensusDraftPrompt` (`internal/app/driver_consensus.go:112-131`) instructing the
drafter toward `## Trade-offs accepted` / `## Dismissed findings` while the protocol's Phase-3
template requires `## Agreed trade-offs` / `## Open items deferred to implementation` /
`## Comparison & blind spots`, and omitting §15.3 `## Verdict conflicts`, §15.5 `## Drafter position
changes` and §15.6 `## Alternatives disposition` entirely. Neither @zcode-1 nor I found this. I
adopt it, and I want the scope argument recorded explicitly because it is the one place where
"fixing a bug" could be mistaken for "adding an obligation": **it adds no duty to any agent.** §15.5
and §15.6 already bind on every track (`## 15.7` per-track table, all three columns `yes`). The
driver currently prompts drafters *away* from duties the protocol already imposes. Repairing it
makes A safe — routing the organizer through the driver must not degrade protocol compliance as the
price of saving organizer tokens. `00-prompt.md` licenses it in terms: *"Preserve §15
dispute/provenance/alternatives duties… the old driver prompt may omit them."*

I adopt your `protocol.RequiredConsensusSections` constant mirroring `RequiredFinalSections`
(`:142-167`), and add one criterion: a **parity test** asserting the constant used to build the
prompt is the same value the consensus gate validates against, so prompt and gate cannot drift apart
again. A constant that only the prompt reads would reproduce the defect a release later.

*I also adopt:* your `unparsed`-visible failure mode; your point that the role map should live in
protocol-governed config so changing it is itself a §7 change; and your framing that the skill copy
drift check is tooling hygiene rather than a protocol obligation.

**Where I disagree, with counter-proposals.**

**K-1 — SKILL.md core ≤ 25 KB.** Counter: **≤ 20,000 B**. You marked this "to be fixed in
consensus", and @zcode-1 and I independently proposed ≤ 20 KB, so I am recording the convergence
rather than arguing it. The reason to prefer the tighter bound is that the core's whole purpose is
to be the post-compaction re-read; at 25 KB it competes with the ~38 KB facilitator packet for the
same budget. If 20 KB proves impossible without dropping a rule, the honest move is to raise it in
FINAL with the measurement attached — not to pre-budget the slack.

**K-2 — validity inside the wait condition, with timeout as the only exit for an invalid artifact.**
Counter: keep validity in the condition (you are right that it belongs there), but add a distinct
early return. Your design and this idea's own history collide: when two round-01 files carried
suffixed `## Existing alternatives (...)` headings, the driver waited rather than reported. I treat
the organizer's note of that as **operational testimony, not a verdict** — but I verified the
mechanism that makes it expensive myself (PC-4: `status --json` has no validity field at all). So:
`parley wait` returns **immediately, with a distinct exit code**, when every expected artifact is
present and at least one fails the validator, printing the validator's error string verbatim next to
the owning agent and path. An organizer that burns a 20-minute wait on a heading typo is paying
exactly the tokens B exists to save.

**K-3 — a driver-written handoff record under `runs/`.** Counter: the organizer's **read path** must
be a recomputed command, not a file. I have no objection to the driver *logging* the same bytes
under `runs/` — §3 already makes `runs/` non-canonical driver state, and you were careful to place
it there. My objection is narrower than my round-01 wording ("no hand-written handoff file")
suggested, and I will state the line I actually care about: (i) no agent may be **obliged** to
author or refresh it, or it becomes a new mandatory obligation and trips the owner's stop rule;
(ii) when the persisted record and the recomputed state disagree, the recomputed state wins, stated
in the schema doc — which is your own "`parley status` recomputation remains the authority", so I
think we already agree and only the emphasis differed. With those two, I withdraw the objection.

**K-4 — "Default timeout 25 min — below the 30-minute provider prompt-cache lifetime named in the
approved brief."** Counter: derive the default from `[defaults.timeouts]`, which already ships at
10–20 minutes. Your attribution is careful — you say the brief *names* the figure, which is true
and is not a provider claim. But a shipped 25-minute constant is longer than every timeout the
toolchain currently configures, and it would encode a number none of us can source. Detail in
`## Current proposal` B-4.

**K-5 — "no other round-01 file existed when I wrote this."** Noted and accepted; same for me. One
non-material correction for the record rather than a §15.3 conflict: your telemetry row says
"claude / codex / opencode and a zcode envelope", which matches what I find; @zcode-1's says
"claude, codex, zcode and no kimi case". `grep -n 'case "' internal/telemetry/usage.go` returns
exactly three: `claude:230`, `codex:252`, `opencode:261`. The material claim — **no kimi case** — is
agreed 3-0 and is what D3 acts on, so I am not escalating this to `DISPUTED` (§15.1 materiality:
the fix is "add kimi" under every version of the list).

### @zcode-1

**Where you are right and I am adopting your position.**

*You found `WriteHandoffPacket` and I did not* (PC-6). *You caught the study's locator error*
independently — `COOPERATION.md:34`, not `:26-33` — which @kimi-1 caught too and which my round-01
cited correctly without flagging the study's version; three independent agreements that the study's
locator is imprecise is worth carrying into FINAL so the next reader does not re-derive it.

*I adopt:* fail-closed `unparsed` over any guess; the digest's next-action line being a **fixed
enumeration**, never generated prose (this is the guardrail that keeps the digest from drifting into
the rejected summarization pattern, and it is stronger than anything I wrote); your argument against
making `status` blocking because it is "read-only fast everywhere and is scripted by other tools";
your **Windows CI leg** for `wait` and `usage`, which I did not raise and which matters because the
release ships Windows assets; and your bounded/size-capped scanning rule with no credentials in the
ledger.

**Where I disagree, with counter-proposals.**

**Z-1 — AC-C1's "≤ 60 KB" is not a binding criterion, and a byte target alone is the wrong shape.**
This is my main disagreement with your file. Measured this session: the shipped `--optimize` at
`--phase 2 --track deliberation` already returns **65,516 B**. A "≤ 60 KB" criterion is therefore an
~8% improvement on an existing flag, and — worse — it can be **passed by cutting the wrong 5.5 KB**.
Counter-proposal: make the binding criterion a **named omission set plus a named retention set**,
with the byte ceiling secondary:

- *Omits* (verified absent from the body): `## 1.`, `## 3.`, `## 8.`, `## 10.`, `## 12.`, `## 13.`,
  `Appendix A`, `### 11.C` — measured source total **27,420 B**.
- *Retains verbatim*: Quickstart, §4, §5, §9, §11.B, §2 (see N-1), and the entire never-cut floor
  including `## 15.` through `### 15.7`.
- *Omission index* names every omitted block with its trigger.
- *Secondary ceiling*: ≤ 45,000 B at phase 2 / deliberation / github-pr. I decline to propose
  anything below ~38 KB, because 65,516 − 27,420 = 38,096 B is the floor reachable by whole-section
  omission and I will not write an acceptance criterion the arithmetic cannot meet.

**Z-2 — `parley usage record` attributing by "slug mentions in the session file".** Counter:
attribute by **explicit `--idea` and `--phase` arguments plus run-record timestamp windows**, and
never by scanning session text for a slug. Two reasons. First, correctness: a slug appears in a
session for many reasons — being discussed, being grepped, being pasted into a brief — and a
228 MB rollout gives that heuristic enormous surface to mis-attribute over. Second, this run has
already produced the failure in a milder form: `organizer-usage.md` on disk now carries an
accounting correction stating that an ambiguous source selector picked the wrong rollout and that
the affected snapshots are superseded. I cite that as **operational testimony from the ledger
artifact, not as a code verdict** — but it is the same class of defect, found the same day. Keep
your `ambiguous` label for the residue that windows cannot resolve; that part I adopt unchanged.

**Z-3 — `parley status --idea <slug> --organizer` as the bootstrap surface.** Counter: put the
bootstrap in a separate verb. Your own Z-argument against blocking `status` applies with equal force
to adding an organizer-shaped output mode to it: `status` has a stable contract other tools script.
This is a small disagreement and I will not block on it; if consensus prefers the flag, the property
I need preserved is that `status`'s existing JSON shape is byte-unchanged when the flag is absent.

**Z-4 — the §15 omission (implicit).** Your facilitator set is "Quickstart, §4, §5, §9, and the
active §11 transport", which does not name §15. Per PC-2 the never-cut floor already forbids
omitting it at phase 2, so I read this as an omission of mention rather than a proposal to cut it —
but since your AC-C1 lists exactly which sections must appear verbatim, please add §15 to that list
explicitly. If it is left implicit, a future role-map edit could try to omit it and the only thing
standing in the way would be a check nobody wrote a test for.

**Z-5 — `--role` in the attestation.** Same counter as K-5's sibling issue: see PC-5 and C-2 below.

## New concerns / questions

**N-1 — Should §2 (roster) be in the facilitator view?** The Quickstart table gives the facilitator
the reading set "Quickstart, §4, §5, §9, your §11 transport" but gives it the *duty* "keep the
roster" in the same row. §2 is 3,978 B. My position: **keep §2**, because a facilitator that cannot
read the roster authority rules cannot discharge a duty the same table assigns it. This is the one
omission I would contest in any facilitator map, and it is why my Z-1 omission set excludes §2. I
would like both of you to say yes or no explicitly, because this is the sort of gap that is invisible
until an organizer silently does the wrong thing with `agents.toml`.

**N-2 — I do not fully understand what `--optimize` actually cuts, and neither peer's design says.**
Measured: at phase 2 the optimized body is 65,516 B, yet `### Phase 5 —` is **present**, and so is
`### 11.C` on a `github-pr` deck. So the optimizer is not doing whole-section phase/transport
scoping the way all three of our round-01 files assumed; the 43,651 B it removes comes from
somewhere I have not located. Consequence for the design: **the audience dimension's interaction
with the existing optimizer must be measured during implementation, not assumed**, and AC-C1 should
be stated against a *measured* `--optimize` baseline captured in the same test run rather than
against a number any of us wrote down in round 1. I flag this as my own incomplete understanding
rather than as a defect claim about the optimizer.

**N-3 — Should `wait` return early on a new unanswered `to-user` escalation?** @zcode-1 raised this
as its Q6 and left it open; @kimi-1 did not address it. My answer: **yes, with its own exit code.**
An organizer blocking 20 minutes on a round that is already stalled behind a question only the owner
can answer is the precise waste B exists to remove, and §14's human brake argues for surfacing the
escalation sooner, not later. `wait` still only *observes* — it must never advance a phase; that
boundary I hold from round-01 concern 4 unchanged.

**N-4 — FINAL should carry an explicit stop-rule ledger.** The owner's stop rule fires if FINAL drops
an A–D item *or adds a mandatory obligation beyond A–D*. Several things in our three files are
genuinely tooling rather than obligations — the skill-copy drift check, the `RequiredConsensusSections`
parity test, the Windows CI leg, the digest determinism test. I propose FINAL contain a short table
listing each such item with an explicit `obligation: no — repo test` / `obligation: no — optional
field` column, so the organizer's pre-Phase-5 inspection is a lookup rather than a judgment call.
This is a drafting discipline, not a new requirement on anyone.

**N-5 — Drafter and implementer.** I am **not claiming** either role. The organizer records kimi-1 as
the default first round-01 writer; I note only that the mechanical record does not establish
ordering (the run's `round.completed` event carries `"reconstructed": true`, and all three round-01
files landed in one commit, `4004a99`), and I do not contest the organizer's operational statement —
it is not decisive for me either way, since I am not claiming. **I am willing** to draft FINAL or to
implement if the quorum prefers it, and would post the normal `inbox/` claim at that point rather
than asserting it here. If neither of you volunteers, the Phase 4 fallback applies unchanged.

## Current proposal

Unchanged in shape from round 1: four thin mechanisms over machinery that mostly already exists. The
deltas below are what cross-review changed. **Every item maps to A, B, C or D; nothing here adds a
mandatory obligation to any agent.**

**A — organizer neither implements nor verifies code.**
A-1. Optional `facilitator: <agent-id>` in `00-prompt.md`; absent → byte-identical behaviour to
v1.48.0 (regression test). A-2. `parley preflight` fails closed when `facilitator:` also appears in
`participants:` unless `facilitator_participates: true`. A-3. The driver resolves implementer,
drafter, reviewer and goal-done checker from participants only, never the declared facilitator, and
**escalates rather than silently falling back** to facilitator implementation (@zcode-1's wording,
adopted). A-4. **Repair `buildConsensusDraftPrompt` from a `protocol.RequiredConsensusSections`
constant, with a parity test against the gate** (@kimi-1's finding, adopted; adds no agent duty —
§15.5/§15.6 already bind on every track). A-5. One permissive sentence each in §4 Phase 5, Phase 6
and the Quickstart facilitator row; one line in §9.0 for the preflight check.

**B — `parley wait` + `PhaseDigest`.**
B-1. `parley wait --idea <slug> --for round|consensus|review|implementation|any [--timeout D]
[--json]`, blocking on the existing append-only event log. B-2. The digest gains exactly `path`,
`owner` (frontmatter `agent:`) and `validity` (shipped validator verdict verbatim) per agent, and
extends past design rounds; everything else already exists and is already persisted as a
`round.digest` event (PC-3). Determinism golden test; a structural test that no field is
model-written; `unparsed` never a guess (@zcode-1); fixed-enumeration next-action line, never prose
(@zcode-1). B-3. **Exit codes:** `0` condition met; `2` timeout (partial digest printed, outstanding
agents named); `3` **all expected artifacts present but ≥1 invalid — returns immediately with the
validator error verbatim** (PC-4 / K-2); `4` new unanswered `to-user` escalation (N-3); `1` usage or
IO error. B-4. **Timeout — the organizer's question, answered.** `[defaults.timeouts]` already ships
in `~/.parley/agents.toml` and is already parsed (`timeoutsBlock`, `internal/config/runtime.go:288`):
`signoff_ms = 600000` (10 min), `round_ms`/`review_ms`/`deep_reasoning_ms = 1200000` (20 min).
`parley wait --timeout` therefore **defaults from the matching existing key** — `round` → `round_ms`,
`consensus` → `signoff_ms`, `review` → `review_ms` — so:

> **Yes: the shipped default is a bounded wait of 10–20 minutes, which is below 30 minutes, and it
> is below it without anyone claiming a provider-wide cache TTL.** The bound comes from configuration
> that already exists, not from a new constant. For belt and braces I propose a hard ceiling
> (`--timeout` values above the §4.0 per-track agent timeout are rejected; §4.0's deliberation column
> already reads "~30 min"). The 30-minute figure remains what it was in my round-01 claim 16:
> `UNVERIFIED`/`RECALL`, sourced to the owner's brief as a *requirement*, never asserted as a
> provider fact — and no document, help text or comment we ship may state it as one.

B-5. `wait` observes only; it never advances a phase. B-6. Raw artifacts stay canonical: every digest
row carries its path, and any `❌` / `unparsed` / `DISPUTED` / adverse validity forces opening the
raw file.

**C — audience-scoped packet, slim SKILL.md, generated brief.**
C-1. `--role facilitator|participant` on `parley protocol packet` (flag name conceded to the 2-1
majority, PC-5). Facilitator-only; the participant pipeline and `context_mode` defaults are
untouched. C-2. **The emitted key is `request.role`, never a top-level `role`**, leaving
`source.role: "source"` unambiguous; a test asserts both appear in one packet with distinct values.
C-3. The role map lives in protocol-governed config so changing it is a §7 change (@kimi-1), and
**the audience dimension may only omit what phase/track/transport already omits — it can never cut
below the never-cut floor**; `packet check` proves it with a negative test. §15 is consequently not
omittable at phase 2 by construction (PC-2), and AC-C1 names it explicitly (Z-4). C-4. AC-C1 is a
**named omission set + retention set + omission index**, with ≤ 45,000 B secondary and the
`--optimize` baseline captured in the same test run (Z-1, N-2). C-5. SKILL.md split, core
**≤ 20,000 B** (K-1), relocation only, a script proving every moved heading lands in exactly one
reference, and the core naming `run`, `continue`, `wait`, `status`, `consensus`, `preflight` — the
direct inverse of the verified v2.12.1 state. C-6. `parley organizer brief --idea <slug>`: computed,
≤ 8,192 B, byte-identical across two runs, writes no file (assertable against a read-only deck).

**D — session per phase + measurement.**
D-1. The brief *is* the handoff; no participant-authored state file. The driver **may** log the same
bytes under `runs/` as non-canonical telemetry provided no agent is obliged to author or refresh it
and the recomputed state is authoritative on disagreement (K-3, converging with @kimi-1 and
@zcode-1). D-2. `parley usage <ingest|record>` — the verb name is the drafter's call; the binding
properties are: **no flag accepts a token count** (a number typed by a model is `RECALL` and
inadmissible; a number the tool parses from the client's own file is `PRIMARY` with a locator);
streaming, ≤ 1 KB stdout on a 228 MB source; the row carries idea, phase, agent, source path and
parser id; attribution by `--idea`/`--phase` plus run-record timestamp windows, **never by slug
scanning** (Z-2), with `ambiguous` for the residue (@zcode-1, adopted); idempotent re-ingest.
D-3. Add the `kimi` telemetry case from a fixture captured live, and structured-output argv only
where the adapter supports it without breaking `-p` semantics (@zcode-1's precision, adopted); if
kimi emits no usage, record `coverage: none` — never a fabricated number (both peers, adopted).

**Retained disagreements after this round.** Three, all narrow and none blocking: (i) SKILL.md core
20 KB vs @kimi-1's 25 KB; (ii) AC-C1 as a named omission set vs @zcode-1's ≤ 60 KB byte target;
(iii) usage attribution by explicit arguments vs @zcode-1's slug-mention scanning. I consider (iii)
the only one with a correctness consequence. **Open for the next round:** N-1 (§2 in the facilitator
set) — I would like an explicit yes/no from each of you.

**Extra mandatory obligations beyond A–D: none, in my proposal.** `facilitator:` is optional; the
§9/§4 protocol lines are permissive; `wait` is an alternative to polling, not a required step; the
drift check, parity test, Windows leg and determinism tests are repo tests. N-4 asks FINAL to say
so in a table so the organizer's stop-rule inspection is a lookup.

## §15 provenance

### Verdicts I issue this round

All `PRIMARY`, produced this session at CLI-worktree HEAD `4004a99` with the command or file quoted.
Round-01's claims 1–16 stand except where a PC entry above amends them.

| # | Claim | Verdict | Tag | Evidence |
|---|---|---|---|---|
| 17 | Full phase-2 packet body is 109,167 B; `--optimize` returns 65,516 B (phase 1: 65,415 B) | CONFIRMED | `PRIMARY` | `parley protocol packet --phase 2 --track deliberation [--optimize] --print \| wc -c` |
| 18 | The optimized phase-2 body still contains `## 1.`, `## 2.`, `## 3.`, `## 8.`, `## 10.`, `## 12.`, `## 13.`, `Appendix A`, `### 11.C`, `### Phase 5 —` | CONFIRMED | `PRIMARY` | `grep -qF` per heading against the captured `--optimize --print` body |
| 19 | Source sizes: §1 3,130 · §2 3,978 · §3 2,337 · §8 1,531 · §10 1,494 · §12 7,544 · §13 5,173 · AppA 2,021 · §11.C 4,190 · §15 8,056 B | CONFIRMED | `PRIMARY` | `awk 'NR>=a && NR<=b' parley-deck/COOPERATION.md \| wc -c`; boundaries from `grep -n '^## \|^### 11\.'` |
| 20 | §15 is in the never-cut list: `## 15.` and 15.1–15.4, 15.7 at `kernel {1,2,3,5,6,7,8}`; 15.5/15.6 at `fullFifteen {1,2,3,6,7}` | CONFIRMED | `PRIMARY` | `internal/protocolpacket/applicability.go:190-191, 208-215` |
| 21 | `parley status --idea <slug> --json` is 1,789 B and carries no path/owner/validity; `state.agents[]` are all `"pending"` while `round_status` is `"completed"`; `terminal: true`, `outcome: "completed"`, `attention: "DONE"` while the next action is to open round-02 | CONFIRMED | `PRIMARY` | command output quoted in full this session; `wc -c` = 1789 |
| 22 | The round digest is persisted to `runs/<id>/events.jsonl` as `round.digest`; per-agent keys are exactly `agent`, `position`, `fell_back`, `present` | CONFIRMED | `PRIMARY` | `json.loads` over the event's embedded digest; top keys `idea, round, total, completed, lines, flag_block, flag_counter, flag_accept, flag_escalate, next` |
| 23 | `parley wait` does not exist at 1.48.0 | CONFIRMED | `PRIMARY` | `parley wait --idea x` → `unknown command: wait`; `parley --version` → `parley 1.48.0` |
| 24 | The packet JSON already emits `source.role: "source"` (authority text `"(protocolRole: source)"`) and a `request` object `{phase, track, transport, optimize}` | CONFIRMED | `PRIMARY` | `parley protocol packet --phase 2 --track deliberation --json` |
| 25 | `[defaults.timeouts]` ships in `~/.parley/agents.toml` with `signoff_ms=600000`, `round_ms=review_ms=deep_reasoning_ms=1200000`, parsed by `timeoutsBlock` | CONFIRMED | `PRIMARY` | file read; `internal/config/runtime.go:288` and the `timeoutsBlock` definition |
| 26 | `internal/telemetry/usage.go` has exactly three adapter cases: `claude:230`, `codex:252`, `opencode:261` — no `kimi`, and no `zcode` *case* | CONFIRMED | `PRIMARY` | `grep -n 'case "' internal/telemetry/usage.go` |
| 27 | `WriteHandoffPacket` exists; `HandoffPacket` is `{InvocationID, Dir, PromptPath, InstructionsPath}` — per-launch instructions, not an organizer phase record | CONFIRMED | `PRIMARY` | `internal/runner/handoff.go`, struct and function read |
| 28 | The round-2 artifact validator requires only matching `agent`/`idea`/`round` frontmatter plus ≥1 `## ` section with content | CONFIRMED | `PRIMARY` | `internal/runner/validation.go:15-45` (`ValidateRoundArtifact`, `anySectionHasContent`) |
| 29 | This run's `round.completed` event carries `"reconstructed": true`; all three round-01 files landed in a single commit `4004a99` | CONFIRMED | `PRIMARY` | `runs/20260923T202501.377412000Z/events.jsonl`; `git log -- round-01/` |
| 30 | "Timeout below the provider's 30-minute prompt-cache lifetime" | UNVERIFIED | `RECALL` | Unchanged from round-01 claim 16. No locator reachable from either worktree; provider- and plan-specific. Design consequence: B-4 derives the default from shipped config, and nothing we ship states the figure as fact |

### Testimony I transcribe but do not verdict

Per §15.1, marked as unverified testimony and not owned by me: (i) every quantitative figure from
`source-context/organizer-token-study.md` §§1–5, unchanged from round-01; (ii) the organizer's
operational notes in `organizer-notes.md` — the round-01 heading stall, the driver's cross-review
refusal, the worktree-registration gap — which I cite in K-2 and Z-2 as *observations that motivated
a design choice*, never as code verdicts. Where a design choice of mine depends on the underlying
mechanism being real, I verified that mechanism myself and recorded it above (claims 21 and 22 for
K-2; the ledger's own on-disk accounting correction is the artifact I cite for Z-2, and it is
testimony too).

### Conflicts, ownership, alternatives

**No `DISPUTED` claims.** The one factual discrepancy between participants — @zcode-1's telemetry
adapter list vs @kimi-1's and mine — is resolved by claim 26 and is **not material** under §15.1: the
recommendation (add a kimi case) is identical under either list, so no decision or acceptance
criterion turns on it. Recorded as a correction, not escalated.

**Ownership.** I own every design proposition in `## Current proposal` and issue no verdict on any of
them (§15.1). Claims 17–30 are about the state of the tree, not about whether my design is good.
I hold no facilitator role; §15.5 does not apply to me unless I later draft, and I am not claiming
the drafter role (N-5).

**§15.6(a) alternatives — delta since round 1.** My round-01 table stands, with one addition and one
correction: **added** `internal/runner/handoff.go` `WriteHandoffPacket` (per-launch handoff packet)
to the "session resume state" row, credited to @zcode-1 (PC-6); **corrected** the "deterministic
round digest" row — the digest is not only TUI-rendered, it is persisted to `runs/<id>/events.jsonl`
(PC-3), which shrinks B's build rather than enlarging it. **Scoped null unchanged:** no shipped
mechanism for a blocking phase-boundary wait (claim 23), an audience-scoped render (claim 24), an
organizer usage ledger, or a generated organizer brief. Sources consulted this round: the `parley`
dispatch surface, `internal/protocolpacket/applicability.go`, `internal/config/runtime.go`,
`internal/runner/{validation,handoff}.go`, `internal/telemetry/usage.go`, this run's event log, and
`~/.parley/agents.toml`.

**§15.4.** I make no exemption claim. The one place I touch a previously-rejected pattern is D2
(client accounting vs "self-reported usage as cost evidence"); my discriminator is mechanical and
unchanged — the ingest cannot accept a number — and Z-2 tightens rather than relaxes it.
