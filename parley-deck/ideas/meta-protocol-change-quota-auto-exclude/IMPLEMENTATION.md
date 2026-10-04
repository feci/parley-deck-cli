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

The first stage-1 invocation timed out after 1800 seconds; its code is preserved and partial. No adapter yet authorizes automatic exclusion. Stage 2 is pending. Both stages of FINAL are required: kickoff and
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

- The three protocol copies retain their pre-existing project/bootstrap headers and §2 tables.
  Identical normative hunks were applied, and the embedded drift guard passes; the skill is NOT
  whole-file byte-identical to the live deck as AC1 literally says. This deviation is open for review
  and owner close; no test or allowlist was weakened.
- AC2 is UNMET in the current prototype: all four adapters are diagnostic-only, and the private
  semantic refinement has no production recognizer caller. A support table with installed-source
  locators is in internal/telemetry/testdata/quota/README.md. The separate claude-1 focused review (review/round-01/claude-1.md, MAJOR-1) concludes that the configured zcode text path cannot establish the required terminal provenance. This is not a
  substitute for AC2, not stage-1 completion, and not permission to release.


The owner's ratification transfers the organizer role from claude-1 to codex-1, superseding FINAL §12's
organizer assignment. No implementation-scope deviation is planned. Exact transition storage, field
names, provenance support and redaction follow FINAL §13.8's implementation choices and will be logged.

## Notes for reviewers

Review FINAL and AC1–AC21 against the complete implementation diff. Attempt counterexamples, particularly
native-vs-content error provenance, saved policy scope, lost vetoes, stale captured membership, committed
but unreconciled transitions, and overlapping runs. Both stages are mandatory. The single-reviewer
configuration does not waive the attended close or permit unattended completion.

## Progress

- 2026-10-04: host `go test ./... -timeout 45m` exited 0 on the current stage-1 code.
  Every package passed; app took 496.176 s and trajectory 548.381 s. The prior host
  command hit the default 10-minute package ceiling in trajectory, with no assertion failure.
  The longer rerun is the current full-suite evidence; log:
  `.parley-runtime/quota-implementation/stage1-full-host-45m.log`. This does not establish
  the missing AC2 recognizer or any not-yet-implemented stage-2 acceptance criterion.

- 2026-10-04: a further read-only installed-source search located native JSONL logs under
  `~/.zcode/cli/log` (singular), including a real `turn.failed` with provider attribution and a
  53-hour quota message. The scrubbed record and source snippets are in
  source-context/provenance-review/codex-1-structured-channel-evidence.md. This differs from AC2's
  exact recorded 49-hour positive. A separate claude-1 supplemental feedback process is checking
  unchanged-invocation capture, complete reset preservation, binding and adversarial negatives,
  plus whether MAJOR-1/MAJOR-2 need owner scope interpretation. No AC verdict is asserted and
  no classifier is enabled. The full Go suite is running on the host.

- 2026-10-03 23:33Z: claude-1 focused reviewer exited 0 after 553.4 seconds. Its own artifact is
  structurally valid (parley wait: 1/1 filed-and-valid, 15656 bytes); it explicitly is NOT a full-scope
  review. MAJOR-1 requires an owner scope decision for AC2; MAJOR-2 identifies the reset parser's
  rejection of the recorded positive body; MINOR-1 asks for the decisive provenance limitations in
  the support table. The support table now cites the review's terminal-binding/framing conclusions.
- 2026-10-03 23:36Z: host `go build ./...`, `go vet ./...`, gofmt listing, diff whitespace check and
  `go test ./internal/app -run '^TestQuota' -count=1` passed. Focused quota/drift cases in runcontrol,
  quota, telemetry, config and protocol passed. runstate/runmanifest matched no quota-named tests,
  so those no-test lines are not coverage evidence. Full current-tree Go suite remains owed.

- 2026-10-03 23:33Z: stage-1 child stopped at its 1800-second limit (exit -15), not a provider failure.
  Preserved shared quota policy/evidence/batch primitives; initial immutable records; presence-aware
  config; kickoff C1 filtering; preflight report-only changes; bare-503 alignment; telemetry capture;
  status/wait/organizer surfaces; associated tests. No independent acceptance verdict is claimed.
  Stage 1 remains partial because AC2 has no supported recognizer and end-to-end coverage is incomplete.
- 2026-10-03 23:33Z: organizer/implementer documentation changes are committed (CLI f1a7f80;
  skill dfad28e/c0d7f58). All three normative rule texts agree; bootstrap-zone discrepancy remains open.
  Phase 0/5/8 packets contain the rule/cross-references; packet check reports ok=true (69 blocks).
  Skill required checks: 55/55. Full npm test: 399 Node and 54 Python tests pass; all manifests match.
- Next concrete work: obtain claude-1 supplemental review of the located native structured log; resolve
  or escalate its AC2 conclusions; finish stage-1 integration/tests; implement full mid-idea transitions,
  serialization, replay, captured-consumer rebinding, known/required signers, preserved findings and gates.
  Keep status in-progress and the protocol stage-2 paragraph not-yet-in-force until delivered.

- 2026-10-03 23:07Z: scoped implementer read FINAL and this plan in full, ratification, handoff, and
  freshly rendered full phase-5 deliberation protocol body with `--flag protocol_change`. Attestation:
  `context_mode=full`, source and packet SHA256
  `b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388`, fallback_reason absent.
  Shared memory recall returned `MCP tool call requires approval, but approval policy is never`; local
  sources govern. No participants launched. Stage 2 and protocol/skill hunks remain for the next invocation.

- 2026-10-03 23:00Z: owner handoff committed; FINAL and controlling brief read in full; full phase-5
  protocol context read. Plan recorded before code. Implementation pending.

## Decision Log

- 2026-10-04 · codex-1: accept the need for an owner decision on claude-1 MAJOR-1. The supported path
  cannot be invented from an adapter name, a diagnostic dump or synthetic policy tests. No Stage-2
  process is launched while deciding whether the required terminal-attributed transport is in scope.
  This is an early implementation blocker, not the attended close and not stage-1-only authorization.
- 2026-10-04 · codex-1: retain MAJOR-2 as an open agreed issue for the next implementation pass.
  Its suggested interpretation of unzoned human-readable text alongside machine reset fields must be
  reconciled with FINAL's strict unparseable-reset rule and independently re-reviewed, never silently
  changed to unknown. No reset correction or false positive fixture is claimed yet.

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

No acceptance criterion is independently claimed complete. Producer checks so far:

- `go test ./internal/quota ./internal/telemetry ./internal/config -count=1` on host: pass.
- `go test ./internal/runner -run 'TestVerifierRecoveryRefusedEvidenceMutationFailsHandles|TestQuota' -count=1`
  on host: pass, including the fixture denied by the subprocess sandbox.
- The subprocess full suite encountered sandbox-denied cache/launch fixtures. The later host full
  suite (`go test ./... -timeout 45m`) passed every package at this stage-1 checkpoint. The first
  host run reached the default 10-minute trajectory timeout; the longer rerun exited 0.
  Full checks must run again after further implementation changes.
- Protocol drift guard and packet checks pass; skill checks above pass. AC1's whole-file comparison is
  still unmet, AC2 is unmet, stage-2 ACs are not implemented, and attended close/signoffs remain owed.

Raw producer evidence is under `.parley-runtime/quota-implementation/` and `/tmp/quota-skill-full.log`.
The focused reviewer owns `review/round-01/claude-1.md`; it is not a full-scope closing review.

## Outcomes & Retrospective

Pending both stages, review and attended close.
