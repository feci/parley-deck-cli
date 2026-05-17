---
agent: claude
idea: roadmap-implementation-plan
round: 2
date: 2026-05-17
---

## Summary

**Accept with reservations.** The facilitator's revised slice 1
(`round-index-artifact-pruning`) is a good convergence: it preserves the
deterministic round-summary index I proposed and folds in Gemini's
hidden-reasoning sanitizer as a sibling helper, both written into the same
runner-side context-helpers module. Both halves are deterministic, local,
write-only-to-new-files, and naturally consumed together by every later
slice (context-pack-wiring, hooks, memory-insights). Codex's downstream
ordering (`repo-map-mvp` → `context-pack-wiring` → `error-classifier-redaction`
→ `hooks-mvp` → `local-memory-insights` → `direct-sdk-backend-pilot`) is
sound and matches the architectural fence I argued for in round 1: nothing
provider-dependent until a direct SDK backend is on the table. My
reservations are about keeping slice 1 reviewable in one PR and not
letting the sanitizer drift into a security-redaction role.

## Position on proposed roadmap

Accept the seven-slice sequence as written.

- Slice 1 (`round-index-artifact-pruning`): accept with the changes listed
  below.
- Slices 2–7: accept as ordering. Slice 4 (`error-classifier-redaction`)
  must own all secret-redaction semantics; slice 1's sanitizer must not
  pretend to cover that case.
- Slice 7 (`direct-sdk-backend-pilot`) is correctly fenced as conditional.
  It is the only path that unlocks real provider prompt caching, semantic
  cache, and in-process sub-agent delegation. The roadmap should not claim
  those wins from earlier slices.

## Required changes before consensus

These are the specific guardrails I need in slice 1 before I can sign off:

1. **Explicit, finite pattern list for the sanitizer.** Document the
   supported hidden-reasoning fences (`<think>`, `<thought>`, `<thinking>`
   and their closing tags) as a closed set in code and in a short doc
   comment. New patterns require a follow-up PR, not a regex tweak in
   review. Prevents slow scope creep.
2. **Sanitizer is context-only, never mutates source artifacts.** The
   function signature should make this obvious (`SanitizeForContext(s
   string) string`), and a test must assert that running it does not write
   to the round directory. This was implicit in codex's note; I want it
   pinned.
3. **Sanitizer is explicitly not redaction.** Add a doc-comment line
   stating that secret redaction lives in the future
   `error-classifier-redaction` slice. Failing to fence this now risks
   future readers assuming the sanitizer handles credentials.
4. **Index output is byte-deterministic for fixed inputs.** Either omit
   `generated_at` from the body, or pin it to the maximum mtime of the
   participating agent artifacts. Golden-file tests should pass on
   repeated runs without manual normalization.
5. **Decide canonical vs derived up front.** I argued canonical in round 1;
   codex's proposal reads canonical (runner writes it into the round
   directory). Confirm in the slice 1 PR description so reviewers know
   whether to commit `_index.md` files or gitignore them. My preference
   remains canonical — replay and consensus tooling benefit.
6. **Index write failure is a warning, not a round failure.** Codex
   already flagged this. The runner should record an `index_write_failed`
   result row and exit 0 if all participant artifacts were written
   successfully.
7. **Skipped/failed agents appear as rows in the index.** A round where
   one agent blocked or errored should still produce a complete
   `_index.md` with explicit status, not omit the agent silently.

If these seven are in the slice 1 PR description and tests, I accept
without further blocking.

## Risks

- **Bundled-slice review surface.** Combining the index writer and the
  sanitizer doubles the review area of what I originally framed as the
  smallest reviewable step. Acceptable because both halves are small,
  share a target package, and each has independent fixture-based tests —
  but the PR description must list both behaviors separately so a reviewer
  can sign off on each.
- **Sanitizer brittleness.** Regex over LLM output is inherently fragile.
  Mitigation: the explicit closed pattern list plus tests covering
  malformed/nested cases. Do not market the sanitizer as complete
  reasoning control; it is a best-effort trimmer for context injection.
- **Drift between slice 1 sanitizer and slice 4 redaction.** If slice 4
  later adds a second pass over the same artifacts for secret patterns,
  the two passes must compose cleanly. Easiest mitigation: have slice 4
  call slice 1's sanitizer as a first step, then apply redaction on top.
- **False progress claim on context efficiency.** The index plus sanitizer
  is a real but bounded win — fewer tokens to re-inject across rounds. It
  is not prompt caching and it is not trajectory compression. The roadmap
  language should keep those wins assigned to slice 3
  (`context-pack-wiring`) and the conditional slice 7, not slice 1.
- **Schema dependence.** The index writer assumes agents emit the
  expected `## Summary` / `## Proposed approach` / … headings. Degrade
  gracefully when absent (capture the agent frontmatter and a "no
  recognized sections" row) rather than failing the index write.
