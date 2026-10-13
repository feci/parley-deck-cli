---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
round: 1
date: 2026-10-13
---

## Protocol attestation

Official-launch context was supplied verbatim at launch: `context_mode: full`,
`source_sha256` = `packet_sha256` = `0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`,
no `fallback_reason`. I re-read the full live protocol supplied in the packet. The
launch context also carried an unapplied shadow-packet diagnostic
(`packet_sha256` `02b0030f…`, 36 included / 33 omitted blocks); it describes no
supplied protocol and I did not rely on it. Transport: `github-pr`.

Independence: I have not read `round-01/codex-1.md` or any other participant
round-1 answer. I read `00-prompt.md`, `source-context/ORGANIZER-BRIEF.md`,
`source-context/ste-source-notes.md`, and the source trees named below. Round 1 of
this idea is exempt from the drafted rules by the kickoff constraint ("use the
drafted rules in this idea from round 2"); I still kept sentences short where that
cost nothing.

## Summary

Make the protocol bind agents to a local, ASD-STE100-derived writing standard for
newly authored English protocol text, not to ASD-STE100 itself. Put the binding
rule in a new on-demand protocol section (§16) plus a one-sentence pointer at the
§6.6 English-only rule, and put the full local rule text in a new skill on-demand
reference `references/STE100.md`. Default on for artifacts written after a deck
adopts the change, per-idea opt-out `ste_output: false`, no rewrite of existing
text, no shipped linter (copy-paste check recipes only), one added prompt line in
each model-facing brief builder. This keeps the packet and SKILL.md inside their
byte budgets, keeps the copyrighted specification out of the repository, and makes
no conformance claim the deck cannot verify.

## Proposed approach

### 1. Frame: a derived local standard, not an ASD-STE100 conformance claim

The protocol adopts **"the deck writing standard (ASD-STE100-derived)"**. Its
authoritative text is the protocol section and skill reference we write; ASD-STE100
Issue 9 (January 2025) is the named, cited inspiration. The protocol must NOT
claim that agents "write ASD-STE100" or "conform to" it: the specification and its
dictionary are copyrighted, the dictionary has not been inspected by any
participant, and a conformance claim would be exactly the kind of assertion §15
forbids presenting as established. This framing honours the owner's direction (the
name, the rules, the default) while staying inside copyright, §15, and the size
limits. The organizer brief itself already forces this shape ("summarize the
writing rules in your own words, cite the source").

### 2. The core rule set (local text, own words)

Binding rules for newly authored English protocol text (rounds, reviews, consensus,
FINAL, IMPLEMENTATION, inbox, PR/MR mirrors, commit messages, generated briefs):

1. **Sentence length.** At most 20 words in a procedural sentence; at most 25 words
   in a descriptive sentence. These are the deck's own limits, chosen to match the
   limits commonly attributed to ASD-STE100 (that attribution is UNVERIFIED — see
   the verification record; nothing material rests on it).
2. **Paragraph length.** At most 6 sentences per paragraph.
3. **One instruction per sentence** in procedures; descriptive sentences carry one
   main statement.
4. **Imperative mood** for procedures ("Write the file.", not "The file should be
   written.").
5. **Active voice by default.** Permitted exception, matching the standard's own
   public description: passive voice in descriptive text when the actor is unknown.
6. **One meaning per word.** Do not vary words for style when one approved word
   serves; avoid words with several meanings where an unambiguous word exists. A
   short avoid-list of common ambiguous or inflated words lives in the reference
   (our own words, ~10 entries, e.g. prefer "before" over "prior to", "use" over
   "utilize", "if" over "in the event of").
7. **Technical names and technical verbs are permitted** (identifiers, file paths,
   product names, domain terms), including `-ing` forms that serve as technical
   nouns. The rule bans `-ing` forms used as verbs, not the suffix itself.
8. **Simple verb forms.** Prefer the simple present for procedures and general
   statements; use the simple past only for completed events. Avoid progressive
   and auxiliary chains in protocol prose.

Limits are **targets enforced by reviewer citation, not machine gates**. A sentence
may exceed a limit when splitting would obscure meaning; reviewers weigh clarity
over style, exactly as they weigh severity today.

### 3. Placement, inside the size limits

- **Protocol — new `## 16. Writing standard for protocol output
  (ASD-STE100-derived)`**, ~1.5–2.5 KB, in `parley-deck/COOPERATION.md`. Content:
  scope (newly authored English protocol text), the eight rules above in compact
  form, the exemption list, the default and `ste_output: false` opt-out, the
  adoption rule, the copyright note, and the Issue 9 citation. Optionally
  `### 16.1 Rules` / `### 16.2 Exemptions, opt-out and adoption` subsections.
- **Packet map** — `parley-deck/meta/packet-applicability.yaml`: one new entry per
  new heading, `include: on-demand`, trigger "authoring, reviewing, or verifying
  any English protocol artifact or generated brief". Precedent: the §13.x
  on-demand entries (map lines 186–201). The map header (lines 14–17) already
  warns that every heading must be classified exactly once, so these entries are
  mandatory, not optional. `on-demand` (not `when`, as §15 uses) keeps the
  facilitator packet's retained-source total under its 70,000 B guardrail
  (`internal/app/facilitator_packet_live_test.go:21`); the measured headroom is
  only ~1.1 KB and a `when`-all-phases section would breach it.
- **Pointer in §6 rule 6** (English-only, `parley-deck/COOPERATION.md` line 798;
  the §6 block is `always` in the map, lines 100–101): one appended sentence —
  English output follows the deck writing standard (§16); code, identifiers,
  quotations, frontmatter and generated text are exempt; per-idea opt-out
  `ste_output: false`. Budget ~250–300 B.
- **Quickstart** reference-appendices list: append "§16" (~30 B, optional).
- **TL;DR (§10)**: one new item (~120 B; §10 sits in the ratified facilitator
  omission set, so this costs the facilitator packet nothing).
- **Skill** — new `skills/parley-deck/references/STE100.md` (~4–6 KB, no byte
  budget applies to references; `ROSTER_AND_PROTOCOL.md` is 39,143 B today): the
  full local rule text, the avoid-list, the check recipes, the copyright notice,
  the citation. `SKILL.md` gains one line near the Core Rule and one bullet in
  "References (on demand)" — together ≤ ~300 B against the measured 765 B headroom
  of the 20,000 B core budget (`test/lean-organizer.test.js:19`, SKILL.md measured
  19,235 B). `references/ARTIFACT_TEMPLATES.md` gains the same participant-prompt
  line the Go builder gains (it mirrors `BuildRoundOnePrompt` today).

### 4. Exemptions (never rewritten, never flagged)

Verbatim quotations of the owner (including the Slovak quotes), of other
participants, and of external sources; code, commands, file paths, identifiers,
API and schema keys, log and error excerpts; YAML frontmatter keys and values;
mechanically generated or validator text (driver digests, phase handoffs, round
index, wait envelopes, telemetry); other agents' text quoted inside a round or
review; proper nouns and product names. Relationship to §6.6: the English-only
rule governs what language protocol files use; the writing standard governs how
the English is written. Owner chat in another language is outside the protocol
file rule entirely — state that explicitly in §16.

### 5. Default, opt-out, adoption

- **Default on** for every English protocol artifact authored after the deck
  adopts the changed protocol. Absence of any key means on.
- **Per-idea opt-out**: `ste_output: false` in `00-prompt.md` frontmatter, set by
  the idea author; adding it to an open idea mid-flight is an editorial (non-quorum)
  addition, noted in that author's next artifact. Boolean-with-false mirrors the
  `quota_auto_exclude: false` precedent in the Phase 0 template. No deck-level
  opt-out in this version — the layered `[defaults]` mechanism exists if a later
  idea wants it; do not build it unasked. No code reads the key in 1.55.0: like
  `roles:`, it is metadata that agents and reviewers honour; that is the right cost
  for a style rule (accepted trade-off, recorded).
- **Older decks / open ideas**: decks on 1.54.0 adopt through the §9.0 freshness
  sync; the change is additive and compatible, so it is a maintenance sync under
  §7's carve-out, not a per-deck meta-idea. Open ideas at sync time: STE applies to
  their next-authored artifact; nothing already written is rewritten. Gradual
  re-authoring of historical artifacts: recommended **against** — deferred
  indefinitely; raw history is testimony, and rewriting it would damage the audit
  trail (§15 values verbatim provenance over style).

### 6. Mechanical check: recommendation — no shipped linter; recipes only

A shipped checker (CLI subcommand or skill script) does **not** earn its cost in
1.55.0: dictionary checking is copyright-barred, passive-voice heuristics are
noisy, nothing gates on the result, and shipped code drags tests, docs, and
release-channel coupling for an advisory nicety — plus it invites gate creep.
What does earn its cost: two copy-paste recipes inside `references/STE100.md`
(sentence word-count; paragraph sentence-count), written to skip code fences,
frontmatter, and blockquotes, plus the ~10-entry avoid-list. Reviewers who want
citations run a recipe and cite real violations as MINOR/NIT. Revisit a committed
script only if dogfooding shows reviewers use the recipe and need reproducibility.

### 7. Prompt coverage (briefs and templates)

One standard line (~45 words, same wording everywhere) added to each model-facing
prompt builder in the CLI:

- `internal/runner/runner.go` — `BuildRoundOnePrompt` (line 959), `BuildRoundPrompt` (line 1111)
- `internal/runner/phase58.go` — `BuildFixupPrompt` (238), `BuildImplementationPrompt` (253), `BuildReviewPrompt` (286), `BuildReviewConsensusPrompt` (406)
- `internal/runner/consult.go` — `BuildConsultPrompt` (60), `BuildGoalCheckPrompt` (78)
- `internal/runner/steer.go` — `BuildSteerPrompt` (296)
- `internal/app/consensus_request_signoffs.go` — `buildConsensusSignoffPrompt` (832)

Driver-mechanical surfaces (`BuildRoundIndex`, phase handoffs, digests) get no
line: their text is generated and exempt. The line names §16 and the exemption
classes, so a headless agent that never opens COOPERATION.md still writes to the
standard. Prompt-pinning tests are updated wherever they assert these builders'
bytes (implementer greps the builder names across `*_test.go`).

### 8. Copyright handling

`references/STE100.md` and §16 carry a notice: rules summarized in our own words
from public descriptions of the standard; ASD-STE100 Issue 9 (January 2025) cited
by name and official URL (www.asd-ste100.org); the specification and dictionary are
copyrighted — do not copy or redistribute them; obtain the full standard through
the official request form; no certification or conformance claim is made or
implied. The avoid-list and recipes are our own words and contain no dictionary
content.

### 9. Verification record (§15)

- **Claim: the current ASD-STE100 issue is 9, released January 2025; issue 10 is
  scheduled for January 2028.** Verdict: **CONFIRMED (PRIMARY)**. I fetched
  https://www.asd-ste100.org/STE_faq.html on 2026-10-13; it states verbatim: "The
  current issue is 9, released in January 2025. The next issue 10 is scheduled for
  January 2028." The downloads page heading reads "Request a copy of ASD-STE100
  (Issue 9, January 2025)" with "Request a free official copy of the current issue
  of the ASD-STE100 standard." — how agents obtain the full specification.
- **Claim: the standard has two parts — writing rules and a controlled dictionary
  of ~900 approved words; technical nouns/verbs are permitted subject-field terms;
  passive voice is permitted in descriptive text when the actor is unknown.**
  Verdict: **CONFIRMED (PRIMARY)**, same fetches; FAQ verbatim: "STE has two parts:
  a set of writing rules (part 1) and a controlled dictionary (part 2)" and
  "approximately 900 approved words"; "We can use the passive voice only if the
  agent (the person or thing that does the action) is unknown".
- **Claim: the 20/25-word and 6-sentence limits are ASD-STE100's actual rules.**
  Verdict: **UNVERIFIED (RECALL)** — no participant has inspected the
  specification's rule text; the public pages I fetched do not state numeric
  limits. Mitigation: FINAL adopts them as the derived standard's own limits, so
  no acceptance criterion depends on the attribution. If a reviewer with spec
  access later confirms the numbers, upgrade the claim then.
- `source-context/ste-source-notes.md` is testimony from codex-1; my own fetches
  independently corroborate every point it makes that FINAL will rely on (issue
  number, two-part structure, technical nouns/verbs, the passive-voice exception,
  the `-ing` technical-noun nuance).
- No claim in this file is owned-and-verdicted by me; the verdicts above are on
  claims whose first canonical appearance is here, and the decisive basis is my own
  fetch, not the organizer's notes.

### 10. Exact changed files and protocol hunks for FINAL

**Protocol hunks** (identical bytes applied to all four targets: deck
`parley-deck/COOPERATION.md`; CLI `internal/protocol/defaults/COOPERATION.md`
(template placeholders preserved); skill `skills/parley-deck/references/COOPERATION.md`
(byte-identical to deck today, sha `0357d504…`); and core 2.19.0 built on
`~/.parley/staging/COOPERATION-2.18.0.md` plus exactly these hunks):

1. §6 rule 6 (deck line 798): append the pointer sentence (~250–300 B).
2. New `## 16` section (+ optional two `###` subsections), ~1.5–2.5 KB, text per
   this round.
3. Quickstart core-vs-reference list: add "§16" (~30 B).
4. §10 TL;DR: one new item (~120 B).

**Packet map**: `parley-deck/meta/packet-applicability.yaml` — on-demand entry with
trigger for §16 and each new subsection (map is protocol; the entries ride this
same idea under §7).

**CLI (Go)**: the ten prompt builders listed in §7 of this round (one line each);
their prompt-pinning tests; one new table test asserting the line appears in every
built prompt. `CHANGELOG.md`; `VERSION` → 1.55.0. No validator change
(`internal/protocol/participantartifact.go` keeps its five required sections; STE
compliance is review-cited, not machine-gated). No `docs/cli-reference.md` change
(no new flag or subcommand).

**Skill**: `SKILL.md` (+~300 B, budget-checked by `test/lean-organizer.test.js`);
new `references/STE100.md`; `references/COOPERATION.md` (synced hunk copy);
`references/ARTIFACT_TEMPLATES.md` (+1 line in the participant prompt shape);
`references/compatibility.json` (`skillVersion` 2.19.0, `recommendedCli`
">=1.55.0"); `package.json` version; `CHANGELOG.md`.

**Acceptance criteria FINAL should state observably**: (a) `go test ./...` green
including the new prompt-line test and the facilitator-packet guardrail with §16
classified on-demand; (b) `parley protocol packet check` green (every new heading
classified); (c) measured byte counts: SKILL.md ≤ 20,000 B, facilitator packet
retained total ≤ 70,000 B; (d) the four protocol copies carry byte-identical hunks
(hash comparison); (e) skill tests green including the 20,000 B core-budget test;
(f) `references/STE100.md` contains the copyright notice, the Issue 9 citation
with URL, and no copied dictionary or specification passage.

## Existing alternatives

§15.6a duty — the mechanisms this proposal builds by hand, and the closest thing
the toolchain already ships, with locators. Sources consulted for the nulls: full
live COOPERATION.md (launch packet), the skill tree listing
(`skills/parley-deck/**`), the CLI Go tree (`internal/**`), `go.mod`, and the
skill `package.json`.

- **ALT-1 — a writing standard expressed as protocol prose.** Hand-built: the §16
  text. Closest shipped: the §6.6 English-only rule itself
  (`parley-deck/COOPERATION.md:798`) — same mechanism (a §6 rule), narrower scope.
  No shipped plain-English/STE style rule exists anywhere in the tree (scoped
  null). Constraint-forced: the rule must bind decks, so it must live in protocol
  text.
- **ALT-2 — an on-demand detailed reference file.** Hand-built:
  `references/STE100.md`. Closest shipped: the entire "References (on demand)"
  pattern (`skills/parley-deck/SKILL.md:259-272`; `references/ROSTER_AND_PROTOCOL.md`
  at 39,143 B proves reference files carry long rule text without a budget).
  Inherited pattern, new file — the right reuse.
- **ALT-3 — packet classification of new headings.** Hand-built: the map entries.
  Closest shipped: `include: on-demand` with a trigger
  (`parley-deck/meta/packet-applicability.yaml:186-201`, §13.x; Appendix A at
  line 131). Constraint-forced: the map header (lines 14–17) makes classification
  of every new heading mandatory or `packet check` fails. `on-demand` chosen over
  §15-style `when`-all-phases on size (would breach the 70,000 B facilitator
  guardrail, `internal/app/facilitator_packet_live_test.go:21`).
- **ALT-4 — a per-idea opt-out key.** Hand-built: `ste_output: false`. Closest
  shipped: optional `00-prompt.md` frontmatter booleans (`strict_gate`,
  `quota_auto_exclude` — §4 Phase 0 template) and `roles:` as the advisory-metadata
  precedent. Inherited mechanism, new key; no CLI change.
- **ALT-5 — a mechanical style check.** Hand-built: two shell/python recipes in the
  reference. Closest shipped: structural artifact validation
  (`internal/protocol/participantartifact.go:80-104` — heading presence and
  non-emptiness, not language) and the LE-4 `checks:` runner (§4.0.1 — runs the
  idea's own commands). No text-lint dependency exists in either tree (scoped
  null: `go.mod` has none — Go stdlib only for this purpose — and skill
  `package.json` declares no linter). External linters (Vale, write-good, textlint)
  exist upstream but none is in the tree and the dictionary-based ones are
  copyright-barred; rejected as new dependencies for an advisory check.
- **ALT-6 — injecting the standard into headless briefs.** Hand-built: the one
  prompt line. Closest shipped: the §15.6a instruction already injected verbatim
  into `BuildRoundOnePrompt` (`internal/runner/runner.go:972-974`). Inherited
  injection pattern; the builder list in this round's §7 is the complete
  model-facing set (verified by grep for `func Build` across `internal/`).
- **ALT-7 — where the binding rule lives.** (a) inline expansion of §6.6 —
  rejected: §6 is `always`-included (map lines 100–101), ~2 KB would breach packet
  headroom; (b) §16 on-demand + skill reference — adopted; (c) rules only in the
  skill reference with a bare protocol pointer — rejected: the protocol is the
  deck's canonical authority and must be self-contained (files-canonical; §6 rule
  4's copy-inclusion principle), so the binding text lives in the protocol and the
  skill reference expands it.
- **ALT-8 — enforcement shape.** (a) validator/CLI gate on sentence length —
  rejected: style is not a §15 verdict regime, and a machine gate invites creep
  plus false positives on exempt text; (b) shipped `parley ste check` subcommand —
  rejected for 1.55.0 (cost without a gate); (c) review-cited targets with
  recipes — adopted.

## Concerns / open questions

- **Packet headroom is the binding constraint.** The brief's ~68,846 B facilitator
  packet figure is codex-1's testimony; my always-block additions are budgeted at
  ~400–450 B total against ~1,154 B headroom. The implementer must re-measure after
  drafting the exact hunk bytes; if over, drop the Quickstart line first, then
  trim the §6.6 pointer to one clause.
- **`on-demand` for an always-applicable rule.** §16 applies to all writing, yet
  packet mode omits it behind a trigger. This is honest only because the §6.6
  pointer (always included) names the standard and where to read it. If reviewers
  think that is too weak, the fallback is a `when`-all-phases entry — but then the
  section must shrink to fit the headroom, which argues for keeping §16 ≤ ~1 KB.
- **The opt-out is social, not mechanical.** No code reads `ste_output:` in
  1.55.0. Accepted (roles: precedent), but FINAL should say so plainly so nobody
  later claims the CLI enforces it.
- **STE strictness vs. review precision.** DISPUTED-claim explanations and verdict
  chains sometimes need longer sentences. The descriptive 25-word limit plus the
  clarity-over-style reviewer judgment covers this, but reviewers should watch for
  compressed nuance during dogfooding from round 2.
- **Owner intent.** My derived-standard frame names and follows ASD-STE100 but
  disclaims certification-grade conformance. The organizer brief's constraints
  (own words, no copying) already imply this frame, so I do not read it as an open
  owner decision; the final handoff note should still state it in one line so the
  owner can redirect if he wanted literal conformance.
- **Map-entry wording.** The trigger string for §16 should name "authoring,
  reviewing, or verifying any English protocol artifact or generated brief"; I
  defer final wording to the drafter so map style stays consistent.

## Risks

- **Byte-budget regression** (SKILL.md 765 B headroom; facilitator packet ~1.1 KB):
  mitigated by on-demand classification, per-hunk byte budgets in this round, and
  observable AC (measured counts) in FINAL; the existing tests
  (`lean-organizer.test.js`, `facilitator_packet_live_test.go`) fail loudly if
  breached.
- **Gate creep** — style findings drifting into blockers: mitigated by no machine
  gate, MINOR/NIT default severity, and the clarity-over-style rule; the strict
  review gate is untouched (style is never an objective code-grounded finding
  under `strict_gate`'s own definition).
- **Copyright exposure** — someone pastes dictionary entries or long spec passages
  into the reference or a round file: mitigated by the notice, the own-words rule,
  and review; a finding would be MAJOR and the text removed before release.
- **Unverified attribution** (20/25/6 as "STE's real limits"): contained — FINAL
  states them as the derived standard's own limits; the UNVERIFIED verdict is
  recorded here so nobody strengthens the claim later without provenance.
- **Correlated adoption friction** — every model's briefs grow by one line and
  must write in a constrained register: negligible per prompt; the real risk is
  slower, stiffer round prose, watched via dogfooding from round 2.
- **Mid-flight idea disruption**: none — no rewrite duty, opt-out available,
  adoption is artifact-granular.
- **Map/protocol drift across the four copies**: mitigated by AC (d) hash
  comparison; the drift guard and `packet check` catch divergence in decks.
