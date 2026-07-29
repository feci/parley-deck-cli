---
idea: readme-skill-catalogue
author: user
created: 2026-07-29
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1]
target-repo: parley-deck-skill
target-file: README.md
status: round-01
roles:
  claude-1: structure & voice (what shape the catalogue takes, what register it speaks in)
  codex-1: truthfulness (every sentence must be checkable against a shipped file)
  hermes-1: the hook (why a stranger keeps reading past line one)
  kimi-1: reader model (who arrives, what they came for, where they bounce)
---

## Problem / idea

`parley-deck-skill` now ships **five skills** — the core `parley-deck` skill plus four
opt-in add-ons — and `README.md` **mentions none of the four**. A reader who installs the
package gets four skills nobody told them about.

The user's ask, verbatim (translated): *"in the description of that skill, put more
interesting descriptions of all the skills that are in there into README.md, and write it so
that it grabs people the moment they start reading."*

So there are two jobs, and they are **not** the same job:

1. **A skill catalogue** — all five skills described so a reader understands what each one
   is *for* and when they would reach for it.
2. **A hook** — the opening of the README has to earn the second paragraph.

Propose both. Say explicitly which parts of the current README you would **cut**, because
adding 200 lines to a 402-line README is not an improvement.

## The constraint that makes this hard

This package ships `parley-design`, whose entire purpose is preventing output that reads as
machine-made, and whose **honesty rule** forbids inventing a metric, benchmark, testimonial,
customer logo, rating, or award.

**A README for that package that reads like AI marketing copy is a self-refuting artifact.**

Concretely, this rules out: "revolutionise", "seamlessly", "supercharge", "unlock",
"game-changing", "in today's fast-paced world", three-adjective stacks, an em-dash-heavy
rhythm that never varies, invented adoption numbers, invented benchmarks, fake quotes, and
any claim of a benefit no shipped file delivers. It also rules out the *shape* that reads as
machine-made: five paragraphs of identical length and identical construction, each opening
with a gerund.

A hook is allowed to be vivid. It is not allowed to be false.

## Ground truth — the five skills (cite these; do not invent)

Read the actual files before writing. Paths are relative to the `parley-deck-skill` repo.

### 1. `parley-deck` — the core skill (`SKILL.md`, 50 KB)

Installs the Parley Deck cooperation protocol into an agent runtime. 8-phase idea lifecycle:
kickoff → independent round-1 → cross-review → consensus → `FINAL.md` → `IMPLEMENTATION.md`
→ code review → fix-up. Non-solo by design: one canonical file per agent per round, no agent
overwrites another. Transports: `local-dir`, `github-pr`, `gitlab-mr`. Conditional rigor via
`track: fast | standard | deliberation`. Installs into Codex, Claude Code, Antigravity,
legacy Gemini, Hermes, Qwen, CodeBuddy, Goose, Kimi, Factory Droid, Vibe, Cursor, OpenCode,
AionRS. Dependency-free Node installer.

### 2. `addons/parley-design` — design doctrine (`SKILL.md` 6.5 KB + 3 references)

Ships the **PDS/1.0** protocol as pure markdown: **zero runtime dependencies**, hard **64 KiB
budget across exactly four files, enforced by a test**. Typed design artifacts
(`DESIGN-BRIEF`, `DIRECTION`, `CRITIQUE`, `VERDICT`, `CONTRACT`, `DESIGN-SYSTEM`, `AUDIT`,
`WAIVERS`), gates G1–G4, conformance levels L1–L4, evidence tiers `T0 ARTIFACT · T1 SOURCE ·
T2 RENDERED · T3 PIXEL`, rule classes `quality` / `slop` / `system`, waivers with expiry and
counter-signature.

Its load-bearing positions, all verifiable in the shipped files:
- **Selection, never averaging** — one direction wins whole; 0–3 grafts; a graft MUST NOT
  modify the winner's token file.
- **No numeric aesthetic score, ever.** It emits a findings ledger, not a verdict. The
  Decider is a human by default; an unattended run records `ABSTAIN` and stops before
  implementation.
- **Ships invariants, not a house look.** "Give me a good default aesthetic" is a request it
  refuses on purpose — an aesthetic guessable from the category is the failure it exists to
  prevent.
- The contrast floor is **system-blind**: widening your own ratified ramp does not satisfy
  it, because otherwise an implementer legalises its own output by editing the system.
- One **literate `RULES.md`** where fenced ```pds-rule``` YAML is the machine source and the
  surrounding prose is the human source — no generated second copy, so nothing can drift.

### 3. `addons/parley-design-check` — the enforcement layer

Node built-ins only. No dependencies, no framework, no network at check time. Findings in one
shape: `rule-id — violation — remedy`, stable and diffable across runs. Verdicts
`PASS | VIOLATION | NEEDS_REVIEW | UNJUDGEABLE`. Exit codes: `0` clean, `1` findings,
`2` the run itself failed, **`3` rule checks refused because no registry was found**.

The properties worth a reader's attention:
- It reads the registry from the installed `parley-design` and **carries no fallback copy** —
  absent that skill it refuses and says so, rather than passing.
- Its capability declaration is **generated by scanning its own detectors**, so it cannot
  claim coverage it does not have.
- v1 judges `T0 ARTIFACT` and `T1 SOURCE`. `T2` and `T3` are reported `UNJUDGEABLE` — named,
  never silently skipped.
- Most of what makes work read as machine-made is not decidable from source, **and it says so
  instead of passing it**.

### 4. `addons/parley-tracker` — tickets that three audiences can read

Epics, user stories, and technical subtasks as canonical markdown, mirrored into Jira, Linear,
GitHub Issues, GitLab, Trello, or a plain board. **The tracker is a mirror; the markdown file
is canonical.** Sync is one-way (file → tracker) by default; `--pull` writes back only fields
the file flags `mirror-owned` (default `status`, `assignee`) — anything else surfaces as a
conflict, not a silent overwrite. Graceful degradation: a field the target tracker lacks is
dropped from the *mirror*, never from the *file*.

One file, three audiences: `## [B] Business`, `## [T] Technical`, `## [A] Agent directives`,
plus a mandatory `## At a glance`. Audience-tagged acceptance criteria; a Definition-of-Done
checklist with one box per AC id and a verify command; a **tool-enforced no-assumption
gap-scan** before work starts; explicit `## Non-goals` to stop AI scope creep.

### 5. `addons/parley-worktrees` — parallel agents that don't corrupt each other

Allocate, name, isolate, merge, and clean up git worktrees so several sessions or Phase-5
implementers work one repo at once. Its own stated claim: the value is not "agents can use
worktrees" but the **branch + worktree + file-set discipline that turns invisible
concurrent-write corruption into an ordinary git merge conflict**. A coordination manifest is
the lock layer; the add-on **refuses** (or demands a recorded override) when two concurrent
sessions declare **intersecting file sets**. Worktrees live in a sibling directory, never
inside `.git/`. Runtime state, ports, caches, and env are isolated per worktree.

### Installer facts (verified from `--help`)

All add-ons install by default alongside the core skill.
`--no-addons` installs the core skill only. `--only <name>[,<name>]` installs the core skill
plus the named add-on(s). Available: `parley-design`, `parley-design-check`,
`parley-tracker`, `parley-worktrees`.

## Current README — what is already there

402 lines / 14 KB. Sections in order: title + pitch blockquote · "What the protocol gives
your agents" · "Inspired by — adopted & adapted" · Install · Use Parley Deck (seven
copy/paste prompt blocks) · Why This Exists · What The Skill Does · Repository Layout ·
Installation Details · Installer Commands · Updating · Local Agent Contract · Transports ·
Relationship To Other Parley Deck Repositories · Status · License.

Known problems beyond the missing catalogue:
- **`Repository Layout` is stale** — it shows no `addons/` directory at all.
- **Seven near-identical prompt blocks** in "Use Parley Deck" is a wall; a newcomer skims past
  all of them.
- `Installation Details` and `Updating` between them contain the same `install --target all
  --force` line many times over.
- The version string in the Windows example is `v1.2.1`; the package is at `1.5.0`.

## What each participant writes in round-01

Write `round-01/<your-agent-id>.md` **independently — do not read the others' files first**.
Cover, in your own order:

1. **The hook.** Write the actual opening — title area through the first ~10 lines, as it
   would ship. Not a description of a hook; the hook.
2. **The catalogue.** Its shape (table / cards / subsections / something else) **and a
   drafted entry for all five skills**, in the register you are arguing for. One-line
   descriptions are not enough to judge a voice; write real sentences.
3. **Cuts.** What comes out, with a line budget. State the target size of the finished README.
4. **Placement.** Where the catalogue sits relative to Install, and why that order serves the
   reader who just arrived.
5. **Your anti-slop test.** Name one concrete check you would apply to your own draft that
   would catch it reading as machine-made — and then report the result of applying it.

## Open forks — take a position, don't hedge

- **F1** Does the catalogue lead (before Install) or follow it? A reader who has not
  installed anything may not care about five skills; a reader who has installed them does not
  know they exist.
- **F2** Table vs. prose entries. A table scans; prose persuades. Five rows of a table cannot
  carry "no numeric aesthetic score, ever". Which loses less?
- **F3** How much of the design add-ons' *ideology* belongs in a README, versus a link to the
  skill? The positions are what make them interesting and also what makes them long.
- **F4** Are the four add-ons peers of the core skill, or satellites? The install default
  (all four, unasked) argues peer; the dependency direction argues satellite.
- **F5** Does the hook lead with the protocol, with the failure mode it prevents, or with the
  artifact trail a reader can go and read?
- **F6** Do we touch `parley-deck-cli/README.md` in the same change, or is that a separate
  idea? Argue scope, don't assume it.
- **F7** Is there a claim in the current README that is no longer true, beyond the four listed
  above? Say which line.

## Non-goals

- No change to any `SKILL.md`, to `COOPERATION.md`, or to installer behaviour. This idea
  changes documentation.
- No new skill, no renamed skill.
- No invented metric, benchmark, testimonial, logo, rating, award, adoption number, or quote.
  A labelled hole is honest; a plausible fabricated number is not a style choice.

## Definition of done

`parley-deck-skill/README.md` describes all five skills accurately and opens in a way that
earns the reader's second paragraph; every factual claim is checkable against a shipped file;
the stale claims listed above are fixed; the README is not longer than it needs to be.
