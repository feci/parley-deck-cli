---
agent: zcode-1
idea: meta-protocol-change-driver-unstall
date: 2026-10-09
product-commit: 134ac40cd178ebe0318a838d9357c4e6f561f925
skill-commit: e46e551871005ecac4225f7cb55a652754fbbb9b
role: independent fresh goal checker (separate process from every review and signoff)
track: deliberation
transport: github-pr
---

# Goal check — zcode-1, fresh process, 2026-10-09 (~12:08Z)

## Role and scope

I am zcode-1, launched as the independent goal checker in a **new process separate from
every review and signoff artifact**. This file is my only canonical write. I signed no
consensus, modified no production code, no prior artifact, and no reviewed plan. Temporary
test executions and one throwaway probe test ran under
`.parley-runtime/driver-unstall/goal-probes/` (the probe file inside `internal/budget/`
was deleted immediately after its run; `git status` restored to exactly the pre-probe
state). No credentials, no worktree declarations or pruning, no real legacy apply, no
terminal allocation, no browser automation. Executions completed in ~10 minutes, inside
the 1800 s hard ceiling. Native usage telemetry was not available to me; I invent no
token/cost figures.

**Verdict scope.** A PASS here is explicitly scoped to **functional source completion** —
AC1–AC8 and AC9's checks/drift portion at product commit 134ac40. AC9's live
channel/install verification remains **PENDING** exactly as signed and reviewed: a
separate fresh Zcode channel process must PASS before any released handoff. No
publication has occurred yet and no local D6 declaration was applied (both verified
below). This check can only withhold a close, never establish one; the final zero-fix
review consensus, both final own ACCEPTs and the brief's attended-close conditions follow
this artifact and are not pre-supplied by it.

## Protocol attestation

Launch supplied the full live protocol verbatim: `context_mode: full`,
`source_sha256 = packet_sha256 = 0357d504982f92b713f2f86604129f46e1ef3276664e4da91d39ecf420253f0c`.
The shadow packet audit in the launch header is an unapplied diagnostic (`would_fallback_reason:
unknown-phase:-1`, packet_bytes 0) and does not describe the supplied full protocol, which I
read. Independently re-verified (PRIMARY): `shasum -a 256` of `parley-deck/COOPERATION.md` =
`0357d504…53f0c`, equal to `meta/version.json` `protocolSha256` and `packagedProtocolSha256`
(`protocolRole: source`) and to the skill copy
`worktrees/driver-unstall-skill/skills/parley-deck/references/COOPERATION.md` (same hash).
Transport `github-pr`; track `deliberation` (protocol-text change; FINAL D2 synchronizes
normative wording). Staged core `~/.parley/staging/COOPERATION-2.18.0.md` exists and differs
from `2.17.0.md` by exactly the two normative hunks (Phase 8 "bounded goal-done check" ceiling
paragraph; §9.0 "goal checks use Phase 8's track/configuration ceiling") — diff read this
session. No release is published from it yet.

## Tree-state verification (verified myself, as instructed)

- CLI worktree HEAD **is exactly** `134ac40cd178ebe0318a838d9357c4e6f561f925`
  (`git rev-parse HEAD`); **zero commits after it** (`git log 134ac40..HEAD` empty).
- Everything after the reviewed product tree is **uncommitted and confined to this idea's
  records**: `IMPLEMENTATION.md` (modified) plus untracked `review/round-02/`,
  `review/round-03/`, seven `source-context/` evidence files, one `inbox/` note to zcode-1,
  and three `parley-deck/runs/` directories (`git status --porcelain`). No `.go`, docs,
  skill, protocol or any non-`parley-deck/` path changed (`git diff 134ac40 --stat --
  ':(exclude)parley-deck'` empty; committed non-deck changes after 134ac40: none).
- Skill worktree clean at `e46e551871005ecac4225f7cb55a652754fbbb9b` (`git status` empty).
- The May history bytes are intact: `parley-deck/runs/20260510T194003Z/events.jsonl` hashes
  `ee4f52b717159e54aa8f144161d14467160cfc6846aed97abbdd6768683c6906`, byte-equal to
  `git show 3ec10ac:parley-deck/runs/20260510T194003Z/events.jsonl`.
- **No local D6 declaration applied**: the git-common budget area
  (`$(git rev-parse --git-common-dir)/parley-launch-budgets/`) contains only the pre-existing
  `cycles-23451469…` entry and **no `legacy-*`** directory anywhere under the common dir.
  D6 remains locally activation-pending until the owner attends the one-time apply.
- Mechanical reviewer gate: `source-context/reviewer-gate-result.json` records
  `membership.SingleReviewerAfterDropout` → `{"allowed":false,"error":""}`. I make **no
  automatic reviewer-relaxation claim**; only the brief's expressly pre-authorized attended
  close path (final review free of open CRITICAL/MAJOR; both final own ACCEPTs; this fresh
  goal PASS; independent current-tree AC evidence) can apply, and its other constituents
  are still outstanding after this file.

## Per-criterion evidence (§15 provenance: PRIMARY = check I executed in this fresh process, with command/result; otherwise exact locator)

### AC1 — refusal without authority; two ideas bind after one attended declaration; bytes retained — MET

- PRIMARY (self-written probe with the **real May bytes** from this repository, not the
  suite's synthetic fixture): temp test `TestGoalProbeRealMayBytesUnstallTwoIdeas` copied
  `parley-deck/runs/20260510T194003Z/events.jsonl` into a two-worktree git fixture, then:
  `EnsureCycleBinding(…, "idea-one", CrossReview, 3, 2, "", "")` before any declaration →
  error containing `lacks idea identity` (the reproduced D6 gate); one library-level
  `ApplyLegacyDeclaration` (inspect first disclosed `History: "unknown-history"`) → both
  `"idea-one"` and `"idea-two"` then bind cleanly with no per-idea declaration; events
  bytes byte-equal before/after. PASS in 0.62 s
  (log: `.parley-runtime/driver-unstall/goal-probes/real-may-probe.log`). Probe file deleted
  after the run.
- PRIMARY: `go test ./internal/budget -run 'TestLegacyDeclaration|TestBudgetLegacy' -count=1`
  → `ok 8.370s` at HEAD (= 134ac40), covering
  `TestLegacyDeclarationUnstallsSuccessiveIdeasWithoutResettingCaps` (incl. cap preservation
  `Reserve` → `ErrLimit` and no-driver-authored-store assertions) and the read-only-preview /
  attendance test (log: `goal-probes/d1-focused.log`).
- The actual real-repository refusal remains on record: `source-context/continue-round02.log`
  line 1: `continue --auto: run round-02: cross-review accounting: historical cycle event
  lacks idea identity` (read this session).

### AC2 — unknown-history disclosure, frozen provenance, replay idempotence, conflict refusal — MET

- PRIMARY (same probe): `InspectLegacyDeclaration` returned `History: "unknown-history"` with
  a manifest digest; identical replay of the same decision ID succeeded; changed reason under
  the same ID → refusal (`conflicting legacy decision ID reuse` path exercised).
- PRIMARY: the focused suite run above included the frozen-adoption, `legacy.json`,
  `Policy.LegacyHistorySHA256`, cycle-inspect provenance and `recorded_at`-tamper assertions
  (`internal/budget/legacy_cycle.go`, `legacy_declaration.go:319-336` — re-read at HEAD).

### AC3 — hostile mutations/copies/aliases/identity/charges refuse; cross-worktree revalidation — MET

- PRIMARY: the same focused run executed
  `TestLegacyDeclarationHostileChangesRefuseNewAndCachedBindings` green at HEAD — 17 modes
  (mutated, missing-copy, all-deleted, cursor, run.json manifest, charge, recovered identity,
  extra-unknown, missing/corrupt record, hidden file, empty directory, symlink, hardlink,
  missing adoption, changed policy reference, unavailable root), each refusing both cached
  binding and new-idea binding except the two correctly policy-local modes.
- Alias/bounds guards re-read at HEAD: `run_identity_inventory.go` (symlink/non-regular
  refusal), per-platform `legacySingleLink` (Nlink/NumberOfLinks hard-link block),
  `legacyScope` parent-chain alias walk (`legacy_declaration.go:58-71`), post-read
  re-digest (`legacy_declaration.go:222-229`).

### AC4 — charges/caps, per-idea migration, unavailable roots unchanged; no driver-created authority — MET

- PRIMARY: covered by the same green focused run (caps not reset; `InspectProtocolMigration`
  still refuses for a fresh idea; unavailable-root mode refuses). No `legacy-*` authority is
  created by inspect or the gate — confirmed live on this real worktree: after my read-only
  probes the git-common budget area still contains no `legacy-*` entry (re-checked after all
  executions).

### AC5 — exact derived bound, default/invalid track, config clamping, beyond-120 s witness — MET

- PRIMARY: `go test ./internal/app -run 'TestGoalTimeoutTrackAndConfiguredBounds|…' -count=1 -v`
  → `TestGoalTimeoutTrackAndConfiguredBounds PASS (0.01s)` at HEAD (table: absent track→15 m,
  fast→5 m, deliberation→30 m, deliberation+240000 ms→4 m, fast+1800000 ms→5 m, −1→15 m,
  max-int→15 m, malformed/empty tracks→`invalid goal-check track` refusal).
- PRIMARY (the offered real witness, in my own process):
  `PARLEY_TEST_LONG_GOAL_CHECK=1 go test ./internal/app -run '^TestGoalUnstallBeyondFormerCeiling$' -count=1 -v`
  → `--- PASS: TestGoalUnstallBeyondFormerCeiling (121.16s)` — a real child sleeping 121 s and
  printing PASS completes under the derived deliberation-track ceiling clamped to 180 s, i.e.
  beyond the former 120 s product ceiling (log: `goal-probes/long-witness.log`).

### AC6 — original plus one attempt at one frozen ceiling across restarts/runs; FAIL finality — MET

- PRIMARY: `TestGoalUnstallRetriesAndReplayRemainBounded` PASS (15.75 s) — valid-fail 1
  attempt, failed-pass 1, malformed 2, child-failure 2, retry-success 2; after
  `RunID="resumed"` + timeout 240000→1000 the replay launches no new child and the frozen
  `participant_timeout_ns` bound is retained.
- PRIMARY: `go test ./internal/runner -run 'TestParticipantGoalCeilingFrozenAcrossInterruptedRetry|TestReviewGateSupervised' -count=1`
  → `ok 8.441s` (interrupted first attempt resumes at the original 240 s ceiling despite the
  raised option; watchdog classes feed the same two-attempt ledger; replay mints no attempt).

### AC7 — failed-process PASS cannot close; hard timeout bounded; cancellation/tamper/control-plane fail closed; buffered hard bounds — MET

- PRIMARY: `TestGoalUnstallBufferedHardDeadlineAndCancellation` PASS (7.55 s) at HEAD with the
  reworked telemetry oracle — exactly two timeout-class terminal records with distinct
  ordinals and the same frozen 1 s ceiling, no third child on replay, cancellation ≤1 start,
  failed close asserted before and after replay. The failed-pass and control-plane refusal
  cases are asserted in `TestGoalUnstallRetriesAndReplayRemainBounded` (same green run).
- No membership/policy widening: production bytes unchanged since the round-1-reviewed
  commit — `git diff c659bc8..134ac40 --name-only -- '*.go'` lists only
  `internal/app/driver_unstall_test.go` (re-verified this session via the tree-state check:
  HEAD = 134ac40, nothing after it).

### AC8 — unchanged LE-5/LE-7/LE-11 gates; close-time process duties — MET as to code; process items correctly still open

- PRIMARY (code, by byte-identity + read): the goal-check close gates, loop budgets and
  close-decision wording are untouched by this idea (only the test file changed in fix-up 1;
  nothing at all changed after 134ac40). Dissent/reservation handling unchanged.
- Process duties — **not yet claimable, and not claimed**: the currently signed
  `review/consensus.md` is cycle 1 with **3 agreed fixes** and cannot be read as a zero-fix
  close; the final zero-fix consensus and both final own ACCEPTs follow this goal check; the
  mechanical single-reviewer exception returned `allowed:false`. My round-02/round-03 reviews
  record no open CRITICAL/MAJOR. This goal check supplies one of the brief's four
  attended-close conditions; it does not supply the others and does not close the idea.

### AC9 — checks/drift pass (this PASS's scope); live channels/install PENDING — MET as scoped; live portion PENDING

- PRIMARY: `go test ./internal/protocol/... -count=1` → `ok 0.214s` (drift/packet tests);
  `go build ./...` OK; `go vet ./internal/budget ./internal/app ./internal/runner` OK;
  `TestVersionFileMatchesBinaryVersion` PASS (version lockstep 1.54.0).
- Protocol-copy identity: deck = skill reference = `meta/version.json` = staged core hunk set
  (all under AC attestation above; staged 2.18.0 remains unpublished).
- Full local suite at 134ac40: independently re-read by me this session —
  `source-context/go-full-fixup01.log`: 34 `ok` lines, **0 FAIL lines**
  (`internal/app 762.022s`, `internal/budget 102.868s`, `internal/runner 162.862s`), started
  after the 134ac40 commit per round-2's verified timing; tree unchanged since. Skill suite:
  `skill-tests-complete.log` read — all skills ok incl. `parley-deck: ok (9 files,
  sha256:26e74a87…)`; skill worktree clean at e46e551.
- CI at 134ac40 (live, `gh pr checks 81`, PR81 head = 134ac40, OPEN):
  **macos-latest ×2 PASS, ubuntu-latest ×2 PASS, windows-latest ×2 FAIL** (jobs 113792656427,
  113792672173) — matching G1-R3's scoped gate: current-tree + macOS + Linux green; Windows
  unresolved/experimental, CLI WinGet held, no all-CI-green claim anywhere.
- **Live channel/install evidence: PENDING** — no release, Homebrew, WinGet or skill-install
  publication has occurred; the separately binding fresh post-publication Zcode channel
  process must PASS before the released handoff. The exact D6 activation request/command and
  the still-unpublished npm/core commands in version order belong to that handoff.

## G1-R3 Windows disposition — verified honestly

- `source-context/windows-final-comparison.json` records both terminal 134ac40 Windows jobs:
  push job 113792656427 `completed/failure` with 356 direct failing names, PR job 113792672173
  `completed/failure` with 355; both include `panic: test timed out after 45m0s` and
  `internal/evidence [build failed]`.
- PRIMARY (re-derived from **raw retained bytes**, not the JSON): extracted direct distinct
  `--- FAIL:` names myself from
  `release-delivery/2026-10-09-driver-unstall/logs/{baseline-windows-ci,windows-fixup-push,windows-fixup-pr}.log`
  → baseline(main 128e30b) 356, push 356, pr 355; **added-name sets vs baseline: empty for
  both** — the tripwire is clean; no new in-scope Windows regression is demonstrated. Commit
  identity anchors confirmed inside each log (fetch of `128e30b…` as main; fetch of
  `134ac40…` as driver-unstall; PR merge `665121e Merge 134ac40… into a9e383d…`).
- Honest limits, retained: name-match does not prove cause equality, does not certify tests
  unreached before the 45-minute timeout, and native Windows remains unresolved/experimental.
- **E1 correction verified (resolves the round-3 NIT; no earlier review rewritten):**
  corrected `windows-ci-comparison.json` now hashes the retained raw bytes and every value
  reproduces — `baseline-windows-ci.log` → `da7c0e4faececbeb70ca3ea0ee95a7dfe01b84b6548d95894c35c6b00648e772`,
  `candidate-windows-ci.log` → `fdefb4d2…` (both re-hashed by me this session, `shasum -a 256`),
  plus the 134ac40 logs `d2fdaf2a…`/`ce898a3c…` matching `windows-final-comparison.json`; it
  stores direct (356/353) and inclusive (357/353) counts, retained log paths and an explicit
  extraction pattern; the old mis-described value `ed35474a…` is now correctly carried as
  `normalized_lf_text_sha256`; the original JSON is preserved byte-identically as
  `windows-ci-comparison-initial.json` (still carrying the old field for audit).

## Windows coverage-limit assessment (launch question) — RA3's phrase is not acceptance evidence

Round-03 RA3's "those packages ran every test" rested on non-verbose logs. I verified the
limitation from source and empirically:

- `internal/app/driver_unstall_test.go:198-201`: `TestGoalUnstallBeyondFormerCeiling` calls
  `t.Skip` unless `PARLEY_TEST_LONG_GOAL_CHECK=1`; CI does not set it (no reference in
  `.github/` or scripts).
- `internal/budget/legacy_declaration_test.go:137-148`: the symlink and hardlink hostile
  modes `t.Skipf` when the operation is unavailable.
- Empirical (PRIMARY): non-verbose `go test` prints **no** skip line for a skipped top-level
  test (local run without the env var: only `ok`, exit 0). Therefore the zero `--- SKIP`
  lines in the Windows logs are fully consistent with the witness having been skipped there
  and prove nothing about per-test execution. Non-verbose logs establish only that the
  packages' binaries ran to their printed summaries with those failures.

Conclusion: no per-test Windows execution claim enters my acceptance rationale; the honest
bound is "package-level completion summaries with no failure lines from this idea's app/budget
fixtures; runner D2 fixtures and the reworked oracle have no completed Windows witness;
coverage unresolved" — exactly what IMPLEMENTATION.md's "Native Windows coverage limit" note
records. All actual Linux/macOS witnesses (including my fresh 121.16 s run) retain their scope.

## Verdict

All functional criteria AC1–AC8 are met with fresh current-tree evidence at product commit
134ac40, and AC9's checks/drift portion passes (local suite, skill suite, drift/packet tests,
build/vet, macOS ×2 + Linux ×2 CI green, protocol copies identical, staged core 2.18.0
correct and unpublished). Windows is recorded unresolved/experimental with a clean
added-name tripwire and no all-green claim. Nothing I attempted broke a criterion, and no
functional criterion is unmet or incomplete within this check's scope. AC9 live
channel/install verification, the final zero-fix review consensus, both final own ACCEPTs,
the owner's one-time D6 activation and all publication acts remain PENDING after this file,
as signed and reviewed.

GOAL-CHECK: PASS
