---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
review-round: 4
date: 2026-08-07
reviewed-commit: 02bb1e3
verdict: FINDINGS
---
## Summary

The per-section matching code in `droppedContent` is finally implemented correctly — verified
with a direct probe that constructs a real cross-section case (a line under the deck's §3 that
survives only in the render's §11). The old global matcher reports "1 line" (masks the deletion);
the new per-section matcher reports "2 lines" (catches it). The symlink refusal on Load works
for both a symlinked store component and a symlinked release directory — verified behaviorally.
The changelog no longer cites `DETECTED-UNATTRIBUTED` as shipped. The render read-error fix and
the check preamble fix both work. All three GOOS targets compile, vet is clean, 20 protocol tests
pass.

But the pattern held: the fix's test outran the fix it claims to pin. `TestProtocolRenderDetectsCrossSectionErasure`
— the test written to prove the per-section fix — is tautological. It passes with the per-section
implementation reverted to the old global multiset. Both implementations produce identical output
for the test's fixture because the "erasure to catch" (the STALE line) is absent from the entire
render, not just from its section. No cross-section coincidence exists in the test. The fix is
real; the test that claims to guard it is not.

Three other cycle-3 claimed fixes have no end-to-end test: the render read-error refusal, the
check preamble rewording, and the release-directory Lstat in Load. All three were revert-verified
to leave the suite green. The §7 "no agent-accessible code path writes a release" overclaim that
kimi-1 flagged in round-03 is still present — the cycle narrowed the TTY refusal description but
left this absolute clause unqualified.

## Prior findings: fixed or not

**[MAJOR hermes-1 + kimi-1 round-03] Cycle 2's "per-section" droppedContent was global.** FIXED
IN CODE, NOT PINNED BY TEST. The implementation at render.go:207 (`rendered := indexBySection(renderedBody)`)
and render.go:225 (`if counts, ok := rendered[current]; ok && counts[t] > 0`) now genuinely
indexes the render per section and matches within the same section. I verified with a direct
probe: a deck with "SHARED-LINE-XYZ" under "## 3. Phases" and a render with the same line only
under "## 11. Appendix" — per-section reports "## 3. Phases — 2 lines not carried forward"
(catches it); the old global multiset reports "1 line" (masks it). The fix works.

But `TestProtocolRenderDetectsCrossSectionErasure` (protocol_test.go:505-521) does not test this.
It adds "Core rules live here." to the deck's §3 — but "Core rules live here." is also in the
render's §3 (it comes from `testCore`). The STALE line is absent from the ENTIRE render, not
just from its section. No cross-section coincidence exists. Both per-section and global produce
"## 3. Phases — 1 line not carried forward" (the STALE line). I reverted `droppedContent` to
the old global multiset and the test passed. The test's own comment admits this: "Core rules
live here. exists in the RENDER, but under ## 3. Phases in the core — while in the deck it sat
under ## 3. Phases too." The "erasure to catch" is a plain deletion, not a cross-section case.

**[MAJOR hermes-1 round-03] assertNoSymlinkComponents was write-only.** FIXED IN CODE; PARTIALLY
PINNED. `Load` (core.go:93) now calls `assertNoSymlinkComponents(s.Root)` and additionally Lstats
the release directory (core.go:97-99). I verified both paths: a symlinked store component is
refused with "is a symlink; the core store must not be reached through one"; a symlinked release
directory is refused with "release directory ... is a symlink". `TestLoadRefusesASymlinkedStore`
pins the store-component refusal — it FAILS when the Load symlink check is reverted. But the
release-directory Lstat (core.go:97-99) is unpinned: reverting just that check leaves all 20
tests green. No test plants a symlink as the release directory itself.

**[MAJOR hermes-1 + kimi-1 round-03] Changelog cited DETECTED-UNATTRIBUTED as shipped.** FIXED.
`parley-deck/meta/protocol-changelog.md:25-27` now says "**`DETECTED-UNATTRIBUTED` and per-idea
pinning are ratified but NOT implemented** — they are ranks 2 and 4." Zero Go references to
DETECTED-UNATTRIBUTED (verified). The changelog no longer presents it as part of "the shipped
guarantee."

**[MINOR kimi-1 round-03] §7 "in force" sentence overclaimed.** PARTIALLY FIXED. The cycle
narrowed the TTY refusal description: "refuses when it cannot see a controlling terminal (which
stops an ordinary agent run, not one that allocates a pty)". But the same sentence still ends
with "and no agent-accessible code path writes a release" — an unqualified absolute. For a
pty-allocating agent, `parley protocol publish` IS an agent-accessible code path that writes a
release, and the sentence itself says the TTY gate does not stop such an agent. The clause
kimi-1 suggested narrowing ("no code path an ordinary agent run can reach writes a release")
was not applied. The overclaim persists in all three copies (deck, embedded, skill). See New
findings.

**[MINOR kimi-1 round-02, claimed fixed in cycle 3] protocol render swallowed deck-file read
errors.** FIXED IN CODE, NOT TESTED. protocol.go:180-186 now captures the read error and exits
non-zero if the deck file exists but is unreadable (distinguishing `os.ErrNotExist`). I verified
behaviorally: a chmod 000 deck file produces exit=1 with "protocol render: ... permission
denied". But no test exercises this path. Reverting the fix (restoring `prior, _ := os.ReadFile`)
leaves all 20 protocol tests green. See New findings.

**[MINOR hermes-1 round-03] protocol check preamble contradicted entries.** FIXED IN CODE, NOT
TESTED. protocol.go:274 now prints "deck content not carried forward by this core:" instead of
"sections present here but not in the core:". I verified behaviorally. But no test asserts the
preamble text. Reverting the fix leaves all 20 tests green. See New findings.

**[NIT kimi-1 round-03] Test count did not reproduce.** FIXED. `grep -c '^func Test' protocol_test.go`
returns 20. IMPLEMENTATION.md:196 says "20 protocol tests." The count is now correct.

## New findings (by severity)

### [MAJOR] TestProtocolRenderDetectsCrossSectionErasure is tautological — the per-section fix has no regression protection

The test written to pin the per-section fix does not test per-section matching. It passes with
the fix reverted to the old global multiset. Both implementations produce identical output for
the test's fixture because the "erasure to catch" (the STALE line) is absent from the entire
render — a plain deletion, not a cross-section coincidence.

Verified by reverting `droppedContent` to the old global multiset (replacing `indexBySection`
with a flat `renderCounts` map) and running the test: PASS. Both outputs contain "## 3. Phases"
which is all the assertion checks.

This is the fourth cycle in a row where a fix's documentation outran its verification. This time
the CODE is correct (verified by direct probe) but the TEST is not. The per-section fix — the
headline fix of this cycle, the one both reviewers caught in round-03 — has no test that would
fail if it were reverted. Any future regression to a global matcher would go undetected.

Suggested fix: replace the test's fixture with a real cross-section case — a deck line under
section A that is ABSENT from render's section A but PRESENT in render's section B — and assert
the line count reflects the per-section match (2, not 1). My probe (TestHermesCrossSectionProbe)
constructs exactly this case and distinguishes the two implementations.

### [MINOR] Three cycle-3 claimed fixes have no end-to-end test

Each of these was claimed as fixed in IMPLEMENTATION.md cycle 3, each works behaviorally, and
reverting each leaves all 20 protocol tests green:

1. **Render read-error refusal** (protocol.go:180-186). Reverting to `prior, _ := os.ReadFile`
   leaves the suite green. The fix works (verified: chmod 000 → exit 1, "permission denied"),
   but no test simulates an unreadable deck file through `runProtocol`.

2. **Check preamble rewording** (protocol.go:274). Reverting to "sections present here but not
   in the core:" leaves the suite green. The fix works (verified), but no test asserts the
   preamble text.

3. **Release-directory Lstat in Load** (core.go:97-99). Reverting just this check (keeping the
   store-component `assertNoSymlinkComponents`) leaves the suite green. The fix works (verified:
   a symlinked release directory is refused with "release directory ... is a symlink"), but
   `TestLoadRefusesASymlinkedStore` only tests a symlinked store component, not a symlinked
   release directory.

All three are G7b violations: claimed fixes with no regression protection. The code is correct;
the tests are missing.

### [MINOR] §7 "no agent-accessible code path writes a release" still overclaims for pty-agents (carried from kimi-1 round-03)

parley-deck/COOPERATION.md:773, internal/protocol/defaults/COOPERATION.md:762, and the skill
copy all carry: "and no agent-accessible code path writes a release." The same sentence says
the TTY refusal "stops an ordinary agent run, not one that allocates a pty." For a pty-allocating
agent, `parley protocol publish` IS an agent-accessible code path that writes a release — the
sentence itself says the gate does not stop such an agent. The absolute clause is false for
pty-agents.

The cycle-3 fix narrowed the TTY refusal description (added the pty qualifier to the parenthetical)
and added "refused through a symlinked store" but did not qualify the "no agent-accessible code
path" clause. kimi-1's suggested fix ("no code path an ordinary agent run can reach writes a
release") was not applied. The changelog carries the same overclaim: "no agent-accessible write
path" (protocol-changelog.md:25).

### [NIT] Carried findings from prior rounds, not claimed as fixed in cycle 3

These are restated so they are not lost; none is a new defect:

- **protocol.go:301 "hash detection"** (kimi-1 round-03 NIT). The refusal text says "The durable
  guarantees today are write-once releases and hash detection." "Hash detection" is ambiguous
  between deck-view drift detection (implemented, tested via `check`) and release-tamper detection
  (DETECTED-UNATTRIBUTED, not implemented). Borderline; one word from honest.
- **TestPublishRefusesExistingReleaseDirAndSymlinks O_NOFOLLOW half is tautological** (kimi-1
  round-02 MINOR). The test pre-creates the 3.0.0/ directory, so `Lstat(dir)` fires before
  `OpenFile` with `O_NOFOLLOW` is reached. The symlink is inside the already-existing directory.
  Removing `|noFollow` from the flags leaves the test green.
- **TestPublishRejectsUnsafeVersions has no padded case** (kimi-1 round-03 NIT). The test list
  is `"", ".", "..", "../x", "a/b", ".hidden"` — no `" 1.0.0"` or `"1.0.0 "`. ValidVersion
  rejects padded input (verified at the API level), but no test pins it. The comment at
  protocol_test.go:354-356 says "see TestPublishRejectsUnsafeVersions" for untrimmed rejection —
  that test does not cover it.
- **IMPLEMENTATION.md:69 "GOOS=windows and GOOS=linux builds are now part of the check"**
  (kimi-1 round-02 open). `grep -rn GOOS scripts/ .github/` returns 0 matches. No automated
  check enforces cross-builds; the claim means the implementer ran them manually.
- **droppedContent false positive on relocated content** (observed this round, inherent to the
  design). A line that moves from deck §A to render §B is reported as "not carried forward"
  under §A even though it IS in the render (just under a different heading). Per-section
  matching trades this false positive for the cross-section masking false negative it fixes.
  Arguably acceptable — G1 is about loss, and the report is conservative — but worth noting
  since the global matcher did not produce this false positive.

## Test-quality assessment

20 protocol tests, all pass. `go build ./...`, `GOOS=windows go build ./...`,
`GOOS=linux go build ./...`, `go vet ./...`, `go test ./...` all exit 0.

Command output (this machine, 02bb1e3, clean tree):

```
$ go build ./...                  EXIT=0
$ GOOS=windows go build ./...     EXIT=0
$ GOOS=linux go build ./...       EXIT=0
$ go vet ./...                    EXIT=0
$ go test ./...                   all 26 packages ok, EXIT=0
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore|TestLoad' -count=1 -v
  20 tests, 20 PASS
$ grep -c '^func Test' internal/app/protocol_test.go
  20
```

### New/changed tests — would each fail if its fix were reverted?

1. **TestProtocolRenderDetectsCrossSectionErasure** (new) — TAUTOLOGICAL. Reverted
   `droppedContent` to the old global multiset: PASS. Both implementations produce "## 3. Phases
   — 1 line not carried forward" because the STALE line is absent from the entire render. No
   cross-section coincidence exists in the fixture. The test's assertion
   (`strings.Contains(out, "## 3. Phases")`) passes either way.

2. **TestLoadRefusesASymlinkedStore** (new) — REAL. Reverted the Load symlink checks
   (both `assertNoSymlinkComponents` and the release-dir Lstat): FAIL with "Load read through a
   symlinked store component." Pins the store-component refusal. Does NOT pin the release-dir
   Lstat (reverting just that check leaves the test green — the store component symlink fires
   first in this test's setup).

### Does any text still claim a guarantee with no end-to-end test?

1. **§7 "no agent-accessible code path writes a release"** — overclaim for pty-agents, all three
   copies. (MINOR above)
2. **Render read-error refusal** — claimed as fixed, works, no test. (MINOR above)
3. **Check preamble rewording** — claimed as fixed, works, no test. (MINOR above)
4. **Release-directory Lstat in Load** — claimed as fixed, works, no test. (MINOR above)
5. **Per-section matching** — works, but the test that claims to pin it is tautological. (MAJOR
   above)
6. **protocol.go:301 "hash detection"** — ambiguous, borderline. (NIT above)
7. **"GOOS builds are now part of the check"** — no CI check. (NIT above)

The following claims ARE tested and survive revert verification:
- Write-once releases → `TestCoreReleasesAreWriteOnce`
- Attended-only publish (TTY gate) → `TestPublishRefusesWithoutATerminal` (fails on revert)
- Store-component symlink refusal on Load → `TestLoadRefusesASymlinkedStore` (fails on revert)
- Path traversal refusal → `TestProtocolRejectsPathTraversalInTheLock` + `TestLoadRefusesToEscapeTheStore`
  (both fail on revert)
- Production dispatch reachability → `TestProtocolIsReachableThroughProductionDispatch`
- Status read-error reporting → `TestProtocolStatusReportsReadErrors` (fails on revert)
