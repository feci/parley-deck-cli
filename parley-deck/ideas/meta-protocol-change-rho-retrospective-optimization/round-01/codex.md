---
agent: codex
idea: meta-protocol-change-rho-retrospective-optimization
round: 1
date: 2026-06-16
---

## Summary

RHO should enter parley-deck as both a small protocol amendment and optional tooling, with a hard boundary: the protocol defines eligibility, evidence, artifacts, and acceptance gates; the tool only mines history and prepares reviewable retro material. A retro pass must never directly edit `COOPERATION.md`, bundled defaults, `SKILL.md`, skills, memory, or helper scripts. It may propose edits through the normal Parley Deck lifecycle, and those edits are accepted only by multi-agent consensus, human approval where protocol or harness behavior changes, and a strict no-regression gate.

From a Go/tooling feasibility lens, the minimal viable implementation is realistic if it is mostly read-only: a `parley retro` command can scan structured Parley artifacts, compute a diverse-hard coreset, generate a diagnostic report and candidate proposal files, and leave all writes behind explicit review steps. The risky part is not parsing or ranking artifacts; it is letting mined self-judgments become persistent behavioral changes without sufficient quarantine and review.

## Proposed approach

Add a concise "Retrospective optimization" section to the protocol and, separately, specify a `parley retro` tool. The protocol section should be normative and small: retrospective optimization is an advisory design workflow that mines prior Parley Deck artifacts to propose improvements; it does not apply improvements. The tooling spec should be non-normative except for safety behavior such as read-only default mode, explicit output paths, provenance capture, and refusal to touch live harness files.

Define our harness in layers:

- Protocol harness: `parley-deck/COOPERATION.md` and the internal default copy. A retro pass may only recommend changes through a meta-protocol-change idea. It must not patch either file directly.
- Runtime harness: Parley Deck skills, embedded instructions, CLI behavior, helper scripts, and workflow defaults. A retro pass may draft a design or implementation proposal, but any change routes through the ordinary idea, implementation, review, and signoff phases.
- Local harness: machine-local memory, ignored headless-agent config, and session launch preferences. A retro pass may report observations about them but must not make them canonical or infer protocol rules from one operator's local setup.
- Evidence corpus: structured artifacts under `parley-deck/ideas/*`, review files, consensus/signoff files, `IMPLEMENTATION.md`, run logs if present, and raw session JSONL only as secondary evidence.

Replace RHO's single-model pairwise self-preference with Parley Deck consensus. A candidate retro proposal should be accepted only when all of these hold: the proposal has a normal idea with participant rounds; all active participants sign off; required human approval is explicit for protocol or harness changes; drift guard passes for protocol-copy lockstep; the relevant test suite or documented checks are green; review finds no agreed blocking regressions. "Mean self-preference is positive" is not an acceptance gate here; at most it is one diagnostic signal that can be attached to the proposal.

Use a structured coreset over our artifacts rather than raw trajectories first. The hard-case score should favor ideas with multiple review rounds, fix-up churn, blocked or abandoned states, escalations to the user, late protocol deviations, dismissed-but-recurring findings, drift-guard failures, long consensus time, or implementation/review mismatch. The diversity score should cover different failure modes: protocol ambiguity, transport mechanics, role/quorum confusion, artifact ownership collisions, review-gate weakness, tooling/runtime friction, and safety/adversarial concerns. A Go implementation can start with deterministic feature extraction and simple embeddings or tags, then optionally add DPP-style selection later; deterministic ranking is easier to review and test.

Tooling should be staged:

- `parley retro scan`: read-only inventory of eligible artifacts and extracted signals.
- `parley retro select`: produce a reviewable coreset manifest with reasons for each selected case.
- `parley retro diagnose`: produce a diagnosis report grouped by recurring workflow or tooling failure.
- `parley retro propose`: draft proposal text or candidate patch descriptions, but only into an explicitly requested idea path and never over an existing participant file.

For Go feasibility, this fits well as ordinary CLI plumbing: frontmatter parsing, Markdown section extraction, filesystem walking with ignore rules, JSON output for auditability, and deterministic tests using fixture idea directories. The first version should avoid replaying agents or editing harness directories. Re-rollout can be a later, opt-in extension because it introduces cost, environment reset requirements, and much more complicated safety semantics.

## Concerns / open questions

The main boundary question is whether `parley retro propose` should create a new idea scaffold or only print instructions for the facilitator. I prefer allowing scaffold creation only when the operator passes an explicit slug and the target does not exist, because predictable file output is testable and reviewable; it must never create participant analysis for another agent.

Raw session JSONL is useful but should not be primary in the first implementation. It is large, provider-specific, more likely to contain injected or sensitive content, and harder to connect to final protocol outcomes. Structured Parley artifacts already encode the safer signal: what agents objected to, what survived consensus, what failed review, and what required human judgment.

The protocol should specify retention and provenance expectations. A retro proposal should cite the source artifacts and extracted features that led to each recommendation, without copying sensitive raw trajectories into canonical files. If a source artifact is adversarial, compromised, or outside the project scope, the retro pass should mark it excluded rather than distill from it.

We also need a practical definition of "no regression" for non-code protocol changes. At minimum it means protocol drift guard green, existing protocol fixtures green if present, no loss of mandatory multi-agent/quorum/human-approval guarantees, and no accepted objection from review that the new wording weakens artifact ownership or signoff semantics.

## Risks

Self-preference bias can become protocol entrenchment if the tool is allowed to rank its own suggestions as sufficient evidence. Mitigation: self-preference output is advisory only; acceptance requires normal Parley Deck consensus and human approval for protocol or harness changes.

Adversarial or noisy trajectories can poison durable instructions. Mitigation: prefer canonical structured artifacts, quarantine raw logs, record provenance, exclude suspicious sources, and require reviewers to inspect the evidence trail before accepting any persistent behavior change.

Selection noise can overfit the protocol to dramatic but rare failures. Mitigation: require diversity across failure types, cap repeated incidents from the same idea family, and keep the first protocol change minimal.

Tooling could accidentally mutate canonical state. Mitigation: default read-only operation, explicit write flags, exact path allowlists, "fail if exists" behavior for generated protocol artifacts, dry-run output, and tests that assert no writes outside the requested retro output path.

Re-rollout can be expensive and unsafe for irreversible tasks. Mitigation: omit re-rollout from the first version; if added later, require resettable fixtures, sandboxed worktrees, no network by default, and a separate implementation proposal.
