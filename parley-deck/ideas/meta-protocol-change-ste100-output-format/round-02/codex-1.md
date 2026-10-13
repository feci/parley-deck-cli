---
agent: codex-1
idea: meta-protocol-change-ste100-output-format
round: 2
date: 2026-10-13
responding-to: [zcode-1/round-01]
---

## Position changes since round 1

I accept `ste_output: false` as the per-idea exception. It needs less metadata than my proposed `output_style` values.
I reject executable check recipes and an avoid-list for this release. A manual review checklist is sufficient.
I retain ASD-STE100 as the named default. Our summary does not certify full conformance.

## Responses to others

### @zcode-1 — round-01

We agree on the core rule set, on-demand detail, unchanged history, and no machine gate.
Your source verification supports Issue 9, January 2025. FINAL can cite your PRIMARY verdict on my source notes.
The exact numeric limits remain local adopted rules. Their attribution to the full standard remains unverified.

The owner asked for ASD-STE100 as the default standard. We must preserve that instruction.
Your phrase "not to ASD-STE100 itself" weakens the requested outcome.
Counter-proposal: name ASD-STE100 as the default and name this document its required working summary.
State that the summary neither replaces the standard nor proves full conformance.
Require the full standard for disputed vocabulary and any full conformance claim.

Use Appendix B, with one heading and no subsections. This matches the owner's suggested appendix and needs one map entry.
Use the same rule body in the skill reference. A larger, different skill rule set could create authority drift.
The appendix keeps global-core consumers self-contained.

Replace section 6 rule 6 with the short proposed hunk. An appended exemption could conflict with its existing translation instruction.
Preserve original source quotes. Add English explanation as needed.
No Quickstart or TL;DR change is necessary; section 6 already binds every packet.

The idea author can set `ste_output: false` at kickoff and record an English reason.
After kickoff, require owner direction or participant consensus. A unilateral style exception could evade a review finding.
Absent, true, empty, or unknown values retain STE. Only the explicit false setting disables it.
Agents interpret the field. The CLI does not parse it or use it as a phase gate.
This exception changes no other protocol duty. Add no deck configuration key.

Apply the word and paragraph limits as rules. Do not add an undocumented clarity exception to a numeric limit.
Preserve meaning through shorter sentences or the documented idea exception.
Reviewers retain the existing severity rules. The new standard creates no automatic close gate.
It also grants no exemption from an existing review gate.

Do not ship executable check recipes. A recipe still needs correct Markdown exemptions and maintenance.
Use a short manual checklist in the common rule body.
Do not ship an avoid-list that readers could mistake for an approved dictionary.

Use the normal authority update for existing decks. Do not promise an automatic compatible sync.
Section 9.0 requires confirmation for a rule-modifying sync. This rule changes prose duties.
The owner's release authority covers this source change and the requested skill installs.
It does not silently replace each project's live authority.

Your prompt inventory misses two private driver builders.
PRIMARY source evidence: `internal/app/driver_consensus.go` defines `buildConsensusDraftPrompt` and `buildFinalDraftPrompt`.
Both currently end their instructions with "English only".
Include these builders with your ten listed builders.
Use one shared instruction constant across all twelve builders. Test their output and preserve embedded input bytes.
This directly meets the brief's template requirement and avoids dependence on a single launch wrapper.

## New concerns / questions

The baseline size record is in `source-context/baseline-sizes.json`.
The no-idea facilitator command produces 68,904 bytes. The skill core has 19,235 bytes.
The meta-idea packet already has 71,888 bytes because it includes section 7.
The existing 70,000-byte test measures its named facilitator fixture, without the protocol-change flag.
Preserve that exact test and its limit. Report both measurements without hiding the pre-existing larger variant.
Do not rewrite unrelated protocol sections to reduce the older variant.

## Current proposal

Use `source-context/codex-1-proposed-hunks.md`, with the changes stated in this round.
The source edit list is:

- Three COOPERATION.md copies: replace rule 6; append Appendix B.
- Packet applicability map: one on-demand Appendix B entry with a read trigger.
- `internal/protocol/outputtext.go`: one shared prompt instruction.
- `internal/runner/runner.go`, `phase58.go`, `consult.go`, and `steer.go`: use that instruction.
- `internal/app/driver_consensus.go` and `consensus_request_signoffs.go`: use the same instruction.
- Prompt coverage tests in `internal/runner` and `internal/app`.
- A skill test for appendix/reference identity and required scope details.
- Skill `SKILL.md`, `references/STE100.md`, and `references/ARTIFACT_TEMPLATES.md`: point to the live policy and shared summary.
- CLI `VERSION`, `internal/app/version.go`, `CHANGELOG.md`, and `README.md`: record 1.55.0 and the new policy.
- Skill `package.json`, `package-lock.json`, `CHANGELOG.md`, `README.md`, and `references/compatibility.json`: record 2.19.0.
- Source `parley-deck/meta/version.json` and `meta/protocol-changelog.md`: record the reviewed protocol change.

Stage core 2.19.0 by applying only the two exact COOPERATION hunks to staged core 2.18.0.
Keep its other bytes unchanged. Record the base, patch, and result hashes.
Release and channel evidence follow the controlling brief.

## User direction and driver fallback

The brief states: "Use the recorded fallback. Never apply the legacy declaration yourself. Never prune or declare worktrees."
The real driver refused round 2: `cross-review accounting: historical cycle event lacks idea identity`.
It started no round-2 child. The original generated escalation will remain unchanged in source-context.
Use the configured CLI fallback. Do not claim a product driver pass or an automatic exclusion.

## Context attestation

`context_mode: full`; source_sha256 and packet_sha256: `0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`.
No fallback reason. The live authority remains unchanged during design.
