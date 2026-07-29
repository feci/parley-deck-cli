---
idea: skills-cli-install-path
round: 01
agent: codex-1
date: 2026-07-29
---

# Empirical verification

## 1. Evidence

### Scope and executable

Both repositories were on the requested branch:

```text
$ git -C /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-cli branch --show-current
readme-skill-catalogue
$ git -C /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-skill branch --show-current
readme-skill-catalogue
```

The literal `@latest` lookup was attempted first and failed because this environment could
not resolve the npm registry:

```text
$ npx -y skills@latest --version
npm error code ENOTFOUND
npm error syscall getaddrinfo
npm error errno ENOTFOUND
npm error network request to https://registry.npmjs.org/skills failed, reason: getaddrinfo ENOTFOUND registry.npmjs.org
npm error network This is a problem related to network connectivity.
npm error network In most cases you are behind a proxy or have bad network settings.
npm error network
npm error network If you are behind a proxy, please make sure that the 'proxy' config is set properly.  See: 'npm help config'
npm error Log files were not written due to an error writing to the directory: /Users/tomasfecko/.npm/_logs
npm error You can rerun the command with `--loglevel=verbose` to see the logs in your terminal
```

A previously resolved, published `skills` executable was present in the local npx cache:

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs --version
1.5.20
```

All discovery findings below therefore apply to `skills` 1.5.20. Whether the registry's
current `skills@latest` tag still resolves to 1.5.20 is **NOT TESTED**. Whether the remote
`feci/parley-deck-skill` clone behaves identically is **NOT TESTED** because both npm and
GitHub network access were unavailable. The same CLI was run against a local copy of the
requested branch.

The source was copied to `/tmp/skills-probe-codex-1.FtB7Od/`, excluding only `.git` and
`node_modules`. Four independent copies were prepared:

```text
$ rsync -a --exclude=.git --exclude=node_modules /Volumes/My\ Shared\ Files/AI_WORKSPACE/parley-deck/parley-deck-skill/ /tmp/skills-probe-codex-1.FtB7Od/baseline/
$ cp -R /tmp/skills-probe-codex-1.FtB7Od/baseline /tmp/skills-probe-codex-1.FtB7Od/skills-container
$ cp -R /tmp/skills-probe-codex-1.FtB7Od/baseline /tmp/skills-probe-codex-1.FtB7Od/root-removed
$ cp -R /tmp/skills-probe-codex-1.FtB7Od/baseline /tmp/skills-probe-codex-1.FtB7Od/marketplace
```

Those commands produced no output.

The `skills-container` variant copied the four existing add-on directories into a standard
container while retaining the root `SKILL.md`:

```text
$ mkdir /tmp/skills-probe-codex-1.FtB7Od/skills-container/skills
$ cp -R /tmp/skills-probe-codex-1.FtB7Od/skills-container/addons/. /tmp/skills-probe-codex-1.FtB7Od/skills-container/skills/
```

The `root-removed` variant changed exactly one discovery input:

```text
$ rm /tmp/skills-probe-codex-1.FtB7Od/root-removed/SKILL.md
```

The `marketplace` variant retained the root and add-ons and added this manifest:

```json
{
  "plugins": [
    {
      "name": "parley-deck",
      "source": "./",
      "skills": [
        "./SKILL.md",
        "./addons/parley-design",
        "./addons/parley-design-check",
        "./addons/parley-tracker",
        "./addons/parley-worktrees"
      ]
    }
  ]
}
```

For readable transcripts, each discovery command pipes output through a Perl expression that
removes ANSI cursor-control bytes only. The text below is the exact output of the displayed
command after that explicit filter.

### Result summary

| Probe | Root `SKILL.md` | Add-on location or declaration | Result |
| --- | --- | --- | --- |
| Baseline | present | `addons/<name>/SKILL.md` | 1: core only |
| Standard container | present | copied under `skills/<name>/SKILL.md` | 1: core only |
| Marketplace manifest | present | all five paths declared | 1: core only |
| Root removed | absent | unchanged `addons/<name>/SKILL.md` | 4: all add-ons |
| Baseline + `--full-depth` | present | unchanged `addons/<name>/SKILL.md` | 5: core plus all add-ons |

The observation selects cause **(b)**. In 1.5.20 the root `SKILL.md` shadows nested
discovery. `addons/` is not a normal priority container, but it is not categorically
unsearched: when the root is removed, recursive fallback discovers all four add-ons there.
The root also suppresses skills copied into the documented `skills/` container and paths
declared through `.claude-plugin/marketplace.json`.

#### Baseline: root plus `addons/`

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs add /tmp/skills-probe-codex-1.FtB7Od/baseline --list 2>&1 | perl -pe 's/\e\[[0-9;?]*[A-Za-z]//g'

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /tmp/skills-probe-codex-1.FtB7Od/baseline
│
◇  Local path validated
│
◇  Found 1 skill

│
◇  Available Skills
│
│    parley-deck
│
│      Run Parley Deck multi-agent idea, implementation, review, or consensus workflows through local CLI agents, using defaulted or user-overridden transport: local files, GitHub PRs, or GitLab MRs. Use when a user wants a task, design, implementation plan, or code review to be independently analyzed by multiple headless or interactive agents according to parley-deck/COOPERATION.md, with each participant writing its own canonical artifacts under parley-deck/ideas/.

│
└  Use --skill <name> to install specific skills
```

#### Standard `skills/` container while retaining the root

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs add /tmp/skills-probe-codex-1.FtB7Od/skills-container --list 2>&1 | perl -pe 's/\e\[[0-9;?]*[A-Za-z]//g'

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /tmp/skills-probe-codex-1.FtB7Od/skills-container
│
◇  Local path validated
│
◇  Found 1 skill

│
◇  Available Skills
│
│    parley-deck
│
│      Run Parley Deck multi-agent idea, implementation, review, or consensus workflows through local CLI agents, using defaulted or user-overridden transport: local files, GitHub PRs, or GitLab MRs. Use when a user wants a task, design, implementation plan, or code review to be independently analyzed by multiple headless or interactive agents according to parley-deck/COOPERATION.md, with each participant writing its own canonical artifacts under parley-deck/ideas/.

│
└  Use --skill <name> to install specific skills
```

#### `.claude-plugin/marketplace.json` while retaining the root

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs add /tmp/skills-probe-codex-1.FtB7Od/marketplace --list 2>&1 | perl -pe 's/\e\[[0-9;?]*[A-Za-z]//g'

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /tmp/skills-probe-codex-1.FtB7Od/marketplace
│
◇  Local path validated
│
◇  Found 1 skill

│
◇  Available Skills
│
│    parley-deck
│
│      Run Parley Deck multi-agent idea, implementation, review, or consensus workflows through local CLI agents, using defaulted or user-overridden transport: local files, GitHub PRs, or GitLab MRs. Use when a user wants a task, design, implementation plan, or code review to be independently analyzed by multiple headless or interactive agents according to parley-deck/COOPERATION.md, with each participant writing its own canonical artifacts under parley-deck/ideas/.

│
└  Use --skill <name> to install specific skills
```

#### Root removed, unchanged `addons/`

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs add /tmp/skills-probe-codex-1.FtB7Od/root-removed --list 2>&1 | perl -pe 's/\e\[[0-9;?]*[A-Za-z]//g'

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /tmp/skills-probe-codex-1.FtB7Od/root-removed
│
◇  Local path validated
│
◇  Found 4 skills

│
◇  Available Skills
│
│    parley-design
│
│      Produce a design system with several independent participants and then apply it, without the result reading as machine-made. Use when a Parley Deck idea creates a new visual world, changes a ratified design rule, or needs an interface audited against a contract instead of against taste. Vendor- and surface-agnostic companion add-on to the parley-deck skill: it ships the PDS/1.0 protocol (typed design artifacts, distinctness and coherence gates, evidence tiers, rule authority, waivers, conformance levels) as pure markdown with zero runtime dependencies, and it never changes canonical artifact ownership.
│
│    parley-design-check
│
│      Run the checkable part of the PDS/1.0 design doctrine against files on disk: design artifacts, DTCG token documents, stylesheets and markup. Use when a design run needs its rules enforced reproducibly rather than argued, when a conformance level claim has to be verified, or when a review wants findings that are stable and diffable across runs. Separable enforcement companion to the parley-design add-on: it reads that skill's rule registry, refuses to check rules when the registry is absent, reports what it cannot judge instead of passing it, and has no runtime dependencies and no network access.
│
│    parley-tracker
│
│      Author epics, user stories, and technical subtasks as canonical markdown files that read well for business people, technical people, AND the AI agents that implement them — then mirror them into any tracker (Jira, Linear, GitHub Issues, GitLab, Trello, a kanban board). Use when a Parley Deck idea or any backlog needs vendor-neutral, no-assumption, AI-implementable tickets with hybrid acceptance criteria and a tool-enforced gap-scan before work starts.
│
│    parley-worktrees
│
│      Allocate, name, isolate, merge, and clean up git worktrees so multiple Parley Deck sessions or parallel Phase-5 implementers can work over one repository without collisions. Use when two or more agents/sessions touch the same repo concurrently, when an implementation needs its own tests/installs/env/local state, or when work splits cleanly along a feature/file-set boundary. Vendor-, tracker-, and runtime-agnostic companion add-on to the parley-deck skill: it teaches worktree mechanics and a file-based claim discipline, and it never changes canonical artifact ownership (FINAL.md and IMPLEMENTATION.md stay authoritative).

│
└  Use --skill <name> to install specific skills
```

#### Unchanged baseline with `--full-depth`

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs add /tmp/skills-probe-codex-1.FtB7Od/baseline --full-depth --list 2>&1 | perl -pe 's/\e\[[0-9;?]*[A-Za-z]//g'

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /tmp/skills-probe-codex-1.FtB7Od/baseline
│
◇  Local path validated
│
◇  Found 5 skills

│
◇  Available Skills
│
│    parley-deck
│
│      Run Parley Deck multi-agent idea, implementation, review, or consensus workflows through local CLI agents, using defaulted or user-overridden transport: local files, GitHub PRs, or GitLab MRs. Use when a user wants a task, design, implementation plan, or code review to be independently analyzed by multiple headless or interactive agents according to parley-deck/COOPERATION.md, with each participant writing its own canonical artifacts under parley-deck/ideas/.
│
│    parley-design
│
│      Produce a design system with several independent participants and then apply it, without the result reading as machine-made. Use when a Parley Deck idea creates a new visual world, changes a ratified design rule, or needs an interface audited against a contract instead of against taste. Vendor- and surface-agnostic companion add-on to the parley-deck skill: it ships the PDS/1.0 protocol (typed design artifacts, distinctness and coherence gates, evidence tiers, rule authority, waivers, conformance levels) as pure markdown with zero runtime dependencies, and it never changes canonical artifact ownership.
│
│    parley-design-check
│
│      Run the checkable part of the PDS/1.0 design doctrine against files on disk: design artifacts, DTCG token documents, stylesheets and markup. Use when a design run needs its rules enforced reproducibly rather than argued, when a conformance level claim has to be verified, or when a review wants findings that are stable and diffable across runs. Separable enforcement companion to the parley-design add-on: it reads that skill's rule registry, refuses to check rules when the registry is absent, reports what it cannot judge instead of passing it, and has no runtime dependencies and no network access.
│
│    parley-tracker
│
│      Author epics, user stories, and technical subtasks as canonical markdown files that read well for business people, technical people, AND the AI agents that implement them — then mirror them into any tracker (Jira, Linear, GitHub Issues, GitLab, Trello, a kanban board). Use when a Parley Deck idea or any backlog needs vendor-neutral, no-assumption, AI-implementable tickets with hybrid acceptance criteria and a tool-enforced gap-scan before work starts.
│
│    parley-worktrees
│
│      Allocate, name, isolate, merge, and clean up git worktrees so multiple Parley Deck sessions or parallel Phase-5 implementers can work over one repository without collisions. Use when two or more agents/sessions touch the same repo concurrently, when an implementation needs its own tests/installs/env/local state, or when work splits cleanly along a feature/file-set boundary. Vendor-, tracker-, and runtime-agnostic companion add-on to the parley-deck skill: it teaches worktree mechanics and a file-based claim discipline, and it never changes canonical artifact ownership (FINAL.md and IMPLEMENTATION.md stay authoritative).

│
└  Use --skill <name> to install specific skills
```

### Installation was also exercised in scratch

`--list` was not the only success signal. This project-scoped installation wrote only below
`/tmp/skills-probe-codex-1.FtB7Od/universal-install-verified/`; it used no `-g` flag:

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs add /tmp/skills-probe-codex-1.FtB7Od/baseline --full-depth --skill '*' --agent claude-code -y 2>&1 | perl -pe 's/\e\[[0-9;?]*[A-Za-z]//g'

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /tmp/skills-probe-codex-1.FtB7Od/baseline
│
◇  Local path validated
│
◇  Found 5 skills
│
●  Installing all 5 skills

│
◇  Installation Summary ─────────────────╮
│                                        │
│  ./.agents/skills/parley-deck          │
│    copy → Claude Code                  │
│                                        │
│  ./.agents/skills/parley-design        │
│    copy → Claude Code                  │
│                                        │
│  ./.agents/skills/parley-design-check  │
│    copy → Claude Code                  │
│                                        │
│  ./.agents/skills/parley-tracker       │
│    copy → Claude Code                  │
│                                        │
│  ./.agents/skills/parley-worktrees     │
│    copy → Claude Code                  │
│                                        │
├────────────────────────────────────────╯
│
◒  Installing skills…◐  Installing skills…◓  Installing skills…◑  Installing skills…◇  Installation complete

│
◇  Installed 5 skills ───────────────────────╮
│                                            │
│  ✓ parley-deck (copied)                    │
│    → ./.claude/skills/parley-deck          │
│  ✓ parley-design (copied)                  │
│    → ./.claude/skills/parley-design        │
│  ✓ parley-design-check (copied)            │
│    → ./.claude/skills/parley-design-check  │
│  ✓ parley-tracker (copied)                 │
│    → ./.claude/skills/parley-tracker       │
│  ✓ parley-worktrees (copied)               │
│    → ./.claude/skills/parley-worktrees     │
│                                            │
├────────────────────────────────────────────╯

│
└  Done!  Review skills before use; they run with full agent permissions.
```

The installed project was then enumerated:

```text
$ node /Users/tomasfecko/.npm/_npx/ac0ed6aa23b37c1e/node_modules/skills/bin/cli.mjs list --agent claude-code --json
[
  {
    "name": "parley-deck",
    "path": "/private/tmp/skills-probe-codex-1.FtB7Od/universal-install-verified/.claude/skills/parley-deck",
    "scope": "project",
    "agents": [
      "Claude Code"
    ],
    "source": "/tmp/skills-probe-codex-1.FtB7Od/baseline",
    "sourceUrl": null,
    "sourceType": "local"
  },
  {
    "name": "parley-design",
    "path": "/private/tmp/skills-probe-codex-1.FtB7Od/universal-install-verified/.claude/skills/parley-design",
    "scope": "project",
    "agents": [
      "Claude Code"
    ],
    "source": "/tmp/skills-probe-codex-1.FtB7Od/baseline",
    "sourceUrl": null,
    "sourceType": "local"
  },
  {
    "name": "parley-design-check",
    "path": "/private/tmp/skills-probe-codex-1.FtB7Od/universal-install-verified/.claude/skills/parley-design-check",
    "scope": "project",
    "agents": [
      "Claude Code"
    ],
    "source": "/tmp/skills-probe-codex-1.FtB7Od/baseline",
    "sourceUrl": null,
    "sourceType": "local"
  },
  {
    "name": "parley-tracker",
    "path": "/private/tmp/skills-probe-codex-1.FtB7Od/universal-install-verified/.claude/skills/parley-tracker",
    "scope": "project",
    "agents": [
      "Claude Code"
    ],
    "source": "/tmp/skills-probe-codex-1.FtB7Od/baseline",
    "sourceUrl": null,
    "sourceType": "local"
  },
  {
    "name": "parley-worktrees",
    "path": "/private/tmp/skills-probe-codex-1.FtB7Od/universal-install-verified/.claude/skills/parley-worktrees",
    "scope": "project",
    "agents": [
      "Claude Code"
    ],
    "source": "/tmp/skills-probe-codex-1.FtB7Od/baseline",
    "sourceUrl": null,
    "sourceType": "local"
  }
]
```

## 2. Recommended layout

I recommend **E: retain the repository layout and make `--full-depth` part of the universal
install contract**.

This is preferable to A-C on the observed evidence:

- A, the marketplace manifest, still found only the root skill.
- B, a copied `skills/` container, still found only the root skill and would add a second
  hand-maintained skill tree plus npm package duplication.
- C would remove the shadow, but the existing installer hard-codes the package-root
  `SKILL.md` as required payload, hashes it for status, copies it for core installs, and
  materializes it as `skills/SKILL.md` for Antigravity. `plugin.json`,
  `gemini-extension.json`, manual installation, and the Codex repository flow also assume
  the root entrypoint. I did not move the real file.
- D is unnecessary because 1.5.20 already has a flag that discovers and installs all five
  in the unchanged tree.

The cost is a mandatory extra flag and a moving dependency on the upstream CLI's
`--full-depth` behavior. Full-depth search also creates a future exposure: any unrelated
`SKILL.md` later added under fixtures, examples, or tests could become installable. The
implementation should therefore add a release/CI probe that asserts the exact discovered
set is these five names. A symlink layout is **NOT TESTED** because it is unnecessary for
this recommendation.

The exact remote command must be run in a networked environment before merge:

```bash
npx -y skills@latest add feci/parley-deck-skill --full-depth --list
```

That owner/repository form is **NOT TESTED** here. The acceptance result should be exactly
five named skills, not merely a count.

## 3. Compatibility matrix

| Existing path | Evidence for the recommended no-restructure approach | Position |
| --- | --- | --- |
| Universal `skills` CLI | Local 1.5.20 discovery found five with `--full-depth`; a scratch-only project install and JSON listing contained all five. Remote owner/repo form: **NOT TESTED**. | Supported only with the required flag. |
| Package's Node/npm installer | Local `install` and `doctor` against a scratch generic destination both returned `ok: true` for all five; full test suite passed. Download through `npx parley-deck-skill@latest`: **NOT TESTED**. | No code or path change required. |
| npm package payload | `npm pack --dry-run --json` showed all four add-on `SKILL.md` files under `addons/` and no `skills/` tree. | No duplicate skill tree is introduced. |
| Homebrew | Actual `brew install` / upgrade: **NOT TESTED**. | No repository path changes are proposed, but runtime compatibility remains unverified here. |
| Standalone Windows binary | **NOT TESTED** on Windows. | No repository path changes are proposed. |
| WinGet `Feci.ParleyDeckSkill` | **NOT TESTED**. | No repository path changes are proposed. |
| Codex `$skill-installer` repository path | **NOT TESTED**. | Root `SKILL.md` remains in place. |
| Antigravity `plugin.json` | Scratch project `--target agy` install and `doctor` returned valid for all five. Actual `agy plugin validate` / runtime loading: **NOT TESTED**. | Existing root manifest and install-time `skills/SKILL.md` materialization remain unchanged. |
| Legacy Gemini `gemini-extension.json` | Scratch project `--target gemini` install and `doctor` returned valid for all five. Actual `gemini extensions install` / runtime loading: **NOT TESTED**. | Existing root context file remains unchanged. |

The package's own installer/doctor evidence:

```text
$ node bin/parley-deck-skill.js install --target generic --dest /tmp/skills-probe-codex-1.FtB7Od/installer-check-2/parley-deck --json | jq '{ok, command, skills: [.actions[].skills[] | {skill, action}]}'
{
  "ok": true,
  "command": "install",
  "skills": [
    {
      "skill": "parley-deck",
      "action": "installed"
    },
    {
      "skill": "parley-design",
      "action": "installed"
    },
    {
      "skill": "parley-design-check",
      "action": "installed"
    },
    {
      "skill": "parley-tracker",
      "action": "installed"
    },
    {
      "skill": "parley-worktrees",
      "action": "installed"
    }
  ]
}
$ node bin/parley-deck-skill.js doctor --target generic --dest /tmp/skills-probe-codex-1.FtB7Od/installer-check-2/parley-deck --json | jq '{ok, command, targets: [.targets[] | {target, status, skills: [.skills[] | {skill, status, missing}]}]}'
{
  "ok": true,
  "command": "doctor",
  "targets": [
    {
      "target": "generic",
      "status": "valid",
      "skills": [
        {
          "skill": "parley-deck",
          "status": "valid",
          "missing": []
        },
        {
          "skill": "parley-design",
          "status": "valid",
          "missing": []
        },
        {
          "skill": "parley-design-check",
          "status": "valid",
          "missing": []
        },
        {
          "skill": "parley-tracker",
          "status": "valid",
          "missing": []
        },
        {
          "skill": "parley-worktrees",
          "status": "valid",
          "missing": []
        }
      ]
    }
  ]
}
```

The test-suite summary:

```text
$ npm test
ℹ tests 247
ℹ suites 0
ℹ pass 247
ℹ fail 0
ℹ cancelled 0
ℹ skipped 0
ℹ todo 0
ℹ duration_ms 2510.085584
```

The package dry-run relevant subset:

```text
$ npm_config_cache=/tmp/skills-probe-codex-1.FtB7Od/npm-cache npm pack --dry-run --json | jq '.[0] | {filename, entryCount, addonSkills: [.files[].path | select(test("^addons/[^/]+/SKILL\\.md$"))], skillsTreeEntries: [.files[].path | select(startswith("skills/"))]}'
{
  "filename": "parley-deck-skill-1.5.0.tgz",
  "entryCount": 153,
  "addonSkills": [
    "addons/parley-design-check/SKILL.md",
    "addons/parley-design/SKILL.md",
    "addons/parley-tracker/SKILL.md",
    "addons/parley-worktrees/SKILL.md"
  ],
  "skillsTreeEntries": []
}
```

The local Antigravity and Gemini payload checks:

```text
$ node bin/parley-deck-skill.js install --scope project --project /tmp/skills-probe-codex-1.FtB7Od/agy-project-2 --target agy --json | jq '{ok, target: .actions[0].target, skills: [.actions[0].skills[] | .skill]}'
{
  "ok": true,
  "target": "agy",
  "skills": [
    "parley-deck",
    "parley-design",
    "parley-design-check",
    "parley-tracker",
    "parley-worktrees"
  ]
}
$ node bin/parley-deck-skill.js doctor --scope project --project /tmp/skills-probe-codex-1.FtB7Od/agy-project-2 --target agy --json | jq '{ok, target: .targets[0].target, status: .targets[0].status, missing: [.targets[0].skills[].missing]}'
{
  "ok": true,
  "target": "agy",
  "status": "valid",
  "missing": [
    [],
    [],
    [],
    [],
    []
  ]
}
$ node bin/parley-deck-skill.js install --scope project --project /tmp/skills-probe-codex-1.FtB7Od/gemini-project-2 --target gemini --json | jq '{ok, target: .actions[0].target, skills: [.actions[0].skills[] | .skill]}'
{
  "ok": true,
  "target": "gemini",
  "skills": [
    "parley-deck",
    "parley-design",
    "parley-design-check",
    "parley-tracker",
    "parley-worktrees"
  ]
}
$ node bin/parley-deck-skill.js doctor --scope project --project /tmp/skills-probe-codex-1.FtB7Od/gemini-project-2 --target gemini --json | jq '{ok, target: .targets[0].target, status: .targets[0].status, missing: [.targets[0].skills[].missing]}'
{
  "ok": true,
  "target": "gemini",
  "status": "valid",
  "missing": [
    [],
    [],
    [],
    [],
    []
  ]
}
```

## 4. README block

This is finished, copy-pasteable Markdown for the beginning of the installation chapter:

````markdown
> [!TIP]
> **Universal installer — install all five Parley Deck skills**
>
> Run this from the project that should receive the skills:
>
> ```bash
> npx -y skills@latest add feci/parley-deck-skill --full-depth --skill "*"
> ```
>
> Choose the target agent or agents when prompted. `--full-depth` is required: without it,
> the installer discovers only the root `parley-deck` skill. The command selects
> `parley-deck`, `parley-design`, `parley-design-check`, `parley-tracker`, and
> `parley-worktrees`.
>
> Restart your agent, then verify the project install:
>
> ```bash
> npx -y skills@latest list
> ```
>
> Prefer Parley Deck's own installer when you also need automatic runtime detection,
> `doctor`, `status`, or `sync-project`:
>
> ```bash
> npx -y parley-deck-skill@latest install --target all
> ```
````

The panel intentionally makes no agent-count claim. It links capability to the command and
the five observed skill names instead.

## 5. Expected user failure mode

The most likely failure is copying the familiar unflagged command and omitting
`--full-depth`. The observed result is a successful-looking install with only
`parley-deck`; there is no error telling the user that four skills were shadowed.

The README pre-empts this in three ways:

1. the required flag is inside the primary copy/paste command;
2. the warning states the exact one-skill consequence of omitting it; and
3. the verification step lets the user compare the installed list with the five explicit
   names in the panel.

## Positions on F1-F6

- **F1 — Co-equal first option.** Put the highlighted universal path first for reach, with
  the native installer in the same panel for detection and lifecycle tooling. Do not replace
  `parley-deck-skill install --target all`.
- **F2 — GitHub alert syntax.** Use `> [!TIP]`. GitHub gives it the requested visual
  emphasis; npm and terminal renderers that do not implement alerts still show a readable
  blockquote with the literal `TIP` marker.
- **F3 — No number.** Do not repeat an upstream agent-count claim in our README. It is not
  needed to explain the benefit, and I did not measure supported agents.
- **F4 — Refuse a misleading feature.** Feature the path only with `--full-depth` and only
  after the exact remote `--list` command returns all five. If a future CLI version again
  installs only the core, remove or demote the panel rather than presenting it as the
  all-five path.
- **F5 — Verification required.** Include `npx -y skills@latest list`. It is the universal
  installer's analogue of `doctor`, even though it does not replace the package-specific
  health checks.
- **F6 — Upstream work is a follow-up, not a prerequisite.** A useful upstream proposal
  would be for explicitly declared marketplace paths to bypass root shadowing, or for the
  default output to warn that nested skills were skipped. No issue or PR was opened, as
  required. Upstream maintainers' response and willingness to change this behavior are
  **NOT TESTED**.
