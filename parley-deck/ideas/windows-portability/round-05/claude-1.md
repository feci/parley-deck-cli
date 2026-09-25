---
agent: claude-1
idea: windows-portability
round: 5
date: 2026-09-25
responding-to: [kimi-1/round-04, zcode-1/round-04, durability-followup-brief, user round-cap authorization]
---

# Round 5 — claude-1: durability resolved, two narrow blockers left

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet sha256 re-verified this session by `shasum -a 256` against the file named in my task; it
matches the attestation. Read this round, in full: all three `round-04` artifacts,
`durability-followup-brief.md`, `inbox/user-to-codex-1_windows-portability_round-cap.md`,
`00-prompt.md` (including the *Authorized focused round 5* section) and `organizer-notes.md`.

HEAD is `6b87cf8` (was `ac1f0f7`/`2659520` when the round-04 files were written).
`git diff --name-only ac1f0f7..HEAD -- '*.go'` = **0 files** (`PRIMARY`); every changed file is an
idea artifact, brief, organizer note or inbox message. All round-01…04 source locators transfer;
I re-read every one I rely on directly at `6b87cf8`.

Environment: darwin/arm64, local GOROOT `go1.27.1`. Fetched myself this round (`PRIMARY`):
go1.26.8 sources for `os/root.go`, `os/root_windows.go`, `internal/syscall/windows/at_windows.go`,
`internal/poll/fd_fsync_windows.go`, `syscall/syscall_windows.go`; and the Microsoft Learn
`CreateFileW`, `FlushFileBuffers` and `MoveFileEx` pages over `curl` — this time extracting the
**whole enclosing paragraphs** rather than the sentences I had already selected, which is what
forces the largest correction below. One throwaway Go control executed under `/tmp`. No browser,
no Chrome.

**Concurrent peer round-05 file, disclosed rather than left to inference.** `round-05/zcode-1.md` appeared on disk at 08:43 local while this file was being drafted; I saved at
08:48 and **did not open it**. Nothing here responds to it, and any position of zcode-1's that this file
argues against may already have moved there. As in rounds 2, 3 and 4 I state this rather than the
stronger and false "it did not exist yet". My cross-review targets this round are the three completed
`round-04` artifacts and the durability brief, as the authorization instructs. `round-05/kimi-1.md` was
absent when I saved.

**Ownership under §15.1.** I own `N1`–`N6`; I issue no verdict on them, only self-corrections,
which take effect immediately when they weaken. Verdicts below are on **peers'** claims. Claims I
introduce this round are flagged owner-asserted with a one-command check attached, because this is
the last authorized round and they cannot receive a non-owner verdict inside Phase 2.

---

## User direction

The owner authorized this round in `parley-deck/inbox/user-to-codex-1_windows-portability_round-cap.md`
(relayed 2026-09-25). The question put to the owner was whether to allow the Windows idea one more
cross-review round over the §4.0 cap, limited to the write-durability dispute.

**The owner's exact English translation of the selected answer, as supplied in that file:**

> "Allow one round — one focused fifth round limited to the prepared durability brief. Then
> consensus, FINAL, implementation and review as before."

**Citation of the original answer.** The owner answered in Slovak and supplied that translation
himself. The original selection quoted in the same file reads: *"Povoliť jedno kolo" — "Jedno
cielené piate kolo obmedzené na pripravený brief o trvácnosti. Potom konsenzus, FINAL, implementácia
a review ako doteraz."* It answers the question recorded there as *"Povoliť Windows idei ďalšie kolo
cross-review nad limit, len na spor o trvácnosť zápisov?"* — "Allow the Windows idea one more
cross-review round over the cap, only for the write-durability dispute?" The Slovak appears here
solely as the cited source string; the operative text for this artifact is the English above.

**How I read the scope.** Exactly ONE additional round, limited to `durability-followup-brief.md`,
all material findings still reportable, not a standing waiver: if durability is unresolved after
this, escalate again — there is no round 6. Implementation, hosted-CI, release-sequencing and
channel-verification gates are unchanged. I therefore treat this file as my **final** Phase-2
artifact and I resolve rather than defer wherever the evidence permits.

---

## Position changes since round 4 — self-corrections first, and the first one is against myself

### 1. `SELF-CORRECTION`, and it reverses my round-04 headline: **`N5` is withdrawn. My peers' tier was right and mine was wrong.**

My round-04 position change 1 claimed the `CreateFileW` page "documents against" write-through
covering a create's directory entry, resting on:

> "Also, the file metadata may still be cached (for example, when creating an empty file). To ensure
> that the metadata is flushed to disk, use the FlushFileBuffers function."

I flagged in round 4 that this sentence sits in the `FILE_FLAG_NO_BUFFERING`-plus-`FILE_FLAG_OVERLAPPED`
paragraph, and then concluded "the create-entry claim fails" anyway. Having now extracted the
paragraph in full (`PRIMARY`, fetched and de-tagged this round), that conclusion does not survive.
The paragraph is:

> "When FILE_FLAG_NO_BUFFERING is combined with FILE_FLAG_OVERLAPPED, the flags give maximum
> asynchronous performance… Also, the file metadata may still be cached (for example, when creating
> an empty file). To ensure that the metadata is flushed to disk, use the FlushFileBuffers function."

That is a statement about a configuration in which **no write-through was requested at all**. It is
not a counterexample to write-through; it is a description of what you get without it. And two
sentences I did not have in round 4 cut the other way:

- `CreateFileW`, the paragraph immediately before: after the NTFS metadata-flush sentence — "For this
  reason, the FILE_FLAG_WRITE_THROUGH flag is often used with the FILE_FLAG_NO_BUFFERING flag **as a
  replacement for calling the FlushFileBuffers function after each write**".
- `FlushFileBuffers`, Remarks: "To open a file for unbuffered I/O, call the CreateFile function with
  the FILE_FLAG_NO_BUFFERING and FILE_FLAG_WRITE_THROUGH flags. This prevents the file contents from
  being cached **and flushes the metadata to disk with each write**."

Both are equivalence statements between write-through and `FlushFileBuffers` for metadata; both are
scoped to **writes**, and neither names directory-entry creation. So the create-entry coverage is
**neither documented nor documented-against** — which is precisely the "documented-mechanism
inference" tier kimi-1 assigned in its N1 verdict and zcode-1 adopted "word for word". **I withdraw
`N5` in full and adopt their tiering.** Consequence: the class-A dispute is over, and it closed
against me.

### 2. `SELF-CORRECTION`: my "13 barrier operations" was itself an undercount. The mechanical count is **15 invocation sites / 16 runtime operations**, and a naive mechanical grep still misses one.

Round 4 I criticised hand-lists and then filed one. Re-derived at `6b87cf8` (`PRIMARY`):

```
grep -rnE 'syncTrajectoryRuntimeParent\(|syncVerificationDirectory\(|syncUnchangedState\(|syncParentRecovery\(|syncRecoveredParent\(|syncIntent\(' \
  --include='*.go' internal | grep -v '_test.go' | grep -v ':func '     → 15 lines
```

Of those 15, `reservation_recovery.go:67` is *inside* `syncIntent`'s own body, not an independent
operation → **14 invocation sites found by name**. Add the one the name-grep **cannot** see:
`parent_recovery.go:311` passes `syncParentRecovery` as a *function value* into
`publishParentRecoveryWithSync`, which invokes it at `:336` as `sync(dir, base)` — the publication
barrier for the `:333` rename is therefore invisible to any name-keyed audit. That gives **15
invocation sites**. At runtime it is **16 operations**, because `trajectory_verify.go:125` sits in a
two-iteration loop over `[filepath.Dir(base), base]` (`:117-128`).

What I got wrong in round 4: I missed the `:118/:125` loop's second iteration and the `:420`
`syncIntent` recovery call, and my class-A locator `trajectory_verify.go:141` is the function
signature — the create is `:146` and the barrier call `:154`.

This is the **fourth** census correction in this idea (kimi-1's 4, zcode-1's ~7, my 6, my 13) and the
third of them mine. It settles the methodology question in zcode-1's favour and against my own
list: **no participant's count should reach FINAL as a headline number.** More than that, the
function-value miss shows a name-keyed grep is *also* insufficient — which is exactly why zcode-1's
`fsutil.SyncDir` named type contract is the right audit: it is enforced by the compiler, not by a
pattern.

### 3. `SELF-CORRECTION`: my round-04 class-B blast radius was overstated by roughly a feature.

I wrote that a class-B refusal means "trajectory verification and budget reservation cannot be
created on Windows at all". Re-checked at HEAD (`PRIMARY`): `openVerificationDirectory` already
refuses on Windows at `verification.go:251-253` ("captured verification journals require a POSIX
execution host"), and `trajectoryRuntime` already refuses at `trajectory_verify.go:106-108`. Every
`internal/app/trajectory_verify.go` operation and every `openVerificationDirectory` operation is
therefore **Windows-unreachable today**. Only **9 of 16** operations are Windows-reachable, and only
**one** class-B operation is: `reservation_recovery.go:36` → `:47`. The real blast radius is
precharge reservation-intent publication, not trajectory verification. kimi-1's reachability
refinement 2 was the correction, and I confirm and extend it below.

### 4. Everything settled outside durability stands and I reopen none of it.

Owner-only protected DACL with refuse-and-instruct; P-A with P-B deferred; blocking missing-`sh`
refusal; universal gate-name encoding plus the raw-ID charset allowlist; the branch-independent
`strict_gate` invariant fix with cause recorded undiagnosed; the ACP split; the CRLF split; the label
gate exactly as the owner set it; `zcode-1` as drafter and single implementation owner with kimi-1
and me reviewing.

---

## Responses to peers' round-04 positions

Both peers wrote round-04 concurrently with me and state they did not open my file. Nothing below
treats that as evasion; it is why `N5` and `N6` reached this round unverdicted, and it is the whole
reason the owner authorized round 5.

### @kimi-1 — round-04

**`CONFIRMED`: your N1 create-entry tiering, and it supersedes my round-04.** Your "rename coverage
is **documented**, create-entry coverage is a **documented-mechanism inference**" is the correct
reading of the page as a whole — see my self-correction 1. Your NTFS scoping is `CONFIRMED` verbatim
(`PRIMARY`, my fetch): the sentence is "also causes **NTFS** to flush any metadata changes", and ReFS
appears in the same page's supported-technology table without being covered by it. Your
strengthening — that on a create the entry results from processing that very handle's request — is
the strongest honest reading and still stops short of the examples given. FINAL discloses it as
inference and never upgrades it.

**`CONFIRMED`, and this is the one that changes my own proposal: your rejection of rename-reshaping
the class-A create sites.** You wrote that it "changes replay-prevention semantics, not just
plumbing". You are right, and my round-04 Current proposal 3 was wrong to put class A in the
restructure column. Read at HEAD (`PRIMARY`): under `O_EXCL` create-in-place, a crash mid-write
leaves the **final** name occupied by a torn artifact, and that is load-bearing by explicit design —
`verification.go:192-193` "*Exclusive writes leave any partial file in place. Its existence prevents
replay; a torn/noncanonical artifact is an unresolved failure, never a retry*", and
`reservation_recovery.go:79` "*A partial publication stays visible and cannot be overwritten by
retry.*" Under stage-then-rename the crash leaves only the **stage** name occupied — and
`parent_recovery.go:328`/`:565` even `defer dir.Remove(stage)` — so the final name stays free and a
retry **succeeds**. That is an anti-replay regression. **I withdraw the class-A half of my round-04
restructure proposal.**

**`CONFIRMED` with an extension: your reachability refinement 2.** Your two gated sites and ungated
reservation path are right. It is stronger than you state: the gate at `trajectory_verify.go:106-108`
covers that file's entire barrier set (`:125` ×2, `:154`, `:193`), because `verifyTrajectoryWithAgent`
returns at `:186-188` when `trajectoryRuntime` fails, so `:190/:193` are downstream of it too. And I
confirm your budget chain by source: `trajectory.Observer.PrepareCycleReservation`
(`reservation_recovery.go:106`) → `openIntentRoot(b, true)` (`:126`) → `:36` `Mkdir` + `:47` barrier →
`publishReservationIntent` (`:131`) → `:75` `O_EXCL` create + `:84` `syncIntent`; invoked from
`internal/budget/cycle_binding.go:269` whenever the observer implements `CycleReservationObserver`;
no `runtime.GOOS` gate anywhere on that chain (`grep` over `internal`: the only Windows gates are
`driver_impl.go:542`, `trajectory_verify.go:106`, `evidence_verify.go:368`, `source.go:252/:400`,
`verification.go:251`, `acp/spawn.go:178`, `acp/shellenv.go:44/:120`).

**`CONFIRMED`: your census file count of 55, and zcode-1's 56 is also right.** Re-executed
(`PRIMARY`): `grep -rl 't\.Skip'` = **56 files**; `grep -rl 't\.Skip(\|t\.Skipf('` = **55 files**. The
56th is `internal/runner/launch_precheck_skip_test.go`, matching only `result.Skipped` at `:58`. You
counted files with call sites, zcode-1 counted files matching the pattern. FINAL must pin both with
the denominator named. Everything else reproduces exactly: 104 escaped lines, 103 call sites, 68
`t.Skip(` + 35 `t.Skipf(` + 0 `t.SkipNow(`.

**`CONFIRMED` and adopted over the Windows-only variant: your `-json` on all three legs.** Uniform
observability, `-json` selects no tests, the unfiltered matrix is untouched, and the enumerated
baseline gives the owner's "no new suppression" real teeth on the still-green legs. Your New
concern 2 (assert the runner volume's filesystem) is `CONFIRMED` as **necessary**, not optional —
the metadata sentence you correctly scoped to NTFS is the only documented text the class-A mechanism
rests on, so an unasserted volume means unmatched documentation.

**Your New concern 1 is right in substance and is where I still press.** You name the rename-site
containment trade-off as "an open implementation choice… the drafter must pick one and record the
trade-off in FINAL". Naming it is better than zcode-1's treatment of it. But the choice is between
keeping a kernel-enforced containment guarantee and keeping a durability barrier on anti-replay code,
on a security surface the owner scope calls out by name — that is a design decision, not an
implementation shape, and my round-04 objection to delegating it survives unchanged. See Blocker 1.

**One gap, stated as a gap rather than a verdict.** Your contract's three branches are keyed on
"rename-publication sites", "create-in-place sites", and "any site where (i)/(ii) cannot be expressed
or proved". The five `Mkdir` barrier operations are in neither (i) nor (ii), so they land in (iii) —
a named blocking refusal — by your contract's own logic. I think that is the right answer. But your
file does not enumerate them or state the blast radius, so a reader of FINAL would not know that
your contract refuses a feature.

### @zcode-1 — round-04

**`CONFIRMED`, and it is the best single idea in either round-04 file: `fsutil.SyncDir` as the
refusal emitter that never returns nil on Windows.** "There is **no branch** in which Windows
`SyncDir` returns nil: converted sites stop calling it on Windows; unconverted sites refuse through
it. Fail-closed by construction, including for future call sites added without conversion." I adopt
this without reservation. My self-correction 2 is independent evidence for it: a name-keyed grep
missed the `:311`/`:336` function-value invocation, so only a compiler-enforced type contract audits
this completely. Your mechanical-audit demand was right and I was the one who violated it.

**`CONFIRMED` and strengthened: your supersession 2, the spike's bar is unsatisfiable.** You argued
directory handles are absent from the `FlushFileBuffers` subject list. It is stronger than absence.
The page states an affirmative precondition (`PRIMARY`, my fetch): "**The file handle must have the
GENERIC_WRITE access right.**" Every handle at all 16 operations is read-only, so the documented
precondition is not merely unaddressed — it is **unmet**. Detail in New finding 1.

**`CONFIRMED`: your partial-publication recovery analysis, with one hole — see New finding 3.**
`O_EXCL` leaving a visible non-replayable artifact, randomized staging names, `defer dir.Remove(stage)`
all verified at HEAD. The hole is the *ordering* of a refusal at the rename sites, which your
contract does not specify.

**`WRONG` — your containment sentence, and this is the brief's named question.** Verdict **`WRONG`**
(`PRIMARY`: go1.26.8 sources fetched this round + HEAD). Your contract §2 says the write-through
replace at rooted sites "must keep both from- and to-paths inside the root — resolved under
`Root.Name()` with the existing sibling-staging check (`replace_windows.go:19-21`) kept — **so the
conversion does not trade away `os.Root` containment**; the helper's signature is code shape, not a
design decision."

- *Locator.* The sibling check is at `replace_windows.go:15-17`; `:19-21` is the `err != nil` return
  for `windows.UTF16PtrFromString(staged)`. kimi-1's `:16-18` is also off. Immaterial to the argument,
  recorded because FINAL carries exact numbers.
- *Substance.* The check is `filepath.Clean(filepath.Dir(staged)) != filepath.Clean(filepath.Dir(path))`.
  `filepath.Clean` is documented as "the shortest path name equivalent to path **by purely lexical
  processing**" — it never touches the filesystem and resolves no reparse point.
- *What it would replace.* `os.Root.Rename` (`os/root.go:220-225`, "Both paths are relative to the
  root") → `rootRename` → `renameat` (`os/root_windows.go:382-383`) → `windows.Renameat`
  (`at_windows.go:364`), which opens **both** sides handle-relative with `FILE_OPEN_REPARSE_POINT`
  (`:376`); and every `os.Root` open goes through `openat` (`os/root_windows.go:147`) with
  `O_NOFOLLOW_ANY`, which sets `OBJ_DONT_REPARSE` on the NT `OBJECT_ATTRIBUTES` (`at_windows.go:109-110`).
  The documented guarantee (`os/root.go:34-43`): "If any component of a file name passed to a method
  of Root references a location outside the root, the method returns an error… symbolic links may not
  reference a location outside the root."
- *Therefore.* The conversion swaps a kernel-enforced, reparse-resistant, handle-relative rename for
  a lexical string comparison followed by a path-based Win32 call that resolves reparse points at
  call time and is TOCTOU-exposed between the check and the move. **That is trading away `os.Root`
  containment.** The brief's sentence — "sibling path strings alone must not be presumed equivalent
  to an `os.Root` guarantee without participant review" — is exactly this, and the review outcome is:
  **not equivalent.**
- *Honest counterweight, which I owe you.* The repo already ships path-based
  `MoveFileEx(MOVEFILE_WRITE_THROUGH)` at `internal/budget/lock_windows.go` (`publishExclusive`, with
  the in-tree comment "*No REPLACE_EXISTING: a competing complete file wins and must match*"), and
  the Unix counterpart is `os.Link` + directory fsync. So the mechanism is in-tree and reviewed on
  both OSes — my round-04 was wrong to imply it was new machinery. But `publishExclusive` operates on
  plain paths and never had an `os.Root` guarantee to lose; the two trajectory rename sites do. In-tree
  precedent for the **mechanism**, none for the **trade**.

**`WRONG` — "no design decision is delegated to implementation".** Three are, and each is visible in
your own file: (a) the disposition of the five `Mkdir` operations, which your §3 fail-closed rule
resolves to refusal without your contract saying so or costing it; (b) the class-C containment trade
just above, which you resolve by an equivalence that does not hold; (c) which class-D sites may
derive a no-op, which you do not enumerate. Your round-03 sentence "no unresolved design decision
remains" was the one I contested in round 4; the round-04 restatement has the same problem.

**No dispute** on your N1–N4 verdicts (they are on my claims; I issue none). I record that both peers
independently confirmed N1 and N3 at go1.26.8, which discharges my round-03 C4 completely, and that
both independently caught the `:64`→`:63/:67` locator error.

### §15.3 status — my round-04 New concern 1 is discharged

I asked that kimi-1's V1, V2 and V4 be struck or carried as conflicts. kimi-1's round-04 position
changes 1, 2 and 3 strike all three independently, before reading my file. **No verdict conflict
survives into consensus** and `consensus.md` needs no `## Verdict conflicts` section on that account.

---

## Operation guarantees versus API success

This is the brief's central question and it has one answer that covers every row of the table: **on
Windows, a returning `SyncFile(<directory handle>)` has never been shown to be a barrier, and a
non-returning one has never been shown to be observed.** Concretely, three separate things have been
conflated across five rounds, and FINAL must keep them apart:

1. **A documented guarantee** — text scoped to the exact operation. We have exactly two:
   `MOVEFILE_WRITE_THROUGH`'s "The function does not return until the file is actually moved on the
   disk" for the same-volume move, and `FILE_FLAG_WRITE_THROUGH`'s NTFS metadata-flush sentence for
   metadata resulting from processing a write request.
2. **A documented-mechanism inference** — the create's directory entry under write-through. Real,
   disclosed, never upgraded. (Both peers' tier; now mine.)
3. **A successful API call** — `err == nil`. This establishes nothing about persistence and may never
   be cited as acceptance evidence. The hosted probe can only ever reach tier 3 plus mechanics.

`FlushFileBuffers` on a directory handle does not reach tier 1, tier 2, or even tier 3 — see below.

---

## New findings

### New finding 1 — the current Windows barrier violates a documented precondition at all 16 operations, and **nobody has observed the result**. (owner-asserted, `N7`)

The chain, every link `PRIMARY` and re-read this round:

- `fsutil.SyncFile(f)` on Windows resolves to `internal/fsutil/sync_other.go:8`, `file.Sync()`.
  There is **no `sync_windows.go`** (`ls internal/fsutil/`); `sync_other.go:1` is `//go:build !darwin`
  and `sync_darwin.go` is filename-constrained.
- `(*os.File).Sync` → `poll.FD.Fsync` (`internal/poll/fd_fsync_windows.go`, go1.26.8) →
  `syscall.Fsync(fd.Sysfd)` → `syscall_windows.go:755-757`, `func Fsync(fd Handle) error { return
  FlushFileBuffers(fd) }`.
- `FlushFileBuffers` documentation: "**The file handle must have the GENERIC_WRITE access right.**"
  Documented subjects: a file, a communications device, the server end of a named pipe, a volume
  (administrative privileges). No directory.
- The handles are read-only at every site. Plain: `os.Open(filepath.Dir(path))`
  (`trajectory_verify.go:133`, `unchanged.go:275`, `parent_recovery.go:304`/`:541`). Rooted:
  `dir.Open(".")` (`verification.go:184`) and `dir.Open("reservation-intents")`
  (`reservation_recovery.go:63`). `os.Root.Open` → `Openat` with `O_RDONLY`, and `at_windows.go:65-68`
  gives `access |= FILE_GENERIC_READ` only — `FILE_GENERIC_WRITE` is added only for `O_WRONLY`,
  `O_RDWR` or `O_CREAT` (`:69-82`).
- In-tree reviewed corroboration (`SECONDARY`, a prior reviewed change in this repo, not verified by
  me on Windows): `internal/fsutil/replace_windows.go:11-12` already states "*Windows directory
  handles opened for reading cannot supply a FlushFileBuffers barrier. Request write-through on the
  replacement operation instead.*"
- Executed darwin control (`PRIMARY`, `/tmp`, this round): both shapes return `err=<nil>` on darwin —
  `os.Open(dir)`+`Sync()` and `OpenRoot(dir)`+`Open(".")`+`Sync()`. `GOOS=windows go vet` of the same
  program is clean. So the identical Go code is a real barrier on Unix and a
  documented-precondition-violating call on Windows, and it compiles either way.

**And here is the part that should have ended this dispute in round 1:** the hosted evidence contains
**no observation of it at all**. `grep -rni 'flushfilebuffers|access is denied|sync'` over
`source-context/release-ci-revalidation-claude-1.md` returns **zero hits**. The recorded Windows
failure classes are snapshot privacy (69 + 2 at `snapshot.go:163`/`:544`), the `internal/evidence`
build failure, the `internal/app` panic chain, 14 file-sharing failures, 10 `/bin/sh` fixtures and one
CRLF mismatch — most plausibly because `snapshot.go:163` refuses 69 times upstream of any publication.

Two consequences, both of which belong in FINAL:

- **Five rounds of durability argument have run on documentation with zero hosted observation of the
  behaviour in dispute.** That is a census-baseline gap, and it is one line to close: the probe must
  execute `SyncFile` on a read-only directory handle on the Windows runner and print the raw error.
- **A named blocking refusal is not a regression.** The documented precondition is already unmet, so
  the most likely current Windows behaviour at these sites is an *unnamed* error. Refusal converts an
  unnamed failure into a named, user-visible one. This materially lowers the cost of the refusal
  branch and I record it in the peers' favour.

*Check (one command, cheap for a non-owner verdict):* the greps and `sed` ranges above at `6b87cf8`,
plus `curl` of the `FlushFileBuffers` page.

### New finding 2 — `N6` re-verified in full; still **`UNVERIFIED`** by a non-owner

Every `N6` sub-claim reproduces verbatim on the `MoveFileEx` page I fetched this round (owner-side,
so no verdict from me): "Moves an existing file **or directory**, including its children"; "When
moving a directory, the destination must be on the same drive"; `MOVEFILE_REPLACE_EXISTING` — "If
lpNewFileName names an existing **directory**, an error is reported"; `MOVEFILE_WRITE_THROUGH` — "The
function does not return until the file is actually moved on the disk. Setting this value guarantees
that a move performed as a copy and delete operation is flushed to disk before the function returns."
Neither peer verdicted `N6`, because both wrote concurrently. It is a strengthening, so under §15.1 it
stays `UNVERIFIED`. It is now **less load-bearing** than in round 4, since I no longer propose
restructuring class A, but FINAL must still label it owner-asserted, and the verdict belongs in
Phase 6 review where kimi-1 and I are the reviewers.

### New finding 3 — a class-C refusal fires **after** the rename. Neither contract orders it. (owner-asserted, `N8`)

`publishParentRecoveryWithSync` (`parent_recovery.go:314-337`) renames at `:333` and **then** calls
`sync(dir, base)` at `:336`. `publishRecoveredParent` renames at `:570` and calls
`syncRecoveredParent` at `:573`. Under zcode-1's §3 — Windows `SyncDir` is the refusal emitter — the
refusal at those two sites fires **after the artifact is already published**, so the product publishes
without a barrier and returns an error. That is precisely the torn state the owner scope forbids from
being silent, produced by the refusal mechanism itself.

The fix is small and belongs in the design, not the implementation: at any operation whose disposition
is refusal, **the refusal must be evaluated before the mutation**, so nothing is published. For class C
that means the Windows disposition is checked before `:333`/`:570`, not at `:336`/`:573`. For class B
it is already safe (`Mkdir` then barrier, and `Mkdir` is idempotent under `!os.IsExist` at `:36`/`:264`).
For class A it is already safe (`O_EXCL` is the anti-replay token by design; a refusal after it leaves
exactly the visible non-replayable artifact the comments require).

*Check:* read `parent_recovery.go:314-337` and `:551-574` at `6b87cf8`.

---

## The complete operation table

Individuation: one row per **barrier operation at runtime**, keyed on (mutation → barrier). 15
invocation sites, 16 runtime operations; derivation and the function-value caveat in self-correction 2.
"Gated" = unreachable on Windows today behind an existing `runtime.GOOS` refusal.

| # | Barrier | Mutation completed | Class | Windows reachability | Required guarantee | Mechanism **or** named blocking refusal | Feature blast radius | Recovery from partial publication |
|---|---|---|---|---|---|---|---|---|
| 1 | `trajectory_verify.go:125` (iter 1) | `:118` `os.Mkdir(filepath.Dir(base))` | B | Gated `:106-108` | parent entry durable | **Refuse** (`SyncDir`), pre-mutation | none today; inherits if gate lifts | `Mkdir` idempotent under `!os.IsExist`; nothing published |
| 2 | `trajectory_verify.go:125` (iter 2) | `:118` `os.Mkdir(base)` | B | Gated | parent entry durable | **Refuse**, pre-mutation | none today | as #1 |
| 3 | `trajectory_verify.go:154` | `:146` `os.OpenFile O_EXCL` | A | Gated | create entry durable | `syscall.O_SYNC` on the open (documented for the write; **inference** for the entry); drop the Windows barrier call | none today | `O_EXCL` leaves torn artifact visible and non-replayable |
| 4 | `trajectory_verify.go:193` | `:190` `os.Mkdir(dir)` | B | Gated | run-dir entry durable | **Refuse**, pre-mutation | none today | as #1 |
| 5 | `verification.go:211` | `:202` `Root.OpenFile O_EXCL` | A | Gated `:251-253` | create entry durable | `O_SYNC` as #3 | none today | `:192-193` invariant |
| 6 | `verification.go:267` | `:264` `parent.Mkdir("trajectory-verifications")` | B | Gated `:251-253` | store entry durable | **Refuse**, pre-mutation | none today | as #1 |
| 7 | `verification.go:286` | `:283` `base.Mkdir(charge)` — **anti-replay token** | B | Gated `:251-253` | reservation entry durable | **Refuse**, pre-mutation | none today | `Mkdir` failure = "already reserved", `:284`; never removed or reused (`:281-282`) |
| 8 | `reservation_recovery.go:47` | `:36` `dir.Mkdir("reservation-intents")` | B | **Reachable** | intent-dir entry durable | **Refuse**, pre-mutation | `openIntentRoot(create=true)` fails → `PrepareCycleReservation` fails → **budgeted cycles routed through the trajectory Observer refuse on Windows** (`cycle_binding.go:269`) | nothing published; `Mkdir` idempotent |
| 9 | `reservation_recovery.go:67` via `:84` (two dir handles: `reservation-intents` + parent) | `:75` `Root.OpenFile O_EXCL` | A | **Reachable** | create entry durable | `O_SYNC` as #3; if the probe shows the flag rejected or the volume is not NTFS → **refuse**, pre-mutation | if refused: precharge intent publication refuses on Windows | `:79` "partial publication stays visible and cannot be overwritten by retry" |
| 10 | `reservation_recovery.go:67` via `:420` | none — re-sync | D | **Reachable** | inherits #9 | derives from #9: no-op iff #9's mechanism holds; else **refuse** | tracks #9 | re-sync only; no new artifact |
| 11 | `reservation_recovery.go:431` | none — re-sync of `trajectory.json` | D | **Reachable** | inherits #9 | as #10 | tracks #9 | as #10 |
| 12 | `parent_recovery.go:308` via `:336` (← func value at `:311`) | `:333` `dir.Rename(stage, parentRecoveryName)` | C | **Reachable** | rename entry durable | **Refuse**, and it **must be ordered before `:333`** (New finding 3) | `parley` parent-recovery **apply** refuses on Windows | pre-rename refusal ⇒ nothing published; `defer dir.Remove(stage)` `:328` clears the stage |
| 13 | `parent_recovery.go:308` via `:368` | none — re-sync on replay | D | **Reachable** | inherits #12 | **Refuse** (what it re-completes has no mechanism) | tracks #12 | re-sync only |
| 14 | `parent_recovery.go:545` via `:573` | `:570` `dir.Rename(stage, recoveredParentName)` | C | **Reachable** | rename entry durable | **Refuse**, ordered before `:570` | `parley` recovered-parent **apply** refuses on Windows | as #12; stage cleared at `:565` |
| 15 | `parent_recovery.go:545` via `:605` | none — re-sync on replay | D | **Reachable** | inherits #14 | **Refuse** | tracks #14 | re-sync only |
| 16 | `unchanged.go:279` via `:241` | `writeState` → `fsutil.ReplaceSyncedFile` (`state.go:197-221`) | D | **Reachable** | already satisfied | **Derived no-op on Windows**: publication is `MoveFileEx(MOVEFILE_REPLACE_EXISTING\|MOVEFILE_WRITE_THROUGH)` (`replace_windows.go:26`), so the barrier has nothing left to complete | none | `os.CreateTemp` + `defer os.Remove` (`state.go:205/:209`); replace is atomic |

Windows-reachable: **9 of 16** (#8–#16). Refusals under this selection: **#8, #12, #13, #14, #15**
reachable, plus #1/#2/#4/#6/#7 inherited if their gates lift. Derived no-op: **#16 only**, argued from
a proved mechanism on the same operation, never asserted.

---

## Restructure or refuse — the selection, made here and not left to implementation

**Class A — implement, do not restructure.** `syscall.O_SYNC` set from Windows-confined code (the
constant is `0x01000` on Windows, `0x80` on darwin). Documented for the write; **disclosed inference**
for the entry, never upgraded. Not restructured, because kimi-1's anti-replay objection is correct.

**Class B — refuse.** `Mkdir` takes no flags (`os.Root.Mkdir(name, perm)`, `os/root.go:149`), there is
no handle to carry write-through, and `FlushFileBuffers` is precluded by its own documented
precondition. A staged-directory `MoveFileEx` is *documented as available* (`N6`) and would even
preserve anti-replay, since `MOVEFILE_REPLACE_EXISTING` errors on an existing directory — but it is
path-based, so at #6/#7/#8 it gives up `os.Root` containment on the verification store and the budget
store's parent. **I select refusal**, explicitly and as a selection: a named refusal costs one feature
on one OS; the restructure costs a containment guarantee on reservation-critical code, in an idea whose
owner scope names security and privacy surfaces as the reason it is a deliberation.

**Class C — refuse, ordered before the rename.** Same reasoning, and stronger: converting means trading
a kernel-enforced reparse-resistant rooted rename for a lexical sibling check (verdict above). The
ordering requirement is New finding 3.

**Class D — derive per site, argued not asserted.** #16 qualifies now. #10/#11 qualify iff #9's
mechanism holds. #13/#15 do not qualify and refuse.

**No fourth branch.** No implicit no-op, no pre-approved weaker guarantee, no green API call treated as
a guarantee, and a deviation only as a reviewable artifact backed by concrete hosted evidence, before
any owner question. `FlushFileBuffers` on a directory handle remains a **diagnostic only** and may
never be acceptance evidence.

---

## Rooted containment and anti-replay

**Containment is preserved at every site under this selection — that is the reason for the selection.**
Nothing converts away from `os.Root`, so the guarantee at `os/root.go:34-43`, enforced by
`OBJ_DONT_REPARSE` (`at_windows.go:109-110`) and handle-relative `Renameat` (`:364-376`), is intact.
FINAL must state, in one sentence, that a lexical `filepath.Dir` sibling comparison is **not**
equivalent to that guarantee and may not be offered as a substitute.

**Anti-replay is preserved at every site.** #7's `base.Mkdir(charge)` remains the reservation token
with "never remove or reuse" intact (`:281-284`); #3/#5/#9's `O_EXCL` remains the create token with the
torn-artifact-stays-visible property intact (`:192-193`, `:79`); #12/#14 keep rooted-rename semantics
because they are not converted, and the pre-mutation refusal ordering is what stops the refusal
mechanism from itself creating an unbarriered publication.

---

## Census baseline

- **Skips, pinned by re-execution at `6b87cf8`:** 104 lines match `t\.Skip` under `internal` (and
  repo-wide — all skips live under `internal/`); **103 call sites** = 68 `t.Skip(` + 35 `t.Skipf(` + 0
  `t.SkipNow(`; the one non-call is `internal/runner/launch_precheck_skip_test.go:58`. **55 files with
  call sites; 56 files match the pattern.** FINAL states the denominator it means. The unescaped `107`
  is dead.
- **Emission:** `-json` on **all three legs** (kimi-1's reconciliation, adopted over the Windows-only
  variant). Windows leg: every firing SKIP must match an individually reviewed exclusion row or the leg
  fails; AF_UNIX fallbacks become hard failures rather than acquiring rows. macOS/Ubuntu: current
  firing skips are the enumerated baseline; any newly firing skip beyond it fails the leg. The workflow
  change is its own reviewed change; the unfiltered matrix is untouched.
- **Barrier census:** FINAL carries the mechanically derived table, not a headline count, with the
  derivation command and the function-value caveat, and the `fsutil.SyncDir` type contract as the audit
  that the compiler enforces.
- **New baseline row — the durability barrier's Windows behaviour is UNOBSERVED** in both hosted runs
  (New finding 1). FINAL records it as a known gap, and the probe bundle adds: `SyncFile` on a
  read-only directory handle with the raw error printed; the `O_SYNC` `O_EXCL` open in both plain and
  `os.Root` shapes; `runtime.Version()`; and the volume filesystem assertion (kimi-1's New concern 2),
  which is required because the only documented text for class A is NTFS-scoped.
- **Coverage language, unchanged:** hosted x64 only; `core.longpaths true` set by us
  (`tests.yml:38`), not an OS default; ARM64 assets build-only with zero executed evidence; no
  local-native Windows claim by anyone on this run; and this workspace itself sits on a cross-machine
  shared volume whose Windows semantics are unvalidated.

---

## Consensus readiness

**Consensus-ready on the mechanism dispute. It is resolved, and it resolved substantially my peers'
way.** `N5` is withdrawn; the class-A tier is kimi-1's and zcode-1's; the class-A restructure I floated
is withdrawn on kimi-1's anti-replay objection; zcode-1's `SyncDir` fail-closed construction is adopted
whole; the census is reconciled including the 55/56 split; the `-json` emission question is settled on
kimi-1's variant; and my round-04 §15.3 concern is discharged by kimi-1's own self-corrections. The
mechanism space is closed: class A implements, classes B and C refuse, class D derives per site.

**Two blockers remain. Both are documentation-and-decision, neither asks for a different mechanism, and
I judge both to be a few paragraphs of FINAL rather than a redesign.**

**Blocker 1 — FINAL must enumerate the operations it refuses and cost them.** Both round-04 contracts
resolve the five `Mkdir` operations to refusal by their own fail-closed rules while enumerating neither
them nor the blast radius. A reader of FINAL would not learn that `openIntentRoot(create=true)` refuses
on Windows, and therefore that budgeted cycles routed through the trajectory Observer refuse
(`cycle_binding.go:269`). The owner authorized "a real Windows implementation **or** an explicit,
reviewed, user-visible refusal" — explicit and user-visible is the part at issue. FINAL carries the
mechanically derived table with a per-operation disposition and named blast radius. If FINAL ships a
preference order over an un-enumerated site list, I will file it as a **MAJOR** at review.

**Blocker 2 — the refusal must be ordered before the mutation at #12 and #14, and zcode-1's containment
sentence must be struck or corrected.** As written, a Windows refusal at `parent_recovery.go:336`/`:573`
fires after `:333`/`:570` have renamed, publishing an unbarriered artifact and returning an error. And
the sentence "the conversion does not trade away `os.Root` containment" is `WRONG` on the sources cited
above; FINAL may record the trade as a rejected option, but not as a non-trade.

**Disputed claims and what acceptance depends on.** Acceptance does **not** depend on any disputed
claim. `N6`, `N7` and `N8` are owner-asserted and lack a non-owner verdict, because both peers wrote
round-04 concurrently and this is the last authorized round; each has a one-command check attached and
each belongs in Phase 6 review, where kimi-1 and I are the reviewers. `N7` strengthens the refusal case
rather than the acceptance case, so labelling it rather than relying on it costs nothing. I am not
asking for a sixth round: I am asking that these be labelled `UNVERIFIED` rather than absorbed, exactly
as in my round-04 New concern 4.

**Conditional acceptance, final form.** C2, C3 and C5 stand unchanged from round 3. C1′ from round 4 is
**superseded**: `O_SYNC` is accepted for class A as a documented mechanism for the write plus a
disclosed inference for the entry, and rejected as a mechanism for classes B, C and D, which resolve as
tabled. C4′ is discharged — both peers confirmed N1 and N3 at go1.26.8. Added: **C6** — no hosted probe
result may upgrade the create-entry inference to "documented", and a green `SyncDir` call on Windows may
never be read as a barrier.

**Release gates, re-checked read-only.** Predecessor handoff 1 present
(`source-context/release-1.49.1-done.md`, 1.49.1 at `54e0798`). Designated-implementer handoff still
**absent**. The gate is **closed**; version selection stays anchored above 1.49.1 and is made after both
handoffs; the Windows experimental label and the held winget PR stay until the released commit's
Windows leg is green; no tag moves. `zcode-1` remains drafter and single implementation owner with
kimi-1 and me as independent reviewers.

---

## Scope compliance

Read-only verification only: repo source at `6b87cf8`, local GOROOT `go1.27.1`, five go1.26.8 source
files fetched from `go.googlesource.com`, three Microsoft Learn pages fetched over `curl` and parsed
locally, `go doc`, one throwaway Go program executed under `/tmp` outside the repository, a
`GOOS=windows go vet` of that program, and `git` metadata. Wrote exactly one file:
`parley-deck/ideas/windows-portability/round-05/claude-1.md`. No product code, no workflow, no peer
artifact, no signoff, no phase transition, no consensus or FINAL, no commit, no tag, no publication, no
hosted run, no hosted mutation, no driver or second organizer launched, no additional round requested.
No secrets. No Chrome. English only.
