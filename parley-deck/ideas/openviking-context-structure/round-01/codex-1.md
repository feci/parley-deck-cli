---
agent: codex-1
idea: openviking-context-structure
round: 1
date: 2026-08-31
---

## Summary

Recommendation: do not adopt OpenViking, its sidecar hierarchy, `viking://`, or a new database; test one narrow idea—progressive loading—by wiring Parley's existing round digest and round index into an opt-in context pack. **PRIMARY (copied OpenViking notes; not independently re-fetched):** OpenViking's relevant mechanism is an optional ~100-token L0 abstract, ~2,000-token L1 overview, and full L2 source, while Parley already has analogous deterministic extracts in `internal/driver/digest.go:BuildRoundDigest` and `internal/runner/round_index.go:BuildRoundIndex`. **PRIMARY (local measurement):** the five active participants' `protocol-generation-bias` round-01 files total 88,703 bytes and round-02 totals 120,267 bytes; because `internal/runner/runner.go:gatherPriorRounds` embeds every earlier `.md` file into every participant prompt, rounds 2 and 3 repeat 1,488,365 raw bytes before wrapper text, about 372,000 tokens only under the repository's explicit bytes/4 heuristic. Cross-deck recall and logical artifact identity are plausible gaps, but their usage and benefit have not been measured; defer them.

## Proposed approach

Adopt an experiment, not a storage migration.

1. Add a read-only command:

   `parley context round-pack [--dir DIR] [--through N] [--tier l0|l1|l2] [--format markdown|json] IDEA`

   - L0 reuses `BuildRoundDigest`: one capped position sentence per participant.
   - L1 reuses `BuildRoundIndex`: status, approximate size, H2 section names, and the first text line per section.
   - L2 is the unchanged participant artifact.
   - The command writes no canonical artifact. The runner may consume the same library behind an explicitly experimental flag; the default remains full L2.

2. Bind every derived entry to its source. Extend runner-owned `round-NN/_index.md` with the artifact path, SHA-256 of the raw source bytes, and sanitizer version. Recompute before packing. A missing index, digest mismatch, parse failure, or truncation ambiguity falls back to L2 and reports the reason. Existing ideas need no migration because missing derived data already selects L2.

3. Preserve gate-bearing text. The tiered pack must include full L2 for an artifact containing protocol markers such as `❌ BLOCK`, `Counter-proposal`, `DISPUTED`, or `ALT-`; summaries remain navigation aids, never evidence for silently discarding a position. This does not eliminate the risk of an important unmarked qualification, so tiered mode must stay outside live canonical rounds until measured and ratified.

4. Use deterministic budgets, not another model call. **PRIMARY (local code):** `internal/runner/runner.go:RunRoundOne` already generates `_index.md` immediately after participant processes finish. Target L0 at most 160 heuristic tokens per participant and L1 at most 512, including locators; mark truncation explicitly. Generation spends zero model tokens, performs one local read/hash pass, and should meet p95 under 100 ms for a 1 MB round on the test machine. Full-source expansions remain ordinary model input and must be counted.

5. Benchmark before adoption on the five largest closed ideas with at least two design rounds.

   - Reconstruct full and tiered prompts from the same commit and frozen roster/model configuration.
   - Record sanitized prompt bytes, bytes/4 estimate, actual provider input tokens when exposed, wall time, L2 fallbacks, and output tokens. Do not label the heuristic as tokenizer output. **PRIMARY (local code):** headless runners do not yet emit `agent.usage` (`internal/driver/loop.go:loopCostUSD` says it is zero in practice), so prompt-byte telemetry is a required addition, not something `parley sessions` already supplies.
   - Build a blind gold set from prior blocks, counter-proposals, verdict conflicts, and adopted/rejected `ALT-` entries. Score whether the next-round output identifies and responds to each item, plus unsupported factual assertions.
   - Pass only if median real input tokens (or separately reported heuristic bytes) fall by at least 50%, every gate-bearing item is retained, material-item coverage loses no more than five percentage points, and no tiered run falsely reports convergence. Any lost blocker or adopted alternative is an automatic failure.

6. If the experiment passes, open a protocol-change idea before allowing an index to substitute for prior artifacts. **PRIMARY (local code/protocol):** the current runner tells participants to read every prior-round artifact, so changing that semantic through a CLI default alone would be dishonest. If it fails, retain `_index.md` for navigation/TUI use and keep full prompts.

Do not adopt the other OpenViking pieces now:

- No `.abstract.md`/`.overview.md` per directory: Parley's `FINAL.md`, `consensus.md`, `_index.md`, and round digest already occupy those summary roles; another generated layer adds a freshness owner without proven retrieval benefit.
- No OpenViking service, vector store, or network dependency: it violates the preferred deployment shape and is unnecessary for the measured prompt problem.
- No MD5-backed logical ID registry. **PRIMARY (copied source note, with direct inference):** OpenViking computes a file ID from `account_id:uri`; changing the URI changes the hash input, so the cited mechanism does not by itself establish identity across rename. Git already gives immutable versioned locators for the audit use case.
- No cross-deck semantic index yet. First collect real recall queries and compare current lexical discovery with manual lookup; a database built before a query workload exists is speculative infrastructure.

## Existing alternatives

The proposal would otherwise build the following components by hand. Coverage statements are **PRIMARY** from the named local code or command surface.

| Component | Load-bearing status | Closest thing that already ships | What it covers / does not cover |
| --- | --- | --- | --- |
| `RoundDigestExtractor` | Merely inherited; reuse it | `internal/driver/digest.go:BuildRoundDigest` and the `round.digest` run event | Produces a deterministic, LLM-free, 120-rune position line per participant plus keyword flags. It is TUI/run-event data, not a prior-round prompt source or durable cross-idea index. |
| `RoundOutlineExtractor` | Merely inherited; reuse it | `internal/runner/round_index.go:BuildRoundIndex`; runner-owned `round-NN/_index.md` | Already records status, bytes/4 token estimates, H2 headings, and first lines capped at 180 characters. `gatherPriorRounds` explicitly skips `_index.md`; it has no source hash and no total per-participant budget. |
| `TieredRoundPackSelector` | Constraint-forced by the measured prompt duplication | `internal/runner/runner.go:gatherPriorRounds`; `parley context repo-map` | The runner already builds a full L2 pack. `repo-map` (`internal/app/context.go`, `internal/repomap/repomap.go`) lists paths, sizes, kinds, and Go symbols with `--max-files`; it does not summarize Markdown, select context tiers, search content, or enforce a token budget. The tier selector is the actual new code. |
| `SourceDigestVerifier` | Constraint-forced by the stale-summary failure mode | Go `crypto/sha256`; `git hash-object`; `git rev-parse <commit>:<path>`; analogous hashing in `internal/runmanifest/manifest.go:RosterRevisionOf` | The primitives ship. What is missing is a raw-source hash plus sanitizer-version field in `_index.md` and a fail-to-L2 freshness check. A hash proves source equality, not summary adequacy. |
| `MaterialMarkerExpander` | Constraint-forced while summaries may hide blockers | `internal/driver/digest.go:stanceFlags`; signoff parsing in `internal/consensus/consensus.go` | Existing code detects hints or structured signoffs for display/gating. It does not force full artifact inclusion in a context pack. The new expander must remain conservative and cannot detect every material nuance. |
| `ContextPackTelemetry` | Constraint-forced by the evaluation brief | `parley sessions list|inspect`, `parley-deck/runs/<run-id>/events.jsonl`, `_index.md`'s bytes/4 estimate, ACP's `agent.acp.usage` event | Sessions expose run/workspace/idea/participant/manifest state, not prompt composition or general headless token usage. Add pack bytes, tier, source hashes, fallbacks, and available provider usage to run events; do not create another protocol artifact. |
| `CrossDeckIdeaIndex` | Rejected; not load-bearing | `parley retro scan|select|diagnose`, `parley learn`, `parley sessions`, `git grep` | `retro` scans one deck's structured failure/churn signals and selects hard cases; it does not query topics. `learn <known-closed-slug>` creates one advisory playbook skeleton and is not discovery. `sessions` can reveal run-indexed workspace roots and idea slugs, but not unrun decks or artifact content. `git grep -n -i -- parley-deck/ideas` provides zero-dependency lexical search inside one repository. None supplies semantic recall across 41 decks; the null result is from `internal/app/retro.go`, `internal/retro/retro.go`, `internal/app/learn.go`, `internal/app/app.go:runSessions`, and Git's shipped commands. The hand-built index is not justified yet. |
| `StableArtifactIDRegistry` | Rejected; not load-bearing | Git object/commit addressing: `git rev-parse <commit>:<path>`, `git show <commit>:<path>`, and `git log --follow -- <path>` | A commit-plus-path retrieves an immutable historical artifact even after a later rename and needs no migration. Git does not provide one logical ID spanning rewrites, copies, repositories, or rewritten history; Parley's protocol already declares idea slugs stable and closed artifacts immutable. No measured workflow currently needs more. |
| `ConsensusSummaryLookup` | Merely inherited; no new component | `parley consensus status [--review] IDEA` and canonical `consensus.md`/`FINAL.md` | The command validates known-idea signoff triage, participants, missing signers, and malformed blocks. The files summarize an already-converged idea. Neither discovers an unknown idea nor reduces context during deliberation. |

This inventory says the hand-built route is correct only for the selector, freshness/fallback logic, and telemetry. The summary generators, structural maps, session index, consensus state, playbook distillation, retrospective scan, and immutable version addressing already ship.

## Concerns / open questions

- Does a summarized pack satisfy the protocol's requirement to read prior artifacts, or does it require a new core version? My position is the latter; benchmark mode must not silently weaken the live rule.
- `_index.md` extracts the first non-empty line under each H2. A participant can put the decisive qualification later, so freshness hashing cannot make the extract semantically complete.
- A model may open every L2 source on demand, moving rather than reducing token cost. Actual provider usage, when available, is the deciding metric.
- Manual and historical rounds may lack `_index.md`; full fallback is required, and derived indexes must remain optional for `cat`, `grep`, and Git users.
- What real questions fail `git grep`, and how often? Without a query log or user sample, "semantic cross-deck recall" has no costed acceptance criterion.
- If a deck or repository moves, should session-index repair be solved separately from artifact identity? `sessions.json` stores absolute workspace roots, so conflating those problems would enlarge scope.

## Risks

- The worst failure is false convergence: a stale or incomplete summary hides dissent, an agent never opens L2, and consensus treats omission as agreement. Hash checks catch staleness, not wrong emphasis; the zero-loss blocker gate and L2 fallback are therefore mandatory.
- Tier selection and marker expansion can become a second protocol encoded only in Go. Any production default needs protocol text, tests that fail when the gate is removed, and plain-Markdown output.
- The bytes/4 estimate is language- and tokenizer-blind. Reporting it as real token savings would overstate evidence.
- Benchmarking on closed ideas can overfit to artifacts whose outcomes are already visible. Freeze inputs, blind the quality rubric, and report model/configuration details.
- A machine-wide recall index can cross project confidentiality boundaries and become stale when absolute workspace paths move. Do not build it as an incidental extension of the prompt experiment.
- Added code and metadata may cost more maintenance than saved tokens on ordinary one-round ideas. Keep full L2 as default and activate tiering only above a measured size threshold if the experiment passes.
