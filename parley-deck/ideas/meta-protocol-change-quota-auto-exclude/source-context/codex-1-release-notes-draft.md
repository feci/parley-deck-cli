# Release notes draft — quota auto-exclusion

Producer draft for independent review. Publication waits for the attended close.

New ideas can continue after a participant exhausts an allowance when the complete failed batch leaves at least two usable non-facilitator participants. The change applies to the current idea's quorum, with one durable transition and an owner notice; machine and deck roster files are unchanged. Existing ideas retain their recorded policy and scope.

The supported zcode recognizer implements the owner's explicit stderr exception: a failed process with no completed artifact or later success, agreeing provider error records containing HTTP 429 exhaustion and a machine reset at least 60 minutes away, and a terminal turn-failure line. All provider error records must agree; missing or contradictory resets and quoted assistant/tool text alone do not qualify. A timezone-free display clock is supplemental only when it agrees with a machine reset in the same record under the bounded offset rule. The owner accepts the remaining possibility that a subagent's quota error accompanies an unrelated root failure. Claude text, Codex and Kimi remain diagnostic-only for automatic exclusion without established native evidence.

Kickoff and mid-idea application are both required for this release. Mid-idea changes preserve historical signoffs, vetoes, disputes, findings, failed invocations and incomplete files. Required signers and dispatch consumers use current membership. A designated or pinned implementer, or a participant that has started a consensus/FINAL draft, cannot be automatically removed. Pending transitions block further mutation until reconciled, while status/wait/organizer views report them without repairing state.

Two fixes apply even when the policy is off: `parley run --yes` no longer retains or dispatches an excluded identity, and a bare preflight HTTP 503 remains a provider gate that `--yes` cannot exclude. There is no quota polling, retry worker, automatic same-idea rejoin or organizer failover.

This idea has one non-facilitator, so the rule cannot reduce its quorum. It retains the separate claude-1 review and owner-confirmed close. Windows remains experimental and CLI winget publication stays held. The legacy run-accounting blocker is a separate owner-requested follow-up (D6).
