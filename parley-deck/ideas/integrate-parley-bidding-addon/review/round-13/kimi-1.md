---
idea: integrate-parley-bidding-addon
review-round: 13
agent: kimi-1
date: 2026-07-31
reviewed-commit: dd8d756
---

## Verdict
ACCEPT

## Outstanding findings — closed or not

All three of my round-12 findings are closed, each re-measured at `dd8d756`:

- **MAJOR (install fleet regression).** Closed. The new `installFleetAtomically` stages the
  whole plan, commits by rename, and reverts every earlier commit on any failure. The shipped
  regression (`install is fleet-wide too: an immovable destination writes nothing anywhere`)
  passes; I additionally measured the sibling arm — `uchg` on a target's *skills directory*
  rather than on a destination — and got a clean preflight block: `ok:false`, 0 units written,
  an untouched target byte-identical, no staging debris.
- **MINOR (non-array `addons` container, shared with codex-1's MAJOR).** Closed. Only an
  absent field and explicit `false` read as core-only now; `"parley-bidding"`, `true`, `null`,
  `{}`, `42` all fail closed through health and both mutations (shipped regression passes; I
  re-measured `42` and `{bogus:true}` by hand).
- **NIT (discrimination table off by one).** Closed. `IMPLEMENTATION.md` cycle 16 now
  separates "proves the fix" from "guards the fix" and says so instead of reporting one number.

Round-12 codex-1's CRITICAL and both MAJORs are likewise covered by shipped regressions that I
re-ran and probed beyond their fixtures (see below).

## New findings

Two MINORs and two NITs. **None blocks the release**; all four are recorded follow-ups.

### MINOR — `install --dry-run` previews a green run that the real command refuses

`installCommand` routes dry runs to `installTarget`, which skips `preflightSkillUnit`
entirely when `dryRun` is set (`lib/installer.js:1357`). Everything that lives only in
preflight is therefore invisible to the preview. Measured at `dd8d756`, each with a real run
for contrast:

- Corrupt recorded selection (`addons: 42` in the core marker): `install --dry-run` returns
  **ok:true**, reporting `replace`/`install` for all six units. The real `install --force`
  blocks every unit with "install marker records an unusable add-on selection".
  `uninstall --dry-run` on the same tree surfaces the block (`blocked`, ok:false) — the two
  commands' previews disagree about the same file.
- Source manifest self-consistent but wrong about the payload: dry-run reports ok:true,
  `install` for all six; the real run fails preflight with "Source payload does not match
  parley-addon.json" and writes nothing.
- The boundary, also measured: the *ownership* gate does apply in dry-run (an unowned
  destination previews `blocked`, ok:false), because that check lives in `installSkillUnit`,
  not in preflight. So the blind set is exactly {recorded-selection problem, impossible
  destination ancestor, source-manifest mismatch}.

Nothing mutates — reality fails closed every time — but a preview whose entire job is "show
what would happen" reporting ok:true plus six mutations for a run that refuses everything is
wrong in the direction that matters, and `--json` automation previewing before executing reads
it as a green light. Fix shape: run the preflight pass in dry-run too and fold its blockers
into the preview, the way `removeFleetAtomically` already does for uninstall.

### MINOR — recorded-selection damage is a recovery dead-end, and no message names the exit

Measured: with `addons: {bogus: true}` in the core marker, `install --force` and
`uninstall --force` both refuse with "neither a list nor false" — the ratified fail-closed
edge, working as intended. But then no installer command can repair the state, and nothing
tells the user what can. The working repair — delete or fix the marker file, then
`install --force` — I verified succeeds (install ok:true), and it is undiscoverable from any
output. Contrast the neighbouring damage messages, which name their repair ("re-install to
repair it"); this class names none, and cycle 16 made the class much broader (any non-array
value lands here). Worth noting for whoever takes it: the install path never constructs paths
from the marker — its units come from discovery/flags — so the block on the install side buys
no path safety; letting install proceed would be self-healing. Exempting install or naming the
manual repair in the message are both acceptable resolutions; the message-only fix is the
smaller one.

### NIT — a recorded selection naming the core itself is authorized by the ownership clause

Measured: `addons: ["parley-deck"]` in the core marker passes `unusableAddonName`, is not
discovered, and is then authorized because the core destination carries our marker claiming
`skill: "parley-deck"` — the clause meant for add-ons dropped from newer packages. The result
is two units pointing at one destination: `doctor` reports the core `valid` and its twin
`malformed` for the same directory (health still fails, fail-closed), and
`uninstall --force` quarantines the core, hits ENOENT on the twin, and rolls back — a clean
block, nothing deleted, before/after listings identical. Harmless but nonsensical: a recorded
selection can never legitimately name the core skill. A one-line guard in the recorded-name
path (reject `CORE_SKILL_NAME`) would make the intent explicit.

### NIT — the install transaction's cleanup paths are not guarded the way uninstall's are

Inspection only; I could not build a deterministic arm. Two asymmetries with the uninstall
transaction hermes-1 attacked: (1) the Phase-1 and Phase-2 leftover `rmSync` calls and the
intra-commit backup restore in `commitStagedUnit` (`lib/installer.js:1746`) run unguarded —
a throw escapes the composed result and lands in the bin's catch-all (one line, exit 1, no
JSON document for `--json` consumers), where the revert loop one screen away is guarded and
reports the backup's location; (2) when the intra-commit restore is skipped or fails (dest
concurrently recreated; permissions flipped in the syscall gap), the composed message
"Destination could not be replaced (CODE); nothing was installed" does not name the surviving
`.bak` path holding the user's previous tree — the revert-failure message does. Both triggers
require mid-run external interference; every deterministic arm I tried (`uchg` destination,
`uchg` parent, pre-existing backup collision) fails closed with accurate messages.

## Release judgement

Releasable as 2.1.0. The cycle-16 claim — stored data is input, not truth — held under every
route I could construct from a file in a destination directory (or any other stored value) to
a filesystem path, a health verdict, or a mutation. The four findings above are follow-ups:
two preview/recovery polish items and two NITs.

## What I verified

Environment: clean scratch copy of the repo at `dd8d756` (`git status` clean before and
after), node v26.5.0, python3 3.14.6. The reviewed tree was never mutated; all fixtures ran
under `os.tmpdir()`.

- Full suite in the scratch copy: **344 node tests, 0 fail**; python leg **54/54 on 3.14**;
  manifest check ok — 47 files, aggregate
  `sha256:7854adf150712e0e3b9cca5618a23855024651670fdacc8392e1860568b95a6d`. Matches the
  facilitator's measurements.
- Stored-data inventory (the completeness attack). Every read of destination-resident data:
  core marker (`name`, `skill`, `addons`, `markerSchema`, `manifest`, `version`), add-on
  markers, installed `parley-addon.json` (keys, digests, aggregate, `runtime.python`),
  `plugin.json`/`gemini-extension.json` (existence only), project protocol/metadata (hash and
  JSON compare only). The only stored→path routes are `marker.addons` names and manifest
  keys; both are validated, confined, and (for marker names) authorized. `runtime.python` is
  regex-gated and anchored to the marker/source manifest hash before it can steer a verdict.
- Marker-name arms, measured: unknown sibling name under `--force` (clean block, sibling
  intact — the codex-1 CRITICAL regression re-run); dropped add-on with own marker
  (uninstallable — guard re-run); self-name `parley-deck` (NIT above); non-array containers
  `42` and `{bogus:true}` (fail closed on health and both mutations); case-variant names
  reasoning confirmed by code path (case-sensitive `discovered.has`, marker `skill` compare).
- Manifest-key arms, measured beyond the shipped fixture: `a\b`, `C:/x`, `a//b`, `a/`, `.`,
  a NUL-byte key, `CON`, and self-declaring keys (`parley-addon.json`, the install marker
  name). All refuse or fail closed; the shipped manifest itself passes.
- Install transaction, measured: `uchg` destination (shipped regression); `uchg` skills-dir
  parent (clean preflight block, zero writes, no debris); two targets sharing one skills dir
  via `CODEX_HOME == KIMI_CODE_HOME` (install ok, uninstall clean fleet rollback, nothing
  deleted, no debris); symlinked add-on destination (`uninstall --force` removes the link,
  the external tree survives); source-manifest mismatch (clean preflight block); dry-run
  boundaries (findings above).
- Uninstall transaction, measured: corrupt marker blocks with the problem named in dry-run
  and real runs; fleet gate intact across the shared-home config; rollback leaves the skills
  directory listing byte-identical.
- Deferred items re-confirmed as open and unchanged: only `parley-bidding` ships a manifest
  (`skills/*/parley-addon.json` — one hit); the `dirExists` discovery guard; debris invisible
  to `doctor` (I saw none created in any arm above, so the deferred item stayed theoretical
  here); `uappnd`/ACL residuals.
