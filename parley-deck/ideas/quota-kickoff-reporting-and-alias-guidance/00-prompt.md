---
idea: quota-kickoff-reporting-and-alias-guidance
author: codex-1
facilitator: claude-1
participants: [codex-1, kimi-1, zcode-1]
created: 2026-10-08
track: deliberation
auto_implement: false
status: candidate
---

## Problem / idea

Carry the accepted/deferred R8-MINOR-1, R8-NIT-1 and R8-NIT-2 findings from
`meta-protocol-change-quota-auto-exclude`. The binding round08-answer relay note
accepts their release limitations and the AC5 escalation-detail / AC15 kickoff
crash-window exceptions. None is fixed or withdrawn by the parent release.

## Scope for later design

- Preserve fail-closed kickoff floor/role gates while reliably recording an
  escalation when the inbox is absent and displaying candidate/arithmetic details
  when it cannot be written. Reproduce present, missing and read-only inbox cases.
- Define and verify kickoff notice crash/recovery semantics between kickoff record,
  event, manifest and publication, including the at-most-once versus eventual-attempt
  boundary. The parent finding is source-derived; inject crashes before claiming a fix.
- Make manual-revision alias refusals contextual. Plain participant edits and
  confirmed exclusions on aliased decks require a physical path even with policy off;
  all signers and driving are blocked until the pending edit can be imported. Do not
  advise disabling a policy already off. Consider preflight ordering without weakening
  physical lease identity or adding unsupported aliases.

## Acceptance direction

An independent reviewer must reproduce the original cases on shared and local
filesystems, verify fixes through real CLI paths and retain quorum, protected-role,
immutable history and serialization gates. Changes to normative notice semantics
require their own protocol-change treatment. Windows directory durability remains
on `windows-portability`; native zcode capture remains a different follow-up.

## Sources and status

- Parent `review/round-08/claude-1.md`, frozen CLI 2705a1e / skill 99b3f3f.
- Parent `source-context/codex-1-round08-disclosure-proposal.md`.
- `../../inbox/user-to-codex-1_meta-protocol-change-quota-auto-exclude_round08-answer.md`.

Backlog only, not cycle 6. No readiness probe, participant process or active
run has started. The future facilitator and explicit participant list are the
owner's global defaults; verify live roster/model settings at that future launch.
They do not change the parent idea's two-participant quorum.
