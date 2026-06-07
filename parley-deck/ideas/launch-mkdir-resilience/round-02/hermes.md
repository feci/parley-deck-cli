---
agent: hermes
idea: launch-mkdir-resilience
round: 2
date: 2026-06-07
---

## Convergence on core design
All four round-01 files converge on `internal/fsutil.MkdirAllResilient` (first os.MkdirAll success = immediate return with zero cost; on error: os.Stat+IsDir check; bounded retries; stdlib only). This is adopted at Store.Append + runner/runmanifest/hitl sites. No source changes here.

## Fork resolutions (addressing each agent)

1. BACKOFF (claude 20/40/80, codex 5/20/50, agy 5/15/45/100 vs hermes zero-sleep): Agree with hermes round-01. The retry itself forces cache revalidation on virtio-fs; any sleep adds latency on an already-failed path with no proven benefit for the observed ENOENT window. Reject compromise (immediate first retry + later backoff) — keep zero sleep to honor "no cost on success path" strictly. Retries: 3 total attempts, no sleep.

2. SCOPE (protocol/workspace.go:67 CreateIdea idea-dir mkdir): Agree with codex/hermes round-01 scoping it OUT. protocol.CreateIdea runs before Append inside runcontrol.Create, but the proven failure and all listed runner sites are under the runs/<runID> tree. Adding workspace idea creation would broaden beyond the diagnosed launch failure and the "run-path only" surgical constraint. Keep runner-only + Store.Append + runmanifest + hitl.

3. EEXIST (codex: never blind success, always Stat+IsDir confirmation vs hermes/agy: treat ErrExist immediate success): Agree with codex. A file collision at the target path must surface the real error; blind ErrExist acceptance risks masking ENOTDIR. Safe rule: always confirm via os.Stat+IsDir after any error (including ErrExist).

4. FAIL-FAST on fs.ErrPermission (agy): Disagree. Permission errors are persistent and must still exhaust the (tiny) retry window so the final returned error is the authentic last MkdirAll error. Fail-fast would hide the real errno from callers. Bounded retries already prevent long hangs.

5. OPEN FILE trailing os.OpenFile in Append: Confirm unchanged. Once MkdirAllResilient has confirmed the directory via Stat, the subsequent O_CREATE|O_APPEND open succeeds; retrying the open risks event duplication. No change needed.

6. TEST PLAN: Confirm exactly the listed fsutil seam tests (transient-then-success; dir-exists-via-Stat without exhausting retries; already-exists; genuine-failure-returns-last-mkdir-error; fail-fast-permission; common-path-no-stat-no-sleep). Use package vars for seams; keep sleep no-op in tests.

## Recommendation
Ship fsutil helper (3 attempts, zero backoff, Stat+IsDir after every error, ErrExist not blind), apply at agreed run-path sites only. This closes all forks decisively while staying minimal.