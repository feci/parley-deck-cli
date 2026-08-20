---
idea: sync-skill-protocol-fallback
author: claude-1
created: 2026-06-24
status: final
participants: [claude-1, codex-1, hermes-1, antigravity-1]
---

## Problem / idea

The npm package **parley-deck-skill** (currently 1.4.0) bundles a portable, vendor-neutral
fallback copy of the protocol at `references/COOPERATION.md`. After the loop-engineering
release (parley-deck-cli 1.31.0), that bundled snapshot is **stale across several protocol
generations**: it has 14 sections / 759 lines, while the CLI's canonical genericized template
(`internal/protocol/defaults/COOPERATION.md`) now has 16 sections / 1037 lines.

The bundled fallback is missing / behind on at least:
- **§13 Retrospective optimization** (shipped earlier, RHO).
- **§14 Automated outer loop — the human brake** (loop-engineering Tier 4).
- The **§0 "Deck bootstrap (one-time)"** paragraph (roster + model + reasoning confirmation).
- The expanded **§4 Phase rules**: Phase 6 refutation-default + `## Refutation attempts` +
  model-diversity, Phase 8 `strict_gate` driver enforcement, the §4 loop-budget invariant
  (LE-5) and close-decision integrity (LE-7/LE-11).
- **§12.11 candidate-remediation** (`status: candidate`).

It has ALSO diverged in placeholder STYLE: the skill ref uses `**Workspace:** \`parley-deck\``,
`**Transport:** \`<transport-choice>\``, `**Created:** \`<YYYY-MM-DD>\``, whereas the CLI
embedded default uses `<workspace-name>` / `github-pr` / `<date> — created by parley init`.

Per the skill's own model the bundled reference is a portability fallback (the live deck file
is canonical, a stale fallback only warns) — so this is hygiene, not a breakage. But the §14
human brake is a binding safety rule a portably-installed agent should carry.

## Decision to make

How to faithfully bring `references/COOPERATION.md` current. Candidate strategies:

- **(A) Single source of truth** — make the skill reference a verbatim copy of the CLI
  embedded default (`internal/protocol/defaults/COOPERATION.md`). Simplest to keep in sync
  forever (could even be drift-guarded), at the cost of adopting the CLI default's header
  placeholder style (`<workspace-name>`, `github-pr`, `parley init` mention).
- **(B) Genericized merge** — bring in every new/amended section but re-genericize them to the
  skill's existing vendor-neutral placeholder style, preserving the skill ref as a CLI-neutral
  artifact (no `parley init` mentions in the header zone).

Pick one (or a hybrid), and specify the exact genericization rules, the complete set of
sections/edits to import, and how the skill's vendor-neutral framing (it is "instructions for
all agents", not tied to the `parley` CLI) is preserved.

## Constraints

- The bundled reference must remain **vendor-neutral** (it ships to agents that may not have
  the `parley` CLI). Decide whether `parley init` / `parley loop tick` literals are acceptable.
- Content must be **faithful** to the canonical protocol — no invented or dropped rules.
- Implementation target is the **parley-deck-skill** repo (`references/COOPERATION.md`), then a
  version bump + `npm publish` + GitHub tag per its `RELEASING.md`.

## Non-goals

- Changing the protocol itself (this only syncs the bundled snapshot to what already shipped).
- A protocol-sync automation script (could be a follow-up; out of scope here).
