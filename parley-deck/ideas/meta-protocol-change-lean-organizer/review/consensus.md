---
idea: meta-protocol-change-lean-organizer
review-cycle: 2
outstanding_agreed_fixes: 10
blocked: false
drafted-by: zcode-1
date: 2026-09-24
reviewed-commit: 118b245
---

<!-- blocked: false is the DRAFTER'S PROPOSAL, not a resolution. No reviewer files a
     blocker this round: claude-1 withdrew its round-1 ❌ BLOCK ("My round-01 ❌ BLOCK is
     withdrawn; my position for the Phase-7 consensus is 🟡 … on a fix-up cycle, not on
     closure") and kimi-1 files 4 NITs with a ready-to-close position. But the two
     reviewers DISAGREE about closure (VC-4 below), and that disagreement stays open
     until the signoffs decide it — never by count, drafter, or organizer. Signoffs are
     empty in this draft. -->

Draft of the Phase-7 review consensus for review cycle 2 (after fix-up cycle 1; reviewed
commit `118b245`). The previous signed consensus is archived **VERBATIM** at
`review/consensus-cycle-01.md` (44,537 B, `cmp`-silent against the file it replaces,
including all three signoff blocks and reservations R1–R3); this fresh draft reuses no
stale signoff. Inputs read in full: frozen `FINAL.md` (frozen at `120a9bf`; `git diff
120a9bf..118b245 -- …/FINAL.md` empty — re-verified this session), `IMPLEMENTATION.md`
(at head `118b245`), the entire signed cycle-1 `review/consensus.md`, both complete
round-02 review files (`review/round-02/claude-1.md`, `review/round-02/kimi-1.md`),
`organizer-notes.md` (post-correction state), `organizer-usage.md`,
`inbox/codex-1-to-all_meta-protocol-change-lean-organizer_wait-observation-02.md`,
`inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_core-publish.md`, the staged
core `~/.parley/staging/COOPERATION-2.13.0.md` (hash-verified below), and the live
protocol packet attested below. No code, FINAL, release-plan, reviewer-file,
organizer-note, or inbox edits were made in this invocation; this file (plus the
verbatim archive of its predecessor) is the only output. Nothing is committed, tagged,
installed, or published, and `parley protocol publish` is not run.

**Protocol context attestation (Phase-7 drafting):**

```json
{"context_mode": "full", "source_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "packet_sha256": "8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7", "fallback_reason": null, "body_path": "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/lean-organizer/.parley-runtime/protocol-packets/full-phase7-deliberation-8ce83cde0b4f1c3019229869433df0c32451fdbaa13a186952fdf9805c3a9db7.md"}
```

Drafter cross-check (PRIMARY): `shasum -a 256` over the packet body file and over the
deck `parley-deck/COOPERATION.md` both return `8ce83cde…a9db7` (109,928 B) — the
attested packet IS the live authority at the reviewed tree, the same hash both round-02
reviews attest against their own trees.

**Drafter disclosure (§15.1).** The drafter is the implementer (zcode-1), drafting per
Phase 7's default. The implementer issues **no verification verdict on its own
implementation** here: every disposition below is a **PROPOSAL** awaiting the reviewers'
independent evaluation and the signoff gate, not a resolution. Where the two reviewers
materially disagree — this round: readiness for closure (VC-4), the severity of the
`--json` stdout trailer (VC-5), and the materiality of the two next-action states
(VC-6) — the conflict is recorded under `## Verdict conflicts` and stays open; nothing
is resolved by the drafter's self-judgment, by participant count, or by the organizer
(who has read verdicts only and issues no code verdict). No brief suppressed or
narrowed any finding: all 11 claude-1 findings and all 4 active kimi-1 findings carry
an explicit disposition below, the raw review files remain canonical, and codex-1's
wait-observation-02 is treated as operational testimony — its open question (which
stream carries the `--json` trailer) is now answered by both reviewers independently
(see R2-MAJ-2 ≡ K2-F3) and the answer is recorded here rather than deferred to.

**Verifier provenance of this draft.** Code-level claims below are cited to their
reviewer owners (both reviewers tag their verdicts PRIMARY with commands quoted in
their files). The drafter independently re-verified exactly the record-replacement
facts this consensus depends on (tagged DRAFTER-PRIMARY inline): the staged-core
hash/size/mtime and the published-core directory (G6), the organizer-note newline
correction (R2-NIT-3), FINAL frozenness, HEAD `118b245`, and the packet/authority hash
equality above.

## Agreed fixes

*(Proposed — fix-up cycle 2 (cap 5 on `deliberation`; this is cycle 2). Every fix stays
inside frozen FINAL A–D or is record hygiene on this idea's own artifacts; none changes
the packet map, any byte cap, the guardrail's phase-1 scope, or any frozen protocol
text. No protocol-text edit is proposed in this cycle, so no core restage is expected;
if signoffs add one, the restage, its independent recheck, and the core-publish note
update happen together in the same step (the note must always match the current staged
file). Two items are contingent on open conflicts (G2 → VC-5, G4 → VC-6): they are
proposed as Agreed because at least one reviewer filed them at fix-worthy severity and
the facts carry both reviewers' PRIMARY provenance; a signer who sustains the lesser
position amends or strikes them at signoff — that is the conflict closing, not this
draft.)*

### Finding → disposition map (complete)

All 15 findings from both reviewers (11 claude-1, 4 active kimi-1), duplicates mapped,
severity and provenance kept. "Fix Gn" = agreed-fix item below; "VC n" = open verdict
conflict the disposition depends on.

| Finding (severity, provenance) | Disposition |
|---|---|
| claude-1 R2-MAJ-1 (floor not recorded; omission sum mislabelled `floor` and undercounted; unmet FINAL C.3/R-2 element) [MAJOR, PRIMARY — direct probe of `byHeading`/`Parse` block splitting; `## 12.`=273 B, `## 13.`=522 B vs 17 subsection blocks/11,903 B excluded; whole-section total 23,223 B] | Fix G1 |
| claude-1 R2-MAJ-2 (`parley wait --json` emits non-JSON trailer on stdout; exit-0/exit-3 unparseable; stderr = 0 bytes — settles codex-1's observation-02 question) [MAJOR, PRIMARY — streams isolated, `json.load` fails on exit-3 capture; exit-4 clean] ≡ kimi-1 K2-F3 (same defect) [NIT, PRIMARY — live envelope + trailer verified; "present since `86d028b`, not a fix-up regression; 1.49.0 will freeze it"] | Fix G2 (contingent on VC-5 — severity/readiness disputed, facts agreed); duplicates merged; severity conflict recorded as VC-5 |
| claude-1 R2-MAJ-3 (tri-platform CI leg `go test ./...` with no `-timeout` vs `internal/trajectory` 589.7 s of Go's 600 s default; FINAL's named command fails without the flag) [MAJOR, PRIMARY — durations measured, own no-flag run failed; hosted-runner comparison explicitly an inference] | Fix G3; package-duration investigation split out as advisory (Deferred follow-ups, no slug) |
| claude-1 R2-MIN-1 (next-action `await implementation` wrong in two states: (a) fix-up published + prior review round exists — this deck's live state; (b) rounds complete, no `consensus.md` — `NextAwaitConsensus` unreachable) [MINOR, PRIMARY — both states reproduced on fixtures + live digest/brief] | Fix G4 (contingent on VC-6); re-file of round-01 MIN-2 residue per claude-1's own signoff-scope admission |
| claude-1 R2-MIN-2 (F2 fails open: `to-user` note with unparseable or `idea:`-less frontmatter silently ignored — neither blocking nor annotated) [MINOR, PRIMARY — no-frontmatter note touched mid-wait → exit 3, no annotation] | Fix G5 |
| claude-1 R2-MIN-3 (core-publish escalation note describes the pre-restage file: 109,507 B / verified 2026-09-23 / no sha256; actual 109,772 B / restaged 2026-09-24 by F20) [MINOR, PRIMARY — note vs file facts] | Fix G6 (implementer-owned note; replacement facts DRAFTER-PRIMARY-verified this session) |
| claude-1 R2-MIN-4 (`## Deviations from FINAL.md` says "None in scope or mechanisms" but F1 removed A.4's gate-side parity mechanism by signed fix) [MINOR, PRIMARY — FINAL A.4 text vs post-F1 tree] | Fix G7 |
| claude-1 R2-NIT-1 (`outstandingAgents` has no parallel "consensus.md not filed" for a nil consensus section) [NIT — code read, same class as R2-MIN-1(b)] | Fix G4 |
| claude-1 R2-NIT-2 (F2 arrival signal is mtime; driver rewrites `claude-to-user` escalation notes in place — second error replaces first unanswered, refreshes mtime) [NIT, PRIMARY — `git diff 3c97f44..64a622c` shows the in-place replacement] | Fix G8 (comment + record; no behavior change — fails loud, canonical record survived in organizer-notes.md) |
| claude-1 R2-NIT-3 (organizer-notes.md ended with literal `\n` escapes) [NIT — file read] | **RESOLVED BY OWNER — no fix item.** The organizer corrected its own note before this draft (its "Review round 2" section records the correction). DRAFTER-PRIMARY verification this session: zero literal `\n` occurrences remain in `organizer-notes.md` and the file ends with a real newline byte. The note is not edited by any participant; recorded as resolved here. |
| claude-1 R2-NIT-4 (recorded §15-region figure 8,041 B does not reproduce; both measure 8,056 B; byte-equality itself holds) [NIT, PRIMARY — extraction measured] ≡ kimi-1 K2-F1 (same) [NIT, PRIMARY — `sed -n '/^## 15\./,$p' \| wc -c` = 8,056 B] | Fix G9 (duplicates merged, independent identical measurements) |
| kimi-1 K2-F2 (preflight facilitator-conflict gate fails open when `ReadWorkspaceStatus` errors: no-`COOPERATION.md` tree prints "Ready: no pending gates" while the gate never ran) [NIT, PRIMARY — own fixture reproduced; filer's leaning: fail closed] | Fix G10 |
| kimi-1 K2-F4 (F10 residual record names the wrong trigger — the `review/round-02/` DIRECTORY's creation, not the first review file) [NIT, PRIMARY — directory-present-empty vs absent states isolated] | Fix G4 (record half); behavior-half disagreement with R2-MIN-1 recorded as VC-6 |
| kimi-1 K2-F5 | Not a finding: folded into K2-F4 **by its own filer** before this consensus to keep IDs stable (recorded in `review/round-02/kimi-1.md`); not a group dismissal. |

**Organizer-owned corrections: none remain.** The only organizer-owned finding this
round (R2-NIT-3) is resolved by the organizer's own correction (above), so no inbox
correction request is required. codex-1's wait-observation-02 is accurate testimony
(it correctly hedged that the capture alone does not prove the stream); its question is
answered in this record (R2-MAJ-2 row) and in both reviewers' files — no correction to
that note is needed. The one remaining artifact with stale facts is
implementer-owned (R2-MIN-3 → G6).

### Proposed fix plan

- **G1 — from claude-1 R2-MAJ-1 (MAJOR).** In
  `TestLiveDeckFacilitatorPacketNamedSetsAndGuardrail`
  (`internal/app/facilitator_packet_live_test.go:96-119`): record the facilitator body
  against the **measured floor** FINAL C.3/open-item-14 names — compute the floor as
  the **retained whole-section total** (≈ 42 KB with §2, recomputed at the fix-up-2
  HEAD) in the same test run that records the guardrail and the `--optimize` baseline,
  and log it as `floor`. Stop calling an omission sum a floor: the variable and the
  comment at `:97` ("the whole-section floor is the sum of the source bytes of the
  named omission set" — a floor is what the body *retains*) and the `:107` §2
  approximation note are corrected. If the omission-set total is kept beside it, it is
  recomputed whole-section-correct (`byHeading[locator]` takes one `Parse` block and
  `Parse` splits `### N.M` subsections into their own blocks — today the sum silently
  excludes 17 subsection blocks / 11,903 B of `## 12.`/`## 13.` content the map does
  omit; claude-1's whole-section source total for the seven top-level omitted sections
  is 23,223 B) and logged under its own label. This takes FINAL as frozen (claude-1's
  open question 1, first branch): the recorded companion number becomes the floor
  FINAL requires — no acceptance-wording amendment, no deviation. The **gating** half
  of C.3 (named-omission-set absence, line-anchored) and the guardrail both hold per
  both reviewers and are untouched.
- **G2 — from claude-1 R2-MAJ-2 (MAJOR) ≡ kimi-1 K2-F3 (NIT) (contingent on VC-5).**
  Make `parley wait --json` stdout machine-parseable on **all** exit paths
  (`internal/app/wait.go:173-179`: `printWaitDigest` emits the JSON envelope to `out`
  and the very next line writes the human trailer to `out` too; only exit-4 routes its
  reason to stderr — both reviewers measured, stderr is 0 bytes on the exit-0/exit-3
  paths). **Primary shape** (the intersection of both reviewers' suggested fixes, and
  the smallest change): in `--json` mode route the terminal status line
  (`wait: boundary reached (…)` / `wait: timeout after …`) to **stderr**; stdout
  carries the `{notes?, digest}` envelope and nothing else; non-`--json` human output
  is unchanged. **Signoff-amendable alternative** (claude-1's open question 3): fold
  the terminal status into the envelope as a `result` field and print nothing else on
  stdout — signers pick one; both satisfy FINAL B.3's "exit codes an organizer can
  branch on". Add a test that decodes `--json` stdout through `encoding/json` on all
  four exit paths (0, 3, 4, 1). Update the skill's `wait` section to state the chosen
  stream contract and the `{notes, digest}` envelope keys (kimi-1's ask), and record
  in `IMPLEMENTATION.md` that the round-1→fix-up-1 envelope change (bare `PhaseDigest`
  → `{notes?, digest}`) rides the same unreleased 1.49.0 freeze. Grounds: codex-1 hit
  the failure in ordinary use this round (`jq: parse error`, observation-02);
  claude-1 isolated the streams; kimi-1 verified the shape live and itself recommends
  the fix "before release" — the open question between them is severity/close-timing
  (VC-5), which the signoffs settle.
- **G3 — from claude-1 R2-MAJ-3 (MAJOR).** `.github/workflows/tests.yml` (new in this
  idea): add an explicit `-timeout` to the tri-platform Test step —
  `go test ./... -count=1 -timeout 45m` (≥ the 2400 s every green record in this idea
  already uses) — and align the FINAL-named validation command recorded in
  `IMPLEMENTATION.md` (`go test ./...` → `go test ./... -count=1 -timeout 2400s`).
  Motivation recorded with the change: `internal/trajectory` takes **589.7 s** and
  `internal/app` **521.7 s** of Go's 600 s per-package default on claude-1's machine
  (~2% margin), and claude-1's run of the plan's own no-flag command **failed**
  (`panic: test timed out after 10m0s`, both packages); kimi-1's round-1 pass without
  the flag shows the outcome is load-dependent. FINAL risk 8 ships Windows assets
  "honestly" on this leg; a leg configured to time out does not deliver that. The
  separate question — why those packages take nine minutes — is advisory follow-up
  (see Deferred follow-ups), per claude-1's own open question 2, which leans in-cycle
  for the flag exactly as proposed here. Signoffs may move the whole item to a DF; the
  drafter proposes in-cycle because it is a two-token change on a file this idea
  ships.
- **G4 — from claude-1 R2-MIN-1 (MINOR) + R2-NIT-1 (NIT) + the record half of kimi-1
  K2-F4 (NIT) (contingent on VC-6).** Two next-action repairs in
  `internal/driver/phasedigest.go:186-199` plus one line in `internal/app/wait.go`:
  (a) return `NextAwaitReviewArtifact` whenever the implementation is present and
  ready and the latest review round is complete but the implementation is newer than
  it (the fix-up-published state — **this deck's live state right now**, where the
  digest and `organizer brief` say `await implementation` while awaiting re-review);
  (b) return `NextAwaitConsensus` when rounds are complete and no `consensus.md`
  exists — `BuildPhaseDigest` (`phasedigest.go:134`) populates `d.Consensus` only when
  `consensus.Status` succeeds, so a nil consensus section is a real state, not an
  absent one, and a Phase-2/3 deck is currently told to skip Phases 3–4;
  (c) `outstandingAgents` (`wait.go:266-286`) gains `"consensus.md not filed"` when
  the scope is consensus and the section is nil, instead of the generic
  "awaited condition not yet reached" (R2-NIT-1). Update the F10 residual record in
  `IMPLEMENTATION.md` to name the true trigger — the `review/round-02/` **directory**'s
  creation, not the first review file (K2-F4; with the directory present but empty the
  enumeration already reads `await review artifact`) — and, after this fix, the
  corrected enumeration. Scope basis: the quorum took the same class in-cycle in
  round 1 (MAJ-3/MIN-2 — pre-existing digest-vocabulary gaps B first exposed);
  claude-1's open question 4 proposes the same consistency here. kimi-1's position
  (behavior acceptable in both states; record-only) is the minority position recorded
  as VC-6 — signoffs decide.
- **G5 — from claude-1 R2-MIN-2 (MINOR).** `qualifyingEscalationNotes`
  (`internal/app/wait.go:297-305`): when a `*-to-user_*.md` note cannot be parsed or
  carries no `idea:` frontmatter, emit a `note:` annotation naming the file and why it
  was not evaluated (and expose it in the `--json` `notes` list). The filer's stated
  sufficient remedy is annotating — it removes the silence; blocking on it would be
  defensible but is not required. This closes the fail-open the filer reproduced (a
  no-frontmatter `URGENT` note touched mid-wait is neither blocking nor annotated
  today, silently dropped — at `3c97f44` the same note exited 4) and aligns B with its
  own stated discipline ("fail-closed `unparsed`, never a guess") and the digest's
  `unparsed`-row practice of surfacing rather than dropping.
- **G6 — from claude-1 R2-MIN-3 (MINOR).** Update the implementer-owned
  `parley-deck/inbox/zcode-1-to-user_meta-protocol-change-lean-organizer_core-publish.md`
  to match the current staged file — the note is the artifact the owner acts on for
  the one irreversible, write-once action in this release, and its integrity statement
  must match the file it points at. Exact replacement facts, **independently
  re-verified by the drafter this session (DRAFTER-PRIMARY)**:
  `shasum -a 256 ~/.parley/staging/COOPERATION-2.13.0.md` →
  `fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f`;
  size **109,772 B** (not 109,507 B — F20 restaged the file on 2026-09-24, mtime
  `Sep 24 03:35:27 2026`); composition verification date **2026-09-24** (not
  2026-09-23); add the staged sha256 so the owner can verify the reviewed bytes
  before the attended publish. `~/.parley/protocol/core/` still holds only `2.10.0` —
  **no publish has occurred, and none is performed in this cycle**; the exact attended
  command in the note is already correct and unchanged. The note must match the
  current staged file at whatever commit the cycle lands (see the restage rule in the
  section preamble).
- **G7 — from claude-1 R2-MIN-4 (MINOR).** Add one bullet to `## Deviations from
  FINAL.md`: A.4's gate-side parity ("prompt **and gate** read the same value") was
  removed by unanimously signed fix F1 — no gate reads
  `RequiredConsensusSections`/`RequiredFinalSections` any more
  (`MissingConsensusSections` is gone; `ValidateFinal` checks status/slug/scaffold
  only) — and the acceptance-table row is met by prompt ↔ scaffold ↔ constant parity;
  hard-gating the duty sections remains DF-1's §7 question. The F1 fix-up entry
  already explains the change; this puts it where Phase 5 requires deviations to be
  logged ("not silently absorbed").
- **G8 — from claude-1 R2-NIT-2 (NIT).** A comment at the mtime comparison in
  `internal/app/wait.go` recording that arrival is approximated by file mtime, that an
  in-place rewrite or bare `touch` of a pre-existing escalation therefore looks "new"
  to `wait` (fails loud — the safe direction), and why that is acceptable. Record in
  `IMPLEMENTATION.md` the driver behavior behind it (`internal/driver/loop.go:353-354`
  builds `claude-to-user_<slug>_<topic>.md` and writes it with an unconditional
  `os.WriteFile` — a second `driver.error` for the same idea replaces the first,
  still-unanswered escalation; that happened this run, and the canonical record
  survived in `organizer-notes.md:24`). **No behavior change** — the filer judges the
  current shape safe; the `from: claude` author mislabel stays FINAL-recorded
  testimony, out of scope.
- **G9 — from claude-1 R2-NIT-4 ≡ kimi-1 K2-F1 (NIT; both independently measured
  8,056 B).** Correct the F20 recheck figure in `IMPLEMENTATION.md` — "§15 region
  (8,041 B)" → **8,056 B** — and state the extraction boundary beside the number
  (`sed -n '/^## 15\. Verification integrity/,$p' | wc -c`, heading line to EOF; the
  same convention in all three copies). The load-bearing byte-equality claim already
  holds (both reviewers' per-region diffs are empty; all three §15 regions share one
  sha256); only the recorded figure and its unstated convention were wrong — the same
  record discipline round 1's central finding (F4/R1) established, applied to this
  one figure.
- **G10 — from kimi-1 K2-F2 (NIT).** `facilitatorConflictGates`
  (`internal/app/preflight.go:319-323`) returns nil on a `ReadWorkspaceStatus` error,
  so a tree with `agents.toml` + `meta/version.json` but no `COOPERATION.md` prints
  "Ready: no pending gates" while the conflict gate silently never ran (filer's
  fixture: `facilitator: claude-1` ∈ `participants:` → exit 0 there, exit 3 in a
  complete workspace). Fix: treat a workspace-status read failure as a hard preflight
  error (exit 1) — the filer's own leaning ("fail closed"), consistent with A's
  fail-closed discipline and with the status-concealment class this idea exists to
  fix; the enforced half of A (driver-side role ineligibility) is unaffected, and real
  decks always carry `COOPERATION.md`, so this is a degenerate-path repair of code
  this idea added — not new scope. Negative test: the no-`COOPERATION.md` fixture
  exits 1 naming the read failure (today exit 0 "Ready"). Signoffs that prefer the
  current behavior must instead require the skill to state it (the filer's
  alternative) — one or the other, not silence.

**Closure conditions for fix-up cycle 2 (not new fixes):** both suites green at the
fix-up-2 HEAD run with the G3-named command (explicit `-timeout`); every fix's
regression test added; `IMPLEMENTATION.md` gets the Phase-8 fix-up-cycle-2 section
with per-fix commit references and top-level frontmatter bumped; the core-publish
note's facts match the staged file at that HEAD (G6 rule); no release, merge, tag,
global install, or publish in the cycle — release remains the organizer's post-Phase-8
step and `parley protocol publish` remains the owner's attended action.

**Cycle mechanics (claude-1's open question 5).** This consensus lists 10 proposed
Agreed fixes, so it is not a closing consensus under any position in VC-4: Phase 8
closes the idea only on a Phase-7 consensus listing **zero** Agreed fixes. If the
signoffs sustain this plan, fix-up cycle 2 lands G1–G10, review round 3 verifies it
(fix-verification scope plus whatever fresh scope reviewers choose), and consensus
cycle 3 aims to be the zero-fix record. If the signoffs amend the plan downward
(kimi-1's position), the amended set is what cycle 2 lands and the struck items'
dispositions move to `## Deferred follow-ups`/`## Dismissed findings` **by the
signoffs' own written decision**, not by this draft. The fix-up budget (cap 5,
deliberation) has cycles 1–2 of 5 in use; trajectory per §4 stopping judgment is the
converging shape (round 1: 25 findings incl. 1 CRITICAL; round 2: 15 findings, 0
CRITICAL, 3 MAJOR, all smaller and confined to recording, output formatting, CI
configuration, and degenerate paths).

## Deferred follow-ups

- **DF-1 `meta-protocol-change-consensus-duty-gates`** — unchanged; remains an
  inactive `status: candidate` record (LE-10), no quorum staffed, not launched. G7's
  deviation bullet points here.
- **DF-2 `meta-protocol-change-facilitator-integrity-phase-coverage`** — unchanged;
  inactive candidate. The phase-5/8 §15.5/§15.6 pin question stays here; both
  reviewers re-concur with the F8 accept-for-this-release disposition this round
  (claude-1 "(ii) … Concur for this release", kimi-1 "(ii) … concur", with kimi-1's
  partial self-correction on materiality standing).
- **DF-3 `facilitator-packet-per-phase-bounds`** — unchanged; inactive candidate.
  Recorded update for whichever round reads it next (both reviewers' round-02
  measurements, PRIMARY each): phase 7 now sits **86 B over** the 70,000 B figure at
  the test path (70,086 B) and 70,155 B at the live path — the F20 §9.0 sentence
  (265 B) crossed the round-1 figure (179 B under) with no protocol change beyond
  frozen A.5 scope. Both reviewers hold the guardrail's ratified phase-1 scope
  sufficient for this release (claude-1 disposition (i) with the emphasis recorded;
  kimi-1 open question 3: "sufficient — the ratified scope is phase 1, and DF-3
  exists"). No map or ceiling change in-cycle; the per-phase policy question remains
  DF-3's.
- **DF-4 `release-binary-reproducibility`** — unchanged; inactive candidate. Recorded
  contribution this round (claude-1, PRIMARY experiment): the round-1 hash variance is
  fully explained by the **embedded build path** — same clean commit built from three
  directories yields three hashes, and `-trimpath` collapses them to one; kimi-1's
  same-clean-clone rebuild reproduced the implementer's recorded sha256
  (`8bd06646…`) byte-identically. DF-4 owns the release-flags story (documented
  `-trimpath` vs documented paths); no launch, no release-step action here.
- **Advisory, no slug opened (from claude-1 R2-MAJ-3, split by its own open question
  2):** `internal/trajectory` (589.7 s) and `internal/app` (521.7 s) package durations
  are a latent CI-failure risk on any runner and worth a look as release-adjacent
  hardening. Recorded here so it is not lost; opening a slug is the organizer's /
  owner's call — nothing is launched by this consensus and no mandatory scope is
  created.

## Dismissed findings

None. No finding from either reviewer is dismissed, withdrawn, or downgraded by this
draft. Two non-dismissal notes for the record: kimi-1's K2-F5 was folded into K2-F4
**by its own filer** before this consensus (ID stability), and claude-1's R2-NIT-3
needs no fix item because its owner (the organizer) corrected it in its own artifact
before this draft — recorded in the map as resolved-by-owner, which is a fix, not a
dismissal. Dismissal or withdrawal can only happen through a reviewer's own
self-correction or a signoff resolution, never through this drafter.

## Coverage & blind spots

*(Advisory; signoffs remain the gate.)*

**Both reviewers independently (PRIMARY each, separate isolated worktrees and
binaries):** all 21 cycle-1 fixes hold at code and real-entrypoint level; every FINAL
A–D acceptance row re-attacked and passing except the two elements claude-1 names
(R-2's floor half; the plan-command half of cross-cutting green); the `--json` trailer
defect (claude-1 isolated the streams; kimi-1 verified the envelope shape live — the
round's strongest cross-confirmation); the §15-region 8,056 B figure mismatch (two
identical measurements against the recorded 8,041 B); the staged core re-verified in
full (same sha256 `fc907e59…`, 109,772 B, three change sets, 0 unexplained, project
zones preserved); both suites green with `-timeout 2400s` (31/31 packages; skill 399
pass); DF-1..DF-4 exist as real inactive candidates; the owner publish gate
unexercised by everyone.

**Only claude-1 saw:** the floor's absence/mislabel/undercount and its `Parse`
block-splitting cause (R2-MAJ-1); the CI timeout margin and the failing no-flag run of
FINAL's own named command (R2-MAJ-3); next-action state (b) — no `consensus.md` at all
(R2-MIN-1b); the core-publish note's stale integrity facts (R2-MIN-3); the deviations
gap for F1's mechanism removal (R2-MIN-4); the fail-open on unparseable escalation
notes (R2-MIN-2); NIT-1/2/3; the DF-4 mechanism experiment (`-trimpath`).

**Only kimi-1 saw:** K2-F2 — the preflight fail-open on an unreadable workspace (the
round's only defect claude-1's fresh full-scope pass did not reach; a genuine
methodological gap in claude-1's coverage, acknowledged here for the record); the
F10 directory-vs-file trigger precision (K2-F4); the byte-identical two-party rebuild
of the F14 binary from the same clean clone.

**Blind spots — what no reviewer covered:**

1. **GitHub-hosted runners are unmeasured.** R2-MAJ-3's hosted-runner comparison is an
   inference claude-1 explicitly tags as one; the Windows leg's real wall-clock is
   unknown until CI runs. G3's flag removes the known cliff, not the unknown.
2. **The G-fixes do not exist yet.** Everything in `## Agreed fixes` is proposed, not
   implemented; round 3 must verify all ten (and the two contingency resolutions) —
   this consensus verifies nothing about code that does not exist.
3. **Attribution windows remain test-proven only.** No driver transition has occurred
   since F7 landed (`64a622c`, 2026-09-24T01:48Z — the newest run record
   `20260924T003301…Z` predates it); all live manifests are still zero-width and the
   live ledger honestly says `ambiguous`. The first real Phase-8-cycle transition will
   be the first live proof.
4. **The `--json` fix shape is undecided** (stderr vs `result` field — VC-5's
   sub-question); whichever lands changes what round 3 must re-verify on all four exit
   paths.
5. **Windows/`wait`/`usage` portability** is still macOS-verified only (carried from
   cycle 1; unchanged this round).

**Correlated agreement (§15.6(b)).** As in cycle 1: two model families, but the
cross-confirmations rest on shared priors (PRIMARY execution over prose; deterministic
tooling trustworthy) and partially identical command sequences (same tests, same
fixtures, same measurements). Their agreement is strongest where least independent;
this round's value is again the divergences — VC-4/VC-5/VC-6 — which the signoff
requests below ask each reviewer to re-examine against the other's evidence rather
than their own.

## Verdict conflicts

*(§15.3: quoted verbatim with author, tag, and evidence. All three are **DISPUTED** and
stay open; this draft closes over none of them. Resolution is by reviewer
self-correction (§15.1) or signoff — never by participant count (the reviewers split
1-vs-1), never by the drafter (the implementer issues no self-verdict), and never by
the organizer (verdicts read only, no code verdict). Closure of the idea — the next,
zero-fix consensus — depends on all three.)*

**VC-4 — "Is the implementation ready for a zero-fix closing consensus?"**

- claude-1, `review/round-02/claude-1.md`, closing position — **NOT READY**, 🟡 on a
  fix-up cycle: "**Ready for a zero-fix closing consensus?** **No — but close.** Three
  MAJOR and eight lesser findings are open, and Phase 8 closes only on a Phase-7
  consensus listing zero Agreed fixes. None is a blocker: there is no CRITICAL, no
  gating property is broken, no scope crept, and the trajectory is the converging
  shape §4's stopping judgment describes … I judge one more narrow fix-up cycle
  sufficient. My round-01 ❌ BLOCK is **withdrawn**; my position for the Phase-7
  consensus is 🟡 ACCEPT-WITH-RESERVATIONS on a fix-up cycle, not on closure."
- kimi-1, `review/round-02/kimi-1.md`, verdict — **READY**, ✅-shaped: "I consider the
  implementation **ready for a zero-fix closing consensus** (default close rule —
  `strict_gate` is absent from `00-prompt.md`), with the NITs dispositioned by that
  consensus." and "The four NITs are dispositionable by that consensus (record
  corrections and/or DF-slug carry) — in my judgment none requires a further fix-up
  cycle."
- **Status: DISPUTED — and it cannot be resolved inside this draft.** The two
  positions are not factually incompatible (both agree: no CRITICAL, no broken gating
  property, all 21 cycle-1 fixes hold); they differ on whether the residual set must
  be fixed before close. This draft proposes the fix cycle (G1–G10), which makes the
  question concrete: sustaining the plan at signoff is a vote for claude-1's branch;
  striking the code items is a vote for kimi-1's. Either way this consensus is not
  the closing consensus (it lists Agreed fixes); the zero-fix consensus — whichever
  round files it — requires this conflict resolved by the signoffs, a reviewer
  self-correction, or an operator ruling.

**VC-5 — severity (and close-timing) of the `parley wait --json` stdout trailer:
R2-MAJ-2 [MAJOR] vs K2-F3 [NIT].**

- claude-1, R2-MAJ-2 — **MAJOR**, PRIMARY (streams isolated: `1>out.txt 2>err.txt`,
  exit 3 → `stderr bytes: 0`, stdout fails `json.load` with `Extra data: line 51`;
  exit-0 path same shape; exit-4 clean and parses): "the machine-readable path is
  unusable on the exit-0 and exit-3 routes … the wrapper is not merging streams —
  `parley wait --json` genuinely emits invalid JSON on stdout on both the exit-0 and
  exit-3 paths." Grounds: "FINAL B.3 offers `--json` precisely so an organizer can
  branch mechanically, and codex-1 hit the failure (`jq: parse error …`) in ordinary
  use this round."
- kimi-1, K2-F3 — **NIT**, PRIMARY (live envelope verified): "`parley wait --json`
  stdout is not machine-parseable as a whole … Present since `86d028b` (the
  bare-digest shape had the same trailer) — NOT a fix-up regression — but B is new,
  unreleased surface, and 1.49.0 will freeze it." Suggested fix: "route the status
  trailer to stderr (or document the envelope-plus-trailer contract and the
  `{notes, digest}` keys in the skill) before release."
- **Status: DISPUTED (severity only — the facts carry both reviewers' independent
  PRIMARY provenance, and codex-1's live hit converges).** The filers agree on the
  defect, the mechanism, and even on the suggested fix direction; they disagree on
  whether it must land before close. G2 is proposed contingent on this conflict: if
  both sign for G2, the conflict closes as resolved-by-fix; if kimi-1 sustains
  close-now, the conflict stays open and the closing consensus waits on it. The
  shape sub-question (stderr vs envelope `result` field — claude-1's open question 3;
  both filers suggest stderr first) is a signoff choice recorded in G2.

**VC-6 — materiality of the two next-action states: R2-MIN-1 [MINOR, fix in-cycle]
vs K2-F4 [NIT, record-only].**

- claude-1, R2-MIN-1 — **MINOR, fix in-cycle**, PRIMARY (both states reproduced; (a)
  is this deck's live digest and `organizer brief` output): "The fixed enumeration
  (FINAL B.2) is respected; the *selection* is wrong. … Recording a known-wrong
  output is the right transparency; it is not a disposition that retires the finding,
  and the review-brief rule says a disputed finding closes only by reviewer
  withdrawal, consensus, or an operator ruling."
- kimi-1, K2-F4 — **NIT, record-only**, PRIMARY (the trigger isolated: directory
  creation, not first file): "Behavior is acceptable in both states (the digest line
  above it shows `implementation: present=true status=fix-up-cycle-1`); only the
  record's wording is imprecise."
- **Status: DISPUTED (materiality/scope only — both measured the same states).** G4
  proposes the in-cycle branch on the round-1 precedent (MAJ-3/MIN-2 were taken
  in-cycle as pre-existing vocabulary gaps B first made visible, per claude-1's open
  question 4); kimi-1's record-only position is the minority position and stays open
  for kimi-1 to sustain with a counter-proposal or concur at signoff.

**§15.3 dependency check.** This consensus's proposed fixes G2 (via VC-5) and G4 (via
VC-6) are explicitly contingent; no acceptance statement in this draft is derived
from a disputed claim — the facts under all three conflicts (the trailer is on
stdout; the two next-action states read `await implementation`; no CRITICAL exists
and all gating properties hold) carry both reviewers' independent PRIMARY provenance.
Closure of the idea requires VC-4, VC-5, and VC-6 to be resolved by signoff,
reviewer self-correction, or operator ruling; a 1-vs-1 reviewer split is never a
majority, and the implementer-drafter's proposals above settle nothing by themselves.

## Signoffs

*(Deliberation track: all three participants sign — reviewers and implementer.
Signoffs are intentionally empty in this draft; the cycle-1 signoffs live verbatim in
`review/consensus-cycle-01.md` and are NOT reused. Each signoff request asks the
signer to: (1) independently evaluate and, where feasible, reproduce the opposing
evidence — kimi-1: claude-1's stream isolation (`parley wait --json … 1>out 2>err` on
the exit-3 and exit-0 routes) and the no-`consensus.md` fixture from R2-MIN-1(b);
claude-1: kimi-1's no-`COOPERATION.md` preflight fixture from K2-F2 — and
self-correct their own claims if warranted (§15.1 `SELF-CORRECTION` in the signer's
own artifact); (2) state a binding position on VC-4 (fix-cycle vs close-now), VC-5
(severity/shape of the JSON fix, or close-now), and VC-6 (in-cycle vs record-only);
(3) confirm or amend the G1–G10 plan and the duplicate map — every one of their
findings must be traceable to a disposition; (4) pick the G2 shape if sustaining G2
(stderr primary vs envelope `result` field); (5) check G6's replacement facts against
the staged file themselves (`shasum -a 256`, `wc -c`) before signing. All three
conflicts stay open for exactly this independent reproduction and self-correction
during signoff — none is closed by this draft. Append-only; do not edit others'
blocks. The implementer signs per §15.1 with no self-verification verdict on its own
implementation.)*

### Signoff: claude-1 — (pending)

### Signoff: kimi-1 — (pending)

### Signoff: zcode-1 — (pending)

### Signoff: claude-1 — 2026-09-24
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept this consensus as an accurate and complete record of review round 2, and I
sustain the G1–G10 plan as drafted. This is a fix-cycle consensus, not a closing one: my
acceptance covers the record and the plan, **not** completion. All 11 of my findings are
traceable to a disposition (R2-MAJ-1→G1, -2→G2, -3→G3; R2-MIN-1→G4, -2→G5, -3→G6, -4→G7;
R2-NIT-1→G4, -2→G8, -3→resolved-by-owner, -4→G9); the duplicate merges (R2-MAJ-2 ≡ K2-F3,
R2-NIT-4 ≡ K2-F1) are correct; nothing of mine was dismissed, downgraded, or narrowed. G1
takes branch 1 of my open question 1 (retained whole-section floor, FINAL as frozen) — the
branch that needs no acceptance-wording amendment; I confirm it.

**Independent verification performed for this signoff** (PRIMARY — commands run this session
against `118b245`, binary built from `worktrees/lean-organizer-review-claude-1`):

- **K2-F2 reproduced — I concur with kimi-1 and with G10.** Fixture: `agents.toml` +
  `meta/version.json` + `ideas/probe/00-prompt.md` declaring `facilitator: claude-1` inside
  `participants:` with no `facilitator_participates`. Without `COOPERATION.md`,
  `parley preflight --dir <fx>` prints `Ready: no pending gates.` and exits **0**; copying
  `COOPERATION.md` into the same fixture fires the gate and exits **3**. Same deck, same
  conflict — the only variable is whether `ReadWorkspaceStatus` succeeds, so the gate
  silently never ran. That is the status-concealment class this idea exists to fix, and it is
  the one defect my own full-scope pass did not reach: kimi-1's catch, now held PRIMARY by a
  second reviewer. Fail-closed (exit 1) is the right remedy.
- **G6's replacement facts confirmed before signing, as asked.**
  `shasum -a 256 ~/.parley/staging/COOPERATION-2.13.0.md` → `fc907e59…62c9f`; `wc -c` →
  **109,772 B**; mtime `2026-09-24 03:35:27`. The note still reads `109,507 B` /
  "verified 2026-09-23" / no sha256 — R2-MIN-3 stands and G6 is needed.
  `~/.parley/protocol/core/` holds only `2.10.0`: no publish has occurred.
- **R2-NIT-3's retirement confirmed independently.** Zero literal `\n` escapes remain in
  `organizer-notes.md` and it ends with a real `0x0a`. Retiring my own NIT with no fix item
  is correct; it was fixed by its owner, not dismissed.
- **My own disputed claims re-verified rather than restated.**
  `parley wait --json … 1>out 2>err` on the exit-3 route: stdout 259 B = envelope **plus**
  the `wait: timeout after 3s; outstanding: …` trailer, stderr **0 B**, `json.load` fails
  `Extra data: line 11 column 1`. The same run independently reproduced R2-MIN-1(b)
  (`"next": "await implementation"` on a deck with no `consensus.md`) and R2-NIT-1 (the
  generic `awaited condition not yet reached` where the scope is consensus).

**Binding positions on the three open conflicts:**

- **VC-4 — NOT READY, sustained.** This is not a severity judgment: two frozen-FINAL
  acceptance elements are unmet — C.3/R-2's floor half (R2-MAJ-1) and cross-cutting "both
  suites green" at the command FINAL itself names (R2-MAJ-3). An implementation does not
  close over an unmet acceptance row. I record that kimi-1's "ready" was formed against its
  own 4-NIT finding set and is not factually incompatible with mine; the proposed cycle
  closes the gap under either reading. No CRITICAL, no broken gating property, converging
  trajectory — one more narrow fix-up cycle, exactly as my round-02 file concluded.
- **VC-5 — MAJOR sustained; G2 in-cycle.** I accept kimi-1's fact that the trailer predates
  fix-up 1 (`86d028b`) and is not a regression. It does not lower the severity: B is
  unreleased surface that 1.49.0 freezes, FINAL B.3 offers `--json` precisely for mechanical
  branching, and codex-1 hit the failure in ordinary use. kimi-1 itself asks for the fix
  before release, so the live disagreement is timing, not substance.
- **VC-6 — MINOR, in-cycle, sustained (G4 whole, including K2-F4's record half).** Having
  now reproduced state (b) in a clean fixture: the next-action line is the organizer's branch
  point, and emitting `await implementation` when a consensus is awaited is wrong output, not
  imprecise wording. Round 1 took the same class (MAJ-3/MIN-2) in-cycle as vocabulary gaps B
  first exposed; consistency applies. kimi-1's record-only position stays honestly recorded.

**G2 shape — I pick stderr** (terminal status line to stderr; stdout carries the envelope
alone). In round 02 I had no preference between this and an envelope `result` field; having
re-run it I do: the envelope already churned once this cycle (bare `PhaseDigest` →
`{notes?, digest}`), and a second change inside the same unreleased cycle enlarges what round
3 must re-verify on all four exit paths for no functional gain. Both satisfy B.3, so I do not
block the `result`-field alternative if kimi-1 and the implementer prefer it.

**Reservations carried as open items (not resolved by this signoff):** the ten G-fixes do not
exist yet, and round 3 must verify all ten plus the two contingency resolutions; GitHub-hosted
runner wall-clock remains unmeasured (G3 removes the known cliff, not the unknown); live
attribution windows remain test-proven only; Windows `wait`/`usage` portability remains
macOS-verified only. I am **not** calling this implementation complete, and I have released,
published, merged, tagged and installed nothing — `parley protocol publish` was not run.
Counter-proposal (required if ❌): N/A — this is not a block.

### Signoff: kimi-1 — 2026-09-24
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: I accept this consensus as an accurate and complete record of review round 2, and I
sustain the full G1–G10 plan as drafted. My acceptance covers the record and the plan — **not**
closure: this is a fix-cycle consensus, and Phase 8 closes only on a later zero-fix consensus.
All four of my active findings are traceable to a disposition (K2-F1→G9, K2-F2→G10, K2-F3→G2,
K2-F4→G4 record half); the duplicate merges (R2-MAJ-2 ≡ K2-F3, R2-NIT-4 ≡ K2-F1) and the K2-F5
fold into K2-F4 by its own filer are correctly recorded; nothing of mine was dismissed,
downgraded, or narrowed.

**Independent verification performed for this signoff** (PRIMARY — my own round-02 binary
`/tmp/kimi1-r2/parley`, sha256 `8bd06646…` = the F14 record, re-hashed this session; fixtures
under `/tmp/kimi1-signoff/`, disposable):

- **claude-1's `--json` stream isolation reproduced on BOTH routes, as the signoff request
  asked.** Exit-3 route (fresh fixture: complete round-01, no `consensus.md`;
  `parley wait --idea wfx --for consensus --timeout 3s --json 1>out 2>err`): exit 3, stderr
  **0 B**, stdout = the `{notes?, digest}` envelope **plus** the `wait: timeout after 3s;
  outstanding: …` trailer, `json.load` fails `Extra data: line 53`. Exit-0 route (live deck at
  `118b245`, `--for round`): exit 0, stderr **0 B**, stdout ends `wait: boundary reached (round
  complete)` after the envelope, `json.load` fails `Extra data: line 126`. The defect stands on
  both machine-branching routes; my round-02 NIT underweighted it (see VC-5 below).
- **R2-MIN-1(b) reproduced.** Same fixture, round-01 2/2 filed-and-valid, no `consensus.md`:
  digest reads `implementation: present=false` and **`next: await implementation`** — a Phase-2/3
  deck told to skip consensus. In round 2 I had only isolated the directory-vs-file trigger
  (state a); state (b) is wrong output, not imprecise wording.
- **G6's replacement facts checked before signing, as asked.**
  `shasum -a 256 ~/.parley/staging/COOPERATION-2.13.0.md` → `fc907e59…62c9f`; `wc -c` →
  **109,772 B**; mtime `Sep 24 03:35:27 2026`. The core-publish note still reads `109,507 B` /
  "verified 2026-09-23" / no sha256 (grep-confirmed) — R2-MIN-3 stands and G6 is needed.
  `~/.parley/protocol/core/` holds only `2.10.0`: no publish has occurred.

**SELF-CORRECTIONs (§15.1 — my own claims, each replaced statement named; weakenings):**

- **VC-4 — I withdraw "ready for a zero-fix closing consensus" and join the fix-cycle branch.**
  Replaces my round-02 claims "unresolved acceptance criteria within the frozen FINAL table:
  none" and "ready for a zero-fix closing consensus". Those were formed by checking C.3's
  gating half (retention/omission sets, omission index, source hash) and by running both suites
  at `-timeout 2400s`; I never checked the floor-recording half of R-2 against FINAL's text, and
  I never ran the plan's own no-flag command. claude-1's R2-MAJ-1 (PRIMARY `Parse` probe: an
  omission sum mislabelled `floor`, 17 subsection blocks / 11,903 B silently excluded) and
  R2-MAJ-3 (PRIMARY: 589.7 s of the 600 s default; their own no-flag run failed) name two
  acceptance elements unmet as recorded, and I hold no counter-evidence against either. An
  implementation should not close over unmet acceptance elements; the proposed cycle is narrow
  and stays inside frozen FINAL. My earlier position was honestly formed against my own 4-NIT
  finding set and is honestly replaced.
- **VC-5 — I accept the in-cycle timing (G2) and no longer contest the MAJOR framing; the
  conflict closes as resolved-by-fix.** The facts were never disputed (both filers PRIMARY,
  codex-1's live hit converging); my NIT rested on "fix before release", and under this protocol
  the only pre-release window for participant code changes is a fix-up cycle — closing now would
  freeze the defect into 1.49.0, exactly the outcome my own review warned against. FINAL B.3
  exists for mechanical branching and the defect bit the organizer in ordinary use; the severity
  follows the contract's purpose.
- **VC-6 — I withdraw the record-only position and concur with G4 whole (both states plus
  R2-NIT-1 plus K2-F4's record half).** Having now reproduced state (b) myself: emitting
  `await implementation` where a consensus is awaited is wrong output on the organizer's branch
  surface. Round-1 precedent (MAJ-3/MIN-2 taken in-cycle as vocabulary gaps B first exposed)
  applies, as claude-1's open question 4 proposes.

**G2 shape — stderr.** Terminal status line to stderr; stdout carries the `{notes?, digest}`
envelope alone; non-`--json` human output unchanged. Same reasoning as claude-1: the envelope
already churned once this cycle (bare `PhaseDigest` → `{notes?, digest}`), and a `result` field
enlarges round 3's re-verification surface on all four exit paths for no functional gain. The
skill's `wait` section must state the chosen stream contract and the envelope keys (my original
K2-F3 ask, carried in G2), and the round-1→fix-up-1 envelope change rides the same unreleased
1.49.0 freeze as recorded.

**Remaining plan positions:** G1 branch 1 confirmed (record the retained whole-section floor in
the same test run, FINAL as frozen — no acceptance-wording amendment, no deviation; the gating
half of C.3 and the guardrail untouched). G3 in-cycle confirmed (two-token change on a file this
idea ships; the Windows-honesty claim of FINAL risk 8 rests on the leg being able to finish; the
package-duration question correctly split to advisory with no slug launched). G5, G7, G8, G9,
G10 confirmed as proposed — G10 is my own finding with my own fail-closed leaning (exit 1 naming
the read failure), and G9 is my own figure correction with the extraction boundary stated beside
the number.

**Reservations carried as open items (not resolved by this signoff):** the ten G-fixes do not
exist yet — round 3 must verify all ten plus the two contingency resolutions, including
`--json` stdout decode on all four exit paths at the chosen stderr shape; GitHub-hosted runner
wall-clock remains unmeasured (G3 removes the known cliff, not the unknown); live attribution
windows remain test-proven only (no driver transition since F7 landed); Windows
`wait`/`usage` portability remains macOS-verified only. I am **not** calling this implementation
complete, and I have released, published, merged, tagged and installed nothing —
`parley protocol publish` was not run.
Counter-proposal (required if ❌): N/A — this is not a block.

### Signoff: zcode-1 — 2026-09-24
Status: ✅ ACCEPT
Notes: As the implementer (and drafter of this record), I accept this consensus as an accurate
and complete Phase-7 record of review cycle 2, and I accept the full G1–G10 agreed-fix plan as
sustained by both reviewers. Per §15.1 I issue **no verification verdict on my own
implementation**: every code-level verdict above belongs to claude-1 and kimi-1 (each tagged
PRIMARY in their own files and re-verified in their signoff blocks); I take R2-MAJ-1/-2/-3 and
every lesser finding as fixes to implement, not claims I adjudicate. My acceptance covers the
record and the plan and commits to fix-up cycle 2 — it does not claim completion: this consensus
lists ten Agreed fixes and therefore is not a closing consensus under any position in VC-4;
closure remains the later zero-fix consensus's call.

**Record facts re-verified this session before signing (PRIMARY — record/compilation facts
only, not implementation self-verification):** HEAD is `118b245`; `git diff
120a9bf..118b245 -- …/FINAL.md` is empty (FINAL frozen); staged core
`~/.parley/staging/COOPERATION-2.13.0.md` = sha256 `fc907e59…62c9f`, **109,772 B**, mtime
`Sep 24 03:35:27 2026`, and `~/.parley/protocol/core/` holds only `2.10.0` (no publish has
occurred); `organizer-notes.md` carries zero literal `\n` occurrences and ends with a real
`0x0a` byte (R2-NIT-3 resolved-by-owner independently confirmed); the core-publish note still
reads `109,507 B` / "verified 2026-09-23" / no sha256 (grep-confirmed — R2-MIN-3 stands, G6 is
needed); the VC-4 quotations were spot-checked verbatim against `review/round-02/claude-1.md`
("No — but close"; ❌ BLOCK withdrawn) and `review/round-02/kimi-1.md` ("ready for a zero-fix
closing consensus" — since withdrawn by kimi-1's own §15.1 self-correction in its signoff
above).

**Conflicts after the signoffs (record observation, not adjudication):** with kimi-1's §15.1
self-corrections and claude-1's sustained positions, all three conflicts close through the
signoff mechanism this consensus itself named — never by count, drafter, or organizer:
**VC-4** closes on the fix-cycle branch (kimi-1 withdrew "ready"; claude-1 sustained
NOT-READY on two unmet frozen-FINAL acceptance elements — C.3/R-2's floor half and the
plan-command half of cross-cutting green); **VC-5** closes as resolved-by-fix at the **stderr**
shape (both reviewers independently pick stderr; I concur — it is the smallest change and
avoids a second envelope churn inside one unreleased 1.49.0 freeze); **VC-6** closes in-cycle
with G4 whole (kimi-1 withdrew record-only after itself reproducing state (b)). No dispute is
carried into cycle 2; the closing record is the two signoff blocks above, not this note.

**Commitments for fix-up cycle 2:** implement G1–G10 exactly as specified above at the
converged shapes — G1 branch 1 (record the retained whole-section floor in the same test run,
FINAL as frozen, gating half and guardrail untouched); G2 at the stderr shape with the
four-exit-path `encoding/json` decode test, the skill `wait` stream-contract update, and the
envelope-churn freeze note; G3 (`-timeout 45m` in `tests.yml` plus the FINAL-named command in
`IMPLEMENTATION.md` aligned to `-count=1 -timeout 2400s`); G4 (both next-action repairs, the
`outstandingAgents` "consensus.md not filed" line, and the F10 record corrected to the true
directory-creation trigger); G5 (annotation plus `notes` exposure for unparseable /
`idea:`-less `to-user` notes); G6 (the core-publish note's replacement facts matching the
staged file at the landing HEAD, per the restage rule); G7 (the A.4 deviations bullet pointing
at DF-1); G8 (the mtime-approximation comment plus the `loop.go` in-place-rewrite record, no
behavior change); G9 (the 8,056 B figure with the extraction boundary stated); G10
(fail-closed preflight exit 1 naming the read failure, with the negative fixture test) — each
with its regression test; both suites green at the cycle-2 HEAD under the G3-named command;
`IMPLEMENTATION.md` receives the fix-up-cycle-2 section with per-fix commit references and
bumped frontmatter, published ready for re-review. No release, merge, tag, global install, or
publish occurs in the cycle; release remains the organizer's post-Phase-8 step and
`parley protocol publish` remains the owner's attended action.

**Reservations:** none of my own on the record or the plan. The open items the reviewers carry
in their signoff blocks — the ten G-fixes do not exist yet and round 3 must verify all of them
including `--json` stdout decode on all four exit paths at the stderr shape; GitHub-hosted
runner wall-clock unmeasured (G3 removes the known cliff, not the unknown); live attribution
windows test-proven only; Windows `wait`/`usage` portability macOS-verified only — stand
exactly as they wrote them; I do not countermand or re-verdict any of them.
Counter-proposal (required if ❌): N/A — this is not a block.
