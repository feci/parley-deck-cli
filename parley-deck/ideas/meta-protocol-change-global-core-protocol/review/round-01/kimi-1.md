---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 1
date: 2026-08-07
reviewed-commit: 4396529
---

## Summary

The slice ships `internal/protocolcore` (release store + pure renderer), the `parley protocol
status|render|check|publish` command, and the §7 blast-radius clause. The happy path is genuinely
good: I rendered THIS deck against the REAL embedded default core in an isolated fixture and the
only change was the synced-stamp; the run is idempotent and `check` then exits 0. The `/dev/null`
TTY defect narrated in IMPLEMENTATION.md is really fixed (verified: refusal, exit 2). `go build
./... && go test ./...` is green on darwin (full output below).

But the adversarial brief was justified. I found the same "documented as landed and wrong at the
call site" pattern this review was told to hunt, in four places:

1. **The Windows build is broken** — IMPLEMENTATION.md claims "a platform without the ioctl
   refuses instead of accepting /dev/null"; in reality a platform without the ioctl **does not
   compile**, and this repo ships Windows binaries (`dist/` v1.37.0–v1.40.1, x64+arm64).
2. **G1's removal report is heading-granular.** Deck-local content under a *shared* heading, and a
   roster table whose header is not exactly `| Agent ID …`, are erased with **no report in preview
   or apply** — the librade-algoTrader failure one level down, reproduced against fixtures.
3. **The TTY gate is bypassed by a stock one-liner** (`python3 -c 'import pty; pty.spawn(…)'` —
   published successfully). §7's "An agent may not — not by publishing one" is wrong at the call
   site, and the refusal has no test (G7b).
4. **§7 ships an unimplemented guarantee.** "An idea that is already open completes under the
   protocol version it was pinned to" (parley-deck/COOPERATION.md:766-767) — rank 2 is not built
   (verified: zero pinning/snapshot code), no test exists, and IMPLEMENTATION.md claims "nothing
   here claims a guarantee it does not implement (G7b)" in the same commit that ships the claim.

Plus a real path-traversal: `Store.Load` does not validate the version string (Publish does), so a
committed deck lock of `core-version: ../../../../trav/x` makes `render`/`check` consume a
COOPERATION.md outside the store — demonstrated end to end.

Counts: 1 CRITICAL, 6 MAJOR, 3 MINOR, 3 NIT.

## Findings

### [CRITICAL] GOOS=windows (and other non-termios platforms) no longer build; IMPLEMENTATION.md's fallback claim is false

`GOOS=windows go build ./...` at 4396529:

```
# parley-deck-cli/internal/app
internal/app/protocol.go:312:21: undefined: unix.IoctlGetTermios
internal/app/protocol.go:312:50: undefined: termiosGet
EXIT=1
```

Only two build-tagged files define `termiosGet`: `internal/app/termios_unix.go:1`
(`darwin || freebsd || netbsd || openbsd`) and `internal/app/termios_linux.go:1` (`linux`). There
is no file for windows, dragonfly, solaris, plan9, or wasm, and `protocol.go` has no build tag —
so the whole `internal/app` package fails on those platforms. `dist/` contains
`parley-v1.37.0`–`v1.40.1` `windows-x64.exe` and `windows-arm64.exe`: Windows is a shipped target,
and this commit breaks its build. (darwin and linux build fine, verified.)

Worse, IMPLEMENTATION.md:61-62 documents the opposite: *"the platform fallback was deleted rather
than kept weak: a platform without the ioctl refuses instead of accepting `/dev/null`."* No such
refusing fallback exists — the platform does not compile, so nothing refuses. This is the exact
documented-as-landed-and-wrong pattern, in the paragraph that narrates fixing another instance of
it. Fix is mechanical: add `termios_unsupported.go` with
`//go:build !darwin && !freebsd && !netbsd && !openbsd && !linux` defining `hasTTY() bool { return false }`
(which *is* the "refuse" posture the text describes).

### [MAJOR] G1: deck content under a shared heading is silently erased — no report in preview or apply

`removedSections` (internal/protocolcore/render.go:184-205) compares only `##`/`###` heading lines
between deck and core. Anything else the deck has and the core does not — prose, lists, tables
under a heading the core also has — is dropped **unreported**.

Fixture (verbatim, run against the built binary): deck whose `## 3. Phases` section adds

```
Project-specific addition: ALWAYS run `make gen` before Phase 5. This sentence is genuine
local protocol content that no other deck has.
```

Preview output reports only `preserved from this deck: Workspace, Transport, Created, §2 roster
table` — nothing about the paragraph. Apply reports the same. Afterwards:
`grep -c "make gen" → 0`. The content is gone and neither preview nor apply said a word.

This is the 2026-08-06 librade-algoTrader erasure the package comment itself cites
(render.go:30-33, 181-183) reproduced at sub-section granularity. G1 reads "MUST report every
block it replaces or removes, in preview and on apply" — the implementation defines "block" as
"heading the core lacks" and never says so. Either removal detection must diff section *content*
(per-section hash), or the gate/CLI output must state plainly that only whole sections are
reported. As shipped, the report creates false confidence exactly where the historical incident
lived.

### [MAJOR] G1: a roster table whose header is not literally `| Agent ID …` is silently emptied

`tableRows` (render.go:159-178) extracts only the first table whose header starts with
`| Agent ID` and contains the marker. Fixture: a deck with a legacy-shaped header
`| Agent          | Workspace dir                       | Role          |` and two member rows.
Render output:

```
preserved from this deck: Workspace, Transport, Created        ← roster table NOT preserved
```

and the written file's §2 is an **empty table** — both `claude-1` and `codex-1` rows gone
(`grep -c "claude-1\|codex-1" → 0`), nothing reported as removed (the `## 2.` heading exists in
the core, so finding #2's heading check does not fire either). The roster idea's own changelog
(1.41.0) states the standard this project set for itself: "§2-only IDs are reported, never
silently erased … `roster render` reports every row it removes — in preview and on apply." This
renderer silently erases the whole roster on a header variant. With "nine distinct rosters"
measured across 40 decks, header variance is not hypothetical. The six identity slots (D3) are
only preserved when the deck's markup matches the heuristic byte-for-prefix.

### [MAJOR] G1: CRLF decks get a false-positive removal report for every section, and mixed-ending output

`heading()` (render.go:207-212) does `strings.TrimRight(line, " \t")` — `\r` is not trimmed. On a
CRLF deck every heading line carries `\r`, matches nothing in the core map, and is reported as
removed. Fixture: a deck byte-identical to the core except CRLF endings + one local section.
Preview (and apply) printed:

```
the following section(s) exist in this deck but NOT in core 1.0.0 and will be REMOVED:
  - ## 2. Active agents (roster)
  - ## 3. Phases
  - ## 4. Reviews
  - ## 9. Local note
```

The first three exist in the core and are NOT removed — the report is wrong about them; only
`## 9.` is a real removal. After apply, `file` reports the output as
`ASCII text, with CRLF, LF line terminators` (identity lines keep their `\r`, core lines are LF).
Byte-idempotence still holds (second apply is a no-op), and the CRLF→mixed conversion arguably
normalizes over time — but G1's report, the thing the gate exists to make trustworthy, cries wolf
on exactly the decks Windows-authored tooling produces. One-line fix: trim `"\r \t"`.

### [MAJOR] G2/G7b: the TTY gate is bypassed by a stock pty one-liner; the refusal guarantee has no test

First, credit: the narrated defect is fixed. `parley protocol publish --version 2.0.0 --from … <
/dev/null` now refuses with exit 2.

But `hasTTY` (protocol.go:310-317) asks for terminal attributes on stdin/stdout, and **any
pseudo-terminal satisfies it**. From this machine, with no user present:

```
$ python3 -c 'import pty; pty.spawn(["parley","protocol","publish","--version","2.0.0","--from","newcore.md"])'
Published core 2.0.0 (175d4c63f7ee) to …/protocol/core/2.0.0
```

The child gets a real controlling terminal; the ioctl succeeds; the release is written (mode 0444,
verified). `pty.spawn` tolerates non-tty stdin, so this works from any agent runtime that can spawn
processes — i.e. all of them — and likewise via `expect`, `unbuffer`, `tmux new-session -d`, or an
agent launched inside tmux where stdout is already a tty. The code comment's premise ("an agent run
(which has no TTY) cannot reach it", protocol.go:269-271) and the §7 text it backs — "An agent may
not — not by editing a release, not by publishing one: releases are write-once and `parley
protocol publish` refuses without a controlling terminal" (parley-deck/COOPERATION.md:760-762) —
are wrong at the call site in the general case. Consensus VC-1 already resolved this honestly:
prevention is real only for sandboxed parley-launched agents (DF-1, unshipped); for unmanaged
agents the model is detection. The shipped protocol text claims the stronger thing VC-1 rejected.

And per G7b: the refusal behavior is documented in the usage string (protocol.go:24), §7, and
IMPLEMENTATION.md:56-62, yet **no test asserts it**. The test process itself has no TTY, so
`runProtocol([]string{"publish", "--version", "v", "--from", f"})` expecting exit 2 is trivially
writable. A guarantee documented as landed without an end-to-end test is precisely what G7b
forbids — doubly so when the guarantee is also bypassable.

### [MAJOR] Load() path traversal: the committed deck lock controls an unvalidated filesystem path

`Store.Publish` validates the version (core.go:97-99 rejects `/`, `\`, `..`). `Store.Load`
(core.go:55-65) does not — it `filepath.Join`s the deck-supplied version straight onto the store
root, and `filepath.Join` cleans but does not contain. The version comes from the committed deck
file `parley-deck/meta/protocol-lock.yaml` (protocol.go:79-93). Demonstrated end to end:

```
core-version: ../../../../trav/x        ← in the deck lock
$ parley protocol render --dir deck --yes
Wrote …/COOPERATION.md from core ../../../../trav/x (7ad244f9097f)
$ grep -c "OUTSIDE-THE-STORE" deck/parley-deck/COOPERATION.md → 1
$ parley protocol check --dir deck
protocol check: in sync with core ../../../../trav/x (7ad244f9097f)   (exit 0)
```

`render` regenerated the deck's protocol from a file living outside the store, and `check` then
certified it "in sync". The blast radius is bounded (read-only; the target must be named
COOPERATION.md; the result is a committed, reviewable file), but the render path is
agent-reachable and the lock is deck content — a malicious or corrupted deck steers a filesystem
read. `Load` should share Publish's validation. The authors knew versions need validating; the
read path was missed.

### [MINOR] Publish's write-once guarantee is per-file, not per-release; version "." and symlinks misbehave

All demonstrated against the real binary (via the pty wrapper above):

- `--version .` is accepted (the `/\..` filter doesn't catch it) and writes `COOPERATION.md`
  directly into the **store root**, mode 0444 — invisible to `Versions()` (which lists dirs only,
  core.go:70-86) but loadable by a deck pinning `core-version: .`.
- The write-once check is `os.Stat(<dir>/COOPERATION.md)` only (core.go:102-105). A pre-existing
  empty `3.0.0/` dir, or one containing a stray `registry.json`, is **reused**: publish succeeds
  and the "release" ends up containing bytes that were never published. D1's model is that a
  version label maps to exactly one immutable content set; the refusal should fire on the *dir*
  existing, not just the file.
- A version dir planted as a **symlink** to an outside directory: Stat misses (no COOPERATION.md
  inside), `MkdirAll` no-ops, `WriteFile` follows the link — the release is written **outside the
  store root**. This project's own convention (`parley learn`: "Lstat symlink guard") exists for
  exactly this. All three need local access via the attended path, hence MINOR — but they are the
  write-once guarantee's edges, and D1 is the load-bearing decision of this idea.

### [MINOR] "all three COOPERATION.md copies" — the third was not changed and no inbox flag tracks it

IMPLEMENTATION.md:44-46 claims the §7 clause was "added to all three `COOPERATION.md` copies".
The diff touches two: `parley-deck/COOPERATION.md` and `internal/protocol/defaults/COOPERATION.md`.
The third class — the external skill fallback snapshots agents actually read — is unsynced,
verified on this machine:

```
~/.claude/skills/parley-deck/references/COOPERATION.md:    0 blast-radius clause(s)
~/.kimi-code/skills/parley-deck/references/COOPERATION.md: 0 blast-radius clause(s)
```

The project has a precedent for exactly this gap (review-gate-honesty flagged the external sync in
`parley-deck/inbox/claude-to-all_review-gate-honesty_external-skill-snapshot-sync.md` and noted it
in its changelog entry). This change does neither: no inbox file, and the G6 changelog entry
(meta/protocol-changelog.md:1-25) does not mention the pending external sync. The sentence in
IMPLEMENTATION.md is wrong as stated.

### [MINOR] `protocol render` ignores deck-file read errors — an unreadable deck file is silently replaced

protocol.go:172: `prior, _ := os.ReadFile(path)`. If the deck file exists but cannot be read
(EPERM, transient I/O error), `prior` is empty, the render proceeds as "fresh deck" (no identity,
`Removed` nil), and `--yes` atomically replaces the file — with no error and no removal report.
The check path handles the same read correctly (protocol.go:224-228 returns 1); render should too.

### [NIT] "Nine tests drive the real command entry points" — eight do

IMPLEMENTATION.md:64-65. `TestCoreReleasesAreWriteOnce` (protocol_test.go:246-262) calls
`Store.Publish` directly; it never touches `runProtocol`. The write-once guarantee it pins is real,
but the sentence overclaims the entry-point coverage — and elides that the publish entry point has
no test at all (see the G2 finding).

### [NIT] Two generators now own the `**Protocol synced:**` line with different formats

The new stamp is `**Protocol synced:** core <version> (<sha12>)` (render.go:87). The legacy §9.0
preflight sync still rewrites any `**Protocol synced:**` line to
`**Protocol synced:** <YYYY-MM-DD> — parley-deck-skill <v> (preflight)` (preflight.go:583-599,
date-stamped from `time.Now()`). A deck that adopts a core and later undergoes the legacy sync
gets its stamp rewritten, after which `protocol check` reports `hand-edited-or-stale` and exit 1.
No deck has adopted yet (DF-2), so this is latent — but the migration order (or retiring the
preflight writer) needs to be decided before DF-2, and the date/skill provenance the old stamp
carried is dropped by the new format with no documented replacement.

### [NIT] Command docs and root changelog not updated

`parley protocol` is wired into dispatch and `--help` (verified in app.go diff and the binary),
but `docs/cli-reference.md` does not mention it and root `CHANGELOG.md` has no entry. The G6
changelog (meta/protocol-changelog.md) IS present and in the §7 template format — that gate
passes; this is the user-facing docs trail only.

## Gate-by-gate assessment

**G1 — FAIL.** Idempotence holds everywhere I tested (fixtures A–D and the real-deck render all
no-op on second apply; the stamp is release-derived, not clock-derived). Removal reporting works
for the whole-section case in BOTH preview and apply (test + fixture C confirm). But the report is
not complete: sub-heading content (finding 2), non-`| Agent ID` tables (finding 3), and CRLF false
positives (finding 4) each produce a wrong or silent report. The gate says "every block it
replaces or removes"; the implementation reports a subset and calls it every.

**G2 — FAIL (partial).** The write-once half holds through `Store.Publish` for the same-version
case (tested). Everything around it leaks: the attended gate is pty-bypassable (finding 5), `Load`
traverses (finding 6), publish's existence check is per-file with symlink and "." edges (finding
8). G2 also demands "a guard test asserts both halves" — only the write-once half has a test.

**G3 — N/A (correctly).** No surface in the new code reports confinement (`grep confinement
internal/` hits only pre-existing `internal/agents` code). IMPLEMENTATION.md and the changelog
state `confinement-unproven`/detection-only as the posture without shipping any confinement claim.
No violation.

**G4 — N/A (correctly).** No overlay code shipped; nothing claims overlay behavior.

**G5, G7, G8 — N/A this cycle (correctly absent in code; overclaimed in protocol text).** Verified
rank 2 does not exist: `grep` for `protocol-snapshot|stale-protocol|pinned protocol|EffectiveHash`
across `internal/` returns nothing. IMPLEMENTATION.md scopes this honestly ("so callers can
implement D8", line 26; "Ranks 2-4 … remain ratified and scheduled", line 15-17). But the §7 text
this same commit ships tells every agent "An idea that is already open completes under the
protocol version it was pinned to" (parley-deck/COOPERATION.md:766-767, mirrored at
internal/protocol/defaults/COOPERATION.md:757-758) — a present-tense behavioral guarantee with no
implementation and no test, in a paragraph that correctly qualifies the overlay with "once that
ships". That is the G7b violation; FINAL.md:53-54 ("Nothing in the shipped slice may claim a
guarantee it does not implement") is binding and is breached by the shipped protocol text itself.

**G6 — PASS.** meta/protocol-changelog.md:1-25 carries the entry in the §7 template format
(`Idea:` path, `Drafted by:`, `Summary:`), plus an enforcement paragraph that is, notably, more
honest than the §7 clause it documents (it scopes prevention to parley-launched participants and
names the unresolved-path limit — the §7 text should read like this).

**G7b — FAIL.** Named guarantees without an end-to-end test of the real entry point:

1. **TTY refusal** — documented in usage text (protocol.go:24), §7 (parley-deck/COOPERATION.md:761-762),
   IMPLEMENTATION.md:56-62; no test exists (trivially writable: the test env has no TTY).
2. **Open-idea pinning** — §7:766-767; no implementation at all, so no test can exist.
3. **"reports every block it replaces or removes"** — tested only for whole-section removals;
   fixtures B and D show the guarantee false at finer granularity, so the existing tests encode
   the happy-case assumption rather than the guarantee.
4. **"the write is atomic and preserves the file mode"** (IMPLEMENTATION.md:37) — the
   mode-preserving branch (protocol.go:205-208) is untested at the entry point (fsutil's atomicity
   is covered elsewhere; the render-side perm handoff is not).
5. **`check --json` non-zero exit on mismatch** — set in code (protocol.go:240-243), asserted
   nowhere; the text-mode exit is tested.

Passes: block-on-missing-release, refuse-to-guess, JSON status shape, check-reports-never-writes,
idempotence — all asserted through `runProtocol` against fixtures, exit codes and file state
included. That half of G7b is genuinely met.

## Test-quality assessment

All nine new tests drive behavior, not implementation trivia; I found no tautology. Each was
checked against "would it fail if its fix were reverted":

- `TestProtocolRenderIsIdempotent` — real; a clock-derived stamp or any non-idempotent write fails
  it, and it pins the no-op message on second apply. It would have caught the class of bug G1
  names.
- `TestProtocolRenderPreservesIdentityAndReplacesCoreText` — real; also pins that the core's
  `example-1` placeholder row does not leak into the deck (the fixture-D case in miniature — but
  only for a *matching* header).
- `TestProtocolRenderReportsRemovedSections` — real for whole sections, preview AND apply.
  **Coverage gap, not a tautology:** it encodes the assumption that removal == deck-only heading,
  so it cannot see findings 2–4.
- `TestProtocolCheckReportsHandEditAndNeverWrites` — real; asserts exit code, message, and file
  bytes unchanged.
- `TestProtocolBlocksWhenPinnedReleaseIsMissing` / `TestProtocolRefusesToGuessAVersion` — real;
  both assert refusal AND that the deck file is untouched.
- `TestCoreReleasesAreWriteOnce` — real, Store-level (mislabeled by IMPLEMENTATION.md as
  entry-point coverage; see NIT).
- `TestProtocolStatusJSON` — parses real CLI JSON and asserts fixture-derived values; thin but
  behavioral.
- `TestProtocolRenderOnAFreshDeck` — real; usefully pins the "empty prior ⇒ no false removal
  report" branch (the fixture-A defect class for an empty deck).

Missing tests, in order of importance: TTY refusal (G7b); `Load` version validation; CRLF
idempotence/report correctness; roster preservation under header variance; sub-heading removal
reporting; publish `--version .`/pre-existing-dir refusal. Every one of these is writable today
against the real entry points with the same fixture pattern the file already uses.

## What I verified that is GOOD (so the fix-up keeps it)

- Real-deck render against the real embedded core (`internal/protocol/defaults/COOPERATION.md` as
  release 1.0.0, this repo's deck copied into a fixture): all six identity slots preserved
  (`Workspace, Transport, Created, §2 roster table, host-handle table`), the ONLY diff is the
  synced-stamp line, second apply no-ops, `check` exits 0.
- `publish < /dev/null` refuses with exit 2 and a clear message — the narrated defect is fixed.
- The renderer is genuinely pure: render.go imports only `fmt` and `strings` — no fs, no clock, no
  env; the stamp is derived from the release (`render.go:87`, `core %s (%s)` with version +
  sha12), and `Load` computes the sha from the bytes read (core.go:64).
- `Publish` rejects `/`, `\`, `..` in versions (core.go:97-99) — the write side is not the
  traversal hole.
- Full suite at 4396529: `go build ./...` OK, `go vet` clean on both touched packages,
  `go test ./...` → all `ok` including `internal/app` (the nine new tests pass individually,
  `-run 'TestProtocol|TestCore' -v`), exit 0. Nothing pre-existing was broken **on darwin/linux**
  (Windows excepted — finding 1).
- D8 blocking is real and tested: missing pinned release names the version and leaves the deck
  file untouched; unpinned deck is refused with "a version is never guessed".

## Open questions

1. What is "a block" for G1 reporting before the rank-3 registry lands — whole section, or section
   content? If whole section, the CLI output and G1's text should say so; if content,
   `removedSections` needs per-section hashing. Today it silently implements the weaker reading
   while documenting the stronger one.
2. Is Windows still a shipped target? If yes, finding 1 is release-blocking for the next version
   cut; if support is deliberately dropped, `dist/` and any release docs need to say so.
3. Should the §7 sentence (parley-deck/COOPERATION.md:760-762) be reworded to VC-1's honest
   posture — "attended-publisher heuristic + detection for unmanaged agents" — instead of "An
   agent may not … publish one", which a pty defeats? The changelog entry's own enforcement
   paragraph already has the right wording.
4. Should `Load` share `Publish`'s version validation (and additionally reject `.`)? The lock is
   committed deck content; anything read from it should be treated as untrusted input.
5. Before DF-2: which writer owns the `**Protocol synced:**` line — preflight's date-stamped
   legacy format or the renderer's `core <v> (<sha>)`? And is the date/skill provenance the old
   stamp carried (e.g. `parley-deck-skill 2.5.1 / parley-deck-cli 1.41.0`) intentionally dropped,
   or should the lock/stamp retain it?
