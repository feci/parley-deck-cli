---
idea: consensus-workflow-cli
cycle: 1
drafted-by: codex
date: 2026-05-12
implementation-pr: https://github.com/feci/parley-deck-cli/pull/10
---

## Agreed fixes

1. Fix the CRITICAL triage bypass: `validateDocument` must use the canonical status returned by `CanonicalStatus` for reservation/block triage and required-field checks. Add a regression test where a manually written `Status: block` signoff cannot become `ready`.
2. Make status parsing more tolerant of protocol/example text:
   - accept status values with trailing parenthetical explanatory text;
   - accept `Counter-proposal:` in addition to `Counter-proposal (required if ❌):`.
3. Fix `reopen` protocol state handling:
   - do not write undocumented `status: discussion`;
   - restore the idea to the latest existing `round-NN` status;
   - name aborted consensus files after the latest round plus an attempt counter, such as `round-02-consensus-aborted-01.md`.
4. Add the missing `### Non-goals` section to the generated `FINAL.md` scaffold and test the required scaffold headings.
5. Add the missing positive test for `reserved` consensus finalization when open items are visibly recorded.
6. Make review consensus draft frontmatter closer to the protocol by using a numeric `cycle` value and including a `reviewed-commit` field, even if the value is empty in this slice.
7. Surface malformed consensus files in the workspace idea list as `consensus=error` instead of silently omitting the label.
8. Sort round directories numerically rather than lexicographically.

## Deferred follow-ups

- Project-level `consensus.*` event storage remains deferred because the current durable event store is run-scoped.
- Cross-process file locking for simultaneous `parley consensus signoff` invocations remains deferred; the slice uses a process-local mutex and append semantics.
- Multi-line `Notes:` parsing remains deferred; the first slice documents and supports one-line structured fields.
- Native GitHub review submission remains deferred to `consensus-github-review-mirror`.
- Automated agent invocation remains deferred to `consensus-request-signoffs`.

## Dismissed findings

- Updating `IMPLEMENTATION.md` automatically on `parley consensus reopen --review` is not required in this slice. Review reopen behavior needs a more complete Phase 8 policy before the CLI mutates implementation status.
- A hard cap on aborted consensus filename attempts is not required once filename generation handles unexpected `os.Stat` errors explicitly.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-05-12
Status: ✅ ACCEPT
Notes: Accept. The agreed fixes address the blocking triage and protocol-state findings while keeping orchestration, GitHub mirroring, and stronger locking as follow-ups.

### Signoff: claude — 2026-05-12
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: Both MAJOR findings I raised (undocumented `status: discussion` on reopen; missing `### Non-goals` in the `FINAL.md` scaffold) are in Agreed fixes, as is Gemini's CRITICAL triage bypass (canonical-status not used for `hasBlock`/`hasReservations`) — that one is correctly the lead item. The MINOR robustness items I care about (tolerant `Counter-proposal:` prefix, positive `reserved + open-items` finalize test, `consensus=error` surfacing in the workspace listing, numeric round sort, review-consensus frontmatter with numeric `cycle` and `reviewed-commit`) are all in. Deferrals for project-level `consensus.*` events, cross-process file locking, and multi-line `Notes:` are reasonable and explicitly logged. Two reservations, neither blocking: (a) the O(N²) workspace-status recomputation in `runStatus` / `consensusTriageLabel` (`internal/app/app.go:447-450`, `:650-656`) is not on the list — fine for the slice but please track as a follow-up since the next slice (`consensus-request-signoffs`) will hit the same path; (b) `Finalize` still has no explicit guard against re-running when `00-prompt.md` `status:` is already `final` (only an implicit `FINAL.md`-exists check) — worth deciding intentional vs. needs-guard in the fix-up cycle or a follow-up.

### Signoff: gemini — 2026-05-12
Status: ✅ ACCEPT
Notes: The agreed fixes address my critical finding regarding triage logic bypass and my major findings on aborted file naming and template omissions. I accept the consensus.
