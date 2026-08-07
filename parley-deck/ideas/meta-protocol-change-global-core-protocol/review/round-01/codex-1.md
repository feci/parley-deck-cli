---
agent: codex-1
idea: meta-protocol-change-global-core-protocol
review-round: 1
date: 2026-08-07
reviewed-commit: 4396529
---

## Summary

This implementation is not ready to merge. G1, G2, and G7b fail at the real CLI boundary. I
reproduced three forms of silent renderer loss/misreporting, used a non-interactive pseudo-terminal
to publish while the child had `stdin </dev/null`, published into an already-existing release
directory, and made `Publish` follow a release-directory symlink outside the core. The committed
tests all pass in isolation, but they call `runProtocol` directly; deleting the production
`app.Run` dispatch still leaves all nine green.

The exact requested command did not reach compilation in this sandbox:

```text
$ go build ./... && go test ./...
cmd/parley/main.go:6:2: open /Users/tomasfecko/Library/Caches/go-build/df/dfb9c55bf55f52cdf3fae47bd3e9f4e30c5848e0fe8a229c659274cdd3153527-a: operation not permitted
pattern ./...: open /Users/tomasfecko/Library/Caches/go-build/df/dfb9c55bf55f52cdf3fae47bd3e9f4e30c5848e0fe8a229c659274cdd3153527-a: operation not permitted
```

Because `&&` short-circuited, `go test` did not run in that invocation. With only `GOCACHE`
redirected to an allowed temporary directory, `go build ./...` succeeded and the clean
`4396529` snapshot test run reached one unrelated existing failure:

```text
--- FAIL: TestDurableKillEndToEndRealProcess (0.02s)
    durablekill_test.go:116: a live attributed process should be killed, got {AgentID:sleeper Killed:false Cleared:false Failed:true SegmentID:segment-0001 Message:process verification failed (no recorded boot id); not killed}
FAIL
FAIL parley-deck-cli/internal/runner 10.459s
```

The same targeted runner test fails identically at base commit `fac2421`, so I do not attribute it
to this implementation. The nine committed protocol tests pass on a clean `4396529` snapshot:

```text
$ GOCACHE=/private/tmp/parley-review-gocache go test ./internal/app -run '^(TestProtocol|TestCore)' -count=1 -v
...
PASS
ok   parley-deck-cli/internal/app 0.398s
```

## Refutation attempts

- G1: rerendered normal, changed-header, two-official-table, extra-table, and CRLF decks; compared
  preview/apply reports and second-run bytes.
- G2: invoked the compiled publisher headlessly, through an allocated PTY, with child stdin set to
  `/dev/null`, against an existing directory, through a symlink, and with traversal strings.
- D8/G7/G8: traced every `protocolcore` caller and searched the run/continue/resume/steer/inspect
  paths for pin/snapshot resolution; only render/check use the minimal lock.
- G7b: removed production dispatch in a temporary clean snapshot and reran all nine new tests; they
  stayed green. I also compared current/base Windows cross-builds and current/base runner failure.
- Purity: traced every import and input used by `Render`, including how production `Load` derives
  the release hash used in the synced stamp.

## Findings

### [CRITICAL] The TTY gate is agent-bypassable and still permits `publish < /dev/null`

`protocolPublish` treats a successful terminal ioctl on either stdin **or stdout** as user
attendance (`internal/app/protocol.go:269-280`, `internal/app/protocol.go:301-316`). A TTY is an I/O
property, not user authorization. An autonomous process can allocate a pseudo-terminal itself; and
because stdout alone is enough, the exact child command with stdin redirected from `/dev/null`
passes. I drove the compiled `parley` entry point in an isolated `PARLEY_HOME`:

```text
$ script -q /dev/null /bin/sh -c 'PARLEY_HOME=... parley protocol publish ... < /dev/null' < /dev/null
Published core stdin-devnull-stdout-pty (6045aa22d988) to .../protocol/core/stdin-devnull-stdout-pty
script_exit=0 exists=yes
```

This directly contradicts both the CLI refusal text, “Only the user may change the global
protocol; an agent run cannot” (`internal/app/protocol.go:277-280`), and the newly shipped §7 claim
that an agent cannot publish (`parley-deck/COOPERATION.md:759-762`). It also violates G2
(`FINAL.md:59-60`). There is no committed TTY/publish call-site test; the test file contains no
reference to `hasTTY` or `protocolPublish`.

Suggested fix: do not use possession of a TTY as an authorization boundary. Remove the core writer
from the agent-reachable binary, or require an out-of-band, trusted user-ratification mechanism
that the launched agent cannot mint or feed through a PTY. Until such a boundary exists, remove the
user-only/prevention claims and report detection-only. Add a compiled-binary negative test using a
pseudo-terminal and a child whose stdin is `/dev/null`.

### [CRITICAL] `Publish` can populate an existing release namespace and follows symlinks

The write-once check stats only `<release>/COOPERATION.md`; it does not reserve or validate the
release directory (`internal/protocolcore/core.go:93-112`). `MkdirAll` and `WriteFile` follow
symlinks, and the stat/write sequence is not exclusive. At the real CLI entry point, after defeating
the TTY gate:

```text
--- pre-existing release directory without COOPERATION.md ---
Published core preexisting-dir (6045aa22d988) to .../protocol/core/preexisting-dir
script_exit=0 exists=yes

--- symlinked release directory ---
Published core symlink-release (6045aa22d988) to .../protocol/core/symlink-release
script_exit=0 outside_exists=yes
```

Literal traversal strings were rejected as intended:

```text
protocol publish: protocolcore: unsafe version "../escape"
dotdot_exit=1
protocol publish: protocolcore: unsafe version "slash/name"
slash_exit=1
```

That string check does not make the resolved path safe. A symlink already inside the store escapes
the root, and an existing version directory without the expected file is modified. Two concurrent
publishers can also both pass `os.Stat` and then open the same file with truncation. Consequently the
“write-once” statements in `core.go:88-92`, §7, and `IMPLEMENTATION.md:21-24` are not true for the
implemented store.

Suggested fix: reject any existing version directory with `Lstat`; reject symlinks in every store
path component; verify the resolved target remains below the resolved store root; create the
version atomically with no-replace semantics; and create files with `O_CREATE|O_EXCL` plus
no-follow semantics. Publish a complete temporary sibling and atomically install it, so crashes and
races cannot leave a republishable partial release. Add real-entry tests for a normal existing
release, pre-existing empty directory, directory symlink, final-component symlink, and concurrent
publish.

### [CRITICAL] The heading heuristic silently erases tables and cannot satisfy G1's block report

The binding design requires permanent registry IDs and explicitly says addressing is “by ID, never
by heading text” (`consensus.md:88-93`); rank 1 repeats that the old heading anchor does not survive
(`consensus.md:222-228`). The shipped store contains no registry, while the renderer recognizes
tables by English header substrings and computes removals by literal Markdown headings
(`internal/protocolcore/render.go:64-80`, `render.go:118-177`, `render.go:180-213`). It never computes
or reports replaced blocks at all.

Changing only the deck's first `| Agent ID` header to `| Participant ID` caused preview to report
no removal, then apply deleted every roster row. The official host-handle table survived, showing
that the lost rows were not merely absent from the fixture:

```text
preserved from this deck: Workspace, Transport, Created, host-handle table
would regenerate .../deck-header/parley-deck/COOPERATION.md ...
Nothing was written. Re-run with --yes to apply.
preserved from this deck: Workspace, Transport, Created, host-handle table
Wrote .../deck-header/parley-deck/COOPERATION.md ...

--- baseline §2 markers ---
133:| Agent ID ... | Workspace dir ... | Role ... |
135:| `claude-1` ...
136:| `codex-1` ...
149:| Agent ID ... | Host handle |

--- changed-header §2 markers after apply ---
133:| Agent ID ... | Workspace dir ... | Role ... |
145:| Agent ID ... | Host handle |
147:| `claude-1` ...
148:| `codex-1` ...
```

An extra table inside the already-existing §2 heading was also dropped without being named:

```text
preserved from this deck: Workspace, Transport, Created, §2 roster table, host-handle table
Wrote .../deck-extra-table/parley-deck/COOPERATION.md ...
extra_table_survived=no
```

CRLF input exposes the opposite reporting error. With byte-identical LF core content and a CRLF
deck, `heading` leaves `\r` attached because it trims only spaces and tabs (`render.go:207-212`), so
preview falsely announced 65 unchanged headings as removals. Apply left 11 carriage returns in an
otherwise LF render. The second apply was byte-idempotent:

```text
the following section(s) exist in this deck but NOT in core pty-agent and will be REMOVED:
  - ## Quickstart — start here (developers & first-timers)
  - ## 0. Choose the transport
  ...
reported_removed_lines=65
carriage_returns_after_apply=      11
protocol render: ... already matches core pty-agent — nothing to do
```

Thus idempotence held for the normal, changed-header, two-official-table, extra-table, and CRLF
inputs I tried, and both official identity tables survive when their headers exactly match. G1
still fails because replacements are never reported, silent row/table loss is possible, and CRLF
reports removals that do not occur.

Suggested fix: implement the binding release registry and render by stable block/slot IDs. Parse
identity zones with exact, versioned schemas; on an unrecognized or duplicated identity table,
fail closed and name the ambiguity rather than treating it as a fresh slot. Normalize line endings
for structural comparison while choosing one deterministic output convention. The returned change
set must enumerate every replaced and removed registry block, not only headings absent from core.

### [MAJOR] The “real command entry point” tests bypass production dispatch

Production reaches the feature through `app.Run` and its `case "protocol"` dispatch
(`cmd/parley/main.go:9-10`, `internal/app/app.go:42-51`, `app.go:100-101`). Every command-oriented new
test calls the private `runProtocol` function directly, e.g. `protocol_test.go:87`,
`protocol_test.go:117`, `protocol_test.go:152`, and `protocol_test.go:268`. The remaining
write-once test calls `Store.Publish` directly (`protocol_test.go:246-260`).

I reverted only the production dispatch in a temporary clean snapshot and reran all nine committed
protocol tests:

```text
dispatch_still_present=no
ok   parley-deck-cli/internal/app 0.387s
```

Therefore `IMPLEMENTATION.md:64-67` is factually wrong when it says the nine tests drive the real
command entry points. The whole command can become unreachable while the advertised call-site gate
stays green. The TTY fix can likewise be reverted with no failing committed test.

Suggested fix: drive at least `app.Run([]string{"protocol", ...})` for every G7b surface, and use a
compiled-binary subprocess where process I/O/TTY behavior matters. Add a test that would fail if the
`app.Run` dispatch or top-level help entry disappeared.

### [MAJOR] Shipped text claims rank-2 continuation/pinning that does not exist

`FINAL.md:52-54` and `IMPLEMENTATION.md:15-17` correctly say per-idea pinning is rank 2 and is not
shipped. No continuation, resume, steer, inspect, run-manifest, or prompt path imports or calls
`protocolcore`; all production references are confined to `internal/app/protocol.go`. There is no
snapshot writer or `protocol-snapshots/` resolver.

Nevertheless, the new §7 text says an open idea “completes under the protocol version it was pinned
to” and the next idea picks up current (`parley-deck/COOPERATION.md:766-767`; the same text is in the
embedded default). That is the exact rank-2 guarantee G7b says must not be documented as landed.
The missing-release test covers only `protocol render` (`protocol_test.go:209-225`); there is no
adoption command and no continuation test.

`IMPLEMENTATION.md` itself does not overclaim that continuation is implemented; its deferral is
clear. The overclaim is in the protocol text shipped by this commit (and in the future-tense code
comment at `internal/protocolcore/core.go:28-32` if read as current behavior).

Suggested fix: remove the open-idea continuation sentence from all shipped protocol copies until
rank 2 lands with G5/G7/G8. For this slice, document only the behavior that exists: an explicitly
written `core-version:` lets `render`/`check` load that release, and a missing release makes those
two commands fail.

### [MAJOR] The fail-closed platform claim is a Windows build regression

`protocol.go` imports `golang.org/x/sys/unix` unconditionally and defines `hasTTY` there, while the
only `termiosGet` definitions are Linux and BSD-family files (`internal/app/protocol.go:12`,
`protocol.go:301-316`, `internal/app/termios_linux.go:1-7`, `internal/app/termios_unix.go:1-10`).
There is no other-platform implementation returning false. This contradicts
`IMPLEMENTATION.md:60-62`, which says a platform without the ioctl “refuses.” It does not run and
refuse; it no longer builds.

```text
# current 4396529
$ GOOS=windows GOARCH=amd64 go build ./internal/app
# parley-deck-cli/internal/app
internal/app/protocol.go:312:21: undefined: unix.IoctlGetTermios
internal/app/protocol.go:312:50: undefined: termiosGet
exit=1

# base fac2421
$ GOOS=windows GOARCH=amd64 go build ./internal/app
exit=0
```

Suggested fix: move the ioctl import and supported implementation behind build tags and add a
complementary `termios_other.go` whose `hasTTY` returns false. That fallback is fail-closed, not the
weak character-device fallback the implementation intended to remove. Add the formerly green
Windows cross-build to CI.

### [MAJOR] The claimed third protocol copy was not updated

The implementation scope requires D11 in all three copies (`FINAL.md:48-50`), and
`IMPLEMENTATION.md:42-46` claims that happened. This repository contains and updates only the live
deck and embedded default; the installed/package fallback remains old:

```text
$ rg --files -g 'COOPERATION.md' -g '!graphify-out/**'
parley-deck/COOPERATION.md
internal/protocol/defaults/COOPERATION.md

$ shasum -a 256 parley-deck/COOPERATION.md internal/protocol/defaults/COOPERATION.md /Users/tomasfecko/.codex/skills/parley-deck/references/COOPERATION.md
6045aa22d988971b219d28195fc0babe88f91621213f4553d1e518e50cd7dd52  parley-deck/COOPERATION.md
a0128d399f79f7408cf9d828f0d20667627191b90786f5336b0bd93f60b200d3  internal/protocol/defaults/COOPERATION.md
8cd0c64340eb990f045947c23312061c7265e2e24cffcdc74f9122955b62b23e  /Users/tomasfecko/.codex/skills/parley-deck/references/COOPERATION.md
```

`rg -n 'Blast radius — a CORE change'` returns no match in the installed fallback. The existing
drift test explicitly compares only the two repository copies (`internal/protocol/drift_test.go:10-20`),
so it cannot catch this omission.

Suggested fix: update and release the packaged skill fallback as the third deliverable, or amend the
ratified scope/implementation record if “three copies” meant something else. Add a release-time
guard covering that artifact; do not claim it was updated before it is present.

## Gate-by-gate assessment

- **G1 — FAIL.** Second renders were byte-identical/no-op for every input tried, including CRLF.
  Both expected identity tables survive with exact headers. However, replacements are not reported;
  a one-token header change and an extra table cause unreported content loss; CRLF causes 65 false
  removal reports and mixed line endings.
- **G2 — FAIL.** A non-interactive agent can mint a PTY and reach `publish`; `stdin </dev/null` still
  passes when stdout is a PTY. Existing directories and symlinked directories defeat write-once/path
  confinement. The only passing part is rejection of literal slash/`..` version strings and the
  normal same-file republish case.
- **G3 — NOT SHIPPED / scope contradiction.** There is no confinement probe or confinement-reporting
  surface in this slice. Consensus defers the sandbox, but `FINAL.md:73` explicitly defers only G5,
  G7, and G8 even though rank 4 is unshipped. If read literally, G3 is unmet rather than N/A.
- **G4 — NOT SHIPPED / scope contradiction.** There is no overlay implementation. As with G3,
  `FINAL.md:73` does not list G4 among deferred gates despite rank 3 being excluded.
- **G5 — DEFERRED by `FINAL.md:73`.** No per-idea protocol pin or later-phase resolver exists.
- **G6 — PASS for the repository artifact.** `parley-deck/meta/protocol-changelog.md` has the required
  `Idea:`, `Drafted by:`, and `Summary:` fields. This does not cure the stale third copy.
- **G7 — DEFERRED by `FINAL.md:73`.** No production call-site pin/continuation implementation exists.
- **G7b — FAIL.** Tests bypass `app.Run`; publish/TTY has no call-site test; renderer change reports
  fail adversarially; rank-2 continuation and user-only publishing are documented; atomic/mode
  preservation, other-platform refusal, and the third-copy claim have no end-to-end proof.
- **G8 — DEFERRED by `FINAL.md:73`.** The minimal lock stores only `core-version:` and verifies no
  core/effective hash. Missing-release render blocking exists, but adoption and snapshot-based
  continuation do not.

For D8 specifically: `render` and `check` fail on a missing pinned release through
`resolveRelease` (`internal/app/protocol.go:103-125`), and the render test verifies the deck is not
modified. Adoption is not implemented. Existing `continue`/`resume` behavior is not the ratified
continuation guarantee because it has no protocol snapshot input; it simply has no relationship to
the new core yet.

The renderer itself is pure as implemented. `internal/protocolcore/render.go` imports only `fmt`
and `strings` (`render.go:3-6`), reads only its two arguments, and builds the stamp from
`rel.Version` and `rel.SHA256` (`render.go:45-93`). Production `Load` derives `SHA256` from the
release bytes (`internal/protocolcore/core.go:55-64`). No time, environment, filesystem, or global
state is read by `Render`.

## Test-quality assessment

The nine tests are useful narrow regression tests, but they are not the G7b end-to-end suite
claimed in `IMPLEMENTATION.md`:

1. `TestProtocolRenderIsIdempotent` tests real byte behavior at `runProtocol` and would catch a
   changing stamp or non-idempotent internal render. It would not catch missing production dispatch.
2. `TestProtocolRenderPreservesIdentityAndReplacesCoreText` tests behavior, but only with a header
   byte-compatible with the core. Its global `strings.Contains` assertions can pass if a value moves
   to the wrong table or appears elsewhere. It does not assert the host-handle identity slot.
3. `TestProtocolRenderReportsRemovedSections` tests output in preview and apply, but its preview
   non-write assertion is tautological: it checks that stdout says `Nothing was written`
   (`protocol_test.go:151-160`) without comparing file bytes. Apply checks that stdout names the
   heading but never asserts the section was actually removed (`protocol_test.go:162-169`).
4. `TestProtocolCheckReportsHandEditAndNeverWrites` is the strongest new test: it asserts exit code,
   message, and unchanged file state. It still bypasses `app.Run`.
5. `TestProtocolBlocksWhenPinnedReleaseIsMissing` tests render failure and unchanged deck state. It
   does not test adoption or continuation, so its D8 comment is broader than its evidence.
6. `TestProtocolRefusesToGuessAVersion` tests the render error path and explanation. It does not
   test any adoption flow because none exists.
7. `TestCoreReleasesAreWriteOnce` is a genuine unit test for the normal existing-file case. It is
   not the required G2 guard test “asserting both halves”: it never drives the CLI/TTY boundary,
   permissions, existing directories, symlinks, or races.
8. `TestProtocolStatusJSON` asserts JSON shape and values at `runProtocol`, not production dispatch.
9. `TestProtocolRenderOnAFreshDeck` tests fresh-deck behavior at `runProtocol`, not production
   dispatch.

Concrete revert check: removing `app.Run`'s `protocol` case leaves all nine tests passing. Reverting
the termios fix to the old character-device check also leaves them passing because no test calls
`hasTTY`. Thus the two fixes most dependent on the real call site have no regression protection.

Missing adversarial cases: mismatched table schema/header, both official tables, duplicate/extra
tables, CRLF, replacement reporting, preview/apply file-state symmetry, PTY allocation, stdout-TTY
with stdin `/dev/null`, existing empty release dir, symlink path components, concurrent publish,
file mode, atomic failure behavior, production dispatch/help, Windows cross-build, and all rank-2
pin/continuation paths when that rank lands.

## Open questions

1. `FINAL.md:73` defers only G5/G7/G8, but the excluded rank-3/rank-4 work is exactly what G4/G3
   govern. Are G3/G4 intended to gate this commit, or should FINAL be superseded/corrected before
   implementation review can close?
2. What trusted event constitutes “explicit user ratification” for publishing? A TTY cannot do so
   because an agent can allocate and feed one. The answer determines whether `publish` must leave
   the agent-facing binary entirely or use a separate OS/user authorization channel.
3. What is the third `COOPERATION.md` deliverable named in FINAL and IMPLEMENTATION? If it is the
   packaged skill fallback, its release/update must be part of the agreed fix rather than recorded
   as already complete.
