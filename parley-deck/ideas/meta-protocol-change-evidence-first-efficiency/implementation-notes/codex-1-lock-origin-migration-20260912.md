---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
status: implemented-awaiting-independent-review
checkpoint-base-commit: 498a21cae3cadccf373208a55df300c62961c525
---

# Same-host lock-origin relocation

A changed cache environment previously refused an established resource because
its original host/path/token remained pinned. The new attended `budget origin
inspect|apply` control can move that synchronization authority while retaining
all accounting and the old permanent kernel object. This is a same-host move
of an unchanged canonical resource path, not repository/scope migration.

## Cutover and recovery contract

Read-only inspection takes the accessible original lock and binds exact origin,
guard-witness presence, destination identity/absence, and protected file hashes
and modes. Apply reacquires and checks the old lock, repeats the protected
inventory after waits, and publishes an immutable decision before creating the
destination identity. A second kernel lock and independent-descriptor exclusion
probe must succeed. Destination contention refuses immediately with the decision
retained, avoiding an indefinite two-lock wait.

Both old and destination descriptors remain held and are checked against their
current paths throughout origin, guard-witness and completion publication. No
ledger byte, reservation, settled charge, unknown cost, ceiling or epoch changes.
`ledger-established` binds accounting scope and stays unchanged;
`guard-established` binds the exact synchronization origin and is updated.

Migrated origins use v3 with an explicit decision reference. V2 binaries refuse
that format. The current reader validates the entire completed chain and pins
exact origin bytes before a kernel wait, including the migration reference. It
never adopts authority read only after acquisition. Missing records/completions,
required witnesses, copied roots, altered identities and missing original objects
refuse. Retained migration history also prevents fresh bootstrap after loss of
origin/witness. Original v2 authority pins a token, not a historical device/inode;
the implementation does not claim otherwise.

Exact apply replay returns the original preview after later writes and migrations.
An interrupted apply resumes only against the original material and both accessible
kernel objects. Before origin cutover, old v2 writers can legitimately change
accounting. The resulting stale prepared decision remains immutable; an attended
fresh inspection and new decision can supersede it without a reset. After v3
cutover, ordinary writers refuse until completion. A changed protected inventory
at that stage, missing old identity, erased journal/history or unavailable host
remains an explicit recovery limit. No forced override or automatic cleanup exists.

The inventory refuses aliases/special files and is limited to 10,000 entries,
16 MiB/member and 64 MiB total. It includes descendants without migrating their
independently locked scopes. Chains permit 127 migrations plus the original v2
entry; capacity is checked before a new cutover can publish unreadable authority.
Large snapshot stores and changed repository/ticket roots need separate recovery.

## Verification and attribution

Final automated validation passed on 385 Go/module files at manifest
`ac1589bbf46569489f94ae29cda45c7f7cff4f3fe3ee04c5e4641d97c7e1e002`:

- Full suite: PASS, 210.427s; all 32 package terminals (31 test packages and the
  command package without tests), with no failed tests.
- Six-package race: PASS, 247.623s; budget, app, driver, runner, trajectory, evidence.
- `go vet ./...`: PASS; Windows budget/app PE amd64 cross-builds: PASS.
  Windows runtime remains unverified.
- Compiled shared-volume budget/app selections: PASS, 6.069s / 0.555s.
- Eleven substantive new top-level scenarios plus their child-process harness
  pass in full/race and the applicable shared selections.
- Five negative overlays removing captured-origin checking, completed authority,
  held-inode validation, protected-material matching or the pre-cutover capacity
  guard fail at their designated assertions.
- The actual built CLI refuses unattended `origin apply` with exit 2 and no
  directory mutation; the platform seam is not the only attendance check.

Logs are captured as complete native bytes after process termination, then copied
to the shared result directory. The final verifier checks matching log hashes,
unchanged source, all new test terminals, both Windows binaries, the actual legacy
process proof and unchanged historical HTML. Evidence is under
`.parley-runtime/lock-origin-migration-final-validation-20260912/`.

The first complete Go suite returned PASS
in 214.932 seconds, with all 32 package terminals and no failed tests, but source
changed during that run to add the pre-cutover chain-capacity guard. Its manifest
check correctly returned false; it is retained as a superseded run and does not
certify the final source. The targeted capacity counterexample passed separately.
The accepted final validation above uses its own frozen source manifest.

A native process proof built the actual old budget package at 498a21c and added
only a test barrier helper. Its production lock.go SHA256 is
`87f0c95fced953be2b79ec37107d19fe086ddae6670e2dcb20b3994f883dc965`, verified against
Git. The old holder held its kernel lock; an old waiter reported a real failed
take; the current implementation paused with both locks held after origin cutover.
A new writer refused incomplete authority, the old waiter refused after release,
and a new writer acquired after completion. The original device/inode remained.
Native/shared logs and the helper hash are retained in
`.parley-runtime/lock-origin-migration-validation-20260912/legacy-process-proof.json`.
This is executable compatibility evidence, not participant-owned acceptance.

No real operator resource was migrated. All activations used synthetic test
resources; no model call, budget grant, policy amendment, global install, release,
merge or deployment occurred. Historical live inventory and HTML remain unchanged.
Independent current-source review, semantic action recovery, live experiments,
final participant signatures, HTML delivery and follow-up observations remain open.
