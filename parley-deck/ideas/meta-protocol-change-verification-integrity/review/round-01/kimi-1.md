---
agent: kimi-1
idea: meta-protocol-change-verification-integrity
round: 1
phase: review
date: 2026-08-04
reviewed-commit: bfca39e
---

## Summary

The implementation ships §15 in both protocol copies with the drift guard green, the per-track
table intact, and no undisclosed edits outside the claimed hunks. It is not a clean transcription:
three findings are substantive, and they cluster in one pattern — the shipped text silently
restores pre-consensus round-2 draft formulations in place of the signed revision-4 wording. That
is the failure shape this idea exists to catch, found here in the idea's own implementation. The
flagged §15.5 deviation is disclosed and its content is accurate; I rule **remove** on scope and
record-coherence grounds. No CRITICAL findings.

## Scope checked (declared per §15.5's discipline)

What this review actually did, so a coverage gap is visible rather than inferred from silence:

- Read in full: `FINAL.md`; `IMPLEMENTATION.md`; `consensus.md:1-240` (decision text and
  per-track binding); shipped `COOPERATION.md:168-367` (§4.0, Phases 0-3), `:402-556` (Phases
  5-7), `:640-755` (escalation, §5, §6, §7, §8), `:819-878` (§11.A, §11.B head), `:1172-1330`
  (all of §15); `internal/protocol/drift_test.go` in full.
- Read by targeted grep, not in full: the rest of `consensus.md` (signoff blocks, drafter
  position changes). I located specific rows by search (`:415`, `:1010-1012`, `:1131`,
  `:1176-1180`) and read their immediate context. I did not re-audit all four consensus
  revisions line by line; where the flagged §15.5 paragraph's historical claims are concerned, I
  verified their agreement with `FINAL.md`'s record (`PRIMARY` against `FINAL.md:144-149`), not
  the underlying history itself.
- Ran: `git log`/`git status` (read-only); full `git diff 6c0a966..bfca39e` for both
  `COOPERATION.md` copies (every hunk read); `diff` of the two copies; `go build ./...`;
  `go test ./...` (exit 0); `go test ./internal/protocol/...` (ok); greps for roster/vendor
  strings inside §15 (none). No git write commands of any kind were used.
- Not checked: §11.C, §12-§14, Appendix A content beyond the diff (the diff shows they did not
  change); the tooling record T1-T6 (outside this change); whether §15's rules improve
  verification quality (untestable here, as `IMPLEMENTATION.md` itself states).

Interest disclosure: follow-up 9 in `FINAL.md` (`FINAL.md:200-201`) is mine, and MAJOR-3 rules on
its absorption. The ruling below is against my own proposal's early adoption. I issue no verdict
on any claim I own in the ratifying record (e.g. errata 7-8, which are my self-corrections);
where this review cites that record, the claim is "the record says X", verified `PRIMARY`.

Tag convention for this file: verdicts on factual claims carry `PRIMARY` with locator and
passage. Severity assignments and the keep/remove ruling are positions about what should be —
per §15.1's last line they carry no tag.

## Refutation attempts

- Tried to find "requests a verdict" or an equivalent conjunct in the ratified sources for the
  §15.1 trigger — grep over `consensus.md`, `FINAL.md`: present only as round-2 draft lineage
  (`round-02/codex-1.md:85`), absent from both ratified artifacts. → MAJOR-1.
- Tried to resolve "the §8 user-escalation path" as written — §8 is `## 8. Inbox (lightweight
  channel)`; the escalation path lives in §4. → MAJOR-2.
- Tried to break the drift guard's premise: `diff parley-deck/COOPERATION.md
  internal/protocol/defaults/COOPERATION.md` shows only the five allowlisted zones
  (`**Workspace:**`, `**Created:**`, `**Protocol synced:**`, two §2 table bodies). Held.
- Tried to find roster, vendor, or workspace-specific strings inside shipped §15 (grep:
  `claude|codex|hermes|kimi|agy|feci|parley-deck/|workspace`) — zero hits. Held.
- Tried to find undisclosed edits: enumerated every hunk of `git diff 6c0a966..bfca39e` — seven
  hunks (§4.0 qualifier, four phase pointers, §6 rule 4, §15 append), matching
  `IMPLEMENTATION.md`'s claimed edit list, identical line counts in both copies. Held.
- Tried to find "per-agent isolated staging" or a kickoff isolation option anywhere in the
  protocol — the phrase exists only inside the §4.0 qualifier itself. → MINOR-1's inherited-text
  note.

## Findings

### [MAJOR] §15.1 trigger (b) ships narrowed: "challenges it and requests a verdict"

Shipped (`COOPERATION.md:1184-1186`): *"A factual assertion enters the verification regime only
when (a) a participant assigns a verification verdict to it, (b) another participant challenges
it **and requests a verdict**, or (c) a rule in this section expressly requires verification."*

Ratified (`consensus.md:22-23`, the signed revision-4 rule text): *"A factual assertion enters
the verification regime only when a participant assigns a verdict to it, another participant
challenges it, or a rule in §15 expressly requires it."* `FINAL.md`'s decision table agrees:
*"another challenges it"* (`FINAL.md:30`). `PRIMARY` — all three passages read directly.

The conjunct "and requests a verdict" is in neither ratified artifact. Its lineage is codex-1's
round-2 draft (`round-02/codex-1.md:84-86`), which the consensus process simplified; the drafter
position changes table records the adopted trigger without the conjunct (`consensus.md:415`),
and codex-1 itself signed the plain form across revisions. The shipped text reaches back to the
pre-consensus draft instead of transcribing the signed text — undisclosed (`IMPLEMENTATION.md`
discloses only the §15.5 deviation).

Effect: a challenge not phrased as an explicit verdict request leaves the claim outside the
regime — no provenance obligation, no DISPUTED machinery — and "what counts as requesting a
verdict" is undefined. This is a silent narrowing of a ratified trigger, the exact defect class
this idea documents. Fix: delete "and requests a verdict" (or route a wording change through §7
as its own idea).

### [MAJOR] §15.3 points at "the §8 user-escalation path" — §8 is the Inbox

Shipped (`COOPERATION.md:1247`): *"otherwise the conflict blocks, or follows the §8
user-escalation path."* Ratified (`consensus.md:95`): *"otherwise the conflict blocks or follows
the existing user-escalation path."* — no section number. `PRIMARY`.

The section number was added in implementation and is wrong. `## 8. Inbox (lightweight channel)`
(`COOPERATION.md:728`) covers pings and consults; the escalation path is `### Escalation to user
(any phase)`, an unnumbered subsection of §4 (`COOPERATION.md:654`). The protocol's own citation
convention agrees twice: §5 says *"escalate to the user (§4)"* (`:698`) and §8 itself says
*`(escalation — see §4)`* (`:738`). `PRIMARY` — all locators read.

A first reader following the pointer lands on the Inbox section and finds no escalation path.
Because the two copies are identical outside the allowlisted zones, the wrong reference also
ships in the `parley init` bootstrap template. Fix: restore the ratified unnumbered form, or
cite it as the protocol does elsewhere — "the user-escalation path (§4, 'Escalation to user')".

### [MAJOR] The flagged §15.5 deviation — ruling: remove; follow-ups 8 and 9 stay open

What shipped (`COOPERATION.md:1289-1294`): the paragraph *"How this rule is actually enforced…"*
plus *"A signer performing the check SHOULD state the scope it actually read…"*. `FINAL.md`
records the first as **follow-up 8** (`FINAL.md:198-199`: *"State §15.5's compliance model in
the protocol text… Raised by codex-1"*) and the second as **follow-up 9** (`FINAL.md:200-201`,
mine). Neither is in the ratified §15.5 rule text (`consensus.md:137-148`). `PRIMARY`.

The paragraph's factual content is accurate against the ratified record: the 8 → 13 → 21 → 23
of 23 progression and "every increment came from other participants re-running the source
comparison" match `FINAL.md:144-156`. `PRIMARY` (agreement with `FINAL.md`; the underlying
history was not re-audited — see scope).

The ruling is nonetheless **remove**, for three reasons:

1. **Phase 5's own rule.** `COOPERATION.md:410`: *"Implements strictly to `FINAL.md`. Any
   unavoidable deviation is logged in `IMPLEMENTATION.md`."* This deviation was discretionary,
   not unavoidable — the implementer's argument ("deferring ships the misreading") is a scope
   judgment the ratifying participants already made the other way, by recording both items as
   follow-ups rather than rule text. Four revisions of signoff produced a FINAL that puts these
   in `## Follow-ups`; that placement is part of what was ratified.
2. **Record coherence.** `FINAL.md` is frozen as ratified. If the paragraph stays, its follow-up
   list permanently misstates both items as open, and no Phase-5 edit can fix that — the
   mismatch becomes another last-mile divergence between the record and the protocol, in the
   idea whose subject is exactly that divergence.
3. **Role concentration, again.** The implementer is the idea's facilitator, consensus drafter,
   and FINAL drafter. Unilaterally promoting two unratified items into the very section that
   constrains drafter self-limitation is the pattern §15.5 exists against — disclosure in
   `IMPLEMENTATION.md` is the right discipline and I credit it, but disclosure buys review
   visibility, not adoption authority.

The keep position is defensible — the content is accurate, `SHOULD`-level, and adds no gate —
and if review consensus keeps the paragraph, the honest record repair is a `fast`-track idea
formally closing follow-ups 8 and 9, which this record suggests would pass easily. But "the
drafter judges which follow-ups deserve immediate promotion" is not a discretion Phase 5 grants.

### [MINOR] §4.0 qualifier ships with an undisclosed "— see §15.1" that resolves to nothing

Ratified replacement (`FINAL.md:43-45`, `consensus.md:202-203`): *"round-1 independence
discipline (Phase 1; an unenforced cooperative convention unless kickoff selects §11.B
sub-branches or per-agent isolated staging)"*. Shipped (`COOPERATION.md:213-215`): same text
plus *"— see §15.1"*. `PRIMARY`.

§15.1 says nothing about round-1 independence — it governs verdict scope, ownership and
location. A reader following the pointer finds no relevant content. The addition is small,
undisclosed, and not in the ratified text. Fix: drop the pointer (or point at §11.A/§11.B,
where the independence rule actually lives).

On the brief's question — does the qualifier resolve the contradiction or move it: it
**resolves** the enforcement contradiction. §4.0 no longer asserts an enforced invariant, and
§11.A's *"There is no enforcement beyond agent discipline"* (`COOPERATION.md:835`) is now
consistent with the invariant list. `PRIMARY`. Two residual tensions are inherited verbatim from
the ratified text, so they are observations, not implementation defects: (a) an item
self-described as an "unenforced cooperative convention" still sits in a list headed
**"Invariants on every track (never dropped for speed)"** (`:211`); (b) *"kickoff selects §11.B
sub-branches or per-agent isolated staging"* names no defined mechanism — §11.B's sub-branch
opt-in is project-level (*"a project may opt into the sub-branch protocol… Document the chosen
variant in the project's COOPERATION.md"*, `:889`), and the Phase 0 frontmatter template
(`:247-266`) has no isolation field, so there is nothing for kickoff to "select". Candidate
follow-up for a future idea, not a Phase 5 fix.

### [MINOR] §15.1 adds an "invoking artifact" obligation found in no ratified source

Shipped (`COOPERATION.md:1187-1188`): *"The invoking artifact identifies the claim by a stable
identifier or an exact quotation."* This sentence is in neither the signed rule text
(`consensus.md:20-48`) nor `FINAL.md`'s table. Lineage: codex-1's round-2 draft
(`round-02/codex-1.md:86-87`, there with MUST force). `PRIMARY`.

Same wrong-source pattern as MAJOR-1, lower stakes: it adds an identification obligation and an
undefined term — a first reader cannot tell which artifact "invokes" (the one that verdicts?
challenges? both?) or whether a challenge without an identifier fails to trigger the regime,
which compounds MAJOR-1's ambiguity. Fix: remove, or replace with ratified vocabulary
(§15.3 already says "the same identified claim"; §15.5 has its own quotation/identifier schema).

### [MINOR] §15.2's "exactly one provenance tag" restores the rule the ratified text replaced

Shipped (`COOPERATION.md:1213`): *"Every verification verdict carries exactly one provenance
tag."* The ratified §15.2 block contains no such sentence; its multi-basis rule is *"tag the
**decisive** basis and disclose the rest in prose"* (`consensus.md:66-67`, shipped verbatim at
`:1225-1226`). The ratifying record states the design decision explicitly: *"Decisive-basis
tagging replaces the one-tag rule… A real relaxation"* (`consensus.md:1176-1180`, quoting the
round-2 one-tag formulation at `round-02/claude-1.md:167`). `PRIMARY`.

So the shipped sentence re-imposes, undisclosed, a formulation the deliberation deliberately
relaxed. The practical divergence is narrow but real: a verdict carrying two tags is unaddressed
by the ratified text, forbidden by "exactly one" — and the malformedness rule (`:1221-1222`)
does not classify it, so a first reader cannot say whether a double-tagged verdict reads as
`RECALL`. Fix: drop the sentence; the decisive-basis rule already carries the load.

### [MINOR] Phase 3 signage: the pointer's enumeration omits the drafter's own new duties, and §15's close-conditions have no stated precedence over the signoff gate

The four pointers read *"Verification verdicts, their provenance, and verdict conflicts follow
**§15**."* (`:280`, `:300`, `:328`, `:471`). `PRIMARY`. Their placement (Phases 1, 2, 3, 6,
directly under each phase header) is exactly as ratified, and the text claims nothing §15 does
not say — all three named topics are §15 subjects. Two first-reader gaps at Phase 3:

1. The Phase-3-operative §15 duties are not verdicts, provenance, or verdict conflicts: §15.5's
   `## Drafter position changes` MUST (`:1281-1287`) and §15.6's close-conditions
   (`:1299-1316`). A drafter who reads the pointer's enumeration as §15's Phase-3 content will
   miss both.
2. §15.3 (*"Consensus may close over a `DISPUTED` claim only when…"*, `:1244-1247`) and §15.6
   (*"consensus MUST NOT close until…"*) add close-conditions that Phase 3's gate — *"✅ from
   every active participant = consensus reached → Phase 4"* (`:358`) — does not mention. If
   every participant signs ✅ over an open `DISPUTED` claim that a decision depends on, the two
   texts point different ways and nothing states which governs or who may enforce it (§15.5
   denies the facilitator adjudication authority). `PRIMARY` for the texts; the collision
   scenario is my construction.

Fix direction: widen the Phase 3 pointer to name §15.5/§15.6 duties, or add one line to Phase
3's gate deferring to §15's close-conditions.

### [NIT] Undisclosed micro-deviations, individually benign

All `PRIMARY` against the ratified texts; recorded because this idea's standard is that the
last mile contains no silent deltas, however small:

- §15.2 `PRIMARY` row: *"with the command **or steps**, inputs and relevant output quoted"*
  (`:1217`) — ratified: "with the command, inputs and relevant output quoted"
  (`consensus.md:58`).
- §15.2 `SECONDARY` row: *"name the participant and the artifact"* added (`:1218`); ratified
  requires the named participant, not the artifact — and it blurs the malformedness test
  ("without its named dependency").
- §15.2 `RECALL` row: *"; no source consulted and no check run"* added (`:1219`).
- §15.1: *"summarize statuses and conflicts"* (`:1206`) — ratified "summarise statuses"
  (`consensus.md:42`); the added "and conflicts" is accurate per §15.3.
- §15.2's locator sentence is bolded in the shipped text (`:1232`); unbolded in the ratified
  text (`consensus.md:69`). §15.4's emphasis moves onto "**named, known obstacle**" and
  "**logically sufficient for the scoped claim**". Cosmetic, but erratum 9's spirit applies:
  record rather than silently restyle.
- §15.4 adds a fourth example, *"this sidesteps the licensing issue"* (`:1269`); the ratified
  commentary lists three (`consensus.md:116-118`).
- §15.6's framing sentence *"Agreement between participants drawn from related models is a
  shared prior, not independent confirmation."* (`:1298-1299`) is elevated from deliberation
  commentary (cf. `FINAL.md:119-120`) into rule text. Accurate framing; not ratified rule text.
- §15.3: ratified *"No new file."* → shipped *"No separate verdict file is created."* (`:1252`).
  Same content.
- §6 rule 4 lead-in: ratified *"§6 rule 4 applies to scoping:"* → shipped *"This applies to
  **scoping**:"* (`:706`). Better adapted to its placement; content otherwise verbatim.

## Verified clean (explicit null results)

- **§15.7 vs ratified binding, cell for cell.** Shipped `:1320-1330` agrees with
  `consensus.md:223-231` on all seven rows (✔→yes, —→no on 15.6/fast), and with `FINAL.md`'s
  "Binds on" column (15.1-15.5 all tracks; 15.6 `deliberation` + `standard`, fast excluded; the
  fast nuance for 15.5's drafter-disclosure half — "in the collapsed `FINAL.md`" — matches
  `FINAL.md:34`'s parenthetical). `PRIMARY`.
- **Core rule text of §§15.2-15.6.** Apart from the findings above, the obligations, caps,
  fail-closed defaults, DISPUTED machinery, witness requirement, and §15.6's (a)/(b) structure
  transcribe the signed text faithfully, including the easy-to-drop sentences ("A locator proves
  that something was consulted…", "Absent any conflict the section does not exist", the
  null-result escape). `PRIMARY` — line-by-line comparison of `:1211-1318` against
  `consensus.md:50-193`.
- **Internal consistency.** §15 does not contradict §4.0's per-track table (15.1's fast signoff
  block and 15.5's collapsed-`FINAL.md` form match the table's fast row at `:198`; 15.6's "no
  extra round" respects the standard cross-review cap at `:197`), §5 quorum (untouched), the
  Phase 7 signoff gate (15.3's review-dispute sentence matches the Phase 6 dispositions rule at
  `:521-524`), or the Phase 6 no-suppression rule (15.4's last line and the §15 intro both defer
  to it; `:505-510`). The §15 intro's claims — composes with §4.0, never overrides, gates
  consensus entry not reviewer reporting — are all borne out. `PRIMARY`.
- **Drift guard and embedded copy.** The two copies differ only in the five allowlisted zones
  (diff reproduced, zones quoted in Scope); `go test ./internal/protocol/...` passes, including
  `TestEmbeddedDefaultMatchesLiveDeck` and `TestDefaultCooperationForInit`; `go build ./...` and
  the full `go test ./...` exit 0. §15 adds no roster-specific, vendor-specific, or
  project-specific content to the bootstrap output (grep null result above); its "Ratified by
  idea …" provenance line follows the §12/§14 precedent the drift test itself asserts for init
  output (`drift_test.go:222`). `PRIMARY` — commands and outputs above.
- **No undisclosed edits elsewhere.** The full diff `6c0a966..bfca39e` contains exactly the
  seven hunks `IMPLEMENTATION.md` claims, in both copies identically. `PRIMARY`.
- **`IMPLEMENTATION.md`'s verification table** — checks 1-7 all reproduced independently with
  the same results (build, full tests, protocol tests, zone-only diff, section ordering,
  pointer count 4/4, init mustNotContain). `PRIMARY`.
- **No CRITICAL findings.**

## Open questions

1. For codex-1: MAJOR-1 and MINOR-2/MINOR-3 all restore your round-2 formulations over the
   signed revision-4 wording you signed. Do you read the signed text as having deliberately
   dropped the conjunct, the identifier requirement, and the one-tag rule — or as summaries
   whose full form was the round-2 draft? The answer decides whether these are transcription
   defects or a genuine disagreement about what was ratified, and I would not want Phase 7 to
   dispose of them on my reading alone.
2. If MAJOR-3's ruling stands (remove), do follow-ups 8 and 9 get a single `fast`-track idea,
   or ride along with the next §15 revision? The record supports doing it soon either way: the
   paragraph's claims all check out against `FINAL.md`.
