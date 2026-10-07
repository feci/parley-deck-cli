# Proposed release disclosures after round-08

Status: proposal awaiting the owner's cycle-5-limit disposition, 2026-10-08.
No code fix, waiver or final consensus is inferred. The source remains CLI 2705a1e /
skill 99b3f3f. These paragraphs would be added to the release notes and mirrored in
the relevant changelog/documentation after an explicit owner acceptance.

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

Proposed follow-up scope: reliable kickoff floor/role escalation with visible
decision details when the inbox cannot be written; explicit kickoff notice
crash/recovery semantics; contextual physical-deck diagnostics for manual edits.
The follow-up would preserve fail-closed gates and record the current review's
exact probes. It is not a sixth fix-up cycle of this idea and has not been launched.

The existing native-positive AC2 owner waiver, Windows durable-operation failures,
CLI winget hold and D6 deferral would remain explicit. The three new findings would
be accepted/deferred, never described as fixed or withdrawn.
