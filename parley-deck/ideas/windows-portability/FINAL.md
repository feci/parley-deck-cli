---
idea: windows-portability
status: final
author: zcode-1
consensus-date: 2026-09-25
participants: [claude-1, kimi-1, zcode-1]
facilitator: codex-1
---

# FINAL — windows-portability

Genuinely correct CLI behavior on native Windows: real Windows implementations where the deck may
select them, and explicit, reviewed, user-visible refusals everywhere it may not — never a disabled
privacy check, never a skipped test standing in for product behavior. This document is
self-contained: a fresh implementer needs no prior rounds, transcripts, or consensus sections;
everything binding is restated here.

## Protocol context attestation (full)

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet: `.parley-runtime/protocol-packets/full-phase4-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md`,
re-verified by `shasum -a 256` this session (hash above). Phase 4 FINAL drafting. Drafter:
zcode-1 (`Drafter: yes` in own signoff); quorum exactly claude-1, kimi-1, zcode-1; codex-1 is
organizer only and never implements, verifies code, authors participant verdicts, or signs off.
Original owner quotations are given in the owner-supplied English translations only; the
originals were answered in Slovak and are cited as such without reproduction.

**Drafter quorum verification (executed this session, recorded as duty):**
`parley consensus status windows-portability` reports Consensus: ready with signoffs claude-1
✅ ACCEPT, kimi-1 ✅ ACCEPT, zcode-1 ✅ ACCEPT — canonical quorum present, every acceptance
unreserved. All participant processes exited 0 (round 6: claude-1 16:41:25 UTC, kimi-1 16:59:00
UTC, zcode-1 17:06:07 UTC). `consensus.md` is signed and is never rewritten by this FINAL; where
this document restates consensus content it does so under the three binding corrections recorded
below.

**Drafter re-verification at HEAD `35b57c3`** (all commands re-executed this session):
`git diff --name-only 6b87cf8..HEAD -- '*.go'` = 0 files, so every round-05/06 locator transfers.
Skip census: **104** `t\.Skip` pattern lines under `internal`; **103 call sites** = 68 `t.Skip(` +
35 `t.Skipf(` + 0 `t.SkipNow(`; **56 files** matching / **55 with call sites** (sole non-call
`internal/runner/launch_precheck_skip_test.go:58`). Inventory: the name-keyed grep over the six
barrier helpers yields 15 call lines; minus `reservation_recovery.go:67` (inside `syncIntent`'s
own body) = 14 name-visible; plus `parent_recovery.go:336` (invoked through a function value
bound at `:311`, invisible to name-keyed grep) = **15 independent sites**; B1 fires twice per run
(`:117-128`) = **16 runtime operations**. Helper `func` lines reproduce exactly:
`trajectory_verify.go:132`, `verification.go:183`, `unchanged.go:267`, `parent_recovery.go:295`,
`:532`, `reservation_recovery.go:55`. **The table in §B is generated from this enumeration,
never transcribed** (the N10 lesson: all three rounds' files mislabeled helper `func` lines — the
fifth census-class error; wherever a number is load-bearing it is generated mechanically).

## Binding FINAL corrections from the signoffs (all three carried, verbatim scope)

All three signoffs attach FINAL-content corrections; all three are binding on this drafter and are
applied throughout this document:

1. **N7 provenance (corroboration breadth).** `consensus.md` §Verdict conflicts item 8 overstated
   N7's corroboration. **N7 carries exactly one round-6 non-owner CONFIRMED — kimi-1's** (round 6,
   layered tags: `PRIMARY` code side, `PRIMARY(r4)` for the FlushFileBuffers subject list,
   `SECONDARY` on zcode-1's round-05 verdict for the `GENERIC_WRITE` precondition);
   `round-06/zcode-1.md` contains no occurrence of N7 and issues no verdict on it. N7's §15.3
   status is unchanged — §15.3 forbids resolution by counting, and **no acceptance criterion in
   this FINAL rests on N7**. N8, N9 and N10 do carry both peers' round-06 verdicts as stated.
2. **Accurate quotations.** The consensus preamble attributed a quoted phrase to all three
   round-06 files that appears in no participant file (the string exists only in `consensus.md`
   itself). The substantive claim is true and is restated here without the fabricated quotation:
   every settled item in §D below carries all three participants' round-06 assent — kimi-1
   round-06 (line 113) and zcode-1 round-06 (line 92) each wrote the settled items "stands un
   reopened"; claude-1 round-06 §4 recorded them as "Unchanged and reopened by nobody". Wherever
   this FINAL cites participant wording, it quotes the file's actual wording.
3. **The refusal branch's benefit side, stated plainly (for the owner's deviation decision).**
   Under the selected refusal branch, **no currently Windows-reachable publication gains a new
   durability mechanism**: of the 9 reachable rows, **8 refuse** (A3, B5, C1, C2, D2, D3, D4a,
   D4b) and **D1 is an already-satisfied derived no-op** (its original publication `writeState`
   `state.go:197-221` already ends in `fsutil.ReplaceSyncedFile`, the reviewed write-through
   replace), while the three restructuring sites A1/B1/B2 are **dormant** behind the existing
   reviewed POSIX refusal at `trajectory_verify.go:106-108` — design-of-record for a gate that
   has not lifted. The branch still **takes away nothing that works today** (trajectory
   verification and captured journals were already refused on Windows) and is the strongest
   position the deck may take on its own authority; the deviation question (§P.1) presents both
   branches' costs side by side so the owner does not have to derive this from the reachability
   column.

## Purpose / user-visible outcome

On native Windows (hosted windows-latest x64 evidence only — §F):

- The CLI **builds and its test matrix runs unfiltered and green** on Windows, macOS and Ubuntu
  on the final released commit, with no new suppression (§E) and only individually reviewed,
  truly-inapplicable exclusions.
- **Snapshot privacy is real**: stores are created owner-only (protected DACL) and verification
  refuses any non-owner grant, naming the trustee (§D.1).
- **Gate names work cross-OS**: universal encoding plus a raw-ID allowlist retire the
  `ERROR_INVALID_NAME` failure class (§D.4).
- **Honest durability**: every publication operation on Windows either carries a proved, scoped
  mechanism or refuses — named, pre-mutation, user-visible, blocking (§A–§C). Two
  Windows-reachable features (precharge reservation-intent publication; parent-recovery and
  recovered-parent apply) refuse rather than weaken `os.Root` containment, pending an explicit
  owner deviation that is **unapproved and nonblocking**.
- **Honest process semantics**: truthful liveness, fail-closed; durable-kill/attribution refuse
  honestly rather than pretend (§D.2). Missing `sh` is a named blocking refusal with an
  actionable prerequisite (§D.3).
- The **Windows "experimental" label is removed if and only if** the released commit's hosted
  Windows leg is green; the held winget PR ships then and only then (§O).

## Context & orientation

- **Repository/worktree:** this `windows-portability` worktree and branch, base commit
  `868825f20bd0988abde1a6c38b9efe405c5123d1`; no development PRs (owner override of deck PR
  mechanics: files canonical, direct integration, main linear); do not mutate other active
  worktrees or main during development. Existing reviewed Linux repairs in the base are
  preserved.
- **Evidence shared with all participants:** `source-context/README.md` with
  `release-ci-revalidation-claude-1.md` and `release-repair-plan-zcode-1.md` (the latter a
  participant proposal, not a decision). Hosted run IDs: **35987916696** (also failed Ubuntu —
  keep per-run claims separate) and **36009912946** (14 red vs 16 ok packages). Reported Windows
  failure classes: internal/evidence test build (`syscall.Mkfifo` undefined), ACP drain guard,
  snapshot privacy/ACLs, invalid gate filenames, cross-process locks, fixture portability.
- **Round history:** rounds 1–3 (2026-09-24/25) merged the six-item tri-participant contract;
  round 4 opened the durability conflict; round 5 (owner-authorized, focused) left it unresolved;
  round 6 (owner-authorized, **sequential** — claude-1, then kimi-1, then zcode-1, each reading
  prior positions) converged all three positions on all three conflicts. No seventh round is
  authorized. The round-6 convergence is sequentially produced by the owner's own design and is
  **not independent evidence** (§N). Consensus drafted by zcode-1 and signed unreserved ACCEPT by
  all three participants; every signoff scopes itself to the refusal branch only.
- **Roles:** zcode-1 is drafter and **single implementation owner** (claim + worktree mapping in
  `IMPLEMENTATION.md` before any code edit; one implementation owner at a time); claude-1 and
  kimi-1 are the independent reviewers; codex-1 organizes only. Channel verification after
  release is by a participant, never codex-1.
- **Owner authority boundary (binding):** any proposed weakening of `os.Root` containment "must
  come back to the owner as an explicit deviation" (owner-supplied translation, original in
  Slovak, round-6 answer). The deck cannot grant it to itself. **Silence is not approval.**

## Final plan / specification

### A. Selected design — the refusal branch, preserving `os.Root` containment

The original owner scope authorizes "a real Windows implementation **or** an explicit, reviewed,
user-visible refusal". The selected design is the **plain/rooted split**, the strongest design the
deck may select on its own authority:

- **Plain sites restructure** (no `os.Root` to lose): **A1, B1, B2** — all currently dormant
  behind the existing reviewed POSIX refusal at `trajectory_verify.go:106-108`, so this is the
  design of record if that gate ever lifts. **D1 derives** (its original publication `writeState`
  `state.go:197-221` already ends in `fsutil.ReplaceSyncedFile`).
- **Rooted sites take the named blocking refusal, pending an explicit owner deviation that has
  not been granted**: **A2, A3, B3, B4, B5, C1, C2**; refusals at C1/C2 are **ordered before the
  rename** (N8: the barrier fires after publication at `parent_recovery.go:336`/`:573`, so the
  refusal must be evaluated before `:333`/`:570`, leaving `defer dir.Remove(stage)` `:328`/`:565`
  to clear the stage with nothing published). **D2/D3/D4a/D4b follow their originals into
  refusal** (the coherence rule: a class-D derivation holds only if its original carries a proved
  mechanism).
- The refusal is **named, pre-mutation, user-visible, and blocking** — never a silent skip, never
  a green-API-call-as-guarantee, never an implicit no-op. There is no fourth branch.
- **The weaker-containment alternative (ALT-13) is carried as UNAPPROVED, not as an
  implementable fallback.** The organizer's nonblocking deviation request is
  `inbox/codex-1-to-user_windows-portability_containment-deviation.md`; its full specification is
  §P.1. No implementation may select it without the owner's explicit approval **and** participant
  signoff on the resulting contract; silence is not approval. If the owner declines, the
  refusal-branch table is final as written — a reviewed outcome, not a failure. **Baseline
  implementation neither waits for nor presumes the answer.**
- All conversions are **Windows-confined**; the Unix path at every row is byte-identical to today.
  Windows `fsutil.SyncDir` remains the fail-closed named-refusal emitter for any site not
  carrying a proved mechanism, including future sites added without conversion — under the
  refusal branch it is also the rooted sites' refusal emitter.

**Durability mechanism for the sites that do convert (class A): staged write-through rename with
deterministic stages.** Stage create (`O_EXCL` + mandatory file fsync) then
`MoveFileEx(MOVEFILE_WRITE_THROUGH)` without `MOVEFILE_REPLACE_EXISTING`, with the
**deterministic-stage condition**: the stage name is a fixed derivation of the final name,
created with `O_EXCL`, and **not removed on the error path**, so an interrupted publication
leaves a visible blocker whose existence prevents replay — exactly the invariant
`verification.go:192-193` ("its existence prevents replay") and `reservation_recovery.go:79` ("a
partial publication stays visible and cannot be overwritten by retry") supply today. The final
name appears only via the atomic move, so readers never see torn bytes; no-REPLACE prevents
clobber of a completed publication; concurrent publishers race on the stage `O_EXCL` exactly as
today they race on the final-name `O_EXCL`.

Two drafter requirements carried with the design: **(a)** the stage derivation must be provably
non-colliding with every valid final name (final names are validated or product-constructed, so
the stage suffix must sit outside every valid-final grammar, pinned by a unit test); **(b)** the
publish-time `validHash` gate at A3 (`reservation_recovery.go:75` builds
`intentName(i.Accounting.EntryKey)` ungated, while the read path gates at `:92-93`; under any
restructure the key feeds both stage and final names).

**Documentation basis, honestly tiered:** the write-through sentence — "does not return until the
file is actually moved on the disk" — is scoped to the exact operation (PRIMARY fetches by all
three participants in round 5, independently); there is **no** documentation scoped to the
parent-directory entry of a create (N5, CONFIRMED non-owner by kimi-1 and zcode-1 in round 5
after claude-1's withdrawal; coverage is "neither documented nor documented-against");
`O_SYNC`→`FILE_FLAG_WRITE_THROUGH` is real for **write** durability only; `FlushFileBuffers` on
our read-only directory handles is foreclosed twice over (no documented directory subject;
documented `GENERIC_WRITE` precondition unmet) and stays **diagnostic-only**. The no-REPLACE
existing-destination error basis is disclosed per the N6 limitation: flag-table contrapositive
plus a hosted existence-probe assertion, never upgraded to "documented". The directory-move
guarantee for class B is a disclosed composition (N6's "existing file **or directory**" subject
sentence plus the write-through sentence), labelled as such here, pinned by mechanics probes
only. **No green API call is ever cited as the guarantee; no power-cycle evidence is obtainable
on hosted runners and this document says so.**

### B. The normalized operation table — 15 IDs / 16 runtime operations

The three historical inventories (14, 15, 16) were the same set under different individuation:
14 + (D4 split) = 15 + (B1 second loop iteration) = 16. Row IDs follow kimi-1's scheme (strict
file-then-line order). The audit of record is the `fsutil.SyncDir` named-type contract
(compiler-enforced fail-closed on Windows; never a grep). Abbreviations: **WT-move** =
`MoveFileEx(MOVEFILE_WRITE_THROUGH)`; **no-REPLACE** = without `MOVEFILE_REPLACE_EXISTING`;
**det-stage** = deterministic stage per §A.

| ID | Site (publication → barrier) | Class | Rooted? | Win-reachable today | Required guarantee | Disposition (refusal branch = selected design) | If deviation granted | Feature blast radius if refused | Partial-publication recovery |
|---|---|---|---|---|---|---|---|---|---|
| **A1** | `trajectory_verify.go:146` `os.OpenFile O_EXCL` → `:154` | A create | plain | no — gated `:106-108` | create entry durable; anti-replay | **Restructure**: det-stage `O_EXCL`+fsync → WT-move no-REPLACE | same | none today; inherits if gate lifts | stage persists, blocks replay; final never torn |
| **A2** | `verification.go:202` `Root.OpenFile O_EXCL` → `:211` | A create | rooted | no — gated `:251-253` | same + `:192-193` invariant | **Refuse**, pre-mutation | det-stage → WT-move no-REPLACE + guard set | none today | `:192-193` invariant intact (unconverted) |
| **A3** | `reservation_recovery.go:75` `Root.OpenFile O_EXCL` → `:84`→`:67` | A create | rooted | **yes** | same + `:79` invariant | **Refuse**, pre-mutation | as A2 + publish-time `validHash` gate | precharge reservation-intent publication refuses → `PrepareCycleReservation` fails (`cycle_binding.go:269`, no GOOS gate) | nothing published; `:79` invariant intact |
| **B1** | `trajectory_verify.go:118` `os.Mkdir` → `:125` — **fires ×2** (`:117-128` loop) | B mkdir | plain | no — gated `:106-108` | base-dir entries durable, idempotent | **Restructure**: stage-`Mkdir` + WT-move + EEXIST-tolerant recheck | same | none today | IsExist tolerance preserved; nothing published |
| **B2** | `trajectory_verify.go:190` `os.Mkdir` → `:193` | B mkdir | plain | no — gated `:106-108` | run-dir entry durable | **Restructure** as B1 | same | none today | as B1 |
| **B3** | `verification.go:264` `parent.Mkdir` (IsExist-tolerant) → `:267` | B mkdir | rooted | no — gated `:251-253` | store entry durable, idempotent | **Refuse**, pre-mutation | stage-`Mkdir` + WT-move + tolerant recheck | none today | as B1 |
| **B4** | `verification.go:283` `base.Mkdir(charge)` — anti-replay token `:281-282` → `:286` | B mkdir | rooted | no — gated `:251-253` | reservation entry durable; never reuse | **Refuse**, pre-mutation | det-stage-`Mkdir` + WT-move no-REPLACE; existing dir = loud either way (N6); EEXIST on stage = loud, matching today's `Mkdir(charge)` failure | none today | `Mkdir` failure = "already reserved" `:284` |
| **B5** | `reservation_recovery.go:36` `dir.Mkdir` (IsExist-tolerant) → `:47` | B mkdir | rooted | **yes** | intent-root entry durable, idempotent | **Refuse**, pre-mutation | stage-`Mkdir` + WT-move + tolerant recheck | `openIntentRoot(create=true)` fails → reservation creation unavailable | nothing published; `Mkdir` idempotent |
| **C1** | `parent_recovery.go:333` `dir.Rename` → `:336` (func value bound `:311`) | C rename | rooted | **yes** | rename entry durable; publish-once | **Refuse, ordered before `:333`** (N8) | WT-move; flag chosen deliberately and disclosed (drafter preference: no-REPLACE), never inherited by pattern-match | parent-recovery **apply** refuses | pre-rename refusal ⇒ nothing published; `defer dir.Remove(stage)` `:328` clears stage |
| **C2** | `parent_recovery.go:570` `dir.Rename` → `:573` | C rename | rooted | **yes** | same | **Refuse, ordered before `:570`** | as C1 | recovered-parent **apply** refuses | as C1; stage cleared `:565` |
| **D1** | `unchanged.go:241` → `:279` (func `:267`) | D re-sync | plain | **yes** | already satisfied | **Derived no-op**, argued per site: `writeState` (`state.go:197-221`) already publishes via `ReplaceSyncedFile` = WT-move | same | none | `os.CreateTemp` + `defer os.Remove` (`state.go:205`/`:209`); replace atomic |
| **D2** | `parent_recovery.go:368` → `:308` | D re-sync | rooted ctx | **yes** | inherits C1 | **Refuse** — derives from a refused publication | derives once C1 converts | tracks C1 | re-sync only; never rewrites |
| **D3** | `parent_recovery.go:605` → `:545` | D re-sync | rooted ctx | **yes** | inherits C2 | **Refuse** | derives once C2 converts | tracks C2 | re-sync only |
| **D4a** | `reservation_recovery.go:420` `syncIntent` | D re-sync | rooted ctx | **yes** | inherits A3/B5 | **Refuse** | derives once A3/B5 convert | tracks A3 | re-sync only; `:426` never rewrites bytes |
| **D4b** | `reservation_recovery.go:431` `SyncFile(trajectory.json)` + `syncVerificationDirectory` | D re-sync | rooted ctx | **yes** | inherits A3/B5 | **Refuse** | derives once A3/B5 convert (`trajectory.json` via `writeState`→`ReplaceSyncedFile` already) | tracks A3 | re-sync only |

**Totals:** 15 IDs; 16 runtime operations; **9 of 15 Windows-reachable today** (A3, B5, C1, C2,
D1, D2, D3, D4a, D4b); 6 dormant behind the two existing reviewed POSIX refusals
(`trajectory_verify.go:106-108`, `verification.go:251-253`). Under the selected branch: 3 plain
rows restructure (A1, B1, B2 — dormant), D1 derives, 7 rooted rows refuse pending the owner
deviation, 4 class-D rows follow their originals. Per correction 3: of the 9 reachable rows, 8
refuse and D1 is already satisfied — no currently reachable Windows publication gains a new
durability mechanism under this branch.

### C. Feature-level refusals (the honest, user-visible cost on Windows)

Under the selected branch, on Windows:

1. **Precharge reservation-intent publication** refuses (A3 + B5, dragging D4a/D4b):
   `openIntentRoot(b, true)` ← `PrepareCycleReservation` (`reservation_recovery.go:126`) ←
   `cycle_binding.go:269`, no GOOS gate. Budget reservation intents cannot be published.
2. **Parent-recovery apply** (C1 + D2) and **recovered-parent apply** (C2 + D3) refuse: the
   recovery publication commands are unavailable.
3. **Already refused today and unchanged** (the refusal branch takes away nothing that currently
   works): trajectory verification (`trajectory_verify.go:106-108`) and captured-verification
   journals (`verification.go:251-253`) carry existing reviewed POSIX refusals on Windows.
4. **Blocking missing-`sh` refusal** (§D.3) and the durable-kill/attribution refusals under P-A
   are pre-existing, designed, user-visible refusals — first-class outcomes under the owner
   scope, each with actionable refusal text.

Every refusal is named, its trigger is before the mutation, its message is user-visible and
actionable, and it is covered by a hosted test pinning the refusal path. This FINAL enumerates
every refused operation with its blast radius and recovery column in the §B table — never a
headline count.

### D. Settled non-durability requirements (rounds 1–3, tri-participant; all three round-06 files record the settled items as standing un-reopened — kimi-1 r06 line 113 and zcode-1 r06 line 92 wrote "stands un reopened"; claude-1 r06 §4 recorded them as "Unchanged and reopened by nobody")

**D.1 Snapshot privacy/ACL (W-PRIV).** Creation sets an **owner-only protected DACL**: exactly
one access-allowed ACE for the token user (`GetCurrentProcessToken().GetTokenUser()`) in a
protected, non-inherited DACL (`PROTECTED_DACL` load-bearing), via
`ACLFromEntries`/`SetNamedSecurityInfo`, all in the already-direct `golang.org/x/sys v0.36.0`
(zero new dependencies; `hectane/go-acl` and `icacls` rejected). Verification walks effective
granted trustees via `GetNamedSecurityInfo` and refuses any non-owner grant, inherited or
explicit, **naming the trustee**; the refusal string is sentence-stable ("snapshot store must be
a private real directory" — the 71× hosted signature) with the trustee appended. FAT/exFAT/no-ACL
volumes refuse, never weaken. Pre-existing stores: **refuse-and-instruct** (move aside,
re-create under the protected DACL; explicit user-invoked documented repair named in the refusal
text) — **no in-place DACL rewrite of user data in this idea**. Adversarial hosted tests:
`BUILTIN\Users` grant → refuse; `Administrators` grant → refuse; inherited-only → refuse; own
creation → pass; create-then-requery readback self-check. `!IsDir()` and reparse/symlink
rejection kept on both platforms; Unix behaviour byte-identical. (kimi-1's owner+SYSTEM
preference is recorded as a minority note, withdrawn r3.)

**D.2 Process liveness P-A; P-B deferred.** Truthful `alive` via
`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` + `GetExitCodeProcess`, `STILL_ACTIVE` sentinel
caveat documented, **fail closed on `ERROR_ACCESS_DENIED`** (unverifiable is alive for refusal
purposes); `bootID` stays `""` (scoped null); `supportsDurableKill()` stays false with
Windows-honest durable-kill/attribution refusal strings; `proclive_windows.go` comment corrected
(x/sys is direct at `go.mod:10`). Fixes the four broken callers (`durablekill.go:29`/`:63`,
`verification.go:142`, `reviewsnapshot.go:244`). **P-B** (Job Objects
`JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `GetProcessTimes` attribution,
`QueryFullProcessImageName`) is a **named follow-up idea** with its own adversarial review: the
Windows image-path is coarser than the Linux full-argv calibration of `procctl.commandMatches`,
and adopting it in-flight would silently loosen a kill-safety gate. The `mvdan.cc/sh/v3`
supervisor port is likewise a separate reviewed idea (new dependency executing untrusted
`checks:` strings).

**D.3 Blocking missing-`sh` refusal (W-SHELL).** On Windows, `checks:`
(`driver_impl.go:266`) and captured verification (`evidence/execute.go:118`/`:122`,
`criterionSupervisor` `:100-114`) call `exec.LookPath("sh")` **once, before any work**; on
failure they return a named, actionable, blocking refusal naming **Git for Windows** as the
documented prerequisite (template: name the context, name the prerequisite, "refusing rather
than executing unverified"). The gate never passes by default; LE-4 stays intact; exits
non-zero; recorded in the run. Hosted tests pin both paths (`sh` forced off `PATH` → message +
non-zero; `sh` present → normal path) — because GitHub's windows-latest ships Git for Windows
`sh.exe` on PATH, an untested leg would only prove the image happened to include it.

**D.4 Universal gate-name encoding + raw-ID charset allowlist (W-NAME).** One universal encoding
on **every OS**, applied inside `GatePath` (`gate.go:104-106`) so the `.tmp` staging name
inherits it by construction (the 9–10× hosted `ERROR_INVALID_NAME` failures were on the staging
name): `net/url.PathEscape` base plus explicit `:` escape plus reserved-device-name and
trailing-dot/trailing-space guards (`block-a->block-b` → `block-a-%3Eblock-b.gate.json`).
Per-OS naming rejected (decks cross OSes; a Unix-written `>` name cannot exist on Windows).
Legacy `a->b.gate.json` read-fallback for one release; shadow retirement on first rewrite (no
stale HITL answer survives); mixed-name deck test; the repo has zero legacy gate names today.
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

**D.5 `strict_gate` branch-independent fix; cause undiagnosed.** The hosted `:179` failure
(product `impl.go:448-472` `reviewRoundHasFindings` returned no-veto on
file-where-directory-belongs) contradicts both source-level mechanisms (excluded at go1.27.1 and
the CI-fetched go1.26.8). Shipped regardless: **stat the round path and veto when it exists and
is not a directory** (+ one `Lstat`), a cross-platform pinning test, a hosted probe printing
`runtime.Version()`, the raw `os.ReadDir(regular file)` error and its `errors.Is(err,
fs.ErrNotExist)` classification at the CI toolchain (read the resolved version from the existing
setup-go log first; `go.mod` has no `toolchain` pin), and the sweep of every
`ErrNotExist`-after-directory-read decision. Cause recorded **undiagnosed**; no branch called
confirmed; not reclassifiable as fixture-only without the probe result.

**D.6 File sharing / cross-process locks (W-LOCK).** The budget locks are already correctly
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

**D.7 ACP split; AF_UNIX; CRLF split.** ACP: **split the two drain tests, do not retune** —
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
comparisons (snapshot content, `TreeDigest`); the load-bearing layer is per-invocation
`-c core.autocrlf=false -c core.eol=lf` pinned on parley's own git invocations (in-tree precedent
`reviewsnapshot.go:138-139`; the only layer that follows parley into arbitrary user repos),
repo-scoped `.gitattributes` second as its own reviewed commit; product parsers comparing
embedded-vs-disk bytes are CRLF-tolerant as a product requirement; runner
`git config --show-origin` dump in the probe bundle.

**D.8 Step-1-first build fixes.** The Mkfifo build tag plus the `app_test.go:155/:157` comma-ok
fix (unchecked type assertions abort the whole test binary — why `wait`/`usage` have no Windows
results) land **first, always**; one diagnostic cycle budgeted for whatever they unblock, green
not assumed. The Ubuntu leg gains the `GOOS=windows go build ./... && go vet ./...`
cross-compile guard — compile-only evidence, never execution evidence. `internal/evidence` was
the sole compile failure in the tree; both windows/amd64 and windows/arm64 build exit 0.

**D.9 Fixture-portability sweep with the reconciliation ledger.** Per-test reconciliation ledger
(hosted log ↔ class ↔ stage ↔ exclusion row) emitted before the sweep lands and refreshed after
step-1 fixes change the denominator. Classes and fixes: `#!/bin/sh`/extensionless fixtures →
test-binary re-exec (`os.Executable()` + guard; in-repo precedents) or `cmd.exe /c` for trivial
scripts; `.bat`/`.cmd` via PATHEXT only where file discovery is the test subject; Chmod-0
fixtures → one `denyRead(t, path)` helper (chmod 0 on Unix, DACL deny-ACE on Windows, sharing
the §D.1 module) — "Chmod cannot make a directory unreadable" is not a legitimate exclusion once
ACLs exist; HOME fixtures set `USERPROFILE` too; `/repo` literals → `t.TempDir()` +
volume-aware joins with per-OS refusal assertions kept; TUI temp-perms reformulated via the ACL
helper, never silently blessing 0666 synthesis as "private". Sweep rule: a Windows exclusion is
admissible only when the behavior under test **cannot exist** on Windows — not when
inconvenient; one reviewed row per exclusion, attackable individually; diff-level check that no
`t.Skip` was added without a row; close-time suppression grep (`Skip`, `-run`,
`continue-on-error`, `|| true`); each ported fixture names the product behavior it still pins.

### E. No-suppression census contract (CI observability)

Re-executed at `35b57c3` by the drafter (commands in the attestation section), identical to all
three round-05 and all three round-06 files: **104** `t\.Skip` pattern lines under `internal`;
**103 call sites** = 68 `t.Skip(` + 35 `t.Skipf(` + 0 `t.SkipNow(`; **56 files** matching /
**55 with call sites** (sole non-call: `internal/runner/launch_precheck_skip_test.go:58`).
Policy: emit `-json` on **all three legs** (the workflow change is its own reviewed change;
`-json` selects nothing; the unfiltered matrix is untouched). **Windows leg:** every firing SKIP
event matched row-by-row to an individually reviewed exclusion-table row (the owner's "truly
inapplicable and reviewed") or the leg fails; AF_UNIX fallback skips convert to hard failures
(§D.7). **macOS/Ubuntu legs:** the currently-firing skips are the enumerated reviewed baseline,
and any **newly firing** skip fails the leg — the only shape that makes the owner's "absence of
new suppression" reviewable on the still-green legs. A firing skip is non-execution evidence,
never behaviour-exercised evidence. CI runs `go test ./... -count=1 -timeout 45m` with no `-v`
today (`tests.yml`), so `--- SKIP` is invisible without `-json` — the census is required, not
optional.

### F. Coverage honesty envelope

Hosted windows-latest x64 only (Windows Server image, `runneradmin` context, Go 1.26.8,
`core.longpaths true` **set by CI at `tests.yml:38`, not an OS default and not a product
behavior**); **ARM64 assets are build-only with zero executed evidence and no ARM64 gate is
added**; no native-Windows claim by anyone on this run; WSL/Wine are not accepted as Windows
evidence; per-run claims kept separate (historical runs 35987916696 — which also failed Ubuntu —
and 36009912946); runner image recorded with every evidence entry, re-baselined if it drifts;
and one documented sentence that this workspace itself sits on a cross-machine shared volume
whose Windows-side semantics are unvalidated. **No hosted probe upgrades an inference**: every
probe in §H is mechanics/diagnostic only — never durability evidence, never acceptance evidence,
and no false durability-probe claim is made anywhere.

### G. Hosted facts register — H1–H7 outcome→behavior mappings (fully expanded)

Provenance note: rounds 3 produced two accepted H-numbering schemes (kimi-1's and zcode-1's;
claude-1 round-04 accepted both, "H1–H7 as written, with the amendments below folded in", and
accepted zcode-1's list likewise). This register is the merged canonical set: every hosted fact
from both schemes and from the consensus probe bundle appears exactly once, with its
pre-agreed outcome→behavior mapping stated for **every** branch — each fact's outcome is mapped
before it is observed; **no result silently changes policy; no probe upgrades an inference.**
(Where a fact's letter was superseded in round 4 — the FlushFileBuffers residual — the
supersession is recorded in the row.)

- **H1 — ACL round-trip and inherited-ACL dumps** (parent directory of the store, and `%TEMP%`;
  merges zcode-1 r3 H1 and kimi-1 r3 H6; the consensus probe bundle's "ACL parent dump").
  Probe: create store under the §D.1 DACL, verify via `GetNamedSecurityInfo`; dump the parent's
  and `%TEMP%`'s effective ACLs. Mappings: round-trip failure ⇒ **diagnosis of the x/sys chain;
  the owner-only policy is unchanged**; permissive parents/`%TEMP%` ⇒ affects **fixture
  placement only** (and pre-existing-store handling stays refuse-and-instruct), **never the
  owner-only policy**.
- **H2 — `FlushFileBuffers` on a read-only directory handle** (kimi-1 r3 H1, **narrowed by
  kimi-1's own round-4 fetch**, folded in by claude-1 r4; half of zcode-1 r3 H2; the probe
  bundle's "`SyncFile` on a read-only directory handle with raw error"). The spike is
  **foreclosed twice over as acceptance evidence** (no documented directory subject — the
  documented subjects are a file, a communications device, a named-pipe server end, and a volume
  handle; the documented `GENERIC_WRITE` precondition is unmet on our read-only handles) and is
  **diagnostic-only forever**: the probe records the raw error; **no outcome can enable a
  mechanism**; the refusal-branch table stands regardless.
- **H3 — Rename mechanics probes and NTFS volume assertion** (the write-through half of
  zcode-1 r3 H2, made concrete by the consensus probe bundle). Probe: directory move accepted;
  no-REPLACE fails on an existing destination; `MOVEFILE_WRITE_THROUGH` accepted; runner volume
  is NTFS. Mappings: all as expected ⇒ the §A disclosure tiers are pinned as written
  (mechanics only — the write-through *guarantee* basis remains the documented sentence, and
  the no-REPLACE basis remains contrapositive-plus-probe, never "documented"); any unexpected
  outcome ⇒ the affected disclosure is re-labelled and escalated to review — **no probe result
  upgrades an inference or rewrites a disposition silently**; a non-NTFS volume voids the
  coverage-envelope claim and stops the leg for re-baselining.
- **H4 — `ReadDir`-on-regular-file error class and toolchain pin** (kimi-1 r3 H2 + zcode-1 r3 H5;
  the probe bundle's "`ReadDir`-on-regular-file error class" and "`runtime.Version()` + toolchain
  pin"). Mappings: the §D.5 fix **ships under every outcome**; the probe (empty-nil vs classified
  error, at the CI-resolved toolchain, version read from the existing setup-go log) only
  **records the mechanism** in `IMPLEMENTATION.md` and selects the sweep's scope; if
  inconclusive, the cause stays recorded **undiagnosed** — never reclassified as fixture-only
  without the probe result.
- **H5 — AF_UNIX hosted creation** (kimi-1 r3 H3 + zcode-1 r3 H3). Mappings: creation works ⇒
  the test **runs on Windows** via the existing `makeUnsupportedSocket` fallback and keeps its
  refusal assertion (`TreeDigest` refuses non-regular/non-symlink entries), and the
  `t.Skipf` fallbacks at `tree_report_test.go:116/:119/:126/:250` become **hard failures**;
  creation fails hosted ⇒ **one individually reviewed exclusion row opens**, with the reason
  recorded and attackable.
- **H6 — CRLF hosted mechanism** (kimi-1 r3 H5 + zcode-1 r3 H4; the probe bundle's
  `git config --show-origin` dump, plus a hexdump of one mismatch file). Mappings: the dump
  **confirms or adjusts the per-invocation pinning layer** (§D.7) before any parser-side
  normalization is finalized; **the split rule is non-negotiable under every outcome** —
  normalize admissible for document-semantic comparisons, forbidden for evidence-integrity
  comparisons; adaptation stays within the two-layer design (per-invocation `-c` pinning
  load-bearing, repo `.gitattributes` second).
- **H7 — First hosted executions, sharing violations, read-only-attribute rename-over, and
  durable-kill under truthful `Alive`** (kimi-1 r3 H4 + H7 + zcode-1 r3 H6 + H7; the probe
  bundle's "read-only-attribute rename-over"). Mappings: `wait`/`usage` first hosted execution
  after the comma-ok fix ⇒ **one diagnostic cycle; green is not assumed; the reconciliation
  ledger is refreshed from that run**. Open-site sharing violations (`.lock` opens) ⇒
  **undocumented cause; structural handle discipline and `os.Root` routing first**; a bounded
  retry ships **only after self-held vs foreign-held classification**, loud, with full
  identity/exclusion verification re-run after any successful retry; the diagnostic cycle logs
  error class + retry outcome. Read-only-attribute rename-over ⇒ the clear-before-replace fix
  (§D.6) is pinned by its probe, mechanics only. Durable-kill failures persisting under P-A ⇒
  **evidence for the Job-Objects follow-up idea's urgency, never an in-flight widening** of this
  idea's scope.

### H. Implementation stages (bounded)

zcode-1's Stage 0–7 frame, unchanged; one implementation owner (zcode-1), claim + worktree
mapping recorded in `IMPLEMENTATION.md` before any code edit:

- **Stage 0** — contract/smoke/ledger: Mkfifo build tag + `app_test.go` comma-ok fixes; Ubuntu
  `GOOS=windows` cross-compile guard; reconciliation ledger emitted; one diagnostic cycle
  budgeted, green not assumed.
- **Stage 1** — snapshot privacy/ACL (§D.1), including the `denyRead` helper shared with the
  sweep.
- **Stage 2** — gate-name encoding + raw-ID allowlist (§D.4), legacy read-fallback + shadow
  retirement.
- **Stage 3** — directory durability: **opens with the §B table and its per-row disposition
  before any code is written**; the `fsutil.SyncDir` named-type contract is the audit of record
  and, under the refusal branch, the rooted sites' refusal emitter.
- **Stage 4** — file sharing (§D.6) + the open-site diagnostic cycle (H7).
- **Stage 5** — P-A liveness (§D.2) + W-SHELL blocking missing-`sh` refusal (§D.3).
- **Stage 6** — fixture-portability sweep (§D.9) + ACP/AF_UNIX/CRLF (§D.7).
- **Stage 7** — validation/release: unfiltered three-leg matrix green, census enforcement (§E),
  coverage envelope statements (§F), release sequence (§O).

**One early hosted probe bundle** (build-tagged, test-only, ~one 20-minute cycle, executed
before/alongside Stage 0–1): rename mechanics (H3), NTFS volume assertion (H3), `SyncFile` on a
read-only directory handle with raw error (H2), `runtime.Version()` + toolchain pin (H4), ACL
parent dump (H1), `git config --show-origin` (H6), `ReadDir`-on-regular-file error class (H4),
read-only-attribute rename-over (H7). **Mechanics only — never durability, never acceptance
evidence; no hosted probe upgrades an inference.**

### I. Observable acceptance criteria (stable IDs; reviewer-checkable)

Each criterion states an observable, reviewer-checkable behavior. A criterion is met only on the
hosted evidence envelope of §F (Windows = hosted windows-latest x64).

**Build**

- **AC-BLD-1** — Step-1 fixes land first: `syscall.Mkfifo` is build-tagged so Windows falls into
  `internal/evidence`'s existing `makeUnsupportedSocket` fallback; `app_test.go:155/:157` use
  comma-ok assertions. Reviewer: hosted Windows log shows `wait` and `usage` packages producing
  results (no whole-binary abort).
- **AC-BLD-2** — Ubuntu leg runs and passes `GOOS=windows go build ./... && GOOS=windows go vet
  ./...` (windows/amd64 and windows/arm64 both build). Reviewer: workflow log; **this is
  compile-only evidence — no ARM64 execution is claimed anywhere.**

**Privacy**

- **AC-PRIV-1** — On Windows, store creation sets the owner-only protected DACL exactly as §D.1.
  Reviewer: create-then-requery readback unit test passes (exactly one access-allowed ACE, the
  token user; protected, non-inherited).
- **AC-PRIV-2** — Verification refuses any non-owner grant, inherited or explicit, naming the
  trustee; refusal string keeps the sentence-stable prefix. Reviewer: refusal-string test.
- **AC-PRIV-3** — Adversarial hosted tests: `BUILTIN\Users` → refuse; `Administrators` →
  refuse; inherited-only → refuse; own creation → pass. Reviewer: hosted log, all four named.
- **AC-PRIV-4** — FAT/exFAT/no-ACL volumes refuse (never weaken). Reviewer: refusal-path test
  where exercisable; design-level refusal otherwise (§N blind spot: not exercisable on the
  runner).
- **AC-PRIV-5** — Pre-existing stores get refuse-and-instruct with the documented user-invoked
  repair named; **no in-place DACL rewrite of user data**. Reviewer: refusal text test; code
  review confirms no rewrite path.
- **AC-PRIV-6** — `!IsDir()` and reparse/symlink rejection kept on both platforms; Unix behavior
  byte-identical. Reviewer: diff shows Unix-path neutrality.

**Gate names**

- **AC-NAME-1** — Universal encoding inside `GatePath` on every OS (PathEscape base + `:`
  escape + device/dot/space guards); `.tmp` inherits by construction. Reviewer:
  `block-a->block-b` → `block-a-%3Eblock-b.gate.json` unit test on all three OSes; hosted
  `ERROR_INVALID_NAME` failures gone.
- **AC-NAME-2** — Legacy `a->b.gate.json` read-fallback works; shadow retirement on first
  rewrite; mixed-name deck test exists. Reviewer: unit tests.
- **AC-NAME-3** — Raw-ID charset allowlist at new pipeline creation; cosmetic-only grandfathering
  on load; unsafe IDs refused at load and at both interpolation backstops on every OS, no legacy
  exemption. Reviewer: refusal tests per class.
- **AC-NAME-4** — Traversal/separators/absolute/drive/ADS-colon/reserved-name/trailing-dot-space
  negative tests run on every OS and validate **on the raw ID** (the executed `IsLocal`
  unsoundness table is the recorded reason). Reviewer: the four previously wrongly-accepted
  traversal IDs now refused, including `../../other-idea`.

**Durability (refusal branch)**

- **AC-DUR-1** — Audit of record: the `fsutil.SyncDir` named-type contract is compiler-enforced
  fail-closed on Windows and is the rooted sites' refusal emitter. Reviewer: mechanical
  type/grep check at the reviewed commit — no directory handle reaches `SyncFile` on Windows;
  every Unix-only directory sync stays behind `//go:build !windows`.
- **AC-DUR-2** — Every rooted row (A2, A3, B3, B4, B5, C1, C2) refuses on Windows pre-mutation
  with a named, user-visible, blocking refusal; C1/C2 refusals ordered before the rename
  (`:333`/`:570`), stage cleared by the existing `defer dir.Remove(stage)` (`:328`/`:565`).
  Reviewer: one hosted test per refusal path (7 rows), each asserting message + non-zero +
  nothing-published state.
- **AC-DUR-3** — Plain-row conversions (A1 det-stage; B1/B2 stage-`Mkdir` + WT-move +
  EEXIST-tolerant recheck) are Windows-confined and dormant behind the POSIX gate; Unix path
  byte-identical. Reviewer: unit tests pin requirement (a) (stage derivation non-colliding with
  every valid final name) and — on the deviation branch only — requirement (b) (publish-time
  `validHash` gate at A3).
- **AC-DUR-4** — **No rooted site adopts a path-based `MoveFileEx`** without an explicit owner
  deviation **and** participant signoff on the converted contract. Reviewer: enumerate every
  `MoveFileEx` call site; each is either a plain-site conversion or deviation-authorized.
  Violation is review-blocking. **Baseline implementation neither waits for nor presumes the
  deviation answer; silence is not approval.**
- **AC-DUR-5** — Unix behavior at every row byte-identical to today. Reviewer: full diff review
  for GOOS confinement; Unix legs green.

**Locks/sharing**

- **AC-LOCK-1** — Share-flag-correct opens: route through `os.Root` (FILE_SHARE_DELETE) where
  possible; raw `windows.CreateFile` helper only where rooting cannot express the site;
  close-before-rename discipline. Reviewer: code review against the enumerated sites.
- **AC-LOCK-2** — Read-only-attribute trap handled (clear `FILE_ATTRIBUTE_READONLY` before
  replace). Reviewer: probe-pinned mechanics test (H7).
- **AC-LOCK-3** — Bounded retry (~250 ms) only on `ERROR_SHARING_VIOLATION`, only for
  third-party holders after self-vs-foreign classification, loud, verification re-run after.
  Reviewer: retry-path tests + log assertions.
- **AC-LOCK-4** — Open-site sharing failures get exactly one hosted diagnostic cycle (error
  class + retry outcome) recorded in `IMPLEMENTATION.md`; cause stays undiagnosed if
  unresolved. Reviewer: diagnostic artifacts present.

**Process liveness**

- **AC-PROC-1** — `alive` truthful via `OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` +
  `GetExitCodeProcess`; STILL_ACTIVE caveat documented; **fail closed on `ERROR_ACCESS_DENIED`**.
  Reviewer: unit tests for alive/dead/unverifiable.
- **AC-PROC-2** — `bootID` stays `""` (documented scoped null); `supportsDurableKill()` stays
  false with honest refusal strings; `proclive_windows.go` comment corrected. Reviewer: tests +
  diff.
- **AC-PROC-3** — The four broken callers fixed (`durablekill.go:29`/`:63`,
  `verification.go:142`, `reviewsnapshot.go:244`). Reviewer: diff touches exactly those call
  shapes.
- **AC-PROC-4** — **P-B is not implemented in this idea** (no Job Objects, no `GetProcessTimes`
  attribution, no `QueryFullProcessImageName` in kill-safety gating); the follow-up idea is named
  in `IMPLEMENTATION.md`. Reviewer: absence check.

**Shell gate**

- **AC-SHELL-1** — `checks:` and captured verification call `exec.LookPath("sh")` once, before
  any work; on failure a named, actionable, blocking refusal names Git for Windows and states
  "refusing rather than executing unverified"; exits non-zero; recorded in the run. Reviewer:
  refusal-path test.
- **AC-SHELL-2** — Hosted tests pin both paths: `sh` forced off `PATH` → message + non-zero;
  `sh` present → normal path. Reviewer: hosted log contains both.

**ACP / AF_UNIX / CRLF**

- **AC-ACP-1** — Two drain tests split per §D.7 ((a) gated ≤2 KiB capacity-independent; (b)
  ungated 16384-write asserting the 8192 ring cap); both survive everywhere; U1 product drain
  ordering untouched. Reviewer: tests present; `spawn.go` Stop/Wait diff-neutral.
- **AC-AFUN-1** — AF_UNIX test runs on Windows via `makeUnsupportedSocket` fallback keeping its
  `TreeDigest` refusal assertion; the four `t.Skipf` fallbacks become hard failures. Reviewer:
  hosted log; no skip at those lines.
- **AC-AFUN-2** — No exclusion row exists for AF_UNIX unless hosted creation itself fails (then:
  one individually reviewed row). Reviewer: ledger.
- **AC-CRLF-1** — Split rule enforced: normalization only for document-semantic comparisons;
  **forbidden** for evidence-integrity comparisons (snapshot content, `TreeDigest`). Reviewer:
  tests on both sides of the rule.
- **AC-CRLF-2** — Per-invocation `-c core.autocrlf=false -c core.eol=lf` pinned on parley's own
  git invocations; repo `.gitattributes` lands as its own reviewed commit; product parsers
  CRLF-tolerant. Reviewer: diff + tests.
- **AC-CRLF-3** — `git config --show-origin` dump present in the probe bundle. Reviewer: bundle
  contents.

**strict_gate**

- **AC-STRICT-1** — Branch-independent fix ships: stat the round path, veto when it exists and
  is not a directory; cross-platform pinning test. Reviewer: test constructs
  file-where-directory-belongs and asserts the veto.
- **AC-STRICT-2** — Probe prints `runtime.Version()`, the raw `os.ReadDir(regular file)` error
  and its `errors.Is(err, fs.ErrNotExist)` classification at the CI toolchain; every
  `ErrNotExist`-after-directory-read decision swept. Cause recorded **undiagnosed**; no branch
  called confirmed without the probe. Reviewer: probe artifacts + `IMPLEMENTATION.md` entry.

**Fixtures & census**

- **AC-FIX-1** — Sweep classes fixed per §D.9 (test-binary re-exec / `cmd.exe /c`; PATHEXT only
  where discovery is the subject; `denyRead` helper for Chmod-0; `USERPROFILE` set;
  `t.TempDir()` + volume-aware joins; TUI temp-perms via ACL helper). Reviewer: per-fixture diff
  names the product behavior each still pins.
- **AC-FIX-2** — Reconciliation ledger emitted before the sweep and refreshed after step-1
  changes; one reviewed row per exclusion; diff-level check no `t.Skip` added without a row;
  close-time suppression grep (`Skip`, `-run`, `continue-on-error`, `|| true`) clean. Reviewer:
  ledger + grep at close.
- **AC-CENSUS-1** — `-json` emitted on **all three legs**; the workflow change is its own
  reviewed change; the unfiltered matrix is untouched. Reviewer: workflow diff.
- **AC-CENSUS-2** — Windows leg: every firing SKIP matched row-by-row to an individually
  reviewed exclusion row or the leg **fails**; AF_UNIX fallback skips are hard failures.
  Reviewer: census-vs-ledger diff on the hosted log.
- **AC-CENSUS-3** — macOS/Ubuntu legs: the enumerated baseline (103 call sites) is pinned; any
  newly firing skip **fails the leg**. Reviewer: baseline artifact + leg outcome.
- **AC-CENSUS-4** — Census numbers pinned at the reviewed commit with their commands
  (104/103 = 68+35+0 / 56/55). Reviewer: re-run commands; numbers match.

**Coverage & matrix**

- **AC-COV-1** — Native-Windows execution evidence is hosted windows-latest **x64 only**; runner
  image recorded with every evidence entry and re-baselined if it drifts; `core.longpaths true`
  documented as CI-set (`tests.yml:38`). Reviewer: evidence entries.
- **AC-COV-2** — ARM64 is build-only (zero executed evidence, no ARM64 gate); WSL/Wine never
  accepted; per-run claims separate; the workspace-shared-volume sentence present; **no false
  durability-probe claim anywhere** (all probes labelled mechanics/diagnostic). Reviewer:
  documentation + absence of contrary claims.
- **AC-MATRIX-1** — On the final released commit: unfiltered build/test **green on Windows,
  macOS and Ubuntu** hosted legs, with participant review of absence of new suppression.
  Reviewer: hosted run on the release commit + review signoff.

**Release**

- **AC-REL-1** — Both predecessor gates exist (drafter-verified this session, read-only:
  `worktrees/lean-organizer/parley-deck/inbox/codex-1-to-user_release-1.49.1_done.md`, 6412
  bytes, reports CLI 1.49.1 at `54e0798`, Windows experimental, winget held; and
  `worktrees/designated-implementer/parley-deck/inbox/codex-1-to-user_meta-protocol-change-designated-implementer_done.md`,
  12075 bytes, mtime 2026-09-25 11:13, `status: agent-controlled-delivery-complete`, releases
  the sequencing hold, reports CLI 1.50.0 and skill 2.14.0 on GitHub and Homebrew). Both checks
  are **existence plus self-testimony, not channel audits**.
- **AC-REL-2** — After reviewed closure, integrate onto **latest origin/main** and select the
  next version **strictly above the latest actually-released version at that time** — at minimum
  above the reported **1.50.0**; verify actual releases/tags at release time and never select
  below an actual release. Reviewer: version arithmetic against the live release list at release
  time.
- **AC-REL-3** — Release includes GitHub Windows assets, Homebrew `parley-deck-cli.rb`, and the
  held CLI winget PR (one application per PR); skill channels only if the skill changed.
- **AC-REL-4** — Windows experimental label removed **iff** the released commit's hosted Windows
  leg is green; otherwise label retained, winget held, outcome reported — **never an incomplete
  release represented as shipped; released tags never move.**
- **AC-REL-5** — A participant (never codex-1) independently verifies **every** released channel
  before completion is reported.
- **AC-REL-6** — Completion handoff `parley-deck/inbox/codex-1-to-user_windows-portability_done.md`
  written with what shipped, hosted evidence, deferrals, `organizer-usage.md`, and owner actions
  left.

**Process**

- **AC-IMPL-1** — `IMPLEMENTATION.md` records the zcode-1 claim + worktree mapping **before any
  code edit**; one implementation owner at a time; claude-1 and kimi-1 review each stage.
- **AC-DEV-1** — The containment-deviation request remains **UNAPPROVED and nonblocking**
  throughout baseline implementation; silence is not approval; if the owner grants it, the
  converted contract (per §P.1) requires its own participant signoff before any conversion code.

## Idempotence & recovery

- **Class A (det-stage, where implemented):** an interrupted publication leaves the deterministic
  stage as a visible blocker whose **existence prevents replay** (`verification.go:192-193`,
  `reservation_recovery.go:79` invariants reproduced one name removed); the final name appears
  only via the atomic move, so readers never see torn bytes; no-REPLACE prevents clobber of a
  completed publication; concurrent publishers race on the stage `O_EXCL` exactly as today they
  race on the final-name `O_EXCL`. Stage litter is a visible, actionable blocker **by design**.
- **Class B (B1/B2 where implemented):** idempotent ensure-exists semantics preserved with
  EEXIST-tolerant recheck after the WT-move; nothing published on failure.
- **Class C under the refusal branch:** refusals fire **before** the rename; the existing
  `defer dir.Remove(stage)` (`:328`/`:565`) clears the stage with nothing published.
- **Class D:** re-sync only, never rewrites bytes (D4a `:426` never rewrites; D1's original
  already atomic via `os.CreateTemp` + `defer os.Remove` at `state.go:205`/`:209`).
- **Refused rows (rooted):** pre-mutation ordering means **nothing is published** when a refusal
  fires; the existing anti-replay tokens (`Mkdir(charge)` "already reserved" at `:284`) keep
  their today semantics.
- **Retry (W-LOCK):** bounded, classified (foreign-transient only), loud, with verification
  re-run — never a silent retry loop.

## Known risks / de-risking

1. **Two reachable Windows features refuse under the selected branch** (reservation-intent
   publication; parent-recovery/recovered-parent apply) in exchange for preserving `os.Root`'s
   kernel-enforced per-open reparse resistance and handle-relative rename at the publication
   operation. The owner reserved exactly this trade; §P.1 puts one costed question to the owner.
2. **Sequential round 6 bought convergence at the cost of independence** — see §N; treat the
   unanimity as a shared prior, not independent confirmation.
3. **What would make the agreed position wrong:** a non-NTFS runner volume; a directory move
   rejected on the runner; no-replace clobbering an existing destination on the runner; the
   deterministic-stage suffix colliding with a valid final name (prevented by requirement (a),
   pinned by unit test); a hosted leg staying red on any of the 14 historically failing packages;
   or an ACL invariant the `runneradmin` context cannot exercise. The owner declining the
   deviation does **not** make the position wrong — the refusal branch stands as a reviewed
   outcome.
4. **Blind spots (no participant addressed; disclosed, not solved):** no power-cycle/real-durability
   evidence is obtainable on hosted runners; ARM64 zero executed evidence (by design); FAT/exFAT
   refusal designed but not exercisable on the runner; non-`runneradmin` contexts and native
   desktops unvalidated; P-B attribution and `bootID` have no Windows analogue in scope (scoped
   null / deferred); refusal strings have no real-user usability testing; five census-class
   errors across rounds argue for mechanical generation wherever a number is load-bearing
   (adopted: §attestation, §B, AC-CENSUS-4).
5. **Census/diagnosis risks:** the `strict_gate` cause and the open-site sharing cause are both
   recorded undiagnosed pending their probes (H4/H7); no fix is conditioned on a diagnosis.

## Provenance & dependency check (§15.3, repeated as required)

FINAL cites the consensus's closed verdict conflicts; per the consensus instruction the
dependency check is repeated: **no acceptance criterion in this document rests on any formerly
disputed claim's contested reading.** Every conflict that arose in rounds 1–6 is closed, each by
owner weakening or by non-owner verdicts, never by counting:

1. **`o_DIRECTORY`/ReadDir chain (strict_gate mechanism)** — both source-level mechanisms
   excluded (go1.27.1 and the CI-fetched go1.26.8, PRIMARY both); hosted phenomenon undiagnosed;
   branch-independent fix ships regardless (AC-STRICT-1/2).
2. **Skip census 107 vs 103** — closed by claude-1's self-correction to 103 call sites
   (68+35), confirmed non-owner by both peers and re-executed by all three participants and by
   this drafter at `35b57c3`.
3. **Create-entry coverage of `O_SYNC`/`FILE_FLAG_WRITE_THROUGH`** — closed by claude-1's round-5
   withdrawal of N5 after his own PRIMARY re-fetch; coverage "neither documented nor
   documented-against"; the staged write-through rename then closed the design (round 6).
4. **"Replay blocked by destination-exists failure" (zcode-1 r5)** — withdrawn by its owner
   after two concordant non-owner WRONG verdicts; the deterministic-stage condition repairs the
   hole.
5. **"Shape-only" change to partial-publication (kimi-1 r5)** — the no-clobber half stands
   confirmed; the second half withdrawn by its owner after claude-1's WRONG verdict.
6. **Containment equivalence (zcode-1 r4)** — verdicted WRONG by claude-1 r5 (`filepath.Clean`
   is lexical-only; `os.Root` enforces reparse-resistant handle-relative operations), unrebutted;
   narrowed then superseded by the owner's authority reservation.
7. **N6 (`MoveFileEx` moves directories; write-through sentence)** — owner-asserted by claude-1
   r4, UNVERIFIED through round 4; CONFIRMED non-owner by kimi-1 and zcode-1 in round 5 from
   their own fetches; basis layers (including the flag-table contrapositive limitation) recorded
   in §A.
8. **N7–N10 (round 6)** — owner-asserted by claude-1 r6. **N7 carries exactly one round-6
   non-owner CONFIRMED — kimi-1's (layered tags: `PRIMARY` code side, `PRIMARY(r4)` for the
   FlushFileBuffers subject list, `SECONDARY` on zcode-1's round-05 verdict for the
   `GENERIC_WRITE` precondition); zcode-1 issued no round-6 verdict on N7 (correction 1).**
   N8, N9 and N10 carry both peers' round-06 non-owner CONFIRMEDs; N9 re-verified read-only by
   the drafter this session (both gate files exist; sizes and frontmatter as stated in
   AC-REL-1).

## Alternatives disposition (§15.6(a)/§15.6(c), ALT-1..ALT-16)

- **ALT-1 — ACL implementation: ADOPT** the direct `golang.org/x/sys v0.36.0` chain (zero new
  deps; full symbol map verified). **REJECT** `hectane/go-acl` (re-encodes the 0700 model being
  abandoned; new dep) and `icacls` (localizable text, fragile, extra process, TOCTOU).
- **ALT-2 — Locking: ADOPT** the in-tree `lock_windows.go` LockFileEx implementation plus
  share-mode/handle discipline. **REJECT** `gofrs/flock` and named mutexes (no lock-acquisition
  library fixes share-mode conflicts; the in-tree design is deliberate and working).
- **ALT-3 — Process scope: ADOPT P-A** (truthful liveness, ~40 lines, fail-closed). **REJECT
  P-B in-idea** (Job Objects/attribution would silently loosen the full-argv kill-safety
  calibration); named follow-up.
- **ALT-4 — ACP drain test: ADOPT the two-test split.** **REJECT** in-`Wait` `cmd.Stderr=io.Writer`
  draining for ACP (the observer must stream non-blockingly during the live session) and
  **REJECT** retuning/chunked-writes (only moves the capacity coupling).
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
- **ALT-11 — Handle-based rename (`SetFileInformationByHandle`/`FileRenameInformation[Ex]):
  REJECT** — carries no write-through flag; adopting it would repeat the inference error N5
  closed. This closure is why containment had no free technical exit.
- **ALT-12 — Documented-nil `SyncDir` / documented-weakening floor: REJECT** (withdrawn r3/r4;
  no silent weakening, no fourth branch).
- **ALT-13 — Restructure everywhere including rooted sites (path-based `MoveFileEx` + lexical
  sibling check): NOT ADOPTED.** A real weakening of `os.Root` containment (TOCTOU-exposed,
  reparse-resolving) that the owner reserved to an explicit deviation; carried fully specified
  and **unapproved** (§P.1). It is not an implementable automatic fallback, and baseline
  implementation is not conditioned on it.
- **ALT-14 — Randomized stage names at class A: REJECT** (loses the replay blocker on the
  crash-before-move interval; A3 is Windows-reachable and budget-charging). Deterministic stages
  adopted. Randomized stages remain correct at C, where the in-tree pattern already uses them on
  every OS and refusal precedes the rename under the selected branch.
- **ALT-15 — REPLACE at C1/C2 inherited from the `ReplaceSyncedFile` contract: DEFERRED to the
  deviation branch** (under the selected branch C1/C2 refuse). If that branch ever activates,
  the flag is a deliberate, disclosed choice — drafter preference no-REPLACE — requiring
  participant signoff, not a pattern-match inheritance.
- **ALT-16 — Windows-leg-only `-json` census: REJECT/SUPERSEDED** by the all-legs variant with
  per-leg baseline enforcement (§E) — the only shape that reviews "no new suppression" on the
  still-green legs.

## Correlated agreement (§15.6(b))

Unanimity among these three related models is a **shared prior, not independent evidence**, and
the round-6 convergence was *sequentially produced by the owner's own design*: one participant at
a time, each reading the prior positions — claude-1 joined the class-A restructure after reading
both round-05 files; kimi-1 joined the plain/rooted split after reading claude-1's round-06;
zcode-1 joined both after reading both round-06 files. The deck holds no dissenting position on
any of the three conflicts, and this FINAL says so rather than read the convergence as
independent confirmation. See §Known risks item 3 for what would make the agreed position wrong.

## O. Roles and release sequence

- **Roles:** zcode-1 — drafter (this FINAL) and **single implementation owner** (claim + worktree
  mapping in `IMPLEMENTATION.md` before any code edit); claude-1 and kimi-1 — independent
  reviewers (Phase-6, refutation-default); codex-1 — organizer only (never implements, verifies
  code, authors participant verdicts, or signs off; organizer commits start
  `[codex-1] windows-portability:`). Channel verification after release is by a participant,
  never codex-1.
- **Release sequence (owner-standing authorization, gated):** (1) reviewed closure of this FINAL
  and its implementation/review phases; (2) both predecessor handoffs present (AC-REL-1 —
  verified); (3) integrate onto **latest origin/main**, keep main linear, no development PRs;
  (4) select the next version **strictly above the latest actually-released version** at that
  time — at minimum above the reported **1.50.0** (AC-REL-2); (5) release GitHub including
  Windows assets, Homebrew `parley-deck-cli.rb`, and the held CLI winget PR (one application per
  PR); skill channels only if the skill changed; (6) remove the Windows experimental label **iff**
  the released commit's hosted Windows leg is green (AC-REL-4); (7) participant channel
  verification (AC-REL-5); (8) completion handoff (AC-REL-6). Released tags never move.

## P. Open items deferred to implementation

1. **The containment-deviation question (owner decision; UNAPPROVED and nonblocking).** The
   fully specified, **unapproved** alternative — recorded by the organizer in
   `inbox/codex-1-to-user_windows-portability_containment-deviation.md` — would convert the seven
   rooted rows to path-based `MoveFileEx(MOVEFILE_WRITE_THROUGH)` with: kimi-1's guard set
   (1)–(4) (root-resolved paths via `Root.Name()`; retained sibling check; every rename name a
   product constant or `validHash`-gated, including the publish-time gate at A3; deterministic
   stages at A / randomized at C; explicit weakened-resistance disclosure), claude-1's items (the
   precise statement of what is given up: kernel-enforced per-open reparse resistance and
   handle-relative rename, replaced by a lexical check plus a path-based call with a TOCTOU
   window; the named but unselected handle-derived-path alternative `GetFinalPathNameByHandle`,
   which narrows but does not close the window; both branches' costs side by side, including the
   §Correction-3 benefit statement), zcode-1's hosted rooted-escape test against the conversion
   helper (mechanics only, never durability evidence), and the C1/C2 flag decision (drafter
   preference, stated as input only: no-REPLACE, matching the publish-once constants — an
   existing final artifact should be a loud replay signal, not an overwrite). If granted, every
   refusal row converts per the already-designed restructure and the D rows derive — **no row is
   added or dropped**; the exact contract still requires participant signoff. If declined, the
   refusal-branch table is final as written. **Implementation of the baseline does not wait for,
   or assume, this answer.**
2. **Hosted probe bundle execution** and the H1–H7 register outcomes (§G) — mappings
   pre-agreed; no result silently changes policy.
3. **`strict_gate` hosted diagnosis** (fix ships regardless; the probe selects the recorded
   mechanism and the sweep's scope).
4. **Open-site sharing-violation diagnosis cycle** (Stage 4).
5. **Reconciliation-ledger refresh** after step-1 fixes change the failing denominator.
6. **Refusal-string wording and `IMPLEMENTATION.md` claim/worktree mapping** — implementation
   owned by zcode-1, reviewed by claude-1 and kimi-1.
7. **Version selection above 1.50.0 at release time** (AC-REL-2), after both gates (present) and
   reviewed closure, per the owner's release sequence.

## References

- Consensus (signed, never rewritten by this FINAL): ./consensus.md — includes the three
  append-only signoffs carrying the three binding FINAL corrections applied above.
- Rounds: ./round-01/ … ./round-06/ (round-06 is the current position of record for all three
  participants); briefs: ./round-02-brief.md, ./round-03-brief.md, ./round-04-brief.md,
  ./round-05-brief.md, ./round-06-brief.md; durability brief: ./durability-followup-brief.md;
  organizer notes: ./organizer-notes.md; usage ledger: ./usage-ledger.jsonl.
- Binding scope/authorization: ./00-prompt.md (status at drafting of this section header set to
  final by this FINAL).
- Evidence shared at launch: ./source-context/README.md with release-ci-revalidation-claude-1.md
  and release-repair-plan-zcode-1.md (participant proposal, not a decision). Hosted run IDs:
  35987916696 (also failed Ubuntu), 36009912946.
- Release-gate files (read-only verified): the two inbox paths in AC-REL-1.
- Deviation request (unapproved): inbox/codex-1-to-user_windows-portability_containment-deviation.md.
