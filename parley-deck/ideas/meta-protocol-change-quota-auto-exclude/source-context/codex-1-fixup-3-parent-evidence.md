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


## Actual baseline differential after W2

The unchanged extended main (`7a87344bd7f4420372a48c711fafdc1b099df1923dcdc97e55f9a6de968b1c12`)
ran against 27e42b8 and product fac40aa on both volumes. D1 through D4 all now match
success/error exit outcomes and consensus triage. D4's missing-signoff diagnostic
additionally retains the declined joiner as missing, unlike the baseline: this is the
explicit G15 requirement not to treat an incomplete joiner as a completed vote. The
canonical BLOCK form, refusal of literal NON-PARTICIPANT status and continued appends
by existing members are unchanged. No closure or membership is granted. D2's extra
missing events-ledger line belongs to the unchanged current probe setup; the permanent
CLI fixtures contain the ledger and test the intended pending diagnostic. Logs are in
`.parley-runtime/quota-implementation/fixup-3/parent-checks/`.

A verification harness initially required raw equality for all three protocol copies.
The embedded bootstrap intentionally substitutes workspace/date/host handles and omits
the project sync line; only that harness assumption failed. Deck/skill remain byte-identical,
and the existing unmodified TestEmbeddedDefaultMatchesLiveDeck passes. No protocol or
normalizer was edited. The first actual crash harness saw its local shell just before
exec; it now waits for the exact verified Python worker. The next attempt looked for
manifest.json instead of the actual run.json. Both failed harness runs remain recorded,
with their owned local workers stopped; no recovery success is claimed for them.
