---
from: codex-1
to: all
idea: windows-portability
phase: implementation
blocking: no
date: 2026-10-08
---

## Quota release handoff: new Windows durable sites

The quota-auto-exclude cycle-5 signed reservation G21 routes these actual failures to
the existing windows-portability branch/idea, locally inspected commit
3526b820d24e460debe3e493e9002a1ae7508927 (IMPLEMENTATION in-progress). That idea is
absent from this branch/main. The remote branch head observed on 2026-10-08 is
9a60390cf05e9135b96253e126479150f9158daf; no claim is made here about its newer work.
The branch/idea source is
https://github.com/feci/parley-deck-cli/blob/windows-portability/parley-deck/ideas/windows-portability/00-prompt.md
and the frozen reference is
https://github.com/feci/parley-deck-cli/blob/3526b820d24e460debe3e493e9002a1ae7508927/parley-deck/ideas/windows-portability/FINAL.md.

Its FINAL predates the quota additions. This handoff does not amend that FINAL or
claim it promises working idea creation: its Windows SyncDir rule explicitly refuses
future unconverted sites without a proved mechanism. Please evaluate these sites in
that track under its own review/scope process:

- internal/quota/history.go: DurableWrite, SyncPath and directory sync.
- internal/quota/record.go: WriteKickoff directory sync.
- internal/protocol/quota.go: CreateIdeaWithQuota ideas-directory sync.
- internal/store/events.go: AppendDurable and Sync.
- internal/pidlease/lease.go: publish and reap (publication links the lease before sync).

Evidence: CLI e04852e, GitHub run 37692987015, Windows Test failure:
https://github.com/feci/parley-deck-cli/actions/runs/37692987015.
Stored CI log `.parley-runtime/quota-implementation/review-cycle-4/windows-ci.log`,
SHA256 b36a9e47732dd37de0d1c637bfe9d28004dce4fcc396545674b4000da609de2d.
Canonical evidence: ideas/meta-protocol-change-quota-auto-exclude/source-context/codex-1-cycle4-platform-ci.md
and review/round-07/claude-1.md plus the cycle-5 signed reservations.

Windows CLI 1.51.0 cannot create new ideas (including policy-off) and default scoped
ideas cannot be driven/signed; manual imports, owner revisions and transitions use
the same failing directory sync. This release discloses the behavior, leaves Windows
experimental and CLI winget held, and makes no Windows durability fix. The skill-only
winget release is distinct. No participant in this handoff is being launched.
