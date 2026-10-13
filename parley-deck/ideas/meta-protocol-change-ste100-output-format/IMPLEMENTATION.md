---
idea: meta-protocol-change-ste100-output-format
status: in-progress
implementer: codex-1
started: 2026-10-13
branch: ste100
head-commit: d51e1c7
design-pr: https://github.com/feci/parley-deck-cli/pull/83
implementation-pr: https://github.com/feci/parley-deck-cli/pull/84
---

## Summary of work

Implement the merged FINAL in the CLI and skill worktrees.
Codex is organizer, participant, drafter, and implementer. Zcode supplies independent review and a fresh goal check.
Source acceptance precedes merge. A separate channel check follows release and installation.

## Implementation plan / checklist

- [x] Apply exact H1 and H2 to all three protocol copies and staged core 2.19.0.
- [x] Add the skill reference, short pointer, template correction, and one on-demand map entry.
- [x] Add one shared instruction to all thirteen builders. Test real paths and preserve input bytes and syntax.
- [x] Update CLI 1.55.0 and skill 2.19.0 release data and source metadata.
- [ ] Run required CLI and skill checks. Record hashes, package contents, and packet sizes.
- [ ] Publish draft implementation PRs. Freeze both source commits for Zcode's independent review.
- [ ] Obtain both review-consensus ACCEPTs and a fresh Zcode source goal PASS with current-tree AC evidence.
- [ ] Merge and publish authorized channels. Install every managed and generic target, then verify hashes.
- [ ] Obtain a separate Zcode channel check. Write the final owner handoff and verify shared memory.

## Deviations from FINAL.md

Add the generated `skills/parley-deck/parley-addon.json` to the skill file list.
The real installer rejects changed core payload bytes until this manifest is regenerated.
`npm run manifest:addons` updated only that manifest. The installer test then passed.
This is required package integrity metadata; no additional behavior or policy changed.
Also update the existing doctor test in `test/installer.test.js` to expect the new reference.
The full suite exposed that fixed file list. This adds no installer behavior.

## Notes for reviewers

Check the exact protocol delta, older authority boundary, exemptions, advisory style treatment, and thirteen prompt paths.
The mechanical single-reviewer exception is false. Use the controlling brief's four explicit close conditions.
Do not count Codex as a code reviewer.
The meta-idea packet already exceeds 70,000 bytes. Preserve the existing fixture guard and report the actual variant.
Native Windows validation remains outside the available environment. CLI WinGet stays held.

## Progress

- 2026-10-13 12:55Z: Design PR 83 merged with both ACCEPT signoffs. No product edit preceded this plan.
- 2026-10-13 13:02Z: Source implementation is committed. CLI PR 84 and skill PR 12 are open.
- Full validation remains active. Zcode will review the frozen source in a separate process. Delivery remains open.

## Decision Log

The owner-provided `ste100` branches remain in use. Do not allocate, declare, or prune worktrees.

## Surprises & Discoveries

The first full skill run lacked the locked `commonmark` dependency and exposed a stale doctor expected-file list.
`npm ci` restored the declared dependencies. The test now includes STE100.md. The full skill preflight was restarted.

The D6 legacy accounting gate refused round 2 before launch.
The recorded native fallback produced Zcode's real independent artifact. No legacy declaration was applied.
The design signoff command worked without fallback.

## Validation evidence

Focused prompt, syntax, drift, and size checks pass. The full CLI and skill checks are in progress.
All three protocol copies and staged core equal their bases plus only H1 and H2.
The skill core is 19,482 bytes. The original packet fixture is 68,874 bytes.
The actual unflagged packet is 68,935 bytes. The meta-idea packet is 71,934 bytes.
Both packet variants classify Appendix B as on-demand and retain the section 6 pointer.
See `source-context/protocol-transform.json` and `source-context/source-sizes.json`.

## Outcomes & Retrospective

Pending source and delivery verification.
