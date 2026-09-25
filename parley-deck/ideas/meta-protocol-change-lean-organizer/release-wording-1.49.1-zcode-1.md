---
idea: meta-protocol-change-lean-organizer
author: zcode-1
role: participant (release-wording finalization; dispatched by organizer codex-1 under the owner's authorization — participant authorship: all edits in the paired commit are zcode-1's, committed under the release-resume's [codex-1] label per dispatch)
artifact: release-wording-1.49.1
date: 2026-09-24
status: wording-only update recorded; no release operations performed
basis: parley-deck/inbox/user-to-codex-1_meta-protocol-change-lean-organizer_windows-scope.md (owner decision, status authorized)
---

# CLI 1.49.1 owner-authorized release wording — exact changes — zcode-1

The owner decided **"Oboje" / "Both"** (2026-09-24 deck inbox, above): release
CLI 1.49.1 now for macOS and Linux through GitHub and Homebrew; keep the
Windows assets explicitly labelled experimental/unvalidated; hold the CLI
winget submission until a separate reviewed windows-portability track fixes
the defects, turns hosted windows-latest CI green, and removes the label.
The skill is unaffected; skill 2.13.0 npm remains with the owner-facing
Claude session and is not a CLI gate. This commit updates the **active**
wording of two owned artifacts to match. Wording only.

## Exact changes

### CHANGELOG.md — 1.49.1 section, "Release scope and limitations", bullet 1

- **Replaced** (active, pre-decision): "Per the organizer's scope escalation
  to the owner, the choice is between deferring the newly exposed
  native-Windows work … Publication is held pending that owner decision."
- **With** (owner authorization): the owner decided 2026-09-24 to do both —
  this release ships now for macOS and Linux through GitHub and Homebrew;
  Windows assets are kept and explicitly labelled experimental/unvalidated on
  the evidence above; the CLI winget submission is held until a separate,
  reviewed windows-portability track fixes the defects, turns hosted
  windows-latest CI green, and removes the label; until that track ships,
  every CLI release keeps this experimental Windows label and makes no CLI
  winget submission; the skill is unaffected.
- The rest of the bullet (the 1.49.0 false-claim correction and the hosted
  Windows failure facts) is unchanged.

### release-candidate-1.49.1-zcode-1.md

| location | change |
|---|---|
| front-matter | added dated `amended-2:` line recording the owner authorization and pointing to § 9 + this report; the original `amended:` line and `status:` untouched |
| intro | "release operations then wait on the required owner Windows scope decision and owner npm verification" → decision received 2026-09-24 (both: macOS/Linux via GitHub+Homebrew now; Windows assets kept labelled experimental/unvalidated; CLI winget held; separate reviewed windows-portability owns the repair); npm handled apart by the owner-facing Claude session, not a CLI gate; this commit performs no release operations |
| § 1 inputs | appended one bullet: the owner-decision inbox file with the decision summarized, the owner-facing release order (release-1.49.1 → meta-protocol-change-designated-implementer → windows-portability, done-file sequencing), and the skill/npm classification; existing npm-context bullet unchanged |
| § 3 Windows correction | tail "owner chooses … Publication is held pending that decision — a process hold…" → marked as the accurate-as-of-§ 8 original framing, then the recorded owner authorization mirroring the new changelog bullet |
| § 5 limits | Windows bullet: "owner choice … pending, and publication held on it; npm … also pending" → decision recorded (both) with its content; npm not a CLI gate |
| § 6 gates | heading dated "(as amended 2026-09-24, after the owner decision)"; gate 1: review marked done (kimi-1 PASS at 481fb65, recorded in 868825f), organizer asset staging still open; gate 2: **RESOLVED 2026-09-24** with the decision and release order; gate 3: **reclassified — not a CLI gate** (owner-attended, separate session, does not gate this release); gate 4: authorized scope — tag/builds/GitHub release/Homebrew for macOS/Linux, Windows assets retained under the label, **no CLI winget submission** until windows-portability ships, `v1.49.0`/`06e563e` never moves |
| § 7 NOT done | added: no CLI winget submission opened or prepared; no `release-1.49.1` done-file written — release operations and completion signaling remain organizer/owner-only |
| § 9 (new) | owner-authorization record: decision summary, what this commit changed, explicit statement that dated historical wording is preserved, pointer to this report |

## Preserved unchanged (dated historical statements)

- § 8 of the candidate report — still describes the `3002782`→`481fb65`
  record corrections as of their date, including "publication is held pending
  that owner decision", which was true then; § 9 records that the decision
  now lifts that hold.
- CHANGELOG 1.49.0 entry and everything below it — byte-identical; the 1.49.1
  edit touches only its own "Release scope and limitations" bullet.
- All peer and closed artifacts — kimi-1 candidate review, claude-1
  acceptance, organizer records (`organizer-notes.md`, `organizer-usage.md`,
  orchestration record), inbox, ledger, `runs/` — untouched; the organizer's
  untracked files (including `release-1.49.1-orchestration-codex-1.md`)
  remain untracked and unmodified.

## NOT done

No push, tag, build, publication, install, or CLI winget submission. No
product, test, workflow, or skill change. No behavior change of any kind. No
edit to peers, closed artifacts, or the owner/organizer inbox. All untracked
files preserved as found.

## Commit contents (this report's paired commit)

Exactly three paths: `CHANGELOG.md`,
`parley-deck/ideas/meta-protocol-change-lean-organizer/release-candidate-1.49.1-zcode-1.md`,
`parley-deck/ideas/meta-protocol-change-lean-organizer/release-wording-1.49.1-zcode-1.md`.
