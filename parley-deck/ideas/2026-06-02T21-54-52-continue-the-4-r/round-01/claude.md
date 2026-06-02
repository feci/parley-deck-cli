---
agent: claude
idea: 2026-06-02T21-54-52-continue-the-4-r
round: 1
date: 2026-06-02
---

## Summary

Three of the four are code (1, 2, 4); item 3 is infra-blocked and must not be faked. Build order: 4 → 2 → 1, with 3 handled as a documented blocker + optional release creation. Item 4 (agy) is the cheapest and unblocks honest multi-agent runs again; item 2 (Phase 8 loop) is medium and self-contained; item 1 (parallel DAG) is the largest because it changes the durable cursor shape. I verified agy CAN write its artifact in `--print` mode, so item 4 is a launch-robustness fix, not a dead end.

## Proposed approach

- **Item 4 (agy) — do first.** Root cause: `--add-dir {root}` on a large repo + a big idea prompt makes agy explore ("I will list the directory contents…") instead of writing. Two complementary, low-risk fixes, both in OUR code: (a) a generic "write-first" preamble in the round/review prompt — "Immediately use your file-writing tool to create exactly <outputPath>; do not explore the workspace first" — safe for every agent; (b) a **stdout-capture fallback** in `runAgent`: if after the process exits the artifact file is absent but stdout starts with `---` (valid frontmatter), persist stdout to the artifact. The agent still authored the content, so the agents-write-markdown rule holds; this also hardens other print-only CLIs. Keep it behind a clear condition (file missing AND stdout looks like the artifact). Test with a fake agent that prints to stdout without writing.
- **Item 2 (Phase 8 loop) — second.** Define the machine-readable contract: `review/consensus.md` frontmatter `outstanding_agreed_fixes: <int>` (already used once) is the authority; add `runner`/`pipeline` helper `ReviewConsensusOutstanding(path) (int, bool)`. Auto loop for an implementation block: RunImplementation (if no IMPLEMENTATION.md) → for cycle in 1..maxCycles: RunReviewRound(round=cycle) → draft review/consensus.md → requestConsensusSignoffs → read outstanding; if 0 → mark IMPLEMENTATION.md complete + advance; if >0 → RunFixup (implementer re-invoked to apply agreed fixes + bump IMPLEMENTATION.md fix-up cycle) → next cycle; stop on BLOCK signoff or maxCycles. Reuse RunReviewRound for round-02+. Test the loop control with fakes + a stubbed outstanding-count.
- **Item 1 (parallel DAG) — third.** Cursor: add additive `ready_blocks[]` + `active_blocks[]` (bump cursor schema_version, default from `current_block` when absent so old runs load). `advanceDAG` already computes ReadyBlocks; generalize the driver to return a SET of ready blocks and have `auto` launch each (bounded concurrency, e.g. the existing parallel cap). Keep single-active as `--max-active 1`. Per-edge boundary gates unchanged; production non-bypassable. Deterministic test: a diamond DAG (spec→{api,ui}→integrate) asserts api+ui both become active after spec, integrate only after both.
- **Item 3 (winget) — blocked.** Confirmed: only tags v1.6–1.9 exist, no GitHub releases, no `.exe` assets (latest release v1.5.4). Safe actionable step: create GitHub releases from the tags (triggers the portable-build CI), then later fill real InstallerSha256 + open the microsoft/winget-pkgs PR. Do NOT invent hashes or open a placeholder PR. Record this as a `to-user` inbox escalation so the decision (create releases now?) is explicit.

## Concerns / open questions

- Cursor schema migration (item 1) is the top backward-compat risk: existing `pipeline-run.json` (current_block only) must keep loading. Mitigation: additive fields + `ready_blocks` defaults to `[current_block]` when empty; round-trip test an old cursor.
- Parallel DAG + per-block engine state: multiple active blocks means multiple block idea-workspaces in flight; ensure no shared-file contention (each block has its own `ideas/<slug>__<block>/`). Bounded concurrency to avoid launching many hosted agents at once.
- Phase 8 loop must be genuinely bounded (max-cycles) and must stop on any BLOCK signoff; never loop forever.
- agy stdout-fallback must not mask real failures: only trigger when the file is absent AND stdout is a plausible artifact; record in the run events that a fallback was used.

## Risks

- **Faking item 3** is the cardinal risk — surface the blocker, don't open a PR with placeholder shas. Mitigation: explicit inbox escalation + no upstream PR until assets exist.
- **Cursor break** (item 1) → existing pipelines fail to resume. Mitigation: additive + schema bump + defaulting + old-cursor round-trip test.
- **Unbounded Phase 8 loop** → runaway agent spend. Mitigation: max-cycles cap + BLOCK stop + explicit stop statuses.
- **agy preamble regressing other agents** → keep it a neutral "write the file first" instruction that helps all and harms none; the stdout fallback is the real agy safety net.
- Smallest shippable if time is short: item 4 (agy fix) + item 3 release-creation, deferring 1 and 2 — but the goal is 1, 2, 4 implemented and 3 unblocked-or-documented.
