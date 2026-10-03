---
idea: meta-protocol-change-quota-auto-exclude
status: in-progress
implementer: codex-1
started: 2026-10-04
branch: quota-auto-exclude
head-commit: c6f8b49
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Implementation has started; no code has changed yet. Both stages of FINAL are required: kickoff and
mid-idea quota auto-exclusion. The owner ratified D1–D5 and authorized implementation, separate claude-1
review, attended close, and release. The handoff is recorded in 00-prompt.md and commit c6f8b49.

Role concentration (§15.5): codex-1 organizes and implements; the separate claude-1 process reviews,
issues independent criterion verdicts, and owns its review and signoff files. codex-1 never reviews itself.

Protocol attestation: context_mode=full;
source_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388;
packet_sha256=b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388;
fallback_reason absent. The complete phase-5 body was read. Transport: owner-authorized local canonical
files on the existing two quota-auto-exclude branches; no development PRs.

## Implementation plan / checklist

- [ ] Establish native terminal-error provenance from installed CLI source/behavior. Record a support
  table and positive/adversarial fixtures. claude/text remains diagnostic-only. Recover incident 6 if
  reachable; otherwise retain the conditional record. Do not infer provenance from adapter names.
- [ ] Stage 1: implement a typed shared classifier in internal/telemetry/quota.go, strict reset parsing,
  scrubbed evidence, native-only recognition, and valid-artifact/later-success precedence. Wire terminal
  results and preflight; align bare-503 provider gates.
- [ ] Stage 1: add presence-aware layered quota_auto_exclude policy, freeze resolved policy/scope at
  kickoff, evaluate the entire readiness batch with fixed floor 2, apply protected-role gates and the
  C1 filter before initial prompt, creation event, manifest and dispatch. Standalone preflight reports only.
- [ ] Stage 2: add idea-scoped serialization only for recorded mid-idea-enabled policies, immutable
  membership history and one checked-durable transition per batch. Reconcile projections before all
  dispatch, signoff evaluation and close; truncated/contradictory history blocks. Keep legacy/off and
  kickoff-only semantics and never widen scope on upgrade.
- [ ] Stage 2: rebind current membership consumers and historical known signers, preserving filed vetoes,
  disputes and findings; reject excluded signoff appends. Re-evaluate role, reviewer, diversity,
  strict-gate and independent goal-check requirements. Preserve failed invocations and incomplete files;
  require a new survivor-validated terminal round evaluation before resuming.
- [ ] Deliver one deduplicated notice or blocking escalation, complete automatic markers and reset hints,
  read-only pending-state reporting, and current membership agreement in status/wait/organizer brief.
- [ ] Apply FINAL §13.1 protocol hunks identically across the three copies and add the changelog entry;
  update skill SKILL.md and ROSTER_AND_PROTOCOL.md. Check packet rendering at phases 0, 5 and 8.
- [ ] Run focused adversarial, batch-permutation, kickoff, replay/fault, serialization and history tests
  corresponding to AC2–AC18; then gofmt, go build ./..., go vet ./..., the full go test ./..., drift guard,
  three-copy comparison, and the skill's installer/lean-organizer checks (plus its required suite).
- [ ] Publish current-tree test evidence mapped to AC1–AC21 and open separate claude-1 review. Process
  binding findings through review consensus and at most five fix-up cycles; no self-issued verdicts.
- [ ] Obtain both review-consensus signoffs and write the blocking owner close note. Only after the owner
  confirms: mark complete, merge, release the current next versions and independently verify channels.

## Deviations from FINAL.md

The owner's ratification transfers the organizer role from claude-1 to codex-1, superseding FINAL §12's
organizer assignment. No implementation-scope deviation is planned. Exact transition storage, field
names, provenance support and redaction follow FINAL §13.8's implementation choices and will be logged.

## Notes for reviewers

Review FINAL and AC1–AC21 against the complete implementation diff. Attempt counterexamples, particularly
native-vs-content error provenance, saved policy scope, lost vetoes, stale captured membership, committed
but unreconciled transitions, and overlapping runs. Both stages are mandatory. The single-reviewer
configuration does not waive the attended close or permit unattended completion.

## Progress

- 2026-10-03 23:00Z: owner handoff committed; FINAL and controlling brief read in full; full phase-5
  protocol context read. Plan recorded before code. Implementation pending.

## Decision Log

- 2026-10-04 · codex-1: keep the existing worktrees and local-file transport override; no worktree pruning,
  declarations or legacy-record migration. D6 is a separate follow-up, outside this implementation.
- 2026-10-04 · codex-1: driver-first attempt `parley continue --auto <slug>` stopped at FINAL because
  auto_implement is false and the per-idea transport override exists only in prose. The normal manual
  Phase 5 path is authorized; use a separate configured codex-1 implementation process to bound context.
  The organizer retains phase transitions and release. Driver use will resume at review boundaries.

## Surprises & Discoveries

- Installed CLI 1.50.0, skill installer and runtime copies 2.14.0. Project source-role metadata is stale
  (2.12.0); dry-run sync reports only a metadata refresh. Do not replace the live source with the packaged
  protocol. Metadata refresh belongs to the authorized release work.
- The graphify query covers design documents, not the code graph; source inspection remains necessary.
- Shared OpenViking scoped recall is available and returned prior release/Windows notes, not this idea's
  current implementation. Current local owner direction and artifacts govern.

## Validation evidence

Pending implementation and independent review. No acceptance criterion is claimed complete yet.

## Outcomes & Retrospective

Pending both stages, review and attended close.
