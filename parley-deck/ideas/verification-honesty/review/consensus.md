---
idea: verification-honesty
review-cycle: 2
drafted-by: claude-1
date: 2026-06-24
outstanding_agreed_fixes: 0
blocked: false
---

## Cycle 2 resolution

All 8 cycle-1 agreed fixes (F1-F8) were applied; round-02 re-review followed.
**round-02:** codex-1 ✅ zero findings, hermes-1 ✅ zero findings (both traced the full
strict_gate close path for a fail-open and found none); antigravity-1 raised two
scan-evasion edge cases + hermes-1 a coverage gap, all resolved in fix-up cycle 2:
- **F9** [MAJOR antigravity-1] — placeholder-title heading with real content below now
  counts as a finding (`scanHasRealFinding` + `headingHasContent`).
- **F10** [MINOR antigravity-1] — scan tolerates any heading level (`##`..`######`);
  bullet/bold/prose shapes documented as out of the backstop's scope by design.
- **F11** [coverage, hermes-1] — added a test for the per-file ReadFile-error fail-closed
  veto (dangling symlink).

Trajectory converged: round-01 had 4 cross-reviewer MAJOR-class issues → round-02 had 2
single-reviewer scan-edge issues, now fixed. Outstanding agreed fixes: 0.

## Agreed fixes

Synthesized from review/round-01 (codex-1, hermes-1, antigravity-1), refutation mode.
Deduped; convergent findings merged.

- **F1 [MAJOR] (hermes-1)** — `reviewRoundHasFindings` (`internal/driver/impl.go`) fails
  OPEN on `os.ReadDir`/`os.ReadFile` errors, contradicting LE-2's "scan can only veto,
  never auto-pass." Fix: on a non-`fs.ErrNotExist` ReadDir error → veto (return true);
  on a per-file ReadFile error for a file `ReviewRoundComplete` confirmed present → veto.
- **F2 [MAJOR] (antigravity-1)** — `scanHasRealFinding` is case-sensitive on the severity
  tag and brittle to spacing (`### [critical]`, `###  [CRITICAL]` evade it → fail open).
  Fix: uppercase-normalize the tag and tolerate whitespace after `###`.
- **F3 [MAJOR] (antigravity-1)** — `scanHasRealFinding` skips a `### [SEV]` heading with
  an empty title (finding text on the next line) → evades the veto. Fix: an empty title
  counts as a finding; ignore ONLY the literal template placeholder `<title>`.
- **F4 [MAJOR] (codex-1, antigravity-1)** — LE-3 acceptance requires a stdout warning
  AND an `agent.model_diversity` event; only stdout is emitted. Fix: append a best-effort
  `agent.model_diversity` event (idea, implementer, reviewers, model, required, action)
  before warn/escalate.
- **F5 [MINOR] (codex-1, hermes-1)** — `ValidateReviewArtifact` uses `strings.Contains`,
  so the phrase in prose passes and an empty section passes. Fix: require a real
  `## Findings` heading line, and a `## Refutation attempts` heading with at least one
  non-blank content line ("show your work").
- **F6 [MINOR] (antigravity-1)** — model comparison is case-sensitive. Fix:
  `strings.EqualFold`. (Unknown implementer model stays conservative no-fire — for a
  warn we cannot assert sameness; documented.)
- **F7 [NIT] (antigravity-1)** — under `strict_gate`, a drafter omitting
  `strict_gate_clean`/`closing_review_round` spins to `MaxFixupCycles`. Fix:
  `DraftReviewConsensus` escalates immediately if the fields are absent under strict_gate.
- **F8 [tests]** — add coverage the round-01 reviewers flagged as missing: the
  `OpenReviewRound` warn + `require_model_diversity` escalate paths and the
  `agent.model_diversity` event (hermes-1 MINOR); scan edge cases (case/empty-title/
  read-error fail-closed); the heading-vs-substring validation fix.

## Deferred follow-ups

- **`sh -c` execution of `checks:` (antigravity-1 [MINOR]).** Split decision: hermes-1
  explicitly classified this as NOT a finding — `checks:` is author-controlled
  `00-prompt.md` frontmatter set at kickoff, the same trust surface as `auto_implement`
  code-writing agents and `go test`, so sandboxing only `checks:` is not meaningful while
  the rest of the workspace is writable. The concern becomes real only when an *automated
  trigger* can author a 00-prompt from untrusted input. Deferred to the automation-trust
  tier → idea `standing-loop-watch-mode` / `automation-human-brake` (LE-8/9), which must
  treat trigger-authored `checks:`/code as untrusted. Documented in IMPLEMENTATION.md.

## Dismissed findings

- **Shared `MaxFixupCycles` budget for strict-close (codex-1, hermes-1 open q).** The
  strict-close loop shares the fix-up budget; FINAL.md only requires "bounded by
  MaxFixupCycles" (termination guaranteed), which holds. A separate counter is a possible
  future tweak, not a defect. Accepted as designed.
- **`consensus.md` bare-name scan asymmetry (hermes-1 [NIT]).** Only matters if an agent
  ID is literally `consensus`/`_index`; practically nil. Accepted/documented.
- **verification-honesty does not set `strict_gate: true` on itself (hermes-1 open q).**
  Intentional: a protocol+code change is verified by `checks: go build && go test`
  (already set), which is the correct gate; strict_gate adds no runtime state to verify
  here. Authoring choice, confirmed.

## Coverage & blind spots

All three reviewers independently converged on the finding-scan robustness (F1/F2/F3)
and the missing model-diversity event (F4) — high-confidence. hermes-1 uniquely caught
the read-error fail-open (F1), the deepest invariant violation. The `checks:` shell
concern was seen by antigravity-1 and explicitly dismissed by hermes-1 → deferred, not
ignored.

## Signoffs

**Outcome: 3/3 ACCEPT (unanimous, no BLOCK) → Phase 8 complete.**

### Signoff: codex-1 — 2026-06-24
Status: ACCEPT
I accept the cycle-2 review consensus: it records zero outstanding agreed fixes after
F9/F10/F11, and my spot-check of `2dd5782..HEAD` shows the scanner now catches
placeholder-with-content and heading-level variants while unreadable review files fail
closed. The narrow verification command `go test ./internal/driver` passes.

### Signoff: hermes-1 — 2026-06-24
Status: ACCEPT
My round-02 re-review reported zero findings: the fail-closed scan (my round-01 MAJOR, F1)
holds, the strict_gate close logic has no remaining fail-open path, and all agreed fixes
F1-F8 are correctly implemented with passing tests. Fix-up cycle 2 closed my sole round-02
open question (per-file ReadFile-error coverage) with F11's dangling-symlink veto test;
`go build ./...` and `go test ./internal/driver/...` are green. With outstanding_agreed_fixes
at 0 and the scan-evasion edge cases (F9/F10) resolved, the verification-honesty
implementation is sound.

### Signoff: antigravity-1 — 2026-06-24
Status: ACCEPT
Both round-02 findings (F9 and F10) have been successfully resolved, and all tests pass
cleanly. The scan-evasion fix (F9) and the tolerance for heading levels ## to ###### (F10)
are correctly implemented in `internal/driver/impl.go`. The documented out-of-scope
decision for bullet, bold, and prose shapes is acceptable as these off-spec formats fall
outside the deterministic backstop's parsing rules.
