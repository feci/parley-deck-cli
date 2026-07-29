---
idea: skills-cli-install-path
round: 01
agent: claude-1
date: 2026-07-29
lens: layout & compatibility
---

## Headline

The shadow rule is real, it is absolute, and **no manifest, container, or flag defeats it.**
There is exactly one way to expose our add-ons through `npx skills`: **the repository root
must not contain a `SKILL.md`.** Everything else in this file follows from that.

---

## 1. Evidence

Eight runs, all against a `--depth 1` clone of the real package at
`/tmp/skills-probe-claude1/base`, restructured per experiment. `--list` only lists; nothing
was installed, `-g/--global` was never passed, and the real repository was never touched.

Tool version resolved by `npx -y skills@latest` on 2026-07-29.

| # | Layout | `Found` |
|---:|---|---:|
| 0 | as published: root `SKILL.md`, add-ons in `addons/` | **1** |
| 1 | + `skills/` container holding the 4 add-ons; root `SKILL.md` kept | **1** |
| 3 | + `.claude-plugin/marketplace.json` declaring all 5; root `SKILL.md` kept | **1** |
| 5 | + `.claude/skills/` container holding the 4 add-ons; root `SKILL.md` kept | **1** |
| 6 | exp3 layout, invoked as `--skill parley-design` | **1** (`parley-deck` only) |
| 2 | `skills/` container holding **all five**; root `SKILL.md` **removed** | **5** |
| 4 | `marketplace.json` declaring all 5; root `SKILL.md` **removed** | **4** |
| 7 | core at `skills/parley-deck/`, add-ons left in `addons/`, manifest declares all 5 | **1** |

Verbatim, the two that matter most:

```text
$ npx -y skills@latest add ./base --list
◇  Local path validated
◇  Found 1 skill
     parley-deck
```

```text
$ npx -y skills@latest add ./exp2 --list
◇  Found 5 skills
     parley-deck
     parley-design
     parley-design-check
     parley-tracker
     parley-worktrees
```

### What the experiments establish

**The brief asked us to distinguish two hypotheses.** (a) `addons/` is simply not a searched
container, versus (b) the root `SKILL.md` shadows everything nested below it.

**Answer: (b), and (a) is true as well but irrelevant.** Experiment 1 puts the add-ons in
`skills/` — a documented level-1 container — and still finds one. Experiment 5 does the same
with `.claude/skills/` and still finds one. Container-ness is not the binding constraint; the
root file is.

**Experiment 6 kills the cheap option.** `--skill parley-design` against a layout that
declares `parley-design` in a manifest returns `parley-deck`. The shadowed skills are not
merely unlisted — they are **unreachable by any flag I could find**. So option D in the brief
("document five commands") is not a floor. It is not an option at all.

**Experiments 3, 4 and 7 together give a model of the manifest.** Compare:

- exp4: no root `SKILL.md`, no container, manifest declaring five → **4** (the four
  `addons/*` paths resolve; the core's `"./"` entry has no `SKILL.md` to point at).
- exp3: root `SKILL.md` present, same manifest → **1**.
- exp7: core moved to `skills/parley-deck/`, add-ons left in `addons/`, same manifest → **1**.

Ordinary discovery succeeded in exp3 and exp7 and failed in exp4, and the manifest was
honoured **only** in exp4. The model that fits all three: **`.claude-plugin/marketplace.json`
is a fallback, consulted only when normal discovery finds nothing.** A manifest cannot
supplement discovery; it can only substitute for it.

I hold this model at **medium confidence**. It fits eight observations and I did not read the
CLI's source. It predicts that a layout with two skills in `skills/` and three more declared
only in the manifest would find two, not five — **a test I did not run, and one I would like
codex-1 to run**, because my recommendation depends on the model being right.

**Untested, stated as untested:** whether the CLI follows symlinks; whether a `skills/`
container of symlinks into `addons/` is discovered; whether `--copy` changes any of this;
behaviour of `skills update` against our repo; whether the published GitHub repo behaves
identically to a local clone (exp0 reproduced the published result, so probably, but I ran
the published probe only in the layout we already ship).

---

## 2. Recommended layout

**Option B, in its complete form — one container, nothing at root, no duplication.**

```text
parley-deck-skill/
├── skills/
│   ├── parley-deck/          # was: root SKILL.md + references/
│   │   ├── SKILL.md
│   │   └── references/
│   ├── parley-design/        # was: addons/parley-design/
│   ├── parley-design-check/
│   ├── parley-tracker/
│   └── parley-worktrees/
├── bin/  lib/  agents/  packaging/  test/
└── package.json  plugin.json  gemini-extension.json  README.md …
```

`addons/` is **renamed**, not copied. This matters more than it looks: the brief forbids a
second hand-maintained copy of any `SKILL.md` without a drift guard, and this repo already
enforces that rule elsewhere by test. A rename has no second copy, so there is nothing to
guard and nothing that can drift. Any option that duplicates a skill tree to satisfy the
scanner is strictly worse than this, however small the diff looks today.

I recommend **not** shipping `.claude-plugin/marketplace.json`. Under my model it would be
dead weight — never consulted, because discovery will now succeed — and a manifest that is
never read is a file that silently rots. If codex-1 refutes the model, this reverses.

### What it costs

The cost is real and I am not going to soften it. Moving the root `SKILL.md` touches every
path assumption in the package:

- `lib/installer.js`: `SKILL.md` appears at lines 116, 127, 572 (`skillSha256`), 994 (the
  Antigravity target, which *synthesises* `skills/SKILL.md` into a temp dir at install time),
  1094, 1096, 1098, 1100 (per-target required-file lists). `ADDON_REQUIRED_FILE` at 138 and
  the add-on discovery glob at 751 assume `addons/<name>/SKILL.md`.
- `package.json` `files`: currently `["SKILL.md", …, "addons/", …, "references/"]` — all
  three entries change, and **`skills/` must be added or the npm tarball ships nothing**.
- `gemini-extension.json`: `"contextFileName": "SKILL.md"`.
- `plugin.json`: declares `"skills": ["skills/SKILL.md"]`. Note this path does **not** exist
  in the repository today — the installer materialises it at line 994. Under the new layout
  the repository and the manifest agree for the first time, but the exact string still needs
  checking against what Antigravity expects.
- `README.md`: every manual path, and the Codex `$skill-installer` instruction.
- The test suite (247 tests today), which is also the thing that makes this survivable.

### What it does not cost

**Nothing about the *installed* layout changes.** The installer already copies files into
per-runtime destinations; what moves here is the source layout inside the package. Users of
`npx parley-deck-skill install`, Homebrew, WinGet and the Windows binaries should see no
difference — but that is a claim to be **proved by running each path**, not asserted, and I
would block a release that asserted it.

---

## 3. Compatibility matrix

| Path | Effect | How it is proved |
|---|---|---|
| `npx -y parley-deck-skill@latest install --target all` | none, if installer constants are updated | run it; `status --target all --json` all `valid` |
| `--no-addons` / `--only <name>` | none; add-on discovery glob changes directory only | existing installer tests, updated |
| npm tarball | **breaks unless `files` is updated** | `npm pack` and list the tarball |
| Homebrew formula | none — it wraps the npm package | `brew upgrade` + `parley-deck-skill --version` |
| WinGet / Windows `.exe` | none — packaging bundles `bin/`+`lib/` | build and run `install --target all --force` |
| Antigravity `plugin.json` | needs re-checking; may improve | `agy plugin validate` |
| legacy Gemini extension | `contextFileName` must change | install to `--target gemini`, confirm the file lands |
| Codex `$skill-installer` | README instruction changes | manual |
| `npx skills add feci/parley-deck-skill` | **1 → 5** | `--list` against the published repo after merge |

---

## 4. The README block

Placement: the first thing under `## Install`. Rendered with a GitHub alert so it is visually
separated on GitHub, and degrades to an ordinary blockquote on npmjs.com and in a terminal.

```markdown
## Install

> [!TIP]
> **One command, most agents.** The universal skill installer from `vercel-labs/skills`
> installs all five skills into whichever coding agents you have:
>
> ```bash
> npx -y skills add feci/parley-deck-skill
> ```
>
> It detects your installed agents and asks which to install into; `--agent <name>` picks
> them explicitly and `--list` shows what the repository offers without installing anything.
> Their tool supports far more agents than our own installer knows about — see
> [vercel-labs/skills](https://github.com/vercel-labs/skills) for the current list.

Our own installer covers fourteen runtimes and does things the universal one does not —
runtime detection, `doctor`, `status`, and project-metadata sync:

```bash
npx -y parley-deck-skill@latest install --target all
npx -y parley-deck-skill@latest doctor --target all
```
```

Two deliberate choices in that copy. It says **"most agents"**, not a number — the count is
their claim about their tool, not a measurement of ours (F3). And it names what our own
installer does *better*, immediately, so the reader is choosing between two paths rather than
being sold one.

**This block must not be written until the layout change is merged and verified.** Today the
first sentence would be false.

---

## 5. The failure mode I most expect

**A user runs `npx skills add feci/parley-deck-skill`, gets five skills, and then runs
`parley-deck-skill doctor` — which reports nothing installed.**

The two installers write to different places and neither knows about the other. Ours writes
`.parley-deck-skill-install.json` markers into destinations it manages; the universal one does
not, so `doctor`, `status` and `sync-project` will not see its work. A user who mixes them has
two installations, one invisible to our tooling and possibly a different version.

The README pre-empts it by saying so in the panel — one sentence, not a warning box — and by
not implying that `doctor` verifies every install path. **Whether `doctor` should learn to
detect foreign installs is a genuine question and I am deliberately not answering it here**:
it is installer behaviour, this idea is scoped to documentation plus layout, and I would
rather record it as a follow-up than smuggle a feature in.

---

## Forks

- **F1 — co-equal, universal path first.** Reach and capability are different axes and the
  reader is not choosing a winner, they are choosing a starting point. Put the universal path
  first because it is one line and works for more people; put ours immediately after with its
  distinguishing verbs (`doctor`, `status`, `sync-project`, `--target all`) visible in the
  same screenful. Do not call either "recommended" — that word would be doing the reader's
  thinking for them, and the honest answer depends on which agent they use.
- **F2 — GitHub alert (`> [!TIP]`).** It is the only syntax that renders as a distinct panel
  on GitHub, which is where this README is overwhelmingly read, and its degradation is a
  plain blockquote with a bold first line — still visibly a callout everywhere else. A fenced
  block would be louder and unreadable; bold-only would not separate from the surrounding
  prose at a glance.
- **F3 — no number of our own.** We may point at their list; we may not restate "70" or "75"
  as a fact in our README, because we have not counted it and it changes without us. "Most
  agents", plus a link, is both honest and more durable.
- **F4 — document the limitation, in the panel.** If the layout change is rejected or fails,
  the panel says plainly that this path installs the core skill only and names the four it
  does not install. Featuring a path while concealing that it delivers one fifth of the
  package is exactly the failure `parley-design`'s honesty rule names, and we would be doing
  it in the README of the package that ships that rule.
- **F5 — yes, one verification line.** `npx skills list` after install. Our own block is
  followed by `doctor`; asymmetry here would read as confidence we have not earned.
- **F6 — follow-up, not a prerequisite, and nothing upstream during this idea.** The layout
  change is entirely within our repository and needs nothing from them. There may be a case
  for asking whether a root `SKILL.md` should suppress a `skills/` container — arguably a
  monorepo-hostile default — but the brief forbids contact during this idea and I agree with
  the brief: we should ship our own fix first and file a report from a position of having
  actually characterised the behaviour.

---

## The thing I would flag to the user before implementing

This is not the small change it looked like when the idea opened. It moves every skill file in
the package to make a third-party tool see them.

I still recommend it — the current state is that four of five skills are invisible on a
distribution channel we are about to advertise, and no cheaper fix exists (exp1, exp3, exp5,
exp6 are all the cheaper fixes, and all four fail). But the honest framing is **"restructure
the package for a third-party scanner"**, not "add an install command to the README", and the
user asked for the second thing believing it was the whole job. That difference is worth
surfacing before Phase 5 rather than after.
