---
agent: codex-1
idea: addon-skills-research
round: 1
date: 2026-06-22
---

## Summary

The `worktree-sessions` skill should be a narrow operational skill: it should teach agents how to allocate, name, merge, and clean up isolated Git worktrees for Parley sessions and Phase-5 implementers, without changing Parley Deck's canonical artifact model. Its highest-value contribution is not "agents can use worktrees" in the abstract; it is the exact branch/worktree discipline that prevents two agents from checking out the same branch, trampling untracked runtime state, or merging half-tested parallel work directly to `main`.

The `ticket-tracker` skill should treat external trackers as mirrors of canonical Parley files, not sources of truth. Tickets need three-layer writing: business value, technical constraints, and an AI implementation contract with explicit acceptance criteria, file/API hints, non-functional requirements, and no hidden assumptions.

## Proposed approach (concrete commands / decision-rules / templates)

### Track A: `worktree-sessions`

#### Core model

Use one base repository with shared Git object storage and multiple isolated working directories:

```bash
git worktree list
git rev-parse --git-common-dir
```

Important mechanics:

- Each worktree has its own working tree, index, `HEAD`, and checked-out branch.
- Worktrees share the same object database and refs through the common `.git` directory.
- One branch cannot be checked out in two worktrees at the same time.
- Local runtime files are per-worktree if they live in the working tree, for example `.env.local`, SQLite files, build output, and package install directories.
- Repository-wide refs such as `refs/stash` are shared; avoid `git stash` as an inter-agent workflow. Prefer WIP commits on the agent branch.

#### Decision rules

Use worktrees when any of these are true:

- More than one Parley session or implementer will modify the same repo concurrently.
- A Phase-5 implementation needs tests, generated files, dependency installs, or local runtime state.
- Agents need different environment variables, ports, databases, or branch heads.
- The work splits cleanly by feature boundary, package, service, or ownership area.

Do not add worktree overhead when:

- The agent is only writing one owned Parley protocol artifact and no code.
- The change is a trivial same-file edit that needs sequential coordination anyway.
- The repo is not clean enough to establish a known base branch.

Branch naming:

```text
idea/<slug>                              # design PR branch, one owner/integration worktree
idea/<slug>/round-01-<agent-id>          # optional isolated round branch
integration/<slug>                       # staging branch for parallel implementation
feature/<slug>/<agent-id>                # per-implementer branch
```

Worktree directory naming:

```text
../worktrees/<slug>-design
../worktrees/<slug>-r01-<agent-id>
../worktrees/<slug>-integration
../worktrees/<slug>-impl-<agent-id>
```

#### Setup for this idea's style of repo

From the main repo directory:

```bash
git status --short
git fetch origin
git worktree list
mkdir -p ../worktrees
```

If the design branch already exists remotely:

```bash
git fetch origin idea/addon-skills-research:refs/remotes/origin/idea/addon-skills-research
git worktree add ../worktrees/addon-skills-research-design origin/idea/addon-skills-research
cd ../worktrees/addon-skills-research-design
git switch -c idea/addon-skills-research
```

If the local design branch already exists and is not checked out anywhere else:

```bash
git worktree add ../worktrees/addon-skills-research-design idea/addon-skills-research
```

If it is already checked out in the main worktree, treat the main worktree as the design integration worktree. Do not try to check out `idea/addon-skills-research` in a second worktree.

#### Round-1 isolation commands

For stronger round-1 independence under the active GitHub PR transport, each agent can work on a sub-branch and later merge into `idea/addon-skills-research`.

Codex worktree:

```bash
git fetch origin
git worktree add -b idea/addon-skills-research/round-01-codex-1 \
  ../worktrees/addon-skills-research-r01-codex-1 \
  origin/idea/addon-skills-research

cd ../worktrees/addon-skills-research-r01-codex-1
git status --short
```

Other participants follow the same pattern:

```bash
git worktree add -b idea/addon-skills-research/round-01-claude-1 \
  ../worktrees/addon-skills-research-r01-claude-1 \
  origin/idea/addon-skills-research

git worktree add -b idea/addon-skills-research/round-01-hermes-1 \
  ../worktrees/addon-skills-research-r01-hermes-1 \
  origin/idea/addon-skills-research

git worktree add -b idea/addon-skills-research/round-01-antigravity-1 \
  ../worktrees/addon-skills-research-r01-antigravity-1 \
  origin/idea/addon-skills-research
```

Each agent writes only its own file, then commits and pushes its sub-branch:

```bash
git add parley-deck/ideas/addon-skills-research/round-01/codex-1.md
git commit -m "[codex-1] addon-skills-research: round-01 analysis"
git push -u origin idea/addon-skills-research/round-01-codex-1
```

The facilitator or design branch owner merges sub-branches into the design branch after all round-1 outputs exist:

```bash
cd ../worktrees/addon-skills-research-design
git fetch origin
git merge --no-ff origin/idea/addon-skills-research/round-01-codex-1 \
  -m "[codex-1] addon-skills-research: merge codex-1 round-01"
git merge --no-ff origin/idea/addon-skills-research/round-01-claude-1 \
  -m "[claude-1] addon-skills-research: merge claude-1 round-01"
git merge --no-ff origin/idea/addon-skills-research/round-01-hermes-1 \
  -m "[hermes-1] addon-skills-research: merge hermes-1 round-01"
git merge --no-ff origin/idea/addon-skills-research/round-01-antigravity-1 \
  -m "[antigravity-1] addon-skills-research: merge antigravity-1 round-01"
git push origin idea/addon-skills-research
```

This keeps round-1 independent while still preserving the one design PR branch.

#### Phase-5 parallel implementation commands

For multiple implementers over one repo, use an integration branch plus per-agent implementation branches. Do not put multiple agents directly on `feature/addon-skills-research`.

Create the integration branch and worktree:

```bash
git fetch origin
git worktree add -b integration/addon-skills-research \
  ../worktrees/addon-skills-research-integration \
  origin/main
cd ../worktrees/addon-skills-research-integration
git push -u origin integration/addon-skills-research
```

Create one worktree per implementer:

```bash
cd /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-cli

git worktree add -b feature/addon-skills-research/codex-1 \
  ../worktrees/addon-skills-research-impl-codex-1 \
  integration/addon-skills-research

git worktree add -b feature/addon-skills-research/claude-1 \
  ../worktrees/addon-skills-research-impl-claude-1 \
  integration/addon-skills-research

git worktree add -b feature/addon-skills-research/hermes-1 \
  ../worktrees/addon-skills-research-impl-hermes-1 \
  integration/addon-skills-research

git worktree add -b feature/addon-skills-research/antigravity-1 \
  ../worktrees/addon-skills-research-impl-antigravity-1 \
  integration/addon-skills-research
```

Each implementer claims a feature boundary in `IMPLEMENTATION.md`, then works only in its worktree:

```bash
cd ../worktrees/addon-skills-research-impl-codex-1
git status --short

# edit code and the owned IMPLEMENTATION.md claim/update
git add <changed-files>
git commit -m "[codex-1] addon-skills-research: implement <boundary>"
git push -u origin feature/addon-skills-research/codex-1
```

Integration branch owner merges and tests sequentially:

```bash
cd ../worktrees/addon-skills-research-integration
git fetch origin

git merge --no-ff origin/feature/addon-skills-research/codex-1 \
  -m "[codex-1] addon-skills-research: integrate codex-1 implementation"

git merge --no-ff origin/feature/addon-skills-research/claude-1 \
  -m "[claude-1] addon-skills-research: integrate claude-1 implementation"

git merge --no-ff origin/feature/addon-skills-research/hermes-1 \
  -m "[hermes-1] addon-skills-research: integrate hermes-1 implementation"

git merge --no-ff origin/feature/addon-skills-research/antigravity-1 \
  -m "[antigravity-1] addon-skills-research: integrate antigravity-1 implementation"

git status --short
git log --oneline --decorate --max-count=12
```

Then run the project checks from `FINAL.md` or `IMPLEMENTATION.md`:

```bash
npm test
npm run lint
```

Replace those with the repo's real checks if it is not a Node project.

For GitHub transport, the merge to `main` should normally be a PR from `integration/addon-skills-research`:

```bash
git push origin integration/addon-skills-research
gh pr create \
  --base main \
  --head integration/addon-skills-research \
  --title "[addon-skills-research] implementation" \
  --body "Canonical implementation record: parley-deck/ideas/addon-skills-research/IMPLEMENTATION.md"
```

If no GitHub CLI is available, push the branch and create the PR in the host UI. The canonical record remains the Parley Deck files.

#### Cleanup commands

Before cleanup, every worktree must be clean or intentionally preserved:

```bash
git -C ../worktrees/addon-skills-research-impl-codex-1 status --short
git -C ../worktrees/addon-skills-research-integration status --short
git worktree list
```

Remove merged per-agent worktrees:

```bash
git worktree remove ../worktrees/addon-skills-research-impl-codex-1
git worktree remove ../worktrees/addon-skills-research-impl-claude-1
git worktree remove ../worktrees/addon-skills-research-impl-hermes-1
git worktree remove ../worktrees/addon-skills-research-impl-antigravity-1
```

Delete merged local branches:

```bash
git branch -d feature/addon-skills-research/codex-1
git branch -d feature/addon-skills-research/claude-1
git branch -d feature/addon-skills-research/hermes-1
git branch -d feature/addon-skills-research/antigravity-1
```

Delete merged remote branches:

```bash
git push origin --delete feature/addon-skills-research/codex-1
git push origin --delete feature/addon-skills-research/claude-1
git push origin --delete feature/addon-skills-research/hermes-1
git push origin --delete feature/addon-skills-research/antigravity-1
```

Prune stale metadata:

```bash
git worktree prune
git worktree list
```

If a worktree directory was moved manually:

```bash
git worktree repair
```

#### Coordination template for `IMPLEMENTATION.md`

The worktree skill should not invent a new canonical claim file when Parley already has `IMPLEMENTATION.md`. Add a compact coordination section there:

```markdown
## Worktree allocation

Base branch: integration/addon-skills-research
Integration worktree: ../worktrees/addon-skills-research-integration

| Boundary | Owner | Branch | Worktree | Status |
| --- | --- | --- | --- | --- |
| <feature/package/service> | codex-1 | feature/addon-skills-research/codex-1 | ../worktrees/addon-skills-research-impl-codex-1 | claimed |
| <feature/package/service> | claude-1 | feature/addon-skills-research/claude-1 | ../worktrees/addon-skills-research-impl-claude-1 | open |

## Integration log

- 2026-06-22 codex-1: merged feature/addon-skills-research/codex-1 at <sha>; checks: <command/result>.
```

Claim rule:

- An agent claims by editing only its own row or an unclaimed row.
- A claim must name the branch, worktree path, and boundary.
- If two agents need the same files, split the work before implementation starts or serialize those edits through the integration branch owner.

#### Sparse checkout for constrained agents

Sparse checkout is useful when an implementer should only see or edit a package boundary:

```bash
cd ../worktrees/addon-skills-research-impl-codex-1
git sparse-checkout init --cone
git sparse-checkout set packages/cli parley-deck/ideas/addon-skills-research
git status --short
```

Rule: sparse checkout is a guardrail, not authority. `FINAL.md` and `IMPLEMENTATION.md` still define scope.

#### Submodule rule

Each worktree needs its own submodule checkout state:

```bash
cd ../worktrees/addon-skills-research-impl-codex-1
git submodule update --init --recursive
git submodule status --recursive
```

If submodule branches are also edited, give each agent a branch inside the submodule too. Do not let agents share a detached submodule worktree for writes.

#### Environment rule

Each worktree gets its own local runtime state:

```bash
cp .env.example .env.local
printf '\nPORT=31%02d\n' 1 >> .env.local
```

Required `.gitignore` entries:

```gitignore
.env.local
.env.*.local
*.sqlite
*.sqlite3
.worktree-state/
```

If the repo uses shared services, assign per-worktree ports, database names, cache dirs, and test temp dirs in `IMPLEMENTATION.md`.

#### Seam with Parley Deck core

Put this in `worktree-sessions`:

- Worktree commands and cleanup commands.
- Branch naming and worktree path naming.
- Integration branch pattern.
- Sparse checkout and submodule handling.
- Environment isolation checklist.
- Troubleshooting for stale worktrees, branch checkout conflicts, and dirty worktrees.

Keep this as only a thin note in Parley Deck core or protocol:

- Phase-5 implementers may isolate work in Git worktrees.
- Worktree use does not change canonical artifact ownership.
- `FINAL.md` and `IMPLEMENTATION.md` remain authoritative.
- External PR/MR and tracker surfaces are mirrors.

### Track B: `ticket-tracker`

#### Core model

The tracker is a mirror. Parley files remain canonical:

```text
parley-deck/ideas/<slug>/00-prompt.md       # problem and participants
parley-deck/ideas/<slug>/FINAL.md           # authoritative design/spec
parley-deck/ideas/<slug>/IMPLEMENTATION.md  # implementation state and validation
```

Ticket sync direction:

```text
FINAL.md -> epic/story/subtask tickets -> comments/status mirrors -> back to Parley only through a new round, inbox escalation, or IMPLEMENTATION.md update
```

Do not let tracker comments silently change requirements.

#### Generic field map

Use tracker-neutral names, then map them to Jira, Linear, GitHub Issues, Trello, or another backend:

```yaml
canonical_source:
  idea: parley-deck/ideas/addon-skills-research/FINAL.md
  revision: <git-sha>
tracker:
  provider: generic
  external_id: <filled-after-create>
  url: <filled-after-create>
item:
  type: epic | story | subtask
  title: <short outcome>
  status: proposed | ready | in_progress | blocked | review | done
  priority: low | normal | high | urgent
  labels: [parley-deck, <slug>, <domain>]
  parent: <epic-or-story-id-or-null>
  owner: <human-or-agent-or-team>
  due: <YYYY-MM-DD-or-null>
ai_contract:
  implementation_agent: <agent-id-or-unassigned>
  file_hints: [<path-or-glob>]
  api_refs: [<openapi-or-doc-url>]
  architecture_refs: [<path-or-url>]
  test_refs: [<path-or-command>]
  assumption_policy: ask-before-assuming
```

#### Epic template

```markdown
# Epic: <business outcome>

## Business context

Who benefits:
Why this matters:
Success metric:
Deadline or timing pressure:

## Scope

In scope:
- <capability or workflow>

Out of scope:
- <explicit non-goal>

## Technical orientation

Systems touched:
Primary constraints:
Security/privacy constraints:
Performance/reliability constraints:
Migration or compatibility constraints:

## Stories

- [ ] <story title or tracker link>
- [ ] <story title or tracker link>

## Acceptance criteria

- [ ] The epic is complete when every linked story is done.
- [ ] The observable success metric is measured by <metric/source>.
- [ ] Documentation/runbook updates are complete.

## Definition of Done

- [ ] All stories meet their acceptance criteria.
- [ ] Required tests/checks pass and are linked.
- [ ] Known limitations are documented.
- [ ] Parley canonical files are updated with final tracker links.
```

#### User story template

```markdown
# Story: <user-visible behavior>

## Business reader

As a <user/persona>,
I want <capability>,
so that <business value>.

Value:
Priority:
Non-goals:

## Technical reader

Systems/components:
Data touched:
Interfaces/APIs:
Dependencies:
Constraints:
NFRs:
- Performance:
- Security/privacy:
- Accessibility:
- Observability:
- Offline/error behavior:

## AI implementation contract

Canonical source: parley-deck/ideas/<slug>/FINAL.md at <git-sha>
Allowed files/areas:
- <path-or-glob>

Do not modify:
- <path-or-glob>

Relevant references:
- API/OpenAPI:
- Architecture:
- Existing tests:
- Similar implementation:

Assumption policy:
- If any required endpoint, schema, UX state, or error behavior is missing, stop and ask through the configured Parley/tracker channel.

## Acceptance criteria

Scenario: <happy path>
Given <initial state>
When <action>
Then <observable result>

Scenario: <error or edge case>
Given <initial state>
When <failure or edge action>
Then <observable fallback/error/result>

Measurable criteria:
- [ ] <specific testable statement>
- [ ] <specific testable NFR>

## Test plan

Commands:
- `<command>`

Manual checks:
- <check>

## Definition of Done

- [ ] Code and tests are complete.
- [ ] Acceptance criteria are demonstrably met.
- [ ] Error, empty, loading, permission, and offline states are handled or explicitly out of scope.
- [ ] Observability/logging is added where required.
- [ ] Documentation or release notes are updated if user-visible.
- [ ] Implementation notes are reflected in `IMPLEMENTATION.md`.
```

#### Technical subtask template

```markdown
# Subtask: <technical change>

Parent story:
Canonical source:
Owner:

## Objective

Implement <specific technical outcome> so that <parent story capability> works.

## Scope

Files likely touched:
- <path-or-glob>

Interfaces/contracts:
- <function/API/schema>

Out of scope:
- <explicit exclusions>

## Steps

- [ ] <small implementation step>
- [ ] <small implementation step>
- [ ] Add or update tests for <behavior>

## Acceptance criteria

- [ ] <specific observable technical result>
- [ ] `<test command>` passes.
- [ ] No changes outside the allowed scope unless recorded in `IMPLEMENTATION.md`.

## AI notes

Known pitfalls:
- <pitfall>

Ask before:
- changing public API
- changing database schema
- broadening file scope
- weakening validation/security behavior
```

#### Ticket decomposition rules

- Create one epic per coherent business outcome, not per technical layer.
- Create one story per independently valuable user-visible behavior.
- Create subtasks for implementation slices that are not independently valuable to the user but help execution and review.
- Every story must have acceptance criteria before it is marked ready.
- Every AI-assigned ticket must include canonical source path, revision, file hints, constraints, test expectations, and assumption policy.
- If a ticket cannot be implemented without interpreting missing product behavior, it is not ready.
- If a ticket only says "improve", "support", "handle", or "integrate" without observable criteria, rewrite it.

#### Mapping from `FINAL.md` to tracker items

```text
FINAL.md Purpose / user-visible outcome -> epic business context
FINAL.md Observable acceptance criteria -> stories and story criteria
FINAL.md Context & orientation -> technical orientation and AI references
FINAL.md Known risks / de-risking -> story risks, ask-before rules, NFRs
IMPLEMENTATION.md plan/checklist -> subtasks
Review consensus deferred follow-ups -> follow-up stories or subtasks
```

Minimum created tracker set:

```text
1 epic: <slug> outcome
N stories: one per observable behavior or acceptance criterion cluster
M subtasks: one per implementation boundary, migration, test suite, or docs task
```

#### Seam with Parley Deck core

Put this in `ticket-tracker`:

- Tracker-neutral epic/story/subtask templates.
- Field mapping for common trackers.
- Readiness lint rules for AI-implementable tickets.
- Sync rules and examples.
- Import/export guidance.

Keep this as only a thin note in Parley Deck core or protocol:

- Trackers are mirrors of canonical files.
- A tracker ticket cannot override `FINAL.md`.
- Requirement changes discovered in a tracker must be brought back through Parley rounds, inbox escalation, or a new idea.

## Concerns / open questions

- The biggest worktree design choice is whether `worktree-sessions` should prescribe an integration branch named `integration/<slug>` or allow each repo to map that to its existing branch convention. I prefer prescribing the pattern as the default and allowing an override because it is easy to explain and avoids multiple implementers pushing to one `feature/<slug>` branch.
- For design rounds, worktrees are useful but not mandatory. The skill should say that one-file protocol artifacts can use ordinary branch discipline, while Phase-5 implementation should strongly prefer worktrees when parallel work is real.
- The active GitHub transport already has a design PR and implementation PR model. The worktree skill should support that model without requiring GitHub; the PR command examples should be optional mirrors, not the core workflow.
- Ticket templates need enough metadata for agents without becoming unreadable to business users. The best compromise is human-readable markdown sections plus a small YAML metadata block for machine fields.
- The ticketing skill needs a "ready for AI" lint checklist. Without it, generic tracker templates will drift back into vague human-only tickets.
- We should decide whether tracker sync creates tickets directly or only emits markdown/JSON payloads for a connector. For vendor neutrality, the skill should define payloads and decision rules first; connector-specific create/update operations can be add-ons.

## Risks

- Worktree branch collision is the main operational risk: if two agents try to check out the same branch, Git blocks one of them; if they instead share a directory, runtime files and generated artifacts can corrupt each other silently.
- Stale worktrees can leave branches undeletable and confuse agents. The skill needs explicit `git worktree list`, `git worktree remove`, and `git worktree prune` cleanup gates.
- Shared repository refs mean `git stash`, branch deletion, and force-pushes affect more than one worktree. The safer pattern is WIP commits on per-agent branches and no force-push without coordination.
- Per-worktree environment drift can produce false test failures or false passes. `IMPLEMENTATION.md` should record ports, databases, caches, and test commands per worktree when they differ.
- Submodules and sparse checkout can create a false sense of isolation. They constrain files, but they do not change the canonical scope in `FINAL.md`.
- Ticket trackers can become a second source of truth if comments, status changes, or hand-edited acceptance criteria are not synced back into Parley files. The skill must repeat the mirror rule everywhere.
- AI-facing tickets can overfit to implementation hints and suppress better designs. The templates should distinguish "allowed files/likely files" from "must modify exactly these files".
- Generic tracker support can become lowest-common-denominator if the skill avoids all structure. The right generic layer is a stable neutral schema, with provider mappings below it.
