---
agent: codex-1
idea: track-aware-driver
round: 1
date: 2026-07-03
---

## Summary

The MVP should be conservative: add deterministic `ReadTrack`, add an explicit `parley classify` command with fail-closed inputs, derive a track policy in the app layer before building `driver.Config`, and hard-reject invariant violations in the driver/review path. Do not make git-diff inference the source of truth in the first slice.

I recommend a hybrid of option (a) and option (c): author/deck files declare `track:`, and `parley classify` provides a script-checkable floor that can reject under-tier declarations. Option (b), diff inference, should come later as an advisory helper only.

There is one important mismatch to resolve deliberately: `00-prompt.md` says `standard`/absent must behave exactly like today, but the §4.0 table says standard fix-up is capped at 2 while the current driver default is 3 (`internal/driver/driver.go:91-93`). My MVP keeps absent/standard compatible and gates the protocol-target cap behind an explicit follow-up or compatibility decision.

## 1. Classifier input model

Use a new command:

```text
parley classify [--json] [--declared fast|standard|deliberation]
  --files N
  --loc N
  --reversibility reversible|irreversible|unknown
  --mechanically-verifiable yes|no|unknown
  [--risk protocol-change]
  [--risk security]
  [--risk auth]
  [--risk secrets]
  [--risk payments]
  [--risk privacy]
  [--risk production-infra]
  [--risk data-migration]
  [--risk destructive]
  [--risk strict-gate]
  [--risk auto-implement]
  [--risk pipeline]
  [--risk action-block]
  [--risk public-api-break]
  [--risk persisted-schema-break]
  [--unknown-risk security|privacy|production-infra|data-migration|pipeline|api-break|schema-break|...]
```

Rules:

- Evaluate deliberation first, matching §4.0: any listed `--risk`, `--files > 15`, or `--loc > 1000` returns `deliberation`.
- Any `--unknown-risk` for a deliberation trigger also returns `deliberation`; this is the "on doubt -> stricter" rule in `parley-deck/COOPERATION.md:185-190`.
- If no deliberation trigger fires, return `fast` only when `--files <= 5`, `--loc <= 300`, `--reversibility reversible`, and `--mechanically-verifiable yes`.
- Otherwise return `standard`.
- With `--declared`, exit 0 when the declared track is at least as strict as the classifier result; exit 4 and print the computed floor when the declaration is under-tiered. Usage errors stay exit 2.
- `--json` should print at least `{ "track": "...", "declared": "...", "valid": true|false, "reasons": [...] }`; plain mode prints only the computed track so scripts can do `track=$(parley classify ...)`.

I would not use git diff inference as MVP truth. Before implementation, there may be no diff; after implementation, the diff may miss external effects, schema semantics, secrets exposure, or deployment risk. A later `parley classify --from-diff BASE...HEAD` can pre-fill the numeric flags and obvious risk flags, but it should still output an auditable classifier input record rather than silently deciding.

Threading:

- Add a `classify` case in `internal/app/app.go:50-99`.
- Put pure classifier logic in a small package, e.g. `internal/track`, so tests do not need CLI plumbing.
- `parley run` can later accept `--track` plus classifier flags; `runcontrol.CreateOptions` in `internal/runcontrol/runcontrol.go:17-25` needs a `Track` field so `protocol.CreateIdeaWithExclusions` can write the frontmatter.
- Current `protocol.CreateIdeaWithExclusions` does not write `track:` (`internal/protocol/workspace.go:133-153`), despite the prompt saying the template already carries it. The MVP should add `track: standard` there for newly created ideas.

## 2. Track -> Config mapping

Add:

```go
type Track string

const (
    TrackFast Track = "fast"
    TrackStandard Track = "standard"
    TrackDeliberation Track = "deliberation"
)

func ReadTrack(ideaDir string) Track
```

`ReadTrack` belongs in `internal/driver/transport.go` beside `ReadAutoImplement`, `ReadStrictGate`, and `ReadRequireModelDiversity` (`internal/driver/transport.go:43-70`). It should use the existing `readFrontmatterField` helper from `internal/driver/cursor.go:282-314`; absent, empty, quoted, or unknown values normalize to `standard`.

Recommended policy shape:

```go
type TrackPolicy struct {
    Track Track
    CrossReviewRounds int
    MaxRounds int
    MaxFixupCycles int
    RequiredReviewers int
    ReviewerMode ReviewerMode // exact-n or all-non-implementers
    RequireReviewerModelDiversity bool
    AgentTimeout time.Duration
}
```

Exact mapping I would implement:

| Track | Existing Config values | Additional policy needed |
| --- | --- | --- |
| `fast` | `CrossReviewRounds: 0`; `MaxFixupCycles: 1`; `MaxRounds`: not enough to prevent a BLOCK reopen cleanly | `RequiredReviewers: 1`; `ReviewerMode: exact-n`; `RequireReviewerModelDiversity: true`; `AgentTimeout: 5m`; collapsed `FINAL.md` mode; fast review consensus = one reviewer |
| `standard` | Compatibility MVP: `CrossReviewRounds: min(ReadCrossReviewRounds, 2)` with absent default 1; `MaxFixupCycles: 0` so `driver.New` keeps today's default 3; `MaxRounds: 0` so current default 4 remains | Protocol target later: `MaxFixupCycles: 2` and `MaxRounds: 2`; `RequiredReviewers: min(2, distinct non-implementers)` to honor the two-participant degradation in `COOPERATION.md:225-226`; `AgentTimeout: 15m`; review signoff set = reviewers who reviewed |
| `deliberation` | Preserve current behavior: `CrossReviewRounds: ReadCrossReviewRounds` default 1 (`internal/driver/transport.go:31-40`); `MaxFixupCycles: 0` -> current default 3; `MaxRounds: 0` -> current default 4 | `ReviewerMode: all-non-implementers`; `RequiredReviewers: len(non-implementers)`; `AgentTimeout: 30m`; review/design signoff set = all participants |

Where this threads:

- App-layer derivation should happen at every current `driver.Config` construction site: `continueAuto` (`internal/app/app.go:1154-1171`), no-TUI auto-drive (`internal/app/app.go:1827-1844`), and TUI auto-drive (`internal/app/app.go:1881-1898`).
- `driver.Config` already has `CrossReviewRounds`, `MaxRounds`, `MaxFixupCycles`, `AutoImplement`, and `StrictGate` (`internal/driver/driver.go:40-75`).
- `driver.New` currently defaults `MaxRounds` to 4, negative `CrossReviewRounds` to 1, and `MaxFixupCycles` to 3 (`internal/driver/driver.go:83-98`).
- Reviewer selection currently lives in `newDriverImplOps`, which dedupes all distinct non-implementers and uses all of them (`internal/app/driver_impl.go:38-63`). That function needs a review policy parameter.
- Reviewer-count completion enforcement is currently hard-coded as `<2` only under `AutoImplement` (`internal/driver/impl.go:232-243`). Replace it with `RequiredReviewers`.
- Timeouts currently flow through `runner.Options.Timeout` (`internal/runner/runner.go:22-29`) and agent `TimeoutMS`; `timeoutForAgent` prefers agent-specific `TimeoutMS` over `Options.Timeout` (`internal/runner/runner.go:1113-1116`). A track timeout therefore must either mutate selected discoveries for this run or become a first-class runner timeout policy.
- `internal/config/runtime.go` has only flat `signoff_ms`, `round_ms`, `review_ms`, and `deep_reasoning_ms` fields (`internal/config/runtime.go:38-43`) and a central default of 20 minutes (`internal/config/runtime.go:282-286`). It has no per-track timeout shape today. Add nested TOML like `[defaults.timeouts.fast] agent_ms = 300000`, `[defaults.timeouts.standard] agent_ms = 900000`, `[defaults.timeouts.deliberation] agent_ms = 1800000`, with the existing flat fields as legacy fallbacks.

§4.0 behavior with no current Config field:

- `track` itself.
- Fast skipping the readiness ping while retaining the non-solo hard stop. Preflight has `NoPing` (`internal/app/preflight.go:231-233`) but no track input.
- Collapsed fast `FINAL.md` with embedded signoffs. Current driver always drafts `consensus.md`, requests signoffs, then drafts `FINAL.md` (`internal/driver/consensus.go:43-88`).
- Standard's "separate, drafted simultaneously" consensus/FINAL behavior. Current flow is sequential.
- Review signer sets by track. `consensus.Status` computes missing signoffs against all idea participants (`internal/consensus/consensus.go:411-415`), including review consensus.
- Fast "one reviewer = review consensus" and standard "reviewers who reviewed sign off"; current review signoff code still uses the consensus package's participant set.
- Per-phase human-gate policy. `driver.Config.Auto` is all-or-nothing (`internal/driver/driver.go:50`, `internal/driver/driver.go:135-144`), and `AutoImplement` only gates code-writing (`internal/driver/impl.go:71-73`).
- Fast model-diverse reviewer as a hard requirement. The current model-diversity check only warns unless `require_model_diversity: true` is set (`internal/app/driver_impl.go:135-164`).

## 3. Invariant enforcement

The driver should validate invariants before taking a phase action, not rely only on `parley run` preflight.

Add `ValidateTrackConfig(cfg Config) error` in `internal/driver` and call it at the top of `Advance` after `Rebuild` (`internal/driver/driver.go:150-158`). On violation, return `ActionEscalated` and an error. This keeps `New`'s signature stable.

Hard rejects:

- Unique `Participants` must be at least 2. Preflight already hard-stops selected solo runs (`internal/app/preflight.go:334-338`), but `continue --auto`, tests, or hand-authored ideas can bypass preflight.
- If `Impl != nil`, `RequiredReviewers` must be at least 1. A review path with zero reviewers is invalid on every track.
- `RequiredReviewers` must be satisfiable by distinct non-implementers after the §4.0 two-participant standard degradation is applied. If not, escalate and ask for a stricter track or a larger roster.
- Fast must not proceed unless the selected reviewer is model-diverse from the implementer when model data is known; if model data is unknown, fail closed by upgrading/rejecting fast rather than pretending diversity.
- Do not add a flag that can disable review refutation. The review prompt already requires refutation-default posture and a `## Refutation attempts` section (`internal/runner/phase58.go:234-270`), and validation rejects missing or empty refutation attempts (`internal/runner/phase58_le_test.go:25-55`). Keep that as a non-optional artifact validator.

Implementation points:

- Extend `driver.Config` with `Track Track` and `RequiredReviewers int`.
- Replace the current LE-11 hard-coded reviewer count in `advanceReview` (`internal/driver/impl.go:232-243`) with `if rs.ReviewerCount < d.cfg.RequiredReviewers`.
- Pass review policy into `newDriverImplOps`; enforce `len(o.reviewers) >= RequiredReviewers` in `OpenReviewRound` before launch. The existing zero-reviewer check at `internal/app/driver_impl.go:220-223` becomes the first case of that broader validation.
- Keep `ReviewRoundComplete` validating every selected review artifact (`internal/app/driver_impl.go:257-269`); because `runner.ValidateReviewArtifact` enforces refutation attempts, a non-refuting review never completes.

## 4. MVP slicing

1. Read and record track without behavior change. Add `ReadTrack`, normalizer tests, `track: standard` in new 00-prompt templates, and show track in run/status output. Value is high, risk is very low.
2. Add `parley classify` pure command and under-tier validation. This gives deterministic, scriptable routing without touching the phase engine.
3. Add invariant validation in the driver: unique participant count, `RequiredReviewers >= 1`, refutation remains mandatory, and reviewer count is policy-driven instead of hard-coded. This is safety-critical and should land before fast behavior.
4. Apply existing-field policy for fast only: `CrossReviewRounds=0`, `MaxFixupCycles=1`, one model-diverse reviewer. This gives the largest speed value while still preserving standard/absent behavior.
5. Apply reviewer selection policies for standard/deliberation. Standard selects two when available; deliberation keeps all non-implementers. Add visible run events listing selected reviewers so reduced reviewer sets are auditable.
6. Add per-track timeout config and seeding in `EnsureCentralDefault` (`internal/config/runtime.go:243-305`). This is useful but not a correctness prerequisite.
7. Implement the higher-risk lifecycle changes: fast collapsed `FINAL.md`, review signoff sets by track, standard simultaneous consensus/FINAL, and per-phase human gates. These require consensus package/API changes and should not be in the first shippable slice.

## 5. Backward compatibility and test plan

Backward compatibility rule: an absent `track:` and `track: standard` must produce the same effective behavior as today's driver in the MVP. That means default cross-review remains 1 (`ReadCrossReviewRounds` default at `internal/driver/transport.go:31-40`) and default fix-up remains 3 through `driver.New` (`internal/driver/driver.go:91-93`) until the standard cap-2 compatibility question is explicitly resolved.

Tests:

- `internal/driver/transport_test.go`: `ReadTrack` returns `standard` for absent, empty, quoted empty, and unknown values; normalizes `fast`, `standard`, `deliberation`.
- Classifier table tests: every deliberation trigger wins over fast-looking numeric inputs; unknown deliberation risk returns `deliberation`; `files=5, loc=300, reversible, mechanically-verifiable` returns `fast`; boundary gaps return `standard`; `--declared fast` under a standard/deliberation floor exits 4.
- App wiring tests around the three `driver.Config` construction sites in `internal/app/app.go:1154-1171`, `internal/app/app.go:1827-1844`, and `internal/app/app.go:1881-1898`, asserting absent/standard policy equals the pre-change config.
- `internal/protocol/workspace_test.go`: newly created prompts include `track: standard` without changing participant parsing.
- `internal/app/preflight_test.go`: fast track sets `NoPing`/presence-only behavior but still hard-fails when the selected set has fewer than two available participants. Existing non-solo tests at `internal/app/preflight_test.go:525-543` should still pass.
- `internal/app/driver_impl_le_test.go`: reviewer dedupe still works; fast chooses exactly one model-diverse reviewer or errors; standard chooses two distinct reviewers when available and degrades to one with only two participants; deliberation chooses all non-implementers.
- `internal/driver/close_integrity_test.go`: reviewer-count guard uses `RequiredReviewers`, so fast can complete with one reviewer, standard cannot complete with one when two are available, and deliberation requires all non-implementers.
- `internal/runner/phase58_le_test.go`: keep and extend refutation validation tests so track policy cannot accidentally allow an empty review.
- Regression tests for `go test ./...` plus focused driver re-entry tests: fast with `CrossReviewRounds=0` does not create `round-02`; standard/absent still follows today's round advancement.

## 6. Risks

- The standard cap conflict is real: §4.0 says standard fix-up cap 2, but today's code defaults to 3. Changing it would violate the prompt's backward-compat constraint; not changing it leaves part of §4.0 unenforced. This should be an explicit consensus decision, not an accidental side effect.
- Fast collapsed `FINAL.md` is easy to get wrong because it cuts across consensus drafting, signoff validation, and finalization. Keep it out of the first safety slice.
- Agent timeout precedence can make track timeouts appear configured but ineffective: agent `TimeoutMS` overrides `runner.Options.Timeout`. Tests must prove the selected run actually uses the per-track value.
- Reviewer reduction can hide expertise. The deterministic selection policy should be visible in event logs and should fail closed when model diversity cannot be established for fast.
- Diff inference can under-classify risk. It should remain advisory until the repo has reliable metadata for migrations, API breaks, infra, and secrets/privacy surfaces.
- `auto_implement` is listed as a deliberation trigger in §4.0, yet §4.0 still defines fast/standard Phase 6-8 behavior. The driver design should treat auto-implemented code as deliberation unless the protocol later clarifies a safe standard/fast auto-implementation path.
