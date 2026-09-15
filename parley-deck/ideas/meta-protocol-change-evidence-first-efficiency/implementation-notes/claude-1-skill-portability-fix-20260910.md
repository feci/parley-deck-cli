---
agent: claude-1
idea: meta-protocol-change-evidence-first-efficiency
date: 2026-09-10
artifact-kind: owner change record (bounded wording + test correction)
not-a-signoff: true
not-a-phase-6-review: true
not-an-acceptance: true
owner-files-changed:
  - <skill-worktree>/skills/parley-deck/SKILL.md (Required Protocol Context: step 3 + closing paragraph)
  - <skill-worktree>/test/packet-context.test.js (two tests' assertions)
---

# Bounded follow-up: bundled snapshot is orientation only; test wording dependency

## Change 1 — SKILL.md, Required Protocol Context

Step 3 previously read as an alternative load path ("load the bundled fallback
snapshot ... and record `context_mode=full-fallback` with reason `bundled-snapshot`"),
which sat ambiguously beside FINAL D4. It now states the snapshot is **LOCAL
ORIENTATION ONLY**: it cannot authorize a protocol task launch, it is not a launch
`context_mode`, and when live authority cannot be established the task stays blocked
(report the blocker instead of launching participants or recording the launch as
attested). The closing paragraph of the section was aligned to the same wording.

Untouched: step 1's explicit `refused` stop rule, step 2's no-CLI path (live file read
in full, `full-fallback` + reason), step 4's stop-and-ask, and the shasum drift check.

## Change 2 — test/packet-context.test.js

The packaged-protocol test matched `with the reason`, while §9 item 1 says `with its
reason`. Replaced with `with (?:the|its) reason` plus semantic assertions that
full-fallback is the disclosed read of the live authority, and that `refused` is a
**stop** and never a substitution onto other authority text. A stale skill-side
assertion (`` `refused` launch ... does not proceed with substituted protocol text ``)
matched no current wording; it was re-expressed against the shipped stop-rule text
rather than dropped. Added coverage for the new step-3 semantics. No assertion removed;
no other test touched.

## Limitations

No shell in this launch: I ran no tests and make no pass claim — the full-suite result
is the facilitator's to check. The add-on manifest still needs regeneration with the
existing script; I invented no hashes. No FINAL, signoff, or acceptance judgment here.
