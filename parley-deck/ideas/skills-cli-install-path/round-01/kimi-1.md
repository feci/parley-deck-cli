---
idea: skills-cli-install-path
round: 01
agent: kimi-1
date: 2026-07-29
---

# kimi-1 — round 01 — lens: risk & failure modes

**Environment.** macOS (Darwin), node v26.5.0, `skills` CLI **1.5.20** (`npm view skills version`, 2026-07-29), our package at 1.5.0 (`package.json`). All experiments ran in `/tmp/skills-probe-kimi-1/` on `rsync` copies of the target repo at branch `readme-skill-catalogue` (`.git`, `node_modules`, `dist` excluded). No real agent directory was ever targeted: every run without `--list` used `HOME=/tmp/skills-probe-kimi-1/fakehome` and a scratch working directory, and `-g`/`--global` was never passed. The pre-existing `~/.claude/skills/parley-*` on this machine predate this session (their marker says `installedAt: 2026-07-29T04:13:43.779Z`, `source: npm:parley-deck-skill@1.5.0`) and were not touched. ANSI cursor-control sequences are stripped from the outputs below for readability; nothing else is edited. Where an output is shortened, the elision is marked.

**One-line verdict.** The operative cause is **(b)**: the root `SKILL.md` does not merely "shadow" — it triggers a hard *early return* in the CLI's `discoverSkills()`, so the container list, the plugin manifests, and the recursive fallback are never consulted. Cause (a) (`addons/` is not a searched container) is also true but is **not** what hides the add-ons: with the root file removed, the recursive fallback finds all four in place.

---

## 1. Evidence

### 1.0 Baseline — the brief's measured fact reproduces

```
$ cd /tmp/skills-probe-kimi-1 && rsync -a --exclude .git --exclude node_modules --exclude dist \
    "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill/" baseline/
$ DO_NOT_TRACK=1 npx -y skills@latest add ./baseline --list

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /private/tmp/skills-probe-kimi-1/baseline
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

One skill, as in the brief. (Side observation with risk relevance: the "Agent detected" banner fires even for `--list`; a bare `add` on a real machine picks detected agents without asking.)

### 1.1 REQUIRED TEST — `skills/` container added, root `SKILL.md` kept → still 1

Setup: `cp -R baseline/addons/<name> with-skills-container/skills/<name>` for the four add-ons; root `SKILL.md` untouched. Verified on disk:

```
with-skills-container/skills/parley-design-check/SKILL.md
with-skills-container/skills/parley-design/SKILL.md
with-skills-container/skills/parley-tracker/SKILL.md
with-skills-container/skills/parley-worktrees/SKILL.md
```

```
$ DO_NOT_TRACK=1 npx -y skills@latest add ./with-skills-container --list

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /private/tmp/skills-probe-kimi-1/with-skills-container
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

A documented, searched container holding four valid skills is **completely ignored** while the root `SKILL.md` exists.

### 1.2 REQUIRED TEST — `.claude-plugin/marketplace.json` added, root `SKILL.md` kept → still 1

Exact manifest used (schema per the CLI's own README, fetched 2026-07-29):

```json
{
  "name": "parley-deck-marketplace",
  "metadata": { "pluginRoot": "." },
  "plugins": [
    {
      "name": "parley-deck",
      "source": ".",
      "skills": [
        "./",
        "./addons/parley-design",
        "./addons/parley-design-check",
        "./addons/parley-tracker",
        "./addons/parley-worktrees"
      ]
    }
  ]
}
```

```
$ DO_NOT_TRACK=1 npx -y skills@latest add ./with-marketplace --list

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /private/tmp/skills-probe-kimi-1/with-marketplace
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

Candidate A is dead **as long as the root `SKILL.md` exists**: manifest-declared paths are never consulted (mechanism in 1.8). Whether a manifest is honored *at all* when it is the only discovery channel was not isolated cleanly (the recursive fallback confounds every in-repo variant I could construct) — **NOT TESTED** in isolation, and moot given the early return.

### 1.3 REQUIRED TEST — root `SKILL.md` removed, nothing else changed → 4 found, in place

```
$ rsync -a baseline/ no-root-skill/ && rm no-root-skill/SKILL.md
$ DO_NOT_TRACK=1 npx -y skills@latest add ./no-root-skill --list

│
●   claude-code_2-1-220_agent  Agent detected — installing non-interactively
│
◇  Source: /private/tmp/skills-probe-kimi-1/no-root-skill
│
◇  Local path validated
│
◇  Found 4 skills

│
◇  Available Skills
│
│    parley-design
│
│      Produce a design system with several independent participants and then apply it, without the result reading as machine-made. …
│
│    parley-design-check
│
│      Run the checkable part of the PDS/1.0 design doctrine against files on disk: design artifacts, DTCG token documents, stylesheets and markup. …
│
│    parley-tracker
│
│      Author epics, user stories, and technical subtasks as canonical markdown files that read well for business people, technical people, AND the AI agents that implement them …
│
│    parley-worktrees
│
│      Allocate, name, isolate, merge, and clean up git worktrees so multiple Parley Deck sessions or parallel Phase-5 implementers can work over one repository without collisions. …
│
└  Use --skill <name> to install specific skills
```

(Long description bodies elided, marked with `…`; names, count and order verbatim.) The four add-ons are discovered **in `addons/`, unchanged** — via the recursive fallback that runs only when the standard locations yield nothing (1.8). Note for risk: exactly four were found, so today the repo carries no stray `SKILL.md` that the fallback would also hoover up; that stops being guaranteed the moment anyone adds an example.

### 1.4 All five under `skills/<name>/`, root `SKILL.md` moved → 5 found (candidate C works)

Setup: `mv SKILL.md skills/parley-deck/SKILL.md`; the four add-ons copied to `skills/<name>/`; `addons/` removed.

```
$ DO_NOT_TRACK=1 npx -y skills@latest add ./all-under-skills --list
…
◇  Found 5 skills
…
│    parley-deck
│    parley-design
│    parley-design-check
│    parley-tracker
│    parley-worktrees
…
└  Use --skill <name> to install specific skills
```

(Full descriptions omitted for length; the five names and the count are verbatim.)

### 1.5 Symlinked `skills/` entries → 0 found (candidate E-by-symlink is dead)

Setup: root `SKILL.md` removed, `addons/` removed, `skills/<name>` are absolute symlinks to `/tmp/skills-probe-kimi-1/external-store/<name>` (outside the repo, so the recursive fallback cannot substitute the real directories):

```
$ ls -la skills-symlinks/skills
lrwxr-xr-x  parley-design -> /tmp/skills-probe-kimi-1/external-store/parley-design
lrwxr-xr-x  parley-design-check -> /tmp/skills-probe-kimi-1/external-store/parley-design-check
lrwxr-xr-x  parley-tracker -> /tmp/skills-probe-kimi-1/external-store/parley-tracker
lrwxr-xr-x  parley-worktrees -> /tmp/skills-probe-kimi-1/external-store/parley-worktrees

$ DO_NOT_TRACK=1 npx -y skills@latest add ./skills-symlinks --list
…
◇  No skills found
│
└  No valid skills found. Skills require a SKILL.md with name and description.
(exit code 1)
```

Discovery does not follow symlinked skill directories (mechanism in 1.8). Independently, our own installer **refuses symlinks in its payload** (`copyRecursive`, `lib/installer.js:1041-1043`: `Refusing to copy symlink in skill payload`). Symlinks are doubly dead, and would also be a Windows-checkout hazard had they worked.

### 1.6 `--full-depth` on the completely untouched repo → 5 found (today's zero-change floor)

```
$ DO_NOT_TRACK=1 npx -y skills@latest add ./baseline --list --full-depth
…
◇  Found 5 skills
│    parley-deck
│    parley-design
│    parley-design-check
│    parley-tracker
│    parley-worktrees
…
└  Use --skill <name> to install specific skills
```

Full discovery of today's published layout with **no repository change at all** — the flag disables the early return and forces the deep pass (1.8). Caveats: it would also surface any future stray `SKILL.md`; it relies on their flag semantics; and installing (not just listing) via `--full-depth` from the remote repo is **NOT TESTED**.

### 1.7 What an install actually writes, and what our tooling makes of it (isolated HOME)

All four runs below used `HOME=/tmp/skills-probe-kimi-1/fakehome`, cwd inside scratch, no `-g`.

**(a) Today's layout, installed by the skills CLI → the whole repository is copied into the agent dir.**

```
$ cd /tmp/skills-probe-kimi-1/fakeproject2
$ HOME=…/fakehome DO_NOT_TRACK=1 npx -y skills@latest add ../baseline --skill parley-deck -a claude-code -y
…
◇  Installed 1 skill ────────────────╮
│  ✓ parley-deck (copied)            │
│    → ./.claude/skills/parley-deck  │
…
$ du -sh .claude/skills/parley-deck
1.7M	.claude/skills/parley-deck
```

The destination contains `bin/`, `lib/`, `test/`, `packaging/`, `package-lock.json`, `.github/`, `addons/` — the entire repo, because a root skill's "skill directory" is the repo root. Observed method was **copy** in this non-interactive run (their README offers symlink interactively; symlink-mode behavior **NOT TESTED**). A `skills-lock.json` is written at the **project root** — their ledger, recording `source`, `sourceType: "local"` and a `computedHash` per skill; **no version number**.

**(b) Post-C layout, installed by the skills CLI → the core skill arrives without its protocol files.**

```
$ cd /tmp/skills-probe-kimi-1/fakeproject
$ HOME=…/fakehome DO_NOT_TRACK=1 npx -y skills@latest add ../all-under-skills --skill '*' -a claude-code -y
…
◇  Installation Summary ─────────────────╮
│  ./.agents/skills/parley-deck          │
│    copy → Claude Code                  │   (repeated for all five)
…
◇  Installed 5 skills ───────────────────────╮
│  ✓ parley-deck (copied)                    │
│    → ./.claude/skills/parley-deck          │   (repeated for all five, ./.claude/skills/…)
```

Two of their own panels **disagree about the destination** (`./.agents/skills/` vs `./.claude/skills/`); the files on disk are in `./.claude/skills/`. A user skimming the summary would look in the wrong place. Each destination holds only its own skill directory: `parley-deck/` contains **only `SKILL.md`** — no `references/COOPERATION.md`, no `references/compatibility.json`, no `agents/` (in my C-variant those still lived at repo root; under a properly executed C they would move with the core — see §2).

**(c) What our `doctor` reports about these foreign installs.**

Against (b) (post-C layout install):

```
$ HOME=…/fakehome node baseline/bin/parley-deck-skill.js doctor --target claude --scope project --project …/fakeproject --json
{
  "ok": false,
  "command": "doctor",
  "targets": [
    {
      "target": "claude",
      "dest": "…/fakeproject/.claude/skills/parley-deck",
      "detected": true,
      "status": "malformed",
      "marker": null,
      "missing": ["references/COOPERATION.md", "references/compatibility.json", "agents/manifest.yaml"],
      "skills": [
        { "skill": "parley-deck",         "status": "malformed", "marker": null, "missing": ["references/COOPERATION.md", "references/compatibility.json", "agents/manifest.yaml"] },
        { "skill": "parley-design",       "status": "valid",     "marker": null, "missing": [] },
        { "skill": "parley-design-check", "status": "valid",     "marker": null, "missing": [] },
        { "skill": "parley-tracker",      "status": "valid",     "marker": null, "missing": [] },
        { "skill": "parley-worktrees",    "status": "valid",     "marker": null, "missing": [] }
      ]
    }
  ]
}
(exit code 1, re-run without pipe)
```

Against (a) (today's layout install): core reports `"status": "valid", "marker": null` (the whole-repo copy happens to contain our required files) and the four add-ons report `"status": "missing"`; top-level `"ok": false`. So: **our `doctor` exits 1 on every skills-CLI install I produced** — once calling a healthy foreign install "malformed", once calling a partial install "valid" while flagging add-ons "missing" that the user never asked the other tool for. `status` shares this code path (`targetStatus`/`skillUnitStatus`, `lib/installer.js:1061-1089`). `uninstall` refusal on unmarked dirs is the same marker check (`uninstallSkillUnit`, `:940-968`) — **NOT TESTED** directly; the install-side equivalent below was.

**(d) The two installers collide on the same destination.**

```
$ HOME=…/fakehome node baseline/bin/parley-deck-skill.js install --target claude --scope project --project …/fakeproject2 --yes
claude/parley-deck: blocked …/fakeproject2/.claude/skills/parley-deck - Destination exists but was not installed by parley-deck-skill. Re-run with --force to replace it.
claude/parley-design: installed …/fakeproject2/.claude/skills/parley-design
claude/parley-design-check: installed …/fakeproject2/.claude/skills/parley-design-check
claude/parley-tracker: installed …/fakeproject2/.claude/skills/parley-tracker
claude/parley-worktrees: installed …/fakeproject2/.claude/skills/parley-worktrees
EXIT=1
```

A user who took the universal path and later follows our README gets a **failed install on the core skill** until they discover `--force`.

**(e) Their `list` sees our installs — as provenance-less directories.**

```
$ cd …/fakeproject2   # core written by skills CLI, add-ons written by OUR installer
$ HOME=…/fakehome DO_NOT_TRACK=1 npx -y skills@latest list
Project Skills

parley-deck           ./.claude/skills/parley-deck
  Agents: Claude Code  Source: /private/tmp/skills-probe-kimi-1/baseline
parley-design         ./.claude/skills/parley-design
  Agents: Claude Code  Source: local
parley-design-check   ./.claude/skills/parley-design-check
  Agents: Claude Code  Source: local
parley-tracker        ./.claude/skills/parley-tracker
  Agents: Claude Code  Source: local
parley-worktrees      ./.claude/skills/parley-worktrees
  Agents: Claude Code  Source: local
```

`skills list` scans agent directories and labels anything it did not install `Source: local`. Both tools see the filesystem; **neither understands the other's ledger** (`skills-lock.json` at project root, keyed by hash; `.parley-deck-skill-install.json` per destination, keyed by npm version — the real-HOME marker shows `version: "1.5.0"`, `source: "npm:parley-deck-skill@1.5.0"`).

### 1.8 Mechanism (source of `skills@1.5.20`, `dist/cli.mjs` from the npx cache) — explanation, not primary evidence

`discoverSkills()` opens with:

```js
if (await hasSkillMd(searchPath)) {
    let skill = await parseSkillAt(searchPath);
    if (skill) {
      …
      skills.push(skill);
      seenNames.add(skill.name);
      if (!options?.fullDepth) return skills;   // root SKILL.md → early return
    }
}
```

The priority-container walk (`skills/`, `skills/.curated`, agent dirs) and the manifest paths (`getPluginSkillPaths`, which reads `.claude-plugin/marketplace.json` / `.claude-plugin/plugin.json`) sit **after** that return. The recursive fallback runs only when the walk found nothing:

```js
if (skills.length === 0 || options?.fullDepth) {
    const allSkillDirs = await findSkillDirs(searchPath);   // depth ≤ 5, skips node_modules/.git/dist/build/__pycache__
    …
}
```

Symlinks: both the container walk and `findSkillDirs` filter on `readdir(…, { withFileTypes: true })` + `entry.isDirectory()`; a `Dirent` for a symlink fails that test, so symlinked skill dirs are never entered. This matches 1.1, 1.2, 1.3, 1.5, 1.6 exactly.

---

## 2. Recommended layout: **C — executed fully — with D+`--full-depth` as the interim honest state**

Scoreboard from §1: **A is dead** (1.2: manifest ignored behind the root early return). **B is dead** (1.1: a `skills/` container is ignored behind the root early return; as a second copy it was also a drift hazard). **Symlinks are dead** (1.5, twice over). **D-plain is false advertising** (1.0). What remains:

- **D+`--full-depth`** (1.6): works today, zero repo change, but pushes a non-obvious flag onto every reader, hoover-finds any future stray `SKILL.md`, and its install path is untested by me. This is the **floor**, and it is a decent floor.
- **C** (1.4): the only repo-side change that delivers plain `npx skills add feci/parley-deck-skill` → all five. **Verified to discover 5/5.**

My position: adopt **C**, and execute it as a *move of the whole core payload*, not just the one file:

```
skills/parley-deck/{SKILL.md, references/, agents/}
skills/parley-design/…  (today's addons/, moved — one home, no second copy, no drift guard needed)
```

The reason is 1.7(b): the skills CLI copies exactly the skill's own directory. Move only `SKILL.md` and every universal install of the core arrives **without the bundled fallback `COOPERATION.md` the skill's own text promises** ("read the live COOPERATION.md protocol, or the bundled fallback") and without `references/compatibility.json` — our `doctor` then calls it "malformed", correctly. C done halfway manufactures the failure mode this idea is meant to remove.

Costs and risks of C (the pricing is mine; the paths are verified in code):

- Our installer hard-codes the package-root layout: `REQUIRED_PAYLOAD_FILES`/`PAYLOAD_ENTRIES` (`lib/installer.js:115-136`), `validatePayload` (`:744-749`), `ADDONS_DIR = "addons"` (`:137`), core hash (`:572`), `packagedProtocolPath` (`:454`), and the Antigravity staging that fabricates `skills/SKILL.md` inside destinations (`:993-995`). All must be re-pointed at `skills/parley-deck/` in the **same release** as the move; an old binary against a new repo layout fails validation.
- `package.json` `files` (`:31-43`) and `pkg.assets` (`:44-55`) must ship `skills/` and stop shipping `addons/` — Homebrew, WinGet and the standalone binaries embed exactly those assets. Portable build not run: **NOT TESTED**.
- `gemini-extension.json` has `"contextFileName": "SKILL.md"` — root-relative. Under C, `gemini extensions install <repo-url>` breaks unless this becomes the nested path, and whether legacy Gemini accepts a nested `contextFileName` is **NOT TESTED**. This is C's sharpest single edge.
- `plugin.json` already declares `"skills": ["skills/SKILL.md"]` — a path that **does not exist in today's repo** (our installer fabricates it inside installed destinations). Under C it becomes a real file at `skills/parley-deck/SKILL.md`; net improvement, but `agy plugin validate` against the new layout is **NOT TESTED**.
- Codex `$skill-installer` from the repo URL (documented in our README): behavior against a rootless layout **NOT TESTED**.
- The npm-tree constraint is satisfied by a *move*: `addons/` disappears, `skills/` appears, `files` lists it once.

Until C ships, the README may offer the universal path only in its D+`--full-depth` form with the caveats of §4 — never as a bare one-liner.

## 3. Compatibility matrix for layout C

| Existing path | Under C | Basis |
|---|---|---|
| `npx parley-deck-skill install --target …` | Breaks until installer is re-pointed at `skills/parley-deck/` and `skills/` for add-ons; must land in the same PR/release as the move | `lib/installer.js:115-136, 454, 572, 744-749, 993-995` (read) |
| `doctor` / `status` / `uninstall` | Same change set; validation kinds (`:1091-1110`) must mirror the new staging | read |
| Homebrew / WinGet / standalone `.exe` | Break silently if `pkg.assets`/`files` still name `addons/` and root `SKILL.md`; build scripts need the same PR | `package.json:31-55` (read); build itself **NOT TESTED** |
| Antigravity `plugin.json` | Improves: declared path becomes real (`skills/parley-deck/SKILL.md`); `agy plugin validate` on new layout **NOT TESTED** | `plugin.json` (read); `installer.js:993-995` |
| Legacy Gemini (`gemini-extension.json`) | Breaks unless `contextFileName` moves; nested value acceptance **NOT TESTED** | `gemini-extension.json` (read) |
| Codex `$skill-installer` from repo URL | Unknown against rootless repo — **NOT TESTED** | README documents the flow |
| `gemini extensions install <url>` | Same as `contextFileName` row | — |
| npm package contents | No duplicate tree: `addons/` removed, `skills/` added to `files` once | constraint satisfied by move |
| Drift guard | No new guard needed (still one copy of each `SKILL.md`); existing doctrine tests unchanged in content | `test/design-addons.test.js` (skimmed) |
| Universal `skills add` | 5/5 discovery (1.4); core installs with its protocol files (fixes 1.7(b)); no more 1.7 MB whole-repo copy (fixes 1.7(a)) | experiments |

## 4. README block (shipping text — valid only after C lands; pre-C it is false)

Place at the top of `## Install`, before the current "Fastest path". Rendering choice per F2 below.

````markdown
> [!TIP]
> **Any agent, one command — the universal `skills` installer.** If your agent is
> not one of the runtimes our own installer targets, the
> [`skills`](https://github.com/vercel-labs/skills) CLI installs all five Parley
> Deck skills (`parley-deck`, `parley-design`, `parley-design-check`,
> `parley-tracker`, `parley-worktrees`) straight from this repository:
>
> ```bash
> npx -y skills@latest add feci/parley-deck-skill
> ```
>
> It asks which skills and which agents to install. Three things to know first:
>
> - **It installs into the current directory by default.** Run it from the
>   project that should carry the skills, or pass `-g` for your user-level
>   agent directories. Non-interactive: append `--skill '*' -a <agent> -y`.
> - **It keeps its own records.** It tracks installs in a `skills-lock.json`
>   that our `doctor`/`status` do not read — they will report these copies as
>   foreign (unmarked), not broken. Verify this path with
>   `npx -y skills@latest list`, not with `doctor`.
> - **One installer per runtime.** If you later run
>   `npx parley-deck-skill install` over the same agent directory, it stops at
>   the foreign copies and asks for `--force` before replacing them — that is
>   the safety latch working, not a bug.
>
> To update, re-run the same command (the CLI also documents a
> `skills update` command), then restart your agent runtime. Our own
> installer below remains the recommended path for the runtimes it targets —
> it adds detection, `doctor`, `status` and project sync that the universal
> installer does not have.
````

Interim variant, honest for **today's** repo (only if C is deferred): same panel, command becomes `npx -y skills@latest add feci/parley-deck-skill --full-depth`, and the first bullet gains "(without `--full-depth`, only `parley-deck` is discovered today — repository change in progress)". `--full-depth` install (vs `--list`) against the remote repo is **NOT TESTED**; verify before printing.

## 5. The failure mode I most expect, and how the README pre-empts it

**The wrong-place, wrong-scope install that no tool will own.** Concretely, from §1 evidence: a user runs the one-liner from their home directory or a random checkout. The CLI defaults to *project* scope (1.7: it wrote `./.claude/skills/` under the cwd) and copies the skills there — detected agents only, no questions asked beyond the banner. Their other agents see nothing. A `skills-lock.json` appears in an unrelated git repo and gets committed by accident. When something feels off they run our `doctor`, which exits 1 — either "malformed" on a healthy foreign core (1.7(c)) or "add-ons missing" they never requested — then they follow our README's install command and it **blocks on the core** with exit 1 until `--force` (1.7(d)). Now two ledgers track overlapping directories with no shared notion of version: theirs stores a `computedHash`, ours stores npm `1.5.0` (1.7(e)); `skills update` re-pulls the repo, our `--force` re-pulls the npm package, and the user can ping-pong between vintages of skills that call each other "companion add-on … to the parley-deck skill" — with `references/compatibility.json`, the file that exists to negotiate exactly that, *absent* from the universal install (1.7(b), 1.7(c)). Every step is something I executed and quoted above, not a hypothetical.

The §4 panel pre-empts it with three bullets that each map to a measured behavior: say where files land (scope), say who keeps records (two ledgers, `skills list` as the verification), say what happens when the installers meet (the `--force` latch). The failure mode is not that the command fails — it is that it *succeeds* somewhere the user did not intend, and every tool they consult disagrees about what happened.

## 6. Forks

- **F1 — recommended or co-equal?** Co-equal second, not *the* recommended install. Our installer keeps top billing for its 14 native targets because it owns the lifecycle it creates (marker, `doctor`, `status`, `sync-project`, `--force` latch) — the universal path, per §5, creates installs nobody owns. The panel's reach argument is real but belongs to users outside those 14 runtimes. Featuring it *first* is acceptable to me only with the caveat bullets intact; replacing our block with it is not.
- **F2 — how to render "highlighted"?** GitHub alert syntax, `> [!TIP]`. On GitHub it renders as a colored callout; on npmjs.com and in terminal `npm view … readme` it degrades to a plain blockquote with a literal `[!TIP]` line — readable, and the bold lead sentence still carries the message (npmjs alert rendering: **NOT TESTED**; the degradation is the designed fallback). A fenced block is louder but reads as code to copy, which invites pasting the whole panel as a command.
- **F3 — may we claim a number?** No number of our own. The "70 more" figure is the `skills` README's claim (their README, fetched 2026-07-29: "Supports OpenCode, Claude Code, Codex, Cursor, and [70 more]"); "75" is the brief's arithmetic on that claim. My §4 panel names no count at all — the strongest honest sentence is "installs into any agent the `skills` CLI supports", linked. If a number is ever printed, it must be "their README advertises …", attributed, dated.
- **F4 — if only the core installs after our change?** Then the universal path does not get featured as a full install. Two honest options: the D+`--full-depth` interim panel (verified 5/5 discovery, 1.6; install untested), or no panel. A bare one-liner that installs 1 of 5 while implying otherwise is the false sentence the brief forbids — and my §5 shows the support burden lands exactly on users least equipped to debug it.
- **F5 — verification step?** Yes: `npx -y skills@latest list` (verified output, 1.7(e)) is the check for the universal path, printed in the panel, with the explicit warning that our `doctor` will report these installs as foreign — otherwise the first curious user files a bug about `doctor` exit 1 (1.7(c)).
- **F6 — upstream contribution?** Yes, as a **follow-up, not a prerequisite**. Two candidate items, both evidence-backed: (i) the root-`SKILL.md` early return skips containers *and* manifests, which contradicts their own README's presentation of both as additive discovery channels (1.1, 1.2, 1.8); (ii) discovery ignores symlinked skill directories (1.5). Neither blocks us — C works repo-side today and `--full-depth` works consumer-side today. Per the non-goals, no upstream contact during this idea without user approval.

## Appendix — experiment safety protocol

All probing happened in `/tmp/skills-probe-kimi-1/`. The only runs without `--list` used `HOME=/tmp/skills-probe-kimi-1/fakehome`, a scratch cwd, explicit `-a claude-code`, and never `-g`. Post-run check: `ls -la ~/.claude/skills | grep -i parley` shows only the pre-existing install carrying our marker (`installedAt` 04:13Z, before this session); `~/.agents/skills` has no parley entries. The real repositories were not modified and nothing was committed anywhere.
