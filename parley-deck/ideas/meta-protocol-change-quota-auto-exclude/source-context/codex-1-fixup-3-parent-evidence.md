# codex-1 cycle-3 parent supplement

Producer implementation evidence only; no independent verdict. The subprocess's original
`codex-1-fixup-3-evidence.md` is preserved, including its W2 stop. Its package boundary was
an organizer delegation boundary, not a limit on owner-authorized G15/W2. After that
process exited, the parent implemented W2 in `internal/consensus` and moved the pending
signer gate out of the CLI facade into the shared API, before any canonical append.

The existing baseline form `--status block --notes '❌ NON-PARTICIPANT' --counter
'Continue without me'` remains a BLOCK in design consensus. Only a pending policy-off
joiner in the requested current set can use the exception. It grants no Known membership,
revision, owner authority or completed vote. Existing signers may append afterward;
the declined joiner stays missing and consensus stays blocked. Literal NON-PARTICIPANT
status remains unsupported, matching the baseline. No new status or join marker exists.

Tests exercise real CLI app.Run and the lower API, exact canonical bytes, continued
existing-member appends, missing quorum, immutable kickoff/history, duplicate handling,
read-path malformed detection, normal accept/reserve/block and quoted-note rejection,
outsider/review/policy-on rejection and known-excluded denial. Existing retained-veto
regressions remain in the broad host suite. The first new test used the wrong expected
counter-proposal label; the assertion was corrected to the actual unchanged baseline
label. No product assertion was relaxed.

The focused `go test ./internal/consensus ./internal/app -run
'TestQuotaCycle3|TestQuotaExcludedVeto' -count=1` passes. Full frozen host and differential
evidence will be recorded separately. Membership and notice behavior is now documented
in CLI guidance and the skill reference; the 20000-byte core skill cap is unchanged.

No provider/native capture, parser relaxation, roster/config/model changes, worktree
operations or edits to Claude's review/signoff bytes were made. Q2 remains native AC2
NOT MET / owner-waived, R5-MAJOR-2 owner-accepted and deferred.
