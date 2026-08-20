---
idea: addon-manifest-coverage
author: claude-1
created: 2026-08-01
participants: [claude-1, codex-1, hermes-1, kimi-1]
track: deliberation
strict_gate: false
status: final
---

## Problem / idea

`doctor` exits 1 and reports five of six skills `malformed` whenever the skills were installed
by anything other than this package's own installer — including the universal
[`vercel-labs/skills`](https://github.com/vercel-labs/skills) CLI, which `README.md` recommends
**first**.

Reproduced on 2026-08-01 against `parley-deck-skill` at repo version 2.1.0
(`HEAD = 23a9856`), in an isolated `HOME`, by copying every packaged skill verbatim with no
install marker — exactly what a foreign installer leaves behind:

```bash
REPO=<repo>/parley-deck-skill
SB=<tmpdir>; mkdir -p "$SB/home/.codex/skills"
for d in "$REPO"/skills/*/; do cp -R "$d" "$SB/home/.codex/skills/$(basename "$d")"; done
cd "$REPO" && HOME="$SB/home" node bin/parley-deck-skill.js doctor --target codex --json
```

Result — `exit=1`, `ok: false`, target status `malformed`:

| skill                 | status            | managed |
| --------------------- | ----------------- | ------- |
| `parley-deck` (core)  | `malformed`       | false   |
| `parley-bidding`      | `valid-unmanaged` | false   |
| `parley-design`       | `malformed`       | false   |
| `parley-design-check` | `malformed`       | false   |
| `parley-tracker`      | `malformed`       | false   |
| `parley-worktrees`    | `malformed`       | false   |

Every `malformed` unit carries exactly one problem:

```
no parley-deck-skill install marker: this directory was not installed by this tool,
or the marker was removed
```

**Mechanism.** In `lib/installer.js`, an expected unit with no install marker is not
immediately malformed — it is routed through `unmanagedButVerified(unit)`, which returns true
only when the unit can be *proven* intact from its packaged source manifest:

```js
function unmanagedButVerified(unit) {
  const source = unit.addon ? unit.addon.root : null;
  if (!source || !addonManifest.hasManifest(source)) {
    return false;
  }
  ...
}
```

Two independent reasons it returns false for five of six units:

1. **Only `parley-bidding` ships a `parley-addon.json`.** The other four add-ons have no
   manifest, so there is nothing to verify against.
2. **The core `parley-deck` unit is not an add-on at all**, so `unit.addon` is null and the
   source is null before the manifest is even consulted. Shipping four more add-on manifests
   would still leave the core skill `malformed`.

The payload bytes are correct in every one of these cases. The tool is reporting *absence of
evidence* as *evidence of defect*, and it does so on the install path documented first.

This was recorded as a deferred follow-up by the `integrate-parley-bidding-addon` review
consensus rather than fixed there, because widening manifest coverage was out of that idea's
ratified scope. `CHANGELOG.md` states it as a known limit. This idea is the follow-up.

## Constraints

- **The B3 invariant must not regress.** A tree gutted down to `SKILL.md` — with or without a
  marker — must still be reported unhealthy. The gutted-tree false green is precisely what the
  marker requirement and the manifest were introduced to close; any proposal that makes
  unmarked trees pass more easily must show, with a regression that fails against the current
  commit, that the gutted tree still fails.
- **Proof stays anchored to the packaged source, not to the installed tree.** Verifying an
  installed payload against whichever manifest happens to sit beside it recognises any
  self-consistent tree. Round 4 of the prior idea established this after a deleted `runtime`
  field silently disabled the interpreter check without rehashing a file.
- **A stale manifest must not be able to turn a correct tree malformed.** If manifests become
  generated artifacts, drift between the manifest and the payload is a new failure mode that
  is strictly worse than today's — it would fire on the installer's own happy path. Say how
  drift is prevented and how it is detected.
- `parley-addon.json` is **defect detection, not tamper resistance**. Do not argue for or
  against a design on the basis of an attacker who can already write to the payload.
- Manifest semantics currently in force: schema `parley-addon/1`; POSIX-relative keys with no
  `..`, no absolute paths, no drive letters, no backslashes; every declared entry must be a
  regular non-symlink file; the manifest excludes itself and the install marker; the aggregate
  digest covers the file list, so renames and additions change it.
- Whatever is decided must hold for all four documented install paths, not only `doctor`:
  install, uninstall, `status`, and `doctor`. One destination must not be simultaneously
  healthy-and-owned by one command and unowned by another.
- English-only for all files under `parley-deck/`. Standard release discipline applies once
  implemented: version bump plus all five channels.

## Non-goals

- Concurrent-installer isolation (locking two installer processes against each other). That is
  a separate recorded follow-up and is out of scope here.
- Changing what the universal `skills` CLI does. This package cannot make a foreign installer
  write our marker; the fix has to live on our side of the boundary.
- Tamper resistance, signing, or any trust model beyond defect detection.
- Reworking the `valid-unselected` / `valid-unmanaged` status vocabulary beyond what this fix
  strictly requires.

## Questions the rounds should answer

1. Should the fix be **manifest coverage** (ship `parley-addon.json` for all six units,
   including the core skill), a **semantic change** to how an unmarked-but-correct tree is
   classified, an **exit-code policy** change, or some combination? State the failure mode each
   option leaves open.
2. If manifests are shipped for all six: are they committed artifacts or generated at pack
   time, and what prevents a stale one from failing a correct install?
3. `unmanagedButVerified` reaches the source through `unit.addon.root`. What is the correct
   source anchor for the core skill, and does that generalise without naming any skill?
4. Is `valid-unmanaged` the right terminal status for a foreign-installed but provably intact
   tree, and should `doctor` exit 0 when every unit is at worst `valid-unmanaged`?
5. What regression proves the fix, and what regression proves the B3 invariant survived it?
