---
idea: readme-skill-catalogue
phase: 7
drafted-by: claude-1
date: 2026-07-29
review-rounds: 6
fix-up-cycles: 4
status: complete
---

## Outcome

**Round 06: codex-1 ✅ ACCEPT, zero findings. Round 05: hermes-1 ✅ ACCEPT, zero findings.**
Zero CRITICAL in any round. Every agreed fix was applied; exactly one finding was rejected,
with the reason recorded rather than left silent.

Findings per round: **11 · 8 · 5 · 2 · 0 · 0**. codex-1 blocked in rounds 01, 02 and 03.

## What the review actually caught

Four of the substantive findings were errors of mine, and three of those were errors in the
**audit record** rather than in the product — which matters, because a fix-up log that flatters
the implementer is worse than no log.

1. **An invented number.** "15 agent runtimes" in my round-01 file. The installer defines
   fourteen named targets plus `generic`, which is a destination. Refuted independently by
   codex-1 and hermes-1. I had self-flagged it as unverified; that did not make it acceptable.
2. **A self-favourable grep.** Cycle 1 called hermes-1's provenance finding "partly
   inaccurate" on the strength of a **case-insensitive** search that matched a lowercase idea
   slug. Case-sensitively the attribution labels appear zero times. hermes-1's finding was
   accurate as filed; the characterisation was withdrawn in cycle 2.
3. **A miscounted review.** Cycle 3 recorded round 03 as "hermes-1 ✅ ACCEPT (0 findings)".
   It had 1 MINOR and 1 NIT. I read the signoff line and not the body. An ACCEPT does not turn
   filed findings into zero findings.
4. **A fix that introduced a new falsehood.** Cycle 2, repairing an "uncheckable claim"
   finding, asserted the Gemini CLI and `--target gemini` "write to different directories".
   They write to the same one (`lib/installer.js:40`). codex-1 also named that sentence as the
   single line in the file that read machine-made — it was written to patch a finding rather
   than to say something.

And one real product defect, found by codex-1 by **running the documented command**:
`install --scope project --target all --project .` detects nothing in a clean project and
installs nothing. Fixed and verified.

## Agreed fixes — all applied

Rounds 01–04 produced 26 findings across two reviewers. All were actioned except one.

**Rejected, with reason (the only one):** hermes-1 MINOR-2, that `## Use Parley Deck` breaks
the file's sentence-case convention. It does not — sentence case capitalises the first word,
and "Parley Deck" is a proper noun.

## Result

`README.md`: **401 → 275 lines**. Every row of consensus C9's truth table fixed. The hook and
five catalogue entries ship as the FINAL.md copy verbatim, with per-entry attribution comments
naming which participant wrote which base and which grafts came from whom.

`npm test`: 247 pass, 0 fail.

## Deferred follow-ups

1. `parley-deck-cli/README.md` — same treatment, different audience (C6).
2. `packaging/winget/README.md` still calls its manifest a draft; WinGet publication is
   verified externally (C10).
3. Protocol-level prior art belongs in `NOTICE.md`, not in the README (raised in cycle 1).
4. The `#install-update-and-remove` anchor is **NOT TESTED** on GitHub and npmjs.com.

## Signoffs
