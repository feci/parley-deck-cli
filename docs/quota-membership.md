# Recorded quota membership and recovery

New ideas record the resolved automatic-exclusion policy and initial membership in `quota-kickoff.json`.
The existing quota_auto_exclude knob now defaults new ideas to
`{"enabled":true,"scope":"kickoff-and-mid-idea","trigger":"participant-failure-v1"}`.
Omitted-trigger saved policies keep quota-only behavior and serialized hashes. Use
`parley run --quota-auto-exclude=false ...` for the new-idea opt-out. Resume and binary
upgrade never widen the trigger; explicit owner-bound revision is required.

Under participant-failure-v1, a dispatched non-protected step gets its original attempt
and one retry after five seconds at the same ceiling. Child/provider failures of any code,
crashes, timeouts and missing/invalid own output qualify after both failures. The durable
idea/agent/logical-step key prevents restart/input/run changes from minting a third attempt.
Valid BLOCK/disagreement wins. Control-plane/operator cancellation and shared-file integrity
failures do not authorize exclusion. Preserve invalid output privately; tampering stops.

Before an idea exists, readiness uses two attempts within the proposed batch. Replaying
that same batch with its stable identity cannot add an attempt. A separate new proposal
probes afresh; there is no persisted proposal/resume/abandon lifecycle across separate
pre-idea commands. After creation, the limit is per idea across driver restarts. Goal-check
and readiness validation recover their original invocation logs after a crash; changed
previously valid evidence blocks instead of launching a replacement.

The usable floor and every precommit gate remain: auto_implement 3→2 blocks with only
one independent reviewer. A permanently dropped ID cannot return to this idea, including
via opt-out/downgrade, owner/manual revision or catch-up; the next idea probes afresh.
Retained dissent needs an owner ruling quoted into the next artifact, or abandonment/v2.
All return/catch-up examples below apply only to IDs without a permanent dropout.

Automatic exclusions append immutable `quota-history/000001.json` revisions; a resume or a newer
binary never widens their policy. The CLI uses current membership from this history. Historical
participants, filed signoffs, retained findings and vetoes remain known. Returning a participant
does not erase an earlier exclusion or count as a fresh signoff.

## Explicit owner revision

Use the normal owner-answer workflow to obtain and commit the decision before running:

```sh
parley quota revise --dir WORKSPACE --idea IDEA --run EXISTING_RUN --request revision.json
```

The run must already have a current manifest bound to this idea. Stop writers first. The command
requires an open idea and serializes the revision on the deck filesystem. It does not commit files,
change a roster or provider, launch participants, close an idea, or waive any quorum or role gate.

The owner's answer must have `from: user` and the exact `idea:` in its frontmatter, under
`parley-deck/inbox/user-to-*.md` (including `inbox/archived/`). Its quoted body must include this
standalone directive, with the exact requested set and policy:

```text
Quota revision: {"participants":["codex-1","kimi-1","zcode-1"],"policy":{"enabled":true,"scope":"kickoff-and-mid-idea"}}
```

A request binds that verbatim quote to a full Git commit ID, the blob ID at that committed path,
and the SHA-256 of the entire answer's bytes. It does not accept a branch name or an uncommitted
answer. For example, replace every placeholder below with the reviewed evidence:

```json
{
  "participants": ["codex-1", "kimi-1", "zcode-1"],
  "policy": {"enabled": true, "scope": "kickoff-and-mid-idea"},
  "authority": {
    "path": "parley-deck/inbox/user-to-all_IDEA-return.md",
    "commit": "FULL_COMMIT_ID",
    "blob": "FULL_BLOB_ID",
    "sha256": "SHA256_OF_COMMITTED_ANSWER_BYTES",
    "quote": "Quota revision: {\"participants\":[\"codex-1\",\"kimi-1\",\"zcode-1\"],\"policy\":{\"enabled\":true,\"scope\":\"kickoff-and-mid-idea\"}}"
  }
}
```

Obtain the blob with `git rev-parse COMMIT:PATH`; hash `git show COMMIT:PATH` with SHA-256.
Use the answer's committed path even after permitted inbox archiving or deletion. The immutable
revision preserves the quote and object identifiers. A missing live inbox file is allowed.
At binding, a contradictory live or archived copy is rejected; a changed digest, wrong idea/author,
missing object or fabricated path is rejected at binding and on every later read.
Keep the committed Git object available for historical verification.

Only previously known identities qualify as re-inclusions. A new identity additionally needs
`"catchup":{"new-id":"round-01/new-id.md"}` and a valid late round-1 artifact. Its frontmatter
must declare `catch-up: true`, `read-priors: [round-01/other-id.md, ...]`, and `join-from: round-02`.
The initial command checks the referenced prior artifacts and captures the catch-up artifact in
immutable history. The owner directive must authorize the new current set. Neither an elapsed
reset time nor a binary upgrade is a join or scope-widening instruction.

Once a ruling is bound, its committed commit/blob, digest and quote remain the authority.
The initial bind rejects a changed live/archived working copy. Later owner annotations or
edits to those inbox copies do not invalidate recorded membership; missing or altered
committed evidence still gates reads and recovery.

## Normal knob-off confirmations and catch-up

For a recorded policy with `enabled: false`, edit the actual `participants:` in
`00-prompt.md`. An exclusion keeps the ordinary §9.0 confirmation:

```text
excluded: [agent-id — reason — confirmed YYYY-MM-DD]
excluded: agent-id — reason — confirmed YYYY-MM-DD
excluded: [agent-id — unavailable — quota exhausted — confirmed YYYY-MM-DD]
```

The trailing ` — confirmed YYYY-MM-DD` is the delimiter. Reasons may contain em
dashes. Use a real calendar date. A note after the date, hyphen separators, an inbox
reference alone, and `unavailable; user confirmed ...` are rejected; rewrite the marker
in the displayed grammar. Markers validate an exclusion; the CLI never subtracts them
from membership.

A known agent returns through a `participants:` edit. During round 1, any policy-off
join or return uses that plain edit, matching the original behavior; keeping or removing
a stale exclusion marker does not affect it. Do not roll status back to round-01 to bypass catch-up.
After round 1, a kickoff-excluded return or new
participant uses the existing §5 catch-up path: read prior rounds, write a valid late
`round-01/agent-id.md`, edit `participants:`, and join from round 2. Its own late artifact
can be dispatched before or after the edit:

```sh
parley agents exec --agent agent-id --artifact parley-deck/ideas/IDEA/round-01/agent-id.md --prompt-file catch-up-prompt.txt --yes
```

An edit made first shows pending catch-up with this command. Existing members can still
append signoffs, but the joiner gains no historical membership and cannot complete the
quorum until its valid late artifact is imported. A failed own incomplete artifact can
be retried in place. This exception never launches a known excluded member.

A pending joiner may instead file the §5 decline using the existing CLI form:

```sh
parley consensus signoff --agent agent-id --status block --notes '❌ NON-PARTICIPANT' --counter 'Continue without me' IDEA
```

The exact note is accepted only in design consensus for a pending policy-off joiner.
It stays a BLOCK and leaves the joiner missing from completed votes; it does not remove
the joiner, change history or authorize closure. The organizer must resolve the decline
through the ordinary protocol. The literal `--status '❌ NON-PARTICIPANT'` is not a
supported status. No `included:` marker, new return/join record,
committed-answer schema, exact directive or new command is required on this path.
Prior-reading and joining from round 2 remain protocol duties; this compatibility path
does not require new frontmatter to attest them.

At the next normal mutation boundary the CLI captures the prompt and late round, if any,
as an immutable **manual revision**. It does not label the edit owner-confirmed, infer
owner authorization, or release a retained veto. Explicit owner confirmation remains a
protocol obligation for re-inclusion, and authority-sensitive gates require additional
proof through the owner ruling/committed authority paths below. Manual revision history,
old artifacts and retained obligations are preserved. An exact historical notice carrying
the former owner-confirmed display label is preserved with a separate manual-authority
clarification; its label grants no authority. Legacy ideas without quota history
keep their pre-existing behavior. Ordinary off, legacy and kickoff-only driving has no
idea-wide singleton lease.

Kickoff notices use the same checked, non-blocking publication, creating a missing safe
inbox. A missing or unwritable inbox does not undo the kickoff, event log or manifest.

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

Policy-off changes against CLI 1.50.0 are explicit: `run --yes` filters confirmed
exclusions, and a bare preflight 503 cannot be excluded with `--yes`. During round 1,
any policy-off join or return uses a plain participant edit, as before;
removing a stale exclusion marker does not change that. After round 1, a new joiner
or kickoff-excluded return must import its late round-1 artifact before signing or
completing quorum; an exact `NON-PARTICIPANT` decline remains in Missing signoffs.
Pending catch-up stops ordinary driving until the manual `agents exec` command
completes it. Membership edits are frozen on final/closed ideas. Legacy ideas without
quota history retain their prior behavior.

For policy-on ideas, an edit after an applied revision is blocked, including an edit back
to an older participant set. Preserve the proposed edit for owner review, restore the
recorded prompt, then use `parley quota revise` with the committed owner decision. The
error names that command. `quota recover` replays only a genuinely pending transition:
the latest durable receipt must be missing or invalid, and the prompt must match that
transition's before/after membership and policy. It does not overwrite unrelated edits.

## Retained vetoes and dispositions

An owner ruling must identify the retained obligation in the consensus disposition paragraph,
with `Disposition: operator-ruling`, `Authority: owner`, a rationale and canonical `Evidence:`
path. Include `Owner-answer`, `Owner-commit`, `Owner-blob`, `Owner-sha256` and `Owner-quote` there.
The independent participant's evidence artifact must quote that same answer under
`## User direction` and identify the obligation. Deleting the inbox copy later does not undo a
committed ruling. Missing or changed committed evidence does block it.

Alternatively, after re-inclusion bound to committed owner authority, the original veto author can file a later round
artifact under their own canonical filename. Bind it with `quota-revision: N` and
`reinclusion: batch-HASH`, and record `Withdrawal: obligation-HASH`. The consensus disposition
uses `Disposition: withdrawn`, `Authority: author-id`, rationale and that `Evidence:` path.
A manual policy-off revision, even with an optional `included:` marker, cannot authorize this withdrawal. Withdrawal removes only the retained veto it identifies. The returning author must append a new
signoff; a later new veto is still a veto. Closed artifacts remain frozen.

These checks establish recorded attribution and exact content. They cannot authenticate the human
behind `from: user`, establish the truth of a statement, prove that priors were actually read, or
stop someone with write access from fabricating an entire cooperative record. The owner and
independent reviewer remain responsible for those judgments.

## Checked receipt recovery and locking

If a `quota-applied/batch-HASH` receipt is missing, empty or contradictory, membership-dependent
actions fail closed. After stopping writers, run:

```sh
parley quota recover --dir WORKSPACE --idea IDEA --run EXISTING_RUN
```

New imports validate the kickoff and every affected run manifest before committing an
immutable revision. If an older interrupted import already committed without its original
manifest, it stays pending. Restore the original manifest from trustworthy recorded
history/backup, then run checked recovery. Do not construct a replacement identity or
rewrite the immutable revision to make recovery pass.

Recovery validates immutable history and the idea-bound run manifests, reconciles the prompt,
all extant idea manifests and round evaluation, attempts non-blocking notice publication where
needed, then replaces mutable receipts with checked durable writes. Existing receipt bytes are
never proof that projections are sound. Repeating recovery does not duplicate a terminal evaluation.
Contradictory manifests require restoration from recorded authority; owner notice contents are
never that authority. Recovery does not invent missing original runs or migrate budgets.

Symlinked deck or idea-scope directories are unsupported where quota leases or scoped
integrity escalations are required. Enabled mid-idea creation refuses such a deck before
writing kickoff state. Use a physical `parley-deck` directory in the workspace, or disable
`quota_auto_exclude` for ordinary driving. Legacy, policy-off and kickoff-only ordinary
driving are unchanged; manual imports and owner revisions still need a physical scope.
Workspace-root and higher ancestor aliases (including macOS `/tmp`) resolve normally and
share the same physical lease. No lease is derived from an unrelated ancestor's name.

Lifetime and projection leases live under the already ignored `.parley-runtime/membership/` on
the deck filesystem. Ownership uses an exclusively published complete PID/token record, host and
boot identity. Projection acquisition waits at most two seconds; it is unrelated to provider retries.
Only a proven dead PID on the exact known host and boot can be reclaimed automatically. Unknown or
different host/boot, partial owners and interrupted reclamation fail closed for operator investigation.
Per-generation reclamation claims prevent a delayed stale reader from deleting a newer winner.
Do not remove a lock while its owner or a possible remote owner may still be alive.

A crashed invocation with `started.json` but no `terminal.json` can be settled by checked
recovery only when its recorded host and boot match this host and boot, its supervisor
and writer PIDs are proven absent, and its recorded supervised process group is absent.
The group must have been the writer's own group. The CLI writes an immutable
`crash-settlement.json` containing invocation identity, start-record digest, proof and time;
replay validates and syncs that record without changing it. It never manufactures a normal
terminal result, provider error, completed artifact or budget accounting settlement.
Live descendants, ambiguous liveness, foreign/unknown identity and older identity-less
records stay blocked. Investigate on the original host and restore authentic terminal
evidence if available; no force-unlock or synthetic owner-authorization route is provided.
Windows automatic crash settlement remains unavailable until runtime proof exists.

## Zcode evidence limit

The stderr recognizer requires complete allowlisted SDK error framing, matching response bodies,
all retry records, reset values, terminal failure and invocation facts. Complete top-level
records and complete SDK retry aggregates may contain decreasing countdowns when their
exhaustion class and absolute resets agree within one second. Each record's duration,
header and countdown must agree. Inferred attempt times must be monotonic within one
second, no later than receipt and no earlier than a recorded invocation start. The fixed
60-minute threshold uses the terminal line's receipt, never an earlier inferred attempt. Arbitrary prose, mixed errors,
unknown fields, incomplete captures and SDK inspect truncation gate. Source-derived SDK fixtures
exercise the retained incident's positive body/reset, but are not native captures. The two retained
native excerpts do not recover the missing 24,833-byte original; full-native AC2 evidence remains
unverified. Other adapters remain diagnostic-only. The owner's accepted subagent-source ambiguity
is unchanged.

Offline inspection of installed zcode 3.7.7-13 locates the SDK stream error sink's default
`console.error(error)` with no inspect options added by the adapter. The located streaming
adapter disables SDK retries (`maxRetries: 0`) and implements its own attempt loop; complete
top-level records are therefore a relevant source-derived shape. Its default-console
reproduction with realistic request messages contains `[Object]` and is rejected. Deep
inspect fixtures and retry aggregates are source-derived tests, not proof of the historical
incident's exact runtime path or a complete native capture.

**Owner-accepted known limitation for this release:** native-positive AC2 is NOT MET
and is owner-waived by round05-answer Q2 (2026-10-06). zcode auto-exclusion may not fire
on real native output; unrecognized failures fall back to the owner-confirmed path.
R5-MAJOR-2 is deferred, not fixed. Complete source-derived positives are not native
verification. The linked follow-up is
[quota-zcode-native-exhaustion-capture](../parley-deck/ideas/quota-zcode-native-exhaustion-capture/00-prompt.md).

Capture is not automatic in this release. On the next ordinary zcode failure, keep
that run's private, unscrubbed `parley-deck/runs/<run-id>/agents/<agent-id>/stderr.log`
in place; never commit or copy it. Its completeness and cleanup are unverified.
The follow-up will establish a complete scrubbed capture before changing the grammar.


## Accepted release limitations

Kickoff reporting and notice replay are repaired in 1.52.0. A missing safe inbox is
created; unavailable publication prints the full blocked decision to stderr. Kickoff
notice publication uses the existing completed-attempt receipt and crash replay,
preserving owner-edited/archived copies. A receipt is not proof of delivery.

1. **Aliased decks and plain edits (R8-NIT-2).** On a symlinked deck, a plain
   `participants:` edit or confirmed exclusion counts as a manual revision and
   requires a physical deck path even when `quota_auto_exclude` is false.
   Until that path is restored, the pending edit blocks all signers and driving.
   Disabling the policy again does not resolve this case; the generic diagnostic's
   suggestion to do so is inapplicable.

Follow-up: [quota-kickoff-reporting-and-alias-guidance](../parley-deck/ideas/quota-kickoff-reporting-and-alias-guidance/00-prompt.md).
