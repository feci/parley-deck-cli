# Release notes draft — quota auto-exclusion

Producer draft for independent review. Publication waits for the attended close.

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
