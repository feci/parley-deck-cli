---
idea: meta-protocol-change-ste100-output-format
status: fix-up-cycle-1
implementer: codex-1
started: 2026-10-13
branch: ste100
head-commit: 99d04b3ecd0740fb834e70ec109c9384c7b7cd4d
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
- [x] Run required CLI and skill checks. Record hashes, package contents, and packet sizes.
- [x] Publish draft implementation PRs. Freeze both source commits for Zcode's independent review.
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

Focused prompt, syntax, drift, and size checks pass. The first full CLI suite exited 1.
`TestLiveMapIsCompleteHonorsNeverCutAndRendersEveryPhase` expected 69 blocks; the new appendix makes 70.
Linux jobs in runs 38401898310 and 38401910409 also failed the Zcode signoff fixture with `argument list too long`.
The full protocol plus launch task crossed the per-argument kernel bound. This is a real transport regression.
CI links: https://github.com/feci/parley-deck-cli/actions/runs/38401898310 and https://github.com/feci/parley-deck-cli/actions/runs/38401910409.

`go build ./...` and `go vet ./...` each exited 0. Their logs are empty because those commands succeeded without output.
Zcode independently ran both and confirmed their success. The suite failure still blocks source acceptance.
The skill suite passed 402 Node tests, 54 Python tests, and six payload manifests.
Package dry-run, native portable build, all-target dry-run, and a fresh six-skill fixture doctor passed.
Doctor against existing 2.18.0 runtimes reports the missing new reference. Delivery will update and recheck them.
All three protocol copies and staged core equal their bases plus only H1 and H2.
The skill core is 19,482 bytes. The original packet fixture is 68,874 bytes.
The actual unflagged packet is 68,935 bytes. The meta-idea packet is 71,934 bytes.
Both packet variants classify Appendix B as on-demand and retain the section 6 pointer.
See `source-context/protocol-transform.json` and `source-context/source-sizes.json`.

## Outcomes & Retrospective

Pending source and delivery verification.

## Review cycle 1

Zcode process `8d194a1c-c163-4c81-a0e1-6e5b64196061` exited 0 after 1094.446 seconds.
Its review requested changes: one CRITICAL, one MAJOR, and one MINOR.
Codex accepted all findings in review consensus. Both participants authored ACCEPT signoffs for the fix plan.
Fix-up cycle 1 now changes the launch transport and packet tests. The signed protocol hunks remain exact.

## Fix-up cycle 1 plan

- [x] Add bounded private-file delivery in new `internal/runner/prompt_argument.go`.
- [x] Use it from `buildAgentInvocation` at both headless launch boundaries in runner.go and launch.go.
- [x] Add real child tests for exact full bytes, SHA256, private mode, short prompts, stdin, and file-write refusal.
- [x] Update the 70-block fixture and test Appendix B classification and the retained section 6 pointer.
- [x] Record full test and CI outcomes. Re-run the required suite, build, and vet.
- [x] Build the candidate CLI. Use it for Zcode's next full source review.

The review consensus authorizes this necessary transport deviation. It adds one helper, launch wiring, and focused tests to the original file list.
No H1/H2 byte, metadata parser, authority rule, packet limit, or roster setting will change.

## Fix-up cycle 1
status: complete
completed: 2026-10-13
head-commit: 99d04b3ecd0740fb834e70ec109c9384c7b7cd4d

### Fixes applied

- Repair the CRITICAL argument-limit regression with `prompt_argument.go` and both launch call sites.
- Preserve the complete prepared prompt in a unique mode-0600 file under its private invocation directory.
- Bound each headless argument below 120 KiB on every host. The same constant retains the existing interactive refusal threshold.
- Keep stdin exact. Require full-file read and SHA256 verification in the short argument instruction.
- Refuse file-publication failures and any remaining oversized argument. Preserve short prompt behavior and embedded templates.
- Add `prompt_argument_test.go` with real child execution, all built-in templates, stdin, exact content, private modes, and refusal cases.
- Update `internal/protocolpacket/packet_test.go` to 70 blocks. Test the on-demand appendix and retained binding pointer.
- Record failures, exit states, CI links, and exact packet request arguments in the canonical record.

### Deviations from agreed fixes

None. The source test run is active; no complete or release claim follows from the cycle's applied fixes.

### Transport trade-off and consult disposition

The threshold applies on all hosts. Large macOS argument launches also use the file, although macOS could accept some inline prompts.
This keeps the model instruction and test path consistent across hosts.
File delivery does not prove the model read the file. The participant must obey the read-and-hash instruction.

A separate fresh Zcode consult supports the same repair. It exited 0 after 543.917 seconds.
Its native stdout and receipt are retained in source-context. It supplies advice, not a review vote.
Adopt its shared threshold and all-adapter tests. Retain unique file creation instead of replacement of a fixed file.
Do not add its optional sidecar: the argument includes the digest, and the private file preserves the exact prepared prompt.
The consult's first bullet confuses the payload hash with the defaults-file hash.
The correct defaults hash remains `ea8f26f11acfe8834061fe286a0812bcf7555ebb14337b501d90006ab893bd05`.
The consult stdout stays unchanged. No acceptance claim uses its mistaken attribution.

### Current validation

Focused transport, prompt, map, existing quota-signoff, and interactive-bound tests exited 0.
The required full CLI suite, build, and vet now run against the fixed source.
The skill source remains at `92e3dc78da2f1c55d35c905cabe57fc69ede0670`; no skill product bytes changed in this cycle.

## Review round 2 and completed check receipt

Zcode process `a2332e7d-807e-47bc-b664-abfd970c901c` exited 0 after 818.215 seconds.
Its complete review accepted the implementation. It resolved the CRITICAL, MAJOR, and MINOR findings from round 1.
It found one NIT: stale `head-commit` metadata. This receipt sets both values to the reviewed fix commit, `99d04b3`.
No product byte changed after that commit. Zcode explicitly agreed to preserve the accepted fix consensus as `review/consensus-cycle-1.md`.

The full local CLI test, build, and vet process exited 0. Each chained command therefore exited 0.
The test log records completed package results. Build and vet succeeded without output.
The receipt records command names, exit states, log hashes, and exec session 98777.
See `source-context/fixup1-validation-status.json` and its three retained logs.
Both Linux jobs and both macOS jobs passed at the reviewed source commit.
Both Windows jobs remain active at this receipt. Their final states will be captured before source close.
CLI Windows support remains experimental. CLI WinGet remains held.

The reviewer read an earlier partial test log during its process. Its local-suite pending statement is historical.
The completed receipt resolves that timing gap. The fresh source goal process must inspect it independently.
`parley consensus reopen --review` refused with `cannot reopen consensus with triage=ready`.
The recorded fallback archives the exact signed bytes before a fresh consensus, as Zcode agreed.
No existing peer artifact or signoff changed.
