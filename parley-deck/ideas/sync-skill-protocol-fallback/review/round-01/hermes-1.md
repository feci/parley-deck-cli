---
agent: hermes-1
idea: sync-skill-protocol-fallback
review-round: 1
date: 2026-06-24
---

## Summary

The sync is correct. I assumed it was wrong and tried to break it across all five
agreed gates; every attempt failed. The shipped skill fallback
(`parley-deck-skill/references/COOPERATION.md` @ ba97ae8) is a byte-identical body
copy of the CLI embedded default (`parley-deck-cli/internal/protocol/defaults/COOPERATION.md`)
with exactly the two neutralized header lines (Transport, Created), the Workspace
literal scrubbed from the stale `parley-deck` to `<workspace-name>`, and version
1.4.1 bumped consistently in package.json + package-lock.json (root + packages[""])
+ references/compatibility.json skillVersion. No protocol rule was invented, dropped,
or altered; no roster/personal identity leaked.

Implementation commit ba97ae8 touches only the four expected files
(COOPERATION.md, package.json, package-lock.json, compatibility.json) and the
diff hunk for COOPERATION.md is exactly the header swap plus the verbatim body
re-import. IMPLEMENTATION.md's self-reported line count (759 -> 1037) matches
`wc -l` on both repos.

## Refutation attempts

(1) Body identity — Ran `diff <(tail -n +7 internal/protocol/defaults/COOPERATION.md)
<(tail -n +7 ../parley-deck-skill/references/COOPERATION.md)`. EXIT 0, EMPTY. The
entire protocol body (line 7 = `**Status:**` onward, 1031 lines) is byte-identical.
No rule invented, dropped, or altered. NOT BROKEN.

(2) Full-file diff — Ran `diff -u`. The only hunk is the header:
  - `**Transport:** `github-pr``  ->  `<transport-choice>` (pick one of local-dir | github-pr | gitlab-mr at deck bootstrap — see §0)`
  - `**Created:** `<date> — created by parley init``  ->  `<YYYY-MM-DD>` (set at deck bootstrap)`
Workspace is `<workspace-name>` in BOTH files (the CLI default already used it;
the commit additionally scrubs the old skill's stale `parley-deck` literal to
`<workspace-name>`). No other full-file hunk. NOT BROKEN. The `parley init` token
that lived in the CLI header's Created line is gone in the skill (accounting for
the 63 vs 62 `parley` literal count — it is a HEADER token removed by the
neutralization, not a body scrub).

(3) Section/anchor presence —
  - `^## 13.` = 1; `^## 14.` = 1. Each present exactly once.
  - `Refutation attempts` = 2 (the §4 Phase-6 review-file template + LE-1 refutation-default paragraph). `Refutation-default` = 1.
  - `strict_gate` = 7. `status: candidate` = 4.
  - `Loop budgets (LE-5)` present at L510; `Close-decision integrity (LE-7/LE-11)` at L519; `fix-up budget` at L488/L494.
  - §9.0 `### 9.0 Pre-idea readiness check` present and full (preflight, freshness sync, roster liveness ping, quorum lock).
  - §12.11 `### 12.11 Monitoring loop-closure` present and full (candidate status, non-solo invariant).
  - §13 Retrospective optimization full (13.1-13.4, `parley retro` tooling, RHO provenance).
  - §14 Automated outer loop full (14.1-14.3, human brake, `parley loop tick`, `status: candidate`).
  All anchors also exist at the same line numbers in the CLI default (970/1000),
  proving the content is canonical, not invented by the implementer. NOT BROKEN.

(4) Roster/personal leak — `grep -in 'feci|claude-1|codex-1|hermes-1|antigravity-1'`
on references/COOPERATION.md → EXIT 1, ZERO matches. No project roster, no
personal identities, no agent IDs. NOT BROKEN.

(5) Version consistency —
  - package.json L3: `"version": "1.4.1"`.
  - package-lock.json root L3: `"version": "1.4.1"`; packages[""] L9: `"version": "1.4.1"`.
  - references/compatibility.json L4: `"skillVersion": "1.4.1"`.
  (package-lock jumped 1.3.1 -> 1.4.1, package.json 1.4.0 -> 1.4.1, compatibility.json
  skillVersion -> 1.4.1 — all via npm version + the compat edit in the same commit.)
  NOT BROKEN.

I also confirmed both files are 1037 lines and that the commit stat
(4 files, +292/-14) is consistent with a header-only swap + full body re-import
on top of a 759-line predecessor.

## Findings

No CRITICAL, MAJOR, MINOR, or NIT findings. The implementation matches FINAL.md
and consensus.md with zero deviations. Body verbatim, neutral header, patch bump,
no leak, version consistent. I could not break it.

## Open questions

OQ1 — parley-token neutrality in §0/§9.0 (my own round-01 reservation, re-examined).
The shipped fallback still carries 62 `parley` literals in the body, including
`parley init` (§0 bootstrap), `parley preflight` (§9.0), `parley run`, `parley loop
tick`, `parley retro`, `parley-deck-skill status`, and `~/.parley/agents.toml`.
For a vendor-neutral fallback that may land with an agent lacking the `parley` CLI,
these read as parley-feature references. I judged this in round-01 and again here:
declining the deeper scrub WAS the right call for THIS idea. Reasons, all still
valid against the shipped artifact:
  (a) Over-scrubbing forks canonical wording — any token substitution inside §0/§9.0
      risks dropping or blunting safety semantics (the quorum-lock, the highest-by-default
      reasoning rule, the consumer auto-sync allowlist). The body-identity diff proves
      the shipped body is the canonical body; a scrubbed body would NOT be byte-identical
      and would need its own review ratification.
  (b) A half-measure (substituting tokens but leaving the surrounding prose) leaves
      §0/§9.0 still reading as parley-feature docs — strictly worse than either
      verbatim or a full tool-agnostic summarize. Summarize is explicitly out of scope
      (consensus deferred it to a follow-up idea).
  (c) The skill's own SKILL.md and tooling already reference the parley-deck-skill /
      parley ecosystem, so body `parley` mentions read as ecosystem tooling names, not
      preconditions that gate a non-parley agent out of the protocol.

So shipping `parley init`/`parley preflight` into a non-parley fallback is NOT a
real defect at this patch level: it is a documented, faithfulness-over-neutrality
trade-off, logged in consensus.md and reserved by me at signoff, with the full
summarize deferred to its own follow-up idea. It becomes a real defect only if the
fallback is later advertised as a standalone tool-agnostic protocol without the
accompanying summarize — which is out of scope here. I recommend the deferred
follow-up (§0/§9.0 summarize + a cross-repo drift test) remain tracked; both are
already recorded as out-of-scope follow-ups in FINAL.md and consensus.md.
