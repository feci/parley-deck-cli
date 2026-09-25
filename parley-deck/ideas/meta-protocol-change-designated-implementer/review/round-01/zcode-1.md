---
agent: zcode-1
idea: meta-protocol-change-designated-implementer
review-round: 1
date: 2026-09-25
reviewed-commit: 0893989 (CLI, baseline e4640bf) + bf7e049 (skill, baseline 8161e5e)
---

## Summary

The implementation is a faithful, unusually well-instrumented rendering of the frozen FINAL: every mechanism rule I traced (four-state parse, four-rank chain, gate siting, tier split, dormant fall-throughs, event identity, review-exclusion parity, drafter preference) matches the signed spec, the unset path is pinned byte-identical by tests that I mutation-tested and independently re-ran, and the three-copy/protocol/skill consistency checks (AC-1/2/3) reproduce exactly. I broke the waiver/reassignment confirmation parsers in three ways (negated confirmations and a mismatched-id waiver are silently accepted as valid owner records) — one MAJOR; the remaining findings are MINOR documentation/robustness items, two of them the already-disclosed IMPLEMENTATION.md metadata defects.

## Protocol context attestation

`parley protocol packet --dir . --phase 6 --track deliberation --idea meta-protocol-change-designated-implementer --json`:

- `context_mode: full`, `source_sha256 = packet_sha256 = c749218255c96c4efeecc8d598abc6192f195a294291eff6c505096f0091568f`, `fallback_reason` ABSENT (top-level keys enumerated).
- Source: `parley-deck/COOPERATION.md` (live source authority, role `source`, transport `github-pr`, 114,771 bytes). Independently re-hashed the working-tree file: `c7492182…568f`, byte-identical to commit `0893989` (working tree clean for that path).
- Rendered body read in full for Phases 5–8, §15, §0, §9.0; the full 1,400-line packet body was the version consulted.

## Review method (independence)

No implementer claim was trusted. All code was read from isolated read-only extractions pinned to the exact commits under `/private/tmp/zcode1-review/` (CLI `git archive 0893989`, skill `git archive bf7e049` — no shared-mount worktree contention with the peer reviewer). Every check below was executed by me on local disk unless marked otherwise. Baseline comparisons use `git diff e4640bf..0893989` (CLI) and `git diff 8161e5e..bf7e049` (skill).

## Refutation attempts

Per FINAL `## Observable acceptance criteria`. Provenance: PRIMARY (command + result) unless stated.

- **AC-1 (three-copy fidelity)** — `tail -n +160` of the deck copy and `tail -n +153` of `internal/protocol/defaults/COOPERATION.md` and the skill `references/COOPERATION.md` all hash to `da9704d57888efe549a2c781ba3bcdf34012c3caa81de9977fff3a1bc7e8eaf7` (matches IMPLEMENTATION's claim). Pairwise `diff` yields exactly the three pre-existing project-zone hunks (header 3 / 5–7; deck 154–159) and nothing else. TRIED TO BREAK: defaults-vs-skill diff for any fourth hunk — only header lines. **HOLDS.**
- **AC-2 (changelog)** — `parley-deck/meta/protocol-changelog.md` head entry carries date, `Idea:`, `Drafted by:`, `Summary:`, honestly marked `**Status: UNRELEASED.**`. **HOLDS.**
- **AC-3 (skill companion consistency)** — `grep -rn "default implementer" skills/ lib/ bin/` in the skill worktree at `bf7e049`: 0 hits. Cross-check: the only remaining Phase-5 "FINAL drafter" restatement in the skill tree outside the guarded copy is the amended `ROSTER_AND_PROTOCOL.md:64` line itself, whose chain now matches the core copy. **HOLDS.**
- **AC-4 (unset-path invariance)** — re-ran, unmodified and passing: `TestResolveImplementerFromRoleMetadata` (app), `TestUnresolvableImplementerExpectsEveryone` and `TestFinalDrafterIsTheFallbackImplementer` (consensus); none of the three test files is in the 16-file implementation diff. `TestUnsetPathIsByteIdentical` drives a REAL dispatch through a fake implementer binary and asserts stdout is exactly `driver: implementing via zz-first ...\n` with no `agent.implementer_resolved` event — I read the driver diff and confirmed the extra line/event are gated on `implDesignated` behind the byte-identical baseline line. `ExpectedRoundParticipants` parity is pinned by `TestExpectedRoundParticipantsUnaffectedByDesignation` across all four states. **HOLDS.**
- **AC-5 (four-state parse)** — `TestImplementerDesignationFourStates` covers absent / empty / whitespace / `none` / casing / padding / quoted forms / the R4 malformed inline-comment value kept literal. I additionally probed quoted `"none"` (→ `DesignationNone`) and quoted whitespace (`"  "` → `DesignationSet` with a garbage id → fails the eligibility gate downstream): both fail closed. **HOLDS.**
- **AC-6 (`none` does not collapse)** — `TestNoneSuppressesTier3` distinguishes `none` (suppress tier 3, source `none`, no gate) from absent (tier 3 fires). **Mutation test (mine):** deleting the `EqualFold("none")` branch from `implementer.go` makes the test suite FAIL (`driver_designation_test.go:245 … got implementer="aa-first" source=""`), then restored from the pinned commit and re-verified green — the pin is real, per T-2's "fails if the rule is removed" requirement. **HOLDS.**
- **AC-7 (config precedence)** — `TestDefaultImplementerLayerPrecedence` spans two REAL layers (central + deck agents.toml): higher non-empty wins; `""` at the higher layer does not clear; `"none"` wins the merge. Matches `mergeDefaults` code read. **HOLDS.**
- **AC-8 (tier-2 gate + exits)** — `TestTier2UnavailabilityGateAndExits` blocks on ping failure via `dispatchErr` escalated by `Implement` AND `Fixup`, and clears via re-designate / confirmed waiver / `none`; unconfirmed waiver rejected. Gate rides dispatch actions only, roleErr stays empty — matching the recorded F2 decision (validity ≠ availability). **HOLDS** (but see **[MAJOR-1]** — the waiver parser's "confirmed" check accepts negations).
- **AC-9 (five dormant tier-3 states)** — `TestDormantTier3StatesFallThrough` covers absent / non-participant / `none` / own declared non-participating facilitator / ping-failed, each asserting no gate, fall-through to today's id, and the distinguishing source (`fall-through-inapplicable` / `fall-through-unavailable` / `none`); absent emits nothing. I verified the R20 tier split both ways (same facilitator id at tier 2 hard-gates). **HOLDS.**
- **AC-10 (invalid tier-3 still fails)** — `TestMalformedTier3HardFails`: `default_implementer = "zz-impl # copied from a comment"` hard-fails via roleErr and blocks `Implement`. The malformed/inapplicable boundary (whitespace-containing vs whitespace-free non-member) is the implementer's recorded F-decision and is applied consistently. **HOLDS.**
- **AC-11 (pin conflict escalates)** — `TestPinDesignationConflictEscalates`: disagreement → roleErr naming the `implementer_reassigned:` exit; unconfirmed and wrong-pair records rejected; confirmed record re-pins; agreement resolves as pin (R27). **HOLDS** (but see **[MAJOR-1]** for the negated-confirmation hole in this same parser).
- **AC-12 (degenerate drafter fallback)** — `TestDrafterSeparationPreference`: designee-only-drafter falls back to the designee rather than failing; empty preference preserves today's order. **HOLDS.**
- **AC-13 (gate scoping)** — `TestDesignationGateIsScopedToTheIdea`: idea A's defect gates A only; no all-ideas iteration anywhere in the delta (the preflight all-ideas trap is not copied — `preflight.go` untouched). **HOLDS.**
- **AC-14 (two-participant warning)** — `TestTwoParticipantDesignatedKickoffWarnsNotBlocks` at tier 2 and tier 3: warns, no gate. **HOLDS.**
- **AC-15 (review exclusion reads the pin)** — `TestReviewExclusionReadsThePinNotTheDesignation`: designation X + pin Y → review excludes Y; FINAL-drafter and unresolvable fallbacks unaffected; exported `ExpectedRoundParticipants` equally designation-blind. Code read confirms the consensus side calls the shared chain with ONLY `LegacyImplementerCandidates()` and the raw participants list — tiers 2/3 cannot enter it. **HOLDS.**
- **AC-16 (inertness)** — `grep -rn "impl-claim" --include="*.go" .` → 0 hits (run by me). New `00-prompt.md` references in the delta are `ReadFrontmatter` READS only (`driver_impl.go:175,184`; `:493` is pre-existing checks-contract code, baseline `:262`); writers remain the pre-existing three. **HOLDS.**
- **AC-17 (whole-tree health)** — independently on local disk at `0893989`: `go build ./...` OK; `go vet ./...` OK; `gofmt -l` empty on all ten changed Go files; `go test ./internal/protocol/ ./internal/config/ ./internal/consensus/` all `ok`; the full designation suite in `internal/app` PASS (18 tests incl. `TestUnsetPathIsByteIdentical` end-to-end). Full `go test ./...` was still running at artifact time under CPU contention with the peer reviewer's concurrent suite (see **Honest limitations**); implementer's 505s/605s shared-mount figures are consistent with what I observed. **PARTIALLY INDEPENDENTLY CONFIRMED — see limitations.**
- **AC-18 (end-to-end designation demonstration)** — `TestTier2DesignationDispatchesDesignee` dispatches the second-listed designee through a real fake-implementer run and verifies the event payload `{idea, implementer, source:"designation"}`; `TestTwoLayerGlobalDefaultDispatchesHigherLayer` sets the key at two real layers and dispatches the higher layer's id with `source:"global-default"`. **HOLDS.**
- **AC-19 (ships UNSET)** — `TestCentralDefaultTemplateShipsImplementerUnset` (commented-out key, active line forbidden) and `TestDeckCreatedByToolingHasNoDefaultImplementer` (no deck TOML carries the key); template text read — the key is emitted `#`-commented with the deviation sentence. The template-scope fix after the prose false-positive is recorded in IMPLEMENTATION. **HOLDS.**
- **AC-20 (attended boundaries)** — the 16-file CLI diff and 2-file skill diff touch no publication command, no pty allocation, no channel action, no exec.Command addition (grep over added lines). `parley protocol publish` code untouched. **HOLDS** for this delta.
- **AC-21 (independent review actually happened)** — this file and claude-1's are the two non-implementer reviews with populated `## Refutation attempts`; the goal-done check is for a fresh non-implementer at close. **In progress by construction — not for the implementer to satisfy.**

**Adversarial break attempts beyond the implementer's tests** (PRIMARY, `go run` scratch harness against the exported parsers at `0893989`): (1) `implementer_waived: zz-impl — NOT confirmed yet, pending owner` → `ImplementerWaived` returns **true**; (2) `implementer_reassigned: aa-first to zz-impl — NOT confirmed` → `ParseImplementerReassignment` returns **true** (pin retired); (3) `implementer_waived: kimi-10 — offline — confirmed 2026-09-25` with designee `kimi-1` → **true** (substring id containment). These are the substance of [MAJOR-1].

## Findings

### [MAJOR] Waiver/reassignment "confirmed" parsers accept negated and mismatched-id records (fail-open)

`internal/protocol/implementer.go` — `ImplementerWaived` (substring `strings.Contains(parts[0], id)` + `strings.Contains(v, "confirmed")`) and `ParseImplementerReassignment` (same marker check) mis-read the §9.0-shaped owner records in three reproducible ways (commands and outputs under Refutation attempts):

1. A record that **explicitly denies confirmation** ("NOT confirmed yet") is accepted as a confirmed waiver — silently converting the R18 tier-2 blocking gate into a fall-through with no notice. The code's own comment ("anything weaker is not a waiver") is contradicted by the code.
2. The same negation retires a pin via the reassignment path — the exact self-appointment-with-paper-trail class R28 exists to prevent; an unconfirmed state is treated as the owner-recorded act.
3. A waiver naming a **different agent** (`kimi-10`) clears the gate for designee `kimi-1` because the id check is substring containment, not equality — a plausible honest-input failure on numbered rosters, and precisely the "typo indistinguishable from intention" failure mode R3 cites as what this idea exists to remove.

Why it blocks-before-merge for me: these parsers are the machine readers of the mechanism's only owner-confirmation records; every other ambiguity in this design was resolved fail-closed, and these two resolve fail-open, silently, on the authorization path. Suggested fix: require the id segment to **equal** the designee id (after trim), require a strict trailing confirmation segment matching `confirmed <date>` (e.g. `(?i)^confirmed \d{4}-\d{2}-\d{2}$` against the last `—` segment), and reject a preceding negation ("not confirmed", "unconfirmed"). Add the three cases above as tests (they currently pass wrongly).

### [MINOR] Tier-3 live read silently swallows layered-config errors

`globalDefaultImplementer` (`internal/app/driver_impl.go`) returns `""` when `config.LoadDefaults` errors, so a malformed `agents.toml` (or a bad `PARLEY_HEADLESS_AGENT_CONFIG`) silently disables the owner's standing default: dispatch falls to the positional tail with **no notice and no event**, the same invisibility class R48 records for the legacy path, now on the designation path the owner opted into. The fail-open choice is documented in a code comment, but nothing surfaces it at dispatch time. Suggested fix: when `LoadDefaults` errors, emit the designation-path one-line notice naming the config error (it already has the R45 line machinery), or at minimum record the cause in the §7 changelog/protocol text as a known fail-open.

### [MINOR] IMPLEMENTATION.md frontmatter `head-commit` names the baseline, not the implementation

`head-commit: e4640bf` is the FINAL-publication commit; the implementation is CLI `0893989` + skill `bf7e049` (both disclosed in my review brief and by the organizer). The Phase-5 template's `head-commit` is the re-entry/review pin reviewers and tooling read; a stale value misleads re-entry tooling and release staging (and the skill-side commit appears only in prose). This is the known documentary discrepancy; my disposition: it is an objective defect in the outcome record, fixable by kimi-1 in fix-up as a one-line author edit (`head-commit: 0893989`, plus the skill commit noted). It does not change the code under review — I reviewed the actual commits.

### [MINOR] Duplicate `## Validation evidence` section in IMPLEMENTATION.md

The file carries TWO `## Validation evidence` headings — a populated one after `## Progress` and a placeholder after `## Surprises & Discoveries`. Beyond template shape, this risks future evidence landing in the wrong section and ambiguity for any goal-done/driver reader that anchors the section. Also fixable by the implementer in fix-up (delete the placeholder).

### [MINOR] Amended protocol rank-4 text re-asserts the imprecise "FINAL drafter" fallback

The new §4 Phase-5 chain sentence (`parley-deck/COOPERATION.md:449`, identical in all copies and mirrored in the ROSTER line) states rank 4 as "the fallback — **the FINAL drafter** (same agent as Phase 4)", while the code (correctly, per R13/R49) implements rank 4 as `FINAL.md{implementer, drafted-by}` → positional tail. The signed FINAL's own Purpose section bars repeating the drafter-implements claim unscoped, and the new sentence replaced the old pre-existing one, so the amended authority now carries the simplification as the *only* rank-4 documentation. Behavior is correct and R49 defers the repair — this is a documentation-precision gap, not a code defect. Suggested: a one-line precision edit in all three copies + the ROSTER line ("(4) today's chain — `FINAL.md`'s recorded implementer/drafter, else the first eligible participant"), or an explicit cross-reference to the deferred follow-up slug; consistent fix-up or a recorded deferral both acceptable.

### [NIT] Re-entry comparison escalates on a source-only change

`checkImplementerReentry` errors when the recorded and current **source** differ even if the implementer id is identical (e.g., the same agent re-designated via `default_implementer` after a per-idea designation was removed). Fails closed, so safe, but it manufactures a blocking escalation where dispatch identity is unchanged. Consider comparing only the implementer id, or documenting the strictness.

## Disposition of disclosed items (independent evaluation)

- **R57 skill companion finding (FINAL's disclosed unreviewed item): I CONCUR.** PRIMARY evidence: the pre-edit line existed at the skill baseline (`git show 8161e5e:skills/parley-deck/references/ROSTER_AND_PROTOCOL.md`, line 64 restated "default implementer is the FINAL drafter unless another participant claims it"), it contradicted the amended core chain, the consensus's AD-14 mandated both identical hunks and the drift grep, and the one-line edit brings the companion into consistency without touching installer, manifest or channels (2-file skill diff verified). This is a correct consequence of the signed consistency requirement, not scope creep; FINAL's own disclosure invited exactly this check.
- **IMPLEMENTATION frontmatter naming e4640bf + duplicate Validation placeholder:** treated as objective MINOR findings above, not as permission to narrow scope; the actual commits `0893989`/`bf7e049` were reviewed.
- **Windows (R52):** the delta adds no build-tagged file and no platform surface (name-only diff checked); the baseline Windows independent-evidence refusal is untouched. Factual residual risk for the deferred `windows-portability` run: the designation mechanism's stdout/event surfaces are untested on Windows terminals, and the slow-fork test profile noted in IMPLEMENTATION will be worse there. No Windows work was expanded here.

## Verdicts (§15)

All verdicts below are mine, written in my own review file, tagged per §15.2:

- AC-1, AC-2, AC-3, AC-5, AC-6, AC-7, AC-9, AC-10, AC-11, AC-12, AC-13, AC-14, AC-15, AC-16, AC-18, AC-19, AC-20 — **CONFIRMED (PRIMARY)** per the commands quoted above, at `0893989`/`bf7e049`.
- AC-4 unset-path invariance — **CONFIRMED (PRIMARY)**: tests re-run unmodified and passing; driver diff read line-by-line; end-to-end unset dispatch verified.
- AC-17 whole-tree health — **build/vet/gofmt + three packages + full designation suite CONFIRMED (PRIMARY)**; whole-suite `go test ./...` independent completion **UNVERIFIED at artifact time** (running under contention — see limitations). The implementer's PASS claim for the two long-runners remains testimony I did not contradict and partially corroborated (designation subset green in 2.4s local).
- AC-21 — not mine to verdict beyond this file's existence; the goal-done check is outstanding by design.
- FINAL VC-4 locator discipline — the implementation's re-measured locators that I spot-checked (`driver_impl.go` roleErr sites, `consensus.go` resolver, template lines) matched the code at `0893989`.

## Honest limitations / unchecked scope

- Full `go test ./...` on `0893989` was launched on local disk but had not completed when this artifact was written; `internal/app` and `internal/trajectory` were contending with the peer reviewer's concurrently running suites on this host. Evidence recorded: build/vet/gofmt green; `protocol`/`config`/`consensus` packages green; the complete designation suite plus the three R33 legacy pins green; a mutation test proving T-2's meaningfulness. If the suite completes after this artifact, that fact belongs to the next round, not retroactively here.
- I did not re-derive the §9.0 ping equivalence of `designeeAvailable` (discovered-agent `Found`) beyond reading it; its approximation is test-consistent but I ran no live-agent ping.
- The TUI (`internal/tui/protosnap.go`) and pipeline dispatch surfaces are FINAL's deferred F3/F10 items; consistent with FINAL, I did not review them beyond confirming the delta does not touch them.
- The advisory `implementation-release-plan-kimi-1.md` was skimmed, not audited (not code, not protocol-authoritative).
- Skill-repo npm-level gates (manifest/installer integrity) were not re-run; the skill delta is two markdown files, which the diff verifies.

## Open questions

1. For Phase 7 discussion: should [MAJOR-1]'s strict `confirmed <date>` shape also require the date to parse, or is the marker-plus-equality fix sufficient for this cycle?
2. For the organizer/release step only (no action for kimi-1): when `default_implementer = "codex-1"` is later set by the owner, seven existing ideas declaring `facilitator: codex-1` will ride the R20 tier-3 fall-through — behavior is correct, but the owner may want the one-line notices surfaced somewhere durable.
