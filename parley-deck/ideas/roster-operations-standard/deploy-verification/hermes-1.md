---
idea: roster-operations-standard
phase: deploy verification
agent: hermes-1
date: 2026-08-06
verdict: DISCREPANCIES
---
# Deploy verification — hermes-1

## Verdict

DISCREPANCIES. The release artifacts (Homebrew formulae, GitHub releases, winget manifests,
npm dist-tag) are all correct and internally consistent. However, the claim of "7 runtime
skill installs at 2.5.1" is incomplete and misleading: while the 7 named targets (codex,
claude, agy, gemini, hermes, kimi, opencode) ARE all at 2.5.1 with correct content, there
are 7 ADDITIONAL stale skill installs on this machine (qwen, cursor, goose, aionrs, vibe,
factory, codebuddy) that remain at a pre-2.5.1 version (804 lines, mtime Jun 13 2026) and
do NOT contain the ratified content. The summary's "7 runtime skill installs" is not wrong
about the 7 it names, but it silently omits 7 more that are stale. The fleet count (35/36)
is correct.

## Channel-by-channel (command + quoted output for each)

### 1. Homebrew formulae — URL and sha256 verification

[PRIMARY] Both formulae exist and are bumped to the claimed versions.

Command: `ls -la "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/homebrew-parley/Formula/"`
Output:
```
-rw-r--r--@  1 tomasfecko  staff   705 Aug  6 22:38 parley-deck-cli.rb
-rw-r--r--@  1 tomasfecko  staff  2251 Aug  6 23:19 parley-deck-skill.rb
```

[PRIMARY] parley-deck-cli.rb formula (read_file):
```
url "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.41.0.tar.gz"
sha256 "381b4b386c899f5c9cd4a7c907c32409c110b7ed1d808aef89d3a2017f524cab"
```

[PRIMARY] parley-deck-skill.rb formula (read_file):
```
url "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.5.1.tar.gz"
sha256 "484f03129751ad0c297efb4b8096b3ccb03cb25278355e767dad7ed42f7202bf"
```

[PRIMARY] Tags exist (gh api):
Command: `gh api repos/feci/parley-deck-cli/git/ref/tags/v1.41.0 --jq '.object.type + " " + .object.sha'`
Output: `tag c7c9d3b8c887e7b9b421cb7bdaa98148f80954b6`
Command: `gh api repos/feci/parley-deck-skill/git/ref/tags/v2.5.1 --jq '.object.type + " " + .object.sha'`
Output: `tag 04ab4ec8b34783e98d8dd17a20eb7094e5b0e7d7`

[PRIMARY] Tarball download and sha256 verification:
Command: `curl -sL "https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.41.0.tar.gz" -o /tmp/cli.tar.gz && shasum -a 256 /tmp/cli.tar.gz`
Output: `381b4b386c899f5c9cd4a7c907c32409c110b7ed1d808aef89d3a2017f524cab` (5018125 bytes)
Formula says: `381b4b386c899f5c9cd4a7c907c32409c110b7ed1d808aef89d3a2017f524cab`
MATCH.

Command: `curl -sL "https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.5.1.tar.gz" -o /tmp/skill.tar.gz && shasum -a 256 /tmp/skill.tar.gz`
Output: `484f03129751ad0c297efb4b8096b3ccb03cb25278355e767dad7ed42f7202bf` (472706 bytes)
Formula says: `484f03129751ad0c297efb4b8096b3ccb03cb25278355e767dad7ed42f7202bf`
MATCH.

[PRIMARY] Homebrew tap git log confirms both bumps committed:
Command: `cd homebrew-parley && git log --oneline -5`
Output:
```
573cb27 parley-deck-skill 2.5.1
04fc606 parley-deck-cli 1.41.0
149187c parley-deck-cli 1.40.1 — review fixes
652d573 parley-deck-cli 1.40.0 + parley-deck-skill 2.5.0 — standardized roster operations
```
Working tree clean (`git status --short` empty).

VERDICT: Both formulae correct. URLs point to existing tags. sha256 hashes match downloaded
tarballs exactly. PASS.

### 2. Binary resolution — Cellar vs npm/working tree

[PRIMARY] Command: `which -a parley-deck-skill`
Output: `/opt/homebrew/bin/parley-deck-skill` (single entry — no shadowing)

[PRIMARY] Command: `ls -la $(which parley-deck-skill)`
Output: `lrwxr-xr-x 1 tomasfecko admin 55 Aug 6 23:19 /opt/homebrew/bin/parley-deck-skill -> ../Cellar/parley-deck-skill/2.5.1/bin/parley-deck-skill`

[PRIMARY] Command: `realpath $(which parley-deck-skill)`
Output: `/opt/homebrew/Cellar/parley-deck-skill/2.5.1/libexec/bin/parley-deck-skill.js`

Resolves into Cellar/ — NOT into node_modules or a working tree. PASS.

[PRIMARY] npm global check — no shadowing:
Command: `ls -la /opt/homebrew/lib/node_modules/ | grep -i parley`
Output: (empty, exit 1 — no npm-installed parley-deck-skill)
Command: `find /opt/homebrew/lib/node_modules -maxdepth 2 -name "parley-deck-skill"`
Output: (empty)

The symlink in /opt/homebrew/bin/ is the brew Cellar symlink; npm's global bin dir is the
same /opt/homebrew/bin, but there is no separate npm package. No shadowing. PASS.

[PRIMARY] brew doctor — no unlinked kegs:
Command: `brew doctor 2>&1 | grep -i "unlink\|parley\|keg"`
Output: (empty, exit 1 — no mentions of unlinked kegs or parley)
Full brew doctor output is 49 lines, all unrelated warnings (tap trust for altfins-com/tap,
PATH ordering). No unlinked keg warnings. PASS.

[PRIMARY] Installed versions:
Command: `brew list --versions parley-deck-cli parley-deck-skill`
Output: `parley-deck-cli 1.41.0` / `parley-deck-skill 2.5.1`
Command: `parley version` → `parley 1.41.0`
Command: `parley-deck-skill --version` → `2.5.1`

[PRIMARY] Cellar kegs:
Command: `ls -la /opt/homebrew/Cellar/parley-deck-skill/`
Output: `drwxr-xr-x 11 tomasfecko admin 352 Aug 6 23:19 2.5.1` (single version)
Command: `ls -la /opt/homebrew/Cellar/parley-deck-cli/`
Output: `drwxr-xr-x 9 tomasfecko admin 288 Aug 6 22:38 1.41.0` (single version)

VERDICT: Binary resolves into Cellar, no npm shadowing, no unlinked kegs. PASS.

### 3. npm dist-tags and GitHub release assets

[PRIMARY] Command: `npm dist-tag ls parley-deck-skill`
Output: `latest: 2.5.1`

[PRIMARY] GitHub release assets — CLI:
Command: `gh api repos/feci/parley-deck-cli/releases/tags/v1.41.0 --jq '.assets[]|{name,digest}'`
Output:
```
{"digest":"sha256:b4d63ee1871c5c443058eb4254836df907ae89149f2d7355a67fe9f4012ec314","name":"parley-v1.41.0-windows-arm64.exe"}
{"digest":"sha256:01472d76fe8c80ee4df9be1e06f4adfea83def065db41a25a546b8915ea03d9a","name":"parley-v1.41.0-windows-x64.exe"}
```
2 assets (arm64 + x64). PASS.

[PRIMARY] GitHub release assets — skill:
Command: `gh api repos/feci/parley-deck-skill/releases/tags/v2.5.1 --jq '.assets[]|{name,digest}'`
Output:
```
{"digest":"sha256:aa4612cfbc9be0267f7dcc0b251337ed4a64bcc0ee7c1a58826d65ac3e67109b","name":"parley-deck-skill-v2.5.1-windows-arm64.exe"}
{"digest":"sha256:18e3145d15cdf2006cd1c0eaaf42c5645962d73299db667f558899567f4e4708","name":"parley-deck-skill-v2.5.1-windows-x64.exe"}
```
2 assets (arm64 + x64). PASS.

VERDICT: npm latest=2.5.1, both GitHub releases exist with 2 Windows assets each. PASS.

### 4. CONTENT check — installed SKILL.md content

[PRIMARY] The installed skill must CONTAIN "Membership is the DECK FILE" and must NOT contain
"roster sync`moves it across" (the stale phrase).

Command (claude): `grep -n "Membership is the DECK FILE" ~/.claude/skills/parley-deck/SKILL.md`
Output: `307:**Membership is the DECK FILE.** The machine layer (~/.parley/agents.toml) seeds *values* for`
Command (claude): `grep -n 'moves it across' ~/.claude/skills/parley-deck/SKILL.md`
Output: (empty — not found)
PASS.

Command (codex): `grep -n "Membership is the DECK FILE" ~/.codex/skills/parley-deck/SKILL.md`
Output: `307:**Membership is the DECK FILE.** The machine layer (~/.parley/agents.toml) seeds *values* for`
Command (codex): `grep -n 'moves it across' ~/.codex/skills/parley-deck/SKILL.md`
Output: (empty — not found)
PASS.

[PRIMARY] All 7 claimed targets verified via `parley-deck-skill doctor` (authoritative path
resolution) + content grep:

| Target | Real path | Has "Membership is the DECK FILE" | Has "moves it across" | Lines | mtime |
|--------|-----------|----------------------------------|-----------------------|-------|-------|
| codex | ~/.codex/skills/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |
| claude | ~/.claude/skills/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |
| agy | ~/.gemini/config/plugins/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |
| gemini | ~/.gemini/extensions/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |
| hermes | ~/.hermes/skills/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |
| kimi | ~/.kimi-code/skills/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |
| opencode | ~/.opencode/skills/parley-deck/SKILL.md | YES (line 307) | NO | 925 | Aug 6 23:20 |

NOTE: `agy` and `gemini` targets install to non-standard paths (~/.gemini/config/plugins/
and ~/.gemini/extensions/ respectively), NOT ~/.agy/skills/ or ~/.gemini/skills/. The
`parley-deck-skill doctor` command reveals the real paths. All 7 claimed targets are
correct content at 925 lines, mtime Aug 6 23:20.

[PRIMARY] `parley-deck-skill doctor` output confirms all 7 targets valid:
```
codex/parley-deck: valid /Users/tomasfecko/.codex/skills/parley-deck
claude/parley-deck: valid /Users/tomasfecko/.claude/skills/parley-deck
agy/parley-deck: valid /Users/tomasfecko/.gemini/config/plugins/parley-deck
gemini/parley-deck: valid /Users/tomasfecko/.gemini/extensions/parley-deck
hermes/parley-deck: valid /Users/tomasfecko/.hermes/skills/parley-deck
kimi/parley-deck: valid /Users/tomasfecko/.kimi-code/skills/parley-deck
opencode/parley-deck: valid /Users/tomasfecko/.opencode/skills/parley-deck
```

VERDICT: All 7 claimed targets have correct content. PASS for the claimed 7.

### 5. winget manifests — InstallerUrl and hash verification

[PRIMARY] winget PRs are open:
Command: `gh pr view 413351 --repo microsoft/winget-pkgs --json number,title,state,headRefName,author`
Output: `{"number":413351,"state":"OPEN","title":"New version: Feci.ParleyDeckCli version 1.41.0","headRefName":"feci-parley-1.41.0","author":{"login":"feci"}}`
Command: `gh pr view 413352 --repo microsoft/winget-pkgs --json number,title,state,headRefName,author`
Output: `{"number":413352,"state":"OPEN","title":"New version: Feci.ParleyDeckSkill version 2.5.1","headRefName":"feci-skill-2.5.1","author":{"login":"feci"}}`

[PRIMARY] CLI installer manifest (from fork branch feci-parley-1.41.0):
```
Installers:
- Architecture: x64
  InstallerUrl: https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-x64.exe
  InstallerSha256: 01472D76FE8C80EE4DF9BE1E06F4ADFEA83DEF065DB41A25A546B8915EA03D9A
- Architecture: arm64
  InstallerUrl: https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-arm64.exe
  InstallerSha256: B4D63EE1871C5C443058EB4254836DF907AE89149F2D7355A67FE9F4012EC314
```

[PRIMARY] SKILL installer manifest (from fork branch feci-skill-2.5.1):
```
Installers:
- Architecture: x64
  InstallerUrl: https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-x64.exe
  InstallerSha256: 18E3145D15CDF2006CD1C0EAAF42C5645962D73299DB667F558899567F4E4708
- Architecture: arm64
  InstallerUrl: https://github.com/feci/parley-deck-skill/releases/download/v2.5.1/parley-deck-skill-v2.5.1-windows-arm64.exe
  InstallerSha256: AA4612CFBC9BE0267F7DCC0B251337ED4A64BCC0EE7C1A58826D65AC3E67109B
```

[PRIMARY] InstallerUrls return HTTP 200:
Command: `curl -sIL https://github.com/feci/parley-deck-cli/releases/download/v1.41.0/parley-v1.41.0-windows-x64.exe | grep "^HTTP"`
Output: `HTTP/2 302` then `HTTP/2 200` (redirects to release-assets.githubusercontent.com, final 200)
All 4 URLs return 302→200. PASS.

[PRIMARY] winget InstallerSha256 vs GitHub release digests (case-insensitive):

| Asset | GitHub digest (lowercase) | winget InstallerSha256 (uppercase) | Match |
|-------|--------------------------|-------------------------------------|-------|
| CLI x64 | 01472d76...03d9a | 01472D76...03D9A | MATCH |
| CLI arm64 | b4d63ee1...ec314 | B4D63EE1...EC314 | MATCH |
| SKILL x64 | 18e3145d...e4708 | 18E3145D...E4708 | MATCH |
| SKILL arm64 | aa4612cf...e109b | AA4612CF...E109B | MATCH |

VERDICT: Both winget PRs open, InstallerUrls return 200, all 4 hashes match GitHub release
digests exactly. PASS.

### 6. Fleet claim — deck protocol sync count

[PRIMARY] Command: `find "/Volumes/My Shared Files/AI_WORKSPACE" -maxdepth 4 -type f -path '*/parley-deck/COOPERATION.md' | wc -l`
Output: `36` total decks found.

[PRIMARY] Command (counting decks WITH ratified marker "NOT authoritative"):
`find . -maxdepth 4 -type f -path '*/parley-deck/COOPERATION.md' -exec /usr/bin/grep -l "NOT authoritative" {} \; | wc -l`
Output: `35`

[PRIMARY] The 1 deck WITHOUT the marker:
Command: `/usr/bin/grep -c "NOT authoritative" "/Volumes/My Shared Files/AI_WORKSPACE/ecb-meeting-2026.05/parley-deck/COOPERATION.md"`
Output: `0`

This deck (ecb-meeting-2026.05) is an older protocol version (missing the "NOT authoritative"
marker phrase). It does not contain "read-only" or "skipped" markers either — it is simply
on an older protocol revision. This is consistent with the claim that it was "skipped:
read-only by design."

VERDICT: 35 of 36 decks synced. 1 unsynced: ecb-meeting-2026.05. Claim is accurate. PASS.

## Discrepancies

### DISCREPANCY 1: 7 additional stale skill installs not reported

[PRIMARY] The claim states "7 runtime skill installs (codex, claude, agy, gemini, hermes,
kimi, opencode) at 2.5.1." This is true for those 7 — all verified with correct content
(925 lines, "Membership is the DECK FILE" present, "moves it across" absent, mtime Aug 6).

However, `find ~ -maxdepth 4 -type f -path '*/skills/parley-deck/SKILL.md'` reveals 12
total skill installs on this machine. The 7 NOT mentioned in the claim are all STALE:

| Target | Path | Lines | mtime | Has "Membership is the DECK FILE" |
|--------|------|-------|-------|----------------------------------|
| qwen | ~/.qwen/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |
| cursor | ~/.cursor/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |
| goose | ~/.goose/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |
| aionrs | ~/.aionrs/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |
| vibe | ~/.vibe/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |
| factory | ~/.factory/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |
| codebuddy | ~/.codebuddy/skills/parley-deck/SKILL.md | 804 | Jun 13 13:02 | NO |

All 7 are at 804 lines (vs 925 for 2.5.1), mtime Jun 13 2026 (vs Aug 6 for the update), and
do NOT contain "Membership is the DECK FILE." They are on a pre-2.5.1 version.

The claim is not false about the 7 it names, but it is incomplete — a deploy report that
says "7 runtime skill installs at 2.5.1" without mentioning 7 more that are stale gives a
false impression of fleet-wide consistency. A reader would reasonably conclude the entire
fleet is at 2.5.1, when in fact only 7 of 14 installs are.

These 7 stale installs are NOT among the claimed 7 targets (codex, claude, agy, gemini,
hermes, kimi, opencode), so the claim is not internally contradictory — but the omission
is material.

### No other discrepancies

All other channels verified clean:
- Homebrew formulae: both URL→tag and sha256→tarball MATCH (downloaded and hashed)
- Binary resolution: Cellar symlink, no npm shadowing, no unlinked kegs
- npm dist-tag: latest=2.5.1
- GitHub releases: both exist with 2 Windows assets each
- winget: both PRs OPEN, InstallerUrls return 200, all 4 hashes match GitHub digests
- Fleet: 35/36 decks synced (ecb-meeting-2026.05 unsynced, as claimed)
- Content: all 7 claimed targets have correct content (marker present, stale phrase absent)
