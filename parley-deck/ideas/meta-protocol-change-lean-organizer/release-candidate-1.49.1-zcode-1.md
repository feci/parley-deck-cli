---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: participant implementer (release-candidate metadata prep)
artifact: release-candidate-1.49.1
date: 2026-09-24
status: CANDIDATE PREPARED — reviewable only, NOT released
base: 519951e (lean-organizer branch; hosted-accepted behavior at 9134c7a + comment-only 519951e, equivalence independently verified)
amended: 2026-09-24 — two record corrections before independent review (see § 8): U2 read-error/cancellation wording (reaped pids exit at once; killed-but-unreaped pids consume the full bound, per accepted C-1/K-1 evidence) and Windows decision framing aligned with the organizer scope inbox (defer-and-label-experimental vs separate reviewed portability track; omission of Windows assets is a deviation the owner must select; publication held pending the decision)
amended-2: 2026-09-24 (later, wording finalization) — owner authorization recorded (see § 9 and release-wording-1.49.1-zcode-1.md): release CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew, keep Windows assets explicitly labelled experimental/unvalidated, hold the CLI winget submission until a separate reviewed windows-portability track removes the label; skill 2.13.0 npm handled separately by the owner-facing Claude session and not a CLI gate; dated historical statements (incl. § 8) unchanged
---

# CLI 1.49.1 release candidate — metadata prep — zcode-1

This prepares a **reviewable release candidate only**. Nothing here publishes,
tags, pushes, builds distribution assets, or installs anything. Per the release
plan's role split, the organizer independently reviews this metadata and stages
assets. The required **owner Windows scope decision** has since been given
(2026-09-24, deck inbox
`user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md`):
release CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew, keep
the Windows assets explicitly labelled experimental/unvalidated, and hold the
CLI winget submission until a separate reviewed windows-portability track
fixes the defects and removes the label; the **owner npm verification** of
the separate skill 2.13.0 tarball is handled by the owner-facing Claude
session apart from this release and is not a CLI gate. This commit updates
the candidate's active wording to that authorization (§ 9); it performs no
release operations, and preparing or amending this candidate does not mark
the release task complete.

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
- Owner decision, received 2026-09-24 after the candidate and its § 8
  corrections: `parley-deck/inbox/user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md`
  — "Oboje" / "Both": release CLI 1.49.1 now for macOS and Linux through
  GitHub and Homebrew; label Windows experimental and hold winget for the
  CLI; open a separate reviewed Windows idea that fixes the defects and
  removes the label. Release order set by the owner-facing session:
  release-1.49.1 first (writes its done file), then
  meta-protocol-change-designated-implementer, then windows-portability. The
  skill is unaffected; the npm item above is handled separately by the
  owner-facing Claude session and is not a CLI gate.

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
  **first** non-empty read (no value shopping), and return immediately on
  read errors: a **reaped** pid's ENOENT exits at once, while a
  **killed-but-unreaped** pid's cmdline reads empty with no error and
  therefore consumes the full 100 ms bound (`cmdlinePublishBound`, 1 ms
  poll) before failing closed **exactly as pre-poll** — cancellation is
  bounded, not immediate (accepted C-1/K-1 evidence).
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
  exercised. The entry originally aligned the decision framing with the
  organizer scope inbox
  (`codex-1-to-user_meta-protocol-change-lean-organizer_scope.md`) — owner
  choice pending, publication held — which was accurate as of this
  candidate's § 8 corrections. The owner has since decided (2026-09-24,
  `user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md`):
  **both** — this release ships now for macOS and Linux through GitHub and
  Homebrew; the Windows assets are kept and explicitly labelled
  experimental/unvalidated on the evidence above; and the CLI winget
  submission is held until a separate reviewed windows-portability track
  fixes the defects, turns hosted windows-latest CI green, and removes the
  label. The entry now records that authorization; until that track ships,
  every CLI release keeps the experimental Windows label and opens no CLI
  winget submission.

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
- `git diff -- CHANGELOG.md` vs pre-candidate HEAD → **0 deleted lines**
  (pure insertion; the frozen 1.49.0 entry untouched in place — re-verified
  at the § 8 record-correction commit, which edits only this 1.49.1
  section's own wording, never the 1.49.0 entry).
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
- Windows: red, uninvestigated here. No longer owner-pending: authorization
  recorded 2026-09-24 (both) — macOS/Linux release now via GitHub and
  Homebrew, Windows assets kept and labelled experimental/unvalidated in the
  notes, CLI winget held until the separate reviewed windows-portability
  track ships and removes the label. The skill 2.13.0 npm item is
  owner-attended in the owner-facing Claude session and is not a CLI gate.
- Optional and non-blocking per the acceptance review: one `-v` hosted ubuntu
  run would confirm the two environment-guarded procctl tests execute there.

## 6. Remaining gates, exactly (as amended 2026-09-24, after the owner decision)

1. **Organizer**: independent review of this metadata + report — done
   (kimi-1, VERDICT PASS at `481fb65`, recorded by the organizer in
   `868825f`); still open on the organizer side: stage release assets from a
   pinned commit (either `9134c7a` as tested, or push the comment-only
   `519951e` first — no new CI needed for it per the acceptance review's
   token-identity verification).
2. **Owner — RESOLVED 2026-09-24** (previously: Windows scope decision
   pending, publication held on it). Decision recorded in
   `user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md`:
   both — release this CLI now for macOS and Linux through GitHub and
   Homebrew; keep the Windows assets explicitly labelled
   experimental/unvalidated; hold the CLI winget submission until a separate
   reviewed windows-portability track fixes the defects, turns hosted
   windows-latest CI green, and removes the label. Release order per the
   owner-facing session: release-1.49.1 (this release) completes first and
   writes its done file, then meta-protocol-change-designated-implementer,
   then windows-portability.
3. **Owner — reclassified, not a CLI gate**: npm verification/publish of the
   skill 2.13.0 tarball (separate artifact; SHA256 `dcf9c75c…d8d802` per the
   inbox record) is owner-attended in the owner-facing Claude session, runs
   apart from this release, and does not gate it; the skill is unaffected by
   the Windows decision.
4. **Organizer/owner only**: tag `v1.49.1`, channel builds, GitHub release,
   distribution — now owner-authorized for macOS and Linux through GitHub
   and Homebrew, with the Windows assets retained under the
   experimental/unvalidated label and **no CLI winget submission** until
   windows-portability ships. `v1.49.0`/`06e563e` never moves.

## 7. NOT done

No push, tag, distribution build, publication, or install. No product source,
test, or workflow change. No rewrite of the frozen `v1.49.0` tag/history, the
1.49.0 changelog entry, closed idea artifacts (`FINAL.md`,
`IMPLEMENTATION.md`, `consensus.md`, `review/`), or skill 2.13.0 material. No
Windows investigation or claim. No CLI winget submission opened or prepared;
no `release-1.49.1` done-file written — release operations and completion
signaling remain organizer/owner-only. Peer reports, organizer records,
inbox, ledger and `runs/` left as found (untracked). Preparing or amending
this candidate is not release approval and does not complete the release
task.

## 8. Record corrections (this commit, before independent review)

Two wording defects in the original candidate records (commit `3002782`),
corrected in place in the 1.49.1 CHANGELOG section and § 3/§ 4/§ 5/§ 6 above;
no behavior, test, peer, or closed-artifact change:

1. **U2 read-error/cancellation accuracy.** The candidate said read errors
   return immediately "(dead processes are never polled)". Per the accepted
   C-1/K-1 evidence that is only the **reaped** case (ENOENT exits at once);
   a **killed-but-unreaped** pid's cmdline reads empty with no error and
   consumes the full 100 ms bound before failing closed exactly as pre-poll.
   Both records now state: cancellation is bounded, not immediate.
2. **Windows decision framing.** The candidate framed the owner decision as
   "fix the Windows track or ship without Windows assets" and added a
   self-imposed "Windows assets must not ship on this evidence" obligation.
   Corrected to the organizer scope inbox
   (`codex-1-to-user_meta-protocol-change-lean-organizer_scope.md`): the
   owner chooses between deferring the newly exposed native-Windows work
   with Windows explicitly labelled experimental/unvalidated, or authorizing
   a separate reviewed Windows-portability track; neither choice is inferred,
   and omitting Windows assets is a further deviation the owner must
   explicitly select. The records now state that publication is **held
   pending that owner decision** — a process hold the owner adjudicates —
   and the invented no-assets prohibition is withdrawn.

## 9. Owner authorization — wording finalization (this commit)

The owner answered the scope question this candidate was held on (2026-09-24,
`parley-deck/inbox/user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md`,
status authorized, relayed by the owner's Claude Code session): **"Oboje" /
"Both"** — release CLI 1.49.1 now for macOS and Linux through GitHub and
Homebrew; label Windows experimental and hold winget for the CLI; at the same
time open a separate reviewed Windows idea that fixes the defects and removes
the label. Until that track ships, every CLI release labels Windows
experimental/unvalidated in its notes, keeps the Windows assets (labelled),
and does not open a CLI winget PR. The skill is unaffected; its npm
verification/publish stays with the owner-facing Claude session and is not a
CLI gate.

This commit updates this report's active wording (intro, § 1, § 3, § 5,
§ 6, § 7, front-matter amendment note) and the CHANGELOG 1.49.1 "Release
scope and limitations" bullet to that authorization. Wording only — no
behavior, test, product, workflow, or skill change, and no release operation
performed. Dated historical statements are preserved: § 8 still describes
what the earlier record corrections said as of their date (publication held
pending the then-open owner decision — a hold the decision above now lifts),
and the pre-decision gate list as PASS-reviewed by kimi-1 at `481fb65` read:
organizer review + staging; owner Windows scope decision with publication
held on it; owner npm verification of the separate skill tarball;
organizer/owner-only tag/builds/release. Exact edit inventory:
`release-wording-1.49.1-zcode-1.md` (same directory).
