---
idea: readme-skill-catalogue
round: 01
agent: claude-1
date: 2026-07-29
lens: structure & voice
---

## Position in one line

The README's problem is not that it lacks a catalogue. It is that **it opens by
describing a package manager**, and the four add-ons are missing because nobody ever
decided what the file is *for*. Fix the second thing and the first fixes itself.

---

## 1. The hook — as it would ship

```markdown
# Parley Deck Skill

> **One model playing four reviewers is still one model.**

Parley Deck is a cooperation protocol built on that assumption. Several *separate*
agents each write their own file, each commits to a position before it is allowed to
read anyone else's, and the disagreements are kept on disk instead of dissolving into
a chat log.

What you get at the end is not a transcript. It is a directory you can diff, resume,
and hand to someone who was not there.

This package installs the protocol — and four add-ons that build on it — into 15 agent
runtimes.
```

Nine lines. Three claims, all checkable: separate agents (§1 non-solo, one file per
agent per round), write-before-read (§4 Phase 1), 15 runtimes (count the installer's
`--target` list: codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi,
droid, vibe, cursor, opencode, aionrs = 14 named runtimes + `generic` = 15 targets;
**I want this number checked by codex-1 before it ships** — "15 runtimes" and "15
targets" are not the same sentence and I may be conflating them).

**What the current opening does wrong.** Line 3 is good — *"not one model role-playing
a committee"* is the sharpest sentence in the file and I would keep the idea. Line 7
then throws it away: *"installs the vendor-neutral Parley Deck cooperation instructions
(and a fallback protocol snapshot) into supported agent runtimes … then helps you check
and sync project metadata."* That is a sentence about metadata sync, positioned where
the reader decides whether to keep reading. A parenthetical about a fallback snapshot
is the fourth thing a reader needs, not the second.

---

## 2. The catalogue — as it would ship

**Shape: prose entries in a fixed four-part frame** — one line of what · the failure it
prevents · one load-bearing rule stated as a claim · the practical note. The frame never
varies between entries, which is the same discipline `PDS.md` applies to its own artifact
definitions. A scannable index sits on top for the reader who is looking for a name.

```markdown
## What's in the box

Installing this package installs five skills. The first is the protocol; the other four
build on it and are opt-in.

| | |
|---|---|
| **`parley-deck`** | the cooperation protocol itself — the core skill |
| **`parley-design`** | building a design system that doesn't look generated |
| **`parley-design-check`** | the part of that a machine can actually check |
| **`parley-tracker`** | tickets a business reader, an engineer, and an agent can all read |
| **`parley-worktrees`** | parallel agents that don't corrupt each other |

### `parley-deck` — the protocol

The core skill. It teaches an agent an eight-phase lifecycle: kickoff, independent
analysis, cross-review, consensus, `FINAL.md`, `IMPLEMENTATION.md`, code review, fix-up.
Every phase is a file, every file has exactly one owner, and no agent may write into
another's.

**The failure it prevents:** the multi-agent run that was really one agent narrating a
committee — and the review that vanished into a scroll-back nobody can resume.

**Its load-bearing rule — write before you read.** A round-1 file drafted after reading
the others is not an independent position, and the protocol treats it as a failed round.

Runs over `local-dir`, `github-pr`, or `gitlab-mr`. The files stay canonical; the pull
request is a mirror of them. Rigor is a dial, not a religion: `track: fast | standard |
deliberation`, because a typo does not need a quorum.

### `parley-design` — a design system that doesn't look generated

Doctrine and a protocol (`PDS/1.0`): typed artifacts from `DESIGN-BRIEF` to `AUDIT`,
four gates, four conformance levels, four evidence tiers, and a rule registry where the
machine-readable YAML and the human-readable prose live in one file — so there is no
second copy to drift out of step.

**The failure it prevents:** the interface you could have guessed from the category
alone. Purple gradient, glass card, three feature columns, a testimonial from nobody.

**Its load-bearing rule — selection, never averaging.** Exactly one direction wins whole.
You may graft at most three things from the ones that lost, and a graft may not touch the
winner's token file. Average four directions together and you get the average of four,
which is the look everything already has.

It ships no house style and emits no score. *"Give me a good default aesthetic"* is a
request this skill refuses on purpose. Pure markdown: no dependencies, four files, a
64 KiB budget that a test enforces rather than a comment.

### `parley-design-check` — the part a machine can check

Runs the checkable rules over files on disk and prints findings in one shape:
`rule-id — violation — remedy`. Node built-ins only. No dependencies, no network — not
at check time, not ever.

**The failure it prevents:** the clean report that means the tool did not understand the
file.

**Its load-bearing rule — it would rather refuse than pass.** Registry missing? Exit 3
and say so, instead of falling back to a bundled copy that has drifted. A rule it has no
detector for is reported `UNJUDGEABLE` by name, never silently skipped. Its capability
list is generated by scanning its own detectors, so it cannot advertise coverage it does
not have.

Most of what makes work read as machine-made is not decidable from source. This tool
says so instead of passing it.

### `parley-tracker` — one ticket, three readers

Epics, stories, and subtasks written so that a business reader, an engineer, and the
agent that implements them can each read the same file without a translation layer:
`[B]`, `[T]` and `[A]` sections, acceptance criteria tagged by audience, and a Definition
of Done that is one checkbox per criterion with the command that verifies it.

**The failure it prevents:** the ticket an agent implements confidently and wrongly,
because the thing everyone in the room already knew was never written down. A gap-scan
runs before work starts and refuses the ticket that leaves it unwritten.

**Its load-bearing rule — the tracker is a mirror.** The markdown file is canonical, sync
is one-way by default, and a field your tracker cannot represent is dropped from the
mirror and never from the file. Migrate from Jira to Linear and you lose a projection,
not a requirement.

### `parley-worktrees` — parallel agents that don't corrupt each other

Allocation, naming, isolation, merge, and cleanup for git worktrees, so several sessions
or implementers can work over one repository at the same time.

**The failure it prevents:** the one that leaves no stack trace — two agents writing the
same file, and a result that was nobody's intention.

**Its load-bearing rule — declare your file set, and it must be disjoint.** An
intersection is refused, or overridden on the record. That single discipline converts
invisible concurrent-write corruption into an ordinary merge conflict, which is a problem
git already knows how to show you.
```

**Why this and not a table (F2).** A five-row table can carry "selection, never
averaging" as three words but cannot carry *why it is true*, and the why is the entire
reason a reader would want the skill. The index table above is deliberately thin — it
exists to answer "which one is that?", not to describe anything. Everything load-bearing
is prose.

**Why one position per add-on and not three (F3).** Each add-on has five or six ideas
worth stating. Stating all of them turns the README into the skill and gives the reader
nothing to click. One is a hook; three is a summary; a summary is what people skip.

---

## 3. Cuts — the budget

The catalogue above is ~95 lines. Adding it to 402 gives ~497, which is worse, not
better. Target for the finished file: **≤330 lines**, so ~165 lines have to go.

| Section | Now | Action | Saved |
|---|---|---|---|
| `Use Parley Deck` — 7 near-identical prompt blocks | ~55 | keep 3, delete 4 | ~30 |
| `Installation Details` + `Updating` | ~120 | merge into one `Install` section; `--force` is one line, not a chapter; the eight-path manual list becomes `parley-deck-skill paths` | ~70 |
| `What The Skill Does` | ~14 | delete — it restates "What the protocol gives your agents" in weaker words | ~14 |
| `Why This Exists` (5 bullets) | ~12 | delete — the catalogue now names a concrete failure per skill, which is the same argument made better | ~12 |
| `Repository Layout` | ~28 | rewrite: it currently shows **no `addons/` at all**, and also omits `test/`, `NOTICE.md`, `RELEASING.md` | ~10 |
| `Status` | ~6 | fold its one real sentence into the hook | ~5 |
| `Local Agent Contract` | ~12 | keep — it is the only place the vendor-neutral contract is stated | 0 |

Two of these are deletions of *good* content. `Why This Exists` is well written. It goes
because after the catalogue it is the third time the reader is told the same thing, and
the third telling is where people close the tab.

---

## 4. Placement (F1)

**Hook → catalogue → install → use → everything else.**

The counter-argument is real: a reader who has installed nothing does not need five
descriptions. But the alternative loses worse. Install-first assumes the reader already
decided, and this README has to *cause* that decision. A three-line `npx` block is not
what anyone bounces off; a page that never says what the thing is, is.

There is also a fact that settles it. The four add-ons **install by default, unasked**.
A reader who follows the install block before reading the catalogue has already put four
skills on their machine that the document has not mentioned. Install-first would make
the README complicit in that.

---

## 5. My anti-slop test, and its result on my own draft

**The test — strike every adjective and see whether a fact survives.** Machine-made copy
collapses under it, because the adjectives were carrying the sentence.

Applied to my hook: *"One model playing four reviewers is still one model."* — no
adjectives present. Survives trivially.

Applied to the entries: `parley-design`'s failure line is *"Purple gradient, glass card,
three feature columns, a testimonial from nobody"* — strike "purple", "glass", "three":
four nouns remain and the sentence still names the thing. Survives.

**Second test, and this one caught me.** No two catalogue entries may open with the same
grammatical construction. My first draft opened `parley-deck` with *"The core skill."*
and `parley-design` with *"Doctrine and a protocol."* — two bare noun phrases in a row,
which is exactly the identical-shape rhythm that reads as generated. I did not notice it
by reading; I noticed it by running the check. Both survived in the draft above only
because I then varied the entries that followed (verb-led, noun-led, gerund-led, article-
led), and I am flagging the pair I did **not** fix so a reviewer can rule on whether two
adjacent noun phrases are acceptable when the entries after them break the pattern.

**Third test — the honesty rule applied to ourselves.** I searched my own draft for a
number I had not counted. I found one: "15 agent runtimes" (see §1). It is in the draft
and it is unverified. I am not going to quietly drop it to look clean; I am marking it
for codex-1.

---

## Forks

- **F1 — catalogue leads.** Argued above. The install-by-default fact decides it.
- **F2 — prose, with a thin index table.** A table cannot carry the reason.
- **F3 — exactly one position per add-on**, stated as a claim, plus the link.
- **F4 — peers of each other, satellites of the protocol.** One heading, "What's in the
  box", core first and explicitly labelled the core. Do not rank the four; they are for
  different jobs, and ranking them invents a hierarchy the package does not have.
- **F5 — failure mode first.** The artifact trail is the better *proof* but the worse
  *opener*: it answers a question the reader has not asked yet. Name the failure, then
  offer the trail as evidence.
- **F6 — separate idea for `parley-deck-cli/README.md`.** The user asked for this
  package. The CLI README has a different reader (someone choosing a binary, not a skill)
  and copying this catalogue there unchanged would be wrong. Record it as a follow-up.
- **F7 — false or stale claims beyond the four already listed:**
  1. Line ~242: *"Until the WinGet manifest is accepted, download the `.exe`"* — a WinGet
     manifest **was** accepted; `Feci.ParleyDeckSkill` is published. This line is now
     false and should state the `winget install` command.
  2. Line ~9: the intro names six runtimes (*"Codex, Claude Code, Antigravity, Gemini,
     Hermes, or a custom skill directory"*) where the installer supports fourteen plus
     `generic`. Not false, but it undersells by more than half.
  3. Line ~186 lists the native targets correctly — so the README **contradicts itself**
     between line 9 and line 186. One of them has to go.
  4. `Repository Layout` omits `addons/` entirely, which is the specific stale claim that
     makes this idea necessary.

---

## One thing outside my lens that I think matters more than my lens

The README is not the only place the add-ons are invisible. I ran the universal skill
installer from `vercel-labs/skills` against our published repository:

```
$ npx -y skills@latest add feci/parley-deck-skill --list
◇  Found 1 skill
     parley-deck
```

**One of five.** That CLI's documented rule is that a `SKILL.md` at the repository root
shadows anything nested beneath it, and ours is at the root. Anyone installing us through
that route today silently gets the core skill and none of the add-ons.

This is a packaging defect, not a README defect, and it belongs to a separate idea rather
than this one. I raise it here only so that nobody writes a README sentence promising
five skills through an install path that delivers one.
