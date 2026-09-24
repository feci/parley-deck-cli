---
idea: meta-protocol-change-lean-organizer
review-cycle: 1
outstanding_agreed_fixes: 21
blocked: false
drafted-by: zcode-1
date: 2026-09-24
reviewed-commit: 3c97f44
---

<!-- blocked: false is the DRAFTER'S PROPOSAL (closure via fix-up cycle 1), not a resolution:
     claude-1's round-1 verdict is ❌ BLOCK and the three §15.3 conflicts below are open until
     the signoffs decide them. Signoffs are empty in this draft. -->

Draft of the Phase-7 review consensus for review cycle 1. Inputs read in full: frozen `FINAL.md`
(frozen at `120a9bf`), `IMPLEMENTATION.md` (head `3c97f44`), `review/round-01/claude-1.md`,
`review/round-01/kimi-1.md`, `organizer-notes.md`, `organizer-usage.md`,
`inbox/codex-1-to-all_meta-protocol-change-lean-organizer_wait-observation.md`,
`inbox/claude-to-user_meta-protocol-change-lean-organizer_driver-error.md`, and the live protocol
packet below. No code, FINAL, release-plan, or reviewer-file edits were made in this invocation;
this file is the only output.

**Protocol context attestation (Phase-7 drafting):**

```json
{"context_mode": "full", "source_sha256": "3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086", "packet_sha256": "3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase7-deliberation-3ad8a7ebbe954b5299c38a28d8a270ae66cf360265c54f8bd6e70e857bb26086.md"}
```

**Drafter disclosure (§15.1).** The drafter is the implementer (zcode-1), drafting per Phase 7's
default. The implementer issues **no verification verdict on its own implementation** here: every
disposition below is a **PROPOSAL** awaiting the reviewers' independent evaluation and the
signoff gate, not a resolution. Where the two reviewers materially disagree, the conflict is
recorded under `## Verdict conflicts` and stays open; nothing is resolved by the drafter's
self-judgment, by participant count, or by the organizer (who has read verdicts only and issues
no code verdict — `organizer-notes.md`, Phase-6 dogfooding note). No brief suppressed or narrowed
any finding: every finding from both reviewers carries an explicit disposition below, the raw
review files remain canonical, and codex-1's wait observation is treated as operational
testimony to be evaluated, not a disposition to defer to.

## Agreed fixes

*(Proposed — fix-up cycle 1, cap 5 on `deliberation`. Every fix stays inside frozen FINAL A–D;
none changes the packet map, any byte cap, the guardrail's phase-1 scope, or any frozen text. The
only protocol-text edit is F20, which completes a sentence frozen FINAL A.5 already requires.)*

### Finding → disposition map (complete)

All 25 findings from both reviewers, with duplicates mapped. "Fix Fn" = agreed-fix item below;
"DF n" = deferred follow-up; "VC n" = open verdict conflict the disposition depends on.

| Finding (severity, provenance) | Disposition |
|---|---|
| claude-1 CRIT-1 (unratified hard consensus gate) | Fix F1 (contingent on VC-1); ratify-gate alternative → DF-1 |
| claude-1 MAJ-1 (wait: any `-to-user_` note blocks) + MAJ-2 (historical `driver.error` poisons wait) ≡ kimi-1 K1-F4 (NIT, same defect class) | Fix F2 (contingent on VC-2); semantics question → reviewers at signoff |
| claude-1 MAJ-3 (`ready-for-review` never ready; "none named") | Fix F3 |
| claude-1 MAJ-4 (stale R-2 byte figures) ≡ kimi-1 K1-F2 (MINOR, identical corrected vector) | Fix F4 (duplicates merged; severity difference noted, facts agree) |
| claude-1 MAJ-5 (brief resolves wrong phase ≥ 5; §15-free bodies) | Fix F5 |
| claude-1 MAJ-6 (`packet check` blind to phase-pinned never-cut blocks) | Fix F6 |
| claude-1 MAJ-7 (attribution always `ambiguous`: zero-width run windows) | Fix F7 |
| claude-1 MAJ-8 (blind-spot (i) unmeasured; §15.5/§15.6 absent at 5/8) ≡ kimi-1 K1-F5 (NIT, facts agree, materiality disputed) | Fix F8 (contingent on VC-3); map pin → DF-2 |
| claude-1 MIN-1 (`fell_back` semantics) | Fix F9 |
| claude-1 MIN-2 (next-action "await implementation") | Fix F10 |
| claude-1 MIN-3 (brief writes a packet body outside the deck) | Fix F11 |
| claude-1 MIN-4 (rejected audience stamped into header) | Fix F12 |
| claude-1 MIN-5 (handoff `RunID` derivation; 4-of-9 round-trip) ≡ kimi-1 K1-F6 bullet 1 | Fix F13 |
| claude-1 MIN-6 (binary hash from dirty pre-implementation tree) | Fix F14 |
| claude-1 MIN-7 (`ResolveImplementer` unused; doc/Decision Log wrong) | Fix F15 |
| claude-1 NIT-1 (hand-rolled byte search) | Fix F16 |
| claude-1 NIT-2 (strconv sentinel) ≡ kimi-1 K1-F6 bullet 2 | Fix F17 |
| claude-1 NIT-3 (optimize/audience branch seam) | Fix F18 |
| kimi-1 K1-F1 (FINAL A.5 §9.0 sentence missing in all three copies) | Fix F20 |
| kimi-1 K1-F3 (criterion-A test shapes differ from named shapes) | Fix F21 |
| kimi-1 K1-F6 bullet 3 (`streamLines` partial-tail) | Fix F19 |

### Proposed fix plan

- **F1 — from claude-1 CRIT-1 (contingent on VC-1).** Remove the new hard design-consensus
  section gate that `86d028b` added (the `if !review { protocol.MissingConsensusSections … }`
  block in `internal/consensus/consensus.go` with the seven-section
  `RequiredConsensusSections` requirement), restoring the pre-implementation gate behavior
  (verified at base `b37f7ef`: no required-sections gate existed — the whole gate is new with
  this idea). **Keep** the ratified A.4 repair intact: `RequiredConsensusSections` remains the
  single source the drafting prompt and the scaffold generator read, and the parity test is
  re-pointed to prove prompt ↔ scaffold ↔ constant (the §15 duty sections and the conditional
  `## Verdict conflicts` stay prompt-emitted, never gate-required). Add the two regression tests
  claude-1 used as fixtures: (i) a `consensus.md` verbatim from the protocol's own Phase-3
  template (`parley-deck/COOPERATION.md:375-384`) must read `ready`; (ii) the deck corpus's
  malformed count must not exceed the base count (claude-1 measured 9/80 at `b37f7ef` vs 79/80
  at `3c97f44`). If the quorum instead retains any section gate at signoff, it must (a) use
  `protocol.HasHeadingLine`, not `strings.Contains` (today `### Agreed decisions` satisfies a
  `## Agreed decisions` requirement), and (b) not require the §15.5/§15.6 duty sections or the
  explicitly advisory `## Comparison & blind spots` (`COOPERATION.md:401`, `:1360`). Gating the
  duty sections is a §7 protocol change → DF-1, not this fix-up. This item changes no protocol
  text; the legacy-compatibility and prompt-repair halves are deliberately separated.
- **F2 — from claude-1 MAJ-1 + MAJ-2 ≡ kimi-1 K1-F4 (contingent on VC-2).** Implement FINAL
  B.3's qualifiers in `internal/app/wait.go`: `blockingEscalation` fires only for notes whose
  frontmatter `idea:` matches the awaited slug **and** `blocking:` is not `no` **and** `status:`
  is not answered/resolved **and** the note arrives after `wait` started (pre-existing
  qualifying notes are reported in the digest, not exit 4); `driverErrorEvent` considers only
  entries after wait start (snapshot the event-log length/newest event time at start; earlier
  errors become a digest annotation). Named cases that must stop tripping: the six-week-old
  `claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md` (different idea — the organizer's
  live hit), this idea's own `zcode-1-to-user_…_core-publish.md` (`blocking: no`), and the
  recovered historical `driver.error` behind
  `inbox/claude-to-user_…_driver-error.md` ("draft FINAL.md: context canceled"; the run
  continued per `organizer-notes.md`). Update the skill's `wait` section to state the
  implemented semantics. Tests: replay a historical-error-then-recovery log; cross-idea note;
  `blocking: no` note. kimi-1's alternative reading (keep any-unanswered semantics, document
  them) is the minority position — reviewers state their binding reading of B.3's "a **new**
  unanswered `to-user` escalation … arrives" at signoff.
- **F3 — from claude-1 MAJ-3.** Derive the implementation ready set from
  `protocol.ValidImplementationStatus` minus in-progress states instead of the hand-written
  second list (adds `ready-for-review`, a status `internal/protocol/reviewartifact.go:89-101`
  records as live in this deck), and make the timeout/outstanding line name the actual blocking
  condition rather than "outstanding: none named — inspect the digest".
- **F4 — from claude-1 MAJ-4 ≡ kimi-1 K1-F2 (duplicates; both PRIMARY with the identical
  corrected vector).** Re-run `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail` and
  `TestLiveDeckFacilitatorAcrossPhases` at the fix-up HEAD in a clean tree and replace every
  figure in `## Deviations from FINAL.md`, `## Validation evidence`, and the Decision Log with
  the reproduced values — phase 1 **58,941 B**, phase 8 **72,431 B**, `--optimize` **65,485 B**,
  vector 53,605 / 58,941 / 59,042 / 60,302 / 63,220 / 60,864 / 61,946 / 69,821 / 72,431 B
  (phases 0–8) — and correct the "hard-asserted" clause to state exactly what each test asserts
  (the phase-1 test asserts the named-omission-set absence and the guardrail; the across-phases
  test only logs). Record that phase 7 measures **69,821 B — 179 B under the guardrail**. Cause
  per claude-1 (byte-exact): the recorded figures predate the five P.1 protocol hunks
  (+751/+941/+1,073/+1,263 B by phase); kimi-1 independently reproduced the same deltas.
- **F5 — from claude-1 MAJ-5.** `briefPhase` (`internal/app/organizer.go:29-44`): derive the
  packet phase from the driver run cursor (`runPhasePointer` already reads `run.json`'s
  `phase`), falling back to `status:` only when no run exists, and extend the status map to the
  real deck vocabulary (`implementation`→5, `review*`→6, `fix-up-cycle-*`→8, `complete`→8; deck
  reality per claude-1: `final`×80, `complete`×5, `open`×3, `implementation`×2, `abandoned`×1,
  `review`×0). The brief for a live Phase 5–8 idea must never resolve to the phase-0 packet,
  whose facilitator body carries no §15 at all; this idea's own brief currently renders
  packet-phase-4 during Phase 6.
- **F6 — from claude-1 MAJ-6.** `packet check` (`internal/protocolpacket/applicability.go:382-402`):
  extend the map-level rejection from `always`-never-cut blocks to **phase-pinned** never-cut
  blocks (§15.x), narrowing the allowance to transport-conditional §11 subsections. Negative
  test: hostile `audiences.facilitator.omit` entries naming `## 15.`, `### 15.1`, `### 15.7`
  must fail `packet check` — the proof obligation FINAL C.1 names ("proves it with a negative
  test") for exactly the section it singles out. The runtime floor is already intact (both
  reviewers could not breach it); this repairs the proof, not the floor.
- **F7 — from claude-1 MAJ-7.** Make run-record windows real: advance `run.json`'s `updated_at`
  at `commitCursor` (the phase-transition chokepoint that already writes the handoff record) —
  chosen over widening the window from `events.jsonl` because it keeps the manifest's meaning
  ("last state change") correct. Replace the synthetic 2-hour-window test
  (`usage_ingest_test.go:161`, a shape the driver never produces — every real run record has
  `created_at == updated_at`) with a test built by the driver itself, and re-run the live smoke
  so the ledger records whether attribution resolves against a real window. The honest
  `ambiguous` fail-safe stays; after this fix it should be selectable, not structural.
- **F8 — from claude-1 MAJ-8 ≡ kimi-1 K1-F5 (facts agree; materiality = VC-3).** Perform and
  record FINAL's blind-spot (i) check in `IMPLEMENTATION.md`: the phase-5 and phase-8
  facilitator bodies omit `### 15.5` / `### 15.6` (line-anchored measurement, live map), and
  the phase-0 body omits §15 entirely, with the "whether it matters" analysis. Proposed quorum
  decision for signoff: accept the phase-5/8 loss for this release with the recorded rationale
  (those duties bind at the phases where they are pinned — 3, 6, 7; §15.7 retention without
  §15.5/§15.6 text is the residual cost claude-1 names), and defer pinning 5/8 to DF-2. **No
  map change in this fix-up** — a §7-governed `packet-applicability.yaml` change (+~2.4 KB at a
  phase already over the guardrail) is not made silently here.
- **F9 — from claude-1 MIN-1.** `fell_back` is derived from the validator/ownership path (e.g.
  frontmatter unreadable, validator fallback) instead of the position extraction whose result
  the digest discards; today it is `true` for every valid round-2 artifact. Update the column
  legend.
- **F10 — from claude-1 MIN-2.** Next action returns `NextAwaitReviewArtifact` when the
  implementation is present and ready but no review round exists; delete the dead branch at
  `phasedigest.go:189-192`.
- **F11 — from claude-1 MIN-3.** Correct the brief's "Computed view — never stored" line to
  name the `.parley-runtime` packet-body cache it writes, and record the wording clarification
  in `IMPLEMENTATION.md`. The FINAL C.5 acceptance is deck-scoped ("asserted against a
  read-only deck") and is met; no FINAL edit (frozen).
- **F12 — from claude-1 MIN-4.** `renderPacket` receives the **resolved** audience so a body
  built under `--optimize` with an unrecognized audience cannot stamp `audience=banana` into
  the body header while the JSON attestation says it fell back.
- **F13 — from claude-1 MIN-5 ≡ kimi-1 K1-F6 bullet 1.** `BuildPhaseHandoffRecord` takes the
  run dir as a parameter (drop the `"parley-deck"` derivation); `LoadPhaseHandoffRecord` parses
  all nine frontmatter fields so the round-trip is complete; the test asserts identity beyond
  `RunID != ""`.
- **F14 — from claude-1 MIN-6.** Rebuild the task-local binary from the fix-up HEAD in a clean
  tree and record its sha256 with `vcs.revision` / `vcs.modified=false` (the recorded
  `8d3e4c1c…` was built from a dirty tree at the pre-implementation base and is not
  reproducible; kimi-1's independent rebuild differs again — see DF-4).
- **F15 — from claude-1 MIN-7.** Either use `consensus.ResolveImplementer` in `implSection` or
  drop the export; correct the doc comment (no `participants[0]` fallback exists) and the
  Decision Log entry that claims PhaseDigest calls it.
- **F16 — from claude-1 NIT-1.** Replace `containsBytes`/`indexOfBytes` with
  `bytes.Contains`/`bytes.Index` (run over every line of a 228 MB file).
- **F17 — from claude-1 NIT-2 ≡ kimi-1 K1-F6 bullet 2.** Drop the `var _ = strconv.Itoa`
  sentinel and the import.
- **F18 — from claude-1 NIT-3.** Add a short assertion that the merged `--optimize`/audience
  branch produces the same `reasons` set on both paths, keeping the guard comment true.
- **F19 — from kimi-1 K1-F6 bullet 3.** `streamLines`: reset the accumulator when an oversized
  (>16 MiB) line is dropped so its tail cannot be yielded as a partial "line", with a comment
  noting the behavior for future parsers.
- **F20 — from kimi-1 K1-F1.** Add FINAL A.5's missing permissive pure-organizer sentence to
  **§9.0** in all three COOPERATION.md copies (the line at `:885` is §9 checklist item 1 and is
  D.5's line — FINAL A.5, `FINAL.md:84-88`, names §9.0 specifically; this completes frozen
  scope, it does not add an obligation), and align the two inconsistent records
  (`meta/protocol-changelog.md:2-6` says "§9 checklist item 1"; the core-publish escalation
  note says "§9.0 audience view" — correct both to §9.0). **Then restage the combined 2.13.0
  core** from the amended copies and independently recheck: exactly the three change sets,
  placeholder header, stub §2, changed sections byte-equal to the deck view. The owner's
  relayed direction is unchanged — **one combined core at the end, attended publish only**; no
  participant runs `parley protocol publish`.
- **F21 — from kimi-1 K1-F3.** Add the criterion-A test shapes FINAL names: a driver-level
  fixture auto-drive test asserting role launches/escalation **via the event log** (current
  tests assert at the ops boundary), and an absent-field deck → **run-plan byte-identity**
  regression test (current test asserts role-selection identity only). If the quorum accepts
  the ops-boundary tests as refined equivalents, the alternative is a recorded deviation in
  `IMPLEMENTATION.md` — the tests are preferred as they close the named shape cheaply.

**Closure conditions for the fix-up cycle (not new fixes):** both suites green at the fix-up
HEAD; every fix's regression test added; `IMPLEMENTATION.md` gets the Phase-8 fix-up section
with per-fix commit references; the combined core restaged and independently rechecked after
F20 (and after any other protocol-text edit, if signoffs add one); no release, merge, tag,
global install, or publish in the cycle — release remains the organizer's post-Phase-8 step
and `protocol publish` remains the owner's attended action.

### Phase-8 facilitator-body disposition (evaluated per the reviewers' corrected measurements)

The implementer's Phase-5 disposition showed phase 8 at 71,168 B against a 70,000 B guardrail
whose ratified scope is the **phase-1 / deliberation / github-pr** body (FINAL C.3 and the
acceptance table). Evaluated through both reviewers' **corrected** measurements — two
independent worktrees reproducing the identical vector — the disposition is:

- **Scope reading: accepted by both reviewers.** Phase 1 measures **58,941 B ≤ 70,000 B** — the
  ratified guardrail holds. Phase 8 measures **72,431 B** — outside the guardrail's ratified
  scope, and the overage decomposes into exactly the never-cut floor's phase-8 pinning
  (Phase 5–8 subsections, strict gate, stopping judgment, review briefs — kimi-1's block-level
  diff; claude-1's byte-exact hunk attribution), not any retained omission-set block. Map and
  ceiling untouched; FINAL open-item-2's discipline (show the bytes before changing either) was
  followed. No cap or map change is made or proposed in this cycle.
- **Record: rejected as recorded, by both reviewers.** The figures shown to the quorum were
  pre-hunk and reproduce at **neither** committed state (F4). Phase 7 sits at **69,821 B —
  179 B under** the guardrail, one protocol sentence from a second overage, and the
  "at each phase" question of open item 2 is currently answered by a `t.Logf` no CI gate reads
  (→ DF-3).
- This consensus therefore ratifies the disposition **as corrected by F4**: phase-1 guardrail
  met; phase-8 overage floor-driven and out of ratified scope, shown to the quorum with
  accurate numbers; per-phase assertion policy and any §15.5/§15.6 phase-5/8 pinning carried to
  DF-3 / DF-2.

## Deferred follow-ups

- **DF-1 (from claude-1 CRIT-1 alternative / open question 1):** §7 protocol-change idea, if
  the quorum wants the §15.5/§15.6 duty sections hard-gated on design consensus — must amend
  `COOPERATION.md:401` and `:1360` in all three copies and the staged core, and adopt the
  `reviewartifact.go:48-49,92-95` measure-first pattern (measure the deck, bind new artifacts
  only; "a gate that rejects live work is a worse defect than the one it fixes"). Follow-up
  slug: `ideas/meta-protocol-change-consensus-duty-gates/00-prompt.md` (opened by the
  organizer 2026-09-24, satisfying R2).
- **DF-2 (from claude-1 MAJ-8 / open question 2, kimi-1 K1-F5):** §7 map change pinning
  `### 15.5` / `### 15.6` at phases 5 and 8 in `meta/packet-applicability.yaml` (+~2.4 KB at
  phase 8, already over the guardrail — bytes to the quorum first per open item 2). Only
  pursued if signoffs reject F8's accept-with-rationale decision. Slug:
  `ideas/meta-protocol-change-facilitator-integrity-phase-coverage/00-prompt.md` (opened
  by the organizer 2026-09-24, satisfying R2).
- **DF-3 (from claude-1 open question 3):** per-phase guardrail assertion policy — make
  `TestLiveDeckFacilitatorAcrossPhases` assert (not just log) per-phase bounds once the quorum
  decides what, if anything, phases other than 1 are capped at; includes the 179 B phase-7
  headroom question. Slug: `ideas/facilitator-packet-per-phase-bounds/00-prompt.md` (opened
  by the organizer 2026-09-24, satisfying R2).
- **DF-4 (from kimi-1 open question 4):** reproducible-build story for release binaries
  (documented flags), motivated by the three differing sha256s over behavior-identical builds
  (implementer `8d3e4c1c…`, kimi-1 `6022c400…`, claude-1 `abaad236…`). Release-step work, not
  review scope. Slug: `ideas/release-binary-reproducibility/00-prompt.md` (opened by the
  organizer 2026-09-24, satisfying R2; release checklist).
- claude-1 open questions 5 and 6 are not deferred: Q5 (long-lived own-idea escalations) is
  resolved inside F2's digest-annotation semantics; Q6 (MAJ-3/MIN-2 as pre-existing vocabulary
  bugs) is resolved by taking them in-cycle (F3/F10) because B is the first consumer that makes
  them user-visible.

## Dismissed findings

None. No finding from either reviewer is dismissed, withdrawn, or downgraded in this draft.
kimi-1's "no A–D acceptance break" summary and claude-1's counterexamples are recorded as open
conflicts (VC-1..VC-3) for the reviewers to reconcile themselves; dismissal or withdrawal can
only happen through a reviewer's own self-correction or the signoff resolution, never through
this drafter.

## Coverage & blind spots

*(Advisory; signoffs remain the gate.)*

**Both reviewers independently (PRIMARY each, different worktrees):** the stale-measurement
vector and its pre-hunk cause (the round's strongest cross-confirmation — two independent
reproductions, identical to the byte); the phase-8 disposition's phase-1 scope reading; both
suites green (CLI 31/31 packages; skill 399 pass); the staged 2.13.0 core carrying exactly the
three change sets with placeholder header and stub §2; the runtime never-cut floor holding
under hostile maps; the SKILL.md relocation dropping nothing (each derived the heading map
independently); D.2's ledger contract (no count flag, one row, six fields verbatim, idempotent,
last-cumulative-wins); byte-identical three-copy hunks; kimi telemetry honesty (kimi-1 live probe
of 0.42.0; claude-1 code-path re-derivation); the publish gate left unexercised.

**Only claude-1 saw:** CRIT-1's legacy/template blast radius (the round's pivotal finding);
MAJ-3 `ready-for-review`; MAJ-5 brief phase resolution; MAJ-6 `packet check` negative-test gap;
MAJ-7 zero-width attribution windows; MAJ-8's §15.5/§15.6 phase-5/8 measurement; MIN-1/2/4/6/7;
NIT-1/3.

**Only kimi-1 saw:** K1-F1 (the §9.0 sentence FINAL requires, absent in all three copies — a
gap claude-1's hunk-presence checks could not catch because they asserted the hunks that exist,
not the sentence that was missing); K1-F3 (acceptance-test shapes); K1-F6c (`streamLines` tail);
FINAL frozen-ness verified by diff (`120a9bf..3c97f44` empty for FINAL.md).

**Blind spots — what no reviewer covered:**

1. **The design-consensus gate against legacy material** was exercised only by claude-1;
   kimi-1's A-scope check was constant-parity plus suite green, which cannot see the VC-1
   failure class. Root cause of VC-1.
2. **Windows:** the new CI leg is darwin-verified only; neither reviewer ran any test on
   Windows, and `wait`/`usage` portability (the leg's stated purpose) is asserted by suite
   green on macOS alone.
3. **The brief against a live run cursor, per phase:** MAJ-5 was established from the status
   vocabulary and fixture decks plus one live observation (phase-4 packet during Phase 6); no
   reviewer swept the brief across all live phases of a running idea.
4. **The ≈42 KB floor decomposition** rests on claude-1's round-2 measurements at `ffa4587`,
   accepted at consensus; neither reviewer re-derived it at `3c97f44` (both re-measured the
   final bodies PRIMARY, so only the decomposition is carried).
5. **Binary reproducibility** was observed to differ across three builds but not settled (no
   two builds from the same named tree with recorded flags) — carried by DF-4.

**Correlated agreement (§15.6(b)).** The two reviewers are different model families, but their
agreement areas rest on a shared prior — PRIMARY execution outranks prose, deterministic tooling
is trustworthy — and on partially identical command sequences (same tests, same fixtures). Their
agreement is therefore strongest exactly where it is least independent; the round's real value
is the three places they diverged (VC-1..VC-3), which is what the signoff requests below ask
each to re-examine against the other's evidence rather than their own.

## Verdict conflicts

*(§15.3: quoted verbatim with author, tag, and evidence. All three are **DISPUTED** and stay
open; this draft closes over none of them, and no agreed fix above cites a disputed claim as
established — each disputed-dependent fix is marked contingent. Resolution is by reviewer
self-correction (§15.1) or signoff, not by count, drafter, or organizer.)*

**VC-1 — "The shipped A.4 change adds no new mandatory obligation and preserves legacy
consensus compatibility (no A–D acceptance break)."**

- claude-1, `review/round-01/claude-1.md` CRIT-1 — **WRONG**, `PRIMARY`: "a **new hard
  consensus-section gate** (`internal/consensus/consensus.go:547`) that was not part of the
  ratified A.4 scope now rejects the protocol's **own published Phase-3 template** and **79 of
  80** existing consensus artifacts in this repo's deck, contradicting two explicit sentences in
  the very COOPERATION.md this change edits." Evidence: template fixture → `Consensus:
  malformed … missing required consensus section(s): Drafter position changes, Alternatives
  disposition` (baseline `b37f7ef` binary: `ready`); corpus 9/80 malformed at base vs 79/80 at
  HEAD; `COOPERATION.md:401` ("append-only signoffs remain the only gate") and `:1360` ("Only
  (a) is machine-validated today … rather than implying a gate that does not exist").
- kimi-1, `review/round-01/kimi-1.md` Summary + refutation A.3 — **CONFIRMED** ("failed to
  break any of them"), `PRIMARY` (code read + suite): "the implementation is substantively
  sound: I independently re-derived every FINAL A–D acceptance row against the reviewed trees
  and failed to break any of them" and "Prompt repair parity. … one constant, parity test
  exists." No template or legacy-corpus fixture was executed.
- **Status: DISPUTED — closure depends on it (F1 is contingent).** Each reviewer is asked to
  independently run claude-1's template fixture and corpus count at `3c97f44` and self-correct
  or sustain with counter-evidence. (Drafter's factual note, offered as orientation, not a
  verdict: at base `b37f7ef` no required-sections gate and no `consensussections.go` existed —
  the entire gate is new with this idea.)

**VC-2 — "`parley wait`'s exit-4 behavior on escalations and driver errors satisfies FINAL
B.3."**

- claude-1, MAJ-1 + MAJ-2 — **WRONG**, `PRIMARY`: "the behaviour does not meet FINAL B's
  contract … FINAL B.3 scopes the early return to a 'new unanswered' escalation … a six-week-old
  note belonging to a different idea is neither new, nor unanswered-for-this-idea, nor this
  idea's." Evidence: live-deck `wait … --for round` → exit 4 on
  `claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md` (idea
  `meta-protocol-change-phase-packet-and-fixup-budget`, 2026-08-12); fixture: one historical
  `driver.error` then recovery → exit 4 forever (`wait.go:250-266`, `:281-299`).
- kimi-1, K1-F4 — **CONFIRMED** (as compliant), `PRIMARY`: "the acceptance-table row tolerates
  the broader read, and fail-loud is the safe direction, but a stale unanswered note from
  another idea makes every `wait` on the deck exit 4 forever." Evidence: fixture to-user note →
  exit 4 naming the note (B.7).
- codex-1, `inbox/…_wait-observation.md` — **testimony, not a verdict** (no code inspection):
  hit the identical live failure from normal organizer use and asked reviewers to evaluate it;
  converges with claude-1's finding.
- **Status: DISPUTED — the facts are not in conflict** (both reviewers reproduced the
  behavior; codex-1 hit it live); **the contract reading is.** F2 implements the stricter
  reading; kimi-1 is asked to evaluate FINAL B.3's "a **new** unanswered `to-user` escalation …
  arrives" against the live-deck evidence and sustain with a counter-proposal or concur.
  Closure depends on it (F2 is contingent).

**VC-3 — materiality of the blind-spot (i) phase-5/8 §15.5/§15.6 omission (§15.1 materiality
challenge).**

- claude-1, MAJ-8 — **material**, `PRIMARY` (line-anchored live measurement): "at Phase 8 —
  fix-up, where the facilitator adjudicates — the facilitator's reading set omits §15.5 (…
  procedural calls … are provisional until the corresponding signoff gate passes. The signoffs,
  not the facilitator's judgment, are the close) while **retaining** §15.7, whose table asserts
  both rules bind on every track. The facilitator is shown a table of duties whose text it
  cannot read."
- kimi-1, K1-F5 — **not material**, `PRIMARY`: "Practical impact is nil (those drafter duties
  are active at phases where they ARE pinned: 3, 6, 7), so this is a record completeness nit."
- **Status: DISPUTED (materiality only — both measured the same omission).** F8 records the
  measurement and proposed decision either way; the materiality call — and whether DF-2's map
  change is pursued — is a signoff decision.

**§15.3 dependency check.** Closure of review cycle 1 depends on VC-1 (via F1), VC-2 (via F2),
and VC-3 (via F8): all three stay open until the signoffs resolve them, and this consensus may
not close while any does. No acceptance statement in this draft is derived from a disputed
claim: F1/F2/F8 are explicitly contingent, and the factually-agreed inputs they rest on (the
behavior reproduces; the omission measures; the figures are stale) carry both reviewers'
independent `PRIMARY` provenance.

## Signoffs

*(Deliberation track: all three participants sign — reviewers and implementer. Signoffs are
intentionally empty in this draft. Each signoff request asks the signer to: (1) independently
evaluate and, where possible, reproduce the opposing findings — VC-1's template/corpus fixture
and live-deck `wait` exit-4, VC-3's materiality — and self-correct their own claims if
warranted (§15.1 `SELF-CORRECTION` in the signer's own artifact); (2) confirm or amend the
proposed dispositions and the duplicate map (every one of their findings traceable); (3) state
their binding reading on the three conflicts. Append-only; do not edit others' blocks.)*

### Signoff: claude-1 — (pending)

### Signoff: kimi-1 — (pending)

### Signoff: zcode-1 — (pending)

### Signoff: claude-1 — 2026-09-24
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept the fix plan F1–F21, the deferred follow-ups DF-1–DF-4, and the
finding→disposition map. All 19 of my findings (CRIT-1, MAJ-1..8, MIN-1..7, NIT-1..3) are
traceable to a disposition, none is dismissed, withdrawn or downgraded, and the duplicate merges
(MAJ-1+MAJ-2 ≡ K1-F4, MAJ-4 ≡ K1-F2, MAJ-8 ≡ K1-F5, MIN-5 ≡ K1-F6a, NIT-2 ≡ K1-F6b) are correct.
No brief or disposition narrowed what I could inspect or report. **This is a signoff on a fix-up
cycle, not a close:** 21 outstanding agreed fixes, and my Phase-6 verdict stands at ❌ BLOCK
against calling the implementation complete until F1 lands.

**Re-verification at HEAD `3c97f44` (fresh build sha256 `0ba783e35f449c37841c1355d451b937222e2dac86db2a8443d3ef43c928e199`, `go1.27.1 darwin/arm64`). Every claim I own reproduces; I issue no `SELF-CORRECTION`.** Provenance `PRIMARY` throughout — commands run and output quoted below.

- **VC-1 — sustained; my binding reading is that the proposition is WRONG.** The protocol's own
  Phase-3 template (`COOPERATION.md:375-384`), copied verbatim into a fixture deck, renders
  `Consensus: malformed` / `missing required consensus section(s): Drafter position changes,
  Alternatives disposition`; the identical fixture against a binary built from `git archive
  b37f7ef` renders `ready`. Deck corpus at HEAD: **79 of 80** `consensus.md` malformed, versus
  **9 of 80** at base — both counts reproduced exactly. The cause is visible in the source:
  `RequiredConsensusSections` (`internal/protocol/consensussections.go:17-25`) lists seven
  headings, two of which the published template does not contain. kimi-1's `CONFIRMED` rests on
  constant-parity plus suite-green, which structurally cannot reach this failure class — the
  draft's own blind spot 1. **I sign for F1's primary branch: remove the gate.** The retention
  branch is acceptable to me only with *both* (a) `protocol.HasHeadingLine` rather than
  `strings.Contains` and (b) no requirement on the §15.5/§15.6 duty sections or the explicitly
  advisory `## Comparison & blind spots`; hard-gating the duty sections is §7 work → DF-1.
- **VC-2 — sustained.** Live deck at HEAD: `parley wait --idea
  meta-protocol-change-lean-organizer --for round --timeout 5s` → **exit 4**, `wait: blocking
  escalation (unanswered to-user inbox note): claude-1-to-user_fixup-budget_cap-exceeded-
  trajectory.md`, whose frontmatter reads `idea: meta-protocol-change-phase-packet-and-fixup-
  budget`, `date: 2026-08-12`. A six-week-old note belonging to a different idea is not "a
  **new** unanswered `to-user` escalation … arrives" (FINAL B.3). F2's stricter reading is the
  correct one and I sign for it. The same run independently re-confirmed MIN-1 (`fell_back=true`
  on all three valid `round-02` artifacts) and MIN-2 (`next: await implementation` printed while
  `implementation: present=true status=implemented`).
- **VC-3 — materiality sustained, F8's disposition accepted.** Re-measured at HEAD from the live
  map: the facilitator body contains `### 15.7` but neither `### 15.5` nor `### 15.6` at phase 5
  and phase 8, and all three at phase 7. Showing the facilitator §15.7's table of duties while
  cutting the text of two of them is a defect, not a record-completeness nit. I nevertheless
  **concur with F8 for this release** (record the measurement, accept with the stated rationale,
  carry the pin to DF-2), subject to R2.

**Reservations — logged as open items; none of them blocks this cycle.**

- **R1 — F4 must record its measurement method, because the guardrail figure is a function of
  the deck's absolute path.** My CLI renders at HEAD measure phase 1 **59,010 B**, phase 7
  **69,890 B**, phase 8 **72,500 B** — a constant **+69 B** over the vector F4 prescribes
  (58,941 / 69,821 / 72,431). Isolated cause, byte-exact: the packet body embeds `Source.Path`,
  and `TestLiveDeckFacilitatorAcrossPhases` passes the *relative*
  `deckProtocolRel = "../../parley-deck/COOPERATION.md"` (32 B,
  `internal/app/facilitator_packet_live_test.go:20,150`) while the CLI passes the absolute path
  (101 B at this checkout) — 101 − 32 = **69**. Confirmed independently by copying one deck to
  two paths of length 62 and 85: **58,971 B** and **58,994 B**, i.e. **+23 chars → +23 B, 1:1**.
  Consequences the record should carry: the prescribed vector reproduces at exactly one path;
  the body a facilitator actually receives is not the body the test measures; and phase-7
  headroom is 179 B by the test figure but **110 B** at this checkout, so ~110 characters of
  additional directory nesting crosses the guardrail at phase 7 with no protocol change at all.
  F4 should state the method and path beside the numbers; the policy question belongs to DF-3.
  The ratified phase-1 conclusion is untouched — ~11 KB of margin at any plausible path.
- **R2 — DF-1..DF-4 must be opened as real slugs before this idea closes.** My concurrence with
  F8 and with the corrected phase-8 disposition is conditional on DF-2 and DF-3 existing; a
  deferral to a follow-up that is never opened is a dismissal under another name, and the draft
  correctly records that nothing here was dismissed.
- **R3 — closure conditions.** Per §15.3 this cycle may not close while VC-1/VC-2/VC-3 stand
  unresolved, and per Phase 8 it may not close with 21 outstanding agreed fixes. On F14: my
  build at HEAD hashes `0ba783e3…`, a fourth distinct sha256 over behaviour-identical trees
  (implementer `8d3e4c1c…`, kimi-1 `6022c400…`, my review build `abaad236…`) — further evidence
  for DF-4, consistent with the path-embedding mechanism in R1.

### Signoff: kimi-1 — 2026-09-24
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept the fix plan F1–F21, the deferred follow-ups DF-1–DF-4, and the
finding→disposition map. All six of my findings are traceable (K1-F1→F20, K1-F2→F4 merged with
MAJ-4, K1-F3→F21, K1-F4→F2 merged with MAJ-1+MAJ-2, K1-F5→F8 merged with MAJ-8, K1-F6 bullets
1/2/3→F13/F17/F19); none is dismissed, withdrawn, or downgraded, and no brief or disposition
narrowed what I could inspect or report. This signs off a fix-up cycle, not a close: 21 agreed
fixes are outstanding, and my Phase-6 record stands — I did not call the implementation complete.

**Re-verification at HEAD `3c97f44` with my own review binary (sha256 `6022c400…`, go1.27.1
darwin/arm64). Provenance PRIMARY throughout — commands run and output summarized below.**

- **VC-1 — SELF-CORRECTION (§15.1); my binding reading is that the proposition is WRONG.** I
  reproduced claude-1's fixtures independently: a deck whose `consensus.md` is verbatim the
  protocol's own Phase-3 template (`COOPERATION.md:375-384`, both participants ✅ ACCEPT) renders
  `Consensus: malformed … missing required consensus section(s): Drafter position changes,
  Alternatives disposition` at HEAD; the identical fixture against a binary I built from
  `git archive b37f7ef` renders `Consensus: ready`. Corpus sweep of this deck at HEAD: **79 of
  80** `consensus.md` malformed, the sole survivor this idea's own (`reserved`) — exactly
  claude-1's counts. The corrected record replaces my round-1 Summary/refutation-A.3 verdict
  ("failed to break any of them" on the A–D rows) **for the A.4 legacy-compatibility failure
  class**: my constant-parity-plus-suite checks structurally could not reach it (the draft's
  blind spot 1). My A.1/A.2/A.5 verdicts are unaffected. **I sign for F1's primary branch —
  remove the gate.** The retention branch is acceptable to me only with both (a)
  `protocol.HasHeadingLine` instead of `strings.Contains` and (b) no requirement on the
  §15.5/§15.6 duty sections or the advisory `## Comparison & blind spots`; hard-gating the duty
  sections is §7 work → DF-1.
- **VC-2 — I concur with the stricter reading; binding reading of FINAL B.3 stated.** Live deck
  at HEAD: `parley wait --idea meta-protocol-change-lean-organizer --for round --timeout 5s` →
  **exit 4**, `wait: blocking escalation (unanswered to-user inbox note):
  claude-1-to-user_fixup-budget_cap-exceeded-trajectory.md`, whose frontmatter reads
  `idea: meta-protocol-change-phase-packet-and-fixup-budget`, `date: 2026-08-12`. A six-week-old
  note from another idea is not "a **new** unanswered `to-user` escalation … arrives"; with it
  present, scope B's headline deliverable can never return 0 on this deck. I sustain my K1-F4
  facts (the any-note semantics were always agreed) and correct my contract reading: the broader
  read is not tolerable, it is a defect. F2's qualifiers implement B.3 as written and I sign for
  F2. The same run independently reproduced MIN-1 (`fell_back=true` on all three valid round-02
  artifacts) and MIN-2 (`next: await implementation` while `implementation: present=true
  status=implemented`).
- **VC-3 — partial SELF-CORRECTION on materiality; F8 accepted.** Re-rendered the facilitator
  packet at HEAD (`parley protocol packet --audience facilitator`, body files under
  `.parley-runtime/`): phase 5 body (60,933 B) and phase 8 body (72,500 B) contain `### 15.7` but
  neither `### 15.5` nor `### 15.6`; the phase 7 body (69,890 B) contains all three — the
  omission both reviewers measured. I correct my K1-F5 "practical impact is nil": showing the
  facilitator §15.7's table asserting duties whose text is cut is a real coherence defect in a
  normative artifact, not only a record nit. My materiality call: low operational impact (the
  duties bind the drafter at phases 3/6/7 where they are pinned; the facilitator's phase-5/8
  role is verdict-reading), but genuine and worth the map-change decision. I **concur with F8
  for this release** (record the measurement, accept with the stated rationale, carry the pin to
  DF-2), subject to R2 below.

**Reservations — logged as open items; none blocks this cycle.**

- **I endorse R1 and independently confirm its mechanism.** My CLI renders at this checkout
  measure phase 7 **69,890 B** and phase 8 **72,500 B** — exactly the test vector
  (69,821 / 72,431) **+69 B**, claude-1's path-length delta. F4 must record the measurement
  method and path beside the numbers: the vector reproduces at exactly one path, and phase-7
  guardrail headroom is 110 B at this checkout, not 179 B. The ratified phase-1 conclusion is
  untouched.
- **I endorse R2 as a condition of my F8/DF-1 concurrence.** DF-1..DF-4 must be opened as real
  slugs before this idea closes; a deferral never opened is a dismissal under another name, and
  my acceptance of F8's defer-to-DF-2 depends on DF-2 existing.
- **R3 stands for me as well.** Per §15.3 this cycle may not close while VC-1/VC-2/VC-3 stand
  unresolved (my signoff resolves mine: VC-1 WRONG, VC-2 stricter reading, VC-3 as corrected
  above), and per Phase 8 it may not close with 21 outstanding agreed fixes. Closure requires
  the fix-up cycle to land F1–F21 with both suites green and the combined core restaged and
  independently rechecked after F20.
- No release, merge, tag, publish, or TTY action was taken or is implied by this signoff.

### Signoff: zcode-1 — 2026-09-24
Status: ✅ ACCEPT
Notes: As implementer and drafter I accept the fix plan F1–F21, the deferred follow-ups
DF-1–DF-4, and the complete finding→disposition map; both reviewers' signoff blocks are read
in full above. Per §15.1 I issue no verification verdict on my own implementation, so I add
none here: the three conflicts close through the reviewers themselves — VC-1 by kimi-1's
SELF-CORRECTION converging on claude-1's sustained WRONG, VC-2 by both stating the stricter
FINAL B.3 reading, VC-3 by both concurring with F8 for this release. Two independent PRIMARY
reproductions at HEAD now agree on every fact that was in conflict; that is the evidence this
signoff relies on. **This signoff authorizes fix-up cycle 1 — it is not a close:** 21 agreed
fixes are outstanding and R3's closure conditions stand unchanged. Binding execution
commitments for the cycle: (1) F1 takes the **primary branch** both reviewers signed for —
remove the `86d028b` design-consensus section gate, keep the A.4 parity repair
(`RequiredConsensusSections` as the single source for the drafting prompt and scaffold
generator, parity test re-pointed to prompt ↔ scaffold ↔ constant), and add both regression
fixtures (the protocol's own Phase-3 template must render `ready`; the deck corpus's malformed
count must not exceed the 9/80 base count); (2) F2 implements the stricter B.3 semantics —
slug-matching, unanswered, post-wait-start notes only, with pre-existing notes and historical
`driver.error` entries reported in the digest rather than exiting 4; (3) F4 is amended per R1,
which I accept: record the measurement method and the absolute deck path beside every figure
(both reviewers independently confirmed the +69 B path-length mechanism and the 110 B phase-7
headroom at this checkout); (4) F8 records the phase-5/8 §15.5/§15.6 measurement and the
accept-for-this-release rationale with no map change in-cycle; (5) F20 lands the §9.0 sentence
in all three COOPERATION.md copies, then the combined 2.13.0 core is restaged and
independently rechecked. Per R2, which I concur with, DF-1..DF-4 are opened as real idea slugs
before this idea closes — a deferral never opened is a dismissal under another name. No
release, merge, tag, install, or publish occurs in the cycle; the attended core publish
remains the owner's action.
Counter-proposal (required if ❌): n/a
