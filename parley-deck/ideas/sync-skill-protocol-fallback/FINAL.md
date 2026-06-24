---
idea: sync-skill-protocol-fallback
status: final
drafter: claude-1
date: 2026-06-24
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Decision

Bring `parley-deck-skill/references/COOPERATION.md` current with the canonical protocol by
copying the CLI embedded default (`internal/protocol/defaults/COOPERATION.md`) **verbatim**,
**except** the two header lines that are active-project state / CLI provenance, which stay in
the skill's vendor-neutral form. Then bump the skill to **1.4.1** and publish to npm + tag.

## Implementation

1. Copy `internal/protocol/defaults/COOPERATION.md` → `parley-deck-skill/references/COOPERATION.md`.
2. In the copied file, replace exactly the `**Transport:**` and `**Created:**` header lines with:
   - `**Transport:** \`<transport-choice>\` (pick one of local-dir | github-pr | gitlab-mr at deck bootstrap — see §0)`
   - `**Created:** \`<YYYY-MM-DD>\` (set at deck bootstrap)`
   Keep `**Workspace:** \`<workspace-name>\`` and every other line as the CLI default has it.
   The body (everything after the header) is byte-identical to the CLI default — including the
   `parley`/`~/.parley` literals in §0, §9.0, §13.4, §14 (ecosystem tooling names, not scrubbed).
3. Bump version to `1.4.1` in `package.json`, the root entry of `package-lock.json`, and
   `references/compatibility.json` `skillVersion` (and any sibling version field that tracks it).

## Verification (gates before publish)

- `diff -u internal/protocol/defaults/COOPERATION.md references/COOPERATION.md` → shows ONLY the
  2-line header hunk (Transport + Created).
- `diff <(tail -n +7 <CLI default>) <(tail -n +7 <skill ref>)` → empty (body identity; lines 1–6
  are the header, line 7 = `**Status:**` onward identical).
- Anchors present exactly once: `^## 13\.`, `^## 14\.`; safety tokens present: `Refutation
  attempts`, `strict_gate`, `Loop budgets`/`loop budget`, `status: candidate`.
- No leakage: `grep -in "feci\|claude-1\|codex-1\|hermes-1\|antigravity-1" references/COOPERATION.md`
  → none.
- `RELEASING.md` preflight: `npm test`, `npm pack --dry-run`,
  `node bin/parley-deck-skill.js install --target all --dry-run`,
  `node bin/parley-deck-skill.js doctor --target all --json`.
- Second-model review of the final diff (the Phase 6 review of this idea).
- Then `npm publish --access public`; `git tag v1.4.1 && git push origin main v1.4.1`.

## Out of scope (follow-ups)

- A skill-repo drift test (assert body == CLI default modulo the header hunk).
- A build-time auto-sync script; §0/§9.0 summarization to a fully tool-agnostic form;
  `compatibility.json` `packagedProtocolSha256` recompute semantics.
