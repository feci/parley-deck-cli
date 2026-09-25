---
agent: claude-1
idea: meta-protocol-change-designated-implementer
review-round: 4
date: 2026-09-25
reviewed-commit: 717f3debe3aabb0f7a702d689e8d10de1e04fa8c
reviewed-commit-skill: a624318dcda02c47ecaa859d987efee08dba3104
implementation-record-commit: ee8849c9ba6f75748fb471fe9c910bf362300e8d
prior-reviewed-commit: 1bad2634380dd785009987a5c50427050efb87c3
prior-record-commit: b850576ba3f943ba9e4c942247e864982d5ccfdc
baseline-commit: e4640bf2840249db0c1a1ecab7493813f4dacfdb
artifact-hashes: IMPLEMENTATION.md@ee8849c sha256 86fd4cb58a5b4f61fd94b46efd61ae488c2db13030f641aac463ae85150c44ba (728 lines / 55,940 B; on-disk identical); review/consensus.md sha256 5da26e8892c5ba42f93e5f0439cf5cf028b73717590b6ee3ed9b4d309ab90923; internal/app/driver_designation_test.go@717f3de sha256 292e16ed2b2b115cfd938f3dc62095622db4b5f5c9a371077262fc0929a02ae8 (@1bad263 504958434e35637c388808ce23e679640fc0c093449df30da0e60369285679a4); COOPERATION.md sha256 b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388
reviewed-against: the signed cycle-3 review/consensus.md (claude-1 / kimi-1 / zcode-1, all ✅ ACCEPT, no BLOCK) + frozen FINAL.md (e4640bf)
---

Independent non-implementer re-review of Phase-8 fix-up **cycle 3** as an ordinary participant.
Source under review: `717f3de` (one Go file, one comment line). Canonical record: `ee8849c`
(`IMPLEMENTATION.md` only, no Go content). Skill unchanged at `a624318` — skill worktree HEAD
re-verified by me at the full SHA. Prior reviewed state `1bad263` / `b850576`. I re-derived AF-16,
AF-17 and the standing bump from the trees and the live CLI, not from the record. The frozen FINAL
and the owner boundaries are controlling; I re-open no settled design. codex-1 is a pure organizer
here and kimi-1 the sole implementer; I sign no consensus in this invocation and commit nothing.

## Protocol context attestation (Phase 8, live, original workspace)

Run by me this session from
`/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/designated-implementer`, parley 1.49.1,
exit 0:

    parley protocol packet --dir . --phase 8 --track deliberation \
      --idea meta-protocol-change-designated-implementer \
      --flag auto_implement --flag protocol_change --audience participant --json

    context_mode  : full
    request       : phase 8, track deliberation, transport github-pr,
                    flags ["auto_implement","protocol_change"], optimize false,
                    audience participant
    source        : parley-deck/COOPERATION.md (role source, transport github-pr, 115,166 B)
    source_sha256 : b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388
    packet_sha256 : b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388
    shadow        : fe0e4c046897…bcf — 86,716 B, 40 included / 29 omitted (not used)

`fallback_reason` **ABSENT** — top-level keys enumerated: `body_path, context_mode, index,
packet_sha256, request, shadow, source, source_sha256`. `source_sha256 == packet_sha256`, so full
mode, no substitution. PRIMARY cross-check: `shasum -a 256 parley-deck/COOPERATION.md` returns the
same digest at 1,400 lines / 115,166 B — **this cycle moved no protocol text**, so the authority is
identical to the one every cycle-2/cycle-3 attestation used. Phase 8 (fix-up, per-publication
frontmatter bump, stopping judgment, LE-7) and §15 were read from this body. `strict_gate` is not
set in `00-prompt.md` (re-read: `facilitator: codex-1`, `participants: [claude-1, kimi-1, zcode-1]`,
`auto_implement: true`), so the **default** Phase-8 close rule applies: a NIT is dispositioned in
full but is not automatically blocking.

## Method and isolation

Isolated local-disk checkout built by me from the commit object — never a peer tree, never the
shared mount: `git archive 717f3de… | tar -x` → `/tmp/parley-claude1-r4-66202/cli`, go1.27.1
darwin/arm64. Packet runs and one shared `parley wait` wrote only to git-ignored `.parley-runtime/`.
I altered no source, record, peer file, signature, FINAL, global setting or release. Prior evidence
is referenced, not reprinted; every reused figure carries its exact commit.

## Findings

**1 NIT. Zero MAJOR, zero MINOR, zero behaviour.** Both agreed fixes landed, and the standing bump
landed with its claimed machine effect. The one finding is new text introduced by AF-16 itself.

### NIT-1 (new) — AF-16's replacement lead-in installs a fresh over-broad universal, falsified by a figure in its own sentence

**Provenance: PRIMARY**, my own read of `ee8849c` plus my own packet measurements below.

`IMPLEMENTATION.md:325-326` now opens the relabelled bullet with:

> **All shadow byte figures in this bullet are rendered WITH `--flag auto_implement --flag
> protocol_change`**:

The bullet contains exactly four shadow byte figures — 86,336 / 86,716 / 86,701 / 80,799 — and the
clause immediately after the colon (`:328-329`) says 86,701 is "the SAME phase-8 packet over the
SAME `d238238` source rendered **WITHOUT** the two flags". So one of the four falsifies the bolded
universal, inside the same sentence.

This is the *same defect shape* AF-16 was written to remove, one dimension over: the pre-AF-16 text
read "**All shadow byte figures in this bullet are PHASE-8 shadows**" and was falsified by the
phase-7 figure later in the same bullet. AF-16 struck that universal and installed a flags-universal
with the identical flaw. Recorded plainly: the shape is carried forward, not invented, but the
flags wording is new at cycle 3, so it is this delta's text.

Why I file it rather than waive it: in this idea the pattern has *demonstrated* propagation — the
unqualified phase claim is what AF-14 then reasoned from into the "did not reproduce / superseded"
error that became my round-03 MINOR-1. A future cycle quoting the bold line alone would repeat it.

Why it is a NIT and not a MINOR: **AF-16's signed observable check passes in full.** Neither
location asserts non-reproducibility or supersession; 86,716, 86,701 and 80,799 each name their flag
state explicitly; and a reader who runs the flagless phase-8 command finds 86,701 *documented*, not
contradicted. The residual harm is limited to the blanket header, which its own sentence corrects,
and to 86,336 — the one figure whose flag state rests solely on the defective universal (its
flagged provenance is independently visible in the command quoted at `:313-315`). The closing rule
"Read the phase and the flags with the figure" is correct and complete.

Repair is two words, no re-measurement: "**Except where stated**, all shadow byte figures in this
bullet are rendered WITH …". I would not block a close over it; a recorded dismissal is an equally
honest disposition.

## Refutation attempts

Each run by me at `717f3de` in my own checkout or against the live CLI. "Held" = the implementation
survived.

**A. Is AF-17's comment actually applied, and is the tree-wide sweep clean?** `grep -rn "any run"
--include="*.go" .` at `717f3de` returns exactly four hits:
`driver_designation_test.go:448` "…on any run **that reaches a role action**." (fixed),
`driver_impl.go:188` "Hard gate on any run that reaches a role action" (AF-13, qualified),
`internal/agents/naming.go:50` and `internal/budget/run_identity_inventory_test.go:351` (unrelated
English). The signed observable check is satisfied exactly. **Held.**

**B. Did AF-17 reuse the AF-13 wording, or quietly expand it?** The two carriers now share the
phrase "any run that reaches a role action" verbatim. No enumeration was added — consistent with the
consensus's own dismissal of my round-03 refutation-A observation. **Held.**

**C. Is the source change truly comment-only — behaviour, assertions and fixtures untouched?**
Mechanical proof, not inspection: I parsed `driver_designation_test.go` at both `1bad263` and
`717f3de` with `go/parser` in mode 0 (**comments dropped**), cleared `Comments`, and printed each
AST. Both hash to `efe7f11aded3650d7d5f8c209094ba19cd209dcaf5f34140deb8afaebe96fdae` — **the
comment-stripped syntax trees are byte-identical**. Zero executable difference. **Held.**

**D. Do the flag labels AF-16 asserts reproduce?** PRIMARY, live source `b273af1e…f388`, parley
1.49.1, one argument varied at a time:

| phase | both `--flag`s | shadow bytes | shadow sha256 | blocks |
|---|---|---|---|---|
| 8 | present | 86,716 | `fe0e4c04…bcf` | 40 incl / 29 om |
| 8 | absent  | **86,701** | `a0845614…8efd` | 40 incl / 29 om |
| 7 | present | 80,799 | `c6d29141…c980` | 41 incl / 28 om |
| 7 | absent  | 80,784 | `70115673…d6de` | 41 incl / 28 om |

Every figure the record labels reproduces at its labelled flag state, including the phase-7 "also
rendered with both flags" claim. **Held.** New corroboration for VC-3.2: I rebuilt a probe deck
carrying the **pre-edit** protocol source (`git show 0893989:parley-deck/COOPERATION.md` →
`c749218…568f`, 1,400 lines / 114,771 B — the record's `:318-319` claim, re-derived and correct) and
the +15 B flag delta holds there too (86,296 flagged / 86,281 flagless). The discriminator is the
flag pair, and it is source-version-independent. This is a third participant-independent route to
the conclusion all three of us signed.

**E. Is 86,336 itself the flagged rendering?** **Not independently re-derivable, and I say so.** My
probe deck reproduces the flag *delta* over the pre-edit source but not the absolute figure
(86,296 ≠ 86,336), because the cycle-1 deck state — other ideas, inbox, runs — contributes to the
packet and is not reconstructible today. Both readings stay arithmetically consistent, so the label
rests on the bullet's own quoted command, which does carry both flags. Named, not waived; it is the
one part of NIT-1's surface I could not close by measurement.

**F. Did the standing bump land, and does it still produce a machine effect?** Frontmatter at
`ee8849c`/worktree: `status: fix-up-cycle-3`, `head-commit: 717f3de`, `skill-commit: a624318`.
`717f3de` is `ee8849c`'s **parent** — PRIOR and knowable, no self-hash claim; the cycle-3 section's
own `record-commit` is correctly left as "named at the next touch". The cycle-2 section's
`record-commit: b850576` is filled in and `b850576` **is** the cycle-2 record commit. Live
confirmation from the original workspace (`parley wait --for review --timeout 15s`):
`implementation: present=true status=fix-up-cycle-3 implementer=kimi-1` / `next: await review
artifact`. **Held.**

**G. Was the deliberately-untouched AF-14 summary touched?** Extracted by content from both record
blobs and compared: the `- **AF-14** — record prose precisions…` bullet is **byte-identical**
(`174b6a8f…b376`) at `b850576` and `ee8849c`, exactly as the signed plan required. The record diff
has no hunk in that region. **Held.**

**H. Is the delta confined, and was anything signed edited to make it fit?**
`git diff --name-only 1bad263 717f3de -- internal/ cmd/` → one file;
`git diff --name-only 717f3de HEAD -- internal/ cmd/` → **empty**. `717f3de` touches only the test
file; `ee8849c` only `IMPLEMENTATION.md`. Hashed against the source commit's parent `770eb9b`:
`FINAL.md`, `00-prompt.md`, `review/consensus.md`, `consensus-cycle-01.md`, `consensus-cycle-02.md`
and all six `review/round-0{1,2,3}/*.md` are **byte-identical** to the worktree. The archives still
hash to the ratified `48a32a02…a15611` (cycle 1) and `dcc1e556…f865` (cycle 2). `FINAL.md` is
unchanged against `e4640bf`. The only later commit (`357d1b5`) is organizer-only
(notes/usage/ledger). **Held.**

**I. AC-4 default-path pins and the shipped default.** `TestUnsetPathIsByteIdentical` is unmodified
(absent from the diff) and PASSes at `717f3de`. No `default_implementer` is set in any `*.toml` in
the tree, nor in `~/.parley/agents.toml` — **the product default still ships UNSET**. **Held.**

**J. Does the record over-claim its evidence?** No. The AC-17 paragraph states that my full suite
stays exact-commit scoped to `d238238` and is not re-claimed, and that the targeted behaviour
evidence stays scoped to `1bad263`. That is accurate and is the honest scoping I would have asked
for. **Held.**

**Considered and declined as findings** (disclosed, not suppressed): (i) "The source hash, line
count and byte count are identical on both invocations" at `:618-619` sits beside two differing
*shadow* byte counts — but "source" heads the coordinated list and the triple is the record's
standard form for 115,166 B / 1,400 lines; not filed. (ii) The cycle-3 sweep parenthetical
enumerates three of the four `any run` hits, omitting the `:448` line it just fixed — but the claim
"leaves only qualified or unrelated hits" is true, the omitted line is the subject of the preceding
sentence, and the enumeration reproduces the signed plan's observable check verbatim (my own
drafting, not the implementer's); not filed.

## Dispositions and concurrences

| Item | Disposition at `717f3de` / `ee8849c` | Evidence |
|---|---|---|
| **AF-16 (1)** `:325-334` | **APPLIED** — supersession clause struck; 86,716 "WITH", 86,701 "WITHOUT", 80,799 "also with both flags"; rule completed to "phase **and the flags**". Residual over-broad lead-in → **NIT-1** | refutations A/D/E |
| **AF-16 (2)** `:613-619` | **APPLIED exactly as signed** — retitled "One plan figure differs by invocation flags, disclosed", non-reproducibility claim dropped, disclosure kept | record read; refutation D |
| **AF-16** untouched AF-14 summary | **HONOURED** — byte-identical | refutation G |
| **AF-17** `driver_designation_test.go:448` | **APPLIED exactly as signed**, wording reused not expanded, comment-only proven | refutations A/B/C |
| **Standing bump** (AF-12's restored duty) | **APPLIED**, all four values correct and PRIOR, machine-confirmed | refutation F |

I **concur** with the signed plan's dispositions and carry them unchanged: **VC-3.1** (zcode-1's
scoped null and my two findings are a scope divergence, not a conflict of fact); **VC-3.2** (the
flag pair is the discriminator — now corroborated by a fourth independent run set of mine over a
*second* protocol source); all three **dismissed findings** (short-SHA convention; pre-existing
fail-closed `confirmed-unconfirmed`; the "role action" non-enumeration I declined to file and still
decline); and **DF-1…DF-4** unchanged, none quietly closed, widened or converted into silent work.
No new deferral and no new owner question arises.

## Limitations (named, not waived)

1. **No broad `go test ./...` at `717f3de`** — deliberately, per the ratified proportional plan for a
   comment/prose-only delta, and none is owed absent a new behaviour concern. NIT-1 is prose, so I
   request none. My independent full suite stays **exact-commit scoped to `d238238`**; the targeted
   behaviour evidence (both reviewers' focused runs, M1=9 / M2=5, the probe sets) stays scoped to
   **`1bad263`**. Neither is re-claimed here. What I *did* run at `717f3de`, in my own isolated
   checkout: `go build ./...` exit 0; `go vet ./...` exit 0, no output; `gofmt -l
   internal/app/driver_designation_test.go` clean; `go test ./internal/app/ -count=1 -run
   'TestMalformedTier2Gates|TestTier2UnavailabilityGateAndExits|TestPinDesignationConflictEscalates|TestUnsetPathIsByteIdentical'`
   → `ok parley-deck-cli/internal/app 0.447s`, all four PASS.
2. **86,336 not absolutely re-derivable** (refutation E).
3. **Self-disclosure, my own frozen artifact:** the block-split column of the table in
   `review/round-03/claude-1.md` prints `40/69` and `41/69`; the measured splits are 40 incl / 29 om
   and 41 incl / 28 om, as the consensus table and the record both correctly carry. A transcription
   slip in a filed, frozen review with no downstream effect — no edit is proposed or authorized.
   (This is a second such artifact alongside the stale `bytes=28911` the consensus already recorded.)
4. **Carried blind spots, unchanged from cycles 2–3:** skill npm-level gates (`npm test`, manifest
   `--check`, `npm pack --dry-run`, installer integrity) still owed to the organizer's release
   preflight; no end-to-end `parley run` against a real designated deck; no live §9.0 ping behind
   `designeeAvailable`; TUI and pipeline-block dispatch (FINAL F3/F10); Windows (DF-4);
   AC-5…AC-15/AC-18/AC-20 at regression level. `ee8849c` holds no Go content, so no suite result
   extends to it — it was reviewed by reading, which is again how the only finding was found.
5. **Pre-existing, reported not blocking:** the unanswered
   `claude-to-user_…_driver-error.md` to-user escalation, and `parley wait` printing
   `reservations=[]` — both the organizer's thread, neither this delta's.

## Trajectory note (Phase-8 stopping judgment)

17 findings (cycle 1) → 8 (cycle 2) → 2 (cycle 3) → **1 NIT** (cycle 4), and this one is a residue
of the cycle-3 repair text rather than new ground. Zero MAJOR, zero MINOR, zero behaviour across the
last two cycles.

## Verdict

**PASS with 1 NIT — no BLOCK, no reservation.** Both agreed fixes and the standing bump landed
exactly as signed; the delta is confined; the frozen FINAL, the signed plan, all three signatures
and every peer artifact are byte-untouched; the AF-14 summary is byte-identical as required; the
default-path pins hold and the product default still ships UNSET. AF-17 is proven comment-only by
AST identity, so nothing shipped changed behaviour.

NIT-1 is **not blocking** under the default Phase-8 close rule and I do not ask for another
implementation cycle on its account: it is a two-word clarification or an explicit recorded
dismissal, and the cycle-4 consensus may take either path. If it is fixed, it is a record-only edit
in one location; if it is dismissed, the dismissal should say so plainly rather than leave the
universal unremarked, because this idea has already watched that exact shape propagate once.

I am a reviewer, not the organizer: this verdict is my own, I sign no consensus in this invocation,
and the zero-fix consensus plus the separate LE-7 goal-done check by a **fresh invocation of an
existing non-implementer quorum member** remain owed and unclaimed. Nothing is committed, merged,
published or released by this act; no global default is mutated.

## Provenance and integrity of this review

Read in full this session: the complete signed cycle-3 `review/consensus.md` (all three ✅ ACCEPT
blocks), my round-03 findings MINOR-1/NIT-1 and their dispositions, the current `IMPLEMENTATION.md`
at `ee8849c` (fix-up cycle 3 section, both AF-16 locations, evidence and metadata), `00-prompt.md`
frontmatter, both repair targets in the tree, and the live phase-8 packet body. Both review commits
were read against their own parents. This file is the only output of my invocation. Temp probes
(`/tmp/parley-claude1-r4-66202`, `/tmp/p8-*-claude1-r4.*`, `/tmp/claude1-r4-*.log`) are removed; no
background process was left running.
