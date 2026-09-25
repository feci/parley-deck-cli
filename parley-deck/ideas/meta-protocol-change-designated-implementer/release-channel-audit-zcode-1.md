---
idea: meta-protocol-change-designated-implementer
artifact: final release-channel audit (published channels)
agent: zcode-1
role: existing non-implementer independent reviewer (participant; not organizer, not implementer)
date: 2026-09-25
cli-tag: v1.50.0 -> 01fc49526ca66b771b549a01cfbbf4addeaa886b
skill-tag: v2.14.0 -> a5664d803f1fb6156ef95ff5921dcf13a3dd2031
verdict: PASS — agent-controlled delivery channels verified end-to-end; npm publication and core protocol publish remain explicitly PENDING OWNER-ONLY; skill WinGet PR pending external Microsoft validation/merge; Windows CI failure owner-deferred (DF-4), reported exactly, never claimed green
evidence: release-delivery/2026-09-25-designated-implementer/audit-zcode-1/ (see INDEX.md)
---

# Final release-channel audit — zcode-1, 2026-09-25

Independent audit of the **actual published channels** (not preparation). Every download,
hash, comparison and runtime command below was executed by me this session; nothing was
published, mutated, committed, or forced. Machine evidence: `audit-zcode-1/` (INDEX.md
itemizes every file). Organizer's inputs were read for candidate identities only; every
channel fact was re-derived from the remote registries/release APIs/disk.

## Verdict by channel

| channel | state | basis (independently verified) |
|---|---|---|
| CLI git: tags + main ancestry | **PASS** | `git ls-remote`: tag `v1.50.0` → `01fc4952…886b` (exact audited candidate, lightweight); main `b6ceea4`; `c49b464 ⊂ 01fc495 ⊂ b775ae7 ⊂ b6ceea4` all fast-forward, no force; candidate→main delta over product paths **empty** (records-only: 6 idea-doc/inbox files) |
| CLI GitHub release (6 binaries + sha256.json) | **PASS** | Fresh `gh release download` of all 7 assets; my SHA256 = API `digest` = published `sha256.json` = organizer `build-provenance.json` on **all six**; `go version -m` on each: `vcs.revision=01fc4952…`, `vcs.modified=false`, `-trimpath`, correct GOOS/GOARCH ×6; darwin-arm64 run → `parley 1.50.0` exit 0 |
| Windows assets labelled | **PASS** | Asset labels on the release: `Windows ARM64 — experimental/unvalidated`, `Windows x64 — experimental/unvalidated`; release notes body states it plainly and cites candidate CI run 36110596716 |
| CLI winget | **HELD (as required)** | winget-pkgs `manifests/f/Feci/ParleyDeckCli` stops at 1.48.0; no 1.49.x/1.50.0 dir; no open feci PR for it |
| Skill git: tag + main | **PASS** | tag `v2.14.0` → `a5664d8…2031` = remote main (fast-forward from 8161e5e); portable workflow run **36113011791** `success` at head `a5664d8…` |
| Skill GitHub release (2 assets) | **PASS** | Exactly 2 assets (the Windows pair), per the **committed** workflow/RELEASING/README requirement (Kimi's classification confirmed against sources; the "5 assets" predecessor count was out-of-band manual uploads). My downloads hash `f3da19cb…9cf` (x64, 91,402,916 B) and `a76a3759…2a5` (arm64, 85,815,972 B) = API digests = winget manifest digests |
| Skill WinGet PR | **OPEN, external validation pending — not merged, not claimed** | PR [441054](https://github.com/microsoft/winget-pkgs/pull/441054) OPEN (not draft), head `26720f091` on fork `feci:feci-skill-2.14.0`, **exactly 3** manifest files under `manifests/f/Feci/ParleyDeckSkill/2.14.0/`, one application/version; committed installer URLs point at the **final uploaded** release assets whose bytes I hashed; `license/cla` success, pipeline steps queued/skipped, `merged_at: null` |
| npm | **PENDING OWNER — NOT published, NOT claimed** | Registry queried live: `dist-tags.latest = 2.13.0`; **no 2.14.0 version exists** (E404). First publish failed E404/auth (log retained); owner login requested in inbox `codex-1-to-user_…_npm-login.md`; exact tarball intact at `delivery/skill/parley-deck-skill-2.14.0.tgz`, sha256 `681b9076…135b` / sha512 `897b1d35…0be7` re-verified by me |
| Homebrew (both formulae) | **PASS** | Remote tap main = `9855366`; installed tap at same commit (clean FF from 52b09f3). Formula hashes equal my **fresh tag-archive downloads**: CLI `2ff8426d…ce6a3` (9,221,841 B, `VERSION`=1.50.0), skill `f43e02f7…feafe` (490,129 B, `package.json` 2.14.0). Installed: `parley-deck-cli 1.50.0`, `parley-deck-skill 2.14.0`; `/opt/homebrew/bin/parley` → Cellar 1.50.0; `parley --version` → `parley 1.50.0`; `parley-deck-skill --version` → `2.14.0`. **I ran `brew audit` (exit 0) and `brew test` on both (exit 0 ×2)** — skill test exercises a sandboxed install+doctor |
| Main CI | **PASS with reported Windows red (owner-deferred DF-4)** | CLI run 36112979464 @ `b775ae7` (= candidate + records-only handoff file; source-identical): ubuntu ✅, macOS ✅, windows ❌ (inherited failure set — Claude's prep audit diffed it to zero newly-failing tests; designation suite never executes there). Skill run 36112998102 @ `a5664d8`: ✅ both jobs. Later records commits (4ecff5b, f512f66, b6ceea4 `[skip ci]`) triggered no CI and change no product path (diff empty) — no drift |
| Runtime skill install (all real sessions) | **PASS** | Full independent enumeration below; 12 detected targets × 6 skills = 72 copies, **every payload file byte-identical to the v2.14.0 shipping source** |
| Product ships UNSET | **PASS** | Tag-archive `internal/config/runtime.go:654` writes `# default_implementer = "agent-id"` **commented**, "SHIPPED UNSET ON PURPOSE"; no active `default_implementer` in any shipped TOML at the tag |
| Owner default `codex-1` | **PASS as post-release owner config** | `~/.parley/agents.toml` `[defaults].default_implementer = "codex-1"`, mtime 2026-09-25T08:43:51Z matches `delivery/owner-default-operation.json` (`previous: null → codex-1`, only that key, after GitHub+Homebrew publication); owner direction `source-context/owner-default-2026-09-25.md` pre-authorizes exactly this and explicitly says it does not rewrite in-flight decks/ideas/signatures |
| Core protocol 2.14.0 | **PENDING OWNER-ONLY — store untouched, never claimed published** | `~/.parley/protocol/core/` contains only 2.10.0 and 2.13.0. Staged `~/.parley/staging/COOPERATION-2.14.0.md` intact: sha256 `51476d69…f67a`, 115,010 B, mode 0444 — byte-identical to the COOPERATION.md bundled in **every** runtime skill copy (verified at all 12 targets). Attended `parley protocol publish` remains owner-only; no TTY allocated by any participant |

## Runtime inventory — derived, not assumed

The installer registry (read from the exact shipped `lib/installer.js`) names **15 targets**:
codex, claude, agy, gemini, hermes, qwen, codebuddy, goose, kimi, droid, vibe, cursor,
opencode, aionrs, zcode. `--target all` installs **detected** targets only; detection =
runtime command on PATH or real (non-marker-only) runtime-home evidence. This host:
**12 detected** (install JSON + my re-run of `doctor`: both say 12, ok, 6 skills each) —
**goose, cursor, aionrs are supported-but-undetected** (no `goose`/`agent`/`aionrs` command;
their skill dirs contain only installer-owned marker content, which the product's own rule
excludes as "not real sessions"). They are correctly **not** called installed.

**Hash audit (all 15 roots + symlink scan + outside-registry scan):** no symlinks anywhere in
any target chain or payload; no skill copies outside the registry roots (`~/parley-deck`,
`~/.local/share/parley-worktrees`, `~/AI_WORKSPACE/bondShift/parley-deck` are deck
workspaces/metadata, no SKILL.md). At the 12 installed targets, all six skills' trees are
**byte-identical to the shipping source** (content diff = 0; the only extras are the
installer's own `.parley-deck-skill-install.json` marker and, for the core skill, the
package-root `plugin.json`/`gemini-extension.json`/`README.md`/`LICENSE` staged by design;
`parley-addon.json` is consumed at install). Shipping SKILL.md sha256 reference values
(matched at every installed target, and at the 3 stale roots for the five unchanged add-ons):
`parley-deck db71eca5…c03f`, `parley-bidding aad7d7d8…9540`, `parley-design 9ccf8e9d…d1bd`,
`parley-design-check 337aafd6…2b1a`, `parley-tracker 25b7a6b2…1198`,
`parley-worktrees 91f94c40…c3a9`. Full per-path table: `audit-zcode-1/runtime/runtime-hash-table.txt`.
Bundled protocol verified beyond SKILL.md: every runtime `parley-deck/references/COOPERATION.md`
hashes `51476d69…f67a` (= staged 2.14.0 core). Chain of custody: runtime == GitHub tag-archive
skills (my diff) and organizer's install-source == exact npm tarball (my diff) which Claude's
prep audit proved 210/210 == git `a5664d8`.

**Reported for organizer awareness (not a delivery defect):** dormant **2.10.0**-era copies of
the six skills exist at `~/.goose/skills/`, `~/.cursor/skills/`, `~/.aionrs/skills/`
(parley-deck there is 2.10.0: `compatibility.json skillVersion 2.10.0`, older SKILL.md
`fe9f0468…34dc`, COOPERATION.md `f41c8e5f…e6b2`). These roots are undetected per the
product's semantics, so `--target all` correctly skipped them; no real session reads them.
Optional discretionary cleanup/refresh is the organizer's choice; no re-audit obligation
arises because no required installation is missing or stale.

## Post-release default nuance — verified without manufacturing eligibility

The frozen idea keeps `facilitator: codex-1`, `participants: [claude-1, kimi-1, zcode-1]`,
and `IMPLEMENTATION.md` re-entry pin `implementer: kimi-1`; both consensuses remain signed
ACCEPT by all three; the idea directory is clean at HEAD `b6ceea4` (no post-signature
mutation). The shipped resolver (`internal/protocol/implementer.go` at the tag) ranks the
re-entry pin **above** the global default, so the new machine default `codex-1` cannot
rewrite this idea. Documented fall-through boundary (source + release notes): when the
global-default designee is unavailable or ineligible for an idea, a one-line notice surfaces
and dispatch falls through to the existing chain (`fall-through-unavailable` /
`fall-through-inapplicable` sources); a malformed/ineligible **per-idea** designation stays
fail-closed. Live read-only `parley status --idea` on this deck resolves around kimi-1/zcode-1
with codex-1 outside the quorum — consistent with no rewrite. No roster changes were made by
anyone (`owner-default-operation.json` records only the defaults key; roster blocks intact).

## Limitations, deferred and owner-only items (never counted complete)

1. **npm 2.14.0 — PENDING OWNER**: owner `npm login`, then `npm publish --access public
   <exact tarball>` + `npm view` verification. Tarball integrity re-verified by me.
2. **Core 2.14.0 — PENDING OWNER-ONLY**: attended `parley protocol publish --version 2.14.0
   --from ~/.parley/staging/COOPERATION-2.14.0.md` (TTY-gated; no participant may run it).
3. **WinGet skill PR — PENDING EXTERNAL**: microsoft/winget-pkgs #441054 validation pipeline
   and merge; local `winget validate/install` not runnable on macOS and never claimed.
4. **Windows CI red — owner-deferred (DF-4, `windows-portability`)**: stated exactly; the
   Windows test binary aborts before the designation suite executes, so Windows provides no
   validation of this feature; assets stay experimental/unvalidated; CLI winget held.
5. Local-only observations, immaterial to channels: worktree's local `main` ref is a stale
   ancestor of `b6ceea4` (no divergence); skill repo working tree checked out at `8161e5e`
   (pre-release) — remote/tag authoritative in both cases.
6. Accepted residual NIT (zcode-1's source-only re-entry NIT, documentation alternative) and
   DF-1…DF-4 remain recorded in `IMPLEMENTATION.md`, untouched.

Nothing was committed by this audit. Final completion/done remains the organizer's to write
after this audit (and after any owner actions they choose to await); this PASS truthfully
covers agent-controlled channels only.
