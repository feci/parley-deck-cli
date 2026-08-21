### opencode-1 — ✅ accept

PRIMARY: Did the fix-up cycle weaken any test? **No.**

Fix-up commits (`git log --oneline -6`):
- `a27a3b6` two FINAL.md files declared a status their own gate rejects — **no `*_test.go`**
- `4903b47` FINAL drafter prompt + pipeline completion test — **two new test files only**
- `0eb83e2` implementer may SIGN a review consensus even though it is not AWAITED — **one new test appended; no existing assertion edited**

No existing test was loosened. No fixture was padded so a change would pass. The three new tests fail without their source fix (verified in a COPY: worktree, `git checkout <parent> --` the three source files, tests left at HEAD).

## Per-file verdicts

### `internal/consensus/roundgate_test.go` — not weakened (`0eb83e2`)

Added `TestTheImplementerMaySignAReviewConsensusEvenThoughItIsNotAwaited`. Existing tests untouched.

Without `consensus.go` known/required split: **FAIL** `the implementer's own signoff was rejected: [line 27: unknown participant impl]`. Then asserts `TriageReady` and that `impl` is still not in `Missing`. New contract, not a loosened old one.

### `internal/app/finalprompt_test.go` — not weakened (`4903b47`, new file)

`TestFinalDraftPromptDescribesWhatTheGateRequires` requires every `protocol.RequiredFinalSections` heading, `idea: my-idea`, `status: final`, and the `N/A` affordance.

Without `driver_consensus.go` prompt rewrite: **FAIL** on six missing headings, the slug, and N/A. `status: final` already lived in the old prompt, so that one assertion does not fire — equal strictness, not a skip.

### `internal/app/pipeline_scaffold_test.go` — not weakened (`4903b47`, new file)

`TestScaffoldWithStatusFinalIsNotACompletedBlock`: unwritten `status: final` scaffold must not be complete; a written FINAL must; no-status must not.

Without `isFinalized` + `FinalIsScaffold`: **FAIL** `an unwritten scaffold declaring status: final was treated as a completed block`. Positive and no-status cases keep polarity.

### `a27a3b6` — N/A for tests

Two `FINAL.md` `status:` values retargeted to `final`. No test file. Legitimate artifact contract, not a test weaken.

SECONDARY: none on this slice.
