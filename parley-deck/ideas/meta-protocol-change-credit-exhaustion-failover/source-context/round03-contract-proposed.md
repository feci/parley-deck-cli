# Focused contract proposal

This is a design draft. It is not a ratified plan.

## Roles and membership

Freeze original participant order and the organizer identity in the new v2 kickoff record.
The optional `organizer:` identifies an ordinary participant organizer.
Otherwise, a declared `facilitator:` identifies the organizer. Conflicting fields block.
An absent organizer remains unresolved; it is never inferred from `author: user`.
Preserve the declared facilitator identity as a permanent floor exclusion.

The implementer can be unresolved at kickoff.
Before the first role transition, resolve it through the existing pin/designation/default/fallback chain.
An absent future Phase-5 pin is normal. A malformed designation or conflicting pin still blocks.
Bind that resolved identity and its source in the immutable transition.
After the first transition, validated role history has priority over the retired pin and designation.
Neither a binary upgrade nor a changed standing default rewrites that history.

Apply the owner's order to the complete settled failure set.
Only proven exhaustion can remove a protected role under v2.
Ordinary transient failures still use v1's two-attempt path and protected-role gates.
The chosen counterpart or first remaining original participant needs positive usability.
An unusable chosen successor blocks; the selector never skips that identity.

A frozen pure facilitator can be outside `participants:`.
Its exhaustion permits a role-only transition: membership stays the same, and the implementer takes the organizer role.
If that pure facilitator takes implementation, append it to current membership in the same typed transition.
Only that already-frozen identity can enter through this exception. No new roster identity enters.
It becomes a known and required signer, but supplies no floor seat or independent review.
Before counts use original membership and the before-implementer.
After counts use current membership and the after-implementer.
The sole reviewer must belong to both sets. The successor must belong to After.
An external successor may be absent from Before only through this frozen-facilitator exception.
This replaces the overly strict requirement that every successor belongs to both sets.

## Handover and artifact preservation

Use one successor-owned `role-handover/<transition-id>.md` for each role transition.
It uses the ordinary execution keys `status:` and `implementer:`.
It also records transition, predecessor, successor, roles, phase, base path, base hash, branch, HEAD, and date.
Use existing living sections: Progress, Decision Log, Validation evidence, and fix-up history.
Add Takeover and Resume state. The successor records the next action and canonical paths.
Before Phase 5, `base: none` is valid. When implementation starts, this same file becomes its execution record.
After Phase 5, preserve the predecessor record byte-for-byte and bind its exact hash.
The handover carries current execution values and inherits earlier history through its base chain.

One resolver selects the current execution record from validated role history.
Every current-state reader and writer uses it: driver, runner, checks, evidence, status, wait, brief, review, goal check, and close.
Missing, corrupt, contradictory, or changed bases block. A directory listing or newest filename never selects authority.
History commits before successor dispatch. Until the handover validates, only its selected author's handover action is admitted.
One notice names the takeover and pending action. Repeated recovery must not duplicate that notice.

A complete consensus or FINAL draft remains unchanged, even when signatures are pending.
A missing draft uses the current eligible selector.
An incomplete shared draft may move byte-exactly only when its frontmatter is valid and no signoff exists.
Record original path, author, hash, and archive path in the transition.
Then the selected successor writes a new draft at the canonical path.
A draft with any signoff, malformed identity, malformed frontmatter, or truncated syntax blocks for integrity repair.
This restriction prevents relocation from dissolving a prior signoff or BLOCK.
Participant rounds and reviews never move. Complete FINAL and closed implementation records stay frozen.

## Terminal evidence and dispatch

Validate the artifact before exhaustion suppression.
Valid completed output wins, including a BLOCK or FAIL verdict.
A terminal-to-validation crash creates a pending admission gate; reconciliation validates existing output without another child.
Persist binding invocation and validation evidence before the no-retry result returns.
The common admission gate suppresses the identity across later steps, runs, prompt changes, restart, and policy changes.
A blocked membership batch never restores dispatch permission.
The CLI prevents its next supervised child. It cannot undo a provider SDK retry before terminal exhaustion became observable.

Keep scope per agent and per idea. Defer peer inference and the machine cache.
Retain the legacy Zcode framing grammar; v2 changes class/reset semantics within that grammar.
Short or invalid resets cannot negate independently valid exhaustion proof.
Malformed framing or contradictory class still rejects. No new gateway waiver exists.
Codex JSON support requires explicit effective JSON argv, exact source-mapped messages, fixtures, and an ordinary JSON smoke test.
Default Codex text, Kimi, and Claude remain diagnostic-only.
Never induce a live credit failure for evidence.
