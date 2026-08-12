verdict: PROBLEMS FOUND

# Independent release verification — 2026-08-12

Verifier: `codex-1`

Scope: `parley-deck-cli` v1.44.0 and `parley-deck-skill` v2.8.0. This was a read-only verification; I did not perform the deployment, push, publish, modify a manifest, or change an installed runtime.

## Result by channel

| Channel | Verdict | Result |
| --- | --- | --- |
| 1. git | **VERIFIED** | Both annotated remote tags peel to the exact pushed `main` commit; the remote trees equal local `HEAD` trees; both worktrees were clean before this report was created. |
| 2. GitHub releases | **VERIFIED** | Both non-draft, non-prerelease releases exist. The CLI has six assets and the skill has five. Independently downloaded Windows x64 assets hash to the values in the winget manifests. |
| 3. npm | **VERIFIED** | `parley-deck-skill@2.8.0` exists and the registry's `latest` dist-tag is `2.8.0`. |
| 4. Homebrew | **VERIFIED** | The pushed tap formulae use v1.44.0/v2.8.0; hashes recomputed from the real source tarball responses match; both versions are installed locally. |
| 5. winget | **BROKEN** | The requested version, URL, HTTP, SHA-256, and YAML-colon subchecks pass, but the PR is still open/unmerged and both locale manifests have stale `ReleaseNotesUrl` values. The winget channel is therefore not live and is stale. |
| 6. Runtime skill snapshots | **BROKEN** | All seven managed runtime targets are valid 2.8.0 and all seven snapshots contain the corrected cap wording, but the same status command reports project metadata still at 2.6.0 with `project-metadata-stale` and protocol mismatch warnings. |
| 7. Shipped CLI binary | **VERIFIED** | The installed binary is 1.44.0 and contains the corrected embedded protocol. `protocol check` could not run because this machine has no installed core/pin, so I used the permitted equivalent: initialize from the installed binary and run the repository's embedded-vs-deck drift guard. |

## 1. git — VERIFIED `[PRIMARY]`

The shell's outbound network is blocked, so the requested `git ls-remote` command was executed but could not resolve GitHub:

```text
$ git ls-remote origin refs/heads/main 'refs/tags/v1.44.0' 'refs/tags/v1.44.0^{}'
fatal: unable to access 'https://github.com/feci/parley-deck-cli.git/': Could not resolve host: github.com

$ git ls-remote origin refs/heads/main 'refs/tags/v2.8.0' 'refs/tags/v2.8.0^{}'
fatal: unable to access 'https://github.com/feci/parley-deck-skill/': Could not resolve host: github.com
```

I did not treat a verifier-sandbox DNS failure as proof about the release. I used GitHub's authoritative Git Data API through the isolated browser network and compared its remote objects to local `git rev-parse` output.

### CLI repository

```text
$ cd parley-deck-cli
$ git rev-parse HEAD
059f7e7793979ddbcd856795e4c95209aec4cae8
$ git rev-parse refs/tags/v1.44.0
8c2bd7b0270f23688a449e681b64f7bb41985131
$ git rev-parse 'v1.44.0^{commit}'
059f7e7793979ddbcd856795e4c95209aec4cae8
$ git rev-parse HEAD^{tree}
9ea213db0e568e9823d5bce2f99b5f766a865646
$ git rev-parse 'v1.44.0^{tree}'
9ea213db0e568e9823d5bce2f99b5f766a865646
$ git status --porcelain=v1 --untracked-files=all
(no output)
```

Authoritative remote API results:

```text
GET https://api.github.com/repos/feci/parley-deck-cli/git/ref/heads/main
branchSha: 059f7e7793979ddbcd856795e4c95209aec4cae8

GET https://api.github.com/repos/feci/parley-deck-cli/git/ref/tags/v1.44.0
tagObject: { sha: 8c2bd7b0270f23688a449e681b64f7bb41985131, type: tag }

GET https://api.github.com/repos/feci/parley-deck-cli/git/tags/8c2bd7b0270f23688a449e681b64f7bb41985131
tag: v1.44.0
object: { sha: 059f7e7793979ddbcd856795e4c95209aec4cae8, type: commit }

GET https://api.github.com/repos/feci/parley-deck-cli/git/commits/059f7e7793979ddbcd856795e4c95209aec4cae8
commit: 059f7e7793979ddbcd856795e4c95209aec4cae8
tree:   9ea213db0e568e9823d5bce2f99b5f766a865646
```

The remote tag object, peeled commit, pushed `main` commit, local tag, local `HEAD`, and tree all agree.

### Skill repository

```text
$ cd parley-deck-skill
$ git rev-parse HEAD
4412a8ae8498575b946450f01495f05ef6191ae0
$ git rev-parse refs/tags/v2.8.0
013301594fabe539887c3487d252f9e967f56417
$ git rev-parse 'v2.8.0^{commit}'
4412a8ae8498575b946450f01495f05ef6191ae0
$ git rev-parse HEAD^{tree}
11c70b6f67ae1406a99646626c4d57a89f51a7c0
$ git rev-parse 'v2.8.0^{tree}'
11c70b6f67ae1406a99646626c4d57a89f51a7c0
$ git status --porcelain=v1 --untracked-files=all
(no output)
```

Authoritative remote API results:

```text
GET https://api.github.com/repos/feci/parley-deck-skill/git/ref/heads/main
branchSha: 4412a8ae8498575b946450f01495f05ef6191ae0

GET https://api.github.com/repos/feci/parley-deck-skill/git/ref/tags/v2.8.0
tagObject: { sha: 013301594fabe539887c3487d252f9e967f56417, type: tag }

GET https://api.github.com/repos/feci/parley-deck-skill/git/tags/013301594fabe539887c3487d252f9e967f56417
tag: v2.8.0
object: { sha: 4412a8ae8498575b946450f01495f05ef6191ae0, type: commit }

GET https://api.github.com/repos/feci/parley-deck-skill/git/commits/4412a8ae8498575b946450f01495f05ef6191ae0
commit: 4412a8ae8498575b946450f01495f05ef6191ae0
tree:   11c70b6f67ae1406a99646626c4d57a89f51a7c0
```

The remote tag object, peeled commit, pushed `main` commit, local tag, local `HEAD`, and tree all agree. The tags are annotated but unsigned; tag signing was not part of the stated release gate.

## 2. GitHub releases — VERIFIED `[PRIMARY]`

The requested `gh release view` calls were attempted, but the same shell network restriction prevented them from reaching `api.github.com`:

```text
$ gh release view v1.44.0 --repo feci/parley-deck-cli --json assets
error connecting to api.github.com
check your internet connection or https://githubstatus.com

$ gh release view v2.8.0 --repo feci/parley-deck-skill --json assets
error connecting to api.github.com
check your internet connection or https://githubstatus.com
```

Direct queries to the authoritative release endpoints returned:

```text
GET https://api.github.com/repos/feci/parley-deck-cli/releases/tags/v1.44.0
tag_name: v1.44.0
target_commitish: main
draft: false
prerelease: false
published_at: 2026-08-12T16:20:11Z
assets (6):
  parley-v1.44.0-darwin-arm64
  parley-v1.44.0-darwin-x64
  parley-v1.44.0-linux-arm64
  parley-v1.44.0-linux-x64
  parley-v1.44.0-windows-arm64.exe
  parley-v1.44.0-windows-x64.exe

GET https://api.github.com/repos/feci/parley-deck-skill/releases/tags/v2.8.0
tag_name: v2.8.0
target_commitish: main
draft: false
prerelease: false
published_at: 2026-08-12T16:20:22Z
assets (5):
  parley-deck-skill-v2.8.0-linux-x64
  parley-deck-skill-v2.8.0-macos-arm64
  parley-deck-skill-v2.8.0-macos-x64
  parley-deck-skill-v2.8.0-windows-arm64.exe
  parley-deck-skill-v2.8.0-windows-x64.exe
```

I downloaded one Windows executable from each release to a temporary directory through Chromium's real download path, then ran `shasum` on the files:

```text
download: parley-v1.44.0-windows-x64.exe
totalBytes: 7425024
receivedBytes: 7425024
state: completed

$ shasum -a 256 /private/tmp/.../parley-v1.44.0-windows-x64.exe
85013dd73b15417f25f8c1ed4c76b008e278d2966e1caec54f588bcb77900e47  /private/tmp/.../parley-v1.44.0-windows-x64.exe

download: parley-deck-skill-v2.8.0-windows-x64.exe
totalBytes: 91392746
receivedBytes: 91392746
state: completed

$ shasum -a 256 /private/tmp/.../parley-deck-skill-v2.8.0-windows-x64.exe
7af28f2fc39270143c3cb586ec941e4743fb5cf606c51108cee4852a23ab52ad  /private/tmp/.../parley-deck-skill-v2.8.0-windows-x64.exe
```

These match the winget x64 declarations exactly (case-insensitive hexadecimal comparison). The temporary downloads were deleted after verification.

## 3. npm — VERIFIED `[PRIMARY]`

The local `npm view` invocation did not return registry data under the restricted shell; its only output was an unrelated local configuration warning:

```text
$ npm view parley-deck-skill@2.8.0 version dist-tags --json
npm warn Unknown user config "store-dir". This will stop working in the next major version of npm.
```

I therefore queried npm's authoritative package and version documents directly:

```text
GET https://registry.npmjs.org/parley-deck-skill
distTags: { latest: 2.8.0 }
hasVersion 2.8.0: true
published: 2026-08-12T16:20:26.634Z

GET https://registry.npmjs.org/parley-deck-skill/2.8.0
name: parley-deck-skill
version: 2.8.0
tarball: https://registry.npmjs.org/parley-deck-skill/-/parley-deck-skill-2.8.0.tgz
integrity: sha512-fLWAOiCf1h4o7NWYZcKymzrmJjCVjyzxdln71l+yTtAmHjckH1Qr/YnO3r6t6PftCs02dp2qDtcHn+ripvaF8A==
```

The package exists and 2.8.0 is `latest`.

## 4. Homebrew — VERIFIED `[PRIMARY]`

The pushed tap `main` ref is commit `592e032a0054fc573e4b3a5e0359823da64d3635`, equal to the clean local tap checkout. The formula fields at that pushed commit are:

```text
Formula/parley-deck-cli.rb
url "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.44.0.tar.gz"
sha256 "00ffd0804af011ffcf2d1a8af38782c456e79c42e8ba5a573bc5aa7d574f8bc3"

Formula/parley-deck-skill.rb
url "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.8.0.tar.gz"
sha256 "e163e0c1397f4204c2c960bb0ef99ad2c63df1877187a69499937def2992a3ef"
```

I streamed the bytes from the actual formula URLs through Chromium and hashed the response streams:

```text
https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.44.0.tar.gz
success=true httpStatusCode=200 bytes=5753122
sha256=00ffd0804af011ffcf2d1a8af38782c456e79c42e8ba5a573bc5aa7d574f8bc3

https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.8.0.tar.gz
success=true httpStatusCode=200 bytes=475473
sha256=e163e0c1397f4204c2c960bb0ef99ad2c63df1877187a69499937def2992a3ef
```

The locally installed formulae are current:

```text
$ brew list --versions parley-deck-cli parley-deck-skill
parley-deck-cli 1.44.0
parley-deck-skill 2.8.0

$ brew info --json=v2 feci/parley/parley-deck-cli feci/parley/parley-deck-skill | jq '{formulae: [.formulae[] | {name,full_name,versions,installed:[.installed[]|{version,installed_as_dependency,installed_on_request}]}]}'
parley-deck-cli stable=1.44.0 installed=1.44.0
parley-deck-skill stable=2.8.0 installed=2.8.0
```

## 5. winget — BROKEN `[PRIMARY]`

### PR publication state

```text
GET https://api.github.com/repos/microsoft/winget-pkgs/pulls/416445
state: open
merged: false
head: feci:feci-1.44.0-skill-2.8.0
head.sha: 95fcc7a7a7912b62fffa4cce7e281891184371b1
base.ref: master
```

The local winget checkout independently matches that exact PR head and is clean:

```text
$ git -C ../winget-pkgs branch --show-current
feci-1.44.0-skill-2.8.0
$ git -C ../winget-pkgs rev-parse HEAD
95fcc7a7a7912b62fffa4cce7e281891184371b1
$ git -C ../winget-pkgs status --short
(no output)
```

The PR is not merged, so these versions are not yet published into the winget community index.

### Requested version, URL, HTTP, and SHA checks

The PR installer manifests declare:

```text
Feci.ParleyDeckCli
PackageVersion: 1.44.0
x64 URL:   https://github.com/feci/parley-deck-cli/releases/download/v1.44.0/parley-v1.44.0-windows-x64.exe
x64 SHA:   85013DD73B15417F25F8C1ED4C76B008E278D2966E1CAEC54F588BCB77900E47
arm64 URL: https://github.com/feci/parley-deck-cli/releases/download/v1.44.0/parley-v1.44.0-windows-arm64.exe
arm64 SHA: FDB347A709894C5068E1DF40E2A08BBEB0C83CC6E0A176A17F829370CEB29B71

Feci.ParleyDeckSkill
PackageVersion: 2.8.0
x64 URL:   https://github.com/feci/parley-deck-skill/releases/download/v2.8.0/parley-deck-skill-v2.8.0-windows-x64.exe
x64 SHA:   7AF28F2FC39270143C3CB586EC941E4743FB5CF606C51108CEE4852A23AB52AD
arm64 URL: https://github.com/feci/parley-deck-skill/releases/download/v2.8.0/parley-deck-skill-v2.8.0-windows-arm64.exe
arm64 SHA: 66F7D4EA595EE1475EA2FA82E8AEBE44A9A60D7BEBAEBAFE75F659A5FD876E1A
```

I streamed and hashed all four actual installer assets, not only the two downloaded in the GitHub-release check:

```text
CLI x64:   HTTP 200,  7425024 bytes, sha256 85013dd73b15417f25f8c1ed4c76b008e278d2966e1caec54f588bcb77900e47
CLI arm64: HTTP 200,  6712832 bytes, sha256 fdb347a709894c5068e1df40e2a08bbeb0c83cc6e0a176a17f829370ceb29b71
Skill x64: HTTP 200, 91392746 bytes, sha256 7af28f2fc39270143c3cb586ec941e4743fb5cf606c51108cee4852a23ab52ad
Skill arm64: HTTP 200, 85805802 bytes, sha256 66f7d4ea595ee1475ea2fa82e8aebe44a9a60d7bebaebafe75f659a5fd876e1a
```

All four URLs resolve with HTTP 200 and all four SHA-256 values match.

### YAML colon-space trap

The remote PR files were scanned for an unquoted YAML scalar whose value contains `: `; both locale files returned no matches. The exact PR-head files were then parsed with Ruby/Psych:

```text
$ for manifest_path in ../winget-pkgs/manifests/f/Feci/ParleyDeckCli/1.44.0/*.yaml ../winget-pkgs/manifests/f/Feci/ParleyDeckSkill/2.8.0/*.yaml; do ruby -e 'require "yaml"; YAML.safe_load(File.read(ARGV[0]), permitted_classes: [], aliases: false); puts "YAML_OK #{ARGV[0]}"' "$manifest_path"; done
YAML_OK .../ParleyDeckCli/1.44.0/Feci.ParleyDeckCli.installer.yaml
YAML_OK .../ParleyDeckCli/1.44.0/Feci.ParleyDeckCli.locale.en-US.yaml
YAML_OK .../ParleyDeckCli/1.44.0/Feci.ParleyDeckCli.yaml
YAML_OK .../ParleyDeckSkill/2.8.0/Feci.ParleyDeckSkill.installer.yaml
YAML_OK .../ParleyDeckSkill/2.8.0/Feci.ParleyDeckSkill.locale.en-US.yaml
YAML_OK .../ParleyDeckSkill/2.8.0/Feci.ParleyDeckSkill.yaml

colon-space trap matches in CLI locale:   []
colon-space trap matches in skill locale: []
```

The long `Description` values are double-quoted, so the known trap is absent.

### Defects that make this channel BROKEN

Both locale manifests point release notes at old tags:

```text
manifests/f/Feci/ParleyDeckCli/1.44.0/Feci.ParleyDeckCli.locale.en-US.yaml
ReleaseNotesUrl: https://github.com/feci/parley-deck-cli/releases/tag/v1.43.1

manifests/f/Feci/ParleyDeckSkill/2.8.0/Feci.ParleyDeckSkill.locale.en-US.yaml
ReleaseNotesUrl: https://github.com/feci/parley-deck-skill/releases/tag/v2.5.1
```

They should identify the releases represented by the manifests: v1.44.0 and v2.8.0. Together with `state: open` / `merged: false`, this means the winget channel is stale and not released even though its installer URLs and hashes are correct.

## 6. Runtime skill snapshots — BROKEN `[PRIMARY]`

The requested runtime deployment checks themselves pass. A compact projection of the exact status command shows seven of seven valid installs at 2.8.0:

```text
$ parley-deck-skill status --target all --project . --json | jq '{ok, installer:{version:.installer.version,source:.installer.source}, runtimeInstalls:[.runtimeInstalls[]|{target,status,version,versionMatchesInstaller,dest}], project:{metadataDeckVersion:.project.metadata.deckVersion,metadataSource:.project.metadata.source,metadataStatus:.project.metadataStatus,metadataMatchesProtocol:.project.metadataMatchesProtocol,protocolMatchesPackaged:.project.protocolMatchesPackaged}, parleyCli, compatibility}'
installer: version=2.8.0 source=npm:parley-deck-skill@2.8.0
runtimeInstalls:
  codex:    status=valid version=2.8.0 versionMatchesInstaller=true
  claude:   status=valid version=2.8.0 versionMatchesInstaller=true
  agy:      status=valid version=2.8.0 versionMatchesInstaller=true
  gemini:   status=valid version=2.8.0 versionMatchesInstaller=true
  hermes:   status=valid version=2.8.0 versionMatchesInstaller=true
  kimi:     status=valid version=2.8.0 versionMatchesInstaller=true
  opencode: status=valid version=2.8.0 versionMatchesInstaller=true
```

All seven bundled protocol snapshots have the same hash:

```text
$ shasum -a 256 <each target>/parley-deck/references/COOPERATION.md
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  codex
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  claude
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  agy
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  gemini
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  hermes
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  kimi
254521eb4e9a149152793f2353ad5748a244a446fdf51e83579565c73d9883a8  opencode
```

Wording scan:

```text
codex    | capped-after-round1=1 | cap-5=1 | unbounded=0
claude   | capped-after-round1=1 | cap-5=1 | unbounded=0
agy      | capped-after-round1=1 | cap-5=1 | unbounded=0
gemini   | capped-after-round1=1 | cap-5=1 | unbounded=0
hermes   | capped-after-round1=1 | cap-5=1 | unbounded=0
kimi     | capped-after-round1=1 | cap-5=1 | unbounded=0
opencode | capped-after-round1=1 | cap-5=1 | unbounded=0

216: ... deliberation ... capped at 3 after round 1, then escalate
220: ... deliberation ... cap 5 cycles; `strict_gate` available
```

However, the same primary status command reports stale project state:

```text
project:
  metadataDeckVersion: 2.6.0
  metadataSource: npm:parley-deck-skill@2.6.0
  metadataStatus: valid
  metadataMatchesProtocol: false
  protocolMatchesPackaged: false
compatibility:
  status: warning
  reasons:
    - project-metadata-stale
    - project-protocol-differs-from-packaged-reference
```

The live `parley-deck/COOPERATION.md` header is also stale:

```text
Protocol synced: 2026-08-06 — parley-deck-skill 2.5.1 / parley-deck-cli 1.41.0
```

Therefore the seven runtime copies are correctly deployed, but the runtime/project status channel is not internally clean. I did not run the suggested write command (`parley-deck-skill sync-project --yes`) because this verification was explicitly read-only.

## 7. Shipped binary — VERIFIED `[PRIMARY]`

```text
$ command -v parley
/opt/homebrew/bin/parley
$ parley --version
parley 1.44.0
```

The preferred check could not run because this machine has no global core release and the deck pins none:

```text
$ parley protocol status
core store : /Users/tomasfecko/.parley/protocol/core
installed  : (none)
deck pins  : —
deck view  : .../parley-deck/COOPERATION.md (b2242c9a3785)

$ parley protocol check
protocol check: no core release is installed in /Users/tomasfecko/.parley/protocol/core and this deck pins none — run `parley protocol publish` first
```

The request permits an equivalent. I tested the installed binary's embedded bootstrap protocol directly by running `parley init` in a temporary directory. Its generated `COOPERATION.md` differed from the release source's embedded template only in the selected transport line (`github-pr` in this repository, `local-dir` in the temporary workspace):

```text
$ parley init -dir /private/tmp/.../embedded-check
Initialized Parley Deck workspace at /private/tmp/.../embedded-check/parley-deck

$ diff -u internal/protocol/defaults/COOPERATION.md /private/tmp/.../embedded-check/parley-deck/COOPERATION.md
-**Transport:** `github-pr`
+**Transport:** `local-dir`
```

The repository's semantic drift guard then confirmed that the embedded template matches this deck copy after its defined project-specific normalization:

```text
$ go test ./internal/protocol -run '^TestEmbeddedDefaultMatchesLiveDeck$' -count=1 -v
=== RUN   TestEmbeddedDefaultMatchesLiveDeck
--- PASS: TestEmbeddedDefaultMatchesLiveDeck (0.00s)
PASS
ok   parley-deck-cli/internal/protocol 0.206s
```

Finally, the actual installed executable contains the corrected rows and no `unbounded` match:

```text
$ strings /opt/homebrew/bin/parley | rg -n 'capped at 3 after round 1|cap 5 cycles|unbounded'
1932:| Cross-review rounds (Phase 2) | ... | capped at 3 after round 1, then escalate |
1938:| Fix-up (Phase 8)              | ... | cap 5 cycles; `strict_gate` available |
```

The temporary initialized deck and build cache were deleted after the check.

## Final answer

Yes. The release is not fully consistent across channels:

1. **winget is broken/stale:** PR 416445 is open and unmerged, and both locale manifests link to old release-note tags (`v1.43.1` and `v2.5.1`). Its package versions, installer URLs, HTTP responses, hashes, and YAML syntax are otherwise correct.
2. **Project runtime metadata is stale:** all seven runtime snapshots are valid 2.8.0, but `parley-deck/meta/version.json` remains at 2.6.0 and status reports both metadata and protocol drift warnings; the live protocol header still names skill 2.5.1 / CLI 1.41.0.

No version or hash disagreement was found among the git tags, GitHub release assets, npm package, Homebrew formulae, installed Homebrew versions, the seven installed runtime snapshots, or the installed `parley` binary.
