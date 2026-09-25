---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
review-round: 3
date: 2026-09-25
reviewed-commit: 1bad2634380dd785009987a5c50427050efb87c3 (CLI source, reviewed baseline d2382388415f8dc120a71ddbf6077d395dd18b8e) + b850576ba3f943ba9e4c942247e864982d5ccfdc (canonical record — IMPLEMENTATION.md only, no Go content; on-disk file byte-identical, sha256 90e1e01b4318b9af6f0c95a9459c7c9a5a4d206c4d5e7c14e5cd2931000f2a8e) + a624318dcda02c47ecaa859d987efee08dba3104 (skill, unchanged this cycle; baseline 8161e5e; skill worktree HEAD re-verified at the full SHA)
---

## Summary

Round-03 re-review of the Phase-8 fix-up cycle 2 delta per the signed cycle-2 plan's
verification-plan step 3. **All four agreed fixes (AF-12, AF-13, AF-14, AF-15) are applied
exactly as signed, including both carried corrections (claude-1's Reservation 1 `:42`
locator fix; elected alternative (a) for AF-15), on my own primary evidence: an isolated
local-disk checkout detached at the exact full SHA, the live phase-8 packet, the live
driver digest, mutation checks M1/M2 re-run by me, and my own novel adversarial regex
probes beyond the filed cases. I filed no new finding — a scoped null, with the scope
named honestly below. The cycle-2 delta is, in my verdict, ready for the zero-fix
consensus step; nothing in it blocks closure, and no duplicate broad suite is owed for
this delta (no new concern arose; the ratified proportional plan governs). This review
does not close the idea, signs no consensus, and waives no coverage. The LE-7
fresh-invocation goal-done check remains separately owed after a zero-fix consensus.

## Protocol context attestation (Phase-8 re-review)

`parley protocol packet --dir . --phase 8 --track deliberation --idea
meta-protocol-change-designated-implementer --flag auto_implement --flag protocol_change
--audience participant --json`, run by me this session in the original CLI workspace (the
designated-implementer worktree), parley 1.49.1, exit 0:

```json
{"context_mode": "full", "source_sha256": "b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388", "packet_sha256": "b273af1e0649a365bf384083d27c319ddb9cb62e0efe1b757eb5e0ccfc22f388", "source": {"role": "source", "transport": "github-pr", "bytes": 115166}, "shadow": {"packet_sha256": "fe0e4c046897ff5d177edd64be71abff44b03c0d6dcb166492bc6a485c9bebcf", "packet_bytes": 86716, "source_bytes": 115166, "included_blocks": 40, "omitted_blocks": 29}}
```

- `fallback_reason` ABSENT (top-level keys enumerated: `body_path, context_mode, index,
  packet_sha256, request, shadow, source, source_sha256`).
- The source authority `b273af1e…f388` (1,400 lines / 115,166 bytes) is UNCHANGED from
  both cycle-2 signoff attestations and equals `git show d238238:parley-deck/COOPERATION.md`
  — cycle 2 moved no protocol text, as the plan requires. Phase 6/7/8, §15.1/§15.2/§15.3,
  §0/§9.0 and the Phase-8 frontmatter-bump sentence were read from the full body; the bump
  sentence re-verified verbatim at body `:622` ("They also update the top-level frontmatter:
  bump `status:` to `fix-up-cycle-N`, update `head-commit:`.").
- `strict_gate` is not set in `00-prompt.md` (frontmatter re-read: `facilitator: codex-1`,
  `participants: [claude-1, kimi-1, zcode-1]`, no `strict_gate` key) — the default Phase-8
  close rule applies.
- **My own phase-8 shadow measurement reproduces kimi-1's, exactly**: 86,716 bytes / 40
  included / 29 omitted, and the identical shadow packet hash `fe0e4c04…ebcf` — see the
  shadow-figure assessment under AF-14.

## Method and provenance

Unique local-disk isolated checkout — `/private/tmp/zcode1-r03-cli`, `git worktree --detach`
at `1bad2634380dd785009987a5c50427050efb87c3`, `git status --porcelain` empty, on local disk
(`/tmp`), never the shared mount — removed after the run. No peer fixtures, run caches or
organizer notes were used as evidence; organizer notes were orientation only. All evidence
below is my own PRIMARY unless another participant is named. Per the ratified proportional
verification plan I ran **no broad suite** (see AC-17 standing below); my executions were:
build/vet/gofmt, the full `internal/protocol` package, three named `internal/app` tests,
mutation checks M1/M2 (applied, run, restored, re-verified green), one scratch regex module
outside the repository (deleted after the run), the phase-8 packet above, and the read-only
driver digest. Nothing was committed, tagged, merged, installed, published or released by
me; no global default was mutated; the product default stays UNSET.

## The four agreed fixes, reassessed

### AF-12 — frontmatter/status/current PRIOR source SHA + cycle-1 metadata repair — APPLIED AS SIGNED, machine effect verified live

- Frontmatter at `b850576` (disk byte-identical to it, sha256 `90e1e01b…2a8e`):
  `status: fix-up-cycle-2`, `head-commit: 1bad263` (the cycle-2 source commit — PRIOR and
  knowable under the source-then-record order), `skill-commit: a624318`. Short-SHA form is
  the file's established convention (`0893989` before it) — noted, not a defect.
- The cycle-1 repair block inside `## Fix-up cycle 1` carries `head-commit: d238238`,
  `skill-commit: a624318`, `record-commit: e4d868a` — all PRIOR commits at repair time, no
  self-naming — plus the explicit admission that the per-publication bump owed at `e4d868a`
  (Phase 8, body :622) was missed there and why (per-cycle close conflated with idea close).
  Exactly the signed three-part shape, including the absorption of my round-02 F-2.
- `## Fix-up cycle 2` names the source commit `1bad263` and states its own record commit as
  following ("named at the next touch") — no self-hash claim, exactly as signed.
- **Machine effect, verified live by me** (the axis the MAJOR framing stood on):
  `parley wait --dir . --idea … --for review --timeout 20s -json` →
  `implementation.status = "fix-up-cycle-2"`, `head_commit = "1bad263"`, `ready_for_review:
  true`, and `next: "await review artifact"` — the digest no longer collapses to
  `await implementation` (claude-1's round-02 live symptom). Exit 3 with
  "outstanding: claude-1 (review artifact), zcode-1 (review artifact)" is the expected
  pre-filing timeout for this very round — the content signal is functioning and now
  correctly awaits the two round-03 re-review artifacts (N=2, M=0 at my run).
- **Commit-prefix provenance, verified rather than trusted:** `git rev-parse bf3336d^{tree}`
  and `git rev-parse 1bad263^{tree}` both return `452ccde32460d78e12947a20cff2329a4a351be7`
  — the superseded `bf3336d` and the corrected `1bad263` are tree-byte-identical, so the
  message-only-rewrite claim holds and every retained test-provenance statement in the
  record is true at the content hash it names. The owner's standing prefix override
  (`parley-deck/inbox/codex-1-to-kimi-1_…_commit-prefix.md`, committed by the organizer in
  `c0fbc2c`) says exactly what the record reports: correct the unpushed message to
  `[codex-1] … (authored by kimi-1)`, do not change content, do not claim old evidence at a
  new hash. Both current subjects carry the `[codex-1]` prefix with kimi-1 authorship
  retained. Accurate and complete; no undisclosed rewrite.

### AF-13 — accurate role-action comments + the `:42` different-struct correction — APPLIED AS SIGNED

- `internal/app/driver_impl.go:148` now reads `gate string // validity / pin-conflict
  gate — roleErr path, escalated by every role action (R16/R25/R28)`; the `:188` block now
  reads "Hard gate on any run that reaches a role action". `grep -n "any run"
  internal/app/driver_impl.go` returns exactly one hit — `:188`, qualified — the
  unqualified claim is gone.
- The carried correction is accurately recorded: `:42` IS `roleErr string //
  declared-facilitator role deadlock; every role action escalates` inside `driverImplOps`
  (re-read by me in the checkout); `:45` is the "All five stay zero on an undesignated
  deck" comment; the two structs are distinct and ~106 lines apart. The record's narrowing
  of the plan's wrong DRAFTER-PRIMARY locator clause — without editing the signed plan or
  any signature — is the §15.2-honest way to carry it.
- Comments only: `go build ./...` exit 0, `go vet ./...` exit 0, `gofmt -l` clean on all
  three touched files; no test changes for AF-13.

### AF-14 — honest pre/post authority labels, phase-labelled shadow figures, pin-source live() corner — APPLIED AS SIGNED

- (1) The cycle-1 Phase-8 attestation bullet is now explicitly PRE-EDIT at `0893989`
  (`c7492182…568f`, 114,771 B, phase-8 shadow 86,336 B), with the POST-edit authority at
  `d238238` (`b273af1e…f388`, 115,166 B) named; **every** shadow figure now carries its
  phase (phase-8 over `d238238`: 86,716 B / 40-29; phase-7 over the same source: 80,799 B /
  41-28), closing claude-1's wording caution with "Read the phase with the figure."
- (2) The AF-11 Surprises & Discoveries sentence now carries the pin-source corner
  parenthetical (after AF-4, `live()` is also true for a pin-source dispatch with tier 3
  set — `present: true`, `source: pin` — so the kickoff `agent.model_diversity` emission
  can fire there too; doubling confined to `live()`-true dispatches). My round-02 NIT-3,
  carried as dismissed-and-recorded. Accurate.
- **The 86,701 → 86,716 shadow figure, assessed openly (no cause assumed):** the signed
  plan's AF-14(1) named an 86,701-byte phase-8 shadow over `d238238` (claude-1's round-02
  measurement, same 40/29 block split). The record reports 86,716 B measured twice under
  parley 1.49.1 and explicitly says the plan's figure "did not reproduce under this binary
  and is superseded by this measurement". My own independent run this session returns
  **86,716 B / 40 included / 29 omitted with the byte-identical shadow packet hash
  `fe0e4c04…ebcf`** — the figure is deterministic under this binary and the record's number
  is the reproducible one. The 15-byte delta against claude-1's round-02 figure has the
  same block split and the same source hash on both sides, so it is not block selection or
  source drift; its cause is not determinable from the evidence available to me (claude-1's
  rendering binary/environment at round-02 was not recorded), and I do not guess one. The
  record's handling — disclose the discrepancy, keep the reproducible figure, keep the
  source hash/line/byte counts that DO reproduce — is exactly the honest treatment; the
  unexplained residual is a property of the measurement history, not a defect in the record.

### AF-15 — enumerated hyphenated/spaced negations including `un[\s-]*confirmed` option (a) — APPLIED AS SIGNED AND ELECTED

- `internal/protocol/implementer.go:185`: `negationMarker` is now
  `(?i)\b(?:un[\s-]*confirmed|not(?:[\s-]+yet)?[\s-]+confirmed)\b` — precisely alternative
  (a) as claude-1 offered and I elected in my consensus signoff: the bare `unconfirmed`
  alternative widened to `un[\s-]*confirmed` (zero-or-more separators, so the enumerated
  verbatim spelling still matches), `not`-branch inner separators `[\s-]+`. No new banned
  word; the completed enumeration is the signed AF-1 clause (iv) one; the comment above the
  regex says exactly that.
- Tests: 14 new cases across BOTH parsers in `TestConfirmationRecordsAreFailClosed` —
  rejects `not-yet-confirmed`, `not-confirmed`, `un-confirmed`, `un confirmed` (both
  parsers) plus cased `UN-Confirmed` (waiver); boundary accepts `cannot-confirmed`,
  `unconfirmedness`, `run confirmed the fix`, em-dash `reason—detail` (NIT-2(a) tolerance
  untouched). The diff to the test file is 19 insertions, 0 deletions — the existing 15
  cases and the T-4/T-5 fixtures are unmodified, as signed.
- **My executions (exit statuses preserved):** `go test ./internal/protocol/ -count=1` →
  `ok`, exit 0; the fail-closed test yields 29 subtests PASS / 0 FAIL (15 prior + 14 new;
  30 PASS lines with the parent) — matching the record's count exactly;
  `go test ./internal/app/ -run 'TestTier2UnavailabilityGateAndExits|TestPinDesignationConflictEscalates|TestUnsetPathIsByteIdentical'
  -count=1` → all three PASS, exit 0 (T-4/T-5 fixtures and the AC-4 byte-identical pin,
  unmodified).
- **Mutation checks, re-run by me (applied → run → restored → re-verified green):**
  **M1** — full separator-class revert to the pre-AF-15 regex → exit 1 with EXACTLY the 9
  new reject subtests failing (5 waiver incl. `UN-Confirmed`, 4 reassignment); **M2** —
  un-branch-only revert (the signed plan's literal shape without correction (a)) → exit 1
  with EXACTLY the 5 (a)-specific subtests failing (3 waiver, 2 reassignment). No boundary
  accept flipped under either mutation. Both of the record's mutation claims reproduce
  exactly; the new tests are mutation-meaningful, not decorative.
- **My own novel probes (scratch module outside the repo, deleted after):** 17 cases beyond
  the filed set — enumerated bases (`unconfirmed`, `UNCONFIRMED`, `not confirmed`,
  `Not  Yet   Confirmed`), mixed separators (`un  -confirmed`, `un-  confirmed`,
  `not - confirmed`, `NOT-YET-confirmed`), and false-positive boundaries (`sun-confirmed`,
  `fun confirmed`, `cannot be confirmed yet`, `not-the-point`/`confirmed-ok` prose,
  bare `confirmed` as a reason, `unconfirm`, em-dash). All behave correctly. One probe
  self-correction, recorded per §15.1: I expected `"confirmed-unconfirmed"` not to match
  and it does — my expectation was wrong, not the code: the trailing `unconfirmed` token is
  matched by the OLD regex too (verified side by side; `\b` holds after the hyphen), so
  this is pre-existing fail-closed behaviour, unchanged by this delta, in the fail-closed
  direction.

## Refutation attempts

- **R-A — "the regex widening must have broken a legal record somewhere."** Refuted three
  ways: the full `internal/protocol` package is green at the exact commit (exit 0); the
  three `internal/app` fixture tests that flow real designation records through the parsers
  are green unmodified; my 17 novel probes found zero false positives (`\b` blocks
  `sun-confirmed`/`fun confirmed`; the trailing `\b` blocks `unconfirmedness`; em-dash is
  outside `[\s-]`). The only behavioural match my probes found beyond the filed set
  (`confirmed-unconfirmed`) is pre-existing under the old regex and correctly fail-closed.
- **R-B — "the tests could be passing vacuously."** Refuted by M1/M2: reverting the
  separator classes fails exactly the 9 new rejects; reverting only the `un`-branch fails
  exactly the 5 (a)-specific rejects. Every new reject case is pinned to the exact widening
  that produces it; boundary accepts flip under neither mutation.
- **R-C — "the prefix rewrite (`bf3336d` → `1bad263`) could hide a content change, or the
  retained provenance could be stale."** Refuted: both commits resolve to the identical
  tree hash `452ccde3…1be7` (PRIMARY, this session); the owner's override instruction and
  the record's provenance note match each other and the git state; the record claims no
  test evidence at any hash other than the ones it names.
- **R-D — "AF-12's bump could be cosmetic — the machine may still misread."** Refuted live:
  the driver digest reports `status=fix-up-cycle-2`, `head_commit=1bad263`, and
  `next: await review artifact` (not `await implementation`), with this round's two
  artifacts correctly listed as outstanding. The content signal works at the boundary
  where it is now in use.
- **R-E — "some unrelated behaviour/roster/model/claim/default/release scope could have
  moved under cover of the record commit."** Refuted by confinement, derived not assumed:
  `d238238..1bad263` touches exactly three Go files (`driver_impl.go` 2 comment lines;
  `implementer.go` regex + comment; `implementer_test.go` 19 added lines, 0 deletions);
  `b850576` touches only `IMPLEMENTATION.md`; the organizer-only commits on the branch
  (`f5a3908`, `2ed601d`, `8d026d4`, `c0fbc2c` — notes/ledger/review-preservation/inbox/
  release-inventory) touch no code and no frozen artifact: `git diff d238238..HEAD --
  internal/ cmd/` is the three reviewed files and nothing else; `FINAL.md` and
  `00-prompt.md` are untouched (`git diff d238238..HEAD` empty for both); the FINAL owner
  boundary sentence "The global product default ships UNSET" stands at `:492`; the skill
  worktree HEAD is still `a624318dcda02c47ecaa859d987efee08dba3104`.
- **R-F — "the cycle-2 pre-edit packet attestation could have been taken against a tree
  that was not content-identical to `d238238`."** Refuted: `git diff d238238..8d026d4 --
  internal/ cmd/` is empty — the record's statement that the packet was rendered while the
  code tree remained content-identical to the reviewed `d238238` is true, and my own
  packet today renders the same source hash over the same bytes.

## Findings

**None — a scoped null.** Within the scope I was assigned (the cycle-2 delta: AF-12 record
metadata and its machine effect; AF-13 comment accuracy; AF-14 labels and figures; AF-15
regex scope and tests; the unmodified default-path pins; confinement of the delta; the
operational provenance of the two commits), I found no defect of any severity: nothing
incorrect, nothing undisclosed, no scope creep, no suppressed reservation. Two
below-finding observations, recorded for traceability only: (1) the frontmatter
`head-commit` uses the 7-char short SHA, consistent with the file's own prior convention
and with what the digest echoes; (2) the 15-byte phase-8 shadow discrepancy against
claude-1's round-02 figure remains unexplained on the evidence available (see AF-14) —
the record discloses it and keeps the reproducible figure, which is all the honesty rule
can ask of it.

## AC-17 whole-suite standing — not re-claimed, not re-run, per the ratified plan

Claude-1's independent full suite at `d238238` (durable log
`/tmp/parley-claude1-r2-1790307993/full-suite.log`, which I read and corroborated in my
consensus signoff) stands **exact-commit scoped to `d238238`**. No broad suite has run at
`1bad263` and none is owed for this delta absent a new concern — I raised none, and the
proportionality premise re-derived in my signoff (the only behaviour change is AF-15's
regex line in `internal/protocol`; the `internal/app` fixtures carry no hyphen-adjacent
reason segments) was re-confirmed by this session's diff read. The record claims the same
scoping in the same words. The LE-7 fresh-invocation goal-done check remains separate and
owed.

## Deferrals — re-verified and concurred

DF-1 (drafter-precheck eligibility divergence → FINAL register F6, a later idea), DF-2
(launch-surfacing → the NAMED, INACTIVE slug
`meta-protocol-change-designation-launch-surfacing` — re-verified absent from the tree by
my own grep this session; opening it remains a post-close owner/organizer act), DF-3
(durable tier-3 fall-through notice surfacing → organizer release/done report and the
owner's post-release configuration step; no participant action), DF-4 (Windows residual →
owner-deferred to `windows-portability`; this delta adds no platform surface — pure
comments, one regex and tests are platform-neutral). **I concur with all four exactly as
carried:** each is outside this idea's ratified scope, has a named carrier, was re-confirmed
by both round-02 reviewers, and none is widened or quietly closed by cycle 2.

## Limitations (named, not waived)

No broad `go test ./...` was run by me at any commit (ratified proportional plan; single
focused-package plus named-test evidence instead). No npm-level skill gates, no end-to-end
`parley run` against a real designated deck, no live §9.0 ping behind `designeeAvailable`,
no TUI/pipeline-block dispatch surfaces, no Windows — the cycle-2 blind-spot set is the
cycle-1 set, unchanged, and none is waived. `b850576` contains no Go content, so no suite
result extends to it; I reviewed it by reading (which is how AF-12/AF-14 were verified).
The eleven cycle-1 fixes were not re-derived from scratch this round: they were
independently re-derived by both reviewers at round-02 and this delta provably does not
touch them (confinement by diff, mutation checks, green focused suites); their standing
evidence is the round-02 record plus claude-1's full suite at `d238238`. The shadow-figure
discrepancy's cause is indeterminate on available evidence. I did not read claude-1's
round-03 artifact (not filed at my run time; I was instructed not to wait on it) and this
review is independent of it.

## Disposition concurrences and verdict

I concur with every disposition of the signed cycle-2 consensus as executed: VC-2.1's
severity resolution, VC-2.2's exact-commit reconciliation, both carried corrections
(Reservation 1 as self-corrected by claude-1; Reservation 2 resolved by elected
alternative (a), which the source now carries verbatim), the four dismissals, DF-1…DF-4,
and the proportional verification plan as applied. The implementer's record is accurate on
every point I checked, including the two places where honesty was the hard choice (the
cycle-1 bump admission; the non-reproducing plan figure superseded openly).

**Verdict: PASS — zero new agreed fixes from zcode-1 for this delta.** From my side the
cycle-2 output is ready for the cycle-3 zero-fix consensus step, after which the LE-7
fresh-invocation goal-done check by a fresh non-implementer invocation of an existing
quorum member remains owed, then the owner/organizer-only release sequence. This verdict
is mine as reviewer; it signs no consensus, closes no idea, and authorizes no release act.
