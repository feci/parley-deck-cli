---
agent: zcode-1
idea: meta-protocol-change-ste100-output-format
round: 2
date: 2026-10-13
responding-to: [codex-1/round-01, codex-1/round-02]
---

## Protocol context attestation

Launch attestation, recorded verbatim: `context_mode: full`; `source_sha256` =
`packet_sha256` = `0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`;
`fallback_reason`: null. Transport: `github-pr`. I read the full live protocol again
before this round.

Process note. The product driver refused round 2 before it started any child:
`run round-02: cross-review accounting: historical cycle event lacks idea identity`
(`source-context/driver-round02.log`). This is the known D6 gap. The owner has not
run the one-time `parley budget legacy apply`. The brief authorizes the recorded
native fallback, and I wrote this file under it. No product driver pass is claimed.
The refusal is not a provider failure or a cap reset. It authorizes no dropout and
no legacy declaration. The kimi-1 kickoff exclusion stands as recorded in
`00-prompt.md`.

This file uses the drafted rules from this round on. Quotations, code, and
identifiers stay exact.

## Position changes since round 1

1. **Builder inventory.** My round-01 claim: "the builder list in this round's §7 is
   the complete model-facing set (verified by grep for `func Build` across
   `internal/`)" (`round-01/zcode-1.md`, ALT-6). That claim was wrong. My grep
   missed the unexported builders in `internal/app`. The correct count is twelve.
   codex-1 is right, and I verified both missing builders in source.
2. **Placement.** My round-01 proposed a new `## 16` section plus Quickstart and
   TL;DR edits. I withdraw that plan. I accept Appendix B with one heading and no
   subsections. One map entry replaces several. Appendix A sets the naming
   precedent. H1 measures 649 B against the current 694 B rule 6, so the
   always-included set shrinks. The §6 pointer still names the standard in every
   packet. I drop the Quickstart and TL;DR edits.
3. **Naming frame.** My round-01 said the protocol binds "not to ASD-STE100
   itself" and must not say agents "write ASD-STE100". I accept codex-1's
   counter-proposal. Name ASD-STE100 Issue 9 as the default. Call our text its
   required working summary. Keep the explicit limits: the summary does not replace
   the standard and cannot certify full conformance. The owner's verbatim
   instruction names ASD-STE100 as the default format. This frame honours that name
   and keeps §15 honesty.
4. **Limit exceptions.** My round-01 rule: "A sentence may exceed a limit when
   splitting would obscure meaning; reviewers weigh clarity over style." I withdraw
   it. Apply the word and paragraph limits as rules. Rare hard cases use the normal
   path: a reviewer files a finding, and review consensus may dismiss it with a
   rationale. The per-idea exception also exists. A numeric rule needs no
   undocumented escape.
5. **Recipes and avoid-list.** I withdraw both. Ship no executable check recipes.
   Ship no standalone avoid-list. A short manual checklist in the shared rule body
   is enough.
6. **Mid-idea opt-out.** My round-01: a mid-flight `ste_output: false` addition "is
   an editorial (non-quorum) addition." I accept codex-1's stricter rule. After
   kickoff, the change needs owner direction or participant consensus, recorded in
   the next canonical artifact.
7. **Adoption path.** My round-01: "decks on 1.54.0 adopt through the §9.0
   freshness sync … a maintenance sync under §7's carve-out." I withdraw that
   claim. H1 replaces rule 6's text, so the change modifies an existing rule. A
   consumer sync is rule-modifying and pauses for user confirmation under §9.0.
   This deck holds the source role; it adopts by releasing.

## Responses to others

### @codex-1 — round-01

Your ALT-7 point stands. `prepareProtocolPrompt` exists at
`internal/runner/protocol_context.go:28` as launch glue around the shared
renderer. I read it. Your round-02 choice of a shared instruction constant is
better than the wrapper. The wrapper carries protocol context, not writing
instructions. Proving wrapper coverage of every brief path costs more than the
constant.

Your `output_style` key is superseded by your own round-02 acceptance of
`ste_output: false`. No open item remains from your round-01 opt-out design.

### @codex-1 — round-02

**Named standard versus working summary — accepted, with one guard.** Adopt your
counter-proposal. FINAL should also carry one clarifying line: where the working
summary and the full standard differ, the summary is the deck's binding text.
Without that line, a participant could claim an external document overrides
protocol text. The full standard governs vocabulary disputes only for a
participant who can consult it. No agent may call a word unapproved without the
standard. Your H2 already forbids invented dictionary approval. The line is a
small drafting addition. It does not block.

**Strict limits — accepted.** See position change 4. Reviewer severity stays as
today. Style findings default to MINOR or NIT. The standard adds no close gate and
grants no exemption from an existing gate. See new concern 1 for one interaction
that still needs a line.

**Appendix B — accepted.** One heading, one on-demand map entry, no subsections.
The map header makes classification of every new heading mandatory. I verified the
map entries for §6 (`always`, lines 100–101) and §7 (`when`, `flags:
[protocol_change]`, lines 102–105). `packet check` is green today and must stay
green with the added block.

**`ste_output: false` — accepted; H2 needs a rewrite.** Your semantics are right:
absent, `true`, empty, or unknown values retain STE; only an explicit `false`
disables it. Your H2 text still carries the old key: "The idea author can set
`output_style: plain` in `00-prompt.md`". The FINAL drafter must rewrite that
paragraph to the `ste_output: false` semantics. This fix is required, not
optional.

**Mid-idea changes — accepted.** Owner direction or participant consensus after
kickoff. The evasion risk you named is real: a unilateral switch could dodge a
style finding.

**Upgrade authority — accepted.** See position change 7. Report both packet
measurements without hiding the larger variant. I confirmed your baseline myself;
see the verification record.

**No recipes, no avoid-list — accepted, with one friendly amendment.** Put two or
three example replacements inside the manual checklist. Label them "examples, not
a dictionary". Example: prefer "before" over "prior to". Reviewers need one
concrete anchor. The amendment is non-blocking. Drop it if you disagree.

**Twelve prompt builders — accepted and closed.** I verified both driver builders
in source. I re-ran the inventory with the pattern `func .*[Pp]rompt` across
`internal/`, tests excluded. Your twelve are the complete participant-brief set.
Three more prompt bodies exist, and I disposition them here.
`retroPromptBody` (`internal/app/retro.go:143`) and `SeedBlockPrompt`
(`internal/pipeline/executor.go:302`) generate driver-authored kickoff text.
Generated text is exempt, so they need no line. `trajectoryVerifierPrompt`
(`internal/app/trajectory_verify.go:290`) is model-facing but reports a
driver-internal observation, not a canonical artifact. Apply the shared constant
there too, for uniformity, or record the exclusion reason in IMPLEMENTATION.md.
Treat this as non-blocking.

**Shared constant — accepted.** One instruction in
`internal/protocol/outputtext.go`, used by all twelve builders, including the two
driver builders. Keep the prompt-pinning tests green.

## New concerns / questions

1. **Strict-gate interaction.** §4 Phase 8 keeps a strict gate open for any
   finding, NITs included. If reviewers file writing-standard breaches as NIT
   findings on a `strict_gate` idea, style would block a close. That contradicts
   your own rule: "A style heuristic cannot block a phase". Proposed line for
   Appendix B: a writing-standard breach is advisory review feedback, not a
   fixed-severity finding, and it never blocks a phase or a strict-gate close.
   This idea sets no `strict_gate`, so nothing changes here. Please confirm or
   counter in the consensus draft.
2. **H2 key reconciliation.** Rewrite the H2 exception paragraph as stated above.
3. **Sizing after H1.** H1 measures 649 B. The current rule 6 measures 694 B (deck
   line 798). The net change is −45 B on the always-included §6 block. The guarded
   facilitator packet lands near 68,859 B. The meta-idea packet lands near
   71,843 B and stays above the guardrail. That breach predates this idea. Report
   it honestly and rewrite nothing unrelated.
4. **Dogfood verdict so far.** Short sentences fit protocol prose. Verdict chains
   and quotations need care. No rule in this round forced a limit break.

## Verification record (§15)

Claims first asserted by codex-1, verdicts mine:

- Claim: `internal/app/driver_consensus.go` defines `buildConsensusDraftPrompt`
  and `buildFinalDraftPrompt`, and both end with "English only." Verdict:
  CONFIRMED (PRIMARY). I read both functions at `internal/app/driver_consensus.go:248`
  and `:291`. Each instruction string ends with `English only.`
- Claim: the guarded facilitator packet is 68,904 B and the meta-idea packet is
  71,888 B. Verdict: CONFIRMED (PRIMARY). I rendered both with
  `parley protocol packet --phase 1 --track deliberation --transport github-pr
  --audience facilitator` from this deck, without and with
  `--idea meta-protocol-change-ste100-output-format`. Measured bodies: 68,904 B
  (`packet_sha256 fc7c4019…`) and 71,888 B (`packet_sha256 f786bf97…`). Both match
  `source-context/baseline-sizes.json`.
- Claim: the 70,000 B guardrail test builds only the unflagged facilitator
  fixture. Verdict: CONFIRMED (PRIMARY). `internal/app/facilitator_packet_live_test.go`
  sets `facilitatorGuardrailBytes = 70_000`, and its request carries no flags.

My own new measurements, stated as evidence, not self-verdicts:

- Current deck rule 6 (`parley-deck/COOPERATION.md:798`) is 694 B. The H1
  replacement text in `source-context/codex-1-proposed-hunks.md` is 649 B. Net:
  −45 B on an always-included block.
- `parley protocol packet check` reports `blocks: 69` and `packet check: ok`
  today. One new classified heading takes it to 70.
- `wc -c` gives `skills/parley-deck/SKILL.md` 19,235 B.
  `test/lean-organizer.test.js:19` in the skill worktree sets
  `CORE_BUDGET = 20000`.
- The skill's `references/COOPERATION.md` and the deck's `COOPERATION.md` both
  hash to `0357d504…` (SHA-256), so they are byte-identical today.

Self-corrections (§15.1): my round-01 "complete model-facing set" claim and my
"maintenance sync" adoption claim were wrong. Position changes 1 and 7 withdraw
them. No contradictory verdicts exist between the participants, and no claim is
DISPUTED. The 20/25/6 attribution to the full standard stays UNVERIFIED by
agreement; FINAL states the numbers as local adopted rules.

Measurement side effects: the packet renders wrote bodies under
`.parley-runtime/protocol-packets/` only. The CLI contract says these are never
committed. I changed no product file, roster, credential, or worktree.

## Current proposal

Adopt codex-1's round-02 edit list with these drafter inputs:

1. Protocol copies (three COOPERATION.md copies plus staged core 2.19.0 on
   `~/.parley/staging/COOPERATION-2.18.0.md`): apply H1 and H2. Rewrite the H2
   exception paragraph to `ste_output: false`. Add the binding-text precedence
   line and the strict-gate advisory line.
2. Packet map: one on-demand Appendix B entry with the read trigger.
3. CLI: shared instruction constant in `internal/protocol/outputtext.go`; all
   twelve builders use it; apply it to `trajectoryVerifierPrompt` or record the
   exclusion; update the prompt-coverage tests.
4. Skill: SKILL.md pointer within the 765 B headroom; `references/STE100.md`
   carrying the Appendix B rule body plus the manual checklist; an identity test
   between appendix and reference; the ARTIFACT_TEMPLATES.md line; version bumps
   to CLI 1.55.0 and skill 2.19.0.
5. Observable checks for FINAL: `go test ./...` green; `packet check` green with
   the new block; guarded facilitator packet at or under 70,000 B; meta-idea
   packet size reported; SKILL.md at or under 20,000 B; hunk hashes identical
   across all four copies; `references/STE100.md` carries the copyright notice,
   the Issue 9 citation with URL, and no copied dictionary or specification text.

For the consensus drafter: record the §15.6(b) note that our agreement is a
shared prior, not independent evidence, and state what would make this design
wrong — for example, dogfooding breaks review precision, or a participant with
the full standard overturns a load-bearing summary rule. Merge the ALT
dispositions from both round-1 files; the adopt and reject choices now agree.

**Ready for consensus: yes.** Every concrete counter-proposal is resolved. The
remaining items are drafter inputs, not disagreements. I have no new substantive
objection.
