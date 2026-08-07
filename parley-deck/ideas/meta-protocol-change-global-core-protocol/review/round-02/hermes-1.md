---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
review-round: 2
date: 2026-08-07
reviewed-commit: 8888e00
verdict: FINDINGS
---
## Summary

Cycle 1 fixed the build and most of the code-level defects I reported. All three GOOS targets
compile, the 15-test suite passes, and the six new tests are real (each fails if its fix reverts).
The symlink/existing-directory hardening, ValidVersion, the CRLF normalization, the content-based
removal report, and the Run-level dispatch test all hold under adversarial probing.

But two of the round-01 findings I and the other reviewers raised were not actually fixed in the
place the overclaim lives, and one was not fixed at all. The §7 protocol text still ships a
rank-2 continuation guarantee as present-tense fact with no implementation and no test. The
changelog still cites DETECTED-UNATTRIBUTED as part of "the shipped guarantee" with zero Go code
behind it. The TTY-gate refusal — the one G7b guarantee all three reviewers called most testable
— still has no test driving `protocolPublish` through any entry point. These are the same G7b
class the project shipped before, and they remain because the fix-up cycle corrected the
IMPLEMENTATION.md narrative while leaving the shipped protocol text and changelog untouched.

## Round-01 findings: fixed or not

**[CRITICAL] GOOS=windows build broken (all 3).** FIXED. `termios_other.go` restored with
`//go:build !darwin && !freebsd && !netbsd && !openbsd && !linux`, `hasTTYSupported = false`,
`platformHasTTY() = false`. Verified: `go build ./...` exits 0; `GOOS=windows GOARCH=amd64 go
build ./...` exits 0; `GOOS=linux GOARCH=amd64 go build ./...` exits 0. The fallback is genuinely
fail-closed: `protocolPublish` checks `hasTTYSupported` first and refuses with exit 2 and a clear
"unavailable on this platform" message before even calling `platformHasTTY`. The `unix` import was
moved out of `protocol.go` into the platform-specific files, so `protocol.go` no longer pulls in
unix on Windows. This is the right fix and it works.

**[CRITICAL] Publish followed symlinks / accepted existing release dir (all 3).** FIXED.
`core.go:135` now `os.Lstat(dir)` — checks the DIRECTORY, not the file, and uses Lstat (does not
follow symlinks). If the directory exists, publish refuses. Files are created with
`os.O_WRONLY|os.O_CREATE|os.O_EXCL|noFollow` (core.go:148). On Unix, `noFollow = syscall.O_NOFOLLOW`
(nofollow_unix.go); on Windows, `noFollow = 0` (nofollow_windows.go) but the Lstat directory check
plus O_EXCL still prevent overwrite. Verified adversarially: a pre-existing empty `2.0.0/`
directory is refused; a symlink planted as `COOPERATION.md` inside a release dir does not modify
the victim file (O_NOFOLLOW refuses to open through it). The test
`TestPublishRefusesExistingReleaseDirAndSymlinks` covers both cases at the Store level.

**[CRITICAL] TTY gate bypassable by a pty (codex-1, kimi-1).** HONESTLY STATED, NOT FIXED. The
refusal text now says plainly: "this stops an ordinary agent run (whose stdin is a pipe or
/dev/null). It does NOT stop an agent that allocates a pty; that case is covered only by the
sandbox follow-up (DF-1), which is not shipped." IMPLEMENTATION.md:71-75 says the same. This is the
honest posture VC-1 resolved. Accepted as the agreed outcome.

**[CRITICAL/MAJOR] Heading-based removal report (all 3).** FIXED. `removedSections` replaced by
`droppedContent` (render.go:187-237), which compares every substantive line of the deck against
the rendered body, grouped under the heading it sat beneath. A deck paragraph under a heading the
core also has is now reported. CRLF is normalized on input (`strings.ReplaceAll(priorDeckBody,
"\r\n", "\n")`) and restored on output. Verified: `TestProtocolRenderReportsContentLostUnderASharedHeading`
asserts content under `## 3. Phases` (a shared heading) is reported; `TestProtocolRenderHandlesCRLFDecks`
asserts no false removal for `## 2. Active agents` and no mixed endings. I probed edge cases
(content moved between sections, duplicate lines): content that moves sections is not reported but
IS carried (acceptable — G1 is about loss, not relocation); duplicate lines reduced from 2 to 1 are
not reported (a limitation, but unrealistic in protocol text). No false positive found.

**[MAJOR] Path traversal via committed deck lock (kimi-1).** FIXED. `ValidVersion` (core.go:48-68)
is called on both `Load` (core.go:83) and `Publish` (core.go:127). The validator rejects `""`,
`"."`, `".."`, any string containing `/`, `\`, or `..`, anything starting with `.`, and any
character outside `[a-zA-Z0-9._-+]`. Verified adversarially through the built binary: `../../../etc`,
`..`, `.`, `.hidden`, `a/b`, `1.0.0\x00evil`, `..a`, `a..`, `a..b` all rejected with exit 1.
`TestProtocolRejectsPathTraversalInTheLock` covers the Load path; `TestPublishRejectsUnsafeVersions`
covers the Publish path.

**[MAJOR] Tests bypassed production dispatch (codex-1).** FIXED. `TestProtocolIsReachableThroughProductionDispatch`
drives `Run([]string{"protocol", "status", ...})` and `Run([]string{"protocol", "render", ...})`,
asserting the protocol code is reached and applies the core. This test would fail if the `protocol`
case were removed from `app.Run`'s dispatch. Real test, not tautological.

**[MAJOR] IMPLEMENTATION.md claimed rank-2 continuation (codex-1).** PARTIALLY CORRECTED. The
IMPLEMENTATION.md now says "Continuation is rank 2 and is not implemented here" (line 26-27) and
"nothing here claims a guarantee it does not implement (G7b)" (line 16-17). The old text "Nine
tests drive the real command entry points" is corrected to "Fifteen tests, including one that
drives production dispatch (Run)." However — codex-1's finding was explicit: "The overclaim is in
the PROTOCOL TEXT shipped by this commit" (codex-1.md:233), not in the IMPLEMENTATION.md. The §7
clause in all three COOPERATION.md copies still says "An idea that is already open completes under
the protocol version it was pinned to" as present-tense fact. The git diff of COOPERATION.md
between 4396529 and 8888e00 is 0 lines — the protocol text was not touched. See New findings.

**[MINOR] --help omitted publish (hermes-1).** FIXED. `app.go:126` now includes
`%s protocol publish --version V --from FILE (attended; requires a terminal)`. Verified in the
built binary's `--help` output.

**[MINOR] protocol status swallowed read errors (hermes-1).** FIXED IN CODE, NOT TESTED.
`protocol.go:134-143` now captures `verErr` and `pinErr` and returns exit 1 with a diagnostic
message. No test exercises this path — an unreadable store or lock file is not simulated. The fix
is correct but has no regression protection.

**[MINOR] Third protocol copy not updated (codex-1, kimi-1, hermes-1).** FIXED. The sibling repo
`parley-deck-skill` committed the §7 clause at `455aafe`. Verified:
`parley-deck-skill/skills/parley-deck/references/COOPERATION.md` now contains the blast-radius
clause at line 745, including the continuation sentence.

## New findings (by severity)

### [MAJOR] G7b: the §7 continuation guarantee is still shipped as present-tense fact with no implementation and no test

All three round-01 reviewers flagged this. codex-1 was explicit: "The overclaim is in the protocol
text shipped by this commit" (codex-1.md:233). The fix-up cycle's response was "IMPLEMENTATION.md
claimed rank-2 continuation. Corrected" (IMPLEMENTATION.md:110). But the IMPLEMENTATION.md already
correctly said rank 2 was not shipped — the old text at line 15-17 said "Ranks 2-4 ... remain
ratified and scheduled; nothing here claims a guarantee it does not implement." The thing that was
wrong was the §7 clause itself.

It is still wrong. All three COOPERATION.md copies (parley-deck, internal/protocol/defaults, and
parley-deck-skill@455aafe) contain at the end of the blast-radius clause:

    An idea that is already open completes under the protocol version it was pinned to; the next
    idea in that deck picks up the current one.

This is present-tense, not hedged. Per-idea pinning is rank 2 (G5/G7), which is not implemented:
zero code for `protocol-snapshots`, `EffectiveHash`, or any continuation resolver. FINAL.md:53-54
says "Nothing in the shipped slice may claim a guarantee it does not implement — G7b makes that
binding." FINAL.md:73 defers G5, G7, G8 — but NOT G7b. The changelog (protocol-changelog.md:23-24)
calls the TTY-gated publisher and write-once releases "the shipped guarantee" and then appends
"detection with DETECTED-UNATTRIBUTED for anything else" — also unimplemented (see below).

The git diff of both in-repo COOPERATION.md copies between 4396529 and 8888e00 is zero lines. The
fix-up did not touch the protocol text. The continuation sentence should either be hedged to
future tense ("An idea that is already open will complete under the protocol version it will be
pinned to, once rank 2 ships") or removed until rank 2 lands with G5/G7/G8.

### [MAJOR] G7b: DETECTED-UNATTRIBUTED still cited as a shipped guarantee with no implementation and no test

`parley-deck/meta/protocol-changelog.md:24` says: "the shipped guarantee is: write-once releases,
an attended TTY-gated publisher, no agent-accessible write path, and detection with
`DETECTED-UNATTRIBUTED` for anything else." The string `DETECTED-UNATTRIBUTED` appears in zero Go
files (verified: `search_files` across `internal/` returns 0 matches). It is not implemented, not
tested, and is presented as part of "the shipped guarantee" in the changelog that G6 passes. The
changelog was not modified in this fix-up cycle (0-line diff). This was flagged in round-01
(hermes-1.md:175-176) and is not listed in the fix-up's "Fixes applied" section — it was missed.

### [MAJOR] G7b: the TTY-gate refusal still has no end-to-end test

All three round-01 reviewers called this the most testable G7b guarantee. The fix-up cycle added
six new tests, but none of them drive `protocolPublish` through `runProtocol` or `Run`. The publish
tests (`TestPublishRefusesExistingReleaseDirAndSymlinks`, `TestPublishRejectsUnsafeVersions`) call
`store.Publish` directly — they test the Store layer, not the CLI's TTY gate. The test environment
has no TTY, so `runProtocol([]string{"publish", "--version", "v", "--from", f})` expecting exit 2
is trivially writable and would directly assert the gate. The refusal text was updated to honestly
state the pty limitation, which is the agreed outcome — but the behavior that IS shipped (refusing
without a TTY) is still documented in the usage string (protocol.go:24), the §7 clause, and
IMPLEMENTATION.md:56-62 with no test. This is a G7b violation on the exact surface G7b was written
for.

### [MINOR] protocol status read-error fix has no test

The fix (protocol.go:134-143) correctly reports unreadable store/lock errors instead of swallowing
them. No test simulates an unreadable store or lock file. The fix is correct but untested — if it
regresses, no test fails.

## Test-quality assessment

15 tests total (9 original + 6 new). All pass. `go test ./...` exits 0. All 15 tests drive behavior,
not implementation trivia.

**New tests — would each fail if its fix were reverted?**

1. `TestProtocolRenderReportsContentLostUnderASharedHeading` — REAL. If `droppedContent` reverts
   to `removedSections` (heading-only), the `## 3. Phases` heading exists in both deck and core, so
   nothing is reported, and the `strings.Contains(out.String(), "## 3. Phases")` assertion fails.

2. `TestProtocolRenderHandlesCRLFDecks` — REAL. If CRLF normalization reverts, the `\r` on each
   heading line prevents matching the core map, producing a false removal report containing
   `## 2. Active agents` → assertion fails. The mixed-ending check (line 340: strip `\r\n`, then
   check for `\n`) also catches a failure to restore CRLF — if restoration is removed, the output
   is pure LF, the `\r\n` strip is a no-op, and the remaining `\n` fails the assertion.

3. `TestProtocolRejectsPathTraversalInTheLock` — REAL. If `ValidVersion` is removed from `Load`,
   `render --yes` succeeds with exit 0 for `../../../etc`, and the `code == 0` assertion fails.

4. `TestPublishRefusesExistingReleaseDirAndSymlinks` — REAL. If `Lstat(dir)` reverts to
   `Stat(path)`, the existing empty `2.0.0/` directory passes (no COOPERATION.md file), publish
   succeeds, and the "published into an existing release directory" assertion fails. If
   `O_NOFOLLOW` is removed, the symlink write-through modifies the victim file, and the
   `string(got) != "original"` assertion fails.

5. `TestPublishRejectsUnsafeVersions` — REAL. If `ValidVersion` is removed from `Publish`, unsafe
   versions are accepted, and the `err == nil` assertion fails.

6. `TestProtocolIsReachableThroughProductionDispatch` — REAL. If the `protocol` case is removed
   from `app.Run`, `Run([]string{"protocol", ...})` returns non-zero, and the exit-code assertion
   fails. This closes codex-1's "tests bypass production dispatch" finding.

**No tautologies found.** Every new test asserts behavior that would break if the corresponding
fix were reverted. The revert check is concrete for each.

**Remaining test gaps (G7b):**
- TTY-gate refusal: no test drives `protocolPublish` through any entry point. The test env has no
  TTY, making this trivially writable.
- DETECTED-UNATTRIBUTED: not implemented, so not testable — but it is cited as shipped.
- §7 continuation guarantee: not implemented, so not testable — but it is shipped as present tense.
- `protocol status` read-error path: not tested.

**droppedContent edge cases probed:**
- Content that moves between sections (appears in deck's §A, in render's §B): NOT reported, but the
  content IS carried — acceptable, G1 is about loss not relocation.
- Duplicate lines in deck (same line twice, once in render): NOT reported. A line count drops from
  2 to 1 silently. Unlikely in protocol text; not a G1 violation in practice.
- No false positive found: no input caused `droppedContent` to report content that IS carried
  forward in the render.
