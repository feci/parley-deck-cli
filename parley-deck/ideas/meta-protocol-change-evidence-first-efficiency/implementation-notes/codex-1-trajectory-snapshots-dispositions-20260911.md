---
agent: codex-1
idea: meta-protocol-change-evidence-first-efficiency
phase: implementation
date: 2026-09-11
---

# Reconstructible trajectory source snapshots

## Outcome and remaining acceptance

Opted-in trajectory capture now retains exact source archives at activation and
terminal publication. The original charge binds the before-archive, and the real
instrumented invocation binds its actual post-source and after-archive. Dirty
modifications, tracked deletions, untracked additions, executable permissions,
binary contents and supported relative symlinks can be reconstructed after the
live files change again. Restore never targets the live worktree or Git metadata.
This completes the source-retention slice of AC-B2, not independent acceptance
or the complete trajectory workflow. Every captured attempt remains pending.

## Implementation and boundaries

Private archives are content-addressed under shared cycle accounting in
trajectory-snapshots. The exact versioned reference pins archive SHA256/length;
the canonical first member pins the original Source observation. Remaining
members are sorted, bounded regular files or relative links. The complete tar
byte stream and evidence.TreeDigest source serialization are checked before
publication and again at inspection/restoration. Re-encoding enforces the exact
archive format including the footer; Go archive/tar's permissive EOF handling
cannot turn missing footer blocks or normalized extension headers into acceptance.

Source inventory comes from actual Git tracked/untracked names. Observed tracked
deletions are explicit exclusions only in trajectory observation; ordinary
completion evidence.TreeDigest remains unchanged. Capture brackets its source
observation and archive write with drift checks. Git observation passes
--no-optional-locks, so read-only status cannot refresh the user's index.
The archive, individual file and inventory limits are 256 MiB, 64 MiB and
100,000 entries. Unsupported or excessive source refuses, never disappears.

Archive files/directories must be private and real. Existing content addresses
are rechecked and never replaced to hide corrupt evidence. Synced publication
uses the existing resource guard and file replacement. Metadata references
contain no raw source, commands, environment or paths; tar bodies contain private
source and must stay out of public telemetry, protocol files and reports.

Restoration uses os.Root confinement in a newly allocated private directory.
Regular files precede links. Component-by-component link resolution preserves
link/.. semantics and rejects missing targets, cycles, parent conflicts and
escape. Duplicate/unordered paths, traversal, Git metadata, unsupported member
or header types, invalid permissions, oversized files, truncation and trailing
bytes refuse. The restored filesystem must independently reproduce the original
evidence.TreeDigest. On failure, only the newly created restore directory is
removed. Existing parent contents and live files remain untouched.

Restoration creates no Git metadata and makes no original-commit/ancestry claim.
The caller supplies a parent outside an existing Git worktree. Host permissions
and symlink semantics must reproduce the archived digest exactly, otherwise the
restore fails. Independent verifier execution still needs separately allocated
Git roots with explicit original-source-to-execution-source bindings; archived
dirty output must not be relabeled as an original clean Git commit.

Persisted trajectory policy/state is now v2. Baseline and before-archive references
are mandatory, and a successful post-source capture needs its exact after-archive.
Unavailable source records source-unavailable; an identifiable source whose
archive cannot be retained records archive-unavailable with the actual terminal
and source observation. Archive loss/corruption refuses inspection, new charges
and both existing completion paths. Legacy v1 digest-only state is refused with
explicit recovery guidance and remains byte-for-byte intact. Current live files
cannot reconstruct historical attempts; no automatic migration is supplied.

These are cooperative runtime controls, not human authentication or resistance
to a same-UID actor consistently fabricating every authority. Snapshot hashing
alone does not establish independent model execution or grant continuation.

## Executed refutation and validation

The first integration run compiled but failed during activation because the
archive source serializer added a newline after regular-file bodies, unlike
existing evidence.TreeDigest. The newline was removed. Round-trip tests now
compare actual restored files through evidence.TreeDigest, including binary
contents and dirty tracked deletions, rather than comparing duplicate literals.

The first new archive fixtures used testing.T.TempDir directories whose mode was
not private enough; explicit 0700 fixture directories corrected that setup.
Negative results from the affected run are not credited as archive-path coverage.
A subsequent malformed-PAX fixture initially requested USTAR with PAX fields;
its explicit PAX format corrected test construction. Retained logs preserve both
fixture errors and the preceding serialization failure. These failures were not
independent acceptance or real model outcomes.

Nine new snapshot tests cover clean/dirty restoration after live edits, unchanged
index/HEAD/refs, long paths/links, supported and refused links, 22 malformed-archive
cases, corruption/loss/publication refusal, oversized source, component-resolution
attacks, simultaneous capture and an actual source change during archive writing.
Three new state tests cover missing/corrupt/rebound/v1 baselines before charging,
retained terminal after archive failure, and lost after-archives. Existing runner
fixtures now restore actual manual/grouped nonzero and interrupted child output
after another live edit. CLI inspection exposes the archive binding and refuses
loss. These are synthetic executable fixtures; no model call or actual operator
accounting/activation action occurred.

Current validation is bound to the exact 366 Go/module files in
.parley-runtime/trajectory-snapshots-validation-20260911/source-manifest.json,
SHA256 7675f7dc655644e7e8d3f1c1e24da57a5c8f52caf2d943c26c6fad11e7eb3fd2.
Full JSON suite passed in 167.333s with all 32 expected package terminal events;
scoped trajectory/evidence/budget/app/driver/runner race passed in 183.107s with
all six package passes. All twelve new and three amended top-level tests have
passed events in both runs. Native/shared logs match exactly, contain no NUL
bytes and have no failed-test events. Vet passed in 0.865s; Windows amd64
trajectory/app cross-builds passed. Compiled shared-volume trajectory, runner,
driver and app fixtures passed in 63.183/13.339/2.332/4.923s respectively.
final-verification.json verifies source identity, complete package/test coverage
and native/shared logs; checksums.json binds forty evidence files. Windows runtime
and independent model acceptance remain unverified.

## Next required integration

Bind each exact retained charge and source archive pair to an independently
selected verifier invocation through the instrumented runner. Prepare isolated
execution roots without altering original-source attribution. The selected agent
must actually invoke the matching helper; parent acceptance must pin invocation,
process marker, request, roots, criteria and retained helper receipt, including
partial failures. Persist observations and derive the expected patch list from
the complete charge inventory before calling the existing AB/BA evaluator.
Retain the first two-confirmation review trigger, implement explicit durable
review/recovery decisions, and enforce them across resume/manual/driver paths.
No reset, proxy-authored acceptance, silent budget reset or quorum reduction.

Safe guard/lock-origin migration, durable semantic action identity/replay, fresh
independent source acceptance, full live launch/concurrency evidence, the exact
packet experiment, twelve-task solo/duo/full-six pilot, blind nonauthor graders,
participant-owned reviews/signatures and actual-delivery follow-ups remain open.
The reported Claude weekly-limit recovery, historical roster/pilot amendment and
experiment ceilings remain unanswered. Ordinary implementation continued without
retrying or replacing the failed participant.

The HTML remains at bd0018865946764cd65cd5cde7860bb5bfd480c9 and does not represent
these source changes. Model inventory remains 34 terminal attempts, 17 unknown
costs, USD 46.1887585 known CLI estimates and unknown total. The preceding source
checkpoint's sparse broad-test log anomaly remains retained and unclassified;
this checkpoint uses complete pipe capture, native retention and complete-byte
shared publication. No final merge/release/deployment/global install/immutable-core
publication or shared-memory write is claimed. The six-area goal stays open.
