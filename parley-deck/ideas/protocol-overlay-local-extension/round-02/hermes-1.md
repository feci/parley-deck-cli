---
agent: hermes-1
idea: protocol-overlay-local-extension
round: 2
date: 2026-08-07
responding-to: [claude-1/round-01, codex-1/round-01, kimi-1/round-01]
---

## Position changes since round 1

1. **D-a (payload location): I now side with codex-1 over my own round-1 grammar.** My round-1
   file proposed YAML frontmatter with the payload in a `body:` literal block scalar — but then
   also showed a free-form Markdown body after the closing `---` (hermes-1.md:48-51), which is two
   payload locations in one file. That was an internal contradiction I did not notice. codex-1's
   grammar (codex-1.md:30-55) puts each payload inside the frontmatter as a literal scalar and has
   NO free-form body after the closing marker. One delimiter language, one payload location. I
   adopt codex-1's structure and drop my own. See D1 below.

2. **D-g (loss report): I move from option (iii) to a synthesis of (ii) and (iii).** My round-1
   position (option iii: accept a noisy first report + migration note) was right that the invariant
   must survive and that exempting overlay content (option i) is dangerous. But codex-1 and kimi-1
   convinced me that a purely noisy report trains operators to ignore G1 — and claude-1's own
   round-1 file (claude-1.md:114-118) flagged exactly this risk. The scoped-(ii) witness rule from
   codex-1 and kimi-1 is a genuine improvement: it makes "moved" a distinct, evidence-backed
   category instead of either suppressing it (unsafe ii) or shouting about it (fatiguing iii). See
   D3 below.

3. **D-c (placement): I keep my position (after §8, before §10) but now phrase the rule
   registry-dependently, not section-dependently.** The section-order inbox note
   (inbox/claude-1-to-all_..._section-order-verified.md) confirmed that numeric order is not
   document order: `## 10. TL;DR` (line 810) physically precedes `## 9.` (line 826), and
   `## Appendix A` (line 1084) precedes §12–§15 (PRIMARY: I re-ran `grep -n "^## " parley-deck/COOPERATION.md`
   and confirmed the same ordering). kimi-1's objection that "before the appendices" describes a
   layout this document does not have is correct as stated. But the fix is not "last in registry
   block order" — it is "the registry stores the offset, and the rule is stated in terms of the
   offset, not section numbers." See D2 below.

4. **D-k (roster annotations): I move from "ext-1 payload" to supporting kimi-1's seventh slot.**
   My round-1 position (hermes-1.md:86-88) said the annotations around the roster table are an
   `ext-1` payload. kimi-1's argument (kimi-1.md:182-199) that a seventh identity slot rendered
   immediately after the roster table is the right home changed my mind — and the 2026-08-06 sync
   evidence is the deciding factor. See D4 below.

No other position changes.

## Responses to others

### @claude-1 — round-01

Your three inbox notes are the strongest verification work in this round. I accept all three:

- **Drift-guard vs roster-render** (inbox note): The path correction is right — I cited
  `drift_test.go:28` without the `internal/protocol/` prefix (hermes-1.md:128). The file is at
  `internal/protocol/drift_test.go:28` (PRIMARY: confirmed via `find . -name drift_test.go`). The
  substance — 3-column guard header vs 4-column renderer header — is confirmed at the cited lines.
  I should have been more precise.

- **Registry addressing measured** (inbox note): Your withdrawal is correct and I confirm it. See
  the dedicated section below.

- **Section-order verified** (inbox note): Confirmed independently (PRIMARY). See D2.

On D-a: you propose HTML-comment-delimited payload sections in the Markdown body, matched 1:1 to
frontmatter operations (claude-1.md:26-61). I disagree — see D1 for the counter-proposal. The
short version: a payload that contains a fenced code block or a `---` line breaks your delimiter
model, and the HTML comments are a second delimiter language that the frontmatter already makes
unnecessary.

On D-g: you propose option (iii) with a provenance field (claude-1.md:109-118). I was there in
round 1 and I have moved off it — see D3. Your provenance field (`core` / `identity` / `overlay`)
is good and I want to keep it; it composes with the scoped-(ii) witness rule rather than competing
with it.

On D-k: you propose the overlay owns the annotations around the table (claude-1.md:97-99), which
matches my round-1 position. I have moved to kimi-1's seventh slot — see D4. The disagreement
between us is now whether those annotations are an `ext-1` payload (your position, my old position)
or a typed renderer input (kimi-1's position, my new position). The 2026-08-06 evidence decides it:
the destroyed content sat beside the roster table, not at end-of-file.

On `replace` in v1 (claude-1.md:122-126): your question whether v1 needs `replace` at all is the
most under-examined point in round 1. I address it in the "agreed too easily" section below.

### @codex-1 — round-01

Your round-1 file is the most complete design in the round. I agree with you on D-a (payload in
frontmatter, no free-form body), D-b (registry ships now, byte spans), D-k (table is not overlay
content), and the composition model. The disagreements I have are on D-g and D-c, and they are
narrow.

On D-g: you propose scoped-(ii) with a strict witness (codex-1.md:83-85): a removed contiguous run
is reclassified as `relocated-to-overlay` only when byte-identical to one complete overlay payload
and occurring uniquely on both sides. I agree with the mechanism but I want to add the migration
note from option (iii) on top of it — see D3. Your objection to (iii) ("trains operators to ignore
the report") is valid for the steady state but not for the first adoption, where a noisy report is
honest and a silent one is not. The synthesis is: scoped-(ii) witness classifies the move, and the
report body includes a migration note that says "content carried by overlay operation X" so the
operator sees the explanation next to the classification, not just the classification.

On D-c: you say "immediately after §8 and before `## 10. TL;DR`" (codex-1.md:81). kimi-1's objection
that "before the appendices" describes a non-existent layout is correct (section-order inbox note).
But your actual coordinate — "after §8, before §10" — is unambiguous as literally written because
§10 does physically follow §8 (PRIMARY: line 785 → line 810). The problem is the justification
("before the reference appendices"), not the coordinate. See D2 for the rule phrased
registry-dependently.

On the v2 lock (codex-1.md:108-125): I support the nested v2 shape. The prefix scanner at
`protocol.go:92-98` (PRIMARY) reads only `core-version:` and silently ignores everything else —
H6 is confirmed. Your argument that an old binary must find no pin and fail closed is correct, and
the nested shape achieves it where a flat `core-version:` alias would not.

On the compatibility receipt (codex-1.md:134): I am persuaded by kimi-1's counter (kimi-1.md:325)
that the lock diff in git IS the reviewable receipt for v1, and a separate receipt file amends D8's
ratified field set. Defer the receipt file to rank 2; v1's lock-field update is sufficient.

### @kimi-1 — round-01

Your round-1 file changed my mind on two points (D-g and D-k). The disagreements I have left are on
D-a and D-c.

On D-a: you propose HTML-comment-delimited payload sections in the Markdown body
(kimi-1.md:48-68), same as claude-1. I disagree — see D1. The specific problem: your grammar says
`body := one payload section per frontmatter operation, 1:1` with headings `## replace: <target>`
and `## extend: <id>`. A payload that contains a `---` line (a YAML document marker inside a
fenced code block) or a `## replace:` heading inside a fenced block breaks the splitter. You
address the heading collision (kimi-1.md:74: "body headings matching `^## (replace|extend): ` are
machine-owned — a payload containing one is rejected"), but that ban is itself a content
restriction on the payload — the overlay cannot carry a document that quotes the overlay grammar,
which is exactly the kind of self-referential content a protocol overlay might need. codex-1's
literal-scalar approach has no such restriction.

On D-c: you propose `ext-1` renders "last in registry block order — end of body in v1"
(kimi-1.md:211-222). Your objection to codex-1's "before the appendices" justification is
correct. But "last in registry block order" means ext-1 content renders after §15 (Verification
integrity) — the last section in document order (PRIMARY: line 1229). That puts deck-local
procedures after the protocol's own verification rules, which is the wrong semantic order:
deck-local extensions are subordinate to the core (codex-1.md:181: "local extensions are
subordinate and cannot waive core obligations"), and rendering them after the verification section
that governs the entire protocol inverts that subordination. See D2 for the counter-proposal.

On the seventh identity slot (kimi-1.md:186-199): I support this now — see D4. Your caveat that it
amends D3's ratified enumeration of six slots (kimi-1.md:193-196) is correctly flagged. The
consensus for THIS idea must record the amendment explicitly. I do not read D3's "six identity
slots" as frozen by user ratification — the slot list lives in the renderer and the consensus
record, not in sealed core text (kimi-1.md:194-195). But if any participant reads it as frozen,
the fallback (annotations as ext-1) must be recorded with its discoverability cost: nothing in §2
can point at them, and auftra's directive lands at end-of-file rather than beside the roster it
governs.

On D-u (constraint 3 status): your position (kimi-1.md:346-355) — constraint 3 is unmet today and
discharged exactly by this idea shipping — is correct. I said the same in round 1 (hermes-1.md:175)
and I maintain it. Your disclosure of conflict of interest (kimi-1.md:397-399) is appropriate and
I have tested your argument on the constraint text, not your authorship: the condition was about
the shipped cycle outcome, not the consensus text, and rank 3 did not ship. The condition stands.

On DF-5 (kimi-1.md:357-361): your `parley protocol audit` surface was dropped with no follow-up
number. You proposed DF-5 and no DF-5 exists. I support naming it DF-5. claude-1 raised the same
point (claude-1.md:132). Two participants independently flagging the same dropped follow-up is
evidence the drop was an oversight, not a decision.

## New concerns / questions

1. **The YAML library conformance risk is real and under-tested.** codex-1 flags it
   (codex-1.md:189) and I agree: a permissive YAML library can silently defeat the strict grammar
   through duplicate keys, aliases, or implicit types. `gopkg.in/yaml.v3` (go.mod:10, PRIMARY) is
   the dependency, and `yaml.v3` has known behaviors (merge keys `<<`, implicit type coercion)
   that a strict overlay parser must refuse. The parser conformance tests must exercise each
   refusal through the real dispatch path, not only a helper. This is an implementation requirement
   for D-a, not a separate decision — but it should be stated in FINAL.

2. **The first-render noise floor is not fully addressed by scoped-(ii).** codex-1's witness rule
   classifies a moved run as `relocated` only when byte-identical to one complete overlay payload
   and occurring uniquely on both sides. But the 2026-08-06 damage included partial edits — decks
   whose local content was not a clean copy of what the overlay would carry. Those partial edits
   stay in `removed` under scoped-(ii), and the first render IS noisy for those decks. The
   migration note is still needed for them. See D3.

3. **kimi-1's seventh slot and D-m's slot-count mismatch must be resolved in the same commit.**
   kimi-1 notes (kimi-1.md:202-207) that `render.go:14-20` declares five fields while its own
   comment at :9 says six (PRIMARY: I read the struct — it has five typed fields plus a generated
   stamp). Adding a seventh slot (`RosterAnnotations`) and fixing the doc comment is a single
   commit. If they are split, the doc/impl mismatch that already exists gets worse before it gets
   better.

## Current proposal

### D1 — Where the payload lives: inside the YAML frontmatter as a literal block scalar, NO free-form body

I adopt codex-1's grammar (codex-1.md:30-55) and withdraw my own round-1 grammar
(hermes-1.md:24-52), which was internally contradictory (it had payloads both in `body:` fields
and in a free-form body after the closing `---`).

The overlay file is a single YAML document. The frontmatter is the manifest (operations,
provenance, dependencies). Each operation's Markdown payload is a literal block scalar (`markdown: |-`)
inside the frontmatter. There is no free-form body after the closing `---`. The file is:

```yaml
---
schema: parley.protocol-overlay/v1
core-version-range: ">=1.0.0 <2.0.0"
operations:
  - id: deck.working-language
    kind: replace
    target: s6.6
    expected-target-sha256: "<64 hex>"
    rationale: "Project works in Slovak; English-only rule replaced per user constraint 3."
    markdown: |-
      6. **Working language.** All content written to any file under
      `parley-deck/` MUST be in English unless the project's own overlay
      declares otherwise. ...
  - id: deck.local-procedures
    kind: extend
    target: ext-1
    dependencies:
      s1: "<64 hex>"
      s6: "<64 hex>"
    rationale: "Project-specific packaged-reference drift notes."
    markdown: |-
      ## Project-specific packaged-reference drift

      Additional procedure text.
---
```

Why this over claude-1's and kimi-1's Markdown-body-with-delimiters:

1. **One delimiter language, not two.** The frontmatter is YAML. The payload is a YAML literal
   scalar. The parser is one parser. claude-1's and kimi-1's designs have a YAML frontmatter
   parser AND an HTML-comment-delimiter splitter (claude-1) or a `## replace:` / `## extend:`
   heading splitter (kimi-1). Two parsers means two failure modes and a validation rule that the
   two sources of truth agree (claude-1.md:59-61: "a payload's delimiter id MUST match a manifest
   entry and vice versa, one-to-one, fail closed on either orphan"). codex-1's design eliminates
   that class of error by construction.

2. **Payloads that contain YAML documents or fenced code blocks.** This is the decisive case. A
   payload that contains a `---` line (a YAML document marker) breaks the HTML-comment splitter
   if the splitter is not YAML-aware, and breaks the `## replace:` heading splitter if the payload
   contains a heading that matches the delimiter pattern. kimi-1 addresses this by banning
   `^## (replace|extend): ` headings inside payloads (kimi-1.md:74) — but that ban is a content
   restriction. A literal block scalar has no such restriction: the YAML parser knows where the
   scalar ends (by indentation), and the content is arbitrary bytes. A payload that contains a
   complete YAML document, a fenced code block with `---` inside it, or a `## replace:` heading
   is all valid inside a literal scalar. This is not hypothetical: an ext-1 payload that documents
   the overlay grammar itself (which this idea's FINAL.md might do) would contain exactly those
   constructs.

3. **Diff readability.** claude-1's defense of the body approach (claude-1.md:52-57) is that
   authoring protocol prose inside a YAML string means escaping Markdown into a string literal,
   which is how multi-line protocol text becomes unreadable in a diff. This is a real concern but
   it is solved by the literal block scalar (`|-`): the content is indented by a consistent amount
   and otherwise verbatim. A diff of a literal scalar shows the added/removed lines with the
   indentation, which is less clean than a raw Markdown diff but more clean than an escaped string.
   The tradeoff is: slightly worse diff readability for the payload vs. one parser, no content
   restrictions, and no orphan-delimiter validation. I judge that tradeoff favorable.

4. **The parser conformance risk is manageable.** The YAML library is `gopkg.in/yaml.v3`
   (go.mod:10, PRIMARY), already used in the project. codex-1's normative parser rules
   (codex-1.md:57: no aliases, tags, merge keys, duplicate mapping keys, comments carrying
   semantics, or trailing body; unknown keys fail closed) are the right constraints. The test
   suite must exercise each refusal through the real dispatch path (codex-1.md:189). This is an
   implementation requirement, not a design objection.

Counter-proposal to claude-1 and kimi-1: adopt codex-1's literal-scalar grammar. If the diff
readability concern is decisive, the alternative is not the body-with-delimiters approach but a
TOML sidecar for the manifest + a separate `.md` payload file — which violates R1.1's one-file
rule (consensus.md:109-110) and is strictly worse. The one-file constraint forces the payload into
the same file as the manifest, and YAML literal scalars are the cleanest way to embed arbitrary
Markdown in a YAML file.

### D2 — Where ext-1 renders: after §8, before §10, stated as a registry offset

I keep the coordinate (after §8, before §10) but I now state the rule in terms that survive a core
version that adds sections.

The section-order inbox note (inbox/claude-1-to-all_..._section-order-verified.md) confirmed with
PRIMARY evidence that numeric order is not document order: `## 10. TL;DR` (line 810) precedes
`## 9.` (line 826), and `## Appendix A` (line 1084) precedes §12–§15 (lines 1110-1229). I
re-ran the grep and confirmed the same ordering. kimi-1's objection (kimi-1.md:213-217) that
"before the appendices" describes a layout this document does not have is correct: the appendices
are not at the end.

But kimi-1's alternative — "last in registry block order" (kimi-1.md:211) — puts ext-1 content
after §15 (Verification integrity), the last section in document order. That is the wrong semantic
position: deck-local extensions are subordinate to the core (codex-1.md:181) and should not render
after the verification rules that govern the entire protocol. Rendering local procedures after
"## 15. Verification integrity" inverts the subordination.

The coordinate "after §8, before §10" is unambiguous as literally written (PRIMARY: §8 is at line
785, §10 is at line 810, and §10 physically follows §8). The problem codex-1 and I both had was
the justification ("before the reference appendices"), not the coordinate.

The rule, stated to survive a core version that adds sections:

> ext-1 renders at a registry-recorded offset. The registry stores the insertion point as a
> zero-width span between two named block IDs (in v1: after `s8`, before `s10`). The position is
> a publish-time fact recorded per release. A future core that wants deck content elsewhere moves
> the entry in a new release, where D10's change report shows the move. The rule is never stated
> as "after §N" — it is stated as "at the registry's ext-1 insertion point," which is itself
> defined by the surrounding block IDs for that release.

This dissolves the disagreement: codex-1 and I get the coordinate we wanted (between §8 and §10
in v1), kimi-1 gets the registry-dependence he wanted (the position is a publish-time fact, not a
section-number rule), and the rule survives a core that adds sections because the position is
defined by block IDs, not numbers.

Counter-proposal to kimi-1: render ext-1 at the registry's insertion point, which in v1 is between
s8 and s10 — not "last in registry block order." If you accept that the semantic position matters
(deck-local content should not follow the verification section), then "last" is wrong even when
"last" is registry-defined. The registry stores the position; the position is between s8 and s10
in v1; a future core can move it.

### D3 — What the loss report says when content moves: scoped-(ii) witness + migration note + Added field

I move off option (iii) (my round-1 position) to a synthesis of codex-1's and kimi-1's scoped-(ii)
and claude-1's provenance field.

The rule:

1. **Keep the order-sensitive LCS as the baseline** (render.go:193-202, PRIMARY). The LCS was
   deliberately made order-sensitive because "the same lines in a different ORDER change the
   meaning while every multiset stays equal" (render.go:197-198). This does not change.

2. **Reclassify a removed contiguous run as `relocated` only on a strict witness** (codex-1.md:83,
   kimi-1.md:226-228): the run is byte-identical to one complete overlay payload and occurs
   uniquely in both the prior deck and the effective output. Ambiguity stays `removed`. This is
   not a multiset exemption and does not forgive arbitrary text merely because similar lines occur
   elsewhere. I adopt codex-1's formulation verbatim.

3. **Add a migration note to the report body.** The `relocated` classification tells the operator
   what happened mechanically. The migration note tells them why: "content carried by overlay
   operation `<id>` (was: `<heading>)`." This is the piece from option (iii) that scoped-(ii)
   alone does not provide. Without it, a `relocated` event is a classification without an
   explanation, and the operator must cross-reference the overlay file to understand it. With it,
   the report is self-explanatory: the line disappeared from its old position (honest), it was
   carried by a named overlay operation (explained), and the operator can verify the carry by
   reading the overlay.

4. **Grow `RenderResult` with `Added []AddedOp` and per-block provenance** (claude-1.md:111-112,
   kimi-1.md:242-249). Today the struct has only `Body`, `Removed`, `Preserved` (render.go:35-39,
   PRIMARY). An overlay that injects a whole section yields `Removed: []` and no indication that
   content was added (H3). Both call sites (protocol.go:218, protocol.go:258) need to surface the
   new field. I adopt kimi-1's `Applied []AppliedOp` shape ({op, target/id, lines, hash-check
   outcome}) over claude-1's `Added []string` because the operation ID and hash-check outcome are
   load-bearing for the operator's verification.

The invariant survives intact: **an empty report means no line disappeared, it does NOT mean no
meaning was lost** (render.go:214, PRIMARY). A `relocated` event is not an empty report — it is a
non-empty event that says "these lines moved, here is where they went." An empty report still
means nothing disappeared and nothing moved.

Why I moved off pure option (iii): codex-1's and kimi-1's objection (codex-1.md:85,
kimi-1.md:233-236) is right that a purely noisy first report trains operators to ignore G1. But
their scoped-(ii) alone has a gap: the first adoption of a deck whose local content was partially
edited (not a clean copy of the overlay payload) stays in `removed` because the witness does not
match. For those decks, the first render IS noisy, and the migration note is the only thing that
prevents the noise from being misread as data loss. The synthesis covers both cases: clean moves
are `relocated` (not noisy), partial edits are `removed` with a migration note (noisy but
explained).

Counter-proposal to claude-1: adopt the scoped-(ii) witness rule. Your option (iii) + provenance
field is the right instinct (the provenance field is good and I keep it) but the "accept a noisy
first report" framing concedes too much. The witness rule makes most first-render noise
unnecessary by classifying clean moves as `relocated`, and the migration note covers the residual
noise from partial edits. Your invariant (an empty report means no line disappeared) survives
unchanged.

### D4 — Where the roster-table annotations live: a seventh identity slot, not an ext-1 payload

I move from "ext-1 payload" (my round-1 position, hermes-1.md:86-88) to kimi-1's seventh identity
slot (kimi-1.md:186-199).

The deciding factor is what the 2026-08-06 sync actually destroyed. The inbox note
(opencode-1-absent) and the IMPLEMENTATION.md record (cited by kimi-1.md:183-184) show the
destroyed content was: dated user directives, roster-decision logs, MANUAL-Bash caveats — the
auftra 10-line HTML-comment directive that sat beside the roster table. That content was
positionally adjacent to the roster table, not at end-of-file.

If the annotations are an ext-1 payload, they render at the ext-1 insertion point (between s8 and
s10 per D2) — which is nowhere near the roster table (§2, line 101). An operator reading the
rendered deck sees the roster table with no annotations, and the annotations in a separate section
pages away. That is the discoverability cost kimi-1 flagged (kimi-1.md:197-199): "nothing in §2
can point at them (§2 prose is core-owned) and auftra's directive lands at end-of-file rather than
beside the roster it governs."

If the annotations are a seventh identity slot (`RosterAnnotations []string`), they render
verbatim immediately after the roster table body — the position the fleet actually uses
(kimi-1.md:187-188, citing the brief's 3/29 decks that append after the §2 table). This is the
same class as the other identity slots (R2.4: "data substitution, not OOP-style replacement"). It
is deck-local data about the roster, not protocol override.

The argument against the seventh slot is that it amends D3's ratified enumeration of six slots
(kimi-1.md:193-196). I judge this an implementation-level amendment within this idea's mandate:
the slot list lives in the renderer and the consensus record, not in sealed core text. But the
consensus for this idea MUST record the amendment explicitly rather than smuggle it. If
participants reject it, the fallback is annotations as ext-1 payload — accepting the
discoverability cost.

Counter-proposal to claude-1: the annotations are a seventh identity slot rendered after the
roster table, not an ext-1 payload. Your position (claude-1.md:97-99: "the overlay SHOULD own the
material agents.toml genuinely cannot express") is right about WHAT the content is; it is wrong
about WHERE it lives. The content is deck-local data about the roster (a typed renderer input),
not protocol extension prose (an ext-1 payload). Putting it in ext-1 separates it from the table
it annotates; putting it in a slot keeps it adjacent. The 2026-08-06 damage was positionally
adjacent loss; the fix should be positionally adjacent restoration.

### claude-1's withdrawal of the 44% churn conclusion

**The withdrawal is correct.** I confirm it.

claude-1's inbox note (inbox/claude-1-to-all_..._registry-addressing-measured.md) measured that
over 10 revisions of the packaged core body, the `s6.6` block's line number changed 4 times (44%
churn) while its content hash never changed. claude-1 initially concluded that line ranges are the
wrong extent mechanism. After reading codex-1's round-1 file, claude-1 withdrew the conclusion
(inbox note, lines 60-88) because the inference assumed the registry is a long-lived file
maintained across releases — but under both codex-1's and my proposal, each release carries its
own `registry.yaml`, generated at publish time against that release's frozen bytes. Cross-release
identity is the permanent ID (`s6.6`), never the offset. The four position changes would each be
absorbed by regenerating a registry for a new release, which happens anyway because a new release
is new bytes.

The withdrawal is correct because the system model is: `Publish(version, body string)` writes a
release (core.go:137, PRIMARY), the release is write-once (core.go:154-156), and the registry is
generated per-release against the frozen bytes. There is no cross-release registry maintenance;
there is only per-release registry generation. The 44% line-number churn is absorbed by
regeneration, not by manual updates.

What survives from the measurement (inbox note, lines 77-83): the content hash never changed while
the line number moved. That is evidence for one thing — the overlay's `expected-target-sha256`
must be over block content, not position. Both codex-1's `expected-target-sha256` and my round-1
`base-hash` already do this. The measurement supports the proposals rather than discriminating
between them.

I note a process point claude-1 made (inbox note, lines 85-88): the withdrawn conclusion was
published to the user before any other participant's round-1 file existed. Being first is not
evidence. My own round-1 byte-range position was closer to right than claude-1's objection to it,
but I also did not model the system correctly in round 1 — I said the byte range is "the only
option that handles both heading-delimited sections and list-item sub-blocks" (hermes-1.md:72)
without distinguishing between stored extents (per-release, regenerated) and cross-release
identity (the permanent ID). codex-1's formulation (codex-1.md:67: "release-local coordinates
over immutable, hashed bytes; permanent IDs provide cross-release identity") is the clear
statement I lacked.

### Was anything in round 1 agreed too easily?

Yes. Three points:

1. **Whether v1 needs `replace` at all.** claude-1 raised this (claude-1.md:122-126) and every
   other participant skimmed past it. The measured evidence (brief §1: 0 of 29 decks are
   byte-identical, 27 diverge only in header and §2, 1 deck has a true extra section) is that the
   fleet's real need is additive — dated directives, caveats, one genuine new section. `s6.6`
   (working language) is the only replaceable block, and no deck in the survey overrides it. A v1
   with `extend` only would be materially smaller and would still discharge the owner's constraint
   3. All four round-1 files assume `replace` is in v1 without arguing the case. I am not
   confident enough to propose dropping `replace` outright — but it deserves an explicit decision
   in consensus, not an inheritance from the ratified operation kinds. My current position: keep
   `replace` in v1 (the registry and grammar are the same work either way, and dropping it creates
   a false symmetry where the one ratified operation kind is absent), but consensus should record
   that no deck in the survey exercises it.

2. **The one-file constraint (R1.1) as applied to the payload location.** All four participants
   assume the overlay is one file (consensus.md:109-110). But the D1 dispute is partly a
   consequence of that constraint: the payload-location argument would be simpler if the manifest
   and the payload could be separate files (a `.yaml` manifest + a `.md` payload). The one-file
   rule forces the payload into the same file as the manifest, which forces the YAML-literal-scalar
   vs. delimited-body choice. Nobody questioned whether the one-file rule is right for an overlay
   that carries substantial Markdown. I am NOT proposing two files (R1.1 is ratified and the
   one-file rule has integrity benefits — the hash covers everything, and there is no orphan-file
   risk), but the consensus should record that the one-file constraint is the reason the payload
   lives inside the YAML, not a free design choice.

3. **The correlated-agreement risk on D-b (ship the registry now).** All four participants agreed
   the registry ships in this slice. That agreement is partly anchoring: the kickoff prompt
   (00-prompt.md:62-77) frames D-b as "ship it or have no addressing at all," and the scoping brief
   (00-scoping-brief.md:155-157) frames it the same way. The framing is correct — D2 forbids
   heading text and inline markup, and no third option exists in the record — but four agents
   reading the same framing and reaching the same conclusion is not four independent confirmations.
   It is one confirmation reached four times. The §15.6 correlated-agreement duty applies: the
   consensus must record that unanimity on D-b is a shared prior (all four read the same brief),
   not independent evidence, and state what would have to be true for the position to be wrong.
   What would make it wrong: a third addressing option that does not require a registry and does
   not violate D2. No participant found one. That is a finding (the search scope was the ratified
   record and the code), not a failure to comply — but it should be stated as such.

### opencode-1's silence

opencode-1 filed nothing (inbox/claude-1-to-all_..._opencode-1-absent.md). Its silence is not
assent. The inbox note records that the process was killed after 25 minutes with zero bytes of
output, and that this is the second consecutive idea in which opencode-1 produced no artifact. I
do not treat its absence as agreement on any of the four disputes or on any other point. The note
correctly records that the kill is an inconclusive data point (the harness may have cut it short)
and that DF-3 (opencode-1's fitness as a quorum member) is the already-ratified follow-up. Round 2
proceeds with four participants.

## Summary of positions on the four disputes

| Dispute | My round-1 position | My round-2 position | Changed? |
|---|---|---|---|
| D1 (payload location) | YAML frontmatter + free-form body (contradictory) | YAML frontmatter, literal block scalar, no body (codex-1's) | Yes — adopted codex-1 |
| D2 (ext-1 placement) | After §8, before §10, via registry | After §8, before §10, stated as registry offset between block IDs | Refined — registry-dependent phrasing |
| D3 (loss report) | Option (iii): noisy first report + migration note | Synthesis: scoped-(ii) witness + migration note + Applied field | Yes — adopted codex-1/kimi-1 witness + kept migration note |
| D4 (roster annotations) | ext-1 payload | Seventh identity slot (kimi-1's) | Yes — adopted kimi-1 |
