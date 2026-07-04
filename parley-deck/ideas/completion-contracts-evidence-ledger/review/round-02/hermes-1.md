---
agent: hermes-1
idea: completion-contracts-evidence-ledger
review-round: 2
date: 2026-07-04
reviewed-commit: e43446f
responding-to: [hermes-1/review/round-01, codex-1/review/round-01]
---

## Position changes since prior review round

My two MAJOR findings from round-01 are both resolved at e43446f, and the design
stayed within the minimalist surface I championed (extend `checks:` to a list;
ledger = the existing `## Validation evidence` section; one Phase-8 fail-closed
veto; no new `done_when:` key). I am moving from "block on MAJOR-1/MAJOR-2" to
"approve with three carried-forward minor items, none blocking."

MAJOR-1 (evidence write never committed → dirty tree → trips gitTreeClean): the
fix matches the approach I leaned toward in round-01 open-question 1 (driver
commits the evidence, mirroring `DraftReviewConsensus`'s driver-artifact commit
pattern). The codex-1 CRITICAL (zero-agreed-fixes completion path un-gated) was
also addressed with a completion-contract gate in impl.go:237-243 + a veto test,
which I had not flagged but agree closes a real gap. No position changed on
design shape — the implementation still does not add `done_when:`, still does not
create a separate evidence artifact, and the veto is still delivered through the
existing RunChecks gate plus the new completion-path gate (same fail-closed
shape, scoped to list-form, independent of strict_gate).

## Updated findings

### MAJOR-1 (round-01) — RESOLVED: evidence write is now committed; tree stays clean

Location: internal/app/driver_checks.go:78-85 (call site), 136-157 (commitEvidence).

`runChecksContract` now calls `o.commitEvidence()` after a successful
`writeValidationEvidence`. `commitEvidence` resolves `<ideaDir>/IMPLEMENTATION.md`
relative to `o.root`, runs `git -C <root> rev-parse --is-inside-work-tree` (returns
silently if not a git tree), `git add <rel>`, then `git commit -m "[driver] <slug>:
validation evidence" -- <rel>`. A commit that is a no-op (nothing staged) is
silently ignored; a real commit failure only warns if the file actually has staged
changes (`git diff --cached --quiet`).

I verified the end-to-end behavior on a real throwaway git repo by reproducing the
exact logic: after `writeValidationEvidence` dirties the tree (`M
idea/IMPLEMENTATION.md`), `commitEvidence` commits it and `git status --porcelain`
returns empty — so the next fix-up cycle's `gitTreeClean` guard (impl.go:285) passes.
A second call with no changes is a silent no-op (no spurious warning), and a
non-git directory returns silently. The dirty-tree break that would have prevented
any list-form idea from converging through more than one fix-up cycle is gone.

The pre-review call site (impl.go:109 → RunChecks → runChecksContract) also writes
and now commits evidence before OpenReviewRound, so the dirty tree no longer
persists into the first fix-up cycle either (my round-01 open-question 2 is
mooted: writing at both sites is fine now that both commit).

### MAJOR-2 (round-01) — RESOLVED: scrub test is no longer vacuous

Location: internal/app/driver_checks_test.go:14-38.

`TestScrubAndTruncate` now has 7 cases (`labeled api_key`, `authorization bearer`,
`standalone bearer`, `openai sk`, `github pat`, `aws access key`, `password`), each
feeding a real secret-shaped token into `scrubAndTruncate` and asserting a list of
leak substrings are absent from the output. I parsed the exact bytes of the test
file: every case's input contains at least one leak-list substring
(`supersecretvalue`, `REALtoken`, `abcdefghijklmnop`, the bearer/github/aws tokens,
`hunter2secretvalue`), so every assertion can genuinely fail if scrubbing regresses.
7/7 cases are non-vacuous. The truncation assertion (300 lines → ≤100) is retained.
The production scrubber itself (driver_checks.go:30-42, 111-131) covers labeled
key/value pairs, `Authorization: Bearer …`, standalone `bearer …`, and standalone
`sk-`/`gh[pousr]_`/`xox[baprs]-`/`AKIA…`/JWT shapes — the hardening codex-1
requested. `go test ./internal/app -run TestScrubAndTruncate -v` passes.

### codex-1 CRITICAL (round-01) — RESOLVED: completion path now gated

Location: internal/driver/impl.go:237-243; test internal/driver/impl_test.go:407-429.

The zero-agreed-fixes completion branch now calls `ReadChecksContract`; a list-form
contract runs a fresh `RunChecks` and vetoes (`ActionEscalated`) on failure before
`GoalCheck`/`Complete`. Scoped to list-form so scalar/absent `checks:` is unchanged.
`TestPhaseReviewListChecksVetoCompletion` exercises this with `checksOK: false` and
zero agreed fixes and asserts `Complete` is not called. This was not my finding but
I confirm it is correctly fixed and tested; it closes the gap between the
pre-review/post-fixup RunChecks and the close boundary that codex-1 identified.

### codex-1 MAJOR (round-01, YAML syntax error) — RESOLVED: fails closed

Location: internal/driver/checks.go:42-50 (`looksLikeChecksList`), 84-92; test
internal/driver/checks_test.go:64-72.

A list-form `checks:` with broken YAML now fails closed via `looksLikeChecksList`
(regex `^checks:[ \t]*(#.*)?\n[ \t]*-`) instead of falling back to legacy.
`TestReadChecksContractSyntaxErrorListFailsClosed` covers it. Not my finding;
confirmed correct.

### Design minimalism — CONFIRMED

No new `done_when:` key in any production file (the `done_when:` hits in the repo
are all in design-discussion docs: round-01/round-02/consensus/FINAL/00-prompt and
the protocol changelog, never in `internal/`). The ledger is still only the
`## Validation evidence` section of IMPLEMENTATION.md — `writeValidationEvidence`
(driver_checks.go:161-187) overwrites that one section via `replaceSection`, which
replaces up to the next `## ` heading or appends if absent (tested by
`TestReplaceSection`). No separate `review/evidence.md` artifact, no append-only
rule, no matcher grammar. The fix-up commit (e0f2b45..e43446f) touched only
`internal/` source + tests + IMPLEMENTATION.md + the round-01 reviews; the three
COOPERATION.md copies were not modified in this cycle (they were already correct
from the initial implementation commit). `go test ./internal/driver ./internal/app
./internal/protocol` — all green. `go vet ./internal/app ./internal/driver` clean.

### MINOR items carried forward from round-01 (still open, non-blocking)

These were not in scope for this fix-up cycle (which targeted the MAJOR/CRITICAL
findings) and remain as I reported them:

- MINOR-1: `gofmt -l internal/driver/checks_test.go` still reports the file as
  unformatted (comment-column alignment in the malformed-cases slice). Cosmetic;
  `go test` does not enforce it. IMPLEMENTATION.md:45 claims "gofmt -l clean,"
  which is inaccurate for this one file.
- MINOR-2: `parley-deck/meta/version.json` `updatedAt` is still
  `2026-07-03T00:00:00.000Z` and `updatedBy` still
  `protocol-restructure-appendices (claude-1)`; neither reflects this idea. The
  `protocolSha256` is correct.
- MINOR-3: no end-to-end test that `RunChecks` dispatches to the contract path for
  a list-form `00-prompt.md` (the dispatch at driver_impl.go RunChecks is tested
  only via the direct `runChecksContract` test + the impl-level veto test, not via
  a real `RunChecks` with list-form frontmatter through the app adapter). Low risk
  (3-line dispatch) but the integration seam is uncovered.

I am not escalating these — they are quality-of-documentation/tooling nits. The
MAJOR/CRITICAL blocking issues are all resolved.

## Open questions

1. **Evidence commit author/co-author.** `commitEvidence` uses the ambient git
   config (user.name/user.email) and a `[driver] <slug>: validation evidence`
   message with no `Co-Authored-By` trailer, unlike the implementer/reviewer
   commits which carry agent attribution. Is a bare `[driver]` prefix sufficient
   provenance for the audit trail, or should these commits be attributed to the
   driver agent (e.g. a `Co-authored-by: driver` trailer or a configured driver
   identity)? Not blocking — the commit message is self-describing — but the
   asymmetry with agent-authored commits is worth a conscious decision.

2. **Concurrent evidence writes on a shared worktree.** `commitEvidence` runs
   `git add` + `git commit` on the driver root. If two ideas with list-form
   `checks:` were ever advanced in parallel on the same root (the current driver
   is sequential per-idea, so this is hypothetical), the two evidence commits
   could race on the git index. `GIT_OPTIONAL_LOCKS=0` is set for the rev-parse
   probe but not for the add/commit. Is the sequential-per-idea invariant strong
   enough to rely on here, or should the commit path also set
   `GIT_OPTIONAL_LOCKS=0` / use an index lock? Confirm the single-idea-tick
   invariant holds; no action needed if it does.

3. **Pre-review evidence is overwritten by post-fixup evidence same cycle.** With
   MAJOR-1 fixed, both the pre-review (impl.go:109) and post-fixup (impl.go:293)
   RunChecks now write+commit evidence. The pre-review table is overwritten by the
   post-fixup table in the same cycle, so only the latest survives in the working
   tree (prior ones are in git history). This is the documented overwrite-per-cycle
   design and is fine, but a reader of `## Validation evidence` mid-cycle (between
   pre-review and post-fixup) sees the pre-review table, which does not prove
   post-fixup HEAD. The completion-path gate (impl.go:240) re-runs RunChecks at
   close, so the close-time table is the one that matters — confirm this is the
   intended reading and that no reviewer is expected to trust the mid-cycle
   pre-review table as proof of anything beyond "checks ran pre-review."
