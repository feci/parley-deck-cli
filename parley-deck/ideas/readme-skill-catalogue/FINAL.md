---
idea: readme-skill-catalogue
status: final
drafted-by: claude-1
date: 2026-07-29
consensus: consensus.md (C1–C13; codex-1 🟡 ACCEPT-WITH-RESERVATIONS, hermes-1 ✅ ACCEPT, kimi-1 ✅ ACCEPT)
participants: [claude-1, codex-1, hermes-1, kimi-1]
implementation-target: parley-deck-skill/README.md
---

## What ships

A rewritten `parley-deck-skill/README.md`: **≤ 300 lines** (from 401), opening with the
agreed hook, followed by a five-skill catalogue, with every claim in C9's truth table fixed.

**No file other than `README.md` changes.** No `SKILL.md`, no installer, no protocol.

## Binding section order (C1)

1. Title + hook
2. `## What's in the box` — scan list + five prose entries
3. `## Install` — fastest path only (~8 lines: one `npx`, one `doctor`)
4. `## Use Parley Deck` — **three** prompt blocks maximum (down from eight)
5. `## Install, update, and remove` — the full apparatus
6. `## Local Agent Contract`
7. `## Transports`
8. `## Repository Layout`
9. `## Relationship to other Parley Deck repositories` + provenance
10. `## License`

## Binding copy

The hook and the five catalogue entries below are **the shipping text**. They are
codex-1's round-02 merged copy, whose hook base is claude-1 per the user ruling in C7. The
`<!-- -->` attribution comments are part of the deliverable and must survive into the file.

```markdown
# Parley Deck Skill

<!-- Hook base: claude-1. Grafts: repository-file proof from kimi-1; five-skill close from codex-1. -->

> **One model playing four reviewers is still one model.**

Parley Deck requires separate participants to write their own files, write round
one before reading the others, and cross-review what the others wrote. Disagreements
stay on disk, and recorded signoffs gate what becomes final.

The working state lives in files in your repository that you can read, diff, and
resume — not a chat log you have to trust.

This package includes five skills: the core cooperation protocol and four add-ons
for design, design enforcement, tracker-ready tickets, and parallel worktrees.

## What's in the box

Installing this package installs five skills. The first is the protocol; the other four
build on it. All five install by default — `--no-addons` takes just the core skill,
`--only <name>[,<name>]` picks specific add-ons.

- **`parley-deck`** — the multi-agent cooperation protocol.
- **`parley-design`** — collaborative design that refuses to read as machine-made.
- **`parley-design-check`** — that doctrine's rules, enforced against files on disk.
- **`parley-tracker`** — tickets a stakeholder, a reviewer, and an agent can all read.
- **`parley-worktrees`** — parallel agents over one repo, without silent corruption.

<!-- Base: codex-1. Graft: the closing "more than one model's first answer" line from kimi-1. -->

### [`parley-deck`](./SKILL.md) — make multi-agent work inspectable

Use the core skill when a design, plan, implementation, or review deserves
independent analysis. Every participant owns its canonical artifact; one agent
does not proxy-write another agent's round, review, or signoff.

The protocol records kickoff, independent round one, cross-review, consensus,
`FINAL.md`, `IMPLEMENTATION.md`, code review, and fix-up. The `fast`, `standard`,
and `deliberation` tracks scale the route to the risk. Canonical files remain
authoritative whether the working surface is a local directory, GitHub pull
requests, or GitLab merge requests. Reach for it when the work is worth more than
one model's first answer.

<!-- Base: kimi-1. Grafts: the alongside/never-instead relationship from hermes-1; the bounded-graft constraint from claude-1. -->

### [`parley-design`](./addons/parley-design/SKILL.md) — choose one visual direction without averaging it away

PDS/1.0 makes participants diverge on directions, critique them, choose one whole,
bind it as a contract, apply it, and audit what shipped. It is markdown doctrine
with no runtime, network, or framework; load it alongside `parley-deck`, never
instead of it.

Its refusals are the point: no numeric aesthetic score, no house look, and no
"good default aesthetic" guessed from the category. One direction wins whole;
zero to three bounded grafts may come from losing directions, but none may modify
the winner's token file. Use it for a new visual world, a changed design rule, or
an audit against a ratified contract instead of taste.

<!-- Base: codex-1. Graft: the "says so instead of passing it" line from claude-1. -->

### [`parley-design-check`](./addons/parley-design-check/SKILL.md) — enforce only what the evidence can prove

This add-on runs the checkable PDS/1.0 rules over design artifacts, DTCG token
documents, stylesheets, and markup. It uses Node built-ins, carries no fallback
registry, and emits stable `rule-id — violation — remedy` findings.

With no registry it refuses rule checks and exits `3`. What it cannot decide is
reported `UNJUDGEABLE`; a run that judged nothing reportable, or left a conformance
claim unverified, exits `4`, not `0`. Its capability declaration is generated from
its detector modules, so it says what it cannot check instead of passing it.

<!-- Base: codex-1. Grafts: "the tracker is a mirror" and the migration consequence from claude-1. -->

### [`parley-tracker`](./addons/parley-tracker/SKILL.md) — write tickets for the business, the builder, and the agent

This skill authors canonical markdown epics, stories, and subtasks with `At a
glance`, `[B] Business`, `[T] Technical`, and `[A] Agent directives` sections.
Acceptance criteria carry audience tags, and the Definition of Done points back
to those criteria with verification commands. Its gap-scan reports the full
readiness list; `claim` refuses to mark a ticket `in-progress` when that scan fails.

The tracker is a mirror; the markdown file is canonical. Sync is one-way by
default, and pull reconciliation may write back only fields declared
`mirror-owned`. The skill defines neutral projections for Jira, Linear, GitHub
Issues, GitLab, Trello, and plain boards; live create/update requires an opt-in
connector. Change trackers and you lose a projection, not a requirement.

<!-- Base: kimi-1. Graft: the no-stack-trace failure framing from claude-1. -->

### [`parley-worktrees`](./addons/parley-worktrees/SKILL.md) — isolate concurrent work before it collides

This is protection against the concurrency failure that leaves no stack trace:
two agents writing the same files and producing a result nobody intended. The
branch + worktree + file-set discipline turns that invisible corruption into a
conflict Git can show.

The worktree-allocation table in `IMPLEMENTATION.md` is the lock manifest. Before
a second concurrent worktree is provisioned, its file set is compared with every
claimed boundary; an intersection is refused unless an explicit override is
recorded. Each implementer gets a sibling worktree and isolated runtime state.
Use it when two or more sessions or Phase-5 implementers work in one repository
at once.
```

## Required fixes (C9) — every one is binding

| Location | Fix |
|---|---|
| `:9` runtime list | complete it, or state it is partial — **"fourteen named targets plus a generic directory"**; the string "15 runtimes" MUST NOT appear |
| `:21-23` "append-only" | only signoff blocks are append-only; `IMPLEMENTATION.md` is a living document |
| `:26-27` "rates confidence by agreement" | delete — the protocol defines no confidence rating |
| `:119` "any capable tier-1 model" | delete the universal claim |
| `:148-176` Repository Layout | rewrite the tree to include `addons/` and its four skills, `test/`, `NOTICE.md`, `RELEASING.md` |
| `:239` `v1.2.1` | make the example versionless or use the current version |
| `:242` WinGet "until accepted" | false — publication verified externally (C10); give the real command |
| `:371` "all discovered installed CLI agents" | bounded set, normally 2–4, at least one non-facilitator |
| `:397` "value should be obvious" | delete; end on the inspectable artifacts |
| `:65-117` eight prompt blocks | keep three |

## Constraints carried from consensus

- **C11** — the tracker copy must not imply live tracker writes without the opt-in connector;
  no unmeasured superlative ships.
- **C13** — this change may not promise five skills through an install path that delivers
  one. The `npx skills` discovery defect belongs to `skills-cli-install-path`.
- **C12** — ≤ 300 lines, verified by `wc -l`, not by estimate.
- No invented number anywhere. Counts that ship must be counted first.

## Definition of done

`wc -l README.md` ≤ 300 · every row of the C9 table fixed · the binding copy present verbatim
including attribution comments · `npm test` green · no occurrence of "15 runtimes",
"append-only" applied to the lifecycle, "tier-1", or "should be obvious".
