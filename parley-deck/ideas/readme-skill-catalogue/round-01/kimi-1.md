---
idea: readme-skill-catalogue
round: 01
agent: kimi-1
date: 2026-07-29
---

Lens: **reader model**. Three readers arrive at this README, and they want
different things from the same 60 seconds:

- **The search reader** typed "multi-agent skill" or followed a link from a
  listicle. They owe this page nothing. They need to know, within one
  screenful, what this thing is and whether it is for them. If the first
  screenful is installation flags, they bounce — not because they dislike
  flags, but because they cannot yet tell what would be installed.
- **The colleague-sent reader** got the URL in a chat with "this is what I
  used for the auth redesign." They arrive pre-sold on *trying* it and want
  the install command immediately, then one confirmation that it does what
  their colleague said. Every paragraph between the title and the `npx` line
  is friction for exactly this reader.
- **The installed reader** already ran the installer and is back looking for
  a flag, an update command, or a manual path. They mostly navigate by
  Ctrl-F — *except* on one question: "what else did I install?" Today the
  README gives that question no answer at all, because it never mentions the
  four add-ons the installer put on their disk by default.

The catalogue exists for the first and third reader. The install block exists
for the second. The ordering below serves all three without choosing one.

## 1. The hook (ship-ready)

```markdown
# Parley Deck Skill

> One agent's answer is a first draft. Several agents' recorded agreement is a
> decision.

`parley-deck-skill` installs the Parley Deck cooperation protocol into your
agent runtimes — Codex, Claude Code, Kimi, Cursor, fourteen targets in all —
so that a design, a plan, or a code review is *worked* by several agents
instead of performed by one. Each participant writes its own files:
independent round-1 analyses, cross-reviews, a gated consensus, then
implementation and review as separate phases. Every step is a file in your
repository you can read, diff, and resume — not a chat log you have to trust.

One install gives you five skills: the core protocol plus add-ons for design
doctrine, design enforcement, ticketing, and parallel worktrees.
```

Notes on the hook, from the reader model: the search reader gets the category
("cooperation protocol… several agents") and the differentiator ("files you
can read, diff, and resume") in the first paragraph. The colleague-sent reader
hits the Install section roughly 15 lines in. The installed reader sees the
sentence that tells them the page now answers their question: *five skills*.
The blockquote is a framing claim, not a factual one; "fourteen targets" is
countable against the target list in `Installation Details`.

## 2. The catalogue (ship-ready)

**Shape: prose entries, not a table** (position on F2 below). Each entry is a
`###` heading with a one-line gloss, two to four sentences, and a closing
"Reach for it when…" line. The section opens with a five-line scan list so the
reader in a hurry gets the map before the territory.

Placement in the document: immediately after a minimal Install block, under
the heading `## What's in the box`.

```markdown
## What's in the box

One `npx` install puts five skills into every detected runtime. The first is
the protocol; the other four are add-ons that load alongside it, never instead
of it. All five install by default — `--no-addons` takes just the core skill,
`--only <name>` picks specific add-ons.

- **`parley-deck`** — the multi-agent cooperation protocol.
- **`parley-design`** — collaborative design that refuses to read as machine-made.
- **`parley-design-check`** — that doctrine's rules, enforced against files on disk.
- **`parley-tracker`** — tickets a stakeholder, a reviewer, and an agent can all read.
- **`parley-worktrees`** — parallel agents over one repo, without silent corruption.

### `parley-deck` — the core skill

Teaches your agents the 8-phase idea lifecycle: kickoff, independent round-1
analyses, cross-review, a recorded consensus, `FINAL.md`, `IMPLEMENTATION.md`,
code review, fix-up. It is non-solo by design — one canonical file per agent
per round, so no participant overwrites another — and it never averages
disagreement away: consensus rates confidence by agreement and surfaces blind
spots instead. Canonical files live in your repository; GitHub PRs or GitLab
MRs are mirrors, not the source of truth. Ideas run at `fast`, `standard`, or
`deliberation` rigor depending on what is at stake. Reach for it whenever a
decision, design, plan, or review is worth more than one model's first answer.

### `parley-design` — design doctrine with refusals built in

Ships the PDS/1.0 protocol as pure markdown — zero runtime dependencies, the
whole thing held under a 64 KiB budget that a test enforces. Its positions are
refusals, and they are the point: one direction wins whole (selection, never
averaging); it emits a findings ledger and never a numeric aesthetic score;
and it declines to hand you "a good default aesthetic," because a look
guessable from the category is exactly the failure it exists to prevent.
Reach for it when a Parley idea creates a new visual world, changes a
ratified design rule, or needs an interface audited against a contract
instead of against taste.

### `parley-design-check` — the doctrine, enforced by a tool

Runs the rule registry from `parley-design` against files on disk and emits
findings in one shape — `rule-id — violation — remedy` — stable across runs
and diffable in review. Node built-ins only, no network at check time. What
it cannot decide it reports `UNJUDGEABLE` instead of passing; with no
registry found it refuses rule checks outright (exit 3) rather than fall back
to a bundled copy, and its capability declaration is generated by scanning
its own detectors, so it cannot claim coverage it does not have. Reach for it
in a pre-commit hook or CI step, to verify a conformance-level claim instead
of accepting it, or to get the mechanical findings into the ledger before a
review round argues about the rest.

### `parley-tracker` — tickets three audiences can read

Authors epics, stories, and subtasks as canonical markdown — one file with
`## [B] Business`, `## [T] Technical`, `## [A] Agent directives`, and a
mandatory `## At a glance` — then mirrors them into Jira, Linear, GitHub
Issues, GitLab, Trello, or a plain board. The tracker is a mirror; the file is
canonical. Sync is one-way (file → tracker) by default, `--pull` writes back
only fields the file flags `mirror-owned`, and a field the tracker lacks is
dropped from the mirror, never from the file. A tool-enforced gap-scan
refuses to let an agent claim a ticket that still has holes. Reach for it
when a `FINAL.md` has to become a backlog — or whenever tickets must survive
a tracker migration.

### `parley-worktrees` — parallel agents that don't corrupt each other

The branch + worktree + file-set discipline that turns invisible
concurrent-write corruption into an ordinary git merge conflict. A claim
manifest in `IMPLEMENTATION.md` is the lock layer; the file-set disjointness
check refuses — or demands a recorded override — when two concurrent sessions
claim intersecting files; runtime state, ports, caches, and env are isolated
per worktree; worktrees live in a sibling directory, never inside `.git/`.
Reach for it when two or more sessions or Phase-5 implementers touch the same
repository at once.
```

## 3. Cuts, with a line budget

Current README: 402 lines. Target finished README: **≤ 300 lines**, of which
the catalogue is ~75. The catalogue is paid for by cuts, not by growth.

| Cut | Lines today | Action | Saved |
|---|---|---|---|
| "Use Parley Deck" seven prompt blocks (~65–117) | ~50 | Keep three blocks (design, implement, continue); drop compare/ship-e2e/PR-transport/quick-architecture — they are variations of the same shape | ~28 |
| "What The Skill Does" (133–146) | 14 | Cut entirely; the core-skill catalogue entry now carries this content | ~14 |
| `install --target all --force` repeated across Installation Details + Updating (178–359) | — | Merge into one "Install, update, remove" section; the command appears once for install, once for update, once per genuinely different path | ~35 |
| "Status" section (393–397) | 5 | Compress to two sentences | ~3 |
| "Repository Layout" (148–176) | 29 | **Fix, don't cut**: add `addons/`, drop per-file bullets that restate names | ~8 |

Net: roughly −88 lines of cuts against +75 lines of catalogue and +8 lines of
corrected layout. Nothing cut is load-bearing for a reader decision; the
prompt blocks being dropped are the ones a newcomer skims past anyway.

## 4. Placement, argued from the reader model

Order: **hook → Install (fastest path + verify only, ~15 lines) → What's in
the box → Use → Why This Exists → everything else** (details, commands,
updating, contract, transports, layout, license).

This is a deliberate middle position on F1, and it is not a hedge — it is a
specific claim about the three readers:

- The colleague-sent reader is the only one whose patience is measured in
  lines, and the full install apparatus (details, per-target commands,
  Windows, Homebrew) is what would bury the catalogue for the other two. So
  the *fastest-path* install goes above the catalogue — 15 lines, ending in
  `doctor` — and the *apparatus* goes below it.
- The search reader then meets the catalogue at roughly line 30: still the
  first screenful-and-a-half, early enough to answer "what is all this"
  before they decide the page is a reference manual.
- The installed reader looking for a flag navigates by Ctrl-F regardless of
  order; what order *can* do for them is surface the add-ons they never knew
  they had, and a catalogue at line 30 does that every time they revisit.

Catalogue-first (before any install) asks the colleague-sent reader to read 90
lines they did not come for. Install-first (all of it) asks the search reader
to scroll past 200 lines of flags to learn what the flags install. Both lose a
reader the page could have kept.

## 5. Anti-slop test, applied to my own draft

The check: **the opening-construction audit plus the banned-lexicon grep.**
Concretely — (a) grep the draft against the brief's banned list
("revolutionise", "seamlessly", "supercharge", "unlock", "game-changing",
"empower", "leverage", "cutting-edge", "robust", "in today's fast-paced
world") and require zero hits; (b) lay the opening clause of each catalogue
entry side by side and require that no two share the same grammatical
construction (the gerund-opening failure mode); (c) require that every
factual sentence is pointable at a file I read in step 1.

Result of applying it to sections 1 and 2 above:

- **(a) Banned lexicon: 0 hits.** Checked by eye against the list, twice.
- **(b) Opening constructions, all five distinct:** core opens with a verb
  ("Teaches your agents…"); design with a verb on a different subject plus an
  em-dash appositive ("Ships the PDS/1.0 protocol as pure markdown —");
  design-check with a verb plus an em-dash finding shape ("Runs the rule
  registry… — `rule-id — violation — remedy` —"); tracker with a verb plus a
  parenthetical ("Authors epics, stories… — then mirrors them");
  worktrees opens with a noun phrase, no main verb in sentence one ("The
  branch + worktree + file-set discipline that turns…"). Two entries (design,
  tracker) share the "appositive between dashes" device — flagged, kept
  deliberately: the appositives carry different content (budget vs. audience
  tags) and the devices are not in adjacent entries.
- **(c) Fact audit: every claim traces to a file read for this round.**
  "Fourteen targets" — counted from `Installation Details`. "64 KiB budget
  enforced by a test" — `parley-design` brief entry, checkable in its
  references. "Exit 3", "`UNJUDGEABLE`", "capability generated by scanning
  detectors" — `parley-design-check/SKILL.md` §Exit codes, §Capability,
  §The registry contract. "`mirror-owned`", "gap-scan refuses claim" —
  `parley-tracker/SKILL.md` §Core Rule, §No-assumption gap-scan.
  "Refuses on intersecting file sets", "sibling directory, never `.git/`" —
  `parley-worktrees/SKILL.md` §4, §2. Track names and transports — core
  `SKILL.md` and current README. No number, quote, or testimonial appears
  that is not countable in a shipped file.

## Forks

**F1 — Catalogue leads or follows Install?** Follows the *fastest-path*
install, precedes the *full* install apparatus (see item 4). If forced to a
binary: leads. The install default means nobody needs to choose skills before
installing, so the catalogue's only pre-install job is persuasion — and the
hook already did that job.

**F2 — Table vs. prose?** Prose entries, with a five-line scan list on top.
A table loses "no numeric aesthetic score, ever" — the cell either truncates
it into meaninglessness or wraps until the table stops scanning. The actual
scanning need is real but smaller than a table: five lines, one clause each,
which the scan list provides. Prose loses less.

**F3 — How much ideology?** Exactly one load-bearing position per add-on,
stated as a refusal or a guarantee, then stop. Design: selection-never-
averaging plus the refused default aesthetic. Design-check: `UNJUDGEABLE`
instead of passing. Tracker: the tracker is a mirror. Worktrees: corruption
becomes an ordinary merge conflict. One idea each is what makes the entries
memorable; the second idea is what would make them long. Everything past one
belongs in the SKILL.md, linked by the runtime's own skill loading.

**F4 — Peers or satellites?** Satellites in framing, peers in visual weight.
The dependency direction is the reader-relevant truth — every add-on's own
SKILL.md says "load alongside `parley-deck`, never instead of it" — so the
section intro says so in one sentence. But the entries themselves get
identical heading level and comparable length, because the install default
(all four, unasked) means the installed reader owns them whether or not the
docs grant them dignity.

**F5 — Hook leads with protocol, failure mode, or artifact trail?** Failure
mode, closed by the artifact trail. "One agent's answer is a first draft" is
the sentence all three readers finish in their own heads; "files you can
read, diff, and resume" is the proof offered before they can doubt it.
Leading with the protocol ("installs the X cooperation protocol") spends the
first line on a category label; leading with the artifact trail alone
("every step is a file") describes a mechanism to a reader who does not yet
know why they want one.

**F6 — Touch `parley-deck-cli/README.md` in the same change?** No. Separate
idea. The brief names one target file; the CLI README has a different job
(orchestration commands, not a skill catalogue) and a different reader.
Bundling them doubles the review surface and guarantees one of the two gets
reviewed lazily. File a follow-up idea after this one lands.

**F7 — Another stale claim?** Yes: the intro paragraph (current README
line 9) says the skill installs into "Codex, Claude Code, Antigravity,
Gemini, Hermes, or a custom skill directory" — five targets — while
`Installation Details` (line 186) lists fourteen native targets including
Qwen, Kimi, Cursor, Goose, and OpenCode. The first factual sentence a new
reader meets under-sells the install surface and contradicts the page's own
later section. Fixed in my hook ("fourteen targets in all").
