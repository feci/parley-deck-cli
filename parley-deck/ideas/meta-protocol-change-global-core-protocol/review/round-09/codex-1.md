---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 9
date: 2026-08-07
reviewed-commit: c7f0c6e
verdict: FINDINGS
---
## Summary

Cycle 8 fixes the reported equal-stamp/before-prose reproduction and the lone-terminal-CR
behavior. The six requested equal/different × before/after/absent stamp cases produce the expected
report counts, a truly stamp-free core/deck is clean, and an independent 50,000-input oracle found
no invalid or non-optimal Hirschberg mask. The required native/Windows/Linux builds and vet are
clean. The 33 protocol/diff tests are green.

This is not CLEAN. The replacement rule still has no proof that the line after `**Created:**` was
actually generated: if a deck has no old generated stamp and genuine stamp-prefixed prose occupies
that slot, the prose is forgiven and would be deleted without a removal report. The cycle-8 claim
that deterministic reversion-sensitive tests were added for all three cycle-7 fixes is also false:
only the stamp reversion fails; compiling whitespace and full-DP reversions leave the 33-test gate
green. The terminal-CR fix is likewise unpinned. The carried G7b documentation/test mismatch and
the false agent-write-path absolute remain unchanged.

`go test ./...` is not green on this host because the same runner boot-ID failure recorded in
rounds 6-8 persists. Cycle 8 did not touch `internal/runner` (`git diff --quiet
c4a8b83..c7f0c6e -- internal/runner` exits 0), so I do not attribute it to this cycle.

Required command evidence at `c7f0c6e91d7e9ee554a65c71a62880f286f18be2` (`PRIMARY`):

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

$ go vet ./...
[no stdout/stderr]
exit 0
```

```text
$ go test ./...
?   parley-deck-cli/cmd/parley [no test files]
ok  parley-deck-cli/internal/acp (cached)
ok  parley-deck-cli/internal/agents (cached)
ok  parley-deck-cli/internal/app 53.156s
ok  parley-deck-cli/internal/config (cached)
ok  parley-deck-cli/internal/consensus (cached)
ok  parley-deck-cli/internal/driver 1.681s
ok  parley-deck-cli/internal/fsutil (cached)
ok  parley-deck-cli/internal/hitl (cached)
ok  parley-deck-cli/internal/loop (cached)
ok  parley-deck-cli/internal/pipeline (cached)
ok  parley-deck-cli/internal/procctl 0.410s
ok  parley-deck-cli/internal/protocol (cached)
ok  parley-deck-cli/internal/protocolcore (cached)
ok  parley-deck-cli/internal/repomap (cached)
ok  parley-deck-cli/internal/retro (cached)
ok  parley-deck-cli/internal/runaction (cached)
ok  parley-deck-cli/internal/runcontrol (cached)
ok  parley-deck-cli/internal/runmanifest (cached)
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL parley-deck-cli/internal/runner 7.901s
ok  parley-deck-cli/internal/runplan (cached)
ok  parley-deck-cli/internal/runstate (cached)
ok  parley-deck-cli/internal/sessionstore (cached)
ok  parley-deck-cli/internal/steer (cached)
ok  parley-deck-cli/internal/store (cached)
ok  parley-deck-cli/internal/track (cached)
ok  parley-deck-cli/internal/tui 0.620s
FAIL
exit 1
```

The implementation's scoped count reproduces. There are 31 top-level tests in
`internal/app/protocol_test.go` and 2 in `internal/protocolcore/render_test.go`:

```text
$ go test -count=1 -v ./internal/app ./internal/protocolcore \
    -run '^(TestProtocol|TestPublish|TestCore|TestLoad|TestDroppedContent|TestLargeDocument|TestLCS|TestStampExemption|TestHirschberg)' \
  | awk '/^--- PASS: Test/ {p++} /^--- FAIL: Test/ {f++} /^--- SKIP: Test/ {s++} END {printf "PROTOCOL_DIFF_TOP_LEVEL pass=%d fail=%d skip=%d total=%d\n", p, f, s, p+f+s}'
PROTOCOL_DIFF_TOP_LEVEL pass=33 fail=0 skip=0 total=33
exit 0
```

## Prior findings: fixed or not

- **[MAJOR] Genuine stamp-prefixed prose before an already-current stamp — FIXED for that exact
  reproduction, but the exemption remains incomplete.** I exercised adjacent prose in every
  requested combination through the real render command. `report` below is the number of lines in
  the header removal group:

  | Deck stamp versus render | Genuine prose | Result |
  | --- | --- | --- |
  | equal | absent | no removal report |
  | equal | immediately before | 1 line reported |
  | equal | immediately after | 1 line reported |
  | different | absent | old stamp forgiven; no removal report |
  | different | immediately before | 1 line reported |
  | different | immediately after | 1 line reported |
  | deck stamp absent | absent | no removal report |
  | no stamp and no `Created` slot in either body | absent | no removal report |

  The equal-before checked-in test is reversion-sensitive. However, in the different-before case
  code inspection shows `headerStamp(prior)` binds the prose in the slot, not the stale generated
  stamp shifted below it. The aggregate heading/count happens to remain 1 because the stale stamp
  is reported in the prose's place. More importantly, when there is no stale generated stamp to
  provide that compensating count, the prose disappears silently; see the new MAJOR.

- **[MINOR] Lone terminal carriage return — FIXED in behavior, UNPINNED.** A direct command-path
  test for `strings.TrimSuffix(minimalCore, "\n") + "\r"` passes at c7f0c6e and fails under the old
  all-elements `TrimSuffix("\r")`. But no such test was committed. Restoring the old implementation
  compiles and all 33 checked-in protocol/diff tests still pass.

- **[MAJOR] Cycle-7 fixes lacked deterministic reversion-sensitive tests — NOT FIXED as claimed.**
  The exact compiling reversion matrix is:

  | Reversion, applied alone in an isolated c7f0c6e archive | Compile | Relevant checked-in result |
  | --- | --- | --- |
  | exact `8ed3c4b` multiset, with the cycle-7 fixture repairs retained | pass | all three repaired round-6 probes fail |
  | replacement predicate back to `deckStamp != renderStamp` | pass | new equal-stamp test fails |
  | restore `strings.TrimSpace(l) == ""` skip | pass | all 33 pass |
  | restore full O(n*m) DP table | pass | `TestHirschbergAgreesWithFullDP` and the 20k test both pass |
  | restore terminal-CR trimming on every split element | pass | all 33 pass |

  The first two rows are good regression tests. The other three refute `IMPLEMENTATION.md:351-354`
  (“Three cycle-7 fixes ... Added, and each proven by reverting its own fix”). The cycle-8 diff adds
  only `TestStampExemptionForgivesNothingWhenNothingWasReplaced` and
  `TestHirschbergAgreesWithFullDP`; it adds no whitespace or terminal-CR test, and the LCS oracle
  checks correctness rather than the memory-bound implementation choice.

  A deterministic whitespace assertion is already available: changing the existing loose
  `"#### Parent — "` assertion to the actual `"#### Parent — 5 lines"` passes current code and
  fails the compiling whitespace revert, which reports 3. A similarly small terminal-CR test also
  passes current code and fails its compiling revert.

- **[MAJOR] Hirschberg correctness — CONFIRMED.** I changed the seed independently and compared
  50,000 randomized pairs of length 0-30 against a full-DP LCS. In addition to optimal length, I
  verified that every marked mask is an actual subsequence of the render. Result:

  ```text
  === RUN   TestHirschbergAgreesWithFullDP
  --- PASS: TestHirschbergAgreesWithFullDP (0.26s)
  PASS
  ok  parley-deck-cli/internal/app 0.560s
  ```

  The modified scratch test executed 50,000 trials; because it stops on the first invalid or
  non-optimal mask, the pass means `invalid_masks=0` and `nonoptimal_masks=0`. This establishes
  correctness after cycle 8. It does not pin linear space: under the compiling full-DP revert the
  checked-in oracle passes, and the 20k test also passes in 0.62s on this host.

- **[MAJOR, carried G7b] Landed guarantees without real-entry-point coverage — NOT FIXED.** Five
  independent, compiling mutants each leave all 33 protocol/diff tests green:

  ```text
  remove final-file O_NOFOLLOW                       PASS 33/33
  validate TrimSpace(version), use original path     PASS 33/33
  remove CRLF-core normalization                     PASS 33/33
  force render mode to 0644                          PASS 33/33
  replace fsutil.WriteFileAtomic with os.WriteFile   PASS 33/33
  ```

  Yet `IMPLEMENTATION.md:21-24,37-38,138-150` documents these as landed behavior, while its line 17
  says nothing claims a guarantee it does not implement/test. The changelog's `O_EXCL|O_NOFOLLOW`
  shipped guarantee remains unpinned. `IMPLEMENTATION.md:295` also still says all four round-6
  cases are permanent tests, although there is no duplicate-heading test.

- **[MINOR, carried] Agent-write-path and hash-detection wording — NOT FIXED.** All three protocol
  copies (`parley-deck/COOPERATION.md:777`, the embedded default at line 768, and the skill-source
  fallback at line 768) end their accurate PTY caveat with the contradictory absolute “no
  agent-accessible code path writes a release.” The changelog says “no agent-accessible write
  path,” and `IMPLEMENTATION.md:74-75` repeats it. A binary built from c7f0c6e, invoked by this agent
  with a PTY, disproves the absolute through the real entry point:

  ```text
  $ PARLEY_HOME=<temp> parley protocol publish --version r9-probe \
      --from internal/protocol/defaults/COOPERATION.md
  Published core r9-probe (73efe37fd237) to <temp>/protocol/core/r9-probe
  exit 0
  ```

  `internal/app/protocol.go:301` also calls “hash detection” a durable guarantee, while
  `Store.Load` explicitly verifies no expected release hash. Deck drift detection exists;
  release-tamper detection does not. Name the narrower guarantee.

- **[MINOR, carried] Phase-8 reviewed HEAD — NOT FIXED.** `IMPLEMENTATION.md:8` names c4a8b83 and
  says cycle 8 lands in the next commit; this review targets c7f0c6e. The record remains one cycle
  behind.

## New findings (by severity, or "none")

### [MAJOR] A deck with no old generated stamp can still lose genuine slot prose silently

`headerStamp` (`internal/protocolcore/render.go:359-377`) does not identify a generated stamp. It
returns any stamp-prefixed line immediately after `**Created:**`. The replacement predicate at
lines 223-225 then sees the current render stamp absent from the deck and enables forgiveness.
Consequently this input has no old generated stamp at all:

```text
**Created:** `<d>`
**Protocol synced:** genuine project prose
```

but dry-run emits no `deck content NOT carried forward` block:

```text
preserved from this deck: Workspace, Transport, Created
would regenerate <fixture>/parley-deck/COOPERATION.md from core 1.0.0 (cd92405d09f3).
Nothing was written. Re-run with --yes to apply.
```

Apply would delete the prose. This is the same silent-erasure safety class as the round-8 MAJOR,
and it also explains why “prose can sit in the slot” was not solved merely by adding
`containsLine(prior, renderStamp)`: that condition proves only that the new stamp is absent, not
that the candidate old line was generated.

Fail closed. At minimum, exempt a slot line only if it satisfies the exact generated-stamp grammar
(`**Protocol synced:** core <version> (<hash>)`); if the line is absent or ambiguous, report it.
Prefer preserving/proving the previous generated stamp as metadata rather than inferring ownership
from a user-editable prefix and position. Add the stamp-absent + genuine-slot case alongside all six
equal/different cases.

## Test-quality assessment

The repaired round-6 fixtures are now good and demonstrably fail under the exact defeated multiset.
The cycle-8 equal-stamp test is also useful and fails under a compiling reversion of the replacement
predicate. The 200-case checked-in full-DP oracle is a sound correctness test; the independent
50,000-case run strengthens confidence that Hirschberg is correct.

The suite still does not justify the implementation record. Whitespace reporting, linear-space
selection, and terminal-CR preservation all survive their own compiling reversions. The six-case
stamp matrix also misses the no-old-stamp slot case that remains silent. The five G7b mutants show
that the three protocol copies, changelog, CLI wording, and IMPLEMENTATION.md continue to outrun
real-entry-point coverage.

I found no cycle-8 regression outside the changed stamp/line-splitting surface. Native, Windows,
and Linux builds plus vet pass; all 33 scoped tests pass; and the only repository-wide failure is
the unchanged, previously recorded runner boot-ID test. The current blockers are the incomplete
stamp ownership rule, the false reversion-coverage claim, and the carried G7b/text discrepancies.
