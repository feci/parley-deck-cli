---
agent: claude-1
idea: windows-portability
round: 4
date: 2026-09-25
responding-to: [kimi-1/round-03, zcode-1/round-03, round-04-brief]
---

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet sha256 re-verified by `shasum -a 256` this session; it matches the attestation. §15 re-read
for this round (15.1 ownership, 15.2 provenance, 15.3 conflicting verdicts). Read this round: all
three `round-03` artifacts in full, `round-04-brief.md`, and `00-prompt.md` including the
**2026-09-25 owner-default update** — which exempts this in-flight run, so facilitator stays
`codex-1`, quorum stays `claude-1`/`kimi-1`/`zcode-1`, launch models unchanged, and zcode-1 remains
drafter/implementer with kimi-1 and me reviewing. Nothing in that update changes my role or my
positions, and I apply no new-run default here.

HEAD moved twice while I worked: `2b1fc8c` → `ac1f0f7` → `2659520`.
`git diff --name-only 2b1fc8c..ac1f0f7 -- '*.go'` = **0 files** and
`git diff --name-only ac1f0f7..2659520 -- '*.go'` = **0 files** (`PRIMARY`); every changed file in
both steps is a `parley-deck/ideas/windows-portability/…` artifact, brief or organizer note. The
only `00-prompt.md` change at `2659520` records that shared memory became callable during round 4 and
returned nothing current — nothing that bears on any position here. Every source locator from rounds
1–3, mine and both peers', transfers unchanged; I spot-checked `verification.go:183-190` and
`gate.go:105` at `2659520` directly. Body verification below was executed at `ac1f0f7`.

Environment: darwin/arm64, local GOROOT `go1.27.1`. **New this round, and it is the substantive
addition:** I fetched the **CI-toolchain go1.26.8 sources** myself from `go.googlesource.com`
(tag `go1.26.8`, `?format=TEXT`, base64-decoded) and re-read the `MoveFileEx` and `CreateFileW`
Microsoft Learn pages over `curl`. That discharges half of my own round-03 C4 and, as set out below,
**refutes the load-bearing half of my own N1**. No browser, no Chrome. No code, workflow, peer
artifact, commit, tag, publication, hosted run or driver touched. Wrote exactly this file.

**Concurrent peer round-04 file, disclosed rather than left to inference.**
`round-04/zcode-1.md` appeared on disk at 00:20 local while this file was being drafted; I saved at
00:24 and **did not open it**. Nothing here responds to it, and any position of zcode-1's that this
file argues against may already have moved in that file. As in rounds 2 and 3 I state this rather
than the stronger and false "it did not exist yet". My cross-review target this round was the three
final `round-03` artifacts, as the brief instructed. `round-04/kimi-1.md` was absent when I saved.

**Ownership note under §15.1.** I own `N1`–`N4`, so I issue **no verdict** on them; where they move
below it is by `SELF-CORRECTION`, and a weakening takes effect immediately. The verdicts I issue are
on **peers'** claims. Two claims I introduce for the first time this round (`N5`, `N6`) are flagged
as owner-asserted and **non-owner-verdict-pending**, with the consequence stated in New concern 4.

---

## Position changes since round 3

**1. `SELF-CORRECTION`, and it is the largest one I have filed: my round-03 position change 1c
over-read the `CreateFileW` documentation, and `O_SYNC` does NOT give create-entry durability.**

My round-03 said "there is a documented real Windows mechanism; implement it", resting on one
sentence, and I labelled the create-entry half a "documented-mechanism inference". Having re-read
the whole *Caching Behavior* section rather than the one sentence I had already selected, that
inference is not merely unproved — **the same page documents against it** (`PRIMARY`,
learn.microsoft.com `CreateFileW`, fetched this round):

> "Also, the file metadata may still be cached (**for example, when creating an empty file**).
> **To ensure that the metadata is flushed to disk, use the FlushFileBuffers function.**"

Microsoft's own named remedy for *creation* metadata is `FlushFileBuffers` — the exact call whose
directory-handle applicability is undocumented (my round-03 position change 2, unchanged). Read with
the sentence I quoted last round — "flush any metadata changes, such as a time stamp update or a
rename operation, **that result from processing the request**" — the honest reading is:
`FILE_FLAG_WRITE_THROUGH` flushes metadata *incidental to the write* (timestamps; a rename performed
as part of the operation). It does **not** document a flush of the parent directory entry created by
the *open*.

Two honest limits on my own refutation, stated rather than buried: the "may still be cached" sentence
sits in the paragraph about `FILE_FLAG_NO_BUFFERING` combined with `FILE_FLAG_OVERLAPPED`, so it is
scoped there and is **not** a flat disproof of write-through-only creates; and Go's `O_SYNC` sets
`FILE_FLAG_WRITE_THROUGH` alone, never `NO_BUFFERING`, so that paragraph's exact configuration is not
ours. What it does establish is that Microsoft treats creation metadata as a separately-cached thing
and prescribes a different call for it. Under the brief's bar — "documentation scoped to the exact
operation" — **the create-entry claim fails**. I withdraw "implement it" and return to the brief's
default for that class.

**What survives of `N1`, and it is worth keeping.** The stdlib mapping itself is real and now
confirmed at the **CI toolchain**, which round-03 left open as my own C4 (`PRIMARY`, go1.26.8 fetched
this round; identical line numbers to go1.27.1):
- `syscall/syscall_windows.go:426-427` — `if flag&O_SYNC != 0 { attrs |= _FILE_FLAG_WRITE_THROUGH }`
- `syscall/types_windows.go:52` — `O_SYNC = 0x01000`
- `internal/syscall/windows/at_windows.go:99-100` — `if flag&syscall.O_SYNC != 0 { options |= FILE_WRITE_THROUGH }`
- `os/root.go:117-122` — `Root.OpenFile` validates only `perm&0o777` and passes `flag` through to
  `rootOpenFileNolog` (`os/root_windows.go:132-138`) → `openat`. So the rooted sites reach it too.

So `O_SYNC` is the right mechanism for **write** durability and for metadata a write carries. It is
**not** a directory-entry mechanism. My round-03 conflated those, and the conflation was in my favour.

**2. `SELF-CORRECTION`: my six-site census counted the wrong unit, and one locator was wrong.**
The exhaustive `grep -rn 'SyncFile(' --include='*.go'` (`PRIMARY`) gives **6** `SyncFile(<directory
handle>)` call sites reachable on Windows — that number holds — plus `internal/budget/lock_unix.go:31`
behind `//go:build !windows` (line 1, verified), unreachable. But:
- my round-03 wrote `reservation_recovery.go:64`; the directory sync is at **`:67`**
  (`f, err = dir.Open("reservation-intents")` at `:63`, `SyncFile(f)` at `:67`). Off by three.
- more importantly, **6 call sites serve 13 distinct barrier operations**, because
  `syncVerificationDirectory` (`internal/trajectory/verification.go:183-190`) has **6** callers
  (`verification.go:211, :267, :286`; `reservation_recovery.go:47, :67, :431`) and
  `syncTrajectoryRuntimeParent` has **3** (`trajectory_verify.go:125, :154, :193`). A per-site audit
  keyed on call sites therefore audits 6 things and silently decides 13.

**3. Consequence, and this is the round's material finding: the mechanism space all three of us
assumed does not cover the operations that exist.** All three round-03 files treat the space as
{write-through publication, directory-flush spike, refusal}. Sorting the 13 operations by *structural
shape* (`PRIMARY`, every line read this round) gives four classes, and write-through is attachable in
one of them:

| Class | Operations | Publication shape | Is write-through attachable? |
|---|---|---|---|
| **A** create-publication | `trajectory_verify.go:141`→barrier `:138`; `verification.go:202`→`:189`; `reservation_recovery.go:75`→`:67`+`:189` | `O_WRONLY\|O_CREATE\|O_EXCL` file create | `O_SYNC` sets on the open, **but** per position change 1 it does not document the parent entry |
| **B** `Mkdir` barrier | `trajectory_verify.go:190`→`:193`; `verification.go:264`→`:267`; `verification.go:283`→`:286`; `reservation_recovery.go:36`→`:47` | `Mkdir`. **No file open, no flags, no rename** | **No.** Structurally impossible — `os.Root.Mkdir`/`os.Mkdir` take no flag argument |
| **C** rooted stage-rename | `parent_recovery.go:333`→`:308`; `parent_recovery.go:570`→`:545` | `dir.Rename(stage, final)` | **No,** not as written — see position change 4 |
| **D** recovery re-sync | `unchanged.go:241`→`:279`; `parent_recovery.go:368`→`:308`; `parent_recovery.go:605`→`:545`; `reservation_recovery.go:431`→`:67`+`:189` | none — reopen an existing file read-only and re-establish a barrier | **No.** There is no publication to attach to |

Class B is the one nobody costed. `verification.go:283`'s `base.Mkdir(charge, 0700)` carries the
comment (`PRIMARY`, `:281-282`) "*Never remove or reuse an existing reservation, even if request
writing was interrupted. Explicit recovery must account for that history.*" The durability of **that
directory's own entry** is the anti-replay invariant. A blocking refusal at class B is therefore not
"one publication path refuses" — it refuses `openVerificationStorage(create=true)` and
`openIntentRoot(create=true)`, i.e. **trajectory verification and budget reservation cannot be
created on Windows at all**. That is authorized by the owner's "real implementation *or* an explicit,
reviewed, user-visible refusal", but it is a feature-level refusal and no round-03 file says so.

**4. `SELF-CORRECTION`: I claimed the write-through rename route was "already in-tree" for the rename
sites. It is in-tree, but the two rename sites cannot reach it without giving up `os.Root`.**
`MOVEFILE_WRITE_THROUGH` is a flag of path-based Win32 `MoveFileEx`, used at
`internal/fsutil/replace_windows.go:26`. Class C publishes with `os.Root.Rename`
(`parent_recovery.go:333`, `:570`), which on Windows is `renameat` →
`SetFileInformationByHandle(FileRenameInformation[Ex])` (`os/root_windows.go:386`;
`internal/syscall/windows/at_windows.go:363-429`, go1.26.8) — an NT-level rename with **no
write-through flag available**. Round-03 I flagged the `os.Root`-containment tension for the `O_SYNC`
route and said `O_SYNC` avoids it; I did not flag that the *rename* route has the opposite problem.
Class D1 is the exception that proves it: `unchanged.go`'s publication is `writeState`
(`state.go:205-221`) which does `CreateTemp` → write → fsync → close → `fsutil.ReplaceSyncedFile`, so
it **is** path-based write-through, and its `:279` barrier is genuinely belt-and-braces. That is one
of 13.

**5. Everything I settled in round 3 outside durability stands, and I reopen none of it.** Owner-only
protected DACL; P-A with P-B deferred; blocking missing-`sh` refusal with no owner question; universal
gate-name encoding with cosmetic-only grandfathering; the reviewed Windows skip census; the ACP
two-test split; the CRLF document-semantic/evidence-integrity split; the label gate exactly as the
owner set it with my extra conditions demoted to FINAL documentation; `zcode-1` as drafter and single
implementation owner, ratified.

---

## Responses to others

### @kimi-1 — round-03

**Accepted without reservation:** your owner-only ACL withdrawal (position change 1) — that closes the
allow-set with all three of us converged, not by majority; your withdrawal of the documented-nil floor
(position change 2), which was zcode-1's single stated ratification condition and mine, so **that
dispute is closed**; your zcode-1 ratification; your no-exclusion-row AF_UNIX position; your ACP split
adoption; your census mechanics sentence (New concern 1) and your pattern-not-symptom sweep framing
(New concern 2); H1–H7 as written, with the amendments below folded in.

**V1 — one half is a verdict I must record as `WRONG`, with provenance.** Your V1 endorses my round-2
masking chain ("*so the CreateFile path ignores the bit*") and then draws your own conclusion, which
you own: "*the source narrows the hosted failure to **exactly one path**: the `ERROR_FILE_NOT_FOUND`-
with-restart-class branch that `break`s with no error (the MS-FSA empty-root quirk), yielding
(empty, nil)*". I own the masking claim and issue no verdict on it (I withdrew it myself as `N3`). The
narrowing conclusion is yours, and my verdict is **`WRONG` (`PRIMARY`, go1.26.8 fetched this round)**:

- `syscall/syscall_windows.go:455-466` handles `o_DIRECTORY` **after** `CreateFile` —
  `GetFileInformationByHandle`, and `if fi.FileAttributes&FILE_ATTRIBUTE_DIRECTORY == 0 { CloseHandle(h);
  return InvalidHandle, ENOTDIR }` at `:462-464`. The masking is real and irrelevant: the bit is
  consumed by a post-open guard.
- Therefore the open **fails** for a regular file, and `os/dir_windows.go`'s readdir loop is never
  entered. The `ERROR_FILE_NOT_FOUND` restart-class `break` you point at is at
  `os/dir_windows.go:134-146` — **inside that loop**. It is unreachable on this input, so it cannot be
  the path the failure narrows to.
- `ENOTDIR` is absent from the `ErrNotExist` set (`syscall_windows.go:201-205`: `ERROR_FILE_NOT_FOUND`,
  `_ERROR_BAD_NETPATH`, `ERROR_PATH_NOT_FOUND`, `ENOENT`), so the product vetoes and the test should
  **pass**.

Your V1 and zcode-1's round-03 verdict are directly contradictory on the same identified claim
(`CONFIRMED` vs `REFUTED at source level`). Under §15.3 that is a verdict conflict; see New concern 1.

**V2 — `four directory-handle sync sites: CONFIRMED`: superseded.** The underlying "four" was my claim
and I withdrew it by `SELF-CORRECTION` in round-03, so I issue no verdict; a weakening takes effect
immediately, which means your `CONFIRMED` now certifies a withdrawn statement and should be struck
rather than left standing. The mechanical result is in position change 2 above. Separately I verdict
zcode-1's independent claim that four is an undercount as **`CONFIRMED`** (below).

**V4 — the skip count is `WRONG`, and the census's denominator is an acceptance artifact, so it
matters.** Verdict **`WRONG` (`PRIMARY`, executed this round at `ac1f0f7`)**. I ran your exact command
and then the escaped form:

```
grep -rn 't.Skip'  --include='*_test.go' internal | wc -l   → 107
grep -rn 't\.Skip' --include='*_test.go' internal | wc -l   → 104
grep -rn 't\.Skip' --include='*_test.go' .        | wc -l   → 104   (all skips live under internal/)
```

The unescaped `.` is a wildcard: the three extra lines are `internal/tui/live_test.go:44`, `:45` and
`internal/runstate/segment_test.go:145`, which match `stateSkipped`/`StateSkipped`, not a skip call. Of
the 104, exactly one is not a call — `internal/runner/launch_precheck_skip_test.go:58`
(`if !result.Skipped`). Breakdown: `t.Skip(` = **68**, `t.Skipf(` = **35**, `t.SkipNow(` = **0**; 68+35
= **103 skip call sites**. Your other V4 premises are **`CONFIRMED` (`PRIMARY`)**: `.github/workflows/
tests.yml:52` is `go test ./... -count=1 -timeout 45m` with no `-v` and no `-json`, and
`git ls-files | grep -c 'gate\.json'` = **0**. Your conclusion is untouched; only the number changes.

**V3, V5, V6, V7, V8 — no dispute; I re-executed V7 and V8 and they hold** (`GatePath` at
`internal/pipeline/gate.go:105`, `tmp := path + ".tmp"` at `:119`, `BlockWorkspace` at
`internal/pipeline/executor.go:30`; 0 `.go` files `2b1fc8c..ac1f0f7`).

**Your H1 needs one amendment, and it is the round's substance.** H1 reads "*if the spike succeeds
**and** authoritative documentation scoped to directory-handle flush exists, implement it; if either
fails, the default is the reviewed per-site refusal*". The ordering is right and I keep it. What H1
cannot see — because you wrote before my round-03 — is that after position changes 1 and 3 above,
**no documented mechanism remains for classes A, B or D**, so H1's "either fails" branch is not a
residual: it is the default for **11 of 13** operations. That is a much larger refusal surface than
your "refusing the whole write path would be disproportionate" caution was aimed at, and I think you
were right to raise that caution. My answer is in Current proposal 3: restructure, rather than either
refuse a feature or weaken a guarantee.

### @zcode-1 — round-03

**Accepted without reservation:** `Drafter: yes` and the single implementation claim recorded in
`IMPLEMENTATION.md` before any code edit; position changes 1 and 4; the mechanical-not-hand-list audit
demand; the per-site disposition table in FINAL (New concern 2); your shell-refusal wording template
(New concern 3) — I have no improvement on it; refuse-and-instruct for existing stores rather than
in-place DACL rewrite (New concern 4) and the actionability note (New concern 5); the census on all
three legs with the Windows leg reviewed row-by-row; `os.Root`-where-already-rooted over raw
`CreateFile`; H1–H7.

**Your go1.26.8 `o_DIRECTORY` refutation: `CONFIRMED` (`PRIMARY`, independently fetched).** I fetched
tag `go1.26.8` myself rather than relying on your fetch, and every line you cite is there at the line
numbers you give: the post-open guard at `syscall/syscall_windows.go:455-466` (`ENOTDIR` at `:464`),
the `ErrNotExist` set at `:201-205`, `openDirNolog` passing `O_RDONLY|windows.O_DIRECTORY` at
`os/file_windows.go:166-167`, `fileFlagsMask = 0xFFF00000` at `syscall/types_windows.go:95`. Your
split verdict — **phenomenon real, mechanism refuted, cause undiagnosed, fix ships regardless** — is
correct, and it is now confirmed at the toolchain CI actually resolves. This discharges the `N3` half
of my round-03 C4.

**Your `MoveFileEx` quote: `CONFIRMED` verbatim (`PRIMARY`, same page fetched this round).** "*The
function does not return until the file is actually moved on the disk. Setting this value guarantees
that a move performed as a copy and delete operation is flushed to disk before the function returns.
The flush occurs at the end of the copy operation.*" Your reading is right: the explicit *guarantees…
flushed* clause is scoped to copy-and-delete, and a same-volume rename is covered only by the weaker
first sentence. Your amendment 2 — that the write-through route meets the *same* documentation bar as
the spike and is "the preferred route, not an automatic proof" — is correct and I adopt it against my
own round-03 phrasing.

**Your undercount amendment: `CONFIRMED`, and your unit is closer to right than mine was
(`PRIMARY`).** Your extra sites are real (`unchanged.go:268-280`, `reservation_recovery.go:55-68`,
reuse at `:431`). Two refinements, offered as precision rather than correction: the directory sync in
`syncIntent` is at **`:67`**, not within `:55-68` generally, and `:431` is not a reuse of the same
*site* but a fourth caller of `syncVerificationDirectory`. Your "~7 directory-sync call points across
5 files" and my "6 sites" are both counting syntax; the number that governs the contract is **13
barrier operations in 4 structural classes** (position change 3). Your instinct that hand-lists have
"been wrong once" is right — they have now been wrong three times, including twice by me.

**Where I now disagree with your Current proposal 3, having agreed with its structure.** You write the
preference order as "write-through-publication → spike (documentation-based bar) → reviewed
user-visible refusal". After this round that order resolves to **refusal for 11 of 13 operations**,
because write-through-publication is structurally unavailable for classes B and D and undocumented for
the parent entry in class A, and the spike is undocumented for directory handles by your own bar. Your
readiness claim — "*no unresolved design decision remains*" — is the one sentence in your file I
cannot co-sign. Choosing between (i) a feature-level Windows refusal for trajectory verification and
budget reservation and (ii) restructuring class A/B/C publication to a write-through rename **is** a
design decision, it lands on anti-replay-critical code, and no probe result selects it. That is my
single remaining blocker and it is narrow; Current proposal 3 states the shape I think FINAL needs.

**One agreement worth making explicit, because it protects your Stage 3.** Your New concern 1 says
that if the probe shows `ReadDir` returning nil+empty on the runner, that "*would itself implicate
every `ReadDir`-based decision on Windows*". Agreed, and I add a cheaper first step in New concern 2:
the discrepancy may be a floating toolchain, which existing logs can settle without a new run.

---

## New concerns / questions

1. **A §15.3 verdict conflict exists now and the drafter must carry it.** kimi-1's V1 assigns
   `CONFIRMED` to the `o_DIRECTORY`-masking ⇒ open-succeeds ⇒ (empty, nil) chain; zcode-1's round-03
   assigns `REFUTED at source level` on both go1.27.1 and go1.26.8; I have independently confirmed
   zcode-1's reading at go1.26.8 and have withdrawn the claim as its owner. Under §15.3 contradictory
   verdicts existing when consensus opens require a `## Verdict conflicts` section in `consensus.md`
   quoting each verdict and its provenance. The cleanest resolution is kimi-1 striking that half of V1;
   absent that, zcode-1 must carry it as a conflict rather than smooth it. Same for V2's `CONFIRMED` on
   a withdrawn claim and V4's 107.
2. **The CI toolchain is not pinned, and that is a cheap, checkable candidate for the `strict_gate`
   discrepancy.** `go.mod:3` is `go 1.26` with **no `toolchain` directive** (grep: none), and
   `.github/workflows/tests.yml:42` is `go-version-file: go.mod` (`PRIMARY`). `actions/setup-go`
   resolves that to the latest `1.26.x` **at the moment the workflow runs** — currently `go1.26.8`
   (`https://go.dev/dl/?mode=json&include=all` lists 1.26.4 … 1.26.8). So the run that produced the
   `:179` failure (36009912946) may have resolved a different patch than the sources zcode-1 and I
   read. Before any new hosted run: **read the resolved version out of the existing setup-go log for
   that run**. If it is not 1.26.8, the source reading may simply describe different code; if it is,
   the divergence is on the runner and zcode-1's escalation in its New concern 1 is warranted. Either
   way the probe must print `runtime.Version()`, as zcode-1 required. I state plainly that **no branch
   is confirmed** and the fix is branch-independent.
3. **`N4`, re-executed, is sharper than I stated and the sharper form should be the FINAL sentence.**
   My round-03 said "four of eight cases show that asymmetry". Precisely (`PRIMARY`, executed this
   round, `filepath.IsLocal` at go1.27.1/darwin): of the five traversal IDs tested, **four flip from
   correctly-rejected to wrongly-accepted** when `IsLocal` is applied to the concatenated element
   `slug__+id` instead of the raw ID — `..`, `../x`, `../../x`, `../../other-idea` all give
   `IsLocal(raw)=false`, `IsLocal("my-idea__"+id)=true`. Only `../../../../x` stays `false`. And
   `../../other-idea` is the worst of them: `BlockWorkspace` lands at `/deck/ideas/other-idea`,
   **another idea's directory**. Separately, `a/b`, `a\b`, `..\..\x`, `C:evil`, `C:\evil`, `a:b`, `CON`
   and `end.` are **all** `IsLocal=true` on darwin, so the predicate does not catch the unsafe class on
   Unix at all. Conclusion unchanged and now better supported: validate the **raw ID** with an explicit
   **charset allowlist**, never a path-shape predicate. `internal/filepathlite/path_windows.go`'s own
   `isLocal` (go1.26.8, read this round) rejects *any* `:` — "*Rejecting any path with a colon is
   conservative but safe*" — and every reserved name; cite it as the allowlist's precedent, do not
   reuse it as the guard.
4. **`N5`/`N6` are mine, introduced in the last scheduled cross-review round, and cannot get a
   non-owner verdict inside Phase 2.** `N5` = the `CreateFileW` creation-metadata refutation (position
   change 1). `N6` = `MoveFileEx` is documented to move directories (position change 3 / Current
   proposal 3). `N5` weakens my own claim, so it takes effect immediately under §15.1. `N6` is a
   **strengthening** supporting a proposal, so it stays `UNVERIFIED` until kimi-1 or zcode-1 verdicts
   it. The deliberation-track cap is 3 cross-review rounds after round 1 — rounds 2, 3, 4 — so there is
   no round 5 to carry it. FINAL must therefore label `N6` owner-asserted/`UNVERIFIED`, and the verdict
   belongs in Phase 6 review, where kimi-1 and I are the reviewers. I am not asking for an extra round;
   I am asking that the gap be labelled rather than absorbed.
5. **Carried unanswered from my round-02 and round-03: pre-existing-store DACL recovery.** zcode-1's
   New concern 4 answers it (refuse-and-instruct, no in-place rewrite) and kimi-1's item 1 agrees. I
   record it as **settled** and withdraw it as a concern. Guard sites re-verified: `snapshot.go:162`
   (directory, `Perm()&0077 != 0`) and `:544` (regular file).

---

## Current proposal

Items 1, 2, 4, 5 and 6 are the merged tri-participant contract already settled in round 3 across all
three files; I restate them by reference rather than re-litigate them: **owner-only protected DACL**
with effective-grant verification, FAT/exFAT refusal, adversarial hosted tests, existing-store
refuse-and-instruct, Unix byte-identical; **P-A only** with truthful `Alive` failing closed on
`ERROR_ACCESS_DENIED`, P-B and `mvdan.cc/sh/v3` as named deferrals, and a **blocking** missing-`sh`
refusal at both entry points exercised hosted with `sh` off `PATH`; **universal gate-name encoding
inside `GatePath`** with `.tmp` inheriting by construction, cosmetic-only grandfathering, the unsafe
class refused at load and at both interpolation backstops on every OS, one-release legacy `->`
read-fallback with shadow retirement on rewrite, and the raw-ID charset allowlist per New concern 3;
the **branch-independent `strict_gate` fix** (stat the round path, veto when it exists and is not a
directory) plus the cross-platform pinning test plus the `ErrNotExist`-after-directory-read sweep, with
the cause recorded as **undiagnosed** and no branch called confirmed; and the **census/ledger/ACP-split/
CRLF-split** acceptance set with the corrected denominator of **103** skip call sites.

**3. Directory durability — the exact conditional contract, replacing my round-03 item 3.**

The contract is per **operation**, not per call site, and the complete inventory is the 13 operations
in the four-class table above. Introduce `fsutil.SyncDir` as a named per-OS contract so no directory
handle reaches `SyncFile`, and prove the conversion by grep at review. Then, per class:

- **Class D (4 operations) — derive, do not weaken.** These reopen an already-published file to
  re-establish a barrier a crash interrupted; they have no publication of their own. Where the original
  publication carries a documented write-through guarantee, the Windows barrier has nothing left to
  complete and `SyncDir` is correctly a no-op **at that site**. That is a *derivation from a proved
  mechanism on the same operation*, not the "documented no-op" the brief forbids, and the distinction
  must be argued per site in FINAL, not asserted. D1 already qualifies (`state.go:221` →
  `ReplaceSyncedFile`). D2–D4 qualify only if their publications (C1, C2, A3) are resolved below.
- **Classes A, B, C (9 operations) — restructure to a write-through rename, or refuse.** No documented
  mechanism exists for the parent entry as these are written today (position changes 1, 3, 4). The one
  documented directory-entry barrier Windows offers is `MoveFileEx(MOVEFILE_WRITE_THROUGH)`, and the
  strongest sentence that applies to a same-volume move is "*does not return until the file is actually
  moved on the disk*" — tier-2 documentation, weaker than the copy-and-delete flush clause, and already
  the reviewed in-tree rationale at `internal/fsutil/replace_windows.go:10-13`. Restructuring is
  **documented as available** (`N6`, `PRIMARY`, `MoveFileEx` page, fetched this round):
  - "*Moves an existing file **or directory**, including its children, with various move options*"; and
    "*When moving a directory, the destination must be on the same drive*" — so a stage-then-rename
    works for class B's directories, on one volume, which is the only case we need.
  - `MOVEFILE_REPLACE_EXISTING`: "*If lpNewFileName names an existing **directory**, an error is
    reported*" — so a rename into place **cannot** clobber an existing reservation. Class B's
    anti-replay invariant ("never remove or reuse an existing reservation") is **preserved**, which is
    the property I would have expected a restructure to break.
  - Cost, stated not hidden: `MoveFileEx` is path-based, so the restructured sites leave `os.Root`
    containment. `ReplaceSyncedFile` already carries the substitute guard — a sibling-directory check
    at `replace_windows.go:15-17` — and that is the pattern to follow, not a fresh invention.
- **The decision itself belongs in FINAL.** Per class, FINAL states one of: **restructure** (with the
  documentation quoted and the `os.Root` trade recorded), or **named blocking refusal** for that
  operation with its blast radius written out — and for class B the blast radius is that trajectory
  verification and budget reservation cannot be **created** on Windows. There is no third branch: no
  implicit no-op, no pre-approved weaker guarantee, and a deviation only as a reviewable artifact
  backed by concrete hosted evidence, before any owner question.
- **`FlushFileBuffers` on a directory handle** remains a **diagnostic only** and may never be cited as
  acceptance evidence (round-03 position change 2, unchanged, and now reinforced: the `CreateFileW`
  page names `FlushFileBuffers` as the remedy for creation metadata while the `fileapi.h`
  `FlushFileBuffers` page documents no directory subject at all).

**Staging.** zcode-1's Stage 0–7 numbering unchanged, with Stage 1 still preceded by the
`syscall.Mkfifo` build tag and the `app_test.go:155/:157` comma-ok fix, and **Stage 3 now opening with
the four-class table and its per-class disposition** before any code is written.

---

## What I accept, what remains blocking, and the exact conditional acceptance

**Consensus-ready: five of six items.** Privacy, process/shell, names, `strict_gate`, and
tests/skips/evidence are settled across all three round-03 files with no live disagreement, and the
last two policy disputes closed this round: kimi-1 withdrew the documented-nil floor (zcode-1's and my
stated condition) and all three of us are on owner-only ACLs. I ratify `zcode-1` as drafter and single
implementation owner, with kimi-1 and me as independent reviewers, and I reopen nothing.

**Remaining blocking counterproposal — one, narrowed and re-aimed.** Round 3 my blocker was "a
documented no-op may not be the default". That is now moot: kimi-1 withdrew it and both peers adopted
refusal-by-default. My blocker this round is different and is the consequence of my own correction:
**FINAL may not ship item 3 with a preference order whose unstated resolution is a feature-level
Windows refusal.** Concretely, FINAL must carry the **four-class, 13-operation table with an explicit
per-class disposition** — restructure-to-write-through-rename or named blocking refusal with blast
radius — decided in the design, not in implementation. If FINAL ships item 3 as
"write-through → spike → refusal" over an un-enumerated site list, I will file it as a **MAJOR** at
review, on the ground that the audit then decides 13 operations while appearing to decide 6, and that
a design decision of this size was delegated to the implementer. This is a documentation-and-decision
demand, not a demand for a different mechanism, and I believe it is a few paragraphs of FINAL rather
than a redesign. I also record that zcode-1's "no unresolved design decision remains" is the specific
sentence I am contesting, so the disagreement is locatable.

**Exact conditional acceptance for the unresolved hosted facts.** I accept the plan subject to these
and nothing else. C1 and C4 from round 3 are superseded; C2, C3, C5 stand.

- **C1′ (replaces C1).** `O_SYNC` write-through is accepted as the mechanism for **write** durability
  wherever a publication open exists, and is **rejected** as a parent-directory-entry mechanism — on
  documentation, not on a probe, so no hosted result reopens it. The **`O_SYNC` proposal from my
  round-03 is therefore: accepted in part, rejected for the purpose it was offered.** For classes A, B
  and C the choice is restructure-or-refuse per Current proposal 3, and a hosted probe may confirm a
  restructure works but may never be the thing that establishes the guarantee.
- **C2 (stands).** If the toolchain-pinned `ReadDir` probe shows an error not classified as
  `fs.ErrNotExist`, I withdraw the fail-open hypothesis entirely and `:179` returns to undiagnosed with
  the invariant fix landing anyway. If it shows nil+empty, or an error classified as `fs.ErrNotExist`,
  the fail-open is confirmed and it is a safety defect. It may not be filed as a fixture without this
  result. **Preceded by** the free check in New concern 2: read the resolved Go version out of the
  existing log for run 36009912946.
- **C3 (stands).** If hosted evidence shows effective-grant verification refusing the default runner
  temp location, FINAL carries the pre-existing-store answer before Stage 1 closes. zcode-1's
  refuse-and-instruct is that answer; the refusal is correct and is not weakened for the runner.
- **C4′ (replaces C4).** The CI-toolchain re-read is **done for `N1` and `N3`** — both at `go1.26.8`,
  by me this round, independently for `N3`. What replaces it: `N5` takes effect now as a weakening;
  **`N6` needs a non-owner verdict from kimi-1 or zcode-1** and, failing that, must reach FINAL labelled
  `UNVERIFIED` per New concern 4. `N2` and `N4` likewise still lack a non-owner verdict; the corrected
  forms are in position change 2 and New concern 3.
- **C5 (stands).** Coverage language in FINAL stays honest: hosted x64 only; `core.longpaths true` set
  by us (`tests.yml:38`) and not an OS default; ARM64 assets build-only with zero executed evidence and
  no ARM64 gate; no local-native Windows claim by anyone on this run; and the one sentence that this
  workspace itself sits on a cross-machine shared volume whose Windows semantics are unvalidated.

**Release gates, factual status, re-read this round.** Predecessor handoff 1 is present
(`source-context/release-1.49.1-done.md`, 1.49.1 at `54e0798`). The designated-implementer handoff is
still absent. The release gate therefore remains **closed**, version selection stays anchored above
1.49.1 and is made after both handoffs, the Windows experimental label and the held winget PR stay as
they are until the released commit's Windows leg is green, and no tag moves.

**Scope compliance.** Read-only verification only: repo source at `ac1f0f7` (locators spot-checked at `2659520`), local GOROOT `go1.27.1`,
go1.26.8 sources fetched from `go.googlesource.com`, two Microsoft Learn pages and the `go.dev`
release index over `curl`, one throwaway `filepath` program executed under `/tmp` outside the
repository, and `git` metadata. Wrote exactly this file:
`parley-deck/ideas/windows-portability/round-04/claude-1.md`. No product code, no workflow, no peer
artifact, no signoff, no phase transition, no commit, no tag, no publication, no hosted run, no
driver or second organizer launched. No secrets. No Chrome.
