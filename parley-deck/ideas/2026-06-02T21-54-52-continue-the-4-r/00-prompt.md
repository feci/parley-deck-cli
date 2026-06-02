---
idea: 2026-06-02T21-54-52-continue-the-4-r
author: user
created: 2026-06-02
participants: [codex, claude, hermes]
status: final
---

## Problem / idea

Continue the 4 remaining COOPERATION.md §12 pipeline follow-ups in parley-deck-cli (on top of the shipped 1.9.0 surface: internal/pipeline {manifest,run,gate,executor,effects,provider,watcher,dag,watch_eval}, internal/runner {round,phase58}, and `parley pipeline validate|start|status|run-block|continue|auto|watch|execute|record-effect|gate`). Additive; preserve linear back-compat, agents-write-markdown/driver-executes, non-bypassable production gates. Give independent Phase-1 analysis + the smallest correct implementation + risks per item, and the recommended build order.

ITEM 1 — Parallel multi-active DAG auto-drive.
Today `execution: dag` is validated and `advanceDAG` runs SINGLE-ACTIVE (one block at a time, dependency-ordered via ReadyBlocks). Goal: drive ALL currently-ready blocks. Proposed: extend the cursor (`pipeline-run.json`) with additive `ready_blocks[]` / `active_blocks[]` (schema bump + zero-value defaulting; keep `current_block` populated for linear back-compat); `parley pipeline auto` (and a new status view) launches/advances every ready block whose inbound edges are complete and gates resolved, concurrently. Boundary gates still per edge; production still non-bypassable. Question: cursor schema evolution without breaking existing runs; how `auto` tracks multiple in-flight blocks + their per-block rounds/consensus; deterministic ordering for tests; whether to bound concurrency.

ITEM 2 — Fully-unattended Phase 8 fix-up loop.
Today `auto` for an implementation block runs Phase 5 (RunImplementation) + review round-01 (RunReviewRound) then STOPS for human review-consensus. Goal: auto-drive Phase 7 (draft review/consensus.md, collect signoffs) + Phase 8 (implementer applies agreed fixes, re-review) until `outstanding_agreed_fixes: 0`, bounded by a max-cycles cap. Proposed: a machine-readable `agreed_fixes` contract in review/consensus.md (e.g. a fenced yaml block or frontmatter count) + a helper `ReviewAgreedFixes(path) int`; a loop: review round → consensus.Draft → requestConsensusSignoffs → if agreed_fixes>0 → RunFixup (implementer) → next review round; stop at zero or max cycles or a BLOCK signoff. Question: exact agreed_fixes machine format + who writes it (drafter vs structured-from-reviews); how a fix-up re-review reuses RunReviewRound (round-02+); testability with fake agents; safety cap; what counts as "applied".

ITEM 3 — WinGet upstream PR to microsoft/winget-pkgs.
RECON FINDING (state it in your analysis): there are NO GitHub *releases* for cli v1.6.0/1.7.0/1.8.0/1.9.0 — only git tags exist; the latest GitHub release is v1.5.4. The portable .exe assets are built by CI on release publication, so no real InstallerSha256 values exist yet. Therefore the upstream microsoft/winget-pkgs PR is BLOCKED on infra: (a) create GitHub releases for the tags (triggers the portable-build CI), (b) wait for .exe assets, (c) fill InstallerSha256 from the immutable assets, (d) validate + open the upstream PR. Question: is creating the GitHub releases (gh release create from existing tags) the right unblocking step to do now, leaving the upstream PR for once assets exist? Do NOT invent hashes or open a PR with placeholder shas.

ITEM 4 — Fix the agy headless launch regression.
RECON FINDING: agy 1.0.4 CAN write its artifact in `--print` mode — verified `agy --print-timeout 120s --dangerously-skip-permissions --add-dir <dir> --print "<runner-style prompt asking to create <path>>"` correctly wrote the file. The real `parley run` failure ("artifact was not created", agy narrated "I will list the directory contents...") most likely occurs because `--add-dir {root}` points at a LARGE repo root and agy explores it instead of writing, especially with a big idea prompt. Proposed fix (in OUR launch spec, agy is third-party): make the agy launch more directive/robust — e.g. an agy-specific prompt preamble instructing it to immediately use its file-write tool for the exact output path before any exploration, and/or a stdout-capture fallback that writes the agent's printed markdown to the artifact when the agent did not create the file itself (validated to start with `---`). Question: where to put agy-specific behavior (discover.go spec note vs runner prompt branch) without special-casing one vendor too deeply; is a generic "write-first" preamble safe for all agents; is the stdout fallback acceptable under the agents-write-markdown rule (the agent still authored the content); how to test without launching real agy.

CONSTRAINTS: additive + backward-compatible; preserve agents-write-markdown/driver-executes + non-bypassable production gates; vendor-neutral (no hard agy/vercel coupling); Go 1.26; full test coverage with fake agents (no reliance on real hosted agents in unit tests); English-only protocol files; do not invent winget hashes or open an upstream PR without real release assets.

DELIVER in round-01: per-item assessment + minimal correct implementation + top risks; explicitly confirm item 3 is infra-blocked and what the safe actionable step is; recommended build order.

## Constraints

- Local-directory transport for this initial run.
- Human-in-the-loop mode by default unless the run was started with auto mode.

## Non-goals

- Do not make unrelated repository changes.
