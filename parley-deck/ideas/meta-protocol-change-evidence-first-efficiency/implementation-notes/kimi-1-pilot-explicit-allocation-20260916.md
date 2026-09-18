# kimi-1 — pilot explicit-table allocation note (2026-09-16)

Status: draft implementation note. The whole design remains pending peer round-02;
nothing here is a ratified schedule, a source freeze, or a preregistration claim.

## Scope (bounded launch task)

- `.parley-runtime/pilot-allocation-candidate-20260916/pilot_analysis.py`
- `.parley-runtime/pilot-allocation-candidate-20260916/test_pilot_analysis.py`
- this note

## Plan

- Keep the six-person seeded `allocation(tasks, roster, seed=...)` path and its output
  unchanged, including rejection of a four-person roster with no table.
- Add keyword-only `table=None`: a four-person run is accepted only with an explicit
  caller-supplied 12-row table (task, execution_index 0..11 in list order, solo/duo/full
  cells with participants/drafter/critic). Validation covers shape, indexes (bool
  rejected), task/category metadata (exactly 4 each of code/review/design), roster (four
  distinct nonempty IDs), role membership, and the stated balance checks. Every violation
  raises ValueError; nothing is silently rebalanced, reseeded, or dropped.
- Explicit-table result: `status: draft-allocation-not-preregistered`, `seed: None`,
  `allocation_method: explicit-table`, freshly built rows (caller objects never mutated).
- Tests hardcode the Codex round-02 proposed table as an unratified fixture: it must pass
  with roles/order preserved; tampered variants must fail; six-person seeded output stays
  schema-identical; caller-input non-mutation is asserted.

Tests were not executed in this launch (no runner permitted in scope).

## Outcome

Implemented as planned. `allocation` gained keyword-only `table=None`; the new private
`_validate_explicit_table` enforces every listed structural, metadata, role, and balance
rule and returns freshly built rows (deep-copy guarantee by construction — no shared
mutable state with caller tasks/roster/table). The six-person seeded branch is untouched
and keeps its exact original output keys. Tests add the hardcoded Codex round-02 fixture
(acceptance with roles/order preserved), 16 tampered-table rejections, metadata/roster
rejections, four-person-without-table rejection, six-person schema invariance, and input
non-mutation. Full-arm participant *ordering* balance is deliberately not required, per
the launch bounds. Design remains pending peer round-02; no ratification is claimed.
