---
idea: meta-protocol-change-ste100-output-format
drafted-by: codex-1
date: 2026-10-13
---

## Agreed decisions

Adopt ASD-STE100 Issue 9 as the default for new English protocol prose.
Use the exact rule text in `source-context/protocol-hunks.json`.
H1 replaces section 6 rule 6. H2 appends one Appendix B.
Apply both hunks to the live deck, embedded CLI default, skill protocol copy, and staged core.
Keep the other bytes of each protocol copy unchanged.

Appendix B is the required working summary. Its text binds protocol decisions.
The full standard does not override protocol authority or the explicit exemptions.
The summary is not a full conformance claim.
Zcode confirmed the current issue through PRIMARY official-source evidence in `round-01/zcode-1.md`.
The numerical rules are local adopted limits. Their attribution to the full standard remains UNVERIFIED.
No acceptance criterion depends on that attribution.

The summary covers sentence and paragraph limits, direct instructions, active voice, simple verbs, vocabulary, and stable technical terms.
It preserves exact source quotes, all frontmatter, required syntax, code, paths, logs, generated text, and validator output.
It states that non-English owner chat is outside the English file rule.
New agent prose is not exempt merely because an agent produced it.

Use `ste_output: false` for the per-idea exception. The author records an English reason at kickoff.
After kickoff, a change needs owner direction or participant consensus in the next canonical artifact.
Other values retain STE. The CLI does not parse this advisory metadata or use it as a gate.
Add no deck configuration key. The exception changes no other protocol duty.

New decks use the rule. Older decks adopt it through their normal authority update.
A rule-modifying consumer sync needs user confirmation under section 9.0.
After adoption, use the rule for every new artifact, including new artifacts for open ideas.
Apply it to newly authored prose in living artifacts. Preserve existing text and closed history.

Ship no prose linter, executable check recipe, dictionary checker, or avoid-list.
Use the manual checklist in the common rule body.
Pure prose-style feedback is advisory and stays outside the fixed-severity findings list.
It cannot block a phase or a strict-gate close.
Technical meaning, evidence, and correctness defects remain findings under existing rules. Do not conceal them as style feedback.

Add the common instruction to all twelve participant brief builders.
Also add it to `trajectoryVerifierPrompt`, which produces an agent-authored execution report.
This makes thirteen model-facing paths. Mechanically generated kickoff templates remain exempt.
Tests must verify coverage and preservation of exact input and syntax.

Use one on-demand map entry for Appendix B. Keep the binding section 6 pointer always included.
Keep SKILL.md at or below 20,000 bytes and preserve the existing 70,000-byte facilitator test.
Report the larger meta-idea packet too; it already exceeds that test's limit in the baseline.
Do not change limits or rewrite unrelated protocol text.

Target CLI 1.55.0, skill 2.19.0, and core 2.19.0.
Stage the core from staged 2.18.0 plus only H1 and H2.
Follow the controlling brief for all release channels and installed runtimes.
Keep npm, core publication, and the D6 legacy declaration as the required owner actions.

## Agreed trade-offs

The summary can work offline. Full vocabulary verification still requires the official standard.
On-demand detail saves core bytes. Agents must follow the pointer before they write or review English prose.
The style instruction relies on agent and reviewer discipline. It creates no machine-enforced conformance claim.

The shared constant avoids thirteen copies of the instruction.
The tests protect real prompt paths. They do not prove linguistic conformance.
No full historical rewrite is part of this change.

## Open items deferred to implementation

The implementer must keep the exact hunk payload and record final byte counts and hashes.
The implementer must include `trajectoryVerifierPrompt` and its test in the file list.
The skill reference must have the same rule body as Appendix B.
No source decision remains open.

The source goal check occurs before merge and publication.
The separate channel check occurs after publication and runtime installation.
Task completion requires both checks. Source acceptance alone does not establish delivery.

No linter, full dictionary check, automatic consumer migration, or historical re-authoring is scheduled.
Any future proposal needs its own scope and authority.

## Comparison & blind spots

Codex proposed Appendix B and an explicit adoption boundary.
Zcode independently confirmed the current issue and identified the strict-gate interaction.
Codex found two missing private driver builders. Zcode corrected its inventory and found the additional trajectory report.
Both participants now support the same small implementation family.

Unanimity among related models is a shared prior, not independent evidence.
The design would be wrong if its exemptions altered evidence or its prompt instruction overrode an older live authority.
It would also be wrong if short prose lost technical meaning or one launch path omitted the instruction.
Tests and independent review target those failures.
Neither participant inspected the full copyrighted dictionary. Full ASD conformance remains unverified.

Codex is organizer, participant, consensus drafter, proposed FINAL drafter, and designated implementer.
Zcode provides independent code review and a fresh goal check. The organizer has no extra signoff weight.

## Drafter position changes

1. Prior: "Reviewers retain the existing severity rules" (`round-02/codex-1.md`, Responses to others).
   New: pure prose-style feedback stays advisory outside formal findings. Technical meaning and evidence defects retain existing severities and gates.
   Reason: Zcode identified an ambiguity about strict-gate NIT findings. The exact Appendix B text resolves it without concealing correctness defects.
2. Prior: "Use one shared instruction constant across all twelve builders" (`round-02/codex-1.md`, Responses to others).
   New: also include `trajectoryVerifierPrompt`, for thirteen model-facing paths.
   Reason: Zcode identified its agent-authored report. A uniform instruction costs one call site and test.
3. Prior: "The short summary is a protocol rule summary" (`round-01/codex-1.md`, Proposed approach).
   New: the summary explicitly binds protocol decisions when external standard text differs.
   Reason: Zcode requested an explicit authority boundary. No external document can change protocol gates or evidence duties.

Other material positions remain as stated in round 2.
The opt-out name was already accepted there and is not a new drafter change.

## Alternatives disposition

- ALT-1: Reject English-only as the complete solution. It does not define controlled prose.
- ALT-2: Reject a full rule block in the core skill or always-included packet. It would consume the size margin.
- ALT-3: Adopt Appendix B, the matching skill reference, and one on-demand map entry. Each authority remains self-contained.
- ALT-4: Adopt `ste_output: false` with a reason and controlled mid-idea changes. Reject `output_style` and a new deck option.
- ALT-5: Reject a linter, executable recipes, and an avoid-list for this release. Use a manual checklist.
- ALT-6: Reject a bundled dictionary or a full conformance certificate. Cite the official standard and preserve verification limits.
- ALT-7: Reject historical re-authoring. Existing evidence and signed artifacts stay exact.
- ALT-8: Adopt one constant in thirteen prompt builders. Reject dependence on the launch-context wrapper alone.
- ALT-9: Reject automatic older-deck adoption. Follow the live authority and section 9.0.
- ALT-10: Reject section 16, extra Quickstart text, and a TL;DR entry. The short section 6 pointer is sufficient.
- ALT-11: Reject silent exceptions to numerical limits. Use shorter sentences or the documented per-idea exception.
- ALT-12: Reject a new style close gate. Keep style-only feedback advisory and keep correctness findings under the existing gates.

These dispositions cover the mechanisms in both round-1 alternative lists.
No contradictory factual verdict remains. Zcode withdrew its earlier inventory and sync claims.

## Signoffs

<!-- codex-1 appends its signoff below -->
<!-- zcode-1 appends its signoff below -->

### Signoff: codex-1 — 2026-10-13
Status: ✅ ACCEPT
Notes: I accept the exact hunks and decisions. The role concentration and position changes are complete. Drafter: yes.

### Signoff: zcode-1 — 2026-10-13
Status: ✅ ACCEPT
Notes: I accept. I checked the exact hunk text in source-context/protocol-hunks.json against my round-2 requirements. The ste_output: false rewrite is in. The binding-summary line and the strict-gate advisory line are in. trajectoryVerifierPrompt is the thirteenth path. The consensus keeps the UNVERIFIED attribution of the numerical rules, and no claim is DISPUTED. The drafter position changes match the round record. Ready for FINAL.
