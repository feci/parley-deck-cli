---
from: codex-1
to: user
idea: meta-protocol-change-lean-organizer
release: CLI 1.49.1
phase: release-complete
status: complete
blocking: no
date: 2026-09-24
recorded_at: 2026-09-24T21:32:53.714580+00:00
---

# CLI 1.49.1 release complete

The owner-authorized release is complete. zcode-1 updated only the candidate/CHANGELOG wording; kimi-1 independently reviewed that exact commit and then independently audited every release channel. Both verdicts **PASS**, with no findings. No A–D deliberation or product-code work was repeated. The organizer performed release operations, not product-code verification.

## Channel evidence

| Channel | Final state and evidence |
|---|---|
| Git / CLI GitHub | Remote main fast-forwarded directly from `9134c7a` to reviewed commit **`54e07981740d495f08cd90842168d9228c58c90f`**. Immutable `v1.49.1` points to exactly that commit. [Release 1.49.1](https://github.com/feci/parley-deck-cli/releases/tag/v1.49.1) published 2026-09-24T21:18:37Z with **six fresh assets** built from a clean checkout of that commit. No development PR. Subsequent closure records are a records-only descendant, not a tag move. |
| Previous tag | `v1.49.0` remains `06e563e8b1fe8094132149b83e14da7be9aae51e`, independently checked locally/remotely; no published tag moved. |
| Windows scope | Release notes and CHANGELOG say **experimental/unvalidated**; both Windows assets are retained and individually labelled `Windows ARM64 — experimental/unvalidated` / `Windows x64 — experimental/unvalidated`. Known Windows defects remain for the separate reviewed `windows-portability` idea. |
| Homebrew CLI | [Tap commit 52b09f3](https://github.com/feci/homebrew-parley/commit/52b09f39c8f4a98612dbf052b1b417bd2f043eb2) changes only CLI formula URL + SHA256. URL: `https://github.com/feci/parley-deck-cli/archive/refs/tags/v1.49.1.tar.gz`. Independently downloaded archive SHA256 **`6c0760611a4c4242418015c70e7de1508eaacbaf763a80c27a5a549a1b7a31f4`**, matching the formula. |
| Local CLI | `brew upgrade feci/parley/parley-deck-cli` completed. `which parley` = `/opt/homebrew/bin/parley`, resolving to `/opt/homebrew/Cellar/parley-deck-cli/1.49.1/bin/parley`; independent participant check reports `parley 1.49.1`, and the final organizer `parley --version` also printed **`parley 1.49.1`**. |
| Skill GitHub/Homebrew | Unaffected at **2.13.0**. GitHub tag `v2.13.0` remains `8161e5e9d3ecdece25b7ae225e2304d3159be79b` with five assets. Skill formula unchanged; its real archive hash independently matches `0d585386219fe1939bd37fd27ba697b25f5b932ae7415657a16833970e44ff44`. Installed skill remains2.13.0. |
| CLI winget | **HELD by owner decision. No CLI1.49.0/1.49.1 PR submitted or open**, independently searched in microsoft/winget-pkgs. No CLI winget mutation by this organizer. |
| Skill winget | [PR440360](https://github.com/microsoft/winget-pkgs/pull/440360) **MERGED**, 2026-09-24T11:16:23Z. Current state reported only. |
| npm skill | Final `npm view parley-deck-skill version dist.integrity dist.tarball --json` still reports **latest2.12.1**. Skill2.13.0 publication remains with the owner-facing Claude session and is **not a CLI release gate**. No new login/publish flow started. |

## Independent evidence

- Wording author: `ideas/meta-protocol-change-lean-organizer/release-wording-1.49.1-zcode-1.md`.
- Exact-commit wording review: `ideas/meta-protocol-change-lean-organizer/release-wording-1.49.1-review-kimi-1.md` — **PASS**.
- Final channel audit: `ideas/meta-protocol-change-lean-organizer/release-1.49.1-audit-kimi-1.md` — **PASS**, no findings.
- Operations record: `ideas/meta-protocol-change-lean-organizer/release-1.49.1-orchestration-codex-1.md`.
- Downloads, clean source checkout, build script/log/provenance, notes, formula/archive evidence, and independent downloads: `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-24-lean-organizer/release-1.49.1` (independent audit subdirectory `audit-kimi-1/`).

All six SHA256 values below independently match the staged assets, the GitHub API digests, and the participant's downloaded bytes. Every downloaded binary embeds `vcs.revision=54e07981740d495f08cd90842168d9228c58c90f` and `vcs.modified=false`, with the correct GOOS/GOARCH. Windows executables were not executed or claimed validated.

| Asset | SHA256 |
|---|---|
| `parley-v1.49.1-darwin-arm64` | `70545193472b05df484d5bab339c3ea14f7f3079f3a452f109185a78cf7f56f5` |
| `parley-v1.49.1-darwin-x64` | `5007e4cc2169bd4407520b371309ec4c9cdf9c67ae0124b09c0c6d20761f6a9d` |
| `parley-v1.49.1-linux-arm64` | `f489bcbab967817f789a83e1dabcff66ee49c06ed8c0b806412d80038c383c7c` |
| `parley-v1.49.1-linux-x64` | `502cef4ba1f54c16555cb4c318a274d292621de14505d0e5b99d96313b71ef15` |
| `parley-v1.49.1-windows-arm64.exe` | `2dcc2c60c5078e8ce1333b445a61b6465bd64a8620067d6c051b441f8dbbe76b` |
| `parley-v1.49.1-windows-x64.exe` | `b6cf7644d6770af446acb5d957eb156befc3271720ab160fc1ae0ca61739cd41` |

## Remaining owner work and release order

The owner-facing Claude session retains responsibility for logging in/publishing the exact audited skill2.13.0 npm tarball. Latest observed registry metadata remains2.12.1, integrity `sha512-45s2sww7d7LVtKgIJ9jaOyHvfl1f70mk1r8/luw+4B+8eU2hT5Wwx5L7HVWQLvmtdENMi0Hp4JczEtISuPGmug==`, tarball `https://registry.npmjs.org/parley-deck-skill/-/parley-deck-skill-2.12.1.tgz`. No further owner input is needed for CLI1.49.1.

This exact completion file releases the sequencing hold for `meta-protocol-change-designated-implementer`. `windows-portability` releases after that idea's completion file. Until Windows portability ships, every CLI release retains the experimental/unvalidated Windows label and keeps CLI winget held. This release did not start either separate idea or alter their work.

## Organizer accounting

Own resume-session usage appended to `ideas/meta-protocol-change-lean-organizer/organizer-usage.md`, sourced from the pinned Codex rollout, not self-estimated. Accounting timestamp **2026-09-24T21:31:03.218Z**: input **3611928**, cached input **3537152**, output **18928**, reasoning **7640**, **62** accounting request events. Snapshot precedes final recording/push and response; it is separate from the previous organizer session. No raw rollout or credentials committed.

Pre-existing untracked runs and usage-ledger artifacts preserved. Google Chrome was never opened or controlled.
