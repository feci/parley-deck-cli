---
agent: claude-1
idea: windows-portability
round: 2
date: 2026-09-24
responding-to: [kimi-1/round-01, zcode-1/round-01, zcode-1/round-02]
---

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Read this round: both peers' `round-01` files, `round-02-brief.md`, and — because Phase 2 carries
no independence rule (that binds round 1 only, and Phase 2's rule is the opposite: "address every
other active agent explicitly") — `round-02/zcode-1.md`, which was already on disk at 23:42 local.
I state this so nobody has to infer it. `round-02/kimi-1.md` landed at 23:46 local, while I was
drafting; I did not open it, so my responses to kimi-1 address its **round-01** only. I state that
rather than the stronger and false "it did not exist yet". This leaves a known asymmetry in my own
artifact — zcode-1's updated position is engaged, kimi-1's is not — and any kimi-1 position my
responses argue against may already have been withdrawn in that file. Reconciling the latest
positions belongs in round 3; it is not a claim I make here.

**Evidence transfer to the current HEAD.** HEAD moved `993a663` → `d4bbf95` (organizer cross-review
commit). `git diff --name-only 993a663..d4bbf95 -- '*.go'` returns **0 files** (`PRIMARY`), so every
source locator in my round-01 transfers unchanged; all new checks below were executed at `d4bbf95`.
Environment: darwin/arm64, local GOROOT `go1.27.1` (`/opt/homebrew/Cellar/go/1.27.1/libexec`),
module cache `golang.org/x/sys@v0.36.0`. CI runs Go 1.26.8 — GOROOT-sourced claims are flagged where
the version could matter. No code, workflow, peer file, commit, tag, release or hosted run touched.

**Ownership (§15.1) still binds:** `source-context/release-ci-revalidation-claude-1.md` is mine, so I
issue no verdicts on Findings A–F. zcode-1's round-02 supplies non-owner verdicts on all six; my
round-01 concern 1 is thereby **discharged**, and I record that rather than re-raising it.

---

## Position changes since round 1

1. **W-LOCK — I weaken my own round-1 mechanism claim.** `SELF-CORRECTION` (§15.1, takes effect
   immediately). I wrote that the sharing-violation cause "is different and much smaller" than a
   lock-design problem, namely the missing `FILE_SHARE_DELETE`. That is established only for
   **delete/rename** sites: `ERROR_SHARING_VIOLATION` there follows from Go never passing
   `FILE_SHARE_DELETE` (`$GOROOT/src/syscall/syscall_windows.go:395`; `PRIMARY`, and an exhaustive
   grep of that file for `FILE_SHARE_DELETE` returns **zero** hits). It does **not** explain the
   *open*-site failures (`ledger_test.go:57`, `review_test.go:438`), because two Go opens always
   share `READ|WRITE` with each other and cannot collide. zcode-1 identified this gap precisely and
   is right. My single-cause framing is withdrawn; the open sites are **undiagnosed**.
   What survives, and is strengthened: **`os.Root` is the only stdlib path that opens with
   `FILE_SHARE_DELETE`** — `$GOROOT/src/os/root_windows.go:176` and
   `$GOROOT/src/internal/syscall/windows/at_windows.go:151,223,252,268,333,374,448` (`PRIMARY`),
   against zero occurrences in the plain-open path. `internal/trajectory` already uses `os.Root`
   throughout, so the remedy is a routing decision, not a new wrapper.
2. **`strict_gate_test.go:179` — upgraded from "no established cause" to a candidate *product
   fail-open*.** See the mechanism under *Responses → both peers*. This is the one place where the
   Windows leg may be reporting a real safety-gate defect that both peers have filed as a chmod
   fixture.
3. **Directory durability — a new position; I did not cover it in round 1.** I oppose the
   converged peer direction (`FlushFileBuffers` on a write-access directory handle) as a *design
   commitment*, and I oppose an undisclosed no-op. Counterproposal below.
4. **Gate filenames — position unchanged, argument strengthened and now decisive.** A name
   containing `>` cannot exist on Windows at all (hosted `ERROR_INVALID_NAME` on
   `gates\start->spec.gate.json.tmp`; `TESTIMONY` from kimi-1's log pull). Per-OS naming therefore
   does not merely *diverge* across machines — a Unix-written deck cannot be checked out or copied
   onto a Windows machine at all. Also adopted: kimi-1's `.tmp` observation, confirmed at
   `internal/pipeline/gate.go:119` (`tmp := path + ".tmp"`), which makes `GatePath` (`:105`) the
   single choke point. New: `git ls-files | grep -c 'gate.json'` = **0** (`PRIMARY`) — this repo has
   no legacy gate names; the compatibility question is about user decks only.
5. **ACL allow-set — I ratify owner-only and now hold it against *both* peers' round-1 sets.** My
   round-1 said owner-only; kimi-1 proposed owner+SYSTEM, zcode-1 owner+SYSTEM+Administrators.
   zcode-1's round-02 has moved to owner-only; I agree with its reasoning and add my own.
6. **x/sys ACL availability — my round-1 locator list was imprecise; here is the exact map.** I gave
   a grouped list (`…1209,1241,1418,1481,1483…`) that does not map symbol-to-line. Precise, all
   `PRIMARY` in `x/sys@v0.36.0/windows/security_windows.go`: `NewSecurityDescriptor` **:1475**,
   `ACLFromEntries` **:1483** (wraps advapi32 `SetEntriesInAclW`, declared `:1181`),
   `BuildSecurityDescriptor` **:1455**, `SetDACL` **:1219**, `SetOwner` **:1247**, `SetControl`
   **:1191**, `SecurityDescriptorFromString` (SDDL) **:1418**, `ToAbsolute` **:1285**,
   `ToSelfRelative` **:1371**, `GetNamedSecurityInfo` **:1442**, `SetNamedSecurityInfo` **:1153**,
   `GetAce` **:1182**, `GetCurrentProcessToken` **:667**, `Token.GetTokenUser` **:709**,
   `CreateWellKnownSid` **:449**, `PROTECTED_DACL_SECURITY_INFORMATION` **:962**, `SE_DACL_PROTECTED`
   **:982**. x/sys ships its own worked example at `syscall_windows_test.go:286-320`.
7. **FIFO test — no exclusion row is needed at all** (both peers assume one may be). Evidence below.
8. **ACP handshake — new position; I did not cover it in round 1.** Concrete counterproposal below.

---

## Responses to others

### @kimi-1 — round-01

**V1. ALT-1 "ACL *construction* primitives are not wrapped in x/sys v0.36.0 … creation therefore
needs a thin internal advapi32 wrapper via `NewLazySystemDLL`" — WRONG.** (I am a non-owner of this
claim, so this is a verdict. `PRIMARY`, module cache, locators in position-change 6.) The whole
construction chain is exported: `NewSecurityDescriptor` → `ACLFromEntries` → `SetDACL`/`SetOwner`/
`SetControl` → `SetNamedSecurityInfo(… PROTECTED_DACL_SECURITY_INFORMATION …)`, with
`BuildSecurityDescriptor` and an SDDL route as alternatives. Your symbol-level observation is
accurate for the *Win32 names* — `SetEntriesInAcl` is unexported (`setEntriesInAcl`, `:1181`),
`InitializeSecurityDescriptor`/`SetSecurityDescriptorDacl` likewise (`:1157`, `:1169`), and
`AddAccessAllowedAce`/`InitializeAcl` genuinely are absent — but they are absent *because* the
exported Go-level API supersedes them. **This retires your own top-listed risk** ("hand-wrapped
advapi32 ACL construction … could produce ACLs that are permissive — privacy hole"): the correct
mitigation is to not hand-wrap. zcode-1 reached the same verdict independently; I record that as
convergence on a check either of us could run, not as a majority.

**V2. WS-C(i) "Go's `os` package does not expose share flags, so use `x/sys/windows.CreateFile`" —
PARTIALLY WRONG.** `os.Root` passes `FILE_SHARE_DELETE` (locators in position-change 1). Where a
call site can be expressed against a root — and `internal/trajectory` already is, e.g.
`snapshot.go:196` — prefer `os.Root` over a bespoke `CreateFile` wrapper: it is stdlib, already in
use in the package with the most Windows failures, and it carries traversal containment for free.
A raw `CreateFile` wrapper remains correct for sites that genuinely cannot be rooted.

**D1. WS-D allow-set "owner + LocalSystem" — I dispute the SYSTEM grant.** Your exclusion of
`BUILTIN\Administrators` is right and I ratify it with your own reasoning (admin-can-take-ownership
is the same residual as root on Unix; document, do not grant). The same argument removes SYSTEM: a
`0700` directory does not grant any second principal, nothing in parley runs as SYSTEM, and the
services that legitimately read the store (backup, Defender) do so through `SeBackupPrivilege`,
which a DACL cannot restrain and which therefore does not need a grant. **Counterproposal: exactly
one access-allowed ACE, for the token user, in a protected DACL.** Because creation sets
`PROTECTED_DACL_SECURITY_INFORMATION`, the create→verify round-trip is deterministic.

**D2. WS-E (Job Objects + creation-time/image-name attribution) — I dispute inclusion in this
idea**, for a reason that is mine and not zcode-1's: the facet you would add is `command`, and on
Windows the only cheap source is `QueryFullProcessImageName`, i.e. an **image path**, while
`procctl.commandMatches` (`procctl.go:164`) is calibrated against Linux **full argv**. Feeding a
coarser facet into an unchanged comparator silently *loosens* a kill-safety gate — the exact
property the base commit's U2 cycle spent a full review hardening. That is a strictness change
wearing portability clothes. **Counterproposal: P-A in this idea** (truthful `alive` via
`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` + `GetExitCodeProcess`, fail closed on
`ERROR_ACCESS_DENIED`; durable-kill refusal retained with Windows-honest strings), **WS-E as a named
deferral with its own adversarial review** (PID reuse, image-path collision, job lifetime). P-A
alone already reaches the four broken callers I located in round 1 — `durablekill.go:29`, `:63`,
`trajectory/verification.go:142`, `reviewsnapshot.go:244`.

**D3. WS-A.1 "a Unix-tagged variant with a recorded reason is a legitimate inapplicability
exclusion" / junction-or-reparse spike — I dispute that any exclusion is needed.** `PRIMARY`:
`$GOROOT/src/net/unixsock_posix.go:5` is `//go:build unix || js || wasip1 || windows`, and Go ships
`net/unixsock_windows_test.go` — so `net.Listen("unix", …)` is supported on Windows, which is what
the file's **already existing** `makeUnsupportedSocket` fallback (`tree_report_test.go:112-129`)
uses. `TreeDigest` refuses anything that is not a regular file or symlink
(`internal/evidence/tree.go:118,132` and `:161,168`), and a Windows AF_UNIX socket file is a reparse
point that `Lstat` does not report as regular. So the test **runs and keeps its assertion on
Windows**; only the `syscall.Mkfifo` *call* needs the build tag. Net effect: one fewer row in the
exclusion table, and no junction spike. **One caveat that cuts the other way:** that fallback ends
in `t.Skipf` (`:126`) — on Windows that skip must be converted to a hard failure, or the test
silently disappears. See New concern 1.

**D4. WS-A.6 "READY first, then parent-paced chunked writes" — counterproposal: split the test
instead.** `PRIMARY`, reading `internal/acp/spawn_test.go:61-137`: the gating observer parks the
copier inside its *first* `Write`, so from that instant every further child stderr byte must fit in
the OS pipe buffer or the child blocks before READY. Chunking does not remove that coupling; it only
moves it. The braid is that one test pins **two** properties at once — full drain before reap
(`observer.bytes == 16384`) and ring truncation (`len(p.Stderr()) == 8192`). Separate them:
(a) a *gated* test where the child writes a small payload (≤ 2 KiB, inside any plausible Windows
anonymous-pipe buffer), announces READY and **exits** before the copier parks — bytes then sit in
the pipe with no child left to block, and the drain-before-reap assertion is capacity-independent;
(b) an *ungated* test that writes 16384 and asserts the 8192 ring cap. Both assertions survive, on
every platform, with no magic number bet.

**A1. Accepted from you.** The `.tmp` staging observation (confirmed, `gate.go:119`). The
`Administrators` exclusion. WS-A.7's cross-compile `GOOS=windows build`+`vet` guard on the ubuntu
leg — cheap, catches the Mkfifo class pre-merge, and honestly labelled compile-only. Your two-layer
CRLF position, with the layers ordered: per-invocation `-c core.autocrlf=false -c core.eol=lf` on
parley's own git calls is the load-bearing layer (it is the only one that follows parley into the
*arbitrary user repos* it snapshots, which is where `hardening_test.go:455` lives), and
`.gitattributes` scoped to this repository is the second (it is the only one that fixes checked-in
fixtures such as the drift anchor and consensus template). **One rule I would add to your "product
parsers must be CRLF-tolerant":** normalization is admissible for *document-semantic* comparisons
(protocol drift, template block extraction) and **forbidden** for *evidence-integrity* comparisons
(snapshot content, `TreeDigest`) — otherwise CRLF tolerance quietly becomes a byte-exactness waiver
in a product whose premise is byte-exact capture.

**Not a dispute:** your 100 distinct failing names vs my ~98 vs zcode-1's 164 `--- FAIL` lines are
parse variance on one log. I support zcode-1's reconciliation ledger as the way to settle it.

### @zcode-1 — round-01 and round-02

**Accepted without reservation.** Your non-owner verdicts on Findings A–F (which discharge my §15.1
concern). Your withdrawal of the `strict_gate_test.go` hedge. Your round-02 concessions on universal
gate naming and owner-only ACLs. Your W3 self-rebasing in round-01, which independently reached the
conclusion my W3 verdict states. Your F.4 correction of my mechanism claim — adopted as a
self-correction above. Your `1.50.0`-class version reasoning (and note the first predecessor handoff
has now landed: `source-context/release-1.49.1-done.md` records 1.49.1 released at `54e0798`, so
"next version above what is released" is now anchored; the second handoff is still absent).

**D5. Round-02 position 2 (directory durability: attempt `CreateFileW(GENERIC_WRITE |
FILE_FLAG_BACKUP_SEMANTICS)` + `FlushFileBuffers`, hosted-decided) — I dispute this as a design
commitment.** This is the one point where you moved toward kimi-1 and I think the move is wrong on
the mechanism. Four facts, all `PRIMARY`:
- There are **four** directory-handle sync sites, not one: `internal/app/trajectory_verify.go:138`
  (`syncTrajectoryRuntimeParent`, `:134-139`), `internal/trajectory/verification.go:189`
  (`syncVerificationDirectory`, `:183-190`), `internal/trajectory/parent_recovery.go:308` and
  `:545`. Any decision must dispose of all four; your Stage 3 names one hosted failure.
- The Go stdlib **deliberately refuses** to open a directory with write access on Windows:
  `$GOROOT/src/syscall/syscall_windows.go:415-419` ("Unix doesn't allow opening a directory with
  O_WRONLY or O_RDWR … which will make CreateFile fail with ERROR_ACCESS_DENIED") and `:447-450`
  (`err = EISDIR`). So the proposal necessarily leaves the stdlib for a raw `CreateFile`, against
  the toolchain's own stated position.
- `os.File.Sync` is `$GOROOT/src/os/file_posix.go:163` → `pfd.Fsync()` → `FlushFileBuffers`, with no
  directory special case anywhere. There is no stdlib route at all.
- **A green spike would be weak evidence.** `FlushFileBuffers` returning `nil` on a directory handle
  does not establish that the *directory entries* created in it are durable; NTFS metadata
  durability comes from the log file and from write-through operations, not from flushing a
  directory handle. Accepting "the call did not error" as proof of the Unix barrier's semantics is
  exactly the over-reading of a green result that my R1 warns about.

**Counterproposal (D5).** Three parts, and none of them is "document a no-op":
   (a) **Name the contract in the API.** Today four sites call `fsutil.SyncFile` with a *directory*
       handle and the signature cannot tell. Introduce `fsutil.SyncDir` as a separate, per-OS
       function. Whatever Windows does then lives at one named, reviewable site instead of
       invisibly inside `SyncFile` — which is precisely what the owner's "documentation alone must
       not silently replace" constraint is protecting against.
   (b) **Prefer the real mechanism that already exists in-tree.** The Windows analogue of "fsync the
       parent after publishing" is to make the publication itself write-through, and the repo
       already does this: `internal/fsutil/replace_windows.go:14-26` and
       `internal/budget/lock_windows.go:37` both use `MoveFileEx(… MOVEFILE_WRITE_THROUGH)`, whose
       documented contract *is* "does not return until the change is on disk". Note the Unix
       directory barrier is already correctly confined to the Unix replace path
       (`fsutil/replace_unix.go:19,31`). So the work is a **per-site audit**: for each of the four
       sites, either route its durability through a write-through publication (real mechanism), or
       declare it cannot be.
   (c) **Only for sites where (b) fails** does the spike question arise, and the acceptance bar for
       calling it a "real implementation" must be a documented guarantee, not an errorless call.
       If it is not met, the choice between a user-visible refusal and a reviewed, disclosed
       deviation goes to the owner as an explicit FINAL line — which is your position 2's ending,
       and I agree with that part.

**D6. Bounded retry on `ERROR_SHARING_VIOLATION` — accepted only with an ordering condition.** A
retry helps exactly one case: a *foreign, transient* holder (AV, indexer). If our own process holds
the conflicting handle, a retry cannot succeed — it converts an instant failure into a 250 ms
failure and manufactures a "flaky" signature that your own round-1 risk 5 warns against. So the
diagnostic must classify **self-held vs foreign-held** before any retry ships, and structural
handle discipline (close before rename/delete; route through `os.Root` where share-delete is what
is needed) is the first fix, retry the narrow last resort.

**D7. Block-ID grandfathering — accepted with one hard carve-out.** Lenient load of existing
manifests plus an interpolation backstop is the right shape. But grandfathering must cover
*cosmetic* characters only: any ID that escapes its intended directory (`..`, a path separator, an
absolute or drive-qualified path, an alternate-data-stream `:`, or a Windows reserved device name)
must be refused at load **and** at the backstop, on every OS, with no legacy exemption — that class
is the security-adjacent surface, not a compatibility surface. `manifest.go:158-164` and
`executor.go:30` are the two interpolation points.

**On your nomination of me:** noted, and I do not treat it as settled — my own nomination is below,
with reasons. I will not turn this into a blocking dispute either way.

### Both peers — `strict_gate_test.go:179` is probably a product fail-open, not a fixture

Neither of you accepts my chmod refutation as *resolved* (zcode-1 now does; kimi-1's round-01 still
files it under "chmod fixtures"), and neither has proposed a cause. Re-confirmed: `grep -c Chmod
internal/driver/strict_gate_test.go` = **0** (`PRIMARY`). Here is the mechanism, built from source:

- Product: `internal/driver/impl.go:448-472`, `reviewRoundHasFindings` returns
  `!errors.Is(err, fs.ErrNotExist)` after `os.ReadDir` — "absent ⇒ no veto, unreadable ⇒ veto".
- Test: `strict_gate_test.go:164-181` writes a **regular file** where the round directory belongs
  and asserts the veto at `:179`.
- Therefore the hosted failure at `:179` means the product returned **false**, which requires
  either (i) `ReadDir` returned a nil error with no entries, or (ii) an error classified as
  `fs.ErrNotExist`. That disjunction is a deduction from the two sources above plus the hosted
  result, not a guess.
- Windows maps `ErrNotExist` only to `ERROR_FILE_NOT_FOUND`, `_ERROR_BAD_NETPATH`,
  `ERROR_PATH_NOT_FOUND`, `ENOENT` (`$GOROOT/src/syscall/syscall_windows.go`, `Errno.Is`) — `ENOTDIR`
  is not among them, which makes (i) the likelier branch. A supporting chain, flagged as a
  source-level **deduction** rather than an executed result: `os.ReadDir` → `openDir` →
  `openFileNolog(name, O_RDONLY|windows.O_DIRECTORY)` (`$GOROOT/src/os/file_windows.go:166`), but
  `O_DIRECTORY` is `0x04000` (`internal/syscall/windows/at_windows.go:17`) while
  `fileFlagsMask` is `0xFFF00000` (`syscall/types_windows.go:95`) — so the CreateFile path
  **ignores the bit** and opening a regular file succeeds; the subsequent
  `GetFileInformationByHandleEx` has a branch at `$GOROOT/src/os/dir_windows.go:134-145` that
  `break`s with **no error** on `ERROR_FILE_NOT_FOUND`, i.e. an empty listing with `err == nil`.

If that holds, a **fail-closed review gate is failing open on Windows** — a safety defect the
Windows leg surfaced, not a portability nit, and it must not be swept into a fixture bucket.
**Proposed fix, correct whichever branch is true:** `reviewRoundHasFindings` stats the path and
vetoes when it exists and is not a directory, making the invariant independent of per-OS `ReadDir`
error behaviour; plus a test that pins the veto for the file-where-directory case on all platforms.
I own this claim, so it needs a peer verdict — and a one-line hosted probe (`os.ReadDir` on a
regular file: error or empty?) settles the branch definitively.

---

## New concerns/questions

1. **The suppression gate as all three of us wrote it cannot catch the largest suppression risk.**
   We each specified "no *new* suppression / grep the diff for added `Skip`". But `grep -rn 't.Skip'
   --include='*_test.go' internal | wc -l` = **107** pre-existing skip sites (`PRIMARY`), and CI runs
   `go test ./... -count=1 -timeout 45m` with **no `-v`** (`.github/workflows/tests.yml`, Test step)
   — so `--- SKIP` is never printed and a skip that newly fires on Windows is **invisible**, with a
   zero-line diff. Two concrete candidates: `tree_report_test.go:126` (the AF_UNIX probe fallback)
   and `:250` (`"symlinks unavailable"`). The owner's constraint is "*a test may be excluded on
   Windows only when truly inapplicable and its reason is reviewed*" — excluded, not newly excluded.
   **Proposal:** add a **Windows skip census** to acceptance — emit `-json` (or `-v`) on the Windows
   leg, enumerate every skipped test, and require each one to carry a reviewed exclusion-table row
   exactly like a newly added one. This is added observability, not suppression, but it is a
   workflow change and should be reviewed as one.
2. **The owner-only DACL may be unverifiable for *pre-existing* stores on the hosted runner.**
   kimi-1's concern 2 is the right question and it now has teeth: our own creations will pass
   (protected DACL, one ACE), but the effective-grant verification zcode-1 and I both want will
   **refuse** any store sitting under a `%TEMP%`/`%LOCALAPPDATA%` path that grants SYSTEM or
   Administrators by inheritance. That is correct fail-closed behaviour and it may also mean the
   default hosted temp location cannot host a store at all. The create→verify round-trip spike must
   therefore also dump the *inherited* ACL of the parent, and FINAL needs a stated answer for users
   whose store predates this change (re-create with a protected DACL, or refuse and instruct).
3. **Does the `strict_gate` mechanism generalize?** If `os.ReadDir` error classification differs on
   Windows, every `errors.Is(err, fs.ErrNotExist)` decision after a directory read is suspect. Worth
   one sweep of that pattern during implementation, not a separate idea.
4. **The Windows shell refusal must not silently disarm the LE-4 gate.** `internal/app/driver_impl.go:266`
   is the `checks:` gate itself. If `sh` is absent and we refuse, `checks:` can never pass on
   Windows — which is the honest outcome, but it must be a loud, blocking refusal, never a
   pass-by-default. I endorse zcode-1's round-02 position (Git for Windows as a documented
   prerequisite + named refusal) as satisfying the owner's "real implementation or explicit,
   reviewed, user-visible refusal" for this release, with the `mvdan.cc/sh/v3` port deferred.
5. **Release-gate status (not a dispute, a fact that changed mid-round):** predecessor handoff 1 has
   arrived (`source-context/release-1.49.1-done.md`, 1.49.1 at `54e0798`); handoff 2
   (designated-implementer) is still absent, so the release gate remains closed.

---

## Current proposal

Unchanged in substance from my round-1 design, with the eight position changes folded in. Where
zcode-1's round-02 staging and mine agree I adopt its Stage numbering rather than restating:
Stage 0 contract + smoke + reconciliation ledger + ubuntu cross-compile guard; Stage 1 owner-only
protected-DACL privacy with adversarial extra-trustee refusal tests; Stage 2 universal gate
encoding inside `GatePath` + legacy read-fallback + block-ID validation; Stage 4 file-sharing after
diagnosis; Stage 5 P-A + shell contract; Stage 6 sweep; Stage 7 validation and release.

My deltas on top of that converged plan, each already argued above:

- **Step 1 first, always:** the `syscall.Mkfifo` build tag plus the comma-ok fix at
  `app_test.go:155/:157`. Highest information per line in the idea — it is what makes `wait`/`usage`
  produce Windows data at all, and it changes the denominator every later count depends on.
- **Directory durability (D5):** `fsutil.SyncDir` as a named per-OS contract; per-site audit of all
  four call sites against the existing `MOVEFILE_WRITE_THROUGH` publication; spike only where that
  fails; owner line in FINAL if no real mechanism survives.
- **Gate names:** universal encoding, plus — beyond both peers — the writer **retires** a legacy
  `a->b.gate.json` for the same edge on first rewrite rather than leaving a duplicate, so no stale
  HITL answer can survive in a shadow file.
- **Test mappings:** no FIFO exclusion row; ACP test split into gated (small, capacity-independent)
  and ungated (ring-cap) halves; `denyRead` via a deny-ACE sharing Stage 1's module; CRLF layered
  with `-c` pinning load-bearing and the normalize-vs-byte-exact rule stated in FINAL.
- **Acceptance:** add the **Windows skip census** (New concern 1) to the gate. Retain honest
  coverage language: hosted x64 `windows-2025-vs2026`, `core.longpaths true` (not the OS default),
  ARM64 asset build-only with zero executed evidence, no ARM64 gate added, no local-native claim.

**Readiness.** Drafting-ready on Stages 0–2, 4, 5, 6 and 7. Five decisions remain open, none of
which blocks starting, all of which must be settled in FINAL before their stage completes:
**(D1)** WS-E in-idea vs deferred — live dispute with kimi-1; my counterproposal is P-A + a
separate reviewed follow-up. **(D2)** directory-durability mechanism and, if needed, the owner
fallback line — live dispute with both peers, counterproposal D5. **(D3)** block-ID grandfathering
with my path-escaping carve-out. **(D4)** the Windows shell contract (owner-facing; §10 escalation
if the deck cannot settle it). **(D5-new)** whether the Windows skip census joins acceptance — I
consider it required by the owner's wording and would file a MAJOR finding at review if a green
leg were accepted without it.

**Nomination for drafting/implementation: `zcode-1`.** zcode-1 nominated me; I decline in its
favour on three grounds rather than modesty. (i) It pulled and classified the hosted Windows log
itself and holds the only stage decomposition covering the classes I missed entirely — directory
durability, worktree path fixtures, the tui perms test. (ii) The dominant risk in this idea is not
analysis, it is a ~33-file sweep quietly laundering product defects into skips (my R2, its risk 4,
kimi-1's "sweep masks product bugs"); the best available predictor of resisting that is
demonstrated willingness to reverse on evidence, and zcode-1 reversed three positions and withdrew
a hedge in full this round. (iii) Verdict economy: I am already barred under §15.1 from verdicting
Findings A–F; with zcode-1 implementing, kimi-1 and I are both unencumbered non-owner verifiers of
the implementation, whereas if I implement, the deck's verdict capacity narrows twice over the same
material. I accept the organizer selecting me instead without further dispute — but the record
should show the reasoning, not a rubber stamp. Model diversity holds either way
(`claude-opus-5[1m]` / `kimi-code/k3` / `zai/glm-5.3`).

**Scope compliance.** Read-only verification only (source at `d4bbf95`, GOROOT `go1.27.1`, x/sys
module cache, `git` metadata); wrote exactly this file. No product code, workflow, peer artifact,
phase transition, commit, tag, publication or hosted run. No driver or second organizer launched.
No secrets. No Chrome.
