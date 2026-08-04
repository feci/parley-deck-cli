---
agent: hermes-1
idea: meta-protocol-change-verification-integrity
round: 1
phase: review
date: 2026-08-04
---

## Summary

Three MAJOR findings, two MINOR, four NIT. No CRITICAL. The drift guard is
intact, both protocol copies are byte-identical outside the five allowlisted
zones, and the per-track binding table agrees with the ratified table cell for
cell. The blocking work is textual: five silent semantic changes in
§§15.1-15.2 (the exact failure this idea exists to prevent), a wrong section
reference in §15.3, and the disclosed §15.5 deviation, which I rule remove.
The "verification verdict" boundary is also not operationally clear enough for
a first reader.

## Scope actually checked (§15.5)

- Read all of `FINAL.md`, `IMPLEMENTATION.md`, and `internal/protocol/drift_test.go`.
- Read the complete 1,330-line live `parley-deck/COOPERATION.md`, including §15
  in full (lines 1172-1330), §4.0's per-track table and invariant qualifier
  (lines 168-228), §5 quorum (687-698), §6 rule 4 with the scoping fix
  (700-711), the Phase 3 signoff gate (326-364), the Phase 6 no-suppression
  rule and review briefs (469-524), the Phase 7 signoff gate (526-554),
  §11.A round-1 independence (823-836), §11.B sub-branches (887-891), and the
  §4 escalation subsection (654-685).
- Read the complete 1,321-line embedded
  `internal/protocol/defaults/COOPERATION.md`, including §15 in full
  (lines 1163-1322) and the §4.0 qualifier (lines 196-209).
- Read `consensus.md` lines 12-576: the adopted §§15.1-15.6 rule text, the
  per-track binding table, the two text fixes, all four verdict conflicts, the
  drafter position changes table (23 rows), the comparison and blind spots
  section, the tooling record, and the follow-ups list. I did not re-audit
  every preserved revision-1 through revision-4 signoff line below line 576.
- Diffed the two protocol copies: differences are confined to the five
  allowlisted zones (Workspace, Created, Protocol synced, two §2 table bodies).
- Executed `go build ./...` (exit 0), `go test -count=1 ./internal/protocol/...`
  (pass), and the two focused drift-guard tests (pass).

## Refutation attempts

1. Compared each shipped subsection 15.1-15.7 against both `FINAL.md`'s
   decision table (lines 28-35) and `consensus.md`'s blockquoted adopted text,
   looking separately for narrowed triggers, widened evidence classes, added
   obligations, dropped clauses, and changed track cells. Found five silent
   semantic changes in §§15.1-15.2 and one wrong section reference in §15.3.
2. Tried to read the §15.5 "How this rule is actually enforced" paragraph as
   descriptive prose only. It contains a normative `SHOULD` (follow-up 9) and
   implements the compliance-model statement (follow-up 8), both of which
   `FINAL.md:198-201` and `consensus.md:559-576` explicitly label as not
   ratified. It also places idea-specific audit history (8 -> 13 -> 21 -> 23)
   into the embedded `parley init` template.
3. Read §15 through each track's actual gate rather than the full-lifecycle
   defaults. No conflict found with §4.0's reviewer/signoff counts, §5 quorum
   as qualified by §4.0, the Phase 3 signoff gate, or the Phase 7 signoff gate.
4. Followed all four phase pointers (Phases 1, 2, 3, 6). They point to §15 and
   do not suppress findings: §15.3 expressly preserves the review-dispute
   close paths, and §15.4 expressly says it does not gate what a reviewer may
   report.
5. Tested whether a first reader could determine when an untagged sentence
   becomes `RECALL`. The text defines entry triggers but never defines what
   language constitutes a "verification verdict," leaving common forms such as
   "the check passed" and "I verified X" unresolved.
6. Followed the §15.3 reference to "the §8 user-escalation path." §8 is
   "Inbox (lightweight channel)" (line 728); the actual escalation procedure
   is under §4, "### Escalation to user (any phase)" (line 654). §8 itself
   says "escalation -- see §4" (line 738), confirming §8 is the filing
   mechanism, not the procedure.

## Findings

### CRITICAL

No findings.

### MAJOR

#### [MAJOR] §§15.1-15.2 contain five silent semantic changes

Target claim: `IMPLEMENTATION.md:11-14` says the six ratified rules were
implemented, and `IMPLEMENTATION.md:51-67` identifies exactly one deviation.

**Verdict: `WRONG` -- `PRIMARY`.** I compared the ratified passages directly
with the shipped passages in both protocol copies:

1. **Challenge trigger narrowed.** `FINAL.md:30` and `consensus.md:22-24` make
   a challenge itself an entry trigger: "(b) another participant challenges
   it." Shipped `COOPERATION.md:1185-1186` adds "and requests a verdict":
   "(b) another participant challenges it and requests a verdict." A material
   challenge without the extra request now falls outside the regime. That
   narrows the scope of the entire verification system.

2. **"Stable identifier" obligation added.** Shipped `COOPERATION.md:1187-1188`
   adds: "The invoking artifact identifies the claim by a stable identifier or
   an exact quotation." Neither `FINAL.md:30` nor `consensus.md:22-48` contains
   this sentence. It is a new procedural obligation on every participant who
   invokes the regime.

3. **"Exactly one provenance tag" reinstated.** Shipped `COOPERATION.md:1213`
   opens §15.2 with "Every verification verdict carries exactly one provenance
   tag." This sentence is not in the consensus blockquoted adopted text
   (`consensus.md:56-69`), which goes straight to the table. Consensus row 18
   (`consensus.md:1176-1180`) documents that "Every verdict carries exactly
   one provenance tag" was the round-2 position, explicitly relaxed in
   revision 2 to "tag the decisive basis and disclose the rest in prose." The
   shipped text has both, which is confusing: a reader cannot tell whether
   "exactly one" is a separate, stricter requirement or a restatement of the
   single-basis case under "decisive basis." The consensus process decided to
   replace the former with the latter; the implementer reinstated the rejected
   text without disclosure.

4. **PRIMARY evidence widened.** `FINAL.md:31` and `consensus.md:58` define
   executed-check `PRIMARY` evidence as "a check the verifier executed, with
   the command, inputs and relevant output quoted." Shipped
   `COOPERATION.md:1217` widens this to "the command or steps, inputs and
   relevant output quoted." The added "or steps" admits multi-step procedures
   as `PRIMARY` evidence without ratification.

5. **SECONDARY evidence strengthened.** `FINAL.md:31` and `consensus.md:59`
   require "a named other participant's non-`RECALL` verdict." Shipped
   `COOPERATION.md:1218` additionally requires: "name the participant and the
   artifact." This is an unratified strengthening: a `SECONDARY` verdict that
   names the participant but not the artifact is now malformed under the
   shipped text, but admissible under the ratified text.

These five clauses control which claims enter the regime and which evidence is
admissible -- the core of §15. They are not editorial transcription. Restore
the adopted wording in both copies. If stable claim IDs, manual check steps,
or mandatory artifact locators are desirable, ratify them explicitly rather
than shipping them as silent scope changes. The drafter was caught doing
exactly this three times inside the consensus process (`consensus.md:369-427`,
the 23-row drafter position changes table), which makes a fourth occurrence in
the implementation phase a pattern, not an accident.

#### [MAJOR] §15.3 references the wrong section for user escalation

**Verdict on `COOPERATION.md:1247` ("follows the §8 user-escalation path"):
`WRONG` -- `PRIMARY`.**

`consensus.md:95` says "the existing user-escalation path" -- no section
number. The shipped text replaces "existing" with "§8." But §8 is "Inbox
(lightweight channel)" (`COOPERATION.md:728`); it describes the filing
mechanism for inbox notes, not the escalation procedure. The actual escalation
procedure is under §4, "### Escalation to user (any phase)"
(`COOPERATION.md:654-685`), which defines when to escalate, what the inbox
file must contain, and how the user's answer is quoted back into the audit
trail. §8 itself confirms the redirect: its `to-user` example says
"escalation -- see §4" (`COOPERATION.md:738`).

A first reader of §15.3 who needs to escalate a `DISPUTED` claim will follow
the pointer to §8, find the inbox filing instructions, and have to notice the
"see §4" cross-reference to reach the actual procedure. Fix by restoring
"the existing user-escalation path" (no section number), or by changing "§8"
to "§4" if a numbered reference is desired.

#### [MAJOR] Remove the §15.5 deviation and leave follow-ups 8 and 9 open

The disposition is **remove**, not keep.

**Verdict on treating `COOPERATION.md:1289-1294` as ratified scope:
`WRONG` -- `PRIMARY`.**

`FINAL.md:198-201` places both the compliance-model statement and the scope
declaration under `## Follow-ups`, not under the ratified decision table.
`consensus.md:559-576` labels the same items "not in scope here." The shipped
paragraph nevertheless adds the `SHOULD` obligation from follow-up 9
(kimi-1's scope-declaration rule) and the compliance-model text from follow-up
8 (codex-1's finding).

Disclosure in `IMPLEMENTATION.md:53-67` is good audit practice, but it does
not supply ratification. Keeping the paragraph would let Phase 5 implement a
genuine protocol rule that §7 requires to go through a meta-protocol-change
idea. The implementer's arguments for keeping it are reasonable -- deferring it
ships the exact misreading the follow-up exists to prevent -- but the process
answer is the same one this idea's own consensus gave: the follow-ups are
open, not ratified.

There is also a concrete defect: the paragraph carries the source idea's
audit narrative ("the disclosure went 8 -> 13 -> 21 -> 23 of 23 material
changes") into the embedded default
(`internal/protocol/defaults/COOPERATION.md:1280-1285`). Every freshly
initialized deck via `parley init` will contain this specific idea's
compliance history, which is neither generic nor useful to a new project. The
`TestDefaultCooperationForInit` `mustNotContain` list checks for agent-family
names but not for idea-specific metrics, so this passed the guard.

Remove the complete "How this rule is actually enforced" paragraph from both
copies and leave follow-ups 8 and 9 open for their own ratification.

#### [MAJOR] "Verification verdict" is undefined, so the regime boundary is not usable

Section 15.1 says the regime begins when a participant "assigns a verification
verdict," and §15.2 says every such verdict "carries exactly one provenance
tag," but no sentence defines what language performs that assignment. A first
reader cannot tell whether only the reserved statuses `CONFIRMED`, `WRONG`,
and `UNVERIFIED` are verdicts, or whether "the test passed," "I verified X,"
"this proves X," and an `IMPLEMENTATION.md` result cell containing `OK` are
verdicts too.

This ambiguity cuts both ways: a broad reading unexpectedly turns routine
factual reporting into untagged `RECALL`; a narrow reading bypasses the
provenance rule merely by avoiding the reserved words. Add a concrete
boundary sentence to §15.1. It should state either the exact reserved syntax
or that equivalent truth-status language also counts, and should say
explicitly whether raw source or command output reported without a truth
classification is evidence rather than a verdict.

### MINOR

#### [MINOR] The §4.0 qualifier adds a false cross-reference

**Verdict on "see §15.1" at `COOPERATION.md:215`: `WRONG` -- `PRIMARY`.** The
ratified replacement at `FINAL.md:43-44` and `consensus.md:202-203` ends at
"per-agent isolated staging." The shipped text adds "-- see §15.1." Section
15.1 (`COOPERATION.md:1182-1209`) discusses scope, ownership, and location of
verification verdicts; it says nothing about round-1 independence, sub-branches,
or staging. The section that describes sub-branches is §11.B
(`COOPERATION.md:889`). Keep the qualifier -- saying that independence is an
unenforced cooperative convention genuinely reconciles §4.0 with §11.A and
§11.B. Remove only "-- see §15.1" from both copies.

#### [MINOR] §15.5 drops "Binds on every track" from the ratified text

**Verdict: `PRIMARY`.** The consensus blockquoted text
(`consensus.md:140`) includes "Binds on every track." at the end of the first
paragraph. The shipped text (`COOPERATION.md:1276-1279`) omits it. The
information is preserved in the §15.7 binding table (all three columns say
"yes" for §15.5), so no binding is lost. But the ratified sentence was dropped
without disclosure. Restore it, or note that §15.7 carries the binding
information.

### NIT

#### [NIT] §15.2 RECALL definition adds descriptive elaboration

Shipped `COOPERATION.md:1219` says "Memory or unsupported reasoning only; no
source consulted and no check run." `consensus.md:60` says "Memory or
unsupported reasoning only." The added clause is descriptive, not normative --
it clarifies what "memory or unsupported reasoning only" means. Not a semantic
change, but a silent addition.

#### [NIT] §15.4 adds a "same shape" paragraph not in ratified text

Shipped `COOPERATION.md:1268-1270` adds: "The same shape appears as 'this
cannot regress X', 'this path is unreachable', 'that case can't happen',
'this sidesteps the licensing issue' -- the rule is about the claim's form,
not its subject matter." This paragraph is not in the consensus blockquoted
text (`consensus.md:120-129`). It is helpful for a first reader and adds no
obligation, but it is an unratified addition.

#### [NIT] §15.6 drops the rationale from clause (b)

`consensus.md:181` says "This clause binds unchanged on both tracks, since
`standard` has a separate `consensus.md`." Shipped `COOPERATION.md:1316` says
"This clause binds unchanged on both tracks." -- the rationale "since
`standard` has a separate `consensus.md`" was dropped. The normative content
is unchanged; only the explanatory note is lost.

#### [NIT] IMPLEMENTATION.md attributes the four pointers to the wrong authority

**Verdict on `IMPLEMENTATION.md:48` ("`FINAL.md` specified pointers"):
`WRONG` -- `PRIMARY`.** A direct search of `FINAL.md` finds no phase-pointer
instruction. `consensus.md:17-18` is the adopted text: "one-line pointers
from Phases 1, 2, 3 and 6." The four shipped pointers are correctly placed and
accurately describe §15. Change "`FINAL.md` specified" to "`consensus.md`
specified" in `IMPLEMENTATION.md`.

## Checks with no finding

- **Verdict on `IMPLEMENTATION.md:25-26` (the copies remain identical outside
  five allowlisted zones): `CONFIRMED` -- `PRIMARY`.** Executed:

      go test -count=1 ./internal/protocol -run '^(TestEmbeddedDefaultMatchesLiveDeck|TestDefaultCooperationForInit)$' -v

  Relevant output:

      --- PASS: TestEmbeddedDefaultMatchesLiveDeck (0.00s)
      --- PASS: TestDefaultCooperationForInit (0.00s)
      PASS

  `drift_test.go:98-129` normalizes only Workspace, Created, Protocol-synced,
  and the two §2 table bodies. `drift_test.go:146-178` separately asserts the
  generic bootstrap shape and empty rosters. The guard was unchanged by
  `bfca39e`. A direct `diff` of the two files confirms differences are confined
  to the five allowlisted zones.

- **Verdict on the shipped §15.7 binding table's fidelity: `CONFIRMED` --
  `PRIMARY`.** I compared `COOPERATION.md:1322-1330` cell-for-cell with
  `FINAL.md:28-35` and `consensus.md:221-231`: 15.1-15.5 bind on all tracks,
  and 15.6 is off on `fast`, an existing round-02 section on `standard`, and
  an assigned round artifact on `deliberation`. All three sources agree.

- **Verdict on `IMPLEMENTATION.md:79` (no agent-family names in bootstrap §15):
  `CONFIRMED` -- `PRIMARY`.** A search of the embedded §15
  (`internal/protocol/defaults/COOPERATION.md` lines 1163-1322) for
  `claude`, `codex`, `hermes`, `kimi`, `feci`, `agy`, `glm`, `gpt`, and `opus`
  produced no output. After the §15.5 deviation is removed, the section also
  contains no ratifying-run metrics.

- **Verdict on internal consistency with §4.0, §5, Phase 3, Phase 7, and the
  Phase 6 no-suppression rule: `CONFIRMED` -- `PRIMARY`.** §15's intro
  paragraph (`COOPERATION.md:1179-1180`) states it composes with §4.0's
  per-track table and the Phase 6 no-suppression rule. §15.3
  (`COOPERATION.md:1254-1255`) preserves the review-dispute close paths.
  §15.4 (`COOPERATION.md:1272`) says it does not gate what a reviewer may
  report. No contradiction found with §5 quorum (as qualified by §4.0's
  per-track table), the Phase 3 signoff gate, or the Phase 7 signoff gate.

- **Verdict on the §4.0 qualifier resolving the round-1 independence
  contradiction: `CONFIRMED` -- `PRIMARY`.** §4.0 previously listed round-1
  independence among invariants "never dropped for speed" while §11.A said
  "There is no enforcement beyond agent discipline" (`COOPERATION.md:835`).
  The shipped qualifier (`COOPERATION.md:213-215`) now says "an unenforced
  cooperative convention unless kickoff selects §11.B sub-branches or
  per-agent isolated staging." This genuinely reconciles the two: §11.B
  sub-branches (`COOPERATION.md:889`) do provide stronger isolation. The
  contradiction is resolved (modulo the false "-- see §15.1" pointer, filed
  as MINOR above).

- `go build ./...` exited 0. `go test -count=1 ./internal/protocol/...`
  returned `ok parley-deck-cli/internal/protocol`.

## Open questions

None. The necessary fixes are textual and do not require new tooling or
ownership changes.
