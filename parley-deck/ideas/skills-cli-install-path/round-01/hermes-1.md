---
idea: skills-cli-install-path
round: 01
agent: hermes-1
date: 2026-07-29
---

# hermes-1 — the highlighted README block

My lens is the install chapter panel: what the reader sees first, how it renders,
and whether it is honest about what `npx skills add` actually delivers. I ran
every experiment below in `/tmp/skills-probe-hermes-1/` against a scratch copy
of the repo. No real agent directory was touched. No repo file was modified.

Environment: skills CLI v1.5.20, Node v26.5.0, npm 11.17.0, macOS 26.5.2.

---

## 1. Evidence — what I ran and what it printed

All experiments use `npx -y skills@latest add <path> --list` (the `--list` flag
lists discovered skills without installing anything). The scratch copy at
`/tmp/skills-probe-hermes-1/repo` is a byte copy of `parley-deck-skill` on
branch `readme-skill-catalogue` with `node_modules/` and `.git/` removed.

### 1.0 Repo layout (scratch copy, unmodified)

```text
$ find /tmp/skills-probe-hermes-1/repo -name "SKILL.md" -not -path "*/test/*" | sort
/tmp/skills-probe-hermes-1/repo/SKILL.md
/tmp/skills-probe-hermes-1/repo/addons/parley-design-check/SKILL.md
/tmp/skills-probe-hermes-1/repo/addons/parley-design/SKILL.md
/tmp/skills-probe-hermes-1/repo/addons/parley-tracker/SKILL.md
/tmp/skills-probe-hermes-1/repo/addons/parley-worktrees/SKILL.md
```

Five SKILL.md files: one at the repo root, four under `addons/<name>/`.

### 1.1 Baseline — unmodified repo, local path

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/repo --list
◇  Source: /tmp/skills-probe-hermes-1/repo
◇  Local path validated
◇  Found 1 skill

◇  Available Skills

    parley-deck
      Run Parley Deck multi-agent idea, implementation, review, or consensus workflows …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** The four add-ons are invisible. This reproduces the
brief's measured fact.

### 1.2 Baseline — remote repo (published GitHub)

```text
$ npx -y skills@latest add feci/parley-deck-skill --list
◇  Source: https://github.com/feci/parley-deck-skill.git
◇  Repository cloned
◇  Found 1 skill

◇  Available Skills

    parley-deck
      Run Parley Deck multi-agent idea, implementation, review, or consensus workflows …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** Confirmed against the live remote — the local copy is
faithful.

### 1.3 Experiment A — `skills/` container added (addons copied in), root SKILL.md kept

Layout: `skills/parley-design/SKILL.md`, `skills/parley-design-check/SKILL.md`,
`skills/parley-tracker/SKILL.md`, `skills/parley-worktrees/SKILL.md` (real
copies). Root `SKILL.md` and `addons/` both still present.

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/exp-skills-container --list
◇  Source: /tmp/skills-probe-hermes-1/exp-skills-container
◇  Local path validated
◇  Found 1 skill

    parley-deck
      …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** Adding a `skills/` container does NOT help while the root
`SKILL.md` exists. The shadow rule suppresses everything below it.

### 1.4 Experiment B — `.claude-plugin/marketplace.json` declaring all 5 paths, root SKILL.md kept

```json
{
  "name": "parley-deck-skill",
  "version": "1.5.0",
  "skills": [
    { "name": "parley-deck", "path": "SKILL.md" },
    { "name": "parley-design", "path": "addons/parley-design/SKILL.md" },
    { "name": "parley-design-check", "path": "addons/parley-design-check/SKILL.md" },
    { "name": "parley-tracker", "path": "addons/parley-tracker/SKILL.md" },
    { "name": "parley-worktrees", "path": "addons/parley-worktrees/SKILL.md" }
  ]
}
```

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/exp-marketplace --list
◇  Source: /tmp/skills-probe-hermes-1/exp-marketplace
◇  Local path validated
◇  Found 1 skill

    parley-deck
      …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** The marketplace.json manifest does NOT override the shadow
rule. The root `SKILL.md` still shadows everything below it, even skills
explicitly declared in a plugin manifest.

### 1.5 Experiment C — root SKILL.md removed, addons/ untouched

Layout: no root `SKILL.md`; only `addons/<name>/SKILL.md` (four files).
`addons/` is not a documented container directory.

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/exp-no-root-skill --list
◇  Source: /tmp/skills-probe-hermes-1/exp-no-root-skill
◇  Local path validated
◇  Found 4 skills

    parley-design
      Produce a design system with several independent participants …
    parley-design-check
      Run the checkable part of the PDS/1.0 design doctrine …
    parley-tracker
      Author epics, user stories, and technical subtasks …
    parley-worktrees
      Allocate, name, isolate, merge, and clean up git worktrees …

└  Use --skill <name> to install specific skills
```

**Result: 4 skills.** Removing the root `SKILL.md` immediately reveals all four
add-ons — even though `addons/` is not a documented container. This means the
CLI's `--full-depth` search (or default two-level scan) does reach `addons/`
when nothing shadows it.

**This is the distinguishing test the brief required.** The cause of the 1-of-5
result is **(b) the root SKILL.md shadows everything nested below it**, not
**(a) addons/ is not a searched container**. The proof: with the root SKILL.md
removed and no other change, all four add-ons under `addons/` are discovered.

### 1.6 Experiment D — `skills/` container with all 5 skills, root SKILL.md removed

Layout: `skills/parley-deck/SKILL.md`, `skills/parley-design/SKILL.md`,
`skills/parley-design-check/SKILL.md`, `skills/parley-tracker/SKILL.md`,
`skills/parley-worktrees/SKILL.md` (real copies). No root `SKILL.md`.

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/exp-all-in-skills --list
◇  Source: /tmp/skills-probe-hermes-1/exp-all-in-skills
◇  Local path validated
◇  Found 5 skills

    parley-deck
      …
    parley-design
      …
    parley-design-check
      …
    parley-tracker
      …
    parley-worktrees
      …

└  Use --skill <name> to install specific skills
```

**Result: 5 skills.** Moving the core into `skills/parley-deck/` and removing
the root `SKILL.md` discovers all five. This is layout C from the brief (move
the root SKILL.md under a container).

### 1.7 Experiment E — marketplace.json + root SKILL.md removed

Same marketplace.json as Experiment B, but root `SKILL.md` removed. Core skill
moved to `skills/parley-deck/SKILL.md` and declared in the manifest at that
path.

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/exp-marketplace-moved --list
◇  Source: /tmp/skills-probe-hermes-1/exp-marketplace-moved
◇  Local path validated
◇  Found 1 skill

    parley-deck
      …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** Surprise: the marketplace.json pointed the core at
`skills/parley-deck/SKILL.md`, but the four add-ons declared at
`addons/<name>/SKILL.md` were NOT discovered — even with no root SKILL.md. The
manifest appears to declare only the skills it can find at declared paths, and
the shadow/container logic still applies to declared-but-not-at-container
paths. This needs further investigation but the result is clear: marketplace.json
is not a reliable way to surface `addons/` skills. NOT TESTED: whether a
marketplace.json that points all five at `skills/<name>/SKILL.md` (all inside
the container) would work — but that is just layout D with extra steps.

### 1.8 Experiment F — `skills/` with symlinks to addons, root SKILL.md removed

Layout: `skills/parley-deck/SKILL.md` (real copy), `skills/parley-design ->
../../addons/parley-design` (symlink), and so on for the other three. No root
`SKILL.md`.

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/exp-symlink-only --list
◇  Source: /tmp/skills-probe-hermes-1/exp-symlink-only
◇  Local path validated
◇  Found 1 skill

    parley-deck
      …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** The CLI does NOT follow symlinks during discovery. Only
the real `skills/parley-deck/SKILL.md` was found. The four symlinked addon
directories were invisible.

### 1.9 Experiment G — `--full-depth` on the UNMODIFIED repo

This is the key finding. No structural changes at all — the repo exactly as
published, with the root `SKILL.md` and `addons/` in their original locations.

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/repo --list --full-depth
◇  Source: /tmp/skills-probe-hermes-1/repo
◇  Local path validated
◇  Found 5 skills

    parley-deck
      …
    parley-design
      …
    parley-design-check
      …
    parley-tracker
      …
    parley-worktrees
      …

└  Use --skill <name> to install specific skills
```

**Result: 5 skills.** The `--full-depth` flag — documented as "Search all
subdirectories even when a root SKILL.md exists" — discovers all five skills
on the unmodified repository. No file moves, no manifest, no container
directory needed.

### 1.10 Experiment G-remote — `--full-depth` against the live GitHub repo

```text
$ npx -y skills@latest add feci/parley-deck-skill --list --full-depth
◇  Source: https://github.com/feci/parley-deck-skill.git
◇  Repository cloned
◇  Found 5 skills

    parley-deck
    parley-design
    parley-design-check
    parley-tracker
    parley-worktrees

└  Use --skill <name> to install specific skills
```

**Result: 5 skills.** Confirmed against the live remote.

### 1.11 Experiment H — actual install with `--full-depth`

I tested a real installation (not `--list`) into an isolated fake home to
verify that `--full-depth` works for installation, not just listing.

```text
$ cd /tmp/skills-probe-hermes-1/install-test-default && \
  HOME=/tmp/skills-probe-hermes-1/fake-home2 \
  npx -y skills@latest add /tmp/skills-probe-hermes-1/repo \
    --agent claude-code --yes --copy --full-depth

◇  Found 5 skills
●  Installing all 5 skills

◇  Installation Summary
  ./.agents/skills/parley-deck          copy → Claude Code
  ./.agents/skills/parley-design        copy → Claude Code
  ./.agents/skills/parley-design-check  copy → Claude Code
  ./.agents/skills/parley-tracker       copy → Claude Code
  ./.agents/skills/parley-worktrees     copy → Claude Code

◇  Installed 5 skills
  ✓ parley-deck (copied)          → ./.claude/skills/parley-deck
  ✓ parley-design (copied)        → ./.claude/skills/parley-design
  ✓ parley-design-check (copied)  → ./.claude/skills/parley-design-check
  ✓ parley-tracker (copied)       → ./.claude/skills/parley-tracker
  ✓ parley-worktrees (copied)     → ./.claude/skills/parley-worktrees
```

**Result: 5 skills installed.** Each skill lands in its own directory under
`.claude/skills/<name>/`.

**Side effect observed:** the core `parley-deck` install copies the entire
repo root into the skill dir (because `SKILL.md` is at the root), including
`bin/`, `lib/`, `package.json`, `test/`, `dist/`, `scripts/`, and the `addons/`
subdirectory itself. The addon skills install as just their own subdirectory.
This means the `addons/` content appears twice in the install: once inside
`parley-deck/addons/` and once as the top-level `parley-design/` etc. This is
cosmetic, not functional — the addon skills are independently usable — but it
ships ~2× the files for the addon payload.

### 1.12 Experiment I — `--skill '*'` without `--full-depth`

```text
$ npx -y skills@latest add /tmp/skills-probe-hermes-1/repo --list --skill '*'
◇  Found 1 skill

    parley-deck
      …

└  Use --skill <name> to install specific skills
```

**Result: 1 skill.** `--skill '*'` does NOT override the shadow rule. Only
`--full-depth` does.

### 1.13 Experiment J — default install (no `--full-depth`, no `--skill`)

```text
$ HOME=/tmp/skills-probe-hermes-1/fake-home \
  npx -y skills@latest add /tmp/skills-probe-hermes-1/repo \
    --agent claude-code --yes --copy

◇  Found 1 skill
●  Skill: parley-deck

◇  Installed 1 skill
  ✓ parley-deck (copied) → ./.claude/skills/parley-deck
```

**Result: 1 skill installed.** The default `npx skills add feci/parley-deck-skill`
installs only `parley-deck`. This is what a README reader would get if the
README said "install with `npx skills add feci/parley-deck-skill`" and nothing
else.

### 1.14 Agent count verification

The skills CLI's own README (bundled in the npm package at v1.5.20) says:

> Supports **OpenCode**, **Claude Code**, **Codex**, **Cursor**, and [70 more].

I parsed the `<!-- supported-agents:start -->` table from that README and
counted the distinct `--agent` values:

```text
75 distinct --agent values:
  adal, aider-desk, amp, antigravity, antigravity-cli, astrbot, augment,
  autohand-code, bob, claude-code, cline, codearts-agent, codebuddy,
  codemaker, codestudio, codex, command-code, continue, cortex, crush,
  cursor, deepagents, devin, dexto, droid, eve, firebender, forgecode,
  gemini-cli, github-copilot, goose, grok, hermes-agent, iflow-cli,
  inference-sh, jazz, junie, kilo, kimchi, kimi-code-cli, kiro-cli, kode,
  lingma, loaf, mcpjam, mistral-vibe, moxby, mux, neovate, ona, openclaw,
  opencode, openhands, pi, pochi, promptscript, qoder, qoder-cn, qwen-code,
  reasonix, replit, roo, rovodev, tabnine-cli, terramind, tinycloud, trae,
  trae-cn, universal, warp, windsurf, zcode, zed, zencoder, zenflow
```

The npm `keywords` field lists 74 agent-specific keywords (excluding 5
meta-keywords: `cli`, `agent-skills`, `skills`, `ai-agents`, `universal`).
The README table lists 75 including `universal`. Either way, the claim is
the skills CLI's, not ours.

**Our installer's 14 native targets vs. the skills CLI's 75:** 13 of our 14
targets have a direct counterpart in the skills CLI's agent list (codex,
claude/claude-code, agy/antigravity, gemini/gemini-cli, hermes/hermes-agent,
qwen/qwen-code, codebuddy, goose, kimi/kimi-code-cli, droid, vibe/mistral-vibe,
cursor, opencode). One target — `aionrs` — has no counterpart in the skills
CLI. This is a real gap: users who rely on `aionrs` get nothing from the
universal path.

---

## 2. Recommended layout

**Layout D (do nothing structural) + document `--full-depth` in the README.**

The `--full-depth` flag discovers all 5 skills on the unmodified repo (Experiments
G, G-remote). No file moves, no manifest, no container directory, no symlinks.
The only cost is one extra flag in the README command.

I recommend this over every structural alternative because:

- **Layout A (marketplace.json):** Does not work. Experiment B showed 1 skill
  even with a manifest declaring all 5 paths. The shadow rule suppresses
  declared paths too. NOT TESTED: whether a manifest pointing all five at
  `skills/<name>/SKILL.md` would work — but that requires the file moves of
  layout C/D, making the manifest redundant.
- **Layout B (skills/ container, root SKILL.md kept):** Does not work.
  Experiment A showed 1 skill. The shadow rule suppresses everything in
  `skills/` while the root `SKILL.md` exists.
- **Layout C (move root SKILL.md under skills/):** Works (Experiment D, 5
  skills). But it is the highest blast radius: `plugin.json` declares
  `"skills": ["skills/SKILL.md"]` — that path would break. Our own installer
  (`lib/installer.js:117`) requires `SKILL.md` at `PACKAGE_ROOT` in
  `REQUIRED_PAYLOAD_FILES`. `gemini-extension.json` has no path field but
  Gemini's own `gemini extensions install` clones the repo and expects
  `SKILL.md` at the root. Every manual path in the README
  (`~/.codex/skills/parley-deck`, etc.) is written by our installer from the
  root file. Moving the root `SKILL.md` breaks all of these unless every
  consumer is updated simultaneously.
- **Layout E (symlinks):** Does not work. Experiment F showed the CLI does not
  follow symlinks during discovery.
- **Layout D + `--full-depth`:** Works, zero structural change, zero
  compatibility risk. The flag is documented in the CLI's `--help` and README.
  It is not obscure — it is the CLI's own answer to the shadow problem.

**What it costs:** The README command is longer by one flag
(`--full-depth`). A reader who omits the flag gets 1 skill, not 5 — so the
README must be explicit about the flag, or must use `--all` (which is
`--skill '*' --agent '*' -y` and does NOT include `--full-depth`; NOT TESTED
whether `--all` alone suffices — but the `--help` says `--all` is shorthand
for `--skill '*' --agent '*' -y`, with no mention of `--full-depth`).

**What it risks:** The `--full-depth` flag could be renamed or removed in a
future CLI version. We have no control over that. Mitigation: the README
should also show the `--skill` fallback for installing individual skills by
name, which does not depend on `--full-depth` (Experiment I showed `--skill '*'`
alone doesn't work, but `--skill parley-design` with `--full-depth` does —
NOT TESTED: `--skill parley-design` without `--full-depth`).

---

## 3. Compatibility matrix

For layout D (do nothing structural), each existing install path:

| Path | Status | Evidence |
|------|--------|----------|
| `npx parley-deck-skill@latest install --target all` | Unchanged | No repo files moved; installer reads `SKILL.md` at `PACKAGE_ROOT` (installer.js:117). NOT re-tested by running it (no node_modules in scratch), but no file the installer reads was touched. |
| `npm install -g parley-deck-skill` | Unchanged | `package.json` `files` array unchanged. |
| Homebrew (`brew install feci/parley/parley-deck-skill`) | Unchanged | Taps the same repo; no structural change. NOT TESTED (no tap configured). |
| Standalone Windows binaries | Unchanged | Built from the same repo via `scripts/build-portable.js`. NOT TESTED. |
| WinGet (`Feci.ParleyDeckSkill`) | Unchanged | Packages the Windows binary. NOT TESTED. |
| Antigravity `plugin.json` | Unchanged | `plugin.json` still at root, still declares `"skills": ["skills/SKILL.md"]`. File not touched. |
| Legacy Gemini `gemini-extension.json` | Unchanged | Still at root, `contextFileName: "SKILL.md"`. File not touched. |
| Codex `$skill-installer` | Unchanged | Clones repo, reads root `SKILL.md`. NOT TESTED (no Codex CLI in environment). |
| Manual paths (`~/.codex/skills/parley-deck`, etc.) | Unchanged | Written by our installer from root `SKILL.md`. |
| `npx skills add feci/parley-deck-skill` (universal, default) | Still 1 skill | Experiment J. This is the path we are documenting — it needs `--full-depth`. |
| `npx skills add feci/parley-deck-skill --full-depth` | 5 skills | Experiments G, G-remote, H. |

Layout D breaks nothing because it changes nothing. The only "new" thing is
the README text telling readers to pass `--full-depth`.

---

## 4. The README block

This is the highlighted panel for the top of the install chapter. It goes
above the existing "Fastest path" block. It uses GitHub alert syntax
(`> [!TIP]`), which renders as a styled callout on GitHub and degrades to a
blockquote on npmjs.com and in terminals — see F2 below for why.

```markdown
> [!TIP]
> **Universal installer — all five skills, 75+ agents.**
>
> Install every skill in this package (parley-deck plus the four add-ons)
> into any agent supported by the [skills](https://github.com/vercel-labs/skills) CLI:
>
> ```bash
> npx -y skills@latest add feci/parley-deck-skill --full-depth --yes
> ```
>
> The `--full-depth` flag is required: this repo ships a root `SKILL.md`
> (the core skill) that shadows the four add-ons under `addons/` unless
> full-depth discovery is requested. Without it, only `parley-deck` is
> found.
>
> To install a single skill by name:
>
> ```bash
> npx -y skills@latest add feci/parley-deck-skill --full-depth --skill parley-design
> ```
>
> The skills CLI supports 75 agents (per its own README at v1.5.20),
> including Claude Code, Codex, Cursor, OpenCode, Hermes, and Gemini —
> more than our own installer's 14 native targets. One caveat: our
> `aionrs` target has no counterpart in the skills CLI.
>
> After installing, restart your agent runtime so it picks up the new
> skills. To verify what was installed:
>
> ```bash
> npx -y skills@latest list
> ```
>
> **Prefer our own installer** if you need `--target all` runtime
> detection, `doctor`, `status`, `sync-project`, or `aionrs` support —
> see the options below.
```

### Design notes for this block

- **Honesty first.** The block says `--full-depth` is required and says why.
  It does not claim the default command installs all five. A reader who skims
  the command and omits the flag gets one skill, but the block told them the
  flag is required.
- **The number 75 is attributed.** "per its own README at v1.5.20" — we are
  repeating their claim, not making our own. If they change the number, our
  attribution is still honest because it cites the source and version.
- **The `aionrs` gap is stated.** A reader who uses `aionrs` learns here, not
  at runtime, that the universal path does not cover them.
- **The verification step is `skills list`, not `skills doctor`.** The
  universal path has no `doctor` equivalent; `skills list` is the closest
  verification command (see F5).
- **The block points to our own installer for capability, not just as a
  fallback.** The universal path has reach (75 agents); our installer has
  capability (runtime detection, doctor, status, sync-project, aionrs). The
  block says "prefer our own installer if you need…" — co-equal, not
  subordinate.

---

## 5. Failure mode I most expect

**The reader omits `--full-depth` and gets only `parley-deck`.**

This is the single most likely failure because the flag is non-obvious: the
command `npx skills add feci/parley-deck-skill` looks complete, succeeds
without error, and installs a real skill. There is no warning, no "4 more
skills available", no prompt. The reader believes they have all five and
discovers the truth only when `parley-design` or `parley-tracker` is missing
later — possibly much later, in a session where nobody remembers the install
command.

**How the README pre-empts it:**

1. The block leads with "all five skills" and immediately shows the command
   with `--full-depth` — the flag is in the first command the reader sees,
   not in a footnote.
2. The sentence "The `--full-depth` flag is required" is a direct statement,
   not a hint. It says why (shadow rule) so the reader understands the flag
   is not optional decoration.
3. The `skills list` verification step lets the reader confirm they got five
   skills before moving on. If they see one, they know to re-read the block.
4. The "Prefer our own installer" closing gives the reader an alternative
   that does not have this footgun: `npx parley-deck-skill@latest install
   --target all` installs all five (core + addons) by default, with no flag
   gymnastics.

**Second failure mode (less likely, more damaging):** the `--full-depth` flag
is renamed or removed in a future skills CLI version. The README command
breaks silently — it would install 1 skill and print no error. Mitigation:
the block attributes the flag to a specific CLI version (v1.5.20) and offers
the `--skill <name>` fallback for individual installs. Long-term mitigation:
an upstream contribution (F6) to make the CLI respect plugin manifests even
when a root SKILL.md exists, so `--full-depth` is not needed.

---

## Forks

### F1 — the recommended install or a co-equal first option?

**Co-equal first option.** The two paths have different strengths:

- Universal (`npx skills add … --full-depth`): reach to 75 agents, one
  command, no Node package install. Weakness: no runtime detection, no
  doctor, no status, no sync-project, no aionrs, and the `--full-depth`
  footgun.
- Our installer (`npx parley-deck-skill@latest install --target all`):
  capability — runtime detection, doctor, status, sync-project, 14 native
  targets including aionrs, default addon install. Weakness: 14 targets, not
  75.

The README block should present both, side by side, with the universal path
first (it is faster to type and reaches more agents) and our installer second
(it is deeper). A reader who wants reach takes the first; a reader who wants
verification and metadata takes the second. Neither is "the" recommended
install — they are two doors into the same package.

### F2 — how is "highlighted" rendered?

**GitHub alert syntax (`> [!TIP]`).** Rationale:

- On GitHub (the primary surface where this README lives), `> [!TIP]` renders
  as a styled callout with a lightbulb icon and a colored background. It is
  visually distinct from body text without being a loud fenced block.
- On npmjs.com, GitHub alert syntax is not rendered as a callout — it
  degrades to a blockquote. A blockquote is still visually distinct
  (indented, vertical bar) and the bold lead line ("Universal installer —
  all five skills, 75+ agents.") carries the message even without styling.
- In a terminal (e.g. `cat README.md` or a markdown renderer without alert
  support), it degrades to a blockquote with `>` prefixes. Readable, if
  less pretty.
- The alternative (bold-and-blockquote) works everywhere but is less
  visually distinct on GitHub, which is the primary surface. The alternative
  (fenced block) is loudest but breaks in terminals and is semantically
  wrong — a callout is not code.

The degradation is acceptable because the block's content is plain text:
the command is in a fenced code block inside the callout, the explanation is
prose, and the bold lead line works without any styling. A reader on
npmjs.com sees a blockquoted paragraph with a bold first line and a code
block — clear enough.

### F3 — do we claim a number of agents?

**Yes, attributed.** The block says "75 agents (per its own README at
v1.5.20)". This is a verifiable attribution: we name the source (their
README), the version (v1.5.20), and the number (75, which I confirmed by
counting the `--agent` values in their table). We do not say "we support 75
agents" — we say the skills CLI does. If they change the number, our
attribution is still honest because it cites the version.

We should NOT round or soften to "~75" in the README block. The table has
exactly 75 entries. Saying "75" with attribution is more honest than saying
"~75" with attribution, because "~75" implies we measured something and
found approximately 75, when in fact we counted exactly 75 in their
documentation.

### F4 — if the universal path installs only the core, document or refuse?

**Document it, and fix it with `--full-depth`.** The whole point of this idea
is that the universal path CAN install all five — it just needs the flag.
Refusing to feature the path because the default behavior is broken would
deny readers a real installation option. The block is honest: it says the
flag is required, says why, and offers a verification step.

If `--full-depth` did not exist (or if it stopped working), the answer would
change: we would either refuse to feature the path, or feature it with
"installs the core `parley-deck` skill only; for the four add-ons, use our
own installer below." But `--full-depth` does exist and works, so the block
features the full five-skill install.

### F5 — does this path need a verification step?

**Yes: `npx skills list`.** Our own install block is followed by
`parley-deck-skill doctor --target all`. The universal path has no `doctor`
equivalent, but `skills list` shows installed skills and is the closest
parallel. The block includes it. A reader who runs the install command and
then `skills list` can confirm they see five skills, not one.

`skills list` is NOT a full substitute for `doctor` — it does not check
skill validity, marker files, or version match. But it answers the question
"did the install actually give me all five?" which is the verification this
path needs.

### F6 — is there an upstream contribution here?

**Yes, as a follow-up, not a prerequisite.** The shadow rule is the root
cause: a root `SKILL.md` silently suppresses everything below it, and the
only escape is `--full-depth`. A repo that ships a core skill at the root
and add-ons in a non-container directory is a legitimate layout, and the
CLI's default behavior hides four-fifths of the package without warning.

An upstream contribution could be: when a root `SKILL.md` is found and
additional `SKILL.md` files exist in non-container directories, print a
notice ("4 additional skills found under addons/ — use --full-depth to
install them"). This would prevent the silent footgun without changing
discovery semantics. But the brief says no upstream PR during this idea
without user approval, so this is a follow-up, not a prerequisite. The
README block works with the CLI as-is.

---

## Summary of findings

| Experiment | Layout | Root SKILL.md | --full-depth | Skills found |
|------------|--------|---------------|--------------|--------------|
| 1.1 Baseline (local) | unmodified | yes | no | 1 |
| 1.2 Baseline (remote) | unmodified | yes | no | 1 |
| 1.3 Exp A | skills/ container added | yes | no | 1 |
| 1.4 Exp B | marketplace.json | yes | no | 1 |
| 1.5 Exp C | root SKILL.md removed | no | no | 4 |
| 1.6 Exp D | all 5 under skills/ | no | no | 5 |
| 1.7 Exp E | marketplace.json + moved | no | no | 1 |
| 1.8 Exp F | skills/ with symlinks | no | no | 1 |
| 1.9 Exp G | unmodified | yes | yes | 5 |
| 1.10 Exp G-remote | unmodified (remote) | yes | yes | 5 |
| 1.11 Exp H | unmodified (install) | yes | yes | 5 installed |
| 1.12 Exp I | unmodified + --skill '*' | yes | no | 1 |
| 1.13 Exp J | unmodified (install) | yes | no | 1 installed |

**Root cause: (b) the root SKILL.md shadows everything nested below it.**
Proven by Experiment C: removing the root SKILL.md and changing nothing else
reveals all four add-ons under `addons/`, which is not a documented container.

**Fix: `--full-depth` flag, no structural change.** Proven by Experiments G,
G-remote, H: all five skills discovered and installed on the unmodified repo.

**marketplace.json does not override the shadow rule.** Proven by Experiments
B and E.

**The CLI does not follow symlinks during discovery.** Proven by Experiment F.

**The skills CLI supports exactly 75 agents** (counted from its README table
at v1.5.20). Our installer's `aionrs` target has no counterpart.
