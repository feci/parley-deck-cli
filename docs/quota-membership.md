# Recorded quota membership and recovery

New ideas record the resolved quota policy and initial membership in `quota-kickoff.json`.
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
revision preserves the quote and object identifiers. A missing live inbox file is allowed; a
contradictory live or archived copy, changed digest, wrong idea/author, missing object or fabricated
path is rejected. Keep the committed Git object available for historical verification.

Only previously known identities qualify as re-inclusions. A new identity additionally needs
`"catchup":{"new-id":"round-01/new-id.md"}` and a valid late round-1 artifact. Its frontmatter
must declare `catch-up: true`, `read-priors: [round-01/other-id.md, ...]`, and `join-from: round-02`.
The initial command checks the referenced prior artifacts and captures the catch-up artifact in
immutable history. The owner directive must authorize the new current set. Neither an elapsed
reset time nor a binary upgrade is a join or scope-widening instruction.

## Normal knob-off confirmations

For a recorded policy with `enabled: false`, existing protocol confirmations still work without
`quota revise`. Record the actual current `participants:` in `00-prompt.md` and the ordinary
`excluded: [agent-id — reason — confirmed YYYY-MM-DD]`. A known return uses an explicit
`included: [agent-id — reason — confirmed YYYY-MM-DD]`. These markers validate the requested
current set; the CLI never subtracts markers to derive membership. At the next normal mutation
boundary, it snapshots and reconciles the confirmation into immutable history. It retains filed
obligations and keeps the policy off. Ordinary off, legacy and kickoff-only driving has no idea-wide
singleton lease.

A new off-mode catch-up join also needs the late round and prior-reading conditions above, plus an
owner answer. Record `owner-answer`, `owner-commit`, `owner-blob`, `owner-sha256` and `owner-quote`
in the late round's frontmatter. The committed quote must contain the standalone directive
`Catch-up join: new-id from round-02`. This uses the existing confirmation path; no new CLI command
is required. Legacy ideas without quota history keep their existing protocol behavior.

## Retained vetoes and dispositions

An owner ruling must identify the retained obligation in the consensus disposition paragraph,
with `Disposition: operator-ruling`, `Authority: owner`, a rationale and canonical `Evidence:`
path. Include `Owner-answer`, `Owner-commit`, `Owner-blob`, `Owner-sha256` and `Owner-quote` there.
The independent participant's evidence artifact must quote that same answer under
`## User direction` and identify the obligation. Deleting the inbox copy later does not undo a
committed ruling. Missing or changed committed evidence does block it.

Alternatively, after authorized re-inclusion, the original veto author can file a later round
artifact under their own canonical filename. Bind it with `quota-revision: N` and
`reinclusion: batch-HASH`, and record `Withdrawal: obligation-HASH`. The consensus disposition
uses `Disposition: withdrawn`, `Authority: author-id`, rationale and that `Evidence:` path.
Withdrawal removes only the retained veto it identifies. The returning author must append a new
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

Recovery validates immutable history and the idea-bound run manifests, reconciles the prompt,
all extant idea manifests, round evaluation and notices, then replaces mutable receipts with checked
durable writes. Existing receipt bytes are never proof that projections are sound. Repeating recovery
does not duplicate a notice or terminal evaluation. Contradictory manifests or notices require
restoration from recorded authority; recovery does not invent missing original runs or migrate budgets.

Lifetime and projection leases live under the already ignored `.parley-runtime/membership/` on
the deck filesystem. Ownership uses an exclusively published complete PID/token record, host and
boot identity. Projection acquisition waits at most two seconds; it is unrelated to provider retries.
Only a proven dead PID on the exact known host and boot can be reclaimed automatically. Unknown or
different host/boot, partial owners and interrupted reclamation fail closed for operator investigation.
Per-generation reclamation claims prevent a delayed stale reader from deleting a newer winner.
Do not remove a lock while its owner or a possible remote owner may still be alive.

## Zcode evidence limit

The stderr recognizer requires complete allowlisted SDK error framing, matching response bodies,
all retry records, reset values, terminal failure and invocation facts. Arbitrary prose, mixed errors,
unknown fields, incomplete captures and SDK inspect truncation gate. Source-derived SDK fixtures
exercise the retained incident's positive body/reset, but are not native captures. The two retained
native excerpts do not recover the missing 24,833-byte original; full-native AC2 evidence remains
unverified. Other adapters remain diagnostic-only. The owner's accepted subagent-source ambiguity
is unchanged.
