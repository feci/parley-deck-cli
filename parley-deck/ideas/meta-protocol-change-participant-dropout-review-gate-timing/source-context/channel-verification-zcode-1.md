# Channel verification — zcode-1 (post-close, AC14 delivery)

- **Agent:** zcode-1 (independent verifier, separate process; not the implementer/reviewer of record)
- **Date:** 2026-10-09 (verification window 09:31–09:43 UTC)
- **Context attestation:** `context_mode=full`, `source_sha256=packet_sha256=acbd4dbc0c0702bc191176bb80bcee32c5093c8b4ebbee42e036df6a9b7d1137`, **no fallback**. I read
  `ORGANIZER-BRIEF.md` (sha256 `50cb660916b882965ff289dc2fc7d5ff9c8f234cac61d2e95e1e5c8849c8722a`, byte-identical to the recorded
  `final-integrity-recheck.json` `brief_sha256`) and the full phase-8 protocol packet end-to-end before verifying.
- **Method:** every claim below was re-derived from the thing itself — live GitHub/GitHub-API state via `gh api`, fresh downloads hashed locally
  (`shasum -a 256`), `go version -m` on the hosted binaries, full `diff -rq` of the skill tarball against the merged source tree, an independent
  re-application of all 14 protocol hunks, and a re-hash of all 114 runtime `SKILL.md` files. Evidence JSONs were treated as claims to refute,
  not as conclusions to trust. No product tests or review rerun (not required); no Claude, credential/model/roster changes, publication, installs,
  commits, GitHub messages, or worktree registration. Only this report file was written.

---

## 1. GitHub channels (CLI PR78, skill PR10, releases, assets, build identity)

**Merges and tags — live-confirmed.**
- CLI PR78 `https://github.com/feci/parley-deck-cli/pull/78`: `MERGED`, mergeCommit `128e30b479065a98437844886d295ae83e9b0e52`, mergedAt
  2026-09T09:20:44Z. Tag `v1.53.0` → same commit; `main` HEAD → same commit.
- Skill PR10 `https://github.com/feci/parley-deck-skill/pull/10`: `MERGED`, mergeCommit `8ce4dec0fde0f563f5110e99d8b7e4550d939770`, mergedAt
  09:20:47Z. `main` HEAD → same commit; release `target_commitish` → same commit.

**Releases and assets — every digest independently reproduced.**
- CLI `v1.53.0` (published 09:29:31Z, not draft/prerelease): 7 assets. Skill `v2.17.0` (published 09:22:13Z): 4 assets.
- I re-hashed all 11 downloaded copies in `hosted-assets/`: all equal the values in `hosted-asset-verification.json`, byte sizes included.
- I fetched the live release API digests: all 11 match my local hashes (CLI: 6 binaries + `sha256.json`, e.g. darwin-arm64
  `9488d1f80ddd6c8e4291f8f43a190eae89d4e23365b104e2d793e006577a45d9`; skill tgz `40b6eb424cfc003864f06fa17065ea1ac997637e47d76f0f86aef0fbcfbae3ca`;
  windows-x64.exe `0bd51b98ab36a520d5dbbd2bf60645f83796a7447b451d6972201ce3d4700821`).
- Hosted `sha256.json` contains exactly the six CLI binary entries; each equals my independently computed hash. (`cli_checksum_manifest_matches_all_six_binaries: true` — re-derived, not trusted.)
- `prepared_asset_matches: null` on the two skill Windows `.exe`s is expected: those assets are built and uploaded by the release workflow, not
  prepared locally; I re-hashed both against the live API digests and the WinGet manifest values instead — all agree.

**Clean-source build proof — independently extracted from the hosted binaries.**
`go version -m` on `parley-v1.53.0-darwin-arm64`, `-linux-x64`, `-windows-arm64.exe` (3 of 6 sampled): each embeds
`vcs.revision=128e30b479065a98437844886d295ae83e9b0e52`, `vcs.time=2026-10-09T09:20:43Z`, **`vcs.modified=false`** (clean source), go1.27.1,
`-trimpath`, `CGO_ENABLED=0` — exactly as recorded for all six in `cli-build-verification.json`.

**Skill tarball/source identity — full byte comparison.**
Extracted the hosted `parley-deck-skill-2.17.0.tgz` (210 packed files, `package.json` 2.17.0) and ran `diff -rq` against the skill worktree at
`8ce4dec`: **zero content differences**; the only asymmetries are repo-only files legitimately absent from the npm package (`.git`, `.github`,
`CHANGELOG.md`, `dist`, `node_modules`, `packaging`, `scripts`, `test`, …). All six `SKILL.md` hashes match `skill/payload-verification.json`
(e.g. `parley-deck` `0b9769e8a58fd83cdd9d1d92be5f0d755f3fe7cbba30be20126d5eea9cd6af0c`).

**CI — reported honestly, nothing pending/cancelled counted as passed.**
- Skill Windows release workflow **37910828042** ("Release portable binaries", release event at `8ce4dec`): `completed/success` — verified live
  and per-job (macos-tests ✓, windows ✓ including `npm run build:portable:windows` and asset upload). Skill merge `Tests` run 37910674348
  (push main at `8ce4dec`): success.
- CLI merge-CI **37910668038** ("Tests", push main at `128e30b`): **still `in_progress` at my final poll 09:42:31Z** — ubuntu-latest success
  (09:24:34Z), macos-latest success, windows-latest mid-Test. **Not a pass; not counted.** Context that keeps this non-blocking for channel
  delivery: the windows-latest Tests job has failed on **every** completed run of this branch (4/4, e.g. 37906220742 windows-only fail:
  TestSpawnStopDrainsStderrBeforeReaping &c.) **and on main runs predating this idea** — baseline `85c5a8f` (CLI 1.52.0) run 37864028901 was
  ubuntu ✓ / macos ✓ / windows ✗ — i.e. a pre-existing main condition, consistent with "Windows experimental / CLI winget held", not a
  regression of this merge. The organizer's two cancellations (37909539566/37909535653 at `8a939761`, 37910628822/37910623311 at `cf1c6398`)
  were metadata-only superseded runs (recorded `product_diff_from_b89e2ab: []`), explicitly logged "Cancellation is not a pass".

## 2. Homebrew (tap feci/parley @ 6964b1188eb0ed81d07f251a4ca8b3852f43dad6)

- Tap HEAD verified both ways: local `git -C $(brew --repository feci/parley) rev-parse HEAD` = `6964b118…` **and** remote `main` via API — identical.
- Both formulae pin source-tag archives. I downloaded both tag archives fresh and hashed them:
  CLI `v1.53.0.tar.gz` → `d47740d8ac528f7fe98ee98d6a0042ba6cf4fcb0c3c93766cc32fd91fcf74a3f`, skill `v2.17.0.tar.gz` →
  `00d25978319e31a91f249547e4cfd67102e59bc5b6a1ee94056222ce6bed41e2` — exactly the in-formula digests.
- Installed state re-checked live: `parley --version` → `parley 1.53.0`; `parley-deck-skill --version` → `2.17.0`;
  `brew list --versions` → `parley-deck-cli … 1.53.0`, `parley-deck-skill … 2.17.0` (matches `homebrew-installed-verification.json`).
- Logs corroborate the narrative: initial unqualified `brew upgrade` exit 1 — "Formulae found in multiple taps" (duplicate tap
  `feci/release-validation-20260918`); qualified `brew upgrade feci/parley/parley-deck-cli feci/parley/parley-deck-skill` succeeded
  (1.52.0→1.53.0, 2.16.0→2.17.0, both built from source); `brew audit` exit 0; `brew style` — "2 files inspected, no offenses detected";
  `brew test` exit 0 for both formulae; test log sha256 re-hashed = `c1ec900e60b12c62e5590c4ce7442bfac34aefe9066bf83acf34b7236117a255` (matches record).
- No tap-trust change; unrelated formulae untouched.

## 3. Runtime targets (15 managed + 4 named generic; 114 SKILL.md)

- **114/114 re-hashed by me from disk**: every path in `runtime-verification.json` exists, and each hash equals both the recorded value **and**
  the packed-payload hash I extracted myself from the hosted tgz. 0 missing, 0 mismatches.
- Target decomposition from the install log: exactly **15 managed targets** (codex, claude, agy→gemini plugins, gemini extensions, hermes, qwen,
  codebuddy, goose, kimi, droid→`.factory`, vibe, cursor, opencode, aionrs, zcode), each `ok: true`; plus the **4 named generic destinations**
  (`~/.hermes/profiles/{ldx,librade,testprofile}/skills/parley-deck`, `~/.config/opencode/skills/parley-deck`), each installed ok.
- Doctor evidence: `runtime-all-doctor.json` ok=true, 15/15 targets `valid`; each of the four generic doctor logs ok=true, `valid`,
  markers at version 2.17.0 / source npm:parley-deck-skill@2.17.0.
- Backup exists: `runtime-skills-before-2.17.0.tar.gz`, sha256 `21815cad4b90752428472ad01791b048e6e00e4b2bce214a5d80ca3028fa04a3` — re-hashed, matches.

## 4. WinGet — skill-only PR 449419

- Live: OPEN, `mergeable`, head `2afe423e66299c133ae5c6b9ae9097c8dbb696a1` (matches brief), branch `feci-parleydeckskill-2.17.0`, author feci, created 09:32:24Z.
- Live diff (`gh pr diff --patch`): **exactly three added files, 48 insertions, 0 deletions** — `Feci.ParleyDeckSkill.installer.yaml`,
  `.locale.en-US.yaml`, `.yaml` (version) under `manifests/f/Feci/ParleyDeckSkill/2.17.0/`, ManifestVersion 1.12.0. Skill-only; no CLI manifest.
- InstallerSha256 `0BD51B98…` (x64) and `E2C20DC7…` (arm64) equal my independently computed hashes of the hosted v2.17.0 Windows release binaries;
  InstallerUrls point at those exact release assets.
- Schema evidence: `winget-validation.json` 3/3 pass against aka.ms 1.12.0 installer/defaultLocale/version schemas; `$schema` headers present in the live manifests.
- Upstream checks at 09:42 UTC: **01–07 SUCCESS** (Pull Request Validation, Manifest Validation, URLs Validation, URL Domain Validation, Manifest
  Policy Validation, Catalog Content Verification, Installers Scan) + `license/cla` SUCCESS; 08 Installation Validation in progress; 09/10 queued.
  Upstream merge is not an AC14 close criterion; CLI WinGet stays held.
- Limitation (disclosed in evidence, holds): native `winget validate`/`install` is unavailable on this Mac.

## 5. Core staging + owner-only actions

- Staged files re-hashed: `COOPERATION-2.16.0.md` = `8c6b95b611901d3d1460ccb8554e26adac67f58d82ea028713791116a3e65a63`;
  `COOPERATION-2.17.0.md` = `6073c311b592653528bc675056b04b8454b491ab2702b37c298402ba2cb87957` (124,976 bytes) — both match `core-staging.json`.
- **Independent reconstruction**: applied the 14 hunks of `source-context/protocol-hunks.json` (sha256 `54d0dd07…`, matches record) to the 2.16.0
  base — each `old` occurs exactly once, substitutions in order — result **byte-identical** to the staged 2.17.0 file, SHA256 `6073c311…` exact.
- Registry/core cross-check (live): npm `dist-tags.latest` = **2.13.0**, no 2.15/2.16/2.17 published; global core installed = **{2.10.0, 2.13.0}**.
  This matches `owner-actions.json`, whose six owner commands (npm + core for 2.15.0 → 2.16.0 → 2.17.0) are therefore all correctly
  still-missing, listed in the required order, `commands_executed: false`. **I executed none of them.**
- The referenced 2.15/2.16 artifacts exist at their recorded paths and re-hash to their recorded values (tgz `64c274a8…`/`80f3172a…`;
  staged md `0d81fd80…`/`8c6b95b6…`).

## 6. Closed-artifacts integrity

All nine artifacts in `closed-artifacts-integrity.json` re-hashed and matching (FINAL.md `c44aac05…`, IMPLEMENTATION.md `5af3e100…`,
goal-check-zcode-1.md `87a9bae3…`, consensus.md `e85d4e0f…`, review/consensus.md `047e0317…`, review/consensus-cycle-01.md `1a95c07a…`,
review/consensus-cycle-01-proposed.md `5e9bacc…`, review/round-01/zcode-1.md `beb3ad15…`, review/round-02/zcode-1.md `49a19f08…`).
Everything frozen; nothing modified by this verification.

## Material gaps and nonblocking limitations (reported freely)

1. **CLI merge-CI 37910668038 had not concluded** by my final poll (09:42:31Z: ubuntu ✓, macos ✓, windows in-Test). It is not counted as passed.
   Mitigating context: windows-latest is a pre-existing failing job on main (pre-idea baseline 85c5a8f failed windows-only), ubuntu+macos are green
   on the merge commit, and the CLI windows channel is expressly experimental/held. If windows fails at completion it is continuation of a known
   main condition, not a delivery regression — but the owner should watch the run to completion.
2. **WinGet upstream checks 08–10 incomplete** at last poll (01–07 + CLA green). Upstream merge remains pending and is outside AC14.
3. **Native winget validation impossible on macOS** — upstream pipeline relied upon instead (now green through step 07).
4. **npm and global-core publication outstanding by design** — six owner-only commands (2.15, 2.16, 2.17 × npm/core), none executed, correctly ordered.
5. Minor: `hosted-asset-verification.json` has no locally-prepared comparison for the two workflow-built skill Windows exes (`prepared_asset_matches: null`);
   closed by my direct hash check against live digests and the WinGet manifest.

---

## Verdict

All authorized channels delivered and independently re-verified: merges and tags at the exact stated SHAs; all 11 hosted assets with live-API /
local-hash / manifest three-way agreement; clean-source VCS proof in the binaries; tarball byte-identical to merged source; tap, formulae,
installed versions and all brew checks corroborated; 114/114 runtime SKILL.md matching the packed payload with valid doctor markers and an intact
backup; a correct skill-only WinGet PR with verified hashes and green upstream checks through step 07; staged core 2.17.0 exactly reconstructible
from 2.16.0 + the 14 reviewed hunks; owner-only commands accurately enumerated, none executed; closed artifacts frozen and the brief byte-identical.
The open items above are nonblocking for AC14 and correctly excluded from its scope.

**CHANNEL-VERIFICATION: PASS**

*(PASS covers authorized channel delivery only — GitHub releases, Homebrew, runtime installs, WinGet PR submission and its submitted content —
and expressly excludes owner npm/global-core publication and the upstream WinGet merge, which remain owner/upstream actions. Process PID,
duration and actual exit are recorded by the organizer; not fabricated here.)*
