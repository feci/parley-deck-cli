---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
review-round: 3
date: 2026-08-07
reviewed-commit: 4a5c447
verdict: FINDINGS
---
## Summary

Cycle 2 fixed the headline G7b violation: all three COOPERATION.md copies now carry an explicit
"Not yet in force" paragraph that hedges per-idea pinning and DETECTED-UNATTRIBUTED as ratified but
not implemented. The TTY-gate refusal and the status read-error path both gained end-to-end tests
that fail when their fixes are reverted. ValidVersion now rejects untrimmed input. CRLF cores are
normalized. The traversal test now asserts the reason and a new TestLoadRefusesToEscapeTheStore
plants a real body outside the store. All three GOOS targets compile, vet is clean, and the full
test suite passes.

But the cycle repeated the exact pattern the brief warned about: each fix introduced or left
something new. The claimed "per-section" matching in droppedContent is not implemented — the match
is still global, so identical text in a differently-governed section is silently erased (the same
CRITICAL codex-1 found in round-02, with only the multiplicity half fixed). assertNoSymlinkComponents
was added to Publish but not Load, so a symlinked store root still lets render consume bytes outside
the store. The changelog still cites DETECTED-UNATTRIBUTED as a shipped guarantee. The protocol check
preamble still contradicts its entries. And the test count in IMPLEMENTATION.md is wrong (18, not 19).

## Round-02 findings: fixed or not

**[MAJOR hermes-1 + kimi-1] §7 text shipped unimplemented guarantees.** FIXED in the protocol text.
All three COOPERATION.md copies (parley-deck, internal/protocol/defaults, parley-deck-skill@455aafe)
now carry:

    **Not yet in force — do not rely on it.** Per-idea version pinning ... and the
    `DETECTED-UNATTRIBUTED` tamper signal are **ratified but not implemented**.

The continuation sentence "An idea that is already open completes under the protocol version it was
pinned to" is gone, replaced by the hedged paragraph. The TTY-gate text now says "it does not stop
an agent that allocates a pty." The "What IS in force today" paragraph names write-once, attended-
only, and no agent-accessible write path — all of which have tests or are honestly hedged. This is
the right fix and it landed in the right place. Verified by reading all three copies and confirming
they are byte-identical in the §7 block.

**[CRITICAL codex-1] droppedContent lost section context and multiplicity.** PARTIALLY FIXED —
multiplicity yes, section context no. The implementation at render.go:201-250 replaced the boolean
`kept` map with an integer `renderCounts` map (consume-a-copy), which fixes the multiplicity
collapse: if the deck has 3 copies of a line and the render has 1, 2 are now reported as dropped.
But the match is still GLOBAL: `renderCounts` is built from every rendered line regardless of
heading, and a deck line is "carried" if it matches ANY rendered line. The comment at render.go:199
says "per-section and multiplicity-aware: within the section a line belongs to, a deck line counts
as carried only if the render still has an unconsumed copy of it" — but the code does not implement
per-section matching.

I verified this with a direct probe: a deck with "- Deploy production." under "## Requires explicit
user approval" and a render with the same line under "## Allowed automatically" reports only "##
Requires explicit user approval — 1 line not carried forward" (the heading itself), not the content
line. The content line is consumed by the copy in the wrong section. This is the exact semantic-
erasure class codex-1 demonstrated in round-02, and it is still present. The existing test
TestProtocolRenderReportsContentLostUnderASharedHeading does not catch this because its added content
("PROJECT RULE: deployments need two approvals.") does not appear anywhere else in the render.

**[MAJOR codex-1 + kimi-1] symlink hardening covered only the final create.** PARTIALLY FIXED —
Publish yes, Load no. assertNoSymlinkComponents (core.go:202-221) checks the store's own three
components (core, protocol, home) and refuses if any is a symlink. I verified it catches a symlink
at home/protocol and at the store root itself. It is correctly scoped to the store's own components,
not the filesystem root (/var symlink on macOS does not trigger it). It is called in Publish
(core.go:138). But it is NOT called in Load. I verified: with the store root as a symlink to an
outside directory containing 1.0.0/COOPERATION.md, Load("1.0.0") reads "SYMLINK-ESCAPED-CONTENT"
from outside the store without error. The round-02 finding explicitly named the Load path
("Store.Load then calls plain os.ReadFile, which follows the release symlink") and the fix only
applied to Publish.

**[MINOR codex-1] ValidVersion validated a trimmed copy, used the untrimmed string.** FIXED.
core.go:51 now checks `if v != strings.TrimSpace(v) { return false }` before any other validation.
Verified: " 1.0.0", "1.0.0 ", "\t1.0.0" are all rejected; "1.0.0", "2.0.0-beta" still pass.

**[MINOR kimi-1] TestProtocolRejectsPathTraversalInTheLock was vacuous.** FIXED. The test now
asserts `strings.Contains(errb.String(), "unsafe version")` for each bad version. I reverted the
ValidVersion call in Load and the test failed with "rejected for the wrong reason: release not
installed" — confirming it now pins the reason, not just the exit code. The new
TestLoadRefusesToEscapeTheStore plants a real COOPERATION.md outside the store and asserts Load
refuses with "unsafe version". Same revert made it fail with "Load escaped the store and read
'PLANTED'". Both tests are real.

**[MAJOR hermes-1] TTY refusal and status read-errors had no e2e tests.** FIXED.
TestPublishRefusesWithoutATerminal drives runProtocol with publish and asserts exit != 0, no
release written, message contains "attended", and message does not overclaim (contains "pty" or
"unavailable on this platform"). I reverted the TTY gate (bypassed both hasTTYSupported and
platformHasTTY checks) and the test failed with "published without a controlling terminal".
TestProtocolStatusReportsReadErrors makes the lock a directory (unreadable as a file) and asserts
status exits non-zero. I reverted the error-handling fix and the test failed with "status reported
success on an unreadable lock". Both tests are real.

**[MINOR kimi-1] CRLF cores unnormalized; synced stamp reported as lost content.** FIXED.
render.go:56 now normalizes rel.Body with `strings.ReplaceAll(rel.Body, "\r\n", "\n")`. I verified
with a direct probe: CRLF core + LF deck produces no lone CR chars and is idempotent across two
renders; CRLF core + CRLF deck produces no CRCRLF sequences. The synced stamp is now skipped in
droppedContent (render.go:228-230) with `if strings.HasPrefix(t, syncedPrefix) { continue }`.

## New findings (by severity)

### [MAJOR] droppedContent still uses a global match, not per-section — the claimed fix is not implemented

The IMPLEMENTATION.md cycle-2 entry says: "Now per-section and multiplicity-aware (consume-a-copy),
so semantic erasure cannot hide behind a coincidental match." The comment at render.go:199 says:
"per-section and multiplicity-aware: within the section a line belongs to, a deck line counts as
carried only if the render still has an unconsumed copy of it."

The code does not implement per-section matching. `renderCounts` at render.go:205 is a single global
map built from every line in the rendered body. A deck line at render.go:231 is carried if
`renderCounts[t] > 0` — regardless of which section the rendered copy sits in. The multiplicity
half (int counts, consume-a-copy) is real; the section-context half is not.

Verified with a direct probe (test written, run, removed):
- Deck: "- Deploy production." under "## Requires explicit user approval"
- Render: "- Deploy production." under "## Allowed automatically"
- droppedContent reports: "## Requires explicit user approval — 1 line not carried forward"
  (the heading only; the content line was consumed by the copy in the wrong section)

This is the exact semantic-erasure class codex-1's round-02 CRITICAL demonstrated: a restriction
under one heading is silently inverted when the same text survives under a differently-governed
heading. The fix claimed to close it; it did not. No test covers this case — the existing
shared-heading test uses unique content that does not appear elsewhere in the render.

Suggested fix: build a per-section renderCounts map (keyed by heading path), and match deck lines
only against counts in their own section. Or, at minimum, add a test that plants identical text in
two differently-governed sections and asserts the dropped content is reported.

### [MAJOR] assertNoSymlinkComponents is not called in Load — a symlinked store root still escapes

core.go:138 calls assertNoSymlinkComponents in Publish. core.go:86-99 (Load) does not. The round-02
finding by codex-1 explicitly named both paths: "Store.Load then calls plain os.ReadFile, which
follows the release symlink." The fix applied only to Publish.

Verified with a direct probe (test written, run, removed):
- Store root (home/protocol/core) is a symlink to an outside directory
- Outside directory contains 1.0.0/COOPERATION.md with "SYMLINK-ESCAPED-CONTENT"
- store.Load("1.0.0") returns Release{Body: "SYMLINK-ESCAPED-CONTENT"} with no error

A committed lock with a valid version can still make render consume bytes outside the release store.
The local-access prerequisite is the same as the original finding, hence MAJOR not CRITICAL.

Suggested fix: call assertNoSymlinkComponents(s.Root) at the top of Load, before joining the version
onto the path.

### [MAJOR] The changelog still cites DETECTED-UNATTRIBUTED as a shipped guarantee

parley-deck/meta/protocol-changelog.md:23-24 says: "the shipped guarantee is: write-once releases,
an attended TTY-gated publisher, no agent-accessible write path, and detection with
`DETECTED-UNATTRIBUTED` for anything else."

DETECTED-UNATTRIBUTED appears in zero Go files (verified: grep across internal/ returns 0 matches).
The §7 protocol text now correctly hedges it as "ratified but not implemented." The changelog was
not touched in cycle 2 (zero-line diff). This was flagged in round-01 (hermes-1.md:175) and round-02
(hermes-1.md:129-137) and is not listed in the cycle-2 fixes. The changelog is the document G6
passes, and it still presents an unimplemented signal as part of "the shipped guarantee."

Suggested fix: update the changelog's enforcement paragraph to match the §7 hedging: name
DETECTED-UNATTRIBUTED as ratified but not implemented, and list only what IS shipped.

### [MINOR] protocol check's preamble still contradicts its content-based entries

protocol.go:268 prints "sections present here but not in the core:" before listing res.Removed
entries. But res.Removed entries are now content-based (e.g., "## 3. Phases — 2 lines not carried
forward"), where the section IS in the core. The render preamble was fixed in this cycle
(protocol.go:194: "deck content NOT carried forward by core %s:") but the check preamble was not.

Verified with a direct probe: after hand-editing a deck to add content under "## 3. Phases" (a
heading the core also has), `protocol check` printed:

    sections present here but not in the core:
      - ## 3. Phases — 2 lines not carried forward

The preamble asserts a falsehood about the section it names. This is the same defect kimi-1 flagged
in round-02 (MINOR) for the render preamble; the render half was fixed, the check half was not.

Suggested fix: change protocol.go:268 to match the render preamble: "deck content NOT carried
forward by core %s:" or similar.

### [MINOR] protocol render silently swallows deck-file read errors (kimi-1 round-02, still open)

protocol.go:180: `prior, _ := os.ReadFile(path)` discards the read error. If the deck file is
unreadable, render proceeds with an empty prior, overwrites the file, and loses content without a
removal report or error. Verified: with the deck file chmod 000, render exited 0, printed "Wrote
...", and overwrote the file. This was kimi-1's round-02 open finding (not claimed as fixed in
cycle 2) and is still present. It is another route around G1's report.

Suggested fix: capture the error and exit non-zero if the deck file exists but is unreadable
(distinguish from "does not exist" which is the fresh-deck path).

### [MINOR] TestPublishRefusesExistingReleaseDirAndSymlinks symlink half is still tautological

kimi-1's round-02 finding: the test pre-creates the 3.0.0/ directory (line 412), so
`Publish("3.0.0")` is refused at the `Lstat(dir)` check before `OpenFile` with `O_NOFOLLOW` is
reached. The symlink at line 415 is inside the already-existing directory. I confirmed: the publish
error is "release 3.0.0 already exists and releases are write-once" — the directory check fires
first. Removing `|noFollow` from the OpenFile flags would leave the test green. This was not
addressed in cycle 2.

Suggested fix: test O_NOFOLLOW at a lower level (create a symlink as the release path itself, not
inside a pre-existing directory), or drop the dead symlink half from the test's name/setup.

### [NIT] IMPLEMENTATION.md says "nineteen protocol tests" — there are eighteen

IMPLEMENTATION.md:159: "Nineteen protocol tests, one of which drives production dispatch."
grep -c '^func Test' internal/app/protocol_test.go returns 18. The count was "fifteen" in cycle 1
(correct at the time), bumped to "nineteen" in cycle 2, but only 3 new tests were added (was 15,
now 18). TestLoadRefusesToEscapeTheStore is the 18th.

## Test-quality assessment

18 protocol tests total. All pass. `go build ./...`, `GOOS=windows go build ./...`,
`GOOS=linux go build ./...`, `go vet ./...`, `go test ./...` all exit 0.

### New/changed tests — would each fail if its fix were reverted?

1. **TestProtocolRejectsPathTraversalInTheLock** — REAL. Reverted ValidVersion in Load: failed
   with "rejected for the wrong reason: release not installed." The reason assertion pins the fix.

2. **TestLoadRefusesToEscapeTheStore** — REAL. Same revert: failed with "Load escaped the store
   and read 'PLANTED'." The planted outside body makes the test fail for the right reason.

3. **TestPublishRefusesWithoutATerminal** — REAL. Reverted the TTY gate (bypassed both checks):
   failed with "published without a controlling terminal." Also asserts no release was written and
   the message does not overclaim.

4. **TestProtocolStatusReportsReadErrors** — REAL. Reverted the error-handling fix (restored `_ =
   ` swallowing): failed with "status reported success on an unreadable lock."

5. **TestPublishRefusesExistingReleaseDirAndSymlinks** — PARTIALLY REAL (unchanged from round-02).
   The directory half catches a revert to Stat-on-file. The symlink half is still tautological:
   the pre-created directory makes Lstat(dir) fire before OpenFile reaches the symlink. Removing
   `|noFollow` leaves the test green. kimi-1 flagged this in round-02; not addressed.

6. **TestPublishRejectsUnsafeVersions** — REAL (unchanged). Reverting Publish's ValidVersion to
   the weaker round-01 check fails it on "." and ".hidden". Does not test padded versions (the
   new untrimmed rejection) — a gap but not a tautology.

No new tautologies introduced in cycle 2. The four new tests are all genuine. The one residual
tautology is carried from round-02 (TestPublishRefusesExistingReleaseDirAndSymlinks symlink half).

### Does any text still claim a guarantee with no end-to-end test?

1. **Changelog: DETECTED-UNATTRIBUTED as a shipped guarantee** — yes, still claimed, still no code,
   still no test. (MAJOR above)
2. **droppedContent: "per-section" matching** — claimed in code comment and IMPLEMENTATION.md, not
   implemented, and no test covers the section-context case. (MAJOR above)
3. **assertNoSymlinkComponents on Load** — the round-02 finding named Load; the fix only applied to
   Publish. No test checks Load against a symlinked store. (MAJOR above)
4. **protocol check preamble** — says "sections present here but not in the core" for content-based
   entries that are under sections IS in the core. (MINOR above)
5. **protocol render read-error swallowing** — no test, and the behavior loses content silently.
   (MINOR above)
6. **"no agent-accessible write path"** — the §7 "What IS in force" paragraph says this. The TTY
   gate is tested. A pty bypass is honestly hedged in the same paragraph. This claim is accurate
   at the code level (protocolPublish is the only write function, and it requires a TTY). Not a
   finding.
7. **write-once releases** — tested by TestCoreReleasesAreWriteOnce. Not a finding.
8. **attended-only publish** — tested by TestPublishRefusesWithoutATerminal. Not a finding.

### What is genuinely fixed and should be kept

- The §7 protocol text in all three copies is now honestly hedged. This was the core G7b violation
  and it is closed in the right place.
- The TTY-gate and status read-error e2e tests are real and revert-sensitive.
- ValidVersion's untrimmed-input rejection is real and tested.
- CRLF core normalization is real and verified.
- The synced-stamp false-positive is fixed (skipped in droppedContent).
- The traversal test now asserts the reason and the new Load escape test plants a real body.
- assertNoSymlinkComponents is correctly scoped (store's own components, not filesystem root) and
  works for Publish.
- The failed-publish cleanup (os.Remove on error) is defensive and honestly not claimed as tested.
