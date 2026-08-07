---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 2
date: 2026-08-07
reviewed-commit: 8888e00
verdict: FINDINGS
---

## Summary

The fix-up repairs the Windows compile regression, direct existing-release/direct-release-symlink
refusal, the two covered renderer cases, literal version traversal, production dispatch, help, and
status read-error handling. It does not reach CLEAN. The new content report is still capable of
silent semantic erasure; symlink handling still does not confine either `Load` or all `Publish`
path components; the authoritative protocol text still claims both agent-proof publishing and
rank-2 continuation; and two of the new regression tests pass when their claimed production fix is
removed.

The required commands produced exactly these results:

```text
$ go build ./...
[no stdout/stderr]
exit 0

$ GOOS=windows go build ./...
[no stdout/stderr]
exit 0

$ GOOS=linux go build ./...
[no stdout/stderr]
exit 0

$ go test ./...
?   	parley-deck-cli/cmd/parley	[no test files]
ok  	parley-deck-cli/internal/acp	(cached)
ok  	parley-deck-cli/internal/agents	(cached)
ok  	parley-deck-cli/internal/app	44.554s
ok  	parley-deck-cli/internal/config	(cached)
ok  	parley-deck-cli/internal/consensus	(cached)
ok  	parley-deck-cli/internal/driver	1.168s
ok  	parley-deck-cli/internal/fsutil	(cached)
ok  	parley-deck-cli/internal/hitl	(cached)
ok  	parley-deck-cli/internal/loop	(cached)
ok  	parley-deck-cli/internal/pipeline	(cached)
ok  	parley-deck-cli/internal/procctl	1.717s
ok  	parley-deck-cli/internal/protocol	(cached)
?   	parley-deck-cli/internal/protocolcore	[no test files]
ok  	parley-deck-cli/internal/repomap	(cached)
ok  	parley-deck-cli/internal/retro	(cached)
ok  	parley-deck-cli/internal/runaction	(cached)
ok  	parley-deck-cli/internal/runcontrol	(cached)
ok  	parley-deck-cli/internal/runmanifest	(cached)
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL	parley-deck-cli/internal/runner	7.815s
ok  	parley-deck-cli/internal/runplan	(cached)
ok  	parley-deck-cli/internal/runstate	(cached)
ok  	parley-deck-cli/internal/sessionstore	(cached)
ok  	parley-deck-cli/internal/steer	(cached)
ok  	parley-deck-cli/internal/store	(cached)
ok  	parley-deck-cli/internal/track	(cached)
ok  	parley-deck-cli/internal/tui	0.403s
FAIL
exit 1
```

`TestDurableKillEndToEndRealProcess` fails identically at `4396529` when run alone, so I do not
attribute it to this delta. The six newly added tests all pass at `8888e00`.

## Round-01 findings: fixed or not

- **Windows build — FIXED for compilation.** All three requested build commands exit 0. The
  unsupported-platform file defines `hasTTYSupported = false`, so the implementation is visibly
  fail-closed in source. There is still no Windows/runtime end-to-end test proving the publisher's
  refusal; a cross-build proves only that it compiles.

- **Existing release directory and direct release-path symlink — FIXED for those exact cases.**
  Through the real publisher under a PTY, a pre-existing `occupied/` directory and a version path
  that is itself a symlink both returned the write-once refusal, and neither wrote a body. The
  broader symlink guarantee is not fixed; see the new MAJOR finding.

- **TTY bypass / overclaim — NOT FIXED as documentation.** The no-TTY refusal now honestly says a
  PTY defeats it, and `IMPLEMENTATION.md` says the same. But the live protocol still says, “An agent
  may not ... publish one” because the command “refuses without a controlling terminal”
  (`parley-deck/COOPERATION.md:759-762`). The embedded copy and the newly committed skill copy at
  `parley-deck-skill@455aafe` repeat that text. `meta/protocol-changelog.md:23-25` and
  `IMPLEMENTATION.md:74-75` also still claim “no agent-accessible write path.” PRIMARY evidence:
  the real binary under `script(1)` published `pty-agent`, printed `Published core pty-agent ...`,
  exited successfully, and created the release. The source comment at
  `internal/app/protocol.go:277-279` is likewise still false. This remains the round-01 CRITICAL
  truthfulness/G2 finding, not an honestly closed disposition.

- **Heading-based removal and CRLF — PARTIALLY FIXED.** The shared-heading regression test fails
  when the old heading implementation is restored, and the CRLF test fails with both the old false
  report and mixed endings. Those two cases are genuine. `droppedContent`, however, is a global
  set-of-lines comparison rather than a block/content comparison and still misses removals; see the
  new CRITICAL finding.

- **Committed-lock traversal — PARTIALLY FIXED.** A real fixture whose lock contained
  `core-version: ../escape` failed with `protocolcore: unsafe version "../escape"`, even though a
  matching outside file existed. A valid version whose release directory is a symlink still escapes
  the store, and the new lock-traversal test does not prove the validator; see below.

- **Production dispatch — FIXED.** `TestProtocolIsReachableThroughProductionDispatch` goes through
  `Run`, exercises both status and render, and fails when only the `case "protocol"` dispatch is
  removed.

- **Rank-2 continuation wording — NOT FIXED globally.** `IMPLEMENTATION.md:15-27` now says rank 2 is
  absent. All three shipped protocol copies still say an open idea completes under its pinned
  version (`parley-deck/COOPERATION.md:766-767`). PRIMARY evidence: `protocolcore` has no production
  callers outside `internal/app/protocol.go`; there is no snapshot writer/resolver or continuation
  integration. No end-to-end continuation test exists. This remains the round-01 MAJOR/G7b
  finding.

- **Third protocol copy — FIXED as delivery, but it carries the two false guarantees above.**
  Sibling commit `455aafe9f99fd6c01223b920a0768af2119e14a3` exists, is HEAD of
  `../parley-deck-skill`, and changes exactly
  `skills/parley-deck/references/COOPERATION.md`. The added clause is byte-for-meaning the stale
  agent-proof-publish and rank-2 text, so committing it closes the missing-artifact finding but
  expands the G7b problem.

- **Top-level help — FIXED.** The real `parley --help` output now includes
  `parley protocol publish --version V --from FILE`.

- **Status read errors — FIXED behaviorally.** With the store directory at mode `000`, the real
  command exited 1 and printed `protocol status: reading ...: permission denied`; an unreadable
  deck lock behaved the same way.

- **Unreadable render input from kimi-1 — NOT FIXED and not listed among the claimed fixes.**
  `protocolRender` still discards the `os.ReadFile` error at `internal/app/protocol.go:180`. In a
  real fixture with an unreadable existing deck, `render --yes` exited 0, replaced the file, and
  lost the deck's Workspace, Transport, and content without a removal report. This remains at least
  MINOR robustness debt and is another route around G1's report.

## New findings (by severity, or "none")

### [CRITICAL] `droppedContent` loses section context and multiplicity, so G1 still permits silent semantic erasure

`droppedContent` builds `map[string]bool` from every trimmed rendered line and treats every matching
deck line as carried (`internal/protocolcore/render.go:195-219`). It does not bind a line to its
heading, position, or occurrence count.

PRIMARY fixture: the core and deck both contained these two headings. The core allowed
`- Deploy production.` under `## Allowed automatically`; the deck also carried that same line under
`## Requires explicit user approval`. Preview printed no removal section at all. Apply removed the
line from `Requires explicit user approval` and retained only the occurrence under
`Allowed automatically`:

```text
preserved from this deck: Workspace, Transport, Created
would regenerate ... from core 1.0.0 (...).
Nothing was written. Re-run with --yes to apply.

8:## Requires explicit user approval
10:## Allowed automatically
12:- Deploy production.
```

The identical text surviving in a different block masks the loss of the binding restriction. A
duplicate line in one section has the same multiplicity bug. This is the exact silent-erasure class
G1 forbids, with a semantic inversion rather than merely an imprecise count.

Suggested fix: compare typed blocks/section occurrences, not a document-global set. Until the rank-1
registry exists, at minimum compute an order- and multiplicity-aware diff within each heading path,
and add an end-to-end fixture where identical text occurs in two differently governed sections.

### [MAJOR] Symlink confinement is applied only to the final create path; `Load` and parent components still follow links

PRIMARY, read path: I created a valid lock `core-version: safe-1`, made
`<store>/safe-1` a symlink to an outside directory containing `COOPERATION.md`, and ran the real
`protocol render --yes`. It exited 0, reported core `safe-1`, and replaced the deck with
`SYMLINK-ESCAPED-CONTENT`. `ValidVersion` validates the string, while `Store.Load` then calls plain
`os.ReadFile` (`internal/protocolcore/core.go:82-94`), which follows the release symlink.

PRIMARY, write path: with `$PARLEY_HOME/protocol` symlinked to an outside directory, the real PTY
publisher printed success for version `escaped` and created
`outside/core/escaped/COOPERATION.md`. `Lstat` checks only the version directory before creation,
and `O_NOFOLLOW` applies only to the final `COOPERATION.md` component
(`internal/protocolcore/core.go:130-159`). It does not reject symlinks in ancestors, and the
Lstat/MkdirAll/Open sequence also leaves an unbound replacement window.

The direct pre-existing version symlink is now refused, so this is narrower than round 1, but the
claim that publishing “can never write through a planted symlink” (`IMPLEMENTATION.md:21-24`) is
still false and a safe-looking committed lock can still make render consume bytes outside the
release store.

Suggested fix: anchor operations to a trusted, resolved store directory and walk/open every
component without following links (for example dirfd/openat-style traversal on supported systems),
apply the same rule to `Load`, and fail closed where the platform cannot prove it. Test a valid
version-directory symlink on Load, an ancestor symlink on Publish, and a directory-replacement race.

### [MINOR] `ValidVersion` validates a trimmed copy but uses the untrimmed path

`ValidVersion` trims only its local variable (`core.go:48-67`); `Publish` later passes the original
string to `ReleaseDir`. PRIMARY evidence: the real PTY publisher accepted `--version ' padded '`,
printed success, and created a directory whose name retained both spaces. The deck-lock reader
trims its value, so the normal real command cannot subsequently name that release exactly. The
write-once publisher can therefore report success for an effectively unusable namespace.

Suggested fix: either reject when `strings.TrimSpace(v) != v`, or canonicalize once and use the
canonical value consistently for validation, storage, return values, display, and lock parsing.

### [MINOR] The Phase-8 implementation record was not advanced to the reviewed fix-up

`IMPLEMENTATION.md` still has top-level `status: implemented` and
`head-commit: (see release commit)` (`IMPLEMENTATION.md:1-10`), and the `## Fix-up cycle 1` section
omits `head-commit`. Phase 8 requires `status: fix-up-cycle-1` and the new SHA. The commit message
calls the cycle ready for re-review, but the machine-readable artifact does not identify
`8888e00` or the current phase.

## Test-quality assessment

The six newly added tests pass at `8888e00`, but the revert audit found two tautologies and several
uncovered fixes:

1. `TestProtocolRenderReportsContentLostUnderASharedHeading` is revert-sensitive: restoring the
   heading-only renderer makes it fail. It checks only that a heading appears in output, not that
   the specific lost content/count is faithfully represented or that apply loses exactly what was
   previewed.
2. `TestProtocolRenderHandlesCRLFDecks` is revert-sensitive: restoring the old renderer makes both
   its false-removal and mixed-ending assertions fail.
3. **Tautology:** `TestProtocolRejectsPathTraversalInTheLock` still passes after removing the
   `ValidVersion` guard from `Store.Load`. Its invalid paths point to no fixture release, and it
   asserts only a non-zero render result, so the later “release not installed” error satisfies the
   test just as well as validation. It must create an outside `COOPERATION.md`, assert the exact
   unsafe-version refusal, assert unchanged deck bytes, and include a valid-version symlink.
4. `TestPublishRefusesExistingReleaseDirAndSymlinks` catches removal of the directory-level Lstat:
   that mutation makes the existing-directory assertion fail. **Its symlink half is tautological
   for the open flags:** the test first creates the version directory, so the new directory Lstat
   returns before `OpenFile` reaches the planted file symlink. Removing `O_NOFOLLOW` leaves the test
   green. It also does not test a version-directory symlink, an ancestor symlink, or the real
   publisher.
5. `TestPublishRejectsUnsafeVersions` is revert-sensitive for Publish: restoring the old
   slash/`..` check makes `.` and `.hidden` fail the test. Despite its comment, it never calls
   `Load`; it also misses padded versions.
6. `TestProtocolIsReachableThroughProductionDispatch` is genuine and revert-sensitive: deleting
   only production dispatch makes it fail with `unknown command: protocol`.

No committed test exercises `protocol publish` through `Run` or a compiled process; none proves
the ordinary no-TTY refusal, the documented PTY limitation, the unsupported-platform refusal, or
the symlink boundary. The Windows cross-build catches the former compile regression but is not an
end-to-end refusal test. The help and status fixes are behaviorally correct but have no new
regression test.

G7b is therefore still violated by text, independently of the new code findings:

- all three protocol copies claim an agent cannot publish, while a PTY publishes;
- all three protocol copies claim open-idea pinned continuation, while rank 2 has no code or test;
- the changelog claims shipped `DETECTED-UNATTRIBUTED`, which appears in no Go implementation;
- `IMPLEMENTATION.md` claims every removed/replaced block is reported, which the scoped-duplicate
  fixture disproves;
- `IMPLEMENTATION.md` claims atomic mode-preserving render and fail-closed unsupported-platform
  publishing without a real-entry end-to-end test for either behavior.

These are guarantees, not merely missing nice-to-have coverage. They must either be made true and
tested at the production boundary or removed/narrowed from the landed text before CLEAN.
