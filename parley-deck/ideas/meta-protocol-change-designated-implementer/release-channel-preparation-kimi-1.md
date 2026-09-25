# Release channel preparation handoff — kimi-1 (sole implementer), 2026-09-25

Scope: PREPARE release channel metadata only. Nothing pushed, published, installed, or mutated
globally; source candidates stay frozen (CLI `01fc495…886b`, skill `a5664d8…2031`). Organizer
has already published GitHub CLI `v1.50.0` (tag 01fc495) and skill `v2.14.0` (tag a5664d8);
skill portable workflow run **36113011791** succeeded and the final uploaded assets are the
Windows x64/arm64 pair. The npm exact-tarball publication failed with E404 access — owner
login requested by the organizer; the tarball prepared in
`release-preparation-kimi-1.md` (sha256 `681b9076…135b`) remains the exact artifact for the
owner's later `npm publish --access public <tarball>`.

## 1. Homebrew tap — BOTH formulae prepared, NOT pushed

Repo `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley` (was clean at
`52b09f3`, `main` == `origin/main`; fetched safely, no dirty work touched).

**Commit `9855366`** (supersedes tree-identical `71f78eb`; identity-only amend to this repo's
local user config):
`[codex] meta-protocol-change-designated-implementer: release CLI 1.50.0 and skill 2.14.0 in Homebrew (authored by kimi-1)`

- `Formula/parley-deck-cli.rb` — url `…/v1.50.0.tar.gz`, sha256
  `2ff8426d257ec4a6b9d92ac34139ef97d92b3a4259d6c2b25d25c424d73ce6a3`
- `Formula/parley-deck-skill.rb` — url `…/v2.14.0.tar.gz`, sha256
  `f43e02f7db2050c0c778c5f232c091c417961ea66c079c34ce061faa8a0feafe`

Hashes are from **freshly downloaded actual GitHub tag archives** this session (CLI archive
9,221,841 B with `VERSION`=1.50.0 inside; skill archive 490,129 B with package.json
`"version": "2.14.0"` inside) — contents sanity-checked, formula mechanics untouched.

**Check results (real, invokable on this host):** `ruby -c` both files Syntax OK;
`brew style` both files — 2 inspected, **no offenses**. `brew audit --strict --online` by
path is disabled in this Homebrew; by name it would audit the separately tapped checkout at
`/opt/homebrew/Library/Taps/feci/homebrew-parley` (the old published formulae), not this
working copy — not invokable without mutating the installed tap, so it is recorded as not run
rather than claimed; the audit's download leg was performed manually via the actual fresh
downloads above. `brew upgrade` / `brew test` remain organizer channel steps post-push.

## 2. Skill WinGet `Feci.ParleyDeckSkill` 2.14.0 — manifests prepared, NOT pushed/PR'd

Isolated sparse clone
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-25-designated-implementer/winget`,
branch `feci-skill-2.14.0` (origin=feci fork, upstream=microsoft, clean master base
`f4aeb1ab5`). One application only; **no CLI winget file or PR** (stays held).

**Commit `26720f091`** — exactly three new files under
`manifests/f/Feci/ParleyDeckSkill/2.14.0/` (installer, locale.en-US, version), pattern copied
from merged 2.13.0 (PR440360), schema 1.12.0 unchanged; updated version/URLs/digests and the
Description's CLI pairing (1.50.0) only.

Digests from the **final release downloads** (never local builds), matched to the release API
twice — API digest comparison and manifest-vs-downloaded-bytes comparison:

| arch | bytes | sha256 (uppercase, as committed) |
|---|---|---|
| x64 | 91,402,916 | `F3DA19CBD2F0FBBF3551899200E07264979F03A6D5F29D2D9CB44C7F736AF9CF` |
| arm64 | 85,815,972 | `A76A375972EF6E0A2236CC98CFB83F41FD8359707E3467613BE0E4A8B26BA2A5` |

**Validation actually run (macOS):** YAML parses of all three; key sets identical to 2.13.0;
cross-file PackageVersion/ManifestType/ManifestVersion consistent; 64-hex uppercase hash
format; v2.14.0 asset URLs (the same URLs the bytes were downloaded from); ReleaseNotesUrl
v2.14.0. **Actual limitation, stated plainly:** `winget validate --manifest` and
`winget install --manifest` are Windows-client commands and are **not invokable on this macOS
host** — not run, not claimed; Windows-side validation remains an organizer step before the
fork PR (per repo AGENTS.md: one package per PR, summarize validation — this paragraph is the
summary).

## 3. GitHub artifact scope — 2 assets is by design, not a defect

Predecessor v2.13.0 shows 5 assets; the v2.14.0 release shows 2 (the Windows pair). Checked
against the actual committed sources at the published tag, not assumed:

- `.github/workflows/release-portable.yml` (at `a5664d8`): job `macos-tests` is validation
  only (`npm test`, no upload); job `windows` runs `npm test` + **`npm run
  build:portable:windows`** and uploads `dist/*` — the workflow's designed output is exactly
  the Windows pair. Run 36113011791 succeeded, so the committed workflow produced everything
  it is designed to produce.
- `RELEASING.md` documents exactly two expected artifacts
  (`dist/parley-deck-skill-vX.Y.Z-windows-x64.exe` / `-windows-arm64.exe`) and names them the
  WinGet inputs.
- No documented install path consumes macOS/Linux portable binaries: npm
  (`npx -y parley-deck-skill@latest install`, README :145-220), Homebrew (builds from the tag
  archive with node), WinGet (Windows pair). README has no portable-binary install path.
- The predecessor's extra 3 (macos-arm64, macos-x64, linux-x64) were out-of-band manual
  uploads, not workflow output, and have no documented consumer.

**Conclusion: no missing-asset defect; no extra builds made; the final Windows assets were not
rebuilt or clobbered.** If the owner later wants the convenience binaries back, that is a new
decision (manual upload or a workflow change), not a release blocker.

## 4. Status and remaining organizer/owner steps (none performed by me)

Done by organizer already: GitHub releases + tags, portable workflow. Pending: npm owner login
then `npm publish --access public <exact tarball>` and `npm view` verification; push Homebrew
tap commit `9855366` then `brew upgrade`/`brew test`/`--version` checks; push winget branch
`26720f091`, Windows-side `winget validate`/`install --manifest`, then the one-application
fork PR; install skill to all sessions + participant runtime SKILL.md hash audit;
owner-attended `parley protocol publish --version 2.14.0 --from
/Users/tomasfecko/.parley/staging/COOPERATION-2.14.0.md` (unchanged, sha256
`51476d69…f67a`); post-release owner-only machine default `default_implementer = "codex-1"`.
OBS-1/3 release notes already updated by the organizer; OBS-2 evidence retained by
`release-evidence-logs-kimi-1.md`.
