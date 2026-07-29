---
idea: skills-cli-install-path
status: final
drafted-by: claude-1
date: 2026-07-29
consensus: consensus.md (S1–S6; user ruling on S4)
participants: [claude-1, codex-1, hermes-1, kimi-1]
implementation-target: parley-deck-skill
---

## What ships

Two things, in one release, in this order:

1. **A layout move** so that plain `npx skills add feci/parley-deck-skill` discovers all five
   skills with no flag — and so the core skill stops shipping the entire repository.
2. **A highlighted panel at the top of the README install chapter** documenting that path.

**The panel may not be written until the move is verified.** Today its first sentence would
be false, and that is the exact failure `parley-design`'s honesty rule names.

## The move

```
skills/parley-deck/          ← SKILL.md + references/ + agents/   (from the repo root)
skills/parley-design/        ← addons/parley-design/
skills/parley-design-check/  ← addons/parley-design-check/
skills/parley-tracker/       ← addons/parley-tracker/
skills/parley-worktrees/     ← addons/parley-worktrees/
```

`addons/` disappears. The repo root keeps `bin/`, `lib/`, `test/`, `packaging/`, `scripts/`,
`package.json`, `plugin.json`, `gemini-extension.json`, `README.md`, `LICENSE`, `NOTICE.md`.

**A move, never a copy (S4b).** No second skill tree, therefore no drift guard needed.

### R1 — the whole payload moves, not just `SKILL.md` (S4a, binding)

The universal CLI copies exactly the skill's own directory. A file-only move would ship the
core skill without the bundled fallback `COOPERATION.md` its own text promises and without
`references/compatibility.json`; our `doctor` would then correctly call it malformed.

### R2 — `plugin.json` and `gemini-extension.json` (hermes-1 and kimi-1, co-signed)

Both signoffs flagged that S4a's list was incomplete. These two are **repo-level manifests,
not skill-internal files**, and `PAYLOAD_ENTRIES` (`lib/installer.js:126-132`) currently
stages them into the core skill destination from `packageRoot`. **They stay at the repo
root**, and the installer's source path for them is updated explicitly rather than left to
follow from the move. The Antigravity validator (`:1098`) requires them in the *destination*,
which is unchanged.

### R3 — installer changes, all in the same commit as the move

| Location | Change |
|---|---|
| `REQUIRED_PAYLOAD_FILES` / `PAYLOAD_ENTRIES` (`:115-136`) | core payload resolves under `skills/parley-deck/`; the two root manifests keep resolving from `packageRoot` |
| `ADDONS_DIR` (`:137`) | `"addons"` → `"skills"`, with the core skill name excluded from add-on discovery |
| `packagedProtocolPath` (`:454`) | `references/COOPERATION.md` → `skills/parley-deck/references/COOPERATION.md` |
| core hash (`:572`) | new `SKILL.md` path |
| `validatePayload` (`:744-749`) | follows `REQUIRED_PAYLOAD_FILES` |
| `discoverAddons` (`:751`) | new root, skip `parley-deck` |
| Antigravity staging (`:993-995`) | new source path for the fabricated `skills/SKILL.md` |
| `validateInstalledPayload` (`:1091-1110`) | **destination** shapes are unchanged — verify, do not edit blindly |
| `package.json` `files` and `pkg.assets` | ship `skills/`, drop `addons/` and the root `SKILL.md` |

**An old binary against a new layout fails validation.** The move, the installer change and
the packaging change are one atomic release.

## The README panel

Placed as the **first thing** under the install chapter, rendered with a GitHub alert.

```markdown
> [!TIP]
> **One command, most agents.** The universal skill installer from
> [`vercel-labs/skills`](https://github.com/vercel-labs/skills) installs all five skills into
> whichever coding agents you have — many more than this package's own installer knows about:
>
> ```bash
> npx -y skills add feci/parley-deck-skill
> npx -y skills list
> ```
>
> It detects your agents and asks which to install into. `--agent <name>` picks them
> explicitly, `--list` shows what the repository offers without installing anything.
```

Then, immediately below and in the same screenful, our own installer with its distinguishing
verbs (`--target all` detection, `doctor`, `status`, `sync-project`).

Binding wording constraints:
- **F3** — no agent count of our own. "most agents" plus the link. Never "70" or "75" as our
  fact.
- **F1** — neither path is labelled "recommended". *(kimi-1 disclosed in signoff that its
  round-01 position gave our installer top billing; it accepted the group wording and asked
  that the difference be visible rather than buried. It is recorded here.)*
- **F5** — the `skills list` verification line is part of the panel, matching the `doctor`
  line that follows our own install block.
- **F4** — if any gate below fails and the move is reverted, the panel MUST say plainly that
  the path installs the core skill only and name the four it does not.

## Gates — every one must pass before release, by running it

1. `install --target all` works; `status --target all --json` reports `valid` for every target.
2. `npm pack` ships `skills/`, and ships no `addons/` and no root `SKILL.md`.
3. The standalone Windows binary builds and installs.
4. `gemini-extension.json` `contextFileName` — **the sharpest single edge (kimi-1)**. It is
   root-relative today. Resolve it and prove the gemini target still installs.
5. `agy plugin validate` against the new layout. `plugin.json` declares
   `"skills": ["skills/SKILL.md"]`, a path the installer fabricates in the destination.
6. `npm test` green, with installer tests updated for the new layout.
7. **G7, strengthened by kimi-1 — an install, not a listing.** `--list` and install diverge:
   kimi-1 measured discovery finding 5 while the installed core arrived payload-less. So:
   run a real `npx skills add feci/parley-deck-skill` from the **merged remote** into an
   isolated `HOME`, and assert the core destination
   - **contains** `references/` and `agents/`,
   - **does not contain** `bin/`, `lib/`, `package.json`,
   - and that all five skills are present **without `--full-depth`**.

   This is the only check that proves S3 fixed rather than asserting it as a side effect.

Codex `$skill-installer` against a rootless layout is **NOT TESTED** and is recorded as a
known unknown, not a gate — no shipped file records what it installs.

## Non-goals

No change to any skill's **content**. No upstream PR, issue, or contact. No new publishing
channel. No claim about agent counts we have not measured.

## Definition of done

`npx skills add feci/parley-deck-skill` discovers and installs all five skills from the
published repository with no flag, proven by G7 · every existing channel verified by running
it · the panel present and true · `npm test` green.
