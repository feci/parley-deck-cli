---
idea: meta-protocol-change-review-gate-honesty
drafted-by: claude
date: 2026-06-12
---

## Agreed decisions

Source: round-01 + round-02 (all four ACCEPT; codex's normative blocks are the
base text; agy's Critical/Advisory split was withdrawn in favour of the
objective-finding boundary).

- **D1 — Phase 6 gains "Review briefs and dispositions".** codex round-01 block
  verbatim: briefs MUST NOT suppress findings (no "do not report", "do not
  re-raise", "ignore", severity floors, or scope narrowing); dispositions
  (rebuttals, accepted trade-offs, sandbox artifacts, deferrals, operator
  rulings) MAY travel in briefs using the disposition shape — which gains
  agy's closing "Reviewer Prompt" line ("Please evaluate whether this
  rationale holds under the current scope. Do you concur?"). The reviewer
  independently states concurrence/dissent per disposition in its review file.
  A disputed finding closes only via reviewer withdrawal, normal review
  consensus, or an operator ruling quoted verbatim into the next artifact
  (reuses the existing escalation channel; Phase 7 triage categories
  untouched).
- **D2 — Phase 8 gains "Strict review gate (optional)".** codex block as base:
  `strict_gate: true` in 00-prompt.md frontmatter; default rule (zero Agreed
  fixes) unchanged when absent/false. Under strict gate, zero Agreed fixes is
  necessary but not sufficient — the gate closes only after a FRESH full-scope
  Phase 6 review round of the current reviewed state reports zero findings of
  any severity or kind, recorded by the subsequent Phase 7 consensus.
  Fix-verification/resumed passes converge but never close. NITs, deferrals,
  and accepted-low-severity findings keep the gate open unless withdrawn by
  the reviewer or closed by an explicit operator ruling. Folded in: hermes's
  full-scope definition ("the complete implementation diff at the time of the
  fresh pass — all files changed since the design FINAL plus all fix-up
  commits"), and the objective-finding boundary replacing agy's withdrawn
  Critical/Advisory split: a finding must be an objective, code-grounded issue
  (correctness, security, robustness, maintainability, or factual
  documentation error) in code the reviewer actually read; subjective
  stylistic preference is never a finding at any severity. Mutability: set at
  kickoff by the idea author; afterwards adding/removing/changing requires
  design/review consensus or explicit operator direction recorded in the idea
  — never silent relaxation (codex's rule; hermes withdrew immutability).
- **D3 — Phase 8 gains "Stopping judgment".** codex block verbatim:
  trajectory, not a pass counter — converging (continue), churning (stop +
  escalate with a trajectory summary), blocked (pause that finding's thread
  for the operator; unrelated fixes continue). agy's concrete triggers join as
  ILLUSTRATIVE examples; the qualitative rules are the norm. MaxFixupCycles
  and driver budgets are escalation thresholds, never close criteria.
- **D4 — §8 (inbox/communication) gains "Consults".** Consult artifacts under
  `parley-deck/consults/` are advisory and non-canonical: never round
  artifacts, signoffs, quorum evidence, or dispositions; promoting a consult's
  conclusion into protocol state requires a normal idea/round/consensus
  artifact. (Command mechanics live in the sibling CLI idea.)
- **D5 — Frontmatter + machine readability.** The 00-prompt frontmatter
  documentation gains the `strict_gate: true|false` line (exact
  case-insensitive "true" enables; absent/other = false — the
  ReadAutoImplement precedent). No driver enforcement in this release; a
  future `ReadStrictGate` + machine-readable close fields (e.g.
  strict_gate_clean + closing_review_round in review consensus) are sibling/
  follow-up work and must not scrape prose.
- **D6 — Copies kept in lockstep.** Live `parley-deck/COOPERATION.md`
  (canonical), `internal/protocol/defaults/COOPERATION.md` (embedded default
  for `parley init`), and a dated entry in
  `parley-deck/meta/protocol-changelog.md`. No in-repo references/
  COOPERATION.md exists (verified). The external parley-deck-skill bundled
  snapshot is flagged via an inbox note in the release step.

## Agreed trade-offs

- A strict gate costs at least one extra full-scope review pass — that is its
  purpose; it is opt-in per idea.
- The objective-finding boundary asks reviewers for judgment about
  "subjective style" — accepted; the disposition mechanism and operator
  rulings resolve boundary disputes.

## Open items deferred to implementation

- Exact placement/heading levels within Phase 6/8 and §8 (drafter's choice,
  reviewed in Phase 6 of the sibling implementation).
- The embedded default also lags the live protocol on §12 (known drift, found
  by the graph pass) — syncing ONLY these amendments there is in scope; the
  broader §12 drift is flagged as a follow-up idea.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: claude — 2026-06-12
Status: ✅ ACCEPT
Notes: Base text per codex with the agreed agy/hermes merges; no driver enforcement this release.

### Signoff: codex — 2026-06-12
Status: ✅ ACCEPT
Notes: Verified against codex round-01/round-02; consensus captures the converged review-gate protocol changes.

### Signoff: hermes — 2026-06-12
Status: ✅ ACCEPT
Notes: Faithful record of round-02 convergence.

### Signoff: agy — 2026-06-12
Status: ✅ ACCEPT
Notes: Verified against agy round-01/round-02; the consensus faithfully captures the objective-finding boundary and strict gate convergence.
