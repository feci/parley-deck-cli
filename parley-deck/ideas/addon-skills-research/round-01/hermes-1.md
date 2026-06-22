---
agent: hermes-1
idea: addon-skills-research
round: 1
date: 2026-06-22
---

## Summary

My lens is the three-audience ticketing problem: a ticket must read well for a
**business person** (value/why), a **technical person** (how/constraints/NFRs),
and an **AI agent** (self-contained, no-assumption, machine-parseable). Those
three readers want different things from the same artifact, and the failure mode
of most templates is that they optimise for one and degrade for the others —
business prose hides the interfaces an engineer needs; engineer shorthand hides
the value a stakeholder needs; and both routinely leave gaps an AI agent will
silently fill with assumptions, producing plausible-but-wrong work.

The thesis of this round, from my lens:

1. **One file, three tagged sections, one machine block.** Each ticket is a
   single markdown file with an audience-tagged body (`[B]`, `[T]`, `[A]`) and a
   YAML frontmatter metadata block. The file is canonical; the tracker (Jira,
   Linear, GitHub Issues, Trello, …) is a **mirror**. This keeps the skill
   vendor-agnostic and keeps AI agents reading the source of truth rather than a
   tracker's lossy projection of it.

2. **Gaps are forbidden, not implicit.** The template enforces "no blank fields" —
   every required slot must be filled or explicitly marked `n/a (reason)`. The
   `ticket-tracker` skill ships a **gap-scan checklist** the agent runs *before*
   claiming a subtask: any required slot missing and not marked `n/a` blocks the
   claim (the agent asks or returns `BLOCKED`, never assumes). This is the single
   highest-leverage rule for AI output quality and it is what most existing
   templates get wrong by omission.

3. **Acceptance criteria are hybrid and audience-split.** Behavioural criteria
   use Given/When/Then (Gherkin); non-functional criteria (perf, security,
   availability) use **measurable bullets** with an explicit number and a
   verification command. Gherkin-for-everything degrades NFRs; this split keeps
   each format where it is honest. Each criterion carries an audience tag so a
   reader knows which lens it serves.

4. **The worktree seam is itself a ticket.** From the ticketing lens, a
   `worktree-sessions` session is a subtask whose metadata records the worktree
   path, branch, base commit, and owner. Claiming a task = setting
   `status: in-progress` + `assignee` + the worktree hint; finishing = merged +
   pruned. The shared task doc the worktree track needs *is* the set of subtask
   tickets. That is the concrete integration seam between the two skills.

## Proposed approach

### 1. Canonical file model — "tracker is a mirror, files stay canonical"

Every ticket lives as a markdown file under `tickets/` in the repo (path mirrors
epic hierarchy: `tickets/<epic>/<story>/<subtask>.md`). The YAML frontmatter is
the machine contract; the body is the human contract. A `sync` command in
`ticket-tracker` projects the file into whichever tracker is configured, mapping
fields onto the tracker's native schema with a **graceful-degradation** rule:
fields the tracker lacks are dropped from the mirror but never from the file.
Re-sync is **one-way (file → tracker)** by default; a `--pull` mode reconciles
tracker-side edits by writing them back to the file only when flagged
`mirror-owned: true` on the field, to prevent silent drift.

Minimum-viable fields every tracker supports (the mirror contract):
`id`, `type` (epic|story|subtask), `title`, `status`, `assignee`, `parent`,
`labels`, `priority`. Extended fields (`files`, `apis`, `nfr`, `dod`, `worktree`)
live in the file and are mirrored only where the tracker allows.

### 2. The three-audience template skeleton (shared by all three ticket types)

```
---
id: E-001 | S-001 | T-001        # tracker-stable id, never reused
type: epic | story | subtask
title: <imperative, <=80 chars>
parent: <parent id or n/a>        # epic has none; story→epic; subtask→story
status: draft | ready | in-progress | review | done | blocked | dropped
assignee: <human | agent:hermes-1 | n/a>
priority: p0 | p1 | p2 | p3
labels: [domain:auth, nfr:security, ...]
estimate: <points or n/a>
files: [src/auth/*.ts]            # path hints; AI uses these to scope reads
apis: [OpenAPI: ./openapi.yaml#/auth, n/a]
arch: [./docs/arch/auth.md, n/a]
worktree: n/a | { path, branch, base }
dod: [ ] criterion refs           # checkbox list of AC ids
mirror-owned: [status, assignee]  # fields the tracker may write back
---

# <title>

## At a glance            # 2-4 lines, all three audiences read this first
<one-sentence outcome · one-sentence scope · one-sentence done-when>

## [B] Business            # value, who benefits, why now, success indicator
## [T] Technical           # how, interfaces touched, constraints, NFRs
## [A] Agent directives    # self-contained scope, explicit do/don't, gaps policy

## Acceptance criteria     # tagged AC-1..AC-N; Gherkin or measurable bullet
## Non-goals               # explicit out-of-scope; prevents AI scope creep
## Dependencies / blocks   # ticket ids and reasons
## Open questions          # each tagged [needs:business|technical|decision]
```

The audience tags are not separate documents — they are sections of one file.
A reader skims their tag; the AI reads all three plus the frontmatter. The
**At a glance** block is the only mandatory cross-audience block and is what
makes a stakeholder willing to open the file at all.

### 3. Epic template (concrete instance)

```markdown
---
id: E-014
type: epic
title: "Add worktree-isolated parallel implementers to Parley Deck"
parent: n/a
status: draft
assignee: n/a
priority: p1
labels: [domain:parley-deck, skill:worktree-sessions]
files: [parley-deck/protocol/**, parley-deck/ideas/**]
apis: [n/a]
arch: [./docs/arch/parley-deck.md]
worktree: n/a
dod: [AC-E1, AC-E2, AC-E3]
mirror-owned: [status]
---

# Add worktree-isolated parallel implementers to Parley Deck

## At a glance
Let multiple AI implementers build in parallel against one repo without
collisions · scope = git worktree lifecycle + a claim protocol · done when two
implementers can run concurrent Phase-5 sessions that merge cleanly.

## [B] Business
Outcome: quorum rounds finish faster because non-conflicting ideas build in
parallel instead of serially. Success indicator: a 4-idea round completes build
+ merge in <=1.5x the time of a single-idea build (measured on the
benchmark-round fixture). Why now: native worktree support in the runtimes
(Claude Code/Cursor/Copilot CLI, 2025-26) makes this cheap; we currently serialise.

## [T] Technical
Mechanics live in the `worktree-sessions` skill (git worktree add/list/remove/
prune, branch-per-session, per-worktree `.env.local`). The core protocol gets a
thin seam only: §6 conflict-avoidance note + Phase-5 "implementers MAY be
isolated in worktrees". Integration branch → test → main. NFR: a worktree
creation must not touch the user's checked-out working tree.

## [A] Agent directives
Do not implement git mechanics in the core protocol — that belongs in the skill.
Do not assume a specific tracker; the claim protocol must be file-based. If a
worktree already exists for a branch, reuse it rather than erroring (record in
`worktree` metadata). Gaps → ask, do not invent commands.

## Acceptance criteria
- AC-E1 [A][T] Given two implementers assigned non-overlapping file sets, When
  each runs in its own worktree, Then both merge into the integration branch
  with zero manual conflict resolution. (Gherkin)
- AC-E2 [B][T] Measurable: total wall-clock for the 4-idea benchmark fixture is
  <=1.5x single-idea build time. Verify: `make bench-round fixture=4-idea`.
- AC-E3 [A][T] Given a stale worktree whose branch was deleted, When
  `worktree-sessions prune` runs, Then the worktree dir and metadata are removed
  and `git worktree list` shows no stale entry. (Gherkin)

## Non-goals
- No new git primitives; use stock `git worktree`.
- No IDE-specific integration (Claude Code/Cursor already do their own).

## Dependencies / blocks
- blocks E-015 (ticket-tracker claim protocol) on the `worktree` metadata field.

## Open questions
- [needs:decision] Is the integration branch per-epic or repo-global?
```

### 4. User story template (concrete instance)

```markdown
---
id: S-042
type: story
title: "As a Parley author, I claim a subtask so no other implementer duplicates it"
parent: E-015
status: ready
assignee: agent:hermes-1
priority: p1
labels: [domain:ticket-tracker, coord:claim]
files: [parley-deck/ideas/**/round-*/]
apis: [n/a]
arch: [n/a]
worktree: { path: ../wt-s042, branch: wt/s042, base: f98507c }
dod: [AC-1, AC-2, AC-3]
mirror-owned: [status, assignee]
---

# As a Parley author, I claim a subtask so no other implementer duplicates it

## At a glance
Prevent two implementers building the same subtask · scope = a claim operation on
the subtask file · done when a second claim on an in-progress subtask is refused.

## [B] Business
Outcome: no wasted parallel work; round cost tracks unique subtasks, not
attempted ones. Success indicator: zero duplicate-build incidents across the
next 3 rounds.

## [T] Technical
Claim = set `status: in-progress` + `assignee` + `worktree` in the subtask's
frontmatter, committed atomically. Conflict detected by `git` on the shared task
index; loser retries against an unclaimed subtask. NFR: claim must be a single
commit, no partial state.

## [A] Agent directives
Read the subtask file fully before claiming. Run gap-scan; if any required slot
is missing and not `n/a`, do NOT claim — return `BLOCKED: <slot>`. Never edit a
subtask you have not claimed. After claim, record `worktree` so others can see
isolation. Re-claim is forbidden; if `status` is already `in-progress` and
`assignee != you`, skip it.

## Acceptance criteria
- AC-1 [A][T] Given an unclaimed subtask, When agent X commits a claim, Then the
  file shows `status: in-progress`, `assignee: agent:X`, `worktree` filled.
- AC-2 [A][T] Given a subtask with `status: in-progress` and `assignee: agent:Y`,
  When agent X attempts to claim, Then the operation is refused with a message
  naming Y and the worktree path.
- AC-3 [T] Measurable: a claim is exactly one commit touching only the claimed
  file's frontmatter. Verify: `git show --stat <claim>` lists one file, one hunk.

## Non-goals
- No locking server; rely on git's last-writer-detects-conflict semantics.

## Dependencies / blocks
- blocked-by T-118 (frontmatter schema for subtask) — must ship first.

## Open questions
- [needs:technical] On a claim conflict, auto-retry how many times before surfacing?
```

### 5. Technical subtask template (concrete instance)

```markdown
---
id: T-118
type: subtask
title: "Define YAML frontmatter schema for ticket files"
parent: S-042
status: ready
assignee: n/a
priority: p1
labels: [domain:ticket-tracker, schema]
files: [parley-deck/ideas/addon-skills-research/**, skills/ticket-tracker/templates/**]
apis: [n/a]
arch: [n/a]
worktree: n/a
dod: [AC-1, AC-2, AC-3, AC-4]
mirror-owned: [status]
---

# Define YAML frontmatter schema for ticket files

## At a glance
A stable, typed frontmatter contract for all ticket files · scope = schema +
validator · done when a bad file fails validation with a field-level message.

## [B] Business
Outcome: tickets are machine-reliable, so tracker sync and AI parsing never
silently misread. Success indicator: 100% of round-02 tickets pass validation.

## [T] Technical
Define required vs optional fields, enums (`type`, `status`, `priority`),
`worktree` as an object or `n/a`, `files`/`apis`/`arch` as path arrays or `n/a`.
Ship a validator (`ticket-tracker validate <file>`) with field-level errors.
NFR: validator runs in <200ms per file; no network.

## [A] Agent directives
Every optional field omitted must be intentional — prefer explicit `n/a` over
absence for `apis`, `arch`, `worktree`; absence is allowed only for truly
non-applicable optional fields listed in the schema doc. Do not add fields not in
the schema. IDs are immutable; `parent` must resolve to an existing file or
`n/a`. Run `ticket-tracker validate` before marking `done`.

## Acceptance criteria
- AC-1 [A][T] Given a file with `type: bogus`, When validated, Then it fails with
  a message naming `type` and the allowed enum. (Gherkin)
- AC-2 [A][T] Given a subtask whose `parent` id does not exist, When validated,
  Then it fails with a message naming the missing parent id.
- AC-3 [T] Measurable: validation of a 50-file tree completes in <10s. Verify:
  `time ticket-tracker validate tickets/**`.
- AC-4 [B][T] Measurable: 100% of round-02 tickets pass `validate`. Verify:
  `ticket-tracker validate tickets/ --strict` exits 0.

## Non-goals
- No JSON Schema draft plumbing; a hand-rolled validator is fine.
- No migration of pre-existing tickets this round.

## Dependencies / blocks
- blocks S-042 (claim operation depends on the schema).
```

### 6. Acceptance-criteria format — the rules

- **Behavioural AC → Gherkin** (`Given/When/Then`), one scenario per AC, named
  `AC-N`. Keep the Given free of implementation detail a business reader can't
  parse; push impl detail into the [T] section and reference it.
- **Non-functional AC → measurable bullet**: `<measurable statement>. Verify:
  <single command>`. The `Verify:` clause is mandatory and is what makes an AC
  AI-verifiable and reviewer-checkable with one paste. If you cannot write a
  `Verify:` command, the AC is not measurable — rewrite it.
- **Every AC carries an audience tag** (`[B]`, `[T]`, `[A]`, or a combination) so
  a reader can filter to their lens and a reviewer can spot a story with only
  `[T]` ACs (a sign the business value is unstated).
- **Required AC categories** for any story/subtask: at least one happy-path
  behavioural AC, **and** an explicit edge/error/offline AC (marked `n/a` only
  with a reason). NFR ACs (perf/security/availability) are required for stories
  touching those concerns and forbidden from being implied.

### 7. Definition of Done — reusable, per ticket type

DoD is a checklist of **AC ids**, not prose, so it is verifiable:

- Epic DoD: all child stories `done`; all epic ACs pass their `Verify:`; non-goals
  respected (no scope creep commits in the diff).
- Story DoD: all ACs pass `Verify:`; `ticket-tracker validate` exits 0; claim
  released (`status: review`); linked subtasks `done`.
- Subtask DoD: all ACs pass `Verify:`; `validate` exits 0; tests green on the
  worktree's integration branch; `dod` checkboxes all ticked with commit refs.

### 8. Gap-scan checklist (the no-assumption enforcement)

Before an agent claims any ticket it must run, in order, and stop at the first
failure:

1. Frontmatter present and `validate` exits 0.
2. All required fields non-empty (`title`, `type`, `status`, `parent` resolves,
   `dod` references existing AC ids).
3. Every AC has an audience tag and either a Gherkin scenario or a `Verify:`
   command.
4. At least one edge/error/offline AC exists or is marked `n/a (reason)`.
5. `[A] Agent directives` section is non-empty and contains at least one
   `Do not` constraint (forces explicit negative scope).
6. `files`/`apis`/`arch` are either populated or explicitly `n/a`.

Any failure → the agent returns `BLOCKED: <reason>` referencing the failing
check, or opens an `Open questions` entry tagged `[needs:...]`. It never proceeds
by guessing.

### 9. Idea / FINAL.md → tickets decomposition rule

An idea's `00-prompt.md` is the epic seed; `FINAL.md` is the epic's design that
unfolds into stories + subtasks. Mapping:

- The idea's "Problem / idea" + non-goals → **epic** `[B]` + `Non-goals`.
- Each research/design decision in `FINAL.md` that produces user-visible value →
  a **story** (one story per distinct value stream).
- Each implementation unit within a story (schema, command, validator, template,
  seam edit) → a **subtask**, with `files`/`apis`/`arch` filled from the design.
- The `FINAL.md` sections map to `[T]` content of the relevant tickets.
- Open questions in `FINAL.md` become `Open questions [needs:...]` entries,
  promoted to tickets only when a decision is needed to unblock work.

This keeps `FINAL.md` canonical for *design* and tickets canonical for *work
state* — two different canonical artifacts, no overlap.

### 10. Track A from the ticketing lens (worktree seam)

I lead with ticketing, but the worktree track has a clean ticketing expression,
and naming it is the concrete integration seam:

- **A worktree session is a subtask.** The `worktree` frontmatter field
  (`{ path, branch, base }`) is the link. Lifecycle maps to ticket status:
  `ready`→(claim)→`in-progress` (worktree created, path recorded)→`review`
  (merged to integration branch)→`done` (pruned). `prune` on a `done` subtask is
  safe; on anything else it is a guarded op requiring confirmation.
- **The shared task doc is the set of subtask files.** No separate claim doc is
  needed: claiming a subtask = the frontmatter claim from S-042 above, and the
  `worktree` field makes isolation visible to all implementers. This is the
  explicit seam the prompt asks for between worktree mechanics and Parley's
  ideas/quorum model: quorum decides *what* (the epic/stories); subtask tickets +
  worktree metadata coordinate *who builds what where*.
- **Decomposition boundary = file sets.** The `files` field on each subtask is
  what makes worktree parallelism safe: two subtasks with disjoint `files` sets
  can run in separate worktrees with zero merge conflict by construction. The
  `worktree-sessions` skill should refuse to start two concurrent sessions whose
  `files` sets intersect (or require explicit override), because that is the
  collision the worktree pattern exists to prevent.
- **Runtime overlap.** Where a runtime (Claude Code/Cursor/Copilot CLI) already
  creates worktrees, the skill should *detect and adopt* the runtime's worktree
  (record it in the subtask's `worktree` field) rather than create a second one.
  The skill adds the **claim + file-set-disjointness + lifecycle-as-ticket-state**
  layer the runtimes do not provide; it does not re-implement `git worktree`.

## Concerns / open questions

- **Gherkin vs bullets weighting.** I propose Gherkin for behavioural and
  measurable bullets for NFRs, but the line is fuzzy for integration-logic ACs.
  Open: should integration ACs default to Gherkin or to a "Given preconditions /
  When command / Then exit+output" bullet form? (Leaning bullet form for
  command-shaped ACs, Gherkin for state-shaped.)
- **Mirror sync direction policy.** One-way (file→tracker) is safe but ignores
  humans who live in Jira/Linear. A `--pull` with `mirror-owned` fields is
  proposed, but the exact reconciliation-on-conflict behaviour (tracker wins?
  file wins? both surface as an `Open questions` entry?) is undecided. This is
  the biggest policy risk and needs a decision before the skill ships.
- **Field coverage across trackers.** Not every tracker supports custom fields
  for `files`/`apis`/`worktree`. Graceful degradation keeps the file whole but
  the mirror loses fidelity. Open: is a degraded mirror acceptable, or do we
  require a minimum tracker capability and refuse to sync otherwise?
- **DoD as AC-id list vs tracker-native checklists.** Some trackers have their
  own DoD checklists that won't accept arbitrary id refs. Open: project DoD as
  a checklist whose items are the AC *titles* (portable) rather than ids (opaque
  in a tracker UI)?
- **AI-directive section ownership.** If `[A] Agent directives` is too
  prescriptive it becomes a spec that humans feel they can't own or edit. Open:
  who is the authoritative editor of `[A]` — the author, or the first agent to
  claim? Leaning author, with agents appending `Do not` constraints only via an
  `Open questions [needs:technical]` that the author promotes.
- **Gap-scan vs author burden.** A strict six-check gap-scan may make tickets
  slow to author by hand, pushing authoring onto AI. That is probably desirable
  (the skill can draft from a brief), but it shifts who writes tickets and needs
  to be stated as a consequence, not hidden.

## Risks

- **Three-audience bloat → each audience skim-reads.** A file with `[B]`, `[T]`,
  `[A]` plus ACs plus metadata is long; a stakeholder may read only `At a glance`
  and miss a `[T]` constraint that changes the value. Mitigation: `At a glance`
  is mandatory and must surface any constraint that materially alters the
  outcome; audience tags let tooling render a filtered view per reader.
- **AI fills gaps anyway.** The gap-scan rule is only as strong as the agent's
  discipline to actually stop. If an agent skips the scan and proceeds, the
  no-assumption guarantee is void. Mitigation: `claim` should *run* gap-scan as
  part of the operation and refuse to set `in-progress` on failure — enforcement
  in the tool, not just the template.
- **Mirror drift corrupts canonical files.** If `--pull` is misconfigured,
  tracker-side edits overwrite carefully authored `[T]`/`[A]` content.
  Mitigation: `mirror-owned` defaults to `[status, assignee]` only; any field not
  listed is never written back; divergence on a non-mirror-owned field surfaces
  as a conflict, not a silent overwrite.
- **Vendor-agnosticism vs useful tracker features.** Refusing tracker-specific
  features (Jira sprints, Linear cycles, GitHub Projects views) to stay generic
  may make the mirror feel anemic and drive users back to native tools,
  abandoning the canonical-file model. Mitigation: generic core + an optional
  per-tracker "projection plugin" layer that adds tracker-specific fields to the
  mirror without polluting the canonical file.
- **Gherkin-for-NFRs awkwardness, if the rule slips.** If authors start forcing
  perf/security into Gherkin to feel "complete", criteria become theatrical and
  unmeasurable. Mitigation: the validator flags a Gherkin AC whose `Then` has no
  machine-checkable assertion as a warning; measurable bullets stay first-class.
- **Worktree-as-subtask coupling is tight.** Tying worktree lifecycle to ticket
  status means a worktree can't outlive its subtask and vice versa without
  special handling (e.g. a subtask paused mid-build). Mitigation: allow a
  `paused` status that keeps the worktree alive without claiming active build
  resources; document it explicitly in the skill.
