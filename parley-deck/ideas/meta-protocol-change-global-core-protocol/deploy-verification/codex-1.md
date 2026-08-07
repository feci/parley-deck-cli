---
idea: meta-protocol-change-global-core-protocol
agent: codex-1
date: 2026-08-07
verdict: DISCREPANCIES
---
# Deploy verification — codex-1

## Verdict

**DISCREPANCIES.** Every claimed distribution channel passed the independent path, content,
HTTP, and byte-hash checks below. The overall deployment is not fully clean for two independently
observed reasons:

1. `parley protocol --help` exits 2 with `unknown subcommand "--help"`. The top-level
   `parley --help` does expose the protocol group, and `parley protocol status` succeeds, but the
   exact requested group-help command does not work.
2. `parley-deck-skill status --target all --project . --json` reports the current checkout's
   project metadata at 2.5.1 and returns compatibility warnings
   `project-metadata-stale` and `project-protocol-differs-from-packaged-reference`. This is local
   project state, not a failure of the published 2.6.0 installer or its 19 installed skill copies.

Shell DNS was blocked in this verifier's command sandbox. The requested `gh api` commands were
attempted and failed; those failures were not treated as evidence. I used the same GitHub REST
endpoints through the isolated read-only ego-browser transport, downloaded tarballs and all four
Windows assets into memory, and hashed the bytes with SHA-256. No download was written to disk.

## Channel-by-channel (command + quoted output for each)

### 1. Homebrew tap: remote formula bytes, real tags, and real tarball hashes

The local tap checkout is clean (blank `git status --short`) and its HEAD is the same commit as
remote `main`:

```sh
tap_repo=$(HOMEBREW_NO_AUTO_UPDATE=1 brew --repository feci/parley)
git -C "$tap_repo" status --short
git -C "$tap_repo" remote -v
git -C "$tap_repo" branch --show-current
git -C "$tap_repo" rev-parse HEAD
```

Output (exit 0):

```text
[no status output: clean checkout]
origin  https://github.com/feci/homebrew-parley (fetch)
origin  https://github.com/feci/homebrew-parley (push)
main
6769f0b2a053fafc74a1ebafc3ed37bb125961c0
```

I fetched remote `main` and both remote formula files directly:

```sh
ego-browser nodejs <<'EOF'
const task = await useOrCreateTaskSpace('verify remote homebrew tap bytes')
const headers = {'Accept':'application/vnd.github+json','User-Agent':'codex-deploy-verifier'}
const branchRes = await fetch('https://api.github.com/repos/feci/homebrew-parley/branches/main',{headers})
const branch = await branchRes.json()
const out = {apiStatus:branchRes.status,remoteHead:branch.commit.sha,formulae:{}}
for (const file of ['parley-deck-cli.rb','parley-deck-skill.rb']) {
  const url = `https://raw.githubusercontent.com/feci/homebrew-parley/${branch.commit.sha}/Formula/${file}`
  const res = await fetch(url,{headers:{'User-Agent':'codex-deploy-verifier'}})
  const body = await res.text()
  out.formulae[file] = {status:res.status,lines:body.split('\n').filter(line=>/^\s*(url|sha256) /.test(line))}
}
cliLog(JSON.stringify(out,null,2))
EOF
```

Output (exit 0):

```text
apiStatus: 200
remoteHead: 6769f0b2a053fafc74a1ebafc3ed37bb125961c0
parley-deck-cli.rb: status 200
  url "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.42.0.tar.gz"
  sha256 "f7c066ede84bd4ff8082df215f32bfed118e98eb7b620ba1d41a8e9551e97924"
parley-deck-skill.rb: status 200
  url "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.6.0.tar.gz"
  sha256 "37934effd8a06c28dedb77447785b0e62cfe561b9f13711961c68b4ba713a80f"
```

For each formula URL, I queried the tag ref, followed the formula URL, downloaded the returned
tarball into memory, and hashed those bytes:

```sh
ego-browser nodejs <<'EOF'
const checks = [
  ['feci/parley-deck-cli','v1.42.0','https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.42.0.tar.gz','f7c066ede84bd4ff8082df215f32bfed118e98eb7b620ba1d41a8e9551e97924'],
  ['feci/parley-deck-skill','v2.6.0','https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.6.0.tar.gz','37934effd8a06c28dedb77447785b0e62cfe561b9f13711961c68b4ba713a80f']
]
for (const [repo,tag,url,formulaSha] of checks) {
  const tagRes = await fetch(`https://api.github.com/repos/${repo}/git/ref/tags/${tag}`)
  const tagBody = await tagRes.json()
  const tarRes = await fetch(url)
  const bytes = await tarRes.arrayBuffer()
  const digest = await crypto.subtle.digest('SHA-256',bytes)
  const actualSha = [...new Uint8Array(digest)].map(b=>b.toString(16).padStart(2,'0')).join('')
  cliLog(JSON.stringify({repo,tagApiStatus:tagRes.status,ref:tagBody.ref,
    tagObject:tagBody.object,tarballStatus:tarRes.status,bytes:bytes.byteLength,
    formulaSha,actualSha,match:actualSha===formulaSha},null,2))
}
EOF
```

Output (exit 0):

```text
feci/parley-deck-cli
  tagApiStatus: 200
  ref: refs/tags/v1.42.0
  tag object: c8184609a728793fd7f99fbdc68c74f805f122c4 (tag)
  tarballStatus: 200
  bytes: 5233256
  formulaSha: f7c066ede84bd4ff8082df215f32bfed118e98eb7b620ba1d41a8e9551e97924
  actualSha:  f7c066ede84bd4ff8082df215f32bfed118e98eb7b620ba1d41a8e9551e97924
  match: true

feci/parley-deck-skill
  tagApiStatus: 200
  ref: refs/tags/v2.6.0
  tag object: 638f1a48099bbaaca4d31b265f8c7256c9139a56 (tag)
  tarballStatus: 200
  bytes: 474418
  formulaSha: 37934effd8a06c28dedb77447785b0e62cfe561b9f13711961c68b4ba713a80f
  actualSha:  37934effd8a06c28dedb77447785b0e62cfe561b9f13711961c68b4ba713a80f
  match: true
```

Result: **PRIMARY CONFIRMED** — both formulae are present on the remote tap, both tag refs exist,
and each formula SHA-256 equals the bytes served by its tag tarball.

### 2. Installed binaries, Cellar resolution, npm shadowing, and Homebrew link health

```sh
command -v parley
which -a parley
parley --version
command -v parley-deck-skill
which -a parley-deck-skill
parley-deck-skill --version
ls -la $(which parley-deck-skill)
realpath $(which parley-deck-skill)
```

Output (exit 0):

```text
/opt/homebrew/bin/parley
/opt/homebrew/bin/parley
parley 1.42.0
/opt/homebrew/bin/parley-deck-skill
/opt/homebrew/bin/parley-deck-skill
2.6.0
lrwxr-xr-x  1 tomasfecko  admin  55 Aug  7 16:25 /opt/homebrew/bin/parley-deck-skill -> ../Cellar/parley-deck-skill/2.6.0/bin/parley-deck-skill
/opt/homebrew/Cellar/parley-deck-skill/2.6.0/libexec/bin/parley-deck-skill.js
```

The version strings above are only corroboration. The decisive evidence is the sole PATH result
and the symlink/realpath into the 2.6.0 Cellar keg.

```sh
HOMEBREW_NO_AUTO_UPDATE=1 brew info --json=v2 parley-deck-cli parley-deck-skill |
  jq '.formulae[] | {name,stable:.versions.stable,url:.urls.stable.url,
      checksum:.urls.stable.checksum,installed:[.installed[].version],linked_keg}'
```

Output (exit 0):

```text
{
  "name": "parley-deck-cli",
  "stable": "1.42.0",
  "installed": ["1.42.0"],
  "linked_keg": "1.42.0"
}
{
  "name": "parley-deck-skill",
  "stable": "2.6.0",
  "installed": ["2.6.0"],
  "linked_keg": "2.6.0"
}
```

```sh
doctor=$(HOMEBREW_NO_AUTO_UPDATE=1 brew doctor 2>&1)
if printf '%s\n' "$doctor" | rg -qi 'unlinked keg'; then
  printf 'unlinked-keg warning found\n'
else
  printf 'no unlinked-keg warning\n'
fi
printf '%s\n' "$doctor" | rg '^Warning:' | head -n 5
```

Output (exit 0 for the reporting wrapper):

```text
no unlinked-keg warning
Warning: The following directories are not writable by your user:
Warning: The staging path /opt/homebrew/Caskroom is not writable by the current user.
Warning: The following taps are not trusted:
```

The writeability warnings are caused by this read-only verifier sandbox; the trust warning names
an unrelated `altfins-com/tap`. Neither formula is unlinked.

```sh
npm prefix -g
npm root -g
npm ls -g --depth=0 parley-deck-skill 2>&1 || true
find ~ /opt/homebrew/lib/node_modules /usr/local/lib/node_modules -maxdepth 7 \
  -type d -path '*/node_modules/parley-deck-skill' -print 2>/dev/null | sort -u
```

Output (exit 0):

```text
/opt/homebrew
/opt/homebrew/lib/node_modules
/opt/homebrew/lib
└── (empty)
/Users/tomasfecko/.npm/_npx/1e90c7c384e3466f/node_modules/parley-deck-skill
/Users/tomasfecko/.npm/_npx/7ba8f2ae3b103a2e/node_modules/parley-deck-skill
```

Result: **PRIMARY CONFIRMED** — no npm-global package is installed and no npm copy shadows the
Homebrew executable. Two inert `_npx` cache entries exist under `~/.npm`; neither is on PATH.

### 3. npm registry dist-tag

```sh
ego-browser nodejs <<'EOF'
const res = await fetch('https://registry.npmjs.org/parley-deck-skill',
  {headers:{Accept:'application/json','User-Agent':'codex-deploy-verifier'}})
const body = await res.json()
cliLog(JSON.stringify({status:res.status,package:body.name,
  distTags:body['dist-tags'],version260Exists:Boolean(body.versions?.['2.6.0'])},null,2))
EOF
```

Output (exit 0):

```text
{
  "status": 200,
  "package": "parley-deck-skill",
  "distTags": {"latest": "2.6.0"},
  "version260Exists": true
}
```

Result: **PRIMARY CONFIRMED** — npm `latest` is 2.6.0.

### 4. GitHub releases and Windows asset bytes

The exact requested shell commands were attempted:

```sh
gh api repos/feci/parley-deck-cli/releases/tags/v1.42.0 --jq '.assets[]|{name,digest}'
gh api repos/feci/parley-deck-skill/releases/tags/v2.6.0 --jq '.assets[]|{name,digest}'
```

Output (both exit 1; not used as evidence):

```text
error connecting to api.github.com
check your internet connection or https://githubstatus.com
```

Equivalent direct GitHub API GETs were then executed through the isolated browser. For each
returned asset I fetched `browser_download_url`, read all bytes into memory, and calculated
SHA-256 with `crypto.subtle.digest`:

```sh
ego-browser nodejs <<'EOF'
for (const [repo,tag] of [['feci/parley-deck-cli','v1.42.0'],
                          ['feci/parley-deck-skill','v2.6.0']]) {
  const res = await fetch(`https://api.github.com/repos/${repo}/releases/tags/${tag}`)
  const rel = await res.json()
  for (const a of rel.assets) {
    const download = await fetch(a.browser_download_url)
    const bytes = await download.arrayBuffer()
    const digest = await crypto.subtle.digest('SHA-256',bytes)
    const actual = [...new Uint8Array(digest)].map(b=>b.toString(16).padStart(2,'0')).join('')
    cliLog(JSON.stringify({repo,apiStatus:res.status,tag:rel.tag_name,
      assetCount:rel.assets.length,name:a.name,digest:a.digest,size:a.size,
      downloadStatus:download.status,actualBytes:bytes.byteLength,actualSha256:actual,
      digestMatches:a.digest===`sha256:${actual}`}))
  }
}
EOF
```

Output (exit 0):

```text
feci/parley-deck-cli v1.42.0: apiStatus=200, assetCount=2
  parley-v1.42.0-windows-arm64.exe
    digest=sha256:98e45df5b9c56ef448dc1e897c116cfdffeca8dbaae827433e97064e2ebafbd4
    size=6693376, downloadStatus=200, actualBytes=6693376, digestMatches=true
  parley-v1.42.0-windows-x64.exe
    digest=sha256:35921b874315ad15b634b239813aaff31850cad46ed577474fc93abdcfb83294
    size=7400960, downloadStatus=200, actualBytes=7400960, digestMatches=true

feci/parley-deck-skill v2.6.0: apiStatus=200, assetCount=2
  parley-deck-skill-v2.6.0-windows-arm64.exe
    digest=sha256:6b882abeebff571d54fb6c76d5afbbc068207b8558d5e71afe84d329734ed064
    size=85805667, downloadStatus=200, actualBytes=85805667, digestMatches=true
  parley-deck-skill-v2.6.0-windows-x64.exe
    digest=sha256:6084b12617996149f547014fabd03a3cee414c0afa721a6a9e3f3669d3a1e92f
    size=91392611, downloadStatus=200, actualBytes=91392611, digestMatches=true
```

Result: **PRIMARY CONFIRMED** — both releases exist and each has exactly two Windows assets;
GitHub's digest equals the independently calculated byte hash for all four.

### 5. winget PR state, manifest content, asset hashes, and URL status

```sh
ego-browser nodejs <<'EOF'
for (const n of [413766,413767]) {
  const prRes = await fetch(`https://api.github.com/repos/microsoft/winget-pkgs/pulls/${n}`)
  const pr = await prRes.json()
  const filesRes = await fetch(`https://api.github.com/repos/microsoft/winget-pkgs/pulls/${n}/files?per_page=100`)
  const files = await filesRes.json()
  cliLog(JSON.stringify({number:n,state:pr.state,draft:pr.draft,merged:pr.merged_at,
    headRepo:pr.head.repo.full_name,headRef:pr.head.ref,headSha:pr.head.sha,
    files:files.map(f=>f.filename)},null,2))
}
EOF
```

Output (exit 0):

```text
PR 413766: state=open, draft=false, merged=null
  head=feci/winget-pkgs:feci-cli-1.42.0
  sha=14243ff320902592f6a3fa33c4c8a24618827e6a
  manifests/f/Feci/ParleyDeckCli/1.42.0/Feci.ParleyDeckCli.installer.yaml
  manifests/f/Feci/ParleyDeckCli/1.42.0/Feci.ParleyDeckCli.locale.en-US.yaml
  manifests/f/Feci/ParleyDeckCli/1.42.0/Feci.ParleyDeckCli.yaml

PR 413767: state=open, draft=false, merged=null
  head=feci/winget-pkgs:feci-skill-2.6.0
  sha=074f8fea7f19134f94a790691775f385abc996e6
  manifests/f/Feci/ParleyDeckSkill/2.6.0/Feci.ParleyDeckSkill.installer.yaml
  manifests/f/Feci/ParleyDeckSkill/2.6.0/Feci.ParleyDeckSkill.locale.en-US.yaml
  manifests/f/Feci/ParleyDeckSkill/2.6.0/Feci.ParleyDeckSkill.yaml
```

I fetched each installer manifest from its immutable PR head SHA, extracted its two installers,
issued a redirect-following HTTP `HEAD` to every InstallerUrl, and compared its uppercase
InstallerSha256 with the independently calculated release-asset hash from the prior check:

```text
PR 413766 manifestStatus=200 PackageVersion=1.42.0
  x64 URL status=200
    InstallerSha256=35921B874315AD15B634B239813AAFF31850CAD46ED577474FC93ABDCFB83294
    actual asset SHA=35921B874315AD15B634B239813AAFF31850CAD46ED577474FC93ABDCFB83294
    uppercase=true shaMatchesActual=true
  arm64 URL status=200
    InstallerSha256=98E45DF5B9C56EF448DC1E897C116CFDFFECA8DBAAE827433E97064E2EBAFBD4
    actual asset SHA=98E45DF5B9C56EF448DC1E897C116CFDFFECA8DBAAE827433E97064E2EBAFBD4
    uppercase=true shaMatchesActual=true

PR 413767 manifestStatus=200 PackageVersion=2.6.0
  x64 URL status=200
    InstallerSha256=6084B12617996149F547014FABD03A3CEE414C0AFA721A6A9E3F3669D3A1E92F
    actual asset SHA=6084B12617996149F547014FABD03A3CEE414C0AFA721A6A9E3F3669D3A1E92F
    uppercase=true shaMatchesActual=true
  arm64 URL status=200
    InstallerSha256=6B882ABEEBFF571D54FB6C76D5AFBBC068207B8558D5E71AFE84D329734ED064
    actual asset SHA=6B882ABEEBFF571D54FB6C76D5AFBBC068207B8558D5E71AFE84D329734ED064
    uppercase=true shaMatchesActual=true
```

Result: **PRIMARY CONFIRMED** — both PRs are open; all four manifest hashes are uppercase and
match the actual release bytes; all four InstallerUrls return 200 after redirects.

### 6. All parley-deck SKILL.md copies under HOME: content, count, and markers

```sh
find ~ -maxdepth 6 -type f -name SKILL.md -path '*parley-deck*' -print 2>/dev/null | sort
```

Output (exit 0):

```text
/Users/tomasfecko/.aionrs/skills/parley-deck/SKILL.md
/Users/tomasfecko/.claude/skills/parley-deck/SKILL.md
/Users/tomasfecko/.codebuddy/skills/parley-deck/SKILL.md
/Users/tomasfecko/.codex/skills/parley-deck/SKILL.md
/Users/tomasfecko/.config/opencode/skills/parley-deck/SKILL.md
/Users/tomasfecko/.cursor/skills/parley-deck/SKILL.md
/Users/tomasfecko/.factory/skills/parley-deck/SKILL.md
/Users/tomasfecko/.gemini/config/plugins/parley-deck/SKILL.md
/Users/tomasfecko/.gemini/config/plugins/parley-deck/skills/SKILL.md
/Users/tomasfecko/.gemini/extensions/parley-deck/SKILL.md
/Users/tomasfecko/.goose/skills/parley-deck/SKILL.md
/Users/tomasfecko/.hermes/profiles/ldx/skills/parley-deck/SKILL.md
/Users/tomasfecko/.hermes/profiles/librade/skills/parley-deck/SKILL.md
/Users/tomasfecko/.hermes/profiles/testprofile/skills/parley-deck/SKILL.md
/Users/tomasfecko/.hermes/skills/parley-deck/SKILL.md
/Users/tomasfecko/.kimi-code/skills/parley-deck/SKILL.md
/Users/tomasfecko/.opencode/skills/parley-deck/SKILL.md
/Users/tomasfecko/.qwen/skills/parley-deck/SKILL.md
/Users/tomasfecko/.vibe/skills/parley-deck/SKILL.md
```

```sh
files=$(find ~ -maxdepth 6 -type f -name SKILL.md -path '*parley-deck*' -print 2>/dev/null)
total=$(printf '%s\n' "$files" | sed '/^$/d' | wc -l | tr -d ' ')
missing=$(printf '%s\n' "$files" | while IFS= read -r f; do
  [ -n "$f" ] && ! rg -Fq 'a global core, a generated deck view' "$f" && printf '%s\n' "$f"
done)
missing_count=$(printf '%s\n' "$missing" | sed '/^$/d' | wc -l | tr -d ' ')
printf 'total=%s\nmissing=%s\n' "$total" "$missing_count"
find ~ -maxdepth 6 -type f -name SKILL.md -path '*parley-deck*' -exec shasum -a 256 {} + \
  2>/dev/null | awk '{print $1}' | sort | uniq -c
```

Output (exit 0):

```text
total=19
missing=0
  19 929d0a38299c01b88c56947ce9240df8cd20b9122124a04c65c6a40e61c21c35
```

The 18 install-root markers (14 named installer targets plus the four generic destinations) all
say 2.6.0; the Antigravity target contributes the extra nested
`.gemini/config/plugins/parley-deck/skills/SKILL.md`, bringing the content-file count to 19:

```sh
find ~ -maxdepth 7 -type f -path '*/parley-deck/.parley-deck-skill-install.json' -print \
  2>/dev/null | while IFS= read -r f; do jq -r '[.target,.version]|@tsv' "$f"; done |
  awk -F '\t' '{count[$1 FS $2]++} END {for (k in count) print k "\t" count[k]}' | sort
```

Output (exit 0):

```text
agy       2.6.0  1
aionrs    2.6.0  1
claude    2.6.0  1
codebuddy 2.6.0  1
codex     2.6.0  1
cursor    2.6.0  1
droid     2.6.0  1
gemini    2.6.0  1
generic   2.6.0  4
goose     2.6.0  1
hermes    2.6.0  1
kimi      2.6.0  1
opencode  2.6.0  1
qwen      2.6.0  1
vibe      2.6.0  1
```

The unsuppressed prescribed `find` also produced normal macOS TCC `Operation not permitted`
messages for protected `~/Library` subtrees and `~/.Trash`. All 19 claimed/discovered skill paths
were accessible and checked, but this process cannot prove that no additional file exists inside
an OS-protected subtree it cannot traverse.

Result: **PRIMARY CONFIRMED for all 19 discovered copies** — every one contains the exact text
`a global core, a generated deck view`, and all 19 SKILL.md byte hashes are identical.

### 7. Installed skill origin and project compatibility status

```sh
parley-deck-skill status --target all --project . --json |
  jq '{installer:{version:.installer.version,packageRoot:.installer.packageRoot,
      executable:.installer.executable},project:{metadataDeckVersion:.project.metadata.deckVersion,
      metadataSource:.project.metadata.source,metadataMatchesProtocol:.project.metadataMatchesProtocol,
      packagedSkillVersion:.project.packaged.compatibilityManifest.skillVersion,
      protocolMatchesPackaged:.project.protocolMatchesPackaged},compatibility}'
```

Output (exit 0):

```text
{
  "installer": {
    "version": "2.6.0",
    "packageRoot": "/opt/homebrew/Cellar/parley-deck-skill/2.6.0/libexec",
    "executable": "/opt/homebrew/bin/parley-deck-skill"
  },
  "project": {
    "metadataDeckVersion": "2.5.1",
    "metadataSource": "npm:parley-deck-skill@2.5.1",
    "metadataMatchesProtocol": false,
    "packagedSkillVersion": "2.6.0",
    "protocolMatchesPackaged": false
  },
  "compatibility": {
    "status": "warning",
    "reasons": [
      "project-metadata-stale",
      "project-protocol-differs-from-packaged-reference"
    ]
  }
}
```

The installer origin is independently confirmed as Cellar 2.6.0. The project warning is retained
as a discrepancy; no sync was performed because this verification is read-only.

### 8. Protocol command exposure and scratch-deck execution

```sh
parley --help | rg -n 'parley protocol|^  protocol$|global core'
```

Output (exit 0):

```text
12:  parley protocol status|render|check [--dir DIR] [--dry-run] [--yes] [--json]
13:  parley protocol publish --version V --from FILE            (attended; requires a terminal)
```

The exact requested group-help invocation fails:

```sh
parley protocol --help
```

Output (exit 2):

```text
protocol: unknown subcommand "--help" (want status|render|check|publish)
```

Because no writes except this report were allowed, I reused a pre-existing isolated scratch deck
and core store. The installed Homebrew binary (the only `parley` on PATH) ran status against it:

```sh
PARLEY_HOME=/private/tmp/codex-probe-r4/home \
  parley protocol status --dir /private/tmp/codex-probe-r4/project --json
```

Output (exit 0):

```text
{
  "deck_protocol": "/private/tmp/codex-probe-r4/project/parley-deck/COOPERATION.md",
  "deck_protocol_read": true,
  "deck_sha256": "4beaec203e77f7fae3cbc0aeac3028117f20da5d0501416c0d30c8b7cf3a32a7",
  "installed": ["1.0.0"],
  "pinned": "1.0.0",
  "store": "/private/tmp/codex-probe-r4/home/protocol/core"
}
```

Result: top-level exposure and scratch-deck status are **PRIMARY CONFIRMED**;
`parley protocol --help` is **PRIMARY WRONG**.

## Discrepancies

1. `parley protocol --help` is not implemented as group help; it exits 2 with an unknown
   subcommand error. Counter-evidence that the command group otherwise shipped: `parley --help`
   lists it and `parley protocol status` exits 0 against a pinned scratch deck.
2. The current repository's `parley-deck/meta/version.json` remains at 2.5.1 and the installed
   2.6.0 status command reports `project-metadata-stale` plus
   `project-protocol-differs-from-packaged-reference`. No remediation was attempted.
