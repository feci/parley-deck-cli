---
agent: claude
idea: meta-protocol-change-agent-teams-patterns
review-round: 1
date: 2026-05-14
reviewed-commit: 181fd08
---

## Summary

The protocol text changes in the four target files (parley-deck-cli COOPERATION.md, parley-deck-skill SKILL.md, parley-deck-skill references/COOPERATION.md, and the installed copies under `~/.codex/skills/parley-deck/`) faithfully render the five user-approved recommendations from FINAL.md: advisory `roles:` map, internal-helper accountability, 2-4 participant sizing guidance, Phase 5 plan-gate checklist, and non-authoritative inbox handoffs. Wording is vendor-neutral, internally consistent, and the protocol-changelog.md entry is appropriately scoped. The protocol-content side of this implementation is sound; the gaps are in the process/audit-trail surface around how the implementation was committed and described.

## Findings

### [MAJOR] IMPLEMENTATION.md head-commit and transport metadata do not match the working tree

`IMPLEMENTATION.md` frontmatter records `branch: /Users/tomasfecko/AI_WORKSPACE/parley-deck/parley-deck-cli#main`, `head-commit: 181fd08`, `design-pr: n/a`, and `implementation-pr: n/a`. But `181fd08` is the head of the unrelated `tui-agent-controls` work (`[codex] tui-agent-controls: implementation`), and `git status` shows the entire implementation for this idea — modified `parley-deck/COOPERATION.md`, modified `parley-deck/meta/protocol-changelog.md`, the new `parley-deck/ideas/meta-protocol-change-agent-teams-patterns/` tree, and the new `parley-deck/inbox/codex-to-all_meta-protocol-change-agent-teams-patterns_prompt-path-repair.md` — uncommitted. The active transport in `COOPERATION.md` is `github-pr`, which under §11.B Phase 5 requires a `feature/<slug>` branch in the code repo and an implementation PR; neither exists.

Why it matters: a core Parley Deck invariant is that the audit trail is durable and verifiable. A `head-commit:` value that points at a different idea's SHA is actively misleading for any future reader (or reviewer pinning a finding to a commit), and "n/a" PRs under a `github-pr` transport with no recorded transport exception leaves a quiet protocol-compliance gap. It also means a working-tree loss right now would lose the entire implementation.

Suggested fix: commit the implementation on a `feature/meta-protocol-change-agent-teams-patterns` branch in `parley-deck-cli`, open the implementation PR (and the design PR if not already merged), and update `IMPLEMENTATION.md` frontmatter to the real SHA and PR URLs. If the user explicitly authorized running this protocol-change idea directly on `main` without a PR, record that exception in `inbox/` and as a brief "Deviations from FINAL.md" entry, and at minimum replace `head-commit: 181fd08` with the SHA that actually contains the implementation (committing `IMPLEMENTATION.md` is a prerequisite for that SHA to exist).

### [MINOR] Phase 5 plan-gate wording creates a timing ambiguity in `IMPLEMENTATION.md` creation

In `parley-deck/COOPERATION.md` §4 Phase 5, the bullet "For multi-file changes or changes outside `parley-deck/`, records a short implementation plan/checklist in `IMPLEMENTATION.md` before making code changes" is immediately followed by "On completion, creates `ideas/<slug>/IMPLEMENTATION.md`:". The first sentence assumes the file exists before code; the second sentence describes it as created on completion. The same tension exists in the SKILL.md Phase 5 prompt block.

Why it matters: FINAL.md is clear that the plan-gate is a pre-code checklist, but a future implementer reading these two adjacent sentences may either (a) postpone the plan until the end (defeating the gate) or (b) create the file twice, once with the plan and once "on completion", and wonder which is correct.

Suggested fix: change "On completion, creates" to "On completion, finalizes" (or "completes") in `COOPERATION.md` Phase 5, and add one half-sentence to the plan-gate bullet noting that `IMPLEMENTATION.md` may be opened with just the plan/checklist section first and the remaining sections filled in at completion. Mirror the wording in the SKILL.md Phase 5 prompt.

### [NIT] "Deviations from FINAL.md: None" understates the SKILL.md / installed-skill scope expansion

`FINAL.md` enumerates only `COOPERATION.md` text changes (and its Tests section is written for "a later implementation PR that edits `COOPERATION.md` or CLI behavior"). The implementation also propagated equivalent edits into `parley-deck-skill/SKILL.md`, `parley-deck-skill/references/COOPERATION.md`, and the two `~/.codex/skills/parley-deck/...` installed copies. The user's instruction to update the active skill is recorded under `IMPLEMENTATION.md` "Notes for reviewers", but the "Deviations from FINAL.md" field reads "None".

Why it matters: future readers scanning the Deviations section will assume the implementation's blast radius matches FINAL.md's enumeration. A reader investigating drift between live protocol and bundled fallback months from now should not have to cross-reference the Notes section to learn that the skill artifacts were touched as part of this idea.

Suggested fix: move (or duplicate) the SKILL.md/installed-skill propagation note from "Notes for reviewers" into "Deviations from FINAL.md" framed as a documented scope extension — e.g., "Scope extended per user instruction to keep `parley-deck-skill/SKILL.md`, `parley-deck-skill/references/COOPERATION.md`, and the installed `~/.codex/skills/parley-deck/` copies in sync with the new protocol text. Protocol intent unchanged."

## Open questions

- Was the implementer asked by the user to ship this idea without a feature branch / PR (i.e. an authorized transport exception for this protocol-change idea), or is the uncommitted working-tree state simply because the review was launched before the close-idea commit? Either way the audit-trail metadata in `IMPLEMENTATION.md` needs to be reconciled.
- Should the bundled `references/COOPERATION.md` snapshots retain the `Transport: github-pr` header that they inherit from the live file, or should the fallback always carry a transport-agnostic header so a new project that bootstraps from the skill is prompted to make the §0 choice deliberately? Out of scope for this review round but worth noting for a future protocol-change idea.
