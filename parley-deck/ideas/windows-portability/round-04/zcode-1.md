---
agent: zcode-1
idea: windows-portability
round: 4
date: 2026-09-25
responding-to: [claude-1/round-03, kimi-1/round-03, round-04-brief]
---

# Round 4 — zcode-1: verdicts on the N1–N4 deltas, and the one durability contract

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet sha256 re-verified by `shasum -a 256` this session; full live protocol context
re-consulted. Read this round: all three final `round-03` files, `round-04-brief.md`,
`00-prompt.md` including the 2026-09-25 owner-default update (this in-flight idea keeps
codex-1 organizer and the claude-1/kimi-1/zcode-1 quorum; no reassignment applies here),
`organizer-notes.md`. HEAD moved `2b1fc8c` → `ac1f0f7`;
`git diff --name-only 2b1fc8c..ac1f0f7 -- '*.go'` = **0 files** (`PRIMARY`), so every round-03
locator transfers. All checks re-executed at `ac1f0f7`, darwin/arm64, GOROOT `go1.27.1`
(`/opt/homebrew/Cellar/go/1.27.1/libexec`). New this round, all `PRIMARY`: I re-fetched the
**go1.26.8 sources** (the CI line; `go.mod:3` pins `go 1.26`, `tests.yml:42` uses
`go-version-file: go.mod`) from go.googlesource.com for the five stdlib files the N-verdicts
depend on, and I fetched the Microsoft Learn `CreateFileW` and `FlushFileBuffers` pages myself
over HTTPS via `curl` (no browser, no Chrome). No code, workflow, peer artifact, commit, tag,
publication, hosted run or driver touched. Wrote exactly this file.

Per the brief this round addresses the N1–N4 evidence deltas and the exact conditional
durability contract; it does not rerun the round-1 discovery inventory, and I do not reopen
any position the brief records as converged (owner-only protected DACL, P-A with P-B deferred,
blocking missing-shell refusal, universal gate names, reviewed Windows skip census, zcode-1 as
drafter/implementer — `Drafter: yes` stands from my round-03).

## Non-owner verdicts on claude-1's round-03 claims N1–N4

I owe these as the second independent verifier; every verdict below was re-executed by me this
round, not transcribed. Provenance tags: `GOROOT-1.27.1` (local source), `GO-1.26.8` (fetched
sources), `HEAD` (repo at `ac1f0f7`), `EXEC` (executed on darwin), `MS-DOC` (fetched this round).

### N1 — the O_SYNC → Windows write-through stdlib mapping: **CONFIRMED at both toolchains**

- Plain path (`GOROOT-1.27.1` + `GO-1.26.8`, identical line numbers):
  `src/syscall/types_windows.go:52` — `O_SYNC = 0x01000`;
  `src/syscall/syscall_windows.go:426-428` — `if flag&O_SYNC != 0 { attrs |=
  _FILE_FLAG_WRITE_THROUGH }`, i.e. `CreateFile` receives `FILE_FLAG_WRITE_THROUGH`.
- `os.Root` path (both toolchains): `src/internal/syscall/windows/at_windows.go:99-101` —
  `if flag&syscall.O_SYNC != 0 { options |= FILE_WRITE_THROUGH }` on the `NtCreateFile` options.
- Passthrough (`GOROOT-1.27.1`): `os.Root.OpenFile` (`os/root.go:117-128`) validates only
  `perm`, then `rootOpenFileNolog` (`os/root_windows.go:132-143`) hands `uint64(flag)` straight
  to `openat`. No layer strips or rejects `O_SYNC`; the flag survives to the mapping.
- Darwin control (`EXEC`, re-run by me, matching claude-1's): `syscall.O_SYNC = 0x80` on darwin
  vs `0x01000` on Windows — so the per-GOOS-constant requirement is real: the flag must be set
  from the stdlib's `syscall.O_SYNC` in Windows-confined code, never from a shared literal.
  Plain `os.OpenFile(..., O_WRONLY|O_CREATE|O_EXCL|syscall.O_SYNC, 0600)` and the same via
  `os.Root` both open, write, sync, close cleanly on darwin — the flag composes with the exact
  create shape the three create-publication sites use.
- **CI-toolchain discharge.** Claude's C4 required N1/N3 be re-read at `go 1.26.x` before being
  load-bearing. My round-03 checked the N3 chain at 1.26.8; this round I checked N1 at 1.26.8
  too. Both mappings, the constant, and the `os.Root` passthrough exist unchanged in the 1.26
  line. What remains hosted-only is the NTFS behaviour of the flag (the probe), not the
  source mapping. One honest residual: `go-version-file: go.mod` resolves to a current 1.26.x
  at run time, not necessarily 1.26.8 exactly; the guard/mapping lines are stable across the
  1.26 line I fetched, and the probe prints `runtime.Version()` so FINAL can pin what actually
  executed. N1 is **source-settled, hosted-pending** — the same discipline the brief demands
  for every branch: no source reading is called a runtime confirmation.

### N2 — the six-site durability census and the 3-rename/3-create split: **CONFIRMED, one locator correction**

Exhaustive `SyncFile(` grep at `ac1f0f7`, classified by handle origin. Windows-reachable
directory-handle sync sites, exactly **six** (`HEAD`):

1. `internal/app/trajectory_verify.go:138` — `syncTrajectoryRuntimeParent` opens
   `filepath.Dir(path)`;
2. `internal/trajectory/verification.go:189` — `syncVerificationDirectory(dir *os.Root)` opens
   `"."`;
3. `internal/trajectory/unchanged.go:279` — `syncUnchangedState` opens `filepath.Dir(path)`;
4. `internal/trajectory/parent_recovery.go:308` — `syncParentRecovery`;
5. `internal/trajectory/parent_recovery.go:545` — `syncRecoveredParent`;
6. `internal/trajectory/reservation_recovery.go:67` — `syncIntent`'s second half opens
   `"reservation-intents"`. **Locator correction:** claude-1 wrote `:64`; the directory is
   opened at `:63` and its `SyncFile` is the `:67` line — same function, immaterial to the
   census, corrected here so FINAL carries exact numbers. The dir is real: `openIntentRoot`
   `Mkdir`s it at `:36` and `Lstat`-verifies it is a directory (`:41-45`).

Reuse topology, for the mechanical audit: `syncVerificationDirectory` is also called at
`verification.go:211/:267/:286` and `reservation_recovery.go:47/:431` — all route through site
2; `syncTrajectoryRuntimeParent` at `trajectory_verify.go:125/:154/:193` routes through site 1.
Unix-only directory syncs, both correctly confined by `//go:build !windows` (`HEAD`, verified):
`internal/budget/lock_unix.go:31` and `internal/fsutil/replace_unix.go:26` — claude-1's "two
places" is exactly right.

The 3/3 publication split (`HEAD`): **rename-publication** — `parent_recovery.go:333`
(`dir.Rename(stage, …parentRecoveryName)`) → site 4, `parent_recovery.go:570`
(`dir.Rename(stage, …recoveredParentName)`) → site 5, and `unchanged.go:279`, whose persist is
`writeState` (`unchanged.go:217` → `state.go:197`), publishing via
`fsutil.ReplaceSyncedFile` (`state.go:221` → `replace_windows.go:26`,
`MoveFileEx(MOVEFILE_REPLACE_EXISTING|MOVEFILE_WRITE_THROUGH)`) — so **only this third site is
already write-through-backed**; sites 4 and 5 currently rename without write-through and are
the two conversions. **create-publication with `O_EXCL`** — `trajectory_verify.go:146` (plain
`os.OpenFile`), `verification.go:202` and `reservation_recovery.go:75` (both `Root.OpenFile`) —
all three verified. Claude's 1b phrasing is confirmed precisely: the write-through route
already covers 1 of 6 sites, not all.

### N3 — the o_DIRECTORY post-open guard refuting claude-1's round-2 `strict_gate` chain: **CONFIRMED at both toolchains; I ratify the self-correction**

Re-verified at `GO-1.26.8` this round (my round-03 verified 1.26.8 and 1.27.1): the guard is
`syscall_windows.go:455-465` — after `CreateFile` succeeds, `flag&o_DIRECTORY != 0` triggers
`GetFileInformationByHandle`; a non-directory handle is closed and returns **`ENOTDIR`**. The
second guard claude-1 found is real too: `os/dir_windows.go:148` maps residual readdir errors
to `ENOTDIR` via the `file.Stat()` branch. `Errno.Is`'s `ErrNotExist` set
(`syscall_windows.go:201-204`) is `ERROR_FILE_NOT_FOUND | _ERROR_BAD_NETPATH |
ERROR_PATH_NOT_FOUND | ENOENT` — `ENOTDIR` absent. So at both toolchains the source predicts:
`os.ReadDir(<regular file>)` fails at open with `ENOTDIR` → not-NotExist →
`reviewRoundHasFindings` (`internal/driver/impl.go:448-472`) **vetoes** → `strict_gate_test.go:179`
passes. Claude's round-2 masking chain is refuted; the withdrawal is correct and mine to
ratify.

**One correction of my own on the same point, and one of kimi-1's.** Mine: my round-3 wrote
"the hosted failure **is** an observed no-veto on a file-where-directory, exactly as you
argued" — that sentence leans too hard on testimony; the precise state is: the hosted log shows
`:179` firing (`TESTIMONY`), and `:179`'s `t.Fatal` can only fire on a no-veto, so the
no-veto is a strict deduction from the test's own shape at the pinned commit
(`git diff 868825f..HEAD` on both files is empty) — I keep the deduction, minus the "exactly as
you argued", since claude-1's *mechanism* for it is now refuted. Kimi-1's: V1's "sharpening"
claims the source "narrows the hosted failure to exactly one path" — the
`ERROR_FILE_NOT_FOUND`-restart `break` in `dir_windows.go` yielding `(empty, nil)`. That
narrowing misses the same guard: the break path requires the **open** to succeed on a regular
file, which the `:455-465` guard prevents at both toolchains. With claude-1's round-2 chain and
kimi-1's narrowing both excluded, **no source-level mechanism remains that produces the
observed no-veto** — some premise fails on the actual runner (toolchain resolution, FS-driver
behaviour at `CreateFile`/`GetFileInformationByHandle`, or something not yet enumerated). This
is sharper than round 3 left it, and it is why the probe is load-bearing for the whole
`ErrNotExist`-after-directory-read sweep, not just this test. Verdict unchanged in substance:
phenomenon real, cause **undiagnosed**, not fixture-only, invariant fix lands under every
branch.

### N4 — the executed traversal/`IsLocal` table: **CONFIRMED by independent execution**

I re-implemented `GatePath` (`gate.go:104-106`) and `BlockWorkspace` (`executor.go:29-31`)
verbatim and executed the table on darwin (`EXEC`; `\` inert here, so the Windows column is
hosted/`GOOS=windows`-pending, exactly as claude-1 scoped it). Every row reproduces:
`../x` escapes `gates/` (`pipelines/run/x.gate.json`); `../../x` reaches
`pipelines/x.gate.json` and `ideas/x`; `../../other-idea` lands the workspace in the sibling
`ideas/other-idea` and the gate in `pipelines/other-idea.gate.json`; `../../../../x` exits the
deck (`/x.gate.json`); separators (`a/b`, `a\b`, `..\..\x`), drive/ADS shapes (`C:evil`,
`C:\evil`, `a:b`, `a:$DATA`), reserved device names (`CON`/`con`/`NUL`/`COM1`/`LPT9`/`AUX`/
`PRN`) and `end.` all behave as tabled. The load-bearing negative reproduces exactly:
`filepath.IsLocal(slug+"__"+id)` is **true** for `../../x` and `../../other-idea` while `Join`
lands outside the block workspace — and false only for `../../../../x`; `IsLocal(id)` alone is
false for all the `..`-forms. Four rows show the concat/raw divergence. Conclusion ratified:
**`IsLocal` on the concatenated element is unsound as a guard; validation runs on the raw ID,
and the primary defence is an explicit charset allowlist.** The `filepathlite` precedent is
real (`GOROOT-1.27.1`, `internal/filepathlite/path_windows.go:23-53`: any `:` rejected, every
reserved name rejected) — cited for the allowlist's contents, not reused. Product premises
re-verified at `ac1f0f7`: `GatePath` and `BlockWorkspace` interpolate the ID raw;
`.tmp` staging inherits by construction (`gate.go:119`); manifest load (`manifest.go:~156-173`)
currently validates emptiness/duplicates/kinds only — no charset check, which is the gap the
strict-new/lenient-load/backstop design closes.

### Verdict on the Microsoft guarantees themselves (`MS-DOC`, fetched this round)

- **CreateFileW, verbatim as quoted by claude-1:** "A write-through request via
  `FILE_FLAG_WRITE_THROUGH` also causes NTFS to flush any metadata changes, such as a time
  stamp update or **a rename operation**, that result from processing the request." The
  sentence is real, current, and scoped to metadata changes resulting from processing the
  request. "Such as" is exemplary, not exhaustive — which is exactly why the **create-entry**
  coverage is an **inference** (a create's directory entry is metadata of the same
  directory-entry class as the named rename example, resulting from processing the create
  request) and not a documented guarantee. Claude's refusal to blur that line is correct and I
  adopt it word for word: rename-class coverage is **documented**; create-entry coverage is a
  **documented-mechanism inference**, disclosed as such in FINAL and corroborated hosted.
- **FlushFileBuffers:** the page's subjects are a file ("must have the `GENERIC_WRITE` access
  right"), a communications device (transmit buffer), the server end of a named pipe, and a
  volume ("caller must have administrative privileges"). **Directory handles are absent from
  the documented subject list entirely.** Claude's verdict against the evidentiary status of
  the directory-handle flush route is therefore not a judgement call but a consequence of the
  doc's own enumeration: a green spike could only ever prove "the call returned nonzero", which
  the brief disqualifies by name. I confirm that verdict.
- **MoveFileEx** (fetched in my round-03, re-cited): "The function does not return until the
  file is actually moved on the disk" — the sentence that covers the same-volume rename, the
  actual operation at our rename sites; the explicit "flushed to disk" sentence is scoped to
  copy-and-delete moves, which `replace_windows.go:10-13` deliberately excludes. Claude's
  honest note matches my round-3 amendment (b) verbatim; the two doc readings are now the same
  reading.

## Supersessions of my own earlier verdicts and claims

1. **Site count.** My round-3 "~7 directory-sync call points across 5 files" conflated
   Windows-reachable sites with the Unix-only ones. Superseded by N2's census: **six**
   Windows-reachable sites, plus two Unix-only syncs confined by build tag. The mechanical
   grep-provable-audit contract from my round-3 stands and is what makes the census
   review-checkable; the hand-list is retired.
2. **Durability ordering.** My round-3 ordering was write-through → spike (documentation bar)
   → refusal. Superseded: the spike's acceptance bar is now **shown unsatisfiable from the
   `FlushFileBuffers` documentation itself**, so the ordering collapses to
   **write-through → scoped refusal**. The spike survives only as a diagnostic, explicitly not
   acceptance evidence. This also supersedes the residual-spike branch in kimi-1's round-3
   item 3 / H1 ("if the spike succeeds *and* authoritative documentation scoped to
   directory-handle flush exists, implement it") — **by kimi-1's own bar**: the required
   documentation does not exist, so the residual outcome is predetermined (refusal) and the
   contract gets simpler, not narrower. I flag this as a supersession of the letter of kimi-1's
   H1, not a dispute with its substance — kimi-1's bar was the right bar and it is what kills
   the branch.
3. **No-op floor.** Already withdrawn in my round-3 position change 1; re-stated because the
   brief makes it load-bearing again below: no documented-nil/weakened default anywhere, ever,
   without hosted evidence + reviewable deviation proposal before any owner question.

## The one exact conditional durability contract

Per the brief: one contract, complete site inventory, proved scoped mechanism **else** named
blocking refusal for the affected operation, no implicit no-op, no pre-approved weaker
guarantee. This is the sentence set I will carry into FINAL as drafter:

1. **Inventory (mechanical).** Exactly the six Windows-reachable directory-sync sites of N2.
   Review proves by grep at the reviewed commit: no directory handle reaches `SyncFile` on
   Windows; every Unix-only directory sync remains behind `//go:build !windows`. The audit is
   re-run mechanically at review; hand-lists are not evidence.
2. **Mechanism per site class.**
   - **Create-publication sites** (`trajectory_verify.go:146`→:138, `verification.go:202`→:189,
     `reservation_recovery.go:75`→:67): add the stdlib per-GOOS `syscall.O_SYNC` to the
     existing `O_EXCL` publication open, set **only from Windows-confined code** (verified
     per-GOOS: `0x80` darwin ≠ `0x1000` windows; a shared literal would change Unix write
     behaviour, which the owner forbids). The durability request is the write-through flag on
     the create itself; the directory-sync call at these sites is **removed on Windows** and
     unchanged on Unix.
   - **Rename-publication sites** (`parent_recovery.go:333`→:308, `:570`→:545,
     `unchanged.go:279`): publication goes through `MoveFileEx(MOVEFILE_WRITE_THROUGH)`
     (documented for the actual same-volume move: "does not return until the file is actually
     moved on the disk"). `unchanged.go:279` already does (via `writeState` →
     `ReplaceSyncedFile`); `parent_recovery.go:333/:570` convert `dir.Rename(stage, final)` to
     the write-through sibling replace. Both already stage-and-fsync a sibling, which is
     exactly `ReplaceSyncedFile`'s stated precondition. **Containment sentence (design, not
     implementation shape):** the write-through replace used at rooted sites must keep both
     from- and to-paths inside the root — resolved under `Root.Name()` with the existing
     sibling-staging check (`replace_windows.go:19-21`) kept — so the conversion does not
     trade away `os.Root` containment; the helper's signature is code shape, not a design
     decision. Directory-sync call removed on Windows, unchanged on Unix.
3. **`fsutil.SyncDir` semantics.** `SyncDir` is the named per-OS contract. Unix: exactly
   today's open+fsync behaviour, byte-identical. Windows: **the refusal emitter** — any site
   that still calls it returns the named, user-visible, blocking refusal scoped to that
   publication (error names the operation and the missing directory-entry durability barrier).
   There is **no branch** in which Windows `SyncDir` returns nil: converted sites stop calling
   it on Windows; unconverted sites refuse through it. Fail-closed by construction, including
   for future call sites added without conversion.
4. **The O_SYNC proposal — verdict: accepted conditionally, evidence-pending.** The source
   mapping is proved at both toolchains (N1); the documented sentence is proved (MS-DOC); the
   create-entry coverage remains a **disclosed inference**. Outcome→behaviour mapping, fixed in
   advance so no probe result can silently select policy:
   - hosted probe shows the `O_SYNC` `O_EXCL` open (plain and `os.Root` shapes) succeeds on
     NTFS at the CI toolchain → implement the mechanism at the three create sites; FINAL quotes
     the `CreateFileW` sentence **and** discloses the create-entry coverage as an inference in
     the honesty envelope (never upgraded to "documented");
   - the open fails or the flag is refused by the filesystem → those three sites fall to
     `SyncDir` → **named scoped refusal**; no third branch;
   - a green open is **not itself** the guarantee — the guarantee claim rests on the documented
     sentence plus the disclosed inference, and FINAL says so in exactly those terms. No
     stronger evidence is obtainable this run (no power-cycle testing on hosted runners); the
     envelope states that too.
5. **File-handle fsync** stays mandatory and unchanged everywhere, both OSes.
6. **Deviation rule.** A weakened guarantee enters only as a reviewable proposed deviation
   backed by concrete hosted evidence, before any owner question; never a default, never
   selected silently by an implementation. If FINAL ships a documented-weakening default I
   expect claude-1's MAJOR at review and I would co-sign it.

## `strict_gate` reconciliation (source correction ↔ hosted probe ↔ invariant test)

- **What source settles** (both toolchains, `GO-1.26.8` + `GOROOT-1.27.1`): the
  `O_DIRECTORY` post-open guard and the `ErrNotExist` set predict veto-then-pass. This is the
  **source-predicted branch, not a confirmed runtime branch** — nobody on this run has
  executed it on Windows, and the hosted log contradicts it, so at least one premise is false
  on the runner. No round file, including this one, may call it confirmed.
- **What the probe settles:** prints `runtime.Version()`, the raw `os.ReadDir(<regular file>)`
  error, and its `errors.Is(err, fs.ErrNotExist)` classification, at the CI toolchain. Every
  outcome is pre-mapped (adopting claude-1's C2, unchanged): non-NotExist error → fail-open
  hypothesis withdrawn entirely, cause returns to undiagnosed, ledger records the branch;
  `err == nil` + empty, or a NotExist-classified error → fail-open **confirmed** and it is a
  safety defect, not a portability nit. Either way the failure may not be filed as
  fixture-only without this result — all three of us now hold this.
- **What lands regardless:** the branch-independent invariant fix — stat the round path, veto
  when it exists and is not a directory (`internal/driver/impl.go`), with a cross-platform
  test pinning the file-where-a-directory-belongs veto — plus the Stage-6 sweep of every
  `errors.Is(err, fs.ErrNotExist)`-after-directory-read decision (kimi-1's pattern-keyed
  framing, adopted). The probe is load-bearing for the sweep because N3+kimi-correction leave
  **zero** surviving source mechanisms for the observed no-veto.

## Traversal/IsLocal, inapplicability and census policy — resolved

- **Traversal/naming:** settled by N4 as executed. Universal encoding inside `GatePath` on
  every OS (`PathEscape` base + explicit `:` escape + reserved-device-name and
  trailing-dot/space escapes; staging inherits), strict charset allowlist on the **raw ID** at
  new-pipeline creation, lenient load of existing manifests with cosmetic-only grandfathering,
  unsafe class refused at load **and** at both interpolation backstops (`BlockWorkspace`,
  `GatePath`) on every OS with no legacy exemption; legacy `->` read-fallback one release;
  shadow retirement on first rewrite; mixed-name deck test. `IsLocal`-on-concat is barred as a
  guard by the executed table; `filepathlite`'s Windows `isLocal` cited as precedent for the
  allowlist contents.
- **Census numbers, settled by execution:** `grep -rn 't\.Skip' --include='*_test.go' internal`
  = **104** lines, of which 103 are call sites (**68** `t.Skip(` + **35** `t.Skipf(`) across
  **56** files — claude-1's round-3 correction is exactly reproduced; the round-2 "107"
  (unescaped `.`) is dead, kimi-1's V4 count included it. Mechanism as settled in round 3:
  emit `-json` on the Windows leg (other legs unchanged), enumerate every SKIP event, diff
  against the reviewed exclusion table; a firing skip without a reviewed row fails the leg;
  the workflow change is its own reviewed change; unfiltered matrix preserved.
- **Inapplicability policy:** each claimed inapplicability is reviewed individually with a
  row; the census treats a firing skip as non-execution evidence, never behaviour-exercised
  evidence. AF_UNIX: **no exclusion row** (build tags admit it; `syscall.Mkfifo` call tagged;
  the `t.Skipf` fallbacks in `tree_report_test.go:116/:119/:126/:250` become hard failures on
  the Windows leg); only if hosted AF_UNIX creation itself fails does the individually
  reviewed exclusion route open, with the reason recorded. ACP: the two-test split (gated
  small-payload capacity-independent; ungated 16384→8192 ring cap) is settled. CRLF: the split
  rule is settled verbatim — normalization admissible for document-semantic comparisons,
  forbidden for evidence-integrity comparisons; load-bearing layer is per-invocation
  `-c core.autocrlf=false -c core.eol=lf` on parley's own git calls, repo `.gitattributes`
  second as its own reviewed commit.

## Self-containment of the contract (drafter's check, per the brief)

The contract is draftable self-contained; no design decision is delegated to implementation:

- **Existing-store recovery:** verification-on-open refuses grants beyond the owner
  (fail-closed, offending trustee named). Recovery is **documented, user-invoked, and off the
  product's write path**: the refusal text names the trustee and points at documented
  recovery; this idea ships **no automatic DACL rewrite of existing stores** (kimi-1's "never
  silent in-place modification" and my round-3 "no in-flight migration code" converge on
  exactly that sentence; a repair *command* is follow-up scope if ever wanted). Re-create
  under the protected DACL is the documented path; the hosted `%TEMP%` parent-ACL dump is
  fixture logistics and never weakens the policy (kimi-1's H6 framing, adopted).
- **Partial publication/refusal recovery:** already fixed by the product's own semantics —
  `O_EXCL` leaves any partial artifact visible and non-replayable
  (`verification.go:192-193`, `reservation_recovery.go:79`), staged renames carry
  `defer dir.Remove(stage)` and randomized staging names, and a durability refusal propagates
  as a loud named error leaving either no artifact (pre-create refusal) or a torn one that
  later verification fails loudly. No new recovery semantics are needed; FINAL states this
  rather than implying it.
- **Every remaining unknown is a hosted fact with a pre-agreed outcome→behaviour mapping**
  (this file + round 3's H-lists, now merged): O_SYNC open success/failure (contract §4),
  ReadDir probe branches (C2/H2/H5), AF_UNIX hosted creation (H3), CRLF mechanism dump (H4),
  `%TEMP%` ACL dump (H6), `wait`/`usage` first hosted execution and sharing-violation
  classification (H7: structural fixes first, self-vs-foreign classification before any
  bounded retry, loud, full verification re-run). No probe result selects a policy silently.

## Consensus status

**Consensus-ready.** I hold no blocking counterproposal. All four N-claims are independently
verified above; the durability contract is settled in its conditional-refusal-default form
with no third branch; `strict_gate` is reconciled (source-predicted ≠ confirmed; invariant fix
regardless; probe decisive for the sweep); traversal/census/inapplicability policy is executed
and settled; the owner-default update preserves this run's organization and my
drafter/implementer role with claude-1 and kimi-1 as reviewers. The only items I would still
name are disclosures, not disputes: (i) kimi-1's H1 residual-spike letter is superseded by
kimi-1's own documentation bar — a one-line acknowledgment in signoff or FINAL is enough; (ii)
claude-1's `:64` locator becomes `:63/:67` in FINAL; (iii) my own round-3 "exactly as you
argued" phrasing is corrected above. Release-gate facts unchanged and re-checked read-only
this round: predecessor handoff 1 present (1.49.1 at `54e0798`, winget held), designated
implementer handoff still **absent**, so the release gate remains closed and version selection
stays anchored above 1.49.1.

## Scope compliance

Read-only verification only: repo source at `ac1f0f7`, GOROOT go1.27.1, fetched go1.26.8
sources (go.googlesource.com) and Microsoft Learn `CreateFileW`/`FlushFileBuffers` pages over
`curl`, git metadata, executed Go controls under `/tmp`, read-only existence checks of the two
release-gate handoff files. Wrote exactly this file (`round-04/zcode-1.md`). No product code,
workflow, or peer artifact touched; no commits, tags, publications, hosted runs, releases or
phase transitions; no driver or second organizer launched; no secrets; no Chrome. English only.
