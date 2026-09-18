---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
source-commit: fd3e3ce68080a0b7b39d168e153af0280221f94e
status: tested-slice-independent-acceptance-pending
---

# Historical launch accounting without erasing execution evidence

## Resulting behavior

A scope with historical invocation records can receive an explicit attended
launch-accounting import through `budget migrate inspect|apply --kind launch`.
This also resolves the reproduced bootstrap dead end where a monetary default
retained a pre-start refusal, then ordinary configuration correctly refused to
ignore that history. The runtime guide documents all concrete flags and units.

Read-only inspection inventories every available linked worktree, requested and
terminal invocation metadata, matching start evidence, historical run events and
cursors, and canonical idea artifacts. It returns content hashes rather than
copying prompts, logs or artifact contents into the decision. Every available run
identity must agree, including explicit event idea/idea_slug and cursor fields.
Unavailable worktrees, unknown/nonterminal records, conflicting copies, aliased
parents/leaves, malformed JSON and inconsistent start/terminal metadata refuse.

The operator supplies the exact history hash, decision identity and reason,
original accounting epoch, explicit additional legacy attempt count, all three
lifetime ceilings and stopped-writers assertion. A positive cost cap also needs
a conservative per-future-launch reservation. The CLI independently checks real
terminal attendance and platform support; no flag or participant field grants it.
The original epoch cannot be later than observed history. Additional attempts
must meet the observed legacy start floor but are never inferred from prose.

Unique recognized attempts become spent entries with their original timestamps;
ordinary failed starts count. Known costs round upward to integer microdollars,
while the original terminal hash and cost basis remain in the inventory. Unknown
costs stay unknown. Historical additional attempts and unobserved handoffs keep
unknown exposure; a future reservation never supplies a retroactive price.
Additional entries use `legacy-migration:<decision-id>:<zero-based-index>` IDs.
Explicit monetary reconciliation is required before unknown exposure becomes
spendable under a monetary cap.

Typed pre-start budget/context refusals with no process or observed work remain
in the inventory but do not fabricate a model execution charge. Conflicting
start events, token/cost, output, activity, exit or truncated-input evidence
refuse that exemption. Source records remain byte-for-byte unchanged.

The immutable import record binds the operator decision, inventory and initial
ledger. The ledger and continuity witness precede the hash-bound policy; a final
activation marker gates readers. Exact replay recovers unchanged partial writes
and preserves later spend, reconciliations and extensions after activation.
Conflicting/stale decisions cannot overwrite existing policies, ledgers or
imports. A missing active policy cannot be reconstructed as first configuration.
Imported charges and the original epoch remain checked on subsequent inspection,
reservation and extension. Required zero-valued fields must be explicitly present
rather than silently normalized into the same digest. Older readers reject the
new migration-reference field.

Sources are rechecked before the publication guard, under it and before final
activation. Detected changes during publication leave the import inactive. The
operator's stopped-writers assertion is not a distributed lock or authentication;
uncooperative external writers and altered inactive imports still need explicit
recovery. Inventory bounds are 10,000 sources/actions and 64 MiB source bytes;
the immutable import record is limited to 8 MiB.

## Executed counterexamples and validation corrections

New tests first reproduced two defects in the unfinished migration candidate:
conflicting event/cursor idea identities could silently exclude historical work,
and omitted or null fields could be normalized into zero/unknown accounting
without invalidating the semantic digest. The tests failed before the fixes and
now pass. The original negative-history.log is retained, SHA256
3f45969077b6f496ae38dc082718f3f9ee0868012ea4479dcd908a1bac415564.

The first full Go suite failed in the existing eight-contender cycle-extension
test: the bounded reader correctly refused an atomic policy replacement during
its initial read, but the test accepted only the later stale-hash refusal. Its
final one-winner, one-grant and unchanged-spend assertions still passed. An
internal typed changed-history sentinel now distinguishes that exact outcome;
the test accepts it without weakening those accounting assertions or accepting
arbitrary errors. The original full.log is retained, SHA256
41571eb198f005f2e82387bcf00eb838be5b69982a0bf5c6d946de5b757ec78a.
The prior candidate's complete scoped race suite passed; it is not presented as
validation of the subsequent sentinel change.

The first shared-volume harness mistakenly placed non-Git fixture roots inside
the integration worktree. They inherited the real Git ancestry and correctly
refused a registered but unavailable historical worktree. The first actual CLI
probe stopped at the same read-only check. The corrected fixtures run outside
any Git repository on the same shared volume; their explicit Git-worktree case
still creates its own disposable repository. The failed logs and separate
harness-observations.md retain this distinction; no real history was pruned and
no production bypass was introduced.

## Current source validation

Source fd3e3ce68080a0b7b39d168e153af0280221f94e matches all 330 Go/module files in the retained
source manifest (SHA256 19e6fc571bfa65ee1869e69e4d1bd3f02f978bf3df5e078d2ca36ae3642eff21).

- `go test -count=1 ./...`: PASS; wall 93.717s.
- `go test -race -count=1 ./internal/budget ./internal/driver ./internal/runner ./internal/evidence ./internal/app`: PASS; wall 107.389s.
- Corrected focused migration/CLI/actual-process/step-history/concurrent-cycle tests
  with race detection: PASS; wall 8.944s.
- Scoped vet and Windows amd64 app cross-build: PASS. Windows runtime untested.
- Isolated shared-volume publication-recovery, shared-worktree and actual-process
  fixtures: PASS; wall 1.695s.
- Actual unattended CLI apply: exit 2 with the explicit attended-control refusal;
  the disposable fixture remains uninitialized.

The fixtures cover before/after failures at all four publication boundaries,
exact replay with later spend/reconciliation/extensions, conflicting concurrent
migration decisions, lost active policy, corrupted witness/reference, stale or
changed sources, malformed/omitted/null fields, legacy lower bounds, aliases,
unavailable worktrees and bounded reads. The actual process counterexample
records one started child and three retained budget refusals across new runs;
its opaque output leaves observed cost unknown while exposure stays reserved.

Full-suite log SHA256: f7925422a05262cef4dc140ce7417cdf3cf4d02de63ef361c5c8ce5cef74cc25.
Race-suite log SHA256: 470c22df6225e916de919b27c500bd7935068853826a04350ac58219406b0317.
Retained evidence: `.parley-runtime/launch-migration-validation-20260911/`,
including the failed runs, commands, timing/status JSON, both source manifests,
refusal-sentinel patch, harness observations and final checksums. These local
ignored logs are not claimed to be published in Git or shared memory.


## Limits and remaining obligations

This is implementer execution evidence, not participant-owned independent
acceptance. No real policy was migrated/extended, and no participant CLI was
invoked for this slice. Windows runtime remains unverified. The CLI's attended
probe does not authenticate a human against another process with filesystem
access. Unknown or changed inactive imports and lock-origin relocation need
separate recovery; do not delete their records to force migration.

Step/cycle legacy import, safe guard/lock-origin recovery, durable semantic action
identity/replay, canonical refusal publication/recovery and opt-in independently
confirmed two-patch regression trajectory remain open. Human evidence-table
integrity and fresh independent source acceptance are not certified here.

The signed exact Phase-1/6 packet experiment, independent real-model concurrency
and closure, twelve-task solo/duo/full-six comparison with frozen equal ceilings
and role rotation, two blind nonauthor graders, participant-owned final reviews
and signatures, final populated report QA, and actual elapsed 14/30-day follow-ups
remain binding. Historical quorum/pilot/funding and Claude recovery choices are
unchanged. Full context remains the default. The actual model inventory stays
34 terminal attempts, 17 unknown costs and USD 46.1887585 known CLI estimates;
total cost is unknown. The all-six goal remains incomplete.


## Report checkpoint

The refreshed English self-contained report is 613084 bytes, SHA256
306c99b7f1cac02452ab724fb999076db00da3f9b82a937a6e40d234f6745b66. All ten report tests pass. The exact file passes
55 live ego-browser DOM/layout/canvas/interaction assertions across 1440x900,
1280x540, 390x844 and 320x720: all panels/documents, expanded technical log,
painted charts, filters/reset, pagination, keyboard focus and print-media panel
visibility. Sidecar: evaluation delivery/2026-09-05/browser-launch-migration-20260911.json.
Raw checks and failed harness logs are retained in its browser-qa subdirectory.

Fresh screenshot capture timed out through raw CDP and the helper, including
after bringing the task tab forward. Coordinate helpers also failed to activate
an offscreen summary/navigation item; the final checks use native DOM click and
state readback, with actual key dispatch for keyboard navigation. Fresh visual
screenshot review and physical click placement are not claimed. The report
content and template were not changed to hide a failing assertion. Physical
printing/PDF pagination and final future populated results remain unverified.
