---
agent: codex-1
idea: meta-protocol-change-verification-integrity
round: 1
phase: review
date: 2026-08-04
---

## Summary

Changes requested: three `MAJOR` findings, one `MINOR`, and one `NIT`. There are no `CRITICAL`
findings. The drift guard is intact, both protocol copies pass its byte-for-byte comparison outside
the five allowlisted zones, and the per-track table agrees with the ratified table. The blocking
work is textual: restore four silently changed clauses, remove the disclosed-but-unratified §15.5
follow-up paragraph, and make the verdict boundary operationally clear.

## Scope actually checked (§15.5)

- Read all of `FINAL.md`, `IMPLEMENTATION.md`, `internal/protocol/drift_test.go`, and the complete
  1,330-line live `parley-deck/COOPERATION.md`.
- Read `consensus.md` lines 12-576: the adopted §§15.1-15.6 text, the 15.7 binding table, the two
  text fixes, conflict resolutions, comparison/blind-spots analysis, and follow-ups. I used locator
  searches in the preserved signoff history but did not re-audit every preserved revision-1 to
  revision-4 signoff line.
- Inspected the embedded template's bootstrap header, empty roster tables, four phase pointers,
  two text fixes, and complete §15. For the unchanged remainder, I relied on the drift guard's
  byte-for-byte normalized comparison rather than manually reading the duplicate a second time.
- Compared read-only archives of `bfca39e^` and `bfca39e`. The archive comparison reported only
  the two protocol copies as changed and `IMPLEMENTATION.md` as added.
- Checked the requested consistency surfaces: §4.0's table and invariant qualifier, §5 quorum,
  Phase 3 and Phase 7 signoff gates, the Phase 6 no-suppression rule, and §11.A/§11.B round-1
  independence mechanics.
- Executed `go build ./...`, the two focused protocol tests, `go test -count=1
  ./internal/protocol/...`, and `go test -count=1 ./...`.

## Refutation attempts

1. Compared every shipped subsection 15.1-15.7 against both `FINAL.md`'s decision table and
   `consensus.md`'s adopted rule text, looking separately for narrower triggers, broader evidence
   classes, stronger obligations, changed locations, and changed track cells.
2. Tried to justify the §15.5 deviation as explanatory prose only. It still contains a normative
   `SHOULD`, implements both follow-ups 8 and 9, and copies idea-specific audit history into every
   freshly initialized deck.
3. Read §15 through each track's actual gate rather than through the full-lifecycle defaults. I
   found no conflict with the §4.0 reviewer/signoff counts, §5 as qualified by §4.0, or the Phase 3
   and Phase 7 gates.
4. Followed all four phase pointers. They do not suppress findings: §15.3 expressly preserves the
   existing review-dispute close paths, and §15.4 expressly says it does not gate what a reviewer
   may report.
5. Tested whether a first reader could determine when an untagged sentence becomes `RECALL`. The
   text defines entry triggers but never defines the phrase “verification verdict,” leaving common
   forms such as “the check passed” and “this proves X” unresolved.

## Findings

### CRITICAL

No findings.

### MAJOR

#### [MAJOR] §§15.1-15.2 contain four silent semantic changes

Target claim: `IMPLEMENTATION.md:11-14` says the six ratified rules were implemented, and
`IMPLEMENTATION.md:51-67` identifies exactly one deviation.

**Verdict: `WRONG` — `PRIMARY`.** I compared the ratified passages directly with the shipped
passages:

- `FINAL.md:30` and `consensus.md:22-24` make a challenge itself an entry trigger. Shipped
  `COOPERATION.md:1184-1188` requires that the challenger also “requests a verdict.” That silently
  narrows the regime and permits a material challenge without the extra request to fall outside it.
- Neither adopted §15.1 text contains a stable-identifier/exact-quotation obligation. Shipped
  `COOPERATION.md:1187-1188` adds one. It may be useful, but it is another unratified obligation.
- `FINAL.md:31` and `consensus.md:58` define executed-check `PRIMARY` evidence with the command,
  inputs, and relevant output. Shipped `COOPERATION.md:1217` widens this to “command or steps.”
- `FINAL.md:31` and `consensus.md:59` require a named participant for `SECONDARY`. Shipped
  `COOPERATION.md:1218` additionally requires the artifact, while the malformed-tag rule at
  `COOPERATION.md:1221-1223` still mentions only a missing named dependency. This is both an
  unratified strengthening and an internal uncertainty about when `SECONDARY` fails closed.

These clauses control which claims enter the regime and which evidence is admissible, so they are
not editorial transcription. Restore the adopted wording in both copies. If stable claim IDs,
manual check steps, or mandatory artifact locators are desirable, ratify them explicitly rather
than shipping them as silent scope changes.

#### [MAJOR] Remove the §15.5 deviation and leave follow-ups 8 and 9 open

The disposition is **remove**, not keep.

**Verdict on treating `COOPERATION.md:1289-1294` as ratified scope: `WRONG` — `PRIMARY`.**
`FINAL.md:198-201` places both the compliance-model statement and scope declaration under
`## Follow-ups`; `consensus.md:559-576` labels the same items “not in scope here.” The shipped
paragraph nevertheless adds the `SHOULD` obligation from follow-up 9 and the compliance-model text
from follow-up 8.

Disclosure in `IMPLEMENTATION.md` is good audit practice, but it does not supply ratification.
Keeping the paragraph would let Phase 5 implement a genuine protocol rule that §7 requires to go
through protocol-change consensus. It also puts the source idea's 8 → 13 → 21 → 23 audit narrative
into the generic `parley init` template. Remove the complete “How this rule is actually enforced”
paragraph from both copies and leave follow-ups 8 and 9 open for their own ratification.

#### [MAJOR] “Verification verdict” is undefined, so the provenance boundary is not usable

Section 15.1 says the regime begins when a participant “assigns a verification verdict,” and §15.2
requires exactly one tag on every such verdict, but no sentence defines what language performs that
assignment. A first reader cannot tell whether only the reserved statuses `CONFIRMED`, `WRONG`, and
`UNVERIFIED` are verdicts, or whether “the test passed,” “I verified X,” “this proves X,” and an
`IMPLEMENTATION.md` result cell containing `OK` are verdicts too.

This ambiguity cuts both ways: a broad reading unexpectedly turns routine factual reporting into
untagged `RECALL`; a narrow reading bypasses the provenance rule merely by avoiding the reserved
words. Add a concrete boundary sentence to §15.1. It should state either the exact reserved syntax
or that equivalent truth-status language also counts, and should say explicitly whether raw source
or command output reported without a truth classification is evidence rather than a verdict.

### MINOR

#### [MINOR] The §4.0 qualifier resolves the contradiction but ends with a false pointer

**Verdict on the relevance of “see §15.1” at `COOPERATION.md:213-215`: `WRONG` — `PRIMARY`.**
The ratified replacement at `FINAL.md:43-44` and `consensus.md:202-203` ends at “per-agent isolated
staging.” Section 15.1 contains no isolation, sub-branch, or staging rule; §11.B is the section that
describes sub-branches.

Keep the qualifier: saying that independence is an unenforced cooperative convention genuinely
reconciles §4.0 with §11.A and §11.B. Remove only “— see §15.1” from both copies.

### NIT

#### [NIT] `IMPLEMENTATION.md` attributes the four pointers to the wrong authority

**Verdict on `IMPLEMENTATION.md:48` (“`FINAL.md` specified pointers”): `WRONG` — `PRIMARY`.** A
direct search finds no phase-pointer instruction in `FINAL.md`; `consensus.md:17-18` is the adopted
text that specifies one-line pointers from Phases 1, 2, 3, and 6. The four shipped pointers are
correctly placed and accurately describe §15. Change “`FINAL.md` specified” to “`consensus.md`
specified” in `IMPLEMENTATION.md`.

## Checks with no finding

- **Verdict on `IMPLEMENTATION.md:25-26` (the copies remain identical outside five allowlisted
  zones): `CONFIRMED` — `PRIMARY`.** Executed:

      go test -count=1 ./internal/protocol -run '^(TestEmbeddedDefaultMatchesLiveDeck|TestDefaultCooperationForInit)$' -v

  Relevant output:

      --- PASS: TestEmbeddedDefaultMatchesLiveDeck (0.00s)
      --- PASS: TestDefaultCooperationForInit (0.00s)
      PASS

  `drift_test.go:98-129` normalizes only Workspace, Created, Protocol-synced, and the two §2 table
  bodies; `drift_test.go:146-178` separately asserts the generic bootstrap shape and empty rosters.
  The guard was unchanged by `bfca39e`.
- **Verdict on the shipped 15.7 table's fidelity: `CONFIRMED` — `PRIMARY`.** I compared
  `COOPERATION.md:1322-1330` cell-for-cell with `FINAL.md:28-35` and
  `consensus.md:221-231`: 15.1-15.5 bind on all tracks, and 15.6 is off on `fast`, an existing
  round-02 section on `standard`, and an assigned round artifact on `deliberation`.
- **Verdict on `IMPLEMENTATION.md:79` (no agent-family names in bootstrap §15): `CONFIRMED` —
  `PRIMARY`.** Searches of the embedded §15 for the roster families and common vendor/model names,
  and for Workspace/Transport/Protocol-synced headers, produced no output. After the §15.5
  deviation is removed, the section also contains no ratifying-run metrics beyond the normal
  protocol-provenance citation.
- `go build ./...` exited 0, and `go test -count=1 ./internal/protocol/...` returned
  `ok parley-deck-cli/internal/protocol`.
- The full survival run was not green in this environment: `go test -count=1 ./...` failed only at
  `internal/runner.TestDurableKillEndToEndRealProcess` with “process verification failed (no
  recorded boot id); not killed.” The same focused test fails identically in the read-only
  `bfca39e^` archive, so I do not attribute it to this protocol-only commit and do not file it as a
  finding here.

## Open questions

None. The necessary fixes are textual and do not require new tooling or ownership changes.
