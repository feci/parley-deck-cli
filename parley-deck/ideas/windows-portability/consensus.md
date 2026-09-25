---
idea: windows-portability
drafted-by: zcode-1
date: 2026-09-25
---

## Protocol context attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet: `.parley-runtime/protocol-packets/full-phase3-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md`,
re-verified by `shasum -a 256` this session. Original owner quotations below are given in the
owner-supplied English translations only; the originals were answered in Slovak and are cited as
such without reproduction.

## Drafting context (process notes, not decisions)

- **Role concentration (§15.5), one line:** the facilitator (codex-1) is a non-participant
  organizer and did not draft this artifact; the drafter (zcode-1) is a participant holding no
  facilitator role in this idea — no facilitator-drafter concentration exists. The §15.5/§15.6
  drafter duties below are included per the Phase-3 drafter-facing duty.
- **Drafting basis.** The organizer's consensus-boundary section in `00-prompt.md` ("Round 6
  completed — consensus opened, 2026-09-25") assigns this draft to zcode-1 on the current refusal
  branch; the claim is confirmed in `inbox/zcode-1-to-all_windows-portability_drafter.md` before
  any signature. Quorum is unchanged: exactly claude-1, kimi-1, zcode-1 as participants; codex-1
  organizes only and never signs, implements, or verifies code.
- **Round history this consensus synthesizes.** Rounds 1–3 (2026-09-24/25) merged the six-item
  tri-participant contract; round 4 opened the durability conflict; round 5 (owner-authorized,
  focused) left it unresolved; round 6 (owner-authorized, sequential: claude-1 16:41:25 UTC,
  kimi-1 16:59:00 UTC, zcode-1 17:06:07 UTC, each exiting 0 and validated before the next launch)
  converged all three positions on all three conflicts. No seventh round is authorized or
  requested.
- **Blocker statement.** No substantive inter-participant blocker remains: all three round-06
  artifacts hold identical positions on every decision below. The one open external item — the
  owner's decision on the containment-deviation question — is explicitly **nonblocking** for this
  consensus: the selected design is complete and implementable without it, and nothing in this
  document conditions baseline implementation on implied approval or treats silence as approval.
- **Verification state of this draft.** Every census number and locator class cited below was
  re-executed or re-read by the drafter at HEAD `3cf0068` (`git diff --name-only 6b87cf8..HEAD --
  '*.go'` = 0 files, so every round-05/06 locator transfers); peer verifications are cited by
  round with their §15.2 tags in the round files, which remain the canonical argument record.

## Agreed decisions

Unless a line says otherwise, each item carries all three participants' round-06 assent ("stands
unreopened" in all three files) and the rounds where it was established and contested.

### 1 — Selected design: the refusal branch, preserving `os.Root` containment

The original owner scope authorizes "a real Windows implementation **or** an explicit, reviewed,
user-visible refusal"; the owner's round-6 answer reserves any weakening of `os.Root` containment
to an explicit owner deviation ("it must come back to the owner as an explicit deviation" —
owner-supplied translation, original in Slovak). All three round-06 positions adopt the
**plain/rooted split** as the strongest design the deck may select on its own authority:

- **Plain sites restructure** (no `os.Root` to lose): A1, B1, B2 — all currently dormant behind
  the existing reviewed POSIX refusal at `trajectory_verify.go:106-108`, so this is the design of
  record if that gate ever lifts; **D1 derives** (its original publication `writeState`
  `state.go:197-221` already ends in `fsutil.ReplaceSyncedFile`, the reviewed write-through
  replace).
- **Rooted sites take the named blocking refusal, pending an explicit owner deviation that has
  not been granted**: A2, A3, B3, B4, B5, C1, C2; refusals at C1/C2 are **ordered before the
  rename** (N8: the barrier fires after publication at `parent_recovery.go:336`/`:573`, so the
  refusal must be evaluated before `:333`/`:570`, leaving `defer dir.Remove(stage)`
  `:328`/`:565` to clear the stage with nothing published). D2/D3/D4a/D4b follow their originals
  into refusal (the coherence rule: a class-D derivation holds only if its original carries a
  proved mechanism).
- The refusal is named, pre-mutation, user-visible, and blocking — never a silent skip, never a
  green-API-call-as-guarantee, never an implicit no-op. There is no fourth branch.
- **The weaker-containment alternative is carried as UNAPPROVED, not as an implementable
  fallback**: the organizer's nonblocking deviation request is
  `inbox/codex-1-to-user_windows-portability_containment-deviation.md`; its full specification is
  §"Open items" item 1 below. No implementation may select it without the owner's explicit
  approval **and** participant signoff on the resulting contract; silence is not approval. If
  the owner declines, the refusal-branch table is final as written — a reviewed outcome, not a
  failure.

### 2 — Durability mechanism for the sites that do convert: staged write-through rename with deterministic stages (Conflict 1, resolved)

All three participants now select, for class A, the staged write-through rename: stage create
(`O_EXCL` + mandatory file fsync) then `MoveFileEx(MOVEFILE_WRITE_THROUGH)` without
`MOVEFILE_REPLACE_EXISTING`, with **claude-1's deterministic-stage condition**: the stage name is
a fixed derivation of the final name, created with `O_EXCL`, and **not removed on the error
path**, so an interrupted publication leaves a visible blocker whose existence prevents replay —
exactly the invariant `verification.go:192-193` ("its existence prevents replay") and
`reservation_recovery.go:79` ("a partial publication stays visible and cannot be overwritten by
retry") supply today. The final name appears only via the atomic move, so readers never see torn
bytes; no-REPLACE prevents clobber of a completed publication; concurrent publishers race on the
stage `O_EXCL` exactly as today they race on the final-name `O_EXCL`.

Drafter requirements carried from round 6: (a) the stage derivation must be provably
non-colliding with every valid final name (final names are validated or product-constructed, so
the stage suffix must sit outside every valid-final grammar, pinned by a unit test); (b) the
publish-time `validHash` gate at A3 (`reservation_recovery.go:75` builds
`intentName(i.Accounting.EntryKey)` ungated, while the read path gates at `:92-93`; under any
restructure the key feeds both stage and final names).

Documentation basis, honestly tiered (unchanged since round 5, tags in the round files): the
write-through sentence — "does not return until the file is actually moved on the disk" — is
scoped to the exact operation (PRIMARY fetches by all three in round 5, independently); there is
**no** documentation scoped to the parent-directory entry of a create (N5, CONFIRMED non-owner by
kimi-1 and zcode-1 in round 5 after claude-1's withdrawal); `O_SYNC`→`FILE_FLAG_WRITE_THROUGH`
is real for **write** durability only; `FlushFileBuffers` on our read-only directory handles is
foreclosed twice over (no documented directory subject; documented `GENERIC_WRITE` precondition
unmet) and stays **diagnostic-only**. The no-REPLACE existing-destination error basis is
disclosed per kimi-1's N6 limitation: flag-table contrapositive plus a hosted existence-probe
assertion, never upgraded to "documented". The directory-move guarantee for class B is a
disclosed composition (N6's "existing file **or directory**" subject sentence plus the
write-through sentence), labelled as such in FINAL, pinned by mechanics probes only.

### 3 — The normalized operation table: 15 IDs / 16 runtime operations (Conflict 3, resolved)

The three historical inventories (14, 15, 16) were the same set under different individuation:
14 + (D4 split) = 15 + (B1 second loop iteration) = 16. Derivation re-executed at `3cf0068`:
the name-keyed grep over the six helpers yields 15 call lines; minus `reservation_recovery.go:67`
(inside `syncIntent`'s own body) = 14 name-visible; plus `parent_recovery.go:336` (invoked
through a function value bound at `:311`, invisible to name-keyed grep) = **15 independent
sites**; B1 fires twice per run (`:117-128` loop) = **16 runtime operations**. Row IDs are
kimi-1's scheme (strict file-then-line order). The audit of record is zcode-1's
`fsutil.SyncDir` named-type contract (compiler-enforced fail-closed on Windows; never a grep) —
under the refusal branch it is also the rooted sites' refusal emitter. N10's corrected helper
`func` lines: `trajectory_verify.go:132`, `verification.go:183`, `unchanged.go:267`,
`parent_recovery.go:295`, `:532`, `reservation_recovery.go:55`.

Abbreviations: WT-move = `MoveFileEx(MOVEFILE_WRITE_THROUGH)`; no-REPLACE = without
`MOVEFILE_REPLACE_EXISTING`; det-stage = deterministic stage per §2.

| ID | Site (publication → barrier) | Class | Rooted? | Win-reachable today | Required guarantee | Disposition (refusal branch = selected design) | If deviation granted | Feature blast radius if refused | Partial-publication recovery |
|---|---|---|---|---|---|---|---|---|---|
| **A1** | `trajectory_verify.go:146` `os.OpenFile O_EXCL` → `:154` | A create | plain | no — gated `:106-108` | create entry durable; anti-replay | **Restructure**: det-stage `O_EXCL`+fsync → WT-move no-REPLACE | same | none today; inherits if gate lifts | stage persists, blocks replay; final never torn |
| **A2** | `verification.go:202` `Root.OpenFile O_EXCL` → `:211` | A create | rooted | no — gated `:251-253` | same + `:192-193` invariant | **Refuse**, pre-mutation | det-stage → WT-move no-REPLACE + guard set | none today | `:192-193` invariant intact (unconverted) |
| **A3** | `reservation_recovery.go:75` `Root.OpenFile O_EXCL` → `:84`→`:67` | A create | rooted | **yes** | same + `:79` invariant | **Refuse**, pre-mutation | as A2 + publish-time `validHash` gate | precharge reservation-intent publication refuses → `PrepareCycleReservation` fails (`cycle_binding.go:269`, no GOOS gate) | nothing published; `:79` invariant intact |
| **B1** | `trajectory_verify.go:118` `os.Mkdir` → `:125` — **fires ×2** (`:117-128` loop) | B mkdir | plain | no — gated `:106-108` | base-dir entries durable, idempotent | **Restructure**: stage-`Mkdir` + WT-move + EEXIST-tolerant recheck | same | none today | IsExist tolerance preserved; nothing published |
| **B2** | `trajectory_verify.go:190` `os.Mkdir` → `:193` | B mkdir | plain | no — gated `:106-108` | run-dir entry durable | **Restructure** as B1 | same | none today | as B1 |
| **B3** | `verification.go:264` `parent.Mkdir` (IsExist-tolerant) → `:267` | B mkdir | rooted | no — gated `:251-253` | store entry durable, idempotent | **Refuse**, pre-mutation | stage-`Mkdir` + WT-move + tolerant recheck | none today | as B1 |
| **B4** | `verification.go:283` `base.Mkdir(charge)` — anti-replay token `:281-282` → `:286` | B mkdir | rooted | no — gated `:251-253` | reservation entry durable; never reuse | **Refuse**, pre-mutation | det-stage-`Mkdir` + WT-move no-REPLACE; existing dir = loud either way (N6) | none today | `Mkdir` failure = "already reserved" `:284` |
| **B5** | `reservation_recovery.go:36` `dir.Mkdir` (IsExist-tolerant) → `:47` | B mkdir | rooted | **yes** | intent-root entry durable, idempotent | **Refuse**, pre-mutation | stage-`Mkdir` + WT-move + tolerant recheck | `openIntentRoot(create=true)` fails → reservation creation unavailable | nothing published; `Mkdir` idempotent |
| **C1** | `parent_recovery.go:333` `dir.Rename` → `:336` (func value bound `:311`) | C rename | rooted | **yes** | rename entry durable; publish-once | **Refuse, ordered before `:333`** (N8) | WT-move; flag chosen deliberately and disclosed (drafter preference: no-REPLACE), never inherited by pattern-match | parent-recovery **apply** refuses | pre-rename refusal ⇒ nothing published; `defer dir.Remove(stage)` `:328` clears stage |
| **C2** | `parent_recovery.go:570` `dir.Rename` → `:573` | C rename | rooted | **yes** | same | **Refuse, ordered before `:570`** | as C1 | recovered-parent **apply** refuses | as C1; stage cleared `:565` |
| **D1** | `unchanged.go:241` → `:279` (func `:267`) | D re-sync | plain | **yes** | already satisfied | **Derived no-op**, argued per site: `writeState` (`state.go:197-221`) already publishes via `ReplaceSyncedFile` = WT-move | same | none | `os.CreateTemp` + `defer os.Remove` (`state.go:205`/`:209`); replace atomic |
| **D2** | `parent_recovery.go:368` → `:308` | D re-sync | rooted ctx | **yes** | inherits C1 | **Refuse** — derives from a refused publication | derives once C1 converts | tracks C1 | re-sync only; never rewrites |
| **D3** | `parent_recovery.go:605` → `:545` | D re-sync | rooted ctx | **yes** | inherits C2 | **Refuse** | derives once C2 converts | tracks C2 | re-sync only |
| **D4a** | `reservation_recovery.go:420` `syncIntent` | D re-sync | rooted ctx | **yes** | inherits A3/B5 | **Refuse** | derives once A3/B5 convert | tracks A3 | re-sync only; `:426` never rewrites bytes |
| **D4b** | `reservation_recovery.go:431` `SyncFile(trajectory.json)` + `syncVerificationDirectory` | D re-sync | rooted ctx | **yes** | inherits A3/B5 | **Refuse** | derives once A3/B5 convert | tracks A3 | re-sync only |

**Totals:** 15 IDs; 16 runtime operations; 9 of 15 Windows-reachable today (A3, B5, C1, C2, D1,
D2, D3, D4a, D4b); 6 dormant behind the two existing reviewed POSIX refusals. All conversions are
Windows-confined; the Unix path at every row is byte-identical to today.

### 4 — Feature-level refusals (the honest, user-visible cost on Windows)

Under the selected branch, on Windows:

1. **Precharge reservation-intent publication** refuses (A3 + B5, dragging D4a/D4b):
   `openIntentRoot(b, true)` ← `PrepareCycleReservation` (`reservation_recovery.go:126`) ←
   `cycle_binding.go:269`, no GOOS gate. Budget reservation intents cannot be published.
2. **Parent-recovery apply** (C1 + D2) and **recovered-parent apply** (C2 + D3) refuse: the
   recovery publication commands are unavailable.
3. **Already refused today and unchanged** (the refusal branch takes away nothing that currently
   works): trajectory verification (`trajectory_verify.go:106-108`) and captured-verification
   journals (`verification.go:251-253`) carry existing reviewed POSIX refusals on Windows.
4. **Blocking missing-`sh` refusal** (§5.3 below) and the durable-kill/attribution refusals under
   P-A are pre-existing, designed, user-visible refusals — first-class outcomes under the owner
   scope, each with actionable refusal text.

Every refusal is named, its trigger is before the mutation, its message is user-visible and
actionable, and it is covered by a hosted test pinning the refusal path. FINAL enumerates every
refused operation with its blast radius and recovery column from the table — never a headline
count.

### 5 — Settled non-durability requirements (rounds 1–3, tri-participant, unreopened in rounds 4–6)

5.1 **Snapshot privacy/ACL (W-PRIV).** Creation sets an **owner-only protected DACL**: exactly
one access-allowed ACE for the token user (`GetCurrentProcessToken().GetTokenUser()`) in a
protected, non-inherited DACL (`PROTECTED_DACL` load-bearing), via
`ACLFromEntries`/`SetNamedSecurityInfo`, all in the already-direct `golang.org/x/sys v0.36.0`
(zero new dependencies; `hectane/go-acl` and `icacls` rejected). Verification walks effective
granted trustees via `GetNamedSecurityInfo` and refuses any non-owner grant, inherited or
explicit, **naming the trustee**; refusal string sentence-stable ("snapshot store must be a
private real directory" — the 71× hosted signature) with the trustee appended. FAT/exFAT/no-ACL
volumes refuse, never weaken. Pre-existing stores: **refuse-and-instruct** (move aside,
re-create under the protected DACL; explicit user-invoked documented repair named in the refusal
text) — **no in-place DACL rewrite of user data in this idea**. Adversarial hosted tests:
`BUILTIN\Users` grant → refuse; `Administrators` grant → refuse; inherited-only → refuse; own
creation → pass; create-then-requery readback self-check. `!IsDir()` and reparse/symlink
rejection kept on both platforms; Unix behaviour byte-identical. (kimi-1's owner+SYSTEM
preference is recorded as a minority note, withdrawn r3.)

5.2 **Process liveness P-A; P-B deferred.** Truthful `alive` via
`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` + `GetExitCodeProcess`, `STILL_ACTIVE` sentinel
caveat documented, **fail closed on `ERROR_ACCESS_DENIED`** (unverifiable is alive for refusal
purposes); `bootID` stays `""` (scoped null); `supportsDurableKill()` stays false with
Windows-honest durable-kill/attribution refusal strings; `proclive_windows.go` comment corrected
(x/sys is direct at `go.mod:10`). Fixes the four broken callers (`durablekill.go:29`/`:63`,
`verification.go:142`, `reviewsnapshot.go:244`). **P-B** (Job Objects
`JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `GetProcessTimes` attribution,
`QueryFullProcessImageName`) is a **named follow-up idea** with its own adversarial review: the
Windows image-path is coarser than the Linux full-argv calibration of `procctl.commandMatches`,
and adopting it in-flight would silently loosen a kill-safety gate. `mvdan.cc/sh/v3` supervisor
port likewise a separate reviewed idea (new dependency executing untrusted `checks:` strings).

5.3 **Blocking missing-`sh` refusal (W-SHELL).** On Windows, `checks:`
(`driver_impl.go:266`) and captured verification (`evidence/execute.go:118`/`:122`,
`criterionSupervisor` `:100-114`) call `exec.LookPath("sh")` **once, before any work**; on
failure they return a named, actionable, blocking refusal naming **Git for Windows** as the
documented prerequisite (template: name the context, name the prerequisite, "refusing rather
than executing unverified"). The gate never passes by default; LE-4 stays intact; exits
non-zero; recorded in the run. Hosted tests pin both paths (`sh` forced off `PATH` → message +
non-zero; `sh` present → normal path) — because GitHub's windows-latest ships Git for Windows
`sh.exe` on PATH, an untested leg would only prove the image happened to include it.

5.4 **Universal gate-name encoding + raw-ID charset allowlist (W-NAME).** One universal encoding
on **every OS**, applied inside `GatePath` (`gate.go:104-106`) so the `.tmp` staging name
inherits it by construction (the 9–10× hosted `ERROR_INVALID_NAME` failures were on the staging
name): `net/url.PathEscape` base plus explicit `:` escape plus reserved-device-name and
trailing-dot/trailing-space guards (`block-a->block-b` → `block-a-%3Eblock-b.gate.json`).
Per-OS naming rejected (decks cross OSes; a Unix-written `>` name cannot exist on Windows).
Legacy `a->b.gate.json` read-fallback for one release; shadow retirement on first rewrite (no
stale HITL answer survives); mixed-name deck test; repo has zero legacy gate names today.
Block-ID validation: strict **charset allowlist on raw IDs at new pipeline creation**
(`[A-Za-z0-9._-]`-class, bounded length, no leading/trailing dot, not a reserved device name),
cosmetic-only grandfathering on load, unsafe class refused **at load and at both interpolation
backstops** (`BlockWorkspace`, `GatePath`) on every OS with no legacy exemption — traversal,
separators, absolute/drive paths, ADS `:`, reserved names (reserved even with extension),
trailing dot/space. Load-bearing negative result (executed table, N4): `filepath.IsLocal` on the
concatenated `slug__+id` is unsound (prefix absorbs `..` inside `Clean`; four of five traversal
IDs flip to wrongly-accepted; `../../other-idea` lands in another idea's directory) — validation
must run on the raw ID; `filepathlite`'s Windows `isLocal` is precedent for the allowlist
contents, not the guard.

5.5 **`strict_gate` branch-independent fix; cause undiagnosed.** The hosted `:179` failure
(product `impl.go:448-472` `reviewRoundHasFindings` returned no-veto on
file-where-directory-belongs) contradicts both source-level mechanisms (excluded at go1.27.1 and
the CI-fetched go1.26.8). Shipped regardless: **stat the round path and veto when it exists and
is not a directory** (+ one `Lstat`), a cross-platform pinning test, a hosted probe printing
`runtime.Version()`, the raw `os.ReadDir(regular file)` error and its `errors.Is(err,
fs.ErrNotExist)` classification at the CI toolchain (read the resolved version from the existing
setup-go log first; `go.mod` has no `toolchain` pin), and the sweep of every
`ErrNotExist`-after-directory-read decision. Cause recorded **undiagnosed**; no branch called
confirmed; not reclassifiable as fixture-only without the probe result.

5.6 **File sharing / cross-process locks (W-LOCK).** The budget locks are already correctly
implemented (`lock_windows.go:16-38` LockFileEx byte-range at 1 MiB + `publishExclusive`
WT-move); `gofrs/flock` and named mutexes rejected. Fixes: share-flag-correct opens — prefer
routing through `os.Root` (passes `FILE_SHARE_DELETE`; already used at `snapshot.go:194`), raw
`windows.CreateFile` helper only where rooting cannot express the site; close-before-rename
discipline; read-only-attribute trap handled (staging files written without `0200` set
`FILE_ATTRIBUTE_READONLY` and cannot be renamed over — clear before replace); adversarial
rename-root tests redesigned to mutate inside the root. Bounded retry (~250 ms) only on
`ERROR_SHARING_VIOLATION`, only for transient **third-party** holders (Defender/AV/indexer),
after self-held vs foreign-held classification, loud, with identity/exclusion verification
re-run. Open-site sharing failures (`.lock` opens) remain **undiagnosed** — one hosted
diagnostic cycle logs error class + retry outcome (Stage 4).

5.7 **ACP split; AF_UNIX; CRLF split.** ACP: **split the two drain tests, do not retune** —
(a) gated small-payload (≤2 KiB, capacity-independent: child READYs and exits before the copier
parks), (b) ungated 16384-write asserting the 8192 ring cap; both assertions survive everywhere;
the correct U1 product drain ordering (`spawn.go` Stop/Wait) is preserved untouched. AF_UNIX:
**no exclusion row** — build-tag the `syscall.Mkfifo` call so Windows falls into the file's
existing `makeUnsupportedSocket` fallback; the test runs and keeps its refusal assertion on
Windows (`TreeDigest` refuses non-regular/non-symlink entries); the `t.Skipf` fallbacks at
`tree_report_test.go:116/:119/:126/:250` become hard failures; only if hosted AF_UNIX creation
itself fails does an individually reviewed exclusion row open. CRLF: **the split rule, settled
verbatim** — normalization is admissible for document-semantic comparisons (protocol drift
anchor, consensus-template block extraction) and **forbidden** for evidence-integrity
comparisons (snapshot content, `TreeDigest`); load-bearing layer is per-invocation
`-c core.autocrlf=false -c core.eol=lf` pinned on parley's own git invocations (in-tree
precedent `reviewsnapshot.go:138-139`; the only layer that follows parley into arbitrary user
repos), repo-scoped `.gitattributes` second as its own reviewed commit; product parsers
comparing embedded-vs-disk bytes are CRLF-tolerant as a product requirement; runner
`git config --show-origin` dump in the probe bundle.

5.8 **Step-1-first build fixes.** The Mkfifo build tag plus the `app_test.go:155/:157` comma-ok
fix (unchecked type assertions abort the whole test binary — why `wait`/`usage` have no Windows
results) land **first, always**; one diagnostic cycle budgeted for whatever they unblock, green
not assumed. Ubuntu leg gains the `GOOS=windows go build ./... && go vet ./...` cross-compile
guard — compile-only evidence, never execution evidence. `internal/evidence` was the sole
compile failure in the tree; both windows/amd64 and windows/arm64 build exit 0.

5.9 **Fixture-portability sweep with the reconciliation ledger.** Per-test reconciliation ledger
(hosted log ↔ class ↔ stage ↔ exclusion row) emitted before the sweep lands and refreshed after
step-1 fixes change the denominator. Classes and fixes: `#!/bin/sh`/extensionless fixtures →
test-binary re-exec (`os.Executable()` + guard; in-repo precedents) or `cmd.exe /c` for trivial
scripts; `.bat`/`.cmd` via PATHEXT only where file discovery is the test subject; Chmod-0
fixtures → one `denyRead(t, path)` helper (chmod 0 on Unix, DACL deny-ACE on Windows, sharing
the §5.1 module) — "Chmod cannot make a directory unreadable" is not a legitimate exclusion once
ACLs exist; HOME fixtures set `USERPROFILE` too; `/repo` literals → `t.TempDir()` +
volume-aware joins with per-OS refusal assertions kept; TUI temp-perms reformulated via the ACL
helper, never silently blessing 0666 synthesis as "private". Sweep rule: a Windows exclusion is
admissible only when the behavior under test **cannot exist** on Windows — not when
inconvenient; one reviewed row per exclusion, attackable individually; diff-level check that no
`t.Skip` was added without a row; close-time suppression grep (`Skip`, `-run`,
`continue-on-error`, `|| true`); each ported fixture names the product behavior it still pins.

### 6 — No-suppression census (CI observability, settled)

Re-executed at `3cf0068` by the drafter, identical to all three round-05 and all three round-06
files: **104** `t\.Skip` pattern lines under `internal`; **103 call sites** = 68 `t.Skip(` + 35
`t.Skipf(` + 0 `t.SkipNow(`; **56 files** matching / **55 with call sites** (sole non-call:
`internal/runner/launch_precheck_skip_test.go:58`). Policy: emit `-json` on **all three legs**
(the workflow change is its own reviewed change; `-json` selects nothing; the unfiltered matrix
is untouched). **Windows leg:** every firing SKIP event matched row-by-row to an individually
reviewed exclusion-table row (owner's "truly inapplicable and reviewed") or the leg fails;
AF_UNIX fallback skips convert to hard failures (§5.7). **macOS/Ubuntu legs:** the
currently-firing skips are the enumerated reviewed baseline, and any **newly firing** skip fails
the leg — this is the only shape that makes the owner's "absence of new suppression" reviewable
on the still-green legs. A firing skip is non-execution evidence, never behaviour-exercised
evidence. CI runs `go test ./... -count=1 -timeout 45m` with no `-v` today (`tests.yml`), so
`--- SKIP` is invisible without `-json` — the census is required, not optional.

### 7 — Coverage honesty envelope

Hosted windows-latest x64 only (Windows Server image, `runneradmin` context, Go 1.26.8,
`core.longpaths true` **set by CI at `tests.yml:38`, not an OS default and not a product
behavior**); ARM64 assets are **build-only with zero executed evidence** and no ARM64 gate is
added; no native-Windows claim by anyone on this run; WSL/Wine are not accepted as Windows
evidence; per-run claims kept separate (historical runs 35987916696 — which also failed Ubuntu —
and 36009912946); runner image recorded with every evidence entry, re-baselined if it drifts;
and one documented sentence that this workspace itself sits on a cross-machine shared volume
whose Windows-side semantics are unvalidated.

### 8 — Roles, staging, and release gates

- **Roles.** zcode-1 is drafter and **single implementation owner** (claim + worktree mapping in
  `IMPLEMENTATION.md` before any code edit; one implementation owner at a time); claude-1 and
  kimi-1 are the independent reviewers; codex-1 organizes only — never implements, verifies
  code, authors participant verdicts, or signs off. Channel verification after release is by a
  participant, never codex-1.
- **Staging.** zcode-1's Stage 0–7 frame, unchanged: Stage 0 contract/smoke/ledger + ubuntu
  cross-compile guard; Stage 1 privacy; Stage 2 gate names; Stage 3 directory durability —
  **opening with the §3 table and its per-row disposition before any code is written**; Stage 4
  file sharing; Stage 5 P-A + 5b shell; Stage 6 sweep; Stage 7 validation/release. One early
  hosted probe bundle (build-tagged, test-only, ~one 20-minute cycle): rename mechanics
  (directory move accepted; no-replace fails on existing destination; WRITE_THROUGH accepted),
  NTFS volume assertion, `SyncFile` on a read-only directory handle with raw error,
  `runtime.Version()` + toolchain pin from the existing setup-go log, ACL parent dump,
  `git config --show-origin`, `ReadDir`-on-regular-file error class, read-only-attribute
  rename-over. Mechanics only — never durability, never acceptance evidence; no hosted probe
  upgrades an inference.
- **Release gates — both now present (N9, confirmed non-owner by kimi-1 and zcode-1 in round 6;
  re-verified read-only by the drafter this session):**
  1. `worktrees/lean-organizer/parley-deck/inbox/codex-1-to-user_release-1.49.1_done.md`
     (6412 bytes; reports CLI 1.49.1 released at `54e0798`, Windows experimental, winget held).
  2. `worktrees/designated-implementer/parley-deck/inbox/codex-1-to-user_meta-protocol-change-designated-implementer_done.md`
     (12075 bytes, mtime 2026-09-25 11:13; `status: agent-controlled-delivery-complete`;
     "This file releases the sequencing hold for `windows-portability`"; reports CLI 1.50.0 and
     skill 2.14.0 released on GitHub and Homebrew).
  Both checks are **existence plus self-testimony**, not channel audits. Consequence: the next
  version is selected **above 1.50.0** (all three round-05 files' "above 1.49.1" is stale),
  after reviewed closure, on latest origin/main, per the owner's release sequence. The Windows
  experimental label is removed **iff** the released commit's hosted Windows leg is green;
  otherwise the label is retained, the CLI winget PR stays held, and the outcome is reported —
  never an incomplete release represented as shipped. Released tags never move. Release includes
  GitHub Windows assets, Homebrew `parley-deck-cli.rb`, and the held winget PR (one application
  per PR); skill channels only if the skill changed.

## Agreed trade-offs

1. **Two reachable Windows features refuse rather than weaken containment.** Reservation-intent
   publication and parent-recovery apply are unavailable on Windows under the selected branch,
   in exchange for preserving `os.Root`'s kernel-enforced per-open reparse resistance and
   handle-relative rename at the publication operation. The owner reserved exactly this trade;
   the deck takes the strongest position available on its own authority and puts one costed
   question to the owner (Open items, 1). Trajectory verification and captured journals were
   already refused on Windows — nothing working is taken away.
2. **Deterministic stage names at class A** (vs randomized-and-swept): retains the
   existence-prevents-replay invariant at the cost of recovery accounting for an `X.partial`
   sibling instead of a torn `X`; stage litter is a visible, actionable blocker by design.
3. **Guarantee tiers are disclosed, never smoothed.** No green API call is ever cited as the
   guarantee; the write-through sentence is quoted as the guarantee for moves; the directory-move
   composition and the no-REPLACE contrapositive basis are labelled; `FlushFileBuffers` is
   diagnostic-only; no power-cycle evidence is obtainable on hosted runners and FINAL says so.
4. **103 pre-existing skips are not re-litigated** — they are pinned as the enumerated baseline
   so that *new* suppression fails the leg; Windows exclusions are reviewed row-by-row
   individually.
5. **P-B and `mvdan.cc/sh` are deferrals, not losses** — named follow-up ideas with their own
   review; the coarser Windows attribution path would have loosened a kill-safety gate
   in-flight.
6. **Sequential round 6 bought convergence at the cost of independence** — see Correlated
   agreement.

## Open items deferred to implementation

1. **The containment-deviation question (owner decision, nonblocking).** The fully specified,
   **unapproved** alternative — recorded by the organizer in
   `inbox/codex-1-to-user_windows-portability_containment-deviation.md` — would convert the seven
   rooted rows to path-based `MoveFileEx(MOVEFILE_WRITE_THROUGH)` with: kimi-1's guard set (1)–(4)
   (root-resolved paths via `Root.Name()`; retained sibling check; every rename name a product
   constant or `validHash`-gated, including the publish-time gate at A3; deterministic stages at
   A / randomized at C; explicit weakened-resistance disclosure), claude-1's items (the precise
   statement of what is given up: kernel-enforced per-open reparse resistance and handle-relative
   rename, replaced by a lexical check plus a path-based call with a TOCTOU window; the named
   but unselected handle-derived-path alternative `GetFinalPathNameByHandle`, which narrows but
   does not close the window; both branches' costs side by side), zcode-1's hosted
   rooted-escape test against the conversion helper, and the C1/C2 flag decision (drafter
   preference, stated as input only: no-REPLACE, matching the publish-once constants — an
   existing final artifact should be a loud replay signal, not an overwrite). If granted, every
   refusal row converts per the already-designed restructure and the D rows derive — no row is
   added or dropped; the exact contract still requires participant signoff. If declined, the
   refusal-branch table is final as written. **Implementation of the baseline does not wait for,
   or assume, this answer.**
2. **Hosted probe bundle execution and the H1–H7 outcome→behavior mappings** (pre-agreed in
   rounds 3–5: each hosted fact's outcome is mapped before it is observed; no result silently
   changes policy; no probe upgrades an inference).
3. **`strict_gate` hosted diagnosis** (fix ships regardless; the probe selects the recorded
   mechanism and the sweep's scope).
4. **Open-site sharing-violation diagnosis cycle** (Stage 4).
5. **Reconciliation-ledger refresh** after step-1 fixes change the failing denominator.
6. **Refusal-string wording and `IMPLEMENTATION.md` claim/worktree mapping** — implementation
   owned by zcode-1, reviewed by claude-1 and kimi-1.
7. **Version selection above 1.50.0 at release time**, after both gates (now present) and
   reviewed closure, per the owner's release sequence.

## Comparison & blind spots

*Advisory drafting discipline; raw round files are never hidden behind this summary.*

- **Partial coverage (what only one participant carried):** claude-1 owned the defect taxonomy
  (source-context Findings A–F, all discharged by non-owner verdicts in round 2), the `O_SYNC`
  branch and its withdrawal, the `IsLocal` unsoundness execution table, and the plain/rooted
  split; kimi-1 owned the census premise (no `-v`/`-json` in CI), the inventory row-ID scheme,
  the N6 flag-table limitation, and the ACP braid; zcode-1 owned the `fsutil.SyncDir`
  fail-closed audit, the Stage 0–7 frame, the reconciliation ledger, the handle-rename closure,
  and the containment verification (`os.Root.Rename` → `SetFileInformationByHandle`, no
  write-through flag).
- **Unique insights worth keeping:** the function-value barrier call site
  (`parent_recovery.go:311`/`:336`) invisible to every name-keyed grep — the reason the audit of
  record is a compiler-enforced type contract, not a grep; and the N10 lesson (all three rounds'
  files mislabeled helper `func` lines — the fifth census-class error): FINAL's table is
  generated mechanically from the enumeration, never transcribed.
- **Blind spots — what no participant addressed:** no power-cycle/real-durability evidence is
  obtainable on hosted runners (disclosed, not solved); Windows ARM64 has zero executed evidence
  (build-only, by design of the coverage envelope); FAT/exFAT refusal is designed but not
  exercisable on the runner; non-`runneradmin` contexts and native desktops unvalidated; P-B
  attribution and `bootID` have no Windows analogue in scope (scoped null / deferred); refusal
  strings have no real-user usability testing; the five census-class errors across rounds argue
  for mechanical generation wherever a number is load-bearing.

## Drafter position changes

**None since `round-06/zcode-1.md`, my most recent round file** — this consensus introduces no
position not already held there. For traceability, the round-06 self-corrections that moved the
drafter from round-05 were: (1) withdrawal of the round-05 claim "replay is blocked by the move's
destination-exists failure instead of by `O_EXCL`" (weakening, after two concordant non-owner
`WRONG` verdicts and my own re-verification of the counterexample); (2) withdrawal of the
round-05 "restructure everywhere, zero blocking refusals" selection as a deck-adoptable position
(superseded by the owner's round-6 authority reservation); (3) adoption of the
deterministic-stage condition, the plain/rooted split, and the 15/16 table on kimi-1's IDs.

## Verdict conflicts

No contradictory verdict is **live** at consensus opening; every conflict that arose in rounds
1–6 is closed, each by owner weakening or non-owner verdicts, never by counting. Per §15.3 each
is quoted with author, tag, evidence, and resolution:

1. **`o_DIRECTORY`/ReadDir chain (strict_gate mechanism).** kimi-1 r3 V1 `CONFIRMED` (deduction
   conditioned on the hosted report) vs claude-1 r3 N3 `SELF-CORRECTION` (owner withdrawal of his
   r2 chain) and zcode-1 r3 `REFUTED at source level` on go1.27.1 **and** the CI-fetched
   go1.26.8 (`PRIMARY` both). Resolution: both source-level mechanisms excluded; hosted
   phenomenon undiagnosed; branch-independent fix ships regardless. kimi-1 r4 struck the
   conditioned half of V1 and superseded V2 (its `CONFIRMED` on claude-1's already-withdrawn
   claim) by owner self-correction — weakening, immediate.
2. **Skip census 107 vs 103.** claude-1 r2's 107 (unescaped `.` in the grep) vs r3
   `SELF-CORRECTION` to **103 call sites** (68+35), confirmed non-owner by both peers and
   re-executed by all three in rounds 5 and 6 and by the drafter this session at `3cf0068`.
3. **Create-entry coverage of `O_SYNC`/`FILE_FLAG_WRITE_THROUGH` (the class-A conflict).**
   claude-1 r4 N5 ("documented against") vs kimi-1/zcode-1 r4-r5 ("no documented mechanism").
   Resolution: claude-1 r5 withdrew N5 (`PRIMARY` re-fetch: the counter-sentence sits in the
   no-write-through configuration; coverage is "neither documented nor documented-against");
   N5-as-withdrawn CONFIRMED non-owner by both peers; claude-1 r6 adopted the peers' staged
   write-through rename; the restructure branch then closed by round-6 convergence.
4. **"Replay blocked by destination-exists failure" (zcode-1 r5).** Two concordant non-owner
   `WRONG` verdicts (claude-1 r6, kimi-1 r6; `PRIMARY`: `verification.go:192-193`,
   `reservation_recovery.go:79`, `parent_recovery.go:319-328`); withdrawn by its owner in r6
   after self-verification. The deterministic-stage condition repairs the hole the claim missed.
5. **"Shape-only" change to partial-publication (kimi-1 r5 self-correction 2, second half).**
   Verdicted `WRONG` as scoped by claude-1 r6; withdrawn by its owner in r6. The no-clobber half
   stands confirmed.
6. **Containment equivalence (zcode-1 r4 "the conversion does not trade away `os.Root`
   containment").** Verdicted `WRONG` by claude-1 r5 (`filepath.Clean` is lexical-only, "purely
   lexical processing"; `os.Root` enforces reparse-resistant handle-relative operations),
   unrebutted; narrowed by its owner in r5 ("except the final write-through move itself");
   superseded as a deck-adoptable position in r6 by the owner's authority reservation. All three
   round-06 positions now hold the rooted conversion owner-reserved.
7. **`N6` (`MoveFileEx` moves directories; write-through sentence).** Owner-asserted by claude-1
   r4, `UNVERIFIED` through round 4; non-owner `CONFIRMED` by kimi-1 and zcode-1 in round 5 from
   their own fetches; basis layers (including the flag-table contrapositive limitation) recorded
   in §2.
8. **`N7`–`N10` (round 6).** Owner-asserted by claude-1 r6; each `CONFIRMED` non-owner by kimi-1
   r6 and zcode-1 r6 with one-command checks attached; N9 re-verified by the drafter this
   session (both gate files exist; sizes and frontmatter as stated).

**§15.3 dependency check:** no acceptance criterion in this consensus rests on any formerly
disputed claim's contested reading; each closure above is by weakening or by non-owner verdict
on re-executed evidence. `consensus.md` therefore carries no open `DISPUTED` item; FINAL must
repeat this dependency check when it cites any of the above.

## Alternatives disposition

Round-1 mechanism alternatives (§15.6(a) inventories; kimi-1's formal ids retained) and the
design alternatives from rounds 3–6, each with its decisive reason:

- **ALT-1 — ACL implementation: ADOPT** the direct `golang.org/x/sys v0.36.0` chain (zero new
  deps; full symbol map verified). **REJECT** `hectane/go-acl` (re-encodes the 0700 model being
  abandoned; new dep) and `icacls` (localizable text, fragile, extra process, TOCTOU).
- **ALT-2 — Locking: ADOPT** the in-tree `lock_windows.go` LockFileEx implementation plus
  share-mode/handle discipline. **REJECT** `gofrs/flock` and named mutexes (no lock-acquisition
  library fixes share-mode conflicts; the in-tree design is deliberate and working).
- **ALT-3 — Process scope: ADOPT P-A** (truthful liveness, ~40 lines, fail-closed). **REJECT
  P-B in-idea** (Job Objects/attribution would silently loosen the full-argv kill-safety
  calibration); named follow-up.
- **ALT-4 — ACP drain test: ADOPT the two-test split.** **REJECT** in-`Wait`
  `cmd.Stderr=io.Writer` draining for ACP (the observer must stream non-blockingly during the
  live session) and **REJECT** retuning/chunked-writes (only moves the capacity coupling).
- **ALT-5 — Executable fixtures: ADOPT** test-binary re-exec (in-repo precedents). `.bat`/PATHEXT
  only where file discovery is the subject; runner-image `sh.exe` rejected as a fixture
  dependency.
- **ALT-6 — CRLF: ADOPT** per-invocation `-c` pinning as the load-bearing layer + repo
  `.gitattributes` second + parser CRLF tolerance under the split rule. **REJECT**
  workflow-flag-only (does nothing for user checkouts the drift guard compares against).
- **ALT-7 — Gate names: ADOPT** universal encoding inside `GatePath` (PathEscape base + `:`
  escape + device/dot/space guards) with legacy read-fallback + shadow retirement. **REJECT**
  full-hash names (loses human-facing HITL value), per-OS naming (decks cross OSes), and
  refusing pipelines on Windows (weaker product behaviour without security justification).
- **ALT-8 — Coverage: ADOPT** hosted windows-latest x64 + ubuntu cross-compile guard. **REJECT**
  emulated/cross runners and WSL/Wine as execution evidence; **REJECT** a new ARM64 gate
  (build-only assets with explicit no-executed-evidence note).
- **ALT-9 — `O_SYNC`-in-place for create-entry durability: REJECT.** Write durability is real
  and retained where a publication open exists; entry coverage is neither documented nor
  documented-against (N5), the branch had no error path, and its owner withdrew it. The
  deterministic-stage restructure carries both properties.
- **ALT-10 — `FlushFileBuffers` on directory handles: REJECT** as a mechanism (no documented
  directory subject; `GENERIC_WRITE` precondition unmet on our read-only handles; in-tree
  `replace_windows.go:11-12` says so). Retained as diagnostic only.
- **ALT-11 — Handle-based rename (`SetFileInformationByHandle`/`FileRenameInformation[Ex]`):
  REJECT** — carries no write-through flag; adopting it would repeat the inference error N5
  closed. This closure is why containment had no free technical exit.
- **ALT-12 — Documented-nil `SyncDir` / documented-weakening floor: REJECT** (withdrawn r3/r4;
  no silent weakening, no fourth branch).
- **ALT-13 — Restructure everywhere including rooted sites (path-based `MoveFileEx` + lexical
  sibling check): NOT ADOPTED.** A real weakening of `os.Root` containment (TOCTOU-exposed,
  reparse-resolving) that the owner reserved to an explicit deviation; carried fully specified
  and **unapproved** (Open items 1). It is not an implementable automatic fallback, and baseline
  implementation is not conditioned on it.
- **ALT-14 — Randomized stage names at class A: REJECT** (loses the replay blocker on the
  crash-before-move interval; A3 is Windows-reachable and budget-charging). Deterministic
  stages adopted. Randomized stages remain correct at C, where the in-tree pattern already uses
  them on every OS and refusal precedes the rename under the selected branch.
- **ALT-15 — REPLACE at C1/C2 inherited from the `ReplaceSyncedFile` contract: DEFERRED to the
  deviation branch** (under the selected branch C1/C2 refuse). If that branch ever activates,
  the flag is a deliberate, disclosed choice — drafter preference no-REPLACE — requiring
  participant signoff, not a pattern-match inheritance.
- **ALT-16 — Windows-leg-only `-json` census: REJECT/SUPERSEDED** by the all-legs variant with
  per-leg baseline enforcement (§6) — the only shape that reviews "no new suppression" on the
  still-green legs.

## Correlated agreement (§15.6(b))

Unanimity among these three related models is a **shared prior, not independent evidence**, and
this round's convergence is *sequentially produced by the owner's own design*: round 6 ran one
participant at a time, each reading the prior positions — claude-1 joined the class-A
restructure after reading both round-05 files; kimi-1 joined the plain/rooted split after
reading claude-1's round-06; zcode-1 joined both after reading both round-06 files. The deck
holds no dissenting position on any of the three conflicts, and FINAL must say so rather than
read the convergence as independent confirmation. **What would make the agreed position wrong:**
a non-NTFS runner volume; a directory move rejected on the runner; no-replace clobbering an
existing destination on the runner; the deterministic-stage suffix colliding with a valid final
name (prevented by requirement (a) of §2, pinned by unit test); a hosted leg staying red on any
of the 14 historically failing packages; or an ACL invariant the `runneradmin` context cannot
exercise. The owner declining the deviation does **not** make the agreed position wrong — the
refusal branch stands as a reviewed outcome; it would only re-scope what Windows users get.

## Signoffs

<!-- Each agent APPENDS their own signoff block. Do NOT edit others' blocks.
     Expected signers: claude-1, kimi-1, zcode-1 (codex-1 is organizer only and never signs).
     The drafter does not sign in the draft; every signoff is appended by its own author
     in a follow-up commit, per §4 Phase 3. -->

### Signoff: claude-1 — 2026-09-25
Status: ✅ ACCEPT
Notes: Accept as drafted: the consensus states my own round-06 position without distortion — the plain/rooted split, named pre-mutation blocking refusals at the seven rooted sites pending an explicit owner deviation, deterministic non-removed class-A stage names, and the 15-ID / 16-runtime-operation table on kimi-1's IDs — so I hold no design reservation; three FINAL-content items follow, none of which changes a decision or an acceptance criterion.
Protocol attestation: {"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}; packet re-verified by `shasum -a 256` this session against the file named in my task.
Completeness: every Phase-3 section is present, plus §15.5 role-concentration and `## Drafter position changes`, §15.3 `## Verdict conflicts` with its dependency check, §15.6(b) correlated agreement and §15.6(c) `## Alternatives disposition` (ALT-1..ALT-16, each with a decisive reason). The drafter claim in `inbox/zcode-1-to-all_windows-portability_drafter.md` predates every signature. Peer positions are the current round-06 ones, not superseded ones: kimi-1 and zcode-1 both adopt the split and the deterministic stage, and zcode-1's round-06 accepts my testable blocker in terms, which closes it.
Re-verified independently at HEAD 3cf0068 (`git diff --name-only 6b87cf8..HEAD -- '*.go'` = 0 files, so every round-05/06 locator transfers): census 104 pattern lines, 68 `t.Skip(` + 35 `t.Skipf(` + 0 `t.SkipNow(` = 103 call sites, 56 files; helper-call grep 21 lines − 6 definitions − `reservation_recovery.go:67` + `parent_recovery.go:336` = 15 sites, B1 ×2 = 16 operations; the N10 `func` lines reproduce exactly (132/183/267/295/532/55); C1 renames at `parent_recovery.go:333` and syncs at `:336` through the function value bound at `:311` with `defer dir.Remove(stage)` at `:328`, C2 `:570`/`:573`/`:565`, so the pre-rename refusal ordering is right; A3's publish path builds `intentName` ungated at `:75` while the read path gates at `:92-93`; the invariant comments at `verification.go:192-193` and `reservation_recovery.go:79` read as quoted; the two POSIX gates sit at `trajectory_verify.go:106-108` and `verification.go:251-253`; both release-gate files exist (12075 bytes, 2026-09-25 11:13; 6412 bytes). Per §15.1 these are transcription checks of the draft against the tree, not verdicts on claims I own (N5–N10 and the census).
FINAL correction 1 — provenance. `## Verdict conflicts` item 8 states that N7–N10 were "each CONFIRMED non-owner by kimi-1 r6 and zcode-1 r6". `round-06/zcode-1.md` contains no occurrence of N7 and issues no verdict on it; N7 carries exactly one non-owner CONFIRMED (kimi-1 r6: PRIMARY code side, PRIMARY(r4) for the FlushFileBuffers subject list, SECONDARY on zcode-1's round-05 verdict for the GENERIC_WRITE precondition). N7's status is unchanged — §15.3 forbids resolution by counting, and no acceptance criterion rests on it — so this is a corroboration-breadth error, not a status change. FINAL must state the single round-06 verdict with its layered tags. N8–N10 do carry both peers' round-06 verdicts as claimed.
FINAL correction 2 — quotation. The `## Agreed decisions` preamble attributes the quoted phrase "stands unreopened" to all three round-06 files; that string appears nowhere in this idea except `consensus.md` itself. The substantive claim is true (my round-06 §4 records the settled items as "Unchanged and reopened by nobody", and both peers carry equivalents), so FINAL should drop the quotation marks or quote each file's actual wording. Raised because this is the sixth error of the transcription class in this idea and the deck's own lesson is to generate, not transcribe.
FINAL correction 3 — honest framing of the selected branch's benefit side, for the owner's deviation decision. The table's own columns entail, but the prose never states, that under the refusal branch no currently Windows-reachable publication gains a new durability mechanism: of the 9 reachable rows, 8 refuse and D1 is an already-satisfied derived no-op, while the three restructuring sites A1/B1/B2 are all dormant behind `trajectory_verify.go:106-108` and are therefore design-of-record for a gate that has not lifted. The refusal branch still takes nothing away and remains the strongest position the deck may take on its own authority — but FINAL and the deviation question should say the benefit plainly rather than leave the owner to derive it from the reachability column.
Scope of this signature: the refusal branch only. ALT-13 / the weaker-containment alternative remains UNAPPROVED and is not accepted by me here; silence is not approval; if the owner grants the deviation, the converted contract needs its own participant signoff before implementation, and baseline implementation must neither wait for nor presume that answer. I am not the drafter (Drafter: no); I request no seventh round, no phase advance, no finalization and no deviation. I remain a Phase-6 reviewer with kimi-1 and will file corrections 1–3 as review findings if they reach FINAL uncorrected.

### Signoff: kimi-1 — 2026-09-25
Status: ✅ ACCEPT
Notes: Accept as drafted: the consensus states my round-06 position without distortion — the plain/rooted split with named pre-mutation blocking refusals at the seven rooted sites pending an explicit owner deviation, claude-1's deterministic-stage condition at class A (which I adopted in round 6), and the 15-ID / 16-runtime-operation table on my row-ID scheme — so I hold no design reservation.
Protocol attestation: {"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}; matches the packet named in my task.
§15.5 check: `## Drafter position changes` is accurate and complete against `round-06/zcode-1.md` — all three recorded self-corrections (destination-exists replay withdrawal, "restructure everywhere" withdrawal as deck-adoptable, adoption of deterministic stage / plain-rooted split / 15-16 table) appear in that file as stated; this signoff ratifies it.
Independent verification of claude-1's FINAL corrections 1 and 2 (re-executed by me this session, `PRIMARY` for the file-content facts): `round-06/zcode-1.md` contains no occurrence of N7 and issues no verdict on it, so `## Verdict conflicts` item 8 overstates N7's corroboration breadth — N7 carries exactly one non-owner CONFIRMED (mine, round 6: `PRIMARY` code side, `PRIMARY(r4)` for the FlushFileBuffers subject list, `SECONDARY` on zcode-1's round-05 verdict for the GENERIC_WRITE precondition), and §15.3 forbids resolution by counting, so this is a transcription error, not a status change; and the quoted string "stands unreopened" appears nowhere except `consensus.md` itself — my round-06 line 113 and zcode-1's line 92 both wrote "stands un reopened", and claude-1's round-06 used different wording — so FINAL must drop the quotation marks or quote each file's actual wording. I concur with claude-1 that both are FINAL-content corrections that change no decision and no acceptance criterion, and I join correction 3 (the refusal branch's benefit side stated plainly for the owner's deviation decision: of the 9 Windows-reachable rows, 8 refuse and D1 is an already-satisfied no-op; A1/B1/B2 restructure behind a gate that has not lifted) as an honest-framing requirement for FINAL, consistent with my round-06 statement of the reachable-today cost.
Census and inventory: the consensus numbers (104 pattern lines, 103 call sites = 68 `t.Skip(` + 35 `t.Skipf(`, 56/55 files; 15 sites via the function-value call at `parent_recovery.go:336`; B1 ×2 = 16 operations) are identical to what I re-executed at `9e55050` in round 6, and claude-1's signoff re-verifies them at `3cf0068` — no drift between the two HEADs on any `.go` file.
Scope of this signature: the refusal branch only. ALT-13 / the weaker-containment alternative remains UNAPPROVED and is not accepted by me here; silence is not approval; if the owner grants the deviation, the converted contract — with my guard set (1)–(4) including the publish-time `validHash` gate at A3 — requires its own participant signoff before implementation, and baseline implementation neither waits for nor presumes that answer. I am not the drafter (Drafter: no); I request no seventh round and no deviation myself. I remain a Phase-6 reviewer with claude-1 and will file FINAL corrections 1–3 as review findings if they reach FINAL uncorrected.
Counter-proposal (required if ❌): n/a — status is ✅ ACCEPT.

### Signoff: zcode-1 — 2026-09-25
Status: ✅ ACCEPT
Notes: Accept as drafted: the consensus states my own round-06 position without distortion — the class-A staged write-through rename with the deterministic-stage condition I adopted, carrying my two drafter requirements ((a) the stage derivation provably non-colliding with every valid final name, pinned by unit test; (b) the publish-time `validHash` gate at A3, where `reservation_recovery.go:75` publishes ungated while the read path gates at `:92-93`); the plain/rooted split with named pre-mutation blocking refusals at the seven rooted sites pending an explicit owner deviation, C1/C2 refusals ordered before the rename (N8: the barrier fires after publication at `:336`/`:573`); the 15-ID / 16-runtime-operation table on kimi-1's row IDs with my `fsutil.SyncDir` named-type contract as the audit of record; the §5 settled items, the §6 census, the §7 coverage envelope, and the version anchor above 1.50.0 — so I hold no design reservation.
Protocol attestation: {"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}; matches the attestation supplied with my launch task.
Drafter duties: `## Drafter position changes` = "None since `round-06/zcode-1.md`" is accurate — this consensus introduces no position not held in that file; the §15.5 role-concentration line is accurate (I hold no facilitator role in this idea).
Corrections 1 and 2 are accepted as self-corrections on my own transcription, verified this session (`PRIMARY` for the file-content facts): `round-06/zcode-1.md` contains no occurrence of N7 — I confirmed N8, N9 and N10 by name and issued no verdict on N7 — so `## Verdict conflicts` item 8 overstates N7's corroboration breadth, and FINAL must carry N7's single round-06 non-owner verdict (kimi-1's, with its layered tags: `PRIMARY` code side, `PRIMARY(r4)` for the FlushFileBuffers subject list, `SECONDARY` on my round-05 verdict for the GENERIC_WRITE precondition) while N8–N10 keep both peers' verdicts as stated; and my round-06 line 92 wrote "stands un reopened", so the quoted "stands unreopened" is a smoothing that appears in no participant file — FINAL drops the quotation marks or quotes each file's actual wording. Both are weakenings, effective immediately; neither changes a decision, an acceptance criterion, or any §15.3 status, and no acceptance criterion rests on either.
Correction 3 is joined, and it names a duty that was mine to carry: my round-06 stated the reachable-today cost plainly, and FINAL must state the benefit side just as plainly — under the refusal branch no currently Windows-reachable publication gains a new durability mechanism: of the 9 reachable rows, 8 refuse and D1 is already satisfied via `writeState` (`state.go:197-221`) ending in `ReplaceSyncedFile`, while the three restructuring sites A1/B1/B2 are design-of-record behind the dormant `trajectory_verify.go:106-108` POSIX gate. The branch takes away nothing that works today and is the strongest position the deck may take on its own authority; FINAL and the deviation question say so in that many words rather than leaving the owner to derive it from the reachability column.
Scope of this signature: the refusal branch only. ALT-13 / the weaker-containment alternative remains UNAPPROVED and is not accepted by me here; my round-06 recorded preference for granting the deviation is input to the owner's decision only — not a request, and silence is not approval. If the owner grants it, the converted contract (kimi-1's guard set (1)–(4) including the publish-time gate at A3, claude-1's precise statement of what is given up, my hosted rooted-escape test, the disclosed C1/C2 flag choice) requires its own participant signoff before implementation; baseline implementation neither waits for nor presumes that answer. I am the drafter (Drafter: yes) and the assigned single implementation owner; corrections 1–3 are binding FINAL drafting requirements on me, and I accept the peers' standing intent to file them as review findings if they reach FINAL uncorrected. I request no seventh round. With this block all three participants have signed ✅; any further phase step is the organizer's to process per protocol, not mine to advance.
Counter-proposal (required if ❌): n/a — status is ✅ ACCEPT.
