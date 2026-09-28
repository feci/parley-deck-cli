---
from: user
to: codex-1
idea: windows-portability
date: 2026-09-28
status: decided
relayed-by: the owner's Claude Code session (not a participant artifact)
---

# Owner answer: containment is preserved — the deviation is NOT approved

This answers `codex-1-to-user_windows-portability_containment-deviation.md`.

Asked "Ako má Windows verzia narábať s ochranou proti úniku z adresára?" ("How should the Windows version handle
the protection against escaping the directory?"), the owner selected:

> "Zachovať ochranu" — "Sedem operácií sa na Windows výslovne odmietne. Dve budgetové funkcie tam nebudú fungovať,
> ale bezpečnostná záruka ostane rovnaká ako na macOS a Linuxe."

"Preserve the protection — the seven operations are explicitly refused on Windows. Two budget features will not
work there, but the security guarantee stays the same as on macOS and Linux."

Keep os.Root containment; refuse the seven rooted sites (A2, A3, B3, B4, B5, C1, C2) and their dependent re-sync
operations before mutation, with a clear user-visible refusal, and disclose the two unavailable Windows features in
the release notes. The path-based MoveFileEx branch is rejected.
