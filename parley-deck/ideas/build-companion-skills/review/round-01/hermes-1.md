---
agent: hermes-1
idea: build-companion-skills
review-round: 1
date: 2026-06-22
reviewed-commit: working-tree
---

## Summary

Reviewed `addons/parley-tracker/SKILL.md` and `templates/{epic,story,subtask}.md`
in `parley-deck-skill/` against the binding design
(`ideas/addon-skills-research/FINAL.md`, Skill B `ticket-tracker`). I also read
`bin/validate.js` and `bin/validate.test.js` because the lens asks whether the
templates actually pass the validate tool, and ran both against the shipped
templates and a placeholder-only fixture.

The skill is well-conceived and the three-audience contract, hybrid AC rules,
field schema, vendor-neutral mapping, and DoD model are faithfully transcribed
from FINAL.md §B. The SKILL.md is clear and internally consistent on the
authoring rules. The templates carry `At a glance` + `[B]`/`[T]`/`[A]` + the
hybrid AC structure, so the *shape* is right.

The problems are in the enforcement layer, which is the design's central thesis
("Enforcement is in the tool, not just the doc"). Concretely: (1) the shipped
epic template fails its own `validate` tool; (2) `validate` passes unfilled
placeholder skeletons, so the "validate exits 0" gate gives a false green and
does not stop an AI proceeding on a ticket it filled with assumptions; (3) the
`claim` enforcement gate described as "the single highest-leverage rule" is not
implemented; (4) `At a glance` and the `[B]`/`[T]` audience sections — the heart
of the three-audience lens — are not enforced. No tracker vendor lock-in was
found; the only couplings are Node (runtime) and a Parley-shaped
`canonical_source` placeholder, both acceptable for an addon.

Evidence is from real execution: `node bin/validate.js templates/epic.md` →
exit 1; `templates/story.md` and `templates/subtask.md` → exit 0 (with
`<...>` placeholders still intact); a synthetic ticket that is 100% placeholders
with empty `At a glance`/`[B]`/`[T]` → exit 0. The 16-test suite passes, but
every test uses a hand-filled fixture, so none of these gaps are guarded.

## Findings

### [CRITICAL] Shipped epic template fails its own validate tool

What is wrong: `templates/epic.md:49` defines `AC-E1 [B][T] The epic is
complete when every linked story is done and its acceptance criteria pass.`
This AC has no `Given/When/Then` and no `Verify:` line. `validate.js` check 4
(`acHasGherkin || acHasVerify`, lines 476-481) therefore flags it, and check 5
(lines 485-492) flags the epic as having no edge/error/offline AC. Running
`node bin/validate.js templates/epic.md` returns exit 1 with exactly those two
errors. `story.md` and `subtask.md` pass only because their placeholder ACs
contain the literal words `Given`/`When`/`Then`.

Why it matters: FINAL.md §B.4 and SKILL.md make `validate` exits 0 the very
first gap-scan step before any ticket is claimed. The template is the skeleton
every author copies. An author who copies the epic template and runs the
mandatory scan is blocked on step 1 before they have done anything wrong — and
the failure message points at structural rules, not at fields they need to fill.
It also means the skill ships a template that contradicts its own lint.

Concrete fix: Make the epic's roll-up criterion measurable with a `Verify:`
command and add an explicit edge waiver, so the template passes as a skeleton.
For example:
`- AC-E1 [B][T] Measurable: every child story under tickets/<epic-slug>/ has
status: done. Verify: \`grep -rL 'status: done' tickets/<epic-slug>/*/story.md\``
(empty output = all done), plus a line
`- Edge/error: n/a (an epic is an aggregate container; edge cases are owned by
its child stories)`. Alternatively, amend the design so epic-level aggregate
ACs are exempted from the per-AC Gherkin/Verify rule — but then `validate`
must encode that exemption, not just the doc.

### [CRITICAL] validate passes unfilled placeholder skeletons (false green)

What is wrong: `validate.js` checks for the *presence* of keywords, not whether
the ticket is filled. `acHasGherkin` (lines 329-334) matches `\bgiven\b`,
`\bwhen\b`, `\bthen\b` anywhere; `acHasVerify` (lines 336-338) matches
`\bverify\s*:`; `isEmptyValue` (lines 347-355) treats any non-empty string —
including `"<placeholder title still here>"` or `"<path-or-glob>"` — as filled.
I verified that a ticket whose title, ACs, and Non-goals are all literal
`<...>` placeholders, and whose `At a glance`/`[B]`/`[T]` sections are empty,
passes `validate` with exit 0. The shipped `story.md` and `subtask.md` also
pass with every `<...>` placeholder still intact.

Why it matters: This is exactly the "AI fills gaps anyway" failure the
no-assumption gap-scan exists to prevent. The design's highest-leverage rule is
that an agent runs the scan and "never proceeds by guessing." But an agent that
copies a template, leaves the placeholders, and runs `validate` gets a green
and proceeds — then fills the placeholders with assumptions during
implementation, unobserved. The gate proves structure, not readiness. The
"Readiness lint rejects vague tickets / AI-unready tickets" promise (SKILL.md
§Quality; FINAL.md §B.5) is not realized for the placeholder case.

Concrete fix: Add a placeholder-leak check to `validate`: flag any required
field, AC line, or `At a glance`/`[B]`/`[T]` body still containing `<...>` or
unreplaced template tokens. A copied template should *fail* until the author
replaces every placeholder (or marks the slot `n/a (reason)`). This makes the
shipped templates fail-by-default until filled, which is the correct
enforcement posture and pairs naturally with the epic fix above.

### [MAJOR] The `claim` enforcement gate is not implemented

What is wrong: SKILL.md (lines 216-219) and FINAL.md §B.4 state "the `claim`
operation RUNS the gap-scan and refuses to set `status: in-progress` on
failure… the single highest-leverage rule for AI output quality." The addon
ships exactly one tool, `bin/validate.js` (confirmed: the addon tree is
`SKILL.md`, `bin/validate.{js,test.js}`, `templates/*` — no `claim` command or
subcommand). There is no operation that atomically runs the scan and gates the
status transition.

Why it matters: The design's whole thesis is that enforcement lives in the
tool, not the doc, precisely because a doc-only rule is something an AI will
skip under pressure. With only a standalone `validate`, an agent can simply
edit `status: in-progress` without ever running the scan, and nothing stops it.
The "highest-leverage rule" is, in practice, a prose aspiration. SKILL.md
describes behavior the addon does not deliver.

Concrete fix: Either ship a `claim` operation (a thin wrapper that runs
`validate`, exits non-zero on failure, and only writes `status: in-progress` +
`assignee` on pass), or reframe SKILL.md to state honestly that `claim` is a
two-step manual procedure (`validate` then edit) and drop the "enforcement in
the tool" claim for the status gate. The former matches the design; the latter
is an honest downgrade.

### [MAJOR] `At a glance` is mandatory in prose but unenforced

What is wrong: SKILL.md (lines 92-96) calls `## At a glance` "MANDATORY" and
"the only mandatory cross-audience block… what makes a busy reader willing to
open the file at all." `validate.js` never checks for it. My placeholder
fixture had an empty `At a glance` and passed exit 0. The gap-scan's six
numbered steps (SKILL.md lines 201-210; FINAL.md §B.4) do not list it either,
so the implementation is faithful to the spec's checklist — but the spec
contradicts itself by calling the block MANDATORY yet omitting it from the
enforced scan.

Why it matters: This is the single most important readability block for the
three-audience lens, and the one the design says a stakeholder may read
*exclusively*. If it is not enforced, an AI can silently drop it and still
claim the ticket, and a business reader opening the file sees no summary. The
"tool-enforced" promise is hollow for the block that matters most to the
[B] audience.

Concrete fix: Add a check to `validate`: a non-empty `## At a glance` section
of 2-4 non-comment lines must be present for every ticket type. Amend the
gap-scan step list in SKILL.md/FINAL.md so "At a glance present and non-empty"
is an explicit enforced step (the implementation and the spec should move
together).

### [MAJOR] `[B]` and `[T]` audience sections are not enforced non-empty

What is wrong: `validate.js` (lines 452-461) enforces only the `[A] Agent
directives` section (non-empty + ≥1 "Do not"). There is no check that
`## [B] Business` or `## [T] Technical` exist or are non-empty. My placeholder
fixture had both empty and passed. The three-audience contract in SKILL.md
(lines 97-105) presents all three as peers.

Why it matters: The lens asks whether each ticket "genuinely serves all three"
audiences. With only `[A]` enforced, a ticket with a strong agent-directive
block and empty business/technical sections passes — i.e., a story whose
business value was never stated is claimable. This is precisely the "story with
only [T] ACs (a sign the business value was never stated)" anti-pattern the
design warns about (SKILL.md lines 183-185), extended to whole sections.

Concrete fix: Extend `validate` to require non-empty `[B]` and `[T]` sections
(mark `n/a (reason)` acceptable where genuinely non-applicable, e.g. some
subtask `[B]` content). Add these to the enforced gap-scan step list.

### [MAJOR] No test exercises the shipped templates

What is wrong: `bin/validate.test.js` has 16 passing tests, all driven by
synthetic hand-filled fixtures (`VALID_STORY`, etc., lines 42-87). No test runs
`validate` against `templates/epic.md`, `templates/story.md`, or
`templates/subtask.md`. The epic-template-fails-validate regression and the
placeholder pass-through are therefore unguarded — the green suite gives false
confidence.

Why it matters: The templates are the user-facing artifact; if they regress
against the tool, every author's first experience is a failure. A green test
suite that does not cover the shipped templates hides exactly the CRITICAL
findings above.

Concrete fix: Add a test that runs `validate` against each shipped template and
asserts the intended outcome. If placeholder-leak detection is added (per the
CRITICAL fix), the test should assert that an unfilled template *fails* and a
filled one *passes*. Add a fixture that is a verbatim copy of a template with
placeholders left in and assert it fails.

### [MINOR] Hybrid behavioural/NFR AC split is not enforced

What is wrong: SKILL.md (lines 166-181) specifies behavioural → Gherkin and
NFR → measurable bullet + `Verify:`. `validate.js` check 4 (lines 472-482)
only requires each AC to have an audience tag AND (`Gherkin OR Verify:`). It
cannot tell a behavioural AC from an NFR, so an NFR written as Gherkin passes,
and a behavioural AC with only a `Verify:` and no scenario passes. This is
faithful to FINAL.md §B.4 check 3 ("Gherkin or `Verify:`"), which collapses the
distinction — but it leaves the design's own Pitfall "Gherkin-for-NFRs
theatre" (SKILL.md lines 349-352) undetected.

Why it matters: The hybrid format is "correct and usable" as authoring
guidance, but its correctness depends on author discipline, not tool
enforcement — undercutting the tool-enforcement thesis for the AC layer
specifically.

Concrete fix: Introduce an AC-kind marker (e.g. `[NFR]` or a `kind:` prefix on
the AC) so `validate` can require NFR ACs to carry `Verify:` and behavioural
ACs to carry Gherkin, and warn when an NFR is expressed only as Gherkin. If the
design prefers not to add a marker, document explicitly that the hybrid split
is author-enforced and downgrade the "tool-enforced" claim for this rule.

### [MINOR] Happy-path AC requirement is stated but not enforced

What is wrong: SKILL.md (lines 175-181) and FINAL.md §B.3 require "≥1
happy-path behavioural AC AND ≥1 edge/error/offline AC." `validate.js` enforces
only the edge/error side (check 5, lines 485-492). A ticket with only
edge/error ACs and no happy path passes.

Why it matters: The paired requirement is what makes a ticket "done-when"
meaningful; enforcing only half of it lets a ticket with no positive-path
criterion be claimed.

Concrete fix: Add a check that at least one AC is a non-edge behavioural
(Gherkin) AC, or is marked as the happy path. Mirror the edge-waiver pattern
(`n/a (reason)`) for tickets that genuinely have no happy path.

### [MINOR] `parent` resolution only enforced in `--strict` mode

What is wrong: Gap-scan step 2 (SKILL.md lines 202-204) requires `parent`
resolves to an existing file. `validate.js` only resolves parents when
`opts.resolveParent` is provided, which happens exclusively in `--strict` mode
via `buildIdResolver` (line 730). The default single-file invocation
`node validate.js <file>` — the usage shown in SKILL.md's gap-scan — accepts
any non-empty, non-`n/a` parent string without checking it exists.

Why it matters: An agent following the SKILL.md gap-scan with the default
command gets no parent-resolution check, so a story pointing at a non-existent
epic passes step 2.

Concrete fix: Either make single-file mode attempt parent resolution against
sibling files in the same directory tree, or document explicitly that
`--strict --dir` is required for the parent-resolution step and direct agents
to use it in the gap-scan instructions.

### [MINOR] `canonical_source` traceability is not enforced

What is wrong: FINAL.md §B.2 and SKILL.md (line 255) state every AI-assigned
ticket carries `canonical_source` (path + revision) — the link that keeps
tickets subservient to `FINAL.md`. `validate.js` does not require it; it is not
in `REQUIRED_FIELDS` (line 55). A ticket omitting `canonical_source` passes.

Why it matters: This link is the mechanism for the core seam "a ticket cannot
override FINAL.md." If it is optional in practice, traceability erodes and the
design's two-canonical-artifacts invariant weakens.

Concrete fix: Require `canonical_source` non-empty (or `n/a (reason)`) for
`story`/`subtask` types at minimum, mirroring the `files`/`apis`/`arch`
populated-or-`n/a` check.

### [NIT] `validate` reports all errors, not "stop at first failure"

What is wrong: SKILL.md (line 199) and FINAL.md §B.4 say the scan runs "in
order, stopping at the first failure." `validateTicket` accumulates every error
and returns them all (line 510).

Why it matters: Minor spec deviation; returning all errors is arguably better
agent UX, but it contradicts the documented behavior and an agent relying on
"the first failure is the blocker" may be surprised by a multi-error report.

Concrete fix: Either keep accumulate-all and update the doc to say "reports all
failures; treat the first as the blocker," or short-circuit `validateTicket` at
the first error to match the doc.

### [NIT] `At a glance` "2-4 lines" guidance is not modeled in the templates

What is wrong: Each template's HTML comment says "MANDATORY 2-4 lines" but the
actual content is a single placeholder line (e.g. `epic.md:25`). With the
enforcement fix above, a naive author could write one line and pass a
"non-empty" check while missing the "2-4 lines" intent.

Why it matters: Low impact, but the templates should model the shape they
prescribe.

Concrete fix: Show a 2-line placeholder example in each template, and have
`validate` count non-comment body lines (2-4) once `At a glance` enforcement
is added.

### [NIT] Story template's example ACs carry no `[B]` tag

What is wrong: `templates/story.md` ACs are tagged `[A][T]`, `[A][T]`, `[T]`
(lines 61-63) — none carry `[B]`. The design (SKILL.md lines 183-185) says a
story with only `[T]` ACs is "a sign the business value was never stated."

Why it matters: The template models the anti-pattern the design warns about;
authors copying the shape may propagate it.

Concrete fix: Tag at least one story AC with `[B]` in the template (e.g.
`AC-1 [B][A][T]`), or add a `[B]`-tagged measurable AC example.

### [NIT] Node-only `validate` and Parley-shaped `canonical_source` placeholder

What is wrong: `validate` is Node-only (dependency-free within Node), and the
`canonical_source` placeholder hardcodes `parley-deck/ideas/<slug>/FINAL.md`.

Why it matters: Low. The skill is a Parley addon inside a Node-based installer,
so both are consistent with the ecosystem and do not constitute tracker vendor
lock-in (the schema is neutral, mapping is at the edge, no credentials in
core). Noting only for completeness against the "vendor lock-in" lens.

Concrete fix: None required. If desired, mention Node as the runtime
requirement for the gap-scan in SKILL.md so non-Node environments know they
need it, and note that `canonical_source` may point at any design doc for
non-Parley use.

## Open questions

1. Should `validate` reject unfilled `<...>` placeholders so copied templates
   fail until filled? This is the design decision behind the top CRITICAL. The
   "validate exits 0" gate implies filled tickets, but templates passing as
   skeletons creates the false-green. Needs agreement: add placeholder
   detection (templates fail until filled) vs. keep `validate` as a
   structure-only check and rely on author discipline.

2. Is the `claim` enforcement gate meant to be a real tool, or is "run
   `validate`, then set `in-progress`" an acceptable two-step manual procedure?
   The design promises tool enforcement; the addon ships none. This changes
   whether SKILL.md's "enforcement in the tool" wording stays or is downgraded.

3. Should `At a glance`, `[B]`, and `[T]` non-emptiness be added to the
   enforced gap-scan (amending FINAL.md §B.4's six-step list), or remain
   prose-only requirements? The spec currently calls `At a glance` MANDATORY
   but omits it from the enforced list — an internal contradiction to resolve.

4. Do epic-level aggregate ACs (e.g. "complete when all child stories done")
   belong under the per-AC Gherkin/Verify rule, or do they need an explicit
   exemption in the design and `validate`? The epic template failure suggests
   the rule does not fit roll-up criteria cleanly.
