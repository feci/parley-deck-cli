---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-09
phase: implementation (Phase 5 slice hardening, continuation of implementation-notes/claude-1.md)
slice: internal/protocolpacket/**, internal/app/protocol_packet.go, internal/app/protocol_packet_test.go
status: written on branch; NOT executed by me (no shell in this session); pending facilitator execution and independent review
---

## What this note is

A continuation of my owned Phase-5 slice, not a new round-1. It hardens packet-body
publication and adds focused tests. Earlier artifacts are preserved: this file is new,
`implementation-notes/claude-1.md` is unchanged, and I edited no other owner's code, no
FINAL/scope/protocol/metadata/global settings, and no Git state. `IMPLEMENTATION.md` is
untouched this session — its claude-1 allocation row still points at the original handoff, and
this file is the current status for the same claimed paths.

I issue no verification verdict on my own slice (§15.1). Everything below is a description for
a non-owner to check against the branch. **No test was executed in this session** — the shell
was deliberately unavailable — so every statement about behaviour is by inspection, not a
measured result. The commands a verifier must run are listed at the end.

(Environment clock reported 2026-09-10 for this session; the requested artifact date
`2026-09-09` is kept in the frontmatter and filename as instructed.)

## The defect I was asked to fix

`WriteBody` published a generated body like this:

```go
name := fmt.Sprintf("%s-phase%d-%s-%s.md", ctx.ContextMode, ctx.Request.Phase,
    safeName(ctx.Request.Track), protocolcore.ShortHash(ctx.PacketSHA256)) // 12 hex chars
os.WriteFile(path, []byte(ctx.Body), 0o600)                               // O_CREATE|O_TRUNC
```

Three problems in one line pair, all of which end in "a launch is handed protocol text that is
not the attested text":

1. **Truncating publication.** `os.WriteFile` opens `O_TRUNC` and writes without atomicity. Two
   renders of the same request race on the same path, and a concurrent reader (a prompt builder
   reading `BodyPath`, or a handoff copying it) can observe an empty or partial protocol body
   while `packet_sha256` in the attestation names the complete one. Fail-open by construction.
2. **Truncated content address.** A 12-hex (48-bit) short hash in the name means the path is not
   a content address: distinct bodies can address one path, and the file's own name cannot be
   used by a reader to verify what it holds.
3. **No path validation.** `MkdirAll` + `WriteFile` follow symlinks. A symlinked
   `.parley-runtime`, a symlinked `protocol-packets`, or a symlinked target path made the
   renderer write *through* it to any file the process can reach; a pre-existing file at the
   target was silently overwritten; a group/world-writable runtime directory let another local
   user swap the body between publication and launch.

## Changes (exact)

### `internal/protocolpacket/source.go`

- Added `runtimeDirPerm = 0o700` / `bodyFilePerm = 0o600` constants with the rationale.
- **`WriteBody`** now, in order: refuses a `refused` context (unchanged); refuses an empty body;
  recomputes `Hash(ctx.Body)` and **refuses to publish a body that does not hash to
  `ctx.PacketSHA256`** (a body mutated after `Build` can no longer be published under an
  attestation that does not describe it); resolves the publication directory; builds the name
  `<mode>-phase<N>-<track>-<FULL 64-hex digest>.md`; publishes. `ctx.BodyPath` is set only after
  a successful publication. Every name component is an `int` or `safeName`-sanitised, so a
  published path cannot escape the runtime directory.
- **`publicationDir` / `ensurePrivateDir`** (new): validate `.parley-runtime` and
  `.parley-runtime/protocol-packets` with `Lstat`. A symlink is an error; a non-directory is an
  error; a missing directory is created with `Mkdir` + explicit `Chmod 0o700` (so umask cannot
  loosen it); an existing directory with `perm&0o077 != 0` is tightened to `0o700` and, if it
  stays group/world accessible, refused. `MkdirAll` is gone — the parent chain is checked
  component by component instead of created blindly. On Windows the POSIX-bit check is skipped
  (`Lstat` there reports synthetic bits regardless of the ACL); the symlink and non-directory
  checks still apply.
- **`publish` / `existingBody` / `verifyPublished` / `writeAndClose`** (new): an existing target
  is `Lstat`ed first — a symlink or non-regular file is an error, a regular file goes to
  `verifyPublished`, which accepts it **only when its bytes are exactly the attested body**
  (idempotent republication) and otherwise returns a "refusing to overwrite a published packet"
  error. A new body is staged with `os.CreateTemp` in the same directory (`O_EXCL`), written,
  `Sync`ed, closed, chmod-ed to `0o600`, then `os.Link`ed into place — atomic for readers and
  `EEXIST` rather than a clobber if another publisher won the race, in which case the existing
  bytes are validated. Where hard links are unavailable (some network/virtualised filesystems —
  this repository sits on one) it falls back to `OpenFile(O_CREATE|O_EXCL)` at the target, which
  still never clobbers; that fallback is not atomic for a reader, which is exactly why the name
  carries the full digest. The staged temporary is removed on every path (`defer`).

### `internal/protocolpacket/packet.go`

- `Parse`: the fence tracker now records **which** marker opened the fence, and only that marker
  closes it. Previously any ``` or `~~~` line toggled, so a `~~~` line inside a ``` block closed
  the fence, turned the following fenced lines into "headings", and could fold a real heading's
  text into the preceding block under the wrong locator. New helper `fenceMarker`. The live
  protocol contains no `~~~` (grep: no match), so the 69-block invariant and the live map are
  unaffected; this is a robustness fix for future protocol text, and it fails closed either way
  (a swallowed heading shows up as `stale-map:` and forces full-fallback).

### Tests added (by inspection — not executed here)

`internal/protocolpacket/packet_test.go` (+ helper `packetContext`, imports `runtime`, `sync`):

- `TestPublishedBodyIsFullDigestAddressedImmutableAndIdempotent` — the name ends in the full
  64-hex body digest; the file holds the attested body; republishing identical content returns
  the same path with no error and leaves exactly one file (no staged leftovers); after the file
  is tampered with on disk, publication **fails** with "different body" and the tampered bytes
  are still there (not overwritten).
- `TestWriteBodyRefusesUnattestedAndEmptyBodies` — a body mutated after `Build` is refused as
  "unattested"; an empty body is refused; neither creates `.parley-runtime`.
- `TestWriteBodyRefusesSymlinkedRuntimePathsAndTargets` — symlinked `.parley-runtime`
  (and nothing written into the link target), symlinked `protocol-packets`, a regular file where
  the runtime directory belongs, and a symlinked target path (the victim file outside the
  workspace is byte-identical afterwards). Skipped on Windows.
- `TestPublicationDirectoryAndBodyArePrivate` — directories `0o700`, body `0o600`; a
  pre-existing `0o777` publication directory is tightened before use. Skipped on Windows.
- `TestConcurrentPublicationPublishesOneCompleteBody` — 8 goroutines released together publish
  the same request: all succeed, all report the same path, the directory holds exactly one file,
  and its bytes are the complete body. Then two goroutines publish two different bodies
  concurrently: distinct paths, and every file in the directory hashes to the digest its own
  name carries. (`Source`/`Map` are built once on the test goroutine; `Build` is pure, so this
  is also a `-race` exercise of concurrent `Build`.)
- `TestMismatchedFenceMarkerDoesNotHideHeadings` — a `~~~` line inside a ``` block does not end
  the fence, and the heading after the fence is still a block.

`internal/app/protocol_packet_test.go`:

- `TestProtocolPacketRepublicationIsIdempotentAndRefusesATamperedBody` — end-to-end through the
  CLI: the published name carries the full digest; a second identical invocation returns the
  same `body_path` and leaves one file; after the body is tampered with, the command exits **1**
  with "refusing to overwrite" and the tampered file is unchanged. This is the fail-closed
  behaviour a launch path depends on.

### What did NOT change (deliberately)

Default remains **full context + shadow packet** (`Optimize` false); the omission index remains
complete and disjoint; never-cut and live-authority rules, `ErrAuthority` fail-closed behaviour
and secret refusal are untouched; no optimization is shipped or enabled, and no experiment was
run, so nothing here rests on the unrun packet trial. No dependency was added (only stdlib
`runtime` in `source.go`, `runtime`/`sync` in the test). No unrelated refactoring. The
`ResolveSource` / `Render` / `BuildProtocolContext` / `Attestation` signatures are unchanged, so
the compact API published in my earlier handoff still holds; only the *file name* of a published
body changed (short hash → full digest), and no code outside my slice reads that name.

## Limitations (honest)

1. **Nothing was executed.** No `go build`, `go test`, `go vet` or `gofmt` ran in this session.
   The code and tests are written by inspection and may contain compile or assertion errors that
   only execution will reveal. Treat every "does X" above as a claim to be verified.
2. The hard-link path and the `O_EXCL` fallback are both implemented, but I could not observe
   which one this filesystem takes. Both are covered by the concurrency test; a verifier can tell
   them apart only by instrumenting, which I did not do.
3. `verifyPublished` compares full bytes, so it is TOCTOU-safe only in the sense that it never
   *writes*: a body replaced after validation and before a reader opens it is still possible on a
   directory that other users can write — which is why the directory permission gate exists. A
   reader that must be certain should re-hash the file and compare it with the name.
4. Bodies accumulate in `.parley-runtime/protocol-packets/` forever; there is no pruning or size
   cap. Immutability makes pruning a policy decision (whose owner is the runtime wiring), not a
   correctness one.
5. `Render` still labels *any* non-`ErrNoMap` map error `malformed-map:` — including an
   unreadable map. The fallback is visible and carries the OS error text, so it is honest, but
   "unreadable" and "malformed" are not distinguished. I left this rather than churn the slice.
6. `ResolveSource` reads `parley-deck/COOPERATION.md` with `os.ReadFile`, which follows a
   symlink. The new symlink rejection covers *publication*, not *authority resolution*. Whether
   a symlinked deck protocol file should be refused is a scope question for the owner group,
   not something I decided unilaterally.
7. Consumer-role authority is still only exercised with temp-dir fixtures (unchanged from my
   first handoff), and `Check` still proves structure, not semantics.

## Commands the facilitator should run

```sh
gofmt -l internal/protocolpacket internal/app
go vet ./internal/protocolpacket/ ./internal/app/
go test ./internal/protocolpacket/                       # all tests, incl. the 6 new ones
go test -race -count=4 ./internal/protocolpacket/ -run 'Concurrent|Published|WriteBody|Publication'
go test ./internal/app/ -run 'ProtocolPacket'
go test ./internal/protocol/                             # drift guard, unchanged by this session
go test ./...                                            # whole-tree regression
```

Worth an explicit adversarial check by the reviewer, since the tests are mine: publish a body,
`chmod 0777` the runtime directory, re-publish, and confirm the mode is tightened; replace a
published body with different bytes and confirm the CLI exits 1 and does not overwrite; point
`.parley-runtime` at `/tmp/somewhere` with a symlink and confirm nothing lands there.

## Unresolved integration dependencies

1. **The renderer is still not on any launch path.** Verified by inspection: no package outside
   `internal/protocolpacket/**` and `internal/app/protocol_packet*.go` imports
   `parley-deck-cli/internal/protocolpacket` (grep over `*.go` for `protocolpacket`). The prompt
   builders in `internal/runner/{runner.go,phase58.go,consult.go,steer.go,handoff.go}` are
   unchanged in that respect, so AC-P1's "across CLI, runner/handoff" half is not met by my slice
   alone. This is Codex-owned wiring, listed in my first handoff and still open.
2. **Attestation into telemetry is not wired.** `internal/telemetry/record.go` now defines
   `Context{context_mode, source_sha256, packet_sha256, fallback_reason}` (Codex's slice), but
   nothing populates it from `protocolpacket.Attestation`, because of (1).
3. **Resolved since my first handoff** (Codex, verified by reading the files): `.gitignore:23`
   contains `/.parley-runtime/`, and `internal/app/protocol.go:41` dispatches
   `case "packet": return runProtocolPacket(rest, stdout, stderr)`, so the command is reachable
   from the binary. I re-verified both this session.
4. Skill-source changes remain source-only in the sibling worktree, unchanged this session.

## Scoped findings in other owners' code (advisory, from files I had to read)

- `internal/app/protocol.go:22` (Codex-serialized) — the group usage line for `protocol packet`
  omits `--transport`, `--idea`, `--flag` and `--optimize`, which the subcommand's own
  `protocolPacketUsage` documents. NIT-level documentation drift; I did not edit that file.
- `internal/runner/handoff.go:65,68` (Codex-owned) — the handoff prompt and instructions are
  written `0o644` with `os.WriteFile`. Same shape of exposure as the defect fixed here, in a
  directory that holds prompt text; whether it matters depends on where `agentDir` lives, which I
  did not audit. Reported for the owner to judge, not changed by me.
- I did not review Hermes's or Kimi's slices this session.
