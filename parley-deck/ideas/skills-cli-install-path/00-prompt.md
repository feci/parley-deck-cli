---
idea: skills-cli-install-path
author: user
created: 2026-07-29
track: standard
participants: [claude-1, codex-1, hermes-1, kimi-1]
target-repo: parley-deck-skill
status: round-01
depends-on: readme-skill-catalogue   # that idea defines the README structure this one writes into
roles:
  claude-1: layout & compatibility (what changes in the repo, and what it breaks)
  codex-1: empirical verification (prove the discovery behaviour; do not reason about it)
  hermes-1: the highlighted README block (what the reader sees first in the install chapter)
  kimi-1: risk & failure modes (what goes wrong for a user who takes this path)
---

## Problem / idea

The user's ask, verbatim (translated): *"also start, via /parley-deck, a new installation
method for the skill; put it at the start of the installation chapter and highlight it —
look at https://github.com/vercel-labs/skills, it is a rather elegant way to install a skill
into 70 agents."*

`vercel-labs/skills` is a universal skill installer. `npx skills add <owner>/<repo>` clones a
repository, discovers the `SKILL.md` files in it, and installs the selected ones into any of
~75 supported coding agents. It would give this package reach far beyond the fourteen
runtimes our own installer knows about, at the cost of one line in the README.

**One line, if it worked. It does not work yet.**

## The measured fact this idea exists to fix

Run against our published repository, today:

```text
$ npx -y skills@latest add feci/parley-deck-skill --list
◇  Repository cloned
◇  Found 1 skill

     parley-deck
       Run Parley Deck multi-agent idea, implementation, review, or consensus workflows …

└  Use --skill <name> to install specific skills
```

**One skill of five.** `parley-design`, `parley-design-check`, `parley-tracker` and
`parley-worktrees` are all invisible to that installer. A reader who follows a README line
saying "install all five with `npx skills add feci/parley-deck-skill`" gets one, and is told
nothing about the other four.

So this idea is not "add a command to the README". It is, in order:

1. **Determine empirically why four skills are invisible.**
2. **Change the repository so all five are discovered** — without breaking the installer,
   the npm package, Homebrew, WinGet, the Codex `$skill-installer` path, `plugin.json`
   (Antigravity) or `gemini-extension.json` (legacy Gemini).
3. **Only then** write the highlighted README block, and only claiming what step 2 delivers.

Advertising the path before fixing discovery would ship a false sentence in the README of a
package whose own design skill forbids exactly that.

## Ground truth about the `skills` CLI (verify; do not trust this summary)

From its README, as read on 2026-07-29:

- Commands: `add`, `use`, `list`/`ls`, `find`, `remove`/`rm`, `update`, `init`.
- `add` flags include `-g/--global`, `-a/--agent <agents...>`, `-s/--skill <name>`,
  `-l/--list`, `--copy`, `-y/--yes`, `--all`.
- Discovery, level 1: the repository root, `skills/`, `skills/.curated/`,
  `skills/.experimental/`, `skills/.system/`, and 30+ agent-specific directories such as
  `.claude/skills/` and `.agents/skills/`.
- Discovery, level 2: one directory deeper inside those same containers.
- **Shadow rule: "A `SKILL.md` discovered at the shallower level shadows anything nested
  below it."**
- A `SKILL.md` at the repository root is itself a discoverable skill.
- Plugin manifests: skills declared in `.claude-plugin/marketplace.json` or
  `.claude-plugin/plugin.json` are discovered at their declared paths.
- Env vars: `INSTALL_INTERNAL_SKILLS`, `DISABLE_TELEMETRY`, `DO_NOT_TRACK`.
- Claims support for ~75 agents ("**OpenCode**, **Claude Code**, **Codex**, **Cursor**, and
  [70 more]").

**Our repository layout today:** `SKILL.md` at the root; the four add-ons at
`addons/<name>/SKILL.md`. `addons/` is not in the documented container list, and the root
`SKILL.md` sits at level 1.

**Both of those could independently explain the result, and the difference matters**, because
one is fixed by adding a container directory and the other is not fixable while a root
`SKILL.md` exists. **Nobody may assert a cause in round-01 without having run a test that
distinguishes them.** Reasoning from the documentation is not evidence here; the documentation
is a summary of a summary.

## Candidate layouts — argue for one, and say what it costs

- **A. `.claude-plugin/marketplace.json`** declaring all five paths. Additive; no file moves.
  Does the shadow rule still suppress them? Unknown — test it.
- **B. Add a `skills/` container** holding the four add-ons (and possibly the core), leaving
  `addons/` in place or replacing it. Note that a second copy of a skill is a **drift
  hazard**, and this repository already enforces an embedded-default drift guard by test for
  exactly that reason. A symlink is not a copy — does the CLI follow one? Test it.
- **C. Move the root `SKILL.md`** under a container so nothing shadows. Highest reach,
  highest blast radius: our own installer, `plugin.json`, `gemini-extension.json`, the Codex
  `$skill-installer` flow and every documented manual path assume the root file.
- **D. Do nothing structural**; document `--skill <name>` and tell readers to run the command
  five times. Cheap, honest, and bad — but it is the floor every other option must beat.
- **E. Something else.** A `skills/` container of symlinks, a build step that materialises the
  container at publish time, or a change upstream. Propose it if you have it.

## Constraints

- **Never break the existing paths.** npm, Homebrew, the standalone Windows binaries, WinGet
  (`Feci.ParleyDeckSkill`), Antigravity `plugin.json`, legacy Gemini `gemini-extension.json`,
  and `npx parley-deck-skill install --target all` must keep working exactly as they do now.
- **No second hand-maintained copy of any `SKILL.md`** unless a test enforces that the copies
  are identical. This repo's rule, already enforced elsewhere: two copies without a drift
  guard is a defect, not a layout.
- **The npm package must not grow a duplicate skill tree** it then ships twice. Check
  `package.json` `files`.
- **No new runtime dependency.** The installer is dependency-free and stays that way.
- English only in every artifact.

## What each participant writes in round-01

Write `round-01/<your-agent-id>.md` **independently — do not read the others' files first**.

1. **Your evidence.** What you ran, and what it printed. Clone the repo to a scratch
   directory, restructure it locally, and run `npx -y skills@latest add <local-or-fork-path>
   --list` against each candidate layout. **Report commands and output, not conclusions.**
   If you could not run something, say so plainly — an untested claim marked as untested is
   fine; an untested claim presented as a finding is not.
2. **Your recommended layout** (A–E), with what it costs and what it risks.
3. **The compatibility matrix** — for your layout, what happens to each existing install path.
4. **The README block.** Actual shipping markdown for the highlighted panel at the top of the
   install chapter. It must be honest about what that path installs.
5. **The failure mode you most expect** for a real user taking this path, and how the README
   pre-empts it.

## Open forks

- **F1** Does this path become *the* recommended install, or a co-equal first option next to
  `npx -y parley-deck-skill@latest install --target all`? Our own installer does things the
  universal one does not: `--target all` detection, `doctor`, `status`, `sync-project`,
  project metadata. Reach versus capability — pick, and justify.
- **F2** How is "highlighted" rendered in GitHub-flavoured markdown, given the README is read
  on GitHub, on npmjs.com, and in a terminal? GitHub alert syntax (`> [!TIP]`) renders on
  GitHub and degrades to a blockquote elsewhere. Bold-and-blockquote works everywhere. A
  fenced block is loudest and least readable. Choose one and say why the degradation is
  acceptable.
- **F3** Do we claim a number of agents? "~75" is their claim, not our measurement. May we
  repeat someone else's number in our README, and if so, attributed how?
- **F4** If the universal path installs only the core skill even after our change, do we
  document that limitation in the panel, or refuse to feature the path at all?
- **F5** Does this path need a verification step for the user (`skills list`), the way our own
  install block is followed by `doctor`?
- **F6** Is there an upstream contribution here — and if so, is it a prerequisite for this
  idea or a follow-up? Do not open anything upstream during this idea without user approval.

## Non-goals

- No change to any `SKILL.md` **content**, to `COOPERATION.md`, or to what the skills do.
- No new publishing channel beyond what is required to be discovered.
- No upstream PR, issue, or contact during this idea.
- No claim about agent counts, adoption, or reach that we have not measured or attributed.

## Definition of done

`npx skills add feci/parley-deck-skill` discovers every skill this package ships — or the
README states plainly and accurately which ones it does not. The install chapter opens with a
highlighted block for this path. Every existing install path still works, proven by running
them, not by assuming.
