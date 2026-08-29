---
from: codex-1
to: all
idea: protocol-generation-bias
phase: release
blocking: no
date: 2026-08-29
---

## Verified

### 7. Locally installed binaries — CONFIRMED (PRIMARY)

Command:

```text
$ command -v parley
/opt/homebrew/bin/parley
$ parley --version
parley 1.47.0
$ command -v parley-deck-skill
/opt/homebrew/bin/parley-deck-skill
$ parley-deck-skill --version
2.11.0
```

Both installed executables report the requested release versions.

### 8. Installed-binary runtime rejection — CONFIRMED (PRIMARY)

I exercised `/opt/homebrew/bin/parley`, not a source helper or `go test`. I initialized an
isolated temporary workspace and `PARLEY_HOME`, configured a local probe command as the
`codex` adapter, and had that command write this otherwise complete round-1 artifact:

```markdown
---
agent: codex
idea: 2026-08-29T21-03-56-runtime-gate-pro
round: 1
date: 2026-08-29
---

## Summary
A complete runtime-probe artifact.

## Proposed approach
Exercise the installed CLI validation path.

## Concerns / open questions
None.

## Risks
None.
```

The heading scan showed exactly the four legacy non-empty sections and no
`## Existing alternatives`:

```text
$ rg -n '^## ' .../round-01/codex.md
8:## Summary
11:## Proposed approach
14:## Concerns / open questions
17:## Risks
```

Installed-binary command:

```text
$ PARLEY_HOME=/private/tmp/parley-release-verification-codex-1.G5odhv/parley-home \
    /opt/homebrew/bin/parley run \
    --dir /private/tmp/parley-release-verification-codex-1.G5odhv/runtime-gate \
    --no-tui --no-auto --no-ping --no-preflight --yes \
    --participants codex 'Runtime gate probe'
```

Relevant output and process result:

```text
Created idea 2026-08-29T21-03-56-runtime-gate-pro and run 20260829T190356.698547000Z
Starting round-01 with participants: codex
  codex    failed: .../round-01/codex.md is missing a non-empty "## Existing alternatives" section (§15.6a)
exit=1
```

A separate installed-binary status read returned:

```json
{
  "terminal": true,
  "outcome": "incomplete",
  "attention": "FAILED",
  "state": {
    "agents": [{
      "id": "codex",
      "state": "failed",
      "error": ".../round-01/codex.md is missing a non-empty \"## Existing alternatives\" section (§15.6a)"
    }],
    "round_status": "incomplete"
  }
}
```

Verdict: the functional release claim is true. The installed Parley 1.47.0 runtime rejects a
round-1 artifact that is otherwise complete but lacks a non-empty
`## Existing alternatives` section. The temporary workspace and probe files were removed after
the check.

## Failed or missing

### Runtime skill copies were not propagated to 2.11.0

This is distinct from channel 7's executable-version result. The installer executable is 2.11.0,
but:

```text
$ parley-deck-skill status --target all --project . --json
installer.version: 2.11.0
parleyCli.version: "parley 1.47.0"
runtimeInstalls[*].version: 2.10.0
runtimeInstalls[*].versionMatchesInstaller: false
compatibility.status: "warning"
```

The command reported 2.10.0 markers for every detected runtime target it listed (Codex, Claude,
Antigravity, Gemini, Hermes, Kimi, OpenCode, and ZCode). Thus the new installer is present, but
the managed runtime skill copies are still the previous version.

## Not checked and why

The shell has no working external network path. I did not promote local tags, release staging
files, cached web pages, release notes, or commit messages into live-channel verdicts.

The common failures were:

```text
$ gh release view v1.47.0 --repo feci/parley-deck-cli
error connecting to api.github.com
check your internet connection or https://githubstatus.com

$ git ls-remote --tags origin ...
fatal: unable to access 'https://github.com/...': Could not resolve host: github.com

$ npm view parley-deck-skill version --fetch-retries=0 --fetch-timeout=5000
npm error code ENOTFOUND
npm error network request to https://registry.npmjs.org/parley-deck-skill failed,
reason: getaddrinfo ENOTFOUND registry.npmjs.org
```

The required ego-browser fallback could not connect from this sandbox because it requires Full
Access; the isolated in-app browser was unavailable. Direct HTTP and the configured MCP Anywhere
gateway also had no working route. Chrome was not used, per project policy.

### 1. GitHub release `feci/parley-deck-cli@v1.47.0` — NOT CHECKED

Both the exact requested `gh release view` command and its JSON form failed before reaching
GitHub. Therefore release existence, six-asset count, asset identities, and notes presence were
not checked.

### 2. GitHub release `feci/parley-deck-skill@v2.11.0` — NOT CHECKED

The corresponding exact and JSON `gh release view` commands failed before reaching GitHub.
Therefore release existence and five-asset count were not checked.

### 3. Origin tags — NOT CHECKED

`git ls-remote --tags origin` failed in both repositories, so origin reachability was not
checked. Local-only evidence was:

```text
parley-deck-cli:
46ba8919ff14c1eaf06590c21eaef973c0711844 refs/tags/v1.47.0
a11fe623ddb5dadb3269c453fb41fbef2237b025 refs/tags/protocol-generation-bias-baseline

parley-deck-skill:
fbca86a9db50e79744a27bf9f4ac6bb0679673d3 refs/tags/v2.11.0
```

The skill repository had no local `protocol-generation-bias-baseline` ref. These local refs do
not prove any tag is pushed or reachable on origin.

### 4. Homebrew tap — NOT CHECKED AS A LIVE CHANNEL

The locally installed tap checkout does point to the requested tags:

```text
parley-deck-cli:
url    https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.47.0.tar.gz
sha256 3f62c443e3b5f3f4a3899213347967ae13d95336ff17ed373008381330b29664

parley-deck-skill:
url    https://github.com/feci/parley-deck-skill/archive/refs/tags/v2.11.0.tar.gz
sha256 b5fad20956ecf6785819ed79f6a3e26c47eb79d0db9da84cd0d429adc75449dd
```

I independently ran `shasum -a 256` over the two Homebrew cached archive files:

```text
3f62c443e3b5f3f4a3899213347967ae13d95336ff17ed373008381330b29664  ...--parley-deck-cli-1.47.0.tar.gz
b5fad20956ecf6785819ed79f6a3e26c47eb79d0db9da84cd0d429adc75449dd  ...--parley-deck-skill-2.11.0.tar.gz
```

Those local bytes match the local formulas. However, I could not freshly download either URL or
read the remote tap tip. Because the task explicitly requires a fresh independent download, this
is partial local evidence only and the Homebrew channel is not marked verified.

### 5. winget PRs #426220 and #426221 — NOT CHECKED

`gh pr view` could not reach GitHub, so I could not confirm either PR's live state, target branch,
or actual live file set. I also could not resolve and download the installer URLs, so no
independent URL-byte SHA-256 comparison was possible.

The local winget branch snapshots each differ from the local `upstream/master` snapshot by three
files under one package/version directory:

```text
Feci.ParleyDeckCli/1.47.0:
  Feci.ParleyDeckCli.installer.yaml
  Feci.ParleyDeckCli.locale.en-US.yaml
  Feci.ParleyDeckCli.yaml

Feci.ParleyDeckSkill/2.11.0:
  Feci.ParleyDeckSkill.installer.yaml
  Feci.ParleyDeckSkill.locale.en-US.yaml
  Feci.ParleyDeckSkill.yaml
```

That is not a substitute for inspecting the live PRs. The declared installer hashes therefore
remain unverified against downloaded release assets.

### 6. npm — NOT CHECKED

The exact `npm view parley-deck-skill version` request could not reach the registry; a bounded
retry returned `ENOTFOUND registry.npmjs.org`, and offline mode returned `ENOTCACHED`.
Therefore I cannot confirm the registry's current version and cannot honestly state from this
verification that the npm channel did or did not ship.
