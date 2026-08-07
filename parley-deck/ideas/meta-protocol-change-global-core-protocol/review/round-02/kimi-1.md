---
agent: kimi-1
idea: meta-protocol-change-global-core-protocol
review-round: 2
date: 2026-08-07
reviewed-commit: 8888e00
verdict: FINDINGS
---

## Summary

Most of cycle 1 is real and I verified it behaviorally: all three builds compile, the full suite
is green on a fresh `-count=1` run, the store now refuses every symlink/existing-dir/traversal
variant I threw at it, the removal report is genuinely content-based, and CRLF *decks* are handled.
The six new tests are mostly real — I reverted each fix in a throwaway tree and watched them fail.

But the cycle has one hole that matters, and it is the same pattern this idea keeps shipping:
**the honesty fix was applied to the wrong documents.** The CLI refusal text and IMPLEMENTATION.md
now admit a pty defeats the gate and that rank 2 is unshipped — but the §7 protocol text, the thing
every agent actually reads, is byte-unchanged in **all three copies** and still says "An agent may
not — not by editing a release, not by publishing one" and "An idea that is already open completes
under the protocol version it was pinned to". IMPLEMENTATION.md:26-27 now even asserts "nothing in
this slice claims it", which is false at the shipped text. Documented-as-corrected, wrong at the
call site — in the fix-up's own honesty paragraph.

Around that: two of the six new tests are provably weaker than their names (one passes with its
fix reverted; one never reaches the code it claims to cover — both demonstrated), and the new CRLF
handling normalizes the deck but never the *core*, producing mixed-ending output, a two-render
convergence, and CRCRLF-mangled files depending on the input pairing.

Command output (this machine, commit 8888e00):

```
$ go build ./...            EXIT=0
$ GOOS=windows go build ./...  EXIT=0
$ GOOS=linux go build ./...    EXIT=0
$ go test ./... -count=1    all ok (no failures), EXIT=0
$ go test ./internal/app -run 'TestProtocol|TestPublish|TestCore' -count=1 -v
  15 tests, 15 PASS
```

Counts: 1 MAJOR, 6 MINOR, 2 NIT.

## Round-01 findings: fixed or not

1. **[CRITICAL] Windows build broken — FIXED.** `GOOS=windows go build ./...` and `GOOS=linux`
   both exit 0 (output above). `termios_other.go` restores the fallback, and it is genuinely
   fail-closed: `hasTTYSupported = false` makes `protocolPublish` exit 2 with "unavailable on this
   platform" before any TTY question is asked (protocol.go:285-289, termios_other.go). The Windows
   refusal behavior is verified by inspection only — I cannot run a Windows binary here — but the
   code path is a compile-time constant and the build is verified.

2. **[CRITICAL] Publish followed symlinks / accepted an existing release dir — FIXED for every
   vector round 1 demonstrated, with two residual edges (new findings 2 and 3).** Verified against
   the real binary via a pty wrapper: pre-existing empty release dir → refused; a *symlink* as the
   release dir → refused (Lstat sees the link itself, `/tmp/.../outside` stayed empty); a regular
   file at the release path → refused; republish of an existing release → refused; `--version .` →
   rejected. The per-release write-once (core.go:135-140) plus `O_EXCL|O_NOFOLLOW` (core.go:148) is
   the right shape.

3. **[CRITICAL] TTY gate pty bypass — as claimed: not fixed, and the honesty correction is real
   but incomplete.** Verified: `< /dev/null` still refuses with exit 2, and the new refusal text
   (protocol.go:290-296) says plainly that a pty defeats it and that confinement is DF-1. A pty
   wrapper still publishes (I re-did it this cycle) — which the text now admits, so that is the
   agreed posture. **But the correction never reached the §7 protocol text** — see new finding 1.
   And the refusal still has no test (round-01 G7b point, unaddressed): nothing in the test file
   drives `publish` through `runProtocol`/`Run`; this is trivially writable since the test process
   has no TTY.

4. **[CRITICAL/MAJOR] Heading-based removal report / CRLF false positives — FIXED, with new
   defects in the replacement (new findings 4, 5, 6).** Verified against fixtures: a local
   paragraph under a heading the core also has is now reported (`## 3. Phases — 1 line not carried
   forward`); the legacy-header roster table's dropped rows are now reported (`## 2. Active agents
   (roster) — 4 lines not carried forward`) instead of vanishing silently; a CRLF deck produces no
   false report, uniform CRLF output (22 CRLF, 0 lone LF/CR), an idempotent second apply, and
   `check` exit 0. A deck line duplicated elsewhere in the output is not reported as dropped — but
   that line is still *in* the document, so nothing is actually lost; I judge the line-set
   semantics acceptable and do not count it as a finding.

5. **[MAJOR] Path traversal via the committed lock — FIXED.** 19 adversarial lock values through
   the real binary (`../../../../etc`, `..`, `.`, `a/b`, `..x`, `x..y`, `...`, `.hidden`, `+`, `-`,
   `1.0.0.`, `CON`, `1%2e0`, unicode digits, `1:0`, `a,b`, whitespace variants): every traversal or
   oddity rejected fail-closed ("unsafe version" / "release not installed"); CRLF and tab lock
   files parse correctly. `ValidVersion` (core.go:48-68) now guards both `Load` and `Publish`.
   **However, the test written to pin this is vacuous** — new finding 7.

6. **[MAJOR] Tests bypassed production dispatch — FIXED and the new test is real.**
   `TestProtocolIsReachableThroughProductionDispatch` (protocol_test.go:405) drives `Run`. Revert
   check: deleting `case "protocol"` from `app.go` in a scratch tree makes it fail
   (`dispatch exit=2: unknown command: protocol`).

7. **[MAJOR] IMPLEMENTATION.md claimed rank-2 continuation — NOT FIXED where it matters.**
   codex-1's finding said explicitly: "IMPLEMENTATION.md itself does not overclaim … The overclaim
   is in the protocol text shipped by this commit", and the suggested fix was to remove the
   sentence from the shipped protocol copies. The cycle corrected IMPLEMENTATION.md and left all
   three protocol copies untouched. See new finding 1.

8. **[MINOR] --help omitted publish — FIXED.** `parley --help` now lists
   `protocol publish --version V --from FILE (attended; requires a terminal)`; verified in binary
   output.

9. **[MINOR] status swallowed read errors — FIXED, no test.** Verified: unreadable store →
   `protocol status: reading …: permission denied`, exit 1; unreadable lock → same, exit 1
   (protocol.go:132-144). No test pins either branch.

10. **[MINOR] Third protocol copy — FIXED.** `parley-deck-skill@455aafe` exists in the sibling
    repo and adds the §7 clause to `skills/parley-deck/references/COOPERATION.md`; the §7
    blast-radius block is byte-identical across the deck copy, the embedded default, and the skill
    copy (diff-verified). (Identical also means the skill copy carries the same overclaim —
    finding 1.)

**Still open from round 01, not claimed by the fix-up** (restating so they are not lost):
`protocol render` still swallows deck-file read errors (protocol.go:180 `prior, _ := os.ReadFile(path)`);
`docs/cli-reference.md` and root `CHANGELOG.md` still silent; preflight.go:583-594 still owns the
legacy `**Protocol synced:**` format (latent until adoption); no test for the mode-preserving write
branch (protocol.go:213-216) or for `check --json`'s non-zero exit.

## New findings (by severity)

### [MAJOR] The §7 protocol text still ships both overclaims the fix-up says were corrected — in all three copies

The fix-up claims the TTY overclaim was handled by making "the refusal text and IMPLEMENTATION.md
say so" and that the rank-2 continuation claim was "corrected". Both corrections landed in
documents *about* the change and skipped the change's primary artifact. Shipped today, verbatim:

- `parley-deck/COOPERATION.md:760-762` and `internal/protocol/defaults/COOPERATION.md:751-753`
  (and the skill copy at 455aafe, byte-identical): "**Only the user may change the global core.**
  An agent may not — not by editing a release, not by publishing one: releases are write-once and
  `parley protocol publish` refuses without a controlling terminal." I published through a stock
  `pty.spawn` one-liner again this cycle; "an agent may not … publish one" remains wrong at the
  call site in the general case, and the sentence names none of the honesty qualifiers the CLI now
  carries.
- `parley-deck/COOPERATION.md:766-767` and the embedded default (`:757-758`): "An idea that is
  already open completes under the protocol version it was pinned to; the next idea in that deck
  picks up the current one." Rank 2 is still unbuilt — zero pinning/snapshot code exists — and no
  test can therefore assert this. This is exactly the sentence codex-1's round-01 MAJOR asked to
  have removed from the shipped copies until rank 2 lands.

Worse, the fix-up commit's own new text re-commits the original sin:
IMPLEMENTATION.md:26-27 — "**Continuation is rank 2 and is not implemented here**, so nothing in
this slice claims it." The slice does claim it, in the protocol text it ships. The changelog
entry's enforcement paragraph (meta/protocol-changelog.md) already has the honest wording; §7
should be brought to that standard, or hedge to future tense, in all three copies — and the
renderer's drift test will keep the two repo copies in lockstep for free.

### [MINOR] CRLF *cores* are never normalized: mixed-ending output, two-render convergence, and CRCRLF mangling

The cycle's CRLF fix normalizes the **deck** (`render.go:52-53`) and restores the deck's convention
(`render.go:96-98`), but the **release body is split raw** (`render.go:57`), so a CRLF-authored
core carries `\r` into the output. Demonstrated against the real binary (published a CRLF core as
2.2.2):

- CRLF core + LF deck: first render writes a mixed file (17 CRLF + 5 lone LF) and is **not
  idempotent** — the second render rewrites it ("Wrote …" again); G1 converges only from render 2.
- CRLF core + CRLF deck: the restore step turns the core's `\r\n` into `\r\r\n` — the output file
  contains **17 CRCRLF sequences**, and `check` then certifies that mangled file "in sync", exit 0.

The fix's own claim was "CRLF decks no longer … get mixed line endings" — true for decks, but the
renderer now owns a line-ending policy and applies it to only one of its two inputs. Normalize
`rel.Body` the same way (or reject a non-LF release at publish). This is precisely the class the
brief told me to hunt: the fix introduced something new.

### [MINOR] The removal report's preamble contradicts its new entries, and every version bump reports the old stamp as "removed project-local content"

Two report-fidelity defects introduced by the content-based rewrite:

1. The preamble was not updated for content semantics: `the following section(s) exist in this
   deck but NOT in core %s and will be REMOVED:` (protocol.go:194; same wording in `check`,
   protocol.go:268) now prefixes entries like `## 3. Phases — 1 line not carried forward` — a
   section that **does** exist in core and is **not** removed. On the shared-heading fixture the
   report literally asserts a falsehood about the section it names, in the report G1 exists to
   make trustworthy.
2. Every legitimate version bump reports the renderer's own stamp as dropped content. Deck pinned
   1.0.0 (stamp `**Protocol synced:** core 1.0.0 (df76…)`), re-rendered against byte-identical core
   1.0.1 → the *only* report line is `(document header) — 1 line not carried forward`, followed by
   `(they are project-local content; the overlay mechanism that will carry them is not shipped
   yet)`. The old stamp is regenerated metadata, not project-local content. On the routine path —
   the one every adoption will take — the report cries wolf and trains users to ignore it.
   `droppedContent` should skip `syncedPrefix` lines (the render loop already skips the core's,
   render.go:61-63).

### [MINOR] A failed publish permanently bricks the version label

`Publish` creates the release dir and file eagerly and cleans up nothing on error. Demonstrated:
`ulimit -f 0` publish → `write …: file too large`, leaving `5.5.5/COOPERATION.md` as a 0-byte 0444
file; republish is then refused — `release 5.5.5 already exists and releases are write-once` — and
a deck pinning 5.5.5 gets `release 5.5.5 has an empty body`. The version is unusable and
unfixable without manual `rm` (chmod first — it's 0444), and the error says nothing about the
leftover state. codex-1's round-01 suggestion (publish a complete temporary sibling, atomically
install) would prevent it; the cycle took the write-once half but not the atomic-install half. At
minimum, remove the partial file/dir on the error path.

### [MINOR] Symlink hardening covers only the last two path components; an intermediate symlink still escapes the store

Verified: with `$PARLEY_HOME/protocol/core` itself a symlink to `/tmp/.../outside2`, publish
succeeds and the release lands at `/tmp/.../outside2/9.9.9/COOPERATION.md` — outside the real
store root. The Lstat guards `<version>` and O_NOFOLLOW guards the final file; nothing resolves
and contains the path above them. codex-1's round-01 fix advice said "reject symlinks in every
store path component; verify the resolved target remains below the resolved store root" — that
part was not taken. Same local-access prerequisite as the original finding, hence MINOR, but the
fix-up sentence ("Now Lstats the release DIRECTORY and opens with O_EXCL|O_NOFOLLOW") reads as
complete and is not.

### [MINOR] `TestProtocolRejectsPathTraversalInTheLock` is vacuous — it passes with the Load fix reverted

Revert check in a scratch tree (`git archive 8888e00` copy, Load's `ValidVersion` call removed):

```
--- PASS: TestProtocolRejectsPathTraversalInTheLock (0.00s)
```

All four bad versions (`../../../etc`, `..`, `.`, `a/b`) resolve to files that do not exist in the
fixture, so `render` exits 1 with "release not installed" whether or not the fix is present, and
the test asserts only `code != 0` (protocol_test.go:347-361). The one test guarding the
lock-traversal fix pins nothing: it would stay green through a full revert. To make it real, plant
a file at the traversed location and assert it is *not* consumed, or assert the error names
"unsafe version". (The fix itself works — the binary refused all 19 of my lock attacks; this is
about the guard test.)

### [MINOR] `TestPublishRefusesExistingReleaseDirAndSymlinks` never exercises O_NOFOLLOW — the flag can be deleted silently

The test's symlink half pre-creates `3.0.0/` and plants a symlink *inside* it
(protocol_test.go:376-390), so `Publish("3.0.0")` is refused at the `Lstat(dir)` check before the
file is ever opened — the O_NOFOLLOW branch is unreachable in the test. Revert check: removing
`|noFollow` from the `OpenFile` flags in a scratch tree, the test still passes. The same test
*does* catch a full revert to the round-01 Stat-on-file check (fails on the pre-existing-dir
half), so the dir-refusal is pinned; the no-follow hardening — the thing the fix-up names first —
is not. (Unreachable in a unit test by construction: the dir-refusal always fires first. Either
test `OpenFile`'s flags at a lower level or drop the dead half of the test's name/setup.)

### [NIT] "GOOS=windows and GOOS=linux builds are now part of the check" — there is no check

IMPLEMENTATION.md:69. No CI workflow, no script, no `go test` hook performs the cross-builds
(`grep GOOS scripts/ .github/` → nothing); they were run manually once, and I re-ran them manually
this cycle (both exit 0). The claim overstates a manual step as an automated gate — small, but
this is the repo that keeps getting bitten by exactly that gap (cycle 1's own fix broke the
Windows build because nothing compiled it).

### [NIT] The new fixes themselves have no tests: status error reporting, Windows fail-closed, TTY refusal

`status`'s new exit-1-on-read-error branches (protocol.go:132-144), the unsupported-platform
refusal (protocol.go:285-289), and — carried from round 01 — the no-TTY refusal all have zero test
coverage. The Windows path is untestable on darwin, but the other two are trivially writable
(`chmod 000` fixture; the test process has no TTY). G7b's standard the project set for itself:
documented behavior that can be tested should be tested.

## Test-quality assessment

Six new tests this cycle. Each revert-checked by actually reverting its fix in a scratch copy of
8888e00 (`git archive`; the project repo untouched):

- `TestProtocolRenderReportsContentLostUnderASharedHeading` (:307) — **real.** Reverting
  `droppedContent` to heading-based `removedSections` fails it.
- `TestProtocolRenderHandlesCRLFDecks` (:324) — **real.** Reverting the normalize/restore fails
  both assertions (false removal report + mixed endings). (My first revert attempt silently
  failed to apply and the test passed — a reminder that revert checks must verify the revert
  actually landed, which I then did by grep before trusting any PASS.)
- `TestProtocolRejectsPathTraversalInTheLock` (:347) — **vacuous**, proven above.
- `TestPublishRefusesExistingReleaseDirAndSymlinks` (:364) — **partially real.** Fails on a full
  revert to Stat-on-file; blind to removal of `noFollow`. The symlink half is dead code as
  written.
- `TestPublishRejectsUnsafeVersions` (:393) — **real.** Reverting Publish's validation to round-01's
  weaker check fails it on `"."` and `".hidden"`.
- `TestProtocolIsReachableThroughProductionDispatch` (:405) — **real.** Deleting the dispatch case
  fails it (`unknown command: protocol`).

The nine pre-existing tests still pass and were assessed in round 01; the new suite's honest
strength is the renderer and dispatch coverage, its weakness is that the store-level guards
(traversal on Load, O_NOFOLLOW) have tests that look like they pin the fix and do not.

**G7b — still FAIL.** Guarantees named in shipped text without an end-to-end test of the real
entry point: the TTY refusal (usage text + refusal text + §7 — no test, trivially writable); the
rank-2 pinning sentence (§7, all three copies — unimplementable today); "an agent may not publish"
(§7 — false as stated); the mode-preserving render write (IMPLEMENTATION.md:38, untested);
`check --json` non-zero exit (untested); the new status-error and platform-refusal branches
(untested); O_NOFOLLOW (unpinned). The atomic-write half and the dispatch reachability now have
real tests — that part of G7b genuinely improved.

## What I verified that is GOOD (so the next cycle keeps it)

- All three builds green; full suite green on `-count=1`; the flaky `internal/runner` failure
  codex-1 saw at 4396529 did not appear here.
- `ValidVersion` refused 19/19 hostile lock values fail-closed; CRLF/tab locks parse; the version
  charset rejects unicode lookalikes, `%`, `:`, `,`, control whitespace.
- The store refused: pre-existing dir, symlink-as-release-dir, regular file at the release path,
  republish, `--version .`, empty version.
- Content-based removal reporting works for the three round-01 erasure classes (shared-heading
  prose, legacy-header roster rows, deck-only section) in preview *and* apply; CRLF decks get
  uniform CRLF output, idempotent re-render, `check` exit 0.
- The refusal text's honesty paragraph is verbatim what VC-1 resolved; the no-TTY refusal still
  exits 2.
- The third copy is committed and byte-identical in its §7 block; the drift test covers the two
  repo copies.
- The `Run`-level dispatch test is real, and IMPLEMENTATION.md's "fifteen tests" count is accurate.
