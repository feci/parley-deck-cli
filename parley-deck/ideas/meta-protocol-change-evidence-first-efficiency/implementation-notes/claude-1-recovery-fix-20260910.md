---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-10
phase: implementation (recovery fix for the two reproduced failures in my Phase-5 slice)
slice: internal/protocolpacket/** (this session touched source.go and packet_test.go only)
status: written on branch; NOT executed by me (no shell this session); pending facilitator execution and independent review
supersedes: nothing — implementation-notes/claude-1.md and claude-1-resume-20260909.md are unchanged and still describe the rest of the slice
---

## Scope of this session

Only `internal/protocolpacket/source.go`, `internal/protocolpacket/packet_test.go`, and this new
file. No Git, no other owner's code, no FINAL/scope/protocol/metadata, no `internal/app/**`, no
`internal/fsutil/**`. Earlier artifacts are byte-preserved. **No test ran** — the shell was
unavailable — so everything below is by inspection, and I issue no verdict on my own slice (§15.1).

## 1. `ensurePrivateDir`: Lstat-missing then Mkdir EEXIST (reproduced failure)

`ensurePrivateDir` returned "cannot create" whenever it lost the create race, which is why
`-run '^TestConcurrentPublication' -count=20` failed repeatedly.

EEXIST is now handled, but not trusted. It proves an entry exists; it does not say the entry is
the private directory this call asked for — a symlink or a plain file planted in that window
produces the same errno. The branch re-`Lstat`s the path and runs it through the *same*
`checkPrivateDir` a pre-existing path goes through: symlink → error, non-directory → error,
group/world bits → tightened, still loose → error. Validation was extracted into
`checkPrivateDir(path, fi)` so both the pre-existing and the raced path provably share it; the
`fi` handed in is always an `Lstat` result, so a symlink is judged as a symlink, never as its
target. Non-EEXIST `Mkdir` errors still fail immediately.

## 2. `.staged-packet` flush on the shared volume (reproduced failure)

`writeAndClose` called `f.Sync()`, which on Darwin is `F_FULLFSYNC` and which this volume rejects
with ENOTTY. It now calls `fsutil.SyncFile(f)` — imported, not copied — which falls back to
`syscall.Fsync` for that errno only and still requires it to succeed. A flush failure remains a
publication failure; nothing is suppressed. New import: `parley-deck-cli/internal/fsutil`.

## 3. The O_EXCL fallback is gone; `verifyPublished` no longer follows symlinks

You are right on both counts, and I withdraw the fallback rather than defend it.

- **Fallback removed.** `O_WRONLY|O_CREATE|O_EXCL` at the target never clobbers, but it publishes
  an empty file and fills it afterwards: a launch reading in that window gets an incomplete body
  under a complete attestation, which is the defect the slice exists to close. Calling that
  "support" for link-less filesystems was weakening publication to claim coverage. Publication is
  now stage + `link()` only; any non-EEXIST link error fails closed with an error naming the
  directory and the reason ("does not support the hard link that makes publication indivisible").
- **`verifyPublished` no longer uses `os.ReadFile`.** It goes through the new `openPublished`:
  `Lstat` (symlink or non-regular → error), `Open`, then `f.Stat()` on the descriptor, requiring a
  regular file *and* `os.SameFile` against the `Lstat` result. A symlink swapped in after the
  `Lstat` is followed by `open()`, so the descriptor identifies a different file and the read is
  refused. Bytes are only ever compared from a descriptor that was, at one instant, the regular
  file at that path.
- **The link-EEXIST caller no longer bypasses validation.** It previously reached `verifyPublished`
  without `existingBody`'s checks; validation now lives inside `verifyPublished`, so both entries
  (pre-check and race) are covered. `existingBody` stays as the early, better-worded pre-check.
- `writeAndClose` now sets `0600` via `f.Chmod` on the descriptor instead of `os.Chmod` on the
  staged path (it is the only caller left). Equivalent exposure, one less path lookup.

**Honest limitation (unchanged by this fix):** the guarantee is against *another user* and against
path/symlink swaps, not against the same UID. Someone who can write inside the `0700` directory —
the same user, or root — can modify the inode's bytes after `openPublished` validated it and
before a launch reads it. No sequence of stat-and-open closes that; only the directory permission
gate and a reader that re-hashes the file against the full digest in its name do.

## 4. Tests added (all in `packet_test.go`, none executed)

- `TestConcurrentPublicationCreatesThePrivateRuntimeDirectoryOnce` — 10 rounds × 12 goroutines
  released together on a fresh root; every `publicationDir` must succeed *and* both directories
  must afterwards be real, non-symlink, non-group/world-accessible. It matches
  `^TestConcurrentPublication`, so your existing `-count=20` command covers it (200 storms).
- `TestVerifyPublishedRefusesASwappedSymlinkAndDifferentBytes` — deterministic: identical bytes
  validate; different bytes are refused and the file is not modified; a symlink is refused **even
  when its target holds exactly the attested body** (this is the case the old `ReadFile` accepted);
  a directory is refused.
- `TestPublicationFailsClosedWhenHardLinksAreUnsupported` — injects a link failure through the new
  `linkFile` seam (`var linkFile = os.Link`, mirroring `fsutil`'s "seams, overridable in tests
  only" idiom) and asserts publication errors, `BodyPath` stays empty, and the runtime directory
  is left with **zero** entries — no partial body, no staged leftover.
- `TestStagedBodyIsFlushedOnThisFilesystem` — calls `writeAndClose` in `t.TempDir()`, so it fails
  on a TMPDIR whose filesystem rejects `F_FULLFSYNC` and passes with the fallback; also pins `0600`
  on the staged file before it is ever linked.

Unchanged: default full context + shadow packet, full 64-hex digest naming, `0700`/`0600`
permissions, the public API, and the parser. No new dependency (only `io` + `internal/fsutil`).

## What is covered by inspection only

1. **Nothing was executed.** Compile and assertion errors are possible.
2. The EEXIST *branch* is exercised probabilistically by the storm test; I did not inject an
   artificial race into `Mkdir`, so its coverage rests on the storm plus the fact that it calls the
   same `checkPrivateDir` the symlink/file tests pin deterministically.
3. Two things could newly misbehave specifically on `/Volumes/My Shared Files/...` and I would
   rather name them than have them look like a mystery: `os.SameFile` needs stable dev/ino across
   `Lstat` and `fstat` (if the mount does not provide that, republication fails with "was replaced
   while it was being validated"), and the `0600` chmod is now `fchmod` rather than path `chmod`
   (a failure surfaces as "cannot restrict .staged-packet-…"). Both are loud and fail closed.
4. Items 4–7 of the limitations in `claude-1-resume-20260909.md` still stand, as do the open
   integration dependencies (the renderer is still not on a launch path; attestation is still not
   wired into telemetry).

## Files

- `internal/protocolpacket/source.go` (modified)
- `internal/protocolpacket/packet_test.go` (modified)
- `parley-deck/ideas/.../implementation-notes/claude-1-recovery-fix-20260910.md` (new, this file)

## Recommended verification

```sh
gofmt -l internal/protocolpacket
go vet ./internal/protocolpacket/
go test ./internal/protocolpacket/ -run '^TestConcurrentPublication' -count=20   # failure 1
TMPDIR='/Volumes/My Shared Files/tmp' go test ./internal/protocolpacket/ \
  -run 'TestPublishedBody|TestConcurrentPublication|TestStagedBody' -count=1     # failure 2
go test ./internal/protocolpacket/ -run 'VerifyPublished|FailsClosed|Private|Symlinked' -count=1
go test -race -count=4 ./internal/protocolpacket/
go test ./internal/protocolpacket/ ./internal/app/ ./internal/protocol/ ./internal/fsutil/
go test ./...
```

Adversarial checks worth doing by hand, since the tests are mine: replace a published body with
different bytes and confirm the CLI exits 1 without overwriting; replace one with a symlink to a
file holding the identical body and confirm publication is still refused; `chmod 0777` the runtime
directory and confirm it is tightened before use.
