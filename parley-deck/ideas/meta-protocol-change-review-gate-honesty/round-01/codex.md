---
agent: codex
idea: meta-protocol-change-review-gate-honesty
round: 1
date: 2026-06-12
---

## Summary

The amendment should be precise about what is normative for humans now and what is only machine-readable metadata for later driver enforcement. P6 belongs in Phase 6 because it governs review briefs; P7 belongs in Phase 8 because it changes the close condition only when an idea opts into `strict_gate: true`.

## Existing collision points

- Phase 6 currently defines reviewer artifacts and severities, but says nothing about brief neutrality or known dispositions (`parley-deck/COOPERATION.md:311-338`). P6 can be added here without changing file ownership or severity tags.
- Phase 7 already has `Agreed fixes`, `Deferred follow-ups`, and `Dismissed findings` (`parley-deck/COOPERATION.md:340-364`). P6 must not remove those categories; it only bans suppressive instructions in later briefs.
- Phase 8 currently closes when Phase 7 lists zero Agreed fixes (`parley-deck/COOPERATION.md:366-383`). P7 should add an opt-in exception: default stays zero Agreed fixes; strict gate additionally needs a fresh full-scope clean review or explicit operator ruling.
- Escalation already requires quoting the user's answer into the next artifact (`parley-deck/COOPERATION.md:397-416`). Reuse that for disputed finding closure; do not invent a second operator-ruling channel.

## P6 wording

Add after the current Phase 6 rules:

```
#### Review briefs and dispositions

Review briefs MUST NOT suppress findings. A facilitator, implementer, or prior
review consensus MAY describe known findings, rebuttals, accepted trade-offs,
sandbox artifacts, deferred follow-ups, and operator rulings as dispositions for
the reviewer to weigh openly. The brief MUST NOT say or imply "do not report",
"do not re-raise", "ignore", "only report above severity X", or otherwise narrow
what the reviewer may inspect or report.

When a brief includes a disposition, it SHOULD use this shape:

- Finding/disposition: <short identifier or summary>
  Prior disposition: rebutted | accepted trade-off | deferred | dismissed |
  operator-ruling
  Rationale: <one or two lines>
  Authority: <review consensus path, follow-up idea, or quoted operator answer>

The reviewer decides independently whether they concur with each disposition and
states that decision in their review file. A disputed finding closes only when
the reviewer withdraws it, the review consensus resolves it through the normal
signoff process, or the operator explicitly rules on it and that ruling is quoted
into the next review artifact.
```

## P7 wording

Add to `00-prompt.md` Phase 0 metadata:

```
strict_gate: true | false       # optional; default false; true uses the strict review-close rule.
```

Add after the existing Phase 8 close paragraph:

```
#### Strict review gate (optional)

An idea may opt into a strict review gate by setting `strict_gate: true` in
`00-prompt.md` frontmatter. If absent or any value other than true, the default
Phase 8 rule remains unchanged: the implementation may complete when Phase 7
consensus lists zero Agreed fixes.

For `strict_gate: true`, zero Agreed fixes is necessary but not sufficient. The
gate closes only after a fresh full-scope Phase 6 review round of the current
reviewed commit produces no findings of any severity or kind, and the subsequent
Phase 7 consensus records that clean result. A fix-verification or resumed pass
may converge the gate by checking prior fixes, but it never closes the gate by
itself. Findings classified as NIT, deferred follow-up, or accepted low severity
still keep the strict gate open unless the reviewer withdraws the finding or the
operator explicitly rules it closed.

`strict_gate` may be set at kickoff by the idea author. After kickoff, adding,
removing, or changing it requires either review/design consensus or explicit
operator direction recorded in the idea. A participant MUST NOT silently relax a
strict gate during implementation or review.
```

Add a trajectory subsection after that:

```
#### Stopping judgment

Review cycles are judged by trajectory, not by a pass counter. If findings are
fewer, lower severity, and confined to code changed by the latest fix-up, continue
within the configured fix-up budget. If fresh CRITICAL/MAJOR findings keep landing
on fix-up code, or the same ground is re-litigated despite open rebuttals, stop and
escalate with a short trajectory summary. If a finding requires an operator
decision, pause that finding's thread until the operator answers; unrelated fixes
may continue.

`MaxFixupCycles` and any driver retry budget are escalation thresholds, not close
criteria. Hitting the budget never marks an implementation complete; it requires
human review of the trajectory and either a new fix-up plan, a recorded operator
ruling, or a decision to abandon/defer the work.
```

## Machine-readable flag

Use the existing frontmatter precedent: `driver.ReadAutoImplement` reads a boolean from `00-prompt.md`, trims quotes, and defaults false (`internal/driver/transport.go:43-49`), backed by the simple top-level frontmatter reader in `internal/driver/cursor.go:282-314`. Define `strict_gate` the same way: exact case-insensitive `true` means enabled; absent/empty/other values mean false for compatibility.

Do not require driver enforcement in this protocol release. The code already completes default gates on `OutstandingAgreedFixes == 0` (`internal/driver/impl.go:186-197`) and escalates at `MaxFixupCycles` (`internal/driver/impl.go:200-204`). The sibling CLI idea may add `ReadStrictGate`, but the protocol amendment should be valid even before automation understands the flag.

If enforcement is later added, it should not scrape prose alone. It needs either a review consensus field such as `strict_gate_clean: true` plus `closing_review_round: N`, or a stricter review artifact verdict schema. That is a follow-up implementation detail, not required to state the protocol norm.

## Copies to sync

- `parley-deck/COOPERATION.md` - live project protocol.
- `internal/protocol/defaults/COOPERATION.md` - embedded default created by `parley init`.
- `parley-deck/meta/protocol-changelog.md` - dated changelog entry required by the live protocol (`parley-deck/COOPERATION.md:435-442`).

I did not find an in-repo `references/COOPERATION.md`; the only `references/` under these ideas is source material for this round. If a packaged skill snapshot outside this repo is part of the release process, update it in the release packaging step, but it is not an in-repo file to edit here.

## Position

Accept P6 and P7 with tight normative language. Keep default review closure unchanged, make `strict_gate` an opt-in boolean in `00-prompt.md`, and state clearly that strict closure is a human/protocol rule now even if the v1.24.0 driver only records or surfaces the flag.
