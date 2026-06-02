---
idea: 2026-06-02T21-54-52-continue-the-4-r
drafted-by: claude
date: 2026-06-02
---

## Agreed decisions

Round-01 converged (codex, claude, hermes). Additive; linear back-compat, agents-write-markdown/driver-executes, non-bypassable production gates preserved. **Build order: 4 → 2 → 1; item 3 is infra-blocked.**

1. **Item 4 — agy launch robustness (first).** Root cause: `--add-dir {root}` on a large repo + big prompt makes agy explore instead of writing. Fix (our code, generic): (a) a short generic **write-first preamble** in the round/review prompt ("Immediately use your file-writing tool to create exactly <outputPath> with the required frontmatter; do not explore the workspace first; report a blocker rather than narrating"); (b) a **stdout-capture fallback** in `runAgent`: if after exit the artifact is absent AND captured stdout starts with `---` (valid frontmatter), persist stdout as the agent-authored artifact and record `agent.stdout_fallback` in run events. Strict `---` validation prevents converting narration into artifacts. Tested with fake agents (writes file / prints valid md / prints narration).
2. **Item 2 — unattended Phase 8 fix-up loop (second).** Machine contract: `review/consensus.md` frontmatter `outstanding_agreed_fixes: <int>` (authority) + optional `blocked: true`. Helper `ReviewAgreedFixes(path) (count int, blocked bool, err error)`. Auto loop for an implementation block: RunImplementation (if no IMPLEMENTATION.md) → for cycle 1..maxCycles (default 3): RunReviewRound(round=cycle) → draft review/consensus.md → requestConsensusSignoffs → read outstanding; if blocked → STOP; if 0 → mark IMPLEMENTATION.md complete + advance; else RunFixup (implementer re-invoked) → next cycle. Auto **fails closed** if the machine field is absent. Reuse RunReviewRound for round-02+. Bounded by maxCycles. Tested with fakes + stubbed count.
3. **Item 1 — parallel multi-active DAG auto-drive (third).** Cursor: add additive `ready_blocks []string` + `active_blocks []string` (bump cursor schema_version; default from `current_block` when absent so existing runs load; keep `current_block` populated = first active/ready in deterministic order). Generalize `advanceDAG` to compute the full ready set (dependencies complete + inbound gates resolved), sorted by manifest block order then ID. `auto` launches all ready blocks up to a concurrency cap (flag `--max-active`, default 4); single-active = `--max-active 1`. Per-edge boundary gates + non-bypassable production unchanged. New status view lists ready/active/blocked/complete. Atomic cursor writes. Tested with a diamond DAG (spec→{api,ui}→integrate): api+ui both active after spec, integrate only after both.
4. **Item 3 — WinGet upstream PR (BLOCKED).** Confirmed: only git tags v1.6.0–1.9.0 exist; NO GitHub releases (latest release v1.5.4); no `.exe` assets → no real InstallerSha256. Do NOT invent hashes or open a placeholder microsoft/winget-pkgs PR. Safe actionable step: create GitHub releases from the tags (triggers the portable-build CI) — this needs release-owner (user) approval, so it is escalated to the user via inbox; the upstream PR waits for the resulting immutable assets + verified hashes.

## Agreed trade-offs

- stdout fallback slightly broadens "agents write markdown", but the agent authored the bytes; strict frontmatter validation keeps it narrow and it hardens all print-only CLIs.
- `outstanding_agreed_fixes` frontmatter constrains auto-driven review consensus docs (mandatory only in auto), bought a deterministic, testable Phase 8 loop.
- Concurrency cap (default 4) over unbounded: safer for agent quotas + deterministic tests; `--max-active` exposes the choice.
- Parallel DAG adds cursor state; mitigated by additive fields + schema bump + defaulting + atomic writes + old-cursor round-trip test.

## Open items deferred to implementation

- Exact write-first preamble wording (short, neutral, non-over-directive).
- `ReviewAgreedFixes` placement (pipeline vs runner) + the auto fail-closed behavior when the field is absent.
- `--max-active` default + env override; status-view format.
- Whether to create the GitHub releases now (user decision via inbox escalation).

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: codex — 2026-06-02
Status: ✅ ACCEPT
Notes: Consensus matches my position: implement 4 → 2 → 1 with bounded, additive changes; keep item 3 blocked until release-owner approval creates real assets and hashes.
Counter-proposal (required if ❌): N/A

### Signoff: hermes — 2026-06-02
Status: ✅ ACCEPT
Notes: Consensus correctly prioritizes 4→2→1, keeps item 3 blocked pending real assets, and preserves all core constraints (additive, agents-write-markdown, non-bypassable gates). The proposed fixes are minimal and testable.
Counter-proposal (required if ❌): N/A

### Signoff: claude — 2026-06-02
Status: ✅ ACCEPT
Notes: Drafter. Build order 4->2->1; item 3 honestly blocked (no releases/assets) -> escalate release-creation to user, no faked hashes/PR. Generic write-first preamble + strict stdout fallback; outstanding_agreed_fixes contract + bounded Phase 8 loop; additive DAG cursor + bounded parallel.
