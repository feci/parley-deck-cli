---
from: claude-1
to: zcode-1
idea: windows-portability
phase: implementation-consult
blocking: no
date: 2026-09-28
---

## Status of this note

**ADVISORY ONLY.** This is a focused implementation consult, not a Phase-6 review and
not a signoff. Nothing here is a verdict. Every claim below must be mirrored into, and
adjudicated by, the canonical implementation record (`IMPLEMENTATION.md`) and the later
formal review before it carries any weight. I hold no implementation claim, I edited no
product code, tests, implementation docs or peer artifacts, and this file is my only
repo write.

Protocol attestation (as supplied to this invocation):
`{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}`

**Read at FIXED COMMIT `9b880d3`** via `git show 9b880d3:<path>` throughout. No
working-tree file was used for any verdict. Go stdlib sources were read from the local
toolchain `go1.27.1 darwin/arm64` (`/opt/homebrew/Cellar/go/1.27.1/libexec`).

**What I could NOT do:** I have no Windows machine and ran no Windows execution. I did
not run the hosted matrix. Nothing below is hosted-verified by me. The mechanism claim
is source-and-doc derived plus failure-signature discrimination; it needs one hosted
Windows confirmation (a cheap one is proposed in §6).

---

## 1. Bottom line

**The implementer's row-30 attribution is, on the primary evidence, wrong in mechanism** —
and the wrong mechanism points at the wrong remediation.

Row 30 (`IMPLEMENTATION.md` line 1008 at `9b880d3`) says:

> on Windows the Lstat→open `os.SameFile` chain does NOT detect a rename+recreate
> replacement (NTFS file-ID semantics; the recreated file presents the same file identity)

Both halves of the parenthetical look false:

* NTFS **does** give a recreated file a distinct 64-bit file index (MFT record number plus
  sequence number), and Go **does** read it — eagerly — whenever the `FileInfo` came from a
  handle. This is not an NTFS capability limit.
* The recreated file does **not** "present the same identity". Go never read the *old*
  file's identity at all.

**The actual mechanism is a Go API shape, not a filesystem property:** on Windows an
`os.Stat` / `os.Lstat` `FileInfo` carries **no captured identity**. It stores the *path*
and defers the identity lookup to the first `os.SameFile` call, which re-opens **by path
at comparison time**. So the "before" operand silently re-resolves to whatever now lives
at that path — which is exactly the thing it is supposed to be compared against.

Consequence: the idiom

```go
before := os.Lstat(p);  f := open(p);  opened := f.Stat()
if !os.SameFile(before, opened) { refuse }   // <- no-op on Windows
```

is **structurally incapable of detecting replacement-at-the-same-path on Windows.** It is
not flaky, not NTFS-dependent, not racy — it is a guaranteed no-op for that attack.

**This is therefore not only a fixture problem.** For the *specific failing test* the
fixture is unrepresentative and should be corrected (§4). But the same idiom appears in
**product** code with a lazy `before` operand at several sites (§3), and those are real,
silent, Windows-only guard no-ops that no current test covers. Fixing the fixture alone
would turn the test green while leaving the product holes in place — which is the worst
available outcome, because it would retire the only signal we have.

---

## 2. Primary evidence

### 2.1 Go stdlib doc (version-independent)

`go doc os.SameFile`:

> SameFile reports whether fi1 and fi2 describe the same file. For example, on Unix this
> means that the device and inode fields of the two underlying structures are identical;
> **on other systems the decision may be based on the path names.**

The documented contract already says the Windows comparison may be *by path*. It does not
promise identity semantics off Unix.

### 2.2 The lazy path — `os.Stat` / `os.Lstat` (GOROOT `src/os/`)

`stat_windows.go:25` `stat()` (implements both `Stat` and `Lstat`), non-reparse-point fast path:

```go
fs := newFileStatFromWin32FileAttributeData(&fa)   // sets times/size/attrs only
if err := fs.saveInfoFromPath(name); err != nil {  // <- stores the PATH
```

`types_windows.go:340` `saveInfoFromPath` sets `fs.path = <abs path>` and nothing else —
`vol`, `idxhi`, `idxlo` are left zero. `GetFileAttributesEx` returns no file index, so
there is nothing to capture.

`types_windows.go:287` `loadFileId()` — the deferred lookup:

```go
if fs.path == "" { return nil }                    // already resolved / eager
...
h, err := syscall.CreateFile(pathp, 0, 0, nil, syscall.OPEN_EXISTING, attrs, 0)
...
err = syscall.GetFileInformationByHandle(h, &i)
fs.path = ""                                       // memoized from here on
fs.vol, fs.idxhi, fs.idxlo = i.VolumeSerialNumber, i.FileIndexHigh, i.FileIndexLow
```

`types_windows.go:353` `sameFile`:

```go
func sameFile(fs1, fs2 *fileStat) bool {
    if e := fs1.loadFileId(); e != nil { return false }
    if e := fs2.loadFileId(); e != nil { return false }
    return fs1.vol == fs2.vol && fs1.idxhi == fs2.idxhi && fs1.idxlo == fs2.idxlo
}
```

`CreateFile(fs.path, ...)` runs **inside `SameFile`**, i.e. after the adversary has acted.

### 2.3 The eager path — handle-derived `FileInfo`

* `stat_windows.go:21` — `(*File).Stat` → `statHandle(file.name, fd)`
* `root_windows.go:227`, `:421` — `os.Root` stat/lstat → `lstatat` → `statHandle(name, fd)`
* `statHandle` → `newFileStatFromGetFileInformationByHandle` (`types_windows.go:46`), which
  sets `vol/idxhi/idxlo` **at stat time** and deliberately leaves `path` empty. Its own comment:

  > `fileStat.path` is used by `os.SameFile` to decide if it needs to fetch vol, idxhi and
  > idxlo. But these are already set, so set `fileStat.path` to `""` to prevent
  > `os.SameFile` doing it again.

So: **`File.Stat`, `Root.Stat`, `Root.Lstat` = eager, real NTFS identity.
`os.Stat`, `os.Lstat` = lazy, deferred path re-lookup.** That asymmetry is the whole bug.

### 2.4 Applying it to the two failing subtests

`snapshot_read_test.go:245` `TestSnapshotRevalidationRejectsDifferentMaterial`:

* `:264` `prior := snapshotStat(t, name)` → `snapshotStat` at `:29` is **`os.Lstat`** → LAZY.
* `root` case `:272`: rename `dir/sub` → `dir/old-sub`, then recreate `dir/sub/source` with
  the *same* bytes. `inode` case `:277`: rename `source` → `old-source`, recreate `source`
  with the same bytes.
* `:303` `verifySnapshotRegular(ctx, root, rel, prior, hash[:])`.
* Inside (`snapshot.go:360`): `entry, _ := fresh.Lstat(name)` → `Root.Lstat` → **EAGER**,
  new file's index. Then `os.SameFile(copiedInfo /* = prior, LAZY */, entry)` →
  `prior.loadFileId()` opens the path **now** → gets the **new** file's index → equal →
  `SameFile` returns **true** → guard does not fire.
* Bytes are identical by construction, size and mode match, ctx is live → `err == nil` →
  `t.Fatal("different material or canceled verification was accepted")`. Exactly the
  reported message.

**Discriminating check:** the other six subtests (`bytes`, `symlink`, `size`, `mode`,
`oversize`, `cancel`) are all caught by non-identity layers (hash, `IsRegular`, size, mode,
`ctx.Err`) and are reported green. The *only* two red cases are precisely the two whose
detection depends on identity differing across a path-preserving replacement. That is what
the lazy-path mechanism predicts, and it is not what "NTFS reuses identities" predicts —
NTFS identity reuse would be nondeterministic across cycles, whereas row 30 itself records
these failures as **identical in 36432547744 / 36434775624 / 36439175112**. Determinism is
evidence *for* a structural API artifact and *against* a filesystem-timing property.

### 2.5 Two secondary facts worth recording

* **Fail-closed on disappearance.** If the path is gone at comparison time, `CreateFile`
  fails, `loadFileId` errors, `sameFile` returns **false** → the guard fires → refusal. So
  *delete-without-recreate* is still caught on Windows. Only *replace-at-the-same-path*
  slips. That bounds the exposure exactly, and is a second reason the failure set is just
  `{root, inode}`.
* **Memoization.** `loadFileId` clears `fs.path` and caches the result, so a lazy operand's
  identity is pinned at its **first** `SameFile` call, not at `Lstat`. `verifySnapshotRegular`
  compares `copiedInfo` three times (`snapshot.go:360`, `:368`, `:397`); all three inherit
  the single post-race resolution. Adding more comparisons cannot recover the lost pin.

---

## 3. The part that matters more than the test: product exposure

I grepped every `os.SameFile` site at `9b880d3` and classified the **pinned ("before")
operand**. The rule is: *lazy pinned operand ⇒ the guard is a Windows no-op against
replace-at-same-path.*

**Lazy pinned operand — product code, believed broken on Windows:**

| Site (at `9b880d3`) | Pin | Compare |
|---|---|---|
| `internal/trajectory/state.go:225` / `:238` | `os.Lstat(path)` | `os.SameFile(info, opened)` |
| `internal/budget/step_history.go:178` / `:191` | `os.Lstat(path)` | `os.SameFile(before, opened)` |
| `internal/runner/action_identity.go:23` / `:40` | `os.Lstat(path)` | `os.SameFile(st, opened)` |
| `internal/budget/lock.go:306` / `:308` | `os.Lstat(path)` | `os.SameFile(info, opened)` |
| `internal/budget/lock.go:387` / `:400` | `os.Lstat(path)` | `os.SameFile(info, opened)` |
| `internal/budget/lock.go:478` / `:491` | `os.Lstat(path)` | `os.SameFile(info, opened)` |
| `internal/driver/cursor.go:120` / `:133` | `os.Lstat(path)` | `os.SameFile(info, opened)` |

Not yet classified by me (same shape, worth the same audit — I ran out of budget, do not
treat these as cleared): `internal/app/evidence_table.go:40`, `internal/app/evidence_verify.go:92`,
`internal/evidence/report.go:120`, `internal/protocolpacket/source.go:385`,
`internal/budget/ledger.go:382`, `internal/runner/cycle_budget.go:106`,
`internal/trajectory/reconcile.go:93`, `internal/trajectory/parent_recovery.go:158/172/181`,
`internal/budget/cycle_history.go:219`, `internal/tui/live.go:2468`.

**Eager pinned operand — believed correct on Windows:**

* `internal/trajectory/verification.go:222` (`dir.Lstat` on an `*os.Root` → eager) / `:235`.
* `internal/trajectory/snapshot.go` **production** path: `captureSnapshotMember` pins via
  `entryRoot.Lstat` (`:284`, eager); `captureSnapshotRegular` pins via `src.Stat()`
  (`:309`, eager) and passes that as `copiedInfo` at `:334`. **Production never feeds
  `verifySnapshotRegular` a lazy operand — only the test does.**
* `verifySnapshotRegular`'s own root pin `origin := root.Stat(".")` (`:349`, eager).

**Useful asymmetry:** where the lazy operand is the *after* side rather than the pin, the
guard still works — a lazy "current" correctly resolves the path *now* and is compared
against an eager pinned handle id. E.g. `action_identity.go:48-49` and
`step_history.go:199-200` (`current, _ := os.Lstat(path)` vs. eager `opened`) are fine in
that direction. So the invariant is specifically about the **pin**, not about `os.Lstat`
per se. That keeps the fix narrow.

---

## 4. Guard or fixture? — both, and they must land together

**The fixture is wrong** for `TestSnapshotRevalidationRejectsDifferentMaterial`: it hands
`verifySnapshotRegular` a weaker operand (`os.Lstat`) than any production caller ever
supplies (`src.Stat()`, `entryRoot.Lstat`). The test is therefore not testing the product's
actual pin. Correcting it is legitimate and is **not** a weakening.

**The product guards are also wrong**, at the seven-plus sites in §3, where the lazy pin is
what production really uses. Those are the genuine security finding.

**Do not fix the fixture alone.** Correcting `snapshotStat` makes the red go green while
every site in §3 stays silently broken, and retires the only symptom. If only one of the two
can land in this stage, land the **product** fix first and leave the test red with a ledger
row, rather than the reverse.

**Do not exclude or relax the test on Windows.** That path is closed by the owner's
constraint ("do not disable privacy checks or turn product failure into skipped tests"),
and row 30's current framing is the main thing pointing at it. The "the verify already
re-reads and hashes" note in row 30's disposition column does not rescue the case: in both
failing subtests the bytes are byte-identical by construction, so the hash layer passes.
`SameFile` is the **only** layer that can catch "same bytes, different file" — which is
verbatim the invariant `snapshot.go:390-391` claims to defend ("a replaced root or directory
entry must not be hidden by the original still-readable inode"). The guard is load-bearing,
not redundant.

---

## 5. Proposed fix

### 5.1 The invariant to adopt

> **The pinned ("before") operand of every `os.SameFile` guard must be handle-derived or
> `os.Root`-derived (`File.Stat`, `Root.Stat`, `Root.Lstat`) — never `os.Stat` / `os.Lstat`.**

This is greppable, testable, platform-neutral, and — importantly — it is a **strengthening**
of `os.Root` containment, not a weakening of it. It needs **no owner deviation** and no
build tags: it is the same source on every platform, it is a no-op semantically on POSIX
(where `os.Lstat` is already eager), and it makes Windows behave the way the signed
invariant already assumes.

### 5.2 Product changes

For the §3 sites, replace the plain pin with a rooted pin. The shape at
`internal/trajectory/state.go:225` becomes, sketching only:

```go
root, err := os.OpenRoot(filepath.Dir(path))   // or the caller's existing Root
...
info, err := root.Lstat(filepath.Base(path))   // EAGER identity + containment
...
f, err := root.Open(filepath.Base(path))       // same root handle, share-delete
opened, err := f.Stat()
if err != nil || !os.SameFile(info, opened) { ... }
```

Two things fall out for free: the pin becomes a real identity, and the open and the pin now
share **one** directory handle instead of resolving the parent path twice — which removes a
second, independent race. Where a caller already holds an `*os.Root` (the `trajectory`
package mostly does), this is a one-line substitution.

Where a plain path genuinely cannot be rooted, the minimum acceptable substitute is to pin
from the handle you are about to use (`f, _ := os.Open(p); pin, _ := f.Stat()`) and derive
every later comparison from `pin` — but note this only pins *after* the open, so it detects
post-open replacement, not pre-open. Prefer the rooted form.

A small named helper in `internal/fsutil` (e.g. `PinIdentity` / `OpenPinned`) would make the
invariant self-documenting and give reviewers one place to audit, matching how `DirHandle`
already makes the durability contract "the audit of record rather than a grep"
(`syncdir.go:43-46`). Your call — the helper is a convenience, the invariant is the fix.

### 5.3 Test changes

`snapshot_read_test.go`: derive `prior` from the root instead of `os.Lstat`, so the fixture
matches production. Note the ordering — `prior` is currently taken at `:264`, *before*
`root := snapshotOpenRoot(t, dir)` at `:265`; the root must move up first. That reorder is
harmless (nothing mutates between the two lines).

Both `TestSnapshotRevalidationRejectsDifferentMaterial` (`:264`) and
`TestSnapshotRevalidationRefusesChangeDuringRead` (`:463`) use the same helper, so check the
second one too — a lazy pin there is a second latent hole in the during-read case.

---

## 6. Cheapest hosted confirmation before you commit to this

I cannot execute on Windows, so the mechanism above is source-derived, not observed. One
tiny Windows-tagged unit test settles it in a single hosted cycle, with no product change:

```
write f; lazy  := os.Lstat(f); h,_ := os.Open(f); eager,_ := h.Stat(); h.Close()
rename f -> g; recreate f with identical bytes
after,_ := os.Lstat(f)
os.SameFile(lazy,  after)   // predict TRUE   <- the bug
os.SameFile(eager, after)   // predict FALSE  <- identity works when pinned eagerly
```

If both predictions hold, row 30's attribution is disproved, the fix in §5 is justified, and
the pair becomes the permanent regression pin for the invariant in §5.1 — including the
negative half, which is what stops the lazy idiom being reintroduced later. If `os.SameFile(eager, after)`
comes back TRUE, then NTFS identity reuse *is* in play after all, my §1 is wrong, and the
problem is genuinely harder — so this check is worth running before any code moves.

**Version caveat, stated plainly:** I read Go **1.27.1** sources. `go.mod` pins `go 1.26`
and `.github/workflows/tests.yml:42` uses `go-version-file: go.mod`, so hosted Windows runs
on 1.26.x. I did **not** verify the 1.26 source bytes — only 1.27.1 was installed locally.
The `os.SameFile` doc quoted in §2.1 is version-independent and the lazy `loadFileId` design
is longstanding, but please re-confirm at the CI-pinned toolchain. The §6 probe does this
implicitly.

---

## 7. `fsutil.OpenSharedDelete` at `9b880d3` — secondary look

`internal/fsutil/syncdir.go:37`. Callers: `internal/budget/step_history.go:185`,
`internal/trajectory/state.go:232`.

**Share-flag claim: CORRECT.** The doc says routing through an `os.Root` "on Windows passes
`FILE_SHARE_DELETE`". Verified: `os/root_windows.go:176` and
`internal/syscall/windows/at_windows.go:151` use
`FILE_SHARE_READ|FILE_SHARE_WRITE|FILE_SHARE_DELETE`, whereas plain `syscall.Open` uses only
`sharemode := uint32(FILE_SHARE_READ | FILE_SHARE_WRITE)` (`syscall/syscall_windows.go:395`).
So the rooted route really does grant delete-share and the plain route really does not. The
AC-LOCK-1 rationale holds.

**Containment claim: PARTIAL, and I'd call the comment overstated.**

```go
root, err := os.OpenRoot(filepath.Dir(name))
...
f, err := root.Open(filepath.Base(name))
```

Only the **final component** is contained. `filepath.Dir(name)` is resolved with full ambient
authority — every symlink, junction, and `..` in the parent chain is followed. Under a race:

* *After* `OpenRoot` the handle is pinned to a directory object, so a later swap of the parent
  cannot redirect the `root.Open`. Good.
* *Before* `OpenRoot`, an adversary controlling any parent component can point the whole thing
  at a different directory. The Root then faithfully contains reads **inside the wrong
  directory**, and nothing notices, because there is **no origin identity pin** — no
  `root.Stat(".")` + `os.SameFile(expectedOrigin, actual)`. Contrast `verifySnapshotRegular`,
  which does exactly that at `snapshot.go:349` and re-pins at `:393`; and
  `captureSnapshotMember`, which re-pins at `snapshot.go:281`. `OpenSharedDelete` has the
  rooted open without the pin that makes rooting meaningful.

So the honest statement is: *"delete-share correct; final-component containment; parent
resolution and directory identity are the caller's responsibility."* Since neither caller
currently supplies an origin pin, that responsibility is at present unmet.

**And both callers then apply the §3 broken guard.** `state.go:225/238` and
`step_history.go:178/191` pin with `os.Lstat` and compare with `os.SameFile`. On Windows those
two sites therefore have **neither** a contained parent **nor** a working identity check —
they are the highest-value places to apply §5.2, and the fix is natural there because they
already need a Root for the share-delete open. Passing that same Root to the pin fixes both
defects with one handle.

**Lower-confidence, flagged not asserted:** the doc requires callers to "have verified the
path is a regular file"; both do so via `os.Lstat` + `IsRegular()` *before* the open, which is
itself TOCTOU. `step_history.go:200` re-checks `current.Mode().IsRegular()` after the read;
`state.go` does not re-check `IsRegular` post-open independently of `SameFile`. With the
`SameFile` pin repaired this is probably moot, but it is worth one look while you are in
there. I did not trace it to a concrete exploit and am not claiming one.

---

## 8. What I am and am not claiming

**Claimed, with primary evidence cited above:** the Go Windows lazy/eager `FileInfo` split
(§2.2/§2.3, GOROOT source at 1.27.1); the `os.SameFile` documented path-based contract
(§2.1); that the failing test's pin is `os.Lstat` and production's is handle/Root-derived
(§2.4, §3, source at `9b880d3`); that `os.Root` opens pass `FILE_SHARE_DELETE` and plain
opens do not (§7, GOROOT source); that `OpenSharedDelete` resolves its parent unrooted and
carries no origin pin (§7, source at `9b880d3`).

**Inference, not observation:** that the lazy/eager split is *the* cause of the two reported
red subtests. It is strongly supported by the failure-set discrimination and by the reported
cross-cycle determinism, but I did not execute on Windows. §6 converts it to observation.

**Not claimed:** any hosted verification, any Windows execution, any assertion about Go
1.26.x source bytes, any exploit for the §7 `IsRegular` note, and any clearance of the
unclassified sites listed in §3. I have not verified the row-30 cycles' logs myself; I took
the run IDs and failure names from `IMPLEMENTATION.md` at `9b880d3` as reported by you.

**Requested of the canonical record:** if you accept §1, row 30's mechanism text needs
correcting — not because attribution tidiness matters, but because "NTFS file-ID semantics"
is an argument for accepting or excluding the failure, and "Go defers the identity lookup to
comparison time" is an argument for fixing seven product guards. The remediation follows the
mechanism, so the mechanism has to be right in the ledger.
