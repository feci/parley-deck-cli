---
agent: claude
idea: embedded-default-protocol-resync
round: 1
date: 2026-06-13
---

## Summary

The drift is small and fully characterized: the embedded default is missing only
§12 (plus the project-deck's new `Protocol synced:` header line). The genuinely
valuable work is not the one-time patch but choosing an anti-drift mechanism so
this stops recurring, and deciding whether the bootstrap template should keep
shipping parley-deck's own roster to brand-new projects.

## Proposed approach

**Axis 1 — §12 propagation.** Carry §12 verbatim into the embedded default,
including its closing "ratified by idea `meta-protocol-change-end-to-end-pipeline`
(2026-06-02)" provenance line — that line is protocol-general (it documents §12's
own ratification, true regardless of project) and dropping it would itself be a
diff the guard must model. Insert after §11's last line with one blank-line
separator, matching the project deck exactly.

**Axis 2 — genericization.** I lean toward **keeping the embedded default a near-
verbatim mirror of the project deck**, NOT genericizing the roster now. Rationale:
(a) the cheapest, most testable invariant is "embedded ≡ deck modulo a tiny
allowlist"; every genericized field becomes a structured exception the guard and
every future amender must reason about; (b) §2's prose already says "The roster is
project-specific. Maintain it as a table here," so a new project is told to edit
it; (c) genericizing roster/workspace/host-handles is a real UX improvement for
`parley init` but it is a *separate, larger* feature (it changes what bootstrap
produces and needs its own consumer changes in `defaultCooperationForInit`), and
folding it in here muddies a clean maintenance fix. **Concrete proposal:** the
ONLY intended differences between embedded and deck are (i) the transport line —
already handled by the init swap, so keep `github-pr` in the embedded source and
let the swap do its job — and (ii) the `Protocol synced:` line, which is
project-sync provenance and should **not** be in the embedded template at all.
So: add §12, do NOT add the `Protocol synced:` line to the embedded copy.

**Axis 3 — anti-drift guard (the deliverable that earns its keep).** Ship a Go
test in `internal/protocol` that embeds/reads both files and asserts they are
identical after normalizing the known allowlist: drop any `**Protocol synced:**`
header line from the deck, and treat the transport line as equal under the
documented `github-pr`↔`local-dir` swap. The test fails loudly the next time an
amendment lands in only one copy. This is ~40 lines, no new deps, and converts an
invisible manual chore into a red build. I'd also add a one-line pointer in §7
("Changing this protocol") that protocol edits must touch both copies and the
test enforces it.

## Concerns / open questions

- Does `defaultCooperationForInit` need to also strip a `Protocol synced:` line if
  one ever leaks into the embedded copy? If we keep that line out of the embedded
  source (my proposal), the consumer needs no change — cleaner.
- Should the guard compare against the *packaged skill reference* too, or only the
  in-repo pair? I say in-repo pair only; the packaged reference is out-of-repo and
  already tracked separately (non-goal).
- Genericization (axis 2) — if the deck disagrees and wants it now, it expands
  scope into `InitWorkspace`/`defaultCooperationForInit` and the guard's allowlist
  grows. I'd rather file it as a follow-up idea than bolt it on.

## Risks

- **Guard brittleness:** if the allowlist is too literal (exact string match) a
  benign reformat breaks the build; mitigate by normalizing on a small, documented
  set of line patterns, not whole-file regex gymnastics.
- **Hidden third copy:** the packaged skill reference is a *third* lagging copy;
  fixing only the in-repo pair leaves it stale. Acceptable for this idea (non-goal)
  but must stay flagged so nobody assumes "all decks synced."
- **Verbatim-mirror decision ages poorly** if we later genericize bootstrap; the
  guard's allowlist would then need real structure. Low near-term risk; revisit
  when/if the genericization follow-up lands.
