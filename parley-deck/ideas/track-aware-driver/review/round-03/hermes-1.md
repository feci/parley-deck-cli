---
agent: hermes-1
idea: track-aware-driver
review-round: 3
date: 2026-07-03
reviewed-commit: ce4ca22
---

## Summary

Fix-up cycle 2 (commit ce4ca22) is a two-hunk diff that closes the two open
items from round 2: (1) the codex-1 BLOCK on the `acquireLock` TOCTOU race in
`internal/driver/loop.go`, and (2) my own round-02 MINOR finding that the
fast-track model-diversity hard gate had no dedicated test. Both are addressed
correctly and minimally — no drive-by changes, both hunks trace directly to the
two findings. I verified the lock fix by reasoning through every interleaving,
ran the exclusivity test 50x, and ran the full suite. Everything is green.

## Verification

### Lock TOCTOU fix (codex-1 round-02 BLOCK) — RESOLVED

The diff at `internal/driver/loop.go:253-260` changes the EEXIST-retry branch:
previously, when `strconv.Atoi` on the lock content FAILED (empty or
unparseable), the code fell through to `os.Remove` + retry — reclaiming the
half-written lock as "stale". Now a parse failure returns immediately with
"driver.lock present but not yet owned ... refusing", and the `os.Remove`
reclaim path is reached ONLY when `Atoi` succeeds AND the parsed PID is dead AND
different from `os.Getpid()`.

Refutation attempt — can two holders still be granted? Trace every interleaving
of two racers A and B on the same path:

1. A wins `O_EXCL`, writes PID, closes; B hits EEXIST, reads a parseable live
   PID → B returns "held by pid A". One holder. Safe.
2. A wins `O_EXCL`, has NOT yet written PID (file empty/zero-length); B hits
   EEXIST, reads empty data → `Atoi("")` fails → B returns "not yet owned,
   refusing". A later completes its token write and holds the lock. One holder.
   Safe. THIS is the window the old code got wrong (B used to remove+recreate);
   the new code closes it by refusing.
3. A wins `O_EXCL`, writes a PARTIAL/garbage token (e.g. "12" of "1234" due to a
   torn read, or non-numeric corruption); B hits EEXIST, reads garbage →
   `Atoi` fails → B refuses. A holds. One holder. Safe.
4. Stale lock from a dead PID P; A and B both hit EEXIST, both read parseable
   "P", both see `processAlive(P)=false`, both attempt `os.Remove`. Only one
   `Remove` succeeds (the other gets `IsNotExist` and the loop continues);
   exactly one subsequent `O_EXCL` wins. One holder. Safe (this was already
   correct; the fix does not regress it).

There is no interleaving in which an empty/unparseable lock body leads to a
reclaim, because the `perr != nil` branch returns before `os.Remove`. The
reclaim path is now guarded by a successful parse, so the create-then-write
window can no longer be exploited by a second racer. The fix is structurally
correct.

Tradeoff note: a genuinely corrupt (non-numeric, non-empty) lock file will now
be refused permanently rather than reclaimed — but refuse is the fail-closed
mode (no double-drive), which is the correct safety posture for a
single-driver lock, and it matches codex-1's counter-proposal exactly. Not a
defect.

Test: `go test ./internal/driver -run TestAcquireLockIsExclusive -count=50`
→ exit 0, "ok parley-deck-cli/internal/driver 0.243s". The test launches 8
concurrent goroutines racing `acquireLock` on one path and asserts exactly one
winner, then verifies re-acquire works after release. 50 iterations (each with
8 racers) produced zero multi-holder grants. The previous TOCTOU would
intermittently surface under `-count`; it does not here.

### Fast model-diversity hard gate test (hermes-1 round-02 MINOR) — RESOLVED

The diff at `internal/app/driver_impl_le_test.go:203-217` adds
`TestCheckModelDiversityHardGateOnFastTrack`: it writes `track: fast` to the
00-prompt (with NO `require_model_diversity`), builds a same-model roster
(claude=m1, codex=m1), and asserts `checkModelDiversity()` returns an error.
This directly exercises the `if t, present := driver.ReadTrack(o.ideaDir);
present && t == track.Fast { required = true }` branch
(`driver_impl.go:154-160`) that my round-02 review identified as untested. A
regression that removes the `t == track.Fast` check would now fail this test
(the same-model roster without the flag would fall to the warn path, not
error). The test is minimal and uses the existing `writePrompt` /
`newOpsWithStore` helpers — no new machinery. Gap closed.

### Full suite — GREEN

`go test ./...` → exit 0. Every package ok (internal/app 9.173s, internal/driver
0.758s, internal/runner cached). No FAIL lines. The standing sandbox exception
(`internal/runner TestDurableKillEndToEndRealProcess` "no recorded boot id")
does not reproduce here — the runner package is cached-green from a prior run in
this environment, confirming it is environmental, not a code defect.

## New findings (if any)

None. Fix-up cycle 2 is scoped exactly to the two round-2 open items and
introduces no new behavior beyond them. The lock fix is a pure control-flow
tightening (parse-failure → refuse instead of reclaim); the test addition
exercises an existing branch. No new surface area, no new findings.

## Signoff

Status: ✅ ACCEPT

The codex-1 round-02 BLOCK (acquireLock TOCTOU) is resolved — the
empty/unparseable lock body is now treated as HELD (refuse), never reclaimed as
stale, and I cannot construct any two-holders interleaving under the new
control flow. 50x `-count` on the exclusivity test passes cleanly. My round-02
MINOR (no test for the fast model-diversity hard gate) is resolved by
`TestCheckModelDiversityHardGateOnFastTrack`. The full suite is green. All five
original track findings remain resolved from round 2. No new findings. The
implementation is ready.
