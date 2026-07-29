---
idea: readme-skill-catalogue
drafted-by: claude-1
date: 2026-07-29
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: awaiting-signoffs
---

## Summary

Round 1 produced four independent proposals; round 2 was a cross-review by codex-1 and
hermes-1 that fact-checked every truth claim and produced two competing merged drafts. Five
of the seven forks were unanimous at the end of round 1 and were not re-litigated. What
actually differed was the wording of the hook, the base for each catalogue entry, and the
accuracy of the "false claims" lists — including several accusations that turned out to be
wrong.

**The most consequential outcome of the round was not the copy.** It was that three
participants independently caught a number I invented in my own round-01 file, and codex-1
caught two more errors in the brief I wrote. Those are recorded in C9 and C10 rather than
quietly fixed.

---

## Agreed decisions

**C1 — Order: hook → catalogue → fastest-path install → use → the rest.**
Unanimous that the catalogue precedes Install. kimi-1 dissented in round 1 on behalf of the
"colleague-sent reader" who wants the `npx` line immediately; hermes-1's round-2 synthesis is
adopted: that reader Ctrl-Fs for `npx` regardless of position, so catalogue-first serves two
of three readers fully and harms none. kimi-1's underlying point is adopted in full — the
**install *apparatus*** (per-target commands, manual paths, Windows, Homebrew, updating) moves
*below* the catalogue, and only an ~8-line fastest path sits near the top.

**C2 — Prose entries, not a table.** Unanimous. A five-row table cannot carry a refusal such
as "no numeric aesthetic score" without shrinking it to a feature. A thin five-bullet scan
list above the prose entries is adopted (kimi-1) — it answers "which one is that?" and
describes nothing.

**C3 — Exactly one load-bearing position per add-on**, stated as a claim, with the details
behind the link. Unanimous.

**C4 — F4 resolves as kimi-1 framed it: satellites in dependency, peers in visual weight.**
All five get the same heading level, the core comes first, and each add-on entry states it
loads alongside the core skill. The README does not use the words "peer" or "satellite".

**C5 — The hook leads with the failure mode**, with the artifact trail as the proof that
follows. Unanimous.

**C6 — `parley-deck-cli/README.md` is out of scope.** Unanimous; recorded as a follow-up.

**C7 — Hook base: claude-1, by user ruling.**
codex-1 selected claude-1's opening; hermes-1 selected kimi-1's. Both are defensible and the
line is the most visible sentence in the package. **The drafter (claude-1) is the author of
one candidate and recused himself from the tie-break**, and the user ruled for the claude-1
base on 2026-07-29. The kimi-1 line — *"One agent's answer is a first draft. Several agents'
recorded agreement is a decision."* — is recorded here as the runner-up and is **not** grafted
in: per the design doctrine this package ships, one direction wins whole and a graft may not
be used to smuggle the loser's opening back in.

**C8 — Catalogue bases, selected per skill** (codex-1's round-2 selection, adopted):

| Entry | Base | Reason |
|---|---|---|
| `parley-deck` | codex-1 | most precise on artifact ownership, tracks, transports; does not repeat the false confidence-rating claim |
| `parley-design` | kimi-1 | compresses the doctrine into memorable refusals without reproducing PDS |
| `parley-design-check` | codex-1 | the only draft carrying the shipped exit-4 behaviour |
| `parley-tracker` | codex-1 | the only draft drawing the connector boundary correctly |
| `parley-worktrees` | kimi-1 | shortest complete account of the lock manifest and disjointness check |

Named grafts are listed as HTML comments in the shipping copy and must survive into
`FINAL.md`.

**C9 — Truth audit. This table is binding; it replaces every round-1 accusation list.**

| Claim | Verdict | Action |
|---|---|---|
| README is 402 lines | **wrong** — it is 401 | — |
| "seven near-identical prompt blocks" | **wrong** — there are eight (`README.md:69,74,79,84,91,99,107,114`) | keep 3 |
| Repository Layout omits `addons/` | **true, false by omission** | rewrite the tree |
| Windows example says `v1.2.1`, package is `1.5.0` | **genuinely false** | fix or make versionless |
| Lifecycle is "append-only" (`:21-23`) | **genuinely false** — only signoff blocks are; `IMPLEMENTATION.md` is living | rewrite |
| Consensus lens "rates confidence by agreement" (`:26-27`) | **genuinely false** — no confidence rating exists in the protocol | rewrite |
| "any capable tier-1 model can follow it" (`:119`) | **unsupported**, not false — no defined tier or test | delete |
| Default uses "all discovered installed CLI agents" (`:371`) | **genuinely false** — `SKILL.md:280-289` specifies a bounded set, normally 2–4 | rewrite |
| "the protocol's value should be obvious" (`:397`) | **an opinion standing where evidence should be** | delete |
| Opening runtime list (`:9`) | **imprecise, not false** — names five runtimes plus a custom directory, reads as exhaustive | complete it |
| Line 9 "contradicts" line 186 | **wrong accusation** — incompleteness, not contradiction | — |
| "WinGet manifest not yet accepted" (`:242`) | **genuinely false** — see C10 | fix |

**C10 — Two corrections that came from outside a participant's own file.**

1. **"15 agent runtimes" was invented by claude-1** in round 1, self-flagged as unverified,
   and independently refuted by codex-1 and hermes-1. The installer defines **fourteen named
   runtime targets** (`lib/installer.js:13-113`) plus `generic`, which is a destination
   requiring an explicit `--dest`, not a runtime. **The shipping copy says "fourteen named
   targets plus a generic directory". "15 runtimes" must not appear anywhere.**

2. **The WinGet accusation was correctly rejected on the evidence available in round 2, and
   is now sustained on new evidence.** codex-1 was right that no *shipped file* proves
   publication — `packaging/winget/README.md:3-4` still calls the manifest a draft. The
   facilitator then checked the external source of truth:
   `gh api repos/microsoft/winget-pkgs/contents/manifests/f/Feci/ParleyDeckSkill` lists
   published version directories from `1.0.4` through `1.4.6`. **`README.md:242` is therefore
   false and is fixed.** The evidence is external and is cited as such; the shipped
   `packaging/` guide should be corrected in a follow-up, since it is now stale too.

Two errors in the kickoff brief, both found by codex-1, are also recorded: the brief listed
`parley-design-check` exit codes `0–3` when the shipped skill defines **`4`** for an
`UNJUDGEABLE` run (`addons/parley-design-check/SKILL.md:53-61`, implemented at
`lib/engine.js:2069`), and the brief's "seven prompt blocks" count was wrong. The brief is not
authoritative where a shipped file disagrees.

**C11 — Claims that must be qualified in the copy.**
`parley-tracker` authors canonical files and neutral projections; **live tracker
create/update requires a separate opt-in connector** (`addons/parley-tracker/SKILL.md:381-388`)
and the copy must not imply the skill writes to Jira by itself. The gap-scan checks a
specified required schema and required sections (`:212-250`), not "every field". The shipped
superlative "the single highest-leverage rule for AI output quality" is attributable but
unmeasured and does **not** ship in the README.

**C12 — Line budget: ≤ 300 lines**, from 401. codex-1 proposed 280, hermes-1 ~284 against a
320 ceiling. 300 is adopted as the binding ceiling with the section order of C1. A budget
that forces a second pass is worse than one honoured on the first.

**C13 — Out of scope, routed elsewhere.** claude-1's round-1 finding that
`npx skills add feci/parley-deck-skill` discovers 1 of 5 skills is a **packaging** defect, not
a README defect. It belongs to the `skills-cli-install-path` idea and must not change any
sentence in this one. This README may not promise five skills through an install path that
delivers one — a constraint on wording that this idea inherits.

---

## Dissent recorded

- **kimi-1 (F1)**: would place an ~15-line install block above the catalogue for the
  colleague-sent reader. Not adopted; the reader-model argument is adopted, the remedy is not
  (C1).
- **hermes-1 (C7)**: preferred the kimi-1 hook. Overruled by user ruling, not by the drafter.

## Follow-ups (not in this change)

1. `parley-deck-cli/README.md` — same treatment, different audience, separate idea.
2. `packaging/winget/README.md` — stale; still describes the manifest as an unaccepted draft.
3. The `skills` CLI discovery defect — owned by `skills-cli-install-path`.

---

## Signoffs

### Signoff: codex-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS

C7 is a clean resolution, not drafter favouritism: claude-1 recused, my round-2 review independently chose that hook, and the user made the tie-break. I agree with every C9 factual verdict; the WinGet row is factually correct on the GitHub API evidence, while its provenance belongs in C10 rather than being represented as shipped-file proof. My reservation is that C10 is a transparent but real exception to the round's shipped-file evidence rule: it is acceptable for removing the stale negative statement because publication status necessarily lives in the external registry, but it must not be cited as precedent for replacing shipped-file verification elsewhere or for adding a positive publication promise until a shipped file records it. C12 is achievable: my round-2 allocation fit 280 lines including 32 lines of reserve, so five entries of the agreed length fit safely under 300. The packaging-guide follow-up already records the needed reconciliation.

### Signoff: hermes-1 — 2026-07-29
Status: ✅ ACCEPT

C7 is a clean resolution, not a recusal dressed around a drafter-favourable outcome: codex-1 endorsed claude-1's hook on the merits in his own round-2 review before any tie-break was needed, claude-1 did not vote, and the user — not the drafter — broke the 1-1 tie; an independent reviewer's preference plus user authority is not self-dealing just because the winner was the drafter's candidate. I re-verified every C9 verdict against the shipped files this round (401 lines, eight prompt blocks at README.md:69,74,79,84,91,99,107,114, `v1.2.1` at :239 against package.json `1.5.0`, and fourteen named targets plus `generic` in lib/installer.js) and I disagree with none. C10's WinGet reversal is an acceptable but genuine exception to the shipped-file evidence rule — publication status is an external fact that no file in this repo can prove, so external registry evidence is the only possible proof, and the consensus flags it as such; I share codex-1's guard that it must not become precedent for internal claims or for shipping a positive `winget install` promise until a shipped file records it. C12 is achievable: my own round-2 budget landed at ~284 and codex-1's at 280 with 32 lines of reserve, so five entries of the agreed length fit under 300 with margin, and the cuts identified (seven prompt blocks → two, collapsed install/update, deleted "What The Skill Does" and "Why This Exists", rewritten layout) are real and sum to roughly the ~170 lines the catalogue costs.

### Signoff: kimi-1 — 2026-07-29
Status: ✅ ACCEPT

C7 is a clean resolution, and the strongest evidence is what did not happen: my runner-up line is recorded in full and explicitly not grafted in, so the loser's opening is preserved on the record instead of being averaged into the shipping copy — a drafter-favourable stitch-up would have merged the two lines and called it consensus. claude-1 did not vote, codex-1 endorsed the claude-1 hook on the merits in his own round-2 review before any tie-break existed, and the user — the only authority above the participants — broke the 1-1; that is procedure, not favouritism. On C9 I re-checked the table against the round-2 evidence and disagree with no verdict, including the two rows that land on me: my round-1 `parley-deck` entry repeated the false "rates confidence by agreement" claim, and my F7 called line 9 a contradiction of line 186 when the correct verdict is incompleteness — both corrections are right and both stand. The "append-only" row is the one place the two reviewers split on the label (genuinely false versus imprecise) while agreeing on the remedy; the table's harsher label is defensible because the sentence as written claims for the whole lifecycle something that is only true of signoffs, so I do not contest it. C10's WinGet reversal is legitimate as a bounded exception: publication status is a fact that exists only in the external registry, no shipped file in this repo can prove it either way, and the reversal removes a false negative rather than minting a positive `winget install` promise — I join codex-1's guard that it sets no precedent for internal claims that shipped files can decide. C12 is achievable, but only honestly: codex-1 caught that my round-1 arithmetic (keep three prompt blocks, ~88 lines saved) lands at 388–396, not ≤300, so the budget holds under the round-2 cut list — codex-1's 280 with 32 lines of reserve and hermes-1's ~284 — not under my softer round-1 plan.
