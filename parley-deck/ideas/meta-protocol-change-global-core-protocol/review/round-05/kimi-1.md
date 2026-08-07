---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 5
date: 2026-08-07
reviewed-commit: 1903e49
verdict: CLEAN
---

## Summary

Cycle 4 is the first cycle where I cannot construct a finding. The diff against 02bb1e3 touches
only `internal/app/protocol_test.go`, IMPLEMENTATION.md, and the two round-04 review files — no
production code changed, so the entire verification burden was the four test claims, and every one
of them survived the test that matters: I reverted each pinned fix in a scratch tree
(`git archive 1903e49` → `/tmp/pd-r5/src`, each revert verified landed by grep/compile before
trusting any result) and watched the corresponding test FAIL, with the other 22 protocol tests
green each time. The tautology that guarded cycle 3's headline fix is gone: the rewritten fixture
puts `Core rules live here.` at the end of the deck's §2 while the render carries it only under
§3, and under a reverted global multiset the report loses the `## 2. Active agents` entry entirely
— the failure output shows the masking happening, not just an assertion flipping.

Command output (this machine, 1903e49, clean tree):

```
$ go build ./...                  BUILD_OK
$ GOOS=windows go build ./...     WIN_OK
$ GOOS=linux go build ./...       LINUX_OK
$ go vet ./...                    VET_OK
$ go test ./... -count=1          25 packages ok, 2 no-test-files (incl. internal/app 36.249s,
                                  internal/runner 11.485s), EXIT=0
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore|TestLoad' -count=1 -v
  --- PASS lines: 23    --- FAIL lines: 0
$ grep -c '^func Test' internal/app/protocol_test.go
  23
```

IMPLEMENTATION.md's "23 protocol tests" reproduces from an actual run. No runner flake.

## Prior findings: fixed or not

Round-04 findings (cycles 2-4 reviewed together, since codex-1's rounds 3-4 were filtered):

1. **[MAJOR hermes-1 + kimi-1] cross-section test was tautological — FIXED AND PROVEN.**
   The rewritten `TestProtocolRenderDetectsCrossSectionErasure` (protocol_test.go:502-523) is a
   genuine cross-section case: the deck's copy sits under §2 (inserted before the `## 3. Phases`
   heading), the render's only copy under §3, and the assertion names the §2 heading — which under
   global matching appears nowhere in the output (the "preserved" line says "§2 roster table", not
   the heading). Revert A (per-section `indexBySection` → cycle-2 global multiset) produced:

   ```
   --- FAIL: TestProtocolRenderDetectsCrossSectionErasure (0.00s)
       protocol_test.go:521: a line dropped from §2 was masked by an identical line surviving in §3:
           ...
           deck content NOT carried forward by core 1.0.0:
             - ## 3. Phases — 1 line not carried forward
             - ## 99. Project-specific rule — 2 lines not carried forward
   (other 22 protocol tests PASS)
   ```

   The §2 entry is absent under the revert — the erasure is masked exactly as documented, and the
   test comment's "Verified: with matching reverted to global, this version FAILS" is now true.
   Cycle 3's headline fix has real regression protection for the first time.

2. **[MINOR kimi-1] three cycle-3 fixes unpinned — ALL THREE PINNED, each pin proven.**
   - Revert B (delete Load's release-dir `Lstat`, core.go:97-99):
     `TestLoadRefusesASymlinkedReleaseDirectory` FAILS — "Load read through a symlinked release
     directory and got \"PLANTED\""; other 22 PASS.
   - Revert C (delete ONLY Publish's `assertNoSymlinkComponents` call, leaving Load's intact —
     verified by grep count 2→1): `TestPublishRefusesASymlinkedStoreComponent` FAILS — "published
     through a symlinked store component"; other 22 PASS, including both Load symlink tests. The
     changelog's "symlinked store components refused on read and write" is now pinned on BOTH
     halves (read by `TestLoadRefusesASymlinkedStore`, write by the new test).
   - Revert D (restore `prior, _ := os.ReadFile(path)`):
     `TestProtocolRenderReportsDeckReadErrors` FAILS on BOTH assertions — "render treated an
     unreadable deck as empty: would regenerate … Nothing was written" plus "render did not report
     the read failure". The `--dry-run` rework is real: with the read error swallowed, the dry run
     reports a full regeneration and exits 0, so the failure cannot be coming from the write path.
     The fix-up's "initially passed for the wrong reason … now uses --dry-run and asserts the
     reported reason" is accurate.

3. **[NIT kimi-1] core-side section rename over-reports — disposition as claimed.** Kept
   deliberately, recorded in IMPLEMENTATION.md:215-217. I re-probed the class this round (below);
   it still errs loud and only loud. Agreed with the call.

**Behavioral re-verification at 1903e49** (binary built from the pristine scratch tree, direct
store setup, release `## Allowed automatically / - Deploy production.` + `## Requires explicit
user approval / - Delete the database.`):

- Round-03 erasure fixture (deck carries `- Deploy production.` ONLY under approval): preview AND
  apply report `## Requires explicit user approval — 1 line not carried forward`; the apply moves
  the line under the core's section; `check` converges (exit 0); re-render "already matches".
- Deck-side duplicate (line under BOTH sections): approval reports exactly 1, allowed reports 0.
- Fully-carried deck: "already matches", no report.
- Deck nesting a carried line under a core-less `### Details` subsection: reports
  `### Details — 2 lines` (heading + the surviving line). This is the rename NIT's class — a
  heading the core drops takes its lines' count with it even when a line survives elsewhere. Loud,
  never silent; same recorded trade-off, not a new finding.

**Still open from earlier rounds, not claimed by this fix-up** (restating so they are not lost;
all re-confirmed present at 1903e49): the §7 absolute "no agent-accessible code path writes a
release" (parley-deck/COOPERATION.md:775, embedded copy identical, the pty admission two clauses
earlier in the same sentence; the changelog carries the adjacent equivalent with the pty limit
stated inline); the O_NOFOLLOW half of `TestPublishRefusesExistingReleaseDirAndSymlinks` remains
unreachable dead code; protocol.go:301's "hash detection" still one word past honest; a deck line
of genuine prose starting with `**Protocol synced:**` is still silently skipped (render.go:220);
`TestPublishRejectsUnsafeVersions` still has no padded case; `docs/cli-reference.md` still never
documents the `protocol` command and the root `CHANGELOG.md` has no entry for it; preflight.go
still owns the legacy stamp format; the mode-preserving write branch and `check --json`'s
non-zero exit remain untested; the check-preamble rewording (cycle 3) is behaviorally verified but
unpinned (hermes-1 round-04 noted the same); IMPLEMENTATION.md frontmatter still says
`head-commit: (see release commit)`; IMPLEMENTATION.md:69 still claims the cross-builds "are now
part of the check" — there is still no check (`grep -rn GOOS scripts/ .github/` → 0 matches).

**Environmental observation, outside this commit's scope:** the third protocol copy on disk —
`~/.kimi-code/skills/parley-deck/references/COOPERATION.md` — no longer carries the §7 core-store
block at all (38 diff lines vs the deck copy; the directory is no longer even a git repo, so the
`parley-deck-skill@455aafe` sync point is gone from disk; the two repo copies remain byte-identical
modulo the roster/stamp placeholders and are drift-pinned by `TestEmbeddedDefaultMatchesLiveDeck`).
No text at 1903e49 claims the skill copy is current, so this is not a finding against the cycle —
but rounds 3-4 asserted three-copy identity, and that is no longer true on disk. Whoever next
touches the skill should re-sync it.

## New findings (by severity, or "none")

none.

## Test-quality assessment

23 protocol tests (+1 drift test in internal/protocol), all green on `-count=1`. Revert method as
in rounds 3-4; each revert verified landed (grep count + compile) before running, one fix at a
time, the file re-extracted from `git archive 1903e49` between reverts.

The four tests cycle 4 stakes its claims on — all reversion-sensitive, demonstrated this round:

- `TestProtocolRenderDetectsCrossSectionErasure` — real for the first time (was a tautology at
  02bb1e3). Fails under global matching with the §2 entry verifiably absent from the report.
- `TestLoadRefusesASymlinkedReleaseDirectory` — fails without the Lstat; no other guard catches
  the escape (it read "PLANTED").
- `TestPublishRefusesASymlinkedStoreComponent` — fails without Publish's guard; surgically
  independent of Load's guard and of the write-once refusal (both stayed green under the revert).
- `TestProtocolRenderReportsDeckReadErrors` — fails on both assertions with the read error
  swallowed; the dry-run isolation means the write path cannot mask the read path.

No test in the cycle passes for a wrong reason that I could construct. The pre-existing suite's
carried weaknesses are unchanged and listed above (dead O_NOFOLLOW half, unpadded version list,
untested mode-branch/`check --json`/check-preamble).

**G7b sweep — does any text still claim what no end-to-end test covers?** Every claim NEW or
re-asserted in cycle 4 reproduces: "reverting to global matching makes it FAIL (verified both
directions)" (Revert A + the green baseline), "each was proven to FAIL with its fix reverted"
(Reverts B/C/D), "23 protocol tests" (counted), "25 packages ok" (quoted). The changelog's
"refused on read and write" moved from half-pinned to fully pinned. What remains overstated is
exactly the carried list above — all of it pre-dates this cycle, none of it claimed as fixed by
this cycle, and the cycle's own texts (including the test comments' self-described verification)
check out word for word.

Four cycles of documentation outrunning the code, twice with the test as the wrong part — this
cycle the documentation is merely true. CLEAN.
