# Proposed rule text for cross-review

This is a design draft by codex-1. The participants must ratify it before any protocol edit.

## H1: replace section 6 rule 6

6. **English only; ASD-STE100 by default.** Write new protocol prose in English. This includes every artifact, inbox note, brief, PR/MR text, and commit message. Use ASD-STE100 Issue 9 as summarized in Appendix B. Read that appendix before you write or review English prose. Keep verbatim source quotations in their original language; add an English explanation when needed. Do not rewrite exact syntax or evidence. Non-English chat with the owner is outside this file rule. STE applies only to English output. Appendix B defines exemptions, adoption, and the per-idea exception. A deck can retain a deliberate, documented working-language override.

## H2: append one on-demand protocol appendix

## Appendix B — ASD-STE100 output text

**Authority and access.** Use ASD-STE100 Issue 9, January 2025, for new English prose under this protocol. The official download page identifies this issue as current on 2026-10-13: https://www.asd-ste100.org/STE_downloads.html#article02-2l. Request the free official standard there. The official FAQ is https://www.asd-ste100.org/STE_faq.html. This summary uses original words. It does not contain the copyrighted dictionary or replace the full standard. A short summary or an automated check cannot certify full conformance.

**Core rules.** Limit an instruction sentence to 20 words and a descriptive sentence to 25 words. Use at most six sentences per paragraph. Give each paragraph one subject. Give one instruction per sentence. Use the imperative for procedures. Use active voice; description can use passive voice if the actor is unknown. Use simple present, simple past, simple future, or imperative forms. Do not use an `-ing` form as a verb. Permitted nouns and adjectives can end in `-ing`. Use approved words with their approved meanings and parts of speech. Use technical names and technical verbs only where the standard permits them. Define project terms and use each term consistently. Consult the full standard for vocabulary disputes; do not invent dictionary approval.

**Scope and exemptions.** Apply these rules to new English prose in rounds, consensus, FINAL, IMPLEMENTATION, reviews, signoffs, inbox notes, briefs, PR/MR mirrors, commits, and other protocol output. Preserve verbatim quotations from the owner, sources, and other agents. Preserve code, commands, paths, identifiers, URLs, log and error excerpts, schema keys, frontmatter, required headings, and verdict or signoff syntax. Preserve generated and validator output. Exempt spans retain their exact content; the surrounding explanation follows the rules. Agent-authored prose is not exempt merely because an agent generated it. Non-English owner chat remains outside the English file rule.

**Default and exception.** Absent `output_style` means `ste100`. The idea author can set `output_style: plain` in `00-prompt.md` and state an English reason there. `output_style: ste100` records the default explicitly. Unknown, empty, or malformed values do not disable the default. After kickoff, a change needs owner direction or participant consensus in the next canonical artifact. The exception affects that idea's prose only. It changes no language, evidence, ownership, or phase gate. There is no new deck configuration key. Agents apply this prose instruction; the CLI does not validate its syntax or use it to pass a phase.

**Adoption and review.** New decks use this rule. Existing decks adopt it through their normal protocol update. A newer CLI or installed skill does not override an older live authority. After adoption, apply the rule to every new artifact, including artifacts for open ideas. Apply it to newly authored prose in a living artifact. Do not rewrite existing text, signatures, closed artifacts, or old protocol prose. A reviewer cites the exact new text and the applicable rule for a violation. Any mechanical prose check is advisory and must skip exempt content. This release adds no prose linter or dictionary checker. A style heuristic cannot block a phase or replace independent review.

## H3: map entry

Use the exact Appendix B heading above with `include: on-demand`.
Trigger: before an agent writes or reviews new English prose under section 6.6.
Keep section 6 in the existing always-included set.

## H4: skill pointer

Before you write or review new English protocol prose, read `references/STE100.md`.
Use the live authority's writing policy and its documented per-idea exception.
The new reference has the same rule body as Appendix B.

## H5: launch instruction

Follow the live protocol's output-text policy. Where it adopts ASD-STE100, read Appendix B before you write English prose. Read this idea's `00-prompt.md` for a documented `output_style: plain` exception. Preserve exact quotations, code, syntax, and evidence. Non-English owner chat is outside the English file rule. An older live authority keeps its own policy.

Apply H5 to all runner and driver participant briefs. The existing shared launch wrapper is a possible single insertion point. Cross-review must prove that it covers every real brief path before selecting it.
