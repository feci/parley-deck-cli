---
idea: integrate-parley-bidding-addon
review-round: 20
agent: kimi-1
date: 2026-07-31
reviewed-commit: 2b7ca3e
---

## Verdict
BLOCK

## Outstanding findings — closed or not

- **Windows root walk (round 19: codex-1 MAJOR, kimi-1 MAJOR) — CLOSED.** `splitAtRoot`
  (lib/installer.js:1343) anchors on `path.parse(resolved).root`. Measured from this POSIX host
  with `path.win32` injected: drive-absolute, UNC (mixed separators), mixed `/`+`\` on drives,
  repeated separators, trailing separator, `\\?\C:\` device paths, and drive-relative `C:` /
  `C:foo` all yield the correct root and parts. The regression exercises the real exported
  function with the win32 impl injected, not `path.win32`'s own behaviour — the self-caught
  weakness from cycle 23 is genuinely fixed. POSIX output is unchanged.
- **Raw link target normalized before its dependencies were seen (round 19, codex-1 MAJOR) —
  CLOSED for the arm as measured, but the family survives — see New findings.**
  `walkRawTarget` (lib/installer.js:1393) records each entry before applying `..`; the
  regression fails at `2b680a2` and passes at `2b7ca3e` (both measured here), and the same arm
  blocks `uninstall --target all` identically (measured: "Destinations depend on each other",
  nothing removed).
- **`canonicalSegment` NFC scope (round 19, kimi-1 NIT) — CLOSED.** It now folds only where
  `CASE_INSENSITIVE_FS` holds (lib/installer.js:1330).
- **Missing regression for the arm that decided the model (round 19, kimi-1 MINOR) — CLOSED.**
  The firmlink pin is present and, as labelled, is a pin rather than evidence: it passes at
  `2b680a2` too (measured: 107 pass / 2 fail there, the pin among the passing).
- **Stale comment naming a deleted caller (round 19 NITs) — CLOSED.** The comment at
  lib/installer.js:2316 now says "the removal transaction".
- One honesty note on the cycle-23 regression table: "chain walk starts at the platform root"
  does fail at `2b680a2`, but as `TypeError: installer.splitAtRoot is not a function` — the old
  code has no export to exercise. The claim "fails" is accurate; the mechanism is absence, not
  misbehaviour. No action needed.

## New findings

### MAJOR — the resolution walk is lexical where the kernel is physical; two measured arms pass the gate green and break the fleet

Cycle 23 fixed the pure-lexical member of this family (`transient/../../../../away`: no links
in the middle, no linked ancestors). The family is wider. `walkRawTarget` and the outer
`resolutionTouchpoints` loop model kernel resolution lexically from the link's **spelling**:
intermediate components of a raw target are recorded by location but never followed, and `..`
is applied with `path.dirname` on the spelling rather than against the link's real location.
`identityChain`'s physical `statSync` anchors mask the gap whenever a link lands exactly on the
other destination's **root** — I measured those shapes and they are all correctly refused
(shared: absolute-target aliasing, physical `..` landing on the other dest, even a POSIX
backslash single-component name) — but not when resolution passes through the tree's
**interior** or through a physically different parent. Both arms reproduce identically at
`2b680a2` and `2b7ca3e`, so this is not a cycle-23 regression; it is a hole in the model I
called complete in round 19, found by the full-scope round that exists to find it.

**Arm 1 — an intermediate component of the raw target is itself a symlink into the other
destination's interior.** Seed kimi (dest `KM/skills/parley-deck` = P, with `P/sub/deep`
created); codex container `B/skills` → raw target `midlink/deep`, where `B/midlink` → `P/sub`.
`walkRawTarget` records `B/midlink` and `B/midlink/deep` by location; the anchor jump in
`identityChain` lands on `dev:ino(P/sub)`, so `dev:ino(P)` never appears in any chain and no
check fires. Measured: `install --target all` → **ok: true**, codex "installed" — its payload
committed *inside kimi's tree* at `P/sub/deep/parley-deck` — and kimi "replaced" later in the
same plan order, so kimi's own commit renamed P aside and discarded the codex payload the same
transaction had just written. `doctor codex` immediately after the green install: **missing**.
The fleet install reported success while leaving the fleet broken. Under uninstall (codex unit
present inside kimi's tree, `uninstall --target all`): **ok: true**, both "removed" — the
quarantine rename moved a directory out of the other unit's tree mid-transaction, the exact
configuration the gate exists to refuse.

**Arm 2 — the container link sits below a symlinked ancestor and its raw target starts with
`..`.** `linkA` → `real/A`; container `linkA/container/skills` →
`../../Btree/skills/parley-deck/inner/deep`, where `inner` is a link *inside* kimi's tree. The
kernel resolves `..` physically from the link's real location (`real/A/container` → `real`);
`walkRawTarget` climbs the spelling (`linkA/container` → `linkA` → home). Measured:
`install --target all` → **ok: true**, codex's payload written *through* kimi's tree to an
outside directory; `uninstall --target kimi` → **ok: true**, and codex's dest immediately
dangles — `doctor codex`: missing, payload orphaned outside. The gate's own message names the
criterion — "resolving one passes through the other" — and this is precisely that, unrefused,
in both the install and the uninstall direction.

Both arms are assembled only from configurations the gate already accepts and supports:
symlinked skills containers (the agy/gemini arm that created the gate), `..` inside link
targets (the cycle-23 arm itself), and symlinked runtime homes. Nothing adversarial is
required, and the consequence is the worst this transaction has: a green fleet mutation that
silently breaks another unit. One root cause covers both arms (and a third, observed but not
measured into harm: the raw-target split on `/[\\/]+/` treats a literal backslash — a legal
POSIX filename character — as a separator, where the kernel sees one component).

**What must change:** the walk must model the kernel physically — follow intermediate symlink
components inside a raw link target, and resolve `..` from the link's real location (e.g.
anchor relative walks at the physical parent of the hop, not the lexical `path.dirname` of its
spelling) — so that pass-through is refused for interior and linked-ancestor traversal, not
only for root-exact landings. Both regressions above belong in the suite; each fails at
`2b7ca3e` today.

### Answer to the assigned parity question

Yes — uninstall's quarantine path is gated identically to install's commit path:
`removeFleetAtomically` consults the same `aliasedDestinations(plan)` before any rename
(lib/installer.js:1643 vs :1528), and I measured parity in both directions — the raw-target arm
and absolute aliasing block uninstall exactly as they block install, and the MAJOR's arms are
missed identically. Parity holds; the blindness is shared, not asymmetric.

### Assigned walkRawTarget probes, for the record

Absolute targets (refused correctly via physical identity), repeated separators (collapsed),
`..` past the root (clamps via `dirname`, matching the kernel), a 2-link loop (terminates,
~5 ms, `seen`-set dedup), and a 70-hop chain (terminates, ~12 ms; macOS ELOOP at 32 binds long
before the 64-hop bound, and the unresolvable container is then refused by the broken-parent
preflight). All behaved.

## Release judgement

Not releasable as 2.1.0. The one thing that must change: the resolution walk must model the
kernel physically — follow intermediate symlinks in raw link targets and resolve `..` from the
link's real location — so the gate refuses the two measured pass-through arms above instead of
returning ok: true while the fleet breaks. Everything else cycle 23 shipped (platform roots,
the raw-target fix as scoped, NFC scope, the pin, both new regressions, uninstall parity) I
verified as correct.

## What I verified

- Full suite at `2b7ca3e` in a `git archive` copy (node_modules symlinked, read-only against
  the working tree): **364 node tests, 0 fail**; python leg **54/54 on 3.14**; manifest check
  ok — 47 files, `sha256:7854adf1…95a6d`, matching the cycle-23 record.
- Regressions against a `git archive` copy of `2b680a2` with the cycle-23 test file: raw-target
  regression fails, splitAtRoot regression fails (TypeError — export absent), firmlink pin
  passes. At `2b7ca3e` all three pass.
- `splitAtRoot` with `path.win32` injected: drive, UNC, mixed separators, repeated separators,
  trailing separator, device path, drive-relative, relative input, POSIX control — all correct.
- New probes (real install/uninstall calls in mkdtemp homes, cleaned up afterwards):
  - Arm 1 (midlink into interior): install ok:true; codex payload committed inside kimi's tree
    and destroyed by kimi's replacement in the same transaction; doctor codex → missing.
    Uninstall with the unit inside the other tree: ok:true, cross-tree rename.
  - Arm 2 (symlinked ancestor + leading `..`): install ok:true, payload written through kimi's
    tree; `uninstall --target kimi` ok:true, codex dest dangling, payload orphaned.
  - Both arms reproduce identically at `2b680a2` — pre-existing model gap, not a c23 regression.
  - Root-exact landings all refused correctly (absolute aliasing install + uninstall, physical
    `..` shared dest, POSIX backslash single-component name — caught as "shared"/"overlaps").
  - Uninstall parity: cycle-23 raw-target arm blocks `uninstall --target all` ("Destinations
    depend on each other", nothing removed); absolute aliasing blocks uninstall ("shared").
  - Termination: `..` past root, 2-link loop, 70-hop chain — all return in milliseconds.
- Working tree untouched: all probing ran against archive copies under /tmp; no writes to the
  repo.
