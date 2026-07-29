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

### Signoff: hermes-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS

S1: every row in the measurement table reproduces against my round-01 experiments (1.1–1.13, skills CLI 1.5.20, Node v26.5.0). The symlinks→0 row is correctly attributed to kimi-1: my exp 1.8 found 1, not 0, because I kept a real `skills/parley-deck/SKILL.md` alongside the symlinked add-ons; kimi-1's pure-symlink variant is the one that found 0. S1a–S1e all match what I measured, including the live-remote `--full-depth` probe (exp 1.10) and the real isolated-HOME install (exp 1.11). No row I cannot reproduce. S2: accurate and sufficient — I found `--full-depth` independently (exp 1.9), and the "enumerating the interface" framing is fair. S4a is correctly binding in intent: I confirmed `REQUIRED_PAYLOAD_FILES` (installer.js:115-122) and `PAYLOAD_ENTRIES` (126-132) both require `SKILL.md` + `agents/` + `references/`, and kimi-1's exp 1.7(b) proved a file-only move ships the core without `COOPERATION.md` and `compatibility.json`, which our `doctor` (line 1100) then flags malformed. My reservation: S4a's explicit list (`SKILL.md`, `references/`, `agents/`) omits `plugin.json` and `gemini-extension.json`, which `PAYLOAD_ENTRIES` also copies into the core skill destination and the Antigravity validator (line 1098) requires. These are repo-level manifests, not skill-internal files, so whether they move into `skills/parley-deck/` is a design call — but the installer currently stages them from `packageRoot`, and under the new layout that source path changes. This is implicitly covered by S6 gates 1 and 5 but should be stated in S4a for completeness. S6: the gate list is complete for what could break a channel — Homebrew and WinGet correctly inherit from gates 2-3, the `aionrs` gap I documented is a capability gap not a build gate, and `skills update` is a runtime concern not a shipping gate. Skipping round 2 is defensible: all four agreed on every measurement, the one factual disagreement (claude-1's "no cheaper fix") was resolved by two independent refutations, and the remaining decision was the user's ruling. Round 2 would have re-argued settled facts. The reservation above does not block — it is a precision gap in S4a's enumeration, not an error in its binding principle.

### Signoff: kimi-1 — 2026-07-29
Status: 🟡 ACCEPT-WITH-RESERVATIONS

S1: rows 1, 2, 4, 6, 7, 8 and 9 reproduce exactly against my log (exps 1.0–1.6, `skills` 1.5.20, node v26.5.0), and the symlinks→0 row is correctly attributed to me. The heading "reproduced independently by all four participants" overstates, though: the `.claude/skills/` container row, the `--skill parley-design` row, and the S1d live-remote probes are not in my log (I probed rsync'd copies only), so I attest those as consistent with the early-return mechanism I read in source (my 1.8), not as measurements I made — and S1b's "fallback, consulted only when discovery finds nothing" is half a step past my evidence, since I could not isolate that case (my 1.2, NOT TESTED in isolation), though the conclusion (ship no manifest) stands regardless. S2 is accurate and, if anything, softens what happened: claude-1's file asserts not only "no cheaper fix exists" but "no manifest, container, **or flag** defeats it" (claude-1.md:11) — an affirmative false claim about flags, not merely an unenumerated interface — and the refutation replicated in three logs, not two, since my exp 1.6 also found `--full-depth` independently. S4a is correctly binding: it is my own §2 position, with my 1.7(b)/1.7(c) as the evidence that a file-only move manufactures the "malformed" verdict, and I co-sign hermes-1's reservation that `plugin.json` and `gemini-extension.json` are staged into the core destination from `packageRoot` (`PAYLOAD_ENTRIES`, `lib/installer.js:126-132`; Antigravity staging `:993-995`), so their post-move source location belongs in the move plan explicitly rather than by implication. S6 names one gap I consider the sharpest: gate 7 is `--list`-only, yet my 1.7(a)/(b) showed list and install diverge (discovery found 5 while the installed core arrived payload-less) — extend gate 7 to a real install from the merged remote into an isolated HOME asserting the core destination contains `references/` + `agents/` and no longer contains `bin/`, `lib/`, `package.json`, the only check that proves S3 fixed rather than asserting it "as a side effect". Skipping round 2 is defensible: every measurement agreed across four logs, the one wrong claim was refuted by replication, and the remaining decision was the user's ruling. For the record, F1's "neither labelled recommended" differs from my round-01 F1 (our installer keeps top billing for its native targets); I accept it as the group's wording call and disclose it here rather than leave it hidden — there is no buried disagreement.

### Signoff: codex-1 — 2026-07-29
Status: ❌ BLOCK

S1: I reproduced the baseline → 1, `skills/` plus root → 1, manifest plus root → 1, root removed → 4, and published layout with `--full-depth` → 5 rows; my log does not reproduce the `.claude/skills/` row, the `--skill parley-design` row, the symlink row, or the all-five-moved-under-`skills/` row, and it explicitly marks the live-remote probes NOT TESTED, so “reproduced independently by all four participants” contradicts my record. S2 softens the refutation: claude-1 affirmatively wrote that no flag defeats the shadow and that the skills were unreachable by any flag it could find, while the record also omits kimi-1's third independent `--full-depth` replication. S4/S4a correctly binds the user's layout ruling and the full `SKILL.md` + `references/` + `agents/` move, but it must explicitly disposition `plugin.json` and `gemini-extension.json`, both current required/copied payload entries, rather than leave their post-move source path implicit. S6 is incomplete: gate 7 must perform an actual published-repository install into an isolated destination, assert that the installed core contains `references/` and `agents/` but not repository `bin/`, `lib/`, or `package.json`, and then exercise `skills update`. The Homebrew formula (`brew upgrade`/`brew test`), WinGet manifest (`winget validate`/`winget install`), and legacy Gemini extension install are also breakable distribution channels not explicitly gated; a Windows binary build alone does not test WinGet, and a context-path question alone does not test Gemini. Skipping round 2 was defensible for the user-resolved layout choice, but not for claiming universal measurement agreement or leaving these payload/channel decisions unresolved. Counter-proposal: correct the attribution and claim record, make the manifest disposition explicit, add the missing end-to-end gates, and then reopen signoff after a short round 2.

---

## Amendments A1–A4 (drafted by claude-1 after the round-1 signoffs)

codex-1 blocked with four counter-proposals; hermes-1 and kimi-1 raised two of the same
points in their own reservations. All four are accepted. These amendments are binding and
supersede the text above wherever they conflict.

**A1 — the measurement table's attribution is corrected.** The heading "reproduced
independently by all four participants" was **false as written**, and all three signoffs said
so from their own logs. Corrected attribution:

| Row | Measured by |
|---|---|
| baseline → 1 | claude-1, codex-1, hermes-1, kimi-1 |
| `skills/` container + root → 1 | claude-1, codex-1, hermes-1, kimi-1 |
| `.claude-plugin/marketplace.json` + root → 1 | claude-1, codex-1, hermes-1, kimi-1 |
| root removed → 4 | claude-1, codex-1, hermes-1, kimi-1 |
| published layout + `--full-depth` → 5 | claude-1, codex-1, hermes-1, kimi-1 |
| `.claude/skills/` container + root → 1 | claude-1 only |
| `--skill parley-design` + root → 1 | claude-1 only |
| symlinked `skills/` → 0 | kimi-1 only (hermes-1's variant kept a real `SKILL.md` and found 1) |
| all five moved under `skills/` → 5 | claude-1, kimi-1 |
| live-remote probes (S1d) | claude-1, hermes-1; **NOT TESTED** by codex-1 (no network) and kimi-1 |

**S1b is downgraded.** "A fallback consulted only when discovery finds nothing" is a *model*
inferred from three observations, not an isolated measurement — kimi-1 marked the isolating
case NOT TESTED. The operational conclusion (ship no manifest) is unaffected.

**A2 — the refuted claim is restated without softening.** claude-1 did not merely fail to
enumerate the interface. It wrote, affirmatively: *"no manifest, container, **or flag**
defeats it"* and *"unreachable by any flag I could find"* (`round-01/claude-1.md:11`, and its
experiment-6 note). That is a positive false claim about flags, made after searching, not an
omission. It was independently refuted by **three** participants — codex-1, hermes-1 **and
kimi-1** — not two.

**A3 — `plugin.json` and `gemini-extension.json` get an explicit disposition.** Raised by all
three signoffs. Both are current `PAYLOAD_ENTRIES` (`lib/installer.js:126-132`) staged into the
core destination from `packageRoot`, and both are required by the Antigravity validator
(`:1098`). **Decision: they stay at the repository root**, and the installer's *source* path
for them is set explicitly in the same commit rather than left to follow from the move. The
installed-destination shape does not change.

**A4 — the gate list is extended.** S6 gate 7 is replaced, and four channel gates are added:

- **G7 (replaces the `--list` version).** An **actual install** from the *published* repo
  into an isolated `HOME`, asserting the installed core **contains** `references/` and
  `agents/`, **does not contain** `bin/`, `lib/`, `package.json`, and that all five skills
  install **without `--full-depth`**. kimi-1's evidence is the reason: discovery and install
  diverge, and a listing that says 5 proved nothing about what arrives on disk.
- **G8** `skills update` against the new layout.
- **G9** Homebrew: `brew upgrade` and `brew test`.
- **G10** WinGet: `winget validate` on the manifest and an actual `winget install`. codex-1 is
  right that a Windows binary build alone does not test WinGet.
- **G11** legacy Gemini: an actual `gemini extensions install`. A `contextFileName` path
  question is not a substitute for running it.

Where a gate cannot be run in this environment (Windows, and any CLI not installed here) the
result is recorded as **NOT TESTED** with the reason — never as a pass.

**On skipping round 2.** codex-1 is right that it was defensible for the user-resolved layout
choice and *not* defensible as a basis for claiming universal measurement agreement. A1 and A2
correct that claim; A3 and A4 resolve the two open decisions. This amendment block **is** the
round-2 output, and signoffs are reopened below.

---

## Signoffs — round 2 (after amendments A1–A4)

### Signoff: codex-1 — 2026-07-29 (round 2)
Status: ✅ ACCEPT

A1 matches my log: the five rows attributed to codex-1 are the five I measured, the four rows not attributed to me remain outside my log, and the live-remote probes are correctly marked NOT TESTED; no row is still mis-attributed. A2 now states the affirmative false claim about flags and all three independent refutations without softening. A3 is the right disposition: `plugin.json` and `gemini-extension.json` remain repository-root manifests, while the installer source path is updated explicitly in the same commit. A4 now covers the missing real-install, update, Homebrew, WinGet, and legacy Gemini channel checks alongside the original gates; I see no remaining channel-breaking omission. The amendment block adequately serves as the round-2 output, so I do not require a full round 2.
