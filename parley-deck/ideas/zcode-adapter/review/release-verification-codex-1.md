---
agent: codex-1
idea: zcode-adapter
kind: release-verification
date: 2026-08-19
---

# Release verification: parley-deck-cli 1.45.0 / parley-deck-skill 2.9.0

[PRIMARY] Provenance convention: PRIMARY marks a command I personally executed and its observed result. SECONDARY marks content I read from a remote formula, release note, or PR manifest.

## 1. Git tags

[PRIMARY] Exact CLI tag command and output:

~~~text
$ git ls-remote --tags https://github.com/feci/parley-deck-cli.git | grep 1.45.0
2b9c10faa22c38db1a95a34bfefcf383ad1963e9	refs/tags/v1.45.0
927425cde53c982f0833e3b06f3ad16d2c461cca	refs/tags/v1.45.0^{}
~~~

[PRIMARY] Exact skill tag command and output:

~~~text
$ git ls-remote --tags https://github.com/feci/parley-deck-skill.git | grep 2.9.0
d08d55f05ce3115e438b2d8ef4e3fb22a178cba1	refs/tags/v2.9.0
b2c7f61c6b877f9d0cc23eeebbe2a5acc4230fe2	refs/tags/v2.9.0^{}
~~~

[PRIMARY] Both annotated tags and their peeled targets were present at the requested remote repositories.

[PRIMARY] Verdict: PASS.

## 2. GitHub release — CLI

[PRIMARY] Exact command and output:

~~~text
$ gh release view v1.45.0 --repo feci/parley-deck-cli --json name,tagName,assets
{"assets":[{"apiUrl":"https://api.github.com/repos/feci/parley-deck-cli/releases/assets/520532952","contentType":"application/octet-stream","createdAt":"2026-08-19T07:24:52Z","digest":"sha256:82272bfdb6cd0aab00f6bb10f146f6385b18295df474da92f2f5f1d370a859d8","downloadCount":0,"id":"RA_kwDOSZovDc4fBrPY","label":"","name":"parley-v1.45.0-darwin-arm64","size":6625154,"state":"uploaded","updatedAt":"2026-08-19T07:24:53Z","url":"https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-darwin-arm64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-cli/releases/assets/520532948","contentType":"application/octet-stream","createdAt":"2026-08-19T07:24:52Z","digest":"sha256:cff17dddf3e974161d3bcb0c6f45947481370273fd55e89f3636e6b33f4c3e05","downloadCount":0,"id":"RA_kwDOSZovDc4fBrPU","label":"","name":"parley-v1.45.0-darwin-x64","size":7204512,"state":"uploaded","updatedAt":"2026-08-19T07:24:54Z","url":"https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-darwin-x64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-cli/releases/assets/520532950","contentType":"application/octet-stream","createdAt":"2026-08-19T07:24:52Z","digest":"sha256:8e70758ed6570c883043b8192a11f605ebec6815ccbf8230473902bde1fef0fc","downloadCount":0,"id":"RA_kwDOSZovDc4fBrPW","label":"","name":"parley-v1.45.0-linux-arm64","size":6553762,"state":"uploaded","updatedAt":"2026-08-19T07:24:53Z","url":"https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-linux-arm64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-cli/releases/assets/520532947","contentType":"application/octet-stream","createdAt":"2026-08-19T07:24:52Z","digest":"sha256:55681566558d44ddef5c709ea862dd40a3a86bd620fb3e0017f88d3ada6e3375","downloadCount":0,"id":"RA_kwDOSZovDc4fBrPT","label":"","name":"parley-v1.45.0-linux-x64","size":7131298,"state":"uploaded","updatedAt":"2026-08-19T07:24:53Z","url":"https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-linux-x64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-cli/releases/assets/520532949","contentType":"application/x-msdownload","createdAt":"2026-08-19T07:24:52Z","digest":"sha256:91f9600226c4a0617277ddcf9f714c8884625750ad7a3dbafdaee728a823c900","downloadCount":0,"id":"RA_kwDOSZovDc4fBrPV","label":"","name":"parley-v1.45.0-windows-arm64.exe","size":6735360,"state":"uploaded","updatedAt":"2026-08-19T07:24:53Z","url":"https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-windows-arm64.exe"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-cli/releases/assets/520532974","contentType":"application/x-msdownload","createdAt":"2026-08-19T07:24:54Z","digest":"sha256:c50d1f51564d0334b145cb522503b545681d792ec89498450afaa47296c9e6ce","downloadCount":0,"id":"RA_kwDOSZovDc4fBrPu","label":"","name":"parley-v1.45.0-windows-x64.exe","size":7448576,"state":"uploaded","updatedAt":"2026-08-19T07:24:55Z","url":"https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-windows-x64.exe"}],"name":"v1.45.0 — the zcode adapter, and a roster with no unknowns","tagName":"v1.45.0"}
~~~

[PRIMARY] The returned tagName is v1.45.0. The returned asset array contains exactly six assets: darwin, linux, and windows, each in arm64 and x64 form.

[PRIMARY] Verdict: PASS.

## 3. GitHub release — skill

[PRIMARY] Exact command and output:

~~~text
$ gh release view v2.9.0 --repo feci/parley-deck-skill --json name,tagName,assets
{"assets":[{"apiUrl":"https://api.github.com/repos/feci/parley-deck-skill/releases/assets/520537219","contentType":"application/octet-stream","createdAt":"2026-08-19T07:29:51Z","digest":"sha256:c34266795493caf833a6fc2032d17e59347c907cf100b3d9bb853eaa57caa96f","downloadCount":0,"id":"RA_kwDOSZipIs4fBsSD","label":"","name":"parley-deck-skill-v2.9.0-linux-x64","size":71194326,"state":"uploaded","updatedAt":"2026-08-19T07:30:08Z","url":"https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-linux-x64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-skill/releases/assets/520537222","contentType":"application/octet-stream","createdAt":"2026-08-19T07:29:51Z","digest":"sha256:5773568d1aa93c802008d7b911399ad6994e66dda09afcca9e1d1afcfb5bb903","downloadCount":0,"id":"RA_kwDOSZipIs4fBsSG","label":"","name":"parley-deck-skill-v2.9.0-macos-arm64","size":65510656,"state":"uploaded","updatedAt":"2026-08-19T07:30:06Z","url":"https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-macos-arm64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-skill/releases/assets/520537223","contentType":"application/octet-stream","createdAt":"2026-08-19T07:29:51Z","digest":"sha256:0c88b4772f8f425965d57195674352818b9e893a01124385cf05679eba66a1fb","downloadCount":0,"id":"RA_kwDOSZipIs4fBsSH","label":"","name":"parley-deck-skill-v2.9.0-macos-x64","size":69269568,"state":"uploaded","updatedAt":"2026-08-19T07:30:08Z","url":"https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-macos-x64"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-skill/releases/assets/520537220","contentType":"application/x-msdownload","createdAt":"2026-08-19T07:29:51Z","digest":"sha256:2265945190dba8e323bfb92298084ab981c6606e6e44965a8f42fd5ed07e2f14","downloadCount":0,"id":"RA_kwDOSZipIs4fBsSE","label":"","name":"parley-deck-skill-v2.9.0-windows-arm64.exe","size":85806534,"state":"uploaded","updatedAt":"2026-08-19T07:30:09Z","url":"https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-windows-arm64.exe"},{"apiUrl":"https://api.github.com/repos/feci/parley-deck-skill/releases/assets/520537218","contentType":"application/x-msdownload","createdAt":"2026-08-19T07:29:51Z","digest":"sha256:06674f43bd3883e5cca526e9f8d37ff459d628ea31f06e54ca1ceed95377123b","downloadCount":0,"id":"RA_kwDOSZipIs4fBsSC","label":"","name":"parley-deck-skill-v2.9.0-windows-x64.exe","size":91393478,"state":"uploaded","updatedAt":"2026-08-19T07:30:09Z","url":"https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-windows-x64.exe"}],"name":"v2.9.0 — zcode install target; config-read roster status terms","tagName":"v2.9.0"}
~~~

[PRIMARY] The returned tagName is v2.9.0 and the returned asset array contains exactly five assets.

[PRIMARY] Verdict: PASS.

## 4. Released binary integrity and version

[PRIMARY] Exact download-and-run command and output:

~~~text
$ curl -sL -o /tmp/p145 https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-darwin-arm64 && chmod +x /tmp/p145 && /tmp/p145 --version
parley 1.45.0
~~~

[PRIMARY] Exact local SHA-256 command and output:

~~~text
$ shasum -a 256 /tmp/p145
82272bfdb6cd0aab00f6bb10f146f6385b18295df474da92f2f5f1d370a859d8  /tmp/p145
~~~

[PRIMARY] The version output is exactly parley 1.45.0. The recomputed SHA-256 equals the digest returned for the darwin-arm64 asset in channel 2.

[PRIMARY] Verdict: PASS.

## 5. npm

[PRIMARY] Exact version command and output:

~~~text
$ npm view parley-deck-skill version
npm warn Unknown user config "store-dir". This will stop working in the next major version of npm. See `npm help npmrc` for supported config options.
2.9.0
~~~

[PRIMARY] Exact dist-tags command and output:

~~~text
$ npm view parley-deck-skill dist-tags
npm warn Unknown user config "store-dir". This will stop working in the next major version of npm. See `npm help npmrc` for supported config options.
{ latest: '2.9.0' }
~~~

[PRIMARY] npm reports package version 2.9.0 and latest points to 2.9.0. The store-dir warning does not alter either returned value.

[PRIMARY] Verdict: PASS.

## 6. Homebrew tap

[PRIMARY] Exact CLI formula read command and relevant output:

~~~text
$ curl -sL https://raw.githubusercontent.com/feci/homebrew-parley/main/Formula/parley-deck-cli.rb | rg '^  (url|sha256) '
  url "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.45.0.tar.gz"
  sha256 "37611030b15bbed3a897f24950a320106479c91bf07974515563532132b8d2fe"
~~~

[SECONDARY] The remote CLI formula pins the v1.45.0 source URL and declares SHA-256 37611030b15bbed3a897f24950a320106479c91bf07974515563532132b8d2fe.

[PRIMARY] Exact CLI source-hash recomputation and output:

~~~text
$ curl -sL https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.45.0.tar.gz | shasum -a 256
37611030b15bbed3a897f24950a320106479c91bf07974515563532132b8d2fe  -
~~~

[PRIMARY] The recomputed CLI hash exactly matches the formula.

[PRIMARY] Exact skill formula read command and relevant output:

~~~text
$ curl -sL https://raw.githubusercontent.com/feci/homebrew-parley/main/Formula/parley-deck-skill.rb | rg '^  (url|sha256) '
  url "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.9.0.tar.gz"
  sha256 "b1872de6bcb5906c49da9b67c898e3ebb07a913c480c9b74bdfcd934d489e3f7"
~~~

[SECONDARY] The remote skill formula pins the v2.9.0 source URL and declares SHA-256 b1872de6bcb5906c49da9b67c898e3ebb07a913c480c9b74bdfcd934d489e3f7.

[PRIMARY] Exact skill source-hash recomputation and output:

~~~text
$ curl -sL https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.9.0.tar.gz | shasum -a 256
b1872de6bcb5906c49da9b67c898e3ebb07a913c480c9b74bdfcd934d489e3f7  -
~~~

[PRIMARY] The recomputed skill hash exactly matches the formula. No Homebrew SHA mismatch was observed.

[PRIMARY] Verdict: PASS.

## 7. winget

[PRIMARY] Exact CLI PR command and output:

~~~text
$ gh pr view 420440 --repo microsoft/winget-pkgs --json state,title,files
{"files":[{"path":"manifests/f/Feci/ParleyDeckCli/1.45.0/Feci.ParleyDeckCli.installer.yaml","additions":15,"deletions":0,"changeType":"ADDED"},{"path":"manifests/f/Feci/ParleyDeckCli/1.45.0/Feci.ParleyDeckCli.locale.en-US.yaml","additions":28,"deletions":0,"changeType":"ADDED"},{"path":"manifests/f/Feci/ParleyDeckCli/1.45.0/Feci.ParleyDeckCli.yaml","additions":6,"deletions":0,"changeType":"ADDED"}],"state":"OPEN","title":"New version: Feci.ParleyDeckCli version 1.45.0"}
~~~

[PRIMARY] Exact skill PR command and output:

~~~text
$ gh pr view 420442 --repo microsoft/winget-pkgs --json state,title,files
{"files":[{"path":"manifests/f/Feci/ParleyDeckSkill/2.9.0/Feci.ParleyDeckSkill.installer.yaml","additions":15,"deletions":0,"changeType":"ADDED"},{"path":"manifests/f/Feci/ParleyDeckSkill/2.9.0/Feci.ParleyDeckSkill.locale.en-US.yaml","additions":27,"deletions":0,"changeType":"ADDED"},{"path":"manifests/f/Feci/ParleyDeckSkill/2.9.0/Feci.ParleyDeckSkill.yaml","additions":6,"deletions":0,"changeType":"ADDED"}],"state":"OPEN","title":"New version: Feci.ParleyDeckSkill version 2.9.0"}
~~~

[PRIMARY] Exact application-directory reductions and outputs:

~~~text
$ gh pr view 420440 --repo microsoft/winget-pkgs --json files --jq '[.files[].path | split("/")[0:4] | join("/")] | unique'
["manifests/f/Feci/ParleyDeckCli"]

$ gh pr view 420442 --repo microsoft/winget-pkgs --json files --jq '[.files[].path | split("/")[0:4] | join("/")] | unique'
["manifests/f/Feci/ParleyDeckSkill"]
~~~

[PRIMARY] Each PR touches exactly one application directory.

[PRIMARY] Exact commands used to read the installer values and their outputs:

~~~text
$ gh pr diff 420440 --repo microsoft/winget-pkgs | grep -E '^\+(PackageIdentifier|PackageVersion|- Architecture:|  InstallerUrl:|  InstallerSha256:)'
+PackageIdentifier: Feci.ParleyDeckCli
+PackageVersion: 1.45.0
+- Architecture: x64
+  InstallerUrl: https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-windows-x64.exe
+  InstallerSha256: C50D1F51564D0334B145CB522503B545681D792EC89498450AFAA47296C9E6CE
+- Architecture: arm64
+  InstallerUrl: https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-windows-arm64.exe
+  InstallerSha256: 91F9600226C4A0617277DDCF9F714C8884625750AD7A3DBAFDAEE728A823C900
+PackageIdentifier: Feci.ParleyDeckCli
+PackageVersion: 1.45.0
+PackageIdentifier: Feci.ParleyDeckCli
+PackageVersion: 1.45.0

$ gh pr diff 420442 --repo microsoft/winget-pkgs | grep -E '^\+(PackageIdentifier|PackageVersion|- Architecture:|  InstallerUrl:|  InstallerSha256:)'
+PackageIdentifier: Feci.ParleyDeckSkill
+PackageVersion: 2.9.0
+- Architecture: x64
+  InstallerUrl: https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-windows-x64.exe
+  InstallerSha256: 06674F43BD3883E5CCA526E9F8D37FF459D628EA31F06E54CA1CEED95377123B
+- Architecture: arm64
+  InstallerUrl: https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-windows-arm64.exe
+  InstallerSha256: 2265945190DBA8E323BFB92298084AB981C6606E6E44965A8F42FD5ED07E2F14
+PackageIdentifier: Feci.ParleyDeckSkill
+PackageVersion: 2.9.0
+PackageIdentifier: Feci.ParleyDeckSkill
+PackageVersion: 2.9.0
~~~

[SECONDARY] PR 420440 declares CLI x64 SHA C50D1F51564D0334B145CB522503B545681D792EC89498450AFAA47296C9E6CE and arm64 SHA 91F9600226C4A0617277DDCF9F714C8884625750AD7A3DBAFDAEE728A823C900. PR 420442 declares skill x64 SHA 06674F43BD3883E5CCA526E9F8D37FF459D628EA31F06E54CA1CEED95377123B and arm64 SHA 2265945190DBA8E323BFB92298084AB981C6606E6E44965A8F42FD5ED07E2F14.

[PRIMARY] Exact released-executable SHA-256 recomputations and outputs:

~~~text
$ curl -sL https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-windows-x64.exe | shasum -a 256
c50d1f51564d0334b145cb522503b545681d792ec89498450afaa47296c9e6ce  -

$ curl -sL https://github.com/feci/parley-deck-cli/releases/download/v1.45.0/parley-v1.45.0-windows-arm64.exe | shasum -a 256
91f9600226c4a0617277ddcf9f714c8884625750ad7a3dbafdaee728a823c900  -

$ curl -sL https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-windows-x64.exe | shasum -a 256
06674f43bd3883e5cca526e9f8d37ff459d628ea31f06e54ca1ceed95377123b  -

$ curl -sL https://github.com/feci/parley-deck-skill/releases/download/v2.9.0/parley-deck-skill-v2.9.0-windows-arm64.exe | shasum -a 256
2265945190dba8e323bfb92298084ab981c6606e6e44965a8f42fd5ed07e2f14  -
~~~

[PRIMARY] All four recomputed hashes match their PR InstallerSha256 values case-insensitively.

[PRIMARY] Both PRs are OPEN. This verifies the submitted manifests and released payloads, but does not claim the packages are already merged into or installable from the winget catalog.

[PRIMARY] Verdict: PASS for the requested PR structure and integrity checks; catalog publication is pending.

## 8. Release self-claim: machine roster has no unknown cells

[PRIMARY] Exact release-note query and output:

~~~text
$ gh release view v1.45.0 --repo feci/parley-deck-cli --json body --jq .body | rg 'reports zero .*unknown.* cells'
Result: `parley roster show --scope machine` reports zero `unknown` cells.
~~~

[SECONDARY] The release notes claim that the machine-scope roster contains zero unknown cells.

[PRIMARY] Exact command run with the downloaded released binary and its output:

~~~text
$ /tmp/p145 roster show --scope machine
AGENT        ADAPTER    STATE    INSTALLED MODEL                  MODEL-FAMILY   MODEL-COMPANY EFFORT   SPEED    AUTO STATUS
claude-1     claude     active   yes       claude-opus-5[1m]      Claude Opus    Anthropic     max      deep     yes  ok
codex-1      codex      active   yes       gpt-5.6-sol            GPT            OpenAI        max      deep     yes  ok
hermes-1     hermes     active   yes       fireworks/inkling      Inkling        Thinking Machines Lab high     deep     yes  ok
kimi-1       kimi       active   yes       kimi-code/k3           Kimi K         Moonshot AI   max      deep     yes  effort-from-config
opencode-1   opencode   active   yes       litellm/xai/grok-4.6   Grok           xAI           xhigh    deep     yes  effort-from-config
zcode-1      zcode      active   yes       zai/glm-5.3            GLM            Zhipu AI      max      deep     yes  model-from-config,effort-from-config
~~~

[PRIMARY] Exact unknown-token count over the same released-binary command and output:

~~~text
$ /tmp/p145 roster show --scope machine | grep -o unknown | wc -l
       0
~~~

[PRIMARY] The actual zcode-1 row reports model zai/glm-5.3, effort max, and STATUS model-from-config,effort-from-config. The complete table contains zero unknown tokens.

[PRIMARY] Verdict: PASS.

## Overall verdict

[PRIMARY] Overall verdict: PASS. Every explicitly requested tag, release, asset-version, npm, Homebrew checksum, winget PR structure/hash, and released-binary roster check passed.

[PRIMARY] Qualification: winget PRs 420440 and 420442 remain OPEN, so winget catalog publication itself is pending and was not claimed as verified.
