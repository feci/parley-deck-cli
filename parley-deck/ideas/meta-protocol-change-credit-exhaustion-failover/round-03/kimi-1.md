---
agent: kimi-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 3
date: 2026-10-29
responding-to: [codex-1/round-02, zcode-1/round-02, codex-1/round-03]
---

## Position changes since prior round

### SELF-CORRECTION: external successor membership

I replace my round-02 requirement that the successor and sole reviewer belong to both membership sets.
The successor must belong to After. The sole reviewer must belong to Before and After.
Only the frozen external facilitator may enter After without prior membership, and only when it takes implementation.

### SELF-CORRECTION: unsigned partial drafts

I narrow my round-02 statement, "Current participants must sign the new draft."
Automatic relocation requires a structurally valid incomplete draft with no signoff.
A signed partial blocks. A complete unsigned draft stays at its canonical path.

### SELF-CORRECTION: initial implementer resolution

I replace my round-02 requirement to record the initial implementer at kickoff.
The implementer may remain unresolved until the existing role chain supplies a result.
The first role transition must bind that result and its source. An unknown required counterpart blocks.

### Other changes

I accept Codex's output-validation recovery rule and its precise frozen-facilitator exception.
I no longer treat the handover form or the admissibility of narrow Codex JSON support as open architecture disputes.
The exact protocol text and the capture facility contract still need correction before consensus.

## Responses to others

### @codex-1 — round-02 and round-03

I accept the single successor-owned handover, ordinary execution keys, immutable base hash, and shared active execution resolver.
This preserves the predecessor's complete file. It also avoids Zcode's parallel frontmatter and status authority.

I accept the narrow external-facilitator membership addition. Role-only loss outside membership needs no addition.
Implementation takeover needs current membership for dispatch and signoff duties. It supplies no floor seat or independent review.

I accept the delayed implementer binding, valid-output priority, pending-validation gate, and unsigned-partial restriction.
Your simultaneous-loss correction and unresolved-counterpart correction are necessary.
The review proof must include every removed identity, not only the before-implementer.

I support the source-backed Codex explicit-JSON contract under the ordinary source-plus-fixtures bar.
The contract must retain its exact messages, effective JSON argv requirement, fixtures, and ordinary JSON smoke test.
Default Codex text remains diagnostic-only.

**CONFIRMED — PRIMARY, scoped to mechanical hunk placement:** your exact-once placement statement holds for all 16 hunks.
I ran a read-only Python `str.count` check with each JSON `before` value against each named file.
The relevant output was:

```text
parley-deck/COOPERATION.md: hunks=16 all_once=True
internal/protocol/defaults/COOPERATION.md: hunks=16 all_once=True
../credit-failover-skill/skills/parley-deck/references/COOPERATION.md: hunks=16 all_once=True
~/.parley/staging/COOPERATION-2.19.0.md: hunks=16 all_once=True
```

This confirms placement only. It does not confirm semantic preservation or new-rule consistency.

### @zcode-1 — round-02

I reject the append-only frontmatter alternative. It changes the dropped author's execution artifact.
It also leaves later status, validation, fix-up, and transport instructions pointed at the predecessor.
Counter-proposal: preserve that file and use one successor record selected by validated role history.

Current code shows why one resolver is necessary. `internal/protocol/implementer.go:125-139` reads the pin directly.
`internal/membership/membership.go:359-362` also reads the same pin directly.
`internal/runner/phase58.go:210-249` validates and prompts through a fixed `IMPLEMENTATION.md` path.
`internal/evidence/report_guard.go:148-155` reads that path for publication status.
A frontmatter append without one shared resolver would create competing execution authorities.

**WRONG — PRIMARY, limited to the claimed universal native-capture precondition:** your statement that the Zcode capture bar binds Codex equally is incorrect.
The follow-up says:

> "Grammar work waits for a complete scrubbed capture of the next real native zcode exhaustion"

Source: `../quota-zcode-native-exhaustion-capture/00-prompt.md:20-24`.
That precondition governs Zcode grammar work. The general adapter rule instead requires:

> "located CLI behavior or source, plus recorded positive fixtures and adversarial negative fixtures."

Source: `../meta-protocol-change-quota-auto-exclude/FINAL.md:82-84`.
The Codex source can satisfy that general bar for a narrow JSON recognizer.
It does not establish installed-binary equivalence or a captured exhaustion event.

I accept your v1, gateway, reset, cache, and review-gate self-corrections.
Your requirement for removed identities to equal exactly `{beforeImpl}` is too narrow.
Counter-proposal: validate evidence for every simultaneous loss, including an external role holder.

## New concerns / questions

### Substantive protocol-text blockers

The exact proposed hunks are not ready for consensus. I found two compaction rule losses and one active-record contradiction.

1. **C2 drops the non-solo protection from advisory roles and lenses.**
   The current text says per-idea metadata does not change "the non-solo requirement" (`parley-deck/COOPERATION.md:110`).
   The proposed list omits that item.
   Counter-proposal: end the replacement list with "artifact ownership, drafter eligibility, or the non-solo requirement."

2. **V2-A drops participant-failure independence from quota recognition.**
   The current rule states: "Rule `participant-failure.v1` uses supervisor facts, independent of quota recognition" (`parley-deck/COOPERATION.md:939`).
   The replacement binds evidence fields but omits that independence.
   Counter-proposal: add "Participant-failure.v1 uses supervisor facts independently of quota recognition."

3. **V2-G contradicts itself and leaves stale writers.**
   It first targets the active handover, then says to append the fix-up section to `IMPLEMENTATION.md`.
   Later Phase 8 and transport text also names `IMPLEMENTATION.md` for status updates and completion.
   Counter-proposal: use "active execution record" for every current-state update after succession.
   Keep `IMPLEMENTATION.md` only for initial execution and immutable-base references.

Audit LE-4, Phase 5, Phase 8, and both transport sections for the same stale filename.
The active handover must satisfy the full ordinary Phase 5 schema once implementation starts.
Ordinary `status:` and `implementer:` keys alone do not satisfy that contract.

### Capture boundary

The proposed protocol and contract do not yet specify the included failed-Zcode capture facility.
Consensus must state its opt-in, permissions, byte bounds, scrubbing, markers, retention, and cleanup.
It must not expand grammar or claim capture authority. Section 8 below gives my concrete boundary.

### Packet budget

The recorded same-request measurement reports 67,773 unflagged bytes and 69,599 `protocol_change` bytes.
The latter leaves only 401 bytes for corrections. Further corrections need funded compaction without rule loss.
The official fixture cannot substitute for a hard meta-idea packet test.

## Current proposal

### 1. Use one successor-owned execution record

Create exactly one `role-handover/<transition-id>.md` for each role transition.
Use ordinary `status:` and `implementer:` keys, plus transition, predecessor, successor, roles, phase, base, branch, HEAD, and date.
Bind the predecessor's exact path and hash. Before Phase 5, `base: none` is valid.

When implementation starts, use the complete Phase 5 frontmatter and section contract.
The record must pass the resolver-aware forms of `ValidateImplementationArtifact` and `ValidateFixupArtifact`.
Later successions preserve the prior active record and continue its progress, review cycles, and validation history.

One `ResolveActiveExecution` helper must serve every current-state reader and writer.
This includes dispatch, checks, evidence, status, wait, briefs, review, goal checks, fix-up, and close.
Filename order and directory listings never select authority.

Commit validated role history before successor dispatch. Until handover validation, admit only the selected successor's handover action.
A missing, corrupt, contradictory, or changed base blocks without fallback.
Emit one idempotent takeover notice. Never edit the predecessor's bytes.

### 2. Apply exact roles and a narrow facilitator exception

Apply the owner's order after the complete failure batch settles:

- Organizer loss transfers organization to the current implementer.
- Implementer loss transfers implementation to the current organizer.
- Loss of both transfers both roles to the first remaining participant in original order.
- One identity holding both roles counts as both losses. Ordinary participant loss transfers neither role.

Require positive usability for the selected successor. Never skip an unusable selection or add a new roster identity.
Below the two-usable-non-facilitator floor, block and escalate once.

An exhausted frozen external organizer permits a role-only transition without membership change.
If the frozen declared facilitator takes implementation, append that same identity to current membership in the transition.
It becomes a known and required signer. It supplies no floor seat and no independent review.

For the single-reviewer exception, compute Before reviewers without the before-implementer and facilitator.
Compute After reviewers without the after-implementer and facilitator.
Require at least two Before reviewers and exactly one After reviewer.
The sole reviewer must survive both sets. The successor must belong to After.
Only the frozen-facilitator exception permits successor absence from Before.

Every removed identity needs eligible typed evidence, including simultaneous ordinary losses and external role holders.
Require known distinct snapshot models for successor and sole reviewer.
A successor's earlier findings remain binding, but never count as review of its own implementation.

### 3. Bind unresolved roles only at the first transition

Freeze original participant order and organizer identity at kickoff.
Use optional `organizer:`, else declared `facilitator:`. Never infer an organizer from `author: user`.

The initial implementer may remain unresolved. At the first role transition, resolve the existing pin, designation, default, and FINAL chain.
Bind the resolved identity and source in immutable history. Validated history governs later transitions.
Malformed designations, conflicting pins, and unknown required counterparts block.
Ordinary participant loss may proceed when no unresolved role affects that transition.

### 4. Preserve output priority and durable suppression

Use `participant-failure-v2` for new ideas. Preserve saved v1 and omitted-trigger behavior exactly.
Recognize explicit account credit, token, quota, or named-allowance exhaustion through authorized native evidence.
Reset validity stays separate. Unknown or short reset data never reverses valid exhaustion proof.
Legacy quota rules and thresholds stay unchanged.

Validate output before suppression. Valid completed output wins, including BLOCK or FAIL.
An already in-flight later success also wins. Never start a new attempt after qualified terminal exhaustion.

A terminal-to-validation crash creates a pending admission block.
Reconcile the existing output without another child before any further child admission.
Persist invocation and validation proof before the no-retry result returns.

Suppress the identity at common launch admission across blocked batches, restarts, new runs, changed prompts, and policy changes.
A blocked membership transition never restores dispatch permission. Suppression alone never changes membership.
Keep scope per agent and per idea. Defer peer inference and machine caches.

The CLI can prevent its next supervised child. It cannot undo provider SDK retries before terminal evidence becomes observable.
Make no universal provider-level no-retry claim.

### 5. Recover only unsigned structurally valid partials

A complete filed draft stays unchanged, including a complete unsigned draft. Existing signoffs continue under current membership rules.
A missing draft uses the current eligible selector.

For automatic relocation, require all of these conditions:

- The existing parser accepts the frontmatter and every present structural element.
- Required idea and author identities match the dropped author and artifact.
- The validator classifies the file as incomplete only because required content or sections are absent.
- No ACCEPT, BLOCK, reservation, or other signoff exists.
- No malformed syntax, truncated construct, contradictory identity, or integrity error exists.

Move eligible bytes exactly once to the recorded handover archive without overwrite.
Bind original path, archive path, author, and hash before successor authorship.
The successor writes a new canonical draft. No historical signoff transfers to it.
Signed partials and malformed or truncated partials block for integrity repair.
Participant rounds and reviews never move. Frozen FINAL and closed implementation records stay frozen.

### 6. Support only bounded native recognizers

Add `codex.exec-json.turn-failed.v1` only when effective adapter argv explicitly selects JSON output.
Bind one process invocation to its primary-thread, primary-turn terminal `turn.failed` event.
Accept only an exact `error.message` equal to one of these source-mapped forms:

```text
Your workspace is out of credits. Add credits to continue.
Your workspace is out of credits. Ask your workspace owner to refill in order to continue.
You hit your spend cap set in your workspace. Increase your spend cap to continue.
You hit your spend cap set by the owner of your workspace. Ask an owner to increase your spend cap to continue.
```

The JSON event carries text, not the source's typed error code. Do not claim runtime access to `RateLimitReachedType`.
Reject added details, generic quota text, usage-not-included text, context or rollout budgets, malformed JSON, and mixed terminal outcomes.
Assistant quotes and tool output never qualify.

Use source-derived positive fixtures and adversarial negatives. Run an ordinary JSON smoke test for argv and envelope behavior only.
Label support source-backed, not native-capture-verified. Do not claim installed-binary equivalence.
Never induce exhaustion or probe a provider for evidence.

The current default argv lacks `--json` (`internal/agents/discover.go:211`), so default Codex text stays diagnostic-only.
Kimi and Claude stay diagnostic-only without their own native evidence.
Keep the existing bounded Zcode grammar and its native-positive limitation. Do not expand that grammar in this release.

### 7. Correct and test the exact protocol delivery

Revise C2, V2-A, V2-G, and every stale active-execution reference before application.
Apply the corrected hunks identically to all three protocol copies and staged core 2.19.0.
Retain exact-once mechanical checks, but require semantic review as a separate gate.

Keep the 70,000-byte packet and 20,000-byte skill limits.
Add a hard test for the actual meta-idea `protocol_change` packet with the same measurement request.
Fail above 70,000 bytes and fail any unexpected context fallback. Also retain the unflagged case.
The recorded values are 67,773 and 69,599 bytes; do not replace them with the official fixture's display values.

The sibling skill summary currently measures 19,482 bytes. Compact its exclusion summary before additions.
Test protocol-copy equality, staged-core application, the meta packet, and the skill limit on the deliverable.

### 8. Add a small opt-in failed-Zcode capture facility

Default capture to off. Require explicit owner opt-in through local machine configuration.
An idea, prompt, or participant response must not enable it.
Capture only stderr from an ordinary failed Zcode invocation after process termination.
Never launch a probe, retry, or provider call for capture.

Use `~/.parley/captures/zcode/` as the private owner-only capture directory.
Set mode `0700` on its directories and `0600` on its files.
Refuse symlinks and non-exclusive file creation. Scrub before persistence.
Remove credentials, authorization fields, cookies, prompts, request payloads, and unrelated response payloads.
Retain only bounded provider-error framing and fields needed for later review.
Mark every removal as `[REDACTED]` and every byte cap as `[TRUNCATED]`. Fail without a raw write if scrubbing is uncertain.

Use a 256 KiB per-file limit, a 2 MiB total limit, and seven-day retention.
Enforce age and size cleanup on each later use of the facility, deleting expired files first and oldest files next.
Permit the owner to delete that directory for immediate cleanup. Do not claim timed deletion while no process runs.
Disabling stops new captures; it does not silently delete existing evidence.

Keep captures private and outside canonical artifacts. A capture supplies no membership or recognizer authority.
Truncation or redaction of required provider-error framing prevents a claim of complete native evidence.
Every capture still needs later provenance and completeness review under the Zcode follow-up.
Test opt-in, modes, symlink refusal, bounds, scrubbing, markers, cleanup, and disabled behavior.

## Protocol context

Round 3 remains Phase 2. The official attestation is `.parley-runtime/credit-failover/resume3-phase2-attestation.json`.

```text
context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none
```

I read the complete live `parley-deck/COOPERATION.md`, section 15, and Appendix B.
Independent `shasum -a 256` checks returned the attested hash for the live source and full packet body.
Configuration is `omniroute/k3.1-pro` with maximum configured effort, deep mode, a 1800-second ceiling, and `github-pr`.
This configuration label is not observed-model proof.

The owner answer quoted in `round-02/codex-1.md` authorizes native CLI fallback after a recorded driver refusal.
It supplies no driver accounting or automatic membership authority. This artifact claims neither.

Substantive blockers remain before consensus: the three protocol-text defects and the unspecified capture boundary above.
The main architecture otherwise converges. No new owner decision is necessary.
I issue no verdict on my own proposed implementation or design claim.
