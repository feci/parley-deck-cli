---
idea: skills-cli-install-path
drafted-by: claude-1
date: 2026-07-29
rounds: 1
participants: [claude-1, codex-1, hermes-1, kimi-1]
status: awaiting-signoffs
---

## Why there is no round 2

All four participants ran the experiments independently and **agreed on every measurement**.
There was no factual disagreement to cross-review, and the one decision that remained was
resolved by the user. The adversarial work that would have gone into round 2 is moved to the
Phase 6 review, where the `NOT TESTED` items in S6 must be discharged against a real build
rather than argued about.

One participant's round-1 conclusion was **wrong and was refuted by the others**: claude-1
wrote *"no cheaper fix exists"* after six experiments. codex-1 and hermes-1 independently
found `--full-depth`, a documented flag that makes the current repository work unchanged.
That is recorded in S2, not quietly dropped.

---

## Measured facts (S1) — reproduced independently by all four participants

Against `skills` 1.5.20, `npx -y skills@latest add <source> --list`:

| Layout | Found |
|---|---:|
| as published — root `SKILL.md`, add-ons in `addons/` | **1** |
| + `skills/` container, root `SKILL.md` kept | **1** |
| + `.claude/skills/` container, root `SKILL.md` kept | **1** |
| + `.claude-plugin/marketplace.json` declaring all five, root kept | **1** |
| `--skill parley-design` against the manifest layout, root kept | **1** (core only) |
| `skills/` container of **symlinks** | **0** (kimi-1) |
| root `SKILL.md` removed, `addons/` otherwise untouched | **4** |
| all five under `skills/<name>/`, root `SKILL.md` moved | **5** |
| **as published, `--full-depth`** | **5** |

**S1a — the cause is the root `SKILL.md`, not the `addons/` directory name.** A `SKILL.md` at
the repository root shadows every nested skill regardless of container or manifest. Remove it
and `addons/` is discovered by recursive fallback without being a documented container.

**S1b — `.claude-plugin/marketplace.json` is a fallback**, consulted only when ordinary
discovery finds nothing. It cannot supplement discovery, only substitute for it. Do not ship
one.

**S1c — symlinks do not work** (kimi-1, twice). A `skills/` container of symlinks discovers
zero skills.

**S1d — verified against the live repository, not just a local copy:**
`npx -y skills@latest add feci/parley-deck-skill --list` → 1 skill;
the same command with `--full-depth` → 5 skills (claude-1, hermes-1).

**S1e — `--full-depth` installs, it does not merely list.** hermes-1 ran a real install into
an isolated `HOME` and got five skills written to `.claude/skills/<name>/`.

## S2 — the refuted claim

claude-1's round-01 stated *"There is exactly one way… no cheaper fix exists"*. This was
**false**, and false because the experiment set was incomplete: claude-1 tested layouts and
never enumerated the CLI's flags. `skills add --help` documents
`--full-depth   Search all subdirectories even when a root SKILL.md exists` — a flag that
exists precisely for this situation. codex-1 and hermes-1 each found it. **Enumerating the
interface would have found in one command what six layout experiments did not.**

## S3 — a second defect, found by hermes-1

Because `SKILL.md` sits at the repository root, the universal installer treats **the whole
repository as the core skill** and copies `bin/`, `lib/`, `test/`, `dist/`, `scripts/`,
`package.json` and `addons/` into the installed `parley-deck` skill directory. The add-on
payload therefore ships twice. This is independent of discovery and is **not** fixed by
`--full-depth`; only relocating the root file fixes it.

## S4 — the decision (user ruling, 2026-07-29)

The user was shown both options and chose **the flag *and* the layout restructure**.

Adopted: **move the whole core payload into `skills/parley-deck/`, and move `addons/*` to
`skills/*`.** After that, plain `npx skills add feci/parley-deck-skill` discovers all five
with no flag, and S3 is fixed as a side effect.

**S4a — move the payload, not just the file (kimi-1, binding).** The CLI copies exactly the
skill's own directory. Moving only `SKILL.md` would ship the core skill **without the bundled
fallback `COOPERATION.md` its own text promises** and without `references/compatibility.json`
— our own `doctor` would then correctly call it malformed. `skills/parley-deck/` must contain
`SKILL.md`, `references/` and `agents/`.

**S4b — a move, never a copy.** `addons/` disappears; `skills/` appears. A duplicated skill
tree would need a drift guard, and this repository already treats an unguarded second copy as
a defect.

## S5 — forks

- **F1 — co-equal, universal path first, neither labelled "recommended".** Reach and
  capability are different axes; the honest answer depends on which agent the reader uses. Our
  installer's distinguishing verbs (`doctor`, `status`, `sync-project`, `--target all`) must be
  visible in the same screenful.
- **F2 — GitHub alert (`> [!TIP]`).** Renders as a panel on GitHub, degrades to a blockquote
  with a bold lead elsewhere. That degradation is acceptable; a fenced block would be louder
  and less readable.
- **F3 — no agent count of our own.** "~70" and "~75" are their claims about their tool. The
  README says "most agents" and links. We do not restate an uncounted number as fact.
- **F4 — if the restructure fails or is reverted, the panel must state plainly that the path
  installs the core skill only, and name the four it does not.** Featuring an install path
  while concealing that it delivers one fifth of the package is precisely the dishonesty
  `parley-design` exists to forbid — in the README of the package that ships it.
- **F5 — yes, one verification line** (`npx skills list`), matching the `doctor` line that
  follows our own install block.
- **F6 — nothing upstream during this idea.** No PR, issue, or contact. Recorded as a
  follow-up: the shadow rule arguably deserves a report now that we have characterised it.

## S6 — what must be proved before this ships, not asserted

Every one of these is currently **NOT TESTED** and is a review gate, not a checkbox:

1. `npx parley-deck-skill install --target all` still installs correctly, and `status --target
   all --json` reports `valid` for every target.
2. `npm pack` ships `skills/` and no longer ships `addons/` or a root `SKILL.md`
   (`package.json` `files` **and** `pkg.assets` — the standalone binaries embed the latter).
3. The Windows binary builds and installs.
4. `gemini-extension.json`'s `contextFileName` — **kimi-1 named this the sharpest single edge**:
   it is root-relative today, and whether legacy Gemini accepts a nested path is unknown.
5. `plugin.json`'s `"skills": ["skills/SKILL.md"]` — today a path that exists only inside
   installed destinations, fabricated by the installer at `lib/installer.js:993-995`. Under the
   new layout it becomes a real repository path; `agy plugin validate` must be run.
6. Codex `$skill-installer` against a rootless layout.
7. `npx skills add feci/parley-deck-skill --list` → **5**, run against the *published* repo
   after merge, with no flag.

Known installer touch points, read by kimi-1 and claude-1:
`lib/installer.js:115-136` (`REQUIRED_PAYLOAD_FILES` / `PAYLOAD_ENTRIES`), `:137`
(`ADDONS_DIR`), `:454` (`packagedProtocolPath`), `:572` (core hash), `:744-749`
(`validatePayload`), `:751` (add-on discovery), `:993-995` (Antigravity staging), `:1091-1110`
(per-target required files); `package.json:31-43` (`files`) and `:44-55` (`pkg.assets`).

**The installer change and the file move must land in the same release.** An old binary against
a new layout fails validation.

## Follow-ups

1. Report the shadow-rule behaviour upstream (needs user approval).
2. `packaging/winget/README.md` is stale (carried from `readme-skill-catalogue` C10).
3. Whether `doctor` should detect installs made by a foreign installer — deliberately unanswered.

---

## Signoffs
