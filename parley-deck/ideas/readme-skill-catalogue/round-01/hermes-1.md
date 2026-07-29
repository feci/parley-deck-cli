---
idea: readme-skill-catalogue
round: 01
agent: hermes-1
date: 2026-07-29
---

## 1. The hook — the actual opening

> Most multi-agent workflows fail in predictable ways: one agent anchors the rest
> before they form their own view, disagreements vanish inside a long chat history,
> implementation starts before there is real consensus, and reviews are informal,
> unowned, and hard to resume. Parley Deck turns that conversation into files — one
> per agent per round, every step explicit, consensus gated, and resumable from the
> repository alone.

Five skills ship in this package, and the current README mentions one of them. The
core skill installs the cooperation protocol: a track-gated eight-step lifecycle
from kickoff through fix-up, where each participant writes its own analysis,
cross-reviews the others, signs off on consensus, and reviews the implementation —
every step a file you can read, diff, and resume. Four add-ons install alongside
it by default: a design doctrine that refuses to ship a house look, a checker that
reports what it cannot judge instead of passing it, a ticket format three audiences
can read, and a worktree discipline that turns invisible concurrent-write corruption
into ordinary merge conflicts.

---

The blockquote above replaces the current one:

> **Install the Parley Deck cooperation protocol into your AI agents** — real
> multi-agent deliberation with a durable audit trail, not one model role-playing
> a committee.

What the current blockquote fails to do: it describes the package from the
seller's side ("install the protocol into your agents") rather than from the
reader's side. It tells you what the product *is*, not what problem arrives at
your door and why you would reach for it. The "not one model role-playing a
committee" clause is a good instinct — it names a failure mode — but it is buried
as a parenthetical after a product claim, and it is the only failure mode named.
A reader who has not yet felt that specific pain reads "real multi-agent
deliberation" as a feature bullet and moves on. The hook needs to surface the
failure modes first, because those are the thing the reader already recognises,
and then show the artifact trail as the answer — not as a claim about the product,
but as a concrete thing the reader can go and read.

The replacement leads with the failure modes (taken verbatim in spirit from the
current README's own "Why This Exists" section, lines 123–129), then states what
Parley Deck does about them, then names all five skills in the second paragraph so
the reader knows the scope before they hit Install. The first blockquote earns the
second paragraph by making the reader think "yes, that is exactly what goes wrong"
before it says "here is what turns it into something you can audit."

---

## 2. The catalogue — shape and all five entries

Shape: five subsections, not a table. A table scans, but five rows cannot carry
"no numeric aesthetic score, ever" or "it refuses rather than passing" without
shrinking those positions to a cell that reads as a feature. Prose entries let
each skill's one load-bearing idea land as a sentence, and the register varies
because the skills vary — the core skill is infrastructure, the design doctrine is
ideological, the checker is operational. Each entry is two to four sentences: what
it is, what it is for, and the one thing about it that a reader would not guess
from the name.

The catalogue sits under a single `## Skills` heading, core skill first, add-ons
after. No "peers vs. satellites" framing in the prose — the order and the one-line
relationship note in each add-on entry does that work without announcing the
decision (see F4 below).

### `parley-deck` — the core skill

The cooperation protocol itself. It runs an eight-step lifecycle — kickoff,
independent analysis, cross-review, consensus, `FINAL.md`, `IMPLEMENTATION.md`,
code review, fix-up — where each participant writes its own file per round and no
agent overwrites another. Three tracks (`fast`, `standard`, `deliberation`) scale
ceremony to risk: a trivial reversible change skips cross-review, while a protocol
or security change runs the full lifecycle with all non-implementers reviewing.
It is transport-agnostic (local files, GitHub PRs, GitLab MRs) and installs into
fourteen runtimes including Codex, Claude Code, Hermes, and Goose, with a
dependency-free Node installer. The skill teaches an agent how to *participate*;
the companion `parley` CLI orchestrates the runs.

### `parley-design` — design doctrine

An opt-in add-on for collaborative design-system work: diverging on directions,
critiquing them, choosing one whole, binding it as a contract, applying it, and
auditing what shipped. It ships the PDS/1.0 protocol as pure markdown with zero
runtime dependencies — four files under a hard 64 KiB ceiling, enforced by a test.
Its load-bearing position is that it emits no score and no ranking: the Decider is
a human by default and receives a findings ledger, not a verdict. A request for "a
good default aesthetic" is one it refuses on purpose, because a look guessable
from the category is the failure it exists to prevent. Load it alongside the core
skill, never instead of it; it is a profile over the protocol's phases, not a
second phase machine.

### `parley-design-check` — the enforcement layer

The companion checker for `parley-design`. It reads the rule registry that skill
ships, runs detectors against files on disk, and emits findings in one shape:
`rule-id — violation — remedy`. Node built-ins only — no dependencies, no
framework, no network. It carries no fallback copy of the registry: absent
`parley-design`, it refuses rule checks and exits 3 rather than passing. Its
capability declaration is generated by scanning its own detectors, so it cannot
claim coverage it does not have. Most of what makes work read as machine-made is
not decidable from source, and it says so — reporting `UNJUDGEABLE` instead of
silently skipping. v1 judges artifacts and source; rendered and pixel evidence are
named `UNJUDGEABLE`, never passed.

### `parley-tracker` — tickets that three audiences can read

An opt-in add-on for authoring epics, user stories, and technical subtasks as
canonical markdown files, then mirroring them into Jira, Linear, GitHub Issues,
GitLab, Trello, or a plain board. The tracker is a mirror; the markdown file is
canonical. Sync is one-way (file → tracker) by default, and a `--pull` writes back
only fields the file flags `mirror-owned` — anything else surfaces as a conflict,
not a silent overwrite. Each ticket is one file with three audience-tagged
sections (`## [B] Business`, `## [T] Technical`, `## [A] Agent directives`) plus a
mandatory `## At a glance` block. A tool-enforced gap-scan runs before an agent can
claim a ticket, refusing `in-progress` status on any incomplete field — the single
highest-leverage rule for AI output quality, because a template alone cannot
guarantee an agent stops.

### `parley-worktrees` — parallel agents that don't corrupt each other

An opt-in add-on for allocating, naming, isolating, merging, and cleaning up git
worktrees so several sessions or Phase-5 implementers can work one repo at once.
Its stated claim is that the value is not "agents can use worktrees" but the
branch + worktree + file-set discipline that turns invisible concurrent-write
corruption into an ordinary git merge conflict. A coordination manifest in
`IMPLEMENTATION.md` is the lock layer; the add-on refuses, or demands a recorded
override, when two concurrent sessions declare intersecting file sets. Worktrees
live in a sibling directory, never inside `.git/`. Runtime state — ports, caches,
env, submodule checkouts — is isolated per worktree. It does not re-implement
`git worktree`; it adds the claim, the disjointness check, and the lifecycle the
runtimes do not provide.

---

## 3. Cuts — what comes out

Target: a README of roughly 280–320 lines (down from 402). The catalogue adds
~80 lines; the cuts below remove ~160–180.

| Section | Lines | Action | Reason |
|---|---|---|---|
| Use Parley Deck (seven prompt blocks) | 65–119 | Cut five of seven blocks; keep one design prompt and one implement prompt | Seven near-identical blocks are a wall; a newcomer skims past all of them. Two examples show the two main verbs (design, implement); the rest are variants the reader can derive. |
| Installation Details | 178–270 | Collapse to ~15 lines: keep `--target all`, `--scope project`, `--dest`, `--include-undetected`, and the manual paths block; cut the single-target permutations and the `install --target codex`, `--target claude`, etc. blocks | The single-target blocks are the same command with one word changed. The manual paths block is the only part a reader cannot derive from `--help`. |
| Updating | 302–359 | Cut to ~10 lines: keep the `install --target all --force` one-liner, the `--dry-run` preview, and the `doctor` check; remove the per-target update blocks and the Homebrew/Gemini/global variants | The per-target update blocks repeat the install blocks with `--force` appended. One update command covers every target. |
| Installer Commands | 274–300 | Fold the useful-flags list into the collapsed Installation Details; remove the standalone section | It duplicates `--help` output. |
| Repository Layout | 148–176 | Rewrite to include `addons/` and its four subdirectories; cut the per-file prose bullets that repeat what the catalogue already says | The current layout is stale (no `addons/`). The fix is a corrected tree, not eight bullets restating the catalogue. |
| What The Skill Does | 133–146 | Cut entirely | The catalogue's core-skill entry and the "Why This Exists" section already cover this. Three sections saying the same thing is what makes a README feel machine-made. |
| Inspired by — adopted & adapted | 31–42 | Keep as-is | This is lineage and attribution; it is short, specific, and honest. It does not read as slop. |

Also fix:
- The stale version string `v1.2.1` on line 239 → `v1.5.0` (verified from `package.json`).

---

## 4. Placement — where the catalogue sits

The catalogue goes **before** Install, immediately after the hook and a one-paragraph
"What the protocol gives your agents" section (kept, trimmed). The order a arriving
reader needs:

1. **Hook** — the failure modes and the artifact trail (the blockquote above).
2. **What the protocol gives your agents** — the four-bullet list from the current
   README (lines 21–29), trimmed. This is the protocol's value proposition in
   scannable form, and it bridges the hook into the catalogue.
3. **Skills** — the five-entry catalogue. This is the new section.
4. **Install** — the fastest-path command, then the collapsed details.
5. **Use Parley Deck** — two prompt blocks, not seven.
6. **Why This Exists** — kept as-is (lines 121–131); it is the longer-form version
   of the hook and earns its place after the reader knows what the package contains.
7. Everything else (Transports, Relationship, Status, License) stays in its
   current order.

Why before Install: a reader who has not installed anything needs to know what
they are installing before they install it. The current README puts Install at
line 44, which means a reader who came in from an `npx` command sees the install
command before they know what the four add-ons are. The catalogue before Install
solves the brief's core problem — "a reader who installs the package gets four
skills nobody told them about" — at the moment of decision, not after it. A reader
who has already installed and is coming back to understand what they got finds the
catalogue in the same place, because nobody scrolls past Install to discover
features.

---

## 5. Anti-slop test

The test: read each catalogue entry aloud and check whether it opens with a gerund
("Shipping", "Providing", "Enabling", "Allowing") or a three-adjective stack, and
whether all five entries have the same sentence-length rhythm. The brief names the
gerund-open and identical-construction patterns as the shape that reads as
machine-made.

Result of applying it:

- `parley-deck` opens with "The cooperation protocol itself." — a noun phrase, no
  gerund.
- `parley-design` opens with "An opt-in add-on for collaborative design-system
  work:" — a noun phrase with a colon, no gerund.
- `parley-design-check` opens with "The companion checker for `parley-design`." —
  a noun phrase, no gerund.
- `parley-tracker` opens with "An opt-in add-on for authoring epics, user stories,
  and technical subtasks" — a noun phrase, no gerund.
- `parley-worktrees` opens with "An opt-in add-on for allocating, naming,
  isolating, merging, and cleaning up git worktrees" — a noun phrase, no gerund.

Two of the five open with "An opt-in add-on for…" — that is a structural repeat.
Varying it: `parley-design` could open "Design-system work as a Parley Deck
profile:" and `parley-tracker` could open "Tickets that survive a tracker
migration." I would apply that variation before shipping. The sentence-length
rhythm varies: the core entry runs long (it covers tracks, transports, and
runtimes), the design entry is mid-length with a short punchy refusal sentence,
the checker entry is the shortest (it is the most operational), the tracker entry
is mid-length, and the worktrees entry ends with a short sentence that names what
it does not do. No entry uses a three-adjective stack. No entry uses "seamlessly",
"revolutionise", "supercharge", "unlock", or "game-changing". No invented numbers
appear — every factual claim traces to a file I read.

---

## Forks

**F1 — catalogue before Install or after?**
Before. A reader who has not installed needs to know what the four add-ons are
before they run `npx … install --target all` and get all of them by default. A
reader who has installed and is coming back to understand what they got needs the
catalogue in the same place. After Install serves neither; before Install serves
both.

**F2 — table vs. prose entries?**
Prose. The brief itself states the test: "Five rows of a table cannot carry 'no
numeric aesthetic score, ever'." That is not a hypothetical — the design entry's
load-bearing position is a refusal, and a refusal shrunk to a table cell reads as
a feature. A table is the right shape for the Cuts section above, where every row
is the same kind of thing (section, lines, action, reason). The catalogue entries
are not the same kind of thing; they are five skills with different registers.

**F3 — how much ideology belongs in the README vs. a link?**
One position per add-on, stated as a sentence, with the rest behind the link. The
design entry names the refusal ("no score, no ranking, a 'good default aesthetic'
refused on purpose") because that refusal is the reason the skill exists and the
thing that makes it interesting. It does not explain the contrast floor's
system-blind property or the effect budget — those belong in the SKILL.md, and the
reader who wants them will click through. The test: if the position is the reason
a reader would install or skip the add-on, it belongs; if it is a detail of how
the add-on works, it does not.

**F4 — peers or satellites?**
Satellites, but the README should not say so. The dependency direction is
one-way: every add-on is loaded "alongside" the core skill and defers to
`COOPERATION.md` for phases, ownership, and quorum. The install default (all four,
unasked) is a packaging convenience, not a peer declaration — `--no-addons`
exists, and the core skill is useful without any of them. But writing "these are
satellites" in the README introduces a hierarchy the reader does not need. The
order (core first, add-ons after) and the one-line relationship note in each
add-on entry ("load it alongside the core skill, never instead of it") does the
work without the label.

**F5 — hook leads with protocol, failure mode, or artifact trail?**
Failure mode first, then artifact trail. The protocol is the mechanism; the
artifact trail is the evidence; the failure mode is the reason the reader is
reading. Leading with the protocol ("Parley Deck is a transport-agnostic
cooperation protocol…") is what the current README does and it is why the current
blockquote does not earn the second paragraph — it describes the product, not the
problem. Leading with the artifact trail ("every step is a file you can read")
is honest but premature: the reader has not yet been told what problem the files
solve. The failure modes are the thing the reader already recognises; the
artifact trail is the answer that makes them keep reading.

**F6 — touch `parley-deck-cli/README.md` in the same change?**
No. This idea's target file is `parley-deck-skill/README.md`; the brief says so
explicitly in its frontmatter. The CLI README is a different package with a
different audience (people running the `parley` CLI, not people installing the
skill). Bundling both into one change mixes two scopes and makes the diff harder
to review. If the CLI README has a parallel problem, that is a separate idea with
its own `00-prompt.md`.

**F7 — a claim in the current README that is no longer true, beyond the four listed?**
Yes. Line 9 of the current README lists the supported runtimes as "Codex, Claude
Code, Antigravity, Gemini, Hermes, or a custom skill directory" — six targets.
The installer `--help` (verified by running it) lists fourteen native targets:
codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe,
cursor, opencode, aionrs. The "Installation Details" section on line 186 does
list the full set, but the opening paragraph on line 9 is stale — it names six
where fourteen now exist. A reader who uses, say, Cursor or Goose reads that line
and concludes the package does not support their runtime.
