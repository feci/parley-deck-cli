# Four-person allocation helper — bounded validation and import

Status: implementation preparation, not amendment ratification, experiment freeze,
preregistration, treatment, full-branch signoff, or completion of the original audit.

Kimi invocation `6e52e998-02ac-4e7a-8b34-cf5e84672d8a` stopped with exit 0 before
coordinator inspection. The participant's original source and own note remain
unchanged in retained evidence. The note is copied byte-for-byte alongside this file.
Its original claim that every malformed input raises ValueError is too broad for
those original bytes: coordinator executions reproduced KeyError for a missing ID
and TypeError for unhashable roster, task-ID/category entries and nonmapping tasks.

## Delivered behavior

The sibling evaluation helper now accepts a four-person roster only through the
explicit `table=` option. It checks twelve unique tasks, four per category;
1/2/4 arm membership; role membership and distinctness; three occurrences per
participant of each named role; one solo per category; all directed duo and full
role pairs; all six arm orders twice; and execution indexes. It returns fresh
structures marked `draft-allocation-not-preregistered`. It does not select or
freeze a schedule. The existing six-person seeded path remains the default.

After stopped ownership, Codex added shape/type checks before hashing/indexing
explicit-table metadata and roster entries, plus boundary regressions, tuple-input
coverage and old-module golden values. These corrections are attributable to Codex,
not retrospective edits of Kimi's own note. Full-arm participant ordering balance
is outside this helper's stated contract; the complete experimental design still
requires peer review and ratification.

## Executed checks and limits

- Original Kimi candidate: 14 tests passed.
- First wider run: 27 tests passed; the grading module failed to import because
  the system Python 3.14 lacked PyYAML. This failed run is retained.
- Isolated Python 3.11.15 with PyYAML 6.0.3: all 35 allocation, grading and
  acceptance tests passed. Command: `python -B -m unittest -v test_pilot_analysis
  test_pilot_grading test_pilot_acceptance`.
- The two new boundary tests against the original Kimi production module failed
  with ten subtest errors; the corrected module passed them. This before-control
  distinguishes the new regressions from tests that merely echo implementation.
- Thirty actual old/new six-person allocation comparisons were equal: five fixed
  seeds, two task orders and three roster orders. Four original-module golden
  hashes are retained in the tests.
- The fixture's twelve task/role/order rows equal the unratified Codex amendment
  round-02 table byte-derived through a separate Markdown parser. This checks
  transcription, not the design's acceptance.

The two exact tested Python files were imported into
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck-evaluation/scripts/` only after
both original hashes and candidate hashes matched. The evaluation directory has
no Git history; its originals are preserved under
`.audit-backups/20260916-explicit-allocation/`. All eight historical pilot input
files and the seventeen other Python scripts were hash-checked unchanged.
The adjacent `.patch` gives the complete reviewable change relative to those
originals; the adjacent validation JSON records commands, times and hashes.

Raw evidence: `.parley-runtime/managed-continuation-20260915/pilot-allocation-validation/`.
No new model call, experiment spending, cycle migration or global installation
occurred during these coordinator checks. Kimi's previously recorded call cost
remains unknown; implementation accounting is separate from experiment spending.

## Outstanding gates

Peer amendment round-02 and consensus remain pending. The proposed historical
cross-review accounting total/epoch still have no user approval. No migration or
round launch occurred. The actual provider route still lacks demonstrated
conservative USD15 total spending enforcement. This bounded helper check does not
close independent current-tree verifier/concurrency execution, whole-source
review/signoffs, the exact 28-call packet experiment, the 12-task/3-arm pilot,
blind nonauthor grading, final HTML delivery or delivery-based 14/30-day follow-ups.
