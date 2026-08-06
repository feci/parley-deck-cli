---
idea: skill-sync-cli-1-39
review-round: 1
drafter: claude-1
reviewers: [codex-1, hermes-1, kimi-1]
reviewed-commit: 661af98
date: 2026-08-06
status: fixes-agreed
---

# Review consensus — round 1

Three reviewers, one MAJOR, three MINOR, two NIT. **The MAJOR is a contradiction the implementer
introduced while fixing a contradiction**, which is the finding that matters most here.

## Reviewer verdicts

| reviewer | CRITICAL | MAJOR | MINOR | NIT |
|---|---|---|---|---|
| codex-1 | 0 | 1 | 1 | 0 |
| hermes-1 | 0 | 0 | 2 | 1 |
| kimi-1 | 0 | 0 | 1 | 1 |

All three independently confirmed that D1's `opencode` row and D2's replacement paragraph are
**verbatim** matches against the adopted text (modulo Markdown wrapping), and all three proved the
D4 guard fires by desyncing `compatibility.json` and restoring it.

## Verdict conflict — D3

**codex-1 filed a MAJOR against D3; kimi-1 explicitly cleared D3** (*"the D3 branch split matches
the CLI's actual launch path"*). hermes-1 called the split "structurally correct" and did not
examine the interaction.

**Resolved for codex-1, on reproduction rather than on the 1-to-1 count** (§15.3). The drafter
reproduced it directly:

- `SKILL.md:381` (Headless Agent Configuration, the **manual** camelCase JSON shape) states
  *"Nothing is appended to it afterwards."*
- `SKILL.md` branch A step 3 states *"Add model/thinking/profile flags only when discovered or
  configured"*, and step 4 delivers the prompt.

Those are contradictory instructions in one document about one activity. kimi-1's clearance
addressed whether the two branches individually describe the CLI correctly — they do — and did not
test the absolute paragraph against branch A. The finding stands.

## Agreed fixes

### AF-1 — MAJOR (codex-1). Scope the absolute paragraph to the write-enabling flag

`skills/parley-deck/SKILL.md:381`. The paragraph applies branch B's CLI rule ("nothing is
appended") to the manual JSON shape, which branch A contradicts two sections later. **Fix:** keep
the part that is true of both — the write-enabling flag belongs inside `headlessArgs`, there is no
separate write-mode list — and remove the unscoped "nothing is appended to it afterwards" from the
manual section, where "complete argv template / nothing appended" belongs to branch B's snake-case
`headless_args` alone. Label the shape as manual-facilitator input.

### AF-2 — MINOR (codex-1 and kimi-1, independently). Remove the `hermes` allusion from CHANGELOG

`CHANGELOG.md`. The release note says the removed list is the model that made *"the `hermes`
regression invisible"*. D6 excluded the incident narrative. **Fix:** state the mechanism without
the incident.

### AF-3 — MINOR (hermes-1). Fold the extracted vendor-flag sentence back into the D2 paragraph

`skills/parley-deck/SKILL.md:254`. Pre-commit, *"a vendor flag change is a config edit, not a skill
revision"* was a clause of the sentence D2 replaced. The implementation extracted it as a standalone
line — same words, different structure, and FINAL decided no structural change. **Fix:** hermes-1's
option (a) — append it as the final sentence of the D2 paragraph, making the replacement drop-in.

### AF-4 — MINOR (hermes-1). Correct the `IMPLEMENTATION.md` narrative

`npm test` is `node --test && … && build-addon-manifest.js --check`, so a single `npm test` run
exercises the guard through **both** callers, not one. **Fix:** documentation accuracy only.

### AF-5 — NIT (kimi-1). Note the `headless_args` → `headlessArgs` adaptation

FINAL's migration rule is written in snake_case; the manual JSON shape is camelCase, and the
implementation adapted it silently. **Fix:** record the adaptation in `IMPLEMENTATION.md`. Correct
as written — the manual shape genuinely uses camelCase — but an undeclared adaptation of ratified
text should be visible.

## Dismissed

**NIT (hermes-1) — the 1124-character D2 line.** Not fixed, and hermes-1 itself says not to:
rewrapping would break the verbatim match all three reviewers just confirmed. The verbatim
property is worth more than the diff ergonomics. Revisit only if the text is ever re-ratified.

## Facilitator process failures recorded

Two, both the facilitator's, both affecting this round rather than the artifact:

1. **codex-1 could not write its own review file.** It was launched with the sandbox root set to
   `parley-deck-skill`, so the deck path in `parley-deck-cli` was outside its writable area. It
   wrote the review to a temp directory and reported the block. The facilitator copied the file
   byte-for-byte without editing; codex-1 authored every word. Recorded because a facilitator
   moving a participant's artifact is exactly the thing the protocol's no-proxy-write rule guards,
   and it happened here because of a facilitator error in the invocation, not by design.
2. **The MAJOR was self-inflicted.** The idea existed because the skill taught a two-list launch
   model; the fix introduced an absolute one-list statement in the wrong section. Reviewers caught
   it; the implementer did not.
