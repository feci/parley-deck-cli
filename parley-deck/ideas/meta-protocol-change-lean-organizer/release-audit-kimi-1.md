# Release Audit — kimi-1 (independent, post-deploy)

Idea: meta-protocol-change-lean-organizer
Release: CLI 1.49.0 (tag commit 06e563e) / skill 2.13.0 (tag commit 8161e5e)
Auditor: kimi-1. Did not perform deployment. All checks re-executed from real artifacts;
my own downloads/hashes in /tmp/kimi1-audit. Read-only except this file.
Audit time: 2026-09-24 ~12:45 local.

## Per-channel verdicts

| Channel | Verdict |
|---|---|
| GitHub CLI release v1.49.0 (6 assets) | PASS |
| GitHub skill release v2.13.0 (5 assets, final Windows) | PASS |
| Homebrew tap feci/homebrew-parley @697901e (both formulae + installed state) | PASS |
| Installed binaries version smoke (Cellar + downloaded release assets) | PASS |
| Runtime SKILL.md propagation, 12 runtimes | PASS |
| winget skill PR microsoft/winget-pkgs#440360 | PASS (open, hashes correct) |
| npm parley-deck-skill 2.13.0 | PENDING (by design — owner reauth) |
| winget CLI 1.49.0 PR | PENDING (by design — held for CLI CI incident) |
| CLI CI run 35987916696 | IN FLIGHT (Claude-owned; reported, not audited as product gate) |

## 1. GitHub releases (assets vs live API digests)

Method: `gh api repos/feci/parley-deck-{cli,skill}/releases/tags/<tag>` for live digests,
then `curl -sSL -O` every asset URL and `shasum -a 256`. 11/11 downloaded assets match
live API digests exactly. Live API digests are identical to organizer evidence
`cli-github-final.json` / `skill-github-final.json`.

Tags (both lightweight, verified via `git/refs/tags`):
- CLI v1.49.0 -> commit 06e563e8b1fe8094132149b83e14da7be9aae51e (= briefed 06e563e)
- skill v2.13.0 -> commit 8161e5e9d3ecdece25b7ae225e2304d3159be79b (= briefed 8161e5e;
  also the headSha of successful portable-binaries run 35983892815)

CLI (6/6 match):
- darwin-arm64 229d06e5a5a510e9bfa8b45241ec9875dd1ce120158580d07ae57c84c0d63126
- darwin-x64   a699f4293e83681082eb856a0688aa032554a2eb528d2e110a7896be352117b4
- linux-arm64  b72f1bf40dc6501ced1f18074bb9b4dcf7bd9bc5b8c97236c002987b0528c294
- linux-x64    becea6ff9d3bf549bbfaf7286bf2970ea11996eeb4828088b2ad2f674b934a01
- windows-arm64.exe ace3d577fb4a4cea6acdd74c0eece02450dfa6fc941d809b270e8412a5c5b0cd
- windows-x64.exe   1aed9c1f1f65c4b5b6232bde01158891e6c401dcffc741b82949c8113e931767

Skill (5/5 match; Windows rows are the FINAL CI-replaced assets):
- linux-x64    0a6eecd32a2ddf284fc0fb2eda785b134c7bb683ce5a493d4fda6b3b3fc7faa5
- macos-arm64  e4346b7d9599392dc8a2ef3a157a1663f5682858f12e666b49b862c0d363b159
- macos-x64    0d1ad83fe5e21e81a4f6de9e2901fd88a2ca6d5dfe6e31ea57d53d858b37c91c
- windows-arm64.exe 7cfa62dc4dc3821db78f71942a0dcd1bf1afd2093a5f32bfbf63d2dd1f0d840f
- windows-x64.exe   fe9f903cb3491ee129ba411d88e5abd9759c7a3202a17df614d72324253aea23

Note (expected, not a defect): evidence `local-sha256.json` lists pre-replacement local
Windows skill builds (ce5532f4…, 96b70401…). Online final assets supersede them;
all three non-Windows skill assets and all six CLI assets match the local build hashes.

## 2. Homebrew (tap @697901e, formulae vs real tag archives)

Commit 697901e9ac2ae667c45e51e3cf42f9aa36d8c30e "[codex-1] … release CLI 1.49.0 and skill 2.13.0".
Independently re-downloaded both archives and hashed:

- CLI formula url github.com/feci/parley-deck-cli/archive/refs/tags/v1.49.0.tar.gz
  sha256 6ebf32c2d6387e7bf69cff66d15032617d75176c23567c8e107536f860941fc0
  -> my fresh download: identical. Evidence cli-v1.49.0.tar.gz: identical.
- Skill formula url github.com/feci/parley-deck-skill/archive/refs/tags/v2.13.0.tar.gz
  sha256 0d585386219fe1939bd37fd27ba697b25f5b932ae7415657a16833970e44ff44
  -> my fresh download: identical. Evidence skill-v2.13.0.tar.gz: identical.

Installed state: `brew list --versions` = parley-deck-cli 1.49.0, parley-deck-skill 2.13.0;
each Cellar contains only that version.
`which parley-deck-skill` -> /opt/homebrew/bin/parley-deck-skill ->
/opt/homebrew/Cellar/parley-deck-skill/2.13.0/libexec/bin/parley-deck-skill.js
`which parley` -> /opt/homebrew/Cellar/parley-deck-cli/1.49.0/bin/parley

## 3. Binary version smokes (isolated empty HOME, temp only)

- Installed Cellar: `parley version` -> "parley 1.49.0"; `parley-deck-skill --version` -> "2.13.0"
- Downloaded release assets (darwin-arm64 CLI, macos-arm64 skill, hash-verified above):
  same outputs. No writes outside the temp HOME observed.

## 4. Runtime SKILL.md (12 runtimes)

Cellar 2.13.0 source /opt/homebrew/Cellar/parley-deck-skill/2.13.0/libexec/skills/parley-deck/SKILL.md
sha256 db71eca5a4f778cb6cb0d6aee0660474b4dbfc0f8cf19c6585edaf8b0010c03f.
I hashed all 12 installed targets myself: codex, claude, agy (.gemini/config/plugins),
gemini (.gemini/extensions), hermes, qwen, codebuddy, kimi, droid (.factory), vibe,
opencode, zcode -> 12/12 identical to source. Matches installed-skill-hashes.json;
skill-doctor-all.json shows 12/12 "valid".

## 5. winget

- Skill PR #440360 "Update: Feci.ParleyDeckSkill to 2.13.0": OPEN, MERGEABLE,
  head feci/winget-pkgs release/parley-deck-skill-2.13.0 @ c9652fde.
  Manifest manifests/f/Feci/ParleyDeckSkill/2.13.0/Feci.ParleyDeckSkill.installer.yaml:
  x64 URL .../v2.13.0/parley-deck-skill-v2.13.0-windows-x64.exe, InstallerSha256 FE9F903C…AEA23;
  arm64 URL .../v2.13.0/parley-deck-skill-v2.13.0-windows-arm64.exe, InstallerSha256 7CFA62DC…D840F
  -> exact match to final release assets (sec. 1). Awaiting Microsoft review/merge (external).
- CLI: no 1.49.0 PR exists (search: latest CLI winget is merged 1.48.0 #438660).
  Intentionally held pending CLI CI incident remediation — PENDING by design, not a defect.

## 6. npm — PENDING (by design)

`npm view parley-deck-skill`: latest = 2.12.1
(integrity sha512-45s2sww7d7LVtKgIJ9jaOyHvfl1f70mk1r8/luw+4B+8eU2hT5Wwx5L7HVWQLvmtdENMi0Hp4JczEtISuPGmug==);
2.13.0 not on registry at audit time (owner reauth pending; E404/E401 per npm-publish.log).
Preserved tarball parley-deck-skill-2.13.0.tgz: sha256 dcf9c75ca84c3d8f0eb475cf745e4e4bd0973269740e334b0cd8c01d7608d802;
my computed integrity sha512-EHqVIXPveetgOD/xfZZ2iy+iE0ksyXcWXOXCMLZJOUZWHCYetul17mJvTwTFMjM6Iubm9HkNl/y+voVC16zYGg==
(consistent with the publish log's prefix/suffix). Nothing registry-side to compare until publish.

## 7. CLI CI run 35987916696 (report only)

"Tests" on push to main @4e0c0699, created 2026-09-24T10:33:29Z — status in_progress at
last check (12:44 local). Known failures are independently remediated by Claude; this audit
does not gate the released product on it (assets above are hash-verified end-to-end).

## Defects found

None blocking. Observations:
1. npm-publish.log line 222 stores an elided integrity string (literal "[...]").
   Real value recoverable from the preserved tarball (sec. 6); keep full logs next time.
2. Evidence dir retains pre-replacement Windows skill binaries/local hashes alongside the
   final API record; authoritative final hashes exist only in skill-github-final.json.

## Follow-up (for delta audit)

1. Owner: npm reauth + publish 2.13.0 -> delta audit compares registry integrity/tarball
   bytes against sec. 6 values.
2. Claude: finish CLI CI remediation (re-check run 35987916696 conclusion), then open CLI
   winget PR 1.49.0 using hashes 1aed9c1f… (x64) / ace3d577… (arm64); delta audit verifies
   the manifest against these final CLI assets.
3. Track winget #440360 to merge.

## Counts

PASS 6 | FAIL 0 | PENDING (by design) 2 | external in-flight 1.
Next: npm publish by owner; CLI winget PR after Claude's CI remediation; delta audit on request.
