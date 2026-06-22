---
agent: codex-1
idea: build-companion-skills
review-round: 1
date: 2026-06-22
reviewed-commit: working-tree
---

## Summary

The addon installer path mostly works: default install lays down core plus both add-ons, `--no-addons` and `--only` select the expected install units, unknown `--only` is rejected through the CLI, and installer-owned add-on directories do not self-detect as runtimes. I ran `npm test` in `parley-deck-skill` and it passed: 48 tests, 0 failures.

The main blockers are correctness gaps around health/status behavior for intentional core-only installs and the tracker validator accepting tickets that do not meet the binding `FINAL.md` readiness contract.

## Findings

### [MAJOR] Default doctor/status treat intentional core-only installs as broken

What is wrong: `selectedAddons()` recomputes the current package default add-on set for every command, and `doctorCommand()` requires every selected unit to be valid. After `install --no-addons`, a plain `doctor` reports the core as valid but the two intentionally omitted add-ons as `missing`, making the overall doctor result fail. This also affects legacy/backward-compatible core-only installs from before add-ons existed. The relevant path is `lib/installer.js:341`, `lib/installer.js:353`, `lib/installer.js:774`, `lib/installer.js:790`, and `lib/installer.js:1000`.

Why it matters: `--no-addons` is a supported mode, and backward-compatible existing core installs should not become unhealthy just because the new installer now ships add-ons by default. Operators will see a false negative unless they remember to pass `--no-addons` to health commands forever.

Concrete fix: persist install selection in the core marker, for example `addons: ["parley-tracker", "parley-worktrees"]` or `addons: false`, and have doctor/status/paths derive the expected add-on set from the installed marker when no selection flag is supplied. Alternatively, make default doctor validate the core and report absent add-ons as `not-installed` warnings, not `ok: false`. Add tests for legacy core-only marker, `install --no-addons` followed by plain doctor/status, and explicit `doctor --only ...`.

### [MAJOR] `validate` accepts non-executable measurable ACs and misses the happy-path requirement

What is wrong: `acHasVerify()` only checks that `Verify:` appears anywhere, so an empty `Verify:` line passes validation. Separately, the validator enforces at least one edge/error/offline AC or waiver, but it never enforces the binding requirement for at least one happy-path AC. A ticket containing only an error-path Gherkin AC passes. See `addons/parley-tracker/bin/validate.js:336` and `addons/parley-tracker/bin/validate.js:484`.

Why it matters: `FINAL.md` requires non-functional criteria to have a concrete single verification command and requires both happy-path and edge/error coverage. Passing empty verification or edge-only tickets undermines the no-assumption gate the tool is supposed to enforce before agents claim work.

Concrete fix: change `acHasVerify()` to require non-whitespace command text after the colon, preferably one shell command line. Track AC classification and require at least one non-edge happy-path AC plus at least one edge/error/offline AC unless there is an explicit `n/a (reason)` waiver. Add tests for `Verify:` with no command and for an edge-only ticket.

### [MAJOR] `validate` accepts tickets missing canonical schema and three-audience sections

What is wrong: the validator only requires `id`, `type`, `title`, `status`, and `parent`, plus separate checks for `files/apis/arch` and `dod`. It does not require several canonical frontmatter fields from the binding template, including `assignee`, `priority`, `labels`, `worktree`, `mirror-owned`, and `canonical_source`. It also does not require `## At a glance`, `## [B] Business`, `## [T] Technical`, `## Non-goals`, or `## Open questions`; only `## [A] Agent directives` is enforced. See `addons/parley-tracker/bin/validate.js:54`, `addons/parley-tracker/bin/validate.js:403`, `addons/parley-tracker/bin/validate.js:441`, and `addons/parley-tracker/bin/validate.js:452`.

Why it matters: the tool can return PASS for a skeletal agent-only ticket that is not readable by all three intended audiences and is not a complete tracker mirror. That is a direct gap against `FINAL.md` section B.2 and B.4.

Concrete fix: expand required frontmatter validation to the canonical template fields, validate `priority` as required rather than optional, require `mirror-owned` to include only supported mirror-owned fields, and require non-empty `At a glance`, `[B]`, `[T]`, `[A]`, `Non-goals`, and `Open questions` sections. Add a regression fixture that currently passes with only minimal fields and `[A]`/AC sections.

### [MINOR] Add-on command coverage misses the cross-command edge cases

What is wrong: the tests cover default install, `--no-addons` install, `--only` install, unknown `--only` through `run(["install", ...])`, default doctor after a default install, and auto-detection ignoring installer-owned add-ons. They do not cover `status`, `paths`, `doctor`, or `uninstall` with `--only` and `--no-addons`, nor the plain doctor/status behavior after a core-only or legacy install. See the new coverage around `test/installer.test.js:385` through `test/installer.test.js:495`.

Why it matters: these are exactly the command surfaces where selection semantics become ambiguous, and the core-only doctor false negative above escaped the current suite.

Concrete fix: add table-driven tests for every command that accepts addon-selection flags: default, `--no-addons`, `--only parley-tracker`, unknown `--only`, and legacy/core-only install state. Include assertions for both structured JSON/action shapes and human-readable output where output is part of the compatibility promise.

### [NIT] Package metadata still uses a personal GitHub namespace

What is wrong: `package.json` points `homepage`, `repository`, and `bugs` at `github.com/feci/...` (`package.json:6`, `package.json:9`, `package.json:12`). The implementation grep checked local names and model IDs, but not this personal namespace.

Why it matters: this is not a runtime or installer bug, but the brief explicitly calls out personal/genericity leakage. If this package is intended to be organization-neutral, the distribution metadata should not be the remaining personal anchor.

Concrete fix: either move these URLs to the intended neutral project/org namespace before release, or document that `feci` is the canonical package owner so reviewers do not keep treating it as a genericity leak.

## Open questions

- Should `validate <single-ticket.md>` be allowed to pass a non-`n/a` parent without proving it resolves, or should parent resolution be mandatory unless the caller explicitly opts out? The current strict mode resolves parents only when validating a tree.
