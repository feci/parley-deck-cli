---
idea: protocol-overlay-local-extension
status: in-progress
implementer: claude-1
started: 2026-08-09
branch: (local-dir transport, working tree)
---

# IMPLEMENTATION — protocol overlay v1

Living execution doc. `FINAL.md` is the static spec; this records what was built, what deviated, and
what is **not** built. Read the "Not implemented" section before reviewing: several binding
decisions are outstanding and must not be reviewed as though they landed.

## Progress

| FINAL item | State |
| --- | --- |
| B2 overlay file + strict grammar | **done** — `internal/protocolcore/overlay.go` |
| B5 `parley.protocol-lock/v2` + hash domain | **done** — `internal/protocolcore/lock.go` |
| B3 terminal boundary placement | **done** — `render.go`, payload appended at end of normalized core body |
| B7 typed change events, source-attributed | **done** — `RenderResult.Events` |
| B7 relocation witness | **done, with one recorded deviation** (below) |
| B1 extend-only | **done** — `replace` is rejected by the grammar with a pointer to the follow-up |
| CLI `protocol overlay show\|validate` | **done** |
| Retire promise in `protocol.go:211` | **done** |
| Lock `core.body-sha256` actually verified | **done** — see "A field the schema demanded and nothing checked" |
| **B6 seventh identity slot (roster annotations)** | **NOT DONE** |
| **G4 / H9 identity zones by declared span** | **NOT DONE** |
| One writer for `COOPERATION.md` | **NOT DONE** |
| Release-format marker | **NOT DONE** |
| Retire promise in `COOPERATION.md:767` | **NOT DONE** |

Build is green on darwin, **linux and windows** (`GOOS=windows go build ./...`) — the Windows break
was a CRITICAL finding last cycle and is explicitly re-checked. `go vet ./...` clean.

## What was built

**The overlay** (`internal/protocolcore/overlay.go`). One strict YAML document at
`parley-deck/protocol-overlay.md`; payloads are literal block scalars. Refuses: unknown keys
(`KnownFields(true)`), a second document, YAML anchors and aliases, any content after the closing
fence, `kind: replace`, an id that is not `deck.<slug>`, an empty rationale or payload, a
`core-sha256` that is not 64 lowercase hex, an empty file, and zero operations.

**The lock** (`internal/protocolcore/lock.go`). Schema v2, nested `core`. A pre-v2 flat lock is
named explicitly and refused with the migration text rather than failing with a generic schema
error, because the fix is a migration and the operator should be told so.

`ReconcileOverlay` blocks before composition on every mismatch: declared-but-absent,
present-but-undeclared, unreadable, empty, and hash mismatch. Unreadable is deliberately never
collapsed into absent — a permissions problem must not silently erase a deck's local content from
the output.

**Composition** (`render.go`). `Render` gained an `*Overlay` parameter and stays pure; the caller
owns the filesystem. The payload is appended at the end of the normalized core body, which is the
terminal boundary *by definition* rather than by lookup — that is why v1 places it with no registry.

## Deviations and decisions

**1. Blank-line edges are excluded before the relocation witness compares.** FINAL B7 says a removed
run is reclassified `relocated` only when it is "byte-identical to exactly one complete decoded
payload". Implemented literally, **the witness never fires**: a block sitting mid-document is bounded
by blank lines belonging to the surrounding sections, so the contiguous non-kept run is
payload-plus-padding and never byte-equal. My first test of it failed for exactly this reason, which
is how it was found.

`trimBlankEdges` narrows the run past leading and trailing blank lines before matching. What is
**not** relaxed: the payload bytes are still compared byte for byte, the match must still be to
exactly one complete payload, and it must still occur exactly once on each side. Interior blank lines
are content and must match.

This is a change to a decision reviewers signed. It is flagged here rather than absorbed, and the
alternative — keeping the literal rule and accepting that the feature silently degrades to the noisy
report consensus rejected — should be considered on review.

**2. A field the schema demanded and nothing checked.** The v2 lock carries `core.body-sha256`, and
nothing compared it to the loaded release. Left that way it would have been decoration: a release
republished under the same version, or a store swapped underneath the deck, would render silently.
`resolveDeck` now blocks on mismatch. This is the "documented but not wired" shape that blocked the
previous cycle, caught here before review rather than during it.

**3. The pre-v2 flat lock is refused outright, not accepted as a legacy form.** No production lock
exists anywhere (the core store is empty and `find` shows no `protocol-lock.yaml` in the fleet), so
there is nothing to be backward compatible with, and accepting both formats would mean a deck could
carry overlay state the older path ignores. All existing tests were migrated to v2.

## Not implemented — do not review these as though they landed

**B6, the seventh identity slot for roster annotations, is NOT built.** This is the outcome of the
D4 dispute and a binding decision. Nothing in the code renders an annotation slot, and
`agents.toml` has no field for it. A deck's roster annotations are therefore still not carried, which
means **the specific content destroyed on 2026-08-06 in `auftra`, `ldx-wt-mail-fixups` and
`librade-algoTrader` is still not protected by this implementation.**

**G4 / the H9 fix is NOT built.** `isTableHeader` still matches prose
(`internal/protocolcore/render.go`), so a core that renames a roster column still empties every
deck's roster. FINAL lists this as an in-slice prerequisite. It is also entangled with an
inconsistency in FINAL itself, which should be resolved on review rather than by me alone:

> B4 says v1 releases hold "the core plus the release-format marker", while the slice requires
> identity zones to be located "by declared span". **There is nowhere for the span to be declared.**
> Fixing H9 needs the release to carry a small fixed slot table — five named spans plus the
> annotation anchor. That is not the per-block registry that was dropped (no permanent block IDs, no
> extents, no tombstones, no per-block policy), but it *is* an addition to the release layout that
> FINAL does not authorize.

**One writer for `COOPERATION.md` is NOT built.** `protocol render`, `roster render` and
`preflight.syncConsumerProtocol` still write it independently, each with its own location logic.
A render can still be undone by the next command.

**The release-format marker is NOT built**, so the fail-closed protection for a future layout change
does not exist yet. The v2 lock provides the stale-binary protection independently and *is* built.

**The promise in the protocol text is NOT retired.** `parley-deck/COOPERATION.md:767` still says the
overlay ships "once that ships". FINAL requires both promises to retire in the same commit; only the
CLI one (`protocol.go:211`) has. Until the other is removed, the CLI and the protocol text disagree.

## Verification

23 tests added or migrated. New coverage in `internal/app/overlay_test.go`: terminal-boundary
placement, idempotence with an overlay, each blocking state as a table, unadopted overlay, pre-v2
flat lock refusal, lock body-hash mismatch, the relocation witness **and its near-miss**, seven
grammar refusals, the whole-file-vs-payload hash distinction, and `protocol overlay` succeeding on a
deck with no overlay.

**Reversion checks (G6).** Five guards were each reverted and the corresponding test confirmed to go
red, with the harness asserting that the revert **compiled** and that the substitution **actually
applied** — the two ways this check failed in the previous cycle. Script:
`scratchpad/revert_check.py`.

```
OK  lock body-sha256 verification        → TestLockBodyHashMismatchBlocks went red
OK  unadopted-overlay block              → TestUnadoptedOverlayBlocks went red
OK  lock hash covers the whole file      → TestLockHashCoversMetadataNotJustThePayload went red
OK  relocation blank-edge trimming       → TestRelocationWitnessIsExactAndNearMissesStayRemoved went red
OK  overlay payload is appended at all   → TestOverlayRendersAtTheTerminalBoundaryAndIsIdempotent went red
```

The suite is green again after every restore, which is itself asserted — a restore that silently
failed would otherwise leave the next check meaningless.

**Not verified.** No end-to-end run against a real deck has been done, because that needs a published
core release and publishing is attended-only. The rehearsal in an isolated `PARLEY_HOME` used during
design is not a substitute for this implementation.
