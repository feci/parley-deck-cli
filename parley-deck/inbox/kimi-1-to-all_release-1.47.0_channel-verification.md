---
from: kimi-1
to: all
idea: protocol-generation-bias
phase: release
blocking: no
date: 2026-08-29
---

# Release 1.47.0 / 2.11.0 — independent channel verification (kimi-1)

Verified on 2026-08-29 against the live systems. Nothing was taken on trust: releases and
tags were read from the live GitHub API/`git ls-remote`, all sha256 values were re-computed
from fresh downloads, and the functional claim was exercised against the installed binary
(`/opt/homebrew/bin/parley`), not the source tree. gh authenticated as `feci`.

## Verified

### 1. GitHub release parley-deck-cli v1.47.0 — VERIFIED

`gh release view v1.47.0 --repo feci/parley-deck-cli --json tagName,isDraft,isPrerelease,publishedAt,assets --jq '{tagName,isDraft,isPrerelease,publishedAt,assetCount:(.assets|length),assets:[.assets[]|{name,size,state}]}'`

```json
{"assetCount":6,"assets":[
 {"name":"parley-v1.47.0-darwin-arm64","size":6977410,"state":"uploaded"},
 {"name":"parley-v1.47.0-darwin-x64","size":7611360,"state":"uploaded"},
 {"name":"parley-v1.47.0-linux-arm64","size":6881440,"state":"uploaded"},
 {"name":"parley-v1.47.0-linux-x64","size":7528608,"state":"uploaded"},
 {"name":"parley-v1.47.0-windows-arm64.exe","size":7040000,"state":"uploaded"},
 {"name":"parley-v1.47.0-windows-x64.exe","size":7810560,"state":"uploaded"}],
 "isDraft":false,"isPrerelease":false,"publishedAt":"2026-08-29T18:50:40Z","tagName":"v1.47.0"}
```

Exactly 6 assets covering darwin/linux/windows × arm64/x64, all `uploaded`, not a draft.
Release notes body is non-trivial (~2.5 KB: §15.6 gate, `runner.ValidateRoundOneArtifact`,
provenance for idea `protocol-generation-bias`).

### 2. GitHub release parley-deck-skill v2.11.0 — VERIFIED

`gh release view v2.11.0 --repo feci/parley-deck-skill --json ...` (same shape):

```json
{"assetCount":5,"assets":[
 {"name":"parley-deck-skill-v2.11.0-linux-x64","size":71195571,"state":"uploaded"},
 {"name":"parley-deck-skill-v2.11.0-macos-arm64","size":65511904,"state":"uploaded"},
 {"name":"parley-deck-skill-v2.11.0-macos-x64","size":69270848,"state":"uploaded"},
 {"name":"parley-deck-skill-v2.11.0-windows-arm64.exe","size":85807779,"state":"uploaded"},
 {"name":"parley-deck-skill-v2.11.0-windows-x64.exe","size":91394723,"state":"uploaded"}],
 "isDraft":false,"isPrerelease":false,"publishedAt":"2026-08-29T18:50:49Z","tagName":"v2.11.0"}
```

Exactly 5 assets, not a draft, non-trivial notes ("Pairs with parley-deck-cli 1.47.0").

### 3. Git tags on origin — VERIFIED

`git ls-remote --tags https://github.com/feci/parley-deck-cli.git` (live origin):

```
a11fe623ddb5dadb3269c453fb41fbef2237b025  refs/tags/protocol-generation-bias-baseline
9d4f45cffcf8bd7ee0ddacfaa38b108c611eb030  refs/tags/protocol-generation-bias-baseline^{}
46ba8919ff14c1eaf06590c21eaef973c0711844  refs/tags/v1.47.0
dd6d2a7d5b03a3e087d5656d468204123637544a  refs/tags/v1.47.0^{}
```

`git ls-remote --tags https://github.com/feci/parley-deck-skill.git`:

```
fbca86a9db50e79744a27bf9f4ac6bb0679673d3  refs/tags/v2.11.0
a9b7bc4aacc2c3f6869e4cd378746e38762cb4f4  refs/tags/v2.11.0^{}
```

- `v1.47.0` on CLI origin (annotated; commit `dd6d2a7d`), `v2.11.0` on skill origin
  (annotated; commit `a9b7bc4a`). Cross-checked via a second API path
  (`gh api repos/.../git/ref/tags/<tag>`) — same SHAs.
- `protocol-generation-bias-baseline` exists — **in the CLI repo only** (annotated; peeled
  commit `9d4f45cf`). Explicit check against the skill repo
  (`git ls-remote --tags ... 'refs/tags/protocol-generation-bias-baseline*'`) returned empty,
  exit 0. The freeze tag is where it is expected to be.
- Local clones match remote exactly: `git -C parley-deck-cli rev-parse v1.47.0
  protocol-generation-bias-baseline` → `46ba8919…` / `a11fe623…`; `git -C parley-deck-skill
  rev-parse v2.11.0` → `fbca86a9…`. Baseline tag is absent from the local skill clone too,
  consistent with the remote.

### 4. Homebrew tap feci/parley — VERIFIED (hashes re-computed, not trusted)

Live repo `feci/homebrew-parley`, default branch `main`, HEAD = `c662107b2f9d`
("parley-deck-cli 1.47.0, parley-deck-skill 2.11.0"). Two formulas in `Formula/`
(**note: the CLI formula is `parley-deck-cli.rb`, there is no `parley.rb`**; it installs a
binary named `parley`).

- `Formula/parley-deck-cli.rb`: url
  `https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.47.0.tar.gz`, declared
  sha256 `3f62c443e3b5f3f4a3899213347967ae13d95336ff17ed373008381330b29664`.
  Re-computed from a fresh `curl -sL` download + `shasum -a 256`:
  `3f62c443e3b5f3f4a3899213347967ae13d95336ff17ed373008381330b29664` → **MATCH**
  (valid gzip, 6,731,948 bytes).
- `Formula/parley-deck-skill.rb`: url
  `https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.11.0.tar.gz`, declared
  sha256 `b5fad20956ecf6785819ed79f6a3e26c47eb79d0db9da84cd0d429adc75449dd`.
  Re-computed: identical → **MATCH** (valid gzip, 480,386 bytes).
- Local clone `parley-deck/homebrew-parley` is in sync with live GitHub (same HEAD, empty
  `git diff HEAD -- Formula/`, blob SHAs match the contents API).

### 5. winget PRs — VERIFIED (installers re-downloaded and re-hashed)

PR #426220 (`gh pr view 426220 --repo microsoft/winget-pkgs --json state,baseRefName,files,...`):
**OPEN**, base **master**, head `c92b9848f282...`. Exactly 3 files, all added, all inside
`manifests/f/Feci/ParleyDeckCli/1.47.0/` (installer / locale.en-US / version yaml) — nothing
outside the one application's folder. Installers (manifest fetched at PR head SHA, then
`curl -sL` + `shasum -a 256`):

- x64 `.../releases/download/v1.47.0/parley-v1.47.0-windows-x64.exe` — HTTP 200,
  declared `BD801939F32A...EE903F77` = actual → **MATCH**
- arm64 `.../parley-v1.47.0-windows-arm64.exe` — HTTP 200,
  declared `12B64DF961...1E7ACBA6` = actual → **MATCH**

PR #426221: **OPEN**, base **master**, head `518b42e8b3bd...`. Exactly 3 files, all added,
all inside `manifests/f/Feci/ParleyDeckSkill/2.11.0/`. Installers:

- x64 `.../releases/download/v2.11.0/parley-deck-skill-v2.11.0-windows-x64.exe` — HTTP 200,
  declared `09CAA3C7CB...D8EC997F` = actual → **MATCH**
- arm64 `.../parley-deck-skill-v2.11.0-windows-arm64.exe` — HTTP 200,
  declared `888EFE35A5...A299B46` = actual → **MATCH**

Both PRs show `mergeStateStatus: BLOCKED` — the normal pre-review state for winget-pkgs
(awaiting moderator validation), not a defect.

### 7. Locally installed binaries — VERIFIED

```
$ which parley parley-deck-skill
/opt/homebrew/bin/parley
/opt/homebrew/bin/parley-deck-skill
$ parley --version
parley 1.47.0
$ parley-deck-skill --version
2.11.0
```

Both report the new versions.

### 8. FUNCTIONAL CLAIM: round-1 artifact without "## Existing alternatives" is rejected at runtime — VERIFIED against the installed binary

There is no standalone `parley validate`; the gate fires inside round-01 execution
(`runner.finalizeExecResult` → `validateArtifactForPhase` → `ValidateRoundOneArtifact`).
It was exercised end-to-end through the real installed CLI: a throwaway deck in
`mktemp -d` (`parley init --dir $TMP`, `PARLEY_HOME=$TMP/...` to redirect all state) with
two fake participant agents (shell scripts, `prompt_mode = "stdin"`) that emit a complete
round-1 artifact with only the alternatives section varied. Real command:
`parley run --dir $TMP --no-tui --no-auto --no-preflight --yes --participants fake-a,fake-b "..."`.
The fake agents exited 0 in all cases — the only variable was the validator's verdict.

- Section **missing entirely** → both participants failed:
  `round-01/fake-a.md is missing a non-empty "## Existing alternatives" section (§15.6a)`,
  run exit code **1**.
- Heading present but body **whitespace-only** → identical rejection, exit code **1**.
- Positive control (same artifact + non-empty alternatives section) →
  `wrote .../round-01/fake-a.md`, run exit code **0**.

Same code path rejects the two non-compliant forms and accepts the compliant one. The claim
holds on the installed 1.47.0 binary.

## Failed or missing

### 6. npm — DID NOT SHIP

```
$ npm view parley-deck-skill version
2.10.0
$ npm whoami
npm error code E401
npm error 401 Unauthorized - GET https://registry.npmjs.org/-/whoami
```

The registry still serves **2.10.0**; the releasing agent's report of being unauthenticated
(E401) is accurate. **This channel did not ship — parley-deck-skill 2.11.0 is not on npm.**
Everything users install via npm is still the old version. All other channels shipped.

## Not checked and why

- **`brew install` / `brew test` of the two formulas** — the tap's url/sha256/tag targets
  were verified and hashes re-computed, but the build stanza and install-time test block were
  not exercised. (Indirect evidence: the locally installed Homebrew binaries report the new
  versions, see channel 7.)
- **`winget validate` / `winget install --manifest` against the PR manifests** — not run
  (requires Windows/winget); URLs, hashes, scope, and PR state were verified instead. Both
  PRs are still `BLOCKED` pending Microsoft's own validation pipeline, which is the normal
  state for fresh winget-pkgs PRs.
- **Auto-drive path behavior for invalid round-1 artifacts** — the loud rejection above is
  the round-01 runner path (where artifacts are ingested). The auto-drive loop
  (`driver.roundComplete`, internal/driver/driver.go:367) treats an invalid round-1 artifact
  as "round incomplete" and waits silently rather than surfacing the error. Not part of the
  release claim; flagged as an observation.
- **GitHub archive byte-stability** — the Homebrew hashes are for GitHub auto-generated
  codeload tarballs. They match today; if GitHub ever changes archive generation the pinned
  hashes would break without any content change. Historical behavior says stable, but it is
  not a hard guarantee.
- **No local installs were (re-)run and no state was mutated**: this verification created no
  files in the workspace other than this report, performed no git mutations, and redirected
  all functional-check state (`--dir`, `PARLEY_HOME`) into a since-deleted temp dir.
