---
idea: addon-skills-research
drafter: claude-1
date: 2026-06-22
status: final
supersedes: consensus.md
---

## Purpose

A research/design doc concrete enough to author two vendor-neutral companion skills
for the Parley Deck ecosystem:
- **`worktree-sessions`** — run multiple Parley sessions / parallel Phase-5
  implementers over one git repo via worktrees, without collisions.
- **`ticket-tracker`** — author epics / user stories / technical subtasks for any
  tracker (Jira, Linear, GitHub Issues, Trello, kanban), readable by business people,
  technical people, AND the AI agents that implement them.

Both follow **core-stays-thin / skill-carries-mechanics** and stay vendor/tracker/
runtime-agnostic. Distribution: core + opt-in plugins, install-all-by-default with
`--no-addons` / `--only <skill>`, independent per-skill versions.

## Context

Synthesized from four round-01 lenses (claude-1 integration seams, codex-1 worktree
mechanics + templates, hermes-1 three-audience ticketing + gap-scan, antigravity-1
standards: INVEST/Gherkin/DoD/story-splitting), grounded by a June-2026 web pass.
Both skills share one seam: **disjoint file sets** are what make parallel work safe —
the ticket's `files` field defines the boundary, the worktree enforces the isolation.

---

# Skill A — `worktree-sessions`

## A.1 Model
One base repo, shared `.git`; each session/implementer gets its own working dir +
branch + index + HEAD. Isolation turns invisible concurrent-write corruption into
ordinary git merge conflicts. A branch can be checked out in only one worktree.

## A.2 Layout & naming (DECISION: sibling dir, never inside `.git/`)
```text
../worktrees/<slug>-design                 branch idea/<slug>
../worktrees/<slug>-r01-<agent-id>         branch idea/<slug>/round-01-<agent-id>  (optional round isolation)
../worktrees/<slug>-integration            branch integration/<slug>
../worktrees/<slug>-impl-<agent-id>        branch feature/<slug>/<agent-id>
```
Root configurable via `$PARLEY_WORKTREES_DIR` (default `../worktrees`). Do NOT nest
checkouts under `.git/worktrees/` — that is git's own admin metadata.

## A.3 Lifecycle commands
```bash
# preflight
git status --short && git fetch origin && git worktree list && mkdir -p ../worktrees

# provision one implementer worktree off the integration branch
git worktree add -b feature/<slug>/<agent-id> ../worktrees/<slug>-impl-<agent-id> integration/<slug>

# isolate runtime + scope files
cd ../worktrees/<slug>-impl-<agent-id>
cp ../../<repo>/.env.example .env.local      # per-worktree env (gitignored)
git sparse-checkout init --cone
git sparse-checkout set <pkg-paths> parley-deck/ideas/<slug>
git submodule update --init --recursive       # if submodules

# integrate (staging → test → main)
cd ../worktrees/<slug>-integration
git merge --no-ff origin/feature/<slug>/<agent-id> -m "[<agent-id>] <slug>: integrate"
<repo test/lint commands>
# under github-pr: PR integration/<slug> -> main

# cleanup
git worktree remove ../worktrees/<slug>-impl-<agent-id>
git branch -d feature/<slug>/<agent-id>
git worktree prune        # after manual dir deletion
git worktree repair       # if a worktree dir was moved
```
Required `.gitignore`: `.env.local`, `.env.*.local`, `*.sqlite*`, `.worktree-state/`.

## A.4 Decision rules
- USE worktrees when: ≥2 sessions/implementers touch the repo concurrently; Phase-5
  needs tests/installs/local state; agents need different env/ports/DB/branch heads;
  work splits by feature boundary.
- SKIP when: an agent only writes one owned protocol artifact (no code); trivial
  sequential same-file edit; base branch not clean.
- The skill **refuses/warns** on two concurrent sessions whose `files` sets intersect
  (explicit override required) — intersecting file sets are the collision worktrees exist to prevent.

## A.5 Coordination manifest (two layers)
- **Parley-native:** a "Worktree allocation" table in `IMPLEMENTATION.md` is the lock
  manifest (reuse the canonical file; don't invent a new claim doc):
  ```markdown
  ## Worktree allocation
  Base branch: integration/<slug>   Integration worktree: ../worktrees/<slug>-integration
  | Boundary (file set) | Owner | Branch | Worktree | Status |
  |---|---|---|---|---|
  | <pkg/paths> | <agent-id> | feature/<slug>/<agent-id> | ../worktrees/<slug>-impl-<agent-id> | claimed |
  ## Integration log
  - <date> <agent-id>: merged <branch> at <sha>; checks: <cmd/result>
  ```
  Claim rule: an agent edits only its own/an unclaimed row, naming branch + worktree +
  boundary; intersecting boundaries are split or serialized through the integration owner.
- **With `ticket-tracker`:** each subtask ticket carries `worktree: {path,branch,base}`;
  claim = frontmatter status/assignee/worktree edit. Same concept, ticket layer.

## A.6 Runtime overlap
If the runtime (Claude Code / Cursor / Copilot CLI) already created a worktree,
**detect and adopt** it (record its path), don't create a second. The skill adds the
claim + file-set-disjointness + lifecycle layer the runtimes don't provide; it does
not re-implement `git worktree`.

## A.7 Pitfalls
Same-branch double-checkout (unique branch per agent); untracked env not carried
(copy templates); stale worktree after `rm -rf` (`prune`); shared refs — avoid
`git stash`/force-push as inter-agent workflow (use WIP commits); submodules need
per-worktree init; IDE dual-index (skill writes editor excludes for `../worktrees/`);
**dependency lock files** (`go.sum`/`package-lock.json`) conflict on integration —
serialize dependency changes or regenerate on the integration branch.

## A.8 Core seam (separate meta-protocol-change, NOT this idea)
One §6/Phase-5 note: "parallel implementers MAY be isolated in git worktrees;
coordination uses the shared `IMPLEMENTATION.md` lock manifest; merge via an
integration branch; worktrees never change canonical artifact ownership."

---

# Skill B — `ticket-tracker`

## B.1 Model
Tracker = **mirror**; markdown ticket files are canonical (`tickets/<epic>/<story>/
<subtask>.md`). Sync **one-way file→tracker** by default; `--pull` writes back only
fields flagged `mirror-owned` (default `[status, assignee]`). Graceful degradation:
a field the tracker lacks is dropped from the mirror, never from the file.

## B.2 One file, three audiences (template skeleton)
```markdown
---
id: E-001 | S-001 | T-001          # immutable, never reused
type: epic | story | subtask
title: <imperative, <=80 chars>
parent: <id or n/a>
status: draft | ready | in-progress | review | done | blocked | paused | dropped
assignee: <human | agent:<id> | n/a>
priority: p0 | p1 | p2 | p3
labels: [domain:<x>, nfr:<y>]
files: [<path-or-glob>]            # the parallelism boundary; AI scopes reads to these
apis: [<openapi/doc ref> | n/a]
arch: [<doc ref> | n/a]
worktree: n/a | {path,branch,base}
dod: [AC-1, AC-2]                  # checklist of AC ids
mirror-owned: [status, assignee]
canonical_source: parley-deck/ideas/<slug>/FINAL.md@<git-sha>
---
# <title>
## At a glance            # MANDATORY 2-4 lines, all audiences: outcome · scope · done-when
## [B] Business           # value, who benefits, why now, success indicator
## [T] Technical          # systems, interfaces, constraints, NFRs
## [A] Agent directives   # self-contained scope, explicit Do / Do-not, assumption policy
## Acceptance criteria    # AC-1..N, tagged, Gherkin or measurable+Verify
## Non-goals              # explicit out-of-scope (prevents AI scope creep)
## Dependencies / blocks
## Open questions         # each tagged [needs:business|technical|decision]
```
Audience tags are sections of one file (not separate docs); a reader skims their tag,
the AI reads all three + frontmatter. "At a glance" must surface any constraint that
materially alters the outcome. (Full epic/story/subtask instances: see codex-1 and
hermes-1 round-01 — used verbatim as the skill's `templates/`.)

## B.3 Acceptance criteria (hybrid, tagged, id'd)
- Behavioural → **Gherkin** (Given/When/Then), one scenario per `AC-N`.
- Non-functional → **measurable bullet + mandatory `Verify: <single command>`**.
  No `Verify:` command ⇒ not measurable ⇒ rewrite.
- Every AC carries `[B]`/`[T]`/`[A]`. Required: ≥1 happy-path AC AND ≥1 edge/error/
  offline AC (or explicit `n/a (reason)`); NFR ACs required where the story touches
  perf/security/availability and forbidden from being implied.

## B.4 No-assumption gap-scan (highest-leverage rule, tool-enforced)
Before claiming any ticket, the agent runs in order, stopping at first failure:
1. frontmatter present + `validate` exits 0; 2. required fields non-empty + `parent`
resolves + `dod` refs exist; 3. every AC has a tag + Gherkin or `Verify:`; 4. ≥1 edge/
error AC or `n/a (reason)`; 5. `[A]` non-empty with ≥1 `Do not`; 6. `files`/`apis`/
`arch` populated or `n/a`. Any failure → `BLOCKED: <check>` or an `Open questions
[needs:…]`; never proceed by guessing. `claim` RUNS the scan and refuses
`in-progress` on failure (enforcement in the tool, not just the doc).

## B.5 Quality & decomposition
- **INVEST** governs stories (Independent → parallel-safe; Small → fits context;
  Testable → clear assertions). Story-splitting by **vertical slice → disjoint file
  sets** (the worktree seam).
- **DoD = checklist of AC ids**, per type (epic: all child stories done + epic ACs
  pass `Verify`; story: all ACs pass + `validate` 0 + claim released; subtask: ACs
  pass + tests green on integration branch).
- **Readiness lint / `validate`:** reject vague tickets ("improve/support/handle"
  with no observable criteria) and AI-unready tickets.

## B.6 FINAL.md → tracker mapping
```text
idea 00-prompt + non-goals      -> epic [B] + Non-goals
FINAL.md design / value streams  -> stories (one per observable behavior)
FINAL.md observable criteria     -> story acceptance criteria
FINAL.md context/orientation     -> [T] + AI references (files/apis/arch)
IMPLEMENTATION.md checklist      -> subtasks
review-consensus deferred items  -> follow-up tickets
```
`FINAL.md` canonical for *design*; tickets canonical for *work-state*. No overlap.

## B.7 Vendor-neutrality
Neutral field schema → mapped to each tracker at the edge; optional per-tracker
"projection plugin" adds native fields (Jira sprints, Linear cycles) to the mirror
without polluting the canonical file. The skill authors text + emits a payload;
actual create/update via MCP/API is an opt-in connector add-on (no auth in core).

## B.8 Core seam (separate meta-protocol-change, NOT this idea)
One note: "external trackers are mirrors; a ticket cannot override `FINAL.md`;
requirement changes found in a tracker return through a Parley round / inbox / new idea."

---

## Validation (what "research done" means here)
- [x] Both skills have: a model, templates, command/field schemas, decision rules,
      pitfalls, runtime-overlap stance, and an explicit parley-deck seam.
- [x] Worktree path correctness resolved (sibling dir, not `.git/`).
- [x] Three-audience + no-assumption + Gherkin/measurable split specified.
- [x] Vendor/tracker/runtime-agnostic throughout.
- Concrete templates to lift into `templates/` live in codex-1 + hermes-1 round-01.

## Idempotence & recovery
This idea produces design only — no code/protocol mutation. Building the two skills
and the two core carve-outs are separate follow-up ideas seeded from this FINAL.

## Outcomes & retrospective (next steps)
1. Author `worktree-sessions` SKILL.md from §A. 2. Author `ticket-tracker` SKILL.md +
`templates/` + `validate` from §B. 3. Installer `--no-addons`/`--only`. 4. One
meta-protocol-change idea for the two thin core seams (§A.8, §B.8). 5. Skill
genericity must-fix (Transport placeholder) can ride the same release.
