---
from: claude-1
to: user
idea: integrate-parley-bidding-addon
phase: implementation
blocking: yes
date: 2026-07-30
---

## Question

`FINAL.md` B3 asks for two things that cannot both be true. Which do you want?

> **B3.11** — "Generic and optional: add-ons without a manifest keep `SKILL.md`-only
> compatibility, so `parley-design`, `parley-design-check`, `parley-tracker` and
> `parley-worktrees` are unaffected."
>
> **B3.13** — "required for `parley-bidding`."

If the manifest is required *for one named add-on*, `lib/installer.js` needs a name-keyed rule —
which re-introduces exactly the hardcoded add-on registry that `discoverAddons`
(`lib/installer.js:778-793`) exists to avoid, and which the 2.0.0 layout move was partly about
removing.

If instead the rule is "required if the add-on ships one" (presence-keyed, genuinely generic),
then **B3.13's own acceptance test becomes unimplementable**: deleting `parley-addon.json`
would silently downgrade `parley-bidding` to `SKILL.md`-only validity and report healthy — the
gutted tree B3 exists to catch.

No measurement settles this. It is a design choice about where the requirement is anchored.

## Context

Everything else in B3 is measured and implementable. The manifest itself, the hashes, the
aggregate digest, and the negative tests (delete `adapter_validate.py`, delete a schema, delete
`references/hitl-and-recovery.md`, flip one byte → `malformed`) all work under either answer.
Only the *anchor* differs.

Today `validateInstalledPayload` (`lib/installer.js:1129-1148`) requires exactly `["SKILL.md"]`
for `kind === "addon"` (line 1132). That single line is why a directory containing nothing but a
`SKILL.md` currently reports `valid` — the defect B3 was written against.

A third option exists that the contract does not mention: anchor the requirement in the
**marker** the installer already writes at install time. The installed marker records which
add-ons were installed and from which package version; it could also record "this add-on
shipped a manifest", making later deletion detectable without naming any add-on in code.

## What I need from you

Pick one:

1. **Name-keyed** — `installer.js` knows `parley-bidding` requires a manifest. Honest to B3.13,
   costs the no-registry property.
2. **Presence-keyed** — any add-on shipping a manifest is validated against it. Keeps the code
   generic; B3.13's deletion test is dropped and recorded as a deferred follow-up.
3. **Marker-anchored** — the installed marker remembers that a manifest was present at install
   time, so deleting it afterwards is detectable and no add-on is named in code. Closes both
   B3.11 and B3.13, but it is a design the reviewers did not see and did not ratify, so it would
   need their signoff before it counts as satisfying B3.

I will not choose this myself: it changes what the ratified contract means, and option 3 is not
in the contract at all.
