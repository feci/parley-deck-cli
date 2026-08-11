---
agent: hermes-1
idea: protocol-read-cost-regression
round: 1
date: 2026-08-10
---

# Review — frontier context selection

## Summary

The implementation wires `frontierContext` into both quadratic paths (design at
`runner.go:933`, review at `phase58.go:494`), keeps the original walkers as the
fallback, leaves the consensus drafter untouched, and rewrites the design-side
instruction to state the owner-disposes rule. `go build ./...`, `go vet`, and the
six frontier tests are green. The consensus drafter (`driver_consensus.go`) diff is
empty — G4 holds at the source level.

The architecture is sound. The two deviations the implementer flagged (derived
ledger, G6 unmet) are real, and I judge them below. But I found three defects the
implementer did not flag: a context-doubling bug on review rounds 1–2 that violates
G1, a marker-list gap that silently drops live objections (the exact Phase 2
"silence = implicit agreement" pathology), and a fallback-completeness shortfall
where most of the uncertainty states my signoff required have no corresponding
code path. The first two should block ship; the third is the question I was asked
to judge hardest.

## Findings

### [CRITICAL] Marker list misses real objection wording — live objections silently vanish

The marker list in `frontier.go:54-74` is a fixed set of case-sensitive substrings.
I checked it against actual round files in this repository. The following real
objection phrasings, drawn from real `round-02/*.md` files in this deck, match
ZERO markers and are silently dropped at the compaction boundary:

- "I reject your on-demand `verdicts.md`" (`addon-manifest-coverage/round-02/codex-1.md`)
- "I reject Claude's vendored registry fallback" (same file)
- "I reject the suggestion that the build wins" (same file)
- "I reject a distinct P0-P6 artifact ladder" (same file)
- "COUNTER on the Critical/Advisory split" (`meta-protocol-change-review-gate-honesty/round-02/claude.md`)
- "COUNTER on immutability" (same file)
- "I now disagree — the core's false green is broader" (`addon-manifest-coverage/round-02/codex-1.md`)
- "I abandon my previous query" (same file)
- "I object to this approach" (common phrasing, not in markers)
- "Objection: the lock is unchecked" (common phrasing, not in markers)
- "CONCERN: the manifest is stale" (common phrasing, not in markers)

The needle "I disagree" does NOT catch "I now disagree" — the substring "I disagree"
is not contiguous in "I now disagree" because "now" sits between "I" and "disagree".

Why this is CRITICAL, not MAJOR: Phase 2 rule 1 is "Silence = implicit agreement."
The protocol converts an omission into consent. If an objection uses any of the
missed phrasings and its author does not restate it in the next round, the objection
is in NO part of the round-N>=3 context — not in the ledger (marker missed it), not
in round N-1 in full (author didn't restate), not in rounds 1..N-2 (compacted away).
The protocol does not record a lost datum; it records agreement that was never given.
This is the exact pathology my signoff (consensus.md:284-291) and claude-1's closing
note (consensus.md:444-449) warned about, and it is the one thing that must not
happen.

The fallback does NOT save this case. The fallback fires only when `len(items)==0`
(all extraction fails across all files) or the previous round is empty. If OTHER
lines in the same or other files match markers, `items` is non-empty, no fallback
fires, and the missed objection is silently dropped while the ledger presents itself
as complete.

Fix: the marker list must be case-insensitive and must include the actual objection
verbs agents use. At minimum add: "reject", "object", "oppose", "abandon",
"counter" (case-insensitive), "concern" (as a heading or lead word), "disagree"
(bare, not just "I disagree"), "withdraw" (already present as "I withdraw" but
should be case-insensitive "withdraw"). Better: since the protocol's review template
uses `### [CRITICAL]` / `### [MAJOR]` etc. and the design template uses `## Concerns`,
`## Remaining disagreements`, `## New concerns`, `## Risks` — carry every line under
those section headings, not just lines matching severity words. Best: do not rely on
substring matching at all for the review path; carry the `## Findings` section
verbatim from every review round, since that is where live findings live.

### [CRITICAL] Review rounds 1–2 double FINAL.md + IMPLEMENTATION.md in the context

`gatherReviewContext` (`phase58.go:486-507`) always prepends `head = FINAL.md +
IMPLEMENTATION.md` (lines 487-493), then calls `frontierContext(round, full)`. For
round <= 2, `frontierContext` returns `full()` = `gatherReviewContextFull`, which
ALSO includes `FINAL.md + IMPLEMENTATION.md` at its start (lines 283-288). The
deduplication guard at line 503 (`if strings.Contains(rounds, "carry-forward
fallback")`) only fires when the fallback string is present — which never happens
for round <= 2, because `frontierContext` returns `full()` directly without wrapping
it in the fallback banner.

Result: for review round 1 and review round 2, `FINAL.md` and `IMPLEMENTATION.md`
appear TWICE in the reviewer's context. This is a context-size regression on the
very rounds G1 promises are "unchanged from today." The old `gatherReviewContext`
returned `FINAL.md + IMPL.md + prior review rounds` once. The new one returns
`[FINAL+IMPL] + [FINAL+IMPL + prior review rounds]`.

This also affects the review-consensus drafter at cycle 1, which calls
`gatherReviewContext(path, roundNumber(opts)+1)` = `gatherReviewContext(path, 2)`
(`runner.go:920`), hitting the same duplication.

Fix: in `gatherReviewContext`, do not prepend `head` when `frontierContext` will
return `full()` (round <= 2). Either skip the head prepend for round <= 2, or
restructure so `frontierContext` does not include `FINAL.md`/`IMPLEMENTATION.md`
in its `full()` for the review path (pass a review-specific full renderer that
excludes them). The cleanest fix: for round <= 2, return `gatherReviewContextFull`
directly without the head prepend, since the full renderer already includes them.

### [MAJOR] Fallback completeness: 5 of 6 required uncertainty states have no code path

My signoff (consensus.md:273-281) and codex-1's (consensus.md:197-201) required
these uncertainty states to trigger fallback to full history:

1. missing
2. invalid
3. ambiguous
4. challenged
5. unresolved hash/locator
6. verdict conflict not marked DISPUTED

Plus codex-1's validation list: duplicate/mutated IDs, silent deletion, unauthorized
transitions, dangling/cyclic links, missing active-participant ledgers.

What the code actually checks (`frontier.go:80-118`):
- `os.ReadDir` / `os.ReadFile` error → propagated as error (not fallback, but blocks)
- `len(items) == 0` → fallback ("extracted no items")
- previous round empty → fallback ("no readable artifacts")

States that DO trigger fallback: empty extraction, empty previous round, filesystem
errors.

States that DO NOT trigger fallback:
- **invalid** — no validation exists. There are no IDs, no hashes, no lifecycle
  states to validate. The `ledgerItem` struct has `Owner`, `Round`, `Kind`, `Line`,
  `Locator` — nothing to invalidate.
- **ambiguous** — no ambiguity detection. A line matching multiple markers gets the
  first match's `Kind`; this is not even flagged.
- **challenged** — no citation/provenance challenge detection. "Challenged" is not
  a concept in the code.
- **unresolved hash/locator** — no hashes exist. The locator is `file:line` but is
  never verified to resolve to an actual line.
- **verdict conflict not marked DISPUTED** — G6, knowingly unmet. The extractor
  carries verdict lines verbatim but does not join claims or detect conflict.
- **missing per-participant ledger** — partial. If ALL participants have no markers,
  fallback fires. If some have markers and others don't, no fallback — the
  participant with no extractable lines gets an empty `--- <owner> ---` header
  that implies they were fully represented.

The derived ledger's simplicity makes most of codex-1's validation list N/A (no
IDs → no duplicate/mutated IDs; no lifecycle → no unauthorized transitions; no
links → no dangling/cyclic). But "challenged", "ambiguous", and "verdict conflict"
are real states that can occur and are not caught. The code's fallback is purely
"extraction found nothing" + "previous round empty" + "filesystem error" — three
cases, where my signoff required six.

Fix: at minimum, add a per-participant check: if any active participant's file
yielded zero items, either trigger fallback or annotate the empty `--- owner ---`
header with "(no marker-matching lines extracted; open the full artifact for this
agent's complete position)". For the verdict-conflict case (G6), see the next
finding.

### [MAJOR] G6 unmet — verdict conflicts cross the compaction boundary as unrelated lines

G6 requires: "The same material claim reworded under a new ID with opposing PRIMARY
verdicts in different rounds must join as DISPUTED or trigger fallback." The
implementer honestly reports this is not implemented and not tested.

I judge this MAJOR, not CRITICAL, for this reason: the derived ledger carries
verdict lines verbatim. If agent A says "PRIMARY: measured 3.3x" in round 1 and
agent B says "PRIMARY: measured 2.1x" in round 2, both lines appear in the ledger
for round 3. An agent reading the ledger sees both verdicts, just not joined as
DISPUTED. The conflict is visible; the system just doesn't label it.

The risk is that a conflict spanning the compaction boundary (verdict in round 1,
opposing verdict in round 3) reaches round 5 as two unrelated ledger lines under
different owner headers, and an agent skimming the ledger may not connect them.
This is a real defect but not a silent-drop: the information is present, just
unstructured. The protocol's §15.6 duty still binds the consensus drafter, who
receives full history.

However, shipping without G6 means the one gate that tests claim-ID forking is
not exercised, and a verdict conflict that spans the boundary is dependent on
agent attention for detection. The fix should be one of:
(a) implement claim joining: match verdict lines by normalized claim text across
rounds and emit a DISPUTED marker when opposing PRIMARY verdicts are found; or
(b) trigger fallback when any round in 1..N-2 contains a PRIMARY verdict and any
other round in 1..N-2 contains a different PRIMARY verdict (overly broad but safe);
or
(c) document G6 as deferred and accept the risk for v1, with the explicit caveat
that a verdict conflict spanning the compaction boundary may not be joined.

I lean toward (c) for v1 with (a) as a fast-follow, but this is the implementer's
call to make explicit, not to leave silent.

### [MAJOR] The derived ledger has no lifecycle — withdrawn objections never disappear

The implementer flags this as deviation 1 and asks reviewers to judge. My judgment:
acceptable for v1 IF AND ONLY IF the marker list is fixed (see CRITICAL above) and
the ledger header's owner-disposes rule is sufficient to prevent false convergence.

The derived ledger carries verbatim lines matching markers, attributed by filename.
There are no lifecycle states (OPEN/RESOLVED/DEFERRED/SUPERSEDED), no transition
history, no supersedes links. A withdrawn objection keeps appearing in every
subsequent round's ledger until its round falls out of scope. The "I withdraw"
marker carries the withdrawal, but the original objection is also carried — the
agent sees two unrelated verbatim lines and must infer the relationship from prose.

This is not "silent deletion" (the opposite problem). It is noise: the ledger grows
monotonically with every round's extracted lines, partially eating the saving the
optimization claims. codex-1 warned about this (consensus.md:212-215): "a cumulative
ledger still grows with newly introduced items, so the implementation should measure
packet growth and avoid claiming a formal complexity reduction it has not
demonstrated." No packet-growth measurement is in the implementation or tests.

kimi-1's re-litigation blindness fixture (consensus.md:356-369) is also unaddressed:
a disposed objection re-raised reworded in a later round cannot be detected as
re-litigation because there are no supersedes links. Both the original and the
rewording appear as unrelated lines.

Why MAJOR, not CRITICAL: the information is present (both the objection and the
withdrawal are carried verbatim), so an attentive agent can reconstruct the
disposition. The risk is noise and re-litigation, not false convergence — unless
the marker miss (CRITICAL above) drops the withdrawal line too, in which case the
objection appears live forever with no visible disposition.

Fix for v1: accept the derived ledger without lifecycle, but (1) fix the marker
list so withdrawals are always carried, (2) add a note to IMPLEMENTATION.md that
packet growth is unmeasured and the complexity reduction is not formally
demonstrated, (3) track kimi-1's re-litigation fixture as a required follow-up.

### [MINOR] Instruction at runner.go:997 blurs the context-optimization / consensus-rule boundary

The new instruction says: "an objection in it is live until ITS OWN OWNER withdraws
it." This is correct (it restates codex-1's rule) and is a reading instruction, not
a validation rule. The code does not enforce it as a consensus-close condition —
the consensus drafter's §15.6 duty governs close, and that path is untouched.

But the instruction's normative language ("is live until ITS OWN OWNER withdraws
it") tells agents how to judge whether an objection is live, which directly affects
whether they consider consensus reached. If an agent treats this as a close
criterion, the optimization has introduced a de facto consensus rule through the
prompt, even though no code enforces it. The boundary my signoff and codex-1's set
was: "implementation-scoped context optimization, not an artifact-validity or
consensus rule."

Fix: frame the instruction as context-reading guidance, not a close condition.
For example: "The ledger is a reading aid, not a consensus rule; the consensus
drafter's full-history audit governs close." This keeps the owner-disposes rule as
a reading instruction while making clear it is not a new close condition.

### [MINOR] Review prompt does not explain the ledger the way the design prompt does

`BuildReviewPrompt` (`phase58.go:236-273`) was not updated to mention the ledger or
the owner-disposes rule. The design-side `BuildRoundPrompt` (`runner.go:990-1019`)
was. The `renderLedger` header (`frontier.go:145-149`) does state the owner-disposes
rule and the "full artifacts on disk" note, so the information reaches the reviewer
via the context blob — but the review prompt template itself doesn't reinforce it.

This is a NIT-to-MINOR asymmetry: the review path (the larger of the two quadratic
paths) relies on the ledger header for interpretation guidance, while the design
path reinforces it in both the prompt and the header. The information is present in
both; the design path is just more explicit.

Fix: add a sentence to `BuildReviewPrompt` noting that older review rounds appear as
a verbatim carry-forward ledger and that the owner-disposes rule applies, mirroring
the design prompt's language.

### [NIT] G7 reversion check for fallback announcement not completed

The implementer honestly reports that 2 of 3 G7 reversion checks were completed and
the third (fallback announces itself) was interrupted when the harness was killed
mid-case, leaving a revert in the working tree that was manually restored.

The test `TestFallbackSelectsFullHistoryAndAnnouncesItself` checks for the
"carry-forward fallback" string, so reverting the banner would make it fail. The
G7 verification is logically sound but not empirically demonstrated for this case.
The implementer's note about the harness needing to restore on kill (not just on
exception) is a valid process observation.

Fix: complete the third reversion check. The test is correct; this is verification
rigor, not a code defect.

### [NIT] Empty `--- owner ---` headers imply full representation

`renderLedger` (`frontier.go:156-164`) emits `--- <owner> ---` for every owner who
had a file, even if that owner had zero extractable marker lines. An agent seeing
`--- d ---` with nothing under it may conclude d had nothing to say, when d actually
had concerns that didn't match markers. This compounds the CRITICAL marker-miss
finding: not only is the objection dropped, but the empty header implies d was
fully represented.

Fix: either omit owners with zero items, or annotate the empty header with "(no
marker-matching lines extracted from this agent's older rounds; open the full
artifact for their complete position)."

## Does the change break `parley run`?

No. `parley run` does not crash or produce invalid prompts:
- Round 2 (design): `frontierContext(2)` returns `full()` = `gatherPriorRounds(path,
  2)` — exactly the old behavior. G1 holds for the design path.
- Review round 1–2: works but doubles `FINAL.md + IMPLEMENTATION.md` in the context
  (CRITICAL finding above — context bloat, not breakage).
- Design consensus drafter: `buildConsensusDraftPrompt` (`driver_consensus.go:110`)
  is untouched, still reads "EVERY round artifact under round-*/". G4 holds.
- Review consensus drafter cycle 1: calls `gatherReviewContext(2)` → same
  duplication as review round 2. Cycle 2+: calls `gatherReviewContext(3+)` →
  compaction path, correct.

No validation failure, no crash, no wrong prompt structure. The duplication is a
byte-count regression, not a logic error.

## Open questions

1. The CRITICAL marker-miss finding and the CRITICAL FINAL/IMPL duplication are both
   fixable in a few lines. Should this review block ship until they are fixed, or
   ship with the fixes as a fast-follow? My position: the marker-miss is a CRITICAL
   that blocks ship (it converts ungiven consent into recorded agreement); the
   duplication is a CRITICAL by severity but is a regression on the "unchanged" path
   and should also block.

2. For G6 (verdict conflict joining), is option (c) — document as deferred and accept
   the risk for v1 — acceptable to the other reviewers? The information is present
   but unstructured. I lean yes for v1 with (a) as a fast-follow.

3. The derived ledger without lifecycle is acceptable for v1 only if the marker list
   is comprehensive enough that no live objection drops. Given the CRITICAL marker
   finding, the derived ledger is not yet safe to ship as-is. Fix the markers first,
   then the derived-without-lifecycle form is acceptable as a v1.
