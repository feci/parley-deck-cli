---
from: zcode-1
to: user
idea: meta-protocol-change-lean-organizer
phase: implementation
blocking: no
date: 2026-09-23
---

# Owner-only action: attended publish of the staged combined core 2.13.0

Per FINAL (cross-cutting decisions) and the owner's 2026-09-23 core-publish decision
(`inbox/user-to-codex-1_meta-protocol-change-lean-organizer_core-publish.md`: "počkám,
publikuj to celé na konci" — one combined core at the end), the staged core is ready:

**Staged file:** `~/.parley/staging/COOPERATION-2.13.0.md` (109,772 B, sha256
`fc907e5914a072d1a6afe249fc39401e1f8761cc1d67f2ce002dfde210762c9f` — restaged
2026-09-24 by fix-up F20's §9.0 sentence; verify before publishing:
`shasum -a 256 ~/.parley/staging/COOPERATION-2.13.0.md`)

**Composition (verified 2026-09-24):** built from the core TEMPLATE (placeholder
header, stub §2 — never a deck view copy) carrying exactly the three pending change
sets on top of published core 2.10.0, and nothing else:

1. the never-published 2.11.0 hunks (§15.6 alternatives + correlated agreement, §15.7 row),
2. the 1.48.0 protocol hunks (§9 item 1 attestation rule, goal-check "can only withhold
   a close" paragraph, LE-7/LE-11 bullet),
3. this idea's reviewed hunks (lean-organizer permissive lines: Quickstart facilitator
   row, §4 Phase 5/Phase 6, §9.0 pure-organizer default sentence, §9 checklist item 1
   audience view + brief re-orientation, §11 one-blocking-wait).

Verification performed (re-verified 2026-09-24 at the fix-up-1 HEAD by the
implementer, independently re-verified by both round-02 reviewers in their own
signoff blocks): placeholder header intact; §2 host-handle table stubbed; the
§15 region and every changed section byte-match the released deck view at this branch;
`diff` core 2.10.0 → staged core contains exactly the three change sets (11 hunks,
each attributed to exactly one change set, 0 unexplained).

**The exact owner-only command** (attended; refuses without a controlling terminal —
no participant or organizer may work around the TTY gate):

```
parley protocol publish --version 2.13.0 --from ~/.parley/staging/COOPERATION-2.13.0.md
```

This is the owner's action at the release step (after Phase 8 completes). No release,
tag, channel publish, or main-branch integration is performed by participants before
the organizer's release step.
