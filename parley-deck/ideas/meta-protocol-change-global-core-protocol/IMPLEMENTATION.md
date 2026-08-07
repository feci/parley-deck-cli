---
idea: meta-protocol-change-global-core-protocol
status: implemented
implementer: claude-1
started: 2026-08-07
completed: 2026-08-07
branch: parley-deck-cli#main
head-commit: (see release commit)
design-pr: n/a
implementation-pr: n/a
---

## Summary of work

Ships **rank 1** of the ratified plan plus the protocol-text change. Ranks 2-4 (per-idea pinning,
the overlay, the detection-layer enforcement) remain ratified and scheduled; nothing here claims a
guarantee it does not implement (G7b).

### New package `internal/protocolcore`

- `core.go` — the immutable release store at `~/.parley/protocol/core/<version>/`. `Publish`
  **refuses to modify an existing release** and writes `0444`; that refusal is D1's write-once
  guarantee and the structural half of G2. `Load` deliberately verifies nothing beyond existence,
  because only the caller knows which hash it expected (G8 belongs there). `ErrNoRelease` is a
  distinct error so callers can implement D8 — a missing release blocks adoption/rendering, never
  continuation.
- `render.go` — `Render(release, priorDeckBody)`, a **pure function**: no filesystem, no clock, the
  synced-stamp derived from the release, so two machines holding the same release render
  byte-identical output. That is what makes G1's idempotence testable at all. It preserves the six
  identity slots and returns `Removed` — the sections the deck has and the core does not.

### New command `parley protocol`

`status | render | check | publish`, wired into dispatch and `--help`.

- `render` reports **preserved** and **removed** blocks in preview AND on apply (G1). Preview is
  the default; `--yes` applies; the write is atomic and preserves the file mode.
- `check` reports a hand-edit and **never overwrites** — the posture `roster show` takes for
  `legacy-roster` — exiting non-zero so a script can act on it.
- `publish` is the attended publisher, G2's sole audited exception.

### Protocol text (D11)

The §7 blast-radius clause added to all three `COOPERATION.md` copies: a CORE change needs the meta
idea and explicit user ratification and is user-only; a DECK change is a normal idea. Changelog
entry added in the §7 template format (G6).

## Deviations from FINAL.md

None in substance. The deck lock is read as a minimal `core-version:` line rather than a full lock
with all five fields — the remaining fields belong to rank 2 (per-idea pinning) and writing them
now would be recording data nothing yet verifies, which is the failure mode G7b forbids.

## Notes for reviewers

**A gate caught a real defect in its own implementation, which is worth reading.** The first
`hasTTY` used the conventional `stat(stdin).Mode()&os.ModeCharDevice` test. Driving the REAL entry
point showed `parley protocol publish … < /dev/null` publishing successfully — `/dev/null` is a
character device, and redirecting stdin that way is exactly what an agent run does. The gate was
documented as landed and was wrong at the call site. It now asks the kernel for terminal attributes
(`unix.IoctlGetTermios`), which succeeds only on a real tty, and the platform fallback was deleted
rather than kept weak: a platform without the ioctl refuses instead of accepting `/dev/null`.

Nine tests drive the real command entry points against fixture decks (G7b), covering idempotence,
identity preservation, removal reporting in both preview and apply, check-reports-never-overwrites,
block-on-missing-pinned-release, refuse-to-guess-a-version, write-once releases, the JSON status
shape, and render onto a fresh deck.
