---
idea: meta-protocol-change-quota-auto-exclude
review-cycle: 2
outstanding_agreed_fixes: 5
blocked: false
drafted-by: codex-1
date: 2026-10-05
reviewed-commit: d8b729af187a18ebe3891a36dfc379c4b04c8b95
product-commit: 906857b9af33306158d3b0e91b94e6ef4b6f7458
skill-commit: dcb7d593130247709cfad19715aa1f53e5cc8be0
---

## Scope and review basis

The separate claude-1 round-04 review is the binding full-scope re-review after fix-up cycle 1.
It completed the single owner-authorized relaunch without quota/auth failure and found 0 CRITICAL,
3 MAJOR, 2 MINOR and 1 NIT. Both earlier CRITICALs and most G1–G9 are resolved within the reviewer's
stated scopes. codex-1 accepts the findings as implementer; it supplies no independent code verdict.
This draft authorizes cycle 2 only after both participants sign. It authorizes no close or release.

Trajectory: 15 round-03 findings became six round-04 findings after one fix-up. Both CRITICALs are
resolved. R4-MAJOR-1 concerns new retry-framing fix code; R4-MAJOR-2 is the remaining G2/R2 compatibility
issue; R4-MAJOR-3 is the previously disclosed native-evidence gap. Continue the bounded unrelated repairs
under Phase 8 stopping judgment, rather than treating either green tests or a pass count as acceptance.
Cycle 2 of the five-cycle cap is proposed. No D6/legacy-accounting change, roster change or excluded-agent
invocation is included.

## Agreed fixes

- **G10 — real retry semantics and actionable recognizer reasons** (round-04/claude-1 R4-MAJOR-1,
  R4-NIT-1): accept complete allowlisted retry reports whose exhaustion class and absolute machine reset
  agree within the owner's one-second bound even when their countdown text changes. Check each record's
  internal duration/header/countdown consistency, distinguish its observation from one aggregated dump's
  receipt time, and apply the fixed 60-minute threshold against the actual observation. Keep exact
  lastError fingerprint equality and all mixed-status/class/contradiction, artifact and later-success
  gates. Add the two retained countdowns and realistic backoff deltas as explicitly source-derived
  fixtures. Return the specific framing/consistency rejection instead of calling a supported recognizer
  unsupported. No change to the owner exception, no relabeling of fixtures as native.
- **G11 — preserve the ordinary knob-off path** (R4-MAJOR-2 option **(b)**, R4-MINOR-1): restore the
  pre-existing protocol's recorded confirmation and catch-up forms for policy-off ideas. Do not require
  a new included marker, a new exact directive, a mandatory committed answer schema or a new command
  for this existing path. Preserve explicit owner-confirmed re-inclusion, late round-1/read-priors/join
  from round-2 catch-up, historical-known/current-required membership, retained obligations and immutable
  revisions. Do not treat an unconfirmed prompt edit as owner authorization. Use the pre-change path
  and existing protocol forms as the compatibility oracle. Make exclusion parsing tolerate reasons
  containing an em dash and report the accepted grammar/action on malformed input. Document the real
  accepted forms in CLI and skill guidance. Choosing (b) avoids proposing a third knob-off behavior
  change; if preserving this contract proves infeasible, report the precise blocker before changing it.
- **G12 — visible policy-on revision/recovery boundary** (R4-MINOR-2): name parley quota revise in the
  contradictory-projection diagnostic. A prompt edit after a fully applied authoritative revision must
  not be silently overwritten merely because it matches an older revision. Distinguish a genuinely
  interrupted pending projection (normal checked recovery) from a later manual edit; the latter gets
  an owner-visible escalation with the required correction/revise path. Do not invent authority from
  the edited prompt or break idempotent pending-transition replay. Clarify the policy-on path in the
  protocol/skill without broadening owner policy.
- **G13 — close the recovery leads with evidence** (round-04 open questions 2 and 3): reproduce a writer
  whose parley process crashed before terminal.json, and manual import with a missing kickoff manifest.
  Establish and document a checked recovery route within FINAL's durability/recovery requirements, or
  fix a confirmed defect with regression coverage. Never infer that an unknown/foreign writer stopped
  from local PID absence, and never dispatch with an unresolved writer. Validate prerequisites before
  immutable commit where possible; already committed incomplete state remains fail-closed until checked
  recovery. Report any action requiring owner judgment rather than adding an unsafe escape hatch.
- **G14 — prepare native-framing evidence; retain the AC2 owner gate** (R4-MAJOR-3 and open question 1):
  locate the actual installed zcode terminal-error sink and inspect settings from source; use local
  source inspection and offline SDK/Node reproductions only. Label source-derived framing, partial native
  tails and complete native capture distinctly. Do not invoke zcode, a provider or another participant,
  loosen the allowlist, or declare AC2 met. The unrecovered complete native capture remains an OPEN
  owner decision after the authorized repair work is concrete. Either owner-authorized native capture
  or an explicit owner ruling is required; no assumption or signoff here substitutes for it.

## Owner-dependent thread

R4-MAJOR-3 / AC2 remains unresolved and is not deferred out of this idea or waived. Only its native-capture
thread is held; unrelated G10–G13 work may proceed after this fix-list consensus under Phase 8 stopping
judgment. Before requesting a decision, finish the authorized repairs and source/evidence preparation
so the owner sees a concrete result. The later blocking note must state exact capture scope or exact
proposed evidence ruling and the remaining risk. No zero-fix/close claim can be based on this draft.

## Deferred follow-ups

- D6 / driver gap 11: separate lasting legacy-accounting fix, follow-up idea TBD after this idea.
- Windows runtime remains unexecuted; compile-only evidence is not a runtime pass. CLI winget stays held.
- Unsupported-adapter native provenance and the Claude JSON-output change remain FINAL's follow-ups.

## Dismissed findings

None of the new round-04 findings is dismissed. The reviewer's prior dispositions stand: G1, G3,
G5–G9 and R1/R3 are resolved in the scopes stated in round-04; G2/G4 and R2 retain the specific new
residues above. Completed-drafter protection and literal AC1 skill/deck equality retain their signed
round-03 dispositions. No deferred item is used to erase a current requirement.

## Coverage & blind spots

claude-1 read the entire current protocol body, all 24 CLI product chunks and both skill product chunks,
and filed coverage over every changed product file. Its 13-package focused checks and actual shared/local
probes are PRIMARY within their recorded scopes. Full host suite/race/vet/skill checks are producer
SECONDARY evidence; Windows runtime and complete native zcode capture remain unexecuted/unavailable.
There is one independent reviewer. codex-1 organizes and implements under the recorded §15.5 role
concentration. A new full-scope re-review, final review consensus with both signoffs and a NEW attended
owner close are still required. No acceptance of code is inferred from agreement to this fix list.

## Signoffs

<!-- Each participant appends only its own block. Signoff accepts this fix list, not the current code. -->

### Signoff: codex-1 — 2026-10-05
Status: ✅ ACCEPT
Notes: I accept G10–G14 as the cycle-2 implementer plan. G11 uses option (b) to preserve existing knob-off behavior without a new owner-policy deviation. AC2/native capture stays open for owner judgment after concrete repairs and source evidence. This is not a code verdict or permission to close, merge or release.
