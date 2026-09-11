---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-11
source-commit: 1a3d512e8909af3c82e418c04d15aae646675ced
status: tested-slice-independent-acceptance-pending
---

# Historical protocol accounting with retained attempt boundaries

## Resulting behavior

`budget migrate inspect|apply --kind step|fixup|cross-review` imports historical
driver steps, fixups and cross-review groups into separate persistent policies.
The runtime guide supplies concrete commands. Inspection is read-only. Apply
requires the actual attended operator control, stopped writers, exact history
hash, decision identity/reason, original accounting epoch, explicit complete
historical action total and lifetime ceilings. Zero totals are explicit. An
imported exhausted state remains exhausted; step zero is unlimited and cycle
zero is forbidden. No real policy was migrated and no provider was invoked.

The source inventory includes all available linked worktrees, requests and
terminal metadata, run events/cursors and canonical artifacts. Nested pipeline
idea paths are supported and hash-bound. Distinct published driver transitions
add across runs; identical observations and copied runs do not duplicate them.
Conflicting copies and run/idea identities refuse. Counter snapshots can overlap
published actions and each other, so they contribute a maximum. Distinct completed
fixup marker paths contribute a union across worktrees. Reported elapsed step
time can move the retained accounting origin earlier than the first event.

The observed lower bound is not the lifetime total. Legacy child requests lack
synchronous operation IDs: typed attempts establish at most one ungrouped action,
not one charge per participant. Unknown phases, unobserved handoffs and pre-start
refusals remain visible for operator reconciliation. A launch refusal may follow
a charged protocol step. A fixup phase event can finish an already charged cycle
after a crash; another publication is not proof of another code-writing attempt.
The explicit operator total must cover failed, partial and otherwise unobserved
work and cannot be smaller than the observed floor. Prose supplies no count.

Publication uses an immutable import record, initial ledger and continuity
witness, a policy referencing the import, then a final activation marker. Exact
replay recovers unchanged partial publication and retains later charges and
finite grants after activation. Existing configured/charged scopes cannot be
replaced. Missing active policies, stale/conflicting decisions, malformed or
unavailable sources and changes before activation refuse. Removing required
zero/nullable fields cannot normalize into valid authority. Runtime checks bind
the original epoch and each imported identity, including cached nested sessions;
equal aggregate counts cannot hide replacement of an imported entry.

Imported entry identities use `protocol-migration:<kind>:<decision-id>:<index>`.
Their time is the operator-declared accounting epoch, not an observed child
start. Zero monetary values represent protocol-count entries only. Model cost
and unknown observed prices remain in the independent launch ledger.

The existing track/runtime gates remain binding. Migration does not parse or
rewrite the track: choose compatible original ceilings before applying. A
mismatched import can be unusable; it cannot bypass a stricter track. Driver
configuration must match the original cycle ceiling; finite extensions retain
that original reference. Fast cross-review refusal is tested through the actual
runner. No unsupported policy rewrite is presented as recovery.

## Executed negative cases and fixture corrections

New history tests failed before correction: five distinct steps in two runs
were reduced to three; conflicting run IDs and unknown action kinds were
accepted; and two fixup/recovery publications were treated as two fresh fixups.
The corrected inventory passes these cases. The original history-negative log
is retained in the validation directory.

The first actual-driver fixture failed because it expected three lifetime steps
after a fresh BLOCK attempted a fourth step and then hit its nested cycle cap.
The runtime correctly kept that fourth step spent. The corrected fixture checks
three steps immediately after the two nested children, then the appropriate
three/four steps after the BLOCK boundary. It also checks unchanged cycle count,
original epoch, exact replay and only two actual child starts. No runtime refund
or relaxed count assertion was introduced. The original driver log is retained.

Actual local helper processes cover pre-import refusal, the final allowed cycle,
grouped participant launches, failed attempts, changed run identity, BLOCK and
post-cycle verification without another fixup charge. A separate driver fixture
exercises the final-cycle completion path with fake verification adapters; it
does not establish independent real-model closure. These local process fixtures
are not additional model attempts or independent participant acceptance.

## Validation checkpoint

All 338 Go/module files match source-manifest.json, SHA256
b111366f3d8145c3b893426fc1c5fdd863a8b40a5712e7fe4df3f520ea51b02d.
Exact logs and manifests are under
`.parley-runtime/protocol-migration-validation-20260911/`.

- Full Go suite: PASS, wall 134.046s.
- Budget/driver/runner/evidence/app race suite: PASS, wall 147.141s.
- Focused migration/CLI/driver/runner tests: PASS across all four modules.
- Scoped vet: PASS. Windows amd64 app cross-build: PASS; runtime unverified.
- Isolated shared-volume publication, worktree and actual-process fixtures:
  PASS, wall 13.366s. Fixtures are outside the real repository's Git ancestry.
- Actual compiled CLI rejects all three unattended migration kinds with exit 2,
  the attended-control diagnostic and no created accounting state.

Full log SHA256: a823173fc269e51b144de271a0feb1a6ab8ff9f07cc2649dd9f3fbf590af9b40.
Race log SHA256: 8583920214c2ac058974fb39b853388b7654cfd65e1ad0c323cbe2af9880ddb3.
Negative history log: 04ee8fe6f5de55c88876c60f5e8f85d148dd3830b996e5af0c869a0f1fc0be85.
Initial driver fixture log: 95d5cfeecc6ef6dcaffbdcec12618e3d7143c53653f26e8ca429a11b6cc6b920.
Prior fd3e3ce validation and report QA remain separately retained.

The updated 617000-byte English offline report, SHA256
c2cbdfb81d5c77641bf7294ca72e0aecab7ce10f75c68d5fbc355487c59b7b07,
passes all ten report tests and 55 live ego-browser assertions at 1440x900,
1280x540, 390x844 and 320x720. All panels/documents, expanded current progress,
painted canvases, search/reset, pagination, keyboard section navigation and
print-media panel visibility pass. The exact sidecar is
`parley-deck-evaluation/delivery/2026-09-05/browser-protocol-migration-20260911.json`
in the sibling evaluation workspace; raw logs/observations/checksums are under
its `browser-qa/protocol-migration-20260911/` directory. Native DOM clicks have
explicit state readback and keyboard checks use real key dispatch. No fresh
screenshot review, physical click-placement or printing/PDF pagination is
certified. This is the current partial report, not future populated-result QA.

## Remaining limits and full objective

Hashes and attended/stopped-writer assertions do not authenticate a human or
coordinate uncooperative same-filesystem/distributed writers. Changed inactive
imports and guard/lock-origin recovery remain unsupported. Import is not durable
semantic operation replay. Windows runtime and fresh independent acceptance
remain unverified. No current owned review or signature is inferred from these
implementer checks; Claude's failed provider-limit attempt remains preserved.

The six audit recommendations remain the full goal. Safe guard recovery, durable
semantic action identity/replay, canonical refusal publication/recovery, the
opt-in independently confirmed two-patch regression trajectory and human
evidence-table integrity remain open. So do independent real-model concurrency
and closure, complete live launch coverage, the exact twelve phase 1/6 AB/BA
packet pairs with three canaries and full control, the twelve-task solo/duo/full-six
pilot with frozen equal ceilings/rotation, two blind nonauthor graders, owned
reviews/signatures, final populated HTML QA and actual 14/30-day observations.
Historical quorum/pilot/funding and Claude recovery decisions remain unanswered.
The known model inventory remains 34 terminal attempts, 17 unknown prices and
USD 46.1887585 in known CLI estimates; the total cost is unknown.

No merge, release, deployment, global install or immutable-core publication is
included. The implementation stays in draft PR #73. Shared OpenViking tools are
unavailable in this session; no shared-memory persistence is claimed.
