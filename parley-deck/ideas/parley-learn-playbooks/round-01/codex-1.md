---
agent: codex-1
idea: parley-learn-playbooks
round: 1
date: 2026-07-04
---

## Summary

`parley learn <closed-idea-slug>` should be a small sibling to `parley retro`, not a new phase runner. It should be wired as a top-level command in `internal/app/app.go`, read one closed idea, and write one advisory artifact at `parley-deck/playbooks/<topic>.md`.

The closest implementation pattern is today's `parley retro`: `app.go` dispatches to `runRetro`, `internal/app/retro.go` owns CLI flags and the narrow write boundary, and `internal/retro` owns read-only artifact mining. `retro propose` is especially relevant because it validates strict slugs, rejects existing targets, uses `Lstat` to avoid symlink redirects, and tests that only the allowed file is written.

## Proposed approach

Add:

```text
parley learn [--dir DIR] [--topic TOPIC] <closed-idea-slug>
```

`--topic` should be optional but strongly useful. If omitted, v1 can derive the topic from the source idea slug, but the output filename should still pass a strict lowercase kebab-case rule. I would not add overwrite/update behavior in v1; if `parley-deck/playbooks/<topic>.md` already exists, fail closed and let a normal idea review any substantive revision.

Implementation shape:

- Add `case "learn": return runLearn(args[1:], stdout, stderr)` in `internal/app/app.go`, plus usage/help text near `retro`.
- Add `internal/app/learn.go` for flag parsing, absolute `--dir`, user-facing errors, and safe file creation.
- Add `internal/learn/learn.go` for read-only extraction and deterministic markdown rendering.
- Share or duplicate the existing strict slug validation from `retro.go`; do not accept path separators, dot paths, uppercase, spaces, or double hyphens.
- Resolve `parley-deck/ideas/<slug>` with `Lstat`; reject missing, symlinked, or non-directory source paths.
- Define "closed" conservatively: accept `IMPLEMENTATION.md` frontmatter `status: complete`, or design-only `00-prompt.md` status `final` with `FINAL.md` present. Anything else should print that the idea is not closed enough to learn from.
- Read only canonical artifacts: `00-prompt.md`, round directories, `consensus.md`, `FINAL.md`, `IMPLEMENTATION.md`, review rounds, review consensus, and run event logs if needed. Do not read raw transcripts.
- Render a playbook with frontmatter like `playbook`, `source_idea`, `created`, `status: advisory`, followed by sections for applicability, phase pattern, checklist, verification evidence, gotchas, and details deliberately not copied.
- Create `parley-deck/playbooks/` with the same defensive posture as `retro propose`: reject a symlinked or non-directory parent, then write `playbooks/<topic>.md` with `O_CREATE|O_EXCL`.
- Print the path written and an advisory reminder: the playbook is optional context and never overrides `COOPERATION.md`.

If the team wants model-quality prose rather than deterministic extraction, I would keep that as a second slice. It can reuse the existing consult/runner machinery with an explicit `--agent AGENT` option, but the default v1 should stay deterministic and reviewable. A model-assisted mode should still validate the final path and write exactly one playbook file.

For usage in new ideas, avoid implicit injection. A later explicit flag such as `parley run --playbook TOPIC TASK` can copy a short reference into `00-prompt.md` and runner prompts. Automatic "brief matches playbook" behavior should only suggest candidates, not mutate prompts without facilitator choice.

Tests should mirror `internal/app/retro_test.go`: success writes only `playbooks/<topic>.md`; invalid source slugs and topics are rejected; open ideas are rejected; existing playbooks are not overwritten; symlinked source/output paths are rejected; and the generated markdown includes advisory status, source idea, and the expected checklist sections.

## Concerns / open questions

The biggest open question is whether `parley learn <slug>` is expected to produce a high-quality generalized playbook by itself. Go can deterministically extract structure and evidence, but "strip idea-specific details" is a semantic task. I would rather ship a conservative deterministic draft than hide an implicit LLM call behind a command that looks local and predictable.

Topic naming also needs a decision. A source slug is often too specific for a reusable playbook, while an inferred topic can collide or be vague. Optional `--topic` gives operators control without breaking the simple command form.

The protocol says playbooks are advisory and maintained through normal ideas when substantively revised. The CLI should therefore refuse silent updates and avoid any auto-application path in `parley run`.

## Risks

The main implementation risk is over-writing. `learn` must not become a broad generator that creates prompts, round files, or protocol edits. The write boundary should be one playbook file, with tests proving that.

There is also a leakage risk: closed ideas can contain project-specific details or accidental secrets. V1 should read only canonical artifacts, include a "details not copied" section, and rely on normal review before the playbook is treated as useful.

Finally, if playbook references are injected automatically into future agent prompts, they could quietly bias new ideas or appear protocol-like. Keep playbooks explicitly advisory, explicitly selected, and subordinate to `COOPERATION.md`.
