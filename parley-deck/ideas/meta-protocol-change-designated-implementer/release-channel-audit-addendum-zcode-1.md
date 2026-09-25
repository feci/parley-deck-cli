---
idea: meta-protocol-change-designated-implementer
artifact: release-channel audit ADDENDUM — 15-root installation delta + report corrections (only new canonical report; original release-channel-audit-zcode-1.md untouched)
agent: zcode-1
role: existing non-implementer independent reviewer (participant; not organizer, not implementer)
date: 2026-09-25
cli-tag: v1.50.0 -> 01fc49526ca66b771b549a01cfbbf4addeaa886b (unchanged)
skill-tag: v2.14.0 -> a5664d803f1fb6156ef95ff5921dcf13a3dd2031 (unchanged)
verdict: PASS — organizer's --include-undetected installation delta independently verified at all 15 registry roots (90/90 SKILL.md, full payloads byte-identical to shipping 2.14.0, 15/15 bundled COOPERATION = staged core, 15/15 markers 2.14.0); 12 detected / 3 dormant split unchanged, no runtime tool became installed; three documentary corrections to my own report recorded below (no source defect implied); external pending items still pending, never claimed
evidence: release-delivery/2026-09-25-designated-implementer/audit-zcode-1/addendum/ (INDEX.md itemizes; original audit-zcode-1/ evidence untouched)
---

# Release-channel audit addendum — zcode-1, 2026-09-25

Narrow follow-up to my final release-channel audit (PASS, filed, intact). Scope: (a) independently
verify the organizer's election to fulfill **every existing runtime copy** — including the three
dormant managed roots — via the exact audited tarball; (b) correct three documentary ambiguities in
my own report; (c) refresh read-only external statuses. No Go/npm suites re-run, no channel
downloads (reused my existing tag archive and original evidence), no commit, no source/tag/config/
roster mutation, no installation or publication by me.

## 1. Installation delta — organizer action verified

**Organizer input (read for identity only):** `node DELIVERY/skill/install-source/package/bin/
parley-deck-skill.js install --target all --include-undetected --json` → exit 0, **15 targets ×
6 skills = 90 actions, all ok/replaced** (`delivery/skill-install-all-including-dormant.json`,
markers written 2026-09-25T09:01:12–13Z). The earlier detected-only run (`skill-install-all.json`)
covered exactly the 12 detected targets — goose/cursor/aionrs absent — so the delta is precisely
those three dormant roots (the 12 detected roots were idempotently re-replaced with identical
content).

**Chain of custody re-checked (no new downloads):** `install-source/package/skills` re-diffed
against my originally downloaded GitHub tag archive (`audit-zcode-1/downloads/skill-v2.14.0.tar.gz`)
— trees identical; npm tarball re-hashed `681b9076…135b` (unchanged). Original chain (runtime ==
tag archive; install-source == npm tarball; tarball == git `a5664d8`) therefore still holds.

**Independent verification — all 15 registry roots, full payloads:**

- **90/90 SKILL.md copies** hashed; every one equals the shipping 2.14.0 reference
  (`parley-deck db71eca5…c03f`, `parley-bidding aad7d7d8…9540`, `parley-design 9ccf8e9d…d1bd`,
  `parley-design-check 337aafd6…2b1a`, `parley-tracker 25b7a6b2…1198`,
  `parley-worktrees 91f94c40…c3a9`).
- **Full payload diff vs shipping source: content=0 at every root.** Only extras are
  installer-owned (`.parley-deck-skill-install.json` marker; staged `LICENSE`/`README.md`/
  `gemini-extension.json`/`plugin.json`; agy also stages the plugin-bundle copy `skills/SKILL.md`,
  itself `db71eca5…c03f`); only missing file is `parley-addon.json` (consumed at install, by
  design). No symlinks anywhere. Zero problems in the machine verdicts.
- **15/15 bundled `parley-deck/references/COOPERATION.md` = `51476d69…f67a`** — byte-equal to the
  staged core `~/.parley/staging/COOPERATION-2.14.0.md` (re-hashed, unchanged, 115,010 B, 0444).
  The dormant roots' pre-refresh 2.10.0-era COOPERATION (`f41c8e5f…e6b2`) is fully superseded.
- **15/15 markers**: `version 2.14.0`, `source npm:parley-deck-skill@2.14.0`, target matches root,
  markerSchema 2; `references/compatibility.json skillVersion 2.14.0` at all 15 (was 2.10.0 at the
  three dormant roots before).
- **Detected/dormant split unchanged — 12 detected sessions + 3 dormant managed roots.** Detection
  inputs did not change: `goose`/`agent`/`aionrs` still not on PATH; `~/.goose`, `~/.cursor`,
  `~/.aionrs` still contain only `skills/` with the installer-staged parley trees. **No runtime
  tool became installed** — the refresh changed managed skill copies at dormant roots, not runtime
  presence; the product's own semantics still classify those roots as not-detected.

Full updated path/hash table (+ markers, compatibility, machine verdicts):
`audit-zcode-1/addendum/runtime-hash-table-15roots.txt` (JSON: `runtime-verdicts-15roots.json`).
My original report's "dormant 2.10.0-era copies" organizer-awareness note is thereby resolved by
the organizer's action; nothing else in the original report's runtime section changes.

## 2. Documentary corrections to my original report (reporting only; no source defect)

1. **Misnamed accepted-residual NIT (original "Limitations" item 6).** My item 6 labeled
   "Accepted residual NIT" as *zcode-1's source-only re-entry NIT*. The closing cycle-4 consensus'
   actually unrepaired NIT is **claude-1's round-04 NIT-1**: the AF-16 replacement lead-in at
   `IMPLEMENTATION.md:325-326` — the bolded "All shadow byte figures in this bullet are rendered
   WITH `--flag auto_implement --flag protocol_change`" — which over-reaches because the same
   bullet states the WITHOUT-flags exception (the flagless 86,701 B figure). Canonical
   `review/consensus.md` (preserved at `804522c`): VC-4.1 reconciles claude-1's NIT-1 with my R-B
   refutation as a scope divergence; disposition "ACCEPTED RESIDUAL, recorded dismissal — the
   finding was real and is NOT fixed"; "the AF-16 lead-in stands as written and unrepaired". My
   source-only re-entry NIT is a **separate, separately-carried older dismissal**
   (`IMPLEMENTATION.md` item 7, "Dismissed-finding documentation … documentation alternative") —
   real and correctly recorded, but not the closing consensus' open item. Both dismissals stand,
   untouched; my report conflated the two labels.
2. **Runtime hash table link.** Original prose cited `audit-zcode-1/runtime/runtime-hash-table.txt`;
   the file's actual location (per INDEX and on disk) is `audit-zcode-1/runtime-hash-table.txt`.
   Correct location hereby confirmed; the 15-root successor table is at
   `audit-zcode-1/addendum/runtime-hash-table-15roots.txt`.
3. **Stale skill working-tree remark (original "Limitations" item 5).** The "skill repo working
   tree checked out at `8161e5e` (pre-release)" refers to the **separate main checkout**
   `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill` (branch `main`, simply not
   pulled past the pre-release ancestor). This task's worktree
   `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer-skill` is
   clean at exactly `a5664d803f1fb6156ef95ff5921dcf13a3dd2031` (= tag v2.14.0), in sync with
   `origin/main`. Remote/tag authoritative in both cases, as originally stated.

## 3. Read-only external status refresh (since the original audit)

| item | state now | change |
|---|---|---|
| winget PR [441054](https://github.com/microsoft/winget-pkgs/pull/441054) | **OPEN, not merged, not claimed merged**; MERGEABLE; steps 01–07 (PR/manifest/URLs/domain/policy/catalog/installers-scan) + `license/cla` now SUCCESS; 08 Installation Validation IN_PROGRESS; 09/10 QUEUED | external validation advanced; merge still pending |
| npm `parley-deck-skill` | `dist-tags.latest = 2.13.0`; **2.14.0 still absent** | unchanged — PENDING OWNER, not published |
| core store | `~/.parley/protocol/core/` still only 2.10.0 + 2.13.0; staged 2.14.0 intact (`51476d69…f67a`) | unchanged — PENDING OWNER-ONLY; publish never run |

## 4. Verdict

**PASS** for this narrow installation delta: all 15 registry roots now hold byte-identical 2.14.0
payloads from the exact audited tarball, with correct markers and bundled core hashes; the
12 detected / 3 dormant classification is unchanged and no runtime tool is newly claimed installed.
The original audit's channel PASS remains valid except as superseded above; npm publication and
core protocol publish remain PENDING OWNER-ONLY; the winget skill PR remains pending external
Microsoft validation/merge; Windows CI failure remains owner-deferred (DF-4), reported exactly.
Nothing committed, published, or mutated by this addendum.
