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

### codex-1/F13 — CONFIRMED
what I ran: `parley classify --auto-implement --declared standard` (source-built binary, PRIMARY); my own probe (re-written, since codex's `cmd/audit-track-probe` does not exist in the tree — PRIMARY grep) calling `driver.New` over a prompt with `track: standard` + `auto_implement: true`, plus `track.PolicyFor(track.Standard, true, 3, true, false)` directly; grep for all `track.Classify` call sites (PRIMARY).
what I got: classify: `deliberation` / `declared track "standard" is under-tiered; the classifier floor is "deliberation" (auto_implement)`, exit 4. Probe: `track=standard maxReviewers=2 minReviewers=2 hardCrossCap=2 fixups=2 err=<nil>`. `track.Classify` is called ONLY from the advisory CLI (internal/app/classify.go:61); `PolicyFor` is what the driver uses (internal/driver/driver.go:133; internal/app/driver_impl.go:59).
verdict reasoning: Both halves reproduce verbatim from one shipped package. `PolicyFor` hard-rejects fast+auto_implement (internal/track/track.go:143-145) but silently applies standard's caps for standard+auto_implement (track.go:167-172). COOPERATION.md:209 makes auto_implement a deliberation trigger and :248-249 permits down-tiering below the floor only with "a recorded user OK" — a mechanism the CLI does not have — so the driver executes a combination the advisory gate rejects, with less rigor on the highest-risk (auto_implement) runs.
corrected severity: (stays MAJOR — a declared track the classifier calls unsafe is silently executed, and no intermediate layer can mask it)
is this a duplicate of another finding? None.

### codex-1/F14 — CONFIRMED
what I ran: my own re-written probe (codex's `TestAuditDefaultTrackPolicy` is not in the tree): two `driver.New` configs, identical 4 participants and `CrossReviewRounds: 9`, differing only in the presence of `track: standard` (PRIMARY).
what I got: `absent: track=standard maxReviewers=0 min=2 cross=9 hardCap=0 fixups=3` vs `explicit standard: track=standard maxReviewers=2 min=2 cross=2 hardCap=2 fixups=2` — codex's numbers exactly.
verdict reasoning: The numbers reproduce, and the protocol does say "(default `standard`)" (COOPERATION.md:200-201). The material mitigation the finding omits: this is a deliberate, commented, test-pinned design invariant — track.go:7-15 ("An ABSENT track reproduces today's behaviour byte-for-byte … only an EXPLICIT track applies §4.0's reduced ceremony" — re-read myself, PRIMARY), pinned by `TestNewAbsentTrackIsLegacy` (internal/driver/track_test.go:70-77). Effect direction is also mixed: legacy keeps ALL non-implementer reviewers (≥ standard's 2) and permits more cross rounds — more ceremony, not less rigor; the real cost is budget and bypass of the "capped at 2, then escalate" contract. Doc-vs-code divergence with the deliberate decision on the code side: the fix is a protocol-text update ("absent = legacy"), not a behavior patch.
corrected severity: MINOR — documented backward-compat posture; protocol text never updated to say so.
is this a duplicate of another finding? Shares one mechanism with codex-1/F15 (Normalize → present=false → legacy branch); F14 is the intended half of that design, F15 the unintended half — one fix does not close both.

### codex-1/F15 — PARTIAL
what I ran: my own probe: `driver.New` over a prompt with `track: standart` (PRIMARY).
what I got: `standart: track=standard maxReviewers=0 min=2 cross=9 hardCap=0 fixups=3 err=<nil>` — no error, no warning, relabeled standard, legacy posture.
verdict reasoning: The core claim reproduces: `track.Normalize` maps any unrecognized value to `(Standard, false)` (internal/track/track.go:36-48), so an explicit typo silently takes the legacy branch — the opposite of COOPERATION.md:211-216's normative "on any doubt … fail closed to the stricter track". But the finding's punchline "a typo quietly buys less enforcement than any valid track" is FALSE as stated: the legacy posture keeps uncapped reviewers (equal to deliberation, stricter than standard's 2) and an uncapped cross budget; it is looser than standard only on fix-ups (3>2) and on the escalation contract. The accurate claim is: "an unrecognized explicit track silently yields the legacy posture instead of an error."
corrected severity: MINOR — the genuine residual defect is the absent error/warning on an unrecognized explicit value; the posture itself is the documented legacy default.
is this a duplicate of another finding? Same root mechanism as codex-1/F14; distinct trigger and distinct fix (error on invalid EXPLICIT value closes only F15).

### codex-1/F16 — CONFIRMED
what I ran: my own probe: driver fixture with round-02 artifacts whose `responding-to:` key is present-but-empty plus valid `### @<other>` headings → `roundComplete(2)`; control with the key absent; cluster probe with empty `responding-to:` AND bodiless `### @x` headings (PRIMARY).
what I got: empty value: `done=true err=<nil>`; control (key absent): `done=false`; empty value + empty heading bodies: `done=true`.
verdict reasoning: Reproduced. `hasRespondingTo` (internal/driver/driver.go:466-469) returns only the presence flag from `readFrontmatterField`, which reports ok=true on a key prefix regardless of value — an empty `responding-to:` satisfies the provenance gate while COOPERATION.md:335 requires a list of prior files. The control proves the gate binds on key presence and leaks only on the empty value. Narrower than "no provenance at all" as reported: the companion `validateCrossReviewBody` still forces one `### @<other>` heading per participant — though (new, PRIMARY) that companion check is itself a bare `strings.Contains(body, "### @"+other)`, so bodiless headings pass too.
corrected severity: MINOR — the per-agent heading requirement still binds; the leak is content-only on both signals.
is this a duplicate of another finding? Same defect CLASS as codex-1/F17/F18/F22/F23 (presence-only structural gates), distinct gate in a distinct file.

### codex-1/F17 — CONFIRMED
what I ran: my own probe: round-01 artifacts with correct agent/idea/round frontmatter, the four required headings, zero body bytes, no `date:` → `roundComplete(1)` (PRIMARY).
what I got: `done=true err=<nil>`; the driver reconstructs a `round.completed` event and advances.
verdict reasoning: Reproduced. `ValidateRoundOneArtifact` (internal/runner/validation.go:42-72) checks only the three identity fields and `strings.Contains` for the four headings — no `date:`, no non-empty-section check — and `roundComplete` (internal/driver/driver.go:344-400) then promotes/drafts. The Phase-1 schema (COOPERATION.md:310-321) includes `date:` and prescribes section substance; no downstream content gate exists before consensus drafting, so header-only scaffolds genuinely auto-advance in auto-drive — which is precisely what auto-drive must not do, since round-01 IS the independent analysis.
corrected severity: (stays MAJOR — the automatic path advances empty analyses with no human in the loop)
is this a duplicate of another finding? Same leniency class as codex-1/F1 (manual draft gate) and F16/F18/F22/F23; distinct gate (round-01 validator). Also note (new, PRIMARY): the heading checks are substring-based, so the four headings may appear anywhere in the file.

### codex-1/F18 — CONFIRMED
what I ran: my own probe: `runner.ValidateReviewArtifact` on an artifact with only agent/idea/review-round frontmatter, a `## Findings` heading, and one refutation sentence — no `reviewed-commit`, `date`, `## Summary`, `## Open questions` (PRIMARY); repo-wide grep for `reviewed-commit` in validation paths (PRIMARY).
what I got: `err=<nil>`; `reviewed-commit` appears in NO validation path — only in the review-consensus template (consensus.go:513) and the manual `--reviewed-commit` flag (internal/app/app.go:572).
verdict reasoning: Reproduced. `ValidateReviewArtifact` (internal/runner/phase58.go:413-441) enforces identity frontmatter, a `## Findings` heading, and a non-empty `## Refutation attempts` section — nothing else — and the driver's review-completion gate (internal/app/driver_impl.go:286-298) relies solely on it, against the Phase-6 schema at COOPERATION.md:505-521. The validator is not a pure rubber stamp (the refutation section binds), so the finding is precisely the missing provenance/date/Summary/Open-questions enforcement — and that is accurate. Nothing proves which code revision was reviewed; a stale or pre-fix review satisfies the automatic close gate.
corrected severity: (stays MAJOR — review provenance is unverifiable on the automatic path)
is this a duplicate of another finding? None (distinct validator; same leniency class).

### codex-1/F19 — CONFIRMED
what I ran: `sed -n '1,8p' <fresh-ws>/parley-deck/COOPERATION.md` on a fresh `init` (PRIMARY).
what I got: line 3 `**Workspace:** \`<workspace-name>\`` and line 6 `**Created:** \`<date> — created by parley init\``; only `**Transport:** \`local-dir\`` substituted.
verdict reasoning: Reproduces verbatim. `cooperationForInit` (internal/protocol/workspace.go:99-108) replaces only the Transport line; the embedded template carries both placeholders (internal/protocol/defaults/COOPERATION.md:3,6), and line 6's own text ("created by parley init") shows substitution was intended. Cosmetic provenance defect.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? Shares the init run with codex-1/F12; distinct defect.

### codex-1/F20 — CONFIRMED
what I ran: fresh idea whose two `### Signoff:` blocks sit under `## Agreed decisions` while `## Signoffs` holds only a comment; then `consensus status` AND `consensus finalize --by alice` (PRIMARY).
what I got: `Consensus: ready` listing both ACCEPT signoffs, exit 0; finalize: `Finalized consensus and created .../FINAL.md`, exit 0.
verdict reasoning: Reproduced, with an extension beyond the finding: the bypass reaches finalization, not just the status printout — `parseDocument` (internal/consensus/consensus.go:338-366) scans every line for `^### Signoff:` with no section tracking, and `Finalize` (consensus.go:204) consumes that same verdict. The protocol template confines signoff blocks to `## Signoffs` (COOPERATION.md:374-382). Unlike F21, the trigger here is plausibly accidental — consensus drafts legitimately carry signoff-shaped example text, so quoted/example blocks anywhere in the body can authorize closure.
corrected severity: (stays MAJOR — an accidentally-plausible body pattern defeats the canonical consensus gate end to end)
is this a duplicate of another finding? One repair cluster with codex-1/F21 (same parser unawareness), independently observable, separable fixes.

### codex-1/F21 — CONFIRMED
what I ran: fresh idea with signoffs correctly placed under `## Signoffs` but frontmatter `idea: unrelated-idea`; then `consensus status` and `consensus finalize --by alice <queried-slug>` (PRIMARY). Tested in a separate fixture from F20 so each defect stands alone.
what I got: `Consensus: ready` with both ACCEPT signoffs, exit 0; finalize succeeded, FINAL.md created for the queried idea.
verdict reasoning: Reproduced: neither `parseDocument` nor `validateDocument` (consensus.go:331-429) ever reads frontmatter, so a consensus artifact declaring a different idea slug authorizes status=ready and finalization of the queried idea. However, the trigger requires an out-of-band file copy/misplacement — the CLI never produces this state itself — and in local-dir transport the canonical files are directly writable anyway, so the missing check is hygiene rather than a boundary (same reasoning as F4).
corrected severity: MINOR — real integrity gap in the canonical trail, but not reachable through any CLI-produced state.
is this a duplicate of another finding? One repair cluster with codex-1/F20, distinct defect.

### codex-1/F22 — CONFIRMED
what I ran: my own probe (codex's `TestAuditFinalWithoutRequiredSchema` is not in the tree): FINAL.md with `status: final`, a WRONG idea slug, only `## Final plan / specification` + two one-word lines + ~300 bytes of padding, none of the required sections → `finalScaffoldReason` (PRIMARY).
what I got: `bytes=392 purpose=false context=false acceptance=false idempotence=false risks=false references=false finalScaffoldReason=""`.
verdict reasoning: Reproduced. `finalScaffoldReason` (internal/driver/consensus.go:164-206) checks only: exists, >250 bytes, frontmatter `status=final`, ≥3 non-comment/non-placeholder lines under the one heading — never Purpose/Context/acceptance-criteria/Idempotence/risks/References/idea-slug, all required by COOPERATION.md:404-420. It is the only FINAL content gate before `setIdeaStatus("final")` (consensus.go:55-65); for design-only ideas nothing downstream re-checks, and for auto_implement the GoalCheck fires only at review close, after implementation ran off the padded spec. The function's own comment scopes it as a scaffold detector — a deliberately minimal gate — but the finding only claims the gate accepts such a file, and it does. (My F5 probe independently shows the CLI's OWN finalize template passes this gate — the two holes compose.)
corrected severity: (stays MAJOR — the automatic close path can mark a padded non-specification final)
is this a duplicate of another finding? Same leniency class as F16/F17/F18/F23; distinct gate. Composes with codex-1/F5.

### codex-1/F23 — PARTIAL
what I ran: my own probes: `ValidateImplementationArtifact` on `status: banana` + empty `## Summary of work` and nothing else; AND the downstream gate: `implReadyForReview`/`implInProgress` on "banana" vs "implemented" (PRIMARY).
what I got: `ValidateImplementationArtifact(banana)=nil`; downstream: `implReadyForReview(banana)=false implInProgress(banana)=false`, `implReadyForReview(implemented)=true`.
verdict reasoning: The validator half reproduces: internal/runner/phase58.go:391-410 requires only a matching idea slug, ANY non-empty status, and the `## Summary of work` substring — against the Phase-5 schema at COOPERATION.md:443-455. But the claimed consequence is refuted for the invalid-status half: `advanceImpl` (internal/driver/impl.go:95-107) re-reads the status via `ImplementationStatus` (internal/app/driver_impl.go:197-203) and fail-closes — "banana" is neither review-ready nor a known in-progress state, so the driver escalates and never launches review with it. A downstream gate codex did not read catches exactly the showcased case. What survives: with a VALID `status: implemented`, an empty-summary, provenance-free artifact passes both gates; only the `checks:`/build gate (impl.go:108-111) stands before review, and it is a no-op for design-only non-Go workspaces.
corrected severity: MINOR — downstream status whitelist + checks gate bound the blast radius; the residual is unchecked provenance fields under a valid status.
is this a duplicate of another finding? None (distinct validator; same leniency class).

### codex-1/F24 — CONFIRMED
what I ran: fresh `init` (scratch PARLEY_HOME), then `sed` the deck COOPERATION.md replacing the multi-agent invariant with `- AUDIT MUTATION: multi-agent execution is disabled.`, then `parley preflight --no-ping --dir <ws>` (PRIMARY); repo-wide grep for any production write path for `protocolSha256`/`packagedProtocolSha256` (PRIMARY — re-run by kimi-1: only READS exist, at internal/app/preflight.go:386-387; :619-620 is a sync-record markdown template; `init` writes only protocolRole/deckVersion/created, internal/protocol/workspace.go:83-88).
what I got: mutation live at line 63; preflight prints `Freshness: consumer — protocol matches packaged skill (in sync)`, `role=consumer deckVersion=(none) classification=in-sync`, `Ready: no pending gates.`, exit 0.
verdict reasoning: Reproduced, and stronger than reported: nothing in this CLI EVER writes either hash field into meta/version.json, so the empty-equals-empty in-sync verdict (preflight.go:419) is permanent for every CLI-created consumer deck, not transient "until some other path populates metadata". (The live deck's version.json DOES carry the hashes — so some external tooling, presumably the npm skill's sync, writes them; the CLI-under-test never does.) The §9.0 freshness gate is fail-open by construction on CLI-created decks.
corrected severity: (stays MAJOR — the deck's drift defense reports "in sync / Ready" immediately after a body-level mutation, on every CLI-created deck, forever)
is this a duplicate of another finding? Distinct from zcode-1/F5 (different field, different deck state, different mandate); see duplicates section.

### zcode-1/F1 — CONFIRMED
what I ran: `cat parley-deck/agents.toml`; `PARLEY_HOME=<scratch-with-roster-blocks> parley roster show --dir <deck-copy>`; `grep -in inherit COOPERATION.md`; `grep -n "adopt-inherited\|inherited-roster" COOPERATION.md` (PRIMARY).
what I got: agents.toml contains zero `[roster.*]` blocks and itself states "ROSTER MEMBERSHIP IS NOT DECLARED HERE — BY OWNER INSTRUCTION (2026-08-19)… the deck INHERITS ~/.parley/agents.toml at read time"; `roster show` exit 0, every row `inherited-roster`; COOPERATION.md:103-104: "The roster's authority is `parley-deck/agents.toml`… Membership … live in `[roster.<id>]` blocks there"; `inherited-roster`/`adopt-inherited`: zero hits in COOPERATION.md.
verdict reasoning: Reproduced. The deck's ratified state contradicts §2:103-106 verbatim, and §2:124-129 documents only the `legacy-roster` alternative — the inherited-roster state appears nowhere in the protocol. Two softeners zcode under-noted: inheritance-as-concept IS named at COOPERATION.md:121 ("`--scope machine` writes `~/.parley/agents.toml` and every deck inherits it"), and agents.toml's own comments (lines 66-83) self-document the state, partially mitigating the "reader concludes the deck has no roster" harm.
corrected severity: (stays MAJOR — the roster section every participant reads describes neither this deck's state nor any path to it; the softeners mitigate but do not close the gap)
is this a duplicate of another finding? Same root cause as zcode-1/F2 (undocumented inherited-roster state), distinct observable — keep both.

### zcode-1/F2 — CONFIRMED
what I ran: `parley roster render --dir <deck-copy>` (source-built 1.45.0, scratch PARLEY_HOME); inspected §2 table before/after (PRIMARY).
what I got: "roster render: this deck declares no roster of its own; the 2 rows shown come from ~/.parley/agents.toml. Writing them into §2 would commit a machine-local roster into a shared file. Declare the roster with `parley roster set <agent> --scope deck --adapter <family>`, or re-run with --adopt-inherited…" — **exit=1** (zcode reported exit 0 — wrong on the source build), no write; the §2 generated view (COOPERATION.md:133-134) still has zero data rows.
verdict reasoning: Reproduced in substance: the one command §2 and Appendix A (:1092) prescribe for generating the view declines under this deck's ratified configuration, and COOPERATION.md:105-106's "overwritten on the next render" describes a render that never happens. The roster remains fully knowable (`roster show` works; agents.toml self-documents), so the loss is the §2 human-readable table, not roster knowledge. Side observations (new, PRIMARY): the refusal message hardcodes `~/.parley/agents.toml` even when PARLEY_HOME redirects home; and zcode's exit-code report diverges from the source build (dist binary vs source, or reporting error — substance unaffected).
corrected severity: MINOR — presentation gap only; the canonical roster answer remains available through the command §2 itself names for it.
is this a duplicate of another finding? Not a duplicate of F1 (prose authority vs command refusal). Complementary pair with claude-1/F2 (same verb, opposite path: refusal vs drift-guard break).

### zcode-1/F3 — CONFIRMED
what I ran: `head -8 parley-deck/COOPERATION.md`; `cat parley-deck/meta/version.json`; `ls parley-deck/meta/` (PRIMARY).
what I got: line 7 "**Protocol synced:** 2026-08-19 — parley-deck-skill 2.9.0 / parley-deck-cli 1.45.0"; version.json `"deckVersion": "2.8.0"`, `"source": "npm:parley-deck-skill@2.8.0"`, `"updatedAt": "2026-08-12"`; meta/ holds exactly one sync record, `protocol-sync_2026-06-13_v1.3.1.md`.
verdict reasoning: Reproduced. §9.0:843-845 makes a consumer sync "update the `Protocol synced:` header line and record `meta/protocol-sync_<ISO-timestamp>.md`" — the header moved, the record and version.json did not. Precision note: the record mandate sits inside the `consumer` bullet; if this deck is source-role the mandate arguably never applied and the header was hand-edited — zcode's own second horn. Either horn is a real inconsistency. Consequence is audit-trail confusion: no operational gate consumes deckVersion-vs-header (and per codex-1/F24 the freshness gate that would consume metadata is fail-open anyway), so the practical blast radius is an auditor's contradictory signals, not a broken transition.
corrected severity: MINOR — false/contradictory sync bookkeeping; real but consequence is confusion, not malfunction.
is this a duplicate of another finding? Related to zcode-1/F5 (both version.json sync metadata) but distinct fields and mandates — keep both.

### zcode-1/F4 — PARTIAL
what I ran (in the repo copy with .git): `git branch -a | wc -l` → 72; `git branch -a | grep -c idea/`; `git log --all --oneline | grep -c "FINAL.md + close idea"`; `git log --all --oneline --merges | grep -c "FINAL.md + close idea"`; `git log --first-parent --oneline main -8 -- parley-deck/`; `git show 211d7de --stat` (PRIMARY — re-run by kimi-1: same results).
what I got: 34 closes / 9 merges match zcode exactly; kickoffs 91dbad1 and f06ccbc are direct commits on main's first-parent chain. BUT `git branch -a | grep idea/` returns **1**, not 0: `remotes/origin/idea/sync-skill-protocol-fallback`, merged via PR #66 ("Merge pull request #66 from feci/idea/sync-skill-protocol-fallback", 211d7de), which carried the full §11.B artifact set (00-prompt, consensus, FINAL, IMPLEMENTATION, review/).
verdict reasoning: The headline counts and the "recent work bypasses §11.B" thrust reproduce, and the header:5 `Transport: github-pr` vs practice gap is real. But "zero `idea/*` branches" is factually wrong and "§11.B's mechanics have never been used" is refuted: PR #66 is a genuine §11.B design-PR close on an `idea/<slug>` branch (and PRs #54/#56 used `meta/<slug>` naming — near-miss mechanics). The accurate statement: "used essentially once out of 34 closes; 25 closes and all recent kickoffs are direct commits." Deleted-branches cannot be fully excluded as an explanation, but the surviving branch + 9 merge commits bound the picture.
corrected severity: MINOR — the deck coherently operates §11.A local-dir mechanics while printing `github-pr`; the printed-rule gap is real, the "never used" framing is not.
is this a duplicate of another finding? None (zcode-1/F11 is a §11.B-internal text contradiction, different lens).

### zcode-1/F5 — CONFIRMED
what I ran: `grep -c protocolRole parley-deck/meta/version.json` → 0; `grep -rl protocolRole parley-deck/meta/` → none; `grep -rn protocolRole COOPERATION.md` → :838, :847; `grep -rn protocolRole internal/` (PRIMARY).
what I got: field absent from version.json and all of meta/; §9.0 gates behavior on it; the CLI DOES implement the switch (internal/app/preflight_test.go:156,200,280 exercises source/consumer/missing).
verdict reasoning: Reproduced exactly, no over-read: the three-way switch §9.0 documents is keyed on a field this deck's metadata never carried, and the prescribed fallback ("ask the user once and backfill") evidently never ran. The mechanism exists in code, so this is purely a deck-metadata gap — exactly what the finding claims. (My own round-01 finding kimi-1/F2 — `sync-project --yes` deletes protocolRole — supplies the root cause; noted for the fix list, not assessed here.)
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? Complements kimi-1/F2 (root cause vs symptom); distinct from zcode-1/F3 (different field, different mandate).

### zcode-1/F6 — CONFIRMED
what I ran: `ls parley-deck/ideas/ | grep -i "end-to-end\|pipeline"` → empty; `git log --all --oneline -- "*end-to-end-pipeline*"` → empty; `head -12 ideas/2026-06-02T12-07-14-meta-protocol-ch/00-prompt.md`; `sed -n '1155p' COOPERATION.md` (PRIMARY).
what I got: no directory, branch, or historical path matches the slug; the real dir is `2026-06-02T12-07-14-meta-protocol-ch` whose 00-prompt reads "Meta-protocol-change: evolve Parley Deck into a full automatic idea-to-monitoring pipeline" and whose FINAL.md contains §12's text; COOPERATION.md:1155: "This section was ratified by idea `meta-protocol-change-end-to-end-pipeline` (2026-06-02)."
verdict reasoning: Reproduced exactly. The ratification citation is dead and never resolves in history; the secondary point also holds — the actual dir violates §3:190 kebab-case and §7:745 naming (same timestamped-dir anomaly as F10c).
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? Same root cause as zcode-1/F10c (timestamped 2026-06-02 dir names vs slug rules), distinct symptom (dead citation vs CLI rejection) — one shared fix.

### zcode-1/F7 — CONFIRMED
what I ran: `grep -n "^## " parley-deck/COOPERATION.md`; greps of the Quickstart block (lines 12-41) for §10/§15; `sed -n '1355,1370p'`; the four phase-header blockquotes (PRIMARY).
what I got: §8@786, **§10@811**, §9@827, §11@877, §15@1230; the Quickstart map (lines 36-39) lists §9, §11–§14 as "the rest" — §10 and §15 appear nowhere in the Quickstart; §15.7 binds 15.1–15.5 on all three tracks including `fast`; §15 is blockquoted at the top of Phases 1, 2, 3, 6.
verdict reasoning: Both claims reproduce exactly. Charitable note zcode didn't make: §10 is a TL;DR, so its omission from a reference-appendix list is arguably intentional; the §15 omission is the substantive half and is fully supported — binding on every track, cited in four phase headers, yet absent from the declared entry point.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None.

### zcode-1/F8 — CONFIRMED
what I ran: `ls parley-deck/` → `COOPERATION.md agents.toml ideas inbox meta runs`; read COOPERATION.md:162-190 and Appendix A's bootstrap skeleton; `parley init --help` text (PRIMARY).
what I got: §3 tree (lines 164-188) lists only COOPERATION.md, ideas/, inbox/, meta/ — no agents.toml, no runs/; §12.12:1153 references `runs/`; init's help says it creates "COOPERATION.md, ideas/, inbox/, meta/, and runs/".
verdict reasoning: Reproduced. Two of six real top-level entries are missing from the canonical layout, one of them the file §2:103 calls the roster's authority. Additional instance zcode didn't cite: Appendix A's step-5 bootstrap skeleton (lines 1095-1100) also omits both — the omission is consistent across the document, not a one-off.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None.

### zcode-1/F9 — CONFIRMED
what I ran: read COOPERATION.md:585-605; spot-checked the templates at 273, 310, 360, 443 (PRIMARY).
what I got: line 588 ends "…they append a new section to `IMPLEMENTATION.md`:" → line 590 is the LE-4 closing-veto paragraph ("When `checks:` is a list… closing additionally requires…") → the promised `## Fix-up cycle N` template appears only at line 592; the LE-4 paragraph's content belongs with the closing rule at line 605; templates at 273/310/360/443 each immediately follow their introducing line.
verdict reasoning: The splice is real and reads as a mis-merged edit; the structural claim and the comparison to other §4 templates both hold. The impact characterization (driver-author confusion) is plausible but speculative — the defect itself is verbatim reproducible.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None.

### zcode-1/F10 — PARTIAL
what I ran: `parley learn automation-outer-loop --dir <ws>/parley-deck` (a); `--dir <ws>` (b); `parley learn 2026-06-02T12-07-14-meta-protocol-ch --dir <ws>` (c); `parley retro scan --dir <deck>` and `roster show --dir <deck>` as sibling controls; read internal/app/learn.go:40-70 (re-read by kimi-1, PRIMARY); `parley learn protocol-and-skill-audit --dir <ws>` as a not-complete control (PRIMARY).
what I got: (a) `learn: no idea at <ws>/parley-deck/**parley-deck**/ideas/automation-outer-loop`, exit 1 — doubling confirmed; siblings accept the deck-dir form (exit 0). (b) workspace-root form exits 0, wrote the playbook — BUT learn.go:57-62 enforces exactly the documented precondition ("Precondition: only a COMPLETED idea is distilled", requiring `IMPLEMENTATION.md` `status: complete`), automation-outer-loop's IMPLEMENTATION.md says `status: complete`, and the control run on a genuinely incomplete idea exits 1 ("idea … is not complete"). (c) exit 2, `learn: invalid idea slug "2026-06-02T12-07-14-meta-protocol-ch" (lowercase kebab-case)`; three such timestamped dirs exist.
verdict reasoning: (a) confirmed: `learn` hard-codes `<dir>/parley-deck/ideas` (learn.go:51) unlike its siblings. (b) REFUTED: the tool DOES check completion, against the protocol's own completion definition (Phase 8, COOPERATION.md:605 — `status: complete` in IMPLEMENTATION.md frontmatter); zcode measured completion by the *00-prompt.md* status (`round-01`), which is F14's staleness data point, not a learn contract violation. "The tool neither checks completion" is factually false. (c) confirmed: kebab-case validation permanently locks out the deck's three oldest completed ideas. So one of zcode's three deviations is wrong; "three of the section's four concrete claims do not hold" overstates — two of three hold.
corrected severity: (stays MINOR, but with sub-claim (b) struck)
is this a duplicate of another finding? (c) shares root cause with zcode-1/F6; (b)'s misused evidence is F14's data point.

### zcode-1/F11 — CONFIRMED
what I ran: `sed -n '977,981p;983,987p;1011,1015p' parley-deck/COOPERATION.md`; §3:177 (PRIMARY).
what I got: line 979 (Phase 5 step 3): "commits `IMPLEMENTATION.md` directly to the integration branch of the parley-deck repo (small, no PR needed)"; line 1013 (Branch protection, recommended): "Require PRs for all changes to `ideas/`."; line 985 (Phase 6) hedges the equivalent commit with "(or via a small PR if branch protection requires)".
verdict reasoning: Reproduced verbatim. IMPLEMENTATION.md lives at `ideas/<slug>/IMPLEMENTATION.md` (§3:177), so under the recommended protection the prescribed Phase-5 direct push is rejected, while Phase 6's equivalent commit carries the hedge. Internal contradiction confirmed.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None.

### zcode-1/F12 — CONFIRMED
what I ran: `head -8 parley-deck/COOPERATION.md`; `sed -n '1085,1100p'`; `head -10 internal/protocol/defaults/COOPERATION.md`; repo-wide grep for "shared channel" / "bootstrapping agent" (PRIMARY).
what I got: live header fields: Workspace, Parley deck, Transport, Created, Protocol synced, Status; Appendix A step 3 (line 1091): "Fill in the header: workspace name, shared channel path, transport, creation date, bootstrapping agent ID."; the embedded default template (`internal/protocol/defaults/COOPERATION.md:3-7`) — the header `parley init` actually writes — also has only Workspace/Parley deck/Transport/Created/Status.
verdict reasoning: Reproduced, and stronger than reported: "shared channel path" and "bootstrapping agent ID" exist in NO header, live or template — so this is not live-deck customization drift; a new adopter running `parley init` and following Appendix A hits the mismatch on the shipped template itself.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None.

### zcode-1/F13 — CONFIRMED
what I ran: `grep -n "excluded\|readiness\|preflight\|liveness"` on `ideas/{agents-verify-hermes-probe,preflight-liveness-false-negative,protocol-and-skill-audit}/00-prompt.md`; read all three frontmatters (PRIMARY).
what I got: matches only where preflight/readiness is the idea's *subject* (evidence/quotation), never a kickoff record; no `excluded:` field, no available/unavailable table, no readiness result in any frontmatter or body.
verdict reasoning: Reproduced. §9.0:831-833 ("records the result in the new idea's `00-prompt.md`") and :855-856 (the `excluded:` shape) are quoted exactly; all three ideas are standard-track with full rosters and carry no record. Scoped correctly to the three newest ideas as claimed.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? Shares only the §9.0 section number with zcode-1/F5; different mandate.

### zcode-1/F14 — CONFIRMED
what I ran: loop over `parley-deck/ideas/*/`: dirs with FINAL.md → 78; non-final/complete/abandoned `status:` in 00-prompt.md → 21; checked the ratifying-section ideas; `grep -rn "^ status:" ideas/*/00-prompt.md` for the claimed leading-space variant; `sed -n '422p;740p' COOPERATION.md` (PRIMARY).
what I got: 21/78 stale — breakdown: 18× `round-01`, 1× `open` (meta-protocol-change-fusion-execplans), 1× `kickoff` (launch-mkdir-resilience), 1× empty (launch-orphan-hardening); five ratifying-section ideas are in the stale set; verification-integrity is `final`. COOPERATION.md:422 ("Update `00-prompt.md` `status: final`") and §6 rule 5 (:740) confirmed. NO leading-space status line exists; zcode's sub-breakdown ("16 say round-01 … one round-01 with a leading space variance") sums to 20, not 21.
verdict reasoning: Headline numbers (21 of 78) and both cited rules reproduce exactly; the finding's consequence reasoning (§6 rule 5 fed false data) is sound. The sub-breakdown discrepancy (20 vs 21; nonexistent leading-space item) is an internal reporting error, immaterial to the finding — correct count is 18× round-01 + open + kickoff + empty = 21.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None (zcode-1/F10b misused one of these stale statuses as evidence, but the findings are independent).

### zcode-1/F15 — CONFIRMED
what I ran: `parley --help | grep -ci learn` → 0; `grep -ci preset` → 0; `parley preset list` → exit 0 (works); `parley learn nonexistent-idea --dir <ws>` → exit 1 (command exists); `parley run --help | grep -i preset` (PRIMARY, source-built 1.45.0).
what I got: neither `learn` nor `preset` appears anywhere in `--help`; both commands work; `run`'s `-preset` flag help reads "named roster preset to expand into participants (see parley preset list)" — cross-referencing a command the same binary never lists.
verdict reasoning: Reproduced on the source build. New adjacent instance zcode missed (PRIMARY): the `--help` Usage block also omits the entire `pipeline` command while the Commands section documents it — a third invisible-command instance.
corrected severity: (stays NIT, as reported)
is this a duplicate of another finding? None.

### claude-1/F1 — PARTIAL
what I ran: the finding's word greps against `parley-deck-skill/skills/parley-deck/SKILL.md` (sibling repo, read-only); `grep -n "parley preflight\|parley retro\|parley loop tick"` against the live COOPERATION.md; `sed -n '215,230p'` for the §4.0 table; the same word greps against the skill-bundled `references/COOPERATION.md`; structural reads of SKILL.md (PRIMARY).
what I got: SKILL.md counts reproduce exactly — `preflight: 0 / retro: 0 / loop tick: 0 / consult: 0 / repo-map: 0 / tui: 0`. Live deck :222 is verbatim `| §9.0 readiness ping | skipped | full | full |`; the deck names `parley preflight` ×1 (:832), `parley retro` ×2 (:491, :1187), `parley loop tick` ×2 (:1196, :1225). BUT SKILL.md:24-27 forbids running from the abbreviated workflow alone and mandates loading a full protocol file ("1. Prefer the live project file `parley-deck/COOPERATION.md`. 2. If the live file is unavailable, load the bundled fallback snapshot `references/COOPERATION.md`"; also SKILL.md:12), and that bundled fallback — which ships inside the skill directory — contains `parley preflight` ×1, `parley retro` ×2, `parley loop tick` ×2, `consult` ×5. Additionally, `repo-map` and `tui` are 0 even in COOPERATION.md (`parley context repo-map`, `parley tui` are CLI-only, never protocol-named) — those two grep entries pad the list.
verdict reasoning: Every literal fact reproduces; the materiality claim does not. "The skill describing obligations it gives no way to discharge" fails because the skill's own design mandates loading a protocol file (live or skill-bundled) that names every one of these commands — the finding's own quote (SKILL.md:34: "`SKILL.md` and `references/COOPERATION.md` are the vendor-neutral instructions") includes the file that discharges them. The residual real gap: SKILL.md's abbreviated workflow names 11 `parley` verbs but not the four protocol-named obligation commands — a cross-reference ergonomics issue.
corrected severity: MINOR — no undischargable obligation; a missing-commands-list nicety in the abbreviated workflow.
is this a duplicate of another finding? None.

### claude-1/F2 — CONFIRMED
what I ran: `grep -n "parley roster render" COOPERATION.md`; baseline `go test ./internal/protocol/...` in a scratch copy (PASS); scratch PARLEY_HOME with three `[roster.*]` blocks; source-built `parley roster render --dir <copy> --yes --adopt-inherited`; re-ran `go test ./internal/protocol/...` and `go test ./internal/...`; then restored the deck file and re-verified (PRIMARY). Re-read myself: internal/app/roster_render.go:68-78 and internal/protocol/drift_test.go:25-62.
what I got: COOPERATION.md:57 verbatim contains "(then regenerate the §2 view with `parley roster render`)"; render prints `Regenerated §2 in <copy>/parley-deck/COOPERATION.md`, exit 0; afterwards `--- FAIL: TestEmbeddedDefaultMatchesLiveDeck` / `drift_test.go:60: live deck: anchor "| Agent ID       | Workspace dir                       | Role          |" appears 0 times, want exactly 1 (drift guard fails closed)` — byte-identical to the claimed output; the only failing test across `./internal/...`; after restore, green again with a clean `git diff`.
verdict reasoning: The mechanism is exactly as described: `roster_render.go:73` unconditionally emits the four-column compact header `| Agent ID | Workspace dir | Role | State |`, while `drift_test.go:28` anchors the three-column padded header (held by `internal/protocol/defaults/COOPERATION.md:132`) and asserts it appears exactly once in the live deck, failing closed. Not reporter-created state: the render is the command COOPERATION.md:57 (also :129, :1092) instructs, run against the repo's own deck. Deeper root (new, PRIMARY): the drift is between two GENERATORS — `parley init`'s embedded default (padded 3-column §2 header) vs `roster render` (compact 4-column); any deck bootstrapped per :57 diverges from the embedded default's shape, and this repo's guard is merely where that becomes fatal. "Breaks the build" = breaks the test suite, not compilation; with CI gating tests, MAJOR stands.
corrected severity: (stays MAJOR — following the protocol's own bootstrap instruction deterministically breaks this repo's test suite)
is this a duplicate of another finding? Not a duplicate of zcode-1/F2 — same verb, opposite path: render REFUSING to write under inheritance without the flag (zcode) vs the success path writing a format the guard rejects (claude). Complementary pair; keep both.

### claude-1/F3 — CONFIRMED
what I ran: `grep -rn "masked-by-env" internal/ cmd/ | grep -v _test`; wider `grep -rn "masked" internal/ cmd/`; read internal/app/roster_set.go:60-110; SKILL.md's vocabulary sentence (`sed -n '280,292p'`); enumerated all `addStatus("` call sites; checked the `--explain` path (PRIMARY).
what I got: `masked-by-env` in non-test code appears exactly twice: roster_set.go:83 (comment) and :88 (advice line inside a stderr warning). SKILL.md:282-285 verbatim lists it in the closed STATUS vocabulary. All 17 `addStatus(` call sites emit 13 distinct codes — `masked-by-env` is not among them. Every OTHER documented term does have an emission path (`ok` via statusOrOK, `stale-snapshot` via rosterSnapshotState, the rest via addStatus).
verdict reasoning: Precise and correct: the term reaches the user only as advice text after `roster set`, never as a STATUS cell, so filtering rows on STATUS can never match it; the roster_set.go:83-84 comment confirms the half-fix the finding describes (the "emitter" emits prose, not the vocabulary value). No doc-vs-code inversion.
corrected severity: (stays MINOR, as reported)
is this a duplicate of another finding? None.

## Tally

42 findings assessed (codex-1 ×24, zcode-1 ×15, claude-1 ×3; own round-01 findings excluded per instructions):

- CONFIRMED: 36 (codex 21, zcode 13, claude 2)
- PARTIAL: 6 (codex F12/F15/F23; zcode F4/F10; claude F1) — in every PARTIAL the headline behavior reproduced but a load-bearing sub-claim was refuted (F23's consequence, F15's punchline, F12's doc reading, F4's "never used", F10's sub-claim (b), F1's materiality).
- REFUTED: 0 outright; UNREPRODUCIBLE: 0. Every reporter ran real commands; the refutation work landed on sub-claims, framing, and severity.
- Severity: of 29 findings filed as MAJOR, 10 survive at MAJOR (codex F2, F6, F13, F17, F18, F20, F22, F24; zcode F1; claude F2). 19 were corrected to MINOR (real but bounded consequences: advisory-only output, hygiene checks in a transport layer that has no authorization boundary, deliberate legacy posture, self-documenting state, fail-closed downstream gates the reporter did not read).

## Findings I could not assess, and why

- None fully blocked. Two honest edges:
  - codex-1/F11's terminal step (watching `parley run "continue <slug>"` actually create a stray idea) was not executed — an isolated-home run stopped at preflight before idea creation. The verdict rests on the built binary's own help text plus the unconditional create path in code (runTask → runcontrol.Create → CreateIdeaFull), which I consider decisive; tagged accordingly.
  - Whether any EXTERNAL tooling (the npm parley-deck-skill sync) writes `protocolSha256`/`packagedProtocolSha256` into meta/version.json: the live deck carries those fields, so something outside this CLI does. This does not weaken codex-1/F24 — the CLI-under-test never writes them, so CLI-created decks are permanently fail-open — but the end-to-end freshness story for skill-synced decks was not exercised (would require running the skill's sync tooling; out of scope, no external writes allowed).
- Not assessed by design (no-secrets / read-only boundaries, unchanged from round 1): live hosted-agent behavior behind any gate (all verdicts judge the gates, not agent output), and GitHub/GitLab transport mutations (§11.B/§11.C host enforcement).

## Duplicates I found across authors

- codex-1/F2 + codex-1/F3 — one root cause (review mode never subtracts the implementer from `idea.Participants`), two gates (draft-time existence at consensus.go:122/486-494; signoff quorum at consensus.go:411-415). One repair cluster, keep both as symptoms.
- codex-1/F7 + codex-1/F8 — one planner root cause (runplan never reads `track:`); F8 additionally asserts the verified-missing collapsed-FINAL driver path, so F8 is not a pure duplicate. codex-1/F11 overlaps F7's printed command but is a separate template defect in runaction.Command.
- codex-1/F14 + codex-1/F15 — one mechanism (track.Normalize unknown/absent → (Standard,false) → legacy branch); F14 is the deliberately documented half, F15 the unintended silent-misparse half. Distinct fixes.
- codex-1/F16/F17/F18/F22/F23 — one defect CLASS (presence-only structural gates vs protocol-mandated content) across five distinct gates in four files. codex-1/F1 is the manual-CLI member of the same class. Worth one hardening pass, not six separate bug fixes.
- codex-1/F5 + codex-1/F22 — compose: the CLI's own finalize template produces a FINAL.md that passes the driver's FINAL gate (verified directly). Distinct findings, one shared acceptor.
- codex-1/F20 + codex-1/F21 — one parser unawareness (parseDocument/validateDocument is section- and frontmatter-blind), two independently observable bypasses with separable fixes.
- codex-1/F24 vs zcode-1/F5 — NOT duplicates: different fields (hashes vs protocolRole), different deck states (fresh CLI-init vs live deck), different mandates. zcode-1/F5's root cause is my own round-01 finding kimi-1/F2 (`sync-project --yes` deletes protocolRole) — cross-reference for the fix list.
- zcode-1/F1 + zcode-1/F2 — one undocumented state (inherited roster), two observables (prose authority false; prescribed render refuses). Keep both.
- zcode-1/F2 + claude-1/F2 — same command (`roster render`), opposite paths: refusal without `--adopt-inherited` (zcode) vs drift-guard break with it (claude). Complementary pair — the command currently cannot succeed safely on this repo's deck at all.
- zcode-1/F3 + zcode-1/F5 — both sync-metadata gaps in version.json, distinct fields/mandates. Keep both.
- zcode-1/F6 + zcode-1/F10c — one root cause (three timestamped 2026-06-02 directory names violating §3:190 kebab-case), two symptoms (dead ratification citation; `learn` slug rejection). One shared fix (rename or alias the dirs).

## Anything round 1 missed that these findings led me to

All PRIMARY (ran/read in isolated copies against source-built 1.45.0):

1. codex-1/F5 and F22 compose (above): `parley consensus finalize`'s own output passes the auto-driver's FINAL gate — the manual and automatic close paths share one hole, not two independent ones.
2. codex-1/F7 bites with NO `cross_review_rounds` key at all (planner default 1) — the defect is the default path, not contingent on the reporter's explicit key.
3. codex-1/F9 control result: a truly EMPTY open-items section IS rejected at finalize (exit 1) — the defect is precisely presence-only checking, which bounds the fix (content check, not existence check).
4. codex-1/F16 is wider than stated: the companion `validateCrossReviewBody` (internal/driver/driver.go:483) is a bare `strings.Contains(body, "### @"+other)`, so bodiless `### @x` headings pass round-2 completion too — both of the cross-review evidence signals are presence-only. Related: `ValidateRoundOneArtifact`'s four heading checks are also substring-based (headings may appear anywhere in the file).
5. codex-1/F24 is permanent, not transient: no production code path in this CLI ever writes `protocolSha256`/`packagedProtocolSha256` (only reads at internal/app/preflight.go:386-387; the :619-620 hits are a sync-record markdown template).
6. zcode-1/F15 has a third instance: `parley --help`'s Usage block omits the entire `pipeline` command even though the Commands section documents it.
7. zcode-1/F2 details: the reported exit code was wrong (source-built 1.45.0 exits 1, not 0 — dist-vs-source divergence or reporting error; substance unaffected), and the refusal message hardcodes `~/.parley/agents.toml` even when PARLEY_HOME redirects home.
8. claude-1/F2's deeper root: two GENERATORS disagree — `parley init`'s embedded default ships the padded 3-column §2 header (internal/protocol/defaults/COOPERATION.md:132) while `roster render` emits compact 4-column (internal/app/roster_render.go:73). Fixing only the drift-test anchor would mask the generator disagreement.
9. claude-1's own "found clean" claim is conditional: the skill-surface drift guard (drift_test.go:276) SKIPS 2 of 4 surfaces and PASSes when the sibling parley-deck-skill checkout is absent — verified by running the test in a sibling-less copy (`skipping bundled skill protocol …`, `NOTE: checked 2 of 4 surfaces … they are NOT enforced by this run`, PASS). A standalone CLI clone gets none of the skill-surface enforcement claude-1 credited.
10. zcode-1/F14's sub-breakdown is internally inconsistent (sums to 20, not 21; the claimed leading-space status line does not exist) — correct breakdown: 18× `round-01` + `open` + `kickoff` + empty = 21. Immaterial, but it is the kind of detail a fix list should quote correctly.
