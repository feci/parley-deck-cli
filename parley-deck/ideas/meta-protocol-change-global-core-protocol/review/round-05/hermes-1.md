---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
review-round: 5
date: 2026-08-07
reviewed-commit: 1903e49
verdict: CLEAN
---
## Summary

Cycle 4 is the first cycle where every claimed fix is real AND every claimed test is real. I
reverted each of the four "proven reversion-sensitive" tests in a scratch copy and confirmed all
four FAIL with their fix reverted. The headline MAJOR from round-04 — the tautological
cross-section test — is genuinely fixed: the fixture now puts "Core rules live here." under the
deck's §2 while the render has it only under §3, a real cross-section coincidence that the global
multiset masks and the per-section matcher catches. The three MINOR tests (Load release-dir
symlink, Publish store-component symlink, render read-error) are all real and all revert-sensitive.
The render read-error test correctly uses --dry-run to avoid the "passed for the wrong reason" trap
the implementer caught on the second pass.

All three GOOS targets compile, vet is clean, 23 protocol tests pass, full suite green. No new
findings. Carried NITs from prior rounds are restated so they are not lost; none is a new defect
and none blocks exit.

## Prior findings: fixed or not

**[MAJOR hermes-1 + kimi-1 round-04] The cross-section test was tautological.** FIXED AND
VERIFIED. The fixture now inserts "Core rules live here." before "## 3. Phases" in the deck,
placing it under the §2 heading; the render has that line only under §3. This is a genuine
cross-section coincidence: the per-section matcher reports it as dropped from §2 (assertion
checks for "## 2. Active agents" in output); the global matcher consumes it from §3 and reports
nothing from §2. I reverted droppedContent to a flat global multiset in a scratch copy and the
test FAILED:

```
--- FAIL: TestProtocolRenderDetectsCrossSectionErasure
    a line dropped from §2 was masked by an identical line surviving in §3:
      - ## 3. Phases — 1 line not carried forward
      - ## 99. Project-specific rule — 2 lines not carried forward
```

With the fix in place: PASS. Both directions confirmed. The headline fix of cycle 3 now has
genuine regression protection.

**[MINOR kimi-1 round-04] Three cycle-3 fixes had no test noticing their absence.** FIXED AND
VERIFIED. Three new tests added, each revert-verified:

1. TestLoadRefusesASymlinkedReleaseDirectory — removes the release-dir Lstat check in Load
   (core.go:97-99): FAIL ("Load read through a symlinked release directory and got 'PLANTED'").

2. TestPublishRefusesASymlinkedStoreComponent — removes assertNoSymlinkComponents from Publish
   (core.go:148): FAIL ("published through a symlinked store component").

3. TestProtocolRenderReportsDeckReadErrors — restores prior, _ := os.ReadFile (swallowing the
   read error): FAIL ("render treated an unreadable deck as empty"). The test uses --dry-run
   (line 600) with a comment explaining why (lines 597-598): without --dry-run the write would
   also fail and the test would pass for the wrong reason. The second pass is genuine.

**[NIT kimi-1 round-04] Core-side section rename reports a surviving line as not carried
forward.** KEPT DELIBERATELY, as recorded. Per-section strictness trades this false positive for
the cross-section masking false negative it fixes. A data-loss report should err loud. This is
inherent to the design and correctly documented in IMPLEMENTATION.md cycle 4. Not a defect.

**Carried findings from earlier rounds, not claimed as fixed in cycle 4 (restated, not new):**

- [NIT] §7 "no agent-accessible code path writes a release" (all three copies, line 775 deck /
  766 embedded). Still an unqualified absolute contradicted by the same sentence's pty
  admission. kimi-1 downgraded to NIT in round-04 because the pty limitation is stated inline.
  The changelog carries the same overclaim ("no agent-accessible write path", line 25).

- [NIT] protocol.go:301 "hash detection" in the refusal text. Still ambiguous between deck-view
  drift detection (implemented, tested via check) and release-tamper detection
  (DETECTED-UNATTRIBUTED, not implemented). One word past honest.

- [NIT] TestPublishRefusesExistingReleaseDirAndSymlinks O_NOFOLLOW half still tautological
  (round-02). The test pre-creates the 3.0.0/ directory, so Lstat(dir) fires before OpenFile
  with O_NOFOLLOW is reached. The changelog cites "O_EXCL|O_NOFOLLOW" — O_EXCL is pinned, the
  O_NOFOLLOW file-level path is not.

- [NIT] TestPublishRejectsUnsafeVersions has no padded case (round-03). ValidVersion rejects
  untrimmed input (core.go:51) but no test pins it.

- [NIT] IMPLEMENTATION.md:69 "GOOS=windows and GOOS=linux builds are now part of the check" —
  still no CI check. grep -rn GOOS scripts/ .github/ returns 0 matches.

- [NIT] IMPLEMENTATION.md:75 "hash detection" and "no agent-accessible write path" in the Notes
  for reviewers. Written in the initial IMPLEMENTATION.md, not updated in later cycles. Same
  overclaims as the §7 text and CLI refusal text.

- [NIT] docs/cli-reference.md still has no protocol command documentation; root CHANGELOG.md
  has no protocol publish/render/check/status entry (kimi-1 round-04 open list).

- [NIT] Skill fallback COOPERATION.md (third copy) is stale — does not contain the §7
  blast-radius clause. The commit was made (parley-deck-skill@455aafe per cycle 1) but the
  locally installed copy predates it. Out of repo scope.

- [NIT] A deck line of genuine prose starting with **Protocol synced:** is silently skipped in
  droppedContent (render.go:220). Very unlikely in practice; the prefix is a machine-generated
  stamp format. kimi-1 round-03 NIT.

## New findings (by severity)

none

## Test-quality assessment

23 protocol tests, all pass. go build ./..., GOOS=windows go build ./...,
GOOS=linux go build ./..., go vet ./..., go test ./... all exit 0.

Command output (this machine, 1903e49, clean tree):

```
$ go build ./...                  EXIT=0
$ GOOS=windows go build ./...     EXIT=0
$ GOOS=linux go build ./...       EXIT=0
$ go vet ./...                    EXIT=0
$ go test ./... -count=1          26 packages: 24 ok, 2 no-test-files, EXIT=0
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore|TestLoad' -count=1 -v
  23 tests, 23 PASS, 0 FAIL
$ grep -c '^func Test' internal/app/protocol_test.go
  23
```

### Revert verification — each claimed reversion-sensitive test

Method: fresh scratch tree from git archive 1903e49 into /tmp/pd-r5, surgical revert per test,
revert verified landed by grep before trusting any result, project repo not modified.

1. **TestProtocolRenderDetectsCrossSectionErasure** (cross-section fixture) — REAL. Reverted
   droppedContent to global multiset (replaced indexBySection with flat map, changed match to
   global): FAIL. The fixture puts "Core rules live here." under §2 in the deck while the render
   has it only under §3. The global matcher consumes it from §3 and does not report §2; the
   per-section matcher reports it as dropped from §2. The assertion (contains "## 2. Active
   agents") distinguishes the two. Verified both directions.

2. **TestLoadRefusesASymlinkedReleaseDirectory** — REAL. Removed the Lstat(dir) symlink check
   in Load (core.go:97-99): FAIL ("Load read through a symlinked release directory and got
   'PLANTED'"). Pins the release-directory half of the Load symlink fix.

3. **TestPublishRefusesASymlinkedStoreComponent** — REAL. Removed assertNoSymlinkComponents
   call from Publish (core.go:148, kept Load's call at core.go:93): FAIL ("published through a
   symlinked store component"). Pins the write half that the changelog claims alongside the
   read half.

4. **TestProtocolRenderReportsDeckReadErrors** — REAL. Restored prior, _ := os.ReadFile
   (swallowed read error): FAIL ("render treated an unreadable deck as empty"). Uses --dry-run
   so the failure cannot come from the write also failing. The second-pass fix is genuine.

All four tests that cycle 4 claims as reversion-sensitive ARE reversion-sensitive. This is the
first cycle in five where no claimed test is tautological, vacuous, or passes for the wrong
reason.

### Does any text still claim a guarantee with no end-to-end test?

The following are carried NITs from prior rounds, restated so they are not lost. None is a new
finding, none was claimed as fixed in cycle 4, and none blocks exit:

1. §7 "no agent-accessible code path writes a release" — overclaim for pty-agents, all three
   copies. The pty limitation IS stated inline in the same sentence. (NIT, carried)
2. Changelog "no agent-accessible write path" — same overclaim. (NIT, carried)
3. protocol.go:301 "hash detection" — ambiguous, one word past honest. (NIT, carried)
4. Changelog "O_EXCL|O_NOFOLLOW" — O_NOFOLLOW file-level path still not pinned by a test (the
   existing test's O_NOFOLLOW half is tautological). (NIT, carried)
5. IMPLEMENTATION.md:69 "GOOS builds are now part of the check" — no CI check exists. (NIT,
   carried)
6. IMPLEMENTATION.md:75 "hash detection" and "no agent-accessible write path" — same overclaims
   in the Notes for reviewers. (NIT, carried)

The following claims ARE tested and survive revert verification:
- Write-once releases → TestCoreReleasesAreWriteOnce
- Attended-only publish (TTY gate) → TestPublishRefusesWithoutATerminal (fails on revert)
- Per-section matching → TestProtocolRenderDetectsCrossSectionErasure (fails on revert)
- Load store-component symlink refusal → TestLoadRefusesASymlinkedStore (fails on revert)
- Load release-dir symlink refusal → TestLoadRefusesASymlinkedReleaseDirectory (fails on revert)
- Publish store-component symlink refusal → TestPublishRefusesASymlinkedStoreComponent (fails on
  revert)
- Render read-error refusal → TestProtocolRenderReportsDeckReadErrors (fails on revert)
- Path traversal refusal → TestProtocolRejectsPathTraversalInTheLock + TestLoadRefusesToEscapeTheStore
  (both fail on revert)
- Production dispatch reachability → TestProtocolIsReachableThroughProductionDispatch
- Status read-error reporting → TestProtocolStatusReportsReadErrors (fails on revert)
- DETECTED-UNATTRIBUTED — correctly hedged as "ratified but not implemented" in all three copies
  and the changelog; 0 Go references; no test claims it.
- 23-test count — reproduces (grep -c '^func Test' = 23).

### Does droppedContent still miss content or report content that IS carried?

The per-section matcher is correct for the cases the reviewers demanded. The one residual false
positive — a core-side section rename that keeps body lines verbatim — is the inherent cost of
per-section strictness, was noted by kimi-1 in round-04, and was kept deliberately (a data-loss
report should err loud). No content is silently missed: a line is reported as dropped only if the
render has no unconsumed copy of it in the same section. The global masking that let a line
deleted from §3 hide behind an identical line in §11 is closed and now genuinely pinned by a
test that fails on revert.

The **Protocol synced:** prefix skip (render.go:220) remains: a deck line of genuine prose
starting with that prefix would be silently skipped. This is a carried NIT (kimi-1 round-03),
very unlikely in practice, and not a new finding.
