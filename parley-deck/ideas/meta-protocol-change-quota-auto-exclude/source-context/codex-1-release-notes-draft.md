# Release notes draft — quota auto-exclusion

Producer draft for independent review. Publication follows the pre-confirmed close conditions.

New ideas can continue after a participant exhausts an allowance when the complete failed batch leaves at least two usable non-facilitator participants. The change applies to the current idea's quorum, with one durable transition and an owner notice; machine and deck roster files are unchanged. Existing ideas retain their recorded policy and scope.

The supported zcode recognizer implements the owner's explicit stderr exception: a failed process with no completed artifact or later success, agreeing provider error records containing HTTP 429 exhaustion and a machine reset at least 60 minutes away, and a terminal turn-failure line. All provider error records must agree; missing or contradictory resets and quoted assistant/tool text alone do not qualify. A timezone-free display clock is supplemental only when it agrees with a machine reset in the same record under the bounded offset rule. The owner accepts the remaining possibility that a subagent's quota error accompanies an unrelated root failure. Claude text, Codex and Kimi remain diagnostic-only for automatic exclusion without established native evidence.

**Owner-accepted known limitation for this release:** native-positive AC2 is NOT MET
and is owner-waived by round05-answer Q2 (2026-10-06). zcode auto-exclusion may not fire
on real native output; unrecognized failures fall back to the owner-confirmed path.
R5-MAJOR-2 is deferred, not fixed. Complete source-derived positives are not native
verification. The linked follow-up is
[quota-zcode-native-exhaustion-capture](../../quota-zcode-native-exhaustion-capture/00-prompt.md).

Capture is not automatic in this release. On the next ordinary zcode failure, keep
that run's private, unscrubbed `parley-deck/runs/<run-id>/agents/<agent-id>/stderr.log`
in place; never commit or copy it. Its completeness and cleanup are unverified.
The follow-up will establish a complete scrubbed capture before changing the grammar.

Kickoff and mid-idea application are both required for this release. Mid-idea changes preserve historical signoffs, vetoes, disputes, findings, failed invocations and incomplete files. Required signers and dispatch consumers use current membership. A designated or pinned implementer, or a participant that has started a consensus/FINAL draft, cannot be automatically removed. Pending transitions block further mutation until reconciled, while status/wait/organizer views report them without repairing state.

Policy-off changes against CLI 1.50.0 are explicit: `run --yes` filters confirmed
exclusions, and a bare preflight 503 cannot be excluded with `--yes`. During round 1,
any policy-off join or return uses a plain participant edit, as before;
removing a stale exclusion marker does not change that. After round 1, a new joiner
or kickoff-excluded return must import its late round-1 artifact before signing or
completing quorum; an exact `NON-PARTICIPANT` decline remains in Missing signoffs.
Pending catch-up stops ordinary driving until the manual `agents exec` command
completes it. Membership edits are frozen on final/closed ideas. Legacy ideas without
quota history retain their prior behavior.

Notices are owner-owned and non-authoritative. After a valid applied receipt, notice
edits, archival, deletion and pathname changes never gate membership, signoff or driving,
and ordinary publication is not repeated. At most one ordinary notice is published per
settled transition, with one when the destination is safe; an unsafe destination can
receive none. Before the receipt, any regular live or archived
copy is preserved without interpreting its contents. If absent, checked exclusive
publication is attempted once; unsafe paths and publication failures produce non-blocking
diagnostics while membership and evaluation receipts still complete. A receipt records
that the publication step finished, not proof of delivery. Genuine receipt/history,
writer and projection corruption still gates. Replay never adds a terminal evaluation.

There is no quota polling, retry worker, automatic same-idea rejoin or organizer failover.

This idea has one non-facilitator, so the rule cannot reduce its quorum. It retains the separate claude-1 review and owner-confirmed close. Windows remains experimental and CLI winget publication stays held. The legacy run-accounting blocker is a separate owner-requested follow-up (D6).

Windows CLI 1.51.0 is known broken for new idea creation: `parley run` fails
  at kickoff directory sync with "Access is denied", including policy-off launches.
  The same directory-sync failure prevents driving/signing ideas with enabled
  mid-idea quota scope and affects manual imports, owner revisions and transitions.
  Legacy ideas without quota records skip these new paths. Windows CLI assets remain
  experimental and CLI winget is held. The affected sites and CI evidence are handed
  to the unmerged `windows-portability` branch/idea; its current design does not
  promise to restore these operations. The skill installer is a separate artifact.

Symlinked deck/idea scopes are explicitly refused where quota leases or scoped
  integrity escalations are required; enabled mid-idea creation refuses before writes.
  Use a physical `parley-deck` directory or disable `quota_auto_exclude` for ordinary
  driving. Legacy, policy-off and kickoff-only ordinary driving keep their prior
  behavior; manual/owner revisions still require a physical scope. Workspace ancestor
  aliases, including macOS `/tmp`, remain supported.

Kickoff notices also use checked, non-blocking publication and create an absent safe inbox.

Candidate validation on 2026-10-08: Windows compilation succeeds, but its test suite
fails, including quota kickoff directory-sync access errors and an unchanged Unix-only
`syscall.Mkfifo` test. Windows assets remain experimental; CLI winget is held. Linux and macOS
candidate CI and the full skill suite pass. The host full Go suite, build, vet and race
checks also pass; independent final review remains separately recorded before publication.


## Accepted round-08 release limitations

1. **Kickoff blocking escalation (R8-MINOR-1).** If the kickoff quota decision
   would leave fewer than two usable participants or excludes a protected role,
   an absent or unwritable `parley-deck/inbox/` prevents the blocking note from
   being saved. The command still stops without applying the exclusion or
   creating an idea, but prints only the file error rather than the candidates
   and arithmetic. Keep a writable inbox directory, including in a fresh clone,
   until this is fixed in the follow-up.
2. **Kickoff notice crash window (R8-NIT-1).** A process crash after writing the
   kickoff run manifest and before publishing its notice can leave the notice
   permanently absent. The kickoff record, marker, status and organizer brief
   still show the exclusion. This case was identified by source review; no crash
   was injected. Mid-idea receipt/replay behavior is unchanged.
3. **Aliased decks and plain edits (R8-NIT-2).** On a symlinked deck, a plain
   `participants:` edit or confirmed exclusion counts as a manual revision and
   requires a physical deck path even when `quota_auto_exclude` is false.
   Until that path is restored, the pending edit blocks all signers and driving.
   Disabling the policy again does not resolve this case; the generic diagnostic's
   suggestion to do so is inapplicable.

These are accepted and deferred, not fixed. The AC5 escalation-detail obligation and
AC15 kickoff crash-window obligation remain unmet with the explicit exceptions in
the round08-answer relay note. Follow-up: [quota-kickoff-reporting-and-alias-guidance](https://github.com/feci/parley-deck-cli/blob/quota-auto-exclude/parley-deck/ideas/quota-kickoff-reporting-and-alias-guidance/00-prompt.md).
