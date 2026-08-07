---
agent: hermes-1
idea: meta-protocol-change-global-core-protocol
review-round: 1
date: 2026-08-07
reviewed-commit: 4396529
---

## Summary

The implementation ships rank 1 (core store, `protocol render|check|status`, the §7 blast-radius
clause) honestly in its core design: the renderer is genuinely pure, the write-once store is
structurally sound for the normal path, and the TTY-gate fix (ioctl instead of ModeCharDevice)
correctly closes the `/dev/null` hole the implementer caught. `go build ./...` and `go test ./...`
both pass on darwin/amd64.

However, adversarial probing found four real defects:

1. A symlink placed at a release path causes `Publish` to write core content THROUGH it to an
   arbitrary location — the write-once guard uses `os.Stat` (follows symlinks) and never checks the
   directory itself.
2. The commit breaks the Windows build: `termiosGet` and `unix.IoctlGetTermios` exist only behind
   `linux` and `darwin || freebsd || netbsd || openbsd` build constraints. There is no
   Windows/other-platform file. The codebase built on `GOOS=windows` before this commit and does not
   after.
3. `protocol render` silently drops a deck table column that the core's table does not have — the
   deck's wider header is replaced by the core's narrower one while the deck's wider rows survive
   underneath it, misaligned and unreported.
4. The changelog and §7 clause claim shipped guarantees (`DETECTED-UNATTRIBUTED`, the TTY-gate
   refusal) that have no end-to-end test — a direct G7b violation on the exact failure mode this
   project has shipped before.

## Findings

### [CRITICAL] Windows build is broken — no termios fallback for non-BSD/non-Linux platforms

`internal/app/protocol.go:12` imports `golang.org/x/sys/unix` and calls `unix.IoctlGetTermios` at
line 312. The constant `termiosGet` is defined only in:

- `internal/app/termios_linux.go` — `//go:build linux`
- `internal/app/termios_unix.go` — `//go:build darwin || freebsd || netbsd || openbsd`

There is no file for Windows or any other platform. Cross-compilation confirms the break:

```
$ GOOS=windows GOARCH=amd64 go build ./...
# parley-deck-cli/internal/app
internal/app/protocol.go:312:21: undefined: unix.IoctlGetTermios
internal/app/protocol.go:312:50: undefined: termiosGet
EXIT=1
```

Before this commit (at fac2421), the same command exits 0:

```
$ git checkout fac2421 && GOOS=windows GOARCH=amd64 go build ./...
EXIT=0
```

The IMPLEMENTATION.md says "a platform without the ioctl refuses instead of accepting `/dev/null`"
(line 62). That is not what happens — on Windows the package does not compile at all, so the
publisher does not "refuse", it fails to build. The project uses `bubbletea` which explicitly
supports Windows (`github.com/erikgeiser/coninput` is in go.mod), so Windows was a viable target
that this commit broke.

Suggested fix: add `internal/app/termios_other.go` with `//go:build !linux && !darwin && !freebsd
&& !netbsd && !openbsd` that provides `const termiosGet = 0` (or a sentinel) and a `hasTTY`
implementation that returns `false` — the fail-closed behavior the IMPLEMENTATION.md claims.

### [CRITICAL] Publish writes through a symlink to an arbitrary location

`internal/protocolcore/core.go:102` checks write-once with `os.Stat(path)` where `path` is
`dir/COOPERATION.md`. `os.Stat` follows symlinks. If a symlink is placed at the release directory
`~/.parley/protocol/core/<version>/` pointing to an arbitrary location, the `os.Stat` on
`<symlink>/COOPERATION.md` returns `os.ErrNotExist` (the target dir has no such file), the
write-once guard passes, `os.MkdirAll` is a no-op on the existing symlink, and `os.WriteFile`
follows the symlink — writing core content to the attacker-chosen location.

Confirmed by test:

```
=== RUN   TestAdversarialPublishSymlink
    protocol_adversarial_test.go:154: Publish into symlink: err=<nil>, rel={Version:1.0.0 ...}
    protocol_adversarial_test.go:159: CRITICAL: Publish wrote THROUGH a symlink to an arbitrary location
--- FAIL: TestAdversarialPublishSymlink
```

The `core-version` string is validated against traversal (`..`, `/`, `\`) at core.go:97, but a
symlink pre-placed at a legitimate version path bypasses that entirely. On a shared machine, any
process that can create a symlink under `~/.parley/protocol/core/` can redirect a future publish.

Suggested fix: use `os.Lstat` instead of `os.Stat` to detect symlinks at both the release directory
and the file path. Refuse if either is a symlink. Alternatively, use `os.MkdirAll` with
`os.FileMode` and check `os.Lstat(dir)` after creation to verify it is a real directory, not a
symlink.

### [MAJOR] Publish does not detect a pre-existing release directory without COOPERATION.md

`internal/protocolcore/core.go:102` checks only whether the COOPERATION.md *file* exists, not
whether the release *directory* exists. If the directory `~/.parley/protocol/core/1.0.0/` is
pre-created (by any process) without a COOPERATION.md file, `os.Stat(path)` returns
`os.ErrNotExist`, the write-once guard passes, and `Publish` writes content into it:

```
=== RUN   TestAdversarialReleaseDirNoFile
    protocol_adversarial_test.go:175: Publish when dir exists but no file: err=<nil>, rel={Version:1.0.0 ...}
    protocol_adversarial_test.go:179: Publish succeeded despite pre-existing release dir
--- FAIL: TestAdversarialReleaseDirNoFile
```

D1 says releases are "write-once" and the error message says "release 1.0.0 already exists and
releases are write-once" — but the check does not actually detect this case. A release directory
that exists but is incomplete (no COOPERATION.md) should be treated as an existing release, not
written into.

Suggested fix: check `os.Stat(dir)` (or `os.Lstat(dir)`) before `os.MkdirAll`. If the directory
already exists, refuse with the same write-once error.

### [MAJOR] protocol render silently drops a deck table column the core does not have

When a deck's §2 roster table has more columns than the core's (e.g., an `Adapter` column), the
renderer replaces the deck's header with the core's narrower header but keeps the deck's wider data
rows. The extra column data dangles under a misaligned header, and the column loss is not reported.

`internal/protocolcore/render.go:64-82`: `isTableHeader` matches the core's `| Agent ID` line,
`tableBodyFor` returns the deck's rows, and the code appends the core's header line (line 66:
`out = append(out, line)`) followed by the deck's rows (line 77: `out = append(out, body...)`).
The deck's own header line is never preserved — only its data rows.

Confirmed by test — a deck with a 4-column table rendered against a 3-column core:

```
line 9:  | Agent ID       | Workspace dir                       | Role          |    ← core's 3-col header
line 10: | -------------- | ----------------------------------- | ------------- |    ← core's 3-col separator
line 11: | `claude-1`     | `./`                                | facilitator   | claude        |  ← deck's 4-col row
```

The `Adapter` column header is gone, the `claude` value dangles as a 4th column under a 3-column
header. `removedSections` (render.go:184) only compares `##` headings, not table structure, so this
loss is invisible to the G1 removal report:

```
=== RUN   TestAdversarialTableHeaderNotReported
    CONFIRMED: the Adapter column loss is NOT reported in preview — silent data loss
```

G1 requires "report every block replaced or removed." A table with a different column set is a
replaced block, and it is not reported. This is the same class of silent erasure the project
already shipped once (the 2026-08-06 fleet sync that destroyed a deck's local section).

Suggested fix: when the deck's table header differs from the core's, either preserve the deck's
header (not the core's) or report the column difference as a replacement. At minimum, compare the
header lines and report a mismatch.

### [MAJOR] G7b violation: claimed guarantees with no end-to-end test

G7b: "Every guarantee named in protocol text or CLI output MUST be asserted by an end-to-end test
driving the REAL command entry point ... A guarantee without such a test MUST NOT be documented as
landed."

Claims without tests:

1. **"`parley protocol publish` refuses without a controlling terminal"** — stated in the §7
   clause (`parley-deck/COOPERATION.md:761`) and the changelog
   (`parley-deck/meta/protocol-changelog.md:23`). There is NO test that drives `protocol publish`
   through `runProtocol` and asserts the TTY-gate refusal. `hasTTY` and `protocolPublish` have zero
   test coverage:

   ```
   $ grep -rn "hasTTY|protocolPublish|\"publish\"" internal/app/*_test.go
   (no matches)
   ```

2. **"detection with `DETECTED-UNATTRIBUTED` for anything else"** — stated in the changelog
   (line 24) as part of "the shipped guarantee." `DETECTED-UNATTRIBUTED` appears in zero Go files.
   It is not implemented, not tested, yet the changelog presents it as shipped.

3. **"An idea that is already open completes under the protocol version it was pinned to"** —
   stated in the §7 clause (`parley-deck/COOPERATION.md:766`). Per-idea pinning (rank 2) is not
   shipped. The claim is present-tense in the protocol text but has no implementation and no test.
   FINAL.md:73 says "Gates G5, G7 and G8 bind the ranks not shipped this cycle" — but the §7 text
   does not hedge; it states the continuation guarantee as fact.

Suggested fix: either add end-to-end tests for each claim, or hedge the text to say "will" rather
than "does." At minimum, the TTY-gate refusal must have a test (it is the most testable of the
three and is the exact "documented as landed and wrong at the call site" pattern this project has
shipped before).

### [MINOR] Top-level --help omits `publish` from the protocol subcommands

`internal/app/app.go:125`:
```
  %s protocol status|render|check [--dir DIR] [--dry-run] [--yes] [--json]
```

The `publish` subcommand is missing from this summary line. The subcommand's own usage text
(`protocol.go:24`) does list it. A user running `parley --help` will not discover `publish`.

### [MINOR] protocol status silently drops directory-read and lock-read errors

`internal/app/protocol.go:134`: `versions, _ := store.Versions()` — if the store root exists but
is unreadable (permissions), the error is silently dropped and `versions` is `nil`, reported as
"(none)". `protocol.go:135`: `pinned, _ := pinnedVersion(root)` — same: a readable-error on the
lock file produces `pinned = ""`, reported as "—". Both are diagnostic surfaces; silently reporting
"no releases" or "no pin" when the real condition is "permission denied" could mislead a user into
re-publishing or re-pinning.

### [MINOR] FINAL.md says "all three COOPERATION.md copies" but only two were updated

FINAL.md:48-49: "The §7 blast-radius clause added to all three COOPERATION.md copies." The diff
touches two: `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`. The
third copy — the skill fallback at
`~/.hermes/skills/parley-deck/references/COOPERATION.md` — does not contain the blast-radius clause
(verified: `grep -c "Blast radius"` returns 0). This is out of repo scope, but the claim "all
three" is inaccurate for a user who has the skill installed.

### [NIT] check --json on a missing pinned release produces a text error, not JSON

When the pinned release is missing, `resolveRelease` fails before the JSON branch, so `check --json`
writes a plain-text error to stderr and returns 1 with empty stdout. A machine consumer expecting
JSON gets none. This is a minor inconsistency for a surface labeled `--json`.

## Gate-by-gate assessment

**G1 (idempotent render, reports removals in preview and apply):** PARTIALLY MET. Idempotence is
tested and holds. Removals (by heading) are reported in both preview and apply. BUT: table column
differences are silently dropped without reporting (see MAJOR finding above). A deck whose §2 table
header differs from the core's loses columns silently. CRLF line endings survive and are
idempotent. Two tables (roster + host handle) — the host handle section is correctly reported as
removed when the core lacks it.

**G2 (no agent-accessible write path into core):** PARTIALLY MET. The TTY gate uses the ioctl
approach (correct fix for the `/dev/null` hole). `Publish` refuses to modify an existing release
file. BUT: the symlink attack and pre-existing-directory bypass (see CRITICAL and MAJOR findings)
are exploitable write paths. The TTY gate itself has no test — the one guarantee G2 makes that is
most testable is untested. Directory traversal in the version string is correctly blocked (`/`,
`\`, `..`).

**D8 (missing pinned release blocks adoption/rendering, not continuation):** MET for the shipped
slice. `resolveRelease` returns `ErrNoRelease` for a missing pinned version, and `protocolRender`
and `protocolCheck` both fail closed (exit 1, no write). `TestProtocolBlocksWhenPinnedReleaseIsMissing`
proves this. Continuation is NOT implemented (rank 2), and IMPLEMENTATION.md does NOT overclaim it —
correct. The `ErrNoRelease` comment (core.go:30-32) correctly states the D8 semantics. The §7 clause
in COOPERATION.md:766 does state continuation as present-tense fact (see G7b finding above), but
IMPLEMENTATION.md itself is honest: "Ranks 2-4 remain ratified and scheduled; nothing here claims a
guarantee it does not implement."

**G7b (no guarantee documented as landed without end-to-end test):** NOT MET. Three claimed
guarantees lack tests (see MAJOR finding above): the TTY-gate refusal, `DETECTED-UNATTRIBUTED`
detection, and per-idea continuation. The changelog explicitly says "the shipped guarantee is"
followed by these items.

**G6 (changelog in §7 template format):** MET. `parley-deck/meta/protocol-changelog.md` has the
entry with `Idea:`, `Drafted by:`, `Summary:`, and ratification context.

**Renderer purity:** MET. `Render` imports only `fmt` and `strings` — no `os`, `time`, `env`, or
`filepath`. The synced-stamp is derived from `rel.Version` and `rel.SHA256` (render.go:87). It is
genuinely a pure function of `(release, priorDeckBody)`.

## Test-quality assessment

Nine tests, all driving the real `runProtocol` entry point. Assessment of each:

1. **TestProtocolRenderIsIdempotent** — REAL. Would fail if `Render` were not deterministic (e.g.,
   if it used `time.Now()` for the stamp). Checks both file equality and "nothing to do" message.

2. **TestProtocolRenderPreservesIdentityAndReplacesCoreText** — REAL. Would fail if identity
   extraction or slot substitution were broken. Checks specific values (`my-project`, `claude-1`)
   and that stale core text is gone.

3. **TestProtocolRenderReportsRemovedSections** — REAL. Would fail if `removedSections` returned
   nil. Checks both preview and apply paths. Does NOT test the column-drop case (see MAJOR finding).

4. **TestProtocolCheckReportsHandEditAndNeverWrites** — REAL. Would fail if `check` overwrote the
   file (the `after != edited` assertion catches this). Checks exit code and message content.

5. **TestProtocolBlocksWhenPinnedReleaseIsMissing** — REAL. Would fail if `resolveRelease` fell
   back to any installed version. Checks exit code, error names the version, and file unchanged.

6. **TestProtocolRefusesToGuessAVersion** — REAL. Would fail if `resolveRelease` picked the newest
   installed version for an unpinned deck. Checks exit code and error message.

7. **TestCoreReleasesAreWriteOnce** — REAL. Would fail if the `os.Stat` existence check were
   removed. Does NOT test the symlink or pre-existing-directory bypass.

8. **TestProtocolStatusJSON** — PARTIAL. Asserts the JSON shape (pinned, installed, deck_sha256).
   Does not assert `store` or `deck_protocol` fields. Would fail if the JSON structure changed
   for the asserted fields.

9. **TestProtocolRenderOnAFreshDeck** — REAL. Would fail if the core body were not applied.
   Checks no false removal report on an empty prior.

**No tautologies found.** Each test asserts behavior through the real entry point and would fail
if the corresponding fix were reverted. The gap is in coverage, not in test quality: the TTY gate,
publish, symlink hardening, and table-column preservation are untested.

## Open questions

1. Is Windows a supported target? If so, the build break is CRITICAL and must be fixed before
   merge. If not, the `bubbletea` Windows support and the lack of a `!linux && !darwin` guard
   should be documented as an intentional platform restriction.

2. Should the symlink and pre-existing-directory hardening be fixed in this cycle, or is the
   threat model "only the user can create files under `~/.parley/protocol/core/`"? If the latter,
   that assumption should be stated, because the write-once guarantee is currently stated
   unconditionally.

3. The §7 clause says "An idea that is already open completes under the protocol version it was
   pinned to." Should this be hedged to future tense ("will complete") until rank 2 ships, or is
   it intended as a protocol-text promise that the implementation will eventually satisfy?
