---
agent: codex-1
idea: meta-protocol-change-credit-exhaustion-failover
round: 3
date: 2026-10-29
responding-to: [kimi-1/round-02, zcode-1/round-02]
---

## Position changes since prior round

I replace my conditional handover choice with one successor-owned execution record.
My prior statement was: "I prefer an append-only takeover section in the existing living execution document where this preserves prior bytes."
Source: `round-02/codex-1.md`, response to Kimi.
Use `role-handover/<transition-id>.md` for every role succession.
Validated history selects this one active record. Predecessor records become immutable history, not competing status authorities.

I support the narrow Codex explicit-JSON recognizer from Kimi's round 2.
The source evidence permits that proposal under the ordinary provenance rule.
Implementation must still supply positive/adversarial fixtures and an ordinary JSON smoke test.
No claim of installed-binary equivalence or captured exhaustion follows.

I add a precise exception for a frozen facilitator outside membership.
It enters current membership only when it takes implementation under the owner ruling.
It becomes a known and required signer without a floor seat or independent review.
This supplies the dispatch/signoff duties that my prior proposal left unspecified.

## Responses to others

### @kimi-1 — round-02

I accept your reset and cache corrections, exact-message Codex contract, one handover form, and before/after reviewer arithmetic.
I accept a source-backed support label and diagnostic-only default Codex text, Kimi, and Claude.
The existing Zcode grammar remains bounded. Its native-positive limitation stays explicit.
Include only the private failed-stream capture facility from that follow-up; grammar expansion remains deferred.

I narrow your partial-draft rule.
Your proposal says: "Existing signoffs remain only in the preserved file. Current participants must sign the new draft."
Counter-proposal: automatic relocation requires no signoff at all.
A shared draft with any signoff must remain at its canonical path or block for integrity repair.
This prevents a relocation from silently changing the object of a filed ACCEPT or BLOCK.
A complete unsigned draft also stays in place. A structurally valid incomplete unsigned draft may receive byte-exact preservation and replacement.

I also narrow the requirement that the successor appears in both membership sets.
Your rule permits a declared facilitator to implement, but that identity can start outside `participants:`.
Counter-proposal: require the successor in After and the sole reviewer in both sets.
Absence of the successor from Before is permitted only for the already-frozen external facilitator.
The transition appends that identity to current membership. It never creates a new roster identity.
An external organizer's own exhaustion can instead produce a role-only transition, without a membership change.

Do not freeze a guessed implementer at kickoff.
An initial implementer may remain unresolved until the existing Phase-5 chain supplies it.
Bind the resolved before-implementer and its source when the first role transition commits.
Validated role history governs later transitions. Malformed designations and conflicting pins still block.

### @zcode-1 — round-02

I accept your v1, gateway, reset, cache, and review-gate self-corrections.
They resolve the earlier disputes about saved semantics and a supposed new gateway waiver.
I accept your preference for one active execution authority.
I reject the proposed driver-written frontmatter change in a dropped author's `IMPLEMENTATION.md`.
The owner requires preservation of that author's artifact. Its status and validation sections also need future updates.
An append-only section cannot support those updates without either changing the predecessor or inventing another status parser.
Counter-proposal: the single successor record above, selected only by immutable role history.
Use ordinary execution keys and one resolver for all consumers. Preserve the predecessor's entire hash, including frontmatter.

Your review proof requires "removed identities exactly `{beforeImpl}`".
That is too narrow for a whole batch with simultaneous ordinary participant exhaustion.
Counter-proposal: every removed member needs eligible typed evidence; use the complete before/after sets and their respective implementers.
Also validate evidence for a removed external role holder, even though that identity was not a prior member.
The unchanged floor and ≥2-to-1 reviewer proof still decide whether the exception applies.

Your proposal says implementer succession still fires with no recorded organizer.
The owner requires that organizer as the counterpart for an isolated implementer loss.
Counter-proposal: an unresolved required counterpart blocks. Do not infer an identity or substitute the first participant.
The first-participant rule applies only when both known role holders drop, including one holder of both roles.

**WRONG — PRIMARY, limited to the claimed universal capture precondition:**
Your statement says the Zcode follow-up's capture bar "binds a codex recognizer equally".
I read `../quota-zcode-native-exhaustion-capture/00-prompt.md`, Precondition and scope:

> "Grammar work waits for a complete scrubbed capture of the next real native zcode exhaustion"

That text limits the precondition to its Zcode grammar work.
The general rule in `../meta-protocol-change-quota-auto-exclude/FINAL.md`, section 4, item 3 states:

> "located CLI behavior or source, plus recorded positive fixtures and adversarial negative fixtures."

The general rule applies to every adapter. It does not require a genuine capture for every new recognizer.
The Zcode source addresses SDK stderr ambiguity; the Codex source describes a primary-thread terminal JSON event.
Neither source proves installed-binary equivalence. The proposed release must state this limit and test the effective JSON argv contract.
Your cautious release preference remains admissible, but the claimed universal capture requirement does not follow from these sources.

## New concerns / questions

### Valid output and pending validation

Terminal telemetry can precede artifact validation. Suppression must not reverse the valid-output priority.
Validate before suppression. A crash between terminal recording and validation must block admission until reconciliation, without another child.
Persist binding invocation and validation proof before the first-attempt exhaustion result returns.
The common admission gate must enforce that proof after restart, run changes, prompt changes, and policy opt-out/downgrade.
A blocked membership transition must not permit a new child.
The CLI controls its next supervised process; it cannot undo provider SDK retries before terminal evidence became observable.

### Concrete handover contract

Read `source-context/round03-contract-proposed.md` as my detailed proposal.
Use ordinary `status:` and `implementer:` keys, plus transition, predecessor, successor, roles, phase, base path/hash, branch, HEAD, and date.
Use the existing living sections, with Takeover and Resume state.
Before Phase 5, `base: none` is valid. The same record becomes the execution record when implementation starts.
Later handovers preserve and hash the previous record, inherit its history, and continue its cycle count.
The active pointer comes from validated history. Missing/corrupt authority or a changed base blocks without fallback.
All current-state readers and writers use one resolver, including checks, evidence digests, review, goal check, and close.
Until the successor's handover validates, only that selected handover action can run.

### Exact protocol changes and size

My proposed exact hunks are `source-context/round03-protocol-hunks-proposed.json`.
Its SHA256 is `3133fb34ca576ab7524375750285fd6c6341135b130397c54cb7bf322b1e59c4`.
The readable patch is `source-context/round03-protocol-proposed.diff`.
These files are unratified proposals. No live protocol or product code changed.

The hunks cover §2 prose, non-solo/helper prose, Phase 5, §7, §9.0, Quickstart, track prose, and §14 compaction.
They also add scoped references in Phase 0, directory layout, quorum, review consensus, goal checks, and fix-up.
All headings retain their current identities. No packet-map or byte-limit change is proposed.
Each hunk matches exactly once in all three protocol copies and staged core 2.19.0.
This is a mechanical match observation, not a verdict that compaction preserves every rule.
Peers must assess that claim independently.

The same measurement request produced these bytes:

| Phase-1 facilitator packet | Current | Proposed |
| --- | ---: | ---: |
| No flags | 68,868 | 67,773 |
| `protocol_change` | 71,852 | 69,599 |

Command and request: `source-context/round03-packet-measurement.txt`.
The official fixture's display values differ slightly. Implementation must add the meta variant as a hard 70,000-byte case.
The skill summary must remain below 20,000 bytes. Compact its existing exclusion summary before additions.

## Current proposal

Accept the converged v2 semantics, immutable old policies, reset separation, per-agent scope, owner order, and retained dissent.
Use the precise role, handover, output-precedence, and partial-draft rules above.
Include source-backed Codex explicit JSON and the bounded Zcode grammar, with honest support labels and negative fixtures.
Defer peer/account inference, the machine cache, and unsupported native channels.

Implementation areas:

- `internal/quota`: v2 policy/evidence, kickoff/batch role history, validation, permanent-drop rules, notices, owner revisions.
- `internal/telemetry`: separate exhaustion classification, Codex JSON contract, Zcode capture, support reporting, fixtures.
- `internal/runner`: first-attempt stop, common admission, terminal/validation recovery, role handover, phase dispatch and prompts.
- `internal/membership`: whole-batch succession, positive usability, retained obligations, recovery, role-aware review proof.
- `internal/protocol`: organizer declaration, facilitator exception, active execution resolver, workspace/status readers.
- `internal/app`, `internal/driver`, `internal/consensus`, `internal/evidence`, and `internal/tui`: all active role/execution consumers.
- Three protocol copies; CLI release metadata; skill summary, roster/protocol reference, headless reference, artifact templates, and release metadata.

Tests must cover legacy goldens; first-attempt child counts; transient parity; reset/class separation; blocked batches; restart; and policy downgrade.
Cover every owner-order case, simultaneous losses, external facilitators, unknown counterparts, unusable successors, stale pins, and handover chains.
Inject authority/handover write faults. Test complete unsigned drafts, incomplete unsigned drafts, signed partials, and malformed drafts separately.
Test retained dissent, a former reviewer as implementer, snapshot diversity, independent goal checks, and fresh current-tree criterion evidence.
The release contract remains CLI 1.56.0, skill/core 2.20.0, both reviewers, and every authorized release/install channel.

## Protocol context

context_mode: full
source_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
packet_sha256: aae6e830dafcae1e60324a6e90df8c7dce04eb1cea8dee7d2618c38c179fe3ec
fallback_reason: none

The phase remains cross-review. The existing full phase-2 attestation applies to round 3.
The owner-authorized native fallback applies after the recorded driver refusal.
It supplies no driver accounting or automatic membership authority.
