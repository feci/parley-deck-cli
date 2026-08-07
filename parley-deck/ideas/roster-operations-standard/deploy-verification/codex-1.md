---
idea: roster-operations-standard
phase: deploy verification
agent: codex-1
date: 2026-08-06
verdict: VERIFIED
---

# Deploy verification — codex-1

## Verdict

**PRIMARY — CONFIRMED:** **VERIFIED.** The independently observed deployment state matches every claim under verification: both Homebrew formulae are published at the claimed versions with real tag refs and matching downloaded-tarball hashes; both installed commands resolve into linked Homebrew Cellar kegs; there is no npm-global skill copy; npm `latest` is 2.5.1; both GitHub releases exist with exactly two Windows assets; both winget PRs are open; all four winget URLs return HTTP 200 and their downloaded bytes match both the uppercase manifest hashes and GitHub asset digests; all seven runtime installs are valid 2.5.1 copies with the ratified content; and the fleet result is 35 synced decks out of 36, with only the read-only `ecb-meeting-2026.05` deck not synced.

**PRIMARY — CONFIRMED:** The local shell sandbox could not resolve `github.com`, so external primary checks used direct GitHub/npm HTTP requests and browser-managed downloads in an isolated Chromium task space. Downloaded files were hashed locally with `/usr/bin/shasum` and deleted afterwards.

```text
$ git -C '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley' ls-remote origin refs/heads/main
fatal: unable to access 'https://github.com/feci/homebrew-parley.git/': Could not resolve host: github.com
```

## Channel-by-channel (command + quoted output for each)

### Homebrew tap and source archives

**PRIMARY — CONFIRMED:** The tap contains exactly the two claimed formulae, and its local `main` is clean and aligned with the published GitHub `main` commit.

```text
$ find '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula' -maxdepth 1 -type f -name '*.rb' -print | LC_ALL=C sort
/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/parley-deck-cli.rb
/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/parley-deck-skill.rb
$ printf 'formula_count='; find '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula' -maxdepth 1 -type f -name '*.rb' -print | wc -l | tr -d ' '
formula_count=2

$ git -C '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley' status --short --branch
## main...origin/main
$ git -C '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley' rev-parse HEAD
573cb27a337a666644d6c871a3abb9c3f093978e

$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
const result = await js(String.raw`(async () => {
  const url='https://api.github.com/repos/feci/homebrew-parley/commits/main'
  const r=await fetch(url,{headers:{Accept:'application/vnd.github+json'}})
  const j=await r.json()
  return {request:'GET '+url,status:r.status,sha:j.sha,commitMessage:j.commit&&j.commit.message}
})()`)
cliLog(JSON.stringify(result,null,2))
EOF
{
  "request": "GET https://api.github.com/repos/feci/homebrew-parley/commits/main",
  "status": 200,
  "sha": "573cb27a337a666644d6c871a3abb9c3f093978e",
  "commitMessage": "parley-deck-skill 2.5.1\n\nCo-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>"
}
```

**PRIMARY — CONFIRMED:** The CLI formula declares v1.41.0 and SHA-256 `381b...4cab`; the skill formula declares v2.5.1 and SHA-256 `484f...02bf`.

```text
$ for f in '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/parley-deck-cli.rb' '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/parley-deck-skill.rb'; do echo "FILE $f"; sed -n '1,120p' "$f"; done
FILE /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/parley-deck-cli.rb
class ParleyDeckCli < Formula
  url "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.41.0.tar.gz"
  sha256 "381b4b386c899f5c9cd4a7c907c32409c110b7ed1d808aef89d3a2017f524cab"
FILE /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/parley-deck-skill.rb
class ParleyDeckSkill < Formula
  url "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.5.1.tar.gz"
  sha256 "484f03129751ad0c297efb4b8096b3ccb03cb25278355e767dad7ed42f7202bf"
```

**PRIMARY — CONFIRMED:** Both formula URLs name tag refs that exist in the authoritative repositories.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
const result = await js(String.raw`(async () => {
  const refs=[['feci/parley-deck-cli','v1.41.0'],['feci/parley-deck-skill','v2.5.1']]
  const out=[]
  for (const [repo,tag] of refs) {
    const url='https://api.github.com/repos/'+repo+'/git/ref/tags/'+tag
    const r=await fetch(url,{headers:{Accept:'application/vnd.github+json'}})
    const j=await r.json()
    out.push({request:'GET '+url,status:r.status,ref:j.ref,objectType:j.object&&j.object.type,objectSha:j.object&&j.object.sha})
  }
  return out
})()`)
cliLog(JSON.stringify(result,null,2))
EOF
[
  {
    "request": "GET https://api.github.com/repos/feci/parley-deck-cli/git/ref/tags/v1.41.0",
    "status": 200,
    "ref": "refs/tags/v1.41.0",
    "objectType": "tag",
    "objectSha": "c7c9d3b8c887e7b9b421cb7bdaa98148f80954b6"
  },
  {
    "request": "GET https://api.github.com/repos/feci/parley-deck-skill/git/ref/tags/v2.5.1",
    "status": 200,
    "ref": "refs/tags/v2.5.1",
    "objectType": "tag",
    "objectSha": "04ab4ec8b34783e98d8dd17a20eb7094e5b0e7d7"
  }
]
```

**PRIMARY — CONFIRMED:** I downloaded both formula URLs themselves. Chromium reported completed downloads of 5,018,125 and 472,706 bytes, and local SHA-256 values exactly equal the formula declarations.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
await cdp('Page.setDownloadBehavior', {behavior:'allow', downloadPath:'/private/tmp/parley-deploy-verify.keANOP'})
for (const url of [
  'https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.41.0.tar.gz',
  'https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.5.1.tar.gz'
]) {
  try { await gotoUrl(url) } catch (e) {}
  await wait(5)
}
cliLog(JSON.stringify(await drainEvents(), null, 2))
EOF
"suggestedFilename": "parley-deck-cli-1.41.0.tar.gz"
"totalBytes": 5018125
"receivedBytes": 5018125
"state": "completed"
"suggestedFilename": "parley-deck-skill-2.5.1.tar.gz"
"totalBytes": 472706
"receivedBytes": 472706
"state": "completed"

$ find /private/tmp/parley-deploy-verify.keANOP -maxdepth 1 -type f -print -exec ls -l {} \; -exec shasum -a 256 {} \;
/private/tmp/parley-deploy-verify.keANOP/parley-deck-cli-1.41.0.tar.gz
-rw-r--r--@ 1 tomasfecko staff 5018125 Aug 6 23:29 /private/tmp/parley-deploy-verify.keANOP/parley-deck-cli-1.41.0.tar.gz
381b4b386c899f5c9cd4a7c907c32409c110b7ed1d808aef89d3a2017f524cab  /private/tmp/parley-deploy-verify.keANOP/parley-deck-cli-1.41.0.tar.gz
/private/tmp/parley-deploy-verify.keANOP/parley-deck-skill-2.5.1.tar.gz
-rw-r--r--@ 1 tomasfecko staff 472706 Aug 6 23:29 /private/tmp/parley-deploy-verify.keANOP/parley-deck-skill-2.5.1.tar.gz
484f03129751ad0c297efb4b8096b3ccb03cb25278355e767dad7ed42f7202bf  /private/tmp/parley-deploy-verify.keANOP/parley-deck-skill-2.5.1.tar.gz
```

### Installed Homebrew commands and npm-global shadowing

**PRIMARY — CONFIRMED:** `parley-deck-skill` has one PATH candidate. The `/opt/homebrew/bin` link enters the 2.5.1 keg, and its fully resolved script is inside that keg's `libexec`; it does not resolve into a working tree or npm `node_modules`. The CLI likewise resolves into the 1.41.0 keg.

```text
$ which -a parley-deck-skill
/opt/homebrew/bin/parley-deck-skill
$ ls -la "$(which parley-deck-skill)"
lrwxr-xr-x  1 tomasfecko  admin  55 Aug  6 23:19 /opt/homebrew/bin/parley-deck-skill -> ../Cellar/parley-deck-skill/2.5.1/bin/parley-deck-skill
$ realpath "$(which parley-deck-skill)"
/opt/homebrew/Cellar/parley-deck-skill/2.5.1/libexec/bin/parley-deck-skill.js

$ which -a parley
/opt/homebrew/bin/parley
$ ls -la "$(which parley)"
lrwxr-xr-x  1 tomasfecko  admin  43 Aug  6 22:38 /opt/homebrew/bin/parley -> ../Cellar/parley-deck-cli/1.41.0/bin/parley
$ realpath "$(which parley)"
/opt/homebrew/Cellar/parley-deck-cli/1.41.0/bin/parley
```

**PRIMARY — CONFIRMED:** Homebrew reports both requested versions installed, linked, and not outdated. The tap commit used by both formulae is the same published commit verified above.

```text
$ HOMEBREW_NO_AUTO_UPDATE=1 brew info --json=v2 parley-deck-cli parley-deck-skill | jq '.formulae[] | {full_name,stable:.versions.stable,url:.urls.stable.url,checksum:.urls.stable.checksum,installed:[.installed[].version],linked_keg,outdated,tap_git_head}'
{
  "full_name": "feci/parley/parley-deck-cli",
  "stable": "1.41.0",
  "url": "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.41.0.tar.gz",
  "checksum": "381b4b386c899f5c9cd4a7c907c32409c110b7ed1d808aef89d3a2017f524cab",
  "installed": ["1.41.0"],
  "linked_keg": "1.41.0",
  "outdated": false,
  "tap_git_head": "573cb27a337a666644d6c871a3abb9c3f093978e"
}
{
  "full_name": "feci/parley/parley-deck-skill",
  "stable": "2.5.1",
  "url": "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.5.1.tar.gz",
  "checksum": "484f03129751ad0c297efb4b8096b3ccb03cb25278355e767dad7ed42f7202bf",
  "installed": ["2.5.1"],
  "linked_keg": "2.5.1",
  "outdated": false,
  "tap_git_head": "573cb27a337a666644d6c871a3abb9c3f093978e"
}
```

**PRIMARY — CONFIRMED:** The active npm global prefix contains no `parley-deck-skill` package and no npm `.bin` link. The only prefix-level command is the Cellar link already shown.

```text
$ npm prefix -g
/opt/homebrew
$ npm root -g
/opt/homebrew/lib/node_modules
$ npm ls -g --depth=0 parley-deck-skill; printf 'exit=%s\n' "$?"
/opt/homebrew/lib
└── (empty)
exit=1
$ find /opt/homebrew/lib/node_modules -maxdepth 2 \( -name 'parley-deck-skill' -o -name 'parley-deck-skill.js' \) -print
$ test -e /opt/homebrew/lib/node_modules/.bin/parley-deck-skill || echo 'ABSENT /opt/homebrew/lib/node_modules/.bin/parley-deck-skill'
ABSENT /opt/homebrew/lib/node_modules/.bin/parley-deck-skill
$ test -e /opt/homebrew/lib/node_modules/parley-deck-skill || echo 'ABSENT /opt/homebrew/lib/node_modules/parley-deck-skill'
ABSENT /opt/homebrew/lib/node_modules/parley-deck-skill
```

**PRIMARY — CONFIRMED:** `brew doctor` produced no unlinked-keg warning.

```text
$ doctor_output=$(HOMEBREW_NO_AUTO_UPDATE=1 brew doctor 2>&1); if printf '%s\n' "$doctor_output" | /usr/bin/grep -Ei 'unlinked|keg'; then :; else echo 'NO unlinked/keg warning in brew doctor output'; fi
NO unlinked/keg warning in brew doctor output
```

### npm distribution tag

**PRIMARY — CONFIRMED:** The registry's authoritative dist-tag endpoint maps `latest` to 2.5.1.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
const url='https://registry.npmjs.org/-/package/parley-deck-skill/dist-tags'
const body=await serverFetch(url)
cliLog('GET '+url)
cliLog(body)
EOF
GET https://registry.npmjs.org/-/package/parley-deck-skill/dist-tags
{"latest":"2.5.1"}
```

### GitHub releases and asset digests

**PRIMARY — CONFIRMED:** Both release-tag APIs return HTTP 200 for non-draft, non-prerelease releases. Each release has exactly two assets, both Windows executables, with GitHub-computed SHA-256 digests.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
const result = await js(String.raw`(async () => {
  const specs=[['feci/parley-deck-cli','v1.41.0'],['feci/parley-deck-skill','v2.5.1']]
  const out=[]
  for (const [repo,tag] of specs) {
    const url='https://api.github.com/repos/'+repo+'/releases/tags/'+tag
    const r=await fetch(url,{headers:{Accept:'application/vnd.github+json'}})
    const j=await r.json()
    out.push({request:'GET '+url,status:r.status,tag_name:j.tag_name,draft:j.draft,prerelease:j.prerelease,asset_count:j.assets.length,assets:j.assets.map(a=>({name:a.name,size:a.size,digest:a.digest,url:a.browser_download_url}))})
  }
  return out
})()`)
cliLog(JSON.stringify(result,null,2))
EOF
[
  {
    "request": "GET https://api.github.com/repos/feci/parley-deck-cli/releases/tags/v1.41.0",
    "status": 200,
    "tag_name": "v1.41.0",
    "draft": false,
    "prerelease": false,
    "asset_count": 2,
    "assets": [
      {"name":"parley-v1.41.0-windows-arm64.exe","size":6648832,"digest":"sha256:b4d63ee1871c5c443058eb4254836df907ae89149f2d7355a67fe9f4012ec314"},
      {"name":"parley-v1.41.0-windows-x64.exe","size":7349760,"digest":"sha256:01472d76fe8c80ee4df9be1e06f4adfea83def065db41a25a546b8915ea03d9a"}
    ]
  },
  {
    "request": "GET https://api.github.com/repos/feci/parley-deck-skill/releases/tags/v2.5.1",
    "status": 200,
    "tag_name": "v2.5.1",
    "draft": false,
    "prerelease": false,
    "asset_count": 2,
    "assets": [
      {"name":"parley-deck-skill-v2.5.1-windows-arm64.exe","size":85804331,"digest":"sha256:aa4612cfbc9be0267f7dcc0b251337ed4a64bcc0ee7c1a58826d65ac3e67109b"},
      {"name":"parley-deck-skill-v2.5.1-windows-x64.exe","size":91391275,"digest":"sha256:18e3145d15cdf2006cd1c0eaaf42c5645962d73299db667f558899567f4e4708"}
    ]
  }
]
```

### winget PRs, fork manifests, live URLs, and downloaded hashes

**PRIMARY — CONFIRMED:** PRs 413351 and 413352 are both open, non-draft, unmerged PRs against `microsoft/winget-pkgs:master`, sourced from the claimed fork. Their head commits contain the expected three-file manifests at 1.41.0 and 2.5.1.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
const result = await js(String.raw`(async () => {
  const out=[]
  for (const pr of [413351,413352]) {
    const url='https://api.github.com/repos/microsoft/winget-pkgs/pulls/'+pr
    const r=await fetch(url,{headers:{Accept:'application/vnd.github+json'}})
    const j=await r.json()
    const fr=await fetch(url+'/files?per_page=100',{headers:{Accept:'application/vnd.github+json'}})
    const files=await fr.json()
    out.push({status:r.status,number:j.number,state:j.state,draft:j.draft,merged:j.merged,headRepo:j.head.repo.full_name,headRef:j.head.ref,headSha:j.head.sha,baseRef:j.base.ref,files:files.map(f=>f.filename)})
  }
  return out
})()`)
cliLog(JSON.stringify(result,null,2))
EOF
[
  {
    "status": 200,
    "number": 413351,
    "state": "open",
    "draft": false,
    "merged": false,
    "headRepo": "feci/winget-pkgs",
    "headRef": "feci-parley-1.41.0",
    "headSha": "2be32fcb1155fcde648a4239ecffc8292679242b",
    "baseRef": "master",
    "files": [
      "manifests/f/Feci/ParleyDeckCli/1.41.0/Feci.ParleyDeckCli.installer.yaml",
      "manifests/f/Feci/ParleyDeckCli/1.41.0/Feci.ParleyDeckCli.locale.en-US.yaml",
      "manifests/f/Feci/ParleyDeckCli/1.41.0/Feci.ParleyDeckCli.yaml"
    ]
  },
  {
    "status": 200,
    "number": 413352,
    "state": "open",
    "draft": false,
    "merged": false,
    "headRepo": "feci/winget-pkgs",
    "headRef": "feci-skill-2.5.1",
    "headSha": "727f7986ffef24243dedb6e4bbf4430a06c507e3",
    "baseRef": "master",
    "files": [
      "manifests/f/Feci/ParleyDeckSkill/2.5.1/Feci.ParleyDeckSkill.installer.yaml",
      "manifests/f/Feci/ParleyDeckSkill/2.5.1/Feci.ParleyDeckSkill.locale.en-US.yaml",
      "manifests/f/Feci/ParleyDeckSkill/2.5.1/Feci.ParleyDeckSkill.yaml"
    ]
  }
]
```

**PRIMARY — CONFIRMED:** Reading the installer manifests from the fork at those exact head SHAs gives the four expected release URLs and uppercase `InstallerSha256` values.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
const specs=[
  {pr:413351,sha:'2be32fcb1155fcde648a4239ecffc8292679242b',path:'manifests/f/Feci/ParleyDeckCli/1.41.0/Feci.ParleyDeckCli.installer.yaml'},
  {pr:413352,sha:'727f7986ffef24243dedb6e4bbf4430a06c507e3',path:'manifests/f/Feci/ParleyDeckSkill/2.5.1/Feci.ParleyDeckSkill.installer.yaml'}
]
for (const s of specs) {
  const url='https://api.github.com/repos/feci/winget-pkgs/contents/'+s.path+'?ref='+s.sha
  const body=JSON.parse(await serverFetch(url,{headers:{Accept:'application/vnd.github+json'}}))
  const manifest=Buffer.from(body.content.replace(/\n/g,''),'base64').toString('utf8')
  cliLog('PR '+s.pr+' GET '+url)
  cliLog(manifest)
}
EOF
PR 413351 GET https://api.github.com/repos/feci/winget-pkgs/contents/manifests/f/Feci/ParleyDeckCli/1.41.0/Feci.ParleyDeckCli.installer.yaml?ref=2be32fcb1155fcde648a4239ecffc8292679242b
PackageIdentifier: Feci.ParleyDeckCli
PackageVersion: 1.41.0
Installers:
- Architecture: x64
  InstallerUrl: https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-x64.exe
  InstallerSha256: 01472D76FE8C80EE4DF9BE1E06F4ADFEA83DEF065DB41A25A546B8915EA03D9A
- Architecture: arm64
  InstallerUrl: https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-arm64.exe
  InstallerSha256: B4D63EE1871C5C443058EB4254836DF907AE89149F2D7355A67FE9F4012EC314

PR 413352 GET https://api.github.com/repos/feci/winget-pkgs/contents/manifests/f/Feci/ParleyDeckSkill/2.5.1/Feci.ParleyDeckSkill.installer.yaml?ref=727f7986ffef24243dedb6e4bbf4430a06c507e3
PackageIdentifier: Feci.ParleyDeckSkill
PackageVersion: 2.5.1
Installers:
- Architecture: x64
  InstallerUrl: https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-x64.exe
  InstallerSha256: 18E3145D15CDF2006CD1C0EAAF42C5645962D73299DB667F558899567F4E4708
- Architecture: arm64
  InstallerUrl: https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-arm64.exe
  InstallerSha256: AA4612CFBC9BE0267F7DCC0B251337ED4A64BCC0EE7C1A58826D65AC3E67109B
```

**PRIMARY — CONFIRMED:** Real GETs to all four `InstallerUrl` values reached the GitHub release-asset host with HTTP 200 and `application/octet-stream` responses.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
await cdp('Network.enable')
await cdp('Page.setDownloadBehavior', {behavior:'deny'})
for (const url of [
  'https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-x64.exe',
  'https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-arm64.exe',
  'https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-x64.exe',
  'https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-arm64.exe'
]) {
  await drainEvents(); try { await gotoUrl(url) } catch (e) {}; await wait(1)
  const events=await drainEvents()
  const response=events.find(e=>e.method==='Network.responseReceived')
  const download=events.find(e=>e.method==='Page.downloadWillBegin')
  cliLog(JSON.stringify({request:'GET '+url,status:response.params.response.status,finalHost:new URL(response.params.response.url).hostname,mimeType:response.params.response.mimeType,filename:download.params.suggestedFilename}))
}
EOF
{"request":"GET https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-x64.exe","status":200,"finalHost":"release-assets.githubusercontent.com","mimeType":"application/octet-stream","filename":"parley-v1.41.0-windows-x64.exe"}
{"request":"GET https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-arm64.exe","status":200,"finalHost":"release-assets.githubusercontent.com","mimeType":"application/octet-stream","filename":"parley-v1.41.0-windows-arm64.exe"}
{"request":"GET https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-x64.exe","status":200,"finalHost":"release-assets.githubusercontent.com","mimeType":"application/octet-stream","filename":"parley-deck-skill-v2.5.1-windows-x64.exe"}
{"request":"GET https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-arm64.exe","status":200,"finalHost":"release-assets.githubusercontent.com","mimeType":"application/octet-stream","filename":"parley-deck-skill-v2.5.1-windows-arm64.exe"}
```

**PRIMARY — CONFIRMED:** I separately downloaded all four asset bodies and hashed their bytes locally. Every actual SHA-256 equals the fork manifest's uppercase `InstallerSha256`; the same values, lowercased and prefixed `sha256:`, equal GitHub's release digests above.

```text
$ ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace(14)
await openOrReuseTab('https://github.com', { wait: true, timeout: 20 })
await cdp('Page.setDownloadBehavior', {behavior:'allow', downloadPath:'/private/tmp/parley-winget-verify.PF4a9L'})
const downloads = [
  ['https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-x64.exe',5],
  ['https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-arm64.exe',5],
  ['https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-x64.exe',12],
  ['https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-arm64.exe',12]
]
for (const [url,seconds] of downloads) {
  try { await gotoUrl(url) } catch (e) {}
  await wait(seconds)
  await drainEvents()
}
EOF
$ find /private/tmp/parley-winget-verify.PF4a9L -maxdepth 1 -type f -print -exec ls -l {} \; -exec shasum -a 256 {} \;
/private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-arm64.exe
-rw-r--r--@ 1 tomasfecko staff 6648832 Aug 6 23:31 /private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-arm64.exe
b4d63ee1871c5c443058eb4254836df907ae89149f2d7355a67fe9f4012ec314  /private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-arm64.exe
/private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-x64.exe
-rw-r--r--@ 1 tomasfecko staff 91391275 Aug 6 23:32 /private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-x64.exe
18e3145d15cdf2006cd1c0eaaf42c5645962d73299db667f558899567f4e4708  /private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-x64.exe
/private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-x64.exe
-rw-r--r--@ 1 tomasfecko staff 7349760 Aug 6 23:31 /private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-x64.exe
01472d76fe8c80ee4df9be1e06f4adfea83def065db41a25a546b8915ea03d9a  /private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-x64.exe
/private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-arm64.exe
-rw-r--r--@ 1 tomasfecko staff 85804331 Aug 6 23:32 /private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-arm64.exe
aa4612cfbc9be0267f7dcc0b251337ed4a64bcc0ee7c1a58826d65ac3e67109b  /private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-arm64.exe

$ check_hash() { file=$1; expected=$2; actual=$(shasum -a 256 "$file" | awk '{print toupper($1)}'); printf '%s\n  declared=%s\n  actual=%s\n  match=%s\n' "$(basename "$file")" "$expected" "$actual" "$( [ "$actual" = "$expected" ] && echo YES || echo NO )"; }
$ check_hash /private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-x64.exe 01472D76FE8C80EE4DF9BE1E06F4ADFEA83DEF065DB41A25A546B8915EA03D9A
parley-v1.41.0-windows-x64.exe
  declared=01472D76FE8C80EE4DF9BE1E06F4ADFEA83DEF065DB41A25A546B8915EA03D9A
  actual=01472D76FE8C80EE4DF9BE1E06F4ADFEA83DEF065DB41A25A546B8915EA03D9A
  match=YES
$ check_hash /private/tmp/parley-winget-verify.PF4a9L/parley-v1.41.0-windows-arm64.exe B4D63EE1871C5C443058EB4254836DF907AE89149F2D7355A67FE9F4012EC314
parley-v1.41.0-windows-arm64.exe
  declared=B4D63EE1871C5C443058EB4254836DF907AE89149F2D7355A67FE9F4012EC314
  actual=B4D63EE1871C5C443058EB4254836DF907AE89149F2D7355A67FE9F4012EC314
  match=YES
$ check_hash /private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-x64.exe 18E3145D15CDF2006CD1C0EAAF42C5645962D73299DB667F558899567F4E4708
parley-deck-skill-v2.5.1-windows-x64.exe
  declared=18E3145D15CDF2006CD1C0EAAF42C5645962D73299DB667F558899567F4E4708
  actual=18E3145D15CDF2006CD1C0EAAF42C5645962D73299DB667F558899567F4E4708
  match=YES
$ check_hash /private/tmp/parley-winget-verify.PF4a9L/parley-deck-skill-v2.5.1-windows-arm64.exe AA4612CFBC9BE0267F7DCC0B251337ED4A64BCC0EE7C1A58826D65AC3E67109B
parley-deck-skill-v2.5.1-windows-arm64.exe
  declared=AA4612CFBC9BE0267F7DCC0B251337ED4A64BCC0EE7C1A58826D65AC3E67109B
  actual=AA4612CFBC9BE0267F7DCC0B251337ED4A64BCC0EE7C1A58826D65AC3E67109B
  match=YES
```

### Seven runtime installs and content-level verification

**PRIMARY — CONFIRMED:** The installed Cellar executable reports seven detected, valid runtime installs. Every marker is 2.5.1, every runtime version matches the installer, and no required file is missing.

```text
$ parley-deck-skill status --target all --project . --json | jq '{ok,installer,runtimeInstallCount:(.runtimeInstalls|length),runtimeInstalls:[.runtimeInstalls[]|{target,dest,detected,status,version,markerVersion:.marker.version,versionMatchesInstaller,missing}],parleyCli,project:{metadataStatus:.project.metadataStatus,metadataMatchesProtocol:.project.metadataMatchesProtocol},compatibility}'
{
  "ok": true,
  "installer": {
    "name": "parley-deck-skill",
    "version": "2.5.1",
    "source": "npm:parley-deck-skill@2.5.1",
    "packageRoot": "/opt/homebrew/Cellar/parley-deck-skill/2.5.1/libexec",
    "executable": "/opt/homebrew/bin/parley-deck-skill"
  },
  "runtimeInstallCount": 7,
  "runtimeInstalls": [
    {"target":"codex","dest":"/Users/tomasfecko/.codex/skills/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]},
    {"target":"claude","dest":"/Users/tomasfecko/.claude/skills/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]},
    {"target":"agy","dest":"/Users/tomasfecko/.gemini/config/plugins/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]},
    {"target":"gemini","dest":"/Users/tomasfecko/.gemini/extensions/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]},
    {"target":"hermes","dest":"/Users/tomasfecko/.hermes/skills/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]},
    {"target":"kimi","dest":"/Users/tomasfecko/.kimi-code/skills/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]},
    {"target":"opencode","dest":"/Users/tomasfecko/.opencode/skills/parley-deck","detected":true,"status":"valid","version":"2.5.1","markerVersion":"2.5.1","versionMatchesInstaller":true,"missing":[]}
  ],
  "parleyCli": {"available":true,"version":"parley 1.41.0"},
  "project": {"metadataStatus":"valid","metadataMatchesProtocol":true},
  "compatibility": {"status":"warning","reasons":["project-protocol-differs-from-packaged-reference"]}
}
```

**PRIMARY — CONFIRMED:** This is not a version-string-only result. Each of the seven installed `SKILL.md` files has the same SHA-256 as the Cellar package, contains `Membership is the DECK FILE`, and lacks the forbidden legacy sentence tested below. The required Claude and Codex checks both pass explicitly.

```text
$ check_skill() {
  target=$1; dest=$2; file="$dest/SKILL.md"; marker="$dest/.parley-deck-skill-install.json"
  echo "TARGET=$target"; shasum -a 256 "$file"
  /usr/bin/grep -Fn 'Membership is the DECK FILE' "$file" || echo REQUIRED_PHRASE_ABSENT
  if /usr/bin/grep -Fn 'roster sync` moves it across' "$file"; then echo FORBIDDEN_PHRASE_PRESENT; else echo FORBIDDEN_PHRASE_ABSENT; fi
  jq '{version,source,target,skill}' "$marker"
}
$ check_skill codex /Users/tomasfecko/.codex/skills/parley-deck
TARGET=codex
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.codex/skills/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT
{"version":"2.5.1","source":"npm:parley-deck-skill@2.5.1","target":"codex","skill":"parley-deck"}
$ check_skill claude /Users/tomasfecko/.claude/skills/parley-deck
TARGET=claude
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.claude/skills/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT
{"version":"2.5.1","source":"npm:parley-deck-skill@2.5.1","target":"claude","skill":"parley-deck"}
$ check_skill agy /Users/tomasfecko/.gemini/config/plugins/parley-deck
TARGET=agy
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.gemini/config/plugins/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT
$ check_skill gemini /Users/tomasfecko/.gemini/extensions/parley-deck
TARGET=gemini
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.gemini/extensions/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT
$ check_skill hermes /Users/tomasfecko/.hermes/skills/parley-deck
TARGET=hermes
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.hermes/skills/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT
$ check_skill kimi /Users/tomasfecko/.kimi-code/skills/parley-deck
TARGET=kimi
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.kimi-code/skills/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT
$ check_skill opencode /Users/tomasfecko/.opencode/skills/parley-deck
TARGET=opencode
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.opencode/skills/parley-deck/SKILL.md
307:**Membership is the DECK FILE.** The machine layer (`~/.parley/agents.toml`) seeds *values* for
FORBIDDEN_PHRASE_ABSENT

$ shasum -a 256 /opt/homebrew/Cellar/parley-deck-skill/2.5.1/libexec/skills/parley-deck/SKILL.md /Users/tomasfecko/.claude/skills/parley-deck/SKILL.md /Users/tomasfecko/.codex/skills/parley-deck/SKILL.md
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /opt/homebrew/Cellar/parley-deck-skill/2.5.1/libexec/skills/parley-deck/SKILL.md
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.claude/skills/parley-deck/SKILL.md
003d034d3775e29dd464c587dd150a19facbcd4eb49a983fed688337971b0d1d  /Users/tomasfecko/.codex/skills/parley-deck/SKILL.md
```

**PRIMARY — CONFIRMED:** The status warning is not payload drift. A direct diff shows the packaged protocol and this source deck differ only in project-specific header and generated roster zones; project metadata itself matches the live protocol.

```text
$ diff -u /opt/homebrew/Cellar/parley-deck-skill/2.5.1/libexec/skills/parley-deck/references/COOPERATION.md parley-deck/COOPERATION.md | sed -n '1,260p'
@@ -1,9 +1,10 @@
-**Workspace:** `<workspace-name>`
+**Workspace:** `parley-deck`
-**Transport:** `<transport-choice>` ...
-**Created:** `<YYYY-MM-DD>` ...
+**Transport:** `github-pr`
+**Created:** 2026-05-09 (initial draft)
+**Protocol synced:** 2026-08-06 — parley-deck-skill 2.5.1 / parley-deck-cli 1.41.0
@@ -131,6 +132,10 @@
+| `claude-1` ... |
+| `codex-1` ... |
+| `hermes-1` ... |
+| `kimi-1` ... |
@@ -143,6 +148,10 @@
+| `claude-1` | `feci` |
+| `codex-1`  | `feci` |
+| `hermes-1` | `feci` |
+| `kimi-1`   | `feci` |
```

### Fleet protocol sync

**PRIMARY — CONFIRMED:** The required `find` scan discovered 36 decks. Testing every file with `/usr/bin/grep -Fq 'NOT authoritative'` found 35 synced decks and exactly one unsynced deck: `./ecb-meeting-2026.05/parley-deck/COOPERATION.md`.

```text
$ /bin/bash -c '
total=0; synced=0; unsynced=0
while IFS= read -r file; do
  total=$((total + 1))
  if /usr/bin/grep -Fq "NOT authoritative" "$file"; then
    synced=$((synced + 1)); printf "SYNCED %s\n" "$file"
  else
    unsynced=$((unsynced + 1)); printf "NOT_SYNCED %s\n" "$file"
  fi
done < <(find . -maxdepth 4 -type f -path "*/parley-deck/COOPERATION.md" -print | LC_ALL=C sort)
printf "TOTAL=%d SYNCED=%d NOT_SYNCED=%d\n" "$total" "$synced" "$unsynced"
'
SYNCED ./BYTE/parley-deck/COOPERATION.md
SYNCED ./Finance/parley-deck/COOPERATION.md
SYNCED ./IGBCE/parley-deck/COOPERATION.md
SYNCED ./IHK_PFALZ/parley-deck/COOPERATION.md
SYNCED ./SU-Group-Prompt_library/parley-deck/COOPERATION.md
SYNCED ./adito-outlook-plugin/parley-deck/COOPERATION.md
SYNCED ./aditoLeads/parley-deck/COOPERATION.md
SYNCED ./adito_jvm _issue/parley-deck/COOPERATION.md
SYNCED ./ai_prezz/parley-deck/COOPERATION.md
SYNCED ./altfins/altfins-patterns/parley-deck/COOPERATION.md
SYNCED ./altfins/parley-deck/COOPERATION.md
SYNCED ./auftra/parley-deck/COOPERATION.md
SYNCED ./cms/parley-deck/COOPERATION.md
SYNCED ./design-mail/design-mail-fe/parley-deck/COOPERATION.md
SYNCED ./design-mail/design-mail/parley-deck/COOPERATION.md
SYNCED ./design-mail/parley-deck/COOPERATION.md
SYNCED ./ecb-ai-prezz/parley-deck/COOPERATION.md
SYNCED ./ecb-api/parley-deck/COOPERATION.md
NOT_SYNCED ./ecb-meeting-2026.05/parley-deck/COOPERATION.md
SYNCED ./igm-app/parley-deck/COOPERATION.md
SYNCED ./ldx-wt-mail-fixups/parley-deck/COOPERATION.md
SYNCED ./ldx/parley-deck/COOPERATION.md
SYNCED ./librade-algoTrader/parley-deck/COOPERATION.md
SYNCED ./lustrator/parley-deck/COOPERATION.md
SYNCED ./millenniumProblems/parley-deck/COOPERATION.md
SYNCED ./paritaetische/parley-deck/COOPERATION.md
SYNCED ./parley-deck/parley-deck-cli/parley-deck/COOPERATION.md
SYNCED ./parley-deck/wt-editor-composer/parley-deck/COOPERATION.md
SYNCED ./parley-deck/wt-learn-playbooks/parley-deck/COOPERATION.md
SYNCED ./parley-deck/wt-roster-presets/parley-deck/COOPERATION.md
SYNCED ./parley-deck/wt-round-summary/parley-deck/COOPERATION.md
SYNCED ./rev-kimi-scratch/parley-deck/COOPERATION.md
SYNCED ./scaleup/parley-deck/COOPERATION.md
SYNCED ./scaleup/scaleup-report/parley-deck/COOPERATION.md
SYNCED ./servers/parley-deck/COOPERATION.md
SYNCED ./zeroTrust/parley-deck/COOPERATION.md
TOTAL=36 SYNCED=35 NOT_SYNCED=1
```

**PRIMARY — CONFIRMED:** The one unsynced file is read-only at the filesystem level (`-r--------`) and is not writable by this user, consistent with the claimed intentional skip.

```text
$ ls -lOe@ '/Volumes/My Shared Files/AI_WORKSPACE/ecb-meeting-2026.05/parley-deck/COOPERATION.md'
-r--------@ 1 tomasfecko staff - 46510 May 12 09:22 /Volumes/My Shared Files/AI_WORKSPACE/ecb-meeting-2026.05/parley-deck/COOPERATION.md
$ stat -f 'mode=%Sp uid=%u gid=%g flags=%Sf file=%N' '/Volumes/My Shared Files/AI_WORKSPACE/ecb-meeting-2026.05/parley-deck/COOPERATION.md'
mode=-r-------- uid=501 gid=20 flags=- file=/Volumes/My Shared Files/AI_WORKSPACE/ecb-meeting-2026.05/parley-deck/COOPERATION.md
$ if [ -w '/Volumes/My Shared Files/AI_WORKSPACE/ecb-meeting-2026.05/parley-deck/COOPERATION.md' ]; then echo WRITABLE; else echo NOT_WRITABLE; fi
NOT_WRITABLE
```

## Discrepancies

none
