---
agent: codex-1
idea: protocol-read-cost-regression
review-round: 4
date: 2026-08-11
reviewed-commit: 41e6cd6 (v1.43.1)
responding-to: [codex-1/review/round-03, claude-1/IMPLEMENTATION-fix-up-cycle-4]
---
verdict: CLEAN

## Summary

v1.43.1 adopts my round-3 counter-proposal in full. `runner.go` and `phase58.go` are
byte-identical to `d4256a2`; `frontier.go` and `frontier_test.go` are absent; and no frontier,
ledger, fallback-banner, or changed round-instruction residue remains in `internal/runner`.

Finding B is closed by deletion, not relocated. The unreachable implementation and its textual
guards no longer exist in executable source, so a future one-line enablement cannot expose the
untested path I objected to. The ledger contract remains only in the idea artifacts, where it carries
no runtime risk.

## Refutation attempts

### 1. Exact pre-idea restoration

PRIMARY — executed against the shipped worktree:

```text
$ git diff --exit-code d4256a2 -- internal/runner/runner.go internal/runner/phase58.go
(no output; exit 0)
```

This is stronger than a functional approximation: both shipped files are exactly the pre-idea
versions. Inspection of the restored paths also shows the original direct calls:

- design cross-review calls `gatherPriorRounds` directly;
- review and review-consensus both call `gatherReviewContext` directly;
- `gatherPriorRounds` and `gatherReviewContext` skip only `_index.md` among Markdown artifacts;
- `BuildRoundPrompt` again says `READ every prior-round artifact below`.

The two active input changes from round 3 are therefore gone. `_ledger.md` is no longer excluded,
and the ledger/banner instruction wording is no longer sent to an agent.

### 2. Finding B deletion and relocation search

PRIMARY — the shipped tree and worktree contain neither frontier file:

```text
$ git ls-tree -r --name-only HEAD -- internal/runner/frontier.go internal/runner/frontier_test.go
(no output)
$ test ! -e internal/runner/frontier.go && test ! -e internal/runner/frontier_test.go
(exit 0)
```

PRIMARY — a search under executable runner source for `frontierContext`, `ledgerFileName`,
`_ledger.md`, the carry-forward-ledger instruction, and the fallback banner returned no matches.
Nothing still carries the constant-false branch, its guards, or its activation risk.

### 3. Agent-visible behavior

Because both input-building files are identical to `d4256a2`, the protocol-read-cost feature changes
neither design-round nor review-round context compared with the pre-idea implementation. This closes
both round-3 input deltas rather than merely disabling compaction around them.

The overlay work shipped in the same release is separate, intentional functionality. Its protocol
description is accurate in the three shipped copies I checked:

- `parley-deck/COOPERATION.md:767-771`;
- `internal/protocol/defaults/COOPERATION.md:758-762`;
- the parley-deck-skill 2.7.0 packaged fallback `references/COOPERATION.md:758-762`.

All three say that the grammar, v2 lock, terminal-boundary composition, and extend-only behavior
exist, while the roster-annotation identity slot and removal of prose-matched zone addressing do
not. They do not promise the full overlay as already complete.

The per-runtime skill snapshots installed on this machine are still version 2.6.0 and retain the
older wording. `parley-deck-skill status` reports that explicitly as runtime version drift; the
released 2.7.0 packaged copy itself is correct. That local cache/update state is not residue in the
v1.43.1 artifact.

### 4. Build and tests

```text
$ go build ./...
PASS

$ go test ./internal/runner -count=1 -run 'Test(BuildRoundPromptIncludesPriorAndRound|GatherPriorRounds|RunRoundCrossReviewWithHeadlessAgent|RunReviewRoundWritesReviews|BuildReviewPromptRefutationDefault)$'
ok  parley-deck-cli/internal/runner

$ go test ./internal/runner ./internal/protocolcore ./internal/app -count=1
internal/protocolcore: PASS
internal/app: PASS
internal/runner: the previously recorded environment-only
TestDurableKillEndToEndRealProcess failure (`no recorded boot id`); no other package failure
```

The targeted tests exercise the restored prompt and full-history paths. The broader runner failure
is the same host-process-attribution limitation recorded in round 3 and is unrelated to this
rollback.

## Findings

None.

## Round-3 disposition

- **Finding A:** closed. Both active prompt-input changes are removed by the exact source restore.
- **Finding B:** closed. The dormant implementation and source-text guards are deleted, with no
  executable relocation.

## Release judgment

There is no code, prompt, protocol-copy, or verification issue in the reviewed shipped state that
should make the owner yank or supersede v1.43.1.
