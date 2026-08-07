---
idea: meta-protocol-change-global-core-protocol
drafted-by: claude-1
date: 2026-08-07
rounds: 2
participants: [claude-1, codex-1, hermes-1, kimi-1]
absent: [opencode-1]
status: revision 2 — awaiting signoffs (rev-2 signoffs in signoffs/rev2/)
---

# Consensus — global core protocol with local override/extension

**Revision 2.** Revision 1: hermes-1 ACCEPT; codex-1 and kimi-1 BLOCK. Both blocks are upheld in
full — see `## Drafter position changes` for what reversed and why.

## 0. Provenance and role concentration (§15)

`PRIMARY` unless tagged. claude-1 is facilitator, drafter, and a participant — three roles in one.
The reviewers' positions decide; where this draft states a fact about code or system behavior it
names the command or `file:line` that establishes it.

**opencode-1 is absent, not agreeing.** It was invoked four times and produced no round-01
artifact; the attempts, what was ruled out, and one facilitator error among them are recorded in
`inbox/claude-1-to-all_meta-protocol-change-global-core-protocol_opencode-timeout.md`. It is not
counted in quorum and its silence must not be read as assent.

## 1. Verdict conflicts (§15.3)

### VC-1 — Is preventing an agent from writing the global core achievable?

> **claude-1** (`round-01/claude-1.md`), tag `PRIMARY` at the time:
>
> An agent runs with the user's own privileges. It can `chmod`, it can write the file with any
> tool. File mode `0444` … stops an *accident*; it does not stop an agent that decides to. … the
> guarantee is not "cannot", it is "cannot do so silently".

hermes-1 and kimi-1 took the same position in round 1.

> **codex-1** (`round-01/codex-1.md`), tag `PRIMARY`:
>
> Real prevention is possible for Parley-managed agent processes only when the launcher applies an
> OS-enforced sandbox that denies writes to the global store and all children inherit that denial.

**Resolution: codex-1 is correct. The three-way majority was wrong, and was not resolved by
counting.** claude-1 tested the claim instead of arguing it (`round-02/claude-1.md`):

```
$ sandbox-exec -f /tmp/sb.sb -D GLOBAL=/private/tmp/fakeparley \
    /bin/sh -c 'echo tampered > /private/tmp/fakeparley/COOPERATION.md'
/bin/sh: /private/tmp/fakeparley/COOPERATION.md: Operation not permitted     # denied
… child process:                                                             # denied, inherited
… rm -f …:                                                                   # denied, file survives
```

A macOS seatbelt profile denies the write, the denial is inherited by children, and the process
cannot escape by deleting the file. Prevention is real for processes parley launches under such a
profile.

**Two scope limits are part of the resolution**, not caveats to it:

1. **The facilitator is outside it.** In this idea the facilitator is a Claude Code process parley
   did not launch — and it is the agent most likely to touch the protocol, because it runs the
   migrations. Prevention does not cover it.
2. **An unresolved path silently disables the profile.** The identical profile built from `/tmp`
   (a symlink to `/private/tmp`) permitted the write with no error. `~/.parley` on a machine whose
   `$HOME` traverses a symlink is exactly that case. A configured sandbox is therefore not evidence
   that anything is confined.

Adopted wording: **prevention for parley-launched participants, detection and attribution for the
facilitator and any unmanaged agent, and no claim of either without a runtime proof.**

## 2. Decisions

### D1 — The core is an immutable, versioned, content-addressed release store in `~/.parley`

`~/.parley/protocol/core/<version>/` holds the exact core Markdown plus its registry, both hashed.
Releases are **write-once**: a release is never edited in place, so a change by the user is a new
version *by construction* rather than by discipline (kimi-1's shape, codex-1's release semantics,
adopted by all four).

### D2 — A registry of permanent block IDs, sealed by default

The release carries a registry mapping permanent, never-reused section IDs (`s1`…`s15`, subsections
such as `s6.6`) to `sealed | replaceable | extension-point`. Deleted IDs get tombstones. Addressing
is by ID, never by heading text: a renumbered or renamed heading must not silently orphan an
override. No inline markup in the body (kimi-1; supersedes claude-1's round-1 HTML-comment markers).

### D3 — v1 open surface is deliberately near-empty

- **replaceable:** `s6.6` (working language) — the one carve-out the protocol already contemplates.
- **extension-point:** exactly one, `ext-1`, rendered at a declared position, deck-namespaced IDs.
- **identity slots:** the six measured zones (`**Workspace:**`, `**Created:**`, `**Transport:**`,
  the `**Protocol synced:**` stamp, the §2 roster table body, the host-handle table body).
- **everything else: sealed**, explicitly including §7 and §15 — a deck must not be able to weaken
  the rules governing how rules change or how claims are verified.

Basis: across 36 decks only one genuine local section existed, and it was sync governance that
belongs in the core (`00-prompt.md:21-24`). kimi-1 retracted a wider five-section proposal on this
evidence.

### D4 — One overlay file per deck, two operations, no empty overlays

`parley-deck/protocol-overlay.md`, committed. Operations: **replace a replaceable block by ID** and
**extend at `ext-1`**. Each operation records a rationale and the **expected hash of the core block
it replaces**. Unknown ID, duplicate provider, sealed target, or a changed base hash **fails
closed**. The *absence* of the file is the only canonical "no customization" state; an empty
overlay file is forbidden, so the fleet cannot accumulate 36 vacuous overlays (codex-1's rule,
adopted by kimi-1 and claude-1).

### D5 — Resolution order is core-owned, not last-writer-wins

1. load and verify the exact core release named by the deck lock;
2. fill the six identity slots;
3. apply the replace operation if provided;
4. append the single `ext-1` payload;
5. validate that all required sealed blocks are present; hash the effective bytes.

Core semantics own each slot's mode. A contradictory extension is *incompatible*, not
higher-precedence (codex-1).

### D6 — The deck keeps a committed, generated `COOPERATION.md`

It is a deterministic **current view**, never the authority. Rationale, agreed by all four:
auditability (the protocol is evidence; a historical idea must stay interpretable) and
self-containment (a checkout without the skill still has a readable protocol; three of four
participants read the deck file directly). `protocol render` regenerates it and **reports every
block it replaces or removes**, as `roster render` does; `protocol check` verifies it and reports a
hand-edit rather than silently overwriting.

### D7 — Per-idea pinning, with the effective snapshot ALWAYS materialized

`00-prompt.md` and the run manifest record core version, core hash, overlay hash, resolver version,
and the effective hash — parallel to `RosterSnapshot`/`RosterRevision`
(`internal/runmanifest/manifest.go:49-58`).

**The rendered body is ALWAYS written**, content-addressed and deduplicated per deck, at
`parley-deck/protocol-snapshots/<effective-sha>.md`, before Phase 0 closes. It is then the **sole
protocol input for every later phase of that idea**.

*Revision 2 — this reverses the drafter's own position.* Revision 1 stored the body only "when the
render is not reproducible from those inputs", on the drafter's round-2 argument that duplicating
~80 KB "buys nothing". codex-1 and kimi-1 blocked it independently and both are right: the inputs
can stop being present. Prune a release, migrate to a new laptop, or clone the repo fresh, and the
open idea's governing bytes are unrecoverable — the body was never stored, the release is gone, D8
forbids substituting another, and "materialize it then" is impossible because there is nothing left
to render from. That fails **user constraint 1** in an ordinary new-laptop scenario. Storage was
traded against the guarantee the pin exists to provide.

Missing or tampered snapshot **blocks continuation**. A missing global release blocks *adoption and
rendering* only — an already-pinned idea continues from its snapshot without the release present.

`sessions inspect` reports `stale-protocol`. An open idea completes under its pinned version; the
next idea in that deck picks up the current one (user constraint 1).

### D8 — Deck lock; a missing pinned release BLOCKS

The deck commits a lock with core version + hash, overlay hash or `none`, resolver version,
effective hash. A machine that does not have the exact pinned release **blocks** rather than
substituting its own same-named or current version (codex-1). This is what makes the per-machine
`~/.parley` safe: divergence surfaces as a refusal and as a diff in the committed render, not as
two machines quietly running different protocols.

### D9 — Enforcement, stated per launch path, never claimed without proof

- **parley-launched participants:** an OS sandbox denying writes to the core subtree, built from
  the **resolved** path. Preflight **proves it** by attempting a real write to a sentinel inside the
  protected subtree; if that write succeeds, the run reports `confinement-unproven` and must not
  claim prevention. This is the fail-closed pattern already used for the AUTO bit.
- **facilitator and unmanaged agents:** **detection**, and attribution only where the change came
  through the supported path. A change made by the attended, TTY-gated publisher is attributable
  because it emits a release receipt. An unexplained hash mismatch is `DETECTED-UNATTRIBUTED` — we
  know the core changed, not who changed it. Revision 1 said "detection and attribution" flatly;
  codex-1 blocked that as an overclaim and is right. Mechanisms: hash verification against the
  declared version, `0444` as a second layer and accident-stopper, no agent-accessible write path.
- **Named limits on any prevention claim:** the preflight probe proves **direct-write** denial
  only. Delegation paths (IPC brokers, helper processes) and inherited writable file descriptors
  are NOT covered by it; a claim of prevention must say so rather than imply completeness
  (codex-1, kimi-1).
- **platforms without the primitive** (no `sandbox-exec`): `confinement-unproven` and
  detection-only. Never a silent downgrade.

Diagnosis must keep working on a tampered core: block *launching participants*, never block
`protocol check`, `roster show`, or reading the committed file (kimi-1's rule, narrowed by
claude-1).

### D10 — Compatibility checking on core adoption

Registry-driven, at adoption time. For a **replace** operation: the target ID must still exist,
still be replaceable, and the expected base-block hash must still match. For an **extension**:
the overlay declares the core blocks it depends on (defaulting to **all sealed blocks**), and any
change among them produces a **reviewable change report requiring reconfirmation** — an extension
written against §7 as it was must not silently ride along after §7 changes. Revision 1 defined no
check for extensions at all (codex-1). The overlay also declares the **core version range** it was
written against, without which the check has no baseline (claude-1's round-2 open question, now
resolved). Mismatch → re-confirm. Missing, tombstoned, or
now-sealed target → **block**. Auto-pass only on a zero-change report. **Never auto-migrate prose.**
An incompatible deck stays pinned on its old core — that pinning *is* the quarantine.

### D11 — §7 gains a blast-radius clause

A **core** change now hits every project at once, so it requires the meta idea *and* explicit user
ratification — the user is the only party permitted to change the global (user constraint 4). A
**deck overlay** change is a smaller act and is allowed through a normal idea in that deck.

## 3. Implementation scope for this cycle, ranked

The ranking below is the synthesis closest to **codex-1's** staging, articulated most explicitly by
hermes-1 and accepted by all four. (Revision 1 credited it to hermes-1 alone — kimi-1's correction,
upheld.)

1. **Core store + `protocol render` + `protocol check` + drift guard.** Converts 36 hand-edited
   copies into 36 generated ones — the bulk of the drift reduction for the least mechanism.
   **The renderer is a NEW pure function** taking (core release, overlay, identity slots) and
   returning bytes, with the synced-stamp derived from the deck lock. `mergePreservingZones`
   (`internal/app/preflight.go:527`) is **zone-extraction scaffolding only** — reused to pull the six
   identity zones out of a legacy deck during migration, never as the renderer. Its `## 3.` heading
   anchor (`preflight.go:539`) does not survive: addressing is by registry ID (kimi-1).
2. **Per-idea pinning + `stale-protocol`** (D7, D8).
3. **The overlay** (D4, D5, D10) — third on purpose: shipping override machinery before the
   generator would build for a case that barely exists.
4. **Detection-layer enforcement** (D9 minus the sandbox): hash verification, `0444`, no write
   path, attended publish.

## 4. Deferred follow-ups

- **DF-1 — the OS sandbox itself** (the prevention half of D9). It is a launcher change, not a
  config edit, and is platform-specific; Linux has no `sandbox-exec`. Until it ships, every surface
  reports `confinement-unproven` and the protocol says detection-only. Ratified as designed.
  **Scope, stated rather than implied** (kimi-1's condition): the preflight probe proves
  **direct-write denial only**. Extending the profile and a conformance suite to delegation paths
  (IPC brokers, helper processes) and inherited writable file descriptors is part of DF-1, not of
  the shipped claim.
- **DF-2 — fleet migration of the 36 flat copies.** A separate attended operation following the
  roster Stage-4 discipline (per-deck dry-run, diff, confirm, rollback). hermes-1 is explicit about
  why: the 2026-08-06 sync destroyed `librade-algoTrader`'s genuine local section, and careless
  migration repeats that. Tooling from items 1–4 must exist and be tested first.
- **DF-3 — opencode-1's fitness as a quorum member.** It completes short tasks and stalls on long
  analytical ones. Needs a timeout/context/invocation review before it is relied on.
- **DF-4 — restoring `librade-algoTrader`'s destroyed section** as an `ext-1` payload once the
  overlay ships. It is preserved in `.parley-protosync-backup-2026-08-06/`.

## 5. Gates (binding on implementation)

- **G1** — `protocol render` MUST be idempotent (byte-identical on a second run) and MUST report
  every block it replaces or removes, in preview and on apply.
- **G2** — no **autonomous or agent-accessible** code path may write into
  `~/.parley/protocol/core/`. The attended, TTY-gated publisher is the sole audited exception, and
  it creates a NEW release rather than modifying one. A guard test asserts both halves. (Revision 1
  said "no code path", which contradicted D9's own publisher — codex-1.)
- **G3** — no surface may report confinement without the preflight write-probe passing; the probe
  MUST use the resolved path.
- **G4** — an overlay targeting a sealed, missing, or tombstoned ID, or whose expected base hash
  does not match, MUST fail closed; an empty overlay file MUST be rejected.
- **G5** — an idea's pinned protocol MUST be what its later phases read; an acceptance test MUST
  prove that changing the core mid-idea does not change what the open idea resolves.
- **G6** — this idea's own protocol change MUST be recorded in `meta/protocol-changelog.md` in the
  §7 template format (`Idea:` path, `Drafted by:`, `Summary:`).
- **G7 — production call-site pin test** (codex-1, verbatim intent). Start an idea under effective
  protocol A, adopt B, then capture the actual prompts/inputs produced by the production entry
  points for round 2, design consensus/signoff, implementation, review, review consensus/signoff,
  fix-up, resume/continue, steer, and inspect. Every entry point MUST resolve and expose snapshot A
  with the recorded hash; none may read the deck-current `COOPERATION.md`, the current core pointer,
  or core B. Instrument filesystem reads (or inject a resolver spy) so any forbidden read fails the
  test. Run the same test **after deleting the global A release**; continuation MUST still succeed
  from snapshot A.
- **G7b — call-site truth** (kimi-1). Every guarantee named in protocol text or CLI output
  (`stale-protocol`, `confinement-unproven`, `drifted-core`, "blocks rather than substituting",
  "reports every block it replaces or removes") MUST be asserted by an end-to-end test driving the
  REAL command entry point against a fixture deck, asserting exit code, reported status, and
  resulting file state — not only unit tests of internals. **A guarantee without such a test MUST
  NOT be documented as landed.**
- **G8 — lock byte-verification** (kimi-1). Adoption and continuation MUST byte-compare the on-disk
  release against the deck lock's core hash and refuse on mismatch. A test MUST install a
  same-version-label / different-bytes release and prove refusal. Without this, D8's lock records a
  hash nothing checks — the documented-not-wired failure this project has shipped before.

### D12 — `protocolRole` is retired, replaced by the deck lock

`meta/version.json`'s `protocolRole` today distinguishes `source` (this repo, which AUTHORS the
protocol and is never auto-written — `internal/app/preflight.go:390`) from `consumer` (every other
deck, which may be synced). Under D1/D8 that distinction is carried by whether a deck has a **deck
lock** and whether it is the release-publishing repo, so the field is redundant. It is **retired**:
kept readable for one release for backwards compatibility, ignored by the resolver, and removed
after migration. Raised as undispatched by hermes-1 in round 2 and by kimi-1 in signoff.

## Drafter position changes

Per §15.5 — claude-1 is facilitator, drafter, and a participant. Changes in the drafter's position
since its most recent round file, `round-02/claude-1.md`:

### DPC-1 — Prevention IS achievable for parley-launched agents (reversal, round 1 → round 2)

> Prior position, verbatim (`round-01/claude-1.md`):
>
> **An agent runs with the user's own privileges.** It can `chmod`, it can write the file with any
> tool. File mode `0444` on `~/.parley/COOPERATION.md` stops an *accident*; it does not stop an
> agent that decides to.

**Prior:** only detection is honest; prevention is unachievable. **New:** prevention is real for
parley-launched participants under a seatbelt profile. **What changed it:** I executed the test
instead of arguing (VC-1). codex-1 was right and the three-way majority, including me, was wrong.

### DPC-2 — The effective snapshot is always materialized (reversal, round 2 → this consensus)

> Prior position, verbatim (`round-02/claude-1.md`):
>
> Store the **hash always, the body only when the render is not reproducible from (core version,
> overlay hash, resolver version)**. … copying ~80 KB into every idea directory buys nothing.

**Prior:** conditional materialization. **New:** always materialize, content-addressed and
deduplicated (D7). **What changed it:** codex-1 and kimi-1 blocked independently and showed the
argument was wrong on its own terms — the inputs can stop being present (pruned release, new
laptop, fresh clone), and then the pinned bytes are unrecoverable, breaking user constraint 1. I
optimized storage against the guarantee the pin exists to provide.

### DPC-3 — Attribution narrowed (this consensus)

Revision 1 of this document said "detection and attribution" for unmanaged agents. codex-1 blocked
it as an overclaim. **New:** attribution only for changes made through the attended publisher; an
unexplained hash mismatch is `DETECTED-UNATTRIBUTED` (D9).

### DPC-4 — Overlay version range and extension dependencies (this consensus)

My round-2 file raised "does the overlay need a version of its own?" as an open question and did not
answer it. It is now decided: the overlay declares the core version range it targets, and an
extension declares the core blocks it depends on, defaulting to all sealed blocks (D10).

**Nothing else qualifies.** Corrections of fact in this document (the ranking's attribution, the
`mergePreservingZones` role, G2's scope) are drafting errors fixed on reviewer conditions, not
changes in the drafter's design position.

## Signoffs

Revision 1: `signoffs/<agent-id>.md` (hermes-1 ACCEPT; codex-1, kimi-1 BLOCK).
Revision 2: `signoffs/rev2/<agent-id>.md`.
