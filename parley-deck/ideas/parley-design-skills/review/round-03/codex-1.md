---
agent: codex-1
idea: parley-design-skills
review-round: 03
date: 2026-07-28
reviewed-commit: 17f6619
---

## Summary

❌ BLOCK. Fix-up cycle 4 genuinely closes every round-02 reproducer I published for AF-1 and AF-2, the empty-path crash, and the fast-path contradiction. It only partially closes alias direction, however, and two other binding inputs remain outside the obligations: DIRECTION token-sidecar identity and the canonical frontmatter lexical subset. Each gap has a current-tree reproducer that reaches `verifiedLevel: L3`, `verdict: PASS`, and process exit 0.

I accept D-1. The four binding files total 65,364 bytes, within the 65,536-byte aggregate cap, and the rebalanced per-file thresholds are documented and tested as early-warning limits. D-2 also remains an honest, non-conformance-affecting deviation.

Round-02 disposition, per finding:

- AF-1's seven published cases are closed: swapped rounds, invented axis, duplicate primary position, false G3 result, false G4 result, quoted-brace CSS evasion, and unknown artifact kind all fail verification now. AF-1's broader certificate-integrity concern is not closed because the new token-sidecar and frontmatter probes below still forge L3 certificates.
- AF-2 is genuinely closed: a counter-signer who authored the waived work is rejected and the finding remains active.
- The fast-path contradiction is genuinely closed in both `SKILL.md` and `PDS.md`: it is outside Parley Deck, carries no Parley verification, and cannot claim above L1.
- The empty `waivers`/`tokens` crash is genuinely closed: the former completes normally and the latter becomes UNJUDGEABLE/exit 4.
- Alias direction is not genuinely closed: upward and primitive-to-primitive edges are covered, but same-tier semantic and component edges are still accepted.

The overlapping MAJOR/CRITICAL round-01 findings from the other reviewers are resolved as follows:

- Hermes AF-7 (undefined banned-slop signature): genuinely closed. The derived T0 signature and sharing rules are specified, and the engine recomputes observed and recorded signatures.
- Kimi AF-8 (G1 remedy omitted C7 conjuncts): genuinely closed. G1 now carries the ban list, category-plus-avoidance requirements, and the ratification/brief-reason rule.
- Kimi AF-9 (U1 assignment unverifiable): its stated rotation/run-id defect is genuinely closed, including the invented-axis and duplicate-primary probes. The DIRECTION token-sidecar defect below is a separate uncovered part of the binding process mapping.

## What I verified (commands run, and their result)

- `git status --short --branch` reported `## parley-design-skills`; `git rev-parse HEAD` reported `17f6619a28851b79f7caec5739f9f91dc2fcd39e`. `git diff 17f6619^..17f6619 --check` produced no errors.
- `npm test` completed with 212 tests passed, 0 failed.
- I reran my round-02 fixtures with the current `check.js`. Results:
  - swapped rounds: exit 1, `pds-check:l2-process-order`;
  - invented axis: exit 1, both directions rejected for an axis absent from the brief;
  - duplicate primary position: exit 1, repeated primary and recomputed-assignment findings;
  - false G3 result: exit 1, the recorded pass refuted by the raw literal;
  - false G4 result: exit 1, the recorded pass refuted by open quality rules;
  - quoted-brace CSS evasion: exit 1, the raw color detected;
  - unknown `TYPO-BRIEF` kind: exit 1, `pds-check:l1-artifact-kind`;
  - waiver signed by the waived work's author: exit 1, waiver rejected and finding retained.
- Empty-path probes: `waivers: ""` completed at L2 with PASS/exit 0 and no crash; `tokens: ""` completed with UNJUDGEABLE/exit 4 and “DIRECTION names no tokens file,” also without a crash.
- A full sound-run check at L3 produced PASS/exit 0. I then used that run for three mutations:
  - semantic `medium -> {semantic.small}` produced L3 PASS/exit 0 and `l3-alias-direction: met`;
  - `x-note: unquoted # trailing comments are forbidden` in a DESIGN-BRIEF produced L3 PASS/exit 0 and `l1-frontmatter-parses: met`;
  - without mutation, both shipped DIRECTION files name the same `../tokens.json`, yet the run produces L3 PASS/exit 0.
- `node addons/parley-design-check/bin/check.js --help` exited 0 and documents exit 4. Checking `NOTICE.md` alone returned UNJUDGEABLE/exit 4.
- `npm pack --dry-run --json` exited 0 with 143 files. The package includes `NOTICE.md`, all four doctrine files, and the checker. The doctrine files total 65,364 bytes; `RULES.md` hashes to the digest recorded by `PDS.md`; the registry exposes 18 unique rule IDs backed by 18 detectors; placeholder scans were empty; and no second `RULES.md` ships under the checker.

## Findings

### [CRITICAL] L2 ignores the required per-direction token sidecars

The binding artifact map requires `round-01/<agent>.md` together with `round-01/<agent>.tokens.json`. The engine's process-home check validates only the DIRECTION markdown locations. Its own sound-run has both `round-01/claude-1.md` and `round-01/codex-1.md` pointing to the same root-level `../tokens.json`; the checker nevertheless certifies L3 with PASS/exit 0.

This matters because independent divergence includes each proposer's token proposal. A shared or misnamed sidecar lets a run replace two independently authored inputs with one file while receiving a process-conformance certificate.

Fix: define token-path resolution unambiguously; require every DIRECTION to resolve to an adjacent, uniquely owned `<agent>.tokens.json` under `round-01`; reject reused paths; correct the PDS example if paths are artifact-relative; and add shared, misnamed, missing, and cross-owned sidecars as negative fixtures.

### [MAJOR] The canonical frontmatter obligation accepts forbidden lexical syntax

PDS §2 permits `#` only inside quotes and permits comments only as whole lines, never trailing. The subset parser accepts a bare value containing `#` and trailing text. Adding `x-note: unquoted # trailing comments are forbidden` to an otherwise sound DESIGN-BRIEF leaves `l1-frontmatter-parses` met and the whole run certified at L3, PASS/exit 0.

This matters because L1 claims that frontmatter parses under the canonical subset, not merely under the checker's more permissive private grammar. Extension keys make this a general certificate bypass even where required fields remain conventional.

Fix: make scalar lexing reject unquoted `,`, `[`, `]`, `{`, `}`, and `#`, edge whitespace, trailing comments, malformed quote pairs, and escapes forbidden by §2; add negative fixtures for each forbidden form.

### [MAJOR] Alias direction still accepts same-tier edges

The cycle-4 condition rejects upward aliases and primitive-to-primitive aliases, but permits semantic-to-semantic and component-to-component aliases. A full sound run with `semantic.medium -> {semantic.small}` is certified at L3, PASS/exit 0, with the alias obligation reported met. PDS G3 requires references to point strictly down the declared primitive → semantic → component order.

This matters because the checker is claiming the exact L3 invariant while accepting edges that do not descend the tier order. It also leaves same-tier cycles structurally possible.

Fix: for every alias whose source and target have recognized tiers, require `targetTier < sourceTier`; reject equality at all three tiers, and add semantic-to-semantic and component-to-component negative fixtures.

## Open questions

None. The three fixes and regression cases are concrete.
