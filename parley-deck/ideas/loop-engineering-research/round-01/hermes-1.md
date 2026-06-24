---
agent: hermes-1
idea: loop-engineering-research
round: 1
date: 2026-06-22
lens: maker/checker & verification
---

# hermes-1 — maker/checker & verification lens

My lens is the maker/checker split (LE block #5) and the "Decide/Verify" rows of
the lifecycle table: does Parley's Phase 5 (maker) vs Phase 6 (checker) actually
*separate generation from verification*, does the verifier *default to
refutation*, is it a *different model*, and what is the *self-grading* risk that
produces a false-green. I trace each of these to exact lines in both
`COOPERATION.md` and the `parley-deck-cli` engine, then end with a prioritized,
concrete recommendation list. I do not address outer-loop/cron discovery or cost
budgets except where they intersect verification — those belong to codex-1 and
antigravity-1.

## 1. What Parley already has through this lens (and the seam)

- **Maker/checker separation by agent, not just phase.** Phase 5 runs one
  implementer (`runner.RunImplementation`, `internal/runner/phase58.go:23`);
  Phase 6 runs every *non-implementer* (`OpenReviewRound`,
  `internal/app/driver_impl.go:132`). The implementer is barred from writing a
  review file (`COOPERATION.md` §Phase 6, line 374). **Seam:** the separation is
  by *agent id*, not by *model* — see §3.
- **Checker isolation is real.** Each reviewer gets a disposable git checkout
  (HEAD or a temp-index snapshot commit) on local tmp, not the live tree
  (`internal/runner/reviewsnapshot.go:52-155`). This is a genuine LE-style
  "isolated verifier" — better than most loop-engineering stacks, which review
  in-place. Keep.
- **The known false-green was already half-fixed (AF3).** The review-consensus
  drafter — the agent that writes the `outstanding_agreed_fixes` integer the
  loop terminates on — MUST be a non-implementer, "so the implementer cannot
  filter reviewer findings out of the consensus (AF3)"
  (`driver_impl.go:43-49`). This is Parley's own prior Phase-6 false-green
  lesson, encoded in code. **The other half is still open — see §2.**
- **Stopping conditions exist** but are iteration/duration only:
  `MaxRounds=4` (`driver.go:73-74`), `MaxFixupCycles=3` (`driver.go:79-80`),
  `roundDeadline=30m` per tick (`internal/driver/loop.go:16`). No max-cost, no
  goal-condition checker — see §4.
- **`strict_gate` exists as protocol prose** (`COOPERATION.md` §Phase 8 lines
  454-477) and is the closest thing to a refutation re-verification pass. **The
  driver ignores it entirely** — see §4.

## 2. The self-grading risk, traced to where the false-green lives

The auto-drive loop terminates on exactly one signal:
`rs.OutstandingAgreedFixes == 0` (`internal/driver/impl.go:186`), plus signoffs
present and triage in `{Ready, Reserved}` (`impl.go:182`). That integer is
authored by **one agent** — `drafter = reviewers[0]`
(`driver_impl.go:46-49`) — via `BuildReviewConsensusPrompt`
(`phase58.go:332-361`), which instructs it to "Set outstanding_agreed_fixes to 0
only when nothing remains to fix." The driver then calls `Complete()` — a
deterministic file write by the orchestrator, "NOT an implementer agent"
(`driver_impl.go:227-229`) — so the implementer can't self-close. Good.

But the *closure judgment* is still a single agent's self-report that "nothing
remains." There is **no objective re-verification** between "drafter wrote 0" and
"loop done":

- `ValidateFixupArtifact` (`phase58.go:162-184`) checks the IMPLEMENTATION.md
  *shape* (frontmatter, a `## Fix-up cycle` section) — not that the agreed fixes
  were actually applied, nor that anything works. The "artifact-wins" override
  (`phase58.go:128-139`, consensus D7) means a fix-up that exits non-zero but
  wrote a valid-shaped file is *success*.
- `RunChecks` (`driver_impl.go:118-130`) runs `go test ./...` **only if `go.mod`
  exists**; otherwise it returns `(true, "no go.mod … no checks to run")`. For
  any non-Go repo, and for every design-only idea, the post-fix-up gate
  (`impl.go:216`, "checks failed after fix-up … before review") is a **no-op
  that always passes.** The guard that's supposed to catch a broken fix-up
  before spending reviewers catches nothing outside Go.
- FINAL.md's "Observable acceptance criteria" (`COOPERATION.md` §Phase 4, line
  280: "behavior a reviewer or the driver can check") are **never evaluated by
  the driver**. The reviewer prompt says reviewers "may cite a criterion"
  (`COOPERATION.md` line 374) — advisory, not a gate.

So AF3 fixed "implementer grades its own consensus." It did **not** fix
"drafter grades the loop closed with no independent check that the fixes hold."
That is the residual false-green: under `auto_implement`, one reviewer agent
says "0 fixes left," everyone signs ✅/🟡, and the loop completes — even if no
agent actually re-ran the code or re-verified a single acceptance criterion. LE
block #5 names this exactly: *"The model that wrote the code is way too nice
grading its own homework."* Moving the grader from implementer to one reviewer
reduces but does not eliminate the optimism bias — especially when reviewer and
drafter are the same single agent (the `len(reviewers)==1` case).

## 3. Refutation mode: present in spirit, absent in the prompt

LE says the verifier **defaults to refutation** and the lifecycle table marks
the Verify row "defaults to refutation." Parley's Phase 6 reviewer prompt
(`BuildReviewPrompt`, `phase58.go:235-267`) is **neutral/confirmatory**:

> "Review the implementation against FINAL.md and IMPLEMENTATION.md … Each
> finding states what is wrong, why it matters, and the concrete fix."

There is no instruction to *assume the implementation is wrong*, to *try to
construct a failing case*, to *run* the acceptance criteria, or to *only* report
"no findings" after actively attempting and failing to break it. The artifact
validator (`ValidateReviewArtifact`, `phase58.go:397-420`) requires a
`## Findings` section — so an empty-findings review is valid with zero effort
expended. That is a confirmatory posture, not an adversarial one. This is the
single highest-leverage gap in my lens: the protocol and the prompt both let a
reviewer green-light by *not noticing* rather than by *failing to refute*.

The `strict_gate` close rule (`COOPERATION.md` §Phase 8) is the protocol's
refutation pass — "a fresh full-scope Phase 6 review round … produces no
findings of any severity." But (a) it is opt-in per idea, not default, and (b)
the driver does not implement it (next section).

## 4. Model diversity: true in practice, not enforced

In this deck the roster is `[claude-1, codex-1, hermes-1, antigravity-1]`, and
`agents.Discovery.Model` is hardcoded per agent (`internal/agents/discover.go`:
claude→`claude-opus-4-8[1m]`, agy→`Gemini 3.5 Flash (High)`,
hermes→`xai/grok-4.3`). So maker and checkers are different models *by
accident of the roster*. Nothing enforces it: a two-agent deck
`[claude-1, claude-2]` would run same-model review and auto-complete. LE block #5
is explicit that the verifier should be "ideally different models"; the prompt
for this idea names "requiring the verifier to use a different model than the
implementer" as in-scope. The `Model` field is already discovered and logged
(`runcontrol.go:161`) — the comparison is a few lines that don't exist yet.

## 5. Conditional rigor: strict_gate is prose; RunChecks is Go-only

Two places where the system's rigor should scale with risk but doesn't:

- **`strict_gate` is advisory-only.** The idea that added it
  (`meta-protocol-change-review-gate-honesty/FINAL.md:48`) states verbatim: "No
  driver enforcement was added (per D5); `strict_gate` is advisory to humans in
  this change." The promised follow-up — `ReadStrictGate` + machine-readable
  close fields `strict_gate_clean` + `closing_review_round` in review consensus
  (`FINAL.md:58`, `consensus.md:59-60`) — was never built. `advanceReview`
  (`impl.go:186-197`) closes on `OutstandingAgreedFixes == 0` whether or not
  `strict_gate: true` is set. So the protocol's one refutation re-verification
  pass is inert under auto-drive — the exact mode where false-green is most
  dangerous.
- **`RunChecks` is `go test ./...` or pass.** A non-Go target repo, or a
  design-only idea with no `go.mod`, gets a free pass through both the
  pre-review gate (`impl.go:94`) and the post-fix-up gate (`impl.go:216`). The
  AF1 guard ("a fix-up that broke the build escalates immediately") is a real
  guard only for Go modules.

## 6. `TriageReserved` closes the auto-loop silently

`ACCEPT-WITH-RESERVATIONS` from all signers, with no BLOCK and none missing,
yields `TriageReserved` (`internal/consensus/consensus.go:423-424`), which
`advanceReview` treats as closeable (`impl.go:182`). A reservation is a soft
"something is off" — under `auto_implement`, completing on it without surfacing
the reservation notes to a human re-creates a quieter false-green: the loop
"passed" while a checker flagged doubt. Under HITL this is fine (the human reads
the reservation); under auto-drive it should not auto-close.

---

## Prioritized recommendations

| # | Recommendation | Protocol / CLI (exact target) | LE principle | Effort | Risk | Call |
|---|---|---|---|---|---|---|
| 1 | Make refutation-default explicit in the Phase 6 reviewer posture and prompt | Protocol: `COOPERATION.md` §Phase 6 (lines 353-374) + CLI: `BuildReviewPrompt` (`internal/runner/phase58.go:235-267`) and `ValidateReviewArtifact` (`phase58.go:397-420`) | Verify row: "defaults to refutation" | S | Low | **Adopt** |
| 2 | Implement `strict_gate` enforcement — `ReadStrictGate` + `strict_gate_clean`/`closing_review_round` close fields the protocol already promised | CLI: `internal/driver/impl.go:advanceReview` (lines 182-197) + new `ReadStrictGate` beside `ReadAutoImplement` (`internal/driver/cursor.go` / `transport.go`) + `driverImplOps.ReviewStatus` (`driver_impl.go:190-206`); Protocol: `COOPERATION.md` §Phase 8 strict gate stays, now machine-backed | Maker/checker; refutation re-verification pass | M | Med (extra review round cost) | **Adopt** |
| 3 | Enforce model diversity between implementer and the drafter/all reviewers; warn or refuse auto-complete when all checkers share the implementer's model | CLI: `internal/app/driver_impl.go:newDriverImplOps` (lines 35-54) comparing `agents.Discovery.Model` of implementer vs reviewers; Protocol: `COOPERATION.md` §Phase 6 add a normative line | Block #5: "ideally different models"; the named "Phase-6 false-green" | S | Low | **Adopt** |
| 4 | Generalize `RunChecks` so the pre/post-fix-up gates actually guard non-Go and design-only ideas; fail-closed for code-writing auto ideas with no checks | CLI: `internal/app/driver_impl.go:RunChecks` (lines 118-130) read a `checks:` command from `00-prompt.md` frontmatter / `~/.parley` `[defaults]`; Protocol: `COOPERATION.md` §Phase 4/5 document the `checks:` field | Stopping conditions / cost control; "verification is still on you" | S-M | Low | **Adopt** |
| 5 | Add an objective goal-condition (acceptance-criteria) check as a termination gate under `auto_implement`/`strict_gate`, reusing the existing consult mechanism as the "separate small model" | CLI: new driver step before `Complete()` (`internal/driver/impl.go:186-197`) evaluating FINAL.md observable criteria; reuse `internal/runner/consult.go` (currently advisory-only, `consult.go:51-63`); Protocol: `COOPERATION.md` §Phase 4 promote acceptance criteria from advisory to gate under auto mode | `/goal` directive; "Decide" row goal condition | M | Med (cost of an extra agent per close) | **Adapt** |
| 6 | Under `auto_implement`, don't auto-complete on `TriageReserved` — escalate reservations to HITL; keep current behavior in HITL mode | CLI: `internal/driver/impl.go:advanceReview` (line 182) split Reserved handling by `AutoImplement`; Protocol: none (signoff semantics unchanged) | Cognitive surrender / HITL fatigue guardrail | S | Low | **Adapt** |
| 7 | Refuse to auto-complete when the reviewer pool is a single agent (`len(reviewers) < 2`); require HITL signoff | CLI: `internal/app/driver_impl.go:newDriverImplOps` (lines 46-49) + `OpenReviewRound` (lines 132-135); Protocol: `COOPERATION.md` §Phase 6 line 376 already forbids *zero* reviewers — extend to "auto-complete requires ≥2" | Maker/checker: no solo checker that also grades itself | S | Low | **Adopt** |
| 8 | Make the "artifact-wins" fix-up override contingent on a real check, not just shape validation | CLI: `internal/runner/phase58.go:RunFixup` (lines 128-139) + `ValidateFixupArtifact` (lines 162-184) — tie the nonzero-exit override to `RunChecks` passing, not merely to a `## Fix-up cycle` section existing | "ship code you confirmed works" | S | Low | **Adapt** |
| 9 | Add a per-reviewer "Refutation attempts" attestation section to the review artifact schema, machine-validated | CLI: `ValidateReviewArtifact` (`phase58.go:397-420`) require a `## Refutation attempts` section; `BuildReviewPrompt` prompt for it; Protocol: `COOPERATION.md` §Phase 6 review file shape (lines 365-372) | Refutation-default; makes "no findings" earned | S | Low | **Adopt** |
| 10 | Do NOT make refutation/default-to-refuse the *only* close path for low-risk design-only ideas | Protocol: `COOPERATION.md` §Phase 6 — keep the default close rule (`OutstandingAgreedFixes==0`) for non-code/design-only ideas; refutation rigor scales with `auto_implement`/`strict_gate` | Avoid HITL fatigue / comprehension debt on trivial ideas | — | — | **Reject** (the over-application), keep conditional |

### Rationale for the top items

**#1 — refutation-default (Adopt).** This is the cheapest, highest-leverage
change and it is the heart of my lens. Today a reviewer can write an empty
`## Findings` and the artifact validates. Two edits: (a) in
`BuildReviewPrompt`, change the posture to "assume the implementation is wrong;
for each FINAL.md acceptance criterion, attempt to construct a failing case or
run the check; only report no findings after stating what you tried that failed
to break it"; (b) in `ValidateReviewArtifact`, require a `## Refutation
attempts` section (item #9) so an empty-findings review must show its work.
Mirror a one-line normative statement in `COOPERATION.md` §Phase 6. Zero
protocol risk; the only cost is reviewers do more work, which is the point.

**#2 — enforce strict_gate (Adopt).** This is the biggest false-green seam. The
protocol already specifies the strict close rule; the implementation already
named the exact fields (`strict_gate_clean`, `closing_review_round`) and the
exact parsing precedent (`ReadAutoImplement`). It was deferred and never done.
Under `strict_gate: true`, `advanceReview` must require a *fresh full-scope*
review round with zero findings before `Complete()` — not just
`OutstandingAgreedFixes == 0`. Until this lands, `strict_gate` is documentation,
not a gate, and the auto-loop's close condition is a single agent's self-report.

**#3 — model diversity (Adopt).** The `Model` field is already discovered and
logged. A guard in `newDriverImplOps` that warns (or, configurable, refuses to
auto-complete) when every reviewer shares the implementer's `Model` closes the
gap that AF3 left open at the model level. Add one normative line to §Phase 6 so
the protocol and CLI agree. In the current four-agent roster this is a no-op; it
matters the day someone runs a 2-agent same-model deck.

**#4 — generalize RunChecks (Adopt).** `RunChecks` returning `true` for every
non-Go workspace means the AF1 post-fix-up guard is fictional outside Go. A
`checks:` frontmatter field (parsed like `auto_implement`) lets each idea name
its verification command (`npm test`, `make ci`, `pytest`); for
`auto_implement` code-writing ideas with no `checks:` and no `go.mod`, fail
closed rather than auto-completing with no verification. This is conditional
rigor: the gate's strength matches the idea's risk.

**#5 — goal-condition check (Adapt, not full Adopt).** LE's `/goal` — "after
every turn a separate small model checks whether you are done" — is valuable but
an always-on separate-model check on every tick is cost/HITL-fatigue risk
(antigravity-1's lane). Adapt it: run the acceptance-criteria check *once,
before close*, and only under `auto_implement`/`strict_gate`. The existing
`consult.go` advisory mechanism (read-only, "NOT a pass/fail gate" today,
`consult.go:58`) is the right primitive to repurpose — promote it to a
structured goal-check for the close decision only.

**#10 — reject over-application (Reject).** Loop engineering's refutation
default is correct for code that ships. Applying it unchanged to low-risk,
design-only, or N/A-acceptance-criteria ideas would burn reviewer budget and
create the comprehension-debt/HITL-fatigue failure modes the corpus warns about.
The refutation rigor in #1/#2/#9 should scale with `auto_implement` and
`strict_gate`; the default `OutstandingAgreedFixes==0` close stays for trivial
ideas. Conditional rigor, not uniform rigor.
