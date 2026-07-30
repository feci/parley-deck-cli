---
idea: integrate-parley-bidding-addon
review-round: 5
agent: kimi-1
date: 2026-07-30
reviewed-commit: 3634cc8
---

## Verdict

PENDING MEASUREMENT (file written first per protocol; code read complete, measurements running.
Provisional direction: ACCEPT, with one suspected MINOR on read-command `--only` semantics —
to be confirmed or withdrawn by measurement below.)

## Round-4 findings — closed or not

My round 4 ended as a PENDING draft; its two named suspects were taken into cycle 4 anyway:

1. **Probe cache key collision** (`\u0001` separator can occur inside env values) — claimed fixed:
   key is now `JSON.stringify` of sorted `[name, value]` pairs over the whole effective
   environment (lib/installer.js:1451). Closure by measurement: PENDING.
2. **`managed` absent on `missing` units** — claimed fixed: the missing shape now carries
   `managed: false` (lib/installer.js:1295). Closure by measurement: PENDING.

My third named suspect (install/uninstall/`--only`/`--no-addons` interactions with the verdict)
overlaps codex-1's F9–F11; those fixes are in the diff and will be re-measured here.

## New findings

PENDING. Suspects under investigation:

- `doctor --only <subset>` on a healthy universal install: the flag overrides the recorded
  selection for read commands (expectedAddonNames, lib/installer.js:864), so the new traversal
  may mark recorded-and-installed add-ons `selected: false` and fail health on a consistent
  tree, with a problem message ("not part of the recorded selection") that is factually wrong
  in that case — they ARE in the recorded selection, just outside the flag.
- `selected: false` traversal vs `--dest` (generic target), project scope, and same-named
  foreign directories in the skills dir.
- Stricter ownership predicate (installerOwnsDestination, lib/installer.js:1610): does any
  previously-legitimate flow now refuse? Checked so far by reading: 2.0.0-era markers carry
  `name` + `skill`, so upgrades remain owned — to be confirmed by measurement.

## What I verified and found correct

PENDING.

## Open questions for the implementer

PENDING.
