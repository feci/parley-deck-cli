---
agent: codex-1
idea: protocol-and-skill-audit
round: 1
date: 2026-08-20
---

## Findings

### F1 — `consensus draft` treats blank participant files as a completed round
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --round 1 --by alice audit-empty 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Drafted consensus at /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-empty/consensus.md
Consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-empty/consensus.md
Missing signoffs: alice,bob
Signoffs: none
exit=0
```
contradicts: `parley-deck/COOPERATION.md:308-322`; `internal/consensus/consensus.go:118-129,486-493`
why it matters: PRIMARY — Both `round-01/alice.md` and `round-01/bob.md` in the isolated fixture contain only a newline, but the command printed “Drafted consensus” and exited 0. Phase 1 requires each artifact's frontmatter and four named sections, while this command's gate checks only whether each expected pathname exists. The manual CLI can therefore advance an idea with no participant analysis at all.

### F2 — `consensus draft --review` requires the implementer to review their own work
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --round 1 --by reviewer-a audit-review 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
consensus draft failed: review/round-01 is incomplete; missing impl.md
exit=1
```
contradicts: `parley-deck/COOPERATION.md:503-525`; `internal/consensus/consensus.go:118-129,486-493`
why it matters: PRIMARY — The isolated `standard` fixture lists `impl`, `reviewer-a`, and `reviewer-b`; both non-implementers supplied valid review files, and `impl` correctly supplied none. The command nevertheless applies the full idea participant list to `review/round-01` and refuses Phase 7. A protocol-compliant manual Phase 6 cannot reach review consensus without fabricating the expressly prohibited implementer review.

### F3 — Standard-track review consensus stays partial after every reviewer accepts
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --round 1 --by reviewer-a audit-review && /tmp/parley-audit-codex-1.DGmaGx/parley consensus signoff --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --agent reviewer-a --status accept audit-review && /tmp/parley-audit-codex-1.DGmaGx/parley consensus signoff --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --agent reviewer-b --status accept audit-review && /tmp/parley-audit-codex-1.DGmaGx/parley consensus status --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review audit-review 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Drafted consensus at /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl,reviewer-a,reviewer-b
Signoffs: none
Appended signoff for reviewer-a to /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl,reviewer-b
Signoffs:
  reviewer-a ✅ ACCEPT
Appended signoff for reviewer-b to /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl
Signoffs:
  reviewer-a ✅ ACCEPT
  reviewer-b ✅ ACCEPT
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl
Signoffs:
  reviewer-a ✅ ACCEPT
  reviewer-b ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:220-235`; `internal/consensus/consensus.go:92-102,370-426`; `internal/app/driver_impl.go:328-355`
why it matters: PRIMARY — The authoritative track table says `standard` review consensus is the reviewers who reviewed signing off. Both selected reviewers accepted, but the shared consensus validator always uses all idea participants and therefore demands `impl`. The auto-driver calls the same validator through `ReviewStatus`, so this is not confined to the manual command: the implemented standard-track quorum override does not bind at the close gate.

### F4 — `consensus finalize --by` accepts a non-participant who has no drafter claim
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus finalize --dir /tmp/parley-audit-codex-1.DGmaGx/repo --by mallory audit-final`
output:
```text
Finalized consensus and created /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/FINAL.md
Consensus: ready
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/consensus.md
Signoffs:
  alice      ✅ ACCEPT
  bob        ✅ ACCEPT
```
contradicts: `parley-deck/COOPERATION.md:398-400`; `internal/consensus/consensus.go:196-243,549-576`
why it matters: PRIMARY — `00-prompt.md` names `alice` as author and lists only `alice` and `bob` as participants; no artifact grants `mallory` a drafter handoff. The command nevertheless exits 0, prints “Finalized,” and writes `author: mallory`. The CLI's finalization authority is therefore an unchecked free-form flag rather than the initiator/accepted-volunteer rule.

### F5 — `consensus finalize` closes an idea with an empty `FINAL.md` outline
severity: MAJOR
evidence: PRIMARY
command:
`sed -n '1,120p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/FINAL.md && sed -n '1,12p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/00-prompt.md`
output:
```text
---
idea: audit-final
status: final
author: mallory
consensus-date: 2026-08-20
participants: [alice, bob]
---

## Final plan / specification

### Goal

### Scope

### Implementation details

### Tests

### Non-goals

### Verification

## References

- Consensus: ./consensus.md
- Rounds: ./round-01/
---
idea: audit-final
author: alice
created: 2026-08-20
track: standard
participants: [alice, bob]
status: final
---

## Problem / idea

Choose the better of two equally viable labels; this is a judgment call.
```
contradicts: `parley-deck/COOPERATION.md:402-430`; `internal/consensus/consensus.go:196-243,549-576`
why it matters: PRIMARY — The immediately preceding finalize command printed success and set `00-prompt.md` to `status: final`, but every substantive field in the generated source-of-truth artifact is empty, and the protocol-required Purpose, Context, Observable acceptance criteria, Idempotence, and Risks sections do not exist. The command's success message therefore overstates its own effect and permanently closes the idea around a scaffold instead of an authoritative specification.

### F6 — A standard-track unanimous judgment closes without §15.6's adversarial alternative
severity: MAJOR
evidence: PRIMARY
command:
`find /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final -maxdepth 1 -type d -name 'round-*' -print | sort && rg -n 'judgment call|Use the label blue|Choose blue|Adversarial alternative|related models|agreed position to be wrong' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/00-prompt.md /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/*.md /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/consensus.md; /tmp/parley-audit-codex-1.DGmaGx/parley status --dir /tmp/parley-audit-codex-1.DGmaGx/repo --idea audit-final 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/00-prompt.md:12:Choose the better of two equally viable labels; this is a judgment call.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/bob.md:10:Use the label blue.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/bob.md:14:Choose blue.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/alice.md:10:Use the label blue.
/tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/round-01/alice.md:14:Choose blue.
Idea: audit-final
Status: final
Participants: alice,bob
Consensus: ready
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-final/consensus.md
Signoffs:
  alice      ✅ ACCEPT
  bob        ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:1339-1358`; `internal/consensus/consensus.go:196-243`; `internal/driver/consensus.go:39-72,162-206`
why it matters: PRIMARY — The fixture explicitly identifies its output as a judgment call; both independent files select blue without disagreement; only `round-01` exists; and the grep finds no adversarial-alternative or correlated-agreement text. Nevertheless the CLI closed the idea. Neither manual finalization nor the auto-driver's FINAL validator checks §15.6, so a mandatory verification-integrity close condition is decorative.

### F7 — `continue` tells a fast-track idea to open the cross-review round that fast skips
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley continue --dir /tmp/parley-audit-codex-1.DGmaGx/repo audit-fast-run 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Run: audit-fast-run
Idea: audit-fast-plan
State: completed
Recommended: Open round-02 (cross-review) before drafting consensus
Command: parley run --auto --dir . "continue audit-fast-plan"
Next actions:
  open-next-round.audit-fast-plan  kind=open-next-round  risk=normal  requires=yes
      Open round-02 (cross-review) before drafting consensus
      command: parley run --auto --dir . "continue audit-fast-plan"
exit=0
```
contradicts: `parley-deck/COOPERATION.md:218-235`; `internal/runplan/runplan.go:112-148,227-279`
why it matters: PRIMARY — The fixture's `00-prompt.md` explicitly says `track: fast`, both round-01 artifacts exist, and the run records round 1 complete. The continuation planner ignores `track:` and consults only the separate `cross_review_rounds:` key, defaulting that key to 1. Its printed recovery instruction therefore drives the user into a phase the binding fast-track route explicitly skips.

### F8 — `continue` routes fast track through a separate `consensus.md` instead of collapsed FINAL
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley continue --dir /tmp/parley-audit-codex-1.DGmaGx/repo audit-fast-run 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Run: audit-fast-run
Idea: audit-fast-plan
State: completed
Recommended: Draft consensus from completed round artifacts
Command: parley consensus draft --round 1 audit-fast-plan
Next actions:
  draft-consensus.audit-fast-plan  kind=draft-consensus  risk=normal  requires=yes
      Draft consensus from completed round artifacts
      command: parley consensus draft --round 1 audit-fast-plan
exit=0
```
contradicts: `parley-deck/COOPERATION.md:218-235`; `internal/runplan/runplan.go:134-148`; `internal/driver/driver.go:301-315`
why it matters: PRIMARY — For this second run the isolated fixture still has `track: fast` but now also carries `cross_review_rounds: 0`, so the planner cannot take F7's extra-round branch. It next recommends the ordinary standalone consensus command. The authoritative fast route instead collapses consensus and FINAL into one `FINAL.md` with embedded signoffs; both the continuation planner and auto-driver always enter the ordinary consensus phase and contain no collapsed-final branch.

### F9 — A reservation need not be logged in deferred items; any filler text passes
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus signoff --dir /tmp/parley-audit-codex-1.DGmaGx/repo --agent alice --status reserve --notes 'Rollback is unresolved and must be designed during implementation.' audit-reserve && /tmp/parley-audit-codex-1.DGmaGx/parley consensus signoff --dir /tmp/parley-audit-codex-1.DGmaGx/repo --agent bob --status accept audit-reserve && /tmp/parley-audit-codex-1.DGmaGx/parley consensus finalize --dir /tmp/parley-audit-codex-1.DGmaGx/repo --by alice audit-reserve && sed -n '/## Open items deferred to implementation/,$p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/consensus.md 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Appended signoff for alice to /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/consensus.md
Consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/consensus.md
Missing signoffs: bob
Signoffs:
  alice      🟡 ACCEPT-WITH-RESERVATIONS — Rollback is unresolved and must be designed during implementation.
Appended signoff for bob to /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/consensus.md
Consensus: reserved
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/consensus.md
Signoffs:
  alice      🟡 ACCEPT-WITH-RESERVATIONS — Rollback is unresolved and must be designed during implementation.
  bob        ✅ ACCEPT
Finalized consensus and created /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/FINAL.md
Consensus: reserved
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-reserve/consensus.md
Signoffs:
  alice      🟡 ACCEPT-WITH-RESERVATIONS — Rollback is unresolved and must be designed during implementation.
  bob        ✅ ACCEPT
## Open items deferred to implementation

None.

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: alice — 2026-08-20
Status: 🟡 ACCEPT-WITH-RESERVATIONS
Notes: Rollback is unresolved and must be designed during implementation.

### Signoff: bob — 2026-08-20
Status: ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:384-392`; `internal/consensus/consensus.go:208-220,626-644`
why it matters: PRIMARY — The protocol permits 🟡 only if that reservation is logged under open items. Here the only section content is the contradictory word “None,” while the signoff names an unresolved rollback design; finalization still succeeds. The check proves merely that some text exists, not that the reservation survives into the implementation handoff.

### F10 — Manual review-consensus output uses a schema the protocol and auto-driver do not accept
severity: MAJOR
evidence: PRIMARY
command:
```text
/tmp/parley-audit-codex-1.DGmaGx/parley consensus draft --dir /tmp/parley-audit-codex-1.DGmaGx/repo --review --round 1 --by reviewer-a audit-review
sed -n '1,24p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md; rg -n 'review-cycle|outstanding_agreed_fixes|blocked:' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md || true
```
output:
```text
Drafted consensus at /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Review consensus: partial
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-review/review/consensus.md
Missing signoffs: impl,reviewer-a,reviewer-b
Signoffs: none
---
idea: audit-review
cycle: 1
drafted-by: reviewer-a
date: 2026-08-20
reviewed-commit:
---

## Agreed fixes

<!-- Review review/round-01 and record fixes that must be implemented. -->

## Deferred follow-ups

## Dismissed findings

## Signoffs

<!-- Each agent APPENDS their signoff block. Do NOT edit others' blocks. -->

### Signoff: reviewer-a — 2026-08-20
Status: ✅ ACCEPT

### Signoff: reviewer-b — 2026-08-20
```
contradicts: `parley-deck/COOPERATION.md:556-584`; `internal/consensus/consensus.go:508-527`; `internal/runner/phase58.go:378-387`; `internal/app/driver_impl.go:328-355`
why it matters: PRIMARY — The generator prints success but writes `cycle: 1`, while the protocol requires `review-cycle: N`; the grep also proves it writes neither `outstanding_agreed_fixes` nor `blocked`. The auto-driver later requires `outstanding_agreed_fixes` to be a non-negative integer. Thus the CLI's manual command produces an artifact that violates the documented schema and cannot be consumed by its own automated Phase 7/8 gate.

### F11 — The printed “Open round-02” command actually creates a new idea at round 1
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley continue --dir /tmp/parley-audit-codex-1.DGmaGx/repo audit-fast-run && /tmp/parley-audit-codex-1.DGmaGx/parley help | sed -n '/^  run$/,/^  resume$/p'`
output:
```text
Run: audit-fast-run
Idea: audit-fast-plan
State: completed
Recommended: Open round-02 (cross-review) before drafting consensus
Command: parley run --auto --dir . "continue audit-fast-plan"
Next actions:
  open-next-round.audit-fast-plan  kind=open-next-round  risk=normal  requires=yes
      Open round-02 (cross-review) before drafting consensus
      command: parley run --auto --dir . "continue audit-fast-plan"
  run
      Create a new idea from TASK and start round-01 with selected agents.
      Auto-drive is ON by default: after round-01 the protocol advances
      automatically (cross-review, consensus, finalize, and — if the idea opts
      into auto_implement — the implementation/fix-up phases), shown live in the
      TUI. Pass --no-auto to stop after round-01 and advance manually; with
      --no-auto and without --yes, Parley asks before launching hosted agents.

  resume
```
contradicts: `parley-deck/COOPERATION.md:324-350`; `internal/runaction/action.go:36-53`; `internal/protocol/workspace.go:121-181`
why it matters: PRIMARY — The CLI labels the action as opening `round-02` of `audit-fast-plan`, then prints a `parley run` command. Its own help confirms that `run` creates a new timestamped idea and starts `round-01`; `CreateIdeaFull` implements exactly that. A user following the recommended recovery command forks an unrelated idea instead of advancing the existing one, leaving the original stalled and creating a false audit trail.

### F12 — `init` reports a completed bootstrap without the mandatory roster/model/effort confirmation
severity: MAJOR
evidence: PRIMARY
command:
`mkdir -p /tmp/parley-audit-codex-1.DGmaGx/init-workspace && env PARLEY_HOME=/tmp/parley-audit-codex-1.DGmaGx/parley-home /tmp/parley-audit-codex-1.DGmaGx/parley init --dir /tmp/parley-audit-codex-1.DGmaGx/init-workspace && find /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck -maxdepth 2 -type f -print | sort && if test -e /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/agents.toml; then printf 'deck agents.toml: present\n'; else printf 'deck agents.toml: missing\n'; fi`
output:
```text
Initialized Parley Deck workspace at /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck
Central agent defaults: /tmp/parley-audit-codex-1.DGmaGx/parley-home/agents.toml (override per-project in /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/agents.toml)
/tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/COOPERATION.md
/tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/meta/version.json
deck agents.toml: missing
```
contradicts: `parley-deck/COOPERATION.md:43-57`; `internal/app/app.go:381-415`; `internal/protocol/workspace.go:37-70`
why it matters: PRIMARY — On a fresh isolated home and workspace, `init` accepted no roster, model, or effort input, wrote no deck `agents.toml`, and still printed “Initialized.” The protocol says this confirmation is a mandatory one-time bootstrap step tied to `parley init`, with persistent picks recorded in the deck roster authority. The success state therefore leaves a new deck before its mandatory bootstrap gate while telling the operator initialization is complete.

### F13 — The driver accepts an explicit standard track that its classifier rejects as under-tiered
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley classify --auto-implement --declared standard 2>&1; audit_classify_rc=$?; printf 'classify exit=%d\n' "$audit_classify_rc"; go run ./cmd/audit-track-probe`
output:
```text
deliberation
declared track "standard" is under-tiered; the classifier floor is "deliberation" (auto_implement)
classify exit=4
driver policy: track=standard apply=true reviewers=2 fixups=2 err=<nil>
```
contradicts: `parley-deck/COOPERATION.md:205-216`; `parley-deck/COOPERATION.md:246-254`; `internal/track/track.go:73-105`; `internal/track/track.go:129-172`; `internal/driver/driver.go:119-157`
why it matters: PRIMARY — Both outputs come from the same shipped `internal/track` package: the public classifier correctly says `auto_implement` raises the floor to deliberation and exits with the documented under-tier result, while the probe calls the policy used by the driver and shows that the driver silently applies standard's smaller reviewer and fix-up caps. A run can therefore be rejected by the advisory command yet accepted and executed under the very track it declares unsafe.

### F14 — Omitting `track:` does not apply the documented standard-track default
severity: MAJOR
evidence: PRIMARY
command:
`go test ./internal/driver -run '^TestAuditDefaultTrackPolicy$' -v`
output:
```text
=== RUN   TestAuditDefaultTrackPolicy
absent: track=standard reviewers=0 min=2 cross=9 cap=0 fixups=3
explicit standard: track=standard reviewers=2 min=2 cross=2 cap=2 fixups=2
--- PASS: TestAuditDefaultTrackPolicy (0.00s)
PASS
ok  	parley-deck-cli/internal/driver	0.304s
```
contradicts: `parley-deck/COOPERATION.md:200-203`; `parley-deck/COOPERATION.md:220-235`; `internal/track/track.go:129-172`; `internal/driver/driver.go:111-160`
why it matters: PRIMARY — The isolated test constructs otherwise identical four-participant drivers, changing only whether `00-prompt.md` explicitly contains `track: standard`. The protocol says an omitted track defaults to standard, so the configurations must be equivalent. Instead, omission enables an uncapped reviewer set, leaves a requested nine cross-review rounds uncapped, and grants three rather than two fix-up cycles. The driver labels the result `standard` while deliberately declining to enforce the standard row.

### F15 — An invalid explicit track silently disables every standard-track cap
severity: MAJOR
evidence: PRIMARY
command:
`set -o pipefail; go test ./internal/driver -run '^TestAuditDefaultTrackPolicy$' -v | sed -n '/unknown:/p'; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
unknown: track=standard reviewers=0 min=2 cross=9 cap=0 fixups=3 error=<nil>
exit=0
```
contradicts: `parley-deck/COOPERATION.md:200-216`; `parley-deck/COOPERATION.md:220-235`; `internal/track/track.go:30-43`; `internal/track/track.go:129-172`; `internal/driver/driver.go:126-160`
why it matters: PRIMARY — The prompt explicitly says `track: standart`. The driver returns no error, relabels it `standard`, but applies none of standard's binding caps: all reviewers remain selected, nine cross-review rounds remain permitted, and the legacy three-cycle fix-up budget survives. This is the opposite of the normative fail-safe rule: a typo quietly buys less enforcement than any valid track.

### F16 — The driver's cross-review gate accepts an empty `responding-to:` list
severity: MAJOR
evidence: PRIMARY
command:
`go test ./internal/driver -run '^TestAuditBlankRespondingTo$' -v`
output:
```text
=== RUN   TestAuditBlankRespondingTo
responding-to:
roundComplete(2)=true error=<nil>
--- PASS: TestAuditBlankRespondingTo (0.00s)
PASS
ok  	parley-deck-cli/internal/driver	0.296s
```
contradicts: `parley-deck/COOPERATION.md:328-349`; `internal/driver/driver.go:344-380`; `internal/driver/driver.go:466-469`; `internal/runner/validation.go:12-34`
why it matters: PRIMARY — Both round-02 artifacts in the isolated driver test have a literally empty `responding-to:` field. They contain the per-agent heading needed by the other half of the gate, but identify no prior artifact at all. The driver reports the round complete because `hasRespondingTo` checks only whether the key exists, while the protocol requires the field to list the prior files being answered. A cross-review can therefore advance without its required provenance.

### F17 — The auto-driver completes round 1 when every required section is empty
severity: MAJOR
evidence: PRIMARY
command:
`go test ./internal/driver -run '^TestAuditEmptyRoundOneSections$' -v`
output:
```text
=== RUN   TestAuditEmptyRoundOneSections
date fields=0 section-body-bytes=0
roundComplete(1)=true error=<nil>
--- PASS: TestAuditEmptyRoundOneSections (0.00s)
PASS
ok  	parley-deck-cli/internal/driver	0.269s
```
contradicts: `parley-deck/COOPERATION.md:308-322`; `internal/runner/validation.go:36-61`; `internal/driver/driver.go:344-380`
why it matters: PRIMARY — Each participant artifact in the isolated fixture has the three identity fields and the four required headings, but no `date:` and zero body content under every heading. `roundComplete(1)` still returns true and reconstructs a completion event. The production auto-driver uses this gate, so header-only scaffolds count as independent analyses and can advance automatically despite providing no analysis, approach, concerns, or risks.

### F18 — A review with no `reviewed-commit` passes the Phase 6 artifact validator
severity: MAJOR
evidence: PRIMARY
command:
`go test ./internal/runner -run '^TestAuditReviewWithoutCommit$' -v`
output:
```text
=== RUN   TestAuditReviewWithoutCommit
date=<missing> reviewed-commit=<missing> summary=<missing> open-questions=<missing>
ValidateReviewArtifact error=<nil>
--- PASS: TestAuditReviewWithoutCommit (0.00s)
PASS
ok  	parley-deck-cli/internal/runner	0.254s
```
contradicts: `parley-deck/COOPERATION.md:503-525`; `internal/runner/phase58.go:413-442`; `internal/app/driver_impl.go:299-325`
why it matters: PRIMARY — The isolated artifact has only agent/idea/review-round identity, one refutation sentence, and a Findings heading; it omits the protocol-required `reviewed-commit`, date, Summary, and Open questions. The shipped validator returns nil, and the driver's review-completion path relies on that validator. Most seriously, nothing proves which code revision was reviewed, so a stale or pre-fix review can satisfy the automatic close gate.

### F19 — `init` leaves its deck identity and creation date as template placeholders
severity: MINOR
evidence: PRIMARY
command:
`sed -n '1,8p' /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/COOPERATION.md; rg -n '^\*\*(Workspace|Created):\*\*.*<' /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/COOPERATION.md`
output:
```text
# COOPERATION.md — Multi-Agent Cooperation Protocol

**Workspace:** `<workspace-name>`
**Parley deck:** `./parley-deck/`
**Transport:** `local-dir`
**Created:** `<date> — created by parley init`
**Status:** Living document — any agent may propose changes via a dedicated idea (see §7).

3:**Workspace:** `<workspace-name>`
6:**Created:** `<date> — created by parley init`
```
contradicts: `../parley-deck-skill/skills/parley-deck/references/COOPERATION.md:3-6`; `internal/protocol/defaults/COOPERATION.md:3-6`; `internal/protocol/workspace.go:42-70`; `internal/protocol/workspace.go:94-105`
why it matters: PRIMARY — This is the fresh workspace from F12. The skill template explicitly labels its date placeholder “set at deck bootstrap,” and the embedded CLI template says the date is created by `parley init`; nevertheless `init` substitutes only the transport. Every new CLI-created protocol starts with false placeholder provenance and no workspace identity while the command reports initialization complete.

### F20 — Signoff-looking headings outside `## Signoffs` satisfy the consensus gate
severity: MAJOR
evidence: PRIMARY
command:
`/tmp/parley-audit-codex-1.DGmaGx/parley consensus status --dir /tmp/parley-audit-codex-1.DGmaGx/repo audit-signoff-scope 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
Consensus: ready
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-signoff-scope/consensus.md
Signoffs:
  alice      ✅ ACCEPT
  bob        ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:358-390`; `internal/consensus/consensus.go:330-366`; `internal/consensus/consensus.go:370-426`
why it matters: PRIMARY — The fixture's two `### Signoff:` blocks are under `## Agreed decisions`; the actual `## Signoffs` section contains only a comment saying nobody signed there. The parser scans every line in the document without tracking the section and reports consensus ready. Consequently a drafter can accidentally or deliberately make quoted/example signoff text count as append-only participant approval, bypassing the protocol's actual consensus gate.

### F21 — Consensus status ignores a frontmatter idea slug that names a different idea
severity: MAJOR
evidence: PRIMARY
command:
`sed -n '1,5p' /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-signoff-scope/consensus.md; /tmp/parley-audit-codex-1.DGmaGx/parley consensus status --dir /tmp/parley-audit-codex-1.DGmaGx/repo audit-signoff-scope 2>&1; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
---
idea: unrelated-idea
drafted-by: alice
date: 2026-08-20
---
Consensus: ready
Path: /tmp/parley-audit-codex-1.DGmaGx/repo/parley-deck/ideas/audit-signoff-scope/consensus.md
Signoffs:
  alice      ✅ ACCEPT
  bob        ✅ ACCEPT
exit=0
```
contradicts: `parley-deck/COOPERATION.md:358-390`; `internal/consensus/consensus.go:330-366`; `internal/consensus/consensus.go:370-426`
why it matters: PRIMARY — The command asks for `audit-signoff-scope`, but the canonical artifact explicitly declares `idea: unrelated-idea`. Status still reports ready because consensus parsing never reads or validates frontmatter. A copied or misplaced consensus can therefore authorize finalization of the wrong idea without any malformed-artifact signal.

### F22 — The auto-driver's FINAL gate accepts three arbitrary padded lines as a complete specification
severity: MAJOR
evidence: PRIMARY
command:
`go test ./internal/driver -run '^TestAuditFinalWithoutRequiredSchema$' -v`
output:
```text
=== RUN   TestAuditFinalWithoutRequiredSchema
purpose=false context=false acceptance-criteria=false idempotence=false risks=false references=false
finalScaffoldReason=""
--- PASS: TestAuditFinalWithoutRequiredSchema (0.00s)
PASS
ok  	parley-deck-cli/internal/driver	0.306s
```
contradicts: `parley-deck/COOPERATION.md:402-430`; `internal/driver/consensus.go:40-64`; `internal/driver/consensus.go:162-206`
why it matters: PRIMARY — The isolated FINAL has `status: final`, the one generic heading, two one-word lines, and a third padding line long enough to cross 250 bytes. It has none of Purpose, Context, observable acceptance criteria, Idempotence, Known risks, or References, and even declares the wrong idea slug. `finalScaffoldReason` returns the empty “acceptable” result used by the automatic consensus gate, so auto-drive can mark this non-specification final.

### F23 — The implementation gate accepts an unknown status and an empty one-heading artifact
severity: MAJOR
evidence: PRIMARY
command:
`go test ./internal/runner -run '^TestAuditImplementationInvalidStatus$' -v`
output:
```text
=== RUN   TestAuditImplementationInvalidStatus
status=banana implementer=<missing> branch=<missing> head-commit=<missing> plan=<missing>
ValidateImplementationArtifact error=<nil>
--- PASS: TestAuditImplementationInvalidStatus (0.00s)
PASS
ok  	parley-deck-cli/internal/runner	0.250s
```
contradicts: `parley-deck/COOPERATION.md:432-495`; `internal/runner/phase58.go:394-411`; `internal/app/driver_impl.go:256-276`
why it matters: PRIMARY — The isolated `IMPLEMENTATION.md` says `status: banana`, has an empty Summary heading, and omits implementer identity, dates, branch, commit, plan, deviations, reviewer notes, progress, decisions, validation, and outcomes. The production validator returns nil because it accepts any non-empty status and merely searches for the Summary substring. The auto-driver can therefore treat a non-reviewable scaffold with no implementation provenance as a finished Phase 5 artifact and launch review.

### F24 — Fresh-deck preflight calls an altered protocol “in sync” because two missing hashes compare equal
severity: MAJOR
evidence: PRIMARY
command:
`rg -n 'AUDIT MUTATION' /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/COOPERATION.md; sed -n '1,8p' /tmp/parley-audit-codex-1.DGmaGx/init-workspace/parley-deck/meta/version.json; set -o pipefail; env PARLEY_HOME=/tmp/parley-audit-codex-1.DGmaGx/parley-home /tmp/parley-audit-codex-1.DGmaGx/parley preflight --no-ping --dir /tmp/parley-audit-codex-1.DGmaGx/init-workspace 2>&1 | sed -n '/^Freshness:/p;/^  role=/p;/^Ready:/p'; audit_rc=$?; printf 'exit=%d\n' "$audit_rc"`
output:
```text
63:- AUDIT MUTATION: multi-agent execution is disabled.
{
  "protocolRole": "consumer",
  "deckVersion": "",
  "created": "2026-08-20"
}
Freshness: consumer — protocol matches packaged skill (in sync)
  role=consumer deckVersion=(none) classification=in-sync
Ready: no pending gates.
exit=0
```
contradicts: `parley-deck/COOPERATION.md:829-846`; `internal/protocol/workspace.go:73-89`; `internal/app/preflight.go:418-443`
why it matters: PRIMARY — This fresh deck's version metadata contains neither `protocolSha256` nor `packagedProtocolSha256`. After a deliberate body-level mutation replacing the mandatory multi-agent invariant, presence-only preflight still prints “protocol matches packaged skill,” “in-sync,” and “Ready.” The classifier equates the two absent JSON fields as empty strings without hashing the live or packaged protocol, so the freshness gate is guaranteed to fail open on every newly initialized deck until some other path happens to populate metadata.

## What I checked and found clean

- PRIMARY — The complete Go suite passed in the isolated tracked snapshot after the audit probes: `go test ./...` exited 0. Both command packages reported `[no test files]`; every package under `internal/` reported `ok`, including `app`, `consensus`, `driver`, `protocol`, `runner`, `runplan`, `track`, and `tui`.

- PRIMARY — The classifier itself produced the documented result and exit 0 for three non-contradictory cases:
  - `parley classify --protocol-change --declared deliberation --json` returned `{"declared":"deliberation","reason":"protocol-change","track":"deliberation","valid":true}`.
  - `parley classify --files 3 --loc 100 --reversible --mechanically-verifiable --declared fast --json` returned fast, `valid: true`, with reason `all fast conditions met (reversible, mechanically verifiable, small)`.
  - `parley classify --declared standard --json` returned standard, `valid: true`, with the documented default reason.

- PRIMARY — Manual signoff input controls bound correctly in four negative tests:
  - unknown `--agent mallory --status accept` returned `consensus signoff failed: unknown participant "mallory"` and exit 1;
  - `--status reserve` without notes returned `notes are required for 🟡 ACCEPT-WITH-RESERVATIONS` and exit 1;
  - `--status block --notes 'Not safe.'` without a counter-proposal returned `counter-proposal is required for ❌ BLOCK` and exit 1;
  - a second signoff from Alice returned `participant alice already signed` and exit 1.

- PRIMARY — The focused driver command below passed all eight named negative-path tests: goal-check failure escalates, consensus BLOCK respects the hard cross-review cap, list checks veto completion, an unclean strict-gate certification is vetoed, strict-gate retries are bounded, and explicit fast/deliberation non-solo contradictions escalate.
  `go test ./internal/driver -run '^(TestFastContradictionEscalates|TestFastNonSoloEscalates|TestExplicitDeliberationNonSoloEscalates|TestBlockedConsensusRespectsTheHardCrossReviewCap|TestStrictGateVetoesUncleanCertification|TestStrictGateBoundedByMaxFixupCycles|TestPhaseReviewListChecksVetoCompletion|TestGoalCheckFailEscalatesUnderAuto)$' -v`

- PRIMARY — `go test ./internal/driver -run '^TestRound02RequiresCrossReviewHeadings$' -v` passed. The driver does reject a round-02 artifact that omits a `### @<other>` heading; F16 is specifically the still-unchecked content of `responding-to:`.

- PRIMARY — `go test ./internal/runner -run '^(TestValidateReviewArtifactRequiresRefutation|TestValidateFixupArtifact)$' -v` passed both tests. An absent, inline-only, or empty `## Refutation attempts` section is rejected, and the fix-up validator's existing status/section cases bind.

- PRIMARY — Fresh `init` did correctly substitute `**Transport:** \`local-dir\`` and wrote `meta/version.json` with `protocolRole: consumer` and the real date `2026-08-20`. F12, F19, and F24 identify the separate bootstrap, header-placeholder, and empty-hash failures.

## What I could not check, and why

- PRIMARY — I did not invoke any hosted agent backend. That would cross the audit's no-secrets boundary by loading user CLI configurations and could incur external calls; all agent-facing behavior here was tested through deterministic CLI fixtures and Go seams in the isolated copy.

- PRIMARY — I did not perform GitHub PR or GitLab MR mutations, native review approvals, label transitions, merges, or transport-side rollback. The repository was explicitly read-only and those checks require external writes; this file therefore makes no PRIMARY claim about §11.B/§11.C host enforcement.

- SECONDARY — Filesystem tests cannot establish that a human or agent identity personally authored a markdown block, nor can they prove round-1 independence before publication. I tested the CLI's observable artifact gates only; actual identity/provenance enforcement would require a signed transport or live-agent evidence that this local-dir fixture does not provide.

- PRIMARY — I did not read or print user agent configuration files, credentials, tokens, or API keys. Consequently I did not test behaviors that depend on their private contents.
