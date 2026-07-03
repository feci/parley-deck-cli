---
agent: antigravity-1
idea: track-aware-driver
round: 1
date: 2026-07-03
---

# Round-01: track-aware-driver design

## 1. Classifier Input Model (MVP)

I recommend a **hybrid approach (c) + (a)** for the MVP: **Author Declares (`track:` in `00-prompt.md`) + CLI Validates/Classifies.**

The source of truth remains the `track:` field in the `00-prompt.md` frontmatter. The CLI should provide a deterministic `parley classify` command that serves two purposes:
1. **Validation:** In CI or pre-round, it checks if the declared track matches the §4.0 triggers.
2. **Recommendation:** For authors, it suggests a track based on objective inputs.

### `parley classify` Interface
```bash
parley classify [IDEA_DIR] \
  [--files N] [--loc N] \
  [--security] [--irreversible] \
  [--protocol-change] [--api-break] \
  [--auto-implement] [--strict-gate] \
  [--json]
```
- **Inputs:** Explicit flags override inference. If no flags are provided, the CLI attempts to infer `files` and `loc` from `git diff` (comparing HEAD to the integration branch).
- **Output:** The recommended track string (`fast`, `standard`, `deliberation`) and a list of triggers fired.
- **Logic:** Implements §4.0 "deliberation-first, then fast, else standard; on doubt → stricter" precisely. This is deterministic and script-checkable.

## 2. Track → Config Mapping

I propose adding `Track string` and `RequiredReviewerCount int` to `driver.Config` (`internal/driver/driver.go`). The mapping from §4.0 is derived in the `app` layer:

| §4.0 Field | `fast` | `standard` | `deliberation` | Threading / Location |
| :--- | :--- | :--- | :--- | :--- |
| `CrossReviewRounds` | `0` | `2` | `-1` (unbounded) | `driver.Config.CrossReviewRounds` |
| `MaxFixupCycles` | `1` | `2` | `3` (current default) | `driver.Config.MaxFixupCycles` |
| `RequiredReviewerCount` | `1` | `2` | `len(Participants)-1` | `driver.Config.RequiredReviewerCount` |
| `AutoAdvance` | Full | Human @ FINAL | Human @ Transitions | `driver.Config.Auto` |
| `Timeout` | `~5m` | `~15m` | `~30m` | `internal/config/runtime.go` |

### Key Code Changes:
- **`internal/driver/driver.go`**:
    - `advanceRound`: Update the check `c.CurrentRound >= 1+d.cfg.CrossReviewRounds` to handle `-1` (unbounded) by never bypassing to consensus.
    - **Fast-Track Collapse**: In `advanceRound` for `fast`, if rounds are complete, skip `PhaseConsensus` and transition directly to `PhaseFinal`, instructing `Consensus.Draft` to produce a `FINAL.md` with embedded signoffs.
- **`internal/driver/impl.go`**:
    - `advanceReview`: Replace the hardcoded `rs.ReviewerCount < 2` with `rs.ReviewerCount < d.cfg.RequiredReviewerCount`.
- **`internal/app/app.go`**:
    - In `continueAuto` and `runTask`, introduce a `deriveTrackConfig` helper that calls `driver.ReadTrack(ideaDir)` and applies the table logic to populate `driver.Config`.
- **`internal/config/runtime.go`**:
    - Seed `[defaults.timeouts]` with track-aware defaults during `parley init`.

### Gaps & New Fields:
- **`RequiredReviewerCount`**: Needed to enforce §4.0 "2 reviewers" vs "all non-implementers".
- **`HumanGatePhases`**: A new slice in `Config` to enforce §4.0 "human gate at each transition" for `deliberation`.

## 3. Invariant Enforcement

The driver must protect the §1 non-solo and §4.0 refutation invariants.

**Where:** `driver.New()` and a new `Config.Validate()` called at the start of `Advance`.
**Logic:**
- **Non-Solo Hard-Stop:** If `RequiredReviewerCount < 1`, return error (prevents solo completion).
- **Refutation Discipline:** The **deterministic finding scan** (`reviewRoundHasFindings` in `internal/driver/impl.go`) must be non-bypassable on all tracks.
- **Track-Quorum Consistency:** If `Track == standard` but `Participants` only contains 2 agents (implementer + 1 reviewer), `RequiredReviewerCount` must be 1 (as per §4.0 trigger: "With only two participants, standard's '2 reviewers' degrades to 1").
- **Safety Escalation:** If `Track` is `deliberation` but `CrossReviewRounds` is manually overridden to `0` in `00-prompt.md`, `driver.New` must reject/escalate.

## 4. MVP Slicing (Safety + Value / Risk)

### Slice 1: Core Routing (High Safety, High Value)
- `ReadTrack` in `transport.go`.
- `Track` and `RequiredReviewerCount` in `driver.Config`.
- Track derivation logic in `app.go`.
- Invariant enforcement in `driver.New`.
- *Outcome:* Deterministic enforcement of §4.0 for manually labeled ideas.

### Slice 2: Classifier CLI (Medium Risk)
- `parley classify` command.
- Git diff inference for `files`/`loc`.
- *Outcome:* Automated track recommendation/validation.

### Slice 3: Refined UI/UX (Low Risk)
- "Collapsed FINAL" implementation (skipping PhaseConsensus for `fast`).
- `HumanGatePhases` enforcement.

## 5. Backward Compatibility & Test Plan

- **Compatibility:** Absence of `track:` defaults to `standard`. `standard` reproduces today's defaults (2 reviewers, 2 fix-up cycles).
- **Test Plan:**
    - `internal/driver/driver_test.go`: Add `TestTrackBehaviors` to verify `fast` skips rounds and `deliberation` respects human gates.
    - `internal/driver/impl_test.go`: Add `TestReviewerQuorum` to verify `RequiredReviewerCount` enforcement.
    - `internal/app/app_test.go`: Verify `deriveTrackConfig` mapping.

## 6. Risks

1. **`auto_implement` Loophole**: §4.0 lists `auto_implement` as a `deliberation` trigger, but `fast` track has code-specific fields. If `auto_implement` always forces `deliberation`, `fast` is effectively unusable for small code fixes.
    - *Resolution*: Treat the classifier triggers as "deliberation if (auto_implement AND NOT fast_eligible)".
2. **Silent Behavior Change**: Existing ideas (now `standard`) that relied on unbounded rounds will now cap at 2.
    - *Mitigation*: Ensure the escalation at round 2 is clear and preserves the option to manually reopen or upgrade to `deliberation`.
3. **Model Diversity**: `fast` track's 1-reviewer minimum increases rubber-stamping risk.
    - *Mitigation*: Keep `LE-3` (model diversity check) as a hard gate if `require_model_diversity: true`.
