---
agent: kimi-1
idea: protocol-and-skill-audit
round: 1
date: 2026-08-20
---

## Findings

All sandbox work was done under `/tmp/pd-audit-kimi1` with `HOME` redirected; the shared working
tree was never modified. `SKILL` below = `/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-skill`
(the sibling npm package, version 2.9.0). Live runs used only read-only commands (`doctor`,
`status`, `paths`, `--check`); all mutations ran against the sandbox.

### F1 — `doctor` does not byte-verify the **managed core skill**, so a modified installed `COOPERATION.md` reports `valid` and exit 0, contra README's "verify any of them byte for byte"

severity: MAJOR
tag: PRIMARY (commands run; code citations SECONDARY)

command:
```
HOME=/tmp/pd-audit-kimi1/home node $SKILL/bin/parley-deck-skill.js install --target kimi   # exit 0, 6 skills
# flip ONE byte at the midpoint of the installed managed core's references/COOPERATION.md (length preserved)
python3 - <flip one byte> /tmp/pd-audit-kimi1/home/.kimi-code/skills/parley-deck/references/COOPERATION.md
HOME=/tmp/pd-audit-kimi1/home node $SKILL/bin/parley-deck-skill.js doctor --target kimi
```

output:
```
kimi/parley-deck: valid /tmp/pd-audit-kimi1/home/.kimi-code/skills/parley-deck
kimi/parley-bidding: valid ...
kimi/parley-design: valid ...
kimi/parley-design-check: valid ...
kimi/parley-tracker: valid ...
kimi/parley-worktrees: valid ...
doctor exit after core byte-flip: 0
```

Controls, same sandbox, proving the asymmetry is specific to the managed core:

- Delete `plugin.json` from the managed core → `kimi/parley-deck: malformed`, `missing: plugin.json`, exit 1. (Presence IS checked.)
- Flip one byte in managed add-on `parley-design/SKILL.md` → `kimi/parley-design: malformed`, `integrity: modified: SKILL.md`, exit 1. (Managed add-ons ARE byte-verified.)
- `cp -R skills/parley-deck` as a markerless foreign copy → `valid-unmanaged`; flip one byte in its `references/COOPERATION.md` → `malformed`, exit 1. (Foreign cores ARE byte-verified.)

contradicts: `parley-deck-skill/README.md:79-82` — "**Every** packaged skill ships a
`parley-addon.json` integrity manifest, so `doctor` can verify any of them byte for byte —
including when another installer put them there". The qualifier presents foreign installs as the
*additional* covered case; in reality the foreign core is the *only* core case that is
byte-verified. Root cause (SECONDARY, read): `lib/installer.js:148-154` `PAYLOAD_ENTRIES` never
copies the core's shipped `parley-addon.json` into a destination, and `lib/installer.js:2397`
runs `manifestProblems` only for `kind === "addon"`, so a managed core is validated on file
presence alone. The package even ships the manifest that would close this
(`skills/parley-deck/parley-addon.json`, 6 files, verified ok by `build-addon-manifest.js
--check`) — it is just never installed or consulted for a managed core.

why it matters: the core skill is the one that carries `references/COOPERATION.md` — the protocol
text every installed agent actually reads. A truncated extraction, a truncated write, or an
accidental edit of the installed protocol passes the health gate green, while the same defect in
any add-on is caught. `doctor` failing open on the highest-value file is exactly the defect class
this audit targets: the printed guarantee ("byte for byte") binds everywhere except the one place
most users actually have — a native, installer-managed core install.

### F2 — `sync-project --yes` silently deletes `protocolRole` from `meta/version.json`, the field §9.0 and the CLI's preflight both gate on — and `status` recommends that very command

severity: MAJOR
tag: PRIMARY (sandbox run + live file inspection); preflight cross-reference SECONDARY

command:
```
# sandbox project with a protocol-complete version.json containing "protocolRole": "source"
node $SKILL/bin/parley-deck-skill.js sync-project --project /tmp/pd-audit-kimi1/proj --yes
cat /tmp/pd-audit-kimi1/proj/parley-deck/meta/version.json
```

output: `wrote /tmp/pd-audit-kimi1/proj/parley-deck/meta/version.json`, exit 0. Before, the file
contained `"protocolRole": "source"`; after, the field is gone — `buildProjectMetadata`
(`lib/installer.js:596-613`) constructs a fresh object from 11 hardcoded keys and
`syncProjectCommand` overwrites the file wholesale. Live corroboration: this deck's own
`parley-deck/meta/version.json` (recorded `"updatedBy": "parley-deck-skill sync-project"`,
2026-08-12) contains **no** `protocolRole`, and `status --json` confirms
`metadata protocolRole present: False`.

contradicts: `parley-deck/COOPERATION.md:835-848` (§9.0 Protocol freshness — "Behaviour depends
on `meta/version.json` `protocolRole`"; "`protocolRole` missing/unknown → **do not auto-write**;
ask the user once and backfill the field"). The CLI enforces that field: `internal/app/preflight.go:384-416`
reads it and gates auto-write on it (`source` → "advisory only, never writes COOPERATION.md"),
`internal/protocol/workspace.go:73-84` writes `"protocolRole": "consumer"` at deck init. The skill
package contains zero references to `protocolRole` (grep over `parley-deck-skill/`: no hits), so
its writer drops the field every time.

why it matters: the role confirmation is a one-time user gate in the protocol. The live `status`
output on this very repo prints `action: Run parley-deck-skill sync-project --project
... --yes to refresh project metadata.` Following the tool's own recommended repair erases the
user's confirmed role; the next `parley preflight` then re-raises the unknown-role gate — or, with
`--yes`, backfills `protocolRole=consumer` (preflight.go:397-407), which on a **source** repo like
this one is the wrong role and flips §9.0 from "never auto-writes" onto the consumer path. Two
packages claim ownership of one file and write incompatible shapes; the message ("refresh project
metadata") misstates its effect ("also deletes protocolRole").

### F3 — README says "fourteen named runtimes" in four places and its `--target` enumeration omits `zcode`; the installer ships 15 targets and the changelog says 15

severity: MINOR
tag: PRIMARY

command:
```
node -e "console.log(require('$SKILL/lib/installer').TARGETS.length)"        # → 15
node $SKILL/bin/parley-deck-skill.js paths --target all --include-undetected | wc -l   # → 90 (15 targets × 6 skills)
grep -n "fourteen" $SKILL/README.md
```

output:
```
TARGETS count: 15   (codex,claude,agy,gemini,hermes,qwen,codebuddy,goose,kimi,droid,vibe,cursor,opencode,aionrs,zcode)
90
75:  ... measured across all fourteen destinations ...
161: This package's own installer covers fourteen named runtimes ...
206: Native targets are **fourteen named runtimes** — Codex, Claude Code, Antigravity CLI (plugin ...
262: `parley-deck-skill paths --target all --include-undetected` for all fourteen. ...
```
`README.md:237` enumerates `--target auto|all|codex|...|aionrs|generic` with no `zcode`.

contradicts: `lib/installer.js:18-128` (15 entries in `TARGETS`) and the package's own
`CHANGELOG.md` under 2.9.0: "The installer gained `zcode` as a target ... bringing the target
count to 15." The CLI's generated `usage()` text does include zcode, so only the README lagged.

why it matters: a user following the README cannot discover the `zcode` target at all, and the
"measured across all fourteen destinations" verification claim (line 75, attached to the bidding
runtime-availability statement) is now stale by construction — the measurement predates the 15th
destination it implicitly covers. Low blast radius, but it is the same printed-vs-actual drift
class, inside the package's own front door.

### F4 — On a source-role project, `status` recommends "adopting packaged protocol updates" — the direction §9.0 forbids — because the installer never reads `protocolRole`

severity: MINOR
tag: PRIMARY (live run on this repo)

command:
```
cd /Volumes/My Shared Files/AI_WORKSPACE/parley-deck/parley-deck-cli   # the protocol's upstream repo
node $SKILL/bin/parley-deck-skill.js status         # and: status --json
```

output (verbatim, live):
```
compatibility: warning
...
action: Review the local COOPERATION.md changes before adopting packaged protocol updates.
```
`status --json` reasons include `project-protocol-differs-from-packaged-reference`. Measured
hashes: live `parley-deck/COOPERATION.md` = `74c8470b…`, packaged copy = `254521eb…` — the live
protocol is the newer one (this repo is the protocol's source; the header records the 2026-08-19
sync to skill 2.9.0).

contradicts: `parley-deck/COOPERATION.md:839-840` — for `protocolRole: source`, the packaged copy
"is older, not newer" and the comparison is "advisory only". `recommendedActions`
(`lib/installer.js:580-594`) emits the same adopt-the-packaged-copy advice unconditionally; nothing
in the skill package reads `protocolRole` (F2), so on the one repository where the advice is
backwards by definition, it is emitted verbatim.

why it matters: the recommended action misstates the situation — following it on this repo means
treating the older packaged snapshot as an update to the newer live protocol. No auto-write
happens from `status` itself (and preflight's role gate currently stands in the way), which is why
this is MINOR; but it is a message that misstates its own effect, and its root cause is the same
missing `protocolRole` concept as F2.

### F5 — Installer comment declares "a marker at schema 2 that omits `manifest` is malformed", but the tool's own core marker is schema 2 with no `manifest` field

severity: NIT
tag: PRIMARY (markers inspected after sandbox install)

command:
```
cat /tmp/pd-audit-kimi1/home2/.gemini/extensions/parley-deck/.parley-deck-skill-install.json
```

output: `"markerSchema": 2`, `"skill": "parley-deck"`, `"addon": false`, an `addons` array — and
**no** `manifest` key. Add-on markers written in the same run do carry `manifest` (verified on
`parley-design`).

contradicts: `lib/installer.js:14-16` — "Schema 2 added `manifest`: a marker at this schema that
omits it is malformed, never treated as legacy." The rule is enforced only on add-on markers
(`manifestProblems` is reached solely via `kind === "addon"`, installer.js:2397), so the tool's
own core marker violates the printed invariant and nothing flags it.

why it matters: no runtime effect today — nothing consumes `manifest` from a core marker. But the
comment states a universal invariant that its own writer violates; anyone later enforcing the
comment literally would misflag every core install the tool has ever made at schema 2. Either the
comment or the core marker is wrong.

## What I checked and found clean

All PRIMARY (ran), in the sandbox or read-only live:

- **Addon manifests all verify**: `node scripts/build-addon-manifest.js --check` → all six skills
  `ok` (parley-bidding 47 files, parley-deck 6, parley-design 4, parley-design-check 126,
  parley-tracker 8, parley-worktrees 1), exit 0.
- **Live `doctor`** on this machine: detected 8 targets (codex, claude, agy, gemini, hermes, kimi,
  opencode, zcode), reported 42 valid units and zcode's 6 correctly `missing`, exit 1 with the
  right stderr summary. Detection of kimi at `~/.kimi-code/skills/parley-deck` matches reality.
- **Managed add-on tamper is caught** (byte flip → `integrity: modified: SKILL.md`, exit 1), and
  **foreign verbatim copies** of the core report `valid-unmanaged` and turn `malformed` on a byte
  flip — the README's `valid-unmanaged` paragraph is accurate for that case.
- **Runtime floor gate works**: with `python3` absent from PATH, `doctor --target agy` →
  `agy/parley-bidding: valid …  unavailable: python3 is not available, but this skill requires
  >=3.10`, exit 1 (manifest declares `python >=3.10`).
- **Antigravity staged shape**: installed plugin root contains `skills/SKILL.md` and
  `agents/manifest.yaml`, so the staged `plugin.json`'s references resolve. **Gemini staged
  manifest** is rewritten to `"contextFileName": "SKILL.md"` for the flat layout. (The repo-root
  `plugin.json` not resolving at repo root is *not* a defect — it is a repo-level manifest staged
  into the destination, and no repo-URL install path is documented for agy.)
- **`sync-project` without `--yes` is a dry run**: prints `would write …/meta/version.json`,
  writes nothing (meta dir stayed empty), exit 0 — matches the usage text.
- **`KIMI_CODE_HOME` redirect works**: `paths --target kimi` with the env var set points at
  `$KIMI_CODE_HOME/skills/parley-deck`.
- **Sandbox uninstall** removes all 6 units; subsequent `doctor` reports them `missing`, exit 1.
- **Version/compat metadata consistent**: `package.json` 2.9.0 = `compatibility.json`
  `skillVersion` 2.9.0 = `status` installer line; CLI `VERSION` 1.45.0 matches
  `COOPERATION.md:7`'s "parley-deck-cli 1.45.0".
- **`status --json` exposes the §9.0 comparison fields** (`project.protocolSha256`,
  `project.packaged.protocolSha256`), so the protocol-freshness comparison §9.0 names is
  mechanically possible.
- **Version-drift detection works**: 2.8.0 marker vs 2.9.0 installer → per-target
  `*-version-drift` reasons plus a reinstall action.
- Live `~/.zcode` exists (with `cli`, `v2`, `workspace`; no `skills` yet) — consistent with the
  zcode target's documented detection behaviour.

## What I could not check, and why

- **`agy plugin validate` end-to-end** — I did not run the Antigravity CLI (only the installer's
  staged file shape was verified, not the runtime's acceptance of it).
- **Whether the zcode runtime actually exposes `~/.zcode/skills`** — the installer comment claims
  verification against the runtime bundle dated 2026-08-19; I confirmed the dir layout exists but
  did not drive the zcode runtime.
- **Gemini CLI accepting the rewritten staged extension manifest** — README.md:257 itself flags
  this as not confirmed end to end; I did not run the Gemini CLI either.
- **The real install conventions of qwen, codebuddy, goose, droid, vibe, cursor, aionrs** — those
  runtimes are not present on this machine; only their detection/destination logic was exercised
  in the sandbox.
- **Windows/portable (`pkg`) binaries** — not exercised; all runs used `node` directly on macOS.
- **`parley preflight` end-to-end** (the CLI half of F2's interaction) — cited from source only
  (SECONDARY); running it against a copy was outside my lens's sandbox and is codex-1's/
  opencode-1's territory.
