---
idea: meta-protocol-change-global-core-protocol
status: fix-up-cycle-6
implementer: claude-1
started: 2026-08-07
completed: 2026-08-07
branch: parley-deck-cli#main
head-commit: 8ed3c4b (cycle 5); cycle 6 lands in the next commit
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

20 tests, including one that drives **production dispatch** (`Run`) rather than `runProtocol`
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

## Fix-up cycle 2

status: complete
completed: 2026-08-07

### Fixes applied

- **[MAJOR, hermes-1 + kimi-1] The §7 protocol text I wrote shipped guarantees the code does not
  implement** — per-idea pinning stated as present fact, and `DETECTED-UNATTRIBUTED` cited as a
  shipped signal. Both are ranks 2 and 4. All three copies now carry an explicit **"Not yet in
  force — do not rely on it"** paragraph naming what IS in force. This is the exact violation of
  G7b, a gate I wrote, inside the protocol section I wrote.
- **[CRITICAL, codex-1] `droppedContent` lost section context and multiplicity.** A flat set of
  trimmed lines made a dropped line look kept when an identical line existed elsewhere in the
  render, and collapsed N dropped copies into one. Now per-section and multiplicity-aware
  (consume-a-copy), so semantic erasure cannot hide behind a coincidental match.
- **[MAJOR, codex-1 + kimi-1] Symlink hardening covered only the final create.** An intermediate
  link at `~/.parley/protocol` redirects the whole store. `assertNoSymlinkComponents` now checks
  the store's own components — deliberately NOT up to the filesystem root, because `/var` is itself
  a symlink on macOS and the first attempt rejected every temp directory on the machine, failing
  eight of its own tests.
- **[MINOR, codex-1] `ValidVersion` validated a trimmed copy while the untrimmed string was used**
  in the path. It now rejects untrimmed input outright.
- **[MINOR, kimi-1] `TestProtocolRejectsPathTraversalInTheLock` was VACUOUS** — it asserted only a
  non-zero exit, which render returns anyway because no release exists at the traversal target. It
  passed with the fix reverted (verified). It now asserts the *reason*, and a new
  `TestLoadRefusesToEscapeTheStore` plants a real readable body outside the store so the refusal
  cannot pass for the wrong reason.
- **[MAJOR, hermes-1] The TTY refusal and the `status` read-error path had no end-to-end tests.**
  Both added, driving the real entry point; the refusal test also asserts the message does not
  overclaim its own strength.
- **[MINOR, kimi-1] CRLF *cores* were never normalized** (mixed endings, non-convergent renders,
  `\r\r\n` mangling) and **the synced stamp was reported as lost project content on every version
  bump**. Both fixed; the report preamble no longer contradicts its entries.
- A failed publish now removes the directory it created so the version label is not bricked. This
  is defensive code, **not claimed as a tested guarantee**: the failure could not be injected on
  this filesystem (a `0555` store root remained writable), and per G7b an untested guarantee is not
  documented as landed.

### Verification

`go build ./...`, `GOOS=windows go build ./...`, `GOOS=linux go build ./...`, `go vet ./...` and
`go test ./...` are all clean. 20 protocol tests (count generated, not asserted from memory — reviewers found the previous
figure did not reproduce), one of which drives production dispatch.

## Fix-up cycle 3

status: complete
completed: 2026-08-07

### Fixes applied

- **[MAJOR, hermes-1 + kimi-1] Cycle 2's "per-section" fix was not per-section.** The counts were
  built over the WHOLE rendered body, so a line deleted from §3 still matched an identical line
  surviving in §11 and reported nothing. Both reviewers caught that the documented fix was not the
  implemented one — the third time this cycle that a claim outran the code, which is why every
  cycle now ends in a re-review rather than a self-assessment. `indexBySection` now indexes the
  render per section and matches within the same section; multiplicity is preserved by consuming a
  copy. A new test constructs exactly the cross-section case.
- **[MAJOR, hermes-1 + MINOR, kimi-1] `assertNoSymlinkComponents` was applied on write only.**
  `Load` now refuses a symlinked store component and a symlinked release directory: a read that
  silently comes from elsewhere is how a deck ends up governed by a protocol nobody published.
- **[MAJOR, hermes-1 + MINOR, kimi-1] The changelog still cited `DETECTED-UNATTRIBUTED` as a
  shipped guarantee** after the protocol copies were corrected. Fixed; it now names exactly what
  ships and marks pinning and the tamper signal as ratified-but-not-implemented.
- **[MINOR, kimi-1] The new §7 "in force" sentence itself overclaimed.** It now says the publisher
  refusal stops an ordinary agent run but not one that allocates a pty, and that releases are
  refused through a symlinked store.
- **[MINOR, kimi-1 round-02, still open] `protocol render` swallowed deck-file read errors**,
  treating an unreadable deck as empty — which would render as though the deck had no content and
  report nothing lost. Now reports and exits non-zero.
- **[MINOR, hermes-1] `protocol check`'s preamble contradicted its content-based entries** ("
  sections present here" vs per-section line counts). Reworded.
- **[NIT, kimi-1] The test count did not reproduce.** It is now generated from an actual run rather
  than written from memory.

### Verification

`go build ./...`, `GOOS=windows`, `GOOS=linux`, `go vet ./...`, `go test ./...` all clean.
20 protocol tests.

## Fix-up cycle 4

status: complete
completed: 2026-08-07

### Fixes applied

- **[MAJOR, hermes-1 + kimi-1] The cross-section test was a tautology.** Its fixture put the line
  in the same section on both sides, so it passed with the per-section fix reverted — verified.
  The headline fix of cycle 3 therefore had zero regression protection. The fixture now places the
  line under §2 in the deck while the render has it only under §3, and reverting to global matching
  makes it FAIL (verified both directions).
- **[MINOR, kimi-1] Three cycle-3 fixes had no test that noticed their absence** — Load's
  release-directory symlink refusal, Publish's store-component refusal, and render's read-error
  report. Each now has a test, and each was proven to FAIL with its fix reverted. The third needed
  a second pass: it initially passed for the wrong reason (the write failed anyway), so it now uses
  `--dry-run` and asserts the reported reason.
- **[NIT, kimi-1] A core-side section RENAME reports a surviving line as not carried forward.**
  Kept as-is deliberately: this is the cost of per-section strictness, and a data-loss report
  should err loud. Recorded rather than silently accepted.

### Reviewer availability

codex-1 declined rounds 3 and 4: the prompt asked it to "try to defeat" the symlink and version
guards, which its policy filter read as offensive security work. That is a reasonable refusal to a
badly-worded request — the task is verifying our own defences in our own repository. Round 5 is
reworded accordingly. codex-1's round-01 and round-02 findings are all addressed; its absence from
3-4 is recorded, not treated as agreement.

### Verification

`go build ./...`, `GOOS=windows`, `GOOS=linux`, `go vet ./...`, `go test ./...` all clean.
23 protocol tests.

## Fix-up cycle 5

status: complete
completed: 2026-08-07

Round 5: hermes-1 CLEAN, kimi-1 CLEAN, codex-1 FINDINGS (reviewing cycles 2-4 together, having
declined 3-4 over wording). All of codex-1's findings upheld.

### Fixes applied

- **[MAJOR] `droppedContent` lost content three more ways**, each proven by codex-1 with a scratch
  probe. This function has now been wrong in four distinct ways across four cycles, which is itself
  the finding: comparing documents for lossless transformation is harder than it looks, and every
  simplification I reached for turned out to be a lossy one.
  - `TrimSpace` on both sides made an INDENTED code line and its unindented prose twin compare
    equal, so a Markdown-meaning-changing edit passed as carried forward. Comparison now trims only
    trailing whitespace.
  - The synced-stamp skip was location-blind, so genuine project prose beginning with that prefix
    vanished unreported. It is now structural: skip only where the render has a regenerated stamp
    in the same section.
  - Only `##` and `###` were headings, so a `####` subsection was not a section boundary and its
    losses were masked by its parent. All ATX levels now count.
  All three of codex-1's probes are permanent tests.
- **[MAJOR] The third protocol copy was still uncommitted.** Committed as
  `parley-deck-skill@4b80468`; all three copies verified to carry the same text.
- **[MINOR] "Only the user may change the global core" contradicted the documented pty behaviour.**
  Reworded across all three copies: it is a rule backed by mechanism, not a proof, and what the
  tooling guarantees is that a change cannot happen quietly or through the normal path.
- **[MINOR] The Phase-8 record named neither the phase nor the reviewed HEAD.** Frontmatter now
  carries `status: fix-up-cycle-5` and the reviewed commit.

### Verification

`go build ./...`, `GOOS=windows`, `GOOS=linux`, `go vet ./...`, `go test ./...` all clean.
26 protocol tests.

## Fix-up cycle 6

status: complete
completed: 2026-08-07

Round 6: codex-1 FINDINGS (one MAJOR). hermes-1 and kimi-1 were interrupted before writing and
re-run in round 7.

### The multiset was abandoned, not patched again

codex-1 showed four more transformations where the report stayed silent while the document's
meaning changed: a lost Markdown hard break (two trailing spaces), prose moved out of an HTML
comment, a rule moved between two sections that share a heading name, and a second
stamp-prefixed line swallowed by the stamp exemption.

The pattern across six cycles is the finding.  was a MULTISET comparison, and it was
defeated five separate times — by content under a shared heading, by a line surviving in another
section, by a per-section claim that was not implemented, by  equating an indented code
block with unindented prose, and now by ORDER. **A multiset cannot see order, so no further patch
would have made it correct.** Each fix was locally reasonable and the approach was wrong.

It is now an **LCS sequence diff over exact lines** (only a CRLF's  is stripped). A moved line
is a removal plus an addition; a lost hard break is a changed line; whitespace is compared as
written. And the claim is narrowed to what it can support: it reports what regeneration will
change so a human can judge — it does not interpret Markdown semantics.

Exactly one stamp is forgiven, the regenerated one; a second stamp-prefixed line is project prose
and is reported. All four of codex-1's cases are permanent tests.

### Verification

`go build ./...`, `GOOS=windows`, `GOOS=linux`, `go vet ./...`, `go test ./...` all clean.
29 protocol tests.
