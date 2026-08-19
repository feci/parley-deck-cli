---
agent: kimi-1
idea: protocol-and-skill-audit
round: 2
date: 2026-08-20
responding-to: [codex-1/round-01, zcode-1/round-01, claude-1/round-01]
---

# Round 02 — kimi-1 verification of codex-1, zcode-1, claude-1 findings

(own round-01 findings excluded per instructions)

## Method

- CLI under test built from HEAD source, never the installed or dist binaries: `go build -o /tmp/pdkimi/parley ./cmd/parley` (reports `parley 1.45.0`, go 1.26.6).
- Five verification lanes, each in its own full copy of the repo under `/tmp/pdkimi/{a1,a2,b,c,d}/repo` (with `.git`), each with a scratch `PARLEY_HOME`; the real tree was never written. Reporter leftovers (`/tmp/parley-audit-codex-1.DGmaGx`) were read for layout only; all fixtures rebuilt fresh.
- codex-1's probe names (`TestAudit*`, `cmd/audit-track-probe`) do NOT exist in the tree (PRIMARY — repo-wide grep finds them only in codex-1.md); equivalent probes were re-written from scratch against the real APIs and run in a copy.
- I additionally re-ran the decisive checks myself (kimi-1, PRIMARY): consensus.go:118-124/405-415/486-494, COOPERATION.md:215-235/575-585, track.go:1-15, roster_render.go:68-78, drift_test.go:25-62, `git branch -a | grep idea/`, the `protocolSha256` write-path grep, and learn.go:40-70. Every spot check matched the lane reports below.

## Verdicts

### codex-1/F1 — CONFIRMED
what I ran: fresh idea `audit-empty` (participants alice,bob) with 1-byte `round-01/alice.md`/`bob.md`; source-built `parley consensus draft --dir <ws> --round 1 --by alice audit-empty` (PRIMARY).
what I got: `Drafted consensus at .../ideas/audit-empty/consensus.md`, `Consensus: partial`, exit 0; `00-prompt.md` flipped to `status: consensus`.
verdict reasoning: Reproduces exactly. `Draft` gates only on `missingRoundArtifacts`, a bare `os.Stat` loop (internal/consensus/consensus.go:122-124, 486-494 — re-read myself, PRIMARY). COOPERATION.md:308-322 does require frontmatter + four named sections, so the protocol reading is accurate. But no doc defines `consensus draft` as the content gate (docs/cli-reference.md:216-224; the content gate lives in the auto-driver path, internal/runner/validation.go) — this is under-enforcement, and drafting alone cannot close anything: signoffs + finalize still stand between scaffold and `status: final`.
corrected severity: MINOR — the command only scaffolds; closing still requires the full signoff chain, and no doc promises a content gate here.
is this a duplicate of another finding? Same leniency class as codex-1/F17 (driver-side round gate) but a distinct gate; not a duplicate.

### codex-1/F2 — CONFIRMED
what I ran: fresh `audit-review` with `participants: [impl, reviewer-a, reviewer-b]`, only `review/round-01/reviewer-a.md` + `reviewer-b.md` (protocol-conformant, no impl.md); `parley consensus draft --dir <ws> --review --round 1 --by reviewer-a audit-review` (PRIMARY).
what I got: `consensus draft failed: review/round-01 is incomplete; missing impl.md`, exit 1.
verdict reasoning: Exact reproduction. The draft gate applies the full `idea.Participants` list to review rounds (consensus.go:122 + 486-494 — re-read myself, PRIMARY), while COOPERATION.md:503 ("every active participant **except the implementer** writes `review/round-01/<agent-id>.md`") and :523 ("The implementer does not write a review-round file") are unambiguous. The documented subcommand is deterministically unusable in the only protocol-compliant state unless the prohibited file is fabricated. Note: the reporter's own leftover fixture contains exactly that fabricated `impl.md` (SECONDARY) — i.e. their run confirms the gate forced the violation; my fresh fixture proves the defect needs no reporter state.
corrected severity: (stays MAJOR — a documented command is hard-blocked on the compliant path)
is this a duplicate of another finding? Shares one root cause with codex-1/F3 (review mode never subtracts the implementer) but is a distinct gate (draft-time existence vs signoff quorum). One repair cluster, two findings.

### codex-1/F3 — CONFIRMED
what I ran: hand-authored protocol-schema `review/consensus.md` (no fabricated impl.md); `consensus signoff --review --agent reviewer-a --status accept`, same for reviewer-b, then `consensus status --review audit-review` (PRIMARY).
what I got: after both reviewers accepted: `Review consensus: partial`, `Missing signoffs: impl`, exit 0.
verdict reasoning: Reproduced with zero reporter state. `validateDocument` computes `Missing` from the unfiltered participant list (consensus.go:411-415 — re-read myself, PRIMARY) and the auto-driver inherits it via `ReviewStatus` (internal/app/driver_impl.go:328-356). BUT the finding over-reads the protocol's clarity: the per-track table says standard review consensus is "reviewers who reviewed sign off" (COOPERATION.md:226) while the Phase 7 template itself says "Each active participant (implementer included) APPENDS their signoff block" (:581-582 — re-read myself, PRIMARY). The override clause (:231-235) resolves the contradiction for the reporter's reading, but the code followed the Phase-7-literal reading — and that reading is executable: the driver can request impl's signoff and impl can sign, so the flow completes with extra ceremony rather than stalling.
corrected severity: MINOR — consequence is one extra demanded signoff the protocol's own Phase-7 text permits, not a stall.
is this a duplicate of another finding? Same root cause as codex-1/F2, distinct enforcement point.

### codex-1/F4 — CONFIRMED
what I ran: `audit-final` (author alice, participants alice,bob) driven to `Consensus: ready`, then `parley consensus finalize --dir <ws> --by mallory audit-final` (PRIMARY).
what I got: `Finalized consensus and created .../FINAL.md`, exit 0; FINAL.md frontmatter contains `author: mallory`.
verdict reasoning: Exact reproduction. `Finalize` (consensus.go:196-243) never inspects `opts.By`; `finalTemplate` (consensus.go:549-576) writes whatever string is passed. COOPERATION.md:398-400 assigns draftership to the prompt author / first submitter / `Drafter: yes` volunteer — mallory is none. The flag is unchecked. However, in local-dir transport the operator can already write the canonical files directly, so the missing check is hygiene, not a boundary; the substantive gate (ready/reserved consensus) still bound.
corrected severity: MINOR — forged attribution only; no authorization boundary exists at this layer to begin with.
is this a duplicate of another finding? None.

### codex-1/F5 — CONFIRMED
what I ran: the F4 finalize, then inspected the produced FINAL.md and `00-prompt.md` (PRIMARY); additionally a Go probe calling `finalScaffoldReason` on that FINAL.md (PRIMARY).
what I got: FINAL.md holds only `## Final plan / specification` with empty `### Goal/Scope/...` subsections and `## References`; prompt now `status: final`; probe returned `finalScaffoldReason=""` — the auto-driver's own FINAL gate accepts this scaffold.
verdict reasoning: Exact reproduction, and stronger than reported: the CLI's `finalTemplate` (consensus.go:549-576) writes none of the protocol-mandated sections (Purpose, Context, Observable acceptance criteria, Idempotence & recovery, Known risks — COOPERATION.md:412-418), yet its heading lines satisfy the driver's `finalScaffoldReason` content-line count (internal/driver/consensus.go:164-206). The manual close path and the automatic close gate compose into the same hole.
corrected severity: MINOR — a scaffold close requires the full signoff chain to have already passed; the missing piece is content validation, not authorization. (Its composition with F22 is what makes the driver gate matter; see F22.)
is this a duplicate of another finding? Composes with codex-1/F22 (same acceptor gate) but the finding itself is about the manual finalize/template — distinct.

### codex-1/F6 — CONFIRMED
what I ran: `find .../audit-final -maxdepth 1 -type d -name 'round-*'`; recursive grep for adversarial-alternative/correlated-agreement text across the idea; `parley status --idea audit-final` (PRIMARY); read of `Finalize` and the driver's consensus advance for any §15.6 logic (SECONDARY).
what I got: only `round-01` exists; no adversarial-alternative text anywhere; `Status: final`, `Consensus: ready`, exit 0; no §15.6 logic in either close path.
verdict reasoning: Reproduced: judgment-call prompt, unanimous round-01, single round, no §15.6 records — yet the idea closed. Caveat the reporter under-states: §15.6's applicability triggers ("no substantive disagreement", "primarily a judgment") are semantic, so a full mechanical gate is genuinely hard; but clause (b)'s consensus.md record and the FINAL one-family statement (COOPERATION.md:1354-1358) are textually checkable and are unchecked. The close condition is indeed decorative in both paths.
corrected severity: (stays MAJOR — the protocol's own defense against correlated agreement does nothing in either close path)
is this a duplicate of another finding? None.

### codex-1/F7 — CONFIRMED
what I ran: fresh fixture `audit-fast-plan` with `track: fast` + run store (`runs/audit-fast-run/events.jsonl`: run.created + round.completed); `parley continue --dir <ws> audit-fast-run` (PRIMARY). Then a stronger variant the reporter did not test: a fast idea with NO `cross_review_rounds` key at all, same run record (PRIMARY).
what I got: both variants print `Recommended: Open round-02 (cross-review) before drafting consensus` with `open-next-round` action, exit 0.
verdict reasoning: Reproduced, and the no-key variant removes the reporter-created-state angle: the planner's `readCrossReviewRounds` defaults to 1 and never reads `track:` (internal/runplan/runplan.go:112-132, 233-280), so a plain fast-track idea gets told to open the round COOPERATION.md:223 says fast skips. Two mitigations the reporter missed: `cross_review_rounds` is not a protocol key at all (no mention in COOPERATION.md), and the executing auto-driver IS track-aware (`track.PolicyFor(Fast)` → `CrossReviewRounds: 0`, internal/track/track.go:149, applied at internal/driver/driver.go:139-142) — so the bad advice is confined to read-only planner output.
corrected severity: MINOR — advisory text only; the component that actually executes honors the fast track (and per F11 the printed command wouldn't open a round anyway).
is this a duplicate of another finding? Same planner root cause as codex-1/F8 (runplan ignores `track:`); F8 adds a second verified-missing component, so not a pure duplicate.

### codex-1/F8 — CONFIRMED
what I ran: same fast-track idea with `cross_review_rounds: 0`; `parley continue --dir <ws> audit-fast-run` (PRIMARY); repo-wide grep for any collapsed-final path in `internal/` (PRIMARY).
what I got: `Recommended: Draft consensus from completed round artifacts` / `Command: parley consensus draft --round 1 audit-fast-plan`, exit 0; `collapsed` appears in `internal/` only inside protocol text strings; the driver's `advanceRound` unconditionally enters the ordinary consensus phase (internal/driver/driver.go:301-315).
verdict reasoning: Reproduced. COOPERATION.md:224 makes fast-track Phase 3–4 "collapsed: one `FINAL.md` with embedded signoffs"; neither the continuation planner nor the auto-driver contains that branch — both always route to the standalone `consensus.md` flow. The consequence is real but mild: a fast idea runs through standard ceremony (extra steps, signoffs in the wrong artifact) and still terminates correctly.
corrected severity: MINOR — wrong-route advisory plus a missing protocol feature; no corruption or stall.
is this a duplicate of another finding? Planner half shares F7's root cause; the absent collapsed-FINAL driver path is a distinct verified component.

### codex-1/F9 — CONFIRMED
what I ran: `audit-reserve` driven to `Consensus: reserved` (alice 🟡 with rollback note, bob ✅); CONTROL: `consensus finalize --by alice` with the open-items section empty; then inserted the filler `None.` under `## Open items deferred to implementation` and re-ran finalize (PRIMARY).
what I got: control: `consensus finalize failed: reserved consensus requires open items deferred to implementation before finalize`, exit 1; with `None.`: `Finalized consensus and created .../FINAL.md`, exit 0 — while the alice signoff still names the unresolved rollback design.
verdict reasoning: Reproduced, and the control bounds the defect precisely: `hasSectionContent` (consensus.go:215, 626-644) requires only some non-comment text, so `None.` satisfies the gate COOPERATION.md:390 attaches to 🟡 ("if the reservation is logged as 'open items deferred to implementation'"). The reservation itself is NOT lost — it survives in the append-only signoff block of the canonical consensus.md that FINAL.md references; the unenforced part is the deferred-items log's correctness, a semantic property tooling can only approximate.
corrected severity: MINOR — presence-only gate; the reservation persists in the signoff record, and the failure mode requires a drafter writing filler that contradicts an existing reservation.
is this a duplicate of another finding? None.

### codex-1/F10 — CONFIRMED
what I ran: past F2's gate, `parley consensus draft --review --round 1 --by reviewer-a audit-review`; inspected the generated file; Go probe calling `runner.ValidateReviewConsensusArtifact` on it (PRIMARY).
what I got: generated frontmatter is `idea / cycle: 1 / drafted-by / date / reviewed-commit:` — no `review-cycle`, no `outstanding_agreed_fixes`, no `blocked:`; probe: `ValidateReviewConsensusArtifact error=... frontmatter missing outstanding_agreed_fixes (the auto fix-up loop fails closed without it)`.
verdict reasoning: Both halves reproduce. The CLI's review template (consensus.go:508-527) writes `cycle:` where the Phase 7 schema requires `review-cycle: N` (COOPERATION.md:562) and omits `outstanding_agreed_fixes`, which the automated Phase 7/8 gate requires (internal/runner/phase58.go:379-388; internal/app/driver_impl.go:337-341) — the manual command's output genuinely cannot feed the tool's own auto path. Two refinements: the `cycle:` vs `review-cycle:` half is cosmetic (nothing in `internal/` reads either key — PRIMARY grep), and the empty `reviewed-commit:` is just the unpassed `--reviewed-commit` flag, not a defect.
corrected severity: MINOR — the auto gate fails closed with a clear actionable error and the repair is one frontmatter line; no silent corruption.
is this a duplicate of another finding? None.

### codex-1/F11 — CONFIRMED
what I ran: `parley continue --dir <ws> audit-fast-run` (prints the command); the built binary's own `parley help` run entry; code read of the run path: `runTask` (internal/app/app.go:1899) → `runcontrol.Create` (internal/runcontrol/runcontrol.go:50) → `protocol.CreateIdeaFull` (internal/protocol/workspace.go:125-131) (PRIMARY output+help, SECONDARY code).
what I got: continue prints `Command: parley run --auto --dir . "continue audit-fast-plan"` for the action labeled `Open round-02 (cross-review)`; help says `run` = `Create a new idea from TASK and start round-01 with selected agents.`; the create path is unconditional with no "continue <slug>" special-casing anywhere.
verdict reasoning: Confirmed. COOPERATION.md:328-336 defines opening round N+1 as creating `round-0(N+1)/<agent-id>.md` files — nothing the printed command does. The code comment at internal/runaction/action.go:50-53 shows the authors know it is a stopgap, but even the stopgap intent is wrong: `run` cannot advance an existing idea. I did not watch a stray idea actually appear (an isolated-home execution stopped at preflight before idea creation), so the terminal step rests on help text plus the unconditional create path — decisive but noted honestly.
corrected severity: MINOR — following the advice visibly creates a stray idea (`Created idea X and run Y`); confusing and audit-trail-noisy, but recoverable and non-destructive.
is this a duplicate of another finding? Overlaps codex-1/F7's printed command (an instance of this defect) but is a distinct component: runaction.Command's template vs runplan.Plan's recommendation.

### codex-1/F12 — PARTIAL
what I ran: `env PARLEY_HOME=<scratch> parley init --dir <fresh-ws>`; `find <ws>/parley-deck -maxdepth 2 -type f`; read the seeded central `<scratch-home>/agents.toml` (PRIMARY).
what I got: `Initialized Parley Deck workspace ...` + `Central agent defaults: <home>/agents.toml (override per-project in .../parley-deck/agents.toml)`; deck contains only COOPERATION.md and meta/version.json; deck agents.toml missing; the seeded central agents.toml DOES carry per-agent `model`/`reasoning` keys (mostly `cli-default`).
verdict reasoning: The behavioral facts reproduce exactly, but the finding over-reads its cited line: COOPERATION.md:43-57 puts the MUST on "the **facilitator**", names `parley roster set` (a separate command) as the recording mechanism, and points at the skill for the interactive flow; the one duty the doc assigns to `parley init` itself — creating the central config — is performed. Doc-vs-code inversion: the code matches the doc's division of labor; the residual gap (a CLI-only user gets no enforcement of the bootstrap confirmation, and "Initialized" reads as completion) is a UX/docs nit.
corrected severity: MINOR — the mandated confirmation was never specified as an init-binary behavior.
is this a duplicate of another finding? Shares the init run with codex-1/F19 but claims a different defect; not a duplicate.

