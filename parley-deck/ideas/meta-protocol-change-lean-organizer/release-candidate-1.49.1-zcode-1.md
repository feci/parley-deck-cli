---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: participant implementer (release-candidate metadata prep)
artifact: release-candidate-1.49.1
date: 2026-09-24
status: CANDIDATE PREPARED — reviewable only, NOT released
base: 519951e (lean-organizer branch; hosted-accepted behavior at 9134c7a + comment-only 519951e, equivalence independently verified)
---

# CLI 1.49.1 release candidate — metadata prep — zcode-1

This prepares a **reviewable release candidate only**. Nothing here publishes,
tags, pushes, builds distribution assets, or installs anything. Per the release
plan's role split, the organizer independently reviews this metadata and stages
assets; release operations then wait on the required **owner Windows scope
decision** and **owner npm verification**. Preparing this candidate does not
mark the release task complete.

## 1. Inputs

- `release-linux-hosted-acceptance-claude-1.md` (independent, not authored by
  me): hosted run **36009912946** at `9134c7a` — **VERDICT PASS** for Linux
  U1/U2 on hosted x86_64 ubuntu (full unsuppressed suite, 32/32 packages, both
  original failure signatures absent, no skips in any previously-failing test,
  probe execution corroborated by procctl leg time 0.106s → 0.593s); macOS leg
  green; **Windows leg FAILED and no Windows claim made**. Commit `519951e`
  (comment-only) verified token-identical excluding comments, so the hosted
  green transfers to this candidate's source.
- Prior records: `release-linux-repair-zcode-1.md` (U1), `bc5ff29` diagnostics
  record (U2 evidence preservation), `release-linux-u2-repair-zcode-1.md`
  (U2 fix + post-review corrections), independent reviews
  `release-linux-u2-review-claude-1.md` / `-kimi-1.md`, and the ledger of
  organizer-recorded approvals at `9134c7a`.
- Owner-gate context: `parley-deck/inbox/codex-1-to-user_..._npm-login.md`
  (skill 2.13.0 npm publish is owner-attended and still pending; that tarball
  is a separate artifact — untouched here).

## 2. Metadata delta (the whole candidate)

| file | change |
|---|---|
| `VERSION` | `1.49.0` → `1.49.1` |
| `internal/app/version.go` | `const version = "1.49.0"` → `"1.49.1"` |
| `CHANGELOG.md` | new top section `## 1.49.1 — 2026-09-24` (+65 lines, **0 deletions** — the 1.49.0 section below is byte-identical) |
| this report | new |

No other file changed (`git diff --stat` = exactly these). No product source,
test, workflow, or skill-2.13.0 change; no closed artifact touched.

## 3. What the changelog entry says, in one place each

- **U1 — stderr drain.** `acp.Process` `Stop`/`Wait` reaped (`cmd.Wait()`)
  before the stderr copier goroutine finished; `os/exec` closes the pipe read
  end when `Wait` returns, so a fast exit could truncate captured stderr to
  zero bytes (hosted run 35987916696: 0 of 16384 bytes observed, exit code
  correct). Fix: drain to EOF, then reap. `Stop` remains bounded — its
  kill-timeout branch closes the stderr read end to interrupt a copier blocked
  by a writer outside the process group; unread bytes on that abandon path are
  explicitly discarded, never silently truncated. Post-fix hosted evidence:
  `acp 2.633s ok` with the drain regression tests unskipped.
- **U2 — Linux argv publication race, bounded fail-closed probe.** Linux can
  stat-expose a process before `execve` publishes argv, so the immediate
  post-`Start` read of `/proc/<pid>/cmdline` could return zero bytes for a
  live process; capture recorded an empty Command and attribution then failed
  closed with `no recorded command` (hosted run 36002574211: exit -1, empty
  output/hash, killed at grace). Fix: the Linux capture probe polls `cmdline`
  — retry **only** zero-byte reads of a live-looking process, return the
  **first** non-empty read (no value shopping), return immediately on read
  errors (dead pids never polled), and on exhaustion of the 100 ms bound
  (`cmdlinePublishBound`, 1 ms poll) fail closed **exactly as pre-poll**.
  All attribution facets, refusal strings, and strictness byte-identical;
  argv publication is monotone and pid swaps remain caught by the exact
  start-time/pgid facets; darwin (ps-based) and Windows untouched.
  Independently reviewed twice and hosted-confirmed on x86_64.
- **Retained executor-error diagnostics.** A non-`ExitError` run failure with
  zero captured output now persists its scrubbed, bounded reason into step
  diagnostics (was: unexplained exit -1 with empty diagnostics).
  `ExitError`-with-empty-output keeps empty diagnostics so the exit code stays
  the discriminator; hashed fields never see the executor text. This is the
  change that made U2's hosted signature *readable*; it is also what the
  hosted green run exercised via the retained-diagnostics tests.
- **Windows claim correction (1.49.0 said: "Windows behavior of `wait`/`usage`
  is exercised by the CI leg").** That claim is false and the 1.49.1 entry
  corrects it in plain words: the first actual hosted Windows execution fails
  — `internal/evidence` does not build (`syscall.Mkfifo` undefined on
  windows), 14 packages fail against 16 ok, the bounded stderr-drain guard
  fires in `internal/acp` — so Windows `wait`/`usage` behavior is **not**
  exercised. The entry states the open **owner scope decision** (fix the
  Windows track vs ship without Windows assets) as neither resolved nor
  waived, and that Windows assets must not ship on this evidence.

## 4. Validation run for this candidate (focused, by me, 2026-09-24, go on darwin/arm64)

Version/packaging checks only — source suites and hosted evidence already
exist independently and were deliberately not repeated:

- `gofmt -l internal/app/version.go` → clean.
- `go build ./...` → ok.
- `go vet ./internal/app/...` → ok.
- `go test ./internal/app/ -run 'TestVersionCommandPrintsSemanticVersion'
  -count=1` → PASS (pins both `parley version` and `parley --version` to
  `versionLine()`).
- `go run ./cmd/parley --version` → **`parley 1.49.1`**.
- Consistency: `VERSION` = `version.go` = CHANGELOG heading = `1.49.1`.
- `git diff -- CHANGELOG.md` → **0 deleted lines** (pure insertion; the frozen
  1.49.0 entry untouched in place).
- `git rev-parse v1.49.0` = `06e563e…` locally **and** at
  `git ls-remote --tags origin` — unmoved; **no `v1.49.1` tag exists** locally
  or remotely.

Not rerun here (evidence already on record): full Go suites (hosted ubuntu
32/32 green + macOS green at run 36009912946; local full-suite pass recorded
at the 1.49.0 metadata commit), U2 container mutation proofs, Windows leg
(failed at 107667649460 — owner-blocked).

## 5. Limits carried into the release notes

- Bound exhaustion is **silent** — degrades to the pre-fix fail-closed
  refusal, never to an unsafe grant. Operational rule recorded in the
  changelog: if `no recorded command` recurs hosted, raise
  `cmdlinePublishBound` before re-diagnosing; observability needs its own
  version.
- The 100 ms bound is measured-tail-based (~1.8× worst of 10,400 spawns,
  claude-1's measurements), not kernel-instrumented.
- Hosted acceptance is one green sample; the U2 claim rests on the
  twice-reproduced mechanism plus mutation evidence, with the hosted run
  closing the x86_64 gap.
- Windows: red, uninvestigated here, owner scope decision pending; npm
  verification for the separate skill 2.13.0 tarball also pending (owner).
- Optional and non-blocking per the acceptance review: one `-v` hosted ubuntu
  run would confirm the two environment-guarded procctl tests execute there.

## 6. Remaining gates, exactly

1. **Organizer**: independently review this metadata + report; stage release
   assets from a pinned commit (either `9134c7a` as tested, or push the
   comment-only `519951e` first — no new CI needed for it per the acceptance
   review's token-identity verification).
2. **Owner**: Windows scope decision — fix the Windows track or ship without
   Windows assets on the corrected claim. Not resolved, not waived.
3. **Owner**: npm verification/publish of the skill 2.13.0 tarball (separate
   artifact; SHA256 `dcf9c75c…d8d802` per the inbox record).
4. **Organizer/owner only**: tag `v1.49.1`, channel builds, GitHub release,
   distribution. `v1.49.0`/`06e563e` never moves.

## 7. NOT done

No push, tag, distribution build, publication, or install. No product source,
test, or workflow change. No rewrite of the frozen `v1.49.0` tag/history, the
1.49.0 changelog entry, closed idea artifacts (`FINAL.md`,
`IMPLEMENTATION.md`, `consensus.md`, `review/`), or skill 2.13.0 material. No
Windows investigation or claim. Peer reports, organizer records, inbox, ledger
and `runs/` left as found (untracked). Preparing this candidate is not release
approval and does not complete the release task.
