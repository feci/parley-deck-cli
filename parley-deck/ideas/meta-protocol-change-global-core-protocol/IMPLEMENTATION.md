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
  refuses when the release **directory** already exists (write-once is per release, not per file),
  opens with `O_EXCL|O_NOFOLLOW` so it can never write through a planted symlink, validates the
  version as a safe path element, and writes `0444`. `Load` validates the version too — it comes
  from a committed lock any contributor can edit. `ErrNoRelease` is distinct so callers can
  implement D8: a missing release blocks adoption and rendering. **Continuation is rank 2 and is
  not implemented here**, so nothing in this slice claims it.
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
(`unix.IoctlGetTermios`), which succeeds only on a real tty.

**And then the fix broke the Windows build, which all three reviewers caught.** Deleting the
fallback file did not make the platform stricter — it removed the constant the build needs, so
`GOOS=windows go build ./...` failed. Since winget ships Windows binaries, that would have shipped
a package that does not compile. The fallback is back, and it **fails closed**: a platform that
cannot prove a terminal reports the publisher unavailable rather than accepting `/dev/null`.
`GOOS=windows` and `GOOS=linux` builds are now part of the check.

**The TTY gate's real strength, stated plainly** (codex-1, kimi-1): it stops an ordinary agent run,
whose stdin is a pipe or `/dev/null`. It does **not** stop an agent that allocates a pty — a stock
one-liner does. The command says so in its own refusal text rather than implying more. The durable
guarantees in this slice are write-once releases, `O_EXCL|O_NOFOLLOW`, no agent-accessible write
path, and hash detection; confinement is the ratified DF-1 follow-up and is not shipped.

**The removal report was heading-based and is now content-based.** Reviewers showed that deck
content under a heading the core ALSO has vanished with nothing reported — the exact silent-erasure
class G1 exists to prevent, and the one that destroyed a real deck's local section during the
2026-08-06 fleet sync. `droppedContent` now compares content, grouped under the heading it sat
beneath. CRLF decks no longer produce a false removal for every section, and no longer get mixed
line endings.

Fifteen tests, including one that drives **production dispatch** (`Run`) rather than `runProtocol`
— codex-1 was right that a test one layer below dispatch can pass while the command is unreachable.
They cover idempotence, identity preservation, removal reporting in preview and apply, content lost
under a shared heading, CRLF decks, check-reports-never-overwrites, block-on-missing-release,
refuse-to-guess-a-version, lock path traversal, write-once per release, symlink write-through,
unsafe versions, the JSON status shape, and render onto a fresh deck.

## Fix-up cycle 1

status: complete
completed: 2026-08-07

### Fixes applied

- **[CRITICAL, all three] Windows build broken.** `termios_other.go` had been deleted, leaving no
  `termiosGet` for non-BSD/non-Linux. Restored and made fail-closed. `GOOS=windows` / `GOOS=linux`
  builds verified.
- **[CRITICAL, codex-1 + hermes-1 + kimi-1] `Publish` symlink write-through and existing release
  namespace.** Now `Lstat`s the release DIRECTORY and opens with `O_EXCL|O_NOFOLLOW`.
- **[CRITICAL, codex-1 + kimi-1] TTY gate overclaim.** A pty defeats it; the refusal text and
  IMPLEMENTATION.md now say so instead of implying confinement.
- **[CRITICAL/MAJOR, all three] Heading-based removal report.** Replaced with content-based
  `droppedContent`; CRLF handled.
- **[MAJOR, kimi-1] Path traversal via the committed deck lock.** `ValidVersion` on both `Load`
  and `Publish`.
- **[MAJOR, codex-1] Tests bypassed production dispatch.** Added a `Run`-level test.
- **[MAJOR, codex-1] IMPLEMENTATION.md claimed rank-2 continuation.** Corrected.
- **[MAJOR/MINOR, all three] The third protocol copy** was edited but uncommitted, in a sibling
  repo. Committed as `parley-deck-skill@455aafe`.
- **[MINOR, hermes-1] `--help` omitted `publish`**; added.
- **[MINOR, hermes-1] `protocol status` swallowed read errors**, rendering an unreadable store as
  an empty one. Now reports and exits non-zero.
