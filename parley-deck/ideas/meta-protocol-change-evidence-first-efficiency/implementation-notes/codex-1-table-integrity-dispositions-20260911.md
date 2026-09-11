---
idea: meta-protocol-change-evidence-first-efficiency
author: codex-1
date: 2026-09-11
status: partial
source-manifest-sha256: 279ec8e3b32c6f32db9ca4eaed72aa5a82497ee46ac0b4e15a05618ff3210800
---

# Human evidence table integrity disposition

Claude's unchanged publication review identified a mutable human evidence table:
its bytes were excluded while only typed results and the remaining document were
bound. Three actual pre-fix counterexamples now reproduce closure with replaced
claims, inflated counts and a deleted table. The correction below is implemented
and tested; this is an implementer's disposition, not independent acceptance or
Phase 7 consensus. Original quorum, signatures and the full objective are intact.

## Current behavior

The managed table is a deterministic projection of original typed criterion
records. A visible caption distinguishes those executions from independent
acceptance. Unknown counts stay unknown, durations are numeric milliseconds, and
status comes from the typed record. Labels escape Markdown/HTML delimiters and
scrub credentials; diagnostics use a fence longer than embedded backtick runs.
JSON-compatible Unicode normalization before and after bounded scrubbing keeps
rendering stable across report persistence, including split multibyte characters.

Goldmark v1.8.6 is pinned in go.mod/go.sum for CommonMark section boundaries.
Actual top-level level-two ATX/setext headings are recognized; comments, HTML
blocks, nested list/blockquote content and fenced/indented examples remain bound
outside the managed section. A following level-one or level-two heading ends it.
Duplicate managed sections refuse. The old textual replacement helper and its
utility-only tests were replaced by production-parser coverage, including all
original example/subheading cases and more heading/context variants. Generation
also verifies that the rendered section is actually a distinct top-level section.

Each report binds exact section bytes under
`<relative IMPLEMENTATION.md path>#validation-evidence/v1`. Helper bindings before,
during and after independent execution, report publication, gate evaluation and
the exact final status-write path reconcile that hash and reconstruct the table
from the original records. Rehashing a forged table alone is insufficient.
Original report/request/receipt and retained independent execution reconciliation
stay in force. Old reports without the extra binding require fresh checks.

The non-evidence remainder is now pinned before commands and compared afterward;
a missing managed section is initialized before taking that snapshot. A command
that modifies implementation scope fails the cycle and leaves the original
binding in the retained report. No new file or section exclusion was added.
CRLF documents retain their endings. The exact independently authorized status
transition remains the only allowed completion change to the bound remainder.

## Executed evidence

All 344 Go/module files match the retained source manifest.
Full `go test -count=1 ./...` PASS (127.709s).
App/evidence/driver/runner race PASS (142.938s).
Scoped vet, module checksum verification and Windows amd64 app cross-build PASS.
Windows runtime remains unverified.

Isolated shared-volume projection tests PASS (6.615s) and actual process fixtures
PASS (12.204s). Real helper/criterion/parent processes refuse table changes before
the helper, during independent execution, after the helper, and after parent
acceptance immediately before Complete. They are local process fixtures, not
new model invocations. The full existing positive completion, negative scope,
refusal recovery, cancellation and publication suites remain passing.

Additional fixtures cover unchanged-rest table deletion, missing legacy binding,
forgery with a recomputed table hash, ambiguous headings, preserved examples,
escaped markup/fences, unknown counts, Unicode round-trip, CRLF completion,
missing-section initialization and command-induced non-evidence scope drift.

Manifest SHA256: `279ec8e3b32c6f32db9ca4eaed72aa5a82497ee46ac0b4e15a05618ff3210800`.
Full log SHA256: `f8d1a3349d57f4d9954b9747cc611db6f5f631fe11ad99d053eac8cfad269062`.
Race log SHA256: `0eeb84b7878b6762b8f98be6de95220a6c283f939bd3b77d2da9317fd9cba98f`.
Original three accepted-table-mutation failures:
`2458d16099c5037fda0f3d0d87e5b260d3d42f4a4050845f199f1f549ccf7136`.
Unicode split-boundary failure:
`2b02f1968762975b3282fa4fd21c35d2b8868f113fb7240d3a8f709849707f56`.
Logs, checksums and exact-source verification are retained under
`.parley-runtime/table-validation-20260911/` and the referenced sibling logs.

## Retained corrections and limits

The first integrated gate returned the new projection error before reporting
missing criterion scope. It now retains both typed and projection reasons;
partial scope is not hidden or weakened. The first Markdown matrix expected a
two-space heading after an open list to be top-level. CommonMark correctly kept
it inside that list; the fixture now ends the example explicitly before testing
top-level heading variants. Both failing logs are retained. The initial Unicode
sample ended on a character boundary and passed; the actual split-boundary case
failed before normalization. No nonexistent earlier failure is claimed.

These checks do not authenticate a same-UID actor able to fabricate every mutually
consistent artifact. The parent CLI enforcing closure must implement the new
binding; an old binary does not enforce an additional map key by itself. Scope
drift detection retains the original report binding, not an invented rollback of
a command's writes. Independent model verification and final acceptance remain
required. The new dependency is pinned and checked, not a global installation.

Safe lock-origin and changed-inactive-import recovery, durable semantic operation
replay, independently confirmed two-patch regression trajectory, full live launch
coverage/concurrency/closure, exact packet trials, the frozen 12-task solo/duo/
full-six pilot, equal ceilings/rotation, blind grading, owned reviews/signatures,
final populated report QA and actual elapsed follow-ups remain open. Historical
quorum/pilot/funding and Claude review recovery decisions remain unanswered.
No real model, operator accounting action, deployment, release, merge, global
install or core publication occurred. The six-part goal is incomplete.
