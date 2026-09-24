---
agent: claude-1
idea: windows-portability
round: 3
date: 2026-09-25
responding-to: [kimi-1/round-02, zcode-1/round-02, round-03-brief]
---

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Read this round: `round-03-brief.md`, both peers' final `round-02` files, `00-prompt.md`,
`source-context/release-1.49.1-done.md`. Round-01 consulted only for ownership and exact quotation,
as instructed. HEAD moved `d4bbf95` → `2b1fc8c`; `git diff --name-only d4bbf95..2b1fc8c -- '*.go'`
returns **0 files** (`PRIMARY`), so every locator from rounds 1–2 transfers. All checks below
re-executed at `2b1fc8c`, darwin/arm64, GOROOT `go1.27.1`, x/sys module cache `v0.36.0`.
**New this round:** authoritative Microsoft Learn documentation fetched over HTTPS via `curl`
(no browser, no Chrome). No code, workflow, peer artifact, commit, tag, publication, hosted run or
driver touched. Wrote exactly this file.

**Concurrent peer round-03 files, disclosed rather than left to inference.**
`round-03/zcode-1.md` (00:06:54 local) and `round-03/kimi-1.md` (00:07:28 local) appeared on disk
while this file was being drafted; I saved at 00:09:00 and **did not open either**. Nothing here
responds to them, and any position of theirs that this file argues against may already have moved in
their own round-03 file. As in round 2 I state that rather than the stronger and false "they did not
exist yet". My cross-review target this round was the final round-02 set, as the brief instructed.

**Non-owner verification owed to me.** Four new load-bearing claims below need a peer verdict:
**N1** the `O_SYNC` → Windows write-through stdlib mapping, **N2** the six-site durability census,
**N3** the `o_DIRECTORY` refutation of my own round-2 `strict_gate` chain, **N4** the executed
traversal/`IsLocal` table. I state them as claims, not as settled facts.

---

## Position changes since round 2

**1. Directory durability — I withdraw the load-bearing half of my own D5 and replace it with a
documented stdlib mechanism.** Three separate corrections, all mine.

*(1a) `SELF-CORRECTION` — my site count was wrong.* I wrote "there are **four** directory-handle
sync sites". There are **six** reachable on Windows (`PRIMARY`, exhaustive `SyncFile(` grep):
`internal/app/trajectory_verify.go:138`, `internal/trajectory/verification.go:189`,
`internal/trajectory/unchanged.go:279`, `internal/trajectory/parent_recovery.go:308` and `:545`,
`internal/trajectory/reservation_recovery.go:64` (`reservation-intents` is a real directory,
created at `:36`). A seventh, `internal/budget/lock_unix.go:31`, is `//go:build !windows` and is not
reachable — which usefully confirms the Unix barrier is already confined by build tag in two places.

*(1b) `SELF-CORRECTION` — my "route durability through the existing write-through publication"
counterproposal covers 3 of the 6 sites, not all of them.* The split is **3 rename-publication / 3
create-publication** (`PRIMARY`). Rename: `parent_recovery.go:333` and `:570` (`dir.Rename`), and
`unchanged.go:279`, whose `persist` is `writeState` (`unchanged.go:217`), which publishes via
`fsutil.ReplaceSyncedFile` (`state.go:221`) — so that third site is **already** backed by in-tree
`MOVEFILE_WRITE_THROUGH` and its directory sync is belt-and-braces. Create-in-place with `O_EXCL`:
`trajectory_verify.go:146`, `verification.go:202`, `reservation_recovery.go:75` — there
`MoveFileEx(MOVEFILE_WRITE_THROUGH)` is structurally unavailable. My round-2 phrasing implied the
write-through route covered every site and was too strong.

*(1c) The replacement, and it is a **real mechanism with documentation scoped to the exact
operation** (`N1`, `PRIMARY`).* Both open routes the six sites use already map POSIX `O_SYNC` onto
Windows write-through, in the stdlib, with no raw `CreateFile` wrapper and **without leaving
`os.Root`**:
- plain path: `$GOROOT/src/syscall/syscall_windows.go:426` — `if flag&O_SYNC != 0 { attrs |=
  _FILE_FLAG_WRITE_THROUGH }`, with `O_SYNC = 0x01000` (`types_windows.go:52`);
- `os.Root` path: `$GOROOT/src/internal/syscall/windows/at_windows.go:99` — `if flag&syscall.O_SYNC
  != 0 { options |= FILE_WRITE_THROUGH }`. `os.Root.OpenFile` (`os/root.go`) validates only `perm`
  and passes `flag` straight through `rootOpenFileNolog` (`root_windows.go:132-137`) to `openat`.
- Executed on darwin: `os.OpenFile(..., O_WRONLY|O_CREATE|O_EXCL|syscall.O_SYNC, 0600)` and
  `Root.OpenFile(same)` both succeed, write/sync/close clean. Note `syscall.O_SYNC` is per-GOOS
  (`0x00080` darwin, `0x01000` windows), so it must be set from a **Windows-only** constant —
  adding it to the Unix path would change Unix write behaviour, which the owner forbids.
- Authoritative documentation, `CreateFileW` (learn.microsoft.com, fetched this round): "*A
  write-through request via `FILE_FLAG_WRITE_THROUGH` also causes NTFS to flush any metadata
  changes, such as a time stamp update or **a rename operation**, that result from processing the
  request.*" That is directory-entry-class metadata, named, for this flag.

So my position moves from "the spike is weak evidence, default to refusal" to: **there is a
documented real Windows mechanism; implement it.** Reviewed refusal stays as the *default*, not the
expectation. Honest residual, stated rather than buried: the doc scopes the flush to metadata
changes "that result from processing the request", and a *create* entry results from the
`NtCreateFile`, not from the subsequent write. The rename example is the closest documented analogue
and it is directory-entry metadata. I therefore call the rename-publication coverage **documented**
and the create-entry coverage a **documented-mechanism inference** that the hosted probe must
corroborate and FINAL must disclose. That distinction is the whole point of the brief's bar and I
will not blur it in my own favour.

**2. The `FlushFileBuffers`-on-a-directory route cannot meet the brief's documentation bar, and I can
now show it from the source rather than argue it.** Authoritative `FlushFileBuffers` documentation
(fileapi.h, fetched this round): the subject is "*a specified **file***"; the parameter is "*A handle
to the open **file**… must have the `GENERIC_WRITE` access right*"; the only other documented
subjects are a communications device (transmit buffer only), the server end of a named pipe, and a
**volume** handle ("*The caller must have administrative privileges*"). **Directory handles are not
documented at all** — not as supported, not as failing. A green spike therefore yields exactly
"the call returned nonzero", which the brief disqualifies by name. Separately, `FILE_FLAG_BACKUP_
SEMANTICS` documentation says only "*You must set this flag to obtain a handle to a directory*" — so
kimi-1's and zcode-1's open mode is **necessary and correctly chosen**; what is missing is the
guarantee, not the handle. This is a verdict against the *evidentiary status* of their route, not
against their reading of the API.

**3. `strict_gate_test.go:179` — `SELF-CORRECTION`, my round-2 mechanism chain is WRONG at
go1.27.1 (`N3`).** I claimed `O_DIRECTORY` is dropped by `fileFlagsMask` so the open of a regular
file succeeds. It is not: `o_DIRECTORY = 0x04000` (`syscall/types_windows.go:54`) is handled
**explicitly after** `CreateFile`, at `syscall_windows.go:455-466` —
`GetFileInformationByHandle` → if `FileAttributes&FILE_ATTRIBUTE_DIRECTORY == 0` → `CloseHandle` +
**`ENOTDIR`**. And `os/dir_windows.go:147-149` carries a second guard I also missed
(`if s,_ := file.Stat(); s != nil && !s.IsDir() { … ENOTDIR }`). At 1.27.1, `os.ReadDir` on a
regular file on Windows fails with `ENOTDIR`, which is **not** `fs.ErrNotExist`, so
`reviewRoundHasFindings` vetoes and the test should pass. My "likelier branch" framing is withdrawn.
I do **not** reclassify the failure as fixture-only — the brief forbids that without evidence and I
have none. What survives, narrowed: the hosted failure at `:179` is real (`TESTIMONY`, peers' log),
its cause is **undiagnosed**, the product's veto still depends on per-OS `ReadDir` error
classification, and CI pins `go 1.26` (`go.mod:3`, `tests.yml:42 go-version-file: go.mod`) while
every line above is 1.27.1. Control executed on darwin: `os.ReadDir(<regular file>)` → `entries=0`,
`err="not a directory"`, `IsNotExist=false`. See New concern 1 for the exact check.

**4. Skip census number corrected.** My round-2 said 107. `grep -rn 't\.Skip'` (escaped) gives
**104** lines, of which one (`internal/runner/launch_precheck_skip_test.go:58`, `result.Skipped`) is
not a skip call: **103 skip call sites** (68 `t.Skip(` + 35 `t.Skipf(`) across 56 files. 107 was an
unescaped-`.` regex artifact. The argument is unaffected; the number was mine and was wrong.

**5. Shell contract — I withdraw the §10 escalation framing (my round-2 D4).** kimi-1 is right that
the owner's "real implementation **or** an explicit, reviewed, user-visible refusal" already
authorizes refusal, and the brief instructs us not to re-ask it. Settled behaviour below instead.

**6. Experimental-label gate — I accept kimi-1's correction in full.** My round-1 concern-2
conditions (ARM64 execution, `core.longpaths`, `sh` presence) become **documentation obligations in
FINAL, not gates**. The gate is exactly the owner's: hosted windows leg green on the released commit.
I withdraw the extra conditions and I do not escalate. This closes the last dispute kimi-1 listed
against me.

**7. Drafter: I ratify `zcode-1`.** My round-2 declination stands on its stated reasons; the
organizer's procedural selection matches it. Both peers nominated me and I record that rather than
silently overriding it. I accept organizer selection either way; I will not contest this further.

---

## Responses to both peers

### @kimi-1 (round-02)

**Accepted, no reservation.** Your ALT-1 withdrawal and self-correction; your WS-E concession (P-A
only, P-B a named follow-up) — that closes D1, which both peers listed as live, and I record it as
**settled, not majority-resolved**; the FIFO simpler shape; the CRLF ordering with product git-pinning
load-bearing; the `proclive_windows.go` premise-comment correction as mandatory and probe adoption as
reviewer's call; your bundled early hosted probe commit (see additions below); the label-gate
correction (position change 6).

**ACL allow-set: settled owner-only.** You wrote "*Either closes the idea; I will not block on
owner-only*" and zcode-1 moved to owner-only. I record the allow-set as **closed: exactly one
access-allowed ACE for the token user, protected DACL**. No SYSTEM, no Administrators. Nobody is
overruled; you declined to block and zcode-1 converged.

**Durability — one remaining disagreement, and it is narrower than it was.** Your floor was "*if it
fails, your documented nil-with-rationale plus a FINAL known-risk entry is the floor — a reviewed,
disclosed weakening*". The brief removes that as an already-authorized option: a documented
no-op/weakened guarantee "is NOT an already authorized third option", and a deviation needs concrete
hosted evidence plus a reviewable proposed deviation **before** any owner question. I accept that
framing. It may also be moot: position change 1c supplies a documented mechanism, so the branch your
floor was written for should not be reached.

**Your "refusing the whole write path would be disproportionate" — agreed, and my refusal is not
that.** Concretely, the refusal is scoped to the publication operation whose durability cannot be
requested: `writeVerificationArtifact` (or the specific `sync*` helper) returns a named error, so the
affected trajectory/verification publication refuses, and everything that does not depend on that
barrier keeps working. It is not "trajectory verification is unusable on Windows" as a whole-feature
kill; it is one publication path failing loudly. If it turns out that scoping is impossible at a
given site, that site's blast radius goes in FINAL explicitly.

**Additions to your probe commit** (all build-tagged, test-only, no product change): (i) the
`O_SYNC` write-through publication probe — create with `O_EXCL|O_SYNC`, assert the open and write
succeed on NTFS and that a `SyncDir` call is no longer needed for the same guarantee; (ii) the
`os.ReadDir(<regular file>)` probe pinned to the **CI** toolchain, printing `len(entries)`, `err`,
and `errors.Is(err, fs.ErrNotExist)`; (iii) `git config --show-origin` dump alongside your hexdump.
Keep your `FlushFileBuffers`-on-directory probe — as a **diagnostic**, explicitly not as acceptance
evidence, per position change 2.

### @zcode-1 (round-02)

**Accepted, no reservation.** All six of your non-owner verdicts on Findings A–F and your verdicts on
my new round-1 claims; your four position changes (universal gate names, real-mechanism-first
durability, owner-only ACLs, product shell in scope) and your two self-corrections; your PathEscape
gap reproduction; your reconciliation-ledger-before-sweep demand; your mixed-name deck test; your
rename-root test redesign; your `1.50.0`-class version reasoning.

**Your D5 status after this round.** (a) `fsutil.SyncDir` as a named per-OS contract — **stands, and
is stronger**: six sites, not four, currently pass a *directory* handle to a function whose signature
says `SyncFile`. (b) write-through publication — **narrowed by me to 2 of 6 sites** (1b). (c) the
spike — **superseded as acceptance evidence** (position change 2) and replaced by the `O_SYNC` route
(1c), which does not cost you `os.Root` containment the way a raw `CreateFile` would.

**One correction to your framing, not to your judgement.** Your New-concern 3 says the fallback
"needs one explicit owner line in FINAL choosing between disclosed deviation and refusal; I recommend
the deviation". Under the brief the FINAL line is not a menu for the owner: the design is
**conditional** — documented real mechanism if proved, **reviewed refusal by default otherwise** —
and a deviation may only be *proposed* after concrete hosted evidence, as a reviewable artifact,
before any owner question. I am not disputing your recommendation on the merits; I am saying it
cannot be the pre-set default and no implementation may select it silently. You already said you
would not treat it as pre-approved, so I think this is agreement with the ordering made explicit.

**Block-ID grandfathering — accepted with my carve-out, now with an executed table (`N4`).**
`filepath.Join` computations at `2b1fc8c`, darwin (`\` is not a separator there — the Windows column
needs a hosted or `GOOS=windows` check):

| block/edge ID | `BlockWorkspace(deck,slug,id)` | `GatePath(deck,slug,id)` | verdict |
|---|---|---|---|
| `plain`, `a b`, `a%b`, `a<b`, `a\|b`, `start->spec` | inside `ideas/slug__…` | inside `gates/` | **grandfather** (cosmetic; encoded on write) |
| `..`, `../x` | `ideas/slug__../x` | `pipelines/slug/x.gate.json` — escapes `gates/` | **reject, all OS** |
| `../../x` | `ideas/x` | `pipelines/x.gate.json` | **reject, all OS** |
| `../../other-idea` | `ideas/other-idea` — **another idea's dir** | `pipelines/other-idea.gate.json` | **reject, all OS** |
| `../../../../x` | `x` — outside `deckDir` | `x.gate.json` | **reject, all OS** |
| `a/b`, `a\b`, `..\..\x` | separator in ID | separator in ID | **reject, all OS** (`\` is inert on Unix, live on Windows — a Unix-authored manifest escapes only once it reaches Windows) |
| `C:evil`, `C:\evil`, `a:b`, `a:$DATA` | drive-relative / ADS | same | **reject, all OS** |
| `CON`, `con`, `NUL`, `COM1`, `LPT9`, `AUX`, `PRN` | device name — reserved **even with an extension**, so `CON.gate.json` is still the console | same | **reject, all OS** |
| `end.` | Windows strips the trailing dot → collides with `end` | same | **reject new; escape on write for legacy** |

**Load-bearing negative result from the same run: `filepath.IsLocal` on the *concatenated* element is
an unsound guard.** `IsLocal("slug__" + "../../other-idea")` returns **true** (the `slug__` prefix
absorbs one `..` inside `Clean`) while `Join` lands in a sibling idea directory;
`IsLocal("../../other-idea")` alone is correctly `false`. Four of eight cases in the table show that
asymmetry. So validation must run on the **raw ID**, and the primary defence should be an explicit
charset allowlist rather than any path-shape predicate. (`filepathlite`'s own Windows `isLocal`,
`path_windows.go:23-53`, rejects *any* `:` and *all* reserved names — a useful precedent for the
allowlist's contents, cited rather than reused.)

**Bounded retry — my ordering condition is unchanged and both of you have now accepted its shape:**
classify self-held vs foreign-held before any retry ships; structural handle discipline first; retry
as the narrow residual for transient third-party holders, loud, with full verification re-run.

---

## New concerns / questions

1. **The `strict_gate` verdict needs one toolchain-pinned check, and I cannot run it.** Both guards
   that refute my round-2 chain are read at go1.27.1; CI resolves `go 1.26` from `go.mod`. There is
   no 1.26.x toolchain in this machine's module cache, and I did not download one. The decisive
   check is two lines on the Windows leg at the CI toolchain: `os.ReadDir` on a regular file →
   `entries`, `err`, `errors.Is(err, fs.ErrNotExist)`. Until then the cause stays open. The fix I
   proposed (`reviewRoundHasFindings` stats the path and vetoes when it exists and is not a
   directory) is correct under **every** branch and costs one `Lstat`, so it should land regardless
   of which branch the probe shows — that is a fix on the invariant, not on a symptom.
2. **Same toolchain caveat, applied honestly to the rest of the deck.** kimi-1 flagged it twice and
   was right both times. The claims that actually move if 1.26.8 differs: perm synthesis
   (`os/types_windows.go`) — discharged, because the 71 hosted failures *are* the 1.26.8 behaviour;
   `FILE_SHARE_DELETE` absence (`syscall_windows.go:395`) — long-stable, low risk; the `O_SYNC`
   mapping (`N1`) and the `o_DIRECTORY` guard (`N3`) — **not discharged**, both must be re-read at
   the CI toolchain before either is load-bearing in FINAL.
3. **Pre-existing-store DACL recovery still has no stated answer.** Carried from my round-2 concern
   2, now grounded: the guard sites are `internal/trajectory/snapshot.go:162` (directory,
   `Perm()&0077 != 0`) and `:544` (file). Our own protected-DACL creations verify deterministically;
   a store already sitting under an inherited `%TEMP%`/`%LOCALAPPDATA%` ACL will be **refused** by
   effective-grant verification. That is correct fail-closed behaviour, and it means FINAL must say
   what a user with a pre-existing store does — re-create under a protected DACL, or refuse and
   instruct. The probe should dump the *inherited* parent ACL so we learn whether the default hosted
   temp location can host a store at all.
4. **The AF_UNIX fallback's three `t.Skipf` calls are the concrete face of the skip-census problem.**
   `internal/evidence/tree_report_test.go:116`, `:119`, `:126` (plus `:250`, `"symlinks
   unavailable"`). On Windows each must become a hard failure or a reviewed exclusion row — a
   silently firing pre-existing skip produces a green leg with a zero-line diff, which is precisely
   what the census exists to catch.

---

## Current proposal — the acceptance contract, per the brief's six items

**1. Privacy (settled).** Windows creation writes exactly **one** access-allowed ACE for the token
user in a **protected** DACL (`ACLFromEntries` → `SetNamedSecurityInfo` with
`PROTECTED_DACL|OWNER|DACL`); no SYSTEM, no Administrators, documented as the admin/root residual.
Verification evaluates **effective granted trustees** via `GetNamedSecurityInfo`, refusing inherited
allowances and unknown trustees, naming the offending trustee in a sentence-stable error. FAT/exFAT/
no-ACL volumes: refuse. Adversarial tests required: grant `BUILTIN\Users` → refuse; grant
`Administrators` → refuse; inherited-only allowance → refuse; our own creation → pass. Existing-store
behaviour must be stated in FINAL per New concern 3. Unix path byte-identical.

**2. Process and shell (settled).** **P-A only**: truthful `Alive` via
`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` + `GetExitCodeProcess`, **fail closed** on
`ERROR_ACCESS_DENIED`; durable kill and attribution keep explicit Windows-honest refusals; the
`proclive_windows.go` premise comment corrected. **P-B (Job Objects, image-name/creation-time
facets) deferred** to a named follow-up with its own adversarial review. Shell, observable behaviour,
no owner question: on Windows, `checks:` and captured verification call `exec.LookPath("sh")` **once,
before any work**; on failure they return a named, actionable refusal naming Git for Windows as the
prerequisite; the refusal is **blocking** — `checks:` fails, it never passes by default, and
`internal/app/driver_impl.go:266` must not be reachable with a disarmed gate. The hosted leg must
exercise the refusal with `sh` forced off `PATH`, or green proves only that Git for Windows is
installed. `mvdan.cc/sh/v3` deferred.

**3. Directory durability (conditional, default refusal).** Introduce `fsutil.SyncDir` as a named
per-OS contract so the six sites stop passing directories to `SyncFile`. Then, per site, in order:
**(i)** request durability through a documented write-through publication — Windows-only `O_SYNC` on
the publication open for the three create-publication sites, `MoveFileEx(MOVEFILE_WRITE_THROUGH)` for
the three rename sites (already in-tree at `internal/fsutil/replace_windows.go:26`, whose reviewed
comment at `:10-13` states this exact rationale and deliberately excludes `COPY_ALLOWED`); **(ii)**
where (i) cannot be expressed, `SyncDir` returns a **named user-visible refusal** scoped to that
publication. There is no third branch. `FlushFileBuffers` on a directory handle may be probed as a
**diagnostic** and may not be cited as acceptance evidence. Accepting a mechanism requires
documentation scoped to the operation, quoted in FINAL, plus hosted confirmation; the create-entry
inference in position change 1c must be disclosed as an inference. Honest note on the in-tree
mechanism: `MOVEFILE_WRITE_THROUGH`'s explicit "guarantees … flushed to disk" sentence is scoped to
a copy-and-delete move, which `replace_windows.go` excludes by design; the sentence that applies to
the same-volume rename is "*does not return until the file is actually moved on the disk*". File-handle
fsync stays mandatory and unchanged everywhere.

**4. Names (settled, with the table above).** Universal encoding on every OS inside `GatePath`
(`internal/pipeline/gate.go:105`), so the `.tmp` staging name at `:119` inherits it by construction:
`PathEscape` base, plus explicit `:` → `%3A`, plus reserved-device-name and trailing-dot/space
escapes. Because separators are escaped, the encoder alone neutralises traversal **for gate
filenames**; block IDs still need validation, because `BlockWorkspace`
(`internal/pipeline/executor.go:30`) interpolates the ID **raw**. Validation therefore: strict
charset allowlist on new pipeline creation, lenient load of existing manifests
(`manifest.go:156-173`), and an enforcement backstop at both interpolation boundaries. Grandfathering
covers **cosmetic characters only**; everything in the reject column is refused at load and at the
backstop, on every OS, with no legacy exemption. One-release legacy `->` read-fallback, plus a
mixed-name deck test. **Legacy shadow files:** on first rewrite the writer **retires** the legacy
name for the same edge, so no stale HITL answer survives in a shadow file — and note the asymmetry
that makes this cheap: the shadow-file risk exists only on Unix (where the legacy name is writable),
while the unreadable-name risk exists only on Windows (where `>` cannot exist at all).

**5. `strict_gate_test.go:179` (open, verdict owed).** Not chmod — `grep -c Chmod` on that file = 0,
independently confirmed by both peers. Cause **undiagnosed**; my round-2 mechanism is refuted at
1.27.1 by my own check (`N3`) and unverified at the CI toolchain. It may not be reclassified as
fixture-only without evidence. Required: the toolchain-pinned `ReadDir` probe (New concern 1) plus a
peer verdict on `N3`. The invariant fix lands either way, with a test pinning the veto for the
file-where-a-directory-belongs case on all platforms.

**6. Tests, skips and evidence (settled).** **Windows skip census** joins acceptance: emit `-json` or
`-v` on the Windows leg, enumerate every skipped test, and require each to carry a reviewed
exclusion-table row exactly like a newly added one — because 103 pre-existing skip sites and a CI
step with no `-v` (`.github/workflows/tests.yml:52`) mean a newly firing skip is invisible with a
zero-line diff. Unfiltered matrix execution preserved; every claimed inapplicability reviewed
individually. **AF_UNIX: no exclusion row** — `net/unixsock_posix.go` is `//go:build unix || js ||
wasip1 || windows` and `TreeDigest` refuses non-regular/non-symlink entries, so the existing
`makeUnsupportedSocket` fallback keeps its assertion on Windows; only the `syscall.Mkfifo` *call*
needs a build tag, and the three `t.Skipf` calls must become hard failures there (New concern 4).
No junction spike unless AF_UNIX creation itself fails hosted. **ACP: split the test**, do not retune
it — `assertGatedChildDrained` (`internal/acp/spawn_test.go:126-138`) pins two properties at once
(`observer.bytes == 16384` drain-before-reap **and** `len(p.Stderr()) == 8192` ring cap) while the
gating observer parks the copier inside its first `Write` (`:68-72`), so every further child byte
must fit the OS pipe buffer; the test's own timeout message already names "pipe capacity below
in-flight bytes". Gated half: small payload, child announces READY and **exits** before the copier
parks — capacity-independent. Ungated half: 16384 bytes, ring-cap assertion. Both survive on every
platform with no magic number. **CRLF:** normalization is admissible for *document-semantic*
comparisons (protocol drift anchor, consensus-template block extraction) and **forbidden** for
*evidence-integrity* comparisons (snapshot content, `TreeDigest`); otherwise CRLF tolerance becomes a
byte-exactness waiver in a product whose premise is byte-exact capture. Load-bearing layer is
per-invocation `-c core.autocrlf=false -c core.eol=lf` on parley's own git calls;
`.gitattributes` scoped to this repo is second, as its own reviewed commit.

**Staging.** zcode-1's Stage 0–7 numbering, unchanged, with Stage 1 always preceded by the
`syscall.Mkfifo` build tag and the `app_test.go:155/:157` comma-ok fix — they change the denominator
every later count depends on.

---

## What I accept, what remains blocking, and the conditional acceptance

**I accept, and will not reopen:** owner-only protected DACL (item 1); P-A scope with P-B deferred
(item 2) — D1 is closed by kimi-1's concession, not by majority; the shell refusal as settled
observable behaviour with no owner question (item 2); universal encoding, grandfathering scope and
the reject list (item 4); the skip census, AF_UNIX no-exclusion, ACP split, CRLF split (item 6);
`zcode-1` as drafter/implementer; the label gate exactly as the owner set it, with my extra
conditions demoted to documentation.

**Remaining blocking counterproposal — one, and it is narrow.** A documented no-op or weakened
directory-durability guarantee may not be the default outcome or a pre-approved option. The contract
must read: documented real mechanism where it can be requested, **named user-visible refusal scoped
to the affected publication otherwise**, with any deviation proposed only as a reviewable artifact
backed by concrete hosted evidence, before any owner question. If FINAL ships item 3 with a
documented-weakening default, I file it as a **MAJOR** finding at review. On present evidence I
expect this branch not to be reached, because position change 1c supplies the mechanism.

**Exact conditional acceptance for the unresolved hosted facts.** I accept the plan as drafted, and
my acceptance is conditional on these and nothing else:
- **C1.** If the hosted probe shows the Windows-only `O_SYNC` publication open succeeding on NTFS at
  all three create-publication sites, I accept it as the real mechanism for item 3 — provided FINAL
  quotes the `CreateFileW` metadata-flush sentence and discloses the create-entry inference as an
  inference. If it fails or is refused by the filesystem, item 3 defaults to scoped refusal.
- **C2.** If the toolchain-pinned `ReadDir` probe shows an error that is not `fs.ErrNotExist`, I
  withdraw the product-fail-open hypothesis entirely and the `:179` failure returns to undiagnosed,
  with the invariant fix landing anyway. If it shows an empty listing with `err == nil`, or an error
  classified as `fs.ErrNotExist`, the fail-open is confirmed and it is a safety defect, not a
  portability nit. Either way it may not be filed as a fixture without this result.
- **C3.** If hosted evidence shows the effective-grant verification refuses the default runner temp
  location, FINAL must carry the pre-existing-store answer (New concern 3) before Stage 1 closes;
  the refusal itself is correct and is not weakened to accommodate the runner.
- **C4.** `N1` and `N3` must be re-read at the CI toolchain (`go 1.26.x`) before either is
  load-bearing in FINAL, and `N1`–`N4` need a non-owner verdict from kimi-1 or zcode-1.
- **C5.** Coverage language stays honest in FINAL: hosted x64 only, `core.longpaths true` set by us
  and not an OS default, ARM64 assets build-only with zero executed evidence, no ARM64 gate, no
  local-native Windows claim by anyone on this run, and the one documented sentence that this
  workspace itself sits on a cross-machine shared volume whose Windows semantics are unvalidated.

**Release gates, factual status:** predecessor handoff 1 is present
(`source-context/release-1.49.1-done.md`, 1.49.1 at `54e0798`); the designated-implementer handoff is
still absent, so the release gate remains closed and version selection stays anchored above 1.49.1.

**Scope compliance.** Read-only verification (source at `2b1fc8c`, GOROOT `go1.27.1`, x/sys module
cache, `git` metadata, three Microsoft Learn pages over `curl`); wrote exactly this file. No product
code, workflow, peer artifact, phase transition, commit, tag, publication or hosted run. No driver or
second organizer launched. No secrets. No Chrome.
