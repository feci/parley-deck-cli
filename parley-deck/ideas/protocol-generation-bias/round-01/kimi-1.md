---
agent: kimi-1
idea: protocol-generation-bias
round: 1
date: 2026-08-28
---

## Summary

The protocol converges on the first plausible frame because nothing makes a second frame exist.
My axis fixes exactly that: before any convergence is permitted, N structurally distinct
candidates must exist on disk, each with a named owner who was structurally obliged to produce
it. The mechanism is a direct port of PDS/1.0's two-part device — deterministic seeded
assignment of distinct positions plus a facilitator-computed distinctness gate between round 1
and round 2 (PRIMARY: `references/PDS.md:316-324,365-376` in the installed `parley-design`
skill — "Diverge in isolation. A proposer MUST NOT read another direction before submitting its
own." / "assignment = rotate(sorted(primary_positions), uint32(sha256(\"PDS/1\" ||
run_id)[0:8]))" / G1 "MUST fail if any pair of directions differs on fewer than two of the
brief's axes … Persistent convergence never auto-passes … short of all three it returns
`ABSTAIN`") — generalized from visual axes to a fixed menu of *generation stances*, plus two
additions PDS does not have: a machine-checkable occupancy floor in frontmatter, and a
`## Dropped candidates` ledger that follows the candidates into later rounds.

- On **B2** (the unproposed option): direct hit. A `REUSE` stance assigns the
  "find the off-the-shelf route" work to a specific agent with an evidence obligation, so the
  `pnpm deploy` class of option is *generated* before convergence, not left to chance.
- On **B1** (the raised-and-lost option): my axis deliberately does less. It guarantees the
  alternative exists and makes its death a recorded, attributable event; it does not stop the
  group from killing it. Changing that outcome is A2/A3 territory, and I name A2 as my
  defection target below.
- **Cost**: net negative. D1+D2 replaces §15.6 in full (~22 lines of rule text that fires on
  none of the failing cases) with ~14 lines that fire on every artifact-class idea, and adds no
  new opt-in flag.

## Proposed approach

### What the port sources actually say (verified this session)

- PDS/1.0 assigns distinct positions deterministically and records them recomputable:
  "Each proposer takes a distinct position on the brief's primary axis by `assignment =
  rotate(sorted(primary_positions), uint32(sha256(\"PDS/1\" || run_id)[0:8]))` … Each DIRECTION
  MUST record what it was given as `assigned`, so the mapping recomputes from the brief and the
  directions alone." (PRIMARY: `references/PDS.md:367-374`)
- PDS/1.0 gates convergence on distinctness, computed by the facilitator with no model call:
  "Between round-01 and round-02; facilitator-computed, no model call. MUST fail if any pair of
  directions differs on fewer than two of the brief's axes … a failed set MUST NOT proceed to
  critique. Remedy: exactly one seeded forced-axis re-diverge." (PRIMARY: `references/PDS.md:316-322`)
- PDS/1.0 isolates proposers: see the "Diverge in isolation" quote above (PRIMARY:
  `references/PDS.md:365-366`).
- The deck already executes the seeded-rotation half cheaply: this very idea assigned six
  distinct generation axes by "agents sorted lexicographically, axes rotated by sha256 of the
  string PDS/1|protocol-generation-bias, first 8 hex chars, mod 6 = 3" (PRIMARY:
  `parley-deck/ideas/protocol-generation-bias/00-prompt.md:7-13`, roles block on disk). The
  rotation is one line of arithmetic. It is already running.
- The deck already executes the isolation half procedurally: my own round-1 dispatch forbids
  reading other agents' round-01 files on pain of invalidating the round (PRIMARY: this
  session's launch envelope for kimi-1).
- What §15.6 currently does instead: it fires only "if round 1 closes with no substantive
  disagreement **and** the idea's output is primarily a judgment rather than a mechanically
  decidable artifact" (PRIMARY: `parley-deck/COOPERATION.md:1346-1350`). Both conjuncts exclude
  the B1/B2 class, and it triggers *after* round 1 — i.e., after the anchor has landed.
- The review vocabulary has no class for "the whole approach is wrong": the concrete-finding
  scanner accepts exactly `case "CRITICAL", "MAJOR", "MINOR", "NIT":` (PRIMARY:
  `internal/driver/impl.go:444-445`).
- The Go review gate already enforces a *verification*-side obligation — a review is rejected
  without "a non-empty '## Refutation attempts' section (refutation-default: a review must
  record what it tried to break)" (PRIMARY: `internal/protocol/reviewartifact.go:41-43`) — and
  it already parses idea/artifact frontmatter (`ReadFrontmatter`, same file, line 18). A
  frontmatter-based generation-side check has an existing enforcement family to join.
- Simplicity vocabulary in the shared rule text: my own count over `COOPERATION.md` for
  `simpler|simplicity|YAGNI|over-engineer|smallest|off-the-shelf|built-in|do nothing` returns
  zero matches (PRIMARY: Grep `count_matches` run this session — "No non-sensitive matches
  found"). §15 has a hard theory of checking and zero vocabulary for *there was a cheaper way*.
- `require_model_diversity:` appears exactly once across the entire `ideas/` tree — this idea's
  own prompt (PRIMARY: Grep `count_matches` over `parley-deck/ideas/` this session, 1 hit).
  Opt-in generation-side machinery is, in this deck, synonymous with absent machinery.

### D1 — Assigned generation stances (the port)

For every idea whose output is declared `artifact_class: mechanical` (or `mixed`) in
`00-prompt.md` frontmatter, round 1 is a **forced-divergence round**:

1. **Stance menu.** The protocol defines four canonical generation stances — structural
   postures toward the problem, not topical positions:
   - `REUSE` — the candidate must be built from an existing, documented, off-the-shelf or
     built-in mechanism; no new machinery. The stance carries an **evidence obligation**: the
     file must list, with §15.2 tags, what sources were searched (vendor docs, package
     managers, man pages, upstream trackers) and which concrete options each source offered.
     A `REUSE` file with only `RECALL` support reads as `UNVERIFIED` by existing §15.2 — the
     current protocol already has the teeth; D1 just points them at generation.
   - `SUBTRACT` — the candidate must achieve the goal by deleting, disabling, or declining
     scope; "do less / do nothing differently" is a legitimate occupant and must be steelmanned,
     not performed.
   - `REPLACE` — a different machine end-to-end; the candidate may not reuse the anchor's
     core mechanism.
   - `MINIMAL` — the smallest diff to the status quo that satisfies the stated need.
2. **Deterministic assignment.** Same formula as PDS §4 rule 2, already dogfooded by this
   idea: sort participants lexicographically, rotate the sorted stance list by
   `sha256(seed || slug)[0:8] mod len(stances)`, map in order. With more participants than
   stances, remaining participants are unassigned (free generation); with fewer, take the
   rotated prefix. Each round-1 file records `assigned:` in frontmatter so the mapping
   recomputes from the brief and the files alone.
3. **Blind round 1.** The existing dispatch discipline (no reading sibling round-1 files)
   becomes part of the mechanism, not a per-launch courtesy. This is the single cheapest
   component and the one with the deepest evidence base (below).
4. **Decline valve.** An assignee may decline its stance with a one-line recorded reason;
   declining does not relax the occupancy floor (PDS's rule, PRIMARY:
   `references/PDS.md:377-378` — "Declining does not relax G1").
5. **Round 2+ is unconstrained.** Amendments, hybrids, defections, convergence — all legal.
   The ceremony is confined to round 1; nothing here slows later rounds.

### D2 — Distinctness floor and the dropped-candidates ledger

1. **Occupancy floor (mechanical).** Before round 2 opens, the facilitator checks: every
   assigned stance has exactly one file, and no two assigned files declare the same stance.
   This is decidable from frontmatter alone — a shell one-liner today, a
   `ValidateRoundArtifact`-style check later, in the same family as the existing review gate
   (PRIMARY for that gate's existence: `reviewartifact.go:17-54`). Unlike §15, which has no
   enforcement surface (SECONDARY: claude-1, 00-prompt.md finding 4), the floor's input is a
   frontmatter key the driver already parses.
2. **Semantic distinctness (judgmental, bounded).** Assigned files must differ structurally,
   not lexically. The only adjudication D2 adds: each round-2 file's existing
   `## Refutation attempts` section must name, per assigned candidate, whether it survives and
   why not. No new judge is created; the burden rides on a section the Go gate already
   requires to be non-empty (PRIMARY: `reviewartifact.go:41-43`).
3. **`## Dropped candidates` ledger.** From round 2 onward, every round file (and `FINAL.md`)
   carries one line per generated candidate: stance, author, round it died in, one-line cause.
   A candidate silently re-absorbed by the anchor is precisely what the ledger makes loud.
   This is the *only* part of my axis that touches post-round-1 artifacts.

### Why this is the right shape — external evidence, tagged

- **Independent generation before interaction beats interaction-first, empirically.**
  Nominal Group Technique: individuals generating ideas alone, pooled afterward, outperform
  interacting groups on quantity and quality (RECALL: Van de Ven, A.H. & Delbecq, A.L., "The
  Effectiveness of Nominal, Delphi, and Interacting Group Decision Making Processes", *Academy
  of Management Journal* 17(4), 1974, 605–621). Production-blocking experiments show the same:
  nominal groups outperform brainstorming groups (RECALL: Diehl, M. & Stroebe, W.,
  "Productivity Loss in Brainstorming Groups: Toward the Solution of a Riddle", *Journal of
  Personality and Social Psychology* 53(3), 1987, 497–509). D1's blind round 1 is NGT's
  written-silent-generation phase with the name filed off.
- **Delphi** institutionalized anonymous, independent, iterated estimates precisely to deny
  anchoring (RECALL: Dalkey, N. & Helmer, O., "An Experimental Application of the Delphi
  Method to the Use of Experts", *Management Science* 9(3), 1963, 458–467).
- **Structured analytic techniques** make the full-hypothesis-set-first move explicit:
  Analysis of Competing Hypotheses requires enumerating alternative hypotheses *before*
  weighing evidence (RECALL: Heuer, R.J. Jr., *Psychology of Intelligence Analysis*, CIA
  Center for the Study of Intelligence, 1999, ch. 8; public PDF on cia.gov). D1's stance menu
  is ACH's hypothesis matrix compressed to four rows.
- **Diversity of solvers has a theoretical basis**, not just a vibe: diverse groups of problem
  solvers can outperform uniformly high-ability groups because distinct heuristics cover
  different parts of the search space (RECALL: Hong, L. & Page, S.E., "Groups of diverse
  problem solvers can outperform groups of high-ability problem solvers", *PNAS* 101(46),
  2004, 16385–16389, DOI 10.1073/pnas.0403723101). Forcing structural diversity is not
  egalitarianism; it is coverage.
- **Without a forced structure, LLM groups drift to agreement.** Five state-of-the-art
  assistants "consistently exhibit sycophancy across four varied free-form text-generation
  tasks", and preference data shows "when a response matches a user's views, it is more likely
  to be preferred" (PRIMARY: Sharma, M. et al., "Towards Understanding Sycophancy in Language
  Models", arXiv:2310.13548, abstract fetched this session). A proposal on the table plays the
  role of "the user's view"; amending it is the path of reward.
- **Reflection does not escape a settled stance** — "once the LLM has established confidence
  in its solutions, it is unable to generate novel thoughts later through reflection even if
  its initial stance is incorrect" (PRIMARY: Liang, T. et al., "Encouraging Divergent Thinking
  in Large Language Models through Multi-Agent Debate", arXiv:2305.19118, abstract fetched
  this session). Their fix is multi-agent debate with a "tit for tat" state and an adaptive
  break — i.e., divergence has to be *maintained by structure*, not requested. Note the
  corollary for D1: divergence confined to round 1 can still re-converge later; that is what
  the ledger is for.
- **Multi-agent debate helps without any of this machinery** — instances "propose and debate
  their individual responses and reasoning processes over multiple rounds to arrive at a
  common final answer", improving reasoning and factuality (PRIMARY: Du, Y. et al.,
  "Improving Factuality and Reasoning in Language Models through Multiagent Debate",
  arXiv:2305.14325, abstract fetched this session). But "a common final answer" from
  undifferentiated proposers is exactly the convergence pattern under indictment; Du-style
  debate samples diversity, it does not *guarantee* it. D1 is the guarantee layer.
- **Self-preference poisons self-adjudication**: LLM evaluators "score their own outputs
  higher than others' while human annotators consider them of equal quality", with
  self-recognition and self-preference linearly correlated (PRIMARY: "LLM Evaluators Recognize
  and Favor Their Own Generations", arXiv:2404.13076, abstract fetched this session). Design
  consequence: D2 must not hand distinctness adjudication to the proposal authors — the floor
  is frontmatter-mechanical, and the semantic half is distributed across rival authors via the
  existing refutation section.
- Sampling-based diversity then marginalization improves reasoning (RECALL: Wang, X. et al.,
  "Self-Consistency Improves Chain of Thought Reasoning in Language Models", ICLR 2023,
  arXiv:2203.11171) — the single-agent analogue of "generate N distinct candidates, then
  select". Diverge-then-converge is also standard design practice (RECALL: UK Design Council,
  "The Double Diamond", designcouncil.org.uk).

### What this deletes (cost accounting, per the ratified net-negative-bytes rule)

- **Deletes §15.6 entirely** (`COOPERATION.md:1346-1367`, ~22 lines incl. both track
  variants). Its judgment-class steelman obligation is subsumed: D1 assigns the
  alternative-generation work up front on the artifact class, and for the judgment class the
  `SUBTRACT`/`REPLACE` stances do the same job with a trigger that actually fires. A rule
  whose two conjuncts both exclude the failing population (PRIMARY: trigger text quoted above)
  is not a safety net; it is a decoration.
- **Shrinks the position-changes theater.** The `## Position changes` apparatus produces
  movement-appearance at high volume (SECONDARY: claude-1, 00-prompt.md finding 8: 1,389
  files, ~40 declaring none). The dropped-candidates ledger records *real* movement — which
  candidate died, when, why — in strictly fewer bytes. I propose the ledger replace, not
  supplement, the per-round position-changes section on D1-activated ideas.
- **Adds no opt-in flag.** D1 binds by `artifact_class:` declaration, defaulting to on for
  `mechanical`/`mixed`. An opt-in version of this mechanism would be the deck's fifth unused
  flag (SECONDARY: claude-1, 00-prompt.md finding 6; partially confirmed PRIMARY: zero
  `require_model_diversity:` uses before this idea).
- Net: roughly −8 to −15 lines of shared rule text, plus deletion of a conditional that never
  fires. New runtime cost: one frontmatter key, one line of hash arithmetic, one shell-checkable
  floor.

### What my axis does on B1

B1: an agent generated the frame-breaking option in round 2, withdrew his own proposal, and
`FINAL.md` shipped the round-01 design anyway (SECONDARY: claude-1, 00-prompt.md benchmarks
section, quoting `FINAL.md:18`; the idea directory is outside this deck root — Glob for
`**/2026-08-14T12-41-49-daily-backup-str/**` under this repository returned no matches this
session, so I could not independently verify the quote).

D1/D2 on B1, honestly: the Proxmox-native-S3 option would have had a named owner in round 1
(`REUSE`: "the backup backend itself may already do this") instead of emerging late and
ownerless; and when it died anyway, the `## Dropped candidates` ledger in `FINAL.md` would
carry the line — stance, author, round, cause — making the frame-restoration a recorded
decision with an attributable rationale instead of an invisible absorption. **That is all.**
If the group decides the anchor wins, D1/D2 lose. B1's fatal step was that the alternative had
no route to a destination once raised; existence was necessary but nowhere near sufficient.
This is the boundary of my axis, and I hold it deliberately: bolting adjudication onto D1
would re-create §15's provenance problems (self-preference, PRIMARY: arXiv:2404.13076 above)
inside a generation mechanism.

### What my axis does on B2

B2 is the case my axis exists for. The failure was not adjudication but *existence*: nobody's
job description included "go look for the documented off-the-shelf route". Under D1, the
`REUSE` assignee in the bundling idea would have been structurally obliged — before seeing
anyone's hand-rolled script — to search the package-manager's own deployment/documentation
surface and file what it found, with §15.2 tags, so `pnpm deploy` (or a documented equivalent,
or a *tagged* null result naming the search scope) enters round 1 as a candidate with an
owner. The obligation cannot guarantee discovery — an obligation to look is not the skill of
finding — but it converts "nobody looked" into "someone looked and reported", which is a
strictly better failure mode and a checkable artifact. Note the compounding detail: with blind
round 1, the `REUSE` assignee searches *uncontaminated* by the anchor's frame, which is
exactly the condition NGT/Delphi evidence says matters (RECALL locators above).

## Concerns / open questions

- **Brief-agreement disclosure (required by the brief's warning).** I agree with the brief's
  core framing, and I say so explicitly: I verified findings 1, 2, 3, the Go half of 5, and 7
  firsthand this session (PRIMARY citations above), and they held up — finding 7 in fact
  understates PDS, which *also* has the isolation rule and the decline valve. Findings 4, 6's
  other counts, 8, and the B1/B2 quotations I took as SECONDARY from claude-1. Where I push
  back on the framing: the brief says the protocol needs to "judge whether a better and
  simpler solution exists", but D1 bets that *generation*, not *judgment*, is the tractable
  half — a stance menu outsources the judgment to coverage. That bet could be wrong; if LLM
  groups can be taught the judgment directly, D1 is scaffolding.
- **The stance menu is itself an anchor.** Four canonical stances decide in advance which
  structures count as distinct. The option B2's critic needed happened to sit inside `REUSE`,
  but the next one may sit outside the menu entirely — a fifth structure the menu trained
  everyone not to imagine. Mitigations considered: a free-axis declaration instead of a menu
  (rejected: uncheckable floor, invites six near-identical "different" candidates); a wildcard
  `FIFTH` stance rotating each idea (attractive; adds entropy but breaks the byte budget
  slightly and needs a rotation rule). Open question, flagged for round 2.
- **Classing rule.** Who declares `artifact_class:`? Proposal: the drafter declares it in
  `00-prompt.md`, challengeable like any other frontmatter; `mixed` and undeclared default to
  D1-on, because the failure mode D1 targets (a mechanically checkable artifact treated as
  judgment) is the expensive direction of error. Not yet stress-tested against the deck's
  idea corpus.
- **Obligation ≠ competence.** The `REUSE` assignee may search shallowly and file a tagged
  null result while `pnpm deploy` sits one doc page deeper. The evidence obligation makes the
  shallowness *visible* (search scope must be named), but does not cure it. A4's standing
  red-team could own stance quality; that is composition, not part of D1.
- **Does G1 work in practice?** PDS's distinctness gate exists and is enforced in JS
  (SECONDARY: claude-1, 00-prompt.md finding 7, "enforced in JS"), but I verified no usage
  data this session — I do not know how often G1 fires, how often re-diverge succeeds, or
  whether facilitators route around it. Porting a mechanism with unknown operational history
  is a measured risk, not a proven one.
- **Interaction with the closed review vocabulary.** D2's ledger creates the record, but a
  reviewer who wants to say "this dropped candidate should have won" still has no finding
  class for it (PRIMARY: `impl.go:444-445`). That is A2's hole and A2 should fix it; D1/D2 are
  designed to compose with exactly that repair.

## Risks

- **Six costumes, one frame.** Assignees may occupy stances nominally while writing the
  anchor's proposal in four dialects. The occupancy floor only checks labels; semantic
  distinctness is adjudicated by the same population that converged. Assigned dissent is
  measurably weaker than authentic dissent (RECALL: Nemeth, C.J., Brown, K. & Rogers, J.,
  "Devil's advocate versus authentic dissent: Stimulating quantity and quality", *European
  Journal of Social Psychology* 31(6), 2001, 707–720), and poorly executed assigned
  alternatives can *increase* confidence in the anchor by giving it strawmen to defeat —
  the intelligence-tradecraft literature warns about exactly this failure of pro-forma
  devil's advocacy (RECALL: *A Tradecraft Primer: Structured Analytic Techniques for
  Improving Intelligence Analysis*, CIA, 2009). Partial mitigation inside my axis: the
  `REUSE` evidence obligation makes at least one stance's shallowness mechanically visible.
- **Re-convergence after round 1.** Divergence confined to round 1 decays: debate participants
  settle, and once confidence sets in, reflection stops producing novel thoughts (PRIMARY:
  arXiv:2305.19118, quoted above). PDS answers this with a *between-rounds* gate and exactly
  one seeded re-diverge (PRIMARY: `references/PDS.md:316-324`); D2 keeps only the ledger,
  dropping the re-diverge to save ceremony. If the deck finds dropped candidates dying
  silently *despite* the ledger, the seeded re-diverge is the piece to re-import — at a byte
  cost I would then defend separately.
- **Judgment-axis gaming.** What counts as "structurally distinct" beyond stance labels is a
  judgment call, and LLM judges self-prefer (PRIMARY: arXiv:2404.13076). Keeping the floor
  mechanical and pushing semantics into rival-authored refutation sections bounds this but
  does not eliminate it.
- **Token cost is real though bounded.** Four owned candidates per idea instead of one anchor
  plus amendments: roughly 4× round-1 generation cost on D1-activated ideas, zero added cost
  afterward. If the deck's corpus shows most ideas never needed a second frame, A6's null
  position outperforms D1 on cost — I take that seriously; the corpus measurement in finding 8
  (SECONDARY: claude-1, 8 of 10 frame reversals human-sourced) is the strongest evidence
  *for* paying the cost, not against it.
- **Adoption risk if this degrades to opt-in.** The deck's opt-in generation machinery sits
  unused (PRIMARY: `require_model_diversity:` count of 1 across all ideas — this idea itself).
  If round 2 waters D1 down to a flag, it is already dead; in that world the honest output is
  A6's deletion-plus-discipline position, not a seventh unused flag.
- **Where my axis fails, and my defection target.** A1 fails outright on B1-class failures:
  when the alternative exists and the group kills it anyway, a generation mechanism has done
  all it can; only a route-and-destination repair changes that outcome. It also fails when the
  structure space is not menu-enumerable (the anchor concern above). If forced to defect, I
  defect to **A2-reframe-vocabulary**: B1's decisive moment was an existing alternative with
  no finding class, no route, and no destination, and that defect survives even a perfect
  generation mechanism. A1 without A2 produces well-documented losses; A2 without A1 produces
  well-routed monotony. They compose, and if only one ships, which one should depend on
  whether the deck weights B1 or B2 higher — I weight B2, because generation failures are
  invisible in the record while adjudication failures at least leave a corpse.
