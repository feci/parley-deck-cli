---
agent: codex-1
idea: sync-skill-protocol-fallback
review-round: 1
date: 2026-06-24
reviewed-commit: ba97ae859e95f92d7131aca846c27b400bbeeb8e
---

## Summary

No findings. I tried to break the sync against the agreed rule: the skill fallback must be the CLI embedded default verbatim after the header, with only the Transport and Created lines neutralized. The committed sibling repo is clean at `ba97ae859e95f92d7131aca846c27b400bbeeb8e`, and the current `references/COOPERATION.md` matches that commit.

Declining the deeper §0/§9.0 `parley` token scrub was the right call for this change. Shipping body text that names `parley init`, `parley preflight`, `parley run`, `parley retro`, and `~/.parley` is not a defect here because the accepted artifact is a protocol fallback for the Parley Deck ecosystem, not a newly generalized non-Parley protocol. Scrubbing those literals would fork canonical body wording and create a harder-to-detect protocol drift risk. The neutral header fixes the active-project/provenance leak; the body literals remain intentional ecosystem/tooling references.

## Refutation attempts

- Read the required design artifacts: `FINAL.md`, `consensus.md`, and `IMPLEMENTATION.md`; also read the live `parley-deck/COOPERATION.md` review rules before writing this Phase 6 artifact.
- Inspected `git -C ../parley-deck-skill show ba97ae8`; the commit changes only `references/COOPERATION.md`, `package.json`, `package-lock.json`, and `references/compatibility.json`.
- Checked sibling repo state: `git -C ../parley-deck-skill rev-parse HEAD` is `ba97ae859e95f92d7131aca846c27b400bbeeb8e`, and `git status --short` / `git diff --stat` are empty.
- Body identity: `diff <(tail -n +7 internal/protocol/defaults/COOPERATION.md) <(tail -n +7 ../parley-deck-skill/references/COOPERATION.md)` produced no output.
- Full-file diff: `diff -u internal/protocol/defaults/COOPERATION.md ../parley-deck-skill/references/COOPERATION.md` showed exactly one header hunk, replacing only `**Transport:** \`github-pr\`` and `**Created:** \`<date> — created by parley init\`` with the agreed neutral lines.
- Header neutrality: `**Workspace:** \`<workspace-name>\`` appears once; the old `**Workspace:** \`parley-deck\`` literal appears zero times; the agreed Transport and Created fallback lines each appear once.
- Section and content anchors: `^## 13\.` count is 1 and `^## 14\.` count is 1. I also located the expected §4 Phase 6 refutation text, Phase 8 strict gate, loop budgets, close-decision integrity, §9.0 readiness check, and §12.11 `status: candidate` remediation text in the shipped fallback.
- Leak check: `grep -Eic 'feci|claude-1|codex-1|hermes-1|antigravity-1' ../parley-deck-skill/references/COOPERATION.md` returned `0`.
- Version consistency: `package.json` is `1.4.1`, `package-lock.json` top-level version is `1.4.1`, `package-lock.json` root package version is `1.4.1`, and `references/compatibility.json` `skillVersion` is `1.4.1`.
- Release preflight spot checks: `npm test` passed with 77 tests and 0 failures; `node bin/parley-deck-skill.js install --target all --dry-run` exited 0; `node bin/parley-deck-skill.js doctor --target all --json` reported `"ok": true`. `npm pack --dry-run` initially failed on this machine's root-owned `~/.npm` cache, then passed when rerun with `npm_config_cache=/tmp/parley-deck-skill-npm-cache`, producing `parley-deck-skill-1.4.1.tgz` with 22 files including the 79.1 kB `references/COOPERATION.md`.

## Findings

None.

## Open questions

None.
