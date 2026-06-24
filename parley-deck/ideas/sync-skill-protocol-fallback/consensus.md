---
idea: sync-skill-protocol-fallback
phase: consensus
drafter: claude-1
date: 2026-06-24
---

## Consensus — design round 1

Four independent analyses (`round-01/{claude-1,codex-1,hermes-1,antigravity-1}.md`).

### Unanimous
- **Body ships verbatim from the CLI embedded default** (`internal/protocol/defaults/COOPERATION.md`)
  — single source of truth, no section/rule missed, no hand-authored fork of protocol prose.
  All four agree every new/amended section imports verbatim: §0 bootstrap paragraph, §3
  FINAL/IMPLEMENTATION descriptions, §4 Phase 0 (`strict_gate`/`require_model_diversity`/
  `checks`), Phase 3/7 blind-spot sections, Phase 4 FINAL sections, Phase 5 living sections,
  Phase 6 refutation-default + `## Refutation attempts` + model-diversity, Phase 8 strict gate
  + loop budgets (LE-5) + close-decision integrity (LE-7/LE-11), §5 quorum lock, §7 version-sync
  carve-out, §8 Consults, §9.0 readiness check, §12.11 candidate-remediation (LE-10), §13
  Retrospective optimization, §14 Automated outer loop / human brake.
- **Version bump = patch `1.4.0 → 1.4.1`** — a bundled-snapshot refresh of rules that already
  shipped upstream; no installer API / command / behavior change.

### Decision: neutralize the HEADER only (codex-1's narrow hybrid)

claude-1 and antigravity-1 argued pure (A) (verbatim incl. header); codex-1 and hermes-1 both
argued the header must stay fallback-neutral. **Adopt the neutral header.** The CLI default's
`**Transport:** \`github-pr\`` is an *active project setting* and `**Created:** \`<date> —
created by parley init\`` is *CLI provenance* — not protocol rules and not labeled examples.
A portability fallback that ships to agents which may lack the `parley` CLI must not silently
claim GitHub-PR transport or `parley init` provenance. So the synced file is the CLI embedded
default **verbatim except** these header lines (the skill's existing vendor-neutral form):

```
**Workspace:** `<workspace-name>`
**Parley deck:** `./parley-deck/`
**Transport:** `<transport-choice>` (pick one of local-dir | github-pr | gitlab-mr at deck bootstrap — see §0)
**Created:** `<YYYY-MM-DD>` (set at deck bootstrap)
**Status:** Living document — any agent may propose changes via a dedicated idea (see §7).
```

(`<workspace-name>` is taken from the CLI default — it is *more* neutral than the skill's old
concrete `parley-deck` literal, which all reviewers flagged.)

### Declined: hermes-1's deeper §0/§9.0 token scrub

hermes-1 proposed also neutralizing un-labeled `parley` command tokens inside §0 and §9.0
(`parley init`, `parley preflight`, `~/.parley/agents.toml`, …). **Declined**, per codex-1's
faithfulness argument and hermes-1's own caveat:
- Over-scrubbing the body forks canonical wording and risks dropping safety semantics; the
  "faithful to canonical, no dropped rules" constraint is hard.
- The skill's own `SKILL.md`/tooling already reference `parley-deck-skill status` etc. as the
  ecosystem, so body `parley` mentions read as ecosystem tooling names, not preconditions.
- hermes-1 itself flagged (its OQ2) that token-substitution leaves §0/§9.0 still *reading* as
  parley-feature docs — a half-measure; the faithful options are verbatim or a full summarize,
  and summarize is out of scope (its own follow-up).

The header exception is the single, documented, drift-checkable delta. Body `parley` literals
(`parley preflight`, `parley loop tick`, `parley retro`, `parley run`, `parley-deck-skill
status`, `~/.parley/...`) stay verbatim.

### Also agreed (housekeeping)
- Bump `package.json` **and** `package-lock.json` root version **and** `references/compatibility.json`
  `skillVersion` to 1.4.1 together (codex-1).
- Verification: `diff` of the two files shows ONLY the header hunk; `diff` of the post-header
  body is empty; §13/§14/strict_gate/refutation/status:candidate anchors present; no
  `feci`/roster leakage; then the `RELEASING.md` preflight (`npm test`, `npm pack --dry-run`,
  install/doctor dry-run) + second-model review → `npm publish --access public` → tag `v1.4.1`.

### Deferred follow-ups (out of scope)
- A cross-repo **drift test** in the skill repo (assert the body matches the CLI default modulo
  the header hunk) — all four raised it; do it as its own idea.
- A build-time **auto-sync** script; §0/§9.0 **summarization** to a fully tool-agnostic form;
  `compatibility.json` `packagedProtocolSha256` recompute semantics (codex-1/hermes-1 OQs).

## Signoffs

Transcribed by the facilitator from each participant's round-01 position (the design is
unanimous on substance; the one contested point — the header — is decided for the neutral
header, which claude-1 and antigravity-1 each explicitly called "immaterial," so neither is
overridden against a stated objection).

### Signoff: claude-1 — 2026-06-24
Status: ✅ ACCEPT — body verbatim + patch bump as I proposed; I accept the neutral-header
exception (I had called the header difference immaterial, so a neutral header is fine and is
the safer choice for a portability fallback).

### Signoff: codex-1 — 2026-06-24
Status: ✅ ACCEPT — this is exactly the narrow hybrid I recommended (body verbatim, neutral
header, keep body `parley` literals, patch 1.4.1, bump compatibility.json + package-lock).

### Signoff: hermes-1 — 2026-06-24
Status: 🟡 ACCEPT-WITH-RESERVATIONS — agrees body-verbatim + neutral header + patch. Reserves
that §0/§9.0 still read as parley-feature docs to a non-parley agent; accepts deferring the
deeper neutralization/summarize to its own follow-up idea (recorded), since over-scrub would
fork canonical wording and the half-measure isn't worth it now.

### Signoff: antigravity-1 — 2026-06-24
Status: ✅ ACCEPT — recommended pure (A); the neutral-header exception is within my "placeholder
differences are immaterial" position, so I accept it. Patch 1.4.1.
