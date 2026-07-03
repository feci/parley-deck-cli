---
agent: codex-1
idea: track-aware-driver
review-round: 3
date: 2026-07-03
reviewed-commit: ce4ca22
---

## Summary

Fix-up cycle 2 resolves my round-02 block. `acquireLock` no longer treats empty or unparseable lock content as stale: a contender that observes a just-created lock file before the owner writes its PID now fails closed instead of removing the file.

I could not construct a normal two-holder interleaving after the change. The stale-PID path can still be raced by multiple contenders, but after removal they all return through the second `O_EXCL` create attempt; only one can win, and the losers then see the winner's live or self PID and refuse. The previous empty-file TOCTOU path is closed because unparseable content is now treated as held.

Hermes's round-02 minor is also resolved: the fast-track model-diversity hard gate is now covered by `TestCheckModelDiversityHardGateOnFastTrack`.

## Verification

- Reviewed my round-02 review, `review/round-03/DIFF-fixups.txt`, live `internal/driver/loop.go`, and live `internal/app/driver_impl_le_test.go`.
- `go test ./internal/driver -run TestAcquireLockIsExclusive -count=50` passed.
- `go test ./internal/app -run TestCheckModelDiversityHardGateOnFastTrack -count=1` passed.
- `go test ./...` failed only at the known codex-sandbox limitation in `internal/runner TestDurableKillEndToEndRealProcess`: `process verification failed (no recorded boot id); not killed`. All other packages passed, including `internal/driver` and `internal/app`.

## New findings (if any)

None.

## Signoff

Status: ✅ ACCEPT
