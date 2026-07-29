---
idea: readme-skill-catalogue
round: 02
agent: codex-1
date: 2026-07-29
---

## Position changes and convergence

I would ship a shorter version than my round-01 draft: one memorable position per
add-on, with operational details left behind the linked `SKILL.md`. I also withdraw
the phrase "every step is a file" as a literal description. The durable state is in
files, but some phases own several files and `consensus.md` has several append-only
signers.

F1, F2, F3, F5, and F6 converged at the decision level, and I agree with that
convergence:

- **F1:** catalogue before Install.
- **F2:** prose entries, not a catalogue table. A thin scan list is optional, not the
  descriptive surface.
- **F3:** one load-bearing position per add-on; details stay in the linked skill.
- **F5:** failure mode first, artifact trail as proof.
- **F6:** do not touch `parley-deck-cli/README.md`; that is a separate idea.

I would not reopen those forks. F4 resolves cleanly as Kimi framed it: the add-ons are
satellites in dependency terms and peers in visual weight. Put the core first, state
that add-ons load alongside it, and give all five the same heading level.

## Cross-review

I am `codex-1`, so the "all except yourself" instruction means the other three active
participants below. Corrections to my own round-01 claims are included in the truth
audit.

### @claude-1

What Claude got right that I did not:

- The hook at line 22, **"One model playing four reviewers is still one model,"** is
  sharper and less procedural than my "can role-play a committee." It is the right
  base.
- The tracker consequence at lines 139–142 — **"Migrate from Jira to Linear and you
  lose a projection, not a requirement"** — makes canonical ownership concrete. My
  entry stated the contract but did not make the benefit memorable.
- "Peers of each other, satellites of the protocol" (lines 240–242) is a useful F4
  distinction, even though Kimi found the cleaner reader-facing formulation.

What is wrong or unsupported:

- **"into 15 agent runtimes"** (lines 32–33) is wrong. `lib/installer.js:13-113`
  defines fourteen native runtime targets. `generic` is a target that requires an
  explicit destination (`lib/installer.js:246-262`), not a fifteenth runtime.
- **"Every phase is a file, every file has exactly one owner"** (lines 78–79) is too
  absolute. A review phase has one file per reviewer, and every participant appends a
  signoff to `consensus.md` (`references/COOPERATION.md:329-345`). Ownership is
  per participant artifact or signoff block, not per whole file in every case.
- **"7 near-identical prompt blocks"** (line 177) repeats the brief's bad count. The
  current README has eight fenced `text` blocks beginning at lines 69, 74, 79, 84,
  91, 99, 107, and 114.
- **"a WinGet manifest was accepted; `Feci.ParleyDeckSkill` is published"**
  (lines 250–252) is not established by a shipped file. The shipped
  `packaging/winget/README.md:3-4` still calls the manifest a draft and lines 21–34
  describe copying it into `winget-pkgs` and opening a PR. The current README may
  need an external release check later, but this round cannot call it false under the
  shipped-file evidence rule.
- The `skills@latest` result at lines 265–276 may expose a real packaging defect, but
  it is evidence from an external tool run, not a shipped file. Keep it for a
  follow-up idea; do not let it alter the copy in this one.

### @hermes-1

What Hermes got right that I did not:

- Lines 10–15 name several recognizable failures — anchoring, disappearing
  disagreement, premature implementation, and unowned review — instead of relying on
  one committee metaphor. That is stronger reader-side reasoning than my first hook.
- Lines 82–83 preserve a useful boundary: **"The skill teaches an agent how to
  participate; the companion `parley` CLI orchestrates the runs."** My draft omitted
  that distinction.
- The design entry at lines 87–96 states the dependency honestly: load the add-on
  alongside the core, never instead of it.

What is wrong or unsupported:

- **"Use Parley Deck (seven prompt blocks)"** (line 148) is the same incorrect count;
  there are eight in `README.md:69-117`.
- **"then mirroring them into Jira, Linear, GitHub Issues, GitLab, Trello, or a plain
  board"** (lines 113–115) overstates what this skill alone performs. It authors the
  canonical file and neutral payload; actual live create/update requires a separate
  opt-in connector (`addons/parley-tracker/SKILL.md:381-388`).
- **"refusing `in-progress` status on any incomplete field"** (lines 120–122) is too
  broad. The shipped gap-scan checks a specified required schema and required
  sections (`addons/parley-tracker/SKILL.md:212-250`), not every possible field.
  **"the single highest-leverage rule for AI output quality"** is present in the
  shipped skill at line 251, so it is attributable, but it remains an unmeasured
  superlative and does not belong in this README.
- The line-9 accusation says the intro names "six targets" (lines 283–290). It names
  five runtimes plus a custom directory. The real defect is incompleteness and
  ambiguous exhaustiveness, not that any named target is unsupported.

### @kimi-1

What Kimi got right that I did not:

- The reader split at lines 11–25 distinguishes discovery, immediate installation,
  and return-reference use. That is the best explanation of why the catalogue must be
  high and the full command reference must be low.
- **"files in your repository you can read, diff, and resume — not a chat log you
  have to trust"** (lines 43–44) is the strongest artifact-trail sentence in the
  round. I use it with the literal "every step" softened to "working state."
- **"Satellites in framing, peers in visual weight"** (lines 246–252) is the cleanest
  resolution of F4.

What is wrong or unsupported:

- **"consensus rates confidence by agreement"** (lines 87–89) repeats a false claim
  from the current README. The protocol asks the drafter to retain contradictions,
  partial coverage, unique insights, and blind spots
  (`references/COOPERATION.md:324-327`); that section is advisory, while signoffs are
  the gate (`references/COOPERATION.md:339-345`). It defines no confidence rating.
- The cut arithmetic does not reach the stated target. Lines 157–159 claim roughly 88
  lines removed against 75 catalogue lines and 8 corrected-layout lines. Starting
  from the actual 401 lines, that yields 396 if both additions are counted, or 388 if
  the layout number was intended as a saving — neither is at most 300.
- **"seven prompt blocks"** (line 151) is again the wrong count, and the proposed
  keep/drop list accounts for only seven of the eight blocks.
- **"then mirrors them into Jira..."** (lines 121–124) needs the same connector
  qualification as Hermes's entry. The canonical authoring and projection contract
  ship here; live tracker writes do not (`addons/parley-tracker/SKILL.md:381-388`).

## Truth audit of the current README and the round-01 accusations

The checked target is the current shipped `parley-deck-skill/README.md`, which ends
at line 401.

| Claim or accusation | Evidence checked | Verdict |
|---|---|---|
| The README has 402 lines. | The current file ends at `README.md:401`. | **Wrong count.** It has 401 lines. |
| "Use Parley Deck" contains seven prompt blocks. | Fenced `text` blocks begin at `README.md:69,74,79,84,91,99,107,114`. | **Wrong accusation.** There are eight. |
| The repository layout is stale because it omits `addons/`. | The displayed tree at `README.md:148-166` has no `addons/`; `package.json:31-42` ships that directory, and it contains four add-on `SKILL.md` files. | **Genuinely stale/incomplete, but false by omission rather than a false listed path.** Rewrite it. |
| Install/update material repeats `install --target all --force`. | The substring occurs at `README.md:239,307,332,346`. | **True redundancy, not a truth claim.** Consolidate it. |
| The Windows example is current. | `README.md:239` names `v1.2.1`; `package.json:3` is `1.5.0`. | **Genuinely false.** Update or make the example versionless. |
| The opening runtime list is false. | `README.md:9` names five supported runtimes plus a custom directory; all are valid. The fourteen native target definitions are at `lib/installer.js:13-113`, and the complete CLI target list appears at `README.md:289`. | **Merely imprecise and stale.** It reads as exhaustive but is only partial. The claims that it names "six runtimes" or directly contradicts line 186 are too strong. |
| The eight-phase lifecycle is "append-only." | The claim is at `README.md:21-23`. Only signoff blocks are append-only (`references/COOPERATION.md:329-345`); `IMPLEMENTATION.md` is explicitly living and updated (`references/COOPERATION.md:424-446`). | **Genuinely false as written.** Replace "append-only" with "file-canonical" or limit it to signoffs. |
| The comparison lens "rates confidence by agreement." | The claim is at `README.md:26-27`; the actual fields are contradictions, partial coverage, unique insights, and blind spots at `references/COOPERATION.md:324-327`, with advisory status at line 345. | **Genuinely false.** No shipped confidence rating exists. |
| "Any capable tier-1 model can follow it." | The claim is at `README.md:119`; neither `SKILL.md` nor the protocol defines "tier-1" or a compatibility test. | **Unsupported, not demonstrated false.** Remove the universal performance claim. |
| The default uses all discovered installed CLI agents. | The claim is at `README.md:371`; `SKILL.md:280-289` specifies a bounded set, normally 2–4, including a non-facilitator when available. | **Genuinely false.** Replace it with the bounded-set rule. |
| "The protocol's value should be obvious." | `README.md:397`. | **Neither true nor false; it is an opinion standing where evidence should be.** End with the inspectable artifacts, not the conclusion. |
| The WinGet-pending sentence is already false. | The sentence is at `README.md:242`; the shipped packaging guide still calls its manifest a draft (`packaging/winget/README.md:3-4`) and instructs the publisher to open a PR (`:21-34`). | **Wrong accusation under this round's evidence rule.** Publication is not verified by a shipped file; do not replace the sentence with a `winget install` promise on this evidence. |

One correction to the brief itself matters for the merged catalogue:
`addons/parley-design-check/SKILL.md:53-61` defines exit code `4` for an
`UNJUDGEABLE` run. The brief's exit-code list stops at `3`; the shipped skill is
authoritative.

## Base selection

This is selection and grafting, not averaging:

- **Hook base: claude-1.** It has the strongest first sentence and gets to
  independent participants immediately.
- **`parley-deck` base: codex-1.** It is the most precise on artifact ownership,
  tracks, and transports, and it does not repeat the false confidence-rating claim.
- **`parley-design` base: kimi-1.** It compresses the doctrine into memorable
  refusals without reproducing the whole PDS protocol.
- **`parley-design-check` base: codex-1.** It alone includes the shipped exit-4
  behavior that the brief omitted.
- **`parley-tracker` base: codex-1.** It alone draws the shipped boundary between
  canonical authoring/projection and a separate live connector.
- **`parley-worktrees` base: kimi-1.** It is the shortest complete account of the
  `IMPLEMENTATION.md` lock manifest and disjointness check.

## Single merged shipping copy

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

## Five skills, one package

The package's installer puts all four add-ons alongside the core skill by default.
Use `--no-addons` for the core alone, or `--only <name>[,<name>]` for the core plus
a selected set.

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

Every product claim in that copy is grounded in these shipped ranges:
`SKILL.md:280-289,463-467`;
`references/COOPERATION.md:267-345,397-470`;
`addons/parley-design/SKILL.md:8-16,34-42`;
`addons/parley-design/references/PDS.md:381-404`;
`addons/parley-design-check/SKILL.md:8-30,53-82`;
`addons/parley-tracker/SKILL.md:212-250,376-388`;
`addons/parley-worktrees/SKILL.md:338-405`; and
`lib/installer.js:782-835`.

## Final cut list

| Current material | Checked range | Shipping action |
|---|---:|---|
| Current pitch, "What the protocol gives your agents," and full provenance block | `README.md:1-42` | Replace the pitch and protocol bullets with the merged hook and catalogue. Move provenance to a compact note beside repository relationships. |
| Install plus Installation Details, Installer Commands, and Updating | `README.md:44-63,178-359` | Merge into one `Install, update, and remove` reference. Keep the fastest path, verification, add-on selectors, target list, project/custom installs, genuinely different platform paths, one update command, dry-run, doctor, and uninstall. |
| Eight usage prompts | `README.md:65-117` | Keep three: start/design, implement, and continue. Express review and transport as short substitutions, not five more fences. |
| Universal "tier-1" claim | `README.md:119` | Delete. |
| Why This Exists | `README.md:121-131` | Delete. Its failure modes are now carried by the hook, the core entry, and the ownership language. |
| What The Skill Does | `README.md:133-146` | Delete. The core catalogue entry replaces it. |
| Repository Layout | `README.md:148-176` | Rewrite as a compact accurate tree including `addons/`; remove prose that repeats the catalogue. |
| Windows example | `README.md:236-242` | Remove the stale `v1.2.1` or update it from a verified release artifact. Do not claim WinGet publication until a shipped source records it. |
| Local Agent Contract | `README.md:361-371` | Keep, but replace "all discovered" with the bounded-set default from `SKILL.md:280-289`. |
| Transports | `README.md:373-381` | Keep with only copy edits. |
| Repository relationships and provenance | `README.md:31-42,383-391` | Merge into one short section; preserve attribution and the skill/CLI boundary. |
| Status | `README.md:393-397` | Delete. It is dated positioning plus an unsupported conclusion. |
| License | `README.md:399-401` | Keep. |

## Final line budget and ordered sections

The final README budget is **at most 280 lines**. The section allocations below
include their headings, fences, and internal blank lines.

| Order | Section | Maximum lines |
|---:|---|---:|
| 1 | Title and hook | 16 |
| 2 | Five skills, one package | 80 |
| 3 | Install, update, and remove | 72 |
| 4 | Use Parley Deck | 24 |
| 5 | Local Agent Contract | 12 |
| 6 | Transports | 9 |
| 7 | Repository Layout | 16 |
| 8 | Relationship and provenance | 16 |
| 9 | License | 3 |
|  | Editorial reserve | 32 |
|  | **Total ceiling** | **280** |
