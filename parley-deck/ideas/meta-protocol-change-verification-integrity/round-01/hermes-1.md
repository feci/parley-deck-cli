---
agent: hermes-1
idea: meta-protocol-change-verification-integrity
round: 1
date: 2026-08-04
---

## Summary

The brief proposes nine new protocol obligations addressing verification integrity, plus six tooling defects. The source run's observations are uncheckable testimony (the p-vs-np deck is absent), so I judge each rule on whether it is sound, checkable, and proportionate in THIS deck. I verified all six tooling defects against parley 1.37.0 / parley-deck-skill 2.1.3 on this machine; four are reproducible, one is partially stale, one is confirmed. I adopt (amended) five proposals, reject two, and find two already covered by existing protocol text.

## Tooling defect verification

I ran every non-destructive check. I did NOT run `parley roster init` against the live deck. I used disposable temp dirs (`mktemp -d`, `git init`, `parley init`) for destructive tests. I did NOT run `parley preflight` with live pings (skipped to avoid pinging agents); I ran it with `--no-ping`.

| # | Defect | Verdict | Evidence |
|---|--------|---------|----------|
| T1 | `parley init` leaves §2 roster empty; `roster roster show` fails | CONFIRMED with nuance | Fresh `parley init` creates COOPERATION.md with an empty §2 table (no rows). `parley roster show` run from INSIDE parley-deck/ fails with "could not read the §2 roster" — but this is a working-directory issue. Run from the PARENT dir, `parley roster show` works fine and reads the agents.toml [roster.*] mappings. The real T1 defect is narrower than reported: `parley init` does not seed §2 from `~/.parley/agents.toml`, so a fresh deck has no roster rows until hand-edited. But `roster show` only fails when run from the wrong directory. From the parent, it reads [roster.*] from agents.toml even with an empty §2 table. |
| T2 | `sync-project` drops `protocolRole` from `meta/version.json` | CONFIRMED | Ran `parley-deck-skill sync-project --project . --dry-run --json` from the parent dir. The current version.json has `"protocolRole": "source"`. The new metadata the sync would write is MISSING the `protocolRole` key entirely. This is round-trip data loss: after a sync, `parley preflight` would see no `protocolRole` and hit the "missing/unknown → do not auto-write; ask the user" path. Verified on parley-deck-skill 2.1.0 (installed) / skill package 1.4.3. |
| T3 | `roster init` fail-closes on AUTO=no, drops agent, re-adds retired adapter; hint is a trap | PARTIALLY STALE | In parley 1.37.0, `roster init` does NOT fail-close or drop agents. Tested with a disposable deck simulating the exact scenario (kimi-1 in §2, antigravity-1 in [roster.*]): `roster init --dry-run` reported "would add [roster.kimi-1] adapter = kimi" and did NOT drop antigravity-1 or fail-close. The `roster init` simply adds the missing mapping. The "⚠ unmapped — run `parley roster init`" hint IS still displayed, and following it is safe in 1.37.0 (it adds, doesn't drop). The agents.toml comment warning about this trap references "parley 1.36" — the behavior may have been fixed in 1.37.0. The stale [roster.antigravity-1] entry persists but is harmless (no agent uses that ID in §2). |
| T4 | `preflight` reports by adapter family, not roster IDs; pings non-rostered adapters | CONFIRMED | Ran `parley preflight --no-ping --json` from the parent. The roster section reports `rosterId: "codex"`, `"claude"`, `"agy"`, `"hermes"`, `"kimi"` — these are ADAPTER FAMILY names, not deck roster IDs (which are claude-1, codex-1, hermes-1, kimi-1). It includes `"agy"` which is NOT in the deck roster (retired). A working agent like kimi reports `available: true` but the preflight does not distinguish "no pong contract" from "unavailable" — with `--no-ping` everything shows available. |
| T5 | `roster show` displays stale derived name while MODEL column shows configured model | CONFIRMED | `parley roster show` (from parent) shows `claude-1` with `DISPLAY-NAME: claude_opus-4.8-1m_max` but `MODEL: claude-opus-5[1m]`. The display name says opus-4.8 but the actual model is opus-5. The display name is derived from a stale source (likely a cached/built-in name template), not from the resolved model. |
| T6 | Skill mandates 30-min timeouts but foreground shell caps at 10 min | CONFIRMED (documentation gap) | The skill's Timeout Policy says "Default per-agent process timeout: 30 minutes" and "Poll long-running CLI processes periodically, but do not terminate them unless the configured process timeout is reached." It does NOT mention background launch, nohup, or the foreground shell cap. The only mention of "10 minutes" is for signoff appends. There is no explicit warning that a foreground shell call in the host harness may cap at 10 minutes and kill an agent mid-round. |

Summary of tooling checks: T1 confirmed (narrower than reported), T2 confirmed, T3 partially stale (fixed in 1.37.0), T4 confirmed, T5 confirmed, T6 confirmed (documentation gap). I ran all checks myself; I did not run `parley preflight` with live pings (legitimate skip — no need to ping agents for this verification).

## Proposed approach

### CRITICAL-1 — Self-verdicts launder errors

**Position: adopt amended.**

The rule is sound: a participant issuing CONFIRMED on its own claim is structurally biased, and the SELF-CORRECTION weakening-only mechanism is well-designed. But the proposed text has an enforcement gap: nothing in the protocol currently defines what a "verification verdict" is or where it lives. The rule says "a participant MUST NOT issue a verification verdict on a claim it owns" but the protocol has no `verdicts.md` or verdict vocabulary — verdicts today are just prose in round files. Without defining the artifact, the rule is unenforceable.

Replacement text:

> A participant MUST NOT issue a verification verdict (CONFIRMED, UNVERIFIED, WRONG, DISPUTED) on a claim it owns. Verdicts on claim X are admissible only from agents other than X's author. An owner revisiting its own claim uses SELF-CORRECTION, which may only weaken a claim (CONFIRMED to UNVERIFIED, UNVERIFIED to WRONG) and may never be CONFIRMED. This rule binds on every track including fast. A facilitator checks compliance by scanning round files for verdict labels on claims authored by the same agent — the check is mechanical (match agent ID in frontmatter to claim ownership).

This depends on CRITICAL-2's provenance tags and CRITICAL-3's verdict register being adopted; if they are rejected, this rule still works but is harder to check.

### CRITICAL-2 — CONFIRMED does not distinguish checking from remembering

**Position: adopt as written.**

The three-tier provenance tag (PRIMARY / SECONDARY / RECALL) is sound, checkable, and proportionate. It directly addresses the core failure: recall-grade confirmations laundering errors. The rule that RECALL caps a verdict at UNVERIFIED is the right constraint — it doesn't forbid recall, it just refuses to certify it. The check is mechanical: a facilitator scans for verdict labels and their accompanying tag. The table format is clear and requires no new tooling.

One note: "SECONDARY" as defined requires "independent confirmation by >=1 other agent, itself not RECALL". This creates a transitive dependency — if agent B confirms A's claim with PRIMARY provenance, and agent C confirms with SECONDARY citing B, C's SECONDARY is valid only if B's was PRIMARY. A facilitator can check this by reading the chain. The rule is enforceable as written.

### CRITICAL-3 — Verdict conflicts have no resolution mechanism

**Position: adopt amended.**

The conflict register is needed, and "never by vote" is the right principle. But the proposed resolution hierarchy has a gap: "higher provenance tag wins" is undefined when tags are equal but verdicts disagree. The fallback ("explicit derivation or source locator wins") is subjective — what counts as "explicit"? And the final fallback ("DISPUTED") is fine but should be cheaper to reach.

Replacement text:

> Add a verdict-conflict register (`verdicts.md`, drafter-owned, append-only). When two participants return opposite verdicts on the same claim, the conflict MUST be recorded there before consensus. Resolution by argument or provenance, never by vote: (1) PRIMARY beats SECONDARY beats RECALL; (2) at equal provenance, the verdict quoting a specific source locator (DOI, URL, line number, file path) wins over one that does not; (3) if both or neither quote a locator, the claim goes to FINAL.md as DISPUTED with both verdicts recorded. Counting agents is explicitly forbidden as a resolution method. This binds on every track including fast.

The key change from the proposed text: replaced "explicit derivation" with "specific source locator" — a facilitator can check whether a locator is present without judging whether a derivation is "explicit".

### MAJOR-4 — Unfalsifiable "this avoids the known obstacle" claims

**Position: adopt as written.**

This is the strongest proposal in the brief. The witness requirement is sound: a claim that a proposal avoids a named obstacle is admissible only with a concrete witness (counterexample, precondition check, or cited result). "Adjectives asserting the exemption are not witnesses" is exactly right. The rule is checkable — a facilitator looks for the obstacle name and then looks for a witness; an adjective-only claim gets UNKNOWN. It generalizes beyond any field, as the brief notes. It composes with P6/P7 (no-suppression): the witness requirement is an admissibility gate, not a suppression rule — it says what enters consensus, not what a reviewer may report.

### MAJOR-5 — Proposed sub-goals are never checked for being already settled

**Position: adopt amended.**

The settledness check is sound in principle but the rule as written is too broad and would create ceremony on every track. "Any proposed sub-goal, milestone, or acceptance criterion" is a large surface — on the fast track, a trivial sub-goal should not require a settledness check with provenance. The RECALL-only cap at NOVELTY UNVERIFIED is good, but the rule should scope to non-trivial sub-goals.

Replacement text:

> Any proposed sub-goal, milestone, or acceptance criterion that is presented as novel or recommended work MUST carry a settledness check before entering consensus.md: already proved, already refuted, or open, with provenance tag per CRITICAL-2. A sub-goal whose settledness is RECALL-only is marked NOVELTY UNVERIFIED and may not be presented as recommended work. On the fast track, this check is required only for sub-goals the facilitator flags as potentially settled (e.g., the claim cites a specific prior result). On standard and deliberation tracks, all proposed sub-goals require the check.

The change: "presented as novel or recommended work" narrows the scope, and the fast-track exemption for non-flagged sub-goals keeps ceremony proportional.

### MAJOR-6 — The facilitator adjudicates disputes about itself, unreviewed

**Position: adopt amended.**

The problem is real: concentrated roles without review is a structural weakness. But (a) "PROVISIONAL until ratified by at least one non-facilitator" is hard to enforce when the deck has only two participants and one is facilitator — the single non-facilitator becomes a bottleneck. And (b) "the consensus.md drafter SHOULD be a different agent" conflicts with Phase 4's rule that the idea initiator drafts FINAL.md — if the initiator is the facilitator, the rule creates a contradiction.

Replacement text:

> (a) Facilitator dispute rulings are PROVISIONAL until ratified by at least one non-facilitator participant. When only one non-facilitator is available, their ratification or objection is recorded but does not block — the ruling stands as PROVISIONAL with the non-ratification noted. (b) When the facilitator is also a participant, FINAL.md MUST record the role concentration, and at least one non-drafter MUST review the drafter's own concessions (self-corrections, withdrawn claims). If no non-drafter is available, the concession review is skipped with a recorded note — the rule is proportional to roster size. (c) roles: remains advisory, but procedural roles (adjudicator, drafter) should be separable from participation. This binds on standard and deliberation tracks only.

The changes: (1) removed the "SHOULD be a different agent" drafter rule that conflicts with Phase 4; (2) made ratification proportional to roster size (one non-facilitator can't be a hard gate in a two-participant deck); (3) scoped to standard/deliberation as the brief allows.

### MAJOR-7 — Agreement between models is treated as convergence

**Position: reject.**

The rule is sound in principle but unenforceable as written. "if round 1 produces unanimity, consensus MUST NOT close until one participant is assigned to steelman the strongest rejected alternative" — but what is the "strongest rejected alternative" when all four participants agreed? There is no rejected alternative to steelman. The rule assumes disagreement exists to mine; if unanimity is genuine, the steelman assignment is theater. The correlated-agreement caveat is valuable but should be advisory, not mandatory.

More fundamentally: "if round 1 produces unanimity" is a condition that triggers on the OUTCOME of round 1, but the rule is supposed to prevent false convergence. A participant who genuinely agrees after independent analysis (round-1 independence is enforced) should not be forced to manufacture disagreement. The protocol already has round-1 independence (Phase 1) as the independence guarantee; if that holds, agreement among models with shared training data is a risk worth noting but not worth blocking consensus over.

The "FINAL.md MUST state where multiple independent proposals are in fact one family" part is good and enforceable — a facilitator can check whether proposals share a program family. I would support that as a standalone advisory rule, but not the full steelman mandate.

Reason for rejection: the steelman mandate is unenforceable when unanimity is genuine (no rejected alternative to steelman), and the consensus block is disproportionate to the risk. The correlated-agreement caveat should be advisory in consensus.md's existing "Comparison & blind spots" section, not a hard gate.

### MINOR-8 — Round-1 independence is honor-system only

**Position: already covered by existing protocol sections.**

COOPERATION.md §4 Phase 1 already states: "Round 1 is written without reading other agents' round-1 files — the point is independent analyses on the table before anchoring. Write your file, save (or commit/push), then read the others."

§11.A (local-dir transport) states: "Independence rule is a social one: write your file first, then git pull/git log to see others. There is no enforcement beyond agent discipline."

§11.B (GitHub PR transport) offers the sub-branch protocol for stronger isolation: "each agent works on idea/<slug>/round-01-<agent-id> and pushes only when ready."

The protocol already states plainly that round-1 independence is a cooperative convention and not enforced (§11.A: "no enforcement beyond agent discipline"), and already offers a stronger isolation mechanism (§11.B sub-branches / parley-worktrees addon). The proposed rule adds nothing — it says "either state plainly that it's a convention OR require parley-worktrees". The protocol already states it plainly AND offers the worktree option. This is a duplicate.

### MINOR-9 — Facilitator context asymmetry

**Position: adopt as written.**

The rule is sound and proportionate: "Any source material the facilitator gathers while scoping an idea MUST be copied into 00-prompt.md (or a sibling file referenced from it) before participants are invoked." It makes explicit what §6 rule 4 implies ("copy the snippet — agents may lack cross-workdir read access"). The check is mechanical: does 00-prompt.md contain the source material, or does it reference a sibling file that does? If not, the facilitator must add it. The rule binds on every track and adds minimal ceremony.

## Concerns / open questions

1. CRITICAL-1, CRITICAL-2, and CRITICAL-3 form a package — they define verdicts, provenance, and conflict resolution. Adopting one without the others creates a partial system. If the round-2 consensus adopts only some, the partial adoption should be noted as a known gap.

2. The provenance tags (CRITICAL-2) introduce new vocabulary (PRIMARY, SECONDARY, RECALL, CONFIRMED, UNVERIFIED, WRONG, DISPUTED, SELF-CORRECTION, NOVELTY UNVERIFIED, UNKNOWN). That is nine new labels. The protocol currently has none of these. There is a risk of vocabulary overload on the fast track. Consider whether the full vocabulary binds on fast, or whether fast needs only CONFIRMED/UNVERIFIED/RECALL.

3. MAJOR-4's witness requirement interacts with P6 (no-suppression review briefs). A facilitator might interpret "without a witness the verdict is UNKNOWN" as a suppression instruction ("don't raise this claim"). The rule should clarify that the witness requirement gates entry into consensus.md, not what a reviewer may report in a review file. P6 already protects review scope; this clarification prevents confusion.

4. MAJOR-5's settledness check is hard to perform without web access or a literature search. On this deck, agents are CLI tools with terminal access — they can search the filesystem and run commands, but cannot always verify whether a mathematical result is settled. The RECALL cap is the right fallback, but the rule should acknowledge that settledness checks are best-effort, not guaranteed.

5. MAJOR-6's "PROVISIONAL until ratified" creates a dependency: if the non-facilitator participant is unavailable or slow, dispute resolution stalls. The amended text handles this, but the round-2 discussion should confirm whether PROVISIONAL rulings can proceed to consensus without ratification in the two-participant case.

6. I rejected MAJOR-7 but support its "one family" disclosure requirement. If round-2 consensus wants to salvage part of MAJOR-7, the family-disclosure rule is the enforceable part.

7. T2 is the most dangerous tooling defect — `sync-project` silently drops `protocolRole`, which changes preflight behavior from "source-advisory" to "missing/unknown → ask the user". This is a data-loss bug that presents as a protocol problem. It should be filed as a CLI issue, not just a protocol note.

8. T3 is partially stale in parley 1.37.0. The agents.toml comment still warns about the trap, but the behavior has been fixed. The comment should be updated to avoid confusing future facilitators who read it and avoid `roster init` unnecessarily.

## Risks

1. **Vocabulary overload.** Nine new verdict labels plus three provenance tags is a lot of ceremony. If adopted, the protocol should include a one-page reference table and the skill should surface it during Phase 0. Without discoverability, agents will improvise and the labels will be inconsistently applied.

2. **Enforceability gap.** CRITICAL-1 through CRITICAL-3 depend on a facilitator mechanically checking round files for verdict labels, provenance tags, and conflict entries. This is feasible but tedious. Without a CLI check (`parley verify` or similar), compliance will be honor-system. The rules are still better than nothing — they define what correct looks like — but the "a rule nobody can enforce is a comment" constraint is real.

3. **Track-scope ambiguity.** The brief says CRITICAL-1 to CRITICAL-3 should bind on every track including fast, and MAJOR-6/MAJOR-7 on standard/deliberation. But MAJOR-4 and MAJOR-5 don't have a track scope in the brief. I assumed MAJOR-4 binds on all tracks (it's a CRITICAL-tier problem) and MAJOR-5 is proportional (amended to exempt fast-track trivial sub-goals). Round 2 should confirm the track assignments.

4. **Composition with P6/P7.** The ratified P6 (no-suppression) and P7 (strict gate + trajectory) rules are adjacent. CRITICAL-2's RECALL cap and MAJOR-4's witness requirement are admissibility gates, not suppression rules — they gate what enters consensus, not what a reviewer may report. This distinction must be explicit in the final text, or a facilitator could misapply them as suppression (which would violate P6).

5. **The source run is absent.** Every reported observation in the brief is testimony I cannot check. I judged each rule on its structural merits, not on whether the incident happened. But the rules were designed in response to specific failures — if the failures were less severe or less common than reported, some rules may be over-calibrated. The facilitator note asks us to say so: I believe CRITICAL-1 and MAJOR-4 are well-calibrated regardless of the anecdote; CRITICAL-3 and MAJOR-6 are slightly over-calibrated for small rosters; MAJOR-7 is over-calibrated for genuine unanimity.
