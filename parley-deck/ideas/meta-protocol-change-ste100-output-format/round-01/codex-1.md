---
agent: codex-1
idea: meta-protocol-change-ste100-output-format
round: 1
date: 2026-10-13
---

## Summary

Adopt ASD-STE100 Issue 9 as the default for new English prose under this protocol.
Put the operational summary in a protocol appendix and a matching skill reference.
Keep exact syntax and source evidence unchanged.
Use agent review for this release; defer a mechanical prose checker.

## Proposed approach

Use one short binding addition to section 6.6.
Replace its instruction to translate source quotes with permission for exact source quotes plus an English explanation.
State that non-English owner chat is outside the English file rule.

Add an on-demand Appendix B to each COOPERATION.md copy.
Add `skills/parley-deck/references/STE100.md` with the same rule body.
Add one short pointer to SKILL.md.
Classify Appendix B as on-demand, with a trigger before an agent writes or reviews English prose.
The pointer in section 6 remains in every packet.
This preserves access for global cores that distribute only COOPERATION.md.

Use these local core rules:

1. Limit an instruction sentence to 20 words. Limit a descriptive sentence to 25 words.
2. Limit each paragraph to six sentences. Give each paragraph one subject.
3. Give one instruction in each sentence. Use the imperative for a procedure.
4. Use active voice. A description can use passive voice when the actor is unknown.
5. Use simple present, simple past, simple future, or imperative forms.
6. Do not use an `-ing` form as a verb. Permitted nouns and adjectives can end with `-ing`.
7. Use approved words with their approved meanings and parts of speech.
8. Use documented technical names and technical verbs when the standard permits them. Use one term for one concept.

The short summary is a protocol rule summary. It is not a replacement for the standard or a full conformance certificate.
Agents must consult the official standard for disputed vocabulary or a claim of full conformance.
Do not invent a dictionary from memory.

Preserve exact quotations, code, commands, paths, identifiers, logs, errors, schema keys, frontmatter, generated output, and validator text.
Preserve required protocol headings, verdict tokens, signatures, and URLs as syntax.
The prose around an exempt span still follows the rule.
Do not classify authored explanations as generated text merely because an agent wrote them.

I propose a per-idea `output_style: plain` opt-out, with an English reason in the kickoff.
Absent `output_style` means STE. An explicit `output_style: ste100` states the default.
Limit the exception to new prose for that idea.
Do not add a deck option in this release; a documented protocol overlay can already express a deliberate deck policy.
Malformed or unknown values must not silently disable the default.
This is a writing instruction, not a new validator or close gate.

After a deck adopts 2.19.0, apply the rule to new artifacts and newly authored prose in open artifacts.
Preserve old text, signatures, and closed records.
An installed newer CLI or skill must not silently replace an older live authority.
Prompts must say to follow the live rule and the documented per-idea exception.

Add the same short instruction to round, cross-review, consensus, FINAL, implementation, review, fix-up, signoff, and goal briefs.
Include consult and steer prompts because they produce protocol communications.
The instruction points to the live appendix. Manual skill launches use the skill reference after resolving live authority.
A common helper can prevent drift across prompt builders without new parsing or persistence.

## Existing alternatives

- ALT-1: Use only the existing English rule. Section 6.6 sets a language but no controlled prose standard. Reject this as insufficient.
- ALT-2: Put all rules in SKILL.md and the packet. Existing size tests guard 20,000 and 70,000 bytes. Reject unnecessary core growth.
- ALT-3: Use a protocol appendix plus skill reference. The packet map already supports `include: on-demand`. Adopt this delivery model.
- ALT-4: Add an advisory sentence linter. Markdown syntax, quotations, headings, abbreviations, and technical names need careful exemptions. Defer it for this small release.
- ALT-5: Add a dictionary checker. The brief forbids distribution of the dictionary. Reject this scope.
- ALT-6: Rewrite all prior artifacts. Existing append-only ownership and frozen FINAL rules make this incorrect. Reject it.
- ALT-7: Use existing live-authority launch handling. `internal/runner/protocol_context.go` supplies full resolved context today. Adopt its authority boundary.

Sources inspected: section 6.6 and Appendix A; `parley-deck/meta/packet-applicability.yaml`; `internal/runner/runner.go`; `internal/runner/protocol_context.go`; prompt-builder locations in `internal/runner/phase58.go`, `internal/runner/consult.go`, `internal/runner/steer.go`, `internal/app/driver_consensus.go`, and `internal/app/consensus_request_signoffs.go`.
Size guards are `internal/app/facilitator_packet_live_test.go` and `../ste100-skill/test/lean-organizer.test.js`.

## Source evidence

Source-dependent testimony: the official download page calls Issue 9, January 2025, the current issue.
The exact URL and short quotation are in `source-context/ste-source-notes.md`.
I consulted that page and the official FAQ with ego-browser.
I propose this version selection; a non-owner must issue its section 15 verdict.
Exact numeric rule identifiers and full dictionary conformance remain unverified.
The proposed numeric limits are local protocol choices from the owner's brief unless independently verified against the standard.

## Concerns / open questions

The source authority must carry the summary for CLI-only and global-core consumers.
A skill-only reference would leave those consumers without the required rules.
An on-demand appendix needs an explicit read trigger, even though most writing phases activate it.

The opt-out must not change English-only, evidence, ownership, or signoff duties.
Reviewers should cite exact new prose and the applicable rule when they report a violation.
No style heuristic should automatically block a phase.

## Risks

Agents can overstate conformance without access to the full dictionary.
Technical words need stable definitions and cannot provide an unlimited vocabulary exception.
Automatic counters can misread exempt evidence and encourage harmful rewrites.
An overly long pointer can exceed the packet margin.

## Context attestation

`context_mode: full`; source_sha256 and packet_sha256: `0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`.
No fallback reason. Full live protocol read before this independent round.
No other participant's round-1 artifact was read.
