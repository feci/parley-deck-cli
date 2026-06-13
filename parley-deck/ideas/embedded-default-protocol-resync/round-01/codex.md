---
agent: codex
idea: embedded-default-protocol-resync
round: 1
date: 2026-06-13
---

## Summary

The embedded default should be treated as the bootstrap template rendered by `parley init`, not as a literal project deck mirror. The immediate protocol drift is straightforward: add Section 12 from the live deck and keep its provenance sentence, while not copying this project's `Protocol synced:` header into bootstrap output. The durable fix is a Go test that compares the embedded template and the live deck after a narrow, documented normalization for intentional bootstrap/project differences.

## Proposed approach

1. Section 12 propagation.

- Append the live `## 12. Pipeline blocks & action stages` section to `internal/protocol/defaults/COOPERATION.md` after the current Section 11, preserving exact text and the final newline.
- Keep the closing provenance line about the ratifying idea. It is protocol history, not runtime project state; dropping it would create a needless special case and weaken auditability.
- Do not add the live `Protocol synced:` header to the embedded default or initialized decks. That line records this repository's project sync state.

2. Bootstrap genericization.

- `Parley deck: ./parley-deck/`: carry verbatim. It is the fixed directory convention used by `DeckDir`.
- `Transport:`: make initialized output explicitly `local-dir`. Prefer making the checked-in bootstrap template say `local-dir` and removing the blind `github-pr` to `local-dir` replace; if consensus keeps the current source line for easier mirroring, cover the transform with a focused test.
- `Workspace:`: project-specific. `InitWorkspace(root)` should fill it from `filepath.Base(root)` after simple display sanitization, or fall back to `<workspace-name>` only when no sane basename exists.
- `Created:`: project-specific. Initialized decks should use the init date with wording like `(created by parley init)`. Keep this testable by passing the date through a helper or parameter rather than burying `time.Now()` in string replacement.
- `Protocol synced:`: omit. It is valid in the live parley-deck project deck, not in generic bootstrap output.
- Section 2 roster rows: project-specific. Remove the hard-coded `codex`, `claude`, `agy`, and `hermes` rows from the embedded template. A new project should not start with false quorum members. If `init` cannot discover agents, leave the roster table structurally present but empty.
- Section 2 host-handle rows: project-specific. Remove the hard-coded agent rows there as well. Empty is better than four `not mapped` entries for agents that are not in the project.
- Other `parley-deck` strings: keep when they name the protocol/product, directory layout, `parley-deck/` paths, or the coordination repo role in transport mechanics. Those are protocol-general terms, not this repository's workspace identity.
- Existing placeholders and examples such as `<idea-slug>`, `<agent-id>`, and `execution-worker-retry-policy` can remain; they are visibly illustrative.

3. Anti-drift guard to ship in this idea.

- Add a Go test in `internal/protocol` that loads `internal/protocol/defaults/COOPERATION.md` and `parley-deck/COOPERATION.md`.
- Compare normalized text, with the allowlist implemented in one helper and named in the failure message. The only allowed differences should be the live-only `Protocol synced:` line, `Workspace:` and `Created:` header values, the `Transport:` bootstrap default if the checked-in template uses `local-dir`, the Section 2 roster table body, and the Section 2 host-handle table body.
- Normalize by exact anchors and table boundaries, not broad regular expressions. If a heading or table separator is missing, the test should fail so the allowlist cannot mask real protocol drift.
- Add a second small test for `defaultCooperationForInit` output: it emits `Transport: local-dir`, contains the derived workspace or placeholder, does not contain this repository's roster rows, and includes Section 12.
- Do not ship a generator in this idea. It may be reasonable later, but it adds a write path and review burden before the simple invariant has proven noisy.
- Do not ship a `parley protocol check` subcommand in this idea. The Go test catches repository drift in CI, which is where this specific failure should have been caught. A subcommand can later wrap the same normalizer for users.
- A Section 7 manual checklist can be added later, but it is not a substitute for the test.

## Concerns / open questions

- Empty roster tables are correct for bootstrap but reduce "ready to run" feel. The better UX is a future `parley init --agent` or discovery-assisted init, not shipping this project's agents as defaults.
- If maintainers want a package provenance marker in initialized files, it should be a distinct template metadata line, not the live project's `Protocol synced:` line. I would leave it out until there is a concrete consumer.
- Deriving `Workspace:` and `Created:` changes the `defaultCooperationForInit` path from a single replace into a small renderer. That is acceptable, but it should be explicit and covered by tests.
- The drift test depends on this repository's `parley-deck/COOPERATION.md`. That is acceptable for this repo's CI; if the canonical deck is missing, the test should fail with a clear message rather than silently skip.

## Risks

- A too-broad normalizer would recreate the silent-drift problem under a test-shaped mask. Keep the allowlist narrow and anchored.
- Genericizing the roster could surprise anyone relying on `parley init` to preload the parley-deck development roster. That reliance is already wrong for new projects, so the compatibility break is acceptable.
- Dynamic `Created:` dates can make tests flaky if not injectable. Avoid hidden clock dependencies in the renderer.
- Copying Section 12 without the guard fixes today's diff but leaves the process vulnerable to the next ratified protocol addition.
