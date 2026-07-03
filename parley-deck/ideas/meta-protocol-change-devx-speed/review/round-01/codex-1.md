---
agent: codex-1
idea: meta-protocol-change-devx-speed
review-round: 1
date: 2026-07-03
reviewed-commit: a224621
---

## Summary

I cannot sign off on this implementation as complete. The drift guard is green, and the
new text does restate the major safety invariants, but the reduced-track table is not
integrated with the existing phase, quorum, and transport rules. A facilitator can now
find contradictory instructions for the same fast or standard idea.

The recorded deviations are also not merely cosmetic for this idea: CLI/script
enforcement, physical core/appendix restructuring, the protocol changelog, and the
protocol hash are all named observable acceptance criteria in FINAL.md. They need either
implementation here or an explicit owner-approved amendment to FINAL.md/review consensus.

## Refutation attempts (what you tried, and the result per acceptance criterion)

- Checked new safety invariants against existing MUST rules: non-solo, refutation-default,
  round-1 independence, append-only signoffs, file-canonical audit trail, English-only,
  no-secrets, and the section 14 human brake are restated in the new Phase 0.0 block. Result:
  the invariants are named, but several older phase rules still conflict with the new
  reductions, especially quorum, consensus, reviewer count, and review signoffs.
- Tried to route risky or ambiguous examples through the classifier. Result: deliberation
  triggers catch obvious protocol/security/data/pipeline/API-break cases, but the fast
  predicates use approximate and subjective terms (`~3-5 files`, `fully reversible`,
  `data surface`, `mechanically verifiable`) with no deterministic fail-closed rule for
  unknown inputs. That does not satisfy "script-checkable" routing.
- Ran the required copy-drift check from the CLI root:
  `GOCACHE=/private/tmp/parley-gocache go test ./internal/protocol/...`. Result:
  `ok parley-deck-cli/internal/protocol 0.313s`. Raw SHA-256 hashes of the live and
  embedded copies differ, but the test passes under the repository's allowlisted zones.
- Checked cross-references and numbering around the new block. Result: references say
  `section 4.0`, but the actual heading is `Phase 0.0`, and the old `Phase 0` heading still
  follows it. That is navigable for humans, but it is a numbering collision introduced by
  the new text.
- Checked the recorded deviations against FINAL.md acceptance criteria. Result: the
  reading-guide replacement for physical appendix relocation leaves acceptance criterion 4
  unmet; CLI/driver enforcement deferral leaves criteria 1-3 partly unmet; acceptance
  criterion 6 is explicitly still pending in IMPLEMENTATION.md.

## Findings

### [CRITICAL] Reduced tracks are not reconciled with the existing phase and quorum rules

`COOPERATION.md:187-195` defines fast/standard reductions: skip or cap Phase 2, collapse
Phase 3-4 for fast, reduce reviewer counts, cap fix-up cycles, and auto-advance. The older
normative text still says:

- `COOPERATION.md:267-289`: cross-review continues until nobody has substantive objections.
- `COOPERATION.md:312-324`: every listed participant appends consensus signoff and every
  active participant must accept.
- `COOPERATION.md:432-456`: every active participant except the implementer writes a review.
- `COOPERATION.md:510-513`: every active participant, implementer included, signs review
  consensus.
- `COOPERATION.md:648-655`: quorum is set by the section 9.0 readiness check, but fast says
  that check is skipped.
- `COOPERATION.md:772+`: transport mechanics still describe the full all-participant flow.

This makes fast and standard non-operational as written: following the table violates the
old MUST rules, while following the old rules fails FINAL.md's reduced-ceremony acceptance
criteria. It also creates a safety ambiguity because a facilitator can choose the looser
interpretation without a single authoritative gate.

Concrete fix: update Phase 2, Phase 3, Phase 6, Phase 7, Phase 8, section 5, section 9.0,
and the section 11 transport mechanics to explicitly say which track-specific rule
overrides the old full-lifecycle default. Include the missing Phase 7 review-consensus row
from FINAL.md (`fast`: reviewer accept equals consensus; `standard`: reviewers who reviewed
sign off; `deliberation`: all sign off), and define the minimal fast-track liveness/quorum
check that preserves non-solo without the full readiness ping.

### [MAJOR] CLI/script enforcement is deferred, leaving acceptance criteria 1-3 unmet

`IMPLEMENTATION.md:48-54` defers classifier command support, `parley init/run` templating,
per-track timeout seeding, driver auto-advance, and validation gates to `track-aware-driver`.
But FINAL.md acceptance criteria require more than prose:

- Criterion 1: absent `track:` defaults to `standard`, and a script can classify objective
  inputs with deliberation triggers forcing deliberation.
- Criterion 2: fast/standard behavior is verifiable on a sample idea end-to-end.
- Criterion 3: removing non-solo/refutation is rejected by driver validation.

The protocol text is useful, but it is not the deterministic routing and validation the
ratified spec required for this idea.

Concrete fix: implement the classifier/defaulting/timeout templating/driver validation in
this implementation, or get an explicit owner-approved amendment that narrows this idea to
protocol prose and moves criteria 1-3 to a named follow-up with adjusted acceptance criteria.

### [MAJOR] The classifier is not objective or unambiguous enough for safe fast routing

The classifier at `COOPERATION.md:181-183` is exhaustive only because `standard` catches
everything else, but fast eligibility is not script-checkable. `<= ~3-5 files`, `fully
reversible`, `no security or data surface`, and `mechanically verifiable` require judgment.
There is also no explicit rule that unknown or unclassifiable inputs must fail closed to
`standard` or `deliberation`.

Risk example: a four- or five-file change touching telemetry, dependency behavior, or a
low-level config surface can plausibly be described as reversible and test-covered, yet still
not be safely "fast" unless the ambiguous predicates are resolved conservatively.

Concrete fix: replace approximate thresholds with exact numeric inputs, define the boolean
classifier fields, and add a fail-closed rule: if any fast predicate is unknown or disputed,
route to `standard`; if any deliberation trigger is unknown but plausible for security,
privacy, production infrastructure, data migration, pipeline/action, API break, or schema
break, route to `deliberation` until disproven.

### [MAJOR] The appendix/core restructuring acceptance criterion is not met

`IMPLEMENTATION.md:34-46` records that the physical appendix relocation and core-before-first
appendix requirement were replaced by a reading guide. That does improve orientation, but it
does not satisfy FINAL.md acceptance criterion 4. In the live file, `section 11` starts at
`COOPERATION.md:772`, the first actual appendix starts at `COOPERATION.md:980`, and sections
12-14 remain after that. The core is therefore not physically <= about 200 lines before the
first appendix, and section 9/11/12/13/14 are not all appendices by position.

Concrete fix: either perform the physical move/renumber with a cross-reference audit, or get
owner-approved review consensus that explicitly accepts the reading-guide substitution and
rewrites acceptance criterion 4.

### [MAJOR] Protocol changelog and protocolSha256 handling are still pending

FINAL.md acceptance criterion 6 requires `meta/protocol-changelog.md` to record the change
and `protocolSha256` to be handled per section 9.0. `IMPLEMENTATION.md:72` marks this as
pending. I found no 2026-07-03 entry for this idea in `meta/protocol-changelog.md`, and
`meta/version.json` still has the old `protocolSha256` value while the current live
`COOPERATION.md` SHA-256 is `bac3711cb206cf3ecf84f6b8e30f3bbd3d2687dfa5770c7a2effb2fa11ee6741`.

Concrete fix: add the protocol changelog entry and update the protocol metadata through the
established sync/version path, or explicitly document why the source deck defers the metadata
write and what release step will perform it before completion.

### [MINOR] The new section is referenced as section 4.0 but headed as Phase 0.0

Quickstart and `track:` comments point readers to `section 4.0`, but the actual heading is
`### Phase 0.0 -- Track selection`, followed by the existing `### Phase 0 -- Kickoff`.
This is a small navigation bug and a numbering collision.

Concrete fix: rename the heading to match the references, for example
`### 4.0 Track selection (conditional rigor)`, or update all references to say
`Phase 0.0` and explain how it relates to Phase 0.

### [MINOR] Quickstart over-compresses the fast track

`COOPERATION.md:21` says "`fast` = one review, then done." That can be read as skipping the
mandatory round-1 artifact, collapsed FINAL/signoff audit, and possible fix-up verification.
The invariants later prohibit that, but Quickstart is the entry point for readers most likely
to miss the nuance.

Concrete fix: change the phrase to something like: "`fast` = round-1, collapsed FINAL
signoff, one refutation-default reviewer, and at most one fix-up cycle."

## Open questions

- Which skill fallback copy was actually re-synced? The installed fallback at
  `/Users/tomasfecko/.codex/skills/parley-deck/references/COOPERATION.md` still differs from
  the live protocol, and I found no in-repo `parley-deck-skill/references/COOPERATION.md`.
- Were the appendix relocation and CLI enforcement deferrals explicitly owner-approved after
  FINAL.md, or are they implementer-proposed deviations awaiting review consensus?
