---
from: claude-1
to: user
idea: meta-protocol-change-participants-optional
phase: implementation
blocking: no
date: 2026-10-31
---

CLI **1.57.0** and skill **2.21.0** are released. Core **2.21.0** is staged for owner publication.
This was a solo run under the recorded owner exception. Independent review: none.

| Channel | Result | Evidence |
| --- | --- | --- |
| CLI GitHub release | Published, `v1.57.0` | https://github.com/feci/parley-deck-cli/releases/tag/v1.57.0; PR 88; merge `d05158a` |
| Skill GitHub release | Published, `v2.21.0` | https://github.com/feci/parley-deck-skill/releases/tag/v2.21.0; PR 14; merge `47f0215` |
| Both Homebrew formulae | Published, upgraded, tested, `brew audit --strict` passed | Tap commit `378135b`; `parley 1.57.0`, skill `2.21.0` |
| Local runtimes | 114 SKILL.md files match the release package (19 parley-deck copies, sha256 `3e714ffd8d30…`) | `install --target all --include-undetected --force`, four generic destinations, `doctor` ok |
| Skill WinGet | PR submitted | https://github.com/microsoft/winget-pkgs/pull/466936 |
| CLI WinGet | Held while CLI Windows is experimental | none |
| npm and global core | Owner-only, pending for 2.15.0–2.21.0 | `runs-handoff/owner-publish-20261013.sh` now covers 2.21.0 |

CLI asset SHA256 values are in the release `sha256.json`. Skill tarball SHA256:
`c9b3960b5dd1dd8437fea56a63754aab355f091931e720ff08e8f4d61365181b`. Staged core 2.21.0 SHA256:
`dcc4b14e1334bf2710d21d276c2fcac5a32763b73cc737b27a57d00c16568e5e`.

Validation: `go test ./... -count=1 -timeout 45m` passed all 34 packages at `9db33f2`. Skill `npm test` passed.
PR CI passed on macOS and Ubuntu; CLI Windows CI stays experimental.

Owner command, in a real terminal (Terminal.app, not the `!` prefix):

```sh
bash '/Users/tomasfecko/AI_WORKSPACE/parley-deck/runs-handoff/owner-publish-20261013.sh'
```

Deferred: an independent review of this solo change by a later idea; a driver auto-close with the
`independent review: none` marker (only if the owner asks); the earlier deferred items of
`meta-protocol-change-credit-exhaustion-failover`.
