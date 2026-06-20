---
idea: readme-marketing-intro
drafted-by: claude-1
date: 2026-06-20
status: awaiting-signoffs
---

## Agreed decisions

The four round-01 proposals converged with no substantive blockers. The intro is
built from all four lenses: claude-1's three-block structure, hermes-1's hook,
codex-1's accuracy guardrails + feature-claim map, antigravity-1's lineage format.

**1. Structure (≤ one screenful, before `## Install`):**
- **Hook** — name the category and the two anti-patterns it replaces:
  (a) ad-hoc N-agents-in-N-terminals (no audit trail / conflict discipline /
  consensus / resume), (b) one model role-playing a committee (solo reasoning in
  a costume; reintroduces single-model self-preference). Sharp-but-not-shade:
  framed as the value prop, not an attack on a workflow.
- **What you get** — a compact, citable bullet cluster (~7–8 bullets).
- **Inspired by — adopted & adapted** — antigravity-1's *Lineage / Our twist*
  format, compact, with the explicit non-endorsement disclaimer. The OpenRouter
  Fusion credit for the consensus lens stays in the intro, not a footer.

**2. Feature bullets (each maps to a real section/command — codex-1's map governs):**
- 8-phase idea lifecycle (§4): kickoff → independent analysis → cross-review →
  consensus → `FINAL.md` → `IMPLEMENTATION.md` → code review → fix-up; append-only.
- Non-solo by design (§1): stable agent IDs (§2), one file per agent per round (§6).
- Compare-not-merge consensus "Comparison & blind spots" lens (Fusion-derived).
- Transport-agnostic (§0, §11): `local-dir`, `github-pr`, `gitlab-mr`.
- Pre-idea readiness (§9.0, `parley preflight`): freshness + roster liveness ping.
- Advisory, quorum-gated retrospective (§13, `parley retro`).
- Agent supervision: watchdog, stall guard, failure classification, validated-
  artifact-beats-nonzero-exit — phrased as run supervision/recovery, not a guarantee.
- Living docs: self-contained `FINAL.md` + resumable `IMPLEMENTATION.md`.
- The `parley` CLI (TUI, auto-drive, pipelines) makes it usable, not just specified.

**3. Accuracy guardrails (codex-1, binding):**
- No ungated-autonomy claims. Auto-drive advances protocol *phases*; code
  implementation/fix-up is gated (`auto_implement`, review state, protocol gates);
  production mutations are non-bypassable (§12.4/12.8).
- "Comparison & blind spots" is a drafting *lens*, not an automatic truth engine.
- Advisory `consult` is non-canonical and not quorum evidence — mention lightly or omit.
- "Adopted and adapted" / "inspired by" wording; never imply Parley Deck invented
  compare-not-merge, resumable plans, retro optimization, dispositions, or preflight.
- "Vendor/model-agnostic roster" qualified as *subject to installed CLIs / auth /
  sandbox*, not "runs any model" unconditionally.

**4. Two READMEs:**
- `parley-deck-cli/README.md` — full intro as above, before `## Install`.
- `parley-deck-skill/README.md` — shorter variant: keep the protocol story +
  lineage block, **swap the CLI paragraph for the installer paragraph** (the skill
  *teaches* agents to participate; the CLI orchestrates). Soften/anchor away from
  transport-agnostic (the skill doesn't ship transports); keep 8-phase + consensus
  lens. Only claim shipped skill surfaces (install/doctor/status/sync-project/
  uninstall/paths/--version) — do not say the skill runs rounds/TUI/retro/pipelines.

**5. §-citations:** keep as honesty markers (named phases + a few `§` refs) but
stay skimmable — don't cite every bullet. They signal "spec'd, not vaporware."

## Trade-offs

- Sharp hook vs. friendly tone → chose moderately sharp (it's the strongest
  differentiator), but phrased as value, with no project named as the "bad" one.
- Citations vs. readability → keep light citations; the protocol doc owns the detail.

## Deferred follow-ups (not blocking this idea)

- Optional external links to the inspiration sources — omit from the intro for now
  (don't draw readers away); can add a "Further reading" line later.
- `docs/cli-reference.md` lags `app.go` for preflight/retro/consult/pipeline
  (codex-1) — separate doc task; the intro points to `parley help` for the live list.

## Signoffs

<!-- Each participant appends its own block: Status + one line. -->

### Signoff: claude-1 — 2026-06-20
Status: ✅ ACCEPT
Structure + guardrails capture all four lenses; the intro stays one screenful and
every claim is citable. Ready to draft FINAL.md.
### Signoff: codex-1 — 2026-06-20
Status: ACCEPT
Technical-accuracy verdict: the feature claims map to real protocol sections or shipped commands and are properly hedged on autonomy, consults, supervision, and model/vendor scope.
### Signoff: hermes-1 — 2026-06-20
Status: ACCEPT
The hook lands: anti-pattern (b) — "one model role-playing a committee" — is the real differentiator (it names the self-preference failure mode that solo-model cosplay reintroduces), and "sharp-but-not-shade" keeps it value-framed without naming targets. Pairing it with (a) audit/resume gaps covers both the multi-terminal and the faux-panel audience in one screenful.

<!-- Facilitator note (claude-1, 2026-06-20): antigravity-1's signoff append hung
(the known agy write-hang; its round-01 file completed normally). Quorum: 3/4
ACCEPT (claude-1, codex-1, hermes-1), zero blockers. antigravity-1's attribution
lens is already fully reflected in FINAL.md — the "Inspired by — adopted & adapted"
lineage block + non-endorsement disclaimer is its own round-01 proposal. Per the
established per-idea agy-waiver pattern and explicit user go-ahead to push, the idea
finalizes with antigravity-1's signoff WAIVED (not fabricated). -->
Status (antigravity-1): WAIVED — signoff append hung; round-01 input incorporated.
