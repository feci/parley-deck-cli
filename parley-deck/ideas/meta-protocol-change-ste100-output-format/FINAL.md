---
idea: meta-protocol-change-ste100-output-format
status: final
author: codex-1
drafted-by: codex-1
implementer: codex-1
consensus-date: 2026-10-13
participants: [codex-1, zcode-1]
---

## Final plan / specification

Adopt ASD-STE100 Issue 9 as the default for new English protocol prose.
Use Appendix B below as the required working summary.
Its rules bind protocol decisions. The external standard cannot change protocol authority or exemptions.
The summary does not certify full conformance.

Apply exactly two protocol hunks. H1 replaces section 6 rule 6.
H2 appends Appendix B, with one blank line before its heading.
Preserve all other bytes in each protocol copy.
The exact payload is `source-context/protocol-hunks.json`.
Its SHA256 is `d304f0a6b944c272ab9b52debf3f52cc570f5977aa930f0ddb215f0a2cbf494b`.

H1 new text:

```markdown
6. **English only; ASD-STE100 by default.** Write new protocol prose in English. This includes artifacts, inbox notes, briefs, PR/MR text, and commit messages. Use ASD-STE100 Issue 9 through the working summary in Appendix B.

   Read Appendix B before you write or review English prose. Preserve exact source quotations and add English explanation when needed. Appendix B defines the exemptions and the per-idea exception.

   Non-English owner chat is outside this file rule. STE applies only to English output. A deck can retain a deliberate, documented working-language override.
```

H2 text:

```markdown
## Appendix B — ASD-STE100 output text

**Authority and access.** ASD-STE100 Issue 9, January 2025, is the default standard for new English protocol prose. This appendix is its required working summary. For protocol decisions, this summary is binding. The external standard does not override protocol authority or the exemptions below.

The official download page identifies Issue 9 as current on 2026-10-13: https://www.asd-ste100.org/STE_downloads.html#article02-2l. Request the free official standard there. The official FAQ is https://www.asd-ste100.org/STE_faq.html.

This summary uses original words. It does not contain the copyrighted dictionary or replace the full standard. Do not copy or redistribute the dictionary or long passages from the standard. A summary or an automated check cannot certify full conformance. Consult the full standard for vocabulary disputes and any full conformance claim.

**Core rules.** Parley adopts these limits as its working rules:

- Limit each procedural sentence to 20 words. Limit each descriptive sentence to 25 words.
- Use at most six sentences per paragraph. Give each paragraph one subject.
- Give one instruction per sentence. Use the imperative for procedures.
- Use active voice. Descriptive text can use passive voice when the actor is unknown.
- Use simple present, simple past, simple future, or imperative forms.
- Do not use an `-ing` form as a verb. Permitted nouns and adjectives can end with `-ing`.
- Use approved words with their approved meanings and parts of speech. Do not invent dictionary approval or rejection.
- Use technical names and technical verbs only where the standard permits them. Define project terms and use each term consistently.

**Scope and exemptions.** Apply these rules to new English prose in rounds, consensus, FINAL, IMPLEMENTATION, reviews, signoffs, inbox notes, briefs, mirrors, and commits. They also apply to other new English output under this protocol.

Preserve verbatim quotations from the owner, sources, and other agents, including quotations in another language. Preserve code, commands, paths, identifiers, URLs, log and error excerpts, schema keys, and all frontmatter. Preserve required headings, verdict tokens, and signoff syntax. Preserve mechanically generated text and validator output. Keep each exempt span exact; apply the rules to its surrounding explanation.

Agent-authored prose is not exempt merely because an agent generated it. Non-English chat with the owner remains outside the English file rule. The prose policy changes no source, evidence, ownership, or language requirement beyond section 6.6.

**Default and exception.** The default is on. The idea author can set `ste_output: false` in `00-prompt.md` at kickoff and state an English reason there. Only an explicit false setting disables the default. Absent, true, empty, or unknown values retain STE.

After kickoff, a change needs owner direction or participant consensus in the next canonical artifact. The exception affects only that idea's prose. It changes no other protocol duty. There is no new deck configuration key. Agents apply this writing instruction; the CLI does not parse the field or use it as a phase gate.

**Adoption.** New decks use this rule. Existing decks adopt it through their normal protocol update, subject to section 9.0. A rule-modifying consumer sync needs user confirmation. A newer CLI or installed skill does not override an older live authority.

After adoption, apply the rule to every new artifact, including artifacts for open ideas. Apply it to newly authored prose in a living artifact. Do not rewrite existing text, signatures, closed artifacts, or old protocol prose.

**Review.** Cite the exact new text and the applicable rule when you give feedback. Pure prose-style feedback is advisory, outside the fixed-severity findings list. It cannot block a phase or a strict-gate close. Defects in technical meaning, evidence, or correctness still follow the existing finding and gate rules. Do not hide those defects as style feedback.

This release adds no prose linter, executable check recipe, or dictionary checker. Any mechanical prose check must remain advisory and skip exempt content. It cannot replace independent review.

Before you publish new prose, do this manual review:

- Count the words in each non-exempt sentence and the sentences in each paragraph.
- Remove ambiguity without a change to technical meaning.
- Confirm that quotations and required syntax remain exact.
- Use the full standard for a vocabulary dispute; state any unverified limit on your claim.
```

The skill reference has the same body as Appendix B, after its title.
Keep a short binding pointer in SKILL.md.
Resolve the artifact template's instruction to translate source quotes. Preserve quotes and add English explanation when needed.
Add one on-demand Appendix B entry to the packet applicability map.
The existing section 6 pointer stays always included.

Use one shared instruction constant in `internal/protocol/outputtext.go`:

```text
Output text: Follow the live protocol's writing policy. Where it adopts ASD-STE100, read Appendix B before you write English prose. Read this idea's 00-prompt.md for a documented ste_output: false exception. Preserve exact quotations, code, syntax, and evidence. Non-English owner chat is outside the English file rule. An older live authority keeps its own policy.
```

Append it to these thirteen model-facing builders:

| File | Builders |
| --- | --- |
| `internal/runner/runner.go` | `BuildRoundOnePrompt`, `BuildRoundPrompt` |
| `internal/runner/phase58.go` | `BuildFixupPrompt`, `BuildImplementationPrompt`, `BuildReviewPrompt`, `BuildReviewConsensusPrompt` |
| `internal/runner/consult.go` | `BuildConsultPrompt`, `BuildGoalCheckPrompt` |
| `internal/runner/steer.go` | `BuildSteerPrompt` |
| `internal/app/driver_consensus.go` | `buildConsensusDraftPrompt`, `buildFinalDraftPrompt` |
| `internal/app/consensus_request_signoffs.go` | `buildConsensusSignoffPrompt` |
| `internal/app/trajectory_verify.go` | `trajectoryVerifierPrompt` |

Preserve every input byte and required output syntax.
Do not parse `ste_output` or add a CLI gate.
Mechanically generated kickoff text is exempt.

Change these remaining product files:

- CLI protocol: `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`.
- CLI map: `parley-deck/meta/packet-applicability.yaml`.
- CLI tests: focused runner and app tests for all builders, preserved input, and required syntax.
- CLI release data: `VERSION`, `internal/app/version.go`, `CHANGELOG.md`, and `README.md`.
- Source metadata: `parley-deck/meta/version.json` and `parley-deck/meta/protocol-changelog.md`.
- Skill: `skills/parley-deck/SKILL.md` and new `skills/parley-deck/references/STE100.md`.
- Skill references: `COOPERATION.md`, `ARTIFACT_TEMPLATES.md`, and `compatibility.json` under `skills/parley-deck/references/`.
- Skill tests: new checks for rule identity, essential semantics, and packaged reference availability.
- Skill release data: `package.json`, `package-lock.json`, `CHANGELOG.md`, and `README.md`.

Preserve the 20,000-byte SKILL.md limit and the existing 70,000-byte facilitator packet test.
Measure the actual meta-idea packet separately. Its baseline is already over 70,000 bytes.
Do not change the limit or rewrite unrelated prose to conceal that baseline.

Target CLI 1.55.0, skill 2.19.0, and staged core 2.19.0.
Build the staged core from `~/.parley/staging/COOPERATION-2.18.0.md` plus only H1 and H2.
Its base SHA256 is `9ead4475df6d6fa500317e7d347a0b6422bf3c202aca49f1d4b57a5fe14dc512`.
Record the base, patch, and result hashes.

## Purpose / user-visible outcome

Agents receive one clear writing policy for new English text.
Exact quotations and evidence remain intact.
Older live authorities retain their existing policy until an authorized update.
The per-idea exception is explicit and documented.
Style feedback stays advisory; correctness and evidence defects retain their existing gates.

## Context & orientation

The CLI worktree is `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100`.
The skill worktree is its sibling `ste100-skill`. Both use the owner-provided branch `ste100`.
The source deck resides in the CLI worktree's `parley-deck/`.
Design transport is GitHub PR 83: https://github.com/feci/parley-deck-cli/pull/83.

Codex is organizer, participant, consensus drafter, FINAL drafter, and implementer.
The owner selected this concentration to save Claude tokens. The organizer has no extra vote.
Zcode supplies independent code review and a fresh goal check.
Kimi had two measured readiness failures. The owner pre-authorized its exclusion.
The mechanical single-reviewer exception returns false.
The controlling brief permits close only with all four recorded conditions:

- Zcode's final review has no open CRITICAL or MAJOR.
- Both current participants author review-consensus ACCEPTs.
- A fresh Zcode goal process returns PASS.
- Current-tree independent evidence covers every source acceptance criterion.

The source goal check occurs before merge and release.
A separate fresh Zcode process verifies delivery after publication and installation.
Task completion requires both checks.

## Observable acceptance criteria

| ID | Source criterion and evidence |
| --- | --- |
| AC1 | Each of the three protocol copies equals its base plus exactly H1 and H2. The normal drift guard passes. |
| AC2 | Appendix B and the skill reference have identical rule bodies. Their text covers limits, exemptions, authority, adoption, exception, and advisory review. |
| AC3 | All thirteen real prompt builders contain the shared instruction once. Tests preserve injected quotations, frontmatter, code, paths, and required report/signoff syntax. |
| AC4 | The map resolves Appendix B as on-demand. Section 6 remains always included. The original 70,000-byte fixture and 20,000-byte skill guard pass. Actual packet sizes are recorded. |
| AC5 | The skill archive contains STE100.md. The template preserves verbatim non-English quotes. All version and protocol hashes identify CLI 1.55.0 and skill 2.19.0. |
| AC6 | Core 2.19.0 equals the verified staged 2.18.0 base plus only H1 and H2. The base, payload, and result hashes are recorded. |
| AC7 | Required CLI tests and build pass. Skill tests, package dry-run, portable build, installer dry-run, and doctor evidence are recorded. Any unrelated CI failure is disclosed. |
| AC8 | No dictionary, linter, executable prose-check recipe, new metadata parser, or historical rewrite enters the product diff. Independent review tests those boundaries. |

Run `go test ./... -count=1 -timeout 45m` and `go build ./...` for the CLI.
Run `npm test`, `npm pack --dry-run`, and `npm run build:portable:current` for the skill.
Use Python 3.10 or newer for its tests.
Run all-target installer dry-run and doctor checks.
Run focused prompt, drift, packet, and skill checks as needed.

Delivery criteria follow source acceptance:

| ID | Delivery criterion |
| --- | --- |
| DC1 | Merge reviewed CLI and skill changes to main. Publish GitHub releases 1.55.0 and 2.19.0 with verified assets. |
| DC2 | Update both formulae in `feci/homebrew-parley`. Verify hosted hashes and local install/test results. |
| DC3 | Open the skill-only WinGet 2.19.0 PR. Verify all three manifests against final hosted assets. Keep CLI WinGet held. |
| DC4 | Install all managed runtime targets, including undetected targets. Install all four generic destinations with the installer. Verify every installed SKILL.md hash. |
| DC5 | A separate fresh Zcode process verifies channel and install evidence. Record pending owner commands and all material limits. |
| DC6 | Write the final release inbox note with channels, deferred items, usage, and ordered owner commands. Save and verify reusable shared memory. |

The generic destinations are `~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck` and `~/.config/opencode/skills/parley-deck`.
Use `parley-deck-skill install --target generic --dest <dir> --force` for each.
List D6 legacy apply first in the final note.
Then list npm and core commands for each still-unpublished version, from 2.15.0 through 2.19.0.
Do not emulate attended publication.

## Idempotence & recovery

Before each hunk, require one exact match of its old text and no existing Appendix B.
On resume, compare hashes before a write. An exact completed result needs no second edit.
Preserve peer artifacts and signoffs byte-for-byte.
Do not overwrite existing release assets or tags without proof that they belong to this release attempt.
Reuse the verified version artifact when a later channel step resumes.
Keep source review and delivery evidence separate.

Use the CLI driver first. Record a real driver gap before the authorized native fallback.
Do not apply the D6 legacy declaration. Do not declare or prune worktrees.
Do not change credentials, keys, or gateway settings.
Limit fix-up to five cycles.
For provider 429/503 or timeout, wait 15 minutes. Allow at most eight relaunches per step.
After a timeout, use 2400 seconds, then 3600 seconds.
A Codex auth or credit failure requires one blocking owner note.
A valid goal FAIL needs resolution; elapsed time never establishes success.

Keep npm and core publication as attended owner commands.
The brief already authorizes GitHub releases, tap updates, skill WinGet PR, and runtime installs.
If an owner-only decision remains open, write one blocking inbox note with a recommendation.
Do not pause at a proposal or permission request that the brief already answers.

## Known risks / de-risking

PRIMARY official evidence identifies Issue 9, January 2025, as current on 2026-10-13.
Zcode independently confirmed the source in round 1.
Attribution of the 20/25/6 numbers to the complete standard remains UNVERIFIED.
They are explicit local working rules. No criterion depends on the unverified attribution.
No participant inspected the full copyrighted dictionary.

A short summary cannot prove linguistic conformance. Manual review protects meaning and exempt evidence.
Prompt tests prove coverage and byte preservation; they do not certify English prose.
The meta-idea facilitator packet exceeds 70,000 bytes before this change. Record the measured value without a wider cleanup.
The CLI Windows CI baseline remains unresolved. Keep CLI WinGet held and state actual validation limits.
The skill release workflow cross-builds Windows assets; this does not establish native Windows validation.

Both participants selected the same small implementation family.
Unanimity among related models is a shared prior, not independent evidence.
Consensus records the three material drafter changes and all twelve alternatives.
It rejects full inline rules, a dictionary, a linter, a new deck key, and historical re-authoring.
Independent tests must challenge exact evidence preservation, old authority, and prompt path coverage.

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/ and ./round-02/
- Controlling brief: ./source-context/ORGANIZER-BRIEF.md
- Exact reviewed payload: ./source-context/protocol-hunks.json
- Source verification: ./source-context/ste-source-notes.md
- Baselines: ./source-context/baseline-sizes.json
- Readiness and mechanical exception receipts: ./source-context/
