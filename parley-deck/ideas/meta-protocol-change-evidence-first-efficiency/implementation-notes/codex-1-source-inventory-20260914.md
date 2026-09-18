---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-14
source-base: ed8b9cbf09efb26ddf3c02f4a640e6873cde3a35
status: implemented-pending-independent-review
---

# Source inventory: local and global Git exclusions

Claude F4 identified a possible untracked-source omission. The original executed
counterexample is retained under .parley-runtime/claude-review-refutations-20260914/:
adding an untracked Go init file hidden by info/exclude kept Observe identical but
changed actual go run output from baseline to hidden plus baseline. That finding
is not withdrawn or signed off by this implementation note.

## Resulting behavior

GitSourceInventory uses Git's own --exclude-per-directory=.gitignore inventory and
requires exact agreement with --exclude-standard. Differences from info/exclude,
core.excludesFile and the default global ignore file now refuse. Evidence digest
and trajectory capture use this shared gate. Refusal carries no private path,
pattern, content or raw Git stderr. No private content is automatically archived
and no operator Git configuration is changed.

Tracked entries remain included even when a rule matches. Nested .gitignore,
negations and linked-worktree exclusions are resolved by Git itself. Project
ignored files remain outside scope; a project-ignored build input can still affect
execution. Untracked project rules are ordinary source when included in that
inventory, and Git can itself exclude such a rule file. This is not complete
build-input discovery or authentication against a same-user actor. Two inventory
commands are not an atomic source/configuration snapshot; trajectory repeats
source observation at its existing capture boundaries.

The checked inventory permits at most 100000 entries and 16 MiB output, with a
30-second subprocess context and caller cancellation. The output buffer is private
rather than embedded: a promoted bytes.Buffer.ReadFrom would let io.Copy bypass
the Write bound. The initial development boundary test caught exactly that bug
and the failing run remains retained (10.028 seconds). The corrected focused
run passed in 11.006 seconds. Neither development run substitutes for final
validation. File-system calls still depend on the host operating system.

Accepted source digest encoding is unchanged. Project-rule changes remain source
changes. Existing valid retained archives restore to the same digest without
creating Git metadata. Historical archives retain their original bytes; this
correction cannot reconstruct or certify source omitted by an old local exclusion.
An ambiguous current worktree fails even when the caller explicitly excludes the
hidden name as an evidence artifact. Unsupported types, links, tracked absences
and snapshot size limits retain their existing semantics.

## Verification scope

Seven new top-level tests cover five private-rule sources/shapes; project/nested
and untracked rules; negation precedence; tracked ignored files; actual Git
worktrees; deleted entries; NUL-delimited unusual names; output/entry bounds;
cancellation; diagnostic privacy; and unchanged Git index/configuration bytes.
The actual hidden-Go test executes both local and global cases, checks each
inventory/digest/Observe/capture gate separately and confirms no archive is
published. Explicit-scope restore exercises dirty tracked deletion and nested
negation against the actual restored filesystem.

Five separate removed-protection overlays fail at intended assertions: private
scope comparison, evidence binding, snapshot binding, output bound and entry
bound. Direct boundary assertions prevent a later rejecting gate from masking
the missing predicate. All overlays change native temporary files only.

Final validation binds 411 Go/module files, manifest SHA256
0c3b811a827ef84e04856ad79f1e2f6ad7a83884a8afa3dc59401f6f65f65949:

- Full Go suite PASS in 378.054 seconds; all 32 package terminals accounted for.
- Six-package race PASS in 407.207 seconds; vet PASS.
- Windows CLI production build and trajectory/app test cross-builds PASS, with
  PE amd64 headers checked. Windows runtime remains unverified.
- The evidence package's Windows test binary does not compile because its
  pre-existing unchanged tree_report_test.go directly references syscall.Mkfifo.
  The failed log is retained. The production Windows CLI includes the changed
  evidence package; this is not a claim that evidence tests passed on Windows.
- Compiled shared-volume evidence, source/snapshot trajectory, and production
  reservation-recovery app selections PASS in 2.134 / 16.075 / 11.470 seconds.
- Seven new top-level tests appear in full/race/shared evidence; all five intended
  negative failures, source stability and native/shared log bytes reconcile.

See .parley-runtime/source-inventory-final-validation-20260914/final-verification.json.
Historical report HTML bytes and Claude's participant-owned review remain intact.
This validates the stated inventory correction, not every independent probe below.

## Additional findings retained during validation

An actual charged quorum-changing patch reproduced F3: request freezing and
verification-ticket preparation accepted other while the activation archive still
named reviewer. A direct call of the unchanged scope routine also accepted later
altered membership. No second complete trajectory or live verifier was executed.
The native-only probe is .parley-runtime/claude-review-refutations-20260914/F3-counterexample.json
(log SHA256 ccd22bb68c33b841b584c94976ddee5e189004e95a60b71523dbc6029338c95b).
F3 remains open and is not hidden by F4's passing validation.

The F6 warm-read probe used a 16782848-byte archive and 0/1/3/5 actually charged,
reconciled unchanged fixture attempts. Native guarded Inspect reads rose from
44.886–46.919 ms at zero to 467.047–505.879 ms at five. One successful shared
metadata-diagnostic run measured 48.892–53.495 ms at zero and 550.792–651.134 ms
at five. These are one-machine observations during other validation, not a
worst-case bound, cold-cache benchmark, or proof of a 30-second deadline violation.

The original shared history probe FAILED, as did one of three diagnostic
repetitions. A complete hash-verified temporary snapshot retained its inode,
private mode and size but its mtime changed by 66835 ns at the final stat. The
production stability guard refused it; the first failure remains an unresolved
availability issue. Diagnostics modified only failure-message text in a native
overlay. No production protection was removed. Underlying host cause is not
proven. The other repetitions passed; they do not erase the failed ones. Exact
samples and the next investigation boundary are under
.parley-runtime/history-cost-probe-20260914/measurement-summary.json,
diagnostic-finding.json and next-snapshot-stability.md.

## Remaining obligations

This is an implementer-owned correction, not participant acceptance. F1 abnormal
terminal and consumed helper-ticket recovery, F3 activation-quorum binding, F6
guarded history cost and remaining workflow-effect recovery stay open. Fresh
participant-owned review/signatures, real launch/concurrency/closure coverage,
the frozen packet AB/BA experiment and twelve-task solo/duo/full-six pilot, final
offline HTML/ego-browser QA and actual-delivery 14/30-day follow-ups remain due.
No treatment call, quorum amendment, model-default change, merge, global install
or immutable protocol-core publication occurred in this source correction.
