---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
round: 1
date: 2026-08-07
---

## Summary

Use one immutable, versioned, content-addressed core under ~/.parley, but do not replace one
uncontrolled Markdown-copy problem with a general patch language. The smallest overlay that
satisfies the decided override-and-extension constraint has only two operations: replace one
stable named block that the core explicitly marks replaceable, or append a namespaced block at
one explicit extension point. Everything is sealed by default. The fleet evidence strongly
supports this restrictive default: almost all observed differences were six generated identity
zones, and the only genuine local section was itself sync governance rather than a project rule
(parley-deck/ideas/meta-protocol-change-global-core-protocol/00-prompt.md:12-31).

Each idea pins a tuple of core version and hash, overlay hash, resolver version, and effective
protocol hash. The exact rendered Markdown it uses is stored as an immutable, content-addressed
deck artifact; later phases read that snapshot, not whatever ~/.parley or the deck's current
COOPERATION.md contains. This follows the roster precedent: the run manifest freezes the effective
launch identity because re-resolution during continuation was unsafe
(internal/runmanifest/manifest.go:49-98).

Keep a committed COOPERATION.md, but make it a generated current-view artifact, never an
authority. It gives humans a conventional entry point and makes a deck portable; per-idea
snapshots preserve simultaneous old and new protocol versions. A committed deck lock containing
version plus content hashes makes two machines deterministic: a machine missing the exact pinned
core blocks rather than substituting its own same-named or current version.

The security claim must be narrow and honest. A process with the user's unrestricted OS identity
cannot be distinguished from the user by Unix ownership, mode bits, or a CLI convention. Real
prevention is possible for Parley-managed agent processes only when the launcher applies an
OS-enforced sandbox that denies writes to the global store and all children inherit that denial.
For an arbitrary unsandboxed same-user agent, the achievable guarantee is integrity detection and,
for changes made through the attended publisher, attribution. A raw direct write is detectable
but not reliably attributable.

## Proposed approach

### 1. Global authority and release model

The authority is an immutable release store under ~/.parley. A logical release contains:

- a unique version;
- the exact core Markdown;
- a stable registry of named block IDs and whether each is sealed, replaceable, or an extension
  point;
- the core SHA-256 and registry SHA-256;
- a release receipt linking the prior version and recording how the user published it.

The serialization of that logical record should be chosen only after the resolution semantics are
agreed, consistent with the prompt's non-goal
(parley-deck/ideas/meta-protocol-change-global-core-protocol/00-prompt.md:73-78).

Released bytes are never edited in place. The user may edit a staging copy and publish it as a new
version. If the user directly edits an already released version, its manifest hash no longer
matches; Parley treats it as damaged or unregistered content and will not use it for a new idea
until the user publishes those bytes as another version. A current-version pointer is convenience,
not sufficient identity: every consumer uses version plus hash.

This separates the same two concerns that LoadRosterScoped separates for the roster: the global
store supplies shared content, while the committed deck record decides what that deck has adopted.
The roster implementation already prevents a machine-local membership set from silently becoming
committed deck state (internal/config/runtime.go:90-118,174-205).

### 2. Minimal resolution model

Do not patch Markdown by heading text, line range, JSON Patch, or unified diff. Those are fragile
under heading renames and make conflict meaning depend on incidental formatting. Core Markdown
instead has stable block IDs independent of its headings. A block is the smallest replaceable
unit, and replacement is whole-block only.

Resolution is deterministic:

1. Load the exact core version and hash selected by the committed deck lock.
2. Verify its release receipt and content hashes.
3. Load either no overlay or one committed deck overlay and hash it.
4. For each replace operation, require an exact target block ID, require that the core marks it
   replaceable, and require the overlay's expected base-block hash to match.
5. For each extension, require a deck-namespaced unique ID and the one explicit extension-point
   target. Extensions cannot redefine core IDs.
6. Inject the six deck data zones (workspace, creation date, transport, sync stamp, roster table,
   and host-handle table). These are renderer inputs, not overrides.
7. Validate required and sealed blocks, render one flat Markdown document, and hash it.

Precedence is therefore not an unrestricted "last writer wins." Generated identity data fills
typed slots; an allowed replacement wins only within its named open block; an extension is
appended only at its declared point; sealed core rules always remain authoritative. If an
extension semantically contradicts a sealed rule, the overlay is incompatible rather than a
higher-precedence exception.

The initial core should expose only:

- one replaceable working-language block, because the current core already explicitly permits a
  project to make that deliberate choice (parley-deck/COOPERATION.md:742-743); and
- one append-only project-extension point for genuinely project-specific rules.

All other normative content starts sealed: track selection and phases, quorum, artifact and
signoff shapes, ownership/conflict rules, the global-change rule itself, the automated-loop human
brake, and verification integrity. Transport choice, repository paths, branch names, roster
membership, and host handles are configuration/data, not prose overrides. Opening another block
later is a global core-version decision. This fulfills override and extension without designing
an inheritance framework for hypothetical demand.

Each overlay operation must record a rationale, owner, target, expected base hash, and a short
statement of which sealed rules it depends on. This is an exceptions ledger, not a partial fork.

### 3. Materialization, pinning, and multi-machine safety

The deck carries a committed logical protocol lock with:

- adopted core version and SHA-256;
- overlay SHA-256, or none;
- resolver/schema version;
- effective SHA-256;
- path to the immutable effective snapshot;
- compatibility-receipt SHA-256 when an overlay exists.

At Phase 0, the same fields are copied into 00-prompt.md frontmatter and into run state. The full
effective Markdown is written once under a content-addressed snapshot name and committed; equal
effective hashes deduplicate. Every phase, manual participant prompt, continuation, inspection,
and review reads the idea's pinned snapshot. The current protocol must never be re-resolved for an
open idea. The run-manifest precedent already hashes a sorted effective snapshot and includes
launch arguments so a value-only change cannot masquerade as current
(internal/runmanifest/manifest.go:63-98).

The next idea resolves the then-adopted new core and gets a new pin. Updating the deck's current
view does not affect the open idea. This permits two ideas at different versions in one session,
as the decided constraint requires
(parley-deck/ideas/meta-protocol-change-global-core-protocol/00-prompt.md:38-46).

COOPERATION.md remains committed and generated. Its provenance banner names the core version and
hash, overlay hash, effective hash, and resolver version. It is the current default view for
humans, not the source used by an already open idea. Rendering must be preview-first, atomic, and
idempotent, paralleling the roster renderer's generated-view and byte-stability rules
(internal/app/roster_render.go:14-28,104-157).

For two machines, the committed lock wins over each machine's current pointer. If machine B lacks
the exact core hash, it must install that release or stop; it may not render from a different
same-version file. Adoption of a newer global version is a compare-and-swap update of the deck
lock, generated COOPERATION.md, compatibility receipt, and new effective snapshot. This is the
minimum needed to make per-machine ~/.parley safe. A semantic version without a content hash is
not sufficient.

### 4. Compatibility checks and overlay drift control

No overlay is created at bootstrap or migration unless reviewed local content actually exists.
An empty/default overlay is forbidden: absence is the canonical no-customization state. This
prevents 36 mechanically generated overlays from becoming the next 36 stale copies.

Every new core version invalidates a nonempty overlay's prior compatibility receipt. Before the
deck adopts that core, the checker must:

- verify every target ID still exists and has the required open mode;
- compare every expected base-block hash;
- identify every changed core block named in the overlay's dependency declarations;
- reject duplicate IDs, an extension outside the extension point, and any attempt to replace or
  omit a sealed block;
- compose the result and validate the required-block registry and effective hash;
- produce a reviewable report of changed targets, changed dependencies, and unresolved semantic
  questions.

Mechanical validity is necessary but not sufficient because prose contradiction is not generally
decidable. For any nonempty overlay, a new core therefore also requires explicit user acceptance
of that report, recorded in a committed compatibility receipt keyed by core hash, overlay hash,
and resolver version. The burden is acceptable precisely because evidence says genuine local
authoring is near-zero. Compatibility must be argued from located changes rather than participant
count, consistent with the protocol's evidence and conflict rules
(parley-deck/COOPERATION.md:1242-1281).

If compatibility is missing or fails, existing ideas continue on their pinned snapshots, but the
next idea is blocked. The tool must not silently keep the next idea on the old core, discard the
overlay, quarantine it out of the effective protocol, or auto-migrate its prose. The user chooses
one of two resolutions: ratify an updated/removed local overlay through the local protocol-change
path, or decline to start the new idea. This preserves both the version-pinning rule and the
mandatory compatibility check.

A fleet audit should report only: decks with overlays, their operation count, targeted block IDs,
last compatible core hash, and status. Repeated local rules are candidates for a global proposal,
never auto-promoted. The existing anti-drift test is a useful negative precedent: it normalizes a
small allowlist and compares everything else byte-for-byte
(internal/protocol/drift_test.go:14-22,98-129). The new guard should compare generated output and
hashes, not maintain another widening prose allowlist.

### 5. Enforcement of global-core immutability

There are three distinct guarantees:

**Managed prevention.** Parley-launched participants run under an OS sandbox/profile whose write
allowlist excludes the global protocol release store and its parent metadata. The denial applies
to child processes. Preflight inspects the effective confinement policy; if exclusion cannot be
proved, the run must not claim that global-core writes are prevented. A headless protocol publish
operation is unavailable from an agent run.

**Attended publication.** The only supported mutation surface is a separate attended user command
that requires a controlling UI/TTY confirmation, consumes a finalized candidate, writes a new
version rather than modifying one, and emits the release receipt. Environment flags and a CLI
refusal are defense in depth, not the security boundary.

**Detection and attribution.** Every read and every pre-idea check recomputes release hashes. A
mismatch blocks new ideas and reports the exact changed path and expected/actual hashes. A valid
publisher receipt attributes a supported publication to the user action and source candidate. A
direct raw write is reported as DETECTED-UNATTRIBUTED. It must not be blamed on an agent without
evidence.

Read-only mode bits help prevent accidents only. They do not enforce the user/agent distinction:
the same owning UID can normally chmod them or write through another process. Absolute prevention
would require a different OS principal, a mandatory sandbox, or a broker holding a capability the
agent cannot access. Managed sandboxing is defeated by launching the agent outside it, a sandbox
escape, or an allowed unsandboxed helper that will write on the agent's behalf. The design should
state these defeat conditions rather than make an unqualified security claim.

### 6. Migration of the 36 flat copies

Migration is an attended fleet operation after the resolver, lock, renderer, pinning, and
compatibility checker exist:

1. Inventory exact deck roots, file hashes, protocol/version metadata, writability, and all open
   ideas. Freeze the source core release hash.
2. Normalize only the six evidenced identity zones and classify every flat copy against known
   ratified cores. An exact normalized match becomes no overlay.
3. For an open legacy idea, create a content-addressed snapshot of the exact flat protocol it
   started under and backfill a legacy version/hash pin. Do not pretend it started under the new
   global version.
4. Treat any remaining prose difference as review-required. Do not automatically translate it
   into an override or extension. In particular, review the one preserved sync-governance section
   as a candidate global rule or a retired workaround; do not resurrect it as a local overlay by
   default.
5. For each clean writable deck, extract identity data, write the lock and snapshot, generate
   COOPERATION.md, run compatibility and hash checks, then atomically apply.
6. Report every deck as applied, unchanged, skipped-read-only, skipped-unreachable,
   skipped-unrecognized, or failed-and-restored.

Use a full dry-run, compare-and-swap on every input hash, file-level backups with verified restore,
small-batch confirmation, per-deck atomicity, resumability, and a machine-readable final report.
Those controls are already the ratified roster-migration precedent
(parley-deck/ideas/roster-operations-standard/FINAL.md:127-131;
parley-deck/ideas/roster-operations-standard/consensus.md:447-468).

Read-only or unreachable decks remain byte-for-byte unchanged and are not marked migrated.
Existing open work can continue only from a preserved local flat copy/snapshot; a new idea blocks
until migration can create a valid lock and pin. A partial fleet operation is a reported state,
not grounds to weaken the new model.

### 7. Change to §7

Section 7 must stop saying that a drafting agent updates COOPERATION.md in place; that is its
current rule (parley-deck/COOPERATION.md:745-758). Replace it with two paths:

- **Deck-local protocol change:** a meta-protocol-change idea may propose and, after normal
  ratification, change only the deck overlay. It cannot write ~/.parley. Its output includes the
  new overlay hash and compatibility result against the adopted core.
- **Global-core proposal:** a full deliberation produces a candidate core diff, migration/impact
  report, and proposed new version. Agents write that candidate only inside the workspace. After
  consensus and human review, the user alone invokes the attended publisher to create and
  activate the new immutable global version.

Publication never rewrites a released version. Existing ideas remain pinned. At the next Phase 0,
each deck runs compatibility, updates its lock/current generated view if valid, and pins the new
idea. A global change therefore becomes available everywhere but does not silently mutate open
work or write into 36 repositories.

## Concerns / open questions

- The exact marker and manifest serialization is intentionally open; the required semantics are
  stable IDs, whole-block replacement, one extension point, sealed-by-default, and content hashes.
- A managed sandbox must be proven for every supported adapter and platform. If one adapter can
  invoke an unsandboxed desktop helper or connector, the prevention claim must be downgraded to
  detection for that launch.
- Requiring explicit user compatibility acceptance for every nonempty overlay/core pair is
  intentionally conservative. If real fleet evidence later shows many legitimate overlays, a
  future core version may introduce declared dependency scopes that allow some no-impact checks
  to auto-pass; the initial design should not assume that demand.
- The attended publisher's user-presence mechanism and release-receipt signing are implementation
  choices. A plain TTY prompt improves safety but is not cryptographic proof of a human.
- The initial replaceable working-language block is supported by an existing explicit exception.
  If participants find no other evidenced override target, they should not open one merely to make
  the mechanism look more general.

## Risks

- **False security language:** same-UID unsandboxed agents can bypass file modes and CLI policy.
  Mitigation: scope prevention to verified managed sandboxes; otherwise report detection only.
- **Semantic conflicts:** structural composition cannot prove that an extension does not weaken a
  sealed rule. Mitigation: sealed core precedence, explicit dependency declarations, changed-block
  report, and per-version user compatibility acceptance.
- **Snapshot growth:** committed effective snapshots duplicate core prose over time. Mitigation:
  content-addressed deduplication and one snapshot per effective hash. The storage cost buys
  simultaneous-version correctness and historical portability.
- **Cross-machine races:** two machines may try to adopt different current cores. Mitigation:
  committed version-plus-hash lock, compare-and-swap, and refusal to substitute a local version.
- **Overlay proliferation:** convenient initialization could recreate fleet drift. Mitigation:
  no overlay by default, no automatic conversion of unknown prose, explicit operation metadata,
  and fleet reporting.
- **Migration ambiguity:** unknown flat copies may not map cleanly to a known core. Mitigation:
  skip and report rather than guessing; preserve exact legacy snapshots for open ideas.
