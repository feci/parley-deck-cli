---
from: user
to: codex-1
idea: meta-protocol-change-lean-organizer
date: 2026-09-24
status: authorized
relayed-by: the owner's Claude Code session (not a participant artifact)
---

# Owner decision: Windows scope for CLI 1.49.1

This answers `codex-1-to-user_meta-protocol-change-lean-organizer_scope.md`.

## The owner's decision (verbatim, Slovak) and translation

Asked how to handle the Windows defects exposed by the new CI (applies to CLI 1.49.1 and to the
designated-implementer release), the owner selected the option:

> "Oboje" — "CLI 1.49.1 hneď vydať pre macOS a Linux cez GitHub a Homebrew, Windows označiť ako
> experimentálne a winget pre CLI podržať. Zároveň otvoriť samostatnú Windows ideu s review, ktorá chyby
> opraví a označenie odstráni."

"Both — release CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew, label Windows as
experimental, and hold winget for the CLI. At the same time open a separate reviewed Windows idea that fixes
the defects and removes the label."

## Release order (to keep main linear; set by the owner-facing session)

1. `release-1.49.1` (resumed lean-organizer release) — writes
   `codex-1-to-user_release-1.49.1_done.md` in the lean-organizer worktree deck inbox when complete.
2. `meta-protocol-change-designated-implementer` — releases only after (1)'s done file exists.
3. `windows-portability` — releases only after (2)'s done file exists; it removes the experimental label
   and submits the held CLI winget PR once hosted windows-latest CI is green.

Until (3) ships, every CLI release labels Windows experimental/unvalidated in its notes, keeps Windows assets
(labelled), and does NOT open a CLI winget PR. The skill is unaffected.
