---
agent: zcode-1
idea: windows-portability
round: 6
date: 2026-09-25
responding-to: [claude-1/round-06, kimi-1/round-06, claude-1/round-05, kimi-1/round-05, round-06-brief, user round5-unresolved authorization]
---

# Round 6 — zcode-1: joins the converged design in full — deterministic stages, plain/rooted split, 15 IDs / 16 operations

## Protocol attestation

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7"}
```

Packet sha256 re-verified this session by `shasum -a 256`; matches. **Sequence position:** I am
launched third and last; both peer round-06 artifacts exist, and I answer them (both peers' latest
positions are round 6). Read this round, in full: `inbox/user-to-codex-1_windows-portability_round5-unresolved.md`
(verbatim), `inbox/codex-1-to-user_windows-portability_round5-unresolved.md`, `round-06-brief.md`,
`00-prompt.md` (all authorization sections), all three `round-05` artifacts, `round-06/claude-1.md`,
`round-06/kimi-1.md`, and §15 of the protocol packet.

HEAD is `9e55050`; `git diff --name-only 6b87cf8..9e55050 -- '*.go'` = **0 files** (`PRIMARY`,
re-executed), so every round-05 locator transfers. Every code claim below was re-verified by me at
`9e55050` (`PRIMARY`): the 15-line name-keyed call-site listing, the six helper `func` lines, the
plain/rooted classification at all fifteen sites, the replay-blocker comments, the C1/C2 ordering
lines, the GOOS-gate census, the skip census, `state.go`'s `writeState`, and the two handoff files.
Microsoft Learn pages were **not** re-fetched this round; documentation basis is tagged
`PRIMARY(r5)` (my own round-5 fetch) or `SECONDARY` (a named peer's round-5 non-owner verdict),
never presented as fresh.

**Ownership under §15.1.** I own my round-05 claims (the destination-exists replay claim, the
restructure-everywhere selection, the "one named ratification outstanding" statement); I issue no
verdict on them — only self-corrections, which take effect immediately when they weaken. Verdicts
below are on peers' claims.

## User direction

From `parley-deck/inbox/user-to-codex-1_windows-portability_round5-unresolved.md` (2026-09-25,
`status: authorized`), answering the organizer's escalation. The question put was *"Povoliť Windows
idei ešte jedno zmierovacie kolo o troch zostávajúcich rozporoch?"* — "Allow the Windows idea one
more reconciliation round on the three remaining conflicts?" The owner selected, in the original
Slovak: *"Áno, postupne" — "Jedno kolo, v ktorom účastníci odpovedajú jeden po druhom a každý vidí
aktuálne pozície predchádzajúcich. Trvá dlhšie, ale znova sa nepretnú. Ak sa nezhodnú, organizátor
sa zastaví."* The owner's exact English translation as supplied in that file:

> "Yes, sequentially — one round in which participants answer one after another, each seeing the
> current positions of those before it. It takes longer, but they will not cross again. If they do
> not agree, the organizer stops."

Three further binding sentences from the same answer, quoted because each controls part of this
file: *"No participant may treat another participant's superseded position as its current
acceptance."* / *"Any proposed weakening of os.Root containment is not accepted by this answer: it
must come back to the owner as an explicit deviation."* / *"If a material conflict remains after
round 6, stop and escalate again; this is not a standing waiver."* Scope as I read it: exactly ONE
sequential reconciliation round on the three named conflicts; no seventh round, no consensus, no
FINAL, no implementation, no hosted mutation, no signoff in this task.

## Position changes since round 5

1. **`SELF-CORRECTION` (weakening, immediate): my round-05 claim "replay is blocked by the move's
   destination-exists failure instead of by `O_EXCL`" is withdrawn — both peers' `WRONG` verdicts
   are correct, and I verified the counterexample myself this round.** The load-bearing invariant
   at the class-A sites is not no-clobber but that **an interrupted publication leaves a blocker
   whose existence prevents replay**: `verification.go:192-193` ("Exclusive writes leave any
   partial file in place. **Its existence prevents replay**; a torn/noncanonical artifact is an
   unresolved failure, never a retry") and `reservation_recovery.go:79` ("A partial publication
   stays visible and cannot be overwritten by retry"), both re-read at `9e55050` (`PRIMARY`). A
   destination-exists failure blocks replay only after a publication **completed**; after a
   crash-before-move with the randomized, swept stage I had selected (the `:319-328` in-tree
   pattern), the final name is free and a retry succeeds — replay, on the one class-A row
   (A3) that is Windows-reachable. My "strictly stronger for readers" framing was true for
   readers and never a statement about replay; I concede the scoped claim.
2. **`SELF-CORRECTION` (superseded by authority, not by merits): my round-05 selection
   "restructure everywhere, zero blocking refusals" is withdrawn as a deck-adoptable position.**
   My own containment sentence — "`os.Root` remains the confinement for every operation **except
   the final write-through move itself**" — conceded the exception, and the excepted operation is
   the publication. The owner's round-6 answer reserves exactly that exception to an explicit
   owner deviation. See Conflict 2.
3. **Adopted: claude-1's deterministic-stage condition for class A** (stage name derived from the
   final name, `O_EXCL` on the stage, not removed on the error path) — it repairs the hole my
   withdrawn claim missed, at the cost of one name attribute. See Conflict 1.
4. **Adopted: the plain/rooted split** (restructure at plain sites; rooted sites take the named
   blocking refusal pending an explicit owner deviation) and **the 15-ID / 16-runtime-operation
   table on kimi-1's row-ID scheme**, both verified independently this round. See Conflicts 2–3.
5. **Closed by kimi-1's round-05, which I could not see:** my round-05 "one named ratification
   outstanding" (kimi must ratify N5 and the restructure) was satisfied by kimi's round-05 file —
   it confirmed N5 non-owner and selected restructure — written nine minutes after mine. The
   crossing the owner observed is closed on this point.

Everything settled outside the three conflicts stands un reopened: owner-only protected DACL with
refuse-and-instruct; P-A with P-B deferred; blocking missing-`sh` refusal; universal gate-name
encoding plus raw-ID charset allowlist; the branch-independent `strict_gate` invariant fix with
cause recorded undiagnosed; the ACP split; the CRLF split; the label gate exactly as the owner set
it; all-legs `-json` with per-leg baseline enforcement; `zcode-1` as drafter and single
implementation owner with claude-1 and kimi-1 reviewing.

## Responses to others

### @claude-1 — round-06 (current position)

**Your class-A concession is accepted, and your deterministic-stage condition is adopted by me in
full.** You moved to the staged write-through rename on the merits — the write-through sentence is
scoped to the exact operation while my `O_EXCL`-entry coverage was an inference — and your
condition closes the one regression my selection carried. I verified the two invariant comments
(`verification.go:192-193`, `reservation_recovery.go:79`) and the `:319-328` randomized-swept-stage
pattern myself this round (`PRIMARY`): with a deterministic stage sibling, `O_EXCL` on the stage
reproduces exactly today's existence-prevents-replay invariant one name removed, the final name
still appears only via the atomic move (readers never see torn bytes), and no-REPLACE still
prevents clobber of a completed publication. Concurrency is unchanged: two concurrent publishers
of the same final name race on the same deterministic stage `O_EXCL`, as today they race on the
final-name `O_EXCL`.

**Your plain/rooted split is `CONFIRMED` by independent verification (`PRIMARY`, every site re-read
at `9e55050`).** Plain: A1 `os.OpenFile` (`trajectory_verify.go:146`), B1 `os.Mkdir` (`:118`),
B2 `os.Mkdir` (`:190`), D1 plain `os.Open` (`unchanged.go:267-279`) — and D1's original
publication `writeState` (`state.go:197-221`) is also plain-path, ending in
`fsutil.ReplaceSyncedFile` (`os.CreateTemp` stage, fsync, replace). Rooted: A2
`dir.OpenFile` on `*os.Root` (`verification.go:202`), A3 (`reservation_recovery.go:75`), B3
`parent.Mkdir` (`verification.go:264`), B4 `base.Mkdir` (`:283`), B5 `dir.Mkdir`
(`reservation_recovery.go:36`), C1 `dir.Rename` (`parent_recovery.go:333`), C2 (`:570`); the D
re-syncs D2/D3/D4a/D4b operate on rooted directory contexts. I also re-ran the GOOS census:
exactly eight `runtime.GOOS == "windows"` sites in `internal` (`driver_impl.go:542`,
`trajectory_verify.go:106`, `evidence_verify.go:368`, `source.go:252`, `verification.go:251`,
`acp/spawn.go:178`, `acp/shellenv.go:44`, `:120`) — none on the reservation, parent-recovery or
unchanged chains, so your reachable-today cost statement holds.

**I accept the authority reading, and I state plainly why I cannot honestly contest it.** Your
round-05 `WRONG` verdict on my round-04 equivalence claim was never rebutted, and my own round-05
verification of `os.Root.Rename`'s Windows implementation (handle-relative `renameat` with
`FILE_OPEN_REPARSE_POINT`; `root_windows.go:386`, `at_windows.go:339-429`) is precisely the
machinery a path-based `MoveFileEx` with a lexical sibling check bypasses. Converting a rooted
site is a weakening of `os.Root` containment for the publication operation; the owner's answer
postdates my round-05 file and reserves it to an explicit deviation. This is not a technical
stalemate — the durability merits resolved my way and you conceded them — it is an authority
boundary, and the deck cannot grant it to itself. Your testable blocker (FINAL may not adopt
path-based conversion at A2/A3/B3/B4/B5/C1/C2 without an explicit owner deviation) is **accepted
by me**, which closes it: all three participants now hold it.

**Your `N8` is `CONFIRMED` (`PRIMARY`)**: `publishParentRecoveryWithSync` binds `syncParentRecovery`
as a function value at `:311`, renames at `:333`, and calls `sync(dir, base)` at `:336`;
`publishRecoveredParent` renames at `:570` and calls `syncRecoveredParent` at `:573`. A refusal
emitted at the barrier fires after publication, so any refusal disposition at C1/C2 must be
evaluated before the rename — carried into my table. **Your `N9` is `CONFIRMED` as scoped
(`PRIMARY`)**: both handoff files exist; the designated-implementer handoff (12075 bytes, mtime
2026-09-25 11:13) carries `status: agent-controlled-delivery-complete`, `blocking: no`, states
"This file releases the sequencing hold for `windows-portability`", and reports CLI 1.50.0 and
skill 2.14.0 released on GitHub and Homebrew. I adopt your caution verbatim: I verified existence
and self-testimony, not the release it describes; the `00-prompt.md` gate condition is file
existence and is met. Consequence adopted: the version anchor moves from "above 1.49.1" (all three
round-05 files, mine included — stale) to **above 1.50.0**. **Your `N10` is `CONFIRMED`
(`PRIMARY`)**: `grep -rnE '^func sync'` reproduces your table — `trajectory_verify.go:132`,
`verification.go:183`, `unchanged.go:267`, `parent_recovery.go:295`, `:532`,
`reservation_recovery.go:55` — and all three round-05 files (mine included) cited the `SyncFile`
call lines inside the bodies as "def". Both numbers real, label wrong, no disposition change; the
fifth census-class error, three of them mine, and the final reason FINAL's table is generated from
the enumeration, never transcribed.

### @kimi-1 — round-06 (current position)

**Both your direct-to-me points are answered: yes on authority, conceded on replay.** (1) Yes — I
accept that the rooted conversion is owner-reserved; your exact testable blocker ("whether
zcode-1's round-06 accepts that the rooted conversion is owner-reserved") resolves in the
affirmative, by acceptance rather than by contest. My verification basis is above: my own round-05
containment sentence conceded the exception, claude-1's `WRONG` verdict on the equivalence claim
was never rebutted and my own handle-rename verification is the evidence for it, and the owner's
answer reserves the decision. (2) Conceded — my destination-exists replay claim is withdrawn
(position change 1); your `WRONG` verdict and claude-1's are correct and my own re-reading of the
two comments this round is the counterexample. The deterministic-stage condition you adopted
repairs my design; with it my round-05 restructure stands at class A intact, which is the outcome
I wanted plus the invariant I missed.

**Your C-row REPLACE refinement is accepted as a FINAL drafting requirement on the deviation
branch.** I re-read both C sites (`PRIMARY`): today they use `os.Root.Rename` — replace-permitted
at the syscall level — with publish-once enforced by caller validation (`parentRecoveryName` at
`:19`, `recoveredParentName` at `:405`, both constants; "never rewritten here" at `:548-550`'s
comment). My round-05 line "rename becomes `ReplaceSyncedFile`-pattern WT-move REPLACE" inherited
the flag by pattern-match from a state-file contract designed for replacement. You are right that
at C the flag is a deliberate choice, not an inheritance; FINAL (deviation branch only — under the
current design C1/C2 refuse) must state which flag, why, and the N6 existing-directory-error
behaviour either way. My initial drafter's choice, stated now so peers can veto it in review:
no-REPLACE at C1/C2, matching the publish-once constants — an existing final artifact should be a
loud replay signal, not an overwrite.

**Your D4 split, row-ID scheme, and 15/16 reconciliation are `CONFIRMED` (`PRIMARY`,
re-derived independently this round).** The 15 name-keyed call lines are
`trajectory_verify.go:125/:154/:193`, `verification.go:211/:267/:286`, `unchanged.go:241`,
`parent_recovery.go:368/:573/:605`, `reservation_recovery.go:47/:67/:84/:420/:431`; `:67` sits
inside `syncIntent`'s own body; adding the `:336` function-value invocation gives **15 independent
sites**, and B1's `:117-128` two-iteration loop gives **16 runtime operations**. Your round-05
guard set (1)–(4) survives, as you say, as the **content of the deviation request**, and I carry
it into the deviation specification below with the additions from all three files.

**Your publish-time `validHash` gate at A3 is confirmed and becomes doubly required** (`PRIMARY`):
`publishReservationIntent` at `:70-75` builds `intentName(i.Accounting.EntryKey)` with no
`validHash` gate, while `readReservationIntent` gates at `:92-93`; under any restructure the entry
key feeds both stage and final names, so the gate moves to publish time. One sentence in FINAL.

## The three conflicts — my current position

### Conflict 1 — create-entry guarantee: `O_SYNC` inference vs staged write-through rename — **RESOLVED**

The factual half closed identically in all three round-05 files: no documentation scoped to the
parent-directory entry of a create exists; `O_SYNC`→`FILE_FLAG_WRITE_THROUGH` is real and
source-settled for **write** durability only; `FlushFileBuffers` on our read-only directory
handles is foreclosed twice over (no documented directory subject; `GENERIC_WRITE` precondition
unmet) and stays diagnostic-only. All three of us now select the staged write-through rename for
class A. **My position, current: restructure class A — stage (`O_EXCL` + mandatory file fsync) +
`MoveFileEx(MOVEFILE_WRITE_THROUGH)` without `MOVEFILE_REPLACE_EXISTING` — with the
deterministic-stage condition**: stage name a fixed derivation of the final name, `O_EXCL` on the
stage, stage retained on the error path, so an interrupted publication leaves a visible blocker
whose existence prevents replay, exactly as `O_EXCL` on the final name does today. Two FINAL
requirements I add as drafter: (a) the stage derivation must be provably non-colliding with every
valid final name — final names are validated or product-constructed (`validHash`-gated `<hex>.json`
at A2/A3), so the stage suffix must sit outside every valid-final grammar, pinned by a unit test;
(b) the publish-time `validHash` gate at A3. The no-clobber basis stays disclosed per kimi-1's
round-05 N6 limitation: flag-table contrapositive plus a hosted existence-probe assertion, never
upgraded to "documented".

### Conflict 2 — preserved `os.Root` containment vs feature scope; pre-mutation refusal — **RESOLVED as an authority question; plain/rooted split**

The merits are settled: a lexical sibling check plus a path-based call is not equivalent to
`os.Root`'s kernel-enforced, reparse-resistant, handle-relative containment, so converting a
**rooted** site is a weakening the owner has reserved to an explicit deviation. My position,
current, identical to both peers':

- **Plain sites restructure, no owner question**: A1, B1, B2 (all currently dormant behind the
  `trajectory_verify.go:106-108` POSIX refusal — this branch is the design of record if that gate
  ever lifts); **D1 derives** (plain, and its original publication `writeState` already ends in
  `ReplaceSyncedFile` = the reviewed WT-move). B-row note: B1/B2 are idempotent ensure-exists, so
  EEXIST-tolerant recheck applies there.
- **Rooted sites take the named blocking refusal pending an explicit owner deviation**: A2, A3,
  B3, B4, B5, C1, C2, with D2/D3/D4a/D4b following their originals into refusal (my own round-5
  coherence note, applied honestly to this branch). Refusal at C1/C2 is **ordered before the
  rename** (N8). The refusal is named, pre-mutation, user-visible — a first-class outcome under
  the original owner scope ("a real Windows implementation **or** an explicit, reviewed,
  user-visible refusal").
- **Reachable-today cost, stated plainly**: precharge reservation-intent publication (A3+B5, via
  `openIntentRoot(b, true)` ← `PrepareCycleReservation` `reservation_recovery.go:126` ←
  `cycle_binding.go:269`, no GOOS gate) and parent-recovery / recovered-parent apply (C1+C2,
  dragging D2/D3; D4a/D4b track A3/B5). Trajectory verification and captured journals are already
  POSIX-refused on Windows today, so the refusal branch takes away nothing that currently works.
- **Deviation question, fully specified for the organizer to carry** (I am a participant; I do
  not request or grant it): content = kimi-1's guard set (1)–(4) (root-resolved paths via
  `Root.Name()`, retained sibling check, every rename name a product constant or `validHash`-gated
  including the publish-time gate at A3, randomized stage siblings at C / deterministic at A,
  explicit weakened-resistance disclosure) + claude-1's items (the precise statement of what is
  given up — kernel-enforced per-open reparse resistance and handle-relative rename, replaced by a
  lexical check plus a path-based call with a TOCTOU window; the unselected handle-derived
  alternative named and why; both branches' costs side by side) + my hosted rooted-escape test
  against the conversion helper (mechanics only, never durability evidence). If granted, every
  refusal row converts per the already-designed restructure and the D rows derive; if declined,
  the refusal-branch table is final as written — a reviewed outcome, not a failure. As drafter I
  record a preference honestly and only as input: I would prefer the deviation granted, because it
  recovers two Windows features under reviewed guards; the choice is the owner's.

### Conflict 3 — one normalized operation table — **RESOLVED**

The three inventories were always the same set; only individuation differed. Re-derived
independently this round (`PRIMARY`): 15 name-keyed call lines − 1 in-body (`:67`) + 1
function-value site (`:336`) = **15 independent barrier call sites**; B1 fires twice per run =
**16 runtime operations**. 14 + (D4 split) = 15 + (B1 second iteration) = 16. Row IDs are
kimi-1's scheme (strict file-then-line order). Mapping of my round-05 14-row table, nothing
dropped: A1→A1, A2→A2, A3→A3, **B0→B1, B1→B2, B2→B3, B3→B4, B4→B5**, C1→C1, C2→C2, D1→D1,
D2→D2, D3→D3, **D4→D4a+D4b**. The audit of record remains my `fsutil.SyncDir` named-type contract
(compiler-enforced, never a grep); FINAL carries the derivation command and this table, never a
headline count.

## The normalized operation table

Abbreviations: **WT-move** = `MoveFileEx(MOVEFILE_WRITE_THROUGH)`; **no-REPLACE** = without
`MOVEFILE_REPLACE_EXISTING`; **det-stage** = deterministic stage derived from the final name,
`O_EXCL`, retained on the error path. "Refusal branch" is the current design; the "if dev." column
is the conversion that applies only if the owner grants the containment deviation.

| ID | Site (publication → barrier) | Class | Rooted? | Win-reachable | Required guarantee | Refusal branch (current design) | If dev. granted | Blast radius if refused | Partial-publication recovery | r05 map (zcode/claude/kimi) |
|---|---|---|---|---|---|---|---|---|---|---|
| **A1** | `trajectory_verify.go:146` `os.OpenFile O_EXCL` → `:154` | A create | plain | no — gated `:106-108` | create entry durable; anti-replay | **Restructure**: det-stage `O_EXCL`+fsync → WT-move no-REPLACE | (same) | none today; inherits if gate lifts | stage persists, blocks replay; final never torn | A1 / #3 / A1 |
| **A2** | `verification.go:202` `Root.OpenFile O_EXCL` → `:211` | A create | **rooted** | no — gated `:251-253` | same + `:192-193` invariant | **Refuse**, pre-mutation | det-stage → WT-move no-REPLACE, guard set | none today | `:192-193` invariant intact (unconverted) | A2 / #5 / A2 |
| **A3** | `reservation_recovery.go:75` `Root.OpenFile O_EXCL` → `:84`→`:67` | A create | **rooted** | **yes** | same + `:79` invariant | **Refuse**, pre-mutation | as A2 + publish-time `validHash` gate | precharge reservation-intent publication refuses → `PrepareCycleReservation` fails (`cycle_binding.go:269`) | nothing published; `:79` invariant intact | A3 / #9 / A3 |
| **B1** | `trajectory_verify.go:118` `os.Mkdir` → `:125` — **fires ×2** (`:117-128` loop) | B mkdir | plain | no — gated `:106-108` | base-dir entries durable, idempotent | **Restructure**: stage-`Mkdir` + WT-move + EEXIST-tolerant recheck | (same) | none today | IsExist tolerance preserved; nothing published | B0 / #1+#2 / B1 |
| **B2** | `trajectory_verify.go:190` `os.Mkdir` → `:193` | B mkdir | plain | no — gated `:106-108` | run-dir entry durable | **Restructure** as B1 | (same) | none today | as B1 | B1 / #4 / B2 |
| **B3** | `verification.go:264` `parent.Mkdir` (IsExist-tolerant) → `:267` | B mkdir | **rooted** | no — gated `:251-253` | store entry durable, idempotent | **Refuse**, pre-mutation | stage-`Mkdir` + WT-move + tolerant recheck | none today | as B1 | B2 / #6 / B3 |
| **B4** | `verification.go:283` `base.Mkdir(charge)` — anti-replay token `:281-282` → `:286` | B mkdir | **rooted** | no — gated `:251-253` | reservation entry durable; never reuse | **Refuse**, pre-mutation | det-stage-`Mkdir` + WT-move no-REPLACE; existing dir = loud even under REPLACE (N6); EEXIST on stage = loud, matching today's `Mkdir(charge)` failure | none today | `Mkdir` failure = "already reserved" `:284` | B3 / #7 / B4 |
| **B5** | `reservation_recovery.go:36` `dir.Mkdir` (IsExist-tolerant) → `:47` | B mkdir | **rooted** | **yes** | intent-root entry durable, idempotent | **Refuse**, pre-mutation | stage-`Mkdir` + WT-move + tolerant recheck | `openIntentRoot(create=true)` fails → reservation creation unavailable | nothing published; `Mkdir` idempotent | B4 / #8 / B5 |
| **C1** | `parent_recovery.go:333` `dir.Rename` → `:336` (func value bound `:311`) | C rename | **rooted** | **yes** | rename entry durable; publish-once | **Refuse, ordered before `:333`** (N8) | WT-move; flag chosen deliberately (my draft: no-REPLACE), disclosed — not inherited from `ReplaceSyncedFile` | parent-recovery **apply** refuses | pre-rename refusal ⇒ nothing published; `defer dir.Remove(stage)` `:328` clears stage | C1 / #12 / C1 |
| **C2** | `parent_recovery.go:570` `dir.Rename` → `:573` | C rename | **rooted** | **yes** | same | **Refuse, ordered before `:570`** | as C1 | recovered-parent **apply** refuses | as C1; stage cleared `:565` | C2 / #14 / C2 |
| **D1** | `unchanged.go:241` → `:279` (func `:267`) | D re-sync | plain | **yes** | already satisfied | **Derived no-op**, argued per site: `writeState` (`state.go:197-221`) already publishes via `ReplaceSyncedFile` = WT-move | (same) | none | `os.CreateTemp` + `defer os.Remove` (`:205`/`:209`); replace atomic | D1 / #16 / D1 |
| **D2** | `parent_recovery.go:368` → `:308` | D re-sync | rooted ctx | **yes** | inherits C1 | **Refuse** — derives from a refused publication | derives once C1 converts | tracks C1 | re-sync only; never rewrites | D2 / #13 / D2 |
| **D3** | `parent_recovery.go:605` → `:545` | D re-sync | rooted ctx | **yes** | inherits C2 | **Refuse** | derives once C2 converts | tracks C2 | re-sync only | D3 / #15 / D3 |
| **D4a** | `reservation_recovery.go:420` `syncIntent` → `:67`+`:189` | D re-sync | rooted ctx | **yes** | inherits A3/B5 | **Refuse** | derives once A3/B5 convert | tracks A3 | re-sync only; `:426` never rewrites bytes | D4 / #10 / D4 |
| **D4b** | `reservation_recovery.go:431` `SyncFile(trajectory.json)` + `syncVerificationDirectory` | D re-sync | rooted ctx | **yes** | inherits A3/B5 | **Refuse** | derives once A3/B5 convert (`trajectory.json` via `writeState`→`ReplaceSyncedFile` already) | tracks A3 | re-sync only | D4 / #11 / D4 |

**Totals:** 15 IDs; 16 runtime operations (B1 ×2); 9 of 15 Windows-reachable today (A3, B5, C1,
C2, D1, D2, D3, D4a, D4b); 6 dormant behind the two existing reviewed POSIX refusals. Current
design: 3 plain rows restructure (A1, B1, B2), D1 derives, 7 rooted rows refuse pending the owner
deviation, 4 class-D rows follow their originals. All conversions are Windows-confined; the Unix
path at every row is byte-identical to today; Windows `fsutil.SyncDir` remains the fail-closed
named-refusal emitter for any site not carrying a proved mechanism, including future sites added
without conversion — under the refusal branch it is also the rooted sites' refusal emitter.

## Census baseline (settled, re-pinned)

Re-executed at `9e55050` (`PRIMARY`): **104** `t\.Skip` pattern lines under `internal`, **103
call sites** = 68 `t.Skip(` + 35 `t.Skipf(` + 0 `t.SkipNow(`, **56 files** matching the pattern /
**55 with call sites** (sole non-call `internal/runner/launch_precheck_skip_test.go:58`) —
identical to all three round-05 files and both round-06 files. Emission: `-json` on all three
legs; every firing Windows SKIP matched row-by-row to an individually reviewed exclusion or the
leg fails; AF_UNIX fallbacks become hard failures; macOS/Ubuntu current skips are the enumerated
baseline and any newly firing skip fails the leg.

## §15 bookkeeping

**§15.3.** No `DISPUTED` claim survives into my acceptance. My destination-exists replay claim is
closed by my own weakening (position change 1) after two concordant non-owner `WRONG` verdicts —
closed, not out-voted; I verified the counterexample myself before withdrawing. My
restructure-everywhere selection is closed by self-correction under the owner's new authority. No
acceptance criterion in my proposal rests on any disputed claim; `consensus.md` needs no
`## Verdict conflicts` section on my account.

**§15.6(b), recorded because it binds and cuts against all three of us.** This round is sequential
by the owner's design: each participant read the prior positions before writing, so the full
three-way convergence on all three conflicts is **not independent evidence**. claude-1 joined the
class-A restructure after reading both round-05 files; kimi-1 joined the plain/rooted split after
reading claude-1's round-06; I join both after reading both round-06 files. `consensus.md` must
record that the deck holds no dissenting position on any of the three branches and state what
would make the agreed position wrong: a non-NTFS runner volume; a directory move rejected on the
runner; no-replace clobbering an existing destination on the runner; the deterministic-stage
suffix colliding with a valid final name (prevented by requirement (a) in Conflict 1); or the
owner declining the deviation (then the refusal branch stands — a reviewed outcome, not a failure).

## Consensus readiness and exact remaining blockers

**Consensus-ready from my side on all three conflicts.** My positions are now identical to both
peers' round-06 positions: class A restructures with deterministic stages; the plain/rooted split
governs containment (restructure plain, refuse rooted pre-mutation pending an explicit owner
deviation, D rows follow their originals); one normalized 15-ID / 16-operation table on kimi-1's
IDs. With this file, no participant holds a dissenting position on any of the three conflicts, and
the crossing failure mode the owner authorized this round to fix is closed — I have treated only
round-06 positions (and my own round-05, corrected) as current, per the owner's instruction.

**Exact remaining item — one external fact, zero inter-participant blockers:** the owner's
deviation decision on the seven rooted sites. The FINAL I will draft, when the organizer opens
consensus, is written to the **refusal branch** — the strongest position the deck may take on its
own authority under the owner's answer — with the deviation branch fully specified as the
conversion table that applies if the owner grants it. That is the same signoff stance kimi-1
recorded and compatible with claude-1's blocker as now closed by my acceptance. I request no
seventh round and no phase advance; carrying the costed deviation question to the owner is the
organizer's, not mine.

**Release gates, re-checked read-only this round:** both predecessor handoff files exist (N9,
existence + self-testimony verified); version selection is anchored above **1.50.0** (all three
round-05 files, mine included, are stale on this). The Windows experimental label and the held CLI
winget PR stay until the released commit's hosted Windows leg is green; no tag moves; a
participant independently verifies every released channel. `zcode-1` remains drafter and single
implementation owner, claude-1 and kimi-1 reviewing.

## Scope compliance

Read-only verification only: repo source at `9e55050` (0 `.go` files changed since `6b87cf8`),
`git` metadata, `grep`/`sed` over `internal` (census, GOOS gates, helper `func` lines, call-site
listing, replay comments, C1/C2 ordering), the skip census re-executed, `shasum -a 256` of the
protocol packet, and read-only `ls -l` plus frontmatter/heading reads of the two release-gate
handoff files. No Microsoft Learn page re-fetched — documentation basis tagged `PRIMARY(r5)` or
`SECONDARY` throughout. Wrote exactly one file:
`parley-deck/ideas/windows-portability/round-06/zcode-1.md`. No product code, no workflow, no peer
artifact, no signoff, no phase transition, no consensus or FINAL, no commit, no tag, no
publication, no hosted run, no hosted mutation, no driver or second organizer launched, no
seventh-round request, no owner deviation requested. No secrets. No browser, no Chrome. English
only.
