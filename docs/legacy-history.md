# Durable declarations of legacy unknown history

A new idea's cycle budget refuses a historical run whose idea identity is missing.
`parley budget legacy` provides a one-time attended scope decision for a narrow class
of readable legacy records. It does not assert zero historical work or change any
cycle count, ceiling, worktree-coverage gate or existing per-idea migration.

First inspect the exact canonical relative run directory:

```sh
parley budget legacy inspect --dir WORKSPACE --run parley-deck/runs/RUN
```

Inspection writes nothing. Its JSON binds the complete manifest digest, all visible
copies and existing decision digests in `preview_sha256`; `history` is explicitly
`unknown-history`. No numeric field counts work in that unknown region.

After stopping all writers, the owner applies that exact preview from an attended
terminal, supplying a unique decision ID and a concrete reason:

```sh
parley budget legacy apply --dir WORKSPACE --run parley-deck/runs/RUN \
  --expected-preview-sha256 PREVIEW_SHA256 --decision-id DECISION --reason REASON \
  --writers-stopped --acknowledge-unknown-history --yes
```

Agents may prepare the preview and command. They must not allocate a terminal or
impersonate owner attendance. A stale preview refuses; inspect again and review the
changed state. Exact request replay returns the original decision, including its
recorded time. Reusing the ID for another request, or declaring the same path with a
second decision, refuses.

Eligibility is deliberately closed: exactly one regular, single-link `events.jsonl`,
nonempty valid strict JSONL, and only identity-absent `run.created` or `run.phase` events.
Creation data may contain descriptive `mode`/`task` strings; phase data must contain a
known non-cycle `action`. Unknown data fields, charged phases, recovered identity,
cursors, hidden files, extra directories, symlinks and hard links refuse. Richer history
requires the existing attended per-idea `budget migrate` accounting decision.

The append-only decisions live at
`<git-common-dir>/parley-launch-budgets/legacy-<scope-hash>/records/<decision-id-hash>.json`,
scoped by the deck's repository prefix; a non-Git deck uses its local runtime origin. First cycle
binding retains adopted payloads and their digest in `legacy.json` and policy. `budget
cycle inspect` discloses them. Every use checks the original decision, every visible
copy and additional unscoped cycle events. Mutations, deleted copies, unavailable
roots, lost/corrupt provenance or extra unknown history refuse before a new charge.
An already-bound idea cannot silently adopt a later declaration; its original adoption
stays frozen. Existing charged counts and caps remain independent of this uncertainty.

Installing the release does not create a declaration. D6 remains blocked locally until
the owner performs this attended activation. Older binaries do not consult the new
ledger, and reject a new policy field rather than acquiring authority from it.
