---
idea: readme-skill-catalogue
round: 01
agent: codex-1
date: 2026-07-29
---

## 1. The hook

This opening leads with the failure mode, then makes the artifact trail concrete. Proposed
README text:

```markdown
# Parley Deck

> **One model can role-play a committee. Parley Deck requires independent
> participants and keeps the evidence.**

Each participant writes its own canonical file. Round one is written before participants
read one another. Later rounds keep disagreements visible, recorded signoffs gate the final
artifact, and implementation receives a separate review.

The protocol records its state in the repository: kickoff, owned round files, `FINAL.md`,
`IMPLEMENTATION.md`, and review artifacts you can read, diff, and resume.

This package ships five skills: the core cooperation protocol and four focused add-ons for
design, design enforcement, tracker-ready tickets, and concurrent git worktrees.
```

Every product claim in this opening is present in `SKILL.md` or the live
`COOPERATION.md`. “Role-play a committee” is not a benchmark claim; it is the exact solo
failure mode that the core skill's Non-Solo Requirement rejects.

## 2. The catalogue

Use prose subsections, not a table. The names remain scannable as headings, while the prose
has enough room to state the consequential limits. Proposed README text:

```markdown
## Five skills, one package

The installer places all four add-ons alongside the core skill by default. Use
`--no-addons` for the core alone, or `--only <name>[,<name>]` for the core plus a selected
set. Installing an add-on makes it available; invoke it when its job fits.

### [`parley-deck`](./SKILL.md) — make multi-agent work inspectable

Use the core skill when a design, plan, implementation, or review deserves independent
analysis. A facilitator discovers available CLI agents, and every participant writes its
own canonical artifact. One agent does not proxy-write another agent's round, review, or
signoff.

The protocol records kickoff, independent round one, cross-review, consensus, `FINAL.md`,
`IMPLEMENTATION.md`, code review, and fix-up. `track: fast | standard | deliberation`
changes how much ceremony each idea carries. Canonical files remain authoritative whether
the working surface is a local directory, GitHub pull requests, or GitLab merge requests.

### [`parley-design`](./addons/parley-design/SKILL.md) — choose a visual direction without averaging it away

A new visual world needs more than a component checklist. PDS/1.0 makes participants
diverge on directions, critique them, choose one whole, bind it as a contract, apply it,
and audit the result. It is Markdown doctrine with no runtime, network, or framework.

The doctrine will not turn taste into a numeric score. A human Decider receives a findings
ledger, not a prettier-than ranking. It also refuses to ship a house look: interaction
states, contrast, effects, rule precedence, evidence, and waivers are governed, while the
actual visual character must come from the brief and the selected direction.

### [`parley-design-check`](./addons/parley-design-check/SKILL.md) — enforce only what the evidence can prove

`parley-design` states the rules; this add-on runs the checkable subset against artifacts,
DTCG token JSON, stylesheets, and markup. It uses Node built-ins, carries no fallback rule
registry, and emits stable `rule-id — violation — remedy` findings. With no registry it
refuses rule checks and exits `3` instead of manufacturing a pass.

Its capability declaration is generated from its detector modules. This version can obtain
`T0 ARTIFACT` and `T1 SOURCE` evidence. Rules needing rendered or pixel evidence appear as
`UNJUDGEABLE`; a run that judged nothing reportable, or failed to verify a level claim,
exits `4`, not `0`.

### [`parley-tracker`](./addons/parley-tracker/SKILL.md) — write tickets for the business, the builder, and the agent

A tracker should mirror the work, not become a second source of truth. This skill authors
canonical Markdown epics, stories, and subtasks with `At a glance`, `[B] Business`,
`[T] Technical`, and `[A] Agent directives` sections. Acceptance criteria carry audience
tags; the Definition of Done points back to acceptance-criterion IDs and records the
verification command and commit.

`validate` reports the whole readiness gap list, and `claim` refuses to mark an incomplete
ticket `in-progress`. Tracker projection is one-way by default. Pull reconciliation may
write back only fields declared `mirror-owned`; other divergence is a conflict. The skill
defines neutral payloads and mappings for Jira, Linear, GitHub Issues, GitLab, Trello, and
plain boards, but live tracker writes require a separate opt-in connector.

### [`parley-worktrees`](./addons/parley-worktrees/SKILL.md) — isolate concurrent work before it collides

Concurrent agents need separate branches, working trees, indexes, and runtime state. This
add-on gives each implementer a unique branch and sibling worktree, then records ownership
in the `IMPLEMENTATION.md` worktree-allocation table. `FINAL.md` and
`IMPLEMENTATION.md` remain authoritative.

Before a second concurrent worktree is provisioned, its declared file set is compared with
every claimed boundary. An intersection is refused unless an explicit override is
recorded; the normal remedies are to split the boundaries or serialize the shared files
through the integration owner. The add-on does not reimplement `git worktree`. It supplies
the claim, disjointness, integration, environment-isolation, and cleanup discipline around
it.
```

This copy deliberately states the checker's exit code `4`. The round brief lists only
codes `0`–`3`, but the shipped `addons/parley-design-check/SKILL.md` defines `4` for an
`UNJUDGEABLE` run. It also narrows the tracker claim: the shipped tracker skill authors a
payload and mapping, while a separate connector is required for live writes.

## 3. Cuts and line budget

The current file is 401 lines, not 402. The finished README should be **at most 300
lines**. My budget is:

| Part | Maximum lines |
|---|---:|
| Hook, five-skill catalogue, short provenance note | 80 |
| Install, update, installer commands, platform notes | 95 |
| Use examples | 25 |
| Core behavior, agent contract, transports | 55 |
| Layout, repository relationships, status, license | 45 |
| **Total** | **300** |

Make these cuts:

1. Replace current lines 1–42 with the hook and catalogue above. Remove “Inspired by —
   adopted & adapted” from the opening; provenance does not earn space before the reader
   knows what the package contains. If retained, reduce it to a short note near the
   repository-relationship section.
2. Replace the eight prompt blocks at lines 69–117 with two: one “design this” prompt and
   one “continue the current workflow” prompt. Design, implementation, review, and
   transport variations can be single-line substitutions beneath them.
3. Delete lines 121–146, “Why This Exists” and “What The Skill Does.” The hook and core
   catalogue entry now do both jobs with less repetition.
4. Replace lines 148–176 with an accurate compact tree that includes `addons/` and its
   four skill directories.
5. Merge lines 44–63, 178–300, and 302–359 into one installation reference. Keep one
   default install, one update, one project/custom-directory example, the target list,
   and the add-on selectors. Do not print `install --target all --force` in several
   separate sections.
6. Delete the closing claim at line 397 that the protocol's value “should be obvious.”
   Tell the reader what files to inspect and stop there.

## 4. Placement

The catalogue goes **after the hook and before Install**. A command is useful only after
the reader knows what it installs. This ordering also fixes the present discovery failure:
the default installer places five skills, so the README should name all five before asking
the reader to run it. Install remains close enough to the top to be reached after one
catalogue scan.

## 5. Anti-slop test

I applied a **first-sentence fingerprint test** to the five entries: take the first two
words of each entry's prose and fail if any pair repeats, because repeated openings are an
early sign that the entries were filled from one template.

Result: **PASS**. The fingerprints are `Use the`, `A new`, `parley-design states`,
`A tracker`, and `Concurrent agents`. The entries also do not begin with five gerunds or
repeat a fixed “what / why / when” paragraph shape.

## Current README truth audit

These claims are false or unsupported in the current `README.md`:

- **Line 9:** the runtime list reads as exhaustive but omits shipped installer targets
  named later on lines 186 and 289, including Qwen, CodeBuddy, Goose, Kimi, Factory Droid,
  Vibe, Cursor, OpenCode, and AionRS.
- **Lines 21–23:** the lifecycle is not “append-only.” Signoff blocks are append-only;
  `00-prompt.md` status and the living `IMPLEMENTATION.md` are updated in place.
- **Lines 26–27:** the consensus lens does not “rate confidence by agreement.” The shipped
  protocol asks for contradictions, partial coverage, unique insights, and blind spots,
  and explicitly labels that section advisory.
- **Line 119:** “any capable tier-1 model can follow it” has no defined tier, test, or
  shipped compatibility criterion. It is a universal performance claim with no evidence.
- **Lines 148–166:** the repository tree is false by omission because the shipped
  `addons/` directory and all four add-on skills are absent.
- **Line 239:** the Windows binary example says `v1.2.1`; `package.json` says `1.5.0`.
- **Line 371:** the default does not simply use “all discovered installed CLI agents.”
  `SKILL.md` defines a bounded participant set, normally 2–4, with at least one
  non-facilitator when available.
- **Line 397:** “the protocol's value should be obvious” is neither factual nor
  checkable. It is a sales conclusion presented in place of evidence.

The prompt's “seven near-identical prompt blocks” count is also stale: current lines
69–117 contain **eight** fenced prompt blocks.

## Open forks

- **F1 — Catalogue before Install.** The default installation includes all five skills;
  readers should know the payload before executing the command.
- **F2 — Prose entries.** Use linked subsections with scannable headings. A table would
  erase the limits that make these skills credible, especially refusal and
  `UNJUDGEABLE` behavior.
- **F3 — Include the load-bearing ideology.** State selection of one direction, no
  numeric taste score, no house look, registry refusal, and explicit unjudgeability.
  Leave schemas and detector details to the linked skill files.
- **F4 — Treat all five as catalogue peers.** Give every installed skill an equal heading.
  State in the prose that `parley-deck` owns the lifecycle and that the add-ons are loaded
  alongside it; dependency does not justify hiding default-installed capabilities.
- **F5 — Lead with the failure mode.** “One model role-playing a committee” is the
  shortest truthful distinction. The next sentences prove the distinction with owned
  files and independent round one.
- **F6 — Do not change `parley-deck-cli/README.md` in this change.** That README has a
  different product, audience, and fact surface. Updating it is a separate idea with its
  own source audit; coupling it here would make this documentation-only scope harder to
  verify.
- **F7 — Yes, there are additional false claims.** Beyond the four problems named in the
  brief, current lines 21–23, 26–27, and 371 materially misstate shipped behavior. Line
  119 is an unsupported universal claim and line 397 is uncheckable marketing; both
  should also be removed.
