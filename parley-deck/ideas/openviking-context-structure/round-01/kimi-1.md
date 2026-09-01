---
agent: kimi-1
idea: openviking-context-structure
round: 1
date: 2026-08-31
---

## Summary

Verdict: **adopt exactly one narrow piece — the L0/L1 tiering pattern — and we already
half-own it.** The runner already writes a deterministic per-round index
(`round-NN/_index.md`, `internal/runner/round_index.go:96`, frontmatter `derived: true`),
and prompt assembly already knows to skip it (`internal/runner/phase58.go:298`). The gap
is not generation but **consumption**: nothing reads the index to shrink a downstream
prompt. The cheapest real change is to make round N+1 prompt assembly index-first with
on-demand full reads, plus a deck-level `INDEX.md` derived from existing frontmatter.

Do **not** adopt: the `viking://` URI scheme (paths + git already address; OpenViking's
own `md5(account_id:uri)` id is *not* rename-stable — rename changes the uri, hence the
hash — notes §4, PRIMARY), the vector store (dependency constraint; `grep` is the
platform answer), and LLM-written semantic `.abstract.md` files (gate behind measurement
below; nothing ships today that writes them deterministically).

Everything claimed about this repo below is PRIMARY (local filesystem / source, checked
2026-08-31). OpenViking-internals claims carry the notes file's tags.

## Proposed approach

Three pressures, three candidate mechanisms, each with cost, measurement, failure mode.
My lens is evidence: nothing here should be adopted without the replay numbers in M1.

**C1 — Consume the existing round index at prompt-assembly time (L0/L1 for rounds).**
- Mechanism: in `gatherPriorRounds` (`internal/runner/phase58.go:295-301`), inline each
  prior round's `_index.md` instead of full artifacts; attach full artifacts only for the
  immediately preceding round, or list them as `path + approx-tokens` for the agent to
  read on demand. No new file format: `_index.md` already carries per-agent status,
  approx tokens (bytes/4 heuristic, `round_index.go:251`), H2 section names, and a
  180-char first-paragraph excerpt per section (`trimSummary`, `round_index.go:258`).
- Cost: near zero. Generation already happens (runner-owned, deterministic, no LLM).
  Change is ~tens of lines in assembly + tests. No protocol change: artifact formats and
  the round flow are untouched; only what gets *inlined into a prompt* changes.
- Measurement (M1, pre-adoption gate): a read-only replay harness (same shape as
  `parley retro scan`) that, per closed idea, computes the bytes each round-N+1
  participant prompt would carry under (a) current full-inline vs (b) index-first.
  Run over the top-quartile expensive ideas. Adopt only if median prompt reduction
  ≥50% on those ideas. Baseline is checkable today: `protocol-generation-bias` is 30
  files / 614,718 bytes total (PRIMARY, measured; the prompt's "22 artifacts / ~491 kB"
  understates it — `reference/` alone adds ~85 kB), round-01 artifacts run 12–30 kB each.
- Failure mode: an agent answers from the index excerpt and misses nuance in the full
  artifact. Mitigation: index is already labeled `derived: true`; keep the rule that
  signoffs/consensus cite artifact paths, never index lines. Measure the failure mode
  directly (M2): for 3 closed ideas, give a fresh agent only the indexes and ask fixed
  probe questions ("which finding was dismissed?", "what was the verdict on X?") whose
  answers are in the full artifacts; require ≥90% answerability before rollout.

**C2 — Deck-level `ideas/INDEX.md` (cross-idea recall, L0 for the deck).**
- Mechanism: a `parley context deck-index` command (or `parley retro scan`-adjacent)
  that regenerates one `INDEX.md` from existing `00-prompt.md` frontmatter (slug, title,
  track, status, participants, date) — pure parse, deterministic, fail-closed, additive.
  No history rewrite; old ideas index fine because frontmatter is already mandatory.
- Cost: zero tokens, one new read-only-ish command (~150 lines, mirrors `internal/retro`).
- Measurement (M3): pick 10 historical "have we deliberated X?" questions; time-to-answer
  with `grep -ri` over `parley-deck/ideas/` (18 MB, ~89 idea dirs — PRIMARY, measured)
  vs with the index. If grep answers all 10 in acceptable time, C2 is convenience, not a
  capability — say so and descope. Cross-deck recall (41 decks) is *not* solved by this;
  that needs a machine-level registry (`cache/projects.json` already maps project path →
  UUID, PRIMARY) and is a separate, later idea.
- Failure mode: index drifts from reality if not regenerated. Mitigation: regenerated on
  every `consensus finalize`/`reopen` (hooks already exist there), and it is disposable —
  delete and rebuild any time; never a source of truth.

**C3 — LLM-written semantic abstracts (OpenViking's `.abstract.md`, ~100 tokens).**
- Mechanism, if M1 shows structural indexes insufficient: the consensus-drafting agent
  emits a ≤100-word `## Abstract` block into `consensus.md` frontmatter/body as part of
  the existing draft step — marginal cost ~150 output tokens once per idea, not per
  artifact, because that agent has already read everything. Record `source-sha256` of the
  summarized artifact in the abstract block so staleness is detectable by the validator.
- Cost: real but bounded (~150 tokens/idea). Per-artifact or per-directory abstracts (the
  full OpenViking shape) are rejected: N×M generation cost, and fix-up cycles edit
  artifacts after the abstract exists.
- Measurement: same M1/M2 harness, comparing structural-index-only vs +abstract. Adopt
  only if abstracts buy a measurable margin over C1 alone.
- Failure mode — the dangerous one: a **stale or wrong abstract acted on as truth**.
  Concretely: a reviewer signs off having read the abstract, the artifact says the
  opposite, and the false consensus propagates into FINAL.md. This is worse than no
  summary. Mitigations: `derived: true` + `source-sha256` + a protocol rule (next core
  version, attended publish — not an edit to COOPERATION.md) that abstracts are
  non-citable as evidence. Without the staleness check, do not ship C3 at all.

## Existing alternatives

Enumeration of what each hand-built piece would duplicate. Sources consulted (all
PRIMARY, local): `parley --help`, `internal/runner/round_index.go`,
`internal/runner/phase58.go`, `internal/retro/retro.go`, `internal/app/learn.go`,
`internal/repomap/repomap.go`, `internal/app/retro.go`, `parley-deck/playbooks/`,
`cache/projects.json`, `git log`. Verdict per element: constraint-forced or inherited.

1. **Per-round L0 index** — closest shipping: `BuildRoundIndex`
   (`internal/runner/round_index.go:96`) → `round-NN/_index.md`. Covers: token estimates,
   section map, excerpts, per-agent status. Does not cover: semantic summary, and it is
   deliberately *excluded* from later prompts (`phase58.go:298`). **Inherited** — already
   built; only consumption (C1) is constraint-forced new code.
2. **Index-first prompt assembly** — closest shipping: `gatherPriorRounds`
   (`internal/runner/phase58.go:295-301`) inlines full prior-round artifacts; no
   substitution mode exists. **Constraint-forced**, but small; no dependency.
3. **Cross-idea mining** — closest shipping: `parley retro scan` (`internal/retro/retro.go`)
   — deterministic failure-density scoring over all ideas' structured artifacts, compact
   JSON. Covers: "which ideas were expensive/painful". Does not cover: content recall
   ("what did we decide about X"). **Partially inherited**; C2 fills the content gap.
4. **Idea distillation** — closest shipping: `parley learn <slug>`
   (`internal/app/learn.go`) → `parley-deck/playbooks/<topic>.md` (observed 1,444 bytes for
   `protocol-generation-bias`). Covers: advisory playbook for COMPLETED ideas, opt-in.
   Does not cover: open ideas, automatic generation, per-round granularity; the generated
   file is mostly template needing human generalization. **Inherited but weak**; not a
   substitute for C1/C2.
5. **Repo structure map** — `parley context repo-map` (`internal/repomap/repomap.go`):
   Go AST symbols of *code*. Irrelevant to deck artifacts. Null result.
6. **Stable addressing** — closest shipping: git (blob SHAs, `git log --follow`) +
   slug-keyed paths + `parley consensus status <IDEA>`. OpenViking's `md5(account_id:uri)`
   (notes §4, PRIMARY) is equally rename-fragile, so the URI scheme buys nothing over
   paths here. **Inherited coverage is sufficient; do not build.**
7. **Semantic/vector search** — nothing ships; platform alternative is `grep -ri` over
   18 MB (fast, measured in M3). A vector DB violates the zero-dependency constraint.
   **Constraint-forced null result: the hand-built route (index + grep) is correct.**
8. **Machine-wide cross-deck recall** — nothing ships; `cache/projects.json` is only a
   path→UUID map. **Constraint-forced** if wanted, but out of this idea's scope.

## Concerns / open questions

- Does M1's ≥50% threshold hold on the ideas that matter? The replay harness must be
  built before any adoption vote; it is itself the deliverable of round 2 if this proceeds.
- Will agents actually *use* on-demand reads, or treat the index as sufficient and get
  lazy? M2 measures exactly this; if answerability is high, laziness is not costly.
- `_index.md` excerpts are first-paragraph-180-char truncations — fine for structure,
  possibly misleading for contentious content. Check whether excerpt quality caps M2.
- Who regenerates `INDEX.md` for the 41 other decks with older CLI versions? Additive
  regeneration on next CLI touch; never a migration.
- Notes caveat (§3): the ~100/~2000-token L0/L1 budgets come from the URI concepts doc,
  not the blog; I used the concepts doc. Treat the exact numbers as illustrative, not
  prescriptive — our budgets should come from M1, not from copying OpenViking's.

## Risks

- **Summary-as-truth contamination** (C3): a stale abstract driving a signoff is a
  protocol-integrity failure, not a UX bug. Gated behind `source-sha256` + non-citable
  rule; without both, reject C3 regardless of M1.
- **Adoption failure**: if the index is ever wrong once, agents learn to bypass it and
  inline everything again — savings evaporate. Deterministic generation (no LLM in C1/C2)
  is the mitigation; this is why C1 precedes C3.
- **Measurement theater**: replaying byte counts is easy to game (e.g. by moving cost to
  agent tool-reads at runtime). M1 must count *total* tokens: inlined + on-demand reads.
- **Scope creep** toward a server, embeddings, or a URI resolver: all vetoed by the
  zero-dependency and `cat`/`grep`/`git` constraints; any round-2 proposal containing
  one should be blocked on that basis.
- **Divergence risk is low**: C1/C2 are additive derived artifacts; worst case is deleting
  them, which restores today's behavior exactly.
